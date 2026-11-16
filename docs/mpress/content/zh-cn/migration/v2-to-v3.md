---
title: "从v2迁移到v3"
description: "将 Wails v2应用迁移到v3的完整指南"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3是一次**彻底重写**，在架构、性能和开发者体验方面均有显著改进。本指南将帮助你把v2应用迁移到v3。

**主要变化：**

- 新的应用结构
- 改进的绑定系统
- 增强的窗口管理
- 更完善的事件系统
- 简化的配置

<strong>迁移时间：</strong>典型应用需要1-4小时

## 破坏性变更

### 应用初始化

在v2中，应用设置、窗口配置和执行全都合并在一次`wails.Run()`调用中。这种单体式做法使创建多个窗口、处理不同阶段的错误或测试应用的各个组件变得困难。

v3将这些职责分为不同阶段：创建应用、创建窗口和执行。这样，你便可以明确控制应用生命周期的每个阶段，并使代码更加模块化且更易于测试。

**v2：**

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

**v3：**

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

**这种方式更好的原因：**

- **多窗口支持**：你可以随时动态创建窗口，而不仅限于启动时
- **更完善的错误处理**：可以通过适当的错误处理分别验证每个阶段
- **更清晰的代码**：职责分离让每个阶段执行的操作一目了然
- **更易于测试**：无需运行事件循环即可测试应用设置
- **更灵活**：在应用的整个生命周期中都可以创建、销毁和重新创建窗口

### 绑定

在v2中，每个绑定的结构体都必须包含上下文字段，并通过`startup(ctx)`方法接收运行时上下文。这使业务逻辑与 Wails 运行时紧密耦合，导致代码更难测试和理解。

v3引入了服务模式。在此模式下，结构体完全独立，无需存储运行时上下文。如果服务需要访问应用实例，则通过依赖注入显式接收该实例，而不是隐式传递上下文。

**v2：**

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

**v3：**

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

**这种方式更好的原因：**

- **没有隐式依赖**：服务是普通的 Go 结构体，不含隐藏的运行时依赖
- **更易于测试**：无需模拟 Wails 上下文即可测试服务方法
- **更清晰的代码**：依赖关系是显式的（作为构造函数参数传入），而不是隐藏在上下文字段中
- **更合理的组织方式**：可以按领域对服务分组，而不必将所有服务都放在单个`App`结构体中
- **规范的初始化**：需要初始化时使用`ServiceStartup()`方法，使初始化过程清晰明确

### 运行时

在v2中，所有运行时操作都必须向`runtime`包中的全局函数传递上下文。这导致整个代码库与上下文对象紧密耦合，也让 API 更像面向过程而非面向对象。

v3不再使用基于上下文的运行时，而是直接调用应用对象和窗口对象的方法。操作直接在其影响的对象上调用，使代码更直观，也更符合面向对象的方式。

**v2：**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3：**

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

**这种方式更好的原因：**

- **面向对象设计**：方法在其影响的对象（窗口、应用、菜单等）上调用
- **意图更明确**：`window.SetTitle()`比`runtime.WindowSetTitle(ctx, ...)`更直观
- **更完善的 IDE 支持**：方法位于对象上时，自动补全可以正常工作
- **多窗口操作更清晰**：使用多个窗口时，你需要明确选择要操作的窗口
- **无需传递上下文**：不必在每个函数之间传递上下文

### 前端绑定

在v2中，绑定按 Go 包和结构体名称组织，通常会生成类似`wailsjs/go/main/App`的路径。这种结构无法反映逻辑分组，也很难找到相关功能。

v3按服务名称和应用模块组织绑定，从而形成更清晰的逻辑结构。绑定会生成到`bindings`目录中，并按应用名称和服务名称组织，让你更容易了解有哪些可用功能。

**v2：**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**这种方式更好的原因：**

- **合理的逻辑组织**：绑定按服务名称分组，而不是按 Go 包结构分组
- **更清晰的导入路径**：路径反映领域逻辑（greetservice），而不是文件结构（main/App）
- **更易于发现**：可以按功能浏览绑定，而不是按技术结构浏览
- **一致的命名**：基于服务的组织方式与后端架构保持一致
- **更简单的路径**：不再需要`../wailsjs/go`前缀，只需使用`./bindings`

### 事件

在v2中，事件使用可变数量的`interface{}`参数，并且每次调用事件函数时都必须传递上下文。事件处理程序接收的是无类型数据，需要手动进行类型断言，导致事件系统容易出错且难以调试。

v3引入了有类型的事件对象，并取消了上下文要求。事件处理程序会接收包含类型化数据的规范事件对象，使事件系统更可靠、更易于使用。

**v2：**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3：**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**这种方式更好的原因：**

- **类型安全**：事件使用规范的事件对象，而不是`...interface{}`
- **更易调试**：事件对象包含事件名称等元数据，使调试更加容易
- **API 更清晰**：`app.Event.On()`和`app.Event.Emit()`比运行时函数更直观
- **无需上下文**：事件可直接在应用对象上使用，无需传递上下文
- **处理程序更简单**：事件处理程序使用明确的签名，而不是可变参数

### 窗口

v2每个应用仅支持一个窗口。该窗口在启动时创建，所有窗口操作都通过运行时函数执行，并隐式作用于这唯一的窗口。

v3将原生多窗口支持作为核心功能引入。每个窗口都是一等对象，拥有自己的方法和生命周期。在应用的整个生命周期内，你可以动态创建、管理和销毁多个窗口。

**v2：**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3：**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**这样做更好的原因：**

- **多窗口应用**：构建拥有多个独立窗口的应用（如仪表板、偏好设置和工具等）
- **显式窗口引用**：每个窗口都是可存储并直接操作的对象
- **动态创建窗口**：可在运行期间随时创建和销毁窗口
- **独立窗口状态**：每个窗口都有自己的事件、属性和生命周期
- **架构更合理**：窗口管理采用面向对象方式，而不是基于上下文

## 迁移步骤

### 步骤1：更新依赖项

**go.mod：**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**更新后：**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### 步骤2：更新 main.go

**v2：**

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

**v3：**

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

### 步骤3：将 App 结构体转换为服务

**v2：**

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

**v3：**

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

### 步骤4：更新运行时调用

**v2：**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3：**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### 步骤5：更新前端

**生成新绑定：**

```bash
wails3 generate bindings
```

**更新导入：**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**更新事件处理：**

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

### 步骤6：更新配置

**v2（wails.json）：**

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

**v3（wails.json）：**

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

## 功能映射

### 对话框

**v2：**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3：**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### 菜单

**v2：**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3：**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### 系统托盘

**v2：**

```go
// Not available in v2
```

**v3：**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## 常见问题

### 问题：找不到绑定

<strong>问题：</strong>迁移后出现导入错误

**解决方案：**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### 问题：上下文错误

**问题：**`ctx`不可用

**解决方案：**

改为存储应用引用：

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### 问题：窗口方法不起作用

**问题：**`runtime.WindowSetTitle()`不存在

**解决方案：**

直接使用窗口方法：

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### 问题：事件未触发

<strong>问题：</strong>事件已注册，但未收到

**解决方案：**

检查事件名称是否完全匹配：

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## 测试迁移结果

### 检查清单

- [ ]应用启动时无错误
- [ ]所有绑定均正常工作
- [ ]事件可以正常发送和接收
- [ ]窗口可以正常打开和关闭
- [ ]菜单正常工作（如适用）
- [ ]对话框正常工作（如适用）
- [ ]系统托盘正常工作（如适用）
- [ ]构建流程正常工作
- [ ]生产构建正常工作

### 测试命令

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## v3的优势

### 性能

- **启动更快** - 优化了初始化过程
- **内存占用更低** - 高效利用资源
- **桥接性能更好** - 调用开销为&lt;1ms

### 功能

- **多窗口** - 原生支持
- **系统托盘** - 内置支持
- **更完善的事件机制** - 类型化且更简洁的 API
- **服务** - 更合理的代码组织

### 开发者体验

- **类型安全** - 全面的 TypeScript 支持
- **更完善的错误处理** - 清晰的错误消息
- **热重载** - 加快开发速度
- **更完善的文档** - 全面的指南

## 获取帮助

### 资源

- [文档](/quick-start/why-wails/)
- [Discord 社区](https://discord.gg/JDdSxwjhGf)
- [GitHub Issue](https://github.com/wailsapp/wails/issues)
- [示例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### 常见问题

<strong>问：我可以同时运行v2和v3吗？</strong>答：可以，它们使用不同的导入路径。

<strong>问：v3是否已可用于生产环境？</strong>Wails v3.0.0是面向 Windows、macOS 和 Linux 桌面应用程序的稳定版本。Android 和 iOS 支持仍处于实验阶段，不在稳定桌面兼容性承诺的范围内。

<strong>问：v2会继续维护吗？</strong>答：会，v2将获得关键更新。

<strong>问：迁移需要多长时间？</strong>答：典型应用程序需要1-4小时。

## 后续步骤

@cards{cols="2"}
🚀 快速入门
开始使用 Wails v3。

[了解更多 →](/quick-start/installation/)

---
★ 核心概念
了解v3架构。

[了解更多 →](/concepts/architecture/)

---
◆ 绑定
了解新的绑定系统。

[了解更多 →](/features/bindings/methods/)

---
📖 示例
查看完整的v3示例。

[查看示例 →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或[提交 Issue](https://github.com/wailsapp/wails/issues)。
