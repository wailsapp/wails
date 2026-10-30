package cmds

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoFmtDryRunPreservesSource(t *testing.T) {
	original := "package fixture\n\nfunc answer( )int{return 42}\n"
	root, output, err := runOfflineGo(t, "go fmt -n ./...", map[string]string{
		"go.mod":  "module example.com/wake-fmt-check\n\ngo 1.25.0\n",
		"main.go": original,
	})
	if err != nil {
		t.Fatalf("dry run failed: %v; output: %s", err, output)
	}
	if !strings.Contains(output, "gofmt") {
		t.Errorf("missing planned formatter command: %s", output)
	}
	data, err := os.ReadFile(filepath.Join(root, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != original {
		t.Errorf("dry run formatted source:\n%s", data)
	}
}
