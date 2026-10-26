---
title: "サーバービルド"
description: "ネイティブ GUI ウィンドウを使用せずに Wails アプリケーションを HTTP サーバーとして実行する"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 はサーバーモードをサポートしており、ネイティブウィンドウを作成せず、GUI の依存関係も必要としない純粋な HTTP サーバーとしてアプリケーションを実行できます。これにより、同じ Wails アプリケーションをサーバー、コンテナ、Web ブラウザーにデプロイできます。

サーバーモードは、次の用途に役立ちます。

- **Docker／コンテナへのデプロイ** - X11／Wayland の依存関係なしで実行
- **サーバーサイドアプリケーション** - ブラウザーからアクセスできる Web サーバーとしてデプロイ
- **Web のみのアクセス** - デスクトップと Web で同じコードベースを共有
- **CI/CD テスト** - ディスプレイサーバーなしで統合テストを実行
- **マイクロサービス** - ヘッドレスのバックエンドサービスで Wails バインディングを使用

## クイックスタート

サーバーモードは、`server` ビルドタグで有効にします。アプリケーションコードはそのままで、タグを指定してビルドするだけです。

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

最小構成の例を次に示します。

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

同じコードを、デスクトップ向けにはタグなしで、サーバーモード向けには `-tags server` 付きでビルドできます。

## 設定

### ServerOptions

`ServerOptions` を使用して HTTP サーバーを設定します。

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## 機能

### ヘルスチェックエンドポイント

ヘルスチェックエンドポイントは `/health` で自動的に利用できます。

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

これは次の用途に役立ちます。

- Kubernetes の liveness／readiness プローブ
- ロードバランサーのヘルスチェック
- 監視システム

### サービスバインディング

すべてのサービスバインディングは、デスクトップモードと同様に動作します。

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

フロントエンドから、標準の Wails ランタイムを使用してこれらのバインディングを呼び出せます。

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### イベント

サーバーモードでは、イベントが双方向に動作します。

- **フロントエンドからバックエンドへ**：ブラウザーから発行されたイベントは HTTP 経由で送信され、Go のイベントハンドラーで受信されます
- **バックエンドからフロントエンドへ**：Go から発行されたイベントは、WebSocket 経由で接続中のすべてのブラウザーにブロードキャストされます

各ブラウザータブは、一意の名前（`browser-1`、`browser-2` など）を持つ「ウィンドウ」として表され、`event.Sender` からアクセスできます。

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

フロントエンドからは、次のように操作します。

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### グレースフルシャットダウン

サーバーは `SIGINT` および `SIGTERM` シグナルを適切に処理します。

1. 新しい接続の受け付けを停止
2. 処理中のリクエストが完了するまで待機（最大 `ShutdownTimeout`）
3. `OnShutdown` フックを実行
4. サービスを逆順でシャットダウン

## デスクトップモードとの違い

| 機能 | デスクトップモード | サーバーモード |
| --- | --- | --- |
| ネイティブウィンドウ | 作成される | ブラウザーウィンドウ（`browser-N`） |
| システムトレイ | 利用可能 | 利用不可 |
| ネイティブダイアログ | 利用可能 | 利用不可 |
| アプリケーションメニュー | 利用可能 | 利用不可 |
| 画面情報 | 利用可能 | ネイティブ画面なし。`GetPrimary()` は `nil` を返す |
| サービスバインディング | 動作する | 動作する |
| イベント | 動作する | 動作する（WebSocket 経由） |
| アセット | WebView 経由 | HTTP 経由 |
| CGO が必要 | はい | いいえ |

### ウィンドウ API の動作

サーバーモードでは、ウィンドウ関連の API が安全に処理されます。

- `app.Window.NewWithOptions(options)` - nil ではない Go のウィンドウオブジェクトを返します。ただし、ネイティブのデスクトップウィンドウは作成されません
- `app.Hide()`／`app.Show()` - 何も実行しません
- `app.Screen.GetPrimary()` - `*application.Screen` を返します。プライマリ画面を利用できない場合、これは `nil` になります。参照を解除する前に確認してください

Go のウィンドウオブジェクトと、ネイティブウィンドウを利用できるかどうかは区別してください。サーバー向けにビルドする場合は、デスクトップ専用の動作をガードしてください。Go のウィンドウオブジェクトが nil でなくても、ネイティブウィンドウや画面が存在するとは限りません。

### 対応プラットフォーム

サーバーモードは純粋な Go で実装されており、CGO を使用せずに、`android` や `ios` を含む、Wails がサポートするすべての `GOOS` 向けにコンパイルできます。モバイルターゲットでは、ヘッドレスサーバーのみがビルドされます。

- ネイティブのモバイル API である `application.Android` と `application.IOS` は利用できません。これらを使用するコードは、`!server` ビルド制約を使用してサーバービルドから除外する必要があります。
- デフォルトのロガーはデスクトップ用です。開発ビルドでは、logcat や WebView コンソールではなく stderr に書き込みます。

## 本番環境向けのビルド

### Task の使用（推奨）

`wails3 init` で作成したプロジェクトには、`build:server` タスクが含まれています。

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### 手動ビルド

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails プロジェクトには、すぐに使用できる Docker 構成が含まれています。アプリケーションをコンテナ内でビルドして実行するには、次のようにします。

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

これで完了です。アプリケーションは `http://localhost:8080` で利用できます。

いくつかのオプションを使用してビルドをカスタマイズできます。

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

生成される `Dockerfile.server` は、distroless をベースとする最小限のイメージを作成します。ネットワークのバインドは自動的に処理されるため、コンテナの外部からアプリケーションにアクセスできます。

### Docker Compose

より複雑なデプロイには、次のヘルスチェック付き Docker Compose 構成を使用できます。

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
ヘルスチェックの例では `wget` を使用しています。distroless ベースイメージを使用する場合は、ヘルスチェック用バイナリをイメージに含めるか、外部のヘルスチェック機構（Docker の `curl` オプションやサイドカーコンテナなど）を使用する必要があります。

@end

### カスタム Dockerfile

より細かく制御する必要がある場合は、独自の Dockerfile を作成できます。重要なのは、サーバーがコンテナ外部からの接続を受け入れるように `WAILS_SERVER_HOST=0.0.0.0` を設定することです。

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## セキュリティ上の考慮事項

サーバーモードのアプリケーションをデプロイする場合は、次の点に注意してください。

1. **デフォルトでは localhost にバインドする** - 必要な場合にのみ `0.0.0.0` を使用してください
2. **本番環境では TLS を使用する** - `ServerOptions.TLS` を設定してください
3. **リバースプロキシの背後に配置する** - セキュリティをさらに強化するために nginx/traefik を使用してください
4. **WebSocket とランタイム呼び出しを同一オリジンに保つ** - `WebSocketOriginPatterns` には信頼できるオリジンのみを追加し、`WebSocketAllowAllOrigins` は避けてください
5. **すべての入力を検証する** - 一般的な Web アプリケーションと同じセキュリティ対策を適用してください

## 例

完全な例は `v3/examples/server/` にあります。

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## 環境変数

コードを変更せずにサーバー設定を上書きする必要があるデプロイ環境向けに、Wails は次の環境変数を認識します。

| 変数 | 説明 | デフォルト |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | バインド先のネットワークインターフェイス | `localhost` |
| `WAILS_SERVER_PORT` | 待ち受けるポート | `8080` |

これらはコード内の `ServerOptions` より優先されます。そのため、Docker の例では `WAILS_SERVER_HOST=0.0.0.0` を設定しています。これにより、アプリケーションを変更せずに、コンテナが外部接続を受け入れられます。

## 関連項目

- [カスタムトランスポート](/guides/custom-transport/) - 高度な IPC のカスタマイズについて
- [サービス](/features/bindings/services/) - サービスバインディングのドキュメント
- [イベント](/guides/events-reference/) - イベントシステムのドキュメント

### ランタイムリクエストのサイズ

`/wails/runtime` へのリクエストは、JSON 処理前のサイズが 64 MiB に制限されます。`Content-Length` がないリクエストを含め、この上限を超える通常のリクエストには HTTP 413 が返されます。チャンク化されたランタイムアップロードには、チャンクごとに 1 MiB、結合後のペイロードに 64 MiB という個別の上限が引き続き適用されます。必要に応じて、アプリケーションのミドルウェアまたはリバースプロキシを使用して、より低い上限を設定してください。

### ランタイムリクエストのオリジン

`/wails/runtime` が受け付けるのは `POST` リクエストのみで、その他のメソッドには HTTP 405 が返されます。ブラウザーからのリクエストでは、`Origin` がリクエスト先ホストと一致せず、かつ `WebSocketOriginPatterns` にも一致しない場合、HTTP 403 が返されます。また、`Origin` がなく、`Sec-Fetch-Site` が `cross-site` または `same-site` である場合も同様です。`curl` などのブラウザー以外のクライアントからのリクエストのように、どちらのヘッダーも含まないリクエストは受け入れられます。リバースプロキシが `Host` ヘッダーを書き換える場合は、公開オリジンを `WebSocketOriginPatterns` に追加してください。

`WebSocketOriginPatterns` は、公開オリジンがアプリケーションから見える `Host` と異なるリバースプロキシ構成のために用意されています。別のオリジンから配信され、`/wails/runtime` を直接呼び出すフロントエンドはサポートされません。このエンドポイントは CORS ヘッダーを出力せず、プリフライトの `OPTIONS` リクエストには HTTP 405 を返します。フロントエンドを Wails サーバーから配信するか、両方を同じプロキシオリジンの背後に配置してください。

### ループバックの Host チェック

サーバーがループバックアドレス（デフォルトの `localhost`、任意の `127.0.0.0/8` アドレス、または `::1`。`ServerOptions.Host` と `WAILS_SERVER_HOST` のどちらで設定した場合も同様）にバインドされている場合、すべてのリクエストには、`localhost`、IP リテラル、バインド先ホスト、または `WebSocketOriginPatterns` エントリのホスト部分に一致するホストを指定した `Host` ヘッダーが必要です。それ以外の `Host` には HTTP 403 が返されます。これにより、悪意のある Web ページが自身のホスト名を `127.0.0.1` に向け、リクエストを前述の同一オリジンチェックに通過させる DNS リバインディングを防止します。ポートは比較されず、`localhost.`（末尾にドットがある形式）や `127.1` などの表記は拒否されます。`WebSocketAllowAllOrigins` が設定されている場合も、このチェックは適用されます。ローカルのリバースプロキシが公開 `Host` を転送する場合は、そのホストを `WebSocketOriginPatterns` に追加してください。コンテナ内の `0.0.0.0` など、その他のアドレスへのバインドには影響しません。
