---
title: "CLI 명령 개요"
description: "프로젝트, 빌드, 서명, 도구 및 서비스 등 용도별 Wails CLI 명령을 찾아보세요"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

용도에 맞는 명령을 찾으려면 여기에서 시작하세요. 명령별 플래그와 사용 예시는 [CLI 플래그 및 예시](/guides/cli/)를 참조하세요.

## 개요

Wails CLI(`wails3`)는 Wails 3 애플리케이션을 생성, 개발, 빌드, 서명, 패키징 및 검사하기 위한 명령줄 진입점입니다. 대부분의 빌드 조정 작업은 프로젝트별 Taskfile(프로젝트의 `build/` 아래에 있음)에 위임됩니다. 많은 `wails3` 명령은 특정 태스크를 호출하는 간단한 래퍼입니다.

명령의 최신 도움말을 확인하려면 다음을 실행하세요:

```bash
wails3 --help
wails3 <command> --help
```

## 프로젝트 수명 주기

| 명령 | 설명 |
| --- | --- |
| `wails3 init` | 템플릿에서 새 프로젝트를 생성합니다. 플래그: `-n`(프로젝트 이름), `-t`(템플릿, 기본값 `vanilla`), `-p`(Go 패키지 이름, 기본값 `main`), `-d`(프로젝트 디렉터리, 기본값 `.`), `-q`(출력 최소화), `-l`(템플릿 목록 표시), `-mod`(Go 모듈 경로), `--git`(Git 저장소 URL), `--skipgomodtidy`, `-s`(원격 템플릿 경고 건너뛰기), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | 프런트엔드 핫 리로드를 사용하여 애플리케이션을 개발 모드로 실행합니다. 플래그: `--config`(기본값 `./build/config.yml`), `--port`(Vite 개발 포트), `-s`(HTTPS 활성화). |
| `wails3 build` | 프로젝트를 빌드합니다. Taskfile의 `build` 태스크를 감싸는 간단한 래퍼입니다. 플래그: `--tags`(`EXTRA_TAGS=`로 전달), `--obfuscated`(Garble로 빌드. [난독화 빌드](/guides/build/obfuscation/) 참조), `--garbleargs`(`build` 하위 명령 앞에서 `garble`에 전달할 추가 플래그). |
| `wails3 package` | 플랫폼별 `package` Taskfile 태스크를 실행합니다. |
| `wails3 task [name]` | 원하는 Taskfile 태스크를 실행합니다. 이름을 지정하지 않으면 `--list`에서 등록된 모든 태스크를 표시합니다. |
| `wails3 mcp` | 프로젝트 MCP 서버를 시작합니다. 에이전트가 시작한 프로세스에는 자동으로 stdio를 사용하고, 대화형 터미널에서 실행할 때는 루프백 Streamable HTTP를 사용합니다. |
| `wails3 doctor` | 환경 진단 보고서를 출력합니다. |
| `wails3 doctor-ng` | `doctor`의 최신 TUI 변형입니다. |
| `wails3 version` | CLI 버전을 출력합니다. |
| `wails3 releasenotes` | 최근 릴리스 정보를 출력합니다. |
| `wails3 docs` | 브라우저에서 문서 사이트를 엽니다. |
| `wails3 sponsor` | 후원 페이지를 엽니다. |

## 생성

`wails3 generate <subcommand>`:

| 하위 명령 | 설명 |
| --- | --- |
| `generate bindings` | Go와 프런트엔드 간 바인딩을 생성합니다. 플래그: `-d`(출력 디렉터리), `-models`, `-index`, `-ts`, `-i`(인터페이스), `-b`(번들), `-names`(`Call.ByName` 출력), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean`(기본값 `true`), `-f`, `-obfuscated`(Garble 빌드를 위한 안정적인 바인딩 ID가 포함된 `wails_obfuscated.gen.go` 생성. [난독화 빌드](/guides/build/obfuscation/) 참조), `-obfuscated-output`(생성된 파일을 저장할 디렉터리. 기본값은 main 패키지 디렉터리). 패키지 패턴(예: `./...`)을 인수로 받을 수 있으며, 지정하지 않으면 현재 디렉터리를 사용합니다. |
| `generate icons` | 원본 PNG를 플랫폼 아이콘 형식으로 변환합니다. 플래그: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | `build/config.yml`에서 `build/` 디렉터리의 콘텐츠(Taskfile 조각, NSIS 파일, `Info.plist`, `.desktop` 템플릿 등)를 생성합니다. |
| `generate runtime` | 웹뷰에 제공되는 사전 빌드된 `/wails/runtime.js`을 다시 생성합니다. |
| `generate syso` | Windows용 `.syso` 리소스 파일(아이콘 + 매니페스트 + 버전 정보)을 생성합니다. |
| `generate webview2bootstrapper` | Windows용 WebView2 부트스트랩 설치 프로그램을 생성합니다. |
| `generate constants` | Go 이벤트 형식에서 JS 이벤트 이름 상수를 생성합니다. |
| `generate template` | 새 프로젝트 템플릿의 기본 구조를 생성합니다. |
| `generate .desktop` | Linux용 `.desktop` 파일(AppImage/DEB/RPM에서 사용)을 생성합니다. |
| `generate appimage` | AppImage 빌드 디렉터리를 생성합니다. |

## 업데이트

`wails3 update <subcommand>`:

| 하위 명령 | 설명 |
| --- | --- |
| `update build-assets` | `build/config.yml`에서 `build/` 디렉터리를 새로 고칩니다(가능한 경우 사용자의 수정 사항을 보존합니다). |
| `update cli` | `wails3` 바이너리를 자체 업데이트합니다. |

## 코드 서명 및 패키징

| 명령 | 설명 |
| --- | --- |
| `wails3 signing` | `build/`에서 감지한 플랫폼의 서명을 구성하는 대화형 마법사입니다. 플래그: `--platform`(반복 지정 가능. 기본적으로 빌드 디렉터리에서 자동 감지). |
| `wails3 entitlements` | macOS 엔타이틀먼트용 대화형 마법사입니다. 플래그: `--output`(경로, 기본값 `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | 현재 OS(또는 `GOOS`을 통해 지정한 OS)에 맞는 플랫폼별 `*:sign` Taskfile 태스크를 실행하는 래퍼입니다. |
| `wails3 tool sign` | 저수준 직접 서명 진입점입니다. 플래그: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

이 릴리스에서는 CLI 파서가 최상위 수준에 `signing` 및 `entitlements`을 노출합니다. `wails3 signing --help` 및 `wails3 entitlements --help`을 확인하세요. `wails3 setup signing`은 서명 마법사를 선택하지 않습니다. 자격 증명을 직접 관리하려면 `xcrun notarytool store-credentials` 또는 `gpg`을 사용하세요.

## 도구

`wails3 tool <subcommand>`:

| 하위 명령 | 설명 |
| --- | --- |
| `tool checkport` | TCP 포트가 열려 있는지 확인합니다(Vite를 기다릴 때 유용함). |
| `tool watcher` | 감시 중인 파일이 변경될 때마다 명령을 실행합니다. |
| `tool cp` | 플랫폼 간에 파일을 복사합니다. |
| `tool buildinfo` | 바이너리에 포함된 Go 빌드 정보를 출력합니다. |
| `tool package` | `build/linux/nfpm`에서 Linux 패키지(`deb`, `rpm`, `archlinux`)를 빌드합니다. |
| `tool version` | 프로젝트의 시맨틱 버전을 올립니다. |
| `tool lipo` | 여러 macOS 아키텍처용 바이너리를 하나의 유니버설 바이너리로 결합합니다. |
| `tool capabilities` | 시스템에서 GTK3/GTK4 및 WebKit을 사용할 수 있는지 조사합니다. |
| `tool sign` | ([코드 서명 및 패키징](#---)을 참조하세요.) |

## 서비스

`wails3 service <subcommand>`:

| 하위 명령 | 설명 |
| --- | --- |
| `service init` | 새 서비스 패키지의 기본 구조를 생성합니다. |

## iOS

`wails3 ios <subcommand>`:

| 하위 명령 | 설명 |
| --- | --- |
| `ios overlay:gen` | iOS 브리지 심용 Go 오버레이를 생성합니다. |
| `ios xcode:gen` | 출력 디렉터리에 Xcode 프로젝트를 생성합니다. |

## 빌드 출력 경로

- 네이티브 바이너리는 `bin/<APP_NAME>`에 저장됩니다(Windows에서는 `bin/<APP_NAME>.exe`). `build/bin/`은 없습니다.
- 패키징된 출력물(`.app`, `.dmg`, NSIS 설치 프로그램, MSIX, DEB/RPM/AppImage)도 마찬가지로 `bin/`에 저장됩니다(또는 관련 Taskfile 작업이 생성한 플랫폼별 하위 디렉터리에 저장됨).

## 공통 플래그

| 플래그 | 적용 대상 | 설명 |
| --- | --- | --- |
| `-nocolour` | `--help`에서 이 플래그를 노출하는 명령(예: `init` 및 `update cli`) | CLI 출력에서 ANSI 색상을 비활성화합니다. `version`을 비롯해 모든 명령에서 허용되는 것은 아닙니다. |

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [예제](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)를 확인하세요.
