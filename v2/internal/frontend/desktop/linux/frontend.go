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
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"unsafe"

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
	debug           bool

	// Assets
	assets   *assetserver.DesktopAssetServer
	startURL string

	// main window handle
	mainWindow      *Window
	bindings        *binding.Bindings
	dispatcher      frontend.Dispatcher
	servingFromDisk bool
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
		startURL:        "file://wails/",
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
	// aren't cached by webkit.

	_assetdir := ctx.Value("assetdir")
	if _assetdir != nil {
		result.servingFromDisk = true
	}

	assets, err := assetserver.NewDesktopAssetServer(ctx, appoptions.Assets, bindingsJSON, result.servingFromDisk)
	if err != nil {
		log.Fatal(err)
	}
	result.assets = assets

	go result.startMessageProcessor()
	go result.startRequestProcessor()

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

func (f *Frontend) Run(ctx context.Context) error {

	f.ctx = context.WithValue(ctx, "frontend", f)

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()

	f.mainWindow.Run()

	return nil
}

func (f *Frontend) WindowCenter() {
	f.mainWindow.Center()
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

func (f *Frontend) WindowSetRGBA(col *options.RGBA) {
	if col == nil {
		return
	}
	f.mainWindow.SetRGBA(col.R, col.G, col.B, col.A)
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
	}
}

//export processURLRequest
func processURLRequest(request unsafe.Pointer) {
	requestBuffer <- request
}

func (f *Frontend) processRequest(request unsafe.Pointer) {
	req := (*C.WebKitURISchemeRequest)(request)
	uri := C.webkit_uri_scheme_request_get_uri(req)
	goURI := C.GoString(uri)

	res, err := common.ProcessRequest(goURI, f.assets, "wails", "", "null")
	if err != nil {
		f.logger.Error("Error processing request '%s': %s (HttpResponse=%s)", goURI, err, res)
	}

	if code := res.StatusCode; code != http.StatusOK {
		message := C.CString(res.StatusText())
		gerr := C.g_error_new_literal(C.g_quark_from_string(message), C.int(code), message)
		C.webkit_uri_scheme_request_finish_error(req, gerr)
		C.g_error_free(gerr)
		C.free(unsafe.Pointer(message))
		return
	}

	var cContent unsafe.Pointer
	bodyLen := len(res.Body)
	var cLen C.long
	if bodyLen > 0 {
		cContent = C.malloc(C.ulong(bodyLen))
		if cContent != nil {
			C.memcpy(cContent, unsafe.Pointer(&res.Body[0]), C.size_t(bodyLen))
			cLen = C.long(bodyLen)
		}
	}

	cMimeType := C.CString(res.MimeType)
	defer C.free(unsafe.Pointer(cMimeType))

	stream := C.g_memory_input_stream_new_from_data(
		cContent,
		cLen,
		(*[0]byte)(C.free))
	C.webkit_uri_scheme_request_finish(req, stream, cLen, cMimeType)
	C.g_object_unref(C.gpointer(stream))
}
