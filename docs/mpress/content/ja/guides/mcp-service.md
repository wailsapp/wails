---
title: "LLM による制御（MCP）"
description: "Model Context Protocol を介して、LLM エージェントによるアプリのテストと制御を可能にします"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="実験的機能"}
組み込み MCP サーバーは実験的な機能であり、今後のリリースで API が変更される可能性があります。

@end

Wails v3 には、LLM エージェント（Claude Code、IDE アシスタント、任意の MCP クライアント）が<strong>実行中</strong>の Wails アプリケーションを検査、テスト、操作できる、組み込みの[Model Context Protocol](https://modelcontextprotocol.io)（MCP）サーバーがあります。

プロジェクトのライフサイクルを自動化するには、別の`wails3 mcp` CLI サーバーを使用してください。このサーバーにより、エージェントはプロジェクトの検査と初期化、診断の実行、ビルドおよび開発ジョブの開始、バインディングの生成、名前付き Taskfile タスクの実行、上限付きジョブ出力の取得を行えます。CLI サーバーはデフォルトで現在のディレクトリ内に制限され、任意のシェルコマンドを実行する機能は公開しません。トランスポート、認証、ツールの詳細については、[CLI MCP ドキュメント](/guides/cli/#mcp)を参照してください。

有効にすると、アプリに接続したエージェントは次の操作を実行できます。

- **ウィンドウの一覧表示と制御** — サイズ、位置、フォーカス、フルスクリーン、開発者ツール、再読み込みなど
- **DOM の検査** — 要素の検索、HTML の取得、構造スナップショットの取得
- **JavaScript の評価** — 任意のウィンドウ内で任意のコードを実行し、結果を取得
- **ユーザー入力のシミュレーション** — マウスの移動、クリック、ドラッグ、スクロールを<strong>アニメーション付きの画面上のカーソル</strong>で表示し、エージェントの操作を確認可能
- **文字入力とキー操作** — React の制御コンポーネントによる入力でも動作する、実際に近い文字単位のイベント
- <strong>バインドされた Go メソッドの呼び出し</strong>と、アプリケーションイベントの発行および待機

## 仕組み

MCP サーバーがアプリケーションにコンパイルされるのは、<strong>`mcp` ビルドタグ</strong>が指定されている場合だけです。このタグがなければ、サーバーコードはバイナリに一切含まれません。実行時のオーバーヘッドも、開いているポートも、攻撃対象領域もありません。

タグが指定されている場合、サーバーは`App.Run()`内で自動的に起動し、デフォルトで`127.0.0.1:9099`にバインドして、そのエンドポイントをログに記録します。ユーザーコードは必要ありません。

## チュートリアル

### ステップ 1 — 通常の Wails アプリケーションを作成する

MCP では、インポートも登録も必要ありません。通常どおりアプリを作成してください。

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### ステップ 2 — `mcp` タグを指定してビルドまたは実行する

@tabs
[Wails CLI（推奨）]
`WAILS_MCP=1`を設定すると、Wails CLI が`mcp`タグを自動的に追加します。

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go を直接使用]
`go run`または`go build`にタグを直接渡します。

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows（PowerShell）]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

起動時に、アプリケーションは MCP エンドポイントをログに記録します。

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### ステップ 3 — クライアントを接続する

サーバーは<strong>MCP Streamable HTTP トランスポート</strong>を使用します。任意の MCP 互換クライアントで接続してください。

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

次に、Claude にアプリの操作を依頼します。

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code（GitHub Copilot）]
`.vscode/settings.json`に追加します。

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[その他のクライアント]
Streamable HTTP トランスポートをサポートする任意の MCP クライアントで、次のアドレスを指定します。

```
http://127.0.0.1:9099/mcp
```

@end

### ステップ 4 — テストセッションを実行する

エージェントにアプリケーションの動作確認を依頼します。プロンプトの例を次に示します。

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## 設定

すべて環境変数で設定でき、コードを変更する必要はありません。

| 環境変数 | デフォルト | 説明 |
| --- | --- | --- |
| `WAILS_MCP` | （未設定） | Wails CLI の使用時に`mcp`ビルドタグを自動的に追加するには、`1`、`true`、`on`、または`yes`に設定します。 |
| `WAILS_MCP_HOST` | `127.0.0.1` | バインド先のインターフェース。ループバック以外へのバインドには`WAILS_MCP_TOKEN`が必要です。 |
| `WAILS_MCP_TOKEN` | 未設定 | ループバックでは任意の Bearer トークン。他のバインドアドレスでは必須です。クライアントは`Authorization: Bearer <token>`を送信します。 |
| `WAILS_MCP_PORT` | `9099` | 待ち受けポート。ランダムに割り当てられる空きポートを使用するには`0`に設定します（ポートはログに出力されます）。 |
| `WAILS_MCP_TIMEOUT` | `30000` | JavaScript 評価のデフォルトのタイムアウト（**ミリ秒**）。 |
| `WAILS_MCP_HIDE_CURSOR` | （未設定） | アニメーション付きカーソルのオーバーレイを無効にするには、`1`または`true`に設定します。 |

例 — カスタムポートと60秒のタイムアウト：

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## 利用可能なツール

| ツール | 用途 |
| --- | --- |
| `app_info` | アプリケーション情報：プラットフォーム、アーキテクチャ、すべてのウィンドウ、MCP エンドポイント |
| `windows_list` | すべてのウィンドウを位置、サイズ、状態とともに一覧表示 |
| `window_control` | フォーカス、サイズ変更、移動、フルスクリーン、開発者ツール、再読み込み、URL の設定など（22個のアクション） |
| `js_eval` | ウィンドウ内で JavaScript を評価（非同期の本体、値には`return`） |
| `dom_html` | ページまたは特定の要素の HTML を取得 |
| `dom_query` | CSS セレクターで要素を検索 — タグ、テキスト、境界、可視性 |
| `screenshot_dom` | 表示中のページの構造スナップショット（DOM ベース、ピクセル情報なし） |
| `mouse_move` | 指定した座標または CSS セレクターまでカーソルをアニメーション移動 |
| `mouse_click` | アニメーションカーソルでクリック（左／右／中央、ダブルクリック、修飾キー） |
| `mouse_drag` | アニメーションカーソルでドラッグ（HTML5 のドラッグ＆ドロップ要素に対応） |
| `mouse_scroll` | 指定した座標または要素でスクロール |
| `keyboard_type` | 実際の操作に近いイベントを発生させながら、テキストを1文字ずつ入力 |
| `keyboard_press` | 単一のキー（Enter、Tab、Escape、ArrowDown など）を、必要に応じて修飾キーとともに押下 |
| `call_bound_method` | バインドされた Go サービスメソッドを呼び出し（例：`main.GreetService.Greet`） |
| `emit_event` | Wails アプリケーションイベントを発行 |
| `wait_for_event` | Wails アプリケーションイベントを待機し、そのデータを返す |

### マルチウィンドウ対応

ウィンドウを操作するすべてのツールでは、ウィンドウの<strong>名前</strong>（`WebviewWindowOptions.Name`で設定）を含む省略可能な`window`引数を指定できます。省略すると、現在フォーカスされているウィンドウが対象になります。フォーカスされているウィンドウがない場合は、最初のウィンドウが対象になります。

```
List all windows, then click the "New" button in the window named "editor".
```

### 要素の選択

マウスおよびキーボードツールでは、次のいずれかを指定できます。

- **CSS セレクター** — `selector: "#submit-btn"`（要素が自動的に表示領域内へスクロールされます）
- **座標** — `x: 400, y: 300`（ビューポートを基準とする CSS ピクセル）

ドラッグ操作では、`from_`および`to_`を先頭に付けます。

```
Drag from selector: ".card" to selector: ".dropzone"
```

## セキュリティ

@note{type="caution"}
MCP サーバーを使用すると、アプリケーションをプログラムから完全に制御できます。そのツールにアクセスできる人は誰でも、DOM の読み取り、JavaScript の評価、ボタンのクリック、Go メソッドの呼び出しが可能です。

@end

- サーバーはデフォルトで`127.0.0.1`にバインドされます。ブラウザーのオリジンは HTTP(S) のループバックオリジンでなければなりません。不透明なオリジン（`null`）、不正な形式のオリジン、外部オリジンは拒否されます。
- ヘッダーを持たないネイティブ MCP クライアントも引き続きサポートされます。`WAILS_MCP_TOKEN`を設定しない場合、ローカルプロセスおよび許可されたローカルブラウザーオリジンは信頼されます。オリジンチェックは認証ではありません。すべての`/mcp`呼び出しで Bearer 認証を必須にするには、十分なエントロピーを持つトークンを設定します。クライアントの`Authorization: Bearer <token>`ヘッダーにも同じトークンを設定してください。プリフライトリクエストにはトークンは不要です。
- `/eval-result`コールバックでは、クライアントの Bearer トークンではなく、評価ごとに予測不能な ID を使用するため、WebView の結果配信との互換性が維持されます。
- 本番ビルドには`mcp`タグを含めるべきでは<strong>ありません</strong>。Wails CLI がこのタグを追加するのは、`WAILS_MCP=1`が明示的に設定されている場合のみです。また、デフォルトの`wails3 build`にはサーバーコードが一切含まれていません。
- ループバック以外のインターフェイスでサーバーを公開する必要がある場合（LAN テストなど）は、`WAILS_MCP_HOST=0.0.0.0`と、十分なエントロピーを持つ`WAILS_MCP_TOKEN`を設定してください。トークンがない場合、起動に失敗します。組み込みリスナーは HTTP を使用するため、信頼できないネットワークでは暗号化トンネルまたは TLS 終端プロキシを使用してください。

## サンプルアプリケーション

すべてのツールのデモを行う完全なプレイグラウンドアプリケーションを[`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp)で利用できます。次の機能が含まれています。

- インクリメント／リセットボタン付きカウンター
- Greet、Add、Shout の各バインドメソッドを備えた名前入力欄
- HTML5 ドラッグ＆ドロップのドラッグ元とドロップ先
- スクロール可能なリスト（50項目）
- イベントログ

次のコマンドで実行します。

```shell
cd v3/examples/mcp
go run -tags mcp .
```

次に、Claude Code または任意の MCP クライアントを`http://127.0.0.1:9099/mcp`に接続し、UI を操作するよう指示します。

## フィードバック

組み込み MCP サーバーは実験的な機能であり、今後の方向性は皆さまからのフィードバックによって決まります。お試しになった場合は、使用したクライアントとツール、期待していた動作と実際に起きたこと、エージェントにアプリを操作させることが有用だったかどうかをぜひお聞かせください。実行した内容を正確に記載した報告が最も役立ちます。[MCP サーバーのフィードバックディスカッション](https://github.com/wailsapp/wails/discussions/5692)でお知らせください。
