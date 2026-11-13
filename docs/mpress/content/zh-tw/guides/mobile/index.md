---
title: "行動平台概覽"
description: "使用與桌面應用程式相同的 Go 程式碼庫建置 iOS 和 Android 應用程式"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3可在**iOS 和 Android**上執行，並沿用您為桌面平台撰寫的相同`main.go`和前端。不需要另外建立行動專案、程式碼共享橋接器，也不必重寫：Go 二進位檔會針對行動平台目標進行編譯，並由原生 WebView 轉譯您現有的前端。

@cards{cols="2"}
iOS
WKWebView + UIKit 宿主。資源透過自訂`wails://`配置提供，不會開啟任何連接埠。需要具備完整版 Xcode 的**macOS**。

[iOS 指南 →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`。Go 透過 NDK 編譯為`libwails.so`。可在 macOS、Linux 和 Windows 上使用。

[Android 指南 →](/guides/mobile/android/)

@end

## 實際執行：Kitchen Sink 範例

若要瞭解能實現哪些功能，最好的方式是查看**Kitchen Sink**：這是一個採用單一程式碼庫，且在 iOS、Android 和桌面平台上執行方式一致的 Wails 應用程式：

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="繫結 · 事件 · 對話方塊 · 觸覺回饋 · 地理位置 · 生物辨識 · 通知 · 安全儲存空間 · 以及更多功能，全都來自同一個 main.go"}
它透過7個分頁示範所有主要的行動 API 介面，而且也能在桌面平台上執行。透過前端的平台檢查，桌面平台會隱藏**Mobile**和**Hardware**分頁；為桌面平台建置時，Go 端不會為`common:*`行動事件註冊任何處理常式。若要使用單一程式碼庫發行至所有平台，建議採用此模式。

| 分頁 | 平台 | 示範內容 |
| --- | --- | --- |
| **繫結** | 全部 | JS → Go 服務呼叫，回傳值、結構和錯誤 |
| **事件** | 全部 | Go → JS 時鐘、JS → Go → JS ping/pong、作業系統事件（電池、網路、佈景主題） |
| **對話方塊** | 全部 | 各平台的原生訊息對話方塊 |
| **系統** | 全部 | 剪貼簿、螢幕度量、裝置資訊 |
| **行動平台** | iOS + Android | 分享面板、保持喚醒、手電筒、亮度、生物辨識、本機通知、安全儲存空間 |
| **硬體** | iOS + Android | 觸覺回饋、地理位置、加速度計、距離感測、文字轉語音 |
| **原生功能** | iOS + Android | iOS：觸覺回饋 + WKWebView 切換 · Android：震動 + Toast 訊息 |

若要自行執行：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 運作方式

所有平台都採用相同的應用程式模型：

1. **Go 後端** — 您的服務、事件處理常式和應用程式邏輯無須變更，即可針對`GOOS=ios`和`GOOS=android`進行編譯。
2. **前端** — 使用完全相同的 HTML/JS/CSS。`@wailsio/runtime`套件的運作方式相同；服務繫結、事件、對話方塊和剪貼簿全都透過相同的處理程序內傳輸進行路由。
3. **WebView 宿主** — 在 iOS 上是位於`UIViewController`內的`WKWebView`；在 Android 上則是位於`Activity`內的`WebView`。Wails 會自動設定訊息橋接器。
4. **處理程序內資源供應** — 資源直接從 Go 記憶體提供，而不是由 localhost 伺服器提供。不會開啟連接埠、不使用回送，也不會增加額外延遲。

平台特定的行為位於受`//go:build ios`或`//go:build android`保護的檔案中，讓共用程式碼保持簡潔。

## 先決條件一覽

| 需求 | iOS | Android |
| --- | --- | --- |
| 作業系統 | 僅限 macOS | macOS、Linux、Windows |
| 工具鏈 | 完整版 Xcode（不能只有 CLI 工具） | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 驗證方式 | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
設定工具鏈後執行`wails3 doctor`，它會精確顯示每個平台已找到和缺少的項目。

@end

## 支援的功能

兩個平台共用相同的核心功能集：

| 功能 | iOS | Android |
| --- | --- | --- |
| 服務繫結（JS → Go） | ✅ | ✅ |
| 事件（雙向） | ✅ | ✅ |
| 訊息對話方塊 | ✅ UIAlertController | ✅ AlertDialog |
| 開啟檔案對話方塊 | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| 儲存檔案對話方塊 | ❌ 改為寫入沙箱 | ❌ 改為寫入沙箱 |
| 剪貼簿 | ✅ UIPasteboard | ✅ ClipboardManager |
| 螢幕／安全區域度量 | ✅ | ✅ |
| 生命週期事件 | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| 觸覺回饋 | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| 裝置資訊 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| 原生分頁（iOS） | ✅ UITabBar | — |
| Toast 訊息（Android） | — | ✅ `Android.Toast.Show` |
| 多個視窗 | ❌ 僅限第一個視窗 | ❌ 僅限第一個視窗 |
| 視窗幾何配置／選單／系統匣 | 刻意不執行任何操作 | 刻意不執行任何操作 |

## 建置標籤規則

撰寫平台條件式程式碼時，請留意兩項重要規則：

- **`ios`隱含`darwin`**——標記為`//go:build darwin`的檔案也會針對 iOS 編譯。若只以 macOS 為目標，請使用`//go:build darwin && !ios`。
- **`android`隱含`linux`**——標記為`//go:build linux`的檔案也會針對 Android 編譯。若只以桌面版 Linux 為目標，請使用`//go:build linux && !android`。

在執行階段，`runtime.GOOS`會分別傳回`"ios"`和`"android"`。

## 執行階段平台偵測

建置標籤適用於只能在特定平台上*編譯*的程式碼。若要在共用程式碼中進行一般條件分支，請使用`application.System`——每個建置版本都可使用它（不需要建置標籤），因此同一個檔案可在所有平台上運作：

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

可用值：`IsMobile()`、`IsDesktop()`、`IsServer()`（`server`建置標籤）以及`IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`。

前端在`@wailsio/runtime`中提供了對應的輔助函式：

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 後續步驟

@cards{cols="2"}
🚀 您的第一個行動應用程式
只需幾分鐘，即可讓桌面版 Wails 應用程式在 iOS 模擬器或 Android 模擬器上執行。

[開始使用 →](/guides/mobile/first-mobile-app/)

---
iOS 指南
完整說明 iOS 工具鏈設定、模擬器、裝置建置、簽署、組態及 API 參考。

[iOS 指南 →](/guides/mobile/ios/)

---
Android 指南
完整說明 Android SDK/NDK 設定、模擬器、APK 簽署、Play Store 封裝及 API 參考。

[Android 指南 →](/guides/mobile/android/)

@end
