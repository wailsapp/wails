---
title: "자동 업데이트되는 Wails 앱"
description: "`wails3 init`부터 서명된 릴리스 검증과 헬퍼 모드 교체까지, GitHub Releases에서 스스로 업데이트하는 Wails v3 애플리케이션을 빌드합니다."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

이 튜토리얼에서는 새 Wails v3 애플리케이션에 앱 내 업데이터를 추가합니다. 완료하면 앱에서 다음 작업을 수행할 수 있습니다.

- 요청 시 GitHub Releases를 확인합니다(선택적으로 타이머를 사용할 수도 있습니다).
- 현재 실행 중인 OS와 아키텍처에 맞는 에셋을 다운로드합니다.
- 다운로드한 바이트를 기준으로 SHA-256 다이제스트를 검증합니다(선택적으로 Ed25519 서명도 검증할 수 있습니다).
- 프레임워크의 기본 업데이트 창에 릴리스 노트를 표시합니다.
- 별도의 헬퍼 실행 파일을 배포하지 않고 실행 중인 바이너리를 교체한 후 다시 실행합니다.

무료이고 별도의 인프라가 필요하지 않으므로 업데이트 소스로 <strong>GitHub Releases</strong>를 사용하겠습니다. 동일한 패턴을 [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen)와 [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast)에도 적용할 수 있습니다. 완료한 후 [업데이터 가이드](/guides/updater/)를 참조하세요.

@note{type="tip" title="사전 요구 사항"}
- Go 1.25 이상
- `wails3` CLI 설치 완료(`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- 릴리스를 푸시할 수 있는 GitHub 저장소
- [QR 코드 서비스 튜토리얼](/tutorials/01-creating-a-service/)을 숙지하면 도움이 되지만 필수는 아닙니다.

@end

<br/>

@steps
### 새 Wails 앱으로 시작하기
vanilla 템플릿으로 새 프로젝트의 기본 구조를 생성합니다.

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

이제 `main.go`, `frontend/` 및 `Taskfile.yml`이 포함된 디렉터리가 있어야 합니다. 빌드되고 실행되는지 확인합니다.

```bash
wails3 task dev
```

빈 Wails 창이 열려야 합니다. 앱을 종료하고 계속 진행하세요.

### 업데이터 import 추가하기
`main.go`을 열고 import에 다음 두 업데이터 패키지를 추가합니다.

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

이 패키지들은 Updater 자체와 GitHub Releases 공급자를 가져옵니다.

### Updater 구성하기
`app.Updater`은 이미 모든 `*application.App`에 연결되어 있으므로 `Init`만 호출하면 됩니다.

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

이 코드를 `application.New` 뒤, `app.Run()` 앞에 배치합니다.

@note{type="note" title="버전 문자열 형식"}
릴리스 태그에 사용하는 것과 동일한 버전을 전달하되, 앞의 `v`는 <strong>제외</strong>합니다. 공급자 측에서는 태그 이름에서 `v`를 제거합니다. 여기의 `1.0.0` ↔ GitHub의 `v1.0.0`입니다.

@end

### 업데이트를 실행하는 메뉴 항목 추가하기
같은 `main.go`에 "업데이트 확인…" 메뉴 항목을 추가합니다.

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall`은 프레임워크의 업데이트 창을 열고 `Check`을 실행한 다음, 릴리스가 발견되면 `DownloadAndInstall`을 자동으로 실행합니다. 새로운 릴리스가 없으면 창은 "최신 버전" 상태로 열린 채 유지되며, 사용자가 **닫기** 버튼으로 닫습니다.

@note{type="caution" title="goroutine에서 실행하기"}
`CheckAndInstall`은 검증과 설치가 완료될 때까지 차단됩니다. 메뉴 클릭에서 직접 호출하면 UI 스레드가 차단됩니다. `go func()`으로 감싸세요.

@end

### 릴리스가 없는 상태에서 한 번 실행하기
```bash
wails3 task dev
```

<strong>앱 → 업데이트 확인…</strong>을 클릭합니다. 업데이트 창이 잠시 열리고 GitHub API를 호출한 뒤 `1.0.0`보다 새로운 릴리스가 없음을 확인하고, 녹색 ✓와 함께 **최신 버전** 상태로 전환되어야 합니다.

여기서 오류가 발생한다면 일반적으로 다음 중 하나가 원인입니다.

| 증상 | 해결 방법 |
| --- | --- |
| `404 Not Found` | `Repository` 필드가 잘못되었습니다. 반드시 `owner/repo`이어야 합니다. |
| `403 rate-limited` | github.Config에 `Token: "ghp_…"`을 추가합니다(`public_repo` 범위가 있는 PAT를 사용하세요). |
| 네트워크 오류 | 실행 중인 앱에서 `api.github.com`에 연결할 수 있는지 확인합니다. |

### 테스트 릴리스 게시하기
`main.go`의 `currentVersion`을 `1.0.0`로 올립니다(또는 그대로 둡니다). 릴리스에 첨부할 바이너리를 얻도록 한 플랫폼용으로 빌드합니다.

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

바이너리와 같은 위치에 `SHA256SUMS` 파일을 생성합니다.

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

다음과 같은 줄이 하나 이상 표시되어야 합니다.

```
abc123…  updater-tutorial-darwin-arm64.zip
```

이제 GitHub 저장소에 이를 <strong>v2.0.0</strong>로 게시합니다.

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="에셋 이름 지정"}
기본 에셋 매처는 파일 이름에 포함된 `GOOS` + `GOARCH` 하위 문자열을 기준으로 선택합니다. 에셋 이름에 `darwin`(또는 `linux` / `windows`)와 `arm64`(또는 `amd64` / `386`)가 포함되어 있으면 매처가 해당 에셋을 찾습니다. 사용자 지정 매처에 대해서는 [업데이터 가이드](/guides/updater/#github-releases--updaterprovidersgithub)를 참조하세요.

@end

### 앱을 실행하고 업데이트 확인하기
`currentVersion`이 여전히 `1.0.0`인 상태에서 앱을 다시 실행합니다.

```bash
wails3 task dev
```

<strong>앱 → 업데이트 확인…</strong>을 클릭합니다. 이번에는 다음과 같은 화면이 표시되어야 합니다.

![업데이트 준비 완료 상태의 기본 업데이터 창으로, 버전 배지, Markdown으로 렌더링된 릴리스 노트 및 기본 작업 버튼인 다시 시작 및 적용이 표시되어 있습니다.](/assets/updater/default-window-ready.png)

- 대표 아이콘이 파란색 ↓("업데이트 사용 가능")에서 녹색 ✓("업데이트 준비 완료")로 바뀝니다.
- 부제에 `v1.0.0 → v2.0.0 · <size>`이 표시됩니다.
- 릴리스 노트 패널은 굵은 글꼴, 인라인 코드 및 표가 포함된 Markdown을 렌더링합니다.
- 다운로드 중에 진행률 표시줄이 채워집니다(바이너리가 작으므로 금방 완료됩니다).

Updater는 새 바이너리를 임시 디렉터리에 준비합니다. 업데이트를 완료하려면 다음과 같이 합니다.

- <strong>다시 시작 및 적용</strong>을 클릭합니다.
- 앱이 종료되고 헬퍼가 바이너리를 교체한 다음 새 바이너리가 다시 실행됩니다.
- 다시 실행된 앱에는 `currentVersion = "1.0.0"`이 표시되지만(하드코딩했기 때문입니다), 디스크의 바이트는 v2.0.0 빌드와 일치합니다.

실제 앱에서는 새 바이너리가 현재 버전이 v2.0.0임을 인식하고 이후 확인 시 업데이트가 없다고 판단하도록, 빌드할 때 `-ldflags`을 통해 `currentVersion`을 설정합니다.

### `currentVersion`을 빌드에 연결하기
상수를 빌드 시점 변수로 교체합니다.

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

그런 다음 빌드 명령에서 다음과 같이 지정합니다.

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

또는 `git describe --tags`에서 값을 가져오도록 `Taskfile.yml`에 `-ldflags`을 추가합니다.

### 암호학적 서명 추가하기(프로덕션 환경에 권장)
SHA256SUMS 방식은 *무결성*(바이트가 GitHub에 저장된 내용과 일치하는지)은 검증하지만, *진위성*(해당 바이트가 탈취된 메인테이너 계정이 아니라 릴리스 파이프라인에서 생성되었는지)은 검증하지 않습니다. 변조를 방지하려면 Ed25519 키로 각 릴리스에 서명하세요.

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

각 릴리스에서 모든 자산의 SHA-256 다이제스트를 비공개 키로 서명하세요. 다음과 같은 간단한 Go 도우미를 사용할 수 있습니다.

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

기본 GitHub 공급자는 현재 별도의 서명 파일을 가져오지 않습니다. 이를 수행하는 [사용자 지정 공급자를 작성](/guides/updater/#writing-your-own-provider)하거나, 모든 아티팩트를 서버 측에서 서명하고 API를 통해 다이제스트와 서명을 모두 제공하는 <strong>keygen.sh</strong>로 전환할 수 있습니다.

공개 키를 앱에 포함하세요.

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

`PublicKey`을 설정하면 `Signature`을 포함해 배포되는 모든 릴리스가 이 키를 사용한 검증을 통과해야 합니다. 릴리스 소스는 자체 키로 바꿔치기할 수 없습니다. 이것이 바로 빌드 시 대역 외 방식으로 키를 고정하는 이유입니다.

### 창 사용자 지정
기본 창은 일반적인 사용 사례를 지원합니다. 더 세밀하게 제어해야 한다면 다음 세 가지 확장 방법 중 원하는 사용자 지정 수준에 따라 하나를 선택하세요.

@tabs
[CSS만 사용]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

전체 변수 목록은 [CSS 변수를 통한 테마 설정](/guides/updater/#theme-via-css-variables) 섹션을 참조하세요.

[사용자 지정 HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

HTML에서 Wails 이벤트 채널을 통해 `updater:*` 이벤트를 구독하고 `updater:user:*` 작업을 내보내야 합니다. JS 심에 대해서는 [템플릿 교체](/guides/updater/#replace-the-template)를 참조하세요.

[자체 창 사용]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

이미 자체 창 인프라가 있고 업데이터가 별도의 창을 여는 대신 해당 창을 제어하도록 하려는 경우에 유용합니다. 기본 창과 동일한 업데이터 이벤트로 구동되는 완전한 사용자 지정 HTML 템플릿은 다음과 같습니다.

![분홍색-주황색 그라데이션 배경과 사용자 지정 둥근 카드 레이아웃을 적용한 자체 업데이터 창으로, 기본 UI를 완전히 교체할 수 있음을 보여 줍니다.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` 필수 설정"}
업데이터의 사용자 지정 HTML 심은 `wails:event:emit:` postMessage 단축 경로를 통해 설치 / 건너뛰기 / 나중에 알림 / 다시 시작을 구동하며, 보안을 위해 이 필드가 설정된 경우에만 해당 단축 경로를 사용할 수 있습니다. 이 설정을 누락하면 버튼을 눌러도 아무 동작 없이 조용히 무시됩니다. 완전히 제어할 수 없는 HTML을 로드하는 창에서는 이 설정을 활성화하지 마세요. 위협 모델은 가이드의 [자체 창 사용](/guides/updater/#bring-your-own-window) 섹션을 참조하세요.

@end

@end

### 백그라운드에서 자동 검사 실행
메뉴 클릭 대신 또는 메뉴 클릭과 함께 타이머로 검사하려면 다음과 같이 설정하세요.

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

타이머가 실행될 때마다 수동으로 클릭했을 때와 동일한 `CheckAndInstall` 흐름이 실행됩니다. 실제로 업데이트가 발견될 때까지 주기적 검사를 표시하지 않으려면 `Window: updater.WindowNone`을 설정하세요. 그런 다음 표시할 UX를 직접 결정할 수 있도록 `EventUpdateAvailable`을 구독하세요.

@end

## 완료되었습니다

이제 Wails 앱에서 다음 기능을 사용할 수 있습니다.

- 요청 시 그리고 타이머에 따라 GitHub Releases에서 업데이트를 확인합니다.
- 세련된 기본 창에서 릴리스 정보를 Markdown으로 렌더링합니다.
- 게시한 SHA-256 다이제스트를 기준으로 다운로드를 검증합니다.
- 선택적으로 빌드 시 포함한 공개 키를 기준으로 Ed25519 서명을 검증합니다.
- 실행 중인 바이너리를 제자리에서 교체하고 자동으로 다시 실행합니다.

## 다음 단계

- [업데이터 가이드](/guides/updater/)에서 전체 API 레퍼런스, 모든 이벤트와 구성 옵션, 도우미 모드의 교체 메커니즘을 확인할 수 있습니다.
- 복제하여 사용할 수 있는 완전한 동작 예제는 [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)에서 확인하세요.
- 테스트 대상 저장소 [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)에서 권장 릴리스 자산 구조를 확인할 수 있습니다.

## 프로덕션 환경에서 주의할 사항

- **macOS에서 코드 서명** — Gatekeeper는 교체되는 바이너리에 서명과 공증을 요구합니다. 릴리스용 ZIP 파일로 압축하기 *전에* `.app` 번들에 서명하세요. 업데이터는 바이트를 그대로 보존하며 아무것도 다시 서명하지 않습니다.
- **Windows의 바이러스 백신** — 인터넷에서 다운로드한 서명되지 않은 `.exe` 파일은 SmartScreen 경고를 유발할 수 있습니다. Authenticode 인증서로 바이너리에 서명하세요. 그러지 않으면 제한이 엄격한 컴퓨터의 사용자가 앱을 허용 목록에 추가해야 할 수 있습니다.
- **원자적 릴리스** — `SHA256SUMS`과 바이너리를 별도 커밋으로 나누지 말고 함께 게시하세요. 업데이터는 바이너리와 별도로 사이드카를 가져오므로, 둘이 일치하지 않으면 다이제스트 검사가 실패 시 차단 방식으로 종료됩니다.
- **건너뛴 버전** — 기본 창의 "이 버전 건너뛰기" 버튼은 건너뛰기 설정을 로컬에 기록합니다. 중요한 보안 업데이트를 배포할 때는 이전 릴리스를 무시한 사용자의 환경에서 자동으로 건너뛰지 않도록 새 버전 번호를 지정하세요.
