---
title: "iOS"
description: "在 iOS 上建置並執行 Wails 應用程式——設定、模擬器、裝置建置、組態及原生功能"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="實驗性功能"}
iOS 支援目前仍屬實驗性質，未來版本可能有所變更。

@end

@note{type="tip"}
剛開始使用 Wails 開發行動應用程式嗎？請先參閱[您的第一個行動應用程式 →](/guides/mobile/first-mobile-app/)，依照逐步指南操作，再回到此處查看完整參考資料。

@end

Wails v3應用程式會在 iOS 上以完全原生的應用程式執行——最棒的是，其運作方式與桌面版*完全*相同。相同的 Go 後端、相同的前端、相同的`@wailsio/runtime`：服務繫結、事件、對話方塊和剪貼簿的行為全都一致，**完全不需要**針對行動平台重新接線。您不需要獨立的行動版程式碼庫、移植層，也不必學習特殊 API——現有的 Wails 應用程式可直接在 iOS 上執行。移植過程真正無縫：將應用程式原封不動地移過來即可發布。

相同的`main.go`可同時為桌面平台和 iOS 建置；iOS 專屬調整則透過`application.Options.IOS`設定。

## 需求

- macOS，並已安裝**完整的 Xcode**（僅有命令列工具並不足夠）——`wails3 doctor`會顯示可找到的 iOS SDK
- Go 1.26+ 和 npm

## 模擬器

從專案目錄執行：

```bash
wails3 task ios:run
```

這會建置您的應用程式；如果尚未執行模擬器，則會啟動一個，然後啟動應用程式。

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

這些是經過最佳化並移除偵錯符號的正式環境建置。

## 裝置建置

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device`會針對實體裝置進行建置。權限項目來自`build/ios/entitlements.plist`，且僅套用於裝置建置——請加入應用程式所需的功能權限鍵。

@note{type="tip"}
若要自動管理簽署、佈建及 App Store 封存，請使用`wails3 task ios:xcode`開啟產生的 Xcode 專案，改由 Xcode 執行建置。

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

啟動選項（`application.Options.IOS`）包括`DisableScroll`、`DisableBounce`、`DisableScrollIndicators`、`DisableInputAccessoryView`、`EnableBackForwardNavigationGestures`、`DisableLinkPreview`、`EnableInlineMediaPlayback`、`EnableAutoplayWithoutUserAction`、`DisableInspectable`、`UserAgent`、`ApplicationNameForUserAgent`、`BackgroundColour`，以及透過`EnableNativeTabs`與`NativeTabsItems`設定的原生底部分頁。

## 原生功能

iOS 專屬功能可透過`application.IOS`使用。請從`//go:build ios`檔案中的 Go 程式碼呼叫，以便讓共用程式碼維持平台無關。Android 透過`application.Android`提供相同的功能集。

單次動作會立即傳回：

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

查詢輔助函式會以 JSON 傳回結果，包括`SafeAreaJSON()`、`AppInfoJSON()`、`PowerJSON()`、`NetworkJSON()`、`StorageJSON()`、`GetOrientation()`、`GetBrightness()`。`StoragePath()`會傳回應用程式 Application Support 目錄的絕對路徑，適合用來存放資料庫及其他持久性檔案（相當於 iOS 上的 Android `getFilesDir()`）。首次存取時會建立此目錄；若無法建立，`StoragePath()`會傳回空字串，因此使用前請檢查是否為`""`。

### 事件

任何稍後才會完成的操作——例如權限提示、感測器串流或相機拍攝——都會以**事件**而非傳回值來傳遞結果；您可以在 Go 或前端中監聽這些事件。與 Android 共用的功能名稱以`common:`為前綴，僅限 iOS 的功能則以`ios:`為前綴。

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

`v3/examples/mobile`下的綜合範例會從頭到尾串接上述所有功能。

## WebView 控制項

您也可以在執行階段從 Go 變更部分 WebView 行為：

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

隨附的`@wailsio/runtime`也公開一個小型的**前端** iOS 命名空間：

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

原生底部分頁的選取結果會以`nativeTabSelected`事件送達`window`。

## 支援狀態

| 領域 | 狀態 |
| --- | --- |
| 前端轉譯與資產 | ✅ |
| 服務繫結、事件（雙向） | ✅ |
| 訊息對話方塊 | ✅ |
| 開啟檔案／多個檔案／目錄對話方塊 | ✅ 以沙箱副本匯入 |
| 儲存檔案對話方塊 | ❌ 請改為寫入應用程式沙箱內 |
| 剪貼簿 | ✅ |
| 螢幕 API | ✅ 包含安全區域的工作區 |
| 生命週期事件 | ✅ |
| 視窗幾何、選單、系統匣 | 在 iOS 上不執行任何操作 |
| 多個視窗 | 只會顯示第一個視窗 |

## 移植注意事項

- 桌面版程式碼無須修改即可為 iOS 建置——視窗、選單和系統匣呼叫只是不會執行任何動作。
- 將儲存檔案對話方塊改為先寫入應用程式的沙盒，再進行分享。
- 以前端響應式設計為原則；安全區域會自動為你處理。
