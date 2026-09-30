---
title: "Redis 查看器"
description: "使用 Wails 构建的桌面 Redis 图形界面"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![RedisViewer 截图](/assets/showcase-images/redisviewer-overview1.webp)

![RedisViewer 截图](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) 是一款使用 Wails 构建的现代 Redis 桌面图形界面。你可以用它检查复杂值、运行命令和分析 Redis 性能，同时保持高质量的交互体验。

它围绕 Wails 的 WebView + Go 架构设计，将大型键空间和大体量载荷保留在 Go 后端，而不是把所有内容都传到前端，从而利用 Go 的内存模型，并规避重度依赖 JS 的客户端中常见的堆内存压力和内存泄漏风险。该界面经过优化，仅渲染屏幕上显示的内容，因此即使浏览海量数据集也能保持响应迅速，体验接近原生桌面应用。

[访问项目网站](https://redisviewer.com/)
