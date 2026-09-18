---
title: "API Dialog"
description: "Referensi lengkap untuk API dialog native"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## Ikhtisar

API Dialog menyediakan metode untuk menampilkan dialog file dan dialog pesan native. Akses dialog melalui pengelola `app.Dialog`.

**Jenis Dialog:**

- **Dialog File** - Dialog Buka dan Simpan
- **Dialog Pesan** - Dialog Info, Kesalahan, Peringatan, dan Pertanyaan

Semua dialog merupakan **dialog native sistem operasi** yang tampilan dan nuansanya sesuai dengan platform.

## Mengakses Dialog

Dialog diakses melalui pengelola `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## Dialog File

### OpenFile()

Membuat dialog untuk membuka file.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**Contoh:**

```go
dialog := app.Dialog.OpenFile()
```

### Metode OpenFileDialogStruct

#### SetTitle()

Menetapkan judul dialog.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**Contoh:**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

Menambahkan filter jenis file.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**Parameter:**

- `displayName` - Deskripsi filter yang ditampilkan kepada pengguna (misalnya, "Gambar", "Dokumen")
- `pattern` - Daftar ekstensi yang dipisahkan dengan titik koma (misalnya, "*.png;*.jpg")

**Contoh:**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

Menetapkan direktori awal.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**Contoh:**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

Mengaktifkan atau menonaktifkan pemilihan direktori.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**Contoh (pemilihan folder):**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

Mengaktifkan atau menonaktifkan pemilihan file.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

Mengaktifkan atau menonaktifkan pembuatan direktori baru.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

Menampilkan atau menyembunyikan file tersembunyi.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

Mengaitkan dialog ke jendela tertentu.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

Menampilkan dialog dan mengembalikan file yang dipilih.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Nilai yang dikembalikan:**

- `string` - Jalur file yang dipilih. Perlakukan string kosong sebagai "tidak ada pilihan" (bergantung pada sistem operasi, implementasi platform dapat mengembalikan string kosong atau kesalahan non-nil ketika operasi dibatalkan).
- `error` - Bernilai non-nil jika dialog itu sendiri gagal ditampilkan.

**Contoh:**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

Menampilkan dialog dan mengembalikan beberapa file yang dipilih.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**Nilai yang dikembalikan:**

- `[]string` - Larik jalur file yang dipilih
- `error` - Kesalahan jika dialog gagal

**Contoh:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

Membuat dialog untuk menyimpan file.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**Contoh:**

```go
dialog := app.Dialog.SaveFile()
```

### Metode SaveFileDialogStruct

#### SetTitle()

Menetapkan judul dialog.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

Menetapkan nama file default.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**Contoh:**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

Menambahkan filter jenis file.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**Contoh:**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

Menetapkan direktori awal.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

Mengaitkan dialog ke jendela tertentu.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

Menampilkan dialog dan mengembalikan jalur penyimpanan.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Contoh:**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### Pemilihan Folder

Tidak ada `SelectFolderDialog` terpisah. Gunakan `OpenFile()` dengan opsi direktori:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## Dialog Pesan

Semua dialog pesan mengembalikan `*MessageDialog` dan menggunakan metode yang sama.

### Info()

Membuat dialog informasi.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**Contoh:**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

Membuat dialog kesalahan.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**Contoh:**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

Membuat dialog peringatan.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**Contoh:**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

Membuat dialog pertanyaan dengan tombol khusus.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**Contoh:**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### Metode MessageDialog

#### SetTitle()

Mengatur judul dialog.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

Mengatur pesan dialog.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

Mengatur ikon khusus untuk dialog.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

Menambahkan tombol ke dialog dan mengembalikan tombol tersebut agar dapat dikonfigurasi.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**Mengembalikan:** `*Button` - Instans tombol untuk konfigurasi lebih lanjut

**Contoh:**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

Menentukan tombol default (diaktifkan dengan menekan Enter).

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**Contoh:**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

Menentukan tombol batal (diaktifkan dengan menekan Escape).

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**Contoh:**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

Mengaitkan dialog ke jendela tertentu.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

Menampilkan dialog. Callback tombol menangani respons pengguna.

```go
func (d *MessageDialog) Show()
```

**Catatan:** `Show()` tidak mengembalikan nilai. Gunakan callback tombol untuk menangani respons pengguna.

### Metode Tombol

#### OnClick()

Mengatur fungsi callback yang dijalankan saat tombol diklik.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

Menandai tombol ini sebagai tombol default.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

Menandai tombol ini sebagai tombol batal.

```go
func (b *Button) SetAsCancel() *Button
```

## Contoh Lengkap

### Contoh Pemilihan File

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Contoh Dialog Konfirmasi

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Dialog Simpan Perubahan

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Pemrosesan Banyak File

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### Penanganan Kesalahan dengan Dialog

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### Default Khusus Platform

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## Praktik Terbaik

### Yang Harus Dilakukan

- **Gunakan dialog native** - Tampilan dan nuansanya sesuai dengan platform
- **Berikan judul yang jelas** - Bantu pengguna memahami tujuannya
- **Atur filter yang sesuai** - Arahkan pengguna ke jenis file yang benar
- **Tangani pembatalan** - Periksa kesalahan (pengguna mungkin membatalkan)
- **Tampilkan konfirmasi untuk tindakan destruktif** - Gunakan dialog Question
- **Berikan umpan balik** - Gunakan dialog Info untuk pesan keberhasilan
- **Tetapkan default yang wajar** - Direktori default, nama file, dan sebagainya
- **Gunakan callback untuk tindakan tombol** - Tangani respons pengguna dengan benar

### Yang Tidak Boleh Dilakukan

- **Jangan abaikan kesalahan** - Pembatalan oleh pengguna mengembalikan kesalahan
- **Jangan gunakan label tombol yang ambigu** - Gunakan label yang spesifik: "Simpan"/"Batal"
- **Jangan terlalu sering menggunakan dialog** - Dialog mengganggu alur kerja
- **Jangan tampilkan kesalahan saat pengguna membatalkan** - Pembatalan adalah tindakan normal
- **Jangan lupakan filter file** - Bantu pengguna menemukan file yang tepat
- **Jangan menulis jalur secara hardcode** - Gunakan os.UserHomeDir() atau fungsi serupa

## Jenis Dialog menurut Platform

### macOS

- Dialog bergeser turun dari bilah judul
- Gaya "sheet" yang terpasang pada jendela induk
- Tampilan asli macOS

### Windows

- Dialog Windows standar
- Mengikuti pedoman desain Windows
- Tampilan Windows modern 10/11

### Linux

- Dialog GTK pada sistem berbasis GTK
- Dialog Qt pada sistem berbasis Qt
- Sesuai dengan lingkungan desktop

#### Perilaku Dialog Linux

Di Linux, build GTK4 default menggunakan **xdg-desktop-portal** untuk dialog file, yang menyediakan integrasi desktop asli, tetapi menyebabkan beberapa opsi tidak berpengaruh. Jalur GTK3 lama (`-tags gtk3`) tetap memberikan kontrol terprogram penuh atas opsi-opsi berikut:

| Opsi | GTK3 (`-tags gtk3`) | GTK4 (bawaan) | Catatan |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ Berfungsi | ❌ Tidak berpengaruh | Pengguna mengaturnya melalui tombol alih pada UI dialog (Ctrl+H atau menu) |
| `CanCreateDirectories()` | ✅ Berfungsi | ❌ Tidak berpengaruh | Selalu diaktifkan di portal |
| `ResolvesAliases()` | ✅ Berfungsi | ❌ Tidak berpengaruh | Portal menangani resolusi symlink |
| `SetButtonText()` | ✅ Berfungsi | ✅ Berfungsi | Teks khusus pada tombol konfirmasi berfungsi |

**Alasan adanya batasan ini:** Dialog berbasis portal GTK4 mendelegasikan kontrol UI kepada lingkungan desktop (GNOME, KDE, dan sebagainya). Ini memang dirancang demikian — portal memberikan UX yang konsisten di seluruh aplikasi dan menghormati preferensi pengguna.

@note{type="info"}
Build GTK4 default menggunakan dialog berbasis portal. Jika aplikasi Anda memerlukan kontrol terprogram penuh atas opsi dialog di atas, lakukan build dengan jalur `-tags gtk3` lama (didukung hingga v3.0.x; dihapus pada v3.1) — lihat [Pengemasan Linux - Dukungan GTK3 Lama](/guides/build/linux/#legacy-gtk3-support).

@end

## Pola Umum

### Pola "Simpan Sebagai"

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Pola "Buka yang Terbaru"

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
