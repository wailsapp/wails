---
title: "설정"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="실험적 기능"}
설정 마법사는 새 기능이며 주로 Linux에서 테스트되었습니다. 문제가 발생하면 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하고 대신 [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요.

@end

## 빠른 시작

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

마법사는 브라우저에서 열리며 종속성 검사, 프로젝트 기본값 설정, 선택적 크로스 플랫폼 빌드 설정 과정을 안내합니다. 위 지침에서는 기본 설정 흐름을 사용하며, 실험적인 Wake 프로젝트 흐름은 아래에서 설명합니다.

이제 첫 번째 프로젝트를 시작할 준비가 되었습니다.

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 수행하는 작업

- **종속성 검사** - Go, npm 및 플랫폼 도구를 확인합니다.
- **기본값 구성** - 작성자 정보, 번들 ID 접두사, 선호 템플릿을 설정합니다.
- **크로스 플랫폼 빌드** - 어느 호스트에서나 빌드할 수 있도록 Docker를 선택적으로 설정합니다.
- **코드 서명** - macOS, Windows 및 Linux용으로 선택적으로 설정합니다.

구성은 `$XDG_CONFIG_HOME/wails/defaults.yaml`에 저장되며, 사용할 수 없으면 `~/.config/wails/defaults.yaml`을 사용하고, `wails3 init`에서 이 구성을 사용합니다. 이는 프로젝트의 `build/config.yml`과 별개입니다.

## Wake 프로젝트 구성

프로젝트 디렉터리에서 [Wake 실험 기능을 활성화](/experimental/yaml-builds/#enable-yaml-builds)한 상태로 설정을 실행하세요.

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

PowerShell에서는 다음과 같이 실행합니다.

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

그러면 프로젝트에 초점을 맞춘 다음 흐름이 열립니다.

1. **프로젝트** — 네이티브 매니페스트가 없으면 프로젝트 메타데이터를 입력합니다.
2. **애플리케이션** — GUI, CLI 또는 서버와 빌드 백엔드를 선택하거나 이미 빌드된 실행 파일을 사용합니다. Pure Go는 명시적으로 선택해야 하는 실험적 옵션입니다.
3. **대상** — 운영 체제와 아키텍처를 선택합니다. iOS와 Android는 여전히 실험적 기능이며, 서명 설정은 YAML에서 구성합니다.
4. **패키지** — 패키지 형식과 서명 여부를 선택합니다.
5. **도구** — 선택한 빌드 및 패키징 경로의 요구 사항을 검토합니다.
6. **서명** — 서명을 선택한 경우 기존 서명 제어 기능을 사용하여 구성합니다.
7. **저장** — 프로젝트 구성을 기록하고 실행할 명령을 표시합니다.

필요한 도구는 선택 사항에 따라 달라집니다. Pure Go 데스크톱 컴파일에는 Docker나 C 컴파일러가 필요하지 않습니다. CGo 크로스 빌드에는 호환되는 컴파일러와 대상 라이브러리 또는 지원되는 컨테이너 도구 체인이 필요합니다. NSIS에는 여전히 `makensis`이 필요하며, 호스트와 다른 플랫폼용 GUI AppImage에는 제공된 Linux 런타임 트리가 필요합니다. 내보낸 서명 ID에는 Wails의 Go 서명 도구를 사용할 수 있지만, 네이티브 키 체인이나 인증서 저장소의 ID에는 플랫폼 도구가 계속 필요합니다. [패키징 및 서명 요구 사항](/experimental/yaml-builds/#package-and-sign-for-distribution)을 참조하세요.

이미 빌드된 실행 파일을 사용하면 Go 컴파일, 프런트엔드 설치 및 빌드, 바인딩, UPX를 건너뜁니다. 그래도 선택한 패키징 및 서명 도구와 유효한 입력 바이너리는 필요합니다. [사전 빌드된 입력 경로](/experimental/yaml-builds/#use-prebuilt-desktop-executables)를 참조하세요. 이 선택 사항은 저장 후 표시되는 명령에 적용되며 영구적인 YAML 빌드 설정이 아닙니다.

저장하면 `wails.yaml`에 기록하고, 애플리케이션 종류, 프런트엔드 활성화 여부 및 선택한 백엔드 설정을 업데이트하며, 전용 `profiles.setup` 대상/패키지/서명 선택 항목을 교체합니다. 다른 프로필, 프런트엔드 명령, 보조 실행 파일, 훅 및 관련 없는 프로젝트 설정은 그대로 유지됩니다. 백엔드 설정은 프로젝트 전체에 적용되므로 이를 상속하는 다른 프로필도 선택한 백엔드를 사용합니다. 새 프로젝트 메타데이터는 저장할 때까지 마법사에만 유지됩니다. 누락된 도구가 있으면 빌드 전에 설치할 수 있도록 알려 줍니다. 유효한 구성을 저장해도 도구를 설치하거나 빌드를 실행하지는 않습니다.

서명된 빌드는 프로젝트 외부의 인증서 파일을 포함하여 누락된 서명 설정을 사용자 수준 기본값에서 재사용합니다. 명시적인 프로젝트 설정이 우선합니다. 설정 과정에서는 개인 키, 비밀번호 또는 전역 자격 증명 경로를 `wails.yaml`에 복사하지 않습니다. 전역으로 구성된 Windows 및 Linux 키의 비밀번호는 기존 자격 증명 저장소 또는 `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`을 통해 제공하세요. 암호화된 Apple ID의 경우 비밀번호가 들어 있는 환경 변수의 이름을 `darwin.signing.credential`에 설정하세요.

저장된 프로필은 `wails3 build setup`과 함께 실행되며, 이미 빌드된 실행 파일 옵션을 선택하면 `--skip-build`이 추가됩니다.

표시된 명령을 실행할 때 Wake를 계속 활성화해 두세요. 활성 `wails.yaml`은 플래그가 설정되지 않은 경우에도 절대 Taskfile로 대체 실행되지 않습니다. 설정 과정에서는 기존 Taskfile이나 레거시 메타데이터를 변환하지 않습니다. 이러한 프로젝트를 Wake로 구성하기 전에 [검토된 마이그레이션 워크플로](/experimental/yaml-builds/#migrate-an-existing-project)를 사용하세요.

## 하위 명령

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 문제가 있나요?

1. 문제를 진단하려면 `wails3 doctor`을 실행하세요.
2. [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요.
3. `wails3 doctor` 출력을 첨부하여 [문제를 보고하세요](https://github.com/wailsapp/wails/issues/4904).
