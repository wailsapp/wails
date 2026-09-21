---
title: "Ringkasan dialog"
description: "Tampilkan dialog sistem native di aplikasi Anda"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Dialog native

Wails menyediakan **dialog sistem native** yang berfungsi di semua platform: dialog pesan (informasi, peringatan, kesalahan, pertanyaan), dialog file (buka, simpan, folder), dan jendela dialog khusus dengan tampilan serta perilaku native platform.

![Dialog pertanyaan Wails di macOS dengan tombol Batal dan Buang](/assets/screenshots/dialog-question-macos.png)

API yang sama ditampilkan sesuai konvensi setiap platform yang didukung. Contoh macOS ini menunjukkan dialog pertanyaan yang ditautkan ke jendela Wails, termasuk tombol default dan tombol batal.

## Mulai Cepat

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

**Selesai!** Dialog native dengan kode minimal.

## Mengakses Dialog

Dialog diakses melalui pengelola `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Jenis Dialog

### Dialog informasi

Tampilkan pesan sederhana:

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**Kasus penggunaan:**

- Pesan keberhasilan
- Pemberitahuan informatif
- Konfirmasi penyelesaian

### Dialog peringatan

Tampilkan peringatan:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Kasus penggunaan:**

- Peringatan nonkritis
- Pemberitahuan deprekasi
- Pesan kehati-hatian

### Dialog kesalahan

Tampilkan kesalahan:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**Kasus penggunaan:**

- Pesan kesalahan
- Notifikasi kegagalan
- Penanganan pengecualian

### Dialog pertanyaan

Ajukan pertanyaan kepada pengguna dan tangani respons melalui callback tombol:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**Kasus penggunaan:**

- Mengonfirmasi tindakan
- Pertanyaan Ya/Tidak
- Pilihan ganda

## Dialog file

### Dialog Buka File

Pilih file yang akan dibuka:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**Pilihan jamak:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### Dialog Simpan File

Pilih lokasi penyimpanan:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### Dialog Pilih Folder

Pilih direktori menggunakan dialog buka file dengan pemilihan direktori diaktifkan:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## Opsi Dialog

### Judul dan Pesan

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Tombol

**Tombol default untuk dialog sederhana:**

Dialog informasi, peringatan, dan kesalahan menampilkan tombol "OK" secara default:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Tombol khusus untuk dialog pertanyaan:**

Gunakan `AddButton()` untuk menambahkan tombol. Metode ini mengembalikan `*Button` yang dapat Anda konfigurasi dengan callback:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Tombol default dan Batal:**

Gunakan `SetDefaultButton()` untuk menentukan tombol yang disorot dan dipicu oleh Enter. Gunakan `SetCancelButton()` untuk menentukan tombol yang dipicu oleh Escape.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### Penautan ke Jendela

Tautkan dialog ke jendela tertentu:

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Perilaku:**

- Dialog muncul di tengah jendela induk
- Jendela induk dinonaktifkan selama dialog ditampilkan
- Dialog bergerak bersama jendela induk (macOS)

## Perilaku Platform

@tabs{sync-key="platform"}
[macOS]
**Dialog macOS:**

- Tampilan NSAlert native
- Mengikuti tema sistem (terang/gelap)
- Mendukung navigasi papan ketik
- Pintasan standar (⌘. untuk Batal)
- Fitur aksesibilitas bawaan
- Bergaya lembar saat ditautkan ke jendela

**Contoh:**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Dialog Windows:**

- Tampilan TaskDialog native
- Mengikuti tema sistem
- Mendukung navigasi papan ketik
- Pintasan standar (Esc untuk Batal)
- Fitur aksesibilitas bawaan
- Modal terhadap jendela induk

**Contoh:**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Dialog Linux:**

- Tampilan dialog GTK
- Mengikuti tema desktop
- Mendukung navigasi papan ketik
- Integrasi dengan lingkungan desktop
- Bervariasi menurut lingkungan desktop (GNOME, KDE, dan sebagainya)

**Contoh:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Pola Umum

### Konfirmasi Sebelum Tindakan Destruktif

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Penanganan Kesalahan dengan Dialog

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### Pemilihan File dengan Validasi

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### Alur Dialog Bertahap

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan dialog native** — UX lebih baik daripada dialog khusus
- **Berikan pesan yang jelas** — Sampaikan secara spesifik
- **Tetapkan judul yang sesuai** — Konteks itu penting
- **Gunakan tombol default secara bijak** — Jadikan opsi aman sebagai default
- **Tangani pembatalan** — Pengguna mungkin membatalkan
- **Validasi file yang dipilih** — Periksa jenis file

### ❌ Jangan Lakukan

- **Jangan gunakan dialog secara berlebihan** — Mengganggu alur kerja
- **Jangan gunakan dialog untuk pesan yang sering muncul** — Gunakan notifikasi
- **Jangan lupa menangani kesalahan** — Pengguna mungkin membatalkan
- **Jangan lakukan pemblokiran jika tidak diperlukan** — Pertimbangkan alternatif
- **Jangan gunakan pesan generik** — Sampaikan secara spesifik
- **Jangan abaikan perbedaan platform** — Uji di semua platform

## Langkah Berikutnya

@cards{cols="2"}
ℹ Dialog pesan
Dialog informasi, peringatan, dan kesalahan.

[Pelajari Lebih Lanjut →](/features/dialogs/message/)

---
📖 Dialog file
Membuka, menyimpan, dan memilih folder.

[Pelajari Lebih Lanjut →](/features/dialogs/file/)

---
◆ Dialog khusus
Buat jendela dialog khusus.

[Pelajari Lebih Lanjut →](/features/dialogs/custom/)

---
▣ Jendela
Pelajari pengelolaan jendela.

[Pelajari Lebih Lanjut →](/features/windows/basics/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh dialog](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
