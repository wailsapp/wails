//go:build windows

package application

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
)

// The recovery decision and the attempt budget are the two pieces of the
// process-failure path that can be exercised without a live WebView2 runtime:
// neither touches COM. The rest (rebuildWebView, the re-navigation itself) needs
// a real controller and is covered by the manual matrix in issue #5733.

func TestWebviewRecoveryActionFor(t *testing.T) {
	tests := []struct {
		name string
		kind edge.COREWEBVIEW2_PROCESS_FAILED_KIND
		want webviewRecoveryAction
	}{
		{
			// A dead browser process invalidates the controller permanently,
			// so only a new controller recovers it.
			name: "browser process exit rebuilds",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED,
			want: webviewRecoveryRebuild,
		},
		{
			name: "render process exit re-navigates",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED,
			want: webviewRecoveryRenavigate,
		},
		{
			name: "unresponsive render process re-navigates",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE,
			want: webviewRecoveryRenavigate,
		},
		{
			// Chromium re-creates a dead out-of-process iframe itself; reloading
			// the whole window over one subframe would be worse than the failure.
			name: "frame render process exit is left alone",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED,
			want: webviewRecoveryNone,
		},
		{
			// Chromium restarts these itself, and when it gives up it exits the
			// browser process, which comes back as BROWSER_PROCESS_EXITED.
			name: "gpu process exit is left alone",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED,
			want: webviewRecoveryNone,
		},
		{
			name: "utility process exit is left alone",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_UTILITY_PROCESS_EXITED,
			want: webviewRecoveryNone,
		},
		{
			name: "sandbox helper process exit is left alone",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND_SANDBOX_HELPER_PROCESS_EXITED,
			want: webviewRecoveryNone,
		},
		{
			// GetProcessFailedKind seeds its out-param with 0xffffffff and a
			// newer runtime may report a kind this build has no constant for.
			// Anything unrecognised must fall through to "leave it alone"
			// rather than trigger a rebuild.
			name: "unknown kind is left alone",
			kind: edge.COREWEBVIEW2_PROCESS_FAILED_KIND(0xffffffff),
			want: webviewRecoveryNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := webviewRecoveryActionFor(tt.kind); got != tt.want {
				t.Errorf("webviewRecoveryActionFor(%d) = %d, want %d", tt.kind, got, tt.want)
			}
		})
	}
}

func TestBeginWebviewRecoveryStopsAtTheLimit(t *testing.T) {
	w := &windowsWebviewWindow{}

	for i := 1; i <= maxWebviewRecoveryAttempts; i++ {
		if !w.beginWebviewRecovery() {
			t.Fatalf("attempt %d refused, want allowed within the budget of %d",
				i, maxWebviewRecoveryAttempts)
		}
		if w.webviewRecoveryAttempts != i {
			t.Fatalf("after attempt %d: webviewRecoveryAttempts = %d, want %d",
				i, w.webviewRecoveryAttempts, i)
		}
	}

	// Past the budget it must keep refusing rather than letting the count (and
	// the WebView2 process spawning behind it) run away.
	for i := 0; i < 3; i++ {
		if w.beginWebviewRecovery() {
			t.Fatalf("attempt %d past the budget of %d was allowed",
				maxWebviewRecoveryAttempts+i+1, maxWebviewRecoveryAttempts)
		}
	}
	if w.webviewRecoveryAttempts != maxWebviewRecoveryAttempts {
		t.Errorf("webviewRecoveryAttempts = %d after refused attempts, want %d",
			w.webviewRecoveryAttempts, maxWebviewRecoveryAttempts)
	}
}

// A completed navigation means recovery worked, so the budget must come back —
// otherwise a long-running app that recovered once would have fewer attempts
// available for an unrelated failure hours later, and eventually none.
//
// This covers resetWebviewRecoveryBudget, not its call site: navigationCompleted
// needs a live controller to invoke.
func TestResetWebviewRecoveryBudget(t *testing.T) {
	w := &windowsWebviewWindow{}

	for i := 0; i < maxWebviewRecoveryAttempts; i++ {
		if !w.beginWebviewRecovery() {
			t.Fatalf("attempt %d refused while filling the budget", i+1)
		}
	}
	if w.beginWebviewRecovery() {
		t.Fatal("budget not exhausted before the reset; test cannot prove the reset works")
	}

	w.resetWebviewRecoveryBudget()

	if w.webviewRecoveryAttempts != 0 {
		t.Errorf("webviewRecoveryAttempts = %d after reset, want 0", w.webviewRecoveryAttempts)
	}
	if !w.beginWebviewRecovery() {
		t.Error("recovery still refused after the budget was reset")
	}
}
