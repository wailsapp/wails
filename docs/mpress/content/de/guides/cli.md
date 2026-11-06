---
title: "CLI-Flags und Beispiele"
description: "Flags und Anwendungsbeispiele für einzelne Wails-CLI-Befehle nachschlagen"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

Verwenden Sie diese Referenz, wenn Sie bereits wissen, welchen Befehl Sie benötigen. Um Befehle nach ihrem Zweck zu durchsuchen, beginnen Sie mit der [Übersicht der CLI-Befehle](/reference/cli/).

Die Wails-CLI bietet umfassende Befehle zum Entwickeln, Erstellen und Warten Ihrer Wails-Anwendungen.

## Kernbefehle

Kernbefehle sind die wichtigsten Befehle zum Erstellen, Entwickeln und Bauen von Projekten.

Alle CLI-Befehle haben das folgende Format: `wails3 <command>`.

### `init`

Initialisiert ein neues Wails-Projekt. Während der Initialisierung wird der Befehl `go mod tidy` ausgeführt, um die Projektpakete zu aktualisieren. Sie können dies umgehen, indem Sie für den Befehl `init` das Flag `-skipgomodtidy` verwenden.

```bash
wails3 init [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-p` | Go-Paketname | `main` |
| `-t` | Vorlagenname oder URL | `vanilla` |
| `-n` | Projektname |  |
| `-d` | Projektverzeichnis | `.` |
| `-q` | Ausgabe unterdrücken | `false` |
| `-l` | Vorlagen auflisten | `false` |
| `-mod` | Go-Modulpfad (wird bei Auslassung aus `-git` berechnet) |  |
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

Führt die Anwendung im Entwicklungsmodus aus. Dadurch erhalten Sie eine Live-Ansicht Ihres Frontend-Codes. Änderungen werden ohne erneutes Erstellen der gesamten Anwendung direkt in der laufenden Anwendung sichtbar. Änderungen an Ihrem Go-Code werden ebenfalls erkannt; die Anwendung wird automatisch neu erstellt und neu gestartet.

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
Dies entspricht der Ausführung von `wails3 task dev` und führt die Aufgabe `dev` im Haupt-Taskfile des Projekts aus. Sie können dies durch Bearbeiten der Datei `Taskfile.yml` anpassen.

@end

### `build`

Erstellt mithilfe der generierten Taskfiles eine Produktionsversion. Verwenden Sie `wails3 dev` für Hot Reload oder geben Sie `DEV=true` für einen Entwicklungsbuild an. Standardmäßig wird für die aktuelle Plattform und Architektur gebaut.

```bash
wails3 build [flags] [CLI variables...]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-tags` | Zusätzliche Go-Build-Tags (durch Kommas getrennt) |  |

Sie können CLI-Variablen übergeben, um den Build anzupassen:

```bash
wails3 build GOOS=linux
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

Die Tags werden als `EXTRA_TAGS` an das zugrunde liegende Taskfile weitergegeben.

@note{type="info"}
Dies entspricht der Ausführung von `wails3 task build`, wodurch die Aufgabe `build` im Haupt-Taskfile des Projekts ausgeführt wird. Alle an `build` übergebenen CLI-Variablen werden an die zugrunde liegende Aufgabe weitergegeben. Sie können den Build-Prozess durch Bearbeiten der Datei `Taskfile.yml` anpassen.

@end

### `package`

Erstellt plattformspezifische Pakete für die Verteilung.

```bash
wails3 package [CLI variables...]
```

Die mitgelieferten Paketierungs-Tasks akzeptieren `ARCH` zur Auswahl der Zielarchitektur:

```bash
wails3 package ARCH=amd64
```

`VERSION` und `OUTPUT` wirken sich nur auf ein benutzerdefiniertes Taskfile aus, das diese Variablen ausdrücklich ausliest; sie überschreiben weder die Paketversion noch den Dateinamen der mitgelieferten Paketierungs-Tasks.

#### Pakettypen

Für jede Plattform sind die folgenden Pakettypen verfügbar:

| Plattform | Pakettyp |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
Dies entspricht `wails3 task package`, das den Task `package` in der zentralen Taskfile des Projekts ausführt. Alle an `package` übergebenen CLI-Variablen werden an den zugrunde liegenden Task weitergeleitet. Sie können den Paketierungsprozess anpassen, indem Sie die Datei `Taskfile.yml` bearbeiten.

@end

Die standardmäßigen Build-Tasks lesen `GOOS`, `ARCH` und `DEV`. `PLATFORM` und `CONFIG` im nachstehenden benutzerdefinierten Task sind benutzerdefinierte Beispiele.

### `task`

Führt Tasks aus, die in der Taskfile.yml Ihres Projekts definiert sind. Dies ist eine eingebettete Version von [Taskfile](https://taskfile.dev), mit der Sie benutzerdefinierte Build-, Test- und Deployment-Tasks definieren und ausführen können.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI-Variablen

Sie können Variablen im Format `KEY=VALUE` an Tasks übergeben:

```bash
wails3 task build GOOS=linux
wails3 task deploy ENV=staging VERSION=1.2.3
```

Auf diese Variablen können Sie in Ihrer Taskfile.yml mit der Go-Template-Syntax zugreifen:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-h` | Zeigt die Verwendung von Task an | `false` |
| `-i` | Erstellt eine neue Taskfile.yml | `false` |
| `-list` | Listet Tasks mit Beschreibungen auf | `false` |
| `-list-all` | Listet alle Tasks auf (mit oder ohne Beschreibungen) | `false` |
| `-json` | Formatiert die Task-Liste als JSON | `false` |
| `-status` | Beendet den Befehl mit einem von null verschiedenen Exit-Code, wenn der Task nicht aktuell ist | `false` |
| `-f` | Erzwingt die Ausführung, auch wenn der Task aktuell ist | `false` |
| `-w` | Aktiviert den Überwachungsmodus für den angegebenen Task | `false` |
| `-v` | Aktiviert den ausführlichen Modus | `false` |
| `-version` | Gibt die Task-Version aus | `false` |
| `-s` | Deaktiviert die Befehlsausgabe | `false` |
| `-p` | Führt Tasks parallel aus | `false` |
| `-dry` | Kompiliert Tasks und gibt sie aus, ohne sie auszuführen | `false` |
| `-summary` | Zeigt eine Zusammenfassung zu einem Task an | `false` |
| `-x` | Reicht den Exit-Code des Tasks durch | `false` |
| `-dir` | Legt das Ausführungsverzeichnis fest |  |
| `-taskfile` | Legt fest, welche Taskfile ausgeführt wird |  |
| `-output` | Legt den Ausgabestil fest: [interleaved|group|prefixed] |  |
| `-c` | Farbige Ausgabe (standardmäßig aktiviert) | `true` |
| `-C` | Begrenzt die Anzahl gleichzeitig ausgeführter Tasks |  |
| `-interval` | Intervall zur Überwachung auf Änderungen (in Sekunden) |  |

#### Beispiele

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build GOOS=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

Startet den MCP-Server des Wails-Projekts für die agentengestützte Projektverwaltung. Dieser ist von dem MCP-Server getrennt, der in eine laufende Anwendung kompiliert wird: `wails3 mcp` verwaltet Projektdateien und Lebenszyklusbefehle, während der MCP-Server der Anwendung die laufende WebView steuert.

```bash
wails3 mcp [flags]
```

Der Transport wird automatisch ausgewählt:

- Wenn ein MCP-Host Wails mit umgeleiteter Standardeingabe und Standardausgabe startet, verwendet der Server **stdio**.
- Bei interaktiver Ausführung in einem Terminal verwendet der Server **Streamable HTTP** auf `127.0.0.1` und fordert beim Betriebssystem einen freien Port an.

Verwenden Sie `--stdio` oder `--http`, um einen Transport explizit auszuwählen. Verwenden Sie `--port 0`, um im HTTP-Modus einen freien Loopback-Port auszuwählen.

#### MCP-Flags

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `--root` | Zulässiges Projektstammverzeichnis. Pfade und symbolische Links, die aus diesem Verzeichnis herausführen, werden abgelehnt. | Aktuelles Verzeichnis |
| `--token` | Sitzungs-/Bearer-Token für verändernde Tools und Tools zur Prozesssteuerung. Fällt auf `WAILS_MCP_TOKEN` zurück. | Sicher generiert |
| `--stdio` | Erzwingt den stdio-Transport. | Automatisch |
| `--http` | Erzwingt den Transport über Streamable HTTP. | Automatisch |
| `--port` | HTTP-Port; `0` wählt einen freien Loopback-Port aus. | `0` |

Im HTTP-Modus gibt Wails den Endpunkt und das Bearer-Token auf stderr aus. Im stdio-Modus ist das Token in den MCP-Initialisierungsanweisungen enthalten. Der Server ermöglicht keine beliebige Ausführung von Shell-Befehlen. Remote-Vorlagen und Git-Remotes erfordern eine ausdrückliche Genehmigung über die Eingabe `allowExternal` des Tools.

### `doctor`

Führt eine Systemprüfung durch und zeigt einen Statusbericht an.

```bash
wails3 doctor
```

## Generierungsbefehle

Generierungsbefehle erstellen verschiedene Projektressourcen wie Bindings, Symbole und Build-Dateien. Alle Generierungsbefehle verwenden den Basisbefehl `wails3 generate <command>`.

### `generate bindings`

Generiert Bindings und Modelle für deinen Go-Code.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-f` | Zusätzliche Go-Build-Optionen |  |
| `-d` | Ausgabeverzeichnis | `frontend/bindings` |
| `-models` | Dateiname für Modelle | `models` |
| `-index` | Dateiname für den Index | `index` |
| `-ts` | TypeScript generieren | `false` |
| `-i` | TS-Schnittstellen verwenden | `false` |
| `-b` | Mitgelieferte Runtime verwenden | `false` |
| `-names` | Namen anstelle von IDs verwenden | `false` |
| `-noindex` | Indexdateien überspringen | `false` |
| `-noevents` | Generierung ereignisbezogener Bindings überspringen | `false` |
| `-dry` | Testlauf | `false` |
| `-silent` | Stiller Modus | `false` |
| `-v` | Debug-Ausgabe | `false` |
| `-clean` | Ausgabeverzeichnis vor der Generierung bereinigen | `true` |

### `generate build-assets`

Generiert Build-Ressourcen für deine Anwendung.

```bash
wails3 generate build-assets [flags]
```

#### Flags

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

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-input` | PNG-Eingabedatei | Erforderlich |
| `-windowsfilename` | Windows-Ausgabedateiname |  |
| `-macfilename` | macOS-Ausgabedateiname |  |
| `-sizes` | Symbolgrößen (durch Kommas getrennt) | `256,128,64,48,32,16` |
| `-example` | Beispielsymbol generieren | `false` |
| `-iconcomposerinput` | Icon-Composer-Eingabedatei (`.icon`) |  |
| `-macassetdir` | Ausgabeverzeichnis für Mac-Ressourcen (Assets.car + icns) |  |

#### Icon Composer (macOS)

Unter macOS 26+ kannst du aus Icon-Composer-Dateien vom Typ `.icon` die Dateien `Assets.car` und `icons.icns` generieren:

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Dabei wird die Datei `.icon` mit Apples Befehl `actool` kompiliert. Erfordert Xcode mit `actool` ab Version 26.

Wenn du Icon Composer verwendest, setze `cfBundleIconName` in deiner Datei `build/config.yml` auf den Dateinamen von `.icon` ohne Erweiterung:

```yaml
info:
  cfBundleIconName: "appicon"
```

Wenn dies nicht festgelegt ist und `Assets.car` vorhanden ist, wird standardmäßig `"appicon"` verwendet.

### `generate syso`

Generiert eine Windows-.syso-Datei.

```bash
wails3 generate syso [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-manifest` | Pfad zur Manifestdatei | Erforderlich |
| `-icon` | Pfad zur Symboldatei | Erforderlich |
| `-info` | Pfad zur Datei mit den Versionsinformationen |  |
| `-arch` | Zielarchitektur | Aktuelle GOARCH |
| `-out` | Name der Ausgabedatei | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Generiert eine Linux-.desktop-Datei.

```bash
wails3 generate .desktop [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-name` | Anwendungsname | Erforderlich |
| `-exec` | Pfad zur ausführbaren Datei | Erforderlich |
| `-icon` | Pfad zum Symbol |  |
| `-categories` | Anwendungskategorien | `Utility` |
| `-comment` | Anwendungskommentar |  |
| `-terminal` | Im Terminal ausführen | `false` |
| `-keywords` | Suchbegriffe |  |
| `-version` | Anwendungsversion |  |
| `-genericname` | Generischer Name |  |
| `-startupnotify` | Startbenachrichtigung anzeigen | `false` |
| `-mimetype` | Unterstützte MIME-Typen |  |
| `-output` | Name der Ausgabedatei | `[name].desktop` |

### `generate runtime`

Generiert die vorgefertigte Version der Runtime.

```bash
wails3 generate runtime
```

### `generate constants`

Generiert JavaScript-Konstanten aus Go-Code.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

Generiert ein Windows-WebView2-Bootstrap-Installationsprogramm zur Verteilung.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

Erstellt das Grundgerüst für ein neues Projektvorlagenverzeichnis.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Generiert ein Linux-AppImage.

```bash
wails3 generate appimage [flags]
```

#### Flags

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-binary` | Pfad zur Binärdatei | Erforderlich |
| `-icon` | Pfad zur Symboldatei | Erforderlich |
| `-desktop` | Pfad zur .desktop-Datei | Erforderlich |
| `-builddir` | Build-Verzeichnis | Temporäres Verzeichnis |
| `-output` | Ausgabeverzeichnis | `.` |

## Dienstbefehle

Dienstbefehle unterstützen bei der Verwaltung von Wails-Diensten. Alle Dienstbefehle verwenden folgenden Basisbefehl: `wails3 service <command>`.

### `service init`

Initialisiert einen neuen Dienst.

```bash
wails3 service init [flags]
```

#### Flags

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

Werkzeugbefehle stellen Hilfsfunktionen für Entwicklung und Debugging bereit. Alle Werkzeugbefehle verwenden folgenden Basisbefehl: `wails3 tool <command>`.

### `tool checkport`

Prüft, ob ein Port geöffnet ist. Nützlich, um zu testen, ob vite ausgeführt wird.

```bash
wails3 tool checkport [flags]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-port` | Zu prüfender Port | `9245` |
| `-host` | Zu prüfender Host | `localhost` |

### `tool watcher`

Überwacht Dateien und führt bei Änderungen einen Befehl aus.

```bash
wails3 tool watcher [flags]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
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

Erhöht anhand der angegebenen Flags eine semantische Version.

```bash
wails3 tool version [flags]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-v` | Zu erhöhende aktuelle Version |  |
| `-major` | Hauptversion erhöhen | `false` |
| `-minor` | Nebenversion erhöhen | `false` |
| `-patch` | Patch-Version erhöhen | `false` |
| `-prerelease` | Vorabversion erhöhen (z. B. von alpha.5 auf alpha.6) | `false` |

Der Befehl verwendet folgende Rangfolge: Hauptversion > Nebenversion > Patch-Version > Vorabversion. Er behält ein vorhandenes Präfix „v“ der Eingabeversion sowie sämtliche Vorabversions- und Metadatenkomponenten bei.

Verwendungsbeispiel:

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.1.0-alpha.5 -prerelease # Output: v3.1.0-alpha.6
```

### `tool package`

Erzeugt Linux-Pakete (deb, rpm, archlinux).

```bash
wails3 tool package [flags]
```

#### Optionen

| Flag | Beschreibung | Standardwert |
| --- | --- | --- |
| `-format` | Paketformat (deb, rpm, archlinux) | `deb` |
| `-name` | Name der ausführbaren Datei | `myapp` |
| `-config` | Pfad zur Konfigurationsdatei |  |
| `-out` | Ausgabeverzeichnis | `.` |

### `tool lipo`

Erstellt eine universelle macOS-Binärdatei, indem architekturspezifische Binärdateien kombiniert werden.

```bash
wails3 tool lipo [flags]
```

#### Flags

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-output` | Pfad zur ausgegebenen Binärdatei |  |

### `tool capabilities`

Prüft die Build-Fähigkeiten des Systems (Verfügbarkeit von GTK4/GTK3 unter Linux).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

Erzeugt Flags zum Einbinden von Docker-Volumes für die Cross-Kompilierung. Gibt `-v`-Flags für den Go-Modul-Cache und alle lokalen `replace`-Direktiven in `go.mod` zur Verwendung in `docker run`-Befehlen des Taskfile aus.

```bash
wails3 tool docker-mounts
```

### `tool has`

Prüft, ob ein Werkzeug oder eine Fähigkeit verfügbar ist, und gibt `true` oder `false` auf stdout aus. Vorgesehen für die Verwendung in `sh:`-Variablen eines Taskfile als plattformübergreifende Alternative zu `command -v`.

Verwende `|`, um zu prüfen, ob eine von mehreren Alternativen verfügbar ist.

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
`wails3 tool has-cc` ist veraltet. Aktualisiere dein Taskfile, sodass es stattdessen `wails3 tool has gcc|clang` verwendet.

@end

Ein abwärtskompatibler Alias für `wails3 tool has gcc|clang`. Prüft, ob `gcc` oder `clang` im PATH verfügbar ist, und gibt `true` oder `false` aus.

```bash
wails3 tool has-cc
```

## Aktualisierungsbefehle

Aktualisierungsbefehle helfen dabei, Projektressourcen zu verwalten und zu aktualisieren. Alle Aktualisierungsbefehle verwenden den Basisbefehl `wails3 update <command>`.

### `update cli`

Aktualisiert die Wails CLI auf eine neue Version.

```bash
wails3 update cli [flags]
```

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-pre` | Auf die neueste Vorabversion aktualisieren | `false` |
| `-version` | Auf eine bestimmte Version aktualisieren |  |
| `-nocolour` | Farbige Ausgabe deaktivieren | `false` |

Mit dem Befehl update cli kannst du deine Wails-CLI-Installation aktualisieren. Standardmäßig wird sie auf die neueste stabile Version aktualisiert. Mit dem Flag `-pre` kannst du auf die neueste Vorabversion aktualisieren oder mit dem Flag `-version` eine bestimmte Version angeben.

Aktualisiere anschließend die Datei go.mod deines Projekts auf dieselbe Version:

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

Aktualisiert die Build-Ressourcen anhand der angegebenen Konfigurationsdatei.

```bash
wails3 update build-assets [flags]
```

#### Flags

| Option | Beschreibung | Standardwert |
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

#### Optionen

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-v` | Version, deren Versionshinweise angezeigt werden sollen |  |
| `-n` | Farbausgabe deaktivieren | `false` |

### `version`

Gibt die aktuelle Wails-Version aus.

```bash
wails3 version
```

### `sponsor`

Öffnet die Wails-Sponsoringseite in Ihrem Standardbrowser.

```bash
wails3 sponsor

```
