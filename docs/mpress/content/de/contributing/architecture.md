---
title: "Architektur von Wails v3"
description: "Detaillierte Diagramme und Erläuterungen zu allen internen Komponenten von Wails v3"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 ist ein **Full-Stack-Desktop-Framework**, das aus einer Go-Laufzeitumgebung, einer JavaScript-Bridge, einer auf Tasks basierenden Toolchain und einer Sammlung von Vorlagen besteht, mit denen Sie native Anwendungen auf Basis moderner Webtechnologien bereitstellen können.

Diese Seite vermittelt das *Gesamtbild* anhand von vier Diagrammen:

1. **Gesamtarchitektur** – wie alle Subsysteme miteinander verbunden sind\
2. **Laufzeitablauf** – was geschieht, wenn JavaScript Go aufruft und umgekehrt\
3. **Entwicklung und Produktion** – die beiden Betriebsmodi des Asset-Servers\
4. **Plattformspezifische Implementierungen** – wo sich betriebssystemspezifischer Code befindet\

---

## 1 · Gesamtarchitektur

**Wails v3 – High-Level-Stack**

**[Platzhalter für das High-Level-Stack-Diagramm]**

---

## 2 · Ablauf von Laufzeitaufrufen

**Laufzeit – Aufrufpfad zwischen JavaScript und Go**

**[Platzhalter für das Diagramm zum Ablauf von Laufzeitaufrufen]**

Wichtige Punkte:

- **Kein HTTP/IPC** – die Bridge nutzt den speicherinternen Kanal der nativen WebView\
- **Methoden-IDs** – ein deterministischer FNV-Hash ermöglicht Nachschlagen in Go mit O(1)\
- **Promises** – Fehler werden als Rejections mit Stack und Code weitergegeben

---

## 3 · Asset-Ablauf in Entwicklung und Produktion

**Asset-Server: Entwicklung ↔ Produktion**

**[Platzhalter für das Asset-Ablaufdiagramm]**

- Im **Entwicklungsmodus** leitet der Server unbekannte Pfade per Proxy an den Live-Reload-Server des Frameworks weiter und stellt statische Assets vom Datenträger bereit.
- Im **Produktionsmodus** wird dieselbe API durch `go:embed` unterstützt, wodurch eine Binärdatei ohne Abhängigkeiten entsteht.

---

## 4 · Plattformspezifische Aufteilung der Laufzeitumgebung

**Laufzeitdateien je Betriebssystem**

**[Platzhalter für das Diagramm zur Plattformaufteilung]**

Jede Funktion folgt diesem Muster:

1. **Gemeinsame Schnittstelle** in `pkg/application`\
2. **Einstiegspunkt des Nachrichtenprozessors** in `pkg/application/messageprocessor_*.go`\
3. **Betriebssystemspezifische Implementierung** in `pkg/application/*_{darwin,linux,windows}.go` (z. B. `webview_window_darwin.go`, `clipboard_linux.go`, `dialogs_windows.go`, `systemtray_*.go`, `mainthread_*.go`), geschützt durch Build-Tags. Unter Linux gibt es zusätzlich die cgo-Bridge in `linux_cgo.go`/`linux_cgo_gtk4.{go,c,h}`.

`internal/runtime/` enthält lediglich den kleinen `runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go`- Verbindungscode für Build-Tags sowie die unter `internal/runtime/desktop/` eingebettete JavaScript-Laufzeitumgebung.

`internal/capabilities/` dient zur Deklaration plattformspezifischer Funktionsumfänge, aber es gibt keinen `ErrCapability`-Sentinelwert – die Verfügbarkeit von Funktionen wird über einfache Build-Tags und plattformspezifische Stub-Rückgaben gesteuert (z. B. `nil` oder funktionsspezifische Fehler).

---

## Zusammenfassung

Diese Diagramme zeigen, **wo sich der Code befindet**, **wie Daten übertragen werden** und **welche Schichten für welche Aufgaben zuständig sind**. Halten Sie sie griffbereit, während Sie die folgenden Detailseiten erkunden – sie dienen Ihnen als Wegweiser durch den Quellbaum von Wails v3.
