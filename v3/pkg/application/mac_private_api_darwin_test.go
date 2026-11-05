//go:build darwin && !ios && !server

package application

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// Compile and exercise the actual Objective-C bridges independently of Go's
// test process so AppKit runs on the main thread. Also inspect production
// objects: selectors in this test's fixtures must not mask binary leakage.
func TestMacPrivateAPINativeVariants(t *testing.T) {
	if _, err := exec.LookPath("xcrun"); err != nil {
		t.Skip("requires Xcode command line tools")
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, private := range []bool{false, true} {
		for _, devtools := range []bool{false, true} {
			name := "public"
			source := "mac_public_api_darwin.go"
			if private {
				name, source = "private", privateMacAPIFile
			}
			if devtools {
				name += "/devtools"
			} else {
				name += "/production"
			}
			t.Run(name, func(t *testing.T) {
				dir := t.TempDir()
				bridge := filepath.Join(dir, "bridge.m")
				native := macNativePreamble(t, source)
				if devtools {
					native = append(native, macNativePreamble(t, "webview_window_darwin_dev.go")...)
				}
				if err := os.WriteFile(bridge, native, 0600); err != nil {
					t.Fatal(err)
				}
				args := []string{"clang", "-Werror", "-I", root}
				if devtools {
					args = append(args, "-DWAILS_MAC_DEVTOOLS")
				}
				object := filepath.Join(dir, "bridge.o")
				macNativeCommand(t, "xcrun", append(slices.Clone(args), "-c", bridge, "-o", object)...)
				data, err := os.ReadFile(object)
				if err != nil {
					t.Fatal(err)
				}
				for selector, want := range map[string]bool{
					"backgroundColor":        private,
					"drawsBackground":        private,
					"setGroupIdentifier:":    private,
					"setGroupName:":          private,
					"setGroupSpacing:":       private,
					"_inspector":             private && devtools,
					"developerExtrasEnabled": private && devtools,
				} {
					if got := bytes.Contains(data, []byte(selector+"\x00")); got != want {
						t.Errorf("compiled bridge contains %q = %v, want %v", selector, got, want)
					}
				}
				executable := filepath.Join(dir, "native-test")
				if private {
					args = append(args, "-DEXPECT_PRIVATE=1")
				}
				args = append(args, object, "testdata/mac_private_apis/main.m", "-framework", "Cocoa", "-framework", "WebKit", "-o", executable)
				macNativeCommand(t, "xcrun", args...)
				macNativeCommand(t, executable)
			})
		}
	}
}

func macNativePreamble(t *testing.T, path string) []byte {
	t.Helper()
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, rest, ok := strings.Cut(string(source), "/*\n")
	if !ok {
		t.Fatalf("%s has no cgo preamble", path)
	}
	preamble, _, ok := strings.Cut(rest, "*/")
	if !ok {
		t.Fatalf("%s has unterminated cgo preamble", path)
	}
	var native strings.Builder
	for _, line := range strings.Split(preamble, "\n") {
		if !strings.HasPrefix(line, "#cgo ") {
			native.WriteString(line + "\n")
		}
	}
	return []byte(native.String())
}

func macNativeCommand(t *testing.T, name string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if output, err := exec.CommandContext(ctx, name, args...).CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, output)
	}
}
