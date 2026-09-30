---
title: "串流"
description: "Go 與 JavaScript 之間的雙向位元組串流，採用 WebSocket 程式設計模型，且不需監聽通訊端"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

串流在 Go 與前端之間提供具名、有序的雙向位元組通道，採用與 WebSocket 相同的程式設計模型，卻<strong>不需繫結 TCP 連接埠</strong>。

WebSocket 無法透過自訂 URL 協定通訊，因此，要在 WebView 內使用 WebSocket，唯一的方法就是執行真正的 HTTP 伺服器並監聽連接埠。對桌面應用程式而言，這代表開放一個本機連接埠，電腦上的任何其他處理程序都能存取；為確保安全，還需要檢查來源並使用權杖，而且使用者執行的每項防火牆與端點安全性產品都能偵測到它。串流可避免這一切：它使用應用程式現有的資產伺服器，而該伺服器已限定來源。

要遷移現有的 WebSocket 實作嗎？請依照[將 WebSocket 遷移至串流](/guides/streams-from-websockets/)操作。這份指南可供逐步照做，並先說明三項會在未發出錯誤的情況下造成故障的差異。

## 快速入門

在 Go 中宣告串流。每個連線都會執行一次處理常式，且各自在其 goroutine 上執行：

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

從前端依名稱連線。此物件會實作`WebSocket`介面：

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)`會像`new WebSocket(url)`一樣，以`readyState === CONNECTING`<strong>同步</strong>傳回，因此你可以在模組範圍建立它：

```js
export const Telemetry = Stream("telemetry");
```

## 訊框是位元組

每個訊框在 Go 中都是`[]byte`，在 JavaScript 中則是`ArrayBuffer`。系統不會強制使用任何結構描述或編碼，你可以依偏好使用 JSON、protobuf、CBOR 或原始位元組。

訊框是<strong>訊息，而不是位元組串流</strong>：它不是完整抵達，就是完全不會抵達，且其長度會隨附其中。任一端都不需要預先知道大小，因此含有`[]byte`欄位的結構在封送處理後，不論結果為何，都會以單一訊框傳送。

## 傳送物件

訊框是位元組，但你通常不會想以位元組的角度思考。兩端都提供可相互搭配的 JSON 便利功能：

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

`JSONStream`與`Stream`是同一個物件，只是在邊界完成編碼；沒有獨立的通訊協定，Go 處理常式也無法分辨差異。不是有效 JSON 的訊框會引發`error`事件並遭到捨棄，而不會中斷連線。

若需要直接處理位元組，請使用一般的`Stream`：例如 protobuf、CBOR、二進位格式，或任何你偏好自行編碼的資料。

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

<strong>處理常式 goroutine 的存續期與連線相同。</strong>從處理常式傳回會關閉連線，因此，只要你想讓連線保持開啟，就應持續封鎖於`Receive`（或`c.Context()`）。這與`gorilla`/`coder` WebSocket 處理常式的形式相同。

錯誤包括`ErrStreamClosed`（對等端已離線）與`ErrStreamFull`（僅由`TrySend`傳回）。

## JavaScript API

`Stream(name)`會傳回一個實作`WebSocket`實用子集的物件：

| 支援項目 | 備註 |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`、`onmessage`、`onclose`、`onerror` | 另加`addEventListener` |
| `send(data)` | 字串、`ArrayBuffer`、型別陣列或`Blob`——請參閱下方的擁有權說明 |
| `JSONStream(name)` | 同一個物件，輸入與輸出皆為物件 |
| `close(code, reason)` |  |
| `binaryType` | **預設為`"arraybuffer"`**，而非`"blob"` |
| `bufferedAmount` | 由`send`排入佇列、但尚未抵達 Go 的位元組 |
| `protocol`、`extensions` | 一律為`""`，不會進行協商 |

`binaryType`的預設值是唯一刻意偏離標準之處：訊框一律為二進位，而`Blob`會迫使每則訊息額外經過一次非同步步驟才能讀取。若要使用標準行為，請將其設為`"blob"`。

允許從同一個或多個視窗，建立多個連至相同串流名稱的連線。每個連線都有自己的`StreamConn`與處理常式 goroutine。

<strong>緩衝區擁有權會因方向而異。</strong>JavaScript 的`send()`會同步建立可變二進位輸入的快照，以符合原生 WebSocket 的行為，因此呼叫端可在`send()`傳回後立即重複使用這些輸入。Go 的`Send`會將其切片的擁有權移交給傳輸層，且不會複製；成功呼叫後，請勿修改或重複使用該切片。若產生器需要重複使用其儲存空間，請移交新的切片。

## 生命週期

串流的行為如同通訊端，而會關閉通訊端的事件也會關閉串流：

| 事件 | 結果 |
| --- | --- |
| 頁面重新載入或導覽 | 連線關閉，處理常式的`Receive`傳回錯誤，而新頁面會建立全新連線 |
| `window.close()`／視窗遭到銷毀 | 該視窗的所有連線都會關閉 |
| JS 中的`s.close()` | 處理常式的`Receive`傳回`ErrStreamClosed` |
| 處理常式傳回 | 前端收到`onclose` |
| 應用程式關閉 | 每個連線的 context 都會被取消 |

系統<strong>不會自動重新連線</strong>，這與`WebSocket`一致。如果應用程式需要自動重新連線，現有的 WebSocket 重新連線邏輯無需修改即可使用，只要在`onclose`中重新建立串流即可。

## 背壓

當前端來不及處理時，`Send`會封鎖，就像通訊端寫入在傳送緩衝區已滿時會封鎖一樣。若你偏好捨棄資料而不是等待，請使用`TrySend`：

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

暫停的前端（例如停在開發人員工具中斷點、視窗遭到隱藏或進入 App Nap）會停止收集資料，之後緩衝區上限便會封鎖產生器。這是刻意設計：它會限制記憶體用量，而不是讓無人讀取的串流無限制增長。

## 伺服器模式

使用`-tags server`建置時，傳輸方式會改為位於`/wails/stream/ws`的<strong>真正 WebSocket</strong>，因為伺服器模式已有可供升級的監聽器。Go 處理常式與前端程式碼完全相同，應用程式不必做任何變更。執行階段會在任何模組程式碼執行前替你選擇傳輸方式。WebSocket 連線預設為同源。若伺服器刻意將前端託管於另一個受信任的來源，可以使用`ServerOptions.WebSocketOriginPatterns`新增該主機。

## 效能

使用`v3/tests/stream-performance`在不節流的情況下測量，約41百萬個訊框中有0次遺失及0次重新排序：

|  | Go→JS 峰值 | JS→Go 峰值 |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

效能曲線比峰值更重要：

- **對小型訊框而言，Go→JS 快得多**——在 macOS 上為634000訊框/秒，反方向則約為6200訊框/秒。單一回應最多會合併256個訊框；JS→Go 也會批次處理等待進行中請求完成時累積的訊框，但每個連線仍須依序執行自己的 POST 請求。若要傳送許多小型訊息，請優先使用 Go→JS，或在向 Go 傳送前先於應用程式層級將其批次處理。
- <strong>在 Windows 上，512 KB 是上傳的最佳大小。</strong>超過此大小的訊框會拆分成多個請求，而4 MB 訊框的測量結果比512 KB 訊框<em>慢</em>。
- **延遲很低且持續維持低水準**：macOS 上的 p99 約為1–2 ms，而且不會隨負載增加而惡化——在每秒20000個訊框下測得的 p99<em>低於</em>每秒100個訊框下的 p99。

各平台的完整表格與測量方法，請參閱此功能隨附的測量記錄。

## 限制

|  | 限制 | 達到限制時會發生什麼情況 |
| --- | --- | --- |
| 每個視窗中已緩衝、等待收集的資料 | 8 MB 或256個訊框，以先達到者為準 | `Send`會阻塞；`TrySend`會傳回`ErrStreamFull` |
| 整個應用程式中已緩衝、等待收集或寫入的資料 | 256 MB 或8192個資料訊框 | 同上 |
| 每個連線中已接收、等待`Receive`的資料 | 8 MB 或256個訊框 | 系統會替你重試前端的`send()`，直到處理常式跟上為止 |
| 整個應用程式中已接收、等待`Receive`的資料 | 256 MB 或8192個訊框 | 同上 |
| 每個視窗的連線數 | 256 | 系統會替你重試開啟操作，直到有可用名額為止 |
| 整個應用程式中尚未關閉的連線數 | 4096 | 同上 |
| 每個視窗的工作階段數 | 16 | 重新載入會取代其本身較舊的工作階段；否則會重試開啟操作 |
| 任一方向的單一訊框 | 64 MB | Go 會傳回`ErrStreamTooLarge`；若來自 JS，串流會引發`error`並關閉 |
| 串流名稱 | 256個 UTF-8位元組 | 開啟操作會遭拒，而且串流會引發`error` |
| 閒置輪詢保持時間 | 20 s | 輪詢會傳回空結果，執行階段隨即重新發出輪詢 |

上述情況都不會無聲地捨棄資料。標示<em>系統會替你重試</em>的兩列是一般的背壓——執行階段會保留訊框，並在短暫退避後重試，因此程式碼看到的是速度較慢的串流，而不是錯誤。會引發`error`的各列代表程式設計錯誤，而非負載問題；這些錯誤會直接呈現，不會被掩蓋。

目前這些是編譯時期常數，不是選項。若需要變更，請參閱內部機制指南。

## 何時不該使用串流

- <strong>對於請求/回應，請使用繫結。</strong>串流適用於連續或未經請求的資料；會傳回值的呼叫以繫結方法實作會更簡單。
- <strong>對於應用程式事件，請使用`Emit`/`On`。</strong>事件會分送給每個接聽程式，且屬於另一套成熟且廣泛使用的系統。串流則是點對點的。
- <strong>不適用於`InitialHTML`視窗。</strong>這些視窗使用`origin === "null"`載入，因此完全無法存取資產伺服器。
