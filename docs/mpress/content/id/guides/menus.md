---
title: "Menu"
description: "Panduan untuk membuat dan menyesuaikan menu di Wails v3"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 menyediakan sistem menu canggih yang memungkinkan Anda membuat menu aplikasi dan menu konteks. Panduan ini akan menjelaskan berbagai fitur dan kemampuan sistem menu tersebut.

## Membuat Menu

Untuk membuat menu baru, gunakan metode `New()` dari pengelola Menus:

```go
menu := app.Menu.New()
```

### Menambahkan Item Menu

Wails mendukung beberapa jenis item menu, masing-masing dengan kegunaan tertentu:

#### Item Menu Biasa

Item menu biasa merupakan komponen dasar menu. Item ini menampilkan teks dan dapat memicu tindakan saat diklik:

```go
menuItem := menu.Add("Click Me")
```

#### Kotak Centang

Item menu kotak centang menyediakan status yang dapat diaktifkan atau dinonaktifkan, sehingga berguna untuk mengaktifkan atau menonaktifkan fitur maupun pengaturan:

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### Grup Tombol Radio

Grup tombol radio memungkinkan pengguna memilih satu opsi dari sekumpulan pilihan yang saling eksklusif. Grup ini dibuat secara otomatis ketika item tombol radio ditempatkan bersebelahan:

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### Pemisah

Pemisah adalah garis horizontal yang membantu mengelompokkan item menu secara logis:

```go
menu.AddSeparator()
```

#### Submenu

Submenu adalah menu bertingkat yang muncul saat penunjuk diarahkan ke item menu atau item tersebut diklik. Submenu berguna untuk menata struktur menu yang kompleks:

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### Menggabungkan menu

Menu dapat ditambahkan ke menu lain dengan menempatkannya di akhir atau di awal.

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
Secara default, `prepend` dan `append` akan berbagi status dengan menu asli. Jika ingin membuat menu baru dengan statusnya sendiri, Anda dapat memanggil `.Clone()` pada menu tersebut.

Contoh: `menu.Append(secondaryMenu.Clone())`

@end

#### Mengosongkan menu

Dalam beberapa kasus, sebaiknya buat menu yang sepenuhnya baru jika Anda menangani jumlah item menu yang berubah-ubah.

Tindakan ini akan menghapus semua item dari menu yang ada dan memungkinkan Anda menambahkan item kembali.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
Mengosongkan menu hanya menghapus item menu pada tingkat teratas. Meskipun submenu tidak akan terlihat, submenu tersebut tetap menggunakan memori, jadi pastikan Anda mengelola menu dengan cermat.

@end

#### Memusnahkan menu

Jika ingin mengosongkan dan melepaskan menu, gunakan metode `Destroy()`:

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### Properti Item Menu

Item menu memiliki beberapa properti yang dapat dikonfigurasi:

| Properti | Metode | Deskripsi |
| --- | --- | --- |
| Label | `SetLabel(string)` | Mengatur teks yang ditampilkan |
| Diaktifkan | `SetEnabled(bool)` | Mengaktifkan atau menonaktifkan item |
| Dicentang | `SetChecked(bool)` | Mengatur status centang (untuk kotak centang atau item tombol radio) |
| Keterangan alat | `SetTooltip(string)` | Mengatur teks keterangan alat |
| Disembunyikan | `SetHidden(bool)` | Menampilkan atau menyembunyikan item |
| Akselerator | `SetAccelerator(string)` | Mengatur pintasan papan ketik |

### Status Item Menu

Item menu dapat memiliki berbagai status yang mengatur visibilitas dan interaktivitasnya:

#### Visibilitas

Item menu dapat ditampilkan atau disembunyikan secara dinamis menggunakan metode `SetHidden()`:

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

Item menu yang disembunyikan akan dihapus sepenuhnya dari menu hingga ditampilkan kembali. Hal ini berguna untuk item menu kontekstual yang sebaiknya hanya muncul pada status aplikasi tertentu.

#### Status Aktif

Item menu dapat diaktifkan atau dinonaktifkan menggunakan metode `SetEnabled()`:

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

Item menu yang dinonaktifkan tetap terlihat, tetapi tampak berwarna abu-abu dan tidak dapat diklik. Hal ini umumnya digunakan untuk menunjukkan bahwa suatu tindakan sedang tidak tersedia, misalnya:

- Menonaktifkan "Simpan" ketika tidak ada perubahan yang perlu disimpan
- Menonaktifkan "Salin" ketika tidak ada yang dipilih
- Menonaktifkan "Urungkan" ketika tidak ada tindakan yang dapat diurungkan

#### Pengelolaan Status Dinamis

Anda dapat menggabungkan status-status ini dengan penangan peristiwa untuk membuat menu dinamis:

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

### Penanganan Peristiwa

Item menu dapat merespons peristiwa klik menggunakan metode `OnClick`:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

Konteks menyediakan informasi tentang item menu yang diklik:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### Item Menu Berbasis Peran

Wails menyediakan sekumpulan peran menu bawaan yang secara otomatis membuat item menu dengan fungsionalitas standar. Berikut adalah peran menu yang didukung:

#### Struktur Menu Lengkap

Peran-peran ini membuat struktur menu lengkap dengan fungsionalitas umum:

| Peran | Deskripsi | Catatan Platform |
| --- | --- | --- |
| `AppMenu` | Menu aplikasi dengan Tentang, Layanan, Sembunyikan/Tampilkan, dan Keluar | Khusus macOS |
| `EditMenu` | Menu Edit standar dengan Urungkan, Ulangi, Potong, Salin, Tempel, dan lain-lain | Semua platform |
| `ViewMenu` | Menu Tampilan dengan kontrol Muat Ulang, Zoom, dan Layar Penuh | Semua platform |
| `WindowMenu` | Kontrol jendela (Minimalkan, Zoom, dan lain-lain) | Semua platform |
| `HelpMenu` | Menu Bantuan dengan tautan "Pelajari Lebih Lanjut" ke situs web Wails | Semua platform |

#### Item Menu Individual

Peran-peran ini dapat digunakan untuk menambahkan item menu individual:

| Peran | Deskripsi | Catatan Platform |
| --- | --- | --- |
| `About` | Tampilkan dialog Tentang aplikasi | Semua platform |
| `Hide` | Sembunyikan aplikasi | Khusus macOS |
| `HideOthers` | Sembunyikan aplikasi lain | Khusus macOS |
| `UnHide` | Tampilkan aplikasi yang disembunyikan | Khusus macOS |
| `CloseWindow` | Tutup jendela saat ini | Semua platform |
| `Minimise` | Minimalkan jendela | Semua platform |
| `Zoom` | Zoom jendela | Khusus macOS |
| `Front` | Bawa jendela ke depan | Khusus macOS |
| `Quit` | Keluar dari aplikasi | Semua platform |
| `Undo` | Urungkan tindakan terakhir | Semua platform |
| `Redo` | Ulangi tindakan terakhir | Semua platform |
| `Cut` | Potong pilihan | Semua platform |
| `Copy` | Salin pilihan | Semua platform |
| `Paste` | Tempel dari papan klip | Semua platform |
| `PasteAndMatchStyle` | Tempel dan sesuaikan gaya | Khusus macOS |
| `SelectAll` | Pilih semua | Semua platform |
| `Delete` | Hapus pilihan | Semua platform |
| `Reload` | Muat ulang halaman saat ini | Semua platform |
| `ForceReload` | Paksa muat ulang halaman saat ini | Semua platform |
| `ToggleFullscreen` | Alihkan mode layar penuh | Semua platform |
| `ResetZoom` | Atur ulang tingkat zoom | Semua platform |
| `ZoomIn` | Perbesar tampilan | Semua platform |
| `ZoomOut` | Perkecil tampilan | Semua platform |

Berikut adalah contoh yang menunjukkan cara menggunakan menu lengkap dan peran individual:

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

## Menu Aplikasi

Menu aplikasi adalah menu yang muncul di bagian atas jendela aplikasi Anda (Windows/Linux) atau di bagian atas layar (macOS).

### Perilaku Menu Aplikasi

Saat Anda menetapkan menu aplikasi menggunakan `app.Menu.Set()`, menu tersebut menjadi menu utama di macOS. Di Windows/Linux, menu ditetapkan untuk setiap jendela.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

Berikut adalah contoh lengkap yang menunjukkan berbagai perilaku menu tersebut:

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

## Menu Konteks

Menu konteks adalah menu pop-up yang muncul saat Anda mengeklik kanan elemen dalam aplikasi. Menu ini memberikan akses cepat ke tindakan yang relevan untuk elemen yang diklik.

### Menu Konteks Bawaan

Menu konteks bawaan adalah menu konteks bawaan webview yang menyediakan operasi tingkat sistem seperti:

- Salin, Potong, dan Tempel untuk memanipulasi teks
- Kontrol pemilihan teks
- Opsi pemeriksaan ejaan

#### Mengontrol Menu Konteks Bawaan

Anda dapat mengontrol kapan menu konteks bawaan muncul menggunakan properti CSS `--default-contextmenu`:

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
Fitur ini hanya akan berfungsi sesuai harapan setelah [runtime frontend siap](/reference/frontend-runtime/).

@end

#### Perilaku Menu Konteks Bertingkat

Saat menggunakan properti `--default-contextmenu` pada elemen bertingkat, aturan berikut berlaku:

1. Elemen anak mewarisi pengaturan menu konteks induknya kecuali jika ditimpa secara eksplisit
2. Pengaturan yang paling spesifik (paling dekat) akan diprioritaskan
3. Nilai `auto` dapat digunakan untuk mengatur ulang ke perilaku bawaan

Contoh perilaku menu konteks bertingkat:

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

### Menu Konteks Kustom

Menu konteks kustom memungkinkan Anda menyediakan tindakan khusus aplikasi yang relevan dengan elemen yang diklik. Menu ini sangat berguna untuk:

- Operasi file dalam pengelola dokumen
- Alat manipulasi gambar
- Tindakan kustom dalam kisi data
- Operasi khusus komponen

#### Membuat Menu Konteks Kustom

Saat membuat menu konteks kustom, Anda memberikan pengenal unik (nama) yang menghubungkan menu tersebut dengan elemen HTML:

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

Parameter nama ("imageMenu" dalam contoh ini) berfungsi sebagai pengenal unik yang akan digunakan untuk:

1. Menghubungkan elemen HTML ke menu konteks khusus ini
2. Menentukan menu yang harus ditampilkan saat pengguna mengeklik kanan
3. Memungkinkan pembaruan dan pembersihan menu

#### Data Konteks

Saat menangani peristiwa menu konteks, Anda dapat mengakses item menu yang diklik beserta data konteks yang terkait dengannya:

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

Data konteks diteruskan dari properti `--custom-contextmenu-data` milik elemen HTML dan tersedia di pengendali klik melalui `ctx.ContextMenuData()`. Hal ini sangat berguna saat:

- Menggunakan daftar atau grid yang setiap itemnya memerlukan identifikasi unik
- Menangani operasi pada komponen atau elemen tertentu
- Meneruskan status atau metadata dari frontend ke backend

#### Pengelolaan Menu Konteks

Setelah mengubah menu konteks, panggil metode `Update()` untuk menerapkan perubahan tersebut:

```go
contextMenu.Update()
```

Jika menu konteks tidak lagi diperlukan, Anda dapat memusnahkannya:

```go
contextMenu.Destroy()
```

@note{type="danger" title="Peringatan"}
Setelah memanggil `Destroy()`, penggunaan kembali referensi menu konteks akan menyebabkan panic.

@end

### Contoh Dunia Nyata: Galeri Gambar

Berikut adalah contoh lengkap penerapan menu konteks kustom untuk galeri gambar:

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

Dalam contoh ini:

1. Menu konteks dibuat dengan pengenal "imageMenu"
2. Setiap kontainer gambar dihubungkan ke menu menggunakan `--custom-contextmenu: imageMenu`
3. Setiap kontainer menyediakan ID gambarnya sebagai data konteks menggunakan `--custom-contextmenu-data`
4. Backend menerima ID gambar di pengendali klik dan dapat menjalankan operasi tertentu
5. Menu yang sama digunakan kembali untuk semua gambar, tetapi data konteks menunjukkan gambar yang akan dioperasikan

Pola ini sangat efektif untuk:

- Grid data yang baris-barisnya memerlukan operasi tertentu
- Pengelola file yang file-filenya memerlukan tindakan sesuai konteks
- Alat desain yang elemen-elemennya memerlukan operasi berbeda
- Komponen apa pun yang menerapkan operasi yang sama pada beberapa instans
