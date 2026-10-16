---
title: "更新程式的 GitHub Release 資產"
description: "Wails 更新程式如何從 GitHub Releases 選取應用程式成品，並避開安裝程式套件。"
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

GitHub Releases 提供者會使用已設定的`AssetMatcher`選取 Release 資產。當`AssetMatcher`為`nil`時，會使用`github.DefaultAssetMatcher`。

## 預設比對

預設比對器會在每個資產的檔名中尋找目前的平台和架構。它可辨識常見的架構別名，包括：

- `amd64`、`x86_64`和`x64`
- `arm64`和`aarch64`
- `386`、`i386`、`x86`和`ia32`

簽章和總和檢查碼等附屬檔案會被忽略。

## 安裝程式資產

一個 GitHub Release 可能同時包含更新程式使用的應用程式二進位檔，以及供首次安裝使用的傳統安裝程式。預設比對器會忽略小寫檔名符合下列條件的資產：

- 包含`-installer.`
- 包含`_installer.`
- 完全等於`installer.exe`

例如，假設有下列 Windows 資產：

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher`會選取`myapp-windows-amd64.exe`並忽略安裝程式。這可防止更新程式以 NSIS 或類似方式封裝的安裝程式可執行檔取代正在執行的應用程式。

此檢查刻意限定在狹窄範圍內。名稱僅包含`installer`一詞的應用程式仍然有效，包括：

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## 自訂命名配置

當 Release 資產未遵循平台與架構命名慣例，或需要不同的安裝程式篩選方式時，請設定`AssetMatcher`：

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

自訂比對器會完全取代`DefaultAssetMatcher`，因此必須負責排除簽章、總和檢查碼、安裝程式，以及不應作為應用程式更新安裝的任何其他資產。
