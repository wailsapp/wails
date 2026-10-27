---
title: "Tests und kontinuierliche Integration"
description: "Wie Wails v3 mit Unit-Tests, Integrations-Testreihen, Race Detection und CI über GitHub Actions die Qualität sicherstellt."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Robuste Desktop-Frameworks erfordern absolut zuverlässige Tests. Wails v3 verfolgt eine **mehrschichtige Strategie**:

| Ebene | Ziel | Werkzeuge |
| --- | --- | --- |
| Unit-Tests | Schnelles Feedback für isolierte Funktionen | `go test ./...` |
| Generator-/CLI-Tests | `wails3 generate bindings` und die CLI-Infrastruktur validieren | `task test:generator`, `task test:cli` |
| Vorlagentests | Sicherstellen, dass sich jede ausgelieferte Vorlage weiterhin bauen lässt | `task test:templates` |
| Race Detection | Datenrennen in Runtime und Bridge erkennen | `go test -race ./...` |
| CI-Matrix | Plattformübergreifende Sicherheit bei jedem PR | GitHub Actions |

Dieses Dokument erläutert, **wo sich die Tests befinden**, **wie sie ausgeführt werden** und **was das Taskfile orchestriert**.

`v3/TESTING.md` dokumentiert plattformübergreifende Beispiel-Builds; es ist kein Katalog zur Unterdrückung von Datenrennen.

---

## 1. Verzeichniskonventionen

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

Richtlinien:

- **Unit-Tests gehören direkt neben den Code** (`foo.go` ↔ `foo_test.go`).
- Verwende für `pkg/`-Pakete den **Black-Box-Stil** (`package application_test`), wenn dies die API sauberer hält.
- Gemeinsam genutzte Fixtures befinden sich dort, wo sie verwendet werden (in diesem Arbeitsbaum gibt es kein zentrales `internal/testutil/`-Paket – binde Hilfsfunktionen stattdessen paketweise ein).

---

## 2. Unit-Tests

### Tests schreiben

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Empfehlungen:

- Verwende [`stretchr/testify`](https://github.com/stretchr/testify) – bereits in `go.mod` enthalten.
- Bevorzuge **tabellengesteuerte** Tests, wenn mehrere Eingaben, Grenzfälle oder erwartete Ergebnisse dasselbe Verhalten prüfen. Gib jedem Testfall einen aussagekräftigen Namen.
- Kapsele plattformspezifische Besonderheiten bei Bedarf mithilfe von Build-Tags (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Erwartete Testabdeckung

Für neue und geänderte Logik wird eine Go-Anweisungsabdeckung von 100 % erwartet. Miss das geänderte Paket, statt dich auf einen Prozentsatz für das gesamte Repository zu verlassen:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Einige Pfade lassen sich in einer normalen Testumgebung nachvollziehbarerweise nicht testen, beispielsweise ausschließlich plattformspezifische Fehler, hardwareabhängiges Verhalten oder ein defensiver Rückfallpfad, der sich nicht gefahrlos auslösen lässt. Beschränke solche Ausnahmen auf das Nötigste und erläutere jeden nicht abgedeckten Pfad in der PR-Beschreibung.

### Lokal ausführen

```bash
cd v3
go test ./... -cover
```

Alternativ über das Taskfile (tatsächliche Targets – es gibt keine `task test`-Kurzform):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. Integrationstests

`v3/tests/` enthält das paketübergreifende Integrationstest-Framework. Die Taskfile-Targets `test:example:*` und `test:examples:*` führen Kompilierungsprüfungen für darwin / windows / linux aus (unter Linux einschließlich Docker-basierter GTK3-/GTK4-Matrizen).

> Ausführbare Beispiele befinden sich in `v3/examples/`. Die Test-Targets wählen die
>
> für die Hostplattform oder CI-Matrix geeigneten Beispiele aus und bauen sie.

Führe die Smoke-Testreihe für die Hostplattform wie folgt aus:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Erkennung von Datenrennen

Datenrennen sind in GUI-Runtimes fatal.

### Ein Datenrennen untersuchen

Untersuche beide in Konflikt stehenden Zugriffe und reproduziere das Problem auf der betroffenen Plattform. Stufe einen Bericht nicht ohne Belege als unbedenklich ein und unterdrücke ihn nicht. `v3/TESTING.md` beschreibt die Kompilierung von Beispielen; erfolgreiche Builds belegen weder die Freiheit von Datenrennen noch das Verhalten der nativen Runtime.

### Lokale Race-Testreihe

```
go test -race ./...
```

> `wails3 dev` besitzt kein `-race`-Flag – seine CLI-Flags sind `--config`, `--port`
>
> und `-s` (aktiviert HTTPS). Um die Runtime unter dem Race Detector zu testen,
>
> baue eine Testanwendung mit `go build -race` und führe sie direkt aus.

---

## 5. GitHub-Actions-Workflows

Tatsächlich vorhandene Workflow-Dateien unter `.github/workflows/` (anhand des Arbeitsbaums verifiziert):

| Datei | Zweck |
| --- | --- |
| `build-and-test-v3.yml` | Go-1.26/1.27-Matrix, Runtime- und Generatorprüfungen, Kompilierung von Vorlagen und Beispielen sowie Prüfungen des Anwendungs-Builds. Linux verwendet standardmäßig GTK4/WebKitGTK 6.0; die Legacy-Variante setzt `BUILD_TAGS=gtk3`. Nativer Start und Interaktion müssen separat geprüft werden. |
| `cross-compile-test-v3.yml` | Plausibilitätsprüfungen der Cross-Kompilierung |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Changelog-Automatisierung |
| `nightly-release-v3.yml` | Nächtliche v3-Release-Artefakte |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Verwaltung der WebView2-Abhängigkeit und -Releases |
| `build-cross-image.yml` | Baut das Container-Image für den Cross-Compiler |
| `publish-npm.yml` | Veröffentlicht die eingebettete `@wailsio/runtime`-JS-Runtime auf npm |
| `pr-master.yml` | PR-seitige Prüfungen gegen den `master`-Branch |
| `semgrep.yml` | Statische Analyse mit Semgrep |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Repository-Pflege und v2-Workflows; die alten Crowdin-Workflows für Synchronisierung und Upload werden nicht mehr verwendet. |

Die CI-Schritte entsprechen den obigen Taskfile-Targets (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), sodass du die CI lokal exakt reproduzieren kannst. Der Smoke-Schritt `wails3 build` in `build-and-test-v3.yml` wird **ohne zusätzliche Flags** aufgerufen – `wails3 build` besitzt kein `-skip-package`-Flag.

---

## 6. Lokale Übereinstimmung mit der CI

Es gibt kein einzelnes übergeordnetes `task ci`-Ziel. Verketten Sie die tatsächlichen Ziele, um die CI lokal nachzubilden:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Fehler bei Tests beheben

| Symptom | Wahrscheinliche Ursache | Behebung |
| --- | --- | --- |
| **Race Condition in `webview_window_darwin.go`** | Der Fensterzustand wird außerhalb des Hauptthreads verändert | Leiten Sie den Aufruf über `application.InvokeAsync` / `Invoke` weiter, damit er im Hauptthread ausgeführt wird |
| **Linux-Test bleibt in einer Headless-CI hängen** | GTK benötigt eine Anzeige | Führen Sie den Test unter `xvfb-run` aus, z. B. mit `xvfb-run task test:examples:linux` |
| **Template-Build schlägt fehl** | Die Lockdatei des Frontends ist veraltet | Führen Sie `wails3 init` erneut für ein leeres Verzeichnis aus, um das Template zu aktualisieren |
| **Coverpkg-Fehler** | Integrationstest importiert `main` | Wechseln Sie zum Build-Tag `//go:build integration` und schützen Sie den Import damit |

---

## 8. Neue Tests hinzufügen

1. **Unit-Test** — erstellen Sie `*_test.go` und führen Sie `go test ./...` aus
2. **Generator / CLI** — erweitern Sie die Testfälle unter `internal/generator/testcases/` oder `internal/commands/*_test.go` und führen Sie `task test:generator` / `task test:cli` erneut aus
3. **Templates / Beispiele** — stellen Sie sicher, dass die ausgelieferten Templates weiterhin mit `task test:templates` gebaut werden können

---

## 9. Übersicht wichtiger Dateien

| Element | Pfad |
| --- | --- |
| Roundtrip-Test des Generators | `internal/generator/generate_test.go` |
| Test der Build-Assets | `internal/commands/build-assets_test.go` |
| Plattformübergreifende Beispiel-Builds | `v3/TESTING.md` |
| Testziele im Taskfile | `v3/Taskfile.yaml` |
| Generator für Ereigniskonstanten | `v3/tasks/events/generate.go` |
| CI-Workflow | `.github/workflows/build-and-test-v3.yml` (Go-1.26/1.27 über `actions/setup-go@v5`) |
| Statische Analyse | `.github/workflows/semgrep.yml` |
| Veröffentlichung des Runtime-npm-Pakets | `.github/workflows/publish-npm.yml` |

---

Qualität ist bei Wails v3 kein nachträglicher Gedanke. Mit Unit-Tests, Generator- und Template-Testsuiten, Race Detection und einer plattformübergreifenden CI-Matrix können Sie zuverlässig beitragen und darauf vertrauen, dass Ihre Änderungen auf jedem unterstützten Betriebssystem fehlerfrei ausgeführt werden. Viel Erfolg beim Testen!
