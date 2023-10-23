//go:build windows
// +build windows

package windows

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/bep/debounce"
	"github.com/wailsapp/go-webview2/pkg/edge"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/win32"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
	wailsruntime "github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/pkg/assetserver"
	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const startURL = "http://wails.localhost/"

var secondInstanceBuffer = make(chan options.SecondInstanceData, 1)

type Screen = frontend.Screen

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	chromium        *edge.Chromium
	debug           bool
	devtoolsEnabled bool

	// Assets
	assets   *assetserver.AssetServer
	startURL *url.URL

	// main window handle
	mainWindow *Window
	bindings   *binding.Bindings
	dispatcher frontend.Dispatcher

	hasStarted bool

	// Windows build number
	versionInfo     *operatingsystem.WindowsVersionInfo
	resizeDebouncer func(f func())
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	// Get Windows build number
	versionInfo, _ := operatingsystem.GetWindowsVersionInfo()

	result := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
		ctx:             ctx,
		versionInfo:     versionInfo,
	}

	if appoptions.Windows != nil {
		if appoptions.Windows.ResizeDebounceMS > 0 {
			result.resizeDebouncer = debounce.New(time.Duration(appoptions.Windows.ResizeDebounceMS) * time.Millisecond)
		}
	}

	// We currently can't use wails://wails/ as other platforms do, therefore we map the assets sever onto the following url.
	result.startURL, _ = url.Parse(startURL)

	if _starturl, _ := ctx.Value("starturl").(*url.URL); _starturl != nil {
		result.startURL = _starturl
		return result
	}

	if port, _ := ctx.Value("assetserverport").(string); port != "" {
		result.startURL.Host = net.JoinHostPort(result.startURL.Host, port)
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
	result.assets = assets

	go result.startSecondInstanceProcessor()

	return result
}

func (f *Frontend) WindowReload() {
	f.ExecJS("runtime.WindowReload();")
}

func (f *Frontend) WindowSetSystemDefaultTheme() {
	f.mainWindow.SetTheme(windows.SystemDefault)
}

func (f *Frontend) WindowSetLightTheme() {
	f.mainWindow.SetTheme(windows.Light)
}

func (f *Frontend) WindowSetDarkTheme() {
	f.mainWindow.SetTheme(windows.Dark)
}

func (f *Frontend) Run(ctx context.Context) error {
	f.ctx = ctx

	f.chromium = edge.NewChromium()

	if f.frontendOptions.SingleInstanceLock != nil {
		SetupSingleInstance(f.frontendOptions.SingleInstanceLock.UniqueId)
	}

	mainWindow := NewWindow(nil, f.frontendOptions, f.versionInfo, f.chromium)
	f.mainWindow = mainWindow

	var _debug = ctx.Value("debug")
	var _devtoolsEnabled = ctx.Value("devtoolsEnabled")

	if _debug != nil {
		f.debug = _debug.(bool)
	}
	if _devtoolsEnabled != nil {
		f.devtoolsEnabled = _devtoolsEnabled.(bool)
	}

	f.WindowCenter()
	f.setupChromium()

	mainWindow.OnSize().Bind(func(arg *winc.Event) {
		if f.frontendOptions.Frameless {
			// If the window is frameless and we are minimizing, then we need to suppress the Resize on the
			// WebView2. If we don't do this, restoring does not work as expected and first restores with some wrong
			// size during the restore animation and only fully renders when the animation is done. This highly
			// depends on the content in the WebView, see https://github.com/wailsapp/wails/issues/1319
			event, _ := arg.Data.(*winc.SizeEventData)
			if event != nil && event.Type == w32.SIZE_MINIMIZED {
				return
			}
		}

		if f.resizeDebouncer != nil {
			f.resizeDebouncer(func() {
				f.mainWindow.Invoke(func() {
					f.chromium.Resize()
				})
			})
		} else {
			f.chromium.Resize()
		}
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
	mainWindow.UpdateTheme()
	return nil
}

func (f *Frontend) WindowClose() {
	if f.mainWindow != nil {
		f.mainWindow.Close()
	}
}

func (f *Frontend) RunMainLoop() {
	_ = winc.RunMainLoop()
}

func (f *Frontend) WindowCenter() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.Center()
}

func (f *Frontend) WindowSetAlwaysOnTop(b bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetAlwaysOnTop(b)
}

func (f *Frontend) WindowSetPosition(x, y int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetPos(x, y)
}
func (f *Frontend) WindowGetPosition() (int, int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return f.mainWindow.Pos()
}

func (f *Frontend) WindowSetSize(width, height int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetSize(width, height)
}

func (f *Frontend) WindowGetSize() (int, int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	return f.mainWindow.Size()
}

func (f *Frontend) WindowSetTitle(title string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetText(title)
}

func (f *Frontend) WindowFullscreen() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if f.frontendOptions.Frameless && f.frontendOptions.DisableResize == false {
		f.ExecJS("window.wails.flags.enableResize = false;")
	}
	f.mainWindow.Fullscreen()
}

func (f *Frontend) WindowReloadApp() {
	f.ExecJS(fmt.Sprintf("window.location.href = '%s';", f.startURL))
}

func (f *Frontend) WindowUnfullscreen() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if f.frontendOptions.Frameless && f.frontendOptions.DisableResize == false {
		f.ExecJS("window.wails.flags.enableResize = true;")
	}
	f.mainWindow.UnFullscreen()
}

func (f *Frontend) WindowShow() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.ShowWindow()
}

func (f *Frontend) WindowHide() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.Hide()
}

func (f *Frontend) WindowMaximise() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
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
	defer runtime.UnlockOSThread()
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
	defer runtime.UnlockOSThread()
	if f.mainWindow.Form.IsFullScreen() {
		return
	}
	f.mainWindow.Restore()
}

func (f *Frontend) WindowMinimise() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if f.hasStarted {
		f.mainWindow.Minimise()
	} else {
		f.frontendOptions.WindowStartState = options.Minimised
	}
}

func (f *Frontend) WindowUnminimise() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if f.mainWindow.Form.IsFullScreen() {
		return
	}
	f.mainWindow.Restore()
}

func (f *Frontend) WindowSetMinSize(width int, height int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetMinSize(width, height)
}
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	f.mainWindow.SetMaxSize(width, height)
}

func (f *Frontend) WindowSetBackgroundColour(col *options.RGBA) {
	if col == nil {
		return
	}

	f.mainWindow.Invoke(func() {
		win32.SetBackgroundColour(f.mainWindow.Handle(), col.R, col.G, col.B)

		controller := f.chromium.GetController()
		controller2 := controller.GetICoreWebView2Controller2()

		backgroundCol := edge.COREWEBVIEW2_COLOR{
			A: col.A,
			R: col.R,
			G: col.G,
			B: col.B,
		}

		// WebView2 only has 0 and 255 as valid values.
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

func (f *Frontend) ScreenGetAll() ([]Screen, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	screens := []Screen{}
	err := error(nil)
	f.mainWindow.Invoke(func() {
		screens, err = GetAllScreens(f.mainWindow.Handle())
		wg.Done()

	})
	wg.Wait()
	return screens, err
}

func (f *Frontend) Show() {
	f.mainWindow.Show()
}

func (f *Frontend) Hide() {
	f.mainWindow.Hide()
}

func (f *Frontend) WindowIsMaximised() bool {
	return f.mainWindow.IsMaximised()
}

func (f *Frontend) WindowIsMinimised() bool {
	return f.mainWindow.IsMinimised()
}

func (f *Frontend) WindowIsNormal() bool {
	return f.mainWindow.IsNormal()
}

func (f *Frontend) WindowIsFullscreen() bool {
	return f.mainWindow.IsFullScreen()
}

func (f *Frontend) Quit() {
	if f.frontendOptions.OnBeforeClose != nil && f.frontendOptions.OnBeforeClose(f.ctx) {
		return
	}
	// Exit must be called on the Main-Thread. It calls PostQuitMessage which sends the WM_QUIT message to the thread's
	// message queue and our message queue runs on the Main-Thread.
	f.mainWindow.Invoke(winc.Exit)
}

func (f *Frontend) WindowPrint() {
	f.ExecJS("window.print();")
}

func (f *Frontend) setupChromium() {
	chromium := f.chromium

	disableFeatues := []string{}
	if !f.frontendOptions.EnableFraudulentWebsiteDetection {
		disableFeatues = append(disableFeatues, "msSmartScreenProtection")
	}

	if opts := f.frontendOptions.Windows; opts != nil {
		chromium.DataPath = opts.WebviewUserDataPath
		chromium.BrowserPath = opts.WebviewBrowserPath

		if opts.WebviewGpuIsDisabled {
			chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, "--disable-gpu")
		}
		if opts.WebviewDisableRendererCodeIntegrity {
			disableFeatues = append(disableFeatues, "RendererCodeIntegrity")
		}
	}

	if len(disableFeatues) > 0 {
		arg := fmt.Sprintf("--disable-features=%s", strings.Join(disableFeatues, ","))
		chromium.AdditionalBrowserArgs = append(chromium.AdditionalBrowserArgs, arg)
	}

	chromium.MessageCallback = f.processMessage
	chromium.WebResourceRequestedCallback = f.processRequest
	chromium.NavigationCompletedCallback = f.navigationCompleted
	chromium.AcceleratorKeyCallback = func(vkey uint) bool {
		if vkey == w32.VK_F12 && f.devtoolsEnabled {
			var keyState [256]byte
			if w32.GetKeyboardState(keyState[:]) {
				// Check if CTRL is pressed
				if keyState[w32.VK_CONTROL]&0x80 != 0 && keyState[w32.VK_SHIFT]&0x80 != 0 {
					chromium.OpenDevToolsWindow()
					return true
				}
			} else {
				f.logger.Error("Call to GetKeyboardState failed")
			}
		}
		w32.PostMessage(f.mainWindow.Handle(), w32.WM_KEYDOWN, uintptr(vkey), 0)
		return false
	}
	chromium.ProcessFailedCallback = func(sender *edge.ICoreWebView2, args *edge.ICoreWebView2ProcessFailedEventArgs) {
		kind, err := args.GetProcessFailedKind()
		if err != nil {
			f.logger.Error("GetProcessFailedKind: %s", err)
			return
		}

		f.logger.Error("WebVie2wProcess failed with kind %d", kind)
		switch kind {
		case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED:
			// => The app has to recreate a new WebView to recover from this failure.
			messages := windows.DefaultMessages()
			if f.frontendOptions.Windows != nil && f.frontendOptions.Windows.Messages != nil {
				messages = f.frontendOptions.Windows.Messages
			}
			winc.Errorf(f.mainWindow, messages.WebView2ProcessCrash)
			os.Exit(-1)
		case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED,
			edge.COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED:
			// => A new render process is created automatically and navigated to an error page.
			// => Make sure that the error page is shown.
			if !f.hasStarted {
				// NavgiationCompleted didn't come in, make sure the chromium is shown
				chromium.Show()
			}
			if !f.mainWindow.hasBeenShown {
				// The window has never been shown, make sure to show it
				f.ShowWindow()
			}
		}
	}

	chromium.Embed(f.mainWindow.Handle())

	if chromium.HasCapability(edge.SwipeNavigation) {
		swipeGesturesEnabled := f.frontendOptions.Windows != nil && f.frontendOptions.Windows.EnableSwipeGestures
		err := chromium.PutIsSwipeNavigationEnabled(swipeGesturesEnabled)
		if err != nil {
			log.Fatal(err)
		}
	}
	chromium.Resize()
	settings, err := chromium.GetSettings()
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreDefaultContextMenusEnabled(f.debug || f.frontendOptions.EnableDefaultContextMenu)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreDevToolsEnabled(f.devtoolsEnabled)
	if err != nil {
		log.Fatal(err)
	}

	if opts := f.frontendOptions.Windows; opts != nil {
		if opts.ZoomFactor > 0.0 {
			chromium.PutZoomFactor(opts.ZoomFactor)
		}
		err = settings.PutIsZoomControlEnabled(opts.IsZoomControlEnabled)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = settings.PutIsStatusBarEnabled(false)
	if err != nil {
		log.Fatal(err)
	}
	err = settings.PutAreBrowserAcceleratorKeysEnabled(false)
	if err != nil {
		log.Fatal(err)
	}

	if f.debug && f.frontendOptions.Debug.OpenInspectorOnStartup {
		chromium.OpenDevToolsWindow()
	}

	// Setup focus event handler
	onFocus := f.mainWindow.OnSetFocus()
	onFocus.Bind(f.onFocus)

	// Set background colour
	f.WindowSetBackgroundColour(f.frontendOptions.BackgroundColour)

	chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
	chromium.Navigate(f.startURL.String())
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
	// Setting the UserAgent on the CoreWebView2Settings clears the whole default UserAgent of the Edge browser, but
	// we want to just append our ApplicationIdentifier. So we adjust the UserAgent for every request.
	if reqHeaders, err := req.GetHeaders(); err == nil {
		useragent, _ := reqHeaders.GetHeader(assetserver.HeaderUserAgent)
		useragent = strings.Join([]string{useragent, assetserver.WailsUserAgentValue}, " ")
		reqHeaders.SetHeader(assetserver.HeaderUserAgent, useragent)
		reqHeaders.Release()
	}

	if f.assets == nil {
		// We are using the devServer let the WebView2 handle the request with its default handler
		return
	}

	//Get the request
	uri, _ := req.GetUri()
	reqUri, err := url.ParseRequestURI(uri)
	if err != nil {
		f.logger.Error("Unable to parse equest uri %s: %s", uri, err)
		return
	}

	if reqUri.Scheme != f.startURL.Scheme {
		// Let the WebView2 handle the request with its default handler
		return
	} else if reqUri.Host != f.startURL.Host {
		// Let the WebView2 handle the request with its default handler
		return
	}

	webviewRequest, err := webview.NewRequest(
		f.chromium.Environment(),
		args,
		func(fn func()) {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			if f.mainWindow.InvokeRequired() {
				var wg sync.WaitGroup
				wg.Add(1)
				f.mainWindow.Invoke(func() {
					fn()
					wg.Done()
				})
				wg.Wait()
			} else {
				fn()
			}
		})

	if err != nil {
		f.logger.Error("%s: NewRequest failed: %s", uri, err)
		return
	}

	f.assets.ServeWebViewRequest(webviewRequest)
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

	if message == "runtime:ready" {
		cmd := fmt.Sprintf("window.wails.setCSSDragProperties('%s', '%s');", f.frontendOptions.CSSDragProperty, f.frontendOptions.CSSDragValue)
		f.ExecJS(cmd)
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
	escaped, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	f.mainWindow.Invoke(func() {
		f.chromium.Eval(`window.wails.Callback(` + string(escaped) + `);`)
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
			win32.ShowWindowMaximised(f.mainWindow.Handle())
		} else {
			win32.ShowWindow(f.mainWindow.Handle())
		}
	case options.Minimised:
		win32.ShowWindowMinimised(f.mainWindow.Handle())
	case options.Fullscreen:
		f.mainWindow.Fullscreen()
		win32.ShowWindow(f.mainWindow.Handle())
	default:
		if f.frontendOptions.Fullscreen {
			f.mainWindow.Fullscreen()
		}
		win32.ShowWindow(f.mainWindow.Handle())
	}

	f.mainWindow.hasBeenShown = true

}

func (f *Frontend) ShowWindow() {
	f.mainWindow.Invoke(func() {
		if !f.mainWindow.hasBeenShown {
			f.mainWindow.hasBeenShown = true
			switch f.frontendOptions.WindowStartState {
			case options.Maximised:
				if !f.frontendOptions.DisableResize {
					win32.ShowWindowMaximised(f.mainWindow.Handle())
				} else {
					win32.ShowWindow(f.mainWindow.Handle())
				}
			case options.Minimised:
				win32.RestoreWindow(f.mainWindow.Handle())
			case options.Fullscreen:
				f.mainWindow.Fullscreen()
				win32.ShowWindow(f.mainWindow.Handle())
			default:
				if f.frontendOptions.Fullscreen {
					f.mainWindow.Fullscreen()
				}
				win32.ShowWindow(f.mainWindow.Handle())
			}
		} else {
			if win32.IsWindowMinimised(f.mainWindow.Handle()) {
				win32.RestoreWindow(f.mainWindow.Handle())
			} else {
				win32.ShowWindow(f.mainWindow.Handle())
			}
		}
		w32.SetForegroundWindow(f.mainWindow.Handle())
		w32.SetFocus(f.mainWindow.Handle())
	})

}

func (f *Frontend) onFocus(arg *winc.Event) {
	f.chromium.Focus()
}

func (f *Frontend) startSecondInstanceProcessor() {
	for secondInstanceData := range secondInstanceBuffer {
		if f.frontendOptions.SingleInstanceLock != nil &&
			f.frontendOptions.SingleInstanceLock.OnSecondInstanceLaunch != nil {
			f.frontendOptions.SingleInstanceLock.OnSecondInstanceLaunch(secondInstanceData)
		}
	}
}
