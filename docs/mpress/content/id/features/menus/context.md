---
title: "Menu Konteks"
description: "Buat menu konteks klik kanan untuk aplikasi Anda"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## Masalah

Pengguna mengharapkan menu klik kanan dengan tindakan yang sesuai dengan konteks. Elemen yang berbeda memerlukan menu yang berbeda:

- **Teks**: Potong, Salin, Tempel
- **Gambar**: Simpan, Salin, Buka
- **Elemen khusus**: Tindakan khusus aplikasi

Membuat menu konteks secara manual berarti harus menangani peristiwa mouse, pemosisian, dan perbedaan antarplatform.

## Solusi Wails

Wails menyediakan **menu konteks deklaratif** menggunakan properti CSS. Kaitkan menu dengan elemen HTML, teruskan data, dan tangani klik—semuanya dengan perilaku native platform.

![Menu konteks khusus Wails yang ditampilkan di atas webview pada macOS](/assets/screenshots/context-menu-macos.png)

Menu tersebut bersifat native bagi platform, sedangkan elemen yang membukanya tetap menjadi bagian dari webview Anda. Tangkapan layar macOS ini menggunakan pendaftaran menu konteks khusus dari contoh di bawah.

## Mulai Cepat

**Kode Go:**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML:**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**Selesai!** Mengklik kanan area teks akan menampilkan menu khusus Anda.

## Membuat Menu Konteks

### Menu Konteks Dasar

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

**ID menu:** Harus unik. Digunakan untuk mengaitkan menu dengan elemen HTML.

### Dengan Submenu

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

### Dengan Kotak Centang dan Grup Tombol Radio

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

Untuk **semua jenis item menu**, lihat [Referensi Menu](/features/menus/reference/).

## Mengaitkan dengan Elemen HTML

Gunakan properti khusus CSS untuk memasang menu konteks:

### Pengaitan Dasar

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**Properti CSS:** `--custom-contextmenu: <menu-id>`

### Dengan Data Konteks

Teruskan data dari HTML ke Go:

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Handler Go:**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**Properti CSS:**

- `--custom-contextmenu: <menu-id>` - Menu yang akan ditampilkan
- `--custom-contextmenu-data: <data>` - Data yang akan diteruskan ke handler

### Data Dinamis

Buat data secara dinamis dalam JavaScript:

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

### Beberapa Elemen, Menu yang Sama

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

**Satu menu, data berbeda untuk setiap elemen.**

## Data Konteks

### Mengakses Data Konteks

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**Jenis data:** Selalu `string`. Uraikan sesuai kebutuhan.

### Meneruskan Data Kompleks

Gunakan JSON untuk data kompleks:

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Handler Go:**

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

@note{type="caution" title="Keamanan"}
**Selalu validasi data konteks** dari frontend. Pengguna dapat memanipulasi properti CSS, jadi perlakukan data sebagai input yang tidak tepercaya.

@end

### Contoh Validasi

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

## Menu Konteks Bawaan

WebView menyediakan menu konteks bawaan untuk operasi standar (salin, tempel, inspeksi). Kendalikan dengan `--default-contextmenu`:

### Sembunyikan Menu Bawaan

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**Kasus penggunaan:** Elemen UI khusus yang tidak cocok menggunakan menu bawaan.

### Tampilkan Menu Bawaan

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**Kasus penggunaan:** Area teks, bidang input, dan konten yang dapat diedit.

### Mode Otomatis (Cerdas)

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**Perilaku bawaan.** Menampilkan menu bawaan ketika:

- Teks dipilih
- Berada dalam bidang input teks
- Berada dalam konten yang dapat diedit (`contenteditable`)

Dalam kondisi lain, menu bawaan disembunyikan.

### Menggabungkan Menu Khusus dan Bawaan

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**Perilaku:**

1. Menu khusus ditampilkan terlebih dahulu
2. Jika menu khusus kosong atau tidak ditemukan, menu bawaan ditampilkan
3. Keduanya dapat digunakan secara bersamaan (bergantung pada platform)

## Menu Konteks Dinamis

Perbarui menu berdasarkan status aplikasi:

### Mengaktifkan/Menonaktifkan Item

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

@note{type="caution" title="Selalu Panggil Update()"}
Setelah mengubah status menu, **panggil `contextMenu.Update()`**. Hal ini sangat penting pada Windows.

Lihat [Referensi Menu](/features/menus/reference/#enabled-state) untuk detail selengkapnya.

@end

### Mengubah Label

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

### Membuat Ulang Menu

Untuk perubahan besar, buat ulang seluruh menu:

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

## Perilaku Platform

Menu konteks bersifat **native pada setiap platform**:

@tabs{sync-key="platform"}
[macOS]
**Menu konteks native macOS:**

- Animasi dan transisi sistem
- Klik kanan = Control+Klik (otomatis)
- Menyesuaikan dengan tampilan sistem (terang/gelap)
- Operasi teks standar dalam menu default
- Pengguliran native untuk menu panjang

**Konvensi macOS:**

- Gunakan kapitalisasi kalimat untuk item menu
- Gunakan elipsis (...) untuk item yang membuka dialog
- Pintasan umum: ⌘C (Salin), ⌘V (Tempel)

[Windows]
**Menu konteks native Windows:**

- Gaya native Windows
- Mengikuti tema Windows
- Operasi standar Windows dalam menu default
- Dukungan input sentuh dan pena

**Konvensi Windows:**

- Gunakan kapitalisasi judul untuk item menu
- Gunakan elipsis (...) untuk item yang membuka dialog
- Pintasan umum: Ctrl+C (Salin), Ctrl+V (Tempel)

[Linux]
**Integrasi lingkungan desktop:**

- Menyesuaikan dengan tema desktop (GTK, Qt, dan sebagainya)
- Perilaku klik kanan mengikuti pengaturan sistem
- Isi menu default berbeda-beda menurut lingkungan
- Penempatan mengikuti konvensi lingkungan desktop

**Pertimbangan untuk Linux:**

- Uji pada lingkungan desktop target
- GTK dan Qt memiliki perilaku yang berbeda
- Beberapa lingkungan desktop menyesuaikan menu konteks

@end

## Contoh Lengkap

**Kode Go:**

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

**HTML:**

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

## Praktik Terbaik

### ✅ Lakukan

- **Jaga agar menu tetap terfokus** - Hanya tampilkan tindakan yang relevan bagi elemen tersebut
- **Validasi data konteks** - Perlakukan sebagai input yang tidak tepercaya
- **Gunakan label yang jelas** - Gunakan "Hapus File", bukan "Hapus"
- **Panggil menu.Update()** - Setelah mengubah status menu
- **Uji pada semua platform** - Perilakunya berbeda-beda
- **Sediakan pintasan papan ketik** - Untuk tindakan umum
- **Kelompokkan item terkait** - Gunakan pemisah

### ❌ Jangan Lakukan

- **Jangan memercayai data konteks** - Selalu lakukan validasi
- **Jangan membuat menu terlalu panjang** - Maksimum 7-10 item
- **Jangan lupa memanggil menu.Update()** - Menu tidak akan berfungsi dengan benar
- **Jangan membuat submenu terlalu bertingkat** - Maksimum 2 tingkat
- **Jangan gunakan jargon** - Pastikan label mudah dipahami pengguna
- **Jangan memblokir handler** - Pastikan handler tetap cepat

## Pemecahan Masalah

### Menu Konteks Tidak Muncul

**Kemungkinan penyebab:**

1. ID menu tidak cocok
2. Kesalahan pengetikan pada properti CSS
3. Runtime belum diinisialisasi

**Solusi:**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### Data Konteks Tidak Diterima

**Kemungkinan penyebab:**

1. Properti CSS belum ditetapkan
2. Data berisi karakter khusus

**Solusi:**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

Atau gunakan JavaScript:

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### Item Menu Tidak Merespons

**Penyebab:** Lupa memanggil `menu.Update()` setelah mengaktifkan

**Solusi:**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## Langkah Selanjutnya

@cards{cols="2"}
📖 Referensi Menu
Referensi lengkap untuk jenis dan properti item menu.

[Pelajari Selengkapnya →](/features/menus/reference/)

---
☰ Menu Aplikasi
Buat bilah menu aplikasi.

[Pelajari Selengkapnya →](/features/menus/application/)

---
★ Menu Baki Sistem
Tambahkan integrasi dengan baki sistem/bilah menu.

[Pelajari Selengkapnya →](/features/menus/systray/)

---
📖 Pola Menu
Pola menu umum dan praktik terbaik.

[Pelajari Selengkapnya →](/guides/menus/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh menu konteks](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus).
