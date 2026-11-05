---
title: "非公開 macOS API"
description: "非公開 macOS API に依存する Wails のすべての機能とオプション、および有効化用コマンドと公開 API のみのビルドで使用される代替処理。"
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails はデフォルトで macOS の公開 API を使用します。単一の Go ビルドタグ `private_mac_apis` を指定すると、このページに記載されている WebKit と AppKit の非公開 API 呼び出しが有効になります。公開されている Go のオプションとメソッドは、どちらのビルドでもすべて引き続き使用できます。このタグを指定しない場合、非公開 API 専用の操作は何も行わず、公開 API による代替手段がある機能はその手段を使用します。

@note{type="caution" title="macOS の非公開動作を有効にする"}
ウィンドウオプションを設定しても、非公開 API は有効になりません。有効にするには、ビルドコマンドに `private_mac_apis` を追加します。このタグが適用されるのは macOS デスクトップビルドのみで、iOS、Android、Windows、Linux、サーバービルドには適用されません。

@end

## 非公開 API を有効にする

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

本番環境向けに直接ビルドする場合は、`go build -tags production,private_mac_apis .` を使用します。フロントエンドのサンプルでは、実行前に各 README の手順に従ってバインディングとアセットをビルドしてください。カスタムまたは旧版の Taskfile では、`EXTRA_TAGS` を Go コンパイラーに渡す必要があります。

## 機能一覧

| 機能または値 | `private_mac_apis` で有効になるもの | タグを指定しない場合 |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | ネイティブウィンドウ上の透明な WKWebView | ネイティブウィンドウは構成されますが、WebView は不透明なままです |
| `Mac.Backdrop: MacBackdropTranslucent` | ネイティブのぼかし効果が透けて見える透明な WKWebView | 不透明な WebView の背後にぼかし効果が構成されます |
| `Mac.Backdrop: MacBackdropLiquidGlass` | 非公開 API による WebView の背景制御を使用した、ガラスレイヤー上の透明な WKWebView | 不透明な WebView の背後にガラスレイヤーが構成され、スタイル設定には公開 API による代替手段が使用されます |
| Liquid Glass のセットアップ時に WebView の背景をクリアする処理 | WebKit の非公開 `backgroundColor` コントロール | macOS 12 以降では公開 `underPageBackgroundColor`、それより古い macOS ではレイヤーの色を使用します。WebView は透明になりません |
| `app.Window.NewNotchWindow(...)` | ノッチ形状のパネル内の透明な WebView | 配置やアニメーションを含め、パネルは引き続き機能しますが、その WebView は不透明なままです |
| `Mac.LiquidGlass.Style` | 文書化されていないダークスタイル値を含む、Wails の既存のネイティブスタイルマッピング | 公開 API の regular/clear スタイルと light/dark の外観を使用します。下記の値の表を参照してください |
| `Mac.LiquidGlass.GroupID` | 空でない識別子に対して非公開 API によるガラスのグループ化を要求します | 無視され、グループ化は要求されません |
| `Mac.LiquidGlass.GroupSpacing` | 0 より大きい値に対して非公開 API によるグループ間隔を要求します | 無視されます |
| `window.OpenDevTools()` と JavaScript の `Window.OpenDevTools()` | macOS 12 以降で WebKit インスペクターをプログラムから開きます | 何も行いません |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | macOS 12 以降で、ウィンドウが初めて表示されたときにインスペクターをプログラムから開くよう要求します | 何も行いません |
| macOS 13.3 より前での従来のインスペクター有効化 | インスペクターのサポートが組み込まれている場合、WebKit の開発者向け拡張機能を有効にします | 何も行いません。Safari の公開 API による検査には macOS 13.3 以降が必要です |

## WebView の透明度と背景

**非公開 API が必要:** `MacBackdropTransparent`、`MacBackdropTranslucent`、`MacBackdropLiquidGlass`、およびノッチウィンドウで使用される WebView の透明化。内部では、Wails が WebKit の非公開キー `drawsBackground` を設定します。HTML または CSS の背景を透明にするだけでは、不透明なネイティブ WKWebView を透明にできません。

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

WebView の背景色を設定する非公開 API の操作では、WebKit のキー `backgroundColor` を使用します。Liquid Glass のセットアップでは、この操作を使用して WebView の背景をクリアします。タグを指定しない場合、この内部操作では macOS 12 以降で公開 `underPageBackgroundColor` を使用し、それより古い macOS ではビューのレイヤーを使用します。これらの代替手段では WebView は透明になりません。

`WebviewWindowOptions.BackgroundColour` と `window.SetBackgroundColour()` は macOS 上で <strong>ネイティブウィンドウ</strong>の色を設定するもので、それ自体は非公開 API を必要としません。同様に、`Frameless` と `Mac.TitleBar.AppearsTransparent` は AppKit の公開 API を使用します。非公開 API に依存するのは WebView の透明化であり、タイトルバーの透明化ではありません。macOS の背景効果には、`BackgroundType` だけに依存せず、`Mac.Backdrop` を構成してください。

[ウィンドウオプション](/features/windows/options/#mac-options)、[フレームレスウィンドウ](/features/windows/frameless/#with-transparent-background)、[ノッチウィンドウ](/features/windows/notch-windows/)を参照してください。

## Liquid Glass の値

ネイティブの `NSGlassEffectView` が利用可能な場合（macOS 26 以降）、次のマッピングが適用されます。文書化されているネイティブスタイル値は、`0`（regular）と `1`（clear）のみです。Go の定数は、どちらのビルドでも既存の値を維持します。

| `MacLiquidGlassStyle` の値 | `private_mac_apis` を指定した場合 | タグを指定しない場合 |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic`（`0`） | ネイティブの regular スタイル（`0`） | ネイティブの regular スタイル（`0`） |
| `LiquidGlassStyleLight`（`1`） | 既存のネイティブスタイルマッピング（`1`、clear） | Aqua 外観を使用したネイティブの regular スタイル（`0`） |
| `LiquidGlassStyleDark`（`2`） | **文書化されていないネイティブスタイル値 `2`** | Dark Aqua アピアランスを使用するネイティブの標準スタイル（`0`） |
| `LiquidGlassStyleVibrant`（`3`） | 既存のライト／ネイティブのクリアスタイル（`1`）にマッピングされる | ネイティブのクリアスタイル（`1`） |

Automatic と Vibrant は文書化されたネイティブスタイル値を使用しますが、ウィンドウ全体に Liquid Glass の背景を表示するには、<strong>WebView の透明化</strong>用のタグが引き続き必要です。Light の外観は、2 つのビルドで異なります。プライベートスタイルのマッピングは、将来の macOS バージョンでも同じ効果が描画されることを保証しません。

**常にプライベート API が必要：**`GroupID`と`GroupSpacing`。Wails はグループ化を要求する前に、プライベートな `setGroupIdentifier:`、`setGroupName:`、`setGroupSpacing:` セレクターの有無を確認します。タグがなければ、これらの操作は何も行いません。タグを有効にしても、実行中の macOS バージョンがこれらのセレクターをサポートするとは限りません。

`MacLiquidGlass.Material`、`CornerRadius`、`TintColor`自体には、プライベート API は必要ありません。ネイティブの Liquid Glass をサポートしない macOS バージョンでは、Wails は半透明のフォールバックを設定します。ただし、それを WebView 越しに表示するには、引き続きタグが必要です。

## Web Inspector

<strong>プライベート API が必要：</strong>アプリケーションから WebKit のインスペクターを開くための `OpenDevTools()` の呼び出し、または `OpenInspectorOnStartup: true` の設定。Wails はプライベートな `_inspector` セレクターを使用します。`private_mac_apis`がなければ、これらの操作は通知なく何も行いません。

インスペクターのサポートもビルドに組み込む必要があります。既存の `production` タグと `devtools` タグの意味は変わりません。

| ビルドタグ | プログラムによるインスペクターの起動 | macOS 13.3 以降での Safari による公開インスペクション |
| --- | --- | --- |
| なし | 何もしない | 有効 |
| `private_mac_apis` | macOS 12 以降で有効 | 有効 |
| `production` | 何もしない | 無効 |
| `production,private_mac_apis` | 何もしない | 無効 |
| `production,devtools` | 何もしない | 有効 |
| `production,devtools,private_mac_apis` | macOS 12 以降で有効 | 有効 |

macOS 13.3 以降では、Wails は公開 API の `WKWebView.inspectable` を使用して Safari によるインスペクションを有効にします。これにはプライベート API は必要ありません。macOS 13.3 より前では、フォールバックがプライベートな `developerExtrasEnabled` 設定を使用するため、ビルドに組み込まれたインスペクターサポートに加えて `private_mac_apis` も必要です。

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## この一覧の保守

プライベートなネイティブ呼び出しはすべて `v3/pkg/application/mac_private_api_darwin.go` に分離されており、デフォルトビルドでは `mac_public_api_darwin.go` が選択されます。この一覧には、透明化、WebView の背景色、Glass スタイル、Glass のグループ化、インスペクターの起動、従来のインスペクターの有効化が含まれます。これらの実装を変更する場合は、このページと影響を受けるオプションのドキュメントを併せて更新してください。
