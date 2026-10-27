---
title: "Übersicht der CLI-Befehle"
description: "Wails-CLI-Befehle nach Zweck finden: Projekte, Builds, Signierung, Werkzeuge und Dienste"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

Hier finden Sie Befehle nach Zweck. Befehlsspezifische Flags und Anwendungsbeispiele finden Sie unter [CLI-Flags und Beispiele](/guides/cli/).

## Übersicht

Die Wails-CLI (`wails3`) ist der Befehlszeileneinstieg zum Erstellen, Entwickeln, Bauen, Signieren, Paketieren und Untersuchen von Wails-3-Anwendungen. Der Großteil der Build-Orchestrierung wird an projektspezifische Taskfiles delegiert, die sich in Ihrem Projekt unter `build/` befinden. Viele `wails3`-Befehle sind dünne Wrapper, die eine bestimmte Task aufrufen.

Führen Sie den folgenden Befehl aus, um die jeweils aktuellste Hilfe zu einem Befehl anzuzeigen:

```bash
wails3 --help
wails3 <command> --help
```

## Projektlebenszyklus

| Befehl | Beschreibung |
| --- | --- |
| `wails3 init` | Erstellt ein neues Projekt aus einer Vorlage. Flags: `-n` (Projektname), `-t` (Vorlage, Standard: `vanilla`), `-p` (Go-Paketname, Standard: `main`), `-d` (Projektverzeichnis, Standard: `.`), `-q` (stille Ausgabe), `-l` (Vorlagen auflisten), `-mod` (Go-Modulpfad), `--git` (URL des Git-Repositorys), `--skipgomodtidy`, `-s` (Warnung zu Remote-Vorlagen überspringen), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Führt die Anwendung im Entwicklungsmodus mit Hot Reload des Frontends aus. Flags: `--config` (Standard: `./build/config.yml`), `--port` (Vite-Entwicklungsport), `-s` (HTTPS aktivieren). |
| `wails3 build` | Baut das Projekt. Dünner Wrapper um die Taskfile-Task `build`. Flags: `--tags` (wird als `EXTRA_TAGS=` weitergereicht), `--obfuscated` (Build mit Garble; siehe [Verschleierte Builds](/guides/build/obfuscation/)), `--garbleargs` (zusätzliche Flags, die vor dem Unterbefehl `build` an `garble` weitergereicht werden). |
| `wails3 package` | Führt die plattformspezifische Taskfile-Task `package` aus. |
| `wails3 task [name]` | Führt eine beliebige Taskfile-Task aus. Ohne Namen zeigt `--list` alle registrierten Tasks an. |
| `wails3 mcp` | Startet den MCP-Server des Projekts. Verwendet für von Agenten gestartete Prozesse automatisch stdio und bei interaktiver Terminalnutzung Streamable HTTP über die Loopback-Schnittstelle. |
| `wails3 doctor` | Gibt einen Diagnosebericht zu Ihrer Umgebung aus. |
| `wails3 doctor-ng` | Neuere TUI-Variante von `doctor`. |
| `wails3 version` | Gibt die CLI-Version aus. |
| `wails3 releasenotes` | Gibt die aktuellen Versionshinweise aus. |
| `wails3 docs` | Öffnet die Dokumentationswebsite im Browser. |
| `wails3 sponsor` | Öffnet die Sponsorenseite. |

## Generierung

`wails3 generate <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `generate bindings` | Generiert Bindings von Go zum Frontend. Flags: `-d` (Ausgabeverzeichnis), `-models`, `-index`, `-ts`, `-i` (Schnittstellen), `-b` (Bundle), `-names` (`Call.ByName` ausgeben), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (Standard: `true`), `-f`, `-obfuscated` (`wails_obfuscated.gen.go` mit stabilen Binding-IDs für Garble-Builds generieren; siehe [Verschleierte Builds](/guides/build/obfuscation/)), `-obfuscated-output` (Verzeichnis für die generierte Datei; standardmäßig das Verzeichnis des Hauptpakets). Akzeptiert Paketmuster (z. B. `./...`); wenn keine angegeben sind, wird das aktuelle Verzeichnis verwendet. |
| `generate icons` | Konvertiert eine Quell-PNG in die plattformspezifischen Symbolformate. Flags: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Generiert aus `build/config.yml` den Inhalt des Verzeichnisses `build/` (Taskfile-Fragmente, NSIS-Dateien, `Info.plist`, `.desktop`-Vorlage usw.). |
| `generate runtime` | Generiert die vorgefertigte Datei `/wails/runtime.js` neu, die an die Webview ausgeliefert wird. |
| `generate syso` | Generiert die Windows-Ressourcendatei `.syso` (Symbol + Manifest + Versionsinformationen). |
| `generate webview2bootstrapper` | Generiert ein WebView2-Bootstrap-Installationsprogramm für Windows. |
| `generate constants` | Generiert JS-Konstanten für Ereignisnamen aus Go-Ereignistypen. |
| `generate template` | Erstellt das Grundgerüst für eine neue Projektvorlage. |
| `generate .desktop` | Generiert eine Linux-Datei `.desktop` (von AppImage/DEB/RPM verwendet). |
| `generate appimage` | Generiert das AppImage-Build-Verzeichnis. |

## Aktualisierung

`wails3 update <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `update build-assets` | Aktualisiert das Verzeichnis `build/` aus `build/config.yml` und behält Benutzeränderungen nach Möglichkeit bei. |
| `update cli` | Aktualisiert die Binärdatei `wails3` selbstständig. |

## Codesignierung und Paketierung

| Befehl | Beschreibung |
| --- | --- |
| `wails3 signing` | Interaktiver Assistent, der die Signierung für die in `build/` erkannten Plattformen konfiguriert. Flags: `--platform` (wiederholbar; standardmäßig automatische Erkennung anhand des Build-Verzeichnisses). |
| `wails3 entitlements` | Interaktiver Assistent für macOS-Berechtigungen. Flag: `--output` (Pfad; Standard: `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Wrapper, der die plattformspezifische Taskfile-Task `*:sign` für das aktuelle Betriebssystem oder das über `GOOS` angegebene Betriebssystem ausführt. |
| `wails3 tool sign` | Direkter Signierungseinstieg auf niedriger Ebene. Flags: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

In dieser Version stellt der CLI-Parser `signing` und `entitlements` auf der obersten Ebene bereit. Prüfen Sie `wails3 signing --help` und `wails3 entitlements --help`; `wails3 setup signing` wählt den Signierungsassistenten nicht aus. Verwenden Sie zur direkten Verwaltung von Anmeldedaten `xcrun notarytool store-credentials` oder `gpg`.

## Werkzeuge

`wails3 tool <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `tool checkport` | Prüft, ob ein TCP-Port geöffnet ist (nützlich, um auf Vite zu warten). |
| `tool watcher` | Führt einen Befehl aus, sobald sich überwachte Dateien ändern. |
| `tool cp` | Kopiert Dateien plattformübergreifend. |
| `tool buildinfo` | Gibt die in einer Binärdatei eingebetteten Go-Buildinformationen aus. |
| `tool package` | Erstellt aus `build/linux/nfpm` ein Linux-Paket (`deb`, `rpm`, `archlinux`). |
| `tool version` | Erhöht die semantische Version eines Projekts. |
| `tool lipo` | Fasst Binärdateien für mehrere macOS-Architekturen zu einer universellen Binärdatei zusammen. |
| `tool capabilities` | Prüft das System auf die Verfügbarkeit von GTK3/GTK4 und WebKit. |
| `tool sign` | (Siehe [Codesignierung und Paketierung](#codesignierung-und-paketierung).) |

## Dienste

`wails3 service <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `service init` | Erstellt das Grundgerüst für ein neues Dienstpaket. |

## iOS

`wails3 ios <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `ios overlay:gen` | Generiert das Go-Overlay für den iOS-Bridge-Shim. |
| `ios xcode:gen` | Generiert ein Xcode-Projekt im Ausgabeverzeichnis. |

## Pfade der Build-Ausgaben

- Native Binärdateien werden in `bin/<APP_NAME>` abgelegt (unter Windows in `bin/<APP_NAME>.exe`). `build/bin/` ist nicht vorhanden.
- Paketierte Ausgaben (`.app`, `.dmg`, NSIS-Installationsprogramm, MSIX, DEB/RPM/AppImage) werden ebenfalls in `bin/` abgelegt (oder in plattformspezifischen Unterverzeichnissen, die der jeweilige Taskfile-Task erstellt).

## Gemeinsame Flags

| Flag | Gilt für | Beschreibung |
| --- | --- | --- |
| `-nocolour` | Befehle, die es in `--help` bereitstellen (zum Beispiel `init` und `update cli`) | Deaktiviert ANSI-Farben in der CLI-Ausgabe. Das Flag wird nicht von jedem Befehl akzeptiert, darunter `version`. |

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) an.
