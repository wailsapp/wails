---
title: "Snippet Expander"
description: "使用 Wails 建置的桌面應用程式"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Snippet Expander 螢幕截圖](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Snippet Expander「選取片段」視窗的螢幕截圖

![Snippet Expander 螢幕截圖](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Snippet Expander「新增片段」畫面的螢幕截圖

![Snippet Expander 螢幕截圖](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Snippet Expander「搜尋與貼上」視窗的螢幕截圖

[Snippet Expander](https://snippetexpander.org)是 Linux 平台上的「您的可展開文字片段小幫手」。

Snippet Expander 包含一個使用 Wails 建置的 GUI 應用程式，用於管理片段和設定，並提供「搜尋與貼上」視窗模式，讓您快速選取並貼上片段。

以 Wails 建置的 GUI、go-lang CLI 和 vala-lang 自動展開常駐程式都會透過 D-Bus 與 go-lang 常駐程式通訊。該常駐程式負責大部分工作，包括管理片段資料庫和共用設定，以及提供展開和貼上片段等服務。

請查看[原始碼](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38)，瞭解 Wails 應用程式如何將訊息從 UI 傳送至後端，再由後端傳送至常駐程式；以及如何訂閱 D-Bus 事件，以監控其他應用程式執行個體或 CLI 對片段所做的變更，並透過 Wails 事件立即將變更顯示在 UI 中。
