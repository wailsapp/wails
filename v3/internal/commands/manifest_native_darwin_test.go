//go:build darwin

package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"howett.net/plist"
)

func TestNativeYAMLAppAndDMGSigningWithCompanion(t *testing.T) {
	if os.Getenv("WAILS_TEST_NATIVE_DMG") != "1" {
		t.Skip("opt in to native codesign and hdiutil integration")
	}
	root := executableFixture(t)
	t.Chdir(root)
	path := filepath.Join(root, "wails.yaml")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	source = []byte(strings.Replace(string(source), "kind: cli", "kind: gui", 1) + "\ndarwin:\n  microphone_usage_description: Record the selected input\n  signing:\n    identity: '-'\nprofiles:\n  app:\n    targets:\n      darwin/" + runtime.GOARCH + ":\n        sign: true\n  disk:\n    targets:\n      darwin/" + runtime.GOARCH + ":\n        formats: [dmg]\n        sign: true\n")
	if err := os.WriteFile(path, source, 0644); err != nil {
		t.Fatal(err)
	}
	command := func(name string, args ...string) []byte {
		t.Helper()
		out, err := exec.CommandContext(t.Context(), name, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, err, out)
		}
		return out
	}
	for _, profile := range []string{"app", "disk"} {
		if _, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", Profile: profile}); err != nil {
			t.Fatal(err)
		}
		if profile == "app" {
			app := filepath.Join(root, "bin", "executables.signed.app")
			command("codesign", "--verify", "--deep", "--strict", app)
			command(filepath.Join(app, "Contents/MacOS/executables"))
			output := command(filepath.Join(app, "Contents/MacOS/helper"))
			if strings.TrimSpace(string(output)) != "companion one" {
				t.Fatalf("companion output: %s", output)
			}
			data, err := os.ReadFile(filepath.Join(app, "Contents/Info.plist"))
			if err != nil {
				t.Fatal(err)
			}
			var values map[string]any
			if _, err := plist.Unmarshal(data, &values); err != nil {
				t.Fatal(err)
			}
			if values["NSMicrophoneUsageDescription"] != "Record the selected input" {
				t.Fatal(values)
			}
		} else {
			dmg := filepath.Join(root, "bin", "executables.signed.dmg")
			command("hdiutil", "verify", dmg)
			command("codesign", "--verify", "--strict", dmg)
			output := command("hdiutil", "attach", "-plist", "-nobrowse", "-readonly", dmg)
			var mounted struct {
				Entities []struct {
					Mount string `plist:"mount-point"`
				} `plist:"system-entities"`
			}
			if _, err := plist.Unmarshal(output, &mounted); err != nil {
				t.Fatal(err)
			}
			found := false
			for _, entity := range mounted.Entities {
				if entity.Mount == "" {
					continue
				}
				found = true
				defer command("hdiutil", "detach", entity.Mount)
				app := filepath.Join(entity.Mount, "executables.app")
				command("codesign", "--verify", "--deep", "--strict", app)
				if _, err := os.Stat(filepath.Join(app, "Contents/MacOS/helper")); err != nil {
					t.Fatal(err)
				}
			}
			if !found {
				t.Fatal("DMG did not mount")
			}
		}
	}
}
