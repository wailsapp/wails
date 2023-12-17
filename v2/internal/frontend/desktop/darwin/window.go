//go:build darwin
// +build darwin

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

	"github.com/wailsapp/wails/v2/pkg/menu"

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

func bool2CboolPtr(value bool) *C.bool {
	v := C.bool(value)
	return &v
}

func NewWindow(frontendOptions *options.App, debug bool, devtools bool) *Window {
	c := NewCalloc()
	defer c.Free()

	frameless := bool2Cint(frontendOptions.Frameless)
	resizable := bool2Cint(!frontendOptions.DisableResize)
	fullscreen := bool2Cint(frontendOptions.Fullscreen)
	alwaysOnTop := bool2Cint(frontendOptions.AlwaysOnTop)
	hideWindowOnClose := bool2Cint(frontendOptions.HideWindowOnClose)
	startsHidden := bool2Cint(frontendOptions.StartHidden)
	devtoolsEnabled := bool2Cint(devtools)
	defaultContextMenuEnabled := bool2Cint(debug || frontendOptions.EnableDefaultContextMenu)
	singleInstanceEnabled := bool2Cint(frontendOptions.SingleInstanceLock != nil)

	var fullSizeContent, hideTitleBar, hideTitle, useToolbar, webviewIsTransparent C.int
	var titlebarAppearsTransparent, hideToolbarSeparator, windowIsTranslucent C.int
	var appearance, title *C.char
	var preferences C.struct_Preferences

	width := C.int(frontendOptions.Width)
	height := C.int(frontendOptions.Height)
	minWidth := C.int(frontendOptions.MinWidth)
	minHeight := C.int(frontendOptions.MinHeight)
	maxWidth := C.int(frontendOptions.MaxWidth)
	maxHeight := C.int(frontendOptions.MaxHeight)
	windowStartState := C.int(int(frontendOptions.WindowStartState))

	title = c.String(frontendOptions.Title)

	singleInstanceUniqueIdStr := ""
	if frontendOptions.SingleInstanceLock != nil {
		singleInstanceUniqueIdStr = frontendOptions.SingleInstanceLock.UniqueId
	}
	singleInstanceUniqueId := c.String(singleInstanceUniqueIdStr)

	enableFraudulentWebsiteWarnings := C.bool(frontendOptions.EnableFraudulentWebsiteDetection)

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

		if mac.Preferences != nil {
			if mac.Preferences.TabFocusesLinks.IsSet() {
				preferences.tabFocusesLinks = bool2CboolPtr(mac.Preferences.TabFocusesLinks.Get())
			}

			if mac.Preferences.TextInteractionEnabled.IsSet() {
				preferences.textInteractionEnabled = bool2CboolPtr(mac.Preferences.TextInteractionEnabled.Get())
			}

			if mac.Preferences.FullscreenEnabled.IsSet() {
				preferences.fullscreenEnabled = bool2CboolPtr(mac.Preferences.FullscreenEnabled.Get())
			}
		}

		windowIsTranslucent = bool2Cint(mac.WindowIsTranslucent)
		webviewIsTransparent = bool2Cint(mac.WebviewIsTransparent)

		appearance = c.String(string(mac.Appearance))
	}
	var context *C.WailsContext = C.Create(title, width, height, frameless, resizable, fullscreen, fullSizeContent,
		hideTitleBar, titlebarAppearsTransparent, hideTitle, useToolbar, hideToolbarSeparator, webviewIsTransparent,
		alwaysOnTop, hideWindowOnClose, appearance, windowIsTranslucent, devtoolsEnabled, defaultContextMenuEnabled,
		windowStartState, startsHidden, minWidth, minHeight, maxWidth, maxHeight, enableFraudulentWebsiteWarnings,
		preferences, singleInstanceEnabled, singleInstanceUniqueId,
	)

	// Create menu
	result := &Window{
		context: unsafe.Pointer(context),
	}

	if frontendOptions.BackgroundColour != nil {
		result.SetBackgroundColour(frontendOptions.BackgroundColour.R, frontendOptions.BackgroundColour.G, frontendOptions.BackgroundColour.B, frontendOptions.BackgroundColour.A)
	}

	if frontendOptions.Mac != nil && frontendOptions.Mac.About != nil {
		title := c.String(frontendOptions.Mac.About.Title)
		description := c.String(frontendOptions.Mac.About.Message)
		var icon unsafe.Pointer
		var length C.int
		if frontendOptions.Mac.About.Icon != nil {
			icon = unsafe.Pointer(&frontendOptions.Mac.About.Icon[0])
			length = C.int(len(frontendOptions.Mac.About.Icon))
		}
		C.SetAbout(result.context, title, description, icon, length)
	}

	if frontendOptions.Menu != nil {
		result.SetApplicationMenu(frontendOptions.Menu)
	}

	if debug && frontendOptions.Debug.OpenInspectorOnStartup {
		showInspector(result.context)
	}
	return result
}

func (w *Window) Center() {
	C.Center(w.context)
}

func (w *Window) Run(url string) {
	_url := C.CString(url)
	C.Run(w.context, _url)
	C.free(unsafe.Pointer(_url))
}

func (w *Window) Quit() {
	C.Quit(w.context)
}

func (w *Window) SetBackgroundColour(r uint8, g uint8, b uint8, a uint8) {
	C.SetBackgroundColour(w.context, C.int(r), C.int(g), C.int(b), C.int(a))
}

func (w *Window) ExecJS(js string) {
	_js := C.CString(js)
	C.ExecJS(w.context, _js)
	C.free(unsafe.Pointer(_js))
}

func (w *Window) SetPosition(x int, y int) {
	C.SetPosition(w.context, C.int(x), C.int(y))
}

func (w *Window) SetSize(width int, height int) {
	C.SetSize(w.context, C.int(width), C.int(height))
}

func (w *Window) SetAlwaysOnTop(onTop bool) {
	C.SetAlwaysOnTop(w.context, bool2Cint(onTop))
}

func (w *Window) SetTitle(title string) {
	t := C.CString(title)
	C.SetTitle(w.context, t)
	C.free(unsafe.Pointer(t))
}

func (w *Window) Maximise() {
	C.Maximise(w.context)
}

func (w *Window) ToggleMaximise() {
	C.ToggleMaximise(w.context)
}

func (w *Window) UnMaximise() {
	C.UnMaximise(w.context)
}

func (w *Window) IsMaximised() bool {
	return (bool)(C.IsMaximised(w.context))
}

func (w *Window) Minimise() {
	C.Minimise(w.context)
}

func (w *Window) UnMinimise() {
	C.UnMinimise(w.context)
}

func (w *Window) IsMinimised() bool {
	return (bool)(C.IsMinimised(w.context))
}

func (w *Window) IsNormal() bool {
	return !w.IsMaximised() && !w.IsMinimised() && !w.IsFullScreen()
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

func (w *Window) IsFullScreen() bool {
	return (bool)(C.IsFullScreen(w.context))
}

func (w *Window) Show() {
	C.Show(w.context)
}

func (w *Window) Hide() {
	C.Hide(w.context)
}

func (w *Window) ShowApplication() {
	C.ShowApplication(w.context)
}

func (w *Window) HideApplication() {
	C.HideApplication(w.context)
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

func (w *Window) GetPosition() (int, int) {
	var _result *C.char = C.GetPosition(w.context)
	temp := C.GoString(_result)
	return parseIntDuo(temp)
}

func (w *Window) Size() (int, int) {
	var _result *C.char = C.GetSize(w.context)
	temp := C.GoString(_result)
	return parseIntDuo(temp)
}

func (w *Window) SetApplicationMenu(inMenu *menu.Menu) {
	mainMenu := NewNSMenu(w.context, "")
	processMenu(mainMenu, inMenu)
	C.SetAsApplicationMenu(w.context, mainMenu.nsmenu)
}

func (w *Window) UpdateApplicationMenu() {
	C.UpdateApplicationMenu(w.context)
}

func (w Window) Print() {
	C.WindowPrint(w.context)
}
