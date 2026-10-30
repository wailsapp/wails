---
title: "Integração com bancos de dados"
description: "Integração do SQLite, PostgreSQL e outros bancos de dados com o Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Esta página é um conteúdo provisório. O conteúdo completo sobre integração com bancos de dados estará disponível em breve.

@end

Os aplicativos Wails v3 são aplicativos Go, portanto você pode usar qualquer driver de banco de dados padrão para Go — incluindo [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx) e [MySQL](https://github.com/go-sql-driver/mysql) — nos seus serviços vinculados.

O [tutorial TODO](/tutorials/02-todo-vanilla/) demonstra um serviço vinculado com armazenamento em memória; ele não implementa persistência com SQLite. Para persistência em arquivos, consulte o [tutorial Notes](/tutorials/03-notes-vanilla/). O gerenciamento de conexões, as migrações e as transações específicas de bancos de dados não são abordados neste conteúdo provisório.
