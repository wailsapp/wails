---
title: "Wails v2 Beta für MacOS"
description: "Versionshinweise und Ankündigungen zu Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![Screenshot von wails-mac](/assets/blog-images/wails-mac.webp)

Heute erscheint die erste Betaversion von Wails v2 für Mac! Es hat eine ganze Weile gedauert, diesen Punkt zu erreichen, und ich hoffe, dass die heutige Version bereits einigermaßen nützlich für Sie ist. Bis hierhin gab es einige Wendungen, und ich hoffe, mit Ihrer Hilfe die verbliebenen Unebenheiten auszubügeln und die Mac-Portierung bis zur endgültigen Veröffentlichung von v2 auf Hochglanz zu bringen.

Das heißt, diese Version ist noch nicht produktionsreif? Für Ihren Anwendungsfall ist sie möglicherweise durchaus bereit, aber es gibt noch einige bekannte Probleme. Behalten Sie daher [dieses Projektboard](https://github.com/wailsapp/wails/projects/7) im Auge. Wenn Sie mitwirken möchten, sind Sie herzlich willkommen!

Was ist also bei Wails v2 für Mac gegenüber v1 neu? Kleiner Hinweis: Es ähnelt stark der Windows-Betaversion :wink:

## Neue Funktionen

![Screenshot von wails-menus-mac](/assets/blog-images/wails-menus-mac.webp)

Native Menüs wurden häufig gewünscht. Wails bietet nun endlich entsprechende Unterstützung. Anwendungsmenüs sind jetzt verfügbar und unterstützen die meisten nativen Menüfunktionen. Dazu gehören Standardmenüeinträge, Kontrollkästchen, Optionsgruppen, Untermenüs und Trennlinien.

In v1 gab es sehr viele Anfragen nach mehr Kontrolle über das Fenster selbst. Ich freue mich, dafür eigens neue Runtime-APIs ankündigen zu können. Sie bieten zahlreiche Funktionen und unterstützen Konfigurationen mit mehreren Monitoren. Außerdem gibt es eine verbesserte Dialog-API: Sie können jetzt moderne, native Dialoge umfassend konfigurieren und damit alle Anforderungen an Ihre Dialoge abdecken.

### Mac-spezifische Optionen

Zusätzlich zu den üblichen Anwendungsoptionen bietet Wails v2 für Mac einige Mac-spezifische Extras:

- Gestalten Sie Ihr Fenster ausgefallen und halbtransparent – wie all die hübschen Swift-Apps!
- Umfangreich anpassbare Titelleiste
- Wir unterstützen die NSAppearance-Optionen für die Anwendung
- Einfache Konfiguration zum automatischen Erstellen eines „Über“-Menüs

### Keine Bündelung von Assets erforderlich

Ein großer Nachteil von v1 war, dass die gesamte Anwendung auf jeweils eine einzige JS- und CSS-Datei reduziert werden musste. Ich freue mich, ankündigen zu können, dass Assets für v2 in keiner Weise gebündelt werden müssen. Möchten Sie ein lokales Bild laden? Verwenden Sie ein `<img>`-Tag mit einem lokalen src-Pfad. Möchten Sie eine tolle Schriftart verwenden? Kopieren Sie sie hinein und tragen Sie den zugehörigen Pfad in Ihr CSS ein.

> Wow, das klingt nach einem Webserver …

Ja, es funktioniert genau wie ein Webserver, ist aber keiner.

> Wie binde ich also meine Assets ein?

Übergeben Sie Ihrer Anwendungskonfiguration einfach ein einzelnes `embed.FS`, das alle Ihre Assets enthält. Sie müssen sich nicht einmal im obersten Verzeichnis befinden – Wails findet selbst heraus, wo sie liegen.

### Neue Entwicklungserfahrung

Da Assets nun nicht mehr gebündelt werden müssen, eröffnet sich eine völlig neue Entwicklungserfahrung. Der neue Befehl `wails dev` erstellt und startet Ihre Anwendung, lädt die Assets jedoch direkt vom Datenträger, anstatt die Assets im `embed.FS` zu verwenden.

Darüber hinaus bietet er folgende Funktionen:

- Hot Reload – Jede Änderung an Frontend-Assets löst ein automatisches Neuladen des Anwendungs-Frontends aus
- Automatischer Neuaufbau – Jede Änderung an Ihrem Go-Code führt dazu, dass Ihre Anwendung neu erstellt und gestartet wird

Zusätzlich wird auf Port 34115 ein Webserver gestartet. Dieser stellt Ihre Anwendung jedem Browser bereit, der eine Verbindung zu ihm herstellt. Alle verbundenen Webbrowser reagieren auf Systemereignisse wie Hot Reload nach einer Asset-Änderung.

In Go arbeiten wir in unseren Anwendungen üblicherweise mit Structs. Häufig ist es sinnvoll, Structs an unser Frontend zu senden und dort als Anwendungszustand zu verwenden. In v1 war dies ein weitgehend manueller Vorgang und eine gewisse Belastung für Entwickler. Ich freue mich, ankündigen zu können, dass in v2 jede im Entwicklungsmodus ausgeführte Anwendung automatisch TypeScript-Modelle für alle Structs erzeugt, die als Eingabe- oder Ausgabeparameter gebundener Methoden verwendet werden. Dies ermöglicht einen nahtlosen Austausch von Datenmodellen zwischen beiden Welten.

Zusätzlich wird dynamisch ein weiteres JS-Modul erzeugt, das alle Ihre gebundenen Methoden umschließt. Es stellt JSDoc für Ihre Methoden bereit und ermöglicht damit Codevervollständigung und Hinweise in Ihrer IDE. Es ist wirklich großartig, wenn Datenmodelle beim Drücken der Tabulatortaste in einem automatisch generierten Modul, das Ihren Go-Code umschließt, automatisch importiert werden!

### Remote-Vorlagen

![Screenshot von remote-mac](/assets/blog-images/remote-mac.webp)

Eine Anwendung schnell zum Laufen zu bringen, war für das Wails-Projekt schon immer ein wichtiges Ziel. Bei der Einführung versuchten wir, viele der damals modernen Frameworks abzudecken: React, Vue und Angular. Die Welt der Frontend-Entwicklung ist stark von unterschiedlichen Meinungen geprägt, entwickelt sich rasant und lässt sich nur schwer vollständig überblicken! Daher stellten wir fest, dass unsere Basisvorlagen ziemlich schnell veralteten, was erheblichen Wartungsaufwand verursachte. Außerdem hatten wir dadurch keine attraktiven modernen Vorlagen für die neuesten und besten Technologie-Stacks.

Mit v2 wollte ich der Community mehr Möglichkeiten geben, indem Sie Vorlagen selbst erstellen und hosten können, anstatt auf das Wails-Projekt angewiesen zu sein. Sie können Projekte jetzt also mit von der Community betreuten Vorlagen erstellen! Ich hoffe, dass dies Entwickler dazu inspiriert, ein lebendiges Ökosystem aus Projektvorlagen aufzubauen. Ich bin wirklich gespannt darauf, was unsere Entwickler-Community erschaffen wird!

### Native M1-Unterstützung

Dank der großartigen Unterstützung von [Mat Ryer](https://github.com/matryer/) unterstützt das Wails-Projekt nun native Builds für M1:

![Screenshot von build-darwin-arm](/assets/blog-images/build-darwin-arm.webp)

Sie können auch `darwin/amd64` als Ziel angeben:

![Screenshot von build-darwin-amd](/assets/blog-images/build-darwin-amd.webp)

Oh, fast hätte ich es vergessen … Sie können auch `darwin/universal` verwenden … :wink:

![Screenshot von build-darwin-universal](/assets/blog-images/build-darwin-universal.webp)

### Cross-Compilation für Windows

Da Wails v2 für Windows vollständig in Go implementiert ist, können Sie Builds für Windows ohne Docker erstellen.

![Screenshot von build-cross-windows](/assets/blog-images/build-cross-windows.webp)  
bu

### WKWebView-Renderer

V1 basierte auf einer inzwischen veralteten WebView-Komponente. V2 verwendet die neueste WKWebKit-Komponente, sodass Sie die neuesten und besten Funktionen von Apple erwarten können.

### Fazit

Wie bereits in den Versionshinweisen für Windows erwähnt, schafft Wails v2 eine neue Grundlage für das Projekt. Ziel dieser Version ist es, Feedback zum neuen Ansatz einzuholen und vor der vollständigen Veröffentlichung sämtliche Fehler zu beheben. Ihr Feedback ist herzlich willkommen! Bitte richten Sie es an das Diskussionsforum [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Abschließend möchte ich allen [Projektsponsoren](/credits/#sponsors) meinen besonderen Dank aussprechen, darunter [JetBrains](https://www.jetbrains.com?from=Wails). Ihre Unterstützung treibt das Projekt auf vielfältige Weise hinter den Kulissen voran.

Ich freue mich darauf zu sehen, was die Menschen in dieser nächsten spannenden Phase des Projekts mit Wails entwickeln!

Lea.

PS: Linux-Nutzer, ihr seid als Nächstes dran!

PPS: Wenn Sie oder Ihr Unternehmen Wails nützlich finden, erwägen Sie bitte, [das Projekt zu unterstützen](https://github.com/sponsors/leaanthony). Vielen Dank!
