---
title: "サーバービルド"
description: "ネイティブGUIウィンドウを使用せずに、WailsアプリケーションをHTTPサーバーとして実行する"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3はサーバーモードをサポートしており、ネイティブウィンドウを作成せず、GUI依存関係も必要とせずに、アプリケーションを純粋なHTTPサーバーとして実行できます。これにより、同じWailsアプリケーションをサーバー、コンテナ、Webブラウザー向けにデプロイできます。

サーバーモードは、次の用途に役立ちます。

- **Docker／コンテナへのデプロイ** — X11／Wayland依存関係なしで実行
- **サーバーサイドアプリケーション** — ブラウザーからアクセスできるWebサーバーとしてデプロイ
- **Webのみからのアクセス** — デスクトップとWebで同じコードベースを共有
- **CI/CDテスト** — ディスプレイサーバーなしで統合テストを実行
- **マイクロサービス** — ヘッドレスのバックエンドサービスでWailsバインディングを使用

## クイックスタート

サーバーモードは、`server`ビルドタグで有効にします。アプリケーションコードを変更する必要はなく、このタグを指定してビルドするだけです。

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

同じコードを、デスクトップ向けにはタグなしで、サーバーモード向けには`-tags server`を指定してビルドできます。

## 設定

### ServerOptions

`ServerOptions`を使用してHTTPサーバーを設定します。

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

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
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

ヘルスチェックエンドポイントは、`/health`で自動的に利用できます。

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

これは、次の用途に役立ちます。

- Kubernetesのliveness／readinessプローブ
- ロードバランサーのヘルスチェック
- 監視システム

### サービスバインディング

すべてのサービスバインディングは、デスクトップモードと同じように動作します。

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

フロントエンドからは、標準のWailsランタイムを使用してこれらのバインディングを呼び出せます。

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### イベント

サーバーモードでも、イベントは双方向に動作します。

- **フロントエンドからバックエンド**：ブラウザーから送出されたイベントはHTTP経由で送信され、Goのイベントハンドラーによって受信されます
- **バックエンドからフロントエンド**：Goから送出されたイベントは、WebSocket経由で接続中のすべてのブラウザーにブロードキャストされます

各ブラウザータブは、一意の名前（`browser-1`、`browser-2`など）を持つ「ウィンドウ」として表され、`event.Sender`からアクセスできます。

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

フロントエンドからは、次のようにします。

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### グレースフルシャットダウン

サーバーは、`SIGINT`および`SIGTERM`シグナルを適切に処理します。

1. 新しい接続の受け付けを停止します
2. 処理中のリクエストが完了するまで待機します（最大`ShutdownTimeout`）
3. `OnShutdown`フックを実行します
4. サービスを逆順でシャットダウンします

## デスクトップモードとの違い

| 機能 | デスクトップモード | サーバーモード |
| --- | --- | --- |
| ネイティブウィンドウ | 作成される | ブラウザーウィンドウ（`browser-N`） |
| システムトレイ | 利用可能 | 利用不可 |
| ネイティブダイアログ | 利用可能 | 利用不可 |
| アプリケーションメニュー | 利用可能 | 利用不可 |
| 画面情報 | 利用可能 | エラーを返す |
| サービスバインディング | 動作する | 動作する |
| イベント | 動作する | 動作する（WebSocket経由） |
| アセット | WebView経由 | HTTP経由 |
| CGOが必要 | はい | いいえ |

### ウィンドウAPIの動作

サーバーモードでは、ウィンドウ関連のAPIは安全に処理されます。

- `app.Window.NewWithOptions()` — 警告をログに記録し、nilを返します
- `app.Hide()`／`app.Show()` — 何も実行しません
- `app.Screen.GetPrimary()` — エラーを返します

これにより、ウィンドウを参照するコードもクラッシュせずに実行できますが、ウィンドウ操作は何の効果も持ちません。

## 本番環境向けのビルド

### Task を使用する（推奨）

`wails3 init` で作成したプロジェクトには、`build:server` タスクが含まれています：

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

Wails プロジェクトには、すぐに使用できる Docker 構成が含まれています。コンテナ内でアプリケーションをビルドして実行するには：

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

以上です。アプリケーションは `http://localhost:8080` で利用できます。

いくつかのオプションを使用してビルドをカスタマイズできます：

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

生成された `Dockerfile.server` は、distroless をベースとする最小限のイメージを作成します。ネットワークへのバインドは自動的に処理されるため、コンテナの外部からアプリケーションにアクセスできます。

### Docker Compose

より複雑なデプロイには、ヘルスチェックを含む次の Docker Compose 構成を使用できます：

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
このヘルスチェックの例では `wget` を使用しています。distroless ベースイメージを使用する場合は、ヘルスチェック用バイナリをイメージに含めるか、外部のヘルスチェック機構（Docker の `curl` オプションやサイドカーコンテナなど）を使用する必要があります。

@end

### カスタム Dockerfile

より細かく制御する必要がある場合は、独自の Dockerfile を作成できます。重要なのは、サーバーがコンテナ外部からの接続を受け入れるように `WAILS_SERVER_HOST=0.0.0.0` を設定することです：

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

サーバーモードのアプリケーションをデプロイする場合：

1. **デフォルトでは localhost にバインドする** - `0.0.0.0` は必要な場合にのみ使用してください
2. **本番環境では TLS を使用する** - `ServerOptions.TLS` を設定してください
3. **リバースプロキシの背後に配置する** - セキュリティを強化するために nginx/traefik を使用してください
4. **WebSocket を同一オリジンに保つ** - `WebSocketOriginPatterns` には信頼できるオリジンのみを追加し、`WebSocketAllowAllOrigins` は避けてください
5. **すべての入力を検証する** - あらゆる Web アプリケーションと同じセキュリティ対策を適用してください

## 例

完全な例は `v3/examples/server/` にあります：

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## 環境変数

コードを変更せずにサーバー構成を上書きする必要があるデプロイ環境向けに、Wails は次の環境変数を認識します：

| 変数 | 説明 | デフォルト |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | バインド先のネットワークインターフェース | `localhost` |
| `WAILS_SERVER_PORT` | 待ち受けるポート | `8080` |

これらはコード内の `ServerOptions` より優先されます。そのため、Docker の例では `WAILS_SERVER_HOST=0.0.0.0` を設定しています。これにより、アプリケーションを変更することなく、コンテナが外部からの接続を受け入れられます。

## 関連項目

- [カスタムトランスポート](/guides/custom-transport/) - 高度な IPC カスタマイズについて
- [サービス](/features/bindings/services/) - サービスバインディングのドキュメント
- [イベント](/guides/events-reference/) - イベントシステムのドキュメント

### ランタイムリクエストのサイズ

`/wails/runtime` へのリクエストは、JSON 処理前の段階で 64 MiB に制限されます。この上限を超える通常のリクエストには、`Content-Length` がないリクエストも含め、HTTP 413 が返されます。チャンク化されたランタイムアップロードには、チャンクあたり 1 MiB、結合後のペイロードに 64 MiB という個別の上限が引き続き適用されます。必要に応じて、アプリケーションミドルウェアまたはリバースプロキシを使用して、より低い上限を設定してください。
