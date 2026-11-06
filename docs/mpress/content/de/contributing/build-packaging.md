---
title: "Build- und Paketierungspipeline"
description: "Was im Hintergrund geschieht, wenn Sie `wails3 build` ausführen, wie plattformübergreifende Binärdateien erstellt werden und wie Installationspakete für jedes Betriebssystem erzeugt werden."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` ist bewusst **schlank** gehalten: Es ist ein Taskfile-Wrapper, der zusätzliche Build-Tags an den Task `build` des Hostprojekts weiterleitet. Die eigentliche Arbeit erfolgt im projekteigenen `build/Taskfile.yml` (erzeugt von `wails3 init`), in `internal/commands/build-assets.go` (zur Verwaltung der beim Build eingebetteten Assets) sowie in `internal/packager` (Linux-Paketierung mit nfpm) und `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go` und `internal/commands/dot_desktop.go` (plattformabhängige Installationsprogramme).

Diese Seite behandelt:

1. Die tatsächlichen CLI-Einstiegspunkte
2. Den Taskfile-gesteuerten Build-Ablauf
3. Das Einbetten von Assets und Einfügen von Build-Informationen
4. Die Paketierungs-Backends der einzelnen Plattformen
5. Das Anpassen der Pipeline
6. Die Fehlerbehebung

---

## 1. Tatsächliche CLI-Einstiegspunkte

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`:

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` stellt nur **ein einziges** Flag bereit: `--tags` (wird als `EXTRA_TAGS=` weitergeleitet). Für `wails3 build` gibt es **keine** Flags namens `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` oder `-clean`. Cross-Kompilierung, Ausgabepfade, Symbole usw. werden in **`Taskfile.yml`** und **`build/config.yml`** sowie über die unterstützenden Befehle `wails3 generate icons` und `wails3 generate build-assets` konfiguriert.

`build/build.json` ist **kein** Bestandteil von v3; die Konfiguration besteht aus `Taskfile.yml` und `build/config.yml`.

---

## 2. Taskfile-gesteuerter Build-Ablauf

Ein neu initialisiertes Projekt enthält ein `build/Taskfile.yml` mit ungefähr folgenden Namensräumen:

| Namensraum | Tasks (Auswahl) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

`wails3 build` ruft standardmäßig den Namensraum `build` des Hostbetriebssystems auf. Das Taskfile des Projekts ruft seinerseits `go build` mit den hostspezifischen Flags auf. Für einen Build für ein anderes Betriebssystem müssen Sie dessen Task direkt ausführen (z. B. `wails3 task darwin:build:universal`), statt ein Flag an `wails3 build` zu übergeben.

Das standardmäßige Ausgabeverzeichnis ist **`bin/<APP_NAME>`** (ohne Präfix `build/bin/`).

---

## 3. Beim Build eingebettete Assets und Build-Informationen

| Bereich | Datei |
| --- | --- |
| Erzeugen/Aktualisieren der Build-Assets | `internal/commands/build-assets.go` |
| Ausgabe der Build-Informationen (CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` – gibt Informationen aus; es ist **kein** Injektor für `ldflags` |
| Produktions-Stub | `internal/assetserver/build_production.go` – `//go:build production` |
| Frontend-Bundles | über `//go:embed` in das anwendungseigene Paket eingebettet (z. B. neben `main.go`) |
| Windows-Ressourcen (`.syso`) | `internal/commands/syso.go` – erzeugt `rsrc_windows_<arch>.syso` |
| Windows-MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS-DMG-Eingaben | `internal/commands/dmg/` |
| Linux-`.desktop` | `internal/commands/dot_desktop.go` |

Die CLI bettet `bundled_assetserver.go` nicht automatisch für Ihre Anwendung ein: `internal/assetserver/bundled_assetserver.go` wird **von Hand geschrieben** und umschließt die unter `bundledassets/` eingebettete JavaScript-Laufzeit.

---

## 4. Paketierungs-Backends

### Linux

Die Linux-Paketierung wird durch **nfpm** gesteuert (nicht durch `fpm`):

- `internal/packager/packager.go` umschließt `github.com/goreleaser/nfpm/v2` und stellt `CreatePackageFromConfig(pkgType, configPath, output)` sowie `CreatePackageFromConfigWriter(...)` bereit.
- Erzeugte Projekte enthalten die Konfigurationen `myapp.DEB`, `myapp.RPM` und `myapp.ARCHLINUX` im nfpm-Stil unter `internal/commands/`; sie werden von `wails3 tool package` verwendet.
- Die AppImage-Erzeugung befindet sich in `internal/commands/appimage.go`. Diese Datei ruft `linuxdeploy` und `linuxdeploy-plugin-gtk` auf (das Plug-in wird unter `internal/commands/linuxdeploy-plugin-gtk.sh` ausgeliefert).

Für `wails3 build` gibt es **kein** Flag `-package deb`/`rpm`. Verwenden Sie `wails3 tool package` oder das plattformspezifische Taskfile-Ziel.

### macOS

- `darwin:package` aus dem Taskfile des Projekts erzeugt das Bundle `.app`.
- Die DMG-Assets befinden sich unter `internal/commands/dmg/`. Ein Projekt kann das Bundle nach Abschluss von `darwin:package` mit `hdiutil` in ein DMG verpacken (das Taskfile neuerer Vorlagen enthält dafür einen `dmg`-Hilfstask).
- CFBundle-Kennungen, Version und Copyright stammen während `wails3 init` aus den Flags für `-product*` sowie aus `build/config.yml`.

### Windows

- Die Windows-Paketierung zielt auf **MSIX** ab (nicht auf WiX/MSI). Den vollständigen Ablauf finden Sie in `internal/commands/msix.go` und `internal/commands/webview2/`.
- Es gibt **kein** `internal/commands/packager.go` und **kein** Verzeichnis `internal/commands/windows_resources/`.
- Die optionale Codesignatur der ausführbaren Datei erfolgt über `wails3 tool sign` (Authenticode); siehe `internal/commands/sign.go`.

---

## 5. Anpassen der Pipeline

| Anforderung | Vorgehensweise |
| --- | --- |
| Zusätzliche Build-Tags | `wails3 build --tags myFeature,otherTag` |
| Linter-/Pre-Build-Schritt | Fügen Sie `build/Taskfile.yml` einen Task hinzu und definieren Sie ihn als Abhängigkeit des betriebssystemspezifischen Tasks `build` |
| Cross-Kompilierung | Führen Sie den Task des betreffenden Betriebssystems aus (z. B. `wails3 task linux:build`); ein Flag `-platform` gibt es nicht |
| Paketierung überspringen | Führen Sie nur den Task `build` aus; `package` ist davon getrennt |
| Benutzerdefinierter Paketierer | Legen Sie eine Konfiguration unter `internal/commands/myapp.*` ab und rufen Sie `wails3 tool package` mit `-config <file>` auf |
| Symbole entfernen | Bearbeiten Sie den `darwin:/windows:/linux:`-Task `build` so, dass `-ldflags "-s -w"` direkt an `go build` übergeben wird — `wails3 build` selbst hat kein `-ldflags`-Flag |

Alle Taskfile-Targets berücksichtigen die von Wails bereitgestellten Umgebungsvariablen (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), sodass benutzerdefinierte Tasks sie verwenden können.

---

## 6. Fehlerbehebung

| Symptom | Wahrscheinliche Ursache | Lösung |
| --- | --- | --- |
| **`ld: framework not found WebKit` (Mac)** | Xcode-CLI-Tools fehlen | `xcode-select --install` |
| **Leeres Fenster im Produktions-Build** | Frontend-Build fehlgeschlagen oder SPA-Routing | Prüfen Sie, ob `frontend/dist/index.html` vorhanden ist und Ihr Asset-Handler darauf zurückfällt |
| **MSIX-Paketierungswerkzeuge fehlen** | `WebView2` SDK/MSIX-Tools sind nicht installiert | Führen Sie `wails3 task install:msix:tools` aus |
| **`linuxdeploy` nicht gefunden** | Plugin fehlt im PATH | Installieren Sie `linuxdeploy` und führen Sie `internal/commands/linuxdeploy-plugin-gtk.sh` über den automatischen Installationsschritt der CLI aus |

`wails3 build` hat kein `-verbose`-Flag. Setzen Sie `TASK_X_VERBOSE=1` (Taskfile) oder prüfen Sie das Task-Target direkt, um die ausgeführten Befehle zu sehen.

---

## 7. Übersicht der wichtigsten Quelldateien

| Bereich | Datei |
| --- | --- |
| Build-Wrapper | `internal/commands/task_wrapper.go` (`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| Generierung der Build-Assets | `internal/commands/build-assets.go` (`GenerateBuildAssets`, `UpdateBuildAssets`) |
| Ausgabe der Build-Informationen | `internal/commands/tool_buildinfo.go` |
| AppImage-Builder | `internal/commands/appimage.go` |
| Linux-Paketierung (nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows-MSIX | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Windows-Ressourcengenerator | `internal/commands/syso.go` |
| macOS-DMG-Assets | `internal/commands/dmg/` |
| `.desktop`-Generator | `internal/commands/dot_desktop.go` |
| Versionskonstanten | `internal/version/version.go` |

Halten Sie diese Tabelle griffbereit, wenn Sie einen Build-Fehler untersuchen.

---

Sie haben nun den vollständigen Überblick vom **Quellcode** bis zum **Installationsprogramm**. Kurz gesagt ist `wails3 build` selbst nur ein schlanker Wrapper — fast alle Anpassungen erfolgen in `Taskfile.yml` / `build/config.yml` des Projekts oder über die expliziten Unterbefehle `wails3 generate …` / `wails3 tool …`. Viel Erfolg bei der Auslieferung!
