---
title: "フレームレスウィンドウ"
description: "フレームレスウィンドウでカスタムウィンドウクロームを作成する"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## フレームレスウィンドウ

Wails は、CSS ベースのドラッグ領域と各プラットフォームにネイティブな動作を備えた<strong>フレームレスウィンドウのサポート</strong>を提供します。プラットフォーム標準のタイトルバーを取り除くことで、ドラッグ、サイズ変更、システムコントロールなどの必須機能を維持しながら、ウィンドウクローム、独自デザイン、固有のユーザー体験を完全に制御できます。

![macOS ネイティブの角丸を使用したフレームレスウィンドウとして動作する、デフォルトの Wails v3 TypeScript スターターアプリ](/assets/screenshots/frameless-v3-native-corners-macos.png)

上の例は、`Frameless: true`を有効にしたデフォルトの Wails v3 TypeScript スターターアプリです。

## クイックスタート

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**ドラッグ可能なタイトルバー用の CSS：**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML：**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

**これだけです！** カスタムタイトルバーが完成しました。

## フレームレスウィンドウの作成

### 角の半径（macOS）

フレームレスウィンドウでは、デフォルトで AppKit 標準の macOS の角丸が維持されます。独自の半径（ポイント単位）を使用するには、`Mac.CornerRadius`を設定します：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

角を直角にするには、`Mac.CornerType`を`MacWindowCornerTypeSquare`に設定します。この場合、`CornerRadius`は無視されます：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### 基本的なフレームレスウィンドウ

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**得られるもの：**

- タイトルバーなし
- ウィンドウ枠なし
- システムボタンなし
- 透明な背景（任意）

**実装が必要なもの：**

- ドラッグ可能な領域
- 閉じる／最小化／最大化ボタン
- サイズ変更ハンドル（サイズ変更可能な場合）

### 透明な背景を使用する

<strong>macOS のプライベート API：</strong>WebView を透明にするには、`Mac.Backdrop: application.MacBackdropTransparent`を設定し、`-tags private_mac_apis`を指定してビルドします。このタグがない場合、HTML/CSS の背景が透明でも、ネイティブ WebView は不透明なままです。`Frameless`と`TitleBar.AppearsTransparent`自体はパブリック API を使用します。[macOS のプライベート API](/guides/build/private-macos-apis/#webview-transparency-and-background)を参照してください。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**ユースケース：**

- 角丸
- カスタム形状
- オーバーレイウィンドウ
- スプラッシュスクリーン

## ドラッグ領域

### CSS ベースのドラッグ

`--wails-draggable` CSS プロパティを使用します：

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**値：**

- `drag` - ドラッグ可能な領域
- `no-drag` - ドラッグできない領域（親要素がドラッグ可能な場合も含む）

### タイトルバーの完全な例

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**ボタン用の JavaScript：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Windows のネイティブ非クライアント領域

Windows では、カスタムタイトルバーの一部をネイティブの非クライアント領域として扱えます。これにより、任意の HTML/CSS デザインでタイトルバーとキャプションボタンを描画しながら、Windows ネイティブの動作を維持できます。キャプション領域ではウィンドウをドラッグでき、最大化ボタンでは Windows 11 の Snap Assist / Snap Layouts を表示でき、最小化、最大化、閉じるボタンではネイティブのヒットテストとマウス状態を利用できます。

以下の動画では、Windows ネイティブのヒットテストを使用するカスタム HTML/CSS タイトルバーを紹介しています。カスタム最大化ボタンでの Windows 11 の Snap Assist / Snap Layouts も含まれます。

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails は、Windows 固有の仕組みを 2 つサポートしています：

- `app-region`：WebView2 のネイティブ非クライアント領域サポートを使用
- `--wails-non-client-region`：カスタムキャプションボタン用の Wails ランタイムによる追跡を使用

### モードの選択

@note{type="caution" title="試験的機能"}
`WebView2CompositionHosting`は、内部でのウィンドウによる WebView2 のホスト方法と操作方法を変更します。Wails は、デフォルトの HWND ホスト型 WebView2 コントローラーの代わりにコンポジションコントローラーによるホストを使用し、入力を明示的に転送します。このモードでは、レンダリング、入力、フォーカス、または WebView2 Runtime の互換性に問題が生じる可能性があります。ネイティブ動作をするカスタムキャプションボタンが必要な場合にのみ有効にし、サポート対象の Windows および WebView2 Runtime の各バージョンでアプリを十分にテストしてください。

@end

Windows で必要な機能に応じて選択してください：

- WebView2 の`app-region: drag`と`app-region: no-drag`を使用して、アプリをシンプルにネイティブドラッグするには、`NonClientRegionSupport`を使用します。
- カスタムの最小化、最大化、閉じるボタンを Windows ネイティブのキャプションボタンと同様に動作させるには、`WebView2CompositionHosting`を使用します。
- 同じウィンドウで WebView2 ネイティブの`app-region`サポートと、Wails が管理するカスタムキャプションボタン領域の両方が必要な場合は、両方を有効にします。

`NonClientRegionSupport`は、Wails の`--wails-draggable`による追跡に代わる軽量なネイティブ方式です。ドラッグ可能な領域とドラッグできない領域を CSS で指定すると、WebView2 がキャプションに属するピクセルを判定し、Wails はヒットテスト時に WebView2 へネイティブ領域を問い合わせます。

現時点で、このモードが提供する機能は以上です。カスタムの最小化、最大化、閉じるボタンを Windows ネイティブのキャプションボタンと同様に動作させるものではなく、カスタム最大化ボタンで Windows 11 の Snap Assist / Snap Layouts を有効にするものでもありません。`--wails-draggable`の追加機構を使わず、アプリをシンプルにネイティブドラッグする必要がある場合に使用してください。

`WebView2CompositionHosting` は、ネイティブに動作するカスタムキャプションボタン用です。Wails は、`--wails-non-client-region` でマークされた DOM の矩形領域を追跡し、それらを `HTMINBUTTON`、`HTMAXBUTTON`、`HTCLOSE` などの Windows ヒットテスト値に対応付け、マウス入力をコンポジションホスト型の WebView2 サーフェスへ転送します。これにより、任意のビジュアルデザインを維持しながら、カスタム最大化ボタンを Windows 11 の Snap Assist / Snap Layouts に対応させることができます。

言い換えると、`NonClientRegionSupport` は WebView2 ネイティブの CSS リージョンサポートです。`WebView2CompositionHosting` では、ホスト側が所有するコンポジションとカスタムの非クライアント領域ヒットテストを Wails が担います。

### WebView2 app-region

ウィンドウで WebView2 ネイティブの非クライアント領域サポートを有効にします。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

次に、ドラッグ可能な領域を CSS の `app-region` プロパティでマークします。

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

ネイティブのキャプションドラッグだけが必要で、タイトルバーのコントロールを通常のフロントエンドクリックで処理する場合に使用します。

このモードは、WebView2 自体の非クライアント領域サポートに制約されます。現在の WebView2 リリースでは、ドラッグ領域と非ドラッグ領域のみを使用できます。ネイティブの最小化、最大化、閉じるという個別の役割を持つ、完全にカスタム化されたフロントエンドのキャプションボタンを表現するためのものではありません。

### ネイティブに動作するカスタムキャプションボタン

システムのキャプションボタンと同様に動作させるカスタムの最小化、最大化、閉じるボタンでは、コンポジションホスティングを有効にします。

@note{type="caution" title="試験的機能"}
`WebView2CompositionHosting` は、DirectComposition を使用した WebView2 コンポジションコントローラーのホスティングを使用します。有効にする前に、[モードの選択](#heading-7)を参照してください。

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

次に、フロントエンドの各領域を `--wails-non-client-region` でマークします。

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

サポートされる `--wails-non-client-region` の値：

- `caption` - ドラッグ可能なキャプション領域
- `minimize` - ネイティブ最小化ボタンのヒットターゲット
- `maximize` - Windows 11 の Snap Assist / Snap Layouts のホバー動作を含む、ネイティブ最大化ボタンのヒットターゲット
- `close` - ネイティブの閉じるボタンのヒットターゲット

Wails ランタイムは DOM、スタイル、サイズ、スクロール、ビューポートの変更を監視し、領域のスナップショットをネイティブウィンドウへ送信します。領域のジオメトリは CSS ピクセル単位で測定され、Windows のヒットテスト用に物理ピクセルへ変換されます。

ビジュアルデザインはすべて自由に決められます。領域は、各矩形が何を意味するかを Windows に伝えるだけです。ボタンの形状、アイコン、色、間隔、ホバースタイル、レイアウトは、引き続きフロントエンドで定義します。

### 両方の組み合わせ

同じウィンドウで WebView2 の `app-region` サポートと Wails が管理するキャプションボタン領域を使用する場合は、両方のオプションを有効にできます。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## システムボタン

### 閉じる／最小化／最大化の実装

**Go 側：**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**JavaScript 側：**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**または、ランタイムメソッドを使用します：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### 最大化状態の切り替え

ボタンアイコン用に最大化状態を追跡します。

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## リサイズハンドル

### CSS ベースのリサイズ

Wails は、フレームレスウィンドウ用の自動リサイズハンドルを提供します。

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**値：**

- `all` - すべての辺からリサイズ
- `top`、`bottom`、`left`、`right` - 特定の辺
- `top-left`、`top-right`、`bottom-left`、`bottom-right` - 四隅
- `none` - リサイズなし

### リサイズハンドルの例

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## プラットフォーム固有の動作

@tabs{sync-key="platform"}
[Windows]
**Windows のフレームレスウィンドウ：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**機能：**

- 自動ドロップシャドウ
- Snap Layouts のサポート（Windows 11）
- Aero Snap のサポート
- DPI スケーリング

**装飾を無効化：**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**Snap Assist：**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

これは、Windows のホットキー経由で Snap Layouts を起動します。ホバー時にネイティブの Snap Layouts を表示するカスタム HTML 最大化ボタンには、代わりに[Windows のネイティブ非クライアント領域](#windows-)を使用してください。

**カスタムタイトルバーの高さ：** Windows は CSS からドラッグ領域を自動的に検出します。

[macOS]
**macOS のフレームレスウィンドウ：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**機能：**

- ネイティブフルスクリーンのサポート
- 信号機ボタン（オプション）
- すりガラス効果
- 透明なタイトルバー

**タイトルバーを完全に非表示にする**（`application` パッケージからエクスポートされているプリセットバリアントを使用してください。`TitleBarStyle` フィールドや `MacTitleBarStyleHidden` 定数はありません）:

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

その他のプリセットには、`MacTitleBarDefault`、`MacTitleBarHiddenInset`、`MacTitleBarHiddenInsetUnified` があります。

**非表示のタイトルバー:** タイトルバーを非表示にしたままドラッグできます。これは、ウィンドウがフレームレスであるか、`AppearsTransparent` を使用している場合にのみ有効です:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Linux のフレームレスウィンドウ:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**機能:**

- 基本的なフレームレス対応
- CSS ドラッグ領域
- デスクトップ環境によって異なる

**デスクトップ環境に関する注意事項:**

- **GNOME:** 良好な対応
- **KDE Plasma:** 良好な対応
- **XFCE:** 基本的な対応
- **タイル型ウィンドウマネージャー:** 限定的な対応

**コンポジターが必要:** 透過にはコンポジターが必要です（最新のデスクトップ環境のほとんどには備わっています）。

@end

## 一般的なパターン

### パターン 1: モダンなタイトルバー

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### パターン 2: スプラッシュスクリーン

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### パターン 3: 角丸ウィンドウ

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### パターン 4: オーバーレイウィンドウ

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## 完全な例

以下は、本番環境で使用できるフレームレスウィンドウです:

**Go:**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS:**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript:**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## ベストプラクティス

### ✅ 推奨事項

- **ドラッグ可能な領域を用意する** - ユーザーがウィンドウを移動するために必要です
- **システムボタンを実装する** - 閉じる、最小化、最大化
- **最小サイズを設定する** - 使用できないレイアウトになることを防ぎます
- **すべてのプラットフォームでテストする** - 動作が異なります
- **ドラッグ領域には CSS を使用する** - 柔軟で保守しやすくなります
- **視覚的なフィードバックを提供する** - ボタンにホバー状態を設けます

### ❌ 禁止事項

- **サイズ変更ハンドルを忘れない** - ウィンドウのサイズを変更可能にする場合は必要です
- **ウィンドウ全体をドラッグ可能にしない** - 操作できなくなります
- **ボタンをドラッグ対象外にすることを忘れない** - 忘れるとボタンが機能しません
- **ドラッグ領域を小さくしすぎない** - つかみにくくなります
- **プラットフォーム間の違いを忘れない** - 十分にテストしてください

## トラブルシューティング

### ウィンドウをドラッグできない

**原因:** `--wails-draggable: drag` が指定されていない

**解決策:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### ボタンが機能しない

**原因:** ボタンがドラッグ可能な領域内にある

**解決策:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### ウィンドウのサイズを変更できない

**原因:** サイズ変更ハンドルがない

**解決策:**

```css
body {
    --wails-resize: all;
}
```

## 次のステップ

@cards{cols="2"}
▣ ウィンドウの基本
ウィンドウ管理の基礎を学びます。

[詳しく見る →](/features/windows/basics/)

---
⚙ ウィンドウオプション
ウィンドウオプションの完全なリファレンスです。

[詳しく見る →](/features/windows/options/)

---
🚀 ウィンドウイベント
ウィンドウのライフサイクルイベントを処理します。

[詳しく見る →](/features/windows/events/)

---
◆ 複数ウィンドウ
マルチウィンドウアプリケーションのパターンです。

[詳しく見る →](/features/windows/multiple/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[フレームレスのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless)を確認してください。
