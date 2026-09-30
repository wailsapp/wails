---
title: "データベース統合"
description: "Wails v3 での SQLite、PostgreSQL、その他のデータベースの統合"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
このページはプレースホルダーです。データベース統合に関する詳しい内容は近日公開予定です。

@end

Wails v3 アプリケーションは Go アプリケーションであるため、バインドされたサービス内で、[SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx)、[MySQL](https://github.com/go-sql-driver/mysql) など、標準的な Go データベースドライバーを自由に使用できます。

実際に動作する例については、完全なアプリケーションに SQLite による永続化を追加する方法を示す[TODO チュートリアル](/tutorials/02-todo-vanilla/)を参照してください。
