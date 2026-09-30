---
title: "檔案拖放"
description: "接受從作業系統拖曳到應用程式中的檔案"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails 可讓使用者將檔案從作業系統（檔案管理員、桌面）拖曳到應用程式中。HTML5 拖放只能在瀏覽器內運作，而此功能可讓您存取磁碟上實際的檔案路徑。

![macOS 上具有外部檔案放置區域的 Wails 拖放範例](/assets/screenshots/file-drop-macos.png)

外部放置區域是 WebView 的一部分，而 Wails 則提供原生作業系統檔案拖放事件。

## 啟用檔案拖放

檔案拖放預設為停用。若要啟用，請在視窗選項中設定`EnableFileDrop: true`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

當`EnableFileDrop`為`false`（預設值）時，從作業系統拖入的檔案會遭到封鎖：它們不會在 WebView 中開啟，也不會觸發任何事件。這可避免使用者將檔案拖曳到應用程式上方時意外導覽。

## 定義放置區域

放置區域會告訴 Wails 哪些元素應接受檔案。放在放置區域之外的檔案會被忽略。

將`data-file-drop-target`屬性新增至任何元素：

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

您可以設定多個放置區域。元素的`id`和 CSS 類別會傳遞至您的 Go 程式碼，讓您能依檔案放置的位置採取不同的處理方式。

## 設定拖曳懸停樣式

當檔案拖曳到放置區域上方時，Wails 會新增`file-drop-target-active`類別。您可藉此提供視覺回饋，讓使用者知道可將檔案放在哪裡：

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

當檔案離開該區域或被放下時，此類別會自動移除。

## 偵測放下的檔案

當檔案放在有效的放置區域時，Wails 會觸發`WindowFilesDropped`事件。事件內容包含所有放下檔案的完整檔案系統路徑：

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

這些路徑是絕對路徑，例如`/home/user/documents/report.pdf`或`C:\Users\Name\Documents\report.pdf`。

## 取得放置目標資訊

若有多個放置區域，您可以使用`DropTargetDetails()`確認是哪一個區域接收了檔案：

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

這可讓您將檔案路由至不同的處理常式：

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## 完整範例

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

## 全視窗拖放

若要讓檔案可放在應用程式中的任何位置，請將該屬性新增至 body 元素：

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

您可以使用 CSS 覆蓋層來表示整個視窗都是放置目標：

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

## 與 HTML 拖放搭配使用

您可以在同一個應用程式中同時使用外部檔案拖放和內部 HTML 拖放。當`EnableFileDrop`為`true`時，Wails 會攔截從外部拖入的檔案，但會讓內部 HTML5 拖曳照常通過。

若要在 HTML 放置區域處理常式中區分兩者，請檢查拖曳內容是否包含檔案：

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

這可確保 HTML 放置處理常式只回應內部拖曳（例如移動清單項目），而 Wails 則透過`WindowFilesDropped`事件另行處理外部檔案拖放。

## 後續步驟

- [HTML 拖放](/features/drag-and-drop/html/)－在應用程式內拖曳元素
- [視窗選項](/features/windows/options/)－所有視窗組態選項
