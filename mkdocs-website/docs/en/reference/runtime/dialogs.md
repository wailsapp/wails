# Dialogs 

The Dialogs API provides a set of common dialog windows that can be used in your Wails application.

## Message Dialogs

Message dialogs are used to display information, ask questions, or show warnings and errors.

You can create one of the 4 supported message dialogs:

=== "Info"

    ```go
    myDialog := application.InfoDialog()  
    ```

=== "Question"

    ```go
    myDialog := application.QuestionDialog()  
    ```

=== "Warning"

    ```go
    myDialog := application.WarningDialog()  
    ```

=== "Error"

    ```go
    myDialog := application.ErrorDialog()  
    ```

Message dialogs have the following methods:

### `SetTitle(title string) *MessageDialog`

The `SetTitle()` method sets the title of the message dialog.

```go
myDialog.SetTitle("My Dialog")
```

### `SetIcon(icon []byte) *MessageDialog`

The `SetIcon()` method sets the icon for the message dialog.

```go
myDialog.SetIcon(iconData)
```

### `AddButton(label string) *Button`

The `AddButton()` method adds a button to the question dialog (Mac and Linux only).

```go
button := myDialog.AddButton("OK")
```

### `SetDefaultButton(button *Button) *MessageDialog`

The `SetDefaultButton()` method sets the default button for the message dialog (Mac and Linux only).

```go
myDialog.SetDefaultButton(button)
```

### `SetCancelButton(button *Button) *MessageDialog`

The `SetCancelButton()` method sets the cancel button for the message dialog (Mac and Linux only).

```go
myDialog.SetCancelButton(button)
```

### `SetMessage(message string) *MessageDialog`

The `SetMessage()` method sets the message text for the dialog.

```go
myDialog.SetMessage("Hello, World!")
```

### `Show()`

The `Show()` method displays the message dialog.

```go
myDialog.Show()
```

## Open File Dialog

The open file dialog allows users to select one or more files.

### `OpenFileDialog() *OpenFileDialogStruct`

The `application.OpenFileDialog()` function creates a new open file dialog.

```go
myDialog := application.newOpenFileDialog()
```

Once the dialog is created, you can configure it using the following methods:

### `CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct`

The `CanChooseFiles()` method sets whether the dialog can choose files.

```go
myDialog.CanChooseFiles(true)
```

### `CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct`

The `CanChooseDirectories()` method sets whether the dialog can choose directories.

```go
myDialog.CanChooseDirectories(false)
```

### `AddFilter(displayName, pattern string) *OpenFileDialogStruct`

The `AddFilter()` method adds a file filter to the dialog.

```go
myDialog.AddFilter("Image Files", "*.jpg;*.png")
```

### `PromptForSingleSelection() (string, error)`

The `PromptForSingleSelection()` method displays the dialog and returns the selected file path.

```go
filePath, err := myDialog.PromptForSingleSelection()
if err != nil {
    // Handle error
}
```

### `PromptForMultipleSelection() ([]string, error)`

The `PromptForMultipleSelection()` method displays the dialog and returns the selected file paths.

```go
filePaths, err := myDialog.PromptForMultipleSelection()
if err != nil {
    // Handle error
}
```

### `OpenFileDialogWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct`

If you prefer to use a more structured approach, you can use the `OpenFileDialogWithOptions()` function.

```go
options := &application.OpenFileDialogOptions{
    CanChooseFiles: true,
    CanChooseDirectories: false,
    Filters: []application.FileFilter{
        {
            DisplayName: "Image Files",
            Pattern: "*.jpg;*.png",
        },
    },
}
myDialog := application.OpenFileDialogWithOptions(options)
```

## Save File Dialog

The save file dialog allows users to select a file to save.

### `application.SaveFileDialog() *SaveFileDialogStruct`

The `application.SaveFileDialog()` function creates a new save file dialog.

```go
myDialog := application.SaveFileDialog()
```

Once the dialog is created, you can configure it using the following methods:

### `AddFilter(displayName, pattern string) *SaveFileDialogStruct`

The `AddFilter()` method adds a file filter to the dialog.

```go
myDialog.AddFilter("Image Files", "*.jpg;*.png")
```

### `PromptForSingleSelection() (string, error)`

The `PromptForSingleSelection()` method displays the dialog and returns the selected file path.

```go
filePath, err := dialog.PromptForSingleSelection()
if err != nil {
    // Handle error
}
```

### `SaveFileDialogWithOptions(options *SaveFileDialogOptions) *SaveFileDialogStruct`

If you prefer to use a more structured approach, you can use the `SaveFileDialogWithOptions()` function.

```go
options := &application.SaveFileDialogOptions{
    Filters: []application.FileFilter{
        {
            DisplayName: "Image Files",
            Pattern: "*.jpg;*.png",
        },
    },
}
myDialog := application.SaveFileDialogWithOptions(options)
```