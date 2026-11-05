---
title: "インストーラーの作成"
description: "配布用にアプリケーションをパッケージ化する"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## 概要

すべてのプラットフォーム向けに、Wailsアプリケーションの本格的なインストーラーを作成します。

## プラットフォーム別インストーラー

@tabs{sync-key="platform"}
[Windows]
### NSISインストーラー

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

MSIインストーラーを作成するための別のツールです。

[macOS]
### DMGの作成

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### コード署名

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

App Storeで配布するにはXcodeを使用します。

[Linux]
### DEBパッケージ

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

### RPMパッケージ

RPMベースのディストリビューションには`rpmbuild`を使用します。

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## パッケージ化の自動化

### GoReleaserの使用

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

## ベストプラクティス

### ✅ 推奨事項

- すべてのプラットフォームでコード署名を行う
- バージョン情報を含める
- アンインストーラーを作成する
- インストール手順をテストする
- 明確なドキュメントを提供する

### ❌ 禁止事項

- コード署名を省略しない
- ファイルの関連付けを忘れない
- パスをハードコードしない
- テストを省略しない

## 次のステップ

- [アプリ内アップデーター](/guides/updater/) - アプリに自己更新機能を追加する
- [クロスプラットフォームビルド](/guides/build/cross-platform/) - 複数のプラットフォーム向けにビルドする
