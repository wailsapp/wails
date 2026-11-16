---
title: "YAML 실행 구성"
description: "wails3 run의 로컬 실행 기본값, 플랫폼별 재정의 및 실행 동작입니다."
---

## 가용성

이 참조 문서는 실험적인 run 확장을 설명합니다. 출시된 CLI에서 사용할 수 있음을 보장하지 않습니다.

`wails3 run`은 로컬 매니페스트 기본값을 사용하여 애플리케이션을 빌드하고 실행합니다. 매니페스트 기반 실행을 사용하려면 [`WAILS_EXP_USE_WAKE=1`에 명시적으로 동의해야 합니다](/experimental/yaml-builds#enable-yaml-builds). 매니페스트를 찾지 못하면 이러한 동의 없이 현재 디렉터리에서 `go run .`을 실행합니다. 유효하지 않거나 읽을 수 없는 매니페스트는 오류이며, 이 경우 대체 동작을 선택하지 않습니다. YAML이 비활성화된 상태에서 매니페스트가 존재하면 명령이 YAML 활성화 방법을 안내합니다.

생성된 [`wails.yaml` 필드 참조](/reference/wails-yaml-fields)에는 필드 형식과 기본값이 나열되어 있습니다. 릴리스에서 네이티브 모바일 및 macOS 번들을 지원한다고 명시하기 전에 해당 호스트에서 번들이 정상적으로 허용되는지 확인해야 합니다.

## 명령 계약

`wails3 run`은 가장 가까운 `wails.yaml`을 찾고, 대상 하나의 구성을 결정한 다음, 기존 캐시 파이프라인을 통해 빌드하고 결과물을 실행합니다. 검색은 작업 디렉터리에서 상위로 진행되며, 각 위치에서 매니페스트를 확인한 후 `go.mod` 경계에 도달하면 중단됩니다. 매니페스트 디렉터리가 프로젝트 루트이자 데스크톱 애플리케이션의 작업 디렉터리입니다.

기본 대상은 호스트 운영 체제와 아키텍처입니다. 모바일 대상은 명시적으로 선택해야 합니다. 크로스 컴파일에 성공하더라도 데스크톱 실행 파일을 호스트에서 실행할 수 있다는 의미는 아닙니다. 호환되지 않는 데스크톱 대상은 빌드 전에 실패합니다. 여기에는 `CGO_ENABLED=0`으로 빌드한 실행 파일도 포함됩니다. Go 전용 크로스 빌드가 성공해도 에뮬레이터가 제공되거나 로컬 실행 계약이 변경되지는 않습니다. 크로스 컴파일에는 [`wails3 build`을 사용하고](/experimental/yaml-builds/#cross-compile-desktop-applications), 생성된 파일은 대상 시스템에서 실행하거나 설치하십시오. Android에는 `--target android/arm64 --device <adb-serial>` 또는 `--target android/arm64 --emulator <avd-name>`를 사용하십시오. 선택한 아키텍처는 기기 또는 에뮬레이터와 일치해야 합니다. iOS 시뮬레이터에는 `--target ios/arm64 --device <udid>`을 사용하고, 실제 기기에는 `--destination device`을 추가하십시오. iOS 계획 수립 및 실행에는 macOS와 Xcode가 필요합니다.

YAML을 사용하면 임베디드 자산과 프로덕션 빌드 경로를 사용하여 빌드하며, 이 경로는 `production` 태그를 자동으로 제공합니다. 라이브 리로드는 계속 `wails3 dev`에서 담당하며, `dev.tags`는 `wails3 run`에 영향을 주지 않습니다.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan`은 빌드, 설치 또는 실행 없이 선택한 대상, 유효한 빌드 태그, 실행 또는 설치 가능한 산출물, 인수 및 작업 디렉터리를 결정하여 표시합니다. 계획에는 환경 재정의 키 이름만 나열되고 값은 생략됩니다.

## CLI 옵션

| 옵션 | 의미 |
| --- | --- |
| `--tags <tags>` | 쉼표로 구분한 추가 Go 빌드 태그입니다. |
| `--skip-build` | YAML 출력 디렉터리에서 기존 데스크톱 실행 파일을 가져온 후, 컴파일하지 않고 실행합니다. |
| `--target <os>/<arch>` | 대상 하나를 선택합니다. 기본값은 호스트 OS와 아키텍처입니다. |
| `--plan` | 빌드하거나 실행하지 않고 결정된 설정을 표시합니다. 크로스 플랫폼 데스크톱 계획도 허용됩니다. |
| `--device <id>` | Android adb 일련번호 또는 iOS 시뮬레이터/기기 식별자입니다. |
| `--emulator <name>` | 시작하거나 재사용할 Android Virtual Device입니다. `--device`과 함께 사용할 수 없습니다. |
| `--destination <kind>` | iOS 대상: `simulator`(기본값) 또는 `device`. |
| `-- <args...>` | 구성된 애플리케이션 인수를 대체합니다. 값이 없는 `--`은 인수를 지웁니다. |
| `--appargs <string>` | 호환성 표기법입니다. POSIX 스타일 인용 규칙으로 문자열 하나를 구문 분석합니다. 반복하거나 `--` 뒤의 인수와 함께 사용할 수 없습니다. |

모바일 선택기를 사용하려면 매니페스트와 일치하는 모바일 대상이 필요합니다.

`--skip-build`을 사용하려면 활성 YAML 프로젝트와 호환되는 데스크톱 실행 파일이 필요합니다. 이 모드에서는 프런트엔드 작업, 바인딩, 컴파일, UPX 및 빌드 훅을 건너뛰지만, 구성된 실행 인수와 환경은 계속 적용됩니다. 가져온 실행 파일은 유효성 검사를 거친 뒤 원본 파일을 수정하지 않고 Wake에서 생성한 작업 공간에 복사됩니다. macOS GUI 애플리케이션은 이러한 입력 파일과 현재 패키지 자산을 사용하여 구성됩니다. [미리 빌드된 입력 경로 및 요구 사항](/experimental/yaml-builds/#use-prebuilt-desktop-executables)을 참조하세요.
빌드 태그로는 미리 빌드된 실행 파일을 변경할 수 없으므로, 이 모드에서는 `--tags` 및 구성된 빌드 태그가 적용되지 않습니다.

## 필드

선택 사항인 `run` 매핑은 공통 실행 기본값을 제공합니다. `targets["<os>"]` 또는 `targets["<os>/<arch>"]` 내부의 선택 사항인 `run` 매핑은 플랫폼 또는 아키텍처별 기본값을 제공합니다. 맵 순서와 관계없이 아키텍처별 설정이 OS 전체 설정보다 우선합니다.

| 필드 | 형식 | 기본값 | 의미 |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | run 빌드에서만 사용하는 추가 Go 태그입니다. |
| `run.args` | `list(string)` | `[]` | 애플리케이션에 전달할 인수입니다. |
| `run.environment` | `map(string)` | `{}` | 컴파일러 환경과 별개인 런타임 환경 재정의입니다. |
| `target["<os>"].run.tags` | `list(string)` | `[]` | 선택한 플랫폼의 추가 run 태그입니다. |
| `target["<os>"].run.args` | `list(string)` | 상속됨 | 명시적으로 설정하면 공통 인수를 대체하며, `[]`은 인수를 지웁니다. |
| `target["<os>"].run.environment` | `map(string)` | `{}` | 공통 런타임 환경 값을 키별로 재정의합니다. |

인수는 셸 확장 없이 인수 벡터로 전달됩니다. 환경 값은 기존 YAML 리터럴 전용 계약에 따라 리터럴 문자열로 처리됩니다.

## 플랫폼 예제

이 조각은 일반 프로젝트 구성에 추가하여 사용합니다. `example_private_api`은 예제에 필요한 실제 비공개 API 태그의 자리표시자입니다.

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Linux 환경 설정은 이 예제에만 해당하며 Wayland 세션이 필요합니다. Wails 기본값은 아닙니다. 비공개 API 태그는 macOS에서만 선택됩니다. 기존 `target.tags`은 run 빌드를 포함한 프로덕션 빌드에 적용됩니다. `target.run.tags`은 추가 태그를 run 빌드로 한정합니다.

아키텍처별 재정의를 사용하여 OS 전체 기본값을 세분화할 수 있습니다.

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

앞의 Linux 매핑과 함께 사용하면 `GDK_BACKEND`은 유지되고 아키텍처별 태그가 추가되며 공통 애플리케이션 인수는 지워집니다.

## 결정 규칙

| 설정 | 결정 순서(가장 낮은 우선순위부터 가장 높은 우선순위까지) |
| --- | --- |
| 태그 | `build.tags`, 선택한 `target.tags`, 자동 `production`, `run.tags`, 선택한 `target.run.tags` 및 CLI `--tags`을 결합하고, 각 태그가 처음 나타난 위치만 유지합니다. |
| 런타임 환경 | 실행 셸의 환경을 상속하고 `run.environment`을 덮어쓴 다음 선택한 `target.run.environment`을 덮어씁니다. 빈 값은 빈 문자열로 설정됩니다. |
| 애플리케이션 인수 | 공통 `run.args`을 명시적으로 설정된 대상 인수로 대체한 다음, `--` 뒤에 명시한 인수로 다시 대체합니다. 값이 없는 `--`은 구성된 인수를 지웁니다. |

선택한 대상 내에서는 OS 전체 설정을 아키텍처별 설정보다 먼저 결정합니다. 태그는 누적되고, 환경 값은 키별로 병합되며, 명시적으로 제공한 인수는 덜 구체적인 목록을 대체합니다. 또한 run 빌드는 해당 모바일 대상에 `android` 또는 `ios`을 제공합니다.

바인딩 생성과 Go 컴파일에는 동일하게 결정된 태그가 전달됩니다. 빌드 캐시 키에는 이 태그들이 포함됩니다. 애플리케이션 인수나 런타임 환경만 변경해도 컴파일 캐시는 무효화되지 않으며, 다음 실행에서 새 값을 사용합니다.

플랫폼 매핑은 특정 플랫폼을 사용자 지정할 뿐, 매핑이 있다고 해서 다른 플랫폼이 지원되지 않는 것은 아닙니다. `project.supported_platforms`은 빌드와 `run`에서 허용되는 플랫폼을 제한합니다. 목록을 생략하거나 비워 두면 모든 플랫폼을 허용합니다. 예를 들어 `project` 내부에 `supported_platforms: [darwin]`을 지정하면 macOS 전용 예제가 Linux에서 조기에 실패합니다.

## 플랫폼 실행 동작

| 플랫폼 | 실행 동작 |
| --- | --- |
| Linux 및 Windows | 표준 입출력을 연결하고 종료 처리를 수행하며 애플리케이션 종료 상태를 전파하도록 실행 파일을 시작합니다. |
| macOS | 앱 번들을 구성하고 표준 입출력을 연결하여 번들의 `Contents/MacOS` 바이너리를 실행합니다. URL 스킴과 파일 연결을 위해 Launch Services에 번들을 등록합니다. 서명과 네이티브 수명 주기 승인에는 macOS에서의 검증이 필요합니다. |
| Android | APK를 빌드하고 기기 또는 에뮬레이터를 선택한 다음 애플리케이션을 설치하고 실행합니다. 빌드, 설치 및 실행 실패를 각각 구분하여 보고합니다. |
| iOS | 앱 번들을 빌드하고 시뮬레이터 또는 기기를 선택한 다음 서명 및 프로비저닝 요구 사항을 충족하여 설치하고 실행합니다. 기기용 빌드에는 적절한 자격 증명과 Apple 도구 체인이 필요합니다. |

Android에서는 배포 전에 애플리케이션 인수와 런타임 환경 재정의를 거부합니다. Android도 지원하는 경우 데스크톱 전용 공유 설정은 데스크톱 대상 매핑에 배치해야 합니다. Android 실행 시 애플리케이션 로그를 연결하며, 연결이 끝나면 애플리케이션을 중지합니다. iOS는 선택한 런처에 인수를 전달합니다. 시뮬레이터에는 `SIMCTL_CHILD_` 변수를 통해 환경 재정의를 전달하고, 실제 기기에는 `devicectl --environment-variables`을 통해 JSON 맵을 전달합니다. 모바일 종료 상태는 이식 가능한 애플리케이션 종료 코드가 아니라 런처 또는 로그 모니터링 결과를 나타냅니다.

## Go 전용 예제 및 대체 동작

기존 또는 임베드된 애셋을 사용하는 Go 전용 예제에서는 프런트엔드 단계를 비활성화합니다.

```yaml
frontend:
  disabled: true
```

이렇게 하면 프런트엔드 종속성 설치, 바인딩 생성 및 프런트엔드 빌드를 건너뜁니다. 생성된 바인딩 클래스가 필요한 예제에서는 대신 `frontend.bindings` 아래에 `interfaces: false`을 설정할 수 있습니다.

매니페스트가 없을 때의 대체 동작은 `production` 또는 `devtools`을 추가하지 않고 `go run [-tags <tags>] . [application arguments]`을 실행합니다. 종료 상태에는 Go 자체 동작을 따릅니다. `--` 뒤의 인수는 `--help` 같은 플래그를 포함하여 애플리케이션에 전달됩니다. `--plan`에 표시되는 런타임 값은 환경 키 이름으로 제한되며 값 자체는 출력되지 않습니다. 대체 계획에는 Go 명령과 작업 디렉터리가 표시됩니다.

## 개발용 애플리케이션 인수

활성 YAML 프로젝트에서 `wails3 dev`은 프로덕션 `run` 기본값과 별개로 애플리케이션 인수를 받습니다. `dev.args`은 공유 기본값을 제공합니다. `targets["<os>"]` 또는 `targets["<os>/<arch>"]`에 중첩된 `dev` 매핑은 해당 대상에서 이 기본값을 대체합니다. 맵 순서와 관계없이 아키텍처 설정이 OS 전체 설정보다 우선합니다. `args: []`을 명시하면 상속된 인수가 지워집니다.

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args`과 `run.tags`은 개발 실행에 사용되지 않습니다. 기존 `dev.tags` 빌드 정책은 변경되지 않습니다. 애플리케이션은 매니페스트 디렉터리를 기준으로 상대 애플리케이션 경로를 해석합니다.

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

`--` 뒤의 인수는 목록이 비어 있는 경우를 포함하여 구성된 인수를 대체합니다. 이 구분자 뒤의 `--help`은 애플리케이션에 속합니다. `--appargs=""`도 기본값을 지웁니다. 호환성 옵션은 따옴표로 묶인 문자열 하나를 받습니다. 이 옵션을 `--`, 위치 기반 애플리케이션 인수 또는 두 번째 `--appargs`과 함께 사용하면 오류가 발생합니다. 따옴표 구문이 잘못된 경우 세션을 시작하기 전에 거부합니다.

호환성 파서는 모든 호스트에서 동일한 POSIX 스타일 따옴표 규칙을 사용합니다. 셸 명령을 실행하거나 변수를 확장하지 않습니다. 다만 호출한 셸에서 Wails에 텍스트를 전달하기 전에 이미 확장했을 수 있습니다. `--` 뒤의 인수는 이미 개별 값이므로 추가로 토큰화하지 않습니다. 예를 들어 PowerShell에서 Windows 경로를 지정할 때는 이 형식을 사용하는 것이 좋습니다.

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

세션은 소스가 다시 빌드되어도 확정된 인수 벡터를 유지합니다. `dev.args` 또는 선택한 대상의 인수만 변경하면 캐시된 실행 파일을 사용하여 백엔드를 다시 시작하며 프런트엔드는 다시 시작하지 않습니다. CLI 재정의는 매니페스트를 다시 로드한 후에도 계속 적용됩니다. 교체에 실패하면 이전 인수를 사용하는 이전 실행 파일을 유지하거나 복원합니다. 애플리케이션 인수는 컴파일러 플래그가 아니며 컴파일 캐시 키에 포함되지 않습니다. `dev --plan`은 실행하지 않고 유한한 빌드 계획과 함께 실제 적용되는 인수를 표시합니다.

Linux, Windows 및 macOS는 인수 벡터를 실행 파일에 직접 전달합니다. iOS는 시뮬레이터 또는 기기 런처를 통해 전달합니다. Android에는 이에 상응하는 지원 인수 채널이 없으므로, 기기를 검색하거나 빌드하기 전에 비어 있지 않은 개발용 애플리케이션 인수를 거부합니다. 데스크톱 전용 기본값은 데스크톱 대상 매핑에 배치하거나 `targets.android.dev` 아래의 `args: []`을 사용하여 지울 수 있습니다.

이러한 개발 옵션을 사용하려면 활성 `wails.yaml`과 YAML 옵트인이 필요합니다. 이 옵션은 레거시 개발 경로의 Taskfile 또는 감시자 셸 명령을 다시 작성하지 않습니다.
