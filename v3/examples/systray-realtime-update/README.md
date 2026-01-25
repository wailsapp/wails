# Systray Realtime Update Demo

This example demonstrates real-time systray menu updates on macOS.

## The Problem (Issue #4719)

Previously, if a background process tried to update the systray menu while the user had the menu open, clicking any menu item would fail with:

```
MenuItem #48 not found
```

This was because `Menu.Update()` would destroy and recreate all menu items with new IDs, but the visible (already-open) menu still referenced the old IDs.

## The Fix

The fix modifies `Menu.Update()` to update existing menu items in-place rather than destroying and recreating them. This preserves the menu item IDs, so clicks continue to work even when the menu is updated while open.

## How to Test

1. Build and run the example:
   ```bash
   go build && ./systray-realtime-update
   ```

2. Click the "RT" systray icon to open the menu

3. **Keep the menu open** and observe:
   - The "Time: HH:MM:SS" item updates every second in real-time
   - The "Status:" item updates based on checkbox states

4. While the menu is still open, click on any item:
   - Counter item should increment and log the click
   - Checkbox items should toggle and log the change
   - No "MenuItem not found" errors should appear

5. Check the console output for logs confirming clicks work correctly.

## What This Example Demonstrates

- **Real-time timestamp updates**: The time is updated every second via a background goroutine
- **Dynamic status based on state**: The status text changes based on checkbox selections
- **Click counter persistence**: Clicking the counter item works correctly even during updates
- **Menu open/close callbacks**: Logs when the menu opens and closes
