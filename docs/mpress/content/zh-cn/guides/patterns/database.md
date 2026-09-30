---
title: "数据库集成"
description: "将 SQLite、PostgreSQL 和其他数据库与 Wails v3 集成"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
此页面为占位页面。完整的数据库集成内容即将推出。

@end

Wails v3 应用程序就是 Go 应用程序，因此你可以在绑定服务中使用任何标准 Go 数据库驱动程序，包括 [SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx) 和 [MySQL](https://github.com/go-sql-driver/mysql)。

有关可运行的示例，请参阅[待办事项教程](/tutorials/02-todo-vanilla/)，其中演示了如何为完整应用程序添加 SQLite 持久化。
