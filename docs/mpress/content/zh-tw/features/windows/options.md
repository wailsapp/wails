---
title: "視窗選項"
description: "視窗選項、預設值與平台限制"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## 視窗設定選項

Wails 提供完整的視窗設定功能，包含數十個用於尺寸、位置、外觀與行為的選項。本參考資料說明 Windows、macOS 與 Linux 共用的常見`WebviewWindowOptions`欄位。以下結構僅摘要列出部分欄位；如需完整型別，請參閱所連結的 Go 原始碼。

## WebviewWindowOptions 結構

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window Window)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions` 沒有 `Parent` 欄位——若要建立父視窗／模態視窗關係，請使用 `parentWindow.AttachModal(childWindow)`。它也沒有 `Assets` 欄位——資產設定位於 `application.Options`（`Assets AssetOptions`）中。

完整原始碼：[`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/pkg/application/webview_window_options.go)。

## 核心選項

### 名稱

<strong>型別：</strong>`string` <strong>預設值：</strong>`window-<numeric ID>` <strong>平台：</strong>全部

```go
Name: "main-window"
```

<strong>用途：</strong>用於稍後尋找視窗的唯一識別碼。

**最佳實務：**

- 使用具描述性的名稱：`"main"`、`"settings"`、`"about"`
- 使用 kebab-case：`"file-browser"`、`"color-picker"`
- 保持簡短且容易記憶

**範例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### 標題

<strong>型別：</strong>`string` <strong>預設值：</strong>應用程式名稱 <strong>平台：</strong>全部

```go
Title: "My Application"
```

<strong>用途：</strong>顯示於標題列和工作列中的文字。

**動態更新：**

```go
window.SetTitle("My Application - Document.txt")
```

### 寬度／高度

<strong>型別：</strong>`int`（像素）<strong>預設值：</strong>800 x 600 <strong>平台：</strong>全部 <strong>限制：</strong>必須為正數

```go
Width:  1200,
Height: 800,
```

<strong>用途：</strong>以邏輯像素指定視窗的初始尺寸。

**注意事項：**

- Wails 會自動處理 DPI 縮放
- 使用邏輯像素，而非實體像素
- 考量最低螢幕解析度（1024x768）

**尺寸範例：**

| 使用情境 | 寬度 | 高度 |
| --- | --- | --- |
| 小型公用程式 | 400 | 300 |
| 標準應用程式 | 1024 | 768 |
| 大型應用程式 | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X／Y

<strong>型別：</strong>`int`（像素）<strong>預設值：</strong>置於螢幕中央 <strong>平台：</strong>全部

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

<strong>用途：</strong>視窗的初始位置。

**座標系統：**

- （0, 0）是主要螢幕的左上角
- X 的正方向為向右
- Y 的正方向為向下

**範例：**

`X`和`Y`僅在設定`InitialPosition: application.WindowXY`後才會生效。若未設定，`InitialPosition`預設為`WindowCentered`，而`X`／`Y`會被忽略。

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

<strong>最佳實務：</strong>如果不在意具體座標，請使用`Center()`在建立視窗後將其置中：

```go
window := app.Window.New()
window.Center()
```

### 最小寬度／最小高度

<strong>型別：</strong>`int`（像素）<strong>預設值：</strong>0（無最小值）<strong>平台：</strong>全部

```go
MinWidth:  400,
MinHeight: 300,
```

<strong>用途：</strong>防止視窗過小。

**使用情境：**

- 避免版面配置損壞
- 確保可用性
- 維持長寬比

**範例：**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### 最大寬度／最大高度

<strong>型別：</strong>`int`（像素）<strong>預設值：</strong>0（無最大值）<strong>平台：</strong>全部

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

<strong>用途：</strong>防止視窗過大。

**使用情境：**

- 固定尺寸的應用程式
- 避免使用過多資源
- 維持設計限制

## 狀態選項

### 隱藏

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>全部

```go
Hidden: true,
```

<strong>用途：</strong>建立視窗但不顯示。

**使用情境：**

- 背景視窗
- 按需顯示的視窗
- 啟動畫面（建立、載入，然後顯示）
- 避免載入內容時閃現白色畫面

**各平台的改進：**

- <strong>Windows：</strong>已修正白色視窗閃現問題——呼叫`Show()`之前，視窗會維持不可見
- <strong>macOS：</strong>完整支援
- <strong>Linux：</strong>完整支援

**建議採用的流暢載入模式：**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**範例：**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### 無框線

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>全部

```go
Frameless: true,
```

<strong>用途：</strong>移除標題列與視窗邊框。

**使用情境：**

- 自訂視窗裝飾
- 啟動畫面
- 資訊站應用程式
- 自訂設計的視窗

<strong>重要：</strong>您需要實作：

- 視窗拖曳
- 關閉／最小化／最大化按鈕
- 調整大小控點（若可調整大小）

**詳情請參閱[無框線視窗](/features/windows/frameless/)。**

### DisableResize

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false`（視窗預設可調整大小）<strong>平台：</strong>全部

```go
DisableResize: true,
```

**用途：**防止調整視窗大小。請注意，此欄位與v2的`Resizable`呈**相反**關係；將`DisableResize: true`設為 true，即可讓視窗無法調整大小。

**使用情境：**

- 固定大小的應用程式
- 啟動畫面
- 對話方塊

<strong>注意：</strong>除非您也透過`MaximiseButtonState`／集合行為停用最大化與全螢幕，否則使用者仍可使用這些功能。

### AlwaysOnTop

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>全部

```go
AlwaysOnTop: true,
```

<strong>用途：</strong>讓視窗保持在所有其他視窗上方。

**使用情境：**

- 浮動工具列
- 通知
- 子母畫面
- 計時器

**平台注意事項：**

- <strong>macOS：</strong>完整支援
- <strong>Windows：</strong>完整支援
- <strong>Linux：</strong>視視窗管理員而定

### StartState

<strong>型別：</strong>`WindowState`列舉<strong>預設值：</strong>`WindowStateNormal` <strong>平台：</strong>全部

```go
StartState: application.WindowStateMaximised,
```

<strong>用途：</strong>視窗顯示時的初始狀態。

**值：**

- `WindowStateNormal`－一般視窗
- `WindowStateMinimised`－最小化
- `WindowStateMaximised`－最大化
- `WindowStateFullscreen`－全螢幕

不存在`WindowStateHidden`常數；若要讓視窗啟動時不可見，請使用`Hidden`布林欄位。

**在執行階段切換全螢幕：**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## 外觀選項

### BackgroundColour

<strong>型別：</strong>`RGBA`結構<strong>預設值：</strong>白色<strong>平台：</strong>全部

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

`RGBA`欄位的型別為`Red, Green, Blue, Alpha`（uint8）。建議使用輔助函式`application.NewRGB(r, g, b)`（Alpha 值為255）或`application.NewRGBA(r, g, b, a)`。

<strong>用途：</strong>內容載入前的視窗背景色彩。

**使用情境：**

- 配合應用程式的佈景主題
- 避免深色佈景主題出現白色閃爍
- 提供流暢的載入體驗

**範例：**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**輔助方法：**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

<strong>型別：</strong>`BackgroundType`列舉<strong>預設值：</strong>`BackgroundTypeSolid` <strong>平台：</strong>macOS、Windows（部分支援）

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**值：**

- `BackgroundTypeSolid`－純色
- `BackgroundTypeTransparent`－完全透明
- `BackgroundTypeTranslucent`－半透明模糊

**平台支援：**

- <strong>macOS：</strong>設定`Mac.Backdrop`；若要讓 WebView 透明，必須設定[`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background)。否則 WebView 仍會保持不透明。
- <strong>Windows：</strong>透明與半透明（Windows 11+）
- <strong>Linux：</strong>僅支援純色

**範例（macOS）：**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup 與 OpenDevTools

<strong>macOS 上的私有 API：</strong>若要以程式方式開啟檢查器，`OpenInspectorOnStartup: true`、Go `window.OpenDevTools()`及 JavaScript `Window.OpenDevTools()`都需要`private_mac_apis`。若未設定，這些操作不會產生任何作用。正式環境組建還需要`devtools`。macOS 13.3+ 上公開的 Safari 檢查功能不需要私有 API；在較舊版本的 macOS 上啟用檢查器則需要。請參閱[Web 檢查器組建矩陣](/guides/build/private-macos-apis/#web-inspector)。

## 內容選項

### URL

<strong>型別：</strong>`string` <strong>預設值：</strong>空白（從 Assets 載入）<strong>平台：</strong>全部

```go
URL: "https://example.com",
```

<strong>用途：</strong>載入外部 URL，而非內嵌資產。

**使用情境：**

- 開發（從開發伺服器載入）
- 網頁式應用程式
- 混合式應用程式

**範例：**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

正式環境情境的應用程式層級程式碼片段：

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

<strong>型別：</strong>`string` <strong>預設值：</strong>空白<strong>平台：</strong>全部

```go
HTML: "<h1>Hello World</h1>",
```

<strong>用途：</strong>直接載入 HTML 字串。

**使用情境：**

- 簡單視窗
- 產生的內容
- 測試

**範例：**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### 資產（僅限應用程式層級）

資產設定**不是**`WebviewWindowOptions`的欄位。前端資產由應用程式本身透過`application.Options.Assets`（`AssetOptions`）提供；每個視窗都會繼承該資產伺服器。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**如需詳細資訊，請參閱[建置系統](/concepts/build-system/)。**

### UseApplicationMenu

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>Windows、Linux（在 macOS 上無作用）

```go
UseApplicationMenu: true,
```

<strong>用途：</strong>讓此視窗使用應用程式選單（透過`app.Menu.Set()`設定）。

在**macOS**上，此選項不會產生任何作用，因為 macOS 一律使用位於畫面頂端的全域應用程式選單。

在**Windows**和**Linux**上，視窗預設不會顯示選單。設定`UseApplicationMenu: true`會讓視窗使用應用程式層級選單，提供簡單的跨平台解決方案。

**範例：**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**注意事項：**

- 如果同時設定`UseApplicationMenu`和視窗專用選單，會以視窗專用選單為優先
- 這樣便不需要在執行階段檢查作業系統，從而簡化跨平台程式碼
- 如需完整的選單文件，請參閱[應用程式選單](/features/menus/application/)

## 輸入選項

### EnableFileDrop

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>全部

```go
EnableFileDrop: true,
```

<strong>用途：</strong>允許將檔案從作業系統拖放至視窗中。

啟用後：

- 可以將從檔案管理員拖曳的檔案放入應用程式
- `WindowFilesDropped`事件會在觸發時提供放入檔案的路徑
- 具有`data-file-drop-target`屬性的元素會提供詳細的放置資訊

**使用情境：**

- 檔案上傳介面
- 文件編輯器
- 媒體匯入工具
- 任何接受檔案的應用程式

**範例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**HTML 放置區域：**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**如需完整文件，請參閱[檔案拖放](/features/drag-and-drop/files/)。**

## 安全性選項

### ContentProtectionEnabled

<strong>型別：</strong>`bool` <strong>預設值：</strong>`false` <strong>平台：</strong>Windows（10+）、macOS

```go
ContentProtectionEnabled: true,
```

<strong>用途：</strong>防止擷取視窗內容的畫面。

**平台支援：**

- <strong>Windows：</strong>Windows 10組建19041以上版本（完整支援），較舊版本（部分支援）
- <strong>macOS：</strong>完整支援
- <strong>Linux：</strong>不支援

**使用情境：**

- 銀行應用程式
- 密碼管理器
- 醫療紀錄
- 機密文件

**重要注意事項：**

1. 無法防止以實體相機拍攝
2. 某些工具可能可以繞過保護
3. 這只是整體安全防護的一部分，不能作為唯一的保護措施
4. 開發人員工具視窗不會自動受到保護

**範例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### 權限

<strong>型別：</strong>`map[PermissionType]Permission` <strong>預設值：</strong>`nil`（平台預設處理方式）<strong>平台：</strong>Linux、Windows（macOS 交由 TCC 處理）

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

<strong>用途：</strong>以宣告方式控制如何處理視窗網頁內容提出的功能權限要求（相機、麥克風、地理位置、通知、剪貼簿讀取），無須編寫平台專用程式碼。

**PermissionType 值：**`PermissionMicrophone`、`PermissionCamera`、`PermissionGeolocation`、`PermissionNotifications`、`PermissionClipboardRead`

**權限值：**

- `PermissionDefault`（0）— 平台的原生處理方式：在 macOS／Windows 上顯示作業系統／WebView2 提示；在 Linux 上允許相機和麥克風，拒絕其他所有項目
- `PermissionAllow`（1）— 不顯示提示即授予權限（Linux：僅實作相機和麥克風；其他類型仍會遭到拒絕）
- `PermissionDeny`（2）— 不顯示提示即拒絕權限

<strong>重要事項 — Windows：</strong>此選項推出前，Wails 會在不提示的情況下授予所有 WebView2 功能權限。現在，只要在`Permissions`中設定任何項目，就會停用這種全面授權。未列出的功能會顯示 WebView2 的原生提示，而不會自動獲得允許。請明確列出應用程式所需的每項功能。

**如需完整指南、平台支援矩陣及範例，請參閱[權限](/features/windows/permissions/)。**

## 視窗生命週期事件

使用`OnWindowEvent`和`RegisterHook`處理視窗生命週期事件。這些方法可精細控制視窗的關閉和銷毀行為。

### 取消關閉視窗

若要防止視窗關閉（例如有尚未儲存的變更），請搭配`WindowClosing`事件使用`RegisterHook`，並呼叫`event.Cancel()`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**重點：**

- `RegisterHook`會在事件發生前攔截事件
- 呼叫`event.Cancel()`可防止視窗關閉
- 取消後，視窗會保持開啟

### 處理視窗關閉

若要在視窗關閉時執行清理，請搭配`WindowClosing`事件使用`OnWindowEvent`：

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**重點：**

- `OnWindowEvent`會處理即將發生的事件
- 清理作業會在視窗銷毀前執行
- 無法在此取消關閉（請使用`RegisterHook`）

### 單例視窗清理模式

對於單例視窗（確保只有一個執行個體），請使用`WindowClosing`清除參照：

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## 平台特定選項

### Mac 選項

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**標題列**（`MacTitleBar`）

- `AppearsTransparent`－使標題列透明，內容延伸至標題列區域
- `Hide`－完全隱藏標題列
- `HideTitle`－僅隱藏標題文字
- `FullSizeContent`－將內容延伸至整個視窗大小

在 macOS 上，WebView 透明效果和以程式方式開啟檢查器需要`private_mac_apis`建置標籤。若未使用該標籤，這些選項仍然有效，但僅限私有 API 的操作不會產生任何效果。Liquid Glass 群組設定也會被忽略，且樣式會改用公開 API 的替代方案。如需建置命令與確切行為，請參閱[私有 macOS API](/guides/build/private-macos-apis/)。

**背景效果**（`MacBackdrop`）

- `MacBackdropNormal`－標準不透明背景
- `MacBackdropTranslucent`－<strong>WebView 透明效果需要私有 API。</strong>若未使用該標籤，原生模糊效果會保留在不透明的 WebView 後方。
- `MacBackdropTransparent`－<strong>WebView 透明效果需要私有 API。</strong>若未使用該標籤，WebView 會保持不透明。
- `MacBackdropLiquidGlass`－<strong>WebView 透明效果需要私有 API。</strong>若未使用該標籤，玻璃圖層會保留在不透明的 WebView 後方；樣式則會改用公開 API 的替代方案。

**Liquid Glass**（`MacLiquidGlass`）

| 欄位或值 | 在 macOS 上對私有 API 的相依性 |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | 原生一般樣式是公開的；背景效果的 WebView 透明效果需要`private_mac_apis`。 |
| `Style: LiquidGlassStyleLight` | 該標籤會保留既有的原生透明樣式對應；若未使用，Wails 會採用搭配 Aqua 外觀的一般玻璃效果。 |
| `Style: LiquidGlassStyleDark` | <strong>私有 API：</strong>未記載於文件的原生樣式值`2`；若未使用該標籤，則採用搭配 Dark Aqua 外觀的一般玻璃效果。 |
| `Style: LiquidGlassStyleVibrant` | 原生透明樣式對應是公開的；背景效果的 WebView 透明效果需要該標籤。 |
| `GroupID` | <strong>私有 API：</strong>非空值會要求進行群組化；若未使用該標籤，則會被忽略。 |
| `GroupSpacing` | <strong>私有 API：</strong>正值會要求設定群組間距；若未使用該標籤，則會被忽略。 |
| `Material`, `CornerRadius`, `TintColor` | 這些值本身不依賴私有 API。 |

如需原生樣式值與作業系統支援情況，請參閱[Liquid Glass 值](/guides/build/private-macos-apis/#liquid-glass-values)。

**不可見標題列高度**（`int`）

- 不可見標題列區域的高度（用於拖曳）
- 僅在原生標題列拖曳區域隱藏時生效，也就是視窗無框線（`Frameless: true`）或使用透明標題列（`AppearsTransparent: true`）時
- 對具有可見標題列的標準視窗沒有作用

**視窗類別**（`MacWindowClass`）

- `MacWindowClassWindow`－標準`NSWindow`行為（預設）
- `MacWindowClassPanel`－永遠不會成為應用程式主視窗的輔助`NSPanel`

`PanelPreferences`僅適用於`MacWindowClassPanel`：

- `NonActivating`會加入`NSWindowStyleMaskNonactivatingPanel`。顯示面板或讓面板取得焦點不會啟用 Wails 應用程式，但面板仍可成為鍵盤輸入的主要視窗，以便操作控制項及輸入文字。
- `FloatingPanel`會啟用 AppKit 的浮動面板行為。
- `BecomesKeyOnlyIfNeeded`僅在所點選的檢視要求鍵盤輸入時取得主要視窗狀態。
- `UtilityWindow`會套用原生工具視窗樣式。

Wails 面板會在應用程式停用時保持可見，並在關閉時釋放，符合`WebviewWindow`預期的生命週期。這些設定是刻意覆寫`NSPanel`中相反的預設值。

視窗類別、層級、啟用原則和集合行為分別用於解決不同問題：

- `WindowClass`會選取`NSWindow`或`NSPanel`，並控制主視窗和鍵盤主要視窗的語意。
- `WindowLevel`控制 Z 軸順序。選單列疊加層請使用`MacWindowLevelPopUpMenu`。
- `MacOptions.ActivationPolicy`控制整個應用程式，包括 Dock 和選單列的呈現方式。非啟用型面板不需要附屬應用程式啟用原則，但應用程式仍可使用該原則隱藏 Dock 圖示。
- `CollectionBehavior`控制視窗參與 Spaces 和全螢幕的方式。

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**視窗層級**（`MacWindowLevel`）

- `MacWindowLevelNormal`－標準視窗層級（預設）
- `MacWindowLevelFloating`－浮動於一般視窗上方
- `MacWindowLevelTornOffMenu`－分離式選單層級
- `MacWindowLevelModalPanel`－強制回應面板層級
- `MacWindowLevelMainMenu`－主選單層級
- `MacWindowLevelStatus`－狀態視窗層級
- `MacWindowLevelPopUpMenu`－快顯選單層級
- `MacWindowLevelScreenSaver`－螢幕保護程式層級

明確設定的`WindowLevel`優先於`AlwaysOnTop`和`PanelPreferences.FloatingPanel`。若未明確設定層級，`AlwaysOnTop`或浮動面板會解析為`MacWindowLevelFloating`；否則層級為`MacWindowLevelNormal`。之後呼叫`SetAlwaysOnTop`仍屬於明確的執行階段變更。

**集合行為**（`MacWindowCollectionBehavior`）

控制視窗在 macOS Spaces 和全螢幕模式下的行為。這些是位元遮罩值，可使用位元 OR（`|`）組合。

**Space 行為：**

- `MacWindowCollectionBehaviorDefault`－使用 FullScreenPrimary（預設，向下相容）
- `MacWindowCollectionBehaviorCanJoinAllSpaces`－視窗會顯示在所有 Spaces 上
- `MacWindowCollectionBehaviorMoveToActiveSpace`－顯示時移至目前使用中的 Space
- `MacWindowCollectionBehaviorManaged`－預設的受管理視窗行為
- `MacWindowCollectionBehaviorTransient`－暫時性／短暫存在的視窗
- `MacWindowCollectionBehaviorStationary`－切換 Space 時保持原位

**視窗循環切換：**

- `MacWindowCollectionBehaviorParticipatesInCycle` - 納入 Cmd+` 循環切換
- `MacWindowCollectionBehaviorIgnoresCycle` - 排除在 Cmd+` 循環切換之外

**全螢幕行為：**

- `MacWindowCollectionBehaviorFullScreenPrimary` - 可進入全螢幕模式
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 可疊加於全螢幕應用程式之上
- `MacWindowCollectionBehaviorFullScreenNone` - 停用全螢幕功能
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - 允許並排顯示（macOS 10.11+）
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - 防止並排顯示（macOS 10.11+）

**範例——類似 Spotlight 的視窗：**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**範例——單一行為：**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode**（`MacWindowTabbingMode`）

控制 macOS 10.12及更新版本的視窗分頁行為。視窗分頁可將多個視窗群組為分頁。

**選項：**

- `MacWindowTabbingModeDefault` - 零值哨兵（未明確設定）。執行階段預設不允許分頁
- `MacWindowTabbingModeAutomatic` - 由系統決定分頁行為
- `MacWindowTabbingModePreferred` - 視窗偏好使用分頁模式
- `MacWindowTabbingModeDisallowed` - 停用視窗分頁

**範例——停用視窗分頁：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**範例——偏好使用視窗分頁：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences**（`MacWebviewPreferences`）

細部控制底層的`WKWebView`設定。所有欄位皆為選用；未設定的欄位不會變更 WebKit 預設值。

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — 設為`true`時，按下 Tab 會將焦點移至連結和表單控制項（預設值：`false`）
- `TextInteractionEnabled` — 設為`true`時，使用者可以選取 WebView 中的文字並與其互動（預設值：`true`）
- `FullscreenEnabled` — 設為`true`時，網頁內容可透過 HTML Fullscreen API 進入全螢幕（預設值：`false`）。需要 macOS 12.3或更新版本。
- `AllowsBackForwardNavigationGestures` — 設為`true`時，水平滑動手勢會觸發上一頁／下一頁導覽（預設值：`false`）
- `AllowsMagnification` — 設為`true`時，WebView 會啟用雙指撥動縮放（預設值：`false`）
- `AllowsAirPlayForMediaPlayback` — 設為`true`時，可將媒體串流至 AirPlay 裝置（預設值：`true`）
- `JavaScriptCanOpenWindowsAutomatically` — 設為`true`時，JavaScript 無須使用者手勢即可開啟新視窗（預設值：`false`）
- `MinimumFontSize` — 以點為單位的最小字型大小。使用`optional.NewVar(12.0)`進行設定。若未設定，則保留 WebKit 預設值。
- `ApplicationNameForUserAgent` — 覆寫 WebKit 使用者代理字串中的應用程式名稱後綴。當網站拒絕預設的`"wails.io"`識別碼時（例如 YouTube 嵌入內容），這項設定很有用。留空可保留預設值。
- `EnableAutoplayWithoutUserAction` — 設為`true`時，音訊和視訊無須使用者手勢即可自動播放。對應至`WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone`（預設值：`false`）

### Windows 選項（個別視窗）

個別視窗的結構為`application.WindowsWindow`，**不是**`WindowsOptions`（後者是*應用程式*層級的結構）。

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon**（`bool`）

- 從標題列移除圖示。

**DisableMenu**（`bool`）

- 停用視窗的選單列。設為`true`時，即使已設定選單列，視窗也不會顯示它。
- 預設值：`false`

**BackdropType**（`BackdropType`）

- `application.Auto` - 系統預設值
- `application.None` - 無背景效果
- `application.Mica` - Mica 材質（Windows 11）
- `application.Acrylic` - Acrylic 材質（Windows 11）
- `application.Tabbed` - 分頁材質（Windows 11）

沒有`WindowsBackdropTypeMica`樣式的常數，請使用`application.Mica`等值。

**CustomTheme**（`ThemeSettings`）

- 值（非指標）。為視窗邊框、標題列文字／背景及選單列自訂深色／淺色模式的色彩。

**DisableFramelessWindowDecorations**（`bool`）

- 停用預設的無框視窗裝飾（Aero 陰影、圓角）。

**NonClientRegionSupport**（`bool`）

- 為無框的自訂標題列啟用 WebView2 原生的`app-region: drag`／`app-region: no-drag`支援。
- 這僅適用於簡單的原生應用程式拖曳。它不提供原生的自訂標題列按鈕行為，也不會為自訂最大化按鈕提供 Windows 11 Snap Assist／Snap Layouts。

**WebView2CompositionHosting**（`bool`）

- 為自訂標題列按鈕啟用由 Wails 管理的`--wails-non-client-region`支援，並提供 Windows 原生行為，包括自訂最大化按鈕上的 Windows 11 Snap Assist／Snap Layouts。
- 實驗性功能。此功能會透過`ICoreWebView2CompositionController`和 DirectComposition 裝載 WebView2，而非使用預設的 HWND 裝載控制器。
- 當視窗同時需要 WebView2 原生的`app-region`支援，以及由 Wails 管理的自訂標題列按鈕區域時，可與`NonClientRegionSupport`搭配使用。

**範例：**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**範例——自訂 Windows 標題列區域：**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

如需詳細行為、取捨及配套 CSS，請參閱[無框 Windows 視窗](/features/windows/frameless/#native-non-client-regions-on-windows)。

### Linux 選項（個別視窗）

個別視窗的結構為`application.LinuxWindow`，**不是**`LinuxOptions`。

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon**（`[]byte`）

- 視窗圖示（PNG 格式）。

**WindowIsTranslucent**（`bool`）

- 需要合成器支援。

**範例：**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## 應用程式層級的 Windows 選項

某些 Windows 特定選項必須在應用程式層級設定，而不能針對個別視窗設定。這是因為 WebView2 會讓每個使用者資料路徑共用單一瀏覽器環境。

### 瀏覽器旗標

WebView2 瀏覽器旗標可控制應用程式中**所有視窗**的實驗性功能與行為。這些旗標必須在`application.Options.Windows`中設定：

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures**（`[]string`）

- 要啟用的 WebView2 功能旗標清單
- 如需可用的旗標，請參閱[WebView2 瀏覽器旗標](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- 範例：`"msWebView2EnableDraggableRegions"`

**DisabledFeatures**（`[]string`）

- 要停用的 WebView2 功能旗標清單
- Wails 會自動停用`msSmartScreenProtection`
- 範例：`"msExperimentalFeature"`

**AdditionalBrowserArgs**（`[]string`）

- 傳遞給瀏覽器處理程序的 Chromium 命令列引數
- 必須包含`--`前綴（例如`"--remote-debugging-port=9222"`）
- 如需可用的引數，請參閱[Chromium 命令列開關](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="重要"}
由於 WebView2 會針對每個使用者資料路徑共用單一瀏覽器環境，因此這些旗標會全域套用至所有視窗。不同視窗無法使用不同的瀏覽器旗標。

@end

**完整範例：**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## 完整範例

以下是可用於生產環境的視窗設定：

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

前端資源是在應用程式層級（`application.Options.Assets`）提供，而不是針對個別視窗提供。

## 後續步驟

- [視窗基礎](/features/windows/basics/)－建立及控制視窗
- [多視窗](/features/windows/multiple/)－多視窗模式
- [無框線視窗](/features/windows/frameless/)－自訂視窗外框
- [視窗事件](/features/windows/events/)－生命週期事件

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[範例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。
