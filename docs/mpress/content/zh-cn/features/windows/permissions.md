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

### macOS（TCC）

macOS 通过其系统隐私框架管理摄像头、麦克风、地理位置和通知访问权限。当网页内容首次请求某项功能时，操作系统会自动显示提示，并在“系统设置”→“隐私与安全性”中按应用记住用户的选择。

无需进行任何`Permissions`配置，此机制即可正常工作。macOS<strong>目前会忽略此映射</strong>——无论你如何设置，所有请求都会交由 TCC 处理。实际缺口是`PermissionDeny`在 macOS 上不起作用：如果 TCC 已在系统层面授予某项功能的使用权限，你无法阻止 WebView 使用该功能。

请确保你的`Info.plist`包含相应的用途说明键：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## 常见模式

### 媒体采集应用

在所有平台上授予摄像头和麦克风权限：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在<strong>Linux</strong>上，这会明确允许使用这两个设备；其他功能仍被拒绝。 在<strong>Windows</strong>上，这会允许使用这两个设备；任何未列出的其他功能都会显示原生提示。 在<strong>macOS</strong>上，这项设置不起作用；一切均由 TCC 处理。

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
| 麦克风 | ✅ | ✅ | 仅由 TCC 处理 |
| 摄像头 | ✅ | ✅ | 仅由 TCC 处理 |
| 地理位置 | ❌ 尚不支持 | ✅ | 仅由 TCC 处理 |
| 通知 | ❌ 尚不支持 | ✅ | 仅由 TCC 处理 |
| 读取剪贴板 | ❌ 尚不支持 | ✅ | 仅由 TCC 处理 |

## 故障排除

**`getUserMedia`升级后在 Linux 上仍然失败**

请确认你没有显式设置`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。默认情况下（未设置），Linux 允许媒体采集。

**Windows 正在提示授予我未配置的权限**

一旦`Permissions`中存在任何条目，Wails 就不再设置统一的`Allow`授权。未列出的功能将显示 WebView2 的原生提示。请为应用使用的每项功能添加明确的`PermissionAllow`条目。

**macOS 权限不起作用**

`Permissions`映射在 macOS 上不起作用。请确保`Info.plist`包含正确的用途说明键（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription`等），并确保用户已在“系统设置 → 隐私与安全性”中授予访问权限。

**地理位置、通知和剪贴板策略在 Linux 上不起作用**

Linux 目前仅处理摄像头和麦克风权限。其他功能类型尚未实现支持——无论你设置何种策略，它们都仍会被拒绝。
