---
title: "Foundation"
description: "Produktionsreifes Gerüst für Wails-v3-Desktop-Apps"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Layout-Einstellungen von Foundation mit neun umschaltbaren Voreinstellungen](/assets/showcase-images/foundation-layouts.webp)

![Startseite von Foundation im dunklen Theme](/assets/showcase-images/foundation-home-dark.webp)

![X-Pro-Diagramme von Foundation im Obsidian-Theme](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** ist ein produktionsreifes **Gerüst für Desktop-Apps mit Wails v3**, gebaut mit **Go, React 19, TypeScript und Mantine 9**. Es liefert einem neuen Projekt die Bausteine, deren Aufbau sonst Wochen dauert — ein geschichtetes Go-Backend, plattformspezifische Fenster, Persistenz, Themes, Internationalisierung und End-to-End-Tests —, damit du direkt mit der Geschäftslogik beginnen kannst.

## Wichtige Funktionen

- **Natives Fenster auf jeder Plattform** — unter Windows und Linux eine rahmenlose, vom Frontend gezeichnete Titelleiste; unter macOS eine immersive Titelleiste, die die nativen Ampel-Schaltflächen behält.
- **Neun umschaltbare Layouts** — Icon-Leiste, Seitenleiste, obere Navigation, in die Titelleiste integrierte Navigation, rechte Leiste, unteres Dock, verschiebbares schwebendes Dock, eingerückte Karte und Fokusmodus, alles auf einer CSS-Grid-Shell; eigene Voreinstellungen lassen sich mit einer Zeile registrieren.
- **Themes und Internationalisierung** — helle, dunkle und Obsidian-Voreinstellungen sowie ein Editor für eigene Farbpaletten; Englisch und Chinesisch sofort verfügbar; Theme und Sprache bleiben zwischen Hauptfenster, Unterfenstern und Tray-Panel synchron.
- **SQLite-Persistenz** — GORM AutoMigrate ohne Migrationsdateien, vom Benutzer wählbarer Speicherort der Datenbank und Diagramme zur Belegung pro Tabelle.
- **Desktop-Integration** — Tray-Panel, Unterfenster, native Datei- und Meldungsdialoge sowie CSV-/XLSX-Export über den nativen „Speichern unter“-Dialog.
- **Robuste Hilfsbibliotheken** — HTTP-Client, an ein Windows Job Object gebundene Unterprozesse, AES-GCM-Verschlüsselung, rotierende Logs und atomare Schreibvorgänge.
- **End-to-End getestet** — Playwright steuert einen Wails-Build im Servermodus (echtes Go-Backend und echtes Frontend), und dieselbe Testsuite erzeugt die Screenshots für die README.

Mit Wails steckt alles in einer einzigen Binärdatei pro Plattform: Die Go-Services werden direkt über generierte Bindings aufgerufen, während das Frontend einer strikten MVVM-Struktur folgt, die auch mit wachsendem Funktionsumfang wartbar bleibt.

Verfügbar für Windows, macOS und Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Releases auf GitHub](https://github.com/FxRayHughes/foundation/releases/latest)
