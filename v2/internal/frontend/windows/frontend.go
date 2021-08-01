package windows

import (
	"context"
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
	"strings"
)

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	chromium        *edge.Chromium

	// Assets
	assets *assetserver.AssetServer

	// main window handle
	mainWindow                               *Window
	minWidth, minHeight, maxWidth, maxHeight int
	bindings                                 *binding.Bindings
	dispatcher                               frontend.Dispatcher
}

func (f *Frontend) Run(ctx context.Context) error {

	mainWindow := NewWindow(nil, f.frontendOptions)
	f.mainWindow = mainWindow

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

	winc.RunMainLoop()
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
	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
	chromium.Navigate("file://index.html")
	f.chromium = chromium
}

func (f *Frontend) processRequest(sender *edge.ICoreWebView2, args *edge.ICoreWebView2WebResourceRequestedEventArgs) uintptr {
	// Get the request
	requestObject, _ := args.GetRequest()
	uri, _ := requestObject.GetURI()

	// Translate URI
	uri = strings.TrimPrefix(uri, "file://index.html")
	if !strings.HasPrefix(uri, "/") {
		return 0
	}

	// Load file from asset store
	content, mimeType, err := f.assets.Load(uri)
	if err != nil {
		return 0
	}

	// Create stream for response
	stream, err := w32.SHCreateMemStream(content)
	if err != nil {
		log.Fatal(err)
	}
	env := f.chromium.Environment()
	var response *edge.ICoreWebView2WebResourceResponse
	err = env.CreateWebResourceResponse(stream, 200, "OK", "Content-Type: "+mimeType, &response)
	if err != nil {
		return 0
	}
	// Send response back
	err = args.PutResponse(response)
	if err != nil {
		return 0
	}
	return 0
}

func (f *Frontend) processMessage(message string) {
	err := f.dispatcher.ProcessMessage(message)
	if err != nil {
		// TODO: Work out what this means
		return
	}
}

func NewFrontend(appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	result := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
	}

	if appoptions.Windows.Assets != nil {
		assets, err := assetserver.NewAssetServer(*appoptions.Windows.Assets)
		if err != nil {
			log.Fatal(err)
		}
		result.assets = assets
	}

	return result
}
