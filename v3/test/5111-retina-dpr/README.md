# Test: WKWebView devicePixelRatio on Retina (#5111)

## Issue
WKWebView reports `devicePixelRatio=1` on macOS Retina displays when content is loaded via custom URL scheme (`wails://`).

## Fix
Added `contentsScale` and `rasterizationScale` configuration to the WKWebView layer during initialization in `windowNew()` (`v3/pkg/application/webview_window_darwin.go`), and updated `windowDidChangeBackingProperties:` in `webview_window_darwin.m` to update the scale when the window moves between displays.

## Manual Test Steps

1. Build a Wails v3 app on a Mac with a Retina display
2. Run the app and open the web inspector console
3. Check `window.devicePixelRatio` — should be `2` (not `1`)
4. Check `window.matchMedia('(-webkit-min-device-pixel-ratio: 2)').matches` — should be `true`
5. Test canvas rendering — text and graphics should be sharp, not blurry
6. Move the window to an external non-Retina display — `devicePixelRatio` should update to `1`
7. Move back to Retina — `devicePixelRatio` should return to `2`

## Expected Code Changes

### webview_window_darwin.go (windowNew)
After WKWebView creation, the layer scale is set:
```objc
[webView setWantsLayer:YES];
if (webView.layer) {
    webView.layer.contentsScale = [[NSScreen mainScreen] backingScaleFactor];
    webView.layer.rasterizationScale = [[NSScreen mainScreen] backingScaleFactor];
    webView.layer.shouldRasterize = YES;
}
```

### webview_window_darwin.m (windowDidChangeBackingProperties)
The delegate method now updates the webview layer scale when the backing properties change (e.g., window moved between Retina and non-Retina displays).
