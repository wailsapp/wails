---
title: "後續步驟"
description: "建立第一個應用程式後的後續學習方向"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

你已建立第一個 Wails 應用程式。請依照下方的學習路徑繼續，若已知道要建立什麼，也可以直接前往相關指南。

## 深入學習

每次透過一個完整的應用程式繼續學習：

1. <strong>[建立 QR Code 服務](/tutorials/01-creating-a-service/)。</strong>練習透過服務將 Go 邏輯公開給前端。
2. <strong>[建立待辦事項清單](/tutorials/02-todo-vanilla/)。</strong>透過服務呼叫新增、修改及刪除應用程式資料。
3. <strong>[建立筆記應用程式](/tutorials/03-notes-vanilla/)。</strong>使用檔案、原生對話方塊及已儲存的資料。

準備好探索版本發布後，請按照[應用程式自動更新教學](/tutorials/04-self-update-a-wails-app/)操作。如需瞭解各課程的學習成果，請參閱[所有教學](/tutorials/overview/)。

## 功能

<strong>已經在建立自己的應用程式了嗎？</strong>請針對目前的工作使用操作指南：

| 我想要…… | 從這裡開始 |
| --- | --- |
| 從前端呼叫 Go | [從 JavaScript 呼叫 Go](/features/bindings/methods/) |
| 開啟或管理視窗 | [視窗基礎](/features/windows/basics/) |
| 新增原生選單或系統匣圖示 | [應用程式選單](/features/menus/application/)或[系統匣選單](/features/menus/systray/) |
| 要求使用者開啟或儲存檔案 | [檔案對話方塊](/features/dialogs/file/) |
| 在元件之間傳送事件 | [事件系統](/features/events/system/) |
| 複製文字或接收拖放的檔案 | [剪貼簿](/features/clipboard/basics/)或[檔案拖放](/features/drag-and-drop/files/) |
| 新增鍵盤快速鍵 | [鍵盤快速鍵](/features/keyboard/shortcuts/) |
| 建置及封裝應用程式 | 先參閱[建置應用程式](/guides/build/building/)，再針對[Windows](/guides/build/windows/)、[macOS](/guides/build/macos/)或[Linux](/guides/build/linux/)進行封裝 |

這些頁面位於側邊欄的**操作指南**下方。你不必依序閱讀整個章節，可直接開啟相關主題。

## 瞭解 Wails 的運作方式

若想瞭解程式碼背後的設計，請參閱**核心概念**：

- [Wails 的運作方式](/concepts/architecture/)說明 Go 後端、前端與原生 WebView 之間的連接。
- [Go 與前端橋接](/concepts/bridge/)說明呼叫如何在 Go 與 JavaScript 之間傳遞。
- [應用程式生命週期](/concepts/lifecycle/)說明啟動及關閉流程。
- [建置系統](/concepts/build-system/)說明專案工作與前端工具的作用。

## 查閱 API

需要查找型別、方法、選項或命令，而不是學習課程時，請使用[API 參考](/reference/overview/)。若要進行命令列操作，請先參閱[CLI 命令概覽](/reference/cli/)；個別命令的詳細資訊則請參閱[CLI 旗標與範例](/guides/cli/)。

## 取得協助

請查看[常見問題](/faq/)、瀏覽[可運作的範例](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples)，或在[Discord](https://discord.gg/JDdSxwjhGf)中提問。各平台特有的問題集中列於側邊欄的**疑難排解**下方。
