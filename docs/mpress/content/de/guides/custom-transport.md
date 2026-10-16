---
title: "Benutzerdefinierte Transportschicht erstellen"
description: "Erfahren Sie, wie Sie eine eigene benutzerdefinierte IPC-Transportschicht für Wails v3 erstellen und anpassen"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Mit Wails v3 können Sie eine benutzerdefinierte IPC-Transportschicht bereitstellen und dabei alle generierten Bindings sowie die Ereigniskommunikation beibehalten. Dadurch können Sie den standardmäßigen, auf HTTP-Fetch basierenden Transport durch WebSockets, benutzerdefinierte Protokolle oder einen beliebigen anderen Transportmechanismus ersetzen.

## Überblick

Standardmäßig verwendet Wails HTTP-Fetch-Anfragen vom Frontend, um über `/wails/runtime` mit dem Backend zu kommunizieren. Mit der API für benutzerdefinierte Transporte können Sie:

- Den HTTP-Transport durch WebSockets, gRPC oder ein beliebiges benutzerdefiniertes Protokoll ersetzen
- Die vollständige Kompatibilität mit der Wails-Codegenerierung beibehalten
- Alle vorhandenen Bindings, Ereignisse, Dialoge und sonstigen Wails-Funktionen beibehalten
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

## Verwendung

### 1. Transportschnittstelle implementieren

Erstellen Sie einen benutzerdefinierten Transport, indem Sie die Schnittstelle `Transport` implementieren:

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Anwendung konfigurieren

Übergeben Sie Ihren benutzerdefinierten Transport an die Anwendungsoptionen:

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Frontend-Runtime anpassen

Wenn Sie einen benutzerdefinierten Transport verwenden, müssen Sie die Frontend-Runtime so anpassen, dass sie Ihren Transport anstelle von HTTP-Fetch verwendet. Implementieren Sie die Schnittstelle `RuntimeTransport`, die zur Verarbeitung von Anfragen verwendet wird:

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## Hinweise

- Der standardmäßige HTTP-Transport funktioniert weiterhin, wenn kein benutzerdefinierter Transport angegeben ist
- Generierte Bindings bleiben unverändert – nur die Transportschicht ändert sich
- Ereignisse, Dialoge, die Zwischenablage und alle anderen Wails-Funktionen funktionieren transparent
- Sie sind für die Fehlerbehandlung, die Wiederverbindungslogik und die Sicherheit Ihres benutzerdefinierten Transports verantwortlich
- Das bereitgestellte WebSocket-Beispiel dient zur Veranschaulichung und muss für den Produktionseinsatz möglicherweise gehärtet werden

## API-Referenz

### Transportschnittstelle

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

Implementieren Sie für browserbasierte Bereitstellungen oder wenn Sie sowohl Assets als auch IPC über Ihren benutzerdefinierten Transport bereitstellen möchten die Schnittstelle `AssetServerTransport`:

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
- Bereitstellen von Assets über HTTP parallel zu Ihrem benutzerdefinierten IPC-Transport
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

Beim Aufruf von `ServeAssets()` stellt assetHandler Folgendes bereit:

- Alle statischen Assets (HTML, CSS, JS, Bilder usw.)
- `/wails/runtime.js` – die Wails-Runtime-Bibliothek

## Siehe auch

- `transport.go` – zentrale Transportschnittstellen und -typen
- `messageprocessor.go` – der zugrunde liegende Nachrichtenprozessor, der die gesamte Wails-IPC verarbeitet
