---
title: "설정"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="실험적 기능"}
설정 마법사는 새로 도입되었으며 주로 Linux에서 테스트되었습니다. 문제가 발생하면 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하고, 대신 [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요.

@end

## 빠른 시작

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

마법사가 브라우저에서 열리고 종속 항목 확인, 프로젝트 기본값 설정, 선택적 크로스 플랫폼 빌드 설정 과정을 안내합니다.

이제 첫 번째 프로젝트를 시작할 준비가 되었습니다.

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 수행하는 작업

- **종속 항목 확인** - Go, npm 및 플랫폼 도구를 확인합니다
- **기본값 구성** - 작성자 정보, 번들 ID 접두사, 선호하는 템플릿
- **크로스 플랫폼 빌드** - 어떤 호스트에서든 빌드할 수 있는 선택적 Docker 설정
- **코드 서명** - macOS, Windows 및 Linux용 선택적 설정

구성은 `~/.config/wails/config.yaml`에 저장되며 `wails3 init`에서 사용됩니다.

## 하위 명령

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## 문제가 있나요?

1. 문제를 진단하려면 `wails3 doctor`을 실행하세요
2. [수동 설치 단계](/getting-started/installation/#platform-specific-dependencies)를 따르세요
3. `wails3 doctor` 출력을 첨부하여 [문제를 보고하세요](https://github.com/wailsapp/wails/issues/4904)
