//go:build qt
// +build qt

package qt

import (
	"context"
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	wailsruntime "github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/assetserver"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"log"
	"net"
	"net/url"
	"time"
	"unsafe"
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
	return "", nil
}

// ClipboardSetText implements frontend.Frontend.
func (f *Frontend) ClipboardSetText(text string) error {
	f.logger.Info("ClipboardSetText")
	return nil
}

// ExecJS implements frontend.Frontend.
func (f *Frontend) ExecJS(js string) {
	f.logger.Info("ExecJS")
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
	return "", nil
}

// Notify implements frontend.Frontend.
func (f *Frontend) Notify(name string, data ...interface{}) {
	f.logger.Info("Notify")
}

// OpenDirectoryDialog implements frontend.Frontend.
func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	f.logger.Info("OpenDirectoryDialog")
	return "", nil
}

// OpenFileDialog implements frontend.Frontend.
func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	f.logger.Info("OpenFileDialog")
	return "", nil
}

// OpenMultipleFilesDialog implements frontend.Frontend.
func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	f.logger.Info("OpenMultipleFilesDialog")
	return []string{}, nil
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

	time.Sleep(3 * time.Second)
	f.WindowSetTitle("New title")

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
	return []frontend.Screen{}, nil
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
	f.logger.Info("WindowHide")
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

// WindowPrint implements frontend.Frontend.
func (f *Frontend) WindowPrint() {
	f.logger.Info("WindowPrint")
}

// WindowReload implements frontend.Frontend.
func (f *Frontend) WindowReload() {
	f.logger.Info("WindowReload")
	C.WebEngineView_reload(f.qWindow.web_engine_view)
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
	defer C.bye(unsafe.Pointer(str))
	C.Window_set_title(f.qWindow.window, str)
}

// WindowShow implements frontend.Frontend.
func (f *Frontend) WindowShow() {
	f.logger.Info("WindowShow")
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
