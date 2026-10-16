package version

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/debug"
)

func TestStringTrimsEmbeddedVersionWhitespace(t *testing.T) {
	originalVersion := versionString
	originalModulePath := debug.LocalModulePath
	t.Cleanup(func() {
		versionString = originalVersion
		debug.LocalModulePath = originalModulePath
	})

	versionString = "v3.0.0-beta.1\r\n"
	debug.LocalModulePath = ""

	if got := String(); got != "v3.0.0-beta.1" {
		t.Errorf("String() = %q, want %q", got, "v3.0.0-beta.1")
	}
}
