package application

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// privateAPIFiles are the only files allowed to reference private macOS APIs.
var privateAPIFiles = map[string]bool{
	"mac_private_api_darwin.go":                   true,
	"mac_private_api_appstore_darwin.go":          true,
	"mac_private_api_devtools_darwin.go":          true,
	"mac_private_api_devtools_appstore_darwin.go": true,
	"mac_private_api_darwin.h":                    true,
	"mac_private_api_test.go":                     true,
}

// privateAPIReferences are undocumented selectors and key paths that must not
// appear anywhere else, so that `-tags appstore` produces a binary with no
// reference to them.
var privateAPIReferences = []string{
	`forKey:@"drawsBackground"`,
	`forKey:@"backgroundColor"`,
	`forKey:@"developerExtrasEnabled"`,
	`_inspector`,
	// Wails has its own "mobile features" files, so match the message send
	// rather than the bare selector name.
	`_features]`,
	`_setEnabled:forFeature:`,
	`setGroupIdentifier:`,
	`setGroupName:`,
	`setGroupSpacing:`,
}

func TestPrivateMacAPIsAreIsolated(t *testing.T) {
	t.Parallel()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || privateAPIFiles[name] {
			continue
		}
		// Only Apple platform sources can reference these APIs.
		if !strings.Contains(name, "darwin") && !strings.Contains(name, "ios") {
			continue
		}
		contents, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, reference := range privateAPIReferences {
			if strings.Contains(string(contents), reference) {
				t.Errorf("%s references the private macOS API %q: move it into mac_private_api_darwin.go and add a public or no-op counterpart to mac_private_api_appstore_darwin.go", name, reference)
			}
		}
	}
}

func TestPrivateMacAPIBuildVariantsAreExclusive(t *testing.T) {
	t.Parallel()

	darwinBuild := build.Default
	darwinBuild.GOOS = "darwin"
	darwinBuild.GOARCH = "arm64"
	darwinBuild.CgoEnabled = true

	appstoreBuild := darwinBuild
	appstoreBuild.BuildTags = []string{"appstore"}

	tests := []struct {
		file            string
		defaultBuild    bool
		appstoreBuild   bool
		devtoolsVariant bool
	}{
		{file: "mac_private_api_darwin.go", defaultBuild: true},
		{file: "mac_private_api_appstore_darwin.go", appstoreBuild: true},
		{file: "mac_private_api_devtools_darwin.go", defaultBuild: true, devtoolsVariant: true},
		{file: "mac_private_api_devtools_appstore_darwin.go", appstoreBuild: true, devtoolsVariant: true},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			got, err := darwinBuild.MatchFile(".", test.file)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.defaultBuild {
				t.Errorf("default darwin build includes %s = %v, want %v", test.file, got, test.defaultBuild)
			}

			got, err = appstoreBuild.MatchFile(".", test.file)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.appstoreBuild {
				t.Errorf("appstore darwin build includes %s = %v, want %v", test.file, got, test.appstoreBuild)
			}
		})
	}
}

// TestPrivateMacAPIVariantsImplementTheSameFunctions guards against adding a
// private implementation without its appstore counterpart, which would only
// fail at link time on a Mac.
func TestPrivateMacAPIVariantsImplementTheSameFunctions(t *testing.T) {
	t.Parallel()

	declared := declaredSeamFunctions(t, "mac_private_api_darwin.h")
	if len(declared) == 0 {
		t.Fatal("no wailsPrivate* functions declared in mac_private_api_darwin.h")
	}

	variants := [][2]string{
		{"mac_private_api_darwin.go", "mac_private_api_appstore_darwin.go"},
		{"mac_private_api_devtools_darwin.go", "mac_private_api_devtools_appstore_darwin.go"},
	}

	implemented := map[string]bool{}
	for _, pair := range variants {
		private := definedSeamFunctions(t, pair[0])
		appstore := definedSeamFunctions(t, pair[1])

		for _, name := range private {
			implemented[name] = true
		}
		if strings.Join(private, ",") != strings.Join(appstore, ",") {
			t.Errorf("%s defines %v but %s defines %v", pair[0], private, pair[1], appstore)
		}
	}

	for _, name := range declared {
		if !implemented[name] {
			t.Errorf("mac_private_api_darwin.h declares %s but no build variant implements it", name)
		}
	}
}

var (
	seamDeclaration = regexp.MustCompile(`(?m)^\s*(?:void|bool|int)\s+(wailsPrivate\w+)\s*\(`)
	seamDefinition  = regexp.MustCompile(`(?m)^(?:void|bool|int)\s+(wailsPrivate\w+)\s*\([^;]*\)\s*\{`)
)

func declaredSeamFunctions(t *testing.T, path string) []string {
	t.Helper()
	return matchSeamFunctions(t, path, seamDeclaration)
}

func definedSeamFunctions(t *testing.T, path string) []string {
	t.Helper()
	return matchSeamFunctions(t, path, seamDefinition)
}

func matchSeamFunctions(t *testing.T, path string, pattern *regexp.Regexp) []string {
	t.Helper()

	contents, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	seen := map[string]bool{}
	for _, match := range pattern.FindAllStringSubmatch(string(contents), -1) {
		if seen[match[1]] {
			continue
		}
		seen[match[1]] = true
		names = append(names, match[1])
	}
	sort.Strings(names)
	return names
}
