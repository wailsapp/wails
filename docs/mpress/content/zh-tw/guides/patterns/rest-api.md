---
title: "REST API"
description: "從 Wails v3 應用程式公開 RESTful API"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
此頁面為預留頁面。完整的 REST API 內容即將推出。

@end

您可以使用標準 Go `net/http`處理常式或 [Gin](/guides/gin-routing/) 等框架，從 Wails v3 應用程式公開 RESTful API。由於 Wails 透過資產處理常式提供前端，因此可以掛載任何 HTTP 處理常式，並從使用者介面或其他用戶端存取。
