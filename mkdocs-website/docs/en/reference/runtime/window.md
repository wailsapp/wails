# Window API

The Window API provides a way to create and manage windows in your Wails application.

## `NewWindow(options WebviewWindowOptions) *WebviewWindow`

The `NewWindow()` function creates a new window with the specified options.

```go
window := application.NewWindow(wails.WebviewWindowOptions{
    Width:  800,
    Height: 600,
    URL:    "/",
})
```

## `SetTitle(title string) Window`

The `SetTitle()` method sets the title of the window.

```go
window.SetTitle("My Window")
```

## `SetSize(width, height int) Window`

The `SetSize()` method sets the size of the window.

```go
window.SetSize(1024, 768)
```

## `SetAlwaysOnTop(b bool) Window`

The `SetAlwaysOnTop()` method sets the window to be always on top.

```go
window.SetAlwaysOnTop(true)
```

## `Show() Window`

The `Show()` method displays the window.

```go
window.Show()
```

## `Hide() Window`

The `Hide()` method hides the window.

```go
window.Hide()
```

## `SetURL(s string) Window`

The `SetURL()` method sets the URL to be displayed in the window.

```go
window.SetURL("https://www.example.com")
```

## `SetZoom(magnification float64) Window`

The `SetZoom()` method sets the zoom level of the window.

```go
window.SetZoom(1.5)
```

## `GetZoom() float64`

The `GetZoom()` method returns the current zoom level of the window.

```go
zoom := window.GetZoom()
```

## `SetResizable(b bool) Window`

The `SetResizable()` method sets whether the window is resizable.

```go
window.SetResizable(true)
```

## `Resizable() bool`

The `Resizable()` method returns whether the window is resizable.

```go
isResizable := window.Resizable()
```

## `SetMinSize(minWidth, minHeight int) Window`

The `SetMinSize()` method sets the minimum size of the window.

```go
window.SetMinSize(400, 300)
```

## `SetMaxSize(maxWidth, maxHeight int) Window`

The `SetMaxSize()` method sets the maximum size of the window.

```go
window.SetMaxSize(1920, 1080)
```

## `ExecJS(js string)`

The `ExecJS()` method executes the given JavaScript in the context of the window.

```go
window.ExecJS("console.log('Hello, World!')")
```

## `Fullscreen() Window`

The `Fullscreen()` method sets the window to fullscreen mode.

```go
window.Fullscreen()
```

## `Flash(enabled bool)`

The `Flash()` method flashes the window's taskbar button/icon to indicate that attention is required (Windows only).

```go
window.Flash(true)
```

## `IsMinimised() bool`

The `IsMinimised()` method returns whether the window is minimised.

```go
isMinimised := window.IsMinimised()
```

## `IsVisible() bool`

The `IsVisible()` method returns whether the window is visible.

```go
isVisible := window.IsVisible()
```

## `IsMaximised() bool`

The `IsMaximised()` method returns whether the window is maximised.

```go
isMaximised := window.IsMaximised()
```

## `IsFocused() bool`

The `IsFocused()` method returns whether the window is currently focused.

```go
isFocused := window.IsFocused()
```

## `IsFullscreen() bool`

The `IsFullscreen()` method returns whether the window is in fullscreen mode.

```go
isFullscreen := window.IsFullscreen()
```

## `SetBackgroundColour(colour RGBA) Window`

The `SetBackgroundColour()` method sets the background color of the window.

```go
window.SetBackgroundColour(wails.RGBA{R: 255, G: 255, B: 255})
```

## `Destroy()`

The `Destroy()` method destroys the window.

```go
window.Destroy()
```

## `Reload()`

The `Reload()` method reloads the page assets.

```go
window.Reload()
```

## `ForceReload()`

The `ForceReload()` method forces the window to reload the page assets.

```go
window.ForceReload()
```

## `OpenDevTools()`

The `OpenDevTools()` method opens the developer tools for the window.

```go
window.OpenDevTools()
```

## `ZoomReset() Window`

The `ZoomReset()` method resets the zoom level of the window to 100%.

```go
window.ZoomReset()
```

## `ZoomIn()`

The `ZoomIn()` method increases the zoom level of the window.

```go
window.ZoomIn()
```

## `ZoomOut()`

The `ZoomOut()` method decreases the zoom level of the window.

```go
window.ZoomOut()
```

## `Close()`

The `Close()` method closes the window.

```go
window.Close()
```

## `SetHTML(html string) Window`

The `SetHTML()` method sets the HTML content of the window.

```go
window.SetHTML("<h1>Hello, World!</h1>")
```

## `SetRelativePosition(x, y int) Window`

The `SetRelativePosition()` method sets the position of the window relative to the screen.

```go
window.SetRelativePosition(100, 100)
```

## `Minimise() Window`

The `Minimise()` method minimises the window.

```go
window.Minimise()
```

## `Maximise() Window`

The `Maximise()` method maximises the window.

```go
window.Maximise()
```

## `UnMinimise()`

The `UnMinimise()` method un-minimises the window.

```go
window.UnMinimise()
```

## `UnMaximise()`

The `UnMaximise()` method un-maximises the window.

```go
window.UnMaximise()
```

## `UnFullscreen()`

The `UnFullscreen()` method un-fullscreens the window.

```go
window.UnFullscreen()
```

## `Restore()`

The `Restore()` method restores the window to its previous state (minimised, maximised, or fullscreen).

```go
window.Restore()
```

## `GetScreen() (*Screen, error)`

The `GetScreen()` method returns the screen that the window is on.

```go
screen, err := window.GetScreen()
if err != nil {
    // Handle error
}
```

## `SetFrameless(frameless bool) Window`

The `SetFrameless()` method sets the window to be frameless (without a title bar or window controls).

```go
window.SetFrameless(true)
```

## `On(eventType events.WindowEventType, callback func(event *WindowEvent)) func()`

The `On()` method registers a callback for the specified window event.

```go
window.On(events.Common.WindowFocus, func(event *application.WindowEvent) {
    // Handle window focus event
})
```

## `RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func()`

The `RegisterHook()` method registers a hook for the specified window event. Hooks are called before the event listeners and can cancel the event.

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Handle window closing event
    event.Cancel()
})
```

## `RegisterContextMenu(name string, menu *Menu)`

The `RegisterContextMenu()` method registers a context menu with the given name.

```go
contextMenu := application.NewMenu()
// Add menu items
window.RegisterContextMenu("my-context-menu", contextMenu)
```

## `NativeWindowHandle() (uintptr, error)`

The `NativeWindowHandle()` method returns the platform-specific native window handle for the window.

```go
handle, err := window.NativeWindowHandle()
if err != nil {
    // Handle error
}
```

The Window API provides a powerful and flexible way to create and manage windows in your Wails application.
```