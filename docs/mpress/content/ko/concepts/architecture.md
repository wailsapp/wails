---
title: "Wails의 작동 방식"
description: "Wails 아키텍처와 네이티브 성능을 구현하는 방식 이해하기"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails는 **백엔드에 Go를 사용하고** **프런트엔드에 웹 기술을 사용하여** 데스크톱 애플리케이션을 빌드하는 프레임워크입니다. 하지만 Electron과 달리 Wails는 브라우저를 번들로 포함하지 않고 **운영 체제의 네이티브 WebView**를 사용합니다.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Wails 앱

  frontend: 프런트엔드
  backend: Go 백엔드
  os: 운영 체제

  Initialisation: 초기화 {
    shape: sequence_diagram
    backend."Serves Static Web App": 정적 웹 앱 제공
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": OS 네이티브 WebView를 통해 사이트 렌더링
  }
  Regular Communication: 일반 통신 {
    shape: sequence_diagram
    frontend."Make API-style call": API 방식 호출
    frontend -> backend.a: JSON
    backend.a."Service processes request": 서비스가 요청 처리
    backend.a -> os: 시스템 API 호출
    backend.a."Generate Response": 응답 생성
    backend.a -> frontend: JSON
    frontend."Process response": 응답 처리
  }
  backend.a.label: a
}
```

**Electron과의 주요 차이점:**

| 항목 | Wails | Electron |
| --- | --- | --- |
| **브라우저** | OS에서 제공하는 WebView | 번들로 포함된 Chromium(~100MB) |
| **백엔드** | Go(컴파일 방식) | Node.js(인터프리트 방식) |
| **통신** | 인메모리 브리지 | IPC(프로세스 간 통신) |
| **번들 크기** | ~15MB | ~150MB |
| **메모리** | ~10MB | ~100MB+ |
| **시작 시간** | &lt;0.5초 | 2-3초 |

## 핵심 구성 요소

### 1. 네이티브 WebView

Wails는 운영 체제에 내장된 웹 렌더링 엔진을 사용합니다.

@tabs{sync-key="platform"}
[Windows]
**WebView2**(Microsoft Edge WebView2)

- Chromium 기반(Edge 브라우저와 동일)
- Windows 10/11에 사전 설치됨
- Windows Update를 통한 자동 업데이트
- 최신 웹 표준 완벽 지원

[macOS]
**WebKit**(Safari의 렌더링 엔진)

- macOS에 내장됨
- Safari 브라우저와 동일한 엔진
- 뛰어난 성능과 배터리 사용 시간
- 최신 웹 표준 완벽 지원

[Linux]
**WebKitGTK**(WebKit의 GTK 포트)

- 패키지 관리자를 통해 설치
- GNOME Web(Epiphany)과 동일한 엔진
- 우수한 표준 지원
- 가볍고 뛰어난 성능

@end

**이것이 중요한 이유:**

- **브라우저를 번들로 포함하지 않음** → 더 작은 앱 크기
- **OS 네이티브** → 더 나은 통합과 성능
- **자동 업데이트** → OS 업데이트를 통해 보안 패치 적용
- **익숙한 렌더링** → 시스템 브라우저와 동일

### 2. Wails 브리지

브리지는 Wails의 핵심으로, Go와 JavaScript 간의 **직접 통신**을 지원합니다.

```d2
direction: down

Frontend: 프런트엔드(JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Wails 브리지 {
  Encoder: JSON 인코더 {
    shape: rectangle
  }

  Router: 메서드 라우터 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 디코더 {
    shape: rectangle
  }
}

Backend: 백엔드(Go) {
  Services: 등록된 서비스 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Go 메서드 호출\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. JSON으로 인코딩\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. 서비스로 라우팅\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. 결과 반환\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. JS로 디코딩\nPromise 이행"
```

**작동 방식:**

1. **프런트엔드가 Go 메서드를 호출합니다**(자동 생성된 바인딩을 통해).
2. **브리지가 호출을 JSON으로 인코딩합니다**(메서드 이름 + 인수).
3. **라우터가 등록된 서비스에서 Go 메서드를 찾습니다**.
4. **Go 메서드가 실행되고** 값을 반환합니다.
5. **브리지가 결과를 디코딩하고** 프런트엔드로 다시 전송합니다.
6. 결과와 함께 JavaScript의 **Promise가 이행됩니다**.

**성능 특성:**

- **인메모리**: 네트워크 오버헤드와 HTTP가 없습니다.
- **제로 카피**: 가능한 경우 적용합니다(대용량 데이터).
- **기본적으로 비동기**: 양쪽 모두에서 블로킹되지 않습니다.
- **타입 안전성**: TypeScript 정의가 자동으로 생성됩니다.

### 3. 서비스 시스템

서비스는 Go 기능을 프런트엔드에 노출하는 데 권장되는 방식입니다.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**서비스 검색:**

- Wails는 시작할 때 **구조체를 스캔합니다**
- **내보낸 메서드**를 프런트엔드에서 호출할 수 있게 됩니다
- TypeScript 바인딩을 위해 **타입 정보**를 추출합니다
- **오류 처리**는 자동으로 이루어집니다(Go 오류 → JS 예외)

**생성된 TypeScript 바인딩:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**서비스를 사용하는 이유**

- **타입 안전성**: 완전한 TypeScript 지원
- **자동 검색**: 메서드를 수동으로 등록할 필요가 없음
- **체계적인 구성**: 관련 기능을 하나로 묶음
- **테스트 가능**: 서비스는 일반 Go 구조체일 뿐임

[서비스에 대해 자세히 알아보기 →](/features/bindings/services/)

### 4. 이벤트 시스템

이벤트를 사용하면 컴포넌트 간에 **발행/구독 통신**을 할 수 있습니다.

```d2
direction: left

Wails Event System: Wails 이벤트 시스템 {
  shape: sequence_diagram

  window1: 창 1
  window2: 창 2
  backend: Go 백엔드

  Event Driver: 이벤트 드라이버 {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "'data-updated' 이벤트 구독"
    window2."Subscribe to 'data-updated' events": "'data-updated' 이벤트 구독"
    backend.a."App Emit('data-updated', data)": "앱 Emit('data-updated', data)"
    backend.a -> window1.a: JSON 이벤트 버스
    backend.a -> window2: JSON 이벤트 버스
    window1.a."Subscriber processes On('data-updated', handler)": "구독자가 On('data-updated', handler) 처리"
    window2."Subscriber processes On('data-updated', handler)": "구독자가 On('data-updated', handler) 처리"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**사용 사례:**

- **창 간 통신**: 한 창이 다른 창에 알림
- **백그라운드 작업**: Go 서비스가 UI에 진행 상황을 알림
- **상태 동기화**: 여러 창의 상태를 동기화
- **느슨한 결합**: 컴포넌트에 직접 참조가 필요하지 않음

**예제:**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[이벤트에 대해 자세히 알아보기 →](/features/events/system/)

## 애플리케이션 수명 주기

수명 주기를 이해하면 리소스를 초기화하고 정리해야 하는 시점을 파악할 수 있습니다.

```d2
direction: down

Start: 애플리케이션 시작 {
  shape: oval
  style.fill: "#10B981"
}

Init: 초기화 {
  Create: 애플리케이션 생성 {
    shape: rectangle
  }

  Register: 서비스 등록 {
    shape: rectangle
  }

  Setup: 창/메뉴 설정 {
    shape: rectangle
  }
}

Run: 이벤트 루프 {
  Events: 이벤트 처리 {
    shape: rectangle
  }

  Messages: 메시지 처리 {
    shape: rectangle
  }

  Render: UI 업데이트 {
    shape: rectangle
  }
}

Shutdown: 종료 {
  Cleanup: 리소스 정리 {
    shape: rectangle
  }

  Save: 상태 저장 {
    shape: rectangle
  }
}

End: 애플리케이션 종료 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: 반복
Run.Events -> Shutdown.Cleanup: 종료 신호
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**수명 주기 훅:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options`에는 `OnStartup` 필드가 없습니다. 시작 작업은 서비스의 `ServiceStartup(ctx, options)`, `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`을 통해 등록한 콜백 또는 단순히 `app.Run()` 이전에 수행해야 합니다.

[수명 주기에 대해 자세히 알아보기 →](/concepts/lifecycle/)

## 빌드 프로세스

Wails가 애플리케이션을 빌드하는 방식을 알아봅니다.

```d2
direction: down

Source: 소스 코드 {
  Go: "Go 코드\n(main.go, 서비스)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "프런트엔드 코드\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: 빌드 프로세스 {
  AnalyseGo: Go 코드 분석 {
    shape: rectangle
  }

  GenerateBindings: 바인딩 생성 {
    shape: rectangle
  }

  BuildFrontend: 프런트엔드 빌드 {
    shape: rectangle
  }

  CompileGo: Go 컴파일 {
    shape: rectangle
  }

  Embed: 애셋 임베드 {
    shape: rectangle
  }
}

Output: 출력 {
  Binary: "네이티브 바이너리\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: 타입 추출
Build.GenerateBindings -> Source.Frontend: TypeScript 바인딩
Source.Frontend -> Build.BuildFrontend: 컴파일(Vite/webpack)
Build.BuildFrontend -> Build.Embed: 번들된 애셋
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**빌드 단계:**

1. **Go 코드 분석**
  - 서비스에서 내보낸 메서드 스캔
  - 매개변수 및 반환 타입 추출
  - 메서드 시그니처 생성


2. **TypeScript 바인딩 생성**
  - 각 서비스의 `.ts` 파일 생성
  - 전체 타입 정의 포함
  - JSDoc 주석 추가


3. **프런트엔드 빌드**
  - 번들러(Vite, webpack 등) 실행
  - 축소 및 최적화
  - `frontend/dist/`에 출력


4. **Go 컴파일**
  - 최적화를 적용하여 컴파일(`-ldflags="-s -w"`)
  - 빌드 메타데이터 포함
  - 플랫폼별 컴파일


5. **애셋 임베드**
  - 프런트엔드 파일을 Go 바이너리에 임베드
  - 애셋 압축
  - 단일 실행 파일 생성


**결과:** 모든 항목이 내장된 단일 네이티브 실행 파일입니다.

[빌드에 대해 자세히 알아보기 →](/guides/build/building/)

## 개발 환경과 프로덕션 환경

Wails는 개발 환경과 프로덕션 환경에서 다르게 동작합니다.

@tabs{sync-key="mode"}
[개발 환경(wails3 dev)]
**특징:**

- **핫 리로드**: 프런트엔드 변경 사항을 즉시 다시 로드합니다
- **소스 맵**: 원본 소스로 디버깅할 수 있습니다
- **DevTools**: 브라우저 DevTools를 사용할 수 있습니다
- **로깅**: 상세 로깅이 활성화됩니다
- **외부 프런트엔드**: 개발 서버(Vite)에서 제공합니다

**작동 방식:**

```d2
direction: right

WailsApp: Wails 앱 {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Vite 개발 서버\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: 요청 프록시
DevServer -> WebView: HMR로 제공
WebView -> WailsApp: Go 메서드 호출
```

**장점:**

- 변경 사항에 대한 즉각적인 피드백
- 완전한 디버깅 기능
- 더 빠른 반복 개발

[프로덕션(wails3 build)]
**특징:**

- **임베드된 자산**: 프런트엔드를 바이너리에 포함합니다
- **최적화**: 축소 및 압축됩니다
- **DevTools 없음**: 기본적으로 비활성화됩니다
- **최소한의 로깅**: 오류만 기록합니다
- **단일 파일**: 모든 요소를 하나의 실행 파일에 포함합니다

**작동 방식:**

```d2
direction: right

Binary: "단일 바이너리\n(myapp.exe)" {
  GoCode: 컴파일된 Go 코드 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "임베드된 자산\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: 메모리에서 제공
WebView -> Binary.GoCode: Go 메서드 호출
```

**장점:**

- 단일 파일로 배포
- 더 작은 크기(축소됨)
- 더 나은 성능
- 외부 종속성 없음

@end

## 메모리 모델

메모리 사용 방식을 이해하면 효율적인 애플리케이션을 빌드하는 데 도움이 됩니다.

**메모리 영역:**

1. **Go 힙**
  - 서비스와 애플리케이션 상태가 저장됩니다
  - Go 가비지 컬렉터가 관리합니다
  - 간단한 앱에서는 일반적으로 5-10MB를 사용합니다


2. **WebView 메모리**
  - DOM, JavaScript 힙, CSS가 저장됩니다
  - WebView 엔진이 관리합니다
  - 간단한 앱에서는 일반적으로 10-20MB를 사용합니다


3. **브리지 메모리**
  - 통신용 메시지 버퍼가 저장됩니다
  - 오버헤드가 매우 적습니다(<1MB)
  - 가능한 경우 대용량 데이터에 제로 카피를 사용합니다


**최적화 팁:**

- **대용량 데이터 전송 피하기**: ID를 전달하고 필요할 때 세부 정보를 가져오세요
- **업데이트에 이벤트 사용하기**: 프런트엔드에서 폴링하지 마세요
- **대용량 파일 스트리밍하기**: 전체 파일을 메모리에 로드하지 마세요
- **리스너 정리하기**: 사용을 마치면 이벤트 리스너를 제거하세요

[성능에 대해 자세히 알아보기 →](/guides/performance/)

## 보안 모델

Wails는 기본적으로 안전한 아키텍처를 제공합니다:

```d2
direction: down

Frontend: 프런트엔드(신뢰할 수 없음) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails 브리지(검증) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: 백엔드(신뢰할 수 있음) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: 메서드 호출
Bridge -> Bridge: "검증:\n- 메서드가 존재하는가?\n- 타입이 올바른가?\n- 접근이 허용되는가?"
Bridge -> Backend: 유효한 경우 실행
Backend -> Bridge: 결과 반환
Bridge -> Frontend: 응답 전송
```

**보안 기능:**

1. **메서드 허용 목록**
  - 내보낸 메서드만 호출할 수 있습니다
  - 비공개 메서드에는 접근할 수 없습니다
  - 서비스를 명시적으로 등록해야 합니다


2. **타입 검증**
  - 인수를 Go 타입과 대조하여 검사합니다
  - 잘못된 타입은 거부됩니다
  - 인젝션 공격을 방지합니다


3. **eval() 없음**
  - 프런트엔드에서 임의의 Go 코드를 실행할 수 없습니다
  - 미리 정의된 메서드만 호출할 수 있습니다
  - 동적 코드를 실행하지 않습니다


4. **컨텍스트 격리**
  - 각 창에는 자체 컨텍스트가 있습니다
  - 서비스에서 호출자의 컨텍스트를 확인할 수 있습니다
  - 창별 권한을 설정할 수 있습니다


**권장 사항:**

- Go에서 **사용자 입력을 검증하세요**(프런트엔드를 신뢰하지 마세요)
- 인증 및 인가에 **컨텍스트를 사용하세요**
- 파일 작업 전에 **파일 경로를 정제하세요**
- 비용이 많이 드는 작업의 **호출 빈도를 제한하세요**

[보안에 대해 자세히 알아보기 →](/guides/security/)

## 다음 단계

**애플리케이션 수명 주기** - 시작, 종료 및 수명 주기 훅을 이해하세요 [자세히 알아보기 →](/concepts/lifecycle/)

**Go-프런트엔드 브리지** - 브리지의 작동 방식을 자세히 살펴보세요 [자세히 알아보기 →](/concepts/bridge/)

**빌드 시스템** - Wails가 애플리케이션을 빌드하는 방식을 이해하세요 [자세히 알아보기 →](/concepts/build-system/)

**빌드 시작하기** - 튜토리얼에서 배운 내용을 적용해 보세요. [튜토리얼 →](/tutorials/03-notes-vanilla/)

---

**아키텍처에 대해 궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [API 참조 문서](/reference/overview/)를 확인하세요.
