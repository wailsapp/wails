---
title: "应用程序菜单"
description: "为桌面应用程序创建原生菜单栏"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## 问题

专业桌面应用程序需要菜单栏，例如“文件”“编辑”“视图”和“帮助”。但菜单在不同平台上的工作方式各不相同：

- **macOS**：位于屏幕顶部的全局菜单栏
- **Windows**：位于窗口标题栏中的菜单栏
- **Linux**：因桌面环境而异

手动构建符合各平台惯例的菜单既繁琐又容易出错。

## Wails 解决方案

Wails 提供一个<strong>统一 API</strong>，可自动创建平台原生菜单。只需编写一次，即可在所有平台上获得原生行为。

![macOS 上的 Wails 应用程序菜单，包含标准项、复选项、单选项和子菜单项](/assets/screenshots/application-menu-macos.png)

在 macOS 上，应用程序菜单位于全局菜单栏中。此截图展示了由 Wails 菜单 API 渲染的原生菜单，其中包含禁用项、复选项、单选项和子菜单项。

## 快速开始

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

<strong>就这么简单！</strong>现在，你已经拥有包含标准菜单项的平台原生菜单。`UseApplicationMenu`选项可确保 Windows 和 Linux 窗口无需额外代码即可显示菜单。

## 创建菜单

### 创建基本菜单

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### 设置菜单

**推荐方法**——使用`UseApplicationMenu`确保跨平台一致性：

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

此方法的行为如下：

- 在<strong>macOS</strong>上：菜单显示在屏幕顶部（标准行为）
- 在<strong>Windows/Linux</strong>上：每个启用`UseApplicationMenu: true`的窗口都会显示应用程序菜单

**各平台的具体情况：**

@tabs{sync-key="platform"}
[macOS]
**全局菜单栏**（每个应用程序一个）：

```go
app.Menu.Set(menu)
```

菜单显示在屏幕顶部，即使所有窗口都已关闭也会保留。由于所有应用程序都使用全局菜单，因此`UseApplicationMenu`选项在 macOS 上不起作用。

[Windows]
**每窗口菜单栏**：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

每个窗口都可以拥有自己的菜单，也可以继承应用程序菜单。菜单显示在窗口的标题栏中。

[Linux]
**每窗口菜单栏**（通常如此）：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

具体行为因桌面环境而异。某些桌面环境（如 Unity）支持全局菜单。

@end

@note{type="tip" title="简化跨平台菜单"}
使用`UseApplicationMenu: true`后，无需再编写如下平台特定代码：

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**每窗口自定义菜单：**

如果某个窗口需要使用不同于应用程序菜单的菜单，请直接为其设置：

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## 菜单角色

Wails 提供<strong>预定义菜单角色</strong>，可自动创建符合各平台惯例的菜单结构。

### 可用角色

| 角色 | 说明 | 平台说明 |
| --- | --- | --- |
| `AppMenu` | 包含“关于”“偏好设置”和“退出”的应用程序菜单 | **仅限 macOS** |
| `FileMenu` | 文件操作（新建、打开、保存等） | 所有平台 |
| `EditMenu` | 文本编辑（撤销、重做、剪切、复制、粘贴） | 所有平台 |
| `WindowMenu` | 窗口管理（最小化、缩放等） | 所有平台 |
| `HelpMenu` | 帮助和信息 | 所有平台 |

### 使用角色

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**你将获得：**

@tabs{sync-key="platform"}
[macOS]
**AppMenu**（包含应用名称）：

- 关于 [应用名称]
- 偏好设置... (⌘,)
- ---
- 服务
- ---
- 隐藏 [应用名称] (⌘H)
- 隐藏其他应用 (⌥⌘H)
- 全部显示
- ---
- 退出 [应用名称] (⌘Q)

**FileMenu**：

- 新建 (⌘N)
- 打开... (⌘O)
- ---
- 关闭窗口 (⌘W)

**EditMenu**：

- 撤销 (⌘Z)
- 重做 (⇧⌘Z)
- ---
- 剪切 (⌘X)
- 复制 (⌘C)
- 粘贴 (⌘V)
- 全选 (⌘A)

**WindowMenu**：

- 最小化 (⌘M)
- 缩放
- ---
- 全部置于最前

**HelpMenu**：

- [应用名称]帮助

[Windows]
**FileMenu**：

- 新建 (Ctrl+N)
- 打开... (Ctrl+O)
- ---
- 退出 (Alt+F4)

**EditMenu**：

- 撤销 (Ctrl+Z)
- 重做 (Ctrl+Y)
- ---
- 剪切 (Ctrl+X)
- 复制 (Ctrl+C)
- 粘贴 (Ctrl+V)
- 全选 (Ctrl+A)

**WindowMenu**：

- 最小化
- 最大化

**HelpMenu**：

- 关于[应用名称]

[Linux]
与 Windows 类似，但键盘快捷键可能因桌面环境而异。

@end

### 自定义角色菜单

`Menu.AddRole(role)`返回<strong>接收者</strong>菜单（顶级菜单），<strong>而不是</strong>该角色的子菜单。要向角色的子菜单添加项目，请使用`FindByRole`查找已插入的角色项目，然后对其调用`GetSubmenu()`：

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## 自定义菜单

为应用程序特有的功能创建自己的菜单：

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

**有关更多菜单项类型**，请参阅[菜单参考](/features/menus/reference/)。

## 动态菜单

根据应用程序状态更新菜单：

### 启用/禁用项目

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="始终调用 menu.Update()"}
更改菜单状态（启用/禁用、标签、选中状态）后，**始终调用`menu.Update()`**。这一点在会重建菜单的 Windows 上尤为重要。

详情请参阅[菜单参考](/features/menus/reference/#enabled-state)。

@end

### 更改标签

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### 重建菜单

如需进行重大更改，请重建整个菜单：

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## 通过菜单控制窗口

菜单项可以控制窗口：

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**获取活动窗口：**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## 平台特定注意事项

### macOS

**菜单栏行为：**

- 显示在<strong>屏幕顶部</strong>（全局）
- 关闭所有窗口后仍然保留
- 第一个菜单<strong>始终是应用程序菜单</strong>
- 使用`menu.AddRole(application.AppMenu)`添加标准项目

**标准位置：**

- **关于**：应用程序菜单
- **偏好设置**：应用程序菜单 (⌘,)
- **退出**：应用程序菜单 (⌘Q)
- **帮助**：帮助菜单

**示例：**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**菜单栏行为：**

- 显示在<strong>窗口标题栏</strong>中
- 每个窗口都有自己的菜单
- 没有应用程序菜单

**标准位置：**

- **退出**：文件菜单 (Alt+F4)
- **设置**：工具或编辑菜单
- **关于**：帮助菜单

**示例：**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**菜单栏行为：**

- 通常每个窗口各有一个菜单（与 Windows 类似）
- 某些桌面环境支持全局菜单（Unity、安装扩展的 GNOME）
- 外观因桌面环境而异

<strong>最佳实践：</strong>遵循 Windows 惯例，并在目标桌面环境中测试。

## 完整示例

以下是一个可用于生产环境的菜单结构：

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## 最佳实践

### ✅ 应该做

- 对标准菜单（文件、编辑等）**使用菜单角色**
- 菜单结构应<strong>遵循平台惯例</strong>
- 为常用操作<strong>添加键盘快捷键</strong>
- 更改菜单状态后，**调用 menu.Update()**
- **在所有平台上进行测试**——行为会有所不同
- **保持菜单层级简洁**——最多 2-3 层
- **使用清晰的标签**——使用“保存项目”，而不是“保存”

### ❌ 不应该做

- **不要硬编码特定平台的快捷键**——请使用 `CmdOrCtrl`
- **不要在 macOS 的“文件”菜单中放置“退出”**——它应位于“应用程序”菜单中
- **不要在 macOS 的“帮助”菜单中放置“关于”**——它应位于“应用程序”菜单中
- **不要忘记调用 menu.Update()**——否则菜单将无法正常工作
- **不要嵌套得太深**——用户容易迷失
- **不要使用行话**——确保标签对用户友好

## 后续步骤

@cards{cols="2"}
📖 菜单参考
菜单项类型和属性的完整参考。

[了解更多 →](/features/menus/reference/)

---
◆ 上下文菜单
创建右键上下文菜单。

[了解更多 →](/features/menus/context/)

---
★ 系统托盘菜单
添加系统托盘/菜单栏集成。

[了解更多 →](/features/menus/systray/)

---
📖 菜单模式
常见的菜单模式和最佳实践。

[了解更多 →](/guides/menus/)

@end

---

<strong>有疑问？</strong>请在 [Discord](https://discord.gg/JDdSxwjhGf) 中提问，或查看[菜单示例](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)。
