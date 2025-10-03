//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
#cgo pkg-config: gtk4 webkitgtk-6.0

#include <jsc/jsc.h>
#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#include "window_webkit6.h"

*/
import "C"
import (
	"log"
	"strings"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

type Window struct {
	appoptions                               *options.App
	debug                                    bool
	devtoolsEnabled                          bool
	gtkApp                                   *C.GtkApplication
	gtkWindow                                unsafe.Pointer
	contentManager                           unsafe.Pointer
	webview                                  unsafe.Pointer
	applicationMenu                          *menu.Menu
	menubar                                  *C.GtkWidget
	webviewBox                               *C.GtkWidget
	vbox                                     *C.GtkWidget
	minWidth, minHeight, maxWidth, maxHeight int
}

func bool2Cint(value bool) C.int {
	if value {
		return C.int(1)
	}
	return C.int(0)
}

var activateWg sync.WaitGroup

//export onActivate
func onActivate() {
	activateWg.Done()
}

func NewWindow(appoptions *options.App, debug bool, devtoolsEnabled bool) *Window {
	validateWebKit2Version(appoptions)

	if appoptions.Linux == nil || appoptions.Linux.ProgramName == "" {
		log.Fatal("App option Linux.ProgramName required for gtk4. ex: org.author.app-name")
		return nil
	}

	result := &Window{
		appoptions:      appoptions,
		debug:           debug,
		devtoolsEnabled: devtoolsEnabled,
		minHeight:       appoptions.MinHeight,
		minWidth:        appoptions.MinWidth,
		maxHeight:       appoptions.MaxHeight,
		maxWidth:        appoptions.MaxWidth,
	}

	activateWg.Add(1)

	appId := C.CString(appoptions.Linux.ProgramName)
	defer C.free(unsafe.Pointer(appId))
	gtkApp := C.createApp(appId)
	result.gtkApp = gtkApp

	go func(gtkApp *C.GtkApplication) {
		C.runApp(gtkApp)
	}(gtkApp)

	activateWg.Wait()

	gtkWindow := C.gtk_application_window_new(gtkApp)
	C.g_object_ref_sink(C.gpointer(gtkWindow))
	result.gtkWindow = unsafe.Pointer(gtkWindow)

	webviewName := C.CString("webview-box")
	defer C.free(unsafe.Pointer(webviewName))
	result.webviewBox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
	C.gtk_widget_set_name(result.webviewBox, webviewName)

	result.vbox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
	C.gtk_window_set_child(result.asGTKWindow(), result.vbox)

	result.contentManager = unsafe.Pointer(C.webkit_user_content_manager_new())
	external := C.CString("external")
	defer C.free(unsafe.Pointer(external))
	C.webkit_user_content_manager_register_script_message_handler(result.cWebKitUserContentManager(), external, nil)
	C.SetupInvokeSignal(result.contentManager)

	var webviewGpuPolicy int
	if appoptions.Linux != nil {
		webviewGpuPolicy = int(appoptions.Linux.WebviewGpuPolicy)
	} else {
		// workaround for https://github.com/wailsapp/wails/issues/2977
		webviewGpuPolicy = int(linux.WebviewGpuPolicyNever)
	}

	webview := C.SetupWebview(
		result.contentManager,
		result.asGTKWindow(),
		bool2Cint(appoptions.HideWindowOnClose),
		C.int(webviewGpuPolicy),
		bool2Cint(appoptions.DragAndDrop != nil && appoptions.DragAndDrop.DisableWebViewDrop),
		bool2Cint(appoptions.DragAndDrop != nil && appoptions.DragAndDrop.EnableFileDrop),
	)
	result.webview = unsafe.Pointer(webview)
	buttonPressedName := C.CString("button-press-event")
	defer C.free(unsafe.Pointer(buttonPressedName))
	C.ConnectButtons(unsafe.Pointer(webview))

	if devtoolsEnabled {
		C.DevtoolsEnabled(unsafe.Pointer(webview), C.int(1), C.bool(debug && appoptions.Debug.OpenInspectorOnStartup))
		// Install Ctrl-Shift-F12 hotkey to call ShowInspector
		C.InstallF12Hotkey(result.gtkApp, result.asGTKWindow())
	}

	if !(debug || appoptions.EnableDefaultContextMenu) {
		C.DisableContextMenu(unsafe.Pointer(webview))
	}

	// Set background colour
	RGBA := appoptions.BackgroundColour
	result.SetBackgroundColour(RGBA.R, RGBA.G, RGBA.B, RGBA.A)

	// Setup window
	result.SetKeepAbove(appoptions.AlwaysOnTop)
	result.SetResizable(!appoptions.DisableResize)
	result.SetDefaultSize(appoptions.Width, appoptions.Height)
	result.SetDecorated(!appoptions.Frameless)
	result.SetTitle(appoptions.Title)
	result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
	result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)
	if appoptions.Linux != nil {
		if appoptions.Linux.Icon != nil {
			result.SetWindowIcon(appoptions.Linux.Icon)
		}
		if appoptions.Linux.WindowIsTranslucent {
			C.SetWindowTransparency(gtkWindow)
		}
	}

	// Menu
	result.SetApplicationMenu(appoptions.Menu)

	return result
}

func (w *Window) asGTKWidget() *C.GtkWidget {
	return C.GTKWIDGET(w.gtkWindow)
}

func (w *Window) asGTKWindow() *C.GtkWindow {
	return C.GTKWINDOW(w.gtkWindow)
}

func (w *Window) asGTKBox() *C.GtkBox {
	return C.GTKBOX(w.gtkWindow)
}

func (w *Window) cWebKitUserContentManager() *C.WebKitUserContentManager {
	return (*C.WebKitUserContentManager)(w.contentManager)
}

func (w *Window) Fullscreen() {
	C.ExecuteOnMainThread(C.Fullscreen, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnFullscreen() {
	if !w.IsFullScreen() {
		return
	}
	C.ExecuteOnMainThread(C.UnFullscreen, C.gpointer(w.asGTKWindow()))
	w.SetMinSize(w.minWidth, w.minHeight)
	w.SetMaxSize(w.maxWidth, w.maxHeight)
}

func (w *Window) Destroy() {
	C.gtk_window_destroy(w.asGTKWindow())
	C.g_object_unref(C.gpointer(w.gtkWindow))
}

func (w *Window) Close() {
	C.gtk_window_close(w.asGTKWindow())
}

func (w *Window) Center() {
	C.ExecuteOnMainThread(C.Center, C.gpointer(w.asGTKWindow()))
}

func (w *Window) SetPosition(x int, y int) {
	invokeOnMainThread(func() {
		C.SetPosition(unsafe.Pointer(w.asGTKWindow()), C.int(x), C.int(y))
	})
}

func (w *Window) Size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	invokeOnMainThread(func() {
		C.gtk_window_get_default_size(w.asGTKWindow(), &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *Window) GetPosition() (int, int) {
	//// TODO: gtk_window_get_position was removed in gtk4
	// var width, height C.int
	// var wg sync.WaitGroup
	// wg.Add(1)
	// invokeOnMainThread(func() {
	// 	C.gtk_window_get_position(w.asGTKWindow(), &width, &height)
	// 	wg.Done()
	// })
	// wg.Wait()
	// return int(width), int(height)
	return 0, 0
}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {
	w.maxHeight = maxHeight
	w.maxWidth = maxWidth
	invokeOnMainThread(func() {
		C.SetMinMaxSize(w.asGTKWindow(), C.int(w.minWidth), C.int(w.minHeight), C.int(w.maxWidth), C.int(w.maxHeight))
	})
}

func (w *Window) SetMinSize(minWidth int, minHeight int) {
	w.minHeight = minHeight
	w.minWidth = minWidth
	invokeOnMainThread(func() {
		C.SetMinMaxSize(w.asGTKWindow(), C.int(w.minWidth), C.int(w.minHeight), C.int(w.maxWidth), C.int(w.maxHeight))
	})
}

func (w *Window) Show() {
	C.ExecuteOnMainThread(C.Show, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Hide() {
	C.ExecuteOnMainThread(C.Hide, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Maximise() {
	C.ExecuteOnMainThread(C.Maximise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnMaximise() {
	C.ExecuteOnMainThread(C.UnMaximise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Minimise() {
	C.ExecuteOnMainThread(C.Minimise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnMinimise() {
	C.ExecuteOnMainThread(C.UnMinimise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) IsFullScreen() bool {
	result := C.IsFullscreen(w.asGTKWidget())
	if result != 0 {
		return true
	}
	return false
}

func (w *Window) IsMaximised() bool {
	result := C.IsMaximised(w.asGTKWidget())
	return result > 0
}

func (w *Window) IsMinimised() bool {
	result := C.IsMinimised(w.asGTKWidget())
	return result > 0
}

func (w *Window) IsNormal() bool {
	return !w.IsMaximised() && !w.IsMinimised() && !w.IsFullScreen()
}

func (w *Window) SetBackgroundColour(r uint8, g uint8, b uint8, a uint8) {
	windowIsTranslucent := false
	if w.appoptions.Linux != nil && w.appoptions.Linux.WindowIsTranslucent {
		windowIsTranslucent = true
	}
	data := C.RGBAOptions{
		r:                   C.uchar(r),
		g:                   C.uchar(g),
		b:                   C.uchar(b),
		a:                   C.uchar(a),
		webview:             w.webview,
		webviewBox:          unsafe.Pointer(w.webviewBox),
		windowIsTranslucent: gtkBool(windowIsTranslucent),
	}
	invokeOnMainThread(func() { C.SetBackgroundColour(unsafe.Pointer(&data)) })
}

func (w *Window) SetWindowIcon(icon []byte) {
	if len(icon) == 0 {
		return
	}
	C.SetWindowIcon(w.asGTKWindow(), (*C.guchar)(&icon[0]), (C.gsize)(len(icon)))
}

func (w *Window) Run(url string) {
	if w.menubar != nil {
		C.gtk_box_prepend(C.GTKBOX(unsafe.Pointer(w.vbox)), w.menubar)
	}

	C.gtk_box_prepend(C.GTKBOX(unsafe.Pointer(w.webviewBox)), C.GTKWIDGET(w.webview))
	C.gtk_box_append(C.GTKBOX(unsafe.Pointer(w.vbox)), w.webviewBox)

	_url := C.CString(url)
	defer C.free(unsafe.Pointer(_url))
	C.LoadIndex(w.webview, _url)
	if w.appoptions.StartHidden {
		w.Hide()
	}

	C.gtk_window_present(w.asGTKWindow())

	w.Center()
	switch w.appoptions.WindowStartState {
	case options.Fullscreen:
		w.Fullscreen()
	case options.Minimised:
		w.Minimise()
	case options.Maximised:
		w.Maximise()
	}
}

func (w *Window) SetKeepAbove(top bool) {
	//// TODO: gtk_window_set_keep_above was removed in gtk4
	// C.gtk_window_set_keep_above(w.asGTKWindow(), gtkBool(top))
}

func (w *Window) SetResizable(resizable bool) {
	C.gtk_window_set_resizable(w.asGTKWindow(), gtkBool(resizable))
}

func (w *Window) SetDefaultSize(width int, height int) {
	C.gtk_window_set_default_size(w.asGTKWindow(), C.int(width), C.int(height))
}

func (w *Window) SetSize(width int, height int) {
	C.gtk_window_set_default_size(w.asGTKWindow(), C.gint(width), C.gint(height))
}

func (w *Window) SetDecorated(frameless bool) {
	C.gtk_window_set_decorated(w.asGTKWindow(), gtkBool(frameless))
}

func (w *Window) SetTitle(title string) {
	// No need to free 'title' C-string here since it's freed in C-side
	C.SetTitle(w.asGTKWindow(), C.CString(title))
}

func (w *Window) ExecJS(js string) {
	jscallback := C.JSCallback{
		webview: w.webview,
		script:  C.CString(js),
	}
	invokeOnMainThread(func() { C.ExecuteJS(unsafe.Pointer(&jscallback)) })
}

func (w *Window) StartDrag() {
	C.StartDrag(w.webview, w.asGTKWindow())
}

func (w *Window) StartResize(edge uintptr) {
	C.StartResize(w.webview, w.asGTKWindow(), C.GdkSurfaceEdge(edge))
}

func (w *Window) Quit() {
	if mainLoop == nil {
		return
	}

	C.g_main_loop_quit(mainLoop)
	C.g_main_loop_unref(mainLoop)

	mainLoop = nil
}

func (w *Window) OpenFileDialog(dialogOptions frontend.OpenDialogOptions, multipleFiles int, action C.GtkFileChooserAction) {

	data := C.OpenFileDialogOptions{
		window:        w.asGTKWindow(),
		title:         C.CString(dialogOptions.Title),
		multipleFiles: C.int(multipleFiles),
		action:        action,
	}

	if len(dialogOptions.Filters) > 0 {
		// Create filter array
		mem := NewCalloc()
		arraySize := len(dialogOptions.Filters) + 1
		data.filters = C.AllocFileFilterArray((C.size_t)(arraySize))
		filters := unsafe.Slice((**C.struct__GtkFileFilter)(unsafe.Pointer(data.filters)), arraySize)
		for index, filter := range dialogOptions.Filters {
			thisFilter := C.gtk_file_filter_new()
			C.g_object_ref(C.gpointer(thisFilter))
			if filter.DisplayName != "" {
				cName := mem.String(filter.DisplayName)
				C.gtk_file_filter_set_name(thisFilter, cName)
			}
			if filter.Pattern != "" {
				for _, thisPattern := range strings.Split(filter.Pattern, ";") {
					cThisPattern := mem.String(thisPattern)
					C.gtk_file_filter_add_pattern(thisFilter, cThisPattern)
				}
			}
			// Add filter to array
			filters[index] = thisFilter
		}
		mem.Free()
		filters[arraySize-1] = nil
	}

	if dialogOptions.CanCreateDirectories {
		data.createDirectories = C.int(1)
	}

	if dialogOptions.ShowHiddenFiles {
		data.showHiddenFiles = C.int(1)
	}

	if dialogOptions.DefaultFilename != "" {
		data.defaultFilename = C.CString(dialogOptions.DefaultFilename)
	}

	if dialogOptions.DefaultDirectory != "" {
		data.defaultDirectory = C.CString(dialogOptions.DefaultDirectory)
	}

	invokeOnMainThread(func() { C.Opendialog(unsafe.Pointer(&data)) })
}

func (w *Window) MessageDialog(dialogOptions frontend.MessageDialogOptions) {

	data := C.MessageDialogOptions{
		window:  w.gtkWindow,
		title:   C.CString(dialogOptions.Title),
		message: C.CString(dialogOptions.Message),
	}
	switch dialogOptions.Type {
	case frontend.InfoDialog:
		data.messageType = C.int(0)
	case frontend.ErrorDialog:
		data.messageType = C.int(1)
	case frontend.QuestionDialog:
		data.messageType = C.int(2)
	case frontend.WarningDialog:
		data.messageType = C.int(3)
	}
	invokeOnMainThread(func() { C.MessageDialog(unsafe.Pointer(&data)) })
}

func (w *Window) ToggleMaximise() {
	if w.IsMaximised() {
		w.UnMaximise()
	} else {
		w.Maximise()
	}
}

func (w *Window) ShowInspector() {
	invokeOnMainThread(func() { C.ShowInspector(w.webview) })
}

// showModalDialogAndExit shows a modal dialog and exits the app.
func showModalDialogAndExit(title, message string) {
	go func() {
		data := C.MessageDialogOptions{
			title:       C.CString(title),
			message:     C.CString(message),
			messageType: C.int(1),
		}

		C.MessageDialog(unsafe.Pointer(&data))
	}()

	<-messageDialogResult
	log.Fatal(message)
}
