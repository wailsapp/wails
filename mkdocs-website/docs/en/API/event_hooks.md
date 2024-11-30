wails3 provides an event system that allows for hooking into application and window events

```go
// Notification of application start
application.RegisterApplicationEventHook(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application started!")
})
```

```go
// Notification of system theme change
application.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
    app.Logger.Info("System theme changed!")
    if event.Context().IsDarkMode() {
        app.Logger.Info("System is now using dark mode!")
    } else {
        app.Logger.Info("System is now using light mode!")
    }
})
```

```go
// Disable window closing by canceling the event
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window 1 Closing? Nope! Not closing!")
    e.Cancel()
})
```

```go
// Notification of window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("[ApplicationEvent] Window focus!")
})
```

### Application Events

Application events are hookable events that can be registered with `application.RegisterApplicationEventHook()`
and `application.OnApplicationEvent()`. These events are based on `events.ApplicationEventType`.

`events.Common.ApplicationStarted`
:   Triggered when the application starts

`events.Common.ThemeChanged`
:   Triggered when the application theme changes


### Window Events

`events.Common.WindowMaximised` 
:   Triggered when the window is maximised

`events.Common.WindowUnmaximised`
:   Triggered when the window is unmaximised

`events.Common.WindowMinimised` 
:   Triggered when the window is minimised

`events.Common.WindowUnminimised`
:   Triggered when the window is unminimised

`events.Common.WindowFullscreen`
:   Triggered when the window is set to fullscreen

`events.Common.WindowUnfullscreen`
:   Triggered when the window is unfullscreened

`events.Common.WindowRestored`
:   Triggered when the window is restored

`events.Common.WindowClosing`
:   Triggered before the window closes

`events.Common.WindowZoom` 
:   Triggered when the window is zoomed

`events.Common.WindowZoomOut`
:   Triggered when the window is zoomed out

`events.Common.WindowZoomIn` 
:   Triggered when the window is zoomed in

`events.Common.WindowZoomReset`
:   Triggered when the window zoom is reset

`events.Common.WindowFocus` 
:   Triggered when the window gains focus

`events.Common.WindowLostFocus` 
:   Triggered when the window loses focus

`events.Common.WindowShow`
:   Triggered when the window is shown

`events.Common.WindowHide` 
:   Triggered when the window is hidden

`events.Common.WindowDPIChanged`
:   Triggered when the window DPI changes

`events.Common.WindowFilesDropped`
:   Triggered when files are dropped on the window

`events.Common.WindowRuntimeReady`
:   Triggered when the window runtime is ready

`events.Common.WindowDidMove` 
:   Triggered when the window is moved

`events.Common.WindowDidResize` 
:   Triggered when the window is resized

### OS-Specific Events
--8<--
./docs/en/API/events_linux.md
./docs/en/API/events_windows.md
./docs/en/API/events_mac.md
--8<--
