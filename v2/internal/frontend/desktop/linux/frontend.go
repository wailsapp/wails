//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"

extern void callDispatchedMethod(int id);

static inline void processDispatchID(gpointer id) {
    callDispatchedMethod(GPOINTER_TO_INT(id));
}

static void gtkDispatch(int id) {
	g_idle_add((GSourceFunc)processDispatchID, GINT_TO_POINTER(id));
}
*/
import "C"
import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
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
	mainWindow *Window
	// minWidth, minHeight, maxWidth, maxHeight int
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

	assets, err := assetserver.NewDesktopAssetServer(ctx, appoptions.Assets, bindingsJSON)
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

func (f *Frontend) WindowSetPos(x, y int) {
	f.mainWindow.SetPos(x, y)
}
func (f *Frontend) WindowGetPos() (int, int) {
	return f.mainWindow.Pos()
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

func (f *Frontend) WindowUnFullscreen() {
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
	f.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) processMessage(message string) {
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
	f.dispatch(func() {
		f.mainWindow.StartDrag()
	})
}

func (f *Frontend) ExecJS(js string) {
	f.dispatch(func() {
		f.mainWindow.ExecJS(js)
	})
}

func (f *Frontend) dispatch(fn func()) {
	dispatchCallbackLock.Lock()
	id := 0
	for fn := dispatchCallbacks[id]; fn != nil; id++ {
	}
	dispatchCallbacks[id] = fn
	dispatchCallbackLock.Unlock()
	C.gtkDispatch(C.int(id))
}

var messageBuffer = make(chan string, 100)

//export processMessage
func processMessage(message *C.char) {
	goMessage := C.GoString(message)
	messageBuffer <- goMessage
}

// Map of functions passed to dispatch()
var dispatchCallbacks = make(map[int]func())
var dispatchCallbackLock sync.Mutex

//export callDispatchedMethod
func callDispatchedMethod(cid C.int) {
	id := int(cid)
	fn := dispatchCallbacks[id]
	if fn != nil {
		fn()
		dispatchCallbackLock.Lock()
		delete(dispatchCallbacks, id)
		dispatchCallbackLock.Unlock()
	} else {
		println("Error: No dispatch method with id", id, cid)
	}
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

	file, match, err := common.TranslateUriToFile(goURI, "wails", "")
	if err != nil {
		// TODO Handle errors
		return
	} else if !match {
		// This should never happen on linux, because we get only called for wails://
		panic("Unexpected host for request on wails:// scheme")
	}

	// Load file from asset store
	content, mimeType, err := f.assets.Load(file)

	// TODO How to return 404/500 errors to webkit?
	//if err != nil {
	//if os.IsNotExist(err) {
	//	f.dispatch(func() {
	//		message := C.CString("not found")
	//		defer C.free(unsafe.Pointer(message))
	//		C.webkit_uri_scheme_request_finish_error(req, C.g_error_new_literal(C.G_FILE_ERROR_NOENT, C.int(404), message))
	//	})
	//} else {
	//	err = fmt.Errorf("Error processing request %s: %w", uri, err)
	//	f.logger.Error(err.Error())
	//	message := C.CString("internal server error")
	//	defer C.free(unsafe.Pointer(message))
	//	C.webkit_uri_scheme_request_finish_error(req, C.g_error_new_literal(C.G_FILE_ERROR_NOENT, C.int(500), message))
	//}
	//return
	//}

	cContent := C.CString(string(content))
	defer C.free(unsafe.Pointer(cContent))
	cMimeType := C.CString(mimeType)
	defer C.free(unsafe.Pointer(cMimeType))
	cLen := C.long(C.strlen(cContent))
	stream := C.g_memory_input_stream_new_from_data(
		unsafe.Pointer(C.g_strdup(cContent)),
		cLen,
		(*[0]byte)(C.g_free))
	C.webkit_uri_scheme_request_finish(req, stream, cLen, cMimeType)
	C.g_object_unref(C.gpointer(stream))
}
