---
title: "YAML-Builds"
description: "Konfigurieren und verwenden Sie das experimentelle native Buildsystem von Wails v3 mit einer einzigen wails.yaml-Datei."
sidebar:
  order: 2
  badge:
    text: "Experimentell"
    variant: caution
---

@note{type="caution" title="Experimentelle Vorschau"}
Die erste unterstützte Version akzeptiert ausschließlich literale Konfiguration. Dieses Experiment erprobt außerdem zeitlich begrenzte Lebenszyklus-Hooks, damit deren Gestaltung bewertet werden kann, bevor sie für eine Veröffentlichung vorgeschlagen werden. Benutzerdefinierte Phasen, YAML-Ausdrücke und typisierte Wails-Toolaufrufe bleiben zurückgestellt. Bestehende Projekte verwenden weiterhin Taskfiles, bis Sie eine geprüfte Migration ausdrücklich aktivieren.
@end

## YAML-Builds aktivieren

Setzen Sie `WAILS_EXP_USE_WAKE` in jedem Terminal, in dem die experimentellen Befehle ausgeführt werden. Für Bash oder Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Für PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Die Variable aktiviert die Funktion allein durch ihre Existenz: Selbst ein leerer Wert oder `false` aktiviert sie. Entfernen Sie die Variable (mit `unset WAILS_EXP_USE_WAKE` in Bash/Zsh oder `Remove-Item Env:WAILS_EXP_USE_WAKE` in PowerShell), um das Experiment zu deaktivieren. Sobald `wails.yaml` aktiv ist, müssen Befehle Wake erneut aktivieren und greifen niemals auf Taskfiles zurück.
Ohne die Variable sind das Build-/Dev-Routing über YAML und die YAML-Projektgenerierung deaktiviert; außerdem sind die Befehle `migrate`, `config` und `eject` nicht verfügbar.
Noch nicht migrierte Projekte werden weiterhin über die separate eingebettete Task-Laufzeit verarbeitet.

## Ein Projekt konfigurieren

Fügen Sie bei aktivierter Opt-in-Funktion `wails.yaml` im Projektstammverzeichnis hinzu. Das Vorhandensein der Datei wählt das native Buildsystem aus und veranlasst Wails, sämtliche Taskfiles in diesem Projekt zu ignorieren. Lassen Sie die Opt-in-Funktion für alle folgenden Befehle aktiviert. Die Schema-Modeline ermöglicht in Editoren mit Unterstützung für YAML-Sprachserver die Vervollständigung und Inline-Validierung.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Führen Sie den standardmäßigen Host-Build oder die vollständige benannte Release-Anforderung aus:

```bash
wails3 build
wails3 build release
```

Wählen Sie für einen einmaligen Build kommagetrennte Ziele und Formate aus, ohne ein Profil hinzuzufügen:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Desktopanwendungen crosskompilieren

Setzen Sie für das experimentelle, ausschließlich auf Go basierende Desktop-Backend `CGO_ENABLED: "0"` in der Buildumgebung. Wake wählt für macOS- und Linux-GUI-Programme das Tag `purego` und verwendet Go direkt auf Windows-, macOS- und Linux-Hosts. In diesem Modus erfordert das Erstellen von Desktop-Programmdateien weder Docker noch Emulation, einen C-Compiler oder ein SDK der Zielplattform. Auch die Anwendungsabhängigkeiten müssen `CGO_ENABLED=0` unterstützen.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

Der zweite Befehl kombiniert beide macOS-Architekturen und erstellt mit Go ein unsigniertes DMG. Er kann auf jedem unterstützten Desktop-Host ausgeführt werden. Codesignierung und Apple-Notarisierung bleiben separate macOS-Vorgänge.

Der aufgelöste Modus wird für jede Programmdatei unabhängig ausgewählt. Ein deklariertes `build.environment.CGO_ENABLED`, eine Zielumgebung oder die Umgebung der Binärdatei überschreibt den Wert des aufrufenden Terminals. Andernfalls verwenden ein ausdrücklich ausgewähltes `purego`-Tag, ein CLI-/Serverprogramm oder ein Windows-Ziel standardmäßig `0`; macOS- und Linux-GUI-Programme verwenden standardmäßig `1`. `--plan --json` meldet für jede Programmdatei den aufgelösten Compiler, den CGo-Modus und die Tags. Bei der Crosskompilierung werden Zieldateien erstellt; zum Ausführen eines inkompatiblen Desktopziels ist weiterhin der entsprechende Host erforderlich.

Mit `CGO_ENABLED=1` verwendet Wake zunächst eine kompatible native Toolchain oder den konfigurierten `CC`-Crosscompiler, anschließend – sofern verfügbar – Zig für einen Windows-Cross-Build und danach, soweit unterstützt, den vorhandenen Docker-/Podman-Fallback. Ein explizites `CC` muss den Zielcompiler, die Bibliotheken und das Sysroot bereitstellen; Zig allein stellt kein Linux-GTK-/WebKit-Sysroot bereit. Für Paketierung und Signierung können wie nachfolgend beschrieben zusätzliche Werkzeuge erforderlich sein. Mobile Ziele erfordern CGo und die jeweiligen Plattform-SDKs.

Bei einer reinen Go-Kompilierung werden die Systembibliotheken des Zielsystems nicht mitgeliefert.
Linux-GUI-Anwendungen einschließlich ausführbarer GUI-Begleitprogramme benötigen weiterhin die durch ihre Build-Tags ausgewählte GTK/WebKitGTK-Laufzeitumgebung. Sie muss über Paketabhängigkeiten oder eine geeignete Bereitstellungsumgebung bereitgestellt werden. Reine CLI-/Serverpakete fügen diese GUI-Abhängigkeiten standardmäßig nicht hinzu; explizite Überschreibungen der Abhängigkeiten sind weiterhin möglich.

## Mehrere Programmdateien erstellen

Das Projekt definiert das primäre Programm. Setzen Sie `project.package` auf dessen Go-Paket (Standard: `.`) und `project.kind` auf `gui`, `cli` oder `server` (Standard: `gui`).
Zusätzliche Programme verwenden benannte Einträge unter `binaries`:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Jedes Programm erhält eigene Kompilierungs- und Cache-Eingaben. Importierte Projektpakete und eingebettete Dateien werden nachverfolgt; nach einer Änderung an einem unabhängigen Programm werden die anderen nicht neu erstellt. Zusätzliche Programmdateien werden auf Desktopplattformen unterstützt und in die standardmäßigen App-Bundles und Installationsprogramme aufgenommen. Paketersatzdateien im Besitz der Benutzer steuern ihren Inhalt selbst.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Ein benanntes Profil ist eine vollständige Anforderung; verwenden Sie diese Auswahlüberschreibungen bei einem anonymen Build. `--jobs` begrenzt die Anzahl der ausführenden Worker. Wenn Sie ausschließlich CLI-/Serverprogramme auswählen, entfällt die Frontend-Verarbeitung.

Konfigurieren Sie sowohl `project.version_file` als auch `project.version_symbol`, um eine Version zu lesen, ohne Anwendungscode auszuführen. Die Datei muss eine statische Go-Zeichenfolgendeklaration enthalten. Das Literal `project.version` kann dann entfallen. Setzen Sie `commit_symbol: main.Commit` für das Projekt oder eine benannte Binärdatei, um die Repository-Revision in die Zeichenfolgenvariable dieser Programmdatei einzufügen. Wenn die Commit-Injektion angefordert wurde, sind fehlende Repository-Metadaten ein Fehler.

`upx_level: 1` bis `9` komprimiert optional Linux- und Windows-Programmdateien für Releases; null deaktiviert die Komprimierung. UPX muss installiert sein. Entwicklungs- und macOS-Builds verzichten auf die Komprimierung. Bei einem Kompilierungs- oder Komprimierungsfehler bleibt das zuvor veröffentlichte Programm erhalten.

## Für die Verteilung paketieren und signieren

Für das Erstellen unsignierter Pakete gelten andere Host-/Werkzeuganforderungen als für das Signieren:

| Ausgabe | Hosts und Werkzeuge für die Erstellung | Signierung |
| --- | --- | --- |
| Windows-MSIX | Windows, macOS oder Linux; Go-Writer, kein Windows SDK erforderlich. | Windows-Host und Signierwerkzeuge. |
| Windows-NSIS-Installationsprogramm | Windows, macOS oder Linux; hostnatives `makensis` erforderlich. | Windows-Host und Signierwerkzeuge. |
| macOS-App und DMG | Windows, macOS oder Linux; Go-Zusammenstellung, Zusammenführung universeller Binärdateien und portabler DMG-Writer. | macOS mit Apple-Werkzeugen und Zugangsdaten. |
| Linux DEB, RPM, Arch und Void | Windows, macOS oder Linux; native Go-Paket-Writer. | DEB/RPM verwenden ihre externen Signierwerkzeuge; Arch verwendet den nachfolgend beschriebenen Go-Signierer. Die Void-Ausgabe ist unsigniert. |
| Linux-AppImage | Übereinstimmender Linux-Host und passende Zielarchitektur; externes linuxdeploy und GTK-Bereitstellungs-Plug-in. | Kein integrierter AppImage-Signierablauf. |

Windows-MSIX-Archive werden direkt in Go erstellt; zum Erstellen eines unsignierten Pakets sind weder das Windows SDK noch Docker erforderlich. Zum Signieren unter Windows sind weiterhin ein Windows-Host und dessen Signierwerkzeuge erforderlich. NSIS verwendet den externen NSIS-Compiler und akzeptiert literale Definitionen über `packages.nsis.defines`. Unter macOS stellt Wails ausschließlich dem untergeordneten NSIS-Compilerprozess eine UTF-8-Locale bereit, wenn dessen geerbte Locale nicht UTF-8 ist; die Umgebungen des Terminals und der Anwendung bleiben unverändert. Das standardmäßige NSIS-Installationsprogramm lässt den WebView2-Bootstrapper weg, wenn alle paketierten ausführbaren Dateien CLI- oder Serverprogramme sind. Bei einem primären GUI-Programm oder einem GUI-Begleitprogramm bleibt WebView2 enthalten. Benutzerdefinierte NSIS-Skripte behalten den Bootstrapper bei, damit sie mit ihrer eigenen Logik zur Installation der Laufzeitumgebung kompatibel bleiben.

Linux unterstützt AppImage-, DEB-, RPM-, Arch-Linux- und Void-Pakete. AppImages behalten die Grenze zu externem linuxdeploy und dem GTK-Plug-in bei. Setzen Sie `packages.appimage.update_info`, um eine zsync-Datei anzufordern, und `packages.appimage.checksum: true` für eine SHA256-Datei. Begleitdateien werden aus dem Paketcache wiederhergestellt und im Artefaktbeleg zusammen mit dem AppImage aufgeführt. Die AppImage-Paketierung erfordert derzeit Linux sowie eine übereinstimmende Host-/Zielarchitektur. Ihre GTK-Erkennung findet keine Bibliotheken, die von einer reinen Go-GUI dynamisch geladen werden; daher unterstützt dieser Adapter die Bereitstellung reiner Go-GUIs noch nicht. Verwenden Sie für diesen Buildmodus eines der Linux-Paketmanagerformate.

Arch-Pakete bewahren Vorabversionen in nativen Paketmetadaten auf; beispielsweise wird `1.2.3-beta.2` zu `1.2.3beta.2-1`, das vor der endgültigen Version einsortiert wird. Void-Abhängigkeiten akzeptieren reine Paketnamen wie `gtk4`, die zu `gtk4>=0` normalisiert werden; explizite Versionsvergleiche und vollständige `name-version_revision`-Werte bleiben erhalten. Verwenden Sie Void-Paketnamen statt Bibliotheksdateinamen.

Verwenden Sie zum Erstellen eines Arch-Pakets und einer separaten Signatur eine exportierte, ASCII-armierte PGP-Datei mit privatem Schlüssel und geben Sie optional den Namen einer Umgebungsvariable für die Passphrase an:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Legen Sie die benannte Umgebungsvariable in Ihrer lokalen Umgebung oder im Secret-Speicher der CI fest und führen Sie dann `wails3 build release` aus. Lassen Sie `credential` bei einem unverschlüsselten Schlüssel weg. Der Go-Signierer läuft auf jedem unterstützten Desktop-Host und benötigt kein GPG. Er veröffentlicht die ursprüngliche Datei `.pkg.tar.zst` und daneben eine Datei `.pkg.tar.zst.sig`; beide werden separat im Artefaktbeleg erfasst. Schlüsselinhalte werden weder bereitgestellt noch im Build-Cache gespeichert. Signaturen und ihre Veröffentlichung werden niemals aus dem Cache wiederverwendet. Schlägt das Signieren fehl, bleibt das zuvor veröffentlichte Dateipaar erhalten. Schlägt die Installation einer der beiden Dateien fehl, macht die Veröffentlichung beide Installationen rückgängig. Dadurch werden zwei Umbenennungen im Dateisystem für gleichzeitige Leser nicht atomar.

macOS-App-Bundles enthalten Begleitprogramme. Legen Sie `darwin.microphone_usage_description` fest, um den Verwendungstext und – sofern keine benutzerdefinierte Entitlements-Datei angegeben ist – das Entitlement für die Audioeingabe hinzuzufügen. Benutzerdefinierte plist- und Entitlements-Quelldateien bleiben unverändert. Signierung und Beglaubigung verwenden unter macOS die vorhandenen Apple-Werkzeuge und Anmeldedaten; für die Beglaubigung ist ein echtes Notarprofil erforderlich. Signierte Ausgaben behalten ihre plattformspezifische Erweiterung bei, beispielsweise `notes.signed.app` und `notes.signed.dmg`.

## Steuerung während der Entwicklung

Geben Sie während `wails3 dev` einen Befehl ein und drücken Sie anschließend die Eingabetaste:

| Befehl | Aktion |
| --- | --- |
| `r` | Erzwingt einen erneuten Build und Neustart, auch bei unverändertem Backend. |
| `b` | Generiert Bindings neu und führt einen erneuten Build aus. |
| `f` | Startet das Frontend neu. |
| `c` | Leert das Terminal. |
| `q` | Beendet die Sitzung und ihre Kindprozesse. |
| `h` | Zeigt die Steuerungsbefehle an. |

Bei fehlgeschlagenen Builds läuft die aktuelle Anwendung weiter. Die Steuerung erfordert ein interaktives Terminal; über eine Pipe weitergeleitete Eingaben ändern das Sitzungsverhalten nicht.

## Arbeit vor der Ausführung prüfen

Die Planprüfung löst denselben unveränderlichen Graphen auf, den auch die Ausführung verwendet. Die Textausgabe ist für Menschen bestimmt; JSON ist ein versionierter, maschinenlesbarer Vertrag.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Konfiguration ohne Build validieren

Verwenden Sie den dafür vorgesehenen Preflight-Befehl in der CI, als Editor-Task oder vor dem Committen der Konfiguration. Ohne Profil validiert er das BasismManifest und jedes benannte Profil; mit einem Profil prüft er nur diese vollständige Anforderung:

```bash
wails3 config check
wails3 config check release
```

Fehler bei Syntax, Schema, Semantik, Pfaden, Zielsystemen, Formaten, Signierung und Planung sowie durch die Konfiguration verursachte Fehler bei Host-Fähigkeiten verwenden denselben Diagnose-Renderer wie Build und Dev. Diagnosen geben das Feld und den Quellbereich an, zeigen die Quellzeile mit einem Einfügezeichen und enthalten, sofern hilfreich, einen Korrekturhinweis. Umgebungswerte werden in Quellauszügen der Diagnose unkenntlich gemacht. Host-Fehler, die ausschließlich durch Befehlszeilenoptionen verursacht werden, bleiben gewöhnliche Befehlsfehler und werden nicht fälschlich als ungültiges YAML gemeldet.

## Projekteigene Assets verwenden

Pfade im Manifest verweisen auf Dateien in Ihrem Besitz, etwa Symbole, Plattformmanifeste, Paketvorlagen und Maintainer-Skripte. Wails bearbeitet diese Quelldateien niemals. Bevor ein Werkzeug sie verarbeitet, kopiert Wails sie in einen temporären `.wails/`-Arbeitsbereich. Die endgültigen Artefakte werden unter `build.output` abgelegt, und ein verifizierter Beleg wird nach `.wails/artifacts/receipt.json` geschrieben.

## Auf Android bereitstellen

Android-Produktionsprofile erzeugen AABs. Für die lokale Bereitstellung erstellt Wails intern ein ABI-spezifisches Debug-APK, installiert es mit ADB und startet die App:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Ohne expliziten Selektor bietet ein interaktives Terminal verbundene Geräte und konfigurierte virtuelle Android-Geräte zur Auswahl an. Bei einer automatisierten Ausführung mit mehreren verfügbaren Zielsystemen muss `--device` oder `--emulator` angegeben werden.

Binden Sie einen auf die Anwendung beschränkten Log-Stream ein und überwachen Sie das Zielsystem, bis die App beendet wird oder Sie Strg+C drücken:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Eingebundene Ausführungen folgen einem Neustart der Anwendung und melden Verbindungsabbrüche oder Offline-Übergänge des Zielsystems. Emulatoren laufen standardmäßig weiter, um schnelle Iterationen zu ermöglichen. `--stop-emulator` muss explizit angegeben werden, erfordert `--logs` und beendet beim Ende der eingebundenen Sitzung nur ein Emulatorziel. Wenn Wails einen neuen Emulator startet, dessen Startbereitschaft jedoch nicht erreicht oder der Vorgang abgebrochen wird, bereinigt Wails diesen Prozess automatisch.

## Begrenzte Lebenszyklus-Hooks ausführen

Die experimentelle Hook-Schnittstelle ruft in jeder deklarierten Phase genau ein projekteigenes Skript auf. Ein Hook ist kein Shell-Befehl: Wails führt die Datei direkt aus, standardmäßig vom Projektstammverzeichnis aus, sofern `directory` nicht festgelegt ist.

Auf Unix-Hosts muss das Skript ausführbar sein und mit einem Shebang beginnen. Windows-Hosts akzeptieren `.cmd`-, `.bat`- und `.ps1`-Skripte über feste Interpreter-Adapter; Wails parst auf keiner der beiden Plattformen einen eingebetteten Shell-Befehl.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

Unterstützte Phasen sind `before_build`, `after_build`, `before_package`, `after_package`, `before_sign` und `after_sign`. Build-Hooks werden je nach Bedarf einmal für das gemeinsame Projekt oder einmal pro Zielsystem ausgeführt. Paketierungs- und Signierungs-Hooks bilden Barrieren um alle für ein Zielsystem angeforderten Formate, während die Formatoperationen innerhalb der Barriere weiterhin gleichzeitig ausgeführt werden.

Hooks erhalten genau eine von Wails verwaltete Umgebungsvariable: `WAILS_HOOK_CONTEXT_FILE`. Sie bezeichnet eine geschützte JSON-Datei mit dem versionierten Aufrufkontext:

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

Der Kontext enthält aufgelöste Pfade, kopiert jedoch niemals die geerbte Prozessumgebung oder Geheimnisse. Wails erstellt ihn atomar und mit Dateiberechtigungen ausschließlich für den Eigentümer unter `.wails/hooks/context/`. Nach erfolgreichem Abschluss oder Abbruch wird er entfernt. Bei einem Fehler bleibt er zur Diagnose erhalten, und sein Pfad wird in die Fehlermeldung aufgenommen. Ein Exitcode ungleich null lässt den regulären Plan-Schritt fehlschlagen; ein Abbruch beendet die Prozessgruppe des Hooks.

Hooks werden standardmäßig bei jedem Aufruf ausgeführt. Aktivieren Sie die Zwischenspeicherung nur, wenn für das Skript sämtliche Eingaben und Ausgaben vollständig deklariert sind:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

Das Skript und sein Ausführbarkeitsmodus sind implizite Eingaben. Mehrere Ausgaben müssen sich ein einziges Verzeichnis außerhalb des Stammverzeichnisses teilen; dieses bildet das begrenzte Artefakt, das Wails erfasst und wiederherstellt. Das Ausgabestammverzeichnis darf weder das Skript noch eine deklarierte Eingabe enthalten. Zwischengespeicherte Ausgabebäume dürfen ausschließlich reguläre Dateien und Verzeichnisse enthalten, niemals symbolische Links oder spezielle Dateien. Ein stets ausgeführter Hook verhindert die Zwischenspeicherung seiner nachgelagerten Operationen, da nicht deklarierte Nebeneffekte nicht sicher in einem Fingerabdruck erfasst werden können. Ein zwischenspeicherbarer Hook erhält die Zwischenspeicherung nachgelagerter Operationen über seine deklarierte Artefaktidentität aufrecht.

## Vorhandenes Projekt migrieren

Erstellen und prüfen Sie einen inaktiven Entwurf:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Konservative Lebenszyklus-Tasks ohne Argumente, die ein projekteigenes Skript aufrufen, werden in Hooks übersetzt. Anderes erreichbares benutzerdefiniertes Taskfile-Verhalten wird als Blocker gemeldet. Aktivieren Sie `wails.migrated.yaml` nach der Prüfung ausdrücklich:

```bash
wails3 migrate --activate
```

Bei der Aktivierung wird `wails.yaml` atomar erstellt. Taskfiles, `build/config.yml`, Skripte oder Assets werden weder gelöscht noch umbenannt oder bearbeitet. Ungültiges YAML schlägt als YAML fehl; es wird niemals auf den Legacy-Pfad zurückgegriffen.

## Alle Standardwerte materialisieren

Schreiben Sie einen inaktiven Snapshot, um die vollständig aufgelöste Konfiguration zu prüfen oder zu kopieren:

```bash
wails3 eject
```

Dadurch wird `wails.ejected.yaml` erstellt. Die aktive Konfiguration wird niemals überschrieben.

## Kompatibilitätszeitraum

`wails3 package` und `wails3 sign` sind veraltete Aliasse für denselben Planer; neue Automatisierungen sollten `wails3 build` mit einem Profil oder Formaten verwenden. Die Einführung ist an Releases gebunden: Das erste Release macht YAML zum Pfad für neue Projekte, das folgende Release warnt bei Legacy-Taskfile-Routing, und das nächste entfernt die Taskfile-Ausführung zusammen mit den Aliassen. Für jeden Übergang muss die Akzeptanzmatrix der nativen Hosts für das jeweilige Release vollständig erfolgreich sein.

Eine Beschreibung aller Felder, Typen, Standardwerte und Anwendbarkeitsregeln finden Sie in der generierten [`wails.yaml`Feldreferenz](/reference/wails-yaml-fields/).

Die [`wails3 run`Konfigurationsreferenz](/reference/wails-yaml-run/) beschreibt die Standardwerte für lokale Starts, plattformspezifische Überschreibungen und den `go run .`-Fallback ohne Manifest.
