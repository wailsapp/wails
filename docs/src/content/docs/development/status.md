---
title: Status
sidebar:
  order: 20
---

Status of features in v3.

:::note

This list is a mixture of public and internal API support.

It is not complete and probably not up to date.

:::

### Legend

- ✅ = Supported
- 🚧 = Under Development
- ❔ = Untested
- ❌ = Not available on the platform

## Known Issues

- Linux is not yet up to feature parity with Windows/Mac

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes |
| ------------------------------------------------------------- | ------- | ----- | --- | ----- |
| run() error                                                   | ✅      | ✅    | ✅  |       |
| destroy()                                                     |         | ✅    | ✅  |       |
| setApplicationMenu(menu \*Menu)                               | ✅      | ✅    | ✅  |       |
| name() string                                                 |         | ✅    | ✅  |       |
| getCurrentWindowID() uint                                     | ✅      | ✅    | ✅  |       |
| showAboutDialog(name string, description string, icon []byte) |         | ✅    | ✅  |       |
| setIcon(icon []byte)                                          | ❌      | ✅    | ✅  |       |
| on(id uint)                                                   |         |       | ✅  |       |
| dispatchOnMainThread(fn func())                               | ✅      | ✅    | ✅  |       |
| hide()                                                        | ✅      | ✅    | ✅  |       |
| show()                                                        | ✅      | ✅    | ✅  |       |
| getPrimaryScreen() (\*Screen, error)                          |         | ✅    | ✅  |       |
| getScreens() ([]\*Screen, error)                              |         | ✅    | ✅  |       |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes                                    |
| -------------------------------------------------- | ------- | ----- | --- | ---------------------------------------- |
| center()                                           | ✅      | ✅    | ✅  |                                          |
| close()                                            | ✅      | ✅    | ✅  |                                          |
| destroy()                                          |         | ✅    | ✅  |                                          |
| execJS(js string)                                  | ✅      | ✅    | ✅  |                                          |
| focus()                                            | ✅      | ✅    |     |                                          |
| forceReload()                                      |         | ✅    | ✅  |                                          |
| fullscreen()                                       | ✅      | ✅    | ✅  |                                          |
| getScreen() (\*Screen, error)                      | ✅      | ✅    | ✅  |                                          |
| getZoom() float64                                  |         | ✅    | ✅  |                                          |
| height() int                                       | ✅      | ✅    | ✅  |                                          |
| hide()                                             | ✅      | ✅    | ✅  |                                          |
| isFullscreen() bool                                | ✅      | ✅    | ✅  |                                          |
| isMaximised() bool                                 | ✅      | ✅    | ✅  |                                          |
| isMinimised() bool                                 | ✅      | ✅    | ✅  |                                          |
| maximise()                                         | ✅      | ✅    | ✅  |                                          |
| minimise()                                         | ✅      | ✅    | ✅  |                                          |
| nativeWindowHandle() (uintptr, error)              | ✅      | ✅    | ✅  |                                          |
| on(eventID uint)                                   | ✅      |       | ✅  |                                          |
| openContextMenu(menu *Menu, data *ContextMenuData) | ✅      | ✅    | ✅  |                                          |
| relativePosition() (int, int)                      | ✅      | ✅    | ✅  |                                          |
| reload()                                           | ✅      | ✅    | ✅  |                                          |
| run()                                              | ✅      | ✅    | ✅  |                                          |
| setAlwaysOnTop(alwaysOnTop bool)                   | ✅      | ✅    | ✅  |                                          |
| setBackgroundColour(color RGBA)                    | ✅      | ✅    | ✅  |                                          |
| setEnabled(bool)                                   |         | ✅    | ✅  |                                          |
| setFrameless(bool)                                 |         | ✅    | ✅  |                                          |
| setFullscreenButtonEnabled(enabled bool)           | ❌      | ✅    | ✅  | There is no fullscreen button in Windows |
| setHTML(html string)                               | ✅      | ✅    | ✅  |                                          |
| setMaxSize(width, height int)                      | ✅      | ✅    | ✅  |                                          |
| setMinSize(width, height int)                      | ✅      | ✅    | ✅  |                                          |
| setRelativePosition(x int, y int)                  | ✅      | ✅    | ✅  |                                          |
| setResizable(resizable bool)                       | ✅      | ✅    | ✅  |                                          |
| setSize(width, height int)                         | ✅      | ✅    | ✅  |                                          |
| setTitle(title string)                             | ✅      | ✅    | ✅  |                                          |
| setURL(url string)                                 | ✅      | ✅    | ✅  |                                          |
| setZoom(zoom float64)                              | ✅      | ✅    | ✅  |                                          |
| show()                                             | ✅      | ✅    | ✅  |                                          |
| size() (int, int)                                  | ✅      | ✅    | ✅  |                                          |
| toggleDevTools()                                   | ✅      | ✅    | ✅  |                                          |
| unfullscreen()                                     | ✅      | ✅    | ✅  |                                          |
| unmaximise()                                       | ✅      | ✅    | ✅  |                                          |
| unminimise()                                       | ✅      | ✅    | ✅  |                                          |
| width() int                                        | ✅      | ✅    | ✅  |                                          |
| zoom()                                             |         | ✅    | ✅  |                                          |
| zoomIn()                                           | ✅      | ✅    | ✅  |                                          |
| zoomOut()                                          | ✅      | ✅    | ✅  |                                          |
| zoomReset()                                        | ✅      | ✅    | ✅  |                                          |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
| ------- | ------- | ----- | --- | ----- |
| Quit    | ✅      | ✅    | ✅  |       |
| Hide    | ✅      | ✅    | ✅  |       |
| Show    | ✅      |       | ✅  |       |

### Dialogs

| Feature  | Windows | Linux | Mac | Notes |
| -------- | ------- | ----- | --- | ----- |
| Info     | ✅      | ✅    | ✅  |       |
| Warning  | ✅      | ✅    | ✅  |       |
| Error    | ✅      | ✅    | ✅  |       |
| Question | ✅      | ✅    | ✅  |       |
| OpenFile | ✅      | ✅    | ✅  |       |
| SaveFile | ✅      | ✅    | ✅  |       |

### Clipboard

| Feature | Windows | Linux | Mac | Notes |
| ------- | ------- | ----- | --- | ----- |
| SetText | ✅      | ✅    | ✅  |       |
| Text    | ✅      | ✅    | ✅  |       |

### ContextMenu

| Feature          | Windows | Linux | Mac | Notes |
| ---------------- | ------- | ----- | --- | ----- |
| OpenContextMenu  | ✅      | ✅    | ✅  |       |
| On By Default    |         |       |     |       |
| Control via HTML | ✅      |       |     |       |

The default context menu is enabled by default for all elements that are
`contentEditable: true`, `<input>` or `<textarea>` tags or have the
`--default-contextmenu: true` style set. The `--default-contextmenu: show` style
will always show the context menu The `--default-contextmenu: hide` style will
always hide the context menu

Anything nested under a tag with `--default-contextmenu: hide` style will not
show the context menu unless it is explicitly set with
`--default-contextmenu: show`.

### Screens

| Feature    | Windows | Linux | Mac | Notes |
| ---------- | ------- | ----- | --- | ----- |
| GetAll     | ✅      | ✅    | ✅  |       |
| GetPrimary | ✅      | ✅    | ✅  |       |
| GetCurrent | ✅      | ✅    | ✅  |       |

### System

| Feature    | Windows | Linux | Mac | Notes |
| ---------- | ------- | ----- | --- | ----- |
| IsDarkMode |         |       | ✅  |       |

### Window

| Feature             | Windows | Linux | Mac | Notes                                                                                |
| ------------------- | ------- | ----- | --- | ------------------------------------------------------------------------------------ |
| Center              | ✅      | ✅    | ✅  |                                                                                      |
| Focus               | ✅      | ✅    |     |                                                                                      |
| FullScreen          | ✅      | ✅    | ✅  |                                                                                      |
| GetZoom             | ✅      | ✅    | ✅  | Get current view scale                                                               |
| Height              | ✅      | ✅    | ✅  |                                                                                      |
| Hide                | ✅      | ✅    | ✅  |                                                                                      |
| Maximise            | ✅      | ✅    | ✅  |                                                                                      |
| Minimise            | ✅      | ✅    | ✅  |                                                                                      |
| RelativePosition    | ✅      | ✅    | ✅  |                                                                                      |
| Screen              | ✅      | ✅    | ✅  | Get screen for window                                                                |
| SetAlwaysOnTop      | ✅      | ✅    | ✅  |                                                                                      |
| SetBackgroundColour | ✅      | ✅    | ✅  | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetEnabled          | ✅      | ❔    | ❌  | Set the window to be enabled/disabled                                                |
| SetMaxSize          | ✅      | ✅    | ✅  |                                                                                      |
| SetMinSize          | ✅      | ✅    | ✅  |                                                                                      |
| SetRelativePosition | ✅      | ✅    | ✅  |                                                                                      |
| SetResizable        | ✅      | ✅    | ✅  |                                                                                      |
| SetSize             | ✅      | ✅    | ✅  |                                                                                      |
| SetTitle            | ✅      | ✅    | ✅  |                                                                                      |
| SetZoom             | ✅      | ✅    | ✅  | Set view scale                                                                       |
| Show                | ✅      | ✅    | ✅  |                                                                                      |
| Size                | ✅      | ✅    | ✅  |                                                                                      |
| UnFullscreen        | ✅      | ✅    | ✅  |                                                                                      |
| UnMaximise          | ✅      | ✅    | ✅  |                                                                                      |
| UnMinimise          | ✅      | ✅    | ✅  |                                                                                      |
| Width               | ✅      | ✅    | ✅  |                                                                                      |
| ZoomIn              | ✅      | ✅    | ✅  | Increase view scale                                                                  |
| ZoomOut             | ✅      | ✅    | ✅  | Decrease view scale                                                                  |
| ZoomReset           | ✅      | ✅    | ✅  | Reset view scale                                                                     |

### Window Options

| Feature                         | Windows | Linux | Mac | Notes                                      |
| ------------------------------- | ------- | ----- | --- | ------------------------------------------ |
| AlwaysOnTop                     | ✅      | ✅    |     |                                            |
| BackgroundColour                | ✅      | ✅    |     |                                            |
| BackgroundType                  |         |       |     | Acrylic seems to work but the others don't |
| CSS                             | ✅      | ✅    |     |                                            |
| DevToolsEnabled                 | ✅      | ✅    | ✅  |                                            |
| DisableResize                   | ✅      | ✅    |     |                                            |
| EnableDragAndDrop               |         | ✅    |     |                                            |
| EnableFraudulentWebsiteWarnings |         |       |     |                                            |
| Focused                         | ✅      | ✅    |     |                                            |
| Frameless                       | ✅      | ✅    |     |                                            |
| FullscreenButtonEnabled         | ✅      |       |     |                                            |
| Height                          | ✅      | ✅    |     |                                            |
| Hidden                          | ✅      | ✅    |     |                                            |
| HTML                            | ✅      | ✅    |     |                                            |
| JS                              | ✅      | ✅    |     |                                            |
| Mac                             | ❌      | ❌    |     |                                            |
| MaxHeight                       | ✅      | ✅    |     |                                            |
| MaxWidth                        | ✅      | ✅    |     |                                            |
| MinHeight                       | ✅      | ✅    |     |                                            |
| MinWidth                        | ✅      | ✅    |     |                                            |
| Name                            | ✅      | ✅    |     |                                            |
| OpenInspectorOnStartup          |         |       |     |                                            |
| StartState                      | ✅      |       |     |                                            |
| Title                           | ✅      | ✅    |     |                                            |
| URL                             | ✅      | ✅    |     |                                            |
| Width                           | ✅      | ✅    |     |                                            |
| Windows                         | ✅      | ❌    | ❌  |                                            |
| X                               | ✅      | ✅    |     |                                            |
| Y                               | ✅      | ✅    |     |                                            |
| Zoom                            |         |       |     |                                            |
| ZoomControlEnabled              |         |       |     |                                            |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
| ------------------------ | ------- | ----- | --- | ----- |
| Default Application Menu | ✅      | ✅    | ✅  |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes                                                                |
| ------------------ | ------- | ----- | --- | -------------------------------------------------------------------- |
| Icon               | ✅      | ✅    | ✅  | Windows has default icons for light/dark mode & supports PNG or ICO. |
| Label              | ❌      | ✅    | ✅  |                                                                      |
| Label (ANSI Codes) | ❌      |       |     |                                                                      |
| Menu               | ✅      | ✅    | ✅  |                                                                      |

### Methods

| Method                        | Windows | Linux | Mac | Notes                              |
| ----------------------------- | ------- | ----- | --- | ---------------------------------- |
| setLabel(label string)        | ❌      | ✅    | ✅  |                                    |
| run()                         | ✅      | ✅    | ✅  |                                    |
| setIcon(icon []byte)          | ✅      | ✅    | ✅  |                                    |
| setMenu(menu \*Menu)          | ✅      | ✅    | ✅  |                                    |
| setIconPosition(position int) | ❌      | ✅    | ✅  |                                    |
| setTemplateIcon(icon []byte)  | ❌      | ✅    | ✅  |                                    |
| destroy()                     | ✅      | ✅    | ✅  |                                    |
| setDarkModeIcon(icon []byte)  | ✅      | ✅    | ✅  | Darkmode isn't handled yet (linux) |

## Cross Platform Events

Mapping native events to cross-platform events.

| Event                    | Windows | Linux | Mac             | Notes |
| ------------------------ | ------- | ----- | --------------- | ----- |
| WindowWillClose          |         |       | WindowWillClose |       |
| WindowDidClose           |         |       |                 |       |
| WindowDidResize          |         |       |                 |       |
| WindowDidHide            |         |       |                 |       |
| ApplicationWillTerminate |         |       |                 |       |

... Add more

## Bindings Generation

Working well.

## Models Generation

Working well.

## Task file

Contains a lot needed for development.

## Theme

| Mode   | Windows | Linux | Mac | Notes |
| ------ | ------- | ----- | --- | ----- |
| Dark   | ✅      |       |     |       |
| Light  | ✅      |       |     |       |
| System | ✅      |       |     |       |

## NSIS Installer

TBD

## Templates

All templates are working.

## Plugins

Built-in plugin support:

| Plugin          | Windows | Linux | Mac | Notes |
| --------------- | ------- | ----- | --- | ----- |
| Browser         | ✅      |       | ✅  |       |
| KV Store        | ✅      | ✅    | ✅  |       |
| Log             | ✅      | ✅    | ✅  |       |
| Single Instance | ✅      |       | ✅  |       |
| SQLite          | ✅      | ✅    | ✅  |       |
| Start at login  | ✅      |       | ✅  |       |
| Server          |         |       |     |       |

TODO:

- Ensure each plugin has a JS wrapper that can be injected into the window.

## Packaging

|                 | Windows | Linux | Mac | Notes |
| --------------- | ------- | ----- | --- | ----- |
| Icon Generation | ✅      |       | ✅  |       |
| Icon Embedding  | ✅      |       | ✅  |       |
| Info.plist      | ❌      |       | ✅  |       |
| NSIS Installer  |         |       | ❌  |       |
| Mac bundle      | ❌      |       | ✅  |       |
| Windows exe     | ✅      |       | ❌  |       |

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes                                          |
| ------- | ------- | ----- | --- | ---------------------------------------------- |
| Resize  | ✅      |       | ✅  |                                                |
| Drag    | ✅      | ✅    | ✅  | Linux - can always drag with `Meta`+left mouse |

## Mac Specific

| Feature      | Mac | Notes |
| ------------ | --- | ----- |
| Translucency | ✅  |       |

### Mac Options

| Feature                 | Default           | Notes                                                |
| ----------------------- | ----------------- | ---------------------------------------------------- |
| Backdrop                | MacBackdropNormal | Standard solid window                                |
| DisableShadow           | false             |                                                      |
| TitleBar                |                   | Standard window decorations by default               |
| Appearance              | DefaultAppearance |                                                      |
| InvisibleTitleBarHeight | 0                 | Creates an invisible title bar for frameless windows |
| DisableShadow           | false             | Disables the window drop shadow                      |

## Windows Specific

| Feature       | Windows | Notes |
| ------------- | ------- | ----- |
| Translucency  | ✅      |       |
| Custom Themes | ✅      |       |

### Windows Options

| Feature                           | Default       | Notes                                       |
| --------------------------------- | ------------- | ------------------------------------------- |
| BackdropType                      | Solid         |                                             |
| DisableIcon                       | false         |                                             |
| Theme                             | SystemDefault |                                             |
| CustomTheme                       | nil           |                                             |
| DisableFramelessWindowDecorations | false         |                                             |
| WindowMask                        | nil           | Makes the window the contents of the bitmap |

## Linux Specific

Implementation details for the functions utilized by the `*_linux.go` files are
located in the following files:

- `linux_cgo.go`: CGo implementation
- `linux_purego.go`: PureGo implementation

### CGO

By default CGO is utilized to compile the Linux port. This prevents easy
cross-compilation and so the PureGo implementation is also being simultaneously
developed.

### Purego

The examples can be compiled using the following command:

```bash
CGO_ENABLED=0 go build -tags purego
```

:::note

Things are currently not working after the refactor

:::
