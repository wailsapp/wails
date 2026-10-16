---
title: "Server-Build"
description: "Wails-Anwendungen als HTTP-Server ohne natives GUI-Fenster ausführen"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 unterstützt einen Servermodus, in dem Sie Ihre Anwendung als reinen HTTP-Server ausführen können, ohne native Fenster zu erstellen oder GUI-Abhängigkeiten zu benötigen. Dadurch lässt sich dieselbe Wails-Anwendung auf Servern, in Containern und in Webbrowsern bereitstellen.

Der Servermodus eignet sich für:

- **Docker-/Container-Bereitstellungen** – Ausführung ohne X11-/Wayland-Abhängigkeiten
- **Serverseitige Anwendungen** – Bereitstellung als Webserver, auf den über einen Browser zugegriffen werden kann
- **Reiner Webzugriff** – Gemeinsame Codebasis für Desktop und Web
- **CI/CD-Tests** – Integrationstests ohne Displayserver ausführen
- **Microservices** – Wails-Bindings in Headless-Backenddiensten verwenden

## Schnellstart

Aktivieren Sie den Servermodus über das Build-Tag `server`. Ihr Anwendungscode bleibt unverändert – erstellen Sie den Build lediglich mit diesem Tag:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Hier ist ein minimales Beispiel:

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Derselbe Code kann für den Desktopmodus (ohne das Tag) oder den Servermodus (mit `-tags server`) gebaut werden.

## Konfiguration

### ServerOptions

Konfigurieren Sie den HTTP-Server mit `ServerOptions`:

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Funktionen

### Endpunkt für Integritätsprüfungen

Unter `/health` ist automatisch ein Endpunkt für Integritätsprüfungen verfügbar:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Dies eignet sich für:

- Liveness-/Readiness-Probes von Kubernetes
- Integritätsprüfungen von Load-Balancern
- Überwachungssysteme

### Service-Bindings

Alle Service-Bindings funktionieren genauso wie im Desktopmodus:

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

Das Frontend kann diese Bindings über die standardmäßige Wails-Laufzeit aufrufen:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Ereignisse

Ereignisse funktionieren im Servermodus bidirektional:

- **Vom Frontend zum Backend**: Vom Browser ausgelöste Ereignisse werden über HTTP gesendet und von Ihren Go-Ereignishandlern empfangen
- **Vom Backend zum Frontend**: Von Go ausgelöste Ereignisse werden über WebSocket an alle verbundenen Browser übertragen

Jeder Browser-Tab wird als „Fenster“ mit einem eindeutigen Namen (`browser-1`, `browser-2` usw.) dargestellt und ist über `event.Sender` zugänglich:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

Vom Frontend aus:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Kontrolliertes Herunterfahren

Der Server verarbeitet die Signale `SIGINT` und `SIGTERM` kontrolliert:

1. Akzeptiert keine neuen Verbindungen mehr
2. Wartet auf den Abschluss aktiver Anfragen (bis zu `ShutdownTimeout`)
3. Führt `OnShutdown`-Hooks aus
4. Fährt Dienste in umgekehrter Reihenfolge herunter

## Unterschiede zum Desktopmodus

| Funktion | Desktopmodus | Servermodus |
| --- | --- | --- |
| Native Fenster | Werden erstellt | Browserfenster (`browser-N`) |
| Infobereich | Verfügbar | Nicht verfügbar |
| Native Dialogfelder | Verfügbar | Nicht verfügbar |
| Anwendungsmenü | Verfügbar | Nicht verfügbar |
| Bildschirminformationen | Verfügbar | Gibt einen Fehler zurück |
| Service-Bindings | Funktionieren | Funktionieren |
| Ereignisse | Funktionieren | Funktionieren (über WebSocket) |
| Assets | Über Webview | Über HTTP |
| CGO erforderlich | Ja | Nein |

### Verhalten der Fenster-API

Im Servermodus werden fensterbezogene APIs sicher behandelt:

- `app.Window.NewWithOptions()` – Protokolliert eine Warnung und gibt nil zurück
- `app.Hide()` / `app.Show()` – Keine Aktion
- `app.Screen.GetPrimary()` – Gibt einen Fehler zurück

Dadurch kann Code, der auf Fenster verweist, ohne Absturz ausgeführt werden, obwohl Fensteroperationen keine Wirkung haben.

## Produktions-Build erstellen

### Task verwenden (empfohlen)

Mit `wails3 init` erstellte Projekte enthalten einen `build:server`-Task:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Manueller Build

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails-Projekte enthalten eine sofort einsatzbereite Docker-Konfiguration. So erstellen Sie Ihre Anwendung und führen sie in einem Container aus:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

Das ist alles! Ihre Anwendung ist unter `http://localhost:8080` verfügbar.

Sie können den Build mit einigen Optionen anpassen:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

Die generierte Datei `Dockerfile.server` erstellt ein minimales Image auf Distroless-Basis. Sie übernimmt die Netzwerkbindung automatisch, sodass Ihre Anwendung von außerhalb des Containers erreichbar ist.

### Docker Compose

Für komplexere Bereitstellungen finden Sie hier eine Docker-Compose-Konfiguration mit Integritätsprüfungen:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
Das Healthcheck-Beispiel verwendet `wget`. Wenn Sie ein Distroless-Basis-Image verwenden, müssen Sie entweder ein Healthcheck-Programm in Ihr Image aufnehmen oder einen externen Mechanismus zur Integritätsprüfung verwenden, beispielsweise die Docker-Option `curl` oder einen Sidecar-Container.

@end

### Benutzerdefiniertes Dockerfile

Wenn Sie mehr Kontrolle benötigen, können Sie ein eigenes Dockerfile erstellen. Wichtig ist, `WAILS_SERVER_HOST=0.0.0.0` festzulegen, damit der Server Verbindungen von außerhalb des Containers akzeptiert:

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## Sicherheitsaspekte

Beachten Sie beim Bereitstellen von Anwendungen im Servermodus Folgendes:

1. **Standardmäßig an localhost binden** – Verwenden Sie `0.0.0.0` nur bei Bedarf
2. **In der Produktion TLS verwenden** – Konfigurieren Sie `ServerOptions.TLS`
3. **Hinter einem Reverse-Proxy betreiben** – Verwenden Sie nginx/traefik für zusätzliche Sicherheit
4. **Für WebSockets denselben Ursprung beibehalten** – Fügen Sie mit `WebSocketOriginPatterns` nur vertrauenswürdige Ursprünge hinzu; vermeiden Sie `WebSocketAllowAllOrigins`
5. **Alle Eingaben validieren** – Es gelten dieselben Sicherheitspraktiken wie für jede andere Webanwendung

## Beispiel

Ein vollständiges Beispiel ist unter `v3/examples/server/` verfügbar:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Umgebungsvariablen

Für Bereitstellungsszenarien, in denen Sie die Serverkonfiguration ohne Codeänderungen überschreiben müssen, erkennt Wails die folgenden Umgebungsvariablen:

| Variable | Beschreibung | Standardwert |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Netzwerkschnittstelle für die Bindung | `localhost` |
| `WAILS_SERVER_PORT` | Port, auf dem Verbindungen entgegengenommen werden | `8080` |

Diese haben Vorrang vor `ServerOptions` in Ihrem Code. Deshalb setzen die Docker-Beispiele `WAILS_SERVER_HOST=0.0.0.0`: Dadurch kann der Container externe Verbindungen akzeptieren, ohne dass Änderungen an Ihrer Anwendung erforderlich sind.

## Siehe auch

- [Benutzerdefinierter Transport](/guides/custom-transport/) – Für erweiterte IPC-Anpassungen
- [Dienste](/features/bindings/services/) – Dokumentation zur Dienstbindung
- [Ereignisse](/guides/events-reference/) – Dokumentation zum Ereignissystem

### Größe von Runtime-Anfragen

Anfragen an `/wails/runtime` sind vor der JSON-Verarbeitung auf 64 MiB begrenzt. Eine gewöhnliche Anfrage oberhalb dieses Grenzwerts erhält HTTP 413; dies gilt auch für Anfragen ohne `Content-Length`. In Blöcke aufgeteilte Runtime-Uploads behalten ihre separaten Grenzwerte von 1 MiB pro Block und 64 MiB für die zusammengesetzten Nutzdaten bei. Verwenden Sie gegebenenfalls Anwendungs-Middleware oder einen Reverse-Proxy, um einen niedrigeren Grenzwert festzulegen.
