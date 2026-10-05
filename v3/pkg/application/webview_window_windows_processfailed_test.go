//go:build windows && !server

package application

import (
	"errors"
	"golang.org/x/sys/windows"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"unsafe"

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
	w := webviewWindowWithExhaustedRecoveryBudget(t)
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

func TestUnsuccessfulNavigationDoesNotResetWebviewRecoveryBudget(t *testing.T) {
	w := webviewWindowWithExhaustedRecoveryBudget(t)

	if shouldResetWebviewRecoveryBudget(false, nil) {
		w.resetWebviewRecoveryBudget()
	}

	if w.webviewRecoveryAttempts != maxWebviewRecoveryAttempts {
		t.Fatalf("webviewRecoveryAttempts = %d after unsuccessful navigation, want %d",
			w.webviewRecoveryAttempts, maxWebviewRecoveryAttempts)
	}
	if w.beginWebviewRecovery() {
		t.Fatal("unsuccessful navigation reset the exhausted recovery budget")
	}
}

func webviewWindowWithExhaustedRecoveryBudget(t *testing.T) *windowsWebviewWindow {
	t.Helper()
	w := &windowsWebviewWindow{}

	for i := 0; i < maxWebviewRecoveryAttempts; i++ {
		if !w.beginWebviewRecovery() {
			t.Fatalf("attempt %d refused while filling the budget", i+1)
		}
	}
	return w
}

func TestProcessFailedDoesNotQueueOverlappingRecovery(t *testing.T) {
	w := &windowsWebviewWindow{webviewRecoveryPending: true, webviewRecoveryAttempts: 1}
	// A nested pump may dispatch another process failure while rebuilding.
	// Ignore it before consulting event args or spending another attempt.
	w.processFailed(nil, nil)
	if w.webviewRecoveryAttempts != 1 {
		t.Fatal("overlapping failure spent another recovery attempt")
	}
}

// Recovery errors must return without invoking the fatal startup handler.
func TestSetupChromiumConfigErrorDuringRecoveryDoesNotTerminate(t *testing.T) {
	previous := globalApplication
	globalApplication = &App{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	t.Cleanup(func() { globalApplication = previous })

	w := &windowsWebviewWindow{}
	if got := w.setupChromiumConfigError(true, errors.New("boom")); got {
		t.Fatal("setupChromiumConfigError(recovering=true, ...) = true, want false")
	}
}

func TestSetupChromiumConfigErrorOnStartupTerminates(t *testing.T) {
	if os.Getenv("WAILS_TEST_FATAL_RECOVERY") == "1" {
		globalApplication = &App{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
		(&windowsWebviewWindow{}).setupChromiumConfigError(false, errors.New("startup failure"))
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(executable, "-test.run=^TestSetupChromiumConfigErrorOnStartupTerminates$")
	cmd.Env = append(os.Environ(), "WAILS_TEST_FATAL_RECOVERY=1")
	output, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		t.Fatalf("startup failure: %v, want exit 1; output: %s", err, output)
	}
}

func TestEnableDevToolsReturnsCOMFailure(t *testing.T) {
	// ICoreWebView2Settings inherits three IUnknown slots. Its DevTools setter
	// is slot 12; use the native ABI so both build-tag implementations are tested.
	vtable := [13]edge.ComProc{}
	vtable[12] = edge.NewComProc(func(uintptr, uintptr) uintptr { return 0x80004005 })
	object := struct{ vtable unsafe.Pointer }{unsafe.Pointer(&vtable)}
	settings := (*edge.ICoreWebViewSettings)(unsafe.Pointer(&object))
	err := (&windowsWebviewWindow{}).enableDevTools(settings)
	runtime.KeepAlive(vtable)
	if !errors.Is(err, windows.Errno(0x80004005)) {
		t.Fatalf("DevTools configuration error = %v, want E_FAIL", err)
	}
}

func TestZoomAfterFailedWebviewRecovery(t *testing.T) {
	w := &windowsWebviewWindow{chromium: edge.NewChromium()}
	w.chromium.Close()
	if zoom := w.getZoom(); zoom != -1 {
		t.Fatalf("zoom without controller = %v, want -1", zoom)
	}
	w.zoomIn()
	w.zoomOut()
}

func TestRebuildWebViewStopsAfterWindowCloses(t *testing.T) {
	for _, w := range []*windowsWebviewWindow{
		{parent: &WebviewWindow{destroyed: true}, hwnd: 1},
		{parent: &WebviewWindow{}, hwnd: 0},
	} {
		w.webviewRecoveryAttempts = 1
		w.rebuildWebView()
		if w.webviewRecoveryAttempts != 1 {
			t.Fatal("closed window spent a recovery attempt")
		}
	}
}
