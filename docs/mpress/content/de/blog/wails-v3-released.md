---
title: "Wails v3.0.0 veröffentlicht"
description: "Wails v3.0.0 ist da: mehrere Fenster, Taskleistenintegration, native Menüs, Dienste und generierte Bindings, ein Updater, Unterstützung für Mobilgeräte und eine neue Dokumentationswebsite."
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Release-Film zu Wails v3.0.0 – am besten im Vollbildmodus ansehen" width="1920" height="1080"}

## Es ist da!

Heute veröffentliche ich Wails v3.0.0.

Ich will ehrlich sein: Es war ein langer Weg. Das erste Alpha-Tag für v3 wurde im Januar 2023 veröffentlicht. Seitdem haben sehr viele Menschen Dinge entwickelt, getestet, beschädigt und repariert. Ich bin euch allen wirklich dankbar. Wenn du die ungeschönte Fassung dieser Reise lesen möchtest, findest du sie in der [v3Beta-Ankündigung](/blog/wails-v3-beta/).

## Was *ist* Wails?

Falls du neu hier bist: Mit Wails können Go-Programmierer Desktopanwendungen mit den ihnen bereits vertrauten Web-Frontend-Werkzeugen erstellen. Dabei verwendet Wails statt eines eingebetteten Browsers die native WebView der jeweiligen Plattform. Du schreibst das Backend in Go und die Benutzeroberfläche in HTML, CSS und JavaScript (oder einem beliebigen Framework); Wails bündelt alles zu einer einzigen nativen Anwendung.

## Warum v3?

Mit Wails v2 ließ sich unkompliziert eine Go-Anwendung mit einem modernen Web-Frontend erstellen, und es hat sich für viele Anwendungen bewährt. Die auf ein einzelnes Fenster und einen Kontext ausgerichtete Laufzeitumgebung erschwerte jedoch einige alltägliche Aufgaben bei Desktopanwendungen unnötig.

v3 basiert auf einem anderen Modell. Anwendungen, Fenster, Dienste, Ereignisse und Plattformfunktionen sind explizite Objekte. Du erstellst eine Anwendung, registrierst Dienste und erstellst Fenster. Anschließend kommunizierst du mit dem Objekt, dem das benötigte Verhalten zugeordnet ist. So bleibt die Anwendung mit zunehmendem Umfang leichter nachvollziehbar, und Funktionen wie mehrere Fenster sind ein regulärer Teil des Modells statt einer Behelfslösung.

## Was bietet v3.0.0?

- **Mehrere Fenster.** Erstelle, verwalte und schließe Fenster zur Laufzeit – jeweils mit eigener Konfiguration und eigenem Lebenszyklus.
- **Taskleistenintegration.** Verleihe deiner Anwendung eine dauerhafte Präsenz – mit Menüs, hellen und dunklen Symbolen, Vorlagensymbolen unter macOS sowie der Möglichkeit, ein Fenster an das Taskleistensymbol anzuhängen.
- **Native Menüs und Dialogfelder.** Anwendungsmenüs, Kontextmenüs und native Dialogfelder auf jeder Plattform.
- **Dienste und generierte Bindings.** Dein Backend bleibt gewöhnlicher Go-Code. Bindings werden durch statische Quellcodeanalyse mit einem einzigen `wails3 generate bindings`-Befehl generiert. Dadurch werden Kommentare und aussagekräftige Parameternamen ins Frontend übernommen.
- **Ein transparentes Buildsystem.** Projekte verwenden statt eines undurchsichtigen Befehls einen auf [Taskfile](https://taskfile.dev) basierenden Build, den du untersuchen, erweitern und debuggen kannst.
- **Plattformübergreifend.** Windows (amd64 und arm64), macOS (Intel und Apple Silicon) sowie Linux (amd64 und arm64), wobei GTK4 und WebKitGTK 6.0 unter Linux standardmäßig verwendet werden. Der bisherige GTK3-Stack ist bis v3.1 weiterhin über `-tags gtk3` verfügbar.
- **iOS und Android.** Derselbe `main.go` und dasselbe Frontend können auf Mobilgeräten ausgeführt werden. Die Unterstützung für Mobilgeräte ist in der Dokumentation weiterhin als experimentell gekennzeichnet. Betrachte sie daher als etwas, das du ausprobieren und zu dem du Feedback geben kannst, nicht als stabile Zusage.
- **Mit Go 1.26 erstellt.** Das Modul v3 wird mit Go 1.26 erstellt.
- **Ein anwendungsinterner Updater.** Austauschbare Anbieter, kryptografische Verifizierung, atomarer Austausch und eine standardmäßige Benutzeroberfläche, die du gestalten oder ersetzen kannst. Weitere Informationen findest du im [Updater-Leitfaden](/guides/updater/).
- **Paketierung und Signierung.** Leitfäden zur [Codesignierung](/guides/build/signing/) auf jeder Plattform, zu [Installationsprogrammen](/guides/installers/) und zur [MSIX-Paketierung](/guides/build/msix/).
- **Server-Builds.** Führe dieselbe Anwendung und dieselben Dienste ohne natives Fenster als HTTP-Server aus.

## Highlights seit der Betaversion

Der Beitrag zur Betaversion erschien am 2. August, und seitdem hat sich im Changelog viel getan. Einige Punkte möchte ich besonders hervorheben:

- **Eine neue Dokumentationswebsite.** Die Dokumentation basiert jetzt auf unserem eigenen Dokumentationswerkzeug [M-Press](https://m-press.me). Die Übersetzungen wurden fertiggestellt und aus der Community übernommen. Die englische Dokumentation bleibt die maßgebliche Quelle.
- **Eine robustere Nutzung unter Windows.** Wails stellt den Betrieb nach Ausfällen von WebView2-Prozessen jetzt wieder her, statt ein leeres Fenster anzuzeigen. Außerdem wurden mehrere Probleme mit der Lebensdauer von WebView2-COM-Objekten behoben.
- **Globale Tastenkürzel, die deine Tastaturbelegung berücksichtigen.** Unter macOS richten sich Buchstaben-Tastenkürzel bei Belegungen wie AZERTY, QWERTZ und Dvorak jetzt nach der Beschriftung der Taste.
- **Härtung des Updaters.** Unter anderem wurden der Austausch von Binärdateien über Dateisystemgrenzen hinweg unter Linux und AppImage-Updates korrigiert.
- **Einrichtung und Werkzeuge.** Der Assistent `wails3 setup` wurde fertiggestellt, und `wails3 dev` beendet jetzt seine Hintergrundprozesse, wenn deine Anwendung geschlossen wird.
- **Natives Liquid Glass unter macOS.** Unterstützung für den Liquid-Glass-Effekt ab macOS 26, mit einer Ausweichlösung für ältere Versionen.

## Migration von v2

v3 ist eine neue Hauptversion. Die Migration ist eine echte Portierung und nicht nur eine Aktualisierung der Versionsnummer. Die größten Änderungen betreffen den Lebenszyklus von Anwendungen und Fenstern, Dienste anstelle kontextgebundener Bindings, direkte Anwendungs- und Fenster-APIs anstelle des Laufzeitpakets v2 sowie neu generierte Frontend-Bindings.

Der [v2-zu-v3-Migrationsleitfaden](/migration/v2-to-v3/) erläutert all diese Änderungen anhand von Vorher-nachher-Beispielen. Portiere deine Anwendung mit Bedacht, teste währenddessen und behalte v2 bei, bis deine v3-Anwendung bereit ist.

**Wails v2 wird nicht eingestellt.** Es erhält weiterhin kritische Aktualisierungen. Verwende für neue Projekte v3.0.0. Die beiden CLIs sind separate Binärdateien (`wails` und `wails3`) mit unterschiedlichen Modulimportpfaden, sodass du beide auf demselben Computer installieren kannst.

## Installation

Installiere die CLI für v3.0.0 mit:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Führe anschließend den geführten Einrichtungsassistenten aus. Er überprüft deine Umgebung und hilft dir bei den von Wails benötigten Abhängigkeiten:

```sh
wails3 setup
```

![Der Einrichtungsassistent von Wails v3 im dunklen Modus](/assets/screenshots/wails3-setup-wizard-dark.png)

Erstelle dann dein erstes Projekt:

```sh
wails3 init
```

- Dokumentation: <https://v3.wails.io/>
- Migrationsleitfaden: <https://v3.wails.io/migration/v2-to-v3/>
- Versionshinweise: [Wails v3.0.0 auf GitHub](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## Support und Releases

Wails ist weit über meine Erwartungen hinausgewachsen, und es war schwierig, mit der Menge an Issues und Pull Requests Schritt zu halten. Releases erfolgten mitunter ohne festen Plan, und Fehlerbehebungen sowie Beiträge brauchten zu lange, um die Benutzer zu erreichen. Wir haben Systeme eingerichtet, die einen größeren Teil der Grundlagen für Triage und Release-Vorbereitung automatisieren. So können wir schneller auf Probleme reagieren und mehr Zeit darauf verwenden, Verbesserungen und neue Funktionen auszuliefern.

Nach v3.0.0 wechseln wir zu **nächtlichen Vorabversionen** mit datierten Tags wie `v3.1.0-nightly.20261026` sowie zu **wöchentlichen Minor-Releases** wie `v3.1.0` und `v3.2.0`. Durch die nächtlichen Builds sind aktuelle Fehlerbehebungen früher verfügbar, während die wöchentlichen Releases einen regelmäßigen Rhythmus getesteter Aktualisierungen bieten.

Wenn du einen reproduzierbaren Fehler findest, melde ihn bitte mit der Ausgabe von `wails3 doctor` und einem Minimalbeispiel. Wenn du eine neue Funktion oder eine Änderung am öffentlichen Verhalten vorschlagen möchtest, erstelle statt eines Feature-Request-Issues einen PR mit einem WEP-Entwurf (Wails Enhancement Proposal).

## Übersetzungen

Allen, die bei der Übersetzung der Wails-v2-Dokumentation geholfen haben: Von Herzen vielen Dank. Ihr habt eure Zeit investiert, damit Entwickler die Dokumentation in ihrer eigenen Sprache lesen konnten, und das war für viele Menschen von großer Bedeutung. Dafür sind wir zutiefst dankbar.

Crowdin erwies sich für uns als besonders schwierig zu verwalten, weshalb wir es nicht mehr verwenden. Ab v3 wird die Dokumentation automatisch übersetzt, und die Übersetzungen befinden sich im Repository direkt neben den englischen Seiten.

Wenn ihr eine Übersetzung in einer Sprache möchtet, die wir noch nicht anbieten, oder in einer vorhandenen Übersetzung einen Fehler entdeckt habt, eröffnet bitte ein Issue und teilt es uns mit. Wir freuen uns, von euch zu hören, und sind für jede Meldung dankbar.

## Würdigung und Dank

Wails v3 gibt es dank der Menschen, die unfertige Builds getestet, schwerwiegende Fehler gemeldet, die Dokumentation übersetzt, Fragen beantwortet und Code beigetragen haben. Vielen Dank. Die vollständige Liste der Mitwirkenden findet ihr auf der [Seite mit den Danksagungen](/credits/). Ich möchte euch ans Herz legen, sie euch anzusehen.

Zunächst ein riesiges Dankeschön an den Wails-Maintainer [Atterpac](https://github.com/atterpac) für die Entwicklungswerkzeuge, die Linux-Desktop-Integration und die Unterstützung der Community, die dazu beigetragen haben, v3 bis an diesen Punkt zu bringen.

Besonderer Dank gilt [Fabio Massaioli (fbbdev)](https://github.com/fbbdev) für die Entwicklung eines so beeindruckenden statischen Quellcodeanalysators und Binding-Generators. Er trägt maßgeblich dazu bei, dass v3 funktioniert: Aus gewöhnlichen Go-Diensten werden JavaScript- und TypeScript-APIs, wobei Typen und Dokumentation bis in euer Frontend übernommen werden.

Unser Dank gilt außerdem [Travis McLane](https://github.com/tmclane) für seine umfangreiche frühe Arbeit an der Plattform, [Stffabi](https://github.com/stffabi) für den Asset-Server und die Grundlagen der nativen Plattformen, [Zach Botterman (popaprozac)](https://github.com/popaprozac) für Benachrichtigungen, die Dock-Integration und weitere Desktop-Funktionen, [Ian VanSchooten (IanVS)](https://github.com/IanVS) für typisierte Ereignisse und Verbesserungen an der Laufzeit sowie [Lyimmi](https://github.com/Lyimmi) für seine Arbeit an Linux. Diese Beiträge sind weiterhin Bestandteil von v3 – zusammen mit der Arbeit aller, die bei der Entwicklung und Weiterentwicklung geholfen haben.

Ein ganz besonders herzliches Dankeschön an die Sponsoren, die Wails durch diesen langen Übergang begleitet haben. Eure Unterstützung gab uns die Zeit, an der Architektur, den Werkzeugen, den Tests und der Dokumentation zu arbeiten, die uns bis hierher gebracht haben. Besonderer Dank gilt außerdem JetBrains und ZSA sowie [CodeRabbit](https://www.coderabbit.ai/oss) für die Unterstützung des Projekts durch Code-Reviews.

Abschließend gilt Valerio Tesei ein besonderer Dank für seine langjährige Unterstützung des Projekts.

## Schlusswort

Es wartet bereits viel Arbeit auf uns. In der Technologiewelt gibt es einige spannende Entwicklungen, und wir sind bereit, neue Paradigmen zu erkunden und Wails um völlig neue Funktionen zu erweitern. Zunächst muss v3 jedoch ein stabiles Fundament sein, auf das ihr euch verlassen könnt. Das hat Vorrang vor der Einführung der nächsten Neuerungen.

Ich bin wirklich stolz auf v3. Ich hoffe, die Entwicklung damit bereitet euch ebenso viel Freude, wie uns die Entwicklung daran bereitet hat. Probiert es aus, entwickelt etwas damit und sagt uns, was ihr davon haltet. Schaut auf [Discord](https://discord.gg/JDdSxwjhGf), in den [GitHub Discussions](https://github.com/wailsapp/wails/discussions) oder im [Subreddit](https://www.reddit.com/r/wails/) vorbei und sagt Hallo.

Viel Spaß!

&dash;&nbsp;Lea

PS: Wenn ihr oder euer Unternehmen Wails nützlich findet, erwägt bitte, [das Projekt zu sponsern](https://github.com/sponsors/leaanthony). Vielen Dank!
