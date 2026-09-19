---
title: "Build- und Paketierungspipeline"
description: "Was im Hintergrund geschieht, wenn Sie `wails3 build` ausführen, wie plattformübergreifende Binärdateien erstellt und wie Installationspakete für jedes Betriebssystem erzeugt werden."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` ist absichtlich **schlank**: Es ist ein Taskfile-Wrapper, der zusätzliche Build-Tags an den Task `build` des Hostprojekts weiterleitet. Die eigentliche Arbeit erfolgt im projekteigenen `build/Taskfile.yml` (erzeugt von `wails3 init`), in `internal/commands/build-assets.go` (verwaltet zur Bake-Zeit eingebundene Assets) sowie in `internal/packager` (Linux-Paketierung mit nfpm) und `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, und `internal/commands/dot_desktop.go` (plattformabhängige Installationspakete).

Diese Seite behandelt:

1. Die tatsächlichen CLI-Einstiegspunkte
2. Den Taskfile-gesteuerten Build-Ablauf
3. Das Einbinden von Assets zur Bake-Zeit und das Einfügen von Build-Informationen
4. Die Paketierungs-Back-Ends der einzelnen Plattformen
5. Das Anpassen der Pipeline
6. Fehlerbehebung

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

`flags.Build` stellt ein **einziges** Flag bereit: `--tags` (wird als `EXTRA_TAGS=` weitergeleitet). Für `wails3 build` gibt es **keine** Flags `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` oder `-clean`. Cross-Kompilierung, Ausgabepfade, Symbole usw. werden in **`Taskfile.yml`**, **`build/config.yml`** und den unterstützenden Befehlen `wails3 generate icons`/`wails3 generate build-assets` konfiguriert.

`build/build.json` ist **nicht** Teil von v3 – die Konfiguration besteht aus `Taskfile.yml` und `build/config.yml`.

---

## 2. Taskfile-gesteuerter Build-Ablauf

Ein neu initialisiertes Projekt enthält ein `build/Taskfile.yml` mit ungefähr den folgenden Namespaces:

| Namespace | Tasks (Auswahl) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

`wails3 build` ruft standardmäßig den Namespace `build` des Hostbetriebssystems auf. Das Taskfile des Projekts führt seinerseits `go build` mit den hostspezifischen Flags aus. Um für ein anderes Betriebssystem zu bauen, führen Sie dessen Task direkt aus (z. B. `wails3 task darwin:build:universal`), statt `wails3 build` ein Flag zu übergeben.

Das standardmäßige Ausgabeverzeichnis ist **`bin/<APP_NAME>`** (ohne Präfix `build/bin/`).

---

## 3. Assets und Build-Informationen zur Bake-Zeit

| Aspekt | Datei |
| --- | --- |
| Erzeugung/Aktualisierung von Build-Assets | `internal/commands/build-assets.go` |
| Ausgabe von Build-Informationen (CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` – gibt Informationen aus; es ist **kein** `ldflags`-Injektor |
| Produktionsstub | `internal/assetserver/build_production.go` – `//go:build production` |
| Frontend-Bundles | über `//go:embed` in das anwendungseigene Paket eingebettet (z. B. neben `main.go`) |
| Windows-Ressourcen (`.syso`) | `internal/commands/syso.go` – erzeugt `rsrc_windows_<arch>.syso` |
| Windows-MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS-DMG-Eingaben | `internal/commands/dmg/` |
| Linux-`.desktop` | `internal/commands/dot_desktop.go` |

Die CLI erzeugt `bundled_assetserver.go` nicht automatisch für Ihre Anwendung – `internal/assetserver/bundled_assetserver.go` ist **von Hand geschrieben** und umschließt die unter `bundledassets/` eingebettete JavaScript-Runtime.

---

## 4. Paketierungs-Back-Ends

### Linux

Die Linux-Paketierung wird von **nfpm** gesteuert (nicht von `fpm`):

- `internal/packager/packager.go` umschließt `github.com/goreleaser/nfpm/v2` und stellt `CreatePackageFromConfig(pkgType, configPath, output)`/`CreatePackageFromConfigWriter(...)` bereit.
- Erzeugte Projekte enthalten unter `internal/commands/` die Konfigurationen `myapp.DEB`, `myapp.RPM` und `myapp.ARCHLINUX` im nfpm-Stil (verwendet von `wails3 tool package`).
- Die AppImage-Erzeugung befindet sich in `internal/commands/appimage.go`, das `linuxdeploy` + `linuxdeploy-plugin-gtk` aufruft (das Plug-in wird unter `internal/commands/linuxdeploy-plugin-gtk.sh` mitgeliefert).

Für `wails3 build` gibt es **kein** Flag `-package deb`/`rpm`. Verwenden Sie `wails3 tool package` oder das plattformspezifische Taskfile-Ziel.

### macOS

- `darwin:package` aus dem Taskfile des Projekts erzeugt das `.app`-Bundle.
- DMG-Assets befinden sich unter `internal/commands/dmg/`. Ein Projekt kann das Bundle nach Abschluss von `darwin:package` mit `hdiutil` in ein DMG verpacken (das Taskfile neuerer Vorlagen enthält dafür den Helfer `dmg`).
- CFBundle-Kennungen, Version und Copyright stammen während `wails3 init` aus den `-product*`-Flags sowie aus `build/config.yml`.

### Windows

- Die Windows-Paketierung ist auf **MSIX** ausgerichtet (nicht auf WiX/MSI). Den vollständigen Ablauf finden Sie unter `internal/commands/msix.go` und `internal/commands/webview2/`.
- Es gibt **kein** `internal/commands/packager.go` und **kein** `internal/commands/windows_resources/`-Verzeichnis.
- Die optionale Codesignierung der ausführbaren Datei erfolgt über `wails3 tool sign` (Authenticode) – siehe `internal/commands/sign.go`.

---

## 5. Pipeline anpassen

| Anforderung | Vorgehensweise |
| --- | --- |
| Zusätzliche Build-Tags | `wails3 build --tags myFeature,otherTag` |
| Linter-/Pre-Build-Schritt | Fügen Sie `build/Taskfile.yml` eine Task hinzu und machen Sie die betriebssystemspezifische `build`-Task davon abhängig |
| Cross-Kompilierung | Führen Sie die passende Betriebssystem-Task aus (z. B. `wails3 task linux:build`) – ein `-platform`-Flag gibt es nicht |
| Paketierung überspringen | Führen Sie nur die `build`-Task aus; `package` ist davon getrennt |
| Benutzerdefiniertes Paketierungswerkzeug | Legen Sie unter `internal/commands/myapp.*` eine Konfiguration ab und rufen Sie `wails3 tool package` mit `-config <file>` auf |
| Symbole entfernen | Bearbeiten Sie die `darwin:/windows:/linux:`-Task `build` so, dass `-ldflags "-s -w"` direkt an `go build` übergeben wird – `wails3 build` selbst hat kein `-ldflags`-Flag |

Alle Taskfile-Targets berücksichtigen die von Wails bereitgestellten Umgebungsvariablen (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), sodass benutzerdefinierte Tasks darauf zurückgreifen können.

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

Damit kennen Sie nun den gesamten Ablauf vom **Quellcode** bis zum **Installationsprogramm**. Kurz gesagt ist `wails3 build` selbst nur ein schlanker Wrapper – nahezu jede Anpassung erfolgt in `Taskfile.yml`/`build/config.yml` des Projekts oder über die expliziten Unterbefehle `wails3 generate …`/`wails3 tool …`. Viel Erfolg bei der Veröffentlichung!
