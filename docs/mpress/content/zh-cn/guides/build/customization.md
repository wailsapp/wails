---
title: "构建自定义"
description: "使用 Task 和 Taskfile.yml 自定义构建流程"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 概述

Wails 构建系统是一款灵活而强大的工具，旨在简化 Wails 应用程序的构建流程。它利用任务运行器[Task](https://taskfile.dev)，让你可以轻松定义和运行任务。虽然 v3 构建系统是默认选项，但 Wails 鼓励采用“自带工具”的方式，允许开发者根据需要自定义构建流程。

有关如何使用 Task 的更多信息，请参阅[官方文档](https://taskfile.dev/usage/)。

## Task：构建系统的核心

[Task](https://taskfile.dev) 是使用 Go 编写的现代 Make 替代方案。它使用 YAML 文件定义任务及其依赖项。在 Wails 构建系统中，[Task](https://taskfile.dev) 在编排构建流程方面发挥着核心作用。

主`Taskfile.yml`位于项目根目录中，而平台特定任务则在`build/<platform>/Taskfile.yml`文件中定义。`build`目录中的通用`Taskfile.yml`文件包含各平台共享的通用任务。

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

项目根目录中的`Taskfile.yml`文件是构建系统的主要入口点。它定义任务及其依赖项。以下是默认的`Taskfile.yml`文件：

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## 平台特定的 Taskfile

每个平台都有自己的 Taskfile，位于`build`目录下对应的平台目录中。这些文件定义该平台的核心任务。每个 Taskfile 都会包含`build/Taskfile.yml`文件中的通用任务。

### Windows

位置：`build/windows/Taskfile.yml`

Windows 特定的 Taskfile 包含在 Windows 上构建、打包和运行应用程序的任务。主要功能包括：

- 使用可选的生产环境标志进行构建
- 生成`.ico`图标文件
- 生成 Windows `.syso`文件
- 创建用于打包的 NSIS 安装程序

### Linux

位置：`build/linux/Taskfile.yml`

Linux 特定的 Taskfile 包含在 Linux 上构建、打包和运行应用程序的任务。主要功能包括：

- 使用可选的生产环境标志进行构建
- 创建 AppImage、deb、rpm 和 Arch Linux 软件包
- 为 Linux 应用程序生成`.desktop`文件

### macOS

位置：`build/darwin/Taskfile.yml`

macOS 特定的 Taskfile 包含在 macOS 上构建、打包和运行应用程序的任务。主要功能包括：

- 为 amd64、arm64 和 universal（两者兼容）架构构建二进制文件
- 生成`.icns`图标文件
- 创建用于分发的`.app`捆绑包
- 对`.app`捆绑包进行临时签名
- 设置 macOS 特定的构建标志和环境变量

## 任务执行和命令别名

`wails3 task`命令是[Taskfile](https://taskfile.dev)的内嵌版本，用于执行`Taskfile.yml`中定义的任务。

`wails3 build`和`wails3 package`命令分别是`wails3 task build`和`wails3 task package`的别名。运行这些命令时，Wails 会在内部将其转换为相应的任务执行命令：

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### 向任务传递参数

你可以使用`KEY=VALUE`格式向任务传递 CLI 变量。这些变量会通过别名命令转发：

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

在`Taskfile.yml`中，你可以使用 Go 模板语法访问这些变量：

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

### Go 构建自定义变量

生成的项目提供用于追加自定义 Go 构建标签、链接器参数以及配置 CGO 的变量。标签值必须是以逗号分隔的名称，不包含 `-tags` 选项；Wails 会将它们与所选构建模式和平台需要的标签组合起来。

| 变量 | 适用范围 |
| --- | --- |
| `APP_TAGS` | 所有通过 Taskfile 执行的构建 |
| `APP_TAGS_LINUX` | Linux 构建 |
| `APP_TAGS_DARWIN` | macOS 构建 |
| `APP_TAGS_WINDOWS` | Windows 构建 |
| `APP_TAGS_ANDROID` | Android 构建 |
| `APP_TAGS_IOS` | iOS 构建 |
| `APP_TAGS_SERVER` | 服务器模式构建 |
| `APP_LDFLAGS` | 为所有通过 Taskfile 执行的构建追加链接器参数 |
| `APP_CGO_ENABLED` | 以 `0` 或 `1` 覆盖桌面和服务器构建的 CGO 设置 |
| `EXTRA_TAGS` | 为单次调用追加标签 |

例如：

```bash
wails3 build APP_TAGS=sqlite_fts5,netgo APP_TAGS_LINUX=myapp_linux
wails3 build EXTRA_TAGS=diagnostics
APP_LDFLAGS='-X example.com/myapp/internal/version.Value=1.2.3' wails3 build
```

如果要设置持久的项目默认值，请在根目录的 `Taskfile.yml` 中设置 `APP_*` 值。随命令传入的 `KEY=value` 优先级最高，会在本次调用中替换该变量。根 Taskfile 中的字面值优先于同名的进程环境变量。如果保留下方所示的自引用默认值表达式，则在未提供命令行值时使用进程环境。若只想为一次构建追加标签而不替换持久的 `APP_TAGS` 值，请使用 `EXTRA_TAGS`。

当 `APP_CGO_ENABLED` 为空时，Linux 和 macOS 默认为 `1`，Windows 默认为 `0`，原生服务器构建保留 Go 在宿主机上的默认值，服务器 Docker 构建默认为 `0`。Android 和 iOS 构建始终使用 CGO，并保持 `CGO_ENABLED=1`；`APP_CGO_ENABLED` 不会覆盖这些移动工具链的设置。基于 Docker 的桌面交叉编译和服务器构建会接收与对应原生 Taskfile 构建相同的适用 `APP_*` 值。

@note{type="info" title="现有项目"}
根目录的 `Taskfile.yml` 由项目维护，因此 `wails3 update build-assets` 不会覆盖它。在引入这些变量之前创建的项目，需要手动将以下条目添加到根 Taskfile 的 `vars` 块中：

```yaml
vars:
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
```

将自引用默认值表达式替换为字面值，即可设置持久的项目默认值。

@end

这些变量由通过 Taskfile 执行的构建使用。生成的 Xcode 项目在构建阶段直接调用 Go，不经过这套 Taskfile 自定义流程，因此从 Xcode 启动的构建目前不会使用这些变量。

## 通用构建流程

在所有平台上，构建流程通常包括以下步骤：

1. 整理 Go 模块
2. 构建前端
3. 生成图标
4. 使用平台特定标志编译 Go 代码
5. 打包应用程序（特定于平台）

## 自定义构建流程

虽然 v3 构建系统提供了可靠的默认配置，但你可以轻松对其进行自定义，以满足项目需求。通过修改`Taskfile.yml`和平台特定的 Taskfile，你可以：

- 添加新任务
- 修改现有任务
- 更改任务执行顺序
- 与其他工具和脚本集成

借助这种灵活性，你可以根据具体需求调整构建流程，同时仍能受益于 Wails 构建系统提供的结构。

@note{type="tip" title="学习 Taskfile"}
强烈建议阅读[Taskfile](https://taskfile.dev)文档，以了解如何有效使用 Taskfile。运行`wails3 task --version`可以查看 Wails CLI 内嵌的 Taskfile 版本。

@end

## 开发模式

Wails 构建系统包含强大的开发模式，可通过实时重新加载和热模块替换提升开发者体验。使用`wails3 dev`命令可启用此模式。

### 工作原理

运行`wails3 dev`时，会执行以下流程：

1. 该命令会检查可用端口；如果未指定端口，则默认使用9245。
2. 它会为前端开发服务器（Vite）设置环境变量。
3. 它会使用[refresh](https://github.com/atterpac/refresh)库启动文件监视器。

[refresh](https://github.com/atterpac/refresh)库负责监视文件变更并触发重新构建。它使用`./build/config.yml`文件中`dev_mode`键下定义的配置。可以将其配置为忽略某些目录和文件、确定要监视的文件，以及检测到变更时要执行的操作。默认配置通常已经很好用，但你可以根据需要进行自定义。

### 配置

下面是其结构示例：

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

此配置文件可用于：

- 设置文件监视的根路径
- 配置日志级别
- 设置文件变更事件的防抖时间
- 忽略特定目录、文件或文件扩展名
- 定义文件发生变更时要执行的命令

### 自定义开发模式

你可以修改`config.yml`文件中的这些值，以自定义开发模式体验。

自定义方式包括：

1. 更改要监视的目录或文件
2. 调整防抖时间，以控制系统响应变更的速度
3. 添加或修改执行命令，以满足项目需求

### 使用浏览器进行开发

虽然 Wails v2 完全支持使用浏览器进行开发，但这造成了许多困惑。由于 WebView 并不提供所有浏览器 API，在浏览器中正常运行的应用程序不一定能在桌面应用程序中正常运行。

对于侧重 UI 的开发工作，在 v3 中仍可灵活使用浏览器：在开发模式下访问位于`http://localhost:9245`的 Vite URL。这样，你在处理样式和布局时便可使用强大的浏览器开发工具。请注意，在此模式下，Go 绑定<em>将无法工作</em>。准备测试绑定和事件等功能时，只需切换到桌面视图，确保所有功能都能在生产环境中正常运行。
