# Status

Status of features in v3. Incomplete - please add as you see fit.

## Application

Application interface methods

| Method                                                        | Windows | Linux | Mac | Notes |
|---------------------------------------------------------------|---------|-------|-----|-------|
| run() error                                                   |         |       | ✅   |       |
| destroy()                                                     |         |       | ✅   |       |
| setApplicationMenu(menu *Menu)                                |         |       | ✅   |       |
| name() string                                                 |         |       | ✅   |       |
| getCurrentWindowID() uint                                     |         |       | ✅   |       |
| showAboutDialog(name string, description string, icon []byte) |         |       | ✅   |       |
| setIcon(icon []byte)                                          |         |       | ✅   |       |
| on(id uint)                                                   |         |       | ✅   |       |
| dispatchOnMainThread(id uint)                                 |         |       | ✅   |       |
| hide()                                                        |         |       | ✅   |       |
| show()                                                        |         |       | ✅   |       |
| getPrimaryScreen() (*Screen, error)                           |         |       | ✅   |       |
| getScreens() ([]*Screen, error)                               |         |       | ✅   |       |

## Webview Window

Webview Window Interface Methods

| Method                                             | Windows | Linux | Mac | Notes |    
|----------------------------------------------------|---------|-------|-----|-------|    
| setTitle(title string)                             |         |       | ✅   |       |
| setSize(width, height int)                         |         |       | ✅   |       |
| setAlwaysOnTop(alwaysOnTop bool)                   |         |       | ✅   |       |
| setURL(url string)                                 |         |       | ✅   |       |
| setResizable(resizable bool)                       |         |       | ✅   |       |
| setMinSize(width, height int)                      |         |       | ✅   |       |
| setMaxSize(width, height int)                      |         |       | ✅   |       |
| execJS(js string)                                  |         |       | ✅   |       |
| restore()                                          |         |       | ✅   |       |
| setBackgroundColour(color *RGBA)                   |         |       | ✅   |       |
| run()                                              |         |       | ✅   |       |
| center()                                           |         |       | ✅   |       |
| size() (int, int)                                  |         |       | ✅   |       |
| width() int                                        |         |       | ✅   |       |
| height() int                                       |         |       | ✅   |       |
| position() (int, int)                              |         |       | ✅   |       |
| destroy()                                          |         |       | ✅   |       |
| reload()                                           |         |       | ✅   |       |
| forceReload()                                      |         |       | ✅   |       |
| toggleDevTools()                                   |         |       | ✅   |       |
| zoomReset()                                        |         |       | ✅   |       |
| zoomIn()                                           |         |       | ✅   |       |
| zoomOut()                                          |         |       | ✅   |       |
| getZoom() float64                                  |         |       | ✅   |       |
| setZoom(zoom float64)                              |         |       | ✅   |       |
| close()                                            |         |       | ✅   |       |
| zoom()                                             |         |       | ✅   |       |
| setHTML(html string)                               |         |       | ✅   |       |
| setPosition(x int, y int)                          |         |       | ✅   |       |
| on(eventID uint)                                   |         |       | ✅   |       |
| minimise()                                         |         |       | ✅   |       |
| unminimise()                                       |         |       | ✅   |       |
| maximise()                                         |         |       | ✅   |       |
| unmaximise()                                       |         |       | ✅   |       |
| fullscreen()                                       |         |       | ✅   |       |
| unfullscreen()                                     |         |       | ✅   |       |
| isMinimised() bool                                 |         |       | ✅   |       |
| isMaximised() bool                                 |         |       | ✅   |       |
| isFullscreen() bool                                |         |       | ✅   |       |
| disableSizeConstraints()                           |         |       | ✅   |       |
| setFullscreenButtonEnabled(enabled bool)           |         |       | ✅   |       |
| show()                                             |         |       | ✅   |       |
| hide()                                             |         |       | ✅   |       |
| getScreen() (*Screen, error)                       |         |       | ✅   |       |
| setFrameless(bool)                                 |         |       | ✅   |       |
| openContextMenu(menu *Menu, data *ContextMenuData) |         |       | ✅   |       |

## Runtime

### Application

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| Quit    |         |       | ✅   |       |
| Hide    |         |       | ✅   |       |
| Show    |         |       | ✅   |       |

### Dialogs

| Feature  | Windows | Linux | Mac | Notes |
|----------|---------|-------|-----|-------|
| Info     |         |       | ✅   |       |
| Warning  |         |       | ✅   |       |
| Error    |         |       | ✅   |       |
| Question |         |       | ✅   |       |
| OpenFile |         |       | ✅   |       |
| SaveFile |         |       | ✅   |       |

### Clipboard

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|-----|-------|
| SetText |         |       | ✅   |       |
| Text    |         |       | ✅   |       |

### ContextMenu

| Feature         | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| OpenContextMenu |         |       | ✅   |       |

### Screens

| Feature    | Windows | Linux | Mac | Notes |
|------------|---------|-------|-----|-------|
| GetAll     |         |       | ✅   |       |
| GetPrimary |         |       | ✅   |       |
| GetCurrent |         |       | ✅   |       |

### Window

| Feature             | Windows | Linux | Mac | Notes                                                                                |
|---------------------|---------|-------|-----|--------------------------------------------------------------------------------------|
| SetTitle            |         |       | ✅   |                                                                                      |
| SetSize             |         |       | ✅   |                                                                                      |
| Size                |         |       | ✅   |                                                                                      |
| SetPosition         |         |       | ✅   |                                                                                      |
| Position            |         |       | ✅   |                                                                                      |
| FullScreen          |         |       | ✅   |                                                                                      |
| UnFullscreen        |         |       | ✅   |                                                                                      |
| Minimise            |         |       | ✅   |                                                                                      |
| UnMinimise          |         |       | ✅   |                                                                                      |
| Maximise            |         |       | ✅   |                                                                                      |
| UnMaximise          |         |       | ✅   |                                                                                      |
| Show                |         |       | ✅   |                                                                                      |
| Hide                |         |       | ✅   |                                                                                      |
| Center              |         |       | ✅   |                                                                                      |
| SetBackgroundColour |         |       | ✅   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetAlwaysOnTop      |         |       | ✅   |                                                                                      |
| SetResizable        |         |       | ✅   |                                                                                      |
| SetMinSize          |         |       | ✅   |                                                                                      |
| SetMaxSize          |         |       | ✅   |                                                                                      |
| Width               |         |       | ✅   |                                                                                      |
| Height              |         |       | ✅   |                                                                                      |
| ZoomIn              |         |       | ✅   | Increase view scale                                                                  |
| ZoomOut             |         |       | ✅   | Decrease view scale                                                                  |
| ZoomReset           |         |       | ✅   | Reset view scale                                                                     |
| GetZoom             |         |       | ✅   | Get current view scale                                                               |
| SetZoom             |         |       | ✅   | Set view scale                                                                       |
| Screen              |         |       | ✅   | Get screen for window                                                                |

### Log

To log or not to log? System logger vs custom logger.

## Menu

| Event                    | Windows | Linux | Mac | Notes |
|--------------------------|---------|-------|-----|-------|
| Default Application Menu |         |       | ✅   |       |

## Tray Menus

| Feature            | Windows | Linux | Mac | Notes |
|--------------------|---------|-------|-----|-------|
| Icon               |         |       | ✅   |       |
| Label              |         |       | ✅   |       |
| Label (ANSI Codes) |         |       |     |       |
| Menu               |         |       | ✅   |       |

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

| Plugin          | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| Dark            |         |       |    |       |
| Light           |         |       |    |       |
| System          |         |       |    |       |

## NSIS Installer

TBD

## Templates

TBD

## Plugins

Built-in plugin support:

| Plugin          | Windows | Linux | Mac | Notes |
|-----------------|---------|-------|-----|-------|
| Browser         |         |       | ✅   |       |
| KV Store        |         |       | ✅   |       |
| Log             |         |       | ✅   |       |
| Single Instance |         |       | ✅   |       |
| SQLite          |         |       | ✅   |       |
| Start at login  |         |       | ✅   |       |
| Server          |         |       |     |       |

## Packaging

|                 | Windows | Linux | Mac | Notes | 
|-----------------|---------|-------|-----|-------| 
| Icon Generation |         |       | ✅   |       | 
| Icon Embedding  |         |       | ✅   |       | 
| Info.plist      |         |       | ✅   |       | 
| NSIS Installer  |         |       |     |       | 
| Mac bundle      |         |       | ✅   |       | 
| Windows exe     |         |       |     |       | 

## Frameless Windows

| Feature | Windows | Linux | Mac | Notes |
|---------|---------|-------|----|-------|
| Resize  |         |       |    |       |
| Drag    |         |       |    |       |

## Mac Specific

- [x] Translucency

## Windows Specific

## Linux Specific
