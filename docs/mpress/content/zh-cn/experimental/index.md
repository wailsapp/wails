---
title: "实验性功能"
description: "集中介绍 Wails v3 中正在进行的实验：它们是什么、为何存在，以及如何提供反馈。"
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="这里是实验区"}
根据定义，本节中的所有内容都是实验性的。这些功能需要主动启用，默认关闭，并且可能在不同版本之间改变形式、重命名，甚至被完全移除。除非已准备好应对频繁变动，否则不要让任何关键功能依赖它们。

@end

## “实验性”的含义

Wails 以公开方式推出实验。实验是我们认为很有前景、值得尽早交到你手中的想法，但我们尚未承诺一定会长期保留。我们之所以发布它，<em>是因为</em>希望在决定它是否成为 Wails 永久且受支持的一部分之前，先从实际使用中吸取经验。

这意味着本节中的所有内容都有以下几个特点：

- <strong>需要主动启用。</strong>实验绝不会改变`wails3`的默认行为。你需要有意启用它们（通常通过环境变量或构建标志）；关闭它们后，你现有的工作流程不会发生任何变化。
- <strong>它可能不会保留下来。</strong>有些实验会发展为稳定功能，另一些则会被彻底改造，甚至放弃。相比只发布我们已经完全确定的内容，我们更愿意公开尝试并快速吸取经验。
- <strong>API 尚未冻结。</strong>在实验逐步定型的过程中，名称、标志、默认值和行为都可能在不同版本之间发生变化。发布说明会明确指出这些变更，但不要期待稳定功能所提供的稳定性保证。

## 我们希望听到你的反馈

这是最重要的部分。实验的去留取决于实际使用者给我们的反馈。如果你尝试了其中一项实验，我们真心希望了解：

- 它是否适用于你的项目？哪里没有达到预期？
- 它是否更快、更清晰、更好用，还是不值得切换？
- 要满足哪些条件，你才会默认使用它？

最有用的反馈应当具体说明：你运行了什么、预期结果是什么，以及实际发生了什么。每项实验在 GitHub Discussions 的<strong>实验</strong>类别中都有自己的讨论串：

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="实验讨论" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="找到你所用实验的讨论串，告诉我们它的实际表现、哪里出了问题，或还缺少什么。"}
@end

## 当前实验

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="一个可替代现有构建运行器、理解 Wails 并可用于现有 Taskfile 的工具。增量构建更快，提供结构化输出，并默认并行执行。"}
@linkcard{title="LLM 控制（MCP）" href="/guides/mcp-service/" description="内置的模型上下文协议服务器，让 LLM 智能体能够检查、测试和操控正在运行的 Wails 应用；无需编写用户代码，通过构建标签启用。"}
@end
