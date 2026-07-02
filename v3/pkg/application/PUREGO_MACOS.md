# CGO-free macOS backend (purego)

This is a port of the Wails v3 **macOS** backend that runs with `CGO_ENABLED=0`
by driving the Objective-C runtime through
[`github.com/ebitengine/purego`](https://github.com/ebitengine/purego) and its
`objc` subpackage instead of compiling Objective-C via cgo.

## Building

```sh
CGO_ENABLED=0 go build -tags purego ./...
```

The `purego` build tag selects this backend; without it the existing cgo
backend is used (unchanged). The two are mutually exclusive: every cgo darwin
file carries `&& !purego`, and every file here carries `&& purego`.

Frameworks (AppKit, WebKit, Foundation, CoreGraphics, ...) are `dlopen`ed at
runtime rather than linked at build time, so a purego binary links only
libSystem/CoreFoundation/Security.

## Design

- `darwin_purego_cocoa.go` — the foundation. A thin `id` wrapper over `objc.ID`
  with `send`/`get[T]`, selector caching, class lookup, NSString/NSURL/NSData
  helpers, CG/NS geometry structs, autorelease pools and a
  `registerDelegateClass` helper.
- Objective-C delegates (NSApplicationDelegate, NSWindowDelegate,
  WKScriptMessageHandler, WKURLSchemeHandler, WKNavigationDelegate, ...) are
  created at runtime with `objc.RegisterClass`; their methods are Go closures
  (IMPs) that push onto the same channels the cgo `//export` callbacks used.
- The main-thread dispatcher uses libdispatch (`dispatch_async` onto
  `_dispatch_main_q`) with an `objc.NewBlock`.

## Verified working (CGO_ENABLED=0)

`examples/window` builds and launches: NSApplication lifecycle fires
`ApplicationDidFinishLaunching` through the runtime-registered delegate, the
WKWebView is created, and the `wails://` scheme handler serves the frontend
(asset-server request observed). JS↔Go bridge (`external` message handler) and
the window-event delegate are wired.

## Parity status

| Area | State |
|------|-------|
| App lifecycle, delegate, theme/power events, terminate, open-file/url | implemented |
| Main-thread dispatch, run loop | implemented |
| Window: create, show/hide/close/destroy, focus, center | implemented |
| Window: size/position/bounds, min/max, resizable, always-on-top | implemented |
| Window: minimise/maximise/fullscreen + state queries | implemented |
| WKWebView: setURL/setHTML/execJS, reload, zoom, devtools | implemented |
| Asset serving (WKURLSchemeHandler), JS message bridge | implemented |
| Window events (~60 NSWindowDelegate notifications) | implemented |
| Drag/drop throttling (frontend drop) | ported (native drag view pending) |
| Clipboard, screens | implemented |
| Menu / menu items / roles, system tray | implemented |
| Dialogs (message/open/save) | implemented (app-modal; sheets pending) |
| Global shortcuts (Carbon), autostart (SMAppService), single-instance | implemented |
| Bundle id (`pkg/mac`), events | implemented |
| Dock badge, notifications services | implemented |
| Window levels, collection behavior, window buttons, content protection, ignore-mouse | implemented |
| Titlebar presets (transparent/hide/full-size/toolbar/style/separator), appearance, backdrop (transparent/translucent) | implemented |
| Context menu, print, attach-modal (sheet), start-drag, disable-size-constraints, frameless toggle, CSS injection, key-event routing | implemented |
| Frameless native title-bar drag (invisible title-bar mouse monitor), show-toolbar-when-fullscreen | implemented |
| Disable-escape-exits-fullscreen, frameless keyboard focus (canBecomeKey/Main) | implemented (NSWindow subclass) |
| Liquid glass backdrop (macOS 26 private NSGlass APIs) | **remaining** (no-op; window renders normally) |
| Native file-drop overlay view (drag-drop of files) | implemented (NSView <NSDraggingDestination>) |
| Window key-binding capture (keyDown: accelerator mapping) | implemented |

## Known gaps

- Writing to an already-stopped `WKURLSchemeTask` cannot convert WebKit's
  `NSException` into `errRequestStopped` without a tiny native shim (no
  `@try/@catch` from pure Go); the exception would otherwise propagate.
- Dialogs are application-modal rather than window sheets.
- Liquid glass backdrop is a no-op (requires macOS 26 private APIs).
