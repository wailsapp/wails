---
title: "ファイルのドロップ"
description: "オペレーティングシステムからアプリケーションへドラッグされたファイルを受け入れる"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails では、ユーザーがオペレーティングシステム（ファイルマネージャー、デスクトップ）からアプリケーションへファイルをドラッグできます。ブラウザー内でしか機能しない HTML5 のドラッグ＆ドロップとは異なり、ディスク上にある実際のファイルパスへアクセスできます。

![macOS 上の外部ファイル用ドロップゾーンを使用した Wails のドラッグ＆ドロップの例](/assets/screenshots/file-drop-macos.png)

外部ファイル用ドロップゾーンは webview の一部であり、Wails はオペレーティングシステムのネイティブなファイルドロップイベントを提供します。

## ファイルドロップを有効にする

ファイルドロップはデフォルトで無効です。有効にするには、ウィンドウオプションで `EnableFileDrop: true` を設定します。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

`EnableFileDrop` が `false`（デフォルト）の場合、OS からドラッグされたファイルはブロックされます。ファイルが webview で開かれたり、イベントが発生したりすることはありません。これにより、ユーザーがアプリ上へファイルをドラッグした際に、意図せずページが移動することを防止できます。

## ドロップゾーンを定義する

ドロップゾーンは、ファイルを受け入れる要素を Wails に示します。ドロップゾーン外にドロップされたファイルは無視されます。

任意の要素に `data-file-drop-target` 属性を追加します。

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

複数のドロップゾーンを設けることができます。要素の `id` と CSS クラスは Go コードに渡されるため、ファイルがドロップされた場所に応じて処理を変えられます。

## ドラッグ中のホバーをスタイル設定する

ファイルがドロップゾーン上へドラッグされると、Wails は `file-drop-target-active` クラスを追加します。これを使用して視覚的なフィードバックを表示し、ドロップ可能な場所をユーザーに示せます。

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

ファイルがゾーンから離れるかドロップされると、このクラスは自動的に削除されます。

## ドロップされたファイルを検出する

有効なドロップゾーンにファイルがドロップされると、Wails は `WindowFilesDropped` イベントを発生させます。イベントコンテキストには、ドロップされたすべてのファイルの完全なファイルシステムパスが含まれます。

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

パスは、`/home/user/documents/report.pdf` や `C:\Users\Name\Documents\report.pdf` のような絶対パスです。

## ドロップ先の情報を取得する

複数のドロップゾーンがある場合は、`DropTargetDetails()` を使用して、どのゾーンがファイルを受け取ったかを確認できます。

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

これにより、ファイルを別々のハンドラーへ振り分けられます。

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## 完全な例

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

## ウィンドウ全体へのドロップ

アプリ内のどこにでもファイルをドロップできるようにするには、body 要素にこの属性を追加します。

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

CSS オーバーレイを使用して、ウィンドウ全体がドロップ先であることを示せます。

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

## HTML ドラッグ＆ドロップとの併用

同じアプリケーションで、外部ファイルのドロップと内部 HTML ドラッグ＆ドロップの両方を使用できます。`EnableFileDrop` が `true` の場合、Wails は外部ファイルのドラッグをインターセプトしますが、内部の HTML5 ドラッグは通常どおり通過させます。

HTML ドロップゾーンのハンドラーで両者を区別するには、ドラッグにファイルが含まれているかどうかを確認します。

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

これにより、HTML のドロップハンドラーは内部のドラッグ（リスト項目の移動など）にのみ応答し、外部ファイルのドロップは Wails が `WindowFilesDropped` イベントを介して別途処理します。

## 次のステップ

- [HTML ドラッグ＆ドロップ](/features/drag-and-drop/html/) — アプリ内で要素をドラッグする
- [ウィンドウオプション](/features/windows/options/) — ウィンドウのすべての設定オプション
