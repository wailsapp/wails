---
title: "Streams — Interna"
description: "Funktionsweise des Stream-Transports, Gründe für seinen Aufbau, Bedeutung der Pufferkonstanten und noch nicht abgeschlossene Arbeiten"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

Referenz für alle — ob Mensch oder Agent —, die den Stream-Transport ändern. Die benutzerseitige API befindet sich unter [Streams](/guides/streams/); auf dieser Seite geht es um die zugrunde liegende Mechanik und die Gründe dafür, da mehrere Entscheidungen willkürlich erscheinen, solange nicht bekannt ist, welche Probleme sie vermeiden.

## Dateien

| Datei | Aufgabe |
| --- | --- |
| `v3/pkg/application/stream.go` | öffentliche API, `StreamConn`, `streamSink`, der Manager und seine Registry |
| `v3/pkg/application/stream_session.go` | ein Seitenaufruf in einem Fenster: ausgehende Warteschlange, Frame-Typen, Verbindungstabelle |
| `v3/pkg/application/stream_transport.go` | die beiden HTTP-Endpunkte, binäres Framing, Zusammensetzen von Chunks, Runtime-Präludium |
| `v3/pkg/application/stream_server.go` | nur `-tags server`: echte WebSocket-Senke |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | wählt den Client-Transport beim Bereitstellen des Bundles aus |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | der wie `WebSocket` aufgebaute Client |
| `v3/tests/stream-performance/` | Lasttest-Harness (`-upload`, `-reloads`, Durchläufe verschiedener Szenarien) |

## Die Struktur

Go→JS und JS→Go verwenden unterschiedliche Mechanismen, und diese Asymmetrie bestimmt das gesamte Design.

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

**Go→JS verwendet einen gehaltenen Poll.** Die Anfrage wartet, bis etwas ausgeliefert werden kann. Absichtlich gibt es weder ein Polling-Intervall noch einen adaptiven Mechanismus: Der Server hält die Anfrage, bis ein Frame vorliegt. Daher beträgt die Auslieferungslatenz bereits ~0, und ein clientseitiges Intervall könnte sie nur erhöhen. Frames, die während einer laufenden Antwort eintreffen, sammeln sich an und werden mit der nächsten übertragen. Dadurch wird die Umlaufzeit selbst zum Bündelungsfenster — es vergrößert sich bei steigender Last, ohne dass etwas es misst. Gemessen wurden 1.0 Frames pro Antwort bei 100/s, weiterhin 1.0 bei 5000/s und 3.4 bei 20000/s; zugleich *sinkt* die p99-Latenz mit steigender Rate.

**JS→Go verwendet gewöhnliche POST-Anfragen.** Sendevorgänge werden pro Verbindung über eine Promise-Kette serialisiert, da gleichzeitige Aufrufe von `fetch` die Reihenfolge nicht beibehalten und Go darauf angewiesen ist, die Sendevorgänge in ihrer ursprünglichen Reihenfolge zu beobachten. Frames, die sich hinter einer laufenden Anfrage ansammeln, werden im nächsten POST gebündelt. Go hängt den akzeptierten Frame oder das akzeptierte Präfix des Batches *vor* dem Antworten an den Posteingang der Verbindung an. Dadurch kann der Client nicht über Bytes hinaus fortfahren, die Go noch nicht in die Warteschlange gestellt hat.

**Pro Fenster ist genau ein Poll aktiv, der alle Verbindungen multiplext.** Dadurch ist die korrekte Reihenfolge konstruktionsbedingt gewährleistet — eine Warteschlange, ein einziger Entleerer und kein zweiter Auslieferungspfad, der den ersten überholen könnte. Dies umgeht außerdem das HTTP/1.1-Limit von sechs Verbindungen pro Host unter Windows, wo es sich um echte Chromium-Netzwerkanfragen an `http://wails.localhost` handelt.

## Gründe für diese konkreten Entscheidungen

Jede dieser Entscheidungen ist eine Narbe aus der Arbeit am Event-Transport. Wird eine davon entfernt, tritt ein gemessener Fehler erneut auf.

**Nichts im Go→JS-Pfad greift auf den Hauptthread zu.** `Send` hängt Daten unter einem Mutex an und kehrt zurück. Früher führten Events ihre Auswertung direkt aus, wenn sie vom Hauptthread emittiert wurden, während ein zuvor von einer Goroutine emittiertes Event noch in der Warteschlange stand — auf allen drei Plattformen waren 4.4 % der Events vertauscht. Bei einer einzigen Warteschlange mit einem einzigen Entleerer kann das nicht passieren.

**Unabhängig von der Größe greift nichts auf `evaluateJavaScript` zu.** Das Einfügen der Nutzdaten in den Auswertungsquelltext hält oberhalb eines plattformspezifischen Schwellenwerts Hostspeicher belegt: 11.6 GB unter macOS und 6.2 GB unter WebKitGTK bei 100 × 1 MB/s. Streams verwenden diesen Pfad nie. Daher verläuft der Durchlauf mit konstanter Byterate bei jeder Frame-Größe flach.

**Steuerdaten werden in Headern übertragen, nie im Body oder im Query-String.** WebKitGTK 6.0 kann bei benutzerdefinierten URI-Schemata POST-Bodys als Query-Parameter übergeben (`transport_http.go` enthält genau dafür einen Fallback), und WebView2 begrenzt die Body-Übertragung auf etwa 2 MB.

**Die Poll-Antwort ist binär, nicht JSON.** Frames sind `[]byte`; Base64 in einem JSON-Umschlag würde für jeden Frame 33 % zusätzlichen Aufwand sowie einen Parsing-Vorgang im UI-Thread verursachen.

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind` ist data / open / close / error. Es gibt weder eine Sequenznummer noch eine Bestätigung: Ein WebSocket führt keine erneute Wiedergabe durch, und beim Abbruch einer Verbindung gehen die gerade übertragenen Daten verloren. Dieses Verhalten nachzubilden ist einfacher und ehrlicher als ein Cursor, dessen Anforderungen der begrenzte Puffer nicht immer erfüllen könnte.

**Das Halten einer Anfrage ist sicher,** da jede Webview-Anfrage bereits eine eigene Goroutine erhält. `dispatchWorkers` in `assetserver_webview.go` ist mit einem Kommentar, der genau diesen Fall nennt, auf 0 festgelegt; bevor dieser Pool aktiviert werden könnte, müsste zunächst die Lebensdauer von Anfragen begrenzt werden.

## Pufferkonstanten

Alle befinden sich in `stream.go`. **Es sind Kompilierzeitkonstanten, keine Optionen** — es gibt weder `Options.Streams` noch eine Einstellung pro Stream. Um sie zu ändern, muss die Datei bearbeitet werden.

| Konstante | Wert | begrenzte Größe |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 MB | pro Fenster gepufferte Bytes, die auf die Abholung warten |
| `streamOutQueueDepth` | 256 | pro Fenster gepufferte Frames |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 MB / 8192 | anwendungsweit gepufferte ausgehende Daten |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 MB / 8192 | anwendungsweit auf `Receive` wartende eingehende Daten |
| `streamMaxConnections` | 256 | aktive Verbindungen sowie in einer Sitzung in die Warteschlange gestellte Schließvorgänge |
| `streamMaxConnectionsGlobal` | 4096 | aktive Verbindungen in der gesamten Anwendung |
| `streamOutCloseDepthGlobal` | 4096 | nicht zugestellte Schließbenachrichtigungen in der gesamten Anwendung |
| `streamMaxSessionsPerWindow` | 16 | Sitzungen, die ein Fenster vorhalten darf, bevor eine neuere Generation eine ältere ersetzen muss |
| `streamMaxSessions` | 1024 | Sitzungen in der gesamten Anwendung |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | wartende Steuerframes außer Close-Frames, pro Sitzung und anwendungsweit |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | unvollständige Uploads pro Sitzung und Teile in einem Upload |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 MB / 4096 | Chunk-Nutzdaten und Metadaten der Teile anwendungsweit |
| `streamMaxChunkIDLen` | 64 Byte | eine vom Client bereitgestellte Chunk-Set-Kennung |
| `streamMaxResponseBytes` | 1 MB | eine Poll-Antwort |
| `streamHoldTimeout` | 20 s | wie lange ein leerer Poll wartet |
| `streamSessionTTL` | 60 s | so lange kein Poll und keine aktiven Verbindungen ⇒ Sitzung ist beendet |
| `streamSessionGrace` | 10 min | so lange kein Poll bei aktiven Verbindungen ⇒ Sitzung ist beendet |
| `streamSessionSweep` | 20 s | wie oft die Bereinigungsroutine nach beendeten Sitzungen sucht |
| `streamMaxFrameBytes` | 64 MB | ein Frame in beliebiger Richtung |
| `streamMaxNameLen` | 256 Byte | ein registrierter oder angeforderter Stream-Name |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 MB | empfangene Frames, die `Receive` noch nicht entnommen hat |

### So wählen Sie die Werte

**`streamOutQueueDepth` ist bewusst nicht `eventQueueCapacity` (64).** Diese Konstante wurde für eine Warteschlange gemessen, die jeweils um einen Eval-Aufruf geleert wird und bei der eine größere Tiefe lediglich die Tail-Latenz erhöhte. Ein Poll leert die Warteschlange stapelweise, daher muss die Tiefe hier die während eines Roundtrips erzeugten Frames aufnehmen können – bei 5000 Frames/s und einem Roundtrip von 5 ms sind das etwa 25 Frames. 256 lässt Spielraum für eine Lastspitze, ohne den Produzenten anzuhalten.

**`streamOutQueueBytes` ist die tatsächlich entscheidende Grenze**, weil 256 Frames mit jeweils 1 MB insgesamt 256 MB ergeben. Sie begrenzt als letzte Schutzmaßnahme den Hostspeicher, wenn ein Frontend keine Daten mehr abholt.

Hier greifen zwei Regeln ineinander, wobei die zweite leicht versehentlich verletzt wird:

- Die Grenzen für Tiefe und Byte-Anzahl beschränken die *Akkumulation*.
- **Eine leere Warteschlange nimmt unabhängig von dessen Größe immer ein Frame an.** Wurde die Byte-Grenze ausnahmslos erzwungen, ließ sich ein Frame, das größer als diese Grenze war, überhaupt nicht senden: Die Wartebedingung konnte nie erfüllt werden, sodass `Send` dauerhaft blockierte und `TrySend` dauerhaft „voll“ meldete. Die Frame-Größe kann nicht immer vom Aufrufer bestimmt werden; eine Struktur mit einem `[]byte`-Feld wird auf die sich aus ihrem Inhalt ergebende Größe serialisiert.

**`streamMaxResponseBytes` existiert wegen Windows.** Der Response-Writer von WebView2 sammelt den gesamten Body im Speicher und übergibt ihn erst in `Finish`. Eine unbegrenzte Antwort führt dort daher zu einer unbegrenzten Speicherallokation. Eine Erhöhung verbessert den Windows-Durchsatz *nicht*: Messungen zeigen, dass der Windows-Engpass pro Byte und nicht pro Antwort entsteht. Über die Frame-Messreihe variiert die Anzahl der Antworten pro Sekunde um den Faktor 4, während der Durchsatz in MB/s konstant bei etwa 90 bleibt.

**Die Eingangsgrenze sorgt dafür, dass das Frontend wartet.** Auf dem Desktop meldet `deliver`, dass die Kapazität erschöpft ist, der Endpunkt antwortet mit `429` und der Client versucht dasselbe Frame oder den nicht angenommenen Rest des Batches mit begrenztem Backoff erneut zu senden. Dadurch wird kein Webview-Anfrage-Slot belegt, während der Handler aufholt. Im Servermodus wartet die Lese-Pump des Sockets und überlässt TCP die Rückstaukontrolle. Ohne diese Grenze könnte ein Handler, der `Receive` nur langsam aufruft, den Hostspeicher unbegrenzt anwachsen lassen.

**Steuerframes umgehen die Datengrenzen, unterliegen jedoch eigenen Lebenszyklusgrenzen.** Geht ein Datenframe durch Rückstau verloren, führt dies zu einer Verlangsamung. Geht dagegen eine Open-Bestätigung verloren, verbleibt das Frontend dauerhaft in `CONNECTING`; geht ein Close-Frame verloren, nimmt es weiterhin an, eine beendete Verbindung sei aktiv. Steuerframes außer Close-Frames besitzen deshalb eine eigene begrenzte Warteschlange, die von der für Close-Frames getrennt ist. So kann eine Häufung abgelehnter Open-Anfragen nicht die Kapazität verbrauchen, die eine angenommene Verbindung zum Melden ihres Endes benötigt. Jede Sitzung reserviert außerdem für jede angenommene Verbindung einen Close-Slot. Ist diese Kapazität belegt, erhält eine neue Open-Anfrage vor ihrer Registrierung ein Backpressure-Signal, nach dem sie erneut versucht werden kann.

**Den sitzungsbezogenen Grenzen entsprechen auch anwendungsweite Grenzen.** Ohne sie könnte jede angenommene Sitzung oder Verbindung gleichzeitig ihre vollständige lokale Kapazität belegen. Für ausgehende und eingehende Daten gilt daher jeweils ein separates, über Desktop- und Servertransporte hinweg gemeinsam genutztes Budget mit Grenzen von 256 MiB und 8192 Frames. Für aktive Verbindungen gilt ein Budget von 4096 Einträgen; für noch nicht zugestellte Close-Benachrichtigungen gilt ein zweites Budget derselben Größe. Ist ein gemeinsames Budget ausgeschöpft, tritt dasselbe blockierende Verhalten von `Send` beziehungsweise nicht blockierende Verhalten von `TrySend` ein wie bei einer lokalen Grenze. Bei jedem Leeren, Empfangen, Schließen, fehlgeschlagenen Schreiben und Herunterfahren wird die jeweilige Reservierung freigegeben.

Diese beiden Budgets sind bewusst voneinander getrennt. Es gibt also nicht ein einziges Kontingent, das eine Verbindung an ihr eigenes Close-Frame weitergibt. Jede Reservierung wird von genau einem Eigentümer freigegeben: Der Slot einer Verbindung durch das einmal ausgeführte `shutdown` und der Slot eines Close-Frames durch den Vorgang, der das Frame verwirft – entweder das Leeren oder der Abbau seiner Sitzung. Wird Eigentum zwischen zwei Beteiligten übertragen, muss dies atomar erfolgen. Eine frühere Revision, in der ein Close-Frame den Slot der Verbindung übernahm, ließ diesen dauerhaft belegt, wenn der Abbau zwischen dem Versuch, ein Close-Frame einzureihen, und dem Fehlschlagen dieses Versuchs ausgeführt wurde.

**Go-Frames übertragen das Eigentum; von JavaScript-Frames wird eine Momentaufnahme erstellt.** Go `Send` behält den Slice des Aufrufers, bis der Transport ihn geschrieben hat. Nach einem erfolgreichen Aufruf dürfen Aufrufer diesen Speicher daher weder verändern noch wiederverwenden. JavaScript `send()` kopiert veränderliche binäre Eingaben vor der Rückkehr und entspricht damit der Eigentumssemantik nativer WebSockets. Diese asymmetrische Regel vermeidet eine zweite vollständige Kopie des Frames in Go und sorgt zugleich dafür, dass sich die Browser-API erwartungsgemäß verhält.

**Das Senden mit JavaScript folgt dem WebSocket-Vertrag für die Pufferung.** `send()` kann nicht blockieren. Eine Anwendung kann Daten daher schneller in die Warteschlange stellen, als der Desktop-Anfragekanal sie annimmt, genauso wie sie einen nativen WebSocket überlasten kann. `bufferedAmount` umfasst jedes von diesem Socket vorgehaltene Byte und dient dem Aufrufer als Rückstausignal; die hostseitigen Warteschlangen bleiben unabhängig davon durch die oben genannten Grenzen beschränkt. Bei einem endgültigen Fehler, dem Schließen durch die Gegenstelle oder einem lokalen `close()` werden die vorgehaltenen Nutzdaten freigegeben. Lokales Schließen bricht außerdem eine Open- oder Datenanfrage ab, die auf `429` wartet, bevor es das reservierte Close-Steuerframe sendet. Rückstau bei der Annahme oder beim Empfänger kann den Socket daher nicht dauerhaft in `CLOSING` festhalten.

**Für die Zusammensetzung von Chunks gilt ein gemeinsames Kontingent für den Hostspeicher.** Jede Sitzung darf einen einzelnen Frame mit bis zu 64 MiB zusammensetzen, dieses Kontingent darf jedoch nicht mit der Anzahl aller zugelassenen Sitzungen multipliziert werden. Unvollständige und wiederholbare Chunk-Sätze teilen sich daher ein Budget von 128 MiB für zugelassene Nutzdaten. Beim Vervollständigen eines Satzes werden kurzzeitig sowohl seine Teile als auch der zusammenhängende zusammengesetzte Frame vorgehalten; selbst die Verdopplung dieses logischen Kontingents bleibt somit innerhalb der effektiven Speicherobergrenze von 256 MiB. Vorgehaltene Teile teilen sich außerdem ein Metadatenkontingent von 4096 Einträgen, damit winzige oder leere Chunks die Maps und die Slice-Verwaltung nicht vergrößern können, ohne sich dem Bytelimit zu nähern. Eine Anfrage, die eines der Kontingente überschreiten würde, erhält ein Backpressure-Signal, nach dem sie erneut versucht werden kann. Bei Zustellung, Ablehnung, Ablauf oder Beendigung der Sitzung werden sowohl die Bytes als auch die Teileinträge an die gemeinsamen Budgets zurückgegeben.

**Beim Polling werden nur behebbare Fehler erneut versucht.** Bei Netzwerkfehlern, Antworten wegen Anfragezeitüberschreitung (`408`), Early-Data-Antworten (`425`), Backpressure (`429`) und Serverfehlern (`5xx`) kommt exponentielles Backoff von 250 ms bis zu 5 Sekunden zum Einsatz. Andere `4xx`-Antworten weisen auf Protokoll- oder Besitzfehler hin und schließen die Streams der Seite sofort; `410` ist das reguläre Abschlusssignal für eine ausgemusterte Sitzung. Beim Schließen der letzten Verbindung wird ein laufender Poll oder Backoff-Timer abgebrochen. Wird während dieses Abbaus eine Verbindung geöffnet, startet sie eine Ersatz-Pollschleife.

**`streamSessionTTL` muss deutlich über `streamHoldTimeout`** bleiben, andernfalls würde eine Sitzung bereinigt, während ihr eigener Poll berechtigterweise wartet.

Bei der Abstimmung auf eine Arbeitslast mit vielen kleinen Nachrichten greift zuerst die Tiefenbegrenzung, bei großen Nutzdaten hingegen die Bytebegrenzung. Für eine typische Anwendung muss keine von beiden geändert werden – unter macOS bewältigen die Standardwerte 634000 Frames/s und 2100 MB/s.

## Lebenszyklus von Verbindungen und Sitzungen

Eine **Sitzung** entspricht dem einmaligen Laden einer Seite in einem Fenster und wird durch eine vom Client erzeugte ID identifiziert, etwa die `clientId` der Runtime. Sitzungen werden bei der zuerst eintreffenden Anfrage verzögert erstellt. Wenn eine Plattform das anfragende Fenster nicht identifizieren kann (`windowID == 0`), ist die Anzahl der Sitzungs-IDs weiterhin global begrenzt, ihre Generationen werden jedoch bewusst nicht verglichen: Sie können zu unabhängigen Browserclients mit voneinander unabhängigen Generationszählern gehören. Diese Sitzungen laufen durch Schließen oder per TTL ab, statt einander zu ersetzen.

Drei Mechanismen schließen Ressourcen, geordnet danach, wie schnell sie den jeweiligen Zustand erkennen:

1. **Ein neuerer Sitzungspoll für ein Fenster ersetzt ältere Generationen.** Beim Neuladen erhält die Seite eine neue Sitzungs-ID, und eine im `sessionStorage` dieses Fensters gespeicherte Generation wird erhöht. Derselbe Wert wird in `window.name` gespiegelt, das bei deaktiviertem Speicher Neuladungen überdauert, und ist an `performance.timeOrigin` (bei älteren Engines an `Date.now()`) gekoppelt, damit das Löschen beider Speicher den Zähler nicht wieder bei eins beginnen lässt. Jede Anfrage enthält die Sitzungs-ID und die Generation. Ein Poll mustert nur niedrigere Seitengenerationen aus, sodass eine durch die Serverplanung verzögerte Anfrage der vorherigen Seite nicht neuer als ihr Ersatz erscheinen kann. Wenn eine Richtlinie sowohl den Speicher als auch `window.name` blockiert, wird für die Reihenfolge ersatzweise die Seitenuhr verwendet; sie hängt dann davon ab, dass spätere Seiten einen späteren Zeitursprung erhalten. Der Manager behält pro Fenster eine Obergrenze der ausgemusterten Generationen bei. Dadurch kann eine bereits laufende Anfrage die alte Seite nicht neu erstellen, ohne dass jede historische Sitzungs-ID gespeichert werden muss. Die Verbindungen der vorherigen Sitzung werden sofort geschlossen.
2. **Beim Zerstören eines Fensters** werden alle Sitzungen dieses Fensters verworfen, wie bei `eventPayloadStore.dropWindow`.
3. **Die TTL-Bereinigung** erfasst alles Übrige – etwa einen abgestürzten Renderer oder einen Rechner im Ruhezustand.

Nur die ersten beiden Mechanismen mustern die Seitengeneration aus. Die TTL-Bereinigung entfernt die inaktive Sitzung, ohne die Obergrenze der ausgemusterten Generationen zu erhöhen: Eine Seite beendet das Polling, sobald ihre letzte Verbindung geschlossen wird, doch dieselbe weiterhin geladene Seite muss später einen weiteren Stream öffnen können. Eine tatsächlich ersetzte Generation bleibt gesperrt, weil der Poll der neueren Seite die Obergrenze erhöht, bevor die alte Sitzung entfernt wird.

**Apple WebViews melden abgebrochene Anfragen.** Unter macOS und iOS bricht der `stopURLSchemeTask`-Callback von WebKit den zugehörigen Anfragekontext ab, sodass ein Poll für eine Seite, von der weg navigiert wurde, sofort freigegeben wird. Die Registry verwendet die beibehaltene native Task-Identität als Schlüssel und entfernt Einträge, wenn sie bei der Anfrageverarbeitung geschlossen werden. Linux und Windows stellen in der aktuellen Bridge weiterhin keinen entsprechenden Callback für einen frühzeitigen Abbruch bereit; dort bleibt eine wartende Anfrage bestehen, bis die Haltefrist abläuft. Aufgrund von Regel 1 wird die *Verbindung* dennoch unverzüglich geschlossen. Ein Abbruch wird ansonsten unter Linux als `EPIPE` und unter Windows erst bei `Finish` sichtbar.

## Transportauswahl

`Stream(name)` fragt `window._wails.streamFactory` ab. Server-Builds installieren eine Implementierung, die ein echtes `WebSocket` zurückgibt; in WebView-Builds bleibt sie ungesetzt, sodass der Poll-Client verwendet wird.

Die Factory **muss** installiert werden, bevor irgendein Modulkörper ausgeführt wird, da generierte Bindings Streams auf Modulebene erstellen. `custom.js` kann dies nicht übernehmen – `loadOptionalScript` sendet eine HEAD-Anfrage und fügt anschließend ein `<script>`-Tag hinzu, sodass es viel zu spät eintrifft. Stattdessen wird die Factory dem ausgelieferten Runtime-Bundle vorangestellt (`stream_prelude_server.go`). Dies ist konstruktionsbedingt synchron: Abhängigkeiten von ES-Modulen werden vor den importierenden Modulen ausgewertet.

Wenn ein dritter Transport hinzugefügt wird, muss er ebenfalls in den Vorspann. Nicht wieder auf `custom.js` zurückgreifen.

## Was noch nicht fertiggestellt ist

|  | Status |
| --- | --- |
| Abbruch von Anfragen durch die Plattformschicht | **Apple fertiggestellt; Linux/Windows ausstehend** – siehe oben |
| Pufferkonstanten als Optionen | nicht fertiggestellt; nur zur Kompilierzeit |
| Typisierte Streams | bewusst nicht fertiggestellt – Frames sind per Entwurfsentscheidung `[]byte` |
| Pipelining (ein zweiter laufender Poll) | nicht fertiggestellt; würde eine geordnete Zusammensetzung in JS erfordern |
| Fairness pro Verbindung | nicht fertiggestellt – Verbindungen in einem Fenster teilen sich eine Warteschlange, sodass eine sie überflutende Verbindung ihre Nachbarverbindungen verlangsamt |
| Zusammenfassung von JS→Go-Frames | **fertiggestellt** – Frames, die sich hinter einer laufenden Anfrage ansammeln, werden in begrenzten Batches gesendet; gering ausgelastete Verbindungen senden weiterhin einen Frame pro POST |
| Durchsatz unter Windows | ~100 MB/s, begrenzt durch das `WebResourceRequested`-Marshalling. Gemeinsame Puffer (`PostSharedBufferToScript`) sind der mögliche Lösungsansatz; Bindings sind unter `internal/webview2/pkg/webview2/` vorhanden, aber noch nicht mit `pkg/edge` verbunden |
| `wails3 dev` / Vite | **funktioniert** – mit einem generierten `vanilla-js`-Projekt verifiziert: Der Vite-Entwicklungsserver leitet über `/` weiter, und `/wails/stream/*` wird vor dem Proxy von der Asset-Server-Middleware abgeglichen, sodass Streams unbeeinflusst bleiben |
| Mehrere Fenster | unter Last nicht getestet, obwohl Sitzungen konstruktionsbedingt auf ein Fenster beschränkt sind |

## Das Frontend-Paket im Entwicklungsmodus

Ein generiertes Projekt importiert `@wailsio/runtime` aus **npm** und nicht aus dem gebündelten `/wails/runtime.js`, das der Asset-Server ausliefert. Unter `wails3 dev` löst Vite es aus `node_modules` auf, sodass eine Anwendung, die gegen eine veröffentlichte Runtime gebaut wurde, clientseitige Ergänzungen aus einem Branch nicht erkennt.

Solange Streams noch nicht veröffentlicht sind, richte eine Test-App auf das Paket dieser Arbeitskopie aus:

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

Dadurch wird zuerst `dist/` neu erstellt, sodass stets die aktuellen Quellen installiert werden. Mache dies im selben Verzeichnis mit  
`npm install @wailsio/runtime@latest` rückgängig.

Beachte, dass es **zwei** Client-Build-Ausgaben gibt und leicht nur eine davon neu erstellt wird:  
`task v3:runtime:build:package` erzeugt das `dist/` des npm-Pakets (das vom Frontend einer App importiert wird), während `task v3:runtime:build:assets`  
`bundledassets/runtime.js` erzeugt (das der Webview vom Asset-Server lädt). Eine Änderung an  
`stream.ts` erfordert beide.

## Tests

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

Entscheidend ist der Test der Reihenfolge: Acht Goroutinen senden gleichzeitig anhand eines Zählers, der unter der Queue-Sperre ausgegeben wird, und die Reihenfolge beim Leeren muss exakt der Annahmereihenfolge entsprechen. Falls der Test jemals fehlschlägt, wurde die Invariante des einzelnen Leerers verletzt.

Lasttest-Harness:

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

Unter Windows muss er in der interaktiven Konsolensitzung ausgeführt werden — ein einfacher SSH-Aufruf wird in Sitzung 0 mit einer Ausgabe der Länge null beendet — und die Binärdatei muss an einem Ort bereitgestellt werden, auf den sowohl das SSH-Konto als auch das Konsolenkonto lesend zugreifen können, da `C:\Users\<user>` per ACL auf seinen Besitzer beschränkt ist.
