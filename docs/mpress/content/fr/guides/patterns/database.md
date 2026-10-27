---
title: "Intégration de bases de données"
description: "Intégration de SQLite, PostgreSQL et d’autres bases de données avec Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Cette page est un espace réservé. Le contenu complet sur l’intégration de bases de données sera bientôt disponible.

@end

Les applications Wails v3 sont des applications Go. Vous pouvez donc utiliser n’importe quel pilote de base de données Go standard — notamment [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx) et [MySQL](https://github.com/go-sql-driver/mysql) — dans vos services liés.

Le [tutoriel TODO](/tutorials/02-todo-vanilla/) présente un service lié avec un stockage en mémoire, mais n’implémente pas la persistance SQLite. Pour la persistance dans un fichier, consultez le [tutoriel Notes](/tutorials/03-notes-vanilla/). La gestion des connexions propre à chaque base de données, les migrations et les transactions ne sont pas abordées dans cet espace réservé.
