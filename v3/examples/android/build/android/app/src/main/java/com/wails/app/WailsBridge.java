package com.wails.app;

import android.content.ClipboardManager;
import android.content.ClipData;
import android.content.Context;
import android.content.Intent;
import android.content.res.Configuration;
import android.net.Uri;
import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.os.VibrationEffect;
import android.os.Vibrator;
import android.os.VibratorManager;
import android.util.DisplayMetrics;
import android.util.Log;
import android.webkit.WebView;
import android.widget.Toast;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
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
    private final Handler mainHandler;
    private WebView webView;
    private Vibrator vibrator;
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
        this.mainHandler = new Handler(Looper.getMainLooper());

        // Initialize vibrator service
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            VibratorManager vibratorManager = (VibratorManager) context.getSystemService(Context.VIBRATOR_MANAGER_SERVICE);
            if (vibratorManager != null) {
                this.vibrator = vibratorManager.getDefaultVibrator();
            }
        } else {
            this.vibrator = (Vibrator) context.getSystemService(Context.VIBRATOR_SERVICE);
        }
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

    // ==================== Android-Specific Features ====================

    /**
     * Trigger haptic feedback (vibration).
     * Called from Go via JNI.
     * @param durationMs Duration of vibration in milliseconds
     */
    @SuppressWarnings("deprecation")
    public void vibrate(int durationMs) {
        Log.d(TAG, "vibrate called: " + durationMs + "ms");

        if (vibrator == null || !vibrator.hasVibrator()) {
            Log.w(TAG, "No vibrator available");
            return;
        }

        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                vibrator.vibrate(VibrationEffect.createOneShot(durationMs, VibrationEffect.DEFAULT_AMPLITUDE));
            } else {
                // Deprecated in API 26, but needed for older devices
                vibrator.vibrate(durationMs);
            }
        } catch (Exception e) {
            Log.e(TAG, "Error triggering vibration", e);
        }
    }

    /**
     * Show a native Android Toast notification.
     * Called from Go via JNI.
     * @param message The message to display
     */
    public void showToast(String message) {
        Log.d(TAG, "showToast called: " + message);

        // Toast must be shown on the main thread
        mainHandler.post(() -> {
            try {
                Toast.makeText(context, message, Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Log.e(TAG, "Error showing toast", e);
            }
        });
    }

    /**
     * Get device information.
     * Called from Go via JNI.
     * @return JSON string with device info
     */
    public String getDeviceInfo() {
        Log.d(TAG, "getDeviceInfo called");

        try {
            JSONObject info = new JSONObject();
            info.put("platform", "android");
            info.put("manufacturer", Build.MANUFACTURER);
            info.put("model", Build.MODEL);
            info.put("brand", Build.BRAND);
            info.put("device", Build.DEVICE);
            info.put("product", Build.PRODUCT);
            info.put("sdkVersion", Build.VERSION.SDK_INT);
            info.put("release", Build.VERSION.RELEASE);

            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                info.put("securityPatch", Build.VERSION.SECURITY_PATCH);
            }

            return info.toString();
        } catch (Exception e) {
            Log.e(TAG, "Error getting device info", e);
            return "{\"platform\":\"android\",\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    /**
     * Open a URL in the default browser using Android's Intent system.
     * Called from Go via JNI.
     * @param url The URL to open
     * @return true if successful, false otherwise
     */
    public boolean openURL(String url) {
        Log.d(TAG, "openURL called: " + url);

        try {
            Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
            intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
            context.startActivity(intent);
            return true;
        } catch (Exception e) {
            Log.e(TAG, "Error opening URL: " + url, e);
            return false;
        }
    }

    /**
     * Set text to clipboard.
     * Called from Go via JNI.
     * @param text The text to copy
     */
    public void setClipboardText(String text) {
        Log.d(TAG, "setClipboardText called");
        mainHandler.post(() -> {
            try {
                ClipboardManager clipboard = (ClipboardManager) context.getSystemService(Context.CLIPBOARD_SERVICE);
                ClipData clip = ClipData.newPlainText("wails", text);
                clipboard.setPrimaryClip(clip);
            } catch (Exception e) {
                Log.e(TAG, "Error setting clipboard text", e);
            }
        });
    }

    /**
     * Get text from clipboard.
     * Called from Go via JNI.
     * @return The clipboard text, or empty string if none
     */
    public String getClipboardText() {
        Log.d(TAG, "getClipboardText called");
        try {
            ClipboardManager clipboard = (ClipboardManager) context.getSystemService(Context.CLIPBOARD_SERVICE);
            if (clipboard != null && clipboard.hasPrimaryClip()) {
                ClipData.Item item = clipboard.getPrimaryClip().getItemAt(0);
                CharSequence text = item.getText();
                return text != null ? text.toString() : "";
            }
        } catch (Exception e) {
            Log.e(TAG, "Error getting clipboard text", e);
        }
        return "";
    }

    /**
     * Set WebView background color.
     * Called from Go via JNI.
     * @param color The ARGB color value
     */
    public void setWebViewBackgroundColor(int color) {
        Log.d(TAG, "setWebViewBackgroundColor called: " + Integer.toHexString(color));
        if (webView != null) {
            mainHandler.post(() -> {
                try {
                    webView.setBackgroundColor(color);
                } catch (Exception e) {
                    Log.e(TAG, "Error setting WebView background color", e);
                }
            });
        }
    }

    /**
     * Check if dark mode is enabled.
     * Called from Go via JNI.
     * @return true if dark mode is enabled
     */
    public boolean isDarkMode() {
        Log.d(TAG, "isDarkMode called");
        try {
            int nightMode = context.getResources().getConfiguration().uiMode & Configuration.UI_MODE_NIGHT_MASK;
            return nightMode == Configuration.UI_MODE_NIGHT_YES;
        } catch (Exception e) {
            Log.e(TAG, "Error checking dark mode", e);
            return false;
        }
    }

    /**
     * Get screen information.
     * Called from Go via JNI.
     * @return JSON string with screen info
     */
    public String getScreenInfo() {
        Log.d(TAG, "getScreenInfo called");
        try {
            DisplayMetrics metrics = context.getResources().getDisplayMetrics();
            JSONObject info = new JSONObject();
            info.put("widthPixels", metrics.widthPixels);
            info.put("heightPixels", metrics.heightPixels);
            info.put("density", metrics.density);
            info.put("densityDpi", metrics.densityDpi);
            info.put("scaledDensity", metrics.scaledDensity);
            info.put("xdpi", metrics.xdpi);
            info.put("ydpi", metrics.ydpi);
            return info.toString();
        } catch (Exception e) {
            Log.e(TAG, "Error getting screen info", e);
            return "{\"widthPixels\":1080,\"heightPixels\":2400}";
        }
    }

    /**
     * Show a message dialog.
     * Called from Go via JNI.
     * @param type Dialog type: "info", "warning", "error", "question"
     * @param title Dialog title
     * @param message Dialog message
     * @param buttons JSON array of button labels, e.g. ["OK"] or ["Yes", "No"]
     * @return The label of the clicked button
     */
    public String showMessageDialog(String type, String title, String message, String buttons) {
        Log.d(TAG, "showMessageDialog: type=" + type + ", title=" + title);

        // Use a blocking approach with CountDownLatch for synchronous result
        final String[] result = new String[1];
        final CountDownLatch latch = new CountDownLatch(1);

        mainHandler.post(() -> {
            try {
                android.app.AlertDialog.Builder builder = new android.app.AlertDialog.Builder(context);
                builder.setTitle(title);
                builder.setMessage(message);

                // Parse buttons JSON
                JSONArray btnArray = new JSONArray(buttons);

                if (btnArray.length() >= 1) {
                    builder.setPositiveButton(btnArray.getString(0), (dialog, which) -> {
                        result[0] = btnArray.optString(0, "OK");
                        latch.countDown();
                    });
                }
                if (btnArray.length() >= 2) {
                    builder.setNegativeButton(btnArray.getString(1), (dialog, which) -> {
                        result[0] = btnArray.optString(1, "Cancel");
                        latch.countDown();
                    });
                }
                if (btnArray.length() >= 3) {
                    builder.setNeutralButton(btnArray.getString(2), (dialog, which) -> {
                        result[0] = btnArray.optString(2, "");
                        latch.countDown();
                    });
                }

                builder.setOnCancelListener(dialog -> {
                    result[0] = "";
                    latch.countDown();
                });

                builder.show();
            } catch (Exception e) {
                Log.e(TAG, "Error showing dialog", e);
                result[0] = "";
                latch.countDown();
            }
        });

        try {
            latch.await(30, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            Log.e(TAG, "Dialog interrupted", e);
            return "";
        }

        return result[0] != null ? result[0] : "";
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
