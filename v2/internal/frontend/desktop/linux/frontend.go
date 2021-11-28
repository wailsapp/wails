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
	gdk_threads_add_idle((GSourceFunc)processDispatchID, GINT_TO_POINTER(id));
}

*/
import "C"
import (
	"context"
	"encoding/json"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"unsafe"
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
	//minWidth, minHeight, maxWidth, maxHeight int
	bindings        *binding.Bindings
	dispatcher      frontend.Dispatcher
	servingFromDisk bool
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {

	// Set GDK_BACKEND=x11 to prevent warnings
	os.Setenv("GDK_BACKEND", "x11")

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
	////dCenter()
}

func (f *Frontend) WindowSetPos(x, y int) {
	////dSetPos(x, y)
}
func (f *Frontend) WindowGetPos() (int, int) {
	////return dPos()
	return 0, 0
}

func (f *Frontend) WindowSetSize(width, height int) {
	////dSetSize(width, height)
}

func (f *Frontend) WindowGetSize() (int, int) {
	////return dSize()
	return 0, 0
}

func (f *Frontend) WindowSetTitle(title string) {
	////dSetText(title)
}

func (f *Frontend) WindowFullscreen() {
	////dSetMaxSize(0, 0)
	////dSetMinSize(0, 0)
	////dFullscreen()
}

func (f *Frontend) WindowUnFullscreen() {
	////dUnFullscreen()
	////dSetMaxSize(f.maxWidth, f.maxHeight)
	////dSetMinSize(f.minWidth, f.minHeight)
}

func (f *Frontend) WindowShow() {
	////dShow()
}

func (f *Frontend) WindowHide() {
	////dHide()
}
func (f *Frontend) WindowMaximise() {
	////dMaximise()
}
func (f *Frontend) WindowUnmaximise() {
	////dUnMaximise()
}
func (f *Frontend) WindowMinimise() {
	//dMinimise()
}
func (f *Frontend) WindowUnminimise() {
	//dUnMinimise()
}

func (f *Frontend) WindowSetMinSize(width int, height int) {
	//f.minWidth = width
	////f.minHeight = height
	//dSetMinSize(width, height)
}
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	//f.maxWidth = width
	////f.maxHeight = height
	//dSetMaxSize(width, height)
}

func (f *Frontend) WindowSetRGBA(col *options.RGBA) {
	if col == nil {
		return
	}
	//
	//f.gtkWindow.Dispatch(func() {
	//	controller := f.chromium.GetController()
	//	controller2 := controller.GetICoreWebView2Controller2()
	//
	//	backgroundCol := edge.COREWEBVIEW2_COLOR{
	//		A: col.A,
	//		R: col.R,
	//		G: col.G,
	//		B: col.B,
	//	}
	//
	//	// Webview2 only has 0 and 255 as valid values.
	//	if backgroundCol.A > 0 && backgroundCol.A < 255 {
	//		backgroundCol.A = 255
	//	}
	//
	//	if f.frontendOptions.Windows != nil && f.frontendOptions.Windows.WebviewIsTransparent {
	//		backgroundCol.A = 0
	//	}
	//
	//	err := controller2.PutDefaultBackgroundColor(backgroundCol)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//})
}

func (f *Frontend) Quit() {
	//winc.Exit()
}

//func (f *Frontend) setupChromium() {
//	chromium := edge.NewChromium()
//	f.chromium = chromium
//	chromium.MessageCallback = f.processMessage
//	chromium.WebResourceRequestedCallback = f.processRequest
//	chromium.NavigationCompletedCallback = f.navigationCompleted
//	chromium.AcceleratorKeyCallback = func(vkey uint) bool {
//		w32.PostMessage(f.gtkWindow.Handle(), w32.WM_KEYDOWN, uintptr(vkey), 0)
//		return false
//	}
//	chromium.Embed(f.gtkWindow.Handle())
//	chromium.Resize()
//	settings, err := chromium.GetSettings()
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutAreDefaultContextMenusEnabled(f.debug)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutAreDevToolsEnabled(f.debug)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutIsZoomControlEnabled(false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutIsStatusBarEnabled(false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutAreBrowserAcceleratorKeysEnabled(false)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = settings.PutIsSwipeNavigationEnabled(false)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Set background colour
//	f.WindowSetRGBA(f.frontendOptions.RGBA)
//
//	chromium.SetGlobalPermission(edge.CoreWebView2PermissionStateAllow)
//	chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
//	chromium.Navigate(f.startURL)
//}

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

//func (f *Frontend) processRequest(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
//	//Get the request
//	uri, _ := req.GetUri()
//
//	// Translate URI
//	uri = strings.TrimPrefix(uri, "file://wails")
//	if !strings.HasPrefix(uri, "/") {
//		return
//	}
//
//	// Load file from asset store
//	content, mimeType, err := f.assets.Load(uri)
//	if err != nil {
//		return
//	}
//
//	env := f.chromium.Environment()
//	headers := "Content-Type: " + mimeType
//	if f.servingFromDisk {
//		headers += "\nPragma: no-cache"
//	}
//	response, err := env.CreateWebResourceResponse(content, 200, "OK", headers)
//	if err != nil {
//		return
//	}
//	// Send response back
//	err = args.PutResponse(response)
//	if err != nil {
//		return
//	}
//	return
//}

func (f *Frontend) processMessage(message string) {
	if message == "drag" {
		//if !f.mainWindow.IsFullScreen() {
		err := f.startDrag()
		if err != nil {
			f.logger.Error(err.Error())
		}
		//}
		return
	}
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
}

func (f *Frontend) Callback(message string) {
	f.ExecJS(`window.wails.Callback(` + strconv.Quote(message) + `);`)
}

func (f *Frontend) startDrag() error {
	//if !w32.ReleaseCapture() {
	//	return fmt.Errorf("unable to release mouse capture")
	//}
	//w32.SendMessage(f.gtkWindow.Handle(), w32.WM_NCLBUTTONDOWN, w32.HTCAPTION, 0)
	return nil
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

//func (f *Frontend) navigationCompleted(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationCompletedEventArgs) {
//	if f.frontendOptions.OnDomReady != nil {
//		go f.frontendOptions.OnDomReady(f.ctx)
//	}
//
//	// If you want to start hidden, return
//	if f.frontendOptions.StartHidden {
//		return
//	}
//
//	f.gtkWindow.Show()
//
//}

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
		go fn()
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

	// Translate URI
	goURI = strings.TrimPrefix(goURI, "wails://")
	if !strings.HasPrefix(goURI, "/") {
		return
	}

	// Load file from asset store
	content, mimeType, err := f.assets.Load(goURI)
	if err != nil {
		return
	}

	cContent := C.CString(string(content))
	defer C.free(unsafe.Pointer(cContent))
	cMimeType := C.CString(mimeType)
	defer C.free(unsafe.Pointer(cMimeType))
	var cLen C.long = (C.long)(C.strlen(cContent))
	stream := C.g_memory_input_stream_new_from_data(unsafe.Pointer(cContent), cLen, nil)
	C.webkit_uri_scheme_request_finish(req, stream, cLen, cMimeType)
}
