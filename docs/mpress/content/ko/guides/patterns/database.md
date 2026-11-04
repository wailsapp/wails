---
title: "데이터베이스 통합"
description: "SQLite, PostgreSQL 및 기타 데이터베이스를 Wails와 통합하기 v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
이 페이지는 자리 표시자입니다. 전체 데이터베이스 통합 콘텐츠는 곧 제공될 예정입니다.

@end

Wails v3 애플리케이션은 Go 애플리케이션이므로 바인딩된 서비스 내에서 [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx), [MySQL](https://github.com/go-sql-driver/mysql)을 비롯한 모든 표준 Go 데이터베이스 드라이버를 사용할 수 있습니다.

[TODO 튜토리얼](/tutorials/02-todo-vanilla/)에서는 메모리 내 저장소를 사용하는 바인딩된 서비스를 보여 주며, SQLite 영속성은 구현하지 않습니다. 파일 기반 영속성에 대해서는 [Notes 튜토리얼](/tutorials/03-notes-vanilla/)을 참조하세요. 이 자리 표시자에서는 데이터베이스별 연결 관리, 마이그레이션 및 트랜잭션을 다루지 않습니다.
