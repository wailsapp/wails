---
title: "Referensi Menu"
description: "Referensi lengkap untuk jenis, properti, dan metode item menu"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Referensi Menu

Referensi lengkap untuk jenis dan properti item menu, serta perilaku dinamis. Buat menu profesional dan responsif dengan kotak centang, grup tombol radio, pemisah, dan pembaruan dinamis.

## Jenis Item Menu

### Item Menu Biasa

Jenis yang paling umum—menampilkan teks dan memicu tindakan:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Gunakan untuk:** Perintah, tindakan, membuka jendela

### Kotak Centang

Item menu yang dapat diaktifkan atau dinonaktifkan dengan status dicentang/tidak dicentang:

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Gunakan untuk:** Pengaturan Boolean, pengalih fitur, opsi tampilan

**Penting:** Status centang berubah secara otomatis saat diklik.

### Grup Tombol Radio

Opsi yang saling eksklusif—hanya satu yang dapat dipilih:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Gunakan untuk:** Pilihan yang saling eksklusif (ukuran, tema, mode)

**Cara kerja pengelompokan:**

- Item radio yang bersebelahan secara otomatis membentuk grup
- Memilih satu item akan membatalkan pilihan item lain dalam grup
- Pisahkan grup dengan pemisah atau item biasa

**Contoh dengan beberapa grup:**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### Submenu

Struktur menu bertingkat untuk pengorganisasian:

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Gunakan untuk:** Mengelompokkan item terkait, mengurangi kepadatan

**Batas tingkat:** Sebagian besar platform mendukung 2-3 tingkat. Hindari tingkat yang lebih dalam.

### Pemisah

Pembatas visual di antara item menu:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Gunakan untuk:** Mengelompokkan item terkait secara visual

**Praktik terbaik:** Jangan mengawali atau mengakhiri menu dengan pemisah.

## Properti Item Menu

### Label

Teks yang ditampilkan untuk item menu:

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**Label dinamis:**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Status Aktif

Atur apakah pengguna dapat berinteraksi dengan item menu:

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Perilaku Menu di Windows"}
Di Windows, menu harus dibuat ulang ketika statusnya berubah. **Selalu panggil `menu.Update()` setelah mengaktifkan/menonaktifkan item menu**, terutama jika item tersebut dibuat dalam keadaan nonaktif.

**Alasannya:** Menu Windows dibangun ulang dari awal ketika diperbarui. Jika Anda tidak memanggil `Update()`, pengendali klik tidak akan dipicu dengan benar.

@end

**Contoh: Pengaktifan/penonaktifan dinamis**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**Pola umum: Aktifkan berdasarkan kondisi**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### Status Centang

Untuk item kotak centang dan tombol radio, atur atau periksa status centangnya:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Perubahan otomatis:** Status kotak centang berubah secara otomatis saat diklik. Anda tidak perlu memanggil `SetChecked()` di pengendali klik.

**Kontrol manual:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### Akselerator (Pintasan Papan Ketik)

Tambahkan pintasan papan ketik ke item menu:

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**Format akselerator:**

- `CmdOrCtrl` - Cmd di macOS, Ctrl di Windows/Linux
- `Shift`, `Alt`, `Option` - Tombol pengubah
- `A-Z`, `0-9` - Tombol huruf/angka
- `F1-F12` - Tombol fungsi
- `Enter`, `Space`, `Backspace`, dll. - Tombol khusus

**Contoh:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**Akselerator khusus platform:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### Tooltip

Tambahkan teks yang muncul saat penunjuk diarahkan ke item menu (dukungan berbeda-beda menurut platform):

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Dukungan platform:**

- **Windows:** ✅ Didukung
- **macOS:** ❌ Tidak didukung (tooltip bukan fitur standar untuk menu)
- **Linux:** ⚠️ Berbeda-beda menurut lingkungan desktop

### Status Tersembunyi

Sembunyikan item menu tanpa menghapusnya:

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**Gunakan untuk:** Opsi debug, flag fitur, fitur bersyarat

## Penanganan Peristiwa

### Pengendali OnClick

Jalankan kode saat item menu diklik:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Konteks menyediakan:**

- `ctx.ClickedMenuItem()` - Item menu yang diklik
- Konteks jendela (jika berasal dari menu jendela)
- Konteks aplikasi

**Contoh: Mengakses item menu dalam handler**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### Beberapa Handler

Anda dapat menetapkan beberapa handler (handler terakhir yang berlaku):

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

**Praktik terbaik:** Tetapkan handler satu kali dan gunakan logika kondisional di dalamnya jika diperlukan.

## Menu Dinamis

### Memperbarui Item Menu

**Aturan utama:** Selalu panggil `menu.Update()` setelah mengubah status menu.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**Mengapa hal ini penting:**

- **Windows:** Menu dibuat ulang saat diperbarui
- **macOS/Linux:** Tidak terlalu penting, tetapi tetap disarankan
- **Handler klik:** Tidak akan dipicu dengan benar tanpa Update()

### Membuat Ulang Menu

Untuk perubahan besar, buat ulang seluruh menu:

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**Kapan harus membuat ulang:**

- Daftar file terbaru berubah
- Menu plugin berubah
- Transisi status utama

**Kapan harus memperbarui:**

- Mengaktifkan/menonaktifkan item
- Mengubah label
- Mengalihkan status kotak centang

### Menu Peka Konteks

Sesuaikan menu berdasarkan status aplikasi:

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## Perbedaan Antarplatform

### Lokasi Bilah Menu

| Platform | Lokasi | Catatan |
| --- | --- | --- |
| **macOS** | Bagian atas layar | Bilah menu global |
| **Windows** | Bagian atas jendela | Menu per jendela |
| **Linux** | Bagian atas jendela | Per jendela (biasanya) |

### Menu Standar

**macOS:**

- Memiliki menu "Aplikasi" (dengan nama aplikasi)
- "Preferensi" berada di menu Aplikasi
- "Keluar" berada di menu Aplikasi

**Windows/Linux:**

- Tidak ada menu Aplikasi
- "Preferensi" berada di menu Edit atau Alat
- "Keluar" berada di menu File

**Contoh: Struktur yang sesuai dengan platform**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### Konvensi Akselerator

**macOS:**

- `Cmd+` untuk sebagian besar pintasan
- `Cmd+,` untuk Preferensi
- `Cmd+Q` untuk Keluar

**Windows:**

- `Ctrl+` untuk sebagian besar pintasan
- `Ctrl+P` atau `Ctrl+,` untuk Preferensi
- `Alt+F4` untuk Keluar (atau `Ctrl+Q`)

**Linux:**

- Umumnya mengikuti konvensi Windows
- Lingkungan desktop dapat menggantinya

## Praktik Terbaik

### ✅ Lakukan

- **Panggil menu.Update()** setelah mengubah status menu (terutama pada Windows)
- **Gunakan grup tombol radio** untuk opsi yang saling eksklusif
- **Gunakan kotak centang** untuk fitur yang dapat diaktifkan atau dinonaktifkan
- **Tambahkan pintasan keyboard** ke tindakan yang umum digunakan
- **Kelompokkan item terkait** dengan pemisah
- **Uji di semua platform** — perilakunya berbeda-beda

### ❌ Jangan lakukan

- **Jangan lupa memanggil menu.Update()** — handler klik tidak akan berfungsi dengan benar
- **Jangan buat submenu terlalu dalam** — maksimum 2-3 tingkat
- **Jangan awali atau akhiri dengan pemisah** — terlihat tidak profesional
- **Jangan gunakan tooltip di macOS** — tidak didukung
- **Jangan tulis pintasan khusus platform secara hardcode** — gunakan `CmdOrCtrl`

## Pemecahan Masalah

### Item Menu Tidak Merespons

**Gejala:** Handler klik tidak dijalankan

**Penyebab:** Lupa memanggil `menu.Update()` setelah mengaktifkan item

**Solusi:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Item Menu Berwarna Abu-abu

**Gejala:** Item menu tidak dapat diklik

**Penyebab:** Item dinonaktifkan

**Solusi:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Pintasan Keyboard Tidak Berfungsi

**Gejala:** Pintasan keyboard tidak memicu item menu

**Penyebab:**

1. Format pintasan keyboard salah
2. Konflik dengan pintasan sistem
3. Jendela tidak memiliki fokus

**Solusi:**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## Langkah Berikutnya

- [Menu Aplikasi](/features/menus/application/) — buat bilah menu aplikasi
- [Menu Konteks](/features/menus/context/) — menu konteks klik kanan
- [Menu Baki Sistem](/features/menus/systray/) — menu baki sistem/bilah menu
- [Pola Menu](/guides/menus/) — pola menu umum dan praktik terbaik

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh menu](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
