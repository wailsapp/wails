package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"howett.net/plist"
)

func TestMicrophoneMetadataOnlyChangesStagedFiles(t *testing.T) {
	for _, custom := range []bool{false, true} {
		t.Run(map[bool]string{false: "generated-entitlements", true: "user-entitlements"}[custom], func(t *testing.T) {
			root := t.TempDir()
			stage := filepath.Join(root, "staged")
			if err := os.MkdirAll(filepath.Join(stage, "darwin"), 0755); err != nil {
				t.Fatal(err)
			}
			original := []byte(`<?xml version="1.0" encoding="UTF-8"?><plist version="1.0"><dict><key>CFBundleName</key><string>Fixture</string></dict></plist>`)
			source := filepath.Join(root, "Info.plist")
			if err := os.WriteFile(source, original, 0444); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(stage, "darwin", "Info.plist")
			if err := os.WriteFile(path, original, 0444); err != nil {
				t.Fatal(err)
			}
			config := manifest.Config{Root: root}
			if custom {
				config.Signing.Darwin.Entitlements = "user.plist"
				if err := os.WriteFile(filepath.Join(root, "user.plist"), []byte("owned"), 0444); err != nil {
					t.Fatal(err)
				}
			}
			handler := manifestHandler{root: root, config: config}
			if err := handler.applyMicrophoneMetadata(stage, pipeline.AssetsSpec{TargetOS: "darwin", MicrophoneUsageDescription: "Record the selected input"}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var info map[string]any
			if _, err := plist.Unmarshal(data, &info); err != nil {
				t.Fatal(err)
			}
			if info["CFBundleName"] != "Fixture" || info["NSMicrophoneUsageDescription"] != "Record the selected input" {
				t.Fatal(info)
			}
			untouched, err := os.ReadFile(source)
			if err != nil || string(untouched) != string(original) {
				t.Fatalf("source changed: %s %v", untouched, err)
			}
			generated := filepath.Join(stage, "signing", "entitlements.plist")
			if custom {
				if _, err := os.Stat(generated); !os.IsNotExist(err) {
					t.Fatalf("user entitlements overridden: %v", err)
				}
				data, err := os.ReadFile(filepath.Join(root, "user.plist"))
				if err != nil || string(data) != "owned" {
					t.Fatal("user entitlements modified")
				}
			} else {
				data, err := os.ReadFile(generated)
				if err != nil {
					t.Fatal(err)
				}
				var entitlements map[string]any
				if _, err := plist.Unmarshal(data, &entitlements); err != nil {
					t.Fatal(err)
				}
				if entitlements["com.apple.security.device.audio-input"] != true {
					t.Fatal(entitlements)
				}
			}
		})
	}
}
