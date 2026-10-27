---
title: "Wails 的運作方式"
description: "瞭解 Wails 架構及其如何實現原生效能"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails 是一個桌面應用程式開發框架，使用<strong>Go 建構後端</strong>，並使用<strong>網頁技術建構前端</strong>。但與 Electron 不同，Wails 不會內附瀏覽器，而是使用<strong>作業系統的原生 WebView</strong>。

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Wails 應用程式

  frontend: 前端
  backend: Go 後端
  os: 作業系統

  Initialisation: 初始化 {
    shape: sequence_diagram
    backend."Serves Static Web App": 提供靜態 Web 應用程式
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": 透過作業系統原生 WebView 呈現網站
  }
  Regular Communication: 一般通訊 {
    shape: sequence_diagram
    frontend."Make API-style call": 發出 API 形式的呼叫
    frontend -> backend.a: JSON
    backend.a."Service processes request": 服務處理要求
    backend.a -> os: 呼叫系統 API
    backend.a."Generate Response": 產生回應
    backend.a -> frontend: JSON
    frontend."Process response": 處理回應
  }
  backend.a.label: a
}
```

**與 Electron 的主要差異：**

| 面向 | Wails | Electron |
| --- | --- | --- |
| **瀏覽器** | 作業系統提供的 WebView | 內附 Chromium（約 100MB） |
| **後端** | Go（編譯執行） | Node.js（直譯執行） |
| **通訊** | 記憶體內橋接器 | IPC（處理程序間通訊） |
| **套件大小** | 約 15MB | 約 150MB |
| **記憶體** | 約 10MB | 約 100MB 以上 |
| **啟動時間** | &lt;0.5 秒 | 2-3 秒 |

## 核心元件

### 1. 原生 WebView

Wails 使用作業系統內建的網頁轉譯引擎：

@tabs{sync-key="platform"}
[Windows]
**WebView2**（Microsoft Edge WebView2）

- 以 Chromium 為基礎（與 Edge 瀏覽器相同）
- 預先安裝於 Windows 10/11
- 透過 Windows Update 自動更新
- 完整支援現代網頁標準

[macOS]
**WebKit**（Safari 的轉譯引擎）

- 內建於 macOS
- 與 Safari 瀏覽器使用相同引擎
- 效能與電池續航力極佳
- 完整支援現代網頁標準

[Linux]
**WebKitGTK**（WebKit 的 GTK 移植版本）

- 透過套件管理員安裝
- 與 GNOME Web（Epiphany）使用相同引擎
- 良好支援網頁標準
- 輕量且效能優異

@end

**這一點為何重要：**

- **不內附瀏覽器** → 應用程式更小
- **作業系統原生** → 整合度與效能更佳
- **自動更新** → 透過作業系統更新取得安全性修補程式
- **熟悉的轉譯結果** → 與系統瀏覽器相同

### 2. Wails 橋接器

橋接器是 Wails 的核心，可讓 Go 與 JavaScript 之間進行<strong>直接通訊</strong>。

```d2
direction: down

Frontend: 前端（JavaScript） {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Wails 橋接器 {
  Encoder: JSON 編碼器 {
    shape: rectangle
  }

  Router: 方法路由器 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 解碼器 {
    shape: rectangle
  }
}

Backend: 後端（Go） {
  Services: 已註冊的服務 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. 呼叫 Go 方法\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. 編碼為 JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. 路由至服務\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. 傳回結果\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. 解碼為 JS\nPromise 完成"
```

**運作方式：**

1. **前端呼叫 Go 方法**（透過自動產生的繫結）
2. <strong>橋接器將呼叫編碼</strong>為 JSON（方法名稱 + 引數）
3. **路由器在已註冊的服務中尋找 Go 方法**
4. <strong>Go 方法執行</strong>並傳回值
5. <strong>橋接器將結果解碼</strong>並傳回前端
6. JavaScript 中的<strong>Promise 以該結果完成</strong>

**效能特性：**

- **記憶體內通訊**：沒有網路額外負擔，也不使用 HTTP
- 盡可能採用<strong>零複製</strong>（適用於大型資料）
- **預設為非同步**：兩端皆不會阻塞
- **型別安全**：自動產生 TypeScript 定義

### 3. 服務系統

建議透過服務將 Go 功能公開給前端。

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**服務探索：**

- Wails 會在啟動時<strong>掃描您的 struct</strong>
- <strong>匯出的方法</strong>可供前端呼叫
- 系統會擷取<strong>型別資訊</strong>，用於 TypeScript 繫結
- <strong>錯誤處理</strong>會自動進行（Go 錯誤 → JS 例外）

**產生的 TypeScript 繫結：**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**為何使用服務？**

- **型別安全**：完整支援 TypeScript
- **自動探索**：無須手動註冊方法
- **條理分明**：將相關功能分組
- **可測試**：服務只是 Go 結構體

[深入瞭解服務 →](/features/bindings/services/)

### 4. 事件系統

事件讓元件之間能透過<strong>發布／訂閱模式進行通訊</strong>。

```d2
direction: left

Wails Event System: Wails 事件系統 {
  shape: sequence_diagram

  window1: 視窗 1
  window2: 視窗 2
  backend: Go 後端

  Event Driver: 事件驅動器 {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "訂閱 'data-updated' 事件"
    window2."Subscribe to 'data-updated' events": "訂閱 'data-updated' 事件"
    backend.a."App Emit('data-updated', data)": "應用程式 Emit('data-updated', data)"
    backend.a -> window1.a: JSON 事件匯流排
    backend.a -> window2: JSON 事件匯流排
    window1.a."Subscriber processes On('data-updated', handler)": "訂閱者透過 On('data-updated', handler) 處理"
    window2."Subscriber processes On('data-updated', handler)": "訂閱者透過 On('data-updated', handler) 處理"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**使用情境：**

- **視窗通訊**：一個視窗通知其他視窗
- **背景工作**：Go 服務向 UI 通知進度
- **狀態同步**：讓多個視窗保持同步
- **鬆散耦合**：元件不需要直接參照彼此

**範例：**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[深入瞭解事件 →](/features/events/system/)

## 應用程式生命週期

瞭解生命週期有助於掌握初始化資源及執行清理的時機。

```d2
direction: down

Start: 應用程式啟動 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初始化 {
  Create: 建立應用程式 {
    shape: rectangle
  }

  Register: 註冊服務 {
    shape: rectangle
  }

  Setup: 設定視窗／選單 {
    shape: rectangle
  }
}

Run: 事件迴圈 {
  Events: 處理事件 {
    shape: rectangle
  }

  Messages: 處理訊息 {
    shape: rectangle
  }

  Render: 更新使用者介面 {
    shape: rectangle
  }
}

Shutdown: 關閉 {
  Cleanup: 清理資源 {
    shape: rectangle
  }

  Save: 儲存狀態 {
    shape: rectangle
  }
}

End: 應用程式結束 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: 迴圈
Run.Events -> Shutdown.Cleanup: 結束訊號
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**生命週期掛鉤：**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options`上沒有`OnStartup`欄位。啟動工作應放在服務的`ServiceStartup(ctx, options)`中、透過`app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`註冊的回呼中，或直接放在`app.Run()`之前。

[深入瞭解生命週期 →](/concepts/lifecycle/)

## 建置流程

瞭解 Wails 如何建置應用程式：

```d2
direction: down

Source: 原始碼 {
  Go: "Go 程式碼\n(main.go, 服務)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "前端程式碼\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: 建置流程 {
  AnalyseGo: 分析 Go 程式碼 {
    shape: rectangle
  }

  GenerateBindings: 產生繫結 {
    shape: rectangle
  }

  BuildFrontend: 建置前端 {
    shape: rectangle
  }

  CompileGo: 編譯 Go {
    shape: rectangle
  }

  Embed: 嵌入資產 {
    shape: rectangle
  }
}

Output: 輸出 {
  Binary: "原生二進位檔\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: 擷取型別
Build.GenerateBindings -> Source.Frontend: TypeScript 繫結
Source.Frontend -> Build.BuildFrontend: 編譯（Vite/webpack）
Build.BuildFrontend -> Build.Embed: 封裝資產
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**建置步驟：**

1. **分析 Go 程式碼**
  - 掃描服務中的匯出方法
  - 擷取參數與傳回型別
  - 產生方法簽章


2. **產生 TypeScript 繫結**
  - 為每個服務建立`.ts`檔案
  - 包含完整的型別定義
  - 新增 JSDoc 註解


3. **建置前端**
  - 執行打包工具（Vite、webpack 等）
  - 縮小並最佳化
  - 輸出至`frontend/dist/`


4. **編譯 Go**
  - 啟用最佳化進行編譯（`-ldflags="-s -w"`）
  - 包含建置中繼資料
  - 針對特定平台進行編譯


5. **嵌入資產**
  - 將前端檔案嵌入 Go 二進位檔
  - 壓縮資產
  - 建立單一可執行檔


<strong>結果：</strong>所有內容均已嵌入的單一原生可執行檔。

[深入瞭解建置 →](/guides/build/building/)

## 開發環境與正式環境

Wails 在開發環境與正式環境中的行為有所不同：

@tabs{sync-key="mode"}
[開發環境（wails3 dev）]
**特性：**

- **熱重載**：前端變更會立即重新載入
- **原始碼對應表**：使用原始程式碼進行偵錯
- **DevTools**：可使用瀏覽器 DevTools
- **記錄**：啟用詳細記錄
- **外部前端**：由開發伺服器（Vite）提供

**運作方式：**

```d2
direction: right

WailsApp: Wails 應用程式 {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Vite 開發伺服器\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: 代理要求
DevServer -> WebView: 透過 HMR 提供內容
WebView -> WailsApp: 呼叫 Go 方法
```

**優點：**

- 立即取得變更回饋
- 完整的偵錯功能
- 加快反覆開發速度

[正式環境（wails3 build）]
**特性：**

- **嵌入式資產**：前端已建置到二進位檔中
- **已最佳化**：經過縮小與壓縮
- **無 DevTools**：預設停用
- **最少量記錄**：僅記錄錯誤
- **單一檔案**：所有內容都在一個可執行檔中

**運作方式：**

```d2
direction: right

Binary: "單一二進位檔\n(myapp.exe)" {
  GoCode: 已編譯的 Go 程式碼 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "嵌入式資產\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: 從記憶體提供內容
WebView -> Binary.GoCode: 呼叫 Go 方法
```

**優點：**

- 單一檔案發佈
- 檔案更小（經過最小化）
- 效能更佳
- 無外部相依項目

@end

## 記憶體模型

瞭解記憶體使用情況，有助於建置高效率的應用程式。

**記憶體區域：**

1. **Go 堆積**
  - 您的服務與應用程式狀態
  - 由 Go 垃圾回收器管理
  - 簡易應用程式通常為 5-10MB


2. **WebView 記憶體**
  - DOM、JavaScript 堆積、CSS
  - 由 WebView 引擎管理
  - 簡易應用程式通常為 10-20MB


3. **橋接層記憶體**
  - 用於通訊的訊息緩衝區
  - 額外負擔極低（<1MB）
  - 大型資料會盡可能採用零複製


**最佳化提示：**

- **避免傳輸大量資料**：傳遞 ID，並視需要擷取詳細資料
- **使用事件進行更新**：不要從前端輪詢
- **以串流方式處理大型檔案**：不要將整個檔案載入記憶體
- **清除監聽器**：使用完畢後移除事件監聽器

[深入瞭解效能 →](/guides/performance/)

## 安全性模型

Wails 提供預設安全的架構：

```d2
direction: down

Frontend: 前端（不受信任） {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails 橋接器（驗證） {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: 後端（受信任） {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: 呼叫方法
Bridge -> Bridge: "驗證：\n- 方法存在嗎？\n- 型別正確嗎？\n- 允許存取嗎？"
Bridge -> Backend: 驗證通過後執行
Backend -> Bridge: 傳回結果
Bridge -> Frontend: 傳送回應
```

**安全性功能：**

1. **方法白名單**
  - 只能呼叫匯出的方法
  - 無法存取私有方法
  - 必須明確註冊服務


2. **型別驗證**
  - 依據 Go 型別檢查引數
  - 拒絕無效的型別
  - 防止注入攻擊


3. **不使用 eval()**
  - 前端無法執行任意 Go 程式碼
  - 只能呼叫預先定義的方法
  - 不會動態執行程式碼


4. **情境隔離**
  - 每個視窗都有自己的情境
  - 服務可以檢查呼叫端的情境
  - 可為每個視窗設定不同權限


**最佳實務：**

- 在 Go 中<strong>驗證使用者輸入</strong>（不要信任前端）
- 使用<strong>情境</strong>進行驗證與授權
- 執行檔案操作前，先<strong>清理檔案路徑</strong>
- 對高成本操作<strong>實施速率限制</strong>

[深入瞭解安全性 →](/guides/security/)

## 後續步驟

**應用程式生命週期** - 瞭解啟動、關閉及生命週期掛鉤 [深入瞭解 →](/concepts/lifecycle/)

**Go 與前端的橋接層** - 深入探討橋接層的運作方式 [深入瞭解 →](/concepts/bridge/)

**建置系統** - 瞭解 Wails 如何建置您的應用程式 [深入瞭解 →](/concepts/build-system/)

**開始建置** - 透過教學實作所學內容 [教學 →](/tutorials/03-notes-vanilla/)

---

<strong>對架構有疑問嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查閱[API 參考](/reference/overview/)。
