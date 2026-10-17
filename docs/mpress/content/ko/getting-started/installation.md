---
title: "설치"
description: "Wails를 설치하고 개발 환경 설정하기"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 지원 플랫폼

- Windows AMD64/ARM64
- macOS 10.15 이상 AMD64(macOS 10.13 이상에 배포 가능)
- macOS 11.0 이상 ARM64
- Ubuntu 24.04 AMD64/ARM64(다른 Linux에서도 작동할 수 있습니다!)

## 종속성

Wails를 설치하기 전에 몇 가지 공통 종속성이 필요합니다.

@note{type="tip"}
Wails CLI를 설치한 후 `wails3 setup`을 실행하면 이러한 종속성을 자동으로 확인하고 설치에 도움을 받을 수 있습니다.

@end

@tabs
[Go(최소 1.24)]
[Go 다운로드 페이지](https://go.dev/dl/)에서 Go를 다운로드하세요.

공식 [Go 설치 지침](https://go.dev/doc/install)을 따르세요. 또한 `PATH` 환경 변수에 `~/go/bin` 디렉터리 경로가 포함되어 있는지 확인해야 합니다. 터미널을 다시 시작한 후 다음을 확인하세요.

- Go가 올바르게 설치되었는지 확인: `go version`
- `~/go/bin`이 PATH 변수에 포함되어 있는지 확인
  - Mac / Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm(선택 사항)]
Wails 자체에는 npm이 필요하지 않지만, 함께 제공되는 템플릿 대부분에는 npm이 필요합니다.

[Node 다운로드 페이지](https://nodejs.org/en/download/)에서 최신 Node 설치 프로그램을 다운로드하세요. 일반적으로 최신 릴리스를 기준으로 테스트하므로 최신 버전을 사용하는 것이 좋습니다.

`npm --version`을 실행하여 확인하세요.

@note{type="info"}
npm이 아닌 다른 패키지 관리자를 선호한다면 자유롭게 사용해도 됩니다. 해당 패키지 관리자를 사용하도록 프로젝트 Taskfile을 수정해야 합니다.

@end

@end

## 플랫폼별 종속성

플랫폼별 종속성도 설치해야 합니다.

@tabs{sync-key="platform"}
[Mac]
Wails를 사용하려면 Xcode 명령줄 도구가 설치되어 있어야 합니다. 다음 명령을 실행하여 설치할 수 있습니다.

```sh
xcode-select --install
```

[Windows]
Wails를 사용하려면 [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)이 설치되어 있어야 합니다. 거의 모든 Windows 환경에는 이미 설치되어 있습니다. `wails doctor` 명령으로 확인할 수 있습니다.

[Linux]
Linux에는 표준 `gcc` 빌드 도구와 `gtk4` 및 `webkitgtk-6.0`가 필요합니다. 설치 후 <code>wails3 doctor</code>를 실행하면 종속성 설치 방법을 확인할 수 있습니다. 기존 GTK3 / WebKit2GTK 4.1 스택은 v3.1까지 `-tags gtk3`을 통해 계속 사용할 수 있습니다([Linux 패키징 - 기존 GTK3 지원](/guides/build/linux/#legacy-gtk3-support) 참조). 사용 중인 배포판이나 패키지 관리자가 지원되지 않는다면 Discord에서 알려 주세요.

@end

## 설치

Go Modules를 사용하여 Wails CLI를 설치하려면 다음 명령을 실행하세요.

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

최신 개발 버전을 설치하려면 다음 명령을 실행하세요.

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

개발 버전을 사용하면 생성되는 모든 프로젝트에 Go의 [replace](https://go.dev/ref/mod#go-mod-file-replace) 지시문이 사용되어 프로젝트에서 Wails 개발 버전을 사용하도록 보장합니다.

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

종속성을 직접 설치하려고 하거나 시스템에서 설정 마법사가 작동하지 않는다면 위의 플랫폼별 지침을 따른 후 다음 명령을 실행하세요.

```shell
wails3 doctor
```

이 명령은 올바른 종속성이 설치되어 있는지 확인하고 누락된 항목을 알려 줍니다.

## `wails3` 명령을 찾을 수 없나요?

시스템에서 `wails3` 명령을 찾을 수 없다고 보고하면 다음 사항을 확인하세요.

- 위의 <strong>Go 설치 가이드</strong>에 나온 안내를 올바르게 따랐는지, 그리고 `go/bin` 디렉터리가 `PATH` 환경 변수에 포함되어 있는지 확인하세요.
- 새 `PATH` 변수를 적용하려면 현재 터미널을 닫았다가 다시 여세요.
