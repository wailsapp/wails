---
title: "Wails 확장하기"
description: "Wails v3에 새로운 기능과 플랫폼을 추가하는 실용 가이드"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails는 **쉽게 수정하고 확장할 수 있도록** 설계되었습니다.
>
> 모든 주요 하위 시스템은 직접 읽고 수정하여 배포할 수 있는 Go 코드로 구현되어 있습니다.
>
> 다음 작업을 수행할 때 *어디서* 시작하고 *어떻게* 크로스 플랫폼을 유지하는지 이 페이지에서 설명합니다.

- **서비스** 추가(알림, KV 저장소, 사용자 정의 IPC 등)
- **새 CLI 명령** 만들기(`wails3 <foo>`)
- **런타임** 확장(창 API, 대화 상자, 이벤트)
- **플랫폼 기능** 도입(Wayland 등)
- `//go:build` 태그에 파묻히지 않으면서 **크로스 플랫폼 호환성** 유지하기

---

## 1. 서비스 추가하기

v3에서 "서비스"란  `application.Options.Services`을 통해 등록하고 생성된 바인딩을 통해 JS에 노출하는 사용자 제공 Go 타입입니다. v3 코드베이스에는 다음이 포함되어 있습니다.

- `internal/service/` — `wails3 generate service`용 스캐폴딩:
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — 지금 바로 등록하여 사용할 수 있는 서비스(알림, kvstore, sqlite, log, fileserver, dock 등).

이전 초안에서 참조한 생성기 및 CLI 파일인 `internal/service/template/template.go`과 `internal/generator/collect/services.go`은 존재하지 않습니다. 스캐폴더는 `internal/service/service.go`이고 진입점은 `service.Install`이며, 서비스의 바인딩 메타데이터는 `internal/generator/collect/service.go`에서 수집됩니다.

### 1.1 서비스 정의하기

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 수명 주기 인터페이스 구현하기(선택 사항)

서비스는 선택적으로 다음 인터페이스(`pkg/application`에 정의됨)를 충족할 수 있습니다.

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **중요:** `ServiceShutdown`은 **인수를 받지 않습니다**. 다음
>
> 시그니처가 있는 메서드 `ServiceShutdown(ctx context.Context) error`은 해당 인터페이스를 **충족하지**
>
> 않으며 아무런 경고 없이 호출되지 않습니다.

### 1.3 애플리케이션에 서비스 등록하기

전역 `services.Register(...)` 호출은 없습니다. 서비스는 런타임에 `application.Options.Services`을 통해 등록합니다.

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

등록하고 나면 `wails3 generate bindings`이 내보낸 메서드를 래핑하는 ES 모듈을 `frontend/bindings/<your import path>/...` 아래에 생성합니다.

### 1.4 JS에서 호출하기

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

v3에는 전역 `window.backend.*`이 없습니다. 호출은 생성된 ES 모듈을 거치며, 이 모듈은 다시 `/wails/runtime.js`의 `Call.ByID(...)`을 호출합니다.

---

## 2. 새 CLI 명령 작성하기

v3 CLI는 cobra가 아닌 <strong>`github.com/leaanthony/clir`</strong>을 사용합니다. 연결 코드는 `v3/cmd/wails3/main.go`에 있습니다.

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

`init()` 기반 자동 등록은 없습니다. 새 하위 명령을 `cmd/wails3/main.go`에 추가하고, 이를 지원하는 함수를 `internal/commands/`에 추가하세요. 옵션을 받는 경우에는 `internal/flags/` 아래에 플래그 구조체도 추가하세요. CLI를 다시 빌드합니다.

```
cd v3
go install ./cmd/wails3
wails3 hello
```

명령에 Taskfile 연결 코드가 필요하다면 `internal/commands/task_wrapper.go`의 도우미(`wrapTask("yourtask", args)`)를 재사용하세요.

---

## 3. 런타임 수정하기

일반적인 수정 이유는 다음과 같습니다.

- 새로운 창 기능(`SetOpacity`, `Shake` 등)
- 추가 대화 상자(`ColorPicker`)
- 시스템 수준 API(화면 밝기)

### 3.1 공개 API

`pkg/application/webview_window.go`에 메서드를 추가하세요. 인터페이스는 `window.go`에 있습니다.

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

호출이 메인 스레드에서 실행되도록 기존 `InvokeSync`/`InvokeAsync` 도우미를 사용하세요.

### 3.2 메시지 프로세서

JS에서 새 메서드를 호출해야 한다면 해당 `pkg/application/messageprocessor_*.go` 파일을 확장하세요. 메시지 프로세서는 전역 `register(...)` 호출이 아니라 `MessageProcessor`의 switch 기반 메서드를 사용합니다.

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

코드베이스에는 `messageprocessor_window_opacity.go` 파일도, `init()` 기반 `register(MsgSetOpacity, ...)` 패턴도 **없습니다**.

### 3.3 플랫폼 구현

`pkg/application/` 아래의 OS별 파일마다 구현을 추가하세요.

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

플랫폼에서 해당 기능을 지원할 수 없다면 아무 작업도 하지 않는 스텁을 작성하세요. 프레임워크에는 `ErrCapability` 센티널이 없습니다. 지원 여부는 문서에 명시하고, 필요한 경우 `Options` 또는 플랫폼별 옵션 구조체의 관련 불리언 필드를 통해 노출하세요.

### 3.4 기능 플래그(선택 사항)

`internal/capabilities/` 패키지는 플랫폼별 기능 집합을 선언하는 데 사용됩니다. 공개 `application.HasCapability` / `application.CapOpacity` API는 없습니다. 런타임에 확인할 수 있는 기능이 필요하다면 `internal/capabilities/` 아래에 추가하고 `pkg/application`에서 타입이 지정된 getter를 노출하세요.

---

## 4. 새로운 플랫폼 기능 추가하기

예: Linux에서 선택적으로 지원하는 Wayland.

1. 관련 `pkg/application/*_linux.go` 파일을 `*_linux_x11.go`(`//go:build linux && !wayland`)과 `*_linux_wayland.go`(`//go:build linux && wayland`)로 분할하세요.
2. 사용자가 `wails3 build --tags wayland`으로 명시적으로 사용하도록 설정하게 하세요. 추가 태그는 `internal/commands/task_wrapper.go`에 있는 기존 `EXTRA_TAGS` 연결 코드를 통해 전달하세요. `dev` 수준의 `--tags wayland` 플래그는 없습니다. `wails3 dev`은 `--config`, `--port`, `-s`만 받습니다.
3. 문서와 `pkg/application/` 아래의 모든 플랫폼별 README를 업데이트하세요.

> 기본 빌드 태그는 최소한으로 유지하고, 특수한 기능에만 선택적 사용 태그를 지정하세요.

---

## 5. 크로스 플랫폼 호환성 체크리스트

| ✅ 단계 | 이유 |
| --- | --- |
| 모든 플랫폼 파일에 **모든** 공개 메서드를 제공하세요(스텁이라도 제공). | 모든 OS에서 빌드가 성공하도록 유지합니다. |
| OS별 단계적 기능 저하 동작을 문서화하세요. | 앱에서 숨겨진 오류 없이 `runtime.GOOS`을 기준으로 분기할 수 있습니다. |
| 먼저 <strong>순수 Go</strong>를 사용하고, 필요한 경우에만 Cgo를 사용하세요. | 크로스 컴파일이 간소화됩니다(Linux에서는 이미 Cgo 사용에 따른 비용을 감수하고 있습니다). |
| `task test:cli`, `task test:generator`, `task test:templates`을 실행하세요. | CI를 로컬에서 재현합니다. |
| 새 빌드 태그를 기여자 문서/템플릿 README에 문서화하세요. | 사용자가 선택적으로 활성화해야 하는 기능을 알고 있어야 합니다. |

---

## 6. 디버그 빌드 및 반복 작업 속도

- 자세한 런타임 활동을 출력하려면 `Options.LogLevel = slog.LevelDebug`(`Options.Logger = slog.Default()`)을 사용하세요. `WAILS_LOG_LEVEL` 환경 변수는 없습니다.
- `wails3 dev` 플래그는 `--config`, `--port`, `-s`입니다. `-race` 또는 `-verbose` 플래그는 없습니다. 레이스 탐지기를 사용하려면 `go test -race ./...`을 실행하거나 앱을 `go build -race`한 후 직접 실행하세요.
- 레이스/Cgo 테스트 가이드는 `v3/TESTING.md`에 있습니다(이전 초안에서는 존재하지 않는 `pkg/application/RACE.md`을 가리켰습니다).

---

## 7. 업스트림 기여

1. 새 기능이나 공개 동작의 변경 사항은 아이디어와 설계를 논의할 수 있도록 **WEP(Wails Enhancement Proposal)** 초안 PR을 여세요. 재현 가능한 버그나 문서 문제에만 이슈를 사용하세요.
2. 위의 기법에 따라 구현하세요.
3. 다음을 추가하세요:
  - 단위 테스트(`*_test.go`)
  - 문서(이 파일 또는 관련 `docs/...` 페이지)
  - 바인딩 생성기를 수정했다면 `internal/generator/testcases/` 아래에 회귀 테스트

4. 푸시하기 전에 로컬에서 `task precommit`과 관련 `task test:*` 대상을 실행하세요.

---

### 빠른 링크

| 영역 | 위치 |
| --- | --- |
| 기본 제공 서비스 | `pkg/services/` |
| 서비스 스캐폴더 | `internal/service/` |
| CLI 연결 | `v3/cmd/wails3/main.go` |
| CLI 명령 본문 | `internal/commands/` |
| OS별 런타임 | `pkg/application/*_{darwin,linux,windows}.go` |
| 기능 선언 | `internal/capabilities/` |
| Taskfile DSL | `v3/Taskfile.yaml` |
| 이벤트 상수 생성기 | `v3/tasks/events/generate.go` |

---

이제 Wails를 원하는 대로 활용하기 위한 <strong>로드맵</strong>을 갖추셨습니다. 서비스를 추가하고, CLI에 마법을 더하고, 런타임을 수정하거나 완전히 새로운 OS 기능을 도입해 보세요. 즐겁게 확장해 보세요!
