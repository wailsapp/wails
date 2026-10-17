---
title: "グローバルショートカット"
description: "アプリケーションにフォーカスがない場合でも発動する、システム全体のキーボードショートカットを登録します"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

グローバルショートカットは、Wails アプリケーションの実行中、現在どのアプリケーションにフォーカスがあるかに関係なく発動する、システム全体のキーボードショートカットです。表示／非表示を切り替えるホットキー、クイックキャプチャツール、メディアコントロールなど、ユーザーがどこからでも利用できることを期待する機能に最適です。

@note{type="info" title="グローバルショートカットとキーバインディングの違い"}
[キーバインディング](/features/keyboard/shortcuts/)（`app.KeyBinding`）は、アプリケーションのいずれかのウィンドウにフォーカスがある間だけ発動します。グローバルショートカット（`app.GlobalShortcut`）は、アプリケーションがバックグラウンドにある場合でも、システム全体で発動します。用途に合う方を使用してください。

@end

グローバルショートカットは各プラットフォームのネイティブ機能上に直接構築されており、サードパーティ依存関係は追加されません。

## グローバルショートカットマネージャーへのアクセス

マネージャーは、アプリケーションインスタンスの `GlobalShortcut` プロパティから利用できます。

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## ショートカットの登録

`Register` は、アクセラレーターとコールバックを受け取ります。ショートカットが押されるたびに、コールバックが専用の goroutine で実行されます。

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

`app.Run` を呼び出す前にショートカットを登録できます。その場合、アプリケーションの起動時にオペレーティングシステムとのバインドが自動的に行われます。

@note{type="tip" title="コールバックからの UI 操作"}
コールバックはメインスレッド外で実行されます。コールバックでウィンドウやその他の UI を操作する必要がある場合、ウィンドウメソッドによる処理は自動的に行われますが、独自のメインスレッド処理を実行する場合は `application.InvokeSync` でラップしてください。

@end

### アクセラレーターの形式

グローバルショートカットでは、メニューアクセラレーターやキーバインディングと同じアクセラレーター形式を使用します。

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` は、macOS では Command、Windows と Linux では Control として解釈されるため、クロスプラットフォームのショートカットに便利です。

## ショートカットの管理

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

登録済みのショートカットはすべて、アプリケーションの終了時に自動的に解放されるため、手動でクリーンアップする必要はありません。

## 同じショートカットを 2 回登録した場合の動作

これには 2 つの異なるケースがあり、Wails はそれぞれを別の方法で処理します。

### 同じアプリケーションがショートカットを 2 回登録する場合

このケースは Wails 自体によって処理され、すべてのプラットフォームで同じ動作になります。2 回目の `Register` 呼び出しはエラーを返し、元のバインドは維持されます（「エラーを返して維持」）。これにより動作の予測可能性が保たれ、有効なショートカットを暗黙的に置き換えることなく、誤りが明らかになります。

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

ショートカットのコールバックを変更するには、まず `Unregister` し、その後もう一度 `Register` してください。

### 別のアプリケーションがすでにショートカットを使用している場合

このケースはオペレーティングシステムによって決定されるため、結果はプラットフォームごとに異なります。

| プラットフォーム | 別のアプリケーションがショートカットを使用している場合の動作 |
| --- | --- |
| **macOS** | 登録は成功します。macOS では複数のアプリケーションが同じホットキーを登録できるため、登録は拒否されず、既存の所有者に加えてコールバックが追加されます。 |
| **Windows** | 登録は失敗し、`Register` はエラーを返します。最初にショートカットを登録したアプリケーションがそのショートカットを引き続き使用します。 |
| **Linux（X11）** | X サーバーが同じキーの組み合わせに対する 2 回目のグラブを拒否するため、登録は失敗し、`Register` はエラーを返します。 |
| **Linux（Wayland）** | コンポジターが調停します。通常、デスクトップのグローバルショートカットダイアログを通じて、ユーザーにバインドの承認または選択を求めます。 |

このような違いがあるため、`Register` が返すエラーを必ず確認し、ショートカットを確保できない場合は代替ショートカットまたはユーザーへのフィードバックを用意してください。

## プラットフォームに関する考慮事項

@tabs
[macOS]
グローバルショートカットでは、Carbon Event Manager のホットキー API を使用します。これは macOS でシステム全体のホットキーを実現する標準的な仕組みであり、アクセシビリティ権限は必要ありません。

ホットキーは物理的なキー位置にバインドされるため、QWERTY 以外の配列では、ショートカットは標準の ANSI/QWERTY 配列で同じ位置にあるキーに対応します。

@note{type="caution" title="非表示ショートカットと `ApplicationShouldTerminateAfterLastWindowClosed`"}
macOS では、`window.Hide()` は `orderOut:` を使用してウィンドウを非表示にします。AppKit は最後の非表示ウィンドウを閉じたものとして扱うため、`Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` を設定し、唯一のウィンドウをグローバルショートカットで非表示にすると、アプリケーションはバックグラウンドに残らず終了します。表示／非表示を切り替えるホットキーを使用する場合は、このオプションを未設定（デフォルト）のままにしてください。これにより、ウィンドウを非表示にした後で再び呼び出せます。

@end

[Windows]
グローバルショートカットでは、Win32 の `RegisterHotKey` API を使用します。自動リピートは抑制されるため、キーを押し続けてもコールバックは繰り返しではなく 1 回だけ実行されます。

別のアプリケーションがすでにそのキーの組み合わせを使用している場合、登録は失敗します。そのため、デフォルトにはあまり一般的でないキーの組み合わせを使用してください。

[Linux]
**X11** セッションでは、Wails は X サーバーからショートカットを直接取得するため、要求したアクセラレーターは指定どおりに割り当てられます。

**Wayland** セッションでは、設計上、アプリケーションがキーを直接取得する方法はありません。代わりに、Wails は XDG Desktop Portal の `org.freedesktop.portal.GlobalShortcuts` インターフェースを使用します。ポータルでは、渡したアクセラレーターは<em>優先</em>トリガーとして扱われ、最終的なキーの組み合わせはコンポジター（そして最終的にはユーザー）が決定します。ショートカットがアクティブになるとコールバックは引き続き実行されますが、実際のキーが要求したものと一致する保証はありません。また、`IsRegistered`/`GetAll` が返すのは、コンポジターが割り当てた内容ではなく、要求した内容です。

ポータルを使用するには、グローバルショートカットポータルを実装したデスクトップ環境（最近の GNOME や KDE Plasma など）が必要です。

@end

## 完全な例

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="重要なシステムショートカットを避ける"}
一部のキーの組み合わせは、オペレーティングシステムまたはデスクトップ環境によって予約されているため、アプリケーションは取得できません。競合しにくいデフォルトを選び、`Register` から返されるエラーを必ず処理してください。

@end
