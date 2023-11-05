# 菜单

可以创建菜单并添加到应用程序中。它们可以用于创建上下文菜单、系统托盘菜单和应用程
序菜单。

要创建一个新菜单，请调用：

```go
    // 创建一个新菜单
    menu := app.NewMenu()
```

然后，`Menu` 类型上可用以下操作：

### 添加

API：`Add(label string) *MenuItem`

此方法以 `string` 类型的 `label` 作为输入，并将具有给定标签的新 `MenuItem` 添加
到菜单中。它返回添加的 `MenuItem`。

### 添加分隔符

API：`AddSeparator()`

此方法将一个新的分隔符 `MenuItem` 添加到菜单中。

### 添加复选框

API：`AddCheckbox(label string, enabled bool) *MenuItem`

此方法以 `string` 类型的 `label` 和 `bool` 类型的 `enabled` 作为输入，并将具有给
定标签和启用状态的新复选框 `MenuItem` 添加到菜单中。它返回添加的 `MenuItem`。

### 添加单选按钮

API：`AddRadio(label string, enabled bool) *MenuItem`

此方法以 `string` 类型的 `label` 和 `bool` 类型的 `enabled` 作为输入，并将具有给
定标签和启用状态的新单选按钮 `MenuItem` 添加到菜单中。它返回添加的 `MenuItem`。

### 更新

API：`Update()`

此方法处理任何单选按钮组，并在菜单未初始化时更新菜单。

### 添加子菜单

API：`AddSubmenu(s string) *Menu`

此方法以 `string` 类型的 `s` 作为输入，并将具有给定标签的新子菜单 `MenuItem` 添
加到菜单中。它返回添加的子菜单。

### 添加角色

API：`AddRole(role Role) *Menu`

此方法以 `Role` 类型的 `role` 作为输入，如果不为 `nil`，则将其添加到菜单中，并返
回 `Menu`。

### 设置标签

API：`SetLabel(label string)`

此方法设置 `Menu` 的 `label`。
