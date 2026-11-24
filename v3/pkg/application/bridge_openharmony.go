//go:build openharmony && cgo && !server

package application

/*
#cgo LDFLAGS: -lace_napi.z

#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include "node_api.h"

static napi_threadsafe_function g_tsfn = NULL;
static pthread_mutex_t g_bridge_lock = PTHREAD_MUTEX_INITIALIZER;
static _Thread_local int g_on_arkts;

typedef struct {
    uint32_t call_id;
    char *name;
    char *args;
} wails_bridge_call_t;

static void wails_tsfn_js_cb(napi_env env, napi_value js_cb, void *context, void *data) {
    wails_bridge_call_t *call = (wails_bridge_call_t *)data;
    if (call == NULL) return;
    if (env != NULL && js_cb != NULL) {
        napi_value argv[3], receiver;
        if (napi_create_uint32(env, call->call_id, &argv[0]) == napi_ok &&
            napi_create_string_utf8(env, call->name, NAPI_AUTO_LENGTH, &argv[1]) == napi_ok &&
            napi_create_string_utf8(env, call->args, NAPI_AUTO_LENGTH, &argv[2]) == napi_ok &&
            napi_get_undefined(env, &receiver) == napi_ok) {
            napi_call_function(env, receiver, js_cb, 3, argv, NULL);
        }
    }
    free(call->name); free(call->args); free(call);
}

static int wails_oh_init_bridge(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value callback, name;
    napi_valuetype type;
    if (napi_get_cb_info(env, info, &argc, &callback, NULL, NULL) != napi_ok || argc != 1 ||
        napi_typeof(env, callback, &type) != napi_ok || type != napi_function ||
        napi_create_string_utf8(env, "wailsBridgeDispatch", NAPI_AUTO_LENGTH, &name) != napi_ok) return 1;
    pthread_mutex_lock(&g_bridge_lock);
    napi_status rc = g_tsfn != NULL ? napi_invalid_arg : napi_create_threadsafe_function(env,
        callback, NULL, name, 64, 1, NULL, NULL, NULL, wails_tsfn_js_cb, &g_tsfn);
    if (rc == napi_ok) {
        g_on_arkts = 1;
        napi_unref_threadsafe_function(env, g_tsfn);
    }
    pthread_mutex_unlock(&g_bridge_lock);
    return rc != napi_ok;
}

static void wails_oh_close_bridge(void) {
    pthread_mutex_lock(&g_bridge_lock);
    if (g_tsfn != NULL && g_on_arkts) {
        napi_release_threadsafe_function(g_tsfn, napi_tsfn_abort);
        g_tsfn = NULL;
        g_on_arkts = 0;
    }
    pthread_mutex_unlock(&g_bridge_lock);
}

static int wails_oh_is_arkts(void) { return g_on_arkts; }

static int wails_oh_call_bridge(uint32_t call_id, char *name, char *args) {
    wails_bridge_call_t *call = (wails_bridge_call_t *)malloc(sizeof(wails_bridge_call_t));
    if (call == NULL) { free(name); free(args); return 1; }
    call->call_id = call_id; call->name = name; call->args = args;
    pthread_mutex_lock(&g_bridge_lock);
    napi_status rc = g_tsfn ? napi_call_threadsafe_function(g_tsfn, call, napi_tsfn_nonblocking) : napi_closing;
    pthread_mutex_unlock(&g_bridge_lock);
    if (rc != napi_ok) { free(name); free(args); free(call); return 1; }
    return 0;
}

static char *wails_oh_get_string(napi_env env, napi_value value) {
    if (value == NULL) return NULL;
    size_t len = 0;
    if (napi_get_value_string_utf8(env, value, NULL, 0, &len) != napi_ok) return NULL;
    char *buf = (char *)malloc(len + 1);
    if (buf == NULL) return NULL;
    size_t copied = 0;
    if (napi_get_value_string_utf8(env, value, buf, len + 1, &copied) != napi_ok) { free(buf); return NULL; }
    buf[copied] = '\0';
    return buf;
}

static int wails_oh_id_arg(napi_env env, napi_callback_info info, uint32_t *out) {
    size_t argc = 1;
    napi_value arg;
    return napi_get_cb_info(env, info, &argc, &arg, NULL, NULL) == napi_ok && argc == 1 &&
        napi_get_value_uint32(env, arg, out) == napi_ok;
}

static char *wails_oh_string_arg(napi_env env, napi_callback_info info, int index) {
    size_t argc = 2;
    napi_value argv[2];
    if (napi_get_cb_info(env, info, &argc, argv, NULL, NULL) != napi_ok || argc <= (size_t)index) return NULL;
    return wails_oh_get_string(env, argv[index]);
}
*/
import "C"

import (
	"encoding/json"
	"sync"
	"time"
	"unsafe"
)

var (
	ohPendingMu   sync.Mutex
	ohPendingNext uint32
	ohPending     = map[uint32]chan string{}
)

// WailsOhNativeInit registers the ArkTS command callback on its owning thread.
//
//export WailsOhNativeInit
func WailsOhNativeInit(env unsafe.Pointer, info unsafe.Pointer) unsafe.Pointer {
	cEnv := C.napi_env(env)
	var result C.napi_value
	C.napi_get_boolean(cEnv, C.wails_oh_init_bridge(cEnv, C.napi_callback_info(info)) == 0, &result)
	return unsafe.Pointer(result)
}

// WailsOhNativeStartMain starts the main function registered through
// RegisterOpenharmonyMain.
//
// The shell calls this once, after the bridge is open: in -buildmode=c-shared
// the application's main() is never invoked automatically, so the host has to
// start it (mirrors Android's nativeInit starting the registered main).
//
//export WailsOhNativeStartMain
func WailsOhNativeStartMain() {
	openharmonyMainLock.Lock()
	mainFunc := openharmonyMainFunc
	openharmonyMainLock.Unlock()
	if mainFunc == nil {
		openharmonyDebugLogf("[NAPI] no main function registered: call application.RegisterOpenharmonyMain(main)")
		return
	}
	go mainFunc()
}

// WailsOhClose releases the bridge and cancels pending replies.
//
//export WailsOhClose
func WailsOhClose(env, _ unsafe.Pointer) unsafe.Pointer {
	C.wails_oh_close_bridge()
	ohPendingMu.Lock()
	for id, ch := range ohPending {
		delete(ohPending, id)
		close(ch)
	}
	ohPendingMu.Unlock()
	var result C.napi_value
	C.napi_get_undefined(C.napi_env(env), &result)
	return unsafe.Pointer(result)
}

func ohOnArkTSThread() bool { return C.wails_oh_is_arkts() != 0 }

// WailsOhOnBridgeResult is invoked from ArkTS on the JS thread to deliver a
// reply for a synchronous bridge call: WailsOhOnBridgeResult(callID, resultJSON).
//
//export WailsOhOnBridgeResult
func WailsOhOnBridgeResult(env unsafe.Pointer, info unsafe.Pointer) unsafe.Pointer {
	cEnv := C.napi_env(env)
	cInfo := C.napi_callback_info(info)

	var id C.uint32_t
	if C.wails_oh_id_arg(cEnv, cInfo, &id) == 0 {
		return nil
	}
	callID := uint32(id)

	var result string
	cResult := C.wails_oh_string_arg(cEnv, cInfo, 1)
	if cResult != nil {
		result = C.GoString(cResult)
		C.free(unsafe.Pointer(cResult))
	}

	ohPendingMu.Lock()
	ch := ohPending[callID]
	delete(ohPending, callID)
	ohPendingMu.Unlock()
	if ch != nil {
		ch <- result
	}
	return nil
}

// ohBridgeCall dispatches a named method (with a JSON arg) onto the JS thread
// and awaits the bridge reply (callID round-trip).
func ohBridgeCall(name, argsJSON string) (string, bool) {
	if ohOnArkTSThread() {
		return "", false
	}
	ohPendingMu.Lock()
	for {
		ohPendingNext++
		if ohPendingNext != 0 && ohPending[ohPendingNext] == nil {
			break
		}
	}
	id := ohPendingNext
	ch := make(chan string, 1)
	ohPending[id] = ch
	ohPendingMu.Unlock()

	if !ohBridgePost(id, name, argsJSON) {
		ohPendingMu.Lock()
		delete(ohPending, id)
		ohPendingMu.Unlock()
		return "", false
	}

	select {
	case result, ok := <-ch:
		return result, ok
	case <-time.After(5 * time.Second):
		ohPendingMu.Lock()
		delete(ohPending, id)
		ohPendingMu.Unlock()
		return "", false
	}
}

func ohBridgePost(id uint32, name, args string) bool {
	// C owns both buffers on every path.
	return C.wails_oh_call_bridge(C.uint32_t(id), C.CString(name), C.CString(args)) == 0
}

func ohBridgeVoid(name, argsJSON string) { ohBridgePost(0, name, argsJSON) }

func ohBridgeBool(name, argsJSON string) bool {
	res, ok := ohBridgeCall(name, argsJSON)
	return ok && res == "true"
}

func ohBridgeString(name, argsJSON string) (string, bool) { return ohBridgeCall(name, argsJSON) }

func ohBridgeVoidInt(name string, v int) {
	payload, _ := json.Marshal(struct {
		ID int `json:"id"`
	}{ID: v})
	ohBridgeVoid(name, string(payload))
}

func ohBridgeVoidIntString(name string, id int, arg string) {
	payload, _ := json.Marshal(struct {
		ID      int    `json:"id"`
		Options string `json:"options"`
	}{ID: id, Options: arg})
	ohBridgeVoid(name, string(payload))
}

func executeJavaScript(js string) {
	if js == "" {
		return
	}
	payload, _ := json.Marshal(struct {
		Script string `json:"script"`
	}{Script: js})
	ohBridgeVoid("executeJavaScript", string(payload))
}
