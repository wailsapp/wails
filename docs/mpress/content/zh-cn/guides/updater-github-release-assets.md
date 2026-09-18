---
title: "用于更新程序的 GitHub Release 资产"
description: "Wails 更新程序如何从 GitHub Release 中选择应用程序构件并避开安装程序包。"
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

GitHub Releases 提供程序使用配置的`AssetMatcher`选择 Release 资产。当`AssetMatcher`为`nil`时，它使用`github.DefaultAssetMatcher`。

## 默认匹配

默认匹配器会在每个资产文件名中查找当前平台和架构。它可识别常见的架构别名，包括：

- `amd64`、`x86_64`和`x64`
- `arm64`和`aarch64`
- `386`、`i386`、`x86`和`ia32`

签名和校验和等附属文件会被忽略。

## 安装程序资产

一个 GitHub Release 可能同时包含更新程序使用的应用程序二进制文件，以及用于首次安装的常规安装程序。默认匹配器会忽略小写文件名符合以下条件的资产：

- 包含`-installer.`
- 包含`_installer.`
- 恰好为`installer.exe`

例如，假设有以下 Windows 资产：

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher`会选择`myapp-windows-amd64.exe`并忽略安装程序。这可以防止更新程序用 NSIS 安装程序或类似打包的安装程序可执行文件替换正在运行的应用程序。

此检查的范围特意设得很窄。名称中仅仅包含单词`installer`的应用程序仍然有效，包括：

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## 自定义命名方案

如果 Release 资产未遵循平台和架构命名约定，或者需要采用不同的安装程序筛选方式，请配置`AssetMatcher`：

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

自定义匹配器会完全取代`DefaultAssetMatcher`，因此它必须负责排除签名、校验和、安装程序，以及不应作为应用程序更新安装的任何其他资产。
