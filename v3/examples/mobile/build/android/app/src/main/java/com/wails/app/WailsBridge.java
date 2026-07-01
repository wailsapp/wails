package com.wails.app;

import android.app.Activity;
import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.content.IntentFilter;
import android.content.res.Configuration;
import android.graphics.Rect;
import android.hardware.Sensor;
import android.hardware.SensorEvent;
import android.hardware.SensorEventListener;
import android.hardware.SensorManager;
import android.hardware.camera2.CameraCharacteristics;
import android.hardware.camera2.CameraManager;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.net.ConnectivityManager;
import android.net.Network;
import android.net.NetworkCapabilities;
import android.net.Uri;
import android.os.BatteryManager;
import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.os.PowerManager;
import android.os.StatFs;
import android.os.VibrationEffect;
import android.os.Vibrator;
import android.provider.Settings;
import android.speech.tts.TextToSpeech;
import android.util.DisplayMetrics;
import android.util.Log;
import android.view.View;
import android.view.WindowInsets;
import android.view.WindowInsetsController;
import android.view.WindowManager;
import android.webkit.WebView;
import android.widget.Toast;

import androidx.appcompat.app.AlertDialog;
import androidx.biometric.BiometricManager;
import androidx.biometric.BiometricPrompt;
import androidx.core.app.NotificationCompat;
import androidx.core.content.ContextCompat;
import androidx.fragment.app.FragmentActivity;
import androidx.security.crypto.EncryptedSharedPreferences;
import androidx.security.crypto.MasterKey;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.Locale;
import java.util.concurrent.Executor;

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

    // Phase D state: sensor listeners, speech engine and keyboard watcher are
    // retained so they can be registered and torn down on demand.
    private SensorEventListener accelListener;
    private SensorEventListener proximityListener;
    private long lastMotionEmit = 0;
    private TextToSpeech tts;
    private View.OnApplyWindowInsetsListener keyboardListener;
    // Battery: remember the user's intent so sensors paused while the app is
    // backgrounded can be restored on foreground; the torch is switched off.
    private boolean motionWanted = false;
    private boolean proximityWanted = false;
    private boolean torchOn = false;

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
        resumeFeaturesForForeground();
    }

    public void onResume() {
        if (initialized) nativeOnResume();
    }

    public void onPause() {
        if (initialized) nativeOnPause();
    }

    public void onStop() {
        if (initialized) nativeOnStop();
        pauseFeaturesForBackground();
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
     * Toggle the camera flash (torch). Emits "common:torch" with the resulting
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
                    emitEvent("common:torch", "{\"on\":false,\"available\":false}");
                    return;
                }
                cm.setTorchMode(flashId, enabled != 0);
                torchOn = enabled != 0;
                emitEvent("common:torch",
                        enabled != 0 ? "{\"on\":true,\"available\":true}"
                                     : "{\"on\":false,\"available\":true}");
            } catch (Exception e) {
                Log.e(TAG, "setTorch failed", e);
                emitEvent("common:torch", "{\"on\":false,\"available\":false}");
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

    // MARK: - Mobile features (Phase C)

    private void emitBiometric(boolean ok, String error) {
        try {
            JSONObject o = new JSONObject().put("ok", ok);
            if (error != null) o.put("error", error);
            emitEvent("common:biometric", o.toString());
        } catch (Exception ignored) {
        }
    }

    /**
     * Show the BiometricPrompt (biometric or device credential). The outcome is
     * emitted as "common:biometric" {ok, error}.
     */
    public void authenticate(final String reason) {
        mainHandler.post(() -> {
            try {
                int allowed = BiometricManager.Authenticators.BIOMETRIC_WEAK
                        | BiometricManager.Authenticators.DEVICE_CREDENTIAL;
                BiometricManager bm = BiometricManager.from(activity);
                if (bm.canAuthenticate(allowed) != BiometricManager.BIOMETRIC_SUCCESS) {
                    emitBiometric(false, "no biometrics or device credential enrolled");
                    return;
                }
                Executor exec = ContextCompat.getMainExecutor(activity);
                BiometricPrompt prompt = new BiometricPrompt((FragmentActivity) activity, exec,
                        new BiometricPrompt.AuthenticationCallback() {
                            @Override
                            public void onAuthenticationSucceeded(BiometricPrompt.AuthenticationResult result) {
                                emitBiometric(true, null);
                            }
                            @Override
                            public void onAuthenticationError(int code, CharSequence err) {
                                emitBiometric(false, err != null ? err.toString() : "error " + code);
                            }
                            // onAuthenticationFailed = a single non-match; prompt stays up, no terminal event.
                        });
                BiometricPrompt.PromptInfo info = new BiometricPrompt.PromptInfo.Builder()
                        .setTitle("Authenticate")
                        .setSubtitle(reason != null && !reason.isEmpty() ? reason : "Confirm it's you")
                        .setAllowedAuthenticators(allowed)
                        .build();
                prompt.authenticate(info);
            } catch (Exception e) {
                Log.e(TAG, "authenticate failed", e);
                emitBiometric(false, "exception");
            }
        });
    }

    /**
     * Post a local notification. json: {"title","body"}. Requests the
     * POST_NOTIFICATIONS runtime permission on Android 13+.
     */
    public void postNotification(final String json) {
        mainHandler.post(() -> {
            try {
                JSONObject opts = new JSONObject(json);
                String title = opts.optString("title", "Notification");
                String body = opts.optString("body", "");
                String channelId = "wails_default";
                NotificationManager nm =
                        (NotificationManager) activity.getSystemService(Context.NOTIFICATION_SERVICE);
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                    NotificationChannel ch = new NotificationChannel(
                            channelId, "General", NotificationManager.IMPORTANCE_DEFAULT);
                    nm.createNotificationChannel(ch);
                }
                if (Build.VERSION.SDK_INT >= 33 && activity.checkSelfPermission(
                        "android.permission.POST_NOTIFICATIONS") != PackageManager.PERMISSION_GRANTED) {
                    activity.requestPermissions(
                            new String[]{"android.permission.POST_NOTIFICATIONS"}, 1001);
                }
                Notification n = new NotificationCompat.Builder(activity, channelId)
                        .setSmallIcon(android.R.drawable.ic_dialog_info)
                        .setContentTitle(title)
                        .setContentText(body)
                        .setAutoCancel(true)
                        .build();
                nm.notify((int) (System.currentTimeMillis() & 0x0fffffff), n);
                emitEvent("common:notification", "{\"ok\":true}");
            } catch (Exception e) {
                Log.e(TAG, "postNotification failed", e);
                emitEvent("common:notification", "{\"ok\":false}");
            }
        });
    }

    /**
     * Backing store for secure storage. Uses EncryptedSharedPreferences (AES via
     * the Android Keystore) on API 23+, falling back to plain prefs below that.
     */
    private SharedPreferences securePrefs() {
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                MasterKey key = new MasterKey.Builder(activity)
                        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
                        .build();
                return EncryptedSharedPreferences.create(activity, "wails_secure", key,
                        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
                        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM);
            }
        } catch (Exception e) {
            Log.e(TAG, "securePrefs failed, using plain prefs", e);
        }
        return activity.getSharedPreferences("wails_secure_plain", Context.MODE_PRIVATE);
    }

    /** Store a value in secure storage. json: {"key","value"}. */
    public void secureSet(final String json) {
        try {
            JSONObject o = new JSONObject(json);
            securePrefs().edit().putString(o.optString("key"), o.optString("value")).apply();
        } catch (Exception e) {
            Log.e(TAG, "secureSet failed", e);
        }
    }

    /** Read a value from secure storage (empty if absent). */
    public String secureGet(final String key) {
        try {
            return securePrefs().getString(key, "");
        } catch (Exception e) {
            return "";
        }
    }

    /** Remove a value from secure storage. */
    public void secureDelete(final String key) {
        try {
            securePrefs().edit().remove(key).apply();
        } catch (Exception e) {
            Log.e(TAG, "secureDelete failed", e);
        }
    }

    // MARK: - Mobile features (Phase D: sensors & hardware)

    /**
     * Play a haptic pattern via the Vibrator. type: impact-light|impact-medium|
     * impact-heavy|success|warning|error|selection.
     */
    public void haptic(final String type) {
        try {
            Vibrator vibrator = (Vibrator) activity.getSystemService(Context.VIBRATOR_SERVICE);
            if (vibrator == null || !vibrator.hasVibrator()) return;
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
                int effect;
                switch (type) {
                    case "impact-heavy": case "error":
                        effect = VibrationEffect.EFFECT_HEAVY_CLICK; break;
                    case "impact-light": case "selection":
                        effect = VibrationEffect.EFFECT_TICK; break;
                    case "success": case "warning": case "impact-medium": default:
                        effect = VibrationEffect.EFFECT_CLICK; break;
                }
                vibrator.vibrate(VibrationEffect.createPredefined(effect));
            } else if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                long ms = "impact-heavy".equals(type) || "error".equals(type) ? 40
                        : "impact-light".equals(type) || "selection".equals(type) ? 10 : 20;
                vibrator.vibrate(VibrationEffect.createOneShot(ms, VibrationEffect.DEFAULT_AMPLITUDE));
            } else {
                vibrator.vibrate(20);
            }
        } catch (Exception e) {
            Log.e(TAG, "haptic failed", e);
        }
    }

    private void emitLocationError(String msg) {
        try {
            emitEvent("common:location", new JSONObject().put("error", msg).toString());
        } catch (Exception ignored) {
        }
    }

    /**
     * Request a one-shot location fix. Emits "common:location"
     * {lat,lng,accuracy} or {error}. Requests ACCESS_FINE_LOCATION on first use.
     */
    public void getLocation() {
        mainHandler.post(() -> {
            try {
                if (activity.checkSelfPermission("android.permission.ACCESS_FINE_LOCATION")
                        != PackageManager.PERMISSION_GRANTED) {
                    activity.requestPermissions(new String[]{
                            "android.permission.ACCESS_FINE_LOCATION",
                            "android.permission.ACCESS_COARSE_LOCATION"}, 1002);
                    emitLocationError("location permission requested — tap again once granted");
                    return;
                }
                LocationManager lm = (LocationManager) activity.getSystemService(Context.LOCATION_SERVICE);
                if (lm == null) { emitLocationError("location unavailable"); return; }
                Location best = null;
                for (String provider : new String[]{LocationManager.GPS_PROVIDER,
                        LocationManager.NETWORK_PROVIDER, LocationManager.PASSIVE_PROVIDER}) {
                    try {
                        Location l = lm.getLastKnownLocation(provider);
                        if (l != null && (best == null || l.getTime() > best.getTime())) best = l;
                    } catch (SecurityException | IllegalArgumentException ignored) {
                    }
                }
                if (best != null) {
                    emitLocation(best);
                    return;
                }
                // No cached fix: request a single update from whichever provider is enabled.
                String provider = lm.isProviderEnabled(LocationManager.GPS_PROVIDER)
                        ? LocationManager.GPS_PROVIDER
                        : lm.isProviderEnabled(LocationManager.NETWORK_PROVIDER)
                        ? LocationManager.NETWORK_PROVIDER : null;
                if (provider == null) { emitLocationError("no location provider enabled"); return; }
                lm.requestSingleUpdate(provider, new LocationListener() {
                    @Override public void onLocationChanged(Location location) { emitLocation(location); }
                    @Override public void onProviderEnabled(String p) {}
                    @Override public void onProviderDisabled(String p) {}
                    @Override public void onStatusChanged(String p, int s, android.os.Bundle e) {}
                }, Looper.getMainLooper());
            } catch (Exception e) {
                Log.e(TAG, "getLocation failed", e);
                emitLocationError("exception");
            }
        });
    }

    private void emitLocation(Location l) {
        try {
            emitEvent("common:location", new JSONObject()
                    .put("lat", l.getLatitude())
                    .put("lng", l.getLongitude())
                    .put("accuracy", l.getAccuracy()).toString());
        } catch (Exception ignored) {
        }
    }

    /** Start (1) / stop (0) accelerometer updates, streamed as "common:motion". */
    public void setMotion(final int enabled) {
        motionWanted = enabled != 0;
        mainHandler.post(() -> {
            SensorManager sm = (SensorManager) activity.getSystemService(Context.SENSOR_SERVICE);
            if (sm == null) return;
            if (enabled != 0) {
                if (accelListener != null) return;
                Sensor accel = sm.getDefaultSensor(Sensor.TYPE_ACCELEROMETER);
                if (accel == null) { emitEvent("common:motion", "{\"available\":false}"); return; }
                accelListener = new SensorEventListener() {
                    @Override public void onSensorChanged(SensorEvent e) {
                        long now = System.currentTimeMillis();
                        if (now - lastMotionEmit < 100) return; // throttle to ~10 Hz
                        lastMotionEmit = now;
                        try {
                            emitEvent("common:motion", new JSONObject()
                                    .put("x", e.values[0]).put("y", e.values[1])
                                    .put("z", e.values[2]).toString());
                        } catch (Exception ignored) {
                        }
                    }
                    @Override public void onAccuracyChanged(Sensor s, int a) {}
                };
                sm.registerListener(accelListener, accel, SensorManager.SENSOR_DELAY_UI);
            } else if (accelListener != null) {
                sm.unregisterListener(accelListener);
                accelListener = null;
            }
        });
    }

    /** Enable (1) / disable (0) the proximity sensor, reported as "common:proximity". */
    public void setProximity(final int enabled) {
        proximityWanted = enabled != 0;
        mainHandler.post(() -> {
            SensorManager sm = (SensorManager) activity.getSystemService(Context.SENSOR_SERVICE);
            if (sm == null) return;
            if (enabled != 0) {
                if (proximityListener != null) return;
                Sensor prox = sm.getDefaultSensor(Sensor.TYPE_PROXIMITY);
                if (prox == null) { emitEvent("common:proximity", "{\"available\":false}"); return; }
                final float max = prox.getMaximumRange();
                proximityListener = new SensorEventListener() {
                    @Override public void onSensorChanged(SensorEvent e) {
                        boolean near = e.values[0] < max;
                        emitEvent("common:proximity", "{\"near\":" + (near ? "true" : "false") + "}");
                    }
                    @Override public void onAccuracyChanged(Sensor s, int a) {}
                };
                sm.registerListener(proximityListener, prox, SensorManager.SENSOR_DELAY_NORMAL);
            } else if (proximityListener != null) {
                sm.unregisterListener(proximityListener);
                proximityListener = null;
            }
        });
    }

    /**
     * Called when the activity leaves the foreground (onStop): stop the
     * accelerometer and proximity sensor and switch the torch off so none of
     * them drain the battery while the app isn't visible. The "wanted" flags are
     * kept so foregrounding can restore the sensors.
     */
    private void pauseFeaturesForBackground() {
        mainHandler.post(() -> {
            SensorManager sm = (SensorManager) activity.getSystemService(Context.SENSOR_SERVICE);
            if (sm != null) {
                if (accelListener != null) { sm.unregisterListener(accelListener); accelListener = null; }
                if (proximityListener != null) { sm.unregisterListener(proximityListener); proximityListener = null; }
            }
            if (torchOn) {
                setTorch(0); // turns the torch off, emits common:torch and clears torchOn
            }
        });
    }

    /** Re-enable on foreground (onStart) any sensors the user had switched on. */
    private void resumeFeaturesForForeground() {
        if (motionWanted && accelListener == null) setMotion(1);
        if (proximityWanted && proximityListener == null) setProximity(1);
    }

    /** Speak text via TextToSpeech (lazily initialised). */
    public void speak(final String text) {
        mainHandler.post(() -> {
            if (text == null || text.isEmpty()) return;
            if (tts == null) {
                tts = new TextToSpeech(activity, status -> {
                    if (status == TextToSpeech.SUCCESS && tts != null) {
                        tts.setLanguage(Locale.US);
                        tts.speak(text, TextToSpeech.QUEUE_FLUSH, null, "wails");
                    }
                });
            } else {
                tts.speak(text, TextToSpeech.QUEUE_FLUSH, null, "wails");
            }
        });
    }

    /** Stop any in-progress speech. */
    public void stopSpeak() {
        mainHandler.post(() -> {
            if (tts != null) tts.stop();
        });
    }

    /** Disk space as {"free":bytes,"total":bytes}. */
    public String getStorageJson() {
        try {
            StatFs stat = new StatFs(activity.getFilesDir().getAbsolutePath());
            long free = stat.getAvailableBytes();
            long total = stat.getTotalBytes();
            return new JSONObject().put("free", free).put("total", total).toString();
        } catch (Exception e) {
            return "{\"free\":0,\"total\":0}";
        }
    }

    /** Absolute path to the app's private internal files directory
     *  (getFilesDir()), suitable for databases and other persistent files. */
    public String getStoragePath() {
        java.io.File dir = activity.getFilesDir();
        return dir != null ? dir.getAbsolutePath() : "";
    }

    /** Battery/power state as {"level":0-1,"charging":bool,"lowPower":bool}. */
    public String getPowerJson() {
        try {
            IntentFilter filter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
            android.content.Intent batt = activity.registerReceiver(null, filter);
            float level = -1;
            boolean charging = false;
            if (batt != null) {
                int raw = batt.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
                int scale = batt.getIntExtra(BatteryManager.EXTRA_SCALE, -1);
                if (raw >= 0 && scale > 0) level = raw / (float) scale;
                int status = batt.getIntExtra(BatteryManager.EXTRA_STATUS, -1);
                charging = status == BatteryManager.BATTERY_STATUS_CHARGING
                        || status == BatteryManager.BATTERY_STATUS_FULL;
            }
            boolean lowPower = false;
            PowerManager pm = (PowerManager) activity.getSystemService(Context.POWER_SERVICE);
            if (pm != null && Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
                lowPower = pm.isPowerSaveMode();
            }
            return new JSONObject().put("level", level)
                    .put("charging", charging).put("lowPower", lowPower).toString();
        } catch (Exception e) {
            return "{\"level\":-1,\"charging\":false,\"lowPower\":false}";
        }
    }

    /** Network status as {"connected":bool,"type":"wifi|cellular|ethernet|none"}. */
    public String getNetworkJson() {
        try {
            ConnectivityManager cm =
                    (ConnectivityManager) activity.getSystemService(Context.CONNECTIVITY_SERVICE);
            boolean connected = false;
            String type = "none";
            if (cm != null && Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                Network net = cm.getActiveNetwork();
                NetworkCapabilities caps = net != null ? cm.getNetworkCapabilities(net) : null;
                if (caps != null) {
                    connected = caps.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET);
                    if (caps.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)) type = "wifi";
                    else if (caps.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR)) type = "cellular";
                    else if (caps.hasTransport(NetworkCapabilities.TRANSPORT_ETHERNET)) type = "ethernet";
                }
            }
            return new JSONObject().put("connected", connected).put("type", type).toString();
        } catch (Exception e) {
            return "{\"connected\":false,\"type\":\"none\"}";
        }
    }

    /**
     * Watch (1) / unwatch (0) the soft keyboard, emitting "common:keyboard"
     * {visible,height} (height in px) via an inset listener on the content view.
     */
    public void setKeyboardWatch(final int enabled) {
        mainHandler.post(() -> {
            final View content = activity.getWindow().getDecorView();
            if (enabled != 0) {
                if (keyboardListener != null) return;
                keyboardListener = (v, insets) -> {
                    int imeHeight = 0;
                    boolean visible;
                    if (Build.VERSION.SDK_INT >= 30) {
                        imeHeight = insets.getInsets(WindowInsets.Type.ime()).bottom;
                        visible = insets.isVisible(WindowInsets.Type.ime());
                    } else {
                        imeHeight = insets.getSystemWindowInsetBottom();
                        visible = imeHeight > 0;
                    }
                    emitKeyboard(visible, imeHeight);
                    return insets;
                };
                content.setOnApplyWindowInsetsListener(keyboardListener);
                content.requestApplyInsets();
            } else if (keyboardListener != null) {
                content.setOnApplyWindowInsetsListener(null);
                keyboardListener = null;
            }
        });
    }

    private void emitKeyboard(boolean visible, int height) {
        try {
            emitEvent("common:keyboard",
                    new JSONObject().put("visible", visible).put("height", height).toString());
        } catch (Exception ignored) {
        }
    }

    /**
     * Toggle FLAG_SECURE (blocks screenshots & screen recording). Reports the new
     * state as "common:screenCapture" {protected}.
     */
    public void setScreenProtect(final int enabled) {
        mainHandler.post(() -> {
            if (enabled != 0) {
                activity.getWindow().addFlags(WindowManager.LayoutParams.FLAG_SECURE);
            } else {
                activity.getWindow().clearFlags(WindowManager.LayoutParams.FLAG_SECURE);
            }
            emitEvent("common:screenCapture",
                    "{\"protected\":" + (enabled != 0 ? "true" : "false") + "}");
        });
    }

    // MARK: - Mobile features (Phase E: camera & background)

    /** Capture a photo with the system camera. Result → "common:capture" event. */
    public void capturePhoto(final String json) {
        mainHandler.post(() -> {
            if (activity instanceof MainActivity) {
                ((MainActivity) activity).launchCameraCapture(false);
            } else {
                emitEvent("common:capture", "{\"error\":\"camera unavailable\"}");
            }
        });
    }

    /** Capture a video with the system camera. Result → "common:capture" event. */
    public void captureVideo(final String json) {
        mainHandler.post(() -> {
            if (activity instanceof MainActivity) {
                ((MainActivity) activity).launchCameraCapture(true);
            } else {
                emitEvent("common:capture", "{\"error\":\"camera unavailable\"}");
            }
        });
    }

    /**
     * Start a foreground service that keeps the process alive for long-running
     * background work (with an ongoing notification). json: {"title","text"}.
     */
    public void startForegroundService(final String json) {
        mainHandler.post(() -> {
            try {
                String title = "Wails", text = "Running in the background";
                try {
                    JSONObject o = new JSONObject(json);
                    title = o.optString("title", title);
                    text = o.optString("text", text);
                } catch (Exception ignored) {
                }
                if (Build.VERSION.SDK_INT >= 33 && activity.checkSelfPermission(
                        "android.permission.POST_NOTIFICATIONS") != PackageManager.PERMISSION_GRANTED) {
                    activity.requestPermissions(new String[]{"android.permission.POST_NOTIFICATIONS"}, 1003);
                }
                Intent i = new Intent(activity, WailsForegroundService.class);
                i.setAction(WailsForegroundService.ACTION_START);
                i.putExtra("title", title);
                i.putExtra("text", text);
                ContextCompat.startForegroundService(activity, i);
                emitEvent("android:foregroundService", "{\"running\":true}");
            } catch (Exception e) {
                Log.e(TAG, "startForegroundService failed", e);
                emitEvent("android:foregroundService", "{\"running\":false,\"error\":\"failed to start\"}");
            }
        });
    }

    /** Stop the foreground service. */
    public void stopForegroundService() {
        mainHandler.post(() -> {
            try {
                activity.stopService(new Intent(activity, WailsForegroundService.class));
                emitEvent("android:foregroundService", "{\"running\":false}");
            } catch (Exception e) {
                Log.e(TAG, "stopForegroundService failed", e);
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
