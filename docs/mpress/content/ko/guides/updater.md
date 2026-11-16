---
title: "업데이터"
description: "Wails v3용 인앱 자동 업데이트 — 교체 가능한 프로바이더, 암호학적 검증, 원자적 교체, 테마를 적용하거나 대체할 수 있는 기본 UI를 제공합니다."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

업데이터를 사용하면 다운로드, 검증, 교체 파이프라인을 직접 구축하지 않고도 인앱 소프트웨어 업데이트를 배포할 수 있습니다. `app.Updater`을 기반으로 작동하며, 하나 이상의 교체 가능한 `Provider`(GitHub Releases, keygen.sh, Sparkle AppCast, 개방형 Wails Update Manifest 프로토콜 또는 자체 구현)을 받아 구성된 공개 키로 다운로드를 인증하고, 실행 중인 바이너리를 안전하게 교체하며, 모든 상태 전환을 표준 Wails 이벤트 버스를 통해 알립니다.

![업데이트 준비 완료 상태의 기본 업데이터 창 — 상태별 아이콘, 버전 표시(v1.0.0 → v2.0.1 · 8.8 MB), GFM 표를 포함해 Markdown으로 렌더링된 릴리스 노트, 단일 기본 작업을 표시합니다.](/assets/updater/default-window-ready.png)

## 빠른 시작

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

그러면 프레임워크의 업데이트 창이 열리고 GitHub를 확인한 후, 플랫폼용 자산을 다운로드하고 검증하여 바이너리를 교체한 다음 사용자가 다시 시작할 때까지 기다립니다.

## 수명 주기

`app.Updater`은 다음 상태를 갖는 상태 머신입니다(`updater.State`):

| 상태 | 시점 |
| --- | --- |
| `unconfigured` | `Init`이 호출되기 전 |
| `idle` | `Init` 이후, 확인 작업 전 |
| `checking` | `Check`이 진행 중 |
| `up-to-date` | 최근 프로바이더 응답에서 호출자의 버전이 최신이라고 확인됨 |
| `available` | 새 릴리스를 찾았지만 아직 다운로드하지 않음 |
| `downloading` | 프로바이더에서 바이트를 스트리밍하는 중 |
| `verifying` | 다운로드가 완료되어 서명/다이제스트를 확인하는 중 |
| `installing` | 검증된 바이트의 압축을 풀고 이름을 변경하여 스테이징 디렉터리에 배치하는 중 |
| `ready` | 업데이트 스테이징 완료. 적용하려면 `Restart` 호출 |
| `error` | 이전 단계 중 하나가 실패함 |

언제든지 `app.Updater.State()`으로 현재 상태를 확인할 수 있습니다. 상태가 전환될 때마다 Wails 이벤트도 발생합니다([이벤트](#heading-2) 참조).

`Restart`은 실행 중인 앱에 종료를 요청하기 전에 헬퍼가 `application.New`에 도달할 때까지 기다립니다. 기본 시작 제한 시간은 30초입니다. 앱이 `application.New` 전에 긴 초기화 작업을 수행한다면 `Config.HelperReadyTimeout`을 `time.Minute`과 같이 더 긴 시간으로 설정하세요. 0을 지정하면 기본값을 사용하며, 음수 시간 값은 거부됩니다. 시작 시간이 초과되면 `Restart`은 `updater.ErrHelperNotReady`을 반환하고 실행 중인 앱을 열린 상태로 유지합니다.

업데이트는 설치된 바이너리 옆의 `wails-update-*` 디렉터리(macOS에서는 `.app` 번들 옆)에 스테이징되므로 해당 디렉터리에 쓰기 권한이 있어야 합니다.

AppImage 내부에서는 읽기 전용 마운트 대신 `$APPIMAGE`이 지정하는 `.AppImage` 파일을 업데이터가 교체하고 다시 실행합니다.

기본 창은 현재 상태를 자동으로 반영합니다. 예를 들어 `Check`이 업그레이드가 없다고 반환하면 사용자에게 다음 화면이 표시되며, **닫기**를 눌러 닫을 수 있습니다.

![최신 상태의 기본 업데이터 창 — 녹색 체크 표시, '최신 버전입니다' 제목, 단일 닫기 버튼을 표시합니다.](/assets/updater/default-window-up-to-date.png)

## 프로바이더

다음 인터페이스를 충족하는 모든 구현이 `Provider`가 될 수 있습니다.

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

트리 내에 네 가지 구현이 포함되어 있습니다.

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

기본 자산 매처는 파일 이름에서 `GOOS` + `GOARCH` 부분 문자열을 기준으로 선택하며, 일반적인 별칭(`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`)을 인식합니다. 사용자 지정 명명 규칙을 사용하려면 다음과 같이 설정하세요.

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset`은 내용이 `<sha256>  <filename>` 행으로 구성된 동일 릴리스의 자산 이름입니다(`sha256sum`과 `shasum -a 256`이 생성하는 형식). 프로바이더는 `Check` 중에 이 자산을 가져와 선택된 아티팩트와 일치하는 행을 찾고, 프레임워크가 다운로드를 검증하도록 `Release.Verification.Digest`을 채웁니다.

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

프로바이더는 keygen.sh의 아티팩트별 SHA-512 체크섬과 Ed25519ph 서명을 프레임워크의 `Release.Verification` 블록에 자동으로 매핑하므로 추가 연결 작업이 필요하지 않습니다.

**토큰 형식:** keygen.sh 토큰에는 역할 접두사(`admi-` / `prod-` / `envi-` / `user-`)가 포함됩니다. 대시보드에 표시되는 원시 UUID는 토큰의 비밀 값이 아니라 *식별자*입니다. 비밀 값은 토큰을 생성할 때만 볼 수 있습니다. 자세한 내용은 keygen.sh의 [인증 문서](https://keygen.sh/docs/api/authentication/)를 참조하세요.

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

기존 Sparkle / WinSparkle 인프라를 변경하지 않고 그대로 사용할 수 있습니다. 피드에서 `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` 및 `sparkle:channel`을 읽습니다.

Sparkle 1의 DSA 서명(`sparkle:dsaSignature`)은 지원되지 않습니다. 해당 서명 체계를 사용하는 프로젝트는 EdDSA(Sparkle 2)로 전환해야 합니다.

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

개방형 [Wails Update Manifest 프로토콜](/reference/update-manifest/)을 사용합니다. 체크섬과 서명을 인라인으로 포함하여 최신 릴리스와 플랫폼별 아티팩트를 설명하는 단일 JSON 문서입니다. 같은 문서를 정적 파일 호스트(S3, GitHub Pages 또는 모든 CDN — 채널별로 모든 플랫폼을 나열하는 매니페스트 하나를 게시)나 동적 업데이트 서버에서 사용할 수 있습니다. 프로바이더는 확인할 때마다 `platform`, `arch`, `version` 및 `channel`을 전송하므로 서버가 정확히 하나의 아티팩트를 반환하거나 라이선스를 기준으로 접근을 제어할 수 있습니다.

URL 자리표시자를 사용하면 정적 레이아웃을 한 줄로 구성할 수 있습니다.

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

구성된 헤더는 모든 매니페스트 요청에 전송됩니다. 아티팩트 다운로드에서는 매니페스트 자체의 호스트이며 `https`에서 `http`로 다운그레이드되지 않는 경우에만 해당 헤더를 재사용합니다. 교차 출처 리디렉션이나 다운그레이드 리디렉션이 발생하면 `Authorization` 헤더를 제거합니다.

게시는 CLI가 담당합니다. `wails3 updater manifest`은 단일 명령으로 릴리스 파일의 다이제스트를 계산하고 서명하며 설명하고, `wails3 updater verify`은 업로드 전에 결과를 다시 확인합니다. [wails3 CLI로 게시하기](/reference/update-manifest/#publishing-with-the-wails3-cli)를 참조하세요.

### 대체 체인

`Config.Providers`에는 순서가 있습니다. 업데이터는 이를 순차적으로 탐색합니다. 릴리스를 반환하는 첫 번째 프로바이더의 결과를 사용하며, "최신 상태"라고 보고하는 첫 번째 프로바이더가 있으면 체인 탐색을 즉시 중단합니다. 대체는 "기본 프로바이더에 연결할 수 없음"을 위한 것이지 "프로바이더 간 의견 불일치"를 위한 것이 아닙니다. 오류가 발생하면 다음 프로바이더로 진행합니다.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### 자체 프로바이더 작성

메서드는 세 개이며, 일반적인 구현에는 약 150줄이 필요합니다. 검증, 원자적 스테이징, 교체 및 창은 Updater가 담당하고, 프로바이더 코드는 다음 릴리스를 결정하여 바이트를 스트리밍합니다.

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

트리 내 프로바이더를 참고하세요. 각 프로바이더는 Go 파일 하나로 구성되어 있습니다.

## 암호학적 검증

릴리스는 `Config.PublicKey`을 신뢰 루트로 사용하는 프레임워크 검증기에 의해 인증됩니다.

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

지원되는 알고리즘(`Release.Verification.SignatureAlgo`):

| 알고리즘 | 서명 대상 | 참고 |
| --- | --- | --- |
| `ed25519` | 아티팩트의 SHA-256 다이제스트 | Sparkle EdDSA에서 사용 |
| `ed25519ph` | Ed25519ph의 사전 해시를 통한 전체 아티팩트(내부적으로 SHA-512 사용) | keygen.sh에서 사용 |
| `ecdsa-p256` | 아티팩트의 SHA-256 다이제스트 | 원시 `r∥s` 서명과 DER 서명을 모두 허용 |

릴리스에 해시는 있지만 서명이 없는 경우에는 다이제스트 전용 방식(`DigestAlgo`: `sha256` / `sha512`)도 지원합니다.

`Config.PublicKey`은 서명 검증의 유일한 신뢰 앵커입니다. 릴리스 소스는 자체 키로 이를 대체할 수 없습니다. `Signature`이 포함된 릴리스는 `Config.PublicKey`가 구성되어 있지 않으면 안전하게 실패합니다. 검증기는 다운로드 중 스트리밍 방식으로 다이제스트를 계산하므로, 수 GB 규모의 업데이트에서도 검증을 위해 디스크를 추가로 순회하지 않습니다.

@note{type="caution" title="다이제스트 전용 방식 ≠ 암호학적 검증"}
`Digest`만 있는 릴리스는 사용자가 제어하는 암호학적 신뢰 루트가 아니라 레지스트리의 TLS와 레지스트리 자체에서 제공하는 무결성 보장을 기준으로 인증됩니다. 비트 손상 감지에는 다이제스트 전용 방식을 사용하고, 침해된 릴리스 파이프라인의 변조에 대비하려면 서명을 사용하세요.

@end

### 서명 키 생성

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

개인 키는 PKCS#8 PEM이고 공개 키는 PKIX PEM입니다. `Config.PublicKey`은 `.pub` 파일을 그대로 허용합니다. 원시 32바이트 키나 그 base64 표현도 허용하며, 인라인 삽입용 값은 `genkey`에서 출력합니다. `wails3 updater manifest -key updater.key ...` 또는 `wails3 updater sign`으로 릴리스에 서명하세요. 자세한 내용은 [wails3 CLI로 게시하기](/reference/update-manifest/#publishing-with-the-wails3-cli)를 참조하세요.

또는 Go에서는 다음과 같이 합니다.

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## 아티팩트 형식

공급자는 게시한 파일의 바이트를 스트리밍하며, 프레임워크는 교체하기 전에 해당 파일의 압축을 풉니다.

- **단일 바이너리**(예: `myapp-linux-amd64`) — 그대로 사용합니다. Linux에서 흔히 사용됩니다.
- **`.zip`** — 제자리에서 압축을 풉니다. 아카이브에는 최상위 항목이 정확히 하나만 있어야 합니다(일반적으로 macOS `.app` 번들 또는 단일 바이너리). macOS에 권장되는 패키징 방식입니다.
- **`.tar.gz`** / **`.tgz`** — 최상위 항목이 하나여야 한다는 동일한 규칙에 따라 제자리에서 압축을 풉니다. 바이너리와 함께 런타임 트리를 배포하는 Linux 배포판에 유용합니다.

최상위 항목이 둘 이상인 아카이브는 거부됩니다. 프레임워크는 디스크상의 단일 대상을 교체하므로, 아카이브에 여러 항목이 들어 있으면 "이 아카이브를 제자리에 교체"하라는 의미가 모호해집니다. v1에서는 `.dmg` 및 `.pkg`(macOS)과 `.msi`(Windows)를 지원하지 않습니다. 대신 번들의 `.zip`을 배포하세요. 압축 해제 과정에는 Zip Slip 방지가 적용되고, 아카이브 루트 밖을 가리키는 심볼릭 링크가 거부되며, 전체 압축 해제 크기(2 GiB)와 항목 수(50 000)가 제한됩니다.

## 기본 창

`app.Updater.CheckAndInstall(ctx)`은 프레임워크가 소유하는 520×540 크기의 창을 열며, 다음 요소를 포함합니다.

- 상태에 따라 달라지는 대표 아이콘(사용 가능/다운로드 중에는 파란색 ↓, 준비 완료/최신 상태에는 녹색 ✓, 오류에는 빨간색 !)
- 버전 필: `v1.0.0 → v2.0.1 · 8.8 MB`
- **렌더링된 Markdown**을 표시하는 스크롤 가능한 릴리스 노트 패널(문단, 굵게/기울임꼴, 목록, GFM 표, 인라인 코드, 펜스 코드 블록, h1–h3, 링크)
- 상태별 단일 기본 동작(설치 / 다시 시작 및 적용 / 다시 시도)
- 고스트 스타일의 보조 동작(이 버전 건너뛰기 / 나중에 알림)
- `prefers-color-scheme`을 통한 다크/라이트 모드
- 전체 크기를 알 수 없을 때 표시되는 불확정 진행 시머 효과

Wails 이벤트 버스에서 `updater:*` 이벤트를 수신하고 `updater:user:*` 동작을 Go로 다시 내보냅니다.

### CSS 변수를 통한 테마 설정

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

기본 스타일시트는 다음 변수를 제공합니다. 필요한 변수를 재정의하세요.

| 변수 | 기본값(라이트) | 기본값(다크) |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | 시스템 스택 | — |

### 템플릿 교체

사용자 정의 HTML을 제공하세요. 이 HTML은 `wails:updater:*` 이벤트를 수신하고 `wails:updater:user:*` 동작을 내보내기만 하면 됩니다.

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML 창은 애셋 서버 원본 없이 로드되므로 `/wails/runtime.js`을 동적으로 가져올 수 없습니다. 이러한 창에서 호스트와 통신하는 방법은 두 가지입니다.

1. **HTML만 작성하세요.** 프레임워크는 `WebviewWindowOptions.AllowSimpleEventEmit = true` 및 `HTML`가 설정된 상태로 열린 모든 창에 최소 `window.wails.Events` 심을 자동으로 삽입합니다. 이는 업데이터의 기본 제공 경로와 BYO 경로에서 사용하는 방식과 정확히 같습니다. 빌드 단계가 필요하지 않습니다. 아래 예제에서도 이 방식을 사용합니다.
2. **선호하는 번들러**(Vite, esbuild, Rollup)로 `@wailsio/runtime`을 번들링하고 빌드 시 사용자 지정 HTML로 가져오세요. `Events.On`은 전적으로 클라이언트 측에서 작동하므로 별도 설정 없이 사용할 수 있습니다. 반면 `Events.Emit`은 런타임의 fetch 전송을 거치며 null 출처에서는 이 전송이 작동하지 않습니다. 따라서 `window._wails.invoke("wails:event:emit:<name>")`을 통해 라우팅하는 런타임의 [`setTransport`](/guides/custom-transport/) 훅을 사용하여 간단한 postMessage 전송을 설치하세요. `window.wails.Events`이 이미 범위에 있으면 프레임워크의 주입은 아무 작업도 하지 않으므로 두 방식이 충돌하지 않습니다.

어느 방식을 사용하든 사용자 지정 HTML에 작성하는 JS는 동일한 `Events.On` / `Events.Emit` API를 호출합니다.

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

이 shim은 단순 이름 이벤트에 필요한 최신 런타임 기능의 일부를 제공합니다. `Events.On(name, cb)`은 구독 해제 함수를 반환하고, `Events.Emit(nameOrEventObject)`은 사용이 제한된 `wails:event:emit:` postMessage 경로를 통해 호스트로 라우팅합니다. 페이지를 불러올 때 인라인 스크립트가 실행되기 전에 한 번만 설치됩니다.

shim을 재정의하고 *싶거나* 다른 방식으로 전체 런타임을 불러오는 경우, 페이지의 첫 번째 `<script>` 태그가 실행되기 전에 `window.wails.Events`을 설정하면 주입을 건너뜁니다.

### 창 장식

HTML을 수정하지 않고 창 옵션(크기, 프레임 없음, 항상 위)을 재정의하세요.

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### 직접 만든 창 사용

직접 생성한 `*application.WebviewWindow`에서 업데이트 흐름을 구동하세요. 업데이터는 해당 창에서 `Show()` / `Close()` / `EmitEvent()`을 호출하며, 무엇을 렌더링할지는 HTML에서 결정합니다.

![분홍색-주황색 그라데이션 배경, 흰색 둥근 카드 하나, 사용자 지정 타이포그래피를 사용하며 동일한 업데이터 이벤트로 표시 상태를 구동하는 ‘직접 만든’ 업데이터 창입니다. 기본 UI를 얼마나 완전하게 대체할 수 있는지 보여 줍니다.](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

HTML에서는 기본 제공 템플릿과 마찬가지로 `window.wails.Events.On` / `Events.Emit`을 사용합니다. 창의 소유자가 프레임워크인지 사용자인지와 관계없이, 프레임워크가 자동 주입하는 shim은 `AllowSimpleEventEmit: true`이 설정된 모든 창에 포함됩니다. API는 [템플릿 교체](#--6)를 참조하세요.

@note{type="caution" title="BYO 업데이터 창에는 `AllowSimpleEventEmit`이 필요합니다"}
보안을 위해 프레임워크는 이 필드가 설정된 경우에만 `wails:event:emit:` postMessage 바로 가기를 허용합니다. 이 필드가 설정되지 않은 창은 호스트 측 사용자 지정 이벤트를 생성할 수 없습니다. 업데이터의 사용자 지정 HTML shim은 이 바로 가기를 통해 `updater:user:*` 이벤트를 내보내므로, BYO 창에서 이 필드를 빠뜨리면 모든 버튼 클릭이 아무 반응 없이 무시됩니다. 사용자가 설치를 클릭해도 아무 일도 일어나지 않습니다.

완전히 제어할 수 없는 HTML(원격 URL, 사용자 제공 콘텐츠)을 불러오는 창에서는 `AllowSimpleEventEmit`을 반드시 **꺼 두세요**. 이 기능을 켜면 XSS 싱크를 포함하여 페이지의 모든 JavaScript가 임의의 `app.Event.On(name, …)` 핸들러를 실행할 수 있습니다. 이 바로 가기는 페이로드 없이 단순 이름만 전달하고 binding/Call 경로에는 접근할 수 없지만, Go 코드의 권한 있는 사용자 지정 이벤트 핸들러가 이벤트 이름만으로 동작한다면 여전히 해당 핸들러를 트리거할 수 있습니다.

프레임워크의 *기본 제공* 업데이터 창에는 이 값이 내부적으로 설정되어 있으므로 BYO 호출자만 기억하면 됩니다.

@end

### 헤드리스

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

창을 전혀 열지 않습니다. 자체 UI(또는 기존 기본 창)에서 `updater:*` 이벤트를 구독하고 버튼 핸들러에서 `app.Updater.CheckAndInstall(ctx)`을 호출하세요. 업데이트를 찾았을 때만 사용자에게 표시해야 하는 주기적 백그라운드 확인이나, 업데이트 흐름을 사용자 지정 설정 패널에 통합하는 앱에 유용합니다.

## 이벤트

Go와 JavaScript 모두 표준 Wails 이벤트 버스를 통해 구독합니다. **전송 문자열을 직접 입력하지 마세요**. 업데이터 패키지(Go) 또는 런타임 패키지(JS)에서 내보낸 상수를 사용하세요. 두 계층은 동일한 이름 집합을 공유하며 회귀 테스트를 통해 동기화됩니다.

### Go에서 사용

상수는 `github.com/wailsapp/wails/v3/pkg/updater`에 있습니다. `app.Event.On(name, fn)`을 통해 구독하세요. 콜백은 `*application.CustomEvent`을 받으며, 이 객체의 `Data` 필드에는 [이벤트 참조](#--8)에 나열된 형식화된 페이로드가 들어 있습니다. JSON으로 디코딩하지 말고 타입 단언을 사용하세요.

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

사용 가능한 모든 Go 상수:

| 상수 | 전송 문자열 |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### JavaScript에서 사용

상수는 `@wailsio/runtime`의 `Updater.Events` 아래에 있습니다. Go와 이름은 같으며, 자동 완성에서 쉽게 찾을 수 있도록 하위 네임스페이스(`User.*`, `Window.*`)별로 구성되어 있습니다.

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

사용자 지정 HTML이 호스트로 *다시* 내보내는 사용자 동작 이벤트는 `Updater.Events.User` 아래에 있습니다.

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### 이벤트 참조

구독 측(호스트 → 페이지):

| 상수(Go) | 상수(JS) | 페이로드 | 발생 시점 |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | 없음 | 각 `Check` 왕복 호출 전 |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check`이 더 최신 릴리스를 찾았을 때 |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | 없음 | `Check`이 최신 상태임을 확인했을 때 |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | 바이트 스트리밍 시작 |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | 다운로드 중 약 10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | 모든 바이트 기록 완료 후, 검증 전 |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | 서명/다이제스트 검사 시작 |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | 압축 해제 및 스테이징 시작 |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | 재시작 대기 중 |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | 어느 단계에서든 실패 |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | 스냅샷 재생 전 세션당 한 번 |

페이지 측(페이지 → 호스트) — 사용자 지정 템플릿을 작성하는 경우 코드에서 구독합니다:

| 상수(Go) | 상수(JS) | 시점 |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | 창 로드가 완료되면 호스트가 현재 상태를 재생 |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | `available` 상태의 기본 동작 |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | `ready` 상태의 기본 동작 |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | "이 버전 건너뛰기" |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | "나중에 다시 알림" |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | 닫기 버튼 |

## API 참조

### `updater.Config`

| 필드 | 타입 | 참고 |
| --- | --- | --- |
| `CurrentVersion` | `string` | **필수.** 릴리스 태그에 사용하는 것과 동일한 문자열(`v` 접두사 제외) |
| `Providers` | `[]updater.Provider` | **필수.** 순서가 지정된 폴백 체인 |
| `PublicKey` | `[]byte` | PEM 또는 원시 바이트. 선택 사항이지만, 서명된 릴리스에 이 값이 없으면 안전하게 실패합니다 |
| `CheckInterval` | `time.Duration` | 0이 아닌 값이면 `CheckAndInstall`을 호출하는 백그라운드 폴링 루프를 시작합니다 |
| `Platform` | `string` | 애셋 선택에 사용할 `runtime.GOOS` 재정의 |
| `Arch` | `string` | 애셋 선택에 사용할 `runtime.GOARCH` 재정의 |
| `Channel` | `string` | 현재는 정보 제공용이며, 공급자별 채널 필터링에 사용됩니다 |
| `Window` | `updater.WindowOption` | `nil`(기본 제공 기본값), `&BuiltinWindow{…}`, `BYOWindow(handle)` 또는 `WindowNone` |

### `*updater.Updater`의 메서드

| 시그니처 | 용도 |
| --- | --- |
| `Init(cfg Config) error` | 구성합니다. 두 번째 호출에서는 `ErrAlreadyConfigured`을 반환합니다 |
| `State() State` | 현재 수명 주기 단계 |
| `CurrentVersion() string` | `Init`에 전달된 버전 |
| `Check(ctx) (*Release, error)` | 공급자 체인을 순회합니다. `(rel, nil)` = 발견, `(nil, nil)` = 최신 상태, `(nil, err)` = 모두 실패 |
| `DownloadAndInstall(ctx) error` | 스트리밍하고 검증한 뒤, 아카이브이면 압축을 해제하고 스테이징합니다. 먼저 `Check`을 호출해야 합니다 |
| `CheckAndInstall(ctx) error` | 편의 메서드: 창을 열고 `Check`을 호출한 다음, 업데이트를 찾으면 `DownloadAndInstall`을 호출합니다 |
| `Restart(ctx) error` | 헬퍼를 생성하고 `Host.Quit`을 호출한 뒤 종료합니다. 헬퍼가 교체 후 다시 실행합니다 |
| `DownloadedPath() string` | 스테이징된 업데이트가 디스크에 있는 위치이며, 없으면 `""`입니다 |
| `SkipVersion(v string)` | `v`을 건너뛴 버전으로 기록합니다. 이후 `Check` 호출에서는 이 버전을 최신 상태로 간주합니다 |
| `SkippedVersion() string` | 현재 건너뛰도록 설정된 버전을 읽습니다 |
| `StopPeriodicCheck()` | `Config.CheckInterval`에서 시작한 타이머를 취소하고 루프가 반환될 때까지 기다립니다 |

### 오류

| 센티널 | 반환 조건 |
| --- | --- |
| `ErrAlreadyConfigured` | 첫 번째 성공 이후의 `Init` |
| `ErrNotConfigured` | `Init` 이전의 모든 작업 |
| `ErrNoPendingRelease` | 사전에 `Check` 없이 실행한 `DownloadAndInstall` |
| `ErrDownloadInProgress` | 다른 작업이 실행 중일 때 호출한 `DownloadAndInstall` |
| `ErrNotReady` | 준비된 업데이트가 없을 때 실행한 `Restart` |

## 교체 작동 방식

`Restart`은 센티널 환경 변수를 설정한 상태로 현재 바이너리를 다시 실행합니다. `application.New`은 시작 시 이를 감지하고 헬퍼 모드로 전환합니다.

1. 헬퍼는 부모 PID가 종료될 때까지 최대 30초 동안 기다립니다(`platformIsAlive`은 Windows에서 `syscall.OpenProcess` + `GetExitCodeProcess`을 통해, Unix에서 `os.FindProcess` + `proc.Signal(syscall.Signal(0))`을 통해 폴링합니다).
2. 헬퍼는 대상을 백업합니다(파일은 복사하고, macOS `.app` 번들 디렉터리는 재귀적으로 복사합니다).
3. 헬퍼는 준비된 아티팩트로 대상을 교체하며, 시도 사이에 500ms의 백오프를 두고 최대 20번 재시도합니다.
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. 이전 inode를 가리키는 열린 파일 디스크립터는 계속 유효합니다.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows에서는 이미지가 아직 매핑된 파일의 이름을 바꿀 수 있지만 삭제할 수는 없습니다. 다음 업데이트에서 헬퍼는 소유 커널 매핑이 해제된 나머지 `.old.*` 형제 항목을 모두 정리합니다.

4. 헬퍼는 새 바이너리에 원래 실행 파일 모드를 복원합니다(다운로드한 파일은 기본 umask로 생성되므로 Unix에서는 `+x`이 제거되며, Windows에서는 아무 작업도 하지 않습니다).
5. 헬퍼는 헬퍼 모드 환경 변수를 제거하고 (이제 교체된) 바이너리를 다시 실행합니다.
6. 헬퍼가 종료됩니다.

실행에 실패하면 헬퍼가 백업을 복원합니다. 부모가 30초 이내에 종료되지 않으면 헬퍼는 대상에 손대기 전에 중단합니다(따라서 종료 대화 상자가 `Quit`을 차단하더라도 사용자는 작동하는 앱을 유지할 수 있습니다).

`.zip`로 배포되는 macOS `.app` 번들(권장 패키징)의 경우, 검증과 준비 완료 사이에 아카이브의 압축을 풀어 헬퍼가 실제 디렉터리를 제자리에 교체할 수 있게 합니다.

## 주기적 확인

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

`CheckInterval > 0`이면 백그라운드 goroutine이 설정된 간격마다 `CheckAndInstall`을 호출합니다. 다른 흐름(확인/다운로드/검증/설치)이 이미 진행 중일 때 발생한 틱은 버립니다. 동시 상태 머신은 지원되지 않습니다.

무언가 발견된 경우에만 표시하는 조용한 백그라운드 폴링을 사용하려면 `Window: updater.WindowNone`을 설정하고 자체 UI에서 `EventUpdateAvailable`에 응답하세요.

## 건너뛰기 및 나중에 알림

기본 창의 "이 버전 건너뛰기" 버튼은 `SkipVersion(rel.Version)`을 통해 사용 가능한 버전을 기록합니다. 이후의 `Check`은 같은 버전을 발견하면 최신 상태로 간주합니다(사용자가 `CurrentVersion`을 업데이트할 때까지이며, 이는 `Restart`이 성공하면 자동으로 수행됩니다). "나중에 알림"은 아무것도 기록하지 않고 창만 닫습니다.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## 배포 체크리스트

업데이터가 설치할 릴리스를 게시하기 전에 다음 사항을 확인하세요.

1. **올바른 아카이브 형식을 선택하세요.** macOS: `.app` 번들의 `.zip`. Linux: 단일 바이너리 또는 `.tar.gz`. Windows: 단일 `.exe` 또는 `.zip`. `.dmg` / `.msi` / `.pkg`은 지원되지 않습니다.
2. **아티팩트에 서명하세요.** `Config.PublicKey`과 일치하는 비공개 키를 사용합니다. 공급자가 게시하는 피드(keygen.sh, AppCast)의 경우 각 공급자의 서명 워크플로를 따르세요. `ChecksumAsset`을 사용하는 GitHub Releases의 경우 `sha256sum` / `shasum -a 256`으로 `SHA256SUMS` 파일을 생성하세요.
3. **버전 문자열을 일치시키세요.** `Config.CurrentVersion`과 릴리스의 버전 태그가 정확히 일치해야 합니다(예: `1.0.0` ↔ 태그 `v1.0.0`. 선행 `v`은 공급자 측에서 제거됩니다).
4. **출시하기 전에 대상 플랫폼에서 교체를 테스트하세요.** 최소 한 번은 테스트해야 합니다. 코드 서명, 공증 및 Gatekeeper 처리는 플랫폼마다 다르며 업데이터 자체에서 처리하지 않습니다.

## 문제 해결

**"서명에 공개 키가 필요하지만 설정되지 않음"** — 릴리스에 `Signature` 필드가 있지만 `Config.PublicKey`이 비어 있습니다. 공개 키를 설정하거나 서명을 포함하지 않도록 릴리스 파이프라인을 변경하세요.

**"다이제스트 불일치"** — 다운로드한 바이트가 공급자가 명시한 내용과 일치하지 않습니다. 일반적으로 네트워크 문제로 인한 불완전한 다운로드 또는 손상된 아티팩트가 원인입니다. 다시 실행하면 해결되는 경우가 많습니다.

**창이 열리지만 즉시 사라지고 Markdown이나 진행 상황이 표시되지 않음** — 사용자 지정 HTML에서 `wails:runtime:ready`을 호출하지 않았습니다. [템플릿 교체](#--6) shim을 참조하세요.

**Windows 업데이트가 완료되지 않고 헬퍼 로그에 "remove old (attempt N): Access is denied"가 표시됨** — 이 PR의 `de764fb` 이전 버전에서만 발생합니다. 현재 구현은 이 문제가 없는 옆으로 이름 바꾸기 방식을 사용합니다. 업그레이드하세요.

**macOS Gatekeeper가 교체된 바이너리를 차단함** — 코드 서명은 전체 과정에서 보존되어야 합니다. 원본 `.app`에 서명하고, 빌드 파이프라인이 업데이트 시점에 엔타이틀먼트를 수정한다면 다시 실행되는 바이너리에도 *반드시* 다시 서명하세요.

## 함께 보기

- 실행 가능한 예제: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- 테스트 데모 저장소: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- 튜토리얼: [Wails 앱에 자동 업데이트 추가하기](/tutorials/04-self-update-a-wails-app/)
