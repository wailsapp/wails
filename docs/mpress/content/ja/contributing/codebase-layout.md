---
title: "コードベースの構成"
description: "Wails v3リポジトリの構成と各要素の連携方法"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3は、フレームワークのランタイム、CLI、 サンプル、ドキュメント、ビルドツールチェーンを含む<strong>モノレポ</strong>です。 このページでは、内部構造を詳しく調べる際に重要となる<em>ディレクトリ構成</em>を説明します。

## 最上位の構成

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

ここからは、<strong>`v3/`</strong>ツリーを詳しく見ていきます。

## `v3/`ルート

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> プロジェクトテンプレートは`internal/templates/`以下に同梱されています（フレームワークごとに1つのフォルダーがあり、
>
> さらに`base/`、`_common/`、`ios/`があります）。最上位には`v3/templates/`
>
> ディレクトリはありません。

### 全体像

1. <strong>`pkg/`</strong>は、<em>アプリケーション開発者がインポートするもの</em>を公開します\
2. <strong>`internal/`</strong>には、<em>内部機能の実装</em>が含まれます\
3. <strong>`cmd/wails3`</strong>は、<em>プロジェクトのライフサイクルとビルド</em>を駆動します\

それ以外のすべては、この3本の柱を支えます。

---

## `cmd/` – コマンド

| パス | 注記 |
| --- | --- |
| `v3/cmd/wails3` | <strong>CLIのエントリーポイント</strong>です。小さな`main.go`が、すべてのロジックを`internal/commands`内のパッケージに委譲します。 |
| `internal/commands/*` | サブコマンド（init、dev、build、doctorなど）。見つけやすいよう、それぞれが個別のファイルに格納されています。 |
| `internal/commands/task_wrapper.go` | CLIフラグとTaskfileビルドパイプラインの橋渡しをします。 |

CLIは次の機能を担います。

- **プロジェクトのスキャフォールディング**（`init`、テンプレート生成）\
- **開発サーバーのオーケストレーション**（`dev`、ライブリロード）\
- **本番ビルドとパッケージング**（`build`、`package`、プラットフォーム別ラッパー）\
- **診断**（`doctor`）\

---

## `internal/` – 中核部分

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### 主要なサブパッケージ

| パッケージ | 役割 | 接続先 |
| --- | --- | --- |
| `runtime` | 小規模な`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`ビルドタグ用の連携コードと、`runtime/desktop/`以下の組み込みJSランタイムが格納されています。OSごとの実際のウィンドウ、クリップボード、ダイアログ、トレイのコードは`pkg/application/*_{darwin,linux,windows}.go`にあります。 | `pkg/application`を介して間接的にインポートされます。 |
| `assetserver` | デュアルモードのファイルサーバー：<br />• 開発時：ディスクから配信し、Viteをプロキシします（`build_dev.go`）<br />• 本番時：`go:embed`を介してアセットを埋め込みます（`build_production.go`） | 起動時に`pkg/application`によって初期化されます。 |
| `generator` | Goソースを解析して<strong>バインディングメタデータ</strong>を構築し、その後、TypeScript/JSスタブファイルとイベント定数を生成します。エントリーポイント：`collect/`と`render/`を基盤とする`generator.Generate` / `generator.Generator`。 | `wails3 generate bindings`によって実行されます。 |
| `packager` | Linux向けの`deb`/`rpm`/`archlinux`アーティファクトを生成するための`nfpm`ラッパーです（`internal/commands/`以下の`myapp.DEB`/`.RPM`/`.ARCHLINUX`用nfpm設定によって駆動されます）。 | `wails3 tool package`によって呼び出されます。macOS DMG / Windows MSIXは`internal/commands/{dmg,msix.go,webview2/}`以下にあります。 |

補助ユーティリティ（例：`s/`、`hash/`、`flags/`）により、内部の関心事が分離されています。

---

## `pkg/` – 公開API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> `pkg/runtime/`、`pkg/options/`、`pkg/menu/`というパッケージはありません。ウィンドウとメニューの
>
> オプションは`pkg/application`と同じ場所にあります（例：`WebviewWindowOptions`、`Menu`、
>
> `MenuItem`）。また、`assetserver/`は`internal/`以下にあります。

`pkg/application`はWailsプログラムを起動します。

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

内部では次の処理を行います。

1. `internal/runtime`ビルドタグ用の連携コードと、`pkg/application/`内のOS別コードを接続します
2. `internal/assetserver`インスタンスをセットアップします
3. バインディング駆動のメッセージプロセッサーを登録します
4. OSのメインスレッドに入ります

---

## `internal/templates/` – スキャフォールディングの設計図

`internal/templates/`には、**基本テンプレート**（`base/`、 `_common/`、`ios/`以下のGo構成）と、**フロントエンドスキン**（`vanilla[-ts]`、`react[-ts]`、 `react-swc[-ts]`、`lit[-ts]`、`preact[-ts]`、`qwik[-ts]`、`solid[-ts]`、 `svelte[-ts]`、`sveltekit[-ts]`、`vue[-ts]`）が同梱されています。

`wails3 init -t react`の実行時、CLIは次の処理を行います。

1. `_common`のGoファイルをコピーします
2. 目的のフロントエンドパックをマージします
3. `go mod tidy`を実行します（`--skipgomodtidy`でスキップ可能）

テンプレートを編集しても、<strong>既存のアプリには</strong>影響せず、今後の`init`にのみ反映されます。公開サンプルは`v3/examples/`以下にありますが、コントリビューター向けドキュメントで説明されている自動テストスイートの代わりにはなりません。

---

## `tasks/` – リリースの自動化

Taskfileは、複雑なクロスコンパイル、バージョン更新、変更履歴の生成をラップします。これらは`internal/commands/task.go`によってプログラムから利用されるため、同じロジックが<strong>CLI</strong>と<strong>CI</strong>の両方を支えています。

---

## 各要素の連携

```d2
direction: down
CLI: wails3 CLI
Generator: internal/generator
AssetDev: assetserver（開発）
Packager: internal/packager
AppRuntime: {
  label: アプリランタイム
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: OS API
}
CLI -> Generator: ビルド／生成
CLI -> AssetDev: 開発
CLI -> Packager: パッケージ化
Generator -> ApplicationPkg: バインディング
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

<em>CLI → generator → runtime</em>が、<strong>ソース</strong>から<strong>実行中のデスクトップアプリ</strong>に至る中心的な経路を形成します。

---

## 概要を把握するためのヒント

| 理解したいもの | 参照先 |
| --- | --- |
| プラットフォーム互換レイヤー | `pkg/application/*_darwin.go`、`*_linux.go`、`*_windows.go`（window、clipboard、dialogs、systray、mainthread、events_common）。Linuxのcgo：`pkg/application/linux_cgo*.go`。 |
| ブリッジプロトコル | `pkg/application/messageprocessor*.go` |
| アセットのワークフロー | `internal/assetserver/`（`build_dev.go`と`build_production.go`の比較） |
| パッケージ化のフロー | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`、`internal/packager/` |
| テンプレートエンジン | `internal/templates/`（`templates.Install`、`templates.GetDefaultTemplates`） |
| 静的解析 | `internal/generator/{generate.go,collect/,render/}` |

---

これで、リポジトリの<strong>全体像</strong>を把握できました。`ripgrep`、IDEの「ファイル／シンボルへ移動」機能、サンプルアプリを併用して、各機能をさらに詳しく調べてください。ハッキングを楽しんでください！
