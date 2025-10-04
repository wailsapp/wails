# Drag-and-Drop Analysis (Windows Focus)

## Overview
- Wails v3 supports two drag-and-drop surfaces when `EnableDragAndDrop` is set: a native window-level file drop channel and HTML drag/drop within the embedded webview. You can see the end-to-end flow in `v3/pkg/application/webview_window_windows.go` (native bridge) and `v3/internal/runtime/desktop/@wailsio/runtime/src/window.ts` (runtime dispatcher).
- Native drops rely on a custom `IDropTarget` implementation that gathers files/coordinates from Win32, pushes them through the Wails event system, and fan out as `events.Common.WindowDropZoneFilesDropped` with enriched context (`DropZoneDetails`, attribute map, coordinates).
- Pure HTML drag/drop (e.g. `v3/examples/html-dnd-api`) is handled entirely by the browser layer; Wails only needs to stay out of the way so standard DOM APIs operate normally.

## Native/Go Pipeline
1. When a window is created with `EnableDragAndDrop`, line ~1948 of `v3/pkg/application/webview_window_windows.go` instantiates `w32.NewDropTarget()` and (optionally) calls `chromium.AllowExternalDrag(false)` to disable WebView2’s built-in file handling.
2. `EnumChildWindows` registers the COM drop target against every current child HWND. The callbacks (`OnEnter`, `OnOver`, `OnLeave`, `OnDrop`) emit Windows-specific events (`events.Windows.WindowDragEnter` etc.) via `w.parent.emit`, so listeners can react even before files are delivered.
3. `DropTarget.Drop` (see `v3/pkg/w32/idroptarget.go:69-140`) extracts filenames from the `IDataObject`, then hands control back to the window impl. Coordinates arrive as screen pixels (`POINT`), so `OnDrop` converts to window-relative coordinates and then calls `convertWindowToWebviewCoordinates` (lines ~1908-1990). Finally `InitiateFrontendDropProcessing` (in `v3/pkg/application/webview_window.go:1484-1515`) formats a JS call: `window.wails.Window.HandlePlatformFileDrop([...], x, y)`.
4. `MessageProcessor` case `WindowDropZoneDropped` (`v3/pkg/application/messageprocessor_window.go:430-488`) decodes the payload that the runtime posts back, wraps it in `DropZoneDetails`, and pushes it through the buffered `windowDragAndDropBuffer`. `App.handleDragAndDropMessage` picks it up and forwards to `WebviewWindow.HandleDragAndDropMessage`, which attaches dropped files + drop-zone metadata to the `WindowEventContext`.
5. The consumer API is the `events.Common.WindowDropZoneFilesDropped` event. The drag-n-drop example shows how to subscribe (`v3/examples/drag-n-drop/main.go:109-158`) and propagate to the frontend via custom events.

## Runtime/JS Behaviour
- `@wailsio/runtime/src/window.ts:538-680` controls the frontend side. `HandlePlatformFileDrop` locates a drop target using `document.elementFromPoint` and `closest([data-wails-dropzone])`; if nothing qualifies it returns early, so native drops that miss a registered dropzone never reach Go.
- The runtime maintains hover styling by tracking `dragenter/over/leave` on `document.documentElement` and toggling `wails-dropzone-hover`. This is how the example achieves live highlighting.
- A legacy helper still exists: `System.HandlePlatformFileDrop` in `@wailsio/runtime/src/system.ts:159-184` marshals a different payload and calls method id `ApplicationFilesDroppedWithContext`, but there is no matching handler in `messageprocessor_application.go`. That means any codepath that invokes it would receive an HTTP 400/500.
- The window runtime bundles `drag.ts`, which manages `--wails-draggable` regions so window dragging/resizing does not swallow pointer events. Developers must ensure dropzones are not also marked draggable.

## Example Insights
- `v3/examples/drag-n-drop/assets/index.html` annotates folders with `data-wails-dropzone` and demonstrates how attributes flow through to Go (`DropZoneDetails.Attributes`). It also shows that the frontend expects `dropX/dropY` in CSS pixels, which helps when validating coordinate transforms.
- `v3/examples/html-dnd-api` confirms standard HTML DnD works without the native bridge; it is a useful regression test when tweaking `drag.ts` so pointer suppression does not break DOM events.

## Potential Bug Hotspots
- **Window-level drops ignored** - `window.ts:554-569` bails if no dropzone is discovered, so the documented `WindowFilesDropped` event never fires. Users expecting “drop anywhere” support lose file payloads entirely.
- **Coordinate scaling** - `convertWindowToWebviewCoordinates` (`webview_window_windows.go:1908-1994`) computes offsets using physical pixels but never converts to WebView2 DIPs. On mixed-DPI or >100% scaling setups, `elementFromPoint` will query the wrong DOM position.
- **Drop target lifecycle** - `EnumChildWindows` runs only once during initialisation. WebView2 can spawn new `Chrome_RenderWidgetHostHWND` instances on navigation or GPU process resets, leaving them unregistered and breaking drops until the app restarts.
- **OLE cleanup** - `DropTarget.Drop` never calls `w32.DragFinish` after `DragQueryFile`. Windows docs recommend doing so to release HDROP resources; skipping it risks leaks on repeated drops.
- **Stale runtime API** - `System.HandlePlatformFileDrop` references a non-existent backend method (`messageprocessor_application.go` lacks case 100). Any future JS that follows the generated docs will fail at runtime.
- **Backpressure risk** - `windowDragAndDropBuffer` (channel size 5) blocks the HTTP handler if event consumers stall. Heavy processing in listeners could cause the runtime call to hang and, on Windows, freeze the drag cursor until the fetch resolves.

## Improvement Opportunities
1. Emit a fallback `WindowFilesDropped` event directly from `DropTarget.OnDrop` when `HandlePlatformFileDrop` declines the payload, preserving drop-anywhere behaviour.
2. Introduce DPI-aware coordinate conversion (use `globalApplication.Screen.PhysicalToDipPoint`) before invoking `elementFromPoint`.
3. Re-run `RegisterDragDrop` whenever a new WebView child window appears (CoreWebView2 `FrameCreated` / `NewBrowserVersionAvailable` callbacks) and on navigation completions.
4. Either wire up the `ApplicationFilesDroppedWithContext` method or remove the `System.HandlePlatformFileDrop` export to avoid misleading integrators.
5. Add integration tests that exercise drops on high-DPI displays and across multiple monitors using the drag-n-drop example as a harness.
6. Consider surfacing drop-target state (e.g., active element id) via diagnostic logging so Windows reports can be correlated without attaching a debugger.

## Open Questions
- Do we need to respect WebView2’s native `AllowExternalDrop(true)` when `EnableDragAndDrop` is disabled, or should we expose both behaviours concurrently?
- How should conflicting CSS states (`--wails-draggable` vs `data-wails-dropzone`) be resolved? Current logic leaves it up to the developer, but documenting or enforcing precedence could prevent accidental suppression.
- Can we guarantee that `elementFromPoint` is safe when overlays or transparent windows are involved, or do we need hit-testing improvements using `elementsFromPoint`?
- Would it be safer to move drop processing entirely to Go (dispatching both window-wide and targeted events) and keep JS solely for hover styling?
