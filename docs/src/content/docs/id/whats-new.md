---
title: Apa yang Baru di Wails v3
description: Temukan peningkatan utama dan fitur baru di Wails v3
---

Wails v3 memperkenalkan perubahan signifikan dari v2. Wails v3 menggantikan
API deklaratif single-window dengan pendekatan prosedural yang lebih fleksibel. Desain
API baru ini meningkatkan keterbacaan kode dan menyederhanakan development, terutama
untuk aplikasi multi-window yang kompleks.

Wails v3 mewakili evolusi substansial dalam cara aplikasi desktop
dapat dibangun menggunakan Go dan teknologi web.

## Multiple Windows

Wails v3 memperkenalkan kemampuan untuk membuat dan mengelola banyak window dalam
satu aplikasi. Fitur ini memungkinkan developer merancang antarmuka pengguna yang lebih kompleks dan
serbaguna, melampaui keterbatasan aplikasi single-window.

Setiap window dapat dikonfigurasi secara independen, memberikan fleksibilitas dalam hal
ukuran, posisi, konten, dan perilaku. Ini memungkinkan pembuatan aplikasi
dengan window terpisah untuk fungsionalitas berbeda, seperti antarmuka utama,
panel pengaturan, atau tampilan tambahan.

Developer dapat membuat, memanipulasi, dan mengelola window ini secara programatis,
memungkinkan antarmuka pengguna dinamis yang beradaptasi terhadap kebutuhan pengguna dan state
aplikasi.

:::tip[Multiple Windows]
<details><summary>Contoh</summary>

```go
package main

import (
   "embed"
   "log"
   
   "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {

   app := application.New(application.Options{
        Name:   "Multi Window Demo",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
   })
   
   window1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 1",
   })
   
   window2 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 2",
   })
   
   // muat html embedded dari embed.FS
   window1.SetURL("/")
   window1.Center()
   
   // Muat URL eksternal
   window2.SetURL("https://wails.io")
   
   err := app.Run()

   if err != nil {
	   log.Fatal(err.Error())
   }
}
```
</details>

:::

## Integrasi System Tray

Wails v3 memperkenalkan dukungan robust untuk fungsionalitas system tray, memungkinkan
aplikasi Anda mempertahankan kehadiran persisten di desktop pengguna. Fitur
ini sangat berguna untuk aplikasi yang perlu berjalan di
background atau menyediakan akses cepat ke fungsi utama.

Fitur utama integrasi system tray Wails v3 meliputi:

1. Window Attachment: Anda dapat mengasosiasikan window dengan ikon system tray. Saat
   diaktifkan, window ini akan di-center relatif terhadap posisi ikon,
   memberikan cara bagus untuk mengakses aplikasi dengan cepat.

2. Dukungan Menu Komprehensif: Buat menu kaya dan interaktif yang dapat diakses pengguna
   langsung dari ikon system tray. Ini memungkinkan aksi cepat
   tanpa perlu membuka window aplikasi penuh.

3. Tampilan Ikon Adaptif: Dukungan ikon light dan dark mode memastikan
   ikon system tray aplikasi tetap terlihat dan estetis
   di berbagai tema sistem. Template icon juga didukung di macOS.

:::tip[Systray]

<details><summary>Contoh</summary>


```go
package main

import (
    "log"
    "runtime"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
    app := application.New(application.Options{
        Name:        "Systray Demo",
        Mac: application.MacOptions{
            ActivationPolicy: application.ActivationPolicyAccessory,
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Width:       500,
        Height:      800,
        Frameless:   true,
        AlwaysOnTop: true,
        Hidden:      true,
        Windows: application.WindowsWindow{
            HiddenOnTaskbar: true,
        },
    })

    systemTray := app.SystemTray.New()

    // Dukungan template icon di macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Dukungan ikon light/dark mode
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Dukungan menu
    myMenu := app.Menu.New()
    myMenu.Add("Hello World!").OnClick(func(_ *application.Context) {
        println("Hello World!")
    })
    systemTray.SetMenu(myMenu)

    // Ini akan center window ke ikon systray dengan offset 5px
    // Window otomatis ditampilkan saat ikon systray diklik
    // dan disembunyikan saat window kehilangan fokus
    systemTray.AttachWindow(window).WindowOffset(5)

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```
</details>
:::

## Pembuatan bindings yang ditingkatkan

Wails v3 memperkenalkan peningkatan signifikan dalam cara bindings di-generate untuk
proyek Anda. Bindings adalah lem yang menghubungkan backend Go ke
frontend, memungkinkan komunikasi seamless antara keduanya.

Pembuatan binding sekarang dilakukan menggunakan static analyzer canggih yang
secara radikal meningkatkan proses pembuatan binding. Analyzer menawarkan kecepatan yang ditingkatkan dan
mempertahankan kualitas kode dengan mempertahankan komentar dan nama parameter.

Proses pembuatan binding telah disederhanakan, hanya memerlukan satu
perintah: `wails3 generate bindings`.

:::tip[Bindings]

<details><summary>Contoh</summary>

```js
// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// Generated layout (excerpt): frontend/bindings/<full-go-import-path>/greetservice.js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * Greet greets a person
 * @param {string} $0
 * @returns {Promise<string>}
 */
export function Greet($0) {
    return $Call.ByID(1411160069, $0);
}

/**
 * GreetPerson greets a person
 * @param {main.Person} $0
 * @returns {Promise<string>}
 */
export function GreetPerson($0) {
    return $Call.ByID(4021313248, $0);
}
```

</details>

:::

## Sistem build yang ditingkatkan

Wails v3 memperkenalkan sistem build yang lebih fleksibel dan transparan, mengatasi
keterbatasan pendahulunya. Di v2, proses build sebagian besar opaque dan
sulit dikustomisasi, yang bisa membuat frustrasi developer yang menginginkan lebih
kontrol atas proses build proyek mereka.

Semua heavy lifting yang dilakukan sistem build v2, seperti pembuatan icon dan
pembuatan manifest, telah ditambahkan sebagai perintah tool di CLI. Kami
menggabungkan [Taskfile](https://taskfile.dev) ke CLI untuk mengorkestrasi
panggilan ini agar memberikan developer experience yang sama seperti v2. Namun, pendekatan ini
membawa keseimbangan ultimate fleksibilitas dan kemudahan penggunaan karena Anda sekarang
dapat menyesuaikan proses build sesuai kebutuhan.

Anda bahkan bisa memakai make jika itu yang Anda suka!

:::tip[Taskfile.yml]

<details><summary>Contoh</summary>

```yaml "Snippet from build/Taskfile.darwin.yml"
darwin:build:
  summary: Builds the application for macOS
  platforms:
    - darwin
  cmds:
    - task: common:go:mod:tidy
    - task: common:build:frontend
    - task: common:generate:icons
    - task: darwin:build:app
  env:
    CGO_CFLAGS: "-mmacosx-version-min=10.15"
    CGO_LDFLAGS: "-mmacosx-version-min=10.15"
    MACOSX_DEPLOYMENT_TARGET: "10.15"
```
</details>

:::

## Events yang ditingkatkan

Wails sekarang mengirim events untuk berbagai operasi runtime dan aktivitas sistem.
Ini memungkinkan aplikasi Anda merespons events ini secara real-time.
Selain itu, events lintas platform (common) tersedia, memungkinkan Anda
menulis metode penanganan event konsisten yang berfungsi di sistem operasi berbeda.

Event hook dapat didaftarkan untuk menangani events spesifik secara sinkron. Berbeda dengan
method `On`, hook ini memungkinkan Anda membatalkan event jika diperlukan. Kasus penggunaan umum
adalah menampilkan dialog konfirmasi sebelum menutup window. Ini memberi
Anda lebih banyak kontrol atas alur event dan pengalaman pengguna.

:::tip[Contoh penanganan event]

<details><summary>Contoh</summary>

```go
package main

import (
    "embed"
    "log"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

    app := application.New(application.Options{
        Name:        "Events Demo",
        Description: "A demo of the Events API",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Penanganan custom event — App.Event.On(name, func(e *CustomEvent))
    app.Event.On("myevent", func(e *application.CustomEvent) {
        log.Printf("[Go] CustomEvent received: %+v\n", e)
    })

    // Application events spesifik OS — App.Event.OnApplicationEvent(eventType, func(e *ApplicationEvent))
    app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
        println("events.Mac.ApplicationDidFinishLaunching fired!")
    })

    // Events platform-agnostic
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
        println("events.Common.ApplicationStarted fired!")
    })

    win1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Takes 3 attempts to close me!",
    })

    var countdown = 3

    // Daftarkan hook untuk membatalkan penutupan window
    win1.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        countdown--
        if countdown == 0 {
            println("Closing!")
            return
        }
        println("Nope! Not closing!")
        e.Cancel()
    })

    win1.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        println("[Event] Window focus!")
    })

    err := app.Run()

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

</details>

:::

## Wails Markup Language (wml)

Fitur eksperimental untuk memanggil method runtime menggunakan plain html, mirip
[htmx](https://htmx.org).

:::tip[Contoh wml]

<details><summary>Contoh</summary>

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Wails ML Demo</title>
  </head>
  <body style="margin-top:50px; color: white; background-color: #191919">
    <h2>Wails ML Demo</h2>
    <p>This application contains no Javascript!</p>
    <button wml-event="button-pressed">Press me!</button>
    <button wml-event="delete-things" wml-confirm="Are you sure?">
      Delete all the things!
    </button>
    <button wml-window="Close" wml-confirm="Are you sure?">
      Close the Window?
    </button>
    <button wml-window="Center">Center</button>
    <button wml-window="Minimise">Minimise</button>
    <button wml-window="Maximise">Maximise</button>
    <button wml-window="UnMaximise">UnMaximise</button>
    <button wml-window="Fullscreen">Fullscreen</button>
    <button wml-window="UnFullscreen">UnFullscreen</button>
    <button wml-window="Restore">Restore</button>
    <div
      style="width: 200px; height: 200px; border: 2px solid white;"
      wml-event="hover"
      wml-trigger="mouseover"
    >
      Hover over me
    </div>
  </body>
</html>
```

</details>

:::

## Contoh

Lebih banyak contoh tersedia di direktori
[examples](https://github.com/wailsapp/wails/tree/master/v3/examples).
Lihat mereka!
