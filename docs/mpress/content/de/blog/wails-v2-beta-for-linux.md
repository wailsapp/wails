---
title: "Wails v2 Beta für Linux"
description: "Versionshinweise und Ankündigungen zu Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-02-22"
slug: "blog/wails-v2-beta-for-linux"
image: "/assets/blog-images/wails-linux.webp"
sourcePath: "blog/wails-v2-beta-for-linux.md"
---

![Screenshot von wails-linux](/assets/blog-images/wails-linux.webp)

Ich freue mich, endlich bekannt geben zu können, dass Wails v2 nun als Beta für Linux verfügbar ist! Es ist etwas ironisch, dass die allerersten Experimente mit v2 unter Linux stattfanden und Linux dennoch als letzte Plattform veröffentlicht wird. Allerdings unterscheidet sich die heutige v2 stark von diesen ersten Experimenten. Kommen wir also ohne weitere Umschweife zu den neuen Funktionen:

## Neue Funktionen

![Screenshot von wails-menus-linux](/assets/blog-images/wails-menus-linux.webp)

Native Menüs wurden vielfach gewünscht. Wails unterstützt sie nun endlich. Anwendungsmenüs sind jetzt verfügbar und unterstützen die meisten nativen Menüfunktionen. Dazu gehören Standardmenüeinträge, Kontrollkästchen, Optionsgruppen, Untermenüs und Trennlinien.

Für v1 wurde sehr häufig mehr Kontrolle über das Fenster selbst gewünscht. Ich freue mich, bekannt geben zu können, dass es dafür neue Runtime-APIs gibt. Sie bieten zahlreiche Funktionen und unterstützen Konfigurationen mit mehreren Monitoren. Außerdem gibt es eine verbesserte Dialog-API: Sie können jetzt moderne, native Dialoge umfassend konfigurieren und so all Ihre Anforderungen an Dialoge abdecken.

### Assets müssen nicht gebündelt werden

Ein großes Problem von v1 war, dass die gesamte Anwendung in jeweils eine einzige JS- und CSS-Datei zusammengefasst werden musste. Ich freue mich, bekannt geben zu können, dass Assets für v2 in keiner Weise gebündelt werden müssen. Möchten Sie ein lokales Bild laden? Verwenden Sie ein `<../../../assets/blog-images>`-Tag mit einem lokalen src-Pfad. Möchten Sie eine coole Schriftart verwenden? Kopieren Sie sie hinein und fügen Sie den entsprechenden Pfad zu Ihrer CSS-Datei hinzu.

> Wow, das klingt nach einem Webserver ...

Ja, es funktioniert genau wie ein Webserver, ist aber keiner.

> Wie binde ich also meine Assets ein?

Übergeben Sie Ihrer Anwendungskonfiguration einfach ein einziges `embed.FS`, das all Ihre Assets enthält. Sie müssen sich nicht einmal im obersten Verzeichnis befinden – Wails kümmert sich darum.

### Neue Entwicklungserfahrung

Da Assets nun nicht mehr gebündelt werden müssen, ergibt sich eine völlig neue Entwicklungserfahrung. Der neue Befehl `wails dev` erstellt und startet Ihre Anwendung, lädt die Assets jedoch nicht aus `embed.FS`, sondern direkt vom Datenträger.

Darüber hinaus bietet er folgende Funktionen:

- Hot Reload – Jede Änderung an Frontend-Assets löst ein automatisches Neuladen des Anwendungs-Frontends aus
- Automatischer Neuaufbau – Jede Änderung an Ihrem Go-Code erstellt und startet Ihre Anwendung neu

Zusätzlich wird auf Port 34115 ein Webserver gestartet. Er stellt Ihre Anwendung jedem Browser bereit, der eine Verbindung zu ihm herstellt. Alle verbundenen Webbrowser reagieren auf Systemereignisse, etwa mit einem Hot Reload bei Änderungen an Assets.

In Go arbeiten wir in unseren Anwendungen regelmäßig mit Structs. Häufig ist es nützlich, Structs an unser Frontend zu senden und sie dort als Anwendungszustand zu verwenden. In v1 war dies ein weitgehend manueller Prozess und eine gewisse Belastung für Entwickler. Ich freue mich, bekannt geben zu können, dass in v2 jede im Entwicklungsmodus ausgeführte Anwendung automatisch TypeScript-Modelle für alle Structs generiert, die als Ein- oder Ausgabeparameter gebundener Methoden dienen. Dadurch können Datenmodelle nahtlos zwischen beiden Welten ausgetauscht werden.

Zusätzlich wird dynamisch ein weiteres JS-Modul generiert, das all Ihre gebundenen Methoden umschließt. Es stellt JSDoc für Ihre Methoden bereit und ermöglicht dadurch Codevervollständigung und Hinweise in Ihrer IDE. Es ist wirklich cool, wenn Sie in einem automatisch generierten Modul, das Ihren Go-Code umschließt, die Tabulatortaste drücken und Datenmodelle automatisch importiert werden!

### Remote-Vorlagen

![Screenshot von remote-linux](/assets/blog-images/remote-linux.webp)

Eine Anwendung schnell zum Laufen zu bringen, war für das Wails-Projekt schon immer ein wichtiges Ziel. Zum Projektstart versuchten wir, viele der damals modernen Frameworks abzudecken: react, vue und angular. Die Welt der Frontend-Entwicklung ist stark von unterschiedlichen Ansichten geprägt, schnelllebig und schwer im Blick zu behalten! Daher veralteten unsere Basisvorlagen ziemlich schnell, was einen hohen Wartungsaufwand verursachte. Außerdem fehlten uns coole, moderne Vorlagen für die neuesten und besten Technologie-Stacks.

Mit v2 wollte ich die Community dazu befähigen, Vorlagen selbst zu erstellen und zu hosten, statt sich auf das Wails-Projekt verlassen zu müssen. Sie können Projekte jetzt also mit von der Community unterstützten Vorlagen erstellen! Ich hoffe, dass dies Entwickler dazu inspiriert, ein lebendiges Ökosystem aus Projektvorlagen aufzubauen. Ich bin wirklich sehr gespannt darauf, was unsere Entwickler-Community erschaffen wird!

### Cross-Kompilierung für Windows

Da Wails v2 für Windows vollständig in Go geschrieben ist, können Sie Builds für Windows ohne docker erstellen.

![Screenshot von build-cross-windows](/assets/blog-images/linux-build-cross-windows.webp)

### Fazit

Wie ich bereits in den Versionshinweisen für Windows erwähnt habe, bildet Wails v2 eine neue Grundlage für das Projekt. Ziel dieser Veröffentlichung ist es, Feedback zum neuen Ansatz zu erhalten und vor der vollständigen Veröffentlichung alle Fehler auszuräumen. Ihre Rückmeldungen sind sehr willkommen! Bitte richten Sie jegliches Feedback an das Diskussionsforum [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Linux zu unterstützen ist **schwierig**. Wir rechnen bei der Beta mit einigen Eigenheiten. Helfen Sie uns bitte, Ihnen zu helfen, indem Sie detaillierte Fehlerberichte einreichen!

Abschließend möchte ich allen [Projektsponsoren](/credits/#sponsors) besonders danken, deren Unterstützung das Projekt auf vielfältige Weise hinter den Kulissen voranbringt.

Ich freue mich darauf zu sehen, was die Menschen in dieser nächsten spannenden Phase des Projekts mit Wails entwickeln werden!

Lea.

PS: Die Veröffentlichung von v2 ist nicht mehr weit entfernt!

PPS: Wenn Sie oder Ihr Unternehmen Wails nützlich finden, erwägen Sie bitte, [das Projekt zu sponsern](https://github.com/sponsors/leaanthony). Vielen Dank!
