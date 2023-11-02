# System Tray

The system tray houses notification area on a desktop environment, which can
contain both icons of currently-running applications and specific system
notifications.

You create a system tray by calling `app.NewSystemTray()`:

```go
    // Create a new system tray
tray := app.NewSystemTray()
```

The following methods are available on the `SystemTray` type:

### SetLabel

API: `SetLabel(label string)`

The `SetLabel` method sets the tray's label.

### Label

API: `Label() string`

The `Label` method retrieves the tray's label.

### PositionWindow

API: `PositionWindow(*WebviewWindow, offset int) error`

The `PositionWindow` method calls both `AttachWindow` and `WindowOffset`
methods.

### SetIcon

API: `SetIcon(icon []byte) *SystemTray`

The `SetIcon` method sets the system tray's icon.

### SetDarkModeIcon

API: `SetDarkModeIcon(icon []byte) *SystemTray`

The `SetDarkModeIcon` method sets the system tray's icon when in dark mode.

### SetMenu

API: `SetMenu(menu *Menu) *SystemTray`

The `SetMenu` method sets the system tray's menu.

### Destroy

API: `Destroy()`

The `Destroy` method destroys the system tray instance.

### OnClick

API: `OnClick(handler func()) *SystemTray`

The `OnClick` method sets the function to execute when the tray icon is clicked.

### OnRightClick

API: `OnRightClick(handler func()) *SystemTray`

The `OnRightClick` method sets the function to execute when right-clicking the
tray icon.

### OnDoubleClick

API: `OnDoubleClick(handler func()) *SystemTray`

The `OnDoubleClick` method sets the function to execute when double-clicking the
tray icon.

### OnRightDoubleClick

API: `OnRightDoubleClick(handler func()) *SystemTray`

The `OnRightDoubleClick` method sets the function to execute when right
double-clicking the tray icon.

### AttachWindow

API: `AttachWindow(window *WebviewWindow) *SystemTray`

The `AttachWindow` method attaches a window to the system tray. The window will
be shown when the system tray icon is clicked.

### WindowOffset

API: `WindowOffset(offset int) *SystemTray`

The `WindowOffset` method sets the gap in pixels between the system tray and the
window.

### WindowDebounce

API: `WindowDebounce(debounce time.Duration) *SystemTray`

The `WindowDebounce` method sets a debounce time. In the context of Windows,
this is used to specify how long to wait before responding to a mouse up event
on the notification icon.

### OpenMenu

API: `OpenMenu()`

The `OpenMenu` method opens the menu associated with the system tray.
