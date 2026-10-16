---
title: "Integração com bancos de dados"
description: "Integração do SQLite, PostgreSQL e outros bancos de dados com o Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Esta página é provisória. O conteúdo completo sobre integração com bancos de dados estará disponível em breve.

@end

Os aplicativos Wails v3 são aplicativos Go, portanto, você pode usar qualquer driver de banco de dados padrão do Go — incluindo [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx) e [MySQL](https://github.com/go-sql-driver/mysql) — nos serviços vinculados.

Para ver um exemplo funcional, consulte o [tutorial de TODO](/tutorials/02-todo-vanilla/), que demonstra como adicionar persistência com SQLite a um aplicativo completo.
