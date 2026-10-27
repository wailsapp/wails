---
title: "数据库集成"
description: "将 SQLite、PostgreSQL 和其他数据库与 Wails 集成v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
此页面为占位页面。完整的数据库集成内容即将推出。

@end

Wails v3应用是 Go 应用，因此你可以在绑定服务中使用任何标准 Go 数据库驱动，包括[SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx)和[MySQL](https://github.com/go-sql-driver/mysql)。

[TODO 教程](/tutorials/02-todo-vanilla/)演示了使用内存存储的绑定服务，但未实现 SQLite 持久化。有关基于文件的持久化，请参阅[Notes 教程](/tutorials/03-notes-vanilla/)。此占位页面不涵盖特定于数据库的连接管理、迁移和事务。
