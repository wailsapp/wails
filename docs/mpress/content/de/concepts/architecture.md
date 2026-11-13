---
title: "Funktionsweise von Wails"
description: "Die Architektur von Wails und ihre native Performance verstehen"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails ist ein Framework zum Erstellen von Desktopanwendungen mit **Go für das Backend** und **Webtechnologien für das Frontend**. Anders als Electron bündelt Wails jedoch keinen Browser, sondern verwendet die **native WebView des Betriebssystems**.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Wails-App

  frontend: Frontend
  backend: Go-Backend
  os: Betriebssystem

  Initialisation: Initialisierung {
    shape: sequence_diagram
    backend."Serves Static Web App": Stellt statische Web-App bereit
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": Rendert die Website über die native WebView des Betriebssystems
  }
  Regular Communication: Reguläre Kommunikation {
    shape: sequence_diagram
    frontend."Make API-style call": API-ähnlichen Aufruf ausführen
    frontend -> backend.a: JSON
    backend.a."Service processes request": Dienst verarbeitet die Anfrage
    backend.a -> os: System-APIs aufrufen
    backend.a."Generate Response": Antwort erzeugen
    backend.a -> frontend: JSON
    frontend."Process response": Antwort verarbeiten
  }
  backend.a.label: a
}
```

**Wesentliche Unterschiede zu Electron:**

| Aspekt | Wails | Electron |
| --- | --- | --- |
| **Browser** | Vom Betriebssystem bereitgestellte WebView | Gebündeltes Chromium (~100 MB) |
| **Backend** | Go (kompiliert) | Node.js (interpretiert) |
| **Kommunikation** | In-Memory-Bridge | IPC (prozessübergreifend) |
| **Paketgröße** | ~15 MB | ~150 MB |
| **Arbeitsspeicher** | ~10 MB | ~100 MB+ |
| **Startzeit** | &lt;0.5 s | 2-3 s |

## Kernkomponenten

### 1. Native WebView

Wails verwendet die integrierte Web-Rendering-Engine des Betriebssystems:

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Basiert auf Chromium (wie der Edge-Browser)
- Unter Windows 10/11 vorinstalliert
- Automatische Updates über Windows Update
- Vollständige Unterstützung moderner Webstandards

[macOS]
**WebKit** (Rendering-Engine von Safari)

- In macOS integriert
- Dieselbe Engine wie der Safari-Browser
- Hervorragende Performance und Akkulaufzeit
- Vollständige Unterstützung moderner Webstandards

[Linux]
**WebKitGTK** (GTK-Portierung von WebKit)

- Installation über den Paketmanager
- Dieselbe Engine wie GNOME Web (Epiphany)
- Gute Unterstützung von Webstandards
- Ressourcenschonend und performant

@end

**Warum das wichtig ist:**

- **Kein gebündelter Browser** → Kleinere App-Größe
- **Betriebssystemnativ** → Bessere Integration und Performance
- **Automatische Updates** → Sicherheitspatches durch Betriebssystemupdates
- **Vertrautes Rendering** → Wie im Systembrowser

### 2. Die Wails Bridge

Die Bridge ist das Herzstück von Wails. Sie ermöglicht die **direkte Kommunikation** zwischen Go und JavaScript.

```d2
direction: down

Frontend: Frontend (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Wails Bridge {
  Encoder: JSON-Encoder {
    shape: rectangle
  }

  Router: Methoden-Router {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON-Decoder {
    shape: rectangle
  }
}

Backend: Backend (Go) {
  Services: Registrierte Dienste {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Go-Methode aufrufen\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Als JSON codieren\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. An Dienst weiterleiten\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Ergebnis zurückgeben\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Für JS decodieren\nPromise wird aufgelöst"
```

**Funktionsweise:**

1. **Das Frontend ruft eine Go-Methode auf** (über ein automatisch generiertes Binding)
2. **Die Bridge codiert den Aufruf** als JSON (Methodenname + Argumente)
3. **Der Router findet die Go-Methode** in den registrierten Diensten
4. **Die Go-Methode wird ausgeführt** und gibt einen Wert zurück
5. **Die Bridge decodiert das Ergebnis** und sendet es an das Frontend zurück
6. **Das Promise wird aufgelöst** und stellt das Ergebnis in JavaScript bereit

**Performance-Eigenschaften:**

- **Im Arbeitsspeicher**: Kein Netzwerk-Overhead, kein HTTP
- **Zero-Copy**, sofern möglich (bei großen Datenmengen)
- **Standardmäßig asynchron**: Auf beiden Seiten nicht blockierend
- **Typsicher**: TypeScript-Definitionen werden automatisch generiert

### 3. Dienstsystem

Dienste sind die empfohlene Methode, um Go-Funktionalität für das Frontend bereitzustellen.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Service-Erkennung:**

- Wails **analysiert beim Start Ihre Struktur**
- **Exportierte Methoden** können vom Frontend aufgerufen werden
- **Typinformationen** werden für TypeScript-Bindings extrahiert
- **Die Fehlerbehandlung** erfolgt automatisch (Go-Fehler → JS-Ausnahmen)

**Generiertes TypeScript-Binding:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Warum Services?**

- **Typsicher**: vollständige TypeScript-Unterstützung
- **Automatische Erkennung**: keine manuelle Registrierung von Methoden
- **Strukturiert**: zusammengehörige Funktionen gruppieren
- **Testbar**: Services sind einfache Go-Strukturen

[Mehr über Services erfahren →](/features/bindings/services/)

### 4. Ereignissystem

Ereignisse ermöglichen die **Pub/Sub-Kommunikation** zwischen Komponenten.

```d2
direction: left

Wails Event System: Wails-Ereignissystem {
  shape: sequence_diagram

  window1: Fenster 1
  window2: Fenster 2
  backend: Go-Backend

  Event Driver: Ereignistreiber {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "'data-updated'-Ereignisse abonnieren"
    window2."Subscribe to 'data-updated' events": "'data-updated'-Ereignisse abonnieren"
    backend.a."App Emit('data-updated', data)": "App Emit('data-updated', data)"
    backend.a -> window1.a: JSON-Ereignisbus
    backend.a -> window2: JSON-Ereignisbus
    window1.a."Subscriber processes On('data-updated', handler)": "Abonnent verarbeitet On('data-updated', handler)"
    window2."Subscriber processes On('data-updated', handler)": "Abonnent verarbeitet On('data-updated', handler)"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**Anwendungsfälle:**

- **Kommunikation zwischen Fenstern**: Ein Fenster benachrichtigt andere
- **Hintergrundaufgaben**: Ein Go-Service informiert die Benutzeroberfläche über den Fortschritt
- **Zustandssynchronisierung**: Mehrere Fenster synchron halten
- **Lose Kopplung**: Komponenten benötigen keine direkten Referenzen

**Beispiel:**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[Mehr über Ereignisse erfahren →](/features/events/system/)

## Anwendungslebenszyklus

Wenn Sie den Lebenszyklus verstehen, wissen Sie, wann Sie Ressourcen initialisieren und bereinigen müssen.

```d2
direction: down

Start: Anwendungsstart {
  shape: oval
  style.fill: "#10B981"
}

Init: Initialisierung {
  Create: Anwendung erstellen {
    shape: rectangle
  }

  Register: Services registrieren {
    shape: rectangle
  }

  Setup: Fenster/Menüs einrichten {
    shape: rectangle
  }
}

Run: Ereignisschleife {
  Events: Ereignisse verarbeiten {
    shape: rectangle
  }

  Messages: Nachrichten verarbeiten {
    shape: rectangle
  }

  Render: Benutzeroberfläche aktualisieren {
    shape: rectangle
  }
}

Shutdown: Beenden {
  Cleanup: Ressourcen bereinigen {
    shape: rectangle
  }

  Save: Zustand speichern {
    shape: rectangle
  }
}

End: Anwendungsende {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: Schleife
Run.Events -> Shutdown.Cleanup: Beendigungssignal
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Lebenszyklus-Hooks:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options` besitzt kein Feld `OnStartup`. Startaufgaben gehören in `ServiceStartup(ctx, options)` eines Services, in einen über `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)` registrierten Callback oder einfach vor `app.Run()`.

[Mehr über den Lebenszyklus erfahren →](/concepts/lifecycle/)

## Build-Prozess

So erstellt Wails Ihre Anwendung:

```d2
direction: down

Source: Quellcode {
  Go: "Go-Code\n(main.go, Services)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "Frontend-Code\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: Build-Prozess {
  AnalyseGo: Go-Code analysieren {
    shape: rectangle
  }

  GenerateBindings: Bindings generieren {
    shape: rectangle
  }

  BuildFrontend: Frontend erstellen {
    shape: rectangle
  }

  CompileGo: Go kompilieren {
    shape: rectangle
  }

  Embed: Assets einbetten {
    shape: rectangle
  }
}

Output: Ausgabe {
  Binary: "Native Binärdatei\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Typen extrahieren
Build.GenerateBindings -> Source.Frontend: TypeScript-Bindings
Source.Frontend -> Build.BuildFrontend: Kompilieren (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Gebündelte Assets
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Build-Schritte:**

1. **Go-Code analysieren**
  - Services nach exportierten Methoden durchsuchen
  - Parameter- und Rückgabetypen extrahieren
  - Methodensignaturen generieren


2. **TypeScript-Bindings generieren**
  - Für jeden Service `.ts`-Dateien erstellen
  - Vollständige Typdefinitionen einfügen
  - JSDoc-Kommentare hinzufügen


3. **Frontend erstellen**
  - Ihren Bundler ausführen (Vite, webpack usw.)
  - Minifizieren und optimieren
  - Nach `frontend/dist/` ausgeben


4. **Go kompilieren**
  - Mit Optimierungen kompilieren (`-ldflags="-s -w"`)
  - Build-Metadaten einfügen
  - Plattformspezifisch kompilieren


5. **Assets einbetten**
  - Frontend-Dateien in die Go-Binärdatei einbetten
  - Assets komprimieren
  - Eine einzelne ausführbare Datei erstellen


**Ergebnis:** Eine einzelne native ausführbare Datei, in die alles eingebettet ist.

[Mehr über das Erstellen erfahren →](/guides/build/building/)

## Entwicklung und Produktion

Wails verhält sich in der Entwicklungs- und Produktionsumgebung unterschiedlich:

@tabs{sync-key="mode"}
[Entwicklung (wails3 dev)]
**Merkmale:**

- **Hot Reload**: Frontend-Änderungen werden sofort neu geladen
- **Source Maps**: Debugging mit dem ursprünglichen Quellcode
- **DevTools**: Browser-DevTools sind verfügbar
- **Protokollierung**: Ausführliche Protokollierung ist aktiviert
- **Externes Frontend**: Wird vom Entwicklungsserver (Vite) bereitgestellt

**Funktionsweise:**

```d2
direction: right

WailsApp: Wails-App {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Vite-Entwicklungsserver\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: Anfragen weiterleiten
DevServer -> WebView: Mit HMR bereitstellen
WebView -> WailsApp: Go-Methoden aufrufen
```

**Vorteile:**

- Sofortiges Feedback zu Änderungen
- Umfassende Debugging-Möglichkeiten
- Schnellere Iteration

[Produktion (wails3 build)]
**Merkmale:**

- **Eingebettete Assets**: Das Frontend ist in die Binärdatei integriert
- **Optimiert**: Minifiziert und komprimiert
- **Keine DevTools**: Standardmäßig deaktiviert
- **Minimale Protokollierung**: Nur Fehler
- **Eine einzige Datei**: Alles befindet sich in einer ausführbaren Datei

**Funktionsweise:**

```d2
direction: right

Binary: "Eine einzige Binärdatei\n(myapp.exe)" {
  GoCode: Kompiliertes Go {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "Eingebettete Assets\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: Aus dem Arbeitsspeicher bereitstellen
WebView -> Binary.GoCode: Go-Methoden aufrufen
```

**Vorteile:**

- Verteilung als einzelne Datei
- Geringere Größe (minifiziert)
- Bessere Leistung
- Keine externen Abhängigkeiten

@end

## Speichermodell

Wenn du die Speichernutzung verstehst, kannst du effiziente Anwendungen entwickeln.

**Speicherbereiche:**

1. **Go-Heap**
  - Deine Services und dein Anwendungszustand
  - Wird vom Go-Garbage-Collector verwaltet
  - Bei einfachen Apps typischerweise 5-10MB


2. **WebView-Speicher**
  - DOM, JavaScript-Heap, CSS
  - Wird von der WebView-Engine verwaltet
  - Bei einfachen Apps typischerweise 10-20MB


3. **Bridge-Speicher**
  - Nachrichtenpuffer für die Kommunikation
  - Minimaler Overhead (<1MB)
  - Für große Daten nach Möglichkeit ohne Kopieren


**Optimierungstipps:**

- **Große Datenübertragungen vermeiden**: IDs übergeben und Details bei Bedarf abrufen
- **Events für Aktualisierungen verwenden**: Nicht vom Frontend aus pollen
- **Große Dateien streamen**: Nicht vollständig in den Arbeitsspeicher laden
- **Listener bereinigen**: Event-Listener entfernen, wenn sie nicht mehr benötigt werden

[Mehr über Leistung erfahren →](/guides/performance/)

## Sicherheitsmodell

Wails bietet eine standardmäßig sichere Architektur:

```d2
direction: down

Frontend: Frontend (nicht vertrauenswürdig) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails Bridge (Validierung) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Backend (vertrauenswürdig) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Methode aufrufen
Bridge -> Bridge: "Validieren:\n- Ist die Methode vorhanden?\n- Sind die Typen korrekt?\n- Ist der Zugriff erlaubt?"
Bridge -> Backend: Bei Gültigkeit ausführen
Backend -> Bridge: Ergebnis zurückgeben
Bridge -> Frontend: Antwort senden
```

**Sicherheitsfunktionen:**

1. **Positivliste für Methoden**
  - Nur exportierte Methoden können aufgerufen werden
  - Auf private Methoden kann nicht zugegriffen werden
  - Eine explizite Service-Registrierung ist erforderlich


2. **Typvalidierung**
  - Argumente werden anhand der Go-Typen geprüft
  - Ungültige Typen werden abgelehnt
  - Verhindert Injection-Angriffe


3. **Kein eval()**
  - Das Frontend kann keinen beliebigen Go-Code ausführen
  - Nur vordefinierte Methoden können aufgerufen werden
  - Keine dynamische Codeausführung


4. **Kontextisolierung**
  - Jedes Fenster hat einen eigenen Kontext
  - Services können den Kontext des Aufrufers prüfen
  - Berechtigungen pro Fenster sind möglich


**Bewährte Verfahren:**

- **Benutzereingaben validieren** – in Go (dem Frontend nicht vertrauen)
- **Kontext verwenden** – für Authentifizierung und Autorisierung
- **Dateipfade bereinigen** – vor Dateioperationen
- **Aufrufrate begrenzen** – bei aufwendigen Operationen

[Mehr über Sicherheit erfahren →](/guides/security/)

## Nächste Schritte

**Anwendungslebenszyklus** – Start, Beenden und Lebenszyklus-Hooks verstehen [Mehr erfahren →](/concepts/lifecycle/)

**Go-Frontend-Bridge** – Die Funktionsweise der Bridge im Detail kennenlernen [Mehr erfahren →](/concepts/bridge/)

**Build-System** – Verstehen, wie Wails deine Anwendung erstellt [Mehr erfahren →](/concepts/build-system/)

**Mit der Entwicklung beginnen** – Wenden Sie das Gelernte in einem Tutorial an: [Tutorials →](/tutorials/03-notes-vanilla/)

---

**Fragen zur Architektur?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder lesen Sie in der [API-Referenz](/reference/overview/) nach.
