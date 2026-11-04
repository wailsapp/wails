---
title: "資料庫整合"
description: "將 SQLite、PostgreSQL 和其他資料庫與 Wails v3整合"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
此頁面是預留頁面。完整的資料庫整合內容即將推出。

@end

Wails v3應用程式是 Go 應用程式，因此你可以在繫結的服務中使用任何標準 Go 資料庫驅動程式，包括 [SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx) 和 [MySQL](https://github.com/go-sql-driver/mysql)。

[TODO 教學](/tutorials/02-todo-vanilla/)示範了使用記憶體內儲存空間的繫結服務，但未實作 SQLite 持久化。如需以檔案為後端的持久化，請參閱[筆記教學](/tutorials/03-notes-vanilla/)。此預留頁面不涵蓋特定資料庫的連線管理、遷移和交易。
