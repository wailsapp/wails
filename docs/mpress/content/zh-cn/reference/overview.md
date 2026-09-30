---
title: "API 参考"
description: "Wails v3 的完整 API 文档"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## 关于本参考文档

这是 Wails v3 的完整 API 参考文档，其中记录了框架提供的所有公共类型、方法和选项。

**组织结构：**

- [应用程序](/reference/application/) - 核心应用程序 API
- [窗口](/reference/window/) - 窗口创建和管理
- [菜单](/reference/menu/) - 应用程序菜单、上下文菜单和系统托盘菜单
- [事件](/reference/events/) - 事件系统和内置事件
- [对话框](/reference/dialogs/) - 文件对话框和消息对话框
- [前端运行时](/reference/frontend-runtime/) - 前端运行时 API
- [CLI](/reference/cli/) - 命令行界面

## API 约定

@details{title="Go API 约定 - 面向 Go 新手开发者"}
### 命名

- <strong></strong>类型<strong></strong>：PascalCase（例如 `WebviewWindow`）
- <strong></strong>方法<strong></strong>：PascalCase（例如 `SetTitle()`）
- <strong></strong>选项<strong></strong>：PascalCase 结构体（例如 `WindowOptions`）
- <strong></strong>常量<strong></strong>：PascalCase（例如 `WindowStartStateMaximised`）

#### 错误处理

大多数可能失败的方法都将`error`作为最后一个返回值。`app.Run()`会阻塞，直到应用退出，并返回任何启动错误：

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

创建窗口不会返回错误——`app.Window.New()`会直接返回`*WebviewWindow`。

#### 上下文

服务生命周期方法会接收一个`context.Context`：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

可通过`app.Context()`访问应用的生命周期上下文。不存在`RunWithContext`——请调用`app.Run()`。

#### 选项模式

配置使用选项结构体：

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

### JavaScript API 约定

#### 命名

- **函数**：camelCase（例如 `setTitle()`）
- **常量**：SCREAMING<em>SNAKE</em>CASE（例如 `WINDOW_EVENT_FOCUS`）

#### 默认异步

所有 Go 方法调用都返回 Promise：

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### 错误处理

Go 错误会转换为 JavaScript 异常：

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### 类型安全

TypeScript 定义会自动生成：

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## 包结构

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

## 导入路径

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

## 类型参考

### 常用类型

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

## 平台差异

某些 API 在不同平台上的行为有所不同：

| 功能 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **应用程序菜单** | 窗口菜单栏 | 全局菜单栏 | 窗口菜单栏 |
| **系统托盘** | 通知区域 | 菜单栏 | 系统托盘 |
| **程序坞** | 不适用 | ✅ 可用 | 不适用 |
| **文件对话框** | 原生 | 原生 | 原生（GTK） |
| **透明效果** | ✅ 完全支持 | 需要[`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ 支持有限 |

每个 API 章节中都记录了特定于平台的行为。

## 版本控制

Wails v3 遵循语义化版本控制：

- **主版本**（v3.x.x）：破坏性变更
- **次版本**（v3.x.x）：新增功能，向后兼容
- **补丁版本**（v3.x.x）：错误修复，向后兼容

<strong>当前状态：</strong>Beta（API 已稳定，仍在持续完善）

## 弃用政策

API 被弃用时：

1. <strong>在文档中标记</strong>并附上弃用通知
2. <strong>提供替代方案</strong>并附上迁移指南
3. **移除前仍会维护 1 个主版本**
4. **编译器警告**（如可实现）

## API 稳定性

### 稳定 API ✅

以下 API 已稳定，可安全用于生产环境：

- 核心应用 API
- 窗口管理
- 菜单系统
- 事件系统
- 文件对话框
- 服务绑定

### 不稳定 API ⚠️

以下 API 在正式发布前可能发生变化：

- 部分高级窗口选项
- 平台特定功能
- 实验性功能

不稳定 API 会在文档中标明。

## 获取帮助

### API 问题

1. **查阅此参考文档** — 完整的 API 文档
2. **查看示例** — [GitHub 示例](https://github.com/wailsapp/wails/tree/master/v3/examples)
3. **搜索 Discord** — [Discord 服务器](https://discord.gg/JDdSxwjhGf)
4. **向社区提问** — Discord #help 频道

### 报告 API 问题

发现了错误或不一致之处？

1. **查看现有问题** — [GitHub issues](https://github.com/wailsapp/wails/issues)
2. **创建详细报告** — 包含代码、错误信息和平台
3. **提供复现步骤** — 提供可重现该问题的最小示例

## 相关文档

- [教程](/tutorials/overview/) — 通过构建真实应用进行学习
- [指南](/guides/architecture/) — 面向常见场景的任务型指南
- [功能](/features/windows/basics/) — 按功能逐项介绍的文档
- [示例](https://github.com/wailsapp/wails/tree/master/v3/examples) — GitHub 上可运行的代码示例

---

<strong>浏览 API：</strong>使用左侧导航探索具体 API。
