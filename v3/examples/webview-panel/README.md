# WebviewPanel example

A `WebviewPanel` embeds a native webview inside a `WebviewWindow`. Each panel has its own page and keyboard focus. Panels overlay the main webview and are positioned in its content coordinate system, using CSS pixels.

Run from this directory:

```sh
go run .
# Legacy Linux GTK3:
go run -tags gtk3 .
```

The demo initially loads a bundled local page. The buttons switch between that page, [Wails](https://wails.io), and [Google](https://www.google.com), or hide/show the panel. A `ResizeObserver` measures the placeholder and updates native bounds; no hard-coded browser measurements are needed.

## Create and manage panels in Go

Create panels before `app.Run()` or after the parent window is running:

```go
panel := window.NewPanel(application.WebviewPanelOptions{
    Name: "browser",
    URL: "/panel.html",
    X: 20, Y: 80, Width: 600, Height: 400,
    Anchor: application.AnchorFill,
    BackgroundColour: application.NewRGB(248, 250, 252),
})

panel.SetURL("https://wails.io")
panel.SetBounds(application.Rect{X: 20, Y: 80, Width: 500, Height: 300})
panel.SetZoom(1.25)
zoom := panel.GetZoom()
_ = zoom
panel.Hide()
panel.Show()
```

Names are unique within a window. `window.NewPanel` returns an existing named panel unchanged; it returns `nil` once window teardown begins. Empty names receive generated identifiers. Use `GetPanel(name)`, `GetPanelByID(id)`, or `GetPanels()` to retrieve panels. `Destroy`, `RemovePanel`, and closing the parent window release native resources; repeated destruction is safe.

`Headers` supplies custom headers on the initial navigation request only. `UserAgent` replaces the panel's browser user agent. Both are supported on desktop platforms. These options are copied when the panel is created. Later `SetURL` calls and reloads do not reapply initial headers; redirects follow the browser's own rules.

Go code owns the destination and headers. Use HTTPS when sending credentials, and only send them to endpoints whose redirect behavior you trust. Panels do not implement a credential store or an application-specific origin policy; custom header names are not automatically classified as secrets.

Local URLs are resolved through the Wails asset server. An empty URL leaves a new panel blank. `SetURL("")` and malformed URLs do not navigate. `URL()` reports the last URL requested by Go, not redirects or links followed inside the page.

## Frontend control and navigation policy

```js
import {Panel} from '@wailsio/runtime';

const panel = Panel.Get('browser');
await panel.SetBounds({x: 20, y: 80, width: 500, height: 300});
await panel.SetZoom(1.25);
const focused = await panel.IsFocused();
await panel.Hide();
await panel.Show();
```

The runtime exposes bounds, stacking order, visibility, zoom, focus, reload, developer tools, name, and destruction. It controls panels already created in Go; it cannot create them.

**Navigation is Go-only.** There are no frontend `SetURL`, `SetHTML`, or `ExecJS` methods, and raw calls to their reserved method IDs are rejected. If an application needs frontend-triggered navigation, expose a Go service with an application-specific allowlist. This example accepts a site key (`local`, `wails`, or `google`), never an arbitrary URL from JavaScript.

Panel pages do not receive the Wails message bridge or runtime automatically. An embedded browser is not a DOM iframe: CSS clipping, transforms, and DOM z-index in the parent do not control it. It also does not disable browser CORS or other web security rules.

## Layout

- `AnchorLeft | AnchorRight` stretches width while retaining edge margins.
- `AnchorTop | AnchorBottom` stretches height while retaining edge margins.
- A single right/bottom anchor moves the panel with that edge.
- `AnchorFill` retains all four margins; set bounds to the window size first for a full-window panel.
- Manual bounds changes establish a new resize baseline. Shrinking clamps dimensions to at least one pixel.

`FillWindow`, `DockLeft`, `DockRight`, `DockTop`, `DockBottom`, and `FillBeside` are one-time positioning helpers; set `Anchor` to maintain that layout on later resizes. `FillBeside` requires a different panel in the same window. For layouts driven by HTML, use `ResizeObserver` as this example does.

`ZIndex` orders panels relative to each other (higher values on top); equal values use creation order. All panels are above the main webview.

## Platform notes

| Platform | Implementation | Notes |
| --- | --- | --- |
| Windows | WebView2 child windows | Uses parent-window DPI for local coordinates. `ForceReload` currently falls back to ordinary reload. |
| macOS | WKWebView child views | Developer tools require a development build or the `devtools` build tag. |
| Linux | WebKitGTK overlay children | GTK4 is the default; `-tags gtk3` selects the legacy backend. |
| Android, iOS, server | No native panel support | Desktop panel APIs are not a portable mobile or server UI. |

`DevToolsEnabled` overrides the application's debug-mode default. An explicit `false` also prevents `OpenDevTools` and inspector-on-startup calls. The native panel is rectangular; use its background and bounds to integrate it into the parent UI.

## Verification

From `v3`:

```sh
go test -race ./pkg/application -run TestPanel
node --test ./examples/webview-panel/layout.test.mjs
go run ./test/manual/webview-panel
# On a Linux graphical desktop, or under xvfb-run:
go run -tags gtk3 ./test/manual/webview-panel
```

The smoke app checks real runtime dispatch, local asset loading, initial HTTP headers, user agent, visibility, bounds, zoom, dynamic creation, and concurrent destruction. It exits with `PASS` when all checks finish. Run the interactive example as well to inspect keyboard focus, stacking, resize behavior, and mixed-DPI displays.

The frontend tests in `src/panel.test.js` compare every exported panel method with the actual Go method-ID table and protect the `IsFocused`/`Destroy` regression. Runtime bundles and the npm package must both be rebuilt after changing TypeScript sources.
