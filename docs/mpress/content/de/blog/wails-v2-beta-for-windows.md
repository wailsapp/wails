---
title: "Wails v2 Beta für Windows"
description: "Versionshinweise und Ankündigungen zu Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-09-27"
slug: "blog/wails-v2-beta-for-windows"
image: "/assets/blog-images/wails.webp"
sourcePath: "blog/wails-v2-beta-for-windows.md"
---

![Wails-Screenshot](/assets/blog-images/wails.webp)

Als ich Wails vor etwas mehr als 2 Jahren aus einem Zug in Sydney erstmals auf Reddit ankündigte, erwartete ich nicht, dass es viel Aufmerksamkeit erhalten würde. Einige Tage später veröffentlichte ein produktiver Tech-Vlogger ein Tutorialvideo und bewertete es positiv. Seitdem ist das Interesse an dem Projekt sprunghaft gestiegen.

Es war klar, dass die Möglichkeit, Go-Projekte um Web-Frontends zu ergänzen, auf große Begeisterung stieß. Fast sofort trieb die Community das Projekt über den von mir erstellten Machbarkeitsnachweis hinaus voran. Damals verwendete Wails das Projekt [webview](https://github.com/webview/webview) für das Frontend, und unter Windows stand nur die IE11-Rendering-Engine zur Verfügung. Viele Fehlerberichte waren auf diese Einschränkung zurückzuführen: mangelhafte JavaScript- und CSS-Unterstützung sowie fehlende Entwicklungswerkzeuge zur Fehlersuche. Das machte die Entwicklung frustrierend, doch es gab kaum Möglichkeiten, das Problem zu beheben.

Ich war lange fest davon überzeugt, dass Microsoft seine Browserprobleme irgendwann lösen müsse. Die Welt entwickelte sich weiter, die Frontend-Entwicklung boomte und IE genügte den Anforderungen nicht mehr. Als Microsoft ankündigte, Chromium zur Grundlage seiner neuen Browserstrategie zu machen, wusste ich, dass Wails diese Technologie nur noch früher oder später nutzen und damit die Entwicklung unter Windows auf ein neues Niveau heben würde.

Heute freue ich mich, **Wails v2 Beta für Windows** anzukündigen! Diese Version bietet enorm viele Neuerungen. Holt euch also etwas zu trinken, setzt euch, und legen wir los …

## Keine CGO-Abhängigkeit!

Nein, das ist kein Scherz: *Keine* *CGO*-*Abhängigkeit* 🤯! Anders als MacOS und Linux enthält Windows standardmäßig keinen Compiler. Außerdem erfordert CGO einen mingw-Compiler, für den es unzählige Installationsmöglichkeiten gibt. Der Wegfall dieser CGO-Anforderung hat die Einrichtung erheblich vereinfacht und erleichtert zudem die Fehlersuche beträchtlich. Ich habe viel Arbeit investiert, um dies zum Laufen zu bringen. Der größte Dank gebührt jedoch [John Chadwick](https://github.com/jchv): Er hat nicht nur mehrere Projekte initiiert, die dies ermöglichen, sondern war auch offen dafür, dass jemand diese Projekte übernimmt und weiterentwickelt. Dank gebührt außerdem [Tad Vizbaras](https://github.com/tadvi), dessen Projekt [winc](https://github.com/tadvi/winc) mich auf diesen Weg brachte.

### WebView2-Chromium-Rendering-Engine

![Screenshot der Entwicklungswerkzeuge](/assets/blog-images/devtools.png)

Endlich erhalten Windows-Entwickler eine erstklassige Rendering-Engine für ihre Anwendungen! Die Zeiten, in denen ihr euren Frontend-Code mühsam verbiegen musstet, damit er unter Windows funktioniert, sind vorbei. Darüber hinaus stehen euch erstklassige Entwicklungswerkzeuge zur Verfügung!

Die WebView2-Komponente setzt allerdings voraus, dass sich `WebView2Loader.dll` neben der Binärdatei befindet. Das macht die Verteilung etwas umständlicher, als wir Gophers es gewohnt sind. Alle mir bekannten Lösungen und Bibliotheken, die WebView2 verwenden, haben diese Abhängigkeit.

Umso mehr freue ich mich, ankündigen zu können, dass für Wails-Anwendungen *keine solche Anforderung besteht*! Dank der Zauberkünste von [John Chadwick](https://github.com/jchv) können wir diese DLL in die Binärdatei einbetten und Windows dazu bringen, sie so zu laden, als läge sie auf dem Datenträger vor.

Freut euch, Gophers! Der Traum von einer einzigen Binärdatei lebt weiter!

### Neue Funktionen

![Screenshot von Wails-Menüs](/assets/blog-images/wails-menus.webp)

Viele von euch wünschten sich Unterstützung für native Menüs. Wails bietet sie nun endlich. Anwendungsmenüs sind jetzt verfügbar und unterstützen die meisten nativen Menüfunktionen. Dazu gehören Standardmenüeinträge, Kontrollkästchen, Optionsgruppen, Untermenüs und Trennlinien.

Für v1 gab es sehr viele Anfragen nach mehr Kontrolle über das Anwendungsfenster. Ich freue mich, dafür neue Runtime-APIs ankündigen zu können. Sie bieten viele Funktionen und unterstützen Konfigurationen mit mehreren Monitoren. Außerdem gibt es eine verbesserte Dialog-API: Ihr könnt nun moderne native Dialoge umfassend konfigurieren und damit all eure Anforderungen an Dialoge abdecken.

Ihr könnt jetzt zusammen mit eurem Projekt eine IDE-Konfiguration erzeugen. Wenn ihr das Projekt in einer unterstützten IDE öffnet, ist sie damit bereits zum Erstellen und Debuggen der Anwendung konfiguriert. Derzeit wird VSCode unterstützt; wir hoffen jedoch, bald weitere IDEs wie Goland zu unterstützen.

![VSCode-Screenshot](/assets/blog-images/vscode.webp)

### Kein Bündeln von Assets erforderlich

Ein großes Ärgernis in v1 war, dass ihr eure gesamte Anwendung in jeweils eine einzige JS- und CSS-Datei zusammenfassen musstet. Ich freue mich, ankündigen zu können, dass Assets in v2 überhaupt nicht mehr gebündelt werden müssen. Möchtet ihr ein lokales Bild laden? Verwendet ein `<img>`-Tag mit einem lokalen src-Pfad. Möchtet ihr eine ansprechende Schriftart verwenden? Kopiert sie in das Projekt und fügt den entsprechenden Pfad in eurem CSS hinzu.

> Wow, das klingt wie ein Webserver …

Ja, es funktioniert genau wie ein Webserver, ist aber keiner.

> Wie binde ich also meine Assets ein?

Übergebt der Anwendungskonfiguration einfach ein einziges `embed.FS`, das all eure Assets enthält. Sie müssen sich nicht einmal im obersten Verzeichnis befinden – Wails ermittelt alles automatisch für euch.

### Neues Entwicklungserlebnis

![Browser-Screenshot](/assets/blog-images/browser.webp)

Da Assets nicht mehr gebündelt werden müssen, ist nun ein völlig neues Entwicklungserlebnis möglich. Der neue Befehl `wails dev` erstellt und startet eure Anwendung, lädt die Assets jedoch direkt vom Datenträger, anstatt die Assets aus `embed.FS` zu verwenden.

Er bietet außerdem folgende zusätzliche Funktionen:

- Hot Reload – Jede Änderung an Frontend-Assets löst ein automatisches Neuladen des Anwendungs-Frontends aus
- Automatischer Neuaufbau – Jede Änderung an eurem Go-Code führt dazu, dass eure Anwendung neu erstellt und gestartet wird

Zusätzlich wird ein Webserver auf Port 34115 gestartet. Er stellt eure Anwendung jedem Browser bereit, der eine Verbindung zu ihm herstellt. Alle verbundenen Webbrowser reagieren auf Systemereignisse wie Hot Reload bei Änderungen an Assets.

In Go arbeiten wir in unseren Anwendungen gewöhnlich mit Structs. Häufig ist es sinnvoll, Structs an das Frontend zu senden und sie dort als Anwendungszustand zu verwenden. In v1 war dies ein weitgehend manueller Prozess und eine gewisse Belastung für Entwickler. Ich freue mich, ankündigen zu können, dass in v2 jede im Entwicklungsmodus ausgeführte Anwendung automatisch TypeScript-Modelle für alle Structs erzeugt, die als Ein- oder Ausgabeparameter gebundener Methoden dienen. Dies ermöglicht einen nahtlosen Austausch von Datenmodellen zwischen beiden Welten.

Zusätzlich wird dynamisch ein weiteres JS-Modul erzeugt, das all eure gebundenen Methoden umschließt. Es stellt JSDoc für eure Methoden bereit und ermöglicht damit Codevervollständigung und Hinweise in eurer IDE. Es ist wirklich beeindruckend, wenn ihr in einem automatisch erzeugten Modul, das euren Go-Code umschließt, die Tabulatortaste drückt und Datenmodelle automatisch importiert werden!

### Remote-Vorlagen

![Screenshot der Remote-Vorlage](/assets/blog-images/remote.webp)

Eine Anwendung schnell zum Laufen zu bringen, war für das Wails-Projekt schon immer ein wichtiges Ziel. Zum Start versuchten wir, viele der damals modernen Frameworks abzudecken: React, Vue und Angular. Die Frontend-Entwicklung ist von starken Meinungen geprägt, entwickelt sich rasant und ist nur schwer vollständig im Blick zu behalten. Daher veralteten unsere Basisvorlagen ziemlich schnell und verursachten erheblichen Wartungsaufwand. Außerdem fehlten uns attraktive moderne Vorlagen für die neuesten und besten Technologie-Stacks.

Mit v2 wollte ich die Community dazu befähigen, Vorlagen selbst zu erstellen und bereitzustellen, statt vom Wails-Projekt abhängig zu sein. Jetzt können Sie also Projekte mit Vorlagen erstellen, die von der Community unterstützt werden. Ich hoffe, dass dies Entwickler dazu anregt, ein lebendiges Ökosystem von Projektvorlagen aufzubauen. Ich bin wirklich gespannt darauf, was unsere Entwickler-Community erschaffen wird!

### Fazit

Wails v2 bildet eine neue Grundlage für das Projekt. Ziel dieser Version ist es, Rückmeldungen zum neuen Ansatz zu erhalten und vor einer vollständigen Veröffentlichung alle Fehler zu beheben. Ihre Rückmeldung ist sehr willkommen. Bitte richten Sie Ihr Feedback an das Diskussionsforum [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Bis zu diesem Punkt gab es viele Wendungen, Kurswechsel und Kehrtwenden. Das lag zum Teil an frühen technischen Entscheidungen, die geändert werden mussten, und zum Teil daran, dass einige grundlegende Probleme, für die wir mit viel Aufwand Umgehungslösungen entwickelt hatten, upstream behoben wurden. Die Einbettungsfunktion von Go ist dafür ein gutes Beispiel. Glücklicherweise fügte sich alles zum richtigen Zeitpunkt zusammen, und heute haben wir die bestmögliche Lösung. Ich glaube, das Warten hat sich gelohnt – noch vor 2 Monaten wäre dies nicht möglich gewesen.

Außerdem gebührt den folgenden Personen ein riesiges Dankeschön :pray:, denn ohne sie würde diese Version schlicht nicht existieren:

- [Misite Bao](https://github.com/misitebao) – Hat bei den chinesischen Übersetzungen unermüdlich gearbeitet und besitzt ein unglaubliches Gespür für Fehler.
- [John Chadwick](https://github.com/jchv) – Seine großartige Arbeit an [go-webview2](https://github.com/jchv/go-webview2) und [go-winloader](https://github.com/jchv/go-winloader) hat die heutige Windows-Version erst möglich gemacht.
- [Tad Vizbaras](https://github.com/tadvi) – Die Experimente mit seinem Projekt [winc](https://github.com/tadvi/winc) waren der erste Schritt auf dem Weg zu einem vollständig in Go entwickelten Wails.
- [Mat Ryer](https://github.com/matryer) – Seine Unterstützung, Ermutigung und Rückmeldungen haben entscheidend dazu beigetragen, das Projekt voranzubringen.

Abschließend möchte ich mich besonders bei allen [Projektsponsoren](/credits/#sponsors) bedanken, darunter [JetBrains](https://www.jetbrains.com?from=Wails). Ihre Unterstützung treibt das Projekt hinter den Kulissen auf vielfältige Weise voran.

Ich freue mich darauf zu sehen, was die Menschen in dieser nächsten spannenden Phase des Projekts mit Wails entwickeln werden!

Lea.

PS: Benutzer von macOS und Linux müssen sich nicht ausgeschlossen fühlen – die Portierung auf diese neue Grundlage ist bereits aktiv im Gange, und der größte Teil der schwierigen Arbeit ist schon erledigt. Bitte haben Sie noch etwas Geduld!

PPS: Wenn Sie oder Ihr Unternehmen Wails nützlich finden, erwägen Sie bitte, [das Projekt zu sponsern](https://github.com/sponsors/leaanthony). Vielen Dank!
