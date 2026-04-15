//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const changelogPath = "v3/UNRELEASED_CHANGELOG.md"

var validSections = map[string]bool{
	"Added": true, "Changed": true, "Fixed": true,
	"Deprecated": true, "Removed": true, "Security": true,
}

func main() {
	prNumber := os.Getenv("PR_NUMBER")
	githubToken := os.Getenv("GITHUB_TOKEN")
	openrouterKey := os.Getenv("OPENROUTER_API_KEY")
	repo := os.Getenv("GITHUB_REPOSITORY")

	if prNumber == "" || githubToken == "" || openrouterKey == "" || repo == "" {
		fmt.Fprintln(os.Stderr, "❌ Required env vars: PR_NUMBER, GITHUB_TOKEN, OPENROUTER_API_KEY, GITHUB_REPOSITORY")
		os.Exit(1)
	}

	pr, err := fetchPR(repo, prNumber, githubToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not fetch PR info: %v\n", err)
		os.Exit(1)
	}

	context, err := fetchCodeRabbitSummary(repo, prNumber, githubToken)
	if err != nil {
		fmt.Printf("⚠️  Could not fetch CodeRabbit summary: %v — falling back to PR title\n", err)
		context = "PR Title: " + pr.Title
	}

	fmt.Printf("📝 Context length: %d chars\n", len(context))

	section, entry, err := generateEntry(context, openrouterKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ LLM error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Section: %s\n", section)
	fmt.Printf("✅ Entry:   %s\n", entry)

	prURL := fmt.Sprintf("https://github.com/%s/pull/%s", repo, prNumber)
	bullet := fmt.Sprintf("- %s in [PR](%s) by @%s", entry, prURL, pr.Author)

	if err := insertEntry(changelogPath, section, bullet); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to update changelog: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Changelog updated.")
}

func fetchCodeRabbitSummary(repo, prNumber, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments?per_page=100", repo, prNumber)
	body, err := githubGet(url, token)
	if err != nil {
		return "", err
	}

	var comments []struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Body string `json:"body"`
	}
	if err := json.Unmarshal(body, &comments); err != nil {
		return "", fmt.Errorf("parse comments: %w", err)
	}

	for _, c := range comments {
		if c.User.Login == "coderabbitai[bot]" && c.Body != "" {
			fmt.Println("✅ Found CodeRabbit summary")
			return c.Body, nil
		}
	}
	return "", fmt.Errorf("no CodeRabbit comment found")
}

type prInfo struct {
	Title  string
	Author string
}

func fetchPR(repo, prNumber, token string) (prInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", repo, prNumber)
	body, err := githubGet(url, token)
	if err != nil {
		return prInfo{}, err
	}
	var raw struct {
		Title string `json:"title"`
		User  struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return prInfo{}, fmt.Errorf("parse PR: %w", err)
	}
	return prInfo{Title: raw.Title, Author: raw.User.Login}, nil
}

func githubGet(url, token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API %s returned %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func generateEntry(context, apiKey string) (string, string, error) {
	prompt := `You are a changelog writer for Wails, a Go framework for building desktop apps with web frontends.

Given the following PR information, output ONLY a raw JSON object (no markdown, no code fences, no explanation) with exactly these two fields:
- "section": one of: Added, Changed, Fixed, Deprecated, Removed, Security
- "entry": a concise description of the change (max 100 chars, no leading dash, no trailing period, present tense, no PR links or usernames)

` + context

	payload, _ := json.Marshal(map[string]any{
		"model":       "google/gemini-2.0-flash-lite-001",
		"temperature": 0.1,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("OpenRouter %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil || len(result.Choices) == 0 {
		return "", "", fmt.Errorf("parse response: %w body: %s", err, body)
	}

	raw := result.Choices[0].Message.Content
	fmt.Printf("🤖 Raw LLM response: %s\n", raw)

	// Extract JSON from response — handles code fences and extra text
	jsonRe := regexp.MustCompile(`\{[^{}]+\}`)
	match := jsonRe.FindString(raw)
	if match == "" {
		return "", "", fmt.Errorf("no JSON object found in response: %s", raw)
	}

	var out struct {
		Section string `json:"section"`
		Entry   string `json:"entry"`
	}
	if err := json.Unmarshal([]byte(match), &out); err != nil {
		return "", "", fmt.Errorf("parse LLM JSON %q: %w", match, err)
	}

	if !validSections[out.Section] {
		fmt.Printf("⚠️  Unknown section %q — defaulting to Changed\n", out.Section)
		out.Section = "Changed"
	}
	if out.Entry == "" {
		return "", "", fmt.Errorf("LLM returned empty entry")
	}

	return out.Section, out.Entry, nil
}

func insertEntry(path, section, bullet string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)

	// Insert bullet after "## <Section>" heading (optionally followed by an HTML comment)
	headingRe := regexp.MustCompile(`(?m)(^## ` + regexp.QuoteMeta(section) + `\n(?:<!--[^\n]*-->\n)*)`)
	if headingRe.MatchString(content) {
		content = headingRe.ReplaceAllString(content, "$1"+bullet+"\n")
	} else {
		// Section not present — add it before the --- footer
		content = strings.Replace(content, "\n---\n", "\n## "+section+"\n"+bullet+"\n\n---\n", 1)
	}

	return os.WriteFile(path, []byte(content), 0644)
}
