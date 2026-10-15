---
title: "通知"
description: "顯示含動作按鈕與文字輸入欄位的原生系統通知"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## 簡介

Wails 為桌面應用程式提供完整的跨平台通知系統。此服務可讓您顯示原生系統通知，並支援：

- 含標題、副標題及本文的基本通知
- 含動作按鈕及文字回覆的互動式通知
- 可重複用於動作的[通知類別](#heading-6)
- 自訂[音效](#heading-10)（預設、靜音或指定名稱）
- [附件](#heading-11)（所有平台皆支援圖片；macOS 支援音訊與視訊）
- 依據`ThreadID`將[相關通知分組](#heading-13)
- 透過`InterruptionLevel`設定[優先順序](#heading-14)（`passive`／`active`／`timeSensitive`／`critical`）
- [排程傳送](#heading-15)（macOS 使用原生機制；Windows 與 Linux 使用處理程序內計時器）
- 依 ID[更新傳送中的通知](#heading-16)

若平台無法支援新加入的選用欄位，該欄位皆會平順降級；各功能的支援對照表請參閱[平台注意事項](#heading-17)。

## 基本用法

### 建立服務

首先，初始化通知服務：

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

## 通知授權

macOS 上的通知需要使用者授權。請要求並檢查授權：

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

在 Windows 和 Linux 上，此項目一律傳回`true`。

## 通知類型

### 基本通知

向使用者傳送基本通知，其中包含唯一 ID、標題、選用的副標題（macOS 和 Linux）及本文文字：

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### 互動式通知

傳送含動作按鈕及文字輸入欄位的通知。這些通知需要先註冊通知類別：

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

## 通知回應

處理使用者與通知的互動：

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

## 自訂通知

### 自訂中繼資料

基本通知和互動式通知可包含自訂資料：

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

### 自訂音效

使用`Sound`控制傳送通知時播放的音訊。保留為`nil`會播放平台預設音效；設為`Silent: true`會停用音效；設定`Name`欄位則會播放指定名稱或隨附的音效。

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

- **macOS** — 將`Name`傳給`[UNNotificationSound soundNamed:]`；音訊檔案必須位於應用程式套件的`Library/Sounds`下。
- **Windows** — 若`Name`已以`ms-winsoundevent:`或`ms-appx:`開頭，便會直接使用；否則會以`ms-winsoundevent:`包裝，作為內建的快顯通知事件名稱使用（請參閱 Microsoft 的快顯通知`<audio>`結構描述文件）。
- **Linux** — 以 freedesktop `sound-name`提示轉送；是否播放取決於目前使用的通知常駐程式及音效佈景主題。

### 附件

`Attachments`會在通知中加入媒體檔案。macOS 支援多個任意媒體類型的附件；Windows 和 Linux 僅支援第一個圖片類型的附件。

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

`Path`必須是絕對檔案系統路徑。macOS 另接受`file://`URL。

#### 附加應用程式隨附的檔案

作業系統會在傳送通知時從磁碟讀取附件，因此`Path`必須解析為終端使用者電腦上實際存在的檔案。對於隨應用程式封裝的資產（使用`go:embed`嵌入的圖示或圖片），因為檔案位於二進位檔內，而非磁碟上的已知位置，所以沒有可寫死的固定絕對路徑。請在啟動時將它寫入可寫入的目錄一次，然後傳入該路徑：

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

使用者提供或下載的檔案已在磁碟上具有實際路徑，因此不必執行此步驟，便可直接將路徑傳給`Path`。未來版本可能會新增直接傳入記憶體內附件位元組的功能。

### 串接與分組

`ThreadID`會將相關通知分組，讓作業系統可在「通知中心」／「行動作業中心」中將它們收合。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### 中斷層級

`InterruptionLevel`控制通知的優先順序。請使用下列其中一個匯出的常數：

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| 常數 | 值 | 含義 |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | 靜默傳送；不會點亮螢幕或播放預設音效 |
| `InterruptionLevelActive` | `"active"` | 預設層級 |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | 在允許的情況下突破「專注模式」／「勿擾模式」 |
| `InterruptionLevelCritical` | `"critical"` | 略過「專注模式」和鈴聲設定；macOS 需要「重大警示」權限（若無此權限，會無提示地降級） |

各平台的對應方式：

- **macOS** — 設定`UNNotificationContent.interruptionLevel`。`critical`需要 macOS 12+ 和「重大警示」權限。
- **Windows** — 對應至快顯通知的`<toast scenario="...">`屬性。
- **Linux** — 對應至 freedesktop 的`urgency`提示。

### 排程傳送

`Schedule`會延後傳送。請只設定`DelaySeconds`（從現在起算的秒數）或`At`（Unix 秒數，UTC）其中之一。

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

@note{type="caution" title="持續性"}
在<strong>macOS</strong>上，排程通知使用原生觸發程序，並可在應用程式重新啟動後保留。在<strong>Windows</strong>和<strong>Linux</strong>上，排程會改用行程內的`time.AfterFunc`計時器；若應用程式在傳送前結束，排程通知將會<strong>遺失</strong>，因為`wintoast`和 freedesktop 規格都未提供延後傳送的基本機制。

@end

### 更新通知

`UpdateNotification`會取代具有相同`ID`且仍有效的既有通知：

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

各平台的行為：

- **macOS** — `UNUserNotificationCenter`會依識別碼自動去除重複項目，因此既有通知會直接原地更新。
- **Linux** — 使用 D-Bus 的`replaces_id`參數取代先前的通知。
- **Windows** — 目前會以新通知重新傳送。若要真正原地取代，需等待上游`wintoast`支援`tag`／`group`。

## 平台注意事項

@tabs
[macOS]
在 macOS 上，通知：

- 需要使用者授權
- 需要將應用程式封裝並簽署（散發時還須經過公證）
- 採用系統標準的通知外觀
- 支援`Subtitle`
- 支援使用者文字輸入（回覆）
- 支援`Destructive`動作選項
- 支援任意媒體類型（圖片、音訊、影片）的多個`Attachments`
- 支援使用`ThreadID`在「通知中心」中分組
- 支援所有`InterruptionLevel`值（`critical`需要「重大警示」權限）
- 支援可在應用程式重新啟動後保留的原生排程傳送
- 依`ID`自動去除重複的`UpdateNotification`呼叫
- 自動處理深色／淺色模式

[Windows]
在 Windows 上，通知：

- 透過`wintoast`後端使用 Windows 系統快顯通知樣式
- 配合 Windows 佈景主題設定
- 支援使用者文字輸入（回覆）
- 支援高 DPI 顯示器
- 不支援`Subtitle`
- 支援單一圖片`Attachment`，其位置提示可設為`hero`、`appLogoOverride`或`inline`（預設為`inline`）
- 支援使用`ThreadID`在「重要訊息中心」中分組
- 透過快顯通知的`scenario`屬性支援`InterruptionLevel`
- 透過行程內計時器支援排程傳送；若應用程式在傳送前結束，**排程通知將會遺失**
- `UpdateNotification`目前會以新通知重新傳送（若要真正原地取代，仍需等待上游`wintoast`支援`tag`／`group`）

[Linux]
在 Linux 上，通知使用 D-Bus 的`org.freedesktop.Notifications`介面。系統<strong>必須執行相容的通知常駐程式</strong>，通知才能運作。

@note{type="caution" title="系統需求：通知常駐程式"}
必須安裝並執行與 freedesktop 相容的通知常駐程式。常見選擇包括：

- **dunst** — 輕量且可高度自訂（`apt install dunst`／`dnf install dunst`）
- **mako** — Wayland 原生（`apt install mako-notifier`）
- **GNOME Shell** — 在 GNOME 43+ 上會自動註冊此介面。在 Ubuntu 24.04（GNOME Shell 46）上，此介面可能不會在工作階段開始時自行註冊；若未顯示通知，請安裝`dunst`作為備用方案。
- **xfce4-notifyd** — 隨 XFCE 桌面環境提供

若沒有執行中的常駐程式，`SendNotification`會傳回 D-Bus 錯誤：`The name org.freedesktop.Notifications was not provided by any .service files`。請在應用程式中處理此錯誤，並通知使用者安裝通知常駐程式。

@end

在 Linux 上，通知：

- 採用桌面環境的佈景主題
- 依桌面環境規則決定位置
- 支援`Subtitle`（對於不會單獨呈現該內容的常駐程式，會將其串接至本文）
- 不支援使用者文字輸入（不屬於 freedesktop 規格的一部分）
- 透過`image-path`提示支援單一圖片`Attachment`
- 支援`ThreadID`（在支援此功能的常駐程式中處理）
- `Sound.Name`會以`sound-name`提示轉送；是否播放取決於目前使用的常駐程式與音效主題
- 將`InterruptionLevel`對應至 freedesktop 的`urgency`提示
- 支援透過行程內計時器排程傳送；如果應用程式在傳送前結束，**已排程的通知將會遺失**
- `UpdateNotification`使用 D-Bus 的`replaces_id`參數，就地取代先前的通知

@end

## 最佳實務

1. 檢查並要求授權：
  - macOS 需要使用者授權


2. 提供清楚簡潔的通知：
  - 使用具描述性的標題、副標題、文字與動作標題


3. 妥善處理通知回應：
  - 檢查通知回應中的錯誤
  - 針對使用者動作提供回饋


4. 考量平台慣例：
  - 遵循各平台特有的通知模式
  - 遵循系統設定


5. 在 Linux 上處理對常駐程式的相依性：
  - 檢查`SendNotification`傳回的錯誤；缺少常駐程式時會產生 D-Bus 錯誤
  - 套件文件或應用程式 README 應註明需要 freedesktop 通知常駐程式


## 範例

瀏覽此範例：

- [通知](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## API 參考

### 服務管理

| 方法 | 說明 |
| --- | --- |
| `New()` | 建立新的通知服務 |

### 通知授權

| 方法 | 說明 |
| --- | --- |
| `RequestNotificationAuthorization()` | 要求顯示通知的權限（macOS） |
| `CheckNotificationAuthorization()` | 檢查目前的通知授權狀態（macOS） |

### 傳送通知

| 方法 | 說明 |
| --- | --- |
| `SendNotification(options NotificationOptions)` | 傳送基本通知 |
| `SendNotificationWithActions(options NotificationOptions)` | 傳送包含動作的互動式通知 |
| `UpdateNotification(options NotificationOptions)` | 依`ID`更新仍有效的既有通知（請參閱[更新通知](#heading-16)） |

### 通知類別

| 方法 | 說明 |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | 註冊可重複使用的通知類別 |
| `RemoveNotificationCategory(categoryID string)` | 移除先前註冊的類別 |

### 管理通知

| 方法 | 說明 |
| --- | --- |
| `RemoveAllPendingNotifications()` | 移除所有待傳送的通知（僅限 macOS 與 Linux） |
| `RemovePendingNotification(identifier string)` | 移除指定的待傳送通知（僅限 macOS 與 Linux） |
| `RemoveAllDeliveredNotifications()` | 移除所有已傳送的通知（僅限 macOS 與 Linux） |
| `RemoveDeliveredNotification(identifier string)` | 移除指定的已傳送通知（僅限 macOS 與 Linux） |
| `RemoveNotification(identifier string)` | 移除通知（Linux 特有） |

### 事件處理

| 方法 | 說明 |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | 註冊處理通知回應的回呼函式 |

### 結構與型別

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

#### InterruptionLevel 常數

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
