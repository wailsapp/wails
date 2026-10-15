---
title: "Migration von v2 zu v3"
description: "Vollständige Anleitung zur Migration Ihrer Wails-v2-Anwendung auf v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 wurde **vollständig neu entwickelt** und bietet erhebliche Verbesserungen bei Architektur, Leistung und Entwicklerfreundlichkeit. Diese Anleitung unterstützt Sie bei der Migration Ihrer v2-Anwendung auf v3.

**Wichtige Änderungen:**

- Neue Anwendungsstruktur
- Verbessertes Bindings-System
- Verbesserte Fensterverwaltung
- Verbessertes Ereignissystem
- Vereinfachte Konfiguration

**Migrationsdauer:** 1-4 Stunden für typische Anwendungen

## Inkompatible Änderungen

### Anwendungsinitialisierung

In v2 waren Anwendungseinrichtung, Fensterkonfiguration und Ausführung in einem einzigen `wails.Run()`-Aufruf zusammengefasst. Dieser monolithische Ansatz erschwerte es, mehrere Fenster zu erstellen, Fehler in verschiedenen Phasen zu behandeln oder einzelne Komponenten Ihrer Anwendung zu testen.

v3 trennt diese Aufgaben in verschiedene Phasen: Anwendungserstellung, Fenstererstellung und Ausführung. Dadurch haben Sie ausdrückliche Kontrolle über jede Phase des Lebenszyklus Ihrer Anwendung, und der Code wird modularer und besser testbar.

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

**Vorteile dieses Ansatzes:**

- **Unterstützung mehrerer Fenster**: Sie können Fenster jederzeit dynamisch erstellen, nicht nur beim Start
- **Bessere Fehlerbehandlung**: Jede Phase kann separat und mit geeigneter Fehlerbehandlung validiert werden
- **Übersichtlicherer Code**: Durch die Trennung ist klar erkennbar, was in jeder Phase geschieht
- **Bessere Testbarkeit**: Sie können die Anwendungseinrichtung testen, ohne die Ereignisschleife auszuführen
- **Mehr Flexibilität**: Fenster können während des gesamten Lebenszyklus der Anwendung erstellt, zerstört und neu erstellt werden

### Bindings

In v2 benötigte jede gebundene Struktur ein Kontextfeld und eine `startup(ctx)`-Methode, um den Laufzeitkontext zu empfangen. Dadurch entstand eine enge Kopplung zwischen Ihrer Geschäftslogik und der Wails-Laufzeit, wodurch der Code schwerer zu testen und zu verstehen war.

v3 führt das Service-Muster ein, bei dem Ihre Strukturen vollständig eigenständig sind und keinen Laufzeitkontext speichern müssen. Benötigt ein Service Zugriff auf die Anwendungsinstanz, erhält er diese ausdrücklich per Dependency Injection statt durch implizite Weitergabe des Kontexts.

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

**Vorteile dieses Ansatzes:**

- **Keine impliziten Abhängigkeiten**: Services sind einfache Go-Strukturen ohne verborgene Laufzeitabhängigkeiten
- **Einfachere Tests**: Sie können Service-Methoden testen, ohne einen Wails-Kontext zu simulieren
- **Übersichtlicherer Code**: Abhängigkeiten sind ausdrücklich angegeben (als Konstruktorargumente) und nicht in einem Kontextfeld verborgen
- **Bessere Organisation**: Services können nach Domäne gruppiert werden, statt sich alle in einer einzigen `App`-Struktur zu befinden
- **Explizite Initialisierung**: Verwenden Sie bei Bedarf die Methode `ServiceStartup()` zur Initialisierung, damit dieser Schritt ausdrücklich erkennbar ist

### Laufzeit

In v2 musste für alle Laufzeitoperationen ein Kontext an globale Funktionen aus dem Paket `runtime` übergeben werden. Dadurch entstand in der gesamten Codebasis eine enge Kopplung an das Kontextobjekt, und die API wirkte eher prozedural als objektorientiert.

v3 ersetzt die kontextbasierte Laufzeit durch direkte Methodenaufrufe auf Anwendungs- und Fensterobjekten. Operationen werden direkt auf den Objekten aufgerufen, die sie betreffen. Dadurch wird der Code intuitiver und objektorientierter.

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

**Vorteile dieses Ansatzes:**

- **Objektorientierter Entwurf**: Methoden werden auf den Objekten aufgerufen, die sie betreffen (Fenster, Anwendung, Menü usw.)
- **Klarere Absicht**: `window.SetTitle()` ist aussagekräftiger als `runtime.WindowSetTitle(ctx, ...)`
- **Bessere IDE-Unterstützung**: Die automatische Vervollständigung funktioniert korrekt, wenn Methoden Objekten zugeordnet sind
- **Mehr Übersicht bei mehreren Fenstern**: Bei mehreren Fenstern wählen Sie ausdrücklich aus, auf welchem Fenster die Operation ausgeführt werden soll
- **Keine Kontextweitergabe**: Sie müssen den Kontext nicht durch jede Funktion weiterreichen

### Frontend-Bindings

In v2 waren Bindings nach Go-Paket und Strukturnamen organisiert, wodurch typischerweise Pfade wie `wailsjs/go/main/App` entstanden. Diese Struktur spiegelte keine logische Gruppierung wider und erschwerte es, zusammengehörige Funktionen zu finden.

v3 organisiert Bindings nach Service-Namen und Anwendungsmodul und schafft so eine übersichtlichere logische Struktur. Die Bindings werden in einem Verzeichnis `bindings` generiert und dort nach dem Namen Ihrer Anwendung und den Service-Namen organisiert. Dadurch lässt sich die verfügbare Funktionalität leichter überblicken.

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

**Vorteile dieses Ansatzes:**

- **Logische Organisation**: Bindings werden nach Service-Namen statt nach der Go-Paketstruktur gruppiert
- **Übersichtlichere Importe**: Der Pfad bildet die Domänenlogik (greetservice) statt der Dateistruktur (main/App) ab
- **Leichter auffindbar**: Bindings lassen sich nach Funktion statt nach technischer Struktur durchsuchen
- **Einheitliche Benennung**: Die servicebasierte Organisation entspricht Ihrer Backend-Architektur
- **Einfachere Pfade**: Kein Präfix `../wailsjs/go` mehr – nur noch `./bindings`

### Ereignisse

In v2 verwendeten Ereignisse variadische `interface{}`-Parameter, und jeder Ereignisfunktion musste ein Kontext übergeben werden. Ereignishandler erhielten untypisierte Daten, für die manuelle Typzusicherungen erforderlich waren. Dadurch war das Ereignissystem fehleranfällig und schwer zu debuggen.

v3 führt typisierte Ereignisobjekte ein und hebt die Kontextanforderung auf. Ereignishandler erhalten ein korrektes Ereignisobjekt mit typisierten Daten. Dadurch ist das Ereignissystem zuverlässiger und einfacher zu verwenden.

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

**Warum dies besser ist:**

- **Typsicherheit**: Ereignisse verwenden korrekte Ereignisobjekte statt `...interface{}`
- **Einfacheres Debugging**: Ereignisobjekte enthalten Metadaten wie den Ereignisnamen, was das Debugging erleichtert
- **Übersichtlichere API**: `app.Event.On()` und `app.Event.Emit()` sind intuitiver als Runtime-Funktionen
- **Kein Kontext erforderlich**: Ereignisse funktionieren direkt mit dem App-Objekt, ohne den Kontext weiterreichen zu müssen
- **Einfachere Handler**: Ereignishandler haben statt variadischer Parameter eine klare Signatur

### Fenster

v2 unterstützte nur ein einziges Fenster pro Anwendung. Das Fenster wurde beim Start erstellt, und alle Fensteroperationen erfolgten über Runtime-Funktionen, die implizit auf dieses eine Fenster abzielten.

v3 führt native Mehrfensterunterstützung als Kernfunktion ein. Jedes Fenster ist ein eigenständiges Objekt mit eigenen Methoden und einem eigenen Lebenszyklus. Während der gesamten Laufzeit Ihrer Anwendung können Sie mehrere Fenster dynamisch erstellen, verwalten und schließen.

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

- **Mehrfensteranwendungen**: Erstellen Sie Apps mit mehreren unabhängigen Fenstern (Dashboards, Einstellungen, Werkzeuge usw.)
- **Explizite Fensterreferenzen**: Jedes Fenster ist ein Objekt, das Sie speichern und direkt bearbeiten können
- **Dynamische Fenstererstellung**: Erstellen und zerstören Sie Fenster jederzeit während der Laufzeit
- **Unabhängiger Fensterzustand**: Jedes Fenster hat eigene Ereignisse, Eigenschaften und einen eigenen Lebenszyklus
- **Bessere Architektur**: Die Fensterverwaltung ist objektorientiert statt kontextbasiert

## Migrationsschritte

### Schritt 1: Abhängigkeiten aktualisieren

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**Aktualisieren:**

```bash
go get github.com/wailsapp/wails/v3@latest
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

### Taskleistenbereich

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

Speichern Sie stattdessen eine Referenz auf die App:

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

Verwenden Sie die Fenstermethoden direkt:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Problem: Ereignisse werden nicht ausgelöst

**Problem:** Ereignisse sind registriert, werden aber nicht empfangen

**Lösung:**

Prüfen Sie, ob die Ereignisnamen exakt übereinstimmen:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Migration testen

### Checkliste

- [ ] Die Anwendung startet fehlerfrei
- [ ] Alle Bindings funktionieren
- [ ] Ereignisse werden gesendet und empfangen
- [ ] Fenster werden korrekt geöffnet und geschlossen
- [ ] Menüs funktionieren (falls zutreffend)
- [ ] Dialoge funktionieren (falls zutreffend)
- [ ] Der Systembereich funktioniert (falls zutreffend)
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

- **Schnellerer Start** – optimierte Initialisierung
- **Geringerer Speicherverbrauch** – effiziente Ressourcennutzung
- **Bessere Bridge** – &lt;1 ms Aufruf-Overhead

### Funktionen

- **Mehrere Fenster** – native Unterstützung
- **Systembereich** – integriert
- **Bessere Ereignisse** – typisierte, einfachere API
- **Dienste** – bessere Codeorganisation

### Entwicklerfreundlichkeit

- **Typsicherheit** – vollständige TypeScript-Unterstützung
- **Bessere Fehlermeldungen** – eindeutige Fehlermeldungen
- **Hot Reload** – schnellere Entwicklung
- **Bessere Dokumentation** – umfassende Anleitungen

## Hilfe erhalten

### Ressourcen

- [Dokumentation](/quick-start/why-wails/)
- [Discord-Community](https://discord.gg/JDdSxwjhGf)
- [GitHub-Issues](https://github.com/wailsapp/wails/issues)
- [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples)

### Häufige Fragen

**F: Kann ich v2 und v3 parallel ausführen?** A: Ja, sie verwenden unterschiedliche Importpfade.

**F: Ist v3 produktionsreif?** A: v3 ist eine Beta-Software mit einer stabilen Desktop-API. Anwendungen werden damit bereits produktiv betrieben. Testen Sie sie jedoch vor der Bereitstellung gründlich. v2 bleibt die derzeitige stabile Version.

**F: Wird v2 weiterhin gepflegt?** A: Ja, v2 wird weiterhin kritische Updates erhalten.

**F: Wie lange dauert die Migration?** A: Bei typischen Anwendungen 1-4 Stunden.

## Nächste Schritte

@cards{cols="2"}
🚀 Schnellstart
Beginnen Sie mit Wails v3.

[Mehr erfahren →](/quick-start/installation/)

---
★ Kernkonzepte
Machen Sie sich mit der Architektur von v3 vertraut.

[Mehr erfahren →](/concepts/architecture/)

---
◆ Bindings
Lernen Sie das neue Binding-System kennen.

[Mehr erfahren →](/features/bindings/methods/)

---
📖 Beispiele
Sehen Sie sich vollständige Beispiele für v3 an.

[Beispiele ansehen →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder [eröffnen Sie ein Issue](https://github.com/wailsapp/wails/issues).
