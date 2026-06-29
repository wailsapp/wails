//go:build android && cgo && !server

package application

/*
#include <jni.h>
#include <stdlib.h>
#include <string.h>
#include <android/log.h>

// Global JavaVM reference for thread attachment
static JavaVM* g_jvm = NULL;

// Global reference to the WailsBridge object (must be a global ref, not local)
static jobject g_bridge = NULL;

// Cached method ID for the hot executeJavaScript path
static jmethodID g_executeJsMethod = NULL;

// Verbose (debug) logging toggle, set from Go at startup
static int g_verbose = 0;
static void wails_set_verbose(int v) { g_verbose = v; }

#define WLOGD(...) do { if (g_verbose) __android_log_print(ANDROID_LOG_DEBUG, "Wails", __VA_ARGS__); } while (0)
#define WLOGE(...) __android_log_print(ANDROID_LOG_ERROR, "Wails", __VA_ARGS__)

static void wails_log(int prio, const char* msg) {
    __android_log_write(prio, "Wails", msg);
}

// Helper to convert Java String to C string
static const char* jstringToC(JNIEnv *env, jstring jstr) {
    if (jstr == NULL) return NULL;
    return (*env)->GetStringUTFChars(env, jstr, NULL);
}

// Helper to release Java String
static void releaseJString(JNIEnv *env, jstring jstr, const char* cstr) {
    if (jstr != NULL && cstr != NULL) {
        (*env)->ReleaseStringUTFChars(env, jstr, cstr);
    }
}

// Helper to create Java byte array from C data
static jbyteArray createByteArray(JNIEnv *env, const void* data, int len) {
    jbyteArray arr = (*env)->NewByteArray(env, len < 0 ? 0 : len);
    if (arr != NULL && data != NULL && len > 0) {
        (*env)->SetByteArrayRegion(env, arr, 0, len, (const jbyte*)data);
    }
    return arr;
}

// Helper to create Java String from C string
static jstring createJString(JNIEnv *env, const char* str) {
    if (str == NULL) return NULL;
    return (*env)->NewStringUTF(env, str);
}

// Store JavaVM and create a global reference to the bridge
static void storeBridgeRef(JNIEnv *env, jobject bridge) {
    if ((*env)->GetJavaVM(env, &g_jvm) != 0) {
        WLOGE("storeBridgeRef: GetJavaVM failed");
        return;
    }
    g_bridge = (*env)->NewGlobalRef(env, bridge);
    if (g_bridge == NULL) {
        WLOGE("storeBridgeRef: NewGlobalRef failed");
        return;
    }
    jclass bridgeClass = (*env)->GetObjectClass(env, g_bridge);
    if (bridgeClass != NULL) {
        g_executeJsMethod = (*env)->GetMethodID(env, bridgeClass, "executeJavaScript", "(Ljava/lang/String;)V");
        (*env)->DeleteLocalRef(env, bridgeClass);
    }
}

// Get a JNIEnv for the current thread, attaching it if necessary.
// Sets *needsDetach when the caller must detach afterwards.
static JNIEnv* wailsGetEnv(int* needsDetach) {
    *needsDetach = 0;
    if (g_jvm == NULL) return NULL;
    JNIEnv* env = NULL;
    jint r = (*g_jvm)->GetEnv(g_jvm, (void**)&env, JNI_VERSION_1_6);
    if (r == JNI_EDETACHED) {
        if ((*g_jvm)->AttachCurrentThread(g_jvm, &env, NULL) != 0) {
            WLOGE("wailsGetEnv: AttachCurrentThread failed");
            return NULL;
        }
        *needsDetach = 1;
    } else if (r != JNI_OK) {
        WLOGE("wailsGetEnv: GetEnv failed: %d", (int)r);
        return NULL;
    }
    return env;
}

static void wailsReleaseEnv(int needsDetach) {
    if (needsDetach && g_jvm != NULL) {
        (*g_jvm)->DetachCurrentThread(g_jvm);
    }
}

static void clearException(JNIEnv* env, const char* where) {
    if ((*env)->ExceptionCheck(env)) {
        WLOGE("Java exception in %s", where);
        (*env)->ExceptionDescribe(env);
        (*env)->ExceptionClear(env);
    }
}

// Call `String name()` on the bridge. Returns a malloc'd C string (caller
// frees) or NULL.
static char* callBridgeStringMethod(const char* name) {
    if (g_bridge == NULL) return NULL;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return NULL;
    char* result = NULL;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "()Ljava/lang/String;");
        if (mid != NULL) {
            jstring jresult = (jstring)(*env)->CallObjectMethod(env, g_bridge, mid);
            clearException(env, name);
            if (jresult != NULL) {
                const char* chars = (*env)->GetStringUTFChars(env, jresult, NULL);
                if (chars != NULL) {
                    result = strdup(chars);
                    (*env)->ReleaseStringUTFChars(env, jresult, chars);
                }
                (*env)->DeleteLocalRef(env, jresult);
            }
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
    return result;
}

// Call `String name(String)` on the bridge. Returns a malloc'd C string
// (caller frees) or NULL.
static char* callBridgeStringStringMethod(const char* name, const char* arg) {
    if (g_bridge == NULL) return NULL;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return NULL;
    char* result = NULL;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "(Ljava/lang/String;)Ljava/lang/String;");
        if (mid != NULL) {
            jstring jarg = (*env)->NewStringUTF(env, arg);
            jstring jresult = (jstring)(*env)->CallObjectMethod(env, g_bridge, mid, jarg);
            clearException(env, name);
            if (jresult != NULL) {
                const char* chars = (*env)->GetStringUTFChars(env, jresult, NULL);
                if (chars != NULL) {
                    result = strdup(chars);
                    (*env)->ReleaseStringUTFChars(env, jresult, chars);
                }
                (*env)->DeleteLocalRef(env, jresult);
            }
            if (jarg != NULL) (*env)->DeleteLocalRef(env, jarg);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
    return result;
}

// Call `void name()` on the bridge.
static void callBridgeVoidMethod(const char* name) {
    if (g_bridge == NULL) return;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "()V");
        if (mid != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, mid);
            clearException(env, name);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
}

// Call `void name(String)` on the bridge.
static void callBridgeVoidString(const char* name, const char* arg) {
    if (g_bridge == NULL) return;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "(Ljava/lang/String;)V");
        if (mid != NULL) {
            jstring jarg = (*env)->NewStringUTF(env, arg);
            (*env)->CallVoidMethod(env, g_bridge, mid, jarg);
            clearException(env, name);
            if (jarg != NULL) (*env)->DeleteLocalRef(env, jarg);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
}

// Call `void name(int)` on the bridge.
static void callBridgeVoidInt(const char* name, int v) {
    if (g_bridge == NULL) return;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "(I)V");
        if (mid != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, mid, (jint)v);
            clearException(env, name);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
}

// Call `void name(int, String)` on the bridge.
static void callBridgeVoidIntString(const char* name, int id, const char* arg) {
    if (g_bridge == NULL) return;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "(ILjava/lang/String;)V");
        if (mid != NULL) {
            jstring jarg = (*env)->NewStringUTF(env, arg);
            (*env)->CallVoidMethod(env, g_bridge, mid, (jint)id, jarg);
            clearException(env, name);
            if (jarg != NULL) (*env)->DeleteLocalRef(env, jarg);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
}

// Call `boolean name()` on the bridge.
static int callBridgeBoolMethod(const char* name) {
    if (g_bridge == NULL) return 0;
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return 0;
    int result = 0;
    jclass cls = (*env)->GetObjectClass(env, g_bridge);
    if (cls != NULL) {
        jmethodID mid = (*env)->GetMethodID(env, cls, name, "()Z");
        if (mid != NULL) {
            result = (*env)->CallBooleanMethod(env, g_bridge, mid) == JNI_TRUE;
            clearException(env, name);
        } else {
            clearException(env, name);
        }
        (*env)->DeleteLocalRef(env, cls);
    }
    wailsReleaseEnv(detach);
    return result;
}

// Execute JavaScript via the bridge - can be called from any thread.
static void executeJavaScriptOnBridge(const char* js) {
    if (g_bridge == NULL || g_executeJsMethod == NULL || js == NULL) {
        WLOGE("executeJavaScriptOnBridge: bridge not ready");
        return;
    }
    int detach = 0;
    JNIEnv* env = wailsGetEnv(&detach);
    if (env == NULL) return;
    jstring jJs = (*env)->NewStringUTF(env, js);
    if (jJs != NULL) {
        (*env)->CallVoidMethod(env, g_bridge, g_executeJsMethod, jJs);
        clearException(env, "executeJavaScript");
        (*env)->DeleteLocalRef(env, jJs);
    }
    wailsReleaseEnv(detach);
}
*/
import "C"

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"
	"unsafe"

	"encoding/json"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	// Global reference to the app for JNI callbacks
	globalApp     *App
	globalAppLock sync.RWMutex

	// Android main function registration
	androidMainFunc func()
	androidMainLock sync.Mutex

	// App ready signal
	appReady     = make(chan struct{})
	appReadyOnce sync.Once
)

// androidLogf logs through logcat (tag "Wails") so messages are visible in
// `adb logcat`. Go's stdout/stderr are not routed anywhere on Android.
func androidLogf(level string, format string, a ...interface{}) {
	var prio C.int
	switch level {
	case "debug":
		prio = 3 // ANDROID_LOG_DEBUG
	case "warn":
		prio = 5 // ANDROID_LOG_WARN
	case "error":
		prio = 6 // ANDROID_LOG_ERROR
	default:
		prio = 4 // ANDROID_LOG_INFO
	}
	msg := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(msg))
	C.wails_log(prio, msg)
}

// androidDebugLogf is for the framework's internal diagnostics. It compiles to
// a no-op in production builds (see android_logging_production.go).
func androidDebugLogf(format string, a ...interface{}) {
	if androidVerboseLogging {
		androidLogf("debug", format, a...)
	}
}

// RegisterAndroidMain registers the main function to be called when the
// Android app starts. Call it from init() in your main package:
//
//	func init() {
//		application.RegisterAndroidMain(main)
//	}
func RegisterAndroidMain(mainFunc func()) {
	androidMainLock.Lock()
	defer androidMainLock.Unlock()
	androidMainFunc = mainFunc
}

// signalAppReady signals that the app is ready to serve requests
func signalAppReady() {
	appReadyOnce.Do(func() {
		close(appReady)
	})
}

// waitForAppReady waits for the app to be ready, with a timeout
func waitForAppReady(timeout time.Duration) bool {
	select {
	case <-appReady:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Go-level bridge call API. These are stubbed in
// application_android_nocgo.go so files shared between the cgo and
// non-cgo builds (dialogs, clipboard, screens, ...) can call them
// unconditionally.

func androidBridgeString(method string) (string, bool) {
	cname := C.CString(method)
	defer C.free(unsafe.Pointer(cname))
	cresult := C.callBridgeStringMethod(cname)
	if cresult == nil {
		return "", false
	}
	defer C.free(unsafe.Pointer(cresult))
	return C.GoString(cresult), true
}

func androidBridgeStringString(method string, arg string) (string, bool) {
	cname := C.CString(method)
	carg := C.CString(arg)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(carg))
	cresult := C.callBridgeStringStringMethod(cname, carg)
	if cresult == nil {
		return "", false
	}
	defer C.free(unsafe.Pointer(cresult))
	return C.GoString(cresult), true
}

func androidBridgeVoid(method string) {
	cname := C.CString(method)
	defer C.free(unsafe.Pointer(cname))
	C.callBridgeVoidMethod(cname)
}

func androidBridgeVoidString(method string, arg string) {
	cname := C.CString(method)
	carg := C.CString(arg)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(carg))
	C.callBridgeVoidString(cname, carg)
}

func androidBridgeVoidInt(method string, v int) {
	cname := C.CString(method)
	defer C.free(unsafe.Pointer(cname))
	C.callBridgeVoidInt(cname, C.int(v))
}

func androidBridgeVoidIntString(method string, id int, arg string) {
	cname := C.CString(method)
	carg := C.CString(arg)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(carg))
	C.callBridgeVoidIntString(cname, C.int(id), carg)
}

func androidBridgeBool(method string) bool {
	cname := C.CString(method)
	defer C.free(unsafe.Pointer(cname))
	return C.callBridgeBoolMethod(cname) != 0
}

// executeJavaScript executes JavaScript in the Android WebView via JNI
func executeJavaScript(js string) {
	if js == "" {
		return
	}
	cJs := C.CString(js)
	defer C.free(unsafe.Pointer(cJs))
	C.executeJavaScriptOnBridge(cJs)
}

func (a *App) platformRun() {
	androidDebugLogf("[application_android.go] platformRun: initialising")

	// Propagate logging verbosity to the native layer
	if androidVerboseLogging {
		C.wails_set_verbose(1)
	}

	// Store global reference for JNI callbacks
	globalAppLock.Lock()
	globalApp = a
	globalAppLock.Unlock()

	// Unblock asset serving
	signalAppReady()

	// Populate the ScreenManager so Screens.* runtime calls return data
	// (desktop platforms do this from their event loop; Android has none).
	if screens, err := getScreens(); err == nil && len(screens) > 0 {
		if err := a.Screen.LayoutScreens(screens); err != nil {
			androidDebugLogf("[application_android.go] LayoutScreens failed: %v", err)
		}
	}

	// Emit the typed launch event from here rather than from nativeInit: by
	// this point setupCommonEvents has registered its listeners, so the event
	// (and its Common.ApplicationStarted mapping) cannot be dropped by
	// startup races.
	applicationEvents <- newApplicationEvent(events.Android.ActivityCreated)

	// Block forever - Android manages the app lifecycle via JNI callbacks
	select {}
}

func (a *App) platformQuit() {
	// Android handles app termination through the Activity lifecycle
}

func (a *App) isDarkMode() bool {
	return androidBridgeBool("isDarkMode")
}

func (a *App) isWindows() bool {
	return false
}

// Platform-specific app implementation for Android
type androidApp struct {
	parent *App
}

func newPlatformApp(app *App) *androidApp {
	return &androidApp{
		parent: app,
	}
}

func (a *androidApp) run() error {
	// Wire common events (e.g. map ActivityCreated → Common.ApplicationStarted)
	a.setupCommonEvents()
	a.parent.platformRun()
	return nil
}

func (a *androidApp) destroy() {
}

func (a *androidApp) setIcon(_ []byte) {
	// Android app icon is set through AndroidManifest.xml
}

func (a *androidApp) name() string {
	return a.parent.options.Name
}

func (a *androidApp) GetFlags(options Options) map[string]any {
	return nil
}

func (a *androidApp) getAccentColor() string {
	return ""
}

func (a *androidApp) getCurrentWindowID() uint {
	return 0
}

func (a *androidApp) hide() {
	// Android manages app visibility
}

func (a *androidApp) isDarkMode() bool {
	return a.parent.isDarkMode()
}

func (a *androidApp) on(eventID uint) {
	registerAndroidListener(eventID)
}

func (a *androidApp) setApplicationMenu(_ *Menu) {
	// Android doesn't have application menus
}

func (a *androidApp) show() {
	// Android manages app visibility
}

func (a *androidApp) showAboutDialog(title string, message string, _ []byte) {
	a.parent.Dialog.Info().SetTitle(title).SetMessage(message).Show()
}

func (a *androidApp) getPrimaryScreen() (*Screen, error) {
	if a.parent.Screen.GetPrimary() == nil {
		if err := a.cacheScreens(); err != nil {
			return nil, err
		}
	}
	return a.parent.Screen.GetPrimary(), nil
}

func (a *androidApp) getScreens() ([]*Screen, error) {
	if len(a.parent.Screen.GetAll()) == 0 {
		if err := a.cacheScreens(); err != nil {
			return nil, err
		}
	}
	return a.parent.Screen.GetAll(), nil
}

func (a *androidApp) cacheScreens() error {
	screens, err := getScreens()
	if err != nil {
		return err
	}
	return a.parent.Screen.LayoutScreens(screens)
}

func (a *App) logPlatformInfo() {
	androidDebugLogf("Platform: Android")
}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"platform": "android",
	}
}

func fatalHandler(errFunc func(error)) {
	// Android fatal handler
}

// androidEventListeners records which native event IDs have at least one
// Go-side listener (mirrors the iOS implementation).
var (
	androidEventListeners     = make(map[uint]bool)
	androidEventListenersLock sync.RWMutex
)

func registerAndroidListener(eventID uint) {
	androidEventListenersLock.Lock()
	defer androidEventListenersLock.Unlock()
	androidEventListeners[eventID] = true
}

// androidFirstWindowID returns the ID of the first (usually only) window.
func androidFirstWindowID(app *App) uint {
	if app == nil {
		return 0
	}
	windows := app.Window.GetAll()
	if len(windows) == 0 {
		return 0
	}
	return windows[0].ID()
}

// emitAndroidApplicationEvent forwards a typed application event to the
// event processing loop.
func emitAndroidApplicationEvent(event events.ApplicationEventType) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	applicationEvents <- newApplicationEvent(event)
}

// androidSystemEventTypes maps the canonical android: event names the Java
// system-event receivers send to their typed application event.
var androidSystemEventTypes = map[string]events.ApplicationEventType{
	"android:BatteryChanged": events.Android.BatteryChanged,
	"android:NetworkChanged": events.Android.NetworkChanged,
	"android:ThemeChanged":   events.Android.ThemeChanged,
	"android:ScreenLocked":   events.Android.ScreenLocked,
	"android:ScreenUnlocked": events.Android.ScreenUnlocked,
}

// emitAndroidApplicationEventWithData forwards a typed application event with an
// optional JSON payload (battery level, theme, …) attached to its context, so
// Go listeners can read it via event.Context().Data() / IsDarkMode().
func emitAndroidApplicationEventWithData(event events.ApplicationEventType, jsonStr string) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	evt := newApplicationEvent(event)
	if jsonStr != "" {
		var m map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &m); err == nil && m != nil {
			evt.Context().setData(m)
		}
	}
	applicationEvents <- evt
}

// JNI Export Functions - Called from Java

//export Java_com_wails_app_WailsBridge_nativeInit
func Java_com_wails_app_WailsBridge_nativeInit(env *C.JNIEnv, obj C.jobject, bridge C.jobject) {
	// Store JavaVM and bridge global reference for JNI callbacks
	C.storeBridgeRef(env, bridge)

	// Start the registered main function in a goroutine
	androidMainLock.Lock()
	mainFunc := androidMainFunc
	androidMainLock.Unlock()

	if mainFunc != nil {
		go mainFunc()
	} else {
		androidLogf("error", "No main function registered! Call application.RegisterAndroidMain(main) in init()")
	}
}

//export Java_com_wails_app_WailsBridge_nativeShutdown
func Java_com_wails_app_WailsBridge_nativeShutdown(env *C.JNIEnv, obj C.jobject) {
	globalAppLock.Lock()
	app := globalApp
	globalAppLock.Unlock()
	if app != nil {
		app.Quit()
	}
}

//export Java_com_wails_app_WailsBridge_nativeOnStart
func Java_com_wails_app_WailsBridge_nativeOnStart(env *C.JNIEnv, obj C.jobject) {
	emitAndroidApplicationEvent(events.Android.ActivityStarted)
}

//export Java_com_wails_app_WailsBridge_nativeOnResume
func Java_com_wails_app_WailsBridge_nativeOnResume(env *C.JNIEnv, obj C.jobject) {
	emitAndroidApplicationEvent(events.Android.ActivityResumed)
}

//export Java_com_wails_app_WailsBridge_nativeOnPause
func Java_com_wails_app_WailsBridge_nativeOnPause(env *C.JNIEnv, obj C.jobject) {
	emitAndroidApplicationEvent(events.Android.ActivityPaused)
}

//export Java_com_wails_app_WailsBridge_nativeOnStop
func Java_com_wails_app_WailsBridge_nativeOnStop(env *C.JNIEnv, obj C.jobject) {
	emitAndroidApplicationEvent(events.Android.ActivityStopped)
}

//export Java_com_wails_app_WailsBridge_nativeOnLowMemory
func Java_com_wails_app_WailsBridge_nativeOnLowMemory(env *C.JNIEnv, obj C.jobject) {
	emitAndroidApplicationEvent(events.Android.ApplicationLowMemory)
}

// Java_com_wails_app_WailsBridge_nativeEmitSystemEvent is the funnel the
// Android system-event receivers (battery, network, lock, theme) call to
// deliver a typed android: application event with its JSON payload.
//
//export Java_com_wails_app_WailsBridge_nativeEmitSystemEvent
func Java_com_wails_app_WailsBridge_nativeEmitSystemEvent(env *C.JNIEnv, obj C.jobject, jname C.jstring, jjson C.jstring) {
	cName := C.jstringToC(env, jname)
	name := C.GoString(cName)
	C.releaseJString(env, jname, cName)

	cJSON := C.jstringToC(env, jjson)
	jsonStr := C.GoString(cJSON)
	C.releaseJString(env, jjson, cJSON)

	if eventType, ok := androidSystemEventTypes[name]; ok {
		emitAndroidApplicationEventWithData(eventType, jsonStr)
	} else {
		androidLogf("warn", "[JNI] nativeEmitSystemEvent: unknown event %q", name)
	}
}

// emitNativeEventToJS forwards a named custom event with an optional JSON
// payload to the frontend. Used by the mobile-feature bridges (torch, biometric,
// secure storage, …) to deliver asynchronous results.
func emitNativeEventToJS(name string, jsonStr string) {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}
	var data map[string]any
	if jsonStr != "" {
		_ = json.Unmarshal([]byte(jsonStr), &data)
	}
	app.Event.Emit(name, data)
}

// Java_com_wails_app_WailsBridge_nativeEmitEvent is the funnel the mobile-feature
// bridges call to deliver an arbitrary custom event (with JSON payload) to JS.
//
//export Java_com_wails_app_WailsBridge_nativeEmitEvent
func Java_com_wails_app_WailsBridge_nativeEmitEvent(env *C.JNIEnv, obj C.jobject, jname C.jstring, jjson C.jstring) {
	cName := C.jstringToC(env, jname)
	name := C.GoString(cName)
	C.releaseJString(env, jname, cName)

	cJSON := C.jstringToC(env, jjson)
	jsonStr := C.GoString(cJSON)
	C.releaseJString(env, jjson, cJSON)

	emitNativeEventToJS(name, jsonStr)
}

//export Java_com_wails_app_WailsBridge_nativeOnPageFinished
func Java_com_wails_app_WailsBridge_nativeOnPageFinished(env *C.JNIEnv, obj C.jobject, jurl C.jstring) {
	cUrl := C.jstringToC(env, jurl)
	url := C.GoString(cUrl)
	C.releaseJString(env, jurl, cUrl)

	androidDebugLogf("[JNI] page finished: %s", url)

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app == nil {
		return
	}

	// The bundled @wailsio/runtime detects the Android JavascriptInterface
	// (window.wails.invoke) at module load and announces itself with
	// "wails:runtime:ready", so no runtime injection is needed here. Just
	// forward the typed window event.
	if windowID := androidFirstWindowID(app); windowID != 0 {
		windowEvents <- &windowEvent{
			WindowID: windowID,
			EventID:  uint(events.Android.WebViewPageFinished),
		}
	}
}

//export Java_com_wails_app_WailsBridge_nativeServeAsset
func Java_com_wails_app_WailsBridge_nativeServeAsset(env *C.JNIEnv, obj C.jobject, jpath C.jstring, jmethod C.jstring, jheaders C.jstring) C.jbyteArray {
	cPath := C.jstringToC(env, jpath)
	cMethod := C.jstringToC(env, jmethod)

	goPath := C.GoString(cPath)
	goMethod := C.GoString(cMethod)

	C.releaseJString(env, jpath, cPath)
	C.releaseJString(env, jmethod, cMethod)

	// Wait for the app to be ready (timeout after 10 seconds)
	if !waitForAppReady(10 * time.Second) {
		androidLogf("error", "[JNI] timeout waiting for app to be ready, dropping request for %s", goPath)
		return C.createByteArray(env, nil, 0)
	}

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app == nil || app.assets == nil {
		androidLogf("error", "[JNI] app or assets not initialized after ready signal")
		return C.createByteArray(env, nil, 0)
	}

	data, err := serveAssetForAndroid(app, goMethod, goPath)
	if err != nil {
		androidLogf("error", "[JNI] error serving asset %s: %v", goPath, err)
		return C.createByteArray(env, nil, 0)
	}

	if len(data) == 0 {
		return C.createByteArray(env, nil, 0)
	}
	return C.createByteArray(env, unsafe.Pointer(&data[0]), C.int(len(data)))
}

//export Java_com_wails_app_WailsBridge_nativeHandleMessage
func Java_com_wails_app_WailsBridge_nativeHandleMessage(env *C.JNIEnv, obj C.jobject, jmessage C.jstring) C.jstring {
	cMessage := C.jstringToC(env, jmessage)
	goMessage := C.GoString(cMessage)
	C.releaseJString(env, jmessage, cMessage)

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app == nil {
		cresp := C.CString(`{"error":"App not initialized"}`)
		defer C.free(unsafe.Pointer(cresp))
		return C.createJString(env, cresp)
	}

	handleMessageForAndroid(app, goMessage)

	cresp := C.CString(`{}`)
	defer C.free(unsafe.Pointer(cresp))
	return C.createJString(env, cresp)
}

//export Java_com_wails_app_WailsBridge_nativeHandleRuntimeCall
func Java_com_wails_app_WailsBridge_nativeHandleRuntimeCall(env *C.JNIEnv, obj C.jobject, jpayload C.jstring) C.jstring {
	cPayload := C.jstringToC(env, jpayload)
	goPayload := C.GoString(cPayload)
	C.releaseJString(env, jpayload, cPayload)

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	var response string
	if app == nil {
		response = `{"ok":false,"error":"App not initialized"}`
	} else {
		response = handleRuntimeCallForAndroid(app, goPayload)
	}

	cresp := C.CString(response)
	defer C.free(unsafe.Pointer(cresp))
	return C.createJString(env, cresp)
}

//export Java_com_wails_app_WailsBridge_nativeGetAssetMimeType
func Java_com_wails_app_WailsBridge_nativeGetAssetMimeType(env *C.JNIEnv, obj C.jobject, jpath C.jstring) C.jstring {
	cPath := C.jstringToC(env, jpath)
	goPath := C.GoString(cPath)
	C.releaseJString(env, jpath, cPath)

	cmime := C.CString(getMimeTypeForPath(goPath))
	defer C.free(unsafe.Pointer(cmime))
	return C.createJString(env, cmime)
}

//export Java_com_wails_app_WailsBridge_nativeDialogCallback
func Java_com_wails_app_WailsBridge_nativeDialogCallback(env *C.JNIEnv, obj C.jobject, callbackID C.jint, buttonIndex C.jint) {
	androidDialogCallback(uint(callbackID), int(buttonIndex))
}

//export Java_com_wails_app_WailsBridge_nativeFilePickerResult
func Java_com_wails_app_WailsBridge_nativeFilePickerResult(env *C.JNIEnv, obj C.jobject, callbackID C.jint, jpath C.jstring) {
	cPath := C.jstringToC(env, jpath)
	goPath := C.GoString(cPath)
	C.releaseJString(env, jpath, cPath)
	androidFilePickerResult(uint(callbackID), goPath)
}

//export Java_com_wails_app_WailsBridge_nativeFilePickerDone
func Java_com_wails_app_WailsBridge_nativeFilePickerDone(env *C.JNIEnv, obj C.jobject, callbackID C.jint) {
	androidFilePickerDone(uint(callbackID))
}

//export Java_com_wails_app_WailsBridge_nativeMainThreadCallback
func Java_com_wails_app_WailsBridge_nativeMainThreadCallback(env *C.JNIEnv, obj C.jobject, callbackID C.jint) {
	androidMainThreadCallback(uint(callbackID))
}

// Helper functions

// androidRuntimeReadyWindows tracks windows for which a synthetic
// "wails:runtime:ready" has been injected (see serveAssetForAndroid).
var androidRuntimeReadyWindows sync.Map

func serveAssetForAndroid(app *App, method string, path string) ([]byte, error) {
	// Runtime calls include a query string that must be preserved
	isRuntimeCall := strings.HasPrefix(path, "/wails/runtime")

	if !isRuntimeCall {
		if path == "" || path == "/" {
			path = "/index.html"
		}
	}

	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	if app.assets == nil {
		return nil, fmt.Errorf("asset server not initialized")
	}

	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequest(method, "https://wails.localhost"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// http.NewRequest leaves Body nil for client requests, but handlers
	// reached via ServeHTTP expect the server guarantee of a non-nil Body
	req.Body = http.NoBody

	// Runtime calls need the window ID and name headers so the
	// MessageProcessor can route the call to the right window.
	if isRuntimeCall {
		windows := app.Window.GetAll()
		if len(windows) > 0 {
			windowID := windows[0].ID()
			req.Header.Set("x-wails-window-id", fmt.Sprintf("%d", windowID))
			req.Header.Set("x-wails-window-name", windows[0].Name())

			// The JavaScript runtime announces itself with a
			// "wails:runtime:ready" message; a call to /wails/runtime proves
			// the runtime is mounted, so treat the first one as an implicit
			// ready signal as a fallback. processMessage handles duplicate
			// ready messages gracefully.
			if _, alreadyReady := androidRuntimeReadyWindows.LoadOrStore(windowID, true); !alreadyReady {
				windowMessageBuffer <- &windowMessage{
					windowId: windowID,
					message:  "wails:runtime:ready",
				}
			}
		}
	}

	recorder := httptest.NewRecorder()
	app.assets.ServeHTTP(recorder, req)

	result := recorder.Result()
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// For runtime calls, return the body even for error responses so the
	// JavaScript side can see the error message.
	if isRuntimeCall {
		if result.StatusCode != http.StatusOK {
			androidDebugLogf("[serveAsset] runtime call returned status %d: %s", result.StatusCode, string(body))
		}
		return body, nil
	}

	if result.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("asset not found: %s", path)
	}

	if result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("asset server error: status %d for %s", result.StatusCode, path)
	}

	return body, nil
}

// The Android transport: the WebView cannot deliver fetch() POST bodies to
// shouldInterceptRequest, so the JavaScript runtime routes runtime calls
// through the JavascriptInterface bridge to nativeHandleRuntimeCall.

var (
	androidMessageProcessor     *MessageProcessor
	androidMessageProcessorOnce sync.Once
)

func androidProcessor(app *App) *MessageProcessor {
	androidMessageProcessorOnce.Do(func() {
		androidMessageProcessor = NewMessageProcessor(app.Logger)
	})
	return androidMessageProcessor
}

type androidRuntimeCallPayload struct {
	Object     *int            `json:"object"`
	Method     *int            `json:"method"`
	WindowName string          `json:"windowName"`
	Args       json.RawMessage `json:"args"`
	ClientID   string          `json:"clientId"`
}

// handleRuntimeCallForAndroid processes a runtime call payload from the
// JavaScript Android transport and returns a response envelope:
// {"ok":true,"data":...} / {"ok":true,"text":"..."} / {"ok":false,"error":"..."}.
func handleRuntimeCallForAndroid(app *App, payload string) string {
	fail := func(msg string) string {
		b, _ := json.Marshal(map[string]any{"ok": false, "error": msg})
		return string(b)
	}

	var call androidRuntimeCallPayload
	if err := json.Unmarshal([]byte(payload), &call); err != nil {
		return fail("unable to parse runtime call: " + err.Error())
	}
	if call.Object == nil {
		return fail("missing object value")
	}
	if call.Method == nil {
		return fail("missing method value")
	}

	resp, err := androidProcessor(app).HandleRuntimeCallWithIDs(context.Background(), &RuntimeRequest{
		Object:            *call.Object,
		Method:            *call.Method,
		Args:              &Args{call.Args},
		WebviewWindowName: call.WindowName,
		ClientID:          call.ClientID,
	})
	if err != nil {
		return fail(err.Error())
	}

	var envelope map[string]any
	if text, ok := resp.(string); ok {
		envelope = map[string]any{"ok": true, "text": text}
	} else {
		envelope = map[string]any{"ok": true, "data": resp}
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return fail("unable to marshal response: " + err.Error())
	}
	return string(b)
}

// handleMessageForAndroid routes a message from JavaScript into the standard
// window message processing pipeline (mirrors HandleJSMessage on iOS).
// Structured payloads carry the message in a "name" or "message" field;
// plain strings (e.g. "wails:runtime:ready") are forwarded as-is.
func handleMessageForAndroid(app *App, message string) {
	if message == "" {
		return
	}

	androidDebugLogf("[JS message] %s", message)

	msg := message
	var msgData map[string]interface{}
	if err := json.Unmarshal([]byte(message), &msgData); err == nil && msgData != nil {
		if name, ok := msgData["name"].(string); ok && name != "" {
			msg = name
		} else if name, ok := msgData["message"].(string); ok && name != "" {
			msg = name
		}
	}

	windowMessageBuffer <- &windowMessage{
		windowId: androidFirstWindowID(app),
		message:  msg,
	}
}

func getMimeTypeForPath(path string) string {
	switch {
	case endsWith(path, ".html"), endsWith(path, ".htm"):
		return "text/html"
	case endsWith(path, ".js"), endsWith(path, ".mjs"):
		return "application/javascript"
	case endsWith(path, ".css"):
		return "text/css"
	case endsWith(path, ".json"), endsWith(path, ".map"):
		return "application/json"
	case endsWith(path, ".png"):
		return "image/png"
	case endsWith(path, ".jpg"), endsWith(path, ".jpeg"):
		return "image/jpeg"
	case endsWith(path, ".gif"):
		return "image/gif"
	case endsWith(path, ".svg"):
		return "image/svg+xml"
	case endsWith(path, ".ico"):
		return "image/x-icon"
	case endsWith(path, ".webp"):
		return "image/webp"
	case endsWith(path, ".woff"):
		return "font/woff"
	case endsWith(path, ".woff2"):
		return "font/woff2"
	case endsWith(path, ".ttf"):
		return "font/ttf"
	case endsWith(path, ".otf"):
		return "font/otf"
	case endsWith(path, ".wasm"):
		return "application/wasm"
	case endsWith(path, ".mp3"):
		return "audio/mpeg"
	case endsWith(path, ".mp4"):
		return "video/mp4"
	case endsWith(path, ".webm"):
		return "video/webm"
	case endsWith(path, ".xml"):
		return "application/xml"
	case endsWith(path, ".txt"):
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
