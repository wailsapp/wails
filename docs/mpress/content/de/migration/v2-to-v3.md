---
title: "Migration von v2 zu v3"
description: "Vollständige Anleitung zur Migration Ihrer Wails-Anwendung von v2 zu v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 wurde **vollständig neu geschrieben** und bietet erhebliche Verbesserungen bei Architektur, Leistung und Entwicklungskomfort. Diese Anleitung unterstützt Sie bei der Migration Ihrer v2-Anwendung zu v3.

**Wichtige Änderungen:**

- Neue Anwendungsstruktur
- Verbessertes Bindings-System
- Verbesserte Fensterverwaltung
- Verbessertes Ereignissystem
- Vereinfachte Konfiguration

**Migrationsdauer:** 1-4 Stunden für typische Anwendungen

## Inkompatible Änderungen

### Anwendungsinitialisierung

In v2 wurden die Anwendungseinrichtung, die Fensterkonfiguration und die Ausführung in einem einzigen Aufruf von `wails.Run()` zusammengefasst. Dieser monolithische Ansatz erschwerte es, mehrere Fenster zu erstellen, Fehler in verschiedenen Phasen zu behandeln oder einzelne Komponenten Ihrer Anwendung zu testen.

v3 unterteilt diese Aufgaben in separate Phasen: Anwendungserstellung, Fenstererstellung und Ausführung. Diese Trennung gibt Ihnen ausdrückliche Kontrolle über jede Phase des Anwendungslebenszyklus und macht den Code modularer und besser testbar.

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

**Vorteile:**

- **Unterstützung mehrerer Fenster**: Sie können Fenster jederzeit dynamisch erstellen, nicht nur beim Start
- **Bessere Fehlerbehandlung**: Jede Phase kann separat und mit angemessener Fehlerbehandlung validiert werden
- **Klarerer Code**: Durch die Trennung ist in jeder Phase eindeutig erkennbar, was geschieht
- **Bessere Testbarkeit**: Sie können die Anwendungseinrichtung testen, ohne die Ereignisschleife auszuführen
- **Mehr Flexibilität**: Fenster können während des gesamten Anwendungslebenszyklus erstellt, geschlossen und neu erstellt werden

### Bindings

In v2 benötigte jede gebundene Struktur ein Kontextfeld und eine `startup(ctx)`-Methode, um den Laufzeitkontext zu empfangen. Dadurch entstand eine enge Kopplung zwischen Ihrer Geschäftslogik und der Wails-Laufzeit, was das Testen und Verstehen des Codes erschwerte.

v3 führt das Service-Muster ein. Dabei sind Ihre Strukturen vollständig eigenständig und müssen keinen Laufzeitkontext speichern. Benötigt ein Service Zugriff auf die Anwendungsinstanz, erhält er diese ausdrücklich per Dependency Injection, statt dass der Kontext implizit weitergereicht wird.

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

**Vorteile:**

- **Keine impliziten Abhängigkeiten**: Services sind einfache Go-Strukturen ohne verborgene Laufzeitabhängigkeiten
- **Einfachere Tests**: Sie können Service-Methoden testen, ohne einen Wails-Kontext zu simulieren
- **Klarerer Code**: Abhängigkeiten sind explizit und werden als Konstruktorargumente übergeben, statt in einem Kontextfeld verborgen zu sein
- **Bessere Organisation**: Services können nach Domäne gruppiert werden, statt alle in einer einzigen `App`-Struktur enthalten zu sein
- **Ordnungsgemäße Initialisierung**: Verwenden Sie bei erforderlicher Initialisierung ausdrücklich die Methode `ServiceStartup()`

### Laufzeit

In v2 musste für alle Laufzeitoperationen ein Kontext an globale Funktionen aus dem Paket `runtime` übergeben werden. Dadurch war der Code in der gesamten Codebasis eng an das Kontextobjekt gekoppelt, und die API wirkte eher prozedural als objektorientiert.

v3 ersetzt die kontextbasierte Laufzeit durch direkte Methodenaufrufe für Anwendungs- und Fensterobjekte. Operationen werden direkt für die betroffenen Objekte aufgerufen, wodurch der Code intuitiver und objektorientierter wird.

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

**Vorteile:**

- **Objektorientierter Entwurf**: Methoden werden für die Objekte aufgerufen, auf die sie wirken (Fenster, Anwendung, Menü usw.)
- **Klarere Absicht**: `window.SetTitle()` ist eindeutiger als `runtime.WindowSetTitle(ctx, ...)`
- **Bessere IDE-Unterstützung**: Die automatische Vervollständigung funktioniert korrekt, wenn Methoden zu Objekten gehören
- **Mehr Übersicht bei mehreren Fenstern**: Bei mehreren Fenstern wählen Sie ausdrücklich aus, auf welchem Fenster die Operation ausgeführt werden soll
- **Keine Kontextweitergabe**: Sie müssen den Kontext nicht durch jede Funktion reichen

### Frontend-Bindings

In v2 waren Bindings nach Go-Paket und Strukturnamen organisiert, was typischerweise zu Pfaden wie `wailsjs/go/main/App` führte. Diese Struktur bildete die logische Gruppierung nicht ab und erschwerte das Auffinden zusammengehöriger Funktionen.

v3 organisiert Bindings nach Service-Namen und Anwendungsmodul und schafft damit eine klarere logische Struktur. Die Bindings werden in einem Verzeichnis `bindings` generiert und dort nach dem Namen Ihrer Anwendung und den Service-Namen gegliedert. So ist leichter erkennbar, welche Funktionen verfügbar sind.

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

**Vorteile:**

- **Logische Organisation**: Bindings werden nach Service-Namen statt nach der Go-Paketstruktur gruppiert
- **Klarere Importe**: Der Pfad bildet die Domänenlogik (greetservice) statt der Dateistruktur (main/App) ab
- **Bessere Auffindbarkeit**: Sie können Bindings nach Funktion statt nach technischer Struktur durchsuchen
- **Einheitliche Benennung**: Die servicebasierte Organisation entspricht Ihrer Backend-Architektur
- **Einfachere Pfade**: Kein Präfix `../wailsjs/go` mehr, sondern nur noch `./bindings`

### Ereignisse

In v2 verwendeten Ereignisse variadische `interface{}`-Parameter, und jeder Ereignisfunktion musste ein Kontext übergeben werden. Ereignishandler empfingen untypisierte Daten, für die manuelle Typzusicherungen erforderlich waren. Dadurch war das Ereignissystem fehleranfällig und schwer zu debuggen.

v3 führt typisierte Ereignisobjekte ein und hebt die Kontextanforderung auf. Ereignishandler empfangen ein ordnungsgemäßes Ereignisobjekt mit typisierten Daten, wodurch das Ereignissystem zuverlässiger und einfacher zu verwenden ist.

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

**Vorteile:**

- **Typsicherheit**: Ereignisse verwenden geeignete Ereignisobjekte anstelle von `...interface{}`
- **Einfacheres Debugging**: Ereignisobjekte enthalten Metadaten wie den Ereignisnamen und erleichtern dadurch das Debugging
- **Übersichtlichere API**: `app.Event.On()` und `app.Event.Emit()` sind intuitiver als Runtime-Funktionen
- **Kein Kontext erforderlich**: Ereignisse funktionieren direkt auf dem App-Objekt, ohne dass Kontext weitergereicht werden muss
- **Einfachere Handler**: Ereignis-Handler haben eine klare Signatur statt variadischer Parameter

### Fenster

v2 unterstützte nur ein einziges Fenster pro Anwendung. Das Fenster wurde beim Start erstellt, und alle Fensteroperationen erfolgten über Runtime-Funktionen, die sich implizit auf dieses einzelne Fenster bezogen.

v3 führt native Unterstützung für mehrere Fenster als Kernfunktion ein. Jedes Fenster ist ein eigenständiges Objekt mit eigenen Methoden und einem eigenen Lebenszyklus. Während der gesamten Laufzeit deiner Anwendung kannst du mehrere Fenster dynamisch erstellen, verwalten und schließen.

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

**Warum dies besser ist:**

- **Anwendungen mit mehreren Fenstern**: Erstelle Apps mit mehreren unabhängigen Fenstern (Dashboards, Einstellungen, Werkzeuge usw.)
- **Explizite Fensterreferenzen**: Jedes Fenster ist ein Objekt, das du speichern und direkt bearbeiten kannst
- **Dynamische Fenstererstellung**: Erstelle und schließe Fenster jederzeit während der Laufzeit
- **Unabhängiger Fensterzustand**: Jedes Fenster hat eigene Ereignisse, Eigenschaften und einen eigenen Lebenszyklus
- **Bessere Architektur**: Die Fensterverwaltung ist objektorientiert statt kontextbasiert

## Migrationsschritte

### Schritt 1: Abhängigkeiten aktualisieren

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**Aktualisieren:**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### Schritt 2: main.go aktualisieren

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

### Schritt 3: App-Struktur in einen Service umwandeln

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

### Schritt 4: Runtime-Aufrufe aktualisieren

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

### Schritt 5: Frontend aktualisieren

**Neue Bindings generieren:**

```bash
wails3 generate bindings
```

**Importe aktualisieren:**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**Ereignisbehandlung aktualisieren:**

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

### Schritt 6: Konfiguration aktualisieren

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

## Funktionszuordnung

### Dialoge

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

### Menüs

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

### Infobereich

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

## Häufige Probleme

### Problem: Bindings nicht gefunden

**Problem:** Importfehler nach der Migration

**Lösung:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### Problem: Kontextfehler

**Problem:** `ctx` ist nicht verfügbar

**Lösung:**

Speichere stattdessen die App-Referenz:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Problem: Fenstermethoden funktionieren nicht

**Problem:** `runtime.WindowSetTitle()` ist nicht vorhanden

**Lösung:**

Verwende Fenstermethoden direkt:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Problem: Ereignisse werden nicht ausgelöst

**Problem:** Ereignisse sind registriert, werden aber nicht empfangen

**Lösung:**

Prüfe, ob die Ereignisnamen exakt übereinstimmen:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Migration testen

### Checkliste

- [ ] Anwendung startet ohne Fehler
- [ ] Alle Bindings funktionieren
- [ ] Ereignisse werden gesendet und empfangen
- [ ] Fenster werden ordnungsgemäß geöffnet und geschlossen
- [ ] Menüs funktionieren (falls zutreffend)
- [ ] Dialoge funktionieren (falls zutreffend)
- [ ] Der Infobereich funktioniert (falls zutreffend)
- [ ] Der Build-Prozess funktioniert
- [ ] Der Produktions-Build funktioniert

### Testbefehle

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Vorteile von v3

### Leistung

- **Schnellerer Start** – Optimierte Initialisierung
- **Geringerer Speicherbedarf** – Effiziente Ressourcennutzung
- **Bessere Bridge** – &lt;1 ms Aufruf-Overhead

### Funktionen

- **Mehrere Fenster** – Native Unterstützung
- **Infobereich** – Integriert
- **Bessere Ereignisse** – typisierte, einfachere API
- **Dienste** – bessere Codeorganisation

### Entwicklungserlebnis

- **Typsicherheit** – vollständige TypeScript-Unterstützung
- **Bessere Fehlermeldungen** – klare Fehlermeldungen
- **Hot Reload** – schnellere Entwicklung
- **Bessere Dokumentation** – umfassende Anleitungen

## Hilfe erhalten

### Ressourcen

- [Dokumentation](/quick-start/why-wails/)
- [Discord-Community](https://discord.gg/JDdSxwjhGf)
- [GitHub-Issues](https://github.com/wailsapp/wails/issues)
- [Beispiele](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### Häufige Fragen

**F: Kann ich v2 und v3 parallel verwenden?** A: Ja, sie verwenden unterschiedliche Importpfade.

**F: Ist v3 produktionsreif?** Wails v3.0.0 ist die stabile Version für Desktopanwendungen unter Windows, macOS und Linux. Die Unterstützung für Android und iOS bleibt experimentell und fällt nicht unter die stabile Kompatibilitätszusage für Desktopanwendungen.

**F: Wird v2 weiterhin gepflegt?** A: Ja, v2 erhält kritische Updates.

**F: Wie lange dauert die Migration?** A: Bei typischen Anwendungen 1-4 Stunden.

## Nächste Schritte

@cards{cols="2"}
🚀 Schnellstart
Erste Schritte mit Wails v3.

[Mehr erfahren →](/quick-start/installation/)

---
★ Kernkonzepte
Die Architektur von v3 verstehen.

[Mehr erfahren →](/concepts/architecture/)

---
◆ Bindings
Das neue Bindingsystem kennenlernen.

[Mehr erfahren →](/features/bindings/methods/)

---
📖 Beispiele
Vollständige Beispiele für v3 ansehen.

[Beispiele ansehen →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

**Noch Fragen?** Frage auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder [öffne ein Issue](https://github.com/wailsapp/wails/issues).
