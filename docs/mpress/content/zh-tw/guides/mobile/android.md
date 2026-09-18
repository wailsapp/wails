---
title: "Android"
description: "在 Android 上建置並執行 Wails 應用程式——工具鏈設定、模擬器、APK 簽署、Play 商店封裝與 API 參考資料"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="實驗性功能"}
Android 支援目前仍屬實驗性質，未來版本可能會有所變更。

@end

@note{type="tip"}
剛開始使用 Wails 開發行動應用程式嗎？請先閱讀[您的第一個行動應用程式 →](/guides/mobile/first-mobile-app/)，依照逐步指南操作，再回到這裡查閱完整參考資料。

@end

Wails v3 應用程式會以原生應用程式的形式在 Android 上執行：Android `WebView` 會呈現前端；資產由 Go 資產伺服器支援的 `WebViewAssetLoader`在<strong>程序內</strong>提供（不使用 localhost 伺服器，也不開放任何 連接埠）；標準`@wailsio/runtime`則可直接沿用——服務 繫結、事件、對話方塊和剪貼簿都會透過 Go 訊息 處理器路由。

同一個`main.go`可針對桌面和 Android 進行建置。Go 程式碼會編譯成 C 共用程式庫（`libwails.so`、`GOOS=android`加上 NDK 工具鏈），並由 小型 Java 主控程式載入。Android 特有的行為位於各平台專用的 Go 檔案中，並由`//go:build android`保護。

## 需求

- **Android SDK**，其中須包含 platform-tools、SDK 平台（API 35）、build-tools，以及<strong>NDK</strong>（26.3.x）——`wails3 doctor`會顯示找到的項目
- 供 Gradle 使用的<strong>JDK</strong>（例如 OpenJDK 21）；如果`java`不在您的`PATH`中，請設定`JAVA_HOME`
- Go 1.25以上版本和 npm
- 指向 SDK 的`ANDROID_HOME`（或`ANDROID_SDK_ROOT`）

使用命令列工具安裝各個 SDK 元件：

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## 在模擬器上執行

在專案目錄中執行：

```bash
wails3 task android:run
```

此操作會在沒有模擬器執行時啟動模擬器，接著產生繫結、建置 前端、將你的 Go 程式碼編譯成適用於模擬器 ABI 的`libwails.so`、 使用 Gradle 組裝偵錯 APK，然後安裝並啟動它。

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

正式環境組建使用`-tags production,android`、會移除符號，並在編譯時 排除框架的內部診斷功能。`wails3 task android:package:fat` 會將`arm64-v8a`和`x86_64`都建置到單一 APK 中。

Google Play 要求新提交的應用程式採用 Android App Bundle（`.aab`）格式， 而且新提交項目必須以 Android 15（API 35）以上版本為目標； 專案範本會在`build/android/app/build.gradle`中，將`compileSdk`和`targetSdk`設為35。 `wails3 task android:bundle:fat`會產生包含兩種 ABI 的 `bin/<AppName>.aab`；Google Play 會從中產生針對各裝置最佳化的 APK， 因此包含所有 ABI 的套件是上傳至商店時應使用的成品。APK 仍是進行本機和模擬器 測試最快的方式，因為`.aab`無法直接使用`adb`安裝。

`android:run`和`android:deploy-emulator`是以模擬器為主的工作。若要在 實體 Android 裝置上執行，偵錯 APK 請使用`android:run:device`， 發行版 APK 則使用`android:deploy-device`。兩者都會針對`arm64`進行建置、從 `adb devices`選取第一個已連線的非模擬器項目、完成安裝，並啟動 `com.wails.app.MainActivity`。傳入`DEVICE_ID=<serial>`可指定特定 裝置。

## 簽署與發行版組建

如果沒有金鑰庫，發行版組建會使用 Android **偵錯** 金鑰庫簽署，以便安裝測試。若要使用自己的金鑰庫簽署，請設定：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

相同的變數也用於簽署 App Bundle：設定這些變數後執行`wails3 task android:bundle:fat`， 即可產生可提交至 Play 的`.aab`。如果未設定，套件會使用 偵錯金鑰庫簽署，而 Google Play 會拒絕該套件，因此此工作會 顯示警告。

@note{type="tip"}
使用[Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756)時， 您在本機簽署所用的金鑰庫就是<strong>上傳金鑰</strong>：Google 會用它驗證 您的上傳內容，再以其管理的應用程式簽署金鑰重新簽署應用程式。另請注意，Google Play 要求 每次上傳都使用更高的`versionCode`；請在`build/android/app/build.gradle`中遞增該值。

@end

## 設定

前端會在執行階段透過`Android`執行階段物件驅動 Android 功能： `Android.Haptics.Vibrate(durationMs)`、`Android.Device.Info()`、 `Android.Toast.Show(message)`。套件名稱由建置工作中的`APP_ID`控制。

## 支援與不支援的功能

| 領域 | 狀態 |
| --- | --- |
| WebView 與程序內資產（`WebViewAssetLoader`） | ✅ |
| 服務繫結、事件（雙向） | ✅ |
| 訊息對話方塊 | ✅ AlertDialog，支援按鈕回呼 |
| 開啟單一／多個檔案對話方塊 | ✅ Storage Access Framework（檔案會以快取複本匯入） |
| 開啟目錄／儲存檔案對話方塊 | ❌ 傳回錯誤——請改為寫入應用程式沙箱 |
| 剪貼簿 | ✅ ClipboardManager |
| 螢幕 API | ✅ WindowMetrics，包括扣除系統列後的工作區域 |
| 生命週期事件（`events.Android.*`） | ✅ |
| 觸覺回饋、裝置資訊、Toast 通知 | ✅ `Android.*`執行階段 API |
| 模擬器與實體裝置組建 | ✅ `android:run`、`android:run:device`、`android:deploy-emulator`、`android:deploy-device` |
| 視窗幾何配置、選單、系統匣 | 刻意不執行任何操作 |
| 多個視窗 | 僅顯示第一個視窗 |

## 移植注意事項

- 桌面程式碼可在`GOOS=android`下直接編譯，不需修改；由於 Android 應用程式會以全螢幕顯示，幾何配置、選單和系統匣呼叫將不執行任何操作。
- `android` **隱含`linux`建置標記**（Android 使用 Linux 核心）：僅適用於桌面 Linux 的檔案需要`//go:build linux && !android`，而在執行階段，`runtime.GOOS`為`"android"`。
- 請將儲存檔案和選擇目錄對話方塊，改為寫入應用程式沙箱並搭配 Intent 分享流程。開啟檔案對話方塊可以使用，且會將所選文件以複本匯入快取目錄，因此您會取得實際的檔案系統路徑。
- 實際的應用程式一律使用`CGO_ENABLED=1`和 NDK 建置；非 cgo 路徑的存在，僅是為了讓`wails3 generate bindings`之類的工具能夠載入套件。
- 前端應採用響應式設計；`Screens`工作區不包含狀態列和導覽列。
