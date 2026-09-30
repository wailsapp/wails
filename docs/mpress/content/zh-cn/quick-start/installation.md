---
title: "安装"
description: "安装 Wails 并做好构建应用程序的准备"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## 快速安装（5 分钟）

@note{type="tip" title="简要说明——面向有经验的开发者"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

也可以使用`wails3 doctor`手动验证。[跳转到第一个应用 →](/quick-start/first-app/)

@end

## 分步安装

@steps
### 安装 Go（必需）
Wails 需要 Go 1.25或更高版本。

@tabs{sync-key="os"}
[Windows]
从<strong>[go.dev/dl](https://go.dev/dl/)</strong>下载 Windows 安装程序并运行。

**验证安装：**

```powershell
go version  # Should show 1.25 or later
```

**检查 PATH：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

如果输出为空，请将`C:\Users\YourName\go\bin`添加到 PATH。

[macOS]
**选项1：官方安装程序**

从<strong>[go.dev/dl](https://go.dev/dl/)</strong>下载 macOS 安装程序（.pkg 文件）并运行。

**选项2：Homebrew**

```bash
brew install go
```

**验证安装：**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

如果 PATH 中没有`~/go/bin`，请将其添加到`~/.zshrc`或`~/.bash_profile`：

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**选项1：官方 Tar 包**

从<strong>[go.dev/dl](https://go.dev/dl/)</strong>下载 Linux tar 包，然后执行：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**选项2：软件包管理器**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**添加到 PATH**（添加到`~/.bashrc`或`~/.zshrc`）：

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**验证：**

```bash
go version
echo $PATH | grep go/bin
```

@end

### 安装平台依赖项
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime**（通常已预装）

Windows 10/11默认包含 WebView2。如果缺失：

- 从[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)下载
- 也可以稍后运行`wails3 doctor`，它会提供引导

<strong>就这么简单！</strong>无需其他依赖项。

@note{type="tip" title="Windows 11性能提示"}
考虑使用[Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)存储项目。Dev Drive 针对开发工作负载进行了优化，可显著缩短构建时间并提高磁盘访问速度，两者的改善幅度均最高可达30%。

@end

[macOS]
**Xcode Command Line Tools**（必需）

```bash
xcode-select --install
```

在出现的对话框中单击“安装”。

**验证：**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

<strong>就这么简单！</strong>macOS 默认包含 WebKit。

[Linux]
**构建工具和 WebKit**

@note{type="caution" title="最低发行版版本"}
Wails v3 默认需要<strong>WebKitGTK 6.0</strong>。仅提供 WebKit2GTK 4.1的发行版——Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x——必须选择启用旧版`-tags gtk3`进行构建。仅提供 WebKit2GTK 4.0的更旧版本（Ubuntu 20.04、Debian 11、RHEL 8）不受支持。

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
默认 GTK4 技术栈要求 Ubuntu 24.04或更高版本，或者 Debian 13或更高版本。

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
添加到`shell.nix`或`devShell`：

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[其他]
安装 Wails 后运行`wails3 doctor`，它会显示你的发行版所需的确切软件包。

@end

@note{type="info" title="旧版 GTK3 技术栈"}
如果目标发行版尚未提供 WebKitGTK 6.0（例如 Ubuntu 22.04 LTS、Debian 12），请改为安装 GTK3 + WebKit2GTK 4.1开发库（在 Debian/Ubuntu 上为`libgtk-3-dev libwebkit2gtk-4.1-dev`；其他发行版请安装对应软件包），并使用`wails3 build -tags gtk3`进行构建。旧版路径在 v3.0.x 系列中受支持，并将在 v3.1中移除。有关详细信息，请参阅[Linux 打包——旧版 GTK3 支持](/guides/build/linux/#legacy-gtk3-support)。

@end

@end

### 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

这会将`wails3`命令安装到`~/go/bin`（Windows 上则安装到`%USERPROFILE%\go\bin`）。

### 运行设置向导（推荐）
```bash
wails3 setup
```

设置向导会检查依赖项，帮助安装缺失的依赖项，并配置项目默认设置。

@note{type="caution" title="实验性功能"}
设置向导是一项新功能，目前主要在 Linux 上经过测试。如果遇到问题，请[报告问题](https://github.com/wailsapp/wails/issues/4904)，并改用`wails3 doctor`。

@end

### 验证安装
```bash
wails3 doctor
```

**预期输出（或类似内容）：**

```
Wails (v3.0.0-dev)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | Your installed version |
| Go Version  | go1.25.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="如果找不到 `wails3` 命令"}
PATH 中没有`~/go/bin`。请参照上面的步骤1修复此问题，然后重新启动终端。

@end

### 安装 npm（可选但推荐）
大多数 Wails 模板使用 npm 作为前端工具。

@tabs{sync-key="os"}
[Windows]
从[nodejs.org](https://nodejs.org/)下载安装程序并运行。

**验证：**

```powershell
npm --version
```

[macOS]
**选项 1：官方安装程序** 从[nodejs.org](https://nodejs.org/)下载

**选项 2：Homebrew**

```bash
brew install node
```

**验证：**

```bash
npm --version
```

[Linux]
**选项 1：NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**选项 2：软件包管理器**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**验证：**

```bash
npm --version
```

@end

@note{type="tip" title="其他软件包管理器"}
更喜欢`pnpm`、`yarn`或`bun`？没问题！只需更新项目中的`Taskfile.yml`，改用你偏好的工具。

@end

@end

## 故障排除

### 找不到`wails3`命令

**原因：**`~/go/bin`（或`%USERPROFILE%\go\bin`）不在 PATH 中。

**解决方法：**

@tabs{sync-key="os"}
[Windows]
1. 打开“环境变量”（在“开始”菜单中搜索）
2. 在“用户变量”下找到`Path`
3. 单击“编辑”→“新建”
4. 添加：`C:\Users\YourName\go\bin`（替换`YourName`）
5. 在所有对话框中单击“确定”
6. **重启终端**

**验证：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
添加到`~/.zshrc`（macOS）或`~/.bashrc`（Linux）：

```bash
export PATH=$PATH:~/go/bin
```

重新加载：

```bash
source ~/.zshrc  # or ~/.bashrc
```

**验证：**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor`报告缺少依赖项

<strong>Linux：</strong>输出会明确告知你需要安装哪些软件包。例如：

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

<strong>Windows：</strong>如果缺少 WebView2：

- 从[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)下载
- 或者，在你首次运行应用时，它会自动安装

<strong>macOS：</strong>如果缺少 Xcode 工具：

```bash
xcode-select --install
```

---

#### Go 版本过旧

Wails v3 要求 Go 1.25 或更高版本。如果你的版本较旧：

@tabs{sync-key="os"}
[Windows/macOS]
从[go.dev/dl](https://go.dev/dl/)下载最新版本并重新安装。

[Linux]
从[go.dev/dl](https://go.dev/dl/)下载最新的 tarball，然后：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## 开发版本（最前沿版本）

想使用主开发分支中的绝对最新代码？这样可以在新功能和修复正式发布前使用它们，但也有遇到错误和破坏性变更的风险。仅建议贡献者或需要测试即将推出功能的用户使用。

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="开发版本"}
- 可能存在错误或破坏性变更
- 创建的项目将使用`replace`指令指向本地 Wails
- 仅建议贡献者或测试新功能时使用

@end

## 后续步骤

<strong>安装完成！</strong>你的系统已准备好进行 Wails 开发。

@cards{cols="1"}
🚀 构建你的第一个应用
在10分钟内创建一个可运行的应用。

[第一个应用教程 →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 探索模板
查看开箱即用的内容。

```bash
wails3 init -l  # List templates
```

@end

---

<strong>遇到问题？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或[提交 issue](https://github.com/wailsapp/wails/issues)。
