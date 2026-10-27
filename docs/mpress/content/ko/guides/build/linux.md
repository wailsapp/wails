---
title: "Linux 패키징"
description: "Linux 배포용 Wails 애플리케이션 패키징"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## 패키지 형식

Linux 배포용으로 앱을 패키징합니다:

```bash
wails3 package GOOS=linux
```

그러면 `bin/` 디렉터리에 다음과 같은 여러 형식이 생성됩니다:

- **AppImage**: 모든 Linux 배포판에서 실행할 수 있는 이식 가능한 형식
- **DEB**: Debian, Ubuntu 및 그 파생 배포판용
- **RPM**: Fedora, RHEL 및 그 파생 배포판용
- **Arch**: Arch Linux 및 그 파생 배포판용

### 개별 형식

특정 형식으로 빌드합니다:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## 패키지 사용자 지정

### 데스크톱 항목

`.desktop` 파일은 애플리케이션 메뉴에 앱이 표시되는 방식을 제어합니다. 이 파일은 `build/linux/Taskfile.yml`의 값으로 생성됩니다:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### 패키지 메타데이터

DEB 및 RPM 패키지를 사용자 지정하려면 `build/linux/nfpm/nfpm.yaml`을 편집합니다:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

AppImage 구성은 `build/linux/appimage/`에 있습니다. 앱 아이콘은 `build/appicon.png`에서 가져옵니다.

AppImage는 WebKit과 그 도우미 프로세스를 함께 묶고, 번들 라이브러리가 해당 프로세스를 AppImage 내부에서 시작하도록 패치합니다. AppImage의 `AppRun`은 내부의 `usr/lib`에서 앱을 시작합니다. 이 위치에서는 패치된 경로가 WebKit의 bubblewrap 샌드박스 외부와 내부 모두에서 올바르게 해석되므로 샌드박스가 계속 활성화됩니다. 이 작업 디렉터리를 그대로 유지할 수 있는 경우에만 자체 `AppRun`을 제공하십시오.

## 패키지 서명

PGP 키로 DEB 및 RPM 패키지에 서명합니다:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

`build/linux/Taskfile.yml`에서 서명을 구성합니다:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

키 암호를 저장합니다:

```bash
wails3 signing
```

자세한 내용은 [애플리케이션 서명](/guides/build/signing/)을 참조하십시오.

## ARM용 빌드

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
x86_64 호스트에서 ARM64용으로 빌드할 때는 CGO 교차 컴파일에 Docker를 사용합니다.

@end

## 레거시 GTK3 지원

Wails v3은 기본적으로 **GTK4와 WebKitGTK 6.0**을 사용해 빌드됩니다. 아직 WebKitGTK 6.0을 제공하지 않는 배포판(Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x)을 위해 레거시 GTK3 / WebKit2GTK 4.1 경로도 계속 사용할 수 있습니다. 레거시 경로는 빌드 태그를 사용해 명시적으로 선택해야 하며 v3.1에서 제거될 예정입니다.

@note{type="caution" title="레거시 경로"}
GTK3 / WebKit2GTK 4.1 경로는 v3.0.x 계열까지 지원됩니다. 대상 배포판의 GTK4 / WebKitGTK 6.0 제공 일정에 맞춰 GTK4로 마이그레이션할 계획을 세우십시오. `-tags gtk3`은 v3.1에서 제거됩니다.

@end

### 종속성

GTK3 및 WebKit2GTK 4.1 개발 라이브러리를 설치합니다:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

필요한 pkg-config 패키지는 `gtk+-3.0` 및 `webkit2gtk-4.1`입니다.

### GTK3로 빌드

`-tags gtk3` 플래그를 사용합니다:

```bash
wails3 build -tags gtk3
```

또는 Go로 직접 빌드합니다:

```bash
go build -tags gtk3 -o myapp .
```

### GTK4와의 알려진 차이점

- **파일 대화 상자**: GTK4는 기본적으로 파일 대화 상자에 `xdg-desktop-portal`을 사용하므로 기본 디렉터리나 사용자 지정 필터 표시 같은 일부 대화 상자 옵션이 GTK3와 다르게 동작합니다. 자세한 내용은 [대화 상자 참조 - Linux 대화 상자 동작](/reference/dialogs/#linux-dialog-behavior)을 참조하십시오.
- **메뉴 스타일**: GTK4는 GNOME HIG에 따라 헤더 바에 햄버거 버튼(☰)을 표시하는 `LinuxMenuStylePrimaryMenu` 옵션을 지원합니다. 이 옵션은 `-tags gtk3` 빌드에는 영향을 주지 않습니다. [Window API - Linux MenuStyle](/reference/window/#linux)을 참조하십시오.
- **DPI 배율 조정**: GTK4는 소수 배율 조정을 지원하기 위해 `gdk_monitor_get_scale`(GTK 4.14+)을 사용합니다.

### 빌드 확인

설정을 확인하려면 `wails3 doctor`을 실행하십시오. 플래그를 지정하지 않으면 기본 구성인 GTK4 / WebKitGTK 6.0을 확인합니다. 레거시 GTK3 / WebKit2GTK 4.1 패키지는 선택 사항으로 표시됩니다.

## 문제 해결

### AppImage가 실행되지 않음

실행 권한을 부여합니다:

```bash
chmod +x MyApp-x86_64.AppImage
```

### 종속성 누락

바이너리의 빌드 변형에 맞는 라이브러리를 사용하십시오. 기본 GTK4/WebKitGTK 6.0 빌드의 경우:

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

`-tags gtk3`을 사용해 명시적으로 빌드한 바이너리의 경우에는 레거시 GTK3/WebKit2GTK 4.1 라이브러리를 대신 설치하십시오:

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

패키지 제공 여부와 이름은 배포판 릴리스에 따라 다릅니다. 필요한 라이브러리를 제공하는 지원 릴리스를 사용하거나 문서에 설명된 교차 빌드/컨테이너 경로를 사용하십시오. 레거시 라이브러리를 설치해도 기본 빌드의 요구 사항은 충족되지 않습니다. AppImage의 이식성은 패키징에 실제로 포함된 항목과 여전히 필요한 호스트 라이브러리에 따라 달라집니다.

### C 컴파일러를 찾을 수 없음

빌드 시스템에서 CGO를 사용하려면 GCC 또는 Clang이 필요합니다:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

또는 `wails3 task setup:docker`을 실행하면 빌드 시스템이 Docker를 자동으로 사용합니다.

### NVIDIA GPU에서 빈 창 또는 흰색 창이 표시됨

NVIDIA 독점 드라이버를 사용하는 Linux에서는 Wails 앱을 시작할 때 빈 창이나 흰색 창이 표시될 수 있습니다. 이는 DMA-BUF 렌더러가 NVIDIA 독점 드라이버에서 `gbm_bo_map()`과 함께 사용될 때 실패하는 WebKitGTK 버그 때문에 발생합니다(X11과 Wayland, 드라이버 버전 377–580+, 10 시리즈 GPU 및 이전 GT 710에 영향).

**Wails는 NVIDIA 커널 모듈(`/sys/module/nvidia`)을 감지하면 `WEBKIT_DISABLE_DMABUF_RENDERER=1`을 자동으로 적용하므로** 대부분의 사용자는 별도로 조치할 필요가 없습니다.

그래도 빈 창이 표시된다면(예: 모듈 경로가 보이지 않는 컨테이너에서 실행하는 경우) 앱을 시작하기 전에 환경 변수를 수동으로 설정하십시오:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

관련 업스트림 버그: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### AppImage strip 호환성

최신 Linux 배포판(Arch Linux, Fedora 39+, Ubuntu 24.04+)에서는 더 효율적인 재배치를 위해 시스템 라이브러리를 `.relr.dyn` ELF 섹션과 함께 컴파일합니다. AppImage를 생성하는 데 사용하는 `linuxdeploy` 도구에는 이러한 최신 섹션을 처리할 수 없는 이전 `strip` 바이너리가 포함되어 있습니다.

Wails는 AppImage를 빌드하기 전에 시스템 GTK 라이브러리를 확인하여 이 상황을 자동으로 감지합니다. 이 상황이 감지되면 호환성을 보장하기 위해 strip 처리를 비활성화합니다(`NO_STRIP=1`).

**이것이 의미하는 사항:**

- 영향을 받는 시스템에서는 AppImage 크기가 약간 더 커집니다(~20-40%).
- 애플리케이션 기능에는 영향이 없습니다.
- 자동으로 처리되므로 별도의 조치가 필요하지 않습니다.

최신 시스템에서 더 작은 AppImage가 필요하다면 최신 `strip` 바이너리를 설치하고, 번들 버전 대신 이 바이너리를 사용하도록 `linuxdeploy`을 구성할 수 있습니다.
