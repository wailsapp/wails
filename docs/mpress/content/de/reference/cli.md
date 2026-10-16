---
title: "CLI-Referenz"
description: "Vollständige Referenz der Wails-CLI-Befehle"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## Übersicht

Die Wails-CLI (`wails3`) ist der Befehlszeileneinstiegspunkt zum Erstellen, Entwickeln, Bauen, Signieren, Paketieren und Untersuchen von Wails-3-Anwendungen. Der Großteil der Build-Orchestrierung wird an projektspezifische Taskfiles delegiert, die sich in Ihrem Projekt unter `build/` befinden. Viele `wails3`-Befehle sind schlanke Wrapper, die einen bestimmten Task aufrufen.

Führen Sie für die aktuellste Hilfe zu einem Befehl Folgendes aus:

```bash
wails3 --help
wails3 <command> --help
```

## Projektlebenszyklus

| Befehl | Beschreibung |
| --- | --- |
| `wails3 init` | Erstellt ein neues Projekt aus einer Vorlage. Flags: `-n` (Projektname), `-t` (Vorlage, Standardwert: `vanilla`), `-p` (Go-Paketname, Standardwert: `main`), `-d` (Projektverzeichnis, Standardwert: `.`), `-q` (stiller Modus), `-l` (Vorlagen auflisten), `-mod` (Go-Modulpfad), `--git` (URL des Git-Repositorys), `--skipgomodtidy`, `-s` (Warnung vor Remote-Vorlagen überspringen), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Führt die Anwendung im Entwicklungsmodus mit Hot Reload des Frontends aus. Flags: `--config` (Standardwert: `./build/config.yml`), `--port` (Vite-Entwicklungsport), `-s` (HTTPS aktivieren). |
| `wails3 build` | Baut das Projekt. Schlanker Wrapper um den Taskfile-Task `build`. Flags: `--tags` (wird als `EXTRA_TAGS=` weitergegeben), `--obfuscated` (Build mit Garble; siehe [Verschleierte Builds](/guides/build/obfuscation/)), `--garbleargs` (zusätzliche Flags, die vor dem Unterbefehl `build` an `garble` weitergegeben werden). |
| `wails3 package` | Führt den plattformspezifischen Taskfile-Task `package` aus. |
| `wails3 task [name]` | Führt einen beliebigen Taskfile-Task aus. Ohne Namen zeigt `--list` alle registrierten Tasks an. |
| `wails3 mcp` | Startet den MCP-Server des Projekts. Verwendet für von Agenten gestartete Prozesse automatisch stdio und bei interaktiver Terminalnutzung Streamable HTTP über die Loopback-Schnittstelle. |
| `wails3 doctor` | Gibt einen Diagnosebericht zu Ihrer Umgebung aus. |
| `wails3 doctor-ng` | Neuere TUI-Variante von `doctor`. |
| `wails3 version` | Gibt die CLI-Version aus. |
| `wails3 releasenotes` | Gibt die neuesten Versionshinweise aus. |
| `wails3 docs` | Öffnet die Dokumentationswebsite in Ihrem Browser. |
| `wails3 sponsor` | Öffnet die Sponsorenseite. |

## Generierung

`wails3 generate <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `generate bindings` | Generiert Bindings von Go zum Frontend. Flags: `-d` (Ausgabeverzeichnis), `-models`, `-index`, `-ts`, `-i` (Schnittstellen), `-b` (Bundle), `-names` (`Call.ByName` ausgeben), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (Standardwert: `true`), `-f`, `-obfuscated` (`wails_obfuscated.gen.go` mit stabilen Binding-IDs für Garble-Builds generieren; siehe [Verschleierte Builds](/guides/build/obfuscation/)), `-obfuscated-output` (Verzeichnis für die generierte Datei; standardmäßig das Verzeichnis des Hauptpakets). Akzeptiert Paketmuster (z. B. `./...`); wenn keine angegeben sind, wird das aktuelle Verzeichnis verwendet. |
| `generate icons` | Konvertiert eine Quell-PNG-Datei in die plattformspezifischen Symbolformate. Flags: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Generiert den Inhalt des Verzeichnisses `build/` (Taskfile-Fragmente, NSIS-Dateien, `Info.plist`, `.desktop`-Vorlage usw.) aus `build/config.yml`. |
| `generate runtime` | Generiert das vorgefertigte, an die Webview ausgelieferte `/wails/runtime.js` erneut. |
| `generate syso` | Generiert die Windows-Ressourcendatei `.syso` (Symbol + Manifest + Versionsinformationen). |
| `generate webview2bootstrapper` | Generiert ein WebView2-Bootstrap-Installationsprogramm für Windows. |
| `generate constants` | Generiert aus Go-Ereignistypen JS-Konstanten für Ereignisnamen. |
| `generate template` | Erstellt das Grundgerüst für eine neue Projektvorlage. |
| `generate .desktop` | Generiert eine Linux-Datei `.desktop` (wird von AppImage/DEB/RPM verwendet). |
| `generate appimage` | Generiert das AppImage-Build-Verzeichnis. |

## Aktualisierung

`wails3 update <subcommand>`:

| Unterbefehl | Beschreibung |
| --- | --- |
| `update build-assets` | Aktualisiert das Verzeichnis `build/` aus `build/config.yml` und behält Benutzeränderungen nach Möglichkeit bei. |
| `update cli` | Aktualisiert die Binärdatei `wails3` selbst. |

## Codesignierung und Paketierung

| Befehl | Beschreibung |
| --- | --- |
| `wails3 setup signing` | Interaktiver Assistent, der die Signierung für die in `build/` erkannten Plattformen konfiguriert. Flags: `--platform` (wiederholbar; standardmäßig automatische Erkennung anhand des Build-Verzeichnisses). |
| `wails3 setup entitlements` | Interaktiver Assistent für macOS-Berechtigungen. Flag: `--output` (Pfad; Standardwert: `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Wrapper, der die plattformspezifische Taskfile-Aufgabe `*:sign` für das aktuelle Betriebssystem ausführt (oder für das über `GOOS` angegebene). |
| `wails3 tool sign` | Direkter Low-Level-Einstiegspunkt für die Signierung. Flags: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

Es gibt **kein** Unterkommando `wails3 signing`. Verwenden Sie für Anmeldedaten im Schlüsselbund `xcrun notarytool store-credentials` und für PGP-Schlüssel direkt `gpg` (der Assistent `wails3 setup signing` automatisiert beides).

## Werkzeuge

`wails3 tool <subcommand>`:

| Unterkommando | Beschreibung |
| --- | --- |
| `tool checkport` | Prüfen, ob ein TCP-Port geöffnet ist (nützlich beim Warten auf Vite). |
| `tool watcher` | Bei jeder Änderung an überwachten Dateien einen Befehl ausführen. |
| `tool cp` | Plattformübergreifendes Kopieren von Dateien. |
| `tool buildinfo` | Die eingebetteten Go-Buildinformationen einer Binärdatei ausgeben. |
| `tool package` | Aus `build/linux/nfpm` ein Linux-Paket (`deb`, `rpm`, `archlinux`) erstellen. |
| `tool version` | Die semantische Version eines Projekts erhöhen. |
| `tool lipo` | Binärdateien für mehrere macOS-Architekturen zu einer universellen Binärdatei zusammenführen. |
| `tool capabilities` | Das System auf die Verfügbarkeit von GTK3/GTK4 und WebKit prüfen. |
| `tool sign` | (Siehe [Codesignierung und Paketierung](#codesignierung-und-paketierung).) |

## Dienste

`wails3 service <subcommand>`:

| Unterkommando | Beschreibung |
| --- | --- |
| `service init` | Das Grundgerüst eines neuen Dienstpakets erstellen. |

## iOS

`wails3 ios <subcommand>`:

| Unterkommando | Beschreibung |
| --- | --- |
| `ios overlay:gen` | Das Go-Overlay für den iOS-Bridge-Shim generieren. |
| `ios xcode:gen` | Ein Xcode-Projekt im Ausgabeverzeichnis generieren. |

## Pfade der Build-Ausgaben

- Native Binärdateien werden in `bin/<APP_NAME>` abgelegt (unter Windows in `bin/<APP_NAME>.exe`). `build/bin/` ist nicht vorhanden.
- Paketierte Ausgaben (`.app`, `.dmg`, NSIS-Installationsprogramm, MSIX, DEB/RPM/AppImage) werden ebenfalls in `bin/` abgelegt (oder in plattformspezifischen Unterverzeichnissen, die von der jeweiligen Taskfile-Aufgabe erstellt werden).

## Globale Flags

| Flag | Gilt für | Beschreibung |
| --- | --- | --- |
| `--no-colour` | Alle Befehle | ANSI-Farben in der CLI-Ausgabe deaktivieren. |

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
