---
title: "API 참조"
description: "Wails 타입, 메서드, 옵션 및 CLI 명령 조회"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## 이 참조 문서 소개

애플리케이션을 개발하면서 Wails 타입, 메서드, 옵션 및 명령을 조회하려면 이 섹션을 사용하세요. 단계별 학습이 필요하면 [튜토리얼](/tutorials/overview/)부터 시작하세요. 작업을 완료하는 절차는 [방법 가이드](/quick-start/next-steps/#features)를 참조하세요. 설계 배경은 [핵심 개념](/concepts/architecture/)을 참조하세요.

**구성:**

- [애플리케이션](/reference/application/) - 핵심 애플리케이션 API
- [창](/reference/window/) - 창 생성 및 관리. [창 옵션](/features/windows/options/)에는 구성 필드가 나열되어 있습니다.
- [메뉴](/reference/menu/) - 애플리케이션, 컨텍스트 및 시스템 트레이 API. [메뉴 항목 참조](/features/menus/reference/)에서는 항목 타입과 속성을 설명합니다.
- [이벤트](/reference/events/) - 이벤트 시스템 및 기본 제공 이벤트
- [대화 상자](/reference/dialogs/) - 파일 및 메시지 대화 상자
- [프런트엔드 런타임](/reference/frontend-runtime/) - 프런트엔드 런타임 API
- [CLI 명령 개요](/reference/cli/) - 용도별로 분류된 명령. [플래그 및 예제](/guides/cli/)에서는 개별 명령을 자세히 설명합니다.
- [업데이트 매니페스트 프로토콜](/reference/update-manifest/) - 업데이트 메타데이터 및 검증 필드
- [모바일 API](/guides/mobile/mobile-api/) - 실험적 모바일 API

## API 규칙

@details{title="Go API 규칙 - Go를 처음 사용하는 개발자용"}
### 명명 규칙

- <strong>타입</strong>: PascalCase(예: `WebviewWindow`)
- <strong>메서드</strong>: PascalCase(예: `SetTitle()`)
- <strong>옵션</strong>: PascalCase 구조체(예: `WindowOptions`)
- <strong>상수</strong>: PascalCase(예: `WindowStartStateMaximised`)

#### 오류 처리

실패할 수 있는 대부분의 메서드는 마지막 반환 값으로 `error`를 반환합니다. `app.Run()`는 애플리케이션이 종료될 때까지 블로킹하며 시작 오류가 있으면 반환합니다.

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

창 생성은 오류를 반환하지 않습니다. `app.Window.New()`는 `*WebviewWindow`를 직접 반환합니다.

#### 컨텍스트

서비스 수명 주기 메서드는 `context.Context`를 받습니다.

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

애플리케이션 수명 주기 컨텍스트는 `app.Context()`를 통해 사용할 수 있습니다. `RunWithContext`는 없으므로 `app.Run()`를 호출하세요.

#### 옵션 패턴

구성에는 옵션 구조체를 사용합니다.

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### JavaScript API 규칙

#### 명명 규칙

- **함수**: camelCase(예: `setTitle()`)
- **상수**: SCREAMING<em>SNAKE</em>CASE(예: `WINDOW_EVENT_FOCUS`)

#### 기본 비동기 방식

모든 Go 메서드 호출은 Promise를 반환합니다.

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### 오류 처리

Go 오류는 JavaScript 예외로 변환됩니다.

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### 타입 안전성

TypeScript 정의는 자동으로 생성됩니다.

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## 패키지 구조

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## 가져오기 경로

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## 타입 참조

### 공통 타입

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## 플랫폼별 차이

일부 API는 플랫폼에 따라 다르게 동작합니다.

| 기능 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **애플리케이션 메뉴** | 창 메뉴 막대 | 전역 메뉴 막대 | 창 메뉴 막대 |
| **시스템 트레이** | 알림 영역 | 메뉴 막대 | 시스템 트레이 |
| **Dock** | 해당 없음 | ✅ 사용 가능 | 해당 없음 |
| **파일 대화 상자** | 네이티브 | 네이티브 | 네이티브(GTK) |
| **투명도** | ✅ 완전 지원 | [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) 필요 | ⚠️ 제한적 |

플랫폼별 동작은 각 API 섹션에 문서화되어 있습니다.

## 버전 관리

Wails v3는 시맨틱 버전 관리를 따릅니다.

- **메이저**(v4.0.0): 호환성을 깨는 변경 사항
- **마이너**(v3.1.0): 하위 호환되는 새로운 기능
- **패치**(v3.0.1): 하위 호환되는 버그 수정

**Wails v3.0.0:** 안정 버전

## 사용 중단 정책

API가 사용 중단될 때는 다음 정책을 적용합니다.

1. **문서에 표시**하고 사용 중단 알림 제공
2. **대체 수단 제공** 및 마이그레이션 가이드 제공
3. 제거하기 전에 **1개의 메이저 버전** 동안 유지
4. **컴파일러 경고**(가능한 경우)

## API 안정성

### 안정 API ✅

다음 API는 안정적이며 프로덕션 환경에서 안전하게 사용할 수 있습니다.

- 핵심 애플리케이션 API
- 창 관리
- 메뉴 시스템
- 이벤트 시스템
- 파일 대화 상자
- 서비스 바인딩

### 실험적 API ⚠️

실험적이라고 명시적으로 문서화된 API만 안정적인 데스크톱 호환성 계약의 적용 대상에서 제외됩니다:

- Android 지원
- iOS 지원
- 그 밖에 실험적이라고 명시된 API

실험적 API는 문서에 표시되어 있습니다.

## 도움말

### API 관련 질문

1. **이 참조 문서 확인** - 전체 API 문서
2. **예제 확인** - [GitHub 예제](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Discord 검색** - [Discord 서버](https://discord.gg/JDdSxwjhGf)
4. **커뮤니티에 질문** - Discord #help 채널

### API 문제 신고

버그나 일관되지 않은 동작을 발견하셨나요?

1. **기존 이슈 확인** - [GitHub 이슈](https://github.com/wailsapp/wails/issues)
2. **상세 보고서 작성** - 코드, 오류, 플랫폼을 포함하세요
3. **재현 방법 제공** - 문제를 보여 주는 최소 예제를 제공하세요

## 관련 문서

- [튜토리얼](/tutorials/overview/) - 실제 애플리케이션을 만들며 배워 보세요
- [가이드](/guides/architecture/) - 일반적인 시나리오를 위한 작업 중심 가이드
- [기능](/features/windows/basics/) - 기능별 문서
- [예제](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) - GitHub에서 제공하는 실행 가능한 코드 예제

---

**API 살펴보기:** 왼쪽 탐색 메뉴에서 각 API를 살펴보세요.
