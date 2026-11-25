#include <node_api.h>
#include <cstdlib>
#include <cstring>
#include <mutex>
#include <string>
#include <thread>
#include "libhomecast_go.h"

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

napi_value Open(napi_env env, napi_callback_info info) {
    // Go owns the TSFN, its queue and the matching close/finalizer lifecycle.
    return reinterpret_cast<napi_value>(WailsBridgeOpen(env, info));
}

napi_value Result(napi_env env, napi_callback_info info) {
    return reinterpret_cast<napi_value>(WailsBridgeResult(env, info));
}

napi_value Close(napi_env env, napi_callback_info info) {
    return reinterpret_cast<napi_value>(WailsBridgeClose(env, info));
}

napi_value Start(napi_env env, napi_callback_info) {
    std::call_once(runtimeStart, [] { std::thread([] { WailsStart(); }).detach(); });
    return Boolean(env, true);
}

napi_value Ready(napi_env env, napi_callback_info) {
    return Boolean(env, WailsReady() != 0);
}

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

napi_value Request(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    std::string id, payload;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) != napi_ok || argc != 2 ||
        !String(env, args[0], id) || !String(env, args[1], payload)) return Boolean(env, false);
    return Boolean(env, WailsRequest(&id[0], &payload[0]) != 0);
}

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

napi_value RunOnMainThread(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    uint32_t id;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) != napi_ok || argc != 1 ||
        napi_get_value_uint32(env, arg, &id) != napi_ok) return Boolean(env, false);
    WailsRunMain(id);
    return Boolean(env, true);
}

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

// FilePickerResult delivers one selected path for a file dialog.
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

// FilePickerDone closes the pending file dialog.
napi_value FilePickerDone(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    uint32_t id = 0;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) == napi_ok && argc == 1 &&
        napi_get_value_uint32(env, arg, &id) == napi_ok) WailsFilePickerDone(id);
    return Undefined(env);
}

napi_value AbilityStarted(napi_env env, napi_callback_info) { WailsAbilityStarted(); return Undefined(env); }
napi_value AbilityResumed(napi_env env, napi_callback_info) { WailsAbilityResumed(); return Undefined(env); }
napi_value AbilityPaused(napi_env env, napi_callback_info) { WailsAbilityPaused(); return Undefined(env); }
napi_value AbilityStopped(napi_env env, napi_callback_info) { WailsAbilityStopped(); return Undefined(env); }
napi_value LowMemory(napi_env env, napi_callback_info) { WailsLowMemory(); return Undefined(env); }

napi_value EmitSystemEvent(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    std::string name, payload;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) == napi_ok && argc == 2 &&
        String(env, args[0], name) && String(env, args[1], payload)) {
        WailsEmitSystemEvent(&name[0], &payload[0]);
    }
    return Undefined(env);
}

napi_value EmitEvent(napi_env env, napi_callback_info info) {
    napi_value args[2];
    size_t argc = 2;
    std::string name, payload;
    if (napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) == napi_ok && argc == 2 &&
        String(env, args[0], name) && String(env, args[1], payload)) {
        WailsEmitEvent(&name[0], &payload[0]);
    }
    return Undefined(env);
}

napi_value PageStarted(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    std::string url;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) == napi_ok && argc == 1 &&
        String(env, arg, url)) WailsPageStarted(&url[0]);
    return Undefined(env);
}

napi_value PageFinished(napi_env env, napi_callback_info info) {
    napi_value arg;
    size_t argc = 1;
    std::string url;
    if (napi_get_cb_info(env, info, &argc, &arg, nullptr, nullptr) == napi_ok && argc == 1 &&
        String(env, arg, url)) WailsPageFinished(&url[0]);
    return Undefined(env);
}

napi_value Shutdown(napi_env env, napi_callback_info) { WailsShutdown(); return Undefined(env); }

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
        {"filePickerDone", nullptr, FilePickerDone, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"abilityStarted", nullptr, AbilityStarted, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"abilityResumed", nullptr, AbilityResumed, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"abilityPaused", nullptr, AbilityPaused, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"abilityStopped", nullptr, AbilityStopped, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"lowMemory", nullptr, LowMemory, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"emitSystemEvent", nullptr, EmitSystemEvent, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"emitEvent", nullptr, EmitEvent, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"pageStarted", nullptr, PageStarted, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"pageFinished", nullptr, PageFinished, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"shutdown", nullptr, Shutdown, nullptr, nullptr, nullptr, napi_default, nullptr}
    };
    return napi_define_properties(env, exports, sizeof(methods) / sizeof(methods[0]), methods) == napi_ok
        ? exports : nullptr;
}

napi_module module = {1, 0, nullptr, Init, "entry", nullptr, {0}};
}

extern "C" __attribute__((constructor)) void RegisterEntryModule() {
    napi_module_register(&module);
}
