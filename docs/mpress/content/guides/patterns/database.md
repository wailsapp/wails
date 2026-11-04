---
title: "Database Integration"
description: "Integrating SQLite, PostgreSQL, and other databases with Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
This page is a placeholder. Full database integration content is coming soon.

@end

Wails v3 applications are Go applications, so you can use any standard Go database driver — including [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx), and [MySQL](https://github.com/go-sql-driver/mysql) — inside your bound services.

The [TODO tutorial](/tutorials/02-todo-vanilla/) demonstrates a bound service with in-memory storage; it does not implement SQLite persistence. For file-backed persistence, see the [Notes tutorial](/tutorials/03-notes-vanilla/). Database-specific connection management, migrations and transactions are not covered by this placeholder.
