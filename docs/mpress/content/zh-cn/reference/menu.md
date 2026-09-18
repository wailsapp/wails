---
title: "菜单 API"
description: "菜单 API 完整参考"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## 概述

菜单 API 提供了创建和管理应用程序菜单、上下文菜单及系统托盘菜单的方法。

**菜单类型：**

- **应用程序菜单**——顶部菜单栏（文件、编辑等）
- **上下文菜单**——右键菜单
- **系统托盘菜单**——系统托盘/通知区域中的菜单

## 创建菜单

### NewMenu()

创建一个新菜单。

```go
func (a *App) NewMenu() *Menu
```

**示例：**

```go
menu := app.NewMenu()
```

## 菜单方法

### Add()

向菜单中添加一个菜单项。

```go
func (m *Menu) Add(label string) *MenuItem
```

**参数：**

- `label`——菜单项显示的文本

<strong>返回值：</strong>创建的菜单项

**示例：**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

向菜单中添加一个子菜单。

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**参数：**

- `label`——子菜单标签

<strong>返回值：</strong>创建的子菜单

**示例：**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

在菜单项之间添加一条可见的分隔线。

```go
func (m *Menu) AddSeparator()
```

**示例：**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

<strong>最佳实践：</strong>使用分隔线对相关菜单项进行分组。

### AddCheckbox()

添加一个可勾选的菜单项。

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**参数：**

- `label`——复选框标签
- `checked`——初始选中状态

**示例：**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

添加一个单选菜单项（互斥组）。

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**参数：**

- `label`——单选按钮标签
- `checked`——初始选中状态

**示例：**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

更新菜单以反映对菜单项所做的任何更改。

```go
func (m *Menu) Update()
```

**示例：**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

<strong>重要：</strong>修改菜单项属性后，务必调用`Update()`。

## 菜单项方法

### OnClick()

为菜单项注册点击处理程序。

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**参数：**

- `callback`——点击菜单项时调用的函数

<strong>返回值：</strong>该菜单项（用于链式调用）

**示例：**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

更改菜单项的标签。

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**示例：**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

启用或禁用菜单项。

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**示例：**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**常见模式：**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

设置复选框/单选菜单项的选中状态。

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**示例：**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

返回当前的选中状态。

```go
func (mi *MenuItem) Checked() bool
```

**示例：**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

为菜单项设置键盘快捷键。

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**参数：**

- `accelerator`——键盘快捷键（例如“Ctrl+S”、“Cmd+Q”）

**快捷键格式：**

- **修饰键：**`Ctrl`、`Cmd`、`Alt`、`Shift`
- **按键：**`A-Z`、`0-9`、`F1-F12`、`Enter`、`Backspace`等
- <strong>平台：</strong>在 macOS 上使用`Cmd`，在 Windows/Linux 上使用`Ctrl`

**示例：**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**平台感知示例：**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

设置将鼠标悬停在菜单项上时显示的工具提示。

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**示例：**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

显示或隐藏菜单项。

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**示例：**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## 应用程序菜单

### app.Menu.Set()

设置应用程序的主菜单栏。

```go
func (mm *MenuManager) Set(menu *Menu)
```

**示例：**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**平台说明：**

- <strong>macOS：</strong>菜单显示在屏幕顶部的菜单栏中
- <strong>Windows/Linux：</strong>菜单显示在窗口标题栏中
- <strong>macOS：</strong>自动添加以应用名称命名的应用程序菜单

## 上下文菜单

### app.ContextMenu.New() / app.ContextMenu.Add()

通过管理器创建`*ContextMenu`，并以指定名称注册它。`ContextMenuManager.Add`接受`*ContextMenu`，**而不是**`*Menu`，并且<strong>没有</strong>`app.RegisterContextMenu`方法。

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

也可以使用包级`application.NewContextMenu(name string) *ContextMenu`一步创建并注册上下文菜单。

**Go：**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS：**

当右键单击的目标元素（或其任意祖先元素）的`--custom-contextmenu` CSS 自定义属性设为某个已注册上下文菜单的名称时，运行时会触发该菜单。可选的`--custom-contextmenu-data`属性会通过`ctx.ContextMenuData()`传递给 Go 回调。要禁用浏览器的默认上下文菜单，请设置`--default-contextmenu: hide`（或`auto`/`show`）。

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

不存在`data-wails-context-menu="..."`属性——运行时从未接入过该属性。

**动态上下文菜单：**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## 系统托盘菜单

### app.SystemTray.New()

创建新的系统托盘图标。

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**示例：**

```go
tray := app.SystemTray.New()
```

### SetIcon()

设置系统托盘图标。

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**示例：**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

设置系统托盘菜单。

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**示例：**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

设置将鼠标悬停在托盘图标上时显示的工具提示。此方法不返回任何内容。

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**示例：**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

处理对托盘图标的左键单击。

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**示例：**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## 完整示例

### 标准应用程序菜单

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### 系统托盘应用程序

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### 动态更新菜单

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## 最佳实践

### ✅ 推荐做法

- **使用标准快捷键**——遵循平台惯例（例如使用 Ctrl+C 复制）
- **更改后调用 Update()**——否则菜单不会反映这些更改
- **将相关项目分组**——使用分隔线组织菜单项
- **禁用不可用的操作**——不要隐藏，应使用 SetEnabled(false) 将其禁用
- **使用清晰的标签**——保持简洁且含义明确
- **遵循平台惯例**——注意 macOS 与 Windows/Linux 的菜单模式差异

### ❌ 避免的做法

- **不要忘记调用 Update()**——这是最常见的错误
- **嵌套不要过深**——菜单最多保持2-3层
- **不要使用含义模糊的标签**——例如“处理”与“处理文档”
- **不要过度复杂化**——保持菜单简洁且重点明确
- **不要混用不同的表达方式**——保持命名和组织方式一致

## 平台特定说明

### macOS

- 自动添加以应用名称命名的应用程序菜单
- 快捷键应使用`Cmd`，而不是`Ctrl`
- 默认情况下，应用程序菜单中包含“关于”、“偏好设置”和“退出”

### Windows/Linux

- 不会自动创建应用程序菜单
- 使用`Ctrl`作为快捷键修饰键
- “退出”通常位于“文件”菜单中
