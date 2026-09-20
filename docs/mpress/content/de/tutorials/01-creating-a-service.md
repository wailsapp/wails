---
title: "QR-Code-Dienst"
description: "Erstellen Sie einen QR-Code-Dienst, um Wails-Dienste kennenzulernen"
slug: "tutorials/01-creating-a-service"
sourcePath: "tutorials/01-creating-a-service.md"
---

Ein **Dienst** in Wails ist eine Go-Struktur, die Geschäftslogik enthält, die Sie Ihrem Frontend zur Verfügung stellen möchten. Dienste strukturieren Ihren Code, indem sie zusammengehörige Funktionen bündeln.

Stellen Sie sich einen Dienst als Sammlung von Methoden vor, die Ihr JavaScript-Code aufrufen kann. Nach dem Generieren der Bindings kann jede öffentliche Methode des Dienstes vom Frontend aus aufgerufen werden.

In diesem Tutorial erstellen wir einen Dienst zur QR-Code-Generierung, um diese Konzepte zu veranschaulichen. Anschließend wissen Sie, wie Sie Dienste erstellen, Abhängigkeiten verwalten und Ihren Go-Code mit dem Frontend verbinden.

<br/>

@steps
### Datei für den QR-Code-Dienst erstellen
Erstellen Sie in Ihrem Anwendungsverzeichnis eine neue Datei namens `qrservice.go`:

```go {title="qrservice.go"}
package main

import (
    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}
```

**Was geschieht hier?**

- `QRService` ist eine leere Struktur, die unsere Methoden zur QR-Code-Generierung aufnehmen wird
- `NewQRService()` ist eine Konstruktorfunktion, die eine neue Instanz unseres Dienstes erstellt
- `Generate()` ist eine Methode, die Text und eine Größe entgegennimmt und anschließend den QR-Code als PNG-Byte-Array zurückgibt
- Die Methode gibt gemäß der Go-Konvention, Fehler als letzten Wert zurückzugeben, `([]byte, error)` zurück
- Für die eigentliche QR-Code-Generierung verwenden wir das Paket `github.com/skip2/go-qrcode`

 <br/>

### Dienst registrieren
Es genügt nicht, einen Dienst zu erstellen – wir müssen ihn bei der Wails-Anwendung **registrieren**, damit sie weiß, dass der Dienst vorhanden ist, und Bindings dafür generieren kann.

Die Registrierung erfolgt beim Erstellen Ihrer Anwendung in `main.go`. Übergeben Sie Ihre Dienstinstanzen an die Option `Services`:

```go {title="main.go" ins="7-9"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewService(NewQRService()),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Was geschieht hier?**

- `application.NewService()` umschließt Ihren Dienst, damit Wails ihn verwalten kann
- Wir rufen `NewQRService()` auf, um eine Instanz unseres Dienstes zu erstellen
- Der Dienst wird dem Slice `Services` in den Anwendungsoptionen hinzugefügt
- Wails durchsucht diesen Dienst nun nach öffentlichen Methoden, die dem Frontend zur Verfügung gestellt werden können

 <br/>

### Abhängigkeiten installieren
Wir haben in unserem Code auf das Paket `github.com/skip2/go-qrcode` verwiesen, es aber noch nicht heruntergeladen. Go muss diese Abhängigkeit kennen und sie in Ihr Projekt herunterladen.

Führen Sie diesen Befehl in Ihrem Projektverzeichnis im Terminal aus:

```bash
go mod tidy
```

**Was geschieht hier?**

- `go mod tidy` durchsucht Ihre Go-Dateien nach Importanweisungen
- Der Befehl lädt alle fehlenden Pakete wie `go-qrcode` herunter und fügt sie zu `go.mod` hinzu
- Außerdem entfernt er alle nicht mehr verwendeten Abhängigkeiten
- Dadurch enthält Ihr Projekt den gesamten Code, den es für eine erfolgreiche Kompilierung benötigt

Sie sollten eine Ausgabe sehen, die bestätigt, dass das QR-Code-Paket heruntergeladen und Ihrem Projekt hinzugefügt wurde.

 <br/>

### Bindings generieren
Damit Sie diese Methoden aus Ihrem Frontend aufrufen können, müssen wir Bindings generieren. Führen Sie dazu `wails generate bindings` im Stammverzeichnis Ihres Projekts aus.

@note{type="info"}
Wenn Sie diesen Befehl zum allerersten Mal in einem Projekt ausführen, analysiert der Binding-Generator Ihren Code und dessen Abhängigkeiten gründlich. Dies kann mitunter etwas länger als erwartet dauern; nachfolgende Ausführungen sind jedoch wesentlich schneller.

@end

Nach der Ausführung sollten Sie in Ihrem Terminal eine Ausgabe ähnlich der folgenden sehen:

```bash
 % wails3 generate bindings
 INFO  Processed: 337 Packages, 1 Service, 1 Method, 0 Enums, 0 Models in 740.196125ms.
 INFO  Output directory: /Users/leaanthony/myproject/frontend/bindings
```

Im Frontend-Verzeichnis sollte nun ein neues Verzeichnis namens `bindings` vorhanden sein:

```bash
frontend/
└── bindings
    └── changeme
        ├── index.js
        └── qrservice.js
```

@note{type="tip" title="Profi-Tipp"}
Wenn Sie Ihre Anwendung mit `wails3 build` erstellen, werden die Bindings automatisch für Sie generiert und auf dem neuesten Stand gehalten.

@end

 <br/>

### Bindings verstehen
Sehen wir uns die generierten Bindings in `bindings/changeme/qrservice.js` an:

```js {title="bindings/changeme/qrservice.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 /**
  * QRService handles QR code generation
  * @module
  */

 // eslint-disable-next-line @typescript-eslint/ban-ts-comment
 // @ts-ignore: Unused imports
 import {Call as $Call, Create as $Create} from "@wailsio/runtime";

 /**
  * Generate creates a QR code from the given text
  * @param {string} text
  * @param {number} size
  * @returns {Promise<string> & { cancel(): void }}
  */
 export function Generate(text, size) {
     let $resultPromise = /** @type {any} */($Call.ByID(3576998831, text, size));
     let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
         return $Create.ByteSlice($result);
     }));
     $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
     return $typingPromise;
 }
```

Wir sehen, dass die Bindings für die Methode `Generate` generiert wurden. Sowohl die Parameternamen als auch die Kommentare wurden beibehalten. Für die Methode wurde außerdem JSDoc generiert, um Ihrer IDE Typinformationen bereitzustellen.

@note{type="info"}
Sie müssen die generierten Bindings nicht vollständig verstehen, sollten aber wissen, wie sie funktionieren.

@end

Die Bindings stellen Folgendes bereit:

- Funktionen, die Ihren Go-Methoden entsprechen
- Automatische Konvertierung zwischen Go- und JavaScript-Typen
- Promise-basierte asynchrone Operationen
- Typinformationen in Form von JSDoc-Kommentaren

@note{type="tip" title="TypeScript"}
Der Binding-Generator unterstützt auch die Generierung von TypeScript-Bindings. Führen Sie dazu `wails3 generate bindings -ts` aus.

@end

Der generierte Dienst wird von einer Datei `index.js` erneut exportiert:

```js {title="bindings/changeme/index.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 import * as QRService from "./qrservice.js";
 export {
     QRService
 };
```

Anschließend können Sie über den vereinfachten Importpfad `./bindings/changeme` darauf zugreifen, der nur aus dem Pfad Ihres Go-Pakets besteht, ohne einen Dateinamen anzugeben.

@note{type="info"}
Vereinfachte Importpfade sind nur bei Verwendung von Frontend-Bundlern verfügbar. Wenn Sie ein Vanilla-Frontend ohne Bundler bevorzugen, müssen Sie entweder `index.js` oder `qrservice.js` manuell importieren.

@end

 <br/>

### Bindings im Frontend verwenden
Jetzt können wir unseren Go-Dienst aus JavaScript aufrufen! Die generierten Bindings machen dies einfach und typsicher.

Aktualisieren Sie `frontend/src/main.js`, um die neuen Bindings zu verwenden:

```js {title="frontend/src/main.js"}
 import { QRService } from './bindings/changeme';

 async function generateQR() {
     const text = document.getElementById('text').value;
     if (!text) {
         alert('Please enter some text');
         return;
     }

     try {
         // Generate QR code as base64
         const qrCodeBase64 = await QRService.Generate(text, 256);

         // Display the QR code
         const qrDiv = document.getElementById('qrcode');
         qrDiv.src = `data:image/png;base64,${qrCodeBase64}`;

     } catch (err) {
         console.error('Failed to generate QR code:', err);
         alert('Failed to generate QR code: ' + err);
     }
 }

 export function initializeQRGenerator() {
     const button = document.getElementById('generateButton');
     button.addEventListener('click', generateQR);
 }
```

**Was geschieht hier?**

- Wir importieren `QRService` aus den generierten Bindings
- `QRService.Generate()` ruft unsere Go-Methode auf – dieser Aufruf gibt ein Promise zurück, daher verwenden wir `await`
- Die Go-Methode gibt `[]byte` zurück, was Wails für JavaScript automatisch in einen Base64-String konvertiert
- Wir erstellen mit dem Base64-String eine Daten-URL, um das PNG-Bild anzuzeigen
- Der Block `try/catch` verarbeitet alle Fehler von der Go-Seite, beispielsweise ungültige Eingaben
- Wenn unser Go-Code einen Fehler zurückgibt, wird das Promise abgelehnt und wir fangen den Fehler hier ab

Aktualisieren Sie nun `index.html`, um die neuen Bindings in der Funktion `initializeQRGenerator` zu verwenden:

```html {title="frontend/src/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>QR Code Generator</title>
            <style>
                body {
                font-family: Arial, sans-serif;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                height: 100vh;
                margin: 0;
            }
                #qrcode {
                margin-bottom: 20px;
                width: 256px;
                height: 256px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
                #controls {
                display: flex;
                gap: 10px;
            }
                #text {
                padding: 5px;
            }
                #generateButton {
                padding: 5px 10px;
                cursor: pointer;
            }
            </style>
</head>
<body>
<img id="qrcode"/>
<div id="controls">
    <input type="text" id="text" placeholder="Enter text">
        <button id="generateButton">Generate QR Code</button>
</div>

<script type="module">
    import { initializeQRGenerator } from './main.js';
    document.addEventListener('DOMContentLoaded', initializeQRGenerator);
</script>
</body>
</html>
```

Führen Sie `wails3 dev` aus, um den Entwicklungsserver zu starten. Nach einigen Sekunden sollte sich die Anwendung öffnen.

Geben Sie Text ein und klicken Sie auf die Schaltfläche „QR-Code generieren“. In der Mitte der Seite sollte ein QR-Code angezeigt werden:

![QR-Code](/assets/qr1.png)

 <br/>

 <br/>

### Alternativer Ansatz: HTTP-Handler
Bislang haben wir die folgenden Bereiche behandelt:

- Erstellen eines neuen Dienstes
- Generieren von Bindings
- Verwenden der Bindings in unserem Frontend-Code

**Warum einen HTTP-Handler verwenden?**

Methodenbindungen eignen sich hervorragend für Datenoperationen, aber zum Bereitstellen von Dateien, Bildern oder anderen Medien gibt es einen alternativen Ansatz. Anstatt alles in Base64 zu konvertieren und über Bindungen zu senden, kann der Dienst als kleiner Webserver fungieren.

Dies ist nützlich, wenn:

- Bilder, Videos oder große Dateien bereitgestellt werden
- Standardmäßige HTML-Tags wie `<img>` oder `<video>` mit `src`-Attributen verwendet werden sollen
- Direkter URL-Zugriff auf Ressourcen erforderlich ist

Wenn der Dienst die Go-Standardmethode `ServeHTTP(w http.ResponseWriter, r *http.Request)` implementiert, kann Wails ihn als HTTP-Endpunkt zugänglich machen. Erweitern wir unseren QR-Code-Dienst um diese Unterstützung:

```go {title="qrservice.go" ins="4-5,37-65"}
package main

import (
    "net/http"
    "strconv"

    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Was hier geschieht:**

- `ServeHTTP` ist die Standardschnittstelle von Go zur Verarbeitung von HTTP-Anfragen
- Wir lesen Abfrageparameter aus der URL aus (`?text=hello&size=256`)
- Wir rufen unsere vorhandene Methode `Generate()` auf, um den QR-Code zu erstellen
- Wir setzen den Inhaltstyp auf `image/png`, damit Browser erkennen, dass es sich um ein Bild handelt
- Wir schreiben die unverarbeiteten PNG-Bytes direkt in die Antwort – Base64 ist nicht erforderlich!

Aktualisiere nun `main.go`, um die Route festzulegen, über die der QR-Code-Dienst erreichbar sein soll:

```go {title="main.go" ins="8-10"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                 Route: "/qrservice",
             }),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Was hier geschieht:**

- Wir fügen `application.ServiceOptions` hinzu, um zu konfigurieren, wie der Dienst bereitgestellt wird
- `Route: "/qrservice"` macht den HTTP-Handler unter `/qrservice` erreichbar
- Nun ruft jede Anfrage an `/qrservice?text=hello` unsere Methode `ServeHTTP` auf
- Wenn `Route` nicht festgelegt ist, ist die HTTP-Handler-Funktion deaktiviert

@note{type="info"}
Wenn die Option `Route` nicht ausdrücklich festgelegt wird, ist der HTTP-Handler vom Frontend aus nicht erreichbar.

@end

Aktualisiere abschließend `main.js`, um anstelle der Base64-Codierung ein einfaches Bild-`src` zu verwenden:

```js {title="frontend/src/main.js"}
async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Make the image source the path to the QR code service, passing the text
    img.src = `/qrservice?text=${encodeURIComponent(text)}`
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Was hier geschieht:**

- Wir haben den Import und den Aufruf von `await QRService.Generate()` entfernt
- Stattdessen setzen wir `img.src` einfach so, dass es auf unseren HTTP-Endpunkt verweist
- `encodeURIComponent()` maskiert Sonderzeichen in der URL sicher
- Der Browser sendet automatisch eine HTTP-GET-Anfrage, wenn wir `src` festlegen
- Für Bilder ist dies einfacher und effizienter – eine Base64-Konvertierung ist nicht erforderlich!

Wenn die Anwendung erneut ausgeführt wird, sollte derselbe QR-Code erscheinen:

![QR-Code](/assets/qr1.png)

 <br/>

 <br/>

### Dynamische Konfigurationen unterstützen
**Das Problem mit fest codierten Routen:**

Im obigen Beispiel haben wir in unserem JavaScript-Code die Route `/qrservice` fest codiert. Dadurch entsteht eine enge Kopplung zwischen der Go-Konfiguration und dem Frontend-Code.

Wenn du `main.go` bearbeitest und die Option `Route` änderst, ohne `main.js` zu aktualisieren, funktioniert die Anwendung nicht mehr:

```go {title="main.go" ins="3"}
        // ...
            application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                Route: "/services/qr",
            }),
        // ...
```

Fest codierte Routen eignen sich für einfache Anwendungen, machen den Code jedoch fragil und erschweren seine Wartung.

**Die Lösung: dynamische Konfiguration**

Methodenbindungen und HTTP-Handler können zusammenarbeiten! Über Bindungen können wir dem Frontend mitteilen, welche Route es verwenden soll. Dadurch wird die Konfiguration dynamisch und der fest codierte Pfad entfällt.

So funktioniert es:

1. Die Lebenszyklusmethode `ServiceStartup` wird beim Start der Anwendung ausgeführt
2. Wir speichern die in den Optionen konfigurierte Route
3. Wir fügen eine Methode `URL()` hinzu, die das Frontend aufrufen kann, um die korrekte Route abzurufen
4. Nun fragt das Frontend den Go-Dienst nach seiner Route, anstatt sie zu erraten

Implementiere zunächst die Schnittstelle `ServiceStartup` und füge eine neue Methode `URL` hinzu:

```go {title="qrservice.go" ins="4,6,10,15,23-27,46-55"}
package main

import (
    "context"
    "net/http"
    "net/url"
    "strconv"

    "github.com/skip2/go-qrcode"
    "github.com/wailsapp/wails/v3/pkg/application"
)

// QRService handles QR code generation
type QRService struct {
    route string
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// ServiceStartup runs at application startup.
func (s *QRService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.route = options.Route
    return nil
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

// URL returns an URL that may be used to fetch
// a QR code with the given text and size.
// It returns an error if the HTTP handler is not available.
func (s *QRService) URL(text string, size int) (string, error) {
    if s.route == "" {
        return "", errors.New("http handler unavailable")
    }

    return fmt.Sprintf("%s?text=%s&size=%d", s.route, url.QueryEscape(text), size), nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Was hier geschieht:**

- Wir haben ein Feld `route` hinzugefügt, um die konfigurierte Route aus `ServiceStartup` zu übernehmen und zu speichern
- `ServiceStartup(ctx, options)` wird beim Start der Anwendung aufgerufen – hier speichern wir die Route
- Die Methode `URL()` erstellt die vollständige URL mit Abfrageparametern
- Wenn keine Route konfiguriert ist, die Route also leer ist, geben wir einen Fehler zurück
- `url.QueryEscape()` codiert den Text sicher für die Verwendung in einer URL
- Diese Methode steht dem Frontend über Bindungen zur Verfügung

Aktualisiere nun `main.js`, sodass anstelle eines fest codierten Pfads die Methode `URL` verwendet wird:

```js {title="frontend/src/main.js" ins="1,11-12"}
import { QRService } from "./bindings/changeme";

async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Invoke the URL method to obtain an URL for the given text.
    img.src = await QRService.URL(text, 256);
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Was hier geschieht:**

- Wir importieren `QRService`, um die Bindungen wieder zu verwenden
- Anstatt `/qrservice` fest zu codieren, rufen wir `await QRService.URL(text, 256)` auf
- Der Go-Dienst erstellt die URL mit der korrekten Route und den korrekten Parametern
- Wenn du nun die Route in `main.go` änderst, verwendet das Frontend automatisch die neue Route
- Die Go-Konfiguration und der Frontend-Code müssen nicht mehr manuell synchronisiert werden!

Es sollte genauso wie das vorherige Beispiel funktionieren, doch eine Änderung der Dienstroute in `main.go` führt nicht mehr dazu, dass das Frontend nicht mehr funktioniert.

@note{type="info"}
Wenn eine Go-Methode einen Fehler ungleich nil zurückgibt, wird das Promise auf der JS-Seite abgelehnt, und await-Anweisungen lösen eine Ausnahme aus.

@end

 <br/>

@end
