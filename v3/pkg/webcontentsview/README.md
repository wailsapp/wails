# WebContentsView for Wails v3

`WebContentsView` is an implementation of Electron's `WebContentsView` (formerly `BrowserView`) for Wails v3. It allows you to embed a fully native, secondary OS-level Webview directly over your Wails application UI.

Unlike a standard HTML `<iframe>`, this native view:
- Bypasses restrictive `X-Frame-Options` and `Content-Security-Policy: frame-ancestors` headers.
- Can have web security (CORS) disabled independently of the main Wails app context.
- Maintains its own session, cookies, and caching behavior.
- Renders with native performance using the OS's underlying web engine (WKWebView on macOS, WebView2 on Windows, WebKitGTK on Linux).

## Architecture

The module is built with a clean separation between the Go API and the platform-specific native implementations.

### 1. Go API Layer
The `webcontentsview` package exposes a structured API for managing the view lifecycle.
- **`NewWebContentsView(options)`**: Initializes the native OS webview but does not display it.
- **`Attach(window)`**: Mounts the webview to the provided Wails `application.Window` using its raw `NativeWindow()` pointer.
- **`SetBounds(rect)`**: Dynamically positions and sizes the view.
- **`SetURL(url)`**: Navigates the view.
- **`ExecJS(js)`**: Evaluates JavaScript inside the context of the secondary view.
- **`Detach()`**: Unmounts and hides the view.

### 2. Platform Specific Implementations
*   **macOS (`webcontentsview_darwin.m`)**: Creates a `WKWebView` via Objective-C. When attached, it gets added as a subview to the `NSWindow`'s `contentView`. To ensure it sits above the main Wails UI, it is backed by a CoreAnimation layer (`wantsLayer = YES`) and assigned an astronomical z-index (`zPosition = 9999.0`). It automatically adjusts web coordinates (top-left) to Cocoa coordinates (bottom-left).
*   **Windows (`webcontentsview_windows.go`)**: Leverages the `github.com/wailsapp/go-webview2/pkg/edge` package to create an `edge.Chromium` instance, embedding it directly into the parent window's `HWND`.
*   **Linux (`webcontentsview_linux.go`)**: Uses CGO and GTK to create a `WebKitSettings` and `GtkWidget` webview, packing it into the main `GtkBox` container.

### 3. Web Preferences
Inspired by Electron, `WebContentsViewOptions` accepts a `WebPreferences` struct. This passes down to `WKPreferences` / `ICoreWebView2Settings` to configure behavior:
- `DevTools`: Enable/disable the web inspector.
- `Javascript`: Enable/disable JS execution.
- `WebSecurity`: Disables cross-origin restrictions and allows local file URL access (crucial for local-dev previewing).
- `ZoomFactor`: Scales the viewport.

---

## Usage Guide

To use `WebContentsView`, you must coordinate between your Go backend and your JavaScript/React frontend.

### 1. Go Backend setup
Add the bridge methods to your Wails `App` struct so the frontend can control the view:

```go
import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/webcontentsview"
)

var browserView *webcontentsview.WebContentsView

func (a *App) InitBrowserView(x, y, width, height int, url string) {
	// ALL UI creation MUST happen on the main thread
	application.InvokeSync(func() {
		browserView = webcontentsview.NewWebContentsView(webcontentsview.WebContentsViewOptions{
			URL: url,
			Bounds: application.Rect{ X: x, Y: y, Width: width, Height: height },
			WebPreferences: webcontentsview.WebPreferences{
				DevTools:    application.Enabled,
				Javascript:  application.Enabled,
				WebSecurity: application.Disabled, // Ideal for bypassing CORS during local dev
			},
		})
		browserView.Attach(a.mainWindow)
	})
}

func (a *App) SetBrowserViewBounds(x, y, width, height int) {
	application.InvokeSync(func() {
		browserView.SetBounds(application.Rect{ X: x, Y: y, Width: width, Height: height })
	})
}
```

### 2. React Frontend setup
Instead of an `<iframe>`, render an empty `<div>` in React to act as a placeholder. Use a `ResizeObserver` to track the exact screen coordinates of the `<div>` and send them to the Go backend. Go will physically move the native window over the empty space.

```tsx
import { useRef, useEffect, useLayoutEffect } from "react";
import { desktopAPI } from "@/features/desktop/api";

export default function BrowserTab({ active, url }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const isInitializedRef = useRef(false);

  const updateBounds = () => {
    if (!containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    
    // Hide by setting dimensions to 0 when inactive
    if (!active || rect.width === 0) {
      desktopAPI.setBrowserViewBounds(0, 0, 0, 0);
      return;
    }

    desktopAPI.setBrowserViewBounds(
      Math.round(rect.x),
      Math.round(rect.y),
      Math.round(rect.width),
      Math.round(rect.height)
    );
  };

  // Initialize on mount
  useLayoutEffect(() => {
    if (!containerRef.current) return;
    
    if (!isInitializedRef.current) {
      const rect = containerRef.current.getBoundingClientRect();
      desktopAPI.initBrowserView(
        Math.round(rect.x), Math.round(rect.y), 
        active ? Math.round(rect.width) : 0, 
        active ? Math.round(rect.height) : 0, 
        url
      ).then(() => { isInitializedRef.current = true; });
    } else {
      updateBounds();
    }
  }, [active]);

  // Track window resizing and layout shifting
  useEffect(() => {
    if (!active) return;
    const observer = new ResizeObserver(() => setTimeout(updateBounds, 10));
    if (containerRef.current) observer.observe(containerRef.current);
    
    window.addEventListener('resize', updateBounds);
    return () => {
      observer.disconnect();
      window.removeEventListener('resize', updateBounds);
    };
  }, [active]);

  return (
    // The native webview will "float" exactly over this transparent div
    <div className="flex-1 w-full relative bg-transparent" ref={containerRef} />
  );
}
```

This pattern ensures the native `WebContentsView` stays perfectly synchronized with your React layout, mimicking the behaviour of a built-in browser component.
