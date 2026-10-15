//go:build windows

// live_smoke_test.go exercises the bindings against the *real* WebView2 stack
// installed on the machine, rather than the fake recording trampolines used by
// the generated *_gen_test.go suite. It is the only layer that proves the
// generator's model of the ABI matches Microsoft's actual DLL end to end.
//
// Three tiers, increasing in what they require from the host:
//
//   - TestLiveCompareBrowserVersions: pure loader logic, no runtime needed.
//   - TestLiveInstalledRuntimeVersion: probes the installed Evergreen runtime;
//     skips cleanly when none is present.
//   - TestLiveEnvironmentCreation: creates a real CoreWebView2 environment and
//     calls a generated binding method (GetBrowserVersionString) through the
//     live vtable. Gated behind WEBVIEW2_LIVE=1 because it spawns a browser
//     process and pumps a message loop, which is too heavy/flaky to ride along
//     with the fast hermetic suite in `go test ./...`.

package webview2_test

import (
	"os"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	webview2 "github.com/wailsapp/wails/v3/internal/webview2/pkg/webview2"
	"github.com/wailsapp/wails/v3/internal/webview2/webviewloader"
)

// TestLiveCompareBrowserVersions exercises the loader's version comparison
// against the real (pure-Go) implementation. Deterministic; no runtime needed.
func TestLiveCompareBrowserVersions(t *testing.T) {
	cases := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0.0.0", "2.0.0.0", -1},
		{"2.0.0.0", "1.0.0.0", 1},
		{"1.0.705.50", "1.0.705.50", 0},
		{"94.0.992.31", "94.0.992.30", 1},
	}
	for _, tc := range cases {
		got, err := webviewloader.CompareBrowserVersions(tc.v1, tc.v2)
		if err != nil {
			t.Fatalf("CompareBrowserVersions(%q,%q): %v", tc.v1, tc.v2, err)
		}
		if got != tc.want {
			t.Errorf("CompareBrowserVersions(%q,%q) = %d, want %d", tc.v1, tc.v2, got, tc.want)
		}
	}
}

// installedRuntimeVersion returns the installed Evergreen runtime version, or
// "" when no runtime is present.
func installedRuntimeVersion(t *testing.T) string {
	t.Helper()
	version, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		t.Fatalf("GetAvailableCoreWebView2BrowserVersionString: %v", err)
	}
	return version
}

// TestLiveInstalledRuntimeVersion probes the installed runtime via the loader.
// This is a genuine integration point: the pure-Go loader locates the real
// runtime DLL on disk (via the registry) and reads its version resource.
func TestLiveInstalledRuntimeVersion(t *testing.T) {
	version := installedRuntimeVersion(t)
	if version == "" {
		t.Skip("no WebView2 runtime installed; skipping live runtime probe")
	}
	t.Logf("installed WebView2 runtime: %s", version)

	// The reported string must be a parseable version (CompareBrowserVersions
	// rejects malformed input), and must be newer than the floor SDK we target.
	if _, err := webviewloader.CompareBrowserVersions(version, version); err != nil {
		t.Fatalf("installed runtime version %q is not parseable: %v", version, err)
	}
	cmp, err := webviewloader.CompareBrowserVersions(version, "86.0.616.0")
	if err != nil {
		t.Fatalf("CompareBrowserVersions: %v", err)
	}
	if cmp < 0 {
		t.Errorf("installed runtime %q is older than the minimum supported 86.0.616.0", version)
	}
}

// liveEnvHandler implements webviewloader.ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler.
//
// The handler is invoked on the pump thread inside DispatchMessage, and the
// loader Release()s the environment immediately after this returns, so we must
// call the generated binding method here and record the result rather than
// stashing the (soon-to-be-freed) pointer.
type liveEnvHandler struct {
	done    bool
	hr      webviewloader.HRESULT
	version string
	verr    error
}

func (h *liveEnvHandler) EnvironmentCompleted(errorCode webviewloader.HRESULT, env *webviewloader.ICoreWebView2Environment) webviewloader.HRESULT {
	h.hr = errorCode
	if errorCode == 0 && env != nil {
		// env already points at an ICoreWebView2Environment (the loader QI'd it
		// for us). webviewloader.ICoreWebView2Environment and the generated
		// ICoreWebView2Environment are both a single vtable-pointer header over
		// the same COM object, so reinterpreting the pointer lets us drive the
		// *generated* binding against the live DLL.
		genEnv := (*webview2.ICoreWebView2Environment)(unsafe.Pointer(env))
		h.version, h.verr = genEnv.GetBrowserVersionString()
	}
	h.done = true
	return 0 // S_OK
}

// win32MSG mirrors the Win32 MSG struct for the message pump.
type win32MSG struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	private uint32
}

// TestLiveEnvironmentCreation creates a real CoreWebView2 environment and calls
// the generated GetBrowserVersionString binding through the live vtable, then
// confirms it agrees with the loader-reported runtime version. Opt-in: it
// spawns a browser process and pumps a message loop.
func TestLiveEnvironmentCreation(t *testing.T) {
	if os.Getenv("WEBVIEW2_LIVE") != "1" {
		t.Skip("set WEBVIEW2_LIVE=1 to run the live environment-creation smoke test")
	}
	runtimeVersion := installedRuntimeVersion(t)
	if runtimeVersion == "" {
		t.Skip("no WebView2 runtime installed; cannot create an environment")
	}

	// WebView2 delivers the environment-created callback on this thread's COM
	// message queue, so pin the goroutine and pump messages.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	const coinitApartmentThreaded = 0x2
	if err := windows.CoInitializeEx(0, coinitApartmentThreaded); err != nil {
		// S_FALSE (already initialized) is reported as a non-nil error by some
		// versions; only a hard failure should abort.
		t.Logf("CoInitializeEx returned %v (continuing)", err)
	}
	defer windows.CoUninitialize()

	h := &liveEnvHandler{}
	if err := webviewloader.CreateCoreWebView2Environment(h); err != nil {
		t.Fatalf("CreateCoreWebView2Environment: %v", err)
	}

	user32 := windows.NewLazySystemDLL("user32.dll")
	peekMessage := user32.NewProc("PeekMessageW")
	translateMessage := user32.NewProc("TranslateMessage")
	dispatchMessage := user32.NewProc("DispatchMessageW")

	const pmRemove = 0x0001
	var msg win32MSG
	deadline := time.Now().Add(30 * time.Second)
	for !h.done {
		if time.Now().After(deadline) {
			t.Fatal("timed out after 30s waiting for the environment-created callback")
		}
		r, _, _ := peekMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, pmRemove)
		if r != 0 {
			translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
			continue
		}
		time.Sleep(5 * time.Millisecond)
	}

	if h.hr != 0 {
		t.Fatalf("environment creation failed: HRESULT 0x%08x", uint32(h.hr))
	}
	if h.verr != nil {
		t.Fatalf("generated GetBrowserVersionString failed: %v", h.verr)
	}
	if h.version == "" {
		t.Fatal("generated GetBrowserVersionString returned an empty version")
	}
	t.Logf("environment BrowserVersionString (via generated binding): %s", h.version)

	// The version reported through the generated binding must match the one the
	// loader found. Channel suffixes can differ, so compare numerically.
	cmp, err := webviewloader.CompareBrowserVersions(h.version, runtimeVersion)
	if err != nil {
		t.Fatalf("CompareBrowserVersions(%q,%q): %v", h.version, runtimeVersion, err)
	}
	if cmp != 0 {
		t.Errorf("binding version %q != loader version %q", h.version, runtimeVersion)
	}
}
