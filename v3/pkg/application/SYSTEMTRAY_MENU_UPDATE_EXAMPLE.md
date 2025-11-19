# System Tray Menu Real-Time Update Example

This example demonstrates how to update system tray menu items in real-time on macOS.

## Issue Fixed

Previously, calling `menu.Update()` would not refresh the menu UI while it was open on macOS.
Updates would only become visible after closing and reopening the menu.

## Solution

The fix implements NSMenuDelegate support, which allows macOS to automatically refresh
the menu when it's opened and when updates are needed.

## Example Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "System Tray Menu Update Example",
    })

    // Create a system tray
    systemTray := app.NewSystemTray()

    // Create a menu
    menu := application.NewMenu()
    statusItem := menu.Add("Loading...")

    // Set the menu on the system tray
    systemTray.SetMenu(menu)
    systemTray.Run()

    // Update the menu label in real-time
    go func() {
        counter := 0
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            counter++
            // Update the menu item label
            statusItem.SetLabel(fmt.Sprintf("Loading %d", counter))

            // Trigger menu update - this will now work in real-time!
            menu.Update()
        }
    }()

    app.Run()
}
```

## How It Works

1. **NSMenuDelegate Setup**: When a menu is attached to a system tray, a MenuDelegate
   is created and attached to the NSMenu object.

2. **Automatic Updates**: The NSMenuDelegate's callback methods are called by macOS:
   - `menuWillOpen:` - Called when the menu is about to be displayed
   - `menuNeedsUpdate:` - Called when macOS determines the menu needs refreshing

3. **Callback to Go**: When these delegate methods are triggered:
   - They call the exported Go function `systrayMenuNeedsUpdate(trayID)`
   - This function calls `menu.Update()` to rebuild the menu structure
   - The NSMenu object is updated in-place by clearing and re-adding items

4. **Real-Time Display**: After the menu is rebuilt:
   - macOS detects the changes to the NSMenu object
   - The displayed menu reflects the latest state
   - This works even while the menu is already open

## Technical Details

### Files Modified

- `v3/pkg/application/systemtray_darwin.h` - Added MenuDelegate interface
- `v3/pkg/application/systemtray_darwin.m` - Implemented MenuDelegate class
- `v3/pkg/application/systemtray_darwin.go` - Added delegate management

### Key Changes

1. **MenuDelegate Class** (Objective-C):
   ```objective-c
   @interface MenuDelegate : NSObject <NSMenuDelegate>
   @property (nonatomic) void* menuPtr;
   @property (nonatomic) long trayID;
   @end
   ```

2. **Delegate Implementation**:
   - `menuNeedsUpdate:` - Called when menu needs refreshing, triggers `systrayMenuNeedsUpdate()`
   - `menuWillOpen:` - Called when menu is about to open, triggers `systrayMenuNeedsUpdate()`

3. **Go Integration**:
   - Added exported function `systrayMenuNeedsUpdate()` that rebuilds the menu
   - Added `nsMenuDelegate` field to `macosSystemTray` struct
   - Updated `setMenu()` to create and attach delegate
   - Updated `run()` to set up delegate during initialization

## Platform Notes

- **macOS**: Real-time updates now work with NSMenuDelegate
- **Linux**: Already supported via DBus signals (PR #4615)
- **Windows**: Menu updates work via native Win32 APIs

## Related

- Issue: #4630
- PR: #4615 (Linux fix for reference)
