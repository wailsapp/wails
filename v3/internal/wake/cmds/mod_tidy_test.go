package cmds

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGoModTidyDiffPreservesModule(t *testing.T) {
	module := "module example.com/wake-tidy-check\n\ngo 1.25.0\n\nrequire example.com/unused v0.0.0\n\nreplace example.com/unused => ./unused\n"
	root, output, err := runOfflineGo(t, "go mod tidy -diff", map[string]string{
		"go.mod":        module,
		"main.go":       "package fixture\n",
		"unused/go.mod": "module example.com/unused\n\ngo 1.25.0\n",
	})
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Errorf("untidy dry run error = %v, want exit 1; output: %s", err, output)
	}
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != module {
		t.Errorf("dry run modified go.mod:\n%s", data)
	}
}
