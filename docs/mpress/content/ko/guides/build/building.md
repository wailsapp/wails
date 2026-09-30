---
title: "애플리케이션 빌드"
description: "Wails 애플리케이션 빌드 및 패키징"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3는 빌드 시스템으로 [Task](https://taskfile.dev)를 사용합니다. `wails3 build` 및 `wails3 package` 명령은 Task를 편리하게 사용할 수 있도록 감싼 래퍼입니다.

## 빌드

현재 플랫폼용으로 빌드합니다:

```bash
wails3 build
```

특정 플랫폼용으로 빌드합니다:

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

출력은 `bin/` 디렉터리에 저장됩니다.

@note{type="tip"}
다른 플랫폼에서 macOS 또는 Linux용으로 크로스 컴파일하려면 Docker가 필요합니다. 설정 방법은 [크로스 플랫폼 빌드](/guides/build/cross-platform/)를 참조하세요.

@end

## 개발

핫 리로드를 사용하여 애플리케이션을 실행합니다:

```bash
wails3 dev
```

파일 감시기가 시작되어 변경 사항이 있을 때 애플리케이션을 다시 빌드하고 재시작합니다. 프런트엔드 개발 서버는 기본적으로 9245 포트에서 실행됩니다.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## 패키징

배포할 애플리케이션을 패키징합니다:

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

이 명령은 다음과 같은 플랫폼별 패키지를 생성합니다:

- **Windows**: NSIS 설치 프로그램 — [Windows 패키징](/guides/build/windows/)을 참조하세요.
- **macOS**: 애플리케이션 번들(`.app`) — [macOS 패키징](/guides/build/macos/)을 참조하세요.
- **Linux**: AppImage, deb 및 rpm — [Linux 패키징](/guides/build/linux/)을 참조하세요.

## 사용자 지정 빌드 태그

`-tags` 플래그로 사용자 지정 Go 빌드 태그를 전달합니다:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

태그는 기본 Taskfile에 `EXTRA_TAGS`로 전달됩니다. 자세한 내용은 [서버 빌드](/guides/server-build/) 및 [Linux 패키징 - 레거시 GTK3 지원](/guides/build/linux/#legacy-gtk3-support)을 참조하세요.

## Task 직접 사용

더 세밀하게 제어하려면 Task를 직접 사용하세요:

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

`linux:create:deb` 또는 `darwin:build:universal` 같은 플랫폼별 작업은 Task를 통해서만 사용할 수 있습니다.

## 에셋 생성

아이콘을 다시 생성하거나 빌드 구성을 업데이트합니다:

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
