---
title: "스트림"
description: "수신 대기 소켓 없이 WebSocket 프로그래밍 모델을 사용하는 Go와 JavaScript 간 양방향 바이트 스트림"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

스트림은 Go와 프런트엔드 사이에 이름이 지정된 순서 보장 양방향 바이트 채널을 제공합니다. 프로그래밍 모델은 WebSocket과 같지만 **TCP 포트를 바인딩하지 않습니다**.

사용자 지정 URL 스킴으로는 WebSocket 통신을 할 수 없으므로, webview 내부에서 WebSocket을 사용하려면 실제 HTTP 서버를 실행하고 포트에서 수신 대기해야 합니다. 데스크톱 앱에서는 머신의 다른 모든 프로세스가 접근할 수 있는 로컬 포트가 열린다는 뜻입니다. 이 포트를 안전하게 사용하려면 오리진 검사와 토큰이 필요하며, 사용자가 실행하는 모든 방화벽과 엔드포인트 보안 제품에도 노출됩니다. 스트림은 이 모든 문제를 피합니다. 앱에서 이미 제공하고 오리진에 바인딩되어 있는 애셋 서버를 통해 통신합니다.

기존 WebSocket 구현을 마이그레이션하시나요? [WebSocket을 스트림으로 마이그레이션하기](/guides/streams-from-websockets/)를 따르세요. 이 문서는 절차대로 적용할 수 있도록 작성되었으며, 눈에 띄는 오류 없이 동작을 중단시키는 세 가지 차이점부터 설명합니다.

## 빠른 시작

Go에서 스트림을 선언하세요. 핸들러는 연결마다 한 번씩 자체 goroutine에서 실행됩니다.

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

프런트엔드에서 이름으로 연결하세요. 이 객체는 `WebSocket` 인터페이스를 구현합니다.

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)`은 `new WebSocket(url)`와 똑같이 `readyState === CONNECTING`인 상태의 객체를 **동기적으로** 반환하므로 모듈 범위에서 생성할 수 있습니다:

```js
export const Telemetry = Stream("telemetry");
```

## 프레임은 바이트입니다

모든 프레임은 Go에서는 `[]byte`이고 JavaScript에서는 `ArrayBuffer`입니다. 정해진 스키마나 강제되는 인코딩은 없습니다. 원하는 대로 JSON, protobuf, CBOR 또는 원시 바이트를 사용하세요.

프레임은 <strong>바이트 스트림이 아니라 메시지</strong>입니다. 프레임 전체가 도착하거나 전혀 도착하지 않으며 길이도 함께 전달됩니다. 어느 쪽도 크기를 미리 알 필요가 없으므로 `[]byte` 필드가 있는 구조체는 마샬링된 결과 그대로 하나의 프레임으로 전송됩니다.

## 객체 보내기

프레임은 바이트이지만, 대부분의 경우 바이트 단위로 생각할 필요는 없습니다. 양쪽에는 서로 연동되는 편리한 JSON 기능이 있습니다.

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

`JSONStream`은 경계에서 인코딩이 수행된 `Stream`과 동일한 객체입니다. 별도의 프로토콜은 없으며 Go 핸들러에서는 그 차이를 구분할 수 없습니다. 유효한 JSON이 아닌 프레임은 연결을 종료하는 대신 `error` 이벤트를 발생시키고 삭제됩니다.

protobuf, CBOR, 바이너리 형식처럼 바이트가 필요하거나 직접 인코딩하려는 경우에는 일반 `Stream`을 사용하세요.

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

**핸들러 goroutine의 수명은 연결의 수명과 같습니다.** 핸들러가 반환되면 연결이 닫히므로 연결을 열어 두려는 동안 `Receive` 또는 `c.Context()`에서 블로킹하세요. 이는 `gorilla`/`coder` WebSocket 핸들러와 같은 형태입니다.

오류는 `ErrStreamClosed`(피어 연결이 끊어짐)와 `ErrStreamFull`(`TrySend`에서만 발생)입니다.

## JavaScript API

`Stream(name)`은 `WebSocket`의 유용한 하위 집합을 구현하는 객체를 반환합니다.

| 지원 항목 | 참고 |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`, `onmessage`, `onclose`, `onerror` | `addEventListener`도 포함 |
| `send(data)` | 문자열, `ArrayBuffer`, 형식화 배열 또는 `Blob` — 아래의 소유권 설명을 참조하세요 |
| `JSONStream(name)` | 동일한 객체이며, 객체를 입력하고 출력합니다 |
| `close(code, reason)` |  |
| `binaryType` | <strong>기본값은 `"blob"`가 아니라 `"arraybuffer"`</strong>입니다 |
| `bufferedAmount` | `send`에 의해 대기열에 추가되었지만 아직 Go에 도달하지 않은 바이트 |
| `protocol`, `extensions` | 항상 `""`이며 협상되지 않습니다 |

`binaryType` 기본값은 표준과 의도적으로 다르게 정한 유일한 부분입니다. 프레임은 항상 바이너리이며 `Blob`을 사용하면 메시지를 읽을 때마다 비동기 단계를 한 번 더 거쳐야 합니다. 표준 동작을 원한다면 `"blob"`으로 설정하세요.

하나 또는 여러 창에서 같은 스트림 이름으로 여러 연결을 만들 수 있습니다. 각 연결에는 자체 `StreamConn`과 자체 핸들러 goroutine이 할당됩니다.

**버퍼 소유권은 전송 방향에 따라 다릅니다.** JavaScript의 `send()`는 네이티브 WebSocket 동작과 동일하게 변경 가능한 바이너리 입력의 스냅샷을 동기적으로 생성하므로 호출자는 `send()`이 반환되는 즉시 이를 재사용할 수 있습니다. Go의 `Send`는 슬라이스를 복사하지 않고 소유권을 전송 계층으로 이전합니다. 호출이 성공한 뒤에는 해당 슬라이스를 변경하거나 재사용하지 마세요. 데이터 생산자가 저장 공간을 재사용해야 한다면 새 슬라이스를 넘기세요.

## 수명 주기

스트림은 소켓처럼 동작하며, 소켓을 닫는 이벤트는 스트림도 닫습니다.

| 이벤트 | 발생하는 동작 |
| --- | --- |
| 페이지 새로고침 또는 탐색 | 연결이 닫히고 핸들러의 `Receive`이 오류를 반환하며 새 페이지가 새로 연결됩니다 |
| `window.close()` / 창 제거 | 해당 창의 모든 연결이 닫힙니다 |
| JS의 `s.close()` | 핸들러의 `Receive`이 `ErrStreamClosed`을 반환합니다 |
| 핸들러 반환 | 프런트엔드에서 `onclose`이 발생합니다 |
| 앱 종료 | 모든 연결의 컨텍스트가 취소됩니다 |

`WebSocket`와 마찬가지로 **자동 재연결 기능은 없습니다**. 앱에 이 기능이 필요하다면 기존 WebSocket 재연결 로직을 변경 없이 사용할 수 있습니다. `onclose`에서 스트림을 다시 생성하세요.

## 백프레셔

프런트엔드가 처리 속도를 따라오지 못하면 `Send`은 송신 버퍼가 가득 찼을 때 소켓 쓰기가 블로킹되는 것처럼 블로킹됩니다. 대기하는 대신 해당 메시지를 버리려면 `TrySend`을 사용하세요:

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

개발자 도구 중단점, 숨겨진 창 또는 App Nap으로 인해 프런트엔드가 일시 중지되면 데이터 수집이 멈추고, 버퍼가 한도에 도달했을 때 데이터 생산자가 블로킹됩니다. 이는 의도된 동작입니다. 읽히지 않는 스트림이 무제한으로 커지게 두지 않고 메모리 사용량을 제한합니다.

## 서버 모드

`-tags server`로 빌드하면 전송 방식이 `/wails/stream/ws`의 <strong>실제 WebSocket</strong>으로 바뀝니다. 서버 모드에는 이미 업그레이드에 사용할 리스너가 있기 때문입니다. Go 핸들러와 프런트엔드 코드는 동일하며 앱에서 변경할 사항은 없습니다. 런타임은 모듈 코드가 실행되기 전에 전송 방식을 자동으로 선택합니다. WebSocket 연결은 기본적으로 동일 출처만 허용합니다. 프런트엔드를 신뢰할 수 있는 다른 출처에서 의도적으로 호스팅하는 서버는 `ServerOptions.WebSocketOriginPatterns`로 해당 호스트를 추가할 수 있습니다.

## 성능

스로틀링 없이 `v3/tests/stream-performance`로 측정한 결과, 약 41백만 프레임에서 0개의 손실과 0개의 순서 변경이 발생했습니다.

|  | Go→JS 최대 처리량 | JS→Go 최대 처리량 |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

최대 처리량보다 성능 양상이 더 중요합니다.

- **작은 프레임에서는 Go→JS가 훨씬 빠릅니다**. macOS에서 초당 634000프레임인 반면, 반대 방향은 초당 약 6200프레임입니다. 하나의 응답은 최대 256개 프레임을 합칩니다. JS→Go도 진행 중인 요청 뒤에 누적되는 프레임을 일괄 처리하지만, 각 연결은 여전히 자체 POST 요청 체인을 직렬화합니다. 작은 메시지를 많이 보낸다면 Go→JS를 사용하거나, 보내기 전에 애플리케이션 수준에서 일괄 처리하세요.
- **Windows에서 업로드에는 512 KB가 가장 효율적입니다.** 이보다 큰 프레임은 여러 요청으로 분할되며, 4 MB 프레임은 512 KB 프레임보다 *느리게* 측정됩니다.
- **지연 시간은 짧고 계속 짧게 유지됩니다**. macOS에서 p99는 약 1~2 ms이며, 부하가 증가해도 저하되지 않습니다. 초당 20000프레임에서 측정한 p99는 초당 100프레임에서보다 *낮았습니다*.

플랫폼별 전체 표와 측정 방법은 이 기능과 함께 제공되는 측정 기록에 있습니다.

## 제한

|  | 제한 | 제한에 도달했을 때의 동작 |
| --- | --- | --- |
| 수집을 기다리며 창별로 버퍼링됨 | 8 MB 또는 256프레임 중 먼저 도달하는 값 | `Send`는 차단되고, `TrySend`는 `ErrStreamFull`을 반환함 |
| 수집/쓰기를 기다리며 애플리케이션 전체에 버퍼링됨 | 256 MB 또는 데이터 프레임 8192개 | 동일 |
| `Receive`을 기다리며 연결별로 수신됨 | 8 MB 또는 256프레임 | 핸들러가 처리 속도를 따라잡을 때까지 프런트엔드의 `send()`을 자동으로 재시도함 |
| `Receive`을 기다리며 애플리케이션 전체에서 수신됨 | 256 MB 또는 8192프레임 | 동일 |
| 창별 연결 수 | 256 | 슬롯이 빌 때까지 열기를 자동으로 재시도함 |
| 애플리케이션 전체의 활성 연결 수 | 4096 | 동일 |
| 창별 세션 수 | 16 | 새로고침하면 해당 창의 이전 세션을 대체하며, 그 외에는 열기를 재시도함 |
| 어느 방향이든 단일 프레임 | 64 MB | Go에서는 `ErrStreamTooLarge`을 반환하고, JS에서는 스트림이 `error`을 발생시킨 후 닫힘 |
| 스트림 이름 | 256 UTF-8 바이트 | 열기가 거부되고 스트림이 `error`을 발생시킴 |
| 유휴 폴 유지 시간 | 20 s | 폴이 빈 상태로 반환되고 런타임이 즉시 다시 요청함 |

위의 어떤 경우에도 데이터가 조용히 손실되지는 않습니다. <em>자동으로 재시도</em>한다고 명시된 두 행은 일반적인 백프레셔입니다. 런타임이 프레임을 보관하고 짧은 백오프를 두고 재시도하므로, 코드에는 오류 대신 느려진 스트림으로 나타납니다. `error`을 발생시킨다고 명시된 행은 부하 문제가 아니라 프로그래밍 오류이며, 감춰지지 않고 드러납니다.

현재 이 값들은 옵션이 아니라 컴파일 시간 상수입니다. 변경해야 한다면 내부 구조 가이드를 참조하세요.

## 스트림을 사용하지 않아야 하는 경우

- **요청/응답에는 바인딩을 사용하세요.** 스트림은 연속적이거나 요청 없이 전달되는 데이터를 위한 것입니다. 값을 반환하는 호출은 바인딩된 메서드로 구현하는 편이 더 간단합니다.
- **애플리케이션 이벤트에는 `Emit`/`On`을 사용하세요.** 이벤트는 모든 리스너에 배포되며, 별도로 오랫동안 검증된 시스템입니다. 스트림은 지점 간 통신입니다.
- **`InitialHTML` 창에서는 사용할 수 없습니다.** 이러한 창은 `origin === "null"`로 로드되므로 애셋 서버에 전혀 접근할 수 없습니다.
