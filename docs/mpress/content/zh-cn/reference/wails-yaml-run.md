---
title: "YAML 运行配置"
description: "wails3 run 的本地运行默认值、平台覆盖项和启动行为。"
---

## 可用性

本参考文档介绍实验性的运行扩展。不保证已发布的 CLI 中会提供此扩展。

`wails3 run`使用本地清单默认值构建并启动应用程序。
由清单驱动的执行需要启用
[`WAILS_EXP_USE_WAKE=1`选择加入项](/experimental/yaml-builds#enable-yaml-builds)。
未找到清单时，该命令会在当前目录中运行`go run .`，
且不要求启用该选择加入项。清单无效或无法读取时会报错；
绝不会因此选择回退方案。如果存在清单但 YAML 已禁用，
该命令会说明如何启用 YAML。

生成的[`wails.yaml`字段参考](/reference/wails-yaml-fields)列出了
字段类型和默认值。在发布版本声称支持原生移动平台和 macOS 应用包之前，
必须在相应主机上验证这些平台能否接受该应用。

## 命令约定

`wails3 run`会查找最近的`wails.yaml`，解析一个目标的配置，
通过现有的缓存管线进行构建，然后启动构建结果。
查找会从工作目录开始向上进行，在遇到`go.mod`边界并停止之前检查清单。
清单所在目录既是项目根目录，也是桌面应用程序的工作目录。

默认目标是主机操作系统和体系结构。移动目标
必须显式选择。交叉编译并不表示桌面可执行文件能在主机上运行；
不兼容的桌面目标会在构建前失败。这也包括使用`CGO_ENABLED=0`构建的可执行文件：
仅使用 Go 成功完成交叉构建并不会提供模拟器，也不会改变本地启动约定。
使用[`wails3 build`进行交叉编译](/experimental/yaml-builds/#cross-compile-desktop-applications)，
并在目标系统上运行或安装生成的文件。
Android 请使用`--target android/arm64 --device <adb-serial>`或`--target android/arm64 --emulator <avd-name>`。
所选体系结构必须与设备或模拟器匹配。iOS 模拟器请使用`--target ios/arm64 --device <udid>`；若使用实体设备，请添加`--destination device`。iOS 的规划和执行需要 macOS 和 Xcode。

使用 YAML 时，构建会采用嵌入式资源和生产构建路径，并自动提供
`production`标签。实时重新加载仍由`wails3 dev`负责；
`dev.tags`不会作用于`wails3 run`。

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan`会解析并显示所选目标、生效的构建标签、
可执行或可安装制品、参数和工作目录，但不会
构建、安装或启动。计划会列出环境覆盖项的键名，但省略其值。

## CLI 选项

| 选项 | 含义 |
| --- | --- |
| `--tags <tags>` | 以逗号分隔的其他 Go 构建标签。 |
| `--skip-build` | 从 YAML 输出目录导入现有桌面可执行文件，然后直接启动，无需编译。 |
| `--target <os>/<arch>` | 选择一个目标；默认为主机操作系统和体系结构。 |
| `--plan` | 显示解析后的设置，但不构建或启动。允许生成跨平台桌面目标的计划。 |
| `--device <id>` | Android adb 序列号或 iOS 模拟器/设备标识符。 |
| `--emulator <name>` | 要启动或复用的 Android 虚拟设备；不能与`--device`组合使用。 |
| `--destination <kind>` | iOS 目标位置：`simulator`（默认）或`device`。 |
| `-- <args...>` | 替换已配置的应用程序参数。单独使用`--`会清除这些参数。 |
| `--appargs <string>` | 兼容拼写形式：使用 POSIX 风格的引号规则解析一个字符串。不能重复使用，也不能与`--`之后的参数组合使用。 |

移动平台选择器要求存在清单和匹配的移动目标。

`--skip-build` 要求当前存在有效的 YAML 项目和兼容的桌面可执行文件。此模式会跳过前端处理、绑定生成、编译、UPX 和构建钩子；但配置的启动参数和环境仍然适用。导入的可执行文件会经过验证并复制到 Wake 生成的工作区中，且不会修改其源文件。macOS GUI 应用程序会使用这些输入文件和当前软件包资源进行组装。请参阅[预构建输入路径和要求](/experimental/yaml-builds/#use-prebuilt-desktop-executables)。
构建标签无法更改预构建的可执行文件，因此在此模式下，`--tags`和配置的构建标签均不起作用。

## 字段

可选的`run`映射提供共享的启动默认值。`targets["<os>"]`或`targets["<os>/<arch>"]`内可选的`run`
映射提供特定于平台或体系结构的默认值。无论映射顺序如何，
特定于体系结构的设置优先于操作系统级设置。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | 仅用于运行构建的其他 Go 标签。 |
| `run.args` | `list(string)` | `[]` | 传递给应用程序的参数。 |
| `run.environment` | `map(string)` | `{}` | 运行时环境覆盖项，与编译器环境分开。 |
| `target["<os>"].run.tags` | `list(string)` | `[]` | 所选平台的其他运行标签。 |
| `target["<os>"].run.args` | `list(string)` | 继承 | 显式设置时替换共享参数；`[]`会清除这些参数。 |
| `target["<os>"].run.environment` | `map(string)` | `{}` | 按键覆盖共享的运行时环境值。 |

参数以参数向量形式传递，不进行 shell 展开。根据现有 YAML 仅允许字面量的约定，
环境值是字面字符串。

## 平台示例

此片段与常规项目配置配合使用。`example_private_api`
是示例所需实际私有 API 标签的占位符。

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Linux 环境设置仅适用于此示例，并且需要 Wayland
会话；它不是 Wails 的默认设置。仅在 macOS 上选择私有 API 标签。
现有的`target.tags`适用于生产构建，包括运行构建。
`target.run.tags`将附加标签限制为仅用于运行构建。

特定于体系结构的覆盖项可以进一步细化操作系统级默认值：

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

与前面的 Linux 映射结合使用时，此设置会保留`GDK_BACKEND`、添加
特定于体系结构的标签，并清除共享的应用程序参数。

## 解析规则

| 设置 | 解析顺序（优先级从低到高） |
| --- | --- |
| 标签 | 合并`build.tags`、所选`target.tags`、自动添加的`production`、`run.tags`、所选`target.run.tags`和 CLI `--tags`；每个标签仅保留首次出现的项。 |
| 运行时环境 | 继承启动命令的 shell 环境，先叠加`run.environment`，再叠加所选`target.run.environment`。空值会设置为空字符串。 |
| 应用程序参数 | 先使用共享的`run.args`，再由显式设置的目标参数替换，最后由`--`之后的显式参数替换。单独使用`--`会清除已配置的参数。 |

在所选目标内，先解析操作系统级设置，再解析特定于体系结构的设置：
标签会累积，环境值按键合并，而显式提供的参数会替换粒度较低的列表。
运行构建还会为相应的移动目标提供`android`或`ios`。

绑定生成和 Go 编译接收相同的解析后标签。构建
缓存键包含这些标签。仅更改应用程序参数或运行时
环境不会使编译缓存失效；下次启动会使用新值。

平台映射用于自定义平台；映射存在并不意味着不支持其他平台。
`project.supported_platforms`会限制构建和`run`所接受的平台；省略该列表或使用空列表则允许所有平台。例如，`project`内的`supported_platforms: [darwin]`会使仅适用于 macOS 的示例在 Linux 上提前失败。

## 平台启动行为

| 平台 | 启动行为 |
| --- | --- |
| Linux 和 Windows | 启动可执行文件，连接标准输入/输出，处理终止事件，并传播应用程序退出状态。 |
| macOS | 组装应用程序包，并在连接标准输入/输出的情况下执行其中的`Contents/MacOS`二进制文件。向 Launch Services 注册该应用程序包，以支持 URL 方案和文件关联。签名和原生生命周期验收需要在 macOS 上进行验证。 |
| Android | 构建 APK，选择设备或模拟器，安装并启动应用程序。分别报告构建、安装和启动失败。 |
| iOS | 构建应用程序包，选择模拟器或设备，满足签名和预置要求，然后安装并启动。设备构建需要适当的凭据和 Apple 工具链。 |

Android 会在部署前拒绝应用程序参数和运行时环境覆盖。如果同时支持 Android，应将仅适用于桌面的共享设置放入桌面目标映射中。Android 运行会附加应用程序日志，并在附加结束时停止应用程序。iOS 会将参数传递给所选的启动器。模拟器通过`SIMCTL_CHILD_`变量接收环境覆盖；物理设备则通过`devicectl --environment-variables`接收 JSON 映射。移动端退出状态描述的是启动器或日志监控结果，而不是可移植的应用程序退出代码。

## 纯 Go 示例和后备方案

对于使用预先存在或嵌入资源的纯 Go 示例，可禁用前端阶段：

```yaml
frontend:
  disabled: true
```

这会跳过前端依赖项安装、绑定生成和前端构建。需要生成绑定类的示例可以改为在`frontend.bindings`下设置`interfaces: false`。

无清单后备方案会执行`go run [-tags <tags>] . [application arguments]`，而不会添加`production`或`devtools`。它遵循 Go 自身的退出状态行为。`--`之后的参数会传递给应用程序，包括`--help`等标志。`--plan`显示的运行时值仅限于环境变量键名；不会打印对应值。后备方案会显示 Go 命令和工作目录。

## 开发环境中的应用程序参数

存在活动 YAML 项目时，`wails3 dev`可接收独立于生产环境`run`默认值的应用程序参数。`dev.args`提供共享默认值；`targets["<os>"]`或`targets["<os>/<arch>"]`中嵌套的`dev`映射会替换该目标的共享默认值。无论映射顺序如何，架构设置的优先级都高于操作系统级设置。显式设置`args: []`会清除继承的参数。

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args`和`run.tags`不会影响开发环境启动。现有的`dev.tags`构建策略保持不变。应用程序会以清单目录为基准解释相对应用程序路径。

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

`--`之后的参数会替换已配置的参数，即使该列表为空也是如此。该分隔符之后的`--help`属于应用程序。`--appargs=""`也会清除默认值。兼容性选项接受一个带引号的字符串；如果将其与`--`、位置应用程序参数或第二个`--appargs`一起使用，则会报错。格式错误的引号会在启动会话前被拒绝。

兼容性解析器在所有主机上都使用相同的 POSIX 风格引号规则。它不会执行 shell 命令或展开变量；调用方 shell 在将文本传递给 Wails 之前可能已经将其展开。`--`之后的参数已经是彼此独立的值，不会再进行分词。对于 Windows 路径，建议使用这种形式，例如在 PowerShell 中：

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

会话会在源代码重新构建期间保留解析后的参数向量。仅更改`dev.args`或所选目标的参数时，会使用缓存的可执行文件重启后端；不会重启前端。重新加载清单后，CLI 覆盖仍然有效。替换失败时，会保留或恢复先前的可执行文件及其原有参数。应用程序参数不是编译器标志，也不会纳入编译缓存键。`dev --plan`会在不启动应用程序的情况下，将有效参数与有限构建计划一并显示。

Linux、Windows 和 macOS 会将参数向量直接传递给可执行文件。iOS 通过模拟器或设备启动器传递该向量。Android 没有受支持的等效参数通道，因此会在发现设备或开始构建之前拒绝非空的开发环境应用程序参数。可以将仅适用于桌面的默认值放入桌面目标映射中，也可以在`targets.android.dev`下设置`args: []`来清除这些默认值。

这些开发选项需要活动的`wails.yaml`并在 YAML 中显式启用。它们不会重写旧版开发路径中的 Taskfile 或监视器 shell 命令。
