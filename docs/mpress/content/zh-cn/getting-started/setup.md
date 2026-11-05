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

该向导会在浏览器中打开，引导你完成依赖项检查、全局项目默认设置、可选的跨平台工具和签名凭据配置。这些首选项可在多个项目中复用。

之后即可开始创建你的第一个项目：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能说明

- **检查依赖项** - 验证 Go、npm 和平台工具
- **配置默认值** - 作者信息、捆绑包 ID 前缀、首选模板
- **跨平台构建** - 可选的 Docker 设置，用于从任意主机进行构建
- **代码签名** - macOS、Windows 和 Linux 的可选设置

配置会保存到`$XDG_CONFIG_HOME/wails/defaults.yaml`；如果不可用，则回退到`~/.config/wails/defaults.yaml`，并由`wails3 init`使用。此配置独立于项目的`build/config.yml`。

## 创建 Wake 项目

启用 Wake 实验性功能后，`wails3 init -ui`会收集项目专属的构建选项，并在项目创建成功时将其保存到`wails.yaml`。按照[Wake 入门教程](/experimental/wake/)操作，在创建桌面项目时选择构建后端、目标平台和软件包格式。`wails3 setup`仍用于管理全局首选项。

现有 Wake 项目在`wails.yaml`中配置。使用[YAML 构建指南](/experimental/yaml-builds/)更改其配置文件和打包设置。现有 Taskfile 项目必须先完成[经过审核的单向迁移](/experimental/wake/#migrate-an-existing-project)，才能使用 Wake；设置过程不会转换或执行其 Taskfile。

## 复用签名首选项

已签名的 Wake 构建会从用户级默认值中补充缺失的签名设置，包括项目外部的证书文件。项目中的显式设置优先。私钥、密码和全局凭据路径不会复制到`wails.yaml`中。

对于全局配置的 Windows 和 Linux 密钥，请通过现有凭据存储或`WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`提供密码。对于加密的 Apple 身份，请将`darwin.signing.credential`设置为包含其密码的环境变量名称。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到问题？

1. 运行`wails3 doctor`以诊断问题
2. 按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作
3. 提交[问题报告](https://github.com/wailsapp/wails/issues/4904)，并附上`wails3 doctor`输出
