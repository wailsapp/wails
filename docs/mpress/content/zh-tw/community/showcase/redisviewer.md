---
title: "Redis 檢視器"
description: "使用 Wails 建置的桌面 Redis 圖形化介面"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![RedisViewer 螢幕截圖](/assets/showcase-images/redisviewer-overview1.webp)

![RedisViewer 螢幕截圖](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) 是一款使用 Wails 建置的現代化 Redis 桌面圖形化介面。您可以用它檢查複雜值、執行命令及分析 Redis 效能，同時兼顧互動品質。

它以 Wails 的 WebView + Go 架構為核心，將大型鍵空間和大量承載資料保留在 Go 後端，而非全部傳送至前端，藉此運用 Go 的記憶體模型，並避免以 JS 為主的用戶端常見的堆積記憶體壓力和記憶體洩漏風險。使用者介面經過調校，只會算繪畫面上的內容，因此即使瀏覽龐大資料集，也能保持靈敏流暢，體驗接近原生桌面應用程式。

[造訪專案網站](https://redisviewer.com/)
