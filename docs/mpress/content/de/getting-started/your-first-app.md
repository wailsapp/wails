---
title: "Ihre erste Anwendung"
description: "Erstellen Sie Schritt für Schritt Ihre erste Wails-Desktopanwendung"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Diese Anleitung zeigt Ihnen, wie Sie Ihre erste Wails-v3-Anwendung erstellen – von der Projekteinrichtung über das Erstellen des Builds bis zum Entwicklungsablauf.

<br/>

<br/>

@steps
### Ein neues Projekt erstellen
Öffnen Sie Ihr Terminal und führen Sie den folgenden Befehl aus, um ein neues Wails-Projekt zu erstellen:

```bash
wails3 init -n myfirstapp
```

Dieser Befehl erstellt ein neues Verzeichnis namens `myfirstapp` mit allen erforderlichen Dateien.

   <video src="/assets/wails_init.mp4" controls></video>

### Die Projektstruktur erkunden
Wechseln Sie in das Verzeichnis `myfirstapp`. Dort finden Sie mehrere Dateien und Ordner:

@filetree
- build/           Enthält die vom Buildprozess verwendeten Dateien
  - appicon.png  Das Anwendungssymbol
  - config.yml   Buildkonfiguration
  - Taskfile.yml Build tasks
  - darwin/      macOS-spezifische Builddateien
    - Info.dev.plist Development configuration
    - Info.plist    Produktionskonfiguration
    - Taskfile.yml  macOS-Buildaufgaben
    - icons.icns    macOS-Anwendungssymbol
  - linux/       Linux-spezifische Builddateien
    - Taskfile.yml  Linux-Buildaufgaben
    - appimage/     AppImage-Paketierung
      - build.sh  AppImage-Buildskript
    - nfpm/        NFPM-Paketierung
      - nfpm.yaml Package configuration
      - scripts/  Buildskripte
  - windows/     Windows-spezifische Builddateien
    - Taskfile.yml        Windows-Buildaufgaben
    - icon.ico           Windows-Anwendungssymbol
    - info.json          Anwendungsmetadaten
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS-Installationsdateien
      - project.nsi                    NSIS-Projektdatei
      - wails_tools.nsh               NSIS-Hilfsskripte
- frontend/        Dateien der Frontend-Anwendung
  - index.html   HTML-Hauptdatei
  - main.js      JavaScript-Hauptdatei
  - package.json NPM package configuration
  - public/      Statische Ressourcen
  - Inter Font License.txt Font license
- .gitignore      Git-Ignore-Datei
- README.md       Projektdokumentation
- Taskfile.yml    Projektaufgaben
- go.mod          Go-Moduldatei
- go.sum          Prüfsummen der Go-Module
- greetservice.go Greeting service
- main.go         Hauptcode der Anwendung
@end

Nehmen Sie sich einen Moment Zeit, um diese Dateien zu erkunden und sich mit der Struktur vertraut zu machen.

@note{type="info"}
Wails v3 verwendet zwar standardmäßig [Task](https://taskfile.dev/) als Buildsystem, Sie können jedoch auch `make` oder ein beliebiges anderes Buildsystem verwenden.

@end

### Ihre Anwendung erstellen
Führen Sie zum Erstellen Ihrer Anwendung Folgendes aus:

```bash
wails3 build
```

Dieser Befehl kompiliert eine Debugversion Ihrer Anwendung und speichert sie in einem neuen Verzeichnis `bin`.

@note{type="info"}
`wails3 build` ist die Kurzform von `wails3 task build` und führt die Aufgabe `build` in `Taskfile.yml` aus.

@end

     <video src="/assets/wails_build.mp4" controls></video>

Nach dem Erstellen können Sie die Anwendung wie jede andere normale Anwendung ausführen:

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

Sie sehen eine einfache Benutzeroberfläche als Ausgangspunkt für Ihre Anwendung. Da es sich um die Debugversion handelt, werden außerdem Protokollmeldungen im Konsolenfenster angezeigt. Dies ist bei der Fehlersuche hilfreich.

### Entwicklungsmodus
Sie können die Anwendung auch im Entwicklungsmodus ausführen. In diesem Modus können Sie Ihren Frontend-Code ändern und die Änderungen direkt in der laufenden Anwendung sehen, ohne die gesamte Anwendung neu erstellen zu müssen.

1. Öffnen Sie ein neues Terminalfenster.
2. Führen Sie `wails3 dev` aus. Die Anwendung wird kompiliert und im Debugmodus ausgeführt.
3. Öffnen Sie `frontend/index.html` in einem Editor Ihrer Wahl.
4. Bearbeiten Sie den Code und ersetzen Sie `Please enter your name below` durch `Please enter your name below!!!`.
5. Speichern Sie die Datei.

Diese Änderung wird sofort in Ihrer Anwendung sichtbar.

Änderungen am Backend-Code lösen einen neuen Build aus:

1. Öffnen Sie `greetservice.go`.
2. Ersetzen Sie in der Zeile mit `return "Hello " + name + "!"` diesen Wert durch `return "Hello there " + name + "!"`.
3. Speichern Sie die Datei.

Die Anwendung wird innerhalb weniger Sekunden aktualisiert.

     <video src="/assets/wails_dev.mp4" controls></video>

### Ihre Anwendung paketieren
Sobald Ihre Anwendung zur Verteilung bereit ist, können Sie plattformspezifische Pakete erstellen:

@tabs{sync-key="platform"}
[Mac]
So erstellen Sie ein `.app`-Bundle:

```bash
wails3 package
```

Dadurch wird ein Produktions-Build erstellt und als `.app`-Bundle im Verzeichnis `bin` paketiert.

[Windows]
So erstellen Sie ein NSIS-Installationsprogramm:

```bash
wails3 package
```

Dadurch wird ein Produktions-Build erstellt und als NSIS-Installationsprogramm im Verzeichnis `bin` paketiert.

[Linux]
Wails unterstützt mehrere Paketformate zur Verteilung der Anwendung unter Linux:

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

Ausführlichere Informationen zu Paketierungsoptionen und deren Konfiguration finden Sie in unserem [Leitfaden zum Erstellen und Paketieren](/guides/build/building/).

### Versionsverwaltung und Modulnamen einrichten
Ihr Projekt wird mit dem Platzhalter-Modulnamen `changeme` erstellt. Es empfiehlt sich, diesen an die URL Ihres Repositorys anzupassen:

1. Erstellen Sie ein neues Repository auf GitHub (oder bei Ihrem bevorzugten Git-Hoster)
2. Initialisieren Sie git in Ihrem Projektverzeichnis:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Legen Sie Ihr Remote-Repository fest (ersetzen Sie den Wert durch die URL Ihres Repositorys):
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Passen Sie Ihren Modulnamen in `go.mod` an die URL Ihres Repositorys an:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Übertragen Sie Ihren Code:
  ```bash
  git push -u origin main
  ```


Dadurch entspricht der Name Ihres Go-Moduls den Namenskonventionen für Go-Module, und Ihr Code lässt sich leichter weitergeben.

@note{type="tip" title="Profi-Tipp"}
Sie können alle Initialisierungsschritte automatisieren, indem Sie beim Erstellen Ihres Projekts das Flag `-git` verwenden:

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

Dabei werden verschiedene Git-URL-Formate unterstützt:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` oder `ssh://git@github.com/username/project`
- Git-Protokoll: `git://github.com/username/project`
- Dateisystem: `file:///path/to/project.git`

@end

@end

## Herzlichen Glückwunsch!

Sie haben soeben Ihre erste Wails-Anwendung erstellt, entwickelt und paketiert. Dies ist erst der Anfang dessen, was Sie mit Wails v3 erreichen können.

## Nächste Schritte

Wenn Sie Wails noch nicht kennen, empfehlen wir Ihnen, als Nächstes unsere Tutorials zu lesen. Sie bieten eine praxisorientierte Einführung in die verschiedenen Funktionen von Wails. Das erste Tutorial heißt [Einen Dienst erstellen](/tutorials/01-creating-a-service/).

Wenn Sie bereits fortgeschrittene Kenntnisse haben, finden Sie im [Leitfaden zum Erstellen und Paketieren](/guides/build/building/) ausführlichere Informationen zur Verwendung von Wails.
