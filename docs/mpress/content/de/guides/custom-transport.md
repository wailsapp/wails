---
title: "Benutzerdefinierte Transportschicht erstellen"
description: "Erfahren Sie, wie Sie eine eigene benutzerdefinierte IPC-Transportschicht für Wails v3 erstellen und anpassen."
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Mit Wails v3 können Sie eine benutzerdefinierte IPC-Transportschicht bereitstellen und dabei alle generierten Bindings sowie die Ereigniskommunikation beibehalten. So können Sie den standardmäßigen, auf HTTP Fetch basierenden Transport durch WebSockets, benutzerdefinierte Protokolle oder einen beliebigen anderen Transportmechanismus ersetzen.

## Übersicht

Standardmäßig verwendet Wails HTTP-Fetch-Anfragen vom Frontend, um über `/wails/runtime` mit dem Backend zu kommunizieren. Mit der API für benutzerdefinierte Transporte können Sie:

- Den HTTP-Transport durch WebSockets, gRPC oder ein beliebiges benutzerdefiniertes Protokoll ersetzen
- Die vollständige Kompatibilität mit der Wails-Codegenerierung beibehalten
- Alle vorhandenen Bindings, Ereignisse, Dialoge und anderen Wails-Funktionen beibehalten
- Eine eigene Verbindungsverwaltung, Authentifizierung und Fehlerbehandlung implementieren

## Architektur

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## Vollständiges WebSocket-Beispiel ausführen

Das [WebSocket-Transportbeispiel](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) enthält beide Seiten des Protokolls: Anfrage-IDs und Antworten, Lebenszyklusverwaltung, den JavaScript-Client und die Ereignisübermittlung. Verwenden Sie diese Dateien als Ausgangspunkt und nicht einen Transport, der nur `Start` und `Stop` implementiert.

Führen Sie in einer Arbeitskopie des Wails-Repositorys, in der die nativen Plattformabhängigkeiten installiert sind, Folgendes aus:

```bash
cd v3/examples/websocket-transport
go run .
```

Das Beispiel öffnet ein natives Fenster und startet einen WebSocket-Endpunkt an Port 9099. Verwenden Sie es in einer vertrauenswürdigen lokalen Entwicklungsumgebung. Bevor Sie eine über das Netzwerk erreichbare Variante bereitstellen, schränken Sie den Listener ein, validieren Sie Ursprünge und authentifizieren Sie Clients. Die Beispielrichtlinie für Ursprünge ist freizügig und stellt keine Sicherheitsrichtlinie für den Produktivbetrieb dar.

### Backend-Vertrag

Ein Transport muss alle drei erforderlichen Methoden implementieren. `WebSocketTransport` aus dem Beispiel implementiert `Start`, `JSClient` und `Stop`; `JSClient` stellt den JavaScript-Shim bereit, den Wails in die Seite injiziert. Fügen Sie Ihrer Implementierung eine Schnittstellenprüfung zur Kompilierzeit hinzu:

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Konfigurieren Sie die Anwendung mit dem Konstruktor aus dem Beispiel:

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

Dieses Feld gehört neben den üblichen Assets und Diensten in `application.Options`. Kopieren Sie sowohl die Implementierungsdateien als auch den Konstruktoraufruf. `Start` empfängt `MessageProcessor`; `Stop` muss Listener, aktive Verbindungen und ausstehende Vorgänge schließen.

### Frontend-Vertrag

`websocket-transport.js` aus dem Beispiel implementiert den Runtime-Transport einschließlich der Zuordnung von Anfragen zu Antworten und der Zurückweisung fehlgeschlagener Aufrufe. Nachdem Sie diese Datei über Ihre Assets bereitgestellt haben, initialisieren Sie sie, bevor Sie generierte Dienst-Bindings aufrufen:

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Generierte Dienst-Bindings bleiben unverändert. Wenn Ihr Transport den Client ersetzt, muss er die von Ihrer Anwendung verwendeten Methoden übertragen und deren Fehlersemantik beibehalten. Für die Ereignisübermittlung ist der optionale Ereignistransportvertrag erforderlich; eine beliebige Implementierung von `Transport` stellt ihn nicht automatisch bereit. Das vollständige WebSocket-Beispiel veranschaulicht ihn.

Verwenden Sie für Remote-Bereitstellungen TLS sowie explizite Richtlinien für Authentifizierung und Ursprünge. Beim ausschließlichen Hosting im Browser müssen zusätzlich die Anwendungs-Assets bereitgestellt werden, wie unten beschrieben.

## API-Referenz

### Transport-Schnittstelle

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport-Schnittstelle (optional)

Implementieren Sie für browserbasierte Bereitstellungen oder wenn Sie sowohl Assets als auch IPC über Ihren benutzerdefinierten Transport bereitstellen möchten, die Schnittstelle `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Wann diese Schnittstelle implementiert werden sollte:**

- Ausführen der Anwendung in einem Browser statt in einer Webview
- Bereitstellen von Assets über HTTP neben Ihrem benutzerdefinierten IPC-Transport
- Erstellen von über das Netzwerk erreichbaren Anwendungen

**Beispielimplementierung:**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

Beim Aufruf von `ServeAssets()` stellt der assetHandler Folgendes bereit:

- Alle statischen Assets (HTML, CSS, JS, Bilder usw.)
- `/wails/runtime.js` – die Wails-Runtime-Bibliothek

## Siehe auch

- `transport.go` – zentrale Transportschnittstellen und -typen
- `messageprocessor.go` – der zugrunde liegende Nachrichtenprozessor, der die gesamte Wails-IPC verarbeitet
