---
title: "데이터베이스 통합"
description: "Wails v3에서 SQLite, PostgreSQL 및 기타 데이터베이스 통합하기"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
이 페이지는 자리 표시자입니다. 전체 데이터베이스 통합 콘텐츠는 곧 제공될 예정입니다.

@end

Wails v3 애플리케이션은 Go 애플리케이션이므로 바인딩된 서비스 내에서 [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx), [MySQL](https://github.com/go-sql-driver/mysql)을 비롯한 모든 표준 Go 데이터베이스 드라이버를 사용할 수 있습니다.

실제로 작동하는 예제는 완성된 애플리케이션에 SQLite 영속성을 추가하는 방법을 보여 주는 [TODO 튜토리얼](/tutorials/02-todo-vanilla/)을 참조하세요.
