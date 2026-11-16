---
title: "YAML-Builds"
description: "Konfigurieren und verwenden Sie das experimentelle native Build-System von Wails v3 mit einer einzigen wails.yaml-Datei."
sidebar:
  order: 2
  badge:
    text: "Experimentell"
    variant: caution
---

@note{type="caution" title="Experimentelle Vorschau"}
Die erste akzeptierte Version unterstützt ausschließlich literale Konfiguration. Dieses Experiment dient außerdem als Prototyp für zeitlich begrenzte Lebenszyklus-Hooks, damit deren Gestaltung vor einem Vorschlag für eine Veröffentlichung bewertet werden kann. Benutzerdefinierte Phasen, YAML-Ausdrücke und typisierte Wails-Toolaufrufe bleiben vorerst zurückgestellt. Bestehende Projekte verwenden weiterhin Taskfiles, bis Sie eine geprüfte Migration ausdrücklich aktivieren.
@end

Neu bei Wake? Beginnen Sie mit dem [Tutorial für den Einstieg](/experimental/wake/).
Es behandelt die Grenzen des Experiments, sämtliche Bildschirme zur Projekt- und Build-Konfiguration, einen ersten lokalen Build und eine geprüfte Migration. Diese Seite beschreibt die einzelnen Workflows für Build, Paketierung und Konfiguration ausführlicher.

## YAML-Builds aktivieren

Setzen Sie `WAILS_EXP_USE_WAKE` in jedem Terminal, in dem Sie die experimentellen Befehle ausführen. Für Bash oder Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Für PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Die Variable aktiviert die Funktion allein durch ihre Existenz: Selbst ein leerer Wert oder `false` aktiviert sie. Entfernen Sie die Variable (mit `unset WAILS_EXP_USE_WAKE` in Bash/Zsh oder `Remove-Item Env:WAILS_EXP_USE_WAKE` in PowerShell), um das Experiment zu deaktivieren. Sobald `wails.yaml` aktiv ist, müssen Befehle Wake erneut aktivieren und greifen nie auf Taskfiles zurück. Ohne die Variable sind das YAML-Build-/Dev-Routing und die YAML-Projekterzeugung deaktiviert; außerdem sind die Befehle `migrate`, `config` und `eject` nicht verfügbar. Noch nicht migrierte Projekte verwenden weiterhin die separate eingebettete Task-Laufzeit.

## Ein Projekt konfigurieren

Um ein Projekt im Browser zu erstellen und zu konfigurieren, führen Sie bei aktivierter Opt-in-Option `wails3 init -ui` aus und folgen Sie dem [Wake-Tutorial](/experimental/wake/). Die Projekteinstellungen befinden sich in `wails.yaml`. Bearbeiten Sie diese Datei, um ein vorhandenes Projekt zu ändern. Mit `wails3 setup` konfigurieren Sie gemeinsame Einstellungen, Werkzeuge und Anmeldedaten für die Signierung.

Fügen Sie bei aktivierter Opt-in-Option im Projektstamm `wails.yaml` hinzu. Das Vorhandensein dieser Datei wählt das native Build-System aus und bewirkt, dass Wails sämtliche Taskfiles in diesem Projekt ignoriert. Lassen Sie die Opt-in-Option für alle folgenden Befehle aktiviert. Die Schema-Modeline ermöglicht in Editoren mit Unterstützung für YAML-Sprachserver die Vervollständigung und Inline-Validierung.

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

Wählen Sie für einen einmaligen Build durch Kommas getrennte Ziele und Formate aus, ohne ein Profil hinzuzufügen:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Desktopanwendungen kreuzkompilieren

Setzen Sie für das experimentelle reine Go-Desktop-Backend `CGO_ENABLED: "0"` in der Build-Umgebung. Wake wählt für macOS- und Linux-GUI-Programme das Tag `purego` aus und verwendet Go direkt auf Windows-, macOS- und Linux-Hosts. Das Erstellen von Desktop-Programmdateien erfordert in diesem Modus weder Docker noch Emulation, einen C-Compiler oder ein SDK der Zielplattform. Auch die Anwendungsabhängigkeiten müssen `CGO_ENABLED=0` unterstützen.

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

Der zweite Befehl kombiniert beide macOS-Architekturen und erstellt mit Go ein unsigniertes DMG. Er kann auf jedem unterstützten Desktop-Host ausgeführt werden.

Der aufgelöste Modus wird für jede Programmdatei unabhängig ausgewählt. Ein deklariertes `build.environment.CGO_ENABLED`, eine Zielumgebung oder eine Binärdateiumgebung überschreibt den Wert des aufrufenden Terminals. Andernfalls verwenden ein ausdrücklich ausgewähltes `purego`-Tag, ein CLI-/Serverprogramm oder ein Windows-Ziel standardmäßig `0`; macOS- und Linux-GUI-Programme verwenden standardmäßig `1`. `--plan --json` meldet für jede Programmdatei den aufgelösten Compiler, den CGo-Modus und die Tags. Die Kreuzkompilierung erstellt Zieldateien; zum Ausführen eines inkompatiblen Desktop-Ziels ist weiterhin der entsprechende Host erforderlich.

Mit `CGO_ENABLED=1` verwendet Wake zuerst eine kompatible native Toolchain oder den konfigurierten `CC`-Cross-Compiler, danach – sofern verfügbar – Zig für einen Windows-Cross-Build und anschließend, soweit unterstützt, den bestehenden Docker-/Podman-Fallback. Ein ausdrücklich angegebenes `CC` muss den Zielcompiler, die Bibliotheken und das Sysroot bereitstellen; Zig allein stellt kein Linux-GTK-/WebKit-Sysroot bereit. Für Paketierung und Signierung können wie nachfolgend beschrieben zusätzliche Werkzeuge erforderlich sein. Mobile Ziele benötigen CGo und die SDKs ihrer Plattformen.

Eine reine Go-Kompilierung bündelt die Systembibliotheken des Zielsystems nicht. Linux-GUI-Anwendungen einschließlich begleitender GUI-Programmdateien benötigen weiterhin die durch ihre Build-Tags ausgewählte GTK-/WebKitGTK-Laufzeit. Paketabhängigkeiten oder eine geeignete Bereitstellungsumgebung müssen sie bereitstellen. Reine CLI-/Serverpakete fügen diese GUI-Abhängigkeiten standardmäßig nicht hinzu; explizite Überschreibungen der Abhängigkeiten bleiben verfügbar.

## Mehrere Programmdateien erstellen

Das Projekt definiert das primäre Programm. Setzen Sie `project.package` auf dessen Go-Paket (Standard: `.`) und `project.kind` auf `gui`, `cli` oder `server` (Standard: `gui`). Zusätzliche Programme verwenden benannte Einträge unter `binaries`:

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

Jedes Programm erhält eigene Kompilierungs- und Cache-Eingaben. Importierte Projektpakete und eingebettete Dateien werden nachverfolgt; eine Änderung an einem unabhängigen Programm führt nicht dazu, dass die anderen neu erstellt werden. Zusätzliche Programmdateien werden auf Desktopplattformen unterstützt und in die standardmäßigen App-Bundles und Installationsprogramme aufgenommen. Vom Benutzer bereitgestellte Ersatzpaketdateien bestimmen ihren Inhalt selbst.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Ein benanntes Profil ist eine vollständige Anforderung; verwenden Sie diese Auswahlüberschreibungen bei einem anonymen Build. `--jobs` begrenzt die Zahl der Ausführungs-Worker. Wenn Sie ausschließlich CLI-/Serverprogramme auswählen, entfallen die Frontend-Arbeiten.

Um eine Version zu lesen, ohne Anwendungscode auszuführen, konfigurieren Sie sowohl `project.version_file` als auch `project.version_symbol`. Die Datei muss eine statische Go-Stringdeklaration enthalten. Das Literal `project.version` kann dann entfallen. Setzen Sie `commit_symbol: main.Commit` für das Projekt oder eine benannte Binärdatei, um die Repository-Revision in die Stringvariable dieser Programmdatei einzufügen. Wenn die Commit-Einfügung angefordert wurde, sind fehlende Repository-Metadaten ein Fehler.

`upx_level: 1` bis `9` komprimiert optional Linux- und Windows-Programmdateien für Releases; null deaktiviert die Komprimierung. UPX muss installiert sein. Entwicklungs- und macOS-Builds verwenden keine Komprimierung. Bei einem Kompilierungs- oder Komprimierungsfehler bleibt das zuvor veröffentlichte Programm erhalten.

## Für die Verteilung paketieren und signieren

### Vorgefertigte Desktop-Programmdateien verwenden

Verwenden Sie `--skip-build`, wenn ein anderes Build-System oder ein CI-Auftrag die Programmdateien bereits erstellt hat:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Lege jede ausführbare Datei in `build.output` (standardmäßig `bin`) ab. Verwende dabei `project.binary_name` für die ausführbare Hauptdatei und für Begleitdateien den Namen des jeweiligen `binaries`-Eintrags. Unter Windows enden die Namen auf `.exe`. Bei einem Aufruf für mehrere Zielplattformen werden die Dateien jeder Zielplattform aus einem `<os>-<arch>`-Unterverzeichnis gelesen. Für eine macOS-GUI liest Wake die unverarbeitete ausführbare Datei aus `bin/<name>`, sofern sie dort vorhanden ist, andernfalls aus `bin/<name>.app/Contents/MacOS/<name>`. Eine `darwin/universal`-Eingabe muss ausführbare Slices sowohl für amd64 als auch für arm64 enthalten.

Wake validiert Dateityp und Zielarchitektur, kopiert die Programmdateien in seinen generierten Arbeitsbereich und verfolgt ihren Inhalt für das Paket-Caching. Fehlende, inkompatible oder als symbolische Links bereitgestellte Programmdateien führen zu einem Fehler, statt neu erstellt oder aus einer alten Eingabe wiederhergestellt zu werden. Bereitgestellte Programmdateien werden weder überschrieben noch mit chmod geändert. Änderungen an einer beliebigen Programmdatei invalidieren die Pakete, die sie verwenden.

Dieser Modus überspringt Frontend-Installation und -Builds, Bindings, Go-Kompilierung, UPX sowie die Hooks `before_build`/`after_build`. Er generiert weiterhin Paketressourcen, stellt macOS-App-Bundles zusammen und führt konfigurierte Paketierungs-/Signierungs-Hooks aus. Für Paketierungs- und Signierungswerkzeuge gelten weiterhin die Anforderungen der nachfolgenden Tabelle; Go- und Frontend-Build-Werkzeuge sind nicht erforderlich. Die Paketmetadaten stammen aus dem aktuellen Manifest, gegebenenfalls einschließlich `version_file`; aus der Binärdatei wird keine Version abgeleitet. Verwenden Sie übereinstimmende Build- und Paketeinstellungen, einschließlich der Art GUI oder CLI.

`wails3 package --skip-build` wird ebenfalls unterstützt. Mit `wails3 run` gelten Startargumente und Umgebung weiterhin, und die Programmdatei muss auf dem Host ausführbar sein. Dieser Modus erfordert ein aktives YAML-Projekt und unterstützt ausschließlich produktive Desktop-Ziele; er gilt weder für Entwicklungssitzungen noch für die Paketierung mobiler Anwendungen.

Für die Erstellung unsignierter Pakete gelten andere Host-/Werkzeuganforderungen als für die Signierung:

IPA-Payloads und ZIP-Dateien für die macOS-Notarisierung werden in Go geschrieben, wobei Bundle-Pfade, Ausführungsberechtigungen und relative symbolische Links erhalten bleiben. Für keinen der beiden Vorgänge ist eine externe ausführbare Datei `zip` oder `ditto` erforderlich. Diese Archive enthalten die Anwendungsdateien statt Finder-Metadaten oder Resource Forks. Die macOS-Signierung entfernt diese Metadaten bereits, da sie signierte Bundles ungültig machen.

| Ausgabe | Hosts und Werkzeuge für die Erstellung | Signierung |
| --- | --- | --- |
| Windows MSIX | Windows, macOS oder Linux; Go-Writer, kein Windows SDK erforderlich. | Go-Signierung mit einer exportierten RSA-PFX-Datei auf jedem Desktop-Host; Identitäten aus Zertifikatspeichern oder von Anbietern erfordern Windows SignTool. |
| Windows-NSIS-Installationsprogramm | Windows, macOS oder Linux; hostnatives `makensis` erforderlich. | Go-PE-Signierung mit einer exportierten RSA-PFX-Datei auf jedem Desktop-Host; Identitäten aus Zertifikatspeichern oder von Anbietern erfordern Windows SignTool. |
| macOS-App und DMG | Windows, macOS oder Linux; Zusammenstellung mit Go, Zusammenführung universeller Binärdateien und portabler DMG-Writer. | Go-Signierung mit einer exportierten RSA-PKCS#12-Identität oder Ad-hoc-Signierung; Schlüsselbundidentitäten verwenden weiterhin macOS-Werkzeuge. Die Notarisierung mit API-Schlüssel und das Anheften von Tickets erfolgen in Go. |
| Linux DEB, RPM, Arch und Void | Windows, macOS oder Linux; native Go-Paket-Writer. | Exportierte OpenPGP-Schlüssel verwenden den nativen Signierer für DEB/RPM/Arch. Ältere DEB/RPM-Schlüssel-IDs erfordern externe Werkzeuge. Void-Ausgaben sind unsigniert. |
| Linux AppImage | Windows, macOS oder Linux; Go-Bereitstellung und SquashFS-Writer sowie eine AppImage-Runtime für die Zielplattform. Das plattformübergreifende Paketieren von GUI-Anwendungen erfordert einen bereitgestellten Linux-Runtime-Verzeichnisbaum. | Kein integrierter Signierablauf für AppImages. |

Für signierte DMGs signiert Wake normalerweise eine private Kopie der Anwendung, bevor der Go-Writer das Image komprimiert, und signiert anschließend das Image. Dadurch muss das DMG nicht mit `hdiutil` eingebunden, erweitert und erneut komprimiert werden. Die Signierung erfolgt bei jeder Anforderung; das signierte Paket wird nicht aus dem Build-Cache wiederverwendet. Projekte mit `after_package`- oder `before_sign`-Hooks behalten den ursprünglichen Ablauf „Paketieren, dann signieren“ einschließlich `hdiutil` bei, sodass Änderungen dieser Hooks am unsignierten Image erhalten bleiben. Dieser mit Hooks kompatible Pfad zur Image-Bearbeitung erfordert macOS und `hdiutil`.

Legen Sie für die portable Apple-Signierung `darwin.signing.identity: "-"` für eine Ad-hoc-Signatur fest oder verwenden Sie `darwin.signing.certificate` für eine exportierte RSA-PKCS#12-Identität. `darwin.signing.credential` gibt den Namen der Umgebungsvariable für das Passwort an. Zertifikatsignaturen verwenden die gemeinsame Go-Implementierung für PFX/CMS; aus einem Schlüsselbund wird nicht automatisch ein privater Schlüssel exportiert. Schlanke und universelle ausführbare Dateien, verschachtelte App-Bundles, Frameworks und DMGs haben mit generierten Testidentitäten die strengen nativen Signaturprüfungen von Apple einschließlich der Ablehnung manipulierter Dateien bestanden.

Für die Notarisierung benötigt eine exportierte Developer-ID-Identität außerdem einen kompatiblen RFC-3161-`darwin.signing.timestamp_server`. Der gemeinsame Prüfer akzeptiert die unterstützten SHA-2-Zeitstempelprofile und lehnt nicht unterstützte ältere Algorithmen ab. Eine Ad-hoc-Signatur kann nicht notarisiert werden. Die native Validierung des Zertifikatsformats ist von der tatsächlichen Vertrauenswürdigkeit der Developer ID und der Annahme einer Einreichung durch Apple getrennt.

Der Go-Stapler authentifiziert Apples Ticketsignatur und prüft, ob sie alle Architekturen abdeckt, bevor er das Ticket an eine private Kopie anheftet. Er verwendet das beobachtete Protokoll zur Ticketbereitstellung, das Apple nicht als unterstützte API veröffentlicht. Er ersetzt nicht die aktuellen Widerrufsprüfungen von Gatekeeper. App- und signierte DMG-Ausgaben haben native `stapler validate`-Prüfungen und strenge Signaturprüfungen bestanden.

Windows-MSIX-Archive werden direkt in Go erstellt; zum Erstellen eines unsignierten Pakets sind weder das Windows SDK noch Docker erforderlich. Exportierte RSA-PFX-Identitäten signieren Windows-Programmdateien und einzelne MSIX-Pakete auf jedem Desktop-Host direkt in Go. Setzen Sie `windows.signing.certificate` auf die projektrelative PFX-Datei und `windows.signing.credential` auf die Umgebungsvariable für deren Passwort. AES/PBKDF2- und ältere 3DES-Exporte werden unterstützt; RC2-Exporte müssen erneut exportiert werden oder erfordern als Rückfalloption das native Windows SignTool. Identitäten aus Zertifikatspeichern oder von Anbietern sowie das Signieren von MSI-Paketen erfordern Windows SignTool. Der standardmäßige Zeitstempeldienst ist DigiCert; Zeitstempel werden vor der Veröffentlichung signierter Ausgaben überprüft. Formatbeschränkungen finden Sie im [Leitfaden zum Signieren](/guides/build/signing/).
NSIS verwendet den externen NSIS-Compiler und akzeptiert literale Definitionen über `packages.nsis.defines`. Unter macOS stellt Wails ausschließlich dem untergeordneten NSIS-Compilerprozess eine UTF-8-Locale bereit, wenn dessen geerbte Locale nicht UTF-8 ist; Terminal- und Anwendungsumgebungen bleiben unverändert. Das standardmäßige NSIS-Installationsprogramm lässt den WebView2-Bootstrapper weg, wenn jede paketierte Programmdatei eine CLI oder ein Server ist. Eine primäre GUI-Anwendung oder ein GUI-Begleitprogramm behält WebView2 bei. Benutzerdefinierte NSIS-Skripte behalten den Bootstrapper bei, damit sie mit ihrer eigenen Logik zur Laufzeitinstallation kompatibel bleiben.

Linux unterstützt AppImage-, DEB-, RPM-, Arch-Linux- und Void-Pakete. Wake erstellt AppImages einschließlich GTK/WebKit-Bereitstellung, SquashFS-Komprimierung und zsync-Generierung direkt in Go. Dabei werden weder linuxdeploy noch das GTK-Plug-in, appimagetool oder mksquashfs aufgerufen. Für die Zielarchitektur wird eine festgelegte und per Prüfsumme verifizierte AppImage-Typ-2-Runtime heruntergeladen; mit `packages.appimage.runtime_path` wählen Sie stattdessen eine projektbezogene Runtime-ELF-Datei aus.

Unter Linux dient standardmäßig das Root-Dateisystem des Hosts als Bibliotheksquelle. Um eine GUI für eine andere Architektur oder unter Windows/macOS zu paketieren, legen Sie `packages.appimage.runtime_root` auf einen projektbezogenen Verzeichnisbaum des Linux-Ziel-Root-Dateisystems fest, der dessen Bibliotheken, kompilierte Schemas und Modul-Caches enthält. Gemeint ist das gesamte Root-Layout einschließlich `usr/`, nicht das Verzeichnis `usr/` selbst. Stellen Sie eine Runtime aus der ältesten Linux-Distribution bereit, die Sie unterstützen möchten. Wake liest ELF-Importe und explizite GTK/WebKit-Profile, ohne Zielprogramme auszuführen. Es bündelt die ausgewählte Abhängigkeitshülle, jedoch weder glibc noch den dynamischen Loader. Inhalts-Hashes von Bibliotheken und Ressourcen invalidieren nach Runtime-Aktualisierungen den Paket-Cache.

Reine Go-GUI-Builds werden unterstützt. Wake wählt GTK 3 für mit `gtk3` erstellte ausführbare GUI-Dateien und andernfalls GTK 4 aus; dies gilt auch für GUI-Begleitanwendungen einer CLI- oder Serveranwendung. Alle ausführbaren GUI-Dateien in einem Image müssen dieselbe GTK-Version verwenden. Für vorgefertigte Binärdateien kann `packages.appimage.gtk_version` explizit `3` oder `4` auswählen; importierte CGo-Bibliotheken müssen mit dieser Auswahl übereinstimmen.

Legen Sie `packages.appimage.update_info` fest, um eine zsync-Datei anzufordern, und `packages.appimage.checksum: true` für eine SHA256-Datei. Begleitdateien werden aus dem Paket-Cache wiederhergestellt und im Artefaktbeleg zusammen mit dem AppImage aufgeführt. Die Paketierung unter Windows und macOS wurde mit einem bereitgestellten GTK/WebKit-Verzeichnisbaum verifiziert; die daraus resultierenden reinen Go-GUI-Images haben Runtime-Prüfungen unter Linux bestanden. Bei einem vergleichbaren CGo-GUI-Build benötigte der native Pfad 8.77 Sekunden, linuxdeploy dagegen 26.19 Sekunden. Das native gzip-Image war 133 MiB groß, gegenüber 114 MiB beim alten zstd-Pfad. Gemessen wird daher der gesamte Paketierungspfad und nicht nur die Komprimierung.

Arch-Pakete bewahren Vorabversionen in den nativen Paketmetadaten; beispielsweise wird `1.2.3-beta.2` zu `1.2.3beta.2-1`, das vor der endgültigen Version einsortiert wird. Void-Abhängigkeiten akzeptieren reine Paketnamen wie `gtk4`, die zu `gtk4>=0` normalisiert werden; explizite Versionsvergleiche und vollständige `name-version_revision`-Werte bleiben erhalten. Verwenden Sie Void-Paketnamen anstelle von Bibliotheksdateinamen.

Für die Signierung von DEB-, RPM- und Arch-Paketen können auf jedem unterstützten Desktop-Host exportierte ASCII-armierte oder binäre private OpenPGP-Schlüssel verwendet werden, ohne GnuPG, `dpkg-sig` oder `rpmsign`. Konfigurieren Sie `linux.signing.certificate` mit der Schlüsseldatei und `credential` mit dem Namen der Umgebungsvariable für deren Passphrase. DEB und RPM betten Signaturen ein; Arch veröffentlicht eine abgetrennte Signatur. Private Schlüssel werden nur während des Signierschritts gelesen und weder in Paketressourcen noch in wiederverwendbare Cache-Artefakte kopiert.

Bestehende DEB/RPM-Konfigurationen, die lediglich eine GnuPG-Schlüsselkennung enthalten, verwenden weiterhin ihren externen Signierer. Um diese Abhängigkeit zu entfernen, exportieren Sie den geheimen Schlüssel und konfigurieren stattdessen dessen Pfad. Verwenden Sie einen Pfad wie `./signing/release.asc`; fehlende, nicht lesbare oder ungültige Schlüsseldateien führen zu einem Fehler, statt auf einen Schlüsselbund zurückzugreifen. DEB-Signaturen behalten das Paketsignaturformat `dpkg-sig` bei; dies ersetzt nicht die Signierung der Release-Metadaten eines apt-Repositorys.

Verwenden Sie zum Erstellen eines Arch-Pakets und einer abgetrennten Signatur eine exportierte, ASCII-armierte PGP-Datei mit privatem Schlüssel und geben Sie optional den Namen einer Umgebungsvariablen für die Passphrase an:

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

Setzen Sie die angegebene Umgebungsvariable in Ihrer lokalen Umgebung oder im Secret-Speicher Ihrer CI und führen Sie anschließend `wails3 build release` aus. Lassen Sie `credential` bei einem unverschlüsselten Schlüssel weg. Das Go-Signierprogramm läuft auf jedem unterstützten Desktop-Host und benötigt kein GPG. Es veröffentlicht das ursprüngliche `.pkg.tar.zst` und eine gleichgeordnete `.pkg.tar.zst.sig`-Datei; beide werden separat im Artefaktbeleg erfasst. Schlüsselinhalte werden weder im Staging-Bereich noch im Build-Cache gespeichert. Signaturen und deren Veröffentlichung werden niemals aus dem Cache wiederverwendet. Wenn das Signieren fehlschlägt, bleibt das zuvor veröffentlichte Dateipaar erhalten. Schlägt die Installation einer der beiden Dateien fehl, setzt die Veröffentlichung beide Dateien zurück. Dadurch werden zwei Dateisystemumbenennungen für nebenläufige Leser nicht atomar.

macOS-App-Bundles enthalten Begleitprogramme. Legen Sie `darwin.microphone_usage_description` fest, um den Nutzungshinweis und, sofern keine benutzerdefinierte Entitlements-Datei bereitgestellt wird, das Entitlement für die Audioeingabe hinzuzufügen. Benutzerdefinierte plist- und Entitlements-Quelldateien bleiben unverändert. Eine exportierte RSA-PKCS#12-Identität verwendet auf jedem Desktop-Host die Go-Signierung; `identity: "-"` wählt die portable Ad-hoc-Signierung aus. Schlüsselbundidentitäten verwenden unter macOS weiterhin `codesign`. Das Anheften von Tickets erfolgt in Go. Für die Notarisierungseinreichung kann ein vorhandenes macOS-Notarisierungsprofil im Schlüsselbund oder ein Satz von App-Store-Connect-Team-API-Zugangsdaten verwendet werden:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

Im API-Modus erfolgen Einreichung, Upload und Statusabfrage ohne `notarytool` in Go. Legen Sie alle drei Felder fest und lassen Sie `credential` weg; nehmen Sie den privaten Schlüssel nicht in die Versionsverwaltung auf. Der Schlüssel wird nur beim Signieren gelesen und weder in generierte Ressourcen kopiert noch zwischengespeichert. Dieser Modus unterstützt Team-API-Schlüssel und Archive bis zu 5 GiB; mehrteilige Uploads und individuelle API-Schlüssel sind nicht implementiert. Bei Ablehnungen bleibt die Einreichungs-ID erhalten, und Apples Diagnoseprotokoll wird in einer privaten temporären Datei gespeichert, deren Pfad in der Fehlermeldung angegeben wird. Eine abgebrochene oder abgelehnte Einreichung ersetzt keine zuvor veröffentlichte Ausgabe. Nach der Annahme ruft Wake Apples Ticket ab, prüft es und heftet es in Go an. Mit einer exportierten Identität und Team-API-Zugangsdaten können Signierung, Einreichung und Anheften unter Linux, Windows oder macOS erfolgen. Schlüsselbundprofile erfordern weiterhin macOS und `notarytool`.

Die portable Zertifikatsignierung unterstützt RSA-Identitäten in AES/PBKDF2- oder älteren 3DES-PKCS#12-Dateien. RC2-verschlüsselte Exporte müssen mit einem unterstützten Verschlüsselungsverfahren erneut exportiert werden. Notarisierte portable Signaturen erfordern einen expliziten SHA-2-RFC-3161-Zeitstempeldienst; der Go-Prüfer lehnt ältere SHA-1-Zeitstempelsignaturen ab. Der Ticketdienst folgt dem beobachteten Protokoll von Apple und ist an Apple Root G3 gebunden; Gatekeeper bleibt die maßgebliche Instanz für Vertrauenswürdigkeit und Widerruf bei der Distribution.

Signierte Ausgaben behalten ihre plattformspezifische Dateinamenserweiterung bei, zum Beispiel `notes.signed.app` und `notes.signed.dmg`.

## Entwicklungssteuerung

Geben Sie während `wails3 dev` einen Befehl ein und drücken Sie anschließend die Eingabetaste:

| Befehl | Aktion |
| --- | --- |
| `r` | Neuen Build und Neustart erzwingen, einschließlich eines unveränderten Backends. |
| `b` | Bindings neu generieren und erneut bauen. |
| `f` | Frontend neu starten. |
| `c` | Terminal leeren. |
| `q` | Sitzung und ihre untergeordneten Prozesse beenden. |
| `h` | Steuerbefehle anzeigen. |

Bei fehlgeschlagenen Builds läuft die aktuelle Anwendung weiter. Die Steuerbefehle erfordern ein interaktives Terminal; weitergeleitete Eingaben ändern das Sitzungsverhalten nicht.

## Arbeit vor der Ausführung prüfen

Die Planprüfung löst denselben unveränderlichen Graphen auf, den auch die Ausführung verwendet. Text ist für Menschen bestimmt; JSON ist ein versionierter, maschinenlesbarer Vertrag.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Konfiguration ohne Build validieren

Verwenden Sie den dedizierten Preflight-Befehl in der CI, als Editor-Aufgabe oder bevor Sie Konfigurationsänderungen committen. Ohne Profil validiert er das Basismanifest und jedes benannte Profil; mit einem Profil prüft er nur diese vollständige Anforderung:

```bash
wails3 config check
wails3 config check release
```

Fehler der Syntax, des Schemas, der Semantik, der Pfade, des Ziels, des Formats, der Signierung, des Planers und der von der Konfiguration festgelegten Host-Fähigkeiten verwenden denselben Diagnose-Renderer wie Build und Dev. Diagnosen nennen das Feld und den Quellbereich, geben die Quellzeile mit einem Markierungszeichen aus und enthalten gegebenenfalls einen hilfreichen Korrekturhinweis. Umgebungswerte werden in den Quellauszügen der Diagnose unkenntlich gemacht. Hostfehler, die ausschließlich durch Befehlszeilenoptionen verursacht werden, bleiben gewöhnliche Befehlsfehler und werden nicht fälschlich als ungültiges YAML gemeldet.

## Projekteigene Assets verwenden

Manifestpfade verweisen auf Dateien, die Ihnen gehören, beispielsweise Symbole, Plattformmanifeste, Paketvorlagen und Maintainer-Skripte. Wails bearbeitet diese Quelldateien niemals. Bevor ein Werkzeug sie verarbeitet, kopiert Wails sie in einen temporären `.wails/`-Arbeitsbereich. Die endgültigen Artefakte werden unter `build.output` abgelegt und ein verifizierter Beleg wird in `.wails/artifacts/receipt.json` geschrieben.

## Auf Android bereitstellen

Android-Produktionsprofile erzeugen AABs. Für die lokale Bereitstellung erstellt Wails intern eine ABI-spezifische Debug-APK, installiert sie mit ADB und startet die App:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Ohne explizite Auswahl bietet ein interaktives Terminal verbundene Geräte und konfigurierte Android Virtual Devices an. Bei der Automatisierung mit mehreren verfügbaren Zielen muss `--device` oder `--emulator` angegeben werden.

Hängen Sie einen auf die Anwendung beschränkten Log-Datenstrom an und überwachen Sie das Ziel, bis die App beendet wird oder Sie Strg+C drücken:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Angehängte Ausführungen folgen einem Neustart der Anwendung und melden, wenn die Verbindung zum Ziel getrennt wird oder das Ziel offline geht. Für schnelle Iterationen laufen Emulatoren standardmäßig weiter. `--stop-emulator` muss explizit angegeben werden, erfordert `--logs` und beendet beim Ende der angehängten Sitzung ausschließlich ein Emulatorziel. Wenn Wails einen neuen Emulator startet, dessen Startbereitschaft jedoch fehlschlägt oder abgebrochen wird, bereinigt Wails diesen Prozess automatisch.

## Begrenzte Lebenszyklus-Hooks ausführen

Die experimentelle Hook-Schnittstelle ruft in jeder deklarierten Phase ein projekteigenes Skript auf. Ein Hook ist kein Shell-Befehl: Wails führt die Datei direkt aus, standardmäßig vom Projektstammverzeichnis aus, sofern `directory` nicht gesetzt ist.

Auf Unix-Hosts muss das Skript ausführbar sein und mit einem Shebang beginnen. Windows-Hosts akzeptieren `.cmd`-, `.bat`- und `.ps1`-Skripte über festgelegte Interpreter-Adapter; Wails analysiert auf keiner der beiden Plattformen einen eingebetteten Shell-Befehl.

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

Unterstützte Phasen sind `before_build`, `after_build`, `before_package`, `after_package`, `before_sign` und `after_sign`. Build-Hooks werden je nach Bedarf einmal für das gemeinsam genutzte Projekt oder einmal pro Ziel ausgeführt. Paketierungs- und Signier-Hooks bilden Barrieren um alle für ein Ziel angeforderten Formate, während die Formatoperationen innerhalb der Barriere weiterhin nebenläufig ausgeführt werden.

Hooks erhalten eine Wails-eigene Umgebungsvariable: `WAILS_HOOK_CONTEXT_FILE`. Sie benennt eine geschützte JSON-Datei, die den versionierten Aufrufkontext enthält:

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

Der Kontext enthält aufgelöste Pfade, kopiert jedoch niemals die geerbte Prozessumgebung oder Secrets. Wails erstellt ihn atomar mit Dateiberechtigungen ausschließlich für den Eigentümer unter `.wails/hooks/context/`. Nach erfolgreicher Ausführung oder einem Abbruch wird er entfernt. Bei einem Fehler bleibt er zur Diagnose erhalten und sein Pfad wird in die Fehlermeldung aufgenommen. Ein Exitcode ungleich null lässt den regulären Plan-Schritt fehlschlagen; bei einem Abbruch wird die Prozessgruppe des Hooks beendet.

Hooks werden standardmäßig bei jedem Aufruf ausgeführt. Aktivieren Sie die Zwischenspeicherung nur, wenn das Skript vollständig deklarierte Ein- und Ausgaben besitzt:

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

Das Skript und sein Ausführbarkeitsmodus sind implizite Eingaben. Mehrere Ausgaben müssen sich in einem gemeinsamen Verzeichnis unterhalb des Stammverzeichnisses befinden. Dieses Verzeichnis ist das begrenzte Artefakt, das Wails erfasst und wiederherstellt. Das Ausgabestammverzeichnis darf weder das Skript noch eine deklarierte Eingabe enthalten. Zwischengespeicherte Ausgabebäume dürfen ausschließlich reguläre Dateien und Verzeichnisse enthalten, niemals symbolische Links oder Spezialdateien. Ein immer ausgeführter Hook verhindert die Zwischenspeicherung nachgelagerter Operationen, da nicht deklarierte Seiteneffekte nicht sicher mit einem Fingerabdruck versehen werden können. Ein zwischenspeicherbarer Hook erhält durch die deklarierte Artefaktidentität die Zwischenspeicherung nachgelagerter Operationen aufrecht.

## Ein vorhandenes Projekt migrieren

Erstellen und prüfen Sie einen inaktiven Entwurf:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Konservative, argumentlose Lebenszyklus-Tasks, die ein projekteigenes Skript aufrufen, werden in Hooks übersetzt. Anderes erreichbares benutzerdefiniertes Taskfile-Verhalten wird als Blocker gemeldet. Aktivieren Sie den Entwurf nach der Prüfung von `wails.migrated.yaml` ausdrücklich:

```bash
wails3 migrate --activate
```

Bei der Aktivierung wird `wails.yaml` atomar erstellt. Taskfiles, `build/config.yml`, Skripte oder Assets werden dabei weder gelöscht noch umbenannt oder bearbeitet. Ungültiges YAML führt zu einem YAML-Fehler; es wird niemals auf den Legacy-Pfad zurückgegriffen.

## Alle Standardwerte materialisieren

Schreiben Sie einen inaktiven Snapshot, um die vollständig aufgelöste Konfiguration zu prüfen oder zu kopieren:

```bash
wails3 eject
```

Dadurch wird `wails.ejected.yaml` erstellt. Die aktive Konfiguration wird niemals neu geschrieben.

## Kompatibilitätszeitraum

`wails3 package` und `wails3 sign` sind veraltete Aliasse für denselben Planer; neue Automatisierungen sollten `wails3 build` mit einem Profil oder mit Formaten verwenden. Die Einführung ist an Releases gebunden: Im ersten Release wird YAML zum Pfad für neue Projekte, das darauffolgende Release warnt bei Legacy-Routing über Taskfiles und das nächste entfernt die Taskfile-Ausführung zusammen mit den Aliassen. Für jeden Übergang muss die Akzeptanzmatrix der nativen Hosts für das jeweilige Release vollständig erfolgreich sein.

Alle Felder, Typen, Standardwerte und Anwendbarkeitsregeln finden Sie in der generierten [`wails.yaml`Feldreferenz](/reference/wails-yaml-fields/).

Die [`wails3 run`Konfigurationsreferenz](/reference/wails-yaml-run/) beschreibt die Standardwerte für den lokalen Start, plattformspezifische Überschreibungen und den `go run .`-Fallback ohne Manifest.
