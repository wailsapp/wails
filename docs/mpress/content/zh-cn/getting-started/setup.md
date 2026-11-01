---
title: "设置"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="实验性功能"}
设置向导是一项新功能，目前主要在 Linux 上进行了测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并改用[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)。

@end

## 快速开始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

向导会在浏览器中打开，引导你检查依赖项、设置项目默认值，并选择是否配置跨平台构建。上述说明使用默认设置流程；实验性的 Wake 项目流程将在下文介绍。

然后就可以创建第一个项目了：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能说明

- **检查依赖项** - 验证 Go、npm 和平台工具链
- **配置默认值** - 作者信息、Bundle ID 前缀和首选模板
- **跨平台构建** - 可选的 Docker 设置，用于从任意宿主系统构建
- **代码签名** - 适用于 macOS、Windows 和 Linux 的可选设置

配置将保存到`$XDG_CONFIG_HOME/wails/defaults.yaml`，并以`~/.config/wails/defaults.yaml`作为后备位置，供`wails3 init`使用。此配置独立于项目的`build/config.yml`。

## 配置 Wake 项目

在项目目录中启用[Wake 实验性功能](/experimental/yaml-builds/#enable-yaml-builds)并运行设置：

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

在 PowerShell 中：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

这会打开以项目为中心的流程：

1. **项目** — 如果不存在原生清单，则输入项目元数据。
2. **应用程序** — 选择 GUI、CLI 或服务器、构建后端，或者已构建的可执行文件。纯 Go 是一个需要明确选择的实验性选项。
3. **目标** — 选择操作系统和体系结构。iOS 和 Android 仍处于实验阶段；请在 YAML 中配置其签名设置。
4. **软件包** — 选择软件包格式以及是否对其签名。
5. **工具** — 查看所选构建和打包路径的要求。
6. **签名** — 如果选择了签名，则使用现有签名控件进行配置。
7. **保存** — 写入项目配置并显示要运行的命令。

工具要求取决于你的选择。纯 Go 桌面编译不需要 Docker 或 C 编译器。CGo 交叉构建需要兼容的编译器和目标库，或者受支持的容器工具链。NSIS 仍需要`makensis`；跨宿主系统构建 GUI AppImage 时，需要提供 Linux 运行时目录树。导出的签名身份可以使用 Wails 的 Go 签名器，而原生钥匙串或证书存储中的身份仍需使用平台工具。请参阅[打包和签名要求](/experimental/yaml-builds/#package-and-sign-for-distribution)。

使用已构建的可执行文件会跳过 Go 编译、前端安装和构建、绑定生成以及 UPX。仍然需要所选的打包和签名工具，以及有效的输入二进制文件。请参阅[预构建输入路径](/experimental/yaml-builds/#use-prebuilt-desktop-executables)。此选项仅适用于保存后显示的命令，并不是永久的 YAML 构建设置。

保存操作会写入`wails.yaml`，更新应用程序类型、前端启用状态和所选后端设置，并替换专用的`profiles.setup`目标、软件包和签名选项。其他配置文件、前端命令、配套可执行文件、钩子和不相关的项目设置均会保留。后端设置适用于整个项目，因此继承这些设置的其他配置文件也会使用所选后端。新项目的元数据会一直保留在向导中，直到执行保存操作。系统会报告缺失的工具，以便你在构建前安装；保存有效配置不会安装工具或运行构建。

签名构建会从用户级默认值中沿用项目缺少的签名设置，包括项目外部的证书文件。显式的项目设置优先。设置过程不会将私钥、密码或全局凭据路径复制到`wails.yaml`。对于全局配置的 Windows 和 Linux 密钥，请通过现有凭据存储或`WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`提供密码。对于加密的 Apple 身份，请将`darwin.signing.credential`设置为包含其密码的环境变量名称。

保存的配置文件通过`wails3 build setup`运行；选择已构建的可执行文件时还会添加`--skip-build`。

运行显示的命令时，请保持启用 Wake。如果未设置该标志，启用的`wails.yaml`绝不会回退到 Taskfile。设置过程不会转换现有 Taskfile 或旧版元数据：使用 Wake 配置这些项目之前，请先执行[已审核的迁移流程](/experimental/yaml-builds/#migrate-an-existing-project)。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到问题？

1. 运行`wails3 doctor`以诊断问题
2. 按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作
3. 提交[问题报告](https://github.com/wailsapp/wails/issues/4904)，并附上`wails3 doctor`的输出
