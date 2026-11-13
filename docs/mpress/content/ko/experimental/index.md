---
title: "실험적 기능"
description: "Wails v3에서 진행 중인 실험의 내용과 도입 이유, 피드백 제공 방법을 안내합니다."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="실험 기능 안내"}
이 섹션의 모든 기능은 말 그대로 실험입니다. 여기 있는 기능은 직접 활성화해야 하며 기본적으로 비활성화되어 있습니다. 릴리스 사이에 형태나 이름이 바뀌거나 완전히 제거될 수도 있습니다. 잦은 변경에 대응할 준비가 되어 있지 않다면 중요한 기능을 이러한 실험에 의존하지 마세요.

@end

## "실험적"이라는 말의 의미

Wails는 실험을 공개적으로 제공합니다. 실험은 조기에 여러분에게 제공할 만큼 유망하다고 판단했지만, 아직 완전히 도입하기로 확정하지 않은 아이디어입니다. 이 기능이 Wails에서 영구적으로 지원되는 기능이 될지 결정하기 전에 실제 사용 경험을 통해 배우고자 *이러한 이유로* 공개합니다.

따라서 이 섹션의 모든 기능에는 다음 사항이 적용됩니다.

- **직접 활성화해야 합니다.** 실험 기능은 `wails3`의 기본 동작을 변경하지 않습니다. 일반적으로 환경 변수나 빌드 플래그를 사용하여 의도적으로 활성화해야 하며, 비활성화한 상태에서는 기존 워크플로가 전혀 바뀌지 않습니다.
- **유지되지 않을 수도 있습니다.** 일부 실험은 안정 기능으로 전환되지만, 완전히 다른 모습으로 재작업되거나 폐기되는 실험도 있습니다. 이미 확신하는 기능만 출시하기보다는 공개적으로 시도하고 빠르게 배우는 편을 택합니다.
- **API는 확정되지 않았습니다.** 실험의 형태를 다듬는 동안 릴리스 사이에 이름, 플래그, 기본값 및 동작이 달라질 수 있습니다. 변경 사항은 릴리스 노트에 명시되지만, 안정 기능과 같은 수준의 안정성 보장은 기대하지 마세요.

## 여러분의 피드백이 필요합니다

이 부분이 중요합니다. 실험의 존속 여부는 실제 사용자의 피드백에 달려 있습니다. 이러한 기능을 사용해 보셨다면 다음 내용을 꼭 알려 주세요.

- 프로젝트에서 제대로 작동했나요? 어디에서 문제가 발생했나요?
- 더 빨랐거나, 명확했거나, 사용하기 좋았나요? 아니면 전환할 가치가 없었나요?
- 이 기능을 기본으로 사용하려면 어떤 조건이 충족되어야 하나요?

어떤 작업을 실행했고, 무엇을 예상했으며, 실제로 어떤 일이 발생했는지 구체적으로 알려 주는 피드백이 가장 유용합니다. GitHub Discussions의 **실험** 카테고리에는 각 실험을 위한 별도 스레드가 있습니다.

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="실험 관련 토론" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="사용 중인 실험의 스레드를 찾아 작동 결과와 발생한 문제 또는 부족한 점을 알려 주세요."}
@end

## 현재 진행 중인 실험

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="네이티브 YAML 빌드 시스템을 사용해 보세요. 하나의 마법사에서 프로젝트를 생성하고 구성한 다음 빌드하고 패키징하거나, Taskfile에서 단방향으로 마이그레이션하는 과정을 검토할 수 있습니다."}
@linkcard{title="CGO 없는 빌드" href="/experimental/#cgo-free-builds-purego" description="CGO_ENABLED=0로 macOS 및 Linux 앱을 빌드합니다. C 도구 체인이 필요 없으며 Docker 없이 크로스 컴파일할 수 있습니다. WAILS_PUREGO=1 또는 purego 빌드 태그로 활성화합니다."}
@linkcard{title="LLM 제어(MCP)" href="/guides/mcp-service/" description="LLM 에이전트가 실행 중인 Wails 앱을 검사하고 테스트하며 제어할 수 있게 해 주는 내장 Model Context Protocol 서버입니다. 사용자 코드가 필요 없으며 빌드 태그로 활성화합니다."}
@end

## CGO 없는 빌드(purego)

@note{type="caution" title="실험적 기능"}
CGO 없는 백엔드는 `purego` 빌드 태그(또는 `WAILS_PUREGO=1`)를 사용해야 활성화되는 선택 기능이며, 기본값이 **아닙니다**. 이를 사용하지 않으면 macOS 및 Linux 빌드는 이전과 똑같이 cgo 백엔드를 사용합니다. 동작은 릴리스 사이에 변경될 수 있습니다.

@end

### 기능 개요

일반적으로 Wails는 macOS와 Linux에서 cgo를 통해 시스템 UI 라이브러리(macOS의 AppKit 및 WKWebView, Linux의 GTK 및 WebKitGTK)와 통신합니다. CGO 없는 백엔드는 런타임에 시스템 라이브러리를 로드하는 [purego](https://github.com/ebitengine/purego)를 사용하여 순수 Go에서 동일한 호출을 수행합니다.

공개 애플리케이션 API는 변경되지 않습니다. 빌드할 때 백엔드를 선택하면 되며 애플리케이션 코드는 변경할 필요가 없습니다. 앱을 배포하기 전에 대상 시스템에서 테스트하세요.

다음과 같은 이점이 있습니다.

- **Go 컴파일에 C 도구 체인이 필요하지 않습니다.** Go 빌드에는 Xcode 명령줄 도구, `gcc` 또는 Linux 개발 헤더가 필요하지 않습니다. 프런트엔드 빌드에는 기존과 마찬가지로 JavaScript 도구가 필요합니다.
- **Docker 없이 크로스 컴파일할 수 있습니다.** 어느 OS에서든 일반 `go build`를 사용하여 macOS 또는 Linux 바이너리를 빌드할 수 있습니다.

Windows는 이미 cgo 없이 빌드되므로 영향을 받지 않습니다.

### 활성화하기 전에

이 실험이 모두 포함된 CLI와 프레임워크 버전을 사용하세요. CLI만 업데이트해도 애플리케이션의 `go.mod`에 있는 이전 프레임워크 종속성에는 이 백엔드가 추가되지 않습니다. 개발 버전을 사용해 보려면 애플리케이션 디렉터리에서 다음 명령을 실행하세요.

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

두 명령은 모두 계속 변경되는 `master` 브랜치를 추적합니다. 재현 가능한 빌드를 위해 CLI와 프레임워크에 동일한 테스트 완료 커밋을 사용하세요. Go 도구 체인 및 CLI 설정 방법은 [설치 가이드](/getting-started/installation/)를 따르세요.

실험적 macOS 및 Linux 백엔드는 `amd64` 및 `arm64`을 지원합니다. 크로스 컴파일은 요청한 대상을 위한 실행 파일을 생성할 뿐, 빌드 머신에서 그 파일을 실행하지는 않습니다. 대상 운영 체제와 아키텍처에서 테스트하세요.

### 사용 방법

다음 방법 중 하나를 선택하여 백엔드를 지정하세요.

#### Wails CLI로 사용하기

이 버전의 CLI로 생성한 프로젝트에서는 `WAILS_PUREGO`을 활성화하면 `build`, `package` 및 `dev`가 실험적 백엔드를 사용합니다. POSIX 셸에서는 다음과 같이 실행합니다.

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

PowerShell에서는 다음과 같이 실행합니다.

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

CLI는 대소문자와 앞뒤 공백을 무시하고 `1`, `true`, `on` 및 `yes`을 인식합니다. 빈 값으로는 이 실험이 활성화되지 않습니다.

또는 단일 빌드나 패키징 명령에 태그를 전달하세요.

```sh
wails3 build -tags purego
wails3 package -tags purego
```

태그를 설정하면 업데이트된 기본 Taskfile이 호스트의 Go 도구 체인을 사용해 `CGO_ENABLED=0`로 빌드합니다. 예를 들어 macOS 또는 Windows에서 Linux ARM64 실행 파일을 빌드하려면 다음과 같이 실행합니다.

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

이전 프로젝트는 기존 Taskfile을 그대로 유지합니다. 사용자 정의 빌드 작업을 업데이트된 템플릿과 비교하거나, 해당 작업을 변경하지 않고 백엔드를 테스트하려면 일반 Go를 사용하세요. 환경 변수는 Wails CLI 명령에서 태그를 선택하지만, 일반 `go build`에는 태그를 명시적으로 지정해야 합니다.

#### 일반 Go로 사용하기

먼저 프런트엔드를 빌드하여 `go:embed`에서 참조하는 모든 에셋이 존재하도록 한 다음 애플리케이션을 컴파일하세요.

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego`도 사용할 수 있습니다. 백엔드는 동일하며, 프로그램의 나머지 부분에서 cgo를 계속 사용할 수 있을 뿐입니다.

크로스 빌드할 때는 `GOOS` 및 `GOARCH`도 설정하세요.

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

Linux에서는 `-tags purego,gtk3,production`로 GTK3 스택을 선택하세요.

기본 백엔드로 돌아가려면 `purego` 태그를 제거하고 `WAILS_PUREGO`의 설정을 해제하세요(POSIX 셸에서는 `unset WAILS_PUREGO`, PowerShell에서는 `Remove-Item Env:WAILS_PUREGO`). macOS 또는 Linux의 cgo 백엔드로 돌아갈 때는 명시적인 `CGO_ENABLED=0` 설정도 제거하세요.

`CGO_ENABLED=0`로 빌드하면서 태그를 빠뜨리면, 해결 방법을 알려 주는 오류 하나와 함께 빌드가 중단됩니다.

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### 런타임 요구 사항

앱은 시작할 때 시스템 라이브러리를 로드하므로 cgo 빌드와 마찬가지로 대상 머신에 해당 라이브러리가 설치되어 있어야 합니다.

| 플랫폼 | 라이브러리 | 지원되는 최소 버전 |
|---|---|---|
| macOS | AppKit, WebKit(OS에 포함) | cgo 백엔드와 동일 |
| Linux(기본값) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40(Debian 12) |
| Linux(`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36(Ubuntu 22.04) |

이 최소 버전보다 최신인 함수는 앱이 시작될 때 조회하며, 설치된 라이브러리에서 해당 함수를 제공할 때만 사용합니다. 따라서 하나의 바이너리를 구형 시스템과 최신 시스템 모두에서 실행할 수 있습니다.

앱은 여전히 네이티브 UI 라이브러리를 사용하며 Linux의 libc 같은 시스템 라이브러리에 동적으로 링크될 수 있습니다. `CGO_ENABLED=0`이 완전한 정적 실행 파일을 의미하지는 않습니다. 필수 런타임 라이브러리가 없으면 대상 시스템에 해당 GTK/WebKit 스택을 설치하세요. 빌드 머신에 C 헤더를 설치해도 누락된 런타임 라이브러리 문제는 해결되지 않습니다.

### 패키징 및 서명

이 백엔드를 사용하면 컴파일 시 C 툴체인이 필요하지 않습니다. 각 패키지 형식에는 고유한 요구 사항이 그대로 적용됩니다. 특히 `purego`을 선택해도 기존 AppImage 배포 도구를 크로스 플랫폼에서 사용할 수 있게 되는 것은 아닙니다. [Linux 패키징](/guides/build/linux/)을 참조하세요.

크로스 빌드한 macOS 애플리케이션도 배포 전에 적절한 서명과 공증을 거쳐야 합니다. Windows 설치 프로그램에도 해당 설치 프로그램 도구가 필요합니다. [macOS 패키징](/guides/build/macos/) 및 [Windows 패키징](/guides/build/windows/)을 참조하세요. cgo를 사용하는 종속 항목의 경우 `CGO_ENABLED=1`에도 C 컴파일러와 대상용 헤더 또는 라이브러리가 필요합니다. Wails의 pure-Go 백엔드를 선택해도 이러한 종속 항목의 요구 사항은 없어지지 않습니다.

### 검증

생성된 바인딩 검사와 `puregocheck`은 네이티브 심볼, 선택적 버전 가드 및 백엔드 적용 범위를 검사합니다. Linux 동등성 워크플로는 두 백엔드 모두에서 GTK4와 GTK3를 테스트하며, 여기에는 에셋 요청, 바이너리 요청 본문 및 네이티브 디스패치가 포함됩니다. 또한 Linux ARM64 대상을 크로스 컴파일합니다. 크로스 컴파일 검사는 코드가 빌드되는지를 확인하고, 런타임 검사는 테스트 머신에 설치된 네이티브 라이브러리를 실제로 사용해 검사합니다.

[백엔드 테스트 가이드](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md)에서는 동등성 프로브와 엔드투엔드 입력 테스트를 설명합니다. 문제를 보고할 때는 컴파일 실패와 런타임 또는 패키징 실패를 구분해 주세요.

### 알려진 차이점

- **Linux의 메뉴 단축키는 애플리케이션 수준에서 등록됩니다.** 서로 다른 메뉴에 충돌하는 단축키를 지정하지 마세요. `SetMenu(nil)`을 사용해 창의 메뉴를 명시적으로 제거하면, 제거된 메뉴의 단축키도 해당 창에서 작동하지 않습니다. 다른 창의 메뉴와 단축키는 그대로 유지됩니다.
- **macOS 열기 및 저장 패널**은 cgo를 사용할 때와 마찬가지로 별도의 시스템 프로세스에서 실행됩니다. 테스트에서는 패널이 정상적으로 열리고 취소되는지 확인합니다.

그 밖의 문제를 발견하셨나요? 알려주세요.

### 피드백

[실험](https://github.com/wailsapp/wails/discussions/categories/experiments) 아래의 **CGO 없는 빌드** 스레드에서 사용 경험을 알려주세요. 사용한 OS와 버전, 빌드 방법, 작동한 부분과 작동하지 않은 부분을 공유해 주세요. 실패 사례와 `wails3 doctor`의 출력을 함께 제공하는 것이 가장 유용합니다.
