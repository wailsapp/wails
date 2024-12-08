---
title: Packaging Your Application
sidebar:
  order: 30
---

This guide explains how to package your Wails application for different
platforms.

## Windows

Windows applications are packaged as `.exe` files. Wails automatically handles
this during the build process, creating a standalone executable that includes
all necessary resources.

## macOS

macOS applications are packaged as `.app` bundles. Wails creates these bundles
automatically during the build process, including proper code signing and
notarization if configured.

## Linux

Linux applications can be packaged in various formats. Wails v3 uses
[nfpm](https://github.com/goreleaser/nfpm), an excellent packaging tool that
makes it easy to create `.deb`, `.rpm`, and Arch Linux packages. nfpm is a
powerful tool that handles the complexities of Linux packaging, making it easy
to create professional-grade packages.

### Package Types

Wails supports creating the following types of Linux packages:

- Debian packages (`.deb`) - for Debian, Ubuntu, and related distributions
- Red Hat packages (`.rpm`) - for Red Hat, Fedora, CentOS, and related
  distributions
- Arch Linux packages - for Arch Linux and related distributions
- AppImage - a distribution-independent package format

### Building Packages

Wails provides several task commands for building Linux packages. These are
defined in `Taskfile.linux.yml` and can be invoked using the `wails3 task`
command:

```bash
# Build all package types (AppImage, deb, rpm, and Arch Linux)
wails3 task linux:package

# Build specific package types
wails3 task linux:create:appimage  # Create an AppImage
wails3 task linux:create:deb       # Create a Debian package
wails3 task linux:create:rpm       # Create a Red Hat package
wails3 task linux:create:aur       # Create an Arch Linux package
```

Each of these tasks will:

1. Build your application in production mode
2. Generate necessary desktop integration files
3. Create the appropriate package using nfpm

### Configuration

The package configuration file should follow the nfpm configuration format and
is typically located at `build/nfpm/nfpm.yaml`. Here's an example:

```yaml
name: "myapp"
arch: "amd64"
version: "v1.0.0"
maintainer: "Your Name <your.email@example.com>"
description: |
  A short description of your application
vendor: "Your Company"
homepage: "https://yourcompany.com"
license: "MIT"
contents:
  - src: ./build/bin/myapp
    dst: /usr/bin/myapp
  - src: ./assets/icon.png
    dst: /usr/share/icons/myapp.png
  - src: ./assets/myapp.desktop
    dst: /usr/share/applications/myapp.desktop
```

For detailed information about all available configuration options, please refer
to the
[nfpm configuration documentation](https://nfpm.goreleaser.com/configuration/).
