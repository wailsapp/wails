---
title: "移动端 API"
description: "跨平台的 application.Mobile 管理器——一个受构建约束保护的入口点，用于访问 iOS 和 Android 共有的原生移动端功能"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="实验性功能"}
移动端支持尚处于实验阶段，未来版本中可能会发生变化。

@end

原生移动端功能通过以下两种方式提供：

- **平台专用管理器**——`application.IOS`（位于`//go:build ios`文件中）和`application.Android`（位于`//go:build android`文件中）。所有平台特定功能都应使用这些管理器。有关各平台提供的完整 API，请参阅[iOS](/guides/mobile/ios/)和[Android](/guides/mobile/android/)参考文档。
- **`application.Mobile`**——一个受构建约束保护的统一管理器，涵盖在两个平台上行为相同的功能子集。如果希望使用一套可在所有平台上编译和运行的代码路径，请使用此管理器。

## `application.Mobile`

`application.Mobile`在 iOS 上分派给`IOS`，在 Android 上分派给`Android`，在桌面端则分派给不执行任何操作的存根。由于它没有构建约束，因此可以在普通的、与平台无关的 Go 代码中调用，无需自行创建`//go:build`文件：

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

`StoragePath()`返回应用私有文件目录的绝对路径——Android 上为`getFilesDir()`，iOS 上为 Application Support 目录——这是存放数据库和其他持久化文件的推荐位置。在桌面端，或者在设备上该目录不可用时（在 iOS 上则为无法创建时），它会返回空字符串，因此在使用前应检查是否为`""`。

@note{type="note"}
在设备之外（桌面端构建中），每个`Mobile`方法都不执行任何操作，每个查询都返回其零值（字符串返回`""`）。因此，跨平台代码可以无条件调用`application.Mobile.*`。如果桌面端也需要实际路径，请根据平台进行分支，并回退到`os.UserConfigDir()`或类似方案。

@end

## 功能

`Mobile`管理器提供在 iOS 和 Android 上具有相同签名的功能：

| 功能 | API | 备注 |
| --- | --- | --- |
| 共享面板 | `Mobile.Share(json)` | `{text, url}` |
| 在外部打开 URL | `Mobile.OpenURL(url)` | 系统浏览器 |
| 保持屏幕唤醒 | `Mobile.SetKeepAwake(bool)` |  |
| 手电筒 | `Mobile.SetTorch(bool)` | → `common:torch` |
| 安全区域边距 | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| 应用信息 | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| 屏幕方向锁定 | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| 状态栏 | `Mobile.SetStatusBar(json)` | 样式和可见性 |
| 存储信息 | `Mobile.StorageJSON()` | `{free,total}`字节 |
| 存储路径 | `Mobile.StoragePath()` | 应用私有文件目录 |
| 电源 / 电池 | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| 网络状态 | `Mobile.NetworkJSON()` | `{connected,type}` |
| 生物识别 | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| 安全存储 | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| 地理位置 | `Mobile.GetLocation()` | 单次获取 → `common:location` |
| 触觉反馈 | `Mobile.Haptic(type)` | 冲击 / 通知 / 选择 |
| 加速度计 | `Mobile.SetMotion(bool)` | → `common:motion` |
| 接近感应 | `Mobile.SetProximity(bool)` | → `common:proximity` |
| 文本转语音 | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| 键盘边距 | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| 屏幕捕获 | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| 相机 | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

异步结果以`common:*`事件的形式到达，与各平台管理器完全相同——有关载荷，请参阅[事件](/guides/mobile/ios/#events)。

## 仍由平台专门提供的功能

iOS 与 Android 之间形式不同的功能<strong>不</strong>位于`Mobile`上；请在带有构建标签的文件中，通过`application.IOS` / `application.Android`调用这些功能：

| 功能 | iOS | Android |
| --- | --- | --- |
| 亮度（设置） | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| 亮度 / 屏幕方向（获取） | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| 本地通知 | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| 安全存储（写入） | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| 后台执行 | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
`MobileManager`接口是`application.Mobile`背后的契约。由于两个平台管理器都必须实现该接口，因此上述任何方法在 iOS 和 Android 上都保证保持完全相同的签名——如果它们出现差异，对应平台的构建将无法编译。

@end
