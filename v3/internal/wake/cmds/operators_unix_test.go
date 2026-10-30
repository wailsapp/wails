//go:build !windows

package cmds

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func installFakeGo(t *testing.T, script string) string {
	t.Helper()
	return installFakeCommand(t, "go", script)
}

func installFakeCommand(t *testing.T, name, script string) string {
	t.Helper()
	base := t.TempDir()
	if err := os.WriteFile(filepath.Join(base, name), []byte("#!/bin/sh\n"+script+"\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", base+string(os.PathListSeparator)+os.Getenv("PATH"))
	return base
}

func captureRoutedOutput(t *testing.T, command, directory string) string {
	t.Helper()
	cmd := Route(command, RouteOptions{Dir: directory})
	var output bytes.Buffer
	cmd.(OutputSetter).SetOutput(&output, &output)
	if err := cmd.Run(); err != nil {
		t.Fatalf("routed command %q failed: %v; output: %s", command, err, output.String())
	}
	return output.String()
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
