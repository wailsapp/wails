---
title: "LLM 控制（MCP）"
description: "讓 LLM 代理透過模型上下文協定測試及控制您的應用程式"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="實驗性功能"}
內建的 MCP 伺服器仍屬實驗性功能，其 API 可能會在未來版本中變更。

@end

Wails v3 內建[模型上下文協定](https://modelcontextprotocol.io)（MCP）伺服器，讓 LLM 代理（例如 Claude Code、IDE 助理或任何 MCP 用戶端）能檢查、測試及操控<strong>執行中</strong>的 Wails 應用程式。

若要自動化專案生命週期，請使用獨立的`wails3 mcp` CLI 伺服器。它可讓代理檢查及初始化專案、執行診斷、啟動建置與開發工作、產生繫結、執行具名的 Taskfile 工作，以及擷取有界限的工作輸出。CLI 伺服器預設僅限存取目前目錄，且不會開放任意 Shell 指令執行功能。如需傳輸、驗證及工具的詳細資訊，請參閱[CLI MCP 文件](/guides/cli/#mcp)。

啟用後，連線至您應用程式的代理可以：

- **列出及控制視窗** — 大小、位置、焦點、全螢幕、開發人員工具、重新載入等
- **檢查 DOM** — 查詢元素、取得 HTML、擷取結構快照
- **執行 JavaScript** — 在任何視窗內執行任意程式碼並取得結果
- **模擬使用者輸入** — 以<strong>螢幕動畫游標</strong>呈現滑鼠移動、點擊、拖曳及捲動，讓您可以觀看代理的操作過程
- **輸入文字及按鍵** — 產生真實的逐字元事件，並可搭配 React 受控輸入元件使用
- **呼叫已繫結的 Go 方法**，以及發出或等待應用程式事件

## 運作方式

只有在使用<strong>`mcp`建置標籤</strong>時，MCP 伺服器才會編譯至您的應用程式中。若未使用該標籤，伺服器程式碼將完全不會出現在二進位檔中 — 不會產生執行階段負擔、不會開啟任何連接埠，也不會增加攻擊面。

使用該標籤時，伺服器會在`App.Run()`內自動啟動，預設繫結至`127.0.0.1:9099`，並將其端點寫入記錄。不需要任何使用者程式碼。

## 教學課程

### 步驟1 — 撰寫一般的 Wails 應用程式

MCP 不需要匯入套件或註冊。請依平常的方式建立應用程式：

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 步驟2 — 使用`mcp`標籤建置或執行

@tabs
[Wails CLI（建議）]
設定`WAILS_MCP=1`後，Wails CLI 會自動為您加入`mcp`標籤：

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[直接使用 Go]
直接將標籤傳給`go run`或`go build`：

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows（PowerShell）]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

應用程式啟動時會記錄 MCP 端點：

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### 步驟3 — 連線用戶端

此伺服器採用<strong>MCP 可串流 HTTP 傳輸</strong>。您可以使用任何相容 MCP 的用戶端連線。

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

接著，請要求 Claude 與您的應用程式互動：

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code（GitHub Copilot）]
加入至`.vscode/settings.json`：

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[其他用戶端]
將任何支援可串流 HTTP 傳輸的 MCP 用戶端指向：

```
http://127.0.0.1:9099/mcp
```

@end

### 步驟4 — 執行測試工作階段

要求代理操作您的應用程式。以下是一些提示詞範例：

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## 設定

所有設定皆透過環境變數完成，不需要變更程式碼。

| 環境變數 | 預設值 | 說明 |
| --- | --- | --- |
| `WAILS_MCP` | （未設定） | 使用 Wails CLI 時，設為`1`、`true`、`on`或`yes`，即可自動加入`mcp`建置標籤。 |
| `WAILS_MCP_HOST` | `127.0.0.1` | 要繫結的網路介面。繫結至非回送介面時需要`WAILS_MCP_TOKEN`。 |
| `WAILS_MCP_TOKEN` | 未設定 | 回送介面可選用持有人權杖；其他繫結位址則必須使用。用戶端傳送`Authorization: Bearer <token>`。 |
| `WAILS_MCP_PORT` | `9099` | 要接聽的連接埠。設為`0`即可隨機指派可用連接埠（會顯示於記錄中）。 |
| `WAILS_MCP_TIMEOUT` | `30000` | 預設的 JS 執行逾時，單位為<strong>毫秒</strong>。 |
| `WAILS_MCP_HIDE_CURSOR` | （未設定） | 設為`1`或`true`可停用動畫游標疊加層。 |

範例 — 自訂連接埠及60秒逾時：

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## 可用工具

| 工具 | 用途 |
| --- | --- |
| `app_info` | 應用程式資訊：平台、架構、所有視窗、MCP 端點 |
| `windows_list` | 列出所有視窗及其幾何資訊與狀態 |
| `window_control` | 聚焦、調整大小、移動、全螢幕、開發人員工具、重新載入、設定 URL 等（22個動作） |
| `js_eval` | 在視窗中執行 JavaScript（非同步主體，使用`return`取得值） |
| `dom_html` | 取得頁面或特定元素的 HTML |
| `dom_query` | 使用 CSS 選擇器尋找元素 — 標籤、文字、邊界、可見性 |
| `screenshot_dom` | 可見頁面的結構快照（以 DOM 為基礎，不含像素） |
| `mouse_move` | 以動畫方式將游標移至某個點或 CSS 選擇器 |
| `mouse_click` | 使用動畫游標點按（左鍵／右鍵／中鍵、按兩下、輔助按鍵） |
| `mouse_drag` | 使用動畫游標拖曳（支援 HTML5 拖放元素） |
| `mouse_scroll` | 在某個點或元素上捲動 |
| `keyboard_type` | 逐字元輸入文字並產生逼真的事件 |
| `keyboard_press` | 按下單一按鍵（Enter、Tab、Escape、ArrowDown……），可選擇搭配輔助按鍵 |
| `call_bound_method` | 呼叫已繫結的 Go 服務方法，例如`main.GreetService.Greet` |
| `emit_event` | 發出 Wails 應用程式事件 |
| `wait_for_event` | 等待 Wails 應用程式事件並傳回其資料 |

### 多視窗支援

所有會對視窗執行操作的工具都接受選用的`window`引數，其中包含視窗的<strong>名稱</strong>（透過`WebviewWindowOptions.Name`設定）。若省略此引數，工具會以目前取得焦點的視窗為目標；若沒有視窗取得焦點，則以第一個視窗為目標。

```
List all windows, then click the "New" button in the window named "editor".
```

### 選取元素

滑鼠和鍵盤工具接受下列任一項：

- **CSS 選擇器** — `selector: "#submit-btn"`（自動捲動至元素所在位置）
- **座標** — `x: 400, y: 300`（相對於檢視區的 CSS 像素）

進行拖曳操作時，請加上`from_`和`to_`前綴：

```
Drag from selector: ".card" to selector: ".dropzone"
```

## 安全性

@note{type="caution"}
MCP 伺服器可讓您以程式完整控制應用程式。任何能夠存取其工具的人都可以讀取 DOM、求值 JavaScript、點按按鈕及呼叫 Go 方法。

@end

- 伺服器預設繫結至`127.0.0.1`。瀏覽器來源必須是 HTTP(S) 回送來源；不透明（`null`）、格式錯誤及外部來源都會遭到拒絕。
- 仍支援不含標頭的原生 MCP 用戶端。未設定`WAILS_MCP_TOKEN`時，系統會信任本機處理程序及獲准的本機瀏覽器來源；來源檢查並非驗證。請設定高熵權杖，要求所有`/mcp`呼叫都必須使用持有人驗證，並在用戶端的`Authorization: Bearer <token>`標頭中設定相同權杖。預檢要求不需要此權杖。
- `/eval-result`回呼使用每次求值時都無法預測的 ID，而非用戶端持有人權杖，因此仍與 WebView 結果傳遞相容。
- 正式環境組建<strong>不應</strong>包含`mcp`標籤。Wails CLI 只會在明確設定`WAILS_MCP=1`時加入此標籤，而預設的`wails3 build`完全不含伺服器程式碼。
- 若需要在非回送介面上公開伺服器（例如進行區域網路測試），請設定`WAILS_MCP_HOST=0.0.0.0`及高熵`WAILS_MCP_TOKEN`；若未設定權杖，啟動將會失敗。由於內建接聽程式使用 HTTP，請在不受信任的網路上使用加密通道或終止 TLS 的 Proxy。

## 範例應用程式

[`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp)提供了展示所有工具的完整遊樂場應用程式。  
其中包括：

- 附有遞增／重設按鈕的計數器
- 附有 Greet、Add、Shout 繫結方法的名稱輸入欄位
- HTML5 拖放來源和目標
- 可捲動清單（50個項目）
- 事件記錄

使用下列命令執行：

```shell
cd v3/examples/mcp
go run -tags mcp .
```

接著，將 Claude Code 或任何 MCP 用戶端連線至`http://127.0.0.1:9099/mcp`，並要求它操作此 UI。

## 意見回饋

內建 MCP 伺服器是一項實驗，其未來方向將由您的意見回饋決定。如果您試用了它，我們很想知道您使用了哪些用戶端和工具、您的預期與實際結果，以及讓代理程式操作您的應用程式是否有幫助 — 最有用的報告會明確說明您實際執行的內容。請到[MCP 伺服器意見回饋討論](https://github.com/wailsapp/wails/discussions/5692)告訴我們。
