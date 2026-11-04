---
title: "カスタムトランスポート層の作成"
description: "独自の Wails v3 カスタム IPC トランスポート層を作成し、カスタマイズする方法を説明します"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 では、生成されたすべてのバインディングとイベント通信を維持したまま、カスタム IPC トランスポート層を使用できます。これにより、デフォルトの HTTP fetch ベースのトランスポートを、WebSocket、カスタムプロトコル、またはその他の任意のトランスポート機構に置き換えられます。

## 概要

デフォルトでは、Wails はフロントエンドからの HTTP fetch リクエストを使用し、`/wails/runtime` を介してバックエンドと通信します。カスタムトランスポート API を使用すると、次のことができます。

- HTTP トランスポートを WebSocket、gRPC、または任意のカスタムプロトコルに置き換える
- Wails のコード生成との完全な互換性を維持する
- 既存のすべてのバインディング、イベント、ダイアログ、およびその他の Wails 機能を維持する
- 独自の接続管理、認証、およびエラー処理を実装する

## アーキテクチャ

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## 完全な WebSocket サンプルの実行

[WebSocket トランスポートのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport)には、リクエスト ID とレスポンス、ライフサイクル管理、JavaScript クライアント、イベント配信という、プロトコルの両側の実装が含まれています。`Start` と `Stop` だけを実装したトランスポートではなく、これらのファイルを出発点にしてください。

ネイティブプラットフォームの依存関係がインストールされている Wails リポジトリの作業コピーで、次を実行します。

```bash
cd v3/examples/websocket-transport
go run .
```

このサンプルはネイティブウィンドウを開き、ポート 9099 で WebSocket エンドポイントを起動します。信頼できるローカル開発環境で使用してください。ネットワークからアクセス可能な派生版をデプロイする前に、リスナーを制限し、オリジンを検証して、クライアントを認証してください。このサンプルのオリジンポリシーは制限が緩く、本番環境向けのセキュリティポリシーではありません。

### バックエンドの規約

トランスポートは、必須の 3 つのメソッドをすべて実装する必要があります。サンプルの `WebSocketTransport` は `Start`、`JSClient`、`Stop` を実装し、`JSClient` は Wails がページに挿入する JavaScript シムを提供します。実装には、コンパイル時のインターフェースチェックを含めてください。

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

サンプルのコンストラクターを使用してアプリケーションを構成します。

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

このフィールドは、通常のアセットおよびサービスと並べて `application.Options` に配置します。コンストラクター呼び出しだけでなく、実装ファイルもコピーしてください。`Start` は `MessageProcessor` を受け取ります。`Stop` は、リスナー、アクティブな接続、および保留中の処理を終了しなければなりません。

### フロントエンドの規約

サンプルの `websocket-transport.js` は、リクエストとレスポンスの対応付けや、失敗した呼び出しの reject 処理を含むランタイムトランスポートを実装しています。そのファイルをアセットから配信した後、生成されたサービスバインディングを呼び出す前に初期化してください。

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

生成されたサービスバインディングは変更されません。トランスポートがクライアントを置き換える場合は、アプリケーションが使用するメソッドを提供し、そのエラーセマンティクスを維持する必要があります。イベント配信には、オプションのイベントトランスポート規約が必要です。任意の `Transport` 実装によって自動的に提供されるものではありません。完全な WebSocket サンプルでその実装を確認できます。

リモート環境へデプロイする場合は、TLS と、明示的な認証ポリシーおよびオリジンポリシーを使用してください。ブラウザのみでホストする場合は、以下で説明するように、アプリケーションアセットの配信も必要です。

## API リファレンス

### Transport インターフェース

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport インターフェース（オプション）

ブラウザベースのデプロイメントで使用する場合、またはカスタムトランスポートを介してアセットと IPC の両方を配信する場合は、`AssetServerTransport` インターフェースを実装してください。

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**このインターフェースを実装する場合：**

- WebView ではなくブラウザでアプリを実行する場合
- カスタム IPC トランスポートと併せて HTTP 経由でアセットを配信する場合
- ネットワークからアクセス可能なアプリケーションを構築する場合

**実装例：**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

`ServeAssets()` が呼び出されると、assetHandler は次のものを提供します。

- すべての静的アセット（HTML、CSS、JS、画像など）
- `/wails/runtime.js` - Wails ランタイムライブラリ

## 関連項目

- `transport.go` - コアトランスポートのインターフェースと型
- `messageprocessor.go` - すべての Wails IPC を処理する基盤のメッセージプロセッサ
