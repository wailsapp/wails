---
title: "Snippet Expander"
description: "使用 Wails 构建的桌面应用程序"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Snippet Expander 截图](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Snippet Expander“选择文本片段”窗口的截图

![Snippet Expander 截图](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Snippet Expander“添加文本片段”界面的截图

![Snippet Expander 截图](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Snippet Expander“搜索并粘贴”窗口的截图

[Snippet Expander](https://snippetexpander.org)是一款适用于 Linux 的“帮你展开文本 片段的小助手”。

Snippet Expander 包含一个使用 Wails 构建的 GUI 应用程序，用于管理 文本片段和设置；它还提供“搜索并粘贴”窗口模式，可快速选择 并粘贴文本片段。

基于 Wails 的 GUI、go-lang CLI 和 vala-lang 自动扩展守护进程均 通过 D-Bus 与 go-lang 守护进程通信。该守护进程承担大部分 工作，包括管理文本片段数据库和通用设置，以及提供 扩展和粘贴文本片段等服务。

查看 [源代码](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38)， 了解 Wails 应用如何将来自 UI 的消息发送到后端，再由后端 将其发送给守护进程；还可了解该应用如何订阅 D-Bus 事件，以监控通过 应用的其他实例或 CLI 对文本片段所做的更改，并通过 Wails 事件立即在 UI 中显示这些更改。
