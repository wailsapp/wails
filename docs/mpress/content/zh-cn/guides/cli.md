---
title: "CLI 标志与示例"
description: "查阅各个 Wails CLI 命令的标志和用法示例"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

如果你已经知道需要使用哪个命令，请查阅此参考文档。若要按用途浏览命令，请从[CLI 命令概览](/reference/cli/)开始。

Wails CLI 提供了一套全面的命令，可帮助你开发、构建和维护 Wails 应用程序。

## 核心命令

核心命令是用于创建、开发和构建项目的主要命令。

所有 CLI 命令均采用以下格式：`wails3 <command>`。

### `init`

初始化新的 Wails 项目。在初始化过程中，会运行`go mod tidy`命令以更新项目软件包。使用`init`命令时，可通过`-skipgomodtidy`标志跳过此步骤。

```bash
wails3 init [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-p` | Go 软件包名称 | `main` |
| `-t` | 模板名称或 URL | `vanilla` |
| `-n` | 项目名称 |  |
| `-d` | 项目目录 | `.` |
| `-q` | 不输出信息 | `false` |
| `-l` | 列出模板 | `false` |
| `-mod` | Go 模块路径（如果省略，则根据`-git`计算） |  |
| `-git` | Git 仓库 URL |  |
| `-s` | 使用远程模板时跳过警告 | `false` |
| `-productname` | 产品名称 | `My Product` |
| `-productdescription` | 产品说明 | `My Product Description` |
| `-productversion` | 产品版本 | `0.1.0` |
| `-productcompany` | 公司名称 | `My Company` |
| `-productcopyright` | 版权声明 | `© now, My Company` |
| `-productcomments` | 文件注释 | `This is a comment` |
| `-productidentifier` | 产品标识符 |  |
| `-skipgomodtidy` | 跳过 go mod tidy | `false` |

`-git`标志接受多种 Git URL 格式：

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project`或`ssh://git@github.com/username/project`
- Git 协议：`git://github.com/username/project`
- 文件系统：`file:///path/to/project.git`

提供此标志后，将执行以下操作：

1. 在项目目录中初始化 Git 仓库
2. 将指定 URL 设置为 origin 远程仓库
3. 更新`go.mod`中的模块名称，使其与仓库 URL 匹配
4. 添加所有文件

### `dev`

以开发模式运行应用程序。你可以实时查看前端代码，并在不重新构建整个应用程序的情况下进行更改并立即在运行中的应用程序里看到效果。系统也会检测 Go 代码的更改，并自动重新构建和启动应用程序。

```bash
wails3 dev [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-config` | 配置文件路径 | `./build/config.yml` |
| `-port` | Vite 开发服务器端口 | `9245` |
| `-s` | 启用 HTTPS | `false` |

@note{type="info"}
这等同于运行`wails3 task dev`，并会执行项目主 Taskfile 中的`dev`任务。你可以编辑`Taskfile.yml`文件来自定义此行为。

@end

### `build`

使用生成的 Taskfile 构建生产版本。使用`wails3 dev`可启用热重载，传入`DEV=true`则可进行开发构建。默认针对当前平台和架构进行构建。

```bash
wails3 build [flags] [CLI variables...]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-tags` | 额外的 Go 构建标签（以逗号分隔） |  |

你可以传入 CLI 变量来自定义构建：

```bash
wails3 build GOOS=linux
```

使用`-tags`标志传入自定义 Go 构建标签：

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

标签会以`EXTRA_TAGS`的形式转发给底层 Taskfile。

@note{type="info"}
这等同于运行`wails3 task build`，后者会执行项目主 Taskfile 中的`build`任务。传给`build`的所有 CLI 变量都会转发给底层任务。你可以编辑`Taskfile.yml`文件来自定义构建流程。

@end

### `package`

创建用于分发的特定平台软件包。

```bash
wails3 package [CLI variables...]
```

标准打包任务接受`ARCH`，用于选择目标架构：

```bash
wails3 package ARCH=amd64
```

`VERSION`和`OUTPUT`仅对显式读取这些变量的自定义 Taskfile 生效；它们不会覆盖标准软件包的版本或文件名。

#### 软件包类型

各平台支持以下软件包类型：

| 平台 | 软件包类型 |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
这等同于`wails3 task package`，后者会运行项目主 Taskfile 中的`package`任务。传递给`package`的所有 CLI 变量都会转发给底层任务。你可以通过编辑`Taskfile.yml`文件来自定义打包过程。

@end

标准构建任务会读取`GOOS`、`ARCH`和`DEV`；下方自定义任务中的`PLATFORM`和`CONFIG`是用户定义的示例。

### `task`

运行项目 Taskfile.yml 中定义的任务。这是[Taskfile](https://taskfile.dev)的嵌入式版本，可用于定义和运行自定义构建、测试及部署任务。

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI 变量

你可以按`KEY=VALUE`格式向任务传递变量：

```bash
wails3 task build GOOS=linux
wails3 task deploy ENV=staging VERSION=1.2.3
```

可以在 Taskfile.yml 中使用 Go 模板语法访问这些变量：

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-h` | 显示 Task 用法 | `false` |
| `-i` | 创建新的 Taskfile.yml | `false` |
| `-list` | 列出带说明的任务 | `false` |
| `-list-all` | 列出所有任务（无论是否带说明） | `false` |
| `-json` | 以 JSON 格式输出任务列表 | `false` |
| `-status` | 如果任务不是最新状态，则以非零状态码退出 | `false` |
| `-f` | 即使任务处于最新状态也强制执行 | `false` |
| `-w` | 为指定任务启用监视模式 | `false` |
| `-v` | 启用详细模式 | `false` |
| `-version` | 输出 Task 版本 | `false` |
| `-s` | 禁用命令回显 | `false` |
| `-p` | 并行执行任务 | `false` |
| `-dry` | 编译并输出任务，但不执行 | `false` |
| `-summary` | 显示任务摘要 | `false` |
| `-x` | 透传任务的退出状态码 | `false` |
| `-dir` | 设置执行目录 |  |
| `-taskfile` | 选择要运行的 Taskfile |  |
| `-output` | 设置输出样式：[interleaved|group|prefixed] |  |
| `-c` | 彩色输出（默认启用） | `true` |
| `-C` | 限制并发运行的任务数量 |  |
| `-interval` | 监视变更的间隔（秒） |  |

#### 示例

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build GOOS=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

启动 Wails 项目的 MCP 服务器，以便通过代理辅助管理项目。它不同于编译到运行中应用程序内的 MCP 服务器：`wails3 mcp`管理项目文件和生命周期命令，而应用程序 MCP 服务器控制运行中的 WebView。

```bash
wails3 mcp [flags]
```

传输方式会自动选择：

- 当 MCP 主机使用管道连接的 stdin/stdout 启动 Wails 时，服务器使用**stdio**。
- 在终端中以交互方式运行时，服务器在`127.0.0.1`上使用**Streamable HTTP**，并请求操作系统分配一个空闲端口。

使用`--stdio`或`--http`可明确选择传输方式。在 HTTP 模式下，使用`--port 0`可选择一个空闲的回环端口。

#### MCP 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `--root` | 允许的项目根目录。位于其外部的路径和符号链接会被拒绝。 | 当前目录 |
| `--token` | 用于修改操作和进程控制工具的会话/持有者令牌。未设置时回退到`WAILS_MCP_TOKEN`。 | 安全生成 |
| `--stdio` | 强制使用 stdio 传输。 | 自动 |
| `--http` | 强制使用 Streamable HTTP 传输。 | 自动 |
| `--port` | HTTP 端口；`0` 表示选择一个空闲的环回端口。 | `0` |

在 HTTP 模式下，Wails 会将端点和持有者令牌输出到 stderr。在 stdio 模式下，令牌包含在 MCP 初始化说明中。服务器不提供任意 shell 命令执行功能。远程模板和 Git 远程仓库需要通过工具的 `allowExternal` 输入显式批准。

### `doctor`

执行系统检查并显示状态报告。

```bash
wails3 doctor
```

## 生成命令

生成命令可帮助创建绑定、图标和构建文件等各种项目资源。所有生成命令都使用基础命令：`wails3 generate <command>`。

### `generate bindings`

为 Go 代码生成绑定和模型。

```bash
wails3 generate bindings [flags] [patterns...]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-f` | 其他 Go 构建标志 |  |
| `-d` | 输出目录 | `frontend/bindings` |
| `-models` | 模型文件名 | `models` |
| `-index` | 索引文件名 | `index` |
| `-ts` | 生成 TypeScript | `false` |
| `-i` | 使用 TS 接口 | `false` |
| `-b` | 使用捆绑的运行时 | `false` |
| `-names` | 使用名称而非 ID | `false` |
| `-noindex` | 跳过索引文件 | `false` |
| `-noevents` | 跳过生成事件相关的绑定 | `false` |
| `-dry` | 试运行 | `false` |
| `-silent` | 静默模式 | `false` |
| `-v` | 调试输出 | `false` |
| `-clean` | 生成前清理输出目录 | `true` |

### `generate build-assets`

为应用程序生成构建资源。

```bash
wails3 generate build-assets [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-name` | 项目名称 |  |
| `-dir` | 输出目录 | `build` |
| `-silent` | 禁止输出 | `false` |
| `-company` | 公司名称 |  |
| `-productname` | 产品名称 |  |
| `-description` | 产品描述 |  |
| `-version` | 产品版本 |  |
| `-identifier` | 产品标识符 | `com.wails.[name]` |
| `-copyright` | 版权声明 |  |
| `-comments` | 文件注释 |  |

### `generate icons`

生成应用程序图标。

```bash
wails3 generate icons [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-input` | 输入 PNG 文件 | 必填 |
| `-windowsfilename` | Windows 输出文件名 |  |
| `-macfilename` | macOS 输出文件名 |  |
| `-sizes` | 图标尺寸（以逗号分隔） | `256,128,64,48,32,16` |
| `-example` | 生成示例图标 | `false` |
| `-iconcomposerinput` | 输入 Icon Composer 文件（`.icon`） |  |
| `-macassetdir` | Mac 资源的输出目录（Assets.car + icns） |  |

#### Icon Composer（macOS）

在 macOS 26 或更高版本中，可以使用 Icon Composer 的 `.icon` 文件生成 `Assets.car` 和 `icons.icns`：

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

此操作使用 Apple 的 `actool` 命令编译 `.icon` 文件。需要安装 Xcode，且 `actool` 版本须为 26 或更高版本。

使用 Icon Composer 时，请在 `build/config.yml` 中设置 `cfBundleIconName`，使其与 `.icon` 文件名（不含扩展名）一致：

```yaml
info:
  cfBundleIconName: "appicon"
```

如果未设置且 `Assets.car` 存在，则默认使用 `"appicon"`。

### `generate syso`

生成 Windows .syso 文件。

```bash
wails3 generate syso [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-manifest` | 清单文件的路径 | 必填 |
| `-icon` | 图标文件的路径 | 必填 |
| `-info` | 版本信息文件的路径 |  |
| `-arch` | 目标架构 | 当前 GOARCH |
| `-out` | 输出文件名 | `rsrc_windows_[arch].syso` |

### `generate .desktop`

生成 Linux .desktop 文件。

```bash
wails3 generate .desktop [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-name` | 应用程序名称 | 必填 |
| `-exec` | 可执行文件路径 | 必填 |
| `-icon` | 图标路径 |  |
| `-categories` | 应用程序类别 | `Utility` |
| `-comment` | 应用程序注释 |  |
| `-terminal` | 在终端中运行 | `false` |
| `-keywords` | 搜索关键字 |  |
| `-version` | 应用程序版本 |  |
| `-genericname` | 通用名称 |  |
| `-startupnotify` | 显示启动通知 | `false` |
| `-mimetype` | 支持的 MIME 类型 |  |
| `-output` | 输出文件名 | `[name].desktop` |

### `generate runtime`

生成预构建版本的运行时。

```bash
wails3 generate runtime
```

### `generate constants`

从 Go 代码生成 JavaScript 常量。

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

生成用于分发的 Windows WebView2 引导安装程序。

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

搭建新的项目模板目录。

```bash
wails3 generate template [flags]
```

### `generate appimage`

生成 Linux AppImage。

```bash
wails3 generate appimage [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-binary` | 二进制文件的路径 | 必填 |
| `-icon` | 图标文件的路径 | 必填 |
| `-desktop` | .desktop 文件的路径 | 必填 |
| `-builddir` | 构建目录 | 临时目录 |
| `-output` | 输出目录 | `.` |

## 服务命令

服务命令用于管理 Wails 服务。所有服务命令均使用基础命令：`wails3 service <command>`。

### `service init`

初始化新服务。

```bash
wails3 service init [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-n` | 服务名称 | `example_service` |
| `-d` | 服务说明 | `Example service` |
| `-p` | 包名称 |  |
| `-o` | 输出目录 | `.` |
| `-q` | 不显示输出 | `false` |
| `-a` | 作者姓名 |  |
| `-v` | 版本 |  |
| `-w` | 网站 URL |  |
| `-r` | 仓库 URL |  |
| `-l` | 许可证 |  |

## 工具命令

工具命令提供用于开发和调试的实用工具。所有工具命令均使用基础命令：`wails3 tool <command>`。

### `tool checkport`

检查端口是否开放。可用于测试 vite 是否正在运行。

```bash
wails3 tool checkport [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-port` | 要检查的端口 | `9245` |
| `-host` | 要检查的主机 | `localhost` |

### `tool watcher`

监视文件，并在文件发生更改时运行命令。

```bash
wails3 tool watcher [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-config` | 配置文件路径 | `./build/config.yml` |
| `-ignore` | 要忽略的模式 |  |
| `-include` | 要包含的模式 |  |

### `tool cp`

复制文件。

```bash
wails3 tool cp
```

### `tool buildinfo`

显示应用程序的构建信息。

```bash
wails3 tool buildinfo
```

### `tool version`

根据提供的标志递增语义化版本号。

```bash
wails3 tool version [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-v` | 要递增的当前版本 |  |
| `-major` | 递增主版本号 | `false` |
| `-minor` | 递增次版本号 | `false` |
| `-patch` | 递增修订版本号 | `false` |
| `-prerelease` | 递增预发布版本号（例如从 alpha.5 递增到 alpha.6） | `false` |

该命令遵循以下优先顺序：主版本号 > 次版本号 > 修订版本号 > 预发布版本号。如果输入版本中存在“v”前缀，则会保留该前缀，同时还会保留所有预发布和元数据组成部分。

用法示例：

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.1.0-alpha.5 -prerelease # Output: v3.1.0-alpha.6
```

### `tool package`

生成 Linux 软件包（deb、rpm、archlinux）。

```bash
wails3 tool package [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-format` | 软件包格式（deb、rpm、archlinux） | `deb` |
| `-name` | 可执行文件名称 | `myapp` |
| `-config` | 配置文件路径 |  |
| `-out` | 输出目录 | `.` |

### `tool lipo`

通过合并特定于各架构的二进制文件，创建 macOS 通用二进制文件。

```bash
wails3 tool lipo [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-output` | 输出二进制文件路径 |  |

### `tool capabilities`

检查系统的构建能力（Linux 上 GTK4/GTK3 的可用性）。

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

生成用于交叉编译的 Docker 卷挂载标志。为 Go 模块缓存以及 `go.mod` 中所有本地 `replace` 指令输出 `-v` 标志，供 Taskfile 的 `docker run` 命令使用。

```bash
wails3 tool docker-mounts
```

### `tool has`

检查工具或功能是否可用，并将 `true` 或 `false` 输出到标准输出。此命令专为在 Taskfile 的 `sh:` 变量中使用而设计，可作为 `command -v` 的跨平台替代方案。

使用 `|` 可检查多个备选项中是否有任意一个可用。

```bash
wails3 tool has <tool>
```

#### 示例

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### 在 Taskfile 中使用

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="已弃用"}
`wails3 tool has-cc` 已弃用。请更新 Taskfile，改用 `wails3 tool has gcc|clang`。

@end

`wails3 tool has gcc|clang` 的向后兼容别名。检查 PATH 中是否存在 `gcc` 或 `clang`，并输出 `true` 或 `false`。

```bash
wails3 tool has-cc
```

## 更新命令

更新命令用于管理和更新项目资源。所有更新命令都使用基础命令：`wails3 update <command>`。

### `update cli`

将 Wails CLI 更新到新版本。

```bash
wails3 update cli [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-pre` | 更新到最新预发布版本 | `false` |
| `-version` | 更新到指定版本 |  |
| `-nocolour` | 禁用彩色输出 | `false` |

update cli 命令可用于更新 Wails CLI 安装。默认情况下，它会更新到最新稳定版本。可以使用 `-pre` 标志更新到最新预发布版本，也可以使用 `-version` 标志指定特定版本。

更新后，请记得更新项目的 go.mod 文件以使用相同版本：

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

使用给定的配置文件更新构建资源。

```bash
wails3 update build-assets [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-config` | 配置文件路径 |  |
| `-dir` | 输出目录 | `build` |
| `-silent` | 禁止输出 | `false` |
| `-company` | 公司名称 |  |
| `-productname` | 产品名称 |  |
| `-description` | 产品说明 |  |
| `-version` | 产品版本 |  |
| `-identifier` | 产品标识符 |  |
| `-copyright` | 版权声明 |  |
| `-comments` | 文件注释 |  |

## 实用命令

实用命令为常见任务提供便捷的快捷方式。请直接与基础命令一起使用这些命令：`wails3 <command>`。

### `docs`

在默认浏览器中打开 Wails 文档。

```bash
wails3 docs
```

### `releasenotes`

显示当前版本或指定版本的发行说明。

```bash
wails3 releasenotes [flags]
```

#### 标志

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-v` | 要显示发行说明的版本 |  |
| `-n` | 禁用彩色输出 | `false` |

### `version`

输出当前的 Wails 版本。

```bash
wails3 version
```

### `sponsor`

在默认浏览器中打开 Wails 赞助页面。

```bash
wails3 sponsor

```
