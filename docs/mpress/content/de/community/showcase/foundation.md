---
title: "Foundation"
description: "Gerüst für Desktop-Apps mit Wails v3.1, React 19 und Mantine 9"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Layout-Einstellungen von Foundation mit neun umschaltbaren Voreinstellungen](/assets/showcase-images/foundation-layouts.webp)

![Startseite von Foundation im dunklen Theme](/assets/showcase-images/foundation-home-dark.webp)

![X-Pro-Diagramme von Foundation im Obsidian-Theme](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** ist ein **Gerüst für Desktop-Apps mit Wails v3.1**, gebaut mit **Go, React 19, TypeScript und Mantine 9**. Als Ausgangspunkt für neue Projekte bietet es ein geschichtetes Go-Backend, plattformspezifische Fenster, SQLite-Persistenz, Themes, Internationalisierung und End-to-End-Tests.

## Funktionen

- **Plattformspezifische Fenster** — unter Windows und Linux ein rahmenloses Fenster mit einer vom Frontend gezeichneten Titelleiste; unter macOS eine immersive Titelleiste, die die nativen Ampel-Schaltflächen behält.
- **Neun Layout-Voreinstellungen** — Icon-Leiste, Seitenleiste, obere Navigation, in die Titelleiste integrierte Navigation, rechte Leiste, unteres Dock, verschiebbares schwebendes Dock, eingerückte Karte und Fokusmodus auf einer CSS-Grid-Shell; eigene Voreinstellungen lassen sich registrieren.
- **Themes und Internationalisierung** — helle, dunkle und Obsidian-Voreinstellungen sowie ein Editor für eigene Farbpaletten; Oberfläche auf Englisch und Chinesisch; Theme und Sprache bleiben zwischen Hauptfenster, Unterfenstern und Tray-Panel synchron.
- **SQLite-Persistenz** — GORM AutoMigrate ohne Migrationsdateien, vom Benutzer wählbarer Speicherort der Datenbank und Diagramme zur Belegung pro Tabelle.
- **Desktop-Integration** — Tray-Panel, Unterfenster, native Datei- und Meldungsdialoge sowie CSV-/XLSX-Export über den nativen „Speichern unter“-Dialog.
- **Hilfsbibliotheken** — HTTP-Client, an ein Windows Job Object gebundene Unterprozesse, AES-GCM-Verschlüsselung, rotierende Logs und atomare Schreibvorgänge.
- **End-to-End-Tests** — Playwright läuft gegen einen Wails-Build im Servermodus (echtes Go-Backend und Frontend); dieselbe Testsuite erzeugt die Screenshots für die README.

Das Frontend ruft die Go-Services über von Wails generierte Bindings auf, und jede Plattform wird zu einer einzigen Binärdatei gebaut. Das Frontend folgt einer MVVM-Struktur, und die README enthält eine Schritt-für-Schritt-Anleitung, um das Gerüst in dein eigenes Projekt umzubenennen.

Läuft unter Windows, macOS und Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Releases auf GitHub](https://github.com/FxRayHughes/foundation/releases/latest)
