---
title: "Wails でウィンドウをカスタマイズする"
description: "Wails アプリケーションのウィンドウの外観と動作をカスタマイズします"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

対象プラットフォーム：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails は、ウィンドウコントロールの外観と機能を制御するための API を提供します。この機能は Windows と macOS で利用できますが、Linux では利用できません。

## ウィンドウボタンの状態を設定する

ボタンの状態は、`ButtonState` 列挙型で定義されます。

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`：ボタンは有効で、表示されます。
- `ButtonDisabled`：ボタンは表示されますが、無効です（グレー表示）。
- `ButtonHidden`：ボタンはタイトルバーに表示されません。

ボタンの状態は、ウィンドウの作成時または実行時に設定できます。

### ウィンドウ作成時にボタンの状態を設定する

新しいウィンドウを作成するときは、`WebviewWindowOptions` 構造体を使用してボタンの初期状態を設定できます。

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

上の例では、最小化ボタンは非表示、最大化ボタンは無効（グレー表示）、閉じるボタンは有効です。

### 実行時にボタンの状態を設定する

`Window` インターフェースの次のメソッドを使用して、実行時にボタンの状態を変更することもできます。

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS：MaximiseButtonState と FullscreenButtonState は同じボタンを共有する

macOS では、緑色の信号ボタン（`NSWindowZoomButton`）が、最大化とフルスクリーンの両方に使用される同一の物理コントロールです。ウィンドウ作成時に `MaximiseButtonState` と `FullscreenButtonState` に異なる値を設定すると、そのままでは警告なしに後から書き込まれた値で上書きされます。

これを避けるため、Wails は初期化時に、`ButtonEnabled` < `ButtonDisabled` < `ButtonHidden` の順序に従って、2 つの状態のうち<strong>より制限の厳しい方</strong>を適用します。

| `MaximiseButtonState` | `FullscreenButtonState` | macOS での実効状態 |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

実行時には、macOS 上の `SetMaximiseButtonState` と `SetFullscreenButtonState` はどちらも `NSWindowZoomButton` を対象とするため、最後の呼び出しが優先されます。

### プラットフォームによる違い

ボタン状態の機能は、Windows と macOS で動作が多少異なります。

|  | Windows | Mac |
| --- | --- | --- |
| 最小化／最大化／閉じるを無効化 | 最小化／最大化／閉じるを無効化 | 最小化／最大化／閉じるを無効化 |
| 最小化を非表示 | 最小化を無効化 | 最小化ボタンを非表示 |
| 最大化を非表示 | 最大化を無効化 | 最大化ボタンを非表示 |
| 閉じるを非表示 | すべてのコントロールを非表示 | 閉じるを非表示 |
| `FullscreenButtonState` | 何もしない | ズーム（緑色）ボタンが対象 |

注: Windows では、最小化ボタンと最大化ボタンを個別に非表示にすることはできません。 ただし、両方を無効にすると両方のコントロールが非表示になり、閉じるボタンだけが表示されます。Windows の標準タイトルバーには専用のフルスクリーンボタンがないため、`FullscreenButtonState` は Windows では何も行いません。

### ウィンドウスタイルの制御（Windows）

Windows でタイトルバーのスタイルを制御するには、`WebviewWindowOptions` 構造体の `ExStyle` フィールドを使用できます。

例:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

ウィンドウの拡張スタイルに影響する次のオプションは、この設定によって上書きされます。

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
