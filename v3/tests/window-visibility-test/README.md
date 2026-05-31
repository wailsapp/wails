# Window Visibility Test - Issue #2861

This example demonstrates and tests the fixes implemented for [Wails v3 Issue #2861](https://github.com/wailsapp/wails/issues/2861) regarding application windows not showing on Windows 10 Pro due to efficiency mode.

## Problem Background

On Windows systems, the "efficiency mode" feature could prevent Wails applications from displaying windows properly. This occurred because:

1. **WebView2 NavigationCompleted events** could be delayed or missed in efficiency mode
2. **Window visibility was gated** behind WebView2 navigation completion
3. **No fallback mechanisms** existed for delayed or failed navigation events

## Solution Implemented

The fix implements a **robust cross-platform window visibility pattern**:

### Windows Improvements
- ✅ **Decouple window container from WebView state** - Windows show immediately
- ✅ **3-second timeout fallback** - Shows WebView if navigation is delayed
- ✅ **Efficiency mode prevention** - Sets WebView2 `IsVisible=true` per Microsoft guidance
- ✅ **Enhanced state tracking** - Proper visibility state management

### Cross-Platform Consistency  
- ✅ **macOS** - Already robust, documented best practices
- ✅ **Linux** - Added missing show/hide methods for both CGO and purego builds

## Test Scenarios

This example provides comprehensive testing for:

### 1. **Basic Window Tests**
- **Normal Window**: Standard window creation - should appear immediately
- **Delayed Content Window**: Simulates heavy content loading (like Vue.js apps)
- **Hidden → Show Test**: Tests delayed showing after initial creation

### 2. **Stress Tests**
- **Multiple Windows**: Creates 3 windows simultaneously
- **Rapid Creation**: Creates windows in quick succession

### 3. **Critical Issue #2861 Test**
- **Efficiency Mode Test**: Specifically designed to reproduce and verify the fix
- Tests window container vs content loading timing
- Includes heavy content simulation

## How to Run

```bash
cd /path/to/wails/v3/examples/window-visibility-test
wails dev
```

## Testing Instructions

### What to Look For
1. **Immediate Window Appearance** - Windows should appear within 100ms of clicking buttons
2. **Progressive Loading** - Content may load progressively, but window container visible immediately  
3. **No Efficiency Mode Issues** - Windows appear even if Task Manager shows "efficiency mode"
4. **Consistent Cross-Platform Behavior** - Similar behavior on Windows, macOS, and Linux

### How to Test
1. **Note the current time** displayed in the app
2. **Click any test button** or use menu items
3. **Immediately observe** if a window appears (should be within 100ms)
4. **Wait for content** to load and check reported timing
5. **Try multiple tests** in sequence to test robustness
6. **Test both buttons and menu items** for comprehensive coverage

### Expected Results
- ✅ Window containers appear immediately upon button click
- ✅ Content loads progressively within 2-3 seconds
- ✅ No blank or invisible windows, even under efficiency mode
- ✅ Activity log shows sub-100ms window creation times
- ✅ All test scenarios work consistently

## Manual Testing Checklist

### Windows 10 Pro (Primary Target)
- [ ] Test with efficiency mode enabled in Task Manager
- [ ] Create windows while system is under load
- [ ] Test rapid window creation scenarios
- [ ] Verify WebView2 content loads after container appears
- [ ] Check activity log for sub-100ms creation times

### Windows 11
- [ ] Verify consistent behavior with Windows 10 Pro fixes
- [ ] Test efficiency mode scenarios
- [ ] Validate timeout fallback mechanisms

### macOS
- [ ] Confirm existing robust behavior maintained
- [ ] Test all window creation scenarios
- [ ] Verify no regressions introduced

### Linux
- [ ] Test both CGO and purego builds
- [ ] Verify new show/hide methods work correctly
- [ ] Test window positioning and timing

## Technical Implementation Details

### Window Creation Flow
```
1. User clicks button → JavaScript calls Go backend
2. Go creates WebviewWindow → Sets properties
3. Go calls window.Show() → IMMEDIATE window container display
4. WebView2 starts navigation → Progressive content loading
5. Timeout fallback ensures WebView shows even if navigation delayed
```

### Key Code Changes
- **Windows**: `/v3/pkg/application/webview_window_windows.go`
- **macOS**: `/v3/pkg/application/webview_window_darwin.go`  
- **Linux**: `/v3/pkg/application/webview_window_linux.go`, `linux_cgo.go`, `linux_purego.go`

## Reporting Test Results

When testing, please report:

1. **Platform & OS Version** (e.g., "Windows 10 Pro 21H2", "macOS 13.1", "Ubuntu 22.04")
2. **Window Creation Timing** (from activity log)
3. **Any Delayed or Missing Windows**
4. **Efficiency Mode Status** (Windows only - check Task Manager)
5. **Content Loading Behavior** (immediate container vs progressive content)
6. **Any Error Messages** in activity log or console

### Sample Test Report Format
```
Platform: Windows 10 Pro 21H2
Efficiency Mode: Enabled
Results:
- Normal Window: ✅ Appeared immediately (<50ms)
- Delayed Content: ✅ Container immediate, content loaded in 2.1s
- Multiple Windows: ✅ All 3 appeared simultaneously
- Critical Test: ✅ Window appeared immediately, content progressive
Notes: No issues observed, all windows visible immediately
```

## Architecture Notes

This example demonstrates the **preferred window visibility pattern** for web-based desktop applications:

1. **Separate Concerns**: Window container vs web content readiness
2. **Immediate Feedback**: Users see window immediately 
3. **Progressive Enhancement**: Content loads and appears when ready
4. **Robust Fallbacks**: Multiple strategies for edge cases
5. **Cross-Platform Consistency**: Same behavior on all platforms
