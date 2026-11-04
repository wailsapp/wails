package taskfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseNeverEvaluatesLegacyCommands(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	source := `version: '3'
vars:
  DYNAMIC: {sh: 'echo evaluated > variable-marker'}
tasks:
  build:
    preconditions: [{sh: 'echo evaluated > precondition-marker'}]
    cmds: ['echo evaluated > command-marker']
`
	path := filepath.Join(root, "Taskfile.yml")
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	parsed, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Vars["DYNAMIC"].Shell == "" || len(parsed.Tasks["build"].Cmds) != 1 {
		t.Fatal("declarations lost during inspection")
	}
	for _, name := range []string{"variable-marker", "precondition-marker", "command-marker"} {
		if _, err := os.Stat(filepath.Join(root, name)); !os.IsNotExist(err) {
			t.Fatalf("migration evaluated %s: %v", name, err)
		}
	}
}

func FuzzParseDeclarations(f *testing.F) {
	for _, source := range []string{"version: '3'\ntasks: {}\n", "version: '3'\ntasks:\n  build: {cmds: [echo hello]}\n", "[]", "version: '3'\nincludes: [one]\n", "version: '3'\nvars: [one]\n"} {
		f.Add(source)
	}
	f.Fuzz(func(t *testing.T, source string) {
		if len(source) > 1<<20 {
			t.Skip()
		}
		path := filepath.Join(t.TempDir(), "Taskfile.yml")
		if err := os.WriteFile(path, []byte(source), 0644); err != nil {
			t.Fatal(err)
		}
		_, _ = Parse(path)
	})
}
