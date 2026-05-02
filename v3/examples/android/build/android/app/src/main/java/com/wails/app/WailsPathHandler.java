package com.wails.app;

import android.net.Uri;
import android.util.Log;
import android.webkit.WebResourceResponse;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.webkit.WebViewAssetLoader;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

/**
 * WailsPathHandler implements WebViewAssetLoader.PathHandler to serve assets
 * from the Go asset server. This allows the WebView to load assets without
 * using a network server, similar to iOS's WKURLSchemeHandler.
 */
public class WailsPathHandler implements WebViewAssetLoader.PathHandler {
    private static final String TAG = "WailsPathHandler";

    private final WailsBridge bridge;

    public WailsPathHandler(WailsBridge bridge) {
        this.bridge = bridge;
    }

    @Nullable
    @Override
    public WebResourceResponse handle(@NonNull String path) {
        Log.d(TAG, "Handling path: " + path);

        // Normalize path
        if (path.isEmpty() || path.equals("/")) {
            path = "/index.html";
        }

        // Get asset from Go
        byte[] data = bridge.serveAsset(path, "GET", "{}");

        if (data == null || data.length == 0) {
            Log.w(TAG, "Asset not found: " + path);
            return null; // Return null to let WebView handle 404
        }

        // Determine MIME type
        String mimeType = bridge.getAssetMimeType(path);
        Log.d(TAG, "Serving " + path + " with type " + mimeType + " (" + data.length + " bytes)");

        // Create response
        InputStream inputStream = new ByteArrayInputStream(data);
        Map<String, String> headers = new HashMap<>();
        headers.put("Access-Control-Allow-Origin", "*");
        headers.put("Cache-Control", "no-cache");

        return new WebResourceResponse(
                mimeType,
                "UTF-8",
                200,
                "OK",
                headers,
                inputStream
        );
    }

    /**
     * Determine MIME type from file extension
     */
    private String getMimeType(String path) {
        String lowerPath = path.toLowerCase();

        if (lowerPath.endsWith(".html") || lowerPath.endsWith(".htm")) {
            return "text/html";
        } else if (lowerPath.endsWith(".js") || lowerPath.endsWith(".mjs")) {
            return "application/javascript";
        } else if (lowerPath.endsWith(".css")) {
            return "text/css";
        } else if (lowerPath.endsWith(".json")) {
            return "application/json";
        } else if (lowerPath.endsWith(".png")) {
            return "image/png";
        } else if (lowerPath.endsWith(".jpg") || lowerPath.endsWith(".jpeg")) {
            return "image/jpeg";
        } else if (lowerPath.endsWith(".gif")) {
            return "image/gif";
        } else if (lowerPath.endsWith(".svg")) {
            return "image/svg+xml";
        } else if (lowerPath.endsWith(".ico")) {
            return "image/x-icon";
        } else if (lowerPath.endsWith(".woff")) {
            return "font/woff";
        } else if (lowerPath.endsWith(".woff2")) {
            return "font/woff2";
        } else if (lowerPath.endsWith(".ttf")) {
            return "font/ttf";
        } else if (lowerPath.endsWith(".eot")) {
            return "application/vnd.ms-fontobject";
        } else if (lowerPath.endsWith(".xml")) {
            return "application/xml";
        } else if (lowerPath.endsWith(".txt")) {
            return "text/plain";
        } else if (lowerPath.endsWith(".wasm")) {
            return "application/wasm";
        } else if (lowerPath.endsWith(".mp3")) {
            return "audio/mpeg";
        } else if (lowerPath.endsWith(".mp4")) {
            return "video/mp4";
        } else if (lowerPath.endsWith(".webm")) {
            return "video/webm";
        } else if (lowerPath.endsWith(".webp")) {
            return "image/webp";
        }

        return "application/octet-stream";
    }
}
