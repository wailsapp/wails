---
title: "將 WebSocket 遷移至 Streams"
description: "逐步將現有 WebSocket 實作轉換為 Wails streams，包括不會立即報錯、卻會造成問題的差異"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

這是一份將目前執行 WebSocket 伺服器的應用程式轉換為 [Streams](/guides/streams/) 的機械式指南。內容可供逐字照做，包括由代理程式執行。

這項遷移的好處是可以刪除監聽器：不再有繫結的 TCP 連接埠、來源檢查或權杖，也沒有任何會引起 防火牆或端點安全性產品疑慮的項目。API 足夠相似，因此大多數前端程式碼都不需要變更——但 **有三項差異不會立即報錯、卻會造成問題**。這些差異列在最前面，因為它們會讓你浪費一整個下午。

## 常見情況：作為因應措施的本機 HTTP 伺服器

Wails 應用程式通常會使用 WebSocket，是因為沒有其他方式能將 連續資料流推送至前端，所以應用程式會在本機連接埠上啟動自己的 `http.Server`， 再由前端連回該伺服器。如果你的架構屬於這種情況，這項遷移會直接刪除伺服器——連同你在其<em>周邊</em>建置的數項機制。

<strong>連接埠探索機制會消失。</strong>你必須以某種方式告訴前端要連線至哪個連接埠：繫結的 `GetServerPort()`、被占用時可改用其他連接埠的固定連接埠、 注入的全域變數，或儲存在 `localStorage` 中的值。這些全都會消失——stream 以名稱定址， 而且兩端的名稱都是編譯期常數。

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

<strong>CORS 設定會消失。</strong>依平台而定，webview 的來源是 `wails://` 或 `http://wails.localhost`，因此本機伺服器需要 `CheckOrigin`、`Access-Control-Allow-Origin` 標頭，或兩者都需要。Streams 會透過載入頁面的資產伺服器傳輸， 因此不需要允許跨來源要求。

<strong>自行設計的任何驗證權杖都會消失。</strong>繫結至 localhost 的連接埠可由機器上的每個處理程序存取， 因此嚴謹的實作會加入權杖或 nonce，防止其他軟體連線。現在已不再有可供存取的連接埠。

<strong>非 WebSocket 端點會移至資產伺服器中介軟體。</strong>這類伺服器很少能一直只處理 WebSocket——檔案下載、圖片端點和健康情況檢查通常會逐漸堆積在 socket 旁。 Streams 不會取代這些端點，但你也不需要為它們另設第二部伺服器。 將相同的處理常式掛載至資產伺服器：

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

接著，前端會以同源相對 URL 呼叫 `/api/...`——不需要主機、不需要連接埠，也不需要 CORS。結合 streams 與這項變更後，本機伺服器便再也沒有任何工作可做。

## 開始之前請先閱讀

### 1. `ev.data` 是 `ArrayBuffer`，絕不是字串

這是最重要的一點。WebSocket 會以字串傳遞文字訊息；stream 則會以位元組傳遞每則 訊息。像這樣的程式碼<strong>可以編譯、可以執行，卻會出現錯誤行為</strong>：

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

請在邊界處修正，而不是在每個呼叫處修正：

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

或者，如果傳輸的是 JSON（通常都是），請使用 `JSONStream`，不要使用 `Stream`， 即可完全避開這個問題：

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

這是遷移 JSON WebSocket 最簡短的方式：替換建構函式、移除 `JSON.parse` 和 `JSON.stringify` 呼叫，其餘處理常式維持不變。若傳輸的不是 JSON，請只包裝一次，讓所有處理常式保持不變——請參閱 [相容性轉接層](#heading-2)。

### 2. 傳送相容；接收不相容

`send()` 接受字串並將其編碼為 UTF-8，因此 `s.send(JSON.stringify(x))` 無須變更即可運作。 只有接收路徑需要編輯。這種不對稱性很容易被忽略，因為一半的程式碼仍可繼續運作。

### 3. 沒有 URL

WebSocket 會在 URL 中攜帶連線參數——路徑、查詢字串、子通訊協定、 驗證權杖。stream 只有名稱。原先透過 URL 傳遞的任何資料，都必須移至 第一個 frame，或移至連線前呼叫的繫結方法。

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go 端

刪除 HTTP 伺服器、upgrader 和連線登錄。每一項都會由處理常式取代。

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
| `upgrader.Upgrade` / `websocket.Accept` | *（不需要任何內容——`HandleStream` 本身就是完整的註冊）* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`，或直接從處理常式返回 |
| 用於廣播的連線登錄 | 保留你自己的連線登錄——請參閱[廣播](#heading-3) |
| `http.ListenAndServe`、mux、來源檢查、權杖 | **刪除** |
| ping/pong 保活機制 | **刪除**——沒有閒置 socket 需要保持連線 |
| `r.Context()` | `c.Context()` |

處理常式 goroutine 的生命週期就是連線的生命週期，與 gorilla 處理常式完全相同， 因此現有迴圈的結構可以原封不動地沿用。

## 前端

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

`readyState`、四個狀態常數、`addEventListener`、`close(code, reason)` 和 `bufferedAmount` 的行為全都與 `WebSocket` 上相同。

### 相容性轉接層

如果完全不想修改處理常式，只要包裝一次建構函式即可。這樣，原本預期接收字串訊息的現有程式碼便可原封不動地運作：

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

如此一來，`const ws = TextStream("feed")`便可直接取代`new WebSocket(url)`。

## 廣播

WebSocket 伺服器通常會維護一份登錄表，以便將訊息分送給多個用戶端。串流沒有內建廣播功能，因此請保留該登錄表，但改為儲存`*StreamConn`，而不是`*websocket.Conn`：

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

有兩項規則值得保留：傳送前先釋放鎖定；進行分送時優先使用`TrySend`，以免單一緩慢的消費端拖住其他所有用戶端。

## 另一種情況：前端直接與訊息代理程式通訊

這在 Wails 應用程式中較少見，但仍值得瞭解。如果前端開啟的 WebSocket **是連至訊息代理程式，而不是您的應用程式**，例如`nats.ws`連至 NATS 伺服器或透過 WebSocket 使用 MQTT，那麼串流便無法直接替代它，因為串流是將前端連至<em>您的 Go 程式碼</em>，而不是第三方。

這項移轉會改變架構，而且通常是有益的改變：

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

將訊息代理程式用戶端移至 Go；原生程式庫比瀏覽器程式庫更好用，再透過串流公開前端所需的部分：

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

請注意訂閱回呼中的`TrySend`：該回呼會在訊息代理程式用戶端的 goroutine 上執行；若將它阻塞，該連線上每項訂閱的訊息傳遞都會停滯。

這樣做的好處是：訊息代理程式的認證資訊絕不會傳至前端、不必向本機開放 WebSocket 連接埠，而且重新連線與退避會由成熟的 Go 用戶端處理，而非瀏覽器用戶端。

## 移轉檢查清單

- [ ] 已為原有的每個 WebSocket 端點註冊`HandleStream`
- [ ] 已轉換讀取迴圈：`ReadMessage`/`Read` → `c.Receive()`
- [ ] 已轉換寫入操作：`WriteMessage` → `c.Send()`；在任何分送或訊息代理程式回呼中則改用`TrySend`
- [ ] HTTP 伺服器、多工器項目、升級器、來源檢查及驗證權杖皆已<strong>刪除</strong>
- [ ] ping/pong 保持連線機制已<strong>刪除</strong>
- [ ] 已將 URL 參數移至第一個訊框或繫結方法
- [ ] **每次讀取`ev.data`時皆已解碼**（`new TextDecoder().decode(ev.data)`），或已採用相容層
- [ ] 重新連線邏輯維持不變（依設計，沒有內建重新連線功能）
- [ ] 若前端原本直接與訊息代理程式通訊，已將其用戶端移至 Go
- [ ] 已檢查`InitialHTML`視窗；這些視窗完全無法使用串流

轉換期間值得搜尋的項目：

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## 預期會出現的行為差異

|  | WebSocket | 串流 |
| --- | --- | --- |
| 訊息類型 | 文字或二進位 | 僅位元組 |
| `ev.data` | 字串或`Blob`/`ArrayBuffer` | 一律為`ArrayBuffer`（除非`binaryType = "blob"`） |
| `binaryType`預設值 | `"blob"` | `"arraybuffer"` |
| 子通訊協定、`extensions` | 經協商決定 | 不支援；一律為`""` |
| 連線參數 | URL 與查詢參數 | 第一個訊框或繫結呼叫 |
| 驗證 | 權杖或 Cookie | 不需要；只有應用程式會呼叫 |
| 保持連線 | ping/pong | 不需要 |
| 自動重新連線 | 無 | 無（相同） |
| 關閉代碼 | 完整範圍 | `1000`正常、`1001`工作階段已關閉、`1002`訊框不符、`1006`錯誤 |
| 背壓 | 核心通訊端緩衝區 | 每個視窗8 MB／256個訊框，之後`Send`會阻塞 |
| 多個連線，一個端點 | 是 | 是 |

## 遷移後

可找出常見錯誤的基本檢查：

1. 反覆重新載入頁面——處理常式每次都應結束，並啟動新的處理常式，絕不能不斷累積。
2. 在兩個方向各傳送一則大於512 KB的訊息。
3. 讓它閒置超過一分鐘；資料傳輸應能恢復，而無需重新連線。
4. 開啟開發者工具，在有負載的情況下於中斷點暫停，然後繼續執行——生產者應阻塞並恢復，而不是遺失資料或無限制地增長。
