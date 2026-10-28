package commands

import (
	"debug/pe"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sysoTestOptions(t *testing.T) *SysoOptions {
	t.Helper()
	root := t.TempDir()
	icon, err := buildAssets.ReadFile("build_assets/windows/icon.ico")
	if err != nil {
		t.Fatal(err)
	}
	iconPath, manifestPath := filepath.Join(root, "icon.ico"), filepath.Join(root, "app.manifest")
	if err := os.WriteFile(iconPath, icon, 0600); err != nil {
		t.Fatal(err)
	}
	manifest := `<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0"><assemblyIdentity name="Example" version="1.0.0.0" type="win32" processorArchitecture="*"/></assembly>`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	return &SysoOptions{Icon: iconPath, Manifest: manifestPath, Out: filepath.Join(root, "output.syso"), Arch: "amd64"}
}

func TestGenerateSysoArchitectures(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	for _, tc := range []struct {
		arch    string
		machine uint16
	}{
		{"amd64", pe.IMAGE_FILE_MACHINE_AMD64},
		{"arm64", pe.IMAGE_FILE_MACHINE_ARM64},
		{"386", pe.IMAGE_FILE_MACHINE_I386},
		{"unsupported", 0},
	} {
		t.Run(tc.arch, func(t *testing.T) {
			options := sysoTestOptions(t)
			options.Arch = tc.arch
			if err := os.WriteFile(options.Out, []byte("previous object"), 0600); err != nil {
				t.Fatal(err)
			}
			err := GenerateSyso(options)
			if tc.machine == 0 {
				if err == nil || !strings.Contains(err.Error(), "not supported") {
					t.Errorf("error=%v, want unsupported architecture", err)
				}
				got, err := os.ReadFile(options.Out)
				if err != nil || string(got) != "previous object" {
					t.Errorf("previous output=%q, error=%v", got, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			object, err := pe.Open(options.Out)
			if err != nil {
				t.Fatal(err)
			}
			defer object.Close()
			if object.Machine != tc.machine {
				t.Errorf("COFF machine=%#x, want %#x", object.Machine, tc.machine)
			}
		})
	}
}
