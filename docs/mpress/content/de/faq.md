---
title: "Häufig gestellte Fragen"
description: "Antworten auf die häufigsten Fragen zur Anwendungsentwicklung mit Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Allgemeines

### Was ist Wails?

Wails ist ein Framework zur Entwicklung von Desktopanwendungen mit Go und Webtechnologien. Sie schreiben die Anwendungslogik in Go und erstellen die Benutzeroberfläche mit HTML, CSS und JavaScript (oder einem beliebigen Frontend-Framework); Wails rendert sie in der nativen WebView des Betriebssystems. Das Ergebnis ist eine kleine, schnelle Anwendung mit nativem Bediengefühl: kein mitgelieferter Browser, geringer Speicherbedarf und eine einzelne Binärdatei, die typischerweise etwa 10 MB groß ist.

### Welche Plattformen unterstützt Wails?

| Plattform | Anforderungen |
| --- | --- |
| Windows | AMD64 und ARM64. Verwendet die [WebView2-Laufzeitumgebung](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15+ auf Intel (Anwendungen können 10.13+ als Zielversion verwenden), 11.0+ auf Apple Silicon. Universal-Binärdateien werden unterstützt. |
| Linux | AMD64 und ARM64. Der Standard-Stack besteht aus GTK4 mit WebKitGTK 6.0 (Ubuntu 24.04+, Debian 13+, Fedora 40+ und vergleichbare Distributionen). Distributionen, die nur WebKit2GTK 4.1 bereitstellen, etwa Ubuntu 22.04, Debian 12 und RHEL 9, werden über den Legacy-Build `-tags gtk3` unterstützt (verfügbar bis v3.1). Distributionen, die ausschließlich WebKit2GTK 4.0 bereitstellen, werden nicht unterstützt. Siehe den [Linux-Build-Leitfaden](/guides/build/linux/). |
| iOS und Android | Experimentell. Siehe die [Leitfäden für Mobilgeräte](/guides/mobile/). |

Mit dem [Server-Build](/guides/server-build/) können Sie Ihre Anwendung auch als reguläre Webanwendung bereitstellen.

Führen Sie jederzeit `wails3 doctor` aus, um Ihr System zu überprüfen und plattformspezifische Installationsanweisungen zu erhalten.

### Was benötige ich für den Einstieg?

- Go 1.25 oder neuer
- Node.js und npm (für den Frontend-Build)
- Plattform-Toolchain: WebView2 unter Windows (auf 10/11 vorinstalliert), Xcode Command Line Tools unter macOS, `gcc` sowie GTK-/WebKit-Entwicklungspakete unter Linux

`wails3 doctor` prüft all dies für Sie und teilt Ihnen genau mit, was fehlt. Eine vollständige Anleitung finden Sie unter [Installation](/quick-start/installation/).

### Ist Wails v3 produktionsreif?

Wails v3 ist Beta-Software mit einer stabilen Desktop-API. Anwendungen werden damit bereits produktiv betrieben, doch während wir den letzten Feinschliff für 3.0 abschließen, sollten Sie vor der Bereitstellung gründlich testen. Den aktuellen Stand finden Sie auf der [Projektstatusseite](/status/). Wails v2 ist die derzeitige stabile Version und erhält weiterhin Fehlerkorrekturen.

## Entwicklung

### Muss ich Go beherrschen?

Grundkenntnisse in Go sind hilfreich, Sie müssen jedoch kein Experte sein. Die Anwendungslogik befindet sich in einfachen Go-Methoden, und die [Tutorials](/tutorials/overview/) führen Sie durch alles Weitere. Viele Entwickler lernen Go bei der Entwicklung ihrer ersten Wails-Anwendung.

### Kann ich mein bevorzugtes Frontend-Framework verwenden?

Ja. Wenn es sich zu HTML, CSS und JavaScript bauen lässt, funktioniert es mit Wails. Vorlagen sind für React, Vue, Svelte und Vanilla JavaScript enthalten (jeweils auch als TypeScript-Variante); alles andere lässt sich in wenigen Minuten einbinden. Siehe [Frontend-Frameworks](/guides/dev/frontend-frameworks/).

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

Bindings werden während `wails3 dev` automatisch oder bei Bedarf mit `wails3 generate bindings` neu generiert. Siehe [Dienste](/features/bindings/services/).

### Kann ich TypeScript verwenden?

Ja. Der Binding-Generator erzeugt TypeScript-Definitionen für Ihre Dienste und deren Typen, sodass Aufrufe von Go-Funktionen vollständig typisiert sind.

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

Ereignisnamen müssen exakt übereinstimmen. Siehe die [Ereignisreferenz](/guides/events-reference/).

### Wie debugge ich meine Anwendung?

Führen Sie `wails3 dev` aus und klicken Sie mit der rechten Maustaste in das Fenster, um die Browser-Entwicklertools genau wie im Web zu öffnen. Der Entwicklungsserver unterstützt außerdem Hot Reload für Ihr Frontend. Siehe [Debugging](/guides/dev/debugging/).

## Build und Distribution

### Wie erstelle ich einen Produktions-Build?

```bash
wails3 build
```

Ihre Binärdatei wird in `bin/` abgelegt. Produktions-Builds verwenden bereits sinnvolle Standardeinstellungen (Build-Tags, `-trimpath`, entfernte Symbole), sodass für eine schlanke Binärdatei keine zusätzlichen Flags erforderlich sind.

### Kann ich plattformübergreifend kompilieren?

Mit Einschränkungen. Reine Go-Cross-Kompilierung ist nicht möglich, da jede Plattform native WebView-Bibliotheken verwendet. Häufige Anwendungsfälle werden jedoch gut unterstützt:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

Für Builds für Linux auf einem anderen Betriebssystem wird eine Docker-basierte Toolchain verwendet. Die vollständige Matrix finden Sie unter [Plattformübergreifende Builds](/guides/build/cross-platform/).

### Wie erstelle ich ein Installationsprogramm oder Paket?

```bash
wails3 package
```

Dadurch wird das native Format der jeweiligen Plattform erzeugt. Der [Leitfaden zu Installationsprogrammen](/guides/installers/) behandelt NSIS unter Windows, `.app`-Bundles und DMGs unter macOS sowie Linux-Pakete.

### Wie signiere ich den Code meiner Anwendung?

Die Signierung unter Windows und macOS einschließlich der Beglaubigung wird im [Signierungsleitfaden](/guides/build/signing/) Schritt für Schritt beschrieben.

## Funktionen

### Kann ich mehrere Fenster erstellen?

Ja, v3 unterstützt mehrere Fenster nativ:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Siehe [Mehrere Fenster](/features/windows/multiple/).

### Unterstützt Wails den Infobereich des Systems?

Ja, einschließlich Menüs und Klick-Handlern:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Siehe [Infobereich des Systems](/features/menus/systray/).

### Kann ich native Dialoge verwenden?

Ja. Datei-, Meldungs- und Fragedialoge verwenden jeweils die nativen Implementierungen:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Siehe [Dialoge](/features/dialogs/overview/).

### Unterstützt Wails automatische Updates?

Ja. Wails v3 enthält einen integrierten Self-Updater (`app.Updater`) mit austauschbaren Providern für GitHub Releases, keygen.sh und Sparkle AppCast, kryptografischer Signaturprüfung sowie einer Standardbenutzeroberfläche, die Sie anpassen oder ersetzen können. Siehe den Leitfaden zum [In-App-Updater](/guides/updater/) und das Tutorial zur [selbstaktualisierenden Wails-App](/tutorials/04-self-update-a-wails-app/).

## Fehlerbehebung

### Etwas funktioniert nicht. Wo fange ich an?

```bash
wails3 doctor
```

Das Tool überprüft Ihre Toolchain, führt fehlende Abhängigkeiten mit Installationsbefehlen auf und gibt die Versionsinformationen aus, die Sie jedem Fehlerbericht beifügen sollten.

### Mein Build schlägt fehl

Probieren Sie der Reihe nach die üblichen Lösungen aus:

1. `go mod tidy`
2. `cd frontend && npm install` (ein fehlendes Verzeichnis `node_modules` ist die häufigste Ursache)
3. Aktualisieren Sie die CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. Prüfen Sie unter Linux mit `wails3 doctor`, ob GTK-/WebKit-Pakete fehlen

### Meine Bindings fehlen oder sind veraltet

```bash
wails3 generate bindings
```

Im Entwicklungsmodus werden Bindings automatisch neu generiert. Wenn Sie außerhalb von `wails3 dev` einen neuen Service hinzugefügt oder Methodensignaturen geändert haben, generieren Sie die Bindings manuell neu.

### Ereignisse werden nicht ausgelöst

Die Ereignisnamen in `app.Event.Emit("name", ...)` in Go und `Events.On("name", ...)` in JavaScript müssen exakt übereinstimmen. Prüfen Sie zuerst, ob Tippfehler oder Unterschiede bei der Groß- und Kleinschreibung vorliegen.

### Ich habe einen Fehler gefunden

Bitte [eröffnen Sie ein Issue](https://github.com/wailsapp/wails/issues) und fügen Sie die Ausgabe von `wails3 doctor` hinzu. Der [Leitfaden für Feedback](/feedback/) erklärt, welche Angaben nötig sind, damit ein Bericht leicht bearbeitet werden kann.

## Migration von v2

### Sollte ich von v2 zu v3 migrieren?

v3 bietet Unterstützung für mehrere Fenster, eine übersichtlichere servicebasierte API, einen integrierten Updater, ein wesentlich flexibleres Build-System und eine bessere Performance. Neue Projekte sollten mit v3 beginnen. Für bestehende Projekte erläutert der [Migrationsleitfaden](/migration/v2-to-v3/) die Unterschiede Schritt für Schritt.

### Wird v2 weiterhin gepflegt?

Ja. v2 erhält weiterhin Fehlerbehebungen, während v3 auf die stabile Veröffentlichung hinarbeitet.

### Kann ich v2 und v3 parallel verwenden?

Ja. Die CLIs sind separate Binärdateien (`wails` und `wails3`), und die Module haben unterschiedliche Importpfade. Daher können Projekte mit verschiedenen Hauptversionen problemlos auf demselben Rechner nebeneinander bestehen.

## Community

### Wie erhalte ich Hilfe?

- [Discord](https://discord.gg/JDdSxwjhGf) für kurze Fragen und Diskussionen
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) für ausführlichere Fragen
- [GitHub Issues](https://github.com/wailsapp/wails/issues) für Fehler

### Wie kann ich beitragen?

Siehe den [Leitfaden für Beiträge](/contributing/). Fehlerbehebungen sind jederzeit willkommen. Neue Funktionen und Änderungen am öffentlichen Verhalten werden über einen als Entwurf markierten PR für einen [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) eingereicht. Eine informelle Diskussion auf Discord oder in GitHub Discussions ist optional.

### Wo finde ich Beispiele?

Das Repository enthält über 60 ausführbare Beispiele zu Fenstern, Dialogen, Ereignissen, dem Infobereich des Systems, Services und weiteren Themen: [v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples).

## Noch Fragen?

Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder [eröffnen Sie eine Diskussion](https://github.com/wailsapp/wails/discussions).
