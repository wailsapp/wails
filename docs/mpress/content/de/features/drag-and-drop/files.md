---
title: "Dateiablage"
description: "Dateien annehmen, die aus dem Betriebssystem in Ihre Anwendung gezogen werden"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Mit Wails können Benutzer Dateien aus dem Betriebssystem (Dateimanager, Desktop) in Ihre Anwendung ziehen. Anders als Drag-and-drop in HTML5, das nur innerhalb des Browsers funktioniert, erhalten Sie dadurch Zugriff auf die tatsächlichen Dateipfade auf dem Datenträger.

![Ein Wails-Beispiel für Drag-and-drop mit einer Ablagezone für externe Dateien unter macOS](/assets/screenshots/file-drop-macos.png)

Die externe Ablagezone ist Teil des Webviews, während Wails die nativen Dateiablagereignisse des Betriebssystems bereitstellt.

## Dateiablage aktivieren

Die Dateiablage ist standardmäßig deaktiviert. Um sie zu aktivieren, legen Sie in Ihren Fensteroptionen `EnableFileDrop: true` fest:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

Wenn `EnableFileDrop` auf `false` gesetzt ist (Standardeinstellung), werden aus dem Betriebssystem gezogene Dateien blockiert: Sie werden weder im Webview geöffnet noch lösen sie Ereignisse aus. Dies verhindert eine versehentliche Navigation, wenn Benutzer Dateien über Ihre Anwendung ziehen.

## Ablagezonen definieren

Ablagezonen teilen Wails mit, welche Elemente Dateien annehmen sollen. Dateien, die außerhalb einer Ablagezone abgelegt werden, werden ignoriert.

Fügen Sie einem beliebigen Element das Attribut `data-file-drop-target` hinzu:

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

Sie können mehrere Ablagezonen verwenden. `id` und die CSS-Klassen des Elements werden an Ihren Go-Code übergeben, sodass Sie Dateiablagen abhängig vom Ablageort unterschiedlich behandeln können.

## Darstellung beim Darüberziehen gestalten

Wenn Dateien über eine Ablagezone gezogen werden, fügt Wails die Klasse `file-drop-target-active` hinzu. So können Sie visuelles Feedback geben, damit Benutzer erkennen, wo sie Dateien ablegen können:

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

Die Klasse wird automatisch entfernt, wenn die Dateien die Zone verlassen oder abgelegt werden.

## Abgelegte Dateien erkennen

Wenn Dateien in einer gültigen Ablagezone abgelegt werden, löst Wails ein `WindowFilesDropped`-Ereignis aus. Der Ereigniskontext enthält die vollständigen Dateisystempfade aller abgelegten Dateien:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

Die Pfade sind absolut, beispielsweise `/home/user/documents/report.pdf` oder `C:\Users\Name\Documents\report.pdf`.

## Informationen zum Ablageziel abrufen

Wenn Sie mehrere Ablagezonen verwenden, können Sie mit `DropTargetDetails()` ermitteln, welche davon die Dateien empfangen hat:

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

So können Sie Dateien an unterschiedliche Handler weiterleiten:

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## Vollständiges Beispiel

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

## Dateiablage im gesamten Fenster

Wenn Dateien überall in Ihrer Anwendung abgelegt werden können sollen, fügen Sie das Attribut dem body-Element hinzu:

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

Mit einem CSS-Overlay können Sie anzeigen, dass das gesamte Fenster als Ablageziel dient:

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

## Mit HTML-Drag-and-drop kombinieren

Sie können in derselben Anwendung sowohl externe Dateiablagen als auch internes HTML-Drag-and-drop verwenden. Wenn `EnableFileDrop` auf `true` gesetzt ist, fängt Wails externes Ziehen von Dateien ab, lässt internes HTML5-Ziehen jedoch wie gewohnt passieren.

Um in den Handlern Ihrer HTML-Ablagezone zwischen beiden zu unterscheiden, prüfen Sie, ob der Ziehvorgang Dateien enthält:

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

Dadurch reagieren Ihre HTML-Ablage-Handler nur auf interne Ziehvorgänge, etwa das Verschieben von Listenelementen, während Wails externe Dateiablagen separat über das `WindowFilesDropped`-Ereignis verarbeitet.

## Nächste Schritte

- [HTML-Drag-and-drop](/features/drag-and-drop/html/) – Elemente innerhalb Ihrer Anwendung ziehen
- [Fensteroptionen](/features/windows/options/) – Alle Optionen zur Fensterkonfiguration
