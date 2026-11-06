---
title: "串流——內部原理"
description: "串流傳輸的運作方式、採用此設計的原因、緩衝區常數的含義，以及尚未完成的部分"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

供任何要修改串流傳輸的人員或代理程式參考。面向使用者的 API 位於[串流](/guides/streams/)；本頁說明其底層機制與設計考量，因為在了解這些決策所避免的問題之前，其中幾項看起來會顯得武斷。

## 檔案

| 檔案 | 用途 |
| --- | --- |
| `v3/pkg/application/stream.go` | 公開 API、`StreamConn`、`streamSink`、管理器及其登錄表 |
| `v3/pkg/application/stream_session.go` | 單一視窗中的一次頁面載入：傳出佇列、訊框種類、連線表 |
| `v3/pkg/application/stream_transport.go` | 兩個 HTTP 端點、二進位訊框封裝、區塊重組、執行階段前置碼 |
| `v3/pkg/application/stream_server.go` | 僅限`-tags server`：真正的 WebSocket 接收端 |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | 在提供套件組合時選擇用戶端傳輸方式 |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | 採用`WebSocket`介面形式的用戶端 |
| `v3/tests/stream-performance/` | 負載測試工具（`-upload`、`-reloads`、情境掃描） |

## 整體架構

Go→JS 與 JS→Go 使用不同的機制，而這種不對稱正是整個設計的核心。

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

<strong>Go→JS 採用掛起式輪詢。</strong>請求會暫停等待，直到有內容可傳送。這裡刻意不設輪詢間隔，也不採用任何自適應機制：伺服器會一直保留請求，直到訊框出現，因此傳送延遲已經約為0，任何用戶端間隔只會增加延遲。回應仍在傳輸時抵達的訊框會累積，並隨下一個回應一起傳送，因此往返時間本身就成為批次處理視窗——負載上升時，它會自行擴大，不需進行任何量測。實測結果：在100/s 時，每個回應包含1.0個訊框；在5000/s 時仍為1.0個；在20000/s 時為3.4個，而且速率越高，p99 延遲反而<em>降低</em>。

<strong>JS→Go 使用一般的 POST。</strong>每條連線都透過 Promise 鏈將傳送作業序列化，因為並行呼叫`fetch`無法維持順序，而 Go 依賴其觀察到的順序與傳送順序一致。進行中請求之後累積的訊框會批次併入下一個 POST。Go 會在回應<em>之前</em>，先將已接受的訊框或批次前綴附加至連線的收件匣，因此用戶端無法越過 Go 尚未排入佇列的位元組繼續前進。

<strong>每個視窗同時只有一個進行中的輪詢，並多工處理所有連線。</strong>這使順序在架構上必然正確——只有一個佇列、一個取出器，也沒有可能超越第一條路徑的第二條傳送路徑。這也避開了 Windows 上 HTTP/1.1每台主機六條連線的限制；在該平台上，這些是對`http://wails.localhost`發出的真正 Chromium 網路請求。

## 為何採用這些特定決策

每項決策都是事件傳輸開發留下的教訓。移除任何一項，都會使實測過的錯誤再次出現。

**Go→JS 路徑中的任何作業都不會觸及主執行緒。**`Send`會在互斥鎖保護下附加資料，然後返回。過去，若事件從主執行緒發出，而較早由 goroutine 發出的事件仍在佇列中，前者會立即執行其 eval——在三個平台上都有4.4% 的事件次序顛倒。單一佇列搭配單一取出器不可能發生此問題。

<strong>無論資料大小為何，任何作業都不會觸及`evaluateJavaScript`。</strong>將承載資料拼接至 eval 原始碼，在超過平台特定的臨界點後會持續占用主機記憶體：以100 × 1 MB/sec 傳輸時，macOS 為11.6 GB，WebKitGTK 為6.2 GB。串流完全不會經過這條路徑，因此固定位元組率掃描在各種訊框大小下都保持平坦。

<strong>控制資料一律透過標頭傳送，絕不放在主體或查詢字串中。</strong>對於自訂 URI 通訊協定，WebKitGTK 的6.0可能會將 POST 主體以查詢參數形式傳送（`transport_http.go`正是為此提供備援），而 WebView2 對主體傳送的限制約為2 MB。

<strong>輪詢回應採用二進位格式，而非 JSON。</strong>訊框是`[]byte`；若在 JSON 封套內使用 base64，每個訊框都會增加33% 的成本，還需在 UI 執行緒上進行解析。

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind`為 data / open / close / error。沒有序號，也沒有 ack：WebSocket 不會重播，而連線中斷時，傳輸中的內容就會遺失。模擬這項行為，比使用有限緩衝區無法始終滿足的游標更簡單，也更忠實。

**掛起請求是安全的**，因為每個 webview 請求原本就會獲得自己的 goroutine。`assetserver_webview.go`中的`dispatchWorkers`固定為0，且註解明確指出正是為了此情況；若要啟用該集區，必須先限制請求存續時間。

## 緩衝區常數

全部位於`stream.go`中。**它們是編譯期常數，不是選項**——沒有`Options.Streams`，也沒有個別串流設定。若要變更，必須編輯該檔案。

| 常數 | 值 | 限制的項目 |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 MB | 每個視窗中已緩衝、等待收集的位元組數 |
| `streamOutQueueDepth` | 256 | 每個視窗緩衝的訊框數 |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 MB / 8192 | 整個應用程式中緩衝的傳出資料 |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 MB / 8192 | 整個應用程式中等待`Receive`的傳入資料 |
| `streamMaxConnections` | 256 | 單一工作階段中的使用中連線加上已排入佇列的關閉作業 |
| `streamMaxConnectionsGlobal` | 4096 | 整個應用程式中的使用中連線 |
| `streamOutCloseDepthGlobal` | 4096 | 整個應用程式中尚未傳送的關閉通知 |
| `streamMaxSessionsPerWindow` | 16 | 單一視窗可保留的工作階段數；超過後，較新的世代必須取代較舊的世代 |
| `streamMaxSessions` | 1024 | 整個應用程式中的工作階段 |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | 每個工作階段及整個應用程式中排隊等候的非關閉控制訊框 |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | 每個工作階段中未完成的上傳，以及單次上傳中的分段數 |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 MB / 4096 | 整個應用程式中的分段承載資料與分段中繼資料 |
| `streamMaxChunkIDLen` | 64 位元組 | 一個由用戶端提供的分段集合識別碼 |
| `streamMaxResponseBytes` | 1 MB | 單次輪詢回應 |
| `streamHoldTimeout` | 20 秒 | 空輪詢暫停等候的時間 |
| `streamSessionTTL` | 60 秒 | 達到這段時間未輪詢且沒有尚未關閉的連線 ⇒ 工作階段已失效 |
| `streamSessionGrace` | 10 分鐘 | 仍有尚未關閉的連線但達到這段時間未輪詢 ⇒ 工作階段已失效 |
| `streamSessionSweep` | 20 秒 | 清理程式搜尋失效工作階段的頻率 |
| `streamMaxFrameBytes` | 64 MB | 任一方向上的單一訊框 |
| `streamMaxNameLen` | 256 位元組 | 單一已註冊或要求的串流名稱 |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 MB | 已收到但尚未由`Receive`取走的影格 |

### 如何選擇這些值

<strong>`streamOutQueueDepth`刻意不設為`eventQueueCapacity`（64）。</strong>該常數是針對一次 eval 僅取出一個項目的佇列測得；在這種情況下，加深佇列只會增加尾端延遲，毫無助益。輪詢會批次取出項目，因此此處的深度必須容納一個往返期間產生的資料量：若速率為每秒5000個影格，往返時間為5毫秒，便約為25個影格。256可容納突發流量，而不會讓產生端停滯。

**`streamOutQueueBytes`才是真正重要的上限**，因為256個1 MB 的影格總計為256 MB。當前端停止取走資料時，這是主機記憶體的最後一道防線。

此處有兩項規則會互相影響，而第二項很容易意外遭到破壞：

- 深度與位元組上限會限制<em>累積量</em>。
- <strong>空佇列一律接受一個影格，不論其大小。</strong>若無條件強制執行位元組上限，大於上限的影格將完全無法傳送：等候條件永遠不可能成立，因此`Send`會永久阻塞，而`TrySend`也會永久回報已滿。影格大小不一定能由呼叫端決定；含有`[]byte`欄位的結構會封送成其實際產生的大小。

<strong>`streamMaxResponseBytes`之所以存在，是因為 Windows。</strong>WebView2 回應寫入器會在記憶體中累積完整的本文，直到`Finish`才交出，因此不設上限的回應在該處會造成不設上限的記憶體配置。提高此值<em>並不會</em>提升 Windows 的輸送量；測量結果顯示，Windows 的瓶頸取決於每個位元組，而非每個回應：掃過不同影格大小時，每秒回應數相差4倍，但 MB/s 維持在約90。

<strong>輸入端上限正是讓前端等候的機制。</strong>在桌面模式中，`deliver`會回報已滿，端點則回應`429`；用戶端會以有上限的退避時間重試相同影格，或重試批次中未獲接受的後綴。如此可避免在處理常式趕上進度時占用 webview 的要求槽位。在伺服器模式中，通訊端讀取幫浦會等候，並讓 TCP 施加反向壓力。如果沒有此上限，遲遲未呼叫`Receive`的處理常式可能會讓主機記憶體無限制成長。

<strong>控制影格會略過資料上限，但有獨立的生命週期上限。</strong>在反向壓力下遺失資料影格只會造成速度變慢；遺失開啟確認會讓前端永遠停留在`CONNECTING`，而遺失關閉通知則會讓前端誤以為已失效的連線仍然存在。因此，非關閉控制影格有自己的有界佇列，與關閉影格使用的佇列分開，確保遭拒開啟要求的突發流量不會耗盡已接受連線回報其已結束所需的容量。每個工作階段也會為每個已接受的連線保留一個關閉槽位。該容量已占用時，新開啟要求會在註冊前收到可重試的反向壓力。

<strong>每個工作階段的上限也有對應的整個應用程式上限。</strong>若沒有這些上限，每個獲准的工作階段或連線都可能同時保留完整的個別配額。因此，輸出與輸入資料分別共用256 MiB／8192個訊框的預算，涵蓋桌面與伺服器傳輸。尚未關閉的連線有4096個項目的預算，尚未送達的關閉通知則有另一份相同大小的預算。達到共用配額時，會採取與達到個別配額時相同的阻塞`Send`或非阻塞`TrySend`行為；每條取出、接收、關閉、寫入失敗及關閉程序路徑都會歸還其保留量。

這兩份預算刻意彼此分開，而不是由連線將一份配額交給自己的關閉影格。每份保留量都只由一個擁有者釋放：連線的槽位由僅執行一次的`shutdown`釋放；關閉影格的槽位則由處置該影格的作業釋放，也就是取出作業或拆除其工作階段。若擁有權在雙方之間移轉，就必須以不可分割的方式完成；較早的修訂版本曾讓關閉影格繼承連線的槽位，但只要拆除發生在嘗試關閉之後、該嘗試失敗之前，就會永久洩漏一個槽位。

<strong>Go 影格會移轉擁有權；JavaScript 影格則會建立快照。</strong>Go 的`Send`會保留呼叫端的切片，直到傳輸層將其寫出，因此呼叫成功後，呼叫端不得修改或重複使用該儲存空間。JavaScript 的`send()`會在回傳前複製可變的二進位輸入，以符合原生 WebSocket 的擁有權語意。這項不對稱規則可避免在 Go 內部再次複製完整影格，同時讓面向瀏覽器的 API 行為符合預期。

**JavaScript 傳送遵循 WebSocket 緩衝合約。**`send()`無法阻塞，因此應用程式將資料排入佇列的速度可能超過桌面要求通道接受資料的速度，就像它也可能超過原生 WebSocket 的處理速度一樣。`bufferedAmount`包含該通訊端保留的每個位元組，是供呼叫端使用的反向壓力訊號；主機端佇列仍獨立受上述限制約束。終止性失敗、對等端關閉或本機`close()`都會釋放保留的承載資料。本機關閉也會先取消正在`429`上等候的開啟或資料要求，再送出保留的關閉控制影格，因此准入或接收端的反向壓力不會讓通訊端卡在`CLOSING`。

<strong>區塊重組共用一個主機記憶體配額。</strong>每個工作階段都可以組合單一大小上限為64 MiB 的訊框，但不能讓每個獲准的工作階段各自占用一份此配額。因此，未完成及可重試的區塊集合共用128 MiB 的獲准承載資料預算。集合完成時會短暫同時保留其各個部分及組合後的連續訊框，因此即使將該邏輯配額加倍，仍不會超過256 MiB 的實際記憶體上限。保留的部分也共用4096個項目的中繼資料配額，因此極小或空白區塊無法在尚未接近位元組限制時，便讓映射和切片簿記無限制增長。任何會超出其中任一配額的要求都會收到可重試的背壓；傳遞、拒絕、到期或關閉工作階段時，位元組與部分項目都會歸還至共用預算。

<strong>輪詢只會重試可復原的失敗。</strong>網路錯誤、要求逾時回應（`408`）、早期資料回應（`425`）、背壓（`429`）及伺服器錯誤（`5xx`）會使用指數退避，從250毫秒逐步增加至5秒。其他`4xx`回應代表通訊協定或擁有權失敗，會立即關閉頁面的 Streams；`410`則是已汰除工作階段的正常終止訊號。關閉最後一個連線會中止進行中的輪詢或退避計時器；若在這段結束過程中開啟連線，則會啟動一個替代輪詢迴圈。

**`streamSessionTTL`必須明顯大於`streamHoldTimeout`**，否則工作階段會在其自身輪詢正合理地停駐等待時遭到回收。

如果要針對大量小型訊息的工作負載進行調校，會先觸及深度上限；對大型承載資料而言，則會先觸及位元組上限。一般應用程式不需要變更兩者中的任何一項——在 macOS 上，預設值可維持634000訊框/秒及2100 MB/秒。

## 連線與工作階段生命週期

一個<strong>工作階段</strong>代表一個視窗中的一次頁面載入，並以用戶端產生的 ID（例如執行階段的`clientId`）作為索引鍵。工作階段由最先到達的要求延遲建立。當平台無法識別發出要求的視窗（`windowID == 0`）時，工作階段 ID 仍受全域上限約束，但系統刻意不比較其世代：它們可能屬於不同的獨立瀏覽器用戶端，且各自使用互不相關的世代計數器。這些工作階段會透過關閉或 TTL 到期，而不會彼此取代。

有三種機制會關閉相關項目，以下依察覺速度由快到慢排列：

1. <strong>同一視窗中較新的工作階段輪詢會取代較舊的世代。</strong>重新載入會為頁面提供新的工作階段 ID，並遞增儲存在該視窗`sessionStorage`中的世代值。同一個值也會鏡像至`window.name`；停用儲存空間時，該處的值在重新載入後仍會保留。此外，此值會以`performance.timeOrigin`（舊版引擎則為`Date.now()`）為基準，因此清除這兩個儲存區也不會讓計數從一重新開始。每個要求都會攜帶工作階段 ID 與世代。輪詢只會汰除頁面中較舊的世代，因此伺服器排程不會讓上一個頁面的延遲要求看起來比其替代頁面更新。如果原則同時封鎖儲存空間與`window.name`，排序便會退回使用頁面時鐘，因此取決於後續頁面是否取得較晚的時間原點。管理器會保留每個視窗的已汰除世代高水位標記，因此無須保留每個歷史工作階段 ID，也能防止已在傳輸中的要求重新建立舊頁面。上一個工作階段的連線會立即關閉。
2. <strong>銷毀視窗</strong>會移除該視窗的所有工作階段，與`eventPayloadStore.dropWindow`相同。
3. <strong>TTL 回收</strong>會處理其他所有情況——例如算繪器當機或電腦進入睡眠。

只有前兩種機制會汰除頁面世代。TTL 清理會移除閒置的工作階段，但不會推進已汰除世代高水位標記：頁面會在其最後一個連線關閉時停止輪詢，但同一個仍保持已載入狀態的頁面之後必須仍能開啟另一個串流。真正已被取代的世代仍會遭到封鎖，因為較新頁面的輪詢會先推進高水位標記，之後才移除舊工作階段。

<strong>Apple WebView 會回報已取消的要求。</strong>在 macOS 和 iOS 上，WebKit 的`stopURLSchemeTask`回呼會取消相符的要求上下文，因此屬於已導覽離開頁面的輪詢會立即解除阻塞。登錄表以保留的原生工作身分作為索引鍵，並在要求處理將其關閉時移除項目。目前的橋接在 Linux 和 Windows 上仍未公開對等的提前中止回呼；在這些平台上，停駐的要求會一直保留到等待期限屆滿。規則1表示無論如何，<em>連線</em>都會迅速關閉。在其他情況下，取消會在 Linux 上顯示為`EPIPE`，而在 Windows 上只會於`Finish`時顯現。

## 傳輸方式選擇

`Stream(name)`會查詢`window._wails.streamFactory`。伺服器組建會安裝一個能傳回實際`WebSocket`的實作；WebView 組建則不設定該項目，並取得輪詢用戶端。

工廠<strong>必須</strong>在任何模組主體執行前完成安裝，因為產生的繫結會在模組範圍建立串流。`custom.js`無法做到這一點——`loadOptionalScript`會發出 HEAD 要求，接著附加一個`<script>`標籤，因此執行時間遠遠太晚。工廠會改為在執行階段套件組合提供服務時前置加入其中（`stream_prelude_server.go`）；依其設計，這會同步進行：ES 模組相依項目會先於其匯入者求值。

如果新增第三種傳輸方式，也請將它放入前置碼。不要再改回`custom.js`。

## 尚未完成的項目

|  | 狀態 |
| --- | --- |
| 來自平台層的要求取消 | **Apple 已完成；Linux/Windows 待處理**——請參閱上文 |
| 將緩衝區常數設為選項 | 尚未完成；僅能在編譯時設定 |
| 具型別的串流 | 刻意不實作——依設計決策，訊框為`[]byte` |
| 管線化（同時有第二個輪詢進行中） | 尚未完成；需要在 JS 中進行依序重組 |
| 各連線間的公平性 | 尚未完成——同一視窗中的連線共用一個佇列，因此大量傳輸的連線會拖慢相鄰連線 |
| JS→Go 訊框合併 | **已完成**——在進行中要求後方累積的訊框會分成有界批次傳送；負載較低的連線仍會在每個 POST 中傳送一個訊框 |
| Windows 輸送量 | 約100 MB/秒，受限於`WebResourceRequested`封送處理。候選修正方案是共用緩衝區（`PostSharedBufferToScript`）；`internal/webview2/pkg/webview2/`下已有繫結，但尚未接入`pkg/edge` |
| `wails3 dev` / Vite | **可運作**——已使用產生的`vanilla-js`專案驗證：Vite 開發伺服器會在`/`進行 Proxy，而資產伺服器中介軟體會在 Proxy 前比對`/wails/stream/*`，因此串流不受影響 |
| 多視窗 | 尚未在負載下測試，但工作階段依設計以視窗為範圍 |

## 開發模式下的前端套件

產生的專案會從<strong>npm</strong>匯入`@wailsio/runtime`，而不是從資產伺服器所提供的套件組合`/wails/runtime.js`匯入。在`wails3 dev`下，Vite 會從`node_modules`解析該項目，因此針對已發布執行階段組建的應用程式，不會看到分支中新增的用戶端功能。

在串流功能尚未發布期間，請讓測試應用程式使用此工作副本中的套件：

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

這會先重新建置`dist/`，因此安裝的永遠是目前的原始碼。若要還原，請在同一目錄中執行`npm install @wailsio/runtime@latest`。

請注意，客戶端有<strong>兩個</strong>建置輸出，很容易只重新建置其中一個而漏掉另一個：`task v3:runtime:build:package`會產生 npm 套件的`dist/`（應用程式前端匯入的內容），而`task v3:runtime:build:assets`則會產生`bundledassets/runtime.js`（webview 從資產伺服器載入的內容）。變更`stream.ts`後，兩者都必須重新建置。

## 測試

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

順序測試至關重要：八個 goroutine 會並行傳送資料，並使用在佇列鎖定期間分配的計數器；排空後的順序必須與接受順序完全一致。只要此測試失敗，就表示單一排空者不變條件已遭破壞。

負載測試工具：

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

在 Windows 上，必須於互動式主控台工作階段中執行；單純透過 SSH 呼叫會在工作階段0中終止，且輸出長度為零。此外，由於`C:\Users\<user>`的 ACL 僅允許其擁有者存取，因此必須將二進位檔暫存到 SSH 帳戶與主控台帳戶都能讀取的位置。
