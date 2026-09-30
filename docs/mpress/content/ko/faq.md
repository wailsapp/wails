---
title: "자주 묻는 질문"
description: "Wails v3로 애플리케이션을 빌드할 때 자주 묻는 질문에 대한 답변"
slug: "faq"
sourcePath: "faq.md"
---

## 일반

### Wails란 무엇인가요?

Wails는 Go와 웹 기술로 데스크톱 애플리케이션을 빌드하기 위한 프레임워크입니다. 애플리케이션 로직은 Go로 작성하고, 인터페이스는 HTML, CSS 및 JavaScript(또는 원하는 프런트엔드 프레임워크)로 빌드하며, Wails는 이를 운영 체제의 네이티브 웹뷰에서 렌더링합니다. 그 결과 브라우저를 번들로 포함하지 않고 메모리 사용량이 적으며, 일반적으로 약 10MB인 단일 바이너리로 구성된 작고 빠르면서 네이티브처럼 느껴지는 애플리케이션을 만들 수 있습니다.

### Wails는 어떤 플랫폼을 지원하나요?

| 플랫폼 | 요구 사항 |
| --- | --- |
| Windows | AMD64 및 ARM64. [WebView2 런타임](https://developer.microsoft.com/microsoft-edge/webview2/)을 사용합니다. |
| macOS | Intel에서는 10.15 이상(애플리케이션은 10.13 이상을 대상으로 지정 가능), Apple Silicon에서는 11.0 이상. 유니버설 바이너리를 지원합니다. |
| Linux | AMD64 및 ARM64. 기본 스택은 GTK4와 WebKitGTK 6.0입니다(Ubuntu 24.04 이상, Debian 13 이상, Fedora 40 이상 및 유사 배포판). Ubuntu 22.04, Debian 12 및 RHEL 9처럼 WebKit2GTK 4.1만 제공하는 배포판은 레거시 `-tags gtk3` 빌드를 통해 지원됩니다(v3.1까지 제공). WebKit2GTK 4.0만 제공하는 배포판은 지원되지 않습니다. [Linux 빌드 가이드](/guides/build/linux/)를 참조하세요. |
| iOS 및 Android | 실험적 지원입니다. [모바일 가이드](/guides/mobile/)를 참조하세요. |

[서버 빌드](/guides/server-build/)를 사용해 애플리케이션을 일반 웹 앱으로 제공할 수도 있습니다.

언제든 `wails3 doctor`을 실행하여 시스템을 확인하고 플랫폼별 설치 지침을 확인하세요.

### 시작하려면 무엇이 필요한가요?

- Go 1.25 이상
- Node.js 및 npm(프런트엔드 빌드용)
- 플랫폼 도구 체인: Windows에서는 WebView2(10/11에 사전 설치됨), macOS에서는 Xcode Command Line Tools, Linux에서는 `gcc`와 GTK/WebKit 개발 패키지

`wails3 doctor`이 이 모든 항목을 확인하고 정확히 무엇이 누락되었는지 알려 줍니다. 전체 안내는 [설치](/quick-start/installation/)를 참조하세요.

### Wails v3는 프로덕션 환경에서 사용할 준비가 되었나요?

Wails v3는 안정적인 데스크톱 API를 갖춘 베타 소프트웨어입니다. 이미 이를 사용해 프로덕션 환경에서 실행되는 애플리케이션이 있지만, 저희가 3.0을 위한 최종 마무리 작업을 진행하는 동안에는 배포 전에 철저히 테스트하는 것이 좋습니다. 현재 상황은 [프로젝트 상태 페이지](/status/)를 참조하세요. Wails v2는 현재 안정 릴리스이며 계속해서 수정 사항이 제공됩니다.

## 개발

### Go를 알아야 하나요?

기본적인 Go 지식이 있으면 도움이 되지만 전문가일 필요는 없습니다. 애플리케이션 로직은 일반 Go 메서드로 작성하며, 그 밖의 모든 내용은 [튜토리얼](/tutorials/overview/)에서 단계별로 안내합니다. 많은 개발자가 첫 Wails 애플리케이션을 빌드하면서 Go를 익힙니다.

### 선호하는 프런트엔드 프레임워크를 사용할 수 있나요?

예. HTML, CSS 및 JavaScript로 빌드할 수 있다면 Wails에서 사용할 수 있습니다. React, Vue, Svelte 및 바닐라 JavaScript용 템플릿이 제공되며(각각 TypeScript 변형 포함), 그 밖의 프레임워크도 몇 분 안에 연동할 수 있습니다. [프런트엔드 프레임워크](/guides/dev/frontend-frameworks/)를 참조하세요.

### JavaScript에서 Go 함수를 호출하려면 어떻게 하나요?

서비스를 등록하면 Wails가 해당 서비스의 형식 지정 바인딩을 생성합니다:

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

바인딩은 `wails3 dev` 중에 자동으로 다시 생성되며, `wails3 generate bindings`을 사용하여 필요할 때 생성할 수도 있습니다. [서비스](/features/bindings/services/)를 참조하세요.

### TypeScript를 사용할 수 있나요?

예. 바인딩 생성기가 서비스와 해당 형식에 대한 TypeScript 정의를 생성하므로 Go 호출에 완전한 형식이 적용됩니다.

### Go와 JavaScript 간에 이벤트를 보내려면 어떻게 하나요?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

이벤트 이름은 정확히 일치해야 합니다. [이벤트 참조 문서](/guides/events-reference/)를 확인하세요.

### 애플리케이션을 디버깅하려면 어떻게 하나요?

`wails3 dev`을 실행한 후 웹에서와 똑같이 창 안을 마우스 오른쪽 버튼으로 클릭하여 브라우저 개발자 도구를 여세요. 개발 서버는 프런트엔드의 핫 리로드도 지원합니다. [디버깅](/guides/dev/debugging/)을 참조하세요.

## 빌드 및 배포

### 프로덕션용으로 빌드하려면 어떻게 하나요?

```bash
wails3 build
```

바이너리는 `bin/`에 생성됩니다. 프로덕션 빌드에는 이미 적절한 기본값(빌드 태그, `-trimpath`, 심벌 제거)이 적용되므로, 용량이 작은 바이너리를 만들기 위해 별도 플래그를 지정할 필요가 없습니다.

### 교차 컴파일할 수 있나요?

제한적으로 가능합니다. 각 플랫폼이 네이티브 웹뷰 라이브러리를 사용하므로 순수 Go 교차 컴파일 방식은 적용되지 않지만, 일반적인 사용 사례는 잘 지원됩니다:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

다른 OS에서 Linux용으로 빌드할 때는 Docker 기반 도구 체인을 사용합니다. 전체 지원 조합은 [크로스 플랫폼 빌드](/guides/build/cross-platform/)를 참조하세요.

### 설치 프로그램이나 패키지를 만들려면 어떻게 하나요?

```bash
wails3 package
```

이렇게 하면 플랫폼 네이티브 형식이 생성됩니다. [설치 프로그램 가이드](/guides/installers/)에서는 Windows의 NSIS, macOS의 `.app` 번들과 DMG, Linux 패키지를 다룹니다.

### 애플리케이션에 코드 서명하려면 어떻게 하나요?

Windows 및 macOS 서명(공증 포함)은 [서명 가이드](/guides/build/signing/)에서 단계별로 설명합니다.

## 기능

### 여러 창을 만들 수 있나요?

예. v3는 다중 창을 기본으로 지원합니다:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

[다중 창](/features/windows/multiple/)을 참조하세요.

### Wails는 시스템 트레이를 지원하나요?

예. 메뉴와 클릭 핸들러도 지원합니다:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

[시스템 트레이](/features/menus/systray/)를 참조하세요.

### 네이티브 대화 상자를 사용할 수 있나요?

예. 파일 대화 상자, 메시지 대화 상자, 질문 대화 상자는 모두 네이티브 구현을 사용합니다:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

[대화 상자](/features/dialogs/overview/)를 참조하세요.

### Wails는 자동 업데이트를 지원하나요?

예. Wails v3에는 GitHub Releases, keygen.sh, Sparkle AppCast용 교체 가능한 공급자, 암호화 서명 검증, 테마를 적용하거나 교체할 수 있는 기본 UI를 갖춘 내장 자체 업데이터(`app.Updater`)가 포함되어 있습니다. [앱 내 업데이터](/guides/updater/) 가이드와 [자동 업데이트 Wails 앱](/tutorials/04-self-update-a-wails-app/) 튜토리얼을 참조하세요.

## 문제 해결

### 무언가 작동하지 않습니다. 어디서부터 시작해야 하나요?

```bash
wails3 doctor
```

이 명령은 툴체인을 검증하고, 누락된 종속성을 설치 명령과 함께 나열하며, 모든 버그 보고서에 포함하는 것이 권장되는 버전 정보를 출력합니다.

### 빌드에 실패합니다

일반적인 해결 방법을 순서대로 살펴보세요:

1. `go mod tidy`
2. `cd frontend && npm install`(`node_modules` 누락이 가장 흔한 원인입니다)
3. CLI를 업데이트하세요: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. Linux에서는 누락된 GTK/WebKit 패키지가 있는지 `wails3 doctor`에서 확인하세요

### 바인딩이 누락되었거나 최신 상태가 아닙니다

```bash
wails3 generate bindings
```

개발 모드에서는 바인딩이 자동으로 다시 생성됩니다. `wails3 dev` 외부에서 새 서비스를 추가하거나 메서드 시그니처를 변경했다면 수동으로 다시 생성하세요.

### 이벤트가 발생하지 않습니다

Go의 `app.Event.Emit("name", ...)`와 JavaScript의 `Events.On("name", ...)`에서 이벤트 이름이 정확히 일치해야 합니다. 먼저 오타와 대소문자 차이를 확인하세요.

### 버그를 발견했습니다

[이슈를 등록](https://github.com/wailsapp/wails/issues)하고 `wails3 doctor` 출력을 포함해 주세요. 조치하기 쉬운 보고서를 작성하는 방법은 [피드백 가이드](/feedback/)에서 설명합니다.

## v2에서 마이그레이션

### v2에서 v3로 마이그레이션해야 하나요?

v3는 다중 창 지원, 더 깔끔한 서비스 기반 API, 내장 업데이터, 훨씬 유연한 빌드 시스템과 향상된 성능을 제공합니다. 새 프로젝트는 v3로 시작하는 것이 좋습니다. 기존 프로젝트의 경우 [마이그레이션 가이드](/migration/v2-to-v3/)에서 차이점을 단계별로 설명합니다.

### v2는 계속 유지 관리되나요?

예. v3가 안정 릴리스를 향해 나아가는 동안에도 v2에는 계속 수정 사항이 제공됩니다.

### v2와 v3를 나란히 실행할 수 있나요?

예. CLI는 별도의 바이너리(`wails` 및 `wails3`)이고 모듈의 임포트 경로도 다르므로, 서로 다른 메이저 버전의 프로젝트가 한 컴퓨터에서 문제없이 공존할 수 있습니다.

## 커뮤니티

### 어디서 도움을 받을 수 있나요?

- 간단한 질문과 토론은 [Discord](https://discord.gg/JDdSxwjhGf)를 이용하세요
- 긴 형식의 질문은 [GitHub Discussions](https://github.com/wailsapp/wails/discussions)를 이용하세요
- 버그는 [GitHub Issues](https://github.com/wailsapp/wails/issues)에 보고하세요

### 어떻게 기여할 수 있나요?

[기여 가이드](/contributing/)를 참조하세요. 버그 수정은 언제든 환영합니다. 새로운 기능과 공개 동작의 변경에는 [WEP(Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) 초안 PR을 사용합니다. Discord 또는 GitHub Discussions에서의 비공식 논의는 선택 사항입니다.

### 예제는 어디에서 찾을 수 있나요?

저장소에는 창, 대화 상자, 이벤트, 시스템 트레이, 서비스 등을 다루는 실행 가능한 예제가 60개 넘게 포함되어 있습니다: [v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples).

## 질문이 더 있으신가요?

[Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [토론을 시작하세요](https://github.com/wailsapp/wails/discussions).
