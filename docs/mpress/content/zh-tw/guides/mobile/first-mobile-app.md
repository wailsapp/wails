---
title: "您的第一個行動應用程式"
description: "幾分鐘內即可在 iOS 模擬器或 Android 模擬器上執行您的 Wails 應用程式"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

本指南將標準 Wails 桌面應用程式置於 iOS 模擬器或 Android 模擬器上執行。<strong>您不需要變更 Go 程式碼。</strong>相同的`main.go`可針對所有目標建置。

<strong>完成時間：</strong>15–30 分鐘（大部分時間用於首次執行時安裝工具鏈）

## 從桌面專案開始

如果還沒有專案，請建立一個全新的專案：

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

請先確認桌面應用程式可正常運作：

```bash
wails3 dev
```

應用程式開啟後，請將它關閉並繼續下一步。所有能在桌面環境執行的內容也能在行動裝置上執行——在本指南中，您不需要修改`main.go`或任何 Go 程式碼。

---

## 選擇平台

@tabs{sync-key="mobile-platform"}
[iOS 模擬器]
### 需求

- **macOS**（只能在 macOS 上建置 iOS 應用程式）
- **完整的 Xcode**——不能只有命令列工具。請從 App Store 安裝，然後執行：
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- <strong>Go 1.25+</strong>和<strong>npm</strong>（如果已執行`wails3 init`，便已安裝）

執行`wails3 doctor`進行驗證；它會列出可找到的 iOS SDK。

### 在模擬器上執行

@steps
### 啟動應用程式
```bash
wails3 task ios:run
```

就是這麼簡單——此操作會建置您的應用程式；如果尚無模擬器在執行，則會啟動一個模擬器，然後啟動應用程式。

@note{type="tip"}
首次執行需要幾分鐘（系統會針對 iOS 編譯並快取 Wails 框架）。之後每次執行都會快得多。

@end

啟動後，未經修改的桌面應用程式便會在 iOS 模擬器上執行——使用相同的`main.go`和相同的前端：

![在 iOS 模擬器上執行的預設 Wails 應用程式](/assets/ios-simulator-first-app.png)

### 串流顯示記錄
在另一個終端機中執行：

```bash
wails3 task ios:logs:dev
```

此操作會持續顯示模擬器記錄，並篩選出您的應用程式。`fmt.Println`和`log.Println`的輸出會顯示在這裡。

### 檢查 WebView
在 Safari 中選擇：**開發 → 模擬器 → 您的應用程式**。您可以使用完整的 Web 檢閱器，包括主控台、偵錯工具、網路面板等所有功能。

### 進行變更
編輯任一前端檔案（`frontend/src/main.js`、`index.html`等），然後重新執行`wails3 task ios:run`。Wails 會重新建置前端並重新啟動應用程式。

如果變更 Go 程式碼，也請重新執行`wails3 task ios:run`。Go 採用增量重新編譯，因此只會重新建置有變更的套件。

@end

### 在 Xcode 中開啟（選用）

```bash
wails3 task ios:xcode
```

此操作會在 Xcode 中開啟`build/ios/`。您可以使用 Xcode 部署至裝置、進行進階效能分析或管理佈建描述檔。Wails 會在每次建置時重新產生 Xcode 專案，因此請勿直接修改產生的檔案。

[Android 模擬器]
### 需求

您需要<strong>Android SDK</strong>、<strong>NDK</strong>和<strong>JDK</strong>。最簡單的方法是使用 Android Studio，也可以使用命令列工具：

@steps
### 安裝 Android 命令列工具
從[developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only)下載，並解壓縮至`~/android-sdk/cmdline-tools/latest/`。

### 安裝 SDK 元件
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### 建立模擬器
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### 設定環境變數
新增至`~/.zshrc`或`~/.bashrc`：

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

重新載入：`source ~/.zshrc`

### 安裝 JDK
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

執行`wails3 doctor`，確認所有項目都能找到。

### 在模擬器上執行

@steps
### 啟動應用程式
```bash
wails3 task android:run
```

首次執行時，此操作會：

- 如果沒有模擬器在執行，則啟動模擬器
- 產生繫結並建置前端
- 透過 NDK 交叉編譯器將 Go 程式碼編譯為`libwails.so`
- 使用 Gradle 組建偵錯 APK
- 在模擬器上安裝並啟動應用程式

@note{type="tip"}
首次建置會下載 Gradle 並編譯 NDK 工具鏈，預計需要5–10 分鐘。後續建置採增量方式，所需時間不到一分鐘。

@end

### 串流顯示記錄
在另一個終端機中執行：

```bash
wails3 task android:logs
```

此操作會執行`adb logcat`，並篩選出您的應用程式。`fmt.Println`的輸出會顯示在這裡。

### 檢查 WebView
開啟 Chrome 並前往`chrome://inspect`。您的應用程式 WebView 會顯示在<strong>遠端目標</strong>下方；按一下<strong>檢查</strong>即可開啟 DevTools。

### 進行變更
編輯任一檔案，然後重新執行`wails3 task android:run`。Gradle 的增量建置只會重新編譯有變更的程式碼。

@end

@end

---

## 瞭解執行過程

您的`main.go`完全沒有變更。所有工作都由 Wails 處理：

- **建置系統** — 專案中的`Taskfile.yml`包含`ios:*`和`android:*`工作，用來驅動平台專屬工具鏈。
- **Go 交叉編譯** — 使用`GOOS=ios`或`GOOS=android`，並搭配適當的`GOARCH`和 sysroot。
- **原生宿主** — 產生的 Xcode 專案（iOS）或 Gradle 專案（Android），會嵌入已編譯的 Go 程式碼並承載 WebView。
- **資產提供** — `frontend/dist/`會嵌入 Go 二進位檔，並由程序內部提供。不需要 localhost 伺服器。

---

## 讓應用程式適應行動平台

應用程式已經可以運作，但在手機螢幕上看起來仍像桌面應用程式。只需幾項小幅調整，就能帶來顯著改善。

### 響應式 CSS

行動裝置螢幕較窄，且採用不同的輸入模式。在`frontend/public/style.css`（或對應檔案）中：

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### 在 Go 中偵測平台

使用建置標籤加入平台專屬行為，避免讓共用程式碼變得雜亂。

建立`mobile_ios.go`來放置僅限 iOS 的程式碼：

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

建立`mobile_android.go`來放置僅限 Android 的程式碼：

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

建立`mobile_desktop.go`作為虛設實作，讓共用程式碼也能在桌面平台上編譯：

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### 在 JavaScript 中偵測平台並限制行動平台專屬 UI

`IOS.*`和`Android.*`執行階段物件只存在於各自的平台。在桌面平台上呼叫它們會擲回例外。Kitchen Sink 採用的正確模式是只偵測一次平台，並完全隱藏僅限行動平台的控制項：

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

接著在 HTML 中：

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

如此一來，僅限行動平台的按鈕絕不會在桌面平台上轉譯，也不必對每個呼叫逐一使用`if (isMobile)`檢查加以防護。

在 Go 端，搭配使用建置標籤的虛設實作，讓事件處理常式只在需要它們的平台上註冊：

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

這正是[Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)採用的模式 — 請參閱`native_features_stub.go`、`native_features_ios.go`和`native_features_android.go`。

@note{type="note" title="原生功能 API 與事件命名"}
有兩項慣例值得瞭解：

- <strong>Go 端的原生功能使用平台管理器。</strong>請透過`application.IOS.*`和`application.Android.*`單例呼叫，例如`application.IOS.Haptic("medium")`或`application.Android.Share(payload)`。每個管理器只存在於其所屬平台，因此相關呼叫應放在`//go:build ios`／`//go:build android`檔案中。
- <strong>事件會依適用範圍劃分命名空間。</strong>兩個平台都能理解的事件使用`common:*`前綴（`common:haptic`、`common:location`……）；只有單一平台能產生或處理的事件則使用`ios:*`或`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）。由於幾乎所有行動平台功能都是共用的，因此前端會在`common:*`下為每種事件保留一個接聽程式。

@end

### 加入觸覺回饋（iOS）

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### 加入震動（Android）

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## 建置正式版本

@tabs{sync-key="mobile-platform"}
[iOS]
**模擬器建置**（用於在模擬器上測試，不需要簽署）：

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**裝置建置**（需要簽署身分與佈建描述檔）：

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**發佈用 IPA**（用於 App Store 或 TestFlight）：

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
若要上傳至 App Store Connect，請使用`wails3 task ios:xcode`，並讓 Xcode 管理簽署與封存；它會自動處理憑證、描述檔和公證等複雜事項。

@end

[Android]
**偵錯 APK**（使用 Android 偵錯金鑰庫簽署，可直接安裝）：

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**發行版 APK**（使用自己的金鑰庫簽署）：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**通用 APK**（單一檔案同時包含 arm64 與 x86_64）：

```bash
wails3 task android:package:fat
```

@note{type="tip"}
若要上傳至 Play Store，請產生`.aab`（Android App Bundle），而不是 APK；在 Android Studio 中開啟`build/android/`，然後使用<strong>Build → Generate Signed Bundle / APK</strong>。

@end

@end

---

## 疑難排解

### `wails3 task ios:run`失敗並顯示「no iOS SDKs found」

必須安裝並選取完整的 Xcode：

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run`失敗並顯示「SDK not found」

請確認已設定並匯出`ANDROID_HOME`。使用以下方式驗證：

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### 模擬器無法啟動

列出可用的模擬器，並手動啟動其中一個：

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect`未顯示任何目標

WebView 必須處於偵錯模式（`android:run`預設採用此模式）。請確認執行的是偵錯建置，而不是正式版本建置。另外也請確認`adb devices`顯示模擬器已連線。

#### 未套用安全區域內距

請確認您的 HTML 包含 viewport meta 標籤：

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## 探索 Kitchen Sink

第一個應用程式執行後，<strong>Kitchen Sink</strong>範例是瞭解其他可能功能的最快方式。這是一個完整的 Wails 應用程式，使用單一程式碼庫即可在 iOS、Android 和桌面平台上執行，涵蓋觸覺回饋、地理位置、生物辨識、本機通知、安全儲存等功能：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

請在[`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)瀏覽原始碼；其中`native_features_ios.go`和`native_features_android.go`檔案特別適合作為複製貼上後開始實作平台特定功能的起點。

## 下一步

@cards{cols="2"}
iOS 指南
完整參考資料：設定選項、原生分頁、WKWebView 切換選項、裝置建置和簽署。

[iOS 指南 →](/guides/mobile/ios/)

---
Android 指南
完整參考資料：設定、快顯通知、Play Store 封裝和 NDK 詳細資訊。

[Android 指南 →](/guides/mobile/android/)

---
📖 Kitchen Sink 原始碼
觸覺回饋、地理位置、生物辨識、通知和安全儲存，全都整合在一個可執行的應用程式中。

[在 GitHub 上檢視 →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
