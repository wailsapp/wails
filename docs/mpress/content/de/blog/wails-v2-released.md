---
title: "Wails v2 veröffentlicht"
description: "Versionshinweise und Ankündigungen für Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![Montage-Screenshot](/assets/blog-images/montage.png)

## Es ist da!

Heute erscheint [Wails](https://wails.io) v2. Seit der ersten v2-Alpha sind etwa 18 Monate vergangen und seit der ersten Beta-Version ungefähr ein Jahr. Ich bin allen, die an der Weiterentwicklung des Projekts beteiligt waren, von Herzen dankbar.

Dass es so lange gedauert hat, lag unter anderem daran, dass wir einen gewissen Grad an Vollständigkeit erreichen wollten, bevor wir die Version offiziell v2 nannten. Tatsächlich gibt es nie den perfekten Zeitpunkt, um eine Version mit einem Tag zu versehen – es bleiben immer offene Probleme oder „nur noch ein“ Feature, das noch aufgenommen werden soll. Eine unvollkommene Hauptversion mit einem Tag zu versehen, bietet den Projektbenutzern jedoch etwas Stabilität und ermöglicht den Entwicklern einen gewissen Neustart.

Diese Version übertrifft alles, was ich je erwartet hatte. Ich hoffe, sie bereitet euch ebenso viel Freude, wie uns ihre Entwicklung bereitet hat.

## Was *ist* Wails?

Falls ihr Wails noch nicht kennt: Mit diesem Projekt können Go-Programmierer mithilfe vertrauter Webtechnologien umfangreiche Frontends für ihre Go-Programme erstellen. Es ist eine schlanke, auf Go basierende Alternative zu Electron. Viele weitere Informationen findet ihr auf der [offiziellen Website](https://v2.wails.io/docs/introduction).

## Was ist neu?

Die Version v2 ist ein gewaltiger Fortschritt für das Projekt und behebt viele Schwachstellen von v1. Falls ihr noch keinen der Blogbeiträge zu den Beta-Versionen für [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) oder [Linux](/blog/wails-v2-beta-for-linux/) gelesen habt, empfehle ich euch, dies nachzuholen, da dort alle wesentlichen Änderungen ausführlicher behandelt werden. Zusammengefasst:

- Webview2-Komponente für Windows mit Unterstützung moderner Webstandards und Debugging-Funktionen.
- [Dunkles/helles Design](https://v2.wails.io/docs/reference/options#theme) und [benutzerdefinierte Designs](https://v2.wails.io/docs/reference/options#customtheme) unter Windows.
- Windows benötigt jetzt kein CGO mehr.
- Integrierte Unterstützung für Projektvorlagen mit Svelte, Vue, React, Preact, Lit und Vanilla.
- [Vite](https://vitejs.dev/)-Integration für eine Entwicklungsumgebung mit Hot Reload für eure Anwendung.
- Native [Menüs](https://v2.wails.io/docs/guides/application-development#application-menu) und [Dialoge](https://v2.wails.io/docs/reference/runtime/dialog) für Anwendungen.
- Native Transparenzeffekte für Fenster unter [Windows](https://v2.wails.io/docs/reference/options#windowistranslucent) und [macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1). Unterstützung für Mica- und Acrylic-Hintergründe.
- Einfache Erstellung eines [NSIS-Installationsprogramms](https://v2.wails.io/docs/guides/windows-installer) für die Bereitstellung unter Windows.
- Eine umfangreiche [Laufzeitbibliothek](https://v2.wails.io/docs/reference/runtime/intro) mit Hilfsmethoden zur Fenstersteuerung sowie für Ereignisse, Dialoge, Menüs und Protokollierung.
- Unterstützung für die [Verschleierung](https://v2.wails.io/docs/guides/obfuscated) eurer Anwendung mit [garble](https://github.com/burrowers/garble).
- Unterstützung für die Komprimierung eurer Anwendung mit [UPX](https://upx.github.io/).
- Automatische Generierung von TypeScript aus Go-Strukturen. Weitere Informationen findet ihr [hier](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Mit eurer Anwendung müssen keine zusätzlichen Bibliotheken oder DLLs ausgeliefert werden. Auf keiner Plattform.
- Frontend-Ressourcen müssen nicht gebündelt werden. Entwickelt eure Anwendung einfach wie jede andere Webanwendung.

## Würdigung und Dank

v2 zu erreichen, war mit enormem Aufwand verbunden. Zwischen der ersten Alpha und der heutigen Veröffentlichung gab es etwa 2200 Commits von 89 Mitwirkenden. Viele, viele weitere Personen haben Übersetzungen, Tests, Feedback und Hilfe in den Diskussionsforen sowie im Issue-Tracker beigesteuert. Ich bin jedem Einzelnen von euch unendlich dankbar. Außerdem möchte ich allen Projektsponsoren ganz besonders für ihre Orientierungshilfe, Ratschläge und Rückmeldungen danken. Alles, was ihr leistet, wissen wir außerordentlich zu schätzen.

Einige Personen möchte ich besonders erwähnen:

Zunächst ein **riesiges** Dankeschön an [@stffabi](https://github.com/stffabi), der so viele Beiträge geleistet hat, von denen wir alle profitieren, und bei zahlreichen Problemen umfassende Unterstützung geboten hat. Er hat einige zentrale Features beigesteuert, darunter die Unterstützung externer Entwicklungsserver. Dadurch wurde unser Entwicklungsmodus grundlegend verbessert, weil wir nun die Superkräfte von [Vite](https://vitejs.dev/) nutzen können. Man kann mit Fug und Recht sagen, dass Wails v2 ohne seine [unglaublichen Beiträge](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04) eine weitaus weniger spannende Version wäre. Vielen Dank, @stffabi!

Ein riesiges Lob geht auch an [@misitebao](https://github.com/misitebao), der unermüdlich die Website gepflegt, chinesische Übersetzungen beigesteuert, Crowdin verwaltet und neuen Übersetzern beim Einstieg geholfen hat. Das ist eine enorm wichtige Aufgabe, und ich bin für all die investierte Zeit und Mühe außerordentlich dankbar! Du bist großartig!

Zu guter Letzt ein riesiges Dankeschön an Mat Ryer, der während der Entwicklung von v2 mit Rat und Unterstützung zur Seite stand. Die gemeinsame Entwicklung von xBar mit einer frühen Alpha-Version von v2 half dabei, die Richtung von v2 zu bestimmen und einige Designfehler der frühen Versionen zu erkennen. Ich freue mich, ankündigen zu können, dass wir ab heute mit der Portierung von xBar auf Wails v2 beginnen. xBar wird zur Vorzeigeanwendung des Projekts. Danke, Mat!

## Gewonnene Erkenntnisse

Auf dem Weg zu v2 haben wir einige Erkenntnisse gewonnen, die unsere künftige Entwicklung prägen werden.

## Kleinere, schnellere und zielgerichtete Versionen

Während der Entwicklung von v2 entstanden viele Features und Fehlerbehebungen auf Ad-hoc-Basis. Dies führte zu längeren Veröffentlichungszyklen und erschwerte das Debugging. Künftig werden wir häufiger Versionen mit weniger Features veröffentlichen. Jede Veröffentlichung umfasst sowohl Aktualisierungen der Dokumentation als auch gründliche Tests. Hoffentlich führen diese kleineren, schnelleren und gezielteren Versionen zu weniger Regressionen und einer hochwertigeren Dokumentation.

## Beteiligung fördern

Als ich dieses Projekt begann, wollte ich sofort allen helfen, die ein Problem hatten. Issues waren für mich „persönlich“, und ich wollte sie so schnell wie möglich lösen. Das ist nicht nachhaltig und wirkt sich letztlich negativ auf die Langlebigkeit des Projekts aus. Künftig werde ich anderen mehr Raum geben, sich an der Beantwortung von Fragen und der Triage von Issues zu beteiligen. Hilfreiche Werkzeuge dafür wären sehr willkommen. Falls ihr Vorschläge habt, beteiligt euch bitte [hier](https://github.com/wailsapp/wails/discussions/1855) an der Diskussion.

## Lernen, Nein zu sagen

Je mehr Menschen sich an einem Open-Source-Projekt beteiligen, desto mehr Anfragen nach zusätzlichen Features entstehen, die für die Mehrheit nützlich sein können oder auch nicht. Die Entwicklung und Fehlersuche für diese Features beansprucht zunächst Zeit und verursacht anschließend fortlaufenden Wartungsaufwand. Ich selbst mache mich dessen am meisten schuldig, weil ich oft „das Meer zum Kochen bringen“ will, statt nur das kleinstmögliche brauchbare Feature bereitzustellen. Künftig müssen wir bei der Aufnahme von Kernfeatures etwas häufiger „Nein“ sagen und unsere Energie darauf konzentrieren, Entwicklern die Möglichkeit zu geben, diese Funktionalität selbst bereitzustellen. Für dieses Szenario prüfen wir Plugins ernsthaft. Damit kann jeder das Projekt nach eigenen Vorstellungen erweitern und zugleich auf einfache Weise dazu beitragen.

## Blick in die Zukunft

Wir erwägen bereits zahlreiche Kernfeatures, die wir Wails im nächsten großen Entwicklungszyklus hinzufügen möchten. Die [Roadmap](https://github.com/wailsapp/wails/discussions/1484) ist voller interessanter Ideen, und ich freue mich darauf, mit der Arbeit daran zu beginnen. Besonders häufig wurde die Unterstützung mehrerer Fenster gewünscht. Das ist knifflig, und um es richtig umzusetzen, müssen wir möglicherweise eine alternative API bereitstellen, da die aktuelle API nicht dafür konzipiert wurde. Ausgehend von einigen ersten Ideen und Rückmeldungen glaube ich, dass euch die geplante Richtung gefallen wird.

Mich persönlich begeistert die Aussicht sehr, Wails-Anwendungen auf Mobilgeräten auszuführen. Wir haben bereits ein Demoprojekt, das zeigt, dass eine Wails-Anwendung unter Android ausgeführt werden kann. Daher möchte ich unbedingt erkunden, welche Möglichkeiten sich daraus ergeben!

Ein letzter Punkt, den ich ansprechen möchte, ist die Funktionsparität. Es war lange ein Grundprinzip, dem Projekt nichts hinzuzufügen, das nicht vollständig plattformübergreifend unterstützt wird. Obwohl sich dies bisher (größtenteils) als machbar erwiesen hat, hat es die Veröffentlichung neuer Funktionen erheblich ausgebremst. Künftig verfolgen wir einen etwas anderen Ansatz: Jede neue Funktion, die nicht sofort für alle Plattformen veröffentlicht werden kann, wird über eine experimentelle Konfiguration oder API bereitgestellt. So können Early Adopters auf bestimmten Plattformen die Funktion ausprobieren und Rückmeldungen geben, die in ihre endgültige Ausgestaltung einfließen. Das bedeutet natürlich, dass die Stabilität der API erst gewährleistet ist, wenn die Funktion von allen Plattformen, auf denen sie unterstützt werden kann, vollständig unterstützt wird. Zumindest wird dadurch aber die Entwicklung nicht länger blockiert.

## Abschließende Worte

Ich bin wirklich stolz darauf, was wir mit dem Release V2 erreicht haben. Es ist beeindruckend zu sehen, was andere bereits mit den bisherigen Betaversionen entwickeln konnten: hochwertige Anwendungen wie [Varly](https://varly.app/), [Surge](https://getsurge.io/) und [October](https://october.utf9k.net/). Ich kann nur empfehlen, sie sich anzusehen.

Dieser Release wurde durch die harte Arbeit vieler Mitwirkender ermöglicht. Wails kann zwar kostenlos heruntergeladen und verwendet werden, ist aber keineswegs ohne Aufwand entstanden. Dieses Projekt hat zweifellos einen erheblichen Preis gefordert. Dazu gehören nicht nur meine Zeit und die Zeit aller Mitwirkenden, sondern auch die Zeit, in der jeder Einzelne seinen Freunden und seiner Familie gefehlt hat. Deshalb bin ich für jede Sekunde außerordentlich dankbar, die diesem Projekt gewidmet wurde. Je mehr Mitwirkende wir haben, desto besser lässt sich dieser Aufwand verteilen und desto mehr können wir gemeinsam erreichen. Ich möchte euch alle ermutigen, eine Sache auszuwählen, zu der ihr beitragen könnt – sei es, den Fehlerbericht einer anderen Person zu bestätigen, eine Fehlerbehebung vorzuschlagen, die Dokumentation zu ändern oder jemandem zu helfen, der Unterstützung benötigt. All diese kleinen Dinge haben eine enorme Wirkung! Es wäre großartig, wenn auch ihr Teil der Geschichte auf dem Weg zu v3 wärt.

Viel Spaß!

&dash; Lea

PS: Wenn Wails für euch oder euer Unternehmen nützlich ist, zieht bitte in Betracht, [das Projekt zu sponsern](https://github.com/sponsors/leaanthony). Vielen Dank!
