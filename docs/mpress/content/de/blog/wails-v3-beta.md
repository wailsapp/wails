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

![Ein natives rahmenloses Wails-v3-Fenster unter macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Heute veröffentlichen wir Wails v3 Beta.

Mit Wails können Go-Entwickler Desktopanwendungen mit den ihnen bereits vertrauten Web-Frontend-Werkzeugen erstellen. Dabei kommt statt eines eingebetteten Browsers die native WebView der jeweiligen Plattform zum Einsatz. v3 ist ein großer Fortschritt: Die Version bietet Anwendungen eine direktere API, ein klareres Build-Modell und ein besseres Fundament für die Desktopanwendungen, deren Unterstützung sich die Wails-Community gewünscht hat.

Dies ist eine Betaversion, nicht die endgültige 3.0-Version. Die Desktop-API ist stabil und Teams setzen v3 bereits produktiv ein. Sie sollten die Anwendung vor der Bereitstellung jedoch gründlich testen. In der Betaphase wollen wir gemeinsam mit der Community die letzten Kompatibilitäts- und Workflow-Probleme erkennen. Wails v2 bleibt die aktuelle stabile Version und wird weiterhin Fehlerkorrekturen erhalten.

## Dokumentation während der Betaphase

Während der Betaphase pflegen wir die englische Dokumentation als maßgebliche Referenz, solange API und Workflows abschließend validiert werden. In dieser Phase nehmen wir keine PRs mit Übersetzungen an. Die Übersetzungsarbeit wird vor der allgemeinen Verfügbarkeit wieder aufgenommen, sobald die Dokumentation stabil genug ist, damit Übersetzer nicht wiederholt Änderungen nachziehen müssen.

## Was v3 enthält

- Eine explizite Anwendungs- und Fenster-API einschließlich erstklassiger Unterstützung für mehrere Fenster
- Go-Dienste mit statischer Quellcodeanalyse, die umfangreichere TypeScript-Bindings generiert und dabei Kommentare sowie aussagekräftige Parameternamen beibehält
- Dienste, die Frontend-Assets und Skripte zusammen mit ihrer Backend-API bündeln können – das Fundament für installierbare, umfangreichere Plugins
- Ein transparentes, Taskfile-basiertes Build-System, das Sie untersuchen, erweitern und debuggen können
- Server-Builds, mit denen dieselbe Anwendung und dieselben Dienste ohne natives Desktopfenster ausgeführt werden können
- Moderne Desktopunterstützung für macOS, Windows und Linux auf Intel, Apple Silicon, amd64 und arm64, sofern unterstützt
- Experimentelle Mobilunterstützung für iOS und Android, die zum Ausprobieren verfügbar ist, aber nicht unter die Kompatibilitätszusage der Desktop-Beta fällt

## Warum v3

Wails v2 vereinfachte die Entwicklung einer Go-Anwendung mit einem modernen Web-Frontend. Die Version hat dem Projekt – und sehr vielen Anwendungen – gute Dienste geleistet. Doch die auf ein einzelnes Fenster und Kontextübergabe ausgelegte Runtime sowie der strikt verwaltete Build-Prozess erschwerten einige übliche Aufgaben bei Desktopanwendungen unnötig.

v3 geht von einem anderen Modell aus. Anwendungen, Fenster, Dienste, Ereignisse und Plattformfunktionen sind explizite Objekte. Dadurch lässt sich das Framework bei wachsenden Anwendungen leichter nachvollziehen. Außerdem werden Funktionen wie mehrere Fenster zu einem normalen Bestandteil des Anwendungsmodells statt zu einer Behelfslösung.

## Was ist neu?

### Eine Anwendungs-API für echte Desktopsoftware

v3 ersetzt den `wails.Run(...)`-Konfigurationsstil von v2 durch einen expliziten Anwendungslebenszyklus. Sie erstellen eine Anwendung, registrieren Dienste und erstellen Fenster. Anschließend interagieren Sie mit den Objekten, denen das jeweils benötigte Verhalten zugeordnet ist.

Dadurch entfällt ein großer Teil der impliziten Kontextübergabe. Fensteroperationen gehören zu Fenstern, anwendungsweite Operationen zur Anwendung. Dieses Modell ist für Anwendungen mit mehreren Fenstern natürlicher und eignet sich besser zum Testen und Warten größerer Codebasen.

Die Rückmeldungen der Personen, die v3 während der Alphaphase verwendet haben, waren überwältigend positiv. Insbesondere das explizite Modell fand bei Entwicklern großen Anklang: Der Code lässt sich leichter nachvollziehen, Zuständigkeiten werden klarer und komplexe Desktopanwendungen können wachsen, ohne gegen das Framework arbeiten zu müssen.

### Erstklassige Unterstützung für mehrere Fenster

Mehrere Fenster sind eine Kernfunktion von v3. Fenster haben einen eigenen Lebenszyklus und können zur Laufzeit erstellt, verwaltet und geschlossen werden. Dadurch ergibt sich ein klarerer Weg zu Desktopsoftware, die Editoren, Inspektoren, Einstellungsdialoge, Werkzeugfenster oder mehrere unabhängige UI-Bereiche benötigt.

### Dienste und generierte Bindings

Go-Dienste ersetzen das bisherige Binding-Modell. Die Anwendungslogik bleibt gewöhnlicher Go-Code, während die Schnittstelle zum Frontend explizit wird. Die Bindings werden in einer Struktur generiert, die die Anwendung und ihre Dienste abbildet. Dadurch lässt sich die für das Frontend bereitgestellte API leichter finden und verwenden.

v3 generiert diese Bindings mithilfe statischer Quellcodeanalyse. Dadurch kann der Generator die von Entwicklern im Code hinterlegten Informationen – einschließlich Kommentaren und aussagekräftigen Parameternamen – beibehalten, statt ein bereits erstelltes Programm durch Reflection zu untersuchen. Das Ergebnis ist eine umfangreichere, nützlichere Frontend-API und ein Generierungsprozess, der leichter zu verstehen und zu warten ist.

Dienste können neben ihrem Go-Code auch Frontend-Assets und Skripte einbinden. Damit erhält eine Funktion einen einzigen zusammenhängenden Ort: ihre Backend-API, das benötigte JavaScript oder UI und den Integrationspunkt für die Hostanwendung. Dies eröffnet die Möglichkeit für Wails-Plugins, die umfangreiche Funktionalität sofort bereitstellen – Plugin installieren, Dienst einbinden und Funktion nutzen –, statt selbst eine lose Sammlung aus Bindings und Frontend-Abhängigkeiten zusammenzustellen. Ein allgemeines Plugin-System ist nicht Teil dieser Beta, doch v3 macht diese Ausrichtung auf eine Weise praktikabel, die mit dem Binding-Modell von v2 nicht möglich war.

### Ein Build-System, das Sie untersuchen und anpassen können

v3 macht die Build-Struktur des Projekts sichtbar. Statt jede Build-Entscheidung in einem einzelnen Befehl zu verbergen, besitzen Projekte eine konventionelle Struktur und eine Taskfile-basierte Build-Konfiguration, die zusammen mit der Anwendung nachvollzogen, erweitert und debuggt werden kann.

### Eine stärkere plattformübergreifende Desktopbasis

Die Beta unterstützt Windows auf amd64 und arm64, macOS auf Intel und Apple Silicon sowie Linux auf amd64 und arm64. GTK4 mit WebKitGTK 6.0 ist der standardmäßige Linux-Stack; GTK3 bleibt während der gesamten v3.0-Reihe als Legacy-Option verfügbar. Die Mobilunterstützung ist vielversprechend, bleibt jedoch experimentell und ist nicht Teil der Kompatibilitätszusage der Desktop-Beta.

Diese Version enthält außerdem die erforderlichen Verbesserungen für einen zuverlässigeren Arbeitsalltag: verbessertes Plattformverhalten, ein leistungsfähigeres Fenstermodell, klarere Diagnosen und Release-Artefakte mit Prüfsummen und Herkunftsnachweisen.

## Migration von v2

v3 ist eine neue Hauptversion, und die Migration ist eine echte Portierung, nicht nur eine Änderung der Versionsnummer. Die wichtigsten konzeptionellen Änderungen sind die Lebenszyklen von Anwendung und Fenstern, Dienste statt kontextgebundener Bindings, direkte Anwendungs- und Fenster-APIs statt des Runtime-Pakets von v2 sowie neu generierte Frontend-Bindings.

Wir haben einen [Leitfaden für die Migration von v2 zu v3](/migration/v2-to-v3/) veröffentlicht, der diese Änderungen erläutert und eine Funktionszuordnung sowie eine Testcheckliste enthält. Dieser manuelle Leitfaden ist der für diese Beta unterstützte Migrationsweg. Erwarten Sie bitte nicht, dass sich jedes v2-Projekt ohne Überprüfung konvertieren lässt: Testen Sie das Ergebnis, portieren Sie Ihre Runtime-Aufrufe bewusst und behalten Sie v2 bei, bis die neue Anwendung bereit ist.

Wir evaluieren außerdem einen experimentellen Migrationsassistenten. Er ist nicht Teil dieser Betaversion, und wir werden ihn erst empfehlen, wenn er anhand repräsentativer realer v2-Projekte validiert wurde.

## Ein offenes Wort zum bisherigen Weg

Das erste Alpha-Tag von v3 wurde am 18. Januar 2023 veröffentlicht. Das ist eine lange Zeit für eine Alphaphase und verdient mehr als eine vage Erwähnung.

Das Projekt hat sich in dieser Zeit enorm verändert. Als im September 2021 die erste Betaversion von Wails v2 für Windows veröffentlicht wurde, hatte das Repository rund 4000 Sterne. Bis zur Veröffentlichung von v2 im September 2022 waren es rund 10300. Heute sind es mehr als 35000. Dieses Wachstum ist ein Privileg, verändert aber auch die Anforderungen an eine verantwortungsvolle Projektpflege: Mehr Benutzer sind von Release-Entscheidungen abhängig, mehr Mitwirkende benötigen klare Wege zur Beteiligung, und ein größerer Teil der Arbeit besteht darin, das Projekt berechenbar zu machen, statt einfach nur die nächste Funktion hinzuzufügen.

Ich habe unsere Prozesse nicht immer so schnell oder so klar angepasst, wie es dieses Wachstum erforderte. Dafür übernehme ich die Verantwortung. Die Antwort darauf besteht nicht in großen Versprechen oder darin, jede Entscheidung zu einem zeremoniellen Akt zu machen. Vielmehr müssen wir klarer kommunizieren, was unterstützt wird, was experimentell ist, wie Entscheidungen getroffen werden und was wir als Nächstes tun.

Diese Arbeit beginnt mit dieser Betaversion. Wir haben nun ausdrücklich definierte Meilensteine für Beta, Release Candidate und allgemeine Verfügbarkeit, klarere Kompatibilitätszusagen, eine aktualisierte Sicherheitsrichtlinie sowie einen WEP-Prozess (Wails Enhancement Proposal) für Änderungen am öffentlichen Verhalten und neue Funktionen. Während Wails wächst, werden wir die Roadmap weiter verbessern und die Projekt-Governance überprüfen. Das Ziel ist ein Projekt, auf das man sich leichter verlassen und zu dem man leichter beitragen kann – nicht eines, das sich schwerer weiterentwickeln lässt.

Außerdem starten wir den [Wails-Subreddit](https://www.reddit.com/r/wails/) neu und werden ihn aktiv als weiteren Ort für praktische Diskussionen, Fragen und Feedback pflegen, während sich v3 auf die allgemeine Verfügbarkeit zubewegt.

## Betaversion installieren und testen

Installieren Sie nach der Veröffentlichung des Releases die neueste v3-CLI mit:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Führen Sie vor dem Erstellen eines Projekts den geführten Einrichtungsassistenten aus. Er überprüft Ihre lokale Entwicklungsumgebung und hilft Ihnen, die von Wails benötigten Abhängigkeiten zu konfigurieren:

```sh
wails3 setup
```

![Der Einrichtungsassistent von Wails v3 im Dunkelmodus](/assets/screenshots/wails3-setup-wizard-dark.png)

### Projekt erstellen

Wenn alles eingerichtet ist, erstellen Sie ein Projekt mit:

```sh
wails3 init
```

- Dokumentation: <https://v3.wails.io/>
- Migrationsleitfaden: <https://v3.wails.io/migration/v2-to-v3/>
- Reddit: <https://www.reddit.com/r/wails/>
- Versionshinweise: [VERSIONSPLATZHALTER](https://github.com/wailsapp/wails/releases)

Wenn Sie einen reproduzierbaren Fehler finden, melden Sie ihn bitte zusammen mit der Ausgabe von `wails3 doctor` und, wenn möglich, einem minimalen Beispiel. Wenn Sie eine neue Funktion oder eine Änderung am öffentlichen Verhalten vorschlagen möchten, erstellen Sie statt eines Feature-Request-Issues einen WEP-PR als Entwurf. Beide Wege helfen uns, klar zu reagieren und die Betaphase voranzubringen.

## Vielen Dank

Wails v3 gibt es dank der Menschen, die unvollständige Builds getestet, schwierige Fehler gemeldet, Dokumentation übersetzt, Fragen beantwortet, Code beigetragen und das Projekt immer wieder dazu angespornt haben, besser zu werden. Vielen Dank.

Ein ganz besonderer, herzlicher Dank gilt außerdem den Sponsoren, die Wails durch diese lange Übergangsphase getragen haben. Ihre Unterstützung hat weit mehr bewirkt, als nur den Betrieb aufrechtzuerhalten: Sie gab uns die Möglichkeit, kontinuierlich Zeit in die Architektur, Werkzeuge, Tests und Dokumentation zu investieren, durch die das Projekt seinen Weg zu v3 beschleunigen konnte. Alle Tester und Mitwirkenden haben dieses Release mitgestaltet, doch die Sponsoren ermöglichten es, dieser Arbeit die Aufmerksamkeit zu widmen, die sie verdiente.

Die Betaversion ist eine Einladung, uns dabei zu helfen, v3 sorgfältig fertigzustellen. Probieren Sie sie aus, entwickeln Sie damit, teilen Sie uns mit, wo Probleme auftreten, und helfen Sie uns, den Weg zu 3.0 kurz und sorgfältig zu gestalten.
