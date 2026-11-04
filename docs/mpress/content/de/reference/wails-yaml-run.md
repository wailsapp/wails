---
title: "YAML-Ausführungskonfiguration"
description: "Lokale Standardeinstellungen, plattformspezifische Überschreibungen und Startverhalten für wails3 run."
---

## Verfügbarkeit

Diese Referenz beschreibt die experimentelle Ausführungserweiterung. Sie garantiert nicht, dass diese in einer veröffentlichten CLI verfügbar ist.

`wails3 run` erstellt und startet die Anwendung mit den lokalen Standardeinstellungen des Manifests. Die manifestgesteuerte Ausführung erfordert die ausdrückliche Aktivierung über [`WAILS_EXP_USE_WAKE=1`](/experimental/yaml-builds#enable-yaml-builds). Wenn kein Manifest gefunden wird, führt der Befehl im aktuellen Verzeichnis `go run .` aus, ohne diese Aktivierung zu erfordern. Ein ungültiges oder nicht lesbares Manifest ist ein Fehler; in diesem Fall wird nie auf die Ausweichlösung zurückgegriffen. Ist ein Manifest vorhanden, während YAML deaktiviert ist, erklärt der Befehl, wie YAML aktiviert wird.

Die generierte [`wails.yaml` Feldreferenz](/reference/wails-yaml-fields) führt Feldtypen und Standardwerte auf. Die Unterstützung nativer Mobilplattformen und von macOS-Bundles muss auf dem jeweiligen Host überprüft werden, bevor eine Version diese Plattformen als einsatzbereit ausweist.

## Befehlsvertrag

`wails3 run` sucht das nächstgelegene `wails.yaml`, löst die Konfiguration für ein Ziel auf, erstellt die Anwendung über die vorhandene zwischengespeicherte Pipeline und startet das Ergebnis. Die Suche beginnt im Arbeitsverzeichnis und wird nach oben fortgesetzt. Dabei wird vor dem Abbruch an einer `go.mod`-Grenze nach einem Manifest gesucht. Das Manifestverzeichnis ist das Projektstammverzeichnis und das Arbeitsverzeichnis der Desktopanwendung.

Standardziel sind Betriebssystem und Architektur des Hosts. Mobile Ziele müssen ausdrücklich ausgewählt werden. Eine Cross-Kompilierung bedeutet nicht, dass eine Desktopprogrammdatei auf dem Host ausgeführt werden kann; inkompatible Desktopziele schlagen bereits vor dem Erstellen fehl. Dies gilt auch für mit `CGO_ENABLED=0` erstellte Programmdateien: Eine erfolgreiche reine Go-Cross-Kompilierung stellt weder einen Emulator bereit noch ändert sie den lokalen Startvertrag. Verwenden Sie [`wails3 build` für die Cross-Kompilierung](/experimental/yaml-builds/#cross-compile-desktop-applications) und führen Sie die resultierenden Dateien auf ihrem Zielsystem aus oder installieren Sie sie dort. Verwenden Sie für Android `--target android/arm64 --device <adb-serial>` oder `--target android/arm64 --emulator <avd-name>`. Die ausgewählte Architektur muss zum Gerät oder Emulator passen. Verwenden Sie für einen iOS-Simulator `--target ios/arm64 --device <udid>` und fügen Sie für ein physisches Gerät `--destination device` hinzu. Planung und Ausführung für iOS erfordern macOS und Xcode.

Mit YAML verwendet der Build eingebettete Assets und den Produktions-Buildpfad, der das Tag `production` automatisch bereitstellt. `wails3 dev` bleibt für Live Reload zuständig; `dev.tags` trägt nicht zu `wails3 run` bei.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan` löst das ausgewählte Ziel, die effektiven Build-Tags, das ausführbare oder installierbare Artefakt, die Argumente und das Arbeitsverzeichnis auf und zeigt sie an, ohne etwas zu erstellen, zu installieren oder zu starten. Der Plan führt die Schlüsselnamen der Umgebungsüberschreibungen auf und lässt deren Werte aus.

## CLI-Optionen

| Option | Bedeutung |
| --- | --- |
| `--tags <tags>` | Zusätzliche, durch Kommas getrennte Go-Build-Tags. |
| `--skip-build` | Vorhandene Desktop-Executables aus dem YAML-Ausgabeverzeichnis importieren und anschließend ohne Kompilierung starten. |
| `--target <os>/<arch>` | Wählt ein Ziel aus; standardmäßig werden Betriebssystem und Architektur des Hosts verwendet. |
| `--plan` | Zeigt die aufgelösten Einstellungen an, ohne etwas zu erstellen oder zu starten. Plattformübergreifende Desktoppläne sind zulässig. |
| `--device <id>` | Android-adb-Seriennummer oder Kennung eines iOS-Simulators beziehungsweise -Geräts. |
| `--emulator <name>` | Zu startendes oder wiederzuverwendendes Android Virtual Device; kann nicht mit `--device` kombiniert werden. |
| `--destination <kind>` | iOS-Ziel: `simulator` (Standard) oder `device`. |
| `-- <args...>` | Ersetzt die konfigurierten Anwendungsargumente. Ein allein stehendes `--` löscht sie. |
| `--appargs <string>` | Kompatible Schreibweise: Analysiert eine Zeichenfolge mit POSIX-artiger Quotierung. Kann weder wiederholt noch mit Argumenten nach `--` kombiniert werden. |

Selektoren für Mobilgeräte erfordern ein Manifest und ein passendes mobiles Ziel.

`--skip-build` erfordert ein aktives YAML-Projekt und ein kompatibles Desktop-Executable. Dabei werden Frontend-Verarbeitung, Bindings, Kompilierung, UPX und Build-Hooks übersprungen; konfigurierte Startargumente und Umgebungsvariablen gelten weiterhin. Importierte Executables werden validiert und in den von Wake generierten Arbeitsbereich kopiert, ohne ihre Quelldateien zu ändern. macOS-GUI-Anwendungen werden aus diesen Eingaben und den aktuellen Paket-Assets zusammengestellt. Siehe [Pfade und Anforderungen für vorkompilierte Eingaben](/experimental/yaml-builds/#use-prebuilt-desktop-executables).
Build-Tags können ein vorkompiliertes Executable nicht ändern. Daher haben `--tags` und konfigurierte Build-Tags in diesem Modus keine Wirkung.

## Felder

Die optionale `run`-Zuordnung stellt gemeinsame Startstandards bereit. Eine optionale `run`-Zuordnung innerhalb von `targets["<os>"]` oder `targets["<os>/<arch>"]` stellt plattform- oder architekturspezifische Standardwerte bereit. Architekturspezifische Einstellungen haben unabhängig von der Reihenfolge der Zuordnungen Vorrang vor betriebssystemweiten Einstellungen.

| Feld | Typ | Standardwert | Bedeutung |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | Zusätzliche Go-Tags, die nur für Ausführungs-Builds verwendet werden. |
| `run.args` | `list(string)` | `[]` | An die Anwendung übergebene Argumente. |
| `run.environment` | `map(string)` | `{}` | Überschreibungen der Laufzeitumgebung, getrennt von der Compilerumgebung. |
| `target["<os>"].run.tags` | `list(string)` | `[]` | Zusätzliche Ausführungs-Tags für die ausgewählte Plattform. |
| `target["<os>"].run.args` | `list(string)` | geerbt | Ersetzt bei ausdrücklicher Festlegung die gemeinsamen Argumente; `[]` löscht sie. |
| `target["<os>"].run.environment` | `map(string)` | `{}` | Überschreibt gemeinsame Werte der Laufzeitumgebung nach Schlüssel. |

Argumente werden ohne Shell-Erweiterung als Argumentvektor übergeben. Umgebungswerte sind gemäß dem bestehenden YAML-Vertrag, der ausschließlich Literale zulässt, literale Zeichenfolgen.

## Plattformbeispiel

Dieses Fragment ergänzt die normale Projektkonfiguration. `example_private_api` ist ein Platzhalter für das tatsächliche Tag der privaten API, das für dieses Beispiel erforderlich ist.

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Die Linux-Umgebungseinstellung gilt speziell für dieses Beispiel und erfordert eine Wayland-Sitzung; sie ist kein Wails-Standard. Das Tag der privaten API wird nur unter macOS ausgewählt. Das vorhandene `target.tags` gilt für Produktions-Builds einschließlich Ausführungs-Builds. `target.run.tags` beschränkt ein zusätzliches Tag auf Ausführungs-Builds.

Eine architekturspezifische Überschreibung kann die betriebssystemweiten Standardwerte präzisieren:

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

Zusammen mit der vorangehenden Linux-Zuordnung behält dies `GDK_BACKEND` bei, fügt das architekturspezifische Tag hinzu und löscht die gemeinsamen Anwendungsargumente.

## Auflösungsregeln

| Einstellung | Auflösung von der niedrigsten zur höchsten Priorität |
| --- | --- |
| Tags | Kombiniert `build.tags`, das ausgewählte `target.tags`, das automatische `production`, `run.tags`, das ausgewählte `target.run.tags` und das CLI-`--tags`; behält jeweils das erste Vorkommen eines Tags bei. |
| Laufzeitumgebung | Übernimmt die startende Shell, überlagert sie mit `run.environment` und anschließend mit dem ausgewählten `target.run.environment`. Ein leerer Wert setzt eine leere Zeichenfolge. |
| Anwendungsargumente | Das gemeinsame `run.args` wird durch ausdrücklich festgelegte Zielargumente und anschließend durch explizite Argumente nach `--` ersetzt. Ein allein stehendes `--` löscht konfigurierte Argumente. |

Innerhalb des ausgewählten Ziels werden betriebssystemweite Einstellungen vor architekturspezifischen Einstellungen aufgelöst: Tags werden gesammelt, Umgebungswerte nach Schlüssel zusammengeführt und ausdrücklich angegebene Argumente ersetzen die weniger spezifische Liste. Ausführungs-Builds stellen für das entsprechende mobile Ziel außerdem `android` oder `ios` bereit.

Die Bindungsgenerierung und die Go-Kompilierung erhalten dieselben aufgelösten Tags. Die Build-Cache-Schlüssel enthalten diese Tags. Werden nur Anwendungsargumente oder die Laufzeitumgebung geändert, wird die Kompilierung nicht ungültig; beim nächsten Start werden die neuen Werte verwendet.

Plattformzuordnungen passen eine Plattform an; ihr Vorhandensein macht andere Plattformen nicht zu nicht unterstützten Plattformen. `project.supported_platforms` beschränkt die von Builds und `run` akzeptierten Plattformen; eine ausgelassene oder leere Liste lässt alle Plattformen zu. Beispielsweise bewirkt `supported_platforms: [darwin]` innerhalb von `project`, dass ein ausschließlich für macOS bestimmtes Beispiel unter Linux frühzeitig fehlschlägt.

## Plattformspezifisches Startverhalten

| Plattform | Startverhalten |
| --- | --- |
| Linux und Windows | Startet die ausführbare Datei mit verbundenen Standardein- und -ausgaben, Behandlung der Prozessbeendigung und Weitergabe des Exit-Status der Anwendung. |
| macOS | Erstellt ein App-Bundle und führt dessen Binärdatei `Contents/MacOS` mit verbundenen Standardein- und -ausgaben aus. Registriert das Bundle für URL-Schemata und Dateizuordnungen bei Launch Services. Signierung und Akzeptanz des nativen Lebenszyklus erfordern eine Überprüfung unter macOS. |
| Android | Erstellt ein APK, wählt ein Gerät oder einen Emulator aus und installiert und startet die Anwendung. Meldet Fehler beim Erstellen, Installieren und Starten getrennt. |
| iOS | Erstellt ein App-Bundle, wählt einen Simulator oder ein Gerät aus, erfüllt die Signierungs- und Bereitstellungsanforderungen und installiert und startet die Anwendung. Builds für physische Geräte erfordern die entsprechenden Anmeldedaten und die Apple-Toolchain. |

Android lehnt Anwendungsargumente und Überschreibungen der Laufzeitumgebung vor der Bereitstellung ab. Wenn auch Android unterstützt wird, gehören gemeinsame, nur für Desktop-Systeme vorgesehene Einstellungen in die Zielzuordnungen für Desktop-Systeme. Bei Android-Ausführungen werden Anwendungsprotokolle angehängt und die Anwendung wird beendet, sobald die Verbindung endet. iOS übergibt Argumente an das ausgewählte Startprogramm. Simulatoren erhalten Umgebungsüberschreibungen über `SIMCTL_CHILD_`-Variablen; physische Geräte erhalten über `devicectl --environment-variables` eine JSON-Zuordnung. Der mobile Exit-Status beschreibt das Ergebnis des Startprogramms oder der Protokollüberwachung und nicht einen portablen Exit-Code der Anwendung.

## Reine Go-Beispiele und Fallback

Ein reines Go-Beispiel mit bereits vorhandenen oder eingebetteten Assets deaktiviert die Frontend-Phasen:

```yaml
frontend:
  disabled: true
```

Dadurch werden die Installation der Frontend-Abhängigkeiten, die Generierung der Bindings und der Frontend-Build übersprungen. Beispiele, die generierte Binding-Klassen benötigen, können stattdessen `interfaces: false` unter `frontend.bindings` festlegen.

Der Fallback ohne Manifest führt `go run [-tags <tags>] . [application arguments]` aus, ohne `production` oder `devtools` hinzuzufügen. Er folgt dem Exit-Status-Verhalten von Go. Argumente nach `--` werden an die Anwendung übergeben, einschließlich Flags wie `--help`. Die von `--plan` angezeigten Laufzeitwerte sind auf die Namen der Umgebungsschlüssel beschränkt; die Werte werden nicht ausgegeben. Der Fallback-Plan zeigt den Go-Befehl und das Arbeitsverzeichnis an.

## Anwendungsargumente für die Entwicklung

Bei einem aktiven YAML-Projekt akzeptiert `wails3 dev` Anwendungsargumente unabhängig von den standardmäßigen Produktionswerten in `run`. `dev.args` stellt gemeinsame Standardwerte bereit; eine verschachtelte `dev`-Zuordnung in `targets["<os>"]` oder `targets["<os>/<arch>"]` ersetzt sie für dieses Ziel. Architektureinstellungen haben unabhängig von der Reihenfolge in der Zuordnung Vorrang vor betriebssystemweiten Einstellungen. Ein explizites `args: []` löscht geerbte Argumente.

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args` und `run.tags` wirken sich nicht auf Entwicklungsstarts aus. Die bestehende Build-Richtlinie `dev.tags` bleibt unverändert. Relative Anwendungspfade werden von der Anwendung ausgehend vom Manifestverzeichnis interpretiert.

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

Argumente nach `--` ersetzen die konfigurierten Argumente, auch wenn die Liste leer ist. `--help` nach diesem Trennzeichen gehört zur Anwendung. `--appargs=""` löscht ebenfalls die Standardwerte. Die Kompatibilitätsoption akzeptiert eine einzelne Zeichenfolge in Anführungszeichen; ihre Verwendung zusammen mit `--`, positionsgebundenen Anwendungsargumenten oder einem zweiten `--appargs` ist ein Fehler. Fehlerhafte Anführungszeichen werden vor dem Start einer Sitzung zurückgewiesen.

Der Kompatibilitätsparser verwendet auf jedem Host dieselben POSIX-artigen Regeln für Anführungszeichen. Er führt keine Shell-Befehle aus und expandiert keine Variablen; die aufrufende Shell hat Text möglicherweise bereits expandiert, bevor sie ihn an Wails übergibt. Argumente nach `--` liegen bereits als separate Werte vor und werden nicht weiter tokenisiert. Dies ist die bevorzugte Form für Windows-Pfade, beispielsweise in PowerShell:

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

Die Sitzung behält den aufgelösten Argumentvektor bei erneuten Builds des Quellcodes bei. Wenn nur `dev.args` oder die Argumente des ausgewählten Ziels geändert werden, wird das Backend mit der zwischengespeicherten ausführbaren Datei neu gestartet; das Frontend wird nicht neu gestartet. CLI-Überschreibungen bleiben nach dem Neuladen des Manifests wirksam. Bei fehlgeschlagenen Ersetzungen wird die vorherige ausführbare Datei mit ihren vorherigen Argumenten beibehalten oder wiederhergestellt. Anwendungsargumente sind keine Compiler-Flags und fließen nicht in die Schlüssel des Kompilierungs-Caches ein. `dev --plan` zeigt die effektiven Argumente zusammen mit dem endlichen Build-Plan an, ohne die Anwendung zu starten.

Linux, Windows und macOS übergeben den Vektor direkt an die ausführbare Datei. iOS übergibt ihn über das Startprogramm des Simulators oder Geräts. Android lehnt nicht leere Anwendungsargumente für die Entwicklung vor der Gerätesuche oder dem Build ab, da kein entsprechender unterstützter Argumentkanal vorhanden ist. Nur für Desktop-Systeme vorgesehene Standardwerte können in den Zielzuordnungen für Desktop-Systeme festgelegt oder mit `args: []` unter `targets.android.dev` gelöscht werden.

Diese Entwicklungsoptionen erfordern ein aktives `wails.yaml` und die explizite Aktivierung über YAML. Sie schreiben im bisherigen Entwicklungspfad weder Taskfile- noch Watcher-Shell-Befehle um.
