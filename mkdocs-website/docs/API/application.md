# Application

The application API assists in creating an application using the Wails framework.

### New

API: `New(appOptions Options) *App`

`New(appOptions Options)` creates a new application using the given application options . It applies default values for unspecified options, merges them with the provided ones, initializes and returns an instance of the application.

In case of an error during initialization, the application is stopped with the error message provided.

It should be noted that if a global application instance already exists, that instance will be returned instead of creating a new one.

### Get

`Get()` returns the global application instance. It's useful when you need to access the application from different parts of your code.

### Capabilities

API: `Capabilities() capabilities.Capabilities`

`Capabilities()` retrieves a map of capabilities that the application currently has. Capabilities can be about different features the operating system provides, like webview features.

### GetPID

API: `GetPID() int`

`GetPID()` returns the Process ID of the application.

### Run

API: `Run() error`

`Run()` starts the execution of the application and its components.


### Quit

API: `Quit()`

`Quit()` quits the application by destroying windows and potentially other components.

### IsDarkMode

API: `IsDarkMode() bool`

`IsDarkMode()` checks if the application is running in dark mode. It returns a boolean indicating whether dark mode is enabled.

### Hide

API: `Hide()`

`Hide()` hides the application window.

### Show

API: `Show()`

`Show()` shows the application window.

### NewWebviewWindow

API: `NewWebviewWindow() *WebviewWindow`

`NewWebviewWindow()` creates a new Webview window with default options, and returns it.

### NewWebviewWindowWithOptions

API: `NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow`

`NewWebviewWindowWithOptions()` creates a new webview window with custom options. The newly created window is added to a map of windows managed by the application.

### OnWindowCreation

API: `OnWindowCreation(callback func(window *WebviewWindow))`

`OnWindowCreation()` registers a callback function to be called when a window is created.

### GetWindowByName

API: `GetWindowByName(name string) *WebviewWindow`

`GetWindowByName()` fetches and returns a window with a specific name.

### CurrentWindow

API: `CurrentWindow() *WebviewWindow`

`CurrentWindow()` fetches and returns a pointer to the currently active window in the application. If there is no window, it returns nil.


### NewSystemTray

API: `NewSystemTray() *SystemTray`

`NewSystemTray()` creates and returns a new system tray instance.



### NewMenu

API: `NewMenu() *Menu`

This method, belonging to App struct and not Menu struct, also initialises and returns a new `Menu`.



### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

`RegisterContextMenu()` registers a context menu with a given name. This menu can be used later in the application.

### SetMenu

API: `SetMenu(menu *Menu)`

`SetMenu()` sets the menu for the application.

### ShowAboutDialog

API: `ShowAboutDialog()`

`ShowAboutDialog()` shows an "About" dialog box. It can show the application's name, description and icon.

### Info

`InfoDialog()` creates and returns a new instance of `MessageDialog` with an `InfoDialogType`. This dialog is typically used to display informational messages to the user.

### Question

`QuestionDialog()` creates and returns a new instance of `MessageDialog` with a `QuestionDialogType`. This dialog is often used to ask a question to the user and expect a response.

### Warning

`WarningDialog()` creates and returns a new instance of `MessageDialog` with a `WarningDialogType`. As the name suggests, this dialog is primarily used to display warning messages to the user.

### Error

`ErrorDialog()` creates and returns a new instance of `MessageDialog` with an `ErrorDialogType`. This dialog is designed to be used when you need to display an error message to the user.

### OpenFile

`OpenFileDialog()` creates and returns a new `OpenFileDialogStruct`. This dialog prompts the user to select one or more files from their file system.

### SaveFile

`SaveFileDialog()` creates and returns a new `SaveFileDialogStruct`. This dialog prompts the user to choose a location on their file system where a file should be saved.

### OpenDirectory

`OpenDirectoryDialog()` creates and returns a new instance of `MessageDialog` with an `OpenDirectoryDialogType`. This dialog enables the user to choose a directory from their file system.


### On

API: `On(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`On()` registers an event listener for specific application events. The callback function provided will be triggered when the corresponding event occurs. The function returns a function that can be called to remove the listener.

### RegisterHook

API: `RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`RegisterHook()` registers a callback to be run as a hook during specific events. These hooks are run before listeners attached with `On()`. The function returns a function that can be called to remove the hook.

### GetPrimaryScreen

API: `GetPrimaryScreen() (*Screen, error)`

`GetPrimaryScreen()` returns the primary screen of the system.

### GetScreens

API: `GetScreens() ([]*Screen, error)`

`GetScreens()` returns information about all screens attached to the system.

This is a brief summary of the exported methods in the provided `App` struct. Do note that for more detailed functionality or considerations, refer to the actual Go code or further internal documentation.

## Options

--8<-- "../../../v3/pkg/application/options_application.go"
