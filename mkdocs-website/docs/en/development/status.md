# Status

Status of features in v3.

!!! note

        This list is a mixture of public and internal API support.<br/>
        It is not complete and probably not up to date.

## Known Issues

- Linux is not yet up to feature parity with Windows/Mac

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes |
| ------------------------------------------------------------- | ------- | ----- | --- | ----- |
| run() error                                                   | Y       | Y     | Y   |       |
| destroy()                                                     |         | Y     | Y   |       |
| setApplicationMenu(menu \*Menu)                               | Y       | Y     | Y   |       |
| name() string                                                 |         | Y     | Y   |       |
| getCurrentWindowID() uint                                     | Y       | Y     | Y   |       |
| showAboutDialog(name string, description string, icon []byte) |         | Y     | Y   |       |
| setIcon(icon []byte)                                          | -       | Y     | Y   |       |
| on(id uint)                                                   |         |       | Y   |       |
| dispatchOnMainThread(fn func())                               | Y       | Y     | Y   |       |
| hide()                                                        | Y       | Y     | Y   |       |
| show()                                                        | Y       | Y     | Y   |       |
| getPrimaryScreen() (\*Screen, error)                          |         | Y     | Y   |       |
| getScreens() ([]\*Screen, error)                              |         | Y     | Y   |       |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes                                    |
| -------------------------------------------------- | ------- | ----- | --- | ---------------------------------------- |
| center()                                           | Y       | Y     | Y   |                                          |
| close()                                            | y       | Y     | Y   |                                          |
| destroy()                                          |         | Y     | Y   |                                          |
| execJS(js string)                                  | y       | Y     | Y   |                                          |
| focus()                                            | Y       | Y     |     |                                          |
| forceReload()                                      |         | Y     | Y   |                                          |
| fullscreen()                                       | Y       | Y     | Y   |                                          |
| getScreen() (\*Screen, error)                      | y       | Y     | Y   |                                          |
| getZoom() float64                                  |         | Y     | Y   |                                          |
| height() int                                       | Y       | Y     | Y   |                                          |
| hide()                                             | Y       | Y     | Y   |                                          |
| isFullscreen() bool                                | Y       | Y     | Y   |                                          |
| isMaximised() bool                                 | Y       | Y     | Y   |                                          |
| isMinimised() bool                                 | Y       | Y     | Y   |                                          |
| maximise()                                         | Y       | Y     | Y   |                                          |
| minimise()                                         | Y       | Y     | Y   |                                          |
| nativeWindowHandle() (uintptr, error)              | Y       | Y     | Y   |                                          |
| on(eventID uint)                                   | y       |       | Y   |                                          |
| openContextMenu(menu *Menu, data *ContextMenuData) | y       | Y     | Y   |                                          |
| relativePosition() (int, int)                      | Y       | Y     | Y   |                                          |
| reload()                                           | y       | Y     | Y   |                                          |
| run()                                              | Y       | Y     | Y   |                                          |
| setAlwaysOnTop(alwaysOnTop bool)                   | Y       | Y     | Y   |                                          |
| setBackgroundColour(color RGBA)                    | Y       | Y     | Y   |                                          |
| setEnabled(bool)                                   |         | Y     | Y   |                                          |
| setFrameless(bool)                                 |         | Y     | Y   |                                          |
| setFullscreenButtonEnabled(enabled bool)           | -       | Y     | Y   | There is no fullscreen button in Windows |
| setHTML(html string)                               | Y       | Y     | Y   |                                          |
| setMaxSize(width, height int)                      | Y       | Y     | Y   |                                          |
| setMinSize(width, height int)                      | Y       | Y     | Y   |                                          |
| setRelativePosition(x int, y int)                  | Y       | Y     | Y   |                                          |
| setResizable(resizable bool)                       | Y       | Y     | Y   |                                          |
| setSize(width, height int)                         | Y       | Y     | Y   |                                          |
| setTitle(title string)                             | Y       | Y     | Y   |                                          |
| setURL(url string)                                 | Y       | Y     | Y   |                                          |
| setZoom(zoom float64)                              | Y       | Y     | Y   |                                          |
| show()                                             | Y       | Y     | Y   |                                          |
| size() (int, int)                                  | Y       | Y     | Y   |                                          |
| toggleDevTools()                                   | Y       | Y     | Y   |                                          |
| unfullscreen()                                     | Y       | Y     | Y   |                                          |
| unmaximise()                                       | Y       | Y     | Y   |                                          |
| unminimise()                                       | Y       | Y     | Y   |                                          |
| width() int                                        | Y       | Y     | Y   |                                          |
| zoom()                                             |         | Y     | Y   |                                          |
| zoomIn()                                           | Y       | Y     | Y   |                                          |
| zoomOut()                                          | Y       | Y     | Y   |                                          |
| zoomReset()                                        | Y       | Y     | Y   |                                          |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
| ------- | ------- | ----- | --- | ----- |
| Quit    | Y       | Y     | Y   |       |
| Hide    | Y       | Y     | Y   |       |
| Show    | Y       |       | Y   |       |

### Dialogs

| Feature  | Windows | Linux | Mac | Notes |
| -------- | ------- | ----- | --- | ----- |
| Info     | Y       | Y     | Y   |       |
| Warning  | Y       | Y     | Y   |       |
| Error    | Y       | Y     | Y   |       |
| Question | Y       | Y     | Y   |       |
| OpenFile | Y       | Y     | Y   |       |
| SaveFile | Y       | Y     | Y   |       |

### Clipboard

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| SetText | Y       | Y     | Y   |       |
| Text    | Y       | Y     | Y   |       |

### ContextMenu

| Feature          | Windows | Linux | Mac | Notes |
|------------------|---------|-------|-----|-------|
| OpenContextMenu  | Y       | Y     | Y   |       |
| On By Default    |         |       |     |       |
| Control via HTML | Y       |       |     |       |

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
| GetAll     | Y       | Y     | Y   |       |
| GetPrimary | Y       | Y     | Y   |       |
| GetCurrent | Y       | Y     | Y   |       |

### System

| Feature    | Windows | Linux | Mac | Notes |
| ---------- | ------- | ----- | --- | ----- |
| IsDarkMode |         |       | Y   |       |

### Window

Y = Supported U = Untested

- = Not available

| Feature             | Windows | Linux | Mac | Notes                                                                                |
| ------------------- | ------- | ----- | --- | ------------------------------------------------------------------------------------ |
| Center              | Y       | Y     | Y   |                                                                                      |
| Focus               | Y       | Y     |     |                                                                                      |
| FullScreen          | Y       | Y     | Y   |                                                                                      |
| GetZoom             | Y       | Y     | Y   | Get current view scale                                                               |
| Height              | Y       | Y     | Y   |                                                                                      |
| Hide                | Y       | Y     | Y   |                                                                                      |
| Maximise            | Y       | Y     | Y   |                                                                                      |
| Minimise            | Y       | Y     | Y   |                                                                                      |
| RelativePosition    | Y       | Y     | Y   |                                                                                      |
| Screen              | Y       | Y     | Y   | Get screen for window                                                                |
| SetAlwaysOnTop      | Y       | Y     | Y   |                                                                                      |
| SetBackgroundColour | Y       | Y     | Y   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetEnabled          | Y       | U     | -   | Set the window to be enabled/disabled                                                |
| SetMaxSize          | Y       | Y     | Y   |                                                                                      |
| SetMinSize          | Y       | Y     | Y   |                                                                                      |
| SetRelativePosition | Y       | Y     | Y   |                                                                                      |
| SetResizable        | Y       | Y     | Y   |                                                                                      |
| SetSize             | Y       | Y     | Y   |                                                                                      |
| SetTitle            | Y       | Y     | Y   |                                                                                      |
| SetZoom             | Y       | Y     | Y   | Set view scale                                                                       |
| Show                | Y       | Y     | Y   |                                                                                      |
| Size                | Y       | Y     | Y   |                                                                                      |
| UnFullscreen        | Y       | Y     | Y   |                                                                                      |
| UnMaximise          | Y       | Y     | Y   |                                                                                      |
| UnMinimise          | Y       | Y     | Y   |                                                                                      |
| Width               | Y       | Y     | Y   |                                                                                      |
| ZoomIn              | Y       | Y     | Y   | Increase view scale                                                                  |
| ZoomOut             | Y       | Y     | Y   | Decrease view scale                                                                  |
| ZoomReset           | Y       | Y     | Y   | Reset view scale                                                                     |

### Window Options

A 'Y' in the table below indicates that the option has been tested and is
applied when the window is created. An 'X' indicates that the option is not
supported by the platform.

| Feature                         | Windows | Linux | Mac | Notes                                      |
|---------------------------------|---------|-------|-----|--------------------------------------------|
| AlwaysOnTop                     | Y       | Y     |     |                                            |
| BackgroundColour                | Y       | Y     |     |                                            |
| BackgroundType                  |         |       |     | Acrylic seems to work but the others don't |
| CSS                             | Y       | Y     |     |                                            |
| DevToolsEnabled                 | Y       | Y     | Y   |                                            |
| DisableResize                   | Y       | Y     |     |                                            |
| EnableDragAndDrop               |         | Y     |     |                                            |
| EnableFraudulentWebsiteWarnings |         |       |     |                                            |
| Focused                         | Y       | Y     |     |                                            |
| Frameless                       | Y       | Y     |     |                                            |
| FullscreenButtonEnabled         | Y       |       |     |                                            |
| Height                          | Y       | Y     |     |                                            |
| Hidden                          | Y       | Y     |     |                                            |
| HTML                            | Y       | Y     |     |                                            |
| JS                              | Y       | Y     |     |                                            |
| Mac                             | -       | -     |     |                                            |
| MaxHeight                       | Y       | Y     |     |                                            |
| MaxWidth                        | Y       | Y     |     |                                            |
| MinHeight                       | Y       | Y     |     |                                            |
| MinWidth                        | Y       | Y     |     |                                            |
| Name                            | Y       | Y     |     |                                            |
| OpenInspectorOnStartup          |         |       |     |                                            |
| StartState                      | Y       |       |     |                                            |
| Title                           | Y       | Y     |     |                                            |
| URL                             | Y       | Y     |     |                                            |
| Width                           | Y       | Y     |     |                                            |
| Windows                         | Y       | -     | -   |                                            |
| X                               | Y       | Y     |     |                                            |
| Y                               | Y       | Y     |     |                                            |
| Zoom                            |         |       |     |                                            |
| ZoomControlEnabled              |         |       |     |                                            |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
| ------------------------ | ------- | ----- | --- | ----- |
| Default Application Menu | Y       | Y     | Y   |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes                                                                |
|--------------------|---------|-------|-----|----------------------------------------------------------------------|
| Icon               | Y       | Y     | Y   | Windows has default icons for light/dark mode & supports PNG or ICO. |
| Label              | -       | Y     | Y   |                                                                      |
| Label (ANSI Codes) | -       |       |     |                                                                      |
| Menu               | Y       | Y     | Y   |                                                                      |

### Methods

| Method                        | Windows | Linux | Mac | Notes |
| ----------------------------- | ------- | ----- | --- | ----- |
| setLabel(label string)        | -       | Y     | Y   |       |
| run()                         | Y       | Y     | Y   |       |
| setIcon(icon []byte)          | Y       | Y     | Y   |       |
| setMenu(menu \*Menu)          | Y       | Y     | Y   |       |
| setIconPosition(position int) | -       | Y     | Y   |       |
| setTemplateIcon(icon []byte)  | -       | Y     | Y   |       |
| destroy()                     | Y       | Y     | Y   |       |
| setDarkModeIcon(icon []byte)  | Y       | Y     | Y   | Darkmode isn't handled yet (linux)      |

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
| Dark   | Y       |       |     |       |
| Light  | Y       |       |     |       |
| System | Y       |       |     |       |

## NSIS Installer

TBD

## Templates

All templates are working.

## Plugins

Built-in plugin support:

| Plugin          | Windows | Linux | Mac | Notes |
| --------------- | ------- | ----- | --- | ----- |
| Browser         | Y       |       | Y   |       |
| KV Store        | Y       | Y     | Y   |       |
| Log             | Y       | Y     | Y   |       |
| Single Instance | Y       |       | Y   |       |
| SQLite          | Y       | Y     | Y   |       |
| Start at login  | Y       |       | Y   |       |
| Server          |         |       |     |       |

TODO:

- Ensure each plugin has a JS wrapper that can be injected into the window.

## Packaging

|                 | Windows | Linux | Mac | Notes |
| --------------- | ------- | ----- | --- | ----- |
| Icon Generation | Y       |       | Y   |       |
| Icon Embedding  | Y       |       | Y   |       |
| Info.plist      | -       |       | Y   |       |
| NSIS Installer  |         |       | -   |       |
| Mac bundle      | -       |       | Y   |       |
| Windows exe     | Y       |       | -   |       |

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes                                          |
| ------- | ------- | ----- | --- | ---------------------------------------------- |
| Resize  | Y       |       | Y   |                                                |
| Drag    | Y       | Y     | Y   | Linux - can always drag with `Meta`+left mouse |

## Mac Specific

- [x] Translucency

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

- [x] Translucency
- [x] Custom Themes

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

- linux_cgo.go: CGo implementation
- linux_purego.go: PureGo implementation

### CGO

By default CGO is utilized to compile the Linux port. This prevents easy
cross-compilation and so the PureGo implementation is also being simultaneously
developed.

### Purego

The examples can be compiled using the following command:

    CGO_ENABLED=0 go build -tags purego

Note: things are currently not working after the refactor
