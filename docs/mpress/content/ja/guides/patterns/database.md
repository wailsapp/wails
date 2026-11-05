---
title: "データベース統合"
description: "SQLite、PostgreSQL、およびその他のデータベースを Wails v3 と統合する方法"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
このページはプレースホルダーです。データベース統合に関する完全なコンテンツは近日公開予定です。

@end

Wails v3 アプリケーションは Go アプリケーションであるため、バインドされたサービス内で、[SQLite](https://github.com/mattn/go-sqlite3)、[PostgreSQL](https://github.com/jackc/pgx)、[MySQL](https://github.com/go-sql-driver/mysql) を含む任意の標準 Go データベースドライバーを使用できます。

[TODO チュートリアル](/tutorials/02-todo-vanilla/)では、インメモリストレージを使用するバインドされたサービスを紹介していますが、SQLite による永続化は実装していません。ファイルベースの永続化については、[Notes チュートリアル](/tutorials/03-notes-vanilla/)を参照してください。このプレースホルダーでは、データベース固有の接続管理、マイグレーション、およびトランザクションについては説明しません。
