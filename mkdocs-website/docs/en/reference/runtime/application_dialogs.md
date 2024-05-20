
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
