---
title: "Android"
description: "在 Android 上构建和运行 Wails 应用程序——工具链设置、模拟器、APK 签名、Play 商店打包和 API 参考"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="实验性功能"}
Android 支持尚处于实验阶段，未来版本中可能会发生变化。

@end

@note{type="tip"}
刚开始使用 Wails 进行移动端开发？请先阅读[您的第一个移动应用 →](/guides/mobile/first-mobile-app/)，按照分步指南操作，然后返回此处查看完整参考。

@end

Wails v3 应用程序在 Android 上作为原生应用运行：Android `WebView` 负责渲染前端；资源由 Go 资源服务器支持的 `WebViewAssetLoader`在<strong>进程内</strong>提供（没有 localhost 服务器，也不开放 端口）；标准`@wailsio/runtime`无需修改即可使用——服务 绑定、事件、对话框和剪贴板均通过 Go 消息 处理器传递。

同一个`main.go`可为桌面端和 Android 构建。Go 代码会编译为 C 共享库（`libwails.so`、`GOOS=android` + NDK 工具链），并由 一个小型 Java 宿主加载。Android 特定行为位于受`//go:build android` 保护的各平台 Go 文件中。

## 要求

- 安装了 platform-tools、SDK 平台（API 35）、build-tools 和<strong>NDK</strong>（26.3.x）的<strong>Android SDK</strong>——`wails3 doctor`会显示它找到的组件
- 供 Gradle 使用的<strong>JDK</strong>（例如 OpenJDK 21）；如果`java`不在`PATH`中，请设置`JAVA_HOME`
- Go 1.25+ 和 npm
- 指向 SDK 的`ANDROID_HOME`（或`ANDROID_SDK_ROOT`）

使用命令行工具安装 SDK 组件：

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## 在模拟器上运行

在项目目录中运行：

```bash
wails3 task android:run
```

此命令会在没有模拟器运行时启动一个模拟器，然后生成绑定、构建 前端、针对模拟器的 ABI 将 Go 代码编译为`libwails.so`、 使用 Gradle 组装调试 APK，最后安装并启动它。

相关实用命令：

```bash
wails3 task android:logs    # stream the app's logcat output
```

在调试构建中，可以通过 Chrome 的`chrome://inspect`检查 WebView。

## 打包

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

生产构建使用`-tags production,android`，会剥离符号，并且在编译时 移除框架的内部诊断功能。`wails3 task android:package:fat` 会将`arm64-v8a`和`x86_64`都构建到单个 APK 中。

Google Play 要求新提交的应用使用 Android App Bundle（`.aab`）格式， 而且新提交的应用必须以 Android 15（API 35）或更高版本为目标； 项目模板在`build/android/app/build.gradle`中将`compileSdk`和`targetSdk`设置为35。 `wails3 task android:bundle:fat`会生成包含两个 ABI 的 `bin/<AppName>.aab`；Google Play 会据此生成针对各设备优化的 APK，因此包含所有 ABI 的 bundle 是上传到商店的正确产物。APK 仍然是进行本地和模拟器 测试的最快方式，因为`.aab`无法直接使用`adb`安装。

`android:run`和`android:deploy-emulator`是面向模拟器的任务。对于 Android 实体设备，请使用`android:run:device`构建调试 APK，或使用 `android:deploy-device`构建发布 APK。两者都会针对`arm64`进行构建，从 `adb devices`中选择第一个已连接的非模拟器条目，安装应用并启动 `com.wails.app.MainActivity`。传入`DEVICE_ID=<serial>`可指定特定 设备。

## 签名和发布构建

如果没有密钥库，发布构建会使用 Android **调试** 密钥库签名，以便安装测试。要使用自己的密钥库签名，请设置：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

同一组变量也用于签署 App Bundle：设置这些变量后运行`wails3 task android:bundle:fat`， 即可生成可提交至 Play 的`.aab`。如果未设置，bundle 会使用 调试密钥库签名，而 Google Play 会拒绝它，因此该任务会 输出警告。

@note{type="tip"}
使用[Play 应用签名](https://support.google.com/googleplay/android-developer/answer/9842756)时， 您在本地用于签名的密钥库就是您的<strong>上传密钥</strong>：Google 会使用它验证 您的上传内容，然后使用其管理的应用签名密钥为应用重新 签名。另请注意，Google Play 要求每次上传使用更高的`versionCode`； 请在`build/android/app/build.gradle`中递增该值。

@end

## 配置

前端在运行时通过`Android`运行时 对象使用 Android 功能：`Android.Haptics.Vibrate(durationMs)`、`Android.Device.Info()`、 `Android.Toast.Show(message)`。软件包名称由构建任务中的`APP_ID` 控制。

## 支持与不支持的功能

| 领域 | 状态 |
| --- | --- |
| WebView + 进程内资源（`WebViewAssetLoader`） | ✅ |
| 服务绑定、事件（双向） | ✅ |
| 消息对话框 | ✅ 带按钮回调的 AlertDialog |
| 打开文件/多个文件对话框 | ✅ Storage Access Framework（文件以缓存副本的形式导入） |
| 打开目录/保存文件对话框 | ❌ 返回错误——请改为写入应用沙盒 |
| 剪贴板 | ✅ ClipboardManager |
| 屏幕 API | ✅ WindowMetrics，包括系统栏之外的工作区域 |
| 生命周期事件（`events.Android.*`） | ✅ |
| 触觉反馈、设备信息、Toast | ✅ `Android.*`运行时 API |
| 模拟器和实体设备构建 | ✅ `android:run`、`android:run:device`、`android:deploy-emulator`、`android:deploy-device` |
| 窗口几何属性、菜单、系统托盘 | 有意设计为空操作 |
| 多个窗口 | 仅显示第一个窗口 |

## 移植说明

- 桌面端代码在`GOOS=android`下无需修改即可编译；由于 Android 应用采用全屏模式，几何属性、菜单和托盘调用会变为空操作。
- `android`**隐含`linux`构建标签**（Android 使用 Linux 内核）：仅用于桌面 Linux 的文件需要`//go:build linux && !android`，而在运行时，`runtime.GOOS`为`"android"`。
- 请将保存文件和选择目录对话框替换为写入应用沙盒并通过 intent 分享的流程。打开文件对话框可以正常工作，并会将所选文档作为副本导入缓存目录，因此您会获得真实的文件系统路径。
- 实际应用始终使用`CGO_ENABLED=1`和 NDK 构建；非 cgo 路径仅用于让`wails3 generate bindings`等工具能够加载该包。
- 采用响应式设计来构建前端；`Screens`工作区不包括状态栏和导航栏。
