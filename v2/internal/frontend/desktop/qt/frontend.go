//go:build qt
// +build qt

package qt

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/url"
	"text/template"
	"time"
	"unsafe"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	wailsruntime "github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/assetserver"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

/*
#cgo linux pkg-config: Qt5Widgets Qt5Core Qt5WebEngineWidgets
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -L/usr/local/lib -lstdc++

#include "lib.hpp"

*/
import "C"

var exitCh = make(chan int)

const startURL = "wails://wails/"

//export appExited
func appExited(retCode C.int) {
	exitCh <- int(retCode)
}

type Frontend struct {
	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	debug           bool
	devtoolsEnabled bool

	// Assets
	assets   *assetserver.AssetServer
	startURL *url.URL

	// main window handle
	//mainWindow *Window
	bindings   *binding.Bindings
	dispatcher frontend.Dispatcher

	qApp    unsafe.Pointer
	qWindow *C.Window
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {
	f := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
		ctx:             ctx,
	}
	f.startURL, _ = url.Parse(startURL)

	if _starturl, _ := ctx.Value("starturl").(*url.URL); _starturl != nil {
		f.startURL = _starturl
	} else {
		if port, _ := ctx.Value("assetserverport").(string); port != "" {
			f.startURL.Host = net.JoinHostPort(f.startURL.Host+".localhost", port)
		}

		var bindings string
		var err error
		if _obfuscated, _ := ctx.Value("obfuscated").(bool); !_obfuscated {
			bindings, err = appBindings.ToJSON()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			appBindings.DB().UpdateObfuscatedCallMap()
		}
		assets, err := assetserver.NewAssetServerMainPage(bindings, appoptions, ctx.Value("assetdir") != nil, myLogger, wailsruntime.RuntimeAssetsBundle)
		if err != nil {
			log.Fatal(err)
		}
		f.assets = assets

		//go f.startRequestProcessor()
	}

	//go f.startMessageProcessor()

	var _debug = ctx.Value("debug")
	var _devtoolsEnabled = ctx.Value("devtoolsEnabled")

	if _debug != nil {
		f.debug = _debug.(bool)
	}
	if _devtoolsEnabled != nil {
		f.devtoolsEnabled = _devtoolsEnabled.(bool)
	}

	//f.mainWindow = NewWindow(appoptions, f.debug, f.devtoolsEnabled)

	C.install_signal_handlers()

	appName := "WailsApp"
	if appoptions.Linux != nil {
		appName = appoptions.Linux.ProgramName
	}
	f.qApp = C.Application_new(C.CString(appName))

	//if appoptions.Linux != nil && appoptions.Linux.ProgramName != "" {
	//	prgname := C.CString(appoptions.Linux.ProgramName)
	//	C.g_set_prgname(prgname)
	//	C.free(unsafe.Pointer(prgname))
	//}

	//go f.startSecondInstanceProcessor()

	return f
}

// BrowserOpenURL implements frontend.Frontend.
func (f *Frontend) BrowserOpenURL(url string) {
	_ = browser.OpenURL(url)
}

// ClipboardGetText implements frontend.Frontend.
func (f *Frontend) ClipboardGetText() (string, error) {
	f.logger.Info("ClipboardGetText")
	cStr := C.Clipboard_get_text(f.qApp)
	return C.GoString(cStr), nil
}

// ClipboardSetText implements frontend.Frontend.
func (f *Frontend) ClipboardSetText(text string) error {
	f.logger.Info("ClipboardSetText")

	cStr := C.CString(text)
	defer C.cfree(unsafe.Pointer(cStr))

	C.Clipboard_set_text(f.qApp, cStr)
	return nil
}

// ExecJS implements frontend.Frontend.
func (f *Frontend) ExecJS(js string) {
	f.logger.Info("ExecJS")
	s := C.CString(js)
	defer C.cfree(unsafe.Pointer(s))
	C.WebEngineView_run_js(f.qWindow.web_engine_view, s)
}

// Hide implements frontend.Frontend.
func (f *Frontend) Hide() {
	f.logger.Info("Hide")
	C.Window_hide(f.qWindow.window)
}

// Show implements frontend.Frontend.
func (f *Frontend) Show() {
	f.logger.Info("Show")
	C.Window_show(f.qWindow.window)
}

// MenuSetApplicationMenu implements frontend.Frontend.
func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.logger.Info("MenuSetApplicationMenu")
}

// MenuUpdateApplicationMenu implements frontend.Frontend.
func (f *Frontend) MenuUpdateApplicationMenu() {
	f.logger.Info("MenuUpdateApplicationMenu")
}

// MessageDialog implements frontend.Frontend.
func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	f.logger.Info("MessageDialog")

	title := C.CString(dialogOptions.Title)
	message := C.CString(dialogOptions.Message)

	defer func() {
		C.cfree(unsafe.Pointer(title))
		C.cfree(unsafe.Pointer(message))
	}()

	var messageType C.int
	switch dialogOptions.Type {
	case frontend.InfoDialog:
		messageType = C.int(0)
	case frontend.ErrorDialog:
		messageType = C.int(1)
	case frontend.QuestionDialog:
		messageType = C.int(2)
	case frontend.WarningDialog:
		messageType = C.int(3)
	}

	ret := C.Window_run_message_dialog(f.qWindow.window, messageType, title, message)
	result := C.GoString(ret)
	f.logger.Info("Message dialog returned code %s", result)

	return result, nil
}

func (f *Frontend) WindowPrint() {
	f.logger.Info("WindowPrint")
	// TODO: webView->page()->printToPdf(fileName);
	//f.ExecJS("window.print();")
}

type EventNotify struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

// Notify implements frontend.Frontend.
func (f *Frontend) Notify(name string, data ...interface{}) {
	notification := EventNotify{
		Name: name,
		Data: data,
	}
	payload, err := json.Marshal(notification)
	if err != nil {
		f.logger.Error(err.Error())
		return
	}
	f.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) openFileDialogCommon(directory bool, multiple bool, dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	j, err := json.Marshal(dialogOptions)
	if err != nil {
		f.logger.Error("Failed to marshal dialogOptions %+v", err)
		return []string{}, err
	}

	s := C.CString(string(j))
	defer C.cfree(unsafe.Pointer(s))

	isDirectory := 0
	if directory {
		isDirectory = 1
	}

	isMultiple := 0
	if multiple {
		isMultiple = 1
	}

	res := C.GoString(C.Window_open_file_dialog(f.qWindow.window, C.int(isDirectory), C.int(isMultiple), s))

	var files []string
	if err := json.Unmarshal([]byte(res), &files); err != nil {
		f.logger.Error("Failed to unmarshal file dialog result %s", err)
		return []string{}, err
	}

	return files, nil
}

// OpenDirectoryDialog implements frontend.Frontend.
func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	f.logger.Info("OpenDirectoryDialog")
	files, err := f.openFileDialogCommon(true, false, dialogOptions)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", nil
	}
	return files[0], nil
}

// OpenFileDialog implements frontend.Frontend.
func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	f.logger.Info("OpenFileDialog")
	files, err := f.openFileDialogCommon(false, false, dialogOptions)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", nil
	}
	return files[0], nil
}

// OpenMultipleFilesDialog implements frontend.Frontend.
func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	f.logger.Info("OpenMultipleFilesDialog")
	return f.openFileDialogCommon(false, true, dialogOptions)
}

// Quit implements frontend.Frontend.
func (f *Frontend) Quit() {
	f.logger.Info("Quit")
	C.Application_quit(f.qApp)
}

// Run implements frontend.Frontend.
func (f *Frontend) Run(ctx context.Context) error {
	f.ctx = ctx

	f.logger.Info("Run")

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()

	//if f.frontendOptions.SingleInstanceLock != nil {
	//	SetupSingleInstance(f.frontendOptions.SingleInstanceLock.UniqueId)
	//}
	//
	//f.mainWindow.Run(f.startURL.String())

	// TODO: Whats up with this?
	if f.startURL.Scheme == "wails" {
		f.startURL.Scheme = "http"
	}

	f.logger.Info("Creating window with url %s", f.startURL)

	f.qWindow = C.Window_new(f.qApp, C.CString(f.startURL.String()))

	return nil
}

// RunMainLoop implements frontend.Frontend.
func (f *Frontend) RunMainLoop() {
	f.logger.Info("RunMainLoop")

	time.Sleep(1 * time.Second)
	_ = f.ClipboardSetText("Foo")
	txt, err := f.ClipboardGetText()
	//res, err := f.OpenMultipleFilesDialog(frontend.OpenDialogOptions{
	//	Title:            "Title",
	//	DefaultDirectory: "/home/ben/Code",
	//	DefaultFilename:  "foo.txt",
	//})
	f.logger.Info("Got clipboard result %s - %s", txt, err)

	<-exitCh

	f.logger.Info("Qt App exited")
}

// SaveFileDialog implements frontend.Frontend.
func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	f.logger.Info("SaveFileDialog")
	return "", nil
}

// ScreenGetAll implements frontend.Frontend.
func (f *Frontend) ScreenGetAll() ([]frontend.Screen, error) {
	f.logger.Info("ScreenGetAll")
	screensJson := C.GoString(C.Application_get_screens(f.qApp))
	var screens []frontend.Screen
	if err := json.Unmarshal([]byte(screensJson), &screens); err != nil {
		f.logger.Error("Failed to unmarshal screens: %s", err)
		return screens, err
	}
	f.logger.Info("Got screens: %+v", screens)
	return screens, nil
}

// WindowCenter implements frontend.Frontend.
func (f *Frontend) WindowCenter() {
	f.logger.Info("WindowCenter")
}

// WindowClose implements frontend.Frontend.
func (f *Frontend) WindowClose() {
	f.logger.Info("WindowClose")
	C.Window_close(f.qWindow.window)
}

// WindowFullscreen implements frontend.Frontend.
func (f *Frontend) WindowFullscreen() {
	f.logger.Info("WindowFullscreen")
	C.Window_fullscreen(f.qWindow.window)
}

// WindowGetPosition implements frontend.Frontend.
func (f *Frontend) WindowGetPosition() (int, int) {
	f.logger.Info("WindowGetPosition")
	return 0, 0
}

// WindowGetSize implements frontend.Frontend.
func (f *Frontend) WindowGetSize() (int, int) {
	f.logger.Info("WindowGetSize")
	return 1, 1
}

// WindowHide implements frontend.Frontend.
func (f *Frontend) WindowHide() {
	f.logger.Info("WindowHide")
	C.Window_hide(f.qWindow.window)
}

// WindowIsFullscreen implements frontend.Frontend.
func (f *Frontend) WindowIsFullscreen() bool {
	f.logger.Info("WindowIsFullscreen")
	return false
}

// WindowIsMaximised implements frontend.Frontend.
func (f *Frontend) WindowIsMaximised() bool {
	f.logger.Info("WindowIsMaximized")
	return false
}

// WindowIsMinimised implements frontend.Frontend.
func (f *Frontend) WindowIsMinimised() bool {
	f.logger.Info("WindowIsMinimized")
	return false
}

// WindowIsNormal implements frontend.Frontend.
func (f *Frontend) WindowIsNormal() bool {
	f.logger.Info("WindowIsNormal")
	return false
}

// WindowMaximise implements frontend.Frontend.
func (f *Frontend) WindowMaximise() {
	f.logger.Info("WindowMaximize")
	C.Window_maximize(f.qWindow.window)
}

// WindowMinimise implements frontend.Frontend.
func (f *Frontend) WindowMinimise() {
	f.logger.Info("WindowMinimize")
	C.Window_hide(f.qWindow.window)
}

// WindowReload implements frontend.Frontend.
func (f *Frontend) WindowReload() {
	f.logger.Info("WindowReload")
	//C.WebEngineView_reload(f.qWindow.web_engine_view)
	f.ExecJS("runtime.WindowReload();")
}

// WindowReloadApp implements frontend.Frontend.
func (f *Frontend) WindowReloadApp() {
	f.logger.Info("WindowReloadApp")
	C.WebEngineView_reload(f.qWindow.web_engine_view)
}

// WindowSetAlwaysOnTop implements frontend.Frontend.
func (f *Frontend) WindowSetAlwaysOnTop(b bool) {
	f.logger.Info("WindowSetAlwaysOnTop")
}

// WindowSetBackgroundColour implements frontend.Frontend.
func (f *Frontend) WindowSetBackgroundColour(col *options.RGBA) {
	f.logger.Info("WindowSetBackgroundColour")
}

// WindowSetDarkTheme implements frontend.Frontend.
func (f *Frontend) WindowSetDarkTheme() {
	f.logger.Info("WindowSetDarkTheme")
}

// WindowSetLightTheme implements frontend.Frontend.
func (f *Frontend) WindowSetLightTheme() {
	f.logger.Info("WindowSetLightTheme")
}

// WindowSetMaxSize implements frontend.Frontend.
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	f.logger.Info("WindowSetMaxSize")
}

// WindowSetMinSize implements frontend.Frontend.
func (f *Frontend) WindowSetMinSize(width int, height int) {
	f.logger.Info("WindowSetMinSize")
}

// WindowSetPosition implements frontend.Frontend.
func (f *Frontend) WindowSetPosition(x int, y int) {
	f.logger.Info("WindowSetPosition")
}

// WindowSetSize implements frontend.Frontend.
func (f *Frontend) WindowSetSize(width int, height int) {
	f.logger.Info("WindowSetSize")
	C.Window_resize(f.qWindow.window, C.int(width), C.int(height))
}

// WindowSetSystemDefaultTheme implements frontend.Frontend.
func (f *Frontend) WindowSetSystemDefaultTheme() {
	f.logger.Info("WindowSetSystemDefaultTheme")
}

// WindowSetTitle implements frontend.Frontend.
func (f *Frontend) WindowSetTitle(title string) {
	f.logger.Info("WindowSetTitle")
	str := C.CString(title)
	defer C.cfree(unsafe.Pointer(str))
	C.Window_set_title(f.qWindow.window, str)
}

// WindowShow implements frontend.Frontend.
func (f *Frontend) WindowShow() {
	f.logger.Info("WindowShow")
	//f.WindowShow()
}

// WindowToggleMaximise implements frontend.Frontend.
func (f *Frontend) WindowToggleMaximise() {
	f.logger.Info("WindowToggleMaximize")
}

// WindowUnfullscreen implements frontend.Frontend.
func (f *Frontend) WindowUnfullscreen() {
	f.logger.Info("WindowUnfullscreen")
}

// WindowUnmaximise implements frontend.Frontend.
func (f *Frontend) WindowUnmaximise() {
	f.logger.Info("WindowUnmaximize")
}

// WindowUnminimise implements frontend.Frontend.
func (f *Frontend) WindowUnminimise() {
	f.logger.Info("WindowUnminimize")
}

var _ frontend.Frontend = &Frontend{}
