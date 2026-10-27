---
title: "CLI 플래그 및 예제"
description: "개별 Wails CLI 명령의 플래그와 사용 예제를 확인합니다"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

필요한 명령을 이미 알고 있다면 이 참조 문서를 사용하세요. 용도별로 명령을 살펴보려면 [CLI 명령 개요](/reference/cli/)부터 확인하세요.

Wails CLI는 Wails 애플리케이션을 개발, 빌드 및 유지 관리하는 데 필요한 포괄적인 명령 모음을 제공합니다.

## 핵심 명령

핵심 명령은 프로젝트 생성, 개발 및 빌드에 사용하는 주요 명령입니다.

모든 CLI 명령은 `wails3 <command>` 형식을 따릅니다.

### `init`

새 Wails 프로젝트를 초기화합니다. 초기화 과정에서 `go mod tidy` 명령을 실행하여 프로젝트 패키지를 최신 상태로 업데이트합니다. `init` 명령에 `-skipgomodtidy` 플래그를 사용하면 이 단계를 건너뛸 수 있습니다.

```bash
wails3 init [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-p` | Go 패키지 이름 | `main` |
| `-t` | 템플릿 이름 또는 URL | `vanilla` |
| `-n` | 프로젝트 이름 |  |
| `-d` | 프로젝트 디렉터리 | `.` |
| `-q` | 출력 숨기기 | `false` |
| `-l` | 템플릿 목록 표시 | `false` |
| `-mod` | Go 모듈 경로(생략하면 `-git`에서 계산) |  |
| `-git` | Git 저장소 URL |  |
| `-s` | 원격 템플릿 사용 시 경고 건너뛰기 | `false` |
| `-productname` | 제품 이름 | `My Product` |
| `-productdescription` | 제품 설명 | `My Product Description` |
| `-productversion` | 제품 버전 | `0.1.0` |
| `-productcompany` | 회사 이름 | `My Company` |
| `-productcopyright` | 저작권 고지 | `© now, My Company` |
| `-productcomments` | 파일 주석 | `This is a comment` |
| `-productidentifier` | 제품 식별자 |  |
| `-skipgomodtidy` | go mod tidy 건너뛰기 | `false` |

`-git` 플래그에는 다음과 같은 다양한 Git URL 형식을 사용할 수 있습니다.

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` 또는 `ssh://git@github.com/username/project`
- Git 프로토콜: `git://github.com/username/project`
- 파일 시스템: `file:///path/to/project.git`

이 플래그를 지정하면 다음 작업을 수행합니다.

1. 프로젝트 디렉터리에 Git 저장소 초기화
2. 지정한 URL을 원격 origin으로 설정
3. 저장소 URL과 일치하도록 `go.mod`의 모듈 이름 업데이트
4. 모든 파일 추가

### `dev`

애플리케이션을 개발 모드로 실행합니다. 프런트엔드 코드를 실시간으로 확인할 수 있으며, 애플리케이션 전체를 다시 빌드하지 않고도 변경 사항이 실행 중인 애플리케이션에 반영되는 것을 확인할 수 있습니다. Go 코드 변경 사항도 감지하여 애플리케이션을 자동으로 다시 빌드하고 실행합니다.

```bash
wails3 dev [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-config` | 구성 파일 경로 | `./build/config.yml` |
| `-port` | Vite 개발 서버 포트 | `9245` |
| `-s` | HTTPS 활성화 | `false` |

@note{type="info"}
이는 `wails3 task dev` 실행과 동일하며, 프로젝트의 기본 Taskfile에 있는 `dev` 태스크를 실행합니다. `Taskfile.yml` 파일을 편집하여 이 동작을 사용자 지정할 수 있습니다.

@end

### `build`

생성된 Taskfile을 사용하여 프로덕션 버전을 빌드합니다. 핫 리로드에는 `wails3 dev`을 사용하고, 개발 빌드에는 `DEV=true`을 전달하세요. 기본적으로 현재 플랫폼과 아키텍처를 대상으로 빌드합니다.

```bash
wails3 build [flags] [CLI variables...]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-tags` | 추가 Go 빌드 태그(쉼표로 구분) |  |

CLI 변수를 전달하여 빌드를 사용자 지정할 수 있습니다.

```bash
wails3 build GOOS=linux
```

사용자 지정 Go 빌드 태그를 전달하려면 `-tags` 플래그를 사용하세요.

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

태그는 기본 Taskfile에 `EXTRA_TAGS` 형식으로 전달됩니다.

@note{type="info"}
이는 `build` 태스크를 프로젝트의 기본 Taskfile에서 실행하는 `wails3 task build` 명령과 동일합니다. `build`에 전달한 모든 CLI 변수는 기본 태스크로 전달됩니다. `Taskfile.yml` 파일을 편집하여 빌드 프로세스를 사용자 지정할 수 있습니다.

@end

### `package`

배포용 플랫폼별 패키지를 생성합니다.

```bash
wails3 package [CLI variables...]
```

기본 패키징 작업에서는 `ARCH`을 사용하여 대상 아키텍처를 선택할 수 있습니다.

```bash
wails3 package ARCH=amd64
```

`VERSION`과 `OUTPUT`은 해당 변수를 명시적으로 읽는 사용자 지정 Taskfile에만 영향을 미치며, 기본 패키지 버전이나 파일 이름을 재정의하지 않습니다.

#### 패키지 유형

각 플랫폼에서 사용할 수 있는 패키지 유형은 다음과 같습니다.

| 플랫폼 | 패키지 유형 |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
이는 프로젝트의 기본 Taskfile에 있는 `package` 태스크를 실행하는 `wails3 task package`과 같습니다. `package`에 전달된 모든 CLI 변수는 내부 태스크로 전달됩니다. `Taskfile.yml` 파일을 편집하여 패키징 프로세스를 사용자 지정할 수 있습니다.

@end

기본 빌드 태스크는 `GOOS`, `ARCH` 및 `DEV`을 읽습니다. 아래 사용자 지정 태스크의 `PLATFORM` 및 `CONFIG`은 사용자가 정의한 예입니다.

### `task`

프로젝트의 Taskfile.yml에 정의된 태스크를 실행합니다. 사용자 지정 빌드, 테스트 및 배포 태스크를 정의하고 실행할 수 있는 [Taskfile](https://taskfile.dev)의 내장 버전입니다.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI 변수

`KEY=VALUE` 형식으로 태스크에 변수를 전달할 수 있습니다.

```bash
wails3 task build GOOS=linux
wails3 task deploy ENV=staging VERSION=1.2.3
```

Taskfile.yml에서 Go 템플릿 구문을 사용하여 이러한 변수에 접근할 수 있습니다.

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-h` | Task 사용법을 표시합니다 | `false` |
| `-i` | 새 Taskfile.yml을 생성합니다 | `false` |
| `-list` | 설명이 있는 태스크를 나열합니다 | `false` |
| `-list-all` | 설명 유무와 관계없이 모든 태스크를 나열합니다 | `false` |
| `-json` | 태스크 목록을 JSON 형식으로 출력합니다 | `false` |
| `-status` | 태스크가 최신 상태가 아니면 0이 아닌 코드로 종료합니다 | `false` |
| `-f` | 태스크가 최신 상태여도 강제로 실행합니다 | `false` |
| `-w` | 지정한 태스크에 감시 모드를 활성화합니다 | `false` |
| `-v` | 상세 출력 모드를 활성화합니다 | `false` |
| `-version` | Task 버전을 출력합니다 | `false` |
| `-s` | 명령어 표시를 비활성화합니다 | `false` |
| `-p` | 태스크를 병렬로 실행합니다 | `false` |
| `-dry` | 태스크를 실행하지 않고 컴파일하여 출력합니다 | `false` |
| `-summary` | 태스크의 요약 정보를 표시합니다 | `false` |
| `-x` | 태스크의 종료 코드를 그대로 반환합니다 | `false` |
| `-dir` | 실행 디렉터리를 설정합니다 |  |
| `-taskfile` | 실행할 Taskfile을 선택합니다 |  |
| `-output` | 출력 스타일을 설정합니다: [interleaved|group|prefixed] |  |
| `-c` | 색상 출력(기본적으로 활성화됨) | `true` |
| `-C` | 동시에 실행할 태스크 수를 제한합니다 |  |
| `-interval` | 변경 사항을 감시할 간격(초) |  |

#### 예제

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build GOOS=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

에이전트 지원 프로젝트 관리를 위한 Wails 프로젝트 MCP 서버를 시작합니다. 이 서버는 실행 중인 애플리케이션에 컴파일되어 포함된 MCP 서버와 별개입니다. `wails3 mcp`은 프로젝트 파일과 수명 주기 명령을 관리하는 반면, 애플리케이션 MCP 서버는 실행 중인 WebView를 제어합니다.

```bash
wails3 mcp [flags]
```

전송 방식은 자동으로 선택됩니다.

- MCP 호스트가 파이프로 연결된 stdin/stdout을 사용하여 Wails를 실행하면 서버는 **stdio**를 사용합니다.
- 터미널에서 대화형으로 실행하면 서버는 `127.0.0.1`에서 **Streamable HTTP**를 사용하고 운영 체제에 사용 가능한 포트를 요청합니다.

전송 방식을 명시적으로 선택하려면 `--stdio` 또는 `--http`을 사용합니다. HTTP 모드에서 사용 가능한 루프백 포트를 선택하려면 `--port 0`을 사용합니다.

#### MCP 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `--root` | 허용된 프로젝트 루트입니다. 이 루트 외부의 경로와 심볼릭 링크는 거부됩니다. | 현재 디렉터리 |
| `--token` | 변경 및 프로세스 제어 도구용 세션/베어러 토큰입니다. 지정하지 않으면 `WAILS_MCP_TOKEN`을 사용합니다. | 안전하게 생성됨 |
| `--stdio` | stdio 전송 방식을 강제로 사용합니다. | 자동 |
| `--http` | Streamable HTTP 전송 방식을 강제로 사용합니다. | 자동 |
| `--port` | HTTP 포트입니다. `0`을 지정하면 사용 가능한 루프백 포트를 선택합니다. | `0` |

HTTP 모드에서는 Wails가 엔드포인트와 전달자 토큰을 stderr에 출력합니다. stdio 모드에서는 MCP 초기화 지침에 토큰이 포함됩니다. 서버는 임의의 셸 명령 실행 기능을 노출하지 않습니다. 원격 템플릿과 Git 원격 저장소를 사용하려면 도구의 `allowExternal` 입력을 통해 명시적으로 승인해야 합니다.

### `doctor`

시스템을 점검하고 상태 보고서를 표시합니다.

```bash
wails3 doctor
```

## 생성 명령

생성 명령을 사용하면 바인딩, 아이콘, 빌드 파일과 같은 다양한 프로젝트 자산을 만들 수 있습니다. 모든 생성 명령은 기본 명령 `wails3 generate <command>`을 사용합니다.

### `generate bindings`

Go 코드의 바인딩과 모델을 생성합니다.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-f` | 추가 Go 빌드 플래그 |  |
| `-d` | 출력 디렉터리 | `frontend/bindings` |
| `-models` | 모델 파일 이름 | `models` |
| `-index` | 인덱스 파일 이름 | `index` |
| `-ts` | TypeScript 생성 | `false` |
| `-i` | TS 인터페이스 사용 | `false` |
| `-b` | 번들 런타임 사용 | `false` |
| `-names` | ID 대신 이름 사용 | `false` |
| `-noindex` | 인덱스 파일 생략 | `false` |
| `-noevents` | 이벤트 관련 바인딩 생성 생략 | `false` |
| `-dry` | 시험 실행 | `false` |
| `-silent` | 무출력 모드 | `false` |
| `-v` | 디버그 출력 | `false` |
| `-clean` | 생성 전에 출력 디렉터리 정리 | `true` |

### `generate build-assets`

애플리케이션의 빌드 자산을 생성합니다.

```bash
wails3 generate build-assets [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-name` | 프로젝트 이름 |  |
| `-dir` | 출력 디렉터리 | `build` |
| `-silent` | 출력 억제 | `false` |
| `-company` | 회사 이름 |  |
| `-productname` | 제품 이름 |  |
| `-description` | 제품 설명 |  |
| `-version` | 제품 버전 |  |
| `-identifier` | 제품 식별자 | `com.wails.[name]` |
| `-copyright` | 저작권 고지 |  |
| `-comments` | 파일 주석 |  |

### `generate icons`

애플리케이션 아이콘을 생성합니다.

```bash
wails3 generate icons [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-input` | 입력 PNG 파일 | 필수 |
| `-windowsfilename` | Windows 출력 파일 이름 |  |
| `-macfilename` | macOS 출력 파일 이름 |  |
| `-sizes` | 아이콘 크기(쉼표로 구분) | `256,128,64,48,32,16` |
| `-example` | 예제 아이콘 생성 | `false` |
| `-iconcomposerinput` | 입력 Icon Composer 파일(`.icon`) |  |
| `-macassetdir` | Mac 자산의 출력 디렉터리(Assets.car + icns) |  |

#### Icon Composer(macOS)

macOS 26 이상에서는 Icon Composer `.icon` 파일을 사용하여 `Assets.car`와 `icons.icns`을 생성할 수 있습니다.

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Apple의 `actool` 명령을 사용하여 `.icon` 파일을 컴파일합니다. `actool` 버전 26 이상이 포함된 Xcode가 필요합니다.

Icon Composer를 사용할 때는 `build/config.yml`의 `cfBundleIconName`을 확장자를 제외한 `.icon` 파일 이름과 일치하도록 설정하십시오.

```yaml
info:
  cfBundleIconName: "appicon"
```

설정하지 않았으며 `Assets.car`이 있으면 기본값으로 `"appicon"`을 사용합니다.

### `generate syso`

Windows .syso 파일을 생성합니다.

```bash
wails3 generate syso [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-manifest` | 매니페스트 파일 경로 | 필수 |
| `-icon` | 아이콘 파일 경로 | 필수 |
| `-info` | 버전 정보 파일 경로 |  |
| `-arch` | 대상 아키텍처 | 현재 GOARCH |
| `-out` | 출력 파일 이름 | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Linux .desktop 파일을 생성합니다.

```bash
wails3 generate .desktop [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-name` | 애플리케이션 이름 | 필수 |
| `-exec` | 실행 파일 경로 | 필수 |
| `-icon` | 아이콘 경로 |  |
| `-categories` | 애플리케이션 카테고리 | `Utility` |
| `-comment` | 애플리케이션 설명 |  |
| `-terminal` | 터미널에서 실행 | `false` |
| `-keywords` | 검색 키워드 |  |
| `-version` | 애플리케이션 버전 |  |
| `-genericname` | 일반 이름 |  |
| `-startupnotify` | 시작 알림 표시 | `false` |
| `-mimetype` | 지원되는 MIME 유형 |  |
| `-output` | 출력 파일 이름 | `[name].desktop` |

### `generate runtime`

미리 빌드된 런타임을 생성합니다.

```bash
wails3 generate runtime
```

### `generate constants`

Go 코드에서 JavaScript 상수를 생성합니다.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

배포용 Windows WebView2 부트스트랩 설치 프로그램을 생성합니다.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

새 프로젝트 템플릿 디렉터리의 기본 구조를 생성합니다.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Linux AppImage를 생성합니다.

```bash
wails3 generate appimage [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-binary` | 바이너리 파일 경로 | 필수 |
| `-icon` | 아이콘 파일 경로 | 필수 |
| `-desktop` | .desktop 파일 경로 | 필수 |
| `-builddir` | 빌드 디렉터리 | 임시 디렉터리 |
| `-output` | 출력 디렉터리 | `.` |

## 서비스 명령어

서비스 명령어는 Wails 서비스를 관리하는 데 사용됩니다. 모든 서비스 명령어는 기본 명령어 `wails3 service <command>`을 사용합니다.

### `service init`

새 서비스를 초기화합니다.

```bash
wails3 service init [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-n` | 서비스 이름 | `example_service` |
| `-d` | 서비스 설명 | `Example service` |
| `-p` | 패키지 이름 |  |
| `-o` | 출력 디렉터리 | `.` |
| `-q` | 출력 억제 | `false` |
| `-a` | 작성자 이름 |  |
| `-v` | 버전 |  |
| `-w` | 웹사이트 URL |  |
| `-r` | 저장소 URL |  |
| `-l` | 라이선스 |  |

## 도구 명령어

도구 명령어는 개발 및 디버깅용 유틸리티를 제공합니다. 모든 도구 명령어는 기본 명령어 `wails3 tool <command>`을 사용합니다.

### `tool checkport`

포트가 열려 있는지 확인합니다. vite가 실행 중인지 테스트할 때 유용합니다.

```bash
wails3 tool checkport [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-port` | 확인할 포트 | `9245` |
| `-host` | 확인할 호스트 | `localhost` |

### `tool watcher`

파일을 감시하고 변경되면 명령을 실행합니다.

```bash
wails3 tool watcher [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-config` | 구성 파일 경로 | `./build/config.yml` |
| `-ignore` | 무시할 패턴 |  |
| `-include` | 포함할 패턴 |  |

### `tool cp`

파일을 복사합니다.

```bash
wails3 tool cp
```

### `tool buildinfo`

애플리케이션의 빌드 정보를 표시합니다.

```bash
wails3 tool buildinfo
```

### `tool version`

지정한 플래그에 따라 유의적 버전을 올립니다.

```bash
wails3 tool version [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-v` | 올릴 현재 버전 |  |
| `-major` | 주 버전 올리기 | `false` |
| `-minor` | 부 버전 올리기 | `false` |
| `-patch` | 패치 버전 올리기 | `false` |
| `-prerelease` | 시험판 버전 올리기(예: alpha.5에서 alpha.6로) | `false` |

이 명령은 major > minor > patch > prerelease 우선순위를 따릅니다. 입력 버전에 "v" 접두사가 있으면 이를 유지하며, 시험판 및 메타데이터 구성 요소도 그대로 유지합니다.

사용 예:

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.1.0-alpha.5 -prerelease # Output: v3.1.0-alpha.6
```

### `tool package`

Linux 패키지(deb, rpm, archlinux)를 생성합니다.

```bash
wails3 tool package [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-format` | 패키지 형식(deb, rpm, archlinux) | `deb` |
| `-name` | 실행 파일 이름 | `myapp` |
| `-config` | 구성 파일 경로 |  |
| `-out` | 출력 디렉터리 | `.` |

### `tool lipo`

아키텍처별 바이너리를 결합하여 macOS 유니버설 바이너리를 생성합니다.

```bash
wails3 tool lipo [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-output` | 출력 바이너리 경로 |  |

### `tool capabilities`

시스템의 빌드 기능을 확인합니다(Linux에서 GTK4/GTK3 사용 가능 여부).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

크로스 컴파일용 Docker 볼륨 마운트 플래그를 생성합니다. Taskfile의 `docker run` 명령에서 사용할 수 있도록 Go 모듈 캐시와 `go.mod`에 있는 모든 로컬 `replace` 지시문에 대한 `-v` 플래그를 출력합니다.

```bash
wails3 tool docker-mounts
```

### `tool has`

도구 또는 기능의 사용 가능 여부를 확인하고 stdout에 `true` 또는 `false`을 출력합니다. Taskfile의 `sh:` 변수에서 `command -v` 대신 사용할 수 있는 크로스 플랫폼 대안으로 설계되었습니다.

여러 대안 중 하나라도 사용 가능한지 확인하려면 `|`을 사용하세요.

```bash
wails3 tool has <tool>
```

#### 예제

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Taskfile에서 사용하기

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="사용 중단됨"}
`wails3 tool has-cc`은 사용 중단되었습니다. 대신 `wails3 tool has gcc|clang`을 사용하도록 Taskfile을 업데이트하세요.

@end

`wails3 tool has gcc|clang`의 이전 버전 호환 별칭입니다. PATH에서 `gcc` 또는 `clang`을 사용할 수 있는지 확인하고 `true` 또는 `false`을 출력합니다.

```bash
wails3 tool has-cc
```

## 업데이트 명령

업데이트 명령은 프로젝트 자산을 관리하고 업데이트하는 데 사용됩니다. 모든 업데이트 명령은 기본 명령 `wails3 update <command>`을 사용합니다.

### `update cli`

Wails CLI를 새 버전으로 업데이트합니다.

```bash
wails3 update cli [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-pre` | 최신 시험판으로 업데이트 | `false` |
| `-version` | 특정 버전으로 업데이트 |  |
| `-nocolour` | 컬러 출력 비활성화 | `false` |

update cli 명령을 사용하면 설치된 Wails CLI를 업데이트할 수 있습니다. 기본적으로 최신 안정 릴리스로 업데이트합니다. `-pre` 플래그를 사용하여 최신 시험판으로 업데이트하거나, `-version` 플래그를 사용하여 특정 버전을 지정할 수 있습니다.

업데이트한 후에는 프로젝트의 go.mod 파일도 같은 버전을 사용하도록 업데이트하세요.

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

지정한 구성 파일을 사용하여 빌드 자산을 업데이트합니다.

```bash
wails3 update build-assets [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-config` | 구성 파일 경로 |  |
| `-dir` | 출력 디렉터리 | `build` |
| `-silent` | 출력 표시 안 함 | `false` |
| `-company` | 회사 이름 |  |
| `-productname` | 제품 이름 |  |
| `-description` | 제품 설명 |  |
| `-version` | 제품 버전 |  |
| `-identifier` | 제품 식별자 |  |
| `-copyright` | 저작권 고지 |  |
| `-comments` | 파일 주석 |  |

## 유틸리티 명령어

유틸리티 명령어는 일반적인 작업에 유용한 단축 기능을 제공합니다. 기본 명령어에 `wails3 <command>`을 직접 사용하세요.

### `docs`

기본 브라우저에서 Wails 문서를 엽니다.

```bash
wails3 docs
```

### `releasenotes`

현재 버전 또는 지정한 버전의 릴리스 정보를 표시합니다.

```bash
wails3 releasenotes [flags]
```

#### 플래그

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-v` | 릴리스 정보를 표시할 버전 |  |
| `-n` | 색상 출력 비활성화 | `false` |

### `version`

현재 Wails 버전을 출력합니다.

```bash
wails3 version
```

### `sponsor`

기본 브라우저에서 Wails 후원 페이지를 엽니다.

```bash
wails3 sponsor

```
