---
title: "擴充 Wails"
description: "在 Wails v3 中新增功能與平台的實務指南"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails 的設計宗旨是讓您能夠<strong>自由改造</strong>。
>
> 每個主要子系統都以 Go 程式碼實作，您可以閱讀、修改並發布這些程式碼。
>
> 當您執行下列工作時，本頁會說明該從<em>何處</em>著手，以及<em>如何</em>維持跨平台支援：

- 新增<strong>服務</strong>（通知、KV 儲存區、自訂 IPC……）
- 建立<strong>新的 CLI 命令</strong>（`wails3 <foo>`）
- 擴充<strong>執行階段</strong>（視窗 API、對話方塊、事件）
- 導入<strong>平台功能</strong>（Wayland……）
- 在不被`//go:build`標籤淹沒的情況下，維持<strong>跨平台相容性</strong>

---

## 1. 新增服務

v3 中的「服務」是使用者提供的 Go 型別，透過`application.Options.Services`註冊，並經由產生的繫結公開給 JS。v3 程式碼庫隨附：

- `internal/service/` — `wails3 generate service`的鷹架：
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — 可立即註冊使用的現成服務（notifications、kvstore、sqlite、log、fileserver、dock……）。

舊版草稿提到的產生器與 CLI 檔案`internal/service/template/template.go`及`internal/generator/collect/services.go`並不存在；鷹架產生器是`internal/service/service.go`（進入點為`service.Install`），服務的繫結中繼資料則收集於`internal/generator/collect/service.go`。

### 1.1 定義服務

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 實作生命週期介面（選用）

服務可選擇實作下列介面（來自`pkg/application`）：

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **重要：**`ServiceShutdown`**不接受任何引數**。具有下列
>
> 簽章`ServiceShutdown(ctx context.Context) error`的方法<strong>不會</strong>實作
>
> 該介面，而且永遠不會被呼叫，也不會顯示任何提示。

### 1.3 向應用程式註冊服務

沒有全域的`services.Register(...)`呼叫。服務是在執行階段透過`application.Options.Services`註冊：

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

註冊後，`wails3 generate bindings`會在`frontend/bindings/<your import path>/...`下產生包裝已匯出方法的 ES 模組。

### 1.4 從 JS 呼叫

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

v3 中沒有全域的`window.backend.*`；呼叫會經過產生的 ES 模組，而這些模組接著會呼叫`/wails/runtime.js`中的`Call.ByID(...)`。

---

## 2. 撰寫新的 CLI 命令

v3 CLI 使用<strong>`github.com/leaanthony/clir`</strong>（而非 cobra）。接線設定位於`v3/cmd/wails3/main.go`：

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

沒有以`init()`為基礎的自動註冊機制。請將新的子命令加入`cmd/wails3/main.go`，並在`internal/commands/`中加入對應的底層函式（若命令接受選項，也要在`internal/flags/`下加入旗標結構）。重新建置 CLI：

```
cd v3
go install ./cmd/wails3
wails3 hello
```

如果您的命令需要 Taskfile 銜接程式碼，請重複使用`internal/commands/task_wrapper.go`中的輔助函式（`wrapTask("yourtask", args)`）。

---

## 3. 修改執行階段

常見原因：

- 新的視窗功能（`SetOpacity`、`Shake`……）
- 額外的對話方塊（`ColorPicker`）
- 系統層級 API（螢幕亮度）

### 3.1 公開 API

將方法加入`pkg/application/webview_window.go`（介面位於`window.go`）：

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

使用現有的`InvokeSync`/`InvokeAsync`輔助函式，確保呼叫在主執行緒上執行。

### 3.2 訊息處理器

如果 JS 需要叫用新方法，請擴充適當的`pkg/application/messageprocessor_*.go`檔案。訊息處理器使用`MessageProcessor`上以 switch 為基礎的方法，而非全域的`register(...)`呼叫：

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

程式碼庫中<strong>沒有</strong>`messageprocessor_window_opacity.go`檔案，也沒有以`init()`為基礎的`register(MsgSetOpacity, ...)`模式。

### 3.3 平台實作

在`pkg/application/`下各作業系統專用的檔案中加入實作：

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

如果某個平台無法支援該功能，請撰寫不執行任何操作的 stub。此框架沒有`ErrCapability`哨兵值；請在文件中說明支援情況，並視需要透過`Options`上的相關布林欄位／平台專用選項結構公開支援狀態。

### 3.4 功能旗標（選用）

`internal/capabilities/`套件用於宣告各平台的功能集合。沒有公開的`application.HasCapability`／`application.CapOpacity` API。如果您希望某項功能可在執行階段檢查，請將它加入`internal/capabilities/`下，並從`pkg/application`公開具型別的 getter。

---

## 4. 新增平台功能

範例：Linux 上的選用 Wayland 支援。

1. 將相關的`pkg/application/*_linux.go`檔案拆分為`*_linux_x11.go`（`//go:build linux && !wayland`）和`*_linux_wayland.go`（`//go:build linux && wayland`）。
2. 讓使用者透過`wails3 build --tags wayland`選擇啟用。請經由`internal/commands/task_wrapper.go`中現有的`EXTRA_TAGS`銜接機制轉送額外標籤。`dev`層級沒有`--tags wayland`旗標；`wails3 dev`僅接受`--config`、`--port`和`-s`。
3. 更新文件以及`pkg/application/`下所有平台專用的 README。

> 將預設建置標籤維持在最低限度；只有小眾功能才使用選擇啟用標籤。

---

## 5. 跨平台相容性檢查清單

| ✅ 步驟 | 原因 |
| --- | --- |
| 在所有平台檔案中提供<strong>每個</strong>公開方法（即使只是虛設常式） | 確保每個作業系統上的建置都能成功 |
| 記錄各作業系統的優雅降級行為 | 應用程式可依據`runtime.GOOS`進行分支處理，而不會出現隱藏錯誤 |
| 優先使用<strong>純 Go</strong>，僅在必要時使用 Cgo | 簡化交叉編譯（Linux 已經承擔了使用 Cgo 的代價） |
| 執行`task test:cli`、`task test:generator`、`task test:templates` | 在本機重現 CI 流程 |
| 在貢獻者文件／範本 README 中記錄新的建置標籤 | 使用者必須知道有哪些功能需要選擇啟用 |

---

## 6. 偵錯建置與迭代速度

- 使用`Options.LogLevel = slog.LevelDebug`（`Options.Logger = slog.Default()`）輸出詳細的執行階段活動。不存在`WAILS_LOG_LEVEL`環境變數。
- `wails3 dev`旗標包括`--config`、`--port`、`-s`。不存在`-race`或`-verbose`旗標；請使用`go test -race ./...`，或對應用程式執行`go build -race`後直接執行，以運用競態偵測器。
- 競態／Cgo 測試指南位於`v3/TESTING.md`（較舊的草稿指向了不存在的`pkg/application/RACE.md`）。

---

## 7. 上游貢獻

1. 若要新增功能或變更公開行為，請建立<strong>WEP（Wails Enhancement Proposal）</strong>草稿 PR，以討論構想與設計。只有可重現的錯誤或文件問題才應使用 issue。
2. 請依照上述方法實作。
3. 請新增：
  - 單元測試（`*_test.go`）
  - 文件（此檔案或相關的`docs/...`頁面）
  - 若修改了繫結產生器，請在`internal/generator/testcases/`下新增迴歸測試

4. 推送前，請先在本機執行`task precommit`以及相關的`task test:*`目標。

---

### 快速連結

| 領域 | 位置 |
| --- | --- |
| 內建服務 | `pkg/services/` |
| 服務鷹架產生器 | `internal/service/` |
| CLI 配線 | `v3/cmd/wails3/main.go` |
| CLI 命令主體 | `internal/commands/` |
| 各作業系統的執行階段 | `pkg/application/*_{darwin,linux,windows}.go` |
| 功能宣告 | `internal/capabilities/` |
| Taskfile DSL | `v3/Taskfile.yaml` |
| 事件常數產生器 | `v3/tasks/events/generate.go` |

---

現在，你已經有一份讓 Wails 隨你所用的<strong>路線圖</strong>：新增服務、 為 CLI 增添巧思、修改執行階段，或加入全新的作業系統功能。 祝擴充順利！
