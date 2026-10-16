---
title: "文件拖放"
description: "接收从操作系统拖入应用程序的文件"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails 允许用户将文件从操作系统（文件管理器、桌面）拖入应用程序。HTML5 拖放只能在浏览器内部使用，而此功能让你能够访问磁盘上文件的实际路径。

![macOS 上带有外部文件放置区的 Wails 拖放示例](/assets/screenshots/file-drop-macos.png)

外部放置区是 WebView 的一部分，而 Wails 则提供原生操作系统文件拖放事件。

## 启用文件拖放

文件拖放默认处于禁用状态。要启用此功能，请在窗口选项中设置`EnableFileDrop: true`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

当`EnableFileDrop`为`false`（默认值）时，从操作系统拖入的文件会被阻止：它们既不会在 WebView 中打开，也不会触发任何事件。这可以防止用户将文件拖到应用程序上时意外发生页面导航。

## 定义放置区

放置区用于告知 Wails 哪些元素应接收文件。在放置区之外放下的文件会被忽略。

将`data-file-drop-target`属性添加到任意元素：

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

你可以设置多个放置区。元素的`id`和 CSS 类会传递给 Go 代码，因此可以根据文件的放置位置采用不同的处理方式。

## 设置拖入悬停样式

当文件被拖到放置区上方时，Wails 会添加`file-drop-target-active`类。你可以借此提供视觉反馈，让用户知道可在何处放下文件：

```css
.drop-zone {
    border: 2px dashed #ccc;
    padding: 40px;
    text-align: center;
    transition: all 0.2s ease;
}

.drop-zone.file-drop-target-active {
    border-color: #007bff;
    background-color: rgba(0, 123, 255, 0.1);
}
```

当文件离开放置区或被放下时，该类会自动移除。

## 检测已放下的文件

当文件被放到有效的放置区时，Wails 会触发`WindowFilesDropped`事件。事件上下文包含所有已放下文件的完整文件系统路径：

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

这些路径是绝对路径，例如`/home/user/documents/report.pdf`或`C:\Users\Name\Documents\report.pdf`。

## 获取放置目标信息

如果设置了多个放置区，可以使用`DropTargetDetails()`确定哪个放置区接收了文件：

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

这样便可将文件分派给不同的处理程序：

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## 完整示例

**Go：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    // Send to frontend
    app.Event.Emit("files-dropped", map[string]any{
        "files":   files,
        "target":  details.ElementID,
    })
})
```

**HTML：**

```html
<div id="images" class="drop-zone" data-file-drop-target>
    Drop images here
</div>

<div id="documents" class="drop-zone" data-file-drop-target>
    Drop documents here
</div>

<style>
    .drop-zone {
        border: 2px dashed #ccc;
        border-radius: 8px;
        padding: 40px;
        text-align: center;
        margin: 20px;
        transition: all 0.2s ease;
    }
    
    .drop-zone.file-drop-target-active {
        border-color: #007bff;
        background-color: rgba(0, 123, 255, 0.1);
    }
</style>
```

## 在整个窗口中放置

如果希望文件可放置在应用程序中的任意位置，请将该属性添加到 body 元素：

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

可以使用 CSS 覆盖层来指示整个窗口都是放置目标：

```css
body.file-drop-target-active::after {
    content: "Drop files anywhere";
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: #007bff;
    background: rgba(255, 255, 255, 0.9);
    pointer-events: none;
}
```

## 与 HTML 拖放结合使用

可以在同一个应用程序中同时使用外部文件拖放和内部 HTML 拖放。当`EnableFileDrop`为`true`时，Wails 会拦截从外部拖入的文件，但允许内部 HTML5 拖放正常通过。

要在 HTML 放置区处理程序中区分两者，请检查拖放内容是否包含文件：

```javascript
zone.addEventListener('dragenter', (e) => {
    // Skip external file drags - Wails handles these
    if (e.dataTransfer?.types.includes('Files')) {
        return;
    }
    // Handle internal HTML5 drags
    zone.classList.add('drag-over');
});

zone.addEventListener('drop', (e) => {
    // Skip external file drops - Wails handles these
    if (e.dataTransfer?.types.includes('Files')) {
        return;
    }
    e.preventDefault();
    zone.classList.remove('drag-over');
    // Handle internal drop
});
```

这样可确保 HTML 放置处理程序仅响应内部拖放（例如移动列表项），而 Wails 则通过`WindowFilesDropped`事件单独处理外部文件拖放。

## 后续步骤

- [HTML 拖放](/features/drag-and-drop/html/)——在应用程序内拖动元素
- [窗口选项](/features/windows/options/)——所有窗口配置选项
