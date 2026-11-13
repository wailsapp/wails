---
title: "Wails v3 Beta: ein neues Fundament für Go-Desktopanwendungen"
description: "Wails v3 Beta führt ein direkteres Anwendungsmodell, umfangreichere Bindings und ein klareres Fundament für Go-Desktopanwendungen ein."
authors: ["leaanthony"]
tags: ["wails","v3","beta"]
date: "2026-08-02"
slug: "blog/wails-v3-beta"
image: "/assets/screenshots/frameless-v3-native-corners-macos.png"
sourcePath: "blog/wails-v3-beta.md"
---

![Ein natives rahmenloses Fenster von Wails v3 unter macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Heute veröffentlichen wir Wails v3 Beta.

Mit Wails können Go-Entwickler Desktopanwendungen mit den ihnen bereits vertrauten Web-Frontend-Werkzeugen erstellen. Dabei kommt statt eines eingebetteten Browsers die native WebView der jeweiligen Plattform zum Einsatz. v3 ist ein bedeutender Fortschritt: Es bietet Anwendungen eine direktere API, ein klareres Build-Modell und ein besseres Fundament für die Desktopanwendungen, deren Unterstützung sich die Wails-Community gewünscht hat.

Dies ist eine Betaversion, nicht die endgültige Version von 3.0. Die Desktop-API ist stabil, und Teams setzen v3 bereits produktiv ein. Vor der Bereitstellung sollten Sie jedoch gründlich testen. Wir nutzen die Betaphase, um gemeinsam mit der Community die letzten Kompatibilitäts- und Workflow-Probleme zu finden. Wails v2 bleibt die aktuelle stabile Version und wird weiterhin Fehlerkorrekturen erhalten.

## Dokumentation während der Betaphase

Während der Betaphase pflegen wir die englische Dokumentation als maßgebliche Quelle, während API und Workflows abschließend validiert werden. In dieser Phase nehmen wir keine PRs für Übersetzungen an. Die Übersetzungsarbeit wird vor der allgemeinen Verfügbarkeit wieder aufgenommen, sobald die Dokumentation stabil genug ist, damit Übersetzer nicht wiederholt Änderungen nachziehen müssen.

## Was v3 enthält

- Eine explizite Anwendungs- und Fenster-API einschließlich erstklassiger Unterstützung mehrerer Fenster
- Go-Dienste mit statischer Quellcodeanalyse, die umfangreichere TypeScript-Bindings erzeugt und dabei Kommentare sowie aussagekräftige Parameternamen beibehält
- Dienste, die Frontend-Assets und Skripte zusammen mit ihrer Backend-API bündeln können – das Fundament für installierbare, umfangreichere Plug-ins
- Ein sichtbares, Taskfile-basiertes Build-System, das Sie untersuchen, erweitern und debuggen können
- Server-Builds, mit denen dieselbe Anwendung und dieselben Dienste ohne natives Desktopfenster ausgeführt werden können
- Moderne Desktopunterstützung für macOS, Windows und Linux auf Intel, Apple Silicon, amd64 und arm64, sofern unterstützt
- Experimentelle Mobilunterstützung für iOS und Android, die zum Ausprobieren verfügbar ist, aber nicht unter die Kompatibilitätszusage der Desktop-Betaversion fällt

## Warum v3

Wails v2 machte es einfach, eine Go-Anwendung mit einem modernen Web-Frontend zu erstellen. Es hat dem Projekt – und sehr vielen Anwendungen – gute Dienste geleistet. Doch seine auf ein einzelnes Fenster und Kontextübergabe ausgelegte Runtime sowie der streng verwaltete Build-Prozess erschwerten einige gängige Aufgaben bei der Desktopentwicklung unnötig.

v3 basiert auf einem anderen Modell. Anwendungen, Fenster, Dienste, Ereignisse und Plattformfunktionen sind explizite Objekte. Dadurch lässt sich das Framework bei wachsenden Anwendungen leichter nachvollziehen, und Funktionen wie mehrere Fenster werden zu einem regulären Bestandteil des Anwendungsmodells statt zu einer Behelfslösung.

## Neuerungen

### Eine Anwendungs-API für echte Desktopsoftware

v3 ersetzt den Konfigurationsstil v2 `wails.Run(...)` durch einen expliziten Anwendungslebenszyklus. Sie erstellen eine Anwendung, registrieren Dienste, erstellen Fenster und interagieren mit den Objekten, denen das benötigte Verhalten zugeordnet ist.

Dadurch entfällt ein großer Teil der impliziten Kontextübergabe. Fensteroperationen gehören zu den Fenstern, anwendungsweite Operationen zur Anwendung. Dieses Modell ist für Anwendungen mit mehreren Fenstern natürlicher und eignet sich besser zum Testen und Warten größerer Codebasen.

Das Feedback der Personen, die v3 während der Alphaphase verwendet haben, war ausgesprochen positiv. Besonders gut kam bei Entwicklern das explizite Modell an: Der Code lässt sich leichter nachvollziehen, Zuständigkeiten werden klarer, und komplexe Desktopanwendungen können wachsen, ohne gegen das Framework arbeiten zu müssen.

### Erstklassige Unterstützung mehrerer Fenster

Mehrere Fenster sind eine Kernfunktion von v3. Fenster besitzen einen eigenen Lebenszyklus und können zur Laufzeit erstellt, verwaltet und geschlossen werden. Dadurch entsteht ein klarerer Weg zu Desktopsoftware, die Editoren, Inspektoren, Einstellungen, Werkzeugfenster oder mehrere unabhängige Oberflächenbereiche benötigt.

### Dienste und generierte Bindings

Go-Dienste ersetzen das ältere Binding-Modell. Sie halten die Anwendungslogik in gewöhnlichem Go-Code und machen die Grenze zum Frontend explizit. Bindings werden in einer Struktur generiert, die die Anwendung und ihre Dienste abbildet. Dadurch lässt sich die für das Frontend bereitgestellte API leichter finden und verwenden.

v3 generiert diese Bindings mithilfe statischer Quellcodeanalyse. Dadurch kann der Generator die von Entwicklern in ihrem Code hinterlegten Informationen – einschließlich Kommentaren und aussagekräftiger Parameternamen – beibehalten, statt ein bereits gebautes Programm mittels Reflexion zu untersuchen. Das Ergebnis sind eine umfangreichere, nützlichere Frontend-API und ein Generierungsprozess, der leichter nachzuvollziehen und zu warten ist.

Dienste können außerdem Frontend-Assets und Skripte zusammen mit ihrem Go-Code einbinden. So erhält eine Funktion einen einzigen, zusammenhängenden Ort: ihre Backend-API, das benötigte JavaScript oder die benötigte Benutzeroberfläche und den Integrationspunkt für die Hostanwendung. Dies eröffnet die Möglichkeit für Wails-Plug-ins, die umfangreiche Funktionen sofort einsatzbereit bereitstellen: Installieren Sie ein Plug-in, binden Sie seinen Dienst ein und verwenden Sie die Funktion, statt selbst eine lose Sammlung von Bindings und Frontend-Abhängigkeiten zusammenzustellen. Ein allgemeines Plug-in-System ist nicht Bestandteil dieser Betaversion, doch mit v3 wird diese Richtung auf eine Weise praktikabel, die das Binding-Modell von v2 nicht ermöglichte.

### Ein Build-System, das Sie untersuchen und anpassen können

v3 macht die Build-Struktur des Projekts sichtbar. Statt jede Build-Entscheidung hinter einem einzigen Befehl zu verbergen, verfügen Projekte über eine konventionelle Struktur und eine Taskfile-basierte Build-Konfiguration, die sich zusammen mit der Anwendung nachvollziehen, erweitern und debuggen lässt.

### Eine stärkere plattformübergreifende Desktopbasis

Die Betaversion unterstützt Windows auf amd64 und arm64, macOS auf Intel und Apple Silicon sowie Linux auf amd64 und arm64. GTK4 mit WebKitGTK 6.0 ist der standardmäßige Linux-Stack; GTK3 bleibt während der gesamten v3.0-Versionsreihe als Legacy-Option verfügbar. Die Mobilunterstützung ist vielversprechend, bleibt jedoch experimentell und ist nicht Teil der Kompatibilitätszusage für die Desktop-Betaversion.

Diese Version umfasst außerdem die nötigen Arbeiten, um die tägliche Nutzung zuverlässiger zu machen: verbessertes Plattformverhalten, ein leistungsfähigeres Fenstermodell, klarere Diagnosen sowie Release-Artefakte mit Prüfsummen und Herkunftsnachweisen.

## Migration von v2

v3 ist eine neue Hauptversion, und die Migration ist eine echte Portierung, nicht lediglich eine Änderung der Versionsnummer. Die wichtigsten konzeptionellen Änderungen betreffen den Lebenszyklus von Anwendung und Fenstern, Dienste anstelle kontextgebundener Bindings, direkte Anwendungs- und Fenster-APIs anstelle des Runtime-Pakets von v2 sowie neu generierte Frontend-Bindings.

Wir haben einen [v2-zu-v3-Migrationsleitfaden](/migration/v2-to-v3/) veröffentlicht, der diese Änderungen erläutert und eine Funktionszuordnung sowie eine Testcheckliste enthält. Dieser manuelle Leitfaden ist der für diese Betaversion unterstützte Migrationsweg. Erwarten Sie bitte nicht, dass sich jedes v2-Projekt ohne Prüfung konvertieren lässt: Testen Sie das Ergebnis, portieren Sie Ihre Runtime-Aufrufe bewusst und behalten Sie v2 bei, bis die neue Anwendung bereit ist.

Wir evaluieren außerdem einen experimentellen Migrationsassistenten. Er ist nicht Bestandteil dieser Betaversion, und wir werden ihn erst empfehlen, nachdem er anhand repräsentativer realer v2-Projekte validiert wurde.

## Ein offenes Wort zum bisherigen Weg

Das erste Alpha-Tag von v3 wurde am 18. Januar 2023 veröffentlicht. Das ist eine lange Alphaphase und verdient mehr als eine vage Erwähnung.

Das Projekt hat sich in dieser Zeit enorm verändert. Als im September 2021 die erste Betaversion von Wails v2 für Windows erschien, hatte das Repository etwa 4000 Sterne. Bei der Veröffentlichung von v2 im September 2022 waren es ungefähr 10300. Heute sind es mehr als 35000. Dieses Wachstum ist ein Privileg, verändert aber auch, was verantwortungsvolle Projektführung bedeutet: Mehr Benutzer sind von Release-Entscheidungen abhängig, mehr Mitwirkende benötigen klare Wege zur Beteiligung, und ein größerer Teil der Arbeit besteht darin, das Projekt berechenbar zu machen, statt lediglich die nächste Funktion hinzuzufügen.

Ich habe unsere Prozesse nicht immer so schnell oder so klar an dieses Wachstum angepasst, wie es nötig gewesen wäre. Dafür übernehme ich die Verantwortung. Die Antwort darauf sind weder große Versprechen noch unnötige Förmlichkeit bei jeder Entscheidung. Stattdessen müssen wir klarer darlegen, was unterstützt wird, was experimentell ist, wie Entscheidungen getroffen werden und was wir als Nächstes tun.

Diese Arbeit beginnt mit dieser Beta. Wir haben nun ausdrücklich definierte Meilensteine für Beta, Release Candidate und allgemeine Verfügbarkeit (GA), klarere Kompatibilitätszusagen, eine aktualisierte Sicherheitsrichtlinie und einen WEP-Prozess (Wails Enhancement Proposal) für Änderungen am öffentlichen Verhalten und neue Funktionen. Während Wails wächst, werden wir die Roadmap weiter verbessern und die Governance des Projekts überprüfen. Ziel ist ein Projekt, auf das man sich leichter verlassen und zu dem man leichter beitragen kann – nicht eines, das sich schwerer voranbringen lässt.

Außerdem starten wir den [Wails-Subreddit](https://www.reddit.com/r/wails/) neu und werden ihn aktiv als weiteren Ort für praxisbezogene Diskussionen, Fragen und Feedback pflegen, während v3 auf die allgemeine Verfügbarkeit zusteuert.

## Beta installieren und ausprobieren

Sobald das Release veröffentlicht ist, installiere die neueste v3-CLI mit:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Führe vor dem Erstellen eines Projekts den geführten Einrichtungsassistenten aus. Er prüft deine lokale Entwicklungsumgebung und hilft dir, die von Wails benötigten Abhängigkeiten zu konfigurieren:

```sh
wails3 setup
```

![Der Einrichtungsassistent von Wails v3 im Dunkelmodus](/assets/screenshots/wails3-setup-wizard-dark.png)

### Projekt erstellen

Wenn alles eingerichtet ist, erstelle ein Projekt mit:

```sh
wails3 init
```

- Dokumentation: <https://v3.wails.io/>
- Migrationsleitfaden: <https://v3.wails.io/migration/v2-to-v3/>
- Reddit: <https://www.reddit.com/r/wails/>
- Versionshinweise: [VERSION PLACEHOLDER](https://github.com/wailsapp/wails/releases)

Wenn du einen reproduzierbaren Fehler findest, melde ihn bitte möglichst zusammen mit der Ausgabe von `wails3 doctor` und einem Minimalbeispiel. Wenn du eine neue Funktion oder eine Änderung am öffentlichen Verhalten vorschlagen möchtest, erstelle statt eines Feature-Request-Issues einen WEP-PR als Entwurf. Beide Wege helfen uns, klar zu reagieren und die Beta weiter voranzubringen.

## Vielen Dank

Wails v3 gibt es dank der Menschen, die unfertige Builds getestet, schwerwiegende Fehler gemeldet, Dokumentation übersetzt, Fragen beantwortet, Code beigetragen und das Projekt immer wieder dazu angespornt haben, besser zu werden. Vielen Dank.

Ein ganz besonderer und herzlicher Dank gilt außerdem den Sponsoren, die Wails durch diesen langen Übergang getragen haben. Eure Unterstützung hat nicht nur den Betrieb aufrechterhalten: Sie gab uns die Möglichkeit, kontinuierlich Zeit in die Architektur, die Werkzeuge, die Tests und die Dokumentation zu investieren, durch die das Projekt schneller auf v3 zusteuern konnte. Alle, die getestet oder Beiträge geleistet haben, haben dieses Release mitgestaltet. Die Sponsoren aber haben es ermöglicht, dieser Arbeit die Aufmerksamkeit zu widmen, die sie verdient hat.

Die Beta ist eine Einladung, uns dabei zu helfen, v3 sorgfältig fertigzustellen. Probiere sie aus, entwickle damit, sage uns, wo sie nicht funktioniert, und hilf uns, den Weg zu 3.0 kurz und sorgfältig zu gestalten.
