---
title: "通知"
description: "アクションボタンとテキスト入力を備えたネイティブシステム通知を表示する"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## はじめに

Wails は、デスクトップアプリケーション向けの包括的なクロスプラットフォーム通知システムを提供します。このサービスでは、次の機能を備えたネイティブシステム通知を表示できます。

- タイトル、サブタイトル、本文を含む基本通知
- アクションボタンとテキスト返信を備えたインタラクティブ通知
- アクションに使用できる再利用可能な[通知カテゴリ](#heading-6)
- カスタム[サウンド](#heading-10)（デフォルト、無音、または名前付き）
- [添付ファイル](#heading-11)（すべてのプラットフォームで画像、macOS では音声と動画にも対応）
- `ThreadID`による[関連通知のグループ化](#heading-13)
- `InterruptionLevel`による[優先度](#heading-14)の指定（`passive` / `active` / `timeSensitive` / `critical`）
- [配信予約](#heading-15)（macOS ではネイティブ機能、Windows と Linux ではプロセス内タイマーを使用）
- ID による[配信中の通知の更新](#heading-16)

新しいオプションフィールドは、その機能にプラットフォームが対応できない場合でも適切にフォールバックします。機能ごとの対応状況については、[プラットフォームに関する考慮事項](#heading-17)の表を参照してください。

## 基本的な使い方

### サービスの作成

最初に、通知サービスを初期化します。

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

## 通知の許可

macOS で通知を使用するには、ユーザーの許可が必要です。許可を要求し、その状態を確認します。

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

Windows と Linux では、常に `true` を返します。

## 通知の種類

### 基本通知

一意の ID、タイトル、省略可能なサブタイトル（macOS と Linux）、本文を含む基本通知をユーザーに送信します。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### インタラクティブ通知

アクションボタンとテキスト入力を備えた通知を送信します。この通知を使用するには、事前に通知カテゴリを登録する必要があります。

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

## 通知への応答

通知に対するユーザー操作を処理します。

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

## 通知のカスタマイズ

### カスタムメタデータ

基本通知とインタラクティブ通知には、カスタムデータを含めることができます。

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

### カスタムサウンド

通知の配信時に再生する音声は、`Sound`で制御します。`nil`のままにするとプラットフォームのデフォルト音が再生されます。音を鳴らさない場合は`Silent: true`を設定し、名前付きまたはバンドル済みのサウンドを再生する場合は`Name`を設定します。

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

`Name`は、プラットフォームごとに次のように解決されます。

- **macOS** — `Name`が`[UNNotificationSound soundNamed:]`に渡されます。音声ファイルは、アプリバンドル内の`Library/Sounds`に配置する必要があります。
- **Windows** — `Name`がすでに`ms-winsoundevent:`または`ms-appx:`で始まっている場合は、そのまま使用されます。それ以外の場合は、組み込みのトーストイベント名として使用できるよう`ms-winsoundevent:`で囲まれます（Microsoft のトースト`<audio>`スキーマのドキュメントを参照してください）。
- **Linux** — freedesktop の`sound-name`ヒントとして転送されます。再生されるかどうかは、使用中の通知デーモンとサウンドテーマによって異なります。

### 添付ファイル

`Attachments`を使用すると、通知にメディアファイルを添付できます。macOS は、あらゆるメディア形式の複数の添付ファイルに対応しています。Windows と Linux では、画像形式として指定された最初の添付ファイルが使用されます。

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

`Path`には、ファイルシステムの絶対パスを指定する必要があります。macOS では、これに加えて`file://` URL も使用できます。

#### アプリに同梱したファイルの添付

OS は通知の配信時にディスクから添付ファイルを読み込むため、`Path`はエンドユーザーのマシン上に実在するファイルに解決される必要があります。アプリケーションにバンドルしたアセット（`go:embed`で埋め込んだアイコンや画像）は、ディスク上の既知の場所ではなくバイナリ内に存在するため、ハードコードできる固定の絶対パスはありません。起動時に一度だけ書き込み可能なディレクトリへ保存し、そのパスを渡します。

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

ユーザーが提供したファイルやダウンロード済みのファイルには、すでにディスク上の実際のパスがあるため、この手順を行わずに`Path`へ直接渡せます。メモリ内の添付ファイルのバイト列を渡す機能は、将来のリリースで追加される可能性があります。

### スレッド化とグループ化

`ThreadID`は関連する通知をグループ化し、OS が通知センターまたはアクションセンターでまとめて表示できるようにします。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### 割り込みレベル

`InterruptionLevel`は通知の優先度を制御します。エクスポートされている次の定数のいずれかを使用します。

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| 定数 | 値 | 意味 |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | 静かに配信します。画面を点灯せず、デフォルト音も再生しません |
| `InterruptionLevelActive` | `"active"` | デフォルトレベル |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | 許可されている場合は、集中モード／応答不可を突破します |
| `InterruptionLevelCritical` | `"critical"` | 集中モードと着信音設定を迂回します。macOS では「重大な通知」のエンタイトルメントが必要です（ない場合は通知されずにフォールバックします） |

プラットフォームごとの対応付けは次のとおりです。

- **macOS** — `UNNotificationContent.interruptionLevel`を設定します。`critical`には macOS 12 以降と「重大な通知」のエンタイトルメントが必要です。
- **Windows** — トーストの`<toast scenario="...">`属性に対応付けられます。
- **Linux** — freedesktop の`urgency`ヒントに対応付けられます。

### 配信予約

`Schedule`は配信を延期します。`DelaySeconds`（現在からの秒数）または`At`（UTC の Unix 秒）のどちらか一方だけを設定してください。

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

@note{type="caution" title="永続性"}
**macOS**では、予約通知にネイティブのトリガーを使用するため、アプリを再起動しても保持されます。**Windows**と**Linux**では、プロセス内の`time.AfterFunc`タイマーにフォールバックするため、**配信前にアプリが終了すると失われます**。`wintoast`にも freedesktop 仕様にも、遅延配信のためのプリミティブはありません。

@end

### 通知の更新

`UpdateNotification`は、同じ`ID`を持つ配信中の通知を置き換えます。

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

プラットフォームごとの動作は次のとおりです。

- **macOS** — `UNUserNotificationCenter`は識別子によって自動的に重複を排除するため、既存の通知がその場で更新されます。
- **Linux** — D-Bus の `replaces_id` パラメーターを使用して、以前の通知を置き換えます。
- **Windows** — 現在は新しい通知として再配信されます。実際にその場で置き換えるには、上流の `wintoast` が `tag` / `group` をサポートする必要があります。

## プラットフォームに関する考慮事項

@tabs
[macOS]
macOS の通知には、次の特徴があります。

- ユーザーの許可が必要です
- アプリをパッケージ化して署名する必要があります（配布時には公証も必要です）
- システム標準の通知表示を使用します
- `Subtitle`をサポートします
- ユーザーによるテキスト入力（返信）をサポートします
- `Destructive` アクションオプションをサポートします
- あらゆるメディア形式（画像、音声、動画）の複数の `Attachments` をサポートします
- 通知センターでのグループ化に `ThreadID` を使用できます
- すべての `InterruptionLevel` 値をサポートします（`critical` には Critical Alert エンタイトルメントが必要です）
- アプリを再起動しても維持されるネイティブの予約配信をサポートします
- `UpdateNotification` の呼び出しを `ID` に基づいて自動的に重複排除します
- ダークモードとライトモードを自動的に処理します

[Windows]
Windows の通知には、次の特徴があります。

- `wintoast` バックエンドを介して、Windows システムのトーストスタイルを使用します
- Windows のテーマ設定に適応します
- ユーザーによるテキスト入力（返信）をサポートします
- 高 DPI ディスプレイをサポートします
- `Subtitle`はサポートしません
- 配置ヒントとして `hero`、`appLogoOverride`、または `inline`（デフォルトは `inline`）を指定した、単一の画像 `Attachment` をサポートします
- アクションセンターでのグループ化に `ThreadID` を使用できます
- トーストの `scenario` 属性を介して `InterruptionLevel` をサポートします
- プロセス内タイマーによる予約配信をサポートします — **配信前にアプリが終了すると、予約した通知は失われます**
- `UpdateNotification` は現在、新しい通知として再配信されます（実際にその場で置き換えるには、上流の `wintoast` による `tag`/`group` のサポートが必要です）

[Linux]
Linux では、通知に D-Bus の `org.freedesktop.Notifications` インターフェースを使用します。通知を機能させるには、互換性のある通知デーモンが**実行されている必要があります**。

@note{type="caution" title="システム要件：通知デーモン"}
freedesktop 互換の通知デーモンをインストールして実行する必要があります。一般的な選択肢は次のとおりです。

- **dunst** — 軽量で高度に設定可能（`apt install dunst` / `dnf install dunst`）
- **mako** — Wayland ネイティブ（`apt install mako-notifier`）
- **GNOME Shell** — GNOME 43 以降では、インターフェースを自動的に登録します。Ubuntu 24.04（GNOME Shell 46）では、セッション開始時にインターフェースが自動登録されない場合があります。通知が表示されない場合は、フォールバックとして `dunst` をインストールしてください。
- **xfce4-notifyd** — XFCE デスクトップに同梱されています

デーモンが実行されていない場合、`SendNotification` は D-Bus エラー `The name org.freedesktop.Notifications was not provided by any .service files` を返します。アプリでこのエラーを処理し、通知デーモンをインストールするようユーザーに案内してください。

@end

Linux の通知には、次の特徴があります。

- デスクトップ環境のテーマに従います
- デスクトップ環境の規則に従って配置されます
- `Subtitle`をサポートします（個別に描画しないデーモンでは本文に連結されます）
- ユーザーによるテキスト入力はサポートしません（freedesktop 仕様には含まれていません）
- `image-path` ヒントを介して、単一の画像 `Attachment` をサポートします
- `ThreadID`をサポートします（対応している場合はデーモンが処理します）
- `Sound.Name` は `sound-name` ヒントとして転送されます。再生されるかどうかは、アクティブなデーモンとサウンドテーマによって異なります
- `InterruptionLevel` を freedesktop の `urgency` ヒントにマッピングします
- プロセス内タイマーによる予約配信をサポートします — **配信前にアプリが終了すると、予約した通知は失われます**
- `UpdateNotification` は D-Bus の `replaces_id` パラメーターを使用し、以前の通知をその場で置き換えます

@end

## ベストプラクティス

1. 許可の状態を確認し、必要に応じて要求します。
  - macOS ではユーザーの許可が必要です


2. 明確で簡潔な通知を作成します。
  - 内容を明確に示すタイトル、サブタイトル、本文、アクションタイトルを使用します


3. 通知への応答を適切に処理します。
  - 通知への応答にエラーがないか確認します
  - ユーザーの操作に対するフィードバックを提供します


4. プラットフォームの規則を考慮します。
  - プラットフォーム固有の通知パターンに従います
  - システム設定を尊重します


5. Linux では、通知デーモンへの依存を処理します。
  - `SendNotification` が返すエラーを確認します — デーモンが存在しない場合は D-Bus エラーが発生します
  - パッケージのドキュメントまたはアプリの README に、freedesktop 通知デーモンが必要であることを明記してください


## 例

次の例を参照してください。

- [通知](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notifications)

## API リファレンス

### サービス管理

| メソッド | 説明 |
| --- | --- |
| `New()` | 新しい通知サービスを作成します |

### 通知の許可

| メソッド | 説明 |
| --- | --- |
| `RequestNotificationAuthorization()` | 通知を表示する許可を要求します（macOS） |
| `CheckNotificationAuthorization()` | 現在の通知許可の状態を確認します（macOS） |

### 通知の送信

| メソッド | 説明 |
| --- | --- |
| `SendNotification(options NotificationOptions)` | 基本的な通知を送信します |
| `SendNotificationWithActions(options NotificationOptions)` | アクション付きのインタラクティブ通知を送信します |
| `UpdateNotification(options NotificationOptions)` | `ID` に基づいて送信中の通知を更新します（[通知の更新](#heading-16)を参照） |

### 通知カテゴリ

| メソッド | 説明 |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | 再利用可能な通知カテゴリを登録します |
| `RemoveNotificationCategory(categoryID string)` | 以前に登録したカテゴリを削除します |

### 通知の管理

| メソッド | 説明 |
| --- | --- |
| `RemoveAllPendingNotifications()` | 保留中の通知をすべて削除します（macOS および Linux のみ） |
| `RemovePendingNotification(identifier string)` | 保留中の特定の通知を削除します（macOS および Linux のみ） |
| `RemoveAllDeliveredNotifications()` | 配信済みの通知をすべて削除します（macOS および Linux のみ） |
| `RemoveDeliveredNotification(identifier string)` | 配信済みの特定の通知を削除します（macOS および Linux のみ） |
| `RemoveNotification(identifier string)` | 通知を削除します（Linux 固有） |

### イベント処理

| メソッド | 説明 |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | 通知への応答を処理するコールバックを登録します |

### 構造体と型

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

#### InterruptionLevel 定数

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
