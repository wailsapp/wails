//go:build openharmony && cgo && !server

package application

/*
#cgo LDFLAGS: -lace_napi.z

#include <stdlib.h>
#include <string.h>
#include "node_api.h"

static napi_ref g_bridge = NULL;
static napi_threadsafe_function g_tsfn = NULL;

typedef struct {
    uint64_t call_id;
    char *name;
    char *args;
} wails_bridge_call_t;

static void wails_tsfn_js_cb(napi_env env, napi_value js_cb, void *context, void *data) {
    wails_bridge_call_t *call = (wails_bridge_call_t *)data;
    if (call == NULL) return;
    if (g_bridge != NULL && call->name != NULL) {
        napi_value bridge_obj = NULL;
        if (napi_get_reference_value(env, g_bridge, &bridge_obj) == napi_ok && bridge_obj != NULL) {
            napi_value name_val = NULL, method = NULL;
            napi_create_string_utf8(env, call->name, NAPI_AUTO_LENGTH, &name_val);
            napi_get_property(env, bridge_obj, name_val, &method);
            napi_valuetype mt = napi_undefined;
            if (method != NULL) napi_typeof(env, method, &mt);
            if (mt == napi_function) {
                napi_value argv[2];
                napi_create_uint32(env, (uint32_t)call->call_id, &argv[0]);
                napi_create_string_utf8(env, call->args ? call->args : "{}", NAPI_AUTO_LENGTH, &argv[1]);
                napi_value result = NULL;
                napi_call_function(env, bridge_obj, method, 2, argv, &result);
            }
        }
    }
    free(call->name); free(call->args); free(call);
}

static int wails_oh_init_bridge(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value argv[1];
    napi_get_cb_info(env, info, &argc, argv, NULL, NULL);
    if (argc < 1) return 1;
    if (g_bridge != NULL) return 0;
    napi_create_reference(env, argv[0], 1, &g_bridge);

    napi_value name = NULL, marker = NULL;
    napi_create_string_utf8(env, "wailsBridgeDispatch", NAPI_AUTO_LENGTH, &name);
    napi_create_function(env, "marker", 0, NULL, NULL, &marker);
    napi_status rc = napi_create_threadsafe_function(env, marker, NULL, name, 0, 1,
                                                     NULL, NULL, NULL, wails_tsfn_js_cb, &g_tsfn);
    if (rc != napi_ok) return 1;
    napi_unref_threadsafe_function(env, g_tsfn);
    return 0;
}

static int wails_oh_call_bridge(uint64_t call_id, char *name, char *args) {
    if (g_tsfn == NULL || name == NULL) { free(name); free(args); return 1; }
    wails_bridge_call_t *call = (wails_bridge_call_t *)malloc(sizeof(wails_bridge_call_t));
    if (call == NULL) { free(name); free(args); return 1; }
    call->call_id = call_id; call->name = name; call->args = args;
    napi_status rc = napi_call_threadsafe_function(g_tsfn, call, napi_tsfn_blocking);
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

static int wails_oh_int_arg(napi_env env, napi_callback_info info, int index, int *out) {
    size_t argc = 0;
    napi_value argv[8];
    napi_get_cb_info(env, info, &argc, argv, NULL, NULL);
    if ((int)argc <= index) return 0;
    return napi_get_value_int32(env, argv[index], out) == napi_ok;
}

static char *wails_oh_string_arg(napi_env env, napi_callback_info info, int index) {
    size_t argc = 0;
    napi_value argv[8];
    napi_get_cb_info(env, info, &argc, argv, NULL, NULL);
    if ((int)argc <= index) return NULL;
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
	ohPendingNext uint64 = 1
	ohPending     = map[uint64]chan string{}
)

// WailsOhNativeInit is exported to ArkTS (registered in the shell's napi_init).
// ArkTS passes its WailsBridge object; the C helper wraps it for dispatch.
//
//export WailsOhNativeInit
func WailsOhNativeInit(env unsafe.Pointer, info unsafe.Pointer) unsafe.Pointer {
	C.wails_oh_init_bridge(C.napi_env(env), C.napi_callback_info(info))
	return nil
}

// WailsOhOnBridgeResult is invoked from ArkTS on the JS thread to deliver a
// reply for a synchronous bridge call: WailsOhOnBridgeResult(callID, resultJSON).
//
//export WailsOhOnBridgeResult
func WailsOhOnBridgeResult(env unsafe.Pointer, info unsafe.Pointer) unsafe.Pointer {
	cEnv := C.napi_env(env)
	cInfo := C.napi_callback_info(info)

	var idInt C.int
	if C.wails_oh_int_arg(cEnv, cInfo, 0, &idInt) == 0 {
		return nil
	}
	callID := uint64(idInt)

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
	ohPendingMu.Lock()
	id := ohPendingNext
	ohPendingNext++
	ch := make(chan string, 1)
	ohPending[id] = ch
	ohPendingMu.Unlock()

	cname := C.CString(name)
	cargs := C.CString(argsJSON)
	if C.wails_oh_call_bridge(C.uint64_t(id), cname, cargs) != 0 {
		C.free(unsafe.Pointer(cname))
		C.free(unsafe.Pointer(cargs))
		ohPendingMu.Lock()
		delete(ohPending, id)
		ohPendingMu.Unlock()
		return "", false
	}

	select {
	case result := <-ch:
		return result, true
	case <-time.After(5 * time.Second):
		ohPendingMu.Lock()
		delete(ohPending, id)
		ohPendingMu.Unlock()
		return "", false
	}
}

func ohBridgeVoid(name, argsJSON string) { _, _ = ohBridgeCall(name, argsJSON) }

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
