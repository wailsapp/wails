# Windows Jump List Example

This example demonstrates how to implement Windows Jump Lists in a Wails v3 application.

## What are Jump Lists?

Jump Lists are a Windows feature introduced in Windows 7 that provide quick access to recently opened files and common tasks. They appear when you right-click on a taskbar button or hover over it.

## Features

- **Custom Categories**: Create custom categories like "Recent Documents" or "Frequent Items"
- **Tasks**: Add application-specific tasks that appear at the bottom of the jump list
- **Runtime Configuration**: Update jump lists dynamically during application runtime
- **Cross-Platform Safe**: The API is designed to be no-op on non-Windows platforms

## Usage

```go
// Create a jump list
jumpList := app.CreateJumpList()

// Add a custom category
category := application.JumpListCategory{
    Name: "Recent Documents",
    Items: []application.JumpListItem{
        {
            Type:        application.JumpListItemTypeTask,
            Title:       "Document.txt",
            Description: "Open Document.txt",
            FilePath:    "C:\\path\\to\\app.exe",
            Arguments:   "--open Document.txt",
            IconPath:    "C:\\path\\to\\app.exe",
            IconIndex:   0,
        },
    },
}
jumpList.AddCategory(category)

// Add tasks (with empty category name)
tasks := application.JumpListCategory{
    Name: "", // Empty name indicates tasks section
    Items: []application.JumpListItem{
        {
            Type:        application.JumpListItemTypeTask,
            Title:       "New Document",
            Description: "Create a new document",
            FilePath:    "C:\\path\\to\\app.exe",
            Arguments:   "--new",
        },
    },
}
jumpList.AddCategory(tasks)

// Apply the jump list
err := jumpList.Apply()
```

## API Reference

### JumpListItem

- `Type`: The type of jump list item (currently only `JumpListItemTypeTask` is supported)
- `Title`: The display title of the item
- `Description`: A tooltip description that appears on hover
- `FilePath`: The path to the executable to launch (usually your app)
- `Arguments`: Command-line arguments to pass when the item is clicked
- `IconPath`: Path to the icon file (can be an .exe, .dll, or .ico file)
- `IconIndex`: Index of the icon if the IconPath contains multiple icons

### JumpListCategory

- `Name`: The category name (use empty string for tasks)
- `Items`: Array of JumpListItem objects

### Methods

- `app.CreateJumpList()`: Creates a new jump list instance
- `jumpList.AddCategory(category)`: Adds a category to the jump list
- `jumpList.ClearCategories()`: Removes all categories
- `jumpList.Apply()`: Applies the jump list to the taskbar

## Platform Support

This feature is only available on Windows. On macOS and Linux, all jump list methods are no-ops and will not cause any errors.

## Building and Running

```bash
# From the example directory
go build -tags desktop
./jumplist.exe
```

Make sure to pin the application to your taskbar to see the jump list in action!

## Notes

- The application ID used for the jump list is taken from the application name in Options
- Jump list items will launch new instances of your application with the specified arguments
- You should handle these arguments in your application's startup code
- Icons can be extracted from executables, DLLs, or standalone .ico files