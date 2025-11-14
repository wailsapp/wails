# macOS Issue #4650 Test Application

This test application demonstrates the fixes for [Issue #4650](https://github.com/wailsapp/wails/issues/4650) - macOS-specific window behavior problems.

## Issues Tested

### Issue #1: System Tray Icon Disappearing When Window is Hidden

**Problem:** When calling `Window.Hide()`, the window disappears entirely, including the tray icon.

**Expected Behavior:** The system tray icon should remain visible even when all windows are hidden, especially when using `ActivationPolicyAccessory`.

### Issue #2: White Flicker When Maximizing with Dark Background

**Problem:** When maximizing, the window first expands as a white screen, then fills correctly. With a dark background, this flicker is visually jarring.

**Expected Behavior:** The window should zoom/maximize smoothly maintaining the background color throughout the animation with no white flash.

## How to Run

```bash
cd v3/examples/macos-issue-4650-test
go run .
```

**Note:** This test application is designed for macOS. The issues being tested are macOS-specific.

## Testing Instructions

### Testing Issue #1 (Tray Icon)

1. Launch the application
2. Look at the macOS menu bar - you should see a tray icon
3. Click the "Hide Window" button in the application
4. **✓ PASS:** The tray icon remains visible in the menu bar
5. **✗ FAIL:** The tray icon disappears from the menu bar
6. Click the tray icon to show the window again

### Testing Issue #2 (Maximize Flicker)

1. Ensure the window is not maximized (click "Restore Window" if needed)
2. Watch the window carefully
3. Click the "Maximize Window" button
4. Observe the zoom/maximize animation
5. **✓ PASS:** The animation is smooth with dark background throughout, no white flash
6. **✗ FAIL:** A white flash appears during the zoom animation

## Technical Details

### Fix for Issue #1

**File:** `v3/pkg/application/webview_window_darwin.go`

When hiding a window, the implementation now checks if there are active system trays. If trays exist, it explicitly activates the application to prevent macOS from hiding the entire app (and its tray icons).

```go
func (w *macosWebviewWindow) hide() {
    C.windowHide(w.nsWindow)

    // Check for active system trays
    globalApplication.systemTraysLock.Lock()
    hasSystemTrays := len(globalApplication.systemTrays) > 0
    globalApplication.systemTraysLock.Unlock()

    if hasSystemTrays {
        // Keep app visible for tray
        C.activateIgnoringOtherApps()
    }
}
```

### Fix for Issue #2

**File:** `v3/pkg/application/webview_window_darwin.m`

Changed window initialization from clear/transparent to opaque with a proper background color. This prevents the white flash during animations while still supporting transparent/translucent backdrops when explicitly requested.

```objc
- (WebviewWindow*) initWithContentRect:... {
    self = [super initWithContentRect:...];
    // Initialize with opaque background (not clear)
    [self setBackgroundColor:[NSColor windowBackgroundColor]];
    [self setOpaque:YES];  // Not NO
    return self;
}
```

## System Information

The test was designed based on the original issue report:
- **OS:** macOS Sequoia 15.5
- **Hardware:** Apple M4 (arm64)
- **Wails Version:** v3.0.0-alpha.36+

## Expected Results Summary

| Test | Expected Result |
|------|----------------|
| Hide Window with Tray | Tray icon remains visible in menu bar |
| Show Window from Tray | Window appears and becomes focused |
| Maximize Window | Smooth animation with no white flicker |
| Restore Window | Window returns to previous size smoothly |

## Additional Notes

- The app uses `ActivationPolicyAccessory` to replicate the exact conditions from the original issue
- The window has a dark background (#1e1e1e) to make any white flicker clearly visible
- The tray icon includes both click handlers and a menu for comprehensive testing
