---
title: "애셋 서버"
description: "Wails v3가 개발 및 프로덕션 환경에서 웹 애셋을 제공하고 임베드하는 방식"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## 개요

모든 Wails 애플리케이션은 다음 요소를 결합한 <strong>단일 네이티브 실행 파일</strong>로 배포됩니다.

1. *Go* 백엔드
2. *Web* 프런트엔드(HTML + JS + CSS)

<strong>애셋 서버</strong>는 이 구성을 가능하게 하는 연결 고리입니다. Go 빌드 태그를 통해 컴파일 시 선택되는 <strong>두 가지 작동 모드</strong>가 있습니다.

| 모드 | 태그 | 용도 |
| --- | --- | --- |
| **개발** | `//go:build !production` | 핫 리로드를 활용한 빠른 반복 개발 |
| **프로덕션** | `//go:build production` | 외부 종속성이 없는 임베드된 애셋 |

구현은 `v3/internal/assetserver/`에 있으며 파일이 명확하게 분리되어 있습니다.

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## 개발 모드

### 수명 주기

1. `wails3 dev`가 시작되고 `build/Taskfile.yml`에 정의된 태스크(일반적으로 `npm run dev`)를 실행하여 **프런트엔드 개발 서버를 생성합니다**(Vite, SvelteKit, React-SWC 등).
2. CLI는 `WAILS_VITE_PORT`을 Wails 개발 포트로 설정하고, `FRONTEND_DEVSERVER_URL`을 실행 중인 프레임워크 개발 서버를 가리키는 **전체** URL(`http://host:port` / `https://host:port`)로 설정합니다. `internal/commands/dev.go`을 참조하세요.
3. 개발용 애셋 서버(`build_dev.go`의 `//go:build !production`을 통해 컴파일에 포함됨)는 `GetDevServerURL()`을 통해 `FRONTEND_DEVSERVER_URL`을 읽고 런타임 이외의 트래픽을 해당 서버로 리버스 프록시합니다.
4. 정적 파일(`/assets/logo.svg`)은 속도를 높이기 위해 `asset_fileserver.go`을 통해 **디스크에서 직접 제공할 수 있으며**, 인식되지 않는 요청은 모두 프레임워크 개발 서버로 **프록시되어** *즉각적인* 핫 모듈 교체를 제공합니다.

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### 기능

- **라이브 리로드** — Vite/SvelteKit 등이 WebSocket을 통해 HMR을 주입하며, 개발용 애셋 서버가 이를 투명하게 프록시합니다.
- **소스 맵 지원** — 애셋이 번들로 묶이지 않으므로 브라우저 개발자 도구에서 오류를 원본 소스에 매핑할 수 있습니다.
- **Go 재컴파일 불필요** — 프런트엔드만 다시 빌드됩니다. `.go` 파일을 변경하기 전까지 Go 코드는 계속 실행됩니다.

### 프레임워크 전환

개발 프록시는 **프레임워크에 종속되지 않습니다**. Wails CLI는 개발 태스크를 시작할 때 환경 변수 두 개를 제공합니다.

| 환경 변수 | 소스 | 의미 |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go`(`wailsVitePort` 상수) | 기본 개발 포트(`--port`이 전달되지 않으면 9245) — Vite 설정에서 이 값을 따르는 것이 좋습니다. |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | Wails가 프록시할 전체 URL입니다. Go에서는 `assetserver.GetDevServerURL()`(`build_dev.go`)을 통해 읽습니다. |

v3 트리에는 `VITE_PORT`, `FRONTEND_DEV_PORT` 또는 `WAILSDEV_VERBOSE` 환경 변수가 없습니다.

새 템플릿을 추가하고 → 해당 개발 태스크를 정의하면 → 애셋 서버가 별도 설정 없이 작동합니다.

---

## 프로덕션 모드

`wails3 build`을 실행하면 파이프라인은 다음 작업을 수행합니다.

1. 프런트엔드 **프로덕션 빌드**(`npm run build`)를 실행하여 `frontend/dist/**`을 생성합니다.
2. 애플리케이션 자체 패키지의 `go:embed`(일반적으로 `main.go` 옆의 `//go:embed all:frontend/dist`)를 통해 해당 디렉터리를 앱에 **임베드합니다**.
3. `-tags production`을 사용하여 Go 바이너리를 컴파일합니다. 이 값은 Taskfile 래퍼가 `EXTRA_TAGS`을 통해 전달합니다.

`internal/assetserver/build_production.go`은 프로덕션 코드 경로로 전환하는 빌드 태그 스텁입니다. `internal/assetserver/bundled_assetserver.go`은 **직접 작성된** 파일입니다. `bundledassets/`에 있는 런타임 JS를 래핑하며 생성된 파일이 아닙니다.

### 요청 처리

실제 핸들러는 `internal/assetserver/assetserver.go` / `asset_fileserver.go`입니다. 개념적으로 다음과 같이 작동합니다.

1. 요청된 경로에서 임베드된 정적 애셋을 찾습니다.
2. SPA 라우팅을 위해 `index.html`으로 폴백합니다.
3. 확장자를 알 수 없으면 콘텐츠 유형을 감지합니다(`content_type_sniffer.go`).
4. 적절한 캐시 헤더를 설정합니다.

- **MIME 감지** — 확장자가 없는 파일은 처음 약 512바이트(`content_type_sniffer.go`)에서 콘텐츠 유형을 감지하고, 그 결과를 `mimecache.go` / `ringqueue.go`에 캐시합니다.
- **보안 헤더** — `file://` 탐색을 허용하지 않고 `nosniff`을 설정합니다.

모든 항목이 임베드되므로 배포되는 바이너리에는 Windows에서도 **외부 종속성이 없습니다**.

---

## 개발 ↔ 프로덕션 연결

`pkg/application`의 관점에서 두 모드는 <strong>동일한 공개 인터페이스</strong>를 제공합니다. 즉, `Handler http.Handler`가 있는 `AssetOptions` 구조체와 `internal/assetserver/` 내부의 미들웨어 및 수명 주기 연결입니다. 개발/프로덕션 전환은 전적으로 Go 빌드 태그를 통해 이루어지므로 두 모드에서 애플리케이션 코드는 동일합니다.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## 프런트엔드 프레임워크 통합 방식

### 템플릿

함께 제공되는 각 템플릿(React, Vue, Svelte, Solid, Vanilla 등)에는 다음 항목이 포함됩니다.

- `build/Taskfile.yml`
- `frontend/vite.config.ts`(또는 이에 상응하는 항목)

각 템플릿의 Vite 또는 이에 상응하는 구성은 `WAILS_VITE_PORT`을 읽고 개발 서버를 해당 포트에 바인딩합니다. 그러면 CLI가 앱 내 프록시에서 사용할 실제 `FRONTEND_DEVSERVER_URL`을 게시합니다.

프레임워크는 Go와 완전히 분리된 상태로 유지됩니다.

- 빌드 시 Wails JS SDK를 가져올 필요가 없습니다. `/wails/runtime.js`은 런타임에 애셋 서버에서 제공됩니다.
- HTTP 개발 서버가 있는 프레임워크라면 무엇이든 연동할 수 있습니다.

---

## 확장 및 사용자 지정

사용자 지정 헤더, 인증 또는 gzip이 필요한가요?

1. `middleware.Middleware`을 정의합니다. 이는 `internal/assetserver/middleware.go`에 선언된 `func(http.Handler) http.Handler`의 별칭입니다.
2. `internal/assetserver/options.go`에서 제공하는 구성을 통해 이를 `application.AssetOptions`에 연결합니다.
3. 개발 환경과 프로덕션 환경에서 동작은 동일하며, 모드별 미들웨어 목록은 없습니다.

---

## 주요 소스 파일

| 파일 | 역할 |
| --- | --- |
| `build_dev.go` / `build_production.go` | 개발 환경과 프로덕션 환경 중 하나를 선택하는 빌드 태그 래퍼 |
| `assetserver.go` / `asset_fileserver.go` | 핵심 HTTP 핸들러 |
| `assetserver_dev.go` | `FRONTEND_DEVSERVER_URL`으로 연결되는 리버스 프록시 |
| `bundled_assetserver.go` | `bundledassets/`을 감싸는 수동 작성 래퍼 |
| `options.go` | `application.AssetOptions` 대상 구성 |
| `mimecache.go` / `ringqueue.go` | MIME 캐시 및 소형 LRU |

---

## 주의 사항 및 디버깅

- **프로덕션 환경에서 흰 화면이 표시됨** — 일반적으로 SPA 라우팅 문제입니다. 개발 서버가 알 수 없는 경로에 대해 `index.html`을 제공하고, 내장된 프로덕션 핸들러의 폴백에 도달하는지 확인하세요.
- 개발 환경에서 **404** — Vite 구성이 `WAILS_VITE_PORT`에 바인딩되지 않았거나, CLI가 개발 서버에 접근하지 못해 `FRONTEND_DEVSERVER_URL`을 채울 수 없었습니다.
- **대용량 애셋** — 임베딩하면 바이너리 크기가 커집니다. 대용량 미디어는 별도의 오리진에서 제공하거나 사용자 지정 `http.Handler`을 통해 스트리밍하세요.

---

이제 Wails <strong>애셋 서버</strong>가 **개발** 환경과 **프로덕션** 환경 모두에서 웹 코드를 네이티브 창에 제공하는 방식을 이해했습니다. 이 계층을 숙지하면 로딩 문제를 디버깅하고, 미들웨어를 추가하거나, 완전히 다른 프런트엔드 도구 체인으로도 자신 있게 교체할 수 있습니다.
