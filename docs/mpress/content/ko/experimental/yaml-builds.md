---
title: "YAML 빌드"
description: "하나의 wails.yaml 파일로 실험적인 Wails v3 네이티브 빌드 시스템을 구성하고 실행합니다."
sidebar:
  order: 2
  badge:
    text: "실험적 기능"
    variant: caution
---

@note{type="caution" title="실험적 미리 보기"}
첫 번째 릴리스에서는 리터럴 구성만 지원합니다. 이 실험에서는 제한된 수명 주기 훅의 형태를 릴리스 대상으로 제안하기 전에 평가할 수 있도록 해당 훅의 프로토타입도 제공합니다. 사용자 지정 단계, YAML 표현식 및 형식이 지정된 Wails 도구 호출은 계속 보류됩니다. 기존 프로젝트는 검토된 마이그레이션을 명시적으로 활성화할 때까지 Taskfile을 계속 사용합니다.
@end

## YAML 빌드 활성화

실험적 명령을 실행하는 모든 터미널에서 `WAILS_EXP_USE_WAKE`을 설정합니다. Bash 또는 Zsh의 경우:

```bash
export WAILS_EXP_USE_WAKE=1
```

PowerShell의 경우:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

이 변수는 존재 여부를 기준으로 참여를 활성화합니다. 따라서 값이 비어 있거나 `false`인 경우에도 활성화됩니다. 실험을 비활성화하려면 변수를 설정 해제합니다(Bash/Zsh에서는 `unset WAILS_EXP_USE_WAKE`, PowerShell에서는 `Remove-Item Env:WAILS_EXP_USE_WAKE`). `wails.yaml`이 활성화된 후에는 명령을 실행할 때 Wake를 다시 활성화해야 하며, Taskfile로 대체 실행되지 않습니다. 이 변수가 없으면 YAML 빌드/개발 라우팅과 YAML 프로젝트 생성이 비활성화되고 `migrate`, `config` 및 `eject` 명령을 사용할 수 없습니다. 마이그레이션하지 않은 프로젝트는 별도의 내장 Task 런타임을 계속 사용합니다.

## 프로젝트 구성

참여 설정을 활성화한 상태에서 프로젝트 루트에 `wails.yaml`을 추가합니다. 이 파일이 있으면 네이티브 빌드 시스템이 선택되고 Wails는 해당 프로젝트의 모든 Taskfile을 무시합니다. 아래의 모든 명령을 실행할 때 참여 설정을 유지하십시오. 스키마 모델라인을 사용하면 YAML 언어 서버를 지원하는 편집기에서 자동 완성과 인라인 유효성 검사를 사용할 수 있습니다.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

기본 호스트 빌드 또는 이름이 지정된 전체 릴리스 요청을 실행합니다:

```bash
wails3 build
wails3 build release
```

일회성 빌드에서는 프로필을 추가하지 않고 쉼표로 구분된 대상과 형식을 선택합니다:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## 데스크톱 애플리케이션 교차 컴파일

실험적인 순수 Go 데스크톱 백엔드를 사용하려면 빌드 환경에서 `CGO_ENABLED: "0"`을 설정합니다. Wake는 macOS 및 Linux GUI 프로그램에 `purego` 태그를 선택하고 Windows, macOS 및 Linux 호스트에서 Go를 직접 사용합니다. 이 모드에서 데스크톱 실행 파일을 빌드하는 데는 Docker, 에뮬레이션, C 컴파일러 또는 대상 플랫폼 SDK가 필요하지 않습니다. 애플리케이션 종속 항목도 `CGO_ENABLED=0`을 지원해야 합니다.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

두 번째 명령은 두 macOS 아키텍처를 결합하고 Go를 사용해 서명되지 않은 DMG를 생성합니다. 지원되는 모든 데스크톱 호스트에서 실행할 수 있습니다. 코드 서명과 Apple 공증은 여전히 별도의 macOS 작업입니다.

확정 모드는 실행 파일마다 독립적으로 선택됩니다. 선언된 `build.environment.CGO_ENABLED`, 대상 환경 또는 바이너리 환경은 명령을 호출한 터미널의 값을 재정의합니다. 그렇지 않으면 명시적으로 선택한 `purego` 태그, CLI/서버 프로그램 또는 Windows 대상에는 기본적으로 `0`이 사용되고, macOS 및 Linux GUI 프로그램에는 기본적으로 `1`이 사용됩니다. `--plan --json`은 각 실행 파일에 대해 확정된 컴파일러, CGo 모드 및 태그를 보고합니다. 교차 컴파일은 대상 파일을 생성하지만, 호환되지 않는 데스크톱 대상을 실행하려면 여전히 적절한 호스트가 필요합니다.

`CGO_ENABLED=1`을 사용하면 Wake는 먼저 호환되는 네이티브 도구 체인 또는 구성된 `CC` 교차 컴파일러를 사용하고, 그다음 사용 가능한 경우 Windows 교차 빌드에 Zig를 사용한 후, 지원되는 환경에서는 기존 Docker/Podman 대체 수단을 사용합니다. 명시적인 `CC`은 대상 컴파일러, 라이브러리 및 sysroot를 제공해야 합니다. Zig만으로는 Linux GTK/WebKit sysroot가 제공되지 않습니다. 아래 설명과 같이 패키징 및 서명에는 추가 도구가 필요할 수 있습니다. 모바일 대상에는 CGo와 해당 플랫폼 SDK가 필요합니다.

Go 전용 컴파일은 대상 시스템의 시스템 라이브러리를 번들로 포함하지 않습니다. GUI 보조 실행 파일을 포함한 Linux GUI 애플리케이션에는 빌드 태그로 선택된 GTK/WebKitGTK 런타임이 여전히 필요합니다. 패키지 종속 항목이나 적절한 배포 환경에서 이를 제공해야 합니다. CLI/서버 전용 패키지는 기본적으로 이러한 GUI 종속 항목을 추가하지 않으며, 명시적인 종속 항목 재정의는 계속 사용할 수 있습니다.

## 여러 실행 파일 빌드

프로젝트에서 기본 프로그램을 정의합니다. `project.package`을 해당 Go 패키지(기본값 `.`)로 설정하고, `project.kind`을 `gui`, `cli` 또는 `server`(기본값 `gui`)로 설정합니다. 추가 프로그램은 `binaries` 아래에 이름이 지정된 항목으로 정의합니다:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

각 프로그램에는 자체 컴파일 입력과 캐시 입력이 사용됩니다. 가져온 프로젝트 패키지와 포함된 파일이 추적되므로, 독립적인 프로그램을 수정해도 다른 프로그램은 다시 빌드되지 않습니다. 데스크톱에서는 추가 실행 파일이 지원되며 기본 앱 번들과 설치 프로그램에 포함됩니다. 사용자가 소유한 대체 패키지 파일의 내용은 사용자가 직접 관리합니다.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

이름이 지정된 프로필은 완전한 요청입니다. 이러한 선택 재정의는 익명 빌드에서 사용하십시오. `--jobs`은 실행 작업자 수를 제한합니다. CLI/서버 프로그램만 선택하면 프런트엔드 작업을 생략합니다.

애플리케이션 코드를 실행하지 않고 버전을 읽으려면 `project.version_file`과 `project.version_symbol`을 모두 구성합니다. 파일에는 정적 Go 문자열 선언이 있어야 합니다. 그러면 리터럴 `project.version`을 생략할 수 있습니다. 저장소 리비전을 해당 실행 파일의 문자열 변수에 주입하려면 프로젝트 또는 이름이 지정된 바이너리에 `commit_symbol: main.Commit`을 설정합니다. 커밋 주입을 요청했는데 저장소 메타데이터가 없으면 오류가 발생합니다.

`upx_level: 1`에서 `9`까지의 값은 릴리스용 Linux 및 Windows 실행 파일을 선택적으로 압축하며, 0으로 설정하면 비활성화됩니다. UPX가 설치되어 있어야 합니다. 개발 빌드와 macOS 빌드에서는 압축을 생략합니다. 컴파일이나 압축에 실패해도 이전에 게시된 프로그램은 보존됩니다.

## 배포용 패키징 및 서명

### 미리 빌드된 데스크톱 실행 파일 사용

다른 빌드 시스템이나 CI 작업에서 실행 파일을 이미 생성한 경우 `--skip-build`을 사용합니다:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

각 실행 파일을 `build.output`(기본값: `bin`)에 배치합니다. 기본 실행 파일에는 `project.binary_name`을 사용하고, 보조 실행 파일에는 각 `binaries` 항목의 이름을 사용합니다. Windows 파일 이름은 `.exe`로 끝납니다. 여러 대상을 한 번에 호출하면 각 대상의 파일을 `<os>-<arch>` 하위 디렉터리에서 읽습니다. macOS GUI의 경우 Wake는 `bin/<name>`에 원시 실행 파일이 있으면 이를 읽고, 없으면 `bin/<name>.app/Contents/MacOS/<name>`을 읽습니다. `darwin/universal` 입력에는 amd64 및 arm64 실행 파일 슬라이스가 모두 포함되어야 합니다.

Wake는 파일 형식과 대상 아키텍처를 검증하고, 실행 파일을 생성된 작업 공간에 복사하며, 패키지 캐싱을 위해 해당 내용을 추적합니다. 실행 파일 입력이 없거나 호환되지 않거나 심볼릭 링크이면 다시 빌드하거나 이전 입력을 복원하는 대신 실패합니다. 제공된 실행 파일은 덮어쓰거나 chmod로 권한을 변경하지 않습니다. 실행 파일이 하나라도 변경되면 이를 사용하는 패키지의 캐시가 무효화됩니다.

이 모드에서는 프런트엔드 설치/빌드, 바인딩, Go 컴파일, UPX 및 `before_build`/`after_build` 훅을 건너뜁니다. 패키지 자산은 계속 생성하고, macOS 앱 번들을 조립하며, 구성된 패키징/서명 훅을 실행합니다. 패키지 및 서명 도구에는 아래 표의 요구 사항이 계속 적용되지만 Go 및 프런트엔드 빌드 도구는 필요하지 않습니다. 패키지 메타데이터는 현재 매니페스트에서 가져오며, 구성된 경우 `version_file`도 포함됩니다. 바이너리에서 버전을 추론하지 않습니다. GUI와 CLI 종류를 비롯해 빌드 설정과 패키지 설정을 일치시키십시오.

`wails3 package --skip-build`도 지원됩니다. `wails3 run`을 사용하면 실행 인수와 환경이 계속 적용되며, 실행 파일은 호스트에서 실행될 수 있어야 합니다. 이 모드에는 활성 YAML 프로젝트가 필요하고 프로덕션 데스크톱 대상만 지원합니다. 개발 세션이나 모바일 패키징에는 적용되지 않습니다.

서명되지 않은 패키지를 생성할 때의 호스트/도구 요구 사항은 서명할 때와 다릅니다:

| 출력 | 생성 호스트 및 도구 | 서명 |
| --- | --- | --- |
| Windows MSIX | Windows, macOS 또는 Linux. Go 작성기를 사용하며 Windows SDK는 필요하지 않습니다. | Windows 호스트 및 서명 도구. |
| Windows NSIS 설치 프로그램 | Windows, macOS 또는 Linux. 호스트 네이티브 `makensis`이 필요합니다. | Windows 호스트 및 서명 도구. |
| macOS 앱 및 DMG | Windows, macOS 또는 Linux. Go 기반 조립, 유니버설 바이너리 병합 및 이식 가능한 DMG 작성기를 사용합니다. | Apple 도구와 자격 증명이 있는 macOS. |
| Linux DEB, RPM, Arch 및 Void | Windows, macOS 또는 Linux. 네이티브 Go 패키지 작성기를 사용합니다. | DEB/RPM은 각각 외부 서명 도구를 사용하고, Arch는 아래의 Go 서명기를 사용합니다. Void 출력은 서명되지 않습니다. |
| Linux AppImage | Linux 호스트와 대상 아키텍처가 일치해야 하며, 외부 linuxdeploy 및 GTK 배포 플러그인이 필요합니다. | AppImage 서명을 위한 기본 제공 흐름은 없습니다. |

Windows MSIX 아카이브는 Go에서 직접 생성되므로 서명되지 않은 패키지를 만드는 데 Windows SDK나 Docker가 필요하지 않습니다. Windows에서 서명하려면 여전히 Windows 호스트와 해당 서명 도구가 필요합니다. NSIS는 외부 NSIS 컴파일러를 사용하며 `packages.nsis.defines`을 통해 리터럴 정의를 받습니다. macOS에서 Wails는 상속된 로캘이 UTF-8이 아닐 때만 NSIS 컴파일러 자식 프로세스에 UTF-8 로캘을 제공합니다. 터미널 및 애플리케이션 환경은 변경되지 않습니다. 패키징된 모든 실행 파일이 CLI 또는 서버이면 기본 NSIS 설치 프로그램에서 WebView2 부트스트래퍼를 제외합니다. 기본 실행 파일이나 보조 실행 파일 중 GUI가 있으면 WebView2를 유지합니다. 사용자 정의 NSIS 스크립트는 자체 런타임 설치 로직과의 호환성을 위해 부트스트래퍼를 유지합니다.

Linux는 AppImage, DEB, RPM, Arch Linux 및 Void 패키지를 지원합니다. AppImage는 외부 linuxdeploy 및 GTK 플러그인 경계를 유지합니다. zsync 파일을 요청하려면 `packages.appimage.update_info`을 설정하고, SHA256 파일을 요청하려면 `packages.appimage.checksum: true`을 설정합니다. 사이드카는 패키지 캐시에서 복원되며 아티팩트 영수증에 AppImage와 함께 기재됩니다. 현재 AppImage를 패키징하려면 Linux가 필요하며 호스트와 대상 아키텍처가 일치해야 합니다. GTK 감지 기능은 순수 Go GUI가 동적으로 로드한 라이브러리를 찾지 못하므로 이 어댑터는 아직 순수 Go GUI 배포를 지원하지 않습니다. 해당 빌드 모드에는 Linux 패키지 관리자 형식 중 하나를 사용하십시오.

Arch 패키지는 시험판 정보를 네이티브 패키지 메타데이터에 보존합니다. 예를 들어 `1.2.3-beta.2`은 `1.2.3beta.2-1`이 되며, 최종 릴리스보다 앞에 정렬됩니다. Void 종속성에는 `gtk4` 같은 단순 패키지 이름을 사용할 수 있으며, 이는 `gtk4>=0`으로 정규화됩니다. 명시적인 버전 비교와 완전한 `name-version_revision` 값은 그대로 보존됩니다. 라이브러리 파일 이름이 아닌 Void 패키지 이름을 사용하십시오.

Arch 패키지와 분리 서명을 만들려면 내보낸 ASCII-armored PGP 개인 키 파일을 사용하고, 선택적으로 암호 환경 변수의 이름을 지정하십시오.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

로컬 환경이나 CI 비밀 저장소에서 지정된 환경 변수를 설정한 후 `wails3 build release`을 실행하십시오. 암호화되지 않은 키에는 `credential`을 생략하십시오. Go 서명기는 지원되는 모든 데스크톱 호스트에서 실행되며 GPG가 필요하지 않습니다. 원본 `.pkg.tar.zst`과 같은 디렉터리의 `.pkg.tar.zst.sig`을 게시하며, 두 파일은 아티팩트 영수증에 각각 기록됩니다. 키 내용은 빌드 캐시에 스테이징되거나 저장되지 않습니다. 서명과 해당 게시 결과는 캐시에서 절대 재사용되지 않습니다. 서명에 실패하면 이전에 게시된 파일 쌍이 보존되며, 어느 한 파일의 설치에 실패하면 게시 과정에서 두 파일을 모두 롤백합니다. 그렇다고 해서 동시 읽기 작업에 대해 두 파일 시스템 이름 변경이 원자적으로 처리되는 것은 아닙니다.

macOS 앱 번들에는 보조 프로그램이 포함됩니다. 사용 안내 문구를 추가하고 사용자 정의 entitlement 파일이 제공되지 않은 경우 오디오 입력 entitlement도 추가하려면 `darwin.microphone_usage_description`을 설정하십시오. 사용자 정의 plist 및 entitlement 소스는 변경되지 않습니다. 서명과 공증에는 macOS의 기존 Apple 도구와 자격 증명을 사용합니다. 공증에는 실제 공증 프로필이 필요합니다. 서명된 출력은 `notes.signed.app` 및 `notes.signed.dmg`과 같은 플랫폼 확장자를 유지합니다.

## 개발 제어 기능

`wails3 dev` 실행 중에 명령을 입력하고 Enter 키를 누르십시오.

| 명령 | 동작 |
| --- | --- |
| `r` | 변경되지 않은 백엔드를 포함해 강제로 다시 빌드하고 재시작합니다. |
| `b` | 바인딩을 다시 생성하고 다시 빌드합니다. |
| `f` | 프런트엔드를 재시작합니다. |
| `c` | 터미널을 지웁니다. |
| `q` | 세션과 해당 자식 프로세스를 중지합니다. |
| `h` | 제어 기능을 표시합니다. |

빌드가 실패해도 현재 애플리케이션은 계속 실행됩니다. 제어 기능을 사용하려면 대화형 터미널이 필요합니다. 파이프로 전달된 입력은 세션 동작을 변경하지 않습니다.

## 실행 전에 작업 검사하기

계획 검사는 실행에서 사용하는 것과 동일한 불변 그래프를 해석합니다. 텍스트는 사람이 읽기 위한 것이며, JSON은 버전이 지정된 기계 판독 가능 계약입니다.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## 빌드하지 않고 구성 검증하기

CI나 편집기 작업에서 또는 구성을 커밋하기 전에 전용 사전 검사 명령을 사용하십시오. 프로필을 지정하지 않으면 기본 매니페스트와 이름이 지정된 모든 프로필을 검증하고, 프로필을 지정하면 해당하는 완전한 요청만 검사합니다.

```bash
wails3 config check
wails3 config check release
```

구문, 스키마, 의미 체계, 경로, 대상, 형식, 서명, 플래너 및 구성에서 결정되는 호스트 기능과 관련된 오류는 Build 및 Dev와 동일한 진단 렌더러를 사용합니다. 진단에는 필드와 소스 범위가 명시되고 소스 줄과 캐럿이 출력되며, 유용한 경우 수정 힌트가 포함됩니다. 진단의 소스 발췌문에서는 환경 값이 가려집니다. 명령줄 선택만으로 발생한 호스트 오류는 잘못된 YAML로 잘못 보고되지 않고 일반 명령 오류로 유지됩니다.

## 프로젝트 소유 에셋 사용하기

매니페스트 경로는 아이콘, 플랫폼 매니페스트, 패키지 템플릿 및 관리자 스크립트처럼 사용자가 소유한 파일을 참조합니다. Wails는 이러한 소스를 절대 편집하지 않습니다. 도구에서 사용하기 전에 일회용 `.wails/` 작업 공간으로 복사합니다. 최종 아티팩트는 `build.output` 아래에 저장되며, 검증된 영수증은 `.wails/artifacts/receipt.json`에 기록됩니다.

## Android에 배포하기

프로덕션 Android 프로필은 AAB를 생성합니다. 로컬 배포를 위해 Wails는 ABI별 디버그 APK를 내부적으로 빌드하고 ADB로 설치한 후 앱을 실행합니다.

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

선택기를 명시하지 않으면 대화형 터미널에 연결된 기기와 구성된 Android Virtual Device가 표시됩니다. 사용 가능한 대상이 여러 개인 자동화 환경에서는 `--device` 또는 `--emulator`을 제공해야 합니다.

애플리케이션 범위의 로그 스트림을 연결하고 앱이 종료되거나 Ctrl+C를 누를 때까지 대상을 모니터링하십시오.

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

연결된 실행은 애플리케이션 재시작을 계속 추적하며 대상 연결 해제 또는 오프라인 전환을 보고합니다. 빠른 반복 작업을 위해 에뮬레이터는 기본적으로 계속 실행됩니다. `--stop-emulator`은 명시적으로 지정해야 하고 `--logs`이 필요하며, 연결된 세션이 끝날 때 대상이 에뮬레이터인 경우에만 중지합니다. Wails가 새 에뮬레이터를 시작했지만 부팅 준비가 실패하거나 취소되면 해당 프로세스를 자동으로 정리합니다.

## 범위가 제한된 수명 주기 후크 실행하기

실험적 후크 인터페이스는 선언된 각 단계에서 프로젝트 소유 스크립트 하나를 호출합니다. 후크는 셸 명령이 아닙니다. `directory`이 설정되지 않은 경우 Wails는 프로젝트 루트에서 파일을 직접 실행합니다.

Unix 호스트에서는 스크립트가 실행 가능해야 하며 shebang으로 시작해야 합니다. Windows 호스트는 고정된 인터프리터 어댑터를 통해 `.cmd`, `.bat` 및 `.ps1` 스크립트를 지원합니다. Wails는 어느 플랫폼에서도 인라인 셸 명령을 구문 분석하지 않습니다.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

지원되는 단계는 `before_build`, `after_build`, `before_package`, `after_package`, `before_sign` 및 `after_sign`입니다. 빌드 후크는 상황에 따라 공유 프로젝트 전체에서 한 번 실행되거나 대상별로 한 번씩 실행됩니다. 패키징 및 서명 후크는 대상에 요청된 모든 형식을 둘러싸는 배리어를 구성하며, 배리어 내부의 형식별 작업은 계속 병렬로 실행됩니다.

후크에는 Wails가 소유하는 환경 변수 하나, 즉 `WAILS_HOOK_CONTEXT_FILE`이 전달됩니다. 이 변수는 버전이 지정된 호출 컨텍스트가 들어 있는 보호된 JSON 파일을 가리킵니다.

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

컨텍스트에는 해석된 경로가 포함되지만 상속된 프로세스 환경이나 비밀은 절대 복사되지 않습니다. Wails는 `.wails/hooks/context/` 아래에 소유자만 접근할 수 있는 파일 권한으로 컨텍스트를 원자적으로 생성합니다. 성공하거나 취소되면 이 파일은 삭제됩니다. 실패하면 진단을 위해 보존되며 오류에 해당 경로가 포함됩니다. 0이 아닌 종료 코드는 일반 Plan 단계를 실패하게 하며, 취소하면 후크 프로세스 그룹이 종료됩니다.

후크는 기본적으로 호출할 때마다 실행됩니다. 스크립트의 입력과 출력이 모두 선언된 경우에만 캐싱을 사용하도록 설정하십시오.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

스크립트와 해당 실행 가능 모드는 암시적 입력입니다. 여러 출력은 루트가 아닌 하나의 디렉터리를 공유해야 하며, 이 디렉터리가 Wails에서 기록하고 복원하는 범위 제한 아티팩트입니다. 출력 루트에는 스크립트나 선언된 입력이 포함될 수 없으며, 캐시된 출력 트리에는 일반 파일과 디렉터리만 포함될 수 있습니다. 심볼릭 링크나 특수 파일은 절대 포함할 수 없습니다. 항상 실행되는 후크는 선언되지 않은 부수 효과의 지문을 안전하게 생성할 수 없으므로 이후 작업을 캐시할 수 없게 합니다. 캐시 가능한 후크는 선언된 아티팩트 ID를 통해 이후 작업의 캐싱을 유지합니다.

## 기존 프로젝트 마이그레이션하기

비활성 초안을 만들고 검토하십시오.

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

인수를 받지 않고 프로젝트 소유 스크립트를 호출하는 보수적인 수명 주기 작업은 훅으로 변환됩니다. 그 밖에 도달 가능한 사용자 정의 Taskfile 동작은 차단 요인으로 보고됩니다. `wails.migrated.yaml`을 검토한 후 명시적으로 활성화하세요:

```bash
wails3 migrate --activate
```

활성화하면 `wails.yaml`이 원자적으로 생성됩니다. Taskfile, `build/config.yml`, 스크립트 또는 애셋을 삭제하거나 이름을 바꾸거나 편집하지 않습니다. 유효하지 않은 YAML은 YAML 오류로 실패하며, 레거시 경로로 대체되지 않습니다.

## 모든 기본값 구체화

완전히 해석된 구성을 검토하거나 복사하려면 비활성 스냅샷을 작성하세요:

```bash
wails3 eject
```

그러면 `wails.ejected.yaml`이 생성됩니다. 활성 구성은 절대 다시 작성되지 않습니다.

## 호환성 기간

`wails3 package`과 `wails3 sign`은 동일한 플래너를 가리키는 사용 중단 예정 별칭입니다. 새 자동화에서는 프로필 또는 형식과 함께 `wails3 build`을 사용해야 합니다. 도입은 릴리스별로 단계적으로 제한됩니다. 첫 번째 릴리스에서는 YAML을 새 프로젝트 경로로 사용하고, 다음 릴리스에서는 레거시 Taskfile 라우팅에 대해 경고하며, 그다음 릴리스에서는 별칭과 함께 Taskfile 실행을 제거합니다. 각 전환을 진행하려면 해당 릴리스의 네이티브 호스트 인수 테스트 매트릭스를 모두 통과해야 합니다.

모든 필드, 형식, 기본값 및 적용 규칙은 생성된 [`wails.yaml` 필드 참조](/reference/wails-yaml-fields/)에서 확인하세요.

[`wails3 run` 구성 참조](/reference/wails-yaml-run/)에서는 로컬 실행 기본값, 플랫폼별 재정의 및 매니페스트가 없을 때의 `go run .` 대체 동작을 설명합니다.
