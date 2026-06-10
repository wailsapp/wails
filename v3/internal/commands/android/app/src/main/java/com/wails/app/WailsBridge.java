package com.wails.app;

import android.content.Context;
import android.util.Log;
import android.webkit.WebView;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * WailsBridge manages the connection between the Java/Android side and the Go native library.
 * It handles:
 * - Loading and initializing the native Go library
 * - Serving asset requests from Go
 * - Passing messages between JavaScript and Go
 * - Managing callbacks for async operations
 */
public class WailsBridge {
    private static final String TAG = "WailsBridge";

    static {
        // Load the native Go library
        System.loadLibrary("wails");
    }

    private final Context context;
    private final AtomicInteger callbackIdGenerator = new AtomicInteger(0);
    private final ConcurrentHashMap<Integer, AssetCallback> pendingAssetCallbacks = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<Integer, MessageCallback> pendingMessageCallbacks = new ConcurrentHashMap<>();
    private WebView webView;
    private volatile boolean initialized = false;

    // Native methods - implemented in Go
    private static native void nativeInit(WailsBridge bridge);
    private static native void nativeShutdown();
    private static native void nativeOnResume();
    private static native void nativeOnPause();
    private static native void nativeOnPageFinished(String url);
    private static native byte[] nativeServeAsset(String path, String method, String headers);
    private static native String nativeHandleMessage(String message);
    private static native String nativeGetAssetMimeType(String path);

    public WailsBridge(Context context) {
        this.context = context;
    }

    /**
     * Initialize the native Go library
     */
    public void initialize() {
        if (initialized) {
            return;
        }

        Log.i(TAG, "Initializing Wails bridge...");
        try {
            nativeInit(this);
            initialized = true;
            Log.i(TAG, "Wails bridge initialized successfully");
        } catch (Exception e) {
            Log.e(TAG, "Failed to initialize Wails bridge", e);
        }
    }

    /**
     * Shutdown the native Go library
     */
    public void shutdown() {
        if (!initialized) {
            return;
        }

        Log.i(TAG, "Shutting down Wails bridge...");
        try {
            nativeShutdown();
            initialized = false;
        } catch (Exception e) {
            Log.e(TAG, "Error during shutdown", e);
        }
    }

    /**
     * Called when the activity resumes
     */
    public void onResume() {
        if (initialized) {
            nativeOnResume();
        }
    }

    /**
     * Called when the activity pauses
     */
    public void onPause() {
        if (initialized) {
            nativeOnPause();
        }
    }

    /**
     * Serve an asset from the Go asset server
     * @param path The URL path requested
     * @param method The HTTP method
     * @param headers The request headers as JSON
     * @return The asset data, or null if not found
     */
    public byte[] serveAsset(String path, String method, String headers) {
        if (!initialized) {
            Log.w(TAG, "Bridge not initialized, cannot serve asset: " + path);
            return null;
        }

        Log.d(TAG, "Serving asset: " + path);
        try {
            return nativeServeAsset(path, method, headers);
        } catch (Exception e) {
            Log.e(TAG, "Error serving asset: " + path, e);
            return null;
        }
    }

    /**
     * Get the MIME type for an asset
     * @param path The asset path
     * @return The MIME type string
     */
    public String getAssetMimeType(String path) {
        if (!initialized) {
            return "application/octet-stream";
        }

        try {
            String mimeType = nativeGetAssetMimeType(path);
            return mimeType != null ? mimeType : "application/octet-stream";
        } catch (Exception e) {
            Log.e(TAG, "Error getting MIME type for: " + path, e);
            return "application/octet-stream";
        }
    }

    /**
     * Handle a message from JavaScript
     * @param message The message from JavaScript (JSON)
     * @return The response to send back to JavaScript (JSON)
     */
    public String handleMessage(String message) {
        if (!initialized) {
            Log.w(TAG, "Bridge not initialized, cannot handle message");
            return "{\"error\":\"Bridge not initialized\"}";
        }

        Log.d(TAG, "Handling message from JS: " + message);
        try {
            return nativeHandleMessage(message);
        } catch (Exception e) {
            Log.e(TAG, "Error handling message", e);
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    /**
     * Inject the Wails runtime JavaScript into the WebView.
     * Called when the page finishes loading.
     * @param webView The WebView to inject into
     * @param url The URL that finished loading
     */
    public void injectRuntime(WebView webView, String url) {
        this.webView = webView;
        // Notify Go side that page has finished loading so it can inject the runtime
        Log.d(TAG, "Page finished loading: " + url + ", notifying Go side");
        if (initialized) {
            nativeOnPageFinished(url);
        }
    }

    /**
     * Execute JavaScript in the WebView (called from Go side)
     * @param js The JavaScript code to execute
     */
    public void executeJavaScript(String js) {
        if (webView != null) {
            webView.post(() -> webView.evaluateJavascript(js, null));
        }
    }

    /**
     * Called from Go when an event needs to be emitted to JavaScript
     * @param eventName The event name
     * @param eventData The event data (JSON)
     */
    public void emitEvent(String eventName, String eventData) {
        String js = String.format("window.wails && window.wails._emit('%s', %s);",
                escapeJsString(eventName), eventData);
        executeJavaScript(js);
    }

    private String escapeJsString(String str) {
        return str.replace("\\", "\\\\")
                .replace("'", "\\'")
                .replace("\n", "\\n")
                .replace("\r", "\\r");
    }

    // Callback interfaces
    public interface AssetCallback {
        void onAssetReady(byte[] data, String mimeType);
        void onAssetError(String error);
    }

    public interface MessageCallback {
        void onResponse(String response);
        void onError(String error);
    }
}
