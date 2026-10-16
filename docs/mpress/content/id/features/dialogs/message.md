---
title: "Dialog pesan"
description: "Tampilkan informasi, peringatan, kesalahan, dan pertanyaan"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## Dialog pesan

Wails menyediakan **dialog pesan native** dengan tampilan yang sesuai untuk setiap platform: dialog informasi, peringatan, kesalahan, dan pertanyaan dengan judul, pesan, serta tombol yang dapat disesuaikan. API sederhana, perilaku native, dan aksesibilitas aktif secara default.

## Membuat Dialog

Dialog pesan diakses melalui pengelola `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

Semua metode mengembalikan `*MessageDialog` yang dapat dikonfigurasi menggunakan perangkaian metode.

## Dialog informasi

Tampilkan pesan informatif:

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**Kasus penggunaan:**

- Konfirmasi keberhasilan
- Pemberitahuan penyelesaian
- Pesan informatif
- Pembaruan status

**Contoh—konfirmasi penyimpanan:**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## Dialog peringatan

Tampilkan peringatan:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Kasus penggunaan:**

- Peringatan nonkritis
- Pemberitahuan penghentian penggunaan
- Pesan peringatan
- Potensi masalah

**Contoh—peringatan ruang disk:**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## Dialog kesalahan

Tampilkan kesalahan:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**Kasus penggunaan:**

- Pesan kesalahan
- Pemberitahuan kegagalan
- Penanganan pengecualian
- Masalah kritis

**Contoh—kesalahan jaringan:**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## Dialog pertanyaan

Ajukan pertanyaan kepada pengguna dan tangani respons melalui fungsi callback tombol:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Kasus penggunaan:**

- Mengonfirmasi tindakan
- Pertanyaan Ya/Tidak
- Pilihan ganda
- Keputusan pengguna

**Contoh—perubahan yang belum disimpan:**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## Opsi dialog

### Judul dan Pesan

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**Praktik terbaik:**

- **Judul:** Singkat dan deskriptif (2-5 kata)
- **Pesan:** Jelas, spesifik, dan dapat ditindaklanjuti
- **Hindari jargon:** Gunakan bahasa yang sederhana

### Tombol

**Satu tombol (Informasi/Peringatan/Kesalahan):**

Dialog informasi, peringatan, dan kesalahan menampilkan tombol "OK" secara default:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

Anda juga dapat menambahkan tombol khusus:

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**Beberapa tombol (Pertanyaan):**

Gunakan `AddButton()` untuk menambahkan tombol. Metode ini mengembalikan `*Button` yang dapat Anda konfigurasi:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

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
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

Anda juga dapat menggunakan metode fluent `SetAsDefault()` dan `SetAsCancel()` pada tombol:

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**Praktik terbaik:**

- **1-3 tombol:** Jangan membebani pengguna dengan terlalu banyak pilihan
- **Label yang jelas:** "Simpan", bukan "OK"
- **Pilihan default yang aman:** Tindakan yang tidak merusak
- **Urutan itu penting:** Dahulukan tindakan yang paling mungkin dipilih (kecuali Batal)

### Ikon Khusus

Tetapkan ikon khusus untuk dialog:

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### Penautan ke Jendela

Tautkan ke jendela tertentu:

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**Manfaat:**

- Dialog muncul di jendela yang tepat
- Jendela induk dinonaktifkan selama dialog ditampilkan
- Pengalaman pengguna yang lebih baik saat menggunakan beberapa jendela

## Contoh Lengkap

### Mengonfirmasi Tindakan Destruktif

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Konfirmasi Keluar

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### Dialog Pembaruan dengan Opsi Unduh

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Pertanyaan dengan Ikon Khusus

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Berikan informasi spesifik** - "File disimpan ke Dokumen", bukan "Berhasil"
- **Gunakan jenis yang sesuai** - Error untuk kesalahan, Warning untuk peringatan
- **Berikan konteks** - Sertakan detail yang relevan
- **Gunakan label tombol yang jelas** - "Hapus", bukan "OK"
- **Tetapkan pilihan default yang aman** - Gunakan tindakan yang tidak merusak
- **Tangani pembatalan** - Pengguna mungkin menutup dialog

### ❌ Jangan Lakukan

- **Jangan gunakan secara berlebihan** - Mengganggu alur kerja
- **Jangan gunakan untuk pembaruan yang sering** - Gunakan notifikasi sebagai gantinya
- **Jangan gunakan pesan umum** - "Kesalahan" tidak memberikan informasi apa pun
- **Jangan abaikan kesalahan** - Tangani kesalahan dari dialog.Show()
- **Jangan memblokir tanpa perlu** - Pertimbangkan alternatif asinkron
- **Jangan gunakan jargon teknis** - Gunakan bahasa yang mudah dipahami

## Perbedaan Antarplatform

### macOS

- Bergaya sheet saat dilampirkan ke jendela
- Pintasan papan ketik standar (⌘. untuk Batal)
- Secara otomatis mengikuti tema sistem
- Fitur aksesibilitas tersedia secara bawaan

### Windows

- Dialog modal
- Tampilan TaskDialog
- Esc untuk Batal
- Mengikuti tema sistem

### Linux

- Dialog GTK
- Bervariasi menurut lingkungan desktop
- Mengikuti tema desktop
- Navigasi papan ketik standar

## Langkah Berikutnya

@cards{cols="2"}
📖 Dialog file
Membuka dan menyimpan file, serta memilih folder.

[Pelajari Selengkapnya →](/features/dialogs/file/)

---
◆ Dialog khusus
Buat jendela dialog khusus.

[Pelajari Selengkapnya →](/features/dialogs/custom/)

---
● Notifikasi
Notifikasi yang tidak mengganggu.

[Pelajari Selengkapnya →](/features/notifications/overview/)

---
★ Event
Gunakan event untuk komunikasi non-pemblokiran.

[Pelajari Selengkapnya →](/features/events/system/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh dialog](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
