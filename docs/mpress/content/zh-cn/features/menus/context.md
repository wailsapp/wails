---
title: "上下文菜单"
description: "为应用程序创建右键上下文菜单"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## 问题

用户希望右键菜单能够提供与上下文相关的操作。不同的元素需要不同的菜单：

- **文本**：剪切、复制、粘贴
- **图像**：保存、复制、打开
- **自定义元素**：应用程序特定操作

手动构建上下文菜单需要处理鼠标事件、定位以及平台差异。

## Wails 的解决方案

Wails 使用 CSS 属性提供<strong>声明式上下文菜单</strong>。可将菜单与 HTML 元素关联、传递数据并处理点击事件，所有操作均采用平台原生行为。

![macOS 上显示在 WebView 上方的 Wails 自定义上下文菜单](/assets/screenshots/context-menu-macos.png)

菜单是平台原生菜单，而打开该菜单的元素仍属于 WebView。此 macOS 截图使用了下方示例中的自定义上下文菜单注册。

## 快速开始

**Go 代码：**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML：**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

<strong>就这么简单！</strong>右键单击文本区域即可显示自定义菜单。

## 创建上下文菜单

### 基本上下文菜单

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

<strong>菜单 ID：</strong>必须唯一，用于将菜单与 HTML 元素关联。

### 使用子菜单

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### 使用复选框和单选组

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

有关<strong>所有菜单项类型</strong>，请参阅[菜单参考](/features/menus/reference/)。

## 与 HTML 元素关联

使用 CSS 自定义属性附加上下文菜单：

### 基本关联

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**CSS 属性：**`--custom-contextmenu: <menu-id>`

### 使用上下文数据

将数据从 HTML 传递到 Go：

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Go 处理程序：**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**CSS 属性：**

- `--custom-contextmenu: <menu-id>` - 要显示的菜单
- `--custom-contextmenu-data: <data>` - 要传递给处理程序的数据

### 动态数据

在 JavaScript 中动态生成数据：

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### 多个元素使用同一菜单

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**同一菜单，每个元素使用不同的数据。**

## 上下文数据

### 访问上下文数据

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

<strong>数据类型：</strong>始终为`string`。请根据需要解析。

### 传递复杂数据

复杂数据请使用 JSON：

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Go 处理程序：**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="安全性"}
请<strong>始终验证</strong>来自前端的上下文数据。用户可以操纵 CSS 属性，因此应将这些数据视为不可信输入。

@end

### 验证示例

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## 默认上下文菜单

WebView 为标准操作（复制、粘贴、检查）提供内置上下文菜单。使用`--default-contextmenu`控制该菜单：

### 隐藏默认菜单

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

<strong>适用场景：</strong>不适合使用默认菜单的自定义 UI 元素。

### 显示默认菜单

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

<strong>适用场景：</strong>文本区域、输入字段、可编辑内容。

### 自动（智能）模式

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

<strong>默认行为。</strong>在以下情况下显示默认菜单：

- 选中了文本
- 位于文本输入字段中
- 位于可编辑内容中（`contenteditable`）

其他情况下隐藏默认菜单。

### 结合使用自定义菜单和默认菜单

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**行为：**

1. 首先显示自定义菜单
2. 如果自定义菜单为空或未找到，则显示默认菜单
3. 两者可以共存（取决于平台）

## 动态上下文菜单

根据应用程序状态更新菜单：

### 启用/禁用菜单项

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="始终调用 Update()"}
更改菜单状态后，请<strong>调用`contextMenu.Update()`</strong>。这在 Windows 上至关重要。

有关详细信息，请参阅[菜单参考](/features/menus/reference/#enabled-state)。

@end

### 更改标签

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### 重新构建菜单

如需进行重大更改，请重新构建整个菜单：

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## 平台行为

上下文菜单是<strong>平台原生的</strong>：

@tabs{sync-key="platform"}
[macOS]
**原生 macOS 上下文菜单：**

- 系统动画和过渡效果
- 右键单击 = Control+单击（自动）
- 适应系统外观（浅色/深色）
- 默认菜单中的标准文本操作
- 长菜单使用原生滚动

**macOS 惯例：**

- 菜单项采用句首字母大写格式
- 对打开对话框的菜单项使用省略号（...）
- 常用快捷键：⌘C（复制）、⌘V（粘贴）

[Windows]
**原生 Windows 上下文菜单：**

- Windows 原生样式
- 遵循 Windows 主题
- 默认菜单中的标准 Windows 操作
- 支持触摸和笔输入

**Windows 惯例：**

- 菜单项采用标题式大小写
- 对打开对话框的菜单项使用省略号（...）
- 常用快捷键：Ctrl+C（复制）、Ctrl+V（粘贴）

[Linux]
**桌面环境集成：**

- 适应桌面主题（GTK、Qt 等）
- 右键单击行为遵循系统设置
- 默认菜单内容因环境而异
- 定位方式遵循桌面环境惯例

**Linux 注意事项：**

- 在目标桌面环境中进行测试
- GTK 和 Qt 的行为不同
- 某些桌面环境会自定义上下文菜单

@end

## 完整示例

**Go 代码：**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML：**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## 最佳实践

### ✅ 应该做

- **让菜单保持专注**——仅提供与该元素相关的操作
- **验证上下文数据**——将其视为不可信输入
- **使用清晰的标签**——使用“删除文件”，而不是“删除”
- **调用 menu.Update()**——更改菜单状态后
- **在所有平台上进行测试**——行为因平台而异
- **提供键盘快捷键**——用于常用操作
- **将相关菜单项分组**——使用分隔符

### ❌ 不应该做

- **不要信任上下文数据**——始终进行验证
- **不要让菜单过长**——最多 7-10 个菜单项
- **不要忘记调用 menu.Update()**——否则菜单将无法正常工作
- **不要嵌套得过深**——最多 2 层
- **不要使用行话**——让标签对用户友好
- **不要阻塞处理程序**——确保处理程序快速运行

## 故障排除

### 上下文菜单未出现

**可能的原因：**

1. 菜单 ID 不匹配
2. CSS 属性拼写错误
3. 运行时尚未初始化

**解决方案：**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### 未收到上下文数据

**可能的原因：**

1. 未设置 CSS 属性
2. 数据包含特殊字符

**解决方案：**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

或者使用 JavaScript：

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### 菜单项无响应

<strong>原因：</strong>启用后忘记调用 `menu.Update()`

**解决方案：**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## 后续步骤

@cards{cols="2"}
📖 菜单参考
菜单项类型和属性的完整参考。

[了解更多 →](/features/menus/reference/)

---
☰ 应用程序菜单
创建应用程序菜单栏。

[了解更多 →](/features/menus/application/)

---
★ 系统托盘菜单
添加系统托盘/菜单栏集成。

[了解更多 →](/features/menus/systray/)

---
📖 菜单模式
常用菜单模式和最佳实践。

[了解更多 →](/guides/menus/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[上下文菜单示例](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus)。
