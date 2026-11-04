#!/usr/bin/env bash
set -euo pipefail

audit_report=""
if [[ $# -ne 0 ]]; then
  if [[ $# -ne 2 || $1 != --external-audit ]]; then
    echo "Usage: $0 [--external-audit REPORT.json]" >&2
    exit 1
  fi
  audit_report=$2
fi

# MPRESS_BIN supports validation with a locally built candidate.
# Otherwise Cloudflare/CI download the pinned release, without Go or Node.js.
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
build_tmp=$(mktemp -d)
trap 'rm -rf "$build_tmp"' EXIT
mpress_bin=${MPRESS_BIN:-}
if [[ -z "$mpress_bin" ]]; then
  # Pin CI and publication to the same reviewed toolchain.
  release_tag=$(cat "$project_root/docs/mpress/mpress-version.txt")
  if [[ ! "$release_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid pinned M-Press version: $release_tag" >&2
    exit 1
  fi
  release_url="https://github.com/leaanthony/mpress/releases/download/$release_tag"
  archive="mpress-linux-amd64.tar.gz"
  echo "Downloading M-Press $release_tag"
  curl --fail --location --silent --show-error --retry 3 "$release_url/$archive" -o "$build_tmp/$archive"
  curl --fail --location --silent --show-error --retry 3 "$release_url/checksums.txt" -o "$build_tmp/checksums.txt"
  awk -v archive="$archive" '$2 == archive {print}' "$build_tmp/checksums.txt" > "$build_tmp/$archive.sha256"
  (cd "$build_tmp" && sha256sum --check --strict "$archive.sha256")
  tar -xzf "$build_tmp/$archive" -C "$build_tmp"
  test "$("$build_tmp/mpress" version)" = "mpress ${release_tag#v}"
  mpress_bin="$build_tmp/mpress"
fi
# A supplied candidate must be absolute because validation changes directory.
if [[ "$mpress_bin" != /* || ! -x "$mpress_bin" ]]; then
  echo "MPRESS_BIN must name an absolute executable path" >&2
  exit 1
fi
cd "$project_root"
python3 docs/mpress/scripts/check_translations.py --mpress "$mpress_bin"
# The imported home-page animation adds CSS classes at runtime.
"$mpress_bin" build --strict --no-purge-css
"$mpress_bin" check
python3 -m unittest discover -s docs/mpress/scripts -p 'test_*.py'
python3 docs/mpress/scripts/check_site.py docs/mpress/site

# Opt-in and independent of publication: report URLs that break after release
# branches disappear. Routine builds do not opt into these HTTP probes.
if [[ -n "$audit_report" ]]; then
  "$mpress_bin" check external --host github.com \
    --url-prefix https://github.com/wailsapp/wails/tree/ \
    --url-prefix https://github.com/wailsapp/wails/blob/ --json > "$audit_report"
fi
