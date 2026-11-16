---
title: "Интеграция с базами данных"
description: "Интеграция SQLite, PostgreSQL и других баз данных с Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Эта страница временно содержит заглушку. Полное руководство по интеграции с базами данных появится в ближайшее время.

@end

Приложения Wails v3 являются приложениями Go, поэтому в связанных сервисах можно использовать любой стандартный драйвер баз данных для Go, включая [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx) и [MySQL](https://github.com/go-sql-driver/mysql).

В [руководстве по приложению TODO](/tutorials/02-todo-vanilla/) показан связанный сервис с хранением данных в памяти; постоянное хранение в SQLite в нём не реализовано. О хранении данных в файле см. в [руководстве по приложению Notes](/tutorials/03-notes-vanilla/). Эта страница-заглушка не рассматривает управление подключениями к конкретным базам данных, миграции и транзакции.
