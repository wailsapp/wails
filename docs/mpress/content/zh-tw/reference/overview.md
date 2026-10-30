---
title: "API 參考資料"
description: "查詢 Wails 型別、方法、選項與 CLI 命令"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## 關於本參考資料

開發應用程式時，可使用本節查詢 Wails 型別、方法、選項與命令。如需引導式教學，請從[教學](/tutorials/overview/)開始。如需完成特定工作的步驟，請參閱[操作指南](/quick-start/next-steps/#features)。如需瞭解設計背景，請參閱[核心概念](/concepts/architecture/)。

**內容架構：**

- [應用程式](/reference/application/)－核心應用程式 API
- [視窗](/reference/window/)－視窗建立與管理；[視窗選項](/features/windows/options/)列出各項設定欄位
- [選單](/reference/menu/)－應用程式、快顯選單與系統匣 API；[選單項目參考資料](/features/menus/reference/)涵蓋項目型別與屬性
- [事件](/reference/events/)－事件系統與內建事件
- [對話方塊](/reference/dialogs/)－檔案與訊息對話方塊
- [前端執行階段](/reference/frontend-runtime/)－前端執行階段 API
- [CLI 命令概覽](/reference/cli/)－依用途分類的命令；[旗標與範例](/guides/cli/)提供各個命令的詳細資訊
- [更新資訊清單通訊協定](/reference/update-manifest/)－更新中繼資料與驗證欄位
- [行動裝置 API](/guides/mobile/mobile-api/)－實驗性行動裝置 API

## API 慣例

@details{title="Go API 慣例－適合剛接觸 Go 的開發人員"}
### 命名

- <strong>型別</strong>：PascalCase（例如`WebviewWindow`）
- <strong>方法</strong>：PascalCase（例如`SetTitle()`）
- <strong>選項</strong>：PascalCase 結構（例如`WindowOptions`）
- <strong>常數</strong>：PascalCase（例如`WindowStartStateMaximised`）

#### 錯誤處理

大多數可能失敗的方法都會將`error`作為最後一個傳回值。`app.Run()`會封鎖至應用程式結束，並傳回任何啟動錯誤：

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

建構視窗不會傳回錯誤；`app.Window.New()`會直接傳回`*WebviewWindow`。

#### 內容

服務生命週期方法會接收`context.Context`：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

可透過`app.Context()`取得應用程式的生命週期內容。並沒有`RunWithContext`；請呼叫`app.Run()`。

#### 選項模式

設定使用選項結構：

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### JavaScript API 慣例

#### 命名

- **函式**：camelCase（例如`setTitle()`）
- **常數**：SCREAMING<em>SNAKE</em>CASE（例如`WINDOW_EVENT_FOCUS`）

#### 預設採用非同步

所有 Go 方法呼叫都會傳回 Promise：

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### 錯誤處理

Go 錯誤會轉換成 JavaScript 例外狀況：

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### 型別安全

TypeScript 定義會自動產生：

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## 套件結構

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## 匯入路徑

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## 型別參考資料

### 常用型別

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## 平台差異

部分 API 在不同平台上的行為有所差異：

| 功能 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **應用程式選單** | 視窗選單列 | 全域選單列 | 視窗選單列 |
| **系統匣** | 通知區域 | 選單列 | 系統匣 |
| **Dock** | 不適用 | ✅ 可用 | 不適用 |
| **檔案對話方塊** | 原生 | 原生 | 原生（GTK） |
| **透明效果** | ✅ 完整支援 | 需要[`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ 支援有限 |

各 API 章節均記載平台特有的行為。

## 版本管理

Wails v3遵循語意化版本控制：

- **主要版本**（v4.0.0）：破壞性變更
- **次要版本**（v3.1.0）：新增功能，向下相容
- **修訂版本**（v3.0.1）：錯誤修正，向下相容

<strong>Wails v3.0.0：</strong>穩定版

## 棄用政策

API 棄用時：

1. **在文件中標示**棄用通知
2. **提供替代方案**及移轉指南
3. **移除前會維護1一個主要版本**
4. **編譯器警告**（若可行）

## API 穩定性

### 穩定版 API ✅

這些 API 已達穩定狀態，可安全用於正式環境：

- 核心應用程式 API
- 視窗管理
- 選單系統
- 事件系統
- 檔案對話方塊
- 服務繫結

### 實驗性 API ⚠️

只有明確記載為實驗性的 API 不受穩定桌面相容性契約保障：

- Android 支援
- iOS 支援
- 其他明確標示為實驗性的 API

實驗性 API 會在文件中標示。

## 取得協助

### API 問題

1. **查閱此參考文件** - 完整的 API 文件
2. **查看範例** - [GitHub 範例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **搜尋 Discord** - [Discord 伺服器](https://discord.gg/JDdSxwjhGf)
4. **詢問社群** - Discord #help 頻道

### 回報 API 問題

發現錯誤或不一致之處？

1. **查看現有問題** - [GitHub Issues](https://github.com/wailsapp/wails/issues)
2. **建立詳細報告** - 請附上程式碼、錯誤訊息和平台資訊
3. **提供重現方式** - 能呈現問題的最小範例

## 相關文件

- [教學](/tutorials/overview/) - 透過建置實際應用程式來學習
- [指南](/guides/architecture/) - 適用於常見情境、以工作為導向的指南
- [功能](/features/windows/basics/) - 逐項介紹各項功能的文件
- [範例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) - GitHub 上可執行的程式碼範例

---

<strong>瀏覽 API：</strong>使用左側導覽列探索特定 API。
