package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	unreleasedChangelogFile = "../../UNRELEASED_CHANGELOG.md"
	versionPrefix           = "v3.0.0-alpha."
)

var (
	checkOnly = flag.Bool("check-only", false, "Only check if there's unreleased content, exit 0 if yes, 1 if no")
	dryRun    = flag.Bool("dry-run", false, "Run in dry-run mode (no actual release)")
)

func main() {
	flag.Parse()

	// Check for unreleased content
	hasContent, entries := checkUnreleasedContent()

	if *checkOnly {
		if hasContent {
			fmt.Println("Found unreleased changelog content")
			os.Exit(0)
		} else {
			fmt.Println("No unreleased changelog content found")
			os.Exit(1)
		}
	}

	if !hasContent {
		fmt.Println("No unreleased changelog content found. Nothing to release.")
		os.Exit(0)
	}

	// Determine the next version
	nextVersion := getNextVersion()
	releaseTag := versionPrefix + strconv.Itoa(nextVersion)

	fmt.Printf("Preparing release: %s\n", releaseTag)

	if *dryRun {
		fmt.Println("[DRY-RUN] Would perform the following actions:")
		fmt.Printf("[DRY-RUN] - Create release tag: %s\n", releaseTag)
		fmt.Printf("[DRY-RUN] - Changelog entries to include:\n")
		for _, entry := range entries {
			fmt.Printf("[DRY-RUN]   %s\n", entry)
		}
		setGitHubOutput("release_version", releaseTag)
		setGitHubOutput("release_tag", releaseTag)
		setGitHubOutput("release_dry_run", "true")
		fmt.Println("[DRY-RUN] Release simulation complete")
		return
	}

	// Perform actual release
	if err := performRelease(releaseTag, entries); err != nil {
		fmt.Printf("Release failed: %v\n", err)
		os.Exit(1)
	}

	// Set GitHub outputs
	setGitHubOutput("release_version", releaseTag)
	setGitHubOutput("release_tag", releaseTag)
	setGitHubOutput("release_dry_run", "false")

	fmt.Printf("Release %s completed successfully!\n", releaseTag)
}

func checkUnreleasedContent() (bool, []string) {
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		fmt.Printf("Warning: Could not read %s: %v\n", unreleasedChangelogFile, err)
		return false, nil
	}

	lines := strings.Split(string(content), "\n")
	var entries []string

	// Look for bullet point entries only after the separator line (---)
	// This avoids matching the template category descriptions
	foundSeparator := false
	bulletRegex := regexp.MustCompile(`^\s*-\s+[^\s]`)

	for _, line := range lines {
		// Look for the separator line
		if strings.TrimSpace(line) == "---" {
			foundSeparator = true
			continue
		}

		// Only match entries after the separator
		if foundSeparator && bulletRegex.MatchString(line) {
			entries = append(entries, strings.TrimSpace(line))
		}
	}

	return len(entries) > 0, entries
}

func getNextVersion() int {
	// Get all existing v3.0.0-alpha.* tags
	cmd := exec.Command("git", "tag", "--list", "v3.0.0-alpha.*")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("No existing alpha tags found, starting with alpha.1")
		return 1
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return 1
	}

	// Extract version numbers and find the highest
	var versions []int
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		// Extract the number after "v3.0.0-alpha."
		parts := strings.Split(tag, "v3.0.0-alpha.")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err == nil {
				versions = append(versions, num)
			}
		}
	}

	if len(versions) == 0 {
		return 1
	}

	sort.Ints(versions)
	return versions[len(versions)-1] + 1
}

func performRelease(releaseTag string, entries []string) error {
	// Reset the unreleased changelog
	if err := resetUnreleasedChangelog(); err != nil {
		return fmt.Errorf("failed to reset unreleased changelog: %w", err)
	}

	// Stage the changes
	if err := runCommand("git", "add", unreleasedChangelogFile); err != nil {
		return fmt.Errorf("failed to stage changelog: %w", err)
	}

	// Create commit
	commitMsg := fmt.Sprintf("chore: release %s", releaseTag)
	if err := runCommand("git", "commit", "-m", commitMsg); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Create annotated tag
	tagMsg := fmt.Sprintf("Release %s\n\nChanges:\n%s", releaseTag, strings.Join(entries, "\n"))
	if err := runCommand("git", "tag", "-a", releaseTag, "-m", tagMsg); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Push commit and tag
	if err := runCommand("git", "push", "origin", "HEAD"); err != nil {
		return fmt.Errorf("failed to push commit: %w", err)
	}

	if err := runCommand("git", "push", "origin", releaseTag); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	return nil
}

func resetUnreleasedChangelog() error {
	template := `# Unreleased Changelog

All notable changes to the v3 alpha will be documented in this file before release.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## Categories

Use the following categories to organize your entries:
- ` + "`Added`" + ` for new features.
- ` + "`Changed`" + ` for changes in existing functionality.
- ` + "`Deprecated`" + ` for soon-to-be removed features.
- ` + "`Removed`" + ` for now removed features.
- ` + "`Fixed`" + ` for any bug fixes.
- ` + "`Security`" + ` in case of vulnerabilities.

---

<!-- Add your changelog entries below this line -->

`
	return os.WriteFile(unreleasedChangelogFile, []byte(template), 0644)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func setGitHubOutput(name, value string) {
	// Write to GITHUB_OUTPUT file if available (GitHub Actions)
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile != "" {
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			defer f.Close()
			writer := bufio.NewWriter(f)
			fmt.Fprintf(writer, "%s=%s\n", name, value)
			writer.Flush()
		}
	}
	// Also print for visibility
	fmt.Printf("::set-output name=%s::%s\n", name, value)
}

// getToday returns today's date in YYYY-MM-DD format (used for release notes)
func getToday() string {
	return time.Now().Format("2006-01-02")
}

// getReleaseDir returns the directory where release artifacts should be stored
func getReleaseDir() string {
	dir, err := filepath.Abs(".")
	if err != nil {
		return "."
	}
	return dir
}
