//go:build windows
// +build windows

package windows

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/leaanthony/go-webview2/pkg/edge"
	"github.com/leaanthony/winc"
	"github.com/leaanthony/winc/w32"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/common"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	chromium        *edge.Chromium
	debug           bool

	// Assets
	assets   *assetserver.DesktopAssetServer
	startURL string

	// main window handle
	mainWindow      *Window
	bindings        *binding.Bindings
	dispatcher      frontend.Dispatcher
	servingFromDisk bool

	hasStarted bool
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	result := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
		ctx:             ctx,
		startURL:        "http://wails.localhost/",
	}

	bindingsJSON, err := appBindings.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	_devServerURL := ctx.Value("devserverurl")
	if _devServerURL != nil {
		devServerURL := _devServerURL.(string)
		if len(devServerURL) > 0 && devServerURL != "http://localhost:34115" {
			result.startURL = devServerURL
			return result
		}
	}

	// Check if we have been given a directory to serve assets from.
	// If so, this means we are in dev mode and are serving assets off disk.
	// We indicate this through the `servingFromDisk` flag to ensure requests
	// aren't cached by WebView2 in dev mode

	_assetdir := ctx.Value("assetdir")
	if _assetdir != nil {
		result.servingFromDisk = true
	}

	assets, err := assetserver.NewDesktopAssetServer(ctx, appoptions.Assets, bindingsJSON, result.servingFromDisk)
	if err != nil {
		log.Fatal(err)
	}
	result.assets = assets

	return result
}

func (f *Frontend) WindowReload() {
	f.ExecJS("runtime.WindowReload();")
}

func (f *Frontend) Run(ctx context.Context) error {

	f.ctx = context.WithValue(ctx, "frontend", f)

	mainWindow := NewWindow(nil, f.frontendOptions)
	f.mainWindow = mainWindow

	var _debug = ctx.Value("debug")
	if _debug != nil {
		f.debug = _debug.(bool)
	}

	f.WindowCenter()
	f.setupChromium()

	f.mainWindow.notifyParentWindowPositionChanged = f.chromium.NotifyParentWindowPositionChanged

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

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()
	mainWindow.Run()
	mainWindow.Close()
	return nil
}

func (f *Frontend) WindowCenter() {
	runtime.LockOSThread()
	f.mainWindow.Center()
}

func (f *Frontend) WindowSetPosition(x, y int) {
	runtime.LockOSThread()
	f.mainWindow.SetPos(x, y)
}
func (f *Frontend) WindowGetPosition() (int, int) {
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
	if f.frontendOptions.Frameless && f.frontendOptions.DisableResize == false {
		f.ExecJS("window.wails.flags.enableResize = false;")
	}
	f.mainWindow.Fullscreen()
}

func (f *Frontend) WindowUnfullscreen() {
	runtime.LockOSThread()
	if f.frontendOptions.Frameless && f.frontendOptions.DisableResize == false {
		f.ExecJS("window.wails.flags.enableResize = true;")
	}
	f.mainWindow.UnFullscreen()
}

func (f *Frontend) WindowShow() {
	runtime.LockOSThread()
	f.ShowWindow()
}

func (f *Frontend) WindowHide() {
	runtime.LockOSThread()
	f.mainWindow.Hide()
}
func (f *Frontend) WindowMaximise() {
	runtime.LockOSThread()
	if f.hasStarted {
		if !f.frontendOptions.DisableResize {
			f.mainWindow.Maximise()
		}
	} else {
		f.frontendOptions.WindowStartState = options.Maximised
	}
}
func (f *Frontend) WindowToggleMaximise() {
	runtime.LockOSThread()
	if !f.hasStarted {
		return
	}
	if f.mainWindow.IsMaximised() {
		f.WindowUnmaximise()
	} else {
		f.WindowMaximise()
	}
}

func (f *Frontend) WindowUnmaximise() {
	runtime.LockOSThread()
	f.mainWindow.Restore()
}
func (f *Frontend) WindowMinimise() {
	runtime.LockOSThread()
	if f.hasStarted {
		f.mainWindow.Minimise()
	} else {
		f.frontendOptions.WindowStartState = options.Minimised
	}
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

func (f *Frontend) WindowSetRGBA(col *options.RGBA) {
	runtime.LockOSThread()
	if col == nil {
		return
	}

	f.mainWindow.Invoke(func() {
		controller := f.chromium.GetController()
		controller2 := controller.GetICoreWebView2Controller2()

		backgroundCol := edge.COREWEBVIEW2_COLOR{
			A: col.A,
			R: col.R,
			G: col.G,
			B: col.B,
		}

		// Webview2 only has 0 and 255 as valid values.
		if backgroundCol.A > 0 && backgroundCol.A < 255 {
			backgroundCol.A = 255
		}

		if f.frontendOptions.Windows != nil && f.frontendOptions.Windows.WebviewIsTransparent {
			backgroundCol.A = 0
		}

		err := controller2.PutDefaultBackgroundColor(backgroundCol)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func (f *Frontend) Quit() {
	if f.frontendOptions.OnBeforeClose != nil && f.frontendOptions.OnBeforeClose(f.ctx) {
		return
	}
	// Exit must be called on the Main-Thread. It calls PostQuitMessage which sends the WM_QUIT message to the thread's
	// message queue and our message queue runs on the Main-Thread.
	f.mainWindow.Invoke(winc.Exit)
}

func (f *Frontend) setupChromium() {
	chromium := edge.NewChromium()
	f.chromium = chromium
	if opts := f.frontendOptions.Windows; opts != nil && opts.WebviewUserDataPath != "" {
		chromium.DataPath = opts.WebviewUserDataPath
	}
	chromium.MessageCallback = f.processMessage
	chromium.WebResourceRequestedCallback = f.processRequest
	chromium.NavigationCompletedCallback = f.navigationCompleted
	chromium.AcceleratorKeyCallback = func(vkey uint) bool {
		w32.PostMessage(f.mainWindow.Handle(), w32.WM_KEYDOWN, uintptr(vkey), 0)
		return false
	}
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
	err = settings.PutIsStatusBarEnabled(false)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreBrowserAcceleratorKeysEnabled(false)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutIsSwipeNavigationEnabled(false)
	if err != nil {
		log.Fatal(err)
	}

	// Setup focus event handler
	onFocus := f.mainWindow.OnSetFocus()
	onFocus.Bind(f.onFocus)

	// Set background colour
	f.WindowSetRGBA(f.frontendOptions.RGBA)

	chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
	chromium.Navigate(f.startURL)
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
	f.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) processRequest(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
	//Get the request
	uri, _ := req.GetUri()

	res, err := common.ProcessRequest(uri, f.assets, "http", "wails.localhost")
	if err == common.ErrUnexpectedScheme {
		// In this case we should let the WebView2 handle the request with its default handler
		return
	} else if err == common.ErrUnexpectedHost {
		// This means file:// to something other than wails, should we prevent this?
		// Maybe we should introduce an AllowList for explicitly allowing schemes and hosts, this could also be interesting
		// for all other platforms to improve security.
		return // Let WebView2 handle the request with its default handler
	} else if err != nil {
		path := strings.Replace(uri, "http://wails.localhost", "", 1)
		f.logger.Error("Error processing request '%s': %s (HttpResponse=%s)", path, err, res)
	}

	headers := []string{}
	if mimeType := res.MimeType; mimeType != "" {
		headers = append(headers, "Content-Type: "+mimeType)
	}
	content := res.Body
	if content != nil && f.servingFromDisk {
		headers = append(headers, "Pragma: no-cache")
	}

	env := f.chromium.Environment()
	response, err := env.CreateWebResourceResponse(content, res.StatusCode, res.StatusText(), strings.Join(headers, "\n"))
	if err != nil {
		f.logger.Error("CreateWebResourceResponse Error: %s", err)
		return
	}
	defer response.Release()

	// Send response back
	err = args.PutResponse(response)
	if err != nil {
		f.logger.Error("PutResponse Error: %s", err)
		return
	}
}

var edgeMap = map[string]uintptr{
	"n-resize":  w32.HTTOP,
	"ne-resize": w32.HTTOPRIGHT,
	"e-resize":  w32.HTRIGHT,
	"se-resize": w32.HTBOTTOMRIGHT,
	"s-resize":  w32.HTBOTTOM,
	"sw-resize": w32.HTBOTTOMLEFT,
	"w-resize":  w32.HTLEFT,
	"nw-resize": w32.HTTOPLEFT,
}

func (f *Frontend) processMessage(message string) {
	if message == "drag" {
		if !f.mainWindow.IsFullScreen() {
			err := f.startDrag()
			if err != nil {
				f.logger.Error(err.Error())
			}
		}
		return
	}
	if strings.HasPrefix(message, "resize:") {
		if !f.mainWindow.IsFullScreen() {
			sl := strings.Split(message, ":")
			if len(sl) != 2 {
				f.logger.Info("Unknown message returned from dispatcher: %+v", message)
				return
			}
			edge := edgeMap[sl[1]]
			err := f.startResize(edge)
			if err != nil {
				f.logger.Error(err.Error())
			}
		}
		return
	}

	go func() {
		result, err := f.dispatcher.ProcessMessage(message, f)
		if err != nil {
			f.logger.Error(err.Error())
			f.Callback(result)
			return
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
	}()
}

func (f *Frontend) Callback(message string) {
	f.mainWindow.Invoke(func() {
		f.chromium.Eval(`window.wails.Callback(` + strconv.Quote(message) + `);`)
	})
}

func (f *Frontend) startDrag() error {
	if !w32.ReleaseCapture() {
		return fmt.Errorf("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until dragging has been finished.
	w32.PostMessage(f.mainWindow.Handle(), w32.WM_NCLBUTTONDOWN, w32.HTCAPTION, 0)
	return nil
}

func (f *Frontend) startResize(border uintptr) error {
	if !w32.ReleaseCapture() {
		return fmt.Errorf("unable to release mouse capture")
	}
	// Use PostMessage because we don't want to block the caller until resizing has been finished.
	w32.PostMessage(f.mainWindow.Handle(), w32.WM_NCLBUTTONDOWN, border, 0)
	return nil
}

func (f *Frontend) ExecJS(js string) {
	f.mainWindow.Invoke(func() {
		f.chromium.Eval(js)
	})
}

func (f *Frontend) navigationCompleted(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationCompletedEventArgs) {
	if f.frontendOptions.OnDomReady != nil {
		go f.frontendOptions.OnDomReady(f.ctx)
	}

	if f.frontendOptions.Frameless && f.frontendOptions.DisableResize == false {
		f.ExecJS("window.wails.flags.enableResize = true;")
	}

	if f.hasStarted {
		return
	}
	f.hasStarted = true

	// Hack to make it visible: https://github.com/MicrosoftEdge/WebView2Feedback/issues/1077#issuecomment-825375026
	err := f.chromium.Hide()
	if err != nil {
		log.Fatal(err)
	}
	err = f.chromium.Show()
	if err != nil {
		log.Fatal(err)
	}

	if f.frontendOptions.StartHidden {
		return
	}

	switch f.frontendOptions.WindowStartState {
	case options.Maximised:
		if !f.frontendOptions.DisableResize {
			f.mainWindow.Maximise()
		} else {
			f.mainWindow.Show()
		}
		f.ShowWindow()

	case options.Minimised:
		f.mainWindow.Minimise()
	case options.Fullscreen:
		f.mainWindow.Fullscreen()
		f.ShowWindow()
	default:
		if f.frontendOptions.Fullscreen {
			f.mainWindow.Fullscreen()
		}
		f.ShowWindow()
	}

}

func (f *Frontend) ShowWindow() {
	f.mainWindow.Invoke(func() {
		f.mainWindow.Restore()
		w32.SetForegroundWindow(f.mainWindow.Handle())
		w32.SetFocus(f.mainWindow.Handle())
	})

}

func (f *Frontend) onFocus(arg *winc.Event) {
	f.chromium.Focus()
}
