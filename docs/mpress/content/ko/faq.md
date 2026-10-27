---
title: "자주 묻는 질문"
description: "Wails로 애플리케이션을 빌드할 때 가장 자주 묻는 질문에 대한 답변 v3"
slug: "faq"
sourcePath: "faq.md"
---

## 일반

### Wails란 무엇인가요?

Wails는 Go와 웹 기술로 데스크톱 애플리케이션을 빌드하기 위한 프레임워크입니다. 애플리케이션 로직은 Go로 작성하고, 인터페이스는 HTML, CSS, JavaScript(또는 원하는 프런트엔드 프레임워크)로 빌드하며, Wails는 이를 운영 체제의 네이티브 웹뷰에 렌더링합니다. 그 결과 브라우저를 번들로 포함하지 않고 메모리 사용량이 적으며, 일반적으로 약 10MB인 단일 바이너리로 구성된 작고 빠르며 네이티브처럼 느껴지는 애플리케이션을 만들 수 있습니다.

### Wails는 어떤 플랫폼을 지원하나요?

| 플랫폼 | 요구 사항 |
| --- | --- |
| Windows | AMD64 및 ARM64. [WebView2 런타임](https://developer.microsoft.com/microsoft-edge/webview2/)을 사용합니다. |
| macOS | Intel에서는 10.15 이상(애플리케이션은 10.13 이상을 대상으로 지정 가능), Apple Silicon에서는 11.0 이상이 필요합니다. 유니버설 바이너리를 지원합니다. |
| Linux | AMD64 및 ARM64. 기본 스택은 WebKitGTK 6.0을 사용하는 GTK4입니다(Ubuntu 24.04 이상, Debian 13 이상, Fedora 40 이상 및 유사 배포판). Ubuntu 22.04, Debian 12, RHEL 9처럼 WebKit2GTK 4.1만 제공하는 배포판은 레거시 `-tags gtk3` 빌드를 통해 지원됩니다(v3.1까지 제공). WebKit2GTK 4.0만 제공하는 배포판은 지원되지 않습니다. [Linux 빌드 가이드](/guides/build/linux/)를 참조하세요. |
| iOS 및 Android | 실험적 지원입니다. [모바일 가이드](/guides/mobile/)를 참조하세요. |

[서버 빌드](/guides/server-build/)를 사용하여 애플리케이션을 일반 웹 앱으로 제공할 수도 있습니다.

언제든 `wails3 doctor`을 실행하여 시스템을 점검하고 플랫폼별 설치 안내를 확인할 수 있습니다.

### 시작하려면 무엇이 필요한가요?

- Go 1.26 이상
- Node.js 및 npm(프런트엔드 빌드용)
- 플랫폼 도구 체인: Windows에서는 WebView2(10/11에 사전 설치됨), macOS에서는 Xcode Command Line Tools, Linux에서는 `gcc`와 GTK/WebKit 개발 패키지

`wails3 doctor`이 이 모든 항목을 점검하고 무엇이 누락되었는지 정확히 알려 줍니다. 전체 절차는 [설치](/quick-start/installation/)를 참조하세요.

### Wails v3은 프로덕션에 사용할 준비가 되었나요?

Wails v3.0.0은 Windows, macOS 및 Linux용 데스크톱 애플리케이션을 위한 안정화 릴리스입니다. Android 및 iOS 지원은 여전히 실험적이며 안정적인 데스크톱 호환성 보장 범위에 포함되지 않습니다. [프로젝트 상태 페이지](/status/)를 참조하세요.

## 개발

### Go를 알아야 하나요?

기본적인 Go 지식이 있으면 도움이 되지만 전문가일 필요는 없습니다. 애플리케이션 로직은 일반 Go 메서드로 작성하며, 나머지는 [튜토리얼](/tutorials/overview/)에서 모두 단계별로 안내합니다. 많은 개발자가 첫 Wails 애플리케이션을 빌드하면서 Go를 익힙니다.

### 선호하는 프런트엔드 프레임워크를 사용할 수 있나요?

예. HTML, CSS 및 JavaScript로 빌드할 수 있다면 Wails에서 사용할 수 있습니다. React, Vue, Svelte 및 순수 JavaScript용 템플릿이 제공되며, 각각 TypeScript 변형도 포함됩니다. 다른 프레임워크도 몇 분이면 연동할 수 있습니다. [프런트엔드 프레임워크](/guides/dev/frontend-frameworks/)를 참조하세요.

### JavaScript에서 Go 함수를 호출하려면 어떻게 하나요?

서비스를 등록하면 Wails가 형식이 지정된 바인딩을 생성합니다:

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

바인딩은 `wails3 dev` 실행 중 자동으로 다시 생성되며, `wails3 generate bindings`을 사용하여 필요할 때 생성할 수도 있습니다. [서비스](/features/bindings/services/)를 참조하세요.

### TypeScript를 사용할 수 있나요?

예. 바인딩 생성기가 서비스와 해당 형식의 TypeScript 정의를 생성하므로 Go 호출에 완전한 형식 정보가 적용됩니다.

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

이벤트 이름은 정확히 일치해야 합니다. [이벤트 레퍼런스](/guides/events-reference/)를 참조하세요.

### 애플리케이션을 디버깅하려면 어떻게 하나요?

`wails3 dev`을 실행한 후 웹에서와 마찬가지로 창 안을 마우스 오른쪽 버튼으로 클릭하여 브라우저 개발자 도구를 여세요. 개발 서버는 프런트엔드 핫 리로드도 지원합니다. [디버깅](/guides/dev/debugging/)을 참조하세요.

## 빌드 및 배포

### 프로덕션용으로 빌드하려면 어떻게 하나요?

```bash
wails3 build
```

바이너리는 `bin/`에 생성됩니다. 프로덕션 빌드에는 이미 적절한 기본값(빌드 태그, `-trimpath`, 심볼 제거)이 적용되므로 간결한 바이너리를 만들기 위해 추가 플래그를 지정할 필요가 없습니다.

### 교차 컴파일할 수 있나요?

제한적으로 가능합니다. 각 플랫폼이 네이티브 웹뷰 라이브러리를 사용하므로 순수 Go 교차 컴파일 방식은 적용되지 않지만, 일반적인 사용 사례는 잘 지원됩니다:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

다른 OS에서 Linux용으로 빌드할 때는 Docker 기반 도구 체인을 사용합니다. 전체 지원 조합은 [크로스 플랫폼 빌드](/guides/build/cross-platform/)를 참조하세요.

### 인스톨러 또는 패키지를 만들려면 어떻게 하나요?

```bash
wails3 package
```

이렇게 하면 플랫폼 네이티브 형식이 생성됩니다. [인스톨러 가이드](/guides/installers/)에서는 Windows의 NSIS, macOS의 `.app` 번들과 DMG, Linux 패키지를 다룹니다.

### 애플리케이션에 코드 서명하려면 어떻게 하나요?

Windows 및 macOS 서명(공증 포함) 절차는 모두 [서명 가이드](/guides/build/signing/)에서 단계별로 설명합니다.

## 기능

### 여러 창을 만들 수 있나요?

예. v3은 다중 창을 네이티브로 지원합니다:

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

예. 파일 대화 상자, 메시지 대화 상자 및 질문 대화 상자는 모두 네이티브 구현을 사용합니다:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

[대화 상자](/features/dialogs/overview/)를 참조하세요.

### Wails는 자동 업데이트를 지원하나요?

예. Wails v3에는 GitHub Releases, keygen.sh 및 Sparkle AppCast용 플러그형 공급자, 암호화 서명 검증, 테마를 적용하거나 교체할 수 있는 기본 UI를 갖춘 내장 자체 업데이트 기능(`app.Updater`)이 포함되어 있습니다. [인앱 업데이터](/guides/updater/) 가이드와 [자동 업데이트 Wails 앱](/tutorials/04-self-update-a-wails-app/) 튜토리얼을 참조하세요.

## 문제 해결

### 무언가 작동하지 않습니다. 어디서부터 확인해야 하나요?

```bash
wails3 doctor
```

이 명령은 도구 체인을 검증하고, 누락된 종속성을 설치 명령과 함께 나열하며, 버그 보고서에 포함해야 할 버전 정보를 출력합니다.

### 빌드에 실패합니다

일반적인 해결 방법을 순서대로 적용하세요:

1. `go mod tidy`
2. `cd frontend && npm install`(`node_modules` 누락이 가장 흔한 원인입니다)
3. CLI 업데이트: `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. Linux에서는 `wails3 doctor`을 확인하여 누락된 GTK/WebKit 패키지가 있는지 점검하세요

### 바인딩이 없거나 최신 상태가 아닙니다

```bash
wails3 generate bindings
```

개발 모드에서는 바인딩이 자동으로 다시 생성됩니다. `wails3 dev` 외부에서 새 서비스를 추가하거나 메서드 시그니처를 변경한 경우에는 바인딩을 수동으로 다시 생성하세요.

### 이벤트가 발생하지 않습니다

이벤트 이름은 Go의 `app.Event.Emit("name", ...)`와 JavaScript의 `Events.On("name", ...)`에서 정확히 일치해야 합니다. 먼저 오타와 대소문자 차이를 확인하세요.

### 버그를 발견했습니다

[이슈를 등록](https://github.com/wailsapp/wails/issues)하고 `wails3 doctor` 출력을 포함해 주세요. [피드백 가이드](/feedback/)에서는 조치하기 쉬운 보고서를 작성하는 방법을 설명합니다.

## v2에서 마이그레이션하기

### v2에서 v3(으)로 마이그레이션해야 하나요?

v3에서는 다중 창 지원, 더 깔끔한 서비스 기반 API, 내장 업데이터, 훨씬 유연한 빌드 시스템과 향상된 성능을 제공합니다. 새 프로젝트는 v3에서 시작하는 것이 좋습니다. 기존 프로젝트의 경우 [마이그레이션 가이드](/migration/v2-to-v3/)에서 버전 간 차이를 단계별로 설명합니다.

### v2는 계속 유지보수되나요?

예. Wails v2에는 계속해서 중요 업데이트가 제공됩니다. 새 프로젝트에서는 Wails v3.0.0을 사용해야 합니다.

### v2와 v3을 나란히 실행할 수 있나요?

예. 각 CLI는 별도의 바이너리(`wails` 및 `wails3`)이고 모듈의 import 경로도 서로 다르므로, 메이저 버전이 다른 프로젝트를 한 컴퓨터에서 문제없이 함께 사용할 수 있습니다.

## 커뮤니티

### 어디에서 도움을 받을 수 있나요?

- 간단한 질문과 토론은 [Discord](https://discord.gg/JDdSxwjhGf)에서 할 수 있습니다
- 긴 형식의 질문은 [GitHub Discussions](https://github.com/wailsapp/wails/discussions)에서 할 수 있습니다
- 버그는 [GitHub Issues](https://github.com/wailsapp/wails/issues)에 보고할 수 있습니다

### 어떻게 기여할 수 있나요?

[기여 가이드](/contributing/)를 참고하세요. 버그 수정은 언제든 환영합니다. 새로운 기능과 공개 동작 변경은 [WEP(Wails 개선 제안)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) 초안 PR을 통해 진행합니다. Discord 또는 GitHub Discussions에서의 비공식 논의는 선택 사항입니다.

### 예제는 어디에서 찾을 수 있나요?

저장소에는 창, 대화 상자, 이벤트, 시스템 트레이, 서비스 등을 다루는 실행 가능한 예제가 60개 넘게 포함되어 있습니다: [v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).

## 아직 궁금한 점이 있나요?

[Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [토론을 시작](https://github.com/wailsapp/wails/discussions)하세요.
