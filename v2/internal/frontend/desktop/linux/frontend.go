//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"text/template"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

const startURL = "wails://wails/"

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	debug           bool

	// Assets
	assets   *assetserver.AssetServer
	startURL *url.URL

	// main window handle
	mainWindow *Window
	bindings   *binding.Bindings
	dispatcher frontend.Dispatcher
}

func init() {
	runtime.LockOSThread()
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	// Set GDK_BACKEND=x11 to prevent warnings
	_ = os.Setenv("GDK_BACKEND", "x11")

	result := &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
		bindings:        appBindings,
		dispatcher:      dispatcher,
		ctx:             ctx,
	}
	result.startURL, _ = url.Parse(startURL)

	if _starturl, _ := ctx.Value("starturl").(*url.URL); _starturl != nil {
		result.startURL = _starturl
	} else {
		bindingsJSON, err := appBindings.ToJSON()
		if err != nil {
			log.Fatal(err)
		}

		assets, err := assetserver.NewAssetServer(ctx, appoptions, bindingsJSON)
		if err != nil {
			log.Fatal(err)
		}
		result.assets = assets

		// Start 10 processors to handle requests in parallel
		for i := 0; i < 10; i++ {
			go result.startRequestProcessor()
		}
	}

	go result.startMessageProcessor()

	C.gtk_init(nil, nil)

	var _debug = ctx.Value("debug")
	if _debug != nil {
		result.debug = _debug.(bool)
	}
	result.mainWindow = NewWindow(appoptions, result.debug)

	return result
}

func (f *Frontend) startMessageProcessor() {
	for message := range messageBuffer {
		f.processMessage(message)
	}
}

func (f *Frontend) WindowReload() {
	f.ExecJS("runtime.WindowReload();")
}

func (f *Frontend) WindowSetSystemDefaultTheme() {
	return
}

func (f *Frontend) WindowSetLightTheme() {
	return
}

func (f *Frontend) WindowSetDarkTheme() {
	return
}

func (f *Frontend) Run(ctx context.Context) error {

	f.ctx = context.WithValue(ctx, "frontend", f)

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()

	f.mainWindow.Run(f.startURL.String())

	return nil
}

func (f *Frontend) WindowCenter() {
	f.mainWindow.Center()
}

func (f *Frontend) WindowSetAlwaysOnTop(b bool) {
	f.mainWindow.SetKeepAbove(b)
}

func (f *Frontend) WindowSetPosition(x, y int) {
	f.mainWindow.SetPosition(x, y)
}
func (f *Frontend) WindowGetPosition() (int, int) {
	return f.mainWindow.GetPosition()
}

func (f *Frontend) WindowSetSize(width, height int) {
	f.mainWindow.SetSize(width, height)
}

func (f *Frontend) WindowGetSize() (int, int) {
	return f.mainWindow.Size()
}

func (f *Frontend) WindowSetTitle(title string) {
	f.mainWindow.SetTitle(title)
}

func (f *Frontend) WindowFullscreen() {
	f.mainWindow.Fullscreen()
}

func (f *Frontend) WindowUnfullscreen() {
	f.mainWindow.UnFullscreen()
}

func (f *Frontend) WindowReloadApp() {
	f.ExecJS(fmt.Sprintf("window.location.href = '%s';", f.startURL))
}

func (f *Frontend) WindowShow() {
	f.mainWindow.Show()
}

func (f *Frontend) WindowHide() {
	f.mainWindow.Hide()
}
func (f *Frontend) WindowMaximise() {
	f.mainWindow.Maximise()
}
func (f *Frontend) WindowToggleMaximise() {
	f.mainWindow.ToggleMaximise()
}
func (f *Frontend) WindowUnmaximise() {
	f.mainWindow.UnMaximise()
}
func (f *Frontend) WindowMinimise() {
	f.mainWindow.Minimise()
}
func (f *Frontend) WindowUnminimise() {
	f.mainWindow.UnMinimise()
}

func (f *Frontend) WindowSetMinSize(width int, height int) {
	f.mainWindow.SetMinSize(width, height)
}
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	f.mainWindow.SetMaxSize(width, height)
}

func (f *Frontend) WindowSetBackgroundColour(col *options.RGBA) {
	if col == nil {
		return
	}
	f.mainWindow.SetBackgroundColour(col.R, col.G, col.B, col.A)
}

func (f *Frontend) Quit() {
	if f.frontendOptions.OnBeforeClose != nil && f.frontendOptions.OnBeforeClose(f.ctx) {
		return
	}
	f.mainWindow.Quit()
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
	f.mainWindow.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) processMessage(message string) {
	if message == "DomReady" {
		if f.frontendOptions.OnDomReady != nil {
			f.frontendOptions.OnDomReady(f.ctx)
		}
		return
	}

	if message == "drag" {
		if !f.mainWindow.IsFullScreen() {
			f.startDrag()
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
	f.ExecJS(`window.wails.Callback(` + strconv.Quote(message) + `);`)
}

func (f *Frontend) startDrag() {
	f.mainWindow.StartDrag()
}

func (f *Frontend) ExecJS(js string) {
	f.mainWindow.ExecJS(js)
}

var messageBuffer = make(chan string, 100)

//export processMessage
func processMessage(message *C.char) {
	goMessage := C.GoString(message)
	messageBuffer <- goMessage
}

var requestBuffer = make(chan unsafe.Pointer, 100)

func (f *Frontend) startRequestProcessor() {
	for request := range requestBuffer {
		f.processRequest(request)
		C.g_object_unref(C.gpointer(request))
	}
}

//export processURLRequest
func processURLRequest(request unsafe.Pointer) {
	// Increment reference counter to allow async processing, will be decremented after the processing
	// has been finished by a worker.
	C.g_object_ref(C.gpointer(request))
	requestBuffer <- request
}

func (f *Frontend) processRequest(request unsafe.Pointer) {
	req := (*C.WebKitURISchemeRequest)(request)
	uri := C.webkit_uri_scheme_request_get_uri(req)
	goURI := C.GoString(uri)

	// WebKitGTK stable < 2.36 API does not support request method, request headers and request.
	// Apart from request bodies, this is only available beginning with 2.36: https://webkitgtk.org/reference/webkit2gtk/stable/WebKitURISchemeResponse.html
	rw := &webKitResponseWriter{req: req}
	defer rw.Close()

	f.assets.ProcessHTTPRequest(
		goURI,
		rw,
		func() (*http.Request, error) {
			req, err := http.NewRequest(http.MethodGet, goURI, nil)
			if err != nil {
				return nil, err
			}

			if req.URL.Host != f.startURL.Host {
				if req.Body != nil {
					req.Body.Close()
				}

				return nil, fmt.Errorf("Expected host '%d' in request, but was '%s'", f.startURL.Host, req.URL.Host)
			}

			return req, nil
		})

}
