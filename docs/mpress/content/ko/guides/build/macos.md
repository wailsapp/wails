---
title: "macOS 패키징"
description: "macOS 배포용으로 Wails 애플리케이션 패키징하기"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## 비공개 macOS API

Wails v3는 기본적으로 공개 macOS API를 사용합니다. 문서화되지 않은 Apple API가 필요한 기능을 사용하려면 단일 Go 빌드 태그 `private_mac_apis`를 지정하여 앱을 빌드하세요:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Go로 직접 빌드하려면 `go build -tags private_mac_apis .`를 사용하세요(프로덕션용 빌드에는 `-tags production,private_mac_apis` 사용). 비공개 동작에 의존하는 기존 애플리케이션은 해당 동작을 유지하려면 이 태그를 추가해야 합니다. 이 태그는 macOS 데스크톱 빌드에만 적용됩니다.

영향을 받는 기능과 옵션 값의 전체 목록, 공개 빌드에서의 정확한 대체 동작, Liquid Glass 스타일 매핑 및 인스펙터 빌드 조합은 [비공개 macOS API](/guides/build/private-macos-apis/)를 참조하세요. 비공개 API 전용 작업은 태그가 없으면 아무 동작도 하지 않으며, 공개 Go API는 변경되지 않습니다.

## 애플리케이션 번들

앱을 표준 macOS `.app` 번들로 패키징하세요:

```bash
wails3 package GOOS=darwin
```

그러면 다음 항목을 포함하는 `bin/<AppName>.app`이 생성됩니다:

- `Contents/MacOS/`에 있는 컴파일된 바이너리
- `Contents/Resources/`에 있는 앱 아이콘(`icons.icns`에서 가져오거나, 애셋 카탈로그 `Assets.car`이 있으면 해당 카탈로그에서 가져옴)
- 앱 메타데이터가 포함된 `Info.plist`

## 번들 리소스

`Contents/Resources/`은 macOS 앱과 함께 제공되는 읽기 전용 파일의 표준 위치입니다. `embed`을 사용해 Go 실행 파일에 컴파일하는 대신 필요할 때 열어야 하는 대용량 템플릿, 초기 데이터, 미디어, 언어 팩 또는 기타 페이로드를 이 위치에 저장하세요.

Wails는 이미 애플리케이션 아이콘을 이 디렉터리에 배치합니다. 자체 파일을 추가하려면 `build/resources/` 같은 소스 디렉터리에 파일을 배치한 다음, `build/darwin/Taskfile.yml`의 `create:app:bundle` 태스크에 복사 단계를 추가하세요:

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Taskfile의 `darwin:run` 태스크를 사용한다면 `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`을 대상으로 하는 동일한 명령을 Taskfile의 `run` 태스크에 추가하세요.

### Go에서 리소스 읽기

macOS 플랫폼 패키지를 임포트하세요:

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

작은 파일에는 `LoadResource`을 사용하세요:

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

큰 파일에는 `ResourceFS`을 사용하세요. 이 함수는 `Contents/Resources`을 루트로 하는 `io/fs.FS`을 반환하므로, 호출자는 리소스 전체를 먼저 Go 바이트 슬라이스로 로드하지 않고도 리소스를 열고 스트리밍할 수 있습니다:

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

리소스 이름은 `Contents/Resources`을 기준으로 하는 슬래시 구분 경로입니다. 실행 파일이 `.app/Contents/MacOS`에서 실행되지 않으면 `ResourceFS`과 `LoadResource`는 `mac.ErrNotInAppBundle`을 반환합니다.

번들 리소스는 변경할 수 없는 것으로 취급하세요. 서명된 애플리케이션 내부의 파일을 변경하면 코드 서명이 무효화됩니다. 다운로드하거나 생성한 데이터 또는 사용자가 편집할 수 있는 데이터는 대신 사용자의 Application Support 디렉터리에 저장하세요.

### 유니버설 바이너리

Apple Silicon Mac과 Intel Mac 모두에서 실행되도록 빌드하세요:

```bash
wails3 task darwin:package:universal
```

그러면 두 아키텍처에서 모두 네이티브로 실행되는 단일 `.app`이 생성됩니다. 유니버설 바이너리는 어느 플랫폼에서든 빌드할 수 있으며, Linux와 Windows에서는 `wails3 tool lipo`이 자동으로 사용됩니다.

## 번들 사용자 지정

다음 항목을 사용자 지정하려면 `build/darwin/Info.plist`을 편집하세요:

- 번들 식별자(`CFBundleIdentifier`)
- 앱 이름 및 버전
- 최소 macOS 버전
- 파일 연결
- URL 스킴

앱 아이콘은 `build/` 디렉터리의 애셋으로 생성됩니다. `generate:icons` 태스크를 사용하세요:

```bash
wails3 task common:generate:icons
```

이 태스크는 `build/appicon.png`을 사용하여 `darwin/icons.icns`과 `windows/icon.ico`을 생성합니다. macOS에서는 `build/appicon.icon`(Icon Composer 형식)도 제공할 수 있습니다. 이 경우 태스크가 `-iconcomposerinput appicon.icon -macassetdir darwin`를 전달하여 `.icon` 파일에서 `Assets.car`과 `darwin/icons.icns`을 생성합니다(macOS 이외의 플랫폼에서는 건너뜀). `Assets.car`이 있으면 `Info.plist`과 `CFBundleIconName`도 그에 맞게 업데이트되도록 `update:build-assets` 태스크를 실행하세요:

```bash
wails3 task common:update:build-assets
```

`build/` 디렉터리에서 아이콘 명령을 수동으로 실행하려면 다음과 같이 하세요:

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## 코드 서명

배포할 앱에 서명하세요:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

`build/darwin/Taskfile.yml`에서 서명을 구성하세요:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### 공증

Mac App Store 외부에서 배포하는 앱은 Apple의 공증을 받아야 합니다:

```bash
wails3 task darwin:sign:notarize
```

먼저 자격 증명을 저장하세요. 대화형 마법사(`wails3 setup signing`)를 실행하거나 `notarytool`을 직접 호출하세요:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

`build/darwin/Taskfile.yml`에서 구성하세요:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

자세한 내용은 [애플리케이션 서명](/guides/build/signing/)을 참조하세요.

## DMG 설치 프로그램

Wails 3와 함께 제공되는 템플릿은 `wails3 task darwin:package:dmg`을 제공합니다. 이 태스크는 먼저 `.app`을 생성한 다음 DMG 라이브러리를 사용하여 스타일이 적용된 DMG를 빌드합니다. 기본적으로 DMG는 빨간색 용 심벌과 WAILS 워드마크가 있는 Wails 브랜드 그라데이션 배경을 사용합니다.

```bash
wails3 task darwin:package:dmg
```

하위 수준의 `darwin:create:dmg` 태스크는 기존 `.app` 번들로 DMG를 생성하며 Taskfile에서 직접 구성할 수 있습니다:

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### 기본 레이아웃

생성된 DMG에는 다음 항목이 포함됩니다:

- 왼쪽의 애플리케이션 번들
- 오른쪽의 `Applications` 링크
- 크기가 540×380픽셀인 Finder 윈도우
- 각 아이콘 아래에 레이블이 표시되는 96포인트 크기의 아이콘
- `build/darwin/dmg-background.png`에서 가져온 Wails 브랜드 배경

애플리케이션 및 `Applications` 아이콘은 구성된 윈도우 크기를 기준으로 배치되므로, `DMG_WINDOW_WIDTH` 또는 `DMG_WINDOW_HEIGHT`을 변경해도 기본 두 아이콘 레이아웃의 간격이 비례하여 유지됩니다. 최상의 결과를 얻으려면 Finder 윈도우와 픽셀 크기가 같은 배경 이미지를 사용하세요.

### DMG 애셋 교체

`build/darwin/` 아래에 생성된 파일은 일반 프로젝트 자산이므로 교체할 수 있습니다:

- `DMG_BACKGROUND`은 Finder 창의 콘텐츠 뒤에 표시되는 이미지를 제어합니다.
- `DMG_VOLUME_ICON`은 마운트된 볼륨에 표시되는 아이콘을 제어합니다.
- `DMG_FILE_ICON`은 생성된 `.dmg` 파일에 Finder가 표시하는 아이콘을 제어합니다.

볼륨 아이콘과 DMG 파일 아이콘은 서로 별개의 리소스입니다. 애플리케이션 아이콘을 교체해도 이 두 아이콘은 자동으로 교체되지 않습니다.

### 추가 파일 넣기

설치 프로그램 스크립트, 릴리스 정보, 라이선스 또는 기타 리소스를 애플리케이션과 함께 포함하려면 `DMG_FILES`을 사용하세요. 값은 쉼표로 구분된 `name=path` 쌍의 목록입니다:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=` 앞의 이름은 DMG 내부에 표시되는 파일 이름입니다. `=` 뒤의 경로는 프로젝트에 있는 원본 파일의 경로입니다. 앞뒤 공백은 무시됩니다.

표시되는 각 이름은 고유해야 합니다. 추가 파일은 애플리케이션 번들이나 `Applications` 항목을 비롯하여 패키저가 이미 생성한 항목을 대체할 수 없습니다. 이름이 충돌하면 손상된 DMG를 생성하는 대신 오류와 함께 패키징이 실패합니다.

@note{type="note"}
DMG 생성에는 macOS의 디스크 이미지 및 Finder 도구가 사용되므로 macOS에서만 지원됩니다. 교차 컴파일된 `.app` 번들은 다른 환경에서도 생성할 수 있지만, 최종 DMG는 Mac에서 만들어야 합니다.

@end

## 문제 해결

### "앱이 손상되어 열 수 없음"

앱이 서명되지 않았습니다. Developer ID 인증서로 앱에 서명하거나, 사용자가 다음과 같이 Gatekeeper를 우회할 수 있습니다:

```bash
xattr -cr /path/to/YourApp.app
```

### 공증 실패

일반적인 문제:

- **잘못된 자격 증명**: `xcrun notarytool store-credentials`(또는 `wails3 setup signing`)을 다시 실행하세요.
- **강화된 런타임 필요**: 필요한 경우 권한 설정에 `com.apple.security.cs.allow-unsigned-executable-memory`이 포함되어 있는지 확인하세요.
- **타임스탬프 누락**: 정상적인 경우 서명 프로세스에서 타임스탬프가 자동으로 포함됩니다.

### 교차 컴파일된 앱이 실행되지 않음

교차 컴파일된 macOS 바이너리는 서명되어 있지 않습니다. 테스트하기 전에 Mac으로 전송하여 서명하세요:

```bash
codesign --force --deep --sign - YourApp.app
```
