---
title: "iOS"
description: "在 iOS 上构建并运行 Wails 应用——涵盖环境设置、模拟器、设备构建、配置和原生功能"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="实验性功能"}
iOS 支持尚处于实验阶段，可能会在未来版本中发生变化。

@end

@note{type="tip"}
刚开始使用 Wails 进行移动开发？请先阅读[您的第一个移动应用 →](/guides/mobile/first-mobile-app/)，按照分步指南完成操作，然后返回此处查阅完整参考。

@end

Wails v3 应用能以完全原生的应用形式在 iOS 上运行——而且最棒的是，其工作方式与桌面版本<em>完全</em>相同。同一个 Go 后端、同一个前端、同一套`@wailsio/runtime`：服务绑定、事件、对话框和剪贴板的行为全都一致，<strong>完全不需要</strong>针对移动端重新接线。无需单独维护移动端代码库，无需移植层，也无需学习特殊 API——现有 Wails 应用可直接在 iOS 上运行。移植过程真正无缝：原样迁移应用即可发布。

同一个`main.go`可同时为桌面端和 iOS 构建；iOS 特有的调整通过`application.Options.IOS`进行配置。

## 要求

- 运行 macOS，并已安装<strong>完整的 Xcode</strong>（仅安装命令行工具还不够）——`wails3 doctor`会显示其能够找到的 iOS SDK
- Go 1.25+ 和 npm

## 模拟器

在项目目录中运行：

```bash
wails3 task ios:run
```

此操作会构建应用；如果尚无模拟器在运行，则会启动一个模拟器，然后启动应用。

实用的配套命令：

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

在调试构建中，可通过 Safari 的“开发”菜单检查 WebView。

## 打包

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

这些是经过优化并移除调试信息的生产构建。

## 设备构建

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device`为实体设备构建。授权配置来自`build/ios/entitlements.plist`，且仅应用于设备构建——请添加应用所需的 capability 键。

@note{type="tip"}
如需自动管理签名、预置和 App Store 归档，请使用`wails3 task ios:xcode`打开生成的 Xcode 项目，然后改为从 Xcode 构建。

@end

## 配置

`build/config.yml`：

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

启动选项（`application.Options.IOS`）包括`DisableScroll`、`DisableBounce`、`DisableScrollIndicators`、`DisableInputAccessoryView`、`EnableBackForwardNavigationGestures`、`DisableLinkPreview`、`EnableInlineMediaPlayback`、`EnableAutoplayWithoutUserAction`、`DisableInspectable`、`UserAgent`、`ApplicationNameForUserAgent`、`BackgroundColour`，以及通过`EnableNativeTabs` + `NativeTabsItems`配置的原生底部标签页。

## 原生功能

iOS 特有的功能可通过`application.IOS`使用。请在`//go:build ios`文件中从 Go 调用它们，以使共享代码保持平台无关。Android 通过`application.Android`提供相同的功能集。

一次性操作会立即返回：

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

查询辅助函数以 JSON 形式返回结果，包括`SafeAreaJSON()`、`AppInfoJSON()`、`PowerJSON()`、`NetworkJSON()`、`StorageJSON()`、`GetOrientation()`、`GetBrightness()`。`StoragePath()`返回应用的 Application Support 目录绝对路径；该目录很适合存放数据库及其他持久化文件（这是 iOS 上与 Android 的`getFilesDir()`所返回目录对应的目录）。首次访问时会创建该目录；如果无法创建，`StoragePath()`会返回空字符串，因此使用前请检查是否为`""`。

### 事件

任何稍后才会完成的操作——例如权限提示、传感器数据流或相机拍摄——都会以<strong>事件</strong>而非返回值的形式提供结果；您可以在 Go 或前端中监听这些事件。与 Android 共享的功能，其事件名称带有`common:`前缀；仅限 iOS 的功能，其事件名称则带有`ios:`前缀。

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

| 事件 | 触发方式 | 载荷 |
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

`v3/examples/mobile`下的综合示例端到端串联了上述所有功能。

## WebView 控制

还可以在运行时通过 Go 更改部分 WebView 行为：

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

随附的`@wailsio/runtime`还公开了一个小型<strong>前端</strong> iOS 命名空间：

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

原生底部标签页的选择操作会作为`nativeTabSelected`事件传入`window`。

## 支持状态

| 领域 | 状态 |
| --- | --- |
| 前端渲染和资源 | ✅ |
| 服务绑定、事件（双向） | ✅ |
| 消息对话框 | ✅ |
| 打开文件/多个文件/目录对话框 | ✅ 作为沙盒副本导入 |
| 保存文件对话框 | ❌ 请改为写入应用沙盒 |
| 剪贴板 | ✅ |
| 屏幕 API | ✅ 包含安全区域内的工作区 |
| 生命周期事件 | ✅ |
| 窗口几何属性、菜单、系统托盘 | 在 iOS 上不执行任何操作 |
| 多个窗口 | 仅显示第一个窗口 |

## 移植说明

- 桌面端代码无需修改即可为 iOS 构建——窗口、菜单和系统托盘调用不会执行任何操作。
- 将保存文件对话框替换为写入应用沙盒后再共享。
- 采用响应式设计构建前端；安全区域已自动处理。
