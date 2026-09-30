---
title: "反馈"
description: "如何为 Wails v3 提供反馈和报告问题"
slug: "feedback"
sourcePath: "feedback.md"
---

我们欢迎（并鼓励）你提供反馈！创建新 issue 或讨论前，请先搜索已有内容。你可以通过以下不同方式参与贡献：

@tabs
[Bug]
如果发现 bug，请使用 bug 报告模板在 GitHub 上[创建 issue](https://github.com/wailsapp/wails/issues/new/choose)。

- 请通过简单且可复现的示例清楚描述 bug。如果文档未明确说明<em>预期</em>行为，也请在报告中指出。
- 请在报告中包含`wails3 doctor`的输出。
- 如果该 bug 表现为与当前文档不一致的行为，还请执行以下操作：
  - 更新`v3/examples`目录中的现有示例，或创建一个能清楚展示该问题的新示例。
  - 创建一个引用该 issue 的[PR](https://github.com/wailsapp/wails/pulls)。


@note{type="caution"}
*请注意*，意外行为不一定是 bug——它可能只是没有按你的预期运行。此类情况请使用`Suggestions`。

@end

也欢迎你在 Discord 的[#v3](https://discord.gg/bdj28QNHmT)频道中讨论 bug。

[修复]
如果你修复了 bug 或改进了文档，请：

- 按照[贡献指南](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)在[Wails 仓库](https://github.com/wailsapp/wails)中创建 pull request。
- 在 PR 描述中引用所有相关 issue。

[增强功能]
新功能和公共行为变更应通过<strong>WEP（Wails Enhancement Proposal）</strong>草案 pull request 提出，而不是创建功能请求 issue。

- 阅读[WEP 流程](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)。
- 复制模板，然后创建一个标题为`[WEP] <title>`的 PR 草案，其中只能包含 WEP 及其支持材料。
- 你可以先在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)或 Discord 的[#v3](https://discord.gg/bdj28QNHmT)频道中非正式地讨论想法，但维护者必须基于 WEP PR 作出决定。

[点赞支持]
- 使用 GitHub 上的 :thumbsup: 表情回应来支持 bug、WEP 和讨论。
- 请<em>不要</em>只添加“+1”或“我也是”之类的评论。
- 如果你有实质性内容要补充，请添加评论，例如“这个 bug 也会影响 ARM 构建”或“另一种方法是……”。

@end

可以在[此处](https://github.com/orgs/wailsapp/projects/6)查看已知问题和正在进行的工作。
