package test4551

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateBindingsUsesStrippedSymbols(t *testing.T) {
	sourceFile := filepath.Join("..", "..", "pkg", "commands", "bindings", "bindings.go")
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Skipf("could not read bindings source: %v", err)
	}

	source := string(data)
	if !strings.Contains(source, `"-ldflags", "-s -w"`) {
		t.Error("expected wailsbindings build command to include -ldflags -s -w for stripped symbols, " +
			"which is required to avoid the Go 1.25 DWARF5+CGO Windows PE linker bug (golang/go#75077)")
	}
}
