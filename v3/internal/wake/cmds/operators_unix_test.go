//go:build !windows

package cmds

import (
	"os"
	"path/filepath"
	"testing"
)

func installFakeGo(t *testing.T, script string) string {
	t.Helper()
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "go"), []byte("#!/bin/sh\n"+script+"\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", base+string(os.PathListSeparator)+os.Getenv("PATH"))
	return base
}

func TestRouteExecutesCommandChain(t *testing.T) {
	base := installFakeGo(t, "exit 0")
	command := Route("go build && echo chain-complete > result.txt", RouteOptions{Dir: base})
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(base, "result.txt")); err != nil || string(data) != "chain-complete\n" {
		t.Errorf("command chain output = %q, %v", data, err)
	}
}
