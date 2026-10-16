---
title: "Windows 打包"
description: "打包 Wails 应用以便在 Windows 上分发"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## NSIS 安装程序

默认打包格式会创建 NSIS 安装程序：

```bash
wails3 package GOOS=windows
```

这会运行`wails3 task windows:package`，它将：

1. 构建应用
2. 生成 WebView2 引导程序
3. 创建 NSIS 安装程序

输出：`build/windows/nsis/<AppName>-installer.exe`

### MSIX 包

用于通过 Microsoft Store 分发或进行现代 Windows 部署：

```bash
wails3 package GOOS=windows FORMAT=msix
```

输出：`bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX 需要`makeappx.exe`（Windows SDK）或独立的 MSIX 工具。Windows Taskfile 将安装任务公开为`wails3 task install:msix:tools`。

@end

## 自定义安装程序

NSIS 配置位于`build/windows/nsis/project.nsi`中。编辑此文件可自定义：

- 安装程序界面和品牌元素
- 安装目录
- 开始菜单和桌面快捷方式
- 文件关联
- 许可协议

应用元数据来自`build/windows/info.json`：

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## 代码签名

对可执行文件和安装程序进行签名，以避免 SmartScreen 警告：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

在`build/windows/Taskfile.yml`中配置签名：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

安全存储证书密码：

```bash
wails3 setup signing
```

有关详细信息，请参阅[应用签名](/guides/build/signing/)。

## 为 ARM 构建

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## 故障排除

### 找不到 makensis

安装 NSIS：

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### SmartScreen 警告

可执行文件未签名。请参阅上文的[代码签名](#heading-1)。

### 缺少 WebView2

安装程序包含 WebView2 引导程序，可在需要时下载运行时。如果需要离线安装，请从 Microsoft 下载 Evergreen Standalone Installer。
