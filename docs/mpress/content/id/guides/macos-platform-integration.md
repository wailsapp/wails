---
title: "Integrasi Platform macOS"
description: "Jendela dokumen, tambahan Dock dan menu, panel native, item status, umpan balik, papan klip kaya format, penyeretan keluar, izin, daya, dan siklus hidup pada macOS"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

Platform yang relevan: macOS

Wails v3 memberi aplikasi Anda perilaku yang diharapkan pengguna macOS dari aplikasi native: jendela dokumen dengan ikon proksi dan penempatan berjenjang, menu dengan simbol dan lencana, menu Dock dengan progres, peringatan dan panel native, item status yang dapat dihapus, umpan balik haptik dan ucapan, papan klip kaya format dengan penyeretan keluar, informasi sistem tentang izin, daya, dan lokal, serta dukungan di menu Services, Handoff, AppleScript, dan Quick Look. Semuanya dikendalikan dari Go melalui paket `application`.

Kode yang sama dapat dikompilasi di Windows dan Linux. Setter menyimpan nilainya, kueri mengembalikan nilai nol, dan operasi yang memerlukan macOS mengembalikan kesalahan terdokumentasi seperti `ErrMacOnly`, `ErrDialogNotSupported`, atau `ErrClipboardNotSupported`. [Catatan platform](#platform-notes) di bawah mencantumkan perilaku setiap area di luar macOS.

Untuk antarmuka jendela native (toolbar, sidebar, panel pemeriksa, aksesori, dan tab jendela), lihat panduan [Antarmuka Jendela Native macOS](/guides/macos-native-chrome).

## Jendela dokumen

Jendela dokumen menampilkan berkas yang diwakilinya pada titlebar, menandai perubahan yang belum disimpan dengan titik pada tombol tutup, dan membuka jendela baru secara berjenjang. Semuanya tersedia sebagai metode pada `WebviewWindow` dan opsi pada `MacWindow`.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Report.md",
    URL:             "/editor",
    InitialPosition: application.WindowCascade,
    Mac: application.MacWindow{
        FrameAutosaveName: "editor.main",
        TitleBar: application.MacTitleBar{
            WindowButtonsOffset: &application.Point{X: 12, Y: 8},
        },
    },
})

window.SetRepresentedFile("/Users/me/Documents/Report.md")
window.SetSubtitle("Documents")
window.SetDocumentEdited(true)
```

- `SetRepresentedFile` menampilkan ikon proksi berkas pada titlebar. Pengguna dapat menyeret ikon ke aplikasi lain atau mengkliknya sambil menekan tombol Command untuk menampilkan jalurnya. Teruskan `""` untuk menghapusnya. `RepresentedFile` membaca nilainya kembali.
- `SetDocumentEdited` menampilkan titik perubahan yang belum disimpan pada tombol tutup dan meredupkan ikon proksi. `IsDocumentEdited` membaca statusnya kembali.
- `SetSubtitle` menampilkan baris kedua di bawah judul pada macOS 11 dan yang lebih baru.
- `InitialPosition: application.WindowCascade` menempatkan jendela di bawah dan di kanan jendela berjenjang terakhir, seperti saat dokumen baru dibuka. `CascadeFrom(other)` melakukan hal yang sama untuk jendela yang sudah ada dan memperbarui titik penempatan berjenjang untuk jendela berikutnya.
- `Mac.FrameAutosaveName` memulihkan posisi dan ukuran yang tersimpan sebelum jendela pertama kali ditampilkan, lalu terus menyimpannya saat jendela dipindahkan. Bingkai yang dipulihkan lebih diutamakan daripada `X`, `Y`, `Width`, `Height`, dan `InitialPosition`. `SetFrameAutosaveName` mengganti namanya pada jendela yang aktif.
- `MacTitleBar.WindowButtonsOffset` menggeser tombol tutup, minimalkan, dan perbesar sebesar sejumlah poin. `SetWindowButtonsOffset` dan `ResetWindowButtonsOffset` mengubahnya saat runtime.

Ketiga setter tersebut dapat dipanggil sebelum jendela native ada; nilainya diterapkan ketika jendela dibuat.

### Permintaan perhatian

`RequestAttention` membuat ikon Dock memantul saat aplikasi berada di latar belakang. Permintaan informasional memantul sekali. Permintaan kritis terus memantul hingga pengguna mengaktifkan aplikasi atau Anda membatalkannya.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` tetap menjadi cara lintas platform untuk meminta perhatian satu kali.

### Pencetakan dan ekspor

`PrintWithOptions` mencetak WebView dengan pengaturan halaman yang eksplisit. Nilai nol menampilkan panel cetak dengan pengaturan cetak bersama. `Print` mempertahankan perilaku lamanya (orientasi lanskap, margin 30 poin).

`ExportPDF` merender halaman menjadi dokumen PDF dan `Snapshot` menangkapnya sebagai PNG. Keduanya menunggu WebKit, jadi panggil dari goroutine dan jangan pernah dari thread aplikasi; pemanggilan di thread tersebut mengembalikan `ErrMacExportOnMainThread`.

```go
err := window.PrintWithOptions(application.PrintOptions{
    Orientation: application.PrintOrientationPortrait,
    Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
    Silent:      false,
})
if err != nil {
    log.Println("print:", err)
}

go func() {
    pdf, err := window.ExportPDF(application.PDFExportOptions{})
    if err != nil {
        log.Println("export:", err)
        return
    }
    if err := os.WriteFile("report.pdf", pdf, 0o644); err != nil {
        log.Println("write:", err)
    }

    png, err := window.Snapshot(application.SnapshotOptions{Width: 800})
    if err != nil {
        log.Println("snapshot:", err)
        return
    }
    if err := os.WriteFile("preview.png", png, 0o644); err != nil {
        log.Println("write:", err)
    }
}()
```

`PrintOptions` juga menerima `PrinterName`, `PaperName` (nama PostScript seperti `"iso-a4"`), dan `Scale`. `PDFExportOptions` dan `SnapshotOptions` menerima `Rect` opsional untuk membatasi area tangkapan dan `Timeout` yang secara bawaan bernilai `DefaultMacExportTimeout` (30 detik).

## Lembar

Lembar adalah jendela kedua yang melekat pada bagian atas jendela induknya, seperti panel Simpan. Setiap `WebviewWindow` dapat ditampilkan sebagai lembar milik jendela lain dengan `PresentSheet`, lalu diakhiri dengan `EndSheet` dan kode respons yang diteruskan ke callback `OnSheetEnd` milik lembar tersebut. Buat jendela lembar dengan `Hidden` diatur agar tidak muncul sekilas di layar sebelum dipasang; AppKit menyembunyikannya lagi ketika selesai, sehingga jendela yang sama dapat ditampilkan berulang kali.

```go
sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Rename",
    URL:    "/rename",
    Width:  420,
    Height: 180,
    Hidden: true,
})
sheet.OnSheetEnd(func(code int) {
    if code == application.MacSheetResponseOK {
        log.Println("renamed")
    }
})

// From a menu item or a bound method:
if err := window.PresentSheet(sheet); err != nil {
    log.Println(err)
}

// From the sheet's own page, through a bound method:
sheet.EndSheet(application.MacSheetResponseOK)
```

`PresentCriticalSheet` menampilkan lembar di depan lembar biasa yang sudah terpasang, alih-alih mengantrekannya di belakang. `PresentNativeSheet` melakukan hal yang sama untuk `NativeWindow`. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet`, dan `HasAttachedSheet` menjelaskan status saat ini. Menutup jendela lembar alih-alih mengakhirinya menghasilkan `MacSheetResponseStop`.

## Popover

`MacPopover` adalah `NSPopover`: panel sementara yang ditambatkan ke persegi panjang dalam jendela, item toolbar, atau item status pada bilah menu. Kontennya berupa bilah kontrol `MacAccessory` native, jenis yang sama dengan yang digunakan panduan [Antarmuka Jendela Native macOS](/guides/macos-native-chrome) untuk aksesori titlebar. Tambahkan semua kontrol sebelum pertama kali menampilkannya.

```go
strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
strip.AddLabel("Sort by")
strip.AddSegmented([]string{"Date", "Title"}, 0)
strip.AddFlexibleSpace()
strip.AddButton("Apply").OnClick(func(*application.Context) {
    // apply the sort
})

popover := application.NewMacPopover(application.MacPopoverOptions{
    Width:    320,
    Behavior: application.MacPopoverBehaviorTransient,
    Content:  strip,
})
popover.OnClose(func() {
    log.Println("popover closed")
})

// Anchored to a rectangle in the page, in window content coordinates:
err := popover.ShowRelativeTo(application.Rect{X: 20, Y: 60, Width: 120, Height: 28}, window, application.MacRectEdgeMaxY)
if err != nil {
    log.Println(err)
}
```

`MacToolbarItem.ShowPopover` dan `SystemTray.ShowPopover` menambatkan popover yang sama ke item toolbar atau item status. `MacPopoverBehaviorTransient` menutup popover saat ada klik di luar popover, `Semitransient` hanya saat ada klik di jendela yang menampilkannya, dan perilaku bawaan mempertahankannya tetap terbuka sampai `Close` dipanggil. `MacRectEdge` memilih sisi tempat popover muncul. `SetContentSize` dan `SetBehavior` menyesuaikan popover yang aktif, sedangkan `Destroy` melepas popover native dan membebaskan bilah kontennya untuk digunakan di tempat lain.

## Pemulihan status

macOS membuka kembali jendela aplikasi setelah aplikasi mengalami crash, dihentikan secara paksa, atau perangkat dimulai ulang, serta setelah aplikasi ditutup secara normal jika "Close windows when quitting an application" dinonaktifkan di System Settings. Berikan `Mac.RestorationID` pada jendela, simpan data yang diperlukan untuk membuatnya kembali dengan `SetRestorationData`, dan daftarkan `app.Window.OnRestore` untuk membangun kembali jendela pada peluncuran berikutnya.

```go
app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
    if id != "editor" {
        return nil
    }
    path := state.Get("path")
    restored := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: filepath.Base(path),
        URL:   "/editor?path=" + path,
        Mac:   application.MacWindow{RestorationID: id},
    })
    restored.SetRestorationData(state.Data)
    return restored
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Report.md",
    URL:   "/editor?path=/Users/me/Documents/Report.md",
    Mac:   application.MacWindow{RestorationID: "editor"},
})
window.SetRestorationData(map[string]string{"path": "/Users/me/Documents/Report.md"})
```

Hanya jendela yang terlihat saat aplikasi berakhir yang disimpan. `RestorationState.Data` adalah map string; gunakan hanya untuk pengidentifikasi, jalur, dan posisi. `InteractionState` mengembalikan daftar navigasi mundur-maju dan posisi gulir WebView sebagai blob opak (macOS 12+) yang diterapkan oleh `RestoreInteractionState` pada jendela yang dibuat kembali, biasanya disimpan dalam data pemulihan dengan enkode base64. `SetRestorationID` dan `RestorationID` mengubah dan membaca pengidentifikasi pada jendela yang aktif.

## Opsi presentasi

`MacPresentationOptions` mencerminkan `NSApplication.presentationOptions`: bitmask yang menyembunyikan Dock atau bilah menu dan menonaktifkan perpindahan proses, penghentian paksa, keluar dari sesi, atau perintah Sembunyikan saat aplikasi aktif. Atur `Mac.PresentationOptions` dalam opsi aplikasi untuk menerapkannya saat peluncuran, atau ubah saat runtime dengan `SetPresentationOptions`. Kombinasi yang tidak valid ditolak sebelum mencapai AppKit dengan kesalahan yang membungkus `ErrMacPresentationOptionsInvalid`.

```go
kiosk := application.MacPresentationHideDock |
    application.MacPresentationHideMenuBar |
    application.MacPresentationDisableProcessSwitching |
    application.MacPresentationDisableForceQuit

if err := app.SetPresentationOptions(kiosk); err != nil {
    log.Println(err)
}
log.Println("presentation:", app.PresentationOptions())

// Restore the standard Dock and menu bar:
_ = app.SetPresentationOptions(application.MacPresentationDefault)
```

Menyembunyikan bilah menu (`HideMenuBar` atau `AutoHideMenuBar`) memerlukan salah satu opsi Dock, sedangkan `AutoHideToolbar` memerlukan `FullScreen` dan `AutoHideMenuBar`. `Validate` melaporkan aturan pertama yang dilanggar suatu nilai dan `Has` menguji setiap flag.

## Menu dan Dock

Item menu mendapat dukungan SF Symbols, lencana, header bagian, palet warna, status centang campuran, item alternatif, dan indentasi. Semuanya adalah metode pada `MenuItem` dan `Menu`, sehingga dapat digunakan dalam menu aplikasi, menu konteks, menu tray, dan menu Dock.

```go
menu := app.NewMenu()
menu.AddRole(application.AppMenu)

fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Open...").SetSymbol("folder").SetAccelerator("CmdOrCtrl+o").
    OnClick(func(*application.Context) {
        // open the file, then note it in the recent list:
        app.Menu.AddRecentDocument("/Users/me/Documents/Report.md")
    })
fileMenu.AddRole(application.OpenRecent)

view := menu.AddSubmenu("View")
view.AddSectionHeader("Mailboxes")
inbox := view.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
inbox.OnClick(func(*application.Context) {
    inbox.ClearBadge()
})
view.Add("Updates").SetBadgeText("New")

view.AddSeparator()
wrap := view.AddCheckbox("Wrap lines", false).SetMixed()
view.Add("Reset").SetIndentationLevel(1).OnClick(func(*application.Context) {
    wrap.SetMixed()
})

view.AddSeparator()
view.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").OnClick(func(*application.Context) {})
view.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).
    OnClick(func(*application.Context) {})

colours := []application.RGBA{
    application.NewRGB(255, 59, 48),
    application.NewRGB(255, 149, 0),
    application.NewRGB(52, 199, 89),
}
view.AddPalette([]string{"tag.fill"}, colours, 0, func(_ *application.Context, index int) {
    log.Println("tag colour", index)
}).SetLabel("Tag colour")

app.Menu.Set(menu)
```

- `SetSymbol` menampilkan SF Symbol di sebelah judul (macOS 11+). Simbol ini menggantikan gambar yang diatur dengan `SetBitmap`.
- `SetBadge` menampilkan jumlah dan `SetBadgeText` menampilkan string pendek setelah judul (macOS 14+). `ClearBadge` menghapusnya; `BadgeCount` dan `BadgeText` membaca nilainya kembali.
- `AddSectionHeader` menambahkan header yang tidak interaktif (macOS 14+). Pada rilis sebelumnya, header menjadi item nonaktif dengan judul yang sama.
- `AddPalette` menambahkan deretan sampel warna yang didukung oleh menu palet `NSMenu` (macOS 14+). Teruskan satu simbol untuk setiap sampel, satu per warna, atau slice kosong untuk lingkaran berisi. Tanpa label, palet tampil langsung dalam menu induk; `SetLabel` menampilkannya sebagai submenu berjudul. `PaletteSelected` mengembalikan indeks yang dipilih.
- `SetMixed` menempatkan kotak centang dalam status campuran yang digambar sebagai garis pendek. Klik akan mengaktifkannya sepenuhnya, sesuai perilaku AppKit.
- `SetAlternate(true)` menampilkan item sebagai pengganti item di atasnya selama tombol pengubah yang membedakannya ditekan. Kedua item harus memiliki tombol yang sama dan berbeda dalam tombol pengubah.
- `SetIndentationLevel` memberi indentasi pada judul hingga 15 tingkat.

### Buka Terkini

`fileMenu.AddRole(application.OpenRecent)` menambahkan submenu Buka Terkini standar. Pada macOS, `NSDocumentController` mengisinya setiap kali dibuka dan menyertakan item Hapus Menu. Tambahkan berkas dengan `app.Menu.AddRecentDocument`, lihat daftarnya dengan `RecentDocuments`, dan kosongkan dengan `ClearRecentDocuments`. Daftar ini tetap ada setelah aplikasi dimulai ulang.

Pemilihan berkas terkini menghasilkan event yang sama seperti berkas yang dibuka dari Finder:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Menu Dock

`app.Menu.SetDockMenu` memasang menu statis yang muncul saat ikon Dock diklik kanan. `OnDockMenu` membuat menu sesuai kebutuhan setiap kali menu akan ditampilkan, sehingga cocok jika itemnya mencerminkan status yang berubah. Pembuat menu lebih diutamakan daripada menu statis.

```go
dockMenu := application.NewMenu()
dockMenu.Add("New Document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
    // create a document
})
app.Menu.SetDockMenu(dockMenu)

app.Menu.OnDockMenu(func() *application.Menu {
    dynamic := application.NewMenu()
    dynamic.AddSectionHeader("Recent")
    for _, path := range app.Menu.RecentDocuments() {
        path := path
        dynamic.Add(filepath.Base(path)).OnClick(func(*application.Context) {
            app.Menu.OpenRecentDocument(path)
        })
    }
    return dynamic
})
```

### Progres Dock

Layanan dock menggambar bilah progres di atas ikon Dock, di samping dukungan lencana yang sudah ada. Daftarkan `dock.New()` sebagai layanan dan panggil `SetProgress` dengan pecahan antara 0 dan 1.

```go
import "github.com/wailsapp/wails/v3/pkg/services/dock"

dockService := dock.New()

app := application.New(application.Options{
    Name: "Exporter",
    Services: []application.Service{
        application.NewService(dockService),
    },
})

app.Event.On("export:progress", func(event *application.CustomEvent) {
    if fraction, ok := event.Data.(float64); ok {
        _ = dockService.SetProgress(fraction)
    }
})
app.Event.On("export:done", func(*application.CustomEvent) {
    _ = dockService.ClearProgress()
})
```

`GetProgress` mengembalikan pecahan saat ini, atau `nil` jika tidak ada bilah yang ditampilkan.

## Dialog

Dialog pesan, buka, dan simpan menerima opsi macOS, sedangkan pengelola dialog mendapat prompt teks serta panel warna dan font sistem.

### Peringatan

`SetSuppression` menambahkan kotak centang "Jangan tampilkan pesan ini lagi" dan `SetHelp` menampilkan tombol bantuan. Baca status kotak centang dengan `Suppressed` dari callback tombol, atau daftarkan `OnSuppression` untuk menerimanya terlebih dahulu.

```go
dialog := app.Dialog.Warning().
    SetTitle("Delete 3 items?").
    SetMessage("The items will be moved to the Bin.").
    SetSuppression("Do not warn me again").
    SetHelp(func() {
        log.Println("help requested")
    }).
    AttachToWindow(window)

dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
    if dialog.Suppressed() {
        // remember not to ask again
    }
})
dialog.AddButton("Cancel").SetAsCancel()
dialog.Show()
```

### Prompt teks

`Prompt` menampilkan peringatan dengan bidang teks dan memblokir hingga peringatan ditutup, jadi panggil dari goroutine atau metode yang diikat. `Secure` mengubah bidang tersebut menjadi bidang kata sandi dan `Window` menampilkan peringatan sebagai lembar.

```go
go func() {
    value, ok, err := app.Dialog.Prompt(application.PromptOptions{
        Title:        "Name this document",
        Message:      "The name is used for the exported file.",
        Placeholder:  "Untitled",
        DefaultValue: "Quarterly report",
        OKLabel:      "Rename",
        Window:       window,
    })
    if err != nil || !ok {
        return
    }
    log.Println("renamed to", value)
}()
```

### Panel berkas

`AddContentType` memfilter berdasarkan pengidentifikasi tipe seragam pada dialog buka maupun simpan. Metode ini digunakan bersama `AddFilter`, sehingga `"public.image"` cocok dengan setiap tipe gambar yang dikenal sistem, sementara filter tetap dapat menangkap PDF berdasarkan ekstensi.

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Choose images").
    AddFilter("PDF", "*.pdf").
    AddContentType("public.image").
    AttachToWindow(window).
    PromptForMultipleSelection()
if err == nil {
    log.Println(paths)
}
```

Panel Simpan mendapat menu tarik turun Format, label khusus untuk bidang nama, dan tag Finder. `SetFormats` mengganti tipe yang diizinkan dan ekstensi pada bidang nama saat pengguna mengubah pilihan dalam menu, sedangkan `SelectedFormat` melaporkan pilihan akhir.

```go
formats := []application.DialogFormat{
    {Label: "PNG image", Extension: "png", UTI: "public.png"},
    {Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
}
save := app.Dialog.SaveFile().
    SetFilename("Quarterly report.png").
    SetNameFieldLabel("Export As:").
    SetTags([]string{"Reports", "Draft"}).
    AttachToWindow(window)
save.SetFormats(formats, 0, nil)

path, err := save.PromptForSingleSelection()
if err == nil && path != "" {
    log.Println("export", path, "as", formats[save.SelectedFormat()].Label)
}
```

### Panel warna dan font

`PickColor` dan `PickFont` membuka panel sistem bersama dan memblokir hingga panel ditutup. `OnChange` mengirimkan setiap pilihan selama panel terbuka, sehingga halaman dapat langsung menampilkan pratinjau pilihan. Hanya satu panel dari setiap jenis yang dapat terbuka pada satu waktu; pemanggilan kedua mengembalikan `ErrDialogInProgress`.

```go
go func() {
    colour, changed, err := app.Dialog.PickColor(application.ColorPickerOptions{
        Initial:    application.NewRGB(52, 120, 246),
        ShowsAlpha: true,
        Title:      "Accent colour",
        OnChange: func(colour application.RGBA) {
            app.Event.Emit("accent:preview", colour)
        },
    })
    if err == nil && changed {
        log.Println("accent", colour)
    }

    font, changed, err := app.Dialog.PickFont(application.FontPickerOptions{
        Family: "Helvetica Neue",
        Size:   18,
    })
    if err == nil && changed {
        log.Println(font.Family, font.Face, font.PostScriptName, font.Size)
    }
}()
```

## Item status dan umpan balik

### Item status

Item tray sistem pada macOS adalah `NSStatusItem`. Item ini dapat digambar dari SF Symbol, memiliki tooltip, dan dihapus oleh pengguna dengan cara yang sama seperti item bawaan.

```go
tray := app.SystemTray.New()
tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
tray.SetTooltip("Recorder. Command-drag to remove.")
tray.SetRemovable(true, "com.example.recorder.tray")
tray.OnVisibilityChange(func(visible bool) {
    log.Println("status item visible:", visible)
})

trayMenu := app.Menu.New()
trayMenu.Add("Show").OnClick(func(*application.Context) {
    window.Show().Focus()
})
tray.SetMenu(trayMenu)
tray.Run()
```

- `SetSymbol` merender simbol sebagai gambar templat agar mengikuti tampilan bilah menu (macOS 11+). `SetSymbolConfiguration` mengatur ukuran dalam poin dan ketebalannya.
- `SetTooltip` mengatur teks saat penunjuk berada di atas item. `Tooltip` membaca nilainya kembali.
- `SetRemovable(true, name)` memungkinkan pengguna menyeret item keluar dari bilah menu sambil menekan tombol Command. Berikan nama penyimpanan otomatis yang stabil agar macOS mengingat penghapusan tersebut setelah aplikasi dijalankan ulang. `Show` atau `SetVisible(true)` menampilkannya kembali.
- `IsVisible` membaca `NSStatusItem.visible`, sehingga bernilai false setelah pengguna menghapus item. `OnVisibilityChange` melaporkan setiap perubahan.

### Umpan balik haptik

`app.Haptics.Perform` memainkan pola pada trackpad Force Touch atau Magic Trackpad selama aplikasi aktif.

```go
app.Haptics.Perform(application.HapticAlignment)
```

Jenisnya adalah `HapticGeneric`, `HapticAlignment` (item terkunci pada posisinya), dan `HapticLevelChange` (titik tahanan atau tahapan klik). `IsSupported` melaporkan apakah platform dapat menghasilkan umpan balik sama sekali.

### Suara

`app.Sound` memainkan suara peringatan, suara sistem bernama, atau berkas audio.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` menerima nama dari `SystemSounds` atau jalur absolut ke berkas apa pun yang dapat didekode oleh Core Audio. `PlayData` memainkan berkas audio lengkap dari memori.

### Ucapan

`app.Speech.Speak` mengantrekan teks untuk diucapkan oleh suara sistem dan mengembalikan `Utterance`. Ucapan diputar satu demi satu; `Stop` membatalkan satu ucapan dan `StopAll` mengosongkan antrean. `Voices` mencantumkan suara yang terpasang beserta pengidentifikasi dan bahasanya.

```go
utterance, err := app.Speech.Speak("Export finished", application.SpeechOptions{
    Voice: "com.apple.voice.compact.en-GB.Daniel",
    Rate:  0.5,
})
if err != nil {
    log.Println(err)
    return
}
utterance.OnFinished(func() {
    log.Println("done, stopped:", utterance.WasStopped())
})
```

`Recognize` mentranskripsikan masukan dari mikrofon bawaan dengan `SFSpeechRecognizer`. Pemanggilan pertama meminta izin mikrofon dan pengenalan ucapan serta memblokir hingga pengguna menjawab, jadi panggil dari goroutine. Transkrip sementara diterima melalui `OnPartial`; `Stop` mengakhiri perekaman dan mengembalikan teks akhir.

```go
go func() {
    session, err := app.Speech.Recognize(application.RecognitionOptions{
        Locale: "en-US",
        OnPartial: func(text string) {
            app.Event.Emit("dictation:partial", text)
        },
    })
    if err != nil {
        if errors.Is(err, application.ErrSpeechRecognitionDenied) {
            _ = app.Permissions.OpenSystemSettings(application.PermissionKindMicrophone)
        }
        log.Println(err)
        return
    }
    time.Sleep(5 * time.Second)
    text, err := session.Stop()
    if err != nil {
        log.Println(err)
        return
    }
    app.Event.Emit("dictation:final", text)
}()
```

Pengenalan ucapan memerlukan aplikasi dalam bundle yang `Info.plist`-nya mendeklarasikan `NSSpeechRecognitionUsageDescription` dan `NSMicrophoneUsageDescription`. Tanpa keduanya, macOS menolak akses dan `Recognize` mengembalikan `ErrSpeechRecognitionUsageDescription`.

## Papan klip dan penyeretan

### Papan klip kaya format

`app.Clipboard` membaca dan menulis gambar, referensi berkas, HTML, RTF, dan data mentah dengan pengidentifikasi tipe seragam apa pun, selain teks biasa. `Types` mencantumkan isi papan klip dan `OnChange` melaporkan perubahan yang dibuat oleh aplikasi mana pun.

```go
if err := app.Clipboard.SetHTML("<p>Rich <b>HTML</b></p>", "Rich HTML"); err != nil {
    log.Println(err)
}
_ = app.Clipboard.SetFiles([]string{"/Users/me/Documents/Report.md"})
_ = app.Clipboard.SetData("com.example.record", []byte(`{"id":42}`))

stop := app.Clipboard.OnChange(func() {
    log.Println("clipboard changed", app.Clipboard.ChangeCount(), app.Clipboard.Types())
    if files, err := app.Clipboard.Files(); err == nil && len(files) > 0 {
        log.Println("files:", files)
    }
})
defer stop()
```

`SetImage` dan `Image` menggunakan byte PNG; gambar yang disalin sebagai TIFF oleh aplikasi lain dikonversi untuk Anda. Tidak ada notifikasi sistem untuk perubahan papan klip, sehingga `OnChange` memeriksa jumlah perubahan setiap 500 ms selama setidaknya ada satu pendengar.

### Menyeret keluar

`StartDrag` memulai penyeretan sistem dari jendela, seolah-olah pengguna mengambil item dari Finder. Operasi ini menawarkan berkas yang sudah ada, janji berkas yang isinya baru dibuat saat tujuan menerima pelepasan, atau teks biasa. Mulai saat gestur mouse: ikat metode Go dan panggil dari penangan `mousedown` atau `pointerdown` pada elemen yang dapat diseret di halaman, dengan atribut HTML `draggable` diatur ke `false` agar WebKit tidak memulai penyeretannya sendiri.

```go
// DragService is bound to the page and called from a mousedown handler.
type DragService struct {
    window *application.WebviewWindow
}

func (s *DragService) DragExport() error {
    return s.window.StartDrag(application.DragItems{
        Promises: []application.DragPromise{{
            Filename: "export.csv",
            Data: func() ([]byte, error) {
                return []byte("id,name\n1,Wails\n"), nil
            },
        }},
        Operations: application.DragOperationCopy,
    })
}
```

```go
window.OnDragEnd(func(operation application.DragOperation) {
    log.Println("drag ended:", operation)
})
```

Memanggil `StartDrag` di luar gestur mengembalikan `ErrDragOutNoGesture`. `DragItems.Image` dan `ImageOffset` mengatur gambar di bawah kursor.

### Item yang dilepas dari aplikasi lain

Pelepasan berkas tetap menggunakan event `WindowFilesDropped`. Untuk menerima teks, URL, atau gambar yang diseret dari aplikasi lain, cantumkan tipenya dalam `DropTypes` dan daftarkan `OnDrop`. Item yang dilepas tersebut dikirim ke Go, bukan ke penangan pelepasan HTML5 milik halaman.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Inbox",
    URL:   "/",
    DropTypes: []application.DropType{
        application.DropFiles,
        application.DropText,
        application.DropURLs,
        application.DropImages,
    },
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    log.Println("files:", event.Context().DroppedFiles())
})

window.OnDrop(func(_ *application.Context, data application.DropData) {
    log.Println("text:", data.Text, "urls:", data.URLs, "images:", len(data.Images), "at", data.X, data.Y)
})
```

## Sistem

### Izin

`app.Permissions` melaporkan dan meminta izin privasi sistem (kamera, mikrofon, perekaman layar, aksesibilitas, lokasi, notifikasi, pemantauan masukan, dan akses disk penuh). `Status` tidak pernah menampilkan permintaan izin. `Request` menampilkan permintaan untuk jenis yang statusnya belum ditentukan dan memblokir hingga pengguna menjawab, jadi panggil dari goroutine. `OpenSystemSettings` membuka panel Privacy & Security yang sesuai.

```go
go func() {
    status := app.Permissions.Status(application.PermissionKindCamera)
    if status == application.PermissionStatusNotDetermined {
        status, _ = app.Permissions.Request(application.PermissionKindCamera)
    }
    if status == application.PermissionStatusDenied {
        _ = app.Permissions.OpenSystemSettings(application.PermissionKindCamera)
    }
}()
```

Akses disk penuh tidak dapat diminta dan menghasilkan `ErrPermissionNotRequestable`; arahkan pengguna ke panel pengaturan. Permintaan yang tidak pernah mendapat jawaban menghasilkan `ErrPermissionRequestTimeout`, yang pada macOS biasanya berarti kunci deskripsi penggunaan untuk jenis tersebut tidak ada dalam `Info.plist`.

Opsi jendela `Permissions` kini diterapkan pada macOS. Opsi ini menentukan cara permintaan `getUserMedia` dari halaman ditangani: `PermissionAllow` melewati permintaan izin milik WebView, `PermissionDeny` menolak tanpa bertanya, dan `PermissionDefault` menampilkan permintaan izin. Permintaan TCC tingkat sistem tetap muncul saat kamera atau mikrofon digunakan untuk pertama kali.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Meeting",
    URL:   "/",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionDefault,
    },
})
```

### Daya

`app.Power.PreventSleep` menjaga sistem tetap aktif, dan dengan `Display` juga menjaga layar tetap aktif, hingga fungsi pelepas yang dikembalikan dipanggil. Penahanan dihitung, sehingga beberapa bagian aplikasi dapat menahan secara bersamaan. Alasannya ditampilkan di Activity Monitor.

```go
release, err := app.Power.PreventSleep("Exporting video", application.PreventSleepOptions{Display: true})
if err != nil {
    log.Println(err)
}
defer release()

state := app.Power.State()
if state.LowPowerMode || state.ThermalState >= application.ThermalStateSerious {
    // trim background work
}
log.Println("battery", state.BatteryLevel, "charging", state.Charging, "on battery", state.OnBattery)
```

Perubahan diterima sebagai `events.Mac.ApplicationDidChangePowerState` (saat Mode Daya Rendah diaktifkan atau dinonaktifkan) dan `events.Mac.ApplicationDidChangeThermalState`.

### Siklus hidup

macOS dapat langsung menghentikan aplikasi yang sedang tidak aktif saat pengguna keluar dari sesi atau sistem dimatikan jika aplikasi mengaktifkan `NSSupportsSuddenTermination`, dan dapat menutup aplikasi yang sedang tidak aktif tanpa jendela jika aplikasi mengaktifkan `NSSupportsAutomaticTermination`. `app.Lifecycle.HoldTermination` menangguhkan keduanya selama bagian penting seperti penyimpanan berkas.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` mengaktifkan atau menonaktifkan penghentian mendadak saat runtime; `SuddenTerminationEnabled` melaporkan status saat ini, yang nilai awalnya berasal dari kunci `Info.plist`.

### Lingkungan

`app.Env` menyediakan tiga kueri tambahan tentang pengaturan pengguna.

```go
a11y := app.Env.Accessibility()
if a11y.ReduceMotion || a11y.ReduceTransparency {
    app.Event.Emit("theme:calm", true)
}

layout := app.Env.KeyboardLayout()
log.Println(layout.ID, layout.Name, layout.Languages)

locale := app.Env.Locale()
log.Println(locale.Identifier, locale.Language, locale.Region, locale.Preferred)
```

- `Accessibility` mencerminkan Kurangi Gerakan, Kurangi Transparansi, Tingkatkan Kontras, Bedakan Tanpa Warna, Balikkan Warna, VoiceOver, dan Switch Control.
- `KeyboardLayout` mengembalikan sumber masukan aktif beserta pengidentifikasi, nama yang dilokalkan, dan bahasanya.
- `Locale` mengembalikan lokal yang dipilih AppKit untuk aplikasi beserta seluruh daftar `Preferred` pengguna dalam urutannya. `Identifier` hanya mencerminkan bahasa yang dideklarasikan bundle dalam `CFBundleLocalizations`; gunakan `Preferred` untuk memilih bahasa sendiri.

### Event

Event aplikasi berikut adalah tambahan baru. Masing-masing dikirim melalui `app.Event.OnApplicationEvent`; kueri pengelola yang sesuai untuk mendapatkan nilai terbaru.

| Event | Dipicu saat | Baca dengan |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | Mode Daya Rendah diaktifkan atau dinonaktifkan | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | tekanan termal berubah | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | pengaturan tampilan aksesibilitas berubah | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | bentuk macOS dari perubahan yang sama | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | sumber masukan berubah | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | lokal berubah | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## Integrasi

### Menu Services

`app.ServicesProvider.Register` menambahkan entri ke submenu Services yang ditampilkan setiap aplikasi macOS untuk teks atau berkas yang dipilih. Penangan menerima papan klip sebagai `ServiceRequest` dan mengembalikan `ServiceResponse` untuk menulis balik; respons kosong membiarkan pilihan tetap seperti semula. Pastikan penangan berjalan cepat, karena AppKit menunggunya pada thread utama.

```go
err := app.ServicesProvider.Register(application.ServiceDefinition{
    Name:          "summarise",
    MenuTitle:     "Summarise with Notes",
    SendTypes:     []string{"public.utf8-plain-text"},
    ReturnTypes:   []string{"public.utf8-plain-text"},
    KeyEquivalent: "S",
    Handler: func(_ *application.Context, request application.ServiceRequest) (application.ServiceResponse, error) {
        return application.ServiceResponse{Text: "Summary: " + request.Text}, nil
    },
})
if err != nil {
    log.Println(err)
}

// Paste this inside the top-level <dict> of Info.plist.
log.Println(app.ServicesProvider.InfoPlistXML())
```

`Name` adalah pesan yang dikirim AppKit dan harus berupa pengidentifikasi biasa. `SendTypes` dan `ReturnTypes` adalah tipe papan klip; sebuah layanan memerlukan setidaknya salah satunya. Pendaftaran saja tidak membuat layanan terlihat: `Info.plist` milik bundle harus mendeklarasikannya di bawah `NSServices`. `InfoPlistXML` mengembalikan blok tersebut yang siap ditempel, sedangkan `InfoPlistEntries` mengembalikan data yang sama sebagai map untuk serialisator plist. `NSPortName` dalam entri tersebut adalah `Name` aplikasi, yang harus cocok dengan `CFBundleName`.

Proyek yang dibuat dengan Wails CLI dapat mendeklarasikan layanan yang sama satu kali dalam `build/config.yml` dan membiarkan proses pengemasan menghasilkan blok `NSServices`:

```yaml
services:
  - name: SummariseText
    menuTitle: Summarise with My Product
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S
```

Setiap entri harus cocok dengan `ServiceDefinition` yang didaftarkan di Go dengan `Name` yang sama. Jalankan `pbs -update` setelah memasang build baru agar menu Services mengambil perubahan tanpa perlu keluar dari sesi.

### Handoff dan aktivitas pengguna

`app.Activity.Publish` menjadikan `NSUserActivity` aktif sehingga pengguna dapat melanjutkannya di perangkat lain, menemukannya di Spotlight, atau menerima saran dari Siri. `PublishedActivity` yang dikembalikan dapat diperbarui saat status berubah dan dibatalkan ketika dokumen ditutup. Menerbitkan aktivitas baru menggantikan aktivitas sebelumnya.

```go
activity, err := app.Activity.Publish(application.UserActivity{
    Type:               "com.example.notes.editing",
    Title:              "Editing Quarterly report",
    UserInfo:           map[string]any{"note": "quarterly-report"},
    WebpageURL:         "https://example.com/notes/quarterly-report",
    EligibleForHandoff: true,
    EligibleForSearch:  true,
    Keywords:           []string{"report", "quarterly"},
})
if err != nil {
    log.Println(err)
    return
}
if err := activity.Update(map[string]any{"note": "quarterly-report", "cursor": 120}); err != nil {
    log.Println(err)
}
// When the document closes:
activity.Invalidate()
```

Aktivitas masuk diterima melalui `OnContinue`. Tautan universal membawa tipe `UserActivityTypeBrowsingWeb` dengan halaman di `WebpageURL`, dan juga dikirim sebagai `events.Common.ApplicationLaunchedWithUrl` agar aplikasi dapat menggunakan satu alur kode URL. `OnWillContinue`, `OnFailed`, dan `OnUpdated` menangani bagian lain dari delegate.

```go
app.Activity.OnContinue(func(_ *application.Context, incoming application.UserActivity) bool {
    if incoming.Type == application.UserActivityTypeBrowsingWeb {
        log.Println("universal link", incoming.WebpageURL)
        return true
    }
    note, _ := incoming.UserInfo["note"].(string)
    log.Println("continue editing", note)
    return true
})
```

Setiap tipe aktivitas harus dicantumkan di bawah `NSUserActivityTypes` dalam `Info.plist`. Tautan universal juga memerlukan entitlement `com.apple.developer.associated-domains` dengan entri `applinks:example.com` dan berkas `apple-app-site-association` yang sesuai di domain tersebut.

### Apple Events

`app.AppleEvents.Handle` mendaftarkan penangan untuk kelas dan ID event, sehingga AppleScript, Shortcuts, dan aplikasi lain dapat mengendalikan aplikasi. Kodenya berupa string empat karakter. Parameter langsung didekode menjadi nilai Go (`string`, `[]string` berisi jalur berkas, `int64`, `float64`, `bool`, `[]any`, atau `AppleEventRawData`), dan `Result` pada balasan menerima jenis yang sama. Penangan berjalan pada goroutine masing-masing selama event ditangguhkan.

```go
err := app.AppleEvents.Handle("WAIL", "note", func(_ *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
    text, _ := event.DirectObject.(string)
    return application.AppleEventReply{Result: "noted: " + text}, nil
})
if err != nil {
    log.Println(err)
}

go func() {
    result, err := app.AppleEvents.Send("com.apple.finder", "misc", "actv", nil)
    log.Println(result, err)
}()

sdef := app.AppleEvents.ScriptingDefinition()
if err := os.WriteFile("Notes.sdef", []byte(sdef), 0o644); err != nil {
    log.Println(err)
}
```

Skrip dapat langsung memanggil penangan dengan sintaks event mentah:

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` menghasilkan `.sdef` minimal yang memberi setiap penangan nama perintah. Sertakan dalam `Contents/Resources` dan rujuk dari `Info.plist` dengan `NSAppleScriptEnabled` dan `OSAScriptingDefinition`; Script Editor kemudian menampilkannya di File > Open Dictionary. Wails sudah menangani event Get URL untuk skema URL khusus; penangan untuk `"GURL"`/`"GURL"` akan dirangkaikan dengannya, sedangkan penangan untuk `"aevt"`/`"odoc"` menggantikan pengiriman Open Documents bawaan. `Send` menargetkan aplikasi yang sedang berjalan berdasarkan pengidentifikasi bundle, memerlukan `NSAppleEventsUsageDescription` dalam aplikasi berbentuk bundle, dan memblokir goroutine pemanggil hingga balasan diterima.

### Quick Look

`app.QuickLook.Preview` membuka panel Quick Look bersama untuk satu atau beberapa berkas; jika ada beberapa jalur, panel menampilkan panah untuk berpindah di antaranya. `Thumbnail` merender berkas melalui penyedia gambar mini sistem dan mengembalikan PNG, sehingga dokumen, gambar, PDF, dan film semuanya didukung.

```go
if err := app.QuickLook.Preview([]string{"/Users/me/Documents/Report.pdf"}); err != nil {
    log.Println(err)
}

go func() {
    png, err := app.QuickLook.Thumbnail("/Users/me/Documents/Report.pdf", application.ThumbnailOptions{
        Width: 256,
        Scale: 2,
    })
    if err != nil {
        log.Println(err)
        return
    }
    if err := os.WriteFile("thumbnail.png", png, 0o644); err != nil {
        log.Println(err)
    }
}()
```

Jalur harus absolut dan berkasnya harus ada. `ClosePreview` dan `IsPreviewOpen` mengelola panel. `ThumbnailOptions.IconMode` menggambar bingkai dokumen bergaya Finder, dan `Scale: 2` menghasilkan gambar Retina. `Thumbnail` memblokir goroutine pemanggil, jadi panggil dari goroutine atau metode yang diikat.

### Utilitas ruang kerja

`app.Browser` menyediakan tiga utilitas tambahan untuk `NSWorkspace`. `OpenWith` membuka berkas dengan aplikasi tertentu yang ditentukan berdasarkan pengidentifikasi bundle atau jalur bundle. `ApplicationsForFile` mencantumkan aplikasi terpasang yang dapat membuka berkas, dengan penangan bawaan di urutan pertama. `ActivateApplication` membawa aplikasi yang sedang berjalan ke depan.

```go
apps := app.Browser.ApplicationsForFile("/Users/me/Documents/Report.md")
for _, info := range apps {
    log.Println(info.Name, info.BundleID, info.Path)
}

if err := app.Browser.OpenWith("/Users/me/Documents/Report.md", "com.apple.TextEdit"); err != nil {
    log.Println(err)
}

if err := app.Browser.ActivateApplication("com.apple.TextEdit"); err != nil {
    log.Println(err)
}
```

### Spotlight

`app.Spotlight.Index` menambahkan konten aplikasi ke indeks pencarian sistem melalui Core Spotlight. Setiap `SearchableItem` memiliki `ID`, `Title`, dan secara opsional `Domain` untuk penghapusan massal, `Description`, `Keywords`, `ContentType`, gambar mini PNG, `URL` tautan langsung, dan waktu kedaluwarsa. `OnOpen` dipanggil saat pengguna memilih salah satu item di Spotlight, atau memilih "Cari di Aplikasi" dengan kueri.

```go
err := app.Spotlight.Index([]application.SearchableItem{{
    ID:          "note:quarterly-report",
    Domain:      "notes",
    Title:       "Quarterly report",
    Description: "Draft for the board meeting",
    Keywords:    []string{"finance", "q3"},
    ContentType: "public.plain-text",
    URL:         "notes://open/quarterly-report",
}})
if err != nil {
    log.Println(err)
}

stop := app.Spotlight.OnOpen(func(_ *application.Context, id string, query string) {
    if query != "" {
        log.Println("search in app:", query)
        return
    }
    log.Println("open item", id)
})
defer stop()

// Later:
_ = app.Spotlight.Delete([]string{"note:quarterly-report"})
_ = app.Spotlight.DeleteDomain("notes")
```

`IsAvailable` melaporkan apakah indeks menerima item. Pengindeksan memerlukan aplikasi berbentuk bundle: item yang diindeks oleh biner `go run` tanpa bundle tidak pernah muncul di Spotlight. `DeleteAll` menghapus semua yang telah diindeks aplikasi.

## Berikutnya

Contoh `mac-windows-extra` yang mencakup lembar, popover, opsi presentasi, dan pemulihan status sedang ditambahkan dan akan ditautkan di sini setelah tersedia.

## Persyaratan versi

Semua yang ada di halaman ini hanya untuk macOS, dan API Go tetap sama di semua platform. Wails menargetkan macOS 10.13 dan yang lebih baru; fitur yang memerlukan rilis lebih baru mengalami penurunan kemampuan seperti dijelaskan.

| Fitur | macOS minimum | Perilaku pada rilis sebelumnya |
|---------|---------------|-------------------------------|
| Status izin kamera dan mikrofon | 10.14 | dilaporkan sebagai diizinkan (rilis sebelumnya tidak membatasi perangkat penangkap) |
| `Speech.Speak` dan `Voices` | 10.14 | `ErrSpeechNotSupported` |
| Izin perekaman layar dan pemantauan masukan | 10.15 | dilaporkan sebagai diizinkan |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | diabaikan dengan entri log debug |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| `SetSymbol` pada item menu dan item status | 11 | tidak ada gambar yang ditampilkan |
| `AddContentType`, `SetFormats` berdasarkan UTI | 11 | pengidentifikasi yang sama diterapkan melalui API tipe berkas yang diizinkan versi lama |
| Ikon janji berkas dalam `StartDrag` | 11 | ikon dokumen generik |
| `PowerState.LowPowerMode` dan event-nya | 12 | selalu false; event tidak pernah dipicu |
| `InteractionState` dan `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| Panel Notifikasi dalam `OpenSystemSettings` | 13 | panel preferensi Notifikasi versi lama terbuka |
| Lencana menu, header bagian, palet | 14 | lencana tidak ditampilkan; header menjadi item nonaktif; palet disembunyikan |
| `MacToolbarItem.ShowPopover` pada item tanpa tampilan khusus | 14 | `ErrMacPopoverAnchorUnavailable` |

Fitur lainnya tidak memiliki persyaratan di luar minimum Wails.

## Kunci Info.plist

Beberapa fitur bergantung pada kunci dalam `Info.plist` aplikasi. Deskripsi penggunaan ditampilkan kepada pengguna dalam permintaan izin; tanpa kunci tersebut, macOS tidak akan menampilkan permintaan izin dan permintaan akan mengalami waktu habis.

| Kunci | Diperlukan oleh |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, akses kamera dari halaman |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, akses mikrofon dari halaman, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | status awal `Lifecycle.SuddenTerminationEnabled`; `HoldTermination` menangguhkannya |
| `NSSupportsAutomaticTermination` | memungkinkan macOS menutup aplikasi yang sedang tidak aktif; `HoldTermination` menangguhkannya |
| `CFBundleLocalizations` | bahasa yang dapat dilaporkan oleh `Env.Locale().Identifier` |
| `NSServices` | satu entri untuk setiap layanan yang didaftarkan dengan `app.ServicesProvider`; hasilkan bloknya dengan `InfoPlistXML` |
| `NSUserActivityTypes` | setiap `UserActivity.Type` yang diterbitkan atau dilanjutkan melalui `app.Activity` |
| `NSAppleScriptEnabled` dan `OSAScriptingDefinition` | menandai aplikasi dapat dikendalikan dengan skrip dan menamai `.sdef` yang ditulis dari `app.AppleEvents.ScriptingDefinition` |
| `NSAppleEventsUsageDescription` | `app.AppleEvents.Send` ke aplikasi lain |

Izin notifikasi dan pengenalan ucapan juga mengharuskan aplikasi berjalan sebagai bundle dengan pengidentifikasi bundle; biner `go run` tanpa bundle melaporkan `PermissionStatusUnsupported` untuk notifikasi.

<a id="platform-notes"></a>

## Catatan platform

Setiap API di halaman ini dapat dikompilasi di Windows dan Linux. Di luar macOS:

- Fitur tambahan jendela: `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName`, dan `SetWindowButtonsOffset` tidak melakukan apa pun. `WindowCascade` berperilaku seperti `WindowCentered`. `RequestAttention` mengembalikan handle yang `Cancel`-nya tidak melakukan apa pun. `PrintWithOptions` memanggil `Print`. `ExportPDF` dan `Snapshot` mengembalikan `ErrMacOnly`.
- Menu: simbol, lencana, status campuran, item alternatif, dan indentasi disimpan tetapi tidak digambar. Header bagian menjadi item nonaktif. Palet disembunyikan. Submenu Buka Terkini diisi dari daftar terkini Go saat menu dibuat. Menu Dock tidak pernah ditampilkan.
- Dialog: `SetSuppression`, `SetHelp`, `SetNameFieldLabel`, dan `SetTags` diabaikan. `AddContentType` dan `SetFormats` menjadi filter ekstensi jika padanannya diketahui. `Prompt`, `PickColor`, dan `PickFont` mengembalikan `ErrDialogNotSupported`.
- Item status: `SetSymbol`, `SetRemovable`, dan `OnVisibilityChange` tidak berpengaruh. `IsVisible` mencerminkan pemanggilan `Show` atau `Hide` terakhir.
- Umpan balik: umpan balik haptik tersedia di iOS dan Android; pada Windows dan Linux tidak melakukan apa pun. `Sound.Beep` dan `Sound.Play` berfungsi di Windows dengan berkas WAV dan alias registri; di tempat lain `Play` mengembalikan `ErrSoundNotSupported`. Ucapan mengembalikan `ErrSpeechNotSupported` dan `ErrSpeechRecognitionNotSupported`.
- Papan klip: metode kaya format mengembalikan `ErrClipboardNotSupported`, `Types` kosong, `ChangeCount` bernilai 0, dan `OnChange` tidak pernah dipicu.
- Penyeretan: `StartDrag` mengembalikan `ErrDragOutUnsupported`. Tipe item yang dilepas selain berkas tidak dikirim; pelepasan berkas tetap berfungsi melalui `WindowFilesDropped`.
- Sistem: `Permissions.Status` melaporkan `PermissionStatusUnsupported` dan `Request` mengembalikan `ErrPermissionsUnsupported`. `PreventSleep` mengembalikan `ErrPreventSleepUnsupported` beserta fungsi pelepas yang tidak melakukan apa pun. `HoldTermination` mengembalikan fungsi pelepas yang tidak melakukan apa pun. Semua nilai `Accessibility` bernilai false, `KeyboardLayout` bernilai nol, dan `Locale` diturunkan dari `LC_ALL`, `LC_MESSAGES`, serta `LANG`. Opsi jendela `Permissions` bersifat lintas platform.
- Integrasi: `ServicesProvider.Register` mengembalikan `ErrServicesUnsupported` dan `Activity.Publish` mengembalikan `ErrActivityUnsupported`; `InfoPlistXML`, `InfoPlistEntries`, dan penangan aktivitas tetap berfungsi. `Handle` dan `Send` pada `AppleEvents` mengembalikan `ErrAppleEventsNotSupported`; `ScriptingDefinition` dihasilkan di semua platform. `QuickLook.Preview` dan `Thumbnail` mengembalikan `ErrQuickLookNotSupported`. Metode pengindeksan `Spotlight` mengembalikan `ErrSpotlightNotSupported` dan `OnOpen` tidak pernah dipicu. `Browser.OpenWith` menjalankan program yang disebutkan dengan jalur sebagai argumennya, `ApplicationsForFile` kosong, dan `ActivateApplication` mengembalikan `ErrApplicationNotRunning`.
- Opsi presentasi: `SetPresentationOptions` mengembalikan `ErrMacOnly` dan `PresentationOptions` bernilai `MacPresentationDefault`.
- Lembar: `PresentSheet`, `PresentCriticalSheet`, dan `PresentNativeSheet` mengembalikan `ErrMacSheetUnsupported`; `EndSheet` tidak melakukan apa pun dan metode kueri melaporkan tidak ada lembar.
- Popover: `NewMacPopover` berfungsi, metode penampilannya mengembalikan `ErrMacPopoverUnsupported`, dan `IsShown` bernilai false.
- Pemulihan status: `SetRestorationID` dan `SetRestorationData` tidak melakukan apa pun, `OnRestore` tidak pernah dipanggil, serta `InteractionState` dan `RestoreInteractionState` mengembalikan `ErrMacOnly`.

## Contoh

Setiap contoh adalah aplikasi lengkap yang dapat dijalankan:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go` menghubungkan titik perubahan, subjudul, ekspor PDF, opsi cetak, dan jendela berjenjang ke editor catatan.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): simbol, lencana, header bagian, status campuran, item alternatif, palet, Buka Terkini, menu Dock dinamis, dan progres Dock.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): kotak centang agar pesan tidak ditampilkan lagi dan tombol bantuan, prompt teks, tipe konten, menu tarik turun Format, tag Finder, serta panel warna dan font.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): item status SF Symbol yang dapat dihapus, umpan balik haptik, suara sistem, sintesis ucapan, dan pengenalan ucapan.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): papan klip kaya format dengan pelacakan perubahan, penyeretan keluar dengan janji berkas, serta pelepasan teks, URL, dan gambar.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): izin, pencegahan tidur, penahanan penghentian, status daya, aksesibilitas, tata letak papan ketik, dan lokal dengan pembaruan langsung.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): entri menu Services, aktivitas Handoff dengan penangan kelanjutan, dan utilitas ruang kerja pada `app.Browser`.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): pengindeksan Spotlight dengan `OnOpen`, pratinjau dan gambar mini Quick Look, serta Apple Event khusus dengan definisi skripnya.
