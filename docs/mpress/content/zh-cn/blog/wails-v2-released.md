---
title: "Wails v2 发布"
description: "Wails 的发行说明和公告"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![蒙太奇截图](/assets/blog-images/montage.png)

## 它来了！

今天，[Wails](https://wails.io) v2 正式发布。距离首个 v2 Alpha 版本发布大约已有18个月，距离首个 Beta 版本发布也已有大约一年。我由衷感谢参与推动项目发展的每一个人。

之所以花了这么长时间，部分原因是我们希望在正式将其称为 v2 之前，先达到某种意义上的完整。事实上，标记发行版本永远没有完美的时机——总会有尚未解决的问题，或是想再塞进去的“最后一个”功能。不过，标记一个并不完美的主版本，确实能为项目用户提供一定的稳定性，也能让开发者重新启程。

这个版本远远超出了我曾经的预期。我希望它能给你带来如同我们开发它时所感受到的那般快乐。

## Wails<em>是什么</em>？

如果你不熟悉 Wails，它是一个让 Go 程序员能够使用熟悉的 Web 技术，为自己的 Go 程序构建丰富前端的项目。它是 Electron 的轻量级 Go 替代方案。你可以在[官方网站](https://wails.io/docs/introduction)上找到更多信息。

## 有哪些新变化？

v2 版本是这个项目的一次巨大飞跃，解决了 v1 的许多痛点。如果你尚未阅读有关 [macOS](/blog/wails-v2-beta-for-mac/)、[Windows](/blog/wails-v2-beta-for-windows/) 或 [Linux](/blog/wails-v2-beta-for-linux/) Beta 版本的博文，我建议你读一读，其中更详细地介绍了所有重大变化。简而言之：

- 面向 Windows 的 WebView2 组件，支持现代 Web 标准和调试功能。
- Windows 上的[深色/浅色主题](https://wails.io/docs/reference/options#theme)和[自定义主题](https://wails.io/docs/reference/options#customtheme)。
- Windows 现在不再需要 CGO。
- 开箱即用地支持 Svelte、Vue、React、Preact、Lit 和 Vanilla 项目模板。
- 集成[Vite](https://vitejs.dev/)，为应用提供支持热重载的开发环境。
- 原生应用[菜单](https://wails.io/docs/guides/application-development#application-menu)和[对话框](https://wails.io/docs/reference/runtime/dialog)。
- 为[Windows](https://wails.io/docs/reference/options#windowistranslucent)和[macOS](https://wails.io/docs/reference/options#windowistranslucent-1)提供原生窗口半透明效果。支持 Mica 和 Acrylic 背景材质。
- 可轻松生成用于 Windows 部署的[NSIS 安装程序](https://wails.io/docs/guides/windows-installer)。
- 功能丰富的[运行时库](https://wails.io/docs/reference/runtime/intro)，为窗口操作、事件处理、对话框、菜单和日志记录提供实用方法。
- 支持使用[garble](https://github.com/burrowers/garble)对应用进行[混淆](https://wails.io/docs/guides/obfuscated)。
- 支持使用[UPX](https://upx.github.io/)压缩应用。
- 根据 Go 结构体自动生成 TypeScript。更多信息请参阅[此处](https://wails.io/docs/howdoesitwork#calling-bound-go-methods)。
- 无论在哪个平台，都无需随应用一起分发额外的库或 DLL。
- 无需打包前端资源。像开发任何其他 Web 应用一样开发你的应用即可。

## 致谢

实现 v2 是一项艰巨的工作。从最初的 alpha 版本到今天正式发布，89 位贡献者共提交了约 2200 次；此外，还有许许多多的人提供了翻译、测试、反馈和帮助，并参与了讨论论坛和问题跟踪器中的交流。我对你们每一位都感激不尽。还要特别感谢所有为项目提供指导、建议和反馈的赞助者。你们所做的一切都弥足珍贵。

我想特别提到几个人：

首先，向 [@stffabi](https://github.com/stffabi) 致以<strong>万分</strong>感谢。他做出了许多让我们所有人受益的贡献，还为许多问题提供了大量支持。他实现了一些关键功能，例如外部开发服务器支持；这让我们能够利用[Vite](https://vitejs.dev/)的强大能力，从而彻底改变了我们的开发模式。可以毫不夸张地说，如果没有他的[杰出贡献](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04)，Wails v2 绝不会是现在这样令人振奋的版本。非常感谢你，@stffabi！

我还要特别感谢 [@misitebao](https://github.com/misitebao)。他不辞辛劳地维护网站、提供中文翻译、管理 Crowdin，并帮助新译者快速上手。这是一项极其重要的工作，我由衷感谢为此投入的所有时间和精力！你太棒了！

最后同样重要的是，我要衷心感谢 Mat Ryer 在 v2 开发期间提供的建议和支持。我们曾使用 v2 的早期 Alpha 版本共同开发 xBar，这不仅帮助确定了 v2 的发展方向，也让我认识到早期版本中的一些设计缺陷。我很高兴地宣布，从今天起，我们将开始把 xBar 移植到 Wails v2，它也将成为本项目的旗舰应用。谢谢你，Mat！

## 经验教训

在迈向 v2 的过程中，我们吸取了许多经验教训，它们将影响今后的开发方向。

## 规模更小、节奏更快、目标更明确的版本

在开发 v2 的过程中，许多功能和错误修复都是临时开展的。这导致发布周期更长，调试也更加困难。今后，我们将更频繁地发布版本，每个版本包含的功能数量会减少。每次发布都将包括文档更新和全面测试。希望这些规模更小、节奏更快、目标更明确的版本能够减少回归问题，并带来质量更高的文档。

## 鼓励参与

刚开始这个项目时，我希望立即帮助每一个遇到问题的人。我把每个问题都当成“自己的事”，想尽快解决。这种做法不可持续，最终反而不利于项目的长久发展。今后，我会为其他人参与回答问题和分类处理议题留出更多空间。如果能有一些工具协助完成这些工作会很有帮助，因此如果你有任何建议，请在[此处](https://github.com/wailsapp/wails/discussions/1855)参与讨论。

## 学会说“不”

参与开源项目的人越多，要求增加新功能的请求就越多，而这些功能对大多数人而言可能有用，也可能没用。开发和调试这些功能需要前期投入时间，此后还会持续产生维护成本。在这方面，我自己最为严重：我常常想要“大包大揽”，而不是先提供最小可行的功能。今后，对于向核心添加功能的请求，我们需要更常说“不”，并把精力集中在如何让开发者自行实现这些功能上。我们正在认真考虑通过插件来满足这种场景。这样，任何人都能按自己的需求扩展项目，同时也能方便地为项目作出贡献。

## 展望未来

我们已经在考虑下一个主要开发周期中要为 Wails 添加的众多核心功能。[路线图](https://github.com/wailsapp/wails/discussions/1484)中充满了有趣的想法，我迫不及待地想开始着手实现。其中呼声很高的一项是多窗口支持。这是个棘手的问题，要想妥善实现，我们可能需要考虑提供一套替代 API，因为当前 API 在设计之初并未考虑这一需求。根据一些初步构想和反馈，我想你会喜欢我们正在探索的方向。

我个人非常期待让 Wails 应用在移动设备上运行。我们已经有一个演示项目，证明 Wails 应用可以在 Android 上运行，因此我非常想探索这方面还能走多远！

最后，我想谈谈功能一致性。长期以来，我们一直坚持一项核心原则：除非某项功能能够获得完整的跨平台支持，否则就不会将其添加到项目中。尽管到目前为止，这一点（基本上）能够做到，但它确实阻碍了项目发布新功能。今后，我们将采用稍有不同的方式：任何无法立即面向所有平台发布的新功能，都将通过实验性配置或 API 发布。这样，特定平台上的早期采用者便可试用该功能并提供反馈，这些反馈将用于完善该功能的最终设计。当然，这意味着在该功能获得所有能够支持它的平台的完整支持之前，API 的稳定性无法得到保证，但至少不会再阻碍开发。

## 最后的话

我为我们在 V2 版本中取得的成果感到无比自豪。看到大家迄今已经使用各个测试版构建出诸如[Varly](https://varly.app/)、[Surge](https://getsurge.io/)和[October](https://october.utf9k.net/)这样的优质应用，实在令人赞叹。我推荐你去体验一下。

此版本凝聚了众多贡献者的辛勤付出。虽然它可以免费下载和使用，但并非毫无成本。毫无疑问，这个项目付出了相当大的代价。不仅包括我和每一位贡献者投入的时间，也包括我们每个人因此无法陪伴亲友所付出的代价。因此，我对投入这个项目的每一分每一秒都深怀感激。贡献者越多，这份工作就越能由更多人分担，我们共同取得的成就也会越大。我想鼓励大家选择一件力所能及的事来贡献，无论是确认他人报告的 bug、提出修复建议、修改文档，还是帮助有需要的人。所有这些小事都能产生巨大的影响！如果你也能成为迈向 v3 这段历程的一员，那就太棒了。

尽情享用吧！

&dash; Lea

附言：如果你或你的公司认为 Wails 很有用，请考虑[赞助本项目](https://github.com/sponsors/leaanthony)。谢谢！
