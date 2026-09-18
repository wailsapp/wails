---
title: "사용자 정의 전송 계층 만들기"
description: "Wails v3용 사용자 정의 IPC 전송 계층을 만들고 맞춤 설정하는 방법을 알아봅니다"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3에서는 생성된 모든 바인딩과 이벤트 통신을 그대로 유지하면서 사용자 정의 IPC 전송 계층을 제공할 수 있습니다. 이를 통해 기본 HTTP fetch 기반 전송을 WebSocket, 사용자 정의 프로토콜 또는 다른 전송 메커니즘으로 대체할 수 있습니다.

## 개요

기본적으로 Wails는 프런트엔드에서 HTTP fetch 요청을 보내 `/wails/runtime`을 통해 백엔드와 통신합니다. 사용자 정의 전송 API를 사용하면 다음 작업을 수행할 수 있습니다.

- HTTP 전송을 WebSocket, gRPC 또는 임의의 사용자 정의 프로토콜로 대체
- Wails 코드 생성과의 완전한 호환성 유지
- 기존의 모든 바인딩, 이벤트, 대화 상자 및 기타 Wails 기능 유지
- 자체 연결 관리, 인증 및 오류 처리 구현

## 아키텍처

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## 사용법

### 1. Transport 인터페이스 구현

`Transport` 인터페이스를 구현하여 사용자 정의 전송을 만드세요.

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. 애플리케이션 구성

사용자 정의 전송을 애플리케이션 옵션으로 전달하세요.

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. 프런트엔드 런타임 수정

사용자 정의 전송을 사용하는 경우 HTTP fetch 대신 해당 전송을 사용하도록 프런트엔드 런타임을 수정해야 합니다. 요청 처리에 사용할 `RuntimeTransport` 인터페이스를 구현하세요.

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## 참고 사항

- 사용자 정의 전송을 지정하지 않으면 기본 HTTP 전송이 계속 작동합니다.
- 생성된 바인딩은 변경되지 않으며 전송 계층만 변경됩니다.
- 이벤트, 대화 상자, 클립보드 및 기타 모든 Wails 기능은 별도의 조정 없이 작동합니다.
- 사용자 정의 전송의 오류 처리, 재연결 로직 및 보안은 직접 구현해야 합니다.
- 제공된 WebSocket 예제는 시연용이며 프로덕션 환경에서 사용하려면 보강이 필요할 수 있습니다.

## API 레퍼런스

### Transport 인터페이스

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport 인터페이스(선택 사항)

브라우저 기반으로 배포하거나 사용자 정의 전송을 통해 애셋과 IPC를 모두 제공하려면 `AssetServerTransport` 인터페이스를 구현하세요.

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**이 인터페이스를 구현해야 하는 경우:**

- WebView 대신 브라우저에서 앱을 실행하는 경우
- 사용자 정의 IPC 전송과 함께 HTTP를 통해 애셋을 제공하는 경우
- 네트워크를 통해 접근할 수 있는 애플리케이션을 빌드하는 경우

**구현 예제:**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

`ServeAssets()`이 호출되면 assetHandler는 다음 항목을 제공합니다.

- 모든 정적 애셋(HTML, CSS, JS, 이미지 등)
- `/wails/runtime.js` - Wails 런타임 라이브러리

## 관련 항목

- `transport.go` - 핵심 전송 인터페이스 및 형식
- `messageprocessor.go` - 모든 Wails IPC를 처리하는 기반 메시지 프로세서
