---
title: "XenSQL"
description: "여러 데이터베이스 지원과 고급 쿼리 도구를 갖춘 로컬 우선 SQL 워크벤치"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![인라인 편집 기능을 갖춘 XenSQL 테이블 편집기](/assets/showcase-images/xensql-1.png) ![테이블 및 열 제안 기능을 갖춘 XenSQL 쿼리 편집기](/assets/showcase-images/xensql-2.png) ![XenSQL 다중 트랜잭션 쿼리](/assets/showcase-images/xensql-3.png)

<strong>[XenSQL](https://github.com/Bare7a/XenSQL)</strong>은 <strong>Go, Wails 및 React</strong>로 구축된 <strong>빠른 로컬 우선 SQL 데스크톱 워크벤치</strong>입니다. PostgreSQL, MySQL/MariaDB 및 SQLite를 하나의 깔끔하고 네이티브 앱처럼 자연스러운 인터페이스에 통합했으며, 클라우드, 원격 분석, 계정을 전혀 사용하지 않습니다.

## 주요 특징

- **강력한 SQL 편집기** - Monaco 기반으로, 스키마를 인식하는 스마트 자동 완성, 다중 문 실행, 결과 스트리밍 및 문별 결과 탭을 제공합니다.
- **고급 데이터 뷰어** - 대화형 JSON 검사기, 구문을 인식하는 셀 편집기(JSON, XML, HTML, 텍스트), 인라인 편집 및 전체 레코드 검사를 제공합니다.
- **매끄러운 데이터 편집** - 테이블 탐색, 인라인 변경 사항 스테이징, 일괄 작업 및 `RETURNING` 지원을 통한 안전한 `INSERT`/`UPDATE`/`DELETE` 기능을 제공합니다.
- **생산성 기능** - 스키마 탐색기, 저장된 쿼리, 쿼리 기록, 빠른 검색(`Ctrl+P`) 및 키보드 중심 워크플로를 제공합니다.
- **내보내기 옵션** - CSV, JSON, Markdown, SQL INSERT 문을 지원합니다.

완전히 오프라인으로 작동하며 휴대할 수 있습니다. 모든 데이터는 앱과 함께 이동하는 단일 `XenSQL-data/` 폴더에 로컬로 저장됩니다.

**지원되는 데이터베이스**: PostgreSQL, MySQL, MariaDB 및 SQLite(읽기 전용 모드와 보안 전송 옵션 지원).

기존 SQL 도구의 불필요하게 무거운 구성 없이 속도, 명확성 및 제어 능력을 원하는 개발자를 위해 설계되었습니다.
