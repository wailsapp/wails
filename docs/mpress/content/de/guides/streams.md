---
title: "Streams"
description: "Bidirektionale Byte-Streams zwischen Go und JavaScript mit dem Programmiermodell von WebSockets, jedoch ohne lauschenden Socket"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

Streams stellen einen benannten, geordneten, bidirektionalen Byte-Kanal zwischen Go und Ihrem Frontend bereit – mit demselben Programmiermodell wie ein WebSocket, jedoch **ohne einen TCP-Port zu binden**.

Ein WebSocket kann nicht über ein benutzerdefiniertes URL-Schema kommunizieren. Die einzige Möglichkeit, einen WebSocket innerhalb eines Webviews zu verwenden, besteht daher darin, einen echten HTTP-Server zu betreiben, der einen Port überwacht. In einer Desktop-Anwendung bedeutet dies einen offenen lokalen Port, den jeder andere Prozess auf dem Rechner erreichen kann. Für einen sicheren Betrieb sind deshalb eine Origin-Prüfung und ein Token erforderlich. Außerdem ist der Port für jede Firewall und jedes Endpoint-Security-Produkt sichtbar, das Ihre Benutzer einsetzen. Streams vermeiden all das: Sie nutzen den Asset-Server, den Ihre Anwendung bereits bereitstellt und der bereits an die Origin gebunden ist.

Migrieren Sie eine vorhandene WebSocket-Implementierung? Folgen Sie [WebSocket zu Streams migrieren](/guides/streams-from-websockets/). Die Anleitung ist für eine schrittweise mechanische Umsetzung ausgelegt und beginnt mit den drei Unterschieden, die unbemerkt zu Fehlern führen.

## Schnellstart

Deklarieren Sie einen Stream in Go. Der Handler wird einmal pro Verbindung in einer eigenen Goroutine ausgeführt:

```go
app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        _ = c.Send(process(frame))  // blocks like a socket write
    }
})
```

Stellen Sie vom Frontend aus über den Namen eine Verbindung her. Das Objekt implementiert die Schnittstelle `WebSocket`:

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` gibt genau wie `new WebSocket(url)` **synchron** mit `readyState === CONNECTING` zurück. Daher können Sie den Stream auf Modulebene erstellen:

```js
export const Telemetry = Stream("telemetry");
```

## Frames bestehen aus Bytes

Jeder Frame ist in Go ein `[]byte` und in JavaScript ein `ArrayBuffer`. Es werden weder ein Schema noch eine Codierung vorgegeben. Verwenden Sie nach Wunsch JSON, Protocol Buffers, CBOR oder Rohbytes.

Ein Frame ist eine **Nachricht und kein Byte-Stream**: Er kommt entweder vollständig oder gar nicht an, und seine Länge wird mitübertragen. Keine Seite muss die Größe im Voraus kennen. Eine Struktur mit einem Feld vom Typ `[]byte` wird daher in die sich ergebende Darstellung serialisiert und als einzelner Frame gesendet.

## Objekte senden

Frames bestehen aus Bytes, doch nur selten möchten Sie auf Byte-Ebene arbeiten. Beide Seiten bieten aufeinander abgestimmte JSON-Komfortfunktionen:

```go
type Reading struct {
    Sensor string  `json:"sensor"`
    Value  float64 `json:"value"`
}

app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    var cmd map[string]any
    if err := c.ReceiveJSON(&cmd); err != nil {
        return
    }

    _ = c.SendJSON(Reading{Sensor: "cpu", Value: 42.5})
})
```

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("telemetry");
s.onopen    = () => s.send({ subscribe: "cpu" });   // stringified for you
s.onmessage = (ev) => console.log(ev.data.value);   // already an object
```

`JSONStream` ist dasselbe Objekt wie `Stream`; die Codierung erfolgt lediglich an der Schnittstelle. Es gibt kein separates Protokoll, und ein Go-Handler kann keinen Unterschied feststellen. Ein Frame ohne gültiges JSON löst ein `error`-Ereignis aus und wird verworfen, statt die Verbindung zu beenden.

Verwenden Sie das einfache `Stream`, wenn Sie die Bytes benötigen, etwa für Protocol Buffers, CBOR, Binärformate oder alles, was Sie lieber selbst codieren möchten.

## Die Go-API

```go
// Register a handler. Runs once per connection, on its own goroutine.
func (a *App) HandleStream(name string, handler StreamHandler)

// The connection.
func (c *StreamConn) Send(data []byte) error        // blocks when the buffer is full
func (c *StreamConn) TrySend(data []byte) error     // ErrStreamFull instead of blocking
func (c *StreamConn) Receive() ([]byte, error)      // blocks until a frame or close
func (c *StreamConn) SendJSON(v any) error          // marshal and send as one frame
func (c *StreamConn) ReceiveJSON(v any) error       // receive one frame and unmarshal
func (c *StreamConn) Context() context.Context      // cancelled on disconnect
func (c *StreamConn) Window() Window                // nil in server mode
func (c *StreamConn) Name() string
func (c *StreamConn) Close() error
```

**Die Lebensdauer der Handler-Goroutine entspricht der Lebensdauer der Verbindung.** Sobald der Handler zurückkehrt, wird die Verbindung geschlossen. Blockieren Sie daher so lange wie gewünscht auf `Receive` (oder auf `c.Context()`), um sie offen zu halten. Dies entspricht der Struktur eines `gorilla`/`coder`-WebSocket-Handlers.

Als Fehler treten `ErrStreamClosed` (die Gegenstelle ist nicht mehr erreichbar) und `ErrStreamFull` (nur von `TrySend`) auf.

## Die JavaScript-API

`Stream(name)` gibt ein Objekt zurück, das die nützliche Teilmenge von `WebSocket` implementiert:

| unterstützt | Hinweise |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`, `onmessage`, `onclose`, `onerror` | zusätzlich `addEventListener` |
| `send(data)` | Zeichenfolge, `ArrayBuffer`, typisiertes Array oder `Blob` – siehe Eigentumsverhältnisse weiter unten |
| `JSONStream(name)` | dasselbe Objekt; Ein- und Ausgabe von Objekten |
| `close(code, reason)` |  |
| `binaryType` | **ist standardmäßig `"arraybuffer"`**, nicht `"blob"` |
| `bufferedAmount` | von `send` eingereihte Bytes, die Go noch nicht erreicht haben |
| `protocol`, `extensions` | immer `""` – wird nicht ausgehandelt |

Der Standardwert von `binaryType` ist die einzige bewusste Abweichung vom Standard: Frames sind immer binär, und ein `Blob` würde zum Lesen jeder Nachricht einen zusätzlichen asynchronen Zwischenschritt erzwingen. Setzen Sie den Wert auf `"blob"`, wenn Sie das Standardverhalten wünschen.

Mehrere Verbindungen mit demselben Stream-Namen sind zulässig, sowohl aus einem als auch aus mehreren Fenstern. Jede erhält einen eigenen `StreamConn` und eine eigene Handler-Goroutine.

**Die Eigentumsverhältnisse für Puffer unterscheiden sich je nach Richtung.** JavaScript-`send()` erstellt synchron eine Momentaufnahme veränderlicher binärer Eingaben und entspricht damit dem Verhalten nativer WebSockets. Aufrufer können diese Eingaben daher wiederverwenden, sobald `send()` zurückkehrt. Go-`Send` überträgt das Eigentum an seinem Slice auf den Transport und kopiert es nicht. Verändern oder verwenden Sie dieses Slice nach einem erfolgreichen Aufruf nicht erneut. Übergeben Sie ein neues Slice, wenn der Produzent seinen Speicher wiederverwenden muss.

## Lebenszyklus

Ein Stream verhält sich wie ein Socket. Ereignisse, die einen Socket schließen, schließen daher auch einen Stream:

| Ereignis | Auswirkung |
| --- | --- |
| Neuladen oder Navigation der Seite | Die Verbindung wird geschlossen, `Receive` des Handlers gibt einen Fehler zurück, und die neue Seite stellt eine neue Verbindung her. |
| `window.close()` / Fenster zerstört | Alle Verbindungen dieses Fensters werden geschlossen. |
| `s.close()` in JS | `Receive` des Handlers gibt `ErrStreamClosed` zurück. |
| Handler kehrt zurück | Das Frontend empfängt `onclose`. |
| Beenden der Anwendung | Der Kontext jeder Verbindung wird abgebrochen. |

Es gibt **keine automatische Wiederherstellung der Verbindung**; dies entspricht `WebSocket`. Falls Ihre Anwendung eine benötigt, funktioniert Ihre vorhandene WebSocket-Logik zur Wiederherstellung der Verbindung unverändert: Erstellen Sie den Stream in `onclose` neu.

## Rückstau

`Send` blockiert, wenn das Frontend nicht Schritt gehalten hat, so wie ein Schreibvorgang auf einem Socket bei vollem Sendepuffer blockiert. Verwenden Sie `TrySend`, wenn Sie Daten lieber verwerfen als warten möchten:

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

Ein pausiertes Frontend – etwa durch einen Haltepunkt in den Entwicklertools, ein ausgeblendetes Fenster oder App Nap – ruft keine weiteren Daten ab. Sobald die Puffergrenze erreicht ist, wird der Produzent blockiert. Das ist beabsichtigt: Dadurch wird der Speicherverbrauch begrenzt, statt einen ungelesenen Stream unbegrenzt wachsen zu lassen.

## Servermodus

Beim Erstellen mit `-tags server` wird der Transport durch einen **echten WebSocket** unter `/wails/stream/ws` ersetzt, da im Servermodus bereits ein Listener vorhanden ist, über den das Upgrade erfolgen kann. Der Go-Handler und der Frontend-Code bleiben identisch – in Ihrer App ändert sich nichts. Die Runtime wählt den Transport für Sie aus, bevor Modulcode ausgeführt wird. WebSocket-Verbindungen verwenden standardmäßig denselben Origin. Ein Server, der sein Frontend bewusst unter einem anderen vertrauenswürdigen Origin bereitstellt, kann diesen Host mit `ServerOptions.WebSocketOriginPatterns` hinzufügen.

## Performance

Gemessen mit `v3/tests/stream-performance`, ohne Drosselung, 0 verworfenen und 0 umsortierten Frames bei insgesamt ~41 Millionen Frames:

|  | Go→JS-Spitzenwert | JS→Go-Spitzenwert |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

Das Profil ist wichtiger als die Spitzenwerte:

- **Go→JS ist bei kleinen Frames deutlich schneller** – 634000 Frames/s unter macOS gegenüber ~6200/s in der Gegenrichtung. Eine einzelne Antwort fasst bis zu 256 Frames zusammen. JS→Go bündelt ebenfalls Frames, die sich hinter einer laufenden Anfrage ansammeln, doch jede Verbindung serialisiert weiterhin ihre eigene POST-Kette. Wenn Sie viele kleine Nachrichten senden, bevorzugen Sie Go→JS oder bündeln Sie sie auf Anwendungsebene, bevor Sie sie an Go senden.
- **Unter Windows sind 512 KB die optimale Größe für Uploads.** Größere Frames werden auf mehrere Anfragen aufgeteilt, und ein Frame mit 4 MB ist laut Messung *langsamer* als einer mit 512 KB.
- **Die Latenz ist niedrig und bleibt niedrig**: ~1–2 ms p99 unter macOS; sie verschlechtert sich auch unter Last nicht – bei 20000 Frames/s wurde eine *niedrigere* p99 gemessen als bei 100 Frames/s.

Vollständige Tabellen für jede Plattform und die Messmethode finden Sie im Messprotokoll zu dieser Funktion.

## Grenzwerte

|  | Grenzwert | Verhalten beim Erreichen des Grenzwerts |
| --- | --- | --- |
| Pro Fenster gepuffert, Abholung ausstehend | 8 MB oder 256 Frames, je nachdem, was zuerst erreicht wird | `Send` blockiert; `TrySend` gibt `ErrStreamFull` zurück |
| Anwendungsweit gepuffert, Abholung/Schreibvorgang ausstehend | 256 MB oder 8192 Datenframes | wie oben |
| Pro Verbindung empfangen, `Receive` ausstehend | 8 MB oder 256 Frames | Das `send()` des Frontends wird automatisch wiederholt, bis der Handler aufgeholt hat |
| Anwendungsweit empfangen, `Receive` ausstehend | 256 MB oder 8192 Frames | wie oben |
| Verbindungen pro Fenster | 256 | Der Öffnungsversuch wird automatisch wiederholt, bis ein Platz frei wird |
| Aktive Verbindungen in der gesamten Anwendung | 4096 | wie oben |
| Sitzungen pro Fenster | 16 | Beim Neuladen wird die eigene ältere Sitzung ersetzt; andernfalls wird der Öffnungsversuch wiederholt |
| Ein einzelner Frame in beliebiger Richtung | 64 MB | Go gibt `ErrStreamTooLarge` zurück; von JS aus löst der Stream `error` aus und wird geschlossen |
| Ein Streamname | 256 UTF-8-Bytes | Der Öffnungsversuch wird abgelehnt und der Stream löst `error` aus |
| Haltezeit einer inaktiven Abfrage | 20 s | Die Abfrage liefert ein leeres Ergebnis zurück und die Runtime sendet sie sofort erneut |

Keiner der obigen Fälle verwirft Daten unbemerkt. Die beiden Zeilen mit dem Hinweis *wird automatisch wiederholt* beschreiben gewöhnlichen Gegendruck: Die Runtime hält den Frame zurück und wiederholt den Versuch nach einer kurzen Wartezeit. Ihr Code sieht daher einen langsameren Stream statt eines Fehlers. Die Zeilen, in denen `error` ausgelöst wird, betreffen Programmierfehler und keine Lastprobleme; diese Fehler werden offengelegt, statt kaschiert zu werden.

Derzeit handelt es sich dabei um Konstanten zur Kompilierzeit, nicht um Optionen. Wenn Sie sie ändern müssen, lesen Sie den Leitfaden zu den Interna.

## Wann Sie keinen Stream verwenden sollten

- **Verwenden Sie für Anfrage-Antwort-Abläufe Bindings.** Streams sind für kontinuierliche oder unaufgeforderte Daten vorgesehen; ein Aufruf, der einen Wert zurückgibt, lässt sich einfacher als gebundene Methode umsetzen.
- **Verwenden Sie für Anwendungsereignisse `Emit`/`On`.** Ereignisse werden an alle Listener verteilt und bilden ein separates, bewährtes System. Streams sind Punkt-zu-Punkt-Verbindungen.
- **Nicht für `InitialHTML`-Fenster verfügbar.** Diese werden mit `origin === "null"` geladen und können daher überhaupt nicht auf den Asset-Server zugreifen.
