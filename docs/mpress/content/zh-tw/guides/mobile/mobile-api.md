---
title: "行動裝置 API"
description: "跨平台的 application.Mobile 管理器——以建置條件保護的單一進入點，用於 iOS 與 Android 共用的原生行動裝置功能"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="實驗性功能"}
行動裝置支援仍屬實驗性功能，未來版本可能有所變更。

@end

原生行動裝置功能透過兩種方式公開：

- **各平台專用管理器**——`application.IOS`（位於`//go:build ios`檔案中）和`application.Android`（位於`//go:build android`檔案中）。凡是平台特有的功能，請使用這些管理器。各平台完整的 API 介面請參閱[iOS](/guides/mobile/ios/)和[Android](/guides/mobile/android/)參考資料。
- **`application.Mobile`**——以建置條件保護的單一管理器，涵蓋兩個平台上行為完全相同的功能子集。若您希望使用一套可在所有平台編譯及執行的程式碼路徑，請使用此管理器。

## `application.Mobile`

`application.Mobile`在 iOS 上會分派至`IOS`，在 Android 上會分派至`Android`，在桌面平台上則分派至不執行任何操作的虛設實作。由於它沒有建置條件限制，您可以從一般、不限定平台的 Go 程式碼呼叫它，無須自行建立`//go:build`檔案：

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()`會傳回應用程式私有檔案目錄的絕對路徑——在 Android 上為`getFilesDir()`，在 iOS 上為 Application Support 目錄——建議將資料庫及其他持久性檔案存放於此。在桌面平台上，或裝置上的目錄無法使用時（在 iOS 上則為無法建立時），它會傳回空字串，因此使用前請檢查是否為`""`。

@note{type="note"}
不在裝置上執行時（桌面版建置），每個`Mobile`方法都不會執行任何操作，而每個查詢都會傳回其零值（字串為`""`）。因此，跨平台程式碼可以無條件呼叫`application.Mobile.*`。如果在桌面平台上也需要實際路徑，請依平台分支處理，並改用`os.UserConfigDir()`或類似方案。

@end

## 功能

`Mobile`管理器會公開在 iOS 與 Android 上具有相同函式簽章的功能：

| 功能 | API | 備註 |
| --- | --- | --- |
| 分享面板 | `Mobile.Share(json)` | `{text, url}` |
| 在應用程式外部開啟 URL | `Mobile.OpenURL(url)` | 系統瀏覽器 |
| 讓螢幕保持喚醒 | `Mobile.SetKeepAwake(bool)` |  |
| 手電筒／閃光燈 | `Mobile.SetTorch(bool)` | → `common:torch` |
| 安全區域內距 | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| 應用程式資訊 | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| 畫面方向鎖定 | `Mobile.SetOrientation(mode)` | `portrait`／`landscape`／`auto` |
| 狀態列 | `Mobile.SetStatusBar(json)` | 樣式與顯示狀態 |
| 儲存空間資訊 | `Mobile.StorageJSON()` | `{free,total}` 位元組 |
| 儲存路徑 | `Mobile.StoragePath()` | 應用程式私有檔案目錄 |
| 電源／電池 | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| 網路狀態 | `Mobile.NetworkJSON()` | `{connected,type}` |
| 生物辨識 | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| 安全儲存空間 | `Mobile.SecureGet(key)`／`Mobile.SecureDelete(key)` | Keychain／`EncryptedSharedPreferences` |
| 地理位置 | `Mobile.GetLocation()` | 單次定位 → `common:location` |
| 觸覺回饋 | `Mobile.Haptic(type)` | 衝擊／通知／選取 |
| 加速度計 | `Mobile.SetMotion(bool)` | → `common:motion` |
| 接近感測 | `Mobile.SetProximity(bool)` | → `common:proximity` |
| 文字轉語音 | `Mobile.Speak(text)`／`Mobile.StopSpeak()` |  |
| 鍵盤內距 | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| 螢幕擷取 | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| 相機 | `Mobile.CapturePhoto()`／`Mobile.CaptureVideo()` | → `common:capture` |

非同步結果會以`common:*`事件的形式傳入，與各平台管理器完全相同；酬載請參閱[事件](/guides/mobile/ios/#events)。

## 仍因平台而異的功能

iOS 與 Android 之間形式不同的功能<strong>不</strong>位於`Mobile`上；請從具有建置標籤的檔案，透過`application.IOS`／`application.Android`呼叫這些功能：

| 項目 | iOS | Android |
| --- | --- | --- |
| 亮度（設定） | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| 亮度／螢幕方向（取得） | `IOS.GetBrightness()`／`IOS.GetOrientation()` | `Android.BrightnessJSON()`／`Android.OrientationJSON()` |
| 本機通知 | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| 安全儲存區（寫入） | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| 背景執行 | `IOS.BeginBackgroundTask`／`EndBackgroundTask` | `Android.StartForegroundService`／`StopForegroundService` |

@note{type="tip"}
`MobileManager`介面是`application.Mobile`背後的契約。由於兩個平台管理器都必須實作此介面，因此可確保上述任何方法在 iOS 與 Android 上都維持完全相同的簽章；如果兩者出現差異，平台建置將無法編譯。

@end
