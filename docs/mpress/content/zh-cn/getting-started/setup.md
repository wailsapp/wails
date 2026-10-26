---
title: "设置"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="实验性功能"}
设置向导是新增功能，目前主要在 Linux 上进行了测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并改为按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速开始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

向导会在浏览器中打开，引导你完成依赖项检查、项目默认设置和可选的跨平台构建设置。

之后，就可以创建你的第一个项目了：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能说明

- **检查依赖项** - 验证 Go、npm 和平台工具链
- **配置默认设置** - 作者信息、Bundle ID 前缀和首选模板
- **跨平台构建** - 可选的 Docker 设置，用于从任意主机进行构建
- **代码签名** - 适用于 macOS、Windows 和 Linux 的可选设置

配置会保存到`$XDG_CONFIG_HOME/wails/defaults.yaml`，并以`~/.config/wails/defaults.yaml`作为后备位置，供`wails3 init`使用。此配置与项目的`build/config.yml`相互独立。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到问题？

1. 运行`wails3 doctor`以诊断问题
2. 按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作
3. 提交[问题报告](https://github.com/wailsapp/wails/issues/4904)时，请附上`wails3 doctor`的输出
