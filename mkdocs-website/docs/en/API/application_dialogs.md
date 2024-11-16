
!!! warning "MacOS Dialogs and Application Lifecycle"

    If you show dialogs during application startup or file open events, you should set `ApplicationShouldTerminateAfterLastWindowClosed` to `false` to prevent the application from terminating when those dialogs close. Otherwise, the application may quit before your main window appears.
    
    ```go
    app := application.New(application.Options{
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
        // ... rest of options
    })
    ```

    Alternatively, you can show startup dialogs after the main window has been displayed:

    ```go
    var filename string
	app.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
		filename = event.Context().Filename()
    })

    window.OnWindowEvent(events.Common.WindowShow, func(event *application.WindowEvent) {
        application.InfoDialog().
            SetTitle("File Opened").
            SetMessage("Application opened with file: " + filename).
            Show()
    })
    ```

### ShowAboutDialog

API: `ShowAboutDialog()`

`ShowAboutDialog()` shows an "About" dialog box. It can show the application's
name, description and icon.

```go
    // Show the about dialog
    app.ShowAboutDialog()
```

### Info

API: `InfoDialog()`

`InfoDialog()` creates and returns a new instance of `MessageDialog` with an
`InfoDialogType`. This dialog is typically used to display informational
messages to the user.

### Question

API: `QuestionDialog()`

`QuestionDialog()` creates and returns a new instance of `MessageDialog` with a
`QuestionDialogType`. This dialog is often used to ask a question to the user
and expect a response.

### Warning

API: `WarningDialog()`

`WarningDialog()` creates and returns a new instance of `MessageDialog` with a
`WarningDialogType`. As the name suggests, this dialog is primarily used to
display warning messages to the user.

### Error

API: `ErrorDialog()`

`ErrorDialog()` creates and returns a new instance of `MessageDialog` with an
`ErrorDialogType`. This dialog is designed to be used when you need to display
an error message to the user.

### OpenFile

API: `OpenFileDialog()`

`OpenFileDialog()` creates and returns a new `OpenFileDialogStruct`. This dialog
prompts the user to select one or more files from their file system.

### SaveFile

API: `SaveFileDialog()`

`SaveFileDialog()` creates and returns a new `SaveFileDialogStruct`. This dialog
prompts the user to choose a location on their file system where a file should
be saved.

### OpenDirectory

API: `OpenDirectoryDialog()`

`OpenDirectoryDialog()` creates and returns a new instance of `MessageDialog`
with an `OpenDirectoryDialogType`. This dialog enables the user to choose a
directory from their file system.
