---
title: "Der Weg zu Wails v3"
description: "Versionshinweise und Ankündigungen für Wails"
authors: ["leaanthony"]
tags: ["wails","v3"]
date: "2023-01-17"
slug: "blog/the-road-to-wails-v3"
image: "/assets/blog-images/multiwindow.webp"
sourcePath: "blog/the-road-to-wails-v3.md"
---

![Screenshot mit mehreren Fenstern](/assets/blog-images/multiwindow.webp)

## Einführung

Wails vereinfacht die Entwicklung plattformübergreifender Desktopanwendungen mit Go. Für das Frontend verwendet es native Webview-Komponenten statt eingebetteter Browser. So bringt Wails die Leistungsfähigkeit des weltweit beliebtesten UI-Systems zu Go und bleibt dabei schlank.

Version 2 wurde am 22. September 2022 veröffentlicht und brachte zahlreiche Verbesserungen mit sich, darunter:

- Live-Entwicklung mithilfe des beliebten Vite-Projekts
- Umfangreiche Funktionen zur Fensterverwaltung und Menüerstellung
- Microsofts WebView2-Komponente
- Generierung von TypeScript-Modellen, die Ihre Go-Strukturen abbilden
- Erstellung eines NSIS-Installationsprogramms
- Verschleierte Builds

Derzeit bietet Wails v2 leistungsfähige Werkzeuge zum Erstellen funktionsreicher, plattformübergreifender Desktopanwendungen.

Dieser Blogbeitrag beleuchtet den aktuellen Stand des Projekts und zeigt, was wir künftig verbessern können.

## Wo stehen wir heute?

Es ist beeindruckend, wie stark die Beliebtheit von Wails seit der Veröffentlichung von v2 gestiegen ist. Die Kreativität der Community und die großartigen Dinge, die sie damit entwickelt, versetzen mich immer wieder in Erstaunen. Mit zunehmender Beliebtheit richtet sich auch mehr Aufmerksamkeit auf das Projekt. Damit nehmen auch Funktionswünsche und Fehlerberichte zu.

Im Laufe der Zeit konnte ich einige der dringendsten Probleme des Projekts sowie mehrere Faktoren identifizieren, die seine Entwicklung hemmen.

## Aktuelle Probleme

Ich habe folgende Bereiche identifiziert, die meiner Ansicht nach die Entwicklung des Projekts hemmen:

- Die API
- Generierung der Bindings
- Das Build-System

### Die API

Die API zum Erstellen einer Wails-Anwendung besteht derzeit aus 2 Teilen:

- Die Anwendungs-API
- Die Runtime-API

Die Anwendungs-API besteht bekanntermaßen aus nur 1 Funktion: `Run()`. Diese übernimmt zahlreiche Optionen, die das Verhalten der Anwendung bestimmen. Das ist zwar sehr einfach zu verwenden, aber auch stark einschränkend. Dieser „deklarative“ Ansatz verbirgt einen großen Teil der zugrunde liegenden Komplexität. Beispielsweise gibt es keinen Handle für das Hauptfenster, sodass Sie nicht direkt damit interagieren können. Dazu müssen Sie die Runtime-API verwenden. Das wird problematisch, sobald Sie komplexere Aufgaben ausführen möchten, etwa mehrere Fenster erstellen.

Die Runtime-API stellt Entwicklern zahlreiche Hilfsfunktionen bereit. Dazu gehören:

- Fensterverwaltung
- Dialogfelder
- Menüs
- Ereignisse
- Protokolle

Mit mehreren Aspekten der Runtime-API bin ich unzufrieden. Erstens muss ein „Kontext“ weitergereicht werden. Das ist frustrierend und verwirrend für neue Entwickler, die einen Kontext übergeben und anschließend einen Laufzeitfehler erhalten.

Das größte Problem der Runtime-API besteht darin, dass sie für Anwendungen mit nur einem Fenster konzipiert wurde. Im Laufe der Zeit ist die Nachfrage nach mehreren Fenstern gestiegen, doch die API eignet sich dafür nicht besonders gut.

### Überlegungen zur v3-API

Wäre es nicht großartig, wenn wir so etwas tun könnten?

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

Dieser programmatische Ansatz ist wesentlich intuitiver und ermöglicht Entwicklern, direkt mit den Anwendungselementen zu interagieren. Alle aktuellen Runtime-Methoden für Fenster würden einfach zu Methoden des Fensterobjekts. Die übrigen Runtime-Methoden könnten wir wie folgt in das Anwendungsobjekt verschieben:

```go
app := wails.NewApplication(options.App{})
app.NewInfoDialog(options.InfoDialog{})
app.Log.Info("Hello World")
```

Diese wesentlich leistungsfähigere API ermöglicht die Entwicklung komplexerer Anwendungen. Außerdem erlaubt sie das Erstellen mehrerer Fenster, [der Funktion mit den meisten positiven Stimmen auf GitHub](https://github.com/wailsapp/wails/issues/1480):

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    myWindow2 := app.NewWindow(options.Window{})
    myWindow2.SetTitle("My Window 2")
    myWindow2.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

### Generierung der Bindings

Eine der wichtigsten Funktionen von Wails ist die Generierung von Bindings für Ihre Go-Methoden, damit diese über JavaScript aufgerufen werden können. Das aktuelle Verfahren dafür ist eine etwas behelfsmäßige Lösung. Dabei wird die Anwendung mit einem speziellen Flag erstellt und anschließend die resultierende Binärdatei ausgeführt, die mittels Reflection ermittelt, welche Elemente gebunden wurden. Das führt zu einem Henne-Ei-Problem: Sie können die Anwendung nicht ohne die Bindings erstellen und die Bindings nicht generieren, ohne die Anwendung zu erstellen. Dafür gibt es viele Umgehungslösungen, doch am besten wäre es, diesen Ansatz gar nicht erst zu verwenden.

Es gab mehrere Versuche, einen statischen Analysator für Wails-Projekte zu entwickeln, die jedoch nicht weit kamen. In jüngerer Zeit ist dies etwas einfacher geworden, da mehr Material zu diesem Thema verfügbar ist.

Im Vergleich zu Reflection ist der AST-Ansatz wesentlich schneller, jedoch erheblich komplexer. Zunächst müssen wir möglicherweise bestimmte Einschränkungen dafür vorgeben, wie Bindings im Code angegeben werden. Ziel ist es, die häufigsten Anwendungsfälle zu unterstützen und die Unterstützung später zu erweitern.

### Das Build-System

Wie der deklarative Ansatz der API wurde auch das Build-System entwickelt, um die Komplexität beim Erstellen einer Desktopanwendung zu verbergen. Wenn Sie `wails build` ausführen, erledigt es im Hintergrund zahlreiche Aufgaben:

- Erstellt die Backend-Binärdatei für die Bindings und generiert die Bindings
- Installiert die Frontend-Abhängigkeiten
- Erstellt die Frontend-Assets
- Ermittelt, ob das Anwendungssymbol vorhanden ist, und bettet es gegebenenfalls ein
- Erstellt die endgültige Binärdatei
- Ist der Build für `darwin/universal` bestimmt, erstellt es 2 Binärdateien: eine für `darwin/amd64` und eine für `darwin/arm64`. Anschließend erzeugt es mit `lipo` eine universelle Binärdatei
- Ist eine Komprimierung erforderlich, komprimiert es die Binärdatei mit UPX
- Ermittelt, ob diese Binärdatei paketiert werden soll, und führt in diesem Fall folgende Schritte aus:
  - Stellt sicher, dass das Symbol und das Anwendungsmanifest in die Binärdatei kompiliert werden (Windows)
  - Erstellt das Anwendungspaket, generiert und kopiert das Symbolpaket und kopiert die Binärdatei sowie Info.plist in das Anwendungspaket (Mac)

- Erstellt bei Bedarf ein NSIS-Installationsprogramm

Dieser gesamte Prozess ist zwar sehr leistungsfähig, aber auch sehr undurchsichtig. Er lässt sich nur sehr schwer anpassen und debuggen.

Um dies in v3 zu beheben, möchte ich auf ein Build-System außerhalb von Wails umsteigen. Nachdem ich [Task](https://taskfile.dev/) eine Weile verwendet habe, bin ich ein großer Fan davon. Es ist ein hervorragendes Werkzeug zum Konfigurieren von Build-Systemen und dürfte allen, die bereits mit Makefiles gearbeitet haben, einigermaßen vertraut sein.

Das Build-System würde über eine `Taskfile.yml`-Datei konfiguriert, die standardmäßig mit jeder unterstützten Vorlage generiert würde. Sie würde alle Schritte enthalten, die für sämtliche derzeitigen Aufgaben wie das Erstellen oder Paketieren der Anwendung erforderlich sind, und so eine einfache Anpassung ermöglichen.

Für dieses Werkzeug wird keine externe Abhängigkeit erforderlich sein, da es Bestandteil der Wails CLI wäre. Das bedeutet, dass Sie `wails build` weiterhin verwenden können und der Befehl alles ausführt, was er heute bereits ausführt. Wenn Sie jedoch den Build-Prozess anpassen möchten, können Sie dazu die `Taskfile.yml`-Datei bearbeiten. Außerdem können Sie so die Build-Schritte leicht nachvollziehen und auf Wunsch Ihr eigenes Build-System verwenden.

Das fehlende Teil im Build-Puzzle sind die atomaren Operationen im Build-Prozess, etwa die Symbolerstellung, Komprimierung und Paketierung. Zahlreiche externe Werkzeuge vorauszusetzen, würde Entwicklern keine gute Erfahrung bieten. Daher wird die Wails CLI all diese Funktionen als Teil der CLI bereitstellen. So funktionieren die Builds weiterhin wie erwartet, ohne zusätzliche externe Werkzeuge; Sie können jedoch jeden Build-Schritt durch ein beliebiges Werkzeug Ihrer Wahl ersetzen.

Dadurch entsteht ein wesentlich transparenteres Build-System, das sich leichter anpassen lässt und viele der dazu gemeldeten Probleme behebt.

## Der Nutzen

Diese positiven Änderungen werden dem Projekt enorme Vorteile bringen:

- Die neue API wird wesentlich intuitiver sein und die Entwicklung komplexerer Anwendungen ermöglichen.
- Die statische Analyse zur Generierung der Bindings wird deutlich schneller sein und einen Großteil der Komplexität des derzeitigen Verfahrens reduzieren.
- Der Einsatz eines etablierten, externen Build-Systems wird den Build-Prozess vollständig transparent machen und umfassende Anpassungen ermöglichen.

Für die Projektbetreuer ergeben sich folgende Vorteile:

- Die neue API wird wesentlich einfacher zu warten und an neue Funktionen und Plattformen anzupassen sein.
- Das neue Build-System wird wesentlich einfacher zu warten und zu erweitern sein. Ich hoffe, dass dadurch ein neues Ökosystem aus von der Community entwickelten Build-Pipelines entsteht.
- Bessere Trennung der Zuständigkeiten innerhalb des Projekts. Dadurch lassen sich neue Funktionen und Plattformen leichter hinzufügen.

## Der Plan

Ein Großteil der Experimente hierfür ist bereits abgeschlossen, und die Ergebnisse sehen gut aus. Derzeit gibt es keinen Zeitplan für diese Arbeiten, aber ich hoffe, dass bis zum Ende des 1. Quartals 2023 eine Alpha-Version für Mac verfügbar sein wird, damit die Community sie testen, damit experimentieren und Feedback geben kann.

## Zusammenfassung

- Die v2-API ist deklarativ, verbirgt vieles vor Entwicklern und eignet sich nicht für Funktionen wie mehrere Fenster. Es wird eine neue API entwickelt, die einfacher, intuitiver und leistungsfähiger sein wird.
- Das Build-System ist undurchsichtig und schwer anzupassen. Daher werden wir auf ein externes Build-System umsteigen, das den gesamten Prozess offenlegt.
- Die Generierung der Bindings ist langsam und komplex. Daher werden wir auf statische Analyse umsteigen und so einen Großteil der Komplexität des derzeitigen Verfahrens beseitigen.

In das Innenleben von v2 ist viel Arbeit geflossen, und es ist solide. Nun ist es an der Zeit, die darüberliegende Schicht zu überarbeiten und die Entwicklungserfahrung erheblich zu verbessern.

Ich hoffe, Sie sind davon genauso begeistert wie ich. Ich freue mich darauf, Ihre Gedanken und Ihr Feedback zu hören.

Viele Grüße,

&dash; Lea

PS: Wenn Sie oder Ihr Unternehmen Wails nützlich finden, erwägen Sie bitte, [das Projekt zu unterstützen](https://github.com/sponsors/leaanthony). Vielen Dank!

PPS: Ja, das ist ein echtes Bildschirmfoto einer mit Wails erstellten Anwendung mit mehreren Fenstern. Es ist kein Mock-up. Es ist echt. Es ist großartig. Es kommt bald.
