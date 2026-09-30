---
title: "权限"
description: "控制网页内容对摄像头、麦克风、地理位置及其他功能的请求"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

调用`navigator.mediaDevices.getUserMedia()`、Geolocation API 或 Notifications API 的网页内容需要宿主应用授予或拒绝这些请求。Wails 在`WebviewWindowOptions`上公开了跨平台的`Permissions`映射，让你能够以声明方式进行控制，无需编写平台特定代码。

## 快速入门

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

来自该窗口网页内容的摄像头和麦克风请求将直接获得授权，不会显示浏览器提示。

## 权限类型

`PermissionType`（uint8）标识网页内容可以请求的功能。

| 常量 | 功能 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 权限值

`Permission`（uint8）是应用于指定类型的策略。

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 使用平台的原生处理方式（见下文） |
| `PermissionAllow` | 1 | 直接授予，无需提示 |
| `PermissionDeny` | 2 | 直接拒绝，无需提示 |

`PermissionDefault`是零值，因此映射中未设置的条目会按默认方式处理。

## 平台行为

由于底层 WebView 的原生行为不同，各平台对`PermissionDefault`的处理方式也不同。

### Linux（WebKitGTK）

WebKitGTK<strong>没有原生权限提示</strong>。如果未附加处理程序，它会静默拒绝所有请求——这就是添加此功能之前`getUserMedia`始终返回`NotAllowedError`的原因。

Wails 目前可在 Linux 上处理<strong>摄像头和麦克风</strong>请求。地理位置、通知和剪贴板读取功能尚未接入，因此无论你设置什么策略，这些请求仍会被拒绝。

| 策略 | 摄像头/麦克风 | 地理位置、通知、剪贴板 |
| --- | --- | --- |
| `PermissionDefault` | **允许**（恢复 getUserMedia） | 始终拒绝 |
| `PermissionAllow` | 允许 | 始终拒绝（尚未实现） |
| `PermissionDeny` | 拒绝 | 始终拒绝 |

### Windows（WebView2）

WebView2 提供原生权限提示和按类型设置权限的 API。它完全支持所有五种功能类型。

| 策略 | 行为 |
| --- | --- |
| `PermissionDefault` | WebView2 显示其原生操作系统/浏览器权限提示 |
| `PermissionAllow` | 静默授予 |
| `PermissionDeny` | 静默拒绝 |

<strong>重要：</strong>在此功能推出之前，Wails 会无条件调用`SetGlobalPermission(Allow)`，静默授予所有功能。现在，只要`Permissions`中存在任何条目，就<strong>不会</strong>设置这种全局授权。未设置的功能将交由 WebView2 显示原生提示，而不是自动获得授权。

这意味着，只要你在 Windows 上配置了`Permissions`，任何未明确列出的功能都会显示提示，而不会被静默允许。请明确设置所需的功能。

### macOS (WKWebView + TCC)

macOS 上需要两层授权。WKWebView 在启动采集会话前询问应用，由 `Permissions` 映射回答该请求。在其下层，系统隐私框架 TCC 控制设备本身：应用首次实际访问摄像头或麦克风时会显示系统提示，用户选择按应用保存在“系统设置 → 隐私与安全性”中。

Wails 在 macOS 12 及更高版本上处理**摄像头和麦克风**请求。地理位置、通知和剪贴板读取没有对应的 `WKUIDelegate` 接口，因此尚未接入，仍交由 TCC 处理。与 Linux 一样，为这些功能设置的策略没有效果。

| 策略 | 摄像头 / 麦克风 | 地理位置、通知、剪贴板 |
| --- | --- | --- |
| `PermissionDefault` | WebKit 显示自己的权限提示 | 仅由 TCC 处理 |
| `PermissionAllow` | 跳过 WebKit 提示 — **TCC 仍然适用** | 仅由 TCC 处理 |
| `PermissionDeny` | 在访问设备之前拒绝 | 仅由 TCC 处理 |

`PermissionAllow` 授权的是 WebView 请求，而不是设备本身。首次采集仍会触发 TCC 提示；用户在系统设置中拒绝的应用仍然被拒绝，应用无法自行授予设备访问权限。`PermissionAllow` 省略的只是前置的 WebKit 提示。

macOS 12 之前不存在此委托方法，因此会忽略映射，所有请求都回退到 WebKit 提示。

请确保你的`Info.plist`包含相应的用途说明键：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="拒绝 Info.plist 未声明用途的功能"}

如果请求在缺少对应用途说明键的情况下到达 AVFoundation，结果不只是请求失败：macOS 会终止应用。

`PermissionDefault` 是零值，因此只设置 `{PermissionMicrophone: PermissionAllow}` 会让摄像头继续使用 WebKit 提示。如果用户接受提示，而应用只声明了 `NSMicrophoneUsageDescription`，应用就会终止。请为每个没有用途说明的功能明确设置 `PermissionDeny`：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionDeny,
},
```

@end

## 常见模式

### 媒体采集应用

在所有平台上授予摄像头和麦克风权限：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在<strong>Linux</strong>上，这会明确允许使用这两个设备；其他功能仍被拒绝。 在<strong>Windows</strong>上，这会允许使用这两个设备；任何未列出的其他功能都会显示原生提示。 在 **macOS** 上，这会在 WebKit 层允许两者，不显示浏览器提示。首次使用时，TCC 仍会询问设备本身的访问权限，且 `Info.plist` 必须包含两个用途说明键。

### 在 Linux 上拒绝媒体采集

Linux 默认允许使用摄像头和麦克风。要停用此行为：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### 在 Windows 上允许所有功能

要在不提示的情况下授予所有功能权限：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### 按窗口设置策略

不同窗口可以使用不同的策略：

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Windows 专用覆盖设置

按窗口设置的`Windows.Permissions`字段（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）仍然有效，并可在应用跨平台映射后覆盖各项功能的设置。当你需要访问没有跨平台对应项的 WebView2 权限类型（例如`CoreWebView2PermissionKindOtherSensors`）时，请使用此字段。

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows 上的求值顺序如下：

1. 跨平台`Permissions`映射（通过`SetPermission`设置各类型的状态）
2. `Windows.Permissions`映射（覆盖各个类型）
3. 对于两个映射均未涵盖的类型：如果配置了策略，则显示 WebView2 的原生提示；如果未配置策略，则自动允许（旧版行为）

## 平台支持矩阵

| 功能 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 麦克风 | ✅ | ✅ | ✅ (macOS 12+) |
| 摄像头 | ✅ | ✅ | ✅ (macOS 12+) |
| 地理位置 | ❌ 尚不支持 | ✅ | ❌ 尚不支持 |
| 通知 | ❌ 尚不支持 | ✅ | ❌ 尚不支持 |
| 读取剪贴板 | ❌ 尚不支持 | ✅ | ❌ 尚不支持 |

macOS 列为 ✅ 时，策略回答 WebKit 请求，TCC 仍在其上控制设备。为 ❌ 时，该功能仅交由 TCC 处理。

## 故障排除

**`getUserMedia`升级后在 Linux 上仍然失败**

请确认你没有显式设置`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。默认情况下（未设置），Linux 允许媒体采集。

**Windows 正在提示授予我未配置的权限**

一旦`Permissions`中存在任何条目，Wails 就不再设置统一的`Allow`授权。未列出的功能将显示 WebView2 的原生提示。请为应用使用的每项功能添加明确的`PermissionAllow`条目。

**macOS 权限不起作用**

`Permissions` 支持 macOS 12 及更高版本的摄像头和麦克风；地理位置、通知和剪贴板读取尚未接入，会忽略策略。TCC 仍控制设备：`PermissionAllow` 只移除 WebKit 提示，不移除系统提示。请检查 `Info.plist` 中的 `NSMicrophoneUsageDescription`、`NSCameraUsageDescription`，以及“系统设置 → 隐私与安全性”中的授权。

**网页内容请求摄像头或麦克风时，macOS 应用退出**

如果请求在缺少对应用途说明键的情况下到达 AVFoundation，结果不只是请求失败：macOS 会终止应用。 请添加对应的键，或为该功能设置 `PermissionDeny`，使请求无法到达 AVFoundation。`PermissionDefault` 会保留用户可以接受的 WebKit 提示。

**地理位置、通知和剪贴板策略在 Linux 上不起作用**

Linux 目前仅处理摄像头和麦克风权限。其他功能类型尚未实现支持——无论你设置何种策略，它们都仍会被拒绝。
