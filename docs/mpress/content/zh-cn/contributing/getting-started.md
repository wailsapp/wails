---
title: "入门指南"
description: "如何开始为 Wails v3 做贡献"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 欢迎，贡献者！

感谢你有兴趣为 Wails 做贡献！本指南将帮助你完成第一次贡献。

## 前提条件

开始之前，请确保你具备：

- 已安装 **Go 1.25+**（[下载](https://go.dev/dl/)）
- **Node.js 20+** 和 **npm**（[下载](https://nodejs.org/)）
- 已使用你的 GitHub 账户配置 **Git**
- 基本熟悉 Go 和 JavaScript/TypeScript

### 特定平台要求

**macOS：**

- Xcode 命令行工具：`xcode-select --install`

**Windows：**

- 建议使用 MSYS2 或类似的类 Unix 环境
- WebView2 运行时（Windows 11上通常已预安装）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（默认 GTK4 技术栈）
- 安装方式：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 对于旧版 `-tags gtk3` 构建路径，还需安装 `libgtk-3-dev` 和 `libwebkit2gtk-4.1-dev`

## 贡献流程概览

典型的贡献工作流包括以下步骤：

1. **派生并克隆**——创建你自己的 Wails 仓库副本
2. **配置环境**——构建 Wails CLI 并验证你的环境
3. **创建分支**——为你的更改创建功能分支
4. **开发**——按照我们的编码规范进行更改
5. **测试**——运行测试以确保一切正常
6. **提交**——使用清晰、符合约定的提交消息进行提交
7. **提交审核**——创建拉取请求以供审核
8. **迭代**——回应反馈并进行调整
9. **合并**——获得批准后，你的更改将成为 Wails 的一部分！

## 分步指南

选择你的贡献类型：

@tabs
[错误修复]
@steps
### 查找或报告错误
- 检查该错误是否已在 [GitHub Issues](https://github.com/wailsapp/wails/issues) 中报告
- 如果尚未报告，请创建新议题并提供复现步骤
- 开始工作前，请等待确认

### 派生并克隆
在 [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) 派生仓库

克隆你的派生仓库：

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
为修复创建分支：

```bash
git checkout -b fix/issue-123-window-crash
```

### 修复错误
- 仅进行修复该错误所需的最少更改
- 不要重构无关代码
- 添加或更新测试以防止回归

```bash
# Make your changes
# Add tests in *_test.go files
```

### 测试修复
运行测试以确保修复有效：

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 提交修复
使用清晰的消息提交：

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### 提交拉取请求
推送并创建 PR：

```bash
git push origin fix/issue-123-window-crash
```

在 PR 描述中：

- 说明错误及其根本原因
- 描述你的修复方案
- 引用议题：“Fixes #123”
- 说明修复前后的行为

### 回应反馈
处理审核意见，并根据需要更新你的 PR。

@end

[WEP（增强功能）]
@steps
### 编写 WEP
- 阅读 [WEP（Wails 增强提案）流程](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)
- 将 WEP 模板复制到 `v3/wep/proposals/<name>/proposal.md`
- 创建一个标题为 `[WEP] <title>` 且仅包含 WEP 的草稿 PR
- 实施前，请等待维护者作出决定

### 派生并克隆
在 [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) 派生仓库

克隆你的派生仓库：

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
- 编写简洁且有文档说明的代码
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

### 为功能编写文档
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
推送并创建 PR：

```bash
git push origin feat/window-transparency-support
```

在 PR 中：

- 描述该功能及其使用场景
- 提供示例或截图
- 列出所有破坏性变更
- 引用已获接受的 WEP PR

### 根据评审意见迭代
维护者可能会要求你进行更改。请耐心沟通、积极协作。

@end

[文档]
欢迎直接提交更正 PR，无需先创建议题。仅涉及文档的更正不需要提供会失败的代码测试。有关 M-Press 的安装、源文件路径、预览、验证和 PR 步骤，请遵循[更正文档](/contributing/documentation/)指南。

@end

## 查找可处理的议题

- 查找[`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)标签
- 查看[`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)议题
- 浏览[开放议题](https://github.com/wailsapp/wails/issues)并请求分配给你

## 获取帮助

- <strong>Discord：</strong>加入[Wails Discord](https://discord.gg/JDdSxwjhGf)
- <strong>讨论区：</strong>在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)中发帖
- <strong>议题：</strong>对于可复现的错误，请创建议题；如有疑问，请使用 Discussions；如需提出增强功能，请提交 WEP PR

## 行为准则

请尊重他人、提供建设性意见并热情友善。我们致力于打造一个友好的社区，共同创造出色的软件。

## 后续步骤

- 设置你的[开发环境](/contributing/setup/)
- 查看我们的[编码规范](/contributing/standards/)
- 浏览[技术文档](/contributing/overview/)
