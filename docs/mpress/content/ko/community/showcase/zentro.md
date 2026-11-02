---
title: "Zentro"
description: "PostgreSQL, MySQL, SQL Server, SQLite를 위한 프로젝트 중심 데스크톱 SQL IDE"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![스키마 인식 자동 완성, 스키마 탐색기, 레코드 패널을 갖춘 Zentro SQL 편집기](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** 는 PostgreSQL, MySQL, SQL Server, SQLite를 위한 **프로젝트 중심 데스크톱 SQL IDE** 이며 **Go, Wails v2, React**로 만들었습니다. 모든 프로젝트는 자신이 어느 환경을 가리키는지 알고 있고, 되돌릴 수 없는 작업을 실행하기 전에 Zentro가 먼저 확인하므로 스테이징용 쿼리가 프로덕션에서 실행되는 일이 없습니다.

## 주요 특징

- **환경 인식** - 프로젝트가 환경(개발, 스테이징, 프로덕션)을 가지며, 프로덕션을 뚜렷하게 표시하고 구문의 위험도를 분류하며 파괴적인 쓰기 전에 영향받는 행을 미리 보여 줍니다
- **SQL 편집기** - 스키마 인식 자동 완성, 인라인 진단, 빠른 수정
- **스키마 탐색기와 ERD** - 테이블, 열, 키, 인덱스를 탐색하고 관계를 시각화합니다
- **Git 기반 추적** - 저장한 스크립트와 스키마 스냅샷을 코드처럼 버전 관리하며 기록과 차이를 볼 수 있습니다
- **SSH 터널과 내장 터미널** - 터널을 통해 데이터베이스에 접속하고 로컬 또는 SSH 터미널을 쿼리 옆에 둘 수 있습니다
- **내 키로 쓰는 AI 어시스턴트** - 직접 설정한 모델 엔드포인트에 연결하며 Zentro가 운영하는 서버를 거치지 않습니다. 외부 에이전트를 위한 MCP 서버도 제공합니다
- **키보드 중심** - 탐색기와 결과 그리드에서의 포커스 이동, 명령 팔레트, 연속 단축키
- **기본적으로 비공개** - 데이터베이스 자격 증명은 운영 체제의 키체인에 보관됩니다

Zentro는 Windows, macOS, Linux용으로 플랫폼마다 하나의 네이티브 바이너리로 제공됩니다. 코어는 GPL-3.0 오픈 소스이며, 선택 사항인 유료 Pro 요금제가 프로덕션 보호 기능 등을 더하고 30일 무료 체험을 제공합니다. Community Edition의 소스 코드는 공개되어 있습니다.

[웹사이트](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[다운로드](https://zentro-ide.com/download)
