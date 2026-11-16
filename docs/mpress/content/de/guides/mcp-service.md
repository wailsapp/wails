---
title: "LLM-Steuerung (MCP)"
description: "LLM-Agenten können Ihre App über das Model Context Protocol testen und steuern"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Der integrierte MCP-Server ist experimentell; die API kann sich in zukünftigen Versionen ändern.

@end

Wails v3 verfügt über einen integrierten [Model Context Protocol](https://modelcontextprotocol.io) (MCP)-Server, über den LLM-Agenten – Claude Code, IDE-Assistenten oder beliebige MCP-Clients – eine **laufende** Wails-Anwendung untersuchen, testen und steuern können.

Verwenden Sie für die Automatisierung des Projektlebenszyklus den separaten `wails3 mcp`-CLI-Server. Damit können Agenten Projekte untersuchen und initialisieren, Diagnosen ausführen, Builds und Entwicklungsaufträge starten, Bindings generieren, benannte Taskfile-Tasks ausführen und begrenzte Auftragsausgaben abrufen. Der CLI-Server ist standardmäßig auf das aktuelle Verzeichnis beschränkt und ermöglicht keine beliebige Shell-Ausführung. Details zu Transport, Authentifizierung und Werkzeugen finden Sie in der [Dokumentation zu CLI-MCP](/guides/cli/#mcp).

Wenn diese Funktion aktiviert ist, kann ein mit Ihrer App verbundener Agent:

- **Fenster auflisten und steuern** – Größe, Position, Fokus, Vollbildmodus, Entwicklertools, Neuladen, …
- **Das DOM untersuchen** – Elemente abfragen, HTML abrufen, einen strukturellen Snapshot erstellen
- **JavaScript auswerten** – beliebigen Code in einem beliebigen Fenster ausführen und das Ergebnis abrufen
- **Benutzereingaben simulieren** – Mausbewegungen, Klicks, Ziehvorgänge und Scrollvorgänge werden mit einem **animierten Bildschirmzeiger** dargestellt, sodass Sie die Arbeit des Agenten beobachten können
- **Text eingeben und Tasten drücken** – realistische zeichenweise Ereignisse, die mit kontrollierten React-Eingabefeldern funktionieren
- **Gebundene Go-Methoden aufrufen** sowie Anwendungsereignisse auslösen und darauf warten

## Funktionsweise

Der MCP-Server wird nur dann in Ihre Anwendung kompiliert, wenn das **`mcp`-Build-Tag** vorhanden ist. Ohne dieses Tag fehlt der Servercode vollständig in der Binärdatei – kein Laufzeitaufwand, keine offenen Ports, keine Angriffsfläche.

Wenn das Tag vorhanden ist, startet der Server automatisch innerhalb von `App.Run()`, bindet sich standardmäßig an `127.0.0.1:9099` und protokolliert seinen Endpunkt. Es ist kein Benutzercode erforderlich.

## Tutorial

### Schritt 1 – eine normale Wails-Anwendung erstellen

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
Setzen Sie `WAILS_MCP=1`; die Wails CLI fügt dann automatisch das Tag `mcp` hinzu:

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

Der Server verwendet den **streamfähigen HTTP-Transport von MCP**. Stellen Sie die Verbindung mit einem beliebigen MCP-kompatiblen Client her.

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

Bitten Sie den Agenten, Ihre Anwendung zu testen. Hier sind einige Beispiel-Prompts:

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
| `WAILS_MCP` | (nicht gesetzt) | Setzen Sie den Wert auf `1`, `true`, `on` oder `yes`, damit bei Verwendung der Wails CLI automatisch das Build-Tag `mcp` hinzugefügt wird. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Schnittstelle, an die der Server gebunden wird. Bindungen an andere Adressen als Loopback erfordern `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | nicht gesetzt | Optionales Bearer-Token bei Loopback; für andere Bindungsadressen erforderlich. Clients senden `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Port, an dem der Server lauscht. Setzen Sie den Wert auf `0`, um einen zufällig zugewiesenen freien Port zu verwenden, der im Protokoll ausgegeben wird. |
| `WAILS_MCP_TIMEOUT` | `30000` | Standardmäßiges Zeitlimit für die JS-Auswertung in **Millisekunden**. |
| `WAILS_MCP_HIDE_CURSOR` | (nicht gesetzt) | Setzen Sie den Wert auf `1` oder `true`, um die animierte Cursor-Einblendung zu deaktivieren. |

Beispiel – benutzerdefinierter Port und ein Zeitlimit von 60 Sekunden:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Verfügbare Werkzeuge

| Werkzeug | Zweck |
| --- | --- |
| `app_info` | Anwendungsinformationen: Plattform, Architektur, alle Fenster, MCP-Endpunkt |
| `windows_list` | Alle Fenster mit Geometrie und Status auflisten |
| `window_control` | Fokussieren, Größe ändern, verschieben, Vollbildmodus, Entwicklertools, neu laden, URL festlegen, … (22 Aktionen) |
| `js_eval` | JavaScript in einem Fenster auswerten (asynchroner Funktionsrumpf, `return` für den Wert) |
| `dom_html` | HTML der Seite oder eines bestimmten Elements abrufen |
| `dom_query` | Elemente anhand eines CSS-Selektors suchen – Tag, Text, Abmessungen, Sichtbarkeit |
| `screenshot_dom` | Struktureller Snapshot der sichtbaren Seite (DOM-basiert, keine Pixel) |
| `mouse_move` | Den Cursor animiert zu einem Punkt oder CSS-Selektor bewegen |
| `mouse_click` | Mit dem animierten Cursor klicken (links/rechts/Mitte, Doppelklick, Zusatztasten) |
| `mouse_drag` | Mit dem animierten Cursor ziehen (unterstützt HTML5-Drag-and-drop-Elemente) |
| `mouse_scroll` | An einem Punkt oder Element scrollen |
| `keyboard_type` | Text mit realistischen Ereignissen Zeichen für Zeichen eingeben |
| `keyboard_press` | Eine einzelne Taste (Enter, Tab, Escape, ArrowDown, …) mit optionalen Modifikatortasten drücken |
| `call_bound_method` | Eine gebundene Methode eines Go-Dienstes aufrufen, z. B. `main.GreetService.Greet` |
| `emit_event` | Ein Wails-Anwendungsereignis auslösen |
| `wait_for_event` | Auf ein Wails-Anwendungsereignis warten und dessen Daten zurückgeben |

### Unterstützung mehrerer Fenster

Alle Werkzeuge, die auf ein Fenster einwirken, akzeptieren ein optionales Argument `window`, das den **Namen** des Fensters enthält (festgelegt über `WebviewWindowOptions.Name`). Wird es weggelassen, richtet sich das Werkzeug an das aktuell fokussierte Fenster oder an das erste Fenster, wenn keines fokussiert ist.

```
List all windows, then click the "New" button in the window named "editor".
```

### Elemente auswählen

Maus- und Tastaturwerkzeuge akzeptieren entweder:

- **CSS-Selektor** — `selector: "#submit-btn"` (das Element wird automatisch in den sichtbaren Bereich gescrollt)
- **Koordinaten** — `x: 400, y: 300` (CSS-Pixel relativ zum Viewport)

Stelle bei Ziehvorgängen `from_` und `to_` voran:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Sicherheit

@note{type="caution"}
Der MCP-Server ermöglicht die vollständige programmgesteuerte Kontrolle über deine Anwendung. Jeder, der auf seine Werkzeuge zugreifen kann, kann das DOM lesen, JavaScript auswerten, Schaltflächen anklicken und Go-Methoden aufrufen.

@end

- Der Server bindet standardmäßig an `127.0.0.1`. Browser-Ursprünge müssen HTTP(S)-Loopback-Ursprünge sein; opake (`null`), fehlerhafte und fremde Ursprünge werden abgelehnt.
- Native MCP-Clients ohne Header werden weiterhin unterstützt. Ohne `WAILS_MCP_TOKEN` wird lokalen Prozessen und zulässigen lokalen Browser-Ursprüngen vertraut; Ursprungsprüfungen sind keine Authentifizierung. Lege ein Token mit hoher Entropie fest, um für alle `/mcp`-Aufrufe eine Bearer-Authentifizierung zu verlangen. Konfiguriere dasselbe Token im Header `Authorization: Bearer <token>` des Clients. Preflight-Anfragen benötigen das Token nicht.
- Der Callback `/eval-result` verwendet statt des Bearer-Tokens des Clients unvorhersehbare IDs pro Auswertung, sodass die Übermittlung von Webview-Ergebnissen kompatibel bleibt.
- Produktions-Builds sollten das Tag `mcp` **nicht** enthalten. Die Wails CLI fügt es nur hinzu, wenn `WAILS_MCP=1` ausdrücklich festgelegt ist, und der standardmäßige `wails3 build` enthält keinen Servercode.
- Wenn du den Server über eine Nicht-Loopback-Schnittstelle bereitstellen musst (z. B. für Tests im LAN), lege `WAILS_MCP_HOST=0.0.0.0` und ein `WAILS_MCP_TOKEN` mit hoher Entropie fest; ohne Token schlägt der Start fehl. Verwende in nicht vertrauenswürdigen Netzwerken einen verschlüsselten Tunnel oder einen TLS-terminierenden Proxy, da der integrierte Listener HTTP verwendet.

## Beispielanwendung

Unter [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mcp) ist eine vollständige Playground-Anwendung verfügbar, die alle Werkzeuge demonstriert. Sie enthält:

- Zähler mit Schaltflächen zum Erhöhen und Zurücksetzen
- Namenseingabe mit den gebundenen Methoden Greet, Add und Shout
- Quelle und Ziel für HTML5-Drag-and-Drop
- Scrollbare Liste (50 Elemente)
- Ereignisprotokoll

Führe sie wie folgt aus:

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Verbinde anschließend Claude Code oder einen beliebigen MCP-Client mit `http://127.0.0.1:9099/mcp` und fordere ihn auf, die Benutzeroberfläche zu testen.

## Feedback

Der integrierte MCP-Server ist ein Experiment, und dein Feedback entscheidet über seine weitere Entwicklung. Wenn du ihn ausprobierst, möchten wir gern erfahren, welchen Client und welche Werkzeuge du verwendet hast, was du erwartet hast und was tatsächlich passiert ist und ob es hilfreich war, einen Agenten deine Anwendung steuern zu lassen — die nützlichsten Berichte beschreiben genau, was ausgeführt wurde. Teile uns deine Erfahrungen in der Diskussion [Feedback zum MCP-Server](https://github.com/wailsapp/wails/discussions/5692) mit.
