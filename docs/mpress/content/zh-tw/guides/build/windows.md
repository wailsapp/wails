---
title: "Windows 封裝"
description: "封裝 Wails 應用程式以供 Windows 發佈"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## NSIS 安裝程式

預設封裝格式會建立 NSIS 安裝程式：

```bash
wails3 package GOOS=windows
```

這會執行`wails3 task windows:package`，其作用如下：

1. 建置應用程式
2. 產生 WebView2 啟動載入器
3. 建立 NSIS 安裝程式

輸出：`build/windows/nsis/<AppName>-installer.exe`

### MSIX 套件

若要透過 Microsoft Store 發佈或使用現代化 Windows 部署：

```bash
wails3 package GOOS=windows FORMAT=msix
```

輸出：`bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX 需要`makeappx.exe`（Windows SDK）或獨立的 MSIX 工具。Windows Taskfile 將安裝工作公開為`wails3 task install:msix:tools`。

@end

## 自訂安裝程式

NSIS 組態位於`build/windows/nsis/project.nsi`。編輯此檔案可自訂：

- 安裝程式使用者介面與品牌識別
- 安裝目錄
- 開始功能表與桌面捷徑
- 檔案關聯
- 授權合約

應用程式中繼資料來自`build/windows/info.json`：

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

## 程式碼簽署

簽署可執行檔與安裝程式，以避免 SmartScreen 警告：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

在`build/windows/Taskfile.yml`中設定簽署：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

安全地儲存憑證密碼：

```bash
wails3 setup signing
```

如需詳細資訊，請參閱[簽署應用程式](/guides/build/signing/)。

## 針對 ARM 建置

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## 疑難排解

### 找不到 makensis

安裝 NSIS：

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### SmartScreen 警告

您的可執行檔尚未簽署。請參閱上方的[程式碼簽署](#heading-1)。

### 缺少 WebView2

安裝程式包含 WebView2 啟動載入器，會在需要時下載執行階段。若需要離線安裝，請從 Microsoft 下載 Evergreen Standalone Installer。
