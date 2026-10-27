---
title: "API-Referenz"
description: "Wails-Typen, -Methoden, -Optionen und CLI-Befehle nachschlagen"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Über diese Referenz

In diesem Abschnitt können Sie beim Arbeiten an einer Anwendung Wails-Typen, -Methoden, -Optionen und -Befehle nachschlagen. Wenn Sie eine schrittweise Einführung benötigen, beginnen Sie mit den [Tutorials](/tutorials/overview/). Anleitungen zum Ausführen einer Aufgabe finden Sie in den [Praxisanleitungen](/quick-start/next-steps/#features). Hintergrundinformationen zum Design finden Sie unter [Kernkonzepte](/concepts/architecture/).

**Gliederung:**

- [Anwendung](/reference/application/) – Zentrale Anwendungs-APIs
- [Fenster](/reference/window/) – Erstellen und Verwalten von Fenstern; unter [Fensteroptionen](/features/windows/options/) sind die Konfigurationsfelder aufgeführt
- [Menü](/reference/menu/) – APIs für Anwendungsmenüs, Kontextmenüs und den Infobereich; die [Menüelementreferenz](/features/menus/reference/) beschreibt Elementtypen und Eigenschaften
- [Ereignisse](/reference/events/) – Ereignissystem und integrierte Ereignisse
- [Dialoge](/reference/dialogs/) – Datei- und Meldungsdialoge
- [Frontend-Laufzeit](/reference/frontend-runtime/) – APIs der Frontend-Laufzeit
- [Übersicht der CLI-Befehle](/reference/cli/) – Nach Zweck gruppierte Befehle; [Flags und Beispiele](/guides/cli/) enthält Details zu den einzelnen Befehlen
- [Protokoll für Aktualisierungsmanifeste](/reference/update-manifest/) – Aktualisierungsmetadaten und Prüffelder
- [Mobile API](/guides/mobile/mobile-api/) – Experimentelle APIs für Mobilgeräte

## API-Konventionen

@details{title="Konventionen der Go-API – Für Entwickler, die mit Go noch nicht vertraut sind"}
### Benennung

- <strong>Typen</strong>: PascalCase (z. B. `WebviewWindow`)
- <strong>Methoden</strong>: PascalCase (z. B. `SetTitle()`)
- <strong>Optionen</strong>: Structs in PascalCase (z. B. `WindowOptions`)
- <strong>Konstanten</strong>: PascalCase (z. B. `WindowStartStateMaximised`)

#### Fehlerbehandlung

Die meisten Methoden, die fehlschlagen können, geben als letzten Rückgabewert `error` zurück. `app.Run()` blockiert, bis die Anwendung beendet wird, und gibt etwaige Startfehler zurück:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Das Erstellen eines Fensters gibt keinen Fehler zurück – `app.Window.New()` gibt direkt `*WebviewWindow` zurück.

#### Kontext

Methoden des Dienstlebenszyklus erhalten einen `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Der Kontext für die Lebensdauer der Anwendung ist über `app.Context()` verfügbar. Es gibt kein `RunWithContext` – rufen Sie `app.Run()` auf.

#### Optionsmuster

Für die Konfiguration werden Options-Structs verwendet:

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### Konventionen der JavaScript-API

#### Benennung

- **Funktionen**: camelCase (z. B. `setTitle()`)
- **Konstanten**: SCREAMING<em>SNAKE</em>CASE (z. B. `WINDOW_EVENT_FOCUS`)

#### Standardmäßig asynchron

Alle Aufrufe von Go-Methoden geben Promises zurück:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Fehlerbehandlung

Go-Fehler werden zu JavaScript-Ausnahmen:

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Typsicherheit

TypeScript-Definitionen werden automatisch generiert:

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Paketstruktur

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Importpfade

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Typenreferenz

### Allgemeine Typen

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Plattformunterschiede

Einige APIs verhalten sich je nach Plattform unterschiedlich:

| Funktion | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Anwendungsmenü** | Fenstermenüleiste | Globale Menüleiste | Fenstermenüleiste |
| **Infobereich** | Benachrichtigungsbereich | Menüleiste | Infobereich |
| **Dock** | Nicht verfügbar | ✅ Verfügbar | Nicht verfügbar |
| **Dateidialoge** | Nativ | Nativ | Nativ (GTK) |
| **Transparenz** | ✅ Vollständig | Erfordert [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Eingeschränkt |

Plattformspezifisches Verhalten ist im jeweiligen API-Abschnitt dokumentiert.

## Versionierung

Wails v3 verwendet semantische Versionierung:

- **Hauptversion** (v4.0.0): Inkompatible Änderungen
- **Nebenversion** (v3.1.0): Neue Funktionen, abwärtskompatibel
- **Patchversion** (v3.0.1): Fehlerbehebungen, abwärtskompatibel

**Wails v3.0.0:** Stabil

## Richtlinie zur Abkündigung

Wenn APIs als veraltet markiert werden:

1. **In der Dokumentation gekennzeichnet**, einschließlich eines Hinweises zur Abkündigung
2. **Alternative bereitgestellt**, einschließlich einer Migrationsanleitung
3. **Vor der Entfernung für 1 Hauptversion** weiter gepflegt
4. **Compilerwarnungen** (soweit möglich)

## API-Stabilität

### Stabile APIs ✅

Diese APIs sind stabil und können sicher in Produktionsumgebungen eingesetzt werden:

- Zentrale Anwendungs-APIs
- Fensterverwaltung
- Menüsystem
- Ereignissystem
- Dateidialoge
- Service-Bindings

### Experimentelle APIs ⚠️

Nur APIs, die ausdrücklich als experimentell dokumentiert sind, fallen nicht unter die stabile Kompatibilitätsgarantie für Desktop-Anwendungen:

- Android-Unterstützung
- iOS-Unterstützung
- Andere ausdrücklich als experimentell gekennzeichnete APIs

Experimentelle APIs sind in der Dokumentation entsprechend gekennzeichnet.

## Hilfe erhalten

### Fragen zur API

1. **Diese Referenz lesen** – Vollständige API-Dokumentation
2. **Beispiele ansehen** – [GitHub-Beispiele](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Discord durchsuchen** – [Discord-Server](https://discord.gg/JDdSxwjhGf)
4. **Die Community fragen** – Discord-Kanal #help

### API-Probleme melden

Einen Fehler oder eine Inkonsistenz gefunden?

1. **Vorhandene Issues prüfen** – [GitHub-Issues](https://github.com/wailsapp/wails/issues)
2. **Detaillierten Bericht erstellen** – Code, Fehlermeldung und Plattform angeben
3. **Reproduktionsbeispiel bereitstellen** – Ein minimales Beispiel, das das Problem veranschaulicht

## Weiterführende Dokumentation

- [Tutorials](/tutorials/overview/) – Durch das Erstellen echter Anwendungen lernen
- [Anleitungen](/guides/architecture/) – Aufgabenorientierte Anleitungen für häufige Szenarien
- [Funktionen](/features/windows/basics/) – Dokumentation der einzelnen Funktionen
- [Beispiele](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) – Funktionsfähige Codebeispiele auf GitHub

---

**API durchsuchen:** Verwenden Sie die Navigation auf der linken Seite, um die einzelnen APIs zu erkunden.
