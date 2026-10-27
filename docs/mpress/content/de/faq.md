---
title: "Häufig gestellte Fragen"
description: "Antworten auf die häufigsten Fragen zur Entwicklung von Anwendungen mit Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Allgemeines

### Was ist Wails?

Wails ist ein Framework zur Entwicklung von Desktopanwendungen mit Go und Webtechnologien. Sie schreiben die Anwendungslogik in Go und erstellen die Benutzeroberfläche mit HTML, CSS und JavaScript oder einem beliebigen Frontend-Framework. Wails stellt sie in der nativen Webview des Betriebssystems dar. Das Ergebnis ist eine kleine, schnelle und nativ wirkende Anwendung: kein mitgelieferter Browser, geringer Speicherbedarf und eine einzelne Binärdatei, die typischerweise etwa 10 MB groß ist.

### Welche Plattformen unterstützt Wails?

| Plattform | Voraussetzungen |
| --- | --- |
| Windows | AMD64 und ARM64. Verwendet die [WebView2-Laufzeit](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15+ auf Intel (Anwendungen können auf 10.13+ abzielen), 11.0+ auf Apple Silicon. Universal Binaries werden unterstützt. |
| Linux | AMD64 und ARM64. Der Standard-Stack ist GTK4 mit WebKitGTK 6.0 (Ubuntu 24.04+, Debian 13+, Fedora 40+ und vergleichbare Distributionen). Distributionen, die nur WebKit2GTK 4.1 bereitstellen, etwa Ubuntu 22.04, Debian 12 und RHEL 9, werden über den Legacy-Build `-tags gtk3` unterstützt, der bis v3.1 verfügbar ist. Distributionen, die nur WebKit2GTK 4.0 bereitstellen, werden nicht unterstützt. Weitere Informationen finden Sie im [Linux-Build-Leitfaden](/guides/build/linux/). |
| iOS und Android | Experimentell. Weitere Informationen finden Sie in den [Leitfäden für Mobilplattformen](/guides/mobile/). |

Mit dem [Server-Build](/guides/server-build/) können Sie Ihre Anwendung auch als reguläre Webanwendung bereitstellen.

Führen Sie jederzeit `wails3 doctor` aus, um Ihr System zu prüfen und plattformspezifische Installationsanweisungen zu erhalten.

### Was benötige ich für den Einstieg?

- Go 1.26 oder neuer
- Node.js und npm (für den Frontend-Build)
- Plattform-Toolchain: WebView2 unter Windows (auf 10/11 vorinstalliert), Xcode Command Line Tools unter macOS sowie `gcc` und GTK-/WebKit-Entwicklungspakete unter Linux

`wails3 doctor` prüft all diese Voraussetzungen und nennt Ihnen genau die fehlenden Komponenten. Eine vollständige Anleitung finden Sie unter [Installation](/quick-start/installation/).

### Ist Wails v3 produktionsreif?

Wails v3.0.0 ist die stabile Version für Desktopanwendungen unter Windows, macOS und Linux. Die Unterstützung für Android und iOS bleibt experimentell und fällt nicht unter die stabile Kompatibilitätszusage für Desktopplattformen. Weitere Informationen finden Sie auf der [Projektstatusseite](/status/).

## Entwicklung

### Muss ich Go beherrschen?

Grundkenntnisse in Go sind hilfreich, Sie müssen jedoch kein Experte sein. Die Anwendungslogik befindet sich in regulären Go-Methoden, und die [Tutorials](/tutorials/overview/) führen Sie durch alle weiteren Schritte. Viele Entwickler lernen Go während der Entwicklung ihrer ersten Wails-Anwendung.

### Kann ich mein bevorzugtes Frontend-Framework verwenden?

Ja. Wenn es sich zu HTML, CSS und JavaScript bauen lässt, funktioniert es mit Wails. Vorlagen sind für React, Vue, Svelte und Vanilla JavaScript verfügbar, jeweils auch in TypeScript-Varianten. Alle anderen Frameworks lassen sich innerhalb weniger Minuten einbinden. Weitere Informationen finden Sie unter [Frontend-Frameworks](/guides/dev/frontend-frameworks/).

### Wie rufe ich Go-Funktionen aus JavaScript auf?

Registrieren Sie einen Dienst; Wails generiert dafür typisierte Bindings:

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

Bindings werden während `wails3 dev` automatisch oder bei Bedarf mit `wails3 generate bindings` neu generiert. Weitere Informationen finden Sie unter [Dienste](/features/bindings/services/).

### Kann ich TypeScript verwenden?

Ja. Der Binding-Generator erzeugt TypeScript-Definitionen für Ihre Dienste und deren Typen, sodass Aufrufe von Go vollständig typisiert sind.

### Wie sende ich Ereignisse zwischen Go und JavaScript?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

Ereignisnamen müssen exakt übereinstimmen. Weitere Informationen finden Sie in der [Ereignisreferenz](/guides/events-reference/).

### Wie debugge ich meine Anwendung?

Führen Sie `wails3 dev` aus und öffnen Sie per Rechtsklick im Fenster die Browser-Entwicklertools, genau wie bei einer Webanwendung. Der Entwicklungsserver unterstützt außerdem Hot Reload für Ihr Frontend. Weitere Informationen finden Sie unter [Debugging](/guides/dev/debugging/).

## Build und Distribution

### Wie erstelle ich einen Produktions-Build?

```bash
wails3 build
```

Ihre Binärdatei wird unter `bin/` abgelegt. Produktions-Builds verwenden bereits sinnvolle Standardeinstellungen (Build-Tags, `-trimpath`, entfernte Symbole), sodass für eine schlanke Binärdatei keine zusätzlichen Flags erforderlich sind.

### Kann ich Cross-Compiling verwenden?

Mit Einschränkungen. Reines Go-Cross-Compiling ist nicht möglich, da jede Plattform native Webview-Bibliotheken verwendet. Häufige Anwendungsfälle werden jedoch gut unterstützt:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

Für Builds für Linux unter einem anderen Betriebssystem wird eine Docker-basierte Toolchain verwendet. Die vollständige Matrix finden Sie unter [Plattformübergreifende Builds](/guides/build/cross-platform/).

### Wie erstelle ich ein Installationsprogramm oder Paket?

```bash
wails3 package
```

Dadurch wird das native Format der jeweiligen Plattform erzeugt. Der [Leitfaden für Installationsprogramme](/guides/installers/) behandelt NSIS unter Windows, `.app`-Bundles und DMGs unter macOS sowie Linux-Pakete.

### Wie signiere ich den Code meiner Anwendung?

Das Signieren unter Windows und macOS einschließlich der Beglaubigung wird im [Leitfaden zur Codesignierung](/guides/build/signing/) Schritt für Schritt beschrieben.

## Funktionen

### Kann ich mehrere Fenster erstellen?

Ja, v3 bietet native Unterstützung für mehrere Fenster:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Weitere Informationen finden Sie unter [Mehrere Fenster](/features/windows/multiple/).

### Unterstützt Wails den Infobereich des Systems?

Ja, einschließlich Menüs und Klick-Handlern:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Weitere Informationen finden Sie unter [Infobereich](/features/menus/systray/).

### Kann ich native Dialogfelder verwenden?

Ja. Datei-, Meldungs- und Fragedialoge verwenden jeweils die nativen Implementierungen:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Weitere Informationen finden Sie unter [Dialogfelder](/features/dialogs/overview/).

### Unterstützt Wails automatische Updates?

Ja. Wails v3 enthält einen integrierten Selbstaktualisierer (`app.Updater`) mit austauschbaren Providern für GitHub Releases, keygen.sh und Sparkle AppCast, kryptografischer Signaturprüfung sowie einer Standardbenutzeroberfläche, die Sie gestalten oder ersetzen können. Weitere Informationen finden Sie im Leitfaden zum [In-App-Updater](/guides/updater/) und im Tutorial zur [selbstaktualisierenden Wails-Anwendung](/tutorials/04-self-update-a-wails-app/).

## Fehlerbehebung

### Etwas funktioniert nicht. Wo fange ich an?

```bash
wails3 doctor
```

Das Tool überprüft Ihre Toolchain, führt fehlende Abhängigkeiten samt Installationsbefehlen auf und gibt die Versionsinformationen aus, die Sie jedem Fehlerbericht beifügen sollten.

### Mein Build schlägt fehl

Probieren Sie diese üblichen Lösungen in der angegebenen Reihenfolge:

1. `go mod tidy`
2. `cd frontend && npm install` (eine fehlende `node_modules` ist die häufigste Ursache)
3. Aktualisieren Sie die CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. Prüfen Sie unter Linux mit `wails3 doctor`, ob GTK-/WebKit-Pakete fehlen

### Meine Bindings fehlen oder sind veraltet

```bash
wails3 generate bindings
```

Im Entwicklungsmodus werden Bindings automatisch neu generiert. Wenn Sie außerhalb von `wails3 dev` einen neuen Service hinzugefügt oder Methodensignaturen geändert haben, generieren Sie sie manuell neu.

### Events werden nicht ausgelöst

Die Event-Namen in `app.Event.Emit("name", ...)` in Go und `Events.On("name", ...)` in JavaScript müssen exakt übereinstimmen. Prüfen Sie zuerst auf Tippfehler und Unterschiede bei der Groß- und Kleinschreibung.

### Ich habe einen Fehler gefunden

Bitte [erstellen Sie ein Issue](https://github.com/wailsapp/wails/issues) und fügen Sie die Ausgabe von `wails3 doctor` bei. Der [Leitfaden für Feedback](/feedback/) erklärt, was einen Bericht leicht umsetzbar macht.

## Migration von v2

### Sollte ich von v2 zu v3 migrieren?

v3 bietet Unterstützung für mehrere Fenster, eine übersichtlichere servicebasierte API, einen integrierten Updater, ein wesentlich flexibleres Build-System und eine bessere Performance. Neue Projekte sollten mit v3 beginnen. Für bestehende Projekte führt der [Migrationsleitfaden](/migration/v2-to-v3/) durch die Unterschiede.

### Wird v2 weiterhin gepflegt?

Ja. Wails v2 erhält weiterhin kritische Updates. Neue Projekte sollten Wails v3.0.0 verwenden.

### Kann ich v2 und v3 parallel verwenden?

Ja. Die CLIs sind separate Binärdateien (`wails` und `wails3`), und die Module haben unterschiedliche Importpfade. Daher können Projekte mit verschiedenen Hauptversionen problemlos auf demselben Rechner nebeneinander bestehen.

## Community

### Wie erhalte ich Hilfe?

- [Discord](https://discord.gg/JDdSxwjhGf) für kurze Fragen und Diskussionen
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) für ausführlichere Fragen
- [GitHub Issues](https://github.com/wailsapp/wails/issues) für Fehlerberichte

### Wie kann ich beitragen?

Lesen Sie den [Leitfaden für Beiträge](/contributing/). Fehlerbehebungen sind jederzeit willkommen. Für neue Funktionen und Änderungen am öffentlichen Verhalten ist ein PR mit einem Entwurf für einen [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) erforderlich. Informelle Diskussionen auf Discord oder in GitHub Discussions sind optional.

### Wo finde ich Beispiele?

Das Repository enthält mehr als 60 ausführbare Beispiele zu Fenstern, Dialogen, Events, der Taskleiste, Services und weiteren Themen: [v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).

## Noch Fragen?

Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder [eröffnen Sie eine Diskussion](https://github.com/wailsapp/wails/discussions).
