---
title: "Ihre erste Anwendung"
description: "Erstellen Sie Schritt für Schritt Ihre erste Wails-Desktopanwendung"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Diese Anleitung zeigt Ihnen, wie Sie Ihre erste Wails-v3-Anwendung erstellen, einschließlich Projekteinrichtung, Build und Entwicklungsablauf.

<br/>

<br/>

@steps
### Neues Projekt erstellen
Öffnen Sie Ihr Terminal und führen Sie den folgenden Befehl aus, um ein neues Wails-Projekt zu erstellen:

```bash
wails3 init -n myfirstapp -t vanilla-js
```

In dieser Lektion wird ausdrücklich die JavaScript-Vanilla-Vorlage ausgewählt. Die einfache Vorlage `vanilla` verwendet TypeScript. Der Befehl erstellt `myfirstapp` mit allen erforderlichen Dateien.

   <video src="/assets/wails_init.mp4" controls></video>

### Projektstruktur erkunden
Wechseln Sie in das Verzeichnis `myfirstapp`. Dort finden Sie mehrere Dateien und Ordner:

@filetree
- build/           Enthält vom Build-Prozess verwendete Dateien
  - appicon.png  Das Anwendungssymbol
  - config.yml   Build-Konfiguration
  - Taskfile.yml Build tasks
  - darwin/      macOS-spezifische Build-Dateien
    - Info.dev.plist Development configuration
    - Info.plist    Produktionskonfiguration
    - Taskfile.yml  macOS-Build-Tasks
    - icons.icns    macOS-Anwendungssymbol
  - linux/       Linux-spezifische Build-Dateien
    - Taskfile.yml  Linux-Build-Tasks
    - appimage/     AppImage-Paketierung
      - build.sh  AppImage-Build-Skript
    - nfpm/        NFPM-Paketierung
      - nfpm.yaml Package configuration
      - scripts/  Build-Skripte
  - windows/     Windows-spezifische Build-Dateien
    - Taskfile.yml        Windows-Build-Tasks
    - icon.ico           Windows-Anwendungssymbol
    - info.json          Anwendungsmetadaten
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS-Installationsdateien
      - project.nsi                    NSIS-Projektdatei
      - wails_tools.nsh               NSIS-Hilfsskripte
- frontend/        Dateien der Frontend-Anwendung
  - index.html   HTML-Hauptdatei
  - src/         Frontend-Quellcode
    - main.js    JavaScript-Hauptdatei
  - package.json NPM package configuration
  - public/      Statische Ressourcen
  - Inter Font License.txt Font license
- .gitignore      Git-Ignore-Datei
- README.md       Projektdokumentation
- Taskfile.yml    Projekt-Tasks
- go.mod          Go-Moduldatei
- go.sum          Prüfsummen der Go-Module
- greetservice.go Greeting service
- main.go         Hauptcode der Anwendung
@end

Nehmen Sie sich einen Moment Zeit, um diese Dateien zu erkunden und sich mit der Struktur vertraut zu machen.

@note{type="info"}
Wails v3 verwendet zwar standardmäßig [Task](https://taskfile.dev/) als Build-System, Sie können jedoch ohne Weiteres `make` oder ein beliebiges anderes Build-System verwenden.

@end

### Anwendung erstellen
Führen Sie zum Erstellen Ihrer Anwendung Folgendes aus:

```bash
wails3 build
```

Mit den generierten Taskfiles erstellt dieser Befehl einen Produktions-Build in `bin`. Verwenden Sie `wails3 dev` für den Entwicklungsmodus oder `wails3 build DEV=true` für einen expliziten Entwicklungs-Build.

@note{type="info"}
`wails3 build` ist die Kurzform von `wails3 task build` und führt den Task `build` in `Taskfile.yml` aus.

@end

     <video src="/assets/wails_build.mp4" controls></video>

Nach dem Build können Sie die Anwendung wie jede andere normale Anwendung ausführen:

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

Die Starter-Benutzeroberfläche wird angezeigt. Bei einem Produktions-Build ist das Verhalten von Entwicklungsinspektor oder Konsole nicht garantiert. Verwenden Sie zum Debuggen den Entwicklungsmodus.

### Entwicklungsmodus
Sie können die Anwendung auch im Entwicklungsmodus ausführen. In diesem Modus können Sie Ihren Frontend-Code ändern und die Änderungen sofort in der laufenden Anwendung sehen, ohne die gesamte Anwendung neu erstellen zu müssen.

1. Öffnen Sie ein neues Terminalfenster.
2. Führen Sie `wails3 dev` aus. Die Anwendung wird kompiliert und im Debug-Modus ausgeführt.
3. Öffnen Sie `frontend/index.html` in einem Editor Ihrer Wahl.
4. Ändern Sie `placeholder="Your name"` des Namenseingabefelds in `placeholder="Your full name"` und speichern Sie die Datei.
5. Speichern Sie die Datei.

Diese Änderung wird sofort in Ihrer Anwendung angezeigt.

Jede Änderung am Backend-Code löst einen neuen Build aus:

1. Öffnen Sie `greetservice.go`.
2. Ändern Sie die Zeile mit `return "Hello " + name + "!"` in `return "Hello there " + name + "!"`.
3. Speichern Sie die Datei.

Die Anwendung wird innerhalb weniger Sekunden aktualisiert.

     <video src="/assets/wails_dev.mp4" controls></video>

### Anwendung paketieren
Sobald Ihre Anwendung zur Verteilung bereit ist, können Sie plattformspezifische Pakete erstellen:

@tabs{sync-key="platform"}
[Mac]
So erstellen Sie ein `.app`-Bundle:

```bash
wails3 package
```

Dadurch wird ein Produktions-Build erstellt und im Verzeichnis `bin` als `.app`-Bundle paketiert.

[Windows]
So erstellen Sie ein NSIS-Installationsprogramm:

```bash
wails3 package
```

Dadurch wird ein Produktions-Build erstellt und im Verzeichnis `bin` als NSIS-Installationsprogramm paketiert.

[Linux]
Wails unterstützt mehrere Paketformate für die Linux-Distribution:

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

Ausführlichere Informationen zu Paketierungsoptionen und Konfiguration finden Sie in unserem [Leitfaden zu Build und Paketierung](/guides/build/building/).

### Versionsverwaltung und Modulnamen einrichten
Ihr Projekt wird mit dem Platzhalter-Modulnamen `changeme` erstellt. Es empfiehlt sich, diesen an die URL Ihres Repositorys anzupassen:

1. Erstellen Sie ein neues Repository auf GitHub oder bei Ihrem bevorzugten Git-Hoster.
2. Initialisieren Sie Git in Ihrem Projektverzeichnis:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Legen Sie Ihr Remote-Repository fest und ersetzen Sie dabei die URL durch die Ihres Repositorys:
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Aktualisieren Sie den Modulnamen in `go.mod`, sodass er der URL Ihres Repositorys entspricht:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Pushen Sie Ihren Code:
  ```bash
  git push -u origin main
  ```


Dadurch entspricht der Name Ihres Go-Moduls den Namenskonventionen für Go-Module, und Ihr Code lässt sich leichter weitergeben.

@note{type="tip" title="Profi-Tipp"}
Sie können alle Initialisierungsschritte automatisieren, indem Sie beim Erstellen Ihres Projekts das Flag `-git` verwenden:

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

Dabei werden verschiedene Git-URL-Formate unterstützt:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` oder `ssh://git@github.com/username/project`
- Git-Protokoll: `git://github.com/username/project`
- Dateisystem: `file:///path/to/project.git`

@end

@end

## Herzlichen Glückwunsch!

Sie haben gerade Ihre erste Wails-Anwendung erstellt, entwickelt und paketiert. Dies ist erst der Anfang dessen, was Sie mit Wails v3 erreichen können.

## Nächste Schritte

Wenn Wails für Sie neu ist, empfehlen wir Ihnen, als Nächstes unsere Tutorials zu lesen. Sie führen Sie praxisnah durch die verschiedenen Funktionen von Wails. Das erste Tutorial ist [Einen Dienst erstellen](/tutorials/01-creating-a-service/).

Wenn Sie bereits fortgeschrittene Kenntnisse haben, finden Sie im [Leitfaden zum Erstellen und Paketieren](/guides/build/building/) ausführlichere Informationen zur Verwendung von Wails.
