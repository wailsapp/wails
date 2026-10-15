---
title: "资源服务器"
description: "Wails v3 如何在开发和生产环境中提供并嵌入 Web 资源"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## 概述

每个 Wails 应用都以一个集成了以下内容的<strong>单一原生可执行文件</strong>交付：

1. *Go* 后端
2. *Web* 前端（HTML + JS + CSS）

<strong>资源服务器</strong>是实现这一点的纽带。它有<strong>两种运行模式</strong>，在编译时通过 Go 构建标签选择：

| 模式 | 标签 | 用途 |
| --- | --- | --- |
| **开发** | `//go:build !production` | 通过热重载快速迭代 |
| **生产** | `//go:build production` | 零依赖的嵌入式资源 |

其实现位于`v3/internal/assetserver/`中，并清晰地拆分为多个文件：

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## 开发模式

### 生命周期

1. `wails3 dev`启动后，会通过运行`build/Taskfile.yml`中定义的任务（通常为`npm run dev`）来<strong>生成前端开发服务器进程</strong>（Vite、SvelteKit、React-SWC 等）。
2. CLI 将`WAILS_VITE_PORT`设为 Wails 开发端口，并将`FRONTEND_DEVSERVER_URL`设为指向框架当前所运行开发服务器的<strong>完整</strong> URL（`http://host:port` / `https://host:port`）。请参阅`internal/commands/dev.go`。
3. 开发资源服务器（通过`build_dev.go`中的`//go:build !production`编译进程序）经由`GetDevServerURL()`读取`FRONTEND_DEVSERVER_URL`，并将非运行时流量反向代理到该地址。
4. 静态文件（`/assets/logo.svg`）可通过`asset_fileserver.go`**直接从磁盘提供**（速度更快），而所有无法识别的请求都会被<strong>代理</strong>到框架开发服务器，从而实现<em>即时</em>热模块替换。

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### 功能

- **实时重载** — Vite、SvelteKit 等通过 WebSocket 注入 HMR；开发资源服务器会透明地代理该连接。
- **源映射支持** — 由于资源未被打包，浏览器开发者工具可将错误映射回原始源代码。
- **无需重新编译 Go** — 只有前端会重新构建；在您更改`.go`文件之前，Go 代码会持续运行。

### 切换框架

开发代理<strong>与框架无关</strong>。Wails CLI 启动开发任务时会发布两个环境变量：

| 环境变量 | 来源 | 含义 |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go`（`wailsVitePort`常量） | 默认开发端口（除非传入`--port`，否则为9245）— Vite 配置应遵循此设置 |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | Wails 将代理到的完整 URL；在 Go 中通过`assetserver.GetDevServerURL()`（`build_dev.go`）读取 |

v3 源代码树中没有`VITE_PORT`、`FRONTEND_DEV_PORT`或`WAILSDEV_VERBOSE`环境变量。

添加新模板 → 定义其开发任务 → 资源服务器即可直接工作。

---

## 生产模式

运行`wails3 build`时，该流水线会：

1. 运行前端<strong>生产构建</strong>（`npm run build`），生成`frontend/dist/**`。
2. 通过应用自身包中的`go:embed`，将该目录<strong>嵌入</strong>应用（通常是位于`main.go`旁边的`//go:embed all:frontend/dist`）。
3. 使用`-tags production`编译 Go 二进制文件（由 Taskfile 包装器通过`EXTRA_TAGS`转发）。

`internal/assetserver/build_production.go`是构建标签存根，用于切换到生产代码路径。`internal/assetserver/bundled_assetserver.go`是<strong>手工编写的</strong>；它封装了位于`bundledassets/`中的运行时 JS，并非生成的文件。

### 请求处理

实际的处理程序是`internal/assetserver/assetserver.go` / `asset_fileserver.go`。从概念上讲，它会：

1. 尝试从请求路径读取嵌入的静态资源。
2. 如果未找到，则回退到`index.html`以支持 SPA 路由。
3. 如果扩展名未知，则探测内容类型（`content_type_sniffer.go`）。
4. 设置合理的缓存标头。

- **MIME 检测** — 对于没有扩展名的文件，根据开头约512字节（`content_type_sniffer.go`）探测其内容类型，并将结果缓存在`mimecache.go` / `ringqueue.go`中。
- **安全标头** — 禁止`file://`导航并设置`nosniff`。

由于所有内容均已嵌入，交付的二进制文件<strong>没有外部依赖</strong>（即使在 Windows 上也是如此）。

---

## 衔接开发与生产模式

从`pkg/application`的角度看，两种模式都公开<strong>相同的公共接口</strong>：一个包含`Handler http.Handler`的`AssetOptions`结构体，以及`internal/assetserver/`中的中间件和生命周期连接逻辑。开发/生产模式完全通过 Go 构建标签切换，因此两种模式下的应用代码完全相同。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## 前端框架的集成方式

### 模板

随附的每个模板（React、Vue、Svelte、Solid、Vanilla 等）都包含：

- `build/Taskfile.yml`
- `frontend/vite.config.ts`（或等效配置）

其 Vite 或等效配置会读取`WAILS_VITE_PORT`，并将开发服务器绑定到该端口。随后，CLI 会发布实际的`FRONTEND_DEVSERVER_URL`，供应用内代理使用。

前端框架与 Go 保持完全解耦：

- 构建时无需导入任何 Wails JS SDK——`/wails/runtime.js`由资源服务器在运行时提供。
- 任何提供 HTTP 开发服务器的框架都可以接入。

---

## 扩展与自定义

需要自定义标头、身份验证或 gzip 吗？

1. 定义一个`middleware.Middleware`（`func(http.Handler) http.Handler`的别名，在`internal/assetserver/middleware.go`中声明）。
2. 通过`internal/assetserver/options.go`公开的配置，将其接入你的`application.AssetOptions`。
3. 开发和生产环境中的行为完全一致——不存在按模式分别设置的中间件列表。

---

## 关键源文件

| 文件 | 作用 |
| --- | --- |
| `build_dev.go` / `build_production.go` | 用于选择开发或生产模式的构建标签封装 |
| `assetserver.go` / `asset_fileserver.go` | 核心 HTTP 处理程序 |
| `assetserver_dev.go` | 指向`FRONTEND_DEVSERVER_URL`的反向代理 |
| `bundled_assetserver.go` | 围绕`bundledassets/`手动编写的封装 |
| `options.go` | 面向`application.AssetOptions`的配置 |
| `mimecache.go` / `ringqueue.go` | MIME 缓存和小型 LRU 缓存 |

---

## 常见问题与调试

- **生产环境中出现白屏**——通常是 SPA 路由问题：请确保开发服务器针对未知路径提供`index.html`，并且能够执行嵌入式生产处理程序的回退逻辑。
- **开发环境中出现404**——你的 Vite 配置未绑定到`WAILS_VITE_PORT`，或者 CLI 无法连接开发服务器，因而未能填充`FRONTEND_DEVSERVER_URL`。
- **大型资源**——嵌入资源会增大二进制文件。请从单独的源提供大型媒体文件，或通过自定义`http.Handler`以流式方式传输。

---

现在，你已经了解 Wails <strong>资源服务器</strong>如何在<strong>开发</strong>和<strong>生产</strong>环境中将 Web 代码提供给原生窗口。掌握这一层后，你就能自信地调试加载问题、添加中间件，甚至换用完全不同的前端工具链。
