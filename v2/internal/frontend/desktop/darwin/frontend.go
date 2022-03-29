//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

#include <stdlib.h>
*/
import "C"
import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"strconv"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/common"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type request struct {
	url *C.char
	ctx unsafe.Pointer
}

var messageBuffer = make(chan string, 100)
var requestBuffer = make(chan *request, 100)
var callbackBuffer = make(chan uint, 10)

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	debug           bool

	// Assets
	assets *assetserver.DesktopAssetServer

	// main window handle
	mainWindow      *Window
	bindings        *binding.Bindings
	dispatcher      frontend.Dispatcher
	servingFromDisk bool
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

	bindingsJSON, err := appBindings.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	assets, err := assetserver.NewDesktopAssetServer(ctx, appoptions.Assets, bindingsJSON, result.servingFromDisk)
	if err != nil {
		log.Fatal(err)
	}
	result.assets = assets

	go result.startMessageProcessor()
	go result.startRequestProcessor()
	go result.startCallbackProcessor()

	return result
}

func (f *Frontend) startMessageProcessor() {
	for message := range messageBuffer {
		f.processMessage(message)
	}
}
func (f *Frontend) startRequestProcessor() {
	for request := range requestBuffer {
		f.processRequest(request)
	}
}
func (f *Frontend) startCallbackProcessor() {
	for callback := range callbackBuffer {
		err := f.handleCallback(callback)
		if err != nil {
			println(err.Error())
		}
	}
}

func (f *Frontend) WindowReload() {
	f.ExecJS("runtime.WindowReload();")
}

func (f *Frontend) Run(ctx context.Context) error {

	f.ctx = context.WithValue(ctx, "frontend", f)

	var _debug = ctx.Value("debug")
	if _debug != nil {
		f.debug = _debug.(bool)
	}

	mainWindow := NewWindow(f.frontendOptions, f.debug)
	f.mainWindow = mainWindow
	f.mainWindow.Center()

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()
	mainWindow.Run()
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
	f.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) processMessage(message string) {

	if message == "DomReady" {
		if f.frontendOptions.OnDomReady != nil {
			f.frontendOptions.OnDomReady(f.ctx)
		}
		return
	}

	//if strings.HasPrefix(message, "systemevent:") {
	//	f.processSystemEvent(message)
	//	return
	//}

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

func (f *Frontend) ExecJS(js string) {
	f.mainWindow.ExecJS(js)
}

func (f *Frontend) processRequest(r *request) {
	uri := C.GoString(r.url)

	res, err := common.ProcessRequest(uri, f.assets, "wails", "wails")
	if err != nil {
		f.logger.Error("Error processing request '%s': %s (HttpResponse=%s)", uri, err, res)
	}

	var content unsafe.Pointer
	var contentLen int
	if _contents := res.Body; _contents != nil {
		content = unsafe.Pointer(&_contents[0])
		contentLen = len(_contents)
	}
	mimetype := C.CString(res.MimeType)
	defer C.free(unsafe.Pointer(mimetype))

	C.ProcessURLResponse(r.ctx, r.url, C.int(res.StatusCode), mimetype, content, C.int(contentLen))
}

//func (f *Frontend) processSystemEvent(message string) {
//	sl := strings.Split(message, ":")
//	if len(sl) != 2 {
//		f.logger.Error("Invalid system message: %s", message)
//		return
//	}
//	switch sl[1] {
//	case "fullscreen":
//		f.mainWindow.DisableSizeConstraints()
//	case "unfullscreen":
//		f.mainWindow.EnableSizeConstraints()
//	default:
//		f.logger.Error("Unknown system message: %s", message)
//	}
//}

//export processMessage
func processMessage(message *C.char) {
	goMessage := C.GoString(message)
	messageBuffer <- goMessage
}

//export processURLRequest
func processURLRequest(ctx unsafe.Pointer, url *C.char) {
	requestBuffer <- &request{
		url: url,
		ctx: ctx,
	}
}

//export processCallback
func processCallback(callbackID uint) {
	callbackBuffer <- callbackID
}
