---
title: "Penjatuhan File"
description: "Terima file yang diseret dari sistem operasi ke aplikasi Anda"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails memungkinkan pengguna menyeret file dari sistem operasi (pengelola file, desktop) ke aplikasi Anda. Berbeda dengan seret dan jatuhkan HTML5 yang hanya berfungsi di dalam browser, fitur ini memberi Anda akses ke jalur file yang sebenarnya pada disk.

![Contoh seret dan jatuhkan Wails dengan zona penjatuhan file eksternal di macOS](/assets/screenshots/file-drop-macos.png)

Zona penjatuhan eksternal merupakan bagian dari webview, sedangkan Wails menyediakan peristiwa penjatuhan file native dari sistem operasi.

## Aktifkan Penjatuhan File

Penjatuhan file dinonaktifkan secara default. Untuk mengaktifkannya, atur `EnableFileDrop: true` dalam opsi jendela Anda:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

Saat `EnableFileDrop` bernilai `false` (nilai default), file yang diseret dari sistem operasi akan diblokir—file tersebut tidak akan dibuka di webview atau memicu peristiwa apa pun. Hal ini mencegah navigasi yang tidak disengaja ketika pengguna menyeret file di atas aplikasi Anda.

## Tentukan Zona Penjatuhan

Zona penjatuhan memberi tahu Wails elemen mana yang harus menerima file. File yang dijatuhkan di luar zona penjatuhan akan diabaikan.

Tambahkan atribut `data-file-drop-target` ke elemen apa pun:

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

Anda dapat memiliki beberapa zona penjatuhan. `id` dan kelas CSS milik elemen diteruskan ke kode Go Anda, sehingga Anda dapat menangani penjatuhan secara berbeda sesuai lokasi file dijatuhkan.

## Atur Gaya Saat Seretan Melintas

Saat file diseret di atas zona penjatuhan, Wails menambahkan kelas `file-drop-target-active`. Dengan demikian, Anda dapat memberikan umpan balik visual agar pengguna mengetahui tempat mereka dapat menjatuhkan file:

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

Kelas tersebut dihapus secara otomatis saat file meninggalkan zona atau dijatuhkan.

## Deteksi File yang Dijatuhkan

Saat file dijatuhkan pada zona penjatuhan yang valid, Wails memicu peristiwa `WindowFilesDropped`. Konteks peristiwa berisi jalur lengkap sistem file untuk semua file yang dijatuhkan:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

Jalur tersebut bersifat absolut, seperti `/home/user/documents/report.pdf` atau `C:\Users\Name\Documents\report.pdf`.

## Dapatkan Informasi Target Penjatuhan

Jika Anda memiliki beberapa zona penjatuhan, Anda dapat mengetahui zona yang menerima file menggunakan `DropTargetDetails()`:

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

Dengan demikian, Anda dapat mengarahkan file ke penangan yang berbeda:

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## Contoh Lengkap

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

## Penjatuhan di Seluruh Jendela

Jika Anda ingin file dapat dijatuhkan di mana saja dalam aplikasi, tambahkan atribut tersebut ke elemen body:

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

Anda dapat menggunakan lapisan CSS untuk menunjukkan bahwa seluruh jendela merupakan target penjatuhan:

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

## Menggabungkan dengan Seret dan Jatuhkan HTML

Anda dapat menggunakan penjatuhan file eksternal dan seret dan jatuhkan HTML internal dalam aplikasi yang sama. Saat `EnableFileDrop` bernilai `true`, Wails mencegat seretan file eksternal, tetapi tetap membiarkan seretan HTML5 internal berjalan seperti biasa.

Untuk membedakan keduanya dalam penangan zona penjatuhan HTML, periksa apakah seretan berisi file:

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

Hal ini memastikan penangan penjatuhan HTML Anda hanya merespons seretan internal (seperti memindahkan item daftar), sementara Wails menangani penjatuhan file eksternal secara terpisah melalui peristiwa `WindowFilesDropped`.

## Langkah Berikutnya

- [Seret dan Jatuhkan HTML](/features/drag-and-drop/html/)—Seret elemen di dalam aplikasi Anda
- [Opsi Jendela](/features/windows/options/)—Semua opsi konfigurasi jendela
