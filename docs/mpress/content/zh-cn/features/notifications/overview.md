---
title: "通知"
description: "显示带操作按钮和文本输入的原生系统通知"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## 简介

Wails 为桌面应用提供了一套全面的跨平台通知系统。该服务可显示原生系统通知，并支持：

- 包含标题、副标题和正文的基本通知
- 带操作按钮和文本回复的交互式通知
- 可复用于操作的[通知类别](#heading-6)
- 自定义[声音](#heading-10)（默认、静音或指定名称）
- [附件](#heading-11)（所有平台均支持图像；macOS 还支持音频和视频）
- 通过`ThreadID`对[相关通知进行分组](#heading-13)
- 通过`InterruptionLevel`设置[优先级](#heading-14)（`passive` / `active` / `timeSensitive` / `critical`）
- [定时送达](#heading-15)（macOS 使用原生机制；Windows 和 Linux 使用进程内定时器）
- 按 ID[更新已发出且仍有效的通知](#heading-16)

如果平台无法支持某个新增的可选字段，系统会平稳降级；有关各项功能的平台支持矩阵，请参阅[平台注意事项](#heading-17)。

## 基本用法

### 创建服务

首先，初始化通知服务：

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## 通知授权

macOS 上的通知需要用户授权。请求并检查授权：

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

在 Windows 和 Linux 上，此项始终返回`true`。

## 通知类型

### 基本通知

向用户发送一条基本通知，其中包含唯一 ID、标题、可选副标题（macOS 和 Linux）以及正文文本：

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### 交互式通知

发送带操作按钮和文本输入框的通知。必须先注册通知类别，才能使用此类通知：

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## 通知响应

处理用户与通知的交互：

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## 自定义通知

### 自定义元数据

基本通知和交互式通知可以包含自定义数据：

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### 自定义声音

使用`Sound`控制通知送达时播放的声音。将其保留为`nil`会播放平台默认声音；设为`Silent: true`可禁用声音；设为`Name`可播放指定名称或随应用捆绑的声音。

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

各平台解析`Name`的方式如下：

- **macOS** — 将`Name`传递给`[UNNotificationSound soundNamed:]`；音频文件必须位于应用程序包的`Library/Sounds`目录下。
- **Windows** — 如果`Name`已以`ms-winsoundevent:`或`ms-appx:`开头，则直接使用；否则将其包装为`ms-winsoundevent:`，作为内置 toast 事件名称使用（请参阅 Microsoft 的 toast `<audio>`架构文档）。
- **Linux** — 作为 freedesktop 的`sound-name`提示转发；是否播放取决于当前使用的通知守护进程和声音主题。

### 附件

`Attachments`可为通知添加媒体文件。macOS 支持任意媒体类型的多个附件；Windows 和 Linux 仅处理第一个图像类型的附件。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path`必须是绝对文件系统路径。macOS 还接受`file://`URL。

#### 附加随应用分发的文件

操作系统在通知送达时从磁盘读取附件，因此`Path`必须能解析为最终用户计算机上的实际文件。对于随应用捆绑的资源（使用`go:embed`嵌入的图标或图像），无法硬编码一个固定的绝对路径，因为该文件位于二进制文件内部，而不是磁盘上的已知位置。请在启动时将其写入可写目录一次，然后传入该路径：

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

用户提供或下载的文件在磁盘上已有实际路径，因此无需执行此步骤，可直接将其传给`Path`。未来版本可能会增加直接传递内存中附件字节的功能。

### 通知串联与分组

`ThreadID`将相关通知归为一组，使操作系统可以在通知中心 / 操作中心中折叠显示。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### 打扰级别

`InterruptionLevel`控制通知优先级。请使用以下任一导出常量：

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | 静默送达；不会点亮屏幕，也不会播放默认声音 |
| `InterruptionLevelActive` | `"active"` | 默认级别 |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | 在允许的情况下突破专注模式 / 勿扰模式 |
| `InterruptionLevelCritical` | `"critical"` | 绕过专注模式和响铃设置；macOS 要求具备“严重警报”权限（缺少该权限时会静默降级） |

各平台的映射方式如下：

- **macOS** — 设置`UNNotificationContent.interruptionLevel`。`critical`要求 macOS 12+ 及“严重警报”权限。
- **Windows** — 映射到 toast 的`<toast scenario="...">`属性。
- **Linux** — 映射到 freedesktop 的`urgency`提示。

### 定时送达

`Schedule`用于推迟送达。请仅设置`DelaySeconds`（从现在起的秒数）或`At`（Unix 秒数，UTC）中的一个。

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="持久性"}
在**macOS**上，定时通知使用原生触发器，并可在应用重启后继续保留。在**Windows**和**Linux**上，定时功能会回退到进程内`time.AfterFunc`定时器；如果应用在送达前退出，通知将**丢失**，因为`wintoast`和 freedesktop 规范均未提供延迟送达原语。

@end

### 更新通知

`UpdateNotification`会替换具有相同`ID`且已发出并仍有效的通知：

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

各平台的行为如下：

- **macOS** — `UNUserNotificationCenter`会按标识符自动去重，因此现有通知将在原位置更新。
- **Linux** — 使用 D-Bus 的`replaces_id`参数替换上一条通知。
- **Windows** — 目前会作为新通知重新投递。真正的原位替换需要上游`wintoast`支持`tag`/`group`。

## 平台注意事项

@tabs
[macOS]
在 macOS 上，通知：

- 需要用户授权
- 要求应用已打包并签名（分发时还需经过公证）
- 使用系统标准的通知外观
- 支持`Subtitle`
- 支持用户输入文本（回复）
- 支持`Destructive`操作选项
- 支持多个任意媒体类型（图像、音频、视频）的`Attachments`
- 支持使用`ThreadID`在通知中心中分组
- 支持所有`InterruptionLevel`值（`critical`需要“严重警报”授权）
- 支持应用重启后仍然有效的原生计划投递
- 自动按`ID`对`UpdateNotification`调用去重
- 自动适配深色/浅色模式

[Windows]
在 Windows 上，通知：

- 通过`wintoast`后端使用 Windows 系统 Toast 样式
- 适配 Windows 主题设置
- 支持用户输入文本（回复）
- 支持高 DPI 显示器
- 不支持`Subtitle`
- 支持单个图像`Attachment`，其位置提示可设为`hero`、`appLogoOverride`或`inline`（默认为`inline`）
- 支持使用`ThreadID`在操作中心中分组
- 通过 Toast 的`scenario`属性支持`InterruptionLevel`
- 通过进程内计时器支持计划投递——如果应用在投递前退出，**计划通知将会丢失**
- `UpdateNotification`目前会作为新通知重新投递（真正的原位替换仍需等待上游`wintoast`支持`tag`/`group`）

[Linux]
在 Linux 上，通知使用 D-Bus 的`org.freedesktop.Notifications`接口。要使通知正常工作，**必须有兼容的通知守护进程正在运行**。

@note{type="caution" title="系统要求：通知守护进程"}
必须安装并运行兼容 freedesktop 的通知守护进程。常见选择包括：

- **dunst** — 轻量且高度可配置（`apt install dunst`/`dnf install dunst`）
- **mako** — 原生支持 Wayland（`apt install mako-notifier`）
- **GNOME Shell** — 在 GNOME 43+ 上会自动注册该接口。在 Ubuntu 24.04（GNOME Shell 46）上，该接口在会话启动时可能不会自行注册；如果通知没有出现，请安装`dunst`作为后备方案。
- **xfce4-notifyd** — 随 XFCE 桌面一同提供

如果没有守护进程在运行，`SendNotification`将返回 D-Bus 错误：`The name org.freedesktop.Notifications was not provided by any .service files`。请在应用中处理此错误，并提示用户安装通知守护进程。

@end

在 Linux 上，通知：

- 遵循桌面环境主题
- 按照桌面环境规则定位
- 支持`Subtitle`（对于不单独呈现它的守护进程，会将其拼接到正文中）
- 不支持用户输入文本（不属于 freedesktop 规范）
- 通过`image-path`提示支持单个图像`Attachment`
- 支持`ThreadID`（如果守护进程支持，则由其处理）
- `Sound.Name`会作为`sound-name`提示转发；是否播放取决于当前使用的守护进程和声音主题
- 将`InterruptionLevel`映射到 freedesktop 的`urgency`提示
- 通过进程内计时器支持计划投递——如果应用在投递前退出，**计划通知将会丢失**
- `UpdateNotification`使用 D-Bus 的`replaces_id`参数原位替换上一条通知

@end

## 最佳实践

1. 检查并请求授权：
  - macOS 需要用户授权


2. 提供清晰简洁的通知：
  - 使用描述清晰的标题、副标题、正文和操作标题


3. 妥善处理通知响应：
  - 检查通知响应中的错误
  - 为用户操作提供反馈


4. 考虑平台惯例：
  - 遵循特定于平台的通知模式
  - 尊重系统设置


5. 在 Linux 上，处理守护进程依赖：
  - 检查`SendNotification`返回的错误——缺少守护进程会产生 D-Bus 错误
  - 软件包文档或应用 README 应说明需要 freedesktop 通知守护进程


## 示例

查看此示例：

- [通知](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notifications)

## API 参考

### 服务管理

| 方法 | 说明 |
| --- | --- |
| `New()` | 创建新的通知服务 |

### 通知授权

| 方法 | 说明 |
| --- | --- |
| `RequestNotificationAuthorization()` | 请求显示通知的权限（macOS） |
| `CheckNotificationAuthorization()` | 检查当前通知授权状态（macOS） |

### 发送通知

| 方法 | 说明 |
| --- | --- |
| `SendNotification(options NotificationOptions)` | 发送基本通知 |
| `SendNotificationWithActions(options NotificationOptions)` | 发送带操作的交互式通知 |
| `UpdateNotification(options NotificationOptions)` | 按`ID`更新正在处理的通知（请参阅[更新通知](#heading-16)） |

### 通知类别

| 方法 | 说明 |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | 注册可复用的通知类别 |
| `RemoveNotificationCategory(categoryID string)` | 移除之前注册的类别 |

### 管理通知

| 方法 | 说明 |
| --- | --- |
| `RemoveAllPendingNotifications()` | 移除所有待处理通知（仅限 macOS 和 Linux） |
| `RemovePendingNotification(identifier string)` | 移除指定的待处理通知（仅限 macOS 和 Linux） |
| `RemoveAllDeliveredNotifications()` | 移除所有已送达的通知（仅限 macOS 和 Linux） |
| `RemoveDeliveredNotification(identifier string)` | 移除指定的已送达通知（仅限 macOS 和 Linux） |
| `RemoveNotification(identifier string)` | 移除通知（仅适用于 Linux） |

### 事件处理

| 方法 | 说明 |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | 注册用于处理通知响应的回调函数 |

### 结构体和类型

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### InterruptionLevel 常量

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
