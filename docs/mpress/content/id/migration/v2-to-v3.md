---
title: "Bermigrasi dari v2 ke v3"
description: "Panduan lengkap untuk memigrasikan aplikasi Wails v2 Anda ke v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 merupakan **penulisan ulang secara menyeluruh** dengan peningkatan signifikan pada arsitektur, performa, dan pengalaman developer. Panduan ini membantu Anda memigrasikan aplikasi v2 ke v3.

**Perubahan utama:**

- Struktur aplikasi baru
- Sistem binding yang ditingkatkan
- Pengelolaan jendela yang disempurnakan
- Sistem peristiwa yang lebih baik
- Konfigurasi yang disederhanakan

**Waktu migrasi:** 1-4 jam untuk aplikasi pada umumnya

## Perubahan yang Tidak Kompatibel

### Inisialisasi Aplikasi

Di v2, penyiapan aplikasi, konfigurasi jendela, dan eksekusi semuanya digabungkan dalam satu pemanggilan `wails.Run()`. Pendekatan monolitik ini menyulitkan pembuatan beberapa jendela, penanganan kesalahan pada tahap yang berbeda, atau pengujian setiap komponen aplikasi secara terpisah.

v3 memisahkan hal-hal tersebut menjadi beberapa fase tersendiri: pembuatan aplikasi, pembuatan jendela, dan eksekusi. Pemisahan ini memberi Anda kendali eksplisit atas setiap tahap siklus hidup aplikasi serta menjadikan kode lebih modular dan mudah diuji.

**v2:**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3:**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**Mengapa ini lebih baik:**

- **Dukungan multi-jendela**: Anda dapat membuat jendela secara dinamis kapan saja, tidak hanya saat aplikasi dimulai
- **Penanganan kesalahan yang lebih baik**: Setiap fase dapat divalidasi secara terpisah dengan penanganan kesalahan yang tepat
- **Kode yang lebih jelas**: Pemisahan ini memperjelas apa yang terjadi pada setiap tahap
- **Lebih mudah diuji**: Anda dapat menguji penyiapan aplikasi tanpa menjalankan event loop
- **Lebih fleksibel**: Jendela dapat dibuat, dihapus, dan dibuat ulang sepanjang siklus hidup aplikasi

### Binding

Di v2, setiap struct yang di-binding memerlukan field konteks dan metode `startup(ctx)` untuk menerima konteks runtime. Hal ini menciptakan keterikatan yang erat antara logika bisnis dan runtime Wails sehingga kode lebih sulit diuji dan dipahami.

v3 memperkenalkan pola layanan, yang memungkinkan struct Anda sepenuhnya mandiri dan tidak perlu menyimpan konteks runtime. Jika suatu layanan memerlukan akses ke instans aplikasi, layanan tersebut menerimanya secara eksplisit melalui injeksi dependensi, bukan melalui penerusan konteks secara implisit.

**v2:**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**Mengapa ini lebih baik:**

- **Tanpa dependensi implisit**: Layanan merupakan struct Go biasa tanpa dependensi runtime tersembunyi
- **Pengujian lebih mudah**: Anda dapat menguji metode layanan tanpa membuat mock konteks Wails
- **Kode yang lebih jelas**: Dependensi dinyatakan secara eksplisit (diteruskan sebagai argumen konstruktor), bukan disembunyikan dalam field konteks
- **Pengorganisasian yang lebih baik**: Layanan dapat dikelompokkan berdasarkan domain, bukan semuanya ditempatkan dalam satu struct `App`
- **Inisialisasi yang tepat**: Gunakan metode `ServiceStartup()` saat Anda memerlukan inisialisasi agar proses tersebut dinyatakan secara eksplisit

### Runtime

Di v2, semua operasi runtime mengharuskan konteks diteruskan ke fungsi global dari paket `runtime`. Hal ini menciptakan keterikatan yang erat dengan objek konteks di seluruh basis kode Anda dan membuat API terasa prosedural, bukan berorientasi objek.

v3 mengganti runtime berbasis konteks dengan pemanggilan metode langsung pada objek aplikasi dan jendela. Operasi dipanggil langsung pada objek yang dipengaruhinya sehingga kode menjadi lebih intuitif dan berorientasi objek.

**v2:**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3:**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**Mengapa ini lebih baik:**

- **Desain berorientasi objek**: Metode dipanggil pada objek yang dipengaruhinya (jendela, aplikasi, menu, dan sebagainya)
- **Maksud yang lebih jelas**: `window.SetTitle()` lebih mudah dipahami daripada `runtime.WindowSetTitle(ctx, ...)`
- **Dukungan IDE yang lebih baik**: Pelengkapan otomatis berfungsi dengan baik ketika metode berada pada objek
- **Pengelolaan multi-jendela yang lebih jelas**: Saat terdapat beberapa jendela, Anda secara eksplisit memilih jendela yang akan dikenai operasi
- **Tanpa penerusan konteks**: Anda tidak perlu meneruskan konteks melalui setiap fungsi

### Binding Frontend

Di v2, binding disusun berdasarkan paket Go dan nama struct, yang biasanya menghasilkan path seperti `wailsjs/go/main/App`. Struktur ini tidak mencerminkan pengelompokan logis dan menyulitkan pencarian fungsionalitas yang saling terkait.

v3 menyusun binding berdasarkan nama layanan dan modul aplikasi sehingga menghasilkan struktur logis yang lebih jelas. Binding dibuat dalam direktori `bindings` yang disusun berdasarkan nama aplikasi dan nama layanan Anda, sehingga fungsionalitas yang tersedia lebih mudah dipahami.

**v2:**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**Mengapa ini lebih baik:**

- **Pengorganisasian logis**: Binding dikelompokkan berdasarkan nama layanan, bukan struktur paket Go
- **Impor yang lebih jelas**: Path mencerminkan logika domain (greetservice), bukan struktur file (main/App)
- **Lebih mudah ditemukan**: Anda dapat menelusuri binding berdasarkan fitur, bukan berdasarkan struktur teknis
- **Penamaan yang konsisten**: Pengorganisasian berbasis layanan selaras dengan arsitektur backend Anda
- **Path yang lebih sederhana**: Tidak ada lagi prefiks `../wailsjs/go`—cukup `./bindings`

### Event

Di v2, event menggunakan parameter variadik `interface{}` dan mengharuskan konteks diteruskan ke setiap fungsi event. Handler event menerima data tanpa tipe yang memerlukan assertion tipe secara manual sehingga sistem event rentan terhadap kesalahan dan sulit di-debug.

v3 memperkenalkan objek event bertipe dan menghapus keharusan penggunaan konteks. Handler event menerima objek event yang tepat dengan data bertipe sehingga sistem event lebih andal dan lebih mudah digunakan.

**v2:**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3:**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**Mengapa ini lebih baik:**

- **Keamanan tipe**: Event menggunakan objek event yang tepat, bukan `...interface{}`
- **Debugging yang lebih mudah**: Objek event berisi metadata seperti nama event sehingga debugging menjadi lebih mudah
- **API yang lebih jelas**: `app.Event.On()` dan `app.Event.Emit()` lebih intuitif daripada fungsi runtime
- **Tidak memerlukan konteks**: Event bekerja langsung pada objek aplikasi tanpa meneruskan konteks
- **Handler yang lebih sederhana**: Handler event memiliki signature yang jelas, bukan parameter variadik

### Jendela

v2 hanya mendukung satu jendela per aplikasi. Jendela dibuat saat aplikasi dimulai dan semua operasi jendela dilakukan melalui fungsi runtime yang secara implisit menargetkan satu-satunya jendela tersebut.

v3 memperkenalkan dukungan multi-jendela native sebagai fitur inti. Setiap jendela merupakan objek mandiri dengan metode dan siklus hidupnya sendiri. Anda dapat membuat, mengelola, dan memusnahkan beberapa jendela secara dinamis sepanjang masa aktif aplikasi.

**v2:**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3:**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**Mengapa ini lebih baik:**

- **Aplikasi multi-jendela**: Buat aplikasi dengan beberapa jendela independen (dasbor, preferensi, alat, dan sebagainya)
- **Referensi jendela yang eksplisit**: Setiap jendela adalah objek yang dapat Anda simpan dan manipulasi secara langsung
- **Pembuatan jendela dinamis**: Buat dan musnahkan jendela kapan saja selama runtime
- **Status jendela independen**: Setiap jendela memiliki event, properti, dan siklus hidupnya sendiri
- **Arsitektur yang lebih baik**: Pengelolaan jendela berorientasi objek, bukan berbasis konteks

## Langkah Migrasi

### Langkah 1: Perbarui Dependensi

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**Perbarui:**

```bash
go get github.com/wailsapp/wails/v3@latest
go mod tidy
```

### Langkah 2: Perbarui main.go

**v2:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### Langkah 3: Ubah Struct App Menjadi Layanan

**v2:**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### Langkah 4: Perbarui Pemanggilan Runtime

**v2:**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3:**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### Langkah 5: Perbarui Frontend

**Buat binding baru:**

```bash
wails3 generate bindings
```

**Perbarui impor:**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**Perbarui penanganan event:**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### Langkah 6: Perbarui Konfigurasi

**v2 (wails.json):**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3 (wails.json):**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## Pemetaan Fitur

### Dialog

**v2:**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3:**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### Menu

**v2:**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3:**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Baki Sistem

**v2:**

```go
// Not available in v2
```

**v3:**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## Masalah Umum

### Masalah: Binding tidak ditemukan

**Masalah:** Kesalahan impor setelah migrasi

**Solusi:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### Masalah: Kesalahan konteks

**Masalah:** `ctx` tidak tersedia

**Solusi:**

Simpan referensi aplikasi sebagai gantinya:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Masalah: Metode jendela tidak berfungsi

**Masalah:** `runtime.WindowSetTitle()` tidak tersedia

**Solusi:**

Gunakan metode jendela secara langsung:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Masalah: Peristiwa tidak dipicu

**Masalah:** Peristiwa telah didaftarkan, tetapi tidak diterima

**Solusi:**

Pastikan nama peristiwa sama persis:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Menguji Migrasi

### Daftar Periksa

- [ ] Aplikasi dimulai tanpa kesalahan
- [ ] Semua binding berfungsi
- [ ] Peristiwa dikirim dan diterima
- [ ] Jendela terbuka dan tertutup dengan benar
- [ ] Menu berfungsi (jika berlaku)
- [ ] Dialog berfungsi (jika berlaku)
- [ ] Baki sistem berfungsi (jika berlaku)
- [ ] Proses build berfungsi
- [ ] Build produksi berfungsi

### Perintah Pengujian

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Manfaat v3

### Performa

- **Startup lebih cepat** - Inisialisasi yang dioptimalkan
- **Penggunaan memori lebih rendah** - Penggunaan sumber daya yang efisien
- **Bridge yang lebih baik** - Overhead panggilan sebesar &lt;1ms

### Fitur

- **Multi-jendela** - Dukungan native
- **Baki sistem** - Bawaan
- **Peristiwa yang lebih baik** - API bertipe yang lebih sederhana
- **Layanan** - Pengorganisasian kode yang lebih baik

### Pengalaman Developer

- **Keamanan tipe** - Dukungan TypeScript penuh
- **Kesalahan yang lebih baik** - Pesan kesalahan yang jelas
- **Hot reload** - Pengembangan lebih cepat
- **Dokumentasi yang lebih baik** - Panduan lengkap

## Mendapatkan Bantuan

### Sumber Daya

- [Dokumentasi](/quick-start/why-wails/)
- [Komunitas Discord](https://discord.gg/JDdSxwjhGf)
- [Issue GitHub](https://github.com/wailsapp/wails/issues)
- [Contoh](https://github.com/wailsapp/wails/tree/master/v3/examples)

### Pertanyaan Umum

**T: Dapatkah saya menjalankan v2 dan v3 secara berdampingan?** J: Ya, keduanya menggunakan jalur impor yang berbeda.

**T: Apakah v3 siap digunakan dalam produksi?** J: v3 merupakan perangkat lunak beta dengan API desktop yang stabil. Sudah ada aplikasi yang menggunakannya di lingkungan produksi, tetapi lakukan pengujian menyeluruh sebelum deployment. v2 tetap merupakan riliis stabil saat ini.

**T: Apakah v2 akan tetap dipelihara?** J: Ya, v2 akan menerima pembaruan kritis.

**T: Berapa lama waktu yang diperlukan untuk migrasi?** J: 1-4 jam untuk aplikasi pada umumnya.

## Langkah Berikutnya

@cards{cols="2"}
🚀 Mulai Cepat
Mulai gunakan Wails v3.

[Pelajari Lebih Lanjut →](/quick-start/installation/)

---
★ Konsep Inti
Pahami arsitektur v3.

[Pelajari Lebih Lanjut →](/concepts/architecture/)

---
◆ Binding
Pelajari sistem binding yang baru.

[Pelajari Lebih Lanjut →](/features/bindings/methods/)

---
📖 Contoh
Lihat contoh lengkap v3.

[Lihat Contoh →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau [buat issue](https://github.com/wailsapp/wails/issues).
