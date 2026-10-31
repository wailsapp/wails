# CGO-free Linux backend (purego)

This is a port of the Wails v3 **Linux** backend (GTK4 + WebKitGTK 6.0) that
runs with `CGO_ENABLED=0` by loading the system libraries at runtime through
[`github.com/ebitengine/purego`](https://github.com/ebitengine/purego)
instead of compiling C via cgo.

## Building

```sh
CGO_ENABLED=0 go build -tags purego ./...       # GTK4 + WebKitGTK 6.0
CGO_ENABLED=0 go build -tags purego,gtk3 ./...  # GTK3 + WebKit2GTK 4.1 (legacy)
CGO_ENABLED=1 go build ./...                    # the cgo backend (default)
```

The purego backend is experimental and opt-in, as on macOS: every purego Linux
file carries `purego` and every cgo Linux file `!purego`, so exactly one backend
is compiled. `CGO_ENABLED=0` without the tag fails with one error
(`internal/buildguard`). `wails3` adds the tag when `WAILS_PUREGO=1` is set.
Cross-compiling from any host works (`GOOS=linux GOARCH=amd64 CGO_ENABLED=0`),
since no C toolchain or Linux sysroot is involved.

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

Minimum versions: GTK 4.8 and WebKitGTK 2.40 — Debian 12, the oldest distro
shipping WebKitGTK 6.0 (see `internal/puregocheck/floors.json`; `puregocheck`
enforces it). Because symbols are resolved by name at runtime, there is no
compile-time SDK ceiling: newer-than-floor functions are resolved with
`registerOptional`/`haveSymbol` and nil-checked before use. Examples:

- `gdk_monitor_get_scale` (GTK 4.14+, fractional scaling) falls back to the
  integer `gdk_monitor_get_scale_factor` on older GTK4 — the cgo build refuses
  to start on those systems; the purego build degrades gracefully.
- GtkFileDialog (GTK 4.10+) is optional. `dialogUsesGtkFileDialog()` picks it
  when it is bound; otherwise the file dialogs fall back to the synchronous
  GtkFileChooserWidget in a modal window (`showFileChooserDialog`), whose
  symbols exist throughout GTK 4. `gtk_file_chooser_dialog_new` is variadic
  and `gtk_dialog_run` spins a nested main loop, so the dialog is assembled by
  hand and the wait iterates the default main context, like the clipboard. The
  window's close-request only records the cancellation and returns TRUE; the
  wait loop is the only code that destroys the window.
- `gtk_file_dialog_set_initial_name` (save dialog's file name) is newer than
  4.8 too, and is bound optionally and guarded like the rest of GtkFileDialog.
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

## Known limitations shared with the cgo backend

- Menu accelerators are app-wide on Linux, in both backends and both stacks:
  GTK4 binds them to `app.*` actions and GTK3 to one accelerator group that
  every window holds. If two windows' own menus bind the same shortcut, the
  item registered first can fire whichever window has focus. Windows keeps
  bindings per window; macOS has a single application menu.

## Verified

Checked on every change by `ops/vm-pg.sh` (lab VMs: Ubuntu 24.04 arm64 and
amd64, both stacks):

- `puregocheck static` and `resolve -gir` report no stubs, no missing
  symbols and no unguarded calls newer than `floors.json`.
- Unit tests with coverage, the native parity probe (`tests/purego`), the
  cgo parity run of the same probe, and the e2e probe (`tests/purego/e2e`,
  real input through XTEST).
- Visual scenarios (`tests/purego/visual`) with screenshots as evidence.
- Every example under `v3/examples` vets in purego GTK4, purego GTK3 and cgo.
