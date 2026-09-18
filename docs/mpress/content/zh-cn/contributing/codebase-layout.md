---
title: "代码库布局"
description: "Wails v3 仓库的组织方式以及各部分如何协同工作"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 位于一个<strong>单体仓库</strong>中，其中包含框架运行时、CLI、示例、文档和构建工具链。 本页将介绍深入研究其内部实现时需要了解的<em>目录结构</em>。

## 顶层概览

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

接下来，我们将深入查看<strong>`v3/`</strong>目录树。

## `v3/`根目录

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> 项目模板位于`internal/templates/`下（每种框架
>
> 技术栈各有一个文件夹，此外还有`base/`、`_common/`和`ios/`）。顶层没有`v3/templates/`
>
> 目录。

### 心智模型

1. <strong>`pkg/`</strong>公开<em>应用开发者导入的内容</em>\
2. <strong>`internal/`</strong>包含<em>底层功能的实现方式</em>\
3. <strong>`cmd/wails3`</strong>驱动<em>项目生命周期和构建</em>\

其他所有内容都为这三大支柱提供支持。

---

## `cmd/`——命令

| 路径 | 说明 |
| --- | --- |
| `v3/cmd/wails3` | **CLI 入口点**。一个精简的`main.go`将所有逻辑委托给`internal/commands`中的包。 |
| `internal/commands/*` | 子命令（init、dev、build、doctor 等）。每个子命令各自位于单独的文件中，便于查找。 |
| `internal/commands/task_wrapper.go` | 在 CLI 标志与 Taskfile 构建流水线之间进行衔接。 |

CLI 负责：

- **项目脚手架**（`init`、模板生成）\
- **开发服务器编排**（`dev`、实时重载）\
- **生产构建和打包**（`build`、`package`、平台封装）\
- **诊断**（`doctor`）\

---

## `internal/`——引擎室

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### 关键子包

| 包 | 职责 | 连接位置 |
| --- | --- | --- |
| `runtime` | 包含少量`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`构建标签衔接代码，以及`runtime/desktop/`下嵌入的 JS 运行时。各操作系统实际使用的窗口、剪贴板、对话框和托盘代码位于`pkg/application/*_{darwin,linux,windows}.go`中。 | 通过`pkg/application`间接导入。 |
| `assetserver` | 双模式文件服务器：<br />• 开发：从磁盘提供文件并代理 Vite（`build_dev.go`）<br />• 生产：通过`go:embed`嵌入资源（`build_production.go`） | 启动期间由`pkg/application`初始化。 |
| `generator` | 解析 Go 源代码以构建<strong>绑定元数据</strong>，随后据此生成 TypeScript/JS 桩文件和事件常量。入口点：基于`collect/`和`render/`的`generator.Generate`/`generator.Generator`。 | 由`wails3 generate bindings`触发。 |
| `packager` | 用于生成 Linux `deb`/`rpm`/`archlinux`产物的`nfpm`封装器（由`internal/commands/`下的`myapp.DEB`/`.RPM`/`.ARCHLINUX` nfpm 配置驱动）。 | 由`wails3 tool package`调用。macOS DMG/Windows MSIX 位于`internal/commands/{dmg,msix.go,webview2/}`下。 |

辅助工具（例如`s/`、`hash/`、`flags/`）使各项内部关注点保持解耦。

---

## `pkg/`——公共 API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> 不存在`pkg/runtime/`、`pkg/options/`或`pkg/menu/`包。窗口和菜单
>
> 选项与`pkg/application`位于同一处（例如`WebviewWindowOptions`、`Menu`、
>
> `MenuItem`），而`assetserver/`位于`internal/`下。

`pkg/application`负责引导启动 Wails 程序：

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

其底层会：

1. 连接`internal/runtime`构建标签衔接代码与`pkg/application/`中各操作系统专用的代码
2. 设置一个`internal/assetserver`实例
3. 注册所有由绑定驱动的消息处理器
4. 进入操作系统主线程

---

## `internal/templates/`——脚手架蓝图

`internal/templates/`提供<strong>基础模板</strong>（Go 布局位于`base/`、 `_common/`和`ios/`下）以及<strong>前端皮肤</strong>（`vanilla[-ts]`、`react[-ts]`、 `react-swc[-ts]`、`lit[-ts]`、`preact[-ts]`、`qwik[-ts]`、`solid[-ts]`、 `svelte[-ts]`、`sveltekit[-ts]`、`vue[-ts]`）。

在执行`wails3 init -t react`时，CLI 会：

1. 复制`_common` Go 文件
2. 合并所需的前端包
3. 运行`go mod tidy`（可通过`--skipgomodtidy`跳过）

编辑模板<strong>不会</strong>影响现有应用，只会影响以后执行的`init`。公开示例位于`v3/examples/`下；它们不能替代贡献者文档中所述的自动化测试套件。

---

## `tasks/` – 发布自动化

Taskfile 封装了复杂的交叉编译、版本号更新和变更日志生成流程。`internal/commands/task.go`以编程方式使用这些 Taskfile，因此<strong>CLI</strong>和<strong>CI</strong>由同一套逻辑驱动。

---

## 各部分如何交互

```d2
direction: down
CLI: wails3 CLI
Generator: internal/generator
AssetDev: assetserver（开发）
Packager: internal/packager
AppRuntime: {
  label: 应用运行时
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: 操作系统 API
}
CLI -> Generator: 构建/生成
CLI -> AssetDev: 开发
CLI -> Packager: 打包
Generator -> ApplicationPkg: 绑定
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

<em>CLI → generator → runtime</em>构成了从<strong>源代码</strong>到<strong>运行中的桌面应用</strong>的核心路径。

---

## 定位提示

| 需要了解…… | 请查看…… |
| --- | --- |
| 平台适配层 | `pkg/application/*_darwin.go`、`*_linux.go`、`*_windows.go`（window、clipboard、dialogs、systray、mainthread、events_common）。Linux cgo：`pkg/application/linux_cgo*.go`。 |
| 桥接协议 | `pkg/application/messageprocessor*.go` |
| 资源工作流 | `internal/assetserver/`（`build_dev.go`与`build_production.go`） |
| 打包流程 | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`、`internal/packager/` |
| 模板引擎 | `internal/templates/`（`templates.Install`、`templates.GetDefaultTemplates`） |
| 静态分析 | `internal/generator/{generate.go,collect/,render/}` |

---

现在，你已经对该仓库有了一份<strong>思维导图</strong>。结合`ripgrep`、IDE 的“转到文件/符号”功能和示例应用，可以进一步探索任何功能。祝你开发愉快！
