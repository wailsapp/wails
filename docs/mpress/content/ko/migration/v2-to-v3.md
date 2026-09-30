---
title: "v2에서 v3로 마이그레이션"
description: "Wails v2 애플리케이션을 v3로 마이그레이션하기 위한 전체 가이드"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3는 아키텍처, 성능 및 개발자 경험을 크게 개선한 <strong>완전한 재작성 버전</strong>입니다. 이 가이드에서는 v2 애플리케이션을 v3로 마이그레이션하는 방법을 안내합니다.

**주요 변경 사항:**

- 새로운 애플리케이션 구조
- 개선된 바인딩 시스템
- 향상된 창 관리
- 개선된 이벤트 시스템
- 간소화된 구성

**마이그레이션 소요 시간:** 일반적인 애플리케이션의 경우 1-4시간

## 호환성을 깨는 변경 사항

### 애플리케이션 초기화

v2에서는 애플리케이션 설정, 창 구성 및 실행이 모두 하나의 `wails.Run()` 호출로 결합되어 있었습니다. 이러한 모놀리식 접근 방식 때문에 여러 창을 만들거나, 단계별로 오류를 처리하거나, 애플리케이션의 개별 구성 요소를 테스트하기가 어려웠습니다.

v3에서는 이러한 작업을 애플리케이션 생성, 창 생성 및 실행이라는 별도의 단계로 분리합니다. 이렇게 분리하면 애플리케이션 수명 주기의 각 단계를 명시적으로 제어할 수 있으며 코드를 더욱 모듈화하고 테스트하기 쉽게 만들 수 있습니다.

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
- **향상된 오류 처리**: 적절한 오류 처리를 통해 각 단계를 개별적으로 검증할 수 있습니다.
- **더 명확한 코드**: 단계를 분리하여 각 단계에서 어떤 작업이 수행되는지 명확하게 알 수 있습니다.
- **향상된 테스트 용이성**: 이벤트 루프를 실행하지 않고도 애플리케이션 설정을 테스트할 수 있습니다.
- **향상된 유연성**: 애플리케이션 수명 주기 전반에 걸쳐 창을 생성하고 제거한 후 다시 생성할 수 있습니다.

### 바인딩

v2에서는 바인딩된 모든 구조체에 컨텍스트 필드와 런타임 컨텍스트를 받는 `startup(ctx)` 메서드가 필요했습니다. 이로 인해 비즈니스 로직과 Wails 런타임이 강하게 결합되어 코드를 테스트하고 이해하기가 더 어려웠습니다.

v3에서는 서비스 패턴을 도입하여 구조체가 완전히 독립적으로 동작하며 런타임 컨텍스트를 저장할 필요가 없습니다. 서비스에서 애플리케이션 인스턴스에 접근해야 하는 경우에는 암시적으로 컨텍스트를 전달하지 않고 종속성 주입을 통해 명시적으로 전달받습니다.

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

- **암시적 종속성 제거**: 서비스는 숨겨진 런타임 종속성이 없는 일반 Go 구조체입니다.
- **더 쉬운 테스트**: Wails 컨텍스트를 모킹하지 않고도 서비스 메서드를 테스트할 수 있습니다.
- **더 명확한 코드**: 종속성이 컨텍스트 필드에 숨겨지지 않고 생성자 인수로 전달되어 명시적으로 드러납니다.
- **향상된 구성**: 모든 서비스를 하나의 `App` 구조체에 넣지 않고 도메인별로 그룹화할 수 있습니다.
- **명시적인 초기화**: 초기화가 필요할 때 `ServiceStartup()` 메서드를 사용하여 이를 명시적으로 드러낼 수 있습니다.

### 런타임

v2에서는 모든 런타임 작업을 수행할 때 `runtime` 패키지의 전역 함수에 컨텍스트를 전달해야 했습니다. 이로 인해 코드베이스 전반이 컨텍스트 객체와 강하게 결합되었고 API가 객체 지향적이라기보다 절차적으로 느껴졌습니다.

v3에서는 컨텍스트 기반 런타임을 애플리케이션 및 창 객체에 대한 직접적인 메서드 호출로 대체합니다. 작업의 영향을 받는 객체에서 메서드를 직접 호출하므로 코드가 더 직관적이고 객체 지향적으로 바뀝니다.

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

- **객체 지향 설계**: 메서드는 해당 작업의 영향을 받는 객체(창, 앱, 메뉴 등)에서 호출됩니다.
- **더 명확한 의도**: `window.SetTitle()`이 `runtime.WindowSetTitle(ctx, ...)`보다 의미가 더 명확합니다.
- **향상된 IDE 지원**: 메서드가 객체에 정의되어 있으면 자동 완성이 제대로 작동합니다.
- **더 명확한 다중 창 처리**: 창이 여러 개일 때 작업할 창을 명시적으로 선택합니다.
- **컨텍스트 전달 불필요**: 모든 함수에 컨텍스트를 전달할 필요가 없습니다.

### 프런트엔드 바인딩

v2에서는 바인딩이 Go 패키지와 구조체 이름을 기준으로 구성되어 일반적으로 `wailsjs/go/main/App`과 같은 경로가 생성되었습니다. 이 구조는 논리적인 그룹화를 반영하지 않아 관련 기능을 찾기가 어려웠습니다.

v3에서는 바인딩을 서비스 이름과 애플리케이션 모듈별로 구성하여 더 명확한 논리 구조를 만듭니다. 바인딩은 애플리케이션 이름과 서비스 이름에 따라 구성된 `bindings` 디렉터리에 생성되므로 어떤 기능을 사용할 수 있는지 더 쉽게 파악할 수 있습니다.

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

- **논리적 구성**: 바인딩은 Go 패키지 구조가 아니라 서비스 이름별로 그룹화됩니다.
- **더 명확한 임포트**: 경로가 파일 구조(main/App)가 아닌 도메인 로직(greetservice)을 반영합니다.
- **검색 용이성 향상**: 기술 구조가 아닌 기능별로 바인딩을 탐색할 수 있습니다
- **일관된 명명 방식**: 서비스 기반 구성이 백엔드 아키텍처와 일치합니다
- **더 간단한 경로**: 더 이상 `../wailsjs/go` 접두사가 필요하지 않으며 `./bindings`만 사용합니다

### 이벤트

v2에서는 이벤트에 가변 인자 `interface{}` 매개변수를 사용했으며 모든 이벤트 함수에 컨텍스트를 전달해야 했습니다. 이벤트 핸들러는 수동으로 타입 단언해야 하는 타입 미지정 데이터를 받았기 때문에 이벤트 시스템에서 오류가 발생하기 쉽고 디버깅하기도 어려웠습니다.

v3에서는 타입이 지정된 이벤트 객체를 도입하고 컨텍스트 요구 사항을 제거했습니다. 이벤트 핸들러는 타입이 지정된 데이터를 포함하는 올바른 이벤트 객체를 받으므로 이벤트 시스템의 안정성과 사용 편의성이 향상됩니다.

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

**개선된 점:**

- **타입 안전성**: 이벤트에서 `...interface{}` 대신 올바른 이벤트 객체를 사용합니다
- **향상된 디버깅**: 이벤트 객체에 이벤트 이름과 같은 메타데이터가 포함되어 디버깅이 더 쉬워집니다
- **더 명확한 API**: `app.Event.On()` 및 `app.Event.Emit()`가 런타임 함수보다 더 직관적입니다
- **컨텍스트 불필요**: 컨텍스트를 전달하지 않고 앱 객체에서 직접 이벤트를 사용할 수 있습니다
- **더 간단한 핸들러**: 이벤트 핸들러가 가변 인자 대신 명확한 시그니처를 사용합니다

### 창

v2에서는 애플리케이션당 하나의 창만 지원했습니다. 시작할 때 창이 생성되었으며, 모든 창 작업은 이 단일 창을 암시적으로 대상으로 하는 런타임 함수를 통해 수행되었습니다.

v3에서는 네이티브 다중 창 지원을 핵심 기능으로 도입했습니다. 각 창은 자체 메서드와 수명 주기를 갖는 일급 객체입니다. 애플리케이션이 실행되는 동안 여러 창을 동적으로 생성하고 관리하고 제거할 수 있습니다.

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

- **다중 창 애플리케이션**: 독립적인 창이 여러 개인 앱을 빌드할 수 있습니다(대시보드, 환경설정, 도구 등)
- **명시적인 창 참조**: 각 창은 저장하고 직접 조작할 수 있는 객체입니다
- **동적 창 생성**: 런타임 중 언제든 창을 생성하고 제거할 수 있습니다
- **독립적인 창 상태**: 각 창에는 자체 이벤트, 속성 및 수명 주기가 있습니다
- **향상된 아키텍처**: 창 관리가 컨텍스트 기반이 아닌 객체 지향 방식으로 이루어집니다

## 마이그레이션 단계

### 1단계: 종속성 업데이트

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**업데이트:**

```bash
go get github.com/wailsapp/wails/v3@latest
go mod tidy
```

### 2단계: main.go 업데이트

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

### 3단계: App 구조체를 서비스로 변환

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

### 4단계: 런타임 호출 업데이트

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

### 5단계: 프런트엔드 업데이트

**새 바인딩 생성:**

```bash
wails3 generate bindings
```

**가져오기 업데이트:**

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

### 6단계: 구성 업데이트

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

**문제:** 마이그레이션 후 가져오기 오류 발생

**해결 방법:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### 문제: 컨텍스트 오류

**문제:** `ctx`을(를) 사용할 수 없음

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

**문제:** `runtime.WindowSetTitle()`이(가) 존재하지 않음

**해결 방법:**

창 메서드를 직접 사용하세요:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### 문제: 이벤트가 발생하지 않음

**문제:** 이벤트를 등록했지만 수신되지 않음

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

- [ ] 애플리케이션이 오류 없이 시작됨
- [ ] 모든 바인딩이 작동함
- [ ] 이벤트가 송수신됨
- [ ] 창이 올바르게 열리고 닫힘
- [ ] 메뉴가 작동함(해당하는 경우)
- [ ] 대화 상자가 작동함(해당하는 경우)
- [ ] 시스템 트레이가 작동함(해당하는 경우)
- [ ] 빌드 프로세스가 작동함
- [ ] 프로덕션 빌드가 작동함

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

- **더 빠른 시작** - 초기화 최적화
- **더 적은 메모리 사용량** - 효율적인 리소스 사용
- **향상된 브리지** - 호출 오버헤드 &lt;1ms

### 기능

- **다중 창** - 네이티브 지원
- **시스템 트레이** - 기본 제공
- **향상된 이벤트** - 타입이 지정된 더 간단한 API
- **서비스** - 향상된 코드 구성

### 개발자 경험

- **타입 안전성** - 완전한 TypeScript 지원
- **향상된 오류 처리** - 명확한 오류 메시지
- **핫 리로드** - 더 빠른 개발
- **향상된 문서** - 포괄적인 가이드

## 도움말 보기

### 리소스

- [문서](/quick-start/why-wails/)
- [Discord 커뮤니티](https://discord.gg/JDdSxwjhGf)
- [GitHub 이슈](https://github.com/wailsapp/wails/issues)
- [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)

### 자주 묻는 질문

**질문: v2와 v3를 함께 실행할 수 있나요?** 답변: 예. 서로 다른 가져오기 경로를 사용합니다.

**질문: v3는 프로덕션 환경에서 사용할 준비가 되었나요?** 답변: v3는 안정적인 데스크톱 API를 갖춘 베타 소프트웨어입니다. 이미 프로덕션 환경에서 이를 사용해 실행되는 애플리케이션이 있지만, 배포하기 전에 철저히 테스트하세요. 현재 안정 버전은 여전히 v2입니다.

**질문: v2는 계속 유지보수되나요?** 답변: 예. v2에는 중요 업데이트가 제공됩니다.

**질문: 마이그레이션에는 얼마나 걸리나요?** 답변: 일반적인 애플리케이션의 경우 1-4시간이 걸립니다.

## 다음 단계

@cards{cols="2"}
🚀 빠른 시작
Wails v3를 시작하세요.

[자세히 알아보기 →](/quick-start/installation/)

---
★ 핵심 개념
v3 아키텍처를 이해하세요.

[자세히 알아보기 →](/concepts/architecture/)

---
◆ 바인딩
새로운 바인딩 시스템을 알아보세요.

[자세히 알아보기 →](/features/bindings/methods/)

---
📖 예제
완전한 v3 예제를 살펴보세요.

[예제 보기 →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

**질문이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [이슈를 등록하세요](https://github.com/wailsapp/wails/issues).
