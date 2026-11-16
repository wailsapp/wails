---
title: "설정"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="실험적 기능"}
설정 마법사는 새 기능이며 주로 Linux에서 테스트되었습니다. 문제가 발생하면 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하고, 대신 [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요.

@end

## 빠른 시작

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

마법사가 브라우저에서 열리고 종속성 확인, 전역 프로젝트 기본값, 선택적 크로스 플랫폼 도구 및 서명 자격 증명 설정을 안내합니다. 이러한 기본 설정은 여러 프로젝트에서 재사용할 수 있습니다.

이제 첫 번째 프로젝트를 시작할 준비가 되었습니다:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 기능

- **종속성 확인** - Go, npm 및 플랫폼 도구를 확인합니다
- **기본값 구성** - 작성자 정보, 번들 ID 접두사, 선호 템플릿
- **크로스 플랫폼 빌드** - 모든 호스트에서 빌드할 수 있도록 선택적으로 Docker를 설정합니다
- **코드 서명** - macOS, Windows 및 Linux용 선택적 설정

구성은 `$XDG_CONFIG_HOME/wails/defaults.yaml`에 저장되며, 사용할 수 없으면 `~/.config/wails/defaults.yaml`을 대신 사용합니다. 이 구성은 `wails3 init`에서 사용됩니다. 이는 프로젝트의 `build/config.yml`과 별개입니다.

## Wake 프로젝트 만들기

Wake 실험 기능을 활성화하면 `wails3 init -ui`에서 프로젝트별 빌드 옵션을 수집하고 프로젝트 생성에 성공했을 때 `wails.yaml`에 저장합니다. 데스크톱 프로젝트를 생성하면서 빌드 백엔드, 대상 및 패키지 형식을 선택하려면 [Wake 시작하기 튜토리얼](/experimental/wake/)을 따르세요. `wails3 setup`은 전역 기본 설정을 관리하는 역할을 그대로 유지합니다.

기존 Wake 프로젝트는 `wails.yaml`에서 구성합니다. 프로필과 패키징 설정을 변경하려면 [YAML 빌드 가이드](/experimental/yaml-builds/)를 따르세요. 기존 Taskfile 프로젝트에서 Wake를 사용하려면 먼저 [검토가 필요한 단방향 마이그레이션](/experimental/wake/#migrate-an-existing-project)을 수행해야 합니다. 설정 과정에서는 해당 Taskfile을 변환하거나 실행하지 않습니다.

## 서명 기본 설정 재사용

서명된 Wake 빌드에서는 프로젝트 외부의 인증서 파일을 포함하여 누락된 서명 설정을 사용자 수준 기본값에서 재사용합니다. 프로젝트에 명시된 설정이 우선합니다. 비공개 키, 비밀번호 및 전역 자격 증명 경로는 `wails.yaml`에 복사되지 않습니다.

전역으로 구성된 Windows 및 Linux 키의 경우 기존 자격 증명 저장소 또는 `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`을 통해 비밀번호를 제공하세요. 암호화된 Apple ID의 경우 비밀번호가 들어 있는 환경 변수의 이름을 `darwin.signing.credential`에 설정하세요.

## 하위 명령

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 문제가 있나요?

1. 문제를 진단하려면 `wails3 doctor`을 실행하세요
2. [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요
3. `wails3 doctor` 출력을 첨부하여 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하세요
