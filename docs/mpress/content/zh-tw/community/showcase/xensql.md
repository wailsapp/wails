---
title: "XenSQL"
description: "支援多種資料庫與進階查詢工具的本機優先 SQL 工作台"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![具備行內編輯功能的 XenSQL 資料表編輯器](/assets/showcase-images/xensql-1.png) ![具備資料表與欄位建議功能的 XenSQL 查詢編輯器](/assets/showcase-images/xensql-2.png) ![XenSQL 多項交易查詢](/assets/showcase-images/xensql-3.png)

**[XenSQL](https://github.com/Bare7a/XenSQL)** 是一款使用<strong>Go、Wails 與 React</strong>打造的<strong>快速、本機優先 SQL 桌面工作台</strong>。它將 PostgreSQL、MySQL/MariaDB 與 SQLite 整合在同一個簡潔、具有原生應用程式質感的介面中，而且完全不使用雲端、不收集遙測資料，也不需要帳號。

## 主要特色

- **強大的 SQL 編輯器** — 以 Monaco 為基礎，提供可感知結構描述的智慧自動完成、多陳述式執行、串流結果，以及各陳述式專屬的結果分頁
- **進階資料檢視器** — 互動式 JSON 檢查器、可感知語法的儲存格編輯器（JSON、XML、HTML、文字）、行內編輯，以及完整記錄檢視
- **流暢的資料編輯體驗** — 瀏覽資料表、以行內方式暫存變更、執行大量操作，以及安全地進行`INSERT`/`UPDATE`/`DELETE`，並支援`RETURNING`
- **生產力功能** — 結構描述瀏覽器、已儲存的查詢、查詢歷程記錄、快速搜尋（`Ctrl+P`），以及鍵盤優先的工作流程
- **匯出選項** — CSV、JSON、Markdown、SQL INSERT 陳述式

完全離線且可攜。所有內容都儲存在單一的本機`XenSQL-data/`資料夾中，並可隨應用程式一起移動。

**支援的資料庫**：PostgreSQL、MySQL、MariaDB 與 SQLite（提供唯讀模式與安全傳輸選項）。

專為追求速度、清晰度與控制權，同時不想承受傳統 SQL 工具臃腫負擔的開發人員而設計。
