---
title: "Dialog file"
description: "Dialog untuk membuka, menyimpan, dan memilih folder"
slug: "features/dialogs/file"
sourcePath: "features/dialogs/file.md"
---

## Dialog file

Wails menyediakan **dialog file native** dengan tampilan yang sesuai untuk setiap platform guna membuka file, menyimpan file, dan memilih folder. API sederhana dengan pemfilteran jenis file, dukungan pemilihan beberapa file, dan lokasi default.

![Pemilih file native macOS yang dibuka oleh aplikasi Wails](/assets/screenshots/file-dialog-macos.png)

API Wails mendelegasikan tugas kepada pemilih milik sistem operasi sehingga perilaku navigasi, pemfilteran, dan pemilihan yang familier tetap dipertahankan.

## Membuat Dialog File

Dialog file diakses melalui pengelola `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Dialog Buka File

Pilih file yang akan dibuka:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

openFile(path)
```

**Kasus penggunaan:**

- Membuka dokumen
- Mengimpor file
- Memuat gambar
- Memilih file konfigurasi

### Pemilihan Satu File

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open Document").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    // User cancelled or error occurred
    return
}

// Use selected file
data, _ := os.ReadFile(path)
```

### Pemilihan Beberapa File

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
    PromptForMultipleSelection()

if err != nil {
    return
}

// Process all selected files
for _, path := range paths {
    processFile(path)
}
```

### Dengan Direktori Default

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open File").
    SetDirectory("/Users/me/Documents").
    PromptForSingleSelection()
```

## Dialog Simpan File

Pilih lokasi penyimpanan:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

saveFile(path, data)
```

**Kasus penggunaan:**

- Menyimpan dokumen
- Mengekspor data
- Membuat file baru
- Menyimpan sebagai...

### Dengan Nama File Default

```go
path, err := app.Dialog.SaveFile().
    SetFilename("export.csv").
    AddFilter("CSV Files", "*.csv").
    PromptForSingleSelection()
```

### Dengan Direktori Default

```go
path, err := app.Dialog.SaveFile().
    SetDirectory("/Users/me/Documents").
    SetFilename("untitled.txt").
    PromptForSingleSelection()
```

### Konfirmasi Penimpaan

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

// Check if file exists
if _, err := os.Stat(path); err == nil {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Overwrite").
        SetMessage("File already exists. Overwrite?")

    overwriteBtn := dialog.AddButton("Overwrite")
    overwriteBtn.OnClick(func() {
        saveFile(path, data)
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
    return
}

saveFile(path, data)
```

## Dialog Pilih Folder

Pilih direktori menggunakan dialog buka file dengan pemilihan direktori yang diaktifkan:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

exportToFolder(path)
```

**Kasus penggunaan:**

- Memilih direktori keluaran
- Memilih ruang kerja
- Memilih lokasi cadangan
- Memilih direktori instalasi

### Dengan Direktori Default

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    SetDirectory("/Users/me/Documents").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Filter File

Gunakan metode `AddFilter()` untuk menambahkan filter jenis file ke dialog. Setiap pemanggilan menambahkan opsi filter baru.

### Filter Dasar

```go
path, _ := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()
```

### Beberapa Ekstensi

Gunakan titik koma untuk menentukan beberapa ekstensi dalam satu filter:

```go
dialog := app.Dialog.OpenFile().
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
    AddFilter("Documents", "*.txt;*.doc;*.docx;*.pdf").
    AddFilter("All Files", "*.*")
```

### Format Pola

Gunakan **titik koma** untuk memisahkan beberapa ekstensi dalam satu filter:

```go
// Multiple extensions separated by semicolons
AddFilter("Images", "*.png;*.jpg;*.gif")
```

## Contoh Lengkap

### Membuka File Gambar

```go
func openImage(app *application.App) (image.Image, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
        PromptForSingleSelection()

    if err != nil {
        return nil, err
    }

    if path == "" {
        return nil, errors.New("no file selected")
    }

    // Open and decode image
    file, err := os.Open(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Open Failed").
            SetMessage(err.Error()).
            Show()
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Image").
            SetMessage("Could not decode image file.").
            Show()
        return nil, err
    }

    return img, nil
}
```

### Menyimpan Dokumen dengan Validasi

```go
func saveDocument(app *application.App, content string) {
    path, err := app.Dialog.SaveFile().
        SetFilename("document.txt").
        AddFilter("Text Files", "*.txt").
        AddFilter("Markdown Files", "*.md").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Validate extension
    ext := filepath.Ext(path)
    if ext != ".txt" && ext != ".md" {
        dialog := app.Dialog.Question().
            SetTitle("Confirm Extension").
            SetMessage(fmt.Sprintf("Save as %s file?", ext))

        saveBtn := dialog.AddButton("Save")
        saveBtn.OnClick(func() {
            doSave(app, path, content)
        })

        cancelBtn := dialog.AddButton("Cancel")
        dialog.SetDefaultButton(cancelBtn)
        dialog.SetCancelButton(cancelBtn)
        dialog.Show()
        return
    }

    doSave(app, path, content)
}

func doSave(app *application.App, path, content string) {
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Saved").
        SetMessage("Document saved successfully!").
        Show()
}
```

### Pemrosesan File secara Batch

```go
func processMultipleFiles(app *application.App) {
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil || len(paths) == 0 {
        return
    }

    // Confirm processing
    dialog := app.Dialog.Question().
        SetTitle("Confirm Processing").
        SetMessage(fmt.Sprintf("Process %d file(s)?", len(paths)))

    processBtn := dialog.AddButton("Process")
    processBtn.OnClick(func() {
        // Process files
        var errs []error
        for i, path := range paths {
            if err := processFile(path); err != nil {
                errs = append(errs, err)
            }

            // Update progress
            // app.Event.Emit("progress", map[string]interface{}{
            //     "current": i + 1,
            //     "total":   len(paths),
            // })
            _ = i // suppress unused variable warning in example
        }

        // Show results
        if len(errs) > 0 {
            app.Dialog.Warning().
                SetTitle("Processing Complete").
                SetMessage(fmt.Sprintf("Processed %d files with %d errors.",
                    len(paths), len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Success").
                SetMessage(fmt.Sprintf("Processed %d files successfully!", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Mengekspor dengan Pemilihan Folder

```go
func exportData(app *application.App, data []byte) {
    // Select output folder
    folder, err := app.Dialog.OpenFile().
        SetTitle("Select Export Folder").
        SetDirectory(getDefaultExportFolder()).
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil || folder == "" {
        return
    }

    // Generate filename
    filename := fmt.Sprintf("export_%s.csv",
        time.Now().Format("2006-01-02_15-04-05"))
    path := filepath.Join(folder, filename)

    // Save file
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Export Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Show success with option to open folder
    dialog := app.Dialog.Question().
        SetTitle("Export Complete").
        SetMessage(fmt.Sprintf("Exported to %s", filename))

    openBtn := dialog.AddButton("Open Folder")
    openBtn.OnClick(func() {
        openFolder(folder)
    })

    dialog.AddButton("OK")
    dialog.Show()
}
```

### Mengimpor dengan Validasi

```go
func importConfiguration(app *application.App) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Import Configuration").
        AddFilter("JSON Files", "*.json").
        AddFilter("YAML Files", "*.yaml;*.yml").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Read Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Validate configuration
    config, err := parseConfig(data)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Configuration").
            SetMessage("File is not a valid configuration.").
            Show()
        return
    }

    // Confirm import
    dialog := app.Dialog.Question().
        SetTitle("Confirm Import").
        SetMessage("Import this configuration?")

    importBtn := dialog.AddButton("Import")
    importBtn.OnClick(func() {
        // Apply configuration
        if err := applyConfig(config); err != nil {
            app.Dialog.Error().
                SetTitle("Import Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        app.Dialog.Info().
            SetTitle("Success").
            SetMessage("Configuration imported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Sediakan filter file** - Bantu pengguna menemukan file
- **Tetapkan judul yang sesuai** - Berikan konteks yang jelas
- **Gunakan direktori default** - Mulai dari lokasi yang logis
- **Validasi pilihan** - Periksa jenis file
- **Tangani pembatalan** - Pengguna mungkin membatalkan
- **Tampilkan konfirmasi** - Untuk tindakan destruktif
- **Berikan umpan balik** - Pesan keberhasilan/kesalahan

### ❌ Jangan Lakukan

- **Jangan lewati validasi** - Periksa jenis file
- **Jangan abaikan kesalahan** - Tangani pembatalan
- **Jangan gunakan filter generik** - Gunakan filter yang spesifik
- **Jangan lupakan "Semua File"** - Selalu sertakan sebagai opsi
- **Jangan tulis jalur secara hardcode** - Gunakan direktori beranda pengguna
- **Jangan berasumsi bahwa file tersedia** - Periksa sebelum membukanya

## Perbedaan Antarplatform

### macOS

- NSOpenPanel/NSSavePanel native
- Bergaya sheet saat dilampirkan ke jendela
- Mengikuti tema sistem
- Mendukung pratinjau Quick Look
- Integrasi dengan tag dan favorit

### Windows

- Dialog Buka/Simpan File native
- Mengikuti tema sistem
- Integrasi file terbaru
- Dukungan lokasi jaringan

### Linux

- Pemilih file GTK
- Bervariasi menurut lingkungan desktop
- Mengikuti tema desktop
- Dukungan file terbaru

## Langkah Berikutnya

@cards{cols="2"}
ℹ Dialog pesan
Dialog informasi, peringatan, dan kesalahan.

[Pelajari Selengkapnya →](/features/dialogs/message/)

---
◆ Dialog khusus
Buat jendela dialog khusus.

[Pelajari Selengkapnya →](/features/dialogs/custom/)

---
🚀 Binding
Panggil fungsi Go dari JavaScript.

[Pelajari Selengkapnya →](/features/bindings/methods/)

---
★ Peristiwa
Gunakan peristiwa untuk memperbarui progres.

[Pelajari Selengkapnya →](/features/events/system/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh dialog file](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
