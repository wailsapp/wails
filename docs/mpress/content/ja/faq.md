---
title: "よくある質問"
description: "Wails v3を使用したアプリケーション開発に関する、よくある質問への回答"
slug: "faq"
sourcePath: "faq.md"
---

## 一般

### Wails とは何ですか？

Wails は、Go と Web 技術を使用してデスクトップアプリケーションを構築するためのフレームワークです。アプリケーションロジックは Go で記述し、インターフェースは HTML、CSS、JavaScript（または任意のフロントエンドフレームワーク）で構築します。Wails はそれをオペレーティングシステムのネイティブ WebView でレンダリングします。その結果、ブラウザーを同梱せず、メモリ使用量が少なく、通常は約 10MB の単一バイナリで構成される、小さく高速でネイティブらしいアプリケーションになります。

### Wails はどのプラットフォームをサポートしていますか？

| プラットフォーム | 要件 |
| --- | --- |
| Windows | AMD64 および ARM64。[WebView2 ランタイム](https://developer.microsoft.com/microsoft-edge/webview2/)を使用します。 |
| macOS | Intel では 10.15 以降（アプリケーションのターゲットは 10.13 以降に設定可能）、Apple Silicon では 11.0 以降。ユニバーサルバイナリもサポートしています。 |
| Linux | AMD64 および ARM64。デフォルトのスタックは GTK4 と WebKitGTK 6.0 です（Ubuntu 24.04 以降、Debian 13 以降、Fedora 40 以降、および同等のディストリビューション）。Ubuntu 22.04、Debian 12、RHEL 9 など、WebKit2GTK 4.1 のみを提供するディストリビューションは、レガシーの `-tags gtk3` ビルド（v3.1 まで提供）でサポートされます。WebKit2GTK 4.0 のみを提供するディストリビューションはサポートされません。[Linux ビルドガイド](/guides/build/linux/)を参照してください。 |
| iOS および Android | 試験的サポートです。[モバイルガイド](/guides/mobile/)を参照してください。 |

[サーバービルド](/guides/server-build/)を使用して、アプリケーションを通常の Web アプリとして配信することもできます。

システムを確認し、プラットフォーム固有のインストール手順を表示するには、いつでも `wails3 doctor` を実行してください。

### 使用を開始するには何が必要ですか？

- Go 1.26 以降
- Node.js と npm（フロントエンドのビルド用）
- プラットフォームのツールチェーン：Windows では WebView2（10/11 にプリインストール済み）、macOS では Xcode Command Line Tools、Linux では `gcc` と GTK/WebKit 開発パッケージ

`wails3 doctor` はこれらをすべて確認し、不足しているものを正確に通知します。手順全体については、[インストール](/quick-start/installation/)を参照してください。

### Wails v3 は本番環境で使用できますか？

Wails v3.0.0 は、Windows、macOS、Linux 向けデスクトップアプリケーションの安定版リリースです。Android と iOS のサポートは引き続き試験的であり、安定版デスクトップの互換性保証の対象外です。[プロジェクトステータスページ](/status/)を参照してください。

## 開発

### Go の知識は必要ですか？

Go の基礎知識があると役立ちますが、専門家である必要はありません。アプリケーションロジックは通常の Go メソッドに記述し、それ以外は[チュートリアル](/tutorials/overview/)で一通り学べます。多くの開発者は、初めての Wails アプリケーションを構築しながら Go を習得しています。

### 好みのフロントエンドフレームワークを使用できますか？

はい。HTML、CSS、JavaScript にビルドできるものであれば、Wails で動作します。React、Vue、Svelte、バニラ JavaScript 用のテンプレート（それぞれに TypeScript 版もあります）が用意されており、その他のフレームワークも数分で接続できます。[フロントエンドフレームワーク](/guides/dev/frontend-frameworks/)を参照してください。

### JavaScript から Go 関数を呼び出すにはどうすればよいですか？

サービスを登録すると、Wails が型付きバインディングを生成します：

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

バインディングは `wails3 dev` の実行中に自動的に再生成されます。また、`wails3 generate bindings` を使用して必要に応じて生成することもできます。[サービス](/features/bindings/services/)を参照してください。

### TypeScript を使用できますか？

はい。バインディングジェネレーターはサービスとその型の TypeScript 定義を生成するため、Go の呼び出しには完全な型情報が付きます。

### Go と JavaScript の間でイベントを送信するにはどうすればよいですか？

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

イベント名は完全に一致している必要があります。[イベントリファレンス](/guides/events-reference/)を参照してください。

### アプリケーションをデバッグするにはどうすればよいですか？

`wails3 dev` を実行し、ウィンドウ内を右クリックすると、Web の場合とまったく同じようにブラウザーの開発者ツールを開けます。開発サーバーはフロントエンドのホットリロードにも対応しています。[デバッグ](/guides/dev/debugging/)を参照してください。

## ビルドと配布

### 本番用にビルドするにはどうすればよいですか？

```bash
wails3 build
```

バイナリは `bin/` に出力されます。本番ビルドには適切なデフォルト設定（ビルドタグ、`-trimpath`、シンボルの除去）がすでに適用されるため、軽量なバイナリを作成するための追加フラグは不要です。

### クロスコンパイルできますか？

制限付きで可能です。各プラットフォームでネイティブの WebView ライブラリを使用するため、純粋な Go のクロスコンパイルは適用できませんが、一般的なケースは十分にサポートされています：

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

別の OS から Linux 向けにビルドする場合は、Docker ベースのツールチェーンを使用します。対応状況の全一覧については、[クロスプラットフォームビルド](/guides/build/cross-platform/)を参照してください。

### インストーラーまたはパッケージを作成するにはどうすればよいですか？

```bash
wails3 package
```

これにより、各プラットフォームのネイティブ形式が生成されます。[インストーラーガイド](/guides/installers/)では、Windows の NSIS、macOS の `.app` バンドルと DMG、Linux パッケージについて説明しています。

### アプリケーションにコード署名するにはどうすればよいですか？

Windows と macOS の署名（公証を含む）は、[署名ガイド](/guides/build/signing/)で手順を追って説明しています。

## 機能

### 複数のウィンドウを作成できますか？

はい。v3 は複数ウィンドウをネイティブでサポートしています：

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

[複数ウィンドウ](/features/windows/multiple/)を参照してください。

### Wails はシステムトレイをサポートしていますか？

はい。メニューとクリックハンドラーもサポートしています：

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

[システムトレイ](/features/menus/systray/)を参照してください。

### ネイティブダイアログを使用できますか？

はい。ファイルダイアログ、メッセージダイアログ、質問ダイアログはすべてネイティブ実装を使用します：

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

[ダイアログ](/features/dialogs/overview/)を参照してください。

### Wails は自動更新をサポートしていますか？

はい。Wails v3 には、GitHub Releases、keygen.sh、Sparkle AppCast 用の差し替え可能なプロバイダー、暗号署名の検証、テーマを変更または置き換えられるデフォルト UI を備えた組み込みのセルフアップデーター（`app.Updater`）が含まれています。[アプリ内アップデーター](/guides/updater/)ガイドと[自己更新型 Wails アプリ](/tutorials/04-self-update-a-wails-app/)チュートリアルを参照してください。

## トラブルシューティング

### 正常に動作しません。まず何を確認すればよいですか？

```bash
wails3 doctor
```

ツールチェーンを検証し、不足している依存関係とそのインストールコマンドを一覧表示したうえで、バグ報告に含めるべきバージョン情報を出力します。

### ビルドに失敗します

一般的な対処方法を順に示します：

1. `go mod tidy`
2. `cd frontend && npm install`（`node_modules` がないことが最も一般的な原因です）
3. CLI を更新します：`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. Linux では、GTK/WebKit パッケージが不足していないか `wails3 doctor` で確認してください

### バインディングが見つからない、または古い

```bash
wails3 generate bindings
```

開発モードではバインディングが自動的に再生成されます。`wails3 dev`以外で新しいサービスを追加した場合や、メソッドシグネチャを変更した場合は、手動で再生成してください。

### イベントが発火しない

Go の`app.Event.Emit("name", ...)`と JavaScript の`Events.On("name", ...)`で、イベント名が完全に一致している必要があります。まず、入力ミスと大文字・小文字の違いを確認してください。

### バグを見つけた

[Issue を作成](https://github.com/wailsapp/wails/issues)し、`wails3 doctor`の出力を記載してください。対応しやすい報告の書き方については、[フィードバックガイド](/feedback/)で説明しています。

## v2からの移行

### v2からv3に移行すべきですか？

v3では、マルチウィンドウ対応、より簡潔なサービスベースの API、組み込みアップデーター、はるかに柔軟なビルドシステム、パフォーマンスの向上が実現されています。新規プロジェクトではv3を使用してください。既存のプロジェクトについては、[移行ガイド](/migration/v2-to-v3/)で相違点を順に説明しています。

### v2は今後もメンテナンスされますか？

はい。Wails v2には、引き続き重大な問題に対する更新が提供されます。新規プロジェクトでは Wails v3.0.0を使用してください。

### v2とv3を同じ環境で併用できますか？

はい。CLI は別々のバイナリ（`wails`と`wails3`）であり、モジュールのインポートパスも異なるため、メジャーバージョンが異なるプロジェクトを同じマシン上で問題なく共存させられます。

## コミュニティ

### サポートを受けるにはどうすればよいですか？

- 簡単な質問やディスカッションには[Discord](https://discord.gg/JDdSxwjhGf)
- 長文の質問には[GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- バグ報告には[GitHub Issues](https://github.com/wailsapp/wails/issues)

### コントリビュートするにはどうすればよいですか？

[コントリビューションガイド](/contributing/)を参照してください。バグ修正はいつでも歓迎します。新機能や公開されている動作の変更には、[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)のドラフト PR を使用します。Discord または GitHub Discussions での非公式な議論は任意です。

### サンプルはどこにありますか？

リポジトリには、ウィンドウ、ダイアログ、イベント、システムトレイ、サービスなどを扱う、実行可能なサンプルが60以上含まれています：[v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。

## ほかにも質問がありますか？

[Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[ディスカッションを作成](https://github.com/wailsapp/wails/discussions)してください。
