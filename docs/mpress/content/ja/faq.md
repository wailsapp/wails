---
title: "よくある質問"
description: "Wails v3 を使用したアプリケーション開発に関する、よくある質問への回答"
slug: "faq"
sourcePath: "faq.md"
---

## 概要

### Wails とは何ですか？

Wails は、Go と Web 技術を使用してデスクトップアプリケーションを構築するためのフレームワークです。アプリケーションロジックは Go で記述し、インターフェースは HTML、CSS、JavaScript（または任意のフロントエンドフレームワーク）で構築します。Wails はそれをオペレーティングシステムのネイティブ WebView でレンダリングします。その結果、ブラウザーを同梱せず、メモリ使用量が少なく、通常は約 10MB の単一バイナリで構成される、小型で高速かつネイティブらしいアプリケーションを実現できます。

### Wails はどのプラットフォームをサポートしていますか？

| プラットフォーム | 要件 |
| --- | --- |
| Windows | AMD64 および ARM64。[WebView2 ランタイム](https://developer.microsoft.com/microsoft-edge/webview2/)を使用します。 |
| macOS | Intel では 10.15 以降（アプリケーションのターゲットは 10.13 以降に設定可能）、Apple Silicon では 11.0 以降。ユニバーサルバイナリをサポートしています。 |
| Linux | AMD64 および ARM64。デフォルトのスタックは、WebKitGTK 6.0 を使用する GTK4 です（Ubuntu 24.04 以降、Debian 13 以降、Fedora 40 以降、および同等のディストリビューション）。Ubuntu 22.04、Debian 12、RHEL 9 など、WebKit2GTK 4.1 のみを提供するディストリビューションは、従来の `-tags gtk3` ビルドでサポートされます（v3.1 まで利用可能）。WebKit2GTK 4.0 のみを提供するディストリビューションはサポートされません。[Linux ビルドガイド](/guides/build/linux/)を参照してください。 |
| iOS および Android | 試験的サポートです。[モバイルガイド](/guides/mobile/)を参照してください。 |

[サーバービルド](/guides/server-build/)を使用して、アプリケーションを通常の Web アプリとして提供することもできます。

システムを確認し、プラットフォーム固有のインストール手順を表示するには、いつでも `wails3 doctor` を実行してください。

### 開始するには何が必要ですか？

- Go 1.25 以降
- Node.js および npm（フロントエンドのビルド用）
- プラットフォーム別のツールチェーン：Windows では WebView2（10/11 にプリインストール済み）、macOS では Xcode Command Line Tools、Linux では `gcc` と GTK/WebKit 開発パッケージ

`wails3 doctor` を実行すると、これらすべてが確認され、不足しているものが正確に示されます。詳しい手順については、[インストール](/quick-start/installation/)を参照してください。

### Wails v3 は本番環境で使用できますか？

Wails v3 は、安定したデスクトップ API を備えたベータ版ソフトウェアです。すでに本番環境で稼働しているアプリケーションもありますが、3.0 に向けた最終調整が完了するまでは、デプロイ前に十分なテストを行うことを推奨します。現在の状況については、[プロジェクトのステータスページ](/status/)を参照してください。Wails v2 は現在の安定版であり、引き続き修正が提供されます。

## 開発

### Go の知識は必要ですか？

Go の基礎知識があると役立ちますが、専門家である必要はありません。アプリケーションロジックは通常の Go メソッドで記述し、それ以外のすべては[チュートリアル](/tutorials/overview/)で順を追って説明します。多くの開発者は、最初の Wails アプリケーションを構築しながら Go を習得しています。

### 好みのフロントエンドフレームワークを使用できますか？

はい。HTML、CSS、JavaScript にビルドできるものであれば、Wails で使用できます。React、Vue、Svelte、プレーン JavaScript 用のテンプレートが用意されており、それぞれに TypeScript 版もあります。その他のフレームワークも数分で連携できます。[フロントエンドフレームワーク](/guides/dev/frontend-frameworks/)を参照してください。

### JavaScript から Go 関数を呼び出すにはどうすればよいですか？

サービスを登録すると、Wails が型付きバインディングを生成します。

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

バインディングは `wails3 dev` の実行中に自動的に再生成されます。また、`wails3 generate bindings` を使用して必要なときに再生成することもできます。[サービス](/features/bindings/services/)を参照してください。

### TypeScript を使用できますか？

はい。バインディングジェネレーターはサービスとその型の TypeScript 定義を生成するため、Go の呼び出しには完全な型情報が付与されます。

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

`wails3 dev` を実行し、ウィンドウ内を右クリックしてブラウザーの開発者ツールを開きます。Web の場合とまったく同じように使用できます。開発サーバーはフロントエンドのホットリロードにも対応しています。[デバッグ](/guides/dev/debugging/)を参照してください。

## ビルドと配布

### 本番環境向けにビルドするにはどうすればよいですか？

```bash
wails3 build
```

バイナリは `bin/` に出力されます。本番ビルドには適切なデフォルト設定（ビルドタグ、`-trimpath`、シンボルの除去）がすでに適用されるため、軽量なバイナリを生成するための追加フラグは不要です。

### クロスコンパイルできますか？

制限付きで可能です。各プラットフォームでネイティブ WebView ライブラリを使用するため、純粋な Go のクロスコンパイルは適用できませんが、一般的なケースは十分にサポートされています。

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

別の OS から Linux 向けにビルドする場合は、Docker ベースのツールチェーンを使用します。対応状況の一覧については、[クロスプラットフォームビルド](/guides/build/cross-platform/)を参照してください。

### インストーラーまたはパッケージを作成するにはどうすればよいですか？

```bash
wails3 package
```

これにより、プラットフォームのネイティブ形式が生成されます。[インストーラーガイド](/guides/installers/)では、Windows の NSIS、macOS の `.app` バンドルと DMG、および Linux パッケージについて説明しています。

### アプリケーションにコード署名するにはどうすればよいですか？

Windows と macOS の両方の署名（公証を含む）について、[署名ガイド](/guides/build/signing/)で手順を追って説明しています。

## 機能

### 複数のウィンドウを作成できますか？

はい。v3 はマルチウィンドウをネイティブにサポートしています。

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

[複数ウィンドウ](/features/windows/multiple/)を参照してください。

### Wails はシステムトレイをサポートしていますか？

はい。メニューとクリックハンドラーもサポートしています。

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

[システムトレイ](/features/menus/systray/)を参照してください。

### ネイティブダイアログを使用できますか？

はい。ファイルダイアログ、メッセージダイアログ、質問ダイアログは、いずれもネイティブ実装を使用します。

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

[ダイアログ](/features/dialogs/overview/)を参照してください。

### Wails は自動更新をサポートしていますか？

はい。Wails v3 には、GitHub Releases、keygen.sh、Sparkle AppCast 用の差し替え可能なプロバイダー、暗号署名の検証、テーマを変更または置き換え可能なデフォルト UI を備えた組み込みの自己更新機能（`app.Updater`）があります。[アプリ内アップデーター](/guides/updater/)ガイドと[自己更新する Wails アプリ](/tutorials/04-self-update-a-wails-app/)チュートリアルを参照してください。

## トラブルシューティング

### 正常に動作しません。何から確認すればよいですか？

```bash
wails3 doctor
```

ツールチェーンを検証し、不足している依存関係とそのインストールコマンドを一覧表示したうえで、バグ報告に含めるべきバージョン情報を出力します。

### ビルドに失敗します

一般的な解決方法を、試す順に示します。

1. `go mod tidy`
2. `cd frontend && npm install`（`node_modules`がないことが最も一般的な原因です）
3. CLI を更新します：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. Linux では、`wails3 doctor`を確認し、不足している GTK/WebKit パッケージがないか調べます

### バインディングがない、または古いままです

```bash
wails3 generate bindings
```

開発モードではバインディングが自動的に再生成されます。`wails3 dev`以外で新しいサービスを追加した場合やメソッドシグネチャを変更した場合は、手動で再生成してください。

### イベントが発生しません

Go の`app.Event.Emit("name", ...)`と JavaScript の`Events.On("name", ...)`で、イベント名が完全に一致している必要があります。まず、入力ミスや大文字と小文字の違いがないか確認してください。

### バグを見つけました

[Issue を作成](https://github.com/wailsapp/wails/issues)し、`wails3 doctor`の出力を含めてください。対応しやすい報告に必要な情報については、[フィードバックガイド](/feedback/)で説明しています。

## v2 からの移行

### v2 から v3 に移行すべきですか？

v3 では、マルチウィンドウのサポート、より整理されたサービスベースの API、組み込みアップデーター、柔軟性が大幅に向上したビルドシステム、パフォーマンスの改善が導入されています。新規プロジェクトでは v3 を使用してください。既存プロジェクトについては、[移行ガイド](/migration/v2-to-v3/)で相違点を順に説明しています。

### v2 は今後もメンテナンスされますか？

はい。v3 が安定版リリースに向けて進む間も、v2 には引き続き修正が提供されます。

### v2 と v3 を併用できますか？

はい。CLI はそれぞれ別のバイナリ（`wails`と`wails3`）であり、モジュールのインポートパスも異なるため、メジャーバージョンが異なるプロジェクトを同じマシン上で問題なく共存させられます。

## コミュニティ

### どこでサポートを受けられますか？

- 簡単な質問や意見交換には[Discord](https://discord.gg/JDdSxwjhGf)
- 詳しい質問には[GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- バグの報告には[GitHub Issues](https://github.com/wailsapp/wails/issues)

### 貢献するにはどうすればよいですか？

[コントリビューションガイド](/contributing/)を参照してください。バグ修正はいつでも歓迎します。新機能や公開されている動作の変更には、[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)のドラフト PR を使用します。Discord または GitHub Discussions での非公式な議論は任意です。

### サンプルはどこにありますか？

リポジトリには、ウィンドウ、ダイアログ、イベント、システムトレイ、サービスなどを扱う、実行可能なサンプルが60以上同梱されています：[v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples)。

## ほかに質問がありますか？

[Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[Discussion を作成](https://github.com/wailsapp/wails/discussions)してください。
