---
title: "Build-Anpassung"
description: "Passe deinen Build-Prozess mit Task und Taskfile.yml an"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Übersicht

Das Wails-Build-System ist ein flexibles und leistungsfähiges Werkzeug, das den Build-Prozess für deine Wails-Anwendungen vereinfacht. Es nutzt [Task](https://taskfile.dev), einen Task-Runner, mit dem du Aufgaben einfach definieren und ausführen kannst. Das v3-Build-System ist zwar die Standardeinstellung, Wails verfolgt jedoch den Ansatz „Bring your own tooling“, sodass Entwickler ihren Build-Prozess nach Bedarf anpassen können.

Weitere Informationen zur Verwendung von Task findest du in der [offiziellen Dokumentation](https://taskfile.dev/usage/).

## Task: Das Herzstück des Build-Systems

[Task](https://taskfile.dev) ist eine moderne, in Go geschriebene Alternative zu Make. Es verwendet eine YAML-Datei, um Aufgaben und ihre Abhängigkeiten zu definieren. Im Wails-Build-System spielt [Task](https://taskfile.dev) eine zentrale Rolle bei der Orchestrierung des Build-Prozesses.

Die zentrale `Taskfile.yml` befindet sich im Projektstammverzeichnis, während plattformspezifische Aufgaben in `build/<platform>/Taskfile.yml`-Dateien definiert sind. Eine gemeinsame `Taskfile.yml`-Datei im Verzeichnis `build` enthält plattformübergreifend verwendete Aufgaben.

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

Die Datei `Taskfile.yml` im Projektstammverzeichnis ist der Haupteinstiegspunkt des Build-Systems. Sie definiert die Aufgaben und ihre Abhängigkeiten. Dies ist die standardmäßige `Taskfile.yml`-Datei:

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

Jede Plattform verfügt über ein eigenes Taskfile, das sich in den Plattformverzeichnissen unterhalb des Verzeichnisses `build` befindet. Diese Dateien definieren die Kernaufgaben für die jeweilige Plattform. Jedes Taskfile bindet gemeinsame Aufgaben aus der Datei `build/Taskfile.yml` ein.

### Windows

Speicherort: `build/windows/Taskfile.yml`

Das Windows-spezifische Taskfile enthält Aufgaben zum Erstellen, Paketieren und Ausführen der Anwendung unter Windows. Zu den wichtigsten Funktionen gehören:

- Erstellen mit optionalen Produktions-Flags
- Erzeugen der `.ico`-Symboldatei
- Erzeugen der Windows-Datei `.syso`
- Erstellen eines NSIS-Installationsprogramms zur Paketierung

### Linux

Speicherort: `build/linux/Taskfile.yml`

Das Linux-spezifische Taskfile enthält Aufgaben zum Erstellen, Paketieren und Ausführen der Anwendung unter Linux. Zu den wichtigsten Funktionen gehören:

- Erstellen mit optionalen Produktions-Flags
- Erstellen eines AppImage sowie von deb-, rpm- und Arch-Linux-Paketen
- Erzeugen der Datei `.desktop` für Linux-Anwendungen

### macOS

Speicherort: `build/darwin/Taskfile.yml`

Das macOS-spezifische Taskfile enthält Aufgaben zum Erstellen, Paketieren und Ausführen der Anwendung unter macOS. Zu den wichtigsten Funktionen gehören:

- Erstellen von Binärdateien für die Architekturen amd64, arm64 und universal (beide)
- Erzeugen der `.icns`-Symboldatei
- Erstellen eines `.app`-Bundles für die Distribution
- Ad-hoc-Signieren von `.app`-Bundles
- Festlegen macOS-spezifischer Build-Flags und Umgebungsvariablen

## Task-Ausführung und Befehlsaliase

Der Befehl `wails3 task` ist eine eingebettete Version von [Taskfile](https://taskfile.dev), die die in deiner `Taskfile.yml` definierten Aufgaben ausführt.

Die Befehle `wails3 build` und `wails3 package` sind Aliase für `wails3 task build` beziehungsweise `wails3 task package`. Wenn du diese Befehle ausführst, übersetzt Wails sie intern in die entsprechende Task-Ausführung:

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Parameter an Tasks übergeben

Du kannst Tasks CLI-Variablen im Format `KEY=VALUE` übergeben. Diese Variablen werden über die Aliasbefehle weitergeleitet:

```bash
# These are equivalent:
wails3 build GOOS=linux
wails3 task build GOOS=linux
```

Die Standard-Tasks verwenden `GOOS`, `ARCH` und `DEV`. Variablen wie `VERSION`, `OUTPUT`, `PLATFORM` oder `CONFIG` wirken sich nur dann auf einen Task aus, wenn dein eigenes Taskfile sie ausliest. Das folgende Beispiel zeigt ein benutzerdefiniertes Taskfile und nicht den standardmäßigen Wails-Build-Task:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

Wenn Ihre benutzerdefinierte Paketierungs-Task `VERSION` und `OUTPUT` ausliest, können Sie beide auf dieselbe Weise übergeben:

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

## Allgemeiner Build-Prozess

Der Build-Prozess umfasst auf allen Plattformen üblicherweise die folgenden Schritte:

1. Bereinigen der Go-Module
2. Erstellen des Frontends
3. Erzeugen von Symbolen
4. Kompilieren des Go-Codes mit plattformspezifischen Flags
5. Paketieren der Anwendung (plattformspezifisch)

## Build-Prozess anpassen

Das v3-Build-System bietet zwar eine solide Standardkonfiguration, du kannst es jedoch problemlos an die Anforderungen deines Projekts anpassen. Durch Änderungen an `Taskfile.yml` und den plattformspezifischen Taskfiles kannst du:

- Neue Tasks hinzufügen
- Vorhandene Tasks ändern
- Die Ausführungsreihenfolge der Tasks ändern
- Andere Werkzeuge und Skripte integrieren

Dank dieser Flexibilität kannst du den Build-Prozess an deine konkreten Anforderungen anpassen und zugleich von der Struktur des Wails-Build-Systems profitieren.

@note{type="tip" title="Taskfile kennenlernen"}
Wir empfehlen ausdrücklich, die Dokumentation zu [Taskfile](https://taskfile.dev) zu lesen, um zu verstehen, wie du Taskfile effektiv verwendest. Mit `wails3 task --version` kannst du ermitteln, welche Taskfile-Version in die Wails-CLI eingebettet ist.

@end

## Entwicklungsmodus

Das Wails-Build-System bietet einen leistungsfähigen Entwicklungsmodus, der durch Live-Neuladen und Hot Module Replacement die Entwicklung verbessert. Du aktivierst diesen Modus mit dem Befehl `wails3 dev`.

### Funktionsweise

Wenn du `wails3 dev` ausführst, läuft der folgende Prozess ab:

1. Der Befehl sucht nach einem verfügbaren Port. Wenn keiner angegeben ist, verwendet er standardmäßig 9245.
2. Er richtet die Umgebungsvariablen für den Frontend-Entwicklungsserver (Vite) ein.
3. Er startet die Dateiüberwachung mit der Bibliothek [refresh](https://github.com/atterpac/refresh).

Die Bibliothek [refresh](https://github.com/atterpac/refresh) überwacht Dateiänderungen und löst erneute Builds aus. Sie verwendet die Konfiguration, die in der Datei `./build/config.yml` unter dem Schlüssel `dev_mode` definiert ist. Du kannst sie so konfigurieren, dass sie bestimmte Verzeichnisse und Dateien ignoriert, festlegt, welche Dateien überwacht werden, und bestimmt, welche Aktionen bei erkannten Änderungen ausgeführt werden. Die Standardkonfiguration funktioniert recht gut, du kannst sie jedoch an deine Anforderungen anpassen.

### Konfiguration

Hier ist ein Beispiel für ihre Struktur:

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

Mit dieser Konfigurationsdatei kannst du:

- Den Stammpfad für die Dateiüberwachung festlegen
- Die Protokollierungsstufe konfigurieren
- Eine Entprellzeit für Dateiänderungsereignisse festlegen
- Bestimmte Verzeichnisse, Dateien oder Dateierweiterungen ignorieren
- Befehle definieren, die bei Dateiänderungen ausgeführt werden

### Entwicklungsmodus anpassen

Sie können den Entwicklungsmodus anpassen, indem Sie diese Werte in der Datei `config.yml` ändern.

Mögliche Anpassungen sind beispielsweise:

1. Überwachte Verzeichnisse oder Dateien ändern
2. Die Entprellzeit anpassen, um zu steuern, wie schnell das System auf Änderungen reagiert
3. Ausführungsbefehle hinzufügen oder ändern, um sie an die Anforderungen Ihres Projekts anzupassen

### Browser für die Entwicklung verwenden

Wails v2 unterstützte zwar uneingeschränkt die Verwendung eines Browsers für die Entwicklung, dies sorgte jedoch häufig für Verwirrung. Anwendungen, die im Browser funktionierten, mussten nicht zwangsläufig auch in der Desktop-Anwendung funktionieren, da nicht alle Browser-APIs in Webviews verfügbar sind.

Bei Entwicklungsarbeiten mit Schwerpunkt auf der Benutzeroberfläche können Sie weiterhin flexibel einen Browser in v3 verwenden, indem Sie im Entwicklungsmodus die Vite-URL unter `http://localhost:9245` aufrufen. So können Sie beim Arbeiten an Styling und Layout die leistungsstarken Browser-Entwicklungswerkzeuge nutzen. Beachten Sie, dass Go-Bindings in diesem Modus *nicht funktionieren*. Wenn Sie Funktionen wie Bindings und Ereignisse testen möchten, wechseln Sie einfach zur Desktop-Ansicht, um sicherzustellen, dass in der Produktionsumgebung alles einwandfrei funktioniert.
