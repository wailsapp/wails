---
title: "執行階段內部機制"
description: "深入瞭解 Wails v3 如何啟動、執行及與作業系統通訊"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

**runtime** 是將一般 Go 函式轉換為跨平台桌面應用程式的層級。  
本文件說明追查原始碼時會遇到的各個運作元件。

---

## 1. 應用程式生命週期

| 階段 | 程式碼路徑 | 執行內容 |
| --- | --- | --- |
| **啟動程序** | `pkg/application/application.go:init()` | 註冊建置時資料，並建立全域的 `application` 單一實例。 |
| **New()** | `application.New(...)` | 驗證 `Options`、啟動 **AssetServer**，並初始化記錄功能。 |
| **Run()** | `application.(*App).Run()` | 1. 呼叫平台的 `mainthread.X()` 以進入作業系統 UI 執行緒。<br />2. 啟動 **runtime**（`internal/runtime`）。<br />3. 阻塞至最後一個視窗關閉或 `Quit()` 被呼叫為止。 |
| **關閉程序** | `application.(*App).Quit()` | 廣播 `application:shutdown` 事件、排清記錄，並拆除視窗與服務。 |

此生命週期嚴格採用<strong>單次進入</strong>模式：您可以建立多個視窗，但應用程式物件本身只會初始化一次。

---

## 2. 視窗管理

### 公開 API

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()`**不接受任何引數**；以下情況請使用`NewWithOptions(...)`：當您需要
>
> 以傳值方式傳遞 `application.WebviewWindowOptions` 結構。

`app.Window.New[WithOptions]()` 會委派給 `pkg/application/webview_window_*.go`，各平台專屬實作均位於該處：

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

每個檔案都會：

1. 建立原生 WebView（WKWebView、WebKitGTK、WebView2）。
2. 註冊<strong>訊息處理器</strong>回呼（`pkg/application/messageprocessor*.go`）。
3. 將 Wails 事件（`WindowDidResize`、`WindowFocus`、`WindowFilesDropped`……）對應至 `pkg/events` 中的常數。

`internal/runtime/` 保留給小型的建置標籤黏合程式碼（`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`），以及 `internal/runtime/desktop/` 下的內嵌 JS 執行階段。

使用 `pkg/application/window_manager.go`／`webview_window.go` 追蹤使用中的視窗。`pkg/application/screenmanager.go` 用於<strong>顯示器</strong>中繼資料（解析度、縮放比例、工作區域），不負責管理視窗。

---

## 3. 訊息處理管線

JavaScript 與 Go 之間的橋接由 `pkg/application/messageprocessor_*.go` 中的<strong>訊息處理器</strong>系列實作。

流程：

1. **JavaScript** 從 `/wails/runtime.js` 呼叫 `Call.ByID(<fnv-id>, ...args)`（於 `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` 中實作）；名稱模式建置則呼叫 `Call.ByName("pkg.Struct.Method", ...args)`。
2. 執行階段輔助程式會封裝呼叫，並透過各平台的原生橋接將其分派至 Go。
3. **Go** 在 `pkg/application/messageprocessor_call.go` 中接收訊息。
4. 處理器會在 `pkg/application/bindings.go`（手寫且以 `reflect` 為基礎）中查詢已繫結的方法，然後加以呼叫。
5. 結果或錯誤會被封送回 JS，由 `Promise` 完成或拒絕。

> 確切的 JSON 封套由 JS 端的執行階段輔助程式
>
> 以及 Go 端的 `messageprocessor_call.go` 定義；本頁較早的草稿
>
> 曾描述一種 `{"t":"c","id":"123","m":"Greet","p":[…]}` 結構，但該結構
>
> 與目前的實作相符。請一併閱讀這兩個檔案，以追查
>
> 傳輸格式錯誤。

專用處理器：

| 檔案 | 用途 |
| --- | --- |
| `messageprocessor_window.go` | 視窗操作（隱藏、最大化……） |
| `messageprocessor_dialog.go` | 原生對話方塊（`OpenFile`、`MessageBox`……） |
| `messageprocessor_clipboard.go` | 剪貼簿讀取／寫入 |
| `messageprocessor_events.go` | 事件訂閱／發出 |
| `messageprocessor_browser.go` | 瀏覽器導覽、開發人員工具 |

處理器是<strong>無狀態的</strong>，會從每則訊息隨附的 `ApplicationContext` 取得所需的一切資料。

---

## 4. 事件系統

事件是具有命名空間的字串，會跨三個層級分派：

1. **應用程式事件**：全域生命週期（`application:ready`、`application:shutdown`）。
2. **視窗事件**：每個視窗各自獨立（`window:focus`、`window:resize`）。
3. **自訂事件**：由使用者定義（`chat:new-message`）。

實作細節：

- 事件常數位於`pkg/events/`（`defaults.go`、`known_events.go`、`events.txt`）。這些常數由`v3/tasks/events/generate.go`產生，並公開為`events.Common.*`、`events.Mac.*`、`events.Windows.*`、`events.Linux.*`。你可以使用`wails3 generate constants`重新產生它們。
- Go 端（應用程式事件）：
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Go 端（視窗事件）：
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Go 端（自訂事件）：
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- JS 端：
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


應用程式、視窗和自訂事件都會流經 `pkg/application/event_manager.go`。視窗事件訂閱的作用域限定於 所屬視窗，因此關閉視窗時會自動取消註冊其處理常式。

---

## 5. 平台特定實作

條件式編譯會隱藏各作業系統的差異，同時讓公開 API 保持一致。

| 關注事項 | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| 主執行緒 | `mainthread_darwin.go`（透過 Cgo 呼叫 Foundation） | `mainthread_linux.go`（GTK） | `mainthread_windows.go`（Win32 `AttachThreadInput`） |
| 對話方塊 | `dialogs_darwin.*`（NSAlert） | `dialogs_linux.go`（GtkFileChooser） | `dialogs_windows.go`（IFileOpenDialog） |
| 剪貼簿 | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| 系統匣圖示 | `systemtray_darwin.*` | `systemtray_linux.go`（DBus） | `systemtray_windows.go`（Shell_NotifyIcon） |

核心原則：

- <strong>macOS 和 Windows</strong>僅少量使用 Cgo（主要透過`pkg/mac/`以及`pkg/w32`中的`w32` Win32 包裝函式）。
- **Linux 因實際需要而大量使用 Cgo**——`pkg/application/linux_cgo.go`（約69 KB）加上`linux_cgo_gtk4.{c,go,h}`（約50 KB 以上）會直接驅動 GTK/WebKitGTK。
- 使用<strong>建置標籤</strong>（`//go:build darwin`、`//go:build linux`……），讓各作業系統專用檔案保持易讀。
- `internal/capabilities/`用於各平台的功能旗標，但此框架<strong>不會</strong>匯出`ErrCapability`哨兵值；功能門控是透過平台特定的虛設實作回傳值完成。

---

## 6. 檔案指南

| 檔案 | 需要修改它的情況 |
| --- | --- |
| `internal/runtime/runtime_*.go` | 變更小型建置標籤虛設實作層（開發與正式環境、作業系統特定的銜接程式碼）。 |
| `pkg/application/webview_window_*.go` | 實作新的視窗提示或行為。 |
| `pkg/application/messageprocessor*.go` | 新增可從 JS 呼叫的橋接命令。 |
| `pkg/events/*.go` | 擴充內建事件定義（接著重新執行`wails3 generate constants`）。 |
| `internal/assetserver/*` | 微調開發／正式環境的資產處理方式。 |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | 編輯內嵌 JS 執行階段（呼叫／事件分派、對話方塊、拖曳等）。 |

---

## 7. 偵錯提示

- 設定`Options.LogLevel`（例如`slog.LevelDebug`）並檢查`Options.Logger`輸出；不存在`WAILS_LOG_LEVEL`環境變數。
- `wails3 dev`旗標包括`--config`、`--port`、`-s`（啟用 HTTPS）以及全域`--no-colour`。不存在`-verbose`旗標。
- 在 macOS 上，於`lldb --`下執行，以便及早捕捉 Objective-C 例外。
- 若要處理 Windows 上的 Chromium 問題，請啟用 WebView2 偵錯記錄：`set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. 擴充執行階段

1. 選用：在`internal/capabilities/`中宣告任何新的功能旗標。
2. 使用建置標籤，在每個`pkg/application/*_{darwin,linux,windows}.go`變體中實作此功能。請在不支援此功能的平台上提供虛設實作。
3. 在`pkg/application`中新增公開 API（介面、具體的`WebviewWindow`方法、選項結構等）。
4. 如果 JS 需要呼叫此功能，請註冊新的訊息處理器方法（`pkg/application/messageprocessor*.go`），並在 JS 執行階段加入對應的輔助函式。
5. 若新增事件，請在`pkg/events/`中宣告其常數，並執行`wails3 generate constants`以更新產生的檔案。

遵循這份檢查清單，即可維持跨平台合約的完整性。

---

## 9. 拖放

所有平台上的檔案拖放都採用<strong>以 JavaScript 優先的方式</strong>。原生層會攔截作業系統的拖曳事件，但實際的放置處理與 DOM 互動是在 JavaScript 中進行。

### 流程

1. 使用者從作業系統將檔案拖曳到 Wails 視窗上方
2. 原生層偵測拖曳，並通知 JavaScript 顯示懸停效果
3. 使用者放下檔案
4. 原生層將檔案路徑與座標傳送給 JavaScript
5. JavaScript 尋找放置目標元素（`data-file-drop-target`）
6. JavaScript 將檔案路徑與元素詳細資料傳送給 Go 後端
7. Go 發出包含完整情境資訊的`WindowFilesDropped`事件

### 各平台的實作

| 平台 | 原生層 | 主要挑戰 |
| --- | --- | --- |
| **Windows** | WebView2 的內建拖曳支援 | 座標以 CSS 像素表示，不需要轉換 |
| **macOS** | NSWindow 拖曳委派 | 將視窗相對座標轉換為 WebView 相對座標 |
| **Linux** | GTK3 拖曳訊號 | 必須區分檔案拖曳與內部 HTML5 拖曳 |

### Linux：區分拖曳類型

GTK 與 WebKit 都會嘗試處理拖曳事件。關鍵在於檢查拖曳目標類型：

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

對於內部拖曳，訊號處理常式會傳回`FALSE`（交由 WebKit 處理）；對於檔案拖曳，則傳回`TRUE`（由我們自行處理）。

### 封鎖檔案放置

當`EnableFileDrop`為`false`時，我們仍需防止瀏覽器導向被放下的檔案。各平台的處理方式不同：

- **Windows**：JavaScript 會在拖曳事件上呼叫`preventDefault()`
- **macOS**：JavaScript 會在拖曳事件上呼叫`preventDefault()`\
- **Linux**：GTK 訊號處理常式會在原生層攔截並拒絕檔案拖曳

### 關鍵檔案

| 檔案 | 用途 |
| --- | --- |
| `pkg/application/linux_cgo.go` | GTK 拖曳訊號處理常式（cgo 前置區段中的 C 程式碼） |
| `pkg/application/webview_window_darwin.go` | macOS 拖曳委派 |
| `pkg/application/webview_window_windows.go` | WebView2 訊息處理 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | JavaScript 放置處理 |

### 偵錯

- **Linux**：在 C 程式碼中加入`printf`（別忘了`fflush(stdout)`）
- **Windows**：使用`globalApplication.debug()`
- **JavaScript**：檢查瀏覽器主控台，並啟用偵錯模式

常見問題：

1. **內部 HTML5 拖曳無法運作**：原生處理常式攔截了該拖曳（對非檔案拖曳傳回`FALSE`）
2. **未顯示懸停效果**：未呼叫 JavaScript 處理常式
3. **座標錯誤**：檢查座標空間轉換

---

現在，您已在引導下概覽執行階段的內部機制。請結合這些知識、 <strong>程式碼庫配置</strong>圖，以及<strong>資產伺服器</strong>文件，從容瀏覽程式碼庫、 做出具影響力的貢獻。祝您開發愉快！
