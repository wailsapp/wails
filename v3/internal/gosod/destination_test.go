package gosod

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestExtractPreservesLiteralDestination(t *testing.T) {
	parent := t.TempDir()
	for _, tt := range []struct{ name, destination string }{
		{"action", "{{.Name}}"},
		{"unclosed action", "unfinished{{"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := filepath.Join(parent, tt.name)
			if err := os.MkdirAll(root, 0755); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, tt.destination)
			files := fstest.MapFS{
				"{{.Name}}/plain.txt":         {Data: []byte("plain")},
				"{{.Name}}/rendered.txt.tmpl": {Data: []byte("Hello {{.Name}}")},
			}
			if err := New(files).Extract(target, struct{ Name string }{"example"}); err != nil {
				t.Fatal(err)
			}
			for name, want := range map[string]string{"plain.txt": "plain", "rendered.txt": "Hello example"} {
				if got := readFile(t, filepath.Join(target, "example", name)); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
			entries, err := os.ReadDir(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 || entries[0].Name() != tt.destination {
				t.Errorf("destination siblings = %v; want only %q", entries, tt.destination)
			}
		})
	}
}
