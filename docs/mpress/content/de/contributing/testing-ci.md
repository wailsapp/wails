---
title: "Tests und Continuous Integration"
description: "Wie Wails v3 mit Unit-Tests, Integrationstests, Race Detection und CI über GitHub Actions die Qualität sicherstellt."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Robuste Desktop-Frameworks erfordern absolut zuverlässige Tests. Wails v3 verfolgt eine **mehrschichtige Strategie**:

| Ebene | Ziel | Werkzeuge |
| --- | --- | --- |
| Unit-Tests | Schnelles Feedback für isolierte Funktionen | `go test ./...` |
| Generator-/CLI-Tests | `wails3 generate bindings` und die CLI-Anbindung validieren | `task test:generator`, `task test:cli` |
| Vorlagentests | Sicherstellen, dass jede ausgelieferte Vorlage weiterhin gebaut werden kann | `task test:templates` |
| Race Detection | Datenrennen in Runtime und Bridge erkennen | `go test -race ./...` |
| CI-Matrix | Verlässlichkeit unter verschiedenen Betriebssystemen bei jedem PR | GitHub Actions |

Dieses Dokument erläutert, **wo sich die Tests befinden**, **wie sie ausgeführt werden** und **welche Aufgaben das Taskfile orchestriert**.

> Der in älteren Entwürfen als `pkg/application/RACE.md` bezeichnete Leitfaden zur Race Detection
>
> befindet sich heute unter `v3/TESTING.md`.

---

## 1. Verzeichniskonventionen

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

Richtlinien:

- **Unit-Tests neben dem Code ablegen** (`foo.go` ↔ `foo_test.go`).
- Verwenden Sie für `pkg/`-Pakete den **Black-Box-Stil** (`package application_test`), wenn dies die API-Hygiene verbessert.
- Gemeinsam genutzte Fixtures befinden sich dort, wo sie verwendet werden. In diesem Arbeitsbaum gibt es kein zentrales `internal/testutil/`-Paket; binden Sie Hilfsfunktionen stattdessen paketweise ein.

---

## 2. Unit-Tests

### Tests schreiben

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Empfehlungen:

- Verwenden Sie [`stretchr/testify`](https://github.com/stretchr/testify) – bereits in `go.mod` enthalten.
- Bevorzugen Sie **tabellengesteuerte** Tests, wenn mehrere Eingaben, Grenzfälle oder erwartete Ergebnisse dasselbe Verhalten prüfen. Geben Sie jedem Fall einen aussagekräftigen Namen.
- Kapseln Sie plattformspezifische Besonderheiten bei Bedarf mithilfe von Build-Tags (`foo_windows_test.go`, `foo_darwin_test.go`, …) in Stubs.

### Erwartete Testabdeckung

Für neue und geänderte Logik wird eine Go-Anweisungsabdeckung von 100 % erwartet. Messen Sie das geänderte Paket, statt sich auf einen Prozentsatz für das gesamte Repository zu verlassen:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Einige Pfade lassen sich in einer normalen Testumgebung berechtigterweise nicht testen, beispielsweise rein plattformspezifische Fehler, hardwareabhängiges Verhalten oder ein defensiver Rückfallpfad, der sich nicht gefahrlos auslösen lässt. Begrenzen Sie solche Ausnahmen eng und erläutern Sie jeden nicht abgedeckten Pfad in der PR-Beschreibung.

### Lokal ausführen

```bash
cd v3
go test ./... -cover
```

Alternativ über das Taskfile (tatsächlich vorhandene Ziele – es gibt kein Kürzel `task test`):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. Integrationstests

`v3/tests/` enthält das paketübergreifende Integrations-Testsystem. Die Taskfile-Ziele `test:example:*` und `test:examples:*` führen Build- und Startprüfungen unter darwin / windows / linux aus, einschließlich Docker-basierter GTK3-/GTK4-Matrizen unter Linux.

> Ausführbare Beispiele befinden sich in `v3/examples/`. Die Testziele wählen die
>
> für die Hostplattform oder CI-Matrix geeigneten Beispiele aus und bauen sie.

Führen Sie die Smoke-Test-Suite für die Hostplattform wie folgt aus:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Erkennung von Datenrennen

Datenrennen sind in GUI-Runtimes fatal.

### Leitfaden zur Race Detection

Informationen zu folgenden Themen finden Sie unter `v3/TESTING.md`:

- Bekannte harmlose Datenrennen und Begründungen für deren Unterdrückung
- Interpretation von Stacktraces über Cgo-Grenzen hinweg (Linux GTK + WebKit2GTK)

### Lokale Race-Test-Suite

```
go test -race ./...
```

> `wails3 dev` besitzt kein Flag `-race` – seine CLI-Flags sind `--config`, `--port`
>
> und `-s` (aktiviert HTTPS). Um die Runtime unter dem Race Detector zu testen,
>
> bauen Sie mit `go build -race` eine Testanwendung und führen Sie sie direkt aus.

---

## 5. GitHub-Actions-Workflows

Tatsächlich vorhandene Workflow-Dateien unter `.github/workflows/` (anhand des Arbeitsbaums geprüft):

| Datei | Zweck |
| --- | --- |
| `build-and-test-v3.yml` | Hauptmatrix für Build und Tests von v3. Verwendet `actions/setup-go@v5` mit `go-version: 1.25`. Führt `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples` (und `BUILD_TAGS=gtk4 task test:examples` für den GTK4-Pfad), `task generator:test:check`, `task install` und anschließend `wails3 build` als Smoke-Test aus. Linux-Jobs installieren `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk` und führen die Test-Suite unter `dbus-run-session -- xvfb-run` aus. |
| `cross-compile-test-v3.yml` | Plausibilitätsprüfungen der Cross-Kompilierung |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Changelog-Automatisierung |
| `nightly-release-v3.yml` | Nächtliche v3-Release-Artefakte |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Verwaltung von WebView2-Abhängigkeiten und -Releases |
| `build-cross-image.yml` | Erstellt das Container-Image für den Cross-Compiler |
| `publish-npm.yml` | Veröffentlicht die eingebettete JS-Laufzeit von `@wailsio/runtime` auf npm |
| `pr-master.yml` | PR-seitige Prüfungen für den Branch `master` |
| `semgrep.yml` | Statische Analyse mit Semgrep |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Repository-Pflege und v2-seitige Abläufe |

In diesem Arbeitsbaum gibt es **keine** `qodana.yaml` und **keine** `runtime.yml` — ältere Entwürfe dieser Seite nannten beide, aber nur `semgrep.yml` deckt die statische Analyse ab, und das JS-Laufzeitpaket wird über `publish-npm.yml` veröffentlicht.

Die CI-Schritte entsprechen den obigen Taskfile-Zielen (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), sodass Sie die CI lokal eins zu eins reproduzieren können. Der Smoke-Schritt `wails3 build` in `build-and-test-v3.yml` wird **ohne zusätzliche Flags** aufgerufen — für `wails3 build` gibt es kein Flag `-skip-package`.

---

## 6. Lokale Übereinstimmung mit der CI

Es gibt kein einzelnes übergeordnetes Ziel `task ci`. Reproduzieren Sie die CI, indem Sie die tatsächlichen Ziele verketten:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Fehlerbehebung bei fehlgeschlagenen Tests

| Symptom | Wahrscheinliche Ursache | Behebung |
| --- | --- | --- |
| **Datenrennen in `webview_window_darwin.go`** | Fensterzustand wird außerhalb des Hauptthreads verändert | Den Aufruf über `application.InvokeAsync` / `Invoke` weiterleiten, damit er im Hauptthread ausgeführt wird |
| **Linux-Test hängt in einer Headless-CI** | GTK benötigt ein Display | Unter `xvfb-run` ausführen, z. B. `xvfb-run task test:examples:linux` |
| **Vorlagen-Build schlägt fehl** | Frontend-Lockdatei ist veraltet | `wails3 init` erneut für ein leeres Verzeichnis ausführen, um die Vorlage zu aktualisieren |
| **Coverpkg-Fehler** | Integrationstest importiert `main` | Zum Build-Tag `//go:build integration` wechseln und den Import damit absichern |

---

## 8. Neue Tests hinzufügen

1. **Unit-Test** — `*_test.go` erstellen und `go test ./...` ausführen
2. **Generator / CLI** — die Fälle unter `internal/generator/testcases/` oder `internal/commands/*_test.go` erweitern und `task test:generator` / `task test:cli` erneut ausführen
3. **Vorlagen / Beispiele** — sicherstellen, dass die ausgelieferten Vorlagen weiterhin mit `task test:templates` gebaut werden können

---

## 9. Übersicht der wichtigsten Dateien

| Zweck | Pfad |
| --- | --- |
| Roundtrip-Test des Generators | `internal/generator/generate_test.go` |
| Test der Build-Assets | `internal/commands/build-assets_test.go` |
| Leitfaden zu Race-Bedingungen und Cgo | `v3/TESTING.md` |
| Taskfile-Testziele | `v3/Taskfile.yaml` |
| Generator für Ereigniskonstanten | `v3/tasks/events/generate.go` |
| CI-Workflow | `.github/workflows/build-and-test-v3.yml` (Go 1.25 über `actions/setup-go@v5`) |
| Statische Analyse | `.github/workflows/semgrep.yml` |
| Veröffentlichung der Laufzeit auf npm | `.github/workflows/publish-npm.yml` |

---

Qualität ist in Wails v3 kein nachträglicher Gedanke. Mit Unit-Tests, Generator- und Vorlagen-Testsuiten, Race-Erkennung und einer plattformübergreifenden CI-Matrix können Sie zuverlässig beitragen, weil Sie wissen, dass Ihre Änderungen auf jedem unterstützten Betriebssystem erfolgreich ausgeführt werden. Viel Erfolg beim Testen!
