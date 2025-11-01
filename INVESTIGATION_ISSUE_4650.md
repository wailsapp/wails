# Investigation Report: Issue #4650 - macOS Window Behavior Issues

**Issue:** [#4650](https://github.com/wailsapp/wails/issues/4650)
**Platform:** macOS 15.5 (Sequoia) on Apple M4
**Wails Version:** v3.0.0-alpha.36
**Branch:** v3-alpha

## Summary

Two macOS-specific window behavior issues reported:
1. **Window.hide() causing tray icon to disappear** - When calling `Window.hide()`, both the window and tray icon disappear
2. **Window.ToggleMaximise white screen flicker** - When maximizing, the window shows a white flash before content fills

## Investigation Findings

### Issue 1: Window.hide() and Tray Icon Behavior

#### Code Analysis

**Window Hide Implementation:**
- Location: `v3/pkg/application/webview_window_darwin.go:930-933`
- Implementation: Uses Objective-C `[(WebviewWindow*)window orderOut:nil]`
- This is the standard macOS method to hide a window

**System Tray Implementation:**
- Location: `v3/pkg/application/systemtray_darwin.go:54-66`
- Implementation: Uses `[statusItem setVisible:YES/NO]`
- Completely independent from window visibility

**Code Flow:**
```
Window.Hide() (webview_window.go:468)
  ↓
window.impl.hide()
  ↓
C.windowHide(w.nsWindow) (webview_window_darwin.go:932)
  ↓
[(WebviewWindow*)window orderOut:nil] (line 689)
```

**Tray Icon Flow:**
```
SystemTray.Hide() (systemtray.go:249)
  ↓
s.impl.Hide()
  ↓
C.systemTrayHide(s.nsStatusItem) (systemtray_darwin.go:65)
  ↓
[statusItem setVisible:NO] (line 26)
```

#### Analysis

**The window hide and tray icon visibility are completely independent operations.** There is no code that links these two actions. Based on the code review:

1. **No direct coupling:** Window visibility and tray icon visibility use separate APIs with no shared state
2. **No event handlers:** No built-in event listeners that hide the tray when a window is hidden
3. **Application delegate:** The app delegate (`application_darwin_delegate.m`) doesn't contain logic that would hide the tray on window hide

#### Possible Causes

1. **User application code:** The user's application might have event listeners that hide the tray when the window hides
   - Check for: `window.OnWindowEvent(events.Common.WindowHidden, ...)` handlers
   - Check if the attached window feature is being used with custom handlers

2. **macOS Sequoia behavior:** macOS 15.5 might have introduced new behaviors for status bar items
   - If the app's activation policy is set in a certain way
   - If there are no visible windows and the app isn't in the dock

3. **Attached Window Feature:** The `SystemTray.AttachWindow()` feature (systemtray.go:281) sets up automatic window toggling, but this shouldn't hide the tray icon itself

4. **NSApplication activation policy:** If the app is running as an accessory (`NSApplicationActivationPolicyAccessory`), hiding all windows might affect tray visibility

#### Recommended Next Steps

1. **Request minimal reproduction:** Ask the user for a minimal code sample showing the issue
2. **Check application setup:** Review how the app is initialized, particularly:
   - `app.SetActivationPolicy()` calls
   - Event listeners on window hide/show
   - System tray initialization
3. **Test on macOS 15.5:** Verify if this is specific to macOS Sequoia
4. **Add debug logging:** Check if both `window.Hide()` AND `systemTray.Hide()` are being called

---

### Issue 2: Window.ToggleMaximise White Screen Flicker

#### Code Analysis

**Maximize Implementation:**
- Location: `v3/pkg/application/webview_window_darwin.go:959-961`
- Implementation: Uses `[(WebviewWindow*)window zoom:nil]` (line 642)
- This is the standard macOS window zoom method

**Window Creation:**
- Location: `v3/pkg/application/webview_window_darwin.go:27-128`
- WebView config: `config.suppressesIncrementalRendering = true` (line 85)
- No explicit window background color set during creation
- Default NSWindow background is white/system color

**Background Handling:**
- Window background color can be set via `windowSetBackgroundColour()` (line 372)
- Webview background can be set transparent via `webviewSetTransparent()` (line 358)
- Backdrop options available: Normal, Transparent, Translucent, LiquidGlass (lines 1243-1253)

#### Root Cause

The white flicker occurs because:

1. **Default window background:** NSWindow uses a white/system default background unless explicitly set
2. **Zoom animation timing:** During the `zoom:nil` animation:
   - The window frame expands to new size
   - The webview content needs to resize and re-render
   - There's a brief moment where the window is larger than the rendered content
3. **Background visibility:** During this gap, the white NSWindow background shows through
4. **Frameless mode:** The issue is more noticeable in frameless mode where users expect a seamless appearance

#### Technical Details

**Relevant Code Locations:**
```
Window creation: webview_window_darwin.go:27-128
Maximize: webview_window_darwin.go:641-643
  → static void windowMaximise(void *window) {
      [(WebviewWindow*)window zoom:nil];
    }
Background color: webview_window_darwin.go:372-374
  → void windowSetBackgroundColour(void* nsWindow, int r, int g, int b, int alpha)
Backdrop: webview_window_darwin.go:1243-1253
```

**Current Background Initialization:**
- If `Backdrop == MacBackdropTransparent`: Sets window and webview transparent
- If `Backdrop == MacBackdropTranslucent`: Adds visual effect view
- If `Backdrop == MacBackdropNormal`: No special handling (uses system default)
- Background color is set via `setBackgroundColour(options.BackgroundColour)` (line 1241)
- Default RGBA is `{0, 0, 0, 0}` (transparent black)

#### Potential Solutions

**Option 1: Set Window Background to Match Content (Recommended)**
```objc
// In windowNew() or during maximize, ensure window background matches content
// If user has dark content, set to black; if light, set appropriately
// For frameless windows, consider matching the webview background
```

**Option 2: Ensure Webview Background is Always Set**
```objc
// Always set webview background color, not just when transparent
// This ensures content fills the entire window during resize
[window.webView setValue:[NSColor colorWithRed:r/255.0 ...] forKey:@"backgroundColor"];
```

**Option 3: Modify Zoom Animation**
```objc
// Disable or customize the zoom animation to be instant or smoother
[NSAnimationContext beginGrouping];
[[NSAnimationContext currentContext] setDuration:0.0]; // Instant
[(WebviewWindow*)window zoom:nil];
[NSAnimationContext endGrouping];
```

**Option 4: Pre-render During Zoom**
```objc
// In windowWillUseStandardFrame delegate method, prepare the webview
// This could involve forcing a render or adjusting the webview size proactively
```

#### Recommended Fix

The most robust solution is **Option 1 + Option 2**:

1. **Always set the window background color** to match the user's content or default to black for frameless windows
2. **Ensure webview background is set** even when not using transparent backdrop
3. **Consider special handling for frameless windows** where visual polish is critical

**Suggested Implementation:**
```c
// In windowNew(), after creating the window:
if (frameless) {
    // For frameless windows, default to black background to avoid flicker
    [(WebviewWindow*)window setBackgroundColor:[NSColor blackColor]];
}

// Also ensure webview always has a background set
// This can be done in setBackgroundColour() to always apply to both window and webview
```

---

## Additional Observations

### Frameless Mode Considerations
- Both issues are more noticeable in frameless mode where users expect seamless behavior
- Frameless windows use `NSWindowStyleMaskBorderless | NSWindowStyleMaskResizable`
- Corner radius of 8.0 is applied to frameless windows (line 55)

### Testing Requirements
1. Test on macOS Sequoia (15.5) specifically
2. Test with frameless windows
3. Test with different backdrop options
4. Test with various background colors
5. Test system tray with and without attached windows

### Files Modified for Potential Fixes

**For Tray Icon Issue:**
- Likely no code changes needed - requires user code review
- Possibly add better documentation about window/tray independence

**For Maximize Flicker:**
- `v3/pkg/application/webview_window_darwin.go` (C code section)
  - Line 27-128: windowNew() function
  - Line 641-643: windowMaximise() function
  - Line 372-374: windowSetBackgroundColour() function
- `v3/pkg/application/webview_window_darwin.m`
  - Line 211-223: windowDidZoom delegate method (if animation control needed)

---

## Conclusion

### Issue 1: Window.hide() / Tray Icon
**Status:** Needs user reproduction case
**Likely Cause:** User application code or macOS Sequoia behavior
**Recommended Action:** Request minimal reproduction and review user's event handlers

### Issue 2: ToggleMaximise Flicker
**Status:** Reproducible issue with known cause
**Root Cause:** Window background color during zoom animation
**Recommended Action:** Implement window background color management, especially for frameless windows
**Priority:** Medium - affects visual polish but not functionality

---

## References

**Key Files:**
- `/home/user/wails/v3/pkg/application/webview_window_darwin.go` - Main macOS window implementation
- `/home/user/wails/v3/pkg/application/webview_window_darwin.m` - Objective-C delegate implementation
- `/home/user/wails/v3/pkg/application/systemtray_darwin.go` - System tray implementation
- `/home/user/wails/v3/pkg/application/webview_window.go` - Generic window interface
- `/home/user/wails/v3/pkg/application/webview_window_options.go` - Window options and configuration

**Relevant macOS APIs:**
- `orderOut:` - Hide window (NSWindow)
- `zoom:` - Maximize/restore window (NSWindow)
- `setVisible:` - Show/hide status item (NSStatusItem)
- `setBackgroundColor:` - Set window background (NSWindow)
