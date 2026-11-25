//go:build openharmony && cgo && !server

package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// This file is the OpenHarmony counterpart of the Android shell's
// main_android.go. It is injected into the application's main package through
// a Go build overlay (see OpenHarmonyOverlayGen), because in -buildmode=c-shared
// the application's main() is never invoked automatically.
//
// The exported functions below form the C ABI between the ArkTS shell and the
// Go backend. The shell's napi_init.cpp calls them through the generated
// cgo header; keep both sides in sync.

var openharmonyRequests = make(chan struct{}, 16)

func init() {
	// In -buildmode=c-shared the application's main() is never invoked
	// automatically, so the host starts it (mirrors RegisterAndroidMain).
	application.RegisterOpenharmonyMain(main)
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
	// The id carries the ArkWeb document session and the request id, so a
	// reply can be matched to the document that asked for it.
	session, requestID, ok := strings.Cut(C.GoString(id), ":")
	if !ok || session == "" || requestID == "" {
		return 0
	}
	select {
	case openharmonyRequests <- struct{}{}:
	default:
		return 0
	}
	call := C.GoString(payload)
	go func() {
		defer func() { <-openharmonyRequests }()
		response := application.WailsOhNativeHandleRuntimeCall(call)
		application.WailsOhNativeResolveRuntimeCall(session, requestID, response)
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
func WailsFilePickerDone(id C.uint) {
	application.WailsOhNativeFilePickerDone(uint(id))
}

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

//export WailsStart
func WailsStart() { application.WailsOhNativeStartMain() }
