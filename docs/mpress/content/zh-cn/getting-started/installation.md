---
title: "安装"
description: "安装 Wails 并设置开发环境"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 支持的平台

- Windows AMD64/ARM64
- macOS 10.15+ AMD64（可部署至 macOS 10.13+）
- macOS 11.0+ ARM64
- Ubuntu 24.04 AMD64/ARM64（其他 Linux 发行版也可能可用！）

## 依赖项

安装 Wails 前，需要先安装一些通用依赖项。

@note{type="tip"}
安装 Wails CLI 后，可以运行`wails3 setup`自动检查这些依赖项并获取安装帮助。

@end

@tabs
[Go（至少为 1.24）]
从[Go 下载页面](https://go.dev/dl/)下载 Go。

请务必遵循官方的[Go 安装说明](https://go.dev/doc/install)。还需要确保`PATH`环境变量包含`~/go/bin`目录的路径。重启终端并执行以下检查：

- 检查 Go 是否已正确安装：`go version`
- 检查`~/go/bin`是否位于 PATH 变量中
  - Mac / Linux：`echo $PATH | grep go/bin`
  - Windows：`$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm（可选）]
虽然 Wails 不要求安装 npm，但大多数内置模板都需要它。

从[Node 下载页面](https://nodejs.org/en/download/)下载最新的 node 安装程序。最好使用最新版本，因为我们通常以该版本进行测试。

运行`npm --version`进行验证。

@note{type="info"}
如果更愿意使用其他包管理器而非 npm，可以自行选择。需要更新项目的 Taskfile 以使用该包管理器。

@end

@end

## 平台特定依赖项

还需要安装平台特定的依赖项：

@tabs{sync-key="platform"}
[Mac]
Wails 要求安装 xcode 命令行工具。可以运行以下命令进行安装：

```sh
xcode-select --install
```

[Windows]
Wails 要求安装[WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)。几乎所有 Windows 系统都已安装该运行时。可以使用`wails doctor`命令进行检查。

[Linux]
Linux 需要标准的`gcc`构建工具，以及`gtk4`和`webkitgtk-6.0`。安装后运行<code>wails3 doctor</code>，即可查看依赖项的安装方法。旧版 GTK3 / WebKit2GTK 4.1技术栈仍可通过`-tags gtk3`使用（请参阅[Linux 打包 - 旧版 GTK3 支持](/guides/build/linux/#legacy-gtk3-support)），直至 v3.1。如果不支持你的发行版或包管理器，请在 discord 上告知我们。

@end

## 安装

要使用 Go Modules 安装 Wails CLI，请运行以下命令：

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

如果要安装最新的开发版本，请运行以下命令：

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

使用开发版本时，所有生成的项目都会使用 Go 的[replace](https://go.dev/ref/mod#go-mod-file-replace)指令，以确保项目使用 Wails 的开发版本。

## 后续步骤

安装 CLI 后，运行设置向导来配置开发环境：

```shell
wails3 setup
```

@note{type="caution" title="实验性功能"}
设置向导是一项新功能，目前主要在 Linux 上进行了测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并按照下方的手动依赖项安装步骤操作。

@end

设置向导将执行以下操作：

- 检查平台依赖项并协助安装
- 配置项目默认值（作者信息、Bundle ID 前缀）
- 可选择设置 Docker 以进行跨平台构建
- 配置代码签名（如有需要）

有关更多详细信息，请参阅[设置指南](/getting-started/setup/)。

## 手动安装依赖项

如果更愿意手动安装依赖项，或者设置向导无法在你的系统上运行，请按照上方针对各平台的说明操作，然后运行：

```shell
wails3 doctor
```

此命令将检查是否已安装正确的依赖项，并提示缺少哪些依赖项。

## 找不到`wails3`命令？

如果系统报告找不到`wails3`命令，请检查以下事项：

- 确保已正确按照上方的<strong>Go安装指南</strong>操作，并且`go/bin`目录位于`PATH`环境变量中。
- 关闭并重新打开当前终端，以加载新的`PATH`变量。
