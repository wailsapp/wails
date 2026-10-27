package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	versionFile          = "../../internal/version/version.txt"
	changelogFile        = "../../../docs/mpress/content/changelog.md"
	defaultReleaseBranch = "master"
	defaultReleaseTitle  = "Wails %s"
	defaultReleaseTarget = "master"
	githubDefaultAPI     = "https://api.github.com"
	githubAPIVersion     = "2022-11-28"
	runtimePackageJSON   = "v3/internal/runtime/desktop/@wailsio/runtime/package.json"
	runtimePackageLock   = "v3/internal/runtime/desktop/@wailsio/runtime/package-lock.json"
)

var (
	unreleasedChangelogFile = "../../UNRELEASED_CHANGELOG.md"
)

type releaseOptions struct {
	version string
	dryRun  bool
	branch  string
	target  string
	// channel selects the release cadence: "nightly" (patch bump), "weekly"
	// (minor bump, patch reset to 0) or "" / "auto" (derive from an explicit
	// version's shape, otherwise nightly).
	channel string
	// tagOnly cuts the release from the checked-out commit (typically a
	// detached, CI-verified master commit older than HEAD): only the version
	// metadata is committed, the branch is NOT pushed (it would not be a
	// fast-forward) and the changelog is left alone. The caller syncs the
	// version files back to the branch afterwards.
	tagOnly bool
}

var errNoUnreleasedContent = errors.New("No unreleased changelog content found.")

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

// getUnreleasedChangelogTemplate returns the template content for UNRELEASED_CHANGELOG.md
func getUnreleasedChangelogTemplate() string {
	return `# Unreleased Changes

<!--
AUTO-GENERATED FILE: do not edit in ordinary contribution PRs.

The Auto Changelog (v3) workflow generates entries after PRs merge.
Describe the user-visible change in the PR title and description instead.
Release tooling moves staged entries into docs/mpress/content/changelog.md
and resets this file. Do not delete or reset it by hand.

For changes to this automation or explicitly requested release-note corrections,
follow the changelog ownership policy in AGENTS.md and preserve existing entries.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
`
}

// clearUnreleasedChangelog clears the UNRELEASED_CHANGELOG.md file and resets it with the template
func clearUnreleasedChangelog() error {
	template := getUnreleasedChangelogTemplate()

	// Write the template back to the file
	err := os.WriteFile(unreleasedChangelogFile, []byte(template), 0o644)
	if err != nil {
		return fmt.Errorf("failed to reset UNRELEASED_CHANGELOG.md: %w", err)
	}

	fmt.Printf("Successfully reset %s with template content\n", unreleasedChangelogFile)
	return nil
}

// extractChangelogContent extracts the actual changelog content from UNRELEASED_CHANGELOG.md
// It returns the content between the section headers and the example section
func extractChangelogContent() (string, error) {
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", unreleasedChangelogFile, err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	var result []string
	var inExampleSection bool
	var inCommentBlock bool
	var hasActualContent bool
	var currentSection string

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Track comment blocks (handle multi-line comments)
		if strings.Contains(line, "<!--") {
			inCommentBlock = true
			// Check if comment ends on same line
			if strings.Contains(line, "-->") {
				inCommentBlock = false
			}
			continue
		}
		if inCommentBlock {
			if strings.Contains(line, "-->") {
				inCommentBlock = false
			}
			continue
		}

		// Skip the main title
		if strings.HasPrefix(trimmedLine, "# Unreleased Changes") {
			continue
		}

		// Check if we're entering the example section
		if strings.HasPrefix(trimmedLine, "---") || strings.HasPrefix(trimmedLine, "### Example Entries") {
			inExampleSection = true
			continue
		}

		// Skip example section content
		if inExampleSection {
			continue
		}

		// Handle section headers
		if strings.HasPrefix(trimmedLine, "##") {
			currentSection = trimmedLine
			// Only include section headers that have content after them
			// We'll add it later if we find content
			continue
		}

		// Handle bullet points
		if strings.HasPrefix(trimmedLine, "-") || strings.HasPrefix(trimmedLine, "*") {
			// Check if this is actual content (not empty)
			content := strings.TrimSpace(trimmedLine[1:])
			if content != "" {
				// If this is the first content in a section, add the section header first
				if currentSection != "" {
					// Only add empty line if this isn't the first section
					if len(result) > 0 {
						result = append(result, "")
					}
					result = append(result, currentSection)
					currentSection = "" // Reset so we don't add it again
				}
				result = append(result, line)
				hasActualContent = true
			}
		} else if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "<!--") {
			// Include other non-empty, non-comment lines that aren't section headers
			if !strings.HasPrefix(trimmedLine, "##") {
				// Check if next line exists and is not a comment placeholder
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if !strings.HasPrefix(nextLine, "<!--") {
						result = append(result, line)
					}
				}
			}
		}
	}

	if !hasActualContent {
		return "", nil
	}

	// Clean up result - remove any trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n"), nil
}

// hasUnreleasedContent checks if UNRELEASED_CHANGELOG.md has actual content beyond the template
func hasUnreleasedContent() (bool, error) {
	content, err := extractChangelogContent()
	if err != nil {
		return false, err
	}
	return content != "", nil
}

// safeFileOperation performs a file operation with backup and rollback capability
func safeFileOperation(filePath string, operation func() error) error {
	// Create backup if file exists
	var backupPath string
	var hasBackup bool

	if _, err := os.Stat(filePath); err == nil {
		backupPath = filePath + ".backup"
		if err := copyFile(filePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of %s: %w", filePath, err)
		}
		hasBackup = true
		defer func() {
			// Clean up backup file on success
			if hasBackup {
				_ = os.Remove(backupPath)
			}
		}()
	}

	// Perform the operation
	if err := operation(); err != nil {
		// Rollback if we have a backup
		if hasBackup {
			if rollbackErr := copyFile(backupPath, filePath); rollbackErr != nil {
				return fmt.Errorf("operation failed and rollback failed: %w (rollback error: %v)", err, rollbackErr)
			}
		}
		return err
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// updateVersion increments the version number properly handling semantic versioning
// Examples:
// v3.0.0-alpha.12 -> v3.0.0-alpha.13
// v3.0.0 -> v3.0.1
// v3.0.0-beta.1 -> v3.0.0-beta.2
func updateVersion() string {
	currentVersionData, err := os.ReadFile(versionFile)
	checkError(err)
	currentVersion := strings.TrimSpace(string(currentVersionData))
	newVersion := computeNextVersion(currentVersion)
	err = os.WriteFile(versionFile, []byte(newVersion), 0o755)
	checkError(err)
	return newVersion
}

// syncRuntimePackageVersion keeps the published runtime package metadata in
// lockstep with the CLI version. This runs for every release, even when no
// runtime source files changed, so a release tag never contains an older
// runtime version than the CLI it ships.
func syncRuntimePackageVersion(repoDir, releaseVersion string) error {
	version := strings.TrimPrefix(strings.TrimSpace(releaseVersion), "v")
	if version == "" {
		return errors.New("runtime package version cannot be empty")
	}

	files := []struct {
		path            string
		versionFields   int
		validatePackage bool
	}{
		{path: runtimePackageJSON, versionFields: 1, validatePackage: true},
		{path: runtimePackageLock, versionFields: 2, validatePackage: false},
	}

	versionField := regexp.MustCompile(`(?m)^([ \t]*"version"[ \t]*:[ \t]*)"[^"]*"([ \t]*,?[ \t]*)$`)
	for _, file := range files {
		path := filepath.Join(repoDir, file.path)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read runtime package metadata %s: %w", file.path, err)
		}

		matches := versionField.FindAllIndex(data, -1)
		if len(matches) < file.versionFields {
			return fmt.Errorf("runtime package metadata %s has %d version fields, want at least %d", file.path, len(matches), file.versionFields)
		}

		// The first version field is the package version. In package-lock.json,
		// the second is the root package entry under packages[""]; dependency
		// versions occur later and must remain untouched.
		replacements := 0
		updated := versionField.ReplaceAllFunc(data, func(field []byte) []byte {
			if replacements >= file.versionFields {
				return field
			}
			submatches := versionField.FindSubmatch(field)
			replacements++
			return []byte(string(submatches[1]) + `"` + version + `"` + string(submatches[2]))
		})

		var metadata struct {
			Version  string `json:"version"`
			Packages map[string]struct {
				Version string `json:"version"`
			} `json:"packages"`
		}
		if err := json.Unmarshal(updated, &metadata); err != nil {
			return fmt.Errorf("updated runtime package metadata %s is invalid JSON: %w", file.path, err)
		}
		if metadata.Version != version {
			return fmt.Errorf("runtime package metadata %s has version %q, want %q", file.path, metadata.Version, version)
		}
		if !file.validatePackage {
			root, ok := metadata.Packages[""]
			if !ok || root.Version != version {
				return fmt.Errorf("runtime package metadata %s root package has version %q, want %q", file.path, root.Version, version)
			}
		}

		if err := os.WriteFile(path, updated, 0o644); err != nil {
			return fmt.Errorf("failed to update runtime package metadata %s: %w", file.path, err)
		}
	}

	fmt.Printf("🔗 Synchronized @wailsio/runtime metadata to %s\n", version)
	return nil
}

// Release model (v3.0.0 and later). Tags are chosen so that Go's "@latest"
// (the highest tag WITHOUT a pre-release suffix) can only ever resolve to a
// weekly or hotfix release, never to a nightly:
//
//	vX.Y.0                    weekly   - minor advances with TIME, not features
//	vX.Y.Z (Z>0)              hotfix   - rare, manual, on top of a weekly
//	vX.(Y+1).0-nightly.DATE   nightly  - preview of the NEXT weekly (UTC date)
//	vX.Y.Z-alpha/beta/rc.N    legacy pre-GA builds
const (
	kindPrerelease = "prerelease"
	kindWeekly     = "weekly"
	kindHotfix     = "hotfix"
	kindNightly    = "nightly"
)

// releasePolicy is everything the pipeline derives from a version's shape.
type releasePolicy struct {
	Kind       string // prerelease | weekly | hotfix | nightly
	Prerelease bool   // GitHub release "prerelease" flag
	MakeLatest string // GitHub release "make_latest": "true" | "false"
	NPMTag     string // npm dist-tag for @wailsio/runtime
}

var versionShape = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(-[0-9A-Za-z.\-]+)?$`)

// nowUTC is replaceable in tests. Nightly tags use the UTC date.
var nowUTC = func() time.Time { return time.Now().UTC() }

// versionKind classifies a version by shape alone.
func versionKind(version string) (string, error) {
	m := versionShape.FindStringSubmatch(strings.TrimSpace(version))
	if m == nil {
		return "", fmt.Errorf("invalid semantic version %q", version)
	}
	switch {
	case strings.HasPrefix(m[4], "-nightly."):
		return kindNightly, nil
	case m[4] != "":
		return kindPrerelease, nil
	case m[3] == "0":
		return kindWeekly, nil
	}
	return kindHotfix, nil
}

// policyFor maps a version to GitHub release flags and the npm dist-tag.
//
//	weekly / hotfix -> normal release, make_latest=true,  npm "latest"
//	nightly         -> prerelease,     make_latest=false, npm "nightly"
//	legacy -beta.N  -> prerelease,     make_latest=false, npm "latest" (unchanged pre-GA behaviour)
func policyFor(version string) (releasePolicy, error) {
	kind, err := versionKind(version)
	if err != nil {
		return releasePolicy{}, err
	}
	switch kind {
	case kindWeekly, kindHotfix:
		return releasePolicy{Kind: kind, Prerelease: false, MakeLatest: "true", NPMTag: "latest"}, nil
	case kindNightly:
		return releasePolicy{Kind: kind, Prerelease: true, MakeLatest: "false", NPMTag: "nightly"}, nil
	}
	return releasePolicy{Kind: kind, Prerelease: true, MakeLatest: "false", NPMTag: "latest"}, nil
}

var releaseTagShape = regexp.MustCompile(`^v3\.\d+\.\d+(-(alpha\d*|beta|rc)\.\d+|-nightly\.\d{8}(\.\d+)?)?$`)

// listReleaseTags returns the v3 release tags (stray tags such as "-tui" are
// ignored). Returns nil if git is unavailable.
func listReleaseTags(repoDir string) []string {
	out, err := exec.Command("git", "-C", repoDir, "tag", "--list", "v3.*").Output()
	if err != nil {
		return nil
	}
	var tags []string
	for _, tag := range strings.Fields(string(out)) {
		if releaseTagShape.MatchString(tag) {
			tags = append(tags, tag)
		}
	}
	return tags
}

// highestOf returns the semver-highest of tags, optionally only those without
// a pre-release suffix ("what Go resolves @latest to").
func highestOf(tags []string, releasesOnly bool) string {
	best := ""
	for _, tag := range tags {
		if releasesOnly && semver.Prerelease(tag) != "" {
			continue
		}
		if best == "" || semver.Compare(tag, best) > 0 {
			best = tag
		}
	}
	return best
}

func highestReleaseTag(repoDir string) string { return highestOf(listReleaseTags(repoDir), true) }

// nextAutoVersion computes the version of an automatic (no explicit version)
// release from the existing tags. The base is the highest RELEASE tag
// (pre-releases ignored), so the result never depends on version.txt being in
// sync with the tags:
//
//	weekly  -> v<M>.<m+1>.0
//	nightly -> v<M>.<m+1>.0-nightly.<UTC yyyymmdd>[.N]
//
// Before any v3 release tag exists (the beta series) a nightly keeps counting
// the legacy pre-release (v3.0.0-beta.34 -> v3.0.0-beta.35).
func nextAutoVersion(versionTxt, channel string, tags []string, now time.Time) (string, error) {
	if channel == "" || channel == "auto" {
		channel = kindNightly
	}
	if channel != kindWeekly && channel != kindNightly {
		return "", fmt.Errorf("unknown channel %q (want auto, nightly or weekly)", channel)
	}
	base := highestOf(tags, true)
	if base == "" {
		kind, err := versionKind(versionTxt)
		if err != nil {
			return "", err
		}
		switch kind {
		case kindPrerelease:
			if channel == kindNightly {
				return computeNextVersion(versionTxt), nil
			}
			return "", fmt.Errorf("weekly release needs a release tag or release version.txt, got %q: pass an explicit --version (e.g. v3.0.0) to leave the pre-release series", versionTxt)
		case kindNightly:
			// version.txt already names the upcoming weekly; use its M.m.0 as target.
			m := versionShape.FindStringSubmatch(strings.TrimSpace(versionTxt))
			minor, _ := strconv.Atoi(m[2])
			base = fmt.Sprintf("v%s.%d.0", m[1], minor-1)
			if minor == 0 {
				base = "v" + m[1] + ".0.0"
			}
		default:
			base = "v" + strings.TrimPrefix(strings.TrimSpace(versionTxt), "v")
		}
	}
	m := versionShape.FindStringSubmatch(base)
	if m == nil {
		return "", fmt.Errorf("invalid base version %q", base)
	}
	minor, _ := strconv.Atoi(m[2])
	target := fmt.Sprintf("v%s.%d.0", m[1], minor+1)
	if channel == kindWeekly {
		return target, nil
	}
	tag := fmt.Sprintf("%s-nightly.%s", target, now.UTC().Format("20060102"))
	existing := map[string]bool{}
	for _, t := range tags {
		existing[t] = true
	}
	if !existing[tag] {
		return tag, nil
	}
	for n := 1; ; n++ {
		if c := fmt.Sprintf("%s.%d", tag, n); !existing[c] {
			return c, nil
		}
	}
}

// checkMonotonic refuses an automatically computed version that is not
// strictly greater (real semver ordering: pre-release < release) than every
// existing tag.
func checkMonotonic(newVersion, highest string) error {
	if highest == "" {
		return nil
	}
	if semver.Compare(newVersion, highest) <= 0 {
		return fmt.Errorf("computed version %s is not greater than existing tag %s: refusing to release", newVersion, highest)
	}
	return nil
}

func computeNextVersion(currentVersion string) string {
	if currentVersion == "" {
		return "v0.0.1"
	}

	if strings.Contains(currentVersion, "-") {
		parts := strings.SplitN(currentVersion, "-", 2)
		baseVersion := parts[0]
		preRelease := parts[1]
		lastDotIndex := strings.LastIndex(preRelease, ".")
		if lastDotIndex != -1 {
			preReleaseTag := preRelease[:lastDotIndex]
			numberStr := preRelease[lastDotIndex+1:]
			if number, err := strconv.Atoi(numberStr); err == nil {
				number++
				return fmt.Sprintf("%s-%s.%d", baseVersion, preReleaseTag, number)
			}
		}
		return computeNextVersion(baseVersion)
	}

	return incrementPatchVersion(currentVersion)
}

// incrementPatchVersion increments the patch version of a semantic version
// e.g., v3.0.0 -> v3.0.1
func incrementPatchVersion(version string) string {
	versionWithoutV := strings.TrimPrefix(version, "v")
	parts := strings.Split(versionWithoutV, ".")
	if len(parts) != 3 {
		fmt.Printf("Warning: Invalid semantic version format: %s\n", version)
		return version
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Printf("Warning: Could not parse patch version: %s\n", parts[2])
		return version
	}

	patch++
	return fmt.Sprintf("v%s.%s.%d", parts[0], parts[1], patch)
}

//func runCommand(name string, arg ...string) {
//	cmd := exec.Command(name, arg...)
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	err := cmd.Run()
//	checkError(err)
//}

//func IsPointRelease(currentVersion string, newVersion string) bool {
//	// The first n-1 parts of the version should be the same
//	if currentVersion[:len(currentVersion)-2] != newVersion[:len(newVersion)-2] {
//		return false
//	}
//	// split on the last dot in the string
//	currentVersionSplit := strings.Split(currentVersion, ".")
//	newVersionSplit := strings.Split(newVersion, ".")
//	// if the last part of the version is the same, it's a point release
//	currentMinor := lo.Must(strconv.Atoi(currentVersionSplit[len(currentVersionSplit)-1]))
//	newMinor := lo.Must(strconv.Atoi(newVersionSplit[len(newVersionSplit)-1]))
//	return newMinor == currentMinor+1
//}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--sync-runtime-version":
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "usage: --sync-runtime-version VERSION")
				os.Exit(1)
			}
			releaseDir, err := os.Getwd()
			if err == nil {
				err = syncRuntimePackageVersion(filepath.Clean(filepath.Join(releaseDir, "..", "..", "..")), args[1])
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		case "--print-policy":
			if len(args) < 2 {
				fmt.Println("usage: --print-policy VERSION")
				os.Exit(1)
			}
			pol, err := policyFor(args[1])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Printf("kind=%s\nprerelease=%t\nmake_latest=%s\nnpm_tag=%s\n", pol.Kind, pol.Prerelease, pol.MakeLatest, pol.NPMTag)
			os.Exit(0)
		case "--check-only":
			if err := handleCheckOnly(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		case "--extract-changelog":
			if err := handleExtractChangelog(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		case "--reset-changelog":
			if err := handleResetChangelog(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		case "--create-release-notes":
			releaseNotesPath := "../../release_notes.md"
			if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
				releaseNotesPath = args[1]
			}
			if err := handleCreateReleaseNotes(releaseNotesPath); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	opts, err := parseReleaseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if err := runRelease(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCheckOnly() error {
	changelogContent, err := extractChangelogContent()
	if err != nil {
		return fmt.Errorf("Error: Failed to extract unreleased changelog content: %w", err)
	}
	if strings.TrimSpace(changelogContent) == "" {
		return errNoUnreleasedContent
	}
	fmt.Println("Found unreleased changelog content.")
	return nil
}

func handleExtractChangelog() error {
	changelogContent, err := extractChangelogContent()
	if err != nil {
		return fmt.Errorf("Error: Failed to extract unreleased changelog content: %w", err)
	}
	if strings.TrimSpace(changelogContent) == "" {
		return errors.New("No changelog content found.")
	}
	fmt.Print(changelogContent)
	return nil
}

func handleResetChangelog() error {
	if err := clearUnreleasedChangelog(); err != nil {
		return fmt.Errorf("Error: Failed to reset changelog: %w", err)
	}
	return nil
}

func handleCreateReleaseNotes(targetPath string) error {
	changelogContent, err := extractChangelogContent()
	if err != nil {
		return fmt.Errorf("Error: Failed to extract unreleased changelog content: %w", err)
	}
	if strings.TrimSpace(changelogContent) == "" {
		return errors.New("Error: No changelog content found in UNRELEASED_CHANGELOG.md")
	}
	if err := os.WriteFile(targetPath, []byte(changelogContent), 0o644); err != nil {
		return fmt.Errorf("Error: Failed to write release notes to %s: %w", targetPath, err)
	}
	fmt.Printf("Successfully created release notes at %s\n", targetPath)
	return nil
}

func parseReleaseArgs(args []string) (releaseOptions, error) {
	fs := flag.NewFlagSet("release", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dryRun := fs.Bool("dry-run", false, "simulate the release without pushing changes or creating a GitHub release")
	branch := fs.String("branch", defaultReleaseBranch, "git branch to push release changes to")
	target := fs.String("target", defaultReleaseTarget, "target reference for the GitHub release (usually the same as branch)")
	channel := fs.String("channel", "auto", "release cadence: auto, nightly (patch bump) or weekly (minor bump, patch 0)")
	tagOnly := fs.Bool("tag-only", false, "commit version metadata, push only the tag (for releasing a detached, CI-verified commit); changelog untouched")
	versionFlag := fs.String("version", "", "explicit release version (overrides automatic increment)")

	if err := fs.Parse(args); err != nil {
		return releaseOptions{}, err
	}

	version := strings.TrimSpace(*versionFlag)
	remaining := fs.Args()
	if version == "" && len(remaining) > 0 {
		version = strings.TrimSpace(remaining[0])
		if len(remaining) > 1 {
			return releaseOptions{}, fmt.Errorf("unexpected argument: %s", strings.Join(remaining[1:], " "))
		}
	} else if version != "" && len(remaining) > 0 {
		return releaseOptions{}, fmt.Errorf("unexpected argument: %s", strings.Join(remaining, " "))
	}

	return releaseOptions{
		version: version,
		dryRun:  *dryRun,
		branch:  *branch,
		target:  *target,
		channel: strings.TrimSpace(*channel),
		tagOnly: *tagOnly,
	}, nil
}

func runRelease(opts releaseOptions) error {
	fmt.Println("🚀 Starting release workflow")

	releaseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine working directory: %w", err)
	}
	repoDir := filepath.Clean(filepath.Join(releaseDir, "..", "..", ".."))

	changelogContent, err := extractChangelogContent()
	if err != nil {
		return fmt.Errorf("failed to extract unreleased changelog content: %w", err)
	}
	changelogContent = strings.TrimSpace(changelogContent)
	if changelogContent == "" && strings.TrimSpace(opts.version) == "" && opts.channel != kindWeekly {
		fmt.Println("ℹ️  UNRELEASED_CHANGELOG.md has no unreleased entries. Skipping release.")
		writeGitHubOutput("release_skipped", "true")
		writeGitHubOutput("release_reason", "no_unreleased_changelog_content")
		writeGitHubOutput("release_outcome", "skipped")
		return nil
	}

	// Resolve the repository and validate the token BEFORE touching any
	// files: discovering a dead token after the version bump and changelog
	// reset have been written (and committed) would leave local side effects
	// even though nothing was pushed — both in CI and when run locally.
	// Dry runs never push, so they need no token at all.
	var repoSlug, token string
	if !opts.dryRun {
		repoSlug, err = resolveRepoSlug(repoDir)
		if err != nil {
			return fmt.Errorf("failed to determine repository slug: %w", err)
		}
		token, err = selectAPIToken(repoSlug)
		if err != nil {
			return err
		}
	}

	originalVersionData, err := os.ReadFile(versionFile)
	if err != nil {
		return fmt.Errorf("failed to read version file: %w", err)
	}
	originalVersion := strings.TrimSpace(string(originalVersionData))
	if originalVersion == "" {
		return errors.New("version file is empty")
	}
	fmt.Printf("📌 Current version: %s\n", originalVersion)

	newVersion := strings.TrimSpace(opts.version)
	if newVersion != "" {
		if !strings.HasPrefix(newVersion, "v") {
			newVersion = "v" + newVersion
		}
		if opts.dryRun {
			fmt.Printf("✏️  (simulated) version override: %s\n", newVersion)
		}
	} else {
		tags := listReleaseTags(repoDir)
		newVersion, err = nextAutoVersion(originalVersion, opts.channel, tags, nowUTC())
		if err != nil {
			return err
		}
		// Automatic versions must never go backwards or collide with a tag.
		if err := checkMonotonic(newVersion, highestOf(tags, false)); err != nil {
			return err
		}
		if opts.dryRun {
			fmt.Printf("🔢 (simulated) next version: %s\n", newVersion)
		}
	}
	policy, err := policyFor(newVersion)
	if err != nil {
		return err
	}
	fmt.Printf("🏷️  Release kind: %s (prerelease=%t, make_latest=%s, npm tag=%s)\n", policy.Kind, policy.Prerelease, policy.MakeLatest, policy.NPMTag)

	if opts.tagOnly || (policy.Kind == kindWeekly && changelogContent == "") {
		prev := highestReleaseTag(repoDir)
		if prev == "" {
			prev = "the previous release"
		}
		changelogContent = fmt.Sprintf("Weekly release: everything merged to master and verified by CI since %s. See the changelog for the individual entries.", prev)
	}

	if !opts.dryRun {
		if err := os.WriteFile(versionFile, []byte(newVersion), 0o755); err != nil {
			return fmt.Errorf("failed to write version file: %w", err)
		}
		fmt.Printf("🔢 Version set to: %s\n", newVersion)
	}

	if !opts.dryRun {
		if err := syncRuntimePackageVersion(repoDir, newVersion); err != nil {
			return err
		}
		if !opts.tagOnly {
			if err := applyChangelogUpdates(newVersion, changelogContent); err != nil {
				return err
			}
		}
	}

	releaseBody := buildReleaseBodyFor(newVersion, changelogContent, policy)
	writeGitHubOutput("release_kind", policy.Kind)
	writeGitHubOutput("release_npm_tag", policy.NPMTag)
	writeGitHubOutput("release_prerelease", strconv.FormatBool(policy.Prerelease))
	writeGitHubOutput("release_version", newVersion)
	writeGitHubOutput("release_tag", newVersion)
	writeGitHubOutput("release_target", opts.target)

	if opts.dryRun {
		writeGitHubOutput("release_dry_run", "true")
		writeGitHubOutput("release_outcome", "dry-run")
		fmt.Println("🧪 Dry run enabled: skipping git commit, push, tagging, and GitHub release creation")
		fmt.Println("\n--- Release Notes Preview ---")
		fmt.Println(releaseBody)
		return nil
	}
	writeGitHubOutput("release_dry_run", "false")

	git := newGitRunner(repoDir)
	if !opts.tagOnly {
		if err := git.ensureOnBranch(opts.branch); err != nil {
			return err
		}
	}

	filesToAdd := []string{
		"v3/internal/version/version.txt",
		runtimePackageJSON,
		runtimePackageLock,
	}
	if !opts.tagOnly {
		filesToAdd = append(filesToAdd, "docs/mpress/content/changelog.md", "v3/UNRELEASED_CHANGELOG.md")
	}
	if err := git.add(filesToAdd...); err != nil {
		return fmt.Errorf("failed to stage release files: %w", err)
	}

	hasChanges, err := git.hasStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to inspect staged changes: %w", err)
	}
	if !hasChanges {
		return errors.New("no changes were staged for commit")
	}

	// Do not add [skip ci] here: version.txt changes must trigger
	// publish-npm.yml so the matching runtime package is published.
	commitMessage := fmt.Sprintf("chore(v3): bump to %s and update changelog", newVersion)
	if err := git.commit(commitMessage); err != nil {
		return fmt.Errorf("failed to commit release changes: %w", err)
	}

	if !opts.tagOnly {
		if err := git.push(opts.branch, repoSlug, token); err != nil {
			return err
		}
	}

	if err := git.createAnnotatedTag(newVersion, fmt.Sprintf("Release %s", newVersion)); err != nil {
		return fmt.Errorf("failed to create git tag: %w", err)
	}

	if err := git.pushTag(newVersion, repoSlug, token); err != nil {
		return err
	}

	releaseTitle := fmt.Sprintf(defaultReleaseTitle, newVersion)
	releaseInfo, err := createGitHubRelease(token, repoSlug, opts.target, newVersion, releaseTitle, releaseBody, policy)
	if err != nil {
		return err
	}
	if releaseInfo.HTMLURL != "" {
		writeGitHubOutput("release_url", releaseInfo.HTMLURL)
	}

	writeGitHubOutput("release_outcome", "success")
	fmt.Println("🎉 Release completed successfully.")
	return nil
}

// selectAPIToken picks the token used for git pushes and the GitHub REST
// API, validating it against the API BEFORE anything is pushed. Discovering
// a dead token after the version bump and tag have been pushed leaves a
// half-released state (a tag with no GitHub release and a reset changelog),
// which is exactly what happened when WAILS_REPO_TOKEN expired: git
// operations still succeeded via the checkout step's persisted credentials,
// then the release API call failed with 401.
//
// If WAILS_REPO_TOKEN fails validation and a different GITHUB_TOKEN is
// available (the workflow passes github.token), it is used as a fallback so
// the nightly release still ships — with a warning, since github.token
// cannot trigger downstream workflows.
func selectAPIToken(repoSlug string) (string, error) {
	primary := strings.TrimSpace(os.Getenv("WAILS_REPO_TOKEN"))
	fallback := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if primary == "" {
		primary, fallback = fallback, ""
	}
	if primary == "" {
		return "", errors.New("WAILS_REPO_TOKEN (or GITHUB_TOKEN) must be set to push and create releases")
	}

	err := validateToken(primary, repoSlug)
	if err == nil {
		return primary, nil
	}
	if fallback != "" && fallback != primary {
		fmt.Printf("⚠️  WAILS_REPO_TOKEN failed API validation (%v) — falling back to GITHUB_TOKEN. The token is likely expired and should be rotated.\n", err)
		if fbErr := validateToken(fallback, repoSlug); fbErr == nil {
			return fallback, nil
		} else {
			return "", fmt.Errorf("WAILS_REPO_TOKEN failed validation (%v) and GITHUB_TOKEN fallback also failed: %w", err, fbErr)
		}
	}
	return "", fmt.Errorf("token failed API validation before any push (nothing was released): %w", err)
}

// validateToken performs a cheap authenticated read against the repository
// to prove the token is alive. An expired fine-grained PAT answers 401
// "Requires authentication" here, the same way the release-creation call
// would fail later.
func validateToken(token, repoSlug string) error {
	apiBase := strings.TrimSpace(os.Getenv("GITHUB_API_URL"))
	if apiBase == "" {
		apiBase = githubDefaultAPI
	}
	apiBase = strings.TrimSuffix(apiBase, "/")

	endpoint := fmt.Sprintf("%s/repos/%s", apiBase, repoSlug)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create token validation request: %w", err)
	}
	setGitHubHeaders(req, token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call GitHub API for token validation: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		return nil
	case resp.StatusCode == http.StatusUnauthorized:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API rejected the token (HTTP 401: %s) — expired or revoked?", strings.TrimSpace(string(body)))
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token validation returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
}

func applyChangelogUpdates(newVersion, changelogContent string) error {
	changelogData, err := os.ReadFile(changelogFile)
	if err != nil {
		return fmt.Errorf("failed to read changelog.md: %w", err)
	}
	changelog := string(changelogData)
	split := strings.Split(changelog, "## [Unreleased]")
	if len(split) != 2 {
		return fmt.Errorf("could not find '## [Unreleased]' section in changelog.md")
	}

	today := time.Now().Format("2006-01-02")
	newChangelog := split[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + "\n\n" + changelogContent + split[1]

	if err := safeFileOperation(changelogFile, func() error {
		return os.WriteFile(changelogFile, []byte(newChangelog), 0o644)
	}); err != nil {
		return fmt.Errorf("failed to update changelog.md: %w", err)
	}
	fmt.Println("📝 Updated docs changelog with new release entry.")

	if err := safeFileOperation(unreleasedChangelogFile, func() error {
		return clearUnreleasedChangelog()
	}); err != nil {
		return fmt.Errorf("failed to reset %s: %w", unreleasedChangelogFile, err)
	}
	fmt.Printf("🧹 Reset %s to template.\n", unreleasedChangelogFile)
	return nil
}

// releaseChannel derives the human-readable release channel from a version
// string, e.g. "v3.0.0-beta.3" -> "Beta". It returns an empty string for
// stable (non-prerelease) versions.
func releaseChannel(version string) string {
	switch {
	case strings.Contains(version, "-beta"):
		return "Beta"
	case strings.Contains(version, "-rc"):
		return "Release Candidate"
	case strings.Contains(version, "-alpha"):
		return "Alpha"
	default:
		return ""
	}
}

func buildReleaseBody(version, changelogContent string) string {
	pol, err := policyFor(version)
	if err != nil {
		pol = releasePolicy{Kind: kindNightly}
	}
	return buildReleaseBodyFor(version, changelogContent, pol)
}

func buildReleaseBodyFor(version, changelogContent string, pol releasePolicy) string {
	trimmed := strings.TrimSpace(changelogContent)
	channel := releaseChannel(version)
	heading := fmt.Sprintf("## Wails v3 Release - %s", version)
	var intro, footer string
	switch pol.Kind {
	case kindPrerelease:
		footer = "🤖 This is an automated pre-release generated from the latest changes on master."
		if channel != "" {
			heading = fmt.Sprintf("## Wails v3 %s Release - %s", channel, version)
		}
	case kindWeekly:
		heading = fmt.Sprintf("## Wails v3 Weekly Release - %s", version)
		footer = "🗓️ This is a **weekly** release. Minor version numbers advance with time, not with features: a new minor is not a feature milestone. Nightly previews are published between weekly releases as `-nightly` pre-releases."
	case kindHotfix:
		heading = fmt.Sprintf("## Wails v3 Hotfix Release - %s", version)
		footer = "🩹 This is a **hotfix** release on top of the previous weekly release."
	case kindNightly:
		heading = fmt.Sprintf("## Wails v3 Nightly - %s", version)
		intro = fmt.Sprintf("Nightly preview of the upcoming %s.", nightlyTarget(version))
		footer = "🤖 This is an automated **nightly** pre-release generated from the latest changes on master. It is not selected by `@latest`; install it by its exact tag."
	}
	sections := []string{heading, ""}
	if intro != "" {
		sections = append(sections, intro, "")
	}
	sections = append(sections,
		trimmed,
		"",
		"---",
		"",
		footer,
		"",
		"**Installation:**",
		"```bash",
		fmt.Sprintf("go install github.com/wailsapp/wails/v3/cmd/wails3@%s", version),
		"```",
		"",
		"**Tip:** always pin an exact version (e.g. `@"+version+"`) in `go.mod` and CI rather than `@latest`.",
	)
	switch channel {
	case "Alpha":
		sections = append(sections, "",
			"**⚠️ Alpha Warning:** This is pre-release software and may contain bugs or incomplete features.")
	case "Beta":
		sections = append(sections, "",
			"**⚠️ Beta Warning:** This is pre-release software. The API is stable, but you may still encounter issues before the final 3.0 release.")
	case "Release Candidate":
		sections = append(sections, "",
			"**⚠️ Release Candidate Warning:** This is pre-release software undergoing final validation before the final release.")
	}
	return strings.Join(sections, "\n")
}

// nightlyTarget returns the weekly a nightly previews: v3.2.0-nightly.20261102 -> v3.2.0.
func nightlyTarget(version string) string {
	if i := strings.Index(version, "-"); i >= 0 {
		return version[:i]
	}
	return version
}

type gitRunner struct {
	repoDir string
}

func newGitRunner(repoDir string) *gitRunner {
	return &gitRunner{repoDir: repoDir}
}

func (g *gitRunner) command(args ...string) *exec.Cmd {
	cmdArgs := append([]string{"-C", g.repoDir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd
}

func (g *gitRunner) run(args ...string) error {
	cmd := g.command(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *gitRunner) runCapture(args ...string) (string, error) {
	cmd := g.command(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func (g *gitRunner) add(files ...string) error {
	if len(files) == 0 {
		return nil
	}
	return g.run(append([]string{"add"}, files...)...)
}

func (g *gitRunner) hasStagedChanges() (bool, error) {
	cmd := g.command("diff", "--cached", "--quiet")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (g *gitRunner) commit(message string) error {
	return g.run("commit", "-m", message)
}

func (g *gitRunner) ensureOnBranch(expected string) error {
	branch, err := g.currentBranch()
	if err != nil {
		return err
	}
	if branch != expected && branch != "HEAD" {
		return fmt.Errorf("release must run from %s branch (current: %s)", expected, branch)
	}
	return nil
}

func (g *gitRunner) currentBranch() (string, error) {
	branch, err := g.runCapture("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to determine current branch: %w", err)
	}
	return branch, nil
}

func (g *gitRunner) push(branch, repoSlug, token string) error {
	fmt.Printf("⬆️  Pushing changes to %s\n", branch)
	authURL := buildAuthURL(repoSlug, token)
	if err := g.run("push", authURL, "HEAD:"+branch); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", branch, err)
	}
	return nil
}

func (g *gitRunner) createAnnotatedTag(tag, message string) error {
	fmt.Printf("🏷️  Creating tag %s\n", tag)
	return g.run("tag", "-a", "-f", tag, "-m", message)
}

func (g *gitRunner) pushTag(tag, repoSlug, token string) error {
	fmt.Printf("⬆️  Pushing tag %s\n", tag)
	authURL := buildAuthURL(repoSlug, token)
	if err := g.run("push", authURL, tag); err != nil {
		return fmt.Errorf("failed to push tag %s: %w", tag, err)
	}
	return nil
}

func resolveRepoSlug(repoDir string) (string, error) {
	if slug := strings.TrimSpace(os.Getenv("GITHUB_REPOSITORY")); slug != "" {
		return slug, nil
	}
	git := newGitRunner(repoDir)
	remote, err := git.runCapture("remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("failed to read origin remote: %w", err)
	}
	slug, err := parseGitRemote(remote)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func parseGitRemote(raw string) (string, error) {
	remote := strings.TrimSpace(raw)
	if remote == "" {
		return "", errors.New("origin remote URL is empty")
	}
	remote = strings.TrimSuffix(remote, ".git")
	if strings.HasPrefix(remote, "git@") {
		parts := strings.SplitN(remote, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("unsupported git remote format: %s", raw)
		}
		return parts[1], nil
	}

	if !strings.Contains(remote, "://") {
		remote = "https://" + remote
	}
	parsed, err := url.Parse(remote)
	if err != nil {
		return "", fmt.Errorf("failed to parse git remote: %w", err)
	}
	path := strings.TrimPrefix(parsed.Path, "/")
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return "", fmt.Errorf("could not extract repository from remote: %s", raw)
	}
	return path, nil
}

func buildAuthURL(repoSlug, token string) string {
	serverURL := strings.TrimSpace(os.Getenv("GITHUB_SERVER_URL"))
	if serverURL == "" {
		serverURL = "https://github.com"
	}
	parsed, err := url.Parse(serverURL)
	if err != nil || parsed.Host == "" {
		parsed = &url.URL{Scheme: "https", Host: "github.com"}
	}
	host := parsed.Host
	return fmt.Sprintf("https://x-access-token:%s@%s/%s.git", token, host, repoSlug)
}

type releaseCreateResponse struct {
	ID      int64  `json:"id"`
	HTMLURL string `json:"html_url"`
}

func createGitHubRelease(token, repoSlug, target, tag, title, body string, policy releasePolicy) (releaseCreateResponse, error) {
	apiBase := strings.TrimSpace(os.Getenv("GITHUB_API_URL"))
	if apiBase == "" {
		apiBase = githubDefaultAPI
	}
	apiBase = strings.TrimSuffix(apiBase, "/")

	payload := map[string]interface{}{
		"tag_name":         tag,
		"target_commitish": target,
		"name":             title,
		"body":             body,
		"draft":            false,
		"prerelease":       policy.Prerelease,
		"make_latest":      policy.MakeLatest,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to marshal GitHub release payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/repos/%s/releases", apiBase, repoSlug)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to create GitHub release request: %w", err)
	}
	setGitHubHeaders(req, token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to call GitHub release API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result releaseCreateResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return releaseCreateResponse{}, fmt.Errorf("failed to decode GitHub release response: %w", err)
		}
		if result.HTMLURL != "" {
			fmt.Printf("✅ GitHub release created: %s\n", result.HTMLURL)
		} else {
			fmt.Println("✅ GitHub release created successfully")
		}
		return result, nil
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		release, err := fetchReleaseByTag(apiBase, token, repoSlug, tag)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return releaseCreateResponse{}, fmt.Errorf("failed to create GitHub release (status %d): %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
		}
		if err := updateGitHubRelease(apiBase, token, repoSlug, release, title, body, policy); err != nil {
			return releaseCreateResponse{}, err
		}
		if release.HTMLURL != "" {
			fmt.Printf("♻️  Updated existing GitHub release: %s\n", release.HTMLURL)
		} else {
			fmt.Println("♻️  Updated existing GitHub release")
		}
		return release, nil
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	return releaseCreateResponse{}, fmt.Errorf("GitHub release API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
}

func fetchReleaseByTag(apiBase, token, repoSlug, tag string) (releaseCreateResponse, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/releases/tags/%s", apiBase, repoSlug, tag)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to create GitHub release lookup request: %w", err)
	}
	setGitHubHeaders(req, token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to query existing GitHub release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return releaseCreateResponse{}, fmt.Errorf("GitHub release lookup returned %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	var result releaseCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return releaseCreateResponse{}, fmt.Errorf("failed to decode release lookup response: %w", err)
	}
	return result, nil
}

func updateGitHubRelease(apiBase, token, repoSlug string, release releaseCreateResponse, title, body string, policy releasePolicy) error {
	if release.ID == 0 {
		return errors.New("cannot update GitHub release: missing release ID")
	}
	payload := map[string]interface{}{
		"name":        title,
		"body":        body,
		"prerelease":  policy.Prerelease,
		"make_latest": policy.MakeLatest,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal GitHub release update payload: %w", err)
	}
	endpoint := fmt.Sprintf("%s/repos/%s/releases/%d", apiBase, repoSlug, release.ID)
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create GitHub release update request: %w", err)
	}
	setGitHubHeaders(req, token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update existing GitHub release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub release update returned %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	return nil
}

func setGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	req.Header.Set("Content-Type", "application/json")
}

func writeGitHubOutput(key, value string) {
	outputPath := strings.TrimSpace(os.Getenv("GITHUB_OUTPUT"))
	if outputPath == "" || key == "" || value == "" {
		return
	}
	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Warning: unable to write %s to GITHUB_OUTPUT: %v\n", key, err)
		return
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "%s=%s\n", key, value); err != nil {
		fmt.Printf("Warning: unable to persist %s to GITHUB_OUTPUT: %v\n", key, err)
	}
}
