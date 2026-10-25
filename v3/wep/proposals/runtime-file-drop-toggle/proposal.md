# Wails Enhancement Proposal (WEP)

## Toggle native file drop at runtime

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-10-25  
**Discussion**: [#4848](https://github.com/wailsapp/wails/pull/4848) (drag-and-drop overhaul, creation-time model), [#5577](https://github.com/wailsapp/wails/pull/5577) and [#5772](https://github.com/wailsapp/wails/pull/5772) (macOS overlay fixes), [#6283](https://github.com/wailsapp/wails/pull/6283) (open, GTK4 motion controller)  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

Add `SetEnableFileDrop(enabled bool) Window` and `IsFileDropEnabled() bool` to
`WebviewWindow` so an application can turn native file drop on and off while a
window is open; today `WebviewWindowOptions.EnableFileDrop` is fixed at
creation. The setter updates the option, the native drop handling and the
page's `_wails.flags.enableFileDrop`, and the option becomes a Go-side gate on
dropped paths. It is Go-only, not in the JS runtime.

## Motivation

Native file drop and the page's own file handling are exclusive, chosen at
creation. With `EnableFileDrop: true` the app gets paths through
`WindowFilesDropped`, but on macOS and Linux the page never sees HTML5 drag
events with `File` objects. With `false` the page gets the drag (except on
GTK3), but the app never learns the paths.

Many windows need both at different times: a project panel taking a dropped
folder as a path, then an editor or upload area reading files through the File
API or its own library, perhaps in one page swapping views. Today that means
visibly recreating the window, losing page state.

Apps cannot do it themselves: the state lives in the native destination (a
macOS overlay view, GTK controllers or signal handlers), the page flag, and
the option it is re-injected from after each navigation. A flag set from JS is
lost on navigation and leaves the native destination, which on macOS cannot
be overridden from outside, in place.

## Detailed Design

```go
// SetEnableFileDrop turns native file drop on or off while the window is
// open, updating WebviewWindowOptions.EnableFileDrop. Not exposed to JS.
func (w *WebviewWindow) SetEnableFileDrop(enabled bool) Window

// IsFileDropEnabled reports whether native file drop is enabled.
func (w *WebviewWindow) IsFileDropEnabled() bool
```

Names follow the option, `SetAlwaysOnTop`/`SetResizable`/`SetContentProtection`
(returning `Window` for chaining) and `IsFullscreen`/`IsFocused`.

**Behaviour.** The setter writes `w.options.EnableFileDrop`, then, if the
window has a platform implementation and is not destroyed, calls the new
`webviewWindowImpl.setEnableFileDrop(enabled)` via `InvokeSync`; before the
window runs only the option changes. Each implementation is idempotent (no
stacked handlers or double delivery), updates the loaded page's flag, and
shares code with window creation.

**Go gate.** While the option is off, `InitiateFrontendDropProcessing` (macOS,
Linux: native → page) hands no paths to the page and
`handleDragAndDropMessage` (page → Go, all platforms) raises no
`WindowFilesDropped`. So the option, not the page flag, is the authority.

**`Window` interface.** Proposed: add both to `Window`, no-ops on
`BrowserWindow` (server mode), like the other setters. Most window accessors
(`WindowManager.Get`, `GetByName`, `GetByID`, `Current`, `GetAll`,
`OnCreate`, the `KeyBindings` and `RawMessageHandler` callbacks) return
`Window`; the setters it lacks (`PhysicalBounds`/`SetPhysicalBounds`,
`RegisterKeyBinding`, `HandleDrag*`) are not a pattern this fits; and its
unexported methods (`handleDragAndDropMessage`,
`shouldUnconditionallyClose`, `cut`, …) mean nothing outside
`pkg/application` implements it. Three lines, not yet in the reference
implementation.

**Open questions:**

1. Add the methods to `Window`?
2. `SetEnableFileDrop` or `SetFileDropEnabled` (see Alternatives)?
3. Should the runtime stop overriding `dropEffect` while the flag is `false`
   (changes existing behaviour)?

## Non-Goals

- A JS API (see Alternatives).
- Native drop and page file drop at once (per element): on macOS a drag has
  one destination.
- Changing the creation-time option or default, `data-file-drop-target`, or
  the hover classes.
- Making GTK3's "off" hand file drags to the page.
- iOS, Android and server mode.

## Platform Considerations

- **macOS.** "On" adds the transparent `WebviewDrag` overlay above the
  `WKWebView` (autoresizing, as since #5577); "off" calls
  `unregisterDraggedTypes`, so an in-flight drag cannot reach it, then removes
  it. AppKit offers a drag only to the frontmost registered view, so the
  overlay must go, not be disabled. It registers only `NSFilenamesPboardType`,
  so in-page element drags (the sortable list in `v3/examples/drag-n-drop`)
  reach the webview either way. Whether AppKit sends `draggingExited:` to a
  view removed mid-drag is undocumented and unchecked, so a hover class could
  stick until the next drag ends; the setter should call the existing exit
  path first.
- **GTK4.** One `GtkDropControllerMotion` and `GtkDropTarget` per webview,
  held under object-data keys with their own reference, attached for "on",
  detached (`gtk_widget_remove_controller`) for "off" so drags reach WebKit.
  A detached controller emits no `leave`, so disabling mid-hover sends one
  first. The open #6283 touches only the motion callbacks; they merge cleanly.
- **GTK3.** `enableDND`'s four accepting handlers (`drag-data-received`,
  `drag-drop`, `drag-motion`, `drag-leave`) and `disableDND`'s two refusing
  ones (`drag-drop`, `drag-motion`) are connected once and exactly one set is
  unblocked (`g_signal_handler_block`/`unblock`, counts never above one).
  Disabling mid-hover sends a drag-leave first. The existing "off" **refuses**
  external file drags (`gdk_drag_status(…, 0, …)`,
  `gtk_drag_finish(…, FALSE, …)`) so WebKit does not navigate to the file,
  passing internal HTML5 drags; kept, as changing it changes creation-time
  `false` too.
- **Windows.** WebView2 always gives file drags to the page; the runtime sends
  `File` objects from `data-file-drop-target` drops to Go via
  `postMessageWithAdditionalObjects` to get paths. `false` for
  `ICoreWebView2Controller4.AllowExternalDrop` breaks all HTML5 drag-and-drop
  (see the comment in `setupChromium`), so the setter only updates the flag;
  off, the runtime shows "no drop".
- **JS runtime "off" (all desktop).** The runtime's bubble-phase
  `dragenter`/`dragover`/`dragleave`/`drop` listeners on
  `document.documentElement` always `preventDefault()` for `Files` drags
  (stopping navigation to the file) and, while the flag is `false`, set
  `dataTransfer.dropEffect = 'none'`. They run after the page's listeners, so
  a page wanting `File` objects must `stopPropagation()` in its `dragover`.
  This already applies to creation-time `false`; it should be documented.
- **iOS, Android, server.** No-ops; they already ignore the option.

## Pros/Cons

### Pros

- Switches modes without recreating the window, keeping the state in step.
- Go-side gate: a late drop or stale page delivers nothing.
- Additive, idempotent, small: one impl method per platform plus re-callable
  Linux helpers.

### Cons

- "Off" differs by platform: GTK3 refuses file drags; elsewhere the page must
  stop propagation to beat the `dropEffect` override.
- macOS sends no leave on mid-hover removal yet.
- The option field is unsynchronised (see Maintenance Plan).

## Alternatives Considered

- **`SetFileDropEnabled`.** Reads better beside the getter, but setters are
  named after their option, and searching `EnableFileDrop` finds both.
- **JS runtime API.** Rejected: a page that could turn it on grants itself
  paths; JS "off" alone is harmless but splits ownership of the state.
- **Always on, filter in Go; or a disabled flag on the macOS overlay.**
  Rejected: the overlay is the drag's only destination (answering
  `NSDragOperationNone` denies the `WKWebView`), and GTK4's drop target takes
  `GDK_TYPE_FILE_LIST` drags in the capture phase. Filtering picks who gets
  paths, not the drag.
- **WebView2 `AllowExternalDrop`.** Rejected: breaks the HTML5 drag-and-drop
  being switched to.
- **Recreate the window.** Today's workaround: loses state, flickers, moves.
- **Per-element native drop.** The macOS overlay would have to track page
  geometry or move into the `WKWebView`'s drag handling; out of scope.

## Backwards Compatibility

Additive; extending `Window` breaks no one (see above). Apps that never call
the setter see no change: the option keeps its meaning and default, Linux
creation uses the restructured but equivalent `enableDND`/`disableDND`, and
the gate only refuses what the native layer and runtime already refuse.

## Security and Privacy

- Only Go can turn file drop on. While off, the Go gate gives an in-flight
  drop, a stale page or a forged `file:drop` message nothing.
- "Off" does not stop HTML5 `File` objects reaching the page where the
  platform delivers them (macOS, GTK4, Windows): standard browser behaviour,
  exposing contents the user chose to drop, not paths.
- No new permissions, storage or network access.

## Test Plan

**Unit** (`v3/pkg/application/webview_window_filedrop_test.go`, present):
`TestFilesDroppedNotDeliveredWhileFileDropDisabled` and
`TestFilesDroppedFollowsSetEnableFileDrop` (on → one delivery; off → getter
false, none; on → next drop delivered once, right path). To add:
`InitiateFrontendDropProcessing` does not call the platform while off.

**Manual** on macOS, GTK4, GTK3 and Windows, toggling repeatedly:

- *On:* exactly one correct `WindowFilesDropped`; hover class clears.
- *Off:* no event; a drop zone stopping propagation gets `File` objects
  (macOS, GTK4, Windows); GTK3 shows "no drop" without navigating.
- *Either:* in-page `draggable="true"` drags work.
- *Mid-hover:* toggling off on a timer clears the hover class (GTK3, GTK4) and
  releasing delivers nothing.
- *Navigation:* after a reload the flag matches `IsFileDropEnabled()`.
- *Creation:* `false`, then `SetEnableFileDrop(true)`, delivers drops.

**Example:** `v3/examples/drag-n-drop` gains a toggle (bound service method),
state readout, and an HTML5 file drop zone for "off"; its README explains
per-platform "off".

## Reference Implementation

https://github.com/juggler-ai/wails/commit/f07c03b12afbede1f609b0f240dc8b821a19c083
(branch `feat/set-enable-file-drop`, one commit on `upstream/master`
`b6167d9b1`): the API, the Go gate, every platform, the unit tests and a
changelog entry. To follow review: the `Window` entries, the macOS mid-hover
leave, the `InitiateFrontendDropProcessing` test, the example and docs. It
merges with #6283 without conflicts.

## Maintenance Plan

The implementor maintains it with the drag-and-drop code. To watch:

- A new navigation hook or platform must inject `enableFileDrop` from the
  options.
- `enableDND`/`disableDND` changes (#6283 is one) must keep them re-callable.
- `options.EnableFileDrop` is unlocked, like other option fields with setters
  (`AlwaysOnTop`, `DisableResize`, `Frameless`); a racing drop is delivered or
  refused, never twice. It should follow any future option synchronisation
  rather than get a lock of its own.

## Conclusion

A window can switch between native paths and page-handled drops without
recreation; the change is additive and Go-gated, so a page cannot grant itself
paths.
