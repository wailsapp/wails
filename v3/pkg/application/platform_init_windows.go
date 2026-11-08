//go:build windows && !server

package application

import (
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// startupThread is the thread the package was initialised on, which
// init_desktop.go keeps the main goroutine on. GetCurrentThreadId is in
// kernel32, which every process has loaded.
var startupThread = w32.GetCurrentThreadId()

// platformInit does what a GUI needs before its first window, when an App is
// built rather than at package init, so a program that links this package but
// never opens a window (a CLI, a daemon) does not load user32, uxtheme and COM.
// newPlatformApp does the same in Run, before services start; this keeps the
// main thread a COM apartment from New on, as package init used to, for code
// that runs between New and Run.
func platformInit() {
	registerWindowMessages()
	w32.InitDarkMode()
	// The COM apartment belongs to a thread: set it up here only on the
	// startup thread. If New runs on another goroutine, newPlatformApp (on the
	// startup thread) does it.
	if w32.GetCurrentThreadId() == startupThread {
		edge.InitCOM()
	}
}
