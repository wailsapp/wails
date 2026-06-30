---
title: Linux Packaging
description: Package your Wails application for Linux distribution
sidebar:
  order: 5
---

import { Aside } from '@astrojs/starlight/components';

## Package Formats

Package your app for Linux distribution:

```bash
wails3 package GOOS=linux
```

This creates multiple formats in the `bin/` directory:
- **AppImage**: Portable, runs on any Linux distribution
- **DEB**: For Debian, Ubuntu, and derivatives
- **RPM**: For Fedora, RHEL, and derivatives
- **Arch**: For Arch Linux and derivatives

### Individual Formats

Build specific formats:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Customizing Packages

### Desktop Entry

The `.desktop` file controls how your app appears in application menus. It's generated from values in `build/linux/Taskfile.yml`:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Package Metadata

Edit `build/linux/nfpm/nfpm.yaml` to customize DEB and RPM packages:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

AppImage configuration is in `build/linux/appimage/`. The app icon comes from `build/appicon.png`.

## Signing Packages

Sign DEB and RPM packages with a PGP key:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Configure signing in `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Store your key password:

```bash
wails3 setup signing
```

See [Signing Applications](/guides/build/signing) for details.

## Building for ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

<Aside type="note">
ARM64 builds from x86_64 hosts use Docker for CGO cross-compilation.
</Aside>

## Legacy GTK3 Support

Wails v3 builds on **GTK4 with WebKitGTK 6.0** by default. A legacy GTK3 / WebKit2GTK 4.1 path is still available for distributions that don't yet ship WebKitGTK 6.0 (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x). The legacy path is opt-in via a build tag and is scheduled for removal in v3.1.

<Aside type="caution" title="Legacy path">
The GTK3 / WebKit2GTK 4.1 path is supported through the v3.0.x line. Plan to migrate to GTK4 in coordination with your target distribution's GTK4 / WebKitGTK 6.0 availability — `-tags gtk3` will be removed in v3.1.
</Aside>

### Dependencies

Install GTK3 and WebKit2GTK 4.1 development libraries:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

The required pkg-config packages are `gtk+-3.0` and `webkit2gtk-4.1`.

### Building with GTK3

Use the `-tags gtk3` flag:

```bash
wails3 build -tags gtk3
```

Or directly with Go:

```bash
go build -tags gtk3 -o myapp .
```

### Known Differences from GTK4

- **File dialogs**: GTK4 uses `xdg-desktop-portal` for file dialogs (the default), which means some dialog options (like default directory, custom filters display) behave differently from GTK3. See [Dialogs Reference - Linux Dialog Behavior](/reference/dialogs#linux-dialog-behavior) for details.
- **Menu style**: GTK4 supports a `LinuxMenuStylePrimaryMenu` option that displays a hamburger button (☰) in the header bar, following GNOME HIG. This option has no effect on `-tags gtk3` builds. See [Window API - Linux MenuStyle](/reference/window#linux).
- **DPI scaling**: GTK4 uses `gdk_monitor_get_scale` (GTK 4.14+) for fractional scaling support.

### Checking Your Build

Run `wails3 doctor` to verify your setup. With no flags it checks for GTK4 / WebKitGTK 6.0 (the default). The legacy GTK3 / WebKit2GTK 4.1 packages are listed as optional.

## Troubleshooting

### AppImage won't run

Make it executable:

```bash
chmod +x MyApp-x86_64.AppImage
```

### Missing dependencies

If the app fails to start, check for missing WebKit dependencies:

```bash
# Debian/Ubuntu
sudo apt install libwebkit2gtk-4.1-0

# Fedora
sudo dnf install webkit2gtk4.1

# Arch
sudo pacman -S webkit2gtk-4.1
```

### No C compiler found

The build system needs GCC or Clang for CGO:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Alternatively, run `wails3 task setup:docker` and the build system will use Docker automatically.

### Blank or white window on NVIDIA GPU

On Linux with NVIDIA proprietary drivers, Wails apps may show a blank or white window on startup. This is caused by a WebKitGTK bug where the DMA-BUF renderer fails with `gbm_bo_map()` against the NVIDIA proprietary driver (affects X11 and Wayland, driver versions 377–580+, GPUs from the 10 series and older GT 710).

**Wails applies `WEBKIT_DISABLE_DMABUF_RENDERER=1` automatically** when it detects the NVIDIA kernel module (`/sys/module/nvidia`), so most users will not need to do anything.

If you still see a blank window (for example in a container where the module path is not visible), set the environment variable manually before launching your app:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Related upstream bugs: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### AppImage strip compatibility

On modern Linux distributions (Arch Linux, Fedora 39+, Ubuntu 24.04+), system libraries are compiled with `.relr.dyn` ELF sections for more efficient relocations. The `linuxdeploy` tool used to create AppImages bundles an older `strip` binary that cannot process these modern sections.

Wails automatically detects this situation by checking system GTK libraries before building the AppImage. When detected, stripping is disabled (`NO_STRIP=1`) to ensure compatibility.

**What this means:**
- AppImages will be slightly larger (~20-40%) on affected systems
- The application functionality is not affected
- This is handled automatically—no action required

If you need smaller AppImages on modern systems, you can install a newer `strip` binary and configure `linuxdeploy` to use it instead of its bundled version.
