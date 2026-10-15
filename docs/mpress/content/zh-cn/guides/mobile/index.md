---
title: "移动端概览"
description: "使用与桌面应用相同的 Go 代码库构建 iOS 和 Android 应用"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 可使用你已经为桌面端编写的同一套`main.go`和前端在<strong>iOS 和 Android</strong>上运行。无需单独创建移动端项目，无需代码共享桥接，也无需重写：Go 二进制文件会针对移动端目标进行编译，并由原生 WebView 渲染现有前端。

@cards{cols="2"}
iOS
WKWebView + UIKit 宿主。通过自定义`wails://`方案提供资源，无需开放端口。 需要装有完整 Xcode 的<strong>macOS</strong>。

[iOS 指南 →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`。Go 通过 NDK 编译为`libwails.so`。 可在 macOS、Linux 和 Windows 上使用。

[Android 指南 →](/guides/mobile/android/)

@end

## 查看实际运行效果：Kitchen Sink 示例

要了解可以实现哪些功能，最好的方式是查看<strong>Kitchen Sink</strong>。这是一个使用单一代码库、在 iOS、Android 和桌面端上以相同方式运行的 Wails 应用：

@linkcard{title="移动端 Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="绑定 · 事件 · 对话框 · 触觉反馈 · 地理位置 · 生物识别 · 通知 · 安全存储 · 以及更多功能——全部来自一个 main.go"}
它通过7个选项卡展示了所有主要的移动端 API 功能，而且也能在桌面端运行。前端通过平台检查在桌面端隐藏<strong>移动端</strong>和<strong>硬件</strong>选项卡；构建桌面端版本时，Go 端不会为`common:*`移动端事件注册任何处理程序。这是在所有平台上交付同一套代码库的推荐模式。

| 选项卡 | 平台 | 展示内容 |
| --- | --- | --- |
| **绑定** | 全部 | JS → Go 服务调用，可返回值、结构体和错误 |
| **事件** | 全部 | Go → JS 时钟、JS → Go → JS ping/pong、操作系统事件（电池、网络、主题） |
| **对话框** | 全部 | 各平台的原生消息对话框 |
| **系统** | 全部 | 剪贴板、屏幕指标、设备信息 |
| **移动端** | iOS + Android | 系统分享面板、保持唤醒、手电筒、亮度、生物识别、本地通知、安全存储 |
| **硬件** | iOS + Android | 触觉反馈、地理位置、加速度计、距离传感器、文本转语音 |
| **原生功能** | iOS + Android | iOS：触觉反馈 + WKWebView 开关 · Android：振动 + toast 提示 |

要自行运行此示例：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 工作原理

所有平台均采用相同的应用模型：

1. **Go 后端**——你的服务、事件处理程序和应用逻辑无需修改，即可为`GOOS=ios`和`GOOS=android`编译。
2. **前端**——使用完全相同的 HTML/JS/CSS。`@wailsio/runtime`包的工作方式也完全相同；服务绑定、事件、对话框和剪贴板都通过同一进程内传输机制进行路由。
3. **WebView 宿主**——在 iOS 上，是位于`UIViewController`中的`WKWebView`；在 Android 上，是位于`Activity`中的`WebView`。Wails 会自动连接消息桥。
4. **进程内资源服务**——资源直接从 Go 内存中提供，而不是通过 localhost 服务器。无需开放端口，无需环回连接，也不会增加额外延迟。

平台特定行为位于受`//go:build ios`或`//go:build android`保护的文件中，从而让共享代码保持整洁。

## 前置条件速览

| 要求 | iOS | Android |
| --- | --- | --- |
| 操作系统 | 仅限 macOS | macOS、Linux、Windows |
| 工具链 | 完整的 Xcode（不能只有 CLI 工具） | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 验证命令 | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
设置好工具链后运行`wails3 doctor`——它会准确显示每个平台已找到和缺少的组件。

@end

## 支持的功能

两个平台具有相同的核心功能集：

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
| 多窗口 | ❌ 仅支持第一个窗口 | ❌ 仅支持第一个窗口 |
| 窗口几何属性/菜单/托盘 | 有意设计为空操作 | 有意设计为空操作 |

## 构建标签规则

编写平台条件代码时，需要了解两条重要规则：

- **`ios` 隐含 `darwin`**——标记为 `//go:build darwin` 的文件也会针对 iOS 进行编译。若要仅针对 macOS，请使用 `//go:build darwin && !ios`。
- **`android` 隐含 `linux`**——标记为 `//go:build linux` 的文件也会针对 Android 进行编译。若要仅针对桌面 Linux，请使用 `//go:build linux && !android`。

在运行时，`runtime.GOOS`会分别返回 `"ios"` 和 `"android"`。

## 运行时平台检测

构建标签适用于只能在特定平台上<em>编译</em>的代码。对于共享代码中的常规分支，请使用 `application.System`——它在每个构建中都可用 （无需构建标签），因此同一文件可在所有平台上使用：

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

可用值：`IsMobile()`、`IsDesktop()`、`IsServer()`（对应 `server` 构建标签）以及 `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`。

前端在 `@wailsio/runtime` 中提供了相应的辅助函数：

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 后续步骤

@cards{cols="2"}
🚀 创建你的第一个移动应用
只需几分钟，即可让桌面版 Wails 应用在 iOS Simulator 或 Android Emulator 上运行。

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
