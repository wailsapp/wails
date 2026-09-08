#!/usr/bin/env bash

# Apply the repository's version and issue-type labels to GitHub issues and
# pull requests. The default mode processes every issue (excluding PRs).
#
# This intentionally uses the GitHub CLI rather than a third-party action so
# that the bulk backfill and the event-driven workflow use the same classifier.

set -Eeuo pipefail

REPO="${GH_REPO:-wailsapp/wails}"
NUMBER=""
DRY_RUN=0
INTERACTIVE=0
DEFAULT_VERSION="v3"

usage() {
    cat <<'EOF'
Usage: label-github-items.sh [options]

Label all issues in wailsapp/wails by default, or one issue/PR when --number
is provided. Existing labels are preserved. Pull requests are included when
--number points to one, which lets the GitHub workflow use this same script.

Options:
  --repo REPOSITORY       GitHub repository (default: wailsapp/wails)
  --number NUMBER         Process one issue or pull request
  --default-version v2|v3 Version used when the item has no clear version
                          signal (default: v3)
  --interactive           Ask for a version/type when the classifier is unsure
  --dry-run               Show changes without modifying GitHub
  -h, --help              Show this help

Requirements: gh, jq, and an authenticated gh session (or GH_TOKEN).
EOF
}

die() {
    echo "error: $*" >&2
    exit 1
}

while (($# > 0)); do
    case "$1" in
        --repo)
            (($# >= 2)) || die "--repo requires a value"
            REPO="$2"
            shift 2
            ;;
        --number)
            (($# >= 2)) || die "--number requires a value"
            NUMBER="$2"
            shift 2
            ;;
        --default-version)
            (($# >= 2)) || die "--default-version requires v2 or v3"
            DEFAULT_VERSION="$2"
            shift 2
            ;;
        --interactive)
            INTERACTIVE=1
            shift
            ;;
        --dry-run)
            DRY_RUN=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            die "unknown option: $1"
            ;;
    esac
done

[[ "$DEFAULT_VERSION" == "v2" || "$DEFAULT_VERSION" == "v3" ]] || \
    die "--default-version must be v2 or v3"
command -v gh >/dev/null || die "gh is required"
command -v jq >/dev/null || die "jq is required"

declare -A REPOSITORY_LABELS=()
while IFS= read -r label; do
    [[ -n "$label" ]] && REPOSITORY_LABELS["$label"]=1
done < <(gh label list --repo "$REPO" --limit 1000 --json name --jq '.[].name')

label_description() {
    case "$1" in
        v2) echo "Applies to Wails v2" ;;
        v3) echo "Applies to Wails v3" ;;
        Bug) echo "Something isn't working" ;;
        Enhancement) echo "New feature or request" ;;
        Security) echo "Security-related bug" ;;
        *) echo "Managed by the Wails issue labeler" ;;
    esac
}

label_color() {
    case "$1" in
        v2) echo "C5DEF5" ;;
        v3) echo "FBCA04" ;;
        Bug) echo "D73A4A" ;;
        Enhancement) echo "A2EEEF" ;;
        Security) echo "B60205" ;;
        *) echo "EDEDED" ;;
    esac
}

ensure_repository_label() {
    local label="$1"

    if [[ -n "${REPOSITORY_LABELS[$label]+x}" ]]; then
        return
    fi

    if ((DRY_RUN)); then
        echo "  would create repository label: $label"
    else
        gh label create "$label" --repo "$REPO" \
            --description "$(label_description "$label")" \
            --color "$(label_color "$label")" >/dev/null
    fi
    REPOSITORY_LABELS["$label"]=1
}

item_has_label() {
    local item="$1"
    local label="$2"
    jq -e --arg label "$label" \
        'any(.labels[]?; .name == $label)' <<<"$item" >/dev/null
}

item_has_any_label() {
    local item="$1"
    shift
    local label
    for label in "$@"; do
        if item_has_label "$item" "$label"; then
            return 0
        fi
    done
    return 1
}

classify_item() {
    local item="$1"
    local title body text files
    title=$(jq -r '.title // ""' <<<"$item")
    body=$(jq -r '.body // ""' <<<"$item")
    text="${title}
${body}"
    text="${text,,}"
    files=""

    # PRs can identify their version through their changed paths even when the
    # title/body does not mention one.
    if jq -e 'has("pull_request")' <<<"$item" >/dev/null; then
        files=$(gh api --paginate "repos/$REPO/pulls/$(jq -r '.number' <<<"$item")/files" \
            --jq '.[].filename')
    fi

    CLASSIFIED_VERSION=""
    if item_has_any_label "$item" "v2" "v2-only"; then
        CLASSIFIED_VERSION="v2"
    elif item_has_label "$item" "v3"; then
        CLASSIFIED_VERSION="v3"
    fi
    local has_v2=0
    local has_v3=0
    if [[ "$text" =~ (^|[^a-z0-9])v2([._:/-]|$) || \
          "$text" =~ version[[:space:]]*2 || "$text" =~ wails[[:space:]]*2 || \
          "$text" =~ master[[:space:]]+branch || "$files" =~ (^|[[:space:]])v2/ ]]; then
        has_v2=1
    fi
    if [[ "$text" =~ (^|[^a-z0-9])v3([._:/-]|$) || \
          "$text" =~ version[[:space:]]*3 || "$text" =~ wails[[:space:]]*3 || \
          "$files" =~ (^|[[:space:]])v3/ ]]; then
        has_v3=1
    fi

    if [[ -n "$CLASSIFIED_VERSION" ]]; then
        :
    elif ((has_v2 && !has_v3)); then
        CLASSIFIED_VERSION="v2"
    elif ((has_v3 && !has_v2)); then
        CLASSIFIED_VERSION="v3"
    elif ((INTERACTIVE)); then
        while [[ "$CLASSIFIED_VERSION" != "v2" && "$CLASSIFIED_VERSION" != "v3" ]]; do
            read -r -p "Version for #$(jq -r '.number' <<<"$item") [$DEFAULT_VERSION]: " answer
            answer="${answer:-$DEFAULT_VERSION}"
            if [[ "$answer" == "v2" || "$answer" == "v3" ]]; then
                CLASSIFIED_VERSION="$answer"
            else
                echo "  enter v2 or v3"
            fi
        done
    else
        CLASSIFIED_VERSION="$DEFAULT_VERSION"
    fi

    CLASSIFIED_TYPE=""
    if item_has_label "$item" "Bug"; then
        CLASSIFIED_TYPE="Bug"
    elif item_has_label "$item" "Enhancement"; then
        CLASSIFIED_TYPE="Enhancement"
    fi
    local has_bug=0
    local has_enhancement=0
    if [[ "$text" =~ (^|[^a-z0-9])(bug|crash|broken|failure|failed|fails|fix|fixes|panic|segfault|sigsegv|regression|incorrect|not[[:space:]]+working|doesn.t|does[[:space:]]+not|don.t|cannot|unable|leak|freeze|hang|wrong|missing|ignored|unexpected)([^a-z0-9]|$) ]]; then
        has_bug=1
    fi
    if [[ "$text" =~ (^|[^a-z0-9])(feature|enhancement|request|add|new|improve|functionality|support[[:space:]]+for|please[[:space:]]+add|would[[:space:]]+be[[:space:]]+nice|implement|allow)([^a-z0-9]|$) ]]; then
        has_enhancement=1
    fi

    if [[ -n "$CLASSIFIED_TYPE" ]]; then
        :
    elif ((has_bug)); then
        CLASSIFIED_TYPE="Bug"
    elif ((has_enhancement)); then
        CLASSIFIED_TYPE="Enhancement"
    elif ((INTERACTIVE)); then
        while [[ "$CLASSIFIED_TYPE" != "Bug" && "$CLASSIFIED_TYPE" != "Enhancement" ]]; do
            read -r -p "Type for #$(jq -r '.number' <<<"$item") [Enhancement]: " answer
            answer="${answer:-Enhancement}"
            case "${answer,,}" in
                bug) CLASSIFIED_TYPE="Bug" ;;
                enhancement|feature) CLASSIFIED_TYPE="Enhancement" ;;
                *) echo "  enter Bug or Enhancement" ;;
            esac
        done
    else
        # Every issue receives one of the two requested type labels. A lack of
        # defect language is treated as a request/change by default.
        CLASSIFIED_TYPE="Enhancement"
    fi

    CLASSIFIED_SECURITY=0
    if [[ "$text" =~ (^|[^a-z0-9])(security|vulnerability|exploit|cve|xss|csrf|injection|path[[:space:]]+traversal|directory[[:space:]]+traversal|arbitrary[[:space:]]+code[[:space:]]+execution|privilege[[:space:]]+escalation|authentication[[:space:]]+bypass|authorization[[:space:]]+bypass)([^a-z0-9]|$) ]]; then
        CLASSIFIED_SECURITY=1
    fi
}

apply_labels() {
    local item="$1"
    local number title
    number=$(jq -r '.number' <<<"$item")
    title=$(jq -r '.title // ""' <<<"$item")
    classify_item "$item"

    local -a labels_to_add=()
    # v2-only is a legacy alias already used by this repository's PR labeler.
    # Add the canonical v2 label while leaving the legacy label in place.
    if ! item_has_any_label "$item" "v2" "v3"; then
        labels_to_add+=("$CLASSIFIED_VERSION")
    fi
    if ! item_has_any_label "$item" "Bug" "Enhancement"; then
        labels_to_add+=("$CLASSIFIED_TYPE")
    fi
    if ((CLASSIFIED_SECURITY)) && [[ "$CLASSIFIED_TYPE" == "Bug" ]] && \
       ! item_has_label "$item" "Security"; then
        labels_to_add+=("Security")
    fi

    if ((${#labels_to_add[@]} == 0)); then
        VERBOSE_COUNT=$((VERBOSE_COUNT + 1))
        return
    fi

    echo "#$number: $title"
    echo "  classification: $CLASSIFIED_VERSION / $CLASSIFIED_TYPE$( ((CLASSIFIED_SECURITY)) && [[ "$CLASSIFIED_TYPE" == "Bug" ]] && echo ' / Security' || true )"

    local label
    for label in "${labels_to_add[@]}"; do
        ensure_repository_label "$label"
    done

    if ((DRY_RUN)); then
        echo "  would add: ${labels_to_add[*]}"
        CHANGED_COUNT=$((CHANGED_COUNT + 1))
        return
    fi

    local payload
    payload=$(printf '%s\n' "${labels_to_add[@]}" | jq -R . | jq -s '{labels: .}')
    gh api --method POST "repos/$REPO/issues/$number/labels" --input - <<<"$payload" >/dev/null
    echo "  added: ${labels_to_add[*]}"
    CHANGED_COUNT=$((CHANGED_COUNT + 1))
}

CHANGED_COUNT=0
VERBOSE_COUNT=0

if [[ -n "$NUMBER" ]]; then
    apply_labels "$(gh api "repos/$REPO/issues/$NUMBER")"
else
    echo "Scanning issues in $REPO..."
    while IFS= read -r item; do
        [[ -n "$item" ]] && apply_labels "$item"
    done < <(gh api --paginate "repos/$REPO/issues?state=all&per_page=100" \
        --jq '.[] | select(.pull_request == null) | @json')
fi

echo "Completed: $CHANGED_COUNT item(s) changed; $VERBOSE_COUNT already classified."
