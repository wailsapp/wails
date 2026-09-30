---
title: "Wails v2 veröffentlicht"
description: "Versionshinweise und Ankündigungen zu Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![Montage aus Screenshots](/assets/blog-images/montage.png)

## Es ist so weit!

Heute wird [Wails](https://wails.io) v2 veröffentlicht. Seit der ersten Alpha-Version von v2 sind etwa 18 Monate und seit der ersten Beta-Version etwa ein Jahr vergangen. Ich bin allen, die an der Weiterentwicklung des Projekts beteiligt waren, von Herzen dankbar.

Ein Grund dafür, dass es so lange gedauert hat, war unser Wunsch, vor der offiziellen Bezeichnung als v2 einen gewissen Grad an Vollständigkeit zu erreichen. Tatsächlich gibt es nie den perfekten Zeitpunkt, um ein Release mit einem Tag zu versehen – es gibt immer noch offene Probleme oder „nur noch ein“ weiteres Feature, das unbedingt hinein soll. Ein nicht perfektes Major-Release mit einem Tag zu versehen, bietet den Nutzern des Projekts jedoch etwas Stabilität und den Entwicklern zugleich die Gelegenheit für einen Neustart.

Dieses Release übertrifft alles, was ich je erwartet hatte. Ich hoffe, es bereitet Ihnen ebenso viel Freude, wie uns seine Entwicklung bereitet hat.

## Was *ist* Wails?

Falls Sie Wails noch nicht kennen: Mit diesem Projekt können Go-Programmierer mithilfe vertrauter Webtechnologien funktionsreiche Frontends für ihre Go-Programme erstellen. Es ist eine schlanke, auf Go basierende Alternative zu Electron. Viele weitere Informationen finden Sie auf der [offiziellen Website](https://wails.io/docs/introduction).

## Was ist neu?

Das v2-Release ist ein gewaltiger Fortschritt für das Projekt und behebt viele der Schwachstellen von v1. Falls Sie noch keinen der Blogbeiträge zu den Beta-Releases für [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) oder [Linux](/blog/wails-v2-beta-for-linux/) gelesen haben, empfehle ich Ihnen, dies nachzuholen, da dort alle wichtigen Änderungen ausführlicher behandelt werden. Zusammengefasst:

- WebView2-Komponente für Windows mit Unterstützung moderner Webstandards und Debugging-Funktionen.
- [Dunkles/helles Design](https://wails.io/docs/reference/options#theme) und [benutzerdefinierte Designs](https://wails.io/docs/reference/options#customtheme) unter Windows.
- Unter Windows ist CGO jetzt nicht mehr erforderlich.
- Integrierte Unterstützung für Projektvorlagen mit Svelte, Vue, React, Preact, Lit und Vanilla.
- [Vite](https://vitejs.dev/)-Integration für eine Entwicklungsumgebung mit Hot Reload für Ihre Anwendung.
- Native [Menüs](https://wails.io/docs/guides/application-development#application-menu) und [Dialoge](https://wails.io/docs/reference/runtime/dialog) für Anwendungen.
- Native Transparenzeffekte für Fenster unter [Windows](https://wails.io/docs/reference/options#windowistranslucent) und [macOS](https://wails.io/docs/reference/options#windowistranslucent-1). Unterstützung für Mica- und Acrylic-Hintergründe.
- Einfache Erstellung eines [NSIS-Installationsprogramms](https://wails.io/docs/guides/windows-installer) für die Bereitstellung unter Windows.
- Eine umfangreiche [Laufzeitbibliothek](https://wails.io/docs/reference/runtime/intro) mit Hilfsmethoden zur Fenstersteuerung sowie für Ereignisse, Dialoge, Menüs und Protokollierung.
- Unterstützung für die [Verschleierung](https://wails.io/docs/guides/obfuscated) Ihrer Anwendung mit [garble](https://github.com/burrowers/garble).
- Unterstützung für die Komprimierung Ihrer Anwendung mit [UPX](https://upx.github.io/).
- Automatische TypeScript-Generierung für Go-Strukturen. Weitere Informationen finden Sie [hier](https://wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Mit Ihrer Anwendung müssen auf keiner Plattform zusätzliche Bibliotheken oder DLLs ausgeliefert werden.
- Frontend-Assets müssen nicht gebündelt werden. Entwickeln Sie Ihre Anwendung einfach wie jede andere Webanwendung.

## Anerkennung und Dank

Der Weg zu v2 war mit enormem Aufwand verbunden. Zwischen der ersten Alpha-Version und der heutigen Veröffentlichung gab es ~2200 Commits von 89 Mitwirkenden. Viele, viele weitere haben Übersetzungen, Tests, Feedback und Hilfe in den Diskussionsforen sowie im Issue-Tracker beigesteuert. Ich bin jedem Einzelnen von euch unendlich dankbar. Ein ganz besonderer Dank gilt außerdem allen Projektsponsoren, die uns mit Orientierung, Rat und Feedback unterstützt haben. Alles, was ihr tut, wissen wir sehr zu schätzen.

Einige Personen möchte ich besonders erwähnen:

Zunächst ein **riesiges** Dankeschön an [@stffabi](https://github.com/stffabi), dessen zahlreiche Beiträge uns allen zugutekommen und der bei vielen Problemen umfangreiche Unterstützung geleistet hat. Er hat einige zentrale Features beigesteuert, darunter die Unterstützung externer Entwicklungsserver. Dadurch wurde unser Entwicklungsmodus grundlegend verbessert, denn wir konnten die Superkräfte von [Vite](https://vitejs.dev/) nutzen. Man kann mit Fug und Recht sagen, dass Wails v2 ohne seine [unglaublichen Beiträge](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04) ein weitaus weniger spannendes Release wäre. Vielen herzlichen Dank, @stffabi!

Ein riesiges Dankeschön geht auch an [@misitebao](https://github.com/misitebao), der unermüdlich die Website pflegt, chinesische Übersetzungen bereitstellt, Crowdin verwaltet und neuen Übersetzern den Einstieg erleichtert. Dies ist eine enorm wichtige Aufgabe, und ich bin für all die investierte Zeit und Mühe äußerst dankbar! Du bist großartig!

Zu guter Letzt ein riesiges Dankeschön an Mat Ryer, der während der Entwicklung von v2 mit Rat und Unterstützung zur Seite stand. Die gemeinsame Entwicklung von xBar mit einer frühen Alpha-Version von v2 half dabei, die Richtung von v2 zu bestimmen, und machte mir zugleich einige Designfehler der frühen Releases bewusst. Ich freue mich, ankündigen zu können, dass wir ab heute mit der Portierung von xBar auf Wails v2 beginnen und xBar zur Referenzanwendung des Projekts wird. Danke, Mat!

## Gewonnene Erkenntnisse

Auf dem Weg zu v2 haben wir einige Erkenntnisse gewonnen, die unsere künftige Entwicklung prägen werden.

## Kleinere, schnellere und fokussierte Releases

Im Laufe der Entwicklung von v2 wurden viele Features und Fehlerkorrekturen ad hoc entwickelt. Dies führte zu längeren Release-Zyklen und erschwerte die Fehlersuche. Künftig werden wir häufiger Releases mit weniger Features erstellen. Jedes Release umfasst Aktualisierungen der Dokumentation sowie gründliche Tests. Hoffentlich führen diese kleineren, schnelleren und fokussierten Releases zu weniger Regressionen und einer besseren Dokumentationsqualität.

## Beteiligung fördern

Als ich dieses Projekt begann, wollte ich jedem, der ein Problem hatte, sofort helfen. Issues waren für mich „persönlich“, und ich wollte sie so schnell wie möglich lösen. Das ist nicht nachhaltig und gefährdet letztlich die Langlebigkeit des Projekts. Künftig werde ich anderen mehr Raum geben, sich an der Beantwortung von Fragen und der Triage von Issues zu beteiligen. Hilfreich wären hierfür geeignete Werkzeuge. Wenn Sie Vorschläge haben, beteiligen Sie sich bitte [hier](https://github.com/wailsapp/wails/discussions/1855) an der Diskussion.

## Lernen, Nein zu sagen

Je mehr Menschen sich an einem Open-Source-Projekt beteiligen, desto mehr Anfragen nach zusätzlichen Features entstehen, die für die Mehrheit der Nutzer möglicherweise nützlich sind oder auch nicht. Die Entwicklung und Fehlersuche für diese Features erfordert zunächst Zeit und verursacht anschließend fortlaufenden Wartungsaufwand. Ich selbst mache mich dessen am meisten schuldig, da ich mir oft zu viel auf einmal vornehme, anstatt ein minimal funktionsfähiges Feature bereitzustellen. Künftig müssen wir beim Hinzufügen von Kernfunktionen etwas häufiger „Nein“ sagen und unsere Energie darauf konzentrieren, Entwicklern die Möglichkeit zu geben, diese Funktionalität selbst bereitzustellen. Für dieses Szenario prüfen wir Plugins ernsthaft. So kann jeder das Projekt nach eigenem Ermessen erweitern und zugleich auf einfache Weise dazu beitragen.

## Ausblick

Es gibt bereits sehr viele Kernfunktionen, die wir Wails im nächsten großen Entwicklungszyklus hinzufügen möchten. Die [Roadmap](https://github.com/wailsapp/wails/discussions/1484) steckt voller interessanter Ideen, und ich kann es kaum erwarten, mit der Arbeit daran zu beginnen. Einer der häufigsten Wünsche ist die Unterstützung mehrerer Fenster. Das ist nicht einfach richtig umzusetzen, und möglicherweise müssen wir dafür eine alternative API bereitstellen, da die aktuelle API nicht für diesen Anwendungsfall konzipiert wurde. Nach ersten Ideen und Rückmeldungen denke ich, dass Ihnen gefallen wird, in welche Richtung wir damit gehen möchten.

Ich persönlich freue mich sehr auf die Möglichkeit, Wails-Apps auf Mobilgeräten auszuführen. Wir haben bereits ein Demoprojekt, das zeigt, dass eine Wails-App unter Android ausgeführt werden kann. Daher möchte ich unbedingt erkunden, welche Möglichkeiten sich daraus ergeben!

Als letzten Punkt möchte ich die funktionale Gleichwertigkeit ansprechen. Seit Langem gilt der Grundsatz, dem Projekt nur Funktionen hinzuzufügen, die vollständig plattformübergreifend unterstützt werden. Obwohl sich dies bisher (größtenteils) als machbar erwiesen hat, hat es die Veröffentlichung neuer Funktionen erheblich verzögert. Künftig werden wir etwas anders vorgehen: Jede neue Funktion, die nicht sofort für alle Plattformen veröffentlicht werden kann, wird über eine experimentelle Konfiguration oder API bereitgestellt. So können Early Adopters auf bestimmten Plattformen die Funktion ausprobieren und Rückmeldungen geben, die in deren endgültige Ausgestaltung einfließen. Das bedeutet natürlich, dass es keine Garantie für die Stabilität der API gibt, solange die Funktion nicht auf allen Plattformen, auf denen sie unterstützt werden kann, vollständig unterstützt wird. Zumindest ermöglicht aber bereits die experimentelle Veröffentlichung, die Entwicklung fortzusetzen.

## Abschließende Worte

Ich bin wirklich stolz darauf, was wir mit der Veröffentlichung von V2 erreicht haben. Es ist beeindruckend zu sehen, was Menschen bereits mit den bisherigen Betaversionen entwickeln konnten – hochwertige Anwendungen wie [Varly](https://varly.app/), [Surge](https://getsurge.io/) und [October](https://october.utf9k.net/). Ich empfehle Ihnen, sie sich anzusehen.

Diese Version ist durch die harte Arbeit vieler Mitwirkender entstanden. Sie kann zwar kostenlos heruntergeladen und verwendet werden, doch ihre Entwicklung war keineswegs kostenlos. Machen Sie sich nichts vor: Dieses Projekt war mit erheblichen Opfern verbunden. Es hat nicht nur meine Zeit und die Zeit aller Mitwirkenden gekostet, sondern auch Zeit, die jeder von ihnen nicht mit Freunden und Familie verbringen konnte. Deshalb bin ich für jede Sekunde außerordentlich dankbar, die der Verwirklichung dieses Projekts gewidmet wurde. Je mehr Mitwirkende wir haben, desto besser lässt sich dieser Aufwand verteilen und desto mehr können wir gemeinsam erreichen. Ich möchte Sie alle ermutigen, sich eine Sache auszusuchen, zu der Sie beitragen können – sei es, den Fehlerbericht einer anderen Person zu bestätigen, eine Korrektur vorzuschlagen, die Dokumentation zu ändern oder jemandem zu helfen, der Unterstützung benötigt. All diese kleinen Dinge haben eine enorme Wirkung! Es wäre großartig, wenn auch Sie Teil der Geschichte auf dem Weg zu v3 wären.

Viel Freude damit!

&dash; Lea

PS: Wenn Sie oder Ihr Unternehmen Wails nützlich finden, erwägen Sie bitte, [das Projekt finanziell zu unterstützen](https://github.com/sponsors/leaanthony). Vielen Dank!
