---
title: "常见问题"
description: "关于使用 Wails 构建应用程序的常见问题解答 v3"
slug: "faq"
sourcePath: "faq.md"
---

## 常规

### 什么是 Wails？

Wails 是一个使用 Go 和 Web 技术构建桌面应用程序的框架。应用程序逻辑使用 Go 编写，界面使用 HTML、CSS 和 JavaScript（或任意前端框架）构建，并由 Wails 在操作系统的原生 WebView 中渲染。最终得到的是小巧、快速且具有原生体验的应用程序：无需捆绑浏览器，内存占用低，而且通常只需一个约 10MB 的二进制文件。

### Wails 支持哪些平台？

| 平台 | 要求 |
| --- | --- |
| Windows | AMD64 和 ARM64。使用 [WebView2 运行时](https://developer.microsoft.com/microsoft-edge/webview2/)。 |
| macOS | Intel 需要 10.15 或更高版本（应用程序可以 10.13 或更高版本为目标），Apple Silicon 需要 11.0 或更高版本。支持通用二进制文件。 |
| Linux | AMD64 和 ARM64。默认技术栈为 GTK4 和 WebKitGTK 6.0（Ubuntu 24.04+、Debian 13+、Fedora 40+ 及类似发行版）。仅提供 WebKit2GTK 4.1 的发行版（例如 Ubuntu 22.04、Debian 12 和 RHEL 9）可通过旧版 `-tags gtk3` 构建获得支持（支持截至 v3.1）。仅提供 WebKit2GTK 4.0 的发行版不受支持。请参阅 [Linux 构建指南](/guides/build/linux/)。 |
| iOS 和 Android | 实验性支持。请参阅[移动端指南](/guides/mobile/)。 |

还可以使用[服务器构建](/guides/server-build/)，将应用程序作为普通 Web 应用提供。

随时运行 `wails3 doctor` 检查系统，并获取针对相应平台的安装说明。

### 开始使用需要准备什么？

- Go 1.26 或更高版本
- Node.js 和 npm（用于前端构建）
- 平台工具链：Windows 上的 WebView2（已预装在 10/11 上）、macOS 上的 Xcode Command Line Tools，以及 Linux 上的 `gcc` 和 GTK/WebKit 开发包

`wails3 doctor` 会为你检查所有这些要求，并明确指出缺少的内容。有关完整步骤，请参阅[安装](/quick-start/installation/)。

### Wails v3 是否已可用于生产环境？

Wails v3.0.0 是面向 Windows、macOS 和 Linux 桌面应用程序的稳定版本。Android 和 iOS 支持仍处于实验阶段，不在稳定桌面版兼容性约定的范围内。请参阅[项目状态页面](/status/)。

## 开发

### 我需要掌握 Go 吗？

具备基本的 Go 知识会有所帮助，但无需成为专家。应用程序逻辑位于普通的 Go 方法中，[教程](/tutorials/overview/)会引导你完成其余所有内容。许多开发者都是在构建第一个 Wails 应用程序的过程中学会 Go 的。

### 可以使用我喜欢的前端框架吗？

可以。只要能够构建为 HTML、CSS 和 JavaScript，就可以与 Wails 配合使用。Wails 提供 React、Vue、Svelte 和原生 JavaScript 模板（每种都有 TypeScript 版本），其他框架也可以在几分钟内接入。请参阅[前端框架](/guides/dev/frontend-frameworks/)。

### 如何从 JavaScript 调用 Go 函数？

注册服务后，Wails 会为其生成类型化绑定：

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

绑定会在 `wails3 dev` 期间自动重新生成，也可以使用 `wails3 generate bindings` 按需生成。请参阅[服务](/features/bindings/services/)。

### 可以使用 TypeScript 吗？

可以。绑定生成器会为服务及其类型生成 TypeScript 定义，因此对 Go 的调用具有完整的类型信息。

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

运行 `wails3 dev`，然后在窗口中右键单击即可打开浏览器开发者工具，就像在 Web 页面中一样。开发服务器还支持前端热重载。请参阅[调试](/guides/dev/debugging/)。

## 构建与分发

### 如何构建生产版本？

```bash
wails3 build
```

二进制文件会生成到 `bin/`。生产构建已经应用合理的默认设置（构建标签、`-trimpath`、剥离符号），因此无需额外标志即可获得精简的二进制文件。

### 可以交叉编译吗？

可以，但有限制。由于每个平台都使用原生 WebView 库，因此无法直接采用纯 Go 的交叉编译方式，不过常见场景得到了良好支持：

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

从其他操作系统为 Linux 构建时，会使用基于 Docker 的工具链。有关完整支持矩阵，请参阅[跨平台构建](/guides/build/cross-platform/)。

### 如何创建安装程序或软件包？

```bash
wails3 package
```

这会生成平台原生格式；[安装程序指南](/guides/installers/)介绍了 Windows 上的 NSIS、macOS 上的 `.app` 包和 DMG，以及 Linux 软件包。

### 如何对应用程序进行代码签名？

[签名指南](/guides/build/signing/)逐步介绍了 Windows 和 macOS 签名，包括公证。

## 功能

### 可以创建多个窗口吗？

可以，v3 原生支持多窗口：

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

支持。Wails v3 内置自更新程序（`app.Updater`），可插拔地支持 GitHub Releases、keygen.sh 和 Sparkle AppCast 提供程序，并提供加密签名验证以及可自定义主题或替换的默认界面。请参阅[应用内更新程序](/guides/updater/)指南和[可自更新的 Wails 应用程序](/tutorials/04-self-update-a-wails-app/)教程。

## 故障排除

### 某些功能无法正常工作，应从哪里开始排查？

```bash
wails3 doctor
```

它会验证你的工具链，列出缺失的依赖项及其安装命令，并输出你应在所有错误报告中附上的版本信息。

### 构建失败

请按以下顺序尝试常见的解决方法：

1. `go mod tidy`
2. `cd frontend && npm install`（最常见的原因是缺少 `node_modules`）
3. 更新 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. 在 Linux 上，检查 `wails3 doctor` 是否报告缺少 GTK/WebKit 软件包

### 我的绑定缺失或已过期

```bash
wails3 generate bindings
```

在开发模式下，绑定会自动重新生成；如果你在`wails3 dev`之外添加了新服务或更改了方法签名，请手动重新生成绑定。

### 事件未触发

Go 中`app.Event.Emit("name", ...)`的事件名称必须与 JavaScript 中`Events.On("name", ...)`的事件名称完全一致。请先检查拼写错误和大小写差异。

### 我发现了一个 bug

请[提交 issue](https://github.com/wailsapp/wails/issues)，并附上`wails3 doctor`的输出。[反馈指南](/feedback/)说明了如何编写便于处理的报告。

## 从v2迁移

### 我应该从v2迁移到v3吗？

v3带来了多窗口支持、更简洁的基于服务的 API、内置更新程序、灵活得多的构建系统以及更高的性能。新项目应使用v3。对于现有项目，[迁移指南](/migration/v2-to-v3/)会逐步介绍二者之间的差异。

### v2会继续维护吗？

会。Wails v2会继续获得关键更新。新项目应使用 Wails v3.0.0。

### 我可以同时运行v2和v3吗？

可以。这两个 CLI 是独立的二进制文件（`wails`和`wails3`），模块的导入路径也不同，因此使用不同主版本的项目可以在同一台计算机上共存。

## 社区

### 如何获得帮助？

- 在[Discord](https://discord.gg/JDdSxwjhGf)中快速提问和讨论
- 在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)中提出需要深入讨论的问题
- 在[GitHub Issues](https://github.com/wailsapp/wails/issues)中报告 bug

### 如何参与贡献？

请参阅[贡献指南](/contributing/)。随时欢迎提交 bug 修复；新增功能和公共行为变更需通过[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)草案 PR 提出。你可以选择在 Discord 或 GitHub Discussions 上进行非正式讨论。

### 在哪里可以找到示例？

仓库附带超过60个可运行示例，涵盖窗口、对话框、事件、系统托盘、服务等内容：[v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。

## 还有其他问题？

请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或[发起讨论](https://github.com/wailsapp/wails/discussions)。
