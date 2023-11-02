//go:build qt
// +build qt

package qt

import (
	"context"
	"fmt"
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
func (*Frontend) BrowserOpenURL(url string) {
	_ = browser.OpenURL(url)
}

// ClipboardGetText implements frontend.Frontend.
func (*Frontend) ClipboardGetText() (string, error) {
	fmt.Println("ClipboardGetText")
	return "", nil
}

// ClipboardSetText implements frontend.Frontend.
func (*Frontend) ClipboardSetText(text string) error {
	fmt.Println("ClipboardSetText")
	return nil
}

// ExecJS implements frontend.Frontend.
func (*Frontend) ExecJS(js string) {
	fmt.Println("ExecJS")
}

// Hide implements frontend.Frontend.
func (*Frontend) Hide() {
	fmt.Println("Hide")
}

// MenuSetApplicationMenu implements frontend.Frontend.
func (*Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	fmt.Println("MenuSetApplicationMenu")
}

// MenuUpdateApplicationMenu implements frontend.Frontend.
func (*Frontend) MenuUpdateApplicationMenu() {
	fmt.Println("MenuUpdateApplicationMenu")
}

// MessageDialog implements frontend.Frontend.
func (*Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	fmt.Println("MessageDialog")
	return "", nil
}

// Notify implements frontend.Frontend.
func (*Frontend) Notify(name string, data ...interface{}) {
	fmt.Println("Notify")
}

// OpenDirectoryDialog implements frontend.Frontend.
func (*Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	fmt.Println("OpenDirectoryDialog")
	return "", nil
}

// OpenFileDialog implements frontend.Frontend.
func (*Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	fmt.Println("OpenFileDialog")
	return "", nil
}

// OpenMultipleFilesDialog implements frontend.Frontend.
func (*Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	fmt.Println("OpenMultipleFilesDialog")
	return []string{}, nil
}

// Quit implements frontend.Frontend.
func (*Frontend) Quit() {
	fmt.Println("Quit")
}

// Run implements frontend.Frontend.
func (f *Frontend) Run(ctx context.Context) error {
	f.ctx = ctx

	fmt.Println("Run")

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

	//<-ctx.Done()
	//return ctx.Err()
	return nil
}

// RunMainLoop implements frontend.Frontend.
func (f *Frontend) RunMainLoop() {
	fmt.Println("RunMainLoop")

	f.qWindow = C.Window_new(f.qApp)

	<-exitCh

	fmt.Println("Qt App exited")
}

// SaveFileDialog implements frontend.Frontend.
func (*Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	fmt.Println("SaveFileDialog")
	return "", nil
}

// ScreenGetAll implements frontend.Frontend.
func (*Frontend) ScreenGetAll() ([]frontend.Screen, error) {
	fmt.Println("ScreenGetAll")
	return []frontend.Screen{}, nil
}

// Show implements frontend.Frontend.
func (f *Frontend) Show() {
	fmt.Println("Show")
}

// WindowCenter implements frontend.Frontend.
func (*Frontend) WindowCenter() {
	fmt.Println("WindowCenter")
}

// WindowClose implements frontend.Frontend.
func (*Frontend) WindowClose() {
	fmt.Println("WindowClose")
}

// WindowFullscreen implements frontend.Frontend.
func (*Frontend) WindowFullscreen() {
	fmt.Println("WindowFullscreen")
}

// WindowGetPosition implements frontend.Frontend.
func (*Frontend) WindowGetPosition() (int, int) {
	fmt.Println("WindowGetPosition")
	return 0, 0
}

// WindowGetSize implements frontend.Frontend.
func (*Frontend) WindowGetSize() (int, int) {
	fmt.Println("WindowGetSize")
	return 1, 1
}

// WindowHide implements frontend.Frontend.
func (*Frontend) WindowHide() {
	fmt.Println("WindowHide")
}

// WindowIsFullscreen implements frontend.Frontend.
func (*Frontend) WindowIsFullscreen() bool {
	fmt.Println("WindowHide")
	return false
}

// WindowIsMaximised implements frontend.Frontend.
func (*Frontend) WindowIsMaximised() bool {
	fmt.Println("WindowIsMaximized")
	return false
}

// WindowIsMinimised implements frontend.Frontend.
func (*Frontend) WindowIsMinimised() bool {
	fmt.Println("WindowIsMinimized")
	return false
}

// WindowIsNormal implements frontend.Frontend.
func (*Frontend) WindowIsNormal() bool {
	fmt.Println("WindowIsNormal")
	return false
}

// WindowMaximise implements frontend.Frontend.
func (*Frontend) WindowMaximise() {
	fmt.Println("WindowMaximize")
}

// WindowMinimise implements frontend.Frontend.
func (*Frontend) WindowMinimise() {
	fmt.Println("WindowMinimize")
}

// WindowPrint implements frontend.Frontend.
func (*Frontend) WindowPrint() {
	fmt.Println("WindowPrint")
}

// WindowReload implements frontend.Frontend.
func (*Frontend) WindowReload() {
	fmt.Println("WindowReload")
}

// WindowReloadApp implements frontend.Frontend.
func (*Frontend) WindowReloadApp() {
	fmt.Println("WindowReloadApp")
}

// WindowSetAlwaysOnTop implements frontend.Frontend.
func (*Frontend) WindowSetAlwaysOnTop(b bool) {
	fmt.Println("WindowSetAlwaysOnTop")
}

// WindowSetBackgroundColour implements frontend.Frontend.
func (*Frontend) WindowSetBackgroundColour(col *options.RGBA) {
	fmt.Println("WindowSetBackgroundColour")
}

// WindowSetDarkTheme implements frontend.Frontend.
func (*Frontend) WindowSetDarkTheme() {
	fmt.Println("WindowSetDarkTheme")
}

// WindowSetLightTheme implements frontend.Frontend.
func (*Frontend) WindowSetLightTheme() {
	fmt.Println("WindowSetLightTheme")
}

// WindowSetMaxSize implements frontend.Frontend.
func (*Frontend) WindowSetMaxSize(width int, height int) {
	fmt.Println("WindowSetMaxSize")
}

// WindowSetMinSize implements frontend.Frontend.
func (*Frontend) WindowSetMinSize(width int, height int) {
	fmt.Println("WindowSetMinSize")
}

// WindowSetPosition implements frontend.Frontend.
func (*Frontend) WindowSetPosition(x int, y int) {
	fmt.Println("WindowSetPosition")
}

// WindowSetSize implements frontend.Frontend.
func (f *Frontend) WindowSetSize(width int, height int) {
	fmt.Println("WindowSetSize")
	C.Window_resize(f.qWindow.window, C.int(width), C.int(height))
}

// WindowSetSystemDefaultTheme implements frontend.Frontend.
func (*Frontend) WindowSetSystemDefaultTheme() {
	fmt.Println("WindowSetSystemDefaultTheme")
}

// WindowSetTitle implements frontend.Frontend.
func (*Frontend) WindowSetTitle(title string) {
	fmt.Println("WindowSetTitle")
}

// WindowShow implements frontend.Frontend.
func (*Frontend) WindowShow() {
	fmt.Println("WindowShow")
}

// WindowToggleMaximise implements frontend.Frontend.
func (*Frontend) WindowToggleMaximise() {
	fmt.Println("WindowToggleMaximize")
}

// WindowUnfullscreen implements frontend.Frontend.
func (*Frontend) WindowUnfullscreen() {
	fmt.Println("WindowUnfullscreen")
}

// WindowUnmaximise implements frontend.Frontend.
func (*Frontend) WindowUnmaximise() {
	fmt.Println("WindowUnmaximize")
}

// WindowUnminimise implements frontend.Frontend.
func (*Frontend) WindowUnminimise() {
	fmt.Println("WindowUnminimize")
}

var _ frontend.Frontend = &Frontend{}
