---
title: "私有 macOS API"
description: "所有依賴私有 macOS API 的 Wails 功能與選項，包括選擇啟用的命令及公開建置的替代行為。"
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails 預設使用公開的 macOS API。只要使用 Go 建置標記`private_mac_apis`，即可啟用本頁列出的私有 WebKit 與 AppKit 呼叫。無論採用哪種建置，所有公開的 Go 選項與方法都仍可使用。未使用此標記時，僅限私有 API 的操作不會執行任何動作；具有公開替代方案的功能則會使用該替代方案。

@note{type="caution" title="選擇啟用 macOS 私有行為"}
設定視窗選項並不會啟用私有 API。請將`private_mac_apis`加入建置命令以啟用這些 API。此標記僅適用於 macOS 桌面建置，不適用於 iOS、Android、Windows、Linux 或伺服器建置。

@end

## 啟用私有 API

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

若要直接建置正式版本，請使用`go build -tags production,private_mac_apis .`。若是前端範例，請依照其 README 先建置繫結與資產，再執行範例。自訂或較舊的 Taskfile 必須將`EXTRA_TAGS`傳遞給 Go 編譯器。

## 功能清單

| 功能或值 | `private_mac_apis`所啟用的功能 | 未使用此標記時 |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | 原生視窗上方的透明 WKWebView | 原生視窗已設定，但 WebView 仍不透明 |
| `Mac.Backdrop: MacBackdropTranslucent` | 透明的 WKWebView，讓原生模糊效果透出 | 模糊效果已設定在不透明的 WebView 後方 |
| `Mac.Backdrop: MacBackdropLiquidGlass` | 玻璃圖層上方的透明 WKWebView，並使用私有 WebView 背景控制 | 玻璃圖層已設定在不透明的 WebView 後方；樣式使用公開的替代方案 |
| 設定 Liquid Glass 時清除 WebView 背景 | 私有 WebKit `backgroundColor`控制 | 在 macOS 12+ 上使用公開的`underPageBackgroundColor`，較舊的 macOS 則使用圖層顏色；這不會讓 WebView 變透明 |
| `app.Window.NewNotchWindow(...)` | 異形瀏海面板內的透明 WebView | 面板仍可運作，包括定位與動畫，但其 WebView 仍不透明 |
| `Mac.LiquidGlass.Style` | Wails 現有的原生樣式對應，包括一個未記載的深色樣式值 | 使用公開的標準／透明樣式與淺色／深色外觀；請參閱下方的值對照表 |
| `Mac.LiquidGlass.GroupID` | 針對非空識別碼要求使用私有玻璃群組功能 | 忽略；不要求建立群組 |
| `Mac.LiquidGlass.GroupSpacing` | 針對大於零的值要求使用私有群組間距 | 忽略 |
| `window.OpenDevTools()`和 JavaScript `Window.OpenDevTools()` | 在 macOS 12+ 上以程式方式開啟 WebKit 檢查器 | 不執行任何動作 |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | 在 macOS 12+ 上，要求於視窗首次顯示時以程式方式開啟檢查器 | 不執行任何動作 |
| macOS 13.3之前版本的舊式檢查器啟用方式 | 若建置時已納入檢查器支援，則啟用 WebKit 開發人員額外功能 | 不執行任何動作；透過公開 API 使用 Safari 檢查功能需要 macOS 13.3+ |

## WebView 透明度與背景

**需要私有 API：**`MacBackdropTransparent`、`MacBackdropTranslucent`、`MacBackdropLiquidGlass`和瀏海視窗所使用的 WebView 透明效果。在內部，Wails 會設定 WebKit 的私有`drawsBackground`鍵。僅將 HTML 或 CSS 背景設為透明，無法讓不透明的原生 WKWebView 變透明。

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

私有的 WebView 背景顏色操作會使用 WebKit 的`backgroundColor`鍵。Liquid Glass 設定會使用此操作清除 WebView 背景。未使用此標記時，這項內部操作會在 macOS 12+ 上使用公開的`underPageBackgroundColor`，較舊的 macOS 則使用檢視的圖層。這些替代方案不會讓 WebView 變透明。

在 macOS 上，`WebviewWindowOptions.BackgroundColour`和`window.SetBackgroundColour()`會設定<strong>原生視窗</strong>的顏色，本身不需要私有 API。同樣地，`Frameless`和`Mac.TitleBar.AppearsTransparent`使用公開的 AppKit API；需要私有 API 的是 WebView 透明度，而不是標題列透明度。若要在 macOS 上使用背景效果，請設定`Mac.Backdrop`，不要只依賴`BackgroundType`。

請參閱[視窗選項](/features/windows/options/#mac-options)、[無邊框視窗](/features/windows/frameless/#with-transparent-background)和[瀏海視窗](/features/windows/notch-windows/)。

## Liquid Glass 值

若原生`NSGlassEffectView`可用（macOS 26+），則採用下列對應關係。只有原生樣式值`0`（標準）和`1`（透明）有文件記載。兩種建置中的 Go 常數都會保留其現有值。

| `MacLiquidGlassStyle`值 | 使用`private_mac_apis`時 | 未使用此標記時 |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic`（`0`） | 原生標準樣式（`0`） | 原生標準樣式（`0`） |
| `LiquidGlassStyleLight`（`1`） | 現有的原生樣式對應（`1`，透明） | 採用 Aqua 外觀的原生標準樣式（`0`） |
| `LiquidGlassStyleDark`（`2`） | **未記載於文件的原生樣式值`2`** | 採用深色 Aqua 外觀的原生一般樣式（`0`） |
| `LiquidGlassStyleVibrant`（`3`） | 對應至現有的淺色／原生透明樣式（`1`） | 原生透明樣式（`1`） |

Automatic 和 Vibrant 使用文件中記載的原生樣式值，但若要讓整個視窗以 Liquid Glass 作為背景，仍須使用標籤才能啟用<strong>網頁檢視透明效果</strong>。Light 在這兩種建置中的外觀不同。私有樣式對應不保證未來的 macOS 版本會呈現相同效果。

**一律屬於私有 API：**`GroupID`和`GroupSpacing`。Wails 在要求進行群組化之前，會先檢查私有的`setGroupIdentifier:`、`setGroupName:`和`setGroupSpacing:`選擇器。若未使用該標籤，這些操作不會執行任何動作。啟用該標籤不保證目前執行的 macOS 版本支援這些選擇器。

`MacLiquidGlass.Material`、`CornerRadius`和`TintColor`本身不需要私有 API。在不具備原生 Liquid Glass 支援的 macOS 版本上，Wails 會設定半透明的備援效果；若要透過網頁檢視看到該效果，仍須使用此標籤。

## Web 檢查器

<strong>需要私有 API：</strong>呼叫`OpenDevTools()`或設定`OpenInspectorOnStartup: true`，以便從應用程式開啟 WebKit 的檢查器。Wails 使用私有的`_inspector`選擇器。若未使用`private_mac_apis`，這些操作會直接略過且不顯示任何訊息。

建置時也必須編譯檢查器支援。現有的`production`和`devtools`標籤仍保有原本的含義：

| 建置標籤 | 以程式開啟檢查器 | macOS 13.3+ 上的公開 Safari 檢查功能 |
| --- | --- | --- |
| 無 | 不執行任何動作 | 已啟用 |
| `private_mac_apis` | 在 macOS 12+ 上啟用 | 已啟用 |
| `production` | 不執行任何動作 | 已停用 |
| `production,private_mac_apis` | 不執行任何動作 | 已停用 |
| `production,devtools` | 不執行任何動作 | 已啟用 |
| `production,devtools,private_mac_apis` | 在 macOS 12+ 上啟用 | 已啟用 |

在 macOS 13.3+ 上，Wails 會使用公開的`WKWebView.inspectable`啟用 Safari 檢查功能；這不需要私有 API。在 macOS 13.3之前的版本中，備援機制會使用私有的`developerExtrasEnabled`偏好設定，因此除了必須在編譯時納入檢查器支援外，也需要`private_mac_apis`。

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## 維護此清單

所有私有原生呼叫都集中在`v3/pkg/application/mac_private_api_darwin.go`中；預設建置會選用`mac_public_api_darwin.go`。此功能清單涵蓋透明效果、網頁檢視背景顏色、玻璃樣式、玻璃群組化、開啟檢查器，以及啟用舊版檢查器。變更這些實作時，應同時更新本頁及受影響的選項文件。
