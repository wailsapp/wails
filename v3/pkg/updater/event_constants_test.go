package updater_test

import (
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

// TestEventConstants_MatchRuntimeTS guards against drift between the
// Go-side event constants and the JS-side mirror in
// internal/runtime/desktop/@wailsio/runtime/src/updater.ts. Both files
// expose the same wire strings; if they diverge, app developers who
// import the JS constants would silently miss events.
func TestEventConstants_MatchRuntimeTS(t *testing.T) {
	const tsPath = "../../internal/runtime/desktop/@wailsio/runtime/src/updater.ts"
	body, err := os.ReadFile(tsPath)
	if err != nil {
		t.Fatalf("read %s: %v", tsPath, err)
	}
	ts := string(body)

	expect := []string{
		updater.EventCheckStarted,
		updater.EventUpdateAvailable,
		updater.EventNoUpdate,
		updater.EventDownloadStarted,
		updater.EventDownloadProgress,
		updater.EventDownloadComplete,
		updater.EventVerifying,
		updater.EventInstalling,
		updater.EventUpdateReady,
		updater.EventError,
		updater.EventMeta,
		updater.EventWindowReady,
		updater.EventUserInstall,
		updater.EventUserSkip,
		updater.EventUserRemind,
		updater.EventUserCancel,
		updater.EventUserRestart,
	}
	for _, s := range expect {
		// Each constant should appear as a quoted string literal in the .ts file.
		if !strings.Contains(ts, `"`+s+`"`) {
			t.Errorf("updater.ts is missing event constant %q — JS mirror has drifted from Go", s)
		}
	}
}
