---
title: "원시 메시지"
description: "성능이 중요한 애플리케이션을 위한 맞춤형 프런트엔드-백엔드 통신 구현"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

원시 메시지는 표준 바인딩 시스템을 우회하여 프런트엔드와 백엔드 사이에 저수준 통신 채널을 제공합니다. 편의성을 낮추는 대신 속도를 높이는 방식입니다.

## 원시 메시지를 사용해야 하는 경우

원시 메시지는 다음과 같은 극단적인 예외 상황에 가장 적합합니다.

- **초고빈도 업데이트** - 매 마이크로초가 중요하며 초당 수천 개의 메시지를 처리해야 하는 경우
- **맞춤형 메시지 프로토콜** - 전송 형식을 완전히 제어해야 하는 경우

@note{type="tip"}
거의 모든 사용 사례에서는 무시해도 될 정도의 오버헤드로 타입 안전성, 자동 직렬화 및 더 나은 개발자 경험을 제공하는 표준 [서비스 바인딩](/features/bindings/services/)을 사용하는 것이 좋습니다.

@end

## 백엔드 설정

애플리케이션 옵션에서 `RawMessageHandler`을 구성하세요.

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### 핸들러 시그니처

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| 매개변수 | 타입 | 설명 |
| --- | --- | --- |
| `window` | `Window` | 메시지를 보낸 창 |
| `message` | `string` | 원시 메시지 내용 |
| `originInfo` | `*application.OriginInfo` | 메시지 출처에 관한 오리진 정보 |

#### OriginInfo 구조체

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| 필드 | 타입 | 설명 |
| --- | --- | --- |
| `Origin` | `string` | 메시지를 보낸 문서의 오리진 URL |
| `TopOrigin` | `string` | 최상위 오리진 URL(iframe에서는 Origin과 다를 수 있음) |
| `IsMainFrame` | `bool` | 메시지가 메인 프레임에서 시작되었는지 여부 |

#### 플랫폼별 제공 여부

- **macOS**: `Origin` 및 `IsMainFrame`이 제공됩니다.
- **Windows**: `Origin` 및 `TopOrigin`이 제공됩니다.
- **Linux**: `Origin`만 제공됩니다.

### 오리진 검증

@note{type="caution"}
메시지가 핸들러에 도착했다는 이유만으로 안전하다고 가정해서는 안 됩니다. 민감한 작업이나 상태를 변경하는 작업을 처리하기 전에 반드시 오리진 정보를 검증해야 합니다.

@end

**수신 메시지를 처리하기 전에 항상 오리진을 확인하세요.** `originInfo` 매개변수는 무단 액세스를 방지하기 위해 반드시 검증해야 하는 중요한 보안 정보를 제공합니다. 악성 콘텐츠, 침해된 콘텐츠 또는 의도하지 않은 스크립트가 원시 메시지를 보낼 수 있습니다. 오리진을 검증하지 않으면 신뢰할 수 없는 출처에서 온 명령을 처리할 수 있습니다. `originInfo`을 사용하여 메시지가 예상한 출처에서 오는지 확인하세요.

### 주요 검증 사항

- **항상 `Origin`** 확인 - 오리진이 예상한 신뢰할 수 있는 출처와 일치하는지 확인하세요(일반적으로 로컬 애셋의 경우 `wails://wails` 또는 `http://wails.localhost`, 그 밖의 경우 애플리케이션의 특정 오리진).
- **`IsMainFrame`** 검증(macOS) - 메시지가 iframe에서 오는지 확인하세요. 이는 서로 다른 보안 컨텍스트를 가진 임베디드 콘텐츠임을 나타낼 수 있습니다.
- **`TopOrigin`** 사용(Windows) - 프레임 콘텐츠를 처리할 때 최상위 오리진을 확인하세요.
- **예상하지 않은 오리진 거부** - 명시적으로 허용하지 않은 오리진의 메시지를 거부하여 안전하게 실패하도록 하세요.

@note{type="info"}
`wails:` 접두사가 붙은 메시지는 Wails 내부 통신용으로 예약되어 있으며 핸들러에 전달되지 않습니다.

@end

## 프런트엔드 설정

`System.invoke()`을 사용하여 원시 메시지를 보내세요.

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### 사전 빌드된 번들 사용

npm을 사용하지 않는 경우 전역 `wails` 객체를 통해 `invoke`에 액세스하세요.

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## 구조화된 메시지

복잡한 데이터는 JSON으로 직렬화하세요.

### 프런트엔드

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### 백엔드

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## 성능 비교

| 방식 | 오버헤드 | 타입 안전성 | 사용 사례 |
| --- | --- | --- | --- |
| 서비스 바인딩 | 높음 | 완전 지원 | 범용 |
| 원시 메시지 | 최소 | 수동 | 고빈도 및 성능이 중요한 경우 |

### 벤치마크 예제

단순한 페이로드의 경우 원시 메시지는 서비스 바인딩보다 초당 훨씬 더 많은 메시지를 처리할 수 있습니다:

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## 전체 예제

다음은 간단한 명령 프로토콜을 구현한 전체 예제입니다:

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## 권장 사례

### 권장 사항

- 실제로 성능이 중요한 경로에 원시 메시지를 사용하세요
- 핸들러에서 적절한 오류 처리를 구현하세요
- 프런트엔드로 응답을 보낼 때 이벤트를 사용하세요
- 구조화된 데이터에는 JSON 사용을 고려하세요
- 블로킹을 방지하도록 메시지를 빠르게 처리하세요

### 피해야 할 사항

- 서비스 바인딩으로 충분한 경우 원시 메시지를 사용하지 마세요
- 수신 메시지의 유효성 검사를 빠뜨리지 마세요
- 장시간 실행되는 작업으로 핸들러를 블로킹하지 마세요(goroutine을 사용하세요)
- 응답을 특정 창으로 보내야 할 때 window 매개변수를 무시하지 마세요

## 다중 창 고려 사항

`window` 매개변수는 메시지를 보낸 창을 식별하므로 다음 작업을 수행할 수 있습니다:

- 올바른 창으로 응답 보내기
- 창별 동작 구현하기
- 디버깅을 위해 메시지 출처 추적하기

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## 다음 단계

- [서비스 바인딩](/features/bindings/services/) - 대부분의 애플리케이션에 사용하는 표준 접근 방식
- [이벤트](/guides/events-reference/) - 백엔드에서 프런트엔드로 통신하기 위한 이벤트 시스템
- [성능](/guides/performance/) - 일반적인 성능 최적화
