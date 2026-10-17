---
title: "ストリーム"
description: "WebSocket のプログラミングモデルを使用し、リッスンソケットを必要としない、Go と JavaScript 間の双方向バイトストリーム"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

ストリームは、Go とフロントエンドの間に、名前付きで順序が保証された双方向のバイトチャネルを提供します。プログラミングモデルは WebSocket と同じですが、**TCP ポートをバインドする必要はありません**。

WebSocket はカスタム URL スキーム上では通信できないため、WebView 内で使用するには、実際の HTTP サーバーを起動してポートをリッスンするしかありません。デスクトップアプリでは、マシン上のほかのどのプロセスからも到達できるローカルポートが開くことになります。安全性を確保するにはオリジンチェックとトークンが必要であり、ユーザーが使用しているすべてのファイアウォール製品やエンドポイントセキュリティ製品からも認識されます。ストリームでは、こうした問題をすべて回避できます。アプリがすでに配信しており、オリジンにもすでにバインドされているアセットサーバーを経由するためです。

既存の WebSocket 実装を移行する場合は、[WebSocket からストリームへの移行](/guides/streams-from-websockets/)に従ってください。機械的に適用できるように書かれており、気付かないまま不具合を引き起こす 3 つの相違点が最初に説明されています。

## クイックスタート

Go でストリームを宣言します。ハンドラーは接続ごとに 1 回、専用の goroutine で実行されます。

```go
app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        _ = c.Send(process(frame))  // blocks like a socket write
    }
})
```

フロントエンドから名前を指定して接続します。このオブジェクトは `WebSocket` インターフェースを実装します。

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` は `new WebSocket(url)` とまったく同様に、`readyState === CONNECTING` を持つオブジェクトを<strong>同期的に</strong>返すため、モジュールスコープで作成できます。

```js
export const Telemetry = Stream("telemetry");
```

## フレームはバイト列

各フレームは、Go では `[]byte`、JavaScript では `ArrayBuffer` です。スキーマやエンコーディングは強制されません。用途に応じて JSON、protobuf、CBOR、または生のバイト列を使用してください。

フレームは<strong>バイトストリームではなくメッセージ</strong>です。フレーム全体が到着するか、まったく到着しないかのどちらかであり、その長さも一緒に伝送されます。どちらの側も事前にサイズを把握する必要はありません。そのため、`[]byte` フィールドを持つ構造体は、マーシャリングされた結果がどのようなサイズでも、1 つのフレームとして送信されます。

## オブジェクトの送信

フレームはバイト列ですが、通常はバイト単位で考える必要はありません。両側には、対になる便利な JSON 機能があります。

```go
type Reading struct {
    Sensor string  `json:"sensor"`
    Value  float64 `json:"value"`
}

app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    var cmd map[string]any
    if err := c.ReceiveJSON(&cmd); err != nil {
        return
    }

    _ = c.SendJSON(Reading{Sensor: "cpu", Value: 42.5})
})
```

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("telemetry");
s.onopen    = () => s.send({ subscribe: "cpu" });   // stringified for you
s.onmessage = (ev) => console.log(ev.data.value);   // already an object
```

`JSONStream` は `Stream` と同じオブジェクトであり、境界でエンコードが行われるだけです。別個のプロトコルは存在せず、Go ハンドラーは両者を区別できません。有効な JSON ではないフレームは、接続を切断する代わりに `error` イベントを発生させて破棄されます。

protobuf、CBOR、バイナリ形式など、バイト列を取得して自分でエンコードしたい場合は、通常の `Stream` を使用してください。

## Go API

```go
// Register a handler. Runs once per connection, on its own goroutine.
func (a *App) HandleStream(name string, handler StreamHandler)

// The connection.
func (c *StreamConn) Send(data []byte) error        // blocks when the buffer is full
func (c *StreamConn) TrySend(data []byte) error     // ErrStreamFull instead of blocking
func (c *StreamConn) Receive() ([]byte, error)      // blocks until a frame or close
func (c *StreamConn) SendJSON(v any) error          // marshal and send as one frame
func (c *StreamConn) ReceiveJSON(v any) error       // receive one frame and unmarshal
func (c *StreamConn) Context() context.Context      // cancelled on disconnect
func (c *StreamConn) Window() Window                // nil in server mode
func (c *StreamConn) Name() string
func (c *StreamConn) Close() error
```

<strong>ハンドラーの goroutine の存続期間は、接続の存続期間と同じです。</strong>ハンドラーから戻ると接続が閉じるため、接続を開いておきたい間は `Receive`（または `c.Context()`）でブロックしてください。これは `gorilla`/`coder` WebSocket ハンドラーと同じ構造です。

エラーは `ErrStreamClosed`（ピアが切断済み）と `ErrStreamFull`（`TrySend` からのみ）です。

## JavaScript API

`Stream(name)` は、`WebSocket` の実用的なサブセットを実装するオブジェクトを返します。

| サポート対象 | 注記 |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`、`onmessage`、`onclose`、`onerror` | さらに `addEventListener` |
| `send(data)` | string、`ArrayBuffer`、型付き配列、または `Blob` — 所有権については後述 |
| `JSONStream(name)` | 同じオブジェクトで、オブジェクトを送受信 |
| `close(code, reason)` |  |
| `binaryType` | **デフォルトは`"blob"`ではなく`"arraybuffer"`** |
| `bufferedAmount` | `send` によってキューに追加され、まだ Go に到達していないバイト数 |
| `protocol`、`extensions` | 常に `""` — ネゴシエーションされない |

`binaryType` のデフォルト値は、標準との唯一の意図的な相違点です。フレームは常にバイナリであり、`Blob` にすると各メッセージを読み取るために非同期処理が 1 段階余分に必要になります。標準の動作が必要な場合は、`"blob"` に設定してください。

1 つまたは複数のウィンドウから、同じストリーム名への複数の接続を確立できます。接続ごとに独自の `StreamConn` とハンドラー goroutine が割り当てられます。

<strong>バッファーの所有権は方向によって異なります。</strong>JavaScript の `send()` は、ネイティブ WebSocket の動作と同様に、変更可能なバイナリ入力のスナップショットを同期的に作成します。そのため、呼び出し元は `send()` が戻り次第、その入力を再利用できます。Go の `Send` はスライスの所有権をトランスポートに移譲し、コピーしません。呼び出しが成功した後は、そのスライスを変更したり再利用したりしないでください。生成側でストレージを再利用する必要がある場合は、新しいスライスを渡してください。

## ライフサイクル

ストリームはソケットと同様に動作し、ソケットを閉じるイベントによってストリームも閉じます。

| イベント | 発生すること |
| --- | --- |
| ページの再読み込みまたはナビゲーション | 接続が閉じ、ハンドラーの `Receive` がエラーを返し、新しいページが新規に接続する |
| `window.close()` / ウィンドウの破棄 | そのウィンドウのすべての接続が閉じる |
| JS の `s.close()` | ハンドラーの `Receive` が `ErrStreamClosed` を返す |
| ハンドラーから戻る | フロントエンドで `onclose` が発生する |
| アプリのシャットダウン | すべての接続のコンテキストがキャンセルされる |

`WebSocket` と同様に、**自動再接続はありません**。アプリで自動再接続が必要な場合は、すでに WebSocket 用に用意している再接続ロジックを変更せずに使用できます。`onclose` でストリームを再作成してください。

## バックプレッシャー

フロントエンドの処理が追いついていない場合、ソケットの送信バッファーがいっぱいになると書き込みがブロックされるのと同様に、`Send` はブロックされます。待機せずに破棄する場合は、`TrySend` を使用してください。

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

DevTools のブレークポイント、非表示のウィンドウ、App Nap などによってフロントエンドが一時停止すると、データの収集も停止し、バッファーの上限に達した時点で生成側がブロックされます。これは意図された動作です。未読のストリームが無制限に増大するのを許すのではなく、メモリ使用量に上限を設けます。

## サーバーモード

`-tags server` を指定してビルドすると、トランスポートが `/wails/stream/ws` の **実際の WebSocket** に切り替わります。サーバーモードには、アップグレードに使用できるリスナーがすでに存在するためです。Go ハンドラーとフロントエンドコードは同一であり、アプリ側の変更は不要です。モジュールコードが実行される前に、ランタイムが適切なトランスポートを選択します。WebSocket 接続はデフォルトで同一オリジンに制限されます。意図的に別の信頼済みオリジンでフロントエンドをホストするサーバーでは、`ServerOptions.WebSocketOriginPatterns` を使用してそのホストを追加できます。

## パフォーマンス

スロットリングなしで `v3/tests/stream-performance` を使用し、約 41 百万フレームを対象に測定した結果、0 件の欠落と 0 件の順序入れ替わりが発生しました。

|  | Go→JS のピーク | JS→Go のピーク |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

ピーク値よりも特性の違いが重要です。

- **小さいフレームでは Go→JS のほうが大幅に高速です**。macOS では 634000 フレーム/秒で、逆方向は約 6200 フレーム/秒です。1 回のレスポンスには最大 256 フレームがまとめられます。JS→Go でも処理中のリクエストの後ろに蓄積したフレームをバッチ処理しますが、各接続では引き続き、それぞれの POST チェーンが直列化されます。小さいメッセージを多数送信する場合は、Go→JS を使用するか、送信前にアプリケーションレベルでまとめてください。
- **Windows では、アップロードに最適なのは 512 KB です。** それより大きいフレームは複数のリクエストに分割され、4 MB のフレームは 512 KB のフレームより<em>低速</em>という測定結果になりました。
- **レイテンシーは低く、その状態が維持されます**。macOS での p99 は約 1～2 ms で、負荷が増えても悪化しません。20000 フレーム/秒で測定した p99 は、100 フレーム/秒の場合より<em>低い</em>値でした。

プラットフォーム別の完全な表と測定方法は、この機能に付属する測定記録に記載されています。

## 制限

|  | 制限 | 上限に達した場合の動作 |
| --- | --- | --- |
| 収集待ちとしてウィンドウごとにバッファリング | 8 MB または 256 フレームのうち、先に達したほう | `Send` はブロックし、`TrySend` は `ErrStreamFull` を返す |
| 収集または書き込み待ちとしてアプリケーション全体でバッファリング | 256 MB または 8192 データフレーム | 同上 |
| `Receive` 待ちとして接続ごとに受信 | 8 MB または 256 フレーム | ハンドラーが追いつくまで、フロントエンドの `send()` が自動的に再試行される |
| `Receive` 待ちとしてアプリケーション全体で受信 | 256 MB または 8192 フレーム | 同上 |
| ウィンドウごとの接続数 | 256 | 空きスロットができるまで、オープンが自動的に再試行される |
| アプリケーション全体の有効な接続数 | 4096 | 同上 |
| ウィンドウごとのセッション数 | 16 | 再読み込み時は同じウィンドウの古いセッションが置き換えられ、それ以外の場合はオープンが再試行される |
| いずれかの方向の単一フレーム | 64 MB | Go は `ErrStreamTooLarge` を返し、JS からの場合はストリームが `error` を発生させて閉じる |
| ストリーム名 | 256 UTF-8 バイト | オープンが拒否され、ストリームが `error` を発生させる |
| アイドル時のポーリング保持時間 | 20 秒 | ポーリングは空の結果を返し、ランタイムが直ちに再発行する |

上記のいずれの場合も、データが通知なく破棄されることはありません。<em>自動的に再試行される</em>と記載された 2 行は通常のバックプレッシャーです。ランタイムはフレームを保持し、短いバックオフを挟んで再試行するため、コードからはエラーではなく速度が低下したストリームとして認識されます。`error` を発生させる行は負荷ではなくプログラミングミスを示しており、問題を覆い隠さずに表面化させます。

現時点では、これらはオプションではなくコンパイル時定数です。変更する必要がある場合は、内部構造ガイドを参照してください。

## ストリームを使用すべきでない場合

- **リクエスト／レスポンスにはバインディングを使用してください。** ストリームは継続的なデータや要求なしに送られるデータに使用します。値を返す呼び出しは、バインドされたメソッドとして実装するほうが簡単です。
- **アプリケーションイベントには `Emit`/`On` を使用してください。** イベントはすべてのリスナーに配信される、独立した実績あるシステムです。ストリームはポイントツーポイントです。
- **`InitialHTML` ウィンドウでは使用できません。** これらは `origin === "null"` を使用して読み込まれるため、アセットサーバーには一切アクセスできません。
