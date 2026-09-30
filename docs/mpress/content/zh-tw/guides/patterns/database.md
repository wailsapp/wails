---
title: "資料庫整合"
description: "將 SQLite、PostgreSQL 及其他資料庫與 Wails v3 整合"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
此頁面目前為預留頁面。完整的資料庫整合內容即將推出。

@end

Wails v3 應用程式就是 Go 應用程式，因此你可以在繫結的服務中使用任何標準 Go 資料庫驅動程式，包括 [SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx) 和 [MySQL](https://github.com/go-sql-driver/mysql)。

如需可實際運作的範例，請參閱[待辦事項教學](/tutorials/02-todo-vanilla/)，其中示範如何為完整的應用程式加入 SQLite 持久化功能。
