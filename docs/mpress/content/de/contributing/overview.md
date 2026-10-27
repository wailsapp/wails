---
title: "Technischer Überblick"
description: "Architekturübersicht und Wegweiser durch die Wails-v3-Codebasis"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Willkommen bei der technischen Dokumentation zu Wails v3

In diesem Abschnitt geht es **nicht** um Community-Richtlinien oder darum, wie Sie einen Pull Request erstellen. Stattdessen wird ausführlich erläutert, **wie Wails v3 aufgebaut ist**, damit Sie sich schnell in der Codebasis zurechtfinden und souverän mit der Entwicklung beginnen können.

Ganz gleich, ob Sie die Runtime korrigieren, die CLI erweitern, neue Vorlagen erstellen oder lediglich die Interna verstehen möchten: Die folgenden Seiten liefern Ihnen den erforderlichen technischen Kontext.

---

## Architekturübersicht

@cards{cols="2"}
◇ Go-Backend
Das Herzstück jeder Wails-App ist Go-Code, der zu einer nativen ausführbaren Datei kompiliert wird. Das Backend übernimmt die Anwendungslogik, die Systemintegration und leistungskritische Operationen.

---
▤ Web-Frontend
Die Benutzeroberfläche wird mit gängigen Webtechnologien (React, Vue, Svelte, Vanilla, …) erstellt und in einer schlanken WebView des Systems gerendert (WebKit unter Linux/macOS, WebView2 unter Windows).

---
◆ Vermittlungsschicht
Eine speicherinterne Zero-Copy-Bridge ermöglicht **Go⇄JavaScript**-Aufrufe mit automatischer Typkonvertierung, Ereignisweitergabe und Fehlerweiterleitung.

---
▸ CLI und Werkzeuge
`wails3` koordiniert die Projekterstellung, den Entwicklungsserver mit Live-Reload, die Bündelung von Assets, die Cross-Kompilierung und die Paketierung (deb, rpm, AppImage, msi, dmg …).

@end

---

## Architektur im Überblick

**Wails v3 – End-to-End-Ablauf**

**[Platzhalter für das End-to-End-Ablaufdiagramm]**

Das Diagramm zeigt den **End-to-End-Ablauf**:

1. Die **CLI** steuert die Generierung, den Entwicklungsserver, die Kompilierung und die Paketierung.\
2. Das **Binding-System** erzeugt Verbindungscode, über den das **Web-Frontend** das **Go-Backend** aufrufen kann.\
3. Während der Entwicklung leitet der **Asset-Server** Anfragen an den Entwicklungsserver des Frameworks weiter; in der Produktion stellt er eingebettete Dateien bereit.\
4. Zur Laufzeit verwaltet die **Desktop-Runtime** Fenster und Betriebssystem-APIs, während die **Bridge** Nachrichten zwischen Go und JavaScript überträgt.

---

## Inhalte dieser Dokumentation

| Thema | Warum es wichtig ist |
| --- | --- |
| **Struktur der Codebasis** | Übersicht über die `/v3`-Verzeichnisse und das Zusammenspiel der Module. |
| **Interna der Runtime** | Fensterverwaltung, System-APIs, Nachrichtenprozessor und plattformspezifische Anpassungsschichten. |
| **Asset- und Entwicklungsserver** | Wie Web-Assets während der Entwicklung bereitgestellt und für die Produktion eingebettet werden. |
| **Build- und Paketierungs-Pipeline** | Taskfile-basierter Workflow, plattformübergreifende Kompilierung und Erstellung von Installationspaketen. |
| **Binding-System** | Pipeline zur statischen Analyse, die typsichere Go⇄TS-Bindings generiert. |
| **Vorlagensystem** | Generatorarchitektur, auf der `wails3 init -t <framework>` basiert. |
| **Tests und CI** | Testumgebung für Unit- und Integrationstests, GitHub Actions und Hinweise zum Race Detector. |
| **Wails erweitern** | Dienste, Vorlagen oder CLI-Unterbefehle hinzufügen. |

Jede der folgenden Seiten vertieft diese Bereiche anhand konkreter Codebeispiele, Diagramme und Verweise auf die relevanten Quelldateien.

---

@note{type="info"}
Voraussetzungen: Sie sollten mit **Go 1.26+**, den Grundlagen von TypeScript und modernen Frontend-Build-Werkzeugen vertraut sein. Wenn Go für Sie neu ist, sollten Sie zunächst die offizielle Tour überfliegen.

@end

Viel Freude beim Erkunden – und willkommen in den Interna von Wails v3!
