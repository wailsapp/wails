---
title: "Foundation"
description: "A production-ready Wails v3 desktop app scaffold"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Foundation layout settings with nine switchable layout presets](/assets/showcase-images/foundation-layouts.webp)

![Foundation home page in the dark theme](/assets/showcase-images/foundation-home-dark.webp)

![Foundation X-Pro charts in the Obsidian theme](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** is a production-ready **desktop app scaffold for Wails v3**, built with **Go, React 19, TypeScript and Mantine 9**. It gives a new project the pieces that usually take weeks to build — a layered Go backend, platform-specific window chrome, persistence, theming, i18n and end-to-end tests — so you can go straight to your business logic.

## Key Highlights

- **Native window on every platform** — a frameless, self-drawn title bar on Windows and Linux; an immersive title bar on macOS that keeps the native traffic lights
- **Nine switchable layouts** — icon rail, sidebar, top nav, merged title bar, right rail, bottom dock, draggable floating dock, inset card and focus mode, all on a CSS Grid shell; register your own preset in one line
- **Themes and i18n** — light, dark and Obsidian presets plus a custom palette editor; English and Chinese out of the box; theme and language stay in sync across the main window, child windows and the tray panel
- **SQLite persistence** — GORM AutoMigrate with no migration files, a user-selectable database location and per-table usage charts
- **Desktop integration** — system tray panel, child windows, native file and message dialogs, and CSV / XLSX export through the native Save As dialog
- **Reliable utilities** — HTTP client, subprocesses bound to a Windows Job Object, AES-GCM encryption, rotating logs and atomic writes
- **Tested end-to-end** — Playwright drives a Wails server-mode build (real Go backend, real frontend), and the same suite generates the README screenshots

Wails keeps it all in a single binary per platform: Go services are called directly through generated bindings, while the frontend follows a strict MVVM structure that stays maintainable as features are added.

Available for Windows, macOS and Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Releases](https://github.com/FxRayHughes/foundation/releases/latest)
