---
title: "Antarmuka Jendela Native macOS"
description: "Bangun toolbar, sidebar, daftar konten, panel pemeriksa, aksesori, dan tab jendela AppKit native di sekitar jendela Wails Anda"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Platform yang relevan: macOS

Wails v3 dapat membungkus WebView Anda dengan antarmuka jendela AppKit asli. Toolbar, sidebar, daftar konten, panel pemeriksa, dan bilah titlebar adalah kontrol native yang dibuat dari Go. Tidak satu pun menggunakan HTML, sehingga semuanya mendapat material, penanganan papan ketik, animasi, dan persistensi AppKit tanpa pekerjaan tambahan, sementara frontend Anda tetap berfokus pada konten.

Setiap API di halaman ini berada dalam paket `application` dan diawali dengan `Mac`. Kode yang sama dapat dikompilasi di Windows dan Linux: konstruktor dan setter berfungsi di semua platform, pemanggilan untuk memasang komponen tidak melakukan apa pun atau mengembalikan kesalahan, dan jendela tetap menggunakan satu WebView seperti biasa.

## Anatomi jendela

Jendela dengan seluruh komponennya memiliki bagian berikut, dari tepi awal hingga tepi akhir:

| Bagian | Tipe | Kelas AppKit |
|------|------|--------------|
| Toolbar | `MacToolbar` | `NSToolbar` |
| Sidebar | `MacSidebar` | daftar sumber `NSOutlineView` dalam item tampilan terbagi sidebar |
| Daftar konten | `MacContentList` | `NSTableView` dalam item tampilan terbagi daftar konten |
| Konten utama | WebView Anda, atau `MacTextEditor` | `WKWebView` atau `NSTextView` |
| Panel pemeriksa | `MacInspector` | kontrol properti native dalam item tampilan terbagi panel pemeriksa |
| Aksesori | `MacAccessory` | `NSTitlebarAccessoryViewController` atau `NSSplitViewItemAccessoryViewController` |

Panel-panel diatur oleh `MacSplitView`, yaitu sebuah `NSSplitViewController`. Buat komponennya terlebih dahulu, tambahkan ke tampilan terbagi dari awal ke akhir, pasang tampilan terbagi ke jendela, lalu pasang toolbar. Anda dapat melakukannya dengan `SetSplitView` dan `SetToolbar`, atau dalam satu langkah melalui opsi jendela `Mac.SplitView` dan `Mac.Toolbar`.

Pengaturan jendela yang umum memasangkan antarmuka ini dengan opsi jendela berikut agar konten dapat digulir di bawah toolbar terpadu:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        ContentLayout: application.MacContentLayoutEdgeToEdge,
        TitleBar: application.MacTitleBar{
            FullSizeContent:      true,
            HideToolbarSeparator: true,
            ToolbarStyle:         application.MacToolbarStyleUnified,
        },
    },
})
```

## Toolbar

`NewMacToolbar` membuat `NSToolbar`. Tambahkan item dengan metode `Add`, rangkai setter dan callback pada handle yang dikembalikan, lalu pasang toolbar dengan `SetToolbar`. Pengidentifikasi dibuat secara otomatis.

```go
toolbar := application.NewMacToolbar().
    SetDisplayMode(application.MacToolbarDisplayModeIconOnly)

// Standard AppKit items. They have no handle because AppKit owns them.
toolbar.AddSidebarToggle()
toolbar.AddSidebarTrackingSeparator()

toolbar.AddButton("New").
    SetSymbol("square.and.pencil").
    SetTooltip("Create a note").
    SetBordered(true).
    OnClick(func(*application.Context) {
        // create a note
    })

toolbar.AddSearch("Search").
    SetSearchPlaceholder("Search notes").
    SetSearchIncremental(true).
    OnSearch(func(_ *application.Context, query string) {
        filterNotes(query)
    })

toolbar.AddFlexibleSpace()

mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(0)
})
mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(1)
})

actions := application.NewMenu()
actions.Add("Export...").OnClick(func(*application.Context) {})
toolbar.AddMenu("Actions", actions).SetSymbol("ellipsis.circle")

window.SetToolbar(toolbar)
```

Jenis itemnya adalah:

- `AddButton` menambahkan tombol tekan. Setiap tombol memerlukan `OnClick` sebelum toolbar dipasang; jika tidak, `SetToolbar` melaporkan kesalahan dan membiarkan toolbar sebelumnya tetap terpasang.
- `AddSearch` menambahkan `NSSearchToolbarItem`. `OnSearch` wajib ada. Gunakan `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` untuk menu pencarian terkini yang persisten, dan `SetSearchMenu` untuk menu khusus di balik ikon kaca pembesar.
- `AddShare` menambahkan item berbagi sistem (lihat di bawah).
- `AddGroup` menambahkan `NSToolbarItemGroup` bersegmen. Tambahkan anggota dengan `AddButton` milik grup dan pilih `ToolbarGroupSelectOne`, `ToolbarGroupMomentary`, atau `ToolbarGroupSelectAny`.
- `AddMenu` menambahkan `NSMenuToolbarItem` tarik turun yang menggunakan `Menu` biasa. `SetShowsIndicator(false)` menyembunyikan tanda panah.
- `AddSpace` dan `AddFlexibleSpace` menambahkan spasi standar.
- `AddSidebarToggle` dan `AddSidebarTrackingSeparator` menambahkan item sidebar AppKit. Pemisah menjaga semua item sebelumnya tetap sejajar di atas pembatas sidebar, sehingga jendela harus memiliki tampilan terbagi dengan panel sidebar.
- `AddInspectorToggle` dan `AddInspectorTrackingSeparator` melakukan hal yang sama untuk panel pemeriksa.

`SetDisplayMode` memilih antara `MacToolbarDisplayModeIconAndLabel` (bawaan), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly`, dan `MacToolbarDisplayModeDefault`.

### Pembaruan langsung

Setiap handle juga dapat digunakan untuk pembaruan langsung. Setter yang dipanggil setelah pemasangan memperbarui item native pada thread aplikasi.

```go
save := toolbar.AddButton("Save").SetSymbol("checkmark.circle")
save.OnClick(func(*application.Context) {
    save.SetBadgeCount(0).SetProminent(false)
})

// Later, when the document changes:
save.SetBadgeCount(1).SetProminent(true)
save.SetVisibilityPriority(application.MacToolbarVisibilityPriorityHigh)

toolbar.Move(save, 0)
toolbar.Remove(save)
```

Setter lain yang berguna adalah `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor`, dan `SetNavigational`, yang mempertahankan item di tepi awal seperti tombol kembali dan maju di Safari. `SetVisibilityPriority` menentukan item mana yang lebih dahulu dipindahkan ke menu luapan ketika jendela menyempit.

### Penyesuaian oleh pengguna

`SetCustomizable` mengaktifkan lembar standar "Sesuaikan Toolbar..." dan menyimpan tata letak pengguna dengan kunci yang Anda berikan. Berikan `SetPersistenceKey` yang stabil pada setiap item agar tata letak tersimpan tetap tersedia setelah aplikasi dijalankan ulang, dan gunakan `SetInDefaultSet(false)` untuk item yang hanya boleh muncul setelah ditambahkan pengguna. Panggil `SetCustomizable` sebelum toolbar dipasang.

```go
toolbar := application.NewMacToolbar().SetCustomizable("myapp.main-toolbar")

newNote := toolbar.AddButton("New").
    SetPersistenceKey("new").
    OnClick(func(*application.Context) {})

// Offered in the palette but hidden until the user adds it.
toolbar.AddButton("Archive").
    SetPersistenceKey("archive").
    SetInDefaultSet(false).
    OnClick(func(*application.Context) {})

toolbar.SetCenteredItems(newNote)

// From a menu item, for example:
toolbar.RunCustomizationPalette()
```

### Berbagi

`AddShare` mengembalikan `MacToolbarShareItem`. Item ini tetap nonaktif sampai `MacShareProvider` menyediakan setidaknya satu representasi. Wails meminta byte dari penyedia hanya saat layanan berbagi memintanya, sehingga ekspor berukuran besar dihasilkan sesuai kebutuhan. `MacShareProviderFunc` mengadaptasi dua fungsi menjadi penyedia; aplikasi yang memiliki state dapat mengimplementasikan antarmukanya secara langsung.

```go
share := toolbar.AddShare("Share")
share.SetProvider(application.MacShareProviderFunc{
    Available: []application.MacShareRepresentation{
        {ContentType: application.MacShareTypePDF},
        {ContentType: application.MacShareTypePlainText},
    },
    Load: func(request application.MacShareRequest) ([]byte, error) {
        if request.ContentType == application.MacShareTypePDF {
            return renderPDF()
        }
        return []byte(currentText()), nil
    },
}).SetSubject("Saturday, slowly").SetSuggestedName("Note")

share.OnShared(func(_ *application.Context, service string) {
    // service is the localised name of the sharing service
})
share.OnShareError(func(_ *application.Context, service string, err error) {
    window.Error("share via %s failed: %s", service, err)
})
```

Tipe konten yang umum adalah `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG`, dan `MacShareTypeJPEG`. String UTI lainnya juga diterima.

### Memasang dan melepas

Satu toolbar hanya dapat dimiliki satu jendela pada satu waktu. `WebviewWindow.SetToolbar` dan `NativeWindow.SetToolbar` menerima toolbar sebelum atau sesudah jendela ada; meneruskan `nil` akan melepasnya dan membuatnya tersedia untuk digunakan di tempat lain. `SetToolbar` pada `WebviewWindow` melaporkan masalah validasi melalui `Window.Error`, sedangkan versi `NativeWindow` mengembalikan kesalahannya. Opsi jendela `Mac.Toolbar` memasang toolbar saat pembuatan; opsi ini diterapkan setelah `Mac.SplitView` agar pemisah pelacak dapat menemukan sidebar yang menjadi acuan kesejajarannya.

## Tampilan terbagi

`MacSplitView` mengatur panel. Tambahkan panel dari awal ke akhir. `AddSidebar`, `AddContentList`, dan `AddInspector` menerima model native yang ditampungnya dan mengembalikan `MacSplitPane` untuk mengatur ukuran dan pelipatan. `AddPrimaryContent` menempatkan WebView yang sudah ada pada jendela dan mengembalikan `MacSplitWebviewPane`, yang menyediakan `SetContentLayout`.

```go
split := application.NewMacSplitView().SetAutosaveName("myapp.main-window")

sidebarPane := split.AddSidebar(sidebar).
    SetMinimumThickness(200).
    SetMaximumThickness(320).
    SetCollapsible(true)

split.AddContentList(list).
    SetMinimumThickness(240).
    SetCollapsible(true)

split.AddPrimaryContent().
    SetContentLayout(application.MacContentLayoutEdgeToEdge)

split.AddInspector(inspector).
    SetPreferredThicknessFraction(0.25).
    SetHoldingPriority(300).
    SetCollapsible(true).
    SetCanCollapseFromWindowResize(false).
    SetCollapsed(true)

sidebarPane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
    // fired for toggles, gestures, menu items and SetCollapsed alike
})

window.SetSplitView(split)

// At runtime:
sidebarPane.Toggle()
```

`SetSplitView` berfungsi sebelum atau sesudah jendela native ada. Jika dipanggil sebelumnya, tata letak dimasukkan ke antrean dan dipasang saat jendela dibuat. Jika dipanggil sesudahnya, misalnya dari callback menu atau tray pada aplikasi yang sedang berjalan, tata letak langsung dipasang: WebView yang sudah ada pada jendela menjadi panel utama, toolbar saat ini dipasang kembali agar pemisah pelacaknya sejajar, dan aksesori yang mengantre dipasang.

Jendela yang dibuat setelah `app.Run` dapat dikonfigurasi dalam satu pemanggilan menggunakan opsi `Mac.SplitView` dan `Mac.Toolbar`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Notes",
    URL:   "/",
    Mac: application.MacWindow{
        SplitView: split,
        Toolbar:   toolbar,
    },
})
```

Aturan tata letak adalah:

- Tata letak memerlukan setidaknya dua panel dan tepat satu panel utama (`AddPrimaryContent` untuk `WebviewWindow`, `AddTextEditor` untuk `NativeWindow`).
- Maksimal satu daftar konten, ditempatkan setelah sidebar dan sebelum panel utama.
- Struktur panel tidak dapat diubah setelah tampilan terbagi dipasang. Pengaturan panel, status terlipat, dan isi sidebar, daftar, serta panel pemeriksa tetap dapat diubah kapan saja.
- Tata letak yang sudah dipasang tidak dapat diganti. Pemanggilan `SetSplitView` kedua pada jendela yang sama melaporkan `ErrMacSplitViewAlreadyInstalled` (melalui `Window.Error` pada `WebviewWindow`, sebagai nilai kembalian pada `NativeWindow`) dan membiarkan jendela tetap seperti semula. Meneruskan `nil` sebelum pemasangan akan menghapus tata letak yang tertunda.
- Sidebar, daftar, panel pemeriksa, atau tampilan terbagi hanya dapat dimiliki satu jendela pada satu waktu.

Setter panel adalah `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed`, dan `OnCollapsedChange`. `SetAutosaveName` menyimpan posisi pembatas agar bertahan saat aplikasi dijalankan ulang.

`SetContentLayout` pada panel utama memilih antara `MacContentLayoutBelowToolbar` dan `MacContentLayoutEdgeToEdge`. `MacContentLayoutAutomatic` mewarisi `MacWindow.ContentLayout`, yang pada gilirannya mengikuti `TitleBar.FullSizeContent`. Tata letak dari tepi ke tepi memungkinkan AppKit menerapkan efek tepi gulir di bawah toolbar pada macOS 26 dan yang lebih baru.

## Sidebar

`MacSidebar` adalah daftar sumber native. Daftar ini memuat baris akar, bagian, dan baris bersarang hingga kedalaman berapa pun. Simpan handle yang dikembalikan untuk memperbarui baris nanti.

```go
sidebar := application.NewMacSidebar()

// A root row above the sections.
sidebar.AddItem("All Notes").
    SetSymbol("tray.full").
    SetBadge(12).
    OnClick(func(*application.Context) {})

notes := sidebar.AddSection("Notes")
draft := notes.AddItem("Saturday, slowly").
    SetSymbol("doc.text").
    SetAccessorySymbol("pin.fill").
    SetTooltip("Field notes").
    SetEditable(true).
    OnRename(func(_ *application.Context, label string) {
        // the row label is already updated
    })

// Nested rows with a tinted symbol.
tags := sidebar.AddSection("Tags")
tint := application.NewRGB(0, 122, 255)
work := tags.AddItem("Work").SetSymbol("briefcase").SetTintColor(&tint)
work.AddItem("Meetings").SetSymbol("tag")
work.SetExpanded(true)

sidebar.SetSelectedItem(draft)
```

Setter baris adalah `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable`, dan `SetExpanded`. `OnClick` dipicu ketika AppKit memilih baris, `OnExpandedChange` ketika pengguna membuka atau menutup baris bersarangnya, dan `OnRename` setelah penggantian nama langsung dikonfirmasi.

### Pemilihan

Secara bawaan, hanya satu item dapat dipilih. `SetSelectedItem` memilih baris tanpa memicu `OnClick`-nya. Jika pemilihan beberapa item diaktifkan, `OnClick` tetap dipicu untuk baris yang diklik dan `OnSelectionChange` melaporkan seluruh kumpulan yang dipilih.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Menu konteks

Klik kanan mencari menu dengan urutan berikut: `SetContextMenu` milik baris, lalu callback `OnContextMenu` milik sidebar, lalu `SetContextMenu` cadangan milik sidebar. Callback berjalan pada thread aplikasi saat AppKit menunggu, jadi pastikan eksekusinya singkat.

```go
sidebar.OnContextMenu(func(_ *application.Context, item *application.MacSidebarItem) *application.Menu {
    if item == nil {
        return nil // fall back to the menu set with SetContextMenu
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { item.Remove() })
    return menu
})

empty := application.NewMenu()
empty.Add("New Note").OnClick(func(*application.Context) {})
sidebar.SetContextMenu(empty)
```

### Mengurutkan ulang dengan seret

`SetReorderable` memungkinkan pengguna menyeret baris di dalam satu bagian, antarbagiannya, serta ke atau dari akar. Model Go diperbarui sebelum `OnMove` dipicu.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Penghapusan

Menghapus baris juga menghapus baris bersarang di dalamnya. Setelah itu, handle tidak lagi berfungsi.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Daftar konten

`MacContentList` adalah kolom tengah pada Finder, Mail, dan penjelajah dokumen. Tanpa kolom, daftar ini menampilkan baris kaya informasi: judul, subjudul, simbol di awal, detail di akhir, dan lencana jumlah.

```go
list := application.NewMacContentList().
    SetStyle(application.MacContentListStyleInset).
    SetEmptyText("No notes match")

row := list.AddRow("Saturday, slowly").
    SetSubtitle("A slow day is still a day well spent.").
    SetDetail("Yesterday").
    SetSymbol("doc.text").
    SetBadge(2)

list.OnSelectionChange(func(_ *application.Context, rows []*application.MacContentListRow) {})
list.OnActivate(func(_ *application.Context, row *application.MacContentListRow) {
    // double-click or Return
})
list.SetSelectedRow(row)
```

Gayanya adalah `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain`, dan `MacContentListStyleFullWidth`. `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible`, dan `SetAllowsMultipleSelection` melengkapi opsi tampilan. Baris mendukung `SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` pada posisi tertentu, `Remove`, dan `RemoveAll`.

### Kolom

`SetColumns` mengalihkan daftar ke mode tabel. Baris kemudian menampilkan nilai `SetCells`, satu per kolom. Kolom yang ditandai `Sortable` menampilkan indikator pengurutan; `SetSortable` mengaktifkan pengurutan melalui klik pada header. Tanpa callback `OnSort`, daftar mengurutkan dirinya berdasarkan teks kolom dengan `SortBy`.

```go
table := application.NewMacContentList().
    SetColumns(
        application.MacContentListColumn{Title: "Name", Width: 220, Sortable: true},
        application.MacContentListColumn{Title: "Size", Width: 80, Alignment: application.MacContentListAlignTrailing},
        application.MacContentListColumn{Title: "Modified", Sortable: true},
    ).
    SetSortable(true).
    SetAlternatingRowBackgrounds(true)

table.AddRow("notes.txt").SetCells("notes.txt", "4 KB", "Today")
table.AddRow("ideas.txt").SetCells("ideas.txt", "1 KB", "Yesterday")

// Without OnSort the list sorts itself by the column text.
table.OnSort(func(_ *application.Context, column int, ascending bool) {
    table.SortRows(func(a, b *application.MacContentListRow) bool {
        return (a.Cells()[column] < b.Cells()[column]) == ascending
    })
})
```

### Menu konteks

Menu konteks ditentukan dengan urutan yang sama seperti sidebar: `SetContextMenu` milik baris, lalu `OnContextMenu`, lalu menu cadangan milik daftar.

```go
list.OnContextMenu(func(_ *application.Context, row *application.MacContentListRow) *application.Menu {
    if row == nil {
        return nil
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { row.Remove() })
    return menu
})
```

## Panel pemeriksa

`MacInspector` adalah panel properti di sisi akhir yang dibangun dari kontrol native dan dikelompokkan ke dalam bagian. Setiap metode `Add` mengembalikan handle `MacInspectorControl` dengan setter dan callback khusus menurut jenisnya.

```go
inspector := application.NewMacInspector()

document := inspector.AddSection("Document")
document.AddTextField("Title", "Saturday, slowly").
    OnTextChange(func(_ *application.Context, value string) {})
document.AddPopup("Category", []string{"Personal", "Work"}, 0).
    OnSelectionChange(func(_ *application.Context, index int, value string) {})
document.AddCheckbox("Pinned", false).
    OnToggle(func(_ *application.Context, checked bool) {})

appearance := inspector.AddSection("Appearance").SetCollapsible(true)
priority := appearance.AddSlider("Priority", 0, 5, 0)
priority.OnValueChange(func(_ *application.Context, value float64) {})
appearance.AddStepper("Indent", 0, 8, 1, 2)
appearance.AddSegmented("Align", []string{"Left", "Centre", "Right"}, 0)
appearance.AddColorWell("Tint", application.NewRGB(0, 122, 255)).
    OnColorChange(func(_ *application.Context, colour application.RGBA) {})
appearance.AddDatePicker("Due", time.Time{}).
    OnDateChange(func(_ *application.Context, t time.Time) {})
appearance.AddButton("Reset").OnClick(func(*application.Context) {
    priority.SetFloatValue(0)
})

statistics := inspector.AddSection("Statistics")
words := statistics.AddLabel("Words", "0")

// Programmatic setters never fire the callbacks.
words.SetValue("128")
```

Jenis kontrol beserta setter-nya:

| Kontrol | Setter | Callback |
|---------|---------|----------|
| `AddLabel` | `SetValue` | tidak ada |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled`, dan `SetHidden` berlaku untuk semua jenis. Bagian dapat dibuat dapat dilipat, dan bagian maupun kontrol dapat dipindahkan atau dihapus kapan saja.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Aksesori

`MacAccessory` adalah bilah kontrol native. Pilih tata letak saat membuatnya: `MacAccessoryLayoutLeading` berada di sebelah tombol jendela, `MacAccessoryLayoutTrailing` di tepi akhir titlebar, `MacAccessoryLayoutBottom` membentang selebar jendela di bawah titlebar dan toolbar, sedangkan `MacAccessoryLayoutTop` berada di bagian atas panel tampilan terbagi. Tambahkan semua kontrol sebelum memasang aksesori.

```go
// Next to the window buttons.
folders := application.NewMacAccessory(application.MacAccessoryLayoutLeading)
folders.AddSegmented([]string{"Inbox", "Starred"}, 0).
    SetSegmentSymbols("tray", "star").
    OnSelectionChange(func(_ *application.Context, index int, label string) {})

// At the trailing edge of the titlebar.
tools := application.NewMacAccessory(application.MacAccessoryLayoutTrailing)
tools.AddSearch("Search mail").
    SetIncremental(true).
    SetWidth(200).
    OnSearch(func(_ *application.Context, query string) {})
tools.AddSymbolButton("square.and.pencil").
    SetTooltip("Compose").
    OnClick(func(*application.Context) {})

// A full-width strip beneath the titlebar and toolbar.
status := application.NewMacAccessory(application.MacAccessoryLayoutBottom).
    SetHeight(30).
    SetPreferredScrollEdgeEffectStyle(application.MacScrollEdgeEffectStyleSoft)
label := status.AddLabel("Saved").SetSymbol("checkmark.circle")
status.AddFlexibleSpace()
status.AddButton("Mark All Read").OnClick(func(*application.Context) {})

for _, accessory := range []*application.MacAccessory{folders, tools, status} {
    if err := window.AddTitlebarAccessory(accessory); err != nil {
        window.Error("titlebar accessory: %s", err)
    }
}

// Live updates and lifecycle.
label.SetText("Edited").SetSymbol("pencil.circle")
status.SetHidden(true)
tools.Remove()
```

Kontrol yang tersedia adalah `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace`, dan, untuk integrasi native, `AddNativeView`. `AddTitlebarAccessory` tersedia pada `WebviewWindow` maupun `NativeWindow`; pemanggilan sebelum jendela ada akan diantrekan dan diterapkan ketika jendela dibuat. `Remove` melepas aksesori agar dapat dipasang kembali di tempat lain, sedangkan `SetHidden` melipatnya di tempat.

### Aksesori panel

Pada macOS 26 dan yang lebih baru, aksesori dapat berada di bagian atas atau bawah panel tampilan terbagi, seperti bidang filter sidebar di Finder. Buat aksesori dengan tata letak `Top` atau `Bottom`, lalu pasang dengan `AddTopAccessory` atau `AddBottomAccessory` pada panel.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` memilih perlakuan `Automatic`, `Soft`, atau `Hard` dari AppKit untuk konten yang bergulir di belakang aksesori. Gaya eksplisit memerlukan macOS 26.1; pada rilis sebelumnya permintaan dilaporkan melalui penangan kesalahan jendela dan gaya otomatis tetap berlaku.

## Menggabungkan semuanya

Program ini membuat jendela tiga panel yang berisi sidebar, WebView, dan panel pemeriksa, serta toolbar yang mengikuti kedua pembatas.

```go title="main.go"
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:   "Notes",
		Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Notes",
		Width:  1100,
		Height: 700,
		URL:    "/",
		Mac: application.MacWindow{
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})

	sidebar := application.NewMacSidebar()
	notes := sidebar.AddSection("Notes")
	notes.AddItem("Saturday, slowly").SetSymbol("doc.text").OnClick(func(*application.Context) {
		app.Event.Emit("note:selected", "saturday")
	})

	inspector := application.NewMacInspector()
	inspector.AddSection("Document").AddTextField("Title", "Saturday, slowly").
		OnTextChange(func(_ *application.Context, value string) {
			app.Event.Emit("note:title", value)
		})

	split := application.NewMacSplitView().SetAutosaveName("notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddPrimaryContent()
	split.AddInspector(inspector).SetMinimumThickness(240).SetCollapsible(true)
	window.SetSplitView(split)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true).
		OnClick(func(*application.Context) {
			notes.AddItem("Untitled").SetSymbol("doc.text")
		})
	toolbar.AddFlexibleSpace()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	window.SetToolbar(toolbar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Antarmuka native dan frontend berkomunikasi melalui event dan layanan Wails seperti biasa. Sidebar mengirim `note:selected`, panel pemeriksa mengirim `note:title`, dan halaman mendengarkan melalui `Events.On` milik runtime.

## Jendela native

@note{type="caution" title="Eksperimental"}
`NativeWindow`, `NativeWindowManager`, dan `MacTextEditor` bersifat eksperimental di v3. API-nya sengaja dibuat terbatas dan dapat berubah ketika API jendela bersama dirancang ulang untuk v4.
@end

`NativeWindow` tidak memiliki WebView. Konten utamanya adalah `MacTextEditor`, yaitu `NSTextView` di dalam `NSScrollView`, dan jendela ini menerima jenis toolbar, tampilan terbagi, serta aksesori yang sama dengan `WebviewWindow`. Buat dengan `app.NativeWindow.New` atau `app.NativeWindow.NewWithOptions`, lalu temukan kembali dengan `Get` atau `GetByID`.

Jendela native baru dibuat setelah memiliki konten. Sediakan tampilan terbagi yang panel utamanya ditambahkan dengan `AddTextEditor`, baik melalui `NativeWindowOptions.SplitView` maupun dengan memanggil `SetSplitView`. Sebelum `app.Run`, tata letak diantrekan; pada aplikasi yang sedang berjalan, `SetSplitView` langsung membuat dan menampilkan jendela (kecuali `Hidden` diatur) serta mengembalikan kesalahan pembuatan, jika ada. Jendela tanpa tata letak tetap ditunda dan `Run` mengembalikan `ErrNativeWindowContentRequired`; tata letak tanpa editor teks ditolak dengan `ErrNativeWindowEditorRequired`. Gunakan `NativeWindowOptions.Toolbar` dan `NativeWindowOptions.SplitView`, bukan bidang `Mac.Toolbar` dan `Mac.SplitView` yang diabaikan oleh jendela native.

```go title="main.go"
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:       "Native Notes",
		NativeOnly: true,
	})

	editor := application.NewMacTextEditor()
	editor.OnChange(func(*application.Context) {
		// mark the document dirty; call editor.Text() only when saving
	})

	sidebar := application.NewMacSidebar()
	sidebar.AddSection("Files").AddItem("README.txt").
		SetSymbol("doc.plaintext").
		OnClick(func(*application.Context) {
			editor.SetText("Hello from AppKit")
		})

	split := application.NewMacSplitView().SetAutosaveName("native-notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddTextEditor(editor).SetMinimumThickness(400)

	window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Title:  "Native Notes",
		Width:  900,
		Height: 600,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				FullSizeContent: true,
				ToolbarStyle:    application.MacToolbarStyleUnified,
			},
		},
	})
	if err := window.SetSplitView(split); err != nil {
		log.Fatal(err)
	}

	toolbar := application.NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Save").SetSymbol("square.and.arrow.down").SetBordered(true).
		OnClick(func(*application.Context) {
			log.Printf("%d bytes", len(editor.Text()))
		})
	if err := window.SetToolbar(toolbar); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Jendela yang sama dapat dibuat dalam satu pemanggilan dari aplikasi yang sedang berjalan dengan meneruskan antarmuka jendela sebagai opsi:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` menyediakan `SetText`, `Text`, `SetEditable`, `OnChange`, dan `Focus`. Pemanggilan `SetText` secara terprogram tidak memicu `OnChange`, sehingga memuat berkas tidak menandainya sebagai berubah. `Text` membaca seluruh dokumen dari AppKit, jadi panggil saat Anda memerlukan isinya, bukan setiap kali ada perubahan.

Dua opsi menjaga aplikasi yang hanya menggunakan komponen native tetap ringan:

- `NativeOnly: true` dalam `application.Options` melewati transport frontend dan server aset saat runtime. Jangan buat `WebviewWindow` ketika opsi ini diatur.
- Tag build `wails_native` menghilangkan kode WebView, frontend, dan pemutakhiran sepenuhnya dari biner saat kompilasi, serta mengatur `NativeOnly` secara otomatis:

```sh
go build -tags wails_native .
```

Dukungan satu instans juga dikecualikan dari build `wails_native`; tambahkan tag `wails_single_instance` jika Anda memerlukannya.

## Tab jendela

macOS dapat mengelompokkan jendela ke dalam tab. Mode tab ditetapkan saat jendela dibuat, jadi atur `Mac.TabbingMode` ke `MacWindowTabbingModePreferred` atau `MacWindowTabbingModeAutomatic` pada setiap jendela yang akan ikut serta. `MacWindowTabbingModeDisallowed` (dan nilai bawaan yang tidak diatur) mengecualikan jendela dari grup tab.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

Operasi tab memerlukan jendela native yang sudah aktif, jadi panggil dari penangan menu, metode layanan, atau kode lain yang berjalan setelah `app.Run` dimulai. `TabGroup` mengembalikan handle `MacWindowTabGroup` yang mencari grup native pada setiap pemanggilan; handle nil aman digunakan dan setiap metodenya mengembalikan nilai nol.

```go
second := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 2",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModeAutomatic,
    },
})
if err := first.AddTab(second, application.MacTabOrderAbove); err != nil {
    app.Logger.Error("add tab", "error", err)
}
second.SetTabTitle("Draft")

if group := first.TabGroup(); group != nil {
    group.SelectNext()
    group.ToggleTabBar()
    for _, member := range group.Windows() {
        app.Logger.Info("tab", "name", member.Name())
    }
}

first.MoveTabToNewWindow()
first.MergeAllWindows()
```

`MacWindowTabGroup` menyediakan `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible`, dan `ToggleTabOverview`. `app.Window.TabGroups` mencantumkan setiap grup yang berisi `WebviewWindow`. `AddNativeTab` menambahkan `NativeWindow` ke grup jendela WebView, dan `SetTabTooltip` mengatur teks saat penunjuk berada di atas tab. Di luar macOS, metode tab mengembalikan `ErrMacWindowTabsUnsupported`.

## Persyaratan versi

Semua yang ada di halaman ini hanya untuk macOS. Pada rilis yang lebih lama, fitur mengalami penurunan kemampuan secara bertahap seperti dijelaskan di bawah; API Go tetap sama di semua platform.

| Fitur | macOS minimum | Perilaku pada rilis sebelumnya |
|---------|---------------|-------------------------------|
| Tab jendela | 10.12 | tidak tersedia |
| Grup toolbar, item berbatas, `AddMenu` | 10.15 | item menu dihilangkan; grup menggunakan tampilan lama |
| SF Symbols (`SetSymbol` pada item toolbar, sidebar, daftar, dan aksesori) | 11 | tidak ada gambar yang ditampilkan |
| `AddSearch` sebagai `NSSearchToolbarItem`, `SetNavigational`, pemisah pelacak sidebar | 11 | pencarian menggunakan bidang pencarian biasa sebagai pengganti; pemisah dihilangkan |
| Peran panel pemeriksa dan daftar konten dalam tampilan terbagi | 11 | panel yang sama ditempatkan dalam item tampilan terbagi biasa |
| Gaya daftar konten | 11 | diabaikan |
| `SetCenteredItems` | 13 | diabaikan |
| Tombol panel pemeriksa dan pemisah pelacaknya | 14 | Wails menyediakan tombol native; pemisah dihilangkan |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` pada item toolbar | 26 | disimpan dan diterapkan jika tersedia |
| Aksesori panel (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Gaya efek tepi gulir eksplisit | 26.1 | dilaporkan melalui penangan kesalahan jendela; gaya otomatis tetap berlaku |

Aksesori titlebar, tampilan terbagi, sidebar, panel pemeriksa, dan editor teks tidak memiliki persyaratan versi di luar minimum Wails.

## Contoh

Setiap contoh adalah aplikasi lengkap yang dapat dijalankan:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): editor catatan dengan toolbar, sidebar, daftar konten, panel pemeriksa, penyedia berbagi, penyesuaian toolbar, dan aksesori panel.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): aksesori titlebar di awal, akhir, dan bawah yang mengendalikan tampilan kotak surat.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): mode tab, `AddTab`, `AddNativeTab`, dan API grup tab dari menu.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): editor teks `NativeWindow` tanpa WebView yang dibuat dengan tag `wails_native`.

## Langkah selanjutnya

Antarmuka jendela di halaman ini adalah salah satu bagian dari aplikasi macOS native. Bagian lainnya adalah perilaku aplikasi: jendela dokumen dengan ikon proksi dan penempatan berjenjang, menu dengan simbol dan lencana, menu Dock dengan progres, peringatan dan panel native, item status yang dapat dihapus, umpan balik haptik dan ucapan, papan klip kaya format dan penyeretan keluar, serta informasi sistem tentang izin, daya, dan lokal. API tersebut dibahas dalam panduan [Integrasi Platform macOS](/guides/macos-platform-integration).
