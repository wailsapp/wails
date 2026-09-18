---
title: "Raw メッセージ"
description: "パフォーマンスが最重要となるアプリケーション向けに、フロントエンドからバックエンドへのカスタム通信を実装する"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

Raw メッセージは、標準のバインディングシステムを介さずに、フロントエンドとバックエンド間の低レベル通信チャネルを提供します。利便性と引き換えに速度を向上させます。

## Raw メッセージを使用する場合

Raw メッセージは、次のような極端なエッジケースに最適です。

- **超高頻度の更新** - 1マイクロ秒単位の差が重要となる、毎秒数千件のメッセージ処理
- **カスタムメッセージプロトコル** - ワイヤ形式を完全に制御する必要がある場合

@note{type="tip"}
ほぼすべてのユースケースでは、標準の[サービスバインディング](/features/bindings/services/)を推奨します。オーバーヘッドを無視できる程度に抑えながら、型安全性、自動シリアル化、より優れた開発者体験を提供するためです。

@end

## バックエンドのセットアップ

アプリケーションオプションで`RawMessageHandler`を設定します。

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### ハンドラーのシグネチャ

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| パラメーター | 型 | 説明 |
| --- | --- | --- |
| `window` | `Window` | メッセージを送信したウィンドウ |
| `message` | `string` | Raw メッセージの内容 |
| `originInfo` | `*application.OriginInfo` | メッセージ送信元のオリジン情報 |

#### OriginInfo 構造体

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `Origin` | `string` | メッセージを送信したドキュメントのオリジン URL |
| `TopOrigin` | `string` | トップレベルのオリジン URL（iframe 内では Origin と異なる場合があります） |
| `IsMainFrame` | `bool` | メッセージがメインフレームから送信されたかどうか |

#### プラットフォーム別の利用可否

- **macOS**：`Origin`と`IsMainFrame`が提供されます
- **Windows**：`Origin`と`TopOrigin`が提供されます
- **Linux**：`Origin`のみが提供されます

### オリジンの検証

@note{type="caution"}
ハンドラーに届いたという理由だけで、メッセージが安全だと判断してはいけません。機密性の高い処理や状態を変更する処理を実行する前に、オリジン情報を検証する必要があります。

@end

<strong>受信したメッセージを処理する前に、必ずそのオリジンを確認してください。</strong>不正アクセスを防ぐには、`originInfo`パラメーターが提供する重要なセキュリティ情報を検証する必要があります。 悪意のあるコンテンツ、侵害されたコンテンツ、または意図しないスクリプトから Raw メッセージが送信される可能性があります。オリジンを検証しなければ、信頼できない送信元からのコマンドを処理してしまう可能性があります。`originInfo`を使用して、メッセージが想定した送信元から届いていることを確認してください。

### 検証の要点

- <strong>必ず`Origin`</strong>を確認する - オリジンが想定する信頼済みの送信元と一致することを確認します（通常、ローカルアセットの場合は`wails://wails`または`http://wails.localhost`、それ以外の場合はアプリ固有のオリジン）
- <strong>`IsMainFrame`</strong>を検証する（macOS） - メッセージが iframe から届いたかどうかを確認します。iframe からのメッセージは、異なるセキュリティコンテキストを持つ埋め込みコンテンツから送信された可能性があります
- <strong>`TopOrigin`</strong>を使用する（Windows） - フレーム化されたコンテンツを扱う場合は、トップレベルのオリジンを確認します
- **想定外のオリジンを拒否する** - 明示的に許可していないオリジンからのメッセージを拒否し、安全側で失敗させます

@note{type="info"}
`wails:`で始まるメッセージは Wails の内部通信用に予約されているため、ハンドラーには渡されません。

@end

## フロントエンドのセットアップ

`System.invoke()`を使用して Raw メッセージを送信します。

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### ビルド済みバンドルの使用

npm を使用していない場合は、グローバルな`wails`オブジェクトを介して`invoke`にアクセスします。

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## 構造化メッセージ

複雑なデータは JSON にシリアル化します。

### フロントエンド

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### バックエンド

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## パフォーマンス比較

| 方式 | オーバーヘッド | 型安全性 | ユースケース |
| --- | --- | --- | --- |
| サービスバインディング | 高い | 完全 | 汎用 |
| 生メッセージ | 最小限 | 手動 | 高頻度、パフォーマンス重視 |

### ベンチマーク例

単純なペイロードの場合、生メッセージ はサービスバインディングと比べて、1 秒あたりに処理できるメッセージ数が大幅に多くなります。

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## 完全な例

以下は、単純なコマンドプロトコルを実装する完全な例です。

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## ベストプラクティス

### 推奨事項

- 本当にパフォーマンスが重要な処理経路には 生メッセージ を使用する
- ハンドラーに適切なエラー処理を実装する
- フロントエンドにレスポンスを返すにはイベントを使用する
- 構造化データには JSON の使用を検討する
- ブロッキングを避けるため、メッセージ処理を高速に保つ

### 禁止事項

- サービスバインディングで十分な場合に 生メッセージ を使用する
- 受信メッセージの検証を忘れる
- 時間のかかる処理でハンドラーをブロックする（goroutine を使用する）
- レスポンスを特定のウィンドウに送信する必要がある場合に、window パラメーターを無視する

## マルチウィンドウに関する考慮事項

`window` パラメーターは、メッセージを送信したウィンドウを識別します。これにより、次のことが可能になります。

- 正しいウィンドウにレスポンスを送信する
- ウィンドウ固有の動作を実装する
- デバッグのためにメッセージの送信元を追跡する

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## 次のステップ

- [サービスバインディング](/features/bindings/services/) - ほとんどのアプリケーションで使用する標準的な方法
- [イベント](/guides/events-reference/) - バックエンドからフロントエンドへの通信に使用するイベントシステム
- [パフォーマンス](/guides/performance/) - 全般的なパフォーマンス最適化
