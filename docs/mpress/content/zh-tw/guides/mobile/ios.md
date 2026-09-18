---
title: "iOS"
description: "在 iOS 上建置並執行 Wails 應用程式，包括設定、模擬器、裝置建置、組態與原生功能"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="實驗性功能"}
iOS 支援仍屬實驗性質，未來版本可能會有所變更。

@end

@note{type="tip"}
第一次使用 Wails 開發行動應用程式嗎？請先參閱[您的第一個行動應用程式 →](/guides/mobile/first-mobile-app/)，按照逐步指南操作，之後再回到這裡查閱完整參考資料。

@end

Wails v3 應用程式會以完全原生的應用程式形式在 iOS 上執行，而且最棒的是，其運作方式與桌面版<em>完全相同</em>。相同的 Go 後端、相同的前端，以及相同的`@wailsio/runtime`：服務繫結、事件、對話方塊和剪貼簿的行為全都一致，<strong>完全不需要</strong>針對行動平台重新接線。不需要獨立的行動版程式碼庫、移植層，也沒有特殊 API 要學習，現有的 Wails 應用程式直接就能在 iOS 上執行。移植過程真正無縫：將應用程式原封不動地移過來，即可發布。

同一個`main.go`可同時為桌面和 iOS 建置；iOS 專屬調整則透過`application.Options.IOS`設定。

## 需求

- 已安裝<strong>完整 Xcode</strong> 的 macOS（只有命令列工具並不足夠）— `wails3 doctor`會顯示其可找到的 iOS SDK
- Go 1.25+ 和 npm

## 模擬器

從專案目錄執行：

```bash
wails3 task ios:run
```

這會建置應用程式；如果尚無模擬器正在執行，則會啟動一個模擬器，然後啟動應用程式。

實用的搭配工具：

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

在偵錯建置中，可從 Safari 的「開發」選單檢查 WebView。

## 封裝

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

這些是經過最佳化並移除偵錯資訊的正式環境建置。

## 裝置建置

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device`會為實體裝置建置。權利設定來自`build/ios/entitlements.plist`，且僅套用至裝置建置；請加入應用程式所需的能力鍵。

@note{type="tip"}
若要自動管理簽署、佈建及 App Store 封存，請使用`wails3 task ios:xcode`開啟產生的 Xcode 專案，改由 Xcode 進行建置。

@end

## 組態

`build/config.yml`：

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

啟動選項（`application.Options.IOS`）包括`DisableScroll`、`DisableBounce`、`DisableScrollIndicators`、`DisableInputAccessoryView`、`EnableBackForwardNavigationGestures`、`DisableLinkPreview`、`EnableInlineMediaPlayback`、`EnableAutoplayWithoutUserAction`、`DisableInspectable`、`UserAgent`、`ApplicationNameForUserAgent`、`BackgroundColour`，以及透過`EnableNativeTabs` + `NativeTabsItems`設定的原生底部分頁。

## 原生功能

iOS 專屬功能可透過`application.IOS`使用；請在`//go:build ios`檔案內從 Go 呼叫，讓共用程式碼不依賴特定平台。Android 透過`application.Android`提供相同的一組功能。

單次動作會立即傳回：

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

查詢輔助函式會以 JSON 傳回結果，包括`SafeAreaJSON()`、`AppInfoJSON()`、`PowerJSON()`、`NetworkJSON()`、`StorageJSON()`、`GetOrientation()`、`GetBrightness()`。`StoragePath()`會傳回應用程式 Application Support 目錄的絕對路徑，適合用來存放資料庫及其他持久性檔案（此函式在 iOS 上提供相當於 Android 的`getFilesDir()`的功能）。首次存取時會建立此目錄；如果無法建立，`StoragePath()`會傳回空字串，因此使用前請先檢查`""`。

### 事件

任何稍後才完成的作業（例如權限提示、感測器串流或相機擷取）都會以<strong>事件</strong>而非傳回值的形式提供結果；您可以在 Go 或前端監聽這些事件。與 Android 共用的功能，其名稱以`common:`為前綴；僅限 iOS 的功能則以`ios:`為前綴。

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| 事件 | 觸發來源 | 承載資料 |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

`v3/examples/mobile`下的綜合範例會以端對端方式串接上述每項功能。

## WebView 控制項

也可以在執行階段從 Go 變更部分 WebView 行為：

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

隨附的`@wailsio/runtime`也公開了一個小型的<strong>前端</strong> iOS 命名空間：

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

原生底部分頁選取會以`nativeTabSelected`事件傳送至`window`。

## 支援狀態

| 項目 | 狀態 |
| --- | --- |
| 前端算繪與資源 | ✅ |
| 服務繫結、事件（雙向） | ✅ |
| 訊息對話方塊 | ✅ |
| 開啟檔案／多個檔案／目錄對話方塊 | ✅ 以沙箱副本形式匯入 |
| 儲存檔案對話方塊 | ❌ 請改為寫入應用程式沙箱內 |
| 剪貼簿 | ✅ |
| 螢幕 API | ✅ 包含安全區域內的工作區域 |
| 生命週期事件 | ✅ |
| 視窗幾何資訊、選單、系統匣 | 在 iOS 上不執行任何操作 |
| 多個視窗 | 僅顯示第一個視窗 |

## 移植注意事項

- 桌面版程式碼無須變更即可針對 iOS 建置——視窗、選單和系統匣呼叫只是不執行任何操作。
- 將儲存檔案對話方塊改為寫入應用程式沙箱，然後分享檔案。
- 請採用回應式設計建構前端；安全區域會自動處理。
