package application

import "runtime"

// Platform identifies a single operating system / build target, for use with
// application.System.IsPlatform.
type Platform int

const (
	PlatformMacOS Platform = iota
	PlatformWindows
	PlatformLinux
	PlatformIOS
	PlatformAndroid
	// PlatformServer is the headless server build (the "server" build tag),
	// regardless of the underlying OS.
	PlatformServer
)

// String returns a human-readable platform name.
func (p Platform) String() string {
	switch p {
	case PlatformMacOS:
		return "macOS"
	case PlatformWindows:
		return "Windows"
	case PlatformLinux:
		return "Linux"
	case PlatformIOS:
		return "iOS"
	case PlatformAndroid:
		return "Android"
	case PlatformServer:
		return "Server"
	default:
		return "unknown"
	}
}

// systemManager is the receiver for runtime platform/mode queries. The
// package-level System singleton is the entry point: application.System.IsMobile(),
// application.System.IsDesktop(), application.System.IsServer(),
// application.System.IsPlatform(...).
type systemManager struct{}

// System reports the platform and mode the application is running under.
//
// Unlike the IOS and Android managers (which only exist on their respective
// builds), System is compiled into every build, so shared code can branch on
// it at runtime without build tags. Each platform is checked explicitly, so
// prefer a switch over an "else" fall-through:
//
//	switch {
//	case application.System.IsMobile():
//	    // iOS / Android
//	case application.System.IsServer():
//	    // headless server
//	case application.System.IsDesktop():
//	    // macOS / Windows / Linux
//	}
//
// Exactly one of IsMobile, IsDesktop and IsServer is true for any given build.
var System systemManager

// IsMobile reports whether the application is running on a mobile OS
// (iOS or Android).
func (systemManager) IsMobile() bool {
	if isServerBuild {
		return false
	}
	return runtime.GOOS == "ios" || runtime.GOOS == "android"
}

// IsServer reports whether the application was built in server mode (the
// "server" build tag) — a headless HTTP server with no native GUI.
func (systemManager) IsServer() bool {
	return isServerBuild
}

// IsDesktop reports whether the application is running as a native desktop GUI
// app (macOS, Windows or Linux). It checks the desktop platforms explicitly
// rather than treating "not mobile" as desktop.
func (systemManager) IsDesktop() bool {
	if isServerBuild {
		return false
	}
	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		return true
	default:
		return false
	}
}

// IsPlatform reports whether the application is running on the given platform —
// a direct test for a single target:
//
//	if application.System.IsPlatform(application.PlatformMacOS) { ... }
//
// PlatformServer matches any OS built with the "server" tag; the OS platforms
// (macOS/Windows/Linux/iOS/Android) only match native, non-server builds.
func (systemManager) IsPlatform(p Platform) bool {
	if p == PlatformServer {
		return isServerBuild
	}
	if isServerBuild {
		return false
	}
	switch p {
	case PlatformMacOS:
		return runtime.GOOS == "darwin"
	case PlatformWindows:
		return runtime.GOOS == "windows"
	case PlatformLinux:
		return runtime.GOOS == "linux"
	case PlatformIOS:
		return runtime.GOOS == "ios"
	case PlatformAndroid:
		return runtime.GOOS == "android"
	default:
		return false
	}
}
