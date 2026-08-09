# WebContentsView for Wails v3

`WebContentsView` embeds a second, native browser surface in a Wails window.
It is useful when a product needs a browser panel that follows a desktop layout
without turning that panel into an iframe inside the application's frontend.

The implementation uses WKWebView on macOS, WebView2 on Windows, and
WebKitGTK on Linux. The default Linux backend is GTK4 + WebKitGTK 6.0; the
legacy GTK3 backend remains available with `-tags gtk3` for one v3 cycle.

## Lifecycle

Create the view in Go, add it to the owning `WebviewWindow`, and update its
native rectangle as the frontend layout changes. A view can be registered
before `App.Run`; native browser creation is deferred until the host window is
ready.

```go
browser := webcontentsview.NewWebContentsView(webcontentsview.WebContentsViewOptions{
	Name:   "documentation",
	URL:    "https://wails.io",
	Bounds: application.Rect{X: 24, Y: 96, Width: 900, Height: 600},
	WebPreferences: webcontentsview.WebPreferences{
		DevTools:   webcontentsview.Enabled,
		Javascript: webcontentsview.Enabled,
	},
})

window.AddChildView(browser)
```

`SetBounds` moves and resizes the native surface. Use `Show` and `Hide` to
switch visibility while retaining its browser session. `RemoveChildView`
detaches it and permits a later re-add; `Destroy` permanently releases its
native resources and cannot be reversed.

```go
browser.SetBounds(application.Rect{X: 300, Y: 96, Width: 600, Height: 600})
browser.Hide()
browser.Show()
window.RemoveChildView(browser)
browser.Destroy()
```

The public API also supports `SetURL`, `SetHTML`, `ExecJS`, `GoBack`,
`GetURL`, and `TakeSnapshot`. `SetURL` and `SetHTML` replace one another as
the pending source. Native child surfaces always render above the host Wails
webview, so frontend DOM stacking cannot overlay a visible `WebContentsView`.

## Platform notes

- **Windows:** each view has an app- and view-specific WebView2 user-data
  directory, avoiding an incompatible environment with the host Wails webview.
- **Linux:** native views are positioned as direct GTK overlay children, so
  uncovered host UI remains interactive.
- **macOS:** WKWebView creation is deferred until an NSWindow is attached.

`WebContentsView` provides a native browser panel with a deliberately small,
platform-neutral lifecycle API.
