/*
 * napi_init.cpp — OpenHarmony native shell bridge.
 *
 * Exposes the Go backend to ArkTS as a NAPI module. The Go side is a
 * c-shared library (libwails_go.so) whose exports are declared in the
 * generated header libwails_go.h.
 *
 * The function names below form the contract with the ArkTS shell; keep them
 * in sync with pages/Index.ets.
 */

#include <node_api.h>
#include <cstdlib>
#include <cstring>
#include <mutex>
#include <string>
#include <thread>

#include "libwails_go.h"

namespace {

std::once_flag runtimeStart;

napi_value Boolean(napi_env env, bool value) {
    napi_value result = nullptr;
    napi_get_boolean(env, value, &result);
    return result;
}

napi_value Undefined(napi_env env) {
    napi_value result = nullptr;
    napi_get_undefined(env, &result);
    return result;
}


bool String(napi_env env, napi_value value, std::string &result) {
    napi_valuetype type;
    size_t length = 0;
    if (napi_typeof(env, value, &type) != napi_ok || type != napi_string ||
        napi_get_value_string_utf8(env, value, nullptr, 0, &length) != napi_ok) return false;
    result.resize(length + 1);
    if (napi_get_value_string_utf8(env, value, &result[0], result.size(), &length) != napi_ok) return false;
    result.resize(length);
    return result.find('\0') == std::string::npos;
}

// Open registers the ArkTS dispatcher with the Go bridge. The Go side owns
// the TSFN, its queue and the matching close lifecycle.
napi_value Open(napi_env env, napi_callback_info info) {
    return reinterpret_cast<napi_value>(WailsBridgeOpen(env, info));
}

// Result delivers the reply of a synchronous bridge call.
napi_value Result(napi_env env, napi_callback_info info) {
    return reinterpret_cast<napi_value>(WailsBridgeResult(env, info));
}

// Close releases the bridge.
napi_value Close(napi_env env, napi_callback_info info) {
    return reinterpret_cast<napi_value>(WailsBridgeClose(env, info));
}

// Start runs the application's main function on a native thread.
napi_value Start(napi_env env, napi_callback_info) {
    std::call_once(runtimeStart, [] { std::thread([] { WailsStart(); }).detach(); });
    return Boolean(env, true);
}

// Ready reports whether the Go runtime can serve assets yet.
napi_value Ready(napi_env env, napi_callback_info) {
    return Boolean(env, WailsReady() != 0);
}

// Asset serves an application asset: the response body and its MIME type.
napi_value Asset(napi_env env, napi_callback_info info) {
    napi_value arg, result = nullptr;
    size_t argc = 1;
    std::string path;
    napi_get_null(env, &result);
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) != napi_ok || argc != 1 ||
        !String(env, arg, path)) return result;
    size_t length = 0;
    char *data = WailsAsset(&path[0], &length);
    if (!data) return result;
    char *mime = WailsMime(&path[0]);
    napi_value object, body, contentType;
    void *buffer = nullptr;
    bool ok = napi_create_object(env, &object) == napi_ok &&
        napi_create_arraybuffer(env, length, &buffer, &body) == napi_ok;
    if (ok && length) std::memcpy(buffer, data, length);
    if (ok) ok = napi_create_string_utf8(env, mime ? mime : "application/octet-stream",
        NAPI_AUTO_LENGTH, &contentType) == napi_ok &&
        napi_set_named_property(env, object, "body", body) == napi_ok &&
        napi_set_named_property(env, object, "mime", contentType) == napi_ok;
    std::free(data);
    std::free(mime);
    return ok ? object : result;
}

// Request starts an asynchronous runtime call. The id carries the document
// session and the request id so the reply reaches the right document.
napi_value Request(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    std::string id, payload;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) != napi_ok || argc != 2 ||
        !String(env, args[0], id) || !String(env, args[1], payload)) return Boolean(env, false);
    return Boolean(env, WailsRequest(&id[0], &payload[0]) != 0);
}

// Message forwards a window message from the page into the Go application.
napi_value Message(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    std::string message;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) == napi_ok && argc == 1 &&
        String(env, arg, message)) WailsMessage(&message[0]);
    napi_value result = nullptr;
    napi_get_undefined(env, &result);
    return result;
}

// RunOnMainThread runs a function that the Go backend queued for the UI thread.
napi_value RunOnMainThread(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    uint32_t id;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) != napi_ok || argc != 1 ||
        napi_get_value_uint32(env, arg, &id) != napi_ok) return Boolean(env, false);
    WailsRunMain(id);
    return Boolean(env, true);
}

// DialogResult delivers the selected button index for a message dialog.
napi_value DialogResult(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    uint32_t id = 0;
    int32_t index = -1;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) == napi_ok && argc == 2 &&
        napi_get_value_uint32(env, args[0], &id) == napi_ok &&
        napi_get_value_int32(env, args[1], &index) == napi_ok) {
        WailsDialogResult(id, static_cast<int>(index));
    }
    return Undefined(env);
}

// FilePickerResult delivers one selected path for an open-file dialog.
napi_value FilePickerResult(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    uint32_t id = 0;
    std::string path;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) == napi_ok && argc == 2 &&
        napi_get_value_uint32(env, args[0], &id) == napi_ok && String(env, args[1], path)) {
        WailsFilePickerResult(id, &path[0]);
    }
    return Undefined(env);
}

// FilePickerDone closes the pending open-file dialog.
napi_value FilePickerDone(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    uint32_t id = 0;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) == napi_ok && argc == 1 &&
        napi_get_value_uint32(env, arg, &id) == napi_ok) WailsFilePickerDone(id);
    return Undefined(env);
}

napi_value Init(napi_env env, napi_value exports) {
    napi_property_descriptor methods[] = {
        {"open", nullptr, Open, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"result", nullptr, Result, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"close", nullptr, Close, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"start", nullptr, Start, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"ready", nullptr, Ready, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"asset", nullptr, Asset, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"request", nullptr, Request, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"message", nullptr, Message, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"runOnMainThread", nullptr, RunOnMainThread, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"dialogResult", nullptr, DialogResult, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"filePickerResult", nullptr, FilePickerResult, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"filePickerDone", nullptr, FilePickerDone, nullptr, nullptr, nullptr, napi_default, nullptr}
    };
    return napi_define_properties(env, exports, sizeof(methods) / sizeof(methods[0]), methods) == napi_ok
        ? exports : nullptr;
}

napi_module module = {1, 0, nullptr, Init, "entry", nullptr, {0}};

}  // namespace

extern "C" __attribute__((constructor)) void RegisterEntryModule() {
    napi_module_register(&module);
}
