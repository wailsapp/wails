---
title: "建立安裝程式"
description: "封裝應用程式以供散布"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## 概觀

為各平台上的 Wails 應用程式建立專業的安裝程式。

## 各平台安裝程式

@tabs{sync-key="platform"}
[Windows]
### NSIS 安裝程式

```bash
# Install NSIS
# Download from: https://nsis.sourceforge.io/

# Create installer script (installer.nsi)
makensis installer.nsi
```

**installer.nsi：**

```nsis
!define APPNAME "MyApp"
!define VERSION "1.0.0"

Name "${APPNAME}"
OutFile "MyApp-Setup.exe"
InstallDir "$PROGRAMFILES\${APPNAME}"

Section "Install"
    SetOutPath "$INSTDIR"
    File "build\bin\myapp.exe"
    CreateShortcut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\myapp.exe"
SectionEnd
```

### WiX Toolset

製作 MSI 安裝程式的另一種選擇。

[macOS]
### 建立 DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### 程式碼簽署

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

使用 Xcode 透過 App Store 散布。

[Linux]
### DEB 套件

```bash
# Create package structure
mkdir -p myapp_1.0.0/DEBIAN
mkdir -p myapp_1.0.0/usr/bin

# Copy binary
cp bin/myapp myapp_1.0.0/usr/bin/

# Create control file
cat > myapp_1.0.0/DEBIAN/control << EOF
Package: myapp
Version: 1.0.0
Architecture: amd64
Maintainer: Your Name
Description: My Application
EOF

# Build package
dpkg-deb --build myapp_1.0.0
```

### RPM 套件

對於採用 RPM 的發行版，請使用`rpmbuild`。

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## 自動化封裝

### 使用 GoReleaser

```yaml
# .goreleaser.yml
project_name: myapp

builds:
  - binary: myapp
    goos:
      - windows
      - darwin
      - linux
    goarch:
      - amd64
      - arm64

archives:
  - format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

nfpms:
  - formats:
      - deb
      - rpm
    vendor: Your Company
    homepage: https://example.com
    description: My Application
```

## 最佳實務

### ✅ 應該做

- 在所有平台上簽署程式碼
- 包含版本資訊
- 建立解除安裝程式
- 測試安裝流程
- 提供清楚的文件

### ❌ 不應該做

- 不要略過程式碼簽署
- 不要忘記檔案關聯
- 不要將路徑寫死
- 不要略過測試

## 後續步驟

- [應用程式內更新程式](/guides/updater/) - 為應用程式加入自動更新功能
- [跨平台建置](/guides/build/cross-platform/) - 為多個平台建置
