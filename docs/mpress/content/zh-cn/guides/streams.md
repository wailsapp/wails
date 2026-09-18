---
title: "流"
description: "Go 与 JavaScript 之间的双向字节流，采用 WebSocket 编程模型且无需监听套接字"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

流在 Go 与前端之间提供一个有名称、有序的双向字节通道，其编程模型与 WebSocket 相同——**但无需绑定 TCP 端口**。

WebSocket 无法通过自定义 URL 方案通信，因此要在 webview 中使用 WebSocket，唯一的方法是运行真正的 HTTP 服务器并监听端口。在桌面应用中，这意味着需要开放一个本地端口，计算机上的任何其他进程都能访问该端口；为确保安全，还需进行来源检查并使用令牌，而且用户运行的所有防火墙和端点安全产品都能看到该端口。流避免了所有这些问题：它们使用应用已经提供的资源服务器，而该服务器已绑定到来源。

要迁移现有的 WebSocket 实现吗？请按照[将 WebSocket 迁移到流](/guides/streams-from-websockets/)操作——该指南可供你按步骤直接实施，并首先说明了三个会在不发出任何提示的情况下导致故障的差异。

## 快速入门

在 Go 中声明一个流。每个连接都会在独立的 goroutine 上运行一次处理程序：

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

在前端按名称连接。该对象实现了`WebSocket`接口：

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)`会<strong>同步</strong>返回一个初始状态为`readyState === CONNECTING`的对象，与`new WebSocket(url)`完全一样，因此可以在模块作用域中创建该对象：

```js
export const Telemetry = Stream("telemetry");
```

## 帧是字节

每个帧在 Go 中都是`[]byte`，在 JavaScript 中都是`ArrayBuffer`。系统不会强制使用任何数据模式或编码方式——你可以按需使用 JSON、protobuf、CBOR 或原始字节。

帧是<strong>消息，而不是字节流</strong>：它要么完整到达，要么完全不到达，并且自身携带长度。双方都无需预先知道大小，因此包含`[]byte`字段的结构体无论被编组为何种结果，都会作为一个帧发送。

## 发送对象

帧是字节，但通常无需以字节为单位思考。双方都提供了相互配套的 JSON 便捷功能：

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

`JSONStream`与`Stream`是同一个对象，只是在边界处完成编码——不存在单独的协议，Go 处理程序也无法分辨二者。不是有效 JSON 的帧会触发`error`事件并被丢弃，而不会导致连接断开。

如果需要直接使用字节，请使用普通的`Stream`：例如 protobuf、CBOR、二进制格式，或任何你希望自行编码的内容。

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

<strong>处理程序 goroutine 的生命周期与连接的生命周期相同。</strong>处理程序返回时会关闭连接，因此只要希望连接保持打开，就应在`Receive`（或`c.Context()`）上阻塞。这与`gorilla`/`coder` WebSocket 处理程序的结构相同。

错误包括`ErrStreamClosed`（对端已断开）和`ErrStreamFull`（仅由`TrySend`返回）。

## JavaScript API

`Stream(name)`返回一个实现了`WebSocket`中实用功能子集的对象：

| 支持的功能 | 说明 |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`、`onmessage`、`onclose`、`onerror` | 以及`addEventListener` |
| `send(data)` | 字符串、`ArrayBuffer`、类型化数组或`Blob`——所有权说明见下文 |
| `JSONStream(name)` | 同一个对象，输入和输出均为对象 |
| `close(code, reason)` |  |
| `binaryType` | **默认为`"arraybuffer"`**，而不是`"blob"` |
| `bufferedAmount` | 由`send`加入队列但尚未到达 Go 的字节数 |
| `protocol`、`extensions` | 始终为`""`——不进行协商 |

`binaryType`的默认值是唯一有意偏离标准的地方：帧始终为二进制，而`Blob`会强制增加一次异步跳转来读取每条消息。如果需要标准行为，请将其设置为`"blob"`。

允许从一个或多个窗口建立多个使用相同流名称的连接。每个连接都有自己的`StreamConn`和处理程序 goroutine。

<strong>缓冲区所有权因传输方向而异。</strong>JavaScript 的`send()`会同步创建可变二进制输入的快照，这与原生 WebSocket 的行为一致，因此调用方可在`send()`返回后立即复用这些输入。Go 的`Send`会将其切片的所有权移交给传输层且不会复制；调用成功后，请勿修改或复用该切片。如果生产者需要复用其存储空间，请移交一个新切片。

## 生命周期

流的行为与套接字类似，会关闭套接字的事件也会关闭流：

| 事件 | 结果 |
| --- | --- |
| 页面重新加载或导航 | 连接关闭，处理程序的`Receive`返回错误，新页面建立全新连接 |
| `window.close()`/窗口被销毁 | 该窗口的所有连接均关闭 |
| JS 中的`s.close()` | 处理程序的`Receive`返回`ErrStreamClosed` |
| 处理程序返回 | 前端收到`onclose` |
| 应用关闭 | 取消每个连接的上下文 |

系统<strong>不会自动重新连接</strong>，这与`WebSocket`一致。如果应用需要自动重连，现有的 WebSocket 重连逻辑无需修改即可使用——在`onclose`中重新创建流。

## 背压

当前端处理不及时，`Send`会阻塞，就像套接字写入操作在发送缓冲区已满时会阻塞一样。如果希望丢弃数据而不是等待，请使用`TrySend`：

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

当前端暂停时——例如停在开发者工具断点、窗口被隐藏或 App Nap 生效——便会停止收集数据，此时缓冲区上限会使生产者阻塞。这是有意设计的行为：它会限制内存用量，而不是让无人读取的流无限增长。

## 服务器模式

使用`-tags server`构建时，传输方式会切换为位于`/wails/stream/ws`的<strong>真正的 WebSocket</strong>，因为服务器模式已有可供升级的监听器。Go 处理程序和前端代码完全相同——应用无需做任何更改。运行时会在任何模块代码运行之前为你选择传输方式。WebSocket 连接默认同源。如果服务器有意将前端托管在另一个受信任的源上，可以使用`ServerOptions.WebSocketOriginPatterns`添加该主机。

## 性能

使用`v3/tests/stream-performance`在不限速条件下测量，跨约41百万个帧时，丢帧0次，乱序0次：

|  | Go→JS 峰值 | JS→Go 峰值 |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

性能曲线比峰值更重要：

- **对于小帧，Go→JS 要快得多**——在 macOS 上可达634000帧/秒，而反方向约为6200帧/秒。单个响应最多可合并256个帧；JS→Go 也会批量处理在请求进行期间积压的帧，但每个连接仍会将自己的 POST 请求链串行化。如果要发送许多小消息，请优先使用 Go→JS，或先在应用层将它们批量合并再上传。
- <strong>在 Windows 上，上传时512 KB 是最佳大小。</strong>超过该大小的帧会拆分为多个请求；测得4 MB 帧比512 KB 帧<em>更慢</em>。
- **延迟很低，并且会一直保持在低水平**：macOS 上的 p99 约为1–2 ms，而且不会随负载增加而恶化——测得20000帧/秒时的 p99 比100帧/秒时<em>更低</em>。

各平台的完整表格和测量方法请参阅此功能随附的测量记录。

## 限制

|  | 限制 | 达到限制时会发生什么 |
| --- | --- | --- |
| 每个窗口中等待收集的缓冲数据 | 8 MB 或256个帧，以先达到者为准 | `Send`会阻塞；`TrySend`返回`ErrStreamFull` |
| 整个应用中等待收集或写入的缓冲数据 | 256 MB 或8192个数据帧 | 同上 |
| 每个连接中已接收并等待`Receive`的数据 | 8 MB 或256个帧 | 系统会为你重试前端的`send()`，直到处理程序跟上 |
| 整个应用中已接收并等待`Receive`的数据 | 256 MB 或8192个帧 | 同上 |
| 每个窗口的连接数 | 256 | 系统会为你重试打开操作，直到有空闲槽位 |
| 整个应用中的活动连接数 | 4096 | 同上 |
| 每个窗口的会话数 | 16 | 重新加载会取代其自身的旧会话；否则将重试打开操作 |
| 任一方向上的单个帧 | 64 MB | Go 返回`ErrStreamTooLarge`；从 JS 发送时，流会引发`error`并关闭 |
| 流名称 | 256个 UTF-8字节 | 打开操作被拒绝，流会引发`error` |
| 空闲轮询保持时间 | 20 s | 轮询返回空结果，运行时会立即重新发起轮询 |

上述情况都不会静默丢弃数据。标有<em>系统会为你重试</em>的两行表示普通的背压——运行时会保留该帧，并在短暂退避后重试，因此你的代码看到的是速度较慢的流，而不是错误。会引发`error`的各行表示编程错误，而不是负载问题；这些错误会显式呈现，而不会被掩盖。

目前这些是编译时常量，而不是选项。如需更改，请参阅内部机制指南。

## 不应使用流的情况

- <strong>对于请求/响应，请使用绑定。</strong>流适用于连续或非请求触发的数据；对于返回值的调用，使用绑定方法更简单。
- <strong>对于应用事件，请使用`Emit`/`On`。</strong>事件会分发给每个监听器，并由一套独立且成熟的系统处理。流是点对点的。
- <strong>不适用于`InitialHTML`窗口。</strong>这些窗口通过`origin === "null"`加载，因此根本无法访问资产服务器。
