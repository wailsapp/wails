---
title: "REST API"
description: "从 Wails v3 应用程序公开 RESTful API"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
此页面为占位页面。完整的 REST API 内容即将推出。

@end

你可以使用标准 Go `net/http`处理程序或[Gin](/guides/gin-routing/)等框架，从 Wails v3 应用程序公开 RESTful API。由于 Wails 通过资源处理程序提供前端内容，因此可以挂载任何 HTTP 处理程序，并从 UI 或其他客户端访问它。
