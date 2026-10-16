---
title: "實驗性功能"
description: "彙整 Wails v3 進行中實驗的專區，說明這些實驗是什麼、為何存在，以及如何提供意見回饋。"
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="此處有實驗性功能"}
本節中的所有內容，顧名思義，都是實驗。這些功能須由使用者主動啟用，預設為停用，而且在不同版本之間可能會改變形式、重新命名，甚至完全移除。除非您已準備好因應頻繁變動，否則請勿以這些功能建構任何關鍵部分。

@end

## 「實驗性」代表什麼

Wails 以公開方式推出實驗性功能。實驗代表我們認為某個構想很有潛力，值得提早交到您手中，但尚未承諾會完整支援。我們之所以將它公開，<em>是因為</em>希望先從實際使用情況中汲取經驗，再決定它是否會成為 Wails 永久且受支援的一部分。

這表示本節中的所有內容都有以下幾項特性：

- <strong>必須主動啟用。</strong>實驗絕不會改變`wails3`的預設行為。您需要刻意將其開啟（通常透過環境變數或建置旗標）；在關閉時，您現有的工作流程不會有任何改變。
- <strong>它可能無法保留下來。</strong>有些實驗會發展為穩定功能；有些則會徹底改造，甚至遭到捨棄。相較於只推出我們已有十足把握的功能，我們更願意公開嘗試並迅速汲取經驗。
- <strong>API 尚未定案。</strong>在實驗逐步定型的過程中，名稱、旗標、預設值和行為都可能隨版本變動。版本資訊會特別說明這些變更，但請勿期待它具備穩定功能所提供的穩定性保證。

## 我們希望收到您的意見回饋

這是最重要的部分。實驗能否繼續發展，取決於實際使用者提供的意見回饋。如果您嘗試其中一項實驗，我們衷心希望了解：

- 它是否適用於您的專案？在哪些地方未能達到預期？
- 它是否更快、更清楚、更好用，還是根本不值得切換？
- 必須具備哪些條件，您才會願意預設使用它？

最有幫助的意見回饋應具體說明：您執行了什麼、預期會發生什麼，以及實際發生了什麼。每項實驗在 GitHub Discussions 的<strong>實驗</strong>類別中都有各自的討論串：

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="實驗討論區" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="找到您所用實驗的討論串，告訴我們實際使用結果、發生了哪些問題，或還缺少哪些功能。"}
@end

## 目前的實驗

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="可搭配現有 Taskfile 使用、能感知 Wails 的替代建置執行器。提供更快速的增量建置、結構化輸出，並預設平行執行。"}
@linkcard{title="LLM 控制（MCP）" href="/guides/mcp-service/" description="內建的 Model Context Protocol 伺服器，可讓 LLM 代理檢查、測試及操控執行中的 Wails 應用程式；不需要使用者程式碼，透過建置標籤即可啟用。"}
@end
