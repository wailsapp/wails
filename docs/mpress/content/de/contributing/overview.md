---
title: "Technischer Überblick"
description: "Überblick über die Architektur und Wegweiser durch die Wails-v3-Codebasis"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Willkommen bei der technischen Dokumentation zu Wails v3

In diesem Abschnitt geht es **nicht** um Community-Richtlinien oder darum, wie Sie einen Pull-Request erstellen. Stattdessen erfahren Sie, **wie Wails v3 aufgebaut ist**, damit Sie sich schnell in der Codebasis zurechtfinden und sicher mit der Entwicklung beginnen können.

Ganz gleich, ob Sie die Runtime korrigieren, die CLI erweitern, neue Vorlagen erstellen oder einfach die Interna verstehen möchten: Die folgenden Seiten liefern Ihnen den erforderlichen technischen Kontext.

---

## Architekturüberblick

@cards{cols="2"}
◇ Go-Backend
Das Herzstück jeder Wails-App ist Go-Code, der zu einer nativen ausführbaren Datei kompiliert wird. Er übernimmt die Anwendungslogik, die Systemintegration und leistungskritische Operationen.

---
▤ Web-Frontend
Die Benutzeroberfläche wird mit gängigen Webtechnologien (React, Vue, Svelte, Vanilla, …) erstellt und von einem schlanken systemeigenen WebView gerendert (WebKit unter Linux/macOS, WebView2 unter Windows).

---
◆ Brückenschicht
Eine kopierfreie In-Memory-Brücke ermöglicht **Go⇄JavaScript**-Aufrufe mit automatischer Typkonvertierung, Ereignisweitergabe und Fehlerweiterleitung.

---
▸ CLI und Werkzeuge
`wails3` steuert die Projekterstellung, den Entwicklungsserver mit Live-Reload, die Bündelung von Assets, die Cross-Kompilierung und die Paketierung (deb, rpm, AppImage, msi, dmg …).

@end

---

## Architekturübersicht

**Wails v3 – durchgängiger Ablauf**

**[Platzhalter für das Diagramm des durchgängigen Ablaufs]**

Das Diagramm zeigt den **durchgängigen Ablauf**:

1. Die **CLI** steuert die Generierung, den Entwicklungsserver, die Kompilierung und die Paketierung.\
2. Das **Binding-System** erzeugt Verbindungscode, über den das **Web-Frontend** das **Go-Backend** aufrufen kann.\
3. Während der Entwicklung fungiert der **Asset-Server** als Proxy für den Entwicklungsserver des Frameworks; in der Produktion stellt er eingebettete Dateien bereit.\
4. Zur Laufzeit verwaltet die **Desktop-Runtime** Fenster und Betriebssystem-APIs, während die **Brücke** Nachrichten zwischen Go und JavaScript überträgt.

---

## Inhalt dieser Dokumentation

| Thema | Warum es wichtig ist |
| --- | --- |
| **Aufbau der Codebasis** | Übersicht über die `/v3`-Verzeichnisse und das Zusammenspiel der Module. |
| **Interna der Runtime** | Fensterverwaltung, System-APIs, Nachrichtenverarbeitung und plattformspezifische Anpassungsschichten. |
| **Asset- und Entwicklungsserver** | Wie Web-Assets während der Entwicklung bereitgestellt und für die Produktion eingebettet werden. |
| **Build- und Paketierungs-Pipeline** | Taskfile-basierter Workflow, plattformübergreifende Kompilierung und Erstellung von Installationsprogrammen. |
| **Binding-System** | Pipeline für statische Analysen, die typsichere Go⇄TS-Bindings erzeugt. |
| **Vorlagensystem** | Generatorarchitektur, auf der `wails3 init -t <framework>` basiert. |
| **Tests und CI** | Testumgebung für Unit- und Integrationstests, GitHub Actions und Hinweise zum Race Detector. |
| **Wails erweitern** | Dienste, Vorlagen oder CLI-Unterbefehle hinzufügen. |

Jede folgende Seite behandelt diese Bereiche ausführlich anhand konkreter Codebeispiele, Diagramme und Verweise auf die relevanten Quelldateien.

---

@note{type="info"}
Voraussetzungen: Sie sollten mit **Go 1.25+**, grundlegendem TypeScript und modernen Frontend-Build-Werkzeugen vertraut sein. Wenn Go für Sie neu ist, empfiehlt es sich, zunächst die offizielle Tour zu überfliegen.

@end

Viel Freude beim Erkunden – und willkommen bei den Interna von Wails v3!
