---
title: "イベント API"
description: "Events API の完全なリファレンス"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## 概要

Events API は、イベントを発行および購読するためのメソッドを提供し、アプリケーションの異なる部分間での通信を可能にします。

**イベントの種類：**

- **アプリケーションイベント** - アプリのライフサイクルイベント（起動、終了）
- **ウィンドウイベント** - ウィンドウの状態変化（フォーカスの取得、喪失、サイズ変更）
- **カスタムイベント** - アプリ固有の通信に使用するユーザー定義イベント

**通信パターン：**

- **Go からフロントエンドへ** - Go からイベントを発行し、JavaScript で購読
- **フロントエンドから Go へ** - 直接は不可（代わりにサービスバインディングを使用）
- **フロントエンドからフロントエンドへ** - Go またはローカルランタイムイベントを経由
- **ウィンドウからウィンドウへ** - 特定のウィンドウを対象にするか、すべてのウィンドウへブロードキャスト

## イベントメソッド（Go）

### app.Event.Emit()

すべてのウィンドウにカスタムイベントを発行します。フックによって発行がキャンセルされた場合は、`true`を返します。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**パラメーター：**

- `name` - イベント名
- `data` - イベントとともに送信する任意のデータ

**例：**

```go
// Emit simple event
app.Event.Emit("user-logged-in")

// Emit with data
app.Event.Emit("data-updated", map[string]interface{}{
    "count": 42,
    "status": "success",
})

// Emit multiple values
app.Event.Emit("progress", 75, "Processing files...")
```

### app.Event.On()

Go でカスタムイベントを購読します。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**パラメーター：**

- `name` - 購読するイベントの名前
- `callback` - イベントの発行時に呼び出される関数

**戻り値：** イベントリスナーを削除するクリーンアップ関数

**例：**

```go
// Listen for events
cleanup := app.Event.On("user-action", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    action := data["action"].(string)
    app.Logger.Info("User action", "action", action)
})

// Later, remove listener
cleanup()
```

### ウィンドウ固有のイベント

特定のウィンドウにイベントを発行します：

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## イベントメソッド（フロントエンド）

### On()

Go からのイベントを購読します。

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**パラメーター：**

- `eventName` - 購読するイベントの名前
- `callback` - イベントの受信時に呼び出される関数

**戻り値：** クリーンアップ関数

**例：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for events
const cleanup = Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
    updateUI(data)
})

// Later, remove listener
cleanup()
```

### Once()

イベントを一度だけ購読します。

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**例：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

1 つ以上のイベントに登録されたすべてのリスナーを削除します。`Off`は可変個のイベント名文字列を受け取り、コールバックは<strong>受け取りません</strong>。単一のリスナーを削除するには、`Events.On(...)`が返す購読解除関数を保持しておき、その関数を呼び出します。

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**例：**

```javascript
import { Events } from '@wailsio/runtime'

// Preferred: keep the unsubscribe fn from On()
const unsubscribe = Events.On('my-event', (data) => {
    console.log('Event received:', data)
})

// Later — remove just this listener
unsubscribe()

// Or: remove every listener for one or more events
Events.Off('my-event', 'another-event')
```

### OffAll()

イベントリスナーを<strong>すべて</strong>削除します。引数は取りません。

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**例：**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

イベントを最大`max`回まで購読し、その後、自動的に購読を解除します。

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**例：**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## アプリケーションイベント

### app.Event.OnApplicationEvent()

アプリケーションのライフサイクルイベントを購読します。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

イベント定数は`events`パッケージにあります。クロスプラットフォームイベントには`events.Common.*`を、プラットフォーム固有のイベントには`events.Mac.*`、`events.Windows.*`、`events.Linux.*`を使用します。

**一般的なアプリケーションイベント：**

- `events.Common.ApplicationStarted` - アプリケーションの起動が完了しました。
- `events.Common.ThemeChanged` - システムテーマがライトとダークの間で切り替わりました。
- `events.Common.ApplicationOpenedWithFile` - ファイルの関連付けを介して起動されました。
- `events.Common.ApplicationLaunchedWithUrl` - URL スキームを介して起動されました。

汎用の「アプリケーション終了」イベント定数は<strong>ありません</strong>。終了時のクリーンアップは、`application.Options.OnShutdown`または`app.OnShutdown(func())`を使用して登録してください。

**例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle application startup
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})

// React to theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Shutdown cleanup is configured on the application, not as an event:
app.OnShutdown(func() {
    database.Close()
    saveSettings()
})
```

## ウィンドウイベント

### OnWindowEvent()

ウィンドウ固有のイベントを購読します。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

ウィンドウイベント定数は`events`パッケージにあります。クロスプラットフォームイベントには`events.Common.*`を使用します（プラットフォーム固有のイベントには`events.Mac.*`、`events.Windows.*`、`events.Linux.*`を使用します）。

**一般的なウィンドウイベント：**

- `events.Common.WindowFocus` - ウィンドウがフォーカスを取得しました。
- `events.Common.WindowLostFocus` - ウィンドウがフォーカスを失いました。
- `events.Common.WindowClosing` - ウィンドウが閉じられようとしています（`RegisterHook`でキャンセル可能）。
- `events.Common.WindowDidResize` - ウィンドウのサイズが変更されました。
- `events.Common.WindowDidMove` - ウィンドウが移動されました。
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`。
- `events.Common.WindowRuntimeReady` - ウィンドウ内のランタイムへ安全にイベントを送出できます。

**例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Handle window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    app.Logger.Info("Window resized", "width", width, "height", height)
})
```

## 一般的なパターン

以下のパターンは、実際のアプリケーションでイベントを使用するための実証済みの手法を示します。各パターンは、Goバックエンドとフロントエンド間の特定の通信上の課題を解決し、応答性に優れた適切な構造のアプリケーションを構築するのに役立ちます。

### リクエスト／レスポンスパターン

データの取得、ファイル処理、バックグラウンドタスクなど、バックエンド処理の完了をフロントエンドに通知する場合に使用します。サービスバインディングはデータを直接返し、イベントはトーストメッセージの表示やリストの更新といったUI更新のための追加通知を提供します。

**Go：**

```go
// Service method
type DataService struct {
    app *application.App
}

func (s *DataService) FetchData(query string) ([]Item, error) {
    items := fetchFromDatabase(query)

    // Emit event when done
    s.app.Event.Emit("data-fetched", map[string]interface{}{
        "query": query,
        "count": len(items),
    })

    return items, nil
}
```

**JavaScript：**

```javascript
import { FetchData } from './bindings/DataService'
import { Events } from '@wailsio/runtime'

// Listen for completion event
Events.On('data-fetched', (data) => {
    console.log(`Fetched ${data.count} items for query: ${data.query}`)
    showNotification(`Found ${data.count} results`)
})

// Call service method
const items = await FetchData("search term")
displayItems(items)
```

### 進捗状況の更新

ファイルのアップロード、バッチ処理、大量データのインポート、動画のエンコードなど、長時間実行される処理に最適です。処理中に進捗イベントを送出してUIのプログレスバー、ステータステキスト、ステップインジケーターを更新し、ユーザーにリアルタイムのフィードバックを提供します。

**Go：**

```go
func (s *Service) ProcessFiles(files []string) error {
    total := len(files)

    for i, file := range files {
        // Process file
        processFile(file)

        // Emit progress event
        s.app.Event.Emit("progress", map[string]interface{}{
            "current": i + 1,
            "total":   total,
            "percent": float64(i+1) / float64(total) * 100,
            "file":    file,
        })
    }

    s.app.Event.Emit("processing-complete")
    return nil
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Update progress bar
Events.On('progress', (data) => {
    progressBar.style.width = `${data.percent}%`
    statusText.textContent = `Processing ${data.file}... (${data.current}/${data.total})`
})

// Handle completion
Events.Once('processing-complete', () => {
    progressBar.style.width = '100%'
    statusText.textContent = 'Complete!'
    setTimeout(() => hideProgressBar(), 2000)
})
```

### 複数ウィンドウ間の通信

設定パネル、ダッシュボード、ドキュメントビューアーなど、複数のウィンドウを持つアプリケーションに最適です。イベントをブロードキャストして、すべてのウィンドウ間で状態（テーマの変更やユーザー設定）を同期したり、特定のウィンドウに対象を限定したイベントを送信して、ウィンドウ固有の更新を行ったりできます。

**Go：**

```go
// Broadcast to all windows
app.Event.Emit("theme-changed", "dark")

// Send to specific window
preferencesWindow.EmitEvent("settings-updated", settings)

// Per-window listener — receive on the global event bus, but
// gate by the event's source-window name (set automatically when
// a window emits via window.EmitEvent).
app.Event.On("request-data", func(e *application.CustomEvent) {
    if e.Sender != window1.Name() {
        return
    }
    window1.EmitEvent("data-response", data)
})
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### 状態の同期

ユーザーセッション、アプリケーション構成、共同作業機能など、フロントエンドとバックエンドの状態を同期させる必要がある場合に使用します。バックエンドの状態が変化したら、接続されているすべてのフロントエンドを更新するイベントを送出し、アプリケーション全体の一貫性を確保します。

**Go：**

```go
type StateService struct {
    app   *application.App
    state map[string]interface{}
    mu    sync.RWMutex
}

func (s *StateService) UpdateState(key string, value interface{}) {
    s.mu.Lock()
    s.state[key] = value
    s.mu.Unlock()

    // Notify all windows
    s.app.Event.Emit("state-updated", map[string]interface{}{
        "key":   key,
        "value": value,
    })
}

func (s *StateService) GetState(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state[key]
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'
import { GetState } from './bindings/StateService'

// Keep local state in sync
let localState = {}

Events.On('state-updated', async (data) => {
    localState[data.key] = data.value
    updateUI(data.key, data.value)
})

// Initialize state
const initialState = await GetState("all")
localState = initialState
```

### イベント駆動型の通知

成功の確認、エラーアラート、情報メッセージなど、ユーザーへのフィードバックを表示する場合に最適です。サービスからUIコードを直接呼び出す代わりに、フロントエンドが一貫した方法で処理する通知イベントを送出します。これにより、通知スタイルの変更や通知履歴などの機能追加が容易になります。

**Go：**

```go
type NotificationService struct {
    app *application.App
}

func (s *NotificationService) Success(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "success",
        "message": message,
    })
}

func (s *NotificationService) Error(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "error",
        "message": message,
    })
}

func (s *NotificationService) Info(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "info",
        "message": message,
    })
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Unified notification handler
Events.On('notification', (data) => {
    const toast = document.createElement('div')
    toast.className = `toast toast-${data.type}`
    toast.textContent = data.message

    document.body.appendChild(toast)

    setTimeout(() => {
        toast.classList.add('fade-out')
        setTimeout(() => toast.remove(), 300)
    }, 3000)
})
```

## 完全な例

**Go：**

```go
package main

import (
    "sync"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type EventDemoService struct {
    app *application.App
    mu  sync.Mutex
}

func NewEventDemoService(app *application.App) *EventDemoService {
    service := &EventDemoService{app: app}

    // Listen for custom events
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        app.Logger.Info("User action received", "data", data)
    })

    return service
}

func (s *EventDemoService) StartLongTask() {
    go func() {
        s.app.Event.Emit("task-started")

        for i := 1; i <= 10; i++ {
            time.Sleep(500 * time.Millisecond)

            s.app.Event.Emit("task-progress", map[string]interface{}{
                "step":    i,
                "total":   10,
                "percent": i * 10,
            })
        }

        s.app.Event.Emit("task-completed", map[string]interface{}{
            "message": "Task finished successfully!",
        })
    }()
}

func (s *EventDemoService) BroadcastMessage(message string) {
    s.app.Event.Emit("broadcast", message)
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    // Handle application lifecycle
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
        app.Logger.Info("Application started!")
    })

    app.OnShutdown(func() {
        app.Logger.Info("Application shutting down...")
    })

    // Register service (RegisterService returns nothing).
    service := NewEventDemoService(app)
    app.RegisterService(application.NewService(service))

    // Create window
    window := app.Window.New()

    // Handle window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "focused")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "blurred")
    })

    window.Show()
    app.Run()
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'
import { StartLongTask, BroadcastMessage } from './bindings/EventDemoService'

// Task events
Events.On('task-started', () => {
    console.log('Task started...')
    document.getElementById('status').textContent = 'Running...'
})

Events.On('task-progress', (data) => {
    const progressBar = document.getElementById('progress')
    progressBar.style.width = `${data.percent}%`
    console.log(`Step ${data.step} of ${data.total}`)
})

Events.Once('task-completed', (data) => {
    console.log('Task completed!', data.message)
    document.getElementById('status').textContent = data.message
})

// Broadcast events
Events.On('broadcast', (message) => {
    console.log('Broadcast:', message)
    alert(message)
})

// Window state events
Events.On('window-state', (state) => {
    console.log('Window is now:', state)
    document.body.dataset.windowState = state
})

// Trigger long task
document.getElementById('startTask').addEventListener('click', async () => {
    await StartLongTask()
})

// Send broadcast
document.getElementById('broadcast').addEventListener('click', async () => {
    const message = document.getElementById('message').value
    await BroadcastMessage(message)
})
```

## 組み込みイベント

Wailsには、アプリケーションとウィンドウのライフサイクルに対応する組み込みのシステムイベントが用意されています。これらのイベントはフレームワークによって自動的に送出されます。

### 共通イベントとプラットフォームネイティブイベントの比較

Wailsは2種類のシステムイベントを提供します。

**共通イベント**（`events.Common.*`）は、macOS、Windows、Linuxで一貫して動作するクロスプラットフォームの抽象化です。移植性を最大限に高めるため、アプリケーションではこれらのイベントを使用してください。

**プラットフォームネイティブイベント**（`events.Mac.*`、`events.Windows.*`、`events.Linux.*`）は、共通イベントへのマッピング元となる、基盤のOS固有イベントです。これらを使用すると、プラットフォーム固有の動作やエッジケースにアクセスできます。

**仕組み：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// ✅ RECOMMENDED: Use Common Events for cross-platform code
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    // This works on all platforms
})

// Platform-specific events for advanced use cases
window.OnWindowEvent(events.Mac.WindowWillClose, func(e *application.WindowEvent) {
    // macOS-specific "will close" event (before WindowClosing)
})

window.OnWindowEvent(events.Windows.WindowClosing, func(e *application.WindowEvent) {
    // Windows-specific close event
})
```

**イベントのマッピング：**

プラットフォームネイティブイベントは、自動的に共通イベントへマッピングされます。

- macOS：`events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows：`events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux：`events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

このマッピングはバックグラウンドで自動的に行われるため、`events.Common.WindowClosing`をリッスンすると、プラットフォームに関係なくそのイベントを受信できます。

**使い分け：**

- **共通イベントを使用する**：アプリケーションコードの99%では、プラットフォーム間で一貫した動作を提供する共通イベントを使用します
- **プラットフォームネイティブイベントを使用する**：共通イベントでは利用できないプラットフォーム固有の機能が必要な場合にのみ使用します（例：macOS固有のウィンドウライフサイクルイベント、Windowsの電源管理イベント）

### アプリケーションイベント

| イベント | 説明 | 送出されるタイミング | キャンセル可能 |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | ファイルを指定してアプリケーションを開いたとき | ファイルの関連付けなどにより、ファイルを指定してアプリが起動されたとき | いいえ |
| `ApplicationStarted` | アプリケーションの起動が完了した | アプリの初期化が完了し、使用可能になった後 | いいえ |
| `ApplicationLaunchedWithUrl` | URLを指定してアプリケーションが起動された | URLスキームを介してアプリが起動されたとき | いいえ |
| `ThemeChanged` | システムテーマが変更された | OSのテーマがライトモードとダークモードの間で切り替わったとき | いいえ |
| `SystemWillSleep` | システムがまもなくサスペンドする | OSがサスペンドする直前（macOS / Windows / logindを使用するLinux） | いいえ |
| `SystemDidWake` | システムがサスペンドから復帰した | スリープから復帰した直後（macOS / Windows / logindを使用するLinux） | いいえ |

**使用例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### ウィンドウイベント

| イベント | 説明 | 発生タイミング | キャンセル可能 |
| --- | --- | --- | --- |
| `WindowClosing` | ウィンドウがまもなく閉じる | ウィンドウが閉じる前（ユーザーがXボタンをクリックしたか、Close()が呼び出されたとき） | はい |
| `WindowDidMove` | ウィンドウが新しい位置に移動した | ウィンドウの位置が変更された後（デバウンス処理あり） | いいえ |
| `WindowDidResize` | ウィンドウのサイズが変更された | ウィンドウのサイズが変更された後 | いいえ |
| `WindowDPIChanged` | ウィンドウのDPIスケーリングが変更された | DPIが異なるモニター間で移動したとき（Windows） | いいえ |
| `WindowFilesDropped` | OSネイティブのドラッグ＆ドロップでファイルがドロップされた | OSからウィンドウにファイルがドロップされた後 | いいえ |
| `WindowFocus` | ウィンドウがフォーカスを取得した | ウィンドウがアクティブになったとき | いいえ |
| `WindowFullscreen` | ウィンドウがフルスクリーンになった | Fullscreen()の呼び出し後、またはユーザーがフルスクリーンに切り替えた後 | いいえ |
| `WindowHide` | ウィンドウが非表示になった | Hide()の呼び出し後、またはウィンドウが遮蔽された後 | いいえ |
| `WindowLostFocus` | ウィンドウがフォーカスを失った | ウィンドウが非アクティブになったとき | いいえ |
| `WindowMaximise` | ウィンドウが最大化された | Maximise()の呼び出し後、またはユーザーが最大化した後 | はい（macOS） |
| `WindowMinimise` | ウィンドウが最小化された | Minimise()の呼び出し後、またはユーザーが最小化した後 | はい（macOS） |
| `WindowRestore` | ウィンドウが最小化または最大化された状態から復元された | Restore()の呼び出し後（主にWindows） | いいえ |
| `WindowRuntimeReady` | Wailsランタイムの読み込みが完了し、使用可能になった | JavaScriptランタイムの初期化が完了したとき | いいえ |
| `WindowShow` | ウィンドウが表示された | Show() の実行後、またはウィンドウが表示されたとき | いいえ |
| `WindowUnFullscreen` | ウィンドウがフルスクリーン状態を終了した | UnFullscreen() の実行後、またはユーザーがフルスクリーンを終了したとき | いいえ |
| `WindowUnMaximise` | ウィンドウが最大化状態を終了した | UnMaximise() の実行後、またはユーザーが最大化を解除したとき | はい（macOS） |
| `WindowUnMinimise` | ウィンドウが最小化状態を終了した | UnMinimise()/Restore() の実行後、またはユーザーがウィンドウを復元したとき | はい（macOS） |
| `WindowZoomIn` | ウィンドウコンテンツのズーム倍率が上がった | ZoomIn() の呼び出し後（主に macOS） | はい（macOS） |
| `WindowZoomOut` | ウィンドウコンテンツのズーム倍率が下がった | ZoomOut() の呼び出し後（主に macOS） | はい（macOS） |
| `WindowZoomReset` | ウィンドウコンテンツのズーム倍率が 100% にリセットされた | ZoomReset() の呼び出し後（主に macOS） | はい（macOS） |
| `WindowDropZoneFilesDropped` | JS で定義されたドロップゾーンにファイルがドロップされた | ドロップゾーンが設定された要素にファイルがドロップされたとき | いいえ |

**使用方法：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window events
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Cancel window close
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    dlg := app.Dialog.Question().SetMessage("Close window?")
    yes := dlg.AddButton("Yes")
    no := dlg.AddButton("No")
    dlg.SetDefaultButton(yes)
    dlg.SetCancelButton(no)
    no.OnClick(func() { e.Cancel() })

    dlg.Show()
})

// Wait for runtime ready
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    app.Logger.Info("Runtime ready, safe to emit events to frontend")
    window.EmitEvent("app-initialized", data)
})
```

**重要な注意事項：**

- **WindowRuntimeReady** は極めて重要です。フロントエンドにイベントを送出する前に、このイベントを待ってください
- **WindowDidMove** と **WindowDidResize** は、イベントの大量発生を防ぐためにデバウンスされます（デフォルトは 50ms）
- <strong>キャンセル可能なイベント</strong>は、`RegisterHook()` ハンドラー内で `event.Cancel()` を呼び出すことでキャンセルできます
- **WindowFilesDropped** は OS ネイティブのファイルドロップ用で、**WindowDropZoneFilesDropped** は Web ベースのドロップゾーン用です
- 一部のイベントはプラットフォーム固有です（例：Windows の WindowDPIChanged、主に macOS のズームイベント）

## イベントの命名規則

```go
// Good - descriptive and specific
app.Event.Emit("user:logged-in", user)
app.Event.Emit("data:fetch:complete", results)
app.Event.Emit("ui:theme:changed", theme)

// Bad - vague and unclear
app.Event.Emit("event1", data)
app.Event.Emit("update", stuff)
app.Event.Emit("e", value)
```

## パフォーマンスに関する考慮事項

### 高頻度イベントのデバウンス

```go
type Service struct {
    app            *application.App
    lastEmit       time.Time
    debounceWindow time.Duration
}

func (s *Service) EmitWithDebounce(event string, data interface{}) {
    now := time.Now()
    if now.Sub(s.lastEmit) < s.debounceWindow {
        return // Skip this emission
    }

    s.app.Event.Emit(event, data)
    s.lastEmit = now
}
```

### イベントのスロットリング

```javascript
import { Events } from '@wailsio/runtime'

let lastUpdate = 0
const throttleMs = 100

Events.On('high-frequency-event', (data) => {
    const now = Date.now()
    if (now - lastUpdate < throttleMs) {
        return // Skip this update
    }

    processUpdate(data)
    lastUpdate = now
})
```
