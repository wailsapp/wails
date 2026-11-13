---
title: "Wails v2 发布"
description: "Wails 的发行说明与公告"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![蒙太奇截图](/assets/blog-images/montage.png)

## 它来了！

今天，[Wails](https://wails.io) v2 正式发布。距离首个 v2 alpha 版本发布已过去约 18 个月，距离首个 beta 版本发布也已有约一年。我由衷感谢参与项目发展历程的每一个人。

之所以耗时这么久，部分原因是我们希望先达到某种意义上的完整状态，再正式将其称为 v2。事实上，标记发行版本从来没有完美的时机——总会有尚未解决的问题，或是想再塞进去“最后一个”功能。不过，即使是不完美的主版本，正式标记并发布也能为项目用户带来一定的稳定性，同时让开发者获得一个重新出发的契机。

这个版本远远超出了我最初的预期。希望它能为你带来和我们开发它时同样多的乐趣。

## Wails *是*什么？

如果你还不了解 Wails，它是一个让 Go 程序员能够使用熟悉的 Web 技术为 Go 程序构建丰富前端的项目。它是 Electron 的轻量级 Go 替代方案。你可以在[官方网站](https://v2.wails.io/docs/introduction)上找到更多信息。

## 有哪些新功能？

v2 版本是项目的一次巨大飞跃，解决了 v1 的许多痛点。如果你还没有读过关于 [macOS](/blog/wails-v2-beta-for-mac/)、[Windows](/blog/wails-v2-beta-for-windows/) 或 [Linux](/blog/wails-v2-beta-for-linux/) beta 版本的博客文章，建议阅读一下，其中更详细地介绍了所有重大变更。简而言之：

- 适用于 Windows 的 Webview2 组件，支持现代 Web 标准和调试功能。
- Windows 上的[深色/浅色主题](https://v2.wails.io/docs/reference/options#theme)和[自定义主题](https://v2.wails.io/docs/reference/options#customtheme)。
- Windows 现在不再需要 CGO。
- 开箱即用地支持 Svelte、Vue、React、Preact、Lit 和 Vanilla 项目模板。
- 集成 [Vite](https://vitejs.dev/)，为应用提供支持热重载的开发环境。
- 原生应用[菜单](https://v2.wails.io/docs/guides/application-development#application-menu)和[对话框](https://v2.wails.io/docs/reference/runtime/dialog)。
- 为 [Windows](https://v2.wails.io/docs/reference/options#windowistranslucent) 和 [macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1) 提供原生窗口半透明效果。支持 Mica 和 Acrylic 背景材质。
- 可轻松为 Windows 部署生成 [NSIS 安装程序](https://v2.wails.io/docs/guides/windows-installer)。
- 功能丰富的[运行时库](https://v2.wails.io/docs/reference/runtime/intro)，为窗口操作、事件处理、对话框、菜单和日志记录提供实用方法。
- 支持使用 [garble](https://github.com/burrowers/garble) [混淆](https://v2.wails.io/docs/guides/obfuscated)应用。
- 支持使用 [UPX](https://upx.github.io/) 压缩应用。
- 自动根据 Go 结构体生成 TypeScript。更多信息请参阅[此处](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods)。
- 在任何平台上，都不需要随应用分发额外的库或 DLL。
- 无需打包前端资源。只需像开发其他 Web 应用一样开发你的应用即可。

## 致谢

走到 v2 这一步，我们付出了巨大的努力。从最初的 alpha 版本到今天正式发布，共有 89 位贡献者提交了约 2200 次 commit；此外，还有许许多多的人贡献了翻译、测试和反馈，并在讨论论坛及问题跟踪器中提供帮助。我对你们每个人都感激不尽。我还要特别感谢所有项目赞助者提供的指导、建议和反馈。你们所做的一切都弥足珍贵。

我想特别提到几个人：

首先，要向 [@stffabi](https://github.com/stffabi) 表示**万分**感谢。他做出了众多惠及所有人的贡献，也为许多问题提供了大量支持。他实现了一些关键功能，例如外部开发服务器支持；这项功能让我们能够利用 [Vite](https://vitejs.dev/) 的强大能力，彻底改变了开发模式的体验。可以毫不夸张地说，如果没有他的[杰出贡献](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04)，Wails v2 版本绝不会如此令人振奋。衷心感谢你，@stffabi！

我还要特别感谢 [@misitebao](https://github.com/misitebao)。他一直不辞辛劳地维护网站、提供中文翻译、管理 Crowdin，并帮助新译者快速上手。这项工作极其重要，我对他为此投入的所有时间和精力深表感激！你太棒了！

最后同样重要的是，非常感谢 Mat Ryer 在 v2 开发期间提供的建议和支持。我们曾使用 v2 的早期 alpha 版本共同开发 xBar，这段经历帮助确定了 v2 的发展方向，也让我认识到早期版本中的一些设计缺陷。我很高兴地宣布，从今天起，我们将开始把 xBar 移植到 Wails v2，它也将成为本项目的旗舰应用。谢谢你，Mat！

## 经验与教训

在实现 v2 的过程中，我们获得了许多经验与教训，它们将影响今后的开发方向。

## 更小、更快、更聚焦的版本

在开发 v2 的过程中，许多功能和错误修复都是临时开发的。这导致发布周期变长，调试也更加困难。今后，我们会更频繁地发布版本，并减少每个版本包含的功能数量。每次发布都会同步更新文档并进行全面测试。希望这些更小、更快、更聚焦的版本能够减少回归问题，并带来更高质量的文档。

## 鼓励参与

刚开始这个项目时，我想立即帮助每一个遇到问题的人。我把每个问题都视为“自己的事”，希望尽快解决。这种做法难以为继，最终还会损害项目的长期发展。今后，我会留出更多空间，让大家参与回答问题和分类处理 issue。如果能有一些工具来协助这项工作就更好了。如果你有任何建议，请在[这里](https://github.com/wailsapp/wails/discussions/1855)参与讨论。

## 学会说“不”

参与开源项目的人越多，要求增加新功能的请求就越多，而这些功能对大多数人可能有用，也可能没用。这些功能不仅需要投入前期时间进行开发和调试，从此还会带来持续的维护成本。在这方面，我自己最需要反省：我经常想要“一口吃成胖子”，而不是先提供最小可行功能。今后，对于向核心添加功能的请求，我们需要更多地说“不”，并集中精力探索如何让开发者自行实现这些功能。针对这种场景，我们正在认真研究插件方案。这样，任何人都能按自己的需求扩展项目，同时也能轻松地为项目作出贡献。

## 展望未来

下一个主要开发周期中，我们已经在考虑为 Wails 添加大量核心功能。[路线图](https://github.com/wailsapp/wails/discussions/1484)中充满了有趣的想法，我迫不及待地想开始实现它们。其中呼声很高的一项功能是多窗口支持。这项功能颇具挑战；为了妥善实现它，我们可能需要考虑提供替代 API，因为当前 API 在设计时并未考虑这一需求。根据一些初步构想和反馈，我认为你会喜欢我们正在探索的方向。

我个人对让 Wails 应用在移动设备上运行的前景感到非常兴奋。我们已经有一个演示项目，证明 Wails 应用可以在 Android 上运行，因此我非常期待探索这项能力未来能发展到什么程度！

最后我想谈的一点是功能对等。长期以来，我们一直坚持一项核心原则：只有在获得完整跨平台支持的情况下，才会向项目添加功能。虽然到目前为止，这一点（基本上）都能做到，但它确实阻碍了项目发布新功能。今后，我们将采用一种略有不同的方式：任何无法立即在所有平台上发布的新功能，都将通过实验性配置或 API 发布。这样，特定平台上的早期采用者就可以试用该功能并提供反馈，这些反馈将用于完善功能的最终设计。当然，这意味着在该功能得到所有能够支持它的平台的完整支持之前，无法保证 API 的稳定性，但至少不会再阻碍开发进程。

## 结语

我对我们在V2版本中取得的成果深感自豪。看到人们使用目前发布的 Beta 版本构建出的成果，实在令人惊叹，例如[Varly](https://varly.app/)、[Surge](https://getsurge.io/)和[October](https://october.utf9k.net/)等高质量应用。建议你去了解一下。

这个版本凝聚了众多贡献者的辛勤付出。虽然它可以免费下载和使用，但绝非毫无成本。毫无疑问，这个项目付出了相当大的代价。这不仅包括我和每一位贡献者投入的时间，也包括大家因此无法陪伴亲友所付出的代价。正因如此，我由衷感谢为实现这个项目所投入的每一分每一秒。贡献者越多，这份工作就越能由更多人分担，我们共同取得的成果也会越多。我想鼓励大家选择一件力所能及的事情来贡献，无论是确认他人报告的 bug、提出修复建议、修改文档，还是帮助有需要的人。所有这些小事都能产生巨大的影响！如果你也能参与推动项目迈向v3的历程，那就太棒了。

尽情享用吧！

&dash; Lea

附言：如果你或你的公司认为 Wails 很有用，请考虑[赞助本项目](https://github.com/sponsors/leaanthony)。谢谢！
