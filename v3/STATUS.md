# Status

Status of features in v3. Incomplete - please add as you see fit.

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes                        |
|---------------------------------------------------------------|---------|-------|-----|------------------------------|
| run() error                                                   |         | Y     | Y   |                              |
| destroy()                                                     |         | Y     | Y   |                              |
| setApplicationMenu(menu *Menu)                                |         |       | Y   |                              |
| name() string                                                 |         |       | Y   |                              |
| getCurrentWindowID() uint                                     |         | Y     | Y   |                              |
| showAboutDialog(name string, description string, icon []byte) |         | Y     | Y   | [linux] No icon possible yet |
| setIcon(icon []byte)                                          |         |       | Y   |                              |
| on(id uint)                                                   |         |       | Y   |                              |
| dispatchOnMainThread(fn func())                               | Y       | Y     | Y   |                              |
| hide()                                                        | Y       |       | Y   |                              |
| show()                                                        | Y       |       | Y   |                              |
| getPrimaryScreen() (*Screen, error)                           |         | Y     | Y   |                              |
| getScreens() ([]*Screen, error)                               |         | Y     | Y   |                              |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes                                    |
|----------------------------------------------------|---------|-------|-----|------------------------------------------|
| center()                                           | Y       | Y     | Y   |                                          |
| close()                                            |         | Y     | Y   |                                          |
| destroy()                                          |         | Y     | Y   |                                          |
| execJS(js string)                                  |         | Y     | Y   |                                          |
| focus()                                            | Y       | Y     |     |                                          |
| forceReload()                                      |         | Y     | Y   |                                          |
| fullscreen()                                       | Y       | Y     | Y   |                                          |
| getScreen() (*Screen, error)                       |         | Y     | Y   |                                          |
| getZoom() float64                                  |         | Y     | Y   |                                          |
| height() int                                       | Y       | Y     | Y   |                                          |
| hide()                                             | Y       | Y     | Y   |                                          |
| isFullscreen() bool                                | Y       | Y     | Y   |                                          |
| isMaximised() bool                                 | Y       | Y     | Y   |                                          |
| isMinimised() bool                                 | Y       | Y     | Y   |                                          |
| maximise()                                         | Y       | Y     | Y   |                                          |
| minimise()                                         | Y       | Y     | Y   |                                          |
| nativeWindowHandle() (uintptr, error)              | Y       | Y     |     |                                          |
| on(eventID uint)                                   |         |       | Y   |                                          |
| openContextMenu(menu *Menu, data *ContextMenuData) |         |       | Y   |                                          |
| position() (int, int)                              | Y       | Y     | Y   |                                          |
| reload()                                           |         | Y     | Y   |                                          |
| run()                                              | Y       | Y     | Y   |                                          |
| setAlwaysOnTop(alwaysOnTop bool)                   | Y       | Y     | Y   |                                          |
| setBackgroundColour(color RGBA)                    | Y       | Y     | Y   |                                          |
| setFrameless(bool)                                 |         | Y     | Y   |                                          |
| setFullscreenButtonEnabled(enabled bool)           | -       | Y     | Y   | There is no fullscreen button in Windows |
| setHTML(html string)                               |         | Y     | Y   |                                          |
| setMaxSize(width, height int)                      | Y       | Y     | Y   |                                          |
| setMinSize(width, height int)                      | Y       | Y     | Y   |                                          |
| setPosition(x int, y int)                          | Y       | Y     | Y   |                                          |
| setResizable(resizable bool)                       | Y       | Y     | Y   |                                          |
| setSize(width, height int)                         | Y       | Y     | Y   |                                          |
| setTitle(title string)                             | Y       | Y     | Y   |                                          |
| setURL(url string)                                 |         | Y     | Y   |                                          |
| setZoom(zoom float64)                              |         | Y     | Y   |                                          |
| show()                                             | Y       | Y     | Y   |                                          |
| size() (int, int)                                  | Y       | Y     | Y   |                                          |
| toggleDevTools()                                   |         | Y     | Y   |                                          |
| unfullscreen()                                     | Y       | Y     | Y   |                                          |
| unmaximise()                                       | Y       | Y     | Y   |                                          |
| unminimise()                                       | Y       | Y     | Y   |                                          |
| width() int                                        | Y       | Y     | Y   |                                          |
| zoom()                                             |         | Y     | Y   |                                          |
| zoomIn()                                           |         | Y     | Y   |                                          |
| zoomOut()                                          |         | Y     | Y   |                                          |
| zoomReset()                                        |         | Y     | Y   |                                          |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Quit    |         | Y     | Y   |       |
| Hide    | Y       |       | Y   |       |
| Show    | Y       |       | Y   |       |

### Dialogs

| Feature  | Windows | Linux | Mac | Notes |
|----------|---------|-------|-----|-------|
| Info     |         |       | Y   |       |
| Warning  |         |       | Y   |       |
| Error    |         |       | Y   |       |
| Question |         |       | Y   |       |
| OpenFile |         |       | Y   |       |
| SaveFile |         |       | Y   |       |

### Clipboard

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| SetText |         |       | Y   |       |
| Text    |         |       | Y   |       |

### ContextMenu

| Feature         | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| OpenContextMenu |         |       | Y   |       |

### Screens

| Feature    | Windows | Linux | Mac | Notes |
|------------|---------|-------|-----|-------|
| GetAll     |         | Y     | Y   |       |
| GetPrimary |         | Y     | Y   |       |
| GetCurrent |         | Y     | Y   |       |

### Window

| Feature             | Windows | Linux | Mac | Notes                                                                                |
|---------------------|---------|-------|-----|--------------------------------------------------------------------------------------|
|                     |         |       |     |                                                                                      |
| SetTitle            | Y       | Y     | Y   |                                                                                      |
| SetSize             | Y       | Y     | Y   |                                                                                      |
| Size                | Y       | Y     | Y   |                                                                                      |
| SetPosition         | Y       | Y     | Y   |                                                                                      |
| Position            | Y       | Y     | Y   |                                                                                      |
| Focus               | Y       | Y     |     |                                                                                      |
| FullScreen          | Y       | Y     | Y   |                                                                                      |
| UnFullscreen        | Y       | Y     | Y   |                                                                                      |
| Minimise            | Y       | Y     | Y   |                                                                                      |
| UnMinimise          | Y       | Y     | Y   |                                                                                      |
| Maximise            | Y       | Y     | Y   |                                                                                      |
| UnMaximise          | Y       | Y     | Y   |                                                                                      |
| Show                | Y       | Y     | Y   |                                                                                      |
| Hide                | Y       | Y     | Y   |                                                                                      |
| Center              | Y       | Y     | Y   |                                                                                      |
| SetBackgroundColour | Y       | Y     | Y   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetAlwaysOnTop      | Y       | Y     | Y   |                                                                                      |
| SetResizable        | Y       | Y     | Y   |                                                                                      |
| SetMinSize          | Y       | Y     | Y   |                                                                                      |
| SetMaxSize          | Y       | Y     | Y   |                                                                                      |
| Width               | Y       | Y     | Y   |                                                                                      |
| Height              | Y       | Y     | Y   |                                                                                      |
| ZoomIn              |         | Y     | Y   | Increase view scale                                                                  |
| ZoomOut             |         | Y     | Y   | Decrease view scale                                                                  |
| ZoomReset           |         | Y     | Y   | Reset view scale                                                                     |
| GetZoom             |         | Y     | Y   | Get current view scale                                                               |
| SetZoom             |         | Y     | Y   | Set view scale                                                                       |
| Screen              |         | Y     | Y   | Get screen for window                                                                |

### Window Options

A 'Y' in the table below indicates that the option has been tested and is applied when the window is created.
An 'X' indicates that the option is not supported by the platform.

| Feature                         | Windows | Linux | Mac | Notes                                             |
|---------------------------------|---------|-------|-----|---------------------------------------------------|
| Name                            |         |       |     |                                                   |
| Title                           | Y       |       |     |                                                   |
| Width                           | Y       | Y     |     |                                                   |
| Height                          | Y       | Y     |     |                                                   |
| AlwaysOnTop                     | Y       | Y     |     |                                                   |
| URL                             |         |       |     |                                                   |
| DisableResize                   | Y       | Y     |     |                                                   |
| Frameless                       |         | Y     |     |                                                   |
| MinWidth                        | Y       | Y     |     |                                                   |
| MinHeight                       | Y       | Y     |     |                                                   |
| MaxWidth                        | Y       | Y     |     |                                                   |
| MaxHeight                       | Y       | Y     |     |                                                   |
| StartState                      | Y       |       |     |                                                   |
| Mac                             | -       | -     |     |                                                   |
| BackgroundType                  |         |       |     | Acrylic seems to work but the others don't        |
| BackgroundColour                | Y       | Y     |     |                                                   |
| HTML                            |         | Y     |     |                                                   |
| JS                              |         | Y     |     |                                                   |
| CSS                             |         | Y     |     |                                                   |
| X                               |         | Y     |     |                                                   |
| Y                               |         | Y     |     |                                                   |
| HideOnClose                     |         | Y     |     |                                                   |
| FullscreenButtonEnabled         |         | ?     |     | [linux] How is this different from DisableResize? |
| Hidden                          | Y       |       |     |                                                   |
| EnableFraudulentWebsiteWarnings |         |       |     |                                                   |
| Zoom                            |         | Y     |     |                                                   |
| EnableDragAndDrop               |         | Y     |     |                                                   |
| Windows                         | Y       | -     | -   |                                                   |
| Focused                         | Y       |       |     |                                                   |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
|--------------------------|---------|-------|-----|-------|
| Default Application Menu |         | Y     | Y   |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes                                                                |
|--------------------|---------|-------|-----|----------------------------------------------------------------------|
| Icon               | Y       |       | Y   | Windows has default icons for light/dark mode & supports PNG or ICO. |
| Label              | -       |       | Y   |                                                                      |
| Label (ANSI Codes) | -       |       |     |                                                                      |
| Menu               |         |       | Y   |                                                                      |

## Cross Platform Events

Mapping native events to cross-platform events.

| Event                    | Windows | Linux | Mac             | Notes |
|--------------------------|---------|-------|-----------------|-------|
| WindowWillClose          |         |       | WindowWillClose |       |
| WindowDidClose           |         |       |                 |       |
| WindowDidResize          |         |       |                 |       |
| WindowDidHide            |         |       |                 |       |
| ApplicationWillTerminate |         |       |                 |       |

... Add more

## Bindings Generation

TBD

## Models Generation

TBD

## Task file

TBD

## Theme

| Mode   | Windows | Linux | Mac | Notes |
|--------|---------|-------|-----|-------|
| Dark   | Y       |       |     |       |
| Light  | Y       |       |     |       |
| System | Y       |       |     |       |

## NSIS Installer

TBD

## Templates

TBD

## Plugins

Built-in plugin support:

| Plugin          | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| Browser         |         |       | Y   |       |
| KV Store        |         | Y     | Y   |       |
| Log             |         | Y     | Y   |       |
| Single Instance |         |       | Y   |       |
| SQLite          |         | Y     | Y   |       |
| Start at login  |         |       | Y   |       |
| Server          |         |       |     |       |

## Packaging

|                 | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| Icon Generation |         |       | Y   |       |
| Icon Embedding  |         |       | Y   |       |
| Info.plist      |         |       | Y   |       |
| NSIS Installer  |         |       | -   |       |
| Mac bundle      |         |       | Y   |       |
| Windows exe     |         |       | -   |       |

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes                                         |
|---------|---------|-------|-----|-----------------------------------------------|
| Resize  |         |       |     |                                               |
| Drag    |         | Y     |     | Linux - can always drag with `Alt`+left mouse |


## Mac Specific

- [x] Translucency

## Windows Specific

- [x] Translucency
- [x] Custom Themes

### Windows Options

| Feature                           | Default | Notes                                       |
|-----------------------------------|---------|---------------------------------------------|
| BackdropType                      |         |                                             |
| DisableIcon                       |         |                                             |
| Theme                             |         |                                             |
| CustomTheme                       |         |                                             |
| DisableFramelessWindowDecorations |         |                                             |
| WindowMask                        | nil     | Makes the window the contents of the bitmap |

    // Select the type of translucent backdrop. Requires Windows 11 22621 or later.
    BackdropType BackdropType
    // Disable the icon in the titlebar
    DisableIcon bool
    // Theme. Defaults to SystemDefault which will use whatever the system theme is. The application will follow system theme changes.
    Theme Theme
    // Custom colours for dark/light mode
    CustomTheme *ThemeSettings

    // Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
    // "Rounded Corners" are only available on Windows 11.
    DisableFramelessWindowDecorations bool

    // WindowMask is used to set the window shape. Use a PNG with an alpha channel to create a custom shape.
    WindowMask []byte

## Linux Specific
