---
title: "CFN Tracker"
description: "Eine mit Wails entwickelte Desktopanwendung"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) – Verfolge die laufenden Matches beliebiger CFN-Profile für Street Fighter 6 oder V. Besuche zum Einstieg [die Website](https://cfn.williamsjokvist.se/).

## Funktionen

- Match-Verfolgung in Echtzeit
- Speicherung von Match-Protokollen und Statistiken
- Unterstützung für die Anzeige von Live-Statistiken in OBS über eine Browserquelle
- Unterstützung für SF6 und SFV
- Benutzer können mit CSS eigene Themes für OBS-Browserquellen erstellen

### Wichtige Technologien neben Wails

- [Task](https://github.com/go-task/task) – umschließt die Wails-CLI, damit häufig verwendete Befehle einfach ausgeführt werden können
- [React](https://github.com/facebook/react) – aufgrund seines umfangreichen Ökosystems ausgewählt (radix, framer-motion)
- [Bun](https://github.com/oven-sh/bun) – wird wegen seiner schnellen Abhängigkeitsauflösung und kurzen Build-Zeit verwendet
- [Rod](https://github.com/go-rod/rod) – Headless-Browser-Automatisierung für die Authentifizierung und das regelmäßige Abrufen von Änderungen
- [SQLite](https://github.com/mattn/go-sqlite3) – wird zum Speichern von Matches, Sitzungen und Profilen verwendet
- [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) – ein HTTP-Stream, der Aktualisierungen der Match-Verfolgung an OBS-Browserquellen sendet
- [i18next](https://github.com/i18next/) – mit Backend-Connector zur Bereitstellung von Lokalisierungsobjekten aus der Go-Schicht
- [xstate](https://github.com/statelyai/xstate) – Zustandsautomaten für die Authentifizierung und Match-Verfolgung
