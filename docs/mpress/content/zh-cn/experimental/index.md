---
title: "实验性功能"
description: "这里汇集了 Wails v3 中正在推进的实验，并介绍它们是什么、为何存在以及如何提供反馈。"
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="这里有实验性功能"}
本节中的所有内容从定义上来说都是实验性的。这些功能需要主动启用，默认处于关闭状态，并且可能在不同版本之间改变形式、重命名，甚至被完全移除。如果尚未做好应对频繁变更的准备，请勿将它们用于关键功能。

@end

## “实验性”意味着什么

Wails 公开发布实验性功能。我们认为某个想法很有前景，值得尽早交到你手中，但尚未完全承诺长期支持时，就会将其作为实验性功能发布。我们之所以发布它，*是因为*希望先从实际使用中学习，再决定是否将其转为 Wails 中永久受支持的组成部分。

这意味着本节中的所有内容都具有以下特点：

- <strong>需要主动启用。</strong>实验性功能绝不会改变`wails3`的默认行为。你需要有意启用它们（通常通过环境变量或构建标志）；关闭它们后，现有工作流不会发生任何变化。
- <strong>它可能不会保留下来。</strong>有些实验性功能会成为稳定功能，另一些则可能被彻底重做或弃用。相比只发布我们已完全确定的功能，我们更愿意公开尝试并快速学习。
- <strong>API 尚未冻结。</strong>在实验性功能逐步成形的过程中，其名称、标志、默认值和行为可能随版本变化。发布说明会明确指出这些变更，但不要期待稳定功能所具备的稳定性保证。

## 我们希望收到你的反馈

这一点至关重要。实验性功能能否继续存在，取决于我们从实际使用者那里收到的反馈。如果你尝试了其中某项功能，我们确实希望了解：

- 它是否适用于你的项目？哪里出了问题？
- 它是否更快、更清晰或更好用，还是不值得切换？
- 要满足哪些条件，你才会默认使用它？

最有用的反馈应当具体说明：你运行了什么、预期结果是什么，以及实际发生了什么。每项实验性功能在 GitHub Discussions 的**实验性功能**类别中都有自己的讨论帖：

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="实验性功能讨论" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="找到你所用实验性功能的讨论帖，告诉我们使用效果如何、哪里出了问题或还缺少什么。"}
@end

## 当前实验性功能

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="试用原生 YAML 构建系统：在一个向导中创建并配置项目、构建并打包项目，或审查从 Taskfile 进行的单向迁移。"}
@linkcard{title="无 CGO 构建" href="/experimental/#cgo-free-builds-purego" description="使用 CGO_ENABLED=0 构建 macOS 和 Linux 应用：无需 C 工具链，也无需 Docker 即可交叉编译。可通过 WAILS_PUREGO=1 或 purego 构建标签启用。"}
@linkcard{title="LLM 控制（MCP）" href="/guides/mcp-service/" description="内置的 Model Context Protocol 服务器，让 LLM 智能体能够检查、测试和操控正在运行的 Wails 应用；无需用户代码，通过构建标签启用。"}
@end

## 无 CGO 构建（purego）

@note{type="caution" title="实验性功能"}
无 CGO 后端是由`purego`构建标签（或`WAILS_PUREGO=1`）控制的可选功能，**并非**默认后端。不启用它时，macOS 和 Linux 构建仍会像以前一样使用 cgo 后端。其行为可能随版本变化。

@end

### 它是什么

在 macOS 和 Linux 上，Wails 通常通过 cgo 与系统 UI 库通信（macOS 上为 AppKit 和 WKWebView，Linux 上为 GTK 和 WebKitGTK）。无 CGO 后端使用[purego](https://github.com/ebitengine/purego)从纯 Go 发起相同调用，并在运行时加载系统库。

公开的应用 API 保持不变。构建时选择后端即可，无需修改应用代码。分发应用前，请在目标系统上进行测试。

你将获得以下优势：

- <strong>Go 编译无需 C 工具链。</strong>Go 构建不需要 Xcode 命令行工具、`gcc`或 Linux 开发头文件。前端构建仍然需要其常规 JavaScript 工具链。
- <strong>无需 Docker 即可交叉编译。</strong>在任何操作系统上使用普通的`go build`即可构建 macOS 或 Linux 二进制文件。

Windows 已经可以在不使用 cgo 的情况下构建，因此不受影响。

### 启用前须知

请确保 CLI 和框架版本都包含此实验性功能。仅更新 CLI 不会将该后端添加到应用`go.mod`中较旧的框架依赖项。要试用开发版本，请在应用目录中运行以下命令：

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

这两个命令都会跟踪不断变化的`master`分支。若要实现可复现构建，请让 CLI 和框架使用同一个经过测试的提交。有关 Go 工具链和 CLI 的设置，请参阅[安装指南](/getting-started/installation/)。

实验性的 macOS 和 Linux 后端支持`amd64`和`arm64`。交叉编译会为指定目标创建可执行文件，但不会在构建机器上运行该文件。请在目标操作系统和架构上进行测试。

### 使用方法

请选择以下任一方式来指定后端。

#### 使用 Wails CLI

在使用此版本 CLI 生成的项目中，启用`WAILS_PUREGO`后，`build`、`package`和`dev`会使用实验性后端。在 POSIX shell 中：

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

在 PowerShell 中：

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

CLI 可识别`1`、`true`、`on`和`yes`，并忽略大小写及首尾空白。空值不会启用此实验性功能。

也可以将该标签传递给单次构建或打包命令：

```sh
wails3 build -tags purego
wails3 package -tags purego
```

设置该标签后，更新后的标准 Taskfile 会使用主机上的 Go 工具链，通过`CGO_ENABLED=0`进行构建。例如，在 macOS 或 Windows 上构建 Linux ARM64 可执行文件：

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

旧项目会保留其现有 Taskfile。请将自定义构建任务与更新后的模板进行比较，或使用普通 Go 测试该后端，而无需更改这些任务。在 Wails CLI 命令中，环境变量会选择该标签；普通的`go build`则需要显式指定标签。

#### 使用普通 Go

请先构建前端，确保`go:embed`引用的所有资源都已存在，然后编译应用：

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego`也可以使用：后端保持不变，只是程序的其余部分仍可使用 cgo。

进行交叉构建时，还需设置`GOOS`和`GOARCH`：

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

在 Linux 上，使用`-tags purego,gtk3,production`选择 GTK3 技术栈。

要恢复默认后端，请移除`purego`标签并取消设置`WAILS_PUREGO`（在 POSIX shell 中使用`unset WAILS_PUREGO`，在 PowerShell 中使用`Remove-Item Env:WAILS_PUREGO`）。恢复使用 macOS 或 Linux 的 cgo 后端时，还应移除显式的`CGO_ENABLED=0`设置。

如果使用`CGO_ENABLED=0`进行构建却忘记添加该标签，构建将停止，并显示一条说明修复方法的错误：

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### 运行时要求

应用启动时会加载系统库，因此与 cgo 构建一样，目标机器需要安装这些库：

| 平台 | 库 | 最低支持版本 |
|---|---|---|
| macOS | AppKit、WebKit（操作系统的一部分） | 与 cgo 后端相同 |
| Linux（默认） | GTK 4、WebKitGTK 6.0 | GTK 4.8、WebKitGTK 2.40（Debian 12） |
| Linux（`gtk3`） | GTK 3、WebKit2GTK 4.1 | GTK 3.24、WebKit2GTK 2.36（Ubuntu 22.04） |

应用启动时会查找高于这些最低版本所提供的函数，并且仅在已安装的库包含这些函数时才使用，因此同一个二进制文件既可在旧系统上运行，也可在新系统上运行。

应用仍会使用原生 UI 库，并可能动态链接 Linux 的 libc 等系统库。`CGO_ENABLED=0`并不意味着生成完全静态的可执行文件。如果缺少必需的运行时库，请在目标系统上安装匹配的 GTK/WebKit 技术栈；在构建机器上安装 C 头文件无法解决运行时库缺失的问题。

### 打包和签名

该后端消除了编译时对 C 工具链的要求。每种软件包格式仍有各自的要求。特别是，选择`purego`不会让现有的 AppImage 部署工具具备跨平台能力。请参阅[Linux 打包](/guides/build/linux/)。

交叉构建的 macOS 应用在分发前仍需进行适当的签名和公证。Windows 安装程序仍需使用相应的安装程序工具。请参阅[macOS 打包](/guides/build/macos/)和[Windows 打包](/guides/build/windows/)。对于使用 cgo 的依赖项，`CGO_ENABLED=1`还需要 C 编译器以及目标平台所需的所有头文件或库；选择 Wails 的纯 Go 后端并不会消除这些依赖项自身的要求。

### 验证

生成绑定检查和`puregocheck`会检查原生符号、可选版本保护机制和后端覆盖情况。Linux 一致性工作流会使用两种后端分别测试 GTK4 和 GTK3，包括资源请求、二进制请求体和原生分派；它还会交叉编译 Linux ARM64 目标。交叉编译检查用于确认代码能够构建，而运行时检查则会测试测试机器上安装的原生库。

[后端测试指南](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md)介绍了一致性探测和端到端输入测试。报告问题时，请区分编译失败与运行时或打包失败。

### 已知差异

- <strong>Linux 上的菜单快捷键采用应用级注册。</strong>请避免为不同菜单分配冲突的快捷键。使用`SetMenu(nil)`显式移除窗口菜单后，已移除菜单的快捷键也不会再在该窗口中触发；其他窗口会保留各自的菜单和快捷键。
- **macOS 打开和保存面板**与使用 cgo 时一样，会在单独的系统进程中运行；测试会检查它们能否正常打开和取消。

还发现了其他问题？请告诉我们。

### 反馈

请在[实验](https://github.com/wailsapp/wails/discussions/categories/experiments)下的**无 CGO 构建**主题中告诉我们使用情况：你的操作系统及版本、构建方式，以及哪些功能正常或未正常工作。包含`wails3 doctor`输出的失败案例，是对我们最有帮助的报告。
