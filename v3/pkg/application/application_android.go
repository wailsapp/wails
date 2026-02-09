//go:build android && cgo && !server

package application

/*
#include <jni.h>
#include <stdlib.h>
#include <string.h>

// Global JavaVM reference for thread attachment
static JavaVM* g_jvm = NULL;

// Global reference to bridge object (must be a global ref, not local)
static jobject g_bridge = NULL;

// Cached method ID for executeJavaScript
static jmethodID g_executeJsMethod = NULL;

// Helper function to convert Java String to C string
static const char* jstringToC(JNIEnv *env, jstring jstr) {
    if (jstr == NULL) return NULL;
    return (*env)->GetStringUTFChars(env, jstr, NULL);
}

// Helper function to release Java String
static void releaseJString(JNIEnv *env, jstring jstr, const char* cstr) {
    if (jstr != NULL && cstr != NULL) {
        (*env)->ReleaseStringUTFChars(env, jstr, cstr);
    }
}

// Helper function to create Java byte array from C data
static jbyteArray createByteArray(JNIEnv *env, const void* data, int len) {
    if (data == NULL || len <= 0) return NULL;
    jbyteArray arr = (*env)->NewByteArray(env, len);
    if (arr != NULL) {
        (*env)->SetByteArrayRegion(env, arr, 0, len, (const jbyte*)data);
    }
    return arr;
}

// Helper function to create Java String from C string
static jstring createJString(JNIEnv *env, const char* str) {
    if (str == NULL) return NULL;
    return (*env)->NewStringUTF(env, str);
}

// Store JavaVM and create global reference to bridge
static void storeBridgeRef(JNIEnv *env, jobject bridge) {
    // Get JavaVM
    if ((*env)->GetJavaVM(env, &g_jvm) != 0) {
        return;
    }

    // Create global reference to bridge object
    g_bridge = (*env)->NewGlobalRef(env, bridge);
    if (g_bridge == NULL) {
        return;
    }

    // Cache the executeJavaScript method ID
    jclass bridgeClass = (*env)->GetObjectClass(env, g_bridge);
    if (bridgeClass != NULL) {
        g_executeJsMethod = (*env)->GetMethodID(env, bridgeClass, "executeJavaScript", "(Ljava/lang/String;)V");
        (*env)->DeleteLocalRef(env, bridgeClass);
    }
}

// Android logging via __android_log_print
#include <android/log.h>
#define LOGD(...) __android_log_print(ANDROID_LOG_DEBUG, "WailsNative", __VA_ARGS__)

// Execute JavaScript via the bridge - can be called from any thread
static void executeJavaScriptOnBridge(const char* js) {
    LOGD("executeJavaScriptOnBridge called, js length: %d", js ? (int)strlen(js) : -1);

    if (g_jvm == NULL) {
        LOGD("executeJavaScriptOnBridge: g_jvm is NULL");
        return;
    }
    if (g_bridge == NULL) {
        LOGD("executeJavaScriptOnBridge: g_bridge is NULL");
        return;
    }
    if (g_executeJsMethod == NULL) {
        LOGD("executeJavaScriptOnBridge: g_executeJsMethod is NULL");
        return;
    }
    if (js == NULL) {
        LOGD("executeJavaScriptOnBridge: js is NULL");
        return;
    }

    JNIEnv *env = NULL;
    int needsDetach = 0;

    // Get JNIEnv for current thread
    jint result = (*g_jvm)->GetEnv(g_jvm, (void**)&env, JNI_VERSION_1_6);
    LOGD("executeJavaScriptOnBridge: GetEnv result = %d", result);
    if (result == JNI_EDETACHED) {
        // Attach current thread to JVM
        LOGD("executeJavaScriptOnBridge: Attaching thread");
        if ((*g_jvm)->AttachCurrentThread(g_jvm, &env, NULL) != 0) {
            LOGD("executeJavaScriptOnBridge: AttachCurrentThread failed");
            return;
        }
        needsDetach = 1;
    } else if (result != JNI_OK) {
        LOGD("executeJavaScriptOnBridge: GetEnv failed with %d", result);
        return;
    }

    // Create Java string and call method
    jstring jJs = (*env)->NewStringUTF(env, js);
    LOGD("executeJavaScriptOnBridge: jJs created: %p", jJs);
    if (jJs != NULL) {
        LOGD("executeJavaScriptOnBridge: Calling Java method");
        (*env)->CallVoidMethod(env, g_bridge, g_executeJsMethod, jJs);
        LOGD("executeJavaScriptOnBridge: Java method called");
        (*env)->DeleteLocalRef(env, jJs);
    }

    // Check for exceptions
    if ((*env)->ExceptionCheck(env)) {
        LOGD("executeJavaScriptOnBridge: Exception occurred!");
        (*env)->ExceptionDescribe(env);
        (*env)->ExceptionClear(env);
    }

    // Detach if we attached
    if (needsDetach) {
        LOGD("executeJavaScriptOnBridge: Detaching thread");
        (*g_jvm)->DetachCurrentThread(g_jvm);
    }

    LOGD("executeJavaScriptOnBridge: Done");
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

	"github.com/wailsapp/wails/v3/internal/runtime"
)

var (
	// Global reference to the app for JNI callbacks
	globalApp     *App
	globalAppLock sync.RWMutex

	// JNI environment and class references
	javaVM       unsafe.Pointer
	bridgeObject unsafe.Pointer

	// Android main function registration
	androidMainFunc func()
	androidMainLock sync.Mutex

	// App ready signal
	appReady     = make(chan struct{})
	appReadyOnce sync.Once
)

func init() {
	androidLogf("info", "🤖 [application_android.go] init() called")
}

// RegisterAndroidMain registers the main function to be called when the Android app starts.
// This should be called from init() in your main.go file for Android builds.
// Example:
//
//	func init() {
//		application.RegisterAndroidMain(main)
//	}
func RegisterAndroidMain(mainFunc func()) {
	androidMainLock.Lock()
	defer androidMainLock.Unlock()
	androidMainFunc = mainFunc
	androidLogf("info", "🤖 [application_android.go] Android main function registered")
}

// signalAppReady signals that the app is ready to serve requests
func signalAppReady() {
	appReadyOnce.Do(func() {
		close(appReady)
		androidLogf("info", "🤖 [application_android.go] App ready signal sent")
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

func androidLogf(level string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	// For now, just use println - we'll connect to Android's Log.* later
	println(fmt.Sprintf("[Android/%s] %s", level, msg))
}

func (a *App) platformRun() {
	androidLogf("info", "🤖 [application_android.go] platformRun() called")

	// Store global reference for JNI callbacks
	globalAppLock.Lock()
	globalApp = a
	globalAppLock.Unlock()

	// Signal that the app is ready to serve requests
	signalAppReady()

	androidLogf("info", "🤖 [application_android.go] App ready, waiting for Android lifecycle...")

	// Block forever - Android manages the app lifecycle via JNI callbacks
	select {}
}

func (a *App) platformQuit() {
	androidLogf("info", "🤖 [application_android.go] platformQuit() called")
	// Android will handle app termination
}

func (a *App) isDarkMode() bool {
	// TODO: Query Android for dark mode status
	return false
}

func (a *App) isWindows() bool {
	return false
}

// Platform-specific app implementation for Android
type androidApp struct {
	parent *App
}

func newPlatformApp(app *App) *androidApp {
	androidLogf("info", "🤖 [application_android.go] newPlatformApp() called")
	return &androidApp{
		parent: app,
	}
}

func (a *androidApp) run() error {
	androidLogf("info", "🤖 [application_android.go] androidApp.run() called")

	// Emit application started event
	a.parent.Event.Emit("ApplicationStarted")

	a.parent.platformRun()
	return nil
}

func (a *androidApp) destroy() {
	androidLogf("info", "🤖 [application_android.go] androidApp.destroy() called")
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

func (a *androidApp) on(_ uint) {
	// Android event handling
}

func (a *androidApp) setApplicationMenu(_ *Menu) {
	// Android doesn't have application menus in the same way
}

func (a *androidApp) show() {
	// Android manages app visibility
}

func (a *androidApp) showAboutDialog(_ string, _ string, _ []byte) {
	// TODO: Implement Android about dialog
}

func (a *androidApp) getPrimaryScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (a *androidApp) getScreens() ([]*Screen, error) {
	return getScreens()
}

func (a *App) logPlatformInfo() {
	// Log Android platform info
	androidLogf("info", "Platform: Android")
}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"platform": "android",
	}
}

func fatalHandler(errFunc func(error)) {
	// Android fatal handler
}

// JNI Export Functions - Called from Java

//export Java_com_wails_app_WailsBridge_nativeInit
func Java_com_wails_app_WailsBridge_nativeInit(env *C.JNIEnv, obj C.jobject, bridge C.jobject) {
	androidLogf("info", "🤖 [JNI] nativeInit called")

	// Store references for later use (legacy - keeping for compatibility)
	javaVM = unsafe.Pointer(env)
	bridgeObject = unsafe.Pointer(bridge)

	// Store JavaVM and bridge global reference for JNI callbacks
	C.storeBridgeRef(env, bridge)
	androidLogf("info", "🤖 [JNI] Bridge reference stored for JNI callbacks")

	// Start the registered main function in a goroutine
	androidMainLock.Lock()
	mainFunc := androidMainFunc
	androidMainLock.Unlock()

	if mainFunc != nil {
		androidLogf("info", "🤖 [JNI] Starting registered main function in goroutine")
		go mainFunc()
	} else {
		androidLogf("warn", "🤖 [JNI] No main function registered! Call application.RegisterAndroidMain(main) in init()")
	}

	androidLogf("info", "🤖 [JNI] nativeInit complete")
}

//export Java_com_wails_app_WailsBridge_nativeShutdown
func Java_com_wails_app_WailsBridge_nativeShutdown(env *C.JNIEnv, obj C.jobject) {
	androidLogf("info", "🤖 [JNI] nativeShutdown called")

	globalAppLock.Lock()
	if globalApp != nil {
		globalApp.Quit()
	}
	globalAppLock.Unlock()
}

//export Java_com_wails_app_WailsBridge_nativeOnResume
func Java_com_wails_app_WailsBridge_nativeOnResume(env *C.JNIEnv, obj C.jobject) {
	androidLogf("info", "🤖 [JNI] nativeOnResume called")

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app != nil {
		app.Event.Emit("ApplicationResumed")
	}
}

//export Java_com_wails_app_WailsBridge_nativeOnPause
func Java_com_wails_app_WailsBridge_nativeOnPause(env *C.JNIEnv, obj C.jobject) {
	androidLogf("info", "🤖 [JNI] nativeOnPause called")

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app != nil {
		app.Event.Emit("ApplicationPaused")
	}
}

//export Java_com_wails_app_WailsBridge_nativeOnPageFinished
func Java_com_wails_app_WailsBridge_nativeOnPageFinished(env *C.JNIEnv, obj C.jobject, jurl C.jstring) {
	cUrl := C.jstringToC(env, jurl)
	defer C.releaseJString(env, jurl, cUrl)
	url := C.GoString(cUrl)

	androidLogf("info", "🤖 [JNI] nativeOnPageFinished called: %s", url)

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app == nil {
		androidLogf("error", "🤖 [JNI] nativeOnPageFinished: app is nil")
		return
	}

	// Inject the runtime into the first window (with proper locking)
	app.windowsLock.RLock()
	windowCount := len(app.windows)
	androidLogf("info", "🤖 [JNI] nativeOnPageFinished: window count = %d", windowCount)
	for id, win := range app.windows {
		androidLogf("info", "🤖 [JNI] Found window ID: %d", id)
		if win != nil {
			androidLogf("info", "🤖 [JNI] Injecting runtime.Core() into window %d", id)
			// Get the runtime core JavaScript
			flags := map[string]any{}
			if app.impl != nil {
				flags = app.impl.GetFlags(app.options)
			}
			runtimeJS := runtime.Core(flags)
			androidLogf("info", "🤖 [JNI] Runtime JS length: %d bytes", len(runtimeJS))
			app.windowsLock.RUnlock()
			// IMPORTANT: We must bypass win.ExecJS because it queues if runtimeLoaded is false.
			// On Android, we need to inject the runtime directly since the runtime hasn't been loaded yet.
			// This is the bootstrap injection that enables the runtime to load.
			androidLogf("info", "🤖 [JNI] Calling executeJavaScript directly (bypassing queue)")
			executeJavaScript(runtimeJS)
			// Emit event
			app.Event.Emit("PageFinished", url)
			return
		}
	}
	app.windowsLock.RUnlock()

	androidLogf("warn", "🤖 [JNI] nativeOnPageFinished: no windows found to inject runtime")
	// Emit event even if no windows
	app.Event.Emit("PageFinished", url)
}

//export Java_com_wails_app_WailsBridge_nativeServeAsset
func Java_com_wails_app_WailsBridge_nativeServeAsset(env *C.JNIEnv, obj C.jobject, jpath C.jstring, jmethod C.jstring, jheaders C.jstring) C.jbyteArray {
	// Convert Java strings to Go strings
	cPath := C.jstringToC(env, jpath)
	cMethod := C.jstringToC(env, jmethod)
	defer C.releaseJString(env, jpath, cPath)
	defer C.releaseJString(env, jmethod, cMethod)

	goPath := C.GoString(cPath)
	goMethod := C.GoString(cMethod)

	androidLogf("debug", "🤖 [JNI] nativeServeAsset: %s %s", goMethod, goPath)

	// Wait for the app to be ready (timeout after 10 seconds)
	if !waitForAppReady(10 * time.Second) {
		androidLogf("error", "🤖 [JNI] Timeout waiting for app to be ready")
		return C.createByteArray(env, nil, 0)
	}

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app == nil || app.assets == nil {
		androidLogf("error", "🤖 [JNI] App or assets not initialized after ready signal")
		return C.createByteArray(env, nil, 0)
	}

	// Serve the asset through the asset server
	data, err := serveAssetForAndroid(app, goPath)
	if err != nil {
		androidLogf("error", "🤖 [JNI] Error serving asset %s: %v", goPath, err)
		return C.createByteArray(env, nil, 0)
	}

	androidLogf("debug", "🤖 [JNI] Serving asset %s (%d bytes)", goPath, len(data))

	// Create Java byte array from the data
	// Handle empty data case to avoid index out of range panic
	if len(data) == 0 {
		return C.createByteArray(env, nil, 0)
	}
	return C.createByteArray(env, unsafe.Pointer(&data[0]), C.int(len(data)))
}

//export Java_com_wails_app_WailsBridge_nativeHandleMessage
func Java_com_wails_app_WailsBridge_nativeHandleMessage(env *C.JNIEnv, obj C.jobject, jmessage C.jstring) C.jstring {
	// Convert Java string to Go string
	cMessage := C.jstringToC(env, jmessage)
	defer C.releaseJString(env, jmessage, cMessage)

	goMessage := C.GoString(cMessage)

	androidLogf("debug", "🤖 [JNI] nativeHandleMessage: %s", goMessage)

	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()

	if app == nil {
		errorResponse := `{"error":"App not initialized"}`
		return C.createJString(env, C.CString(errorResponse))
	}

	// Parse and handle the message
	response := handleMessageForAndroid(app, goMessage)
	return C.createJString(env, C.CString(response))
}

//export Java_com_wails_app_WailsBridge_nativeGetAssetMimeType
func Java_com_wails_app_WailsBridge_nativeGetAssetMimeType(env *C.JNIEnv, obj C.jobject, jpath C.jstring) C.jstring {
	// Convert Java string to Go string
	cPath := C.jstringToC(env, jpath)
	defer C.releaseJString(env, jpath, cPath)

	goPath := C.GoString(cPath)
	mimeType := getMimeTypeForPath(goPath)
	return C.createJString(env, C.CString(mimeType))
}

// Helper functions

func serveAssetForAndroid(app *App, path string) ([]byte, error) {
	// Check if this is a runtime call (includes query string)
	isRuntimeCall := strings.HasPrefix(path, "/wails/runtime")

	// Normalize path for regular assets (not runtime calls)
	if !isRuntimeCall {
		if path == "" || path == "/" {
			path = "/index.html"
		}
	}

	// Ensure path starts with /
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	// Check if asset server is available
	if app.assets == nil {
		return nil, fmt.Errorf("asset server not initialized")
	}

	// Create a fake HTTP request
	fullURL := "https://wails.localhost" + path
	androidLogf("debug", "🤖 [serveAssetForAndroid] Creating request for: %s", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// For runtime calls (/wails/runtime), we need to add the window ID header
	// This is required by the MessageProcessor to route the call correctly
	if isRuntimeCall {
		// Get the first window (on Android, there's typically only one)
		windows := app.Window.GetAll()
		androidLogf("debug", "🤖 [serveAssetForAndroid] Runtime call, found %d windows", len(windows))
		if len(windows) > 0 {
			// Use the first window's ID
			windowID := windows[0].ID()
			req.Header.Set("x-wails-window-id", fmt.Sprintf("%d", windowID))
			androidLogf("debug", "🤖 [serveAssetForAndroid] Added window ID header: %d", windowID)
		} else {
			androidLogf("warn", "🤖 [serveAssetForAndroid] No windows available for runtime call")
		}
	}

	// Use httptest.ResponseRecorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request through the asset server
	app.assets.ServeHTTP(recorder, req)

	// Check response status
	result := recorder.Result()
	defer result.Body.Close()

	// Read the response body
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	androidLogf("debug", "🤖 [serveAssetForAndroid] Response status: %d, body length: %d", result.StatusCode, len(body))

	// For runtime calls, we need to return the body even for error responses
	// so the JavaScript can see the error message
	if isRuntimeCall {
		if result.StatusCode != http.StatusOK {
			androidLogf("warn", "🤖 [serveAssetForAndroid] Runtime call returned status %d: %s", result.StatusCode, string(body))
		}
		// Return the body regardless of status - the JS will handle errors
		return body, nil
	}

	// For regular assets, check status codes
	if result.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("asset not found: %s", path)
	}

	if result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("asset server error: status %d for %s", result.StatusCode, path)
	}

	return body, nil
}

func handleMessageForAndroid(app *App, message string) string {
	if strings.HasPrefix(message, "wails:") {
		windows := app.Window.GetAll()
		if len(windows) > 0 {
			windows[0].HandleMessage(message)
		}
		return marshalAndroidInvokeResponse(nil, nil)
	}

	var payload androidRuntimeMessage
	if err := json.Unmarshal([]byte(message), &payload); err != nil {
		return marshalAndroidInvokeResponse(nil, fmt.Errorf("invalid message: %w", err))
	}

	if payload.Type != "runtime" {
		return marshalAndroidInvokeResponse(nil, fmt.Errorf("unsupported message type: %s", payload.Type))
	}

	processor := NewMessageProcessor(app.Logger)
	request := &RuntimeRequest{
		Object:            payload.Object,
		Method:            payload.Method,
		Args:              &Args{payload.Args},
		WebviewWindowName: payload.WindowName,
		ClientID:          payload.ClientID,
	}

	response, err := processor.HandleRuntimeCallWithIDs(context.Background(), request)
	return marshalAndroidInvokeResponse(response, err)
}

type androidRuntimeMessage struct {
	Type       string          `json:"type"`
	Object     int             `json:"object"`
	Method     int             `json:"method"`
	Args       json.RawMessage `json:"args,omitempty"`
	WindowName string          `json:"windowName,omitempty"`
	ClientID   string          `json:"clientId,omitempty"`
}

type androidInvokeResponse struct {
	Ok    bool   `json:"ok"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func marshalAndroidInvokeResponse(data any, err error) string {
	response := androidInvokeResponse{Ok: err == nil}
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Data = data
	}

	payload, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return fmt.Sprintf(`{"ok":false,"error":"%s"}`, marshalErr.Error())
	}
	return string(payload)
}

func getMimeTypeForPath(path string) string {
	// Simple MIME type detection based on extension
	switch {
	case endsWith(path, ".html"), endsWith(path, ".htm"):
		return "text/html"
	case endsWith(path, ".js"), endsWith(path, ".mjs"):
		return "application/javascript"
	case endsWith(path, ".css"):
		return "text/css"
	case endsWith(path, ".json"):
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
	case endsWith(path, ".woff"):
		return "font/woff"
	case endsWith(path, ".woff2"):
		return "font/woff2"
	case endsWith(path, ".ttf"):
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// executeJavaScript executes JavaScript code in the Android WebView via JNI callback
func executeJavaScript(js string) {
	androidLogf("info", "🤖 executeJavaScript called, length: %d", len(js))
	if js == "" {
		androidLogf("warn", "🤖 executeJavaScript: empty JS string")
		return
	}

	// Convert Go string to C string and call the JNI bridge
	androidLogf("info", "🤖 executeJavaScript: calling C.executeJavaScriptOnBridge")
	cJs := C.CString(js)
	defer C.free(unsafe.Pointer(cJs))

	C.executeJavaScriptOnBridge(cJs)
	androidLogf("info", "🤖 executeJavaScript: done")
}
