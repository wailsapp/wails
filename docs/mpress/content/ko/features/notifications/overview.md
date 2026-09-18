---
title: "알림"
description: "작업 버튼과 텍스트 입력을 지원하는 네이티브 시스템 알림 표시"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## 소개

Wails는 데스크톱 애플리케이션을 위한 포괄적인 크로스 플랫폼 알림 시스템을 제공합니다. 이 서비스를 사용하면 다음 기능을 지원하는 네이티브 시스템 알림을 표시할 수 있습니다.

- 제목, 부제목 및 본문이 포함된 기본 알림
- 작업 버튼과 텍스트 답장을 지원하는 대화형 알림
- 작업에 재사용할 수 있는 [알림 카테고리](#--5)
- 사용자 지정 [소리](#---2)(기본, 무음 또는 이름 지정)
- [첨부 파일](#--7)(모든 플랫폼의 이미지, macOS의 오디오/동영상)
- `ThreadID`을 기준으로 [관련 알림 그룹화](#---3)
- `InterruptionLevel`을 통한 [우선순위](#--8) 지정(`passive` / `active` / `timeSensitive` / `critical`)
- [예약 전송](#--9)(macOS에서는 네이티브 방식, Windows 및 Linux에서는 프로세스 내 타이머 사용)
- ID를 사용하여 [전송 중인 알림 업데이트](#--10)

새로 추가된 각 선택적 필드를 플랫폼에서 지원할 수 없는 경우에도 기능은 문제없이 축소 적용됩니다. 기능별 지원 매트릭스는 [플랫폼 고려 사항](#---4)을 참조하세요.

## 기본 사용법

### 서비스 생성

먼저 알림 서비스를 초기화하세요.

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

## 알림 권한

macOS에서 알림을 사용하려면 사용자의 권한이 필요합니다. 권한을 요청하고 확인하세요.

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

Windows와 Linux에서는 항상 `true`을 반환합니다.

## 알림 유형

### 기본 알림

고유 ID, 제목, 선택적 부제목(macOS 및 Linux), 본문 텍스트가 포함된 기본 알림을 사용자에게 보내세요.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### 대화형 알림

작업 버튼과 텍스트 입력이 포함된 알림을 보내세요. 이러한 알림을 사용하려면 먼저 알림 카테고리를 등록해야 합니다.

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

## 알림 응답

알림에 대한 사용자 상호작용을 처리하세요.

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

## 알림 사용자 지정

### 사용자 지정 메타데이터

기본 알림과 대화형 알림에 사용자 지정 데이터를 포함할 수 있습니다.

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

### 사용자 지정 소리

알림이 전달될 때 재생되는 오디오는 `Sound`으로 제어하세요. `nil`으로 두면 플랫폼 기본 소리가 재생됩니다. 소리를 끄려면 `Silent: true`으로 설정하고, 이름이 지정되었거나 번들에 포함된 소리를 재생하려면 `Name` 필드를 설정하세요.

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

플랫폼별 `Name` 해석 방식은 다음과 같습니다.

- **macOS** — `Name`이 `[UNNotificationSound soundNamed:]`에 전달됩니다. 오디오 파일은 앱 번들의 `Library/Sounds` 아래에 있어야 합니다.
- **Windows** — `Name`이 이미 `ms-winsoundevent:` 또는 `ms-appx:`으로 시작하면 그대로 사용합니다. 그렇지 않으면 기본 제공 토스트 이벤트 이름으로 사용할 수 있도록 `ms-winsoundevent:`으로 감쌉니다(Microsoft의 토스트 `<audio>` 스키마 문서 참조).
- **Linux** — freedesktop의 `sound-name` 힌트로 전달됩니다. 재생 여부는 활성 알림 데몬과 사운드 테마에 따라 달라집니다.

### 첨부 파일

`Attachments`은 알림과 함께 미디어 파일을 추가합니다. macOS는 모든 미디어 유형의 첨부 파일을 여러 개 지원하며, Windows와 Linux는 이미지 유형의 첫 번째 첨부 파일만 처리합니다.

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

`Path`은 절대 파일 시스템 경로여야 합니다. macOS에서는 `file://` URL도 사용할 수 있습니다.

#### 앱과 함께 제공되는 파일 첨부

알림이 전송될 때 OS가 디스크에서 첨부 파일을 읽으므로 `Path`은 최종 사용자의 컴퓨터에 실제로 존재하는 파일로 해석되어야 합니다. 애플리케이션에 번들로 포함한 에셋(`go:embed`을 사용하여 임베드한 아이콘이나 이미지)은 알려진 디스크 위치가 아니라 바이너리 내부에 있으므로, 하드 코딩할 수 있는 고정 절대 경로가 없습니다. 시작할 때 쓰기 가능한 디렉터리에 한 번 저장한 후 해당 경로를 전달하세요.

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

사용자가 제공했거나 다운로드한 파일은 이미 디스크의 실제 경로를 가지므로 이 단계 없이 `Path`에 직접 전달할 수 있습니다. 메모리에 있는 첨부 파일 바이트를 전달하는 기능은 향후 릴리스에 추가될 수 있습니다.

### 스레딩 및 그룹화

`ThreadID`은 관련 알림을 함께 그룹화하여 OS가 알림 센터/관리 센터에서 이를 접을 수 있도록 합니다.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### 방해 수준

`InterruptionLevel`은 알림 우선순위를 제어합니다. 내보낸 상수 중 하나를 사용하세요.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| 상수 | 값 | 의미 |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | 조용히 전송되며 화면을 켜거나 기본 소리를 재생하지 않음 |
| `InterruptionLevelActive` | `"active"` | 기본 수준 |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | 허용되는 경우 집중 모드/방해 금지를 무시함 |
| `InterruptionLevelCritical` | `"critical"` | 집중 모드와 벨소리 설정을 우회합니다. macOS에서는 Critical Alert 권한이 필요하며, 이 권한이 없으면 알림 수준이 별도 안내 없이 낮아집니다. |

플랫폼별 매핑:

- **macOS** — `UNNotificationContent.interruptionLevel`을 설정합니다. `critical`에는 macOS 12 이상과 Critical Alert 권한이 필요합니다.
- **Windows** — 토스트의 `<toast scenario="...">` 속성에 매핑합니다.
- **Linux** — freedesktop의 `urgency` 힌트에 매핑합니다.

### 예약 전송

`Schedule`은 전송을 지연합니다. `DelaySeconds`(현재부터의 초 단위 시간) 또는 `At`(UTC 기준 Unix 초) 중 정확히 하나만 설정하세요.

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

@note{type="caution" title="지속성"}
<strong>macOS</strong>에서는 예약 알림이 네이티브 트리거를 사용하므로 앱을 다시 시작해도 유지됩니다. <strong>Windows</strong>와 <strong>Linux</strong>에서는 프로세스 내 `time.AfterFunc` 타이머를 대신 사용하므로, **전송 전에 앱이 종료되면 예약이 소실됩니다**. `wintoast`과 freedesktop 명세 모두 지연 전송 기능을 제공하지 않기 때문입니다.

@end

### 알림 업데이트

`UpdateNotification`은 동일한 `ID`을 가진 기존 알림을 교체합니다:

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

플랫폼별 동작:

- **macOS** — `UNUserNotificationCenter`이 식별자를 기준으로 자동 중복 제거하므로 기존 알림이 그 자리에서 업데이트됩니다.
- **Linux** — 이전 알림을 교체하는 데 D-Bus의 `replaces_id` 매개변수를 사용합니다.
- **Windows** — 현재는 새 알림으로 다시 전송합니다. 실제 제자리 교체를 구현하려면 업스트림 `wintoast`에서 `tag` / `group`을 지원해야 합니다.

## 플랫폼별 고려 사항

@tabs
[macOS]
macOS의 알림은 다음과 같습니다:

- 사용자 승인이 필요합니다
- 앱을 패키징하고 서명해야 합니다(배포하려면 공증도 필요합니다)
- 시스템 표준 알림 모양을 사용합니다
- `Subtitle`을 지원합니다
- 사용자 텍스트 입력(답장)을 지원합니다
- `Destructive` 작업 옵션을 지원합니다
- 모든 미디어 유형(이미지, 오디오, 동영상)의 `Attachments`을 여러 개 지원합니다
- 알림 센터에서 그룹화할 수 있도록 `ThreadID`을 지원합니다
- 모든 `InterruptionLevel` 값을 지원합니다(`critical`에는 Critical Alert 권한이 필요합니다)
- 앱을 다시 시작해도 유지되는 네이티브 예약 전송을 지원합니다
- `UpdateNotification` 호출의 중복을 `ID`을 기준으로 자동 제거합니다
- 다크 모드와 라이트 모드를 자동으로 처리합니다

[Windows]
Windows의 알림은 다음과 같습니다:

- `wintoast` 백엔드를 통해 Windows 시스템 토스트 스타일을 사용합니다
- Windows 테마 설정에 맞게 조정됩니다
- 사용자 텍스트 입력(답장)을 지원합니다
- 고DPI 디스플레이를 지원합니다
- `Subtitle`을 지원하지 않습니다
- 배치 힌트 `hero`, `appLogoOverride` 또는 `inline`(기본값: `inline`)을 사용하는 이미지 `Attachment` 하나를 지원합니다
- 알림 센터에서 그룹화할 수 있도록 `ThreadID`을 지원합니다
- 토스트의 `scenario` 속성을 통해 `InterruptionLevel`을 지원합니다
- 프로세스 내 타이머를 통한 예약 전송을 지원하지만, **전송 전에 앱이 종료되면 예약 알림이 소실됩니다**
- `UpdateNotification`은 현재 새 알림으로 다시 전송합니다(실제 제자리 교체는 업스트림 `wintoast`의 `tag`/`group` 지원을 기다리고 있습니다)

[Linux]
Linux에서 알림은 D-Bus `org.freedesktop.Notifications` 인터페이스를 사용합니다. 알림이 작동하려면 호환되는 알림 데몬이 **반드시 실행 중이어야 합니다**.

@note{type="caution" title="시스템 요구 사항: 알림 데몬"}
freedesktop 호환 알림 데몬이 설치되어 실행 중이어야 합니다. 일반적인 선택지는 다음과 같습니다:

- **dunst** — 가볍고 구성 가능성이 매우 높습니다(`apt install dunst` / `dnf install dunst`)
- **mako** — Wayland 네이티브입니다(`apt install mako-notifier`)
- **GNOME Shell** — GNOME 43 이상에서는 인터페이스를 자동으로 등록합니다. Ubuntu 24.04(GNOME Shell 46)에서는 세션 시작 시 인터페이스가 자동 등록되지 않을 수 있습니다. 알림이 표시되지 않으면 대체 수단으로 `dunst`을 설치하세요.
- **xfce4-notifyd** — XFCE 데스크톱에 번들로 제공됩니다

실행 중인 데몬이 없으면 `SendNotification`은 D-Bus 오류 `The name org.freedesktop.Notifications was not provided by any .service files`을 반환합니다. 앱에서 이 오류를 처리하고 사용자에게 알림 데몬을 설치하도록 안내하세요.

@end

Linux의 알림은 다음과 같습니다:

- 데스크톱 환경의 테마를 따릅니다
- 데스크톱 환경의 규칙에 따라 배치됩니다
- `Subtitle`을 지원합니다(이를 별도로 렌더링하지 않는 데몬에서는 본문에 이어 붙입니다)
- 사용자 텍스트 입력을 지원하지 않습니다(freedesktop 명세에 포함되지 않음)
- `image-path` 힌트를 통해 이미지 `Attachment` 하나를 지원합니다
- `ThreadID` 지원(지원되는 경우 데몬에서 처리)
- `Sound.Name`은 `sound-name` 힌트로 전달되며, 재생 여부는 활성 데몬과 사운드 테마에 따라 달라집니다.
- `InterruptionLevel`을 freedesktop `urgency` 힌트에 매핑합니다.
- 프로세스 내 타이머를 통한 예약 전송을 지원합니다. 단, **전송 전에 앱이 종료되면 예약된 알림이 사라집니다**.
- `UpdateNotification`은 D-Bus `replaces_id` 매개변수를 사용하여 이전 알림을 그 자리에서 대체합니다.

@end

## 모범 사례

1. 권한을 확인하고 요청하세요.
  - macOS에서는 사용자 권한이 필요합니다.


2. 명확하고 간결한 알림을 제공하세요.
  - 내용을 잘 나타내는 제목, 부제목, 본문 및 작업 제목을 사용하세요.


3. 알림 응답을 적절히 처리하세요.
  - 알림 응답에 오류가 있는지 확인하세요.
  - 사용자 작업에 대한 피드백을 제공하세요.


4. 플랫폼 관례를 고려하세요.
  - 플랫폼별 알림 패턴을 따르세요.
  - 시스템 설정을 준수하세요.


5. Linux에서는 데몬 종속성을 처리하세요.
  - `SendNotification`에서 반환된 오류를 확인하세요. 데몬이 없으면 D-Bus 오류가 발생합니다.
  - 패키지 문서나 앱 README에 freedesktop 알림 데몬이 필요하다고 명시하는 것이 좋습니다.


## 예제

다음 예제를 살펴보세요.

- [알림](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## API 참조

### 서비스 관리

| 메서드 | 설명 |
| --- | --- |
| `New()` | 새 알림 서비스를 생성합니다. |

### 알림 권한

| 메서드 | 설명 |
| --- | --- |
| `RequestNotificationAuthorization()` | 알림 표시 권한을 요청합니다(macOS). |
| `CheckNotificationAuthorization()` | 현재 알림 권한 상태를 확인합니다(macOS). |

### 알림 보내기

| 메서드 | 설명 |
| --- | --- |
| `SendNotification(options NotificationOptions)` | 기본 알림을 보냅니다. |
| `SendNotificationWithActions(options NotificationOptions)` | 작업이 포함된 대화형 알림을 보냅니다. |
| `UpdateNotification(options NotificationOptions)` | `ID`을 기준으로 현재 활성 상태인 알림을 업데이트합니다([알림 업데이트](#--10) 참조). |

### 알림 카테고리

| 메서드 | 설명 |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | 재사용 가능한 알림 카테고리를 등록합니다. |
| `RemoveNotificationCategory(categoryID string)` | 이전에 등록한 카테고리를 제거합니다. |

### 알림 관리

| 메서드 | 설명 |
| --- | --- |
| `RemoveAllPendingNotifications()` | 대기 중인 모든 알림을 제거합니다(macOS 및 Linux 전용). |
| `RemovePendingNotification(identifier string)` | 대기 중인 특정 알림을 제거합니다(macOS 및 Linux 전용). |
| `RemoveAllDeliveredNotifications()` | 전달된 모든 알림을 제거합니다(macOS 및 Linux 전용). |
| `RemoveDeliveredNotification(identifier string)` | 전달된 특정 알림을 제거합니다(macOS 및 Linux 전용). |
| `RemoveNotification(identifier string)` | 알림을 제거합니다(Linux 전용). |

### 이벤트 처리

| 메서드 | 설명 |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | 알림 응답을 처리할 콜백을 등록합니다. |

### 구조체 및 형식

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

#### InterruptionLevel 상수

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
