package com.wails.app;

import android.app.Activity;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.Intent;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageInfo;
import android.content.res.Configuration;
import android.graphics.Rect;
import android.hardware.camera2.CameraCharacteristics;
import android.hardware.camera2.CameraManager;
import android.net.Uri;
import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.os.VibrationEffect;
import android.os.Vibrator;
import android.provider.Settings;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.view.WindowInsets;
import android.view.WindowInsetsController;
import android.view.WindowManager;
import android.webkit.WebView;
import android.widget.Toast;

import androidx.appcompat.app.AlertDialog;

import org.json.JSONArray;
import org.json.JSONObject;

/**
 * WailsBridge manages the connection between the Java/Android side and the Go
 * native library. It handles:
 * - Loading and initializing the native Go library
 * - Serving asset requests from Go
 * - Passing messages between JavaScript and Go
 * - Native facilities the Go side calls via JNI (dialogs, clipboard,
 *   screen/device info, toasts, vibration, main-thread dispatch)
 */
public class WailsBridge {
    private static final String TAG = "WailsBridge";
    private static final boolean DEBUG = BuildConfig.DEBUG;

    static {
        // Load the native Go library
        System.loadLibrary("wails");
    }

    private final Activity activity;
    private final Handler mainHandler = new Handler(Looper.getMainLooper());
    private WebView webView;
    private volatile boolean initialized = false;

    // Native methods - implemented in Go
    private static native void nativeInit(WailsBridge bridge);
    private static native void nativeShutdown();
    private static native void nativeOnStart();
    private static native void nativeOnResume();
    private static native void nativeOnPause();
    private static native void nativeOnStop();
    private static native void nativeOnLowMemory();
    private static native void nativeOnPageFinished(String url);
    private static native byte[] nativeServeAsset(String path, String method, String headers);
    private static native String nativeHandleMessage(String message);
    private static native String nativeHandleRuntimeCall(String payload);
    private static native String nativeGetAssetMimeType(String path);
    private static native void nativeDialogCallback(int callbackID, int buttonIndex);
    private static native void nativeFilePickerResult(int callbackID, String path);
    private static native void nativeFilePickerDone(int callbackID);
    private static native void nativeMainThreadCallback(int callbackID);
    private static native void nativeEmitSystemEvent(String name, String json);
    private static native void nativeEmitEvent(String name, String json);

    public WailsBridge(Activity activity) {
        this.activity = activity;
    }

    /**
     * Initialize the native Go library
     */
    public void initialize() {
        if (initialized) {
            return;
        }
        try {
            nativeInit(this);
            initialized = true;
            Log.i(TAG, "Wails bridge initialized");
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
        try {
            nativeShutdown();
            initialized = false;
        } catch (Exception e) {
            Log.e(TAG, "Error during shutdown", e);
        }
    }

    /**
     * The WebView used for JavaScript execution. Must be set before the page
     * loads.
     */
    public void setWebView(WebView webView) {
        this.webView = webView;
    }

    // Lifecycle forwarding

    public void onStart() {
        if (initialized) nativeOnStart();
    }

    public void onResume() {
        if (initialized) nativeOnResume();
    }

    public void onPause() {
        if (initialized) nativeOnPause();
    }

    public void onStop() {
        if (initialized) nativeOnStop();
    }

    public void onLowMemory() {
        if (initialized) nativeOnLowMemory();
    }

    /**
     * Notify Go that the page finished loading.
     */
    public void onPageFinished(String url) {
        if (initialized) nativeOnPageFinished(url);
    }

    /**
     * Emit a "system:*" event (battery, network, lock, theme, lifecycle) to JS.
     * Called from the system-event receivers registered by MainActivity.
     */
    public void emitSystemEvent(String name, String json) {
        if (initialized) nativeEmitSystemEvent(name, json);
    }

    /**
     * Emit an arbitrary custom event with a JSON payload to JS. Used by the
     * mobile-feature bridges to deliver asynchronous results.
     */
    public void emitEvent(String name, String json) {
        if (initialized) nativeEmitEvent(name, json);
    }

    /**
     * Serve an asset from the Go asset server
     */
    public byte[] serveAsset(String path, String method, String headers) {
        if (!initialized) {
            Log.w(TAG, "Bridge not initialized, cannot serve asset: " + path);
            return null;
        }
        if (DEBUG) Log.d(TAG, "Serving asset: " + path);
        try {
            return nativeServeAsset(path, method, headers);
        } catch (Exception e) {
            Log.e(TAG, "Error serving asset: " + path, e);
            return null;
        }
    }

    /**
     * Get the MIME type for an asset
     */
    public String getAssetMimeType(String path) {
        if (!initialized) {
            return "application/octet-stream";
        }
        try {
            String mimeType = nativeGetAssetMimeType(path);
            return mimeType != null ? mimeType : "application/octet-stream";
        } catch (Exception e) {
            return "application/octet-stream";
        }
    }

    /**
     * Handle a message from JavaScript
     */
    public String handleMessage(String message) {
        if (!initialized) {
            Log.w(TAG, "Bridge not initialized, cannot handle message");
            return "{\"error\":\"Bridge not initialized\"}";
        }
        if (DEBUG) Log.d(TAG, "Message from JS: " + message);
        try {
            return nativeHandleMessage(message);
        } catch (Exception e) {
            Log.e(TAG, "Error handling message", e);
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    /**
     * Handle a runtime call from JavaScript (the Android transport).
     * The payload and response are JSON strings.
     */
    public String handleRuntimeCall(String payload) {
        if (!initialized) {
            return "{\"ok\":false,\"error\":\"Bridge not initialized\"}";
        }
        if (DEBUG) Log.d(TAG, "Runtime call: " + payload);
        try {
            return nativeHandleRuntimeCall(payload);
        } catch (Exception e) {
            Log.e(TAG, "Error in runtime call", e);
            return "{\"ok\":false,\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    /**
     * Execute JavaScript in the WebView. Called from Go via JNI (any thread).
     */
    public void executeJavaScript(final String js) {
        final WebView view = webView;
        if (view == null) {
            Log.w(TAG, "executeJavaScript: no WebView attached");
            return;
        }
        mainHandler.post(() -> view.evaluateJavascript(js, null));
    }

    // Facilities called from Go via JNI

    /**
     * Screen metrics as JSON: hardware pixels, density and system bar insets.
     */
    public String getScreenInfoJson() {
        try {
            JSONObject result = new JSONObject();

            int insetTop = 0, insetBottom = 0, insetLeft = 0, insetRight = 0;
            int widthPx, heightPx;
            float density;

            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
                android.view.WindowMetrics wm = activity.getWindowManager().getCurrentWindowMetrics();
                Rect bounds = wm.getBounds();
                widthPx = bounds.width();
                heightPx = bounds.height();
                density = activity.getResources().getDisplayMetrics().density;
                android.graphics.Insets insets = wm.getWindowInsets()
                        .getInsetsIgnoringVisibility(WindowInsets.Type.systemBars());
                insetTop = insets.top;
                insetBottom = insets.bottom;
                insetLeft = insets.left;
                insetRight = insets.right;
            } else {
                DisplayMetrics metrics = new DisplayMetrics();
                activity.getWindowManager().getDefaultDisplay().getRealMetrics(metrics);
                widthPx = metrics.widthPixels;
                heightPx = metrics.heightPixels;
                density = metrics.density;
            }

            result.put("widthPx", widthPx);
            result.put("heightPx", heightPx);
            result.put("density", density);
            result.put("insetTop", insetTop);
            result.put("insetBottom", insetBottom);
            result.put("insetLeft", insetLeft);
            result.put("insetRight", insetRight);
            return result.toString();
        } catch (Exception e) {
            Log.e(TAG, "getScreenInfoJson failed", e);
            return "";
        }
    }

    /**
     * Device information as JSON.
     */
    public String getDeviceInfoJson() {
        try {
            JSONObject result = new JSONObject();
            result.put("manufacturer", Build.MANUFACTURER);
            result.put("brand", Build.BRAND);
            result.put("model", Build.MODEL);
            result.put("device", Build.DEVICE);
            result.put("version", Build.VERSION.RELEASE);
            result.put("sdkInt", Build.VERSION.SDK_INT);
            return result.toString();
        } catch (Exception e) {
            return "";
        }
    }

    public boolean isDarkMode() {
        int mode = activity.getResources().getConfiguration().uiMode
                & Configuration.UI_MODE_NIGHT_MASK;
        return mode == Configuration.UI_MODE_NIGHT_YES;
    }

    public boolean isMainThread() {
        return Looper.myLooper() == Looper.getMainLooper();
    }

    /**
     * Post a Go callback onto the Android main thread.
     */
    public void runOnMainThread(final int callbackID) {
        // Guard against the callback firing after shutdown() tore down the Go side
        mainHandler.post(() -> {
            if (initialized) nativeMainThreadCallback(callbackID);
        });
    }

    // Clipboard (note: reads on Android 10+ require input focus)

    public void setClipboardText(String text) {
        try {
            ClipboardManager cm = (ClipboardManager) activity.getSystemService(Context.CLIPBOARD_SERVICE);
            cm.setPrimaryClip(ClipData.newPlainText("wails", text));
        } catch (Exception e) {
            Log.e(TAG, "setClipboardText failed", e);
        }
    }

    public String getClipboardText() {
        try {
            ClipboardManager cm = (ClipboardManager) activity.getSystemService(Context.CLIPBOARD_SERVICE);
            ClipData clip = cm.getPrimaryClip();
            if (clip != null && clip.getItemCount() > 0) {
                CharSequence text = clip.getItemAt(0).coerceToText(activity);
                return text != null ? text.toString() : "";
            }
        } catch (Exception e) {
            Log.e(TAG, "getClipboardText failed", e);
        }
        return "";
    }

    public void showToast(final String message) {
        mainHandler.post(() -> Toast.makeText(activity, message, Toast.LENGTH_SHORT).show());
    }

    public void vibrate(int durationMs) {
        try {
            Vibrator vibrator = (Vibrator) activity.getSystemService(Context.VIBRATOR_SERVICE);
            if (vibrator == null || !vibrator.hasVibrator()) {
                return;
            }
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                vibrator.vibrate(VibrationEffect.createOneShot(durationMs, VibrationEffect.DEFAULT_AMPLITUDE));
            } else {
                vibrator.vibrate(durationMs);
            }
        } catch (Exception e) {
            Log.e(TAG, "vibrate failed", e);
        }
    }

    // MARK: - Mobile features (Phase A)

    /**
     * Present the Android share chooser. json: {"text": "...", "url": "..."}.
     */
    public void share(final String json) {
        mainHandler.post(() -> {
            try {
                JSONObject opts = new JSONObject(json);
                String text = opts.optString("text", "");
                String url = opts.optString("url", "");
                StringBuilder body = new StringBuilder();
                if (!text.isEmpty()) body.append(text);
                if (!url.isEmpty()) {
                    if (body.length() > 0) body.append("\n");
                    body.append(url);
                }
                if (body.length() == 0) return;
                Intent send = new Intent(Intent.ACTION_SEND);
                send.setType("text/plain");
                send.putExtra(Intent.EXTRA_TEXT, body.toString());
                Intent chooser = Intent.createChooser(send, null);
                chooser.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                activity.startActivity(chooser);
            } catch (Exception e) {
                Log.e(TAG, "share failed", e);
            }
        });
    }

    /**
     * Open a URL in the system browser.
     */
    public void openURL(final String url) {
        mainHandler.post(() -> {
            try {
                Intent view = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
                view.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                activity.startActivity(view);
            } catch (Exception e) {
                Log.e(TAG, "openURL failed", e);
            }
        });
    }

    /**
     * Keep the screen on (1) or release the hold (0) via FLAG_KEEP_SCREEN_ON.
     */
    public void setKeepAwake(final int enabled) {
        mainHandler.post(() -> {
            if (enabled != 0) {
                activity.getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
            } else {
                activity.getWindow().clearFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
            }
        });
    }

    /**
     * Toggle the camera flash (torch). Emits "native:torch" with the resulting
     * state and availability.
     */
    public void setTorch(final int enabled) {
        mainHandler.post(() -> {
            try {
                CameraManager cm = (CameraManager) activity.getSystemService(Context.CAMERA_SERVICE);
                String flashId = null;
                for (String id : cm.getCameraIdList()) {
                    Boolean hasFlash = cm.getCameraCharacteristics(id)
                            .get(CameraCharacteristics.FLASH_INFO_AVAILABLE);
                    if (Boolean.TRUE.equals(hasFlash)) {
                        flashId = id;
                        break;
                    }
                }
                if (flashId == null) {
                    emitEvent("native:torch", "{\"on\":false,\"available\":false}");
                    return;
                }
                cm.setTorchMode(flashId, enabled != 0);
                emitEvent("native:torch",
                        enabled != 0 ? "{\"on\":true,\"available\":true}"
                                     : "{\"on\":false,\"available\":true}");
            } catch (Exception e) {
                Log.e(TAG, "setTorch failed", e);
                emitEvent("native:torch", "{\"on\":false,\"available\":false}");
            }
        });
    }

    // MARK: - Mobile features (Phase B)

    /**
     * System-bar insets as JSON {"top","bottom","left","right"} in px.
     */
    public String getSafeAreaJson() {
        try {
            int top = 0, bottom = 0, left = 0, right = 0;
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
                android.graphics.Insets insets = activity.getWindowManager()
                        .getCurrentWindowMetrics().getWindowInsets()
                        .getInsetsIgnoringVisibility(WindowInsets.Type.systemBars());
                top = insets.top; bottom = insets.bottom; left = insets.left; right = insets.right;
            }
            return new JSONObject()
                    .put("top", top).put("bottom", bottom)
                    .put("left", left).put("right", right).toString();
        } catch (Exception e) {
            return "{\"top\":0,\"bottom\":0,\"left\":0,\"right\":0}";
        }
    }

    /**
     * Set window brightness, 0-100. A negative value restores the system default.
     */
    public void setBrightness(final int pct) {
        mainHandler.post(() -> {
            try {
                WindowManager.LayoutParams lp = activity.getWindow().getAttributes();
                lp.screenBrightness = pct < 0 ? WindowManager.LayoutParams.BRIGHTNESS_OVERRIDE_NONE
                                              : Math.max(0.01f, Math.min(1f, pct / 100f));
                activity.getWindow().setAttributes(lp);
            } catch (Exception e) {
                Log.e(TAG, "setBrightness failed", e);
            }
        });
    }

    /**
     * Current brightness as {"value": 0.0-1.0}. Falls back to the system
     * brightness setting when the window has not overridden it.
     */
    public String getBrightnessJson() {
        try {
            float v = activity.getWindow().getAttributes().screenBrightness;
            if (v < 0) {
                int sys = Settings.System.getInt(activity.getContentResolver(),
                        Settings.System.SCREEN_BRIGHTNESS, 128);
                v = sys / 255f;
            }
            return new JSONObject().put("value", v).toString();
        } catch (Exception e) {
            return "{\"value\":-1}";
        }
    }

    /**
     * App info as JSON {"name","version","build","bundleId"}.
     */
    public String getAppInfoJson() {
        try {
            PackageInfo pi = activity.getPackageManager()
                    .getPackageInfo(activity.getPackageName(), 0);
            long code = Build.VERSION.SDK_INT >= Build.VERSION_CODES.P
                    ? pi.getLongVersionCode()
                    : pi.versionCode;
            CharSequence label = activity.getApplicationInfo()
                    .loadLabel(activity.getPackageManager());
            return new JSONObject()
                    .put("name", label != null ? label.toString() : "")
                    .put("version", pi.versionName != null ? pi.versionName : "")
                    .put("build", String.valueOf(code))
                    .put("bundleId", activity.getPackageName())
                    .toString();
        } catch (Exception e) {
            return "{}";
        }
    }

    /**
     * Lock orientation to "portrait", "landscape" or "auto".
     */
    public void setOrientation(final String mode) {
        mainHandler.post(() -> {
            int o = ActivityInfo.SCREEN_ORIENTATION_UNSPECIFIED;
            if ("portrait".equals(mode)) o = ActivityInfo.SCREEN_ORIENTATION_PORTRAIT;
            else if ("landscape".equals(mode)) o = ActivityInfo.SCREEN_ORIENTATION_LANDSCAPE;
            activity.setRequestedOrientation(o);
        });
    }

    /**
     * Current orientation as {"orientation":"portrait"|"landscape"}.
     */
    public String getOrientationJson() {
        int o = activity.getResources().getConfiguration().orientation;
        String s = o == Configuration.ORIENTATION_LANDSCAPE ? "landscape" : "portrait";
        try {
            return new JSONObject().put("orientation", s).toString();
        } catch (Exception e) {
            return "{\"orientation\":\"" + s + "\"}";
        }
    }

    /**
     * Set status-bar appearance. json: {"style":"light|dark|default","hidden":bool}.
     * "light" = light (white) icons; "dark" = dark icons.
     */
    public void setStatusBar(final String json) {
        mainHandler.post(() -> {
            try {
                JSONObject opts = new JSONObject(json);
                String style = opts.optString("style", "default");
                boolean hidden = opts.optBoolean("hidden", false);
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
                    WindowInsetsController c = activity.getWindow().getInsetsController();
                    if (c != null) {
                        if ("dark".equals(style)) {
                            c.setSystemBarsAppearance(
                                    WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS,
                                    WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS);
                        } else if ("light".equals(style)) {
                            c.setSystemBarsAppearance(0,
                                    WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS);
                        }
                        if (hidden) c.hide(WindowInsets.Type.statusBars());
                        else c.show(WindowInsets.Type.statusBars());
                    }
                } else {
                    int vis = activity.getWindow().getDecorView().getSystemUiVisibility();
                    if ("dark".equals(style)) vis |= View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                    else if ("light".equals(style)) vis &= ~View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                    if (hidden) vis |= View.SYSTEM_UI_FLAG_FULLSCREEN;
                    else vis &= ~View.SYSTEM_UI_FLAG_FULLSCREEN;
                    activity.getWindow().getDecorView().setSystemUiVisibility(vis);
                }
            } catch (Exception e) {
                Log.e(TAG, "setStatusBar failed", e);
            }
        });
    }

    /**
     * Show a message dialog. optionsJson:
     * {"title": "...", "message": "...",
     *  "buttons": [{"label": "...", "isCancel": bool, "isDefault": bool}]}
     * Calls nativeDialogCallback with the index of the pressed button in the
     * original buttons array, or -1 when dismissed without a matching button.
     */
    public void showMessageDialog(final int callbackID, final String optionsJson) {
        mainHandler.post(() -> {
            try {
                JSONObject options = new JSONObject(optionsJson);
                String title = options.optString("title", "");
                String message = options.optString("message", "");
                JSONArray buttons = options.optJSONArray("buttons");

                AlertDialog.Builder builder = new AlertDialog.Builder(activity);
                builder.setTitle(title);
                builder.setMessage(message);

                int count = buttons != null ? buttons.length() : 0;
                int cancelIndex = -1;
                for (int i = 0; i < count; i++) {
                    if (buttons.getJSONObject(i).optBoolean("isCancel", false)) {
                        cancelIndex = i;
                        break;
                    }
                }
                final int dismissIndex = cancelIndex;
                builder.setOnCancelListener(d -> dialogCallback(callbackID, dismissIndex));

                if (count == 0) {
                    builder.setPositiveButton(android.R.string.ok,
                            (d, w) -> dialogCallback(callbackID, -1));
                } else if (count <= 3) {
                    // Map buttons to AlertDialog slots: the default (or last)
                    // button is positive, the cancel button negative, any
                    // remaining button neutral.
                    int positive = -1;
                    for (int i = 0; i < count; i++) {
                        if (buttons.getJSONObject(i).optBoolean("isDefault", false)) {
                            positive = i;
                            break;
                        }
                    }
                    if (positive == -1) {
                        positive = count - 1;
                        if (positive == cancelIndex && count > 1) {
                            positive = count - 2;
                        }
                    }
                    int negative = cancelIndex;
                    if (negative == -1) {
                        for (int i = count - 1; i >= 0; i--) {
                            if (i != positive) {
                                negative = i;
                                break;
                            }
                        }
                    }
                    int neutral = -1;
                    for (int i = 0; i < count; i++) {
                        if (i != positive && i != negative) {
                            neutral = i;
                            break;
                        }
                    }

                    final int positiveIdx = positive, negativeIdx = negative, neutralIdx = neutral;
                    builder.setPositiveButton(buttons.getJSONObject(positive).optString("label", "OK"),
                            (d, w) -> dialogCallback(callbackID, positiveIdx));
                    if (negative != -1) {
                        builder.setNegativeButton(buttons.getJSONObject(negative).optString("label", "Cancel"),
                                (d, w) -> dialogCallback(callbackID, negativeIdx));
                    }
                    if (neutral != -1) {
                        builder.setNeutralButton(buttons.getJSONObject(neutral).optString("label", ""),
                                (d, w) -> dialogCallback(callbackID, neutralIdx));
                    }
                } else {
                    // More than three buttons: show as a list
                    String[] labels = new String[count];
                    for (int i = 0; i < count; i++) {
                        labels[i] = buttons.getJSONObject(i).optString("label", "");
                    }
                    builder.setItems(labels, (d, which) -> dialogCallback(callbackID, which));
                }

                builder.show();
            } catch (Exception e) {
                Log.e(TAG, "showMessageDialog failed", e);
                dialogCallback(callbackID, -1);
            }
        });
    }

    /**
     * Show the system document picker. optionsJson: {"multiple": bool}.
     * Results flow back through filePickerResult/filePickerDone.
     */
    public void showFilePicker(final int callbackID, final String optionsJson) {
        boolean multiple = false;
        try {
            multiple = new JSONObject(optionsJson).optBoolean("multiple", false);
        } catch (Exception ignored) {
        }
        final boolean allowMultiple = multiple;
        mainHandler.post(() -> {
            if (activity instanceof MainActivity) {
                ((MainActivity) activity).launchFilePicker(callbackID, allowMultiple);
            } else {
                Log.e(TAG, "showFilePicker: activity is not a MainActivity");
                nativeFilePickerDone(callbackID);
            }
        });
    }

    /** Forward a picked file path to Go (package-private, used by MainActivity). */
    void filePickerResult(int callbackID, String path) {
        // The picker completes on a background thread that may outlive shutdown()
        if (initialized) nativeFilePickerResult(callbackID, path);
    }

    /** Signal the end of a file picking session (package-private). */
    void filePickerDone(int callbackID) {
        if (initialized) nativeFilePickerDone(callbackID);
    }

    /** Dialog button callback that no-ops once the Go side is gone. */
    private void dialogCallback(int callbackID, int buttonIndex) {
        if (initialized) nativeDialogCallback(callbackID, buttonIndex);
    }
}
