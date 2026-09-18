---
title: "프로젝트 상태"
description: "Wails v3 베타 호환성, 보안 지원 및 업그레이드 안내"
slug: "status"
sourcePath: "status.md"
---

## 현재 상태: 베타

최신 상태는 [변경 로그](/changelog/)에서 확인하세요.

목표는 안정적인 v3.0 릴리스입니다. Wails v2는 현재 안정 버전으로 유지되며 계속 수정 사항을 받습니다. 배포 전에 애플리케이션에서 베타 릴리스를 테스트하세요.

## 베타 호환성 보장

v3 베타 호환성 계약은 데스크톱 애플리케이션에 적용됩니다.

| 플랫폼 | 지원 대상 | 요구 사항 및 참고 사항 |
| --- | --- | --- |
| Windows | amd64 및 arm64 | WebView2 런타임 |
| macOS | Intel 및 Apple Silicon | 설치 가이드에 명시된 macOS 및 WebKit 버전 |
| Linux | amd64 및 arm64 | 기본값은 GTK4 + WebKitGTK 6.0입니다. GTK3 + WebKit2GTK 4.1는 v3.0.x까지 `-tags gtk3` 레거시 옵션으로 유지되며 v3.1에서 제거됩니다. |

모든 대상을 개발하려면 Go 1.25 이상이 필요합니다. Android 및 iOS 지원은 실험 단계이며 데스크톱 베타 출시를 막지 않습니다. 베타 API는 안정성을 목표로 하지만, 시험판 결함과 명시적으로 공지된 변경 사항은 v3.0.0 이전에도 수정될 수 있습니다.

## 기여 방법

- 최신 베타 릴리스를 테스트하고 재현 가능한 버그를 보고하세요.
- 문서와 예제 작성에 기여하세요.
- 토론에 참여하고 WEP 초안에 대한 피드백을 제공하세요.
- 버그 수정, 문서 또는 승인된 WEP에 대한 풀 리퀘스트를 제출하세요.

커뮤니티의 기여를 환영합니다. 이러한 목표 달성을 돕고 싶다면 커뮤니티 토론에 참여하세요. 새로운 기능에 대한 제안은 기능 요청 이슈가 아니라 WEP PR로 제출해야 합니다.

## 피드백 및 업데이트

재현 가능한 문제는 이슈로 보고하고, 새로운 기능은 [WEP(Wails 개선 제안)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) PR을 통해 제안하세요.

## 베타 사용

`latest`를 추적하는 대신 CLI, Go 모듈, 프런트엔드 런타임의 정확한 버전을 고정하세요. 기존 알파 프로젝트는 [알파에서 베타로 업그레이드하는 안내서](/migration/alpha-to-beta/)를 따르세요.

[보안 정책](https://github.com/wailsapp/wails/blob/master/SECURITY.md)은 v3 베타 릴리스를 지원 대상으로, 알파 릴리스를 지원 대상이 아닌 것으로 명시합니다. 취약점은 공개 이슈가 아닌 [비공개 취약점 보고](https://github.com/wailsapp/wails/security/advisories/new)를 통해 신고하세요.

## 추적 중인 작업

- [v3 라벨이 있는 열린 버그](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [P0 또는 P1 라벨이 있는 열린 v3 이슈](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [릴리스 마일스톤](https://github.com/wailsapp/wails/milestones)

이 실시간 검색은 이슈 라벨에 의존하며, 전체 목록이나 릴리스 날짜 또는 범위에 대한 약속이 아닙니다. 이슈를 읽고 프로젝트에 미치는 영향을 평가하세요.
