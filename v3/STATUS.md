# Status

Status of features in v3. Incomplete - please add as you see fit.


## Major Blockers

- [ ] Dev Support - What will it look like?
- [ ] Translucency on Windows doesn't work right
- [ ] Events - cross-platform events, simplified event handling?
- [ ] Error handling - needs to be revised. Centralised system error handling?
- [ ] Documentation - needs to be written

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes |
|---------------------------------------------------------------|---------|-------|-----|-------|
| run() error                                                   | Y       |       | Y   |       |
| destroy()                                                     |         |       | Y   |       |
| setApplicationMenu(menu *Menu)                                | Y       |       | Y   |       |
| name() string                                                 |         |       | Y   |       |
| getCurrentWindowID() uint                                     | Y       |       | Y   |       |
| showAboutDialog(name string, description string, icon []byte) |         |       | Y   |       |
| setIcon(icon []byte)                                          | -       |       | Y   |       |
| on(id uint)                                                   |         |       | Y   |       |
| dispatchOnMainThread(fn func())                               | Y       |       | Y   |       |
| hide()                                                        | Y       |       | Y   |       |
| show()                                                        | Y       |       | Y   |       |
| getPrimaryScreen() (*Screen, error)                           |         |       | Y   |       |
| getScreens() ([]*Screen, error)                               |         |       | Y   |       |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes                                    |
|----------------------------------------------------|---------|-------|-----|------------------------------------------|
| center()                                           | Y       |       | Y   |                                          |
| close()                                            | Y       |       | Y   |                                          |
| destroy()                                          |         |       | Y   |                                          |
| execJS(js string)                                  | Y       |       | Y   |                                          |
| focus()                                            | Y       |       | Y   |                                          |
| forceReload()                                      |         |       | Y   |                                          |
| fullscreen()                                       | Y       |       | Y   |                                          |
| getScreen() (*Screen, error)                       | Y       |       | Y   |                                          |
| getZoom() float64                                  |         |       | Y   |                                          |
| height() int                                       | Y       |       | Y   |                                          |
| hide()                                             | Y       |       | Y   |                                          |
| isFullscreen() bool                                | Y       |       | Y   |                                          |
| isMaximised() bool                                 | Y       |       | Y   |                                          |
| isMinimised() bool                                 | Y       |       | Y   |                                          |
| maximise()                                         | Y       |       | Y   |                                          |
| minimise()                                         | Y       |       | Y   |                                          |
| nativeWindowHandle() (uintptr, error)              | Y       |       | Y   |                                          |
| on(eventID uint)                                   | Y       |       | Y   |                                          |
| openContextMenu(menu *Menu, data *ContextMenuData) | Y       |       | Y   |                                          |
| position() (int, int)                              | Y       |       | Y   |                                          |
| reload()                                           | Y       |       | Y   |                                          |
| run()                                              | Y       |       | Y   |                                          |
| setAlwaysOnTop(alwaysOnTop bool)                   | Y       |       | Y   |                                          |
| setBackgroundColour(color RGBA)                    | Y       |       | Y   |                                          |
| setFrameless(bool)                                 |         |       | Y   |                                          |
| setFullscreenButtonEnabled(enabled bool)           | -       |       | Y   | There is no fullscreen button in Windows |
| setHTML(html string)                               | Y       |       | Y   |                                          |
| setMaxSize(width, height int)                      | Y       |       | Y   |                                          |
| setMinSize(width, height int)                      | Y       |       | Y   |                                          |
| setPosition(x int, y int)                          | Y       |       | Y   |                                          |
| setResizable(resizable bool)                       | Y       |       | Y   |                                          |
| setSize(width, height int)                         | Y       |       | Y   |                                          |
| setTitle(title string)                             | Y       |       | Y   |                                          |
| setURL(url string)                                 | Y       |       | Y   |                                          |
| setZoom(zoom float64)                              | Y       |       | Y   |                                          |
| show()                                             | Y       |       | Y   |                                          |
| size() (int, int)                                  | Y       |       | Y   |                                          |
| toggleDevTools()                                   | Y       |       | Y   |                                          |
| unfullscreen()                                     | Y       |       | Y   |                                          |
| unmaximise()                                       | Y       |       | Y   |                                          |
| unminimise()                                       | Y       |       | Y   |                                          |
| width() int                                        | Y       |       | Y   |                                          |
| zoom()                                             |         |       | Y   |                                          |
| zoomIn()                                           | Y       |       | Y   |                                          |
| zoomOut()                                          | Y       |       | Y   |                                          |
| zoomReset()                                        | Y       |       | Y   |                                          |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Quit    | Y       |       | Y   |       |
| Hide    | Y       |       | Y   |       |
| Show    | Y       |       | Y   |       |

### Dialogs

| Feature  | Windows | Linux | Mac | Notes |
|----------|---------|-------|-----|-------|
| Info     | Y       |       | Y   |       |
| Warning  | Y       |       | Y   |       |
| Error    | Y       |       | Y   |       |
| Question | Y       |       | Y   |       |
| OpenFile | Y       |       | Y   |       |
| SaveFile | Y       |       | Y   |       |

### Clipboard

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| SetText | Y       |       | Y   |       |
| Text    | Y       |       | Y   |       |

### ContextMenu

| Feature         | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| OpenContextMenu | Y       |       | Y   |       |

### Screens

| Feature    | Windows | Linux | Mac | Notes |
|------------|---------|-------|-----|-------|
| GetAll     | Y       |       | Y   |       |
| GetPrimary | Y       |       | Y   |       |
| GetCurrent | Y       |       | Y   |       |

### Window

| Feature             | Windows | Linux | Mac | Notes                                                                                |
|---------------------|---------|-------|-----|--------------------------------------------------------------------------------------|
| SetTitle            | Y       |       | Y   |                                                                                      |
| SetSize             | Y       |       | Y   |                                                                                      |
| Size                | Y       |       | Y   |                                                                                      |
| SetPosition         | Y       |       | Y   |                                                                                      |
| Position            | Y       |       | Y   |                                                                                      |
| Focus               | Y       |       | Y   |                                                                                      |
| FullScreen          | Y       |       | Y   |                                                                                      |
| UnFullscreen        | Y       |       | Y   |                                                                                      |
| Minimise            | Y       |       | Y   |                                                                                      |
| UnMinimise          | Y       |       | Y   |                                                                                      |
| Maximise            | Y       |       | Y   |                                                                                      |
| UnMaximise          | Y       |       | Y   |                                                                                      |
| Show                | Y       |       | Y   |                                                                                      |
| Hide                | Y       |       | Y   |                                                                                      |
| Center              | Y       |       | Y   |                                                                                      |
| SetBackgroundColour | Y       |       | Y   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetAlwaysOnTop      | Y       |       | Y   |                                                                                      |
| SetResizable        | Y       |       | Y   |                                                                                      |
| SetMinSize          | Y       |       | Y   |                                                                                      |
| SetMaxSize          | Y       |       | Y   |                                                                                      |
| Width               | Y       |       | Y   |                                                                                      |
| Height              | Y       |       | Y   |                                                                                      |
| ZoomIn              | Y       |       | Y   | Increase view scale                                                                  |
| ZoomOut             | Y       |       | Y   | Decrease view scale                                                                  |
| ZoomReset           | Y       |       | Y   | Reset view scale                                                                     |
| GetZoom             | Y       |       | Y   | Get current view scale                                                               |
| SetZoom             | Y       |       | Y   | Set view scale                                                                       |
| Screen              | Y       |       | Y   | Get screen for window                                                                |

### Window Options

A 'Y' in the table below indicates that the option has been tested and is applied when the window is created.
An 'X' indicates that the option is not supported by the platform.

| Feature                         | Windows | Linux | Mac | Notes                                      |
|---------------------------------|---------|-------|-----|--------------------------------------------|
| Name                            | Y       |       |     |                                            |
| Title                           | Y       |       |     |                                            |
| Width                           | Y       |       |     |                                            |
| Height                          | Y       |       |     |                                            |
| AlwaysOnTop                     | Y       |       |     |                                            |
| URL                             | Y       |       |     |                                            |
| DisableResize                   | Y       |       |     |                                            |
| Frameless                       | Y       |       |     |                                            |
| MinWidth                        | Y       |       |     |                                            |
| MinHeight                       | Y       |       |     |                                            |
| MaxWidth                        | Y       |       |     |                                            |
| MaxHeight                       | Y       |       |     |                                            |
| StartState                      | Y       |       |     |                                            |
| Mac                             | -       | -     |     |                                            |
| BackgroundType                  |         |       |     | Acrylic seems to work but the others don't |
| BackgroundColour                | Y       |       |     |                                            |
| HTML                            | Y       |       |     |                                            |
| JS                              | Y       |       |     |                                            |
| CSS                             | Y       |       |     |                                            |
| X                               | Y       |       |     |                                            |
| Y                               | Y       |       |     |                                            |
| HideOnClose                     | Y       |       |     |                                            |
| FullscreenButtonEnabled         | Y       |       |     |                                            |
| Hidden                          | Y       |       |     |                                            |
| EnableFraudulentWebsiteWarnings |         |       |     |                                            |
| Zoom                            |         |       |     |                                            |
| ZoomControlEnabled              |         |       |     |                                            |
| OpenInspectorOnStartup          |         |       |     |                                            |
| EnableDragAndDrop               |         |       |     |                                            |
| Windows                         | Y       | -     | -   |                                            |
| Focused                         | Y       |       |     |                                            |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
|--------------------------|---------|-------|-----|-------|
| Default Application Menu | Y       |       | Y   |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes                                                                |
|--------------------|---------|-------|-----|----------------------------------------------------------------------|
| Icon               | Y       |       | Y   | Windows has default icons for light/dark mode & supports PNG or ICO. |
| Label              | -       |       | Y   |                                                                      |
| Label (ANSI Codes) | -       |       |     |                                                                      |
| Menu               | Y       |       | Y   |                                                                      |

### Methods

| Method                        | Windows | Linux | Mac | Notes |
|-------------------------------|---------|-------|-----|-------|
| setLabel(label string)        | -       |       | Y   |       |
| run()                         | Y       |       | Y   |       |
| setIcon(icon []byte)          | Y       |       | Y   |       |
| setMenu(menu *Menu)           | Y       |       | Y   |       |
| setIconPosition(position int) | -       |       | Y   |       |
| setTemplateIcon(icon []byte)  | -       |       | Y   |       |
| destroy()                     | Y       |       | Y   |       |
| setDarkModeIcon(icon []byte)  | Y       |       | Y   |       |

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
| Browser         | Y       |       | Y   |       |
| KV Store        | Y       |       | Y   |       |
| Log             | Y       |       | Y   |       |
| Single Instance | Y       |       | Y   |       |
| SQLite          | Y       |       | Y   |       |
| Start at login  |         |       | Y   |       |
| Server          |         |       |     |       |


TODO: 

- Ensure each plugin has a JS wrapper that can be injected into the window.

## Packaging

|                 | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| Icon Generation | Y       |       | Y   |       |
| Icon Embedding  | Y       |       | Y   |       |
| Info.plist      | -       |       | Y   |       |
| NSIS Installer  |         |       | -   |       |
| Mac bundle      | -       |       | Y   |       |
| Windows exe     | Y       |       | -   |       |

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Resize  |         |       |     |       |
| Drag    |         |       |     |       |

## Mac Specific

- [x] Translucency

### Mac Options

| Feature                 | Default           | Notes                                                |
|-------------------------|-------------------|------------------------------------------------------|
| Backdrop                | MacBackdropNormal | Standard solid window                                |
| DisableShadow           | false             |                                                      |
| TitleBar                |                   | Standard window decorations by default               |
| Appearance              | DefaultAppearance |                                                      |
| InvisibleTitleBarHeight | 0                 | Creates an invisible title bar for frameless windows |
| DisableShadow           | false             | Disables the window drop shadow                      |

## Windows Specific

- [ ] Translucency
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

## Linux Specific

## Examples

| Example      | Windows                                                | Linux | Mac |
|--------------|--------------------------------------------------------|-------|-----|
| binding      | NO                                                     |       |     |
| build        | Yes (Debug + Prod)                                     |       |     |
| clipboard    | Yes                                                    |       |     |
| contextmenus | Yes                                                    |       |     |
| dialogs      | Almost                                                 |       |     |
| drag-n-drop  | NO                                                     |       |     |
| events       | NO                                                     |       |     |
| frameless    | Needs resizing                                         |       |     |
| kitchensink  | Yes                                                    |       |     |
| menu         | Yes                                                    |       |     |
| plain        | Yes                                                    |       |     |
| plugins      | Needs StartAtLogin                                     |       |     |
| screen       | Yes                                                    |       |     |
| systray      | Yes                                                    |       |     |
| window       | Size of initial window not right with application menu |       |     |
| windowjs     | Example not complete                                   |       |     |
| wml          | Yes                                                    |       |     |

# Alpha Release TODO

- [ ] Check all runtime methods are available in the JS runtime

# Beta Release TODO

- [ ] Make better looking examples