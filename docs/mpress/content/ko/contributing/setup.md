---
title: "개발 환경 설정"
description: "Wails v3 개발을 위한 개발 환경을 설정합니다"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## 개발 환경 설정

이 가이드에서는 Wails v3 작업에 필요한 완전한 개발 환경을 설정하는 방법을 안내합니다.

## 필수 도구

### Go 개발

1. **Go 1.25 이상을 설치하세요.**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Go 환경을 구성하세요.**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **유용한 Go 도구를 설치하세요.**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js 및 npm

프런트엔드 통합 예제에만 필요합니다.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### 플랫폼별 종속성

**macOS:**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows:**

1. Unix 계열 환경을 사용하려면 [MSYS2](https://www.msys2.org/)를 설치하세요
2. WebView2 Runtime(Windows 11에는 사전 설치되어 있으며, Windows 10용은 [다운로드](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)할 수 있음)
3. 선택 사항: [Git for Windows](https://git-scm.com/download/win)를 설치하세요

**Linux(Debian/Ubuntu):**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux(Fedora/RHEL):**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux(Arch):**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## 저장소 설정

### 복제 및 구성

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Wails CLI 빌드

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### PATH에 추가(선택 사항)

**Linux/macOS:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows:**

시스템 속성에서 Wails 디렉터리를 PATH 환경 변수에 추가하세요.

## IDE 설정

### VS Code(권장)

1. **VS Code를 설치하세요.** [다운로드](https://code.visualstudio.com/)

2. **확장 프로그램을 설치하세요.**
  - Go(Google의 Go Team 제공)
  - ESLint
  - Prettier
  - MDX(문서용)


3. **워크스페이스 설정을 구성하세요**(`.vscode/settings.json`).
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **GoLand를 설치하세요.** [다운로드](https://www.jetbrains.com/go/)

2. **다음과 같이 구성하세요.**
  - Go 모듈 지원을 활성화하세요
  - `goimports`용 파일 감시자를 설정하세요
  - 프로젝트 규칙에 맞게 코드 스타일을 구성하세요


## 설정 확인

다음 명령을 실행하여 모든 항목이 정상적으로 작동하는지 확인하세요.

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

테스트 앱이 빌드되고 실행되면 환경 설정이 완료된 것입니다.

## 테스트 실행

### 단위 테스트

```bash
cd v3
go test ./...
cd ..
```

### 특정 패키지 테스트

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### 커버리지 측정과 함께 실행

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 경쟁 상태 감지기와 함께 실행

```bash
cd v3
go test ./... -race
```

## 문서 작업

Wails v3 문서는 M-Press로 작성됩니다. 영어 원본 파일은  
`docs/mpress/content/`에 있으며, 번역 파일은 `fr/` 및 `id/` 같은 언어별 디렉터리에 있습니다.

저장소 루트에서 문서 변경 사항을 미리 보고 검증하세요:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

프로덕션 사이트는 정적 사이트입니다. 로컬에서 문서 작업을 할 때는 Node.js, 번역 제공업체 및 Cloudflare  
자격 증명이 필요하지 않습니다.

## 디버깅

### Go 코드 디버깅

**VS Code:**

`.vscode/launch.json`을 생성하세요.

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**명령줄:**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### 플랫폼 코드 디버깅

플랫폼별 디버깅에는 해당 플랫폼용 도구가 필요합니다.

- **macOS:** Xcode Instruments
- **Windows:** Visual Studio Debugger
- **Linux:** GDB

## 일반적인 문제

### "command not found: wails3"

Wails 디렉터리를 PATH에 추가하거나 프로젝트 루트에서 `./wails3`을 사용하세요.

### "webkitgtk-6.0 not found" 또는 "webkit2gtk not found"(Linux)

빌드 대상으로 사용하는 스택에 맞는 개발 패키지를 설치하세요.

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Go 모듈 오류로 빌드 실패

```bash
cd v3
go mod tidy
go mod download
```

### Windows에서 발생하는 "CGO_ENABLED" 오류

C 컴파일러(MSYS2를 통한 MinGW-w64)가 PATH에 포함되어 있는지 확인하세요.

## 다음 단계

- [코딩 표준](/contributing/standards/)을 검토하세요.
- [기술 문서](/contributing/)를 살펴보세요.
- 작업할 이슈를 찾아보세요: [초보자에게 적합한 이슈](https://github.com/wailsapp/wails/labels/good%20first%20issue)
