---
title: "你的第一个移动应用"
description: "只需几分钟，即可在 iOS 模拟器或 Android 模拟器上运行 Wails 应用"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

本指南将指导你在 iOS 模拟器或 Android 模拟器上运行标准的 Wails 桌面应用。<strong>无需更改 Go 代码。</strong>同一个`main.go`可为所有目标平台构建。

<strong>完成所需时间：</strong>15–30 分钟（大部分时间用于首次运行时安装工具链）

## 从桌面项目开始

如果还没有项目，请新建一个：

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

首先确认桌面应用能够正常运行：

```bash
wails3 dev
```

应用打开后，将其退出并继续后续步骤。桌面端可运行的所有内容也都能在移动端运行——本指南无需修改`main.go`或任何 Go 代码。

---

## 选择平台

@tabs{sync-key="mobile-platform"}
[iOS 模拟器]
### 要求

- **macOS**（只能在 macOS 上构建 iOS 应用）
- **完整的 Xcode**——不能只有命令行工具。请从 App Store 安装，然后运行：
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.26+** 和**npm**（如果运行过`wails3 init`，则已安装）

运行`wails3 doctor`进行验证——它会列出能够找到的 iOS SDK。

### 在模拟器上运行

@steps
### 启动应用
```bash
wails3 task ios:run
```

就这么简单——该操作会构建应用；如果模拟器尚未运行，则将其启动；然后启动应用。

@note{type="tip"}
首次运行需要几分钟（系统正在为 iOS 编译并缓存 Wails 框架）。之后每次运行都会快得多。

@end

启动后，未经修改的桌面应用便会在 iOS 模拟器上运行——使用相同的`main.go`和相同的前端：

![在 iOS 模拟器上运行的默认 Wails 应用](/assets/ios-simulator-first-app.png)

### 实时查看日志
在另一个终端中运行：

```bash
wails3 task ios:logs:dev
```

该命令会持续输出模拟器日志，并筛选出应用的相关内容。`fmt.Println`和`log.Println`的输出会显示在这里。

### 检查 WebView
在 Safari 中，选择**开发 → 模拟器 → 你的应用**。完整的 Web 检查器均可使用，包括控制台、调试器和网络面板等所有功能。

### 进行更改
编辑任意前端文件（`frontend/src/main.js`、`index.html`等），然后重新运行`wails3 task ios:run`。Wails 会重新构建前端并重新启动应用。

更改 Go 代码后，也要重新运行`wails3 task ios:run`。Go 采用增量重新编译，因此只会重新构建已更改的软件包。

@end

### 在 Xcode 中打开（可选）

```bash
wails3 task ios:xcode
```

该命令会在 Xcode 中打开`build/ios/`。你可以使用 Xcode 将应用部署到设备、进行高级性能分析或管理预置描述文件。Wails 会在每次构建时重新生成 Xcode 项目，因此不要直接修改生成的文件。

[Android 模拟器]
### 要求

你需要**Android SDK**、**NDK**和**JDK**。最简单的方式是使用 Android Studio，也可以使用命令行工具：

@steps
### 安装 Android 命令行工具
从[developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only)下载，并解压到`~/android-sdk/cmdline-tools/latest/`。

### 安装 SDK 组件
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### 创建模拟器
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### 设置环境变量
将以下内容添加到`~/.zshrc`或`~/.bashrc`：

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

重新加载：`source ~/.zshrc`

### 安装 JDK
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

运行`wails3 doctor`，确认已找到所有组件。

### 在模拟器上运行

@steps
### 启动应用
```bash
wails3 task android:run
```

首次运行时，该命令会：

- 如果没有模拟器正在运行，则启动模拟器
- 生成绑定并构建前端
- 通过 NDK 交叉编译器将 Go 代码编译为`libwails.so`
- 使用 Gradle 组装调试版 APK
- 在模拟器上安装并启动 APK

@note{type="tip"}
首次构建会下载 Gradle 并编译 NDK 工具链——预计需要5–10 分钟。后续构建采用增量方式，只需不到一分钟。

@end

### 实时查看日志
在另一个终端中运行：

```bash
wails3 task android:logs
```

该命令会运行`adb logcat`，并筛选出应用的相关内容。`fmt.Println`的输出会显示在这里。

### 检查 WebView
打开 Chrome 并转到`chrome://inspect`。应用的 WebView 会显示在**远程目标**下——点击**检查**即可打开 DevTools。

### 进行更改
编辑任意文件，然后重新运行`wails3 task android:run`。Gradle 采用增量构建，因此只会重新编译已更改的代码。

@end

@end

---

## 了解背后的过程

你的`main.go`完全没有变化。所有工作都由 Wails 完成：

- **构建系统**——项目中的`Taskfile.yml`包含用于驱动平台特定工具链的`ios:*`和`android:*`任务。
- **Go 交叉编译**——使用`GOOS=ios`或`GOOS=android`，并配合相应的`GOARCH`和 sysroot。
- **原生宿主**——生成的 Xcode 项目（iOS）或 Gradle 项目（Android），用于嵌入已编译的 Go 代码并承载 WebView。
- **资源提供**——`frontend/dist/`会嵌入 Go 二进制文件，并在进程内提供，无需 localhost 服务器。

---

## 让应用适配移动端

应用已经可以运行，但在手机屏幕上看起来仍像桌面应用。只需进行几项小改动，就能显著改善体验。

### 响应式 CSS

移动设备的屏幕更窄，输入方式也有所不同。在`frontend/public/style.css`（或等效文件）中：

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

### 在 Go 中检测平台

使用构建标签添加平台特定行为，避免让共享代码变得杂乱。

为仅限 iOS 的代码创建`mobile_ios.go`：

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

为仅限 Android 的代码创建`mobile_android.go`：

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

创建`mobile_desktop.go`作为存根，使共享代码也能在桌面平台上编译：

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### 在 JavaScript 中检测平台并限制仅限移动端的 UI

`IOS.*`和`Android.*`运行时对象仅存在于各自对应的平台上。在桌面端调用它们会引发异常。正确的模式（Kitchen Sink 采用的模式）是只检测一次平台，并完全隐藏仅限移动端的控件：

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

然后在 HTML 中：

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

这样，移动端专用按钮绝不会在桌面端渲染，也无需为每个调用分别添加 `if (isMobile)` 检查。

在 Go 端，配合使用构建标签存根，确保只在需要这些事件处理程序的平台上注册它们：

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

这正是 [Kitchen Sink](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) 使用的模式——请参阅 `native_features_stub.go`、`native_features_ios.go` 和 `native_features_android.go`。

@note{type="note" title="原生功能 API 与事件命名"}
有两个约定值得了解：

- <strong>Go 端的原生功能使用平台管理器。</strong>请通过 `application.IOS.*` 和 `application.Android.*` 单例调用它们，例如 `application.IOS.Haptic("medium")` 或 `application.Android.Share(payload)`。每个管理器仅存在于其对应平台上，因此相关调用应放在 `//go:build ios` / `//go:build android` 文件中。
- <strong>事件按适用范围划分命名空间。</strong>两个平台都支持的事件使用 `common:*` 前缀（`common:haptic`、`common:location`……）；只有一个平台能够产生或处理的事件使用 `ios:*` 或 `android:*`（例如 `ios:backgroundTask`、`android:foregroundService`）。由于几乎所有移动端功能都是共享的，因此前端只需在 `common:*` 下为每个事件保留一个监听器。

@end

### 添加触觉反馈（iOS）

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### 添加振动（Android）

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## 构建生产版本

@tabs{sync-key="mobile-platform"}
[iOS]
**模拟器构建**（用于在模拟器上测试，无需签名）：

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**设备构建**（需要签名身份和预置描述文件）：

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**分发 IPA**（用于 App Store 或 TestFlight）：

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
如需上传到 App Store Connect，请使用 `wails3 task ios:xcode`，并让 Xcode 管理签名和归档——它会自动处理证书、描述文件和公证等复杂事项。

@end

[Android]
**调试 APK**（使用 Android 调试密钥库签名，可直接安装）：

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**发布 APK**（使用你自己的密钥库签名）：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**通用 APK**（在一个文件中包含 arm64 和 x86_64）：

```bash
wails3 task android:package:fat
```

@note{type="tip"}
如需上传到 Play Store，请生成 `.aab`（Android App Bundle），而不是 APK——在 Android Studio 中打开 `build/android/`，然后选择 **Build → Generate Signed Bundle / APK**。

@end

@end

---

## 故障排除

### `wails3 task ios:run` 失败并显示“no iOS SDKs found”

必须安装并选择完整版 Xcode：

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` 失败并显示“SDK not found”

确保已设置并导出 `ANDROID_HOME`。使用以下命令验证：

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### 模拟器无法启动

列出可用的模拟器，并手动启动一个：

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` 未显示任何目标

WebView 必须处于调试模式（`android:run` 默认启用此模式）。请确保运行的是调试构建，而不是生产构建。还要确认 `adb devices` 显示模拟器已连接。

#### 安全区域内边距未生效

确保 HTML 中包含 viewport 元标记：

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## 探索 Kitchen Sink

第一个应用运行后，通过 **Kitchen Sink** 示例可以最快了解其他可实现的功能。这是一个完整的 Wails 应用，通过同一套代码库在 iOS、Android 和桌面端运行，涵盖触觉反馈、地理位置、生物识别、本地通知、安全存储等功能：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

在 [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) 浏览源代码——其中 `native_features_ios.go` 和 `native_features_android.go` 文件特别适合作为实现平台特定功能时复制粘贴的起点。

## 后续步骤

@cards{cols="2"}
iOS 指南
完整参考：配置选项、原生标签页、WKWebView 开关、设备构建和签名。

[iOS 指南 →](/guides/mobile/ios/)

---
Android 指南
完整参考：配置、Toast、Play Store 打包和 NDK 详细信息。

[Android 指南 →](/guides/mobile/android/)

---
📖 Kitchen Sink 源代码
触觉反馈、地理位置、生物识别、通知和安全存储——全部包含在一个可运行的应用中。

[在 GitHub 上查看 →](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile)

@end
