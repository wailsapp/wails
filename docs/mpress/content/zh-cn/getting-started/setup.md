---
title: "设置"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="实验性功能"}
设置向导是一项新功能，目前主要在 Linux 上进行了测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并改为按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速开始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

该向导会在浏览器中打开，并引导你完成依赖项检查、项目默认设置和可选的跨平台构建设置。

然后就可以创建第一个项目了：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能

- **检查依赖项**——验证 Go、npm 和平台工具链
- **配置默认值**——作者信息、捆绑包 ID 前缀和首选模板
- **跨平台构建**——可选择配置 Docker，以便从任意主机进行构建
- **代码签名**——可选择为 macOS、Windows 和 Linux 进行设置

配置将保存到`~/.config/wails/config.yaml`，并由`wails3 init`使用。

## 子命令

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## 遇到问题？

1. 运行`wails3 doctor`以诊断问题
2. 按照[手动安装步骤](/getting-started/installation/#platform-specific-dependencies)操作
3. 提交[问题报告](https://github.com/wailsapp/wails/issues/4904)时，请附上`wails3 doctor`的输出
