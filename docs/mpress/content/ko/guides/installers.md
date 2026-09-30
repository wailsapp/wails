---
title: "인스톨러 만들기"
description: "배포할 애플리케이션 패키징하기"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## 개요

모든 플랫폼에서 Wails 애플리케이션을 위한 전문적인 인스톨러를 만드세요.

## 플랫폼별 인스톨러

@tabs{sync-key="platform"}
[Windows]
### NSIS 인스톨러

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

MSI 인스톨러를 위한 대안입니다.

[macOS]
### DMG 만들기

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### 코드 서명

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

App Store 배포에는 Xcode를 사용하세요.

[Linux]
### DEB 패키지

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

### RPM 패키지

RPM 기반 배포판에서는 `rpmbuild`을 사용하세요.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## 자동 패키징

### GoReleaser 사용하기

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

## 모범 사례

### ✅ 권장 사항

- 모든 플랫폼에서 코드 서명하기
- 버전 정보 포함하기
- 제거 프로그램 만들기
- 설치 과정 테스트하기
- 명확한 문서 제공하기

### ❌ 금지 사항

- 코드 서명을 생략하지 마세요
- 파일 연결을 빠뜨리지 마세요
- 경로를 하드코딩하지 마세요
- 테스트를 생략하지 마세요

## 다음 단계

- [인앱 업데이터](/guides/updater/) - 앱에 자체 업데이트 기능 추가하기
- [크로스 플랫폼 빌드](/guides/build/cross-platform/) - 여러 플랫폼용으로 빌드하기
