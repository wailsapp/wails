---
title: "维护者分类处理"
description: "以一致的方式分类处理错误、文档问题报告和 WEP"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## 目的

Issue 用于记录可复现的错误和文档问题。WEP（Wails 增强提案）拉取请求用于记录对某项功能或公开行为变更的提议。

## Issue 分类处理

- 确认错误报告包含已发布版本、平台、复现步骤、预期行为、实际行为以及`wails3 doctor`输出。
- 为已确认的报告添加`Bug`标签，并添加相关的版本和平台标签。必要时要求提供最小复现示例。
- 如果文档问题报告指出了具体的文档缺陷，请将其保持为开放状态；如果报告者能够完成更改，则鼓励其提交 PR。
- 将功能请求引导至 WEP 指南，然后关闭。自动重定向工作流会处理新添加 enhancement 标签的 Issue；对于较早的 Issue，请使用相同的措辞。
- 将问题和支持请求转至 GitHub Discussions 或 Discord。

## WEP 分类处理

1. 检查该 PR 是否为标题为`[WEP] <title>`的草稿，并且仅包含 WEP 及其支持材料。
2. 检查该 PR 是否使用 WEP 模板、明确实施者，并涵盖兼容性、平台、测试、维护以及安全与隐私。
3. 将技术讨论保留在 WEP PR 中。Discussions 可提供有用的上下文，但并非决策记录。
4. 在 PR 评论中记录维护者的决定：接受、拒绝或撤回，并附上简短理由。
5. 对于已接受的 WEP，为其分配编号、更新 WEP 索引、合并 WEP PR，并要求实施 PR 链接回该 WEP。

## 现有功能增强 Issue

不要静默删除历史功能增强 Issue。对于每个仍然相关的 请求，请留下重定向评论并将其关闭；感兴趣的贡献者可以 发起 WEP。关闭重复 Issue 时，请附上指向现有 WEP 或决定的链接。
