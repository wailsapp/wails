---
title: "Wails 的工作原理"
description: "了解 Wails 的架构及其如何实现原生性能"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails 是一个桌面应用程序开发框架，使用<strong>Go 构建后端</strong>，并使用<strong>Web 技术构建前端</strong>。但与 Electron 不同，Wails 不会捆绑浏览器，而是使用<strong>操作系统的原生 WebView</strong>。

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Wails 应用

  frontend: 前端
  backend: Go 后端
  os: 操作系统

  Initialisation: 初始化 {
    shape: sequence_diagram
    backend."Serves Static Web App": 提供静态 Web 应用
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": 通过操作系统原生 WebView 渲染站点
  }
  Regular Communication: 常规通信 {
    shape: sequence_diagram
    frontend."Make API-style call": 发起 API 风格的调用
    frontend -> backend.a: JSON
    backend.a."Service processes request": 服务处理请求
    backend.a -> os: 调用系统 API
    backend.a."Generate Response": 生成响应
    backend.a -> frontend: JSON
    frontend."Process response": 处理响应
  }
  backend.a.label: a
}
```

**与 Electron 的主要区别：**

| 方面 | Wails | Electron |
| --- | --- | --- |
| **浏览器** | 操作系统提供的 WebView | 捆绑的 Chromium（约 100MB） |
| **后端** | Go（编译执行） | Node.js（解释执行） |
| **通信** | 内存桥接 | IPC（进程间通信） |
| **包体积** | 约 15MB | 约 150MB |
| **内存** | 约 10MB | 约 100MB+ |
| **启动时间** | &lt;0.5 秒 | 2-3s |

## 核心组件

### 1. 原生 WebView

Wails 使用操作系统内置的 Web 渲染引擎：

@tabs{sync-key="platform"}
[Windows]
**WebView2**（Microsoft Edge WebView2）

- 基于 Chromium（与 Edge 浏览器相同）
- 预装在 Windows 10/11上
- 通过 Windows Update 自动更新
- 全面支持现代 Web 标准

[macOS]
**WebKit**（Safari 的渲染引擎）

- 内置于 macOS
- 与 Safari 浏览器使用相同的引擎
- 性能出色且续航表现优异
- 全面支持现代 Web 标准

[Linux]
**WebKitGTK**（WebKit 的 GTK 移植版本）

- 通过软件包管理器安装
- 与 GNOME Web（Epiphany）使用相同的引擎
- 良好的标准支持
- 轻量且高性能

@end

**这为何重要：**

- **不捆绑浏览器** → 应用体积更小
- **操作系统原生** → 集成度和性能更好
- **自动更新** → 通过操作系统更新获取安全补丁
- **熟悉的渲染效果** → 与系统浏览器一致

### 2. Wails 桥接器

桥接器是 Wails 的核心，它支持 Go 与 JavaScript 之间的<strong>直接通信</strong>。

```d2
direction: down

Frontend: 前端（JavaScript） {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Wails 桥接层 {
  Encoder: JSON 编码器 {
    shape: rectangle
  }

  Router: 方法路由器 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 解码器 {
    shape: rectangle
  }
}

Backend: 后端（Go） {
  Services: 已注册的服务 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. 调用 Go 方法\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. 编码为 JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. 路由到服务\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. 返回结果\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. 解码到 JS\nPromise 得到兑现"
```

**工作原理：**

1. **前端调用 Go 方法**（通过自动生成的绑定）
2. <strong>桥接器将调用编码</strong>为 JSON（方法名 + 参数）
3. **路由器在已注册的服务中查找 Go 方法**
4. <strong>Go 方法执行</strong>并返回值
5. <strong>桥接器解码结果</strong>并将其发回前端
6. JavaScript 中的<strong>Promise 使用该结果完成兑现</strong>

**性能特征：**

- **内存通信**：无网络开销，不使用 HTTP
- **尽可能实现零拷贝**（适用于大型数据）
- **默认异步**：两端均不阻塞
- **类型安全**：自动生成 TypeScript 定义

### 3. 服务系统

服务是向前端公开 Go 功能的推荐方式。

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

**服务发现：**

- Wails 在启动时<strong>扫描你的结构体</strong>
- <strong>导出的方法</strong>可由前端调用
- 提取<strong>类型信息</strong>以生成 TypeScript 绑定
- <strong>错误处理</strong>会自动完成（Go 错误 → JS 异常）

**生成的 TypeScript 绑定：**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**为什么使用服务？**

- **类型安全**：完整支持 TypeScript
- **自动发现**：无需手动注册方法
- **组织有序**：将相关功能归为一组
- **可测试**：服务只是 Go 结构体

[详细了解服务 →](/features/bindings/services/)

### 4. 事件系统

事件支持组件之间进行<strong>发布/订阅通信</strong>。

```d2
direction: left

Wails Event System: Wails 事件系统 {
  shape: sequence_diagram

  window1: 窗口 1
  window2: 窗口 2
  backend: Go 后端

  Event Driver: 事件驱动器 {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "订阅 'data-updated' 事件"
    window2."Subscribe to 'data-updated' events": "订阅 'data-updated' 事件"
    backend.a."App Emit('data-updated', data)": "应用发出 Emit('data-updated', data)"
    backend.a -> window1.a: JSON 事件总线
    backend.a -> window2: JSON 事件总线
    window1.a."Subscriber processes On('data-updated', handler)": "订阅者通过 On('data-updated', handler) 处理"
    window2."Subscriber processes On('data-updated', handler)": "订阅者通过 On('data-updated', handler) 处理"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**用例：**

- **窗口通信**：一个窗口通知其他窗口
- **后台任务**：Go 服务向 UI 通知进度
- **状态同步**：使多个窗口保持同步
- **松耦合**：组件无需直接引用彼此

**示例：**

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

[详细了解事件 →](/features/events/system/)

## 应用生命周期

了解生命周期有助于确定何时初始化和清理资源。

```d2
direction: down

Start: 应用启动 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初始化 {
  Create: 创建应用 {
    shape: rectangle
  }

  Register: 注册服务 {
    shape: rectangle
  }

  Setup: 设置窗口/菜单 {
    shape: rectangle
  }
}

Run: 事件循环 {
  Events: 处理事件 {
    shape: rectangle
  }

  Messages: 处理消息 {
    shape: rectangle
  }

  Render: 更新 UI {
    shape: rectangle
  }
}

Shutdown: 关闭 {
  Cleanup: 清理资源 {
    shape: rectangle
  }

  Save: 保存状态 {
    shape: rectangle
  }
}

End: 应用结束 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: 循环
Run.Events -> Shutdown.Cleanup: 退出信号
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**生命周期钩子：**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options`上没有`OnStartup`字段。启动工作应放在服务的`ServiceStartup(ctx, options)`中、通过`app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`注册的回调中，或者直接放在`app.Run()`之前。

[详细了解生命周期 →](/concepts/lifecycle/)

## 构建流程

了解 Wails 如何构建应用：

```d2
direction: down

Source: 源代码 {
  Go: "Go 代码\n(main.go, 服务)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "前端代码\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: 构建流程 {
  AnalyseGo: 分析 Go 代码 {
    shape: rectangle
  }

  GenerateBindings: 生成绑定 {
    shape: rectangle
  }

  BuildFrontend: 构建前端 {
    shape: rectangle
  }

  CompileGo: 编译 Go 代码 {
    shape: rectangle
  }

  Embed: 嵌入资源 {
    shape: rectangle
  }
}

Output: 输出 {
  Binary: "原生二进制文件\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: 提取类型
Build.GenerateBindings -> Source.Frontend: TypeScript 绑定
Source.Frontend -> Build.BuildFrontend: 编译（Vite/webpack）
Build.BuildFrontend -> Build.Embed: 已打包的资源
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**构建步骤：**

1. **分析 Go 代码**
  - 扫描服务中的导出方法
  - 提取参数类型和返回类型
  - 生成方法签名


2. **生成 TypeScript 绑定**
  - 为每个服务创建`.ts`文件
  - 包含完整的类型定义
  - 添加 JSDoc 注释


3. **构建前端**
  - 运行打包工具（Vite、webpack 等）
  - 缩小并优化
  - 输出到`frontend/dist/`


4. **编译 Go**
  - 启用优化进行编译（`-ldflags="-s -w"`）
  - 包含构建元数据
  - 针对特定平台进行编译


5. **嵌入资源**
  - 将前端文件嵌入 Go 二进制文件
  - 压缩资源
  - 创建单个可执行文件


<strong>结果：</strong>一个嵌入了所有内容的原生可执行文件。

[详细了解构建 →](/guides/build/building/)

## 开发环境与生产环境

Wails 在开发环境和生产环境中的行为有所不同：

@tabs{sync-key="mode"}
[开发环境（wails3 dev）]
**特性：**

- **热重载**：前端更改会立即重新加载
- **源映射**：使用原始源代码进行调试
- **DevTools**：可使用浏览器 DevTools
- **日志记录**：启用详细日志记录
- **外部前端**：由开发服务器（Vite）提供

**工作原理：**

```d2
direction: right

WailsApp: Wails 应用 {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Vite 开发服务器\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: 代理请求
DevServer -> WebView: 通过 HMR 提供内容
WebView -> WailsApp: 调用 Go 方法
```

**优势：**

- 立即获得更改反馈
- 完整的调试能力
- 更快地迭代

[生产环境（wails3 build）]
**特性：**

- **嵌入式资源**：前端内置于二进制文件中
- **已优化**：经过缩小和压缩
- **无 DevTools**：默认禁用
- **最少日志记录**：仅记录错误
- **单文件**：所有内容都包含在一个可执行文件中

**工作原理：**

```d2
direction: right

Binary: "单个二进制文件\n(myapp.exe)" {
  GoCode: 已编译的 Go 代码 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "嵌入式资源\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: 从内存提供内容
WebView -> Binary.GoCode: 调用 Go 方法
```

**优势：**

- 单文件分发
- 体积更小（已压缩）
- 性能更佳
- 无外部依赖

@end

## 内存模型

了解内存使用情况有助于构建高效的应用程序。

**内存区域：**

1. **Go 堆**
  - 服务和应用程序状态
  - 由 Go 垃圾回收器管理
  - 简单应用程序通常为 5-10MB


2. **WebView 内存**
  - DOM、JavaScript 堆、CSS
  - 由 WebView 引擎管理
  - 简单应用程序通常为 10-20MB


3. **桥接层内存**
  - 用于通信的消息缓冲区
  - 开销极小（<1MB）
  - 尽可能对大型数据采用零拷贝


**优化建议：**

- **避免传输大量数据**：传递 ID，按需获取详细信息
- **使用事件进行更新**：不要从前端轮询
- **以流式方式处理大文件**：不要将整个文件加载到内存中
- **清理监听器**：使用完毕后移除事件监听器

[详细了解性能 →](/guides/performance/)

## 安全模型

Wails 提供默认安全的架构：

```d2
direction: down

Frontend: 前端（不受信任） {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails 桥接层（验证） {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: 后端（受信任） {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: 调用方法
Bridge -> Bridge: "验证：\n- 方法是否存在？\n- 类型是否正确？\n- 是否允许访问？"
Bridge -> Backend: 验证通过后执行
Backend -> Bridge: 返回结果
Bridge -> Frontend: 发送响应
```

**安全功能：**

1. **方法白名单**
  - 只能调用已导出的方法
  - 无法访问私有方法
  - 必须显式注册服务


2. **类型验证**
  - 根据 Go 类型检查参数
  - 拒绝无效类型
  - 防止注入攻击


3. **不使用 eval()**
  - 前端无法执行任意 Go 代码
  - 只能调用预定义的方法
  - 不执行动态代码


4. **上下文隔离**
  - 每个窗口都有自己的上下文
  - 服务可以检查调用方上下文
  - 可以为每个窗口设置权限


**最佳实践：**

- 在 Go 中<strong>验证用户输入</strong>（不要信任前端）
- 使用<strong>上下文</strong>进行身份验证和授权
- 执行文件操作前<strong>清理文件路径</strong>
- 对开销较大的操作实施<strong>速率限制</strong>

[详细了解安全性 →](/guides/security/)

## 后续步骤

**应用程序生命周期** - 了解启动、关闭和生命周期钩子 [了解更多 →](/concepts/lifecycle/)

**Go-前端桥接** - 深入了解桥接的工作原理 [了解更多 →](/concepts/bridge/)

**构建系统** - 了解 Wails 如何构建应用程序 [了解更多 →](/concepts/build-system/)

**开始构建** - 在教程中运用所学知识 [教程 →](/tutorials/03-notes-vanilla/)

---

<strong>对架构有疑问？</strong>请前往 [Discord](https://discord.gg/JDdSxwjhGf) 提问，或查阅[API 参考](/reference/overview/)。
