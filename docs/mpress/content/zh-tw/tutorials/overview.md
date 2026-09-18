---
title: "教學"
description: "透過建置應用程式來學習 Wails"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

這些逐步教學會透過建置完整的應用程式，帶你學習 Wails 的各項概念。每篇教學都包含可運作的程式碼、說明與實用模式。

@note{type="tip" title="剛接觸 Go？"}
開始教學前，請先完成[Go 導覽](https://go.dev/tour/)。

@end

## QR Code 服務

![QR Code 範例](/assets/qr1.png)

透過建置 QR Code 產生器，學習 Wails 服務的基礎知識。本教學將介紹如何把應用程式邏輯組織成可重複使用的服務等核心概念。

**你將學到：**

- 如何建立及組織 Wails 服務
- 管理外部 Go 相依套件
- 將 Go 方法繫結至前端
- 在 Go 與 JavaScript 之間傳遞資料
- 組織程式碼以提升可維護性

<strong>最適合：</strong>想瞭解服務架構的 Wails 初次使用者

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 待辦事項清單

![待辦事項清單應用程式](/assets/todo-app.png)

建置具備美觀現代介面的完整待辦事項清單應用程式。這篇實作教學使用原生 JavaScript，透過實際的真實應用程式教你掌握 Wails 的核心模式。

**你將學到：**

- 採用執行緒安全狀態管理的服務式架構
- CRUD 操作（建立、讀取、更新、刪除）
- Go 與 JavaScript 之間的型別安全繫結
- 在沒有框架複雜度的情況下建置現代化 UI
- 妥善的錯誤處理與驗證模式

<strong>完成時間：</strong>約 20 分鐘

<strong>最適合：</strong>你的第一個完整 Wails 應用程式——非常適合在加入框架的複雜度前先瞭解基礎知識

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 筆記

![筆記應用程式](/assets/notes-app.png)

建置具有原生檔案對話方塊與自動儲存功能、風格類似 Apple 備忘錄的應用程式。本教學將示範檔案操作、原生對話方塊及專業 UI 模式等桌面環境特有功能。

**你將學到：**

- 原生檔案對話方塊（儲存、開啟、資訊）
- 以 JSON 為基礎的資料持久化
- 採用防彈跳機制的自動儲存模式
- 專業的雙欄桌面版面配置
- 在 Go 中處理檔案系統操作

<strong>完成時間：</strong>約 30 分鐘

<strong>最適合：</strong>學習檔案操作與作業系統原生對話方塊等桌面環境特有功能

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 可自行更新的 Wails 應用程式

從全新的`wails3 init`開始，為 Wails 應用程式加入應用程式內自我更新功能，並完成簽署版本驗證及輔助程式模式的二進位檔交換。使用 GitHub Releases 作為更新來源。

**你將學到：**

- 如何將`app.Updater`整合至 Wails 應用程式
- 設定 GitHub Releases 提供者
- 使用`SHA256SUMS`發布版本以進行摘要驗證
- 加入 Ed25519 簽署以防止竄改
- 透過 CSS、自訂 HTML 或 BYO 自訂預設視窗
- 使用`CheckInterval`定期在背景檢查

<strong>完成時間：</strong>約 25 分鐘

<strong>最適合：</strong>發布可更新的桌面應用程式——涵蓋完整的發布管線，而不只是 API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
