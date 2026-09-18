---
title: "API-Referenz"
description: "Vollständige API-Dokumentation für Wails v3"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Über diese Referenz

Dies ist die vollständige API-Referenz für Wails v3. Sie dokumentiert alle öffentlichen Typen, Methoden und Optionen des Frameworks.

**Gliederung:**

- [Anwendung](/reference/application/) – Kern-APIs der Anwendung
- [Fenster](/reference/window/) – Erstellen und Verwalten von Fenstern
- [Menüs](/reference/menu/) – Anwendungs-, Kontext- und Infobereichsmenüs
- [Ereignisse](/reference/events/) – Ereignissystem und integrierte Ereignisse
- [Dialogfelder](/reference/dialogs/) – Datei- und Meldungsdialogfelder
- [Frontend-Laufzeit](/reference/frontend-runtime/) – APIs der Frontend-Laufzeit
- [CLI](/reference/cli/) – Befehlszeilenschnittstelle

## API-Konventionen

@details{title="Go-API-Konventionen – Für Entwickler, die neu in Go sind"}
### Benennung

- <strong></strong>Typen<strong></strong>: PascalCase (z. B. `WebviewWindow`)
- <strong></strong>Methoden<strong></strong>: PascalCase (z. B. `SetTitle()`)
- <strong></strong>Optionen<strong></strong>: Structs in PascalCase (z. B. `WindowOptions`)
- <strong></strong>Konstanten<strong></strong>: PascalCase (z. B. `WindowStartStateMaximised`)

#### Fehlerbehandlung

Die meisten Methoden, die fehlschlagen können, geben `error` als letzten Rückgabewert zurück. `app.Run()` blockiert, bis die Anwendung beendet wird, und gibt etwaige Startfehler zurück:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Das Erstellen eines Fensters gibt keinen Fehler zurück – `app.Window.New()` gibt direkt `*WebviewWindow` zurück.

#### Kontext

Methoden für den Lebenszyklus von Diensten erhalten einen `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Der Lebenszeitkontext der Anwendung ist über `app.Context()` verfügbar. `RunWithContext` ist nicht vorhanden – rufen Sie `app.Run()` auf.

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

## Typreferenz

### Häufig verwendete Typen

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
| **Dateidialogfelder** | Nativ | Nativ | Nativ (GTK) |
| **Transparenz** | ✅ Vollständig | Erfordert [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Eingeschränkt |

Plattformspezifisches Verhalten ist im jeweiligen API-Abschnitt dokumentiert.

## Versionierung

Wails v3 folgt der semantischen Versionierung:

- **Hauptversion** (v3.x.x): Inkompatible Änderungen
- **Nebenversion** (v3.x.x): Neue, abwärtskompatible Funktionen
- **Patchversion** (v3.x.x): Abwärtskompatible Fehlerbehebungen

**Aktueller Status:** Beta (API stabil, weitere Verbesserungen laufen)

## Richtlinie zur Kennzeichnung veralteter APIs

Wenn APIs als veraltet markiert werden:

1. **In der Dokumentation gekennzeichnet** mit einem Hinweis auf den veralteten Status
2. **Alternative bereitgestellt**, einschließlich einer Migrationsanleitung
3. **Wird vor der Entfernung noch 1 Hauptversion lang unterstützt**
4. **Compilerwarnungen** (soweit möglich)

## API-Stabilität

### Stabile APIs ✅

Diese APIs sind stabil und für den produktiven Einsatz geeignet:

- Zentrale Anwendungs-APIs
- Fensterverwaltung
- Menüsystem
- Ereignissystem
- Dateidialoge
- Service-Bindings

### Instabile APIs ⚠️

Diese APIs können sich vor der endgültigen Veröffentlichung noch ändern:

- Einige erweiterte Fensteroptionen
- Plattformspezifische Funktionen
- Experimentelle Funktionen

Instabile APIs sind in der Dokumentation gekennzeichnet.

## Hilfe erhalten

### Fragen zur API

1. **Diese Referenz lesen** – vollständige API-Dokumentation
2. **Beispiele ansehen** – [GitHub-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples)
3. **Discord durchsuchen** – [Discord-Server](https://discord.gg/JDdSxwjhGf)
4. **Die Community fragen** – Discord-Kanal #help

### API-Probleme melden

Einen Fehler oder eine Unstimmigkeit gefunden?

1. **Vorhandene Issues prüfen** – [GitHub-Issues](https://github.com/wailsapp/wails/issues)
2. **Detaillierten Bericht erstellen** – Code, Fehlermeldung und Plattform angeben
3. **Reproduktionsbeispiel bereitstellen** – minimales Beispiel, das das Problem veranschaulicht

## Weiterführende Dokumentation

- [Tutorials](/tutorials/overview/) – durch die Entwicklung realer Anwendungen lernen
- [Anleitungen](/guides/architecture/) – aufgabenorientierte Anleitungen für häufige Szenarien
- [Funktionen](/features/windows/basics/) – Dokumentation der einzelnen Funktionen
- [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) – ausführbare Codebeispiele auf GitHub

---

**API durchsuchen:** Verwenden Sie die Navigation auf der linken Seite, um bestimmte APIs zu erkunden.
