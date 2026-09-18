---
title: "Build-Anpassung"
description: "Passen Sie Ihren Build-Prozess mit Task und Taskfile.yml an"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Übersicht

Das Wails-Build-System ist ein flexibles und leistungsfähiges Werkzeug, das den Build-Prozess für Ihre Wails-Anwendungen vereinfacht. Es nutzt [Task](https://taskfile.dev), einen Task-Runner, mit dem Sie Tasks einfach definieren und ausführen können. Obwohl das v3-Build-System standardmäßig verwendet wird, unterstützt Wails den Ansatz, eigene Werkzeuge einzusetzen, sodass Entwickler ihren Build-Prozess nach Bedarf anpassen können.

Weitere Informationen zur Verwendung von Task finden Sie in der [offiziellen Dokumentation](https://taskfile.dev/usage/).

## Task: Das Herzstück des Build-Systems

[Task](https://taskfile.dev) ist eine moderne, in Go geschriebene Alternative zu Make. Task verwendet eine YAML-Datei, um Tasks und deren Abhängigkeiten zu definieren. Im Wails-Build-System spielt [Task](https://taskfile.dev) eine zentrale Rolle bei der Orchestrierung des Build-Prozesses.

Die zentrale `Taskfile.yml` befindet sich im Stammverzeichnis des Projekts, während plattformspezifische Tasks in `build/<platform>/Taskfile.yml`-Dateien definiert sind. Eine gemeinsame `Taskfile.yml`-Datei im Verzeichnis `build` enthält allgemeine Tasks, die von allen Plattformen gemeinsam genutzt werden.

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

Die Datei `Taskfile.yml` im Stammverzeichnis des Projekts ist der Haupteinstiegspunkt des Build-Systems. Sie definiert die Tasks und deren Abhängigkeiten. Dies ist die standardmäßige `Taskfile.yml`-Datei:

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## Plattformspezifische Taskfiles

Jede Plattform verfügt über ein eigenes Taskfile in den Plattformverzeichnissen unterhalb des Verzeichnisses `build`. Diese Dateien definieren die zentralen Tasks für die jeweilige Plattform. Jedes Taskfile bindet die allgemeinen Tasks aus der Datei `build/Taskfile.yml` ein.

### Windows

Speicherort: `build/windows/Taskfile.yml`

Das Windows-spezifische Taskfile enthält Tasks zum Erstellen, Paketieren und Ausführen der Anwendung unter Windows. Zu den wichtigsten Funktionen gehören:

- Erstellen mit optionalen Produktions-Flags
- Generieren der Symboldatei `.ico`
- Generieren der Windows-Datei `.syso`
- Erstellen eines NSIS-Installationsprogramms zur Paketierung

### Linux

Speicherort: `build/linux/Taskfile.yml`

Das Linux-spezifische Taskfile enthält Tasks zum Erstellen, Paketieren und Ausführen der Anwendung unter Linux. Zu den wichtigsten Funktionen gehören:

- Erstellen mit optionalen Produktions-Flags
- Erstellen eines AppImage sowie von deb-, rpm- und Arch-Linux-Paketen
- Generieren der Datei `.desktop` für Linux-Anwendungen

### macOS

Speicherort: `build/darwin/Taskfile.yml`

Das macOS-spezifische Taskfile enthält Tasks zum Erstellen, Paketieren und Ausführen der Anwendung unter macOS. Zu den wichtigsten Funktionen gehören:

- Erstellen von Binärdateien für die Architekturen amd64, arm64 und universal (beide)
- Generieren der Symboldatei `.icns`
- Erstellen eines `.app`-Bundles für die Verteilung
- Ad-hoc-Signieren von `.app`-Bundles
- Festlegen macOS-spezifischer Build-Flags und Umgebungsvariablen

## Task-Ausführung und Befehlsaliase

Der Befehl `wails3 task` ist eine eingebettete Version von [Taskfile](https://taskfile.dev), die die in Ihrer `Taskfile.yml` definierten Tasks ausführt.

Die Befehle `wails3 build` und `wails3 package` sind Aliase für `wails3 task build` beziehungsweise `wails3 task package`. Wenn Sie diese Befehle ausführen, übersetzt Wails sie intern in die entsprechende Task-Ausführung:

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Parameter an Tasks übergeben

Sie können CLI-Variablen im Format `KEY=VALUE` an Tasks übergeben. Diese Variablen werden über die Aliasbefehle weitergeleitet:

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

In Ihrer `Taskfile.yml` können Sie mit der Go-Template-Syntax auf diese Variablen zugreifen:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

### Variablen zur Anpassung von Go-Builds

Generierte Projekte stellen additive Variablen für eigene Go-Build-Tags, Linker-Flags und CGO bereit. Tag-Werte müssen kommagetrennte Namen ohne die Option `-tags` sein. Wails kombiniert sie mit den Tags, die der gewählte Build-Modus und die Plattform benötigen.

| Variable | Gilt für |
| --- | --- |
| `APP_TAGS` | Alle über Taskfiles ausgeführten Builds |
| `APP_TAGS_LINUX` | Linux-Builds |
| `APP_TAGS_DARWIN` | macOS-Builds |
| `APP_TAGS_WINDOWS` | Windows-Builds |
| `APP_TAGS_ANDROID` | Android-Builds |
| `APP_TAGS_IOS` | iOS-Builds |
| `APP_TAGS_SERVER` | Builds im Servermodus |
| `APP_LDFLAGS` | Zusätzliche Linker-Flags für alle über Taskfiles ausgeführten Builds |
| `APP_CGO_ENABLED` | Überschreibt CGO mit `0` oder `1` für Desktop- und Server-Builds |
| `EXTRA_TAGS` | Zusätzliche Tags für einen einzelnen Aufruf |

Zum Beispiel:

```bash
wails3 build APP_TAGS=sqlite_fts5,netgo APP_TAGS_LINUX=myapp_linux
wails3 build EXTRA_TAGS=diagnostics
APP_LDFLAGS='-X example.com/myapp/internal/version.Value=1.2.3' wails3 build
```

Legen Sie einen `APP_*`-Wert in der Stammdatei `Taskfile.yml` fest, wenn er dauerhaft als Projektstandard gelten soll. Ein mit dem Befehl übergebenes `KEY=value` hat die höchste Priorität und ersetzt die Variable für diesen Aufruf. Ein fest eingetragener Wert im Stamm-Taskfile hat Vorrang vor der gleichnamigen Umgebungsvariable des Prozesses. Bleibt der unten gezeigte Ausdruck mit Selbstreferenz und Standardwert erhalten, wird die Prozessumgebung verwendet, sofern kein Kommandozeilenwert angegeben ist. Verwenden Sie `EXTRA_TAGS`, um Tags für einen Build hinzuzufügen, ohne den dauerhaften `APP_TAGS`-Wert zu ersetzen.

Ist `APP_CGO_ENABLED` leer, verwenden Linux und macOS standardmäßig `1`, Windows verwendet `0`, native Server-Builds behalten den Go-Standardwert des Hosts und Server-Builds mit Docker verwenden `0`. Android- und iOS-Builds benötigen immer CGO und behalten `CGO_ENABLED=1`; `APP_CGO_ENABLED` überschreibt diese mobilen Toolchains nicht. Docker-basierte Desktop-Cross-Compilation und Server-Builds erhalten dieselben jeweils anwendbaren `APP_*`-Werte wie ihre nativen Taskfile-Entsprechungen.

@note{type="info" title="Bestehende Projekte"}
Die Stammdatei `Taskfile.yml` gehört zum Projekt und wird daher von `wails3 update build-assets` nicht überschrieben. Projekte, die vor Einführung dieser Variablen erstellt wurden, müssen die folgenden Einträge manuell zu ihrem `vars`-Block im Stamm-Taskfile hinzufügen:

```yaml
vars:
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
```

Ersetzen Sie einen Ausdruck mit Selbstreferenz und Standardwert durch einen festen Wert, um ihn dauerhaft als Projektstandard festzulegen.

@end

Diese Variablen werden von Taskfile-gesteuerten Builds verwendet. Die Build-Phase im generierten Xcode-Projekt ruft Go direkt auf und liegt außerhalb dieser Taskfile-Anpassung. Aus Xcode gestartete Builds verwenden diese Variablen daher derzeit nicht.

## Allgemeiner Build-Prozess

Der Build-Prozess umfasst auf allen Plattformen üblicherweise die folgenden Schritte:

1. Bereinigen der Go-Module
2. Erstellen des Frontends
3. Generieren von Symbolen
4. Kompilieren des Go-Codes mit plattformspezifischen Flags
5. Paketieren der Anwendung (plattformspezifisch)

## Build-Prozess anpassen

Das v3-Build-System bietet zwar eine solide Standardkonfiguration, lässt sich jedoch problemlos an die Anforderungen Ihres Projekts anpassen. Durch Änderungen an der `Taskfile.yml` und den plattformspezifischen Taskfiles können Sie:

- Neue Tasks hinzufügen
- Vorhandene Tasks ändern
- Die Reihenfolge der Task-Ausführung ändern
- Andere Werkzeuge und Skripte integrieren

Dank dieser Flexibilität können Sie den Build-Prozess auf Ihre spezifischen Anforderungen zuschneiden und zugleich von der Struktur des Wails-Build-Systems profitieren.

@note{type="tip" title="Taskfile kennenlernen"}
Wir empfehlen nachdrücklich, die [Taskfile](https://taskfile.dev)-Dokumentation zu lesen, um zu erfahren, wie Sie Taskfile effektiv einsetzen. Mit `wails3 task --version` können Sie ermitteln, welche Version von Taskfile in die Wails CLI eingebettet ist.

@end

## Entwicklungsmodus

Das Wails-Build-System umfasst einen leistungsfähigen Entwicklungsmodus, der durch automatisches Neuladen und Hot Module Replacement für eine bessere Entwicklungserfahrung sorgt. Dieser Modus wird mit dem Befehl `wails3 dev` aktiviert.

### Funktionsweise

Wenn Sie `wails3 dev` ausführen, läuft der folgende Prozess ab:

1. Der Befehl sucht nach einem verfügbaren Port und verwendet standardmäßig 9245, wenn keiner angegeben wurde.
2. Er richtet die Umgebungsvariablen für den Frontend-Entwicklungsserver (Vite) ein.
3. Er startet die Dateiüberwachung mit der Bibliothek [refresh](https://github.com/atterpac/refresh).

Die Bibliothek [refresh](https://github.com/atterpac/refresh) überwacht Dateiänderungen und löst neue Builds aus. Sie verwendet die Konfiguration unter dem Schlüssel `dev_mode` in der Datei `./build/config.yml`. Sie können festlegen, welche Verzeichnisse und Dateien die Bibliothek ignoriert, welche Dateien sie überwacht und welche Aktionen sie bei erkannten Änderungen ausführt. Die Standardkonfiguration funktioniert recht gut, Sie können sie aber jederzeit an Ihre Anforderungen anpassen.

### Konfiguration

Hier sehen Sie ein Beispiel für ihre Struktur:

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

Mit dieser Konfigurationsdatei können Sie:

- Den Stammpfad für die Dateiüberwachung festlegen
- Die Protokollierungsstufe konfigurieren
- Eine Entprellzeit für Dateiänderungsereignisse festlegen
- Bestimmte Verzeichnisse, Dateien oder Dateierweiterungen ignorieren
- Befehle definieren, die bei Dateiänderungen ausgeführt werden

### Entwicklungsmodus anpassen

Sie können den Entwicklungsmodus anpassen, indem Sie diese Werte in der Datei `config.yml` ändern.

Zu den Anpassungsmöglichkeiten gehören:

1. Die überwachten Verzeichnisse oder Dateien ändern
2. Die Entprellzeit anpassen, um zu steuern, wie schnell das System auf Änderungen reagiert
3. Die auszuführenden Befehle ergänzen oder ändern, um sie an die Anforderungen Ihres Projekts anzupassen

### Einen Browser für die Entwicklung verwenden

Wails v2 unterstützte zwar vollständig die Verwendung eines Browsers für die Entwicklung, dies führte jedoch häufig zu Verwirrung. Anwendungen, die im Browser funktionierten, funktionierten nicht zwangsläufig auch als Desktopanwendung, da in Webviews nicht alle Browser-APIs verfügbar sind.

Für Entwicklungsarbeiten mit Schwerpunkt auf der Benutzeroberfläche können Sie auch in v3 weiterhin einen Browser verwenden, indem Sie im Entwicklungsmodus die Vite-URL unter `http://localhost:9245` aufrufen. Dadurch stehen Ihnen beim Arbeiten an Styling und Layout leistungsfähige Browser-Entwicklungswerkzeuge zur Verfügung. Beachten Sie, dass Go-Bindings in diesem Modus *nicht funktionieren*.  
Wenn Sie Funktionen wie Bindings und Ereignisse testen möchten, wechseln Sie einfach zur Desktopansicht, um sicherzustellen, dass in der Produktionsumgebung alles einwandfrei funktioniert.
