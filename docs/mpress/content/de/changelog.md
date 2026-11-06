---
title: "Änderungsprotokoll"
description: "Versionsverlauf und Versionshinweise für Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Legende:

-  – macOS
- ⊞ – Windows
- 🐧 – Linux

/_-- Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

Das Format basiert auf [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), und dieses Projekt folgt der [semantischen Versionierung](https://semver.org/spec/v2.0.0.html).

- `Added` für neue Funktionen.
- `Changed` für Änderungen an bestehenden Funktionen.
- `Deprecated` für Funktionen, die demnächst entfernt werden.
- `Removed` für inzwischen entfernte Funktionen.
- `Fixed` für sämtliche Fehlerbehebungen.
- `Security` bei Sicherheitslücken.

_/

/_   * BITTE AKTUALISIEREN SIE DIESE DATEI NICHT *   Ergänzen Sie Aktualisierungen in `v3/UNRELEASED_CHANGELOG.md`   Vielen Dank! _/

## [Unveröffentlicht]

## v3.0.0-beta.21 - 2026-09-13

## Hinzugefügt

- Wails-v3-Dokumentation mit M-Press in [PR](https://github.com/wailsapp/wails/pull/6116) von @leaanthony bereitstellen

## Behoben

- JSON-Slug-Werte aus dem MPD-Frontmatter für die Changelog-Generierung parsen, in [PR](https://github.com/wailsapp/wails/pull/6118) von @leaanthony
- Der Updater löscht Hilfsumgebungsvariablen und startet nach fehlgeschlagenen Sicherungen das ursprüngliche Ziel neu, in [PR](https://github.com/wailsapp/wails/pull/6080) von @cnmax
- Standard-Signalhandler während App.Run starten, in [PR](https://github.com/wailsapp/wails/pull/6098) von @leaanthony
- Das Windows-Menü verarbeitet nil-Menüs, gibt ersetzte Ressourcen frei und zeichnet die Menüleiste neu, in [PR](https://github.com/wailsapp/wails/pull/6112) von @taliesin-ai
- MSIX-Paketierung für neue Projekte mit gemeinsamer YAML-Konfiguration wiederherstellen, in [PR](https://github.com/wailsapp/wails/pull/6115) von @leaanthony
- Beheben, dass generierte JavaScript- und TypeScript-Bindings nicht geladen werden, wenn Ersteller generischer Modelle auf später deklarierte Hilfsfunktionen verweisen, und Stacküberläufe beim Erstellen voneinander abhängiger generischer Modelle verhindern (#6062)

## v3.0.0-beta.20 – 2026-09-10

## Geändert

- Links zum Clave-Showcase auf die aktuelle Website und das aktuelle Repository aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/6082) von @01xR4in

## Behoben

- Abgebrochene Windows-Asset-Anfragen einschließlich Worker-Anfragen werden beendet, während Keepalive-Handler über Navigationen hinweg erhalten bleiben. Native Anfragekontexte werden auf Apple-Plattformen über den Anwendungs-Wrapper weitergeleitet. (#5963, #5969)
- Changelog-Einträge bleiben bei konkurrierenden Pushes durch Wiederholungsversuche erhalten, in [PR](https://github.com/wailsapp/wails/pull/6094) von @leaanthony
- Fehler behoben, durch den `go mod vendor` auf jeder Plattform mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug: Die Einbettungen, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen, wurden entfernt. Dadurch wurden [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) behoben, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002

## Entfernt

- Native WebView2-Loader-Unterstützung entfernt, die durch den reinen Go-Loader ersetzt wurde. Dadurch entfallen die eingebetteten `WebView2Loader.dll`-Binärdateien und die Abhängigkeit `github.com/jchv/go-winloader`. Das Build-Tag `native_webview2loader` wird weiterhin akzeptiert und verursacht keinen Fehler mehr, hat jedoch keine Auswirkung auf v3-Builds, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Nicht verwendete Build-Tags und die FPS-Option aus dem macOS-API-Leitfaden entfernt, in [PR](https://github.com/wailsapp/wails/pull/6097) von @leaanthony

## v3.0.0-beta.19 – 2026-09-09

## Hinzugefügt

- Private macOS-APIs durch Build-Tags abgesichert, sodass ihre Verwendung explizit aktiviert werden muss – siehe [Dokumentation](https://v3.wails.io/features/browser/integration) und [Dokumentation](https://v3.wails.io/features/environment/info) und [Dokumentation](https://v3.wails.io/features/windows/basics) und [Dokumentation](https://v3.wails.io/features/windows/frameless) und [Dokumentation](https://v3.wails.io/features/windows/notch-windows) und [Dokumentation](https://v3.wails.io/features/windows/options) und [Dokumentation](https://v3.wails.io/guides/build/macos) und [Dokumentation](https://v3.wails.io/guides/build/private-macos-apis) und [Dokumentation](https://v3.wails.io/reference/overview), in [PR](https://github.com/wailsapp/wails/pull/6087) von @leaanthony

## Behoben

- Runtime-Anfragen über 64 MiB werden mit HTTP 413 abgelehnt, in [PR](https://github.com/wailsapp/wails/pull/6091) von @leaanthony

## Sicherheit

- MCP-Ursprünge und Remotezugriff durch Token-Authentifizierung abgesichert, in [PR](https://github.com/wailsapp/wails/pull/6092) von @leaanthony

## v3.0.0-beta.18 – 2026-09-08

## Behoben

- Calloc-Speicherleck unter Linux und Darwin durch die Verwendung von Pointer-Receivern behoben, in [PR](https://github.com/wailsapp/wails/pull/6083) von @4RH1T3CT0R7

## v3.0.0-beta.17 – 2026-09-06

## Behoben

- Windows: Ein fehlgeschlagenes oder nil-wertiges `GetRequest` im WebResourceRequested-Handler beendet den Prozess nicht mehr (`log.Fatal` / Panic durch nil-Dereferenzierung). Stattdessen wird die Anfrage verworfen und protokolliert, in [PR](https://github.com/wailsapp/wails/pull/6006) von @midagedev

## v3.0.0-beta.16 – 2026-08-29

## Geändert

- Das Notarisierungspasswort wird in einem neuen Terminalfenster abgefragt, in [PR](https://github.com/wailsapp/wails/pull/6029) von @leaanthony

## Behoben

- Klicktypen des Infobereichsymbols unter macOS werden korrekt verarbeitet, in [PR](https://github.com/wailsapp/wails/pull/5919) von @ChewbaccaCookie
- Die CI entfernt vor der Aktualisierung nicht verwendete Microsoft-APT-Repositorys, in [PR](https://github.com/wailsapp/wails/pull/6041) von @Grantmartin2002

## v3.0.0-beta.15 – 2026-08-27

## Behoben

- Zeitüberschreitung für die WebView2-Einbettung auf 60 Sekunden erhöht, in [PR](https://github.com/wailsapp/wails/pull/6043) von @Grantmartin2002

## v3.0.0-beta.14 – 2026-08-26

## Behoben

- Tastendrücke mit Strg-Buchstaben werden unter macOS korrekt benannt, in [PR](https://github.com/wailsapp/wails/pull/6032) von @taliesin-ai
- ICO-Symbole im Infobereich korrigiert und an das Taskleisten-Design unter Windows angepasst, in [PR](https://github.com/wailsapp/wails/pull/6016) von @nik9play

## v3.0.0-beta.13 – 2026-08-25

## Behoben

- Die Verarbeitung von Aufgaben im Hauptthread unter macOS wird auch während einer modalen Ereignisschleife fortgesetzt, in [PR](https://github.com/wailsapp/wails/pull/6026) von @leaanthony
- Sicherer Speicher auf Mobilgeräten kann nun fehlschlagen und verweigert den Zugriff bei einem Fehler, in [PR](https://github.com/wailsapp/wails/pull/5923) von @mortenolsrud
- Ereignis-Hooks der Anwendung werden auch ausgeführt, wenn kein Listener registriert ist, in [PR](https://github.com/wailsapp/wails/pull/5999) von @archy-rock3t-cloud
- Tippfehler in Kommentaren und lokalisierter Dokumentation korrigiert, in [PR](https://github.com/wailsapp/wails/pull/6023) von @haoku123
- Die unter `v3/examples` eingecheckten vorkompilierten macOS-Binärdateien entfernt, in [PR](https://github.com/wailsapp/wails/pull/6025) von @4RH1T3CT0R7

## v3.0.0-beta.12 – 2026-08-21

## Hinzugefügt

- macOS-Benachrichtigungsfenster für die Displayaussparung mit Lebenszyklus- und Telemetriebeispiel hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/windows/notch-windows), in [PR](https://github.com/wailsapp/wails/pull/6010) von @leaanthony
- Unterstützung für macOS-NSPanel-Fenster mit neuen Optionen und nativer Integration hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/windows/options) in [PR](https://github.com/wailsapp/wails/pull/6008) von @leaanthony

## Behoben

- Verhindert, dass SQLite Prepare bei gleichzeitigen Aufrufen hängen bleibt, in [PR](https://github.com/wailsapp/wails/pull/5998) von @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Entfernt

- Veralteten Implementierungs-Tracker aus der Dokumentation entfernt, in [PR](https://github.com/wailsapp/wails/pull/6005) von @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Behoben

- Behoben, dass der GTK4-Linux-Host beim Start über benutzerdefinierte Protokolle und Dateizuordnungen übergebene Argumente verwirft, in [PR](https://github.com/wailsapp/wails/pull/6000) von @midagedev
- Gelöschte Zeilen und Korrekturen aus derselben Quelle werden bei der Changelog-Validierung nun korrekt verarbeitet, in [PR](https://github.com/wailsapp/wails/pull/5993) von @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Hinzugefügt

- Sicheren wails3-MCP-Server für agentengestützte Projektverwaltung hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5896) von @leaanthony
- Dokumentation für Modelle in Bindings hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/bindings/models) in [PR](https://github.com/wailsapp/wails/pull/5988) von @taliesin-ai
- Unterstützung für die Installation mit rpm-ostree auf atomaren Linux-Systemen hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5987) von @leaanthony
- Native wöchentliche Generierung und Veröffentlichung von Diagrammen zum Sternverlauf hinzugefügt – siehe [Dokumentation](https://v3.wails.io/credits), [Dokumentation](https://v3.wails.io/de/credits), [Dokumentation](https://v3.wails.io/fr/credits), [Dokumentation](https://v3.wails.io/id/credits), [Dokumentation](https://v3.wails.io/ja/credits), [Dokumentation](https://v3.wails.io/ko/credits), [Dokumentation](https://v3.wails.io/pt/credits), [Dokumentation](https://v3.wails.io/ru/credits), [Dokumentation](https://v3.wails.io/zh-cn/credits) und [Dokumentation](https://v3.wails.io/zh-tw/credits) in [PR](https://github.com/wailsapp/wails/pull/5986) von @leaanthony
- Nur für Darwin verfügbares mac-Paket zum Auflösen von Ressourcen aus Anwendungspaketen hinzugefügt – siehe [Dokumentation](https://v3.wails.io/guides/build/macos) in [PR](https://github.com/wailsapp/wails/pull/5965) von @leaanthony
- Showcase-Seite und Indexeintrag für Condui hinzugefügt – siehe [Dokumentation](https://v3.wails.io/community/showcase/condui) und [Dokumentation](https://v3.wails.io/community/showcase) in [PR](https://github.com/wailsapp/wails/pull/5962) von @mgueregath
- Showcase-Seite für Redis Viewer mit Screenshots und Projektlink hinzugefügt – siehe [Dokumentation](https://v3.wails.io/community/showcase) und [Dokumentation](https://v3.wails.io/community/showcase/redisviewer) in [PR](https://github.com/wailsapp/wails/pull/5984) von @redisviewer

## Geändert

- GTK-Anwendungs-Flags unter Linux auf G<em>APPLICATION</em>NON_UNIQUE aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5971) von @overlordtm
- Fehlende Fensterereignisse werden nun auf Debug- statt auf Warnstufe protokolliert, in [PR](https://github.com/wailsapp/wails/pull/5914) von @julianstorer

## Behoben

- Registrierte macOS-Tastenkürzel erhalten nun Vorrang vor der Webview, in [PR](https://github.com/wailsapp/wails/pull/5902) von @julianstorer
- Defekte Links in der Seitenleiste der Dokumentation repariert, in [PR](https://github.com/wailsapp/wails/pull/5937) von @northes
- Kontexte von Asset-Anfragen unter macOS und iOS werden abgebrochen, wenn WebKit die zugehörige Aufgabe des benutzerdefinierten Schemas abbricht (#5963)
- Nachricht WindowSetFullscreenButtonEnabled wird nun verarbeitet, in [PR](https://github.com/wailsapp/wails/pull/5976) von @archy-rock3t-cloud
- Fragment im preact-ts-Template importiert, um den Build-Fehler zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5979) von @haoku123
- Verhindert Abstürze älterer reiner GTK3-Dienstanwendungen, wenn die Bildschirmerkennung ausgeführt wird, bevor ein aktives Fenster oder Display verfügbar ist (#5966)
- Release-Läufe mit expliziter Version können nun fortgesetzt werden, wenn der Changelog für unveröffentlichte Änderungen leer ist (#5977)

## Sicherheit

- nanoid-Lockfiles der Website auf die korrigierte Version 3.3.18 aktualisiert, um Sicherheitshinweise zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5985) von @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Hinzugefügt

- Generierung von Dokumentations-URLs für automatische Changelog-Einträge hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5957) von @taliesin-ai
- Streams hinzugefügt: bidirektionale Byte-Streams zwischen Go und JavaScript mit dem WebSocket-Programmiermodell und ohne lauschenden Socket. Deklarieren Sie einen Stream in Go mit `app.HandleStream(name, handler)` und stellen Sie aus dem Frontend mit `Stream(name)` eine Verbindung her; dies gibt ein Objekt in der Form von `WebSocket` zurück. Go→JS wird über den Asset-Server durch eine gehaltene Poll-Anfrage pro Fenster übertragen, JS→Go durch einen normalen POST. Dabei wird weder ein TCP-Port gebunden noch etwas über `evaluateJavaScript` geleitet. In Server-Builds (`-tags server`) wird derselbe Handler stattdessen über einen echten WebSocket bereitgestellt, sodass der Anwendungscode in allen Builds identisch ist. Von @leaanthony
- Mailbox-Changelog-Eintrag nach „Unveröffentlicht“ verschoben, in [PR](https://github.com/wailsapp/wails/pull/5935) von @leaanthony

## Geändert

- Automatische Generierung der Dokumentationsseitenleiste und Ableitung des Blogautoren-Typs aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5938) von @leaanthony

## Behoben

- Die WebView2-Initialisierung verwendet nun eine Frist und eine Nachrichtenschleife, in [PR](https://github.com/wailsapp/wails/pull/5952) von @leaanthony
- Der WebView2-Cookie-Test wird in CI übersprungen, sofern er nicht ausdrücklich aktiviert wurde, und seine Ausführung wird an den aktuellen Betriebssystem-Thread gebunden, in [PR](https://github.com/wailsapp/wails/pull/5951) von @leaanthony
- Windows-Menü-Builder stellen Befehls-IDs für übergeordnete Untermenüelemente wieder her, in [PR](https://github.com/wailsapp/wails/pull/5944) von @gilad-ch
- Offizielles Cross-Compilation-Image an die GTK-4.14+-Baseline für die Linux-Unterstützung angeglichen (#5928)
- iOS-Xcode-Projekt so konfiguriert, dass geerbte Linker-Flags beibehalten werden, und -ObjC hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5915) von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im `wails3 dev`-Asset-Proxy bei großen Frontends behoben, der die kurzlebigen Ports des Hosts erschöpfen und dazu führen konnte, dass nicht zugehörige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Fensterspezifisches Ereignis-JavaScript wird für geordnete Zustellung und Rückstaukontrolle in eine Warteschlange gestellt, in [PR](https://github.com/wailsapp/wails/pull/5934) von @leaanthony

## Entfernt

- Release-Pipeline für Desktop-Binärdateien entfernt: v3-Releases bestehen nur aus Tags, und die `wails3`-CLI wird mit `go install` installiert. Entfernt `release-v3.yml` sowie den Nightly-Schritt, der sie ausgelöst hat, in [PR](https://github.com/wailsapp/wails/pull/5946) von @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Hinzugefügt

- macOS-Einstellung für automatische Wiedergabe hinzugefügt, mit der sich die Anforderung einer Benutzeraktion für die Medienwiedergabe deaktivieren lässt, in [PR](https://github.com/wailsapp/wails/pull/5512) von @Eyalm321
- Mailbox-Changelog-Eintrag nach „Unveröffentlicht“ verschoben in [PR](https://github.com/wailsapp/wails/pull/5935) von @leaanthony

## Geändert

- Die macOS-Zoomanimation verwendet CADisplayLink oder NSTimer für eine flüssigere Darstellung in [PR](https://github.com/wailsapp/wails/pull/5945) von @savely-krasovsky

## Behoben

- iOS-Xcode-Projekt so konfiguriert, dass geerbte Linker-Flags beibehalten werden, und -ObjC hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5915) von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im Asset-Proxy `wails3 dev` bei großen Frontends behoben, der die kurzlebigen Ports des Hosts erschöpfen und dazu führen konnte, dass nicht zugehörige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Fensterspezifisches Ereignis-JavaScript für geordnete Zustellung und Rückstaukontrolle in eine Warteschlange gestellt in [PR](https://github.com/wailsapp/wails/pull/5934) von @leaanthony

### Hinzugefügt

- Generische asynchrone FIFO-Mailbox für die geordnete Ereigniszustellung implementiert in [PR](https://github.com/wailsapp/wails/pull/5851) von @savely-krasovsky und @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Hinzugefügt

- Begrenzten hostseitigen Speicher für übergroße Ereignisse und geordnete JavaScript-Zustellung implementiert in [PR](https://github.com/wailsapp/wails/pull/5930) von @leaanthony
- Springen des macOS-Dock-Symbols als Umsetzung des Fensterblinkens implementiert in [PR](https://github.com/wailsapp/wails/pull/5921) von @julianstorer

## Behoben

- Der Asset-Server bewahrt beim Flushen Fehler der Inhaltstyperkennung und noch nicht geschriebene Präfixe auf in [PR](https://github.com/wailsapp/wails/pull/5931) von @leaanthony
- Absturz von macOS-Anwendungen beim Ersetzen des Anwendungsmenüs aus einem Wails-Callback verhindert
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) behoben. Die uxtheme-Exporte für den Dunkelmodus wurden erst ab Build 18334 verwendet, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel dargestellt, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, wodurch dunkler Text auf dunklem Hintergrund erschien. Die Ordinalwerte sind bereits ab 17763 vorhanden; die Mindestversion wurde entsprechend angepasst.
- Behoben, dass `w32.GetStockObject` statt `GetStockObject` die Funktion `GetDeviceCaps` aufrief, wodurch für jedes Standardobjekt 0 zurückgegeben wurde.
- Fehlerbehandlung und -meldung beim Herunterladen des WebView2-Bootstrappers verbessert in [PR](https://github.com/wailsapp/wails/pull/5924) von @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Behoben

- Die Aktivierung von macOS-Apps berücksichtigt die Aktivierungsrichtlinie jetzt nur für reguläre Apps in [PR](https://github.com/wailsapp/wails/pull/5897) von @julianstorer
- Nicht initialisierte GTK-Fenster in Linux-Builds abgesichert in [PR](https://github.com/wailsapp/wails/pull/5898) von @julianstorer
- Explizite deckende Hintergrundfarbe für Linux-WebKit-Fenster vor dem Laden der URL festgelegt in [PR](https://github.com/wailsapp/wails/pull/5899) von @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Geändert

- Android-Build-Tasks verwenden standardmäßig arm64, und deploy-emulator wählt die Hostarchitektur aus in [PR](https://github.com/wailsapp/wails/pull/5890) von @mortenolsrud

## Behoben

- Zoomzustand von macOS-Fenstern beim Ziehen beibehalten und Bewegung reduziert in [PR](https://github.com/wailsapp/wails/pull/5900) von @leaanthony
- Build im Windows-Servermodus durch Hinzufügen von `!server` zur Build-Einschränkung der Datei `webview_window_windows_nonclient.go` behoben

## v3.0.0-beta.3 - 2026-08-03

## Hinzugefügt

- Abschluss der Betaverifizierung für Phase 10 in den Implementierungsdetails dokumentiert in [PR](https://github.com/wailsapp/wails/pull/5881) von @leaanthony

## Behoben

- Fensterhandle an die Windows-Dunkelmodus-API übergeben und Argumente validiert in [PR](https://github.com/wailsapp/wails/pull/5877) von @leaanthony
- Ermittlung des Zustands der macOS-Titelleistenschaltflächen für rahmenlose Fenster zentralisiert in [PR](https://github.com/wailsapp/wails/pull/5870) von @taliesin-ai
- Unlesbaren Text in nativen Menüs verhindert, wenn eine Windows-Anwendung den Dunkelmodus anfordert, während das Windows-App-Design hell ist. Das Menü verwendet nun den passenden hellen nativen Hintergrund, bis Windows dunklen Menütext darstellen kann.
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) behoben. Die uxtheme-Exporte für den Dunkelmodus wurden erst ab Build 18334 verwendet, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel dargestellt, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, wodurch dunkler Text auf dunklem Hintergrund erschien. Die Ordinalwerte sind ab 17763 vorhanden, daher entspricht die Mindestversion nun diesem Stand.

## v3.0.0-beta.2 - 2026-08-02

## Geändert

- v3 von Alpha auf Beta hochgestuft
- Intelligente Standardwerte und das automatische Ausblenden von Pop-ups im Infobereich dokumentiert, einschließlich Regressionstests für die Auswahl des Klick-Handlers (#5840).
- Der GitHub-Updater schließt Windows-Installer-Assets standardmäßig aus in [PR](https://github.com/wailsapp/wails/pull/5861) von @leaanthony
- Unterstützung für abgerundete und eckige Ecken sowie Ecken mit benutzerdefiniertem Radius bei rahmenlosen macOS-Fenstern hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Aktuelle GTK4-Fenstergrößen gemeldet und Ereignisse für Größenänderung, Maximierung, Minimierung und Vollbildstatus von der konfigurierten Oberfläche ausgelöst (#5830).
- Absturz von Linux-WebKit beim Senden von Blob oder FormData in fetch-Anfragen behoben in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai
- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Hinzugefügt

## Geändert

- Unterstützung für abgerundete und eckige Ecken sowie Ecken mit benutzerdefiniertem Radius bei rahmenlosen macOS-Fenstern hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Hinzugefügt

- Unterstützung für die macOS-DMG-Paketierung mit neuen Optionen und Build-Tasks hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5857) von @leaanthony

## Geändert

- Der GitHub-Updater schließt Windows-Installer-Assets standardmäßig aus in [PR](https://github.com/wailsapp/wails/pull/5861) von @leaanthony

## Behoben

- Linux-WebKit-Absturz beim Senden von Blob oder FormData in fetch-Anfragen behoben in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Hinzugefügt

- Doppelklickaktionen auf die macOS-Titelleiste zum Maximieren oder Minimieren implementiert in [PR](https://github.com/wailsapp/wails/pull/5853) von @taliesin-ai

## Geändert

- Tutorial zum QR-Dienst auf NewServiceWithOptions umgestellt und Abstände hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5849) von @jeongkyu

## Behoben

- Reaktionsfähigkeit von WKWebView beim Zoomen unter macOS aufrechterhalten in [PR](https://github.com/wailsapp/wails/pull/5856) von @leaanthony
- GTK4-Abfragen der Fenstergröße korrigiert und Ereignisse für Größenänderungen sowie Zustandsereignisse für Maximierung, Minimierung und Vollbildmodus vom konfigurierten `GdkSurface` ausgelöst.

## v3.0.0-alpha2.119 - 2026-07-27

## Behoben

- Mehrsprachige Dokumentation um Architekturdiagramme ergänzt in [PR](https://github.com/wailsapp/wails/pull/5833) von @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Hinzugefügt

- Standardpfade für Ein- und Ausgaben der Symbolgenerierung bereitgestellt in [PR](https://github.com/wailsapp/wails/pull/5825) von @taliesin-ai
- Quell-Einstiegsmodule zum Feld sideEffects in der Datei package.json des Runtime-Pakets hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5797) von @savely-krasovsky
- Abschnitt zu Lizenz und Herkunft zum Leitfaden für Mitwirkende hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5816) von @taliesin-ai

## Behoben

- GTK4-CSS mit begrenztem Geltungsbereich für rahmenlose Fenster angewendet, um den Eckenradius zu entfernen, in [PR](https://github.com/wailsapp/wails/pull/5800) von @savely-krasovsky
- Fehler beim Ermitteln der Cursorposition unter Windows für Popup-Menüs und die Bildschirmauflistung ordnungsgemäß behandelt in [PR](https://github.com/wailsapp/wails/pull/5789) von @wayneforrest
- Der macOS-Dialog zum Öffnen von Dateien filtert Erweiterungen korrekt und validiert zulässige Dateien anhand ihres Suffixes in [PR](https://github.com/wailsapp/wails/pull/5678) von @phergul
- Initialisierung des Dunkelmodus unter Windows gegen nil-API-Aufrufe abgesichert in [PR](https://github.com/wailsapp/wails/pull/5793) von @roachadam
- Fehler beim 32-Bit-Build des Updaters behoben: Die Konstante `maxArchiveTotalSize` (2 GiB) überschritt beim Übergeben an `fmt.Errorf` auf `GOARCH=386` den Wertebereich des plattformspezifischen Typs `int`. Sie ist jetzt explizit als `int64` typisiert.
- nil-Zeiger-Panic beim Start behoben, wenn ein Fenster in Windows-Builds, welche die uxtheme-APIs für den Dunkelmodus nicht laden, eine dunkle (oder systembedingt dunkle) Titelleiste verwendet, etwa unter Windows 10 1809 / Windows Server 2019 (Build 17763). Die `AllowDarkModeForWindow`-Aufrufe beim Einrichten des Fensterdesigns sind jetzt wie bereits in `w32.SetMenuTheme` gegen nil abgesichert.

## v3.0.0-alpha2.117 - 2026-07-08

## Hinzugefügt

- Benutzerdefinierte Hit-Test-Logik für Nicht-Client-Bereiche unter Windows implementiert in [PR](https://github.com/wailsapp/wails/pull/5462) von @savely-krasovsky

## Geändert

- WebView2-Erkennung der Monitorskalierung anhand von UseVisualHosting konfiguriert in [PR](https://github.com/wailsapp/wails/pull/5761) von @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Hinzugefügt

- FAQ-Dokumentation auf Funktionen und Anleitungen für Wails v3 ausgerichtet in [PR](https://github.com/wailsapp/wails/pull/5763) von @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Behoben

- Behoben, dass `Menu.Update()` das native Menü unter GTK4 Linux nicht neu aufbaut (#5659, von @puneetdixit200 unabhängig diagnostiziert und in #5539 behoben)
- Absturz bei der Auflistung von macOS-Bildschirmen nach einer Anzeigeänderung behoben, indem ID- und Namenszeichenfolgen der Bildschirme kopiert und die Anzahl als Momentaufnahme erfasst wird (#5565, von @x-haose unabhängig diagnostiziert und in #5584 behoben)
- Absturz unter Windows behoben, wenn `WM_ERASEBKGND` während eines Minimierungs-/Wiederherstellungsübergangs einen einfarbigen Hintergrund zeichnet und `GetClientRect` dabei nil zurückgibt (Absicherung von @sinspired in #5636 gemeldet)
- Behoben, dass Fehler von Frontend-Bindings immer als Text geparst werden, von @mbaklor in #5690
- Windows-Buildfehler bei Verwendung des Build-Tags `server` behoben. Ursache war, dass in den Windows-GUI-Dateien die Build-Bedingung `!server` fehlte, die ihre macOS- und Linux-Entsprechungen bereits besitzen (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Hinzugefügt

- Update-Manifest-Protokoll und Endpoint-Provider implementiert in [PR](https://github.com/wailsapp/wails/pull/5720) von @taliesin-ai

## Geändert

- Das Binding `webview2` als `v3/internal/webview2` in das v3-Modul integriert und dabei das eigenständige Modul, dessen Workflows für nächtliche Releases und Synchronisierung sowie den Versionswechsel in go.mod entfernt (v3 ist der einzige Nutzer) in [PR](https://github.com/wailsapp/wails/pull/5711) von @taliesin-ai

## Behoben

- WebView2-Erkennung der Monitorskalierung und Korrektur der Host-Neusynchronisierung bei DPI-Änderungen in den Abschnitt „Unveröffentlicht“ verschoben in [PR](https://github.com/wailsapp/wails/pull/5750) von @taliesin-ai
- WebView2-COM-Marshalling für float64- und BOOL-Parameter aktualisiert in [PR](https://github.com/wailsapp/wails/pull/5741) von @wayneforrest
- Panic und nil-Dereferenzierung beim Aktualisieren und Zerstören von Symbolen im Windows-Infobereich verhindert in [PR](https://github.com/wailsapp/wails/pull/5703) von @wayneforrest
- Behoben, dass ausgeblendete Fenster unter Windows nicht korrekt erneut ausgeblendet werden in [PR](https://github.com/wailsapp/wails/pull/5743) von @wayneforrest
- Sichtbarkeit des WebView2-Controllers mit dem Minimieren, Maximieren und Wiederherstellen des Fensters synchronisiert in [PR](https://github.com/wailsapp/wails/pull/5742) von @wayneforrest

### Behoben

- WebView2-Erkennung der Monitorskalierung wieder aktiviert und Host-Neusynchronisierung auf Fälle mit geänderter DPI beschränkt in [PR](https://github.com/wailsapp/wails/pull/5734) von @taliesin-ai, basierend auf der von @randalmurphal validierten Korrektur, mit Überprüfung der Grundursache durch @eleclin und Hardwaretests durch @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Hinzugefügt

- Warnung beim Erstellen eines Release-AAB ohne gesetztes `ANDROID_KEYSTORE_FILE` hinzugefügt (Google Play lehnt mit Debug-Signatur signierte Bundles ab) und die Paketierung sowie Signierung von App Bundles in [PR](https://github.com/wailsapp/wails/pull/5730) durch @taliesin-ai dokumentiert
- Dokumentation Why Wails in mehreren Sprachen in [PR](https://github.com/wailsapp/wails/pull/5739) durch @taliesin-ai hinzugefügt
- Zuordnung von Go time.Time zu JS Date oder string in Bindings in [PR](https://github.com/wailsapp/wails/pull/5398) durch @fbbdev unterstützt
- Aufgaben für die Paketierung von Android App Bundles (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) zur Einreichung im Play Store hinzugefügt — APK-Aufgaben bleiben für Tests auf lokalen Geräten und Emulatoren erhalten; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5728) durch @mortenolsrud (behebt [#5726](https://github.com/wailsapp/wails/issues/5726))
- Aufgabenziele für physische Android-Geräte sowie die Wiederaufnahme der Kamera- und Standortberechtigungen in [PR](https://github.com/wailsapp/wails/pull/5735) durch @taliesin-ai hinzugefügt

## Geändert

- `webview2` auf v1.0.28 aktualisiert ([Versionshinweise](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- `compileSdk`/`targetSdk` der Android-Vorlage von 34 auf 35 aktualisiert, wie von Google Play für die Einreichung neuer Apps verlangt; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5730) durch @taliesin-ai

## Behoben

- Fest eingebettete Avatar-Masken in sponsorkit in [PR](https://github.com/wailsapp/wails/pull/5745) durch @leaanthony korrigiert
- Fehler behoben, durch den bei der automatischen Erstellung eines Android-AVD aufgrund lexikografischer Versionssortierung das falsche System-Image oder die falsche cmdline-tools-Version ausgewählt wurde; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5730) durch @taliesin-ai
- Fehler behoben, durch den der Einrichtungsassistent eine veraltete Android-NDK-Version vorgeschlagen hat (jetzt 26.3.11579264, entsprechend der dokumentierten Anforderung); umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5730) durch @taliesin-ai
- Französische Dokumentation für SvelteKit und Optionen in [PR](https://github.com/wailsapp/wails/pull/5744) durch @leaanthony aktualisiert
- SIGSEGV bei der Bildschirmaufzählung unter macOS während Anzeigeänderungen in [PR](https://github.com/wailsapp/wails/pull/5516) durch @flofreud behoben

## v3.0.0-alpha2.112 - 2026-07-03

## Hinzugefügt

- Go-basierten SVG-Generator für Mitwirkende hinzugefügt und die Credits-Seiten der Dokumentation und Website in [PR](https://github.com/wailsapp/wails/pull/5724) durch @taliesin-ai aktualisiert
- Zuordnung von Go time.Time zu JS Date oder string in Bindings in [PR](https://github.com/wailsapp/wails/pull/5398) durch @fbbdev unterstützt

## Geändert

- Node-basierte Pipeline für Sponsorenbilder in [PR](https://github.com/wailsapp/wails/pull/5719) durch einen Go-Generator ersetzt; umgesetzt durch @taliesin-ai

## Behoben

- Installationsskript für Abhängigkeiten der Android-Build-Assets in [PR](https://github.com/wailsapp/wails/pull/5729) durch @taliesin-ai korrigiert
- Steuerzeichen U+0085 (NEXT LINE) in `ValidateAndSanitizeURL` zurückgewiesen und damit die Abdeckung von Leerraumzeichen durch den URL-Validator vervollständigt
- DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster in [PR](https://github.com/wailsapp/wails/pull/4785) durch @leaanthony neu berechnet
- Fehler behoben, durch den die Erkennung von DnD-Ablagezonen unter Windows bei einer Skalierung ungleich 100 % fehlschlug; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/4632) durch @yulesxoxo
- Explizite Objective-C-Speicherverwaltung für Cocoa-Objekte in Dialogen, Menüs, Menüleistensymbolen und Benachrichtigungen unter Darwin in [PR](https://github.com/wailsapp/wails/pull/5714) durch @taliesin-ai hinzugefügt
- Fehler im Linux-CGO-Backend und Probleme mit dem Taskleistensymbol in [PR](https://github.com/wailsapp/wails/pull/5718) durch @taliesin-ai behoben

## v3.0.0-alpha2.111 - 2026-07-01

## Hinzugefügt

- HappyTools in [PR](https://github.com/wailsapp/wails/pull/5061) durch @Aliuyanfeng zur Community-Präsentation hinzugefügt
- Unterstützung für die indonesische Lokalisierung und umfassende Dokumentation in [PR](https://github.com/wailsapp/wails/pull/5643) durch @triadmoko hinzugefügt
- Option DisableMenu in [PR](https://github.com/wailsapp/wails/pull/4813) durch @leaanthony zu WindowsWindow hinzugefügt

## Geändert

- Taskfile-Vorlage und CLI aktualisiert, um Build-/Paketierungsaufgaben mit GOOS und ARCH aufzurufen; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5617) durch @leaanthony

## Behoben

- Problem mit Fenster-Tabs unter macOS in [PR](https://github.com/wailsapp/wails/pull/5708) durch @taliesin-ai behoben

## Entfernt

- Deutsche MDX-Übersetzungen aus den Bereichen Mitwirken, Funktionen und Anleitungen in [PR](https://github.com/wailsapp/wails/pull/5702) durch @taliesin-ai entfernt

## v3.0.0-alpha2.110 - 2026-06-30

## Hinzugefügt

- Neuladen und erzwungenes Neuladen der WebView unter macOS implementiert sowie Wiederherstellung nach Beendigung des WebContent-Prozesses hinzugefügt; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/5129) durch @wayneforrest
- Umfassende deutsche Dokumentation für die Bereiche Mitwirken, Funktionen und Anleitungen in [PR](https://github.com/wailsapp/wails/pull/5396) durch @leaanthony hinzugefügt
- Benachrichtigungen um Töne, Anhänge, Zeitplanung und eine API zum Aktualisieren von Benachrichtigungen erweitert in [PR](https://github.com/wailsapp/wails/pull/5333) von @popaprozac

## Behoben

- DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster in [PR](https://github.com/wailsapp/wails/pull/4785) durch @leaanthony neu berechnet
- Fehler behoben, durch den die Erkennung von DnD-Ablagezonen unter Windows bei einer Skalierung ungleich 100 % fehlschlug; umgesetzt in [PR](https://github.com/wailsapp/wails/pull/4632) durch @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Hinzugefügt

- Codebeispiele zur EventsEmit-Dokumentation in [PR](https://github.com/wailsapp/wails/pull/5026) durch @iamhabbeboy hinzugefügt
- Option für visuelles Hosting mit Windows WebView2 in [PR](https://github.com/wailsapp/wails/pull/5380) durch @MerIijn hinzugefügt
- Klustr in [PR](https://github.com/wailsapp/wails/pull/5536) durch @SametKUM zur Dokumentation der Community-Präsentation hinzugefügt
- Kira mit neuen Seiten und einem Eintrag im Änderungsprotokoll in [PR](https://github.com/wailsapp/wails/pull/5685) durch @thiennguyen93 zur Community-Präsentation hinzugefügt
- Feedback-Abschnitt in [PR](https://github.com/wailsapp/wails/pull/5694) durch @taliesin-ai zum Leitfaden für den MCP-Dienst hinzugefügt

## Geändert

- Der Servermodus verfügt nun über einen vollwertigen Produktions-Build, der mit den Desktop-Build-Tasks übereinstimmt (#5693). `task build:server` erstellt standardmäßig eine Produktionsbinärdatei (`-tags server,production`, `-trimpath`, ohne Symbole) und akzeptiert `DEV=true` (Entwicklungsserver), `OBFUSCATED=true` (Garble) sowie `EXTRA_TAGS`. `task run:server` führt einen Entwicklungsserver aus. `Dockerfile.server` / `task build:docker` erstellen zuerst den Produktionsserver (`-tags server,production`) und das Produktions-Frontend. Das Image verwendet standardmäßig einen statischen Pure-Go-Build auf distroless/static; `CGO_ENABLED`, `GO_IMAGE` und `RUNTIME_IMAGE` werden als überschreibbare Build-Argumente für CGO-Anwendungen bereitgestellt.

## Behoben

- Absturz beim Schließen eines Fensters mit ausstehenden asynchronen Aufrufen verhindert, in [PR](https://github.com/wailsapp/wails/pull/4435) von @leaanthony
- Aktivierung des Fensters beim Öffnen ausgeblendeter Anwendungen unter Windows verhindert, in [PR](https://github.com/wailsapp/wails/pull/5249) von @leaanthony
- Sichergestellt, dass Metadaten von WebKit-Anfragen, der Abschluss von Antworten und die Verarbeitung des Body-Streams im GTK-Hauptthread ausgeführt werden, in [PR](https://github.com/wailsapp/wails/pull/5668) von @taliesin-ai
- Fehler behoben, durch den `Menu.Update()` das native Menü unter GTK4 Linux nicht neu erstellte (#5659, unabhängig diagnostiziert und in #5539 von @puneetdixit200 behoben)
- Absturz beim Auflisten der macOS-Bildschirme nach einer Anzeigeänderung behoben, indem die Zeichenfolgen für Bildschirm-ID und -Namen kopiert und die Anzahl als Momentaufnahme erfasst werden (#5565, unabhängig diagnostiziert und in #5584 von @x-haose behoben)
- Fehler behoben, durch den WebView2-Inhalte nach dem Ziehen eines Fensters über Monitore mit unterschiedlichen DPI-Einstellungen unter Windows zunächst schrumpften und dann verschwanden. Dazu werden die Controller-Grenzen im `WM_DPICHANGED`-Handler erneut gesetzt, analog zur erneuten DPI-Synchronisierung beim Wiederherstellen aus dem minimierten Zustand (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Hinzugefügt

- Globale (systemweite) Tastenkürzel über `app.GlobalShortcut` hinzugefügt (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Die Tastenkürzel werden auch ausgelöst, wenn die Anwendung nicht fokussiert ist. Sie sind für jede Plattform nativ und ohne Drittanbieterabhängigkeiten implementiert: Carbon-Hotkeys unter macOS, `RegisterHotKey` unter Windows, `XGrabKey` unter X11 und die Schnittstelle für globale Tastenkürzel des XDG Desktop Portal unter Wayland.
- Integrierten MCP-Server hinzugefügt: einen Model-Context-Protocol-Server, der automatisch startet, wenn die Anwendung mit dem Tag `mcp` erstellt wird. Dadurch können LLM-Agenten eine laufende Wails-Anwendung testen und steuern – einschließlich Fenstersteuerung, DOM-Inspektion, JavaScript-Auswertung, Aufrufen gebundener Methoden, Ereignissen und simulierter Maus-/Tastatureingaben, die mit einem animierten Bildschirmcursor dargestellt werden. Kein Benutzercode erforderlich: Das Tag `mcp` wird von `wails3 build`/`wails3 dev` automatisch hinzugefügt, wenn `WAILS_MCP=1` gesetzt ist. Die Konfiguration erfolgt vollständig über Umgebungsvariablen (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Behoben

- Fehler behoben, durch den `Menu.Update()` das native Menü unter GTK4 Linux nicht neu erstellte (#5659, unabhängig diagnostiziert und in #5539 von @puneetdixit200 behoben)
- Absturz beim Auflisten der macOS-Bildschirme nach einer Anzeigeänderung behoben, indem die Zeichenfolgen für Bildschirm-ID und -Namen kopiert und die Anzahl als Momentaufnahme erfasst werden (#5565, unabhängig diagnostiziert und in #5584 von @x-haose behoben)

## v3.0.0-alpha2.107 - 2026-06-27

## Hinzugefügt

- Experimentelle Wake-Dokumentation mit Seitenleistennavigation hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5613) von @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Geändert

- `webview2` auf v1.0.27 aktualisiert.
  - ci(webview2): Release-Build korrigiert (Windows-Cross-Kompilierung + vollständige go.sum) (#5671)\

  **Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- go vet aus der Cross-Kompilierung im WebView2-Release-Workflow entfernt, in [PR](https://github.com/wailsapp/wails/pull/5672) von @taliesin-ai
- OpenRouter-Modell des automatischen Changelogs auf google/gemini-2.5-flash-lite aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5670) von @taliesin-ai
- `webview2` auf v1.0.26 aktualisiert.

### Korrekturen

- **Wiederherstellung nach vorübergehenden COM-Laufzeitfehlern, statt die Anwendung zu beenden** (#5658, #5580). `Chromium.errorCallback` rief zuvor bei *jedem* COM-Fehler `os.Exit(1)` auf, sodass eine behebbare Störung nach dem Start die gesamte Anwendung beendete. Laufzeitpfade (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) protokollieren den Fehler nun und stellen den Betrieb wieder her. Insbesondere wird eine fehlerhafte oder nicht vertrauenswürdige Webnachricht in `MessageReceived` nun verworfen, statt den Prozess zu beenden. Dies behebt die Klasse von Abstürzen beim Verschieben über Monitore mit unterschiedlichen DPI-Einstellungen (#5544, #5650). Fehler in Pfaden zur Erstellung von Umgebung und Controller bleiben fatal.\

**Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Behoben

- release-webview2-Workflow für die korrekte Verarbeitung von go.sum-Dateien korrigiert, in [PR](https://github.com/wailsapp/wails/pull/5671) von @taliesin-ai
- Menüaktualisierungen unter Linux mit GTK4 korrigiert, indem das native Menü geleert und neu erstellt wird, in [PR](https://github.com/wailsapp/wails/pull/5659) von @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Hinzugefügt

- `application.System` zur Laufzeiterkennung der Plattform aus gemeinsam genutztem Code hinzugefügt: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (das Build-Tag `server`) und `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` zum direkten Testen eines einzelnen Ziels. Es lässt sich für jedes Ziel kompilieren, sodass Verzweigungen ohne Build-Tags möglich sind. Entsprechende Frontend-Hilfsfunktionen (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) sind in `@wailsio/runtime` verfügbar
- Leitfaden „Andere Frontend-Frameworks verwenden“ hinzugefügt, der zeigt, wie ein eigenes Vite-Projekt in `frontend/` eingefügt wird (behandelt Solid, Preact, Lit, SvelteKit, Qwik, Angular usw.)
- Der Assistent `wails3 setup` prüft nun die mobile Toolchain für iOS/Android – Xcode und die iOS-Simulator-Laufzeit, JDK, Android SDK/NDK und Emulator – und bietet, sofern anwendbar, eine Installation mit einem Klick sowie kopierbare Korrekturen für die Shell-Konfiguration
- Generierte Projekte enthalten eine `frontend/.npmrc`, die für `minimum-release-age` 7 Tage festlegt, um die Gefährdung durch neu veröffentlichte und möglicherweise kompromittierte Pakete zu verringern (wird von pnpm und bun berücksichtigt; von npm ohne Auswirkungen ignoriert)

## Geändert

- Alle integrierten Starter-Vorlagen mit einem neuen Hero-Design im Neongebirgsstil überarbeitet (Web, iOS und Android)
- **TypeScript ist nun die Standardsprache für Starter-Vorlagen und verwendet den unveränderten Vorlagennamen.** `wails3 init` (ohne `-t`) erzeugt das Grundgerüst eines TypeScript-Projekts; `-t vanilla`, `-t react`, `-t vue` und `-t svelte` verwenden TypeScript. JavaScript-Varianten sind unter `-t vanilla-js`, `-t react-js`, `-t vue-js` und `-t svelte-js` verfügbar. Integrierte Vorlagen deklarieren ihre Sprache mit `typescript:` in `template.yaml`; Community-Vorlagen mit dem Suffix `-ts` funktionieren weiterhin als Fallback
- Den Assistenten `wails3 setup` mit dem Neon-Design „digital Wails“ überarbeitet (Milchglas-Vibranz vor einer Gebirgskulisse)

## Behoben

- Absturz unter Windows beim Wiederherstellen einer App behoben, die so lange minimiert war, dass WebView2 angehalten oder dessen Render-/GPU-Prozess beendet und neu erstellt wurde. Die DPI-Neusynchronisierung beim Minimieren/Wiederherstellen (#5544) greift jetzt nur noch auf den WebView2-Controller zu, wenn sich die DPI des Fensters tatsächlich geändert hat. Dadurch werden beim üblichen Wiederherstellen mit unveränderter DPI fatale COM-Aufrufe an einen angehaltenen Controller vermieden (#5605)
- Wiederholte native `SIGABRT`/`SIGSEGV`-Abstürze in lange laufenden Linux-Apps bei häufigem Laden von Assets oder Medien behoben (typischerweise innerhalb von `g_object_unref` während der GTK-Hauptschleife). Der Asset-Server schloss `WebKitURISchemeRequest`s aus Worker-Goroutinen ab und rief dadurch nicht threadsichere WebKit2GTK-Funktionen außerhalb des GTK-Hauptthreads auf. Der Abschluss (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) erfolgt jetzt im Hauptthread. Vervollständigt die teilweise Fehlerbehebung aus #5566. Betrifft sowohl GTK3- als auch GTK4/WebKitGTK-6.0-Builds (#5631, #5557)
- Sporadisches `fatal error: invalid pointer found on stack` in `setupSignalHandlers` unter Linux/GTK3 behoben. Als Signal-`user_data` übergebene Fenster-IDs wurden in einer lokalen Go-Variable vom Typ `unsafe.Pointer` gehalten. Daher brach der Garbage Collector ab, als er den Wert, der kein Zeiger war, während einer Stack-Kopie untersuchte. Die ID behält auf der Go-Seite jetzt einen Ganzzahltyp (`uintptr_t`). Damit wird dieselbe Fehlerbehebung, die #4958 auf den GTK4-Pfad anwendete, auf den älteren GTK3-Pfad zurückportiert. Beim GTK4-Pfad wurden die C-Signalfunktionen auf `uintptr_t` umgestellt, um `-race`-/checkptr-Fehler zu beseitigen (#5631)

## Entfernt

- Die Starter-Vorlagen `react-swc`, `preact`, `lit`, `solid`, `qwik` und `sveltekit` sowie ihre `-ts`-Varianten wurden entfernt. Die unterstützten integrierten Vorlagen sind jetzt `vanilla`, `react`, `vue` und `svelte`. Sie verwenden standardmäßig jeweils TypeScript; JavaScript-Varianten sind unter `-js` verfügbar. Jedes andere Framework kann weiterhin verwendet werden, indem Sie [Ihr eigenes Frontend einbinden](https://v3.wails.io/guides/dev/frontend-frameworks) oder eine benutzerdefinierte Vorlage nutzen

## v3.0.0-alpha2.104 - 2026-06-18

## Behoben

- iOS-Absturz (SIGABRT) behoben, wenn eine gebundene Go-Dienstmethode eine leere Zeichenfolge zurückgibt. Der Writer für iOS-Asset-Antworten prüfte den Body-Zeiger mit `buf != nil` statt dessen Länge, sodass ein Body mit der Länge null eine Panik in `&buf[0]` auslöste. Jetzt wird die Länge geprüft, wie bei den Desktop-Writern

## v3.0.0-alpha2.103 - 2026-06-15

## Geändert

- Native iOS- und Android-Funktionen wurden in Plattformmanager verschoben: Rufen Sie sie über `application.IOS.*` und `application.Android.*` auf (z. B. `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) statt über die bisherigen freien Funktionen `application.IOS*`/`application.Android*` (#5602)
- Mobile-Bridge-Ereignisse wurden umbenannt: Plattformübergreifende Ereignisse verwenden jetzt das Präfix `common:*` (z. B. `common:haptic`, `common:location`), plattformspezifische Ereignisse dagegen `ios:*` bzw. `android:*` (z. B. `ios:backgroundTask`, `android:foregroundService`). Das Präfix `native:*` wird nicht mehr verwendet (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Hinzugefügt

- Experimenteller `wails3 setup`-Assistent zur interaktiven Projekteinrichtung und Abhängigkeitsprüfung hinzugefügt
- Flag `--json` für maschinenlesbare Ausgaben zu `wails3 doctor` hinzugefügt
- Abschnitt zum Signierungsstatus zum Befehl `wails3 doctor` hinzugefügt

## Behoben

- npm-Erkennung unter Linux korrigiert, sodass neben dem Paketmanager auch PATH geprüft wird

## v3.0.0-alpha.101 - 2026-06-13

## Hinzugefügt

- iOS: native Meldungsdialoge (UIAlertController) sowie Dialoge zum Öffnen einer Datei, mehrerer Dateien oder eines Verzeichnisses (UIDocumentPickerViewController); Speicherdialoge geben einen expliziten Fehler zurück
- iOS: Zwischenablagenunterstützung über UIPasteboard
- iOS: tatsächliche Bildschirmmetriken über UIScreen (Punkte, Pixel, Skalierung, Arbeitsbereich innerhalb der Safe Area)
- iOS: Geräte-Builds (`IOS_PLATFORM=device`), Unterstützung für Codesignierungsidentität, Bereitstellungsprofil und Berechtigungen, `.ipa`-Paketierung sowie `deploy-device` über devicectl
- iOS: konfigurierbare iOS-Mindestversion (`ios.minIOSVersion` in build/config.yml)
- iOS: `wails3 doctor` meldet unter macOS die Verfügbarkeit von Xcode und dem iOS SDK
- iOS: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und knappen Arbeitsspeicher werden als `events.IOS.*`- und plattformneutrale `events.Common.*`-Anwendungsereignisse bereitgestellt
- iOS: native Bridge für Mobilfunktionen (exportiertes `application.IOS*`) – Teilen-Dialog, URL öffnen, Aktivhalten, Taschenlampe, Safe-Area-Abstände, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (Face ID/Touch ID), lokale Benachrichtigungen und sichere Speicherung im Keychain
- iOS: Sensoren und Hardware – Haptik, einmalige Standortbestimmung, Beschleunigungsmesser, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Erkennung von Bildschirmaufnahmen
- iOS: Dokumentation (IOS.md und ein Leitfaden auf der Dokumentationswebsite)
- Android: native Meldungsdialoge (AlertDialog) sowie Dialoge zum Öffnen einer oder mehrerer Dateien (Storage Access Framework; als Kopien in den Cache importiert); Dialoge zum Öffnen eines Verzeichnisses und zum Speichern geben einen expliziten Fehler zurück
- Android: Zwischenablagenunterstützung über ClipboardManager
- Android: tatsächliche Bildschirmmetriken über WindowMetrics/DisplayMetrics (dp, Pixel, Skalierung, verfügbarer Arbeitsbereich unter Berücksichtigung der Systemleisten)
- Android: Laufzeitmethoden für Haptik (`Android.Haptics.Vibrate`), Geräteinformationen (`Android.Device.Info`) und Toast-Meldungen (`Android.Toast.Show`)
- Android: typisierte Lebenszyklusereignisse (`events.Android.*`, aus events.txt generiert), wobei `ActivityCreated` auf `Common.ApplicationStarted` abgebildet wird
- Android: Die Build-Pipeline erzeugt installierbare Debug- und Release-APKs (`android:run`, `android:package`, `android:package:fat`). Release-Builds werden standardmäßig mit dem Debug-Keystore signiert oder über `ANDROID_KEYSTORE_*`-Umgebungsvariablen mit einem echten Keystore
- Android: `wails3 doctor` meldet Android SDK, NDK und JDK
- Android: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und knappen Arbeitsspeicher werden als `events.Android.*`- und plattformneutrale `events.Common.*`-Anwendungsereignisse bereitgestellt
- Android: native Bridge für Mobilfunktionen (exportiertes `application.Android*`) – Teilen, URL öffnen, Aktivhalten, Taschenlampe, Safe-Area-Abstände, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (BiometricPrompt), lokale Benachrichtigungen und sichere Speicherung mit EncryptedSharedPreferences
- Android: Sensoren und Hardware – Haptik, einmalige Standortbestimmung, Beschleunigungsmesser, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Blockierung von Bildschirmaufnahmen durch FLAG_SECURE
- Android: Dokumentation (ANDROID.md und ein Leitfaden auf der Dokumentationswebsite)
- Beispiel: Das `mobile`-Kitchen-Sink-Beispiel erhält die Registerkarten „Mobil“ und „Hardware“, die die native Funktions-Bridge unter iOS und Android demonstrieren (pill-förmige Registerkarten werden auf mehrere Zeilen umgebrochen)
- Mobilgeräte: Akku – Beschleunigungsmesser, Näherungssensor, Taschenlampe und die periodische Uhr des Beispiels werden pausiert, wenn die App in den Hintergrund wechselt, und bei der Rückkehr wiederhergestellt. Android lässt den Prozess im Hintergrund weiterlaufen, und unter iOS bleibt der Hardwarezustand der Taschenlampe bestehen. Android-Systemereignisempfänger werden außerdem nur registriert, solange sich die App im Vordergrund befindet
- iOS: Kameraaufnahme — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → ein `native:capture`-Ereignis mit einem Base64-Vorschaubild)
- iOS: Hintergrundausführung — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (ein UIApplication-Hintergrundaufgaben-Zeitfenster) und eine konfigurierbare Option `ios.backgroundModes` (build/config.yml), die über die Vorlage `UIBackgroundModes` in die generierte Info.plist einfügt
- Android: Kameraaufnahme — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (Systemkamera über FileProvider → ein `native:capture`-Ereignis)
- Android: Vordergrunddienst — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (ein `WailsForegroundService` mit einer dauerhaften Benachrichtigung hält den Prozess für lang laufende Hintergrundaufgaben aktiv)
- Beispiel: Eine Registerkarte „Kamera“ demonstriert die Foto-/Videoaufnahme und Hintergrundausführung (Vordergrunddienst unter Android, Hintergrundaufgaben-Zeitfenster unter iOS)

## Behoben

- Behoben, dass `getUserMedia` unter Linux immer mit `NotAllowedError` fehlschlug: WebKitGTK lehnt Berechtigungsanfragen ab, die von niemandem verarbeitet werden, und das Signal `permission-request` war nicht verbunden. Kamera und Mikrofon werden jetzt anhand einer neuen plattformübergreifenden `WebviewWindowOptions.Permissions`-Zuordnung (`map[PermissionType]Permission`) behandelt, die sowohl unter Linux (WebKitGTK) als auch unter Windows (WebView2) berücksichtigt wird. Unter Linux, wo es keine native Abfrage gibt, sind Kamera und Mikrofon standardmäßig zugelassen (wodurch `getUserMedia` wiederhergestellt wird) und können mit `PermissionDeny` deaktiviert werden (#5552)
- iOS: `GOOS=ios` wird wieder kompiliert (exportiertes `events.IOS`, Stubs für mobile Methodennamen); außerdem werden Builds mit Produktions-Tag kompiliert (Korrekturen der Build-Tags in pkg/application und mehreren Diensten)
- iOS: Go→JS-Ereignisse und ExecJS funktionieren jetzt — die Seite wird beim Start nicht mehr zweimal geladen, und der `wails:runtime:ready`-Handshake kann nicht mehr verloren gehen
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` verursachen beim App-Start keine Race-Condition mehr; die feste Startverzögerung von 2 Sekunden wurde entfernt
- iOS: Ein C-String-Speicherleck bei jeder Go→JS-JavaScript-Ausführung wurde behoben
- iOS: `hasListeners` gibt jetzt die tatsächliche Listener-Registrierung wieder
- iOS: Das Debug-Logging des Frameworks wird aus Produktions-Builds herauskompiliert
- Android: `GOOS=android` wird wieder kompiliert — `events.Android` wurde definiert, das außerhalb der Grenzen liegende `events_android.go`-Listener-Array entfernt, der Stub für mobile Methodennamen hinzugefügt und verhindert, dass Desktop-Linux-Dateien (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) in Android-Builds einfließen
- Android: JS→Go-Bindings funktionieren jetzt — die WebView kann `fetch()`-POST-Bodys nicht an `shouldInterceptRequest` übergeben. Daher werden Runtime-Aufrufe über einen JavascriptInterface-Transport (`nativeHandleRuntimeCall`) geleitet, statt wegen eines nil-Anfragebodys abzustürzen
- Android: `Screens.*`-Runtime-Aufrufe geben echte Daten zurück — der ScreenManager wird jetzt beim Start befüllt (er war nie angebunden, sodass `GetAll` nil zurückgab)
- Android: Das Debug-Logging des Frameworks wird aus Produktions-Builds herauskompiliert und in Debug-Builds unter dem Tag `Wails` an logcat geleitet
- Android: echte `hasListeners`-Registry, Behandlung von JNI-Referenzen und -Ausnahmen sowie ein Seitenlebenszyklus mit einmaligem Laden (keine doppelte Navigation)
- Behoben, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug, indem generierte Dateien in das Ausgabeverzeichnis synchronisiert werden, statt es durch Umbenennen zu überschreiben (#5515)
- Ein sporadischer schwerwiegender Absturz unter macOS beim Lesen von Bildschirminformationen nach einer Anzeigeänderung wurde behoben: Die gespeicherten Zeiger für Bildschirm-ID und -Name verwiesen auf automatisch freigegebene `UTF8String`-Puffer, die freigegeben werden konnten, bevor Go sie kopierte (Use-after-free). Die Strings werden jetzt mit `strdup` dupliziert und nach der Konvertierung freigegeben. Außerdem läuft die Bildschirmauflistung in einem expliziten Autorelease-Pool, sodass bei Aufrufen aus Go-Goroutinen kein Speicher mehr verloren geht (#5556)
- Ein sporadischer SIGSEGV unter Linux beim Schließen eines `WebKitURISchemeRequest` durch den Assetserver wurde behoben: Das abschließende `g_object_unref` lief in der Assetserver-Goroutine und finalisierte dadurch ein WebKit-GObject außerhalb des GTK-Hauptthreads. Die unref-Operation wird jetzt über `g_main_context_invoke` an den GTK-Hauptkontext übergeben (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Hinzugefügt

- `MacWebviewPreferences` um zusätzliche WKWebView-Konfigurationsoptionen erweitert: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` und `ApplicationNameForUserAgent` (#5549)

## Behoben

- Behoben, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug, indem generierte Dateien in das Ausgabeverzeichnis synchronisiert werden, statt es durch Umbenennen zu überschreiben (#5561)
- Behoben, dass JS-Ereignisse zur Größenänderung bei rahmenlosen Fenstern unter Linux nicht ausgelöst wurden; außerdem wurde die Erkennung des Scrollleistenrands bei rahmenlosen Fenstern korrigiert (#5368)
- Behoben, dass der Updater unter Windows mit „invalid cross-device link“ fehlschlug, wenn sich das temporäre Verzeichnis auf einem anderen Volume als das Installationsverzeichnis befand (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Behoben

- Behoben, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug, indem generierte Dateien in das Ausgabeverzeichnis synchronisiert werden, statt es durch Umbenennen zu überschreiben (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Behoben

- Einfrieren der WebKit-Benutzeroberfläche unter Linux im Leerlauf (z. B. bei geöffnetem Inspektor) behoben, indem `SA_ONSTACK` nicht mehr für `SIGUSR1` erzwungen wird, da dies die Synchronisierung des GC-Threads von JavaScriptCore beeinträchtigte (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Hinzugefügt

- Debugging-Seite und Arbeiten mit `runtime/trace` hinzugefügt

## Geändert

- Einige unnötige `_ "embed"`-Importe entfernt und den Code dadurch etwas bereinigt

## Behoben

- Behoben, dass die Mindestbreite und -höhe unter Windows nach dem Aufheben der Fenstermaximierung nicht erzwungen wurden (#4593)
- Durchreichen von Mausklicks im Vollbildmodus bei den Fensteroptionen Frameless + Transparent behoben (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Hinzugefügt

- Unterstützung für Garble-Verschleierung hinzugefügt ([#4563](https://github.com/wailsapp/wails/issues/4563)): stabile Methoden-IDs für Bindings, Build-/Taskfile-Integration (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) und JSON-Struct-Tags für jede Runtime-seitige Nutzlast (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), damit das Übertragungsformat die Umbenennung exportierter Felder durch Garble übersteht.

## v3.0.0-alpha.95 - 2026-05-20

## Hinzugefügt

- Fehlende Seite zur Projektstruktur hinzugefügt

## Geändert

- Dokumentation: Einige Diagramme auf der Architekturseite wurden in Sequenzdiagramme geändert, um eine übersichtlichere Darstellung zu erzielen
- Dokumentation: Hinweis aufnehmen, dass D2 als Voraussetzung für die Ausführung installiert werden muss

## Behoben

- Fehler bei `wails3 generate appimage` mit der GTK4-Standardeinstellung behoben: Der Bundler erkennt nun den GTK-Stack aus der Binärdatei, bevor er nach Laufzeitdateien sucht. Dadurch wählt er für GTK4-Builds `libwebkitgtkinjectedbundle.so` (unter `webkitgtk-6.0/`) und für `-tags gtk3`-Builds `libwebkit2gtkinjectedbundle.so` (unter `webkit2gtk-4.1/`) aus. Die `.relr.dyn`-Prüfung berücksichtigt außerdem `libgtk-4.so.1`, sodass das Stripping bei modernen Toolchains unabhängig vom Stack korrekt deaktiviert wird. (#5475)
- Fehler behoben, durch den `wails3 generate appimage` bei Aufruf mit einem relativen `-builddir` fehlschlug: Der Bundler löst nun `-binary`, `-icon`, `-desktopfile`, `-builddir` und `-outputdir` von Anfang an in absolute Pfade auf, sodass das während des Ablaufs ausgeführte `s.CD` weder die AppRun-Download-Goroutine noch die nach dem Kopieren ausgeführte `ldd`-Prüfung beeinträchtigt.
- Fehler behoben, durch den `wails3 generate appimage` das fertige AppImage nicht nach `-outputdir` verschob, wenn das Feld `Name=` der Desktop-Datei nicht mit dem Basisnamen der Binärdatei übereinstimmte: Der Bundler erzwingt nun über die Umgebungsvariable `OUTPUT`, dass das AppImage-Plug-in von linuxdeploy das AppImage unter `<binary>-<arch>.AppImage` statt unter dem aus der Desktop-Datei abgeleiteten Namen speichert.
- Fehler behoben, durch den `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht ausgelöst wurden, nachdem der Stack aus GTK4 und WebKitGTK 6.0 in alpha.93 zur Standardeinstellung gemacht worden war. Der neue standardmäßige `application_linux.go`-`run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an ihre `Common.*`-Entsprechungen weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer zur Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Buildpfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Behoben

- Fehler behoben, durch den `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht ausgelöst wurden, nachdem der Stack aus GTK4 und WebKitGTK 6.0 in alpha.93 zur Standardeinstellung gemacht worden war. Der neue standardmäßige `application_linux.go`-`run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an ihre `Common.*`-Entsprechungen weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer zur Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Buildpfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Hinzugefügt

- `XDG_SESSION_TYPE` zur `wails3 doctor`-Ausgabe unter Linux hinzugefügt, von @leaanthony

## Behoben

- Absturz des Fenstermenüs unter Wayland behoben, der durch den Zugriff von appmenu-gtk-module auf ein noch nicht realisiertes Fenster verursacht wurde (#4769), von @leaanthony
- Absturz der GTK-Anwendung behoben, wenn der Anwendungsname ungültige Zeichen wie Leerzeichen oder Klammern enthält, von @leaanthony
- Fehler „Nicht genügend Arbeitsspeicher“ beim Initialisieren von Drag-and-drop unter Windows behoben (#4701), von @overlordtm
- Race-Condition im Callback-Speicher des Hauptthreads behoben, bei der zum Löschen aus der Map fälschlicherweise RLock verwendet wurde (Linux, macOS, iOS) (#4424), von @leaanthony
- Variablenverarbeitung bei der Übergabe von Befehlszeilenargumenten an Tasks korrigiert. Als KEY=VALUE-Paare angegebene CLI-Variablen werden nun ordnungsgemäß initialisiert und während der gesamten Task-Ausführung weitergegeben.
- Konflikt bei NSWindowZoomButton unter macOS behoben: `MaximiseButtonState` und `FullscreenButtonState` wenden nun sowohl beim Start als auch zur Laufzeit den restriktiveren Zustand an; keiner der beiden Setter kann den anderen unbemerkt überschreiben (#5319)
- Mehrere bereits vorhandene Fehler im veralteten GTK3-Buildpfad (`-tags gtk3`) behoben, die CodeRabbit bei #5463 aufgedeckt hat: Durch Dateizuordnungen gestartete Anwendungen überspringen die Startup-Handler nicht mehr; `getTheme` ist nun grenzen- und typsicher; `appName` gibt keinen GLib-eigenen Speicher mehr frei; `clipboardGet` gibt nun den von GTK zurückgegebenen `gchar*` ordnungsgemäß frei; `Calloc` verwendet nun Pointer-Receiver (und `NewCalloc` gibt `*Calloc` zurück), sodass der Pool Allokationen tatsächlich erfasst; `zoomOut` verwendet den Kehrwert von `zoomInFactor` statt eines negativen Multiplikators, der auf 1.0 begrenzt wurde; `execJS` verwendet den vorab allokierten leeren World-Namen wieder, statt bei jedem Aufruf einen `C.CString("")` ohne anschließende Freigabe zu allokieren; ein für die Entwicklung vorgesehenes `fmt.Println` wurde aus `menuItem.setAccelerator` entfernt. Behebt #5465.
- Dasselbe durch einen Value-Receiver verursachte `Calloc`-Speicherleck im standardmäßigen GTK4-Buildpfad (`linux_cgo.go`) behoben: Pointer-Receiver und `NewCalloc() *Calloc` stellen sicher, dass die fensterspezifischen `c.String(...)`-Allokationen tatsächlich erfasst und freigegeben werden.

## v3.0.0-alpha.92 - 2026-05-15

## Hinzugefügt

- Taskfiles so geändert, dass sich der verwendete Frontend-Paketmanager über die Option `PACKAGE_MANAGER` steuern lässt
- Vorlagendaten um `{{.Opn}}` und `{{.Cls}}` erweitert, damit sich Taskfile-Vorlagen vorhersehbarer erstellen lassen

## Geändert

- Einige der vorhandenen Taskfiles so geändert, dass sie `{{.Opn}} and {{.Cls}}` verwenden

## Behoben

- Fatalen `concurrent map read and map write`-Laufzeitfehler in `linuxSystemTray` behoben, der auftrat, wenn das Tray-Menü aktualisiert wurde, während das Panel es las.
- Für WebView2-Fehler und Stacktrace-Ausgaben wird nun `log` statt `fmt` verwendet, damit Meldungen nicht verloren gehen, wenn die Anwendung unter Windows ohne angehängte Konsole ausgeführt wird.

## v3.0.0-alpha.91 - 2026-05-12

## Geändert

- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5414) von `@github-actions[bot]` aktualisiert
- **NICHT ABWÄRTSKOMPATIBEL (macOS):** Das macOS-Koordinatensystem wurde vereinheitlicht, sodass `GetScreens`, `Position` und `SetPosition` denselben Raum verwenden: logische Punkte, eine nach unten verlaufende Y-Achse und `(0,0)` oben links auf dem primären Bildschirm. Dies entspricht Windows, GTK sowie den öffentlichen APIs von Electron und des Webs. Bildschirme, die sich physisch oberhalb des primären Bildschirms befinden, melden nun negative `Bounds.Y`-Werte statt wie zuvor positive. `Position()`- und `SetPosition()`-Werte werden nun in logischen Punkten statt in `points × primaryScale` angegeben. Die verlustfreie Umwandlung `Position()` → `SetPosition()` bleibt erhalten; absolute Werte aus Protokollen früherer Alpha-Builds oder manuell berechnete Workarounds, etwa die Multiplikation mit `primaryScale` oder das Spiegeln von Y anhand einer Bildschirmhöhe, müssen aktualisiert werden. Behebt [#5117](https://github.com/wailsapp/wails/issues/5117).

## Behoben

- DBus-Signalnamen und die Länge der Signaldaten werden nun vorsorglich validiert, um Panics zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5416) von @leaanthony
- Problem mit der Speichersicherheit bei der GTK-Menüverarbeitung unter Linux in [PR](https://github.com/wailsapp/wails/pull/5363) behoben, von @leaanthony
- Erkennung von NVIDIA-GPUs und Deaktivierung des DMA-BUF-Renderers unter Linux in [PR](https://github.com/wailsapp/wails/pull/5295) hinzugefügt, von @leaanthony
- Bildschirmübergreifende Y-Umrechnung von `SetPosition` unter macOS korrigiert: Die Höhe des primären Bildschirms dient nun als globale Referenz, sodass Fenster auf Monitoren, die gegenüber dem primären Bildschirm vertikal versetzt sind, an der richtigen Position erscheinen, in [#5117](https://github.com/wailsapp/wails/issues/5117)
- Git-PR-Vorlage so korrigiert, dass sie in [PR](https://github.com/wailsapp/wails/pull/5109) auf die richtige Feedback-URL verweist, von @wayneforrest
- Behebt eine Reihe von Abstürzen durch Windows-Infobereich-`SetMenu`, die durch einen fehlerhaften `DestroyMenu`-Systemaufruf verursacht wurden. Dieser übergab vier Argumente statt eines, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Gibt außerdem bei der Neuerstellung von Menüs HMENU- und HBITMAP-Handles frei, einschließlich der zur Laufzeit über `MenuItem.SetBitmap` zugewiesenen Handles, setzt veraltete Kontrollkästchen-/Optionsfeld-Zuordnungen in `Win32Menu.Update` zurück und entfernt einen redundanten `Update()`-Aufruf in `systemtray.updateMenu`, der die Anzahl der Zuweisungen verdoppelte. Lang laufende Infobereich-Anwendungen verlieren nun nicht mehr bei jeder Neuerstellung des Menüs GDI-/USER-Objekte.

## v3.0.0-alpha.90 - 2026-05-11

## Hinzugefügt

- Konfigurierbaren Anwendungsnamen für den WKWebView-User-Agent unter macOS hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5261) von @vinhvoit225
- Indirekte Abhängigkeit github.com/coder/websocket zum gin-service-Beispiel hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5400) von @taliesin-ai
- Unterstützung für Tiefengleichheitsvergleiche zu den Tests der Build-Assets hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5402) von @leaanthony

## Geändert

- Build-Ausgabe im assets-Verzeichnis zusammengeführt in [PR](https://github.com/wailsapp/wails/pull/5401) von @taliesin-ai
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5399) von `@github-actions[bot]` aktualisiert

## Behoben

- Benachrichtigungsobjekt für die Nachricht des Single-Instance-Mechanismus unter macOS verwendet in [PR](https://github.com/wailsapp/wails/pull/5289) von @overlordtm
- Windows-Callbacks gebündelt, um bei hoher Last den Verlust von Promises zu verhindern in [PR](https://github.com/wailsapp/wails/pull/5383) von @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Hinzugefügt

- Job go<em>test</em>results zum Zusammenfassen der Go-Testergebnisse hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5316) von @leaanthony

## Geändert

- Große RPC-Nutzlasten bei Bedarf in gestückelte POST-Anfragen aufgeteilt in [PR](https://github.com/wailsapp/wails/pull/5369) von @leaanthony
- Vite in allen Frontend-Vorlagen von 5.x.x auf 8.0.0 aktualisiert in [PR](https://github.com/wailsapp/wails/pull/5386) von @leaanthony
- Portkonfiguration des Vite-Entwicklungsservers auf Umgebungsvariablen umgestellt in [PR](https://github.com/wailsapp/wails/pull/5365) von @leaanthony
- Vite-Entwicklungsserver in allen Vorlagen für die Bindung an 127.0.0.1 konfiguriert in [PR](https://github.com/wailsapp/wails/pull/5361) von @leaanthony
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5384) von `@github-actions[bot]` aktualisiert

## Behoben

- Info.plist-Vorlagenfragmente während der Aktualisierung der Build-Assets bereinigt in [PR](https://github.com/wailsapp/wails/pull/5312) von @leaanthony
- Veralteten Zustand in macOS-Menüs behoben, indem die Menüelement-Mutatoren (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) synchron im Hauptthread angewendet werden. Dadurch wird die `dispatch_async`-Race-Condition beseitigt, aufgrund derer Menüs beim schnellen erneuten Öffnen den vorherigen Zustand anzeigten (#5002)
- `*_test.go`-Dateien im Entwicklungsmodus ignoriert, um unnötige Neuerstellungen zu verhindern in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony
- Segmentierungsfehler durch Menu.Update() verhindert, wenn die Anwendung nicht ausgeführt wird, in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667
- Neuzeichnen der Menüleiste unter Windows mit lastSizeWParam gesteuert in [PR](https://github.com/wailsapp/wails/pull/5382) von @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Geändert

- HiddenOnTaskbar auf die Verwendung von WS<em>EX</em>TOOLWINDOW umgestellt in [PR](https://github.com/wailsapp/wails/pull/5371) von @leaanthony
- Abhängigkeiten neu angeordnet und die webview2-replace-Direktive aus go.mod entfernt in [PR](https://github.com/wailsapp/wails/pull/5370) von @atterpac
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5358) von `@github-actions[bot]` aktualisiert

## Behoben

- Generische Indirektionsaliase entfernt und Map-Schlüsseltypen zusammengeführt in [PR](https://github.com/wailsapp/wails/pull/5331) von @fbbdev

## Entfernt

- PR-master-Workflow einschließlich Dokumentation, Go-Tests und Skip-Tests gelöscht in [PR](https://github.com/wailsapp/wails/pull/5377) von @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Hinzugefügt

- Koreanische Dokumentation für Wails v3 hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5352) von @leaanthony
- Französische Dokumentation für Installation und Schnellstart hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5354) von @leaanthony
- Portugiesische Dokumentation für Schnellstart, Konzepte und Community hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5355) von @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Hinzugefügt

- Französische Dokumentationslokalisierung hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5328) von @leaanthony
- Deutsche Sprachversion zur Dokumentationswebsite hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5343) von @leaanthony

## Geändert

- Alle 8 übersetzten Sprachversionen in der Dokumentationskonfiguration registriert in [PR](https://github.com/wailsapp/wails/pull/5347) von @leaanthony
- Verschiedene Windows-bezogene Dateien für WebView2 aktualisiert in [PR](https://github.com/wailsapp/wails/pull/5317) von @leaanthony

## Behoben

- Dialogweiterleitung unter Linux zwischen GTK3 und GTK4 aufgeteilt in [PR](https://github.com/wailsapp/wails/pull/5340) von @leaanthony
- Sichergestellt, dass Dialog-Callbacks im GTK-Thread ausgeführt werden, wodurch Segmentierungsfehler behoben werden, in [PR](https://github.com/wailsapp/wails/pull/5339) von @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Hinzugefügt

- URL der PR-Vorlage zum Repository hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5179) von @leaanthony
- Deutsche Dokumentation für Wails v3 hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5330) von @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Hinzugefügt

- Option hinzugefügt, mit der sich unter macOS verhindern lässt, dass die Escape-Taste den Vollbildmodus beendet, in [PR](https://github.com/wailsapp/wails/pull/5307) von @leaanthony
- Option hinzugefügt, mit der sich unter macOS verhindern lässt, dass die Escape-Taste den Vollbildmodus beendet, in [PR](https://github.com/wailsapp/wails/pull/5310) von @leaanthony
- Dokumentation für Pausa zur Community-Galerie in [PR](https://github.com/wailsapp/wails/pull/5288) von @yuseferi hinzugefügt

## Geändert

- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5308) von `@github-actions[bot]` aktualisiert
- Befehl zur Symbolgenerierung für den Umgang mit nicht unterstützten Plattformen in [PR](https://github.com/wailsapp/wails/pull/5309) von @leaanthony aktualisiert
- Boolesche Vollbild-API durch den dreistufigen ButtonState ersetzt und Plattformbindungen in [PR](https://github.com/wailsapp/wails/pull/5224) von @leaanthony implementiert

## Behoben

- WebView2-Fokusoperationen in [PR](https://github.com/wailsapp/wails/pull/5315) von @leaanthony gegen einen nil-Controller-Zustand abgesichert
- GitHub-Actions-Workflow in [PR](https://github.com/wailsapp/wails/pull/5313) von @leaanthony aktualisiert, damit er den Basis-Branch des PR korrekt referenziert
- `*_test.go`-Dateien im Entwicklungsmodus ignoriert, um unnötige Neubuilds zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony
- Segmentierungsfehler durch Menu.Update() verhindert, wenn die Anwendung nicht ausgeführt wird, in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Hinzugefügt

- InstallScope-Flag und Build-Option für die Installation pro Computer oder Benutzer in [PR](https://github.com/wailsapp/wails/pull/5094) von @symball hinzugefügt
- Wirkungslose SetScreen-Methode zu BrowserWindow hinzugefügt, um das Window-Interface zu erfüllen, in [PR](https://github.com/wailsapp/wails/pull/5294) von @leaanthony

## Behoben

- Erkennung von NVIDIA-GPUs und Deaktivierung des DMA-BUF-Renderers unter Linux in [PR](https://github.com/wailsapp/wails/pull/5295) von @leaanthony hinzugefügt
- Git-PR-Vorlage in [PR](https://github.com/wailsapp/wails/pull/5109) von @wayneforrest korrigiert, sodass sie auf die richtige Feedback-URL verweist
- Eine Reihe von Abstürzen beim Aufruf von `SetMenu` in der Windows-Systray-Implementierung behoben, die durch einen fehlerhaften `DestroyMenu`-Systemaufruf verursacht wurden: Dieser übergab vier Argumente statt eines, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Außerdem werden bei der Neuerstellung von Menüs HMENU- und HBITMAP-Handles freigegeben, einschließlich der zur Laufzeit über `MenuItem.SetBitmap` allokierten Handles, veraltete Checkbox-/Optionsfeld-Zuordnungen in `Win32Menu.Update` zurückgesetzt und ein redundanter `Update()`-Aufruf in `systemtray.updateMenu` entfernt, der die Anzahl der Allokationen verdoppelte. Lang laufende Infobereichsanwendungen verursachen nun nicht mehr bei jeder Neuerstellung des Menüs Ressourcenlecks durch nicht freigegebene GDI-/USER-Objekte.

## v3.0.0-alpha.82 - 2026-05-01

## Behoben

- Generierung der Desktop-Datei in [PR](https://github.com/wailsapp/wails/pull/5232) von @leaanthony korrigiert, damit der Desktop-Name richtig verarbeitet wird

## v3.0.0-alpha.81 - 2026-04-30

## Geändert

- Zeitplan für Nightly-Releases in [PR](https://github.com/wailsapp/wails/pull/5286) von @leaanthony auf 15:00 UTC angepasst

## Behoben

- Halbierte Werte für Screen Bounds, WorkArea und Size auf Retina-Macs behoben –  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Geändert

- Dokumentationsabhängigkeiten und Loader für Inhaltssammlungen in [PR](https://github.com/wailsapp/wails/pull/5285) von @leaanthony aktualisiert

## v3.0.0-alpha.79 - 2026-04-29

## Hinzugefügt

- Berechtigung actions: write für den trigger-release-Job in [PR](https://github.com/wailsapp/wails/pull/5270) von @leaanthony erteilt

## Geändert

- Release-Task verwendet standardmäßig den master-Branch und aktualisiert die Formulierung des Änderungsprotokolls in [PR](https://github.com/wailsapp/wails/pull/5283) von @leaanthony
- Auto-Changelog-Workflow in [PR](https://github.com/wailsapp/wails/pull/5282) von @leaanthony zur Verwendung der neuesten Version aktualisiert
- Workflow-Effizienz durch Hinzufügen von Pfadfiltern und Entfernen nicht mehr verwendeter Workflows in [PR](https://github.com/wailsapp/wails/pull/5280) von @leaanthony verbessert
- Dokumentation in [PR](https://github.com/wailsapp/wails/pull/5274) von @leaanthony aktualisiert, sodass Beispiellinks auf den master-Branch verweisen
- Dokumentation und Beispiele für v3 in [PR](https://github.com/wailsapp/wails/pull/5272) von @leaanthony aktualisiert

## Behoben

- Reverse-Proxy in [PR](https://github.com/wailsapp/wails/pull/5265) von @AkagiYui um Wiederholungslogik und erzwungene IPv4-Nutzung für die Entwicklung erweitert
- Trigger-Workflow für das unveröffentlichte Änderungsprotokoll in [PR](https://github.com/wailsapp/wails/pull/5281) von @leaanthony neu geschrieben

## Entfernt

- Shell-Testskripte für verschiedene Testzwecke in [PR](https://github.com/wailsapp/wails/pull/5267) von @leaanthony entfernt
- Workflow zur Bereitstellung der v3-alpha-Dokumentation und CNAME-Eintrag in [PR](https://github.com/wailsapp/wails/pull/5266) von @leaanthony gelöscht

### Hinzugefügt

- Eintrag „Frontend-Routing“ zur Seitenleistennavigation in [PR](https://github.com/wailsapp/wails/pull/5196) von @leaanthony hinzugefügt
- Leitfaden zum Frontend-Routing mit Framework-spezifischen Empfehlungen in [PR](https://github.com/wailsapp/wails/pull/5185) von @leaanthony hinzugefügt
- Unterstützung für modale Sheets (macOS) hinzugefügt
- ghw-Version für eine bessere Unterstützung von Apple-Geräten durch @leaanthony aktualisiert (#4977)
- Methode `GetBadge` zum Dock-Dienst hinzugefügt
- Flag `-tags` zum Befehl `wails3 build` hinzugefügt, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Dokumentation zur automatischen Enum-Generierung im Binding-Generator hinzugefügt, einschließlich einer eigenen Seite „Enums“ und der Seitenleistennavigation (#4972)
- Flag `-tags` zum Befehl `wails3 build` hinzugefügt, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Web-API-Beispiele in `v3/examples/web-apis/` hinzugefügt, die 41 Browser-APIs demonstrieren, darunter Speicher (localStorage, sessionStorage, IndexedDB, Cache API), Netzwerk (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Medien (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Gerät (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Leistung (Performance API, Mutation Observer, Intersection/Resize Observer), Benutzeroberfläche (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) und weitere
- Beispiel für eine WebView-API-Kompatibilitätsprüfung (`v3/examples/webview-api-check/`) hinzugefügt, die plattformübergreifend mehr als 200 Browser-APIs testet
- Paket `internal/libpath` zum Ermitteln der Pfade nativer Bibliotheken unter Linux hinzugefügt, mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix
- **WIP:** Experimentelle Unterstützung für WebKitGTK 6.0 / GTK4 unter Linux hinzufügen, verfügbar über `-tags gtk4` (GTK3/WebKit2GTK 4.1 bleibt die Standardeinstellung)
- Hinweis: Bei Tiling-Fenstermanagern (z. B. Hyprland, Sway) funktionieren Minimieren und Maximieren möglicherweise nicht wie erwartet, da der Fenstermanager die Fenstergeometrie steuert.
- Anleitung zu **einmalig ausgeführten Ereignishandlern** in der Dokumentation zum **Lauschen auf Ereignisse in JavaScript** von @AbdelhadiSeddar hinzugefügt
- Option `UseApplicationMenu` zu `WebviewWindowOptions` hinzufügen, damit Fenster unter Windows/Linux das über `app.Menu.Set()` festgelegte Anwendungsmenü übernehmen können, von @leaanthony
- Unterstützung für `.icon`-Dateien (Apple-Icon-Composer-Format) zur Erzeugung von Liquid-Glass-Symbolen und Asset-Katalogen unter macOS hinzufügen (#4934), von @wimaha
- Experimentellen Servermodus für Headless-/Web-Bereitstellungen hinzufügen (`-tags server`). Damit können Wails-Apps ohne native GUI-Abhängigkeiten als HTTP-Server ausgeführt werden. Mit `wails3 task build:server` bauen. Weitere Informationen finden Sie unter `examples/server`.
- Paket `internal/libpath` zum Ermitteln nativer Bibliothekspfade unter Linux mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix hinzufügen
- Option `CollectionBehavior` zu `MacWindow` hinzufügen, um das Fensterverhalten über macOS-Spaces und den Vollbildmodus hinweg zu steuern (#4756), von @leaanthony
- Unit-Tests für pkg/application hinzufügen, von @leaanthony
- Unterstützung benutzerdefinierter Protokolle zur MSIX-Paketierung hinzufügen, von @leaanthony
- Erkennung der Desktop-Umgebung unter Linux hinzufügen, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Methode `Window.Print()` zur JavaScript-Laufzeit hinzufügen, um den Druckdialog vom Frontend aus zu öffnen (#4290), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Weitere WebKit2-Ereignisse für Änderungen des Ladezustands unter Linux hinzufügen: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Datei `.desktop` bereits beim Linux-Build und nicht erst bei der Paketierung erzeugen (#4575)
- Dokumentation der Linux-Laufzeitabhängigkeiten mit distributionsspezifischen Paketnamen und Beispielen für die nfpm-Paketierung hinzufügen (#4339), von @leaanthony
- Informationen zur NVIDIA-Treiberversion zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Ursprung zum Handler für Rohdaten-Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Transportschicht für Bindings überarbeiten, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- aria-label-Kennungen zu den helloworld-Vorlagen hinzufügen, damit Appium-Testclients die Beispiel-App einfach testen können, von @chinenual in [PR](https://github.com/wailsapp/wails/pull/4760)
- Ursprung zum Handler für Rohdaten-Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Transportschicht für Bindings überarbeiten, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- Typisierte Ereignisse von @fbbdev und @ianvs in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Beispiel `systray-clock` hinzufügen, das ein Headless-Taskleistensymbol mit live aktualisierten Tooltips zeigt (#4653).
- NSIS-Protokollvorlage für Windows hinzugefügt, von @Tolfx in #4510
- Tests für build-assets hinzugefügt, von @Tolfx in #4510
- macOS: Zeigt native Fenstersteuerelemente in der Menüleiste, in [#4588](https://github.com/wailsapp/wails/pull/4588) von @nidib
- macOS-Dock-Dienst zum Ausblenden/Einblenden des App-Symbols im Dock hinzufügen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- macOS-Dock-Dienst zum Ausblenden/Einblenden des App-Symbols im Dock hinzufügen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- Unterstützung für den nativen Liquid-Glass-Effekt unter macOS mit NSGlassEffectView (macOS 15.0+) und NSVisualEffectView als Fallback hinzufügen, einschließlich umfassender Optionen zur Materialanpassung, von @leaanthony in [#4534](https://github.com/wailsapp/wails/pull/4534)
- Bereinigung von Browser-URLs durch @leaanthony in [#4500](https://github.dev/wailsapp/wails/pull/4500). Basiert auf [#4484](https://github.com/wailsapp/wails/pull/4484) von @APShenkin.
- Inhaltsschutz unter Windows/macOS durch [@leaanthony](https://github.com/leaanthony) hinzufügen, basierend auf der ursprünglichen Arbeit von [@Taiterbase](https://github.com/Taiterbase) in diesem [PR](https://github.com/wailsapp/wails/pull/4241)
- Unterstützung für die Übergabe von CLI-Variablen an Task-Befehle über die Aliasse `wails3 build` und `wails3 package` hinzufügen (#4422), von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4488)
- Unterstützung für Dropzones, bei denen Ereignisse Daten abgelegter Elemente bereitstellen, von [@atterpac](https://github.com/atterpac) in [#4318](https://github.com/wailsapp/wails/pull/4318)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snap-Assist-Funktion von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snap-Assist-Funktion von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- Windows-Implementierung von `getAccentColor` durch [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427) hinzufügen
- Windows-Implementierung von `getAccentColor` durch [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427) hinzufügen
- Menüs und Menüleiste mit dunklem Design unter Windows. Von @leaanthony in [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Integrierte Dienste für eindeutigere JS-/TS-Bindings umbenennen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` zum Abrufen der Akzentfarbe des Systems eines Benutzers. Funktioniert unter macOS. Von [@etesam913](https://github.com/etesam913)
- API `window.ToggleFrameless()` von [@atterpac](https://github.com/atterpac) in [#4137](https://github.com/wailsapp/wails/pull/4137) hinzufügen
- Distributionsspezifische Build-Abhängigkeiten für Linux hinzugefügt von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Leitfaden zu Bindings hinzugefügt von @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
- **Organisierte Testinfrastruktur**: Docker-Testdateien in ein eigenes `test/docker/`-Verzeichnis mit optimierten Images verschoben und die Zuverlässigkeit des Builds verbessert von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Verbesserte Muster zur Ressourcenverwaltung**: Ordnungsgemäße Bereinigung von Event-Handlern und kontextabhängige Verwaltung von Goroutinen in Beispielen hinzugefügt von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- Unterstützung für aarch64-AppImage-Builds hinzugefügt von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Diagnoseabschnitt zu `wails doctor` hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Beim Aufruf einer Servicemethode das Fenster zum Kontext hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- `window-call`-Beispiel hinzugefügt, das zeigt, wie sich ermitteln lässt, welches Fenster einen Service aufruft, von [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Verbesserte Behandlung von Panics von [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Dokumentationskommentare für die Service-API hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Funktion `application.NewServiceWithOptions` zum Initialisieren von Services mit zusätzlicher Konfiguration hinzugefügt von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Verbesserte Menüsteuerung von [@FalcoG](https://github.com/FalcoG) und [@leaanthony](https://github.com/leaanthony) in [#4031](https://github.com/wailsapp/wails/pull/4031)
- Weitere Dokumentation von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für das Abbrechen von Events in Standard-Event-Listenern hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`Hide`, -`Show` und -`Destroy` von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`SetTooltip` von [@leaanthony](https://github.com/leaanthony). Ursprüngliche Idee von [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Paketpfad in Warnungen des Binding-Generators zu nicht unterstützten Typen ausgegeben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für generische Aliasse zum Binding-Generator hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für das JSON-Flag `omitzero` zum Binding-Generator hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:ignore` hinzugefügt, um die Generierung von Bindings für ausgewählte Servicemethoden zu verhindern, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:internal` für Services und Modelle hinzugefügt, um Typen zuzulassen, die in Go, aber nicht in JS/TS exportiert werden, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für Konstanten eines Aliastyps zum Binding-Generator hinzugefügt, um schwach typisierte Enums zu ermöglichen, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tests des Binding-Generators für Go-1.24-Features hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4068](https://github.com/wailsapp/wails/pull/4068)
- Unterstützung für macOS 15 „Sequoia“ zu `OSInfo.Branding` hinzugefügt, um die Erkennung der Betriebssystemversion zu verbessern, in [#4065](https://github.com/wailsapp/wails/pull/4065)
- Hook `PostShutdown` zum Ausführen von benutzerdefiniertem Code nach Abschluss des Herunterfahrens hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Struct `FatalError` hinzugefügt, um die Erkennung schwerwiegender Fehler in benutzerdefinierten Fehler-Handlern zu unterstützen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Reihenfolge beim Starten und Herunterfahren von Services standardisiert und dokumentiert von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Testgerüst für die Start-/Herunterfahrsequenz der Anwendung sowie Tests zum Starten und Herunterfahren von Services hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Methode `RegisterService` zum Registrieren von Services nach dem Erstellen der Anwendung hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Feld `MarshalError` in den Anwendungs- und Serviceoptionen zur benutzerdefinierten Fehlerbehandlung bei Binding-Aufrufen hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Abbrechbaren Promise-Wrapper hinzugefügt, der Abbruchanforderungen durch Promise-Ketten weiterleitet, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Möglichkeit hinzugefügt, den Abbruch eines Binding-Aufrufs an ein `AbortSignal` zu koppeln, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Unterstützung für `data-wml-*`-Attribute in WML neben den üblichen `wml-*`-Attributen hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Methode `Configure` für alle Services zur nachträglichen Konfiguration und dynamischen Neukonfiguration hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `fileserver` sendet im unkonfigurierten Zustand eine 503-Service-Unavailable-Antwort von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `kvstore` stellt im unkonfigurierten Zustand standardmäßig einen In-Memory-Schlüssel-Wert-Speicher bereit von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Load` für den Service `kvstore` hinzugefügt, um Daten nach Konfigurationsänderungen erneut aus der Datei zu laden, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Clear` für den Service `kvstore` hinzugefügt, um alle Schlüssel zu löschen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Typ `Level` im Service `log` hinzugefügt, um Konstanten für Protokollierungsstufen auf der JS-Seite bereitzustellen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Log` für den Service `log` hinzugefügt, um die Protokollierungsstufe dynamisch festzulegen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `sqlite` stellt im unkonfigurierten Zustand standardmäßig eine In-Memory-Datenbank bereit von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Close` zum manuellen Schließen der Datenbank zum Dienst `sqlite` hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Unterstützung für den Abbruch von Abfragemethoden des Dienstes `sqlite` hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Unterstützung für vorbereitete Anweisungen mit JS-Bindings zum Dienst `sqlite` hinzugefügt von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Gin-Unterstützung von [Lea Anthony](https://github.com/leaanthony) in [PR](https://github.com/wailsapp/wails/pull/3537), basierend auf der ursprünglichen Arbeit von [@AnalogJ](https://github.com/AnalogJ) in diesem [PR](https://github.com/wailsapp/wails/pull/3537)
- Fehler behoben, durch den automatisches Speichern und automatisches Speichern von Passwörtern immer aktiviert waren, von [@oSethoum](https://github.com/osethoum) in [#4134](https://github.com/wailsapp/wails/pull/4134)
- `SetMenu()` zu Fenstern hinzugefügt, um für ein Fenster ein Menü festlegen zu können, von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Benachrichtigungen hinzugefügt von [@popaprozac](https://github.com/popaprozac) in [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Unterstützung für Dateizuordnungen auf dem Mac hinzugefügt von [@wimaha](https://github.com/wimaha) in [#4177](https://github.com/wailsapp/wails/pull/4177)
- `wails3 tool version` zum Erhöhen semantischer Versionsnummern hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Badges unter macOS und Windows hinzugefügt von [@popaprozac](https://github.com/popaprozac) in [#](https://github.com/wailsapp/wails/pull/4234)
- Unterstützung für registrierte/streng typisierte Ereignisse hinzugefügt von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- Möglichkeit zur Registrierung von Hooks für benutzerdefinierte Ereignisse hinzugefügt von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` zum Öffnen des Pfads `path` im Dateimanager des Systems, mit optionaler Hervorhebung über `selectFile`, von [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Neues Flag `-git` für den Befehl `wails3 init` von [@leaanthony](https://github.com/leaanthony)
- Neuer Befehl `wails3 generate webview2bootstrapper` von [@leaanthony](https://github.com/leaanthony)
- Methode `init()` zur manuellen Initialisierung der Runtime zur Runtime hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Option `WindowDidMoveDebounceMS` zu WindowOptions von Window hinzugefügt von [@leaanthony](https://github.com/leaanthony)
- Single-Instance-Funktion hinzugefügt von [@leaanthony](https://github.com/leaanthony). Basierend auf dem [v2-PR](https://github.com/wailsapp/wails/pull/2951) von @APshenkin.
- Befehl `wails3 generate template` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 releasenotes` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 update cli` von [@leaanthony](https://github.com/leaanthony)
- Option `-clean` für den Befehl `wails3 generate bindings` von [@leaanthony](https://github.com/leaanthony)
- AppImage-Builds für Linux auf aarch64 (arm64) ermöglicht von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Hyperlink für Sponsor hinzugefügt von @ansxuman in [#3958](https://github.com/wailsapp/wails/pull/3958)
- Unterstützung für die Linux-Paketierung als deb-, rpm- und Arch-Linux-Pakete von
- Unterstützung für universelle Darwin-Builds und -Pakete hinzugefügt von
- Ereignisdokumentation zur Website hinzugefügt von
- Vorlagen für sveltekit und sveltekit-ts, die für die Entwicklung ohne SSR konfiguriert sind
- Build-Assets mit dem neuen Befehl `wails3 update build-assets` aktualisiert von
- Beispiel zum Testen der HTML Drag and Drop API von
- Unterstützung für Dateizuordnungen von [leaanthony](https://github.com/leaanthony) in
- Neuer Befehl `wails3 generate runtime` von
- Neue Option `InitialPosition`, um anzugeben, ob das Fenster zentriert werden soll oder
- Methoden `Path` und `Paths` zum Paket `application` hinzugefügt von
- Windows-Optionen `GeneralAutofillEnabled` und `PasswordAutosaveEnabled` hinzugefügt
- Möglichkeit hinzugefügt, das Fenster abzurufen, das eine Dienstmethode aufruft, von
- Optionen `EnabledFeatures` und `DisabledFeatures` für WebView2 hinzugefügt von
- ⊞ Neues DIP-System für verbesserte Unterstützung von Monitoren mit hoher DPI von
- ⊞ Option für den Fensterklassennamen von [windom](https://github.com/windom/) in
- Dienste wurden um Plugin-Funktionalität erweitert. Von
- 🐧 Ereignisse WindowDidMove / WindowDidResize in
- ⊞ Ereignis WindowDidResize in
-  Ereignis ApplicationShouldHandleReopen hinzugefügt, um Dock-Ereignisse verarbeiten zu können
-  getPrimaryScreen/getScreens zur Implementierung hinzugefügt von @tmclane in
-  Option zum Anzeigen der Symbolleiste im Vollbildmodus unter macOS hinzugefügt von
- 🐧 onKeyPress-Logik hinzugefügt, um Linux-Tastendrücke in einen Accelerator umzuwandeln
- 🐧 Task `run:linux` hinzugefügt von
- Methode `SetIcon` exportiert von [@almas-x](https://github.com/almas-x) in
- `OnShutdown` verbessert von [@almas-x](https://github.com/almas-x) in
- Methode `ToggleMaximise` in der Schnittstelle `Window` wiederhergestellt von
- Weitere Informationen zu `Environment()` hinzugefügt. Von @leaanthony in
- Methode `WebviewWindow.IsFocused` über die Schnittstelle `Window` verfügbar gemacht von
- Mehrere durch Leerzeichen getrennte auslösende Ereignisse im WML-System unterstützt von
- ESM-Exporte aus dem gebündelten JS-Runtime-Skript hinzugefügt von
- Flag für den Binding-Generator hinzugefügt, um das gebündelte JS-Runtime-Skript zu verwenden anstelle von
- `setIcon` unter Linux implementiert von [@abichinger](https://github.com/abichinger)
- Flag `-port` zum dev-Befehl hinzugefügt und Umgebungsvariable unterstützt
- Tests für Aufrufe gebundener Methoden hinzugefügt von
- ⊞ `SetIgnoreMouseEvents` für ein bereits erstelltes Fenster hinzugefügt von
-  Möglichkeit zum Festlegen der Stapelebene (Reihenfolge) eines Fensters hinzugefügt von

### Behoben

- Halbierte `Screen.Bounds`, `WorkArea` und `Size` auf Macs mit Retina-Display behoben, indem NSScreen-Punktwerte in den `Physical*`-Feldern in Gerätepixel umgerechnet und die übergeordneten Felder `Screen.X`/`Y` befüllt werden, sodass die Erkennung sich berührender Monitore und die Platzierung im Arbeitsbereich in [PR](https://github.com/wailsapp/wails/pull/5168) korrekt funktionieren, von @wayneforrest
- Datenrennen im ScreenManager behoben, das bei einer Änderung der Bildschirmkonfiguration einen WebKit-DisplayLink-Deadlock verursacht (z. B. beim Anschließen eines externen Monitors während des Ruhezustands/Aufwachens)
- CFBundleIconName wird direkt auf appicon gesetzt, wenn Assets.car vorhanden ist, in [PR](https://github.com/wailsapp/wails/pull/5154) von @symball
- Behoben, dass `wails3 doctor` unter Fedora, openSUSE, Arch und NixOS falsche WebKitGTK-Pakete meldete. Die 4.0-Fallback-Einträge wurden entfernt, da v3 die 4.1-API zur Kompilierzeit benötigt (#5071)
- Paketname von webkit2gtk in Doctor für openSUSE korrigiert (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, der korrekte openSUSE-Paketname) (#5071)
- `Unexpected token '<'`-Fehler behoben, wenn `/wails/custom.js` im Desktop-Entwicklungsmodus fehlt. Ein expliziter 404-Handler für `/wails/custom.js` sowie eine Validierung von `Content-Type` ohne Berücksichtigung der Groß-/Kleinschreibung in `loadOptionalScript` wurden hinzugefügt, damit HTML-SPA-Fallbacks nicht als JavaScript eingeschleust werden. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Hervorhebungsstatus des Systemleistenmenüs unter macOS behoben: Das Symbol zeigt nun den ausgewählten Zustand an, wenn das Menü geöffnet ist (#4910)
- Behoben, dass ein an die Systemleiste angeheftetes Fenster unter macOS hinter anderen Fenstern erschien; es verwendet nun die korrekte Fensterebene für Pop-ups (#4910)
- Falsche `@wailsio/runtime`-Importbeispiele in der gesamten Dokumentation korrigiert (#4989)
- Behoben, dass rahmenlose Fenster unter Darwin nicht minimiert werden konnten (#4294)
- 20-30 Minuten dauernde Hänger während `wails3 build` und `wails3 dev` behoben, indem `node_modules/` von der Aktualitätsprüfung von go-task ausgeschlossen wurde. Zuvor führte das Glob-Muster `sources: "**/*"` dazu, dass go-task jede Datei in `node_modules/` auflistete und ihre Prüfsumme berechnete (50000-100000+ Dateien bei umfangreichen Abhängigkeiten wie MUI); dies war unter Windows/NTFS besonders langsam (#4939)
- GTK4-Buildfehler behoben, der durch eine Kollision des C-`Screen`-Typedefs mit X11 Xlib.h verursacht wurde (#4957)
- Konsistenz der Dock-Badge-Methoden unter macOS korrigiert
- Behoben, dass `InvisibleTitleBarHeight` auf alle macOS-Fenster statt nur auf rahmenlose Fenster oder Fenster mit transparenter Titelleiste angewendet wurde (#4960)
- Ruckeln/Zittern des Fensters bei der Größenänderung über die oberen Ecken mit aktiviertem `InvisibleTitleBarHeight` behoben, indem die Einleitung des Ziehvorgangs in der Nähe der Fensterränder übersprungen wird (#4960)
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Behoben, dass Drag-and-drop von Dateien unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Behoben, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Windows aktiviert war
- Behoben, dass Koordinaten beim Ablegen von Dateien unter Windows im falschen Pixelraum angegeben wurden (physische statt CSS-Pixel)
- Behoben, dass Drag-and-drop von Dateien unter Linux mit Hover-Effekten nicht zuverlässig funktionierte
- Behoben, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Linux aktiviert war
- Behoben, dass Fenster beim Ein-/Ausblenden unter Linux/GTK4 manchmal minimiert wiederhergestellt wurden; dazu wird `gtk_window_present()` verwendet (#4957)
- Behoben, dass das Abrufen/Festlegen der Fensterposition unter Linux/GTK4 stets 0,0 zurückgab; dazu wurde bedingte X11-Unterstützung über `XTranslateCoordinates`/`XMoveWindow` hinzugefügt (#4957)
- Behoben, dass die maximale Fenstergröße unter Linux/GTK4 nicht erzwungen wurde; dazu wurde eine signalbasierte Größenbegrenzung als Ersatz für das entfernte `gtk_window_set_geometry_hints` hinzugefügt (#4957)
- DPI-Skalierung unter Linux/GTK4 behoben, indem die korrekte Berechnung von PhysicalBounds und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` implementiert wurden (GTK 4.14+)
- Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4 behoben
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Behoben, dass Drag-and-drop von Dateien unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Behoben, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Windows aktiviert war
- Behoben, dass Koordinaten beim Ablegen von Dateien unter Windows im falschen Pixelraum angegeben wurden (physische statt CSS-Pixel)
- Behoben, dass Drag-and-drop von Dateien unter Linux mit Hover-Effekten nicht zuverlässig funktionierte
- Behoben, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Linux aktiviert war
- DPI-Skalierung unter Linux/GTK4 behoben, indem die korrekte Berechnung von PhysicalBounds und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` implementiert wurden (GTK 4.14+)
- Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4 behoben
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Problem mit „Geisterfenstern“ unter macOS behoben, das dadurch verursacht wurde, dass AppKit-APIs in App.Window.Current() nicht über den Hauptthread aufgerufen wurden (#4947) von @wimaha
- Behoben, dass HTML-`<input type="file">` unter macOS nicht funktionierte; dazu wurde WKUIDelegate runOpenPanelWithParameters implementiert (#4862)
- Behoben, dass natives Drag-and-drop von Dateien bei Verwendung des npm-Moduls `@wailsio/runtime` unter macOS/Linux nicht funktionierte (#4953) von @leaanthony
- Binding-Generierung für paketübergreifende Typaliase korrigiert (#4578) von @fbbdev
- Absturz von OpenFileDialog unter Linux aufgrund einer Verletzung der GTK-Threadsicherheit behoben (#3683) von @ddmoney420
- SIGSEGV-Absturz beim Aufruf von `Focus()` für ein ausgeblendetes oder zerstörtes Fenster behoben (#4890) von @ddmoney420
- Mögliche Panic beim Festlegen eines leeren Symbols oder Bitmaps unter Linux behoben (#4923) von @ddmoney420
- Absturz von ErrorDialog beim Aufruf aus einem Service-Binding unter macOS behoben (#3631) von @leaanthony
- Menüs werden unter Windows in `v3\examples\dialogs` angezeigt, von @ndianabasi
- Race-Condition behoben, die beim Neuladen der Seite einen TypeError verursachte (#4872) von @ddmoney420
- Falsche Ausgabe der Binding-Generator-Tests korrigiert, indem globaler Zustand aus der Methode `Collector.IsVoidAlias()` entfernt wurde (#4941) von @fbbdev
- Behoben, dass die `<input type="file">`-Dateiauswahl unter macOS nicht funktionierte (#4862) von @leaanthony
- Behebt inkonsistente Koordinatensysteme bei `Position()` und `SetPosition()` unter macOS, die beim Speichern und Wiederherstellen des Zustands zu einer Verschiebung der Fensterposition führten (#4816), von @leaanthony
- Behebt den SetProcessDpiAwarenessContext-Fehler „Zugriff verweigert“, wenn die DPI-Awareness bereits über das Anwendungsmanifest festgelegt ist (#4803)
- Aktualisiert die Dokumentationsseite zu Tastenkürzeln und korrigiert den Typ des Callback-Parameters für `KeyBinding.Add`, von @ndianabasi
- Korrigiert die Dokumentation zur Generierung benutzerdefinierter Bindings: Es muss `-d String` statt `-o String` verwendet werden
- Behebt, dass beim Aufruf von `menu.Update()` die untergeordneten Menüelemente nicht entfernt werden
- Korrigiert veraltete Verweise auf die Manager-API in der Dokumentation (31-Dateien wurden auf das neue Muster mit `app.Window.New()`, `app.Event.Emit()` usw. umgestellt), von @leaanthony
- Behebt einen Linux-Absturz bei einer Panic in an JavaScript gebundenen Go-Methoden, der durch das Überschreiben von Signal-Handlern durch WebKit verursacht wurde (#3965), von @leaanthony
- Behebt, dass SaveFileDialog.SetFilename() unter Linux keine Wirkung hatte (#4841), von @samstanier
- Behebt, dass die Ablagekoordinaten im Drag-and-drop-Beispiel als undefined angezeigt wurden
- Behebt, dass die Erstellung eines macOS-App-Bundles fehlschlug, wenn APP_NAME Leerzeichen enthielt (Problem mit der Klammererweiterung)
- Behebt eine Index-out-of-bounds-Panic unter Windows beim Aufruf von Dienstmethoden (Zurücksetzen von goccy/go-json)
- Behebt, dass Drag-and-drop von Dateien unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Behebt, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Windows aktiviert war
- Behebt, dass sich die Koordinaten beim Ablegen von Dateien unter Windows im falschen Pixelraum befanden (physische Pixel statt CSS-Pixel)
- Behebt, dass Drag-and-drop von Dateien unter Linux in Verbindung mit Hover-Effekten nicht zuverlässig funktionierte
- Behebt, dass internes HTML5-Drag-and-drop nicht funktionierte, wenn das Ablegen von Dateien unter Linux aktiviert war
- Aktualisiert alle Befehle in Taskfile.yml-Dateien für sämtliche Betriebssysteme, damit Leerzeichen in Variablen wie `APP_NAME` unterstützt werden, von @ndianabasi
- Behebt einen Fehler bei Befehlsargumenten während der Ausführung des Tasks 'build:universal:lipo:go' unter Linux, von @wux1an
- Behebt den Docker-Fehler „undefined symbol: **<em>ubsan</em>handle_xxxxxxx“ bei der Ausführung von 'wails3 build GOOS=darwin GOARCH=arm64' unter Linux, von @wux1an
- Fasst die Dokumentation zu benutzerdefinierten Protokollen zusammen und ergänzt Abschnitte zu Universal Links, von @leaanthony
- Behebt einen Absturz des Windows-Menüs im Infobereich beim wiederholten Klicken auf das Symbol durch eine Schutzvorkehrung gegen gleichzeitige TrackPopupMenuEx-Aufrufe (#4151), von @leaanthony
- Verhindert einen Anwendungsabsturz, wenn systray.Run() vor app.Run() aufgerufen wird, von @leaanthony
- Behebt einen macOS-Absturz beim Umschalten der Fenstersichtbarkeit über Hide()/Show(), wenn ApplicationShouldTerminateAfterLastWindowClosed aktiviert ist (#4389), von @leaanthony
- Behebt ein Speicherleck in Kontextmenüs unter macOS und Windows beim wiederholten Öffnen der Menüs (#4012), von @leaanthony
- Behebt, dass native Ressourcen von Kontextmenüs unter macOS nicht wiederverwendet wurden, wodurch bei jeder Anzeige ein neues Menü erstellt wurde (#4012), von @leaanthony
- Behebt, dass ein Klick auf das macOS-Dock-Symbol ausgeblendete Fenster nicht anzeigte, wenn die Anwendung mit `Hidden: true` gestartet wurde (#4583), von @leaanthony
- Behebt, dass der macOS-Druckdialog aufgrund eines falschen Fensterzeigertyps im CGO-Aufruf nicht geöffnet wurde (#4290), von @leaanthony
- Behebt einen Absturz des Fenstermenüs unter Wayland, der durch den Zugriff von appmenu-gtk-module auf ein noch nicht realisiertes Fenster verursacht wurde (#4769), von @leaanthony
- Behebt einen Absturz der GTK-Anwendung, wenn der Anwendungsname ungültige Zeichen enthält (Leerzeichen, Klammern usw.), von @leaanthony
- Behebt den Fehler „Nicht genügend Arbeitsspeicher“ bei der Initialisierung von Drag-and-drop unter Windows (#4701), von @overlordtm
- Behebt, dass der Dateiexplorer unter Linux aufgrund fehlerhafter URI-Kodierung das falsche Verzeichnis öffnete (#4397), von @leaanthony
- Behebt AppImage-Buildfehler auf modernen Linux-Distributionen (Arch, Fedora 39+, Ubuntu 24.04+), indem `.relr.dyn`-ELF-Abschnitte automatisch erkannt und das Entfernen von Symbolen deaktiviert wird (#4642), von @leaanthony
- Behebt, dass `wails doctor` WebKit-Pakete auf Fedora- und DNF-basierten Systemen fälschlicherweise als installiert meldete (#4457), von @leaanthony
- Behoben, dass die standardmäßig bereitgestellte Datei `config.yml` `wails3 dev` mit einem Produktions-Build ausführte, von @mbaklor
- Behebt Buildfehler durch iOS-Dienst-Stubs, die ein nicht vorhandenes Paket importierten, von @leaanthony
- Behebt, dass strukturiertes Logging in Debug-/Info-Methoden Fehler des Typs „no formatting directives“ verursachte, von @leaanthony
- Entfernt temporäre Debug-Ausgaben, die versehentlich aus dem Merge der mobilen Plattformen übernommen wurden, von @leaanthony
- Behebt einen WebKitGTK-Absturz unter Wayland mit NVIDIA-GPUs (Fehler 71 Protocol error), indem der DMA-BUF-Renderer automatisch deaktiviert wird, von @leaanthony
- Behebt, dass der Alphawert in `application.WebviewWindowOptions.BackgroundColour` unter Linux ignoriert wurde ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Behebt, dass für das Windows-Symbol im Infobereich nicht standardmäßig das Anwendungssymbol verwendet wurde, wenn kein benutzerdefiniertes Symbol angegeben war (#4704)
- Verfolgt die Eigentümerschaft von `HICON`, sodass nur vom Benutzer erstellte Handles zerstört werden, und verhindert dadurch Abstürze beim Neustart von Explorer (#4653).
- Gibt beim Zerstören den Listener für das Windows-Systemdesign und zurückgehaltene Infobereichssymbole frei, um Ressourcenlecks durch nicht freigegebene Goroutinen und Gerätekontexte zu verhindern (#4653).
- Kürzt Tooltips von Infobereichssymbolen auf 127 UTF-16-Einheiten, um die Beschädigung von Surrogatpaaren und Mehrbyte-Glyphen zu vermeiden (#4653).
- Behebt einen Fehler im Windows-Paket-Task (#4667)
- Korrigiert die Linux-AppImage-Variable appicon im Linux-Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Behebt einen Windows-Buildfehler, der durch eine Signaturänderung in go-webview2 v1.0.22 verursacht wurde (#4513, #4645)
- Korrigiert die Linux-AppImage-Variable appicon im Linux-Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Protokollbereich in der Linux-Datei desktop.tmpl durch Ändern des Vorlagenverweises von `<.Info.Protocol>` zu `<.Protocol>` korrigiert, von @Tolfx in #4510
- Behebt den Fehler durch eine Neudefinition in der Liquid-Glass-Demo in [#4542](https://github.com/wailsapp/wails/pull/4542), von @Etesam913
- Aktualisierung des Menüs im Infobereich unter Linux korrigiert [#4604](https://github.com/wailsapp/wails/issues/4604) von [@JackDoan](https://github.com/JackDoan)
- Behebt, dass unter Windows beim Erstellen eines ausgeblendeten Fensters ein weißes Fenster angezeigt wurde, von @leaanthony in [#4612](https://github.com/wailsapp/wails/pull/4612)
- Korrigiert den Importpfad des Benachrichtigungspakets in der Dokumentation, von @rxliuli in [#4617](https://github.com/wailsapp/wails/pull/4617)
- Fehler behoben, durch den Drag-and-drop bei Verwendung des npm-Pakets @wailsio/runtime nicht funktionierte (#4489), von @leaanthony in #4616
- Windows: Flackern des Fensters beim Start und fälschlicherweise angezeigte ausgeblendete Fenster behoben in [PR](https://github.com/wailsapp/wails/pull/4600) von @leaanthony.
- Probleme beim Maximieren der Fenstergröße unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429), von [@samstanier](https://github.com/samstanier)
- Probleme beim Maximieren der Fenstergröße unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429), von [@samstanier](https://github.com/samstanier)
- Fehler durch erneute Definition in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- Problem behoben, durch das AssetServer unter MacOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576) von @jghiloni
- Kompilierungsproblem beim Build mit NextJs behoben. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Pipelines für Nightly-Releases korrigiert in [#4597](https://github.com/wailsapp/wails/pull/4597) von @riadafridishibly
- Fehler durch erneute Definition in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- Problem behoben, durch das AssetServer unter MacOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576) von @jghiloni
- Kompilierungsproblem beim Build mit NextJs behoben. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Pipelines für Nightly-Releases korrigiert in [#4597](https://github.com/wailsapp/wails/pull/4597) von @riadafridishibly
- Fehler durch erneute Definition in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- SetBackgroundColour unter Windows korrigiert von @PPTGamer in [PR](https://github.com/wailsapp/wails/pull/4492)
- Dokumentation an die Änderungen durch das Manager-API-Refactoring angepasst von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Variable appicon der Linux-.desktop-Datei im Linux-Taskfile korrigiert, [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Dokumentation an die Änderungen durch das Manager-API-Refactoring angepasst von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- In [#4456](https://github.com/wailsapp/wails/issues/4456) gemeldeten Fehler durch Dereferenzierung eines nil-Zeigers unter Windows behoben von @leaanthony in [#4460](https://github.com/wailsapp/wails/pull/4460)
- Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzugefügt, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Problem behoben, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren, von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die anfängliche Untersuchung.
- Fehler behoben, durch den der Vite-Server nach einem fehlgeschlagenen Build nicht bereinigt wurde (#4403)
- Panic beim Schließen oder Abbrechen eines `SaveFileDialog` unter Windows behoben. Behoben in [PR](https://github.com/wailsapp/wails/pull/4284) von @hkhere
- Drag-and-drop auf HTML-Ebene unter Windows behoben von [@mbaklor](https://github.com/mbaklor) in [#4259](https://github.com/wailsapp/wails/pull/4259)
- Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzugefügt, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Problem behoben, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren, von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die anfängliche Untersuchung.
- Fehler behoben, durch den der Vite-Server nach einem fehlgeschlagenen Build nicht bereinigt wurde (#4403)
- Analyse von Benachrichtigungen unter Windows korrigiert von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4450)
- Befehl doctor korrigiert, sodass er Abhängigkeiten des Windows SDK prüft, von [@kodumulo](https://github.com/kodumulo) in [#4390](https://github.com/wailsapp/wails/issues/4390)
- Dereferenzierung eines nil-Zeigers in processURLRequest auf dem Mac behoben von [@etesam913](https://github.com/etesam913) in [#4366](https://github.com/wailsapp/wails/pull/4366)
- Linux-Fehler behoben, der gefilterte Dialoge verhinderte, von [@bh90210](https://github.com/bh90210) in [#4287](https://github.com/wailsapp/wails/pull/4287)
- Probleme mit dem Menü Bearbeiten unter Windows und Linux behoben von [@leaanthony](https://github.com/leaanthony) in [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Mindestsystemversion in macOS-.plist-Dateien von 10.13.0 auf 10.15.0 aktualisiert von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Problem mit übersprungenen Fenster-IDs behoben von [@leaanthony](https://github.com/leaanthony)
- Problem mit einem nil-Menü beim Aufruf von RegisterContextMenu behoben von [@leaanthony](https://github.com/leaanthony)
- Abhängigkeitszyklen in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- Use-before-define-Fehler in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- Build-Flags an den Binding-Generator übergeben von [@fbbdev](https://github.com/fbbdev) in [#4023](https://github.com/wailsapp/wails/pull/4023)
- Pfade im Windows-Taskfile auf Schrägstriche umgestellt, damit es auf Nicht-Windows-Plattformen funktioniert, von [@leaanthony](https://github.com/leaanthony)
- Mac- und Mac-JS-Ereignisse korrigiert von [@leaanthony](https://github.com/leaanthony)
- Ereignis-Deadlock unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- Fehler `Parameter incorrect` bei der Fensterinitialisierung unter Windows behoben, wenn HTML, aber kein JS bereitgestellt wurde, von [@leaanthony](https://github.com/leaanthony)
- Größe des Antwortpräfixes für die Erkennung des Inhaltstyps im Asset-Server korrigiert von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049)
- Verarbeitung von Nicht-404-Antworten im Root-Index-Pfad des Asset-Servers korrigiert von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049)
- Undefiniertes Verhalten im Binding-Generator beim Prüfen von Eigenschaften generischer Typen behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Modelle korrigiert, wenn der zugrunde liegende Typ nicht dieselben Eigenschaften wie der benannte Wrapper besitzt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Map-Schlüsseltypen und die Vorverarbeitung korrigiert von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Structs korrigiert, die Marshaler-Schnittstellen implementieren, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Erkennung von Typzyklen mit generischen Typen im Binding-Generator korrigiert von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ungültige Verweise auf nicht exportierte Modelle in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Injizierten Code an das Ende der Servicedateien verschoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Fehlerbehandlung beim Schließen von Dateien im Binding-Generator behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Warnungen für Services unterdrückt, die Lebenszyklus- oder HTTP-Methoden, aber keine anderen gebundenen Methoden definieren, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Fehler behoben, durch den Nicht-React-Vorlagen bei Verwendung des hellen Systemfarbschemas die „Hello World“-Fußzeile nicht anzeigten, von [@marcus-crane](https://github.com/marcus-crane) in [#4056](https://github.com/wailsapp/wails/pull/4056)
- Ausgeblendete Menüeinträge unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- Verarbeitung und Formatierung von Fehlern in Nachrichtenprozessoren behoben von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Übersprungenes Herunterfahren von Services beim Beenden der Anwendung behoben von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Sichergestellt, dass Menüaktualisierungen im Hauptthread erfolgen, von [@leaanthony](https://github.com/leaanthony)
- Der Mechanismus zum Ziehen und Ändern der Größe ist jetzt robuster und entspricht genauer dem erwarteten Plattformverhalten, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- [#4097](https://github.com/wailsapp/wails/issues/4097) behoben: Webpack/Angular verwirft den Initialisierungscode der Runtime, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Anfangs ausgeblendete Menüeinträge behoben von [@IanVS](https://github.com/IanVS) in [#4116](https://github.com/wailsapp/wails/pull/4116)
- Fehler behoben, durch den assetFileServer bei einer Anfrage ohne Dateierweiterung keine `.html`-Dateien auslieferte, wenn `[request]` nicht existiert, `[request].html` jedoch schon
- Pfade für die Symbolgenerierung korrigiert von [@robin-samuel](https://github.com/robin-samuel) in [#4125](https://github.com/wailsapp/wails/pull/4125)
- Fehler behoben, durch den die Ereignisse `fullscreen`, `unfullscreen`, `unminimise` und `unmaximise` nicht ausgelöst wurden, von [@oSethoum](https://github.com/osethoum) in [#4130](https://github.com/wailsapp/wails/pull/4130)
- NSIS-Fehler aufgrund eines falschen Präfixes bei der Standardversion in der Konfiguration behoben von [@robin-samuel](https://github.com/robin-samuel) in [#4126](https://github.com/wailsapp/wails/pull/4126)
- Fehler behoben, durch den die Runtime-Funktion Dialogs unter Windows maskierte Pfade zurückgab, von [TheGB0077](https://github.com/TheGB0077) in [#4188](https://github.com/wailsapp/wails/pull/4188)
- Erkennungspfad für Webview2 in HKCU korrigiert von [@leaanthony](https://github.com/leaanthony).
- Eingabeproblem unter macOS behoben von [@leaanthony](https://github.com/leaanthony).
- Dateinamen der Taskdatei für die Windows-Symbolgenerierung korrigiert von [@yulesxoxo](https://github.com/yulesxoxo) in [#4219](https://github.com/wailsapp/wails/pull/4219).
- Transparenzproblem bei rahmenlosen Fenstern behoben von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Fokusaufrufe bei deaktiviertem oder minimiertem Fenster korrigiert von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Fehler behoben, durch den Infobereichssymbole nach einem Neustart der Taskleiste nicht angezeigt wurden, von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Fehlende Implementierung von Flush() durch fallbackResponseWriter behoben in [#4245](https://github.com/wailsapp/wails/pull/4245)
- Fehlende Implementierung von Flush() durch fallbackResponseWriter behoben von [@superDingda] in [#4236](https://github.com/wailsapp/wails/issues/4236)
- Absturz beim Schließen eines macOS-Fensters mit einem ausstehenden asynchronen Aufruf einer an Go gebundenen Funktion behoben von [@joshhardy](https://github.com/joshhardy) in [#4354](https://github.com/wailsapp/wails/pull/4354)
- Race-Condition beim Start des Windows-Effizienzmodus behoben von [@leaanthony](https://github.com/leaanthony)
- Bereinigung von Windows-Symbolhandles korrigiert von [@leaanthony](https://github.com/leaanthony).
- `OpenFileManager` unter Windows behoben von [@PPTGamer](https://github.com/PPTGamer) in [#4375](https://github.com/wailsapp/wails/pull/4375).
- Optionen für minimale/maximale Breite unter Linux korrigiert von @atterpac in [#3979](https://github.com/wailsapp/wails/pull/3979)
- Typdefinitionen für Typescript-Vorlagen durch Erhöhung der npm-Version korrigiert von @atterpac in [#3966](https://github.com/wailsapp/wails/pull/3966)
- CSS-Verweis in der Sveltekit-Vorlage korrigiert von @atterpac in [#3945](https://github.com/wailsapp/wails/pull/3945)
- Sichergestellt, dass wichtige Callbacks in window run() im Hauptthread aufgerufen werden, von [@leaanthony](https://github.com/leaanthony)
- Beispiele für die Verzeichnisauswahl in Dialogen korrigiert von [@leaanthony](https://github.com/leaanthony)
- Neue chinesische Fehlerseite für den Fall erstellt, dass index.html fehlt, von [@leaanthony](https://github.com/leaanthony)
-  Sichergestellt, dass der Callback `windowDidBecomeKey` im Hauptthread ausgeführt wird, von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt von [@leaanthony](https://github.com/leaanthony)
-  Logik zum Zerstören von Fenstern verbessert von [@leaanthony](https://github.com/leaanthony)
-  Logik für die Fensterposition bei Anbindung an Infobereichssymbole korrigiert von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt von [@leaanthony](https://github.com/leaanthony)
- Ereignisbehandlung korrigiert von [@leaanthony](https://github.com/leaanthony)
- Logik zum Herunterfahren von Fenstern korrigiert von [@leaanthony](https://github.com/leaanthony)
- Das gemeinsame Taskfile generiert nun standardmäßig Typescript-Bindings für Typescript-Vorlagen, von [@leaanthony](https://github.com/leaanthony)
- Beenden der Anwendung bei einer WM_CLOSE-Nachricht korrigiert, wenn keine Fenster geöffnet sind oder nur ein Infobereichssymbol vorhanden ist, von [@mmalcek](https://github.com/mmalcek) in [#3990](https://github.com/wailsapp/wails/pull/3990)
- garble-Build korrigiert von @5aaee9 in [#3192](https://github.com/wailsapp/wails/pull/3192)
- Windows-NSIS-Builds korrigiert von [@leaanthony](https://github.com/leaanthony)
- Deadlock im Linux-Dialog für Mehrfachauswahlen behoben, verursacht durch nicht geschlossene
- Plattformübergreifende Bereinigung von .syso-Dateien während des Windows-Builds korrigiert von
- Kompilierung des amd64-AppImage korrigiert von @atterpac in
- Aktualisierung der Build-Assets korrigiert von @ansxuman in
- Linux-Implementierung der Systray-Funktionen `OnClick` und `OnRightClick` korrigiert von @atterpac
- Fehler behoben, durch den `AlwaysOnTop` auf dem Mac nicht funktionierte, von
-  Duplikat in `application.NewEditMenu` behoben
- 🐧 aarch64-Kompilierung korrigiert
- ⊞ Menüeinträge für Optionsgruppen korrigiert von
- Fehler beim Erstellen einer ausführbaren .app unter macOS behoben, wenn sich 'name' und 'outputfilename'
- Fehler bei der Verwendung von customEventProcessor im Drag-and-drop-Beispiel behoben von
- 🐧 Durch das Hinzufügen von IgnoreMouseEvents verursachten Linux-Kompilierungsfehler behoben von
- ⊞ Fehler bei der Erzeugung der syso-Symboldatei behoben von
- 🐧 Korrektur für die native Ausführung unter Wayland übernommen aus
- Interne Dienstmethoden nicht binden in
- ⊞ Panic beim Start des Infobereichsymbols behoben in
- Interne Dienstmethoden nicht binden in
- ⊞ Panic beim Start des Infobereichsymbols behoben in
- Umfangreiche Überarbeitung der Menüeinträge und Ereignisbehandlung. Verbessert vorerst hauptsächlich macOS. Von
- Tests nach der Überarbeitung der Plug-ins und Ereignisse korrigiert in
- ⊞ Warnung zu `Failed to unregister class Chrome_WidgetWin_0` behoben. Von
- Modulprobleme
- Nachrichtenübertragung für Größenänderungsereignisse korrigiert von [atterpac](https://github.com/atterpac) in
- 🐧 Fehler bei der Theme-Verarbeitung unter NixOS behoben von
- Projektinstallation über Laufwerksgrenzen hinweg unter Windows korrigiert von
- CSS der React-Vorlage korrigiert, damit die Fußzeile angezeigt wird, von
- Zombieprozesse im Entwicklungsmodus durch Aktualisierung auf die neueste refresh-Version behoben
- Einbindung der WebKit-Datei für AppImage korrigiert von [Atterpac](https://github.com/atterpac)
- Prüfung von apt-Paketen in Doctor korrigiert von [Atterpac](https://github.com/Atterpac) in
- Einfrieren der Anwendung beim Beenden unter Darwin behoben von @5aaee9 in
- Hintergrundfarben der Beispiele unter Windows korrigiert von
- Standardkontextmenüs korrigiert von [mmghv](https://github.com/mmghv) in
- Hexadezimalwerte für Pfeiltasten unter Darwin korrigiert von
- Drag-and-drop unter Windows funktionsfähig gemacht. Hinzugefügt von
- Linux-Fehler in Doctor behoben, der auftrat, wenn geeignete Treiber fehlten
- DPI-Skalierung beim Start unter Windows korrigiert. Geändert von [@almas-x](https://github.com/almas-x) in
- Ersetzungszeile in `go.mod` auf relative Pfade umgestellt. Korrigiert Windows-Pfade mit
- Klickverarbeitung des macOS-Systrays ohne zugeordnetes Fenster korrigiert von
- Fehlschlagenden Windows-Build aufgrund einer unbekannten Option korrigiert von
- Absturz unter Windows beim Linksklick auf das Systray-Symbol behoben, wenn kein
- Falsche baseURL beim zweimaligen Öffnen eines Fensters korrigiert von @5aaee9 in PR
- Reihenfolge der if-Zweige in der Methode `WebviewWindow.Restore` korrigiert von
- `startURL` über mehrere Aufrufe von `GetStartURL` hinweg korrekt berechnet, wenn
- JS-Typ der Struktur `Screen` an ihr Go-Gegenstück angeglichen von
- Methode `WML.Reload` korrigiert, um die ordnungsgemäße Bereinigung registrierter Ereignisse sicherzustellen
- Sofortiges Schließen benutzerdefinierter Kontextmenüs unter Linux behoben von
- Ausgabepfad und Erweiterung der von der Bindung erzeugten Modelldateien korrigiert
- Importpfade der Modelldateien im von der Bindung erzeugten JS-Code korrigiert
- Drag-and-drop auf einigen Linux-Distributionen korrigiert von
- Fehlenden Task unter macOS bei Verwendung von `wails3 task dev` ergänzt von
- nil-map-assignment beim Registrieren von Ereignissen behoben von
- Unmarshaling der Parameter gebundener Methoden korrigiert von
- Verarbeitung mehrerer Rückgabewerte gebundener Methoden korrigiert von
- Doctor-Erkennung einer npm-Installation korrigiert, die nicht mit dem Systempaketmanager installiert wurde
- Fehlende MicrosoftEdgeWebview2Setup.exe ergänzt. Dank an
- Zufälligen Absturz unter Linux aufgrund der Verarbeitung von Fenster-IDs behoben von @leaanthony. Basierend auf
- Absturz von systemTray.setIcon unter Linux behoben von
- Sichergestellt, dass der Fensterrahmen beim ersten Aufruf der Funktion `setFrameless` angewendet wird unter

### Geändert

- **NICHT ABWÄRTSKOMPATIBEL**: Map-Schlüssel in generierten JS-/TS-Bindungen sind jetzt als optional gekennzeichnet, um die Semantik von Go-Maps korrekt abzubilden. Der Zugriff auf Map-Werte in TypeScript gibt jetzt `T | undefined` statt `T` zurück und erfordert daher Nullprüfungen oder Assertions (#4943) von `@fbbdev`
- Verwendung von `Event` gemäß den Änderungen in `@wailsio/runtime` durch `Events` ersetzt und die entsprechenden Funktionsaufrufe in der Dokumentation in `Features/Events/Event System` angepasst von @AbdelhadiSeddar
- `EnabledFeatures`, `DisabledFeatures` und `AdditionalBrowserArgs` aus den fensterspezifischen Optionen nach `Options.Windows` auf Anwendungsebene verschoben (#4559) von @leaanthony
- README für das Beispiel `Drag N Drop` aktualisiert und hervorgehoben, dass `Internal Drag and Drop` in diesem Beispiel demonstriert wird, von @ndianabasi
- Verschiedene Debug-Protokollmeldungen von Info auf Debug umgestellt (von @mbaklor)
- **NICHT ABWÄRTSKOMPATIBEL:** `EnableDragAndDrop` in den Fensteroptionen in `EnableFileDrop` umbenannt
- **NICHT ABWÄRTSKOMPATIBEL:** `DropZoneDetails` im Ereigniskontext in `DropTargetDetails` umbenannt
- **NICHT ABWÄRTSKOMPATIBEL:** Methode `DropZoneDetails()` auf `WindowEventContext` in `DropTargetDetails()` umbenannt
- **NICHT ABWÄRTSKOMPATIBEL:** Ereignis `WindowDropZoneFilesDropped` entfernt; stattdessen `WindowFilesDropped` verwenden
- **BREAKING:** HTML-Attribut von `data-wails-dropzone` in `data-file-drop-target` ändern
- **BREAKING:** CSS-Hover-Klasse von `wails-dropzone-hover` in `file-drop-target-active` ändern
- **BREAKING:** Optionen `DragEffect`, `OnEnterEffect` und `OnOverEffect` aus Windows entfernen (waren Teil des entfernten IDropTarget)
- Für die gesamte JSON-Verarbeitung zur Laufzeit (Methodenbindungen, Ereignisse, WebView-Anfragen, Benachrichtigungen, kvstore) zu goccy/go-json wechseln, wodurch sich die Leistung um 21-63 % verbessert und die Speicherallokationen um 40-60 % sinken
- Layout der BoundMethod-Struktur optimieren und das Flag isVariadic zwischenspeichern, um den Aufwand pro Aufruf zu reduzieren
- Für Methoden mit `<=8` Argumenten einen auf dem Stack allozierten Argumentpuffer verwenden, um Heap-Allokationen zu vermeiden
- Ergebniserfassung bei Methodenaufrufen optimieren, um bei einzelnen Rückgabewerten eine Slice-Allokation zu vermeiden
- sync.Map für den MIME-Typ-Cache verwenden, um die Leistung bei gleichzeitigem Zugriff zu verbessern
- Puffer-Pool zum Lesen des Anfragekörpers beim HTTP-Transport verwenden
- CloseNotify-Kanal bei der Inhaltstyperkennung verzögert allozieren, um die Allokationen pro Anfrage zu reduzieren
- CSS-Debugprotokollierung aus dem Asset-Server entfernen
- Zuordnung der MIME-Typ-Erweiterungen auf über 50 gängige Webformate erweitern (Schriftarten, Audio, Video usw.)
- Dokumentation für die Window-Optionen `X/Y` aktualisieren @ruhuang2001
- Dokumentation zu `Frontend Runtime` um weitere Optionen zum Generieren von Frontend-Bindungen erweitern, von @ndianabasi
- Dokumentationsseite für den Wails-v3-Asset-Server aktualisieren, von @ndianabasi
- **BREAKING**: Dialogfunktionen auf Paketebene entfernen (`application.InfoDialog()`, `application.QuestionDialog()` usw.). Stattdessen den Manager `app.Dialog` verwenden: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Dialogdokumentation an die tatsächliche API anpassen: `app.Dialog.*`, `AddButton()` mit Callbacks (nicht `SetButtons()`), `SetDefaultButton(*Button)` (kein String), `AddFilter()` (nicht `SetFilters()`), `SetFilename()` (nicht `SetDefaultFilename()`) und `app.Dialog.OpenFile().CanChooseDirectories(true)` für die Ordnerauswahl verwenden
- **BREAKING**: Produktions-Builds sind jetzt die Standardeinstellung. Für Entwicklungs-Builds `DEV=true` in den Taskfiles festlegen. Für Beispiele ein neues Projekt generieren, von @leaanthony
- Beim Auslösen eines benutzerdefinierten Ereignisses mit keinem oder einem Datenargument wird der Datenwert direkt dem Feld Data zugewiesen, ohne ihn in ein Slice einzuschließen, von [@fbbdev](https://github.com/fbbdev) in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows-Tray-Icons berücksichtigen jetzt `SystemTray.Show()`/`Hide()`, indem sie `NIS_HIDDEN` umschalten, sodass Apps tatsächlich verschwinden und wieder angezeigt werden können (#4653).
- Die Tray-Registrierung verwendet aufgelöste Symbole erneut, legt `NOTIFYICON_VERSION_4` einmalig fest und aktiviert `NIF_SHOWTIP`, sodass QuickInfos nach einem Neustart des Explorers wiederhergestellt werden (#4653).
- macOS: Für die Fensterzentrierung `visibleFrame` anstelle von `frame` verwenden, um die Bereiche von Menüleiste und Dock auszuschließen
- macOS: Für die Fensterzentrierung `visibleFrame` anstelle von `frame` verwenden, um die Bereiche von Menüleiste und Dock auszuschließen
- Wenn `wails3 update build-assets` mit dem Parameter `-config` ausgeführt wird, werden die über die Parameter `-product*` festgelegten Werte
- `window.NativeWindowHandle()` -> `window.NativeWindow()`, von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- Interne Fensterverwaltung überarbeiten, von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` und `application.WindowNameKey` entfernt (durch `application.WindowKey` ersetzt), von [@leaanthony](https://github.com/leaanthony)
- ContextMenuData gibt jetzt einen String statt any zurück, von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindungen werden Klassenfelder mit Array-Typen fester Länge jetzt mit der erwarteten Länge initialisiert, statt leer zu sein, von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData gibt jetzt einen String statt any zurück, von [@leaanthony](https://github.com/leaanthony)
- `application.NewService` akzeptiert keine Optionen mehr als optionalen Parameter (stattdessen `application.NewServiceWithOptions` verwenden), von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Abhängigkeit `nanoid` entfernt, von [@leaanthony](https://github.com/leaanthony)
- Window-Beispiel für Mica-/Acrylic-/Tabbed-Fensterstile aktualisiert, von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindungen wurden die Modelldateien `internal.js/ts` entfernt; alle Modelle befinden sich jetzt in `models.js/ts`, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindungen werden benannte Typen nie als Aliasse für andere benannte Typen gerendert; das bisherige Verhalten gilt jetzt nur noch für Aliasse, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindungen werden im Klassenmodus Strukturfelder, deren Typ ein Typparameter ist, als optional gekennzeichnet und nie automatisch initialisiert, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint aus den Vorlagen entfernen, von [@IanVS](https://github.com/IanVS) in [#4059](https://github.com/wailsapp/wails/pull/4059)
- Copyright-Jahr auf 2025 aktualisieren, von [@IanVS](https://github.com/IanVS) in [#4037](https://github.com/wailsapp/wails/pull/4037)
- Dokumentation für event.Sender hinzufügen, von [@IanVS](https://github.com/IanVS) in [#4075](https://github.com/wailsapp/wails/pull/4075)
- Unterstützung für Go 1.24, von [@leaanthony](https://github.com/leaanthony)
- `ServiceStartup`-Hooks werden jetzt beim Aufruf von `App.Run` ausgeführt, nicht in `application.New`, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- `ServiceStartup`-Fehler werden jetzt von `App.Run` zurückgegeben, statt den Prozess zu beenden, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Bindungs- und Dialogaufrufe aus JS werden jetzt mit Fehlerobjekten statt mit Strings abgelehnt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Positionierung des Systray-Menüs unter Windows verbessert, von [@leaanthony](https://github.com/leaanthony)
- Die JS-Runtime wurde nach TypeScript portiert, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Die Runtime wird unmittelbar beim Import initialisiert; es ist nicht erforderlich, auf das Laden des Fensters zu warten, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Die Runtime exportiert keine init-Methode mehr. Zur Initialisierung kann ein Import mit Seiteneffekten verwendet werden, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Gebundene Methoden geben jetzt ein `CancellablePromise` zurück, das bei einem Abbruch mit einem `CancelError` abgelehnt wird. Das eigentliche Ergebnis des Aufrufs wird verworfen, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Integrierte Diensttypen heißen jetzt einheitlich `Service`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Funktionen zum Erstellen integrierter Dienste mit Optionen heißen jetzt einheitlich `NewWithConfig`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Die Methode `Select` des Dienstes `sqlite` heißt jetzt `Query`, um mit den Go-APIs übereinzustimmen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Vorlagen: Runtime nach „dependencies“ verschoben und package.json-Dateien organisiert, von [@IanVS](https://github.com/IanVS) in [#4133](https://github.com/wailsapp/wails/pull/4133)
- Erstellt und signiert App-Bundles während der Entwicklung ad hoc, um bestimmte macOS-APIs zu aktivieren, von [@popaprozac](https://github.com/popaprozac) in [#4171](https://github.com/wailsapp/wails/pull/4171)
- Build-Assets in plattformspezifische Verzeichnisse verschoben, von [@leaanthony](https://github.com/leaanthony)
- Taskfiles in plattformspezifische Verzeichnisse verschoben und umbenannt, von [@leaanthony](https://github.com/leaanthony)
- Die Benutzerführung bei fehlendem `index.html` erheblich verbessert, von [@leaanthony](https://github.com/leaanthony)
- [Windows] Leistung beim Minimieren und Wiederherstellen verbessert, von [@leaanthony](https://github.com/leaanthony). Basiert auf dem ursprünglichen [PR](https://github.com/wailsapp/wails/pull/3955) von [562589540](https://github.com/562589540)
- Option `ShouldClose` entfernt (registrieren Sie stattdessen einen Hook für events.Common.WindowClosing), von [@leaanthony](https://github.com/leaanthony)
- [Windows] Flackern beim Öffnen eines Fensters reduziert, von [@leaanthony](https://github.com/leaanthony)
- `Window.Destroy` entfernt, da dies als interne Funktion vorgesehen war, von [@leaanthony](https://github.com/leaanthony)
- `WindowClose`-Ereignisse in `WindowClosing` umbenannt, von [@leaanthony](https://github.com/leaanthony)
- Frontend-Builds verwenden jetzt abhängig vom Build-Typ die Vite-Umgebung „development“ oder „production“, von [@leaanthony](https://github.com/leaanthony)
- Auf go-webview2 v1.19 aktualisiert, von [@leaanthony](https://github.com/leaanthony)
- Sichergestellt, dass der Fork von taskfile verwendet wird, von @leaanthony
- Fork von Taskfile aktualisiert, um Versionsprobleme bei der Installation mit folgendem Befehl zu beheben:
- Fork von Taskfile verwendet, um Versionsprobleme bei der Installation mit folgendem Befehl zu beheben:
- `service.OnStartup` beendet die Anwendung jetzt bei einem Fehler und führt Folgendes aus:
- Nachrichtenübermittlung für Klicks im Infobereich überarbeitet, um Benutzerinteraktionen besser abzubilden, von
- `all:frontend/dist` in die Asset-Einbettung aufgenommen, um Frameworks zu unterstützen, die Folgendes generieren:
- Taskfile überarbeitet, von [leaanthony](https://github.com/leaanthony) in
- Upgrade auf `go-webview2` v1.0.16, von
- Typ `Screen` korrigiert, sodass er `ID` statt `Id` enthält, von
- Wails-Version von `go.mod.tmpl` aktualisiert, um `application.ServiceOptions` zu unterstützen, von
- Ermittlung des Dienstnamens korrigiert, von [windom](https://github.com/windom/) in
- mkdocs serve verwendet jetzt Docker, von [leaanthony](https://github.com/leaanthony)
- Entwicklungskonfiguration in `config.yml` zusammengeführt, von
- Der Infobereich-Dialog verwendet jetzt standardmäßig das Anwendungssymbol, sofern verfügbar (Windows), von
- GPU und Arbeitsspeicher unter macOS werden besser gemeldet, von
- `WebviewGpuIsDisabled` und `EnableFraudulentWebsiteWarnings` entfernt
- Änderung der Events-API: `On`/`Emit` -> Benutzerereignisse, `OnApplicationEvent` ->
- Events-API unter Linux korrigiert, von [TheGB0077](https://github.com/TheGB0077) in
- [CI] Actions verbessert und deren Ausführung auch in Forks ermöglicht und
- `AbsolutePosition()` in `Position()` umbenannt, von
- Linux-WebKit-Abhängigkeit von webkitgtk2-4.0 auf webkit2gtk-4.1 aktualisiert, um
- Das gebündelte JS-Runtime-Skript ist jetzt ein ESM-Modul: Script-Tags, die es importieren,
- Das Paket `@wailsio/runtime` veröffentlicht seine API nicht im `window.wails`
- Das Window-API-Modul `@wailsio/runtime/src/window` stellt jetzt das enthaltende
- Die JS-Window-API wurde an das aktuelle Go-`WebviewWindow` angepasst
- Der Binding-Generator verwendet jetzt standardmäßig Aufrufe anhand der ID. Die CLI-Option `-id`
- Neues Layout für Binding-Code: Ausgabedateien waren zuvor in Ordnern organisiert
- Das Struct-Feld `application.Options.Bind` wurde umbenannt in
- Neue Syntax für das Binden von Diensten: Dienstinstanzen müssen jetzt gekapselt werden in ein
- Spinner in Nicht-Terminal- oder CI-Umgebungen deaktiviert, von

### Entfernt

- **BREAKING**: `EnabledFeatures`, `DisabledFeatures` und `AdditionalLaunchArgs` aus den fensterspezifischen `WindowsWindow`-Optionen entfernt. Verwenden Sie stattdessen die anwendungsweiten Optionen `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` und `Options.Windows.AdditionalBrowserArgs`. Diese Flags gelten global für die gemeinsam genutzte WebView2-Umgebung (#4559), von @leaanthony
- Native Implementierung von `IDropTarget` unter Windows zugunsten eines JavaScript-basierten Ansatzes entfernt (entspricht dem Verhalten von v2)
- Abhängigkeit github.com/wailsapp/mimetype zugunsten einer erweiterten Erweiterungszuordnung und http.DetectContentType aus der Standardbibliothek entfernt; dadurch reduziert sich die Binärgröße um etwa 1.2 MB
- Abhängigkeit gopkg.in/ini.v1 durch Implementierung eines minimalen Parsers für .desktop-Dateien im Linux-Datei-Explorer entfernt; dadurch werden etwa 45 KB eingespart
- samber/lo durch Verwendung des Pakets slices aus der Go-1.21+-Standardbibliothek und minimaler interner Hilfsfunktionen aus dem Laufzeitcode entfernen, wodurch etwa 310 KB eingespart werden
- Debug-printf-Anweisungen aus dem Darwin-Handler für URL-Schemata entfernen (#4834)
- **INKOMPATIBLE ÄNDERUNG**: Das Ereignis `linux:WindowLoadChanged` entfernen; stattdessen mit `linux:WindowLoadFinished` erkennen, wann das Laden der WebView abgeschlossen ist (#3896), von @leaanthony

### Inkompatible Änderungen

- **Überarbeitung der Manager-API**: Die Anwendungs-API wurde zur besseren Codeorganisation und leichteren Auffindbarkeit von einer flachen Struktur in übersichtliche Manager umstrukturiert, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Service-Methoden umbenannt: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, von [@leaanthony](https://github.com/leaanthony)
- Die Methoden `Path` und `Paths` wurden in das Paket `application` verschoben, von [@leaanthony](https://github.com/leaanthony)
- Das Anwendungsmenü ist jetzt ausschließlich unter macOS verfügbar, von [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Hinzugefügt

## Behoben

## v3.0.0-alpha.77 - 2026-04-18

## Behoben

## v3.0.0-alpha.76 - 2026-04-17

## Behoben

## v3.0.0-alpha.75 - 2026-04-16

## Behoben

## v3.0.0-alpha.74 - 2026-03-01

## Hinzugefügt

## Behoben

## v3.0.0-alpha.73 - 2026-02-27

## Behoben

## v3.0.0-alpha.72 - 2026-02-16

## Behoben

## v3.0.0-alpha.71 - 2026-02-10

## Hinzugefügt

## Behoben

## v3.0.0-alpha.70 - 2026-02-09

## Hinzugefügt

## Behoben

## v3.0.0-alpha.69 - 2026-02-08

## Hinzugefügt

## Behoben

## v3.0.0-alpha.68 - 2026-02-07

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.67 - 2026-02-04

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.66 - 2026-02-03

## Hinzugefügt

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.65 - 2026-02-01

## Hinzugefügt

## v3.0.0-alpha.64 - 2026-01-26

## Hinzugefügt

## v3.0.0-alpha.63 - 2026-01-25

## Behoben

## v3.0.0-alpha.62 - 2026-01-22

## Behoben

## v3.0.0-alpha.61 - 2026-01-20

## Behoben

## v3.0.0-alpha.60 - 2026-01-14

## Behoben

## v3.0.0-alpha.59 - 2026-01-11

## Geändert

## v3.0.0-alpha.58 - 2026-01-09

## Behoben

## v3.0.0-alpha.57 - 2026-01-05

## Geändert

## Behoben

## v3.0.0-alpha.56 - 2026-01-04

## Hinzugefügt

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.55 - 2026-01-02

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.54 - 2025-12-29

## Hinzugefügt

## Behoben

## Entfernt

## v3.0.0-alpha.53 - 2025-12-27

## Hinzugefügt

## Behoben

## v3.0.0-alpha.52 - 2025-12-26

## Behoben

## v3.0.0-alpha.51 - 2025-12-23

## Behoben

## v3.0.0-alpha.50 - 2025-12-21

## Geändert

## v3.0.0-alpha.49 - 2025-12-18

## Geändert

## v3.0.0-alpha.48 - 2025-12-16

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.47 - 2025-12-15

## Hinzugefügt

## Behoben

## v3.0.0-alpha.46 - 2025-12-14

## Hinzugefügt

## Entfernt

## v3.0.0-alpha.45 - 2025-12-13

## Hinzugefügt

## Behoben

## v3.0.0-alpha.44 - 2025-12-12

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.43 - 2025-12-11

## Hinzugefügt

## v3.0.0-alpha.42 - 2025-12-10

## Hinzugefügt

## v3.0.0-alpha.41 - 2025-11-23

## Behoben

## v3.0.0-alpha.40 - 2025-11-13

## Behoben

## v3.0.0-alpha.39 - 2025-11-12

## Hinzugefügt

## Geändert

## v3.0.0-alpha.38 - 2025-11-04

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.37 - 2025-11-02

## Behoben

## v3.0.0-alpha.36 - 2025-10-15

## Behoben

## v3.0.0-alpha.35 - 2025-10-14

## Behoben

## v3.0.0-alpha.34 - 2025-10-06

## Hinzugefügt

## Behoben

## v3.0.0-alpha.33 - 2025-10-04

## Behoben

## v3.0.0-alpha.32 - 2025-10-02

## Behoben

## v3.0.0-alpha.31 - 2025-09-27

## Behoben

## v3.0.0-alpha.30 - 2025-09-26

## Behoben

## v3.0.0-alpha.29 - 2025-09-25

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.29 - 2025-09-25

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.27 - 2025-09-07

## Behoben

## v3.0.0-alpha.26 - 2025-08-24

## Hinzugefügt

## v3.0.0-alpha.25 - 2025-08-16

## Geändert

werden nicht mehr ignoriert und überschreiben den Konfigurationswert.

## v3.0.0-alpha.24 - 2025-08-13

## Hinzugefügt

## v3.0.0-alpha.23 - 2025-08-11

## Behoben

## v3.0.0-alpha.22 - 2025-08-10

## Hinzugefügt

## Geändert

+ Zu weit gefasste Linux-Paketabhängigkeiten und veraltete RPM-Abhängigkeiten behoben.

## v3.0.0-alpha.21 - 2025-08-07

## Behoben

## v3.0.0-alpha.20 - 2025-08-06

## Behoben

## v3.0.0-alpha.19 - 2025-08-05

## Hinzugefügt

## Behoben

## v3.0.0-alpha.18 - 2025-08-03

## Hinzugefügt

## Behoben

## v3.0.0-alpha.17 - 2025-07-31

## Behoben

## v3.0.0-alpha.16 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.15 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.14 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.12 - 2025-07-15

### Hinzugefügt

### Behoben

## v3.0.0-alpha.11 - 2025-07-12

## Hinzugefügt

## v3.0.0-alpha.10 - 2025-07-06

### Inkompatible Änderungen

### Hinzugefügt

### Behoben

### Geändert

## v3.0.0-alpha.9 - 2025-01-13

### Hinzugefügt

### Behoben

### Geändert

## v3.0.0-alpha.8.3 - 2024-12-07

### Geändert

## v3.0.0-alpha.8.2 - 2024-12-07

### Geändert

`go install` von @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Geändert

`go install` von @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Hinzugefügt

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   von [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   von [leaanthony](https://github.com/leaanthony) an der angegebenen X-/Y-Position platziert in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) und   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   von [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Geändert

`service.OnShutdown`für alle zuvor gestarteten Dienste, von @atterpac in   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac in [#3907](https://github.com/wailsapp/wails/pull/3907)   Unterordner von @atterpac in   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) in   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) in   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (durch die Optionen `EnabledFeatures` und `DisabledFeatures` ersetzt) von   [leaanthony](https://github.com/leaanthony)

### Behoben

Kanalvariable von @michael-freling in   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) in   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   in [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) in   [#3841](https://github.com/wailsapp/wails/pull/3841)   `PasteAndMatchStyle`-Rolle im Menü „Bearbeiten“ unter Darwin von   [johnmccabe](https://github.com/johnmccabe) in   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) in   [#3854](https://github.com/wailsapp/wails/pull/3854) von   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   unterscheiden sich. Von @nickisworking in   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Hinzugefügt

[mmghv](https://github.com/mmghv) in   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) und   [leaanthony](https://github.com/leaanthony) in   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Geändert

Anwendungsereignisse `OnWindowEvent` -> Fensterereignisse, von   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   Branches mit dem Präfix `v3/` oder `v3-` von   [stendler](https://github.com/stendler) in   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Behoben

[etesam913](https://github.com/etesam913) in   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) in   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) in   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) in   [#3614](https://github.com/wailsapp/wails/pull/3614) von   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) von   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) von   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) von   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) von   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) von   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Behoben

## v3.0.0-alpha.5 - 2024-07-30

### Hinzugefügt

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   Symbolklick von @5aaee9 in [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) in   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   in [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) in   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) in   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   basierend auf [PR](https://github.com/wailsapp/wails/pull/2044) von @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   das npm-Paket von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   in [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` von [@abichinger](https://github.com/abichinger) in   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) in   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) in   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Behoben

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) in   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) in   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) in   [#3477](https://github.com/wailsapp/wails/pull/3477)   von [Atterpac](https://github.com/atterpac) in   [#3320](https://github.com/wailsapp/wails/pull/3320).   in [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) in   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) in   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) in   [PR](https://github.com/wailsapp/wails/pull/3039)   installiert. Hinzugefügt von [@pylotlight](https://github.com/pylotlight) in   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   Leerzeichen – @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) in PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) in PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   angehängtes Fenster, [tw1nk](https://github.com/tw1nk) in PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) in   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` ist vorhanden.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Listener von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) in   [#3330](https://github.com/wailsapp/wails/pull/3330)   Generator von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   Generator von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) in   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) in   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) in   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   von [@pekim](https://github.com/pekim) in   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) von   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) in   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows von [@bruxaodev](https://github.com/bruxaodev/) in   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Geändert

[mmghv](https://github.com/mmghv) in   [#3611](https://github.com/wailsapp/wails/pull/3611)   Unterstützung für Ubuntu 24.04 LTS von [atterpac](https://github.com/atterpac) in   [#3461](https://github.com/wailsapp/wails/pull/3461)   muss das Attribut `type="module"` aufweisen. Von   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Objekt und startet das WML-System nicht. Diese Änderung verbessert die   Kapselung. Das WML-System kann bei Bedarf manuell durch Aufruf der neuen   Methode `WML.Enable` gestartet werden. Das gebündelte JS-Laufzeitskript führt weiterhin beide   Vorgänge automatisch aus. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Fensterobjekt als Standardexport. Einzelne Methoden können nicht mehr   über benannte ESM-Importe oder ESM-Namespace-Importe importiert werden.   API. Einige Methoden haben einen anderen Namen oder Prototyp, insbesondere: Aus `Screen`   wird `GetScreen`; aus `GetZoomLevel`/`SetZoomLevel` wird `GetZoom`/`SetZoom`;   `GetZoom`, `Width` und `Height` geben Werte jetzt direkt zurück, statt sie   in Objekte einzuschließen. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   wurde entfernt. Verwenden Sie die CLI-Option `-names`, um wieder Aufrufe nach Namen zu verwenden.   Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   wurden nach ihrem enthaltenden Paket benannt; jetzt werden vollständige Go-Importpfade   einschließlich des Modulpfads verwendet. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   Aufruf von `application.NewService`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) in   [#3574](https://github.com/wailsapp/wails/pull/3574)
