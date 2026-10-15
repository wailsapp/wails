---
title: "Перетаскивание файлов"
description: "Принимайте в приложении файлы, перетаскиваемые из операционной системы"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails позволяет пользователям перетаскивать в приложение файлы из операционной системы — например, из файлового менеджера или с рабочего стола. В отличие от перетаскивания HTML5, которое работает только в браузере, этот механизм предоставляет доступ к фактическим путям к файлам на диске.

![Пример перетаскивания в Wails с внешней зоной приёма файлов в macOS](/assets/screenshots/file-drop-macos.png)

Внешняя зона приёма является частью веб-представления, а Wails предоставляет нативные события операционной системы о перетаскивании файлов.

## Включение приёма файлов

По умолчанию приём файлов отключён. Чтобы включить его, задайте `EnableFileDrop: true` в параметрах окна:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

Если `EnableFileDrop` имеет значение `false` (по умолчанию), перетаскиваемые из ОС файлы блокируются: они не открываются в веб-представлении и не вызывают никаких событий. Это предотвращает случайную навигацию, когда пользователи перетаскивают файлы над приложением.

## Определение зон приёма

Зоны приёма указывают Wails, какие элементы должны принимать файлы. Файлы, сброшенные за пределами зоны приёма, игнорируются.

Добавьте атрибут `data-file-drop-target` к любому элементу:

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

Можно создать несколько зон приёма. Значение `id` элемента и его классы CSS передаются в код Go, поэтому файлы можно обрабатывать по-разному в зависимости от того, куда они были сброшены.

## Оформление при наведении перетаскиваемого файла

Когда файлы перетаскивают над зоной приёма, Wails добавляет класс `file-drop-target-active`. С его помощью можно предоставить визуальную обратную связь, чтобы пользователи понимали, куда можно сбросить файлы:

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

Класс удаляется автоматически, когда файлы покидают зону или сбрасываются в неё.

## Обнаружение сброшенных файлов

Когда файлы сбрасывают в допустимую зону приёма, Wails инициирует событие `WindowFilesDropped`. Контекст события содержит полные пути в файловой системе ко всем сброшенным файлам:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

Пути являются абсолютными, например `/home/user/documents/report.pdf` или `C:\Users\Name\Documents\report.pdf`.

## Получение сведений о целевой зоне

Если у вас несколько зон приёма, с помощью `DropTargetDetails()` можно определить, какая из них получила файлы:

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

Это позволяет направлять файлы разным обработчикам:

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## Полный пример

**Go:**

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

**HTML:**

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

## Приём файлов во всём окне

Чтобы файлы можно было сбрасывать в любой области приложения, добавьте атрибут к элементу body:

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

С помощью наложения CSS можно показать, что целевой зоной приёма является всё окно:

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

## Совместное использование с перетаскиванием HTML

В одном приложении можно использовать как перетаскивание внешних файлов, так и внутреннее перетаскивание HTML. Если `EnableFileDrop` имеет значение `true`, Wails перехватывает перетаскивание внешних файлов, но пропускает внутренние операции перетаскивания HTML5 в обычном режиме.

Чтобы различать эти операции в обработчиках зоны приёма HTML, проверьте, содержит ли перетаскиваемый объект файлы:

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

Благодаря этому обработчики сброса HTML реагируют только на внутреннее перетаскивание (например, перемещение элементов списка), а Wails отдельно обрабатывает сбрасывание внешних файлов с помощью события `WindowFilesDropped`.

## Дальнейшие действия

- [Перетаскивание HTML](/features/drag-and-drop/html/) — перетаскивание элементов внутри приложения
- [Параметры окна](/features/windows/options/) — все параметры конфигурации окна
