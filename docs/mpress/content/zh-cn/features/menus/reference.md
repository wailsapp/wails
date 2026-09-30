---
title: "菜单参考"
description: "菜单项类型、属性和方法的完整参考"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## 菜单参考

菜单项类型、属性和动态行为的完整参考。使用复选框、单选组、分隔符和动态更新构建专业且响应迅速的菜单。

## 菜单项类型

### 常规菜单项

最常见的类型——显示文本并触发操作：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

<strong>适用于：</strong>命令、操作、打开窗口

### 复选框

可切换选中/未选中状态的菜单项：

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

<strong>适用于：</strong>布尔设置、功能开关、视图选项

<strong>重要：</strong>单击时，选中状态会自动切换。

### 单选组

互斥选项——只能选择其中一个：

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

<strong>适用于：</strong>互斥选项（大小、主题、模式）

**分组方式：**

- 相邻的单选菜单项会自动组成一组
- 选择其中一项会取消选中组内其他项
- 使用分隔符或常规菜单项分隔不同的组

**包含多个组的示例：**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### 子菜单

用于组织菜单项的嵌套菜单结构：

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

<strong>适用于：</strong>对相关菜单项分组、减少杂乱

<strong>嵌套限制：</strong>大多数平台支持2-3层。请避免更深层级的嵌套。

### 分隔符

菜单项之间的视觉分隔线：

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

<strong>适用于：</strong>在视觉上对相关菜单项分组

<strong>最佳实践：</strong>不要在菜单开头或末尾放置分隔符。

## 菜单项属性

### 标签

菜单项显示的文本：

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**动态标签：**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### 启用状态

控制菜单项是否可交互：

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Windows 菜单行为"}
在 Windows 上，菜单状态发生变化时需要重新构建。启用或禁用菜单项后，**务必调用`menu.Update()`**，尤其是该菜单项创建时处于禁用状态的情况。

<strong>原因：</strong>Windows 菜单更新时会从头重新构建。如果不调用`Update()`，单击处理程序将无法正常触发。

@end

**示例：动态启用/禁用**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**常见模式：满足条件时启用**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### 选中状态

对于复选框和单选菜单项，可控制或查询其选中状态：

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

<strong>自动切换：</strong>单击复选框时，其状态会自动切换。无需在单击处理程序中调用`SetChecked()`。

**手动控制：**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### 加速键（键盘快捷键）

为菜单项添加键盘快捷键：

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**加速键格式：**

- `CmdOrCtrl`——在 macOS 上为 Cmd，在 Windows/Linux 上为 Ctrl
- `Shift`、`Alt`、`Option`——修饰键
- `A-Z`、`0-9`——字母键/数字键
- `F1-F12`——功能键
- `Enter`、`Space`、`Backspace`等——特殊键

**示例：**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**平台特定的加速键：**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### 工具提示

为菜单项添加悬停文本（平台支持情况各异）：

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**平台支持情况：**

- **Windows：**✅ 支持
- **macOS：**❌ 不支持（工具提示并非菜单的标准功能）
- **Linux：**⚠️ 因桌面环境而异

### 隐藏状态

隐藏菜单项而不将其移除：

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

<strong>适用于：</strong>调试选项、功能标志、条件性功能

## 事件处理

### OnClick 处理程序

单击菜单项时执行代码：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**上下文提供：**

- `ctx.ClickedMenuItem()` — 被单击的菜单项
- 窗口上下文（如果来自窗口菜单）
- 应用程序上下文

**示例：在处理程序中访问菜单项**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### 多个处理程序

可以设置多个处理程序（最后设置的生效）：

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

<strong>最佳实践：</strong>只设置一次处理程序，并根据需要在其中使用条件逻辑。

## 动态菜单

### 更新菜单项

<strong>黄金法则：</strong>更改菜单状态后，始终调用`menu.Update()`。

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**这为何很重要：**

- <strong>Windows：</strong>更新菜单时会重新构建菜单
- <strong>macOS/Linux：</strong>重要性较低，但仍建议这样做
- <strong>单击处理程序：</strong>如果不调用 Update()，将无法正常触发

### 重新构建菜单

如需进行重大更改，请重新构建整个菜单：

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**何时重新构建：**

- 最近使用的文件列表发生变化
- 插件菜单发生变化
- 重大状态转换

**何时更新：**

- 启用或禁用菜单项
- 更改标签
- 切换复选框状态

### 上下文相关菜单

根据应用程序状态调整菜单：

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## 平台差异

### 菜单栏位置

| 平台 | 位置 | 说明 |
| --- | --- | --- |
| **macOS** | 屏幕顶部 | 全局菜单栏 |
| **Windows** | 窗口顶部 | 每个窗口各有一个菜单 |
| **Linux** | 窗口顶部 | 每个窗口各有一个（通常如此） |

### 标准菜单

**macOS：**

- 有“应用程序”菜单（以应用名称命名）
- “偏好设置”位于“应用程序”菜单中
- “退出”位于“应用程序”菜单中

**Windows/Linux：**

- 没有“应用程序”菜单
- “偏好设置”位于“编辑”或“工具”菜单中
- “退出”位于“文件”菜单中

**示例：符合平台惯例的结构**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### 快捷键惯例

**macOS：**

- 大多数快捷键使用`Cmd+`
- “偏好设置”使用`Cmd+,`
- “退出”使用`Cmd+Q`

**Windows：**

- 大多数快捷键使用`Ctrl+`
- “偏好设置”使用`Ctrl+P`或`Ctrl+,`
- “退出”使用`Alt+F4`（或`Ctrl+Q`）

**Linux：**

- 通常遵循 Windows 惯例
- 桌面环境可能会覆盖这些设置

## 最佳实践

### ✅ 推荐做法

- 更改菜单状态后（尤其是在 Windows 上），**调用 menu.Update()**
- 对互斥选项<strong>使用单选组</strong>
- 对可切换的功能<strong>使用复选框</strong>
- **为常用操作添加快捷键**
- **使用分隔符对相关项目进行分组**
- **在所有平台上测试**——行为因平台而异

### ❌ 不要这样做

- **不要忘记调用 menu.Update()**——否则点击处理程序将无法正常工作
- **不要嵌套得太深**——最多 2-3 层
- **不要以分隔符开头或结尾**——这样显得不专业
- **不要在 macOS 上使用工具提示**——不受支持
- **不要硬编码特定平台的快捷键**——请使用`CmdOrCtrl`

## 故障排除

### 菜单项无响应

<strong>症状：</strong>点击处理程序未触发

<strong>原因：</strong>启用菜单项后忘记调用`menu.Update()`

**解决方案：**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### 菜单项呈灰色

<strong>症状：</strong>无法点击菜单项

<strong>原因：</strong>菜单项已被禁用

**解决方案：**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### 快捷键不起作用

<strong>症状：</strong>键盘快捷键无法触发菜单项

**原因：**

1. 快捷键格式不正确
2. 与系统快捷键冲突
3. 窗口未获得焦点

**解决方案：**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## 后续步骤

- [应用程序菜单](/features/menus/application/)——创建应用程序菜单栏
- [上下文菜单](/features/menus/context/)——右键上下文菜单
- [系统托盘菜单](/features/menus/systray/)——系统托盘/菜单栏菜单
- [菜单模式](/guides/menus/)——常用菜单模式和最佳实践

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[菜单示例](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)。
