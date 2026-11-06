---
title: "Foundation"
description: "Wails v3.1, React 19, Mantine 9 기반 데스크톱 앱 스캐폴드"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![9가지 레이아웃 프리셋을 전환할 수 있는 Foundation 레이아웃 설정 화면](/assets/showcase-images/foundation-layouts.webp)

![다크 테마의 Foundation 홈 화면](/assets/showcase-images/foundation-home-dark.webp)

![Obsidian 테마의 Foundation X-Pro 차트](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** 프로젝트는 **Go, React 19, TypeScript, Mantine 9**로 만든 **Wails v3.1 데스크톱 앱 스캐폴드**입니다. 계층화된 Go 백엔드, 플랫폼별 창 외형, SQLite 영속성, 테마, 국제화, 엔드 투 엔드 테스트를 새 프로젝트의 출발점으로 제공합니다.

## 기능

- **플랫폼별 창** — Windows와 Linux에서는 프레임 없는 창과 프런트엔드가 직접 그리는 제목 표시줄을, macOS에서는 기본 신호등 버튼을 유지하는 몰입형 제목 표시줄을 사용합니다.
- **9가지 레이아웃 프리셋** — 아이콘 레일, 사이드바, 상단 내비게이션, 제목 표시줄 통합, 오른쪽 레일, 하단 Dock, 드래그 가능한 플로팅 Dock, 인셋 카드, 포커스 모드를 CSS Grid 셸 위에서 제공하며, 사용자 정의 프리셋도 등록할 수 있습니다.
- **테마와 국제화** — 라이트, 다크, Obsidian 프리셋과 사용자 정의 팔레트 편집기를 제공하며 UI는 영어와 중국어를 지원합니다. 메인 창, 하위 창, 트레이 패널의 테마와 언어가 동기화됩니다.
- **SQLite 영속성** — GORM AutoMigrate로 마이그레이션 파일이 필요 없습니다. 데이터베이스 위치는 사용자가 선택할 수 있고 테이블별 사용량을 차트로 확인할 수 있습니다.
- **데스크톱 통합** — 시스템 트레이 패널, 하위 창, 네이티브 파일·메시지 대화상자를 제공하며 네이티브 "다른 이름으로 저장" 대화상자로 CSV / XLSX를 내보낼 수 있습니다.
- **유틸리티** — HTTP 클라이언트, Windows Job Object에 묶인 하위 프로세스, AES-GCM 암호화, 순환 로그, 원자적 쓰기를 제공합니다.
- **엔드 투 엔드 테스트** — Playwright가 Wails 서버 모드 빌드(실제 Go 백엔드와 프런트엔드)를 대상으로 실행되며, 같은 테스트 모음으로 README 스크린샷도 생성합니다.

프런트엔드는 Wails가 생성한 바인딩을 통해 Go 서비스를 호출하며, 각 플랫폼마다 단일 바이너리로 빌드됩니다. 프런트엔드는 MVVM 구조를 따르고, README에는 스캐폴드를 자신의 프로젝트 이름으로 바꾸는 단계별 안내가 있습니다.

Windows, macOS, Linux에서 실행됩니다.

[GitHub](https://github.com/FxRayHughes/foundation) |
[릴리스](https://github.com/FxRayHughes/foundation/releases/latest)
