---
title: "Foundation"
description: "프로덕션에 바로 쓸 수 있는 Wails v3 데스크톱 앱 스캐폴드"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![9가지 레이아웃 프리셋을 전환할 수 있는 Foundation 레이아웃 설정 화면](/assets/showcase-images/foundation-layouts.webp)

![다크 테마의 Foundation 홈 화면](/assets/showcase-images/foundation-home-dark.webp)

![Obsidian 테마의 Foundation X-Pro 차트](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** 프로젝트는 **Go, React 19, TypeScript, Mantine 9**로 만든, 프로덕션에 바로 쓸 수 있는 **Wails v3 데스크톱 앱 스캐폴드**입니다. 계층화된 Go 백엔드, 플랫폼별 창 외형, 영속성, 테마, 국제화, 엔드 투 엔드 테스트처럼 보통 구축에 몇 주가 걸리는 부분을 미리 갖추고 있어 곧바로 비즈니스 로직 개발을 시작할 수 있습니다.

## 주요 특징

- **모든 플랫폼에서 네이티브 창** — Windows와 Linux에서는 프런트엔드가 직접 그리는 프레임 없는 제목 표시줄을, macOS에서는 기본 신호등 버튼을 유지하는 몰입형 제목 표시줄을 사용합니다.
- **전환 가능한 9가지 레이아웃** — 아이콘 레일, 사이드바, 상단 내비게이션, 제목 표시줄 통합, 오른쪽 레일, 하단 Dock, 드래그 가능한 플로팅 Dock, 인셋 카드, 포커스 모드를 모두 CSS Grid 셸 위에서 제공합니다. 사용자 정의 프리셋도 한 줄로 등록할 수 있습니다.
- **테마와 국제화** — 라이트, 다크, Obsidian 프리셋과 사용자 정의 팔레트 편집기를 제공하며 영어와 중국어를 기본 지원합니다. 메인 창, 하위 창, 트레이 패널의 테마와 언어가 실시간으로 동기화됩니다.
- **SQLite 영속성** — GORM AutoMigrate로 마이그레이션 파일이 필요 없습니다. 데이터베이스 위치는 사용자가 선택할 수 있고 테이블별 사용량을 차트로 확인할 수 있습니다.
- **데스크톱 통합** — 시스템 트레이 패널, 하위 창, 네이티브 파일·메시지 대화상자를 제공하며 네이티브 "다른 이름으로 저장" 대화상자로 CSV / XLSX를 내보낼 수 있습니다.
- **견고한 유틸리티** — HTTP 클라이언트, Windows Job Object에 묶인 하위 프로세스, AES-GCM 암호화, 순환 로그, 원자적 쓰기를 제공합니다.
- **엔드 투 엔드 테스트** — Playwright가 Wails 서버 모드 빌드(실제 Go 백엔드 + 실제 프런트엔드)를 구동하며, 같은 테스트 모음으로 README 스크린샷도 생성합니다.

Wails 덕분에 모든 것이 플랫폼별 단일 바이너리에 담깁니다. Go 서비스는 생성된 바인딩을 통해 직접 호출되고, 프런트엔드는 기능이 늘어나도 유지보수하기 쉬운 엄격한 MVVM 구조를 따릅니다.

Windows, macOS, Linux를 지원합니다.

[GitHub](https://github.com/FxRayHughes/foundation) |
[릴리스](https://github.com/FxRayHughes/foundation/releases/latest)
