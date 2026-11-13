---
title: "移动端概览"
description: "使用与桌面应用相同的 Go 代码库构建 iOS 和 Android 应用"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3可使用你为桌面端编写的同一套`main.go`和前端，在**iOS 和 Android**上运行。无需单独创建移动端项目，无需代码共享桥接，也无需重写：Go 二进制文件会针对移动端目标进行编译，并由原生 WebView 渲染现有前端。

@cards{cols="2"}
iOS
WKWebView + UIKit 宿主。资源通过自定义`wails://`方案提供，无需开放端口。需要装有完整版 Xcode 的**macOS**。

[iOS 指南 →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`。Go 通过 NDK 编译为`libwails.so`。支持 macOS、Linux 和 Windows。

[Android 指南 →](/guides/mobile/android/)

@end

## 查看实际运行效果：Kitchen Sink 示例

了解其功能的最佳方式是查看**Kitchen Sink**：这是一个基于单一代码库构建的 Wails 应用，可在 iOS、Android 和桌面端以相同方式运行：

@linkcard{title="移动端 Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="绑定 · 事件 · 对话框 · 触觉反馈 · 地理定位 · 生物识别 · 通知 · 安全存储 · 更多功能——全部集中在一个 main.go 中"}
它通过7个标签页展示了所有主要的移动端 API 功能，而且也能在桌面端运行。前端通过平台检查在桌面端隐藏**移动端**和**硬件**标签页；针对桌面端构建时，Go 端不会为`common:*`移动端事件注册任何处理程序。这是使用一套代码库面向所有平台发布应用的推荐模式。

| 标签页 | 平台 | 展示内容 |
| --- | --- | --- |
| **绑定** | 全部 | JS → Go 服务调用，返回值、结构体和错误 |
| **事件** | 全部 | Go → JS 时钟、JS → Go → JS ping/pong、操作系统事件（电池、网络、主题） |
| **对话框** | 全部 | 各平台的原生消息对话框 |
| **系统** | 全部 | 剪贴板、屏幕指标、设备信息 |
| **移动端** | iOS + Android | 共享面板、保持唤醒、手电筒、亮度、生物识别、本地通知、安全存储 |
| **硬件** | iOS + Android | 触觉反馈、地理定位、加速度计、距离感应、文本转语音 |
| **原生功能** | iOS + Android | iOS：触觉反馈 + WKWebView 开关 · Android：振动 + Toast 消息 |

自行运行：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 工作原理

所有平台都采用相同的应用模型：

1. **Go 后端**——你的服务、事件处理程序和应用逻辑无需更改，即可针对`GOOS=ios`和`GOOS=android`进行编译。
2. **前端**——使用完全相同的 HTML/JS/CSS。`@wailsio/runtime`包的工作方式完全相同；服务绑定、事件、对话框和剪贴板都通过同一进程内传输机制进行路由。
3. **WebView 宿主**——iOS 上是在`UIViewController`内运行`WKWebView`；Android 上是在`Activity`内运行`WebView`。Wails 会自动连接消息桥。
4. **进程内资源服务**——资源直接由 Go 内存提供，而不是通过 localhost 服务器提供。无需开放端口，无需环回连接，也没有额外延迟。

平台特定行为位于受`//go:build ios`或`//go:build android`保护的文件中，从而保持共享代码整洁。

## 前置条件概览

| 要求 | iOS | Android |
| --- | --- | --- |
| 操作系统 | 仅限 macOS | macOS、Linux、Windows |
| 工具链 | 完整版 Xcode（不能只有 CLI 工具） | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 验证命令 | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
设置好工具链后运行`wails3 doctor`——它会准确显示各平台已找到和缺少的组件。

@end

## 支持的功能

两个平台共享同一组核心功能：

| 功能 | iOS | Android |
| --- | --- | --- |
| 服务绑定（JS → Go） | ✅ | ✅ |
| 事件（双向） | ✅ | ✅ |
| 消息对话框 | ✅ UIAlertController | ✅ AlertDialog |
| 打开文件对话框 | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| 保存文件对话框 | ❌ 改为写入沙盒 | ❌ 改为写入沙盒 |
| 剪贴板 | ✅ UIPasteboard | ✅ ClipboardManager |
| 屏幕/安全区域指标 | ✅ | ✅ |
| 生命周期事件 | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| 触觉反馈 | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| 设备信息 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| 原生标签页（iOS） | ✅ UITabBar | — |
| Toast 消息（Android） | — | ✅ `Android.Toast.Show` |
| 多个窗口 | ❌ 仅限第一个窗口 | ❌ 仅限第一个窗口 |
| 窗口几何属性 / 菜单 / 托盘 | 有意不执行任何操作 | 有意不执行任何操作 |

## 构建标签规则

编写平台条件代码时，需要了解两条重要规则：

- **`ios`隐含`darwin`**——带有`//go:build darwin`标签的文件也会为 iOS 编译。若仅以 macOS 为目标，请使用`//go:build darwin && !ios`。
- **`android`隐含`linux`**——带有`//go:build linux`标签的文件也会为 Android 编译。若仅以桌面 Linux 为目标，请使用`//go:build linux && !android`。

在运行时，`runtime.GOOS`会分别返回`"ios"`和`"android"`。

## 运行时平台检测

构建标签适用于只能在给定平台上*编译*的代码。对于共享代码中的常规分支，请使用`application.System`——它在每个构建中都可用（无需构建标签），因此同一个文件可在所有平台上运行：

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

可用值：`IsMobile()`、`IsDesktop()`、`IsServer()`（`server`构建标签）和`IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`。

前端在`@wailsio/runtime`中提供了对应的辅助函数：

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 后续步骤

@cards{cols="2"}
🚀 你的第一个移动应用
只需几分钟，即可将桌面 Wails 应用运行在 iOS 模拟器或 Android 模拟器中。

[开始使用 →](/guides/mobile/first-mobile-app/)

---
iOS 指南
完整介绍 iOS 工具链设置、模拟器、设备构建、签名、配置和 API 参考。

[iOS 指南 →](/guides/mobile/ios/)

---
Android 指南
完整介绍 Android SDK/NDK 设置、模拟器、APK 签名、Play Store 打包和 API 参考。

[Android 指南 →](/guides/mobile/android/)

@end
