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
- **CI/CD-Tests** – Integrationstests ohne Display-Server ausführen
- **Microservices** – Wails-Bindings in Headless-Backend-Diensten verwenden

## Schnellstart

Der Servermodus wird über das Build-Tag `server` aktiviert. Ihr Anwendungscode bleibt unverändert – Sie führen den Build lediglich mit diesem Tag aus:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Hier ist ein Minimalbeispiel:

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

Derselbe Code kann für den Desktopmodus (ohne Tag) oder den Servermodus (mit `-tags server`) gebaut werden.

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

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Funktionen

### Endpunkt für Zustandsprüfungen

Unter `/health` ist automatisch ein Endpunkt für Zustandsprüfungen verfügbar:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Dies eignet sich für:

- Liveness- und Readiness-Probes von Kubernetes
- Zustandsprüfungen von Load-Balancern
- Überwachungssysteme

### Dienst-Bindings

Alle Dienst-Bindings funktionieren genauso wie im Desktopmodus:

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

Das Frontend kann diese Bindings über die standardmäßige Wails-Runtime aufrufen:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Ereignisse

Ereignisse funktionieren im Servermodus bidirektional:

- **Frontend an Backend**: Vom Browser ausgelöste Ereignisse werden über HTTP gesendet und von Ihren Go-Event-Handlern empfangen
- **Backend an Frontend**: Von Go ausgelöste Ereignisse werden über WebSocket an alle verbundenen Browser übertragen

Jeder Browser-Tab wird als „Fenster“ mit einem eindeutigen Namen dargestellt (`browser-1`, `browser-2` usw.) und ist über `event.Sender` zugänglich:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

Im Frontend:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Geordnetes Herunterfahren

Der Server verarbeitet die Signale `SIGINT` und `SIGTERM` ordnungsgemäß:

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
| Bildschirminformationen | Verfügbar | Kein nativer Bildschirm; `GetPrimary()` gibt `nil` zurück |
| Dienst-Bindings | Funktionieren | Funktionieren |
| Ereignisse | Funktionieren | Funktionieren (über WebSocket) |
| Assets | Über WebView | Über HTTP |
| CGO erforderlich | Ja | Nein |

### Verhalten der Fenster-API

Im Servermodus werden fensterbezogene APIs sicher behandelt:

- `app.Window.NewWithOptions(options)` – Gibt ein Go-Fensterobjekt zurück, das nicht nil ist; dadurch wird kein natives Desktopfenster erstellt
- `app.Hide()` / `app.Show()` – Keine Wirkung
- `app.Screen.GetPrimary()` – Gibt `*application.Screen` zurück; wenn kein primärer Bildschirm verfügbar ist, ist dies `nil`. Prüfen Sie den Wert vor der Dereferenzierung

Unterscheiden Sie Go-Fensterobjekte von der Verfügbarkeit nativer Fenster. Sichern Sie beim Erstellen eines Servers ausschließlich für den Desktop vorgesehenes Verhalten ab; ein Go-Fensterobjekt, das nicht nil ist, belegt nicht, dass ein natives Fenster oder ein Bildschirm vorhanden ist.

### Unterstützte Plattformen

Der Servermodus ist vollständig in Go implementiert und lässt sich ohne CGO für jedes von Wails unterstützte `GOOS` kompilieren, einschließlich `android` und `ios`. Für mobile Zielplattformen wird lediglich ein Headless-Server gebaut:

- Die nativen mobilen APIs `application.Android` und `application.IOS` sind nicht verfügbar. Code, der sie verwendet, muss durch eine `!server`-Build-Bedingung von Server-Builds ausgeschlossen werden.
- Standardmäßig wird der Desktop-Logger verwendet, der in Entwicklungs-Builds in stderr schreibt und nicht in logcat oder die WebView-Konsole.

## Build für den Produktivbetrieb

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

Wails-Projekte enthalten eine sofort einsatzbereite Docker-Konfiguration. So bauen Sie Ihre Anwendung und führen sie in einem Container aus:

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

Für komplexere Bereitstellungen finden Sie hier eine Docker-Compose-Konfiguration mit Zustandsprüfungen:

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
Das Healthcheck-Beispiel verwendet `wget`. Wenn Sie ein Distroless-Basis-Image verwenden, müssen Sie entweder ein Healthcheck-Binärprogramm in Ihr Image aufnehmen oder einen externen Mechanismus für Zustandsprüfungen verwenden (z. B. Dockers Option `curl` oder einen Sidecar-Container).

@end

### Benutzerdefiniertes Dockerfile

Wenn Sie mehr Kontrolle benötigen, können Sie ein eigenes Dockerfile erstellen. Entscheidend ist, `WAILS_SERVER_HOST=0.0.0.0` festzulegen, damit der Server Verbindungen von außerhalb des Containers akzeptiert:

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

1. **Standardmäßig an localhost binden** – `0.0.0.0` nur bei Bedarf verwenden
2. **In der Produktion TLS verwenden** – `ServerOptions.TLS` konfigurieren
3. **Hinter einem Reverse-Proxy betreiben** – nginx/traefik für zusätzliche Sicherheit verwenden
4. **WebSockets und Runtime-Aufrufe auf demselben Origin belassen** – mit `WebSocketOriginPatterns` ausschließlich vertrauenswürdige Origins hinzufügen; `WebSocketAllowAllOrigins` vermeiden
5. **Alle Eingaben validieren** – dieselben Sicherheitspraktiken wie bei jeder Webanwendung anwenden

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
| `WAILS_SERVER_HOST` | Netzwerkschnittstelle, an die gebunden werden soll | `localhost` |
| `WAILS_SERVER_PORT` | Port, auf dem eingehende Verbindungen angenommen werden | `8080` |

Diese Variablen haben Vorrang vor `ServerOptions` in Ihrem Code. Deshalb setzen die Docker-Beispiele `WAILS_SERVER_HOST=0.0.0.0`: So kann der Container externe Verbindungen akzeptieren, ohne dass Änderungen an Ihrer Anwendung erforderlich sind.

## Siehe auch

- [Benutzerdefinierter Transport](/guides/custom-transport/) – für erweiterte IPC-Anpassungen
- [Dienste](/features/bindings/services/) – Dokumentation zur Dienstbindung
- [Ereignisse](/guides/events-reference/) – Dokumentation zum Ereignissystem

### Größe von Runtime-Anfragen

Anfragen an `/wails/runtime` sind vor der JSON-Verarbeitung auf 64 MiB begrenzt. Eine gewöhnliche Anfrage oberhalb dieses Grenzwerts erhält HTTP 413; dies gilt auch für Anfragen ohne `Content-Length`. Für in Blöcke aufgeteilte Runtime-Uploads gelten weiterhin die separaten Grenzwerte von 1 MiB pro Block und 64 MiB für die zusammengesetzten Nutzdaten. Verwenden Sie gegebenenfalls Anwendungs-Middleware oder einen Reverse-Proxy, um einen niedrigeren Grenzwert festzulegen.

### Origin von Runtime-Anfragen

`/wails/runtime` akzeptiert nur `POST`-Anfragen; andere Methoden erhalten HTTP 405. Eine Browseranfrage, deren `Origin` weder dem Host der Anfrage noch `WebSocketOriginPatterns` entspricht, erhält HTTP 403. Dasselbe gilt für eine Anfrage ohne `Origin`, deren `Sec-Fetch-Site` den Wert `cross-site` oder `same-site` hat. Anfragen, die keinen der beiden Header enthalten, etwa von `curl` oder anderen Nicht-Browser-Clients, werden akzeptiert. Wenn ein Reverse-Proxy den Header `Host` umschreibt, fügen Sie den öffentlichen Origin zu `WebSocketOriginPatterns` hinzu.

`WebSocketOriginPatterns` ist für Bereitstellungen mit Reverse-Proxy vorgesehen, bei denen sich der öffentliche Origin von dem für die Anwendung sichtbaren `Host` unterscheidet. Frontends, die von einem anderen Origin bereitgestellt werden und `/wails/runtime` direkt aufrufen, werden nicht unterstützt: Der Endpunkt gibt keine CORS-Header aus und beantwortet Preflight-Anfragen mit `OPTIONS` mit HTTP 405. Stellen Sie das Frontend über den Wails-Server bereit oder betreiben Sie beide hinter demselben Proxy-Origin.

### Prüfung des Loopback-Hosts

Wenn der Server an eine Loopback-Adresse gebunden ist (den Standardwert `localhost`, eine beliebige `127.0.0.0/8`-Adresse oder `::1`, unabhängig davon, ob sie in `ServerOptions.Host` oder `WAILS_SERVER_HOST` festgelegt wurde), muss jede Anfrage einen `Host`-Header enthalten, der `localhost`, ein IP-Literal, den Bind-Host oder einen Host angibt, der dem Host-Teil eines `WebSocketOriginPatterns`-Eintrags entspricht. Jeder andere Wert für `Host` erhält HTTP 403. Dies verhindert DNS-Rebinding, bei dem eine schädliche Webseite ihren eigenen Hostnamen auf `127.0.0.1` verweisen lässt, sodass ihre Anfragen die oben beschriebenen Same-Origin-Prüfungen bestehen. Der Port wird nicht verglichen, und Schreibweisen wie `localhost.` (mit abschließendem Punkt) oder `127.1` werden abgelehnt. Die Prüfung gilt auch, wenn `WebSocketAllowAllOrigins` festgelegt ist. Wenn ein lokaler Reverse-Proxy seinen öffentlichen `Host` weiterleitet, fügen Sie diesen Host zu `WebSocketOriginPatterns` hinzu. Bindungen an andere Adressen, beispielsweise `0.0.0.0` in Containern, sind davon nicht betroffen.
