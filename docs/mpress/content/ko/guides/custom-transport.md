---
title: "사용자 정의 전송 계층 만들기"
description: "Wails v3 사용자 정의 IPC 전송 계층을 만들고 맞춤 설정하는 방법을 알아봅니다"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3에서는 생성된 모든 바인딩과 이벤트 통신을 그대로 유지하면서 사용자 정의 IPC 전송 계층을 제공할 수 있습니다. 이를 통해 기본 HTTP fetch 기반 전송을 WebSocket, 사용자 정의 프로토콜 또는 다른 전송 메커니즘으로 대체할 수 있습니다.

## 개요

기본적으로 Wails는 프런트엔드에서 HTTP fetch 요청을 사용하여 `/wails/runtime`을 통해 백엔드와 통신합니다. 사용자 정의 전송 API를 사용하면 다음 작업을 수행할 수 있습니다.

- HTTP 전송을 WebSocket, gRPC 또는 사용자 정의 프로토콜로 대체
- Wails 코드 생성과 완전한 호환성 유지
- 기존의 모든 바인딩, 이벤트, 대화 상자 및 기타 Wails 기능 유지
- 연결 관리, 인증 및 오류 처리를 직접 구현

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

## 전체 WebSocket 예제 실행

[WebSocket 전송 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport)에는 요청 ID와 응답, 수명 주기 관리, JavaScript 클라이언트 및 이벤트 전달을 포함한 프로토콜의 양쪽 구현이 모두 들어 있습니다. `Start`과 `Stop`만 구현하는 전송 대신 해당 파일을 기반으로 시작하세요.

네이티브 플랫폼 종속 항목이 설치된 Wails 저장소 작업 복사본에서 다음을 실행하세요.

```bash
cd v3/examples/websocket-transport
go run .
```

이 예제는 네이티브 창을 열고 9099 포트에서 WebSocket 엔드포인트를 시작합니다. 신뢰할 수 있는 로컬 개발 환경에서 사용하세요. 네트워크를 통해 접근할 수 있는 변형을 배포하기 전에 리스너를 제한하고, 오리진을 검증하고, 클라이언트를 인증하세요. 예제의 오리진 정책은 허용 범위가 넓으며 프로덕션용 보안 정책이 아닙니다.

### 백엔드 계약

전송은 필수 메서드 세 개를 모두 구현해야 합니다. 예제의 `WebSocketTransport`은 `Start`, `JSClient` 및 `Stop`을 구현하며, `JSClient`은 Wails가 페이지에 삽입하는 JavaScript shim을 제공합니다. 구현에 컴파일 타임 인터페이스 검사를 포함하세요.

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

예제의 생성자를 사용하여 애플리케이션을 구성하세요.

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

이 필드는 일반 자산 및 서비스와 함께 `application.Options`에 속합니다. 생성자 호출뿐만 아니라 구현 파일도 복사하세요. `Start`은 `MessageProcessor`을 전달받으며, `Stop`은 리스너, 활성 연결 및 대기 중인 작업을 닫아야 합니다.

### 프런트엔드 계약

예제의 `websocket-transport.js`은 요청과 응답의 연결 및 실패한 호출의 거부를 포함하여 런타임 전송을 구현합니다. 자산에서 해당 파일을 제공한 후 생성된 서비스 바인딩을 호출하기 전에 초기화하세요.

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

생성된 서비스 바인딩은 변경되지 않습니다. 전송이 클라이언트를 대체하는 경우 애플리케이션에서 사용하는 메서드를 전달하고 해당 오류 의미 체계를 보존해야 합니다. 이벤트 전달에는 선택적 이벤트 전송 계약이 필요하며, 임의의 `Transport` 구현에서 자동으로 제공되지 않습니다. 전체 WebSocket 예제에서 이를 확인할 수 있습니다.

원격 배포에는 TLS와 명시적인 인증 및 오리진 정책을 사용하세요. 브라우저 전용 호스팅의 경우 아래 설명과 같이 애플리케이션 자산도 제공해야 합니다.

## API 참조

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

브라우저 기반 배포를 사용하거나 사용자 정의 전송을 통해 자산과 IPC를 모두 제공하려면 `AssetServerTransport` 인터페이스를 구현하세요.

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

- 웹뷰 대신 브라우저에서 앱 실행
- 사용자 정의 IPC 전송과 함께 HTTP를 통해 자산 제공
- 네트워크를 통해 접근할 수 있는 애플리케이션 빌드

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

`ServeAssets()`이 호출되면 assetHandler는 다음을 제공합니다.

- 모든 정적 자산(HTML, CSS, JS, 이미지 등)
- `/wails/runtime.js` - Wails 런타임 라이브러리

## 관련 문서

- `transport.go` - 핵심 전송 인터페이스 및 형식
- `messageprocessor.go` - 모든 Wails IPC를 처리하는 기반 메시지 프로세서
