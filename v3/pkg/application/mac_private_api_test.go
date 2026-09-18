package application

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// Only this translation unit may contain private macOS selectors and keys.
const privateMacAPIFile = "mac_private_api_darwin.go"

var privateMacAPIReferences = []*regexp.Regexp{
	regexp.MustCompile(`@"(?:drawsBackground|backgroundColor|developerExtrasEnabled|groupSpacing|groupIdentifier|groupName)"`),
	regexp.MustCompile(`\b(?:_WKInspector|_inspector|_features|_setEnabled)\b`),
	regexp.MustCompile(`\bsetGroup(?:Identifier|Name|Spacing)\s*:`),
}

func TestPrivateMacAPIsAreIsolated(t *testing.T) {
	t.Parallel()
	// Scan all v3 native and cgo sources, including helpers outside application.
	err := filepath.WalkDir("../..", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".m" && ext != ".h" && ext != ".mm" {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") || strings.Contains(filepath.ToSlash(path), "/testdata/") {
			return nil
		}
		if filepath.Clean(path) == filepath.Join("..", "..", "pkg", "application", privateMacAPIFile) {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, pattern := range privateMacAPIReferences {
			if match := pattern.Find(contents); match != nil {
				t.Errorf("%s references private macOS API %q; move it into %s", path, match, privateMacAPIFile)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrivateMacAPIBuildVariantsAreExclusive(t *testing.T) {
	t.Parallel()
	for _, platform := range []string{"darwin", "linux", "windows", "ios"} {
		for _, arch := range []string{"arm64", "amd64"} {
			for _, mode := range []string{"", "production", "devtools", "production,devtools", "server", "ios", "appstore", "noprivateapis"} {
				for _, private := range []bool{false, true} {
					tags := strings.FieldsFunc(mode, func(r rune) bool { return r == ',' })
					if private {
						tags = append(tags, "private_mac_apis")
					}
					t.Run(platform+"/"+arch+"/"+strings.Join(tags, ","), func(t *testing.T) {
						ctx := build.Default
						ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled, ctx.BuildTags = platform, arch, true, tags
						desktop := platform == "darwin" && !slices.Contains(tags, "ios") && !slices.Contains(tags, "server")
						for file, want := range map[string]bool{
							privateMacAPIFile:              desktop && private,
							"mac_public_api_darwin.go":     desktop && !private,
							"webview_window_darwin_dev.go": desktop && (!slices.Contains(tags, "production") || slices.Contains(tags, "devtools")),
						} {
							got, err := ctx.MatchFile(".", file)
							if err != nil {
								t.Fatal(err)
							}
							if got != want {
								t.Errorf("%s selected = %v, want %v", file, got, want)
							}
						}
					})
				}
			}
		}
	}
}

func TestPrivateMacAPIVariantsImplementTheSameFunctions(t *testing.T) {
	t.Parallel()
	declared := macSeamFunctions(t, "mac_private_api_darwin.h", `(?m)^void\s+(wailsPrivate\w+)\([^;]*\);`)
	if len(declared) == 0 {
		t.Fatal("no private API bridge declarations")
	}
	for _, name := range []string{privateMacAPIFile, "mac_public_api_darwin.go"} {
		defined := macSeamFunctions(t, name, `(?m)^void\s+(wailsPrivate\w+)\([^;]*\)\s*\{`)
		if !slices.Equal(declared, defined) {
			t.Errorf("%s implements %v, header declares %v", name, defined, declared)
		}
	}
}

func macSeamFunctions(t *testing.T, path, pattern string) []string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, match := range regexp.MustCompile(pattern).FindAllStringSubmatch(string(contents), -1) {
		names = append(names, match[1])
	}
	slices.Sort(names)
	return names
}
