---
title: "Wails erweitern"
description: "Praxisleitfaden zum Hinzufügen neuer Funktionen und Plattformen zu Wails v3"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails ist so konzipiert, dass es sich **leicht anpassen** lässt.
>
> Jedes wichtige Subsystem besteht aus Go-Code, den Sie lesen, ändern und ausliefern können.
>
> Diese Seite zeigt, *wo* Sie beginnen und *wie* Sie die Plattformkompatibilität wahren, wenn Sie:

- einen **Dienst** hinzufügen (Benachrichtigungen, KV-Speicher, benutzerdefiniertes IPC, …)
- einen **neuen CLI-Befehl** erstellen (`wails3 <foo>`)
- die **Runtime** erweitern (Fenster-API, Dialoge, Ereignisse)
- eine **Plattformfunktion** einführen (Wayland, …)
- die **plattformübergreifende Kompatibilität** wahren, ohne in `//go:build`-Tags zu ertrinken

---

## 1. Einen Dienst hinzufügen

Ein „Dienst“ ist in v3 ein vom Benutzer bereitgestellter Go-Typ, der über `application.Options.Services` registriert und über generierte Bindings für JS verfügbar gemacht wird. Die v3-Codebasis enthält:

- `internal/service/` — Grundgerüst für `wails3 generate service`:
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — einsatzbereite Dienste, die Sie direkt registrieren können (Benachrichtigungen, kvstore, sqlite, log, fileserver, dock, …).

Die in älteren Entwürfen als `internal/service/template/template.go` und `internal/generator/collect/services.go` bezeichneten Generator- und CLI-Dateien existieren nicht — das Grundgerüst wird von `internal/service/service.go` erzeugt (Einstiegspunkt: `service.Install`), und die Binding-Metadaten für Dienste werden in `internal/generator/collect/service.go` erfasst.

### 1.1 Dienst definieren

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 Lebenszyklus-Schnittstellen implementieren (optional)

Ein Dienst kann optional diese Schnittstellen erfüllen (aus `pkg/application`):

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **Wichtig:** `ServiceShutdown` akzeptiert **keine Argumente**. Eine Methode mit der
>
> Signatur `ServiceShutdown(ctx context.Context) error` erfüllt die Schnittstelle **nicht**
>
> und wird ohne Fehlermeldung niemals aufgerufen.

### 1.3 Dienst bei der Anwendung registrieren

Es gibt keinen globalen `services.Register(...)`-Aufruf. Dienste werden zur Laufzeit über `application.Options.Services` registriert:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

Nach der Registrierung erzeugt `wails3 generate bindings` unter `frontend/bindings/<your import path>/...` ES-Module, welche die exportierten Methoden umschließen.

### 1.4 Aufruf aus JS

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

In v3 gibt es kein globales `window.backend.*` — Aufrufe erfolgen über die generierten ES-Module, die ihrerseits `Call.ByID(...)` aus `/wails/runtime.js` aufrufen.

---

## 2. Einen neuen CLI-Befehl schreiben

Die v3-CLI verwendet **`github.com/leaanthony/clir`** (nicht cobra). Die Verdrahtung befindet sich in `v3/cmd/wails3/main.go`:

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

Es gibt keine auf `init()` basierende automatische Registrierung. Fügen Sie Ihren neuen Unterbefehl zu `cmd/wails3/main.go` hinzu, außerdem eine zugehörige Funktion in `internal/commands/` (und eine Flag-Struktur unter `internal/flags/`, falls er Optionen akzeptiert). Erstellen Sie die CLI neu:

```
cd v3
go install ./cmd/wails3
wails3 hello
```

Wenn Ihr Befehl die Taskfile-Integration benötigt, verwenden Sie die Hilfsfunktionen in `internal/commands/task_wrapper.go` (`wrapTask("yourtask", args)`) wieder.

---

## 3. Die Runtime ändern

Häufige Gründe:

- Neue Fensterfunktion (`SetOpacity`, `Shake`, …)
- Zusätzlicher Dialog (`ColorPicker`)
- API auf Systemebene (Bildschirmhelligkeit)

### 3.1 Öffentliche API

Fügen Sie die Methode zu `pkg/application/webview_window.go` hinzu (die Schnittstelle befindet sich in `window.go`):

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

Verwenden Sie die vorhandenen Hilfsfunktionen `InvokeSync`/`InvokeAsync`, damit der Aufruf im Hauptthread ausgeführt wird.

### 3.2 Nachrichtenprozessor

Wenn JS die neue Methode aufrufen muss, erweitern Sie die entsprechende `pkg/application/messageprocessor_*.go`-Datei. Der Nachrichtenprozessor verwendet Switch-basierte Methoden für `MessageProcessor`, keinen globalen `register(...)`-Aufruf:

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

In der Codebasis gibt es **keine** `messageprocessor_window_opacity.go`-Datei und kein auf `init()` basierendes `register(MsgSetOpacity, ...)`-Muster.

### 3.3 Plattformspezifische Implementierung

Fügen Sie die Implementierung zu jeder betriebssystemspezifischen Datei unter `pkg/application/` hinzu:

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

Wenn eine Plattform die Funktion nicht unterstützen kann, schreiben Sie einen No-op-Stub. Das Framework verfügt über keinen `ErrCapability`-Sentinel — dokumentieren Sie die Unterstützung und machen Sie sie bei Bedarf über das entsprechende boolesche Feld in `Options` beziehungsweise in der plattformspezifischen Optionsstruktur verfügbar.

### 3.4 Funktions-Flag (optional)

Das Paket `internal/capabilities/` dient zur Deklaration plattformspezifischer Funktionssätze. Es gibt keine öffentliche `application.HasCapability`-/ `application.CapOpacity`-API. Wenn eine Funktion zur Laufzeit überprüfbar sein soll, fügen Sie sie unter `internal/capabilities/` hinzu und stellen Sie über `pkg/application` einen typisierten Getter bereit.

---

## 4. Neue Plattformfunktionen hinzufügen

Beispiel: optionale Wayland-Unterstützung unter Linux.

1. Teilen Sie die betreffende `pkg/application/*_linux.go`-Datei in `*_linux_x11.go` (`//go:build linux && !wayland`) und `*_linux_wayland.go` (`//go:build linux && wayland`) auf.
2. Lassen Sie Benutzer die Funktion mit `wails3 build --tags wayland` aktivieren. Leiten Sie zusätzliche Tags über die vorhandene `EXTRA_TAGS`-Verdrahtung in `internal/commands/task_wrapper.go` weiter. Auf `dev`-Ebene gibt es kein `--tags wayland`-Flag — `wails3 dev` akzeptiert nur `--config`, `--port` und `-s`.
3. Aktualisieren Sie die Dokumentation und alle plattformspezifischen README-Dateien unter `pkg/application/`.

> Halten Sie die standardmäßigen Build-Tags auf ein Minimum beschränkt; verwenden Sie Opt-in-Tags nur für Nischenfunktionen.

---

## 5. Checkliste für plattformübergreifende Kompatibilität

| ✅ Schritt | Warum |
| --- | --- |
| Implementieren Sie in allen plattformspezifischen Dateien **jede** öffentliche Methode (auch nur als Stub). | Sorgt dafür, dass der Build auf jedem Betriebssystem erfolgreich bleibt |
| Dokumentieren Sie für jedes Betriebssystem die kontrollierte Einschränkung der Funktionalität. | Apps können anhand von `runtime.GOOS` verzweigen, ohne dass Fehler verborgen bleiben |
| Verwenden Sie zuerst **reines Go** und Cgo nur bei Bedarf. | Vereinfacht Cross-Compiles (unter Linux fällt der Cgo-Aufwand ohnehin an) |
| Führen Sie `task test:cli`, `task test:generator` und `task test:templates` aus. | Reproduziert die CI lokal |
| Dokumentieren Sie neue Build-Tags in der Dokumentation für Mitwirkende bzw. in der README der Vorlage. | Benutzer müssen über optionale Aktivierungen informiert sein |

---

## 6. Debug-Builds und Iterationsgeschwindigkeit

- Verwenden Sie `Options.LogLevel = slog.LevelDebug` (`Options.Logger = slog.Default()`), um ausführliche Runtime-Aktivitäten auszugeben. Eine Umgebungsvariable namens `WAILS_LOG_LEVEL` gibt es nicht.
- Die Flags für `wails3 dev` sind `--config`, `--port` und `-s`. Ein Flag namens `-race` oder `-verbose` gibt es nicht. Führen Sie den Race Detector mit `go test -race ./...` aus oder bauen Sie Ihre App mit `go build -race` und führen Sie sie direkt aus.
- Die Anleitung für Race-/Cgo-Tests befindet sich unter `v3/TESTING.md` (ältere Entwürfe verwiesen auf das nicht vorhandene `pkg/application/RACE.md`).

---

## 7. Beiträge zum Upstream-Projekt

1. Öffnen Sie für neue Funktionalität oder eine Änderung des öffentlichen Verhaltens einen PR im Entwurfsstatus für einen **WEP (Wails Enhancement Proposal)**, um die Idee und das Design zu besprechen. Erstellen Sie nur für einen reproduzierbaren Fehler oder ein Dokumentationsproblem ein Issue.
2. Verwenden Sie für die Implementierung die oben beschriebenen Verfahren.
3. Fügen Sie Folgendes hinzu:
  - Unit-Tests (`*_test.go`)
  - Dokumentation (diese Datei oder die entsprechende Seite unter `docs/...`)
  - Einen Regressionstest unter `internal/generator/testcases/`, falls Sie den Bindings-Generator geändert haben

4. Führen Sie vor dem Pushen `task precommit` und die relevanten `task test:*`-Targets lokal aus.

---

### Direktlinks

| Bereich | Speicherort |
| --- | --- |
| Integrierte Services | `pkg/services/` |
| Service-Scaffolder | `internal/service/` |
| CLI-Verdrahtung | `v3/cmd/wails3/main.go` |
| Implementierungen der CLI-Befehle | `internal/commands/` |
| Betriebssystemspezifische Runtime | `pkg/application/*_{darwin,linux,windows}.go` |
| Deklarationen von Fähigkeiten | `internal/capabilities/` |
| Taskfile-DSL | `v3/Taskfile.yaml` |
| Generator für Ereigniskonstanten | `v3/tasks/events/generate.go` |

---

Damit haben Sie nun einen **Fahrplan**, um Wails nach Ihren Vorstellungen anzupassen: Fügen Sie Services hinzu, verleihen Sie der CLI etwas Magie, modifizieren Sie die Runtime oder integrieren Sie völlig neue Betriebssystemfunktionen. Viel Freude beim Erweitern!
