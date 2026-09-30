---
title: "Aufbau der Codebasis"
description: "Wie das Wails-v3-Repository organisiert ist und wie seine Bestandteile zusammenspielen"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 befindet sich in einem **Monorepo**, das die Framework-Laufzeit, die CLI, Beispiele, die Dokumentation und die Build-Toolchain enthält. Diese Seite erläutert die *Verzeichnisstruktur*, die für alle relevant ist, die sich mit den Interna befassen.

## Übersicht auf oberster Ebene

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

Im Folgenden betrachten wir den **`v3/`**-Verzeichnisbaum genauer.

## Stammverzeichnis `v3/`

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> Projektvorlagen werden unter `internal/templates/` bereitgestellt (je ein Ordner pro Framework-
>
> Stack sowie `base/`, `_common/` und `ios/`). Auf der obersten Ebene gibt es kein `v3/templates/`-
>
> Verzeichnis.

### Grundkonzept

1. **`pkg/`** stellt bereit, *was Anwendungsentwickler importieren*\
2. **`internal/`** enthält, *wie die Magie implementiert ist*\
3. **`cmd/wails3`** steuert *Projektlebenszyklus und Builds*\

Alles andere unterstützt diese drei Säulen.

---

## `cmd/` – Befehle

| Pfad | Hinweise |
| --- | --- |
| `v3/cmd/wails3` | Der **CLI-Einstiegspunkt**. Eine kleine `main.go` delegiert die gesamte Logik an Pakete in `internal/commands`. |
| `internal/commands/*` | Unterbefehle (init, dev, build, doctor, …). Jeder befindet sich in einer eigenen Datei und ist dadurch leicht auffindbar. |
| `internal/commands/task_wrapper.go` | Stellt die Verbindung zwischen CLI-Flags und der Taskfile-Build-Pipeline her. |

Die CLI ist zuständig für:

- **Projektgerüste** (`init`, Vorlagengenerierung)\
- **Orchestrierung des Entwicklungsservers** (`dev`, Live-Reload)\
- **Produktions-Builds und Paketierung** (`build`, `package`, plattformspezifische Wrapper)\
- **Diagnose** (`doctor`)\

---

## `internal/` – Der Maschinenraum

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### Wichtige Unterpakete

| Paket | Aufgabe | Anbindung |
| --- | --- | --- |
| `runtime` | Enthält den kleinen `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`-Build-Tag-Verbindungscode sowie die unter `runtime/desktop/` eingebettete JS-Laufzeit. Der eigentliche betriebssystemspezifische Code für Fenster, Zwischenablage, Dialoge und System-Tray-Symbole befindet sich in `pkg/application/*_{darwin,linux,windows}.go`. | Wird indirekt über `pkg/application` importiert. |
| `assetserver` | Dateiserver mit zwei Betriebsmodi:<br />• Entwicklung: stellt Dateien vom Datenträger bereit und leitet Anfragen an Vite weiter (`build_dev.go`)<br />• Produktion: bettet Assets über `go:embed` ein (`build_production.go`) | Wird beim Start von `pkg/application` initialisiert. |
| `generator` | Analysiert Go-Quellcode, um **Binding-Metadaten** zu erstellen, aus denen später TypeScript-/JS-Stubdateien und Ereigniskonstanten erzeugt werden. Einstiegspunkte: `generator.Generate` / `generator.Generator` auf Basis von `collect/` + `render/`. | Wird durch `wails3 generate bindings` ausgelöst. |
| `packager` | `nfpm`-Wrapper zum Erzeugen von Linux-Artefakten im Format `deb`/`rpm`/`archlinux` (gesteuert durch die nfpm-Konfigurationen `myapp.DEB`/`.RPM`/`.ARCHLINUX` unter `internal/commands/`). | Wird von `wails3 tool package` aufgerufen. macOS-DMG und Windows-MSIX befinden sich unter `internal/commands/{dmg,msix.go,webview2/}`. |

Unterstützende Hilfsprogramme (z. B. `s/`, `hash/`, `flags/`) halten interne Belange voneinander entkoppelt.

---

## `pkg/` – Öffentliche API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> Es gibt kein Paket namens `pkg/runtime/`, `pkg/options/` oder `pkg/menu/`. Optionen für Fenster und Menüs
>
> befinden sich neben `pkg/application` (z. B. `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), und `assetserver/` befindet sich unter `internal/`.

`pkg/application` initialisiert ein Wails-Programm:

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

Intern geschieht dabei Folgendes:

1. Bindet den `internal/runtime`-Build-Tag-Verbindungscode und den betriebssystemspezifischen Code in `pkg/application/` ein
2. Richtet eine `internal/assetserver`-Instanz ein
3. Registriert alle durch Bindings gesteuerten Nachrichtenprozessoren
4. Wechselt in den Hauptthread des Betriebssystems

---

## `internal/templates/` – Vorlagen für Projektgerüste

`internal/templates/` stellt **Basisvorlagen** (Go-Struktur unter `base/`, `_common/`, `ios/`) und **Frontend-Designs** (`vanilla[-ts]`, `react[-ts]`, `react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`, `svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`) bereit.

Bei `wails3 init -t react` führt die CLI Folgendes aus:

1. Kopiert die Go-Dateien aus `_common`
2. Führt das gewünschte Frontend-Paket zusammen
3. Führt `go mod tidy` aus (kann mit `--skipgomodtidy` übersprungen werden)

Das Bearbeiten von Vorlagen wirkt sich **nicht** auf bestehende Apps aus, sondern nur auf zukünftige `init`s. Die öffentlichen Beispiele befinden sich unter `v3/examples/`; sie ersetzen nicht die in der Dokumentation für Mitwirkende beschriebenen automatisierten Testsuiten.

---

## `tasks/` – Release-Automatisierung

Taskfiles kapseln die komplexe Cross-Kompilierung, Versionsanhebungen und die Generierung von Änderungsprotokollen. `internal/commands/task.go` verwendet sie programmgesteuert, sodass dieselbe Logik sowohl die **CLI** als auch die **CI** steuert.

---

## Zusammenspiel der Komponenten

```d2
direction: down
CLI: wails3-CLI
Generator: internal/generator
AssetDev: assetserver (Entwicklung)
Packager: internal/packager
AppRuntime: {
  label: App-Laufzeit
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: Betriebssystem-APIs
}
CLI -> Generator: erstellen/generieren
CLI -> AssetDev: Entwicklung
CLI -> Packager: paketieren
Generator -> ApplicationPkg: Bindings
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

*CLI → Generator → Laufzeitumgebung* bildet den zentralen Pfad vom **Quellcode** zur **ausgeführten Desktop-App**.

---

## Tipps zur Orientierung

| Zu verstehen … | Siehe … |
| --- | --- |
| Plattformspezifische Adapter | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go` (Fenster, Zwischenablage, Dialoge, Infobereich, Hauptthread, events_common). Linux-cgo: `pkg/application/linux_cgo*.go`. |
| Bridge-Protokoll | `pkg/application/messageprocessor*.go` |
| Asset-Workflow | `internal/assetserver/` (`build_dev.go` gegenüber `build_production.go`) |
| Paketierungsablauf | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| Vorlagen-Engine | `internal/templates/` (`templates.Install`, `templates.GetDefaultTemplates`) |
| Statische Analyse | `internal/generator/{generate.go,collect/,render/}` |

---

Sie verfügen nun über eine **mentale Übersicht** des Repositorys. Nutzen Sie sie zusammen mit `ripgrep`, der Funktion „Zu Datei/Symbol wechseln“ Ihrer IDE und den Beispiel-Apps, um sich eingehender mit beliebigen Funktionen zu befassen. Viel Spaß beim Hacken!
