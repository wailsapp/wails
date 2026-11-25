package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"embed"
	"encoding/json"
	"io/fs"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend
var embedded embed.FS

var start sync.Once
var requests = make(chan struct{}, 16)
var window *application.WebviewWindow

type GreetService struct{}

func (*GreetService) Greet(name string) string { return "Hello " + name + "!" }

// LifecycleProbe completes after the host replaces the requesting document.
func (*GreetService) LifecycleProbe(phase string) string {
	if phase == "old" {
		time.Sleep(6 * time.Second)
	} else {
		time.Sleep(8 * time.Second)
	}
	return phase
}

type NativeProbe struct{}

// AccentColor exercises the ArkTS resource-manager bridge with a real reply.
func (*NativeProbe) AccentColor() string { return application.Get().Env.GetAccentColor() }

// AppInfo exercises an OpenHarmony info getter with a real reply.
func (*NativeProbe) AppInfo() string { return application.Openharmony.AppInfoJSON() }

// SecureRoundTrip exercises asset storage through the ArkTS bridge.
func (*NativeProbe) SecureRoundTrip() string {
	application.Openharmony.SecureSet("m3-native-key", "m3-native-value")
	value, _ := application.Openharmony.SecureGet("m3-native-key")
	application.Openharmony.SecureDelete("m3-native-key")
	return value
}

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
			Services: []application.Service{application.NewService(&GreetService{}), application.NewService(&NativeProbe{})},
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
	session, requestID, ok := strings.Cut(C.GoString(id), ":")
	if !ok || session == "" || requestID == "" {
		return 0
	}
	select {
	case requests <- struct{}{}:
	default:
		return 0
	}
	call := C.GoString(payload)
	go func() {
		defer func() { <-requests }()
		response := application.WailsOhNativeHandleRuntimeCall(call)
		sessionJSON, _ := json.Marshal(session)
		idJSON, _ := json.Marshal(requestID)
		payloadJSON, _ := json.Marshal(response)
		// ExecJS resolves the current document, not the one that started the call.
		window.ExecJS("if(window.wailsOpenHarmonySession===" + string(sessionJSON) + "){window.wailsOpenHarmonyResolve?.(" + string(idJSON) + "," + string(payloadJSON) + ")}")
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

//export WailsDialogResult
func WailsDialogResult(id C.uint, buttonIndex C.int) {
	application.WailsOhNativeDialogCallback(uint(id), int(buttonIndex))
}

//export WailsFilePickerResult
func WailsFilePickerResult(id C.uint, path *C.char) {
	application.WailsOhNativeFilePickerResult(uint(id), C.GoString(path))
}

//export WailsFilePickerDone
func WailsFilePickerDone(id C.uint) { application.WailsOhNativeFilePickerDone(uint(id)) }

//export WailsAbilityStarted
func WailsAbilityStarted() { application.WailsOhNativeOnStart() }

//export WailsAbilityResumed
func WailsAbilityResumed() { application.WailsOhNativeOnResume() }

//export WailsAbilityPaused
func WailsAbilityPaused() { application.WailsOhNativeOnPause() }

//export WailsAbilityStopped
func WailsAbilityStopped() { application.WailsOhNativeOnStop() }

//export WailsLowMemory
func WailsLowMemory() { application.WailsOhNativeOnLowMemory() }

//export WailsEmitSystemEvent
func WailsEmitSystemEvent(name *C.char, payload *C.char) {
	application.WailsOhNativeEmitSystemEvent(C.GoString(name), C.GoString(payload))
}

//export WailsEmitEvent
func WailsEmitEvent(name *C.char, payload *C.char) {
	application.WailsOhNativeEmitEvent(C.GoString(name), C.GoString(payload))
}

//export WailsPageStarted
func WailsPageStarted(url *C.char) { application.WailsOhNativeOnPageStarted(C.GoString(url)) }

//export WailsPageFinished
func WailsPageFinished(url *C.char) { application.WailsOhNativeOnPageFinished(C.GoString(url)) }

//export WailsShutdown
func WailsShutdown() { application.WailsOhNativeShutdown() }

func main() {}
