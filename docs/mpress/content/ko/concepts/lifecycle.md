---
title: "애플리케이션 수명 주기"
description: "시작부터 종료까지 Wails 애플리케이션의 수명 주기 이해하기"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## 애플리케이션 수명 주기 이해하기

데스크톱 애플리케이션에는 시작부터 종료까지 이어지는 수명 주기가 있습니다. Wails v3에서는 **서비스**, **이벤트**, <strong>훅</strong>을 통해 이 수명 주기를 효과적으로 관리할 수 있습니다.

## 수명 주기 단계

```d2
direction: down

Start: 애플리케이션 시작 {
  shape: oval
  style.fill: "#10B981"
}

Init: 초기화 {
  Parse: 옵션 파싱 {
    shape: rectangle
  }
  Register: 서비스 등록 {
    shape: rectangle
  }
  Setup: 런타임 설정 {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: 서비스 시작 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: 이벤트 루프 {
  Process: 이벤트 처리 {
    shape: rectangle
  }
  Handle: 메시지 처리 {
    shape: rectangle
  }
  Update: UI 업데이트 {
    shape: rectangle
  }
}

QuitSignal: 종료 신호 {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: ShouldQuit 확인 {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: OnShutdown 콜백 {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: 서비스 종료 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: 정리 {
  Close: 창 닫기 {
    shape: rectangle
  }
  Release: 리소스 해제 {
    shape: rectangle
  }
}

End: 애플리케이션 종료 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: 반복
EventLoop.Process -> QuitSignal: 사용자가 종료
QuitSignal -> ShouldQuit: 종료 허용 여부 확인
ShouldQuit -> EventLoop.Process: 거부됨
ShouldQuit -> OnShutdown: 허용됨
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. 애플리케이션 생성

`application.New()`을 사용하여 애플리케이션을 생성합니다:

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**처리 내용:**

1. 옵션을 구문 분석하고 검증합니다
2. 서비스를 등록합니다(아직 시작하지는 않음)
3. 애셋 서버를 구성합니다
4. 런타임을 설정합니다

### 2. 애플리케이션 실행

애플리케이션을 시작하려면 `app.Run()`을 호출합니다:

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**처리 내용:**

1. 서비스를 등록 순서대로 시작합니다
2. 이벤트 리스너를 활성화합니다
3. 창을 생성할 수 있습니다
4. 이벤트 루프를 시작합니다

### 3. 이벤트 루프

애플리케이션은 대부분의 시간을 보내는 이벤트 루프로 진입합니다:

- OS 이벤트를 처리합니다(마우스, 키보드, 창 이벤트)
- Go에서 JS로 보내는 메시지를 처리합니다
- JS에서 Go로 보내는 호출을 실행합니다
- UI 업데이트를 렌더링합니다

### 4. 종료

애플리케이션을 종료할 때 다음 작업을 수행합니다:

1. `ShouldQuit` 콜백이 설정되어 있으면 확인합니다
2. `OnShutdown` 콜백을 실행합니다
3. 서비스를 등록의 역순으로 종료합니다
4. 창을 닫습니다
5. 리소스를 해제합니다

## 서비스 수명 주기

서비스는 Wails v3에서 수명 주기를 관리하는 기본 수단입니다. 인터페이스를 통해 시작 및 종료 훅을 제공합니다. 서비스에 관한 전체 문서는 [서비스 가이드](/features/bindings/services/)를 참조하세요.

### 서비스 생성

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### 서비스 등록

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**핵심 사항:**

- 서비스는 등록 순서대로 시작됩니다
- 서비스는 등록 순서의 <strong>역순</strong>으로 종료됩니다
- 서비스의 `ServiceStartup`이 오류를 반환하면 애플리케이션이 중단됩니다
- `ServiceStartup`에 전달된 `ctx`은 종료가 시작될 때 취소됩니다

### 애플리케이션 컨텍스트 사용

`ServiceStartup`에 전달되는 컨텍스트는 애플리케이션의 전체 수명 주기 동안 유효합니다.

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

애플리케이션 인스턴스에서도 컨텍스트에 접근할 수 있습니다.

```go
app := application.Get()
ctx := app.Context()
```

## 애플리케이션 수준 훅

`application.Options`의 편의 콜백을 사용하면 완전한 서비스를 만들지 않고도 애플리케이션 수명 주기에 훅을 연결할 수 있습니다. 간단한 정리 작업이나 종료 확인 또는 종료 절차의 특정 시점에 코드를 실행해야 할 때 유용합니다.

시작 로직, 의존성 주입 또는 상태를 유지하는 리소스를 포함한 더 복잡한 수명 주기 관리가 필요하면 대신 [서비스](#---1)를 사용하세요.

### ShouldQuit

마지막 창을 닫거나 Cmd+Q(macOS)/Alt+F4(Windows)를 누르는 등 사용자가 종료를 요청하거나, 코드에서 `app.Quit()`을 호출할 때마다 `ShouldQuit` 콜백이 호출됩니다.

**반환 값:**

- 종료를 계속 진행하도록 허용하려면 `true`을 반환하세요(애플리케이션이 종료됩니다).
- 종료를 취소하려면 `false`을 반환하세요(애플리케이션이 계속 실행됩니다).

여기에서 종료 요청을 가로채고 필요에 따라 막을 수 있습니다. 예를 들어 저장하지 않은 변경 사항이 있음을 사용자에게 알리고 확인을 요청할 수 있습니다.

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

`ShouldQuit`이 설정되어 있지 않으면 요청 즉시 애플리케이션이 종료됩니다.

**ShouldQuit이 호출되는 경우:**

- 사용자가 마지막 창을 닫는 경우(`DisableQuitOnLastWindowClosed`이 설정되어 있지 않은 경우)
- 사용자가 macOS에서 Cmd+Q를 누르는 경우
- 사용자가 Windows에서 Alt+F4를 누르는 경우(마지막 창에 포커스가 있을 때)
- 코드에서 `app.Quit()`을 호출하는 경우

**ShouldQuit이 호출되지 않는 경우:**

- 프로세스가 강제 종료되는 경우(SIGKILL, 작업 관리자의 강제 종료)
- `os.Exit()`이 직접 호출되는 경우

### OnShutdown

애플리케이션 종료가 확정되면 `OnShutdown` 콜백이 호출됩니다(`ShouldQuit`이 설정된 경우 이 콜백이 `true`을 반환한 후). 상태 저장, 데이터베이스 연결 닫기 또는 리소스 해제와 같은 정리 작업에 사용하세요.

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

애플리케이션 수명 주기 중 언제든지 코드에서 종료 콜백을 추가로 등록할 수도 있습니다.

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

여러 콜백은 등록된 순서대로 실행됩니다. 모든 콜백이 완료될 때까지 종료 절차가 대기합니다.

**중요:** 종료 콜백은 빠르게(1초 미만) 실행되도록 유지하세요. 종료하는 데 너무 오래 걸리는 애플리케이션은 운영 체제에서 강제로 종료할 수 있으며, 이 경우 정리 작업이 중단되어 데이터가 손실될 수 있습니다.

### PostShutdown

모든 종료 작업이 완료된 후 프로세스가 끝나기 직전에 `PostShutdown` 콜백이 호출됩니다. 이 시점에는 창이 닫히고 서비스가 종료되며 리소스가 해제되었으므로 애플리케이션 인스턴스를 더 이상 사용할 수 없습니다.

주로 다음 용도로 유용합니다.

- 다른 모든 정리 작업 후에 수행해야 하는 최종 로깅
- 종료 동작 테스트 및 디버깅
- `app.Run()`이 반환되지 않는 플랫폼(콜백을 통해 코드 실행 보장)

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**참고:** `PostShutdown`에서는 애플리케이션 기능(창, 대화 상자 등)을 사용하지 마세요. 해당 기능은 더 이상 사용할 수 없습니다.

## 이벤트 기반 수명 주기

Wails는 창 열기, 애플리케이션 시작, 테마 변경 등 애플리케이션에서 일어나는 일을 알려 주는 이벤트 시스템을 제공합니다. 이러한 이벤트를 수신하면 수명 주기 변경을 차단하거나 가로채지 않고도 그에 대응할 수 있습니다.

창 이벤트의 경우 `OnWindowEvent` 대신 `RegisterHook`을 사용하여 동작을 가로채고 취소할 수도 있습니다. 예를 들어 창이 닫히지 않도록 할 수 있습니다. 아래의 [창 훅](#----2)을 참조하세요.

이벤트 시스템에 관한 전체 문서는 [이벤트 가이드](/features/events/system/)를 참조하세요.

### 애플리케이션 이벤트

애플리케이션 수명 주기 이벤트를 수신하세요.

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

플랫폼별 이벤트도 사용할 수 있습니다.

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### 창 이벤트

창 수명 주기 이벤트를 수신하세요.

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### 창 훅(취소 가능한 이벤트)

이벤트를 <strong>취소</strong>해야 할 때는 `OnWindowEvent` 대신 `RegisterHook`을 사용하세요.

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**OnWindowEvent와 RegisterHook의 차이:**

- `OnWindowEvent`: 이벤트 발생을 알립니다(취소할 수 없음).
- `RegisterHook`: 이벤트를 가로채고 필요한 경우 취소할 수 있습니다.

## 창 수명 주기

창에는 생성부터 소멸까지 자체 수명 주기가 있습니다. 각 창은 프런트엔드 콘텐츠를 독립적으로 로드하며 언제든지 표시하거나 숨기거나 닫을 수 있습니다. 사용자가 창을 닫으려고 할 때 `RegisterHook`으로 이를 가로채 확인을 요청하거나, 창을 소멸하는 대신 숨길 수 있습니다.

창에 관한 전체 문서는 [창 가이드](/features/windows/basics/)를 참조하세요.

```d2
direction: down

Create: 창 생성 {
  shape: oval
  style.fill: "#10B981"
}

Load: 프런트엔드 로드 {
  shape: rectangle
}

Show: 창 표시 {
  shape: rectangle
}

Active: 창 활성 상태 {
  Events: 이벤트 처리 {
    shape: rectangle
  }
}

CloseRequest: 닫기 요청 {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: WindowClosing 훅 {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: 창 제거 {
  shape: rectangle
}

End: 창 닫힘 {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: 반복
Active.Events -> CloseRequest: 사용자가 닫음
CloseRequest -> Hook
Hook -> Active.Events: 취소됨
Hook -> Destroy: 허용됨
Destroy -> End
```

### 창 생성

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### 창 닫기 방지

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### 닫는 대신 숨기기

시스템 트레이 앱에서 흔히 사용하는 패턴은 다음과 같습니다.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## 다중 창 수명 주기

창이 여러 개인 경우:

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**기본 동작은 플랫폼에 따라 다릅니다.**

| 플랫폼 | 마지막 창이 닫힐 때의 기본 동작 |
| --- | --- |
| macOS | 앱이 계속 실행됨(메뉴 막대 유지) |
| Windows | 앱 종료 |
| Linux | 앱 종료 |

macOS는 창이 하나도 없어도 일반적으로 애플리케이션이 메뉴 막대에서 활성 상태로 유지되는 플랫폼 고유의 관례를 따릅니다. Windows와 Linux에서는 기본적으로 앱이 종료됩니다.

**모든 플랫폼에서 마지막 창이 닫히면 앱이 종료되도록 설정:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**모든 플랫폼에서 마지막 창이 닫혀도 앱이 계속 실행되도록 설정:**

이는 시스템 트레이 애플리케이션이나 백그라운드에서 계속 실행되어야 하는 앱에 유용합니다.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## 일반적인 패턴

### 패턴 1: 데이터베이스 서비스

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### 패턴 2: 구성 서비스

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### 패턴 3: 백그라운드 워커

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## 수명 주기 참조

| 훅/인터페이스 | 호출 시점 | 취소 가능 여부 | 용도 |
| --- | --- | --- | --- |
| `ServiceStartup` | `app.Run()` 도중, 이벤트 루프 시작 전 | 불가(중단하려면 오류 반환) | 초기화 |
| `ServiceShutdown` | 종료 과정 중, `OnShutdown` 이후 | 불가 | 정리 |
| `OnShutdown` | 종료가 확인되었을 때 | 불가 | 애플리케이션 정리 |
| `ShouldQuit` | 종료가 요청되었을 때 | 가능(false 반환) | 종료 확인 |
| `RegisterHook(WindowClosing)` | 창 닫기가 요청되었을 때 | 가능(`e.Cancel()`) | 창 닫기 방지 |
| `OnWindowEvent` | 이벤트가 발생했을 때 | 불가 | 이벤트에 대응 |
| `OnApplicationEvent` | 이벤트가 발생했을 때 | 불가 | 이벤트에 대응 |

## 플랫폼별 차이

### macOS

- 창이 없어도 <strong>애플리케이션 메뉴</strong>가 유지됩니다
- <strong>Cmd+Q</strong>를 누르면 종료가 시작됩니다(`ShouldQuit`을 거침)
- 숨기지 않는 한 <strong>Dock 아이콘</strong>이 유지됩니다
- 종료 동작을 제어하려면 `ApplicationShouldTerminateAfterLastWindowClosed`을 사용하세요

### Windows

- 창이 없으면 **애플리케이션 메뉴도 없습니다**
- <strong>Alt+F4</strong>를 누르면 창이 닫힙니다(`RegisterHook`을 사용해 방지할 수 있음)
- <strong>시스템 트레이</strong>를 사용하면 앱을 계속 실행할 수 있습니다

### Linux

- 데스크톱 환경에 따라 **동작이 달라집니다**
- **일반적으로 Windows와 유사합니다**

## 수명 주기 문제 디버깅

### 문제: 애플리케이션이 종료되지 않음

**원인:**

1. `ShouldQuit`이 `false`을 반환함
2. `OnShutdown`에 시간이 너무 오래 걸림
3. 백그라운드 goroutine이 중지되지 않음

**해결 방법:**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### 문제: 서비스 시작 실패

**해결 방법:** 설명이 명확한 오류를 반환하세요:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

오류가 로그에 기록되고 애플리케이션은 시작되지 않습니다.

## 모범 사례

### 권장 사항

- **수명 주기 관리에 서비스를 사용하세요** - 서비스는 적절한 시작 및 종료 훅을 제공합니다
- **빠르게 종료하세요** - 모든 정리 작업에 걸리는 시간이 1초 미만이 되도록 목표를 잡으세요
- **취소 처리에 context를 사용하세요** - 백그라운드 작업을 올바르게 중지하세요
- **시작 중 발생하는 오류를 처리하세요** - 오류를 반환하여 문제없이 중단하세요
- **수명 주기 이벤트를 로그에 기록하세요** - 디버깅에 도움이 됩니다

### 금지 사항

- **서비스 시작을 블로킹하지 마세요** - 초기화를 빠르게 완료하세요(2초 미만)
- **종료 중에는 대화 상자를 표시하지 마세요** - 앱이 종료 중이므로 UI가 작동하지 않을 수 있습니다
- **context를 무시하지 마세요** - goroutine에서는 항상 `ctx.Done()`을 확인하세요
- **리소스를 누수하지 마세요** - 항상 `ServiceShutdown`을 구현하세요

## 다음 단계

**서비스** - 서비스 시스템에 대해 자세히 알아보세요 [자세히 알아보기 →](/features/bindings/services/)

**이벤트 시스템** - 이벤트를 사용해 통신하세요 [자세히 알아보기 →](/features/events/system/)

**창 관리** - 창을 만들고 관리하세요 [자세히 알아보기 →](/features/windows/basics/)

---

**수명 주기에 관해 궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
