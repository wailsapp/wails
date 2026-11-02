---
title: "Zentro"
description: "Project-centric desktop SQL IDE for PostgreSQL, MySQL, SQL Server and SQLite"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![Zentro SQL editor with schema-aware autocomplete, the schema explorer and a record panel](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** is a **project-centric desktop SQL IDE** for PostgreSQL, MySQL, SQL Server and SQLite, built with **Go, Wails v2 and React**. Every project knows which environment it points at, and Zentro asks before anything irreversible runs, so the query meant for staging does not land on production.

## Key Highlights

- **Environment awareness** - projects carry their environments (development, staging, production) with a clear production overlay, risk classification of statements, and an affected-rows preview before a destructive write
- **SQL editor** - schema-aware autocomplete, inline diagnostics and quick fixes
- **Schema explorer and ERD** - browse tables, columns, keys and indexes and visualize relationships
- **Git-backed tracking** - saved scripts and schema snapshots versioned like code, with history and diffs
- **SSH tunnels and an in-app terminal** - reach a database through a tunnel, and keep a local or SSH terminal next to your queries
- **AI assistant with your own key** - connects to the model endpoint you configure, never through a Zentro-run server, plus an MCP server for external agents
- **Keyboard-first** - roving focus through the explorer and result grid, a command palette and chord shortcuts
- **Private by default** - database credentials are kept in your operating system's keychain

Zentro ships as a single native binary per platform for Windows, macOS and Linux. The core is open source under GPL-3.0 and an optional paid Pro tier adds production guardrails and more, with a 30-day free trial; the source of the Community Edition is public.

[Website](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[Download](https://zentro-ide.com/download)
