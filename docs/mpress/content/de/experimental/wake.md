---
title: "Wake"
description: "Erstelle dein erstes Wake-Projekt, durchlaufe den Assistenten zur Projekterstellung, baue es, führe es aus und migriere ein bestehendes Taskfile-Projekt."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake ist das experimentelle Build-System von Wails. Eine `wails.yaml` beschreibt deine Anwendung, das Frontend, die Zielplattformen, Pakete und Signierungseinstellungen. Wake verwendet diese Konfiguration, um Builds zu planen und auszuführen, unveränderte Arbeit zwischenzuspeichern und die erzeugten Dateien anzugeben.

In diesem Tutorial erstellst du eine kleine Desktop-Anwendung, konfigurierst sie im Browser, baust sie für deinen eigenen Computer und führst sie im Entwicklungsmodus aus. Beginne mit einem lokalen Ziel und einem unsignierten Build. Füge Cross-Kompilierung und Distribution hinzu, nachdem dieser erste Build funktioniert.

Du hast bereits eine Anwendung? Gehe zu [Bestehendes Projekt migrieren](#ein-vorhandenes-projekt-migrieren). Informationen zu einzelnen Feldern und fortgeschrittenen Workflows findest du im [YAML-Build-Leitfaden](/experimental/yaml-builds/) und in der [Manifestreferenz](/reference/wails-yaml-fields/).

## Was du von diesem Experiment erwarten kannst

@note{type="caution" title="Verwende ein Testprojekt oder einen separaten Branch"}
Wake muss ausdrücklich aktiviert werden. Seine Konfiguration, Befehle und sein Verhalten können sich zwischen Vorschau-Builds ändern. Es bietet noch nicht die Kompatibilitätsgarantien eines stabilen Build-Systems. Bewahre eine nachweislich funktionierende Revision deines Projekts auf und notiere die getesteten CLI- und Framework-Versionen. Teste Pakete vor der Verteilung auf ihren Zielsystemen.
@end

Wake und das **Pure-Go**-Backend sind voneinander unabhängige Optionen. Wenn du Wake aktivierst, wird Pure Go nicht automatisch aktiviert. Diese Anleitung verwendet das Standard-Backend; GUI-Builds für macOS und Linux benötigen daher die üblichen nativen Entwicklungswerkzeuge.

Der Assistent prüft die Voraussetzungen und speichert die Konfiguration. Das Speichern ist kein Build und kein Nachweis dafür, dass jedes ausgewählte Ziel einsatzbereit ist. Fehlende Werkzeuge oder Anmeldedaten kannst du später ergänzen, beispielsweise auf einem CI-Host. Ungültige Kombinationen aus Ziel und Paket musst du vor dem Speichern korrigieren.

Die Migration ist eine **Einbahnstraße**. Wake führt keine Taskfiles aus, expandiert keine Task-Vorlagen und übersetzt YAML nicht zurück in Taskfiles. Nicht migrierte Projekte verwenden weiterhin die separate ältere Task-Laufzeit. Eine aktive `wails.yaml` greift niemals stillschweigend darauf zurück.

## Vorbereitung

Du benötigst Go, Node.js mit npm für das Starter-Frontend, die [Plattformabhängigkeiten](/getting-started/installation/#platform-specific-dependencies) für deinen Computer sowie einen Browser für den Assistenten. Beginne in dem Verzeichnis, in dem der neue Projektordner erstellt werden soll.

Verwende eine CLI- und Framework-Revision, die diese Vorschau enthält. Durch die Installation der stabilen `v3.0.0`-CLI werden die hier beschriebenen Wake-Ansichten zur Projekterstellung nicht hinzugefügt. Mit den folgenden Befehlen legst du die Vorschau aus [PR Nr. 6531](https://github.com/wailsapp/wails/pull/6531) fest:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@b1daf95f3
wails3 version
```

Stelle sicher, dass sich das Verzeichnis der Go-Programmdateien in `PATH` befindet; der [Installationsleitfaden](/getting-started/installation/) erläutert dies. Später im Tutorial legst du für die Go-Abhängigkeit der Anwendung dieselbe Revision fest. Wenn du eine neuere Vorschau testest, aktualisiere beide Revisionsargumente gemeinsam.

### Wake im Terminal aktivieren

In Bash oder Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Lasse dieses Terminal für den Rest des Tutorials geöffnet. Setze die Variable in jedem neuen Terminal oder CI-Job erneut, der Wake-Befehle ausführt. Ihre **Existenz** aktiviert Wake: Der Wert `false` oder eine leere Zeichenfolge deaktiviert es nicht. Der alte Name `WAILS_USE_WAKE` wird ignoriert.

## Anwendung erstellen

Starte den Assistenten zur Projekterstellung:

```sh
wails3 init -ui
```

Im Browser wird ein Projektassistent geöffnet. Wenn Wake aktiviert ist, erfasst er die Quell- und Build-Optionen in einem einzigen Ablauf. Bis zum abschließenden Klick auf **Projekt erstellen** wird nichts erstellt. `wails3 setup` ist ein separater Assistent für globale Einstellungen, Werkzeuge und Signierungsanmeldedaten; er konfiguriert nicht das Wake-Profil dieses Projekts.

| Ansicht | Auswahl für dieses Tutorial | Nächster Schritt |
| --- | --- | --- |
| **Projekt – Wails-Projekt erstellen** | Gib `wake-demo` ein und klicke dann auf **Weiter**. | Das Projekt wird in einem neuen Ordner namens `wake-demo` unterhalb des angezeigten Verzeichnisses erstellt. |
| **Framework – Framework auswählen** | Wähle **Vanilla** und klicke dann auf **Weiter**. | Du erhältst ein kleines Frontend, ohne dich zuerst in ein bestimmtes Framework einarbeiten zu müssen. |
| **Sprache – Spracheinstellung** | Wähle **JavaScript** und klicke dann auf **Weiter**. | Die ausschließlich für TypeScript vorgesehene Ansicht **Bindings** wird übersprungen. |
| **Details – Projektdetails** | Verwende den Produktnamen `Wake Demo`, die Bundle-ID `com.example.wakedemo` und die Version `0.1.0`. Klicke auf **Weiter**. | Deine Metadaten bleiben im Assistenten erhalten, während du die Build-Einstellungen auswählst. |

Verwende für diese Anleitung nicht **Globale Standardeinstellungen verwenden**, da dadurch die Auswahl von Framework und Sprache übersprungen wird. Wenn du in einem anderen Projekt TypeScript auswählst, kannst du in der Ansicht **Bindings** zwischen Schnittstellen und Klassen für die generierten Go-Datenbindungen wählen.

## Build-Einstellungen vor der Erstellung auswählen

Fahre im selben Browserfenster fort. Die folgenden Screenshots zeigen diesen Ablauf unter Linux. Unter Windows oder macOS unterscheiden sich dein Host-Ziel und die benötigten Werkzeuge.

### Anwendung – Was erstellst du?

Behalte **Desktop-Anwendung** bei und wähle **Standard-Backend**. Klicke auf **Weiter**.

![Anwendungsansicht im Assistenten zur Projekterstellung mit ausgewählter Desktop-Anwendung und ausgewähltem Standard-Backend.](/assets/screenshots/wake/application.png)

In dieser Anleitung erstellst du eine GUI mit einem Web-Frontend. **Pure Go (experimentell)** muss separat aktiviert werden. Für die unterstützte Desktop-Kompilierung entfällt damit die Anforderung eines C-Compilers, sofern sich auch deine Abhängigkeiten ohne CGo bauen lassen. Ziellaufzeitbibliotheken sowie sämtliche Werkzeuge für Paketierung und Signierung werden dadurch nicht überflüssig.

Das Backend ist eine **projektweite Einstellung**. Andere Profile, die diese Einstellungen erben, verwenden ebenfalls das ausgewählte Backend.

### Ziele – Wo soll die Anwendung ausgeführt werden?

Behalte nur das Ziel für den Computer bei, den du gerade verwendest, und klicke dann auf **Weiter**. Der Assistent wählt zunächst das Betriebssystem und die Architektur des Hosts aus.

![Zielansicht, in der für den Tutorial-Host nur linux/amd64 ausgewählt ist.](/assets/screenshots/wake/targets.png)

Zielnamen kombinieren ein Betriebssystem mit einer CPU-Architektur, beispielsweise `windows/amd64`, `darwin/arm64` für einen Apple-Silicon-Mac oder `linux/amd64`. `darwin` steht für macOS. Wenn du ein anderes Ziel auswählst, forderst du eine ausführbare Datei für diesen Computer an; dadurch kann sie nicht auf deinem aktuellen Computer ausgeführt werden.

Zusätzliche Architekturen, universelle macOS-Builds und experimentelle Mobilziele kannst du später hinzufügen. Ergänze sie nach dem ersten lokalen Build in `wails.yaml`.

### Pakete – Release-Ausgaben auswählen

Lasse die Kontrollkästchen für Paketformate deaktiviert und **Release-Ausgaben signieren** ausgeschaltet. Klicke auf **Weiter**.

![Paketansicht ohne zusätzlich ausgewählte Linux-Paketformate und mit deaktivierter Signierung.](/assets/screenshots/wake/packages.png)

Wenn keine zusätzlichen Formate ausgewählt sind, erzeugt Wake die Standardausgabe des Ziels: eine Desktop-Programmdatei oder für eine macOS-GUI ein macOS-App-Bundle. Das genügt, um die Anwendung lokal zu prüfen. Durch die Auswahl von Formaten werden, sofern unterstützt, Distributionsausgaben wie DMG, MSIX, NSIS, DEB oder AppImage hinzugefügt.

Für diesen ersten lokalen Build ist die Signierung optional. Ein erfolgreicher unsignierter Build bedeutet nicht, dass die Anwendung zur öffentlichen Verteilung bereit ist. Konfigurieren Sie für ein signiertes Projekt vor der Erstellung wiederverwendbare Anmeldedaten mit `wails3 setup` oder anschließend die Signierung in `wails.yaml`. Der Erstellungsassistent verwendet Ihre konfigurierten Anmeldedaten; Ihre globalen Signierungseinstellungen ändert er nicht.

Wake kann konfigurierte Zertifikatsdateien außerhalb des Projekts verwenden, ohne private Schlüssel oder Passwörter nach `wails.yaml` zu kopieren. Explizite Signierungseinstellungen des Projekts haben Vorrang. Befolgen Sie vor der Vorbereitung einer öffentlichen Veröffentlichung die Anleitung zum [Einrichten der Signierung](/getting-started/setup/) und den [Leitfaden zum Paketieren und Signieren](/experimental/yaml-builds/#package-and-sign-for-distribution). Die mobile Signierung wird in YAML konfiguriert.

### Tools — Werkzeuge für Ihren Build

Lesen Sie die Anforderungen für Ihre Auswahl. Die Liste stammt aus dem Build-Plan von Wake; sie fordert nicht dazu auf, jedes Werkzeug auf jeder Plattform zu installieren.

![Tools-Ansicht mit den ermittelten Anforderungen für den lokalen Linux-Build aus dem Quellcode.](/assets/screenshots/wake/tools.png)

Für diese Lektion benötigen Sie voraussichtlich Go und die Frontend-Werkzeuge sowie die nativen Entwicklungswerkzeuge für das standardmäßige GUI-Backend Ihres Hostsystems. Docker ist keine allgemeine Voraussetzung. Ein Cross-Build mit CGo kann einen Compiler und Zielbibliotheken oder eine unterstützte Container-Toolchain erfordern.

Beheben Sie Einträge mit dem Status **Nicht bereit**, **Konfiguration erforderlich** oder **Aktualisierung erforderlich**, bevor Sie den lokalen Build versuchen. Verwenden Sie die angezeigten Anweisungen, sofern verfügbar, oder den [Installationsleitfaden für die Plattform](/getting-started/installation/) und klicken Sie anschließend auf **Anforderungen erneut prüfen**. Wenn eine Installation den `PATH` des Terminals geändert hat, starten Sie `wails3 init -ui` aus einem aktualisierten Terminal neu und wiederholen Sie Ihre Auswahl.

Der Assistent kann auch dann ein gültiges Projekt erstellen, wenn auf diesem Computer Werkzeuge oder Anmeldedaten für die Signierung fehlen. Das ist für die CI-Konfiguration nützlich, doch für dieses Tutorial benötigen Sie einen Build, der hier ausgeführt werden kann. Eine ungültige Konfiguration muss vor der Erstellung korrigiert werden.

### Projekt erstellen

Prüfe deine Auswahl: Desktop-Anwendung, Standard-Backend, Zielsystem deines Hosts und deaktivierte Signierung. Klicke unter **Tools** auf **Projekt erstellen** und warte auf **Projektauswahl übernommen**. Kehre zum Terminal zurück und warte, bis die CLI die Projekterstellung abgeschlossen hat. Die Bestätigung im Browser übernimmt lediglich deine Auswahl.

![Bildschirm „Projektauswahl übernommen“ mit dem vorgesehenen Anwendungspfad und den Befehlen, die nach erfolgreicher Erstellung verwendet werden sollen.](/assets/screenshots/wake/save.png)

Der Assistent erstellt den Quellcode und speichert die Ziel-, Paket- und Signierungsauswahl unter `profiles.setup` in `wails.yaml`. Außerdem legt er den Anwendungstyp, die Frontend-Aktivierung und das projektweite Backend fest. Sie können den Browser schließen und zum Terminal zurückkehren. Warten Sie, bis die Projekterstellung abgeschlossen ist, und führen Sie dann Folgendes aus:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@b1daf95f3
```

Sie sollten nun über `main.go`, ein Go-Modul, ein Verzeichnis `frontend` und eine native `wails.yaml` verfügen. Die integrierte Vorlage verwendet Wake direkt; dieses neue Projekt müssen Sie nicht migrieren. Eine Community-Vorlage kann benutzerdefinierte Legacy-Tasks enthalten und stattdessen die unten beschriebene Migrationsprüfung erfordern.

Die Abschlussansicht zeigt `wails3 build setup` und `wails3 dev`. Der Profilname ist wichtig: `wails3 build setup` verwendet die soeben gespeicherte Auswahl. Einfaches `wails3 build` wählt das Profil `setup` nicht automatisch aus. Der Profilname bedeutet nicht, dass Sie `wails3 setup` ausführen müssen.

## Ihre erste Anwendung erstellen und ausführen

Führen Sie den vom Assistenten angezeigten Befehl aus:

```sh
wails3 build setup
```

Der erste Build installiert Frontend-Abhängigkeiten, generiert Bindings und kompiliert die Anwendung. Spätere Builds können unveränderte Arbeit wiederverwenden. Warten Sie auf **Build erfolgreich** und lesen Sie die von der CLI ausgegebenen Pfade. Das Standardausgabeverzeichnis ist `bin`; eine macOS-GUI wird zu einem `.app`-Bundle zusammengestellt.

Starten Sie den Entwicklungsmodus auf Ihrem eigenen Computer:

```sh
wails3 dev
```

Nun sollte das Fenster der Starteranwendung angezeigt werden. Geben Sie Ihren Namen ein und probieren Sie die Begrüßung aus. Ändern Sie den Platzhalter des Namensfelds in `frontend/index.html` und speichern Sie die Datei; das Frontend sollte aktualisiert werden. Bearbeiten Sie die Begrüßung in `greetservice.go`, um einen erneuten Go-Build auszulösen. Drücken Sie abschließend im Terminal **Ctrl+C**.

Die Entwicklung ist ein lokaler Host-Workflow. Ein Release-Profil kann mehrere Cross-Build-Ziele enthalten, ohne zu ändern, auf welchem Computer das Entwicklungsfenster ausgeführt werden kann. Verwenden Sie `wails3 run`, um eine Produktionsanwendung lokal zu erstellen und zu starten; siehe die [Referenz zu run](/reference/wails-yaml-run/).

Sie verfügen nun über ein Quellprojekt, ein geprüftes Build-Profil und einen funktionierenden lokalen Entwicklungszyklus. Committen Sie den Quellcode und `wails.yaml` in die Versionsverwaltung. Generierte `.wails/`-Arbeitsdateien sind nicht die Quelle Ihrer Konfiguration.

## Nächste Schritte

| Ziel | Wake-Unterstützung | Zuerst zu prüfen |
| --- | --- | --- |
| Für einen anderen Desktop erstellen | Ziele für Windows, macOS und Linux, einschließlich beider macOS-Architekturen und universeller Ausgabe. | Pure Go und CGo stellen unterschiedliche Compileranforderungen. Testen Sie auf dem Zielbetriebssystem und der Zielarchitektur. |
| Installationsprogramme oder Pakete erstellen | MSIX, NSIS, macOS-App/DMG, DEB, RPM, Arch, Void und AppImage. | NSIS benötigt `makensis`. AppImage benötigt eine Ziellaufzeit; die plattformübergreifende GUI-Paketierung benötigt außerdem einen bereitgestellten Linux-Laufzeitverzeichnisbaum. |
| Eine Veröffentlichung signieren | Portable Signierung mit exportierten Schlüsseln für Windows, macOS und DEB/RPM/Arch; native Identitätsspeicher bleiben verfügbar. | Jede Identität, jedes Format und jeder Host unterliegt Einschränkungen. AppImage und Void besitzen keinen integrierten Signierungsablauf. |
| Eine macOS-App notarisieren | Übermittlung mit einem Team-API-Schlüssel und Anheften des Tickets oder ein macOS-Schlüsselbundprofil. | Eine geeignete Developer-ID-Signatur und ein Zeitstempel sind weiterhin erforderlich. Dass von Go erzeugte Testsignaturen native Prüfungen bestehen, beweist nicht, dass Apple Ihre Veröffentlichung akzeptiert. |
| Mehrere Programme erstellen | Ein primäres GUI-, CLI- oder Serverprogramm sowie benannte Begleitprogramme. | Deklarieren Sie Begleitpakete und -typen in YAML. Der Assistent generiert deren Quellcode nicht. |
| Andernorts erstellte Dateien paketieren | `wails3 build setup --skip-build`. | Geben Sie in `build.output` übereinstimmende Namen und Architekturen der Programmdateien an. Dieses Flag ist eine einmalige Befehlsoption und keine dauerhafte Profileinstellung. |
| Release-Automatisierung hinzufügen | Benannte Profile, begrenzte Lebenszyklus-Hooks, Caching und parallele Paketverarbeitung. | Verwenden Sie literales YAML und deklarierte, projekteigene Skripte. Beliebiges Taskfile-Verhalten ist kein Erweiterungsmechanismus für Wake. |

Bearbeiten Sie `profiles.setup`, andere benannte Profile und erweiterte Einstellungen in `wails.yaml`. Verwenden Sie `wails3 setup`, um gemeinsame Einstellungen und Anmeldedaten zu ändern. Vollständige Beispiele und Paketanforderungen finden Sie im [Leitfaden für YAML-Builds](/experimental/yaml-builds/). GUI-Anwendungen benötigen weiterhin die Systembibliotheken des Zielsystems; Pure Go macht sie nicht zu vollständig statischen Anwendungen.

iOS und Android sind weiterhin experimentell, benötigen ihre jeweiligen Plattform-SDKs sowie CGo und liegen außerhalb dieser ersten Desktop-Lektion. Der Assistent macht mobile Builds nicht zwischen beliebigen Hosts portierbar.

## Ein vorhandenes Projekt migrieren

Arbeiten Sie in einem separaten Branch oder einer Kopie des Arbeitsbaums eines **Wails-v3**-Projekts. Dies ist eine Migration der Build-Konfiguration und kein Upgrade von Wails v2. Bewahren Sie Ihr vorhandenes Build-Ergebnis und Ihre Tests auf, damit Sie das Verhalten nach der Aktivierung vergleichen können.

@note{type="note" title="Mit einer geprüften Migration beginnen"}
Neue Projekte, die aus den aktuellen integrierten Legacy-Vorlagen generiert wurden, können migriert werden, bevor ihr Frontend gebaut wurde. Wake erkennt deren unveränderte Task-Selektoren und Entwicklungslebenszyklus, ohne Taskfiles auszuführen oder benutzerdefinierte Ausdrücke auszuwerten.

Benutzerdefinierte Tasks, Variablenüberschreibungen und Entwicklungsrichtlinien müssen weiterhin geprüft werden. Wenn die Analyse Ihr Projekt nicht abbilden kann, belassen Sie es beim Legacy-Build und melden Sie die Diagnoseinformationen. Entfernen Sie kein benötigtes Verhalten, nur damit die Aktivierung erfolgreich ist. Das obige Tutorial für neue Projekte erfordert keine Migration.
@end

### 1. Analysieren, ohne den aktiven Build zu ändern

Aktivieren Sie Wake und verwenden Sie die Vorschau-CLI vom Anfang dieser Seite. Führen Sie im Stammverzeichnis des vorhandenen Projekts Folgendes aus:

```sh
wails3 migrate --dry-run
```

Prüfen Sie die gemeldeten unterstützten Einstellungen und Blocker. Der Analysator liest Legacy-Deklarationen, führt deren Befehle jedoch nicht aus. Ein Probelauf erstellt kein aktives Manifest. Verwenden Sie für maschinenlesbare Diagnoseinformationen `wails3 migrate --json`. Wenn außerdem kein Entwurf geschrieben werden soll, verwenden Sie `--dry-run --json`.

### 2. Inaktiven Vorschlag erstellen und prüfen

```sh
wails3 migrate
```

Dadurch wird **`wails.migrated.yaml`** erstellt, nicht `wails.yaml`. Prüfen Sie die Datei zusammen mit Ihren vorhandenen Taskfiles und `build/config.yml`. Kontrollieren Sie Projektmetadaten, Frontend-Befehle, Zielarchitekturen, Paketinhalte, Assets und Signierungseinstellungen. Bewahren Sie Geheimnisse in Umgebungsvariablen oder Ihrem Anmeldedatenspeicher auf.

Bekannte integrierte Tasks können nativen Einstellungen zugeordnet werden. Konservative Lebenszyklus-Tasks ohne Argumente, die ein projekteigenes Skript aufrufen, können zu Hooks werden. Anderes erreichbares Verhalten benutzerdefinierter Tasks ist ein Blocker; Wake verwirft es weder stillschweigend noch führt es ersatzweise über Task aus.

Wenn bereits ein Entwurf vorhanden ist, überschreibt die Migration Ihre geprüften Änderungen nicht. Bewahren Sie eine Kopie auf, bevor Sie ihn absichtlich ersetzen oder neu generieren. Benennen Sie einen blockierten Entwurf nicht in `wails.yaml` um, um die Aktivierungsprüfungen zu umgehen. Reduzieren oder ersetzen Sie nicht unterstütztes Legacy-Verhalten und führen Sie die Analyse vor der Aktivierung erneut aus. Wenn sich das Projekt noch nicht abbilden lässt, behalten Sie den Legacy-Build bei und melden Sie die Lücke.

### 3. Explizit aktivieren

Nachdem Sie den Entwurf geprüft und die Blocker behoben haben, führen Sie Folgendes aus:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

Der erste Befehl prüft die Aktivierung, ohne sie anzuwenden. Der zweite prüft die Legacy-Deklarationen erneut und validiert den Entwurf, bevor er `wails.yaml` atomar erstellt. Nach erfolgreicher Aktivierung wird die inaktive Datei `wails.migrated.yaml` entfernt. Taskfiles, `build/config.yml`, Skripte und Assets bleiben unverändert. Ab diesem Zeitpunkt steuert das aktive Manifest die Wake-Builds. Die alten Dateien sind keine zweite Konfigurationsquelle.

Einige ältere Projekte enthalten bereits eine reine Metadatendatei namens `wails.yaml` mit einem Abschnitt `info:` auf oberster Ebene. Obwohl sie kein natives Wake-Manifest ist, belegt sie dennoch das Aktivierungsziel. Bewahren Sie diese Legacy-Datei auf und verschieben Sie sie vor der Aktivierung bewusst. Das Migrationswerkzeug überschreibt oder benennt sie nicht für Sie um.

### 4. Migrierte Anwendung überprüfen

```sh
wails3 config check
wails3 build
wails3 dev
```

Vergleichen Sie die Anwendung und die Paketinhalte einschließlich Ihres benutzerdefinierten Release-Verhaltens mit dem vorherigen Build. Führen Sie Ihre Projekttests aus. Prüfen Sie jedes migrierte Profil explizit, beispielsweise `wails3 config check release` und `wails3 build release`, wenn der Entwurf ein Profil `release` enthält.

Bearbeiten Sie die migrierten Profile direkt in `wails.yaml`. `wails3 init -ui` erstellt neue Projekte; `wails3 setup` verwaltet globale Einstellungen. Keines von beiden dient zum Bearbeiten des Release-Profils eines vorhandenen Projekts.

### Wenn Sie das Experiment beenden müssen

Wenn Sie `WAILS_EXP_USE_WAKE` entfernen, wird die explizite Aktivierung deaktiviert. Ein aktives `wails.yaml` blockiert jedoch weiterhin die automatische Weiterleitung an Taskfiles. Ungültiges YAML führt ebenfalls zum Abbruch mit einer Fehlermeldung. Es gibt keinen Befehl für eine Rückmigration.

Um zum alten Workflow zurückzukehren, stellen Sie in der Versionsverwaltung die vollständige Projektversion vor der Migration wieder her, einschließlich des fehlenden aktiven Manifests, und verwenden Sie die dazu kompatible CLI und die entsprechenden Abhängigkeiten. Sichern Sie zuvor alle neuen Quellcodeänderungen, die Sie benötigen. Ein separater Test-Branch erleichtert diese bewusste Wiederherstellung erheblich.

## Problem diagnostizieren und Feedback geben

| Symptom | Erste Prüfung |
| --- | --- |
| Bei der Projekterstellung werden Anwendung, Ziele, Pakete und Werkzeuge übersprungen | Setzen Sie `WAILS_EXP_USE_WAKE` in dem Terminal, aus dem Sie den Vorgang starten, und prüfen Sie, welche ausführbare Datei `wails3` über `PATH` gefunden wird. Das globale `wails3 setup` behält immer seinen Ablauf für Einstellungen bei. |
| Eine Community-Vorlage muss migriert werden | Prüfen Sie vor der Aktivierung ihre Legacy-Konfiguration mit den Migrationsbefehlen. |
| Ein gespeicherter Build meldet fehlende Werkzeuge | Beim Erstellen eines Projekts wird die beabsichtigte Konfiguration gespeichert. Installieren Sie vor dem Build die gemeldeten Anforderungen für das Hostsystem. Bei der Erstellung werden nicht alle ausgewählten Werkzeuge installiert. |
| Eine für ein anderes Ziel erstellte ausführbare Datei startet auf diesem System nicht | Führen Sie sie auf dem ausgewählten Betriebssystem und der ausgewählten Architektur aus. Verwenden Sie für die lokale Entwicklung das Host-Ziel. |
| Im Prebuilt-Modus wird keine Binärdatei gefunden | Prüfen Sie die [Dateinamen und Pfade der vorgefertigten Binärdateien](/experimental/yaml-builds/#use-prebuilt-desktop-executables) einschließlich der Unterverzeichnisse für mehrere Ziele. |
| Eine Manifeständerung schlägt fehl | Führen Sie `wails3 config check` für das betreffende Profil aus und befolgen Sie die Diagnosehinweise für das Feld. Wake greift nicht ersatzweise auf Taskfiles zurück. |

Prüfen Sie einen Build vor der Ausführung mit `wails3 build setup --plan --json`. `wails3 config check` validiert die Konfiguration, und `wails3 eject` schreibt eine inaktive Momentaufnahme `wails.ejected.yaml`, ohne das aktive Manifest zu ändern. Prüfen Sie die Befehlsausgabe, bevor Sie sie weitergeben: Projektpfade und Speicherorte von Anmeldedaten können Rückschlüsse auf Ihre lokale Umgebung ermöglichen.

Geben Sie Feedback in den [Diskussionen zu Experimenten](https://github.com/wailsapp/wails/discussions/categories/experiments) oder [melden Sie einen reproduzierbaren Fehler](https://github.com/wailsapp/wails/issues). Geben Sie die Versionen Ihrer CLI und Ihres Frameworks, Betriebssystem und Architektur des Hostsystems, das ausgewählte Ziel, den genauen Befehl, die Auswahl von Pure Go sowie die relevanten Diagnoseinformationen an. Nennen Sie, ob der Fehler beim Kompilieren, Paketieren, Signieren oder Ausführen der Anwendung auftrat. Fügen Sie bei Bedarf ein minimales Manifest ohne Geheimnisse sowie die Ausgabe von `wails3 doctor` bei.
