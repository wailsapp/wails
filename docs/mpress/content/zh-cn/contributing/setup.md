---
title: "开发环境设置"
description: "设置用于 Wails v3 开发的开发环境"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## 开发环境设置

本指南将指导你搭建用于开发 Wails v3 的完整开发环境。

## 必备工具

### Go 开发环境

1. **安装 Go 1.25或更高版本：**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **配置 Go 环境：**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **安装实用的 Go 工具：**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js 和 npm

仅前端集成示例需要。

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### 平台特定依赖项

**macOS：**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows：**

1. 安装[MSYS2](https://www.msys2.org/)以获得类 Unix 环境
2. WebView2 Runtime（Windows 11已预装；若使用 Windows 10，请[下载](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)）
3. 可选：安装[Git for Windows](https://git-scm.com/download/win)

**Linux（Debian/Ubuntu）：**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux（Fedora/RHEL）：**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux（Arch）：**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## 仓库设置

### 克隆并配置

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### 构建 Wails CLI

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### 添加到 PATH（可选）

**Linux/macOS：**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows：**

通过“系统属性”将 Wails 目录添加到 PATH 环境变量。

## IDE 设置

### VS Code（推荐）

1. **安装 VS Code：**[下载](https://code.visualstudio.com/)

2. **安装扩展：**
  - Go（由 Google Go 团队开发）
  - ESLint
  - Prettier
  - MDX（用于文档）


3. **配置工作区设置**（`.vscode/settings.json`）：
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **安装 GoLand：**[下载](https://www.jetbrains.com/go/)

2. **配置：**
  - 启用 Go Modules 支持
  - 为`goimports`设置文件监视器
  - 配置代码样式，使其符合项目约定


## 验证环境设置

运行以下命令，验证所有组件是否正常工作：

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

如果测试应用能够成功构建并运行，开发环境就已准备就绪！

## 运行测试

### 单元测试

```bash
cd v3
go test ./...
cd ..
```

### 特定软件包测试

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### 运行覆盖率测试

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 使用竞态检测器运行测试

```bash
cd v3
go test ./... -race
```

## 处理文档

Wails v3 文档使用 M-Press 编写。英文源文件位于  
`docs/mpress/content/`；译文位于各语言目录中，例如  
`fr/` 和 `id/`。

在仓库根目录中预览并验证文档更改：

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

生产站点是静态站点。在本地处理文档时，不需要 Node.js、翻译服务提供商或 Cloudflare  
凭据。

## 调试

### 调试 Go 代码

**VS Code：**

创建`.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**命令行：**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### 调试平台代码

调试平台特定代码需要使用相应平台的工具：

- <strong>macOS：</strong>Xcode Instruments
- <strong>Windows：</strong>Visual Studio Debugger
- <strong>Linux：</strong>GDB

## 常见问题

### “command not found: wails3”

将 Wails 目录添加到 PATH，或在项目根目录中使用`./wails3`。

### “未找到 webkitgtk-6.0”或“未找到 webkit2gtk”（Linux）

根据要构建所针对的技术栈，安装相应的开发包：

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### 构建因 Go 模块错误而失败

```bash
cd v3
go mod tidy
go mod download
```

### Windows 上的“CGO_ENABLED”错误

确保 PATH 中包含 C 编译器（通过 MSYS2 安装的 MinGW-w64）。

## 后续步骤

- 查阅[编码规范](/contributing/standards/)
- 浏览[技术文档](/contributing/)
- 查找可参与处理的问题：[适合初次贡献的问题](https://github.com/wailsapp/wails/labels/good%20first%20issue)
