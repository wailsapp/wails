---
title: "변경 로그"
description: "Wails v3의 버전 기록 및 릴리스 정보"
slug: "changelog"
sourcePath: "changelog.md"
---

범례:

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 이 프로젝트의 모든 주요 변경 사항은 이 파일에 기록됩니다.

이 형식은 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)를 기반으로 하며, 이 프로젝트는 [유의적 버전 관리](https://semver.org/spec/v2.0.0.html)를 따릅니다.

- 새 기능은 `Added`에 기록합니다.
- 기존 기능의 변경 사항은 `Changed`에 기록합니다.
- 곧 제거될 기능은 `Deprecated`에 기록합니다.
- 제거된 기능은 `Removed`에 기록합니다.
- 모든 버그 수정은 `Fixed`에 기록합니다.
- 취약점은 `Security`에 기록합니다.

_/

/_   * 이 파일을 업데이트하지 마세요 *   업데이트는 `v3/UNRELEASED_CHANGELOG.md`에 추가하세요.   감사합니다! _/

## [미출시]

## v3.0.0-beta.21 - 2026-09-13

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/6116)에서 Wails v3 문서를 M-Press로 제공 - @leaanthony

## 수정됨

- 변경 로그 생성을 위해 MPD 프런트매터의 JSON 슬러그 값을 파싱하도록 수정([PR](https://github.com/wailsapp/wails/pull/6118), @leaanthony)
- 백업 실패 후 업데이터가 헬퍼 환경 변수를 지우고 원래 대상을 다시 실행하도록 수정([PR](https://github.com/wailsapp/wails/pull/6080), @cnmax)
- App.Run 중 기본 시그널 핸들러를 시작하도록 수정([PR](https://github.com/wailsapp/wails/pull/6098), @leaanthony)
- Windows 메뉴가 nil 메뉴를 처리하고, 교체된 리소스를 해제하며, 메뉴 모음을 다시 그리도록 수정([PR](https://github.com/wailsapp/wails/pull/6112), @taliesin-ai)
- 공유 YAML 구성을 사용하는 새 프로젝트의 MSIX 패키징을 복원([PR](https://github.com/wailsapp/wails/pull/6115), @leaanthony)
- 제네릭 모델 생성기가 뒤에 선언된 헬퍼를 참조할 때 생성된 JavaScript 및 TypeScript 바인딩을 불러오지 못하는 문제를 수정하고, 상호 의존적인 제네릭 모델을 생성할 때 스택 오버플로가 발생하지 않도록 수정(#6062)

## v3.0.0-beta.20 - 2026-09-10

## 변경

- @01xR4in이 [PR](https://github.com/wailsapp/wails/pull/6082)에서 Clave 쇼케이스 링크를 현재 웹사이트와 저장소로 업데이트했습니다.

## 수정

- 탐색 간에 keepalive 핸들러는 유지하면서 worker 요청을 포함하여 중단된 Windows 에셋 요청을 취소합니다. Apple 플랫폼에서는 네이티브 요청 컨텍스트를 애플리케이션 래퍼를 통해 전달합니다. (#5963, #5969)
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6094)에서 재시도를 사용해 푸시가 경합할 때도 변경 로그 항목이 보존되도록 수정했습니다.
- @Grantmartin2002가 [PR](https://github.com/wailsapp/wails/pull/6031)에서 모듈에 포함된 적이 없는 바이너리를 참조하던 embed를 제거하여 모든 플랫폼에서 `pattern arm64/WebView2Loader.dll: no matching files found` 오류로 실패하던 `go mod vendor` 문제를 수정했습니다. 이로써 [#5782](https://github.com/wailsapp/wails/issues/5782) 및 [#5376](https://github.com/wailsapp/wails/issues/5376) 문제도 해결되었습니다.

## 제거

- 순수 Go 로더로 대체된 네이티브 WebView2 로더 지원을 제거했습니다. 이에 따라 내장 `WebView2Loader.dll` 바이너리와 `github.com/jchv/go-winloader` 종속성이 제거됩니다. `native_webview2loader` 빌드 태그는 계속 허용되며 더 이상 오류를 발생시키지 않지만, v3 빌드에는 영향을 주지 않습니다. @Grantmartin2002의 [PR](https://github.com/wailsapp/wails/pull/6031)에서 변경되었습니다.
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6097)에서 macOS API 가이드의 사용하지 않는 빌드 태그와 FPS 옵션을 제거했습니다.

## v3.0.0-beta.19 - 2026-09-09

## 추가

- 선택적으로 사용할 수 있도록 비공개 macOS API를 빌드 태그로 제한했습니다. [문서](https://v3.wails.io/features/browser/integration), [문서](https://v3.wails.io/features/environment/info), [문서](https://v3.wails.io/features/windows/basics), [문서](https://v3.wails.io/features/windows/frameless), [문서](https://v3.wails.io/features/windows/notch-windows), [문서](https://v3.wails.io/features/windows/options), [문서](https://v3.wails.io/guides/build/macos), [문서](https://v3.wails.io/guides/build/private-macos-apis) 및 [문서](https://v3.wails.io/reference/overview)를 참조하세요. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/6087)에서 추가되었습니다.

## 수정

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6091)에서 64 MiB를 초과하는 런타임 요청을 HTTP 413로 거부하도록 변경했습니다.

## 보안

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6092)에서 토큰 인증을 사용해 MCP origin과 원격 액세스를 강화했습니다.

## v3.0.0-beta.18 - 2026-09-08

## 수정

- @4RH1T3CT0R7이 [PR](https://github.com/wailsapp/wails/pull/6083)에서 포인터 리시버를 사용해 Linux 및 Darwin의 Calloc 메모리 누수를 수정했습니다.

## v3.0.0-beta.17 - 2026-09-06

## 수정

- Windows: 이제 WebResourceRequested 핸들러에서 실패했거나 nil인 `GetRequest` 때문에 프로세스가 종료되지 않습니다(`log.Fatal` / nil 역참조 패닉). 대신 요청을 폐기하고 로그에 기록합니다. @midagedev의 [PR](https://github.com/wailsapp/wails/pull/6006)에서 수정되었습니다.

## v3.0.0-beta.16 - 2026-08-29

## 변경

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6029)에서 새 터미널 창에 공증 암호 입력 프롬프트가 표시되도록 변경했습니다.

## 수정

- @ChewbaccaCookie가 [PR](https://github.com/wailsapp/wails/pull/5919)에서 macOS의 시스템 트레이 클릭 유형이 올바르게 처리되도록 수정했습니다.
- @Grantmartin2002가 [PR](https://github.com/wailsapp/wails/pull/6041)에서 CI가 업데이트 전에 사용하지 않는 Microsoft apt 저장소를 제거하도록 변경했습니다.

## v3.0.0-beta.15 - 2026-08-27

## 수정

- @Grantmartin2002가 [PR](https://github.com/wailsapp/wails/pull/6043)에서 WebView2 삽입 제한 시간을 60초로 늘렸습니다.

## v3.0.0-beta.14 - 2026-08-26

## 수정

- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6032)에서 macOS의 Control-문자 키 입력 이름이 올바르게 지정되도록 수정했습니다.
- @nik9play가 [PR](https://github.com/wailsapp/wails/pull/6016)에서 Windows의 ICO 트레이 아이콘을 수정하고 작업 표시줄 테마를 따르도록 변경했습니다.

## v3.0.0-beta.13 - 2026-08-25

## 수정

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6026)에서 macOS에서 모달 루프가 실행되는 동안에도 메인 스레드 작업을 계속 처리하도록 수정했습니다.
- @mortenolsrud가 [PR](https://github.com/wailsapp/wails/pull/5923)에서 모바일 보안 저장소 작업이 실패할 수 있도록 하고, 실패 시 안전하게 차단되도록 수정했습니다.
- @archy-rock3t-cloud가 [PR](https://github.com/wailsapp/wails/pull/5999)에서 등록된 리스너가 없어도 애플리케이션 이벤트 훅이 실행되도록 수정했습니다.
- @haoku123이 [PR](https://github.com/wailsapp/wails/pull/6023)에서 주석과 현지화된 문서의 오타를 수정했습니다.
- @4RH1T3CT0R7이 [PR](https://github.com/wailsapp/wails/pull/6025)에서 `v3/examples` 아래에 커밋된 사전 컴파일 macOS 바이너리를 제거했습니다.

## v3.0.0-beta.12 - 2026-08-21

## 추가

- 수명 주기 및 텔레메트리 예제가 포함된 macOS 노치 알림 창을 추가했습니다. [문서](https://v3.wails.io/features/windows/notch-windows)를 참조하세요. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/6010)에서 추가되었습니다.
- 새 옵션과 네이티브 통합을 제공하는 macOS NSPanel 창 지원 추가 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/6008)에서 구현, [문서](https://v3.wails.io/features/windows/options) 참조

## 수정됨

- 동시 호출 시 SQLite Prepare가 멈추는 문제 방지 — @archy-rock3t-cloud의 [PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 제거됨

- 문서에서 더 이상 사용하지 않는 구현 추적기 제거 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 수정됨

- GTK4 Linux 호스트가 사용자 정의 프로토콜 및 파일 연결을 통한 실행 인수를 누락하는 문제 수정 — @midagedev의 [PR](https://github.com/wailsapp/wails/pull/6000)
- 변경 로그 검증에서 삭제된 줄과 동일한 소스의 수정 사항을 올바르게 처리 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 추가됨

- 에이전트 지원 프로젝트 관리를 위한 안전한 wails3 mcp 서버 추가 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5896)
- 바인딩의 모델 관련 문서 추가 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5988)에서 구현, [문서](https://v3.wails.io/features/bindings/models) 참조
- 원자적 Linux 시스템에서 rpm-ostree 설치 지원 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5987)
- 주간 star-history 차트를 네이티브 방식으로 생성하고 게시하는 기능 추가 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5986)에서 구현, [문서](https://v3.wails.io/credits), [문서](https://v3.wails.io/de/credits), [문서](https://v3.wails.io/fr/credits), [문서](https://v3.wails.io/id/credits), [문서](https://v3.wails.io/ja/credits), [문서](https://v3.wails.io/ko/credits), [문서](https://v3.wails.io/pt/credits), [문서](https://v3.wails.io/ru/credits), [문서](https://v3.wails.io/zh-cn/credits) 및 [문서](https://v3.wails.io/zh-tw/credits) 참조
- 애플리케이션 번들 리소스의 경로를 찾기 위한 Darwin 전용 mac 패키지 추가 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5965)에 있는 [문서](https://v3.wails.io/guides/build/macos) 참조
- Condui 쇼케이스 페이지와 색인 항목 추가 — @mgueregath의 [PR](https://github.com/wailsapp/wails/pull/5962)에서 구현, [문서](https://v3.wails.io/community/showcase/condui) 및 [문서](https://v3.wails.io/community/showcase) 참조
- 스크린샷과 프로젝트 링크가 포함된 Redis Viewer 쇼케이스 페이지 추가 — @redisviewer의 [PR](https://github.com/wailsapp/wails/pull/5984)에서 구현, [문서](https://v3.wails.io/community/showcase) 및 [문서](https://v3.wails.io/community/showcase/redisviewer) 참조

## 변경됨

- Linux용 GTK 애플리케이션 플래그를 G<em>APPLICATION</em>NON_UNIQUE로 업데이트 — @overlordtm의 [PR](https://github.com/wailsapp/wails/pull/5971)
- 누락된 창 이벤트를 경고 수준이 아닌 디버그 수준으로 기록 — @julianstorer의 [PR](https://github.com/wailsapp/wails/pull/5914)

## 수정됨

- 등록된 macOS 키보드 단축키가 webview보다 우선하도록 수정 — @julianstorer의 [PR](https://github.com/wailsapp/wails/pull/5902)
- 깨진 문서 사이드바 링크 복구 — @northes의 [PR](https://github.com/wailsapp/wails/pull/5937)
- WebKit이 일치하는 사용자 정의 스킴 작업을 중단할 때 macOS 및 iOS의 에셋 요청 컨텍스트 취소 (#5963)
- WindowSetFullscreenButtonEnabled 메시지 처리 — @archy-rock3t-cloud의 [PR](https://github.com/wailsapp/wails/pull/5976)
- 빌드 실패를 해결하도록 preact-ts 템플릿에서 Fragment 가져오기 — @haoku123의 [PR](https://github.com/wailsapp/wails/pull/5979)
- 활성 창이나 디스플레이를 사용할 수 있기 전에 화면 검색이 실행될 경우 레거시 GTK3 서비스 전용 애플리케이션이 충돌하는 문제 방지 (#5966)
- 미출시 변경 로그가 비어 있어도 명시적 버전의 릴리스 실행이 계속 진행되도록 수정 (#5977)

## 보안

- 보안 권고를 해결하도록 웹사이트의 nanoid 잠금 파일을 패치된 3.3.18 버전으로 업데이트 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 추가됨

- 자동 변경 로그 항목에 문서 URL 생성 기능 추가 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5957)
- Streams 추가: 수신 소켓 없이 WebSocket 프로그래밍 모델을 사용하여 Go와 JavaScript 간에 양방향 바이트 스트림을 제공합니다. Go에서 `app.HandleStream(name, handler)`을 사용해 스트림을 선언하고 프런트엔드에서 `Stream(name)`을 사용해 연결하면 `WebSocket` 형태의 객체가 반환됩니다. Go→JS는 에셋 서버를 통해 창마다 하나의 대기 상태 폴링으로 전달되고, JS→Go는 일반 POST로 전달됩니다. TCP 포트에는 아무것도 바인딩되지 않으며 `evaluateJavaScript`을 통과하는 것도 없습니다. 서버 빌드(`-tags server`)에서는 동일한 핸들러가 실제 WebSocket을 통해 대신 제공되므로 애플리케이션 코드는 모든 빌드에서 동일합니다. 작성자: @leaanthony
- 메일박스 변경 로그 항목을 Unreleased로 이동 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5935)

## 변경됨

- 문서 사이드바 자동 생성 및 블로그 작성자 타입 파생 방식 업데이트 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5938)

## 수정됨

- WebView2 초기화에 기한과 메시지 펌프 사용 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5952)
- WebView2 쿠키 테스트는 명시적으로 활성화하지 않으면 CI에서 건너뛰며, 실행을 현재 OS 스레드에 고정 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5951)
- Windows 메뉴 빌더에서 하위 메뉴의 상위 항목에 대한 명령 ID 복원 — @gilad-ch의 [PR](https://github.com/wailsapp/wails/pull/5944)
- 공식 크로스 컴파일 이미지가 GTK 4.14+ Linux 지원 기준에 부합하도록 수정 (#5928)
- 상속된 링커 플래그를 유지하고 -ObjC를 추가하도록 iOS Xcode 프로젝트 구성 — @mortenolsrud의 [PR](https://github.com/wailsapp/wails/pull/5915)
- 대규모 프런트엔드에서 `wails3 dev` 에셋 프록시의 과도한 TCP 연결 교체 문제 수정. 이 문제로 인해 호스트의 임시 포트가 고갈되고 관련 없는 프로세스가 `EADDRNOTAVAIL` 오류와 함께 실패할 수 있었음
- 순서가 보장된 디스패치와 역압력을 위해 창별 이벤트 JavaScript를 대기열에 추가 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5934)

## 제거됨

- 데스크톱 바이너리 릴리스 파이프라인 제거: v3 릴리스는 태그만 사용하며 `wails3` CLI는 `go install`을 사용해 설치합니다. `release-v3.yml` 및 이를 디스패치하던 nightly 단계를 삭제 — @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5946)

## v3.0.0-beta.7 - 2026-08-11

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5512)에서 미디어 재생 시 사용자 동작 요구 사항을 비활성화하는 macOS 자동 재생 환경설정을 추가함(@Eyalm321)
- [PR](https://github.com/wailsapp/wails/pull/5935)에서 메일박스 변경 로그 항목을 미출시 섹션으로 이동함(@leaanthony)

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5945)에서 macOS 확대/축소 애니메이션이 CADisplayLink 또는 NSTimer를 사용하도록 하여 성능을 더욱 부드럽게 개선함(@savely-krasovsky)

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5915)에서 상속된 링커 플래그를 유지하고 -ObjC를 추가하도록 iOS Xcode 프로젝트를 구성함(@mortenolsrud)
- 대규모 프런트엔드에서 `wails3 dev` 애셋 프록시가 TCP 연결을 지나치게 자주 교체하여 호스트의 임시 포트를 고갈시키고 관련 없는 프로세스에서 `EADDRNOTAVAIL` 오류를 일으킬 수 있는 문제를 수정함
- [PR](https://github.com/wailsapp/wails/pull/5934)에서 순차적 디스패치와 백프레셔를 위해 창별 이벤트 JavaScript를 대기열에 추가함(@leaanthony)

### 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5851)에서 순차적 이벤트 전달을 위한 제네릭 비동기 FIFO 메일박스를 구현함(@savely-krasovsky, @DevLumuz)

## v3.0.0-beta.6 - 2026-08-09

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5930)에서 크기가 지나치게 큰 이벤트를 위한 용량 제한 호스트 측 저장소와 순차적 JavaScript 전달을 구현함(@leaanthony)
- [PR](https://github.com/wailsapp/wails/pull/5921)에서 창을 깜박일 때 macOS Dock 아이콘이 튀도록 구현함(@julianstorer)

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5931)에서 애셋 서버가 플러시 중 콘텐츠 유형 감지 오류와 아직 기록되지 않은 접두부를 보존하도록 수정함(@leaanthony)
- Wails 콜백에서 애플리케이션 메뉴를 교체할 때 macOS 애플리케이션이 충돌하지 않도록 수정함
- Windows 10 1809 / Windows Server 2019(빌드 17763)에서 네이티브 메뉴를 읽을 수 없던 문제를 수정함. 다크 모드 uxtheme 내보내기가 빌드 18334 이상으로 제한되어 해당 호스트에서는 앱 수준의 다크 모드 옵트인이 실행되지 않았습니다. 이에 따라 메뉴 배경은 어둡게 그려졌지만 Windows는 메뉴 텍스트를 라이트 테마로 계속 그려 어두운 배경에 어두운 텍스트가 표시되었습니다. 해당 서수는 17763부터 존재하므로 이제 조건이 이에 맞게 설정됩니다.
- `w32.GetStockObject`이 `GetStockObject` 대신 `GetDeviceCaps`을 호출하여 모든 스톡 객체에 대해 0을 반환하던 문제를 수정함
- [PR](https://github.com/wailsapp/wails/pull/5924)에서 WebView2 부트스트래퍼 다운로드 오류 처리 및 보고를 개선함(@jannskiee)

## v3.0.0-beta.5 - 2026-08-07

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5897)에서 macOS 앱 활성화가 일반 앱에 대해서만 활성화 정책을 따르도록 수정함(@julianstorer)
- [PR](https://github.com/wailsapp/wails/pull/5898)에서 Linux 빌드의 초기화되지 않은 GTK 창을 보호함(@julianstorer)
- [PR](https://github.com/wailsapp/wails/pull/5899)에서 URL을 불러오기 전에 Linux WebKit 창에 명시적인 불투명 배경색을 설정함(@julianstorer)

## v3.0.0-beta.4 - 2026-08-05

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5890)에서 Android 빌드 작업의 기본 아키텍처를 arm64로 변경하고 deploy-emulator가 호스트 아키텍처를 선택하도록 함(@mortenolsrud)

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5900)에서 드래그하는 동안 macOS 창의 확대/축소 상태를 유지하고 움직임을 줄임(@leaanthony)
- `webview_window_windows_nonclient.go` 빌드 제약 조건에 `!server`을 추가하여 Windows 서버 모드 빌드를 수정함

## v3.0.0-beta.3 - 2026-08-03

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5881)에서 구현 세부 정보에 Phase 10 베타 검증 완료를 문서화함(@leaanthony)

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5877)에서 창 핸들을 Windows 다크 모드 API에 전달하고 인수를 검증함(@leaanthony)
- [PR](https://github.com/wailsapp/wails/pull/5870)에서 프레임 없는 창의 macOS 제목 표시줄 버튼 상태 결정을 중앙화함(@taliesin-ai)
- Windows 앱 테마가 라이트 모드인 상태에서 Windows 애플리케이션이 다크 모드를 요청할 때 네이티브 메뉴 텍스트를 읽을 수 없던 문제를 수정했습니다. 이제 Windows가 다크 모드용 메뉴 텍스트를 렌더링할 수 있을 때까지 메뉴에 그에 맞는 밝은 네이티브 배경을 사용합니다.
- Windows 10 1809 / Windows Server 2019(빌드 17763)에서 네이티브 메뉴를 읽을 수 없던 문제를 수정함. 다크 모드 uxtheme 내보내기가 빌드 18334 이상으로 제한되어 해당 호스트에서는 앱 수준의 다크 모드 옵트인이 실행되지 않았습니다. 이에 따라 메뉴 배경은 어둡게 그려졌지만 Windows는 메뉴 텍스트를 라이트 테마로 계속 그려 어두운 배경에 어두운 텍스트가 표시되었습니다. 해당 서수는 17763부터 존재하므로 이제 조건이 이에 맞게 설정됩니다.

## v3.0.0-beta.2 - 2026-08-02

## 변경됨

- v3를 알파에서 베타로 승격함
- 시스템 트레이의 스마트 기본값과 팝업 자동 숨김 동작을 문서화하고 클릭 핸들러 선택에 대한 회귀 테스트를 추가함(#5840).
- [PR](https://github.com/wailsapp/wails/pull/5861)에서 GitHub 업데이터가 기본적으로 Windows 설치 프로그램 애셋을 제외하도록 변경함(@leaanthony)
- [PR](https://github.com/wailsapp/wails/pull/5866)에서 macOS 프레임 없는 창에 둥근 모서리, 직각 모서리 및 사용자 지정 반경 모서리 지원을 추가함(@leaanthony)

## 수정됨

- 현재 GTK4 창 크기를 보고하고 구성된 서피스에서 크기 조정, 최대화, 최소화 및 전체 화면 상태 이벤트를 발생시키도록 수정함(#5830).
- [PR](https://github.com/wailsapp/wails/pull/5854)에서 fetch 요청으로 Blob 또는 FormData를 보낼 때 Linux WebKit이 충돌하는 문제를 수정함(@taliesin-ai)
- [PR](https://github.com/wailsapp/wails/pull/5865)에서 누락된 Blob/FormData 헤더에 대해 fetch shim이 undefined를 전달하도록 수정함(@leaanthony)

## v3.0.0-alpha2.122 - 2026-08-01

## 추가됨

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5866)에서 macOS 프레임 없는 창에 둥근 모서리, 직각 모서리 및 사용자 지정 반경 모서리 지원을 추가함(@leaanthony)

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5865)에서 누락된 Blob/FormData 헤더에 대해 fetch shim이 undefined를 전달하도록 수정함(@leaanthony)

## v3.0.0-alpha2.121 - 2026-07-31

## 추가됨

- 새 옵션과 빌드 작업을 포함한 macOS DMG 패키징 지원 추가: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5857)

## 변경됨

- GitHub 업데이터가 기본적으로 Windows 설치 프로그램 자산을 제외하도록 변경: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5861)

## 수정됨

- fetch 요청에서 Blob 또는 FormData를 전송할 때 발생하는 Linux WebKit 충돌 수정: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5854)

## v3.0.0-alpha2.120 - 2026-07-31

## 추가됨

- macOS 제목 표시줄을 두 번 클릭할 때 최대화하거나 최소화하는 동작 구현: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5853)

## 변경됨

- QR 서비스 튜토리얼에서 NewServiceWithOptions를 사용하도록 업데이트하고 여백 추가: @jeongkyu의 [PR](https://github.com/wailsapp/wails/pull/5849)

## 수정됨

- macOS에서 확대/축소하는 동안 WKWebView가 계속 응답하도록 수정: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5856)
- GTK4 창 크기 조회를 수정하고 구성된 `GdkSurface`에서 크기 조절, 최대화, 최소화 및 전체 화면 상태 이벤트를 발생시키도록 수정했습니다.

## v3.0.0-alpha2.119 - 2026-07-27

## 수정됨

- 여러 언어의 문서에 아키텍처 다이어그램을 포함하도록 업데이트: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5833)

## v3.0.0-alpha2.118 - 2026-07-26

## 추가됨

- 아이콘 생성 입력 및 출력의 기본 경로 제공: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5825)
- 런타임 package.json의 sideEffects에 소스 진입점 모듈 추가: @savely-krasovsky의 [PR](https://github.com/wailsapp/wails/pull/5797)
- 기여 가이드에 라이선스 및 출처 섹션 추가: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5816)

## 수정됨

- 테두리 반경을 제거하도록 범위가 지정된 GTK4 프레임리스 CSS 적용: @savely-krasovsky의 [PR](https://github.com/wailsapp/wails/pull/5800)
- Windows에서 팝업 메뉴 및 화면 열거 시 커서 위치 확인에 실패해도 정상적으로 처리하도록 수정: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5789)
- macOS 파일 열기 대화 상자가 확장자를 올바르게 필터링하고 접미사를 기준으로 허용된 파일인지 검증하도록 수정: @phergul의 [PR](https://github.com/wailsapp/wails/pull/5678)
- Windows 다크 모드 초기화 시 nil API 호출을 방지하도록 보호 로직 추가: @roachadam의 [PR](https://github.com/wailsapp/wails/pull/5793)
- 업데이터의 32비트 빌드 실패 수정: `GOARCH=386`에서 `fmt.Errorf`에 전달할 때 `maxArchiveTotalSize` 상수(2 GiB)가 플랫폼의 `int` 범위를 초과했습니다. 이제 `int64` 형식으로 명시적으로 지정됩니다.
- 다크 모드 uxtheme API를 로드하지 않는 Windows 빌드(예: Windows 10 1809 / Windows Server 2019(빌드 17763))에서 창이 Dark(또는 시스템 다크) 제목 표시줄을 사용할 때 시작 시 발생하는 nil 포인터 패닉을 수정했습니다. 이제 창 테마 설정의 `AllowDarkModeForWindow` 호출에 nil 보호 로직이 적용되며, 이는 `w32.SetMenuTheme`에서 이미 사용하는 보호 로직과 동일합니다.

## v3.0.0-alpha2.117 - 2026-07-08

## 추가됨

- Windows의 비클라이언트 영역을 위한 사용자 정의 히트 테스트 로직 구현: @savely-krasovsky의 [PR](https://github.com/wailsapp/wails/pull/5462)

## 변경됨

- UseVisualHosting을 기준으로 WebView2 모니터 배율 감지를 구성하도록 변경: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5761)

## v3.0.0-alpha2.116 - 2026-07-07

## 추가됨

- Wails v3 기능과 지침에 중점을 두도록 FAQ 문서 업데이트: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5763)

## v3.0.0-alpha2.115 - 2026-07-06

## 수정됨

- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제 수정(#5659, @puneetdixit200이 #5539에서 별도로 진단하고 수정)
- 디스플레이 변경 시 macOS 화면을 열거하는 과정에서 발생하는 충돌을 화면 ID/이름 문자열을 복사하고 개수를 스냅샷으로 저장하는 방식으로 수정(#5565, @x-haose가 #5584에서 별도로 진단하고 수정)
- 최소화/복원 전환 중 `GetClientRect`이 nil을 반환하는 상황에서 `WM_ERASEBKGND`이 단색 배경을 그릴 때 Windows에서 발생하는 충돌 수정(@sinspired가 #5636에서 보호 로직을 보고)
- 프런트엔드 바인딩 오류가 항상 텍스트로 파싱되는 문제 수정: #5690의 @mbaklor
- macOS 및 Linux의 해당 파일에는 이미 있는 `!server` 빌드 제약 조건이 Windows GUI 파일에 누락되어 `server` 빌드 태그 사용 시 발생하는 Windows 빌드 실패 수정(#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## 추가됨

- Update Manifest 프로토콜 및 엔드포인트 공급자 구현: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5720)

## 변경됨

- `webview2` 바인딩을 `v3/internal/webview2`로 v3 모듈에 통합하고, 독립 실행형 모듈과 해당 모듈의 나이틀리 릴리스/동기화 워크플로 및 go.mod 버전 조정 절차를 제거했습니다(v3만 이 모듈을 사용함): @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5711)

## 수정됨

- WebView2 모니터 배율 감지 및 DPI 변경 시 호스트 재동기화 수정 사항을 미출시 섹션으로 이동: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5750)
- float64 및 BOOL 매개변수에 대한 WebView2 COM 마샬링 업데이트: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5741)
- Windows 시스템 트레이 아이콘 업데이트 및 제거 시 패닉과 nil 역참조 방지: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5703)
- Windows에서 숨겨진 창이 올바르게 다시 숨겨지지 않는 문제 수정: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5743)
- 창 최소화/최대화/복원 시 WebView2 컨트롤러의 표시 상태를 동기화: @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5742)

### 수정됨

- WebView2 모니터 배율 감지를 다시 활성화하고 DPI 변경 시에만 호스트가 재동기화되도록 제한: @randalmurphal이 검증한 수정 사항을 바탕으로 @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5734), @eleclin이 근본 원인을 검증하고 @qq540491950이 하드웨어 테스트 수행

## v3.0.0-alpha2.113 - 2026-07-04

## 추가됨

- `ANDROID_KEYSTORE_FILE`이 설정되지 않은 상태에서 릴리스 AAB를 빌드하면 경고하도록 추가하고(Google Play는 디버그 서명된 번들을 거부함), @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5730)에서 App Bundle 패키징 및 서명 방법을 문서화함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5739)에서 여러 언어로 된 Why Wails 문서를 추가함
- @fbbdev가 작성한 [PR](https://github.com/wailsapp/wails/pull/5398)에서 바인딩의 Go time.Time을 JS Date 또는 string으로 매핑하도록 지원함
- @mortenolsrud가 작성한 [PR](https://github.com/wailsapp/wails/pull/5728)에서 Play Store 제출용 Android App Bundle(AAB) 패키징 작업(`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`)을 추가함. APK 작업은 로컬/에뮬레이터 테스트용으로 유지함([#5726](https://github.com/wailsapp/wails/issues/5726) 해결)
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5735)에서 Android 실제 기기용 작업 대상을 추가하고 카메라/위치 권한을 재개함

## 변경됨

- `webview2`을 v1.0.28로 업데이트함([릴리스 정보](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5730)에서 Android 템플릿의 `compileSdk`/`targetSdk`을 34에서 신규 앱 제출 시 Google Play가 요구하는 35으로 업데이트함

## 수정됨

- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/5745)에서 sponsorkit에 내장된 아바타 마스크를 수정함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5730)에서 사전식 버전 정렬로 인해 Android AVD 자동 생성 시 잘못된 시스템 이미지 또는 cmdline-tools 버전이 선택되는 문제를 수정함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5730)에서 설정 마법사가 오래된 Android NDK 버전을 제안하는 문제를 수정함(이제 문서화된 요구 사항과 일치하는 26.3.11579264을 제안함)
- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/5744)에서 SvelteKit 및 옵션에 관한 프랑스어 문서를 업데이트함
- @flofreud가 작성한 [PR](https://github.com/wailsapp/wails/pull/5516)에서 디스플레이 변경 중 macOS 화면 열거 시 발생하는 SIGSEGV를 수정함

## v3.0.0-alpha2.112 - 2026-07-03

## 추가됨

- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5724)에서 Go 기반 기여자 SVG 생성기를 추가하고 문서/웹사이트의 크레딧 페이지를 업데이트함
- @fbbdev가 작성한 [PR](https://github.com/wailsapp/wails/pull/5398)에서 바인딩의 Go time.Time을 JS Date 또는 string으로 매핑하도록 지원함

## 변경됨

- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5719)에서 Node 기반 스폰서 이미지 파이프라인을 Go 생성기로 교체함

## 수정됨

- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5729)에서 Android 빌드 애셋 종속성 설치 스크립트를 수정함
- `ValidateAndSanitizeURL`에서 U+0085(NEXT LINE) 제어 문자를 거부하도록 하여 URL 유효성 검사기의 공백 문자 검사를 완성함
- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/4785)에서 프레임 없는 창의 DPI가 변경될 때 DWM 프레임을 다시 계산하도록 수정함
- @yulesxoxo가 작성한 [PR](https://github.com/wailsapp/wails/pull/4632)에서 Windows 배율이 100%가 아닐 때 DnD 드롭 영역 감지가 실패하는 문제를 수정함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5714)에서 Darwin의 대화 상자, 메뉴, 트레이 및 알림 전반에 사용되는 Cocoa 객체에 명시적인 Objective-C 메모리 관리를 추가함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5718)에서 Linux CGO 백엔드 버그와 시스템 트레이 문제를 수정함

## v3.0.0-alpha2.111 - 2026-07-01

## 추가됨

- @Aliuyanfeng이 작성한 [PR](https://github.com/wailsapp/wails/pull/5061)에서 HappyTools를 커뮤니티 쇼케이스에 추가함
- @triadmoko가 작성한 [PR](https://github.com/wailsapp/wails/pull/5643)에서 인도네시아어 로케일 지원과 포괄적인 문서를 추가함
- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/4813)에서 WindowsWindow에 DisableMenu 옵션을 추가함

## 변경됨

- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/5617)에서 GOOS와 ARCH를 사용해 build/package 작업을 디스패치하도록 Taskfile 템플릿과 CLI를 업데이트함

## 수정됨

- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5708)에서 mac 창 탭 기능 문제를 수정함

## 제거됨

- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5702)에서 기여, 기능 및 가이드의 독일어 번역 MDX 파일을 제거함

## v3.0.0-alpha2.110 - 2026-06-30

## 추가됨

- @wayneforrest가 작성한 [PR](https://github.com/wailsapp/wails/pull/5129)에서 macOS WebView 새로고침과 강제 새로고침을 구현하고 WebContent 프로세스 종료 시 복구 기능을 추가함
- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/5396)에서 기여, 기능 및 가이드에 관한 포괄적인 독일어 문서를 추가함
- @popaprozac이 작성한 [PR](https://github.com/wailsapp/wails/pull/5333)에서 소리, 첨부 파일, 예약 및 업데이트 API로 알림 기능을 강화함

## 수정됨

- @leaanthony가 작성한 [PR](https://github.com/wailsapp/wails/pull/4785)에서 프레임 없는 창의 DPI가 변경될 때 DWM 프레임을 다시 계산하도록 수정함
- @yulesxoxo가 작성한 [PR](https://github.com/wailsapp/wails/pull/4632)에서 Windows 배율이 100%가 아닐 때 DnD 드롭 영역 감지가 실패하는 문제를 수정함

## v3.0.0-alpha2.109 - 2026-06-29

## 추가됨

- @iamhabbeboy가 작성한 [PR](https://github.com/wailsapp/wails/pull/5026)에서 EventsEmit 문서에 코드 예제를 추가함
- @MerIijn이 작성한 [PR](https://github.com/wailsapp/wails/pull/5380)에서 Windows WebView2 비주얼 호스팅 옵션을 추가함
- @SametKUM이 작성한 [PR](https://github.com/wailsapp/wails/pull/5536)에서 Klustr를 커뮤니티 쇼케이스 문서에 추가함
- @thiennguyen93이 작성한 [PR](https://github.com/wailsapp/wails/pull/5685)에서 새 페이지 및 변경 로그 항목과 함께 Kira를 커뮤니티 쇼케이스에 추가함
- @taliesin-ai가 작성한 [PR](https://github.com/wailsapp/wails/pull/5694)에서 MCP 서비스 가이드에 피드백 섹션을 추가함

## 변경됨

- 서버 모드에서 이제 데스크톱 빌드 작업과 일관된 일급 프로덕션 빌드를 지원합니다(#5693). `task build:server`은 기본적으로 프로덕션 바이너리(`-tags server,production`, `-trimpath`, 심볼 제거)를 빌드하며, `DEV=true`(개발 서버), `OBFUSCATED=true`(garble), `EXTRA_TAGS`을 허용합니다. `task run:server`은 개발 서버를 실행합니다. `Dockerfile.server` / `task build:docker`는 먼저 프로덕션 서버(`-tags server,production`)와 프로덕션 프런트엔드를 빌드합니다. 이미지는 기본적으로 distroless/static 기반의 순수 Go 정적 빌드를 사용하며, CGO 앱을 위해 `CGO_ENABLED`, `GO_IMAGE`, `RUNTIME_IMAGE`을 재정의 가능한 빌드 인수로 제공합니다.

## 수정됨

- 대기 중인 비동기 호출이 있는 창을 닫을 때 발생하는 충돌 방지: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/4435)
- Windows에서 숨겨진 앱을 열 때 창이 활성화되지 않도록 수정: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5249)
- WebKit 요청 메타데이터, 응답 완료 및 본문 스트림 처리가 GTK 메인 스레드에서 실행되도록 보장: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5668)
- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제 수정(#5659, @puneetdixit200도 #5539에서 독립적으로 진단하고 수정)
- 화면 ID/이름 문자열을 복사하고 개수의 스냅샷을 생성하여 디스플레이 변경 시 macOS 화면을 열거할 때 발생하는 충돌 수정(#5565, @x-haose도 #5584에서 독립적으로 진단하고 수정)
- Windows에서 혼합 DPI 모니터 사이로 창을 드래그한 후 WebView2 콘텐츠가 축소되었다가 사라지는 문제를 수정했습니다. 최소화 해제 시의 DPI 재동기화와 동일하게 `WM_DPICHANGED` 핸들러에서 컨트롤러 경계를 다시 설정합니다(#5677).

## v3.0.0-alpha2.108 - 2026-06-28

## 추가됨

- `app.GlobalShortcut`을 통해 전역(시스템 전체) 키보드 단축키를 추가했습니다(`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). 애플리케이션에 포커스가 없어도 단축키가 작동합니다. 타사 종속성 없이 플랫폼별로 네이티브 구현을 사용합니다. macOS에서는 Carbon 핫키, Windows에서는 `RegisterHotKey`, X11에서는 `XGrabKey`, Wayland에서는 XDG Desktop Portal 전역 단축키 인터페이스를 사용합니다.
- 내장 MCP 서버를 추가했습니다. 애플리케이션을 `mcp` 태그로 빌드하면 자동으로 시작되는 Model Context Protocol 서버로, LLM 에이전트가 실행 중인 Wails 애플리케이션을 테스트하고 제어할 수 있습니다. 창 제어, DOM 검사, JavaScript 평가, 바인딩된 메서드 호출, 이벤트, 애니메이션 화면 커서로 표시되는 마우스/키보드 입력 시뮬레이션을 지원합니다. 사용자 코드는 필요하지 않습니다. `WAILS_MCP=1`이 설정되면 `wails3 build`/`wails3 dev`이 `mcp` 태그를 자동으로 추가합니다. 환경 변수(`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`)만으로 구성합니다.

## 수정됨

- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제 수정(#5659, @puneetdixit200도 #5539에서 독립적으로 진단하고 수정)
- 화면 ID/이름 문자열을 복사하고 개수의 스냅샷을 생성하여 디스플레이 변경 시 macOS 화면을 열거할 때 발생하는 충돌 수정(#5565, @x-haose도 #5584에서 독립적으로 진단하고 수정)

## v3.0.0-alpha2.107 - 2026-06-27

## 추가됨

- 사이드바 탐색 기능이 포함된 실험적 Wake 문서 추가: @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5613)

## v3.0.0-alpha2.106 - 2026-06-24

## 변경됨

- `webview2`을 v1.0.27로 업데이트했습니다.
  - ci(webview2): 릴리스 빌드 수정(Windows 교차 컴파일 + 완전한 go.sum)(#5671)\

  **전체 차이:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- webview2 릴리스 워크플로의 교차 컴파일에서 go vet 제거: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5672)
- 자동 변경 로그의 OpenRouter 모델을 google/gemini-2.5-flash-lite로 업데이트: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5670)
- `webview2`을 v1.0.26로 업데이트했습니다.

### 수정 사항

- **일시적인 런타임 COM 오류가 발생해도 종료하지 않고 복구하도록 수정**(#5658, #5580). 이전에는 `Chromium.errorCallback`이 *모든* COM 오류에 대해 `os.Exit(1)`을 호출했으므로, 시작 후 복구 가능한 일시적 오류가 발생해도 애플리케이션 전체가 종료되었습니다. 이제 런타임 경로(`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`)는 오류를 기록하고 복구합니다. 특히 `MessageReceived`의 잘못 구성되었거나 신뢰할 수 없는 웹 메시지는 이제 프로세스를 종료시키는 대신 폐기됩니다. 이 변경으로 혼합 DPI 모니터 사이를 이동할 때 발생하는 충돌 유형을 해결합니다(#5544, #5650). 환경/컨트롤러 생성 경로의 오류는 계속 치명적 오류로 처리됩니다.\

**전체 차이:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 수정됨

- go.sum 파일을 올바르게 처리하도록 release-webview2 워크플로 수정: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5671)
- 네이티브 메뉴를 지우고 다시 빌드하여 Linux GTK4 메뉴 업데이트 수정: @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5659)

## v3.0.0-alpha2.105 - 2026-06-21

## 추가됨

- 공유 코드에서 런타임 플랫폼을 감지할 수 있도록 `application.System`을 추가했습니다. `System.IsMobile()`(iOS/Android), `System.IsDesktop()`(macOS/Windows/Linux), `System.IsServer()`(`server` 빌드 태그), 단일 대상을 직접 테스트하는 `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`을 제공합니다. 모든 대상에서 컴파일되므로 빌드 태그 없이 분기할 수 있습니다. 이에 대응하는 프런트엔드 헬퍼(`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`)는 `@wailsio/runtime`에서 사용할 수 있습니다.
- 자체 Vite 프로젝트를 `frontend/`에 추가하는 방법을 보여 주는 "다른 프런트엔드 프레임워크 사용" 가이드를 추가했습니다(Solid, Preact, Lit, SvelteKit, Qwik, Angular 등 포함).
- 이제 `wails3 setup` 마법사가 모바일(iOS/Android) 툴체인인 Xcode와 iOS Simulator 런타임, JDK, Android SDK/NDK 및 에뮬레이터를 검사하며, 해당되는 경우 원클릭 설치와 복사 가능한 셸 구성 수정 사항을 제공합니다.
- 생성된 프로젝트에는 새로 게시된(잠재적으로 침해된) 패키지에 대한 노출을 줄이기 위해 7일의 `minimum-release-age`를 설정하는 `frontend/.npmrc`이 포함됩니다(pnpm과 bun은 이를 준수하며 npm은 문제없이 무시합니다).

## 변경됨

- 모든 내장 스타터 템플릿을 새로운 네온 산 배경의 히어로 디자인으로 개편했습니다(web, iOS 및 Android).
- **이제 TypeScript가 스타터 템플릿의 기본 언어이며 접미사가 없는 템플릿 이름을 사용합니다.** `wails3 init`(`-t` 없음)은 TypeScript 프로젝트를 스캐폴딩합니다. `-t vanilla`, `-t react`, `-t vue`, `-t svelte`은 TypeScript 템플릿이며, JavaScript 변형은 `-t vanilla-js`, `-t react-js`, `-t vue-js`, `-t svelte-js`입니다. 내장 템플릿은 `template.yaml`의 `typescript:`으로 언어를 선언합니다. `-ts` 접미사를 사용하는 커뮤니티 템플릿은 대체 방식으로 계속 작동합니다.
- 산 배경 위에 반투명 유리 효과의 생동감을 더한 네온 "digital Wails" 테마로 `wails3 setup` 마법사를 개편했습니다.

## 수정

- WebView2가 일시 중단되거나 렌더링/GPU 프로세스가 재활용될 만큼 오랫동안 최소화된 앱을 Windows에서 복원할 때 발생하는 충돌을 수정했습니다. 이제 최소화/복원 시 DPI 재동기화(#5544)는 창의 DPI가 실제로 변경된 경우에만 WebView2 컨트롤러에 접근하므로, 일반적인 동일 DPI 복원 과정에서 일시 중단된 컨트롤러를 대상으로 치명적인 COM 호출이 발생하지 않습니다(#5605).
- 자산/미디어를 빈번하게 로드하며 장시간 실행되는 Linux 앱에서 네이티브 `SIGABRT`/`SIGSEGV` 충돌이 반복되는 문제를 수정했습니다. 이 충돌은 일반적으로 GTK 메인 루프 중 `g_object_unref` 내부에서 발생했습니다. 자산 서버가 작업자 고루틴에서 `WebKitURISchemeRequest`를 완료하면서 스레드 안전성이 보장되지 않는 WebKit2GTK 함수를 GTK 메인 스레드 외부에서 호출했습니다. 이제 완료 처리(`webkit_uri_scheme_request_finish_with_response`/`finish_error`)는 메인 스레드에서 실행됩니다. #5566의 부분적인 수정 사항을 완성합니다. GTK3 및 GTK4/WebKitGTK 6.0 빌드 모두에 영향을 줍니다(#5631, #5557).
- Linux/GTK3의 `setupSignalHandlers`에서 간헐적으로 발생하는 `fatal error: invalid pointer found on stack`을 수정했습니다. 신호 `user_data`로 전달된 창 ID가 Go의 로컬 `unsafe.Pointer`에 보관되어, 스택 복사 중 가비지 컬렉터가 이 비포인터 값을 스캔할 때 실행을 중단했습니다. 이제 Go 측에서는 ID를 정수형(`uintptr_t`)으로 유지합니다. 이는 #4958에서 GTK4 경로에 적용한 것과 동일한 수정 사항을 기존 GTK3 경로에 백포트한 것입니다. GTK4 경로에서는 `-race`/checkptr 오류를 없애기 위해 C 신호 함수를 `uintptr_t`로 변경했습니다(#5631).

## 제거

- `react-swc`, `preact`, `lit`, `solid`, `qwik`, `sveltekit` 시작 템플릿과 해당 `-ts` 변형을 제거했습니다. 이제 지원되는 기본 제공 템플릿은 `vanilla`, `react`, `vue`, `svelte`입니다. 각 템플릿은 기본적으로 TypeScript를 사용하며, `-js` JavaScript 변형도 제공합니다. 그 밖의 프레임워크도 계속해서 [자체 프런트엔드를 가져오는 방식](https://v3.wails.io/guides/dev/frontend-frameworks)이나 사용자 지정 템플릿을 통해 사용할 수 있습니다.

## v3.0.0-alpha2.104 - 2026-06-18

## 수정

- 바인딩된 Go 서비스 메서드가 빈 문자열을 반환할 때 발생하는 iOS 충돌(SIGABRT)을 수정했습니다. iOS 자산 응답 작성기가 본문 길이 대신 `buf != nil`으로 본문 포인터를 검사하여, 길이가 0인 본문에서 `&buf[0]` 패닉이 발생했습니다. 이제 데스크톱 작성기와 동일하게 길이를 검사합니다.

## v3.0.0-alpha2.103 - 2026-06-15

## 변경

- iOS 및 Android 네이티브 기능을 플랫폼 관리자로 이동했습니다. 이제 기존 `application.IOS*`/`application.Android*` 독립 함수 대신 `application.IOS.*` 및 `application.Android.*`을 통해 호출합니다(예: `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`)(#5602).
- 모바일 브리지 이벤트의 이름을 변경했습니다. 이제 크로스 플랫폼 이벤트는 `common:*` 접두사를 사용하고(예: `common:haptic`, `common:location`), 플랫폼 전용 이벤트는 `ios:*` / `android:*`를 사용합니다(예: `ios:backgroundTask`, `android:foregroundService`). `native:*` 접두사는 더 이상 사용하지 않습니다(#5602).

## v3.0.0-alpha.102 - 2026-06-14

## 추가

- 대화형 프로젝트 설정 및 종속성 검사를 위한 실험적 `wails3 setup` 마법사를 추가했습니다.
- 기계 판독 가능 출력을 위해 `wails3 doctor`에 `--json` 플래그를 추가했습니다.
- `wails3 doctor` 명령에 서명 상태 섹션을 추가했습니다.

## 수정

- Linux에서 npm을 감지할 때 패키지 관리자뿐 아니라 PATH도 확인하도록 수정했습니다.

## v3.0.0-alpha.101 - 2026-06-13

## 추가

- iOS: 네이티브 메시지 대화 상자(UIAlertController)와 파일/여러 파일/디렉터리 열기 대화 상자(UIDocumentPickerViewController)를 추가했습니다. 저장 대화 상자는 명시적인 오류를 반환합니다.
- iOS: UIPasteboard를 통한 클립보드 지원을 추가했습니다.
- iOS: UIScreen을 통한 실제 화면 측정값(포인트, 픽셀, 배율, 안전 영역 내 작업 영역)을 추가했습니다.
- iOS: 기기 빌드(`IOS_PLATFORM=device`), 코드 서명 ID/프로비저닝 프로파일/entitlements 지원, `.ipa` 패키징 및 devicectl을 통한 `deploy-device`를 추가했습니다.
- iOS: 설정 가능한 최소 iOS 버전(build/config.yml의 `ios.minIOSVersion`)을 추가했습니다.
- iOS: macOS에서 `wails3 doctor`이 Xcode 및 iOS SDK의 사용 가능 여부를 보고합니다.
- iOS: 배터리, 네트워크, 테마, 화면 잠금 및 메모리 부족 시스템 이벤트를 `events.IOS.*` 및 플랫폼 중립적인 `events.Common.*` 애플리케이션 이벤트로 제공합니다.
- iOS: 네이티브 모바일 기능 브리지(내보낸 `application.IOS*`)를 추가했습니다. 공유 시트, URL 열기, 절전 방지, 손전등, 안전 영역 인셋, 밝기, 앱 정보, 방향 잠금, 상태 표시줄, 생체 인증(Face ID/Touch ID), 로컬 알림 및 Keychain 보안 저장소를 지원합니다.
- iOS: 센서 및 하드웨어 기능으로 햅틱, 일회성 위치 확인, 가속도계, 근접 센서, 텍스트 음성 변환, 저장소 정보, 전원/배터리 상태, 네트워크 상태, 키보드 인셋 및 화면 캡처 감지를 추가했습니다.
- iOS: 문서(IOS.md 및 문서 사이트 가이드)를 추가했습니다.
- Android: 네이티브 메시지 대화 상자(AlertDialog)와 파일/여러 파일 열기 대화 상자(Storage Access Framework, 캐시 복사본으로 가져옴)를 추가했습니다. 디렉터리 열기 및 저장 대화 상자는 명시적인 오류를 반환합니다.
- Android: ClipboardManager를 통한 클립보드 지원을 추가했습니다.
- Android: WindowMetrics/DisplayMetrics를 통한 실제 화면 측정값(dp, 픽셀, 배율, 시스템 표시줄을 제외한 작업 영역)을 추가했습니다.
- Android: 햅틱(`Android.Haptics.Vibrate`), 기기 정보(`Android.Device.Info`) 및 토스트(`Android.Toast.Show`) 런타임 메서드를 추가했습니다.
- Android: 형식이 지정된 수명 주기 이벤트(`events.Android.*`, events.txt에서 생성)를 추가하고, `ActivityCreated`을 `Common.ApplicationStarted`에 매핑했습니다.
- Android: 빌드 파이프라인에서 설치 가능한 디버그 및 릴리스 APK(`android:run`, `android:package`, `android:package:fat`)를 생성합니다. 릴리스 서명에는 기본적으로 디버그 키 저장소를 사용하며, `ANDROID_KEYSTORE_*` 환경 변수를 통해 실제 키 저장소를 사용할 수도 있습니다.
- Android: `wails3 doctor`이 Android SDK, NDK 및 JDK를 보고합니다.
- Android: 배터리, 네트워크, 테마, 화면 잠금 및 메모리 부족 시스템 이벤트를 `events.Android.*` 및 플랫폼 중립적인 `events.Common.*` 애플리케이션 이벤트로 제공합니다.
- Android: 네이티브 모바일 기능 브리지(내보낸 `application.Android*`)를 추가했습니다. 공유, URL 열기, 절전 방지, 손전등, 안전 영역 인셋, 밝기, 앱 정보, 방향 잠금, 상태 표시줄, 생체 인증(BiometricPrompt), 로컬 알림 및 EncryptedSharedPreferences 보안 저장소를 지원합니다.
- Android: 센서 및 하드웨어 기능으로 햅틱, 일회성 위치 확인, 가속도계, 근접 센서, 텍스트 음성 변환, 저장소 정보, 전원/배터리 상태, 네트워크 상태, 키보드 인셋 및 FLAG_SECURE 화면 캡처 차단을 추가했습니다.
- Android: 문서(ANDROID.md 및 문서 사이트 가이드)를 추가했습니다.
- 예제: `mobile` 종합 예제에 iOS와 Android 전반의 네이티브 기능 브리지를 보여 주는 Mobile 및 Hardware 탭을 추가했습니다. 알약 모양 탭은 여러 행으로 줄바꿈됩니다.
- 모바일: 배터리 절약을 위해 앱이 백그라운드로 전환되면 가속도계, 근접 센서, 손전등 및 예제의 주기적 시계를 일시 중지하고, 앱으로 돌아오면 복원합니다(Android는 백그라운드에서도 프로세스를 계속 실행하며, iOS에서 손전등은 지속되는 하드웨어 상태입니다). 또한 Android 시스템 이벤트 수신기는 앱이 포그라운드에 있는 동안에만 등록됩니다.
- iOS: 카메라 캡처 — `application.IOSCapturePhoto`/`IOSCaptureVideo`(UIImagePickerController → base64 썸네일이 포함된 `native:capture` 이벤트)
- iOS: 백그라운드 실행 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`(UIApplication 백그라운드 작업 실행 시간) 및 생성된 Info.plist에 `UIBackgroundModes`을 템플릿으로 삽입하는 구성 가능한 `ios.backgroundModes`(build/config.yml)
- Android: 카메라 캡처 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`(FileProvider를 통한 시스템 카메라 → `native:capture` 이벤트)
- Android: 포그라운드 서비스 — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`(지속 알림을 사용하는 `WailsForegroundService`가 장시간 실행되는 백그라운드 작업을 위해 프로세스를 계속 실행함)
- 예제: 사진/동영상 캡처와 백그라운드 실행을 보여 주는 Camera 탭(Android에서는 포그라운드 서비스, iOS에서는 백그라운드 작업 실행 시간 사용)

## 수정됨

- Linux에서 `getUserMedia`이 항상 `NotAllowedError` 오류로 실패하던 문제를 수정했습니다. WebKitGTK는 아무도 처리하지 않는 권한 요청을 거부하며, `permission-request` 신호가 연결되어 있지 않았습니다. 이제 새로운 크로스 플랫폼 `WebviewWindowOptions.Permissions` 맵(`map[PermissionType]Permission`)에 따라 카메라와 마이크를 처리하며, Linux(WebKitGTK)와 Windows(WebView2) 모두에서 이 설정을 적용합니다. 네이티브 권한 프롬프트가 없는 Linux에서는 카메라와 마이크를 기본적으로 허용하여 `getUserMedia`을 복원하며, `PermissionDeny`으로 비활성화할 수 있습니다(#5552).
- iOS: `GOOS=ios`이 다시 컴파일됩니다(`events.IOS` 내보내기, 모바일 메서드 이름 스텁 추가). 또한 pkg/application과 여러 서비스의 빌드 태그를 수정하여 프로덕션 태그가 지정된 빌드도 컴파일됩니다.
- iOS: 이제 Go→JS 이벤트와 ExecJS가 작동합니다. 시작 시 페이지가 더 이상 두 번 로드되지 않으며 `wails:runtime:ready` 핸드셰이크도 더 이상 유실되지 않습니다.
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted`이 더 이상 앱 시작 과정과 경합하지 않으며, 고정된 2초의 시작 대기 시간을 제거했습니다.
- iOS: Go→JS JavaScript를 실행할 때마다 발생하던 C 문자열 누수를 수정했습니다.
- iOS: 이제 `hasListeners`이 실제 리스너 등록 상태를 반영합니다.
- iOS: 프로덕션 빌드에서는 프레임워크 디버그 로깅이 컴파일에서 제외됩니다.
- Android: `GOOS=android`이 다시 컴파일됩니다. `events.Android`을 정의하고, 범위를 벗어나던 `events_android.go` 리스너 배열을 제거했으며, 모바일 메서드 이름 스텁을 추가하고, 데스크톱 Linux 파일(`linux_cgo.*`, `events_linux.*`, `environment_linux.go`)이 Android 빌드에 포함되지 않도록 했습니다.
- Android: 이제 JS→Go 바인딩이 작동합니다. WebView가 `fetch()` POST 본문을 `shouldInterceptRequest`에 전달할 수 없으므로, 런타임 호출은 nil 요청 본문으로 인해 충돌하는 대신 JavascriptInterface 전송 계층(`nativeHandleRuntimeCall`)을 통해 라우팅됩니다.
- Android: 이제 `Screens.*` 런타임 호출이 실제 데이터를 반환합니다. 시작 시 ScreenManager를 채우도록 변경했습니다(이전에는 연결되지 않아 `GetAll`이 nil을 반환했습니다).
- Android: 프로덕션 빌드에서는 프레임워크 디버그 로깅이 컴파일에서 제외되며, 디버그 빌드에서는 `Wails` 태그로 logcat에 기록됩니다.
- Android: 실제 `hasListeners` 레지스트리, JNI 참조/예외 처리 및 단일 로드 페이지 수명 주기(중복 탐색 없음)를 구현했습니다.
- Windows에서 Vite 개발 서버가 실행 중일 때 `wails3 generate bindings`이 "액세스가 거부되었습니다" 오류로 실패하던 문제를 수정했습니다. 이름 변경 작업으로 기존 출력 디렉터리를 덮어쓰는 대신 생성된 파일을 출력 디렉터리에 동기화합니다(#5515).
- macOS에서 디스플레이 변경 후 화면 정보를 읽을 때 간헐적으로 발생하던 치명적인 충돌을 수정했습니다. 화면 ID와 이름이 Go에서 복사하기 전에 해제될 수 있는 자동 해제 `UTF8String` 버퍼의 포인터를 저장하여 해제 후 사용 문제가 발생했습니다. 이제 문자열에 `strdup`을 적용하고 변환 후 해제하며, 화면 열거를 명시적인 자동 해제 풀에서 실행하므로 Go 고루틴에서 호출해도 더 이상 누수가 발생하지 않습니다(#5556).
- Linux에서 assetserver가 `WebKitURISchemeRequest`을 닫을 때 간헐적으로 발생하던 SIGSEGV를 수정했습니다. 마지막 `g_object_unref`이 assetserver 고루틴에서 실행되어 GTK 메인 스레드 외부에서 WebKit GObject를 종료했습니다. 이제 `g_main_context_invoke`을 통해 unref를 GTK 메인 컨텍스트로 마샬링합니다(#5557).

## v3.0.0-alpha.100 - 2026-06-13

## 추가됨

- `MacWebviewPreferences`에 WKWebView 구성 옵션 `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` 및 `ApplicationNameForUserAgent`을 추가했습니다(#5549).

## 수정됨

- Windows에서 Vite 개발 서버가 실행 중일 때 `wails3 generate bindings`이 "액세스가 거부되었습니다" 오류로 실패하던 문제를 수정했습니다. 이름 변경 작업으로 기존 출력 디렉터리를 덮어쓰는 대신 생성된 파일을 출력 디렉터리에 동기화합니다(#5561).
- Linux의 프레임 없는 창에서 JS 크기 조정 이벤트가 발생하지 않던 문제와 프레임 없는 창의 스크롤바 가장자리 감지 문제를 수정했습니다(#5368).
- Windows에서 임시 디렉터리와 설치 디렉터리가 서로 다른 볼륨에 있을 때 업데이터가 "invalid cross-device link" 오류로 실패하던 문제를 수정했습니다(#5560).

## v3.0.0-alpha.99 - 2026-06-10

## 수정됨

- Windows에서 Vite 개발 서버가 실행 중일 때 `wails3 generate bindings`이 "액세스가 거부되었습니다" 오류로 실패하던 문제를 수정했습니다. 이름 변경 작업으로 기존 출력 디렉터리를 덮어쓰는 대신 생성된 파일을 출력 디렉터리에 동기화합니다(#5515).

## v3.0.0-alpha.98 - 2026-06-03

## 수정됨

- Linux에서 유휴 상태일 때(예: 인스펙터가 열려 있을 때) WebKit UI가 멈추던 문제를 수정했습니다. JavaScriptCore의 GC 스레드 동기화를 손상시키던 `SIGUSR1`의 `SA_ONSTACK` 강제 설정을 제거했습니다(#5527).

## v3.0.0-alpha.97 - 2026-05-31

## 추가됨

- 디버깅 페이지와 `runtime/trace` 사용 방법을 추가했습니다.

## 변경됨

- 불필요한 `_ "embed"` import 일부를 제거하여 코드를 정리했습니다.

## 수정됨

- Windows에서 창의 최대화를 해제한 후 최소 너비/높이 제약 조건이 적용되지 않던 문제를 수정했습니다(#4593).
- Frameless + Transparent 창 옵션을 사용할 때 전체 화면 모드에서 마우스 클릭이 창을 통과하던 문제를 수정했습니다(#4408).

## v3.0.0-alpha.96 - 2026-05-25

## 추가됨

- Garble 난독화 지원([#4563](https://github.com/wailsapp/wails/issues/4563))을 추가했습니다. 여기에는 안정적인 바인딩 메서드 ID, 빌드/Taskfile 연동(`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), 그리고 Garble이 내보낸 필드 이름을 변경해도 전송 형식이 유지되도록 모든 런타임 대상 페이로드에 추가한 JSON 구조체 태그(`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`)가 포함됩니다.

## v3.0.0-alpha.95 - 2026-05-20

## 추가됨

- 누락된 프로젝트 구조 페이지를 추가했습니다.

## 변경됨

- 문서: 더 깔끔하게 표시되도록 아키텍처 페이지의 다이어그램 두 개를 시퀀스 다이어그램으로 변경했습니다.
- 문서: 실행 전제 조건으로 D2 설치에 관한 참고 사항 추가

## 수정됨

- GTK4가 기본값일 때 발생하는 `wails3 generate appimage` 문제 수정: 이제 번들러가 런타임 파일을 검색하기 전에 바이너리에서 GTK 스택을 감지하므로, GTK4 빌드에는 `webkitgtk-6.0/` 아래의 `libwebkitgtkinjectedbundle.so`을 선택하고 `-tags gtk3` 빌드에는 `webkit2gtk-4.1/` 아래의 `libwebkit2gtkinjectedbundle.so`을 선택합니다. 또한 `.relr.dyn` 검사에서 `libgtk-4.so.1`도 확인하므로 스택과 관계없이 최신 툴체인에서 스트리핑이 올바르게 비활성화됩니다. (#5475)
- 상대 경로인 `-builddir`로 호출할 때 `wails3 generate appimage`이 실패하는 문제 수정: 이제 번들러가 `-binary`, `-icon`, `-desktopfile`, `-builddir` 및 `-outputdir`을 처음부터 절대 경로로 해석하므로, 처리 도중의 `s.CD`로 인해 AppRun 다운로드 고루틴이나 복사 후 `ldd` 검사가 중단되지 않습니다.
- 데스크톱 `Name=` 필드가 바이너리 기본 이름과 일치하지 않을 때 `wails3 generate appimage`이 최종 AppImage를 `-outputdir`로 이동하지 못하는 문제 수정: 이제 번들러가 `OUTPUT` 환경 변수를 통해 linuxdeploy의 appimage 플러그인에서 데스크톱 파일로부터 파생된 이름 대신 `<binary>-<arch>.AppImage`에 AppImage를 쓰도록 강제합니다.
- alpha.93에서 GTK4 + WebKitGTK 6.0 스택이 기본값으로 승격된 후 Linux에서 `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` 및 `Common.SystemDidWake`이 발생하지 않는 문제 수정. 새로운 기본 `application_linux.go` `run()`에서 `setupCommonEvents()`(`Linux.*` 이벤트를 해당 `Common.*` 이벤트로 전달) 또는 `monitorPowerEvents()`을 호출하지 않았습니다. 이제 `application_linux_dbus.go`을 통해 GTK3와 GTK4 빌드 경로에서 DBus 전원 모니터 도우미를 공유합니다. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## 수정됨

- alpha.93에서 GTK4 + WebKitGTK 6.0 스택이 기본값으로 승격된 후 Linux에서 `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` 및 `Common.SystemDidWake`이 발생하지 않는 문제 수정. 새로운 기본 `application_linux.go` `run()`에서 `setupCommonEvents()`(`Linux.*` 이벤트를 해당 `Common.*` 이벤트로 전달) 또는 `monitorPowerEvents()`을 호출하지 않았습니다. 이제 `application_linux_dbus.go`을 통해 GTK3와 GTK4 빌드 경로에서 DBus 전원 모니터 도우미를 공유합니다. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## 추가됨

- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE` 추가

## 수정됨

- @leaanthony가 appmenu-gtk-module에서 실체화되지 않은 창에 접근하여 Wayland에서 창 메뉴가 충돌하는 문제 수정 (#4769)
- @leaanthony가 앱 이름에 잘못된 문자(공백, 괄호 등)가 포함될 때 GTK 애플리케이션이 충돌하는 문제 수정
- @overlordtm이 Windows에서 드래그 앤 드롭을 초기화할 때 발생하는 "메모리 부족" 오류 수정 (#4701)
- @leaanthony가 mainthread 콜백 저장소에서 맵 삭제에 잘못된 RLock을 사용하여 발생하는 경쟁 조건 수정(Linux, macOS, iOS) (#4424)
- 작업에 명령줄 인수를 전달할 때의 변수 처리 문제를 수정했습니다. 이제 KEY=VALUE 쌍으로 지정한 CLI 변수가 올바르게 초기화되고 작업 실행 전반에 전달됩니다.
- macOS에서 발생하는 NSWindowZoomButton 충돌 수정: 이제 `MaximiseButtonState`과 `FullscreenButtonState`이 시작 시점과 런타임 모두에서 더 제한적인 상태를 적용하며, 어느 setter도 다른 setter의 설정을 암묵적으로 재정의할 수 없습니다. (#5319)
- #5463에서 CodeRabbit이 발견한 기존 GTK3 빌드 경로(`-tags gtk3`)의 여러 버그 수정: 파일 연결을 통한 실행에서 더 이상 시작 핸들러를 건너뛰지 않습니다. `getTheme`은 이제 경계 및 타입 안전성을 보장합니다. `appName`은 더 이상 GLib 소유 메모리를 해제하지 않습니다. `clipboardGet`는 더 이상 GTK가 반환한 `gchar*`을 누수하지 않습니다. 이제 `Calloc`에서 포인터 리시버를 사용하고 `NewCalloc`이 `*Calloc`을 반환하므로 풀이 실제로 할당을 추적합니다. `zoomOut`은 1.0로 제한되는 음수 승수 대신 `zoomInFactor`의 역수를 사용합니다. `execJS`은 호출할 때마다 `C.CString("")`을 누수하는 대신 미리 할당된 빈 world-name을 재사용합니다. `menuItem.setAccelerator`에서 개발용 `fmt.Println`를 제거했습니다. #5465 해결.
- 기본 GTK4 빌드 경로(`linux_cgo.go`)에서 동일하게 발생하는 `Calloc` 값 리시버 누수 수정: 포인터 리시버와 `NewCalloc() *Calloc`을 사용하여 창별 `c.String(...)` 할당이 실제로 추적되고 해제되도록 했습니다.

## v3.0.0-alpha.92 - 2026-05-15

## 추가됨

- `PACKAGE_MANAGER` 옵션을 통해 사용할 프런트엔드 패키지 관리자를 제어할 수 있도록 Taskfile 수정
- Taskfile 템플릿을 더 예측 가능하게 작성할 수 있도록 템플릿 데이터에 `{{.Opn}}` 및 `{{.Cls}}` 추가

## 변경됨

- 기존 Taskfile 중 일부에서 `{{.Opn}} and {{.Cls}}`을 사용하도록 수정

## 수정됨

- 패널이 트레이 메뉴를 읽는 도중 메뉴를 업데이트할 때 `linuxSystemTray`에서 발생하는 `concurrent map read and map write` 런타임 치명적 오류 수정.
- Windows에서 연결된 콘솔 없이 앱을 실행할 때 메시지가 유실되지 않도록 WebView2 오류 및 스택 추적 출력에 `fmt` 대신 `log` 사용.

## v3.0.0-alpha.91 - 2026-05-12

## 변경됨

- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5414)에서 후원자 SVG를 업데이트했습니다.
- **호환성이 깨지는 변경 사항(macOS):** `GetScreens`, `Position` 및 `SetPosition`가 모두 같은 공간, 즉 논리적 포인트 단위이며 Y축이 아래쪽을 향하고 기본 화면의 왼쪽 위가 `(0,0)`인 좌표계를 사용하도록 macOS 좌표계를 정규화했습니다. 이는 Windows, GTK, Electron 및 웹의 공개 API와 일치합니다. 이제 기본 화면보다 물리적으로 위에 있는 화면은 음수 `Bounds.Y` 값을 보고하며(이전에는 양수), `Position()`/`SetPosition()` 값은 `points × primaryScale` 대신 논리적 포인트 단위를 사용합니다. `Position()` → `SetPosition()` 왕복 변환은 유지됩니다. 이전 알파 빌드에서 기록한 절대 좌표값이나 직접 계산한 우회 방법(예: `primaryScale`를 곱하거나 화면 높이를 기준으로 Y축을 반전하는 방법)은 업데이트해야 합니다. [#5117](https://github.com/wailsapp/wails/issues/5117) 해결.

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5416)에서 패닉을 방지하도록 DBus 신호 이름과 본문 길이를 방어적으로 검증
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5363)에서 Linux의 GTK 메뉴 처리에 발생하는 메모리 안전성 문제 수정
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5295)에서 NVIDIA GPU를 감지하고 Linux에서 DMA-BUF 렌더러 비활성화
- macOS에서 `SetPosition`의 화면 간 Y 좌표 변환 문제 수정: 기본 화면 높이를 전역 기준으로 사용하여 기본 디스플레이에서 세로 방향으로 오프셋된 모니터에서도 창이 올바른 위치에 배치되도록 했습니다. [#5117](https://github.com/wailsapp/wails/issues/5117)
- @wayneforrest가 [PR](https://github.com/wailsapp/wails/pull/5109)에서 올바른 피드백 URL을 가리키도록 git PR 템플릿 수정
- 잘못된 `DestroyMenu` 시스템 호출에서 인수 하나가 아닌 네 개를 전달하여 모든 호출이 FALSE를 반환하고 아무것도 해제하지 못해 발생하던 Windows 시스템 트레이 `SetMenu` 충돌 문제를 일괄 수정했습니다. 또한 메뉴를 다시 빌드할 때 HMENU 및 HBITMAP 핸들(런타임에 `MenuItem.SetBitmap`을 통해 할당된 핸들 포함)을 해제하고, `Win32Menu.Update`에서 오래된 체크박스/라디오 맵을 초기화하며, 할당을 두 배로 늘리던 `systemtray.updateMenu`의 불필요한 `Update()` 호출을 제거했습니다. 이제 장시간 실행되는 시스템 트레이 앱에서 메뉴를 다시 빌드할 때마다 GDI/USER 객체가 누수되지 않습니다.

## v3.0.0-alpha.90 - 2026-05-11

## 추가됨

- macOS에서 WKWebView User-Agent에 사용할 애플리케이션 이름을 구성할 수 있도록 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5261), @vinhvoit225
- gin-service 예제에 간접 종속성 github.com/coder/websocket을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5400), @taliesin-ai
- 빌드 자산 테스트에 심층 동등성 비교 지원을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5402), @leaanthony

## 변경됨

- 빌드 출력을 assets 디렉터리로 통합했습니다. [PR](https://github.com/wailsapp/wails/pull/5401), @taliesin-ai
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5399)에서 후원자 SVG를 업데이트했습니다.

## 수정됨

- macOS 단일 인스턴스 메시지에 알림 객체를 사용하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5289), @overlordtm
- 부하가 높을 때 프로미스가 유실되지 않도록 Windows 콜백을 일괄 처리합니다. [PR](https://github.com/wailsapp/wails/pull/5383), @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## 추가됨

- Go 테스트 결과를 집계하는 go<em>test</em>results 작업을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5316), @leaanthony

## 변경됨

- 대용량 RPC 페이로드를 조건에 따라 청크 단위의 POST 요청으로 분할하도록 변경했습니다. [PR](https://github.com/wailsapp/wails/pull/5369), @leaanthony
- 모든 프런트엔드 템플릿에서 Vite를 5.x.x에서 8.0.0(으)로 업그레이드했습니다. [PR](https://github.com/wailsapp/wails/pull/5386), @leaanthony
- Vite 개발 서버 포트 구성을 환경 변수로 이전했습니다. [PR](https://github.com/wailsapp/wails/pull/5365), @leaanthony
- 모든 템플릿에서 Vite 개발 서버가 127.0.0.1에 바인딩되도록 구성했습니다. [PR](https://github.com/wailsapp/wails/pull/5361), @leaanthony
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5384)에서 후원자 SVG를 업데이트했습니다.

## 수정됨

- build-assets 업데이트 중 Info.plist 템플릿 스텁을 정리하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5312), @leaanthony
- 메뉴 항목 변경 메서드(`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`)를 메인 스레드에서 동기적으로 적용하여 macOS 메뉴의 오래된 상태 문제를 수정했습니다. 이에 따라 메뉴를 빠르게 다시 열었을 때 이전 상태가 렌더링되던 `dispatch_async` 경쟁 조건이 제거되었습니다(#5002).
- 불필요한 재빌드를 방지하기 위해 개발 모드에서 `*_test.go` 파일을 무시하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5203), @leaanthony
- 앱이 실행 중이 아닐 때 Menu.Update()에서 세그멘테이션 오류가 발생하지 않도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5291), @wucm667
- Windows에서 메뉴 모음 다시 그리기 여부를 제어하는 데 lastSizeWParam을 사용하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5382), @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## 변경됨

- HiddenOnTaskbar가 WS<em>EX</em>TOOLWINDOW를 사용하도록 변경했습니다. [PR](https://github.com/wailsapp/wails/pull/5371), @leaanthony
- 종속성 순서를 변경하고 go.mod에서 webview2 replace 지시문을 제거했습니다. [PR](https://github.com/wailsapp/wails/pull/5370), @atterpac
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5358)에서 후원자 SVG를 업데이트했습니다.

## 수정됨

- 제네릭 간접 참조 별칭을 제거하고 맵 키 유형을 통합했습니다. [PR](https://github.com/wailsapp/wails/pull/5331), @fbbdev

## 제거됨

- PR-master 워크플로를 삭제하여 문서, Go 테스트 및 테스트 건너뛰기를 제거했습니다. [PR](https://github.com/wailsapp/wails/pull/5377), @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## 추가됨

- Wails v3 한국어 문서를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5352), @leaanthony
- 설치 및 빠른 시작에 대한 프랑스어 문서를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5354), @leaanthony
- 빠른 시작, 개념 및 커뮤니티에 대한 포르투갈어 문서를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5355), @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## 추가됨

- 프랑스어 문서 현지화를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5328), @leaanthony
- 문서 사이트에 독일어 로캘을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5343), @leaanthony

## 변경됨

- 번역된 로캘 8개를 모두 문서 구성에 등록했습니다. [PR](https://github.com/wailsapp/wails/pull/5347), @leaanthony
- WebView2를 위해 여러 Windows 관련 파일을 업데이트했습니다. [PR](https://github.com/wailsapp/wails/pull/5317), @leaanthony

## 수정됨

- Linux의 대화 상자 디스패치를 GTK3와 GTK4로 분리했습니다. [PR](https://github.com/wailsapp/wails/pull/5340), @leaanthony
- 대화 상자 콜백이 GTK 스레드에서 실행되도록 보장하여 세그멘테이션 오류를 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5339), @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## 추가됨

- 저장소에 PR 템플릿 URL을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5179), @leaanthony
- Wails v3 독일어 문서를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5330), @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## 추가됨

- macOS에서 Escape 키로 전체 화면을 종료하지 못하게 하는 옵션을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5307), @leaanthony
- macOS에서 Escape 키로 전체 화면을 종료하지 못하게 하는 옵션을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5310), @leaanthony
- @yuseferi가 [PR](https://github.com/wailsapp/wails/pull/5288)에서 Pausa 커뮤니티 쇼케이스 문서를 추가

## 변경됨

- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5308)에서 후원자 SVG를 업데이트했습니다.
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5309)에서 지원되지 않는 플랫폼을 처리하도록 아이콘 생성 명령을 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5224)에서 불리언 전체 화면 API를 3상태 ButtonState로 교체하고 플랫폼 바인딩을 구현

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5315)에서 컨트롤러 상태가 nil일 때 WebView2 포커스 작업을 실행하지 않도록 보호
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5313)에서 PR의 기준 브랜치를 올바르게 참조하도록 GitHub Actions 워크플로를 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5203)에서 불필요한 재빌드를 방지하기 위해 개발 모드에서 `*_test.go` 파일을 무시하도록 수정
- @wucm667이 [PR](https://github.com/wailsapp/wails/pull/5291)에서 앱이 실행 중이 아닐 때 Menu.Update()에서 세그멘테이션 오류가 발생하지 않도록 수정

## v3.0.0-alpha.83 - 2026-05-02

## 추가됨

- @symball이 [PR](https://github.com/wailsapp/wails/pull/5094)에서 머신/사용자 설치를 위한 InstallScope 플래그와 빌드 옵션을 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5294)에서 Window 인터페이스를 충족하도록 아무 작업도 하지 않는 SetScreen 메서드를 BrowserWindow에 추가

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5295)에서 NVIDIA GPU를 감지하고 Linux에서 DMA-BUF 렌더러를 비활성화하도록 수정
- @wayneforrest가 [PR](https://github.com/wailsapp/wails/pull/5109)에서 올바른 피드백 URL을 가리키도록 Git PR 템플릿을 수정
- 하나가 아닌 네 개의 인수를 전달해 모든 호출이 FALSE를 반환하고 아무것도 해제하지 못하게 했던 손상된 `DestroyMenu` 시스템 호출로 인해 발생하는 일련의 Windows 시스템 트레이 `SetMenu` 충돌을 수정. 또한 메뉴를 다시 빌드할 때 HMENU 및 HBITMAP 핸들(런타임에 `MenuItem.SetBitmap`을 통해 할당된 핸들 포함)을 해제하고, `Win32Menu.Update`에서 오래된 체크박스/라디오 맵을 초기화했으며, 할당을 두 배로 늘리던 `systemtray.updateMenu`의 중복 `Update()` 호출을 제거. 이제 장시간 실행되는 시스템 트레이 앱에서 메뉴를 다시 빌드할 때마다 GDI/USER 객체가 누수되지 않음.

## v3.0.0-alpha.82 - 2026-05-01

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5232)에서 데스크톱 이름을 올바르게 처리하도록 데스크톱 파일 생성을 수정

## v3.0.0-alpha.81 - 2026-04-30

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5286)에서 나이틀리 릴리스 일정을 15:00 UTC로 조정

## 수정됨

- Retina Mac에서 Screen Bounds, WorkArea 및 Size 값이 절반으로 줄어들던 문제 수정 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5285)에서 문서 종속성과 콘텐츠 컬렉션 로더를 업데이트

## v3.0.0-alpha.79 - 2026-04-29

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5270)에서 trigger-release 작업에 actions: write 권한을 부여

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5283)에서 릴리스 작업의 기본 브랜치를 master로 변경하고 변경 로그 문구를 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5282)에서 최신 버전을 사용하도록 자동 변경 로그 워크플로를 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5280)에서 경로 필터를 추가하고 사용되지 않는 워크플로를 제거하여 워크플로 효율성을 개선
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5274)에서 예제 링크가 master 브랜치를 참조하도록 문서를 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5272)에서 v3용 문서와 예제를 업데이트

## 수정됨

- @AkagiYui가 [PR](https://github.com/wailsapp/wails/pull/5265)에서 개발 환경용 재시도 로직과 IPv4 강제 사용 기능을 추가하여 리버스 프록시를 개선
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5281)에서 미출시 변경 로그 트리거 워크플로를 다시 작성

## 제거됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5267)에서 다양한 테스트 용도의 셸 테스트 스크립트를 제거
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5266)에서 v3-alpha 문서 배포 워크플로와 CNAME 레코드를 삭제

### 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5196)에서 사이드바 탐색 메뉴에 프런트엔드 라우팅 항목을 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5185)에서 프레임워크별 권장 사항을 포함한 프런트엔드 라우팅 가이드를 추가
- 모달 시트 지원 추가(macOS)
- Apple 기기 지원 개선을 위해 @leaanthony가 ghw 버전을 상향 (#4977)
- Dock 서비스에 `GetBadge` 메서드 추가
- 사용자 지정 Go 빌드 태그(예: `wails3 build -tags gtk4`)를 전달할 수 있도록 `wails3 build` 명령에 `-tags` 플래그 추가 (#4957)
- 전용 열거형 페이지와 사이드바 탐색 메뉴를 포함하여 바인딩 생성기의 자동 열거형 생성에 관한 문서 추가 (#4972)
- 사용자 지정 Go 빌드 태그(예: `wails3 build -tags gtk4`)를 전달할 수 있도록 `wails3 build` 명령에 `-tags` 플래그 추가 (#4957)
- Storage(localStorage, sessionStorage, IndexedDB, Cache API), Network(Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media(Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Device(Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performance(Performance API, Mutation Observer, Intersection/Resize Observer), UI(Web Components, Pointer Events, Selection, Dialog, Drag and Drop) 등을 비롯한 41 브라우저 API를 보여 주는 Web API 예제를 `v3/examples/web-apis/`에 추가
- 플랫폼 전반에서 200개 이상의 브라우저 API를 테스트하는 WebView API 호환성 검사기 예제(`v3/examples/webview-api-check/`) 추가
- 병렬 검색과 캐싱을 지원하고 Flatpak/Snap/Nix와 호환되며 Linux에서 네이티브 라이브러리 경로를 찾는 `internal/libpath` 패키지 추가
- **WIP:** Linux용 실험적 WebKitGTK 6.0 / GTK4 지원을 추가했습니다. `-tags gtk4`을 통해 사용할 수 있습니다(GTK3/WebKit2GTK 4.1은 기본값으로 유지됨).
- 참고: 타일링 창 관리자(예: Hyprland, Sway)에서는 WM이 창의 위치와 크기를 제어하므로 최소화/최대화 작업이 예상대로 작동하지 않을 수 있습니다.
- @AbdelhadiSeddar가 **JavaScript에서 이벤트 수신하기** 문서에 <strong>일회성 핸들러</strong>를 구현하는 방법을 추가했습니다.
- @leaanthony가 Windows/Linux의 창이 `app.Menu.Set()`을 통해 설정된 애플리케이션 메뉴를 상속할 수 있도록 `WebviewWindowOptions`에 `UseApplicationMenu` 옵션을 추가했습니다.
- @wimaha가 Liquid Glass 아이콘 및 애셋 카탈로그(macOS)를 생성할 때 `.icon` 파일(Apple Icon Composer 형식)을 사용할 수 있도록 지원을 추가했습니다(#4934).
- 헤드리스/웹 배포를 위한 실험적 서버 모드(`-tags server`)를 추가했습니다. 네이티브 GUI 종속성 없이 Wails 앱을 HTTP 서버로 실행할 수 있습니다. `wails3 task build:server`로 빌드하십시오. 자세한 내용은 `examples/server`에서 확인하십시오.
- 병렬 검색과 캐싱을 사용하고 Flatpak/Snap/Nix를 지원하여 Linux에서 네이티브 라이브러리 경로를 찾는 `internal/libpath` 패키지를 추가했습니다.
- @leaanthony가 macOS Spaces 및 전체 화면에서의 창 동작을 제어할 수 있도록 `MacWindow`에 `CollectionBehavior` 옵션을 추가했습니다(#4756).
- @leaanthony가 pkg/application의 단위 테스트를 추가했습니다.
- @leaanthony가 MSIX 패키징에 사용자 지정 프로토콜 지원을 추가했습니다.
- Linux 데스크톱 환경 감지를 추가했습니다. [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- @leaanthony가 프런트엔드에서 인쇄 대화 상자를 표시할 수 있도록 JavaScript 런타임에 `Window.Print()` 메서드를 추가했습니다(#4290).
- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE`을 추가했습니다.
- @leaanthony가 Linux용 WebKit2 로드 상태 변경 이벤트 `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished`을 추가했습니다(#3896).
- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE`을 추가했습니다.
- 패키징할 때뿐만 아니라 Linux 빌드 중에도 `.desktop` 파일을 생성하도록 했습니다(#4575).
- @leaanthony가 배포판별 패키지 이름과 nfpm 패키징 예제가 포함된 Linux 런타임 종속성 문서를 추가했습니다(#4339).
- @leaanthony가 Linux의 `wails3 doctor` 출력에 NVIDIA 드라이버 버전 정보를 추가했습니다.
- @APshenkin이 원시 메시지 핸들러에 origin을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4710)
- @APshenkin이 macOS용 유니버설 링크 지원을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4712)
- @APshenkin이 바인딩 전송 계층을 리팩터링했습니다. [PR](https://github.com/wailsapp/wails/pull/4702)
- @chinenual이 Appium 테스트 클라이언트에서 예제 앱을 쉽게 테스트할 수 있도록 helloworld 템플릿에 aria-label 식별자를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4760)
- @APshenkin이 원시 메시지 핸들러에 origin을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4710)
- @APshenkin이 macOS용 유니버설 링크 지원을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4712)
- @APshenkin이 바인딩 전송 계층을 리팩터링했습니다. [PR](https://github.com/wailsapp/wails/pull/4702)
- @fbbdev와 @ianvs가 형식 지정 이벤트를 추가했습니다. [#4633](https://github.com/wailsapp/wails/pull/4633)
- 실시간 툴팁 업데이트를 사용하는 헤드리스 트레이를 보여 주는 `systray-clock` 예제를 추가했습니다(#4653).
- @Tolfx가 Windows용 NSIS Protocol 템플릿을 추가했습니다(#4510).
- @Tolfx가 build-assets 테스트를 추가했습니다(#4510).
- macOS: @nidib가 메뉴 막대에 네이티브 창 컨트롤을 표시하도록 했습니다. [#4588](https://github.com/wailsapp/wails/pull/4588)
- @popaprozac이 Dock에서 앱 아이콘을 숨기거나 표시하는 macOS Dock 서비스를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4451)
- @popaprozac이 Dock에서 앱 아이콘을 숨기거나 표시하는 macOS Dock 서비스를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4451)
- @leaanthony가 NSGlassEffectView(macOS 15.0 이상) 및 NSVisualEffectView 대체 구현을 사용하는 macOS용 네이티브 Liquid Glass 효과 지원을 추가했습니다. 광범위한 재질 사용자 지정 옵션도 포함됩니다. [#4534](https://github.com/wailsapp/wails/pull/4534)
- @leaanthony가 브라우저 URL의 안전하지 않은 내용을 정제하는 기능을 추가했습니다. [#4500](https://github.dev/wailsapp/wails/pull/4500). @APShenkin의 [#4484](https://github.com/wailsapp/wails/pull/4484)을 기반으로 합니다.
- [@Taiterbase](https://github.com/Taiterbase)의 원래 작업을 기반으로 [@leaanthony](https://github.com/leaanthony)가 Windows/Mac용 콘텐츠 보호 기능을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4241)
- @leaanthony가 `wails3 build` 및 `wails3 package` 별칭을 통해 CLI 변수를 Task 명령에 전달할 수 있도록 지원을 추가했습니다(#4422). [PR](https://github.com/wailsapp/wails/pull/4488)
- 이벤트를 통해 드롭된 요소의 데이터를 제공하는 드롭존을 지원합니다. [@atterpac](https://github.com/atterpac), [#4318](https://github.com/wailsapp/wails/pull/4318)
- WebView2 브라우저에 추가 명령줄 인수를 전달할 수 있도록 `WindowsWindow` 옵션에 `AdditionalLaunchArgs`을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4467)
- wails init 실행 후 go mod tidy를 자동으로 실행하도록 추가했습니다. [@triadmoko](https://github.com/triadmoko), [PR](https://github.com/wailsapp/wails/pull/4286)
- @leaanthony가 Windows Snap Assist 기능을 추가했습니다. [PR](https://github.dev/wailsapp/wails/pull/4463)
- WebView2 브라우저에 추가 명령줄 인수를 전달할 수 있도록 `WindowsWindow` 옵션에 `AdditionalLaunchArgs`을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4467)
- wails init 실행 후 go mod tidy를 자동으로 실행하도록 추가했습니다. [@triadmoko](https://github.com/triadmoko), [PR](https://github.com/wailsapp/wails/pull/4286)
- @leaanthony가 Windows Snap Assist 기능을 추가했습니다. [PR](https://github.dev/wailsapp/wails/pull/4463)
- [@almas-x](https://github.com/almas-x)가 Windows용 `getAccentColor` 구현을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4427)
- [@almas-x](https://github.com/almas-x)가 Windows용 `getAccentColor` 구현을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/4427)
- Windows 다크 테마 메뉴 및 메뉴 모음. @leaanthony가 [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)에서 추가했습니다.
- @popaprozac이 JS/TS 바인딩을 더 명확하게 만들기 위해 기본 제공 서비스의 이름을 변경했습니다. [PR](https://github.com/wailsapp/wails/pull/4405)
- 사용자 시스템의 강조색을 가져오는 `app.Env.GetAccentColor`입니다. MacOS에서 작동합니다. [@etesam913](https://github.com/etesam913)
- [@atterpac](https://github.com/atterpac)이 `window.ToggleFrameless()` API를 추가했습니다. [#4137](https://github.com/wailsapp/wails/pull/4137)
- @leaanthony이 [PR](https://github.com/wailsapp/wails/pull/4345)에서 Linux 배포판별 빌드 종속성을 추가했습니다.
- @atterpac이 [PR](https://github.com/wailsapp/wails/pull/4404)에서 바인딩 가이드를 추가했습니다.
- **테스트 인프라 정리**: [@leaanthony](https://github.com/leaanthony)이 [#4359](https://github.com/wailsapp/wails/pull/4359)에서 Docker 테스트 파일을 전용 `test/docker/` 디렉터리로 옮기고 이미지를 최적화하여 빌드 안정성을 높였습니다.
- **리소스 관리 패턴 개선**: [@leaanthony](https://github.com/leaanthony)이 [#4359](https://github.com/wailsapp/wails/pull/4359)에서 예제에 올바른 이벤트 핸들러 정리와 컨텍스트를 인식하는 goroutine 관리를 추가했습니다.
- [@AkshayKalose](https://github.com/AkshayKalose)가 [#3981](https://github.com/wailsapp/wails/pull/3981)에서 aarch64 AppImage 빌드를 지원하도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)이 `wails doctor`에 진단 섹션을 추가했습니다.
- [@leaanthony](https://github.com/leaanthony)이 서비스 메서드를 호출할 때 컨텍스트에 창을 추가하도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)이 어느 창에서 서비스를 호출하는지 확인하는 방법을 보여 주는 `window-call` 예제를 추가했습니다.
- [@leaanthony](https://github.com/leaanthony)이 새로운 메뉴 가이드를 추가했습니다.
- [@leaanthony](https://github.com/leaanthony)이 panic 처리를 개선했습니다.
- [@leaanthony](https://github.com/leaanthony)이 새로운 메뉴 가이드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4024](https://github.com/wailsapp/wails/pull/4024)에서 Service API에 문서 주석을 추가했습니다.
- [@leaanthony](https://github.com/leaanthony)이 [#4024](https://github.com/wailsapp/wails/pull/4024)에서 추가 구성으로 서비스를 초기화하는 `application.NewServiceWithOptions` 함수를 추가했습니다.
- [@FalcoG](https://github.com/FalcoG)와 [@leaanthony](https://github.com/leaanthony)이 [#4031](https://github.com/wailsapp/wails/pull/4031)에서 메뉴 제어를 개선했습니다.
- [@leaanthony](https://github.com/leaanthony)이 문서를 추가로 작성했습니다.
- [@leaanthony](https://github.com/leaanthony)이 표준 이벤트 리스너에서 이벤트 취소를 지원하도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)이 Systray의 `Hide`, `Show` 및 `Destroy` 지원을 추가했습니다.
- [@leaanthony](https://github.com/leaanthony)이 Systray의 `SetTooltip` 지원을 추가했습니다. 최초 아이디어는 [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)이 제안했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 지원되지 않는 형식에 관한 바인딩 생성기 경고에 패키지 경로를 표시하도록 했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 바인딩 생성기에 제네릭 별칭 지원을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 바인딩 생성기에 `omitzero` JSON 플래그 지원을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 선택한 서비스 메서드의 바인딩 생성을 방지하는 `//wails:ignore` 지시문을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 Go에는 내보내지만 JS/TS에는 내보내지 않는 형식을 허용하도록 서비스와 모델에 `//wails:internal` 지시문을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 약한 형식의 열거형을 사용할 수 있도록 별칭 형식 상수에 대한 바인딩 생성기 지원을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4068](https://github.com/wailsapp/wails/pull/4068)에서 Go 1.24 기능에 대한 바인딩 생성기 테스트를 추가했습니다.
- [#4065](https://github.com/wailsapp/wails/pull/4065)에서 OS 버전 감지를 개선하기 위해 `OSInfo.Branding`에 macOS 15 "Sequoia" 지원을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 종료 절차가 완료된 후 사용자 지정 코드를 실행하는 `PostShutdown` 훅을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 사용자 지정 오류 핸들러가 치명적 오류를 감지할 수 있도록 `FatalError` 구조체를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 서비스 시작 및 종료 순서를 표준화하고 문서화했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 애플리케이션 시작/종료 순서와 서비스 시작/종료 테스트를 위한 테스트 하네스를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 애플리케이션 생성 후 서비스를 등록하는 `RegisterService` 메서드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 바인딩 호출의 사용자 지정 오류 처리를 위해 애플리케이션 및 서비스 옵션에 `MarshalError` 필드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 promise 체인을 통해 취소 요청을 전파하는 취소 가능 promise 래퍼를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 바인딩 호출 취소를 `AbortSignal`에 연결할 수 있도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)이 WML에서 일반적인 `wml-*` 특성과 함께 `data-wml-*` 특성도 지원하도록 했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 지연 구성 및 동적 재구성을 위해 모든 서비스에 `Configure` 메서드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 `fileserver` 서비스가 구성되지 않은 경우 503 Service Unavailable 응답을 보내도록 했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 `kvstore` 서비스가 구성되지 않은 경우 기본적으로 메모리 내 키-값 저장소를 제공하도록 했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 구성 변경 후 파일의 데이터를 다시 불러오도록 `kvstore` 서비스에 `Load` 메서드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 모든 키를 삭제하도록 `kvstore` 서비스에 `Clear` 메서드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 JS 측 로그 수준 상수를 제공하도록 `log` 서비스에 `Level` 형식을 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 로그 수준을 동적으로 지정하도록 `log` 서비스에 `Log` 메서드를 추가했습니다.
- [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 `sqlite` 서비스가 구성되지 않은 경우 기본적으로 메모리 내 DB를 제공하도록 했습니다.
- DB를 수동으로 닫을 수 있도록 `sqlite` 서비스에 `Close` 메서드 추가 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- `sqlite` 서비스의 쿼리 메서드에 취소 지원 추가 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- `sqlite` 서비스에 JS 바인딩을 포함한 준비된 문 지원 추가 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- Gin 지원 — [@AnalogJ](https://github.com/AnalogJ)가 이 [PR](https://github.com/wailsapp/wails/pull/3537)에서 수행한 원래 작업을 바탕으로 [PR](https://github.com/wailsapp/wails/pull/3537)에서 [Lea Anthony](https://github.com/leaanthony)
- 자동 저장과 비밀번호 자동 저장이 항상 활성화되는 문제 수정 — [#4134](https://github.com/wailsapp/wails/pull/4134)에서 [@oSethoum](https://github.com/osethoum)
- 창에 메뉴를 설정할 수 있도록 창에 `SetMenu()` 추가 — [@leaanthony](https://github.com/leaanthony)
- 알림 지원 추가 — [#4098](https://github.com/wailsapp/wails/pull/4098)에서 [@popaprozac](https://github.com/popaprozac)
-  mac용 파일 연결 지원 추가 — [#4177](https://github.com/wailsapp/wails/pull/4177)에서 [@wimaha](https://github.com/wimaha)
- 시맨틱 버전 증가를 위한 `wails3 tool version` 추가 — [@leaanthony](https://github.com/leaanthony)
- macOS 및 Windows용 배지 지원 추가 — [#](https://github.com/wailsapp/wails/pull/4234)에서 [@popaprozac](https://github.com/popaprozac)
- 등록된/엄격한 형식의 이벤트 지원 추가 — [#4161](https://github.com/wailsapp/wails/pull/4161)에서 [@fbbdev](https://github.com/fbbdev) 및 [@IanVS](https://github.com/IanVS)
- 사용자 정의 이벤트용 훅을 등록하는 기능 추가 — [#4161](https://github.com/wailsapp/wails/pull/4161)에서 [@fbbdev](https://github.com/fbbdev) 및 [@IanVS](https://github.com/IanVS)
- `path` 경로에서 시스템 파일 관리자를 열고 `selectFile`을 통해 선택적으로 강조 표시하는 `app.OpenFileManager(path string, selectFile bool)` — [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- `wails3 init` 명령에 새로운 `-git` 플래그 추가 — [@leaanthony](https://github.com/leaanthony)
- 새로운 `wails3 generate webview2bootstrapper` 명령 — [@leaanthony](https://github.com/leaanthony)
- 런타임을 수동으로 초기화할 수 있도록 런타임에 `init()` 메서드 추가 — [@leaanthony](https://github.com/leaanthony)
- Window의 WindowOptions에 `WindowDidMoveDebounceMS` 옵션 추가 — [@leaanthony](https://github.com/leaanthony)
- 단일 인스턴스 기능 추가 — [@leaanthony](https://github.com/leaanthony). @APshenkin의 [v2 PR](https://github.com/wailsapp/wails/pull/2951)을 기반으로 함.
- `wails3 generate template` 명령 — [@leaanthony](https://github.com/leaanthony)
- `wails3 releasenotes` 명령 — [@leaanthony](https://github.com/leaanthony)
- `wails3 update cli` 명령 — [@leaanthony](https://github.com/leaanthony)
- `wails3 generate bindings` 명령용 `-clean` 옵션 — [@leaanthony](https://github.com/leaanthony)
- aarch64(arm64) AppImage Linux 빌드 지원 — [#3981](https://github.com/wailsapp/wails/pull/3981)에서 [@AkshayKalose](https://github.com/AkshayKalose)
- 후원자용 하이퍼링크 추가 — [#3958](https://github.com/wailsapp/wails/pull/3958)에서 @ansxuman
- deb, rpm 및 Arch Linux 패키저 빌드를 위한 Linux 패키징 지원 —
- Darwin 유니버설 빌드 및 패키지 지원 추가 —
- 웹사이트에 이벤트 문서 추가 —
- 비 SSR 개발용으로 설정된 sveltekit 및 sveltekit-ts 템플릿
- 새로운 `wails3 update build-assets` 명령을 사용하여 빌드 자산 업데이트 —
- HTML Drag and Drop API를 테스트하는 예제 —
- 파일 연결 지원 — 다음에서 [leaanthony](https://github.com/leaanthony):
- 새로운 `wails3 generate runtime` 명령 —
- 창을 가운데에 배치할지 또는 다음과 같이 할지 지정하는 새로운 `InitialPosition` 옵션:
- `application` 패키지에 `Path` 및 `Paths` 메서드 추가 —
- `GeneralAutofillEnabled` 및 `PasswordAutosaveEnabled` Windows 옵션 추가
- 서비스 메서드를 호출한 창을 가져오는 기능 추가 —
- WebView2용 `EnabledFeatures` 및 `DisabledFeatures` 옵션 추가 —
- ⊞ 고DPI 모니터 지원을 개선하는 새로운 DIP 시스템 —
- ⊞ 창 클래스 이름 옵션 — 기여자: [windom](https://github.com/windom/), 관련 참조:
- 플러그인 기능을 제공하도록 서비스를 확장함. 작성자:
- 🐧 다음의 WindowDidMove / WindowDidResize 이벤트:
- ⊞ 다음의 WindowDidResize 이벤트:
-  Dock을 처리할 수 있도록 ApplicationShouldHandleReopen 이벤트 추가
-  구현에 getPrimaryScreen/getScreens 추가 — 다음에서 @tmclane:
-  macOS 전체 화면 모드에서 도구 모음을 표시하는 옵션 추가 —
- 🐧 Linux 키 입력을 단축키로 변환하는 onKeyPress 로직 추가
- 🐧 `run:linux` 태스크 추가 —
- `SetIcon` 메서드 내보내기 — 다음에서 [@almas-x](https://github.com/almas-x):
- `OnShutdown` 개선 — 다음에서 [@almas-x](https://github.com/almas-x):
- `Window` 인터페이스에 `ToggleMaximise` 메서드 복원 —
- `Environment()`에 추가 정보 보강. 다음에서 @leaanthony:
- `Window` 인터페이스에 `WebviewWindow.IsFocused` 메서드 공개 —
- WML 시스템에서 공백으로 구분된 여러 트리거 이벤트 지원 —
- 번들 JS 런타임 스크립트에서 ESM 내보내기 추가 —
- 다음 대신 번들 JS 런타임 스크립트를 사용하도록 바인딩 생성기 플래그 추가:
- Linux에서 `setIcon` 구현 — [@abichinger](https://github.com/abichinger)
- dev 명령에 `-port` 플래그를 추가하고 환경 변수 지원
- 바인딩된 메서드 호출 테스트 추가 — 기여자:
- ⊞ 이미 생성된 창에 `SetIgnoreMouseEvents` 추가 — 기여자:
-  창의 스태킹 레벨(순서)을 설정하는 기능 추가 — 기여자:

### 수정됨

- NSScreen 포인트 값을 `Physical*` 필드의 장치 픽셀로 변환하여 Retina Mac에서 `Screen.Bounds`, `WorkArea`, `Size` 값이 절반으로 표시되던 문제를 수정하고, [PR](https://github.com/wailsapp/wails/pull/5168)에서 다중 모니터의 맞닿음 감지와 작업 영역 배치가 올바르게 작동하도록 최상위 `Screen.X`/`Y` 값을 채움 — @wayneforrest
- 디스플레이 구성이 변경될 때(예: 절전 모드 진입/해제 중 외부 모니터 핫플러그) WebKit DisplayLink 교착 상태를 일으키는 ScreenManager의 데이터 경합 수정
- [PR](https://github.com/wailsapp/wails/pull/5154)에서 Assets.car가 있으면 CFBundleIconName을 appicon으로 직접 설정 — @symball
- Fedora, openSUSE, Arch, NixOS에서 `wails3 doctor`이 잘못된 WebKitGTK 패키지를 보고하던 문제 수정 — v3에서는 컴파일 시 4.1 API가 필요하므로 4.0 대체 항목 제거(#5071)
- openSUSE의 webkit2gtk doctor 패키지 이름 수정(`webkit2gtk4_1-devel` → 올바른 openSUSE 패키지 이름인 `webkit2gtk3-devel`)(#5071)
- 데스크톱 개발 모드에서 `/wails/custom.js`이 없을 때 발생하는 `Unexpected token '<'` 오류 수정. HTML SPA 대체 응답이 JavaScript로 삽입되지 않도록 `/wails/custom.js`용 명시적 404 핸들러와 `loadOptionalScript`의 대소문자를 구분하지 않는 `Content-Type` 검증 추가. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- macOS의 시스템 트레이 메뉴 강조 상태 수정 — 이제 메뉴가 열려 있으면 아이콘이 선택된 상태로 표시됨(#4910)
- macOS에서 시스템 트레이에 연결된 창이 다른 창 뒤에 나타나던 문제 수정 — 이제 적절한 팝업 창 레벨 사용(#4910)
- 문서 전반의 잘못된 `@wailsio/runtime` 가져오기 예제 수정(#4989)
- darwin에서 프레임 없는 창을 최소화할 수 없던 문제 수정(#4294)
- go-task의 최신 상태 검사에서 `node_modules/`을 제외하여 `wails3 build` 및 `wails3 dev` 실행 중 20-30분 동안 멈추던 현상을 수정했습니다. 이전에는 `sources: "**/*"` glob 패턴 때문에 go-task가 `node_modules/`의 모든 파일(MUI 같은 대규모 의존성을 사용하는 경우 50000-100000개 이상)을 열거하고 각 파일의 체크섬을 계산했으며, 특히 Windows/NTFS에서 속도가 느렸습니다(#4939).
- C `Screen` typedef가 X11 Xlib.h와 충돌하여 발생하던 GTK4 빌드 실패 수정(#4957)
- macOS에서 Dock 배지 메서드의 일관성 수정
- `InvisibleTitleBarHeight`이 프레임이 없거나 제목 표시줄이 투명한 창에만 적용되지 않고 모든 macOS 창에 적용되던 문제 수정(#4960)
- `InvisibleTitleBarHeight`이 활성화된 상태에서 위쪽 모서리로 창 크기를 조절할 때 창이 흔들리던 문제를 창 가장자리 근처에서 드래그 시작을 건너뛰도록 하여 수정(#4960)
- JS/TS 바인딩에서 enum 키를 사용하는 매핑된 타입의 생성 문제 수정(#4437) — @fbbdev
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간(물리적 픽셀과 CSS 픽셀)에 있던 문제 수정
- Linux에서 호버 효과를 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않던 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- `gtk_window_present()`을 사용하여 Linux/GTK4에서 창을 표시하거나 숨길 때 가끔 최소화된 상태로 복원되던 문제 수정(#4957)
- `XTranslateCoordinates`/`XMoveWindow`를 통한 X11 조건부 지원을 추가하여 Linux/GTK4에서 창 위치 가져오기/설정이 항상 0,0을 반환하던 문제 수정(#4957)
- 제거된 `gtk_window_set_geometry_hints`을 대체하는 신호 기반 크기 제한을 추가하여 Linux/GTK4에서 최대 창 크기가 적용되지 않던 문제 수정(#4957)
- 올바른 PhysicalBounds 계산과 `gdk_monitor_get_scale`을 통한 소수 배율 조정 지원을 구현하여 Linux/GTK4의 DPI 배율 조정 문제 수정(GTK 4.14 이상)
- Linux/GTK4에서 새 창을 만들 때 메뉴 항목이 중복되던 문제 수정
- JS/TS 바인딩에서 enum 키를 사용하는 매핑된 타입의 생성 문제 수정(#4437) — @fbbdev
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간(물리적 픽셀과 CSS 픽셀)에 있던 문제 수정
- Linux에서 호버 효과를 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않던 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- 올바른 PhysicalBounds 계산과 `gdk_monitor_get_scale`을 통한 소수 배율 조정 지원을 구현하여 Linux/GTK4의 DPI 배율 조정 문제 수정(GTK 4.14 이상)
- Linux/GTK4에서 새 창을 만들 때 메뉴 항목이 중복되던 문제 수정
- JS/TS 바인딩에서 enum 키를 사용하는 매핑된 타입의 생성 문제 수정(#4437) — @fbbdev
- App.Window.Current()에서 Main Thread를 통해 AppKit API에 접근하지 않아 macOS에서 발생하던 "고스트 창" 문제 수정(#4947) — @wimaha
- WKUIDelegate runOpenPanelWithParameters를 구현하여 macOS에서 HTML `<input type="file">`이 작동하지 않던 문제 수정(#4862)
- macOS/Linux에서 `@wailsio/runtime` npm 모듈을 사용할 때 네이티브 파일 드래그 앤 드롭이 작동하지 않던 문제 수정(#4953) — @leaanthony
- 패키지 간 타입 별칭의 바인딩 생성 문제 수정(#4578) — @fbbdev
- GTK 스레드 안전성 위반으로 인해 Linux에서 발생하던 OpenFileDialog 충돌 수정(#3683) — @ddmoney420
- 숨겨졌거나 제거된 창에서 `Focus()`을 호출할 때 발생하던 SIGSEGV 충돌 수정(#4890) — @ddmoney420
- Linux에서 빈 아이콘이나 비트맵을 설정할 때 발생할 수 있던 패닉 수정(#4923) — @ddmoney420
- macOS에서 서비스 바인딩을 통해 호출할 때 발생하던 ErrorDialog 충돌 수정(#3631) — @leaanthony
- `v3\examples\dialogs`에서 Windows OS에 메뉴가 표시되도록 수정 — @ndianabasi
- 페이지를 다시 불러오는 동안 TypeError를 일으키던 경합 상태 수정(#4872) — @ddmoney420
- `Collector.IsVoidAlias()` 메서드에서 전역 상태를 제거하여 바인딩 생성기 테스트의 잘못된 출력 수정(#4941) — @fbbdev
- macOS에서 `<input type="file">` 파일 선택기가 작동하지 않던 문제 수정(#4862) — @leaanthony
- macOS에서 `Position()`과 `SetPosition()`이 서로 다른 좌표계를 사용하여 상태 저장 및 복원 시 창 위치가 어긋나는 문제 수정(#4816), @leaanthony
- 애플리케이션 매니페스트를 통해 DPI 인식이 이미 설정된 경우 발생하는 SetProcessDpiAwarenessContext "Access is denied" 오류 수정(#4803)
- 키보드 단축키 문서 페이지를 업데이트하고 `KeyBinding.Add`의 콜백 매개변수 타입 수정, @ndianabasi
- 사용자 정의 바인딩 생성 관련 문서를 수정하여 `-o String` 대신 `-d String`을 사용하도록 변경
- `menu.Update()`에서 메뉴의 자식 항목이 지워지지 않는 문제 수정
- 문서의 오래된 Manager API 참조 수정(31 파일에서 `app.Window.New()`, `app.Event.Emit()` 등과 같은 새 패턴을 사용하도록 업데이트), @leaanthony
- WebKit이 시그널 핸들러를 재정의하여 JS에 바인딩된 Go 메서드에서 패닉 발생 시 Linux에서 애플리케이션이 충돌하는 문제 수정(#3965), @leaanthony
- Linux에서 SaveFileDialog.SetFilename()이 적용되지 않는 문제 수정(#4841), @samstanier
- 드래그 앤 드롭 예제에서 드롭 좌표가 undefined로 표시되는 문제 수정
- APP_NAME에 공백이 있으면 macOS 앱 번들 생성에 실패하는 문제 수정(중괄호 확장 문제)
- Windows에서 서비스 메서드 호출 시 발생하는 인덱스 범위 초과 패닉 수정(goccy/go-json 되돌림)
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않는 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않는 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간으로 전달되는 문제 수정(물리적 픽셀과 CSS 픽셀)
- Linux에서 호버 효과와 함께 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않는 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않는 문제 수정
- `APP_NAME` 같은 변수에 포함된 공백을 처리할 수 있도록 모든 운영 체제용 Taskfile.yml 파일의 모든 명령 업데이트, @ndianabasi
- Linux에서 'build:universal:lipo:go' 태스크 실행 시 발생하는 명령 인수 오류 수정, @wux1an
- Linux에서 'wails3 build GOOS=darwin GOARCH=arm64' 실행 시 발생하는 Docker 오류 "undefined symbol: **<em>ubsan</em>handle_xxxxxxx" 수정, @wux1an
- 사용자 정의 프로토콜 문서를 통합하고 Universal Links 섹션 추가, @leaanthony
- TrackPopupMenuEx 동시 호출을 방지하는 가드를 추가하여 아이콘을 반복해서 클릭할 때 Windows 시스템 트레이 메뉴가 충돌하는 문제 수정(#4151), @leaanthony
- app.Run()보다 먼저 systray.Run()을 호출할 때 앱이 충돌하지 않도록 수정, @leaanthony
- ApplicationShouldTerminateAfterLastWindowClosed가 활성화된 상태에서 Hide()/Show()로 창 표시 여부를 전환할 때 macOS에서 충돌하는 문제 수정(#4389), @leaanthony
- macOS와 Windows에서 메뉴를 반복해서 열 때 컨텍스트 메뉴에서 발생하는 메모리 누수 수정(#4012), @leaanthony
- macOS에서 컨텍스트 메뉴의 네이티브 리소스가 재사용되지 않아 메뉴를 표시할 때마다 새 메뉴가 생성되는 문제 수정(#4012), @leaanthony
- 앱을 `Hidden: true`으로 시작한 경우 macOS Dock 아이콘을 클릭해도 숨겨진 창이 표시되지 않는 문제 수정(#4583), @leaanthony
- CGO 호출에서 잘못된 창 포인터 타입을 사용하여 macOS 인쇄 대화 상자가 열리지 않는 문제 수정(#4290), @leaanthony
- appmenu-gtk-module이 아직 실체화되지 않은 창에 접근하여 Wayland에서 창 메뉴가 충돌하는 문제 수정(#4769), @leaanthony
- 앱 이름에 유효하지 않은 문자(공백, 괄호 등)가 포함된 경우 GTK 애플리케이션이 충돌하는 문제 수정, @leaanthony
- Windows에서 드래그 앤 드롭 초기화 시 발생하는 "not enough memory" 오류 수정(#4701), @overlordtm
- 잘못된 URI 이스케이프 처리로 인해 Linux 파일 탐색기에서 엉뚱한 디렉터리가 열리는 문제 수정(#4397), @leaanthony
- `.relr.dyn` ELF 섹션을 자동 감지하고 스트리핑을 비활성화하여 최신 Linux 배포판(Arch, Fedora 39+, Ubuntu 24.04+)에서 AppImage 빌드가 실패하는 문제 수정(#4642), @leaanthony
- Fedora/DNF 기반 시스템에서 `wails doctor`이 webkit 패키지가 설치되었다고 잘못 보고하는 문제 수정(#4457), @leaanthony
- 기본 `config.yml`이 프로덕션 빌드와 함께 `wails3 dev`을 실행하던 문제 수정, @mbaklor
- 존재하지 않는 패키지를 임포트하여 iOS 서비스 스텁에서 빌드가 실패하는 문제 수정, @leaanthony
- debug/info 메서드의 구조화된 로깅에서 "no formatting directives" 오류가 발생하는 문제 수정, @leaanthony
- 모바일 플랫폼 병합 과정에서 실수로 포함된 임시 디버그 출력문 제거, @leaanthony
- DMA-BUF 렌더러를 자동으로 비활성화하여 NVIDIA GPU를 사용하는 Wayland에서 WebKitGTK가 충돌하는 문제 수정(Error 71 Protocol error), @leaanthony
- Linux에서 `application.WebviewWindowOptions.BackgroundColour`의 알파 값이 무시되는 문제 해결([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- 사용자 정의 아이콘을 제공하지 않은 경우 Windows 시스템 트레이 아이콘이 애플리케이션 아이콘을 기본값으로 사용하지 않는 문제 수정(#4704)
- 사용자가 생성한 핸들만 제거하도록 `HICON` 소유권을 추적하여 Explorer 재시작 시 발생하는 충돌을 방지했습니다(#4653).
- destroy 중 Windows 시스템 테마 리스너와 유지 중인 트레이 아이콘을 해제하여 goroutine 및 디바이스 컨텍스트 누수 방지(#4653).
- 서로게이트 쌍과 멀티바이트 글리프가 손상되지 않도록 트레이 도구 설명을 127 UTF-16 단위에서 잘라냄(#4653).
- Windows 패키지 태스크 실패 문제 수정(#4667)
- Linux taskfile의 Linux appimage appicon 변수 수정 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- go-webview2 v1.0.22의 시그니처 변경으로 발생하는 Windows 빌드 오류 수정(#4513, #4645)
- Linux taskfile의 Linux appimage appicon 변수 수정 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- `<.Info.Protocol>`을 `<.Protocol>`로 변경하여 Linux desktop.tmpl의 프로토콜 범위 수정, @Tolfx, #4510
- liquid glass 데모의 재정의 오류 수정, [#4542](https://github.com/wailsapp/wails/pull/4542), @Etesam913
- Linux의 시스템 트레이 메뉴 업데이트 수정 [#4604](https://github.com/wailsapp/wails/issues/4604), [@JackDoan](https://github.com/JackDoan)
- Windows에서 숨겨진 창을 생성할 때 흰색 창이 나타나는 문제 수정, @leaanthony, [#4612](https://github.com/wailsapp/wails/pull/4612)
- 문서의 notifications 패키지 임포트 경로 수정, @rxliuli, [#4617](https://github.com/wailsapp/wails/pull/4617)
- npm 패키지 @wailsio/runtime 사용 시 드래그 앤 드롭이 작동하지 않는 문제 수정(#4489, @leaanthony, #4616)
- Windows: 시작 시 창이 깜박이고 숨겨진 창이 잘못 표시되는 문제 수정([PR](https://github.com/wailsapp/wails/pull/4600), @leaanthony).
- Wayland에서 창 크기를 최대화할 때 발생하는 문제 수정(https://github.com/wailsapp/wails/issues/4429, [@samstanier](https://github.com/samstanier))
- Wayland에서 창 크기를 최대화할 때 발생하는 문제 수정(https://github.com/wailsapp/wails/issues/4429, [@samstanier](https://github.com/samstanier))
- 리퀴드 글래스 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542), @Etesam913)
- MacOS에서 AssetServer가 충돌할 수 있는 문제 수정([#4576](https://github.com/wailsapp/wails/pull/4576), @jghiloni)
- NextJs로 빌드할 때 발생하는 컴파일 문제 수정([#4585](https://github.com/wailsapp/wails/pull/4585), @rev42)
- 나이틀리 릴리스 파이프라인 수정([#4597](https://github.com/wailsapp/wails/pull/4597), @riadafridishibly)
- 리퀴드 글래스 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542), @Etesam913)
- MacOS에서 AssetServer가 충돌할 수 있는 문제 수정([#4576](https://github.com/wailsapp/wails/pull/4576), @jghiloni)
- NextJs로 빌드할 때 발생하는 컴파일 문제 수정([#4585](https://github.com/wailsapp/wails/pull/4585), @rev42)
- 나이틀리 릴리스 파이프라인 수정([#4597](https://github.com/wailsapp/wails/pull/4597), @riadafridishibly)
- 리퀴드 글래스 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542), @Etesam913)
- Windows의 SetBackgroundColour 수정(@PPTGamer, [PR](https://github.com/wailsapp/wails/pull/4492))
- Manager API 리팩터링의 변경 사항을 반영하도록 문서 업데이트(@yulesxoxo, [PR #4476](https://github.com/wailsapp/wails/pull/4476))
- Linux Taskfile에 있는 Linux .desktop 파일의 appicon 변수 수정([PR #4477](https://github.com/wailsapp/wails/pull/4477))
- Manager API 리팩터링의 변경 사항을 반영하도록 문서 업데이트(@yulesxoxo, [PR #4476](https://github.com/wailsapp/wails/pull/4476))
- [#4456](https://github.com/wailsapp/wails/issues/4456)에서 보고된 Windows의 nil 포인터 역참조 버그 수정(@leaanthony, [#4460](https://github.com/wailsapp/wails/pull/4460))
- 두 손가락 스와이프 탐색 제스처를 활성화하도록 macOS WKWebView에서 `allowsBackForwardNavigationGestures` 지원 추가(#1857)
- 처음에 비활성화된 상태로 설정된 메뉴 항목에서 onClick이 작동하지 않는 문제 수정(@leaanthony, [PR #4469](https://github.com/wailsapp/wails/pull/4469)). 초기 조사에 도움을 준 @IanVS에게 감사드립니다.
- 빌드 실패 시 Vite 서버가 정리되지 않는 문제 수정(#4403)
- Windows에서 `SaveFileDialog`을 닫거나 취소할 때 발생하는 패닉 수정([PR](https://github.com/wailsapp/wails/pull/4284), @hkhere)
- Windows의 HTML 수준 드래그 앤 드롭 수정([@mbaklor](https://github.com/mbaklor), [#4259](https://github.com/wailsapp/wails/pull/4259))
- 두 손가락 스와이프 탐색 제스처를 활성화하도록 macOS WKWebView에서 `allowsBackForwardNavigationGestures` 지원 추가(#1857)
- 처음에 비활성화된 상태로 설정된 메뉴 항목에서 onClick이 작동하지 않는 문제 수정(@leaanthony, [PR #4469](https://github.com/wailsapp/wails/pull/4469)). 초기 조사에 도움을 준 @IanVS에게 감사드립니다.
- 빌드 실패 시 Vite 서버가 정리되지 않는 문제 수정(#4403)
- Windows의 알림 구문 분석 수정(@popaprozac, [PR](https://github.com/wailsapp/wails/pull/4450))
- Windows SDK 종속성을 확인하도록 doctor 명령 수정([@kodumulo](https://github.com/kodumulo), [#4390](https://github.com/wailsapp/wails/issues/4390))
- Mac용 processURLRequest의 nil 포인터 역참조 수정([@etesam913](https://github.com/etesam913), [#4366](https://github.com/wailsapp/wails/pull/4366))
- 필터가 적용된 대화 상자를 사용할 수 없게 하던 Linux 버그 수정([@bh90210](https://github.com/bh90210), [#4287](https://github.com/wailsapp/wails/pull/4287))
- Windows 및 Linux의 편집 메뉴 문제 수정([@leaanthony](https://github.com/leaanthony), [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062))
- macOS .plist 파일의 최소 시스템 버전을 10.13.0에서 10.15.0(으)로 업데이트([@AkshayKalose](https://github.com/AkshayKalose), [#3981](https://github.com/wailsapp/wails/pull/3981))
- 창 ID 건너뛰기 문제 수정([@leaanthony](https://github.com/leaanthony))
- RegisterContextMenu 호출 시 발생하는 nil 메뉴 문제 수정([@leaanthony](https://github.com/leaanthony))
- 바인딩 생성기 출력의 종속성 순환 수정([@fbbdev](https://github.com/fbbdev), [#4001](https://github.com/wailsapp/wails/pull/4001))
- 바인딩 생성기 출력의 정의 전 사용 오류 수정([@fbbdev](https://github.com/fbbdev), [#4001](https://github.com/wailsapp/wails/pull/4001))
- 바인딩 생성기에 빌드 플래그 전달([@fbbdev](https://github.com/fbbdev), [#4023](https://github.com/wailsapp/wails/pull/4023))
- Windows가 아닌 플랫폼에서도 작동하도록 windows Taskfile의 경로를 슬래시 형식으로 변경([@leaanthony](https://github.com/leaanthony))
- Mac 및 Mac JS 이벤트 수정([@leaanthony](https://github.com/leaanthony))
- macOS의 이벤트 교착 상태 수정([@leaanthony](https://github.com/leaanthony))
- Windows에서 HTML은 제공되었지만 JS는 제공되지 않았을 때 Window 초기화 중 발생하는 `Parameter incorrect` 오류 수정([@leaanthony](https://github.com/leaanthony))
- 에셋 서버에서 콘텐츠 유형 감지에 사용하는 응답 접두부의 크기 수정([@fbbdev](https://github.com/fbbdev), [#4049](https://github.com/wailsapp/wails/pull/4049))
- 에셋 서버의 루트 인덱스 경로에서 404이 아닌 응답을 처리하는 방식 수정([@fbbdev](https://github.com/fbbdev), [#4049](https://github.com/wailsapp/wails/pull/4049))
- 바인딩 생성기에서 제네릭 타입의 속성을 검사할 때 발생하는 정의되지 않은 동작 수정([@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045))
- 기반 타입의 속성이 명명된 래퍼와 동일하지 않은 경우 모델에 대한 바인딩 생성기 출력 수정([@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045))
- 맵 키 타입 및 전처리에 대한 바인딩 생성기 출력 수정([@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045))
- marshaler 인터페이스를 구현하는 구조체에 대한 바인딩 생성기 출력 수정([@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045))
- 바인딩 생성기에서 제네릭 타입이 관련된 타입 순환 감지 수정([@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045))
- 바인딩 생성기 출력에서 내보내지 않은 모델에 대한 잘못된 참조를 수정했습니다. [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- 삽입된 코드를 서비스 파일 끝으로 이동했습니다. [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- 바인딩 생성기에서 파일 닫기 작업의 오류 처리를 수정했습니다. [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- 수명 주기 또는 http 메서드는 정의하지만 그 밖의 바인딩된 메서드는 정의하지 않는 서비스에 대한 경고를 표시하지 않도록 했습니다. [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- React를 사용하지 않는 템플릿에서 밝은 시스템 색 구성표를 사용할 때 Hello World 바닥글이 표시되지 않던 문제를 수정했습니다. [@marcus-crane](https://github.com/marcus-crane), [#4056](https://github.com/wailsapp/wails/pull/4056)
- macOS에서 숨겨진 메뉴 항목과 관련된 문제를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- 메시지 프로세서의 오류 처리 및 서식을 수정했습니다. [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
-  애플리케이션 종료 시 서비스 종료 절차를 건너뛰던 문제를 수정했습니다. [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
-  메뉴 업데이트가 메인 스레드에서 실행되도록 했습니다. [@leaanthony](https://github.com/leaanthony)
- 드래그 및 크기 조절 메커니즘의 안정성을 높이고 예상되는 플랫폼 동작에 더 가깝게 개선했습니다. [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- [#4097](https://github.com/wailsapp/wails/issues/4097) Webpack/angular가 런타임 초기화 코드를 제거하는 문제를 수정했습니다. [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- 처음에 숨겨진 메뉴 항목과 관련된 문제를 수정했습니다. [@IanVS](https://github.com/IanVS), [#4116](https://github.com/wailsapp/wails/pull/4116)
- 확장자가 없는 요청에서 `[request]`은 존재하지 않지만 `[request].html`은 존재할 때 assetFileServer가 `.html` 파일을 제공하지 않던 문제를 수정했습니다.
- 아이콘 생성 경로를 수정했습니다. [@robin-samuel](https://github.com/robin-samuel), [#4125](https://github.com/wailsapp/wails/pull/4125)
- `fullscreen`, `unfullscreen`, `unminimise` 및 `unmaximise` 이벤트가 발생하지 않던 문제를 수정했습니다. [@oSethoum](https://github.com/osethoum), [#4130](https://github.com/wailsapp/wails/pull/4130)
- 구성의 기본 버전에 잘못된 접두사가 지정되어 발생하던 NSIS 오류를 수정했습니다. [@robin-samuel](https://github.com/robin-samuel), [#4126](https://github.com/wailsapp/wails/pull/4126)
- Windows에서 Dialogs 런타임 함수가 이스케이프된 경로를 반환하던 문제를 수정했습니다. [TheGB0077](https://github.com/TheGB0077), [#4188](https://github.com/wailsapp/wails/pull/4188)
- HKCU의 Webview2 감지 경로를 수정했습니다. [@leaanthony](https://github.com/leaanthony).
- macOS의 입력 문제를 수정했습니다. [@leaanthony](https://github.com/leaanthony).
- Windows 아이콘 생성 작업 파일 이름을 수정했습니다. [@yulesxoxo](https://github.com/yulesxoxo), [#4219](https://github.com/wailsapp/wails/pull/4219).
- @kron의 작업을 바탕으로 프레임 없는 창의 투명도 문제를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- @kron의 작업을 바탕으로 창이 비활성화되거나 최소화되었을 때의 포커스 호출을 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- @kron의 작업을 바탕으로 작업 표시줄을 다시 시작한 후 시스템 트레이가 표시되지 않던 문제를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- fallbackResponseWriter가 Flush()를 구현하지 않던 문제를 수정했습니다. [#4245](https://github.com/wailsapp/wails/pull/4245)
- fallbackResponseWriter가 Flush()를 구현하지 않던 문제를 수정했습니다. [@superDingda], [#4236](https://github.com/wailsapp/wails/issues/4236)
- 비동기 Go 바인딩 함수 호출이 대기 중일 때 macOS 창을 닫으면 충돌하던 문제를 수정했습니다. [@joshhardy](https://github.com/joshhardy), [#4354](https://github.com/wailsapp/wails/pull/4354)
- Windows 효율성 모드 시작 시 발생하던 경쟁 상태를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- Windows 아이콘 핸들 정리를 수정했습니다. [@leaanthony](https://github.com/leaanthony).
- Windows에서 `OpenFileManager` 문제를 수정했습니다. [@PPTGamer](https://github.com/PPTGamer), [#4375](https://github.com/wailsapp/wails/pull/4375).
- Linux의 최소/최대 너비 옵션을 수정했습니다. @atterpac, [#3979](https://github.com/wailsapp/wails/pull/3979)
- npm 버전을 올려 Typescript 템플릿의 타입 정의를 수정했습니다. @atterpac, [#3966](https://github.com/wailsapp/wails/pull/3966)
- Sveltekit 템플릿의 CSS 참조를 수정했습니다. @atterpac, [#3945](https://github.com/wailsapp/wails/pull/3945)
- window run()의 주요 콜백이 메인 스레드에서 호출되도록 했습니다. [@leaanthony](https://github.com/leaanthony)
- 대화 상자의 디렉터리 선택기 예제를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- index.html이 없을 때 표시할 새로운 중국어 오류 페이지를 만들었습니다. [@leaanthony](https://github.com/leaanthony)
-  `windowDidBecomeKey` 콜백이 메인 스레드에서 실행되도록 했습니다. [@leaanthony](https://github.com/leaanthony)
-  프레임 없는 창에서 전체 화면을 지원합니다. [@leaanthony](https://github.com/leaanthony)
-  창 제거 로직을 개선했습니다. [@leaanthony](https://github.com/leaanthony)
-  시스템 트레이에 연결되었을 때의 창 위치 로직을 수정했습니다. [@leaanthony](https://github.com/leaanthony)
-  프레임 없는 창에서 전체 화면을 지원합니다. [@leaanthony](https://github.com/leaanthony)
- 이벤트 처리를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- 창 종료 로직을 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- 이제 공통 taskfile은 기본적으로 Typescript 템플릿의 Typescript 바인딩을 생성합니다. [@leaanthony](https://github.com/leaanthony)
- 열린 창이 없거나 시스템 트레이만 사용하는 경우 WM_CLOSE 메시지에서 애플리케이션이 종료되도록 수정했습니다. [@mmalcek](https://github.com/mmalcek), [#3990](https://github.com/wailsapp/wails/pull/3990)
- garble 빌드를 수정했습니다. @5aaee9, [#3192](https://github.com/wailsapp/wails/pull/3192)
- Windows nsis 빌드를 수정했습니다. [@leaanthony](https://github.com/leaanthony)
- Linux의 다중 선택 대화 상자에서 발생하던 교착 상태 수정. 원인: 닫히지 않은
- Windows 빌드 중 .syso 파일을 플랫폼 간에 정리하는 동작을 수정했습니다.
- amd64 appimage 컴파일을 수정했습니다. @atterpac
- 빌드 자산 업데이트를 수정했습니다. @ansxuman
- @atterpac이 Linux 시스템 트레이의 `OnClick` 및 `OnRightClick` 구현을 수정
- Mac에서 `AlwaysOnTop`이 작동하지 않는 문제를 수정한 사람:
-  `application.NewEditMenu`이 중복 항목을 포함하는 문제 수정
- 🐧 aarch64 컴파일 문제 수정
- ⊞ 라디오 그룹 메뉴 항목을 수정한 사람:
- MacOS에서 실행 가능한 .app을 빌드할 때 발생하는 오류 수정. 발생 조건: 'name'과 'outputfilename'이
- 드래그 앤 드롭 예제에서 customEventProcessor 사용 시 발생하는 버그를 수정한 사람:
- 🐧 IgnoreMouseEvents 추가로 인해 발생한 Linux 컴파일 오류를 수정한 사람:
- ⊞ syso 아이콘 파일 생성 버그를 수정한 사람:
- 🐧 Wayland에서 네이티브로 실행되도록 하는 수정 사항을 다음에서 반영:
- 다음에서 내부 서비스 메서드를 바인딩하지 않도록 수정:
- ⊞ 다음에서 시스템 트레이 시작 시 발생하는 패닉 수정:
- 다음에서 내부 서비스 메서드를 바인딩하지 않도록 수정:
- ⊞ 다음에서 시스템 트레이 시작 시 발생하는 패닉 수정:
- 메뉴 항목과 이벤트 처리를 대대적으로 리팩터링했습니다. 현재는 주로 macOS가 개선되었습니다. 작업자:
- 다음에서 플러그인 및 이벤트 리팩터링 후 테스트 수정:
- ⊞ `Failed to unregister class Chrome_WidgetWin_0` 경고를 수정했습니다. 작업자:
- 모듈 문제
- 다음에서 [atterpac](https://github.com/atterpac)이 크기 조정 이벤트 메시징을 수정:
- 🐧 NixOS에서 발생하는 테마 처리 오류를 수정한 사람:
- Windows에서 볼륨이 서로 다른 경우의 프로젝트 설치를 수정한 사람:
- 바닥글이 표시되도록 React 템플릿 CSS를 수정한 사람:
- refresh를 최신 버전으로 업데이트하여 개발 모드에서 작업할 때 발생하는 좀비 프로세스 문제 수정
- [Atterpac](https://github.com/atterpac)이 AppImage의 WebKit 파일 소싱 문제를 수정
- 다음에서 [Atterpac](https://github.com/Atterpac)이 Doctor의 apt 패키지 검증 문제를 수정:
- 다음에서 @5aaee9이 종료 시 애플리케이션이 멈추는 문제(Darwin)를 수정:
- Windows에서 예제의 배경색을 수정한 사람:
- 다음에서 [mmghv](https://github.com/mmghv)이 기본 컨텍스트 메뉴를 수정:
- Darwin에서 화살표 키의 16진수 값을 수정한 사람:
- Windows에서 드래그 앤 드롭이 작동하도록 수정했습니다. 추가한 사람:
- 사용자에게 적절한 드라이버가 없을 때 Linux의 Doctor에서 발생하는 버그 수정
- 시작 시 DPI 배율 조정 문제(Windows)를 수정했습니다. 다음에서 [@almas-x](https://github.com/almas-x)이 변경:
- 상대 경로를 사용하도록 `go.mod`의 치환 줄을 수정했습니다. 다음이 포함된 Windows 경로 문제를 해결:
- 연결된 창이 없을 때 MacOS 시스템 트레이 클릭 처리 문제를 수정한 사람:
- 알 수 없는 옵션으로 인해 Windows 빌드가 실패하는 문제를 수정한 사람:
- 다음이 없을 때 Windows에서 시스템 트레이 아이콘을 왼쪽 클릭하면 발생하는 충돌 수정:
- 창을 두 번 열 때 baseURL이 잘못되는 문제를 @5aaee9이 PR에서 수정:
- `WebviewWindow.Restore` 메서드에서 if 분기 순서 문제를 수정한 사람:
- 다음 조건에서 여러 차례의 `GetStartURL` 호출에 걸쳐 `startURL`을 올바르게 계산하도록 수정:
- `Screen` 구조체의 JS 타입이 Go의 대응 타입과 일치하도록 수정한 사람:
- 등록된 이벤트가 올바르게 정리되도록 `WML.Reload` 메서드 수정
- Linux에서 사용자 지정 컨텍스트 메뉴가 즉시 닫히는 문제를 수정한 사람:
- 바인딩에서 생성하는 모델 파일의 출력 경로와 확장자 수정
- 바인딩에서 생성하는 JS 코드 내 모델 파일의 import 경로 수정
- 일부 Linux 배포판에서 드래그 앤 드롭이 작동하지 않는 문제를 수정한 사람:
- `wails3 task dev` 사용 시 macOS용 태스크가 누락되는 문제를 수정한 사람:
- 이벤트 등록 시 nil 맵에 값을 대입하여 발생하는 오류를 수정한 사람:
- 바인딩된 메서드 매개변수의 언마샬링 문제를 수정한 사람:
- 바인딩된 메서드의 여러 반환값 처리 문제를 수정한 사람:
- 시스템 패키지 관리자로 설치하지 않은 npm을 Doctor가 감지하지 못하는 문제 수정
- 누락된 MicrosoftEdgeWebview2Setup.exe 문제 수정. 도움을 주신 분:
- 창 ID 처리로 인해 Linux에서 무작위로 발생하는 충돌을 @leaanthony가 수정했습니다. 다음을 기반으로 함:
- Linux에서 systemTray.setIcon 호출 시 충돌하는 문제를 수정한 사람:
- 다음 플랫폼의 `setFrameless` 함수에서 첫 호출 시 창 프레임이 확실히 적용되도록 수정:

### 변경됨

- **호환성을 깨는 변경 사항**: Go 맵 의미 체계를 정확히 반영하도록 생성된 JS/TS 바인딩의 맵 키를 이제 선택 사항으로 표시합니다. 이제 Typescript에서 맵 값에 접근하면 `T` 대신 `T | undefined`을 반환하므로 null 검사 또는 단언이 필요합니다(#4943). 작업자: `@fbbdev`
- `@wailsio/runtime`의 변경 사항에 따라 `Event` 사용을 `Events` 사용으로 변경하고, @AbdelhadiSeddar가 `Features/Events/Event System`의 문서에서 적절한 함수를 호출하도록 수정
- @leaanthony가 창별 옵션이었던 `EnabledFeatures`, `DisabledFeatures` 및 `AdditionalBrowserArgs`를 애플리케이션 수준의 `Options.Windows`로 이동(#4559)
- `Drag N Drop` 예제의 README를 업데이트하고, 이 예제에서 `Internal Drag and Drop`을 시연한다는 점을 강조함 — @ndianabasi
- 여러 디버그 로그의 수준을 Info에서 Debug로 변경(@mbaklor)
- **호환성을 깨는 변경 사항:** 창 옵션에서 `EnableDragAndDrop`의 이름을 `EnableFileDrop`로 변경
- **호환성을 깨는 변경 사항:** 이벤트 컨텍스트에서 `DropZoneDetails`의 이름을 `DropTargetDetails`로 변경
- **호환성을 깨는 변경 사항:** `WindowEventContext`의 `DropZoneDetails()` 메서드 이름을 `DropTargetDetails()`로 변경
- **호환성을 깨는 변경 사항:** `WindowDropZoneFilesDropped` 이벤트를 제거하고 대신 `WindowFilesDropped` 사용
- **호환성 중단:** HTML 속성을 `data-wails-dropzone`에서 `data-file-drop-target`(으)로 변경
- **호환성 중단:** CSS 호버 클래스를 `wails-dropzone-hover`에서 `file-drop-target-active`(으)로 변경
- **호환성 중단:** Windows에서 `DragEffect`, `OnEnterEffect`, `OnOverEffect` 옵션 제거(제거된 IDropTarget의 일부였음)
- 모든 런타임 JSON 처리(메서드 바인딩, 이벤트, webview 요청, 알림, kvstore)에 goccy/go-json을 사용하도록 전환하여 성능을 21-63% 향상하고 메모리 할당을 40-60% 절감
- 호출당 오버헤드를 줄이도록 BoundMethod 구조체 레이아웃을 최적화하고 isVariadic 플래그를 캐시
- 힙 할당을 방지하도록 인수가 `<=8`개인 메서드에 스택 할당 인수 버퍼 사용
- 반환 값이 하나인 경우 슬라이스 할당을 방지하도록 메서드 호출의 결과 수집 최적화
- 동시 실행 성능을 개선하도록 MIME 형식 캐시에 sync.Map 사용
- HTTP 전송 요청 본문을 읽을 때 버퍼 풀 사용
- 요청당 할당을 줄이도록 콘텐츠 유형 스니퍼에서 CloseNotify 채널을 지연 할당
- 애셋 서버에서 디버그 CSS 로깅 제거
- 50개 이상의 일반적인 웹 형식(글꼴, 오디오, 비디오 등)을 포함하도록 MIME 형식 확장자 맵 확대
- Window `X/Y` 옵션 문서 업데이트 @ruhuang2001
- 프런트엔드 바인딩 생성 옵션을 더 추가하여 `Frontend Runtime` 문서 업데이트, 작성자: @ndianabasi
- Wails v3 Asset Server 문서 페이지 업데이트, 작성자: @ndianabasi
- **호환성 중단**: 패키지 수준 대화 상자 함수(`application.InfoDialog()`, `application.QuestionDialog()` 등) 제거. 대신 `app.Dialog` 관리자를 사용: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- 실제 API와 일치하도록 대화 상자 문서 업데이트: `app.Dialog.*` 사용, `AddButton()`을 콜백과 함께 사용(`SetButtons()` 아님), `SetDefaultButton(*Button)` 사용(문자열 아님), `AddFilter()` 사용(`SetFilters()` 아님), `SetFilename()` 사용(`SetDefaultFilename()` 아님), 폴더 선택에는 `app.Dialog.OpenFile().CanChooseDirectories(true)` 사용
- **호환성 중단**: 이제 프로덕션 빌드가 기본값입니다. 개발 빌드를 만들려면 Taskfile에서 `DEV=true`을 설정하십시오. 설정 예제를 보려면 새 프로젝트를 생성하십시오. 작성자: @leaanthony
- 데이터 인수가 없거나 하나인 사용자 지정 이벤트를 내보낼 때 데이터 값은 슬라이스로 래핑되지 않고 Data 필드에 직접 할당됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4633](https://github.com/wailsapp/wails/pull/4633)
- 이제 Windows 트레이는 `NIS_HIDDEN`을 전환하여 `SystemTray.Show()`/`Hide()`을 준수하므로 앱이 실제로 사라졌다가 다시 나타날 수 있음(#4653).
- 트레이 등록 시 확인된 아이콘을 재사용하고 `NOTIFYICON_VERSION_4`을 한 번 설정하며 `NIF_SHOWTIP`을 활성화하므로 Explorer가 다시 시작된 후에도 도구 설명이 복구됨(#4653).
- macOS: 메뉴 막대와 Dock 영역을 제외하고 창을 가운데에 배치하도록 `frame` 대신 `visibleFrame` 사용
- macOS: 메뉴 막대와 Dock 영역을 제외하고 창을 가운데에 배치하도록 `frame` 대신 `visibleFrame` 사용
- `-config` 매개변수와 함께 `wails3 update build-assets`을 실행하면 `-product*` 매개변수를 통해 설정한 값은
- `window.NativeWindowHandle()` → `window.NativeWindow()`, 작성자: @leaanthony, [#4471](https://github.com/wailsapp/wails/pull/4471)
- 내부 창 처리 리팩터링, 작성자: @leaanthony, [#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` 및 `application.WindowNameKey` 제거(`application.WindowKey`(으)로 대체), 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 ContextMenuData가 any 대신 string을 반환함, 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 JS/TS 바인딩에서 고정 길이 배열 형식의 클래스 필드는 빈 상태가 아니라 예상 길이로 초기화됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4001](https://github.com/wailsapp/wails/pull/4001)
- 이제 ContextMenuData가 any 대신 string을 반환함, 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 `application.NewService`은 options를 선택적 매개변수로 받지 않음(대신 `application.NewServiceWithOptions` 사용), 작성자: [@leaanthony](https://github.com/leaanthony), [#4024](https://github.com/wailsapp/wails/pull/4024)
- `nanoid` 종속성 제거, 작성자: [@leaanthony](https://github.com/leaanthony)
- mica/acrylic/tabbed 창 스타일을 위한 Window 예제 업데이트, 작성자: [@leaanthony](https://github.com/leaanthony)
- JS/TS 바인딩에서 `internal.js/ts` 모델 파일이 제거되었으며, 이제 모든 모델은 `models.js/ts`에서 찾을 수 있음, 작성자: [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS 바인딩에서 명명된 형식은 더 이상 다른 명명된 형식의 별칭으로 렌더링되지 않으며, 이제 기존 동작은 별칭에만 적용됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS 바인딩의 클래스 모드에서 형식 매개변수를 형식으로 사용하는 구조체 필드는 선택 사항으로 표시되며 자동으로 초기화되지 않음, 작성자: [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- 템플릿에서 ESLint 제거, 작성자: [@IanVS](https://github.com/IanVS), [#4059](https://github.com/wailsapp/wails/pull/4059)
- 저작권 연도를 2025(으)로 업데이트, 작성자: [@IanVS](https://github.com/IanVS), [#4037](https://github.com/wailsapp/wails/pull/4037)
- event.Sender 문서 추가, 작성자: [@IanVS](https://github.com/IanVS), [#4075](https://github.com/wailsapp/wails/pull/4075)
- Go 1.24 지원, 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 `ServiceStartup` 훅은 `application.New`에서가 아니라 `App.Run`이 호출될 때 실행됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- 이제 `ServiceStartup` 오류는 프로세스를 종료하는 대신 `App.Run`에서 반환됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- 이제 JS에서 수행한 바인딩 및 대화 상자 호출은 문자열 대신 오류 객체로 거부됨, 작성자: [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Windows에서 시스템 트레이 메뉴 위치 지정 개선, 작성자: [@leaanthony](https://github.com/leaanthony)
- JS 런타임을 TypeScript로 포팅, 작성자: [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- 런타임을 가져오는 즉시 초기화되므로 창이 로드될 때까지 기다릴 필요가 없습니다. [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- 런타임은 더 이상 init 메서드를 내보내지 않습니다. 부수 효과 가져오기를 사용하여 초기화할 수 있습니다. [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- 이제 바인딩된 메서드는 취소될 경우 `CancelError`와 함께 거부되는 `CancellablePromise`을 반환합니다. 실제 호출 결과는 폐기됩니다. [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- 이제 기본 제공 서비스 형식의 이름을 일관되게 `Service`으로 사용합니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- 이제 옵션이 있는 기본 제공 서비스 생성 함수의 이름을 일관되게 `NewWithConfig`으로 사용합니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Go API와 일관성을 유지하기 위해 `sqlite` 서비스의 `Select` 메서드 이름을 `Query`으로 변경했습니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- 템플릿: 런타임을 "dependencies"로 이동하고 package.json 파일을 정리했습니다. [@IanVS](https://github.com/IanVS), [#4133](https://github.com/wailsapp/wails/pull/4133)
- 특정 macOS API를 사용할 수 있도록 개발 환경에서 앱 번들을 생성하고 임시 서명합니다. [@popaprozac](https://github.com/popaprozac), [#4171](https://github.com/wailsapp/wails/pull/4171)
- 빌드 자산을 플랫폼별 디렉터리로 이동했습니다. [@leaanthony](https://github.com/leaanthony)
- Taskfile을 플랫폼별 디렉터리로 이동하고 이름을 변경했습니다. [@leaanthony](https://github.com/leaanthony)
- `index.html`이 없을 때의 사용 환경을 크게 개선했습니다. [@leaanthony](https://github.com/leaanthony)
- [Windows] 최소화 및 복원 성능을 개선했습니다. [@leaanthony](https://github.com/leaanthony). [562589540](https://github.com/562589540)의 원본 [PR](https://github.com/wailsapp/wails/pull/3955)을 기반으로 합니다.
- `ShouldClose` 옵션을 제거했습니다(대신 events.Common.WindowClosing 이벤트에 대한 훅을 등록하세요). [@leaanthony](https://github.com/leaanthony)
- [Windows] 창을 열 때 발생하는 깜박임을 줄였습니다. [@leaanthony](https://github.com/leaanthony)
- 내부 함수로 사용하려던 `Window.Destroy`을 제거했습니다. [@leaanthony](https://github.com/leaanthony)
- `WindowClose` 이벤트의 이름을 `WindowClosing`으로 변경했습니다. [@leaanthony](https://github.com/leaanthony)
- 이제 프런트엔드 빌드에서는 빌드 유형에 따라 vite 환경 "development" 또는 "production"을 사용합니다. [@leaanthony](https://github.com/leaanthony)
- go-webview2 v1.19으로 업데이트했습니다. [@leaanthony](https://github.com/leaanthony)
- taskfile 포크를 사용하도록 보장했습니다. @leaanthony
- 다음을 사용하여 설치할 때 발생하는 버전 문제를 해결하도록 Taskfile 포크 업데이트
- 다음을 사용하여 설치할 때 발생하는 버전 문제를 해결하기 위해 Taskfile 포크 사용
- 이제 `service.OnStartup`은 오류 발생 시 애플리케이션을 종료하고 다음을 실행합니다
- 사용자 상호작용에 더 잘 부합하도록 시스템 트레이 클릭 메시징을 리팩터링했습니다. 작성자:
- 다음을 생성하는 프레임워크를 지원하도록 자산 임베드에 `all:frontend/dist` 포함
- Taskfile 리팩터링: [leaanthony](https://github.com/leaanthony), 다음에서:
- `go-webview2` v1.0.16으로 업그레이드. 작성자:
- `Screen` 형식이 `Id`이 아닌 `ID`을 포함하도록 수정했습니다. 작성자:
- `application.ServiceOptions`을 지원하도록 `go.mod.tmpl` wails 버전 업데이트. 작성자:
- 서비스 이름 결정 방식을 수정했습니다. [windom](https://github.com/windom/), 다음에서:
- 이제 mkdocs serve는 docker를 사용합니다. [leaanthony](https://github.com/leaanthony)
- 개발 구성을 `config.yml`으로 통합했습니다. 작성자:
- 이제 시스템 트레이 대화 상자는 사용 가능한 경우 애플리케이션 아이콘을 기본값으로 사용합니다(Windows). 작성자:
- macOS의 GPU 및 메모리 보고를 개선했습니다. 작성자:
- `WebviewGpuIsDisabled` 및 `EnableFraudulentWebsiteWarnings` 제거
- Events API 변경: `On`/`Emit` -> 사용자 이벤트, `OnApplicationEvent` ->
- Linux의 Events API 수정: [TheGB0077](https://github.com/TheGB0077), 다음에서:
- [CI] Actions를 개선하고 포크에서도 Actions를 실행할 수 있도록 했으며
- `AbsolutePosition()`의 이름을 `Position()`으로 변경했습니다. 작성자:
- 다음을 위해 Linux WebKit 종속성을 webkitgtk2-4.0 대신 webkit2gtk-4.1으로 업데이트했습니다
- 이제 번들로 제공되는 JS 런타임 스크립트는 ESM 모듈입니다. 이를 가져오는 script 태그는
- `@wailsio/runtime` 패키지는 `window.wails`에 API를 게시하지 않습니다
- 이제 Window API 모듈 `@wailsio/runtime/src/window`은 이를 포함하는 항목을 노출합니다
- 현재 Go `WebviewWindow`과 일치하도록 JS window API를 업데이트했습니다
- 이제 바인딩 생성기는 기본적으로 ID별 호출을 사용합니다. `-id` CLI 옵션은
- 새로운 바인딩 코드 레이아웃: 이전에는 출력 파일이 폴더에 정리되어 있었습니다
- 구조체 필드 `application.Options.Bind`의 이름을 다음으로 변경했습니다
- 바인딩 서비스를 위한 새 구문: 이제 서비스 인스턴스를 다음으로 감싸야 합니다
- 비터미널 또는 CI 환경에서 스피너를 비활성화했습니다. 작성자:

### 제거됨

- **호환성 중단**: 창별 `WindowsWindow` 옵션에서 `EnabledFeatures`, `DisabledFeatures` 및 `AdditionalLaunchArgs`을 제거합니다. 대신 애플리케이션 수준의 `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` 및 `Options.Windows.AdditionalBrowserArgs`을 사용하세요. 이 플래그는 공유 WebView2 환경 전체에 적용됩니다(#4559). 작성자: @leaanthony
- Windows의 네이티브 `IDropTarget` 구현을 제거하고 JavaScript 기반 방식을 사용합니다(v2 동작과 동일).
- github.com/wailsapp/mimetype 종속성을 제거하고 확장된 확장자 맵과 표준 라이브러리의 http.DetectContentType을 사용하여 바이너리 크기를 약 1.2MB 줄입니다.
- Linux 파일 탐색기용 최소 .desktop 파일 파서를 구현하여 gopkg.in/ini.v1 종속성을 제거하고 약 45KB를 절약합니다.
- Go 1.21+ 표준 라이브러리의 slices 패키지와 최소한의 내부 헬퍼를 사용하여 런타임 코드에서 samber/lo를 제거하고 약 310KB를 절감
- Darwin URL 스킴 핸들러에서 디버그 printf 문 제거(#4834)
- **호환성이 깨지는 변경 사항**: `linux:WindowLoadChanged` 이벤트를 제거하고 WebView 로딩 완료를 감지할 때는 대신 `linux:WindowLoadFinished` 사용(#3896), 작성자: @leaanthony

### 호환성이 깨지는 변경 사항

- **Manager API 리팩터링**: 코드 구성과 검색 편의성을 개선하기 위해 애플리케이션 API를 평면 구조에서 체계적인 manager 구조로 재구성. [#4359](https://github.com/wailsapp/wails/pull/4359)에서 [@leaanthony](https://github.com/leaanthony)가 기여
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Service 메서드 이름 변경: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`. 작성자: [@leaanthony](https://github.com/leaanthony)
- `Path` 및 `Paths` 메서드를 `application` 패키지로 이동. 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 애플리케이션 메뉴는 macOS에서만 사용할 수 있음. 작성자: [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## 추가됨

## 수정됨

## v3.0.0-alpha.77 - 2026-04-18

## 수정됨

## v3.0.0-alpha.76 - 2026-04-17

## 수정됨

## v3.0.0-alpha.75 - 2026-04-16

## 수정됨

## v3.0.0-alpha.74 - 2026-03-01

## 추가됨

## 수정됨

## v3.0.0-alpha.73 - 2026-02-27

## 수정됨

## v3.0.0-alpha.72 - 2026-02-16

## 수정됨

## v3.0.0-alpha.71 - 2026-02-10

## 추가됨

## 수정됨

## v3.0.0-alpha.70 - 2026-02-09

## 추가됨

## 수정됨

## v3.0.0-alpha.69 - 2026-02-08

## 추가됨

## 수정됨

## v3.0.0-alpha.68 - 2026-02-07

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.67 - 2026-02-04

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.66 - 2026-02-03

## 추가됨

## 변경됨

## 수정됨

## 제거됨

## v3.0.0-alpha.65 - 2026-02-01

## 추가

## v3.0.0-alpha.64 - 2026-01-26

## 추가

## v3.0.0-alpha.63 - 2026-01-25

## 수정

## v3.0.0-alpha.62 - 2026-01-22

## 수정

## v3.0.0-alpha.61 - 2026-01-20

## 수정

## v3.0.0-alpha.60 - 2026-01-14

## 수정

## v3.0.0-alpha.59 - 2026-01-11

## 변경

## v3.0.0-alpha.58 - 2026-01-09

## 수정

## v3.0.0-alpha.57 - 2026-01-05

## 변경

## 수정

## v3.0.0-alpha.56 - 2026-01-04

## 추가

## 변경

## 수정

## 제거

## v3.0.0-alpha.55 - 2026-01-02

## 변경

## 수정

## 제거

## v3.0.0-alpha.54 - 2025-12-29

## 추가

## 수정

## 제거

## v3.0.0-alpha.53 - 2025-12-27

## 추가

## 수정

## v3.0.0-alpha.52 - 2025-12-26

## 수정

## v3.0.0-alpha.51 - 2025-12-23

## 수정

## v3.0.0-alpha.50 - 2025-12-21

## 변경

## v3.0.0-alpha.49 - 2025-12-18

## 변경

## v3.0.0-alpha.48 - 2025-12-16

## 추가

## 변경

## 수정

## v3.0.0-alpha.47 - 2025-12-15

## 추가

## 수정

## v3.0.0-alpha.46 - 2025-12-14

## 추가

## 제거

## v3.0.0-alpha.45 - 2025-12-13

## 추가

## 수정

## v3.0.0-alpha.44 - 2025-12-12

## 추가

## 변경

## 수정

## v3.0.0-alpha.43 - 2025-12-11

## 추가

## v3.0.0-alpha.42 - 2025-12-10

## 추가

## v3.0.0-alpha.41 - 2025-11-23

## 수정

## v3.0.0-alpha.40 - 2025-11-13

## 수정

## v3.0.0-alpha.39 - 2025-11-12

## 추가

## 변경

## v3.0.0-alpha.38 - 2025-11-04

## 추가

## 변경

## 수정

## v3.0.0-alpha.37 - 2025-11-02

## 수정

## v3.0.0-alpha.36 - 2025-10-15

## 수정

## v3.0.0-alpha.35 - 2025-10-14

## 수정

## v3.0.0-alpha.34 - 2025-10-06

## 추가

## 수정

## v3.0.0-alpha.33 - 2025-10-04

## 수정됨

## v3.0.0-alpha.32 - 2025-10-02

## 수정됨

## v3.0.0-alpha.31 - 2025-09-27

## 수정됨

## v3.0.0-alpha.30 - 2025-09-26

## 수정됨

## v3.0.0-alpha.29 - 2025-09-25

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.29 - 2025-09-25

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.27 - 2025-09-07

## 수정됨

## v3.0.0-alpha.26 - 2025-08-24

## 추가됨

## v3.0.0-alpha.25 - 2025-08-16

## 변경됨

더 이상 무시되지 않으며 구성 값을 재정의합니다.

## v3.0.0-alpha.24 - 2025-08-13

## 추가됨

## v3.0.0-alpha.23 - 2025-08-11

## 수정됨

## v3.0.0-alpha.22 - 2025-08-10

## 추가됨

## 변경됨

+ 지나치게 광범위한 Linux 패키지 종속성을 수정하고 오래된 RPM 종속성을 수정했습니다.

## v3.0.0-alpha.21 - 2025-08-07

## 수정됨

## v3.0.0-alpha.20 - 2025-08-06

## 수정됨

## v3.0.0-alpha.19 - 2025-08-05

## 추가됨

## 수정됨

## v3.0.0-alpha.18 - 2025-08-03

## 추가됨

## 수정됨

## v3.0.0-alpha.17 - 2025-07-31

## 수정됨

## v3.0.0-alpha.16 - 2025-07-25

## 추가됨

## v3.0.0-alpha.15 - 2025-07-25

## 추가됨

## v3.0.0-alpha.14 - 2025-07-25

## 추가됨

## v3.0.0-alpha.12 - 2025-07-15

### 추가됨

### 수정됨

## v3.0.0-alpha.11 - 2025-07-12

## 추가됨

## v3.0.0-alpha.10 - 2025-07-06

### 호환성을 깨는 변경 사항

### 추가됨

### 수정됨

### 변경됨

## v3.0.0-alpha.9 - 2025-01-13

### 추가됨

### 수정됨

### 변경됨

## v3.0.0-alpha.8.3 - 2024-12-07

### 변경됨

## v3.0.0-alpha.8.2 - 2024-12-07

### 변경됨

`go install`, 작성자: @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### 변경됨

`go install`, 작성자: @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### 추가됨

@atterpac이 [#3909](https://github.com/wailsapp/wails/3909)에서   [ansxuman](https://github.com/ansxuman)이   [#3902](https://github.com/wailsapp/wails/pull/3902)에서   [atterpac](https://github.com/atterpac)이   [#3867](https://github.com/wailsapp/wails/pull/3867)에서   [atterpac](https://github.com/atterpac)이   [#3829](https://github.com/wailsapp/wails/pull/3829)에서   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000)이   [#3856](https://github.com/wailsapp/wails/pull/3856)에서   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   지정된 X/Y 위치에 배치함.   [leaanthony](https://github.com/leaanthony)가   [#3885](https://github.com/wailsapp/wails/pull/3885)에서 구현   [ansxuman](https://github.com/ansxuman)과   [leaanthony](https://github.com/leaanthony)가   [#3823](https://github.com/wailsapp/wails/pull/3823)에서   [leaanthony](https://github.com/leaanthony)가   [#3766](https://github.com/wailsapp/wails/pull/3766)에서 구현   [leaanthony](https://github.com/leaanthony)가   [#3888](https://github.com/wailsapp/wails/pull/3888)에서 구현   [leaanthony](https://github.com/leaanthony). -

### 변경됨

이전에 시작된 모든 서비스에 대한 `service.OnShutdown` 관련 변경, @atterpac이   [#3920](https://github.com/wailsapp/wails/pull/3920)에서 변경   @atterpac이 [#3907](https://github.com/wailsapp/wails/pull/3907)에서 변경   하위 폴더 관련 변경, @atterpac,   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)이   [#3778](https://github.com/wailsapp/wails/pull/3778)에서 변경   [northes](https://github.com/northes)가   [#3836](https://github.com/wailsapp/wails/pull/3836)에서 변경   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (`EnabledFeatures` 및 `DisabledFeatures` 옵션으로 대체됨)   [leaanthony](https://github.com/leaanthony)가 변경

### 수정됨

@michael-freling이 채널 변수를   [#3925](https://github.com/wailsapp/wails/pull/3925)에서 수정   [ansxuman](https://github.com/ansxuman)이   [#3924](https://github.com/wailsapp/wails/pull/3924)에서 수정   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   [#3886](https://github.com/wailsapp/wails/pull/3886)에서 수정   [leaanthony](https://github.com/leaanthony)가   [#3841](https://github.com/wailsapp/wails/pull/3841)에서 수정   Darwin의 편집 메뉴에 있는 `PasteAndMatchStyle` 역할을   [johnmccabe](https://github.com/johnmccabe)가   [#3839](https://github.com/wailsapp/wails/pull/3839)에서 수정   [#3840](https://github.com/wailsapp/wails/issues/3840)을   [#3854](https://github.com/wailsapp/wails/pull/3854)에서   [kodflow](https://github.com/kodflow)가 수정   [@leaanthony](https://github.com/leaanthony)   서로 다릅니다. @nickisworking이   [#3789](https://github.com/wailsapp/wails/pull/3789)에서 수정

## v3.0.0-alpha.7 - 2024-09-18

### 추가됨

[mmghv](https://github.com/mmghv)가   [#3665](https://github.com/wailsapp/wails/pull/3665)에서 추가   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)과   [leaanthony](https://github.com/leaanthony)가   [#3570](https://github.com/wailsapp/wails/pull/3570)에서 추가

### 변경됨

애플리케이션 이벤트 `OnWindowEvent`를 창 이벤트로 변경.   [leaanthony](https://github.com/leaanthony)가 변경   [#3734](https://github.com/wailsapp/wails/pull/3734)   `v3/` 또는 `v3-` 접두사가 붙은 브랜치를   [stendler](https://github.com/stendler)가   [#3747](https://github.com/wailsapp/wails/pull/3747)에서 변경

### 수정됨

[etesam913](https://github.com/etesam913)이   [#3742](https://github.com/wailsapp/wails/pull/3742)에서 수정   [atterpac](https://github.com/atterpac)이   [#3721](https://github.com/wailsapp/wails/pull/3721)에서 수정   [atterpac](https://github.com/atterpac)이   [#3675](https://github.com/wailsapp/wails/pull/3675)에서 수정   [#1811](https://github.com/wailsapp/wails/pull/1811)을   [#3614](https://github.com/wailsapp/wails/pull/3614)에서   [@stendler](https://github.com/stendler)가 수정   [#3720](https://github.com/wailsapp/wails/pull/3720)에서   [leaanthony](https://github.com/leaanthony)가 수정   [#3693](https://github.com/wailsapp/wails/issues/3693)에서   [@DeltaLaboratory](https://github.com/DeltaLaboratory)가 수정   [#3720](https://github.com/wailsapp/wails/pull/3720)에서   [leaanthony](https://github.com/leaanthony)가 수정   [#3693](https://github.com/wailsapp/wails/issues/3693)에서   [@DeltaLaboratory](https://github.com/DeltaLaboratory)가 수정   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)에서   [@stendler](https://github.com/stendler)가 수정   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 수정됨

## v3.0.0-alpha.5 - 2024-07-30

### 추가됨

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   @5aaee9가 아이콘 클릭을 [#2991](https://github.com/wailsapp/wails/pull/2991)에서 추가   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev)가   [#3282](https://github.com/wailsapp/wails/pull/3282)에서 추가   @[Atterpac](https://github.com/Atterpac)이   [#3022](https://github.com/wailsapp/wails/pull/3022])에서 추가   [@marcus-crane](https://github.com/marcus-crane)이   [#3146](https://github.com/wailsapp/wails/pull/3146)에서 추가   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev)가   [#3281](https://github.com/wailsapp/wails/pull/3281)에서 추가   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   @Mai-Lapyst의 [PR](https://github.com/wailsapp/wails/pull/2044)을 기반으로 함   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 추가   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 추가   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 추가   npm 패키지를 [@fbbdev](https://github.com/fbbdev)가   [#3334](https://github.com/wailsapp/wails/pull/3334)에서 추가   [#3354](https://github.com/wailsapp/wails/pull/3354)에서 추가   `WAILS_VITE_PORT`을 [@abichinger](https://github.com/abichinger)가   [#3429](https://github.com/wailsapp/wails/pull/3429)에서 추가   [@abichinger](https://github.com/abichinger)가   [#3431](https://github.com/wailsapp/wails/pull/3431)에서 추가   [@bruxaodev](https://github.com/bruxaodev)가   [#3667](https://github.com/wailsapp/wails/pull/3667)에서 추가   [@OlegGulevskyy](https://github.com/OlegGulevskyy)가   [#3674](https://github.com/wailsapp/wails/pull/3674)에서 추가

### 수정됨

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane)이 다음에서 기여:   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac)이 다음에서 기여:   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac)이 다음에서 기여:   [#3477](https://github.com/wailsapp/wails/pull/3477)   [Atterpac](https://github.com/atterpac)이 다음에서 기여:   [#3320](https://github.com/wailsapp/wails/pull/3320).   다음에서: [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv)가 다음에서 기여:   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave)가 다음에서 기여:   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight)가 다음에서 기여:   [PR](https://github.com/wailsapp/wails/pull/3039)   설치됨. [@pylotlight](https://github.com/pylotlight)가 다음에서 추가:   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   공백 - @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal)이 PR에서 기여:   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal)이 PR에서 기여:   [#3208](https://github.com/wailsapp/wails/pull/3208)   연결된 창, [tw1nk](https://github.com/tw1nk)가 PR에서 기여:   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL`이 있습니다.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3295](https://github.com/wailsapp/wails/pull/3295)   리스너, [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger)가 다음에서 기여:   [#3330](https://github.com/wailsapp/wails/pull/3330)   생성기, [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3334](https://github.com/wailsapp/wails/pull/3334)   생성기, [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger)가 다음에서 기여:   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy)가 다음에서 기여:   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy)가 다음에서 기여:   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@pekim](https://github.com/pekim)이 다음에서 기여:   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   [@5aaee9](https://github.com/5aaee9)의 PR [#3466](https://github.com/wailsapp/wails/pull/3622).   [@windom](https://github.com/windom/)이 다음에서 기여:   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows, [@bruxaodev](https://github.com/bruxaodev/)가 다음에서 기여:   [#3691](https://github.com/wailsapp/wails/pull/3691).

### 변경됨

[mmghv](https://github.com/mmghv)가 다음에서 기여:   [#3611](https://github.com/wailsapp/wails/pull/3611)   Ubuntu 24.04 LTS 지원, [atterpac](https://github.com/atterpac)이 다음에서 기여:   [#3461](https://github.com/wailsapp/wails/pull/3461)   `type="module"` 속성이 있어야 합니다. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3295](https://github.com/wailsapp/wails/pull/3295)   객체이며 WML 시스템을 시작하지 않습니다. 캡슐화를 개선하기 위해 이렇게 변경했습니다. 원하는 경우 새 `WML.Enable` 메서드를 호출하여 WML 시스템을 수동으로 시작할 수 있습니다. 번들로 제공되는 JS 런타임 스크립트는 여전히 두 작업을 모두 자동으로 수행합니다. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3295](https://github.com/wailsapp/wails/pull/3295)   window 객체를 기본 내보내기로 제공합니다. 이제 ESM 명명된 가져오기 또는 네임스페이스 가져오기 구문을 통해 개별 메서드를 가져올 수 없습니다.   API. 일부 메서드의 이름이나 프로토타입이 변경되었습니다. 구체적으로 `Screen`은 `GetScreen`이 되고, `GetZoomLevel`/`SetZoomLevel`은 `GetZoom`/`SetZoom`이 됩니다. 이제 `GetZoom`, `Width` 및 `Height`은 값을 객체로 래핑하지 않고 직접 반환합니다. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3295](https://github.com/wailsapp/wails/pull/3295)   제거되었습니다. 이름으로 호출하는 방식으로 되돌리려면 `-names` CLI 옵션을 사용하세요. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3468](https://github.com/wailsapp/wails/pull/3468)   이전에는 포함된 패키지의 이름을 따랐지만, 이제 모듈 경로를 포함한 전체 Go 가져오기 경로를 사용합니다. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.NewService` 호출. [@fbbdev](https://github.com/fbbdev)가 다음에서 기여:   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory)가 다음에서 기여:   [#3574](https://github.com/wailsapp/wails/pull/3574)
