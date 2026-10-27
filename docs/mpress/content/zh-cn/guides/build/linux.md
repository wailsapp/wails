---
title: "Linux 打包"
description: "将 Wails 应用打包以便在 Linux 上分发"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## 软件包格式

将应用打包以便在 Linux 上分发：

```bash
wails3 package GOOS=linux
```

这会在`bin/`目录中创建多种格式的软件包：

- **AppImage**：可移植，可在任何 Linux 发行版上运行
- **DEB**：适用于 Debian、Ubuntu 及其衍生发行版
- **RPM**：适用于 Fedora、RHEL 及其衍生发行版
- **Arch**：适用于 Arch Linux 及其衍生发行版

### 单独构建各格式

构建指定格式：

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## 自定义软件包

### 桌面条目

`.desktop`文件控制应用在应用程序菜单中的显示方式。该文件根据`build/linux/Taskfile.yml`中的值生成：

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### 软件包元数据

编辑`build/linux/nfpm/nfpm.yaml`以自定义 DEB 和 RPM 软件包：

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

AppImage 配置位于`build/linux/appimage/`中。应用图标来自`build/appicon.png`。

AppImage 会捆绑 WebKit 及其辅助进程，并修补捆绑的库，使其从 AppImage 内部启动这些进程。其`AppRun`会从 AppImage 内的`usr/lib`启动应用；修补后的路径在 WebKit 的 bubblewrap 沙箱内外均可正确解析，因此沙箱会保持启用状态。仅当自定义的`AppRun`能够保留该工作目录时，才应提供它。

## 签名软件包

使用 PGP 密钥签名 DEB 和 RPM 软件包：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

在`build/linux/Taskfile.yml`中配置签名：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

存储密钥密码：

```bash
wails3 signing
```

有关详细信息，请参阅[应用签名](/guides/build/signing/)。

## 为 ARM 构建

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
在 x86_64 主机上进行 ARM64 构建时，使用 Docker 完成 CGO 交叉编译。

@end

## 旧版 GTK3 支持

默认情况下，Wails v3基于**GTK4 和 WebKitGTK 6.0**构建。对于尚未提供 WebKitGTK 6.0的发行版（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x），仍可使用旧版 GTK3 / WebKit2GTK 4.1构建路径。必须通过构建标签显式启用该旧版路径，并计划在v3.1中移除。

@note{type="caution" title="旧版路径"}
GTK3 / WebKit2GTK 4.1路径在v3.0.x 系列中受支持。请根据目标发行版对 GTK4 / WebKitGTK 6.0的支持情况规划迁移到 GTK4，因为`-tags gtk3`将在v3.1中移除。

@end

### 依赖项

安装 GTK3 和 WebKit2GTK 4.1开发库：

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

所需的 pkg-config 软件包为`gtk+-3.0`和`webkit2gtk-4.1`。

### 使用 GTK3 构建

使用`-tags gtk3`标志：

```bash
wails3 build -tags gtk3
```

也可以直接使用 Go：

```bash
go build -tags gtk3 -o myapp .
```

### 与 GTK4 的已知差异

- **文件对话框**：GTK4 默认使用`xdg-desktop-portal`显示文件对话框，因此部分对话框选项（例如默认目录、自定义筛选器的显示方式）的行为与 GTK3 不同。有关详细信息，请参阅[对话框参考 - Linux 对话框行为](/reference/dialogs/#linux-dialog-behavior)。
- **菜单样式**：GTK4 支持`LinuxMenuStylePrimaryMenu`选项，该选项遵循 GNOME HIG，在标题栏中显示汉堡菜单按钮（☰）。此选项对`-tags gtk3`构建无效。请参阅[窗口 API - Linux MenuStyle](/reference/window/#linux)。
- **DPI 缩放**：GTK4 使用`gdk_monitor_get_scale`（GTK 4.14+）支持非整数倍缩放。

### 检查构建环境

运行`wails3 doctor`以验证环境配置。不带任何标志时，它会检查默认使用的 GTK4 / WebKitGTK 6.0。旧版 GTK3 / WebKit2GTK 4.1软件包会列为可选项。

## 故障排除

### AppImage 无法运行

为其添加可执行权限：

```bash
chmod +x MyApp-x86_64.AppImage
```

### 缺少依赖项

请安装与二进制文件构建变体相匹配的库。对于默认的 GTK4/WebKitGTK 6.0构建：

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

对于明确使用`-tags gtk3`构建的二进制文件，请改为安装旧版 GTK3/WebKit2GTK 4.1库：

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

软件包的可用性及名称取决于发行版版本。请使用包含所需库且受支持的发行版，或使用文档中说明的交叉构建/容器构建路径；安装旧版库无法满足默认构建的依赖要求。AppImage 的可移植性取决于其打包内容实际捆绑了哪些组件，以及仍需使用哪些主机库。

### 未找到 C 编译器

构建系统需要 GCC 或 Clang 来使用 CGO：

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

也可以运行`wails3 task setup:docker`，构建系统将自动使用 Docker。

### 使用 NVIDIA GPU 时窗口空白或全白

在使用 NVIDIA 专有驱动程序的 Linux 系统上，Wails 应用启动时可能显示空白或全白窗口。这是由 WebKitGTK 的一个错误引起的：DMA-BUF 渲染器使用`gbm_bo_map()`时无法配合 NVIDIA 专有驱动程序正常工作（影响 X11 和 Wayland、377–580+ 版本的驱动程序、10系列及更早的 GT 710 GPU）。

**当 Wails 检测到 NVIDIA 内核模块（`/sys/module/nvidia`）时，会自动应用`WEBKIT_DISABLE_DMABUF_RENDERER=1`**，因此大多数用户无需执行任何操作。

如果仍然看到空白窗口（例如在无法看到模块路径的容器中），请在启动应用前手动设置环境变量：

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

相关上游错误：[WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607)、[WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739)。

### AppImage 的 strip 兼容性

在现代 Linux 发行版（Arch Linux、Fedora 39+、Ubuntu 24.04+）中，系统库使用`.relr.dyn` ELF 节进行编译，以提高重定位效率。用于创建 AppImage 的`linuxdeploy`工具捆绑了较旧的`strip`二进制文件，无法处理这些现代节。

Wails 会在构建 AppImage 前检查系统 GTK 库，从而自动检测这种情况。检测到后，会禁用剥离（`NO_STRIP=1`）以确保兼容性。

**这意味着：**

- 在受影响的系统上，AppImage 会略大（约20-40%）
- 应用功能不受影响
- 系统会自动处理，无需执行任何操作

如果需要在现代系统上生成更小的 AppImage，可以安装较新的`strip`二进制文件，并配置`linuxdeploy`使用该文件，而不是其捆绑的版本。
