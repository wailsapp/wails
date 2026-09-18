---
title: "创建安装程序"
description: "将应用程序打包以供分发"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## 概述

为所有平台上的 Wails 应用程序创建专业的安装程序。

## 各平台安装程序

@tabs{sync-key="platform"}
[Windows]
### NSIS 安装程序

```bash
# Install NSIS
# Download from: https://nsis.sourceforge.io/

# Create installer script (installer.nsi)
makensis installer.nsi
```

**installer.nsi:**

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

用于创建 MSI 安装程序的替代方案。

[macOS]
### 创建 DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### 代码签名

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

使用 Xcode 在 App Store 上分发。

[Linux]
### DEB 软件包

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

### RPM 软件包

对于基于 RPM 的发行版，请使用`rpmbuild`。

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## 自动打包

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

## 最佳实践

### ✅ 应做

- 在所有平台上进行代码签名
- 包含版本信息
- 创建卸载程序
- 测试安装流程
- 提供清晰的文档

### ❌ 不应做

- 不要跳过代码签名
- 不要忘记配置文件关联
- 不要硬编码路径
- 不要跳过测试

## 后续步骤

- [应用内更新程序](/guides/updater/) - 为应用添加自我更新功能
- [跨平台构建](/guides/build/cross-platform/) - 为多个平台构建
