package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"embed"
	"encoding/json"
	"io/fs"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"unsafe"
)

//go:embed frontend
var embedded embed.FS

var start sync.Once
var requests = make(chan struct{}, 16)
var window *application.WebviewWindow

type GreetService struct{}

func (*GreetService) Greet(name string) string { return "Hello " + name + "!" }

//export WailsBridgeOpen
func WailsBridgeOpen(env, info unsafe.Pointer) unsafe.Pointer {
	return application.WailsOhNativeInit(env, info)
}

//export WailsBridgeResult
func WailsBridgeResult(env, info unsafe.Pointer) unsafe.Pointer {
	return application.WailsOhOnBridgeResult(env, info)
}

//export WailsBridgeClose
func WailsBridgeClose(env, info unsafe.Pointer) unsafe.Pointer {
	return application.WailsOhClose(env, info)
}

//export WailsStart
func WailsStart() {
	start.Do(func() {
		assets, err := fs.Sub(embedded, "frontend")
		if err != nil {
			panic(err)
		}
		app := application.New(application.Options{
			Name:     "Native HelloWails",
			Assets:   application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
			Services: []application.Service{application.NewService(&GreetService{})},
		})
		window = app.Window.NewWithOptions(application.WebviewWindowOptions{Name: "main"})
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for n := 1; ; n++ {
				select {
				case <-app.Context().Done():
					return
				case <-ticker.C:
					app.Event.Emit("hello:tick", n)
				}
			}
		}()
		if err := app.Run(); err != nil {
			panic(err)
		}
	})
}

//export WailsReady
func WailsReady() C.int {
	if application.OpenharmonyReady() {
		return 1
	}
	return 0
}

//export WailsAsset
func WailsAsset(path *C.char, length *C.size_t) *C.char {
	data, ok := application.WailsOhNativeServeAsset(C.GoString(path), "GET", "")
	if !ok {
		return nil
	}
	*length = C.size_t(len(data))
	return (*C.char)(C.CBytes(data))
}

//export WailsMime
func WailsMime(path *C.char) *C.char {
	name := C.GoString(path)
	if name == "/" {
		name = "/index.html"
	}
	return C.CString(application.WailsOhNativeGetAssetMimeType(name))
}

//export WailsRequest
func WailsRequest(id, payload *C.char) C.int {
	if !application.OpenharmonyReady() {
		return 0
	}
	select {
	case requests <- struct{}{}:
	default:
		return 0
	}
	requestID, call := C.GoString(id), C.GoString(payload)
	go func() {
		defer func() { <-requests }()
		response := application.WailsOhNativeHandleRuntimeCall(call)
		idJSON, _ := json.Marshal(requestID)
		payloadJSON, _ := json.Marshal(response)
		window.ExecJS("window.wailsOpenHarmonyResolve?.(" + string(idJSON) + "," + string(payloadJSON) + ")")
	}()
	return 1
}

//export WailsMessage
func WailsMessage(message *C.char) {
	value := C.GoString(message)
	go application.WailsOhNativeHandleMessage(value)
}

//export WailsRunMain
func WailsRunMain(id C.uint) { application.WailsOhNativeMainThreadCallback(uint(id)) }

func main() {}
