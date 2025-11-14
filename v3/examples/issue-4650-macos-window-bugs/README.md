# Issue #4650: macOS Window Behavior Bugs

This example demonstrates two macOS-specific window behavior issues in Wails v3 alpha that don't occur on Windows.

## Issues

### Issue 1: Window.hide() Removes Tray Icon

**Expected Behavior (Windows):**
- When `Window.Hide()` is called, the window hides
- The tray icon remains visible and accessible
- User can restore the window by clicking the tray icon

**Actual Behavior (macOS):**
- When `Window.Hide()` is called, the window disappears
- **The tray icon also disappears**, making the application inaccessible
- This may be related to `NSApplicationActivationPolicy` (specifically `ActivationPolicyAccessory` behavior on macOS Sequoia)

### Issue 2: ToggleMaximise White Flicker

**Expected Behavior:**
- Window smoothly transitions to maximized state
- Background color remains consistent during animation

**Actual Behavior (macOS with Frameless=true):**
- Window first expands showing a **white screen/background**
- Then fills with the correct dark background content
- Very noticeable with dark backgrounds
- This is due to the default white NSWindow background showing during the zoom animation

## Technical Details

### Environment
- **OS:** macOS Sequoia 15.5 (24F74)
- **Hardware:** Apple M4 (ARM64)
- **Wails Version:** v3.0.0-alpha.36+
- **Go Version:** 1.21+
- **Configuration:** `Frameless = true`, `ActivationPolicyAccessory`

### Root Causes

1. **Tray Icon Issue**: Related to NSApplicationActivationPolicy handling on macOS
2. **Maximize Flicker**: NSWindow's default white background is visible during zoom animation before the window content is rendered

## How to Run This Example

### Build and Run

```bash
cd v3/examples/issue-4650-macos-window-bugs
go run .
```

Or use the task runner from the v3 directory:

```bash
# From v3 directory
task example NAME=issue-4650-macos-window-bugs
```

### Testing Issue #1 (Window.hide())

1. Launch the application
2. Observe the tray icon in the system tray
3. Click the "Hide Window" button
4. **Bug**: Both the window AND tray icon disappear
5. **Expected**: Only the window should hide, tray icon should remain

### Testing Issue #2 (ToggleMaximise Flicker)

1. Launch the application (dark background should be visible)
2. Click the "Toggle Maximize" button
3. **Bug**: White flash appears before dark content fills the window
4. **Expected**: Dark background should be maintained throughout the animation

## Files in This Example

- `main.go` - Main application with both systray and frameless window configuration
- `assets/index.html` - Interactive UI to test both issues
- `README.md` - This file

## Related Information

- **Issue:** https://github.com/wailsapp/wails/issues/4650
- **Reported:** October 17, 2025
- **Investigation commits:** 7b9cfa0, 9d6e894

## Additional Notes

- Both issues are macOS-specific and do not occur on Windows
- The flicker is particularly noticeable with dark UI themes
- The tray icon issue makes the application completely inaccessible once hidden
- Using `ActivationPolicyAccessory` is required for tray-only applications but seems to trigger the hide/show bug

## Platform Comparison

| Feature | macOS (Sequoia 15.5) | Windows | Expected |
|---------|---------------------|---------|----------|
| Window.hide() with tray | Tray disappears ❌ | Tray remains ✅ | Tray remains |
| ToggleMaximise (frameless) | White flicker ❌ | Smooth ✅ | Smooth |
| Tray menu accessibility | Lost after hide ❌ | Always accessible ✅ | Always accessible |
| Background color consistency | Flickers white ❌ | Consistent ✅ | Consistent |
