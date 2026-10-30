---
title: "入门指南"
description: "如何开始为 Wails 做贡献 v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 欢迎参与贡献！

感谢你有兴趣为 Wails 做贡献！本指南将帮助你完成第一次贡献。

## 前置条件

开始之前，请确保你已具备：

- 已安装**Go 1.26+**（[下载](https://go.dev/dl/)）
- 已安装**Node.js 20+**和**npm**（[下载](https://nodejs.org/)）
- 已使用你的 GitHub 账户配置**Git**
- 基本熟悉 Go 和 JavaScript/TypeScript

### 特定平台的要求

**macOS：**

- Xcode 命令行工具：`xcode-select --install`

**Windows：**

- 建议使用 MSYS2 或类似的类 Unix 环境
- WebView2 运行时（通常已预装在 Windows 11上）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（默认 GTK4 技术栈）
- 安装命令：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 对于旧版`-tags gtk3`构建路径，还需安装`libgtk-3-dev`和`libwebkit2gtk-4.1-dev`

## 贡献流程概览

典型的贡献工作流程包括以下步骤：

1. **复刻并克隆**——创建一份属于你自己的 Wails 仓库副本
2. **设置环境**——构建 Wails CLI 并验证你的环境
3. **创建分支**——为你的更改创建功能分支
4. **开发**——按照我们的编码规范进行更改
5. **测试**——运行测试以确保一切正常
6. **提交**——使用清晰且符合约定的提交消息提交更改
7. **发起请求**——创建拉取请求以供评审
8. **迭代**——回应反馈并进行调整
9. **合并**——获得批准后，你的更改便会成为 Wails 的一部分！

## 分步指南

选择你的贡献类型：

@tabs
[错误修复]
@steps
### 查找或报告错误
- 检查[GitHub Issues](https://github.com/wailsapp/wails/issues)中是否已报告该错误
- 如果尚未报告，请创建新议题并提供复现步骤
- 等待确认后再开始处理

### 复刻并克隆
在[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)复刻仓库

克隆你的复刻仓库：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 构建并验证
构建 Wails，并验证你能否复现该错误：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### 创建错误修复分支
为修复创建一个分支：

```bash
git checkout -b fix/issue-123-window-crash
```

### 修复错误
- 仅进行修复该错误所需的最小更改
- 不要重构无关代码
- 添加或更新测试以防止回归

```bash
# Make your changes
# Add tests in *_test.go files
```

### 测试你的修复
运行测试以确保修复有效：

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 提交你的修复
使用清晰的提交消息提交更改：

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### 提交拉取请求
推送更改并创建 PR：

```bash
git push origin fix/issue-123-window-crash
```

在 PR 描述中：

- 说明错误及其根本原因
- 说明你的修复方案
- 引用相关议题：“Fixes #123”
- 说明修复前后的行为

### 回应反馈
处理评审意见，并根据需要更新你的 PR。

@end

[WEP（增强功能）]
@steps
### 编写 WEP
- 阅读[WEP（Wails 增强提案）流程](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- 将 WEP 模板复制到`v3/wep/proposals/<name>/proposal.md`
- 创建一个标题为`[WEP] <title>`且仅包含 WEP 的草稿 PR
- 等待维护者作出决定后再实施

### 复刻并克隆
在[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)复刻仓库

克隆你的复刻仓库：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 设置开发环境
构建 Wails 并验证你的环境：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 创建功能分支
创建一个名称清晰的分支：

```bash
git checkout -b feat/window-transparency-support
```

### 实现功能
- 遵循我们的[编码规范](/contributing/standards/)
- 确保更改聚焦于该功能
- 编写简洁且有完善文档的代码
- 添加全面的测试

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 全面测试
测试你的功能：

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### 为你的功能编写文档
- 为所有公共 API 添加文档注释
- 更新`/docs/mpress/content/`中的相关文档
- 如适用，请添加示例

### 按约定提交
使用约定式提交：

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### 提交拉取请求
推送更改并创建 PR：

```bash
git push origin feat/window-transparency-support
```

在你的 PR 中：

- 说明该功能及其使用场景
- 提供示例或屏幕截图
- 列出所有破坏性更改
- 引用已获接受的 WEP PR

### 根据评审意见迭代
维护者可能会要求你进行修改。请耐心沟通并积极协作。

@end

[文档]
欢迎直接提交更正 PR，无需先创建议题。仅更正文档时，无需提供失败的代码测试。有关 M-Press 的安装、源文件路径、预览、验证及 PR 步骤，请参阅[使用 M-Press 编写文档](/contributing/documentation/)。

@end

## 查找可处理的议题

- 查找带有[`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)标签的议题
- 查看[`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)议题
- 浏览[开放的议题](https://github.com/wailsapp/wails/issues)，并请求将其分配给你

## 获取帮助

- <strong>Discord：</strong>加入[Wails Discord](https://discord.gg/JDdSxwjhGf)
- <strong>讨论区：</strong>在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)中发帖
- <strong>议题：</strong>对于可复现的错误，请创建议题；如有疑问，请使用 Discussions；如需提出增强功能，请提交 WEP PR

## 行为准则

请保持尊重、提供建设性意见，并友善地接纳他人。我们正在共同建设一个友好的社区，携手打造出色的软件。

## 后续步骤

- 设置你的[开发环境](/contributing/setup/)
- 查看我们的[编码规范](/contributing/standards/)
- 浏览[技术文档](/contributing/overview/)
