package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeNextVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: "v0.0.1"},
		{name: "alpha iteration", input: "v3.0.0-alpha.12", expected: "v3.0.0-alpha.13"},
		{name: "alpha without counter", input: "v3.0.0-alpha", expected: "v3.0.1"},
		{name: "without leading v", input: "3.2.5", expected: "v3.2.6"},
		{name: "invalid semver", input: "v3", expected: "v3"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := computeNextVersion(tc.input)
			if got != tc.expected {
				t.Fatalf("computeNextVersion(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestBuildReleaseBody(t *testing.T) {
	releaseContent := "## Added\n- Feature A\n\n## Fixed\n- Bug B"
	wantFragments := []string{
		"## Wails v3 Alpha Release - v3.0.0-alpha.42",
		"Feature A",
		"Bug B",
		"go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.42",
	}

	body := buildReleaseBody("v3.0.0-alpha.42", releaseContent)

	for _, fragment := range wantFragments {
		if !strings.Contains(body, fragment) {
			t.Fatalf("release body missing fragment %q\nbody: %s", fragment, body)
		}
	}
}

func TestParseReleaseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		expect    releaseOptions
	}{
		{
			name:   "defaults",
			args:   nil,
			expect: releaseOptions{version: "", dryRun: false, branch: defaultReleaseBranch, target: defaultReleaseTarget},
		},
		{
			name:   "dry run flag",
			args:   []string{"--dry-run"},
			expect: releaseOptions{version: "", dryRun: true, branch: defaultReleaseBranch, target: defaultReleaseTarget},
		},
		{
			name:   "explicit version",
			args:   []string{"--version", "v3.0.0-alpha.99"},
			expect: releaseOptions{version: "v3.0.0-alpha.99", dryRun: false, branch: defaultReleaseBranch, target: defaultReleaseTarget},
		},
		{
			name:   "positional version",
			args:   []string{"v3.0.0-alpha.99"},
			expect: releaseOptions{version: "v3.0.0-alpha.99", dryRun: false, branch: defaultReleaseBranch, target: defaultReleaseTarget},
		},
		{
			name:   "custom branch and target",
			args:   []string{"--branch", "release", "--target", "release"},
			expect: releaseOptions{version: "", dryRun: false, branch: "release", target: "release"},
		},
		{
			name:      "unexpected extra arg",
			args:      []string{"--version", "v3.0.0", "oops"},
			expectErr: true,
		},
		{
			name:      "missing flag value",
			args:      []string{"--branch"},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseReleaseArgs(tc.args)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil (opts=%+v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseReleaseArgs returned error: %v", err)
			}
			if got.version != tc.expect.version || got.dryRun != tc.expect.dryRun || got.branch != tc.expect.branch || got.target != tc.expect.target {
				t.Fatalf("parseReleaseArgs mismatch:\n got  %+v\n want %+v", got, tc.expect)
			}
		})
	}
}

func TestParseGitRemote(t *testing.T) {
	tests := []struct {
		name      string
		remote    string
		expected  string
		expectErr bool
	}{
		{name: "https", remote: "https://github.com/wailsapp/wails.git", expected: "wailsapp/wails"},
		{name: "ssh", remote: "git@github.com:wailsapp/wails.git", expected: "wailsapp/wails"},
		{name: "bare host", remote: "github.com/wailsapp/wails", expected: "wailsapp/wails"},
		{name: "with trailing slash", remote: "https://github.com/wailsapp/wails/", expected: "wailsapp/wails"},
		{name: "invalid", remote: "not a remote", expectErr: true},
		{name: "empty", remote: "", expectErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseGitRemote(tc.remote)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for %q but got none", tc.remote)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseGitRemote(%q) error: %v", tc.remote, err)
			}
			if got != tc.expected {
				t.Fatalf("parseGitRemote(%q) = %q, want %q", tc.remote, got, tc.expected)
			}
		})
	}
}

func TestBuildAuthURL(t *testing.T) {
	tests := []struct {
		name      string
		serverURL string
		expected  string
	}{
		{name: "default server", serverURL: "", expected: "https://x-access-token:TOKEN@github.com/org/repo.git"},
		{name: "custom server", serverURL: "https://ghe.example.com", expected: "https://x-access-token:TOKEN@ghe.example.com/org/repo.git"},
		{name: "invalid server", serverURL: ":://bad host", expected: "https://x-access-token:TOKEN@github.com/org/repo.git"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serverURL == "" {
				t.Setenv("GITHUB_SERVER_URL", "")
			} else {
				t.Setenv("GITHUB_SERVER_URL", tc.serverURL)
			}
			got := buildAuthURL("org/repo", "TOKEN")
			if got != tc.expected {
				t.Fatalf("buildAuthURL() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestWriteGitHubOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")
	if err := os.WriteFile(outputFile, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to prepare output file: %v", err)
	}

	t.Setenv("GITHUB_OUTPUT", outputFile)

	writeGitHubOutput("version", "v3.0.0-alpha.1")
	writeGitHubOutput("tag", "v3.0.0-alpha.1")

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	expected := []string{"version=v3.0.0-alpha.1", "tag=v3.0.0-alpha.1"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d (%v)", len(expected), len(lines), lines)
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Fatalf("line %d = %q, want %q", i, line, expected[i])
		}
	}
}
