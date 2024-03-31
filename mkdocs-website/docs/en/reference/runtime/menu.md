# Menu

The Menu API provides a way to create and manage application menus.

## `NewMenu() *Menu`

The `NewMenu()` function creates a new Menu instance.

```go
menu := application.NewMenu()
```

### `Add(label string) *MenuItem`

The `Add()` method adds a new menu item to the menu and returns the menu item.

```go
menuItem := menu.Add("File")
```

### `AddSeparator()`

The `AddSeparator()` method adds a separator to the menu.

```go
menu.AddSeparator()
```

### `AddCheckbox(label string, enabled bool) *MenuItem`

The `AddCheckbox()` method adds a checkbox menu item to the menu.

```go
checkboxItem := menu.AddCheckbox("Option", true)
```

### `AddRadio(label string, enabled bool) *MenuItem`

The `AddRadio()` method adds a radio menu item to the menu.

```go
radioItem := menu.AddRadio("Option", true)
```

### `AddSubmenu(label string) *Menu`

The `AddSubmenu()` method adds a new submenu to the menu.

```go
submenu := menu.AddSubmenu("Submenu")
```

### `AddRole(role Role) *Menu`

The `AddRole()` method adds a predefined menu role to the menu.

```go
menu.AddRole(application.RoleAbout)
```

### `Update()`

The `Update()` method updates the menu, applying any changes made to the menu items.

```go
menu.Update()
```

### `SetLabel(label string)`

The `SetLabel()` method sets the label for the menu.

```go
menu.SetLabel("My Menu")
```

### `Append(in *Menu)`

The `Append()` method appends the items of one menu to another.

```go
menu.Append(anotherMenu)
```

The Menu API provides a flexible and extensible way to create and manage application menus in your Wails application.
```

```
# MenuItem

The MenuItem API provides a way to create and manage menu items in your Wails application. This is
used in conjunction with the Menu API.

## `SetTooltip(tooltip string) *MenuItem`

The `SetTooltip()` method sets the tooltip for the menu item.

```go
menuItem.SetTooltip("Open a file")
```

### `SetLabel(label string) *MenuItem`

The `SetLabel()` method sets the label for the menu item.

```go
menuItem.SetLabel("Open")
```

### `SetEnabled(enabled bool) *MenuItem`

The `SetEnabled()` method sets the enabled state of the menu item.

```go
menuItem.SetEnabled(true)
```

### `SetBitmap(bitmap []byte) *MenuItem`

The `SetBitmap()` method sets the bitmap for the menu item.

```go
menuItem.SetBitmap(bitmapData)
```

### `SetChecked(checked bool) *MenuItem`

The `SetChecked()` method sets the checked state of the menu item.

```go
menuItem.SetChecked(true)
```

### `SetHidden(hidden bool) *MenuItem`

The `SetHidden()` method sets the hidden state of the menu item.

```go
menuItem.SetHidden(false)
```

### `OnClick(callback func(*Context)) *MenuItem`

The `OnClick()` method sets the click handler for the menu item.

```go
menuItem.OnClick(func(ctx *application.Context) {
// Handle click
})
```

### `SetAccelerator(shortcut string) *MenuItem`

The `SetAccelerator()` method sets the accelerator (keyboard shortcut) for the menu item.

```go
menuItem.SetAccelerator("Ctrl+O")
```

