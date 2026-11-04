---
title: "설치"
description: "Wails를 설치하고 개발 환경을 설정합니다"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 지원 플랫폼

- Windows AMD64/ARM64
- macOS AMD64/ARM64. 생성된 빌드 태스크는 macOS 12.0 이상을 대상으로 합니다. 이는 표준 빌드의 배포 대상이며, 모든 사용자 지정 빌드 구성이 이전 시스템을 지원한다는 의미는 아닙니다.
- Ubuntu 24.04 AMD64/ARM64(다른 Linux에서도 작동할 수 있습니다!)

## 종속성

Wails를 설치하기 전에 몇 가지 공통 종속성을 설치해야 합니다.

@note{type="tip"}
Wails CLI를 설치한 후 `wails3 setup`을 실행하면 이러한 종속성을 자동으로 확인하고 설치에 대한 도움을 받을 수 있습니다.

@end

@tabs
[Go(최소 1.26)]
Wails v3은 Go 1.26 이상의 툴체인을 선택합니다. 이전 버전의 Go 런처는 `GOTOOLCHAIN=auto`을 사용하여 해당 툴체인을 자동으로 다운로드할 수 있습니다. 그렇지 않으면 명시적으로 설치하세요. 시작 모듈에 선언된 더 낮은 Go 버전은 Wails 종속성의 요구 사항을 재정의하지 않습니다.

[Go 다운로드 페이지](https://go.dev/dl/)에서 Go를 다운로드하세요.

공식 [Go 설치 지침](https://go.dev/doc/install)을 정확히 따르세요. 또한 `PATH` 환경 변수에 `~/go/bin` 디렉터리 경로가 포함되어 있는지 확인해야 합니다. 터미널을 다시 시작한 후 다음을 확인하세요.

- Go가 올바르게 설치되었는지 확인: `go version`
- `~/go/bin`이 PATH 변수에 포함되어 있는지 확인
  - Mac / Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm(선택 사항)]
Wails 자체에는 npm이 필요하지 않지만, 기본 제공 템플릿 대부분에는 npm이 필요합니다.

[Node 다운로드 페이지](https://nodejs.org/en/download/)에서 최신 Node 설치 프로그램을 다운로드하세요. 일반적으로 최신 릴리스를 기준으로 테스트하므로 최신 릴리스를 사용하는 것이 가장 좋습니다.

`npm --version`을 실행하여 확인하세요.

@note{type="info"}
npm 대신 다른 패키지 관리자를 사용하고 싶다면 자유롭게 사용해도 됩니다. 해당 패키지 관리자를 사용하도록 프로젝트 Taskfile을 수정해야 합니다.

@end

@end

## 플랫폼별 종속성

플랫폼별 종속성도 설치해야 합니다.

@tabs{sync-key="platform"}
[Mac]
Wails를 사용하려면 Xcode 명령줄 도구를 설치해야 합니다. 다음 명령을 실행하여 설치할 수 있습니다.

```sh
xcode-select --install
```

[Windows]
Wails를 사용하려면 [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)을 설치해야 합니다. 거의 모든 Windows 환경에는 이미 설치되어 있습니다. `wails doctor` 명령으로 확인할 수 있습니다.

[Linux]
Linux에는 표준 `gcc` 빌드 도구와 `gtk4` 및 `webkitgtk-6.0`이 필요합니다. 설치 후 <code>wails3 doctor</code>을 실행하면 종속성 설치 방법을 확인할 수 있습니다. 기존 GTK3 / WebKit2GTK 4.1 스택은 v3.1까지 `-tags gtk3`을 통해 계속 사용할 수 있습니다([Linux 패키징 - 기존 GTK3 지원](/guides/build/linux/#legacy-gtk3-support) 참조). 사용 중인 배포판 또는 패키지 관리자가 지원되지 않는다면 Discord에서 알려주세요.

@end

## 설치

Go Modules를 사용하여 Wails CLI를 설치하려면 다음 명령을 실행하세요.

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

최신 개발 버전을 설치하려면 다음 명령을 실행하세요.

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

개발 버전을 사용하면 생성되는 모든 프로젝트에서 Go의 [replace](https://go.dev/ref/mod#go-mod-file-replace) 지시문을 사용하여 해당 프로젝트가 Wails 개발 버전을 사용하도록 합니다.

## 다음 단계

CLI를 설치한 후 설정 마법사를 실행하여 개발 환경을 구성하세요.

```shell
wails3 setup
```

@note{type="caution" title="실험적 기능"}
설정 마법사는 새 기능이며 주로 Linux에서 테스트되었습니다. 문제가 발생하면 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하고 아래의 수동 종속성 설치 단계를 따르세요.

@end

설정 마법사는 다음 작업을 수행합니다.

- 플랫폼 종속성을 확인하고 설치 지원
- 프로젝트 기본값 구성(작성자 정보, 번들 ID 접두사)
- 선택적으로 크로스 플랫폼 빌드용 Docker 설정
- 코드 서명 구성(필요한 경우)

자세한 내용은 [설정 가이드](/getting-started/setup/)를 참조하세요.

## 수동 종속성 설치

종속성을 직접 설치하고 싶거나 시스템에서 설정 마법사가 작동하지 않는다면 위의 플랫폼별 지침을 따른 후 다음 명령을 실행하세요.

```shell
wails3 doctor
```

이 명령은 올바른 종속성이 설치되어 있는지 확인하고 누락된 항목을 안내합니다.

## `wails3` 명령이 보이지 않나요?

시스템에서 `wails3` 명령이 없다고 보고한다면 다음 사항을 확인하세요.

- 위의 **Go 설치 가이드**를 올바르게 따랐는지, `go/bin` 디렉터리가 `PATH` 환경 변수에 포함되어 있는지 확인하세요.
- 새 `PATH` 변수를 적용하려면 현재 터미널을 닫았다가 다시 여세요.
