---
title: "Rohnachrichten"
description: "Benutzerdefinierte Frontend-Backend-Kommunikation für leistungskritische Anwendungen implementieren"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

Rohnachrichten stellen einen systemnahen Kommunikationskanal zwischen Frontend und Backend bereit und umgehen dabei das standardmäßige Bindungssystem. Dadurch wird Komfort zugunsten höherer Geschwindigkeit aufgegeben.

## Wann Rohnachrichten sinnvoll sind

Rohnachrichten eignen sich am besten für extreme Sonderfälle:

- **Extrem hochfrequente Aktualisierungen** – Tausende Nachrichten pro Sekunde, bei denen jede Mikrosekunde zählt
- **Benutzerdefinierte Nachrichtenprotokolle** – Wenn Sie vollständige Kontrolle über das Übertragungsformat benötigen

@note{type="tip"}
Für nahezu alle Anwendungsfälle werden standardmäßige [Service-Bindungen](/features/bindings/services/) empfohlen, da sie Typsicherheit, automatische Serialisierung und eine bessere Entwicklererfahrung bei vernachlässigbarem Zusatzaufwand bieten.

@end

## Backend einrichten

Konfigurieren Sie `RawMessageHandler` in den Optionen Ihrer Anwendung:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### Handler-Signatur

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| Parameter | Typ | Beschreibung |
| --- | --- | --- |
| `window` | `Window` | Das Fenster, das die Nachricht gesendet hat |
| `message` | `string` | Der unverarbeitete Nachrichteninhalt |
| `originInfo` | `*application.OriginInfo` | Ursprungsinformationen zur Quelle der Nachricht |

#### OriginInfo-Struktur

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| Feld | Typ | Beschreibung |
| --- | --- | --- |
| `Origin` | `string` | Die Ursprungs-URL des Dokuments, das die Nachricht gesendet hat |
| `TopOrigin` | `string` | Die Ursprungs-URL der obersten Ebene (kann in iframes von Origin abweichen) |
| `IsMainFrame` | `bool` | Gibt an, ob die Nachricht aus dem Hauptframe stammt |

#### Plattformspezifische Verfügbarkeit

- **macOS**: `Origin` und `IsMainFrame` werden bereitgestellt
- **Windows**: `Origin` und `TopOrigin` werden bereitgestellt
- **Linux**: Nur `Origin` wird bereitgestellt

### Ursprungsvalidierung

@note{type="caution"}
Gehen Sie niemals davon aus, dass eine Nachricht sicher ist, nur weil sie Ihren Handler erreicht. Die Ursprungsinformationen müssen validiert werden, bevor vertrauliche oder zustandsändernde Vorgänge verarbeitet werden.

@end

**Überprüfen Sie stets den Ursprung eingehender Nachrichten, bevor Sie sie verarbeiten.** Der Parameter `originInfo` enthält wichtige Sicherheitsinformationen, die validiert werden müssen, um unbefugten Zugriff zu verhindern. Schädliche oder kompromittierte Inhalte sowie unbeabsichtigt ausgeführte Skripte könnten Rohnachrichten senden. Ohne Ursprungsvalidierung verarbeiten Sie möglicherweise Befehle aus nicht vertrauenswürdigen Quellen. Stellen Sie mit `originInfo` sicher, dass Nachrichten aus den erwarteten Quellen stammen.

### Wichtige Validierungspunkte

- **Prüfen Sie stets `Origin`** – Überprüfen Sie, ob der Ursprung den erwarteten vertrauenswürdigen Quellen entspricht (bei lokalen Assets in der Regel `wails://wails` oder `http://wails.localhost` beziehungsweise dem spezifischen Ursprung Ihrer Anwendung)
- **Validieren Sie `IsMainFrame`** (macOS) – Achten Sie darauf, ob die Nachricht aus einem iframe stammt, da dies auf eingebettete Inhalte mit einem anderen Sicherheitskontext hinweisen kann
- **Verwenden Sie `TopOrigin`** (Windows) – Überprüfen Sie bei Inhalten in Frames den Ursprung der obersten Ebene
- **Unerwartete Ursprünge ablehnen** – Sorgen Sie für ein sicheres Fehlerverhalten, indem Sie Nachrichten aus Ursprüngen ablehnen, die Sie nicht ausdrücklich zulassen

@note{type="info"}
Nachrichten mit dem Präfix `wails:` sind für die interne Kommunikation von Wails reserviert und werden nicht an Ihren Handler weitergeleitet.

@end

## Frontend einrichten

Senden Sie Rohnachrichten mit `System.invoke()`:

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### Vorgefertigtes Bundle verwenden

Wenn Sie npm nicht verwenden, greifen Sie über das globale Objekt `wails` auf `invoke` zu:

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## Strukturierte Nachrichten

Serialisieren Sie komplexe Daten als JSON:

### Frontend

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### Backend

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## Leistungsvergleich

| Ansatz | Zusatzaufwand | Typsicherheit | Anwendungsfall |
| --- | --- | --- | --- |
| Service-Bindings | Höher | Vollständig | Allgemeiner Einsatz |
| Rohnachrichten | Minimal | Manuell | Hochfrequent, leistungskritisch |

### Benchmark-Beispiel

Bei einfachen Nutzdaten können Rohnachrichten deutlich mehr Nachrichten pro Sekunde verarbeiten als Service-Bindings:

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## Vollständiges Beispiel

Das folgende vollständige Beispiel implementiert ein einfaches Befehlsprotokoll:

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## Bewährte Vorgehensweisen

### Empfohlen

- Verwenden Sie Rohnachrichten für tatsächlich leistungskritische Ausführungspfade
- Implementieren Sie in Ihrem Handler eine angemessene Fehlerbehandlung
- Verwenden Sie Events, um Antworten an das Frontend zurückzusenden
- Ziehen Sie JSON für strukturierte Daten in Betracht
- Verarbeiten Sie Nachrichten schnell, um Blockierungen zu vermeiden

### Nicht empfohlen

- Verwenden Sie keine Rohnachrichten, wenn Service-Bindings ausreichen würden
- Vergessen Sie nicht, eingehende Nachrichten zu validieren
- Blockieren Sie den Handler nicht durch lang laufende Operationen (verwenden Sie goroutines)
- Ignorieren Sie den Fensterparameter nicht, wenn Antworten an bestimmte Fenster gesendet werden müssen

## Überlegungen bei mehreren Fenstern

Der Parameter `window` gibt an, welches Fenster die Nachricht gesendet hat. Dadurch können Sie:

- Antworten an das richtige Fenster senden
- Fensterspezifisches Verhalten implementieren
- Nachrichtenquellen für die Fehlersuche nachverfolgen

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## Nächste Schritte

- [Service-Bindings](/features/bindings/services/) – Standardansatz für die meisten Anwendungen
- [Events](/guides/events-reference/) – Eventsystem für die Kommunikation vom Backend zum Frontend
- [Performance](/guides/performance/) – Allgemeine Leistungsoptimierung
