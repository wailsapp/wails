# Menu

### `type Menu struct`

The `Menu` struct holds information about a menu, including which items it contains and its label.

### `func NewMenu() *Menu`

This function initializes a new Menu.

### Add

API: `Add(label string) *MenuItem`

This method takes a `label` of type `string` as an input and adds a new `MenuItem` with the given label to the menu. It returns the `MenuItem` added.

### AddSeparator

API: `AddSeparator()`

This method adds a new separator `MenuItem` to the menu.

### AddCheckbox

API: `AddCheckbox(label string, enabled bool) *MenuItem`

This method takes a `label` of type `string` and `enabled` of type `bool` as inputs and adds a new checkbox `MenuItem` with the given label and enabled state to the menu. It returns the `MenuItem` added.

### AddRadio

API: `AddRadio(label string, enabled bool) *MenuItem`

This method takes a `label` of type `string` and `enabled` of type `bool` as inputs and adds a new radio `MenuItem` with the given label and enabled state to the menu. It returns the `MenuItem` added.

### Update

API: `Update()`

This method processes any radio groups and updates the menu if a menu implementation is not initialized.

### AddSubmenu

API: `AddSubmenu(s string) *Menu`

This method takes a `s` of type `string` as input and adds a new submenu `MenuItem` with the given label to the menu. It returns the submenu added.

### AddRole

API: `AddRole(role Role) *Menu`

This method takes `role` of type `Role` as input, adds it to the menu if it is not `nil` and returns the `Menu`.

### SetLabel

API: `SetLabel(label string)`

This method sets the `label` of the `Menu`.
