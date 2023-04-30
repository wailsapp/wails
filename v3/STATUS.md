# Status

Status of features in v3. Incomplete - please add as you see fit.

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes |
|---------------------------------------------------------------|---------|-------|-----|-------|
| run() error                                                   |         |       | Y   |       |
| destroy()                                                     |         |       | Y   |       |
| setApplicationMenu(menu *Menu)                                |         |       | Y   |       |
| name() string                                                 |         |       | Y   |       |
| getCurrentWindowID() uint                                     |         |       | Y   |       |
| showAboutDialog(name string, description string, icon []byte) |         |       | Y   |       |
| setIcon(icon []byte)                                          |         |       | Y   |       |
| on(id uint)                                                   |         |       | Y   |       |
| dispatchOnMainThread(fn func())                               | Y       |       | Y   |       |
| hide()                                                        |         |       | Y   |       |
| show()                                                        |         |       | Y   |       |
| getPrimaryScreen() (*Screen, error)                           |         |       | Y   |       |
| getScreens() ([]*Screen, error)                               |         |       | Y   |       |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes |    
|----------------------------------------------------|---------|-------|-----|-------|    
| setTitle(title string)                             | Y       |       | Y   |       |
| setSize(width, height int)                         | Y       |       | Y   |       |
| setAlwaysOnTop(alwaysOnTop bool)                   | Y       |       | Y   |       |
| setURL(url string)                                 |         |       | Y   |       |
| setResizable(resizable bool)                       | Y       |       | Y   |       |
| setMinSize(width, height int)                      |         |       | Y   |       |
| setMaxSize(width, height int)                      |         |       | Y   |       |
| execJS(js string)                                  |         |       | Y   |       |
| restore()                                          |         |       | Y   |       |
| setBackgroundColour(color RGBA)                    | Y       |       | Y   |       |
| run()                                              | Y       |       | Y   |       |
| center()                                           | Y       |       | Y   |       |
| size() (int, int)                                  |         |       | Y   |       |
| width() int                                        | Y       |       | Y   |       |
| height() int                                       | Y       |       | Y   |       |
| position() (int, int)                              | Y       |       | Y   |       |
| destroy()                                          |         |       | Y   |       |
| reload()                                           |         |       | Y   |       |
| forceReload()                                      |         |       | Y   |       |
| toggleDevTools()                                   |         |       | Y   |       |
| zoomReset()                                        |         |       | Y   |       |
| zoomIn()                                           |         |       | Y   |       |
| zoomOut()                                          |         |       | Y   |       |
| getZoom() float64                                  |         |       | Y   |       |
| setZoom(zoom float64)                              |         |       | Y   |       |
| close()                                            |         |       | Y   |       |
| zoom()                                             |         |       | Y   |       |
| setHTML(html string)                               |         |       | Y   |       |
| setPosition(x int, y int)                          |         |       | Y   |       |
| on(eventID uint)                                   |         |       | Y   |       |
| minimise()                                         | Y       |       | Y   |       |
| unminimise()                                       | Y       |       | Y   |       |
| maximise()                                         | Y       |       | Y   |       |
| unmaximise()                                       | Y       |       | Y   |       |
| fullscreen()                                       |         |       | Y   |       |
| unfullscreen()                                     |         |       | Y   |       |
| isMinimised() bool                                 | Y       |       | Y   |       |
| isMaximised() bool                                 | Y       |       | Y   |       |
| isFullscreen() bool                                |         |       | Y   |       |
| disableSizeConstraints()                           |         |       | Y   |       |
| setFullscreenButtonEnabled(enabled bool)           |         |       | Y   |       |
| show()                                             | Y       |       | Y   |       |
| hide()                                             | Y       |       | Y   |       |
| getScreen() (*Screen, error)                       |         |       | Y   |       |
| setFrameless(bool)                                 |         |       | Y   |       |
| openContextMenu(menu *Menu, data *ContextMenuData) |         |       | Y   |       |
| nativeWindowHandle() (uintptr, error)              | Y       |       |     |       |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Quit    |         |       | Y   |       |
| Hide    |         |       | Y   |       |
| Show    |         |       | Y   |       |

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
| GetAll     |         |       | Y   |       |
| GetPrimary |         |       | Y   |       |
| GetCurrent |         |       | Y   |       |

### Window

| Feature             | Windows | Linux | Mac | Notes                                                                                |
|---------------------|---------|-------|-----|--------------------------------------------------------------------------------------|
| SetTitle            |         |       | Y   |                                                                                      |
| SetSize             |         |       | Y   |                                                                                      |
| Size                |         |       | Y   |                                                                                      |
| SetPosition         |         |       | Y   |                                                                                      |
| Position            |         |       | Y   |                                                                                      |
| FullScreen          |         |       | Y   |                                                                                      |
| UnFullscreen        |         |       | Y   |                                                                                      |
| Minimise            |         |       | Y   |                                                                                      |
| UnMinimise          |         |       | Y   |                                                                                      |
| Maximise            |         |       | Y   |                                                                                      |
| UnMaximise          |         |       | Y   |                                                                                      |
| Show                |         |       | Y   |                                                                                      |
| Hide                |         |       | Y   |                                                                                      |
| Center              |         |       | Y   |                                                                                      |
| SetBackgroundColour |         |       | Y   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetAlwaysOnTop      |         |       | Y   |                                                                                      |
| SetResizable        |         |       | Y   |                                                                                      |
| SetMinSize          |         |       | Y   |                                                                                      |
| SetMaxSize          |         |       | Y   |                                                                                      |
| Width               |         |       | Y   |                                                                                      |
| Height              |         |       | Y   |                                                                                      |
| ZoomIn              |         |       | Y   | Increase view scale                                                                  |
| ZoomOut             |         |       | Y   | Decrease view scale                                                                  |
| ZoomReset           |         |       | Y   | Reset view scale                                                                     |
| GetZoom             |         |       | Y   | Get current view scale                                                               |
| SetZoom             |         |       | Y   | Set view scale                                                                       |
| Screen              |         |       | Y   | Get screen for window                                                                |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
|--------------------------|---------|-------|-----|-------|
| Default Application Menu |         |       | Y   |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes |
|--------------------|---------|-------|-----|-------|
| Icon               |         |       | Y   |       |
| Label              |         |       | Y   |       |
| Label (ANSI Codes) |         |       |     |       |
| Menu               |         |       | Y   |       |

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

| Plugin | Windows | Linux | Mac | Notes |
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
| KV Store        |         |       | Y   |       |
| Log             |         |       | Y   |       |
| Single Instance |         |       | Y   |       |
| SQLite          |         |       | Y   |       |
| Start at login  |         |       | Y   |       |
| Server          |         |       |     |       |

## Packaging

|                 | Windows | Linux | Mac | Notes | 
|-----------------|---------|-------|-----|-------| 
| Icon Generation |         |       | Y   |       | 
| Icon Embedding  |         |       | Y   |       | 
| Info.plist      |         |       | Y   |       | 
| NSIS Installer  |         |       |     |       | 
| Mac bundle      |         |       | Y   |       | 
| Windows exe     |         |       |     |       | 

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Resize  |         |       |     |       |
| Drag    |         |       |     |       |

## Mac Specific

- [x] Translucency

## Windows Specific

- [x] Translucency
- [x] Custom Themes

## Linux Specific
