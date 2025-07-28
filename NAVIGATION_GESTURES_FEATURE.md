# WKWebView Navigation Gestures Support

This document describes the implementation of `allowsBackForwardNavigationGestures` support in Wails v3 for macOS.

## Problem

The GitHub issue [#1857](https://github.com/wailsapp/wails/issues/1857) requested support for Mac's two-finger swipe navigation gestures in WKWebView. Without this feature, users couldn't use the native macOS horizontal swipe gestures to navigate back and forward in the webview.

## Solution

Added support for the `allowsBackForwardNavigationGestures` property by:

1. **Extended MacWebviewPreferences**: Added `AllowsBackForwardNavigationGestures u.Bool` to the Mac webview preferences struct
2. **Updated C struct**: Added the corresponding field to the C struct that bridges Go and Objective-C
3. **Implemented WKWebView configuration**: Added code to set the property on the WKWebView during initialization
4. **Updated preference handling**: Extended the preference processing to handle the new setting

## Usage

To enable navigation gestures in your Wails v3 application:

```go
package main

import (
    "github.com/leaanthony/u"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    window := app.NewWebviewWindow()
    window.SetOptions(application.WebviewWindowOptions{
        Title: "Navigation Gestures Demo",
        Mac: application.MacWindow{
            WebviewPreferences: application.MacWebviewPreferences{
                // Enable horizontal swipe gestures for back/forward navigation
                AllowsBackForwardNavigationGestures: u.True(),
            },
        },
        // Your other options...
    })

    // ... rest of your app setup
}
```

## How It Works

1. When `AllowsBackForwardNavigationGestures` is set to `u.True()`, the preference is passed to the native code
2. During WKWebView initialization, the `allowsBackForwardNavigationGestures` property is set on the webview
3. macOS automatically handles the two-finger horizontal swipe gestures to trigger back/forward navigation
4. The gestures work with the webview's navigation history

## Files Modified

- `v3/pkg/application/webview_window_options.go`: Added the new preference field
- `v3/pkg/application/webview_window_darwin.go`: 
  - Updated C struct definition
  - Added WKWebView property setting
  - Extended preference processing

## Compatibility

- **Platform**: macOS only (WKWebView specific)
- **macOS Version**: All versions that support WKWebView
- **Wails Version**: v3+

## Testing

The implementation has been tested by:
1. Building the modified code successfully
2. Creating example code that demonstrates the feature
3. Verifying the C bridge code compiles without errors

## Benefits

- **Native UX**: Provides the standard macOS navigation experience users expect
- **No JavaScript required**: Uses native WebKit functionality instead of custom JS implementations
- **Better performance**: Native gestures feel more responsive than JS alternatives
- **Accessibility**: Maintains native accessibility features

## Related

- Original issue: https://github.com/wailsapp/wails/issues/1857
- Apple Documentation: [allowsBackForwardNavigationGestures](https://developer.apple.com/documentation/webkit/wkwebview/1414995-allowsbackforwardnavigationgestu)