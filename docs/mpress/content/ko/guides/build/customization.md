---
title: "빌드 사용자 지정"
description: "Task와 Taskfile.yml을 사용하여 빌드 프로세스 사용자 지정"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 개요

Wails 빌드 시스템은 Wails 애플리케이션의 빌드 프로세스를 간소화하도록 설계된 유연하고 강력한 도구입니다. 작업을 쉽게 정의하고 실행할 수 있는 태스크 러너인 [Task](https://taskfile.dev)를 활용합니다. v3 빌드 시스템이 기본값이지만, Wails는 개발자가 필요에 따라 빌드 프로세스를 사용자 지정할 수 있도록 "원하는 도구를 직접 사용하는" 방식을 권장합니다.

Task 사용 방법에 대한 자세한 내용은 [공식 문서](https://taskfile.dev/usage/)를 참조하세요.

## Task: 빌드 시스템의 핵심

[Task](https://taskfile.dev)는 Go로 작성된 현대적인 Make 대체 도구입니다. YAML 파일을 사용하여 작업과 작업 간 종속성을 정의합니다. Wails 빌드 시스템에서는 [Task](https://taskfile.dev)가 빌드 프로세스를 조율하는 핵심 역할을 합니다.

기본 `Taskfile.yml`은 프로젝트 루트에 있으며, 플랫폼별 작업은 `build/<platform>/Taskfile.yml` 파일에 정의됩니다. `build` 디렉터리의 공통 `Taskfile.yml` 파일에는 여러 플랫폼에서 공유하는 공통 작업이 포함됩니다.

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

프로젝트 루트의 `Taskfile.yml` 파일은 빌드 시스템의 기본 진입점입니다. 이 파일에서 작업과 작업 간 종속성을 정의합니다. 기본 `Taskfile.yml` 파일은 다음과 같습니다.

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## 플랫폼별 Taskfile

각 플랫폼에는 `build` 디렉터리 아래의 해당 플랫폼 디렉터리에 자체 Taskfile이 있습니다. 이 파일들은 해당 플랫폼의 핵심 작업을 정의합니다. 각 Taskfile은 `build/Taskfile.yml` 파일의 공통 작업을 포함합니다.

### Windows

위치: `build/windows/Taskfile.yml`

Windows 전용 Taskfile에는 Windows에서 애플리케이션을 빌드하고 패키징하고 실행하는 작업이 포함됩니다. 주요 기능은 다음과 같습니다.

- 선택적 프로덕션 플래그를 사용한 빌드
- `.ico` 아이콘 파일 생성
- Windows `.syso` 파일 생성
- 패키징용 NSIS 설치 프로그램 생성

### Linux

위치: `build/linux/Taskfile.yml`

Linux 전용 Taskfile에는 Linux에서 애플리케이션을 빌드하고 패키징하고 실행하는 작업이 포함됩니다. 주요 기능은 다음과 같습니다.

- 선택적 프로덕션 플래그를 사용한 빌드
- AppImage, deb, rpm 및 Arch Linux 패키지 생성
- Linux 애플리케이션용 `.desktop` 파일 생성

### macOS

위치: `build/darwin/Taskfile.yml`

macOS 전용 Taskfile에는 macOS에서 애플리케이션을 빌드하고 패키징하고 실행하는 작업이 포함됩니다. 주요 기능은 다음과 같습니다.

- amd64, arm64 및 유니버설(두 아키텍처 모두) 아키텍처용 바이너리 빌드
- `.icns` 아이콘 파일 생성
- 배포용 `.app` 번들 생성
- `.app` 번들 임시 서명
- macOS 전용 빌드 플래그 및 환경 변수 설정

## 작업 실행 및 명령어 별칭

`wails3 task` 명령어는 [Taskfile](https://taskfile.dev)의 내장 버전으로, `Taskfile.yml`에 정의된 작업을 실행합니다.

`wails3 build` 및 `wails3 package` 명령어는 각각 `wails3 task build` 및 `wails3 task package`의 별칭입니다. 이러한 명령어를 실행하면 Wails가 내부적으로 적절한 작업 실행 명령으로 변환합니다.

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### 작업에 매개변수 전달

`KEY=VALUE` 형식을 사용하여 CLI 변수를 작업에 전달할 수 있습니다. 이러한 변수는 별칭 명령어를 통해 전달됩니다.

```bash
# These are equivalent:
wails3 build GOOS=linux
wails3 task build GOOS=linux
```

기본 제공 작업은 `GOOS`, `ARCH` 및 `DEV`를 사용합니다. `VERSION`, `OUTPUT`, `PLATFORM` 또는 `CONFIG` 같은 변수는 자체 Taskfile에서 해당 변수를 읽는 경우에만 작업에 영향을 줍니다. 다음은 사용자 지정 Taskfile의 예시이며, Wails에서 기본 제공하는 빌드 작업이 아닙니다.

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

사용자 지정 패키징 작업에서 `VERSION` 및 `OUTPUT`을 읽는 경우, 같은 방식으로 전달할 수 있습니다:

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

## 공통 빌드 프로세스

모든 플랫폼에서 빌드 프로세스에는 일반적으로 다음 단계가 포함됩니다.

1. Go 모듈 정리
2. 프런트엔드 빌드
3. 아이콘 생성
4. 플랫폼별 플래그를 사용한 Go 코드 컴파일
5. 애플리케이션 패키징(플랫폼별)

## 빌드 프로세스 사용자 지정

v3 빌드 시스템은 안정적인 기본 구성을 제공하지만, 프로젝트의 요구 사항에 맞게 쉽게 사용자 지정할 수 있습니다. `Taskfile.yml` 및 플랫폼별 Taskfile을 수정하여 다음 작업을 수행할 수 있습니다.

- 새 작업 추가
- 기존 작업 수정
- 작업 실행 순서 변경
- 다른 도구 및 스크립트와 통합

이러한 유연성 덕분에 Wails 빌드 시스템이 제공하는 구조의 이점을 유지하면서 구체적인 요구 사항에 맞게 빌드 프로세스를 조정할 수 있습니다.

@note{type="tip" title="Taskfile 학습"}
Taskfile을 효과적으로 사용하는 방법을 이해하려면 [Taskfile](https://taskfile.dev) 문서를 읽어보시기를 적극 권장합니다. `wails3 task --version`을 실행하면 Wails CLI에 내장된 Taskfile 버전을 확인할 수 있습니다.

@end

## 개발 모드

Wails 빌드 시스템에는 라이브 리로드와 핫 모듈 교체를 제공하여 개발자 경험을 개선하는 강력한 개발 모드가 포함되어 있습니다. `wails3 dev` 명령어를 사용하면 이 모드가 활성화됩니다.

### 작동 방식

`wails3 dev`을 실행하면 다음 프로세스가 진행됩니다.

1. 명령어가 사용 가능한 포트를 확인하며, 포트를 지정하지 않으면 기본값으로 9245을 사용합니다.
2. 프런트엔드 개발 서버(Vite)의 환경 변수를 설정합니다.
3. [refresh](https://github.com/atterpac/refresh) 라이브러리를 사용하여 파일 감시기를 시작합니다.

[refresh](https://github.com/atterpac/refresh) 라이브러리는 파일 변경 사항을 모니터링하고 재빌드를 트리거합니다. `./build/config.yml` 파일의 `dev_mode` 키 아래에 정의된 구성을 사용합니다. 특정 디렉터리와 파일을 무시하고, 감시할 파일과 변경이 감지되었을 때 수행할 작업을 결정하도록 구성할 수 있습니다. 기본 구성도 상당히 잘 작동하지만, 필요에 맞게 자유롭게 사용자 지정할 수 있습니다.

### 구성

구조의 예시는 다음과 같습니다.

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

이 구성 파일에서는 다음 항목을 설정할 수 있습니다.

- 파일 감시의 루트 경로 설정
- 로깅 수준 구성
- 파일 변경 이벤트의 디바운스 시간 설정
- 특정 디렉터리, 파일 또는 파일 확장자 무시
- 파일 변경 시 실행할 명령어 정의

### 개발 모드 사용자 지정

`config.yml` 파일에서 다음 값을 수정하여 개발 모드 환경을 사용자 지정할 수 있습니다.

다음과 같은 방식으로 사용자 지정할 수 있습니다.

1. 감시할 디렉터리 또는 파일 변경
2. 시스템이 변경 사항에 응답하는 속도를 제어하도록 디바운스 시간 조정
3. 프로젝트 요구 사항에 맞게 실행 명령 추가 또는 수정

### 개발에 브라우저 사용

Wails v2에서는 개발에 브라우저를 사용하는 방식을 완전히 지원했지만, 이로 인해 많은 혼란이 발생했습니다. 일부 브라우저 API는 웹뷰에서 사용할 수 없으므로, 브라우저에서 작동하는 애플리케이션이 데스크톱 애플리케이션에서도 반드시 작동하는 것은 아니었습니다.

UI 중심의 개발 작업에서는 개발 모드로 `http://localhost:9245`의 Vite URL에 접속하여 v3에서 브라우저를 사용할 수 있습니다. 따라서 스타일과 레이아웃을 작업하면서 강력한 브라우저 개발자 도구를 활용할 수 있습니다. 이 모드에서는 Go 바인딩이 *작동하지 않는다는 점*에 유의하세요. 바인딩 및 이벤트 같은 기능을 테스트할 준비가 되면 데스크톱 보기로 전환하여 프로덕션 환경에서 모든 기능이 제대로 작동하는지 확인하세요.
