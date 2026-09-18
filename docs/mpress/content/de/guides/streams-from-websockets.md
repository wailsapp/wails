---
title: "Migration eines WebSockets zu Streams"
description: "Schrittweise Umstellung einer bestehenden WebSocket-Implementierung auf Wails Streams, einschließlich der Unterschiede, die unbemerkt zu Fehlern führen"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

Eine mechanische Anleitung zur Umstellung einer App, die derzeit einen WebSocket-Server betreibt, auf [Streams](/guides/streams/). Sie ist dafür geschrieben, wörtlich befolgt zu werden, auch von einem Agenten.

Der Vorteil besteht darin, den Listener zu entfernen: kein gebundener TCP-Port, keine Origin-Prüfung, kein Token, nichts, was eine Firewall oder ein Endpoint-Security-Produkt beanstanden könnte. Die API ist ähnlich genug, dass der Großteil des Frontend-Codes unverändert bleibt — aber **drei Unterschiede führen unbemerkt zu Fehlern**. Sie werden zuerst aufgeführt, weil sie Sie sonst einen Nachmittag kosten werden.

## Der häufigste Fall: ein lokaler HTTP-Server als Behelfslösung

Der übliche Grund für einen WebSocket in einer Wails-App ist, dass es keine andere Möglichkeit gab, einen kontinuierlichen Datenstrom an das Frontend zu übertragen. Daher startet die App einen eigenen `http.Server` an einem lokalen Port, mit dem sich das Frontend anschließend verbindet. Wenn dies Ihrer Architektur entspricht, entfernt diese Migration den Server vollständig — und damit auch mehrere Komponenten, die Sie *um* ihn herum aufgebaut haben.

**Der Mechanismus zur Portermittlung entfällt.** Irgendetwas muss dem Frontend mitteilen, mit welchem Port es sich verbinden soll: ein gebundener `GetServerPort()`, ein fester Port mit einer Ausweichlösung, falls er belegt ist, eine injizierte globale Variable oder ein in `localStorage` gespeicherter Wert. All das entfällt — ein Stream wird über seinen Namen adressiert, und dieser Name ist auf beiden Seiten eine Konstante zur Kompilierzeit.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**Die CORS-Konfiguration entfällt.** Der Origin der Webview ist je nach Plattform `wails://` oder `http://wails.localhost`. Daher benötigt ein lokaler Server `CheckOrigin`, einen `Access-Control-Allow-Origin`-Header oder beides. Streams verwenden den Asset-Server, von dem die Seite geladen wurde, sodass keine ursprungsübergreifende Anfrage zugelassen werden muss.

**Jedes selbst entwickelte Authentifizierungstoken entfällt.** Ein an localhost gebundener Port ist für jeden Prozess auf dem Rechner erreichbar. Eine sorgfältige Implementierung fügt daher ein Token oder eine Nonce hinzu, um Verbindungen durch andere Software zu verhindern. Nun gibt es keinen erreichbaren Port mehr.

**Endpunkte, die keine WebSocket-Endpunkte sind, werden in die Middleware des Asset-Servers verschoben.** Diese Server bleiben selten auf WebSockets beschränkt — neben dem Socket sammeln sich häufig ein Dateidownload, ein Bildendpunkt und eine Zustandsprüfung an. Streams ersetzen diese Endpunkte nicht, aber Sie benötigen dafür auch keinen zweiten Server. Binden Sie dieselben Handler in den Asset-Server ein:

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

Das Frontend ruft anschließend `/api/...` als relative URL desselben Ursprungs auf — ohne Host, Port oder CORS. Mit Streams und dieser Änderung hat der lokale Server keine Aufgabe mehr.

## Lesen Sie dies, bevor Sie beginnen

### 1. `ev.data` ist ein `ArrayBuffer`, niemals ein String

Dies ist der wichtigste Unterschied. Ein WebSocket liefert Textnachrichten als Strings; ein Stream liefert jede Nachricht als Bytes. Code wie dieser **wird kompiliert und ausgeführt, verhält sich aber fehlerhaft**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Beheben Sie das Problem an der Schnittstelle statt an jeder Aufrufstelle:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Wenn Ihre Daten im JSON-Format vorliegen — was normalerweise der Fall ist — verwenden Sie alternativ `JSONStream` anstelle von `Stream` und umgehen Sie das Problem vollständig:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

Dies ist die kürzeste Migration für einen JSON-WebSocket: Tauschen Sie den Konstruktor aus und entfernen Sie die Aufrufe von `JSON.parse` und `JSON.stringify`; der übrige Handler bleibt unverändert. Legen Sie für Nicht-JSON-Verkehr einmalig einen Wrapper um den Stream und lassen Sie alle Handler unverändert — siehe [Kompatibilitäts-Wrapper](#kompatibilitts-wrapper).

### 2. Das Senden ist kompatibel, das Empfangen nicht

`send()` akzeptiert einen String und codiert ihn als UTF-8, sodass `s.send(JSON.stringify(x))` unverändert funktioniert. Nur der Empfangspfad muss bearbeitet werden. Diese Asymmetrie ist leicht zu übersehen, weil die Hälfte Ihres Codes weiterhin funktioniert.

### 3. Es gibt keine URL

Ein WebSocket überträgt Verbindungsparameter in seiner URL — Pfad, Abfragezeichenfolge, Subprotokoll und Authentifizierungstoken. Ein Stream hat nur einen Namen. Alles, was Sie über die URL übergeben haben, muss in den ersten Frame oder in eine gebundene Methode verschoben werden, die vor dem Verbindungsaufbau aufgerufen wird.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go-Seite

Entfernen Sie den HTTP-Server, den Upgrader und die Verbindungsregistrierung. Jedes dieser Elemente wird durch einen Handler ersetzt.

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Stream |
| --- | --- |
| `upgrader.Upgrade` / `websocket.Accept` | *(nichts — `HandleStream` übernimmt die vollständige Registrierung)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()` oder einfach den Handler beenden |
| Verbindungsregistrierung für Broadcasts | eine eigene beibehalten — siehe [Broadcasting](#broadcasting) |
| `http.ListenAndServe`, Mux, Origin-Prüfung, Token | **entfernen** |
| Ping/Pong-Keepalive | **entfernen** — es gibt keinen inaktiven Socket, der aktiv gehalten werden müsste |
| `r.Context()` | `c.Context()` |

Die Lebensdauer der Handler-Goroutine entspricht genau wie bei einem gorilla-Handler der Lebensdauer der Verbindung. Die Struktur Ihrer bestehenden Schleife kann daher unverändert übernommen werden.

## Frontend-Seite

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`, die vier Zustandskonstanten, `addEventListener`, `close(code, reason)` und `bufferedAmount` verhalten sich alle wie bei einem `WebSocket`.

### Kompatibilitäts-Wrapper

Wenn Sie die Handler gar nicht ändern möchten, kapseln Sie den Konstruktor einmal. Vorhandener Code, der String-Nachrichten erwartet, funktioniert dann unverändert:

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

Dann ist `const ws = TextStream("feed")` ein direkter Ersatz für `new WebSocket(url)`.

## Broadcasting

Ein WebSocket-Server führt normalerweise eine Registrierung, um Nachrichten an mehrere Empfänger zu verteilen. Streams bieten kein integriertes Broadcasting. Behalten Sie die Registrierung bei, speichern Sie darin aber `*StreamConn` statt `*websocket.Conn`:

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

Zwei Regeln sollten Sie beibehalten: Geben Sie die Sperre vor dem Senden frei und bevorzugen Sie bei einer Verteilung an mehrere Empfänger `TrySend`, damit ein einzelner langsamer Consumer nicht alle anderen Clients aufhält.

## Der andere Fall: Das Frontend kommuniziert direkt mit einem Broker

Das ist in einer Wails-App seltener, sollte Ihnen aber bekannt sein. Wenn das Frontend einen WebSocket **zu einem Broker statt zu Ihrer App** öffnet – `nats.ws` zu einem NATS-Server oder für MQTT über WebSocket –, ist ein Stream kein direkter Ersatz: Ein Stream verbindet das Frontend mit *Ihrem Go-Code* und nicht mit einem Drittanbieter.

Die Migration ist eine Architekturänderung und in der Regel eine sinnvolle:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Verschieben Sie den Broker-Client nach Go, wo die native Bibliothek besser als die Browserbibliothek ist, und stellen Sie die vom Frontend benötigten Teile über einen Stream bereit:

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

Beachten Sie `TrySend` im Subscription-Callback: Dieser Callback läuft in der Goroutine des Broker-Clients. Würde er blockiert, käme die Zustellung für jede Subscription dieser Verbindung zum Stillstand.

Ihre Vorteile: Die Broker-Anmeldedaten gelangen nie ins Frontend, auf dem Rechner wird kein WebSocket-Port zugänglich gemacht und der ausgereifte Go-Client übernimmt Wiederverbindung und Backoff statt des Browser-Clients.

## Migrationscheckliste

- [ ] `HandleStream` für jeden bisherigen WebSocket-Endpunkt registriert
- [ ] Leseschleife umgestellt: `ReadMessage`/`Read` → `c.Receive()`
- [ ] Schreibvorgänge umgestellt: `WriteMessage` → `c.Send()` oder `TrySend` bei jeder Verteilung an mehrere Empfänger und in jedem Broker-Callback
- [ ] HTTP-Server, Mux-Eintrag, Upgrader, Origin-Prüfung und Authentifizierungstoken **gelöscht**
- [ ] Ping/Pong-Keepalive **gelöscht**
- [ ] URL-Parameter in den ersten Frame oder eine gebundene Methode verschoben
- [ ] **Jeden gelesenen `ev.data`-Wert decodiert** – `new TextDecoder().decode(ev.data)` – oder den Shim übernommen
- [ ] Wiederverbindungslogik unverändert beibehalten (eine integrierte Wiederverbindung gibt es absichtlich nicht)
- [ ] Broker-Clients nach Go verschoben, falls das Frontend direkt mit einem Broker kommuniziert hat
- [ ] Auf `InitialHTML`-Fenster geprüft – sie können Streams überhaupt nicht verwenden

Nach diesen Elementen sollten Sie während der Umstellung suchen:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Zu erwartende Verhaltensunterschiede

|  | WebSocket | Stream |
| --- | --- | --- |
| Nachrichtentyp | Text oder Binärdaten | nur Bytes |
| `ev.data` | String oder `Blob`/`ArrayBuffer` | immer `ArrayBuffer` (außer bei `binaryType = "blob"`) |
| Standardwert für `binaryType` | `"blob"` | `"arraybuffer"` |
| Subprotokolle, `extensions` | ausgehandelt | nicht unterstützt; immer `""` |
| Verbindungsparameter | URL und Abfrageparameter | erster Frame oder ein gebundener Aufruf |
| Authentifizierung | Token oder Cookie | nicht erforderlich – die App ist der einzige Aufrufer |
| Keepalive | Ping/Pong | nicht erforderlich |
| Automatische Wiederverbindung | keine | keine (identisch) |
| Schließcodes | vollständiger Bereich | `1000` normal, `1001` Sitzung geschlossen, `1002` Framing stimmt nicht überein, `1006` Fehler |
| Rückstau | Socket-Puffer des Kernels | 8 MB / 256 Frames pro Fenster, danach blockiert `Send` |
| Viele Verbindungen, ein Endpunkt | ja | ja |

## Nach der Migration

Plausibilitätsprüfungen, die häufige Fehler aufdecken:

1. Laden Sie die Seite wiederholt neu – der Handler sollte jedes Mal beendet und ein neuer gestartet werden, ohne dass sich Handler ansammeln.
2. Senden Sie in jede Richtung eine Nachricht, die größer als 512 KB ist.
3. Lassen Sie die Verbindung länger als eine Minute inaktiv; der Datenverkehr sollte ohne erneuten Verbindungsaufbau fortgesetzt werden.
4. Öffnen Sie die Entwicklertools, halten Sie die Ausführung unter Last an einem Haltepunkt an und setzen Sie sie anschließend fort – der Producer sollte blockieren und sich danach wieder erholen, statt Daten zu verlieren oder unbegrenzt zu wachsen.
