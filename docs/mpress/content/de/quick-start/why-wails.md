---
title: "Warum Wails?"
description: "Erfahren Sie, warum Wails die richtige Wahl für Ihre Desktopanwendung ist"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

Wails verbindet **die Leistung und Einfachheit von Go** mit **der Flexibilität moderner Webbenutzeroberflächen**. So können Sie mit den Ihnen bereits vertrauten Werkzeugen ansprechende, native Desktopanwendungen entwickeln.

## Leistung, die Benutzer bemerken

**Wails-Anwendungen:**

- **Binärdateien mit ca. 15 MB** (gegenüber 150 MB bei Electron)
- **Basisarbeitsspeicherbedarf von ca. 10 MB** (gegenüber mehr als 100 MB bei Electron)
- **&lt;0.5 s Startzeit** (gegenüber 2-3 s bei Electron)
- **Natives Rendering** mit der vom Betriebssystem bereitgestellten WebView

Benutzer erleben Ihre Anwendung als schnell, schlank und professionell.

## Entwicklungserlebnis

**Einmal schreiben, überall ausführen:**

- Eine Go-Codebasis für Windows, macOS und Linux
- Beliebige Webframeworks verwenden (React, Vue, Svelte, Vanilla JS)
- Hot Reload während der Entwicklung
- Automatisch aus Go-Code generierte TypeScript-Bindings

Liefern Sie schneller aus und reduzieren Sie den zu wartenden Code.

## Produktionsreife Funktionen

**Alles, was Sie benötigen:**

- Mehrere Fenster mit unabhängigen Lebenszyklen
- Native Menüs (Anwendungs- und Kontextmenüs sowie Menüs im Infobereich)
- Dateidialoge mit plattformnativer Benutzeroberfläche
- Systemintegration (Benachrichtigungen, Zwischenablage, Tastenkombinationen)
- Codesignierung und Paketierung für alle Plattformen

Entwickeln Sie professionelle Anwendungen statt Prototypen.

## Schnellere Entwicklung

- **Eine Codebasis, drei Plattformen** – Einmal schreiben und für Windows, macOS und Linux erstellen
- **Vorhandene Kenntnisse nutzen** – Go für das Backend, HTML/CSS/JS für die Benutzeroberfläche
- **Sofortiges Feedback** – Hot Reload während der Entwicklung und Kompilierungszeiten von wenigen Sekunden
- **Kleine Binärdateien** – Anwendungen mit 15 MB ermöglichen schnellere Builds, Downloads und Iterationen

## Wann Sie Wails wählen sollten

**Wails eignet sich hervorragend für:**

- **Geschäftsanwendungen** (CRM, Bestandsverwaltung, Dashboards, Administrationswerkzeuge)
- **Entwicklungswerkzeuge** (Datenbankclients, API-Tester, Bereitstellungswerkzeuge)
- **Produktivitätsanwendungen** (Notiz-, Aufgabenverwaltungs- und Zeiterfassungsanwendungen)
- **Kreativwerkzeuge** (Bildbearbeitungsprogramme, Videoverarbeitungsprogramme, Designwerkzeuge)
- **Interne Werkzeuge** (unternehmensspezifische Anwendungen, Automatisierungswerkzeuge)

## Erfolgsgeschichten aus der Praxis

@note{type="tip" title="Produktivanwendungen"}
Wails bildet die Grundlage realer Anwendungen, die von Tausenden Benutzern eingesetzt werden:

- **Datenbankverwaltungswerkzeuge** mit komplexen Benutzeroberflächen
- **Finanz-Dashboards** zur Verarbeitung von Echtzeitdaten
- **Videobearbeitungswerkzeuge** mit nativer Leistung
- **Entwicklungswerkzeuge** für Entwicklungsteams

[Beispiele ansehen →](/community/showcase/)

@end

## Funktionsweise von Wails

Anders als Electron, das einen vollständigen Browser und die Node.js-Laufzeit bündelt, verfolgt Wails einen grundlegend anderen Ansatz: Ihr Go-Code wird zu einer nativen Binärdatei kompiliert, während Ihre Benutzeroberfläche in der integrierten WebView des Betriebssystems ausgeführt wird. Diese Architektur sorgt für die kleinen Binärdateien, den schnellen Start und den geringen Arbeitsspeicherbedarf, durch die sich Wails-Anwendungen nativ anfühlen.

### Architektur

Wails-Anwendungen bestehen aus zwei Hauptkomponenten, die nahtlos miteinander kommunizieren: einem Go-Backend für Geschäftslogik und Systemoperationen sowie einem webbasierten Frontend für Ihre Benutzeroberfläche. Die vom Betriebssystem bereitgestellte WebView rendert Ihre Benutzeroberfläche, ohne einen Browser zu bündeln. Gleichzeitig ermöglicht die Binding-Schicht die typsichere Kommunikation zwischen Go und JavaScript.

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Wails-Architektur: ein Go-Backend und Ihre Web-Benutzeroberfläche, kompiliert zu einer einzigen nativen Binärdatei, verbunden durch generierte Bindings und gerendert in der WebView des Betriebssystems" style="max-width: 640px; width: 100%;" />
</div>

Dank dieser einfachen Architektur kann JavaScript-Code Go-Funktionen direkt über automatisch generierte Bindings aufrufen, während Go Ereignisse und Daten an das Frontend zurücksenden kann. Beide Schichten kommunizieren über eine effiziente In-Memory-Bridge mit einem Overhead von weniger als einer Millisekunde.

**So erzielt Wails seine hohe Leistung:**

1. **Keine gebündelte Laufzeit** – Verwendet die kompilierte Go-Binärdatei
2. **Native WebView** – Vom Betriebssystem bereitgestellte Rendering-Engine
3. **Direkte Go-↔-JS-Bridge** – In-Memory-Kommunikation ohne Netzwerk-Overhead
4. **Kompilierte Binärdatei** – Sofortiger Start ohne JIT-Kompilierung

## Nächste Schritte

Nachdem Sie nun wissen, was Wails bietet, richten wir Ihre Entwicklungsumgebung ein:

1. **Wails installieren** – Richten Sie Ihre Entwicklungsumgebung in 5 Minuten ein [Installationsanleitung →](/quick-start/installation/)

2. **Ihre erste App erstellen** – Erstellen Sie eine funktionsfähige Anwendung und lernen Sie die Grundlagen kennen [Tutorial zur ersten App →](/quick-start/first-app/)

3. **Funktionen erkunden** – Entdecken Sie, was Wails für Ihre Anwendung leisten kann [Funktionsübersicht →](/quick-start/next-steps/)

---

**Haben Sie noch Fragen?** Treten Sie unserer [Discord-Community](https://discord.gg/JDdSxwjhGf) bei und fragen Sie das Team direkt.
