---
title: "菜单"
description: "在 Wails v3 中创建和自定义菜单的指南"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 提供了功能强大的菜单系统，可用于创建应用程序菜单和上下文菜单。本指南将介绍菜单系统的各种功能。

## 创建菜单

要创建新菜单，请使用 Menus 管理器的`New()`方法：

```go
menu := app.Menu.New()
```

### 添加菜单项

Wails 支持多种菜单项，每种都有特定用途：

#### 常规菜单项

常规菜单项是菜单的基本组成部分。它们显示文本，并可在点击时触发操作：

```go
menuItem := menu.Add("Click Me")
```

#### 复选框

复选框菜单项提供可切换的状态，适合用于启用或禁用功能或设置：

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### 单选组

单选组允许用户从一组互斥选项中选择一个。将单选菜单项相邻放置时，会自动创建单选组：

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### 分隔线

分隔线是用于将菜单项整理成多个逻辑分组的水平线：

```go
menu.AddSeparator()
```

#### 子菜单

子菜单是嵌套菜单，将鼠标悬停在菜单项上或点击菜单项时会显示。它们适合用于组织复杂的菜单结构：

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### 合并菜单

可以通过追加或前置的方式，将一个菜单添加到另一个菜单中。

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
默认情况下，`prepend`和`append`会与原菜单共享状态。如果要创建具有独立状态的新菜单，可以对该菜单调用`.Clone()`。

例如：`menu.Append(secondaryMenu.Clone())`

@end

#### 清空菜单

如果处理的菜单项数量不固定，在某些情况下，构建一个全新的菜单会更合适。

这会清除现有菜单中的所有菜单项，之后可以重新添加菜单项。

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
清空菜单只会清除顶层菜单项。虽然子菜单将不再可见，但仍会占用内存，因此请务必谨慎管理菜单。

@end

#### 销毁菜单

如果要清空并释放菜单，请使用`Destroy()`方法：

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### 菜单项属性

菜单项有多个可配置的属性：

| 属性 | 方法 | 说明 |
| --- | --- | --- |
| 标签 | `SetLabel(string)` | 设置显示文本 |
| 已启用 | `SetEnabled(bool)` | 启用或禁用菜单项 |
| 已选中 | `SetChecked(bool)` | 设置选中状态（用于复选框或单选菜单项） |
| 工具提示 | `SetTooltip(string)` | 设置工具提示文本 |
| 已隐藏 | `SetHidden(bool)` | 显示或隐藏菜单项 |
| 快捷键 | `SetAccelerator(string)` | 设置键盘快捷键 |

### 菜单项状态

菜单项可以处于不同状态，这些状态控制其可见性和交互能力：

#### 可见性

可以使用`SetHidden()`方法动态显示或隐藏菜单项：

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

隐藏的菜单项会从菜单中完全移除，直至再次显示。这适用于只应在应用程序处于特定状态时出现的上下文相关菜单项。

#### 启用状态

可以使用`SetEnabled()`方法启用或禁用菜单项：

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

禁用的菜单项仍然可见，但会呈灰色且无法点击。这通常用于表明某项操作当前不可用，例如：

- 没有可保存的更改时禁用“保存”
- 未选择任何内容时禁用“复制”
- 没有可撤销的操作时禁用“撤销”

#### 动态状态管理

可以将这些状态与事件处理程序结合起来，创建动态菜单：

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### 事件处理

菜单项可以使用`OnClick`方法响应点击事件：

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

上下文提供了被点击菜单项的相关信息：

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### 基于角色的菜单项

Wails 提供了一组预定义的菜单角色，可自动创建具备标准功能的菜单项。以下是支持的菜单角色：

#### 完整菜单结构

这些角色会创建包含常用功能的完整菜单结构：

| 角色 | 说明 | 平台说明 |
| --- | --- | --- |
| `AppMenu` | 包含“关于”“服务”“隐藏/显示”和“退出”的应用菜单 | 仅限 macOS |
| `EditMenu` | 包含“撤销”“重做”“剪切”“复制”“粘贴”等功能的标准“编辑”菜单 | 所有平台 |
| `ViewMenu` | 包含“重新加载”“缩放”和“全屏”控件的“视图”菜单 | 所有平台 |
| `WindowMenu` | 窗口控件（最小化、缩放等） | 所有平台 |
| `HelpMenu` | 包含指向 Wails 网站的“了解更多”链接的“帮助”菜单 | 所有平台 |

#### 单个菜单项

可以使用这些角色添加单个菜单项：

| 角色 | 说明 | 平台说明 |
| --- | --- | --- |
| `About` | 显示应用的“关于”对话框 | 所有平台 |
| `Hide` | 隐藏应用 | 仅限 macOS |
| `HideOthers` | 隐藏其他应用 | 仅限 macOS |
| `UnHide` | 显示隐藏的应用 | 仅限 macOS |
| `CloseWindow` | 关闭当前窗口 | 所有平台 |
| `Minimise` | 最小化窗口 | 所有平台 |
| `Zoom` | 缩放窗口 | 仅限 macOS |
| `Front` | 将窗口置于最前端 | 仅限 macOS |
| `Quit` | 退出应用 | 所有平台 |
| `Undo` | 撤销上一个操作 | 所有平台 |
| `Redo` | 重做上一个操作 | 所有平台 |
| `Cut` | 剪切所选内容 | 所有平台 |
| `Copy` | 复制所选内容 | 所有平台 |
| `Paste` | 从剪贴板粘贴 | 所有平台 |
| `PasteAndMatchStyle` | 粘贴并匹配样式 | 仅限 macOS |
| `SelectAll` | 全选 | 所有平台 |
| `Delete` | 删除所选内容 | 所有平台 |
| `Reload` | 重新加载当前页面 | 所有平台 |
| `ForceReload` | 强制重新加载当前页面 | 所有平台 |
| `ToggleFullscreen` | 切换全屏模式 | 所有平台 |
| `ResetZoom` | 重置缩放级别 | 所有平台 |
| `ZoomIn` | 放大 | 所有平台 |
| `ZoomOut` | 缩小 | 所有平台 |

以下示例展示了如何同时使用完整菜单和各个角色：

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## 应用程序菜单

应用程序菜单是显示在应用程序窗口顶部（Windows/Linux）或屏幕顶部（macOS）的菜单。

### 应用程序菜单行为

使用`app.Menu.Set()`设置应用程序菜单后，它会成为 macOS 上的主菜单。 在 Windows/Linux 上，菜单按窗口分别设置。

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

以下完整示例展示了这些不同的菜单行为：

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## 上下文菜单

上下文菜单是在右键单击应用程序中的元素时显示的弹出菜单。它们可用于快速访问与所单击元素相关的操作。

### 默认上下文菜单

默认上下文菜单是 WebView 内置的上下文菜单，可提供如下系统级操作：

- 用于文本操作的复制、剪切和粘贴
- 文本选择控件
- 拼写检查选项

#### 控制默认上下文菜单

可以使用`--default-contextmenu` CSS 属性控制何时显示默认上下文菜单：

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
只有在[前端运行时准备就绪](/reference/frontend-runtime/)后，此功能才能按预期工作。

@end

#### 嵌套上下文菜单的行为

在嵌套元素上使用`--default-contextmenu`属性时，适用以下规则：

1. 除非显式覆盖，否则子元素会继承其父元素的上下文菜单设置
2. 最具体（距离最近）的设置优先
3. 可以使用`auto`值恢复默认行为

嵌套上下文菜单行为示例：

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### 自定义上下文菜单

自定义上下文菜单可提供与所单击元素相关的应用程序专用操作。它们尤其适用于：

- 文档管理器中的文件操作
- 图像处理工具
- 数据网格中的自定义操作
- 组件特定操作

#### 创建自定义上下文菜单

创建自定义上下文菜单时，需要提供一个唯一标识符（名称），用于将菜单与 HTML 元素关联起来：

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

name 参数（本例中为“imageMenu”）用作唯一标识符，可用于：

1. 将 HTML 元素关联到此特定上下文菜单
2. 确定右键单击时应显示哪个菜单
3. 支持菜单更新和清理

#### 上下文数据

处理上下文菜单事件时，可以访问被单击的菜单项及其关联的上下文数据：

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

上下文数据从 HTML 元素的`--custom-contextmenu-data`属性传入，并可通过`ctx.ContextMenuData()`在单击处理程序中获取。这在以下情况中特别有用：

- 处理列表或网格，且其中每个项目都需要唯一标识
- 处理针对特定组件或元素的操作
- 将状态或元数据从前端传递到后端

#### 上下文菜单管理

更改上下文菜单后，调用`Update()`方法以应用更改：

```go
contextMenu.Update()
```

不再需要某个上下文菜单时，可以将其销毁：

```go
contextMenu.Destroy()
```

@note{type="danger" title="警告"}
调用`Destroy()`后，再次使用该上下文菜单引用将导致 panic。

@end

### 实际示例：图片库

下面是为图片库实现自定义上下文菜单的完整示例：

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

在此示例中：

1. 使用标识符“imageMenu”创建上下文菜单
2. 使用`--custom-contextmenu: imageMenu`将每个图片容器关联到该菜单
3. 每个容器都使用`--custom-contextmenu-data`将其图片 ID 作为上下文数据提供
4. 后端在单击处理程序中接收图片 ID，并可执行相应的特定操作
5. 所有图片复用同一个菜单，但上下文数据会指明要操作的图片

此模式特别适用于：

- 需要对各行执行特定操作的数据网格
- 需要对文件执行上下文相关操作的文件管理器
- 不同元素需要不同操作的设计工具
- 对多个实例应用相同操作的任何组件
