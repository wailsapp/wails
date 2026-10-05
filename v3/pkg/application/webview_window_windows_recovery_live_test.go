//go:build windows && !server && webview2_recovery_test

package application

import (
	"fmt"
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Run with powershell -File scripts/test-webview2-recovery.ps1 on an interactive Windows desktop.
// The script overlays HRESULT injection points without changing repository files.
func TestLiveWebviewRecoveryFailures(t *testing.T) {
	if os.Getenv("WAILS_RECOVERY_TEST_CHILD") == "1" {
		w := newTestMenuWindow(t)
		if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED); err != nil {
			t.Fatal(err)
		}
		defer windows.CoUninitialize()
		globalApplication = &App{Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)), options: Options{
			Windows:      WindowsOptions{WebviewUserDataPath: filepath.Join(os.TempDir(), "wails-6167-live-fixed")},
			ErrorHandler: func(err error) { fmt.Fprintln(os.Stderr, "ERROR_HANDLER", err) },
		}}
		w.parent.options.BackgroundColour = RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}
		w.parent.options.HTML = "<p>Recovered</p>"
		w.chromium = w.newChromium()
		if !w.setupChromium(true) {
			t.Fatal("initial setup failed")
		}
		previous := w.chromium
		mode := os.Getenv("WAILS_RECOVERY_TEST_MODE")
		once := os.Getenv("WAILS_RECOVERY_TEST_ONCE") == "1"
		os.Setenv("WAILS_RECOVERY_TEST_FAULT", mode)
		w.beginWebviewRecovery()
		w.webviewRecoveryPending = true
		w.rebuildWebView()
		if !w.webviewRecoveryPending {
			t.Fatal("pending guard cleared while rebuilding")
		}
		if previous.GetController() != nil || previous.IsReady() {
			t.Fatal("abandoned controller not closed")
		}
		if once {
			if w.webviewRecoveryAttempts != 2 || !w.chromium.IsReady() {
				t.Fatalf("transient failure: attempts=%d ready=%v", w.webviewRecoveryAttempts, w.chromium.IsReady())
			}
			navigated := false
			w.chromium.NavigationCompletedCallback = func(_ *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationCompletedEventArgs) {
				ok, err := args.GetIsSuccess()
				navigated = ok && err == nil
			}
			deadline := time.Now().Add(15 * time.Second)
			for !navigated && time.Now().Before(deadline) {
				var msg w32.MSG
				for w32.PeekMessage(&msg, 0, 0, 0, 1) {
					w32.TranslateMessage(&msg)
					w32.DispatchMessage(&msg)
				}
				time.Sleep(5 * time.Millisecond)
			}
			if !navigated {
				t.Fatal("replacement never navigated successfully")
			}
		} else {
			if w.webviewRecoveryAttempts != 3 || w.chromium.IsReady() || w.chromium.GetController() != nil {
				t.Fatalf("exhaustion cleanup failed attempts=%d", w.webviewRecoveryAttempts)
			}
		}
		w.chromium.Close()
		w.chromium.Close()
		fmt.Fprintf(os.Stderr, "RECOVERY_SURVIVED mode=%s transient=%v attempts=%d\n", mode, once, w.webviewRecoveryAttempts)
		return
	}
	for _, mode := range []string{"background", "devtools", "resource", "embed"} {
		for _, once := range []string{"0", "1"} {
			t.Run(mode+"/transient="+once, func(t *testing.T) {
				exe, err := os.Executable()
				if err != nil {
					t.Fatal(err)
				}
				cmd := exec.Command(exe, "-test.run=^TestLiveWebviewRecoveryFailures$", "-test.v")
				cmd.Env = append(os.Environ(), "WAILS_RECOVERY_TEST_CHILD=1", "WAILS_RECOVERY_TEST_MODE="+mode, "WAILS_RECOVERY_TEST_ONCE="+once)
				out, err := cmd.CombinedOutput()
				t.Logf("%s", out)
				if err != nil || !strings.Contains(string(out), "RECOVERY_SURVIVED") || !strings.Contains(string(out), "INJECT_HRESULT "+mode) {
					t.Fatalf("recovery failed: %v", err)
				}
			})
		}
	}
}
