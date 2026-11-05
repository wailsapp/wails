---
title: "維護者分流處理"
description: "以一致方式分流處理錯誤、文件問題回報及 WEP"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## 目的

Issue 用於記錄可重現的錯誤及文件問題。WEP（Wails 增強提案）Pull Request 用於記錄所提議的功能或公開行為變更。

## Issue 分流處理

- 確認錯誤報告包含已發布的版本、平台、可重現問題的步驟、預期行為、實際行為，以及 `wails3 doctor` 輸出。
- 為已確認的報告加上 `Bug` 標籤，並套用相關的版本及平台標籤。必要時，要求提供最小重現範例。
- 若文件問題回報指出具體的文件缺陷，請保持其開啟；若回報者能進行修改，鼓勵其提出 PR。
- 將功能要求引導至 WEP 指南，然後將其關閉。自動重新引導工作流程會處理新加上 enhancement 標籤的 Issue；對較舊的 Issue 使用相同措辭。
- 將問題及支援要求移至 GitHub Discussions 或 Discord。

## WEP 分流處理

1. 確認該 PR 是標題為 `[WEP] <title>` 的草稿，且僅包含 WEP 及其支援資料。
2. 確認該 PR 使用 WEP 範本、指明實作者，並涵蓋相容性、平台、測試、維護及安全性／隱私權。
3. 將技術討論保留在 WEP PR 中。Discussions 可提供有用的脈絡，但不是決策紀錄。
4. 在 PR 留言中記錄維護者的決定：接受、拒絕或撤回，並附上簡短理由。
5. 對於已接受的 WEP，請指派其編號、更新 WEP 索引、合併 WEP PR，並要求實作 PR 連結回該 WEP。

## 既有的增強功能 Issue

請勿在未告知的情況下刪除歷史增強功能 Issue。對於每項仍具相關性的 要求，請留下重新引導留言並關閉該 Issue；有興趣的貢獻者可以 提出 WEP。關閉重複的 Issue 時，請附上既有 WEP 或決策的連結。
