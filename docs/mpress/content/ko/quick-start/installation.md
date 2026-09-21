---
title: "설치"
description: "Wails를 설치하고 애플리케이션 빌드 준비하기"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## 빠른 설치(5분)

@note{type="tip" title="요약 - 숙련된 개발자용"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

또는 `wails3 doctor` 명령으로 직접 확인하세요. [첫 번째 앱으로 이동 →](/quick-start/first-app/)

@end

## 단계별 설치

@steps
### Go 설치(필수)
Wails에는 Go 1.25 이상이 필요합니다.

@tabs{sync-key="os"}
[Windows]
<strong>[go.dev/dl](https://go.dev/dl/)</strong>에서 Windows 설치 프로그램을 다운로드하여 실행하세요.

**설치 확인:**

```powershell
go version  # Should show 1.25 or later
```

**PATH 확인:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

비어 있으면 `C:\Users\YourName\go\bin`을 PATH에 추가하세요.

[macOS]
**옵션 1: 공식 설치 프로그램**

<strong>[go.dev/dl](https://go.dev/dl/)</strong>에서 macOS 설치 프로그램(.pkg 파일)을 다운로드하여 실행하세요.

**옵션 2: Homebrew**

```bash
brew install go
```

**설치 확인:**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

`~/go/bin`이 PATH에 없으면 `~/.zshrc` 또는 `~/.bash_profile`에 추가하세요:

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**옵션 1: 공식 tarball**

<strong>[go.dev/dl](https://go.dev/dl/)</strong>에서 Linux tarball을 다운로드한 후 다음을 실행하세요:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**옵션 2: 패키지 관리자**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**PATH에 추가**(`~/.bashrc` 또는 `~/.zshrc`에 추가):

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**확인:**

```bash
go version
echo $PATH | grep go/bin
```

@end

### 플랫폼 종속 항목 설치
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime**(대개 사전 설치됨)

Windows 10/11에는 기본적으로 WebView2가 포함되어 있습니다. 없는 경우:

- [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)에서 다운로드하세요.
- 또는 나중에 `wails3 doctor`을 실행하면 안내를 받을 수 있습니다.

**이것으로 끝입니다!** 다른 종속 항목은 필요하지 않습니다.

@note{type="tip" title="Windows 11 성능 팁"}
프로젝트 저장 위치로 [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)를 사용하는 것을 고려해 보세요. Dev Drive는 개발자 워크로드에 최적화되어 있으며 빌드 시간과 디스크 액세스 속도를 최대 30%까지 크게 개선할 수 있습니다.

@end

[macOS]
**Xcode Command Line Tools**(필수)

```bash
xcode-select --install
```

표시되는 대화 상자에서 "Install"을 클릭하세요.

**확인:**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**이것으로 끝입니다!** macOS에는 기본적으로 WebKit이 포함되어 있습니다.

[Linux]
**빌드 도구 및 WebKit**

@note{type="caution" title="최소 배포판 버전"}
Wails v3에는 기본적으로 <strong>WebKitGTK 6.0</strong>가 필요합니다. WebKit2GTK 4.1만 제공하는 배포판(Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x)은 레거시 `-tags gtk3` 옵션을 활성화하여 빌드해야 합니다. WebKit2GTK 4.0만 제공하는 더 오래된 릴리스(Ubuntu 20.04, Debian 11, RHEL 8)는 지원되지 않습니다.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
기본 GTK4 스택을 사용하려면 Ubuntu 24.04 이상 또는 Debian 13 이상이 필요합니다.

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
`shell.nix` 또는 `devShell`에 추가하세요:

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[기타]
Wails를 설치한 후 `wails3 doctor`을 실행하세요. 사용 중인 배포판에 필요한 정확한 패키지가 표시됩니다.

@end

@note{type="info" title="레거시 GTK3 스택"}
대상 배포판에서 아직 WebKitGTK 6.0를 제공하지 않는 경우(예: Ubuntu 22.04 LTS, Debian 12), 대신 GTK3 + WebKit2GTK 4.1 개발 라이브러리(Debian/Ubuntu에서는 `libgtk-3-dev libwebkit2gtk-4.1-dev`, 다른 배포판에서는 이에 해당하는 패키지)를 설치하고 `wails3 build -tags gtk3` 옵션으로 빌드하세요. 레거시 경로는 v3.0.x 계열까지 지원되며 v3.1에서 제거됩니다. 자세한 내용은 [Linux 패키징 - 레거시 GTK3 지원](/guides/build/linux/#legacy-gtk3-support)을 참조하세요.

@end

@end

### Wails CLI 설치
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

이렇게 하면 `wails3` 명령이 `~/go/bin`(Windows에서는 `%USERPROFILE%\go\bin`)에 설치됩니다.

### 설정 마법사 실행(권장)
```bash
wails3 setup
```

설정 마법사는 종속 항목을 확인하고, 누락된 항목의 설치를 도우며, 프로젝트 기본값을 구성합니다.

@note{type="caution" title="실험적 기능"}
설정 마법사는 새로 도입되었으며 주로 Linux에서 테스트되었습니다. 문제가 발생하면 [문제를 보고](https://github.com/wailsapp/wails/issues/4904)하고 대신 `wails3 doctor`을 사용하세요.

@end

### 설치 확인
```bash
wails3 doctor
```

**예상 출력(또는 이와 유사한 출력):**

```
Wails (v3.0.0-dev)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | Your installed version |
| Go Version  | go1.25.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="`wails3` 명령을 찾을 수 없는 경우"}
`~/go/bin`이 PATH에 없습니다. 위의 1단계를 따라 문제를 해결한 다음 터미널을 다시 시작하세요.

@end

### npm 설치(선택 사항이지만 권장)
대부분의 Wails 템플릿은 프런트엔드 도구로 npm을 사용합니다.

@tabs{sync-key="os"}
[Windows]
[nodejs.org](https://nodejs.org/)에서 다운로드한 후 설치 프로그램을 실행하세요.

**확인:**

```powershell
npm --version
```

[macOS]
**옵션 1: 공식 설치 프로그램** [nodejs.org](https://nodejs.org/)에서 다운로드하세요.

**옵션 2: Homebrew**

```bash
brew install node
```

**확인:**

```bash
npm --version
```

[Linux]
**옵션 1: NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**옵션 2: 패키지 관리자**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**확인:**

```bash
npm --version
```

@end

@note{type="tip" title="대체 패키지 관리자"}
`pnpm`, `yarn` 또는 `bun`을(를) 선호하시나요? 문제없습니다! 프로젝트의 `Taskfile.yml`을(를) 업데이트하여 원하는 도구를 사용하세요.

@end

@end

## 문제 해결

### `wails3` 명령을 찾을 수 없음

**원인:** `~/go/bin`(또는 `%USERPROFILE%\go\bin`)이 PATH에 없습니다.

**해결 방법:**

@tabs{sync-key="os"}
[Windows]
1. "환경 변수"를 여세요(시작 메뉴에서 검색).
2. "사용자 변수"에서 `Path`을(를) 찾으세요.
3. "편집" → "새로 만들기"를 클릭하세요.
4. `C:\Users\YourName\go\bin`을(를) 추가하세요(`YourName` 교체).
5. 모든 대화 상자에서 "확인"을 클릭하세요.
6. **터미널을 다시 시작하세요**

**확인:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
`~/.zshrc`(macOS) 또는 `~/.bashrc`(Linux)에 추가하세요:

```bash
export PATH=$PATH:~/go/bin
```

다시 불러오세요:

```bash
source ~/.zshrc  # or ~/.bashrc
```

**확인:**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor`에서 누락된 종속성을 보고함

**Linux:** 출력에 설치해야 할 패키지가 정확히 표시됩니다. 예:

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows:** WebView2가 없는 경우:

- [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)에서 다운로드하세요.
- 또는 첫 번째 앱을 실행할 때 자동으로 설치됩니다.

**macOS:** Xcode 도구가 없는 경우:

```bash
xcode-select --install
```

---

#### Go 버전이 너무 오래됨

Wails v3에는 Go 1.25 이상이 필요합니다. 이전 버전을 사용 중인 경우:

@tabs{sync-key="os"}
[Windows/macOS]
[go.dev/dl](https://go.dev/dl/)에서 최신 버전을 다운로드하여 다시 설치하세요.

[Linux]
[go.dev/dl](https://go.dev/dl/)에서 최신 tarball을 다운로드한 후 다음을 실행하세요:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## 개발 버전(최신 개발판)

주 개발 브랜치의 가장 최신 코드를 사용하고 싶으신가요? 릴리스 전에 새로운 기능과 수정 사항을 사용할 수 있지만, 버그와 호환성을 깨는 변경 사항이 포함될 위험이 있습니다. 기여자 또는 출시 예정 기능을 테스트해야 하는 사용자에게만 권장합니다.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="개발 버전"}
- 버그 또는 호환성을 깨는 변경 사항이 있을 수 있음
- 생성된 프로젝트는 `replace` 지시문을 사용하여 로컬 Wails를 가리킵니다.
- 기여자 또는 새로운 기능을 테스트하는 경우에만 권장

@end

## 다음 단계

**설치가 완료되었습니다!** 이제 시스템에서 Wails 개발을 시작할 수 있습니다.

@cards{cols="1"}
🚀 첫 번째 앱 빌드하기
10분 만에 작동하는 애플리케이션을 만드세요.

[첫 번째 앱 튜토리얼 →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 템플릿 살펴보기
기본으로 제공되는 항목을 확인하세요.

```bash
wails3 init -l  # List templates
```

@end

---

**문제가 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [이슈를 등록하세요](https://github.com/wailsapp/wails/issues).
