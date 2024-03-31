# Window

To create a window, use
[Application.NewWebviewWindow](application.md#newwebviewwindow) or
[Application.NewWebviewWindowWithOptions](application.md#newwebviewwindowwithoptions).
The former creates a window with default options, while the latter allows you to
specify custom options.

These methods are callable on the returned WebviewWindow object:

### SetTitle

API: `SetTitle(title string) *WebviewWindow`

This method updates the window title to the provided string. It returns the
WebviewWindow object, allowing for method chaining.

### Name

API: `Name() string`

This function returns the name of the WebviewWindow.

### SetSize

API: `SetSize(width, height int) *WebviewWindow`

This method sets the size of the WebviewWindow to the provided width and height
parameters. If the dimensions provided exceed the constraints, they are adjusted
appropriately.

### SetAlwaysOnTop

API: `SetAlwaysOnTop(b bool) *WebviewWindow`

This function sets the window to stay on top based on the boolean flag provided.

### Show

API: `Show() *WebviewWindow`

`Show` method is used to make the window visible. If the window is not running,
it first invokes the `run` method to start the window and then makes it visible.

### Hide

API: `Hide() *WebviewWindow`

`Hide` method is used to hide the window. It sets the hidden status of the
window to true and emits the window hide event.

### SetURL

API: `SetURL(s string) *WebviewWindow`

`SetURL` method is used to set the URL of the window to the given URL string.

### SetZoom

API: `SetZoom(magnification float64) *WebviewWindow`

`SetZoom` method sets the zoom level of the window content to the provided
magnification level.

### GetZoom

API: `GetZoom() float64`

`GetZoom` function returns the current zoom level of the window content.

### GetScreen

API: `GetScreen() (*Screen, error)`

`GetScreen` method returns the screen on which the window is displayed.

### SetFrameless

API: `SetFrameless(frameless bool) *WebviewWindow`

This function is used to remove the window frame and title bar. It toggles the
framelessness of the window according to the boolean value provided (true for
frameless, false for framed).

### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

This function is used to register a context menu and assigns it the given name.

### NativeWindowHandle

API: `NativeWindowHandle() (uintptr, error)`

This function is used to fetch the platform native window handle for the window.

### Focus

API: `Focus()`

This function is used to focus the window.

### SetEnabled

API: `SetEnabled(enabled bool)`

This function is used to enable/disable the window based on the provided boolean
value.

### SetAbsolutePosition

API: `SetAbsolutePosition(x int, y int)`

This function sets the absolute position of the window in the screen.
