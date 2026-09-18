---
title: "Installationsprogramme erstellen"
description: "Anwendung für die Verteilung paketieren"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## Übersicht

Erstellen Sie professionelle Installationsprogramme für Ihre Wails-Anwendung auf allen Plattformen.

## Plattformspezifische Installationsprogramme

@tabs{sync-key="platform"}
[Windows]
### NSIS-Installationsprogramm

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

Alternative für MSI-Installationsprogramme.

[macOS]
### DMG-Erstellung

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### Codesignierung

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

Verwenden Sie Xcode für die Verteilung über den App Store.

[Linux]
### DEB-Paket

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

### RPM-Paket

Verwenden Sie `rpmbuild` für RPM-basierte Distributionen.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## Automatisierte Paketierung

### GoReleaser verwenden

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- Code auf allen Plattformen signieren
- Versionsinformationen angeben
- Deinstallationsprogramme erstellen
- Installationsprozess testen
- Klare Dokumentation bereitstellen

### ❌ Zu vermeiden

- Codesignierung nicht überspringen
- Dateizuordnungen nicht vergessen
- Pfade nicht fest codieren
- Tests nicht überspringen

## Nächste Schritte

- [In-App-Updater](/guides/updater/) – Ermöglichen Sie Ihrer Anwendung, sich selbst zu aktualisieren
- [Plattformübergreifendes Erstellen](/guides/build/cross-platform/) – Erstellen Sie Builds für mehrere Plattformen
