---
title: "常见问题"
description: "有关使用 Wails v3 构建应用程序的常见问题解答"
slug: "faq"
sourcePath: "faq.md"
---

## 常规

### 什么是 Wails？

Wails 是一个使用 Go 和 Web 技术构建桌面应用程序的框架。应用程序逻辑使用 Go 编写，界面使用 HTML、CSS 和 JavaScript（或任意前端框架）构建，Wails 则通过操作系统的原生 WebView 进行渲染。最终得到的是一个小巧、快速且具有原生体验的应用程序：无需捆绑浏览器、内存占用低，并且通常只需一个约 10MB 的二进制文件。

### Wails 支持哪些平台？

| 平台 | 要求 |
| --- | --- |
| Windows | AMD64 和 ARM64。使用[WebView2 运行时](https://developer.microsoft.com/microsoft-edge/webview2/)。 |
| macOS | Intel 上要求10.15或更高版本（应用程序可将10.13或更高版本设为目标），Apple Silicon 上要求11.0或更高版本。支持通用二进制文件。 |
| Linux | AMD64 和 ARM64。默认技术栈为 GTK4 和 WebKitGTK 6.0（Ubuntu 24.04或更高版本、Debian 13或更高版本、Fedora 40或更高版本以及类似发行版）。对于仅提供 WebKit2GTK 4.1的发行版，例如 Ubuntu 22.04、Debian 12和 RHEL 9，可通过旧版`-tags gtk3`构建提供支持（支持至 v3.1）。不支持仅提供 WebKit2GTK 4.0的发行版。请参阅[Linux 构建指南](/guides/build/linux/)。 |
| iOS 和 Android | 实验性支持。请参阅[移动端指南](/guides/mobile/)。 |

你还可以使用[服务器构建](/guides/server-build/)将应用程序作为常规 Web 应用提供。

随时运行`wails3 doctor`即可检查系统并获取针对当前平台的安装说明。

### 开始使用需要准备什么？

- Go 1.25或更高版本
- Node.js 和 npm（用于前端构建）
- 平台工具链：Windows 上的 WebView2（已预装在10/11上）、macOS 上的 Xcode Command Line Tools，以及 Linux 上的`gcc`和 GTK/WebKit 开发包

`wails3 doctor`会为你检查所有这些要求，并明确指出缺少哪些组件。完整步骤请参阅[安装](/quick-start/installation/)。

### Wails v3 是否已可用于生产环境？

Wails v3 是测试版软件，但桌面 API 已保持稳定。已有应用程序使用它在生产环境中运行，但在我们完成3.0的最后完善工作期间，你仍应在部署前进行全面测试。有关当前状况，请参阅[项目状态页面](/status/)。Wails v2 是当前的稳定版本，并会继续获得修复。

## 开发

### 我需要掌握 Go 吗？

具备基本的 Go 知识会有所帮助，但你不必成为专家。应用程序逻辑由普通的 Go 方法实现，[教程](/tutorials/overview/)会引导你完成其他所有内容。许多开发者都是在构建第一个 Wails 应用程序的过程中学会 Go 的。

### 我可以使用自己喜欢的前端框架吗？

可以。只要能构建为 HTML、CSS 和 JavaScript，就能与 Wails 配合使用。Wails 随附 React、Vue、Svelte 和原生 JavaScript 模板（每种模板均有 TypeScript 变体），其他框架也可以在几分钟内完成接入。请参阅[前端框架](/guides/dev/frontend-frameworks/)。

### 如何从 JavaScript 调用 Go 函数？

注册一个服务，Wails 就会为其生成类型化绑定：

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

在`wails3 dev`期间，绑定会自动重新生成；也可以使用`wails3 generate bindings`按需生成。请参阅[服务](/features/bindings/services/)。

### 我可以使用 TypeScript 吗？

可以。绑定生成器会为你的服务及其类型生成 TypeScript 定义，因此对 Go 的调用具有完整的类型信息。

### 如何在 Go 和 JavaScript 之间发送事件？

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

事件名称必须完全匹配。请参阅[事件参考](/guides/events-reference/)。

### 如何调试应用程序？

运行`wails3 dev`，然后在窗口中单击右键打开浏览器开发者工具，就像在 Web 上操作一样。开发服务器还支持前端热重载。请参阅[调试](/guides/dev/debugging/)。

## 构建与分发

### 如何构建生产版本？

```bash
wails3 build
```

生成的二进制文件位于`bin/`中。生产构建已采用合理的默认设置（构建标签、`-trimpath`、移除符号），因此无需额外标志即可获得精简的二进制文件。

### 可以交叉编译吗？

可以，但存在限制。由于每个平台都使用原生 WebView 库，因此无法直接采用纯 Go 交叉编译，不过常见场景都得到了良好支持：

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

从其他操作系统为 Linux 构建时，会使用基于 Docker 的工具链。有关完整的支持矩阵，请参阅[跨平台构建](/guides/build/cross-platform/)。

### 如何创建安装程序或软件包？

```bash
wails3 package
```

这会生成平台原生格式；[安装程序指南](/guides/installers/)介绍了 Windows 上的 NSIS、macOS 上的`.app`捆绑包和 DMG，以及 Linux 软件包。

### 如何对应用程序进行代码签名？

[签名指南](/guides/build/signing/)分步骤介绍了 Windows 和 macOS 签名，包括公证。

## 功能

### 我可以创建多个窗口吗？

是的，v3 原生支持多窗口：

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

请参阅[多窗口](/features/windows/multiple/)。

### Wails 支持系统托盘吗？

支持，包括菜单和点击处理程序：

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

请参阅[系统托盘](/features/menus/systray/)。

### 可以使用原生对话框吗？

可以。文件对话框、消息对话框和询问对话框均使用原生实现：

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

请参阅[对话框](/features/dialogs/overview/)。

### Wails 支持自动更新吗？

支持。Wails v3 内置了自更新程序（`app.Updater`），提供适用于 GitHub Releases、keygen.sh 和 Sparkle AppCast 的可插拔提供程序、加密签名验证，以及可自定义主题或替换的默认 UI。请参阅[应用内更新程序](/guides/updater/)指南和[可自更新的 Wails 应用](/tutorials/04-self-update-a-wails-app/)教程。

## 故障排除

### 遇到功能异常时，应从哪里着手？

```bash
wails3 doctor
```

它会验证工具链，列出缺失的依赖项及其安装命令，并输出应随错误报告一并提供的版本信息。

### 构建失败

请依次尝试以下常见解决方法：

1. `go mod tidy`
2. `cd frontend && npm install`（最常见的原因是缺少`node_modules`）
3. 更新 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. 在 Linux 上，检查`wails3 doctor`是否报告缺少 GTK/WebKit 软件包

### 绑定缺失或已过时

```bash
wails3 generate bindings
```

在开发模式下，绑定会自动重新生成；如果您在`wails3 dev`之外添加了新服务或更改了方法签名，请手动重新生成绑定。

### 事件未触发

Go 中`app.Event.Emit("name", ...)`使用的事件名称必须与 JavaScript 中`Events.On("name", ...)`使用的事件名称完全一致。请先检查拼写错误和大小写差异。

### 发现了错误

请[创建 issue](https://github.com/wailsapp/wails/issues)，并附上`wails3 doctor`的输出。[反馈指南](/feedback/)说明了如何提交便于处理的报告。

## 从 v2 迁移

### 应该从 v2 迁移到 v3 吗？

v3 提供多窗口支持、更简洁的基于服务的 API、内置更新程序、灵活得多的构建系统以及更出色的性能。新项目应从 v3 开始。对于现有项目，[迁移指南](/migration/v2-to-v3/)会逐步介绍各项差异。

### v2 会继续维护吗？

会。在 v3 迈向稳定版本的同时，v2 会继续获得修复。

### 可以同时运行 v2 和 v3 吗？

可以。这两个 CLI 是独立的二进制文件（`wails`和`wails3`），模块的导入路径也不同，因此采用不同主版本的项目可以在同一台计算机上顺利共存。

## 社区

### 如何获取帮助？

- 如需快速提问和讨论，请使用[Discord](https://discord.gg/JDdSxwjhGf)
- 如需提出需要详细说明的问题，请使用[GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- 如需报告错误，请使用[GitHub Issues](https://github.com/wailsapp/wails/issues)

### 如何贡献？

请参阅[贡献指南](/contributing/)。随时欢迎提交错误修复；新增功能和公共行为变更需通过[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)草案 PR 提交。可选择在 Discord 或 GitHub Discussions 上进行非正式讨论。

### 在哪里可以找到示例？

代码仓库随附超过60个可运行示例，涵盖窗口、对话框、事件、系统托盘、服务等内容：[v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples)。

## 仍有疑问？

请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或[发起讨论](https://github.com/wailsapp/wails/discussions)。
