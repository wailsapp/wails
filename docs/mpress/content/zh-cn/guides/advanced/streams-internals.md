---
title: "流——内部机制"
description: "流传输的工作原理、采用这种设计的原因、缓冲区常量的含义，以及尚未完成的部分"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

供所有修改流传输的人或代理参考。面向用户的 API 位于[流](/guides/streams/)中；本页介绍其底层机制和背后的设计考量，因为有些决策在不了解其所规避的问题时会显得武断。

## 文件

| 文件 | 作用 |
| --- | --- |
| `v3/pkg/application/stream.go` | 公共 API、`StreamConn`、`streamSink`、管理器及其注册表 |
| `v3/pkg/application/stream_session.go` | 一个窗口中的一次页面加载：出站队列、帧类型、连接表 |
| `v3/pkg/application/stream_transport.go` | 两个 HTTP 端点、二进制成帧、分块重组、运行时前置代码 |
| `v3/pkg/application/stream_server.go` | 仅限`-tags server`：真正的 WebSocket 接收端 |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | 在提供 bundle 时选择客户端传输机制 |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | 形如`WebSocket`的客户端 |
| `v3/tests/stream-performance/` | 负载测试工具（`-upload`、`-reloads`、场景遍历） |

## 整体结构

Go→JS 和 JS→Go 使用不同的机制，而这种不对称正是整个设计的核心。

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

<strong>Go→JS 使用挂起式轮询。</strong>请求会挂起，直到有内容可供传送。这里刻意没有设置轮询间隔，也不采用任何自适应机制：服务器会一直保持请求，直到出现帧，因此传送延迟已经约为0，任何客户端间隔都只会增加延迟。在响应仍在传输时到达的帧会累积起来，随下一个响应一并发送，因此往返过程本身就成为批处理窗口——负载上升时，它无需任何测量便会自行扩大。实测结果：速率为100/s 时，每个响应包含1.0帧；速率达到5000/s 时仍为1.0帧；速率达到20000/s 时为3.4帧；而且随着速率上升，p99 延迟反而<em>下降</em>。

<strong>JS→Go 使用普通的 POST。</strong>每个连接的发送操作通过 Promise 链串行执行，因为并发调用`fetch`无法保持顺序，而 Go 依赖于它观察到的顺序与发送顺序一致。在进行中的请求之后累积的帧会成批放入下一个 POST。Go 会在响应<em>之前</em>，将已接受的帧或批次前缀追加到连接的收件箱中，因此客户端无法越过 Go 尚未排队的字节继续推进。

<strong>每个窗口同时只有一个轮询请求在进行，并复用该请求传输所有连接的数据。</strong>这使顺序从结构上就能保证正确——只有一个队列和一个排空器，不存在可能超越第一条路径的第二条传送路径。它还规避了 Windows 上 HTTP/1.1每个主机最多六个连接的限制；在 Windows 上，这些请求是针对`http://wails.localhost`发起的真正 Chromium 网络请求。

## 为何采用这些具体设计

这些设计都是事件传输开发过程中留下的经验教训。移除其中任何一项，都会使一个经测量确认的问题再次出现。

**Go→JS 路径中的任何操作都不会触及主线程。**`Send`会在互斥锁保护下追加内容，然后返回。过去，当早先由 goroutine 发出的事件仍在队列中时，从主线程发出的事件会以内联方式执行其 eval——在所有三个平台上，均有4.4% 的事件发生顺序颠倒。只有一个排空器的单一队列不会出现这种情况。

<strong>无论数据大小如何，任何操作都不会触及`evaluateJavaScript`。</strong>将有效载荷拼接进 eval 源代码，会在超过特定于平台的临界点后持续占用宿主内存：以100 × 1 MB/秒运行时，macOS 上为11.6 GB，WebKitGTK 上为6.2 GB。流完全不会使用它，因此固定字节率遍历在所有帧大小下都保持平坦。

<strong>控制数据通过标头传输，绝不放入正文或查询字符串。</strong>对于自定义 URI 方案，WebKitGTK 的6.0可能会把 POST 正文作为查询参数传送（`transport_http.go`正是为此提供了回退机制），而 WebView2 对正文传送的限制约为2 MB。

<strong>轮询响应采用二进制格式，而非 JSON。</strong>帧为`[]byte`；如果将 base64 放入 JSON 封装中，每个帧都会增加33% 的开销，并且还需在 UI 线程上进行一次解析。

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind`包括 data / open / close / error。没有序列号，也没有 ack：WebSocket 不会重放，而连接断开时会丢失传输中的内容。模拟这种行为，比使用一个有界缓冲区无法始终满足的游标更简单，也更如实。

**保持请求是安全的**，因为每个 webview 请求本来就会获得自己的 goroutine。`assetserver_webview.go`中的`dispatchWorkers`固定为0，并有注释明确指出正是为了这种情况；若要启用该池，必须先为请求生命周期设定上限。

## 缓冲区常量

全部位于`stream.go`中。**它们是编译时常量，而非选项**——没有`Options.Streams`，也没有按流设置的配置。要更改它们，必须编辑该文件。

| 常量 | 值 | 限制的内容 |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 MB | 每个窗口中等待收集的缓冲字节数 |
| `streamOutQueueDepth` | 256 | 每个窗口中缓冲的帧数 |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 MB / 8192 | 整个应用中缓冲的出站数据 |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 MB / 8192 | 整个应用中等待`Receive`的入站数据 |
| `streamMaxConnections` | 256 | 一个会话中的活动连接数与已排队关闭数之和 |
| `streamMaxConnectionsGlobal` | 4096 | 整个应用中的活动连接数 |
| `streamOutCloseDepthGlobal` | 4096 | 整个应用中尚未传送的关闭通知数 |
| `streamMaxSessionsPerWindow` | 16 | 一个窗口可以保留的会话数；超过此数后，较新代必须取代较旧代 |
| `streamMaxSessions` | 1024 | 整个应用中的会话数 |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | 每个会话以及整个应用中排队等待的非关闭控制帧 |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | 每个会话中未完成的上传数，以及单次上传中的分片数 |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 MB / 4096 | 整个应用中的分片有效负载和分片元数据 |
| `streamMaxChunkIDLen` | 64 字节 | 一个由客户端提供的分片集标识符 |
| `streamMaxResponseBytes` | 1 MB | 单个轮询响应 |
| `streamHoldTimeout` | 20 秒 | 空轮询保持挂起的时长 |
| `streamSessionTTL` | 60 秒 | 达到此时长未轮询且没有活跃连接 ⇒ 会话已失效 |
| `streamSessionGrace` | 10 分钟 | 存在活跃连接但达到此时长未轮询 ⇒ 会话已失效 |
| `streamSessionSweep` | 20 秒 | 清理程序检查失效会话的频率 |
| `streamMaxFrameBytes` | 64 MB | 任一方向上的单个帧 |
| `streamMaxNameLen` | 256 字节 | 一个已注册或请求的流名称 |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 MB | 已收到但尚未被`Receive`取走的帧 |

### 如何选择这些值

<strong>`streamOutQueueDepth`有意没有采用`eventQueueCapacity`（64）。</strong>该常量是针对每次 eval 仅取出一个元素的队列测得的；在这种情况下，增加深度只会增加尾部延迟。轮询会批量取出元素，因此此处的深度必须足以容纳一次往返期间产生的帧——按每秒5000帧、往返耗时5毫秒计算，约为25帧。256还能容纳突发流量，而不会阻塞生产者。

**`streamOutQueueBytes`才是真正重要的上限**，因为256个1 MB 的帧合计为256 MB。当前端停止取走数据时，它是限制主机内存占用的最后一道防线。

这里有两条相互作用的规则，其中第二条很容易被意外破坏：

- 深度和字节上限共同限制<em>累积量</em>。
- <strong>空队列始终接受一个帧，无论该帧有多大。</strong>若无条件强制执行字节上限，大于该上限的帧将根本无法发送——等待条件永远不可能满足，因此`Send`会永久阻塞，而`TrySend`会一直报告队列已满。帧大小并不总能由调用方决定；带有`[]byte`字段的结构体在编组后有多大就是多大。

<strong>`streamMaxResponseBytes`因 Windows 而存在。</strong>WebView2 响应写入器会在内存中累积整个响应体，直到`Finish`时才将其交出，因此无上限的响应会在该平台导致无上限的内存分配。提高这个值<em>并不能</em>提升 Windows 的吞吐量——测量表明，Windows 的瓶颈取决于每字节开销，而非每响应开销：在遍历各种帧大小时，每秒响应数相差4倍，而 MB/s 稳定在约90。

<strong>入站上限会使前端等待。</strong>在桌面模式下，`deliver`报告队列已满，端点返回`429`，客户端则以有上限的退避重试同一帧或批次中未被接受的后缀。这样可以避免在处理程序追赶进度时占用 WebView 请求槽。在服务器模式下，套接字读取泵会等待，并让 TCP 施加背压。若无此上限，迟迟不调用`Receive`的处理程序可能导致主机内存无限增长。

<strong>控制帧不受数据上限约束，但有独立的生命周期上限。</strong>因背压丢失数据帧只会造成减速；丢失打开确认会使前端永远停留在`CONNECTING`，而丢失关闭帧则会使其误以为已失效的连接仍然活跃。因此，非关闭控制帧使用自己的有界队列，与关闭帧所用的队列分开，防止一批遭拒的打开请求耗尽已接受连接报告其结束状态所需的容量。每个会话还会为每个已接受的连接预留一个关闭槽位。当该容量已被占用时，新的打开请求会在注册前收到可重试的背压响应。

<strong>每会话上限也有对应的应用级全局上限。</strong>如果没有全局上限，每个获准的会话或连接都可能同时占满其全部本地配额。因此，在桌面和服务器传输之间，出站与入站数据分别共享独立的256 MiB / 8192帧预算。活跃连接共享一份包含4096个条目的预算，尚未送达的关闭通知则共享另一份大小相同的预算。达到共享配额时，行为与达到本地配额时相同：执行阻塞式`Send`或非阻塞式`TrySend`；每条取出、接收、关闭、写入失败和终止运行路径都会归还其预留量。

这两份预算有意彼此独立，并未设计成由连接将自己的配额转交给其关闭帧。每份预留量都只由一个所有者释放：连接槽位由只运行一次的`shutdown`释放；关闭帧槽位则由处置该帧的一方释放——即取出操作，或被拆除的所属会话。若所有权在双方之间转移，就必须以原子方式完成；较早的一个修订版本曾让关闭帧继承连接槽位，但只要拆除发生在尝试关闭与该尝试失败之间，就会永久泄漏一个槽位。

<strong>Go 帧会转移所有权；JavaScript 帧则会创建快照。</strong>Go 的`Send`会保留调用方的切片，直到传输层将其写出，因此调用成功后，调用方不得修改或复用该存储空间。JavaScript 的`send()`会在返回前复制可变二进制输入，从而与原生 WebSocket 的所有权语义保持一致。这种非对称规则既避免了在 Go 内部再次完整复制帧，又使面向浏览器的 API 行为符合预期。

**JavaScript 发送遵循 WebSocket 缓冲契约。**`send()`不能阻塞，因此应用排入数据的速度可能超过桌面请求通道的接收速度，就像它也可能超过原生 WebSocket 的处理速度一样。`bufferedAmount`包含该套接字保留的每一个字节，是提供给调用方的背压信号；主机端队列仍独立受上述限制约束。发生终止性故障、对端关闭或本地`close()`时，会释放保留的有效负载。本地关闭还会先取消正在`429`上等待的打开请求或数据请求，然后再发送已预留的关闭控制帧，因此准入背压或接收方背压不会使套接字卡在`CLOSING`状态。

<strong>分块重组共享一项宿主内存配额。</strong>每个会话可以组装一个最大为64 MiB 的帧，但不能为每个获准接入的会话分别分配这项配额。因此，未完成和可重试的分块集共享128 MiB 的已接纳载荷预算。分块集完成时会短暂地同时保留其各个分块和组装后的连续帧，因此，即使将该逻辑配额加倍，仍不会超过256 MiB 的实际内存上限。保留的分块还共享4096个条目的元数据配额，以免微小或空分块在尚未接近字节限制时就让映射和切片簿记数据无限增长。会超出任一配额的请求将收到可重试的背压；交付、拒绝、过期或会话关闭后，所占字节和分块条目都会归还到共享预算中。

<strong>轮询仅重试可恢复的故障。</strong>网络错误、请求超时响应（`408`）、早期数据响应（`425`）、背压（`429`）和服务器错误（`5xx`）采用指数退避，时间从250 ms 增加到最多5秒。其他`4xx`响应属于协议或所有权故障，会立即关闭页面的 Streams；`410`则是已停用会话正常终止的信号。关闭最后一个连接会中止正在进行的轮询或退避计时器；如果在此收尾期间打开连接，则会启动一个替代轮询循环。

**`streamSessionTTL`必须明显大于`streamHoldTimeout`**，否则会话自身的轮询在正常挂起期间，该会话就会被清理。

如果要针对大量小消息的工作负载进行调优，会先触及深度上限；对于大型载荷，则会先触及字节上限。典型应用无需更改这两个上限——在 macOS 上，默认值可维持634000帧/秒和2100 MB/秒。

## 连接和会话生命周期

一个<strong>会话</strong>对应一个窗口中的一次页面加载，以客户端生成的 ID（例如运行时的`clientId`）为键。会话由最先到达的请求延迟创建。当平台无法识别发出请求的窗口时（`windowID == 0`），会话 ID 的总数仍受全局限制，但系统有意不比较其代次：这些会话可能属于彼此独立、代次计数器互不相关的浏览器客户端。它们通过关闭或 TTL 过期，而不是通过相互取代来终止。

有三种机制会关闭相关对象，按发现速度从快到慢依次为：

1. <strong>窗口中较新会话的轮询会取代较旧的代次。</strong>页面重新加载后会获得新的会话 ID，并递增该窗口`sessionStorage`中存储的代次。相同的值还会镜像到`window.name`中，以便在存储被禁用时仍能跨重新加载保留，并以`performance.timeOrigin`（旧版引擎上为`Date.now()`）为基准，因此即使清除这两个存储，计数也不会从一重新开始。每个请求都携带会话 ID 和代次。轮询只会停用代次较低的页面，因此服务器调度不会让上一页面延迟到达的请求显得比替代它的新页面更新。如果策略同时阻止存储和`window.name`，排序将回退到页面时钟，因此依赖于后加载的页面获得更晚的时间原点。管理器会为每个窗口保留已停用代次的高水位标记，这样，无需保留所有历史会话 ID，也能防止已在处理中的请求重新创建旧页面。上一会话的连接会立即关闭。
2. <strong>窗口销毁</strong>会移除该窗口的所有会话，行为与`eventPayloadStore.dropWindow`类似。
3. <strong>TTL 清理</strong>会处理其他所有情况，例如渲染器崩溃或计算机进入睡眠。

只有前两种机制会停用页面代次。TTL 清理会移除空闲会话，但不会推进已停用代次的高水位标记：页面的最后一个连接关闭后便会停止轮询，但这个仍处于加载状态的页面以后必须能够再次打开流。真正已被取代的代次仍会被阻止，因为新页面的轮询会在移除旧会话前推进高水位标记。

<strong>Apple WebView 会报告已取消的请求。</strong>在 macOS 和 iOS 上，WebKit 的`stopURLSchemeTask`回调会取消匹配的请求上下文，因此属于已离开页面的轮询会立即解除阻塞。注册表以保留的原生任务标识为键，并在请求处理关闭条目时将其移除。当前桥接在 Linux 和 Windows 上仍未提供等效的提前中止回调；在这些平台上，挂起的请求会一直保留到等待期限届满。规则1确保<em>连接</em>无论如何都会迅速关闭。在其他情况下，取消会在 Linux 上表现为`EPIPE`，而在 Windows 上只有到`Finish`时才会显现。

## 传输方式选择

`Stream(name)`会查询`window._wails.streamFactory`。服务器构建会安装一个返回真实`WebSocket`的实现；WebView 构建则不设置它，因而使用轮询客户端。

工厂<strong>必须</strong>在任何模块主体运行前安装，因为生成的绑定会在模块作用域创建流。`custom.js`无法完成此操作——`loadOptionalScript`会先发出 HEAD 请求，然后追加`<script>`标签，因此执行得太晚。工厂会改为在提供运行时包时添加到其开头（`stream_prelude_server.go`）；按照构造方式，此过程是同步的：ES 模块依赖项会先于导入它们的模块求值。

如果添加第三种传输方式，也要将其放入前导代码中。不要再改回`custom.js`。

## 尚未完成的工作

|  | 状态 |
| --- | --- |
| 来自平台层的请求取消 | **Apple 已完成；Linux/Windows 尚待完成**——见上文 |
| 将缓冲区常量设为选项 | 尚未完成；仅能在编译时设置 |
| 类型化流 | 有意未实现——根据设计决定，帧为`[]byte` |
| 流水线处理（同时进行第二个轮询） | 尚未完成；需要在 JS 中进行有序重组 |
| 连接间的公平性 | 尚未完成——一个窗口中的连接共享同一队列，因此大量发送数据的连接会拖慢相邻连接 |
| JS→Go 帧合并 | **已完成**——在处理中的请求之后积累的帧会按有界批次发送；负载较轻的连接仍会在每次 POST 中发送一帧 |
| Windows 吞吐量 | 约100 MB/s，受`WebResourceRequested`封送处理限制。共享缓冲区（`PostSharedBufferToScript`）是候选修复方案；`internal/webview2/pkg/webview2/`下已有绑定，但尚未接入`pkg/edge` |
| `wails3 dev` / Vite | **可用**——已使用生成的`vanilla-js`项目验证：Vite 开发服务器在`/`处进行代理，而`/wails/stream/*`会在代理之前由资产服务器中间件匹配，因此流不受影响 |
| 多窗口 | 尚未经过负载测试，但会话按构造方式限定在窗口作用域内 |

## 开发模式下的前端包

生成的项目从<strong>npm</strong>导入`@wailsio/runtime`，而不是从资产服务器提供的捆绑`/wails/runtime.js`中导入。在`wails3 dev`下，Vite 会从`node_modules`解析它，因此，使用已发布运行时构建的应用不会看到分支中新增的客户端功能。

在流功能尚未发布期间，将测试应用指向此工作副本中的软件包：

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

该命令会先重新构建`dist/`，因此始终安装当前源代码。要撤销此操作，请在同一目录中运行`npm install @wailsio/runtime@latest`。

请注意，客户端有<strong>两个</strong>构建输出，很容易只重新构建其中一个而漏掉另一个：`task v3:runtime:build:package`生成 npm 软件包的`dist/`（应用前端导入的内容），而`task v3:runtime:build:assets`生成`bundledassets/runtime.js`（WebView 从资源服务器加载的内容）。更改`stream.ts`后，两者都需要重新构建。

## 测试

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

顺序测试至关重要：八个 goroutine 并发发送数据，所用计数器在持有队列锁时分配；排空后的顺序必须与接受顺序完全一致。如果该测试失败，就说明单一排空器不变量已被破坏。

负载测试工具：

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

在 Windows 上，该工具必须在交互式控制台会话中运行——通过普通 SSH 调用时，它会在会话0中终止，并产生零长度输出——而且必须将二进制文件暂存到 SSH 账户和控制台账户都能读取的位置，因为`C:\Users\<user>`的 ACL 仅允许其所有者访问。
