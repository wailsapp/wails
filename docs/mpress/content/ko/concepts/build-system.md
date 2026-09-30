---
title: "빌드 시스템"
description: "Wails가 애플리케이션을 빌드하고 패키징하는 방식 이해하기"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## 통합 빌드 시스템

Wails는 Go 코드를 컴파일하고, 프런트엔드 애셋을 번들링하고, 모든 항목을 단일 실행 파일에 임베드하며, 플랫폼별 빌드를 처리하는 <strong>통합 빌드 시스템</strong>을 제공합니다. 이 모든 작업을 명령 하나로 수행할 수 있습니다.

```bash
wails3 build
```

**출력:** 모든 항목이 임베드된 네이티브 실행 파일.

## 빌드 프로세스 개요

**[빌드 프로세스 다이어그램 자리표시자]**

## 빌드 단계

### 1. 분석 단계

Wails는 서비스 구조를 파악하기 위해 Go 코드를 스캔합니다:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Wails가 추출하는 항목:**

- 서비스 이름: `GreetService`
- 메서드 이름: `Greet`
- 매개변수 타입: `string`
- 반환 타입: `string`

**용도:** TypeScript 바인딩 생성

### 2. 생성 단계

#### TypeScript 바인딩

Wails는 타입 안전성이 보장되는 바인딩을 생성합니다:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**이점:**

- 완전한 타입 안전성
- IDE 자동 완성
- 컴파일 시점 오류
- JSDoc 주석

#### 프런트엔드 빌드

프런트엔드 번들러가 실행됩니다(Vite, webpack 등):

```bash
# Vite example
vite build --outDir dist
```

**수행되는 작업:**

- JavaScript/TypeScript 컴파일
- CSS 처리 및 최소화
- 애셋 최적화
- 소스 맵 생성(개발 환경만 해당)
- 출력 위치: `frontend/dist/`

### 3. 컴파일 단계

#### Go 컴파일

Go 코드는 최적화를 적용해 컴파일됩니다:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**플래그:**

- `-s`: 심볼 테이블 제거
- `-w`: DWARF 디버깅 정보 제거
- 결과: 더 작은 바이너리(약 30% 감소)

**플랫폼별 형식:**

- Windows: 아이콘이 임베드된 `.exe`
- macOS: `.app` 번들 구조
- Linux: ELF 바이너리

#### 애셋 임베딩

프런트엔드 애셋이 Go 바이너리에 임베드됩니다:

```go
//go:embed frontend/dist
var assets embed.FS
```

**결과:** 모든 항목이 포함된 단일 실행 파일.

### 4. 출력

**단일 네이티브 바이너리:**

- Windows: `myapp.exe`(약 15MB)
- macOS: `myapp.app`(약 15MB)
- Linux: `myapp`(약 15MB)

**종속성 없음**(시스템 WebView 제외).

## 개발 빌드와 프로덕션 빌드 비교

@tabs{sync-key="mode"}
[개발 환경(wails3 dev)]
**속도에 최적화:**

```bash
wails3 dev
```

**수행되는 작업:**

1. 프런트엔드 개발 서버 시작(기본적으로 9245 포트에서 Vite 실행)
2. 최적화 없이 Go 컴파일
3. 개발 서버를 가리키도록 앱 실행
4. 핫 리로드 활성화
5. 소스 맵 포함

**특징:**

- **빠른 재빌드**(프런트엔드 변경 시 &lt;1초)
- **애셋 임베딩 없음**(개발 서버에서 제공)
- **디버그 심볼** 포함
- **소스 맵** 활성화
- **상세 로깅**

**파일 크기:** 더 큼(디버그 심볼 포함 시 약 50MB)

[프로덕션(wails3 build)]
**크기와 성능에 맞게 최적화:**

```bash
wails3 build
```

**처리 과정:**

1. 프로덕션용 프런트엔드 빌드(최소화)
2. 최적화를 적용하여 Go 컴파일
3. 디버그 심볼 제거
4. 애셋 임베드
5. 단일 바이너리 생성

**특징:**

- **최적화된 코드**(최소화 및 트리 셰이킹 적용)
- **애셋 임베드**(외부 파일 없음)
- **디버그 심볼 제거**
- **소스 맵 없음**
- **최소한의 로깅**

**파일 크기:** 더 작음(약 15MB)

@end

## 빌드 명령어

### 기본 빌드

```bash
wails3 build
```

**출력:** `bin/<APP_NAME>`(Windows에서는 `bin/<APP_NAME>.exe`). `bin/` 디렉터리는 프로젝트 루트에 있습니다.

`wails3 build`은 `wails3 task build`을 얇게 감싼 래퍼입니다. 전달하는 유일한 빌드 시점 플래그는 `--tags`이며, 이 플래그는 `EXTRA_TAGS` Taskfile 변수가 됩니다.

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build`에는 `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags` 또는 `-package` 플래그가 없습니다. 크로스 컴파일, 출력 경로, 아이콘 및 패키징은 프로젝트의 Taskfile(`Taskfile.yml` + `build/config.yml`)을 통해 제어합니다.

### 크로스 플랫폼 및 플랫폼별 빌드

플랫폼 빌드는 `darwin:` / `windows:` / `linux:` 네임스페이스 아래의 Taskfile 태스크로 제공됩니다(`build/Taskfile.<platform>.yml`에 정의됨). 예:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

현재 프로젝트에서 사용할 수 있는 모든 태스크를 확인하려면 다음을 실행하세요.

```bash
wails3 task --list
```

### 아이콘 및 패키징

소스 PNG에서 플랫폼 아이콘(`build/icons.icns`, `build/icon.ico` 등)을 생성하려면 다음을 실행하세요.

```bash
wails3 generate icons -input appicon.png
```

플랫폼별 설치 프로그램/패키지를 빌드하려면 다음을 실행하세요.

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## 빌드 구성

### Taskfile.yml

Wails 3 프로젝트는 [Taskfile](https://taskfile.dev/)을 빌드 오케스트레이터로 사용합니다. 루트 `Taskfile.yml`에는 `build/`의 플랫폼별 태스크 파일이 포함됩니다.

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

`wails3 task <name>` 또는 `task <name>`로 태스크를 실행하세요.

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### 프로젝트 구성: `build/config.yml`

프로젝트 메타데이터(이름, 식별자, 버전, info-plist 값, NSIS 설정, `.desktop` 필드, 사용자 지정 프로토콜 등)는 `build/config.yml`에 있습니다. Taskfile은 아이콘, 매니페스트, 설치 프로그램 등을 생성할 때 이 파일을 읽습니다. Wails 3에는 **`build/build.json` 파일이 없습니다**.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

이 구성에서 플랫폼별 빌드 애셋을 새로 생성하려면 `wails3 generate build-assets`(또는 `wails3 update build-assets`)을 실행하세요.

## 애셋 임베딩

### 작동 방식

Wails는 Go의 `embed` 패키지를 사용합니다.

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**빌드 시:**

1. 프런트엔드를 `frontend/dist/`에 빌드
2. `//go:embed` 지시문으로 파일 포함
3. 파일을 바이너리에 컴파일
4. 바이너리에 모든 항목 포함

**런타임 시:**

1. 앱 시작
2. 메모리에서 애셋 제공
3. 애셋에 대한 디스크 I/O 없음
4. 빠른 로딩

### 사용자 지정 애셋

추가 파일을 임베드하려면 다음과 같이 하세요.

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## 빌드 최적화

### 프런트엔드 최적화

**Vite(기본값):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**결과:**

- JavaScript 최소화(약 70% 감소)
- CSS 최소화(약 60% 감소)
- 이미지 최적화
- 트리 셰이킹 적용

### Go 최적화

**컴파일러 플래그:**

```bash
-ldflags="-s -w"
```

- `-s`: 심볼 테이블 제거(약 10% 감소)
- `-w`: DWARF 디버그 정보 제거(약 20% 감소)

**추가 최적화:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: 빌드 시 변수 값 설정
- 버전 번호와 빌드 날짜에 유용합니다

### 바이너리 압축

**UPX(선택 사항):**

```bash
# After building
upx --best bin/myapp.exe
```

**결과:**

- 크기 약 50% 감소
- 시작 시간이 약간 느려짐(약 100ms)
- macOS에는 권장하지 않습니다(코드 서명 문제)

## 플랫폼별 빌드

### Windows

**출력:** `myapp.exe`

**포함 항목:**

- 애플리케이션 아이콘
- 버전 정보
- 매니페스트(UAC 설정)

**아이콘:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

그런 다음 Windows의 `tool package` 단계에서 생성된 `.ico`을 실행 파일에 포함합니다.

**매니페스트:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**출력:** `myapp.app`(애플리케이션 번들)

**구조:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**유니버설 바이너리:**

macOS Taskfile은 두 아키텍처를 모두 빌드한 후 `wails3 tool lipo`을 통해 결합하는 `darwin:build:universal`(및 `darwin:package:universal`) 태스크를 제공합니다:

```bash
wails3 task darwin:build:universal
```

### Linux

**출력:** `myapp`(ELF 바이너리)

**종속성:**

- GTK3
- WebKitGTK

**데스크톱 파일:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**설치:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## 빌드 성능

### 일반적인 빌드 시간

| 단계 | 시간 | 참고 |
| --- | --- | --- |
| 분석 | &lt;1s | Go 코드 스캔 |
| 바인딩 생성 | &lt;1s | TypeScript 생성 |
| 프런트엔드 빌드 | 5-30s | 프로젝트 크기에 따라 다름 |
| Go 컴파일 | 2-10s | 코드 크기에 따라 다름 |
| 애셋 포함 | &lt;1s | 프런트엔드 포함 |
| **합계** | **10-45s** | 첫 빌드 |
| **증분 빌드** | **5-15s** | 후속 빌드 |

### 빌드 속도 향상

**1. 빌드 캐시 사용:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. 필요한 작업만 실행:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. 병렬 빌드(여러 머신/CI):**

v3에서 Linux/Windows/macOS 간 교차 컴파일은 일반적으로 Docker `wails-cross` 컨테이너 또는 플랫폼별 전용 러너에서 수행되며, `wails3 build` 자체는 호스트 OS를 대상으로 합니다. 지원되는 워크플로는 [크로스 플랫폼 빌드](/guides/build/cross-platform/)를 참조하세요.

**4. 더 빠른 도구를 사용하세요:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## 문제 해결

### 빌드 실패

**증상:** `wails3 build`가 오류와 함께 종료됩니다

**일반적인 원인:**

1. **Go 컴파일 오류**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **프런트엔드 빌드 오류**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **누락된 종속성**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### 바이너리가 너무 큼

**증상:** 바이너리 크기가 50MB를 초과합니다

**해결 방법:**

1. **디버그 심벌 제거**(제공되는 Taskfile은 이미 `go build`에 `-ldflags="-s -w"`를 전달합니다).

2. **임베드된 애셋 확인**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **UPX 압축 사용**
  ```bash
  upx --best bin/myapp.exe
  ```


### 느린 빌드

**증상:** 빌드에 1분 넘게 걸립니다

**해결 방법:**

1. **빌드 캐시 사용**
  - Go 캐시는 자동으로 사용됩니다
  - 프런트엔드 캐시(Vite)는 자동으로 사용됩니다


2. **필요한 작업만 실행**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **프런트엔드 빌드 최적화**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## 권장 사항

### ✅ 해야 할 일

- **개발 중에는 `wails3 dev` 사용** - 빠른 반복 개발
- **릴리스에는 `wails3 build` 사용** - 최적화된 출력
- **빌드에 버전 지정** - 버전을 임베드하려면 `-ldflags` 사용
- **대상 플랫폼에서 빌드 테스트** - 크로스 컴파일은 완벽하지 않습니다
- **프런트엔드 빌드 속도를 빠르게 유지** - 번들러 구성 최적화
- **빌드 캐시 사용** - 이후 빌드 속도가 빨라집니다

### ❌ 하지 말아야 할 일

- **`build/` 디렉터리를 커밋하지 마세요** - `.gitignore`에 추가하세요
- **빌드 테스트를 건너뛰지 마세요** - 릴리스 전에 항상 테스트하세요
- **불필요한 애셋을 임베드하지 마세요** - 바이너리 크기를 작게 유지하세요
- **프로덕션에 디버그 빌드를 사용하지 마세요** - 최적화된 빌드를 사용하세요
- **코드 서명을 잊지 마세요** - 배포에 필요합니다

## 다음 단계

**애플리케이션 빌드** - 빌드 및 패키징 상세 가이드 [자세히 알아보기 →](/guides/build/building/)

**크로스 플랫폼 빌드** - 한 대의 머신에서 모든 플랫폼용으로 빌드 [자세히 알아보기 →](/guides/build/cross-platform/)

**인스톨러 만들기** - 최종 사용자용 인스톨러 만들기 [자세히 알아보기 →](/guides/installers/)

---

**빌드에 관해 궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [빌드 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/build)를 확인하세요.
