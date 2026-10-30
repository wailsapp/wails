---
title: "教學"
description: "透過建置應用程式學習 Wails"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

想透過實作來學習 Wails 時，請選擇引導式課程。每堂課都會帶你完成一個專案，達成具體成果。

如果這是你第一次造訪，請先完成[安裝](/quick-start/installation/)和[你的第一個應用程式](/quick-start/first-app/)。

## 選擇教學

完成第一個應用程式後，請依照以下順序進行桌面應用程式課程：

| 課程 | 你將練習的內容 |
| --- | --- |
| [QR Code 服務](/tutorials/01-creating-a-service/) | 透過服務公開 Go 邏輯，並從前端呼叫 |
| [待辦事項清單](/tutorials/02-todo-vanilla/) | 建立、更新及刪除應用程式資料 |
| [筆記應用程式](/tutorials/03-notes-vanilla/) | 使用原生對話方塊並將資料儲存至檔案 |
| [可自動更新的應用程式](/tutorials/04-self-update-a-wails-app/) | 從 GitHub Releases 檢查並套用更新 |

若要選修行動應用程式課程，請參閱[你的第一個行動應用程式](/guides/mobile/first-mobile-app/)。行動平台支援仍處於實驗階段。

如果你正在處理現有應用程式，並需要完成特定工作，請改用[操作指南](/quick-start/next-steps/#features)。若要瞭解設計背景，請參閱[核心概念](/concepts/architecture/)。

@note{type="tip" title="剛接觸 Go？"}
開始教學前，請先完成[Go 導覽](https://go.dev/tour/)。

@end

## QR Code 服務

![QR Code 範例](/assets/qr1.png)

建置 QR Code 產生器，學習 Wails 服務的基礎知識。本教學將帶你認識如何將應用程式邏輯整理成可重複使用服務的核心概念。

**你將學到：**

- 如何建立 Wails 服務並安排其結構
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

建置具備美觀現代化介面的完整待辦事項清單應用程式。這套實作教學使用原生 JavaScript，透過實用的真實應用程式帶你學習 Wails 的核心模式。

**你將學到：**

- 採用執行緒安全狀態管理的服務式架構
- CRUD 操作（建立、讀取、更新、刪除）
- Go 與 JavaScript 之間具型別安全性的繫結
- 不受框架複雜度影響地建置現代化 UI
- 妥善的錯誤處理與驗證模式

<strong>完成時間：</strong>約20分鐘

<strong>最適合：</strong>你的第一個完整 Wails 應用程式——非常適合在加入框架的複雜度之前，先瞭解基礎知識

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 筆記

![筆記應用程式](/assets/notes-app.png)

建置一個 Apple Notes 風格的應用程式，加入原生檔案對話方塊與自動儲存功能。本教學示範檔案操作、原生對話方塊及專業 UI 模式等桌面平台專屬功能。

**你將學到：**

- 原生檔案對話方塊（儲存、開啟、資訊）
- 以 JSON 為基礎的資料持久化
- 具防彈跳機制的自動儲存模式
- 專業的雙欄桌面版面配置
- 在 Go 中處理檔案系統操作

<strong>完成時間：</strong>約30分鐘

<strong>最適合：</strong>學習檔案操作和作業系統原生對話方塊等桌面平台專屬功能

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 可自動更新的 Wails 應用程式

從全新的`wails3 init`開始，為 Wails 應用程式加入應用程式內自動更新功能，完成簽署版本驗證及輔助模式的二進位檔案置換。使用 GitHub Releases 作為更新來源。

**你將學到：**

- `app.Updater`如何整合至 Wails 應用程式
- 設定 GitHub Releases 提供者
- 使用`SHA256SUMS`發佈版本，以進行摘要驗證
- 加入 Ed25519 簽署以防止竄改
- 透過 CSS、自訂 HTML 或自備介面來自訂預設視窗
- 使用`CheckInterval`定期在背景檢查

<strong>完成時間：</strong>約25分鐘

<strong>最適合：</strong>發佈可更新的桌面應用程式——涵蓋完整的發佈流程，而不只是 API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>開始使用</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
