package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

#include <stdlib.h>
*/
import "C"
import (
	"log"
	"runtime"
	"strconv"
	"strings"
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
	debug := bool2Cint(true)
	alpha := C.int(frontendOptions.RGBA.A)
	red := C.int(frontendOptions.RGBA.R)
	green := C.int(frontendOptions.RGBA.G)
	blue := C.int(frontendOptions.RGBA.B)

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
	var context *C.WailsContext = C.Create(title, width, height, frameless, resizable, fullscreen, fullSizeContent, hideTitleBar, titlebarAppearsTransparent, hideTitle, useToolbar, hideToolbarSeparator, webviewIsTransparent, alwaysOnTop, hideWindowOnClose, appearance, windowIsTranslucent, debug)

	C.free(unsafe.Pointer(title))
	if appearance != nil {
		C.free(unsafe.Pointer(appearance))
	}

	C.SetRGBA(unsafe.Pointer(context), red, green, blue, alpha)

	return &Window{
		context: unsafe.Pointer(context),
	}
}

func (w *Window) Center() {
	C.Center(w.context)
}

func (w *Window) Run() {
	C.Run(w.context)
	println("I exited!")
}

func (w *Window) Quit() {
	C.Quit(w.context)
}

func (w *Window) SetRGBA(r uint8, g uint8, b uint8, a uint8) {
	C.SetRGBA(w.context, C.int(r), C.int(g), C.int(b), C.int(a))
}

func (w *Window) ExecJS(js string) {
	_js := C.CString(js)
	C.ExecJS(w.context, _js)
	C.free(unsafe.Pointer(_js))
}

func (w *Window) SetPos(x int, y int) {
	C.SetPosition(w.context, C.int(x), C.int(y))
}

func (w *Window) SetSize(width int, height int) {
	C.SetSize(w.context, C.int(width), C.int(height))
}

func (w *Window) SetTitle(title string) {
	t := C.CString(title)
	C.SetTitle(w.context, t)
	C.free(unsafe.Pointer(t))
}

func (w *Window) Maximise() {
	C.Maximise(w.context)
}

func (w *Window) UnMaximise() {
	C.UnMaximise(w.context)
}

func (w *Window) Minimise() {
	C.Minimise(w.context)
}

func (w *Window) UnMinimise() {
	C.UnMinimise(w.context)
}

func (w *Window) SetMinSize(width int, height int) {
	C.SetMinSize(w.context, C.int(width), C.int(height))
}

func (w *Window) SetMaxSize(width int, height int) {
	C.SetMaxSize(w.context, C.int(width), C.int(height))
}

func (w *Window) Fullscreen() {
	C.Fullscreen(w.context)
}

func (w *Window) UnFullscreen() {
	C.UnFullscreen(w.context)
}

func (w *Window) Show() {
	C.Show(w.context)
}

func (w *Window) Hide() {
	C.Hide(w.context)
}

func parseIntDuo(temp string) (int, int) {
	split := strings.Split(temp, ",")
	x, err := strconv.Atoi(split[0])
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal(err)
	}
	return x, y
}

func (w *Window) Pos() (int, int) {
	var _result *C.char = C.GetPos(w.context)
	temp := C.GoString(_result)
	return parseIntDuo(temp)
}

func (w *Window) Size() (int, int) {
	var _result *C.char = C.GetSize(w.context)
	temp := C.GoString(_result)
	return parseIntDuo(temp)
}
