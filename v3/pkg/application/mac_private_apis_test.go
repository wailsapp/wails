package application

import (
	"go/build"
	"os"
	"strings"
	"testing"
)

func TestMacPrivateAPIBuildVariantsAreExclusive(t *testing.T) {
	t.Parallel()

	defaultBuild := build.Default
	defaultBuild.GOOS = "darwin"
	defaultBuild.CgoEnabled = true

	publicBuild := defaultBuild
	publicBuild.BuildTags = []string{"noprivateapis"}

	tests := []struct {
		name           string
		file           string
		defaultMatches bool
		publicMatches  bool
	}{
		{
			name:          "window public implementation",
			file:          "webview_window_darwin_public_apis.go",
			publicMatches: true,
		},
		{
			name:           "window private implementation",
			file:           "webview_window_darwin_private_apis.go",
			defaultMatches: true,
		},
		{
			name:          "devtools public implementation",
			file:          "webview_window_darwin_dev.go",
			publicMatches: true,
		},
		{
			name:           "devtools private implementation",
			file:           "webview_window_darwin_dev_private_apis.go",
			defaultMatches: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defaultMatches, err := defaultBuild.MatchFile(".", test.file)
			if err != nil {
				t.Fatal(err)
			}
			if defaultMatches != test.defaultMatches {
				t.Fatalf("default build match = %v, want %v", defaultMatches, test.defaultMatches)
			}

			publicMatches, err := publicBuild.MatchFile(".", test.file)
			if err != nil {
				t.Fatal(err)
			}
			if publicMatches != test.publicMatches {
				t.Fatalf("noprivateapis build match = %v, want %v", publicMatches, test.publicMatches)
			}
		})
	}
}

func TestPrivateMacAPIReferencesStayOutOfNoPrivateAPIFiles(t *testing.T) {
	t.Parallel()

	allowedFiles := map[string]bool{
		"webview_window_darwin_private_apis.go":     true,
		"webview_window_darwin_dev_private_apis.go": true,
	}
	privateReferences := []string{
		"_inspector",
		"developerExtrasEnabled",
		`forKey:@"drawsBackground"`,
		`forKey:@"backgroundColor"`,
		"setGroupIdentifier:",
		"setGroupName:",
		"setGroupSpacing:",
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "webview_window_darwin") || strings.HasSuffix(name, "_test.go") || allowedFiles[name] {
			continue
		}
		contents, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, privateReference := range privateReferences {
			if strings.Contains(string(contents), privateReference) {
				t.Errorf("%s contains private API reference %q", name, privateReference)
			}
		}
	}
}
