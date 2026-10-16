---
title: "WebSocket을 Streams로 마이그레이션하기"
description: "기존 WebSocket 구현을 Wails streams로 단계별 변환하는 방법과 드러나지 않게 오작동을 일으키는 차이점"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

현재 WebSocket 서버를 실행하는 앱을 [Streams](/guides/streams/)로 변환하기 위한 단계별 가이드입니다. 에이전트도 그대로 따라 할 수 있도록 작성되었습니다.

이렇게 하면 리스너를 삭제할 수 있습니다. 바인딩된 TCP 포트도, 오리진 검사도, 토큰도 없어지며 방화벽이나 엔드포인트 보안 제품이 문제 삼을 요소도 사라집니다. API가 충분히 비슷해 프런트엔드 코드 대부분은 그대로 유지할 수 있지만, <strong>드러나지 않게 오작동을 일으키는 세 가지 차이점</strong>이 있습니다. 이 차이로 한나절을 허비할 수 있으므로 먼저 설명합니다.

## 일반적인 사례: 우회 수단으로 사용하는 로컬 HTTP 서버

Wails 앱에서 WebSocket을 사용하는 일반적인 이유는 프런트엔드에 연속 데이터 피드를 푸시할 다른 방법이 없었기 때문입니다. 이에 따라 앱이 로컬 포트에서 자체 `http.Server`을(를) 구동하고 프런트엔드가 다시 여기에 연결합니다. 이런 구조라면 이번 마이그레이션으로 서버 자체를 삭제할 수 있으며, 서버 *주변에* 구축한 여러 요소도 함께 사라집니다.

**포트 검색 메커니즘이 사라집니다.** 바인딩된 `GetServerPort()`, 사용 중일 때 대체 포트를 사용하는 고정 포트, 주입된 전역 값 또는 `localStorage`에 저장한 값 등 어떤 방식으로든 프런트엔드에 연결할 포트를 알려야 했습니다. 이제 이 모든 것이 사라집니다. 스트림은 이름으로 지정하며, 이 이름은 양쪽 모두에서 컴파일 타임 상수입니다.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**CORS 구성이 사라집니다.** 웹뷰의 오리진은 플랫폼에 따라 `wails://` 또는 `http://wails.localhost`이므로 로컬 서버에는 `CheckOrigin`, `Access-Control-Allow-Origin` 헤더 또는 둘 다 필요합니다. Streams는 페이지를 로드한 애셋 서버를 통해 동작하므로 허용해야 할 교차 오리진 요청이 없습니다.

**직접 만든 인증 토큰이 모두 사라집니다.** localhost에 바인딩된 포트에는 컴퓨터의 모든 프로세스가 접근할 수 있으므로, 신중하게 구현하려면 다른 소프트웨어의 연결을 막기 위한 토큰이나 nonce를 추가해야 합니다. 이제 접근할 포트 자체가 없습니다.

**WebSocket이 아닌 엔드포인트는 애셋 서버 미들웨어로 이동합니다.** 이런 서버는 순수한 WebSocket 서버로 남는 경우가 드뭅니다. 소켓 옆에 파일 다운로드, 이미지 엔드포인트, 상태 확인 등이 쌓이는 경향이 있습니다. Streams가 이를 대체하지는 않지만, 이를 위해 별도의 서버를 둘 필요도 없습니다. 같은 핸들러를 애셋 서버에 마운트하세요:

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

그러면 프런트엔드는 호스트, 포트, CORS 없이 `/api/...`을(를) 동일 오리진의 상대 URL로 호출합니다. Streams와 이 방식을 함께 사용하면 로컬 서버에 남는 역할이 없습니다.

## 시작하기 전에 읽어 보세요

### 1. `ev.data`은(는) 문자열이 아니라 항상 `ArrayBuffer`입니다

가장 중요한 차이입니다. WebSocket은 텍스트 메시지를 문자열로 전달하지만 스트림은 모든 메시지를 바이트로 전달합니다. 다음과 같은 코드는 **컴파일되고 실행되지만 잘못 동작합니다**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

모든 호출 지점에서 수정하지 말고 경계에서 수정하세요:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

또는 트래픽이 JSON이라면(대부분 그렇습니다) `Stream` 대신 `JSONStream`을(를) 사용하여 이 문제를 완전히 피하세요:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

JSON WebSocket을 마이그레이션하는 가장 간단한 방법입니다. 생성자를 교체하고 `JSON.parse` 및 `JSON.stringify` 호출을 제거하면 나머지 핸들러는 변경하지 않아도 됩니다. JSON이 아닌 트래픽의 경우 한 번만 래핑하고 모든 핸들러를 그대로 두세요. [호환성 shim](#-shim)을 참조하세요.

### 2. 보내기는 호환되지만 받기는 호환되지 않습니다

`send()`은(는) 문자열을 받아 UTF-8로 인코딩하므로 `s.send(JSON.stringify(x))`은(는) 변경 없이 작동합니다. 수신 경로만 수정하면 됩니다. 코드의 절반은 계속 작동하기 때문에 이러한 비대칭을 놓치기 쉽습니다.

### 3. URL이 없습니다

WebSocket은 경로, 쿼리 문자열, 서브프로토콜, 인증 토큰 등의 연결 매개변수를 URL에 담습니다. 스트림에는 이름만 있습니다. URL로 전달했던 모든 값은 첫 번째 프레임으로 옮기거나 연결 전에 호출하는 바인딩된 메서드로 옮겨야 합니다.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go 측

HTTP 서버, 업그레이더, 연결 레지스트리를 삭제하세요. 각각은 핸들러로 대체됩니다.

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
| `upgrader.Upgrade` / `websocket.Accept` | *(없음 — `HandleStream` 자체가 전체 등록 과정임)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()` 또는 핸들러에서 바로 반환 |
| 브로드캐스트용 연결 레지스트리 | 직접 유지 — [브로드캐스트](#heading) 참조 |
| `http.ListenAndServe`, mux, 오리진 검사, 토큰 | **삭제** |
| ping/pong 연결 유지 | **삭제** — 연결 상태를 유지할 유휴 소켓이 없음 |
| `r.Context()` | `c.Context()` |

핸들러 고루틴의 수명은 gorilla 핸들러와 마찬가지로 연결의 수명과 정확히 일치하므로 기존 루프 구조를 변경 없이 그대로 사용할 수 있습니다.

## 프런트엔드 측

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

`readyState`, 네 가지 상태 상수, `addEventListener`, `close(code, reason)` 및 `bufferedAmount`은(는) 모두 `WebSocket`에서와 동일하게 동작합니다.

### 호환성 shim

핸들러를 전혀 수정하고 싶지 않다면 생성자를 한 번 래핑하세요. 그러면 문자열 메시지를 기대하는 기존 코드가 그대로 작동합니다:

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

그러면 `const ws = TextStream("feed")`을 `new WebSocket(url)` 대신 그대로 사용할 수 있습니다.

## 브로드캐스트

WebSocket 서버는 일반적으로 여러 대상으로 전송할 수 있도록 레지스트리를 유지합니다. 스트림에는 브로드캐스트 기능이 내장되어 있지 않습니다. 레지스트리는 그대로 유지하되 `*websocket.Conn` 대신 `*StreamConn`을 저장하세요:

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

다음 두 가지 규칙을 지키는 것이 좋습니다. 전송하기 전에 잠금을 해제하고, 단일 소비자의 느린 처리로 인해 다른 모든 클라이언트가 중단되지 않도록 여러 대상으로 전송할 때는 `TrySend`을 사용하는 편이 좋습니다.

## 다른 경우: 프런트엔드가 브로커와 직접 통신하는 경우

Wails 앱에서는 흔하지 않지만 알아 둘 가치가 있습니다. 프런트엔드가 **앱이 아닌 브로커에** WebSocket을 여는 경우(예: `nats.ws`을 사용해 NATS 서버에 연결하거나 WebSocket을 통해 MQTT에 연결하는 경우), 스트림을 그대로 대신 사용할 수 없습니다. 스트림은 프런트엔드를 타사가 아니라 <em>개발자의 Go 코드</em>에 연결하기 때문입니다.

이 마이그레이션에는 아키텍처 변경이 필요하며, 대개 바람직한 변경입니다:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

브라우저용 라이브러리보다 네이티브 라이브러리가 더 나은 Go로 브로커 클라이언트를 옮기고, 프런트엔드에 필요한 부분을 스트림을 통해 노출하세요:

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

구독 콜백 안의 `TrySend`에 유의하세요. 이 콜백은 브로커 클라이언트의 고루틴에서 실행되므로, 이를 블로킹하면 해당 연결의 모든 구독에 대한 전달이 중단됩니다.

다음과 같은 이점이 있습니다. 브로커 자격 증명이 프런트엔드에 전달되지 않고, 머신에 WebSocket 포트를 노출하지 않으며, 브라우저 클라이언트 대신 성숙한 Go 클라이언트가 재연결과 백오프를 처리합니다.

## 마이그레이션 체크리스트

- [ ] 기존의 각 WebSocket 엔드포인트에 `HandleStream` 등록
- [ ] 읽기 루프 변환: `ReadMessage`/`Read` → `c.Receive()`
- [ ] 쓰기 변환: `WriteMessage` → `c.Send()`, 여러 대상으로 전송하거나 브로커 콜백에서 쓸 때는 `TrySend` 사용
- [ ] HTTP 서버, mux 항목, 업그레이더, 오리진 검사 및 인증 토큰 **삭제**
- [ ] ping/pong 연결 유지 기능 **삭제**
- [ ] URL 매개변수를 첫 번째 프레임 또는 바인딩된 메서드로 이동
- [ ] **모든 `ev.data` 읽기를 디코딩** — `new TextDecoder().decode(ev.data)` — 또는 호환성 shim 도입
- [ ] 재연결 로직은 그대로 유지(의도적으로 재연결 기능이 내장되어 있지 않음)
- [ ] 프런트엔드가 브로커와 직접 통신했다면 브로커 클라이언트를 Go로 이동
- [ ] `InitialHTML` 창이 있는지 확인 — 해당 창에서는 스트림을 전혀 사용할 수 없음

변환 중에 grep으로 찾아볼 만한 항목:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## 예상되는 동작 차이

|  | WebSocket | 스트림 |
| --- | --- | --- |
| 메시지 유형 | 텍스트 또는 바이너리 | 바이트만 |
| `ev.data` | 문자열 또는 `Blob`/`ArrayBuffer` | 항상 `ArrayBuffer`(`binaryType = "blob"`인 경우 제외) |
| `binaryType` 기본값 | `"blob"` | `"arraybuffer"` |
| 서브프로토콜, `extensions` | 협상됨 | 지원되지 않음. 항상 `""` |
| 연결 매개변수 | URL 및 쿼리 | 첫 번째 프레임 또는 바인딩된 호출 |
| 인증 | 토큰 또는 쿠키 | 필요 없음 — 앱만 호출할 수 있음 |
| 연결 유지 | ping/pong | 필요 없음 |
| 자동 재연결 | 없음 | 없음(동일) |
| 종료 코드 | 전체 범위 | `1000` 정상, `1001` 세션 종료, `1002` 프레이밍 불일치, `1006` 오류 |
| 백프레셔 | 커널 소켓 버퍼 | 창당 8 MB / 256개 프레임, 이후 `Send` 호출은 대기 |
| 여러 연결, 하나의 엔드포인트 | 예 | 예 |

## 마이그레이션 후

흔한 실수를 찾아내는 기본 점검 항목은 다음과 같습니다:

1. 페이지를 반복해서 새로 고침하세요. 매번 핸들러가 종료되고 새 핸들러가 시작되어야 하며, 핸들러가 누적되어서는 안 됩니다.
2. 각 방향으로 512 KB보다 큰 메시지를 보내세요.
3. 1분 넘게 유휴 상태로 두세요. 다시 연결하지 않아도 트래픽이 재개되어야 합니다.
4. 개발자 도구를 열고 부하가 걸린 상태에서 중단점에서 일시 중지한 다음 재개하세요. 프로듀서는 데이터를 잃거나 제한 없이 커지는 대신 차단되었다가 복구되어야 합니다.
