---
title: "使用 M-Press 编写文档"
description: "下载并安装 M-Press，预览 Wails 文档，然后提交文档贡献。"
sourcePath: "contributing/documentation.md"
---

Wails v3使用**M-Press**将 Markdown 转换为此文档站点。你可以修正拼写错误、改进说明、添加示例或指南，也可以修正翻译。仅涉及文档的修正无需事先创建议题，也无需先出现失败的代码测试。

你需要使用 Git 准备贡献，并需要 GitHub 账户来提交拉取请求。可下载的 M-Press 可执行文件包含文档生成器和图表渲染器；预览文档不需要 Go、Node.js、Wails CLI 或桌面应用构建依赖项。仓库中的附加验证脚本需要 Python 3。

## 从文档页面开始

在侧边栏中打开**贡献 → 文档**以进入本指南。要处理特定页面，请选择该页面标题栏中的**编辑图标**，然后选择**在本地编辑此页面**。复制适用于你的操作系统的命令，并在终端中运行。

如果已安装 M-Press，该命令会直接使用它；否则，安装程序会下载官方版本并验证校验和。它会准备一个用于贡献的本地工作副本，并在 M-Press 中打开所选页面。在运行仓库检查之前，请对照下方固定的版本检查`mpress version`：自动安装程序可能会下载更新的版本。

本地预览重新构建时，你可以使用自己的编辑器编辑源文件。向导还将**审阅更改**、**提交更改**、**推送贡献分支**和**创建草稿拉取请求**作为独立操作提供。提交前，请遵循下文的检查和拉取请求指南。生成翻译是可选操作，并且可能需要配置提供商；编辑和预览 Markdown 则不需要。

如需可复现的手动设置，请按照以下步骤操作。

## 下载并安装 M-Press

仓库构建脚本会下载**最新的 M-Press 稳定版**。请在本地使用同一版本：选择下面的一种安装方式，运行`mpress version`，然后将其与发布页面上显示的版本进行比较。

### 下载可执行文件

打开[最新的 M-Press 稳定版](https://github.com/leaanthony/mpress/releases/latest)发布页面，或在下面下载对应的压缩包。ARM64 包括 Apple Silicon；AMD64 是 Intel/AMD PC 常用的架构。

| 计算机 | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

从同一个发布版本下载[checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt)。在下载目录中，计算压缩包的 SHA-256 哈希值，并将其与`checksums.txt`中该文件名所在行的值进行比较：

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

如果选择了其他架构，请使用实际下载的文件名。确认哈希值匹配后，解压并安装可执行文件：

@tabs
[macOS]
在下载目录中，使用 Apple Silicon 归档文件：

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

对于 Intel Mac，请替换为`mpress-darwin-amd64.tar.gz`。将`export PATH`这一行添加到`~/.zshrc`，以便在新终端中继续使用该命令。

[Linux]
在下载目录中：

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

对于 ARM64，请替换为`mpress-linux-arm64.tar.gz`。将`export PATH`这一行添加到 shell 启动文件（例如`~/.bashrc`）中，以便在新终端中继续使用该命令。

[Windows PowerShell]
在下载目录中：

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

对于 ARM64，请替换为`mpress-windows-arm64.zip`。在 Windows 的**环境变量**中，将你的用户`bin`目录添加到用户**Path**，然后重新打开终端，以便在后续会话中继续使用该命令。
@end

### 使用 Go 安装

如果你已在使用 Go，也可以改为自行构建最新的稳定版。请查看其`go.mod`以了解所需的工具链；Go 的自动工具链选择功能可能会下载该工具链。本轮文档更新所使用的版本要求 Go 1.26.10或更高版本。

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go 会将该命令安装到`GOBIN`；如果未设置`GOBIN`，则安装到`go env GOPATH`下的`bin`目录。如果找不到`mpress`，请将该目录添加到`PATH`，然后重新打开终端。

## 获取文档源文件

[复刻 wailsapp/wails](https://github.com/wailsapp/wails/fork)，然后运行以下命令，并将`YOUR_USERNAME`替换为你的 GitHub 用户名：

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

这会基于上游的**master**分支创建你的文档分支。请用你自己的描述性分支名称替换`docs/my-change`。

## 运行实时预览

在包含`mpress.yaml`的 Wails 仓库根目录中：

```sh
mpress dev
```

在浏览器中打开[http://127.0.0.1:3000](http://127.0.0.1:3000)。要查看本指南，请访问[http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/)。编辑期间请让该命令保持运行；保存源文件更改后会重新构建预览。在运行下方的验证命令之前，按**Ctrl+C**停止该命令。

仓库已包含 M-Press 配置，无需初始化新站点。要使用其他端口，请运行`mpress dev --port 3001`并打开对应的 localhost 地址。

## 编辑或添加页面

英文源页面位于`docs/mpress/content/`中。翻译位于语言文件夹下，例如`fr/`和`id/`。编辑图标中的贡献说明会指出当前页面的源文件。

对于现有页面，请编辑其 Markdown，并保留顶部的 YAML 元数据。使用标题、段落、列表、链接和围栏代码块。请保留成对的 M-Press 组件，例如`@note` / `@end`和`@tabs` / `@end`。[M-Press 编写文档](https://m-press.me/authoring/)介绍了支持的语法。

要创建新指南，请新建一个文件，例如`docs/mpress/content/guides/my-guide.md`：

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

将页面添加到`docs/mpress/content/_nav.yaml`中的适当部分，除非该部分已使用`autogenerate`。检查读者能否通过导航以及相关页面中的链接访问该页面。新的英文页面还需要提供完整的各语言副本，覆盖率检查才能通过；请与维护者或翻译贡献者协调完成。

生成的输出位于`docs/mpress/site/`中。请提交源页面以及有意进行的导航或静态资源更改；生成的站点文件和本地缓存不应包含在拉取请求中。

## 检查你的贡献

在仓库根目录中：

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

在 Windows 上，如果 Python 3是以这种方式安装的，请使用`python`。如果已安装[Task](https://taskfile.dev/)，则`task docs:build`会使用 PATH 中的`mpress`运行相同的检查。

翻译检查会在翻译文件缺失、为空或成为孤立文件时失败。默认情况下，过时或损坏的译文段落会以警告形式报告；`--strict`还会使这些问题导致检查失败。修改现有英文正文时，即使翻译同步尚未跟上，也可以提交。修正某一种语言的翻译时，可以只更改该语言。请保持命令、API 名称、链接、代码示例和图表连接不变；不要将英文占位内容用作译文。

请在浏览器中阅读已更改的页面，并运行你修改过的所有代码示例。如果检查因未改动的页面而失败，请在 PR 中附上输出和相关上下文，而不要降低检查要求。

## 提交拉取请求

检查差异，提交源文件更改，然后将贡献分支推送到你的 fork。向`wailsapp/wails`中的**master**分支发起 PR。说明所做的更改，链接受影响的页面，并列出你运行过的检查。对于布局更改，请附上更改前后的屏幕截图；对于有改动的代码示例，请提供平台和版本详情。

生产站点是静态站点；本地预览和验证不需要部署凭据或翻译服务提供商。维护者会在审核后负责发布。

有关代码贡献和功能提案，请参阅[为 Wails 做贡献](/contributing/)。有关框架内部机制，请参阅[技术概览](/contributing/overview/)。
