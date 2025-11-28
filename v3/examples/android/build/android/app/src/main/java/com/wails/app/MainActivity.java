package com.wails.app;

import android.annotation.SuppressLint;
import android.os.Bundle;
import android.util.Log;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.webkit.WebViewAssetLoader;
import com.wails.app.BuildConfig;

/**
 * MainActivity hosts the WebView and manages the Wails application lifecycle.
 * It uses WebViewAssetLoader to serve assets from the Go library without
 * requiring a network server.
 */
public class MainActivity extends AppCompatActivity {
    private static final String TAG = "WailsActivity";
    private static final String WAILS_SCHEME = "https";
    private static final String WAILS_HOST = "wails.localhost";

    private WebView webView;
    private WailsBridge bridge;
    private WebViewAssetLoader assetLoader;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // Initialize the native Go library
        bridge = new WailsBridge(this);
        bridge.initialize();

        // Set up WebView
        setupWebView();

        // Load the application
        loadApplication();
    }

    @SuppressLint("SetJavaScriptEnabled")
    private void setupWebView() {
        webView = findViewById(R.id.webview);

        // Configure WebView settings
        WebSettings settings = webView.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setDatabaseEnabled(true);
        settings.setAllowFileAccess(false);
        settings.setAllowContentAccess(false);
        settings.setMediaPlaybackRequiresUserGesture(false);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);

        // Enable debugging in debug builds
        if (BuildConfig.DEBUG) {
            WebView.setWebContentsDebuggingEnabled(true);
        }

        // Set up asset loader for serving local assets
        assetLoader = new WebViewAssetLoader.Builder()
                .setDomain(WAILS_HOST)
                .addPathHandler("/", new WailsPathHandler(bridge))
                .build();

        // Set up WebView client to intercept requests
        webView.setWebViewClient(new WebViewClient() {
            @Nullable
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                String url = request.getUrl().toString();
                Log.d(TAG, "Intercepting request: " + url);

                // Handle wails.localhost requests
                if (request.getUrl().getHost() != null &&
                        request.getUrl().getHost().equals(WAILS_HOST)) {

                    // For wails API calls (runtime, capabilities, etc.), we need to pass the full URL
                    // including query string because WebViewAssetLoader.PathHandler strips query params
                    String path = request.getUrl().getPath();
                    if (path != null && path.startsWith("/wails/")) {
                        // Get full path with query string for runtime calls
                        String fullPath = path;
                        String query = request.getUrl().getQuery();
                        if (query != null && !query.isEmpty()) {
                            fullPath = path + "?" + query;
                        }
                        Log.d(TAG, "Wails API call detected, full path: " + fullPath);

                        // Call bridge directly with full path
                        byte[] data = bridge.serveAsset(fullPath, request.getMethod(), "{}");
                        if (data != null && data.length > 0) {
                            java.io.InputStream inputStream = new java.io.ByteArrayInputStream(data);
                            java.util.Map<String, String> headers = new java.util.HashMap<>();
                            headers.put("Access-Control-Allow-Origin", "*");
                            headers.put("Cache-Control", "no-cache");
                            headers.put("Content-Type", "application/json");

                            return new WebResourceResponse(
                                "application/json",
                                "UTF-8",
                                200,
                                "OK",
                                headers,
                                inputStream
                            );
                        }
                        // Return error response if data is null
                        return new WebResourceResponse(
                            "application/json",
                            "UTF-8",
                            500,
                            "Internal Error",
                            new java.util.HashMap<>(),
                            new java.io.ByteArrayInputStream("{}".getBytes())
                        );
                    }

                    // For regular assets, use the asset loader
                    return assetLoader.shouldInterceptRequest(request.getUrl());
                }

                return super.shouldInterceptRequest(view, request);
            }

            @Override
            public void onPageFinished(WebView view, String url) {
                super.onPageFinished(view, url);
                Log.d(TAG, "Page loaded: " + url);
                // Inject Wails runtime
                bridge.injectRuntime(webView, url);
            }
        });

        // Add JavaScript interface for Go communication
        webView.addJavascriptInterface(new WailsJSBridge(bridge, webView), "wails");
    }

    private void loadApplication() {
        // Load the main page from the asset server
        String url = WAILS_SCHEME + "://" + WAILS_HOST + "/";
        Log.d(TAG, "Loading URL: " + url);
        webView.loadUrl(url);
    }

    /**
     * Execute JavaScript in the WebView from the Go side
     */
    public void executeJavaScript(final String js) {
        runOnUiThread(() -> {
            if (webView != null) {
                webView.evaluateJavascript(js, null);
            }
        });
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (bridge != null) {
            bridge.onResume();
        }
    }

    @Override
    protected void onPause() {
        super.onPause();
        if (bridge != null) {
            bridge.onPause();
        }
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        if (bridge != null) {
            bridge.shutdown();
        }
        if (webView != null) {
            webView.destroy();
        }
    }

    @Override
    public void onBackPressed() {
        if (webView != null && webView.canGoBack()) {
            webView.goBack();
        } else {
            super.onBackPressed();
        }
    }
}
