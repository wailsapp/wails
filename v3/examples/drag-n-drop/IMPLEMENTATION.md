# Drag-and-Drop Implementation Details

This document explains how file drag-and-drop works across platforms in Wails v3. This is intended for developers working on the Wails codebase, not end users.

## Architecture Overview

File drag-and-drop in Wails v3 uses a **JavaScript-first approach** on all platforms. The native layer intercepts OS drag events, but the actual drop handling and DOM interaction happens in JavaScript. This ensures consistent behavior and proper coordinate handling across platforms.

### Flow

1. User drags files from OS (file manager, desktop) over the Wails window
2. Native layer detects the drag and notifies JavaScript for hover effects
3. User drops files
4. Native layer sends file paths + coordinates to JavaScript
5. JavaScript finds the drop target element (`data-file-drop-target`)
6. JavaScript sends file paths + element details to Go backend
7. Go emits `WindowFilesDropped` event with full context

## Platform Implementations

### Windows (`webview_window_windows.go`)

Windows uses WebView2's built-in file drop support via `chrome.webview.postMessageWithAdditionalObjects`.

**Setup:**
No special setup is needed. WebView2's `AllowExternalDrop` is enabled by default, which is what we want. The JavaScript runtime handles file drops natively.

**How it works:**
1. User drags files from the OS into the WebView2 window
2. JavaScript `dragenter`/`dragover`/`drop` events fire (WebView2 allows external drops by default)
3. JavaScript calls `event.preventDefault()` to stop the browser from navigating to the file
4. JavaScript collects `File` objects and calls `chrome.webview.postMessageWithAdditionalObjects`
5. WebView2 resolves `File` objects to actual file paths
6. Go receives paths via `processMessageWithAdditionalObjects` handler

**Key files:**
- `v3/pkg/application/webview_window_windows.go` - Setup code
- `v3/internal/runtime/desktop/@wailsio/runtime/src/window.ts` - JS handling

**Coordinate handling:**
- WebView2 provides coordinates in CSS pixels (no DPI conversion needed)
- Drop coordinates come from JavaScript `drop` event

### macOS (`webview_window_darwin.go`)

macOS uses native `NSWindow` drag-and-drop with Objective-C.

**Setup (in Objective-C):**
```objc
[window registerForDraggedTypes:@[NSPasteboardTypeFileURL]];
```

**How it works:**
1. `draggingEntered:` - Called when drag enters window, returns `NSDragOperationCopy`
2. `draggingUpdated:` - Called during drag movement, notifies JS for hover effects
3. `draggingExited:` - Called when drag leaves window
4. `performDragOperation:` - Called on drop, extracts file URLs and sends to JS

**Key files:**
- `v3/pkg/application/webview_window_darwin.go` - Go bindings
- Objective-C code embedded in cgo comments

**Coordinate handling:**
- macOS provides coordinates in window-relative points
- Need to convert to webview-relative coordinates
- May need to account for title bar height

**JavaScript notification:**
```objc
// On drag enter, call JS to show drag entered state
macosOnDragEnter(self.windowId);  // Calls window._wails.handleDragEnter()

// On drag over, notify position for hover effects  
macosOnDragOver(self.windowId, x, y);  // Calls window._wails.handleDragOver(x,y)

// On drag exit, clean up state
macosOnDragExit(self.windowId);  // Calls window._wails.handleDragLeave()

// On drop, send files to JS for processing
processDragItems(self.windowId, cArray, (int)count, x, y);
```

**Performance Optimizations (macOS):**

1. **Zero-allocation drag updates**: Pre-allocated buffer for JS strings
2. **Window caching**: Cache window implementation to avoid map lookups
3. **5-pixel threshold**: Only send updates if cursor moves >5 pixels
4. **50ms debounce**: Limits updates to max 20/sec while maintaining responsiveness
5. **Main thread handling**: Timer callbacks use InvokeSync for UI updates

### Linux (`linux_cgo.go`)

Linux uses GTK3 drag-and-drop signals. This is the most complex implementation because GTK intercepts drag events before WebKit sees them.

**Key insight:** GTK and WebKit both want to handle drag events. We must:
1. Handle external file drags at the GTK level
2. Let WebKit handle internal HTML5 drags
3. Distinguish between them using drag target types

**Setup:**
```c
static void enableDND(GtkWidget *widget, gpointer data) {
    g_signal_connect(widget, "drag-data-received", G_CALLBACK(on_drag_data_received), data);
    g_signal_connect(widget, "drag-drop", G_CALLBACK(on_drag_drop), data);
    g_signal_connect(widget, "drag-motion", G_CALLBACK(on_drag_motion), data);
    g_signal_connect(widget, "drag-leave", G_CALLBACK(on_drag_leave), data);
}
```

**How to distinguish file drags from HTML5 drags:**
```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

**Signal handlers:**

`on_drag_drop`:
```c
static gboolean on_drag_drop(GtkWidget *widget, GdkDragContext *context, 
                             gint x, gint y, guint time, gpointer data) {
    if (!is_file_drag(context)) {
        return FALSE;  // Let WebKit handle internal HTML5 drags
    }
    // Request file data
    GdkAtom target = gdk_atom_intern("text/uri-list", FALSE);
    gtk_drag_get_data(widget, context, target, time);
    return TRUE;  // We're handling this
}
```

`on_drag_data_received`:
```c
static void on_drag_data_received(GtkWidget *widget, GdkDragContext *context,
                                  gint x, gint y, GtkSelectionData *data,
                                  guint target_type, guint time, gpointer user_data) {
    // target_type 2 = text/uri-list (file drop)
    // Other types are internal WebKit drags
    if (target_type != 2) {
        return;  // Don't interfere with internal drags
    }
    
    // Parse URIs and send to Go
    gchar **uris = gtk_selection_data_get_uris(data);
    // ... convert to file paths and call Go
    
    gtk_drag_finish(context, TRUE, FALSE, time);
}
```

`on_drag_motion` (for hover effects):
```c
static gboolean on_drag_motion(GtkWidget *widget, GdkDragContext *context,
                               gint x, gint y, guint time, gpointer data) {
    if (!is_file_drag(context)) {
        return FALSE;  // Don't interfere with internal drags
    }
    
    // Notify JavaScript for hover effects
    // Uses execJSDragOver() which writes to a preallocated buffer
    
    gdk_drag_status(context, GDK_ACTION_COPY, time);
    return TRUE;
}
```

**Key files:**
- `v3/pkg/application/linux_cgo.go` - All C code in cgo preamble

**Coordinate handling:**
- GTK provides coordinates relative to the widget (webview)
- No conversion needed - coordinates are already in the right space

**Important notes:**
- Return `FALSE` from handlers to let WebKit process internal drags
- Return `TRUE` and call `gtk_drag_finish()` for file drops
- The `target_type` check in `on_drag_data_received` is crucial

## JavaScript Runtime (`window.ts`)

The JavaScript runtime handles the frontend side of drag-and-drop.

**Key constants:**
```typescript
const DROP_TARGET_ATTRIBUTE = 'data-file-drop-target';
const DROP_TARGET_ACTIVE_CLASS = 'file-drop-target-active';
```

**Hover effect handlers (called from native code on Linux/macOS):**
```typescript
function handleDragEnter(): void {
    nativeDragActive = true;  // Renamed from linuxDragActive for clarity
}

function handleDragOver(x: number, y: number): void {
    const targetElement = document.elementFromPoint(x, y);
    const dropTarget = getDropTargetElement(targetElement);
    
    // Remove class from previous target
    if (currentDropTarget && currentDropTarget !== dropTarget) {
        currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
    }
    
    // Add class to new target
    if (dropTarget) {
        dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
        currentDropTarget = dropTarget;
    }
}

function handleDragLeave(): void {
    if (currentDropTarget) {
        currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
        currentDropTarget = null;
    }
    nativeDragActive = false;
}
```

**Drop handling (called from native code):**
```typescript
HandlePlatformFileDrop(filenames: string[], x: number, y: number): void {
    const element = document.elementFromPoint(x, y);
    const dropTarget = getDropTargetElement(element);

    if (!dropTarget) {
        return;  // Drop outside valid target - ignore
    }

    const elementDetails = {
        id: dropTarget.id,
        classList: Array.from(dropTarget.classList),
        attributes: { /* ... */ },
    };

    // Send to Go backend
    this[callerSym](FilesDropped, {
        filenames,
        x,
        y,
        elementDetails,
    });
}
```

**Finding drop targets:**
```typescript
function getDropTargetElement(element: Element | null): Element | null {
    if (!element) return null;
    return element.closest(`[${DROP_TARGET_ATTRIBUTE}]`);
}
```

## Go Backend

**Window event context (`application.go`):**
```go
type DropTargetDetails struct {
    X          int               `json:"x"`
    Y          int               `json:"y"`
    ElementID  string            `json:"id"`
    ClassList  []string          `json:"classList"`
    Attributes map[string]string `json:"attributes,omitempty"`
}
```

**Event emission:**
```go
// In processFileDrop or equivalent
ctx := newWindowEventContext()
ctx.setDroppedFiles(files)
ctx.setDropTargetDetails(details)
window.emit(events.Common.WindowFilesDropped, ctx)
```

## Debugging Tips

### Linux
Add `printf` statements in C code (remember to `fflush(stdout)`):
```c
printf("DND: drag-drop, is_file_drag=%d\n", is_file_drag(context));
fflush(stdout);
```

### Windows
Use `globalApplication.debug()`:
```go
globalApplication.debug("[DragDrop] Received files", "count", len(files))
```

### JavaScript
Check browser console. Enable debug mode for verbose logging.

### Common issues

1. **Internal HTML5 drag not working**: Native handler is intercepting it. Make sure to return `FALSE`/`false` for non-file drags.

2. **Hover effects not showing**: JavaScript handlers not being called. Check native code is calling `execJS` or evaluating JS.

3. **Wrong coordinates**: Check coordinate space conversions. CSS pixels vs physical pixels vs window-relative.

4. **Drop ignored**: Element doesn't have `data-file-drop-target` attribute. The JS code ignores drops outside valid targets.

## Implementing Drag-Over Updates for Windows

When implementing drag-over hover effects for Windows, consider these approaches based on what we learned from macOS/Linux:

### Approach 1: WebView2 Native Events (Recommended)
If WebView2 provides drag-over events for external files:
- Use the native WebView2 drag events if available
- JavaScript already handles the hover effects via standard DOM events
- No additional work needed if WebView2 passes through drag events

### Approach 2: Win32 Drag-Drop with Notifications
If you need to intercept at the Win32 level (like macOS/Linux):

```cpp
// In your IDropTarget implementation
HRESULT DragOver(DWORD grfKeyState, POINTL pt, DWORD* pdwEffect) {
    // Convert screen coordinates to window coordinates
    POINT windowPt = {pt.x, pt.y};
    ScreenToClient(hwnd, &windowPt);
    
    // Notify JavaScript (similar to macOS/Linux)
    // Use a debouncer to limit update frequency
    if (ShouldSendUpdate(windowPt)) {  // 5-pixel threshold + 50ms debounce
        NotifyJavaScript(windowPt.x, windowPt.y);
    }
    
    *pdwEffect = DROPEFFECT_COPY;
    return S_OK;
}
```

### Key Considerations:

1. **Coordinate Systems**: 
   - Win32 uses screen coordinates, convert to window-relative
   - WebView2 might need DPI scaling adjustments
   - Test with different DPI settings

2. **Performance Optimizations**:
   - **5-pixel threshold**: Reduce events by ~90%
   - **50ms debounce timer**: Cap at 20 updates/sec
   - **Pre-allocated buffers**: For JavaScript strings if using ExecuteScript
   - **Caching**: Cache window/WebView2 references

3. **Threading**:
   - Win32 drag callbacks may come from different threads
   - Use PostMessage or similar to marshal to UI thread
   - WebView2 ExecuteScript must be called from UI thread

4. **Distinguishing Drag Types**:
   - Check IDataObject format to distinguish file drags
   - Let WebView2 handle internal HTML5 drags if possible
   - Similar to Linux's target type checking

### Example Debouncer Implementation:
```cpp
class DragDebouncer {
    POINT lastPoint;
    DWORD lastTime;
    UINT_PTR timerId;
    
    bool ShouldSendImmediate(POINT pt) {
        // 5-pixel threshold
        return abs(pt.x - lastPoint.x) >= 5 || 
               abs(pt.y - lastPoint.y) >= 5;
    }
    
    void OnDragOver(POINT pt) {
        if (ShouldSendImmediate(pt)) {
            SendUpdate(pt);
            lastPoint = pt;
            // Start 50ms timer for next update
            SetTimer(hwnd, DRAG_TIMER_ID, 50, nullptr);
        } else {
            // Update pending position for timer
            pendingPoint = pt;
        }
    }
};
```

### Testing Recommendations:
1. Test with high-frequency mouse polling (gaming mice)
2. Verify UI updates without mouse movement after timer fires
3. Test with multiple monitors and DPI settings
4. Ensure no memory leaks in long drag sessions

## Testing

Run the example:
```bash
cd v3/examples/drag-n-drop
go build && ./drag-n-drop
```

Test cases:
1. Drag file from OS → drop on drop zone → should categorize file
2. Drag file from OS → drop outside drop zone → should be ignored
3. Drag task item → drop on priority column → should move (HTML5 drag)
4. Drag file while HTML5 draggable is visible → both should work independently
