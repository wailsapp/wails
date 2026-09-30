---
title: "從 v2 遷移至 v3"
description: "將 Wails v2 應用程式遷移至 v3 的完整指南"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 是一次<strong>徹底重寫</strong>，在架構、效能和開發者體驗方面都有顯著改進。本指南協助你將 v2 應用程式遷移至 v3。

**主要變更：**

- 全新的應用程式結構
- 改進的繫結系統
- 強化的視窗管理
- 更完善的事件系統
- 簡化的設定

<strong>遷移時間：</strong>一般應用程式約需 1-4 小時

## 破壞性變更

### 應用程式初始化

在 v2 中，應用程式設定、視窗設定和執行全都合併在單一 `wails.Run()` 呼叫中。這種單體式做法讓建立多個視窗、在不同階段處理錯誤，或測試應用程式的個別元件變得困難。

v3 將這些關注點分為不同階段：建立應用程式、建立視窗和執行。這種分離方式讓你能明確控制應用程式生命週期的每個階段，並使程式碼更模組化且更容易測試。

**v2：**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3：**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**這樣做更好的原因：**

- **支援多視窗**：你可以隨時動態建立視窗，不再僅限於啟動時
- **更完善的錯誤處理**：每個階段都能分別驗證並妥善處理錯誤
- **更清楚的程式碼**：分離各階段後，每個階段正在執行的作業一目瞭然
- **更容易測試**：你可以在不執行事件迴圈的情況下測試應用程式設定
- **更有彈性**：在應用程式的整個生命週期中，都可以建立、銷毀及重新建立視窗

### 繫結

在 v2 中，每個已繫結的結構都必須有一個 context 欄位和一個 `startup(ctx)` 方法，以接收執行階段 context。這導致業務邏輯與 Wails 執行階段緊密耦合，使程式碼更難測試和理解。

v3 引入服務模式，讓你的結構完全獨立，不需要儲存執行階段 context。如果服務需要存取應用程式執行個體，會透過相依性注入明確接收，而不是以隱含方式逐層傳遞 context。

**v2：**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3：**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**這樣做更好的原因：**

- **沒有隱含相依性**：服務是普通的 Go 結構，不含隱藏的執行階段相依性
- **更容易測試**：你可以測試服務方法，無須模擬 Wails context
- **更清楚的程式碼**：相依性會明確顯示（以建構函式引數傳入），而不是隱藏在 context 欄位中
- **更完善的組織方式**：服務可以依領域分組，不必全都放在單一 `App` 結構中
- **妥善的初始化**：需要初始化時使用 `ServiceStartup()` 方法，讓初始化作業明確可見

### 執行階段

在 v2 中，所有執行階段作業都必須將 context 傳給 `runtime` 套件中的全域函式。這導致整個程式碼庫都與 context 物件緊密耦合，也讓 API 感覺更偏向程序式，而非物件導向。

v3 以直接呼叫應用程式和視窗物件的方法，取代以 context 為基礎的執行階段。作業會直接在其所影響的物件上呼叫，使程式碼更直覺，也更符合物件導向設計。

**v2：**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3：**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**這樣做更好的原因：**

- **物件導向設計**：方法會在其所影響的物件（視窗、應用程式、選單等）上呼叫
- **意圖更清楚**：`window.SetTitle()` 比 `runtime.WindowSetTitle(ctx, ...)` 更直觀
- **更完善的 IDE 支援**：當方法定義在物件上時，自動完成功能便能正常運作
- **多視窗操作更清楚**：有多個視窗時，你會明確選擇要操作的視窗
- **無須逐層傳遞 context**：你不需要將 context 傳入每個函式

### 前端繫結

在 v2 中，繫結依 Go 套件和結構名稱組織，通常會產生 `wailsjs/go/main/App` 之類的路徑。這種結構無法反映邏輯分組，也使相關功能難以尋找。

v3 依服務名稱和應用程式模組組織繫結，建立更清楚的邏輯結構。繫結會產生至 `bindings` 目錄，並依應用程式名稱和服務名稱組織，讓你更容易瞭解有哪些可用功能。

**v2：**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**這樣做更好的原因：**

- **邏輯化組織**：繫結依服務名稱分組，而不是依 Go 套件結構分組
- **更清楚的匯入路徑**：路徑反映領域邏輯（greetservice），而不是檔案結構（main/App）
- **更容易探索**：您可以依功能瀏覽繫結，而非依技術結構瀏覽
- **一致的命名方式**：以服務為基礎的組織方式與後端架構一致
- **更簡潔的路徑**：不再需要`../wailsjs/go`前綴，只需使用`./bindings`

### 事件

在v2中，事件使用可變參數`interface{}`，且每個事件函式都必須傳入context。事件處理常式收到的是無型別資料，需要手動進行型別斷言，導致事件系統容易出錯且難以偵錯。

v3引入具型別的事件物件，並移除必須傳入context的要求。事件處理常式會收到包含具型別資料的適當事件物件，使事件系統更可靠且更容易使用。

**v2：**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3：**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**這種方式更好的原因：**

- **型別安全**：事件使用正式的事件物件，而非`...interface{}`
- **更容易偵錯**：事件物件包含事件名稱等中繼資料，讓偵錯更容易
- **更清楚的 API**：`app.Event.On()`和`app.Event.Emit()`比執行階段函式更直覺
- **不需要context**：事件直接在應用程式物件上運作，不必逐層傳遞context
- **更簡潔的處理常式**：事件處理常式具有明確的簽章，不再使用可變參數

### 視窗

v2 的每個應用程式僅支援單一視窗。視窗會在啟動時建立，而所有視窗操作都透過執行階段函式執行，並隱含地以該單一視窗為目標。

v3 將原生多視窗支援作為核心功能引入。每個視窗都是一級物件，具有自己的方法和生命週期。您可以在應用程式的整個生命週期中動態建立、管理及銷毀多個視窗。

**v2：**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3：**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**這種方式更好的原因：**

- **多視窗應用程式**：建置具有多個獨立視窗的應用程式（儀表板、偏好設定、工具等）
- **明確的視窗參照**：每個視窗都是可儲存並直接操作的物件
- **動態建立視窗**：可在執行期間隨時建立及銷毀視窗
- **獨立的視窗狀態**：每個視窗都有自己的事件、屬性及生命週期
- **更完善的架構**：視窗管理採用物件導向方式，而非以context為基礎

## 移轉步驟

### 步驟1：更新相依套件

**go.mod：**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**更新：**

```bash
go get github.com/wailsapp/wails/v3@latest
go mod tidy
```

### 步驟2：更新 main.go

**v2：**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3：**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 步驟3：將 App 結構轉換為服務

**v2：**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3：**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### 步驟4：更新執行階段呼叫

**v2：**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3：**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### 步驟5：更新前端

**產生新的繫結：**

```bash
wails3 generate bindings
```

**更新匯入：**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**更新事件處理：**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### 步驟6：更新設定

**v2（wails.json）：**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3（wails.json）：**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## 功能對照

### 對話方塊

**v2：**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3：**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### 選單

**v2：**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3：**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### 系統匣

**v2：**

```go
// Not available in v2
```

**v3：**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## 常見問題

### 問題：找不到繫結

<strong>問題：</strong>遷移後發生匯入錯誤

**解決方案：**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### 問題：Context 錯誤

<strong>問題：</strong>無法使用`ctx`

**解決方案：**

改為儲存 app 參照：

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### 問題：視窗方法無法運作

**問題：**`runtime.WindowSetTitle()`不存在

**解決方案：**

直接使用視窗方法：

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### 問題：事件未觸發

<strong>問題：</strong>事件已註冊，但未收到

**解決方案：**

檢查事件名稱是否完全相符：

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## 測試遷移結果

### 檢查清單

- [ ] 應用程式啟動時沒有錯誤
- [ ] 所有繫結皆可正常運作
- [ ] 事件可正常傳送及接收
- [ ] 視窗可正常開啟及關閉
- [ ] 選單可正常運作（如適用）
- [ ] 對話方塊可正常運作（如適用）
- [ ] 系統匣可正常運作（如適用）
- [ ] 建置流程可正常運作
- [ ] 正式環境建置可正常運作

### 測試命令

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## v3 的優點

### 效能

- **啟動更快**－經過最佳化的初始化流程
- **記憶體用量更低**－高效率的資源使用方式
- **更完善的橋接機制**－呼叫額外負擔為&lt;1毫秒

### 功能

- **多視窗**－原生支援
- **系統匣**－內建功能
- **更完善的事件機制**－具型別且更簡潔的 API
- **服務**－更完善的程式碼組織方式

### 開發人員體驗

- **型別安全**－完整支援 TypeScript
- **更完善的錯誤處理**－清楚的錯誤訊息
- **熱重新載入**－加快開發速度
- **更完善的文件**－完整的指南

## 取得協助

### 資源

- [文件](/quick-start/why-wails/)
- [Discord 社群](https://discord.gg/JDdSxwjhGf)
- [GitHub 議題](https://github.com/wailsapp/wails/issues)
- [範例](https://github.com/wailsapp/wails/tree/master/v3/examples)

### 常見問題

**問：可以同時執行 v2 和 v3 嗎？** 答：可以，兩者使用不同的匯入路徑。

**問：v3 已可用於正式環境嗎？** 答：v3 是測試版軟體，並提供穩定的桌面 API。已有應用程式使用它在 正式環境中執行，但部署前仍須進行完整測試。v2 仍是 目前的穩定版本。

**問：v2會繼續維護嗎？** 答：會，v2將持續收到關鍵更新。

**問：遷移需要多久？** 答：一般應用程式需要1-4小時。

## 後續步驟

@cards{cols="2"}
🚀 快速入門
開始使用 Wails v3。

[深入瞭解 →](/quick-start/installation/)

---
★ 核心概念
瞭解 v3 架構。

[深入瞭解 →](/concepts/architecture/)

---
◆ 繫結
瞭解新的繫結系統。

[深入瞭解 →](/features/bindings/methods/)

---
📖 範例
查看完整的 v3 範例。

[查看範例 →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或[建立議題](https://github.com/wailsapp/wails/issues)。
