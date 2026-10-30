---
title: "Android"
description: "在 Android 上建置及執行 Wails 應用程式——工具鏈設定、模擬器、APK 簽署、Play Store 封裝及 API 參考"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="實驗性功能"}
Android 支援目前仍屬實驗性功能，未來版本可能會有所變更。

@end

@note{type="tip"}
第一次使用 Wails 開發行動應用程式嗎？請先閱讀[您的第一個行動應用程式 →](/guides/mobile/first-mobile-app/)，依照逐步指南操作，再回到這裡查閱完整參考資料。

@end

Wails v3應用程式會以原生應用程式的形式在 Android 上執行：Android `WebView`負責呈現前端；資產由 Go 資產伺服器支援的`WebViewAssetLoader`透過**程序內**方式提供（不使用 localhost 伺服器，也不開放任何連接埠）；標準`@wailsio/runtime`則可原封不動地運作——服務繫結、事件、對話方塊及剪貼簿都會經由 Go 訊息處理器路由。

同一個`main.go`可為桌面平台和 Android 建置。Go 程式碼會編譯為 C 共用程式庫（`libwails.so`、`GOOS=android`加上 NDK 工具鏈），再由小型 Java 主程式載入。Android 特有的行為位於各平台專用的 Go 檔案中，並由`//go:build android`保護。

## 需求

- **Android SDK**，並包含 platform-tools、SDK 平台（API 35）、build-tools 及**NDK**（26.3.x）——`wails3 doctor`會顯示它找到的元件
- 供 Gradle 使用的**JDK**（例如 OpenJDK 21）；若`java`不在您的`PATH`中，請設定`JAVA_HOME`
- Go 1.26以上版本及 npm
- 指向 SDK 的`ANDROID_HOME`（或`ANDROID_SDK_ROOT`）

使用命令列工具安裝 SDK 元件：

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## 在模擬器上執行

從專案目錄執行：

```bash
wails3 task android:run
```

如果沒有正在執行的模擬器，此操作會啟動一個模擬器，接著產生繫結、建置前端、將 Go 程式碼編譯成適用於模擬器 ABI 的`libwails.so`、使用 Gradle 組裝偵錯 APK，然後安裝並啟動它。

其他實用指令：

```bash
wails3 task android:logs    # stream the app's logcat output
```

在偵錯組建中，可透過 Chrome 的`chrome://inspect`檢查 WebView。

## 封裝

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

正式環境組建使用`-tags production,android`、會移除符號，並在編譯時排除框架的內部診斷功能。`wails3 task android:package:fat`會將`arm64-v8a`和`x86_64`一併建置至單一 APK 中。

Google Play 要求新提交的應用程式使用 Android App Bundle（`.aab`）格式，而且必須以 Android 15（API 35）或更高版本為目標；專案範本會在`build/android/app/build.gradle`中將`compileSdk`和`targetSdk`設為35。`wails3 task android:bundle:fat`會產生包含兩種 ABI 的`bin/<AppName>.aab`；Google Play 會據此產生針對各裝置最佳化的 APK，因此包含所有 ABI 的套件是上傳商店時應使用的成品。APK 仍是本機及模擬器測試最快的方式，因為`.aab`無法直接使用`adb`安裝。

`android:run`和`android:deploy-emulator`是以模擬器為目標的工作。若要使用實體 Android 裝置，請使用`android:run:device`建置偵錯 APK，或使用`android:deploy-device`建置發行 APK。兩者都會針對`arm64`建置，從`adb devices`中選取第一個已連線且非模擬器的項目，安裝 APK，然後啟動`com.wails.app.MainActivity`。若要指定特定裝置，請傳入`DEVICE_ID=<serial>`。

## 簽署與發行組建

若未提供金鑰庫，發行組建會使用 Android **偵錯**金鑰庫簽署，以便安裝測試。若要使用自己的金鑰庫簽署，請設定：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

相同的變數也可用於簽署 App Bundle：設定這些變數後執行`wails3 task android:bundle:fat`，即可產生可上傳至 Play 的`.aab`。若未設定，套件會使用偵錯金鑰庫簽署，而 Google Play 會拒絕該套件，因此此工作會顯示警告。

@note{type="tip"}
使用[Play 應用程式簽署](https://support.google.com/googleplay/android-developer/answer/9842756)時，您在本機用於簽署的金鑰庫就是**上傳金鑰**：Google 會使用它驗證您上傳的內容，再以其管理的應用程式簽署金鑰重新簽署應用程式。另請注意，Google Play 要求每次上傳都使用更高的`versionCode`；請在`build/android/app/build.gradle`中遞增該值。

@end

## 設定

前端在執行階段透過`Android`執行階段物件驅動 Android 功能：`Android.Haptics.Vibrate(durationMs)`、`Android.Device.Info()`、`Android.Toast.Show(message)`。套件名稱由建置工作中的`APP_ID`控制。

## 可用與不可用的功能

| 範圍 | 狀態 |
| --- | --- |
| WebView 與程序內資產（`WebViewAssetLoader`） | ✅ |
| 服務繫結、事件（雙向） | ✅ |
| 訊息對話方塊 | ✅ 具有按鈕回呼的 AlertDialog |
| 開啟檔案／多個檔案對話方塊 | ✅ Storage Access Framework（檔案會匯入為快取副本） |
| 開啟目錄／儲存檔案對話方塊 | ❌ 傳回錯誤——請改為寫入應用程式沙箱 |
| 剪貼簿 | ✅ ClipboardManager |
| 螢幕 API | ✅ WindowMetrics，包括扣除系統列後的工作區域 |
| 生命週期事件（`events.Android.*`） | ✅ |
| 觸覺回饋、裝置資訊、快顯通知 | ✅ `Android.*`執行階段 API |
| 模擬器與實體裝置組建 | ✅ `android:run`、`android:run:device`、`android:deploy-emulator`、`android:deploy-device` |
| 視窗幾何、選單、系統匣 | 刻意設計為無操作 |
| 多個視窗 | 僅顯示第一個視窗 |

## 移植注意事項

- 桌面程式碼可在`GOOS=android`下原封不動地編譯；由於 Android 應用程式採全螢幕顯示，幾何、選單及系統匣呼叫會變成無操作。
- `android` **隱含`linux`建置標籤**（Android 使用 Linux 核心）：僅供桌面 Linux 使用的檔案需要`//go:build linux && !android`，而在執行階段，`runtime.GOOS`為`"android"`。
- 請以寫入應用程式沙箱並搭配 intent 分享流程，取代儲存檔案及選擇目錄對話方塊。開啟檔案對話方塊可以正常運作，並會將所選文件匯入為快取目錄中的副本，因此您可以取得實際的檔案系統路徑。
- 實際的應用程式一律使用`CGO_ENABLED=1`和 NDK 建置；非 cgo 路徑僅供`wails3 generate bindings`等工具載入套件。
- 請採用響應式前端設計；`Screens`工作區域不包含狀態列和導覽列。
