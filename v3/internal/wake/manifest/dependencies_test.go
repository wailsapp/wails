package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExplicitEmptyDependenciesSurviveProfilesResolutionAndEjection(t *testing.T) {
	root := t.TempDir()
	source := `version: 3
project:
  name: dependencies
  product_name: Dependencies
  identifier: org.example.dependencies
  version: 1.0.0
packages:
  deb:
    dependencies: []
  rpm:
    dependencies: [custom]
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm]
`
	path := filepath.Join(root, Filename)
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	for _, profile := range []string{"", "release"} {
		loaded, err := Load(root, profile)
		if err != nil {
			t.Fatal(err)
		}
		resolved, err := ResolvePackageFormat(loaded.Config.Package, "linux", "deb")
		if err != nil {
			t.Fatal(err)
		}
		if resolved.Dependencies == nil || len(resolved.Dependencies) != 0 {
			t.Fatal("explicit empty dependency declaration was lost")
		}
		rpm, err := ResolvePackageFormat(loaded.Config.Package, "linux", "rpm")
		if err != nil {
			t.Fatal(err)
		}
		rpm.Dependencies[0] = "mutated"
		if loaded.Config.Package.Linux.RPM.Dependencies[0] != "custom" {
			t.Fatal("resolved dependencies alias the source configuration")
		}
		encoded, err := EncodeEjectedYAML(loaded.Config, "test")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(encoded), "dependencies: []") {
			t.Fatalf("ejection omitted opt-out:\n%s", encoded)
		}
		ejected := filepath.Join(root, "ejected.yaml")
		if err := os.WriteFile(ejected, encoded, 0600); err != nil {
			t.Fatal(err)
		}
		reloaded, err := LoadFile(root, ejected, profile)
		if err != nil {
			t.Fatal(err)
		}
		result, err := ResolvePackageFormat(reloaded.Config.Package, "linux", "deb")
		if err != nil || result.Dependencies == nil || len(result.Dependencies) != 0 {
			t.Fatalf("reload lost opt-out: %+v %v", result, err)
		}
	}
}

func TestDependenciesRemainOptionalInEditorSchema(t *testing.T) {
	data, err := JSONSchema()
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatal(err)
	}
	var visit func(any)
	count := 0
	visit = func(value any) {
		switch value := value.(type) {
		case map[string]any:
			if properties, ok := value["properties"].(map[string]any); ok {
				if dependency, ok := properties["dependencies"].(map[string]any); ok {
					count++
					if dependency["type"] != "array" {
						t.Fatalf("dependency schema=%v", dependency)
					}
					if required, ok := value["required"].([]any); ok {
						for _, key := range required {
							if key == "dependencies" {
								t.Fatal("dependencies became required")
							}
						}
					}
				}
			}
			for _, child := range value {
				visit(child)
			}
		case []any:
			for _, child := range value {
				visit(child)
			}
		}
	}
	visit(root)
	if count == 0 {
		t.Fatal("no dependency schema was checked")
	}
}
