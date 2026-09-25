---
title: "アセットサーバー"
description: "Wails v3 が開発環境と本番環境で Web アセットを配信し、埋め込む仕組み"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## 概要

すべての Wails アプリケーションは、次の要素を組み合わせた<strong>単一のネイティブ実行ファイル</strong>として配布されます。

1. *Go* バックエンド
2. *Web* フロントエンド（HTML + JS + CSS）

これを実現する橋渡し役が<strong>アセットサーバー</strong>です。Go のビルドタグによってコンパイル時に選択される、<strong>2 つの動作モード</strong>があります。

| モード | タグ | 目的 |
| --- | --- | --- |
| **開発** | `//go:build !production` | ホットリロードによる高速な反復開発 |
| **本番** | `//go:build production` | 外部依存関係のない埋め込みアセット |

実装は`v3/internal/assetserver/`にあり、ファイルごとに役割が明確に分けられています。

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## 開発モード

### ライフサイクル

1. `wails3 dev`が起動し、`build/Taskfile.yml`で定義されたタスク（通常は`npm run dev`）を実行して、**フロントエンド開発サーバーを生成します**（Vite、SvelteKit、React-SWC など）。
2. CLI は、`WAILS_VITE_PORT`を Wails の開発用ポートに設定し、`FRONTEND_DEVSERVER_URL`を、実行中のフレームワーク開発サーバーを指す<strong>完全な</strong> URL（`http://host:port` / `https://host:port`）に設定します。`internal/commands/dev.go`を参照してください。
3. 開発用アセットサーバー（`build_dev.go`の`//go:build !production`によって組み込まれます）は、`GetDevServerURL()`を介して`FRONTEND_DEVSERVER_URL`を読み取り、ランタイム以外のトラフィックをそこへリバースプロキシします。
4. 静的ファイル（`/assets/logo.svg`）は、高速化のため`asset_fileserver.go`を介して<strong>ディスクから直接配信</strong>できます。一方、不明なものはすべてフレームワークの開発サーバーへ<strong>プロキシされ</strong>、<em>即時</em>のホットモジュール置換を利用できます。

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### 機能

- **ライブリロード** — Vite、SvelteKit などは WebSocket 経由で HMR を挿入し、開発用アセットサーバーが透過的にプロキシします。
- **ソースマップ対応** — アセットはバンドルされないため、ブラウザーの開発者ツールでエラーを元のソースに対応付けられます。
- **Go の再コンパイル不要** — 再ビルドされるのはフロントエンドだけです。`.go`ファイルを変更するまで、Go コードは実行され続けます。

### フレームワークの切り替え

開発用プロキシは<strong>フレームワークに依存しません</strong>。Wails CLI は、開発タスクの起動時に 2 つの環境変数を公開します。

| 環境変数 | ソース | 意味 |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go`（`wailsVitePort`定数） | 既定の開発用ポート（`--port`を渡さない場合は9245）— Vite の設定ではこの値を使用することを推奨します |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | Wails がプロキシ先として使用する完全な URL。Go では`assetserver.GetDevServerURL()`（`build_dev.go`）を介して読み取ります |

v3 ツリーには、`VITE_PORT`、`FRONTEND_DEV_PORT`、`WAILSDEV_VERBOSE`という環境変数はありません。

新しいテンプレートを追加し、その開発タスクを定義すれば、アセットサーバーはそのまま動作します。

---

## 本番モード

`wails3 build`を実行すると、パイプラインは次の処理を行います。

1. フロントエンドの<strong>本番ビルド</strong>（`npm run build`）を実行し、`frontend/dist/**`を生成します。
2. アプリケーション独自のパッケージ（通常は`main.go`の隣にある`//go:embed all:frontend/dist`）の`go:embed`を介して、そのディレクトリをアプリに<strong>埋め込みます</strong>。
3. `-tags production`を指定して Go バイナリをコンパイルします（Taskfile ラッパーによって`EXTRA_TAGS`を介して渡されます）。

`internal/assetserver/build_production.go`は、本番用コードパスに切り替える ビルドタグ用スタブです。`internal/assetserver/bundled_assetserver.go`は <strong>手書き</strong>で、`bundledassets/`内のランタイム JS をラップしており、 生成ファイルではありません。

### リクエスト処理

実際のハンドラーは`internal/assetserver/assetserver.go` / `asset_fileserver.go`です。概念的には、次のように動作します。

1. リクエストされたパスにある埋め込み静的アセットを試します。
2. SPA ルーティングでは`index.html`にフォールバックします。
3. 拡張子が不明な場合は、コンテンツタイプを判定します（`content_type_sniffer.go`）。
4. 適切なキャッシュヘッダーを設定します。

- **MIME 検出** — 拡張子のないファイルでは、先頭約512バイト（`content_type_sniffer.go`）からコンテンツタイプを判定し、その結果を`mimecache.go` / `ringqueue.go`にキャッシュします。
- **セキュリティヘッダー** — `file://`ナビゲーションを禁止し、`nosniff`を設定します。

すべてが埋め込まれるため、配布されるバイナリには<strong>外部依存関係がありません</strong>（Windows 上でも同様です）。

---

## 開発環境と本番環境の橋渡し

`pkg/application`から見ると、両方のモードが公開するのは<strong>同じ公開インターフェース</strong>です。つまり、`Handler http.Handler`を持つ`AssetOptions`構造体に加え、`internal/assetserver/`内のミドルウェアとライフサイクル配線です。開発モードと本番モードの切り替えは Go のビルドタグのみで行われるため、どちらのモードでもアプリケーションコードは同一です。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## フロントエンドフレームワークの統合方法

### テンプレート

同梱される各テンプレート（React、Vue、Svelte、Solid、Vanilla など）には、次のものが含まれています。

- `build/Taskfile.yml`
- `frontend/vite.config.ts`（または同等のもの）

各テンプレートの Vite または同等の設定は `WAILS_VITE_PORT` を読み取り、開発サーバーをそのポートにバインドします。その後、CLI はアプリ内プロキシが使用する実際の `FRONTEND_DEVSERVER_URL` を公開します。

フレームワークは Go から完全に分離されたままです。

- ビルド時に Wails JS SDK をインポートする必要はありません。`/wails/runtime.js` は実行時にアセットサーバーから配信されます。
- HTTP 開発サーバーを備えたフレームワークであれば、どれでも統合できます。

---

## 拡張／カスタマイズ

カスタムヘッダー、認証、gzip が必要ですか？

1. `middleware.Middleware`（`func(http.Handler) http.Handler` のエイリアスで、`internal/assetserver/middleware.go` で宣言）を定義します。
2. `internal/assetserver/options.go` が公開する設定を介して、これを `application.AssetOptions` に組み込みます。
3. 動作は開発環境と本番環境で同一です。モードごとに個別のミドルウェアリストはありません。

---

## 主要なソースファイル

| ファイル | 役割 |
| --- | --- |
| `build_dev.go` / `build_production.go` | 開発環境と本番環境を選択するビルドタグ用ラッパー |
| `assetserver.go` / `asset_fileserver.go` | 中核となる HTTP ハンドラー |
| `assetserver_dev.go` | `FRONTEND_DEVSERVER_URL` へのリバースプロキシ |
| `bundled_assetserver.go` | `bundledassets/` を包む手書きのラッパー |
| `options.go` | `application.AssetOptions` 向けの設定 |
| `mimecache.go` / `ringqueue.go` | MIME キャッシュと小規模な LRU |

---

## 注意点とデバッグ

- **本番環境で画面が真っ白になる** — 通常は SPA ルーティングが原因です。開発サーバーが不明なパスに対して `index.html` を配信し、埋め込みの本番用ハンドラーのフォールバックまで処理が到達することを確認してください。
- **404が開発環境で発生する場合** — Vite の設定が `WAILS_VITE_PORT` にバインドされていないか、CLI が開発サーバーに到達できず、`FRONTEND_DEVSERVER_URL` を設定できなかったことが原因です。
- **大容量のアセット** — 埋め込むとバイナリが肥大化します。大容量のメディアは別のオリジンから配信するか、カスタム `http.Handler` を介してストリーミングしてください。

---

これで、Wails の <strong>アセットサーバー</strong>が、<strong>開発環境</strong>と<strong>本番環境</strong>の両方で、Web コードをネイティブウィンドウに提供する仕組みを理解できました。このレイヤーを習得すれば、読み込みの問題をデバッグしたり、ミドルウェアを追加したり、さらにはまったく別のフロントエンドツールチェーンに自信を持って切り替えたりできます。
