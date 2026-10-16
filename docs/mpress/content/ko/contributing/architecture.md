---
title: "Wails v3 아키텍처"
description: "Wails v3 내부의 모든 구성 요소를 자세히 살펴보는 다이어그램과 설명"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3는 Go 런타임, JavaScript 브리지, 작업 기반 도구 체인과 최신 웹 기술로 구동되는 네이티브 애플리케이션을 출시할 수 있게 해 주는 템플릿 모음으로 구성된 <strong>풀스택 데스크톱 프레임워크</strong>입니다.

이 페이지에서는 다음 네 가지 다이어그램으로 <em>전체 구조</em>를 보여 줍니다:

1. **전체 아키텍처** – 모든 하위 시스템이 연결되는 방식\
2. **런타임 흐름** – JS가 Go를 호출하거나 Go가 JS를 호출할 때 일어나는 과정\
3. **개발 환경과 프로덕션 환경** – 애셋 서버의 두 가지 모드\
4. **플랫폼별 구현** – OS별 코드가 위치하는 곳\

---

## 1 · 전체 아키텍처

**Wails v3 – 상위 수준 스택**

**[상위 수준 스택 다이어그램 자리표시자]**

---

## 2 · 런타임 호출 흐름

**런타임 – JavaScript ⇄ Go 호출 경로**

**[런타임 호출 흐름 다이어그램 자리표시자]**

핵심 사항:

- **HTTP / IPC 없음** – 브리지는 네이티브 WebView의 인메모리 채널을 사용합니다\
- **메서드 ID** – 결정적 FNV 해시를 사용하여 Go에서 O(1) 조회가 가능합니다\
- **Promise** – 오류는 스택과 코드를 포함한 거부로 전파됩니다

---

## 3 · 개발 환경과 프로덕션 환경의 애셋 흐름

**개발 ↔ 프로덕션 애셋 서버**

**[애셋 흐름 다이어그램 자리표시자]**

- **개발** 환경에서는 서버가 알 수 없는 경로를 프레임워크의 실시간 리로드 서버로 프록시하고 디스크에서 정적 애셋을 제공합니다.
- **프로덕션** 환경에서는 동일한 API가 `go:embed`의 지원을 받아 의존성이 전혀 없는 바이너리를 생성합니다.

---

## 4 · 플랫폼별 런타임 분리

**OS별 런타임 파일**

**[플랫폼 분리 다이어그램 자리표시자]**

모든 기능은 다음 패턴을 따릅니다:

1. `pkg/application`의 **공통 인터페이스**\
2. `pkg/application/messageprocessor_*.go`의 **메시지 프로세서** 진입점\
3. `pkg/application/*_{darwin,linux,windows}.go`의 **OS별 구현**(예: `webview_window_darwin.go`, `clipboard_linux.go`, `dialogs_windows.go`, `systemtray_*.go`, `mainthread_*.go`)은 빌드 태그로 보호됩니다. Linux에는 `linux_cgo.go` / `linux_cgo_gtk4.{go,c,h}`의 cgo 브리지도 있습니다.

`internal/runtime/`에는 소규모 `runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` 빌드 태그 연결 코드와 `internal/runtime/desktop/` 아래에 포함된 JS 런타임만 있습니다.

`internal/capabilities/`는 플랫폼별 기능 집합을 선언하기 위해 존재하지만 `ErrCapability` 센티널은 없습니다. 기능 게이팅은 일반 빌드 태그와 플랫폼별 스텁 반환(예: `nil` 또는 기능별 오류)을 통해 수행됩니다.

---

## 요약

이 다이어그램들은 **코드가 있는 위치**, **데이터가 이동하는 방식**, 그리고 <strong>각 계층이 담당하는 책임</strong>을 개괄적으로 보여 줍니다. 이어지는 상세 페이지를 살펴볼 때 가까이 두고 참고하세요. Wails v3 소스 트리를 안내하는 지도 역할을 합니다.
