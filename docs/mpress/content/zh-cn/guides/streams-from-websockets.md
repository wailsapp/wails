---
title: "将 WebSocket 迁移到 Streams"
description: "将现有 WebSocket 实现逐步转换为 Wails streams，包括那些会悄无声息地导致故障的差异"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

本指南按固定步骤说明如何将当前运行 WebSocket 服务器的应用转换为 [Streams](/guides/streams/)。编写方式适合严格照做，包括由智能体执行。

这样做的好处是可以删除监听器：不再绑定 TCP 端口，不需要源检查，不需要令牌，也不会再有任何内容引起防火墙或端点安全产品的异议。两者的 API 足够接近，因此大多数前端代码都无需改动——但有<strong>三个差异会悄无声息地导致故障</strong>。本文首先列出这些差异，因为它们很可能会耗掉你一个下午。

## 常见情况：用本地 HTTP 服务器作为变通方案

Wails 应用使用 WebSocket，通常是因为没有其他办法将连续数据流推送到前端，所以应用会在本地端口上启动自己的`http.Server`，再由前端连接回来。如果你的架构属于这种情况，此次迁移会彻底删除该服务器，同时也会移除围绕<em>它</em>构建的若干内容。

<strong>端口发现机制将被移除。</strong>必须通过某种方式告诉前端要连接哪个端口：已绑定的`GetServerPort()`、被占用时可回退的固定端口、注入的全局变量，或存放在`localStorage`中的值。这些都会消失——stream 通过名称寻址，而且两端的名称都是编译时常量。

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

<strong>CORS 配置将被移除。</strong>根据平台不同，webview 的源是`wails://`或`http://wails.localhost`，因此本地服务器需要`CheckOrigin`、`Access-Control-Allow-Origin`标头，或同时需要两者。Streams 使用加载页面的同一资产服务器，因此没有需要放行的跨源请求。

<strong>你自行设计的所有身份验证令牌都将被移除。</strong>绑定在 localhost 上的端口可被机器上的任何进程访问，因此严谨的实现会添加令牌或 nonce，阻止其他软件连接。现在已没有可供访问的端口。

<strong>非 WebSocket 端点将迁移到资产服务器中间件。</strong>这些服务器很少始终只处理 WebSocket——旁边往往会逐渐增加文件下载、图像端点和健康检查。Streams 不会取代这些端点，但你也不必为它们另设服务器。将相同的处理程序挂载到资产服务器上：

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

随后，前端会将`/api/...`作为同源相对 URL 调用——无需主机、端口或 CORS。结合 streams 与这种处理方式，本地服务器就已无事可做。

## 开始前请先阅读

### 1. `ev.data`是`ArrayBuffer`，绝不是字符串

这是最重要的一点。WebSocket 将文本消息作为字符串传送；stream 则将每条消息都作为字节传送。如下代码<strong>可以编译和运行，但行为不正确</strong>：

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

应在边界处修复，而不是在每个调用点修复：

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

或者，如果传输的是 JSON（通常确实如此），请使用`JSONStream`而不是`Stream`，从而彻底避开这个问题：

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

对于 JSON WebSocket，这是最简短的迁移方式：替换构造函数，移除`JSON.parse`和`JSON.stringify`调用，处理程序的其余代码无需改动。对于非 JSON 流量，只需封装一次，即可让所有处理程序保持不变——请参阅[兼容性适配层](#heading-2)。

### 2. 发送兼容，接收不兼容

`send()`接受字符串并将其编码为 UTF-8，因此`s.send(JSON.stringify(x))`无需改动即可工作。只有接收路径需要编辑。这种不对称很容易被忽略，因为一半的代码仍可正常工作。

### 3. 没有 URL

WebSocket 通过 URL 携带连接参数——路径、查询字符串、子协议和身份验证令牌。stream 只有名称。原先通过 URL 传递的所有内容都必须移到第一帧中，或移到连接前调用的绑定方法中。

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go 端

删除 HTTP 服务器、升级器和连接注册表。它们各自都由处理程序取代。

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
| `upgrader.Upgrade` / `websocket.Accept` | *（无需任何内容——`HandleStream`就是完整的注册过程）* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`，或直接从处理程序返回 |
| 用于广播的连接注册表 | 自行维护——请参阅[广播](#heading-3) |
| `http.ListenAndServe`、mux、源检查、令牌 | **删除** |
| ping/pong 保活 | **删除**——没有空闲套接字需要保持活动 |
| `r.Context()` | `c.Context()` |

处理程序 goroutine 的生命周期就是连接的生命周期，与 gorilla 处理程序完全相同，因此现有循环的结构可以原样沿用。

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

`readyState`、四个状态常量、`addEventListener`、`close(code, reason)`和`bufferedAmount`的行为都与在`WebSocket`上一样。

### 兼容性适配层

如果完全不想改动处理程序，只需封装一次构造函数。这样，原本要求字符串消息的现有代码就可以原样运行：

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

这样，`const ws = TextStream("feed")`就可以直接替代`new WebSocket(url)`。

## 广播

WebSocket 服务器通常会维护一个注册表，以便向多个客户端分发消息。流没有内置广播功能——保留该注册表，只需存储`*StreamConn`而不是`*websocket.Conn`：

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

有两条规则值得保留：发送前释放锁；在向多个客户端分发消息时优先使用`TrySend`，以免单个缓慢的消费者阻塞其他所有客户端。

## 另一种情况：前端直接与消息代理通信

这种情况在 Wails 应用中不太常见，但值得了解。如果前端打开 WebSocket，**连接的是消息代理而不是你的应用**——例如`nats.ws`连接 NATS 服务器或通过 WebSocket 使用 MQTT——流就不能直接替代它，因为流将前端连接到<em>你的 Go 代码</em>，而不是第三方。

这种迁移会改变架构，而且通常是有益的改变：

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

将消息代理客户端移入 Go；其原生库优于浏览器库，然后通过流公开前端所需的部分：

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

请注意订阅回调中的`TrySend`：该回调在消息代理客户端的 goroutine 上运行；阻塞它会使该连接上所有订阅的消息传递停滞。

这样做的好处是：消息代理凭据永远不会传到前端；无需向本机开放 WebSocket 端口；重新连接和退避由成熟的 Go 客户端处理，而不是由浏览器客户端处理。

## 迁移检查清单

- [ ] 为原有的每个 WebSocket 端点注册`HandleStream`
- [ ] 转换读取循环：`ReadMessage`/`Read` → `c.Receive()`
- [ ] 转换写入：`WriteMessage` → `c.Send()`；在任何多客户端分发或消息代理回调中，则转换为`TrySend`
- [ ] HTTP 服务器、mux 条目、升级器、来源检查和身份验证令牌均已<strong>删除</strong>
- [ ] ping/pong 保活机制已<strong>删除</strong>
- [ ] 将 URL 参数移入第一帧或绑定方法
- [ ] **每次读取`ev.data`后都进行解码**——`new TextDecoder().decode(ev.data)`——或者采用兼容层
- [ ] 保持重新连接逻辑不变（此处特意没有内置重新连接功能）
- [ ] 如果前端原本直接与消息代理通信，则将消息代理客户端移入 Go
- [ ] 检查是否存在`InitialHTML`窗口——它们完全无法使用流

转换期间值得用 grep 搜索的内容：

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## 预期的行为差异

|  | WebSocket | 流 |
| --- | --- | --- |
| 消息类型 | 文本或二进制 | 仅字节 |
| `ev.data` | 字符串或`Blob`/`ArrayBuffer` | 始终为`ArrayBuffer`（除非使用`binaryType = "blob"`） |
| `binaryType`默认值 | `"blob"` | `"arraybuffer"` |
| 子协议、`extensions` | 协商确定 | 不支持；始终为`""` |
| 连接参数 | URL 和查询参数 | 第一帧或绑定调用 |
| 身份验证 | 令牌或 Cookie | 无需验证——应用是唯一调用方 |
| 保活机制 | ping/pong | 不需要 |
| 自动重新连接 | 无 | 无（相同） |
| 关闭代码 | 完整范围 | `1000`表示正常，`1001`表示会话已关闭，`1002`表示帧格式不匹配，`1006`表示错误 |
| 背压 | 内核套接字缓冲区 | 每个窗口 8 MB / 256 帧，之后 `Send` 阻塞 |
| 多个连接，一个端点 | 是 | 是 |

## 迁移后

可发现常见错误的基本检查：

1. 反复重新加载页面——处理程序每次都应退出并启动一个新的处理程序，绝不能不断累积。
2. 在每个方向上各发送一条大于 512 KB 的消息。
3. 让其空闲超过一分钟；流量应能恢复，而无需重新连接。
4. 打开开发者工具，在负载下命中断点时暂停，然后继续执行——生产者应阻塞并恢复，而不是丢失数据或无限增长。
