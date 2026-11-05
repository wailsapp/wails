---
title: "WebSocket から Streams への移行"
description: "既存の WebSocket 実装を Wails streams に段階的に移行する手順と、気付かないまま不具合を引き起こす相違点"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

現在 WebSocket サーバーを実行しているアプリを [Streams](/guides/streams/) に移行するための機械的な手順です。エージェントを含め、この手順どおりに作業できるように記述されています。

移行の利点は、リスナーを削除できることです。バインドされた TCP ポートも、オリジンチェックも、トークンもなくなり、ファイアウォールやエンドポイントセキュリティ製品が問題視するものは何も残りません。API は十分に似ているため、フロントエンドコードの大半は変更不要です。しかし、<strong>気付かないまま不具合を引き起こす相違点が 3 つ</strong>あります。これらは午後の時間を丸ごと奪いかねないため、最初に説明します。

## 一般的なケース：回避策としてのローカル HTTP サーバー

Wails アプリで WebSocket が使われる一般的な理由は、継続的なフィードをフロントエンドにプッシュする方法がほかになかったことです。そのため、アプリはローカルポートで独自の `http.Server` を起動し、フロントエンドからそこへ接続します。この構成であれば、移行によってサーバー自体を完全に削除でき、それを<em>取り巻く</em>複数の仕組みも一緒に削除できます。

<strong>ポート検出の仕組みは不要になります。</strong>接続先のポートを、何らかの方法でフロントエンドに伝える必要がありました。バインドされた `GetServerPort()`、使用中の場合にフォールバックする固定ポート、注入されたグローバル変数、または `localStorage` に保存した値などです。これらはすべて不要になります。stream は名前で指定し、その名前は両側でコンパイル時定数となります。

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

<strong>CORS 設定は不要になります。</strong>webview のオリジンは、プラットフォームに応じて `wails://` または `http://wails.localhost` となるため、ローカルサーバーには `CheckOrigin`、`Access-Control-Allow-Origin` ヘッダー、またはその両方が必要です。Streams はページの読み込み元であるアセットサーバーを経由するため、許可すべきクロスオリジンリクエストはありません。

<strong>独自に設けた認証トークンは不要になります。</strong>localhost にバインドされたポートには、そのマシン上のすべてのプロセスから到達できます。そのため、慎重な実装では、ほかのソフトウェアからの接続を防ぐためにトークンまたは nonce を追加します。移行後は、到達可能なポート自体がなくなります。

<strong>WebSocket 以外のエンドポイントは、アセットサーバーのミドルウェアへ移します。</strong>この種のサーバーが WebSocket 専用のままであることはほとんどなく、ファイルのダウンロード、画像エンドポイント、ヘルスチェックなどがソケットと並んで増えていく傾向があります。Streams はそれらを置き換えませんが、そのために別のサーバーを用意する必要もありません。同じハンドラーをアセットサーバーにマウントします：

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

その後、フロントエンドは `/api/...` を同一オリジンの相対 URL として呼び出します。ホストもポートも CORS も不要です。Streams とこの方法を組み合わせれば、ローカルサーバーが担う処理は何も残りません。

## 開始前に必ずお読みください

### 1. `ev.data` は文字列ではなく、必ず `ArrayBuffer` です

これが最も重要な相違点です。WebSocket はテキストメッセージを文字列として渡しますが、stream はすべてのメッセージをバイト列として渡します。次のようなコードは、**コンパイルも実行もできますが、正しく動作しません**：

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

各呼び出し箇所ではなく、境界部分で修正します：

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

また、通常そうであるように通信内容が JSON なら、`Stream` の代わりに `JSONStream` を使用すれば、この問題を完全に回避できます：

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

これが JSON WebSocket を移行する最短の方法です。コンストラクターを置き換え、`JSON.parse` と `JSON.stringify` の呼び出しを削除すれば、ハンドラーの残りの部分は変更不要です。JSON 以外の通信では、一度だけラップして、すべてのハンドラーを変更せずに残します。[互換性 shim](#-shim)を参照してください。

### 2. 送信には互換性がありますが、受信にはありません

`send()` は文字列を受け取り、UTF-8 としてエンコードするため、`s.send(JSON.stringify(x))` は変更せずに動作します。編集が必要なのは受信処理だけです。コードの半分はそのまま動作し続けるため、この非対称性は見落としやすい点です。

### 3. URL はありません

WebSocket は、パス、クエリ文字列、サブプロトコル、認証トークンなどの接続パラメーターを URL に含めます。stream が持つのは名前だけです。URL で渡していたものは、すべて最初のフレームへ移すか、接続前に呼び出すバインド済みメソッドへ移す必要があります。

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go 側

HTTP サーバー、upgrader、接続レジストリを削除します。それぞれハンドラーに置き換わります。

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Stream |
| --- | --- |
| `upgrader.Upgrade` / `websocket.Accept` | *（なし — `HandleStream` が登録処理のすべてです）* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`、または単にハンドラーから return |
| ブロードキャスト用の接続レジストリ | 独自のものを維持 — [ブロードキャスト](#heading-2)を参照 |
| `http.ListenAndServe`、mux、オリジンチェック、トークン | **削除** |
| ping/pong によるキープアライブ | **削除** — 維持すべきアイドル状態のソケットはありません |
| `r.Context()` | `c.Context()` |

ハンドラー goroutine の存続期間は接続の存続期間と一致します。gorilla のハンドラーとまったく同じなので、既存のループ構造を変更せずに引き継げます。

## フロントエンド側

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`、4 つの状態定数、`addEventListener`、`close(code, reason)`、`bufferedAmount` はすべて、`WebSocket` の場合と同じように動作します。

### 互換性 shim

ハンドラーにまったく手を加えたくない場合は、コンストラクターを一度ラップします。これで、文字列メッセージを前提とする既存のコードをそのまま使用できます。

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

これで、`const ws = TextStream("feed")`を`new WebSocket(url)`の代わりにそのまま使用できます。

## ブロードキャスト

WebSocket サーバーは通常、複数の接続先へ配信できるようにレジストリを保持します。ストリームにはブロードキャスト機能が組み込まれていないため、レジストリは残し、`*websocket.Conn`の代わりに`*StreamConn`を格納します。

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

守るべきルールは2つあります。送信前にロックを解放することと、単一の低速なコンシューマーがほかのすべてのクライアントを停止させないよう、複数の接続先への配信では`TrySend`を優先することです。

## もう1つのケース：フロントエンドがブローカーと直接通信する場合

Wails アプリではあまり一般的ではありませんが、知っておく価値があります。フロントエンドが<strong>アプリではなくブローカーへの</strong> WebSocket 接続を開く場合（`nats.ws` を使った NATS サーバーへの接続や、WebSocket 経由の MQTT など）、ストリームをそのまま代用することはできません。ストリームがフロントエンドを接続する先はサードパーティーではなく、<em>自分の Go コード</em>だからです。

この移行ではアーキテクチャの変更が必要ですが、通常は望ましい変更です。

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

ブラウザー用ライブラリよりネイティブライブラリの方が優れている Go 側へブローカークライアントを移し、フロントエンドに必要な部分をストリーム経由で公開します。

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

サブスクリプションのコールバック内では`TrySend`を使用している点に注意してください。このコールバックはブローカークライアントの goroutine で実行されるため、ここでブロックすると、その接続上のすべてのサブスクリプションへの配信が停止します。

この構成には、ブローカーの認証情報がフロントエンドに渡らない、マシンに WebSocket ポートを公開しなくてよい、再接続とバックオフをブラウザー用クライアントではなく成熟した Go クライアントで処理できる、という利点があります。

## 移行チェックリスト

- [ ] 既存の各 WebSocket エンドポイントに`HandleStream`を登録した
- [ ] 読み取りループを変換した：`ReadMessage`/`Read` → `c.Receive()`
- [ ] 書き込みを変換した：`WriteMessage` → `c.Send()`。複数の接続先への配信またはブローカーのコールバックでは`TrySend`を使用した
- [ ] HTTP サーバー、mux エントリー、アップグレーダー、オリジンチェック、認証トークンを<strong>削除した</strong>
- [ ] ping/pong によるキープアライブを<strong>削除した</strong>
- [ ] URL パラメーターを最初のフレームまたはバインドされたメソッドへ移した
- [ ] **すべての `ev.data` 読み取り結果をデコードした**（`new TextDecoder().decode(ev.data)`）、または互換性シムを導入した
- [ ] 再接続ロジックをそのまま残した（設計上、再接続機能は組み込まれていない）
- [ ] フロントエンドがブローカーと直接通信していた場合は、ブローカークライアントを Go 側へ移した
- [ ] `InitialHTML`ウィンドウがないか確認した。このウィンドウではストリームを一切使用できない

変換時に grep で検索しておくべきもの：

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## 想定される動作の違い

|  | WebSocket | ストリーム |
| --- | --- | --- |
| メッセージ型 | テキストまたはバイナリ | バイト列のみ |
| `ev.data` | 文字列または`Blob`/`ArrayBuffer` | 常に`ArrayBuffer`（`binaryType = "blob"`を使用する場合を除く） |
| `binaryType`のデフォルト値 | `"blob"` | `"arraybuffer"` |
| サブプロトコル、`extensions` | ネゴシエーションされる | 未対応。常に`""` |
| 接続パラメーター | URL とクエリ | 最初のフレームまたはバインドされた呼び出し |
| 認証 | トークンまたは Cookie | 不要 — 呼び出し元はアプリだけ |
| キープアライブ | ping/pong | 不要 |
| 自動再接続 | なし | なし（同じ） |
| クローズコード | 全範囲 | `1000` 正常、`1001` セッション終了、`1002` フレーミング不一致、`1006` エラー |
| バックプレッシャー | カーネルのソケットバッファ | ウィンドウあたり 8 MB / 256 フレーム、その後 `Send` がブロックされる |
| 複数の接続、単一のエンドポイント | はい | はい |

## 移行後

よくある間違いを検出するための健全性チェック：

1. ページを繰り返し再読み込みします。再読み込みのたびにハンドラーが終了して新しいハンドラーが起動し、決して蓄積しないことを確認してください。
2. 各方向に 512 KB を超えるメッセージを送信します。
3. 1 分を超えてアイドル状態にしてから、再接続せずに通信が再開することを確認してください。
4. 開発者ツールを開き、負荷がかかった状態でブレークポイントにより一時停止してから再開します。プロデューサーがデータを失ったり際限なく増大したりせず、ブロックされた後に回復することを確認してください。
