---
title: "Membuat Penginstal"
description: "Kemas aplikasi Anda untuk didistribusikan"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## Gambaran Umum

Buat penginstal profesional untuk aplikasi Wails Anda di semua platform.

## Penginstal Platform

@tabs{sync-key="platform"}
[Windows]
### Penginstal NSIS

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

Alternatif untuk penginstal MSI.

[macOS]
### Pembuatan DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### Penandatanganan Kode

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

Gunakan Xcode untuk distribusi melalui App Store.

[Linux]
### Paket DEB

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

### Paket RPM

Gunakan `rpmbuild` untuk distribusi berbasis RPM.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## Pengemasan Otomatis

### Menggunakan GoReleaser

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

## Praktik Terbaik

### ✅ Lakukan

- Tandatangani kode di semua platform
- Sertakan informasi versi
- Buat program penghapus instalasi
- Uji proses instalasi
- Sediakan dokumentasi yang jelas

### ❌ Jangan Lakukan

- Jangan lewatkan penandatanganan kode
- Jangan lupakan asosiasi file
- Jangan menulis jalur secara langsung dalam kode
- Jangan lewatkan pengujian

## Langkah Berikutnya

- [Pembaruan Dalam Aplikasi](/guides/updater/) - Tambahkan kemampuan agar aplikasi Anda dapat memperbarui dirinya sendiri
- [Build Lintas Platform](/guides/build/cross-platform/) - Lakukan build untuk beberapa platform
