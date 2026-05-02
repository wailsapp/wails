# Testing Guide - Window Visibility Issue #2861

## Quick Start

1. **Build and run the application:**
   ```bash
   cd v3/examples/window-visibility-test
   ./build.sh
   # OR
   wails dev
   ```

2. **Main testing interface:**
   - The app opens with a comprehensive testing dashboard
   - Contains multiple test scenarios accessible via buttons
   - Also provides menu-based testing (File, Tests, Help menus)
   - Real-time activity logging with precise timing

## Critical Test Cases

### 🎯 **Issue #2861 Reproduction Test** (Most Important)
**Button:** "Efficiency Mode Test"
**Expected:** Window container appears immediately, content loads progressively
**Watch for:** 
- Window visible within 100ms of button click
- Content loading message appears initially
- Content completes loading after 2-3 seconds
- No blank or invisible windows

### ⏳ **Delayed Content Simulation**
**Button:** "Create Delayed Content Window"  
**Expected:** Tests navigation completion timing
**Watch for:**
- Window container appears immediately
- Loading spinner visible initially
- Content loads after 3-second delay
- Window remains visible throughout

### 🔄 **Hidden → Show Robustness**
**Button:** "Hidden → Show Test"
**Expected:** Tests delayed show() calls
**Watch for:**
- Initial response in activity log
- Window appears after exactly 2 seconds
- No timing issues or failures

## Platform-Specific Testing

### Windows 10 Pro (Primary Target)
**Enable Efficiency Mode Testing:**
1. Open Task Manager → Processes tab
2. Find the test application process
3. Right-click → "Efficiency mode" (if available)
4. Run all test scenarios
5. Verify windows still appear immediately

**Key Metrics:**
- Window creation: < 100ms
- Content loading: 2-3 seconds
- No invisible windows under any conditions

### Windows 11
**Similar to Windows 10 Pro but also test:**
- New Windows 11 efficiency features
- Multiple monitor scenarios
- High DPI scaling

### macOS
**Focus on consistency:**
- All scenarios should work identical to Windows
- No regressions in existing robust behavior
- Test across different macOS versions if possible

### Linux
**Test both build variants:**
```bash
# CGO build (default)
wails dev

# Purego build  
CGO_ENABLED=0 wails dev
```
- Verify both variants behave identically
- Test across different Linux distributions

## Success Criteria

### ✅ **Pass Conditions**
- All windows appear within 100ms of button click
- Activity log shows consistent sub-100ms timing
- Content loads progressively without blocking window visibility
- No blank, invisible, or delayed windows under any test scenario
- Efficiency mode (Windows) does not prevent window appearance
- Menu and button testing yield identical results

### ❌ **Fail Conditions**  
- Any window takes >200ms to appear
- Blank or invisible windows under any condition
- Window visibility blocked by content loading
- Efficiency mode prevents window appearance
- Inconsistent behavior between test methods
- Platform-specific failures

## Reporting Results

**Please provide this information:**

```
Platform: [Windows 10 Pro/Windows 11/macOS/Linux distro + version]
Build Type: [CGO/Purego] (Linux only)
Efficiency Mode: [Enabled/Disabled/N/A] (Windows only)

Test Results:
- Normal Window: [✅ Pass / ❌ Fail] - [timing in ms]
- Delayed Content: [✅ Pass / ❌ Fail] - [container timing / content timing]  
- Hidden→Show: [✅ Pass / ❌ Fail] - [notes]
- Multiple Windows: [✅ Pass / ❌ Fail] - [notes]
- Efficiency Mode Test: [✅ Pass / ❌ Fail] - [critical timing results]

Notes:
[Any additional observations, error messages, or unexpected behavior]
```

## Advanced Testing Scenarios

### **Rapid Stress Testing**
1. Click "Rapid Creation Test" multiple times quickly
2. Use keyboard shortcuts to rapidly access menu items
3. Create multiple windows then close them rapidly
4. Test system under load (other applications running)

### **Edge Case Testing**
1. Test during system startup (high load)
2. Test with multiple monitors
3. Test with different DPI scaling settings
4. Test while other WebView2 applications are running

### **Timing Verification**
1. Use browser dev tools (F12) to check console timing
2. Compare activity log timing with system clock
3. Test on slower/older hardware if available
4. Verify timing consistency across multiple runs

## Troubleshooting

### **Common Issues**
- **Blank window**: Check activity log for error messages
- **Slow timing**: Verify system isn't under heavy load
- **Build failures**: Ensure Wails v3 CLI is latest version
- **Import errors**: Run `go mod tidy` in example directory

### **Debug Information**
The application provides extensive logging:
- Browser console (F12) shows JavaScript timing
- Activity log shows backend call timing  
- Go application logs show window creation details
- Check system Task Manager for process efficiency mode status

This comprehensive testing should validate that the window visibility fixes successfully resolve issue #2861 across all supported platforms.