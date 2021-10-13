package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func init() {
	runtime.LockOSThread()
}

type Window struct {
	context unsafe.Pointer
}

func bool2Cint(value bool) C.int {
	if value {
		return C.int(1)
	}
	return C.int(0)
}

func NewWindow(frontendOptions *options.App) *Window {

	frameless := bool2Cint(frontendOptions.Frameless)
	resizable := bool2Cint(!frontendOptions.DisableResize)
	fullscreen := bool2Cint(frontendOptions.Fullscreen)
	alwaysOnTop := bool2Cint(frontendOptions.AlwaysOnTop)
	webviewIsTransparent := bool2Cint(frontendOptions.AlwaysOnTop)
	hideWindowOnClose := bool2Cint(frontendOptions.HideWindowOnClose)

	var fullSizeContent, hideTitleBar, hideTitle, useToolbar C.int
	var titlebarAppearsTransparent, hideToolbarSeparator, windowIsTranslucent C.int
	var appearance, title *C.char

	width := C.int(frontendOptions.Width)
	height := C.int(frontendOptions.Height)

	title = C.CString(frontendOptions.Title)

	if frontendOptions.Mac != nil {
		mac := frontendOptions.Mac
		if mac.TitleBar != nil {
			fullSizeContent = bool2Cint(mac.TitleBar.FullSizeContent)
			hideTitleBar = bool2Cint(mac.TitleBar.HideTitleBar)
			hideTitle = bool2Cint(mac.TitleBar.HideTitle)
			useToolbar = bool2Cint(mac.TitleBar.UseToolbar)
			titlebarAppearsTransparent = bool2Cint(mac.TitleBar.TitlebarAppearsTransparent)
			hideToolbarSeparator = bool2Cint(mac.TitleBar.HideToolbarSeparator)
		}
		windowIsTranslucent = bool2Cint(mac.WindowIsTranslucent)
		appearance = C.CString(string(mac.Appearance))
	}
	var context *C.WailsContext = C.Create(title, width, height, frameless, resizable, fullscreen, fullSizeContent, hideTitleBar, titlebarAppearsTransparent, hideTitle, useToolbar, hideToolbarSeparator, webviewIsTransparent, alwaysOnTop, hideWindowOnClose, appearance, windowIsTranslucent)

	C.free(unsafe.Pointer(title))
	if appearance != nil {
		C.free(unsafe.Pointer(appearance))
	}

	return &Window{
		context: unsafe.Pointer(context),
	}
}

//func (w *Window) Center() {
//	C.Center(w.wailsApplication)
//}

func (w *Window) Run() {
	C.Run(w.context)
	println("I exited!")
}
