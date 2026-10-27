---
title: "Criação de instaladores"
description: "Empacote seu aplicativo para distribuição"
slug: "guides/installers"
sourcePath: "guides/installers.md"
---

## Visão geral

Crie instaladores profissionais para seu aplicativo Wails em todas as plataformas.

## Instaladores por plataforma

@tabs{sync-key="platform"}
[Windows]
### Instalador NSIS

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

Alternativa para instaladores MSI.

[macOS]
### Criação de DMG

```bash
# Create DMG
hdiutil create -volname "MyApp" -srcfolder bin/MyApp.app -ov -format UDZO MyApp.dmg
```

### Assinatura de código

```bash
# Sign application
codesign --deep --force --verify --verbose --sign "Developer ID" MyApp.app

# Notarize
xcrun notarytool submit MyApp.dmg --apple-id "email" --password "app-password"
```

### App Store

Use o Xcode para distribuir pela App Store.

[Linux]
### Pacote DEB

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

### Pacote RPM

Use `rpmbuild` para distribuições baseadas em RPM.

### AppImage

```bash
# Use appimagetool
appimagetool myapp.AppDir
```

@end

## Empacotamento automatizado

### Uso do GoReleaser

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

## Práticas recomendadas

### ✅ Faça

- Assine o código em todas as plataformas
- Inclua informações de versão
- Crie desinstaladores
- Teste o processo de instalação
- Forneça documentação clara

### ❌ Não faça

- Não deixe de assinar o código
- Não se esqueça das associações de arquivos
- Não fixe caminhos diretamente no código
- Não deixe de testar

## Próximas etapas

- [Atualizador integrado ao aplicativo](/guides/updater/) — Adicione ao seu aplicativo a capacidade de atualizar a si mesmo
- [Compilação multiplataforma](/guides/build/cross-platform/) — Compile para várias plataformas
