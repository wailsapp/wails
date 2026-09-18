---
title: "MSIX 打包"
description: "将 Wails v3 应用程序打包为 MSIX 包"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX 是现代 Windows 应用程序打包格式。Wails 可以在 Windows 构建过程中生成 MSIX 包。

MSIX 打包说明详见[Windows 打包](/guides/build/windows/#msix-package)指南。

## 直接使用 CLI 打包

在 Windows 上运行 MSIX 工具，CI 环境也一样。默认后端使用 `MakeAppx.exe`；签名还需要 `signtool.exe`。两者均包含在 Windows SDK 中。安装助手会打开 Microsoft Store，并在需要时打开 SDK 下载页面；请先完成安装，再进行打包。

在 `build/config.yml` 中设置标识信息，然后构建并打包可执行文件：

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

直接运行此命令会生成 `MyApp.msix`，文件位于当前目录。[Windows 打包任务](/guides/build/windows/#msix-package)则会指定自己的输出路径。

## CLI 选项

| 选项 | 含义 |
| --- | --- |
| `--config` | 配置文件；默认为 `build/config.yml`。 |
| `--executable`, `--name` | 现有可执行文件的路径及其在包内的文件名；两者均为必填项。 |
| `--out` | 输出文件；默认为 `<ProductName>.msix`。 |
| `--arch` | 包的架构：`x64`（默认）、`x86`、`arm`、`arm64`、`x86a64` 或 `neutral`。也接受 Go 的别名 `amd64` 和 `386`。请与可执行文件的架构保持一致。 |
| `--publisher` | 发布者标识；默认为 `CN=<companyName>`。 |
| `--cert`, `--cert-password` | 用于签名的 PFX 证书路径和密码。 |
| `--use-makeappx` | 使用 Windows SDK 的默认打包工具。 |
| `--use-msix-tool` | 显式选择 `MsixPackagingTool.exe`，该工具必须位于 `PATH` 中。 |

## 签名与 CI

在 Store 之外分发时，请使用目标计算机信任的证书进行签名。证书的 Subject 必须与 `--publisher` 完全匹配。提供 `--cert` 后，MakeAppx 后端会调用 SignTool 并使用 SHA256。请参阅 [Microsoft 签名指南](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide)。

此 Windows 工作流步骤假设已安装 Wails 和 SDK，且前一步已将 PFX 文件安全地放置在 `CERT_PATH` 指定的位置。仅包含路径的机密不会自动上传证书：

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## 文件关联与资源

在 `build/config.yml` 中添加不带前导点的扩展名；生成的清单会添加点：

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

按照[文件关联](/guides/file-associations/)中的说明处理运行时的文件打开操作。目前 MakeAppx 后端只复制可执行文件并生成透明占位图像，不会导入项目中的 `Assets/` 文件，也不会转换 `iconName` 图标。如需品牌图像或额外的 DLL，请使用自定义打包工作流。

| 生成的资源 | 尺寸（像素） |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

只有配置了文件关联时才会生成 `FileIcon.png`。这些文件位于包的 `Assets/` 目录中。

## 提交到 Store 与故障排查

在 [Partner Center 门户](https://partner.microsoft.com/dashboard) 中预留应用，并在准备提交时使用其中的包标识和发布者。Store 会在提交过程中为 MSIX 包签名，因此通过此途径分发时无需购买签名证书。请参阅 [Microsoft 包要求](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements)。

如果找不到 `MakeAppx.exe` 或 `signtool.exe`，请安装或修复 Windows SDK。Wails 会搜索 `PATH` 和标准 SDK 位置。签名失败时，请检查证书的 Subject、有效期以及目标计算机是否信任该证书；请参阅 [MSIX 故障排查](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide)。
