# CGO-free Linux backend (purego)

This is a port of the Wails v3 **Linux** backend (GTK4 + WebKitGTK 6.0) that
runs with `CGO_ENABLED=0` by loading the system libraries at runtime through
[`github.com/ebitengine/purego`](https://github.com/ebitengine/purego)
instead of compiling C via cgo.

## Building

```sh
CGO_ENABLED=0 go build -tags purego ./...
```

The `purego` build tag selects this backend; without it the existing cgo
backend is used (unchanged). The two are mutually exclusive: every cgo Linux
file carries `&& !purego`, and every file here carries `&& purego`.
Cross-compiling from any host works (`GOOS=linux GOARCH=amd64 CGO_ENABLED=0`),
since no C toolchain or Linux sysroot is involved.

The GTK3 variant (`-tags gtk3`) is **not** supported in combination with
`purego` — the purego backend targets the same default stack as the default
cgo build: GTK4 and WebKitGTK 6.0.

## Design

- `linux_purego_lib.go` — the foundation. `dlopen(3)`s GLib/GObject/Gio/
  GTK4/WebKitGTK-6.0/JavaScriptCore/libsoup (with per-distro soname
  fallbacks) and binds every C function to a typed Go function variable via
  `purego.RegisterFunc`. A missing library or symbol produces one aggregate,
  actionable error (which package to install per distro) instead of a
  dynamic-linker one-liner or a nil-pointer crash.
- `linux_purego_callbacks.go` — the pure-Go port of `linux_cgo.c`: GTK signal
  trampolines (`purego.NewCallback`, a fixed set created once — never per
  window/item/call), main-thread dispatch via `g_idle_add`, the GAction menu
  machinery, GTK4 file/message dialogs, drag-and-drop controllers, clipboard,
  and the X11 helpers (window move/position, always-on-top) resolved with
  `dlsym(RTLD_DEFAULT)` from GTK's own X11 backend — no libX11 link, no-ops on
  Wayland, exactly like the cgo backend.
- `linux_purego.go` — the port of `linux_cgo.go`: the full shim function
  surface the shared Linux files (`webview_window_linux.go`,
  `menu_linux.go`, `dialogs_linux.go`, …) compile against.
- `application_linux_purego.go` — the port of `application_linux.go`
  (GApplication lifecycle, dbus theme monitoring, screens cache).
- `global_shortcut_linux_x11_purego.go` — XGrabKey global shortcuts via
  dlopen'ed libX11 (the Wayland portal backend was already pure Go and is now
  shared between backends via a `(cgo || purego)` tag, as are the dbus theme
  monitor and permission helpers).
- `internal/assetserver/webview/*_linux_purego.go` — the WebKit URI-scheme
  request/response plumbing, preserving the #5631/#5668 main-thread
  confinement design (every WebKit/GObject touch hops to the GTK main loop;
  after the loop stops, dispatch is disabled and completions run inline).
- The signal-handler fix (SA_ONSTACK re-application, issue #5527, including
  the deliberate SIGUSR1 exemption for JavaScriptCore's GC) is preserved by
  calling libc `sigaction` through purego.

## Runtime requirements and capability guards

Runtime-loaded libraries (sonames tried in order):

| Library | Soname | Debian/Ubuntu package |
|---|---|---|
| GLib/GObject/Gio | `libglib-2.0.so.0`, `libgobject-2.0.so.0`, `libgio-2.0.so.0` | `libglib2.0-0` |
| GTK4 | `libgtk-4.so.1` | `libgtk-4-1` |
| WebKitGTK | `libwebkitgtk-6.0.so.4` | `libwebkitgtk-6.0-4` |
| JavaScriptCore | `libjavascriptcoregtk-6.0.so.1` | `libjavascriptcoregtk-6.0-1` |
| libsoup | `libsoup-3.0.so.0` | `libsoup-3.0-0` |
| libX11 (optional) | `libX11.so.6` | only needed for X11 global shortcuts |

Minimum versions: GTK 4.10 (GtkFileDialog) and WebKitGTK 2.40, matching the
cgo backend's compile floor. Because symbols are resolved by name at runtime,
there is no compile-time SDK ceiling: newer-than-floor functions are resolved
with `registerOptional`/`haveSymbol` and nil-checked before use. Current
examples:

- `gdk_monitor_get_scale` (GTK 4.14+, fractional scaling) falls back to the
  integer `gdk_monitor_get_scale_factor` on older GTK4 — the cgo build refuses
  to start on those systems; the purego build degrades gracefully.
- The GDK X11 symbols (`gdk_x11_display_get_xdisplay`, …) are optional: they
  don't exist in Wayland-only GTK builds, and every X11 helper no-ops without
  them.

When adding new library calls, follow the conventions documented at the top
of `linux_purego_lib.go` — in particular: purego cannot call variadic C
functions (use the `_value`/`_with_properties` variants), and any symbol newer
than the floor must be registered as optional and guarded.

## Bugs fixed relative to the cgo backend

The port is behaviour-parity, not bug-parity: defects found in the cgo
backend while translating it are fixed here and catalogued in
[BUGS_FOUND.md](BUGS_FOUND.md) (invalid free of a GLib-owned string, NULL
GList dereference, clipboard reentrancy, unsynchronised menu maps, a
file-dialog main-loop deadlock, a message-dialog zombie-window/UAF, and a
vestigial `import "C"` that silently made the system tray cgo-only).

## Verified

- `GOOS=linux GOARCH=amd64|arm64 CGO_ENABLED=0 go build -tags purego ./...`
- The default cgo build is unaffected (`go build ./...` on a Linux box).
- Runtime validation on a real Linux desktop: see the session notes / PR
  description for the exact checks performed.
