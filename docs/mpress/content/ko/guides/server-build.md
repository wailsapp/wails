---
title: "서버 빌드"
description: "네이티브 GUI 창 없이 Wails 애플리케이션을 HTTP 서버로 실행합니다"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3은 서버 모드를 지원하므로 네이티브 창을 생성하거나 GUI 종속성을 요구하지 않고 애플리케이션을 순수 HTTP 서버로 실행할 수 있습니다. 따라서 동일한 Wails 애플리케이션을 서버와 컨테이너에 배포하고 웹 브라우저에서 사용할 수 있습니다.

서버 모드는 다음 용도에 유용합니다:

- **Docker/컨테이너 배포** - X11/Wayland 종속성 없이 실행합니다
- **서버 측 애플리케이션** - 브라우저에서 접근할 수 있는 웹 서버로 배포합니다
- **웹 전용 접근** - 데스크톱과 웹에서 동일한 코드베이스를 공유합니다
- **CI/CD 테스트** - 디스플레이 서버 없이 통합 테스트를 실행합니다
- **마이크로서비스** - 헤드리스 백엔드 서비스에서 Wails 바인딩을 사용합니다

## 빠른 시작

`server` 빌드 태그로 서버 모드를 활성화합니다. 애플리케이션 코드는 그대로 유지되며, 이 태그를 지정해 빌드하기만 하면 됩니다:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

다음은 최소 구성 예제입니다:

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

동일한 코드를 데스크톱 모드(태그 없음) 또는 서버 모드(`-tags server` 사용)로 빌드할 수 있습니다.

## 구성

### ServerOptions

`ServerOptions`으로 HTTP 서버를 구성합니다:

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## 기능

### 상태 확인 엔드포인트

`/health`에서 상태 확인 엔드포인트를 자동으로 사용할 수 있습니다:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

다음 용도에 유용합니다:

- Kubernetes 활성/준비 상태 프로브
- 로드 밸런서 상태 확인
- 모니터링 시스템

### 서비스 바인딩

모든 서비스 바인딩은 데스크톱 모드와 동일하게 작동합니다:

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

프런트엔드에서는 표준 Wails 런타임을 사용해 이러한 바인딩을 호출할 수 있습니다:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### 이벤트

서버 모드에서 이벤트는 양방향으로 작동합니다:

- **프런트엔드에서 백엔드로**: 브라우저에서 발생한 이벤트가 HTTP를 통해 전송되어 Go 이벤트 핸들러에 수신됩니다
- **백엔드에서 프런트엔드로**: Go에서 발생한 이벤트가 WebSocket을 통해 연결된 모든 브라우저로 브로드캐스트됩니다

각 브라우저 탭은 고유한 이름(`browser-1`, `browser-2` 등)을 가진 "창"으로 표현되며, `event.Sender`을 통해 접근할 수 있습니다:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

프런트엔드에서는 다음과 같이 사용합니다:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### 정상 종료

서버는 `SIGINT` 및 `SIGTERM` 신호를 정상적으로 처리합니다:

1. 새 연결 수락을 중지합니다
2. 활성 요청이 완료될 때까지 기다립니다(최대 `ShutdownTimeout`)
3. `OnShutdown` 훅을 실행합니다
4. 서비스를 역순으로 종료합니다

## 데스크톱 모드와의 차이점

| 기능 | 데스크톱 모드 | 서버 모드 |
| --- | --- | --- |
| 네이티브 창 | 생성됨 | 브라우저 창(`browser-N`) |
| 시스템 트레이 | 사용 가능 | 사용할 수 없음 |
| 네이티브 대화 상자 | 사용 가능 | 사용할 수 없음 |
| 애플리케이션 메뉴 | 사용 가능 | 사용할 수 없음 |
| 화면 정보 | 사용 가능 | 네이티브 화면 없음. `GetPrimary()`은 `nil`을 반환함 |
| 서비스 바인딩 | 작동함 | 작동함 |
| 이벤트 | 작동함 | 작동함(WebSocket 사용) |
| 애셋 | 웹뷰를 통해 제공 | HTTP를 통해 제공 |
| CGO 필요 여부 | 예 | 아니요 |

### 창 API 동작

서버 모드에서는 창 관련 API를 안전하게 처리합니다:

- `app.Window.NewWithOptions(options)` - nil이 아닌 Go 창 객체를 반환하지만 네이티브 데스크톱 창을 생성하지는 않습니다
- `app.Hide()` / `app.Show()` - 아무 작업도 하지 않습니다
- `app.Screen.GetPrimary()` - `*application.Screen`을 반환합니다. 기본 화면을 사용할 수 없으면 이 값은 `nil`이므로 역참조하기 전에 확인하십시오

Go 창 객체와 네이티브 창의 가용성을 구분하십시오. 서버를 빌드할 때는 데스크톱 전용 동작을 조건부로 실행해야 합니다. Go 창 객체가 nil이 아니더라도 네이티브 창이나 화면이 존재한다는 의미는 아닙니다.

### 지원 플랫폼

서버 모드는 순수 Go로 구현되며 CGO 없이 `android` 및 `ios`을 비롯해 Wails가 지원하는 모든 `GOOS`용으로 컴파일됩니다. 모바일 대상에서는 헤드리스 서버만 빌드합니다:

- 네이티브 모바일 API인 `application.Android` 및 `application.IOS`은 사용할 수 없습니다. 이러한 API를 사용하는 코드는 `!server` 빌드 제약 조건으로 서버 빌드에서 제외해야 합니다.
- 기본 로거는 데스크톱용 로거이며, 개발 빌드에서는 logcat이나 WebView 콘솔이 아니라 stderr에 기록합니다.

## 프로덕션용 빌드

### Task 사용(권장)

`wails3 init`으로 생성한 프로젝트에는 `build:server` 태스크가 포함됩니다:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### 수동 빌드

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails 프로젝트에는 바로 사용할 수 있는 Docker 구성이 포함되어 있습니다. 컨테이너에서 애플리케이션을 빌드하고 실행하려면 다음과 같이 하십시오:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

이것으로 끝입니다! `http://localhost:8080`에서 애플리케이션을 사용할 수 있습니다.

몇 가지 옵션으로 빌드를 사용자 지정할 수 있습니다:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

생성된 `Dockerfile.server`은 distroless 기반의 최소 이미지를 만듭니다. 네트워크 바인딩을 자동으로 처리하므로 컨테이너 외부에서 애플리케이션에 접근할 수 있습니다.

### Docker Compose

더 복잡한 배포에는 다음과 같이 상태 확인이 포함된 Docker Compose 구성을 사용할 수 있습니다:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
상태 확인 예제에서는 `wget`을 사용합니다. distroless 기반 이미지를 사용하는 경우 이미지에 상태 확인 바이너리를 포함하거나 외부 상태 확인 메커니즘(예: Docker의 `curl` 옵션 또는 사이드카 컨테이너)을 사용해야 합니다.

@end

### 사용자 지정 Dockerfile

더 세밀하게 제어해야 한다면 자체 Dockerfile을 만들 수 있습니다. 서버가 컨테이너 외부의 연결을 수락하도록 `WAILS_SERVER_HOST=0.0.0.0`을 설정해야 한다는 점이 핵심입니다.

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## 보안 고려 사항

서버 모드 애플리케이션을 배포할 때는 다음 사항을 준수하세요.

1. **기본적으로 localhost에 바인딩** - 필요한 경우에만 `0.0.0.0` 사용
2. **프로덕션 환경에서 TLS 사용** - `ServerOptions.TLS` 구성
3. **리버스 프록시 뒤에 배치** - 보안을 강화하려면 nginx/traefik 사용
4. **WebSocket과 런타임 호출을 동일 출처로 유지** - `WebSocketOriginPatterns`에는 신뢰할 수 있는 출처만 추가하고 `WebSocketAllowAllOrigins`은 사용하지 않기
5. **모든 입력 검증** - 다른 웹 애플리케이션과 동일한 보안 관행 적용

## 예제

전체 예제는 `v3/examples/server/`에서 확인할 수 있습니다.

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## 환경 변수

코드를 변경하지 않고 서버 구성을 재정의해야 하는 배포 환경을 위해 Wails는 다음 환경 변수를 인식합니다.

| 변수 | 설명 | 기본값 |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | 바인딩할 네트워크 인터페이스 | `localhost` |
| `WAILS_SERVER_PORT` | 수신 대기할 포트 | `8080` |

이러한 환경 변수는 코드의 `ServerOptions`보다 우선합니다. 따라서 Docker 예제에서는 `WAILS_SERVER_HOST=0.0.0.0`을 설정합니다. 이렇게 하면 애플리케이션을 변경하지 않고도 컨테이너가 외부 연결을 수락할 수 있습니다.

## 관련 문서

- [사용자 지정 전송 방식](/guides/custom-transport/) - 고급 IPC 사용자 지정
- [서비스](/features/bindings/services/) - 서비스 바인딩 문서
- [이벤트](/guides/events-reference/) - 이벤트 시스템 문서

### 런타임 요청 크기

`/wails/runtime`에 대한 요청은 JSON 처리 전에 64 MiB로 제한됩니다. `Content-Length`이 없는 요청을 포함하여 이 제한을 초과하는 일반 요청에는 HTTP 413가 반환됩니다. 청크 분할된 런타임 업로드에는 청크당 1 MiB 및 조립된 페이로드당 64 MiB라는 별도의 제한이 그대로 적용됩니다. 필요한 경우 애플리케이션 미들웨어나 리버스 프록시를 사용하여 더 낮은 제한을 적용하세요.

### 런타임 요청 출처

`/wails/runtime`은 `POST` 요청만 수락하며, 다른 메서드에는 HTTP 405를 반환합니다. 브라우저 요청의 `Origin`이 요청 호스트와 같지도 않고 `WebSocketOriginPatterns`와 일치하지도 않으면 HTTP 403가 반환됩니다. `Origin`이 없고 `Sec-Fetch-Site`이 `cross-site` 또는 `same-site`인 요청에도 동일하게 반환됩니다. `curl`이나 기타 비브라우저 클라이언트에서 보내는 요청처럼 두 헤더가 모두 없는 요청은 수락됩니다. 리버스 프록시가 `Host` 헤더를 다시 작성하는 경우 공개 출처를 `WebSocketOriginPatterns`에 추가하세요.

`WebSocketOriginPatterns`은 공개 출처가 애플리케이션에 표시되는 `Host`과 다른 리버스 프록시 배포를 위해 제공됩니다. 다른 출처에서 제공되면서 `/wails/runtime`을 직접 호출하는 프런트엔드는 지원되지 않습니다. 이 엔드포인트는 CORS 헤더를 내보내지 않으며 사전 요청 `OPTIONS`에 HTTP 405로 응답합니다. 프런트엔드를 Wails 서버에서 제공하거나 프런트엔드와 서버를 동일한 프록시 출처 뒤에 배치하세요.

### 루프백 Host 검사

서버가 루프백 주소(기본값 `localhost`, 임의의 `127.0.0.0/8` 주소 또는 `::1`)에 바인딩되면, 이를 `ServerOptions.Host`에서 설정했든 `WAILS_SERVER_HOST`에서 설정했든 관계없이 모든 요청에는 `localhost`, IP 리터럴, 바인딩 호스트 또는 `WebSocketOriginPatterns` 항목의 호스트 부분과 일치하는 호스트를 지정하는 `Host` 헤더가 있어야 합니다. 그 밖의 모든 `Host`에는 HTTP 403가 반환됩니다. 이를 통해 악성 웹 페이지가 자체 호스트 이름을 `127.0.0.1`으로 지정하여 요청이 위의 동일 출처 검사를 통과하게 하는 DNS 리바인딩을 차단합니다. 포트는 비교하지 않으며, `localhost.`(후행 점 포함) 또는 `127.1` 같은 표기는 거부됩니다. `WebSocketAllowAllOrigins`을 설정한 경우에도 이 검사는 계속 적용됩니다. 로컬 리버스 프록시가 공개 `Host`를 전달하는 경우 해당 호스트를 `WebSocketOriginPatterns`에 추가하세요. 컨테이너의 `0.0.0.0`처럼 다른 주소에 대한 바인딩은 영향을 받지 않습니다.
