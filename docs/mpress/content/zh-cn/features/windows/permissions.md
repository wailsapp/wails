---
title: "权限"
description: "控制网页内容对摄像头、麦克风、地理位置及其他功能的请求"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

调用`navigator.mediaDevices.getUserMedia()`、地理位置 API 或通知 API 的网页内容，需要由宿主应用授予或拒绝这些请求。Wails 在`WebviewWindowOptions`上提供跨平台的`Permissions`映射，让你能以声明方式进行控制，无需编写平台特定代码。

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

来自该窗口网页内容的摄像头和麦克风请求将被授予，且不会显示浏览器提示。

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

`Permission`（uint8）是应用于给定类型的策略。

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 使用平台的原生处理方式（见下文） |
| `PermissionAllow` | 1 | 不提示，直接授予 |
| `PermissionDeny` | 2 | 不提示，直接拒绝 |

`PermissionDefault`是零值，因此映射中未设置的条目按默认策略处理。

## 平台行为

由于底层 WebView 的原生行为不同，各平台对`PermissionDefault`的处理方式也不同。

### Linux（WebKitGTK）

WebKitGTK**没有原生权限提示**。如果未附加处理程序，它会静默拒绝所有请求；这就是为什么在添加此功能之前，`getUserMedia`始终返回`NotAllowedError`。

目前，Wails 在 Linux 上处理**摄像头和麦克风**请求。地理位置、通知和剪贴板读取功能尚未接入，因此无论设置何种策略，这些请求仍会被拒绝。

| 策略 | 摄像头/麦克风 | 地理位置、通知、剪贴板 |
| --- | --- | --- |
| `PermissionDefault` | **允许**（恢复 getUserMedia） | 始终拒绝 |
| `PermissionAllow` | 允许 | 始终拒绝（尚未实现） |
| `PermissionDeny` | 拒绝 | 始终拒绝 |

### Windows（WebView2）

WebView2 提供原生权限提示和按类型设置权限的 API。五种功能类型均受完整支持。

| 策略 | 行为 |
| --- | --- |
| `PermissionDefault` | WebView2 显示其原生操作系统/浏览器权限提示 |
| `PermissionAllow` | 静默授予 |
| `PermissionDeny` | 静默拒绝 |

**重要：**在此功能推出之前，Wails 会无条件调用`SetGlobalPermission(Allow)`，从而静默授予所有功能。现在，只要`Permissions`中存在任何条目，就**不会**设置这种一概授予。未设置的功能将转由 WebView2 显示原生提示，而不会自动授予。

这意味着，只要在 Windows 上配置了`Permissions`，任何未明确列出的功能都会显示提示，而不是被静默允许。请明确设置所需的功能。

### macOS（WKWebView + TCC）

在 macOS 上，这两个层级都必须同意。WKWebView 在启动捕获会话前会询问应用，而`Permissions`映射负责响应该请求。在它之下，系统隐私框架（TCC）控制设备本身：应用首次真正访问摄像头或麦克风时会显示操作系统提示，系统会按应用记住用户的选择，相关设置位于“系统设置”→“隐私与安全性”。

Wails 在 macOS 12及更高版本上处理**摄像头和麦克风**请求。地理位置、通知和剪贴板读取没有对应的`WKUIDelegate`，因此尚未接入，仅由 TCC 处理；与 Linux 一样，为这些功能设置的策略不会生效。

| 策略 | 摄像头/麦克风 | 地理位置、通知、剪贴板 |
| --- | --- | --- |
| `PermissionDefault` | WebKit 显示自身的权限提示 | 仅由 TCC 处理 |
| `PermissionAllow` | 跳过 WebKit 的提示，但**TCC 仍然适用** | 仅由 TCC 处理 |
| `PermissionDeny` | 由 WebKit 拒绝；系统授权单独处理 | 仅由 TCC 处理 |

`PermissionAllow`授予的是 WebView 的请求，而不是设备访问权限。系统授权仍独立生效；如果应用在“系统设置”中被拒绝，它仍会被拒绝。只有 WebKit 自身的提示会被跳过。

在低于 macOS 12的版本中，该委托方法不存在，因此映射会被忽略，所有请求都回退到 WebKit 的提示。

请确保`Info.plist`包含适当的用途说明键：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="必须提供用途说明"}

WebKit 可能会在调用 Wails 权限委托之前请求系统授权。`PermissionDeny`会阻止 WebKit 请求，但不能保证不会出现 TCC 提示或 AVFoundation 授权请求。

请为应用所请求的设备添加适当的用途说明；如果应用在沙盒中运行，还需添加匹配的设备权利声明。不要用`PermissionDeny`替代这些声明。

@end

## 常见模式

### 媒体捕获应用

在所有平台上授予摄像头和麦克风权限：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在**Linux**上，此设置会明确允许使用这两种设备；其他能力仍会被拒绝。在**Windows**上，此设置会允许这两种设备；任何未列出的其他能力都会触发原生提示。在**macOS**上，此设置会在 WebKit 层允许这两种设备，因此不会出现浏览器提示；首次使用时，TCC 仍会请求设备访问权限，并且`Info.plist`中必须包含这两项用途说明键。

### 在 Linux 上拒绝媒体采集

Linux 默认允许使用摄像头和麦克风。要选择停用此功能：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### 在 Windows 上允许所有能力

要在不提示的情况下授予所有能力：

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

按窗口设置的`Windows.Permissions`字段（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）仍然有效，并可在应用跨平台映射后覆盖各项能力。需要访问没有跨平台对应项的 WebView2 权限类型（例如`CoreWebView2PermissionKindOtherSensors`）时，请使用此字段。

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
3. 对于两种映射均未涵盖的任何类型：已配置策略时显示 WebView2 原生提示；未配置策略时自动允许（旧版行为）

## 平台支持矩阵

| 能力 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 麦克风 | ✅ | ✅ | ✅（macOS 12 及更高版本） |
| 摄像头 | ✅ | ✅ | ✅（macOS 12 及更高版本） |
| 地理位置 | ❌ 尚未支持 | ✅ | ❌ 尚未支持 |
| 通知 | ❌ 尚未支持 | ✅ | ❌ 尚未支持 |
| 读取剪贴板 | ❌ 尚未支持 | ✅ | ❌ 尚未支持 |

在 macOS 一栏标为 ✅ 时，策略会响应 WebKit 的请求；除此之外，设备访问仍受 TCC 限制。标为 ❌ 时，该能力仅由 TCC 决定。

## 故障排除

**`getUserMedia`升级后在 Linux 上仍然失败**

请检查是否明确设置了`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。默认值（未设置）允许在 Linux 上采集媒体。

**Windows 正在提示我授予未配置的权限**

只要`Permissions`中存在任何条目，Wails 就不再设置全局`Allow`授权。未列出的能力会触发 WebView2 原生提示。请为应用使用的每项能力添加明确的`PermissionAllow`条目。

**macOS 权限不起作用**

在 macOS 12及更高版本上，`Permissions`涵盖摄像头和麦克风；地理位置、通知和剪贴板读取功能尚未接入，因此会忽略该策略。在策略生效的情况下，设备访问还会受到 TCC 限制：`PermissionAllow`会消除 WebKit 的提示，但不会消除系统提示。请确保`Info.plist`中包含正确的用途说明键（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription`），并确认用户已在“系统设置”→“隐私与安全性”中授予访问权限。

**网页内容请求使用摄像头或麦克风时，我的 macOS 应用会退出**

请检查`Info.plist`，确保针对应用请求的每种设备包含`NSCameraUsageDescription`或`NSMicrophoneUsageDescription`。缺少用途说明可能会导致请求被拒绝；在缺少必需用途说明的情况下访问 AVFoundation，可能会导致应用终止。`PermissionDeny`控制 WebKit 的决定，不能替代这些键。

**地理位置、通知和剪贴板设置在 Linux 上不起作用**

Linux 目前仅处理摄像头和麦克风。尚未实现对其他能力类型的支持——无论如何设置策略，它们都会保持拒绝状态。
