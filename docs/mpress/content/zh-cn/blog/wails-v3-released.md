---
title: "Wails v3.0.0 正式发布"
description: "Wails v3.0.0 已正式推出：支持多窗口、系统托盘、原生菜单、服务和生成式绑定、更新程序、移动端，并启用了全新的文档网站。"
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Wails v3.0.0 发布影片——全屏观看效果最佳" width="1920" height="1080"}

## 它来了！

今天，我要发布 Wails v3.0.0。

坦白说，这是一段漫长的旅程。首个 v3 alpha 标签于 2023 年 1 月发布，从那以后，许多人一路参与构建、测试、发现问题并修复问题。我由衷感谢你们每一个人。如果你想了解这段旅程不加修饰的版本，可以阅读我写的 [v3 beta 发布公告](/blog/wails-v3-beta/)。

## Wails *是*什么？

如果你刚接触 Wails：它让 Go 程序员能够使用自己熟悉的 Web 前端工具构建桌面应用，并采用各平台的原生 WebView，而不是嵌入式浏览器。你可以使用 Go 编写后端，使用 HTML、CSS 和 JavaScript（或任何你喜欢的框架）编写 UI，再由 Wails 将这一切打包成单个原生应用。

## 为什么推出 v3？

Wails v2 让使用现代 Web 前端构建 Go 应用变得简单直接，也一直很好地服务于众多应用。但它的单窗口、上下文驱动运行时，使一些日常桌面开发工作变得比应有的更困难。

v3 采用了不同的模型。应用、窗口、服务、事件和平台能力都是显式对象。你可以创建应用、注册服务、创建窗口，并直接与负责所需行为的对象交互。随着应用不断发展，这种模型更容易理解；多窗口等功能也成为模型的常规组成部分，而不再是权宜之计。

## v3.0.0 包含哪些内容

- <strong>多窗口。</strong>在运行时创建、管理和关闭窗口，每个窗口都拥有自己的配置和生命周期。
- <strong>系统托盘。</strong>让应用常驻系统托盘，支持菜单、浅色和深色图标、macOS 模板图标，还可以将窗口附加到托盘图标。
- <strong>原生菜单和对话框。</strong>在各个平台上提供应用菜单、上下文菜单和原生对话框。
- <strong>服务和生成式绑定。</strong>后端仍是普通的 Go 代码。绑定通过静态源代码分析生成，因此只需一条 `wails3 generate bindings` 命令，就能将注释和有意义的参数名称一并传递到前端。
- <strong>清晰可见的构建系统。</strong>项目使用基于 [Taskfile](https://taskfile.dev) 的构建系统，你可以检查、扩展和调试它，而不再依赖一条不透明的命令。
- <strong>跨平台。</strong>支持 Windows（amd64 和 arm64）、macOS（Intel 和 Apple Silicon）以及 Linux（amd64 和 arm64）；Linux 默认使用 GTK4 和 WebKitGTK 6.0。在 v3.1 之前，仍可通过 `-tags gtk3` 使用旧版 GTK3 技术栈。
- <strong>iOS 和 Android。</strong>同一套 `main.go` 和前端可以在移动设备上运行。文档仍将移动端支持标记为实验性功能，因此请将其视为可供探索和反馈的功能，而非稳定承诺。
- <strong>使用 Go 1.26 构建。</strong>v3 模块使用 Go 1.26 构建。
- <strong>应用内更新程序。</strong>支持可插拔提供程序、加密验证、原子替换，以及可自定义主题或替换的默认 UI。请参阅[更新程序指南](/guides/updater/)。
- <strong>打包和签名。</strong>提供各个平台的[代码签名](/guides/build/signing/)、[安装程序](/guides/installers/)和[MSIX 打包](/guides/build/msix/)指南。
- <strong>服务器构建。</strong>无需原生窗口，即可将同一应用和服务作为 HTTP 服务器运行。

## Beta 版发布以来的亮点

Beta 版文章于 2 月发布，此后变更日志一直在不断更新。以下是我想特别介绍的几项内容：

- <strong>全新的文档网站。</strong>文档现已由我们自己的文档工具[M-Press](https://m-press.me)驱动，译文也已完成并从社区迁移过来。英文文档仍是唯一权威来源。
- <strong>更稳健的 Windows 使用体验。</strong>Wails 现在能从 WebView2 进程故障中恢复，不会再留下空白窗口；同时还修复了多个 WebView2 COM 生命周期问题。
- <strong>尊重键盘布局的全局快捷键。</strong>在 macOS 上，字母快捷键现在会遵循 AZERTY、QWERTZ 和 Dvorak 等布局中标示的按键。
- <strong>增强更新程序的稳健性。</strong>其中包括修复在 Linux 上跨文件系统替换二进制文件以及 AppImage 更新等问题。
- **设置和工具。**`wails3 setup` 向导现已完成；应用退出时，`wails3 dev` 现在也会停止其后台进程。
- <strong>macOS 原生 Liquid Glass。</strong>支持 macOS 26 及更高版本中的 Liquid Glass 效果，并为更早版本提供回退方案。

## 从 v2 迁移

v3 是一个新的主版本，迁移需要真正进行移植，而不只是升级版本号。主要变化包括应用和窗口的生命周期、以服务取代与上下文绑定的绑定、以直接的应用和窗口 API 取代 v2 运行时包，以及重新生成前端绑定。

[v2 到 v3 的迁移指南](/migration/v2-to-v3/)通过迁移前后的示例，逐一讲解所有这些变化。请有计划地进行移植，在过程中持续测试，并在 v3 应用准备就绪之前保留 v2。

<strong>Wails v2 不会消失。</strong>它会继续获得关键更新。新项目请使用 v3.0.0。这两个 CLI 是独立的二进制文件（`wails` 和 `wails3`），模块导入路径也不同，因此可以同时安装在同一台计算机上。

## 安装

使用以下命令安装 v3.0.0 CLI：

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

然后运行引导式设置向导。它会检查你的环境，并帮助处理 Wails 所需的依赖项：

```sh
wails3 setup
```

![深色模式下的 Wails v3 设置向导](/assets/screenshots/wails3-setup-wizard-dark.png)

接着创建你的第一个项目：

```sh
wails3 init
```

- 文档：<https://v3.wails.io/>
- 迁移指南：<https://v3.wails.io/migration/v2-to-v3/>
- 发行说明：GitHub 上的[Wails v3.0.0](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## 支持和发布

Wails 的发展规模远超我的预期，要跟上大量议题和拉取请求一直很困难。有时发布缺乏固定安排，修复和社区贡献也要等待太久才能送达用户。我们已经建立了一些系统，以便将议题分类和发布准备中的更多基础工作自动化，从而更快地响应问题，并投入更多时间交付改进和新功能。

v3.0.0 之后，我们将推出带日期标签（如 `v3.1.0-nightly.20261026`）的**每夜预发布版**，同时发布**每周次版本**，例如 `v3.1.0` 和 `v3.2.0`。每夜版可以让你更早获得近期修复，而每周发布则会定期提供经过测试的更新。

如果发现可复现的错误，请附上 `wails3 doctor` 的输出和最小示例进行报告。如果想提议新增功能或更改公开行为，请创建一份 WEP（Wails Enhancement Proposal，Wails 增强提案）草稿 PR，而不是功能请求议题。

## 翻译

向所有帮助翻译 Wails v2 文档的人致以最衷心的感谢。你们付出了自己的时间，让开发者能够用母语阅读文档，这对许多人意义重大。我们深表感激。

我们发现 Crowdin 尤其难以管理，因此已停止使用。从 v3 开始，文档将自动翻译，译文会与英文页面一同存放在代码仓库中。

如果你希望我们提供尚未支持的语言译文，或发现现有译文中有错误，请提交 issue 告诉我们。我们非常期待你的反馈，也由衷感谢每一份报告。

## 致谢

Wails v3 的诞生离不开所有测试未完成版本、报告棘手 bug、翻译文档、解答问题和贡献代码的人。谢谢大家。完整的贡献者名单位于[致谢页面](/credits/)，诚邀你前往查看。

首先，要特别感谢 Wails 维护者[Atterpac](https://github.com/atterpac)。他在开发工具、Linux 桌面集成和社区支持方面所做的工作，帮助 v3 走到了今天。

特别感谢[Fabio Massaioli (fbbdev)](https://github.com/fbbdev)打造了如此出色的静态源代码分析器和绑定生成器。它是 v3 得以运行的关键组成部分：普通的 Go 服务可以转化为 JavaScript 和 TypeScript API，并将类型和文档一并传递到前端。

还要感谢[Travis McLane](https://github.com/tmclane)在早期开展的大量平台工作；感谢[Stffabi](https://github.com/stffabi)奠定了资源服务器和原生平台基础；感谢[Zach Botterman (popaprozac)](https://github.com/popaprozac)实现了通知、Dock 集成及其他桌面功能；感谢[Ian VanSchooten (IanVS)](https://github.com/IanVS)实现了类型化事件并改进了运行时；也感谢[Lyimmi](https://github.com/Lyimmi)在 Linux 方面所做的工作。这些贡献，以及所有帮助构建和完善它的人所做的工作，至今仍是 v3 的一部分。

向在 Wails 漫长转型期间一路支持它的赞助者致以最衷心的感谢。你们的支持让我们有时间完善架构、工具、测试和文档，最终走到今天。还要特别感谢 JetBrains 和 ZSA，以及通过代码审查支持本项目的[CodeRabbit](https://www.coderabbit.ai/oss)。

最后，特别感谢 Valerio Tesei 多年来对本项目给予的全力支持。

## 最后的话

我们已经有大量工作蓄势待发。整个科技领域正涌现出一些令人振奋的发展，而我们也已准备好探索新的范式，为 Wails 带来全新的功能。但首先，我们需要让 v3 成为一个稳定、值得你信赖的基础。先把这一点做好，再推出后续内容。

我真的为 v3 感到自豪。希望你使用它进行开发时，能像我们打造它时一样享受其中。试试它，构建些东西，然后告诉我们你的想法。欢迎到 [Discord](https://discord.gg/JDdSxwjhGf)、[GitHub Discussions](https://github.com/wailsapp/wails/discussions) 或[subreddit](https://www.reddit.com/r/wails/) 上来打个招呼。

尽情享受吧！

&dash;&nbsp;Lea

附言：如果你或你的公司认为 Wails 很有用，请考虑[赞助本项目](https://github.com/sponsors/leaanthony)。谢谢！
