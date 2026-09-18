---
title: "Yang Baru di Wails v3"
description: "Temukan peningkatan utama dan fitur baru di Wails v3"
slug: "whats-new"
sourcePath: "whats-new.md"
---

Wails v3 menghadirkan perubahan signifikan dibandingkan v2. API deklaratif berjendela tunggal digantikan dengan pendekatan prosedural yang lebih fleksibel. Desain API baru ini meningkatkan keterbacaan kode dan menyederhanakan pengembangan, terutama untuk aplikasi multijendela yang kompleks.

Wails v3 merupakan evolusi besar dalam cara membangun aplikasi desktop menggunakan Go dan teknologi web.

## Beberapa Jendela

Wails v3 menghadirkan kemampuan untuk membuat dan mengelola beberapa jendela dalam satu aplikasi. Fitur ini memungkinkan pengembang merancang antarmuka pengguna yang lebih kompleks dan serbaguna, melampaui keterbatasan aplikasi berjendela tunggal.

Setiap jendela dapat dikonfigurasi secara independen, sehingga memberikan fleksibilitas dalam hal ukuran, posisi, konten, dan perilaku. Hal ini memungkinkan pembuatan aplikasi dengan jendela terpisah untuk berbagai fungsi, seperti antarmuka utama, panel pengaturan, atau tampilan tambahan.

Pengembang dapat membuat, memanipulasi, dan mengelola jendela tersebut secara terprogram, sehingga memungkinkan antarmuka pengguna dinamis yang menyesuaikan diri dengan kebutuhan pengguna dan status aplikasi.

@note{type="tip" title="Beberapa Jendela"}
@details{title="Contoh"}
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
   
   // load the embedded html from the embed.FS
   window1.SetURL("/")
   window1.Center()
   
   // Load an external URL
   window2.SetURL("https://wails.io")
   
   err := app.Run()

   if err != nil {
	   log.Fatal(err.Error())
   }
}
```

@end

@end

## Integrasi Baki Sistem

Wails v3 menghadirkan dukungan tangguh untuk fungsionalitas baki sistem, sehingga aplikasi Anda dapat terus hadir di desktop pengguna. Fitur ini sangat berguna bagi aplikasi yang perlu berjalan di latar belakang atau menyediakan akses cepat ke fungsi utama.

Fitur utama integrasi baki sistem Wails v3 meliputi:

1. Penautan Jendela: Anda dapat mengaitkan jendela dengan ikon baki sistem. Saat diaktifkan, jendela ini akan dipusatkan relatif terhadap posisi ikon, sehingga memberikan cara praktis untuk mengakses aplikasi Anda dengan cepat.

2. Dukungan Menu Lengkap: Buat menu interaktif yang kaya dan dapat diakses pengguna langsung dari ikon baki sistem. Dengan demikian, tindakan cepat dapat dilakukan tanpa perlu membuka jendela aplikasi secara penuh.

3. Tampilan Ikon Adaptif: Dukungan ikon untuk mode terang dan gelap memastikan ikon baki sistem aplikasi Anda tetap terlihat dan tampak menarik pada berbagai tema sistem. Ikon templat juga didukung di macOS.

@note{type="tip" title="Baki Sistem"}
@details{title="Contoh"}
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

    // Support for template icons on macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Support for light/dark mode icons
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Support for menu
    myMenu := app.Menu.New()
    myMenu.Add("Hello World!").OnClick(func(_ *application.Context) {
        println("Hello World!")
    })
    systemTray.SetMenu(myMenu)

    // This will center the window to the systray icon with a 5px offset
    // It will automatically be shown when the systray icon is clicked
    // and hidden when the window loses focus
    systemTray.AttachWindow(window).WindowOffset(5)

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

@end

@end

## Peningkatan Pembuatan Binding

Wails v3 menghadirkan peningkatan signifikan pada cara binding dibuat untuk proyek Anda. Binding merupakan penghubung antara backend Go dan frontend Anda, yang memungkinkan komunikasi lancar di antara keduanya.

Pembuatan binding kini dilakukan menggunakan penganalisis statis canggih yang secara drastis meningkatkan proses tersebut. Penganalisis ini memberikan kecepatan yang lebih tinggi dan mempertahankan kualitas kode dengan menjaga komentar serta nama parameter.

Proses pembuatan binding telah disederhanakan dan kini hanya memerlukan satu perintah: `wails3 generate bindings`.

@note{type="tip" title="Binding"}
@details{title="Contoh"}
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

@end

@end

## Peningkatan Sistem Build

Wails v3 menghadirkan sistem build yang lebih fleksibel dan transparan untuk mengatasi keterbatasan pendahulunya. Di v2, proses build sebagian besar tidak transparan dan sulit disesuaikan, sehingga dapat menyulitkan pengembang yang menginginkan kendali lebih besar atas proses build proyek mereka.

Semua pekerjaan berat yang sebelumnya dilakukan oleh sistem build v2, seperti pembuatan ikon dan manifes, telah ditambahkan sebagai perintah alat di CLI. Kami telah mengintegrasikan [Taskfile](https://taskfile.dev) ke dalam CLI untuk mengatur pemanggilan tersebut dan memberikan pengalaman pengembang yang sama seperti v2. Namun, pendekatan ini memberikan keseimbangan terbaik antara fleksibilitas dan kemudahan penggunaan karena kini Anda dapat menyesuaikan proses build dengan kebutuhan Anda.

Anda bahkan dapat menggunakan make jika itu pilihan Anda!

@note{type="tip" title="Taskfile.yml"}
@details{title="Contoh"}
```yaml {title="build/Taskfile.darwin.yml"}
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

@end

@end

## Peningkatan Peristiwa

Wails kini memancarkan peristiwa untuk berbagai operasi runtime dan aktivitas sistem. Hal ini memungkinkan aplikasi Anda merespons peristiwa tersebut secara real-time. Selain itu, tersedia peristiwa lintas platform (umum), sehingga Anda dapat menulis metode penanganan peristiwa yang konsisten dan berfungsi di berbagai sistem operasi.

Hook peristiwa dapat didaftarkan untuk menangani peristiwa tertentu secara sinkron. Berbeda dengan metode `On`, hook ini memungkinkan Anda membatalkan peristiwa jika diperlukan. Kasus penggunaan yang umum adalah menampilkan dialog konfirmasi sebelum menutup jendela. Hal ini memberi Anda kendali lebih besar atas alur peristiwa dan pengalaman pengguna.

@note{type="tip" title="Contoh penanganan peristiwa"}
@details{title="Contoh"}
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

    // Custom event handling — App.Event.On(name, func(e *CustomEvent))
    app.Event.On("myevent", func(e *application.CustomEvent) {
        log.Printf("[Go] CustomEvent received: %+v\n", e)
    })

    // OS-specific application events — App.Event.OnApplicationEvent(eventType, func(e *ApplicationEvent))
    app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
        println("events.Mac.ApplicationDidFinishLaunching fired!")
    })

    // Platform-agnostic events
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
        println("events.Common.ApplicationStarted fired!")
    })

    win1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Takes 3 attempts to close me!",
    })

    var countdown = 3

    // Register a hook to cancel the window closing
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

@end

@end

## Wails Markup Language (wml)

Fitur eksperimental untuk memanggil metode runtime menggunakan HTML biasa, serupa dengan [htmx](https://htmx.org).

@note{type="tip" title="Contoh wml"}
@details{title="Contoh"}
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

@end

@end

## Contoh

Contoh lainnya tersedia di direktori [examples](https://github.com/wailsapp/wails/tree/master/v3/examples). Silakan lihat!
