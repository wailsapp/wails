### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

`RegisterContextMenu()` registers a context menu with a given name. This menu
can be used later in the application.

```go

    // Create a new menu
    ctxmenu := app.NewMenu()

    // Register the menu as a context menu
    app.RegisterContextMenu("MyContextMenu", ctxmenu)
```

### SetMenu

API: `SetMenu(menu *Menu)`

`SetMenu()` sets the menu for the application. On Mac, this will be the global
menu. For Windows and Linux, this will be the default menu for any new window
created.

```go
    // Create a new menu
    menu := app.NewMenu()

    // Set the menu for the application
    app.SetMenu(menu)
```
