---
title: "反馈"
description: "如何为 Wails v3提供反馈和报告问题"
slug: "feedback"
sourcePath: "feedback.md"
---

我们欢迎（并鼓励）您提供反馈！创建新议题或讨论前，请先搜索是否已有相关内容。您可以通过以下方式参与贡献：

@tabs
[缺陷]
如果发现缺陷，请使用缺陷报告模板在 GitHub 上[创建议题](https://github.com/wailsapp/wails/issues/new/choose)。

- 请通过简单且可复现的示例清楚描述该缺陷。如果文档没有明确说明*应该*出现什么结果，也请在报告中说明这一点。
- 请在报告中附上`wails3 doctor`的输出。
- 如果该缺陷表现为行为与当前文档不符，还请执行以下操作：
  - 更新`v3/examples`目录中的现有示例，或创建一个能清楚展示该问题的新示例。
  - 创建一个引用该议题的[PR](https://github.com/wailsapp/wails/pulls)。


@note{type="caution"}
*请注意*，非预期行为不一定是缺陷——它可能只是没有按您的预期运行。遇到这种情况，请使用`Suggestions`。

@end

也欢迎您在 Discord 的[#v3](https://discord.gg/bdj28QNHmT)频道中讨论缺陷。

[修复]
如果您修复了某个缺陷或改进了文档，请：

- 按照[贡献指南](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)，在[Wails 仓库](https://github.com/wailsapp/wails)中创建拉取请求。
- 请在 PR 描述中引用所有相关议题。

[功能增强]
新功能和公共行为的变更应通过<strong>WEP（Wails Enhancement Proposal）</strong>草案拉取请求提出，而不是通过功能请求议题提出。

- 请阅读[WEP 流程](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md)。
- 复制模板，然后创建标题为`[WEP] <title>`的草案 PR，其中仅包含 WEP 及其支持材料。
- 您可以先在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)或 Discord 的[#v3](https://discord.gg/bdj28QNHmT)频道中非正式地讨论想法，但维护者必须依据 WEP PR 作出决定。

[投票支持]
- 请在 GitHub 上使用 :thumbsup: 表情回应来支持缺陷、WEP 和讨论。
- 请*不要*只添加“+1”或“我也是”之类的评论。
- 如果您有实质性内容可以补充，请发表评论，例如“此缺陷也会影响 ARM 构建”或“另一种方法是……”。

@end

可在[此处](https://github.com/orgs/wailsapp/projects/6)查看已知问题和正在进行的工作。
