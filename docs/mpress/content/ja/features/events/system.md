---
title: "イベントシステム"
description: "イベントシステムを使用してコンポーネント間で通信する"
slug: "features/events/system"
sourcePath: "features/events/system.md"
---

## イベントシステム

Wails は、Pub/Sub 通信用の**統合イベントシステム**を提供します。Go から JavaScript、JavaScript から Go、ウィンドウ間など、どこからでもイベントを発行し、どこからでもリッスンできます。これにより、型付きイベントとライフサイクルフックを備えた疎結合アーキテクチャを実現できます。

## クイックスタート

**Go（発行）：**

```go
app.Event.Emit("user-logged-in", map[string]interface{}{
    "userId": 123,
    "name": "Alice",
})
```

**JavaScript（リッスン）：**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("user-logged-in", (event) => {
    console.log(`User ${event.data.name} logged in`)
})
```

**これだけです！** 言語間で Pub/Sub 通信を行えます。

## イベントの種類

### カスタムイベント

アプリケーション固有のイベント：

```go
// Emit from Go
app.Event.Emit("order-created", order)
app.Event.Emit("payment-processed", payment)
app.Event.Emit("notification", message)
```

```javascript
// Listen in JavaScript
Events.On("order-created", handleOrder)
Events.On("payment-processed", handlePayment)
Events.On("notification", showNotification)
```

### システムイベント

組み込みの OS およびアプリケーションイベント：

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Application lifecycle
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})
```

### ウィンドウイベント

ウィンドウ固有のイベント：

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

## イベントの発行

### Go から

**基本的な発行：**

```go
app.Event.Emit("event-name", data)
```

**さまざまなデータ型を使用する場合：**

```go
// String
app.Event.Emit("message", "Hello")

// Number
app.Event.Emit("count", 42)

// Struct
app.Event.Emit("user", User{ID: 1, Name: "Alice"})

// Map
app.Event.Emit("config", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// Array
app.Event.Emit("items", []string{"a", "b", "c"})
```

**特定のウィンドウに発行する場合：**

```go
window.EmitEvent("window-specific-event", data)
```

### JavaScript から

```javascript
import { Events } from '@wailsio/runtime'

// Emit to Go
Events.Emit("button-clicked", { buttonId: "submit" })

// Emit to all windows
Events.Emit("broadcast-message", "Hello everyone")
```

## イベントのリッスン

### Go で

**アプリケーションイベント：**

```go
app.Event.On("custom-event", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})
```

**型アサーションを使用する場合：**

```go
app.Event.On("user-updated", func(e *application.CustomEvent) {
    user := e.Data.(User)
    app.Logger.Info("User updated", "name", user.Name)
})
```

**複数のハンドラー：**

```go
// All handlers will be called
app.Event.On("order-created", logOrder)
app.Event.On("order-created", sendEmail)
app.Event.On("order-created", updateInventory)
```

**配信順序：** `On`、`Once`、または `OnMultiple` で登録された各ハンドラーは、一度に 1 つのイベントを処理します。そのため、同じハンドラーが自身と並行して実行されることはなく、連続する `Emit` 呼び出しを順番に受け取ります。ある `Emit` が完了してから次の呼び出しが行われた場合、すべてのハンドラーはその順序でイベントを受け取ります。
異なる goroutine から `Emit` を並行して呼び出した場合、それらの間の順序は定義されません。各ハンドラーは独立して実行されます。処理の遅いハンドラーが他のハンドラーを遅延させることはなく、ハンドラー間の順序も保証されません。`Emit` はハンドラーの完了を待ちません。ハンドラーを削除すると、そのハンドラーのキューに残っているイベントは破棄されます。

### JavaScript で

**基本的なリスナー：**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("event-name", (event) => {
    console.log("Event received:", event.data)
})
```

**クリーンアップを行う場合：**

```javascript
const unsubscribe = Events.On("event-name", handleEvent)

// Later, stop listening
unsubscribe()
```

**複数のハンドラー：**

```javascript
Events.On("data-updated", updateUI)
Events.On("data-updated", saveToCache)
Events.On("data-updated", logChange)
```

**一度だけ実行されるハンドラー：**

```javascript
Events.Once("data-updated", updateVariable)
```

## システムイベント

### アプリケーションイベント

**共通イベント（クロスプラットフォーム）：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Application started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("App started")
})

// Theme changed
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    isDark := e.Context().IsDarkMode()
    app.Event.Emit("theme-changed", isDark)
})

// System about to suspend
app.Event.OnApplicationEvent(events.Common.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})

// System resumed from suspend
app.Event.OnApplicationEvent(events.Common.SystemDidWake, func(e *application.ApplicationEvent) {
    reconnectSockets()
    refreshState()
})

// File opened
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(e *application.ApplicationEvent) {
    // Single-file: e.Context().Filename() returns ""; for multi-file launches
    // (Finder "Open With..." → multiple selection) use OpenedFiles().
    filePath := e.Context().Filename()
    openFile(filePath)
})
```

**プラットフォーム固有のイベント：**

@tabs{sync-key="platform"}
[macOS]
```go
// Application became active
app.Event.OnApplicationEvent(events.Mac.ApplicationDidBecomeActive, func(e *application.ApplicationEvent) {
    app.Logger.Info("App became active")
})

// Application will terminate
app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(e *application.ApplicationEvent) {
    cleanup()
})

// System about to sleep / resumed
app.Event.OnApplicationEvent(events.Mac.ApplicationWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})

// Displays sleeping/waking — distinct from system sleep (e.g. lid lowered)
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidSleep, func(e *application.ApplicationEvent) {
    pauseRendering()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidWake, func(e *application.ApplicationEvent) {
    resumeRendering()
})
```

[Windows]
```go
// Power status changed
app.Event.OnApplicationEvent(events.Windows.APMPowerStatusChange, func(e *application.ApplicationEvent) {
    app.Logger.Info("Power status changed")
})

// System suspending
app.Event.OnApplicationEvent(events.Windows.APMSuspend, func(e *application.ApplicationEvent) {
    saveState()
})

// System resumed. APMResumeAutomatic always fires on resume;
// APMResumeSuspend fires after it when the wake was triggered by
// user input. Prefer Common.SystemDidWake if you don't need to
// distinguish.
app.Event.OnApplicationEvent(events.Windows.APMResumeAutomatic, func(e *application.ApplicationEvent) {
    refreshState()
})
```

[Linux]
```go
// Application startup
app.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(e *application.ApplicationEvent) {
    app.Logger.Info("App starting")
})

// Theme changed
app.Event.OnApplicationEvent(events.Linux.SystemThemeChanged, func(e *application.ApplicationEvent) {
    updateTheme()
})

// System sleep/wake — sourced from systemd-logind's PrepareForSleep
// signal on the system bus. Requires logind/elogind exposing
// org.freedesktop.login1, so it may be unavailable on distros/setups
// without it (Alpine, Void, some Devuan setups).
app.Event.OnApplicationEvent(events.Linux.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Linux.SystemDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})
```

@end

### ウィンドウイベント

**共通のウィンドウイベント：**

```go
// Window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Window blur
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window blurred")
})

// Window closing
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel()  // Prevent close
    }
})

// Window closed
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    cleanup()
})
```

## イベントフック

フックは標準リスナーより**前に**実行され、イベントを**キャンセル**できます：

```go
// Hook - runs first, can cancel
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        result := showConfirmdialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            e.Cancel()  // Prevent window close
        }
    }
})

// Standard listener - runs after hooks
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

**主な違い：**

| 機能 | フック | 標準リスナー |
| --- | --- | --- |
| 実行順序 | 最初に登録順で実行 | フックの後に実行、順序の保証なし |
| ブロッキング | 同期実行、次のフックをブロック | 非同期実行、ノンブロッキング |
| キャンセル可否 | 可 | 不可（すでに伝播済み） |
| ユースケース | 制御フロー、検証 | ログ記録、副作用 |

## イベントパターン

### Pub/Sub パターン

```go
// Publisher (service)
type OrderService struct {
    app *application.App
}

func (o *OrderService) CreateOrder(items []Item) (*Order, error) {
    order := &Order{Items: items}
    
    if err := o.saveOrder(order); err != nil {
        return nil, err
    }
    
    // Publish event
    o.app.Event.Emit("order-created", order)
    
    return order, nil
}

// Subscribers
app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    sendConfirmationEmail(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    updateInventory(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    logOrder(order)
})
```

### リクエスト／レスポンスパターン

```go
// Frontend requests data
Emit("get-user-data", { userId: 123 })

// Backend responds
app.Event.On("get-user-data", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    userId := int(data["userId"].(float64))
    
    user := getUserFromDB(userId)
    
    // Send response
    app.Event.Emit("user-data-response", user)
})
```

```javascript
// Frontend receives response
Events.On("user-data-response", (event) => {
    const user = event.data
    displayUser(user)
})
```

**注：** リクエスト／レスポンスには、**バインディングの方が適しています**。通知にはイベントを使用してください。

### ブロードキャストパターン

```go
// Broadcast to all windows
app.Event.Emit("global-notification", "System update available")

// Each window handles it
Events.On("global-notification", (event) => {
    const message = event.data
    showNotification(message)
})
```

### イベント集約

```go
type EventAggregator struct {
    events []Event
    mu     sync.Mutex
}

func (ea *EventAggregator) Add(event Event) {
    ea.mu.Lock()
    defer ea.mu.Unlock()
    
    ea.events = append(ea.events, event)
    
    // Emit batch every 100 events
    if len(ea.events) >= 100 {
        app.Event.Emit("event-batch", ea.events)
        ea.events = nil
    }
}
```

## 完全な例

**Go：**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type NotificationService struct {
    app *application.App
}

func (n *NotificationService) Notify(message string) {
    // Emit to all windows
    n.app.Event.Emit("notification", map[string]interface{}{
        "message":   message,
        "timestamp": time.Now(),
    })
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })
    
    notifService := &NotificationService{app: app}
    
    // System events
    app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
        isDark := e.Context().IsDarkMode()
        app.Event.Emit("theme-changed", isDark)
    })
    
    // Custom events from frontend
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        action := data["action"].(string)
        
        app.Logger.Info("User action", "action", action)
        
        // Respond
        notifService.Notify("Action completed: " + action)
    })
    
    // Window events
    window := app.Window.New()
    
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        app.Event.Emit("window-focused", window.Name())
    })
    
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        // Confirm before close
        app.Event.Emit("confirm-close", nil)
        e.Cancel()  // Wait for confirmation
    })
    
    app.Run()
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for notifications
Events.On("notification", (event) => {
    showNotification(event.data.message)
})

// Listen for theme changes
Events.On("theme-changed", (event) => {
    const isDark = event.data
    document.body.classList.toggle('dark', isDark)
})

// Listen for window focus
Events.On("window-focused", (event) => {
    const windowName = event.data
    console.log(`Window ${windowName} focused`)
})

// Handle close confirmation
Events.On("confirm-close", (event) => {
    if (confirm("Close window?")) {
        Events.Emit("close-confirmed", true)
    }
})

// Emit user actions
document.getElementById('button').addEventListener('click', () => {
    Events.Emit("user-action", { action: "button-clicked" })
})
```

## ベストプラクティス

### ✅ 推奨事項

- **通知にはイベントを使用する** - 一方向通信
- **リクエストにはバインディングを使用する** - 双方向通信
- **イベント名に一貫性を持たせる** - kebab-case を使用する
- **イベントデータを文書化する** - どのフィールドが含まれるかを明記する
- **不要になったら購読を解除する** - メモリリークを防ぐ
- **検証にはフックを使用する** - イベントフローを制御する

### ❌ 非推奨事項

- **RPC にイベントを使用しない** - 代わりにバインディングを使用する
- **イベントを高頻度で発行しない** - 必要に応じてバッチ処理する
- **ハンドラー内でブロックしない** - 高速に処理できるようにする
- **購読解除を忘れない** - メモリリークの原因になる
- **大きなデータにイベントを使用しない** - バインディングを使用する
- **イベントループを作成しない** - A が B を発行し、B が A を発行する

## 次のステップ

@cards{cols="2"}
★ イベントガイド
イベントパターンと型安全なイベント生成について学びます。

[詳しく見る →](/guides/events-reference/)

---
◆ イベント API
アプリケーションとウィンドウのイベント API 全体を確認します。

[詳しく見る →](/reference/events/)

---
🚀 バインディング
リクエスト／レスポンスにはバインディングを使用します。

[詳しく見る →](/features/bindings/methods/)

---
▣ ウィンドウイベント
ウィンドウのライフサイクルイベントを処理します。

[詳しく見る →](/features/windows/events/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[イベントの例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/events)を確認してください。
