---
title: "v2에서 v3(으)로 마이그레이션"
description: "Wails v2 애플리케이션을 v3(으)로 마이그레이션하는 전체 가이드"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3은 아키텍처, 성능, 개발자 경험이 크게 개선된 **완전한 재작성 버전**입니다. 이 가이드에서는 v2 애플리케이션을 v3(으)로 마이그레이션하는 방법을 안내합니다.

**주요 변경 사항:**

- 새로운 애플리케이션 구조
- 개선된 바인딩 시스템
- 향상된 창 관리
- 개선된 이벤트 시스템
- 간소화된 구성

**마이그레이션 소요 시간:** 일반적인 애플리케이션의 경우 1-4시간

## 호환성을 깨는 변경 사항

### 애플리케이션 초기화

v2에서는 애플리케이션 설정, 창 구성, 실행이 모두 하나의 `wails.Run()` 호출에 결합되어 있었습니다. 이러한 모놀리식 접근 방식 때문에 여러 창을 만들거나, 단계별로 오류를 처리하거나, 애플리케이션의 개별 구성 요소를 테스트하기가 어려웠습니다.

v3에서는 이러한 작업을 애플리케이션 생성, 창 생성, 실행이라는 별도의 단계로 분리합니다. 이렇게 분리하면 애플리케이션 수명 주기의 각 단계를 명시적으로 제어할 수 있으며, 코드가 더욱 모듈화되고 테스트하기 쉬워집니다.

**v2:**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3:**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**이 방식이 더 나은 이유:**

- **다중 창 지원**: 시작할 때뿐만 아니라 언제든지 동적으로 창을 만들 수 있습니다.
- **개선된 오류 처리**: 적절한 오류 처리를 통해 각 단계를 개별적으로 검증할 수 있습니다.
- **더 명확한 코드**: 단계를 분리하여 각 단계에서 어떤 작업이 이루어지는지 명확하게 알 수 있습니다.
- **향상된 테스트 용이성**: 이벤트 루프를 실행하지 않고도 애플리케이션 설정을 테스트할 수 있습니다.
- **향상된 유연성**: 애플리케이션 수명 주기 전반에 걸쳐 창을 생성하고 제거한 후 다시 생성할 수 있습니다.

### 바인딩

v2에서는 바인딩된 모든 구조체에 컨텍스트 필드와 런타임 컨텍스트를 받는 `startup(ctx)` 메서드가 필요했습니다. 이로 인해 비즈니스 로직과 Wails 런타임이 강하게 결합되어 코드를 테스트하고 이해하기가 더 어려웠습니다.

v3에서는 구조체가 완전히 독립적이며 런타임 컨텍스트를 저장할 필요가 없는 서비스 패턴을 도입합니다. 서비스가 애플리케이션 인스턴스에 접근해야 하는 경우, 암시적으로 컨텍스트를 전달하는 대신 의존성 주입을 통해 명시적으로 인스턴스를 받습니다.

**v2:**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**이 방식이 더 나은 이유:**

- **암시적 의존성 없음**: 서비스는 숨겨진 런타임 의존성이 없는 일반 Go 구조체입니다.
- **더 쉬운 테스트**: Wails 컨텍스트를 모의 객체로 만들지 않고도 서비스 메서드를 테스트할 수 있습니다.
- **더 명확한 코드**: 의존성이 컨텍스트 필드에 숨겨지지 않고 생성자 인수로 전달되어 명시적으로 드러납니다.
- **개선된 구성**: 모든 서비스를 하나의 `App` 구조체에 넣는 대신 도메인별로 그룹화할 수 있습니다.
- **적절한 초기화**: 초기화가 필요할 때 `ServiceStartup()` 메서드를 사용하여 이를 명시적으로 나타냅니다.

### 런타임

v2에서는 모든 런타임 작업에서 `runtime` 패키지의 전역 함수에 컨텍스트를 전달해야 했습니다. 이로 인해 코드베이스 전반이 컨텍스트 객체와 강하게 결합되었으며 API가 객체 지향적이기보다 절차적으로 느껴졌습니다.

v3에서는 컨텍스트 기반 런타임을 애플리케이션 및 창 객체의 직접적인 메서드 호출로 대체합니다. 작업이 영향을 미치는 객체에서 직접 작업을 호출하므로 코드가 더 직관적이고 객체 지향적으로 바뀝니다.

**v2:**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3:**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**이 방식이 더 나은 이유:**

- **객체 지향 설계**: 메서드는 해당 메서드가 영향을 미치는 객체(창, 앱, 메뉴 등)에서 호출됩니다.
- **더 명확한 의도**: `runtime.WindowSetTitle(ctx, ...)`보다 `window.SetTitle()`의 의미가 더 명확합니다.
- **향상된 IDE 지원**: 메서드가 객체에 속해 있으면 자동 완성이 올바르게 작동합니다.
- **명확한 다중 창 처리**: 창이 여러 개일 때 작업할 창을 명시적으로 선택합니다.
- **컨텍스트 전달 불필요**: 모든 함수에 컨텍스트를 전달할 필요가 없습니다.

### 프런트엔드 바인딩

v2에서는 바인딩이 Go 패키지와 구조체 이름을 기준으로 구성되어 일반적으로 `wailsjs/go/main/App` 같은 경로가 생성되었습니다. 이러한 구조는 논리적 그룹화를 반영하지 않아 관련 기능을 찾기 어려웠습니다.

v3에서는 서비스 이름과 애플리케이션 모듈을 기준으로 바인딩을 구성하여 논리적 구조를 더 명확하게 만듭니다. 바인딩은 애플리케이션 이름과 서비스 이름에 따라 구성된 `bindings` 디렉터리에 생성되므로 사용할 수 있는 기능을 더 쉽게 파악할 수 있습니다.

**v2:**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**이 방식이 더 나은 이유:**

- **논리적 구성**: 바인딩이 Go 패키지 구조가 아닌 서비스 이름을 기준으로 그룹화됩니다.
- **더 명확한 가져오기 경로**: 경로가 파일 구조(main/App)가 아닌 도메인 로직(greetservice)을 반영합니다.
- **향상된 검색 용이성**: 기술적 구조가 아닌 기능을 기준으로 바인딩을 탐색할 수 있습니다.
- **일관된 명명**: 서비스 기반 구성이 백엔드 아키텍처와 일치합니다.
- **더 간단한 경로**: 더 이상 `../wailsjs/go` 접두사를 사용하지 않고 `./bindings`만 사용합니다.

### 이벤트

v2에서는 이벤트가 가변 인자 `interface{}` 매개변수를 사용했으며 모든 이벤트 함수에 컨텍스트를 전달해야 했습니다. 이벤트 핸들러는 수동으로 타입 단언해야 하는 타입 없는 데이터를 받았으므로 이벤트 시스템에서 오류가 발생하기 쉽고 디버깅하기도 어려웠습니다.

v3에서는 타입이 지정된 이벤트 객체를 도입하고 컨텍스트 요구 사항을 제거합니다. 이벤트 핸들러가 타입이 지정된 데이터를 포함한 올바른 이벤트 객체를 받으므로 이벤트 시스템의 안정성과 사용 편의성이 향상됩니다.

**v2:**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3:**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**이 방식이 더 나은 이유:**

- **타입 안전성**: 이벤트는 `...interface{}` 대신 적절한 이벤트 객체를 사용합니다.
- **더 쉬운 디버깅**: 이벤트 객체에는 이벤트 이름과 같은 메타데이터가 포함되어 있어 디버깅이 더 쉽습니다.
- **더 명확한 API**: `app.Event.On()` 및 `app.Event.Emit()`는 런타임 함수보다 더 직관적입니다.
- **컨텍스트 불필요**: 이벤트는 스레딩 컨텍스트 없이 앱 객체에서 직접 작동합니다.
- **더 간단한 핸들러**: 이벤트 핸들러는 가변 인자 대신 명확한 시그니처를 사용합니다.

### 창

v2에서는 애플리케이션당 하나의 창만 지원했습니다. 창은 시작 시 생성되었으며, 모든 창 작업은 암시적으로 해당 단일 창을 대상으로 하는 런타임 함수를 통해 수행되었습니다.

v3에서는 네이티브 다중 창 지원을 핵심 기능으로 도입합니다. 각 창은 자체 메서드와 수명 주기를 갖는 일급 객체입니다. 애플리케이션이 실행되는 동안 여러 창을 동적으로 생성하고 관리하며 제거할 수 있습니다.

**v2:**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3:**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**개선된 점:**

- **다중 창 애플리케이션**: 서로 독립적인 여러 창이 있는 앱을 구축할 수 있습니다(대시보드, 환경설정, 도구 등).
- **명시적인 창 참조**: 각 창은 직접 저장하고 조작할 수 있는 객체입니다.
- **동적 창 생성**: 런타임 중 언제든지 창을 생성하고 제거할 수 있습니다.
- **독립적인 창 상태**: 각 창에는 자체 이벤트, 속성 및 수명 주기가 있습니다.
- **더 나은 아키텍처**: 창 관리는 컨텍스트 기반이 아니라 객체 지향 방식입니다.

## 마이그레이션 단계

### 단계 1: 종속성 업데이트

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**업데이트:**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### 단계 2: main.go 업데이트

**v2:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 단계 3: App 구조체를 서비스로 변환

**v2:**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### 단계 4: 런타임 호출 업데이트

**v2:**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3:**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### 단계 5: 프런트엔드 업데이트

**새 바인딩 생성:**

```bash
wails3 generate bindings
```

**임포트 업데이트:**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**이벤트 처리 업데이트:**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### 단계 6: 구성 업데이트

**v2 (wails.json):**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3 (wails.json):**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## 기능 매핑

### 대화상자

**v2:**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3:**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### 메뉴

**v2:**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3:**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### 시스템 트레이

**v2:**

```go
// Not available in v2
```

**v3:**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## 일반적인 문제

### 문제: 바인딩을 찾을 수 없음

**문제:** 마이그레이션 후 임포트 오류가 발생합니다.

**해결 방법:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### 문제: 컨텍스트 오류

**문제:** `ctx`을(를) 사용할 수 없습니다.

**해결 방법:**

대신 앱 참조를 저장하세요:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### 문제: 창 메서드가 작동하지 않음

**문제:** `runtime.WindowSetTitle()`이(가) 존재하지 않습니다.

**해결 방법:**

창 메서드를 직접 사용하세요:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### 문제: 이벤트가 발생하지 않음

**문제:** 이벤트가 등록되었지만 수신되지 않습니다.

**해결 방법:**

이벤트 이름이 정확히 일치하는지 확인하세요:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## 마이그레이션 테스트

### 체크리스트

- [ ] 애플리케이션이 오류 없이 시작됩니다.
- [ ] 모든 바인딩이 작동합니다.
- [ ] 이벤트가 전송되고 수신됩니다.
- [ ] 창이 올바르게 열리고 닫힙니다.
- [ ] 메뉴가 작동합니다(해당하는 경우).
- [ ] 대화상자가 작동합니다(해당하는 경우).
- [ ] 시스템 트레이가 작동합니다(해당하는 경우).
- [ ] 빌드 프로세스가 작동합니다.
- [ ] 프로덕션 빌드가 작동합니다.

### 테스트 명령어

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## v3의 이점

### 성능

- **더 빠른 시작** - 최적화된 초기화
- **더 적은 메모리 사용량** - 효율적인 리소스 사용
- **개선된 브리지** - 호출 오버헤드 &lt;1ms

### 기능

- **다중 창** - 네이티브 지원
- **시스템 트레이** - 기본 제공
- **개선된 이벤트** - 타입이 지정된 더 간단한 API
- **서비스** - 더 나은 코드 구성

### 개발자 경험

- **타입 안전성** - 완전한 TypeScript 지원
- **개선된 오류 처리** - 명확한 오류 메시지
- **핫 리로드** - 더 빠른 개발
- **개선된 문서** - 포괄적인 가이드

## 도움말

### 리소스

- [문서](/quick-start/why-wails/)
- [Discord 커뮤니티](https://discord.gg/JDdSxwjhGf)
- [GitHub 이슈](https://github.com/wailsapp/wails/issues)
- [예제](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### 자주 묻는 질문

**Q: v2과 v3을 함께 실행할 수 있나요?** A: 예. 서로 다른 import 경로를 사용합니다.

**Q: v3은 프로덕션 환경에서 사용할 준비가 되었나요?** Wails v3.0.0은 Windows, macOS, Linux용 데스크톱 애플리케이션의 안정 버전입니다. Android 및 iOS 지원은 여전히 실험 단계이며 안정적인 데스크톱 호환성 계약에 포함되지 않습니다.

**Q: v2은 계속 유지보수되나요?** A: 예. v2에는 중요 업데이트가 제공됩니다.

**Q: 마이그레이션에는 얼마나 걸리나요?** A: 일반적인 애플리케이션은 1-4시간이 걸립니다.

## 다음 단계

@cards{cols="2"}
🚀 빠른 시작
Wails v3을 시작해 보세요.

[자세히 알아보기 →](/quick-start/installation/)

---
★ 핵심 개념
v3 아키텍처를 이해해 보세요.

[자세히 알아보기 →](/concepts/architecture/)

---
◆ 바인딩
새로운 바인딩 시스템을 알아보세요.

[자세히 알아보기 →](/features/bindings/methods/)

---
📖 예제
완전한 v3 예제를 살펴보세요.

[예제 보기 →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [이슈를 등록하세요](https://github.com/wailsapp/wails/issues).
