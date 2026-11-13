---
title: "Foundation"
description: "Desktop app scaffold for Wails v3.1 with React 19 and Mantine 9"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Foundation layout settings with nine switchable layout presets](/assets/showcase-images/foundation-layouts.webp)

![Foundation home page in the dark theme](/assets/showcase-images/foundation-home-dark.webp)

![Foundation X-Pro charts in the Obsidian theme](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** is a **desktop app scaffold for Wails v3.1**, built with **Go, React 19, TypeScript and Mantine 9**. It provides a layered Go backend, platform-specific window chrome, SQLite persistence, themes, i18n and end-to-end tests as a starting point for new projects.

## Features

- **Platform-specific windows** — a frameless window with a frontend-drawn title bar on Windows and Linux; on macOS, an immersive title bar that keeps the native traffic lights
- **Nine layout presets** — icon rail, sidebar, top nav, merged title bar, right rail, bottom dock, draggable floating dock, inset card and focus mode on a CSS Grid shell; custom presets can be registered
- **Themes and i18n** — light, dark and Obsidian presets and a custom palette editor; English and Chinese UI; theme and language are kept in sync across the main window, child windows and the tray panel
- **SQLite persistence** — GORM AutoMigrate without migration files, a user-selectable database location and per-table usage charts
- **Desktop integration** — system tray panel, child windows, native file and message dialogs, and CSV / XLSX export through the native Save As dialog
- **Utilities** — HTTP client, subprocesses bound to a Windows Job Object, AES-GCM encryption, rotating logs and atomic writes
- **End-to-end tests** — Playwright runs against a Wails server-mode build (real Go backend and frontend); the same suite generates the README screenshots

Go services are called from the frontend through Wails-generated bindings, and each platform builds to a single binary. The frontend uses an MVVM structure, and the README includes a step-by-step guide for renaming the scaffold into your own project.

Runs on Windows, macOS and Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Releases](https://github.com/FxRayHughes/foundation/releases/latest)
