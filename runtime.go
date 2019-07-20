package wails

import (
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/runtime"
)

// CustomLogger type alias
type CustomLogger = logger.CustomLogger

// Runtime is the Wails Runtime Interface, given to a user who has defined the WailsInit method
type Runtime struct {
	Events     *runtime.Events
	Log        *runtime.Log
	Dialog     *runtime.Dialog
	Window     *runtime.Window
	Browser    *runtime.Browser
	FileSystem *runtime.FileSystem
}

// NewRuntime creates a new Runtime struct
func NewRuntime(eventManager interfaces.EventManager, renderer interfaces.Renderer) *Runtime {
	return &Runtime{
		Events:     runtime.NewEvents(eventManager),
		Log:        runtime.NewLog(),
		Dialog:     runtime.NewDialog(renderer),
		Window:     runtime.NewWindow(renderer),
		Browser:    runtime.NewBrowser(),
		FileSystem: runtime.NewFileSystem(),
	}
}
