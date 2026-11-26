package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/features"
)

func TestHarmonyOSCLICommandIsRegistered(t *testing.T) {
	t.Setenv(features.WakeEnv, "1")
	app := newApplication(nil)

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	runErr := app.Run("harmonyos")
	writer.Close()
	os.Stdout = original
	if runErr != nil {
		t.Fatal(runErr)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	help := string(output)
	for _, command := range []string{"devices", "deploy", "logs", "overlay:gen"} {
		if !strings.Contains(help, command) {
			t.Errorf("harmonyos help does not list %q:\n%s", command, help)
		}
	}
}
