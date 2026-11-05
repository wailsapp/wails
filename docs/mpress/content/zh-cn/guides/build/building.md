---
title: "构建应用程序"
description: "构建并打包 Wails 应用程序"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 使用[Task](https://taskfile.dev)作为构建系统。`wails3 build`和`wails3 package`命令是对 Task 的便捷封装。

## 构建

为当前平台构建：

```bash
wails3 build
```

为指定平台构建：

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

输出将写入`bin/`目录。

@note{type="tip"}
从其他平台交叉编译到 macOS 或 Linux 需要 Docker。有关设置方法，请参阅[跨平台构建](/guides/build/cross-platform/)。

@end

## 开发

运行应用程序并启用热重载：

```bash
wails3 dev
```

这会启动一个文件监视器，在文件发生更改时重新构建并重启应用程序。前端开发服务器默认在端口9245上运行。

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## 打包

打包应用程序以供分发：

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

这会创建特定于平台的软件包：

- **Windows**：NSIS 安装程序——请参阅[Windows 打包](/guides/build/windows/)
- **macOS**：应用程序包（`.app`）——请参阅[macOS 打包](/guides/build/macos/)
- **Linux**：AppImage、deb 和 rpm——请参阅[Linux 打包](/guides/build/linux/)

## 自定义构建标签

使用`-tags`标志传递自定义 Go 构建标签：

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

标签将以`EXTRA_TAGS`的形式转发给底层 Taskfile。有关详细信息，请参阅[服务器构建](/guides/server-build/)和[Linux 打包 - 旧版 GTK3 支持](/guides/build/linux/#legacy-gtk3-support)。

## 直接使用 Task

如需更多控制，请直接使用 Task：

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

`linux:create:deb`或`darwin:build:universal`等特定于平台的任务只能通过 Task 使用。

## 生成资源

重新生成图标或更新构建配置：

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
