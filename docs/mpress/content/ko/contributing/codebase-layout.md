---
title: "코드베이스 구조"
description: "Wails v3 저장소의 구성 방식과 각 부분이 서로 연동되는 방식"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3는 프레임워크 런타임, CLI, 예제, 문서 및 빌드 도구 체인을 포함하는 <strong>모노레포</strong>에 있습니다. 이 페이지에서는 내부 구조를 살펴보려는 개발자에게 중요한 <em>디렉터리 구조</em>를 안내합니다.

## 최상위 구조

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

이제부터 **`v3/`** 트리를 자세히 살펴보겠습니다.

## `v3/` 루트

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> 프로젝트 템플릿은 `internal/templates/` 아래에 제공됩니다(프레임워크
>
> 스택마다 폴더 하나씩, 그리고 `base/`, `_common/`, `ios/` 포함). 최상위에는 `v3/templates/`
>
> 디렉터리가 없습니다.

### 개념 모델

1. <strong>`pkg/`</strong>는 <em>애플리케이션 개발자가 가져오는 항목</em>을 공개합니다\
2. <strong>`internal/`</strong>에는 <em>내부 기능의 구현 방식</em>이 들어 있습니다\
3. <strong>`cmd/wails3`</strong>는 <em>프로젝트 수명 주기와 빌드</em>를 구동합니다\

나머지는 모두 이 세 축을 지원합니다.

---

## `cmd/` – 명령어

| 경로 | 참고 |
| --- | --- |
| `v3/cmd/wails3` | <strong>CLI 진입점</strong>입니다. 작은 `main.go`가 모든 로직을 `internal/commands`의 패키지에 위임합니다. |
| `internal/commands/*` | 하위 명령어(init, dev, build, doctor 등)입니다. 쉽게 찾을 수 있도록 각각 별도의 파일에 있습니다. |
| `internal/commands/task_wrapper.go` | CLI 플래그와 Taskfile 빌드 파이프라인을 연결합니다. |

CLI는 다음 기능을 담당합니다.

- **프로젝트 스캐폴딩**(`init`, 템플릿 생성)\
- **개발 서버 오케스트레이션**(`dev`, 실시간 다시 로드)\
- **프로덕션 빌드 및 패키징**(`build`, `package`, 플랫폼 래퍼)\
- **진단**(`doctor`)\

---

## `internal/` – 핵심 구현부

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### 주요 하위 패키지

| 패키지 | 역할 | 연결 지점 |
| --- | --- | --- |
| `runtime` | 작은 `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go` 빌드 태그 연결 코드와 `runtime/desktop/` 아래의 임베디드 JS 런타임을 포함합니다. 실제 OS별 창, 클립보드, 대화 상자 및 트레이 코드는 `pkg/application/*_{darwin,linux,windows}.go`에 있습니다. | `pkg/application`를 통해 간접적으로 가져옵니다. |
| `assetserver` | 이중 모드 파일 서버:<br />• 개발: 디스크에서 제공하고 Vite를 프록시함(`build_dev.go`)<br />• 프로덕션: `go:embed`를 통해 에셋을 임베드함(`build_production.go`) | 시작할 때 `pkg/application`에서 초기화합니다. |
| `generator` | Go 소스를 파싱하여 <strong>바인딩 메타데이터</strong>를 생성하며, 이 메타데이터는 이후 TypeScript/JS 스텁 파일과 이벤트 상수를 생성합니다. 진입점: `collect/` + `render/`를 기반으로 하는 `generator.Generate` / `generator.Generator`. | `wails3 generate bindings`에서 실행합니다. |
| `packager` | Linux `deb`/`rpm`/`archlinux` 아티팩트를 생성하는 데 사용하는 `nfpm` 래퍼입니다(`internal/commands/` 아래의 `myapp.DEB`/`.RPM`/`.ARCHLINUX` nfpm 설정으로 구동됨). | `wails3 tool package`에서 호출합니다. macOS DMG / Windows MSIX는 `internal/commands/{dmg,msix.go,webview2/}` 아래에 있습니다. |

지원 유틸리티(예: `s/`, `hash/`, `flags/`)는 내부 관심사가 서로 결합되지 않도록 합니다.

---

## `pkg/` – 공개 API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> `pkg/runtime/`, `pkg/options/` 또는 `pkg/menu/` 패키지는 없습니다. 창/메뉴
>
> 옵션은 `pkg/application`와 같은 위치에 있으며(예: `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), `assetserver/`는 `internal/` 아래에 있습니다.

`pkg/application`는 Wails 프로그램을 부트스트랩합니다.

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

내부적으로는 다음 작업을 수행합니다.

1. `internal/runtime` 빌드 태그 연결 코드와 `pkg/application/`의 OS별 코드를 연결합니다.
2. `internal/assetserver` 인스턴스를 설정합니다.
3. 바인딩 기반 메시지 프로세서를 모두 등록합니다.
4. OS 메인 스레드에 진입합니다.

---

## `internal/templates/` – 스캐폴딩 청사진

`internal/templates/`는 **기본 템플릿**(`base/`, `_common/`, `ios/` 아래의 Go 구조)과 **프런트엔드 스킨**(`vanilla[-ts]`, `react[-ts]`, `react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`, `svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`)을 제공합니다.

`wails3 init -t react`에서 CLI는 다음 작업을 수행합니다.

1. `_common` Go 파일을 복사합니다.
2. 원하는 프런트엔드 팩을 병합합니다.
3. `go mod tidy` 실행(`--skipgomodtidy`로 건너뛸 수 있음)

템플릿을 편집해도 기존 앱에는 **영향을 주지 않으며**, 이후 `init` 실행 시에만 적용됩니다. 공개 예제는 `v3/examples/` 아래에 있으며, 기여자 문서에서 설명하는 자동화된 테스트 스위트를 대신하지 않습니다.

---

## `tasks/` – 릴리스 자동화

Taskfile은 복잡한 교차 컴파일, 버전 업데이트 및 변경 로그 생성을 래핑합니다. `internal/commands/task.go`에서 이를 프로그래밍 방식으로 사용하므로 동일한 로직이 <strong>CLI</strong>와 <strong>CI</strong>를 구동합니다.

---

## 구성 요소의 상호작용 방식

```d2
direction: down
CLI: wails3 CLI
Generator: internal/generator
AssetDev: assetserver(개발)
Packager: internal/packager
AppRuntime: {
  label: 앱 런타임
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: OS API
}
CLI -> Generator: 빌드/생성
CLI -> AssetDev: 개발
CLI -> Packager: 패키징
Generator -> ApplicationPkg: 바인딩
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

<em>CLI → generator → runtime</em>은 <strong>소스</strong>에서 <strong>실행 중인 데스크톱 앱</strong>으로 이어지는 핵심 경로를 형성합니다.

---

## 구조 파악을 위한 팁

| 파악하려는 항목 | 살펴볼 위치 |
| --- | --- |
| 플랫폼 호환 계층 | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go`(window, clipboard, dialogs, systray, mainthread, events_common). Linux cgo: `pkg/application/linux_cgo*.go`. |
| 브리지 프로토콜 | `pkg/application/messageprocessor*.go` |
| 애셋 워크플로 | `internal/assetserver/`(`build_dev.go`와 `build_production.go` 비교) |
| 패키징 흐름 | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| 템플릿 엔진 | `internal/templates/`(`templates.Install`, `templates.GetDefaultTemplates`) |
| 정적 분석 | `internal/generator/{generate.go,collect/,render/}` |

---

이제 저장소의 <strong>전체 구조</strong>를 파악했습니다. 이 구조를 `ripgrep`, IDE의 “파일/기호로 이동”, 예제 앱과 함께 활용하여 원하는 기능을 더 깊이 탐색해 보세요. 즐겁게 개발하세요!
