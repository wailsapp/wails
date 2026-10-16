---
title: "カスタムトランスポート層の作成"
description: "独自の Wails v3 カスタム IPC トランスポート層を作成し、カスタマイズする方法を説明します"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 では、生成されたすべてのバインディングとイベント通信を維持したまま、カスタム IPC トランスポート層を使用できます。これにより、デフォルトの HTTP fetch ベースのトランスポートを、WebSocket、カスタムプロトコル、またはその他の任意のトランスポート方式に置き換えられます。

## 概要

デフォルトでは、Wails はフロントエンドからの HTTP fetch リクエストを使用し、`/wails/runtime` を介してバックエンドと通信します。カスタムトランスポート API では、次のことが可能です。

- HTTP トランスポートを WebSocket、gRPC、または任意のカスタムプロトコルに置き換える
- Wails のコード生成との完全な互換性を維持する
- 既存のすべてのバインディング、イベント、ダイアログ、およびその他の Wails 機能を維持する
- 独自の接続管理、認証、エラー処理を実装する

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

## 使用方法

### 1. Transport インターフェースを実装する

`Transport` インターフェースを実装して、カスタムトランスポートを作成します。

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. アプリケーションを設定する

カスタムトランスポートをアプリケーションのオプションに渡します。

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. フロントエンドランタイムを変更する

カスタムトランスポートを使用する場合は、HTTP fetch の代わりにそのトランスポートを使用するよう、フロントエンドランタイムを変更する必要があります。リクエストの処理に使用される `RuntimeTransport` インターフェースを実装します。

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## 注意事項

- カスタムトランスポートを指定しなければ、デフォルトの HTTP トランスポートが引き続き動作します
- 生成されたバインディングは変更されず、変更されるのはトランスポート層のみです
- イベント、ダイアログ、クリップボード、およびその他のすべての Wails 機能は透過的に動作します
- カスタムトランスポートにおけるエラー処理、再接続ロジック、セキュリティは、実装者が担う必要があります
- 提示されている WebSocket の例はデモ用であり、本番環境で使用するには堅牢化が必要になる場合があります

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

### AssetServerTransport インターフェース（任意）

ブラウザベースでデプロイする場合、またはカスタムトランスポートを介してアセットと IPC の両方を提供する場合は、`AssetServerTransport` インターフェースを実装します。

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**このインターフェースを実装するケース：**

- WebView ではなくブラウザでアプリを実行する場合
- カスタム IPC トランスポートと併用して、HTTP 経由でアセットを提供する場合
- ネットワーク経由でアクセス可能なアプリケーションを構築する場合

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
