//go:build android && cgo

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

// Cached method IDs for Android-specific features
static jmethodID g_vibrateMethod = NULL;
static jmethodID g_showToastMethod = NULL;
static jmethodID g_getDeviceInfoMethod = NULL;
static jmethodID g_openURLMethod = NULL;
static jmethodID g_setWebViewBackgroundColorMethod = NULL;
static jmethodID g_isDarkModeMethod = NULL;
static jmethodID g_getScreenInfoMethod = NULL;
static jmethodID g_setClipboardTextMethod = NULL;
static jmethodID g_getClipboardTextMethod = NULL;
static jmethodID g_showMessageDialogMethod = NULL;
static jmethodID g_setHTMLMethod = NULL;
static jmethodID g_setURLMethod = NULL;

// Helper function to get JNIEnv for current thread
static JNIEnv* getEnv(int *needsDetach) {
    *needsDetach = 0;
    if (g_jvm == NULL) return NULL;

    JNIEnv *env = NULL;
    jint result = (*g_jvm)->GetEnv(g_jvm, (void**)&env, JNI_VERSION_1_6);
    if (result == JNI_EDETACHED) {
        if ((*g_jvm)->AttachCurrentThread(g_jvm, &env, NULL) != 0) {
            return NULL;
        }
        *needsDetach = 1;
    } else if (result != JNI_OK) {
        return NULL;
    }
    return env;
}

// Helper function to detach from JVM if needed
static void releaseEnv(int needsDetach) {
    if (needsDetach && g_jvm != NULL) {
        (*g_jvm)->DetachCurrentThread(g_jvm);
    }
}

// Cache method IDs for Android-specific features
static void cacheAndroidMethods(JNIEnv *env) {
    if (g_bridge == NULL) return;

    jclass bridgeClass = (*env)->GetObjectClass(env, g_bridge);
    if (bridgeClass == NULL) return;

    // Cache method IDs if not already cached
    if (g_vibrateMethod == NULL) {
        g_vibrateMethod = (*env)->GetMethodID(env, bridgeClass, "vibrate", "(I)V");
    }
    if (g_showToastMethod == NULL) {
        g_showToastMethod = (*env)->GetMethodID(env, bridgeClass, "showToast", "(Ljava/lang/String;)V");
    }
    if (g_getDeviceInfoMethod == NULL) {
        g_getDeviceInfoMethod = (*env)->GetMethodID(env, bridgeClass, "getDeviceInfo", "()Ljava/lang/String;");
    }
    if (g_openURLMethod == NULL) {
        g_openURLMethod = (*env)->GetMethodID(env, bridgeClass, "openURL", "(Ljava/lang/String;)Z");
    }
    if (g_setWebViewBackgroundColorMethod == NULL) {
        g_setWebViewBackgroundColorMethod = (*env)->GetMethodID(env, bridgeClass, "setWebViewBackgroundColor", "(I)V");
    }
    if (g_isDarkModeMethod == NULL) {
        g_isDarkModeMethod = (*env)->GetMethodID(env, bridgeClass, "isDarkMode", "()Z");
    }
    if (g_getScreenInfoMethod == NULL) {
        g_getScreenInfoMethod = (*env)->GetMethodID(env, bridgeClass, "getScreenInfo", "()Ljava/lang/String;");
    }
    if (g_setClipboardTextMethod == NULL) {
        g_setClipboardTextMethod = (*env)->GetMethodID(env, bridgeClass, "setClipboardText", "(Ljava/lang/String;)V");
    }
    if (g_getClipboardTextMethod == NULL) {
        g_getClipboardTextMethod = (*env)->GetMethodID(env, bridgeClass, "getClipboardText", "()Ljava/lang/String;");
    }
    if (g_showMessageDialogMethod == NULL) {
        g_showMessageDialogMethod = (*env)->GetMethodID(env, bridgeClass, "showMessageDialog", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;");
    }
    if (g_setHTMLMethod == NULL) {
        g_setHTMLMethod = (*env)->GetMethodID(env, bridgeClass, "setHTML", "(Ljava/lang/String;)V");
    }
    if (g_setURLMethod == NULL) {
        g_setURLMethod = (*env)->GetMethodID(env, bridgeClass, "setURL", "(Ljava/lang/String;)V");
    }

    (*env)->DeleteLocalRef(env, bridgeClass);
}

// Call Android Vibrator service
static void androidVibrate(int durationMs) {
    LOGD("androidVibrate called: %dms", durationMs);

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidVibrate: env or bridge is NULL");
        return;
    }

    // Ensure method is cached
    if (g_vibrateMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_vibrateMethod != NULL) {
        (*env)->CallVoidMethod(env, g_bridge, g_vibrateMethod, (jint)durationMs);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        }
    }

    releaseEnv(needsDetach);
}

// Show Android Toast
static void androidShowToastNative(const char* message) {
    LOGD("androidShowToastNative called: %s", message ? message : "null");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || message == NULL) {
        LOGD("androidShowToastNative: env, bridge, or message is NULL");
        return;
    }

    // Ensure method is cached
    if (g_showToastMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_showToastMethod != NULL) {
        jstring jMessage = (*env)->NewStringUTF(env, message);
        if (jMessage != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, g_showToastMethod, jMessage);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            }
            (*env)->DeleteLocalRef(env, jMessage);
        }
    }

    releaseEnv(needsDetach);
}

// Get Android device info
static const char* androidGetDeviceInfoNative() {
    LOGD("androidGetDeviceInfoNative called");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidGetDeviceInfoNative: env or bridge is NULL");
        return NULL;
    }

    // Ensure method is cached
    if (g_getDeviceInfoMethod == NULL) {
        cacheAndroidMethods(env);
    }

    const char* result = NULL;
    if (g_getDeviceInfoMethod != NULL) {
        jstring jResult = (jstring)(*env)->CallObjectMethod(env, g_bridge, g_getDeviceInfoMethod);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        } else if (jResult != NULL) {
            const char* cResult = (*env)->GetStringUTFChars(env, jResult, NULL);
            if (cResult != NULL) {
                // Make a copy since we need to release the JNI string
                result = strdup(cResult);
                (*env)->ReleaseStringUTFChars(env, jResult, cResult);
            }
            (*env)->DeleteLocalRef(env, jResult);
        }
    }

    releaseEnv(needsDetach);
    return result;
}

// Open URL via Android Intent
static int androidOpenURLNative(const char* url) {
    LOGD("androidOpenURLNative called: %s", url ? url : "null");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || url == NULL) {
        LOGD("androidOpenURLNative: env, bridge, or url is NULL");
        return 0;
    }

    // Ensure method is cached
    if (g_openURLMethod == NULL) {
        cacheAndroidMethods(env);
    }

    int result = 0;
    if (g_openURLMethod != NULL) {
        jstring jUrl = (*env)->NewStringUTF(env, url);
        if (jUrl != NULL) {
            jboolean success = (*env)->CallBooleanMethod(env, g_bridge, g_openURLMethod, jUrl);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            } else {
                result = success ? 1 : 0;
            }
            (*env)->DeleteLocalRef(env, jUrl);
        }
    }

    releaseEnv(needsDetach);
    LOGD("androidOpenURLNative result: %d", result);
    return result;
}

// Set WebView background color
static void androidSetWebViewBackgroundColorNative(int color) {
    LOGD("androidSetWebViewBackgroundColorNative called: 0x%08x", color);

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidSetWebViewBackgroundColorNative: env or bridge is NULL");
        return;
    }

    // Ensure method is cached
    if (g_setWebViewBackgroundColorMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_setWebViewBackgroundColorMethod != NULL) {
        (*env)->CallVoidMethod(env, g_bridge, g_setWebViewBackgroundColorMethod, (jint)color);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        }
    }

    releaseEnv(needsDetach);
}

// Check if dark mode is enabled
static int androidIsDarkModeNative() {
    LOGD("androidIsDarkModeNative called");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidIsDarkModeNative: env or bridge is NULL");
        return 0;
    }

    // Ensure method is cached
    if (g_isDarkModeMethod == NULL) {
        cacheAndroidMethods(env);
    }

    int result = 0;
    if (g_isDarkModeMethod != NULL) {
        jboolean isDark = (*env)->CallBooleanMethod(env, g_bridge, g_isDarkModeMethod);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        } else {
            result = isDark ? 1 : 0;
        }
    }

    releaseEnv(needsDetach);
    LOGD("androidIsDarkModeNative result: %d", result);
    return result;
}

// Get Android screen info
static const char* androidGetScreenInfoNative() {
    LOGD("androidGetScreenInfoNative called");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidGetScreenInfoNative: env or bridge is NULL");
        return NULL;
    }

    // Ensure method is cached
    if (g_getScreenInfoMethod == NULL) {
        cacheAndroidMethods(env);
    }

    const char* result = NULL;
    if (g_getScreenInfoMethod != NULL) {
        jstring jResult = (jstring)(*env)->CallObjectMethod(env, g_bridge, g_getScreenInfoMethod);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        } else if (jResult != NULL) {
            const char* cResult = (*env)->GetStringUTFChars(env, jResult, NULL);
            if (cResult != NULL) {
                // Make a copy since we need to release the JNI string
                result = strdup(cResult);
                (*env)->ReleaseStringUTFChars(env, jResult, cResult);
            }
            (*env)->DeleteLocalRef(env, jResult);
        }
    }

    releaseEnv(needsDetach);
    return result;
}

// Call Android ClipboardManager to set text
static void androidSetClipboardTextNative(const char* text) {
    LOGD("androidSetClipboardTextNative called: %s", text ? text : "null");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || text == NULL) {
        LOGD("androidSetClipboardTextNative: env, bridge, or text is NULL");
        return;
    }

    // Ensure method is cached
    if (g_setClipboardTextMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_setClipboardTextMethod != NULL) {
        jstring jText = (*env)->NewStringUTF(env, text);
        if (jText != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, g_setClipboardTextMethod, jText);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            }
            (*env)->DeleteLocalRef(env, jText);
        }
    }

    releaseEnv(needsDetach);
}

// Get text from Android ClipboardManager
static const char* androidGetClipboardTextNative() {
    LOGD("androidGetClipboardTextNative called");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL) {
        LOGD("androidGetClipboardTextNative: env or bridge is NULL");
        return NULL;
    }

    // Ensure method is cached
    if (g_getClipboardTextMethod == NULL) {
        cacheAndroidMethods(env);
    }

    const char* result = NULL;
    if (g_getClipboardTextMethod != NULL) {
        jstring jResult = (jstring)(*env)->CallObjectMethod(env, g_bridge, g_getClipboardTextMethod);
        if ((*env)->ExceptionCheck(env)) {
            (*env)->ExceptionDescribe(env);
            (*env)->ExceptionClear(env);
        } else if (jResult != NULL) {
            const char* cResult = (*env)->GetStringUTFChars(env, jResult, NULL);
            if (cResult != NULL) {
                // Make a copy since we need to release the JNI string
                result = strdup(cResult);
                (*env)->ReleaseStringUTFChars(env, jResult, cResult);
            }
            (*env)->DeleteLocalRef(env, jResult);
        }
    }

    releaseEnv(needsDetach);
    return result;
}

// Show message dialog via Android AlertDialog
static const char* androidShowMessageDialogNative(const char* dialogType, const char* title, const char* message, const char* buttons) {
    LOGD("androidShowMessageDialogNative called: type=%s, title=%s", dialogType ? dialogType : "null", title ? title : "null");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || dialogType == NULL || title == NULL || message == NULL || buttons == NULL) {
        LOGD("androidShowMessageDialogNative: env, bridge, or parameters are NULL");
        return NULL;
    }

    // Ensure method is cached
    if (g_showMessageDialogMethod == NULL) {
        cacheAndroidMethods(env);
    }

    const char* result = NULL;
    if (g_showMessageDialogMethod != NULL) {
        jstring jDialogType = (*env)->NewStringUTF(env, dialogType);
        jstring jTitle = (*env)->NewStringUTF(env, title);
        jstring jMessage = (*env)->NewStringUTF(env, message);
        jstring jButtons = (*env)->NewStringUTF(env, buttons);

        if (jDialogType != NULL && jTitle != NULL && jMessage != NULL && jButtons != NULL) {
            jstring jResult = (jstring)(*env)->CallObjectMethod(env, g_bridge, g_showMessageDialogMethod, jDialogType, jTitle, jMessage, jButtons);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            } else if (jResult != NULL) {
                const char* cResult = (*env)->GetStringUTFChars(env, jResult, NULL);
                if (cResult != NULL) {
                    // Make a copy since we need to release the JNI string
                    result = strdup(cResult);
                    (*env)->ReleaseStringUTFChars(env, jResult, cResult);
                }
                (*env)->DeleteLocalRef(env, jResult);
            }
            (*env)->DeleteLocalRef(env, jDialogType);
            (*env)->DeleteLocalRef(env, jTitle);
            (*env)->DeleteLocalRef(env, jMessage);
            (*env)->DeleteLocalRef(env, jButtons);
        }
    }

    releaseEnv(needsDetach);
    return result;
}

// Set HTML content in WebView
static void androidSetHTMLNative(const char* html) {
    LOGD("androidSetHTMLNative called");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || html == NULL) {
        LOGD("androidSetHTMLNative: env, bridge, or html is NULL");
        return;
    }

    // Ensure method is cached
    if (g_setHTMLMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_setHTMLMethod != NULL) {
        jstring jHtml = (*env)->NewStringUTF(env, html);
        if (jHtml != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, g_setHTMLMethod, jHtml);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            }
            (*env)->DeleteLocalRef(env, jHtml);
        }
    }

    releaseEnv(needsDetach);
}

// Set URL in WebView
static void androidSetURLNative(const char* url) {
    LOGD("androidSetURLNative called: %s", url ? url : "null");

    int needsDetach = 0;
    JNIEnv *env = getEnv(&needsDetach);
    if (env == NULL || g_bridge == NULL || url == NULL) {
        LOGD("androidSetURLNative: env, bridge, or url is NULL");
        return;
    }

    // Ensure method is cached
    if (g_setURLMethod == NULL) {
        cacheAndroidMethods(env);
    }

    if (g_setURLMethod != NULL) {
        jstring jUrl = (*env)->NewStringUTF(env, url);
        if (jUrl != NULL) {
            (*env)->CallVoidMethod(env, g_bridge, g_setURLMethod, jUrl);
            if ((*env)->ExceptionCheck(env)) {
                (*env)->ExceptionDescribe(env);
                (*env)->ExceptionClear(env);
            }
            (*env)->DeleteLocalRef(env, jUrl);
        }
    }

    releaseEnv(needsDetach);
}

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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/browser"
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
	// Register the Android OpenURL implementation with the browser package
	browser.OpenURLFunc = AndroidOpenURL
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
	result := C.androidIsDarkModeNative()
	return result != 0
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
			runtimeJS := runtime.Core()
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
	androidLogf("debug", "🤖 [handleMessageForAndroid] Received message: %s", message)

	// Check for special system messages that don't need JSON parsing
	// These are sent directly as strings from the frontend runtime
	if strings.HasPrefix(message, "wails:") {
		androidLogf("info", "🤖 [handleMessageForAndroid] System message detected: %s", message)

		// Route to the first window's HandleMessage (Android typically has only one window)
		app.windowsLock.RLock()
		for _, win := range app.windows {
			if win != nil {
				app.windowsLock.RUnlock()
				androidLogf("info", "🤖 [handleMessageForAndroid] Routing to window.HandleMessage: %s", message)
				win.HandleMessage(message)
				return `{"success":true}`
			}
		}
		app.windowsLock.RUnlock()

		androidLogf("warn", "🤖 [handleMessageForAndroid] No windows available to handle message")
		return `{"error":"No windows available"}`
	}

	// Try to parse as JSON for other message types
	var msg map[string]interface{}
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		// Not JSON and not a system message - treat as an unknown message
		androidLogf("warn", "🤖 [handleMessageForAndroid] Unknown message format: %s", message)
		return fmt.Sprintf(`{"error":"Unknown message format: %s"}`, err.Error())
	}

	// Handle JSON-based messages (e.g., method calls, events from JS)
	androidLogf("debug", "🤖 [handleMessageForAndroid] JSON message: %v", msg)

	// Check for event emission from JS to Go
	if eventName, ok := msg["event"].(string); ok {
		androidLogf("info", "🤖 [handleMessageForAndroid] Event from JS: %s", eventName)
		// Get event data if present
		eventData := msg["data"]
		app.Event.Emit(eventName, eventData)
		return `{"success":true}`
	}

	// For other JSON messages, return success for now
	// Additional message types can be added here as needed
	return `{"success":true}`
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

// ==================== Android Platform Features ====================
// These Go wrapper functions call the C/JNI functions and can be
// called from other Go files in the package (e.g., messageprocessor_android.go)

// AndroidVibrate triggers haptic feedback via JNI
func AndroidVibrate(durationMs int) {
	androidLogf("debug", "AndroidVibrate: %dms", durationMs)
	C.androidVibrate(C.int(durationMs))
}

// AndroidShowToast shows a native Android toast notification via JNI
func AndroidShowToast(message string) {
	androidLogf("debug", "AndroidShowToast: %s", message)
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.androidShowToastNative(cMessage)
}

// AndroidGetDeviceInfo returns device information as a JSON string via JNI
func AndroidGetDeviceInfo() string {
	androidLogf("debug", "AndroidGetDeviceInfo called")
	cResult := C.androidGetDeviceInfoNative()
	if cResult == nil {
		return `{"platform":"android","model":"Unknown","version":"Unknown"}`
	}
	result := C.GoString(cResult)
	C.free(unsafe.Pointer(cResult))
	return result
}

// AndroidOpenURL opens a URL using Android's Intent system via JNI
func AndroidOpenURL(url string) error {
	androidLogf("debug", "AndroidOpenURL: %s", url)
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	result := C.androidOpenURLNative(cURL)
	if result == 0 {
		return fmt.Errorf("failed to open URL: %s", url)
	}
	return nil
}

// AndroidGetScreenInfo returns screen information as a JSON string via JNI
func AndroidGetScreenInfo() string {
	androidLogf("debug", "AndroidGetScreenInfo called")
	cResult := C.androidGetScreenInfoNative()
	if cResult == nil {
		return `{"widthPixels":1080,"heightPixels":2400,"density":2.0}`
	}
	result := C.GoString(cResult)
	C.free(unsafe.Pointer(cResult))
	return result
}

// AndroidSetWebViewBackgroundColor sets the WebView background color via JNI
// Android uses ARGB format: (alpha << 24) | (red << 16) | (green << 8) | blue
func AndroidSetWebViewBackgroundColor(r, g, b, a uint8) {
	androidLogf("debug", "AndroidSetWebViewBackgroundColor: rgba(%d,%d,%d,%d)", r, g, b, a)
	// Convert RGBA to Android ARGB int format
	color := (int(a) << 24) | (int(r) << 16) | (int(g) << 8) | int(b)
	androidLogf("debug", "AndroidSetWebViewBackgroundColor: color=0x%08x", color)
	C.androidSetWebViewBackgroundColorNative(C.int(color))
}

// AndroidSetClipboardText sets clipboard text via JNI
func AndroidSetClipboardText(text string) {
	androidLogf("debug", "AndroidSetClipboardText: %s", text)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.androidSetClipboardTextNative(cText)
}

// AndroidGetClipboardText gets clipboard text via JNI
func AndroidGetClipboardText() string {
	androidLogf("debug", "AndroidGetClipboardText called")
	cResult := C.androidGetClipboardTextNative()
	if cResult == nil {
		return ""
	}
	result := C.GoString(cResult)
	C.free(unsafe.Pointer(cResult))
	return result
}

// AndroidShowMessageDialog shows a message dialog via JNI
// Returns the label of the button that was clicked, or empty string if cancelled
func AndroidShowMessageDialog(dialogType, title, message, buttons string) string {
	androidLogf("debug", "AndroidShowMessageDialog: type=%s, title=%s", dialogType, title)
	cDialogType := C.CString(dialogType)
	cTitle := C.CString(title)
	cMessage := C.CString(message)
	cButtons := C.CString(buttons)
	defer C.free(unsafe.Pointer(cDialogType))
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cMessage))
	defer C.free(unsafe.Pointer(cButtons))

	cResult := C.androidShowMessageDialogNative(cDialogType, cTitle, cMessage, cButtons)
	if cResult == nil {
		return ""
	}
	result := C.GoString(cResult)
	C.free(unsafe.Pointer(cResult))
	androidLogf("debug", "AndroidShowMessageDialog result: %s", result)
	return result
}

// AndroidSetHTML loads HTML content into the WebView via JNI
func AndroidSetHTML(html string) {
	androidLogf("debug", "AndroidSetHTML called")
	cHtml := C.CString(html)
	defer C.free(unsafe.Pointer(cHtml))
	C.androidSetHTMLNative(cHtml)
}

// AndroidSetURL loads a URL into the WebView via JNI
func AndroidSetURL(url string) {
	androidLogf("debug", "AndroidSetURL: %s", url)
	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))
	C.androidSetURLNative(cUrl)
}
