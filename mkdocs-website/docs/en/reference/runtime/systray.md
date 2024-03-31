# System Tray API

The System Tray API provides a way to create and manage system tray icons for your Wails application.

## `newSystemTray(id uint) *SystemTray`

The `newSystemTray()` function creates a new system tray instance.

```go
systemTray := application.newSystemTray(1)
```

### `SetLabel(label string)`

The `SetLabel()` method sets the label for the system tray icon.

```go
systemTray.SetLabel("My App")
```

### `Run()`

The `Run()` method starts the system tray. This happens automatically when the application is run and only needs
to be called if the system tray is dynamically created after the main application is running.

```go
systemTray.Run()
```

### `SetIcon(icon []byte) *SystemTray`

The `SetIcon()` method sets the icon for the system tray. Should be in PNG format.

```go
systemTray.SetIcon(iconData)
```

### `SetDarkModeIcon(icon []byte) *SystemTray`

The `SetDarkModeIcon()` method sets the icon to be used in dark mode. Should be in PNG format.

```go
systemTray.SetDarkModeIcon(darkModeIconData)
```

### `SetMenu(menu *Menu) *SystemTray`

The `SetMenu()` method sets the context menu for the system tray.

```go
menu := application.NewMenu()
systemTray.SetMenu(menu)
```

### `SetIconPosition(iconPosition int) *SystemTray`

The `SetIconPosition()` method sets the position of the icon relative to the label (Mac Only).

```go
systemTray.SetIconPosition(application.NSImageLeading)
```

### `SetTemplateIcon(icon []byte) *SystemTray`

The `SetTemplateIcon()` method sets a template icon that can be used to change the icon color (Mac Only). 

```go
systemTray.SetTemplateIcon(templateIconData)
```

### `Destroy()`

The `Destroy()` method destroys the system tray.

```go
systemTray.Destroy()
```

### `OnClick(handler func()) *SystemTray`

The `OnClick()` method sets the click handler for the system tray icon.

```go
systemTray.OnClick(func() {
    // Handle click
})
```

### `AttachWindow(window *WebviewWindow) *SystemTray`

The `AttachWindow()` method attaches a window to the system tray. The window will be shown centered to the system tray 
when the system tray icon is clicked.

```go
window := application.NewWebviewWindow()
systemTray.AttachWindow(window)
```

