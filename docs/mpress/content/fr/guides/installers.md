---
title: "Création de programmes d’installation"
description: "Empaquetez votre application pour la distribuer"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## Vue d’ensemble

Créez des programmes d’installation professionnels pour votre application Wails sur toutes les plateformes.

## Programmes d’installation par plateforme

@tabs{sync-key="platform"}
[Windows]
### Programme d’installation NSIS

```bash
# Install NSIS
# Download from: https://nsis.sourceforge.io/

# Create installer script (installer.nsi)
makensis installer.nsi
```

**installer.nsi :**

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

### Ensemble d’outils WiX

Alternative pour les programmes d’installation MSI.

[macOS]
### Création d’un DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### Signature du code

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

Utilisez Xcode pour la distribution sur l’App Store.

[Linux]
### Paquet DEB

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

### Paquet RPM

Utilisez `rpmbuild` pour les distributions basées sur RPM.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## Création automatisée de paquets

### Utilisation de GoReleaser

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

## Bonnes pratiques

### ✅ À faire

- Signez le code sur toutes les plateformes
- Incluez les informations de version
- Créez des programmes de désinstallation
- Testez le processus d’installation
- Fournissez une documentation claire

### ❌ À ne pas faire

- Ne négligez pas la signature du code
- N’oubliez pas les associations de fichiers
- Ne codez pas les chemins en dur
- Ne négligez pas les tests

## Étapes suivantes

- [Mise à jour intégrée](/guides/updater/) — Permettez à votre application de se mettre à jour elle-même
- [Compilation multiplateforme](/guides/build/cross-platform/) — Compilez pour plusieurs plateformes
