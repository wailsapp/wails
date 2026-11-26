#include <node_api.h>
#include <cstdint>
#include "libhomecast_go.h"

static napi_threadsafe_function session = nullptr;

static void Deliver(napi_env env, napi_value callback, void *, void *payload) {
    if (!env || !callback) return;
    uintptr_t packed = reinterpret_cast<uintptr_t>(payload);
    napi_value args[2], receiver;
    if (napi_create_uint32(env, packed >> 32, &args[0]) != napi_ok ||
        napi_create_uint32(env, packed & UINT32_MAX, &args[1]) != napi_ok ||
        napi_get_undefined(env, &receiver) != napi_ok) return;
    napi_call_function(env, receiver, callback, 2, args, nullptr);
}

static napi_value Boolean(napi_env env, bool result) {
    napi_value value = nullptr;
    napi_get_boolean(env, result, &value);
    return value;
}

static napi_value Open(napi_env env, napi_callback_info info) {
    if (session) return Boolean(env, false);
    size_t argc = 1;
    napi_value callback, name;
    napi_valuetype type;
    if (napi_get_cb_info(env, info, &argc, &callback, nullptr, nullptr) != napi_ok || argc != 1 ||
        napi_typeof(env, callback, &type) != napi_ok || type != napi_function ||
        napi_create_string_utf8(env, "wailsM2", NAPI_AUTO_LENGTH, &name) != napi_ok) return Boolean(env, false);
    return Boolean(env, napi_create_threadsafe_function(env, callback, nullptr, name, 16, 1,
        nullptr, nullptr, nullptr, Deliver, &session) == napi_ok);
}

static napi_value Request(napi_env env, napi_callback_info info) {
    size_t argc = 2;
    napi_value args[2];
    uint32_t id, delay;
    if (!session || napi_get_cb_info(env, info, &argc, args, nullptr, nullptr) != napi_ok || argc != 2 ||
        napi_get_value_uint32(env, args[0], &id) != napi_ok || id == 0 || id > 65535 ||
        napi_get_value_uint32(env, args[1], &delay) != napi_ok || delay > 2000 ||
        napi_acquire_threadsafe_function(session) != napi_ok) return Boolean(env, false);
    StartProbe(session, id, delay);
    return Boolean(env, true);
}

static napi_value Close(napi_env env, napi_callback_info) {
    auto fn = session;
    session = nullptr;
    return Boolean(env, fn && napi_release_threadsafe_function(fn, napi_tsfn_abort) == napi_ok);
}

static napi_value Init(napi_env env, napi_value exports) {
    napi_property_descriptor desc[] = {
        {"open", nullptr, Open, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"request", nullptr, Request, nullptr, nullptr, nullptr, napi_default, nullptr},
        {"close", nullptr, Close, nullptr, nullptr, nullptr, napi_default, nullptr}
    };
    return napi_define_properties(env, exports, 3, desc) == napi_ok ? exports : nullptr;
}

static napi_module module = {1, 0, nullptr, Init, "entry", nullptr, {0}};
extern "C" __attribute__((constructor)) void RegisterEntryModule() {
    napi_module_register(&module);
}
