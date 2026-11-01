---
title: "Wake"
description: "Erstelle dein erstes Wake-Projekt, durchlaufe den Einrichtungsassistenten, baue es, führe es aus und migriere ein bestehendes Taskfile-Projekt."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake ist das experimentelle Build-System von Wails. Eine `wails.yaml` beschreibt deine Anwendung, das Frontend, die Zielplattformen, Pakete und Signierungseinstellungen. Wake verwendet diese Konfiguration, um Builds zu planen und auszuführen, unveränderte Arbeit zwischenzuspeichern und die erzeugten Dateien aufzulisten.

In diesem Tutorial erstellst du eine kleine Desktop-Anwendung, konfigurierst sie im Browser, baust sie für deinen eigenen Computer und führst sie im Entwicklungsmodus aus. Beginne mit einem lokalen Ziel und einem unsignierten Build. Füge Cross-Kompilierung und Distribution hinzu, nachdem dieser erste Build funktioniert.

Du hast bereits eine Anwendung? Gehe zu [Ein bestehendes Projekt migrieren](#ein-vorhandenes-projekt-migrieren). Informationen zu einzelnen Feldern und fortgeschrittenen Arbeitsabläufen findest du im [YAML-Build-Leitfaden](/experimental/yaml-builds/) und in der [Manifestreferenz](/reference/wails-yaml-fields/).

## Was du von diesem Experiment erwarten kannst

@note{type="caution" title="Verwende ein Testprojekt oder einen separaten Branch"}
Wake muss ausdrücklich aktiviert werden. Seine Konfiguration, Befehle und sein Verhalten können sich zwischen Vorschau-Builds ändern. Es bietet noch nicht die Kompatibilitätsgarantien eines stabilen Build-Systems. Bewahre eine nachweislich funktionierende Revision deines Projekts auf und dokumentiere die getesteten CLI- und Framework-Versionen. Teste Pakete vor der Distribution auf ihren Zielsystemen.
@end

Wake und das **Pure-Go**-Backend sind voneinander unabhängige Optionen. Wenn du Wake aktivierst, wird Pure Go nicht automatisch aktiviert. In diesem Tutorial wird das Standard-Backend verwendet; GUI-Builds für macOS und Linux benötigen daher weiterhin die üblichen nativen Entwicklungswerkzeuge.

Der Assistent prüft die Anforderungen und speichert die Konfiguration. Das Speichern ist weder ein Build noch ein Nachweis dafür, dass jedes ausgewählte Ziel bereit ist. Fehlende Werkzeuge oder Anmeldedaten können später ergänzt werden, beispielsweise auf einem CI-Host. Ungültige Kombinationen aus Zielen und Paketen müssen vor dem Speichern korrigiert werden.

Die Migration ist eine **Einbahnstraße**. Wake führt keine Taskfiles aus, expandiert keine Task-Vorlagen und übersetzt YAML nicht zurück in Taskfiles. Nicht migrierte Projekte verwenden weiterhin die separate alte Task-Laufzeit. Eine aktive `wails.yaml` greift niemals stillschweigend darauf zurück.

## Bevor du beginnst

Du benötigst Go, Node.js mit npm für das Starter-Frontend und die [Plattformabhängigkeiten](/getting-started/installation/#platform-specific-dependencies) für deinen Computer. Beginne in einem Verzeichnis, in dem der neue Projektordner erstellt werden soll. Für den Assistenten benötigst du außerdem einen Browser.

Verwende eine CLI- und Framework-Revision, die diese Vorschau enthält. Durch die Installation der stabilen `v3.0.0`-CLI werden die hier beschriebenen Wake-Einrichtungsseiten nicht hinzugefügt. Die folgenden Befehle fixieren die Vorschau aus [PR Nr. 6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@0243da963
wails3 version
```

Stelle sicher, dass sich das Verzeichnis mit der ausführbaren Go-Datei in `PATH` befindet; die [Installationsanleitung](/getting-started/installation/) erklärt dies. Später im Tutorial fixierst du die Go-Abhängigkeit der Anwendung auf dieselbe Revision. Aktualisiere beim Testen einer neueren Vorschau immer beide Revisionsargumente gemeinsam.

### Wake in deinem Terminal aktivieren

In Bash oder Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Lasse dieses Terminal für den Rest des Tutorials geöffnet. Setze die Variable in jedem neuen Terminal und jedem CI-Job erneut, der Wake-Befehle ausführt. Bereits ihre **Existenz** aktiviert Wake: Wenn du sie auf `false` oder eine leere Zeichenfolge setzt, wird Wake nicht deaktiviert. Der alte Name `WAILS_USE_WAKE` wird ignoriert.

## Anwendung erstellen

Führe den Assistenten zur Projekterstellung aus:

```sh
wails3 init -ui
```

Im Browser wird ein Projektassistent geöffnet. Dieser erstellt die Quelldateien; mit dem separaten `wails3 setup`-Assistenten weiter unten konfigurierst du, wie sie gebaut werden.

| Seite | Auswahl für dieses Tutorial | Was als Nächstes geschieht |
| --- | --- | --- |
| **Projekt – Wails-Projekt erstellen** | Gib `wake-demo` ein und klicke dann auf **Weiter**. | Das Projekt wird unterhalb des angezeigten Verzeichnisses in einem neuen Ordner namens `wake-demo` erstellt. |
| **Framework – Framework auswählen** | Wähle **Vanilla** und klicke dann auf **Weiter**. | Du erhältst ein kleines Frontend, ohne dich zunächst in ein bestimmtes Framework einarbeiten zu müssen. |
| **Sprache – Spracheinstellung** | Wähle **JavaScript** und klicke dann auf **Weiter**. | Die ausschließlich für TypeScript vorgesehene Seite **Bindungen** wird übersprungen. |
| **Details – Projektdetails** | Verwende `Wake Demo` als Produktnamen, `com.example.wakedemo` als Bundle-ID und `0.1.0` als Version. Klicke auf **Projekt erstellen**. | Wenn Wake aktiviert ist, werden die Details in `wails.yaml` geschrieben. |

Verwende für dieses Tutorial nicht **Globale Standardeinstellungen verwenden**, da dadurch die Auswahl von Framework und Sprache übersprungen wird. Wenn du in einem anderen Projekt TypeScript auswählst, kannst du auf der Seite **Bindungen** zwischen Schnittstellen und Klassen für generierte Go-Datenbindungen wählen.

Kehre zum Terminal zurück und warte, bis die Projekterstellung abgeschlossen ist. Führe dann Folgendes aus:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@0243da963
```

Du solltest nun über `main.go`, ein Go-Modul, ein `frontend`-Verzeichnis und eine native `wails.yaml` verfügen. Die integrierte Vorlage verwendet Wake direkt; du musst dieses neue Projekt nicht migrieren. Eine Community-Vorlage kann eigene alte Tasks enthalten und stattdessen die weiter unten beschriebene Migrationsprüfung erfordern.

## Build Seite für Seite konfigurieren

Führe in `wake-demo` Folgendes aus:

```sh
wails3 setup
```

Im Browser wird nun **WAKE-EINRICHTUNG** angezeigt. Wenn stattdessen der allgemeine Einrichtungsablauf erscheint, prüfe, ob `WAILS_EXP_USE_WAKE` in dem Terminal gesetzt ist, aus dem du den Befehl gestartet hast, und ob du die Vorschauversion der CLI ausführst.

Die folgenden Screenshots zeigen diesen Ablauf unter Linux. Dein Host-Ziel und die erforderlichen Werkzeuge unterscheiden sich unter Windows oder macOS.

### Projekt – nur bei fehlenden Metadaten

Bei der gerade erstellten Anwendung wird diese Seite **übersprungen**, da `init` die Metadaten bereits geschrieben hat. Das ist zu erwarten.

Wenn du ein bestehendes Go-Quellverzeichnis ohne natives Manifest oder alte Build-Konfiguration konfigurierst, fragt **Über dein Projekt** nach Produktname, Projektname, Anwendungs-ID, Version, Unternehmen und Beschreibung. Diese Werte verbleiben bis zum **Speichern** im Assistenten. Die Einrichtung erstellt keine Anwendungsquelldateien. Wenn Taskfiles oder alte Metadaten gefunden werden, führe zuerst die Migration durch.

### Anwendung – Was möchtest du bauen?

Wähle **Desktop-Anwendung**. Lasse **Vorhandene Binärdateien paketieren** deaktiviert und wähle **Standard-Backend**. Klicke auf **Weiter**.

![Anwendungsseite mit ausgewählter Desktop-Anwendung und ausgewähltem Standard-Backend; das Paketieren vorab erstellter Binärdateien ist deaktiviert.](/assets/screenshots/wake/application.png)

- **Desktop-Anwendung** baut eine GUI mit einem Web-Frontend.
- **Befehlszeilenanwendung** und **Server** bezeichnen Programme ohne Desktop-Fenster. Wenn du diese Auswahl änderst, wird der Quellcode deiner Anwendung nicht umgeschrieben. Verwende ein Quellprojekt, das dem gewählten Typ entspricht.
- **Pure Go (experimentell)** muss separat aktiviert werden. Für unterstützte Desktop-Kompilierungen entfällt damit die Anforderung eines C-Compilers, sofern sich auch deine Abhängigkeiten ohne CGo bauen lassen. Ziellaufzeitbibliotheken sowie sämtliche Werkzeuge zum Paketieren und Signieren werden dadurch nicht überflüssig.
- **Vorhandene Binärdateien paketieren** bereitet einen Befehl mit `--skip-build` vor.
  Dafür sind bereits erstellte Binärdateien erforderlich; lassen Sie diese Option bei einem neuen Projekt deaktiviert.

Das Backend ist eine **projektweite Einstellung**. Wenn Sie sie hier ändern, kann sich das Verhalten anderer Profile ändern, die diese Einstellungen erben.

### Zielplattformen — Wo soll die Anwendung ausgeführt werden?

Behalten Sie nur die Zielplattform des verwendeten Computers bei und klicken Sie dann auf **Weiter**.
Der Assistent wählt zunächst das Betriebssystem und die Architektur des Hostsystems aus.

![Ansicht „Zielplattformen“, in der für das Hostsystem des Tutorials nur linux/amd64 ausgewählt ist.](/assets/screenshots/wake/targets.png)

Zielnamen kombinieren ein Betriebssystem mit einer CPU-Architektur, zum Beispiel `windows/amd64`, `darwin/arm64` für einen Apple-Silicon-Mac oder `linux/amd64`.
`darwin` steht für macOS. Wenn Sie eine andere Zielplattform auswählen, fordern Sie eine ausführbare Datei für diesen Computer an; dadurch wird die Datei nicht auf Ihrem aktuellen Computer ausführbar.

Zusätzliche Architekturen, universelle macOS-Builds und experimentelle mobile Zielplattformen können Sie bei einem späteren Setup-Durchlauf auswählen.

### Pakete — Release-Ausgaben auswählen

Lassen Sie die Kontrollkästchen für Paketformate leer und **Release-Ausgaben signieren** deaktiviert.
Klicken Sie auf **Weiter**.

![Ansicht „Pakete“ ohne ausgewählte zusätzliche Linux-Paketformate und mit deaktivierter Signierung.](/assets/screenshots/wake/packages.png)

Wenn keine zusätzlichen Formate ausgewählt sind, erzeugt Wake die Standardausgabe der Zielplattform: eine ausführbare Desktopdatei oder bei einer macOS-GUI ein macOS-Anwendungspaket. Dies genügt, um die Anwendung lokal zu prüfen. Wenn Sie später Formate auswählen, werden – sofern unterstützt – Distributionsausgaben wie DMG, MSIX, NSIS, DEB oder AppImage hinzugefügt.

Für diesen ersten lokalen Build ist die Signierung optional. Ein erfolgreicher unsignierter Build bedeutet nicht, dass die Anwendung zur öffentlichen Verteilung bereit ist.

### Werkzeuge — Werkzeuge für Ihren Build

Lesen Sie die Anforderungen für Ihre Auswahl. Die Liste stammt aus dem Build-Plan von Wake; sie fordert nicht dazu auf, jedes Werkzeug auf jeder Plattform zu installieren.

![Ansicht „Werkzeuge“ mit den ermittelten Anforderungen für den lokalen Linux-Build aus dem Quellcode.](/assets/screenshots/wake/tools.png)

Für diese Lektion benötigen Sie voraussichtlich Go und die Frontend-Werkzeuge sowie die nativen Entwicklungswerkzeuge, die das Standard-GUI-Backend Ihres Hostsystems erfordert. Docker ist keine pauschale Voraussetzung. Ein Cross-Build mit CGo kann einen Compiler und Zielbibliotheken oder eine unterstützte Container-Toolchain erfordern.

Beheben Sie vor dem lokalen Build alle Zeilen mit dem Status **Nicht bereit**, **Konfiguration erforderlich** oder **Aktualisierung erforderlich**. Verwenden Sie, sofern verfügbar, die angezeigten Anweisungen oder die [Installationsanleitung für die Plattform](/getting-started/installation/) und klicken Sie anschließend auf **Erneut prüfen**. Möglicherweise müssen Sie das Setup neu starten, wenn eine Installation den Wert `PATH` des Terminals geändert hat. Wenn die Anforderungen des Hostsystems erfüllt sind, klicken Sie auf **Weiter**.

Der Assistent kann auch dann eine gültige Konfiguration speichern, wenn auf diesem Computer Werkzeuge fehlen. Das ist für die CI-Konfiguration nützlich; in diesem Tutorial benötigen Sie jedoch einen Build, der hier ausgeführt werden kann.

### Signierung — nur bei Auswahl

Diese Ansicht wird **übersprungen**, weil Sie die Signierung deaktiviert gelassen haben. Wenn Sie bei einem späteren Durchlauf **Release-Ausgaben signieren** aktivieren, werden die vorhandenen Signierungsoptionen für die ausgewählten Desktopplattformen angezeigt. **Signierung später einrichten** deaktiviert die Signierung für das gespeicherte Setup-Profil. Die mobile Signierung wird in YAML und nicht über diese Optionen konfiguriert.

Die Standardeinstellungen für die Signierung befinden sich in Ihren Benutzereinstellungen. Wake kann konfigurierte Zertifikatdateien außerhalb des Projekts verwenden, ohne private Schlüssel oder Passwörter nach `wails.yaml` zu kopieren. Explizite Signierungseinstellungen des Projekts haben Vorrang. Befolgen Sie vor der Vorbereitung eines öffentlichen Releases die [Anleitung zur Einrichtung der Signierung](/getting-started/setup/) und den [Leitfaden zur Paketierung und Signierung](/experimental/yaml-builds/#package-and-sign-for-distribution).

### Speichern — Build-Konfiguration speichern

Prüfen Sie die Zusammenfassung: Desktopanwendung, Standard-Backend, Build aus dem Quellcode, Ihre Host-Zielplattform und deaktivierte Signierung. Klicken Sie auf **wails.yaml speichern**.

![Ansicht „Speichern“ mit einer Zusammenfassung der GUI-Anwendung, des Standard-Backends, des Builds aus dem Quellcode und der unsignierten Zielplattform linux/amd64.](/assets/screenshots/wake/save.png)

Nun sollten **Build-Konfiguration gespeichert**, der Pfad zu Ihrem Manifest und Folgendes angezeigt werden:

```sh
wails3 build setup
```

Klicken Sie auf **Fertigstellen** und kehren Sie zum Terminal zurück. Das Setup hat die Auswahl der Zielplattformen und Pakete unter `profiles.setup` gespeichert. Außerdem aktualisiert es den Anwendungstyp sowie die Einstellungen für die Frontend-Aktivierung und das Backend. Vorhandene Frontend-Befehle, Begleitprogramme, Hooks und nicht zugehörige Konfigurationen bleiben erhalten.

Der Profilname ist wichtig: `wails3 build setup` verwendet diese gespeicherte Auswahl.
Der einfache Befehl `wails3 build` wählt das Profil `setup` nicht automatisch aus.

## Erste Anwendung erstellen und ausführen

Führen Sie den vom Assistenten angezeigten Befehl aus:

```sh
wails3 build setup
```

Beim ersten Build werden die Frontend-Abhängigkeiten installiert, Bindings generiert und die Anwendung kompiliert. Spätere Builds können unveränderte Arbeit wiederverwenden. Warten Sie auf **Build erfolgreich** und beachten Sie die von der CLI ausgegebenen Ausgabepfade. Das Standard-Ausgabeverzeichnis ist `bin`; eine macOS-GUI wird als `.app`-Paket erstellt.

Starten Sie den Entwicklungsmodus auf Ihrem eigenen Computer:

```sh
wails3 dev
```

Nun sollte das Fenster der Starteranwendung angezeigt werden. Geben Sie Ihren Namen ein und testen Sie die Begrüßung. Ändern Sie den Platzhalter des Namenseingabefelds in `frontend/index.html` und speichern Sie die Datei; das Frontend sollte aktualisiert werden. Bearbeiten Sie die Begrüßung in `greetservice.go`, um einen erneuten Go-Build auszulösen. Drücken Sie anschließend im Terminal **Ctrl+C**.

Die Entwicklung erfolgt lokal auf dem Hostsystem. Ein Release-Profil kann mehrere Cross-Build-Zielplattformen enthalten, ohne zu ändern, auf welchem Computer das Entwicklungsfenster ausgeführt werden kann. Verwenden Sie `wails3 run`, um eine Produktionsanwendung lokal zu erstellen und zu starten; weitere Informationen finden Sie in der [Referenz zum Ausführen](/reference/wails-yaml-run/).

Sie verfügen nun über ein Quellprojekt, ein geprüftes Build-Profil und einen funktionierenden lokalen Entwicklungszyklus. Übernehmen Sie den Quellcode und `wails.yaml` in die Versionsverwaltung. Die generierten `.wails/`-Arbeitsdateien sind nicht die Quelle Ihrer Konfiguration.

## Nächste Schritte

| Ziel | Unterstützung durch Wake | Zuerst zu prüfen |
| --- | --- | --- |
| Für einen anderen Desktop erstellen | Zielplattformen für Windows, macOS und Linux, einschließlich beider macOS-Architekturen und universeller Ausgabe. | Pure Go und CGo haben unterschiedliche Compileranforderungen. Testen Sie auf dem Zielbetriebssystem und der Zielarchitektur. |
| Installationsprogramme oder Pakete erstellen | MSIX, NSIS, macOS-App/DMG, DEB, RPM, Arch, Void und AppImage. | NSIS benötigt `makensis`. AppImage benötigt eine Ziel-Laufzeitumgebung; die plattformübergreifende GUI-Paketierung erfordert außerdem einen bereitgestellten Linux-Laufzeitverzeichnisbaum. |
| Ein Release signieren | Portable Signierung mit exportierten Schlüsseln für Windows, macOS und DEB/RPM/Arch; native Identitätsspeicher bleiben verfügbar. | Für jede Identität, jedes Format und jedes Hostsystem gelten Einschränkungen. Für AppImage und Void gibt es keinen integrierten Signierungsablauf. |
| Eine macOS-App notarisieren | Übermittlung mit einem Team-API-Schlüssel und Anheften des Tickets oder ein macOS-Schlüsselbundprofil. | Eine geeignete Developer-ID-Signatur und ein Zeitstempel sind weiterhin erforderlich. Dass von Go generierte Testsignaturen native Prüfungen bestehen, belegt nicht, dass Apple Ihr Release akzeptiert. |
| Mehrere Programme erstellen | Ein primäres GUI-, CLI- oder Serverprogramm sowie benannte Begleitprogramme. | Deklarieren Sie die Pakete und Typen der Begleitprogramme in YAML. Der Assistent generiert deren Quellcode nicht. |
| Andernorts erstellte Dateien paketieren | `wails3 build setup --skip-build`. | Geben Sie in `build.output` übereinstimmende Namen und Architekturen der ausführbaren Dateien an. Dieses Flag ist eine einmalige Befehlsoption und keine dauerhafte Konfigurationseinstellung. |
| Release-Automatisierung hinzufügen | Benannte Profile, begrenzte Lebenszyklus-Hooks, Caching und parallele Paketverarbeitung. | Verwenden Sie literales YAML und deklarierte, zum Projekt gehörende Skripte. Beliebiges Taskfile-Verhalten ist kein Erweiterungsmechanismus für Wake. |

Öffnen Sie `wails3 setup` erneut, um das Profil `setup` zu ändern. Andere benannte Profile und erweiterte Einstellungen bearbeiten Sie in `wails.yaml`. Vollständige Beispiele und Paketanforderungen finden Sie im [Leitfaden für YAML-Builds](/experimental/yaml-builds/). GUI-Anwendungen benötigen weiterhin die Systembibliotheken des Zielsystems; Pure Go macht sie nicht zu vollständig statischen Anwendungen.

iOS und Android sind weiterhin experimentell, benötigen ihre Plattform-SDKs und CGo und werden in dieser ersten Desktop-Einführung nicht behandelt. Der Assistent macht mobile Builds nicht auf beliebigen Hosts portierbar.

## Ein vorhandenes Projekt migrieren

Arbeiten Sie in einem separaten Branch oder mit einer Kopie eines funktionsfähigen **Wails-v3**Projekts. Dies ist eine Migration der Build-Konfiguration und kein Upgrade von Wails v2. Bewahren Sie das Ergebnis Ihres vorhandenen Builds und Ihre Tests auf, damit Sie das Verhalten nach der Aktivierung vergleichen können.

@note{type="note" title="Mit einer geprüften Migration beginnen"}
Neue Projekte, die mit den aktuellen integrierten Legacy-Vorlagen erzeugt wurden, können migriert werden, bevor ihr Frontend gebaut wurde. Wake erkennt deren unveränderte Task-Selektoren und Entwicklungslebenszyklus, ohne Taskfiles auszuführen oder benutzerdefinierte Ausdrücke auszuwerten.

Benutzerdefinierte Tasks, Variablenüberschreibungen und Entwicklungsrichtlinien müssen weiterhin geprüft werden. Wenn die Analyse Ihr Projekt nicht abbilden kann, verwenden Sie dafür weiterhin das Legacy-Buildsystem und melden Sie die Diagnoseergebnisse. Entfernen Sie keine benötigte Funktionalität, nur damit die Aktivierung erfolgreich ist. Für das obige Tutorial zu neuen Projekten ist keine Migration erforderlich.
@end

### 1. Analysieren, ohne den aktiven Build zu ändern

Aktivieren Sie Wake und verwenden Sie die am Anfang dieser Seite angegebene Vorschau-CLI. Führen Sie im Stammverzeichnis des vorhandenen Projekts Folgendes aus:

```sh
wails3 migrate --dry-run
```

Lesen Sie die gemeldeten unterstützten Einstellungen und Blocker. Das Analyseprogramm liest Legacy-Deklarationen, führt deren Befehle jedoch nicht aus. Ein Probelauf erstellt kein aktives Manifest. Verwenden Sie für maschinenlesbare Diagnosedaten `wails3 migrate --json`. Wenn außerdem kein Entwurf geschrieben werden soll, verwenden Sie `--dry-run --json`.

### 2. Inaktiven Vorschlag erstellen und prüfen

```sh
wails3 migrate
```

Dadurch wird **`wails.migrated.yaml`** erstellt, nicht `wails.yaml`. Prüfen Sie die Datei zusammen mit Ihren vorhandenen Taskfiles und `build/config.yml`. Kontrollieren Sie Projektmetadaten, Frontend-Befehle, Zielarchitekturen, Paketinhalte, Assets und Signierungseinstellungen. Speichern Sie Geheimnisse in Umgebungsvariablen oder in Ihrem Anmeldedatenspeicher.

Bekannte integrierte Tasks können nativen Einstellungen zugeordnet werden. Konservative Lebenszyklus-Tasks ohne Argumente, die ein zum Projekt gehörendes Skript aufrufen, können zu Hooks werden. Anderes erreichbares Verhalten benutzerdefinierter Tasks ist ein Blocker; Wake verwirft es nicht stillschweigend und führt es auch nicht ersatzweise über Task aus.

Wenn bereits ein Entwurf vorhanden ist, überschreibt die Migration Ihre geprüften Änderungen nicht. Erstellen Sie eine Kopie, bevor Sie ihn absichtlich ersetzen oder neu generieren. Benennen Sie einen blockierten Entwurf nicht in `wails.yaml` um, um die Aktivierungsprüfungen zu umgehen. Reduzieren oder ersetzen Sie nicht unterstütztes Legacy-Verhalten und führen Sie die Analyse vor der Aktivierung erneut aus. Wenn das Projekt noch nicht abgebildet werden kann, verwenden Sie weiterhin seinen Legacy-Build und melden Sie die Lücke.

### 3. Explizit aktivieren

Führen Sie nach der Prüfung des Entwurfs und der Behebung der Blocker Folgendes aus:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

Der erste Befehl prüft die Aktivierung, ohne sie zu veröffentlichen. Der zweite prüft die Legacy-Deklarationen erneut und validiert den Entwurf, bevor `wails.yaml` atomar erstellt wird. Bei erfolgreicher Aktivierung wird die inaktive Datei `wails.migrated.yaml` entfernt. Taskfiles, `build/config.yml`, Skripte und Assets bleiben unverändert. Ab diesem Zeitpunkt steuert das aktive Manifest die Wake-Builds. Die alten Dateien sind keine zweite Konfigurationsquelle.

Einige ältere Projekte enthalten bereits eine reine Metadatendatei `wails.yaml` mit einem Abschnitt `info:` auf oberster Ebene. Obwohl sie kein natives Wake-Manifest ist, belegt sie dennoch das Aktivierungsziel. Bewahren Sie diese Legacy-Datei auf und verschieben Sie sie vor der Aktivierung bewusst. Das Migrationsprogramm überschreibt oder benennt sie nicht für Sie um.

### 4. Migrierte Anwendung überprüfen

```sh
wails3 config check
wails3 build
wails3 dev
```

Vergleichen Sie die Anwendung und die Paketinhalte einschließlich Ihres benutzerdefinierten Release-Verhaltens mit dem vorherigen Build. Führen Sie Ihre Projekttests aus. Prüfen Sie jedes migrierte Profil ausdrücklich, beispielsweise `wails3 config check release` und `wails3 build release`, wenn der Entwurf ein Profil `release` enthält.

Sie können nun `wails3 setup` ausführen, um das dedizierte Profil `setup` hinzuzufügen oder zu ersetzen. Ein vorhandenes Profil `release` wird dadurch weder automatisch ausgewählt noch bearbeitet.

### Wenn Sie das Experiment beenden müssen

Wenn `WAILS_EXP_USE_WAKE` nicht mehr gesetzt ist, wird die explizite Aktivierung deaktiviert. Ein aktives `wails.yaml` blockiert jedoch weiterhin das automatische Routing über Taskfiles. Ungültiges YAML führt ebenfalls zu einem Abbruch mit einer Fehlermeldung. Es gibt keinen Befehl für eine Rückmigration.

Um zum alten Arbeitsablauf zurückzukehren, stellen Sie in der Versionsverwaltung die vollständige Projektversion vor der Migration wieder her, einschließlich des fehlenden aktiven Manifests, und verwenden Sie die dazu kompatible CLI und die entsprechenden Abhängigkeiten. Sichern Sie zuvor alle neuen Quellcodeänderungen, die Sie benötigen. Ein separater Test-Branch erleichtert diese bewusste Wiederherstellung erheblich.

## Probleme diagnostizieren und Feedback geben

| Symptom | Erste Prüfung |
| --- | --- |
| Setup zeigt den allgemeinen Assistenten an | Setzen Sie `WAILS_EXP_USE_WAKE` im Terminal, aus dem Sie den Prozess starten, und prüfen Sie, welche ausführbare Datei `wails3` über `PATH` gefunden wird. |
| Setup fordert zur Migration auf | Sie haben ein Legacy-Projekt oder eine Community-Vorlage mit Legacy-Konfiguration geöffnet. Prüfen Sie die Migration vor der Aktivierung. |
| Ein gespeicherter Build meldet fehlende Tools | Beim Speichern wird die beabsichtigte Konfiguration festgehalten. Kehren Sie zu **Tools** zurück, erfüllen Sie die Host-Anforderungen und prüfen Sie erneut. |
| Eine für ein anderes Ziel erstellte ausführbare Datei lässt sich hier nicht starten | Führen Sie sie auf dem ausgewählten Betriebssystem und der ausgewählten Architektur aus. Verwenden Sie für die lokale Entwicklung das Host-Ziel. |
| Im Prebuilt-Modus wird eine Binärdatei nicht gefunden | Prüfen Sie die [Dateinamen und Pfade der vorgefertigten Binärdateien](/experimental/yaml-builds/#use-prebuilt-desktop-executables), einschließlich der Unterverzeichnisse für mehrere Ziele. |
| Eine Manifeständerung schlägt fehl | Führen Sie `wails3 config check` für das betreffende Profil aus und folgen Sie der feldbezogenen Diagnose. Wake greift nicht auf Taskfiles zurück. |

Prüfen Sie einen Build vor der Ausführung mit `wails3 build setup --plan --json`.
`wails3 config check` validiert die Konfiguration, und `wails3 eject` schreibt einen
inaktiven `wails.ejected.yaml`-Snapshot, ohne das aktive Manifest zu ändern.
Prüfen Sie die Befehlsausgabe, bevor Sie sie weitergeben: Projektpfade und Speicherorte
von Anmeldedaten können Rückschlüsse auf Ihre lokale Umgebung zulassen.

Geben Sie Feedback in den [Diskussionen zu Experimenten](https://github.com/wailsapp/wails/discussions/categories/experiments) oder [melden Sie einen reproduzierbaren Fehler](https://github.com/wailsapp/wails/issues). Geben Sie die Revisionen Ihrer CLI und Ihres Frameworks, Betriebssystem und Architektur des Hosts, das ausgewählte Ziel, den genauen Befehl, die Auswahl von Pure Go sowie die relevante Diagnose an. Nennen Sie, ob der Fehler beim Kompilieren, Paketieren, Signieren oder Ausführen der Anwendung auftrat. Fügen Sie ein minimales Manifest ohne Geheimnisse und, falls hilfreich, die Ausgabe von `wails3 doctor` bei.
