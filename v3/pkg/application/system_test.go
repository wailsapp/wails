package application

import (
	"runtime"
	"testing"
)

// Tests build without the "server" tag, so isServerBuild is false here.
func TestSystemModes(t *testing.T) {
	wantMobile := runtime.GOOS == "ios" || runtime.GOOS == "android"

	if got := System.IsServer(); got {
		t.Errorf("IsServer() = true; want false in a non-server build")
	}
	if got := System.IsMobile(); got != wantMobile {
		t.Errorf("IsMobile() = %v; want %v for GOOS=%q", got, wantMobile, runtime.GOOS)
	}
	if got := System.IsDesktop(); got != !wantMobile {
		t.Errorf("IsDesktop() = %v; want %v for GOOS=%q", got, !wantMobile, runtime.GOOS)
	}
}

// Exactly one of IsMobile/IsDesktop/IsServer must be true for any build.
func TestSystemModesMutuallyExclusive(t *testing.T) {
	count := 0
	for _, b := range []bool{System.IsMobile(), System.IsDesktop(), System.IsServer()} {
		if b {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one of IsMobile/IsDesktop/IsServer to be true, got %d", count)
	}
}

// IsPlatform must match the current GOOS (non-server build) and exactly one
// platform should be true.
func TestIsPlatform(t *testing.T) {
	want := map[Platform]bool{
		PlatformMacOS:   runtime.GOOS == "darwin",
		PlatformWindows: runtime.GOOS == "windows",
		PlatformLinux:   runtime.GOOS == "linux",
		PlatformIOS:     runtime.GOOS == "ios",
		PlatformAndroid: runtime.GOOS == "android",
		PlatformServer:  false, // non-server test build
	}
	trueCount := 0
	for p, w := range want {
		if got := System.IsPlatform(p); got != w {
			t.Errorf("IsPlatform(%s) = %v; want %v (GOOS=%q)", p, got, w, runtime.GOOS)
		}
		if want[p] {
			trueCount++
		}
	}
	if trueCount != 1 {
		t.Errorf("expected exactly one platform to match GOOS=%q, got %d", runtime.GOOS, trueCount)
	}
}
