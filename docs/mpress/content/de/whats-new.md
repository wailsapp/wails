---
title: "Neuerungen in Wails v3"
description: "Entdecken Sie die wichtigsten Verbesserungen und neuen Funktionen in Wails v3"
slug: "whats-new"
sourcePath: "whats-new.md"
---

Wails v3 bringt gegenüber v2 wesentliche Änderungen mit sich. Die deklarative API für ein einzelnes Fenster wird durch einen flexibleren prozeduralen Ansatz ersetzt. Dieses neue API-Design verbessert die Lesbarkeit des Codes und vereinfacht die Entwicklung, insbesondere bei komplexen Anwendungen mit mehreren Fenstern.

Wails v3 stellt eine bedeutende Weiterentwicklung beim Erstellen von Desktopanwendungen mit Go und Webtechnologien dar.

## Mehrere Fenster

Mit Wails v3 können Sie mehrere Fenster innerhalb einer einzigen Anwendung erstellen und verwalten. Dadurch lassen sich komplexere und vielseitigere Benutzeroberflächen entwickeln, die über die Einschränkungen von Anwendungen mit nur einem Fenster hinausgehen.

Jedes Fenster kann unabhängig konfiguriert werden. Größe, Position, Inhalt und Verhalten lassen sich flexibel festlegen. So können Anwendungen separate Fenster für unterschiedliche Funktionen bereitstellen, etwa Hauptoberflächen, Einstellungsbereiche oder ergänzende Ansichten.

Entwickler können diese Fenster programmgesteuert erstellen, bearbeiten und verwalten. Dadurch entstehen dynamische Benutzeroberflächen, die sich an die Anforderungen der Benutzer und den Zustand der Anwendung anpassen.

@note{type="tip" title="Mehrere Fenster"}
@details{title="Beispiel"}
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

## Systemleistenintegration

Wails v3 bietet umfassende Unterstützung für Funktionen der Systemleiste, sodass Ihre Anwendung dauerhaft auf dem Desktop des Benutzers verfügbar bleibt. Diese Funktion ist besonders für Anwendungen nützlich, die im Hintergrund ausgeführt werden müssen oder schnellen Zugriff auf zentrale Funktionen bieten sollen.

Die Systemleistenintegration von Wails v3 bietet unter anderem folgende Hauptfunktionen:

1. Fensterzuordnung: Sie können dem Symbol in der Systemleiste ein Fenster zuordnen. Bei der Aktivierung wird dieses Fenster relativ zur Position des Symbols zentriert und ermöglicht so einen schnellen Zugriff auf Ihre Anwendung.

2. Umfassende Menüunterstützung: Erstellen Sie umfangreiche interaktive Menüs, auf die Benutzer direkt über das Symbol in der Systemleiste zugreifen können. So lassen sich schnelle Aktionen ausführen, ohne das vollständige Anwendungsfenster öffnen zu müssen.

3. Adaptive Symboldarstellung: Die Unterstützung für Symbole im hellen und dunklen Modus stellt sicher, dass das Systemleistensymbol Ihrer Anwendung bei verschiedenen Systemdesigns sichtbar und ansprechend bleibt. Unter macOS werden außerdem Vorlagensymbole unterstützt.

@note{type="tip" title="Systemleiste"}
@details{title="Beispiel"}
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

## Verbesserte Bindings-Generierung

Wails v3 verbessert die Generierung von Bindings für Ihr Projekt erheblich. Bindings verbinden Ihr Go-Backend mit Ihrem Frontend und ermöglichen eine nahtlose Kommunikation zwischen beiden.

Die Bindings werden jetzt mit einem leistungsfähigen statischen Analysewerkzeug generiert, das den Generierungsprozess grundlegend verbessert. Es erhöht die Geschwindigkeit und bewahrt die Codequalität, indem Kommentare und Parameternamen erhalten bleiben.

Die Generierung der Bindings wurde vereinfacht und erfordert nur noch einen einzigen Befehl: `wails3 generate bindings`.

@note{type="tip" title="Bindings"}
@details{title="Beispiel"}
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

## Verbessertes Build-System

Wails v3 führt ein flexibleres und transparenteres Build-System ein, das die Einschränkungen seines Vorgängers beseitigt. In v2 war der Build-Prozess weitgehend undurchsichtig und nur schwer anpassbar. Das konnte Entwickler frustrieren, die mehr Kontrolle über den Build-Prozess ihres Projekts benötigten.

Alle aufwendigen Aufgaben, die das Build-System von v2 übernahm, etwa das Generieren von Symbolen und das Erstellen von Manifesten, wurden der CLI als Werkzeugbefehle hinzugefügt. Zur Orchestrierung dieser Aufrufe haben wir [Taskfile](https://taskfile.dev) in die CLI integriert, um dieselbe Entwicklererfahrung wie in v2 zu bieten. Dieser Ansatz schafft jedoch die bestmögliche Balance zwischen Flexibilität und Benutzerfreundlichkeit, da Sie den Build-Prozess jetzt an Ihre Anforderungen anpassen können.

Sie können sogar make verwenden, wenn Ihnen das lieber ist!

@note{type="tip" title="Taskfile.yml"}
@details{title="Beispiel"}
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

## Verbesserte Ereignisse

Wails löst jetzt Ereignisse für verschiedene Laufzeitvorgänge und Systemaktivitäten aus. Dadurch kann Ihre Anwendung in Echtzeit auf diese Ereignisse reagieren. Zusätzlich stehen plattformübergreifende (gemeinsame) Ereignisse zur Verfügung, sodass Sie einheitliche Methoden zur Ereignisbehandlung schreiben können, die unter verschiedenen Betriebssystemen funktionieren.

Zur synchronen Behandlung bestimmter Ereignisse können Ereignis-Hooks registriert werden. Anders als bei der Methode `On` können Sie mit diesen Hooks das Ereignis bei Bedarf abbrechen. Ein typischer Anwendungsfall ist die Anzeige eines Bestätigungsdialogs, bevor ein Fenster geschlossen wird. Dadurch haben Sie mehr Kontrolle über den Ereignisablauf und die Benutzerführung.

@note{type="tip" title="Beispiel für die Ereignisbehandlung"}
@details{title="Beispiel"}
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

Eine experimentelle Funktion zum Aufrufen von Laufzeitmethoden mit einfachem HTML, ähnlich wie [htmx](https://htmx.org).

@note{type="tip" title="Beispiel für wml"}
@details{title="Beispiel"}
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

## Beispiele

Weitere Beispiele finden Sie im Verzeichnis [examples](https://github.com/wailsapp/wails/tree/master/v3/examples). Sehen Sie sie sich an!
