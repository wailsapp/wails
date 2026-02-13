package com.wails.app;

import android.app.Activity;
import android.content.Context;
import android.os.Build;
import android.util.Log;
import android.view.Menu;
import android.view.MotionEvent;
import android.view.View;
import android.webkit.WebSettings;
import android.webkit.WebView;

import com.google.android.material.bottomnavigation.BottomNavigationView;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
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
    private BottomNavigationView nativeTabsView;
    private final CopyOnWriteArrayList<String> nativeTabTitles = new CopyOnWriteArrayList<>();
    private volatile boolean nativeTabsEnabled = false;
    private volatile boolean scrollEnabled = true;
    private volatile boolean bounceEnabled = true;
    private volatile boolean scrollIndicatorsEnabled = true;
    private volatile boolean backForwardGesturesEnabled = false;
    private volatile boolean linkPreviewEnabled = true;
    private volatile String customUserAgent = null;
    private volatile float touchStartX = 0f;
    private volatile float touchStartY = 0f;
    private volatile int swipeThresholdPx = 120;
    private volatile View.OnTouchListener touchListener;
    private final View.OnLongClickListener blockLongClickListener = v -> true;

    private static final List<String> DEFAULT_NATIVE_TAB_TITLES = Collections.emptyList();

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
        applyWebViewSettings();
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

    public String getDeviceInfoJson() {
        JSONObject info = new JSONObject();
        try {
            info.put("platform", "android");
            info.put("model", Build.MODEL);
            info.put("version", Build.VERSION.RELEASE);
            info.put("manufacturer", Build.MANUFACTURER);
            info.put("brand", Build.BRAND);
            info.put("device", Build.DEVICE);
            info.put("product", Build.PRODUCT);
            info.put("sdkInt", Build.VERSION.SDK_INT);
        } catch (JSONException e) {
            Log.w(TAG, "Failed to build device info JSON", e);
        }
        return info.toString();
    }

    public void setScrollEnabled(boolean enabled) {
        scrollEnabled = enabled;
        applyWebViewSettings();
    }

    public void setBounceEnabled(boolean enabled) {
        bounceEnabled = enabled;
        applyWebViewSettings();
    }

    public void setScrollIndicatorsEnabled(boolean enabled) {
        scrollIndicatorsEnabled = enabled;
        applyWebViewSettings();
    }

    public void setBackForwardGesturesEnabled(boolean enabled) {
        backForwardGesturesEnabled = enabled;
        applyWebViewSettings();
    }

    public void setLinkPreviewEnabled(boolean enabled) {
        linkPreviewEnabled = enabled;
        applyWebViewSettings();
    }

    public void setCustomUserAgent(String userAgent) {
        customUserAgent = userAgent;
        applyWebViewSettings();
    }

    private void applyWebViewSettings() {
        if (webView == null) {
            return;
        }

        webView.post(() -> {
            if (webView == null) {
                return;
            }

            WebSettings settings = webView.getSettings();
            if (customUserAgent == null || customUserAgent.trim().isEmpty()) {
                settings.setUserAgentString(null);
            } else {
                settings.setUserAgentString(customUserAgent);
            }

            webView.setVerticalScrollBarEnabled(scrollIndicatorsEnabled);
            webView.setHorizontalScrollBarEnabled(scrollIndicatorsEnabled);
            webView.setOverScrollMode(bounceEnabled ? View.OVER_SCROLL_IF_CONTENT_SCROLLS : View.OVER_SCROLL_NEVER);

            if (!linkPreviewEnabled) {
                webView.setOnLongClickListener(blockLongClickListener);
            } else {
                webView.setOnLongClickListener(null);
            }

            float density = webView.getResources().getDisplayMetrics().density;
            swipeThresholdPx = (int) (density * 120f);
            updateTouchListener();
        });
    }

    private void updateTouchListener() {
        if (webView == null) {
            return;
        }

        if (!scrollEnabled || backForwardGesturesEnabled) {
            if (touchListener == null) {
                touchListener = (v, event) -> {
                    if (webView == null) {
                        return false;
                    }

                    switch (event.getAction()) {
                        case MotionEvent.ACTION_DOWN:
                            touchStartX = event.getX();
                            touchStartY = event.getY();
                            return false;
                        case MotionEvent.ACTION_MOVE:
                            return !scrollEnabled;
                        case MotionEvent.ACTION_UP:
                            if (backForwardGesturesEnabled) {
                                float dx = event.getX() - touchStartX;
                                float dy = event.getY() - touchStartY;
                                if (Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > swipeThresholdPx) {
                                    if (dx > 0 && webView.canGoBack()) {
                                        webView.goBack();
                                        return true;
                                    }
                                    if (dx < 0 && webView.canGoForward()) {
                                        webView.goForward();
                                        return true;
                                    }
                                }
                            }
                            return false;
                        default:
                            return false;
                    }
                };
            }
            webView.setOnTouchListener(touchListener);
        } else {
            webView.setOnTouchListener(null);
        }
    }

    /**
     * Enable or disable native tabs on Android.
     */
    public void setNativeTabsEnabled(boolean enabled) {
        nativeTabsEnabled = enabled;
        applyNativeTabs();
    }

    /**
     * Configure native tab items via JSON array: [{"Title":"..."}]
     */
    public void setNativeTabsItemsJson(String json) {
        List<String> titles = new ArrayList<>();
        if (json != null && !json.trim().isEmpty()) {
            try {
                JSONArray arr = new JSONArray(json);
                for (int i = 0; i < arr.length(); i++) {
                    Object entry = arr.get(i);
                    if (!(entry instanceof JSONObject)) {
                        continue;
                    }
                    JSONObject obj = (JSONObject) entry;
                    String title = obj.optString("Title", "");
                    titles.add(title);
                }
            } catch (JSONException e) {
                Log.w(TAG, "Failed to parse native tabs JSON", e);
            }
        }

        nativeTabTitles.clear();
        nativeTabTitles.addAll(titles);
        if (!nativeTabTitles.isEmpty()) {
            nativeTabsEnabled = true;
        }
        applyNativeTabs();
    }

    /**
     * Programmatically select a native tab index.
     */
    public void selectNativeTabIndex(int index) {
        Activity activity = getActivity();
        if (activity == null) {
            return;
        }

        activity.runOnUiThread(() -> {
            BottomNavigationView tabs = ensureNativeTabsView(activity);
            if (tabs == null) {
                return;
            }
            int count = tabs.getMenu().size();
            if (index < 0 || index >= count) {
                return;
            }
            tabs.setSelectedItemId(index);
        });
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

    private void applyNativeTabs() {
        Activity activity = getActivity();
        if (activity == null) {
            Log.w(TAG, "Native tabs unavailable: context is not an Activity");
            return;
        }

        activity.runOnUiThread(() -> {
            BottomNavigationView tabs = ensureNativeTabsView(activity);
            if (tabs == null) {
                Log.w(TAG, "Native tabs view not found in layout");
                return;
            }

            List<String> titles = nativeTabTitles;
            if (titles.isEmpty()) {
                titles = DEFAULT_NATIVE_TAB_TITLES;
            }

            boolean shouldShow = nativeTabsEnabled && !titles.isEmpty();
            if (!shouldShow) {
                tabs.setVisibility(View.GONE);
                return;
            }

            Menu menu = tabs.getMenu();
            menu.clear();
            for (int i = 0; i < titles.size(); i++) {
                String title = titles.get(i);
                if (title == null) {
                    title = "";
                }
                menu.add(Menu.NONE, i, i, title).setIcon(android.R.drawable.ic_menu_view);
            }

            tabs.setOnItemSelectedListener(null);
            tabs.setSelectedItemId(0);
            tabs.setOnItemSelectedListener(item -> {
                dispatchNativeTabSelected(item.getItemId());
                return true;
            });

            tabs.setVisibility(View.VISIBLE);
        });
    }

    private void dispatchNativeTabSelected(int index) {
        String js = "window.dispatchEvent(new CustomEvent('nativeTabSelected',{detail:{index:" + index + "}}));";
        executeJavaScript(js);
    }

    private BottomNavigationView ensureNativeTabsView(Activity activity) {
        if (nativeTabsView != null) {
            return nativeTabsView;
        }
        nativeTabsView = activity.findViewById(R.id.native_tabs);
        return nativeTabsView;
    }

    private Activity getActivity() {
        if (context instanceof Activity) {
            return (Activity) context;
        }
        return null;
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
