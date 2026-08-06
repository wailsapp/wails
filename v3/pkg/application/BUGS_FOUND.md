# Bugs found in the cgo Linux backend during the purego port

The purego port's goal is behavioural parity with the cgo GTK4/WebKitGTK-6.0
backend — but not bug-for-bug parity. Translating `linux_cgo.go`/`linux_cgo.c`
line by line surfaced the defects below. Each is FIXED in the purego backend;
the cgo backend still has them (fixing it is follow-up work, tracked here so
the fixes aren't lost).

## 1. `appName()` frees a GLib-owned string (undefined behaviour)

`linux_cgo.go:129-133`: `g_get_application_name()` returns a pointer owned by
GLib ("the returned string is owned by GLib and must not be modified or
freed"), but the cgo code does `defer C.free(unsafe.Pointer(name))`. Freeing
memory allocated by GLib's allocator with libc `free()` — and freeing memory
that GLib still owns — is a use-after-free waiting to happen (GLib returns the
same pointer on the next call). It also crashes outright when the application
name was never set (`g_get_application_name` may return NULL, and the
subsequent `C.GoString(nil)` masks it while `free(NULL)` is a no-op — the
real damage is when it is non-NULL).

**purego fix:** copy the string, never free it (`linux_purego.go: appName`).

## 2. `getWindows()` dereferences a NULL GList head

`linux_cgo.go:207-217`: the loop over `gtk_application_get_windows()` reads
`windows.data` before the nil check. With zero windows the function returns
NULL and the first dereference crashes. Reachable via `App.Hide()`/`Show()`
(hideAllWindows/showAllWindows) and `linuxApp.isVisible()` when called before
any window exists or after all have closed.

**purego fix:** standard nil-checked loop (`linux_purego.go: getWindows`).

## 3. Clipboard sync-read machinery is not reentrant

`linux_cgo.c:924-956`: `clipboard_get_text_sync()` parks the result in two
STATIC globals (`clipboard_sync_result`, `clipboard_sync_done`) while spinning
a nested `g_main_context_iteration` loop. If a second clipboard read starts
while the first is still iterating (the nested loop dispatches arbitrary main
-loop sources — including another dispatched `clipboardGet`), the two calls
overwrite each other's flags: one caller can return the other's text, or spin
forever after its result was consumed.

**purego fix:** per-call state keyed by a handle in a mutex-guarded map
(`linux_purego_callbacks.go: clipboardGetTextSync`).

## 4. `gtkSignalToMenuItem` map is read/written without synchronisation

`linux_cgo.go:80,269-290,405-408`: `attachMenuHandler` writes the map from
whatever goroutine builds the menu; `menuActionActivated` reads it on the GTK
main thread. Unsynchronised concurrent map access is a fatal runtime error
("concurrent map read and map write") — a menu rebuilt via `Menu.Update()`
while the user clicks an item can kill the process.

**purego fix:** RWMutex around all accesses (`linux_purego.go`).

## 5. `menuItemActions` / `menuItemCounters` maps likewise unguarded

`linux_cgo.go:257-267,327-348`: same pattern as #4 — `generateActionName`
writes `menuItemActions` during menu construction, `menuItemSetChecked`/
`menuItemSetDisabled`/`setMenuItemAccelerator` read it from API goroutines
(only `menuItemIds` got a mutex in cgo). Same fatal-crash class.

**purego fix:** one RWMutex (`menuItemsLock`) guards all three maps.

## 6. File-dialog callback can deadlock the GTK main loop on >100 selections

`linux_cgo.go:1767-1794`: `fileDialogCallback` runs on the GTK main thread
(GtkFileDialog async completion) and sends each selected path into a channel
with a fixed buffer of 100 (`runChooserDialog`). Selecting more than 100 files
blocks the main thread on the 101st send until the consumer drains — and if
the consumer is itself waiting on anything main-thread-bound, the app
deadlocks permanently.

**purego fix:** hand the collected paths to a goroutine for delivery, so the
main thread never blocks (`linux_purego.go: fileDialogCallback`).

## 7. System tray silently requires cgo for no reason

`systemtray_linux.go` (shared file) carried a vestigial `import "C"` with zero
`C.` usages. It's a pure-Go dbus StatusNotifier implementation, but the import
made the file cgo-only, which would have silently dropped the tray from any
CGO_ENABLED=0 build.

**Fix (applies to both backends):** removed the import — this one is fixed in
the shared file itself, not just in the purego twin.

## 8. Message dialog: closing via the titlebar leaves a zombie window (and a use-after-free)

`linux_cgo.c:787-793`: `on_message_dialog_close` (the "close-request" handler)
delivers the cancel result, frees the `MessageDialogData`, and returns `TRUE`.
For GTK4's `close-request`, returning TRUE means "handled — do NOT close", so
the dialog window stays on screen forever. Worse, its buttons are still wired
to the now-freed `MessageDialogData`; clicking one after the X-button is a
use-after-free.

**purego fix:** the handler returns FALSE so GTK's default destroys the
window, and the dialog state lives in a Go registry keyed by handle — a late
button click after teardown resolves to nil and is ignored
(`linux_purego_callbacks.go: onMessageDialogClosePtr`).

## Parity observation (crash present in BOTH backends, not fixed here)

On a headless X server without DRI (Xvfb; `/dev/dri/*` inaccessible, WebKit on
the software/EGL-fallback path), closing a second window while the first stays
open kills the process with an X error — `BadDrawable, request_code 14
(X_GetGeometry)` — followed in the cgo build by glibc `free(): corrupted
unsorted chunks`. Verified on Ubuntu 26.04 / GTK 4.22.2 / WebKitGTK 2.52.3
with byte-identical reproduction steps against both the purego and the cgo
binaries: both die the same way, so this is an upstream WebKitGTK teardown
issue in the no-DRI rendering path, not a port defect. Main-window lifecycle
(open → interact → close → clean exit 0) is unaffected. Worth re-testing on a
real GPU-backed session before chasing it in Wails.

## Non-bug deltas (deliberate)

- **Fractional monitor scale on older GTK4:** the cgo build hard-requires
  `gdk_monitor_get_scale` (GTK 4.14+) at link time, so it cannot run against
  older GTK4. The purego backend resolves it optionally and falls back to the
  integer `gdk_monitor_get_scale_factor` at runtime (`monitorScale()`),
  extending the supported range downward instead of crashing at startup.
- **Missing-library UX:** a cgo binary fails at exec time with the dynamic
  linker's terse "cannot open shared object file". The purego backend reports
  every missing library/symbol with per-distro install hints
  (`linux_purego_lib.go: loadLinuxLibraries`).
- **Per-menu-item `MenuItemData` heap blocks leak in cgo** (`linux_cgo.c`,
  `g_new0(MenuItemData,1)` never freed; freed only implicitly at exit). The
  purego port passes the item id as the callback's data word, so the
  allocation doesn't exist. Listed as a delta, not a fix, because the leak is
  bounded by menu size.
