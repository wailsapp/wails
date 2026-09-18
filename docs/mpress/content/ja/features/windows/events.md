---
title: "ウィンドウイベント"
description: "ウィンドウのライフサイクルイベントと状態変更イベントを処理する"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## ウィンドウイベント

Wails は、単一の API を通じてウィンドウのライフサイクルイベントと状態変更イベントをディスパッチします。イベントを受け取るだけのリスナーには `OnWindowEvent` を、デフォルトのアクションをキャンセルできるフックには `RegisterHook` を使用します。

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listener — observes the event, cannot cancel it.
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) { /* ... */ })

// Hook — runs before listeners; can call e.Cancel() to suppress the default action.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel() // prevent the window from closing
    }
})
```

どちらの呼び出しも、ハンドラーを削除するために呼び出せる `unsubscribe func()` を返します。

クロスプラットフォームイベントは `events.Common.*` にあります。プラットフォーム固有のイベントは `events.Mac.*`、`events.Windows.*`、`events.Linux.*` にあります。完全な一覧は `v3/pkg/events/events.go` に生成されます。

## ライフサイクルイベント

### ウィンドウの作成

`app.Window.OnCreate` を使用すると、ウィンドウが作成されるたびにコールバックを実行できます。

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s (ID: %d)\n", window.Name(), window.ID())

    // Configure all new windows
    window.SetMinSize(400, 300)

    // Register a hook for this window
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if !confirmClose() {
            e.Cancel()
        }
    })
})
```

コールバックは `application.Window`（インターフェース）を受け取ります。`OnCreate` コールバックは、ウィンドウのランタイムが初期化された後、ウィンドウごとに一度呼び出されます。

### WindowClosing

ユーザーがウィンドウを閉じようとしたとき（X のクリック、⌘W、Alt+F4 など）に発生します。

**フック**（`RegisterHook`）を使用してください。終了をキャンセルできるのはフックだけです。リスナーはイベントを監視できますが、終了を阻止することはできません。

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**重要：**

- `WindowClosing` は、ユーザーが終了を試みたときにディスパッチされます。
- `RegisterHook` コールバックは `e.Cancel()` を呼び出して、ウィンドウを開いたままにできます。
- `WindowClosing` に対する `OnWindowEvent` コールバックは監視用です。呼び出されますが、キャンセルはできません。

**プログラムによる終了：** `window.Close()` を呼び出します。`window.Destroy()` は存在しません。

### WindowRuntimeReady

ウィンドウ内ランタイムの初期化が完了すると発生します。この時点から、ウィンドウの JS コンテキストを安全に呼び出せます。

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## フォーカスイベント

### WindowFocus

ウィンドウがフォーカスを得たときに呼び出されます。

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

ウィンドウがフォーカスを失ったときに呼び出されます。

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**例：フォーカス状態に応じた UI：**

```go
type FocusAwareWindow struct {
    window  *application.WebviewWindow
    focused bool
}

func (fw *FocusAwareWindow) Setup() {
    fw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        fw.focused = true
        fw.updateAppearance()
    })

    fw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        fw.focused = false
        fw.updateAppearance()
    })
}

func (fw *FocusAwareWindow) updateAppearance() {
    if fw.focused {
        fw.window.EmitEvent("update-theme", "active")
    } else {
        fw.window.EmitEvent("update-theme", "inactive")
    }
}
```

## 状態変更イベント

### WindowMinimise / WindowUnMinimise

```go
window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
    pauseRendering()
    saveWindowState()
})

window.OnWindowEvent(events.Common.WindowUnMinimise, func(e *application.WindowEvent) {
    resumeRendering()
    refreshContent()
})
```

### WindowMaximise / WindowUnMaximise

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "maximised")
})

window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "normal")
})
```

### WindowFullscreen / WindowUnFullscreen

```go
window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", false)
    window.EmitEvent("layout-mode", "fullscreen")
})

window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", true)
    window.EmitEvent("layout-mode", "normal")
})
```

`window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` を使用して、プログラムからフルスクリーンへの移行または解除を行います。`SetFullscreen(bool)` は存在しません。状態は `window.IsFullscreen() bool` で確認します。

## 位置とサイズのイベント

### WindowDidMove

```go
window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
    x, y := window.Position()
    fmt.Printf("Window moved to: %d, %d\n", x, y)
    saveWindowPosition(x, y)
})
```

### WindowDidResize

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    fmt.Printf("Window resized to: %dx%d\n", width, height)
    saveWindowSize(width, height)
    window.EmitEvent("window-size", map[string]int{
        "width":  width,
        "height": height,
    })
})
```

`WindowDidResize` と `WindowDidMove` のイベント自体には座標が含まれません。コールバック内で `window.Size()` / `window.Position()` を使用してウィンドウの情報を取得してください。

**例：レスポンシブレイアウト：**

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, _ := window.Size()
    var layout string
    switch {
    case width < 600:
        layout = "compact"
    case width < 1200:
        layout = "normal"
    default:
        layout = "wide"
    }
    window.EmitEvent("layout-changed", layout)
})
```

## 完全な例

すべてのイベント処理を備えた、本番環境で使用できるウィンドウの例です。

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type WindowState struct {
    X          int  `json:"x"`
    Y          int  `json:"y"`
    Width      int  `json:"width"`
    Height     int  `json:"height"`
    Maximised  bool `json:"maximised"`
    Fullscreen bool `json:"fullscreen"`
}

type ManagedWindow struct {
    app    *application.App
    window *application.WebviewWindow
    state  WindowState
    dirty  bool
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    mw := &ManagedWindow{app: app}
    mw.CreateWindow()
    mw.LoadState()
    mw.SetupEventHandlers()

    app.Run()
}

func (mw *ManagedWindow) CreateWindow() {
    mw.window = mw.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Event Demo",
        Width:  800,
        Height: 600,
    })
}

func (mw *ManagedWindow) SetupEventHandlers() {
    // Focus events
    mw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", true)
    })

    mw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", false)
    })

    // State change events
    mw.window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
        mw.SaveState()
    })

    mw.window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = false
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = false
        mw.dirty = true
    })

    // Position and size events
    mw.window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
        mw.state.X, mw.state.Y = mw.window.Position()
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
        mw.state.Width, mw.state.Height = mw.window.Size()
        mw.dirty = true
    })

    // Cancellable close — use a hook, not a listener.
    mw.window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if mw.dirty {
            mw.SaveState()
        }
    })
}

func (mw *ManagedWindow) LoadState() {
    data, err := os.ReadFile("window-state.json")
    if err != nil {
        return
    }

    if err := json.Unmarshal(data, &mw.state); err != nil {
        return
    }

    // Restore window state
    mw.window.SetPosition(mw.state.X, mw.state.Y)
    mw.window.SetSize(mw.state.Width, mw.state.Height)

    if mw.state.Maximised {
        mw.window.Maximise()
    }

    if mw.state.Fullscreen {
        mw.window.Fullscreen()
    }
}

func (mw *ManagedWindow) SaveState() {
    data, err := json.Marshal(mw.state)
    if err != nil {
        return
    }

    os.WriteFile("window-state.json", data, 0644)
    mw.dirty = false

    fmt.Println("Window state saved")
}
```

## イベントの連携

### ウィンドウ間イベント

アプリケーションのイベントバスを使用して、複数のウィンドウ間で連携します。

```go
// In main window
mainWindow.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Event.Emit("main-window-focused", nil)
})

// In other windows
app.Event.On("main-window-focused", func(event *application.CustomEvent) {
    updateRelativeToMain()
})
```

### イベントチェーン

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### デバウンスされたイベント

```go
var resizeTimer *time.Timer

window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    if resizeTimer != nil {
        resizeTimer.Stop()
    }

    resizeTimer = time.AfterFunc(500*time.Millisecond, func() {
        w, h := window.Size()
        saveWindowSize(w, h)
    })
})
```

## ベストプラクティス

### ✅ 推奨事項

- **キャンセルにはフックを使用する** — `e.Cancel()` を呼び出せるのは `RegisterHook` コールバックだけです。
- **終了時に状態を保存する** — 次回起動時にウィンドウの位置とサイズを復元します。
- **頻繁に発生するイベントをデバウンスする** — `WindowDidResize` と `WindowDidMove` は短い間隔で繰り返し発生します。
- **フォーカスの変更を処理する** — UI を適切に更新します。
- **イベントを介して連携する** — ウィンドウ間のメッセージ送信には `app.Event.Emit` を使用します。
- **購読を解除する** — ハンドラーが不要になったら解除してください。`OnWindowEvent` と `RegisterHook` は、どちらも購読解除用の `func()` を返します。

### ❌ 非推奨事項

- **イベントハンドラーをブロックしない** — 処理を短時間で完了させてください。
- **`OnWindowEvent`** からキャンセルしようとしないでください。`RegisterHook` を使用します。
- **`window.Destroy()`** を使用しないでください。これは存在しないため、`window.Close()` を使用します。
- **イベントのたびに保存しない** — 先にデバウンスしてください。
- **イベントに座標が含まれると想定しない** — `window.Position()` / `window.Size()` を呼び出してください。

## トラブルシューティング

### WindowClosing フックで終了を阻止できない

**原因：** `RegisterHook` ではなく `OnWindowEvent` を使用しています。

**解決策：** `e.Cancel()` を呼び出せるのは `RegisterHook` コールバックだけです。`OnWindowEvent` コールバックは監視用であり、キャンセルはできません。

```go
// ❌ Cannot cancel — this is a listener.
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel() // no effect
})

// ✅ Can cancel — this is a hook.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel()
})
```

### イベントが発生しない

<strong>原因：</strong>イベントの発生後にハンドラーを登録しています。

<strong>解決策：</strong>ウィンドウの作成直後（または `app.Window.OnCreate` コールバック内）にハンドラーを登録します。

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### メモリリーク

<strong>原因：</strong>有効期間の短い状態が、有効期間の長いハンドラーを保持しています。

**解決策：**`OnWindowEvent`/`RegisterHook` が返す購読解除用の `func()` を保持し、クリーンアップ時に呼び出します。

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## 次のステップ

**ウィンドウの基本** - ウィンドウ管理の基礎を学びます [詳細を見る →](/features/windows/basics/)

**複数ウィンドウ** - マルチウィンドウアプリケーションのパターンを学びます [詳細を見る →](/features/windows/multiple/)

**イベントシステム** - イベントシステムを詳しく学びます [詳細を見る →](/features/events/system/)

**アプリケーションのライフサイクル** - アプリケーションのライフサイクルについて理解します [詳細を見る →](/concepts/lifecycle/)

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
