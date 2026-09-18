---
title: "LLM-Steuerung (MCP)"
description: "LLM-Agenten können Ihre App über das Model Context Protocol testen und steuern"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Der integrierte MCP-Server ist experimentell und die API kann sich in zukünftigen Versionen ändern.

@end

Wails v3 verfügt über einen integrierten [Model Context Protocol](https://modelcontextprotocol.io) (MCP)-Server, mit dem LLM-Agenten – Claude Code, IDE-Assistenten oder beliebige MCP-Clients – eine **laufende** Wails-Anwendung untersuchen, testen und steuern können.

Verwenden Sie für die Automatisierung des Projektlebenszyklus den separaten `wails3 mcp` CLI-Server. Damit können Agenten Projekte untersuchen und initialisieren, Diagnosen ausführen, Builds und Entwicklungsaufträge starten, Bindings generieren, benannte Taskfile-Tasks ausführen und begrenzte Auftragsausgaben abrufen. Der CLI-Server ist standardmäßig auf das aktuelle Verzeichnis beschränkt und ermöglicht keine beliebige Shell-Ausführung. Einzelheiten zu Transport, Authentifizierung und Tools finden Sie in der [CLI-MCP-Dokumentation](/guides/cli/#mcp).

Wenn die Funktion aktiviert ist, kann ein mit Ihrer App verbundener Agent:

- **Fenster auflisten und steuern** – Größe, Position, Fokus, Vollbildmodus, Entwicklertools, Neuladen, …
- **Das DOM untersuchen** – Elemente abfragen, HTML abrufen und eine strukturelle Momentaufnahme erstellen
- **JavaScript auswerten** – beliebigen Code in einem beliebigen Fenster ausführen und das Ergebnis abrufen
- **Benutzereingaben simulieren** – Mausbewegungen, Klicks, Zieh- und Bildlaufaktionen werden mit einem **animierten Bildschirmcursor** dargestellt, sodass Sie die Arbeit des Agenten beobachten können
- **Text und Tasten eingeben** – realistische zeichenweise Ereignisse, die mit kontrollierten React-Eingabefeldern funktionieren
- **Gebundene Go-Methoden aufrufen** sowie Anwendungsereignisse auslösen und abwarten

## Funktionsweise

Der MCP-Server wird nur dann in Ihre Anwendung kompiliert, wenn das **`mcp` Build-Tag** vorhanden ist. Ohne dieses Tag fehlt der Servercode vollständig in der Binärdatei – keine Laufzeitbelastung, keine offenen Ports, keine Angriffsfläche.

Wenn das Tag vorhanden ist, startet der Server automatisch in `App.Run()`, bindet sich standardmäßig an `127.0.0.1:9099` und protokolliert seinen Endpunkt. Benutzercode ist nicht erforderlich.

## Tutorial

### Schritt 1 – eine normale Wails-Anwendung schreiben

MCP erfordert weder Importe noch eine Registrierung. Erstellen Sie Ihre App wie gewohnt:

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Schritt 2 – mit dem Tag `mcp` bauen oder ausführen

@tabs
[Wails CLI (empfohlen)]
Setzen Sie `WAILS_MCP=1`; die Wails CLI fügt das Tag `mcp` automatisch hinzu:

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Direkt mit Go]
Übergeben Sie das Tag direkt an `go run` oder `go build`:

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows (PowerShell)]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

Beim Start protokolliert die Anwendung den MCP-Endpunkt:

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### Schritt 3 – einen Client verbinden

Der Server verwendet den **streamfähigen HTTP-Transport von MCP**. Stellen Sie mit einem beliebigen MCP-kompatiblen Client eine Verbindung her.

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

Bitten Sie Claude anschließend, mit Ihrer App zu interagieren:

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code (GitHub Copilot)]
Fügen Sie Folgendes zu `.vscode/settings.json` hinzu:

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[Andere Clients]
Richten Sie einen beliebigen MCP-Client, der den streamfähigen HTTP-Transport unterstützt, auf folgende Adresse:

```
http://127.0.0.1:9099/mcp
```

@end

### Schritt 4 – eine Testsitzung ausführen

Fordern Sie den Agenten auf, Ihre Anwendung zu testen. Hier sind einige Beispiel-Prompts:

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## Konfiguration

Die gesamte Konfiguration erfolgt über Umgebungsvariablen – Codeänderungen sind nicht erforderlich.

| Umgebungsvariable | Standardwert | Beschreibung |
| --- | --- | --- |
| `WAILS_MCP` | (nicht gesetzt) | Setzen Sie den Wert auf `1`, `true`, `on` oder `yes`, damit bei Verwendung der Wails CLI das Build-Tag `mcp` automatisch hinzugefügt wird. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Schnittstelle, an die der Server gebunden wird. Bindungen an Nicht-Loopback-Adressen erfordern `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | nicht gesetzt | Optionales Bearer-Token bei Loopback; für andere Bindungsadressen erforderlich. Clients senden `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Port, auf dem der Server lauscht. Setzen Sie den Wert auf `0`, um einen zufällig zugewiesenen freien Port zu verwenden (wird im Protokoll ausgegeben). |
| `WAILS_MCP_TIMEOUT` | `30000` | Standardzeitlimit für die JavaScript-Auswertung in **Millisekunden**. |
| `WAILS_MCP_HIDE_CURSOR` | (nicht gesetzt) | Setzen Sie den Wert auf `1` oder `true`, um die animierte Cursor-Einblendung zu deaktivieren. |

Beispiel – benutzerdefinierter Port und ein Zeitlimit von 60 Sekunden:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Verfügbare Tools

| Tool | Zweck |
| --- | --- |
| `app_info` | Anwendungsinformationen: Plattform, Architektur, alle Fenster, MCP-Endpunkt |
| `windows_list` | Alle Fenster mit Geometrie und Status auflisten |
| `window_control` | Fokussieren, Größe ändern, verschieben, Vollbildmodus, Entwicklertools, neu laden, URL festlegen, … (22 Aktionen) |
| `js_eval` | JavaScript in einem Fenster auswerten (asynchroner Funktionsrumpf, `return` für den Wert) |
| `dom_html` | HTML der Seite oder eines bestimmten Elements abrufen |
| `dom_query` | Elemente per CSS-Selektor finden — Tag, Text, Begrenzungsrechteck, Sichtbarkeit |
| `screenshot_dom` | Strukturelle Momentaufnahme der sichtbaren Seite (DOM-basiert, ohne Pixel) |
| `mouse_move` | Cursor zu einem Punkt oder CSS-Selektor animieren |
| `mouse_click` | Mit dem animierten Cursor klicken (links/rechts/Mitte, Doppelklick, Modifikatortasten) |
| `mouse_drag` | Mit dem animierten Cursor ziehen (unterstützt HTML5-Drag-and-Drop-Elemente) |
| `mouse_scroll` | An einem Punkt oder Element scrollen |
| `keyboard_type` | Text Zeichen für Zeichen mit realistischen Ereignissen eingeben |
| `keyboard_press` | Einzelne Taste (Enter, Tab, Escape, ArrowDown, …) mit optionalen Modifikatortasten drücken |
| `call_bound_method` | Methode eines gebundenen Go-Dienstes aufrufen, z. B. `main.GreetService.Greet` |
| `emit_event` | Wails-Anwendungsereignis auslösen |
| `wait_for_event` | Auf ein Wails-Anwendungsereignis warten und dessen Daten zurückgeben |

### Unterstützung mehrerer Fenster

Alle Werkzeuge, die auf ein Fenster einwirken, akzeptieren ein optionales `window`-Argument mit dem **Namen** des Fensters (festgelegt über `WebviewWindowOptions.Name`). Wird es weggelassen, verwendet das Werkzeug das derzeit fokussierte Fenster oder das erste Fenster, falls keines fokussiert ist.

```
List all windows, then click the "New" button in the window named "editor".
```

### Elemente auswählen

Maus- und Tastaturwerkzeuge akzeptieren wahlweise:

- **CSS-Selektor** — `selector: "#submit-btn"` (das Element wird automatisch in den sichtbaren Bereich gescrollt)
- **Koordinaten** — `x: 400, y: 300` (CSS-Pixel relativ zum Viewport)

Stellen Sie bei Ziehvorgängen `from_` und `to_` voran:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Sicherheit

@note{type="caution"}
Der MCP-Server ermöglicht die vollständige programmgesteuerte Kontrolle über Ihre Anwendung. Jeder, der auf seine Werkzeuge zugreifen kann, kann das DOM lesen, JavaScript auswerten, auf Schaltflächen klicken und Go-Methoden aufrufen.

@end

- Der Server bindet standardmäßig an `127.0.0.1`. Browser-Ursprünge müssen HTTP(S)-Loopback-Ursprünge sein; opake (`null`), fehlerhafte und fremde Ursprünge werden abgelehnt.
- Native MCP-Clients ohne Header werden weiterhin unterstützt. Ohne `WAILS_MCP_TOKEN` gelten lokale Prozesse und zulässige lokale Browser-Ursprünge als vertrauenswürdig; Ursprungsprüfungen sind keine Authentifizierung. Legen Sie ein Token mit hoher Entropie fest, um für alle `/mcp`-Aufrufe eine Bearer-Authentifizierung zu verlangen. Konfigurieren Sie dasselbe Token im `Authorization: Bearer <token>`-Header des Clients. Preflight-Anfragen benötigen das Token nicht.
- Der `/eval-result`-Callback verwendet statt des Bearer-Tokens des Clients unvorhersehbare IDs für jede einzelne Auswertung, sodass die Übermittlung von Webview-Ergebnissen kompatibel bleibt.
- Produktions-Builds sollten das `mcp`-Tag **nicht** enthalten. Die Wails-CLI fügt es nur hinzu, wenn `WAILS_MCP=1` ausdrücklich festgelegt wurde, und das standardmäßige `wails3 build` enthält keinerlei Servercode.
- Wenn Sie den Server über eine Schnittstelle verfügbar machen müssen, die keine Loopback-Schnittstelle ist (z. B. für Tests im LAN), legen Sie `WAILS_MCP_HOST=0.0.0.0` und ein `WAILS_MCP_TOKEN` mit hoher Entropie fest; ohne Token schlägt der Start fehl. Verwenden Sie in nicht vertrauenswürdigen Netzwerken einen verschlüsselten Tunnel oder einen TLS-terminierenden Proxy, da der integrierte Listener HTTP verwendet.

## Beispielanwendung

Eine vollständige Playground-Anwendung, die alle Werkzeuge demonstriert, ist unter [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp) verfügbar. Sie umfasst:

- Zähler mit Schaltflächen zum Erhöhen und Zurücksetzen
- Namenseingabe mit den gebundenen Methoden Greet, Add und Shout
- Quelle und Ziel für HTML5-Drag-and-Drop
- Scrollbare Liste (50 Elemente)
- Ereignisprotokoll

Führen Sie sie wie folgt aus:

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Verbinden Sie anschließend Claude Code oder einen beliebigen MCP-Client mit `http://127.0.0.1:9099/mcp` und weisen Sie ihn an, die Benutzeroberfläche zu testen.

## Feedback

Der integrierte MCP-Server ist ein Experiment, und Ihr Feedback entscheidet über seine weitere Entwicklung. Wenn Sie ihn ausprobieren, möchten wir gern erfahren, welchen Client und welche Werkzeuge Sie verwendet haben, was Sie erwartet haben und was tatsächlich geschehen ist und ob es hilfreich war, einen Agenten Ihre Anwendung steuern zu lassen — die hilfreichsten Berichte beschreiben genau, was ausgeführt wurde. Berichten Sie uns davon in der [Feedback-Diskussion zum MCP-Server](https://github.com/wailsapp/wails/discussions/5692).
