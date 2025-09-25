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
	"strconv"
	"strings"
	"time"
)

const (
	versionFile          = "../../internal/version/version.txt"
	changelogFile        = "../../../docs/src/content/docs/changelog.mdx"
	defaultReleaseBranch = "v3-alpha"
	defaultReleaseTitle  = "Wails %s"
	defaultReleaseTarget = "v3-alpha"
	githubDefaultAPI     = "https://api.github.com"
	githubAPIVersion     = "2022-11-28"
)

var (
	unreleasedChangelogFile = "../../UNRELEASED_CHANGELOG.md"
)

type releaseOptions struct {
	version string
	dryRun  bool
	branch  string
	target  string
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
This file is used to collect changelog entries for the next v3-alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
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

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new ` + "`SetWindowIcon()`" + ` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
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
	if changelogContent == "" {
		return errNoUnreleasedContent
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
		} else {
			if err := os.WriteFile(versionFile, []byte(newVersion), 0o755); err != nil {
				return fmt.Errorf("failed to write version file: %w", err)
			}
			fmt.Printf("✏️  Version override applied: %s\n", newVersion)
		}
	} else {
		if opts.dryRun {
			newVersion = computeNextVersion(originalVersion)
			fmt.Printf("🔢 (simulated) next version: %s\n", newVersion)
		} else {
			newVersion = updateVersion()
			fmt.Printf("🔢 Auto-incremented version to: %s\n", newVersion)
		}
	}

	if !opts.dryRun {
		if err := applyChangelogUpdates(newVersion, changelogContent); err != nil {
			return err
		}
	}

	releaseBody := buildReleaseBody(newVersion, changelogContent)
	writeGitHubOutput("release_version", newVersion)
	writeGitHubOutput("release_tag", newVersion)
	writeGitHubOutput("release_target", opts.target)

	if opts.dryRun {
		writeGitHubOutput("release_dry_run", "true")
		fmt.Println("🧪 Dry run enabled: skipping git commit, push, tagging, and GitHub release creation")
		fmt.Println("\n--- Release Notes Preview ---")
		fmt.Println(releaseBody)
		return nil
	}
	writeGitHubOutput("release_dry_run", "false")

	git := newGitRunner(repoDir)
	if err := git.ensureOnBranch(opts.branch); err != nil {
		return err
	}

	filesToAdd := []string{
		"v3/internal/version/version.txt",
		"docs/src/content/docs/changelog.mdx",
		"v3/UNRELEASED_CHANGELOG.md",
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

	commitMessage := fmt.Sprintf("chore(v3): bump to %s and update changelog [skip ci]", newVersion)
	if err := git.commit(commitMessage); err != nil {
		return fmt.Errorf("failed to commit release changes: %w", err)
	}

	repoSlug, err := resolveRepoSlug(repoDir)
	if err != nil {
		return fmt.Errorf("failed to determine repository slug: %w", err)
	}

	token := strings.TrimSpace(os.Getenv("WAILS_REPO_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}
	if token == "" {
		return errors.New("WAILS_REPO_TOKEN (or GITHUB_TOKEN) must be set to push and create releases")
	}

	if err := git.push(opts.branch, repoSlug, token); err != nil {
		return err
	}

	if err := git.createAnnotatedTag(newVersion, fmt.Sprintf("Release %s", newVersion)); err != nil {
		return fmt.Errorf("failed to create git tag: %w", err)
	}

	if err := git.pushTag(newVersion, repoSlug, token); err != nil {
		return err
	}

	releaseTitle := fmt.Sprintf(defaultReleaseTitle, newVersion)
	releaseInfo, err := createGitHubRelease(token, repoSlug, opts.target, newVersion, releaseTitle, releaseBody)
	if err != nil {
		return err
	}
	if releaseInfo.HTMLURL != "" {
		writeGitHubOutput("release_url", releaseInfo.HTMLURL)
	}

	fmt.Println("🎉 Release completed successfully.")
	return nil
}

func applyChangelogUpdates(newVersion, changelogContent string) error {
	changelogData, err := os.ReadFile(changelogFile)
	if err != nil {
		return fmt.Errorf("failed to read changelog.mdx: %w", err)
	}
	changelog := string(changelogData)
	split := strings.Split(changelog, "## [Unreleased]")
	if len(split) != 2 {
		return fmt.Errorf("could not find '## [Unreleased]' section in changelog.mdx")
	}

	today := time.Now().Format("2006-01-02")
	newChangelog := split[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + "\n\n" + changelogContent + split[1]

	if err := safeFileOperation(changelogFile, func() error {
		return os.WriteFile(changelogFile, []byte(newChangelog), 0o644)
	}); err != nil {
		return fmt.Errorf("failed to update changelog.mdx: %w", err)
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

func buildReleaseBody(version, changelogContent string) string {
	trimmed := strings.TrimSpace(changelogContent)
	sections := []string{
		fmt.Sprintf("## Wails v3 Alpha Release - %s", version),
		"",
		trimmed,
		"",
		"---",
		"",
		"🤖 This is an automated nightly release generated from the latest changes in the v3-alpha branch.",
		"",
		"**Installation:**",
		"```bash",
		fmt.Sprintf("go install github.com/wailsapp/wails/v3/cmd/wails3@%s", version),
		"```",
		"",
		"**⚠️ Alpha Warning:** This is pre-release software and may contain bugs or incomplete features.",
	}
	return strings.Join(sections, "\n")
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

func createGitHubRelease(token, repoSlug, target, tag, title, body string) (releaseCreateResponse, error) {
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
		"prerelease":       true,
		"make_latest":      "false",
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
		if err := updateGitHubRelease(apiBase, token, repoSlug, release, title, body); err != nil {
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

func updateGitHubRelease(apiBase, token, repoSlug string, release releaseCreateResponse, title, body string) error {
	if release.ID == 0 {
		return errors.New("cannot update GitHub release: missing release ID")
	}
	payload := map[string]interface{}{
		"name":        title,
		"body":        body,
		"prerelease":  true,
		"make_latest": "false",
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
