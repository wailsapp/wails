---
title: "CLI-Referenz"
description: "Vollständige Referenz der Wails-CLI-Befehle"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

Die Wails-CLI bietet umfassende Befehle zum Entwickeln, Erstellen und Warten Ihrer Wails-Anwendungen.

## Kernbefehle

Kernbefehle sind die wichtigsten Befehle zum Erstellen, Entwickeln und Bauen von Projekten.

Alle CLI-Befehle haben das folgende Format: `wails3 <command>`.

### `init`

Initialisiert ein neues Wails-Projekt. Während der Initialisierung wird der Befehl `go mod tidy` ausgeführt, um die Projektpakete zu aktualisieren. Dies lässt sich umgehen, indem Sie beim Befehl `init` das Flag `-skipgomodtidy` verwenden.

```bash
wails3 init [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-p` | Name des Go-Pakets | `main` |
| `-t` | Vorlagenname oder URL | `vanilla` |
| `-n` | Projektname |  |
| `-d` | Projektverzeichnis | `.` |
| `-q` | Ausgabe unterdrücken | `false` |
| `-l` | Vorlagen auflisten | `false` |
| `-mod` | Go-Modulpfad (wird bei fehlender Angabe aus `-git` berechnet) |  |
| `-git` | URL des Git-Repositorys |  |
| `-s` | Warnung bei Verwendung einer Remote-Vorlage überspringen | `false` |
| `-productname` | Produktname | `My Product` |
| `-productdescription` | Produktbeschreibung | `My Product Description` |
| `-productversion` | Produktversion | `0.1.0` |
| `-productcompany` | Unternehmensname | `My Company` |
| `-productcopyright` | Urheberrechtshinweis | `© now, My Company` |
| `-productcomments` | Dateikommentare | `This is a comment` |
| `-productidentifier` | Produktkennung |  |
| `-skipgomodtidy` | go mod tidy überspringen | `false` |

Das Flag `-git` akzeptiert verschiedene Git-URL-Formate:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` oder `ssh://git@github.com/username/project`
- Git-Protokoll: `git://github.com/username/project`
- Dateisystem: `file:///path/to/project.git`

Wenn dieses Flag angegeben wird, geschieht Folgendes:

1. Ein Git-Repository im Projektverzeichnis initialisieren
2. Die angegebene URL als Remote-Repository „origin“ festlegen
3. Den Modulnamen in `go.mod` an die Repository-URL anpassen
4. Alle Dateien hinzufügen

### `dev`

Führt die Anwendung im Entwicklungsmodus aus. Dadurch erhalten Sie eine Live-Ansicht Ihres Frontend-Codes. Änderungen werden direkt in der laufenden Anwendung angezeigt, ohne dass Sie die gesamte Anwendung neu bauen müssen. Änderungen an Ihrem Go-Code werden ebenfalls erkannt; die Anwendung wird automatisch neu gebaut und gestartet.

```bash
wails3 dev [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-config` | Pfad zur Konfigurationsdatei | `./build/config.yml` |
| `-port` | Port des Vite-Entwicklungsservers | `9245` |
| `-s` | HTTPS aktivieren | `false` |

@note{type="info"}
Dies entspricht der Ausführung von `wails3 task dev` und führt die Task `dev` im Haupt-Taskfile des Projekts aus. Sie können dies durch Bearbeiten der Datei `Taskfile.yml` anpassen.

@end

### `build`

Erstellt eine Debug-Version Ihrer Anwendung. Standardmäßig wird sie für die aktuelle Plattform und Architektur erstellt.

```bash
wails3 build [flags] [CLI variables...]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-tags` | Zusätzliche Go-Build-Tags (durch Kommas getrennt) |  |

Sie können den Build mit CLI-Variablen anpassen:

```bash
wails3 build PLATFORM=linux CONFIG=production
```

Verwenden Sie das Flag `-tags`, um benutzerdefinierte Go-Build-Tags zu übergeben:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

Die Tags werden als `EXTRA_TAGS` an das zugrunde liegende Taskfile weitergeleitet.

@note{type="info"}
Dies entspricht der Ausführung von `wails3 task build`, wodurch der Task `build` im Haupt-Taskfile des Projekts ausgeführt wird. Alle an `build` übergebenen CLI-Variablen werden an den zugrunde liegenden Task weitergeleitet. Sie können den Build-Prozess anpassen, indem Sie die Datei `Taskfile.yml` bearbeiten.

@end

### `package`

Erstellt plattformspezifische Pakete zur Distribution.

```bash
wails3 package [CLI variables...]
```

Sie können die Paketerstellung mit CLI-Variablen anpassen:

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

#### Pakettypen

Für jede Plattform sind die folgenden Pakettypen verfügbar:

| Plattform | Pakettyp |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
Dies entspricht `wails3 task package`, wodurch der Task `package` im Haupt-Taskfile des Projekts ausgeführt wird. Alle an `package` übergebenen CLI-Variablen werden an den zugrunde liegenden Task weitergeleitet. Sie können die Paketerstellung anpassen, indem Sie die Datei `Taskfile.yml` bearbeiten.

@end

### `task`

Führt die in der Taskfile.yml Ihres Projekts definierten Tasks aus. Dies ist eine eingebettete Version von [Taskfile](https://taskfile.dev), mit der Sie benutzerdefinierte Build-, Test- und Deployment-Tasks definieren und ausführen können.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI-Variablen

Sie können Variablen im Format `KEY=VALUE` an Tasks übergeben:

```bash
wails3 task build PLATFORM=linux CONFIG=production
wails3 task deploy ENV=staging VERSION=1.2.3
```

Auf diese Variablen können Sie in Ihrer Taskfile.yml mithilfe der Go-Template-Syntax zugreifen:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-h` | Zeigt die Verwendung von Task an | `false` |
| `-i` | Erstellt eine neue Taskfile.yml | `false` |
| `-list` | Listet Tasks mit Beschreibungen auf | `false` |
| `-list-all` | Listet alle Tasks auf (mit oder ohne Beschreibungen) | `false` |
| `-json` | Formatiert die Task-Liste als JSON | `false` |
| `-status` | Beendet das Programm mit einem von null verschiedenen Exit-Code, wenn der Task nicht aktuell ist | `false` |
| `-f` | Erzwingt die Ausführung, auch wenn der Task aktuell ist | `false` |
| `-w` | Aktiviert den Überwachungsmodus für den angegebenen Task | `false` |
| `-v` | Aktiviert den ausführlichen Modus | `false` |
| `-version` | Gibt die Task-Version aus | `false` |
| `-s` | Deaktiviert die Anzeige der ausgeführten Befehle | `false` |
| `-p` | Führt Tasks parallel aus | `false` |
| `-dry` | Kompiliert Tasks und gibt sie aus, ohne sie auszuführen | `false` |
| `-summary` | Zeigt eine Zusammenfassung zu einem Task an | `false` |
| `-x` | Übernimmt den Exit-Code des Tasks | `false` |
| `-dir` | Legt das Ausführungsverzeichnis fest |  |
| `-taskfile` | Auszuführende Taskfile auswählen |  |
| `-output` | Legt den Ausgabestil fest: [interleaved|group|prefixed] |  |
| `-c` | Farbige Ausgabe (standardmäßig aktiviert) | `true` |
| `-C` | Anzahl der gleichzeitig ausgeführten Tasks begrenzen |  |
| `-interval` | Intervall für die Prüfung auf Änderungen (in Sekunden) |  |

#### Beispiele

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build PLATFORM=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

Startet den MCP-Server des Wails-Projekts für die agentengestützte Projektverwaltung. Dieser ist vom MCP-Server getrennt, der in eine laufende Anwendung kompiliert ist: `wails3 mcp` verwaltet Projektdateien und Lebenszyklusbefehle, während der MCP-Server der Anwendung die laufende WebView steuert.

```bash
wails3 mcp [flags]
```

Der Transport wird automatisch ausgewählt:

- Wenn ein MCP-Host Wails mit weitergeleiteter Standardeingabe und Standardausgabe startet, verwendet der Server **stdio**.
- Bei interaktiver Ausführung in einem Terminal verwendet der Server **Streamable HTTP** auf `127.0.0.1` und fordert beim Betriebssystem einen freien Port an.

Verwenden Sie `--stdio` oder `--http`, um einen Transport explizit auszuwählen. Verwenden Sie `--port 0`, um im HTTP-Modus einen freien Loopback-Port auszuwählen.

#### MCP-Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `--root` | Zulässiges Projektstammverzeichnis. Pfade und symbolische Links außerhalb dieses Verzeichnisses werden abgelehnt. | Aktuelles Verzeichnis |
| `--token` | Sitzungs-/Bearer-Token für Tools, die Änderungen vornehmen oder Prozesse steuern. Verwendet ersatzweise `WAILS_MCP_TOKEN`. | Sicher generiert |
| `--stdio` | Erzwingt den stdio-Transport. | Automatisch |
| `--http` | Erzwingt den Transport über Streamable HTTP. | Automatisch |
| `--port` | HTTP-Port; `0` wählt einen freien Loopback-Port aus. | `0` |

Im HTTP-Modus gibt Wails den Endpunkt und das Bearer-Token über stderr aus. Im stdio-Modus ist das Token in den MCP-Initialisierungsanweisungen enthalten. Der Server ermöglicht keine beliebige Ausführung von Shell-Befehlen. Remote-Vorlagen und Git-Remotes erfordern eine ausdrückliche Genehmigung über die Eingabe `allowExternal` des Tools.

### `doctor`

Führt eine Systemprüfung durch und zeigt einen Statusbericht an.

```bash
wails3 doctor
```

## Generierungsbefehle

Mit Generierungsbefehlen können verschiedene Projektressourcen wie Bindings, Symbole und Build-Dateien erstellt werden. Alle Generierungsbefehle verwenden den Basisbefehl `wails3 generate <command>`.

### `generate bindings`

Generiert Bindings und Modelle für Ihren Go-Code.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-f` | Zusätzliche Go-Build-Flags |  |
| `-d` | Ausgabeverzeichnis | `frontend/bindings` |
| `-models` | Dateiname für Modelle | `models` |
| `-index` | Dateiname für den Index | `index` |
| `-ts` | TypeScript generieren | `false` |
| `-i` | TS-Schnittstellen verwenden | `false` |
| `-b` | Gebündelte Runtime verwenden | `false` |
| `-names` | Namen anstelle von IDs verwenden | `false` |
| `-noindex` | Indexdateien überspringen | `false` |
| `-noevents` | Generierung ereignisbezogener Bindings überspringen | `false` |
| `-dry` | Testlauf | `false` |
| `-silent` | Stiller Modus | `false` |
| `-v` | Debug-Ausgabe | `false` |
| `-clean` | Ausgabeverzeichnis vor der Generierung bereinigen | `true` |

### `generate build-assets`

Generiert Build-Ressourcen für Ihre Anwendung.

```bash
wails3 generate build-assets [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-name` | Projektname |  |
| `-dir` | Ausgabeverzeichnis | `build` |
| `-silent` | Ausgabe unterdrücken | `false` |
| `-company` | Unternehmensname |  |
| `-productname` | Produktname |  |
| `-description` | Produktbeschreibung |  |
| `-version` | Produktversion |  |
| `-identifier` | Produktkennung | `com.wails.[name]` |
| `-copyright` | Urheberrechtshinweis |  |
| `-comments` | Dateikommentare |  |

### `generate icons`

Generiert Anwendungssymbole.

```bash
wails3 generate icons [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-input` | PNG-Eingabedatei | Erforderlich |
| `-windowsfilename` | Name der Windows-Ausgabedatei |  |
| `-macfilename` | Name der macOS-Ausgabedatei |  |
| `-sizes` | Symbolgrößen (durch Kommas getrennt) | `256,128,64,48,32,16` |
| `-example` | Beispielsymbol generieren | `false` |
| `-iconcomposerinput` | Icon-Composer-Eingabedatei (`.icon`) |  |
| `-macassetdir` | Ausgabeverzeichnis für Mac-Ressourcen (Assets.car + icns) |  |

#### Icon Composer (macOS)

Unter macOS 26+ können Sie mit Icon Composer erstellte `.icon`-Dateien verwenden, um `Assets.car` und `icons.icns` zu generieren:

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Dadurch wird die `.icon`-Datei mit Apples Befehl `actool` kompiliert. Erfordert Xcode mit `actool` in Version 26 oder neuer.

Wenn Sie Icon Composer verwenden, legen Sie `cfBundleIconName` in Ihrer `build/config.yml` so fest, dass der Wert dem Namen der `.icon`-Datei (ohne Erweiterung) entspricht:

```yaml
info:
  cfBundleIconName: "appicon"
```

Wenn die Option nicht festgelegt ist und `Assets.car` vorhanden ist, wird standardmäßig `"appicon"` verwendet.

### `generate syso`

Erzeugt eine Windows-.syso-Datei.

```bash
wails3 generate syso [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-manifest` | Pfad zur Manifestdatei | Erforderlich |
| `-icon` | Pfad zur Symboldatei | Erforderlich |
| `-info` | Pfad zur Datei mit Versionsinformationen |  |
| `-arch` | Zielarchitektur | Aktuelle GOARCH |
| `-out` | Name der Ausgabedatei | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Erzeugt eine Linux-.desktop-Datei.

```bash
wails3 generate .desktop [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-name` | Anwendungsname | Erforderlich |
| `-exec` | Pfad zur ausführbaren Datei | Erforderlich |
| `-icon` | Pfad zur Symboldatei |  |
| `-categories` | Anwendungskategorien | `Utility` |
| `-comment` | Anwendungskommentar |  |
| `-terminal` | Im Terminal ausführen | `false` |
| `-keywords` | Suchbegriffe |  |
| `-version` | Anwendungsversion |  |
| `-genericname` | Allgemeiner Name |  |
| `-startupnotify` | Startbenachrichtigung anzeigen | `false` |
| `-mimetype` | Unterstützte MIME-Typen |  |
| `-output` | Name der Ausgabedatei | `[name].desktop` |

### `generate runtime`

Erzeugt die vorgefertigte Version der Runtime.

```bash
wails3 generate runtime
```

### `generate constants`

Erzeugt JavaScript-Konstanten aus Go-Code.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

Erzeugt ein Windows-WebView2-Bootstrap-Installationsprogramm zur Verteilung.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

Erstellt das Grundgerüst eines neuen Projektvorlagenverzeichnisses.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Erzeugt ein Linux-AppImage.

```bash
wails3 generate appimage [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-binary` | Pfad zur Binärdatei | Erforderlich |
| `-icon` | Pfad zur Symboldatei | Erforderlich |
| `-desktop` | Pfad zur .desktop-Datei | Erforderlich |
| `-builddir` | Build-Verzeichnis | Temporäres Verzeichnis |
| `-output` | Ausgabeverzeichnis | `.` |

## Dienstbefehle

Dienstbefehle dienen zur Verwaltung von Wails-Diensten. Alle Dienstbefehle verwenden den Basisbefehl `wails3 service <command>`.

### `service init`

Initialisiert einen neuen Dienst.

```bash
wails3 service init [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-n` | Dienstname | `example_service` |
| `-d` | Dienstbeschreibung | `Example service` |
| `-p` | Paketname |  |
| `-o` | Ausgabeverzeichnis | `.` |
| `-q` | Ausgabe unterdrücken | `false` |
| `-a` | Name des Autors |  |
| `-v` | Version |  |
| `-w` | Website-URL |  |
| `-r` | Repository-URL |  |
| `-l` | Lizenz |  |

## Werkzeugbefehle

Werkzeugbefehle stellen Hilfsfunktionen für Entwicklung und Debugging bereit. Alle Werkzeugbefehle verwenden den Basisbefehl `wails3 tool <command>`.

### `tool checkport`

Prüft, ob ein Port geöffnet ist. Dies ist nützlich, um zu testen, ob vite ausgeführt wird.

```bash
wails3 tool checkport [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-port` | Zu prüfender Port | `9245` |
| `-host` | Zu prüfender Host | `localhost` |

### `tool watcher`

Überwacht Dateien und führt bei Änderungen einen Befehl aus.

```bash
wails3 tool watcher [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-config` | Pfad zur Konfigurationsdatei | `./build/config.yml` |
| `-ignore` | Zu ignorierende Muster |  |
| `-include` | Einzuschließende Muster |  |

### `tool cp`

Kopiert Dateien.

```bash
wails3 tool cp
```

### `tool buildinfo`

Zeigt Build-Informationen zur Anwendung an.

```bash
wails3 tool buildinfo
```

### `tool version`

Erhöht eine semantische Version entsprechend den angegebenen Optionen.

```bash
wails3 tool version [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-v` | Zu erhöhende aktuelle Version |  |
| `-major` | Hauptversion erhöhen | `false` |
| `-minor` | Nebenversion erhöhen | `false` |
| `-patch` | Patch-Version erhöhen | `false` |
| `-prerelease` | Vorabversion erhöhen (z. B. von alpha.5 auf alpha.6) | `false` |

Der Befehl verwendet folgende Rangfolge: Hauptversion > Nebenversion > Patch-Version > Vorabversion. Er behält das Präfix „v“ bei, sofern es in der Eingabeversion vorhanden ist, ebenso wie alle Bestandteile der Vorabversion und Metadaten.

Anwendungsbeispiel:

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.0.0-alpha.5 -prerelease # Output: v3.0.0-alpha.6
```

### `tool package`

Erzeugt Linux-Pakete (deb, rpm, archlinux).

```bash
wails3 tool package [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-format` | Paketformat (deb, rpm, archlinux) | `deb` |
| `-name` | Name der ausführbaren Datei | `myapp` |
| `-config` | Pfad zur Konfigurationsdatei |  |
| `-out` | Ausgabeverzeichnis | `.` |

### `tool lipo`

Erstellt ein universelles macOS-Binärprogramm, indem architekturspezifische Binärprogramme kombiniert werden.

```bash
wails3 tool lipo [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-output` | Pfad zum Ausgabebinärprogramm |  |

### `tool capabilities`

Prüft die Build-Fähigkeiten des Systems (Verfügbarkeit von GTK4/GTK3 unter Linux).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

Erzeugt Docker-Optionen zum Einbinden von Volumes für die Cross-Kompilierung. Gibt `-v`-Optionen für den Go-Modul-Cache und alle lokalen `replace`-Direktiven in `go.mod` aus, die in Taskfile-`docker run`-Befehlen verwendet werden können.

```bash
wails3 tool docker-mounts
```

### `tool has`

Prüft, ob ein Werkzeug oder eine Fähigkeit verfügbar ist, und gibt `true` oder `false` auf der Standardausgabe aus. Der Befehl ist für Taskfile-`sh:`-Variablen als plattformübergreifende Alternative zu `command -v` vorgesehen.

Verwenden Sie `|`, um zu prüfen, ob mindestens eine von mehreren Alternativen verfügbar ist.

```bash
wails3 tool has <tool>
```

#### Beispiele

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Verwendung in einem Taskfile

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="Veraltet"}
`wails3 tool has-cc` ist veraltet. Aktualisieren Sie Ihr Taskfile und verwenden Sie stattdessen `wails3 tool has gcc|clang`.

@end

Ein abwärtskompatibler Alias für `wails3 tool has gcc|clang`. Prüft, ob `gcc` oder `clang` im PATH verfügbar ist, und gibt `true` oder `false` aus.

```bash
wails3 tool has-cc
```

## Aktualisierungsbefehle

Mit Aktualisierungsbefehlen können Sie Projektressourcen verwalten und aktualisieren. Alle Aktualisierungsbefehle verwenden folgenden Basisbefehl: `wails3 update <command>`.

### `update cli`

Aktualisiert die Wails-CLI auf eine neue Version.

```bash
wails3 update cli [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-pre` | Auf die neueste Vorabversion aktualisieren | `false` |
| `-version` | Auf eine bestimmte Version aktualisieren |  |
| `-nocolour` | Farbige Ausgabe deaktivieren | `false` |

Mit dem Befehl update cli können Sie Ihre Wails-CLI-Installation aktualisieren. Standardmäßig wird sie auf die neueste stabile Version aktualisiert. Mit der Option `-pre` können Sie auf die neueste Vorabversion aktualisieren oder mit der Option `-version` eine bestimmte Version angeben.

Aktualisieren Sie anschließend auch die Datei go.mod Ihres Projekts auf dieselbe Version:

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

Aktualisiert die Build-Ressourcen anhand der angegebenen Konfigurationsdatei.

```bash
wails3 update build-assets [flags]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-config` | Pfad zur Konfigurationsdatei |  |
| `-dir` | Ausgabeverzeichnis | `build` |
| `-silent` | Ausgabe unterdrücken | `false` |
| `-company` | Unternehmensname |  |
| `-productname` | Produktname |  |
| `-description` | Produktbeschreibung |  |
| `-version` | Produktversion |  |
| `-identifier` | Produktkennung |  |
| `-copyright` | Urheberrechtshinweis |  |
| `-comments` | Dateikommentare |  |

## Hilfsbefehle

Hilfsbefehle bieten praktische Kurzbefehle für häufige Aufgaben. Verwenden Sie diese Befehle direkt mit dem Basisbefehl: `wails3 <command>`.

### `docs`

Öffnet die Wails-Dokumentation in Ihrem Standardbrowser.

```bash
wails3 docs
```

### `releasenotes`

Zeigt die Versionshinweise für die aktuelle oder angegebene Version an.

```bash
wails3 releasenotes [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-v` | Version, deren Versionshinweise angezeigt werden sollen |  |
| `-n` | Farbausgabe deaktivieren | `false` |

### `version`

Gibt die aktuelle Version von Wails aus.

```bash
wails3 version
```

### `sponsor`

Öffnet die Sponsoring-Seite von Wails in Ihrem Standardbrowser.

```bash
wails3 sponsor

```
