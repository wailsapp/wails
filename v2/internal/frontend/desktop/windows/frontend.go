//go:build windows

package windows

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jchv/go-webview2/pkg/edge"
	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"log"
	"runtime"
	"strconv"
	"strings"
)

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	chromium        *edge.Chromium
	debug           bool

	// Assets
	assets *assetserver.DesktopAssetServer

	// main window handle
	mainWindow                               *Window
	minWidth, minHeight, maxWidth, maxHeight int
	bindings                                 *binding.Bindings
	dispatcher                               frontend.Dispatcher
	servingFromDisk                          bool
}

func (f *Frontend) Run(ctx context.Context) error {

	mainWindow := NewWindow(nil, f.frontendOptions)
	f.mainWindow = mainWindow

	var _debug = ctx.Value("debug")
	if _debug != nil {
		f.debug = _debug.(bool)
	}

	f.WindowCenter()

	if !f.frontendOptions.StartHidden {
		mainWindow.Show()
	}

	f.setupChromium()

	mainWindow.OnSize().Bind(func(arg *winc.Event) {
		f.chromium.Resize()
	})

	mainWindow.OnClose().Bind(func(arg *winc.Event) {
		if f.frontendOptions.HideWindowOnClose {
			f.WindowHide()
		} else {
			f.Quit()
		}
	})

	// TODO: Move this into a callback from frontend
	go func() {
		ctx := context.WithValue(ctx, "frontend", f)
		f.frontendOptions.Startup(ctx)
	}()

	mainWindow.Run()
	return nil
}

func (f *Frontend) WindowCenter() {
	runtime.LockOSThread()
	f.mainWindow.Center()
}

func (f *Frontend) WindowSetPos(x, y int) {
	runtime.LockOSThread()
	f.mainWindow.SetPos(x, y)
}
func (f *Frontend) WindowGetPos() (int, int) {
	runtime.LockOSThread()
	return f.mainWindow.Pos()
}

func (f *Frontend) WindowSetSize(width, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetSize(width, height)
}

func (f *Frontend) WindowGetSize() (int, int) {
	runtime.LockOSThread()
	return f.mainWindow.Size()
}

func (f *Frontend) WindowSetTitle(title string) {
	runtime.LockOSThread()
	f.mainWindow.SetText(title)
}

func (f *Frontend) WindowFullscreen() {
	runtime.LockOSThread()
	f.mainWindow.Fullscreen()
}

func (f *Frontend) WindowUnFullscreen() {
	runtime.LockOSThread()
	f.mainWindow.UnFullscreen()
}

func (f *Frontend) WindowShow() {
	runtime.LockOSThread()
	f.mainWindow.Show()
}

func (f *Frontend) WindowHide() {
	runtime.LockOSThread()
	f.mainWindow.Hide()
}
func (f *Frontend) WindowMaximise() {
	runtime.LockOSThread()
	f.mainWindow.Maximise()
}
func (f *Frontend) WindowUnmaximise() {
	runtime.LockOSThread()
	f.mainWindow.Restore()
}
func (f *Frontend) WindowMinimise() {
	runtime.LockOSThread()
	f.mainWindow.Minimise()
}
func (f *Frontend) WindowUnminimise() {
	runtime.LockOSThread()
	f.mainWindow.Restore()
}

func (f *Frontend) WindowSetMinSize(width int, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetMinSize(width, height)
}
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetMaxSize(width, height)
}

func (f *Frontend) WindowSetColour(colour int) {
	runtime.LockOSThread()
	// TODO: Set webview2 background to this colour
}

func (f *Frontend) Quit() {
	winc.Exit()
}

func (f *Frontend) setupChromium() {
	chromium := edge.NewChromium()
	chromium.MessageCallback = f.processMessage
	chromium.WebResourceRequestedCallback = f.processRequest
	chromium.Embed(f.mainWindow.Handle())
	chromium.Resize()
	settings, err := chromium.GetSettings()
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreDefaultContextMenusEnabled(f.debug)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreDevToolsEnabled(f.debug)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutIsZoomControlEnabled(false)
	if err != nil {
		log.Fatal(err)
	}
	//c2, err := chromium.GetWebView2Controller2()
	//err = c2.PutDefaultBackgroundColor(edge.COREWEBVIEW2_COLOR{R: 255, G: 0, B: 0, A: 255})
	//if err != nil {
	//	log.Fatal(err)
	//}
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
	chromium.Navigate("file://wails/")
	f.chromium = chromium
}

type EventNotify struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

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
	f.mainWindow.Dispatch(func() {
		f.chromium.Eval(`window.wails.EventsNotify('` + string(payload) + `');`)
	})
}

func (f *Frontend) processRequest(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
	//Get the request
	uri, _ := req.GetUri()

	// Translate URI
	uri = strings.TrimPrefix(uri, "file://wails")
	if !strings.HasPrefix(uri, "/") {
		return
	}

	// Load file from asset store
	content, mimeType, err := f.assets.Load(uri)
	if err != nil {
		return
	}

	env := f.chromium.Environment()
	headers := "Content-Type: " + mimeType
	if f.servingFromDisk {
		headers += "\nPragma: no-cache"
	}
	response, err := env.CreateWebResourceResponse(content, 200, "OK", headers)
	if err != nil {
		return
	}
	// Send response back
	err = args.PutResponse(response)
	if err != nil {
		return
	}
	return
}

func (f *Frontend) processMessage(message string) {
	if message == "drag" {
		err := f.startDrag()
		if err != nil {
			f.logger.Error(err.Error())
		}
		return
	}
	result, err := f.dispatcher.ProcessMessage(message)
	if err != nil {
		f.logger.Error(err.Error())
	}
	if result == "" {
		return
	}

	switch result[0] {
	case 'c':
		// Callback from a method call
		f.Callback(result[1:])
	default:
		f.logger.Info("Unknown message returned from dispatcher: %+v", result)
	}
}

func (f *Frontend) Callback(message string) {
	f.mainWindow.Dispatch(func() {
		f.chromium.Eval(`window.wails.Callback(` + strconv.Quote(message) + `);`)
	})
}

func (f *Frontend) startDrag() error {
	if !w32.ReleaseCapture() {
		return fmt.Errorf("unable to release mouse capture")
	}
	w32.SendMessage(f.mainWindow.Handle(), w32.WM_NCLBUTTONDOWN, w32.HTCAPTION, 0)
	return nil
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	result := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
		ctx:             ctx,
	}

	// Check if we have been given a directory to serve assets from.
	// If so, this means we are in dev mode and are serving assets off disk.
	// We indicate this through the `servingFromDisk` flag to ensure requests
	// aren't cached by WebView2 in dev mode
	_assetdir := ctx.Value("assetdir")
	if _assetdir != nil {
		result.servingFromDisk = true
	}

	if appoptions.Assets != nil {
		bindingsJSON, err := appBindings.ToJSON()
		if err != nil {
			log.Fatal(err)
		}
		assets, err := assetserver.NewDesktopAssetServer(ctx, *appoptions.Assets, bindingsJSON)
		if err != nil {
			log.Fatal(err)
		}
		result.assets = assets
	}

	return result
}
