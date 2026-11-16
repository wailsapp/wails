package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNativeLinuxPackageVersionRetainsPrerelease(t *testing.T) {
	for _, test := range []struct{ format, version, want string }{
		{"archlinux", "1.2.3-beta.2+git.9", "1.2.3beta.2"},
		{"archlinux", "1.2.3-alpha-12", "1.2.3alpha12"},
		{"archlinux", "1.2.3", "1.2.3"},
		{"deb", "1.2.3-beta.2", "1.2.3-beta.2"},
		{"rpm", "1.2.3-beta.2", "1.2.3-beta.2"},
	} {
		t.Run(test.format+"/"+test.version, func(t *testing.T) {
			root := t.TempDir()
			spec := packageTestSpec("linux", "amd64", test.format)
			spec.Project.Version = test.version
			assets := filepath.Join(root, spec.Assets, "linux", "nfpm")
			if err := os.MkdirAll(assets, 0755); err != nil {
				t.Fatal(err)
			}
			source := filepath.Join(assets, "nfpm.yaml")
			original := []byte("name: example\nversion: 0.0.0\ncontents: []\n")
			if err := os.WriteFile(source, original, 0644); err != nil {
				t.Fatal(err)
			}
			path, err := (&manifestHandler{root: root}).prepareLinuxPackageConfig(spec)
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "version: "+test.want+"\n") || strings.Contains(string(data), "version_schema: none\n") != (test.format == "archlinux") {
				t.Fatalf("incorrect native version policy:\n%s", data)
			}
			unchanged, err := os.ReadFile(source)
			if err != nil || string(unchanged) != string(original) {
				t.Fatalf("mutated platform assets: %s, %v", unchanged, err)
			}
		})
	}
}
