package commands

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestLinuxPackageRuntimeDependencies(t *testing.T) {
	for _, test := range []struct {
		name, kind, companion string
		dependencies, want    []string
	}{
		{name: "CLI", kind: "cli"},
		{name: "server", kind: "server"},
		{name: "GUI", kind: "gui", want: []string{"libgtk-4-1", "libwebkitgtk-6.0-4"}},
		{name: "GUI companion", kind: "cli", companion: "gui", want: []string{"libgtk-4-1", "libwebkitgtk-6.0-4"}},
		{name: "explicit headless", kind: "cli", dependencies: []string{"custom-runtime"}, want: []string{"custom-runtime"}},
		{name: "explicit GUI", kind: "gui", dependencies: []string{"custom-runtime"}, want: []string{"custom-runtime"}},
		{name: "explicit empty GUI", kind: "gui", dependencies: []string{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler, spec := linuxPackageFixture(t)
			spec.Project.Kind = test.kind
			if test.companion != "" {
				spec.Companions[0].Kind = test.companion
			}
			spec.Config.Dependencies = test.dependencies
			if _, err := handler.packageLinux(spec); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(handler.root, spec.Output))
			if err != nil {
				t.Fatal(err)
			}
			controlArchive := tar.NewReader(bytes.NewReader(debTarMember(t, data, "control.tar.gz")))
			var control string
			for {
				entry, err := controlArchive.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				if strings.TrimPrefix(entry.Name, "./") == "control" {
					contents, err := io.ReadAll(controlArchive)
					if err != nil {
						t.Fatal(err)
					}
					control = string(contents)
					break
				}
			}
			if control == "" {
				t.Fatal("missing DEB control file")
			}
			var got []string
			for _, line := range strings.Split(control, "\n") {
				if value, ok := strings.CutPrefix(line, "Depends: "); ok && value != "" {
					got = strings.Split(value, ", ")
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("dependencies=%v want=%v\n%s", got, test.want, control)
			}
		})
	}
}

func TestLinuxPackageDependencyOverridesAreFormatScoped(t *testing.T) {
	for _, format := range []string{"rpm", "archlinux"} {
		for _, test := range []struct {
			name, kind   string
			dependencies []string
			defaults     bool
		}{
			{name: "headless defaults", kind: "cli"},
			{name: "GUI defaults", kind: "gui", defaults: true},
			{name: "explicit GUI", kind: "gui", dependencies: []string{"declared"}},
			{name: "explicit empty GUI", kind: "gui", dependencies: []string{}},
		} {
			t.Run(format+" "+test.name, func(t *testing.T) {
				selected := map[string]any{"depends": []string{"stock-GTK"}, "custom_setting": "kept"}
				other := map[string]any{"depends": []string{"other-format"}}
				config := map[string]any{"depends": []string{"base-GTK"}, "overrides": map[string]any{format: selected, "other": other}}
				spec := pipeline.PackageSpec{Format: format, Project: manifest.Project{Kind: test.kind}, Config: manifest.PackageFormat{Dependencies: test.dependencies}}
				setLinuxPackageDependencies(config, spec)
				if test.defaults {
					if !reflect.DeepEqual(selected["depends"], []string{"stock-GTK"}) {
						t.Fatal("changed GUI defaults")
					}
				} else {
					if _, exists := selected["depends"]; exists {
						t.Fatal("selected stock override masks dependency decision")
					}
					if !reflect.DeepEqual(config["depends"], append([]string{}, test.dependencies...)) {
						t.Fatalf("lost explicit dependency policy: %v", config)
					}
				}
				if selected["custom_setting"] != "kept" || !reflect.DeepEqual(other["depends"], []string{"other-format"}) {
					t.Fatal("changed unrelated configuration")
				}
			})
		}
	}
}

func TestLinuxPackageCustomTemplateRemainsOpaque(t *testing.T) {
	handler, spec := linuxPackageFixture(t)
	source := []byte("# user-owned\ndepends: [custom]\ncontents: []\n")
	if err := os.WriteFile(filepath.Join(handler.root, "custom.yaml"), source, 0600); err != nil {
		t.Fatal(err)
	}
	spec.Config.Template = "custom.yaml"
	path, err := handler.prepareLinuxPackageConfig(spec)
	if err != nil {
		t.Fatal(err)
	}
	copied, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(copied, source) {
		t.Fatalf("template changed: %v", err)
	}
}

func TestLinuxDependencyOptOutHasDistinctCacheIdentity(t *testing.T) {
	spec := pipeline.PackageSpec{Format: "deb", Config: manifest.PackageFormat{Dependencies: nil}}
	inherited, err := cache.ActionKey(string(pipeline.PackageArtifact), spec, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	spec.Config.Dependencies = []string{}
	optedOut, err := cache.ActionKey(string(pipeline.PackageArtifact), spec, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if inherited == optedOut {
		t.Fatal("explicit empty dependencies reused the GUI-default package cache key")
	}
}
