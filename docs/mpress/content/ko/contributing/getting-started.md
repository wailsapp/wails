---
title: "시작하기"
description: "Wails 기여를 시작하는 방법 v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 기여자 여러분, 환영합니다!

Wails에 관심을 갖고 기여해 주셔서 감사합니다! 이 가이드에서는 첫 기여를 하는 방법을 안내합니다.

## 사전 요구 사항

시작하기 전에 다음 사항을 준비했는지 확인하세요.

- **Go 1.26+** 설치([다운로드](https://go.dev/dl/))
- **Node.js 20+** 및 **npm** 설치([다운로드](https://nodejs.org/))
- GitHub 계정과 연동되도록 구성된 **Git**
- Go 및 JavaScript/TypeScript에 대한 기본 지식

### 플랫폼별 요구 사항

**macOS:**

- Xcode 명령줄 도구: `xcode-select --install`

**Windows:**

- MSYS2 또는 이와 유사한 Unix 계열 환경을 권장합니다.
- WebView2 런타임(일반적으로 Windows 11에 사전 설치됨)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev`(기본 GTK4 스택)
- 설치 방법: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`(Debian/Ubuntu)
- 레거시 `-tags gtk3` 빌드 경로를 사용하려면 `libgtk-3-dev` 및 `libwebkit2gtk-4.1-dev`도 설치하세요.

## 기여 절차 개요

일반적인 기여 작업 흐름은 다음과 같습니다.

1. **포크 및 복제** - Wails 저장소의 사본을 직접 만듭니다.
2. **설정** - Wails CLI를 빌드하고 환경을 확인합니다.
3. **브랜치** - 변경 사항을 위한 기능 브랜치를 만듭니다.
4. **개발** - 코딩 표준에 따라 변경합니다.
5. **테스트** - 모든 기능이 정상적으로 작동하는지 테스트합니다.
6. **커밋** - 명확한 Conventional Commit 메시지로 커밋합니다.
7. **제출** - 검토를 받을 풀 리퀘스트를 엽니다.
8. **반복 개선** - 피드백에 답변하고 변경 사항을 조정합니다.
9. **병합** - 승인을 받으면 변경 사항이 Wails의 일부가 됩니다!

## 단계별 가이드

기여 유형을 선택하세요.

@tabs
[버그 수정]
@steps
### 버그 찾기 또는 신고하기
- [GitHub Issues](https://github.com/wailsapp/wails/issues)에 해당 버그가 이미 신고되었는지 확인하세요.
- 신고되지 않았다면 재현 단계를 포함하여 새 이슈를 만드세요.
- 작업을 시작하기 전에 확인을 기다리세요.

### 포크 및 복제
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)에서 저장소를 포크하세요.

포크한 저장소를 복제하세요.

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 빌드 및 확인
Wails를 빌드하고 버그를 재현할 수 있는지 확인하세요.

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### 버그 수정 브랜치 만들기
수정 작업을 위한 브랜치를 만드세요.

```bash
git checkout -b fix/issue-123-window-crash
```

### 버그 수정하기
- 버그 수정에 필요한 최소한의 변경만 적용하세요.
- 관련 없는 코드는 리팩터링하지 마세요.
- 회귀를 방지하도록 테스트를 추가하거나 업데이트하세요.

```bash
# Make your changes
# Add tests in *_test.go files
```

### 수정 사항 테스트하기
테스트를 실행하여 수정 사항이 제대로 작동하는지 확인하세요.

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 수정 사항 커밋하기
명확한 메시지로 커밋하세요.

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### 풀 리퀘스트 제출하기
푸시하고 PR을 만드세요.

```bash
git push origin fix/issue-123-window-crash
```

PR 설명에 다음 내용을 포함하세요.

- 버그와 근본 원인을 설명하세요.
- 수정 사항을 설명하세요.
- 이슈를 참조하세요: "Fixes #123"
- 수정 전후의 동작을 포함하세요.

### 피드백에 대응하기
검토 의견을 반영하고 필요에 따라 PR을 업데이트하세요.

@end

[WEP(개선 사항)]
@steps
### WEP 작성하기
- [WEP(Wails 개선 제안) 절차](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)를 읽어 보세요.
- WEP 템플릿을 `v3/wep/proposals/<name>/proposal.md`에 복사하세요.
- WEP만 포함하고 제목이 `[WEP] <title>`인 초안 PR을 여세요.
- 구현하기 전에 메인테이너의 결정을 기다리세요.

### 포크 및 복제
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)에서 저장소를 포크하세요.

포크한 저장소를 복제하세요.

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 개발 환경 설정
Wails를 빌드하고 환경을 확인하세요.

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 기능 브랜치 만들기
목적을 명확히 나타내는 브랜치를 만드세요.

```bash
git checkout -b feat/window-transparency-support
```

### 기능 구현하기
- [코딩 표준](/contributing/standards/)을 따르세요.
- 변경 사항은 해당 기능에만 집중하세요.
- 깔끔하고 문서화된 코드를 작성하세요.
- 포괄적인 테스트를 추가하세요.

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 철저하게 테스트하기
기능을 테스트하세요.

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### 기능 문서화하기
- 모든 공개 API에 문서 주석을 추가하세요.
- `/docs/mpress/content/`의 관련 문서를 업데이트하세요.
- 해당하는 경우 예제를 추가하세요.

### 규칙에 따라 커밋하기
Conventional Commits 규칙을 사용하세요.

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### 풀 리퀘스트 제출하기
푸시하고 PR을 만드세요.

```bash
git push origin feat/window-transparency-support
```

PR에 다음 내용을 포함하세요.

- 기능과 사용 사례를 설명하세요.
- 예제 또는 스크린샷을 제시하세요.
- 호환성을 깨는 변경 사항을 모두 나열하세요.
- 승인된 WEP PR을 참조하세요.

### 검토에 따라 반복 개선하기
유지관리자가 변경을 요청할 수 있습니다. 인내심을 갖고 협력해 주세요.

@end

[문서]
수정 PR은 이슈를 먼저 열지 않아도 환영합니다. 문서만 수정하는 경우에는 실패하는 코드 테스트가 필요하지 않습니다. M-Press 설치, 소스 경로, 미리 보기, 검증 및 PR 절차는 [M-Press로 문서 작성하기](/contributing/documentation/)를 따르세요.

@end

## 작업할 이슈 찾기

- [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue) 레이블을 찾아보세요
- [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted) 이슈를 확인하세요
- [열린 이슈](https://github.com/wailsapp/wails/issues)를 살펴보고 담당자로 지정해 달라고 요청하세요

## 도움받기

- **Discord:** [Wails Discord](https://discord.gg/JDdSxwjhGf)에 참여하세요
- **Discussions:** [GitHub Discussions](https://github.com/wailsapp/wails/discussions)에 게시하세요
- **Issues:** 재현 가능한 버그는 이슈를 여세요. 질문에는 Discussions를, 개선 사항에는 WEP PR을 사용하세요

## 행동 강령

서로 존중하고 건설적인 태도로 모두를 환영해 주세요. 우리는 함께 훌륭한 소프트웨어를 만드는 데 집중하는 친근한 커뮤니티를 만들어 가고 있습니다.

## 다음 단계

- [개발 환경](/contributing/setup/)을 설정하세요
- [코딩 표준](/contributing/standards/)을 검토하세요
- [기술 문서](/contributing/overview/)를 살펴보세요
