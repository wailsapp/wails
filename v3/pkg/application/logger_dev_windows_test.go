//go:build windows && !production

package application

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

func TestDefaultLoggerRedirectedStderr(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	previous := os.Stderr
	os.Stderr = file
	t.Cleanup(func() {
		os.Stderr = previous
		file.Close()
	})

	var level slog.LevelVar
	logger := DefaultLogger(&level)
	logger.Debug("filtered")
	logger.Info("visible", "key", "value")
	level.Set(slog.LevelDebug)
	logger.Debug("later")

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)
	if strings.Contains(output, "\x1b[") || strings.Contains(output, "filtered") {
		t.Fatalf("unexpected colors or filtered record: %q", output)
	}
	for _, want := range []string{"INF visible key=value", "DBG later"} {
		if !strings.Contains(output, want) {
			t.Errorf("output %q does not contain %q", output, want)
		}
	}
}

func TestEnableConsoleColorsPreservesModes(t *testing.T) {
	file, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err != nil {
		t.Skip("no console attached")
	}
	defer file.Close()
	handle := windows.Handle(file.Fd())
	var original uint32
	if err := windows.GetConsoleMode(handle, &original); err != nil {
		t.Fatal(err)
	}
	defer windows.SetConsoleMode(handle, original)
	if !enableConsoleColors(file) {
		t.Skip("console does not accept native ANSI output")
	}
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		t.Fatal(err)
	}
	want := original | windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	if mode != want {
		t.Errorf("mode = %#x, want %#x", mode, want)
	}
}
