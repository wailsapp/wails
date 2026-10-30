package cmds

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func runOfflineGo(t *testing.T, command string, files map[string]string) (string, string, error) {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		target := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := Route(command, RouteOptions{
		Dir: root,
		Env: append(os.Environ(), "GOWORK=off", "GOPROXY=off", "GOSUMDB=off", "GOTOOLCHAIN=local", "GOFLAGS=", "GOBIN="+filepath.Join(root, "bin")),
	})
	var output bytes.Buffer
	cmd.(OutputSetter).SetOutput(&output, &output)
	err := cmd.Run()
	return root, output.String(), err
}
