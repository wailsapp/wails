---
title: "通知"
description: "アクションボタンとテキスト入力を備えたネイティブのシステム通知を表示します"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## はじめに

Wails は、デスクトップアプリケーション向けの包括的なクロスプラットフォーム通知システムを提供します。このサービスでは、次の機能を備えたネイティブのシステム通知を表示できます。

- タイトル、サブタイトル、本文を含む基本通知
- アクションボタンとテキスト返信を備えたインタラクティブ通知
- アクション用に再利用可能な[通知カテゴリ](#heading-6)
- カスタム[サウンド](#heading-10)（デフォルト、無音、または名前付き）
- [添付ファイル](#heading-11)（すべてのプラットフォームで画像、macOS では音声と動画にも対応）
- `ThreadID`による[関連通知のグループ化](#heading-13)
- `InterruptionLevel`による[優先度](#heading-14)の設定（`passive` / `active` / `timeSensitive` / `critical`）
- [配信のスケジュール設定](#heading-15)（macOS ではネイティブ機能、Windows と Linux ではプロセス内タイマーを使用）
- ID による[処理中の通知の更新](#heading-16)

新しいオプションフィールドはいずれも、プラットフォームが対応できない場合には適切に機能を縮退させます。機能別のサポート表については、[プラットフォームに関する考慮事項](#heading-17)を参照してください。

## 基本的な使い方

### サービスの作成

まず、通知サービスを初期化します。

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

Windows と Linux では、常に`true`が返されます。

## 通知の種類

### 基本通知

一意の ID、タイトル、省略可能なサブタイトル（macOS および Linux）、本文テキストを含む基本通知をユーザーに送信します。

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

通知の配信時に再生する音声は、`Sound`で制御します。`nil`のままにするとプラットフォームのデフォルト音が再生されます。サウンドを無効にするには`Silent: true`を設定し、名前付きまたはバンドル済みのサウンドを再生するには`Name`を設定します。

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

- **macOS** — `Name`は`[UNNotificationSound soundNamed:]`に渡されます。音声ファイルは、アプリバンドルの`Library/Sounds`配下に配置する必要があります。
- **Windows** — `Name`がすでに`ms-winsoundevent:`または`ms-appx:`で始まっている場合は、そのまま使用されます。それ以外の場合は、組み込みのトーストイベント名として使用できるように`ms-winsoundevent:`で囲まれます（Microsoft のトースト`<audio>`スキーマのドキュメントを参照してください）。
- **Linux** — freedesktop の`sound-name`ヒントとして渡されます。再生されるかどうかは、動作中の通知デーモンとサウンドテーマによって決まります。

### 添付ファイル

`Attachments`を使用すると、通知にメディアファイルを添付できます。macOS は、あらゆるメディアタイプの複数の添付ファイルに対応します。Windows と Linux では、画像タイプとして指定された最初の添付ファイルが使用されます。

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

`Path`には、ファイルシステム上の絶対パスを指定する必要があります。macOS では、`file://` URL も使用できます。

#### アプリに同梱したファイルの添付

OS は通知の配信時にディスクから添付ファイルを読み取るため、`Path`はエンドユーザーのマシン上に実在するファイルとして解決される必要があります。アプリケーションにバンドルするアセット（`go:embed`で埋め込んだアイコンや画像）は、ディスク上の既知の場所ではなくバイナリ内に存在するため、ハードコードできる固定の絶対パスはありません。起動時に一度だけ書き込み可能なディレクトリへ保存し、そのパスを渡します。

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

ユーザーが指定したファイルやダウンロードしたファイルには、すでにディスク上の実在するパスがあるため、この手順を実行せずに`Path`へ直接渡せます。添付ファイルのバイト列をメモリ内で渡す機能は、将来のリリースで追加される可能性があります。

### スレッド化とグループ化

`ThreadID`は関連する通知をまとめ、OS が Notification Center / Action Center で折りたためるようにします。

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
| `InterruptionLevelPassive` | `"passive"` | 控えめに配信され、画面を点灯させず、デフォルトのサウンドも再生しません |
| `InterruptionLevelActive` | `"active"` | デフォルトレベル |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | 許可されている場合、集中モード／応答不可モードを突破して通知します |
| `InterruptionLevelCritical` | `"critical"` | 集中モードと着信音設定を無視します。macOS では Critical Alert エンタイトルメントが必要です（ない場合は警告などを出さずに機能が低下します） |

プラットフォームごとのマッピング：

- **macOS** — `UNNotificationContent.interruptionLevel`を設定します。`critical`には macOS 12以降と Critical Alert エンタイトルメントが必要です。
- **Windows** — トーストの`<toast scenario="...">`属性にマッピングされます。
- **Linux** — freedesktop の`urgency`ヒントにマッピングされます。

### 配信のスケジュール

`Schedule`は配信を延期します。`DelaySeconds`（現在からの秒数）または`At`（UTC の Unix 秒）のいずれか一方だけを設定してください。

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
<strong>macOS</strong>では、スケジュールされた通知にネイティブトリガーが使用され、アプリを再起動しても維持されます。<strong>Windows</strong>と<strong>Linux</strong>では、スケジューリングはプロセス内の`time.AfterFunc`タイマーにフォールバックするため、**配信前にアプリが終了すると失われます**。`wintoast`にも freedesktop 仕様にも、遅延配信のプリミティブはありません。

@end

### 通知の更新

`UpdateNotification`は、同じ`ID`を持つ処理中の通知を置き換えます：

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

プラットフォームごとの動作：

- **macOS** — `UNUserNotificationCenter`が識別子に基づいて自動的に重複を排除するため、既存の通知がその場で更新されます。
- **Linux** — D-Bus の`replaces_id`パラメーターを使用して、以前の通知を置き換えます。
- **Windows** — 現在は新しい通知として再配信されます。既存の通知をその場で置き換えるには、上流の`wintoast`による`tag` / `group`のサポートが必要です。

## プラットフォームに関する考慮事項

@tabs
[macOS]
macOS の通知には、次の特性があります：

- ユーザーの許可が必要
- アプリのパッケージ化と署名が必要（配布する場合は公証も必要）
- システム標準の通知表示を使用
- `Subtitle`をサポート
- ユーザーによるテキスト入力（返信）をサポート
- `Destructive`アクションオプションをサポート
- 任意のメディア形式（画像、音声、動画）の複数の`Attachments`をサポート
- 通知センターでのグループ化に`ThreadID`をサポート
- すべての`InterruptionLevel`値をサポート（`critical`には Critical Alert エンタイトルメントが必要）
- アプリを再起動しても維持されるネイティブなスケジュール配信をサポート
- `UpdateNotification`呼び出しの重複を`ID`に基づいて自動的に排除
- ダークモードとライトモードを自動的に処理

[Windows]
Windows の通知には、次の特性があります：

- `wintoast`バックエンドを介して Windows システムのトーストスタイルを使用
- Windows のテーマ設定に適応
- ユーザーによるテキスト入力（返信）をサポート
- 高 DPI ディスプレイをサポート
- `Subtitle`はサポートしない
- 配置ヒントとして`hero`、`appLogoOverride`、`inline`（デフォルトは`inline`）を指定できる、単一の画像`Attachment`をサポート
- アクションセンターでのグループ化に`ThreadID`をサポート
- トーストの`scenario`属性を介して`InterruptionLevel`をサポート
- プロセス内タイマーによるスケジュール配信をサポート。ただし、**配信前にアプリが終了すると、スケジュールされた通知は失われます**
- `UpdateNotification`は現在、新しい通知として再配信されます（既存の通知をその場で置き換えるには、上流の`wintoast`による`tag`/`group`のサポートが必要です）

[Linux]
Linux では、通知に D-Bus の`org.freedesktop.Notifications`インターフェースを使用します。通知を機能させるには、互換性のある通知デーモンが<strong>実行中でなければなりません</strong>。

@note{type="caution" title="システム要件：通知デーモン"}
freedesktop 互換の通知デーモンをインストールし、実行しておく必要があります。一般的な選択肢は次のとおりです：

- **dunst** — 軽量で高度な設定が可能（`apt install dunst` / `dnf install dunst`）
- **mako** — Wayland ネイティブ（`apt install mako-notifier`）
- **GNOME Shell** — GNOME 43以降では、インターフェースが自動的に登録されます。Ubuntu 24.04（GNOME Shell 46）では、セッション開始時にインターフェースが自動登録されない場合があります。通知が表示されない場合は、フォールバックとして`dunst`をインストールしてください。
- **xfce4-notifyd** — XFCE デスクトップに同梱

デーモンが実行されていない場合、`SendNotification`は D-Bus エラー`The name org.freedesktop.Notifications was not provided by any .service files`を返します。アプリでこのエラーを処理し、通知デーモンをインストールするようユーザーに案内してください。

@end

Linux の通知には、次の特性があります：

- デスクトップ環境のテーマに従う
- デスクトップ環境の規則に従って配置
- `Subtitle`をサポート（個別に表示しないデーモンでは本文に連結）
- ユーザーによるテキスト入力はサポートしない（freedesktop 仕様に含まれないため）
- `image-path`ヒントを介して単一の画像`Attachment`をサポート
- `ThreadID` をサポート（対応している場合はデーモンが処理）
- `Sound.Name` は `sound-name` ヒントとして転送されます。再生されるかどうかは、稼働中のデーモンとサウンドテーマによって異なります
- `InterruptionLevel` を freedesktop の `urgency` ヒントにマッピング
- プロセス内タイマーによるスケジュール配信をサポート — **配信前にアプリが終了すると、スケジュールされた通知は失われます**
- `UpdateNotification` は D-Bus の `replaces_id` パラメーターを使用して、以前の通知をその場で置き換えます

@end

## ベストプラクティス

1. 通知の許可状態を確認し、許可を要求します：
  - macOS ではユーザーの許可が必要です


2. 明確で簡潔な通知を提供します：
  - 内容が分かりやすいタイトル、サブタイトル、テキスト、アクションタイトルを使用します


3. 通知への応答を適切に処理します：
  - 通知への応答にエラーがないか確認します
  - ユーザーの操作に対するフィードバックを提供します


4. プラットフォームの慣例を考慮します：
  - プラットフォーム固有の通知パターンに従います
  - システム設定を尊重します


5. Linux では、デーモンへの依存を適切に処理します：
  - `SendNotification` が返すエラーを確認します。デーモンが存在しない場合は D-Bus エラーが発生します
  - パッケージのドキュメントまたはアプリの README に、freedesktop 通知デーモンが必要であることを明記してください


## 例

次の例を参照してください：

- [通知](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## API リファレンス

### サービス管理

| メソッド | 説明 |
| --- | --- |
| `New()` | 新しい通知サービスを作成します |

### 通知の許可

| メソッド | 説明 |
| --- | --- |
| `RequestNotificationAuthorization()` | 通知を表示する許可を要求します（macOS） |
| `CheckNotificationAuthorization()` | 現在の通知の許可状態を確認します（macOS） |

### 通知の送信

| メソッド | 説明 |
| --- | --- |
| `SendNotification(options NotificationOptions)` | 基本的な通知を送信します |
| `SendNotificationWithActions(options NotificationOptions)` | アクション付きの対話型通知を送信します |
| `UpdateNotification(options NotificationOptions)` | 送信済みで処理中の通知を `ID` によって更新します（[通知の更新](#heading-16)を参照） |

### 通知カテゴリー

| メソッド | 説明 |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | 再利用可能な通知カテゴリーを登録します |
| `RemoveNotificationCategory(categoryID string)` | 以前に登録したカテゴリーを削除します |

### 通知の管理

| メソッド | 説明 |
| --- | --- |
| `RemoveAllPendingNotifications()` | 保留中の通知をすべて削除します（macOS および Linux のみ） |
| `RemovePendingNotification(identifier string)` | 指定した保留中の通知を削除します（macOS および Linux のみ） |
| `RemoveAllDeliveredNotifications()` | 配信済みの通知をすべて削除します（macOS および Linux のみ） |
| `RemoveDeliveredNotification(identifier string)` | 指定した配信済みの通知を削除します（macOS および Linux のみ） |
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
