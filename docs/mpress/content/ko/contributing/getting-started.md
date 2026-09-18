---
title: "시작하기"
description: "Wails v3 기여를 시작하는 방법"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 기여자 여러분, 환영합니다!

Wails 기여에 관심을 가져 주셔서 감사합니다! 이 가이드에서는 첫 기여를 완료하는 과정을 안내합니다.

## 사전 요구 사항

시작하기 전에 다음 항목이 준비되어 있는지 확인하세요.

- **Go 1.25+** 설치([다운로드](https://go.dev/dl/))
- **Node.js 20+** 및 **npm**([다운로드](https://nodejs.org/))
- GitHub 계정으로 구성된 **Git**
- Go 및 JavaScript/TypeScript에 대한 기본 지식

### 플랫폼별 요구 사항

**macOS:**

- Xcode Command Line Tools: `xcode-select --install`

**Windows:**

- MSYS2 또는 이와 유사한 Unix 계열 환경 권장
- WebView2 런타임(일반적으로 Windows 11에 사전 설치됨)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev`(기본 GTK4 스택)
- 다음 명령으로 설치: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- 레거시 `-tags gtk3` 빌드 경로를 사용하는 경우 `libgtk-3-dev` 및 `libwebkit2gtk-4.1-dev`도 설치하세요.

## 기여 절차 개요

일반적인 기여 워크플로는 다음 단계로 진행됩니다.

1. **포크 및 복제** - Wails 저장소의 자체 사본 만들기
2. **설정** - Wails CLI를 빌드하고 환경 확인하기
3. **브랜치** - 변경 사항을 위한 기능 브랜치 만들기
4. **개발** - 코딩 표준에 따라 변경하기
5. **테스트** - 모든 기능이 작동하는지 테스트 실행하기
6. **커밋** - 명확하고 규약에 맞는 커밋 메시지로 커밋하기
7. **제출** - 검토를 위한 풀 리퀘스트 열기
8. **반복 개선** - 피드백에 답변하고 변경 사항 조정하기
9. **병합** - 승인되면 변경 사항이 Wails에 반영됩니다!

## 단계별 가이드

기여 유형을 선택하세요.

@tabs
[버그 수정]
@steps
### 버그 찾기 또는 보고하기
- [GitHub Issues](https://github.com/wailsapp/wails/issues)에 버그가 이미 보고되었는지 확인하세요.
- 보고되지 않았다면 재현 절차를 포함해 새 이슈를 만드세요.
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
- 버그 수정에 필요한 최소한의 변경만 수행하세요.
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

- 버그와 근본 원인 설명
- 수정 사항 설명
- 이슈 참조: "Fixes #123"
- 수정 전후의 동작 포함

### 피드백에 대응하기
검토 의견을 반영하고 필요에 따라 PR을 업데이트하세요.

@end

[WEP(개선 제안)]
@steps
### WEP 작성하기
- [WEP(Wails Enhancement Proposal) 절차](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)를 읽으세요.
- WEP 템플릿을 `v3/wep/proposals/<name>/proposal.md`에 복사하세요.
- WEP만 포함하고 제목이 `[WEP] <title>`인 초안 PR을 여세요.
- 구현을 시작하기 전에 메인테이너의 결정을 기다리세요.

### 포크 및 복제
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)에서 저장소를 포크하세요.

포크한 저장소를 복제하세요.

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 개발 환경 설정
Wails를 빌드하고 환경이 올바르게 설정되었는지 확인하세요:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 기능 브랜치 생성
내용을 잘 나타내는 이름으로 브랜치를 생성하세요:

```bash
git checkout -b feat/window-transparency-support
```

### 기능 구현
- [코딩 표준](/contributing/standards/)을 준수하세요
- 변경 사항은 해당 기능에 집중하세요
- 깔끔하고 문서화된 코드를 작성하세요
- 포괄적인 테스트를 추가하세요

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 철저한 테스트
기능을 테스트하세요:

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

### 기능 문서화
- 모든 공개 API에 문서화 문자열을 추가하세요
- `/docs/mpress/content/`의 관련 문서를 업데이트하세요
- 해당하는 경우 예제를 추가하세요

### 규칙에 맞게 커밋
Conventional Commits 규칙을 따르세요:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Pull Request 제출
푸시하고 PR을 생성하세요:

```bash
git push origin feat/window-transparency-support
```

PR에 다음 내용을 포함하세요:

- 기능과 사용 사례를 설명하세요
- 예제 또는 스크린샷을 제시하세요
- 호환성을 깨뜨리는 변경 사항이 있다면 모두 나열하세요
- 승인된 WEP PR을 참조하세요

### 검토 의견에 따른 수정
유지관리자가 변경을 요청할 수 있습니다. 인내심을 갖고 협력해 주세요.

@end

[문서]
먼저 이슈를 열지 않아도 문서 수정 PR을 제출할 수 있습니다. 문서만 수정하는 경우에는 실패하는 코드 테스트가 필요하지 않습니다. M-Press 설치, 소스 경로, 미리보기, 검증 및 PR 절차는 [문서 수정](/contributing/documentation/)을 따르세요.

@end

## 작업할 이슈 찾기

- [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue) 레이블을 찾아보세요
- [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted) 이슈를 확인하세요
- [열린 이슈](https://github.com/wailsapp/wails/issues)를 살펴보고 담당자로 지정해 달라고 요청하세요

## 도움받기

- **Discord:** [Wails Discord](https://discord.gg/JDdSxwjhGf)에 참여하세요
- **토론:** [GitHub Discussions](https://github.com/wailsapp/wails/discussions)에 글을 올리세요
- **이슈:** 재현 가능한 버그는 이슈를 여세요. 질문은 Discussions를 이용하고, 개선 사항은 WEP PR로 제안하세요

## 행동 강령

서로를 존중하고 건설적으로 소통하며 모두를 환영해 주세요. 훌륭한 소프트웨어를 함께 만드는 데 집중하는 친근한 커뮤니티를 만들어 가고 있습니다.

## 다음 단계

- [개발 환경](/contributing/setup/)을 설정하세요
- [코딩 표준](/contributing/standards/)을 검토하세요
- [기술 문서](/contributing/overview/)를 살펴보세요
