// Package changelog generates release-local notes from merged PRs in a tested Git range.
package changelog

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const Archive = "docs/mpress/content/changelog.md"

var sections = []string{"Added", "Changed", "Fixed", "Deprecated", "Removed", "Security"}
var versionPattern = regexp.MustCompile(`^v3\.\d+\.\d+(?:-(?:alpha\d*|beta|rc)\.\d+|-nightly\.\d{8}(?:\.\d+)?)?$`)
var conventional = regexp.MustCompile(`(?i)^([a-z]+)(?:\([^)]*\))?!?:\s*(.+)$`)
var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)
var authorPattern = regexp.MustCompile(`^[A-Za-z0-9-]+(?:\[bot\])?$`)
var releaseHeading = regexp.MustCompile(`(?im)^## Release notes[ \t]*\r?\n`)
var comments = regexp.MustCompile(`(?s)<!--.*?-->`)
var archiveHeading = regexp.MustCompile(`(?m)^## v3\.\d+\.\d+`)

type Entry struct {
	Number      int    `json:"number"`
	Section     string `json:"section"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	MergeCommit string `json:"merge_commit"`
}

type Receipt struct {
	NotesSHA256 string  `json:"notes_sha256"`
	Repository  string  `json:"repository"`
	Target      string  `json:"target"`
	Baseline    string  `json:"baseline"`
	BaselineRef string  `json:"baseline_ref"`
	Entries     []Entry `json:"entries"`
	SkippedPRs  []int   `json:"skipped_prs"`
	CommitCount int     `json:"commit_count"`
}

type Options struct {
	Repository string
	Target     string
	Channel    string
	Version    string
	Since      string
}

// API returns every item from a paginated GitHub list endpoint.
type API func(string) ([]json.RawMessage, error)
type pullRequest struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	MergedAt    string `json:"merged_at"`
	MergeCommit string `json:"merge_commit_sha"`
	Base        struct {
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"base"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

func git(root string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(output)), nil
}

func GitHubPages(endpoint string) ([]json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "gh", "api", "--method", "GET", "--paginate", "--slurp", endpoint, "-H", "Accept: application/vnd.github+json").Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, fmt.Errorf("GitHub %s: %w: %s", endpoint, err, bytes.TrimSpace(exitError.Stderr))
		}
		return nil, fmt.Errorf("GitHub %s: %w", endpoint, err)
	}
	var pages []json.RawMessage
	if err := json.Unmarshal(output, &pages); err != nil || pages == nil {
		return nil, fmt.Errorf("expected paginated lists from %s", endpoint)
	}
	items := []json.RawMessage{}
	for _, page := range pages {
		var batch []json.RawMessage
		if err := json.Unmarshal(page, &batch); err != nil || batch == nil {
			return nil, fmt.Errorf("expected paginated lists from %s", endpoint)
		}
		items = append(items, batch...)
	}
	return items, nil
}

func resolve(root, ref string) (string, error) {
	return git(root, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
}

func baseline(root, target string, stable bool) (string, error) {
	output, err := git(root, "tag", "--list", "v3.*")
	if err != nil {
		return "", err
	}
	args := []string{"describe", "--tags", "--abbrev=0"}
	hasEligibleTag := false
	for _, tag := range strings.Fields(output) {
		if versionPattern.MatchString(tag) && (!stable || !strings.Contains(tag, "-")) {
			args = append(args, "--match", tag)
			hasEligibleTag = true
		}
	}
	if !hasEligibleTag {
		return "", fmt.Errorf("no eligible v3 release tag; use --since for an explicit bootstrap baseline")
	}
	result, err := git(root, append(args, target)...)
	if err != nil {
		return "", fmt.Errorf("no reachable v3 release tag; use --since to bootstrap explicitly: %w", err)
	}
	return result, nil
}

func category(pr pullRequest) (string, string, error) {
	labels := map[string]bool{}
	for _, label := range pr.Labels {
		labels[strings.ToLower(label.Name)] = true
	}
	if labels["skip-changelog"] {
		return "", "", nil
	}
	kind, title := "", pr.Title
	if match := conventional.FindStringSubmatch(strings.TrimSpace(title)); match != nil {
		kind, title = strings.ToLower(match[1]), match[2]
	}
	section := ""
	for _, candidate := range sections {
		if labels["changelog:"+strings.ToLower(candidate)] {
			if section != "" {
				return "", "", fmt.Errorf("PR #%d has conflicting changelog categories", pr.Number)
			}
			section = candidate
		}
	}
	if section == "" {
		switch kind {
		case "build", "chore", "ci", "docs", "style", "test":
			return "", "", nil
		case "feat":
			section = "Added"
		case "fix":
			section = "Fixed"
		case "security":
			section = "Security"
		case "deprecate":
			section = "Deprecated"
		case "remove":
			section = "Removed"
		default:
			section = "Changed"
		}
	}
	if loc := releaseHeading.FindStringIndex(pr.Body); loc != nil {
		summary := pr.Body[loc[1]:]
		if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(summary); end != nil {
			summary = summary[:end[0]]
		}
		summary = strings.TrimSpace(comments.ReplaceAllString(summary, ""))
		if summary != "" {
			title = strings.TrimPrefix(strings.TrimPrefix(summary, "- "), "* ")
		}
	}
	title = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(strings.Join(strings.Fields(title), " "))
	if title == "" {
		return "", "", fmt.Errorf("PR #%d has no release-note title", pr.Number)
	}
	return section, title, nil
}

func collect(root, repo, start, target string, api API) ([]Entry, []int, []string, error) {
	output, err := git(root, "rev-list", "--reverse", start+".."+target)
	if err != nil {
		return nil, nil, nil, err
	}
	commits := strings.Fields(output)
	included := map[string]bool{}
	for _, sha := range commits {
		included[sha] = true
	}
	prs, err := api("repos/" + repo + "/pulls?state=closed&per_page=100")
	if err != nil {
		return nil, nil, nil, err
	}
	selected := map[int]pullRequest{}
	for _, raw := range prs {
		var pr pullRequest
		if err := json.Unmarshal(raw, &pr); err != nil {
			return nil, nil, nil, err
		}
		if pr.MergedAt != "" && included[pr.MergeCommit] && strings.EqualFold(pr.Base.Repo.FullName, repo) {
			selected[pr.Number] = pr
		}
	}
	numbers := []int{}
	for number := range selected {
		numbers = append(numbers, number)
	}
	sort.Ints(numbers)
	entries := []Entry{}
	skipped := []int{}
	for _, number := range numbers {
		pr := selected[number]
		section, title, err := category(pr)
		if err != nil {
			return nil, nil, nil, err
		}
		if section == "" {
			skipped = append(skipped, number)
			continue
		}
		files, err := api(fmt.Sprintf("repos/%s/pulls/%d/files?per_page=100", repo, number))
		if err != nil {
			return nil, nil, nil, err
		}
		if len(files) >= 3000 {
			return nil, nil, nil, fmt.Errorf("PR #%d reaches the GitHub changed-files limit", number)
		}
		relevant := false
		for _, raw := range files {
			var file struct {
				Filename string `json:"filename"`
				Previous string `json:"previous_filename"`
			}
			if err := json.Unmarshal(raw, &file); err != nil {
				return nil, nil, nil, err
			}
			for _, name := range []string{file.Filename, file.Previous} {
				if strings.HasPrefix(name, "v3/") || strings.HasPrefix(name, "docs/mpress/") {
					relevant = true
				}
			}
		}
		if !relevant {
			skipped = append(skipped, number)
			continue
		}
		if !authorPattern.MatchString(pr.User.Login) {
			return nil, nil, nil, fmt.Errorf("invalid author for PR #%d", number)
		}
		entries = append(entries, Entry{
			Number: number, Section: section, Title: title,
			Author: pr.User.Login, MergeCommit: pr.MergeCommit,
		})
	}
	return entries, skipped, commits, nil
}

func migratePending(root, repo string, commits []string, entries []Entry) ([]Entry, error) {
	data, err := os.ReadFile(filepath.Join(root, "v3/tasks/release/changelog-migration.json"))
	if os.IsNotExist(err) {
		return entries, nil
	}
	if err != nil {
		return nil, err
	}
	var seed struct {
		Repository string  `json:"repository"`
		Entries    []Entry `json:"entries"`
	}
	if err := json.Unmarshal(data, &seed); err != nil {
		return nil, err
	}
	if seed.Repository != repo {
		return entries, nil
	}
	included := map[string]bool{}
	for _, sha := range commits {
		included[sha] = true
	}
	current := map[int]Entry{}
	for _, entry := range entries {
		current[entry.Number] = entry
	}
	for _, entry := range seed.Entries {
		if included[entry.MergeCommit] {
			current[entry.Number] = entry
		}
	}
	result := []Entry{}
	for _, entry := range current {
		result = append(result, entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Number < result[j].Number })
	return result, nil
}

func render(repo string, entries []Entry) string {
	var notes strings.Builder
	notes.WriteString("# Unreleased Changes\n\n")
	for _, section := range sections {
		heading := false
		for _, entry := range entries {
			if entry.Section != section {
				continue
			}
			if !heading {
				fmt.Fprintf(&notes, "## %s\n\n", section)
				heading = true
			}
			fmt.Fprintf(&notes, "- %s in [PR](https://github.com/%s/pull/%d) by @%s\n", entry.Title, repo, entry.Number, entry.Author)
		}
		if heading {
			notes.WriteString("\n")
		}
	}
	return notes.String()
}

func Generate(root string, options Options, api API) (*Receipt, error) {
	for _, name := range []string{"release_changelog.md", "release_changelog.json"} {
		if err := os.Remove(filepath.Join(root, "v3", name)); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	if !repositoryPattern.MatchString(options.Repository) {
		return nil, fmt.Errorf("expected an explicit owner/repository")
	}
	if options.Target == "" {
		options.Target = "HEAD"
	}
	if options.Channel == "" {
		options.Channel = "nightly"
	}
	if options.Channel != "nightly" && options.Channel != "weekly" {
		return nil, fmt.Errorf("expected nightly or weekly channel")
	}
	target, err := resolve(root, options.Target)
	if err != nil {
		return nil, err
	}
	head, err := resolve(root, "HEAD")
	if err != nil {
		return nil, err
	}
	if target != head {
		return nil, fmt.Errorf("target must equal the checked-out, tested HEAD")
	}
	if options.Version != "" && !versionPattern.MatchString("v"+strings.TrimPrefix(options.Version, "v")) {
		return nil, fmt.Errorf("expected a v3 release version")
	}
	stable := options.Channel == "weekly"
	if options.Version != "" {
		stable = !strings.Contains(options.Version, "-")
	}
	startRef := options.Since
	if startRef == "" {
		startRef, err = baseline(root, target, stable)
		if err != nil {
			return nil, err
		}
	}
	start, err := resolve(root, startRef)
	if err != nil {
		return nil, err
	}
	if _, err = git(root, "merge-base", "--is-ancestor", start, target); err != nil {
		return nil, err
	}
	if api == nil {
		api = GitHubPages
	}
	entries, skipped, commits, err := collect(root, options.Repository, start, target, api)
	if err != nil {
		return nil, err
	}
	entries, err = migratePending(root, options.Repository, commits, entries)
	if err != nil {
		return nil, err
	}
	notes := render(options.Repository, entries)
	receipt := &Receipt{
		NotesSHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(notes))),
		Repository:  options.Repository, Target: target,
		Baseline: start, BaselineRef: startRef,
		Entries: entries, SkippedPRs: skipped, CommitCount: len(commits),
	}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return nil, err
	}
	receiptPath := filepath.Join(root, "v3/release_changelog.json")
	if err := os.WriteFile(receiptPath, append(data, '\n'), 0644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(root, "v3/release_changelog.md"), []byte(notes), 0644); err != nil {
		_ = os.Remove(receiptPath)
		return nil, err
	}
	return receipt, nil
}

func archiveSection(content, tag string) string {
	heading := regexp.MustCompile(`(?m)^## ` + regexp.QuoteMeta(tag) + ` - [^\n]+\n`)
	loc := heading.FindStringIndex(content)
	if loc == nil {
		return ""
	}
	end := len(content)
	if next := archiveHeading.FindStringIndex(content[loc[1]:]); next != nil {
		end = loc[1] + next[0]
	}
	return strings.TrimSpace(content[loc[0]:end])
}

func SyncArchive(root, tag string) error {
	if !versionPattern.MatchString(tag) {
		return fmt.Errorf("expected a v3 release tag")
	}
	tagged, err := git(root, "show", "refs/tags/"+tag+":"+Archive)
	if err != nil {
		return err
	}
	section := archiveSection(tagged+"\n", tag)
	if section == "" {
		return fmt.Errorf("tag %s has no release changelog section", tag)
	}
	path := filepath.Join(root, Archive)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	current := string(data)
	if existing := archiveSection(current, tag); existing != "" {
		if existing != section {
			return fmt.Errorf("existing changelog section for %s differs from the tag", tag)
		}
		return nil
	}
	const heading = "## [Unreleased]\n"
	if strings.Count(current, heading) != 1 {
		return fmt.Errorf("expected exactly one Unreleased heading")
	}
	return os.WriteFile(path, []byte(strings.Replace(current, heading, heading+"\n"+section+"\n", 1)), 0644)
}
