---
title: "Datenbankintegration"
description: "Integration von SQLite, PostgreSQL und anderen Datenbanken in Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Diese Seite ist ein Platzhalter. Ausführliche Inhalte zur Datenbankintegration folgen in Kürze.

@end

Wails-v3Anwendungen sind Go-Anwendungen. Daher können Sie in Ihren gebundenen Diensten jeden standardmäßigen Go-Datenbanktreiber verwenden, darunter [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx) und [MySQL](https://github.com/go-sql-driver/mysql).

Das [TODO-Tutorial](/tutorials/02-todo-vanilla/) zeigt einen gebundenen Dienst mit In-Memory-Speicherung; SQLite-Persistenz wird darin nicht implementiert. Informationen zur dateibasierten Persistenz finden Sie im [Notes-Tutorial](/tutorials/03-notes-vanilla/). Datenbankspezifische Verbindungsverwaltung, Migrationen und Transaktionen werden in diesem Platzhalter nicht behandelt.
