---
title: "Binding-System"
description: "Wie Go und JavaScript mit Wails v3 ohne Boilerplate-Code gegenseitig Funktionen aufrufen können"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> „Bindings“ bilden den **typsicheren Vertrag**, mit dem Sie Folgendes schreiben können:

```go
msg, err := chatService.Send("Hello")
```

in Go *und*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

in TypeScript, **ohne IPC-Verbindungscode manuell schreiben zu müssen**. Dieses Dokument erläutert, *wie* dies funktioniert: von der **statischen Analyse** zur Build-Zeit über die **Codegenerierung** bis zur **Laufzeit-Bridge**, die Bytes über die WebView überträgt.

> Eine maßgebliche ausführliche Beschreibung der Generator-Pipeline finden Sie unter [`contributing/architecture/bindings`](/contributing/architecture/bindings/). Diese Seite bietet
>
> einen Überblick
>
> für Mitwirkende.

---

## 1. 30-Sekunden-Überblick

| Phase | Komponente | Ausgabe |
| --- | --- | --- |
| **Erfassen/Analysieren** | `internal/generator/collect/`, `internal/generator/analyse.go` | In-Memory-Modell exportierter Go-Dienste, Methoden, Parameter, Rückgabetypen und Modelle |
| **Generierung** | `internal/generator/render/templates/*.tmpl` (`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | ES-Modul je Dienst unter `frontend/bindings/<full Go import path>/...` |
| **Laufzeit** | `pkg/application/messageprocessor*.go` und die eingebettete JS-Laufzeit unter `internal/runtime/desktop/@wailsio/runtime/src/` (`calls.ts`, `events.ts`, …) | Aufruf- und Ereignisnachrichten über die native Bridge der WebView |

Der Befehl `wails3 generate bindings` koordiniert den Ablauf und führt dabei `generator.Generate` (definiert in `internal/generator/generate.go`) für eine Gruppe von Go-Paketen aus.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. Statische Analyse

### Einstiegspunkt

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

Der Collector-Durchlauf durchläuft jedes geladene Paket und erfasst:

- `collect.ServiceInfo` – eine Instanz je exportierter, gebundener Go-Struktur.
- `collect.ServiceMethodInfo`/`collect.MethodInfo` – Signaturinformationen je Methode (Name, Parameter, Ergebnisse, Fehlerposition, Empfänger, Dokumentation).
- `collect.ModelInfo`/`collect.StructInfo` – werden als TS-/JS-Modelle ausgegeben.
- Direktivenkommentare wie `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore` und `//wails:id <hex>` (siehe `internal/generator/collect/directive.go`).

Nicht unterstützte Typen verursachen einen Generatorfehler, sodass Fehler bereits zur Build-Zeit statt erst zur Laufzeit sichtbar werden.

### Modellbezeichner

Der Aufruf-Umschlag der Laufzeit identifiziert eine Methode anhand eines **deterministischen FNV-1a-Hashes** ihres vollständig qualifizierten Namens (`pkg.Struct.Method`). In generierten Bindings erscheint dieser als `$Call.ByID(<numeric-id>, …)` oder als `$Call.ByName("pkg.Struct.Method", …)`, wenn die Generierung mit `-names` ausgeführt wird.

---

## 3. Codegenerierung

### Vorlagen

`internal/generator/render/templates/`:

| Vorlage | Zweck |
| --- | --- |
| `service.js.tmpl` | Ein JS-Modul je gebundenem Dienst |
| `service.ts.tmpl` | Zugehörige TypeScript-Datei (mit `-ts`) |
| `models.js.tmpl` | Ausgabe der Modellklassen (je Paket) |
| `models.ts.tmpl` | Ausgabe der Modell-`.d.ts` (je Paket) |
| `index.tmpl` | `index.{js,ts}`-Barrel-Re-Exporte je Paket |
| `eventcreate.js.tmpl`/`eventdata.d.ts.tmpl` | Ereigniskonstruktor / Typdefinitionen für Nutzdaten |
| `newline.tmpl` | Normalisierung des abschließenden Zeilenumbruchs |

Die Ausgabe wird unter `frontend/bindings/<full Go import path>/...` abgelegt. Beispielsweise wird ein in `github.com/you/yourapp/services/chat` definierter Dienst unter `frontend/bindings/github.com/you/yourapp/services/chat/` abgelegt. In v3 gibt es kein Verzeichnis `frontend/src/wailsjs/`.

### JavaScript-Ausgabe

Generierte Bindings sind ES-Module, die die Laufzeit-Hilfsfunktionen aus `/wails/runtime.js` importieren:

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

Wird die Generierung mit `-names` ausgeführt, wird stattdessen `$Call.ByName("pkg.Struct.Method", ...)` ausgegeben – stets **vollständig qualifiziert**, niemals nur `"Method"`.

Generierte Modellklassen verwenden ein `$$source`-Konstruktormuster mit `if (!("X" in $$source))`-Standardwerten je Feld, Feldnamen in Anführungszeichen und einer `static createFrom(...)`, die für String-Eingaben `JSON.parse` ausführt.

### Wichtige Typzuordnungen

Anhand von `internal/generator/render/` verifiziert:

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V` (`K` ist kein String-Typ) | `{ [_ in K]?: V }` (weder `Map<K, V>` noch `Record<K, V>`) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string` (JSON-ISO-8601) |
| `error` (Rückgabeposition) | abgelehntes Promise |

### Hinweis zur Reflection

`pkg/application/bindings.go` ist **handgeschrieben** und verwendet `reflect`, um den Methodenaufruf anhand einer `BoundMethod`-Registry zu steuern. Nehmen Sie ältere Aussagen wie „keine Reflection zur Laufzeit“ nicht zu wörtlich: Der Generator vermeidet Reflection, der Laufzeit-Dispatcher verwendet sie jedoch.

---

## 4. Protokoll für Laufzeitaufrufe

### JavaScript-Seite

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

Die Laufzeit-Hilfsfunktionen befinden sich in `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` (Aufrufweiterleitung), `events.ts` (Ereignisse) und zugehörigen Dateien – `invoke.ts` oder `errors.ts` gibt es in diesem Arbeitsbaum nicht. Das genaue Übertragungsformat wird auf der JS-Seite von `calls.ts` codiert und auf der Go-Seite von `pkg/application/messageprocessor_call.go` decodiert. Ziehen Sie beim Debuggen der Bridge beide Dateien gemeinsam heran.

### Go-Seite

1. `pkg/application/messageprocessor_call.go` empfängt die Aufrufnachricht.
2. Sucht in `pkg/application/bindings.go` die gebundene Methode anhand ihrer ID oder ihres Namens (gesteuert von `reflect`).
3. Ruft die gebundene Methode auf und serialisiert `{result, error}` zurück an JS.

### Fehlerzuordnung

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise` wird mit dem Ergebnis erfüllt |
| `error != nil` | `Promise` wird mit einem `Error` abgelehnt, dessen `message` die Go-Fehlerzeichenfolge enthält |

---

## 5. JavaScript aus Go aufrufen

Der Binding-Generator arbeitet nur in eine Richtung (Go-Methoden werden für JS bereitgestellt). Verwenden Sie für die Kommunikation von Go zu JS den Event-Bus oder führen Sie JS in einem Fenster aus:

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

Abonnieren Sie auf der JS-Seite mit `Events.On(name, cb)` aus `/wails/runtime.js`.

---

## 6. Erweitern und Fehlerbehebung

### Fehler wegen eines nicht unterstützten Typs

```
error: field "Client" uses unsupported type: chan struct{}
```

→ Kapseln Sie den Channel hinter einer Methoden-API oder markieren Sie das Feld mit `//wails:internal`, damit der Generator es überspringt.

### Veraltete Bindings

Die generierte Ausgabe wird bei jedem `wails3 generate bindings` / `wails3 dev` / `wails3 build` überschrieben. Wenn die IntelliSense der IDE veraltete Stubs anzeigt, löschen Sie `frontend/bindings/` und führen Sie den Generator erneut aus. Das Flag `-clean` (in aktuellen Builds standardmäßig `true`) leert das Binding-Verzeichnis vor jedem Durchlauf.

### Tipps zur Performance

- Streamen Sie keine großen Byte-Slices über die Bridge, sondern stellen Sie sie über den Asset-Server bereit.
- Fassen Sie mehrere schnelle Aufrufe in einer Methode zusammen, wenn die Latenz entscheidend ist.
- Verwenden Sie für kleine Parameter-Structs bevorzugt Wertempfänger, um Allokationen zu reduzieren.

---

## 7. Übersicht der wichtigsten Dateien

| Aufgabe | Datei |
| --- | --- |
| Generator-Orchestrierung | `internal/generator/generate.go` |
| Semantische Prüfungen | `internal/generator/analyse.go` |
| Erfassung (Services, Methoden, Modelle) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| Rendering-Templates | `internal/generator/render/templates/*.tmpl` |
| Speicherort der generierten Bindings | `frontend/bindings/<full Go import path>/...` |
| Dispatcher auf der Go-Seite | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| JS-Laufzeit | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

Halten Sie diese Kurzübersicht griffbereit, wenn Sie einen Bridge-Fehler untersuchen.

---

## 8. Zusammenfassung

1. **Collector** durchsucht Ihren Go-Code → semantisches Modell im Arbeitsspeicher.
2. **Templates** erzeugen ES-Module pro Service sowie Modell- und Indexdateien pro Paket.
3. **Message Processor** leitet Aufrufe auf der Go-Seite über die Binding-Registry weiter.
4. **JS Runtime** kapselt alles in idiomatische Promises mit Abbruchmöglichkeit.

Und das alles, ohne dass Sie eine einzige Zeile IPC-Boilerplate schreiben müssen. Das ist das Binding-System von Wails v3. Legen Sie los und erstellen Sie Bindings!
