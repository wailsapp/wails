---
title: "Создание установщиков"
description: "Упаковка приложения для распространения"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## Обзор

Создавайте профессиональные установщики для приложения Wails на всех платформах.

## Установщики для разных платформ

@tabs{sync-key="platform"}
[Windows]
### Установщик NSIS

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

### Набор инструментов WiX

Альтернативный вариант для установщиков MSI.

[macOS]
### Создание DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### Подписание кода

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

Используйте Xcode для распространения через App Store.

[Linux]
### Пакет DEB

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

### Пакет RPM

Используйте `rpmbuild` для дистрибутивов на базе RPM.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## Автоматизированная упаковка

### Использование GoReleaser

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

## Рекомендации

### ✅ Следует

- Подписывайте код на всех платформах
- Добавляйте информацию о версии
- Создавайте средства удаления
- Тестируйте процесс установки
- Предоставляйте понятную документацию

### ❌ Не следует

- Не пропускайте подписание кода
- Не забывайте об ассоциациях файлов
- Не задавайте пути жёстко в коде
- Не пропускайте тестирование

## Дальнейшие действия

- [Встроенное средство обновления](/guides/updater/) — добавьте в приложение возможность обновлять себя
- [Кроссплатформенная сборка](/guides/build/cross-platform/) — выполняйте сборку для нескольких платформ
