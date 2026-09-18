---
title: "업데이터용 GitHub Release 에셋"
description: "Wails 업데이터가 GitHub Releases에서 애플리케이션 아티팩트를 선택하고 설치 프로그램 패키지를 제외하는 방법입니다."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

GitHub Releases 공급자는 구성된 `AssetMatcher`을 사용하여 릴리스 에셋을 선택합니다. `AssetMatcher`이 `nil`이면 `github.DefaultAssetMatcher`을 사용합니다.

## 기본 일치 방식

기본 매처는 각 에셋 파일 이름에서 현재 플랫폼과 아키텍처를 찾습니다. 다음을 비롯한 일반적인 아키텍처 별칭을 인식합니다.

- `amd64`, `x86_64` 및 `x64`
- `arm64` 및 `aarch64`
- `386`, `i386`, `x86` 및 `ia32`

서명 및 체크섬과 같은 사이드카 파일은 무시합니다.

## 설치 프로그램 에셋

GitHub Release에는 업데이터가 사용하는 애플리케이션 바이너리와 최초 설치용으로 제공되는 일반 설치 프로그램이 모두 포함될 수 있습니다. 기본 매처는 소문자로 변환한 파일 이름이 다음 조건에 해당하는 에셋을 무시합니다.

- `-installer.`을 포함함
- `_installer.`을 포함함
- 정확히 `installer.exe`임

예를 들어 다음과 같은 Windows 에셋이 있다고 가정합니다.

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher`은 `myapp-windows-amd64.exe`을 선택하고 설치 프로그램은 무시합니다. 이렇게 하면 업데이터가 실행 중인 애플리케이션을 NSIS 또는 이와 유사한 방식으로 패키징된 설치 프로그램 실행 파일로 대체하는 것을 방지할 수 있습니다.

이 검사는 의도적으로 제한된 범위에만 적용됩니다. 이름에 `installer`이라는 단어가 단순히 포함된 애플리케이션은 다음을 비롯해 계속 유효합니다.

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## 사용자 지정 명명 규칙

릴리스 에셋이 플랫폼 및 아키텍처 명명 규칙을 따르지 않거나 설치 프로그램을 다른 방식으로 필터링해야 하는 경우 `AssetMatcher`을 구성하세요.

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

사용자 지정 매처는 `DefaultAssetMatcher`을 완전히 대체하므로, 애플리케이션 업데이트로 설치해서는 안 되는 서명, 체크섬, 설치 프로그램 및 기타 모든 에셋을 직접 제외해야 합니다.
