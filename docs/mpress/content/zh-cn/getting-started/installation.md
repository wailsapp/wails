---
title: "安装"
description: "安装 Wails 并设置开发环境"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 支持的平台

- Windows AMD64/ARM64
- macOS AMD64/ARM64。生成的构建任务以 macOS 12.0 或更高版本为目标；这是标准构建的部署目标，并不表示每种自定义构建配置都支持较旧的系统。
- Ubuntu 24.04 AMD64/ARM64（其他 Linux 发行版也可能适用！）

## 依赖项

安装 Wails 前需要满足一些通用依赖项。

@note{type="tip"}
安装 Wails CLI 后，可以运行 `wails3 setup` 来自动检查这些依赖项并协助安装。

@end

@tabs
[Go（至少为 1.26）]
Wails v3 选择 Go 1.26 或更高版本的工具链。较旧的 Go 启动器可能会通过 `GOTOOLCHAIN=auto` 自动下载该工具链；否则请显式安装。入门模块中较低版本的 Go 声明不会覆盖 Wails 依赖项的要求。

从[Go 下载页面](https://go.dev/dl/)下载 Go。

请务必遵循官方的[Go 安装说明](https://go.dev/doc/install)。还需确保 `PATH` 环境变量中包含 `~/go/bin` 目录的路径。重启终端，然后执行以下检查：

- 检查 Go 是否已正确安装：`go version`
- 检查 `~/go/bin` 是否在 PATH 变量中
  - Mac / Linux：`echo $PATH | grep go/bin`
  - Windows：`$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm（可选）]
尽管 Wails 不要求安装 npm，但大多数内置模板都需要使用它。

从[Node 下载页面](https://nodejs.org/en/download/)下载最新的 Node 安装程序。最好使用最新版本，因为我们通常针对该版本进行测试。

运行 `npm --version` 进行验证。

@note{type="info"}
如果你更喜欢 npm 以外的其他包管理器，可以自行使用。你需要更新项目的 Taskfile 以改用该包管理器。

@end

@end

## 平台特定依赖项

还需要安装平台特定的依赖项：

@tabs{sync-key="platform"}
[Mac]
Wails 要求安装 Xcode 命令行工具。可通过运行以下命令来安装：

```sh
xcode-select --install
```

[Windows]
Wails 要求安装 [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)。几乎所有 Windows 系统都已安装它。可以使用 `wails doctor` 命令进行检查。

[Linux]
Linux 需要标准的 `gcc` 构建工具，以及 `gtk4` 和 `webkitgtk-6.0`。安装后运行 <code>wails3 doctor</code>，即可查看依赖项的安装方法。在 v3.1 之前，仍可通过 `-tags gtk3` 使用旧版 GTK3 / WebKit2GTK 4.1 技术栈（请参阅[Linux 打包——旧版 GTK3 支持](/guides/build/linux/#legacy-gtk3-support)）。如果你的发行版或包管理器不受支持，请在 Discord 上告知我们。

@end

## 安装

要使用 Go Modules 安装 Wails CLI，请运行以下命令：

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

如果想安装最新的开发版本，请运行以下命令：

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

使用开发版本时，所有生成的项目都会使用 Go 的 [replace](https://go.dev/ref/mod#go-mod-file-replace) 指令，以确保项目使用 Wails 的开发版本。

## 后续步骤

安装 CLI 后，运行设置向导来配置开发环境：

```shell
wails3 setup
```

@note{type="caution" title="实验性功能"}
设置向导是一项新功能，目前主要在 Linux 上经过测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并按照下方的依赖项手动安装步骤操作。

@end

设置向导将：

- 检查平台依赖项并协助安装
- 配置项目默认值（作者信息、捆绑包 ID 前缀）
- 可选择为跨平台构建设置 Docker
- 配置代码签名（如有需要）

有关更多详细信息，请参阅[设置指南](/getting-started/setup/)。

## 手动安装依赖项

如果你希望手动安装依赖项，或者设置向导无法在你的系统上运行，请先按照上方适用于相应平台的说明操作，然后运行：

```shell
wails3 doctor
```

这将检查是否安装了正确的依赖项，并提示缺少哪些依赖项。

## 找不到 `wails3` 命令？

如果系统报告找不到 `wails3` 命令，请检查以下事项：

- 确保已正确遵循上方的**Go 安装指南**，并且 `go/bin` 目录已包含在 `PATH` 环境变量中。
- 关闭并重新打开当前终端，使新的 `PATH` 变量生效。
