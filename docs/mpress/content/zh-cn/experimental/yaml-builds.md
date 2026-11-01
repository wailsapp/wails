---
title: "YAML 构建"
description: "使用一个 wails.yaml 文件配置并运行实验性 Wails v3 原生构建系统。"
sidebar:
  order: 2
  badge:
    text: "实验性"
    variant: caution
---

@note{type="caution" title="实验性预览"}
首个获准发布的版本仅支持字面量配置。此实验还对有界生命周期钩子进行原型验证，以便在提议发布这些钩子之前评估其形式。自定义阶段、YAML 表达式和类型化 Wails 工具调用仍将延期实现。在明确启用经审查的迁移方案之前，现有项目继续使用 Taskfile。
@end

## 启用 YAML 构建

在每个运行实验性命令的终端中设置`WAILS_EXP_USE_WAKE`。对于 Bash 或 Zsh：

```bash
export WAILS_EXP_USE_WAKE=1
```

对于 PowerShell：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

此变量按是否存在来决定是否选择启用：即使值为空或为`false`，也会启用该功能。取消设置该变量（在 Bash/Zsh 中使用`unset WAILS_EXP_USE_WAKE`，在 PowerShell 中使用`Remove-Item Env:WAILS_EXP_USE_WAKE`）可禁用此实验。启用`wails.yaml`后，命令要求重新启用 Wake，并且绝不会回退到 Taskfile。
如果未设置该变量，YAML 构建/开发路由和 YAML 项目生成功能将被禁用，且`migrate`、`config`和`eject`命令不可用。尚未迁移的项目继续使用独立的嵌入式 Task 运行时。

## 配置项目

选择启用后，在项目根目录添加`wails.yaml`。该文件的存在会选用原生构建系统，并使 Wails 忽略此项目中的所有 Taskfile。运行以下所有命令时都应保持选择启用状态。模式 modeline 可在支持 YAML 语言服务器的编辑器中启用补全和行内验证。

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

运行默认的宿主机构建或完整的具名发布请求：

```bash
wails3 build
wails3 build release
```

对于一次性构建，可以选择以逗号分隔的目标和格式，而无需添加配置文件：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## 交叉编译桌面应用程序

要使用实验性的纯 Go 桌面后端，请在构建环境中设置`CGO_ENABLED: "0"`。Wake 为 macOS 和 Linux GUI 程序选择`purego`标签，并在 Windows、macOS 和 Linux 宿主机上直接使用 Go。在此模式下，构建桌面可执行文件不需要 Docker、模拟器、C 编译器或目标平台 SDK。应用程序依赖项也必须支持`CGO_ENABLED=0`。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

第二条命令会合并两种 macOS 架构，并使用 Go 创建未签名的 DMG。它可以在任何受支持的桌面宿主机上运行。

系统会分别为每个可执行文件选择解析后的模式。声明的`build.environment.CGO_ENABLED`、目标环境或二进制文件环境会覆盖调用命令的终端中的值。否则，显式选择的`purego`标签、CLI/服务器程序或 Windows 目标默认为`0`；macOS 和 Linux GUI 程序默认为`1`。`--plan --json`会报告每个可执行文件解析后的编译器、CGo 模式和标签。交叉编译会创建目标文件；运行不兼容的桌面目标仍需使用相应的宿主机。

使用`CGO_ENABLED=1`时，Wake 会首先使用兼容的原生工具链或已配置的`CC`交叉编译器；进行 Windows 交叉构建时，如果 Zig 可用，则接着使用 Zig；然后在支持的平台上使用现有的 Docker/Podman 回退方案。显式的`CC`必须提供目标编译器、库和 sysroot；仅有 Zig 并不能提供 Linux GTK/WebKit sysroot。打包和签名可能还需要下文所述的其他工具。移动端目标需要 CGo 及其平台 SDK。

纯 Go 编译不会捆绑目标系统的系统库。Linux GUI 应用程序（包括 GUI 配套可执行文件）仍需要其构建标签所选择的 GTK/WebKitGTK 运行时。必须由软件包依赖项或适当的部署环境提供该运行时。默认情况下，仅含 CLI/服务器的包不会添加这些 GUI 依赖项；仍可使用显式的依赖项覆盖设置。

## 构建多个可执行文件

项目定义主程序。将`project.package`设置为其 Go 包（默认为`.`），并将`project.kind`设置为`gui`、`cli`或`server`（默认为`gui`）。其他程序使用`binaries`下的具名条目：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

每个程序都有自己的编译输入和缓存输入。系统会跟踪导入的项目包和嵌入文件；编辑某个独立程序不会导致重新构建其他程序。桌面平台支持其他可执行文件，并会将其纳入默认的应用程序包和安装程序。由用户提供的替换包文件自行控制其内容。

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

具名配置文件表示一个完整请求；请将这些选择覆盖项用于匿名构建。`--jobs`用于限制执行工作进程的数量。仅选择 CLI/服务器程序时会省略前端工作。

要在不运行应用程序代码的情况下读取版本，请同时配置`project.version_file`和`project.version_symbol`。该文件必须包含静态 Go 字符串声明。随后可以省略字面量`project.version`。在项目或具名二进制文件上设置`commit_symbol: main.Commit`，可将仓库修订版本注入该可执行文件的字符串变量中。如果请求注入提交信息，而仓库元数据缺失，则会报错。

`upx_level: 1`到`9`可选择性地压缩发布版 Linux 和 Windows 可执行文件；设置为零会禁用压缩。必须安装 UPX。开发构建和 macOS 构建不进行压缩。编译或压缩失败时，会保留此前已发布的程序。

## 为分发进行打包和签名

### 使用预构建的桌面可执行文件

如果其他构建系统或 CI 作业已经生成可执行文件，请使用`--skip-build`：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

将每个可执行文件放入`build.output`（默认为`bin`），主可执行文件使用`project.binary_name`，配套可执行文件则使用各`binaries`条目的名称。Windows 文件名以`.exe`结尾。多目标调用会从`<os>-<arch>`子目录中读取各目标的文件。对于 macOS GUI，如果`bin/<name>`存在，Wake 会从中读取原始可执行文件，否则从`bin/<name>.app/Contents/MacOS/<name>`读取。`darwin/universal`输入必须同时包含 amd64 和 arm64 可执行文件切片。

Wake 会验证文件类型和目标架构，将可执行文件复制到其生成的工作区，并跟踪其内容以用于包缓存。可执行文件输入缺失、不兼容或为符号链接时，操作会失败，而不会重新构建或恢复旧输入。系统不会覆盖提供的可执行文件，也不会对其执行 chmod。任何可执行文件发生变化，都会使使用该文件的包缓存失效。

此模式会跳过前端安装/构建、绑定生成、Go 编译、UPX 以及`before_build`/`after_build`钩子。它仍会生成包资源、组装 macOS 应用程序包，并运行已配置的打包/签名钩子。包和签名工具仍受下表所列要求约束；不需要 Go 和前端构建工具。包元数据来自当前清单，包括已配置的`version_file`；不会从二进制文件推断版本。请使用相匹配的构建和打包设置，包括 GUI 与 CLI 类型。

也支持`wails3 package --skip-build`。使用`wails3 run`时，启动参数和环境设置仍然适用，而且可执行文件必须能在宿主机上运行。此模式要求存在已启用的 YAML 项目，并且仅支持生产环境的桌面目标；它不适用于开发会话或移动端打包。

创建未签名包与进行签名对宿主机和工具的要求不同：

| 输出 | 创建所需的宿主机和工具 | 签名 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS 或 Linux；使用 Go 写入器，无需 Windows SDK。 | Windows 宿主机和签名工具。 |
| Windows NSIS 安装程序 | Windows、macOS 或 Linux；需要宿主机原生的`makensis`。 | Windows 宿主机和签名工具。 |
| macOS 应用程序和 DMG | Windows、macOS 或 Linux；使用 Go 进行组装、通用二进制文件合并，并通过可移植 DMG 写入器生成 DMG。 | 需要配备 Apple 工具和凭据的 macOS。 |
| Linux DEB、RPM、Arch 和 Void | Windows、macOS 或 Linux；使用原生 Go 包写入器。 | DEB/RPM 使用各自的外部签名工具；Arch 使用下文所述的 Go 签名器。Void 输出不签名。 |
| Linux AppImage | Linux 主机与目标架构相同；需要外部 linuxdeploy 和 GTK 部署插件。 | 没有内置的 AppImage 签名流程。 |

Windows MSIX 归档直接使用 Go 创建；创建未签名的软件包不需要 Windows SDK 或 Docker。Windows 签名仍需要 Windows 主机及其签名工具。NSIS 使用外部 NSIS 编译器，并通过`packages.nsis.defines`接受字面量定义。在 macOS 上，仅当 NSIS 编译器继承的区域设置不是 UTF-8 时，Wails 才会为该编译器子进程提供 UTF-8 区域设置；终端和应用程序环境不会改变。当打包的所有可执行文件都是 CLI 或服务器程序时，标准 NSIS 安装程序会省略 WebView2 引导程序。若主程序或伴随程序包含 GUI，则会保留 WebView2。自定义 NSIS 脚本会保留引导程序，以便与其自身的运行时安装逻辑兼容。

Linux 支持 AppImage、DEB、RPM、Arch Linux 和 Void 软件包。AppImage 仍以外部 linuxdeploy 和 GTK 插件为边界。设置`packages.appimage.update_info`可请求生成 zsync 文件，设置`packages.appimage.checksum: true`可请求生成 SHA256 文件。伴随文件会从软件包缓存中恢复，并与 AppImage 一同列入产物收据。目前，AppImage 打包要求使用 Linux，并且主机与目标架构必须相同。其 GTK 检测无法发现纯 Go GUI 动态加载的库，因此此适配器尚不支持部署纯 Go GUI。对于这种构建模式，请使用一种 Linux 软件包管理器格式。

Arch 软件包会在原生软件包元数据中保留预发布版本，例如`1.2.3-beta.2`会变为`1.2.3beta.2-1`，其排序位置位于最终版本之前。Void 依赖项接受诸如`gtk4`的裸软件包名称，并将其规范化为`gtk4>=0`；显式版本比较和完整的`name-version_revision`值会予以保留。请使用 Void 软件包名称，而不是库文件名。

要创建 Arch 软件包及其分离式签名，请使用导出的 ASCII 封装 PGP 私钥文件，并可选择指定一个口令环境变量：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

通过本地环境或 CI 密钥存储区设置指定的环境变量，然后运行`wails3 build release`。如果密钥未加密，请省略`credential`。Go 签名器可在所有受支持的桌面主机上运行，且不需要 GPG。它会发布原始`.pkg.tar.zst`及同目录下的`.pkg.tar.zst.sig`；两者会分别记录在产物收据中。密钥内容不会暂存或存入构建缓存。签名及其发布结果绝不会从缓存中复用。签名失败时会保留之前发布的一对文件；如果安装任一文件失败，发布过程会回滚这两个文件。这并不能使两次文件系统重命名对并发读取者呈现原子性。

macOS 应用程序包包含伴随程序。设置`darwin.microphone_usage_description`可添加用途说明，并在未提供自定义权利文件时添加音频输入权利。自定义 plist 和权利源文件保持不变。签名和公证在 macOS 上使用现有的 Apple 工具与凭据；公证需要真实的公证配置文件。签名后的输出会保留其平台扩展名，例如`notes.signed.app`和`notes.signed.dmg`。

## 开发控制命令

运行`wails3 dev`期间，输入命令并按 Enter：

| 命令 | 操作 |
| --- | --- |
| `r` | 强制重新构建并重启，包括未发生更改的后端。 |
| `b` | 重新生成绑定并重新构建。 |
| `f` | 重启前端。 |
| `c` | 清空终端。 |
| `q` | 停止会话及其子进程。 |
| `h` | 显示控制命令。 |

构建失败时，当前应用程序会继续运行。控制命令需要交互式终端；通过管道传入内容不会改变会话行为。

## 运行前检查工作计划

计划检查所解析的不可变图与执行时使用的相同。文本供人阅读；JSON 是带版本的机器可读契约。

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## 不执行构建，仅验证配置

请在 CI、编辑器任务中或提交配置前使用专用的预检命令。未指定配置文件时，它会验证基础清单及每个具名配置文件；指定配置文件时，它只检查该项完整请求：

```bash
wails3 config check
wails3 config check release
```

语法、架构、语义、路径、目标、格式、签名、规划器以及由配置决定的主机能力错误，均使用与 Build 和 Dev 相同的诊断渲染器。诊断信息会标明字段和源码范围，输出源码行及插入符，并在有帮助时提供修正提示。诊断信息的源码摘录会隐去环境变量值。仅由命令行选项导致的主机错误仍作为普通命令错误处理，而不会被误报为无效 YAML。

## 使用项目自有的资源

清单中的路径指向您拥有的文件，例如图标、平台清单、软件包模板和维护脚本。Wails 绝不会编辑这些源文件。工具使用这些文件前，Wails 会将其复制到一次性的`.wails/`工作区。最终产物存放在`build.output`下，经过验证的收据会写入`.wails/artifacts/receipt.json`。

## 部署到 Android

生产 Android 配置文件会生成 AAB。对于本地部署，Wails 会在内部构建针对特定 ABI 的调试 APK，使用 ADB 安装并启动应用：

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

未显式指定选择器时，交互式终端会列出已连接的设备和已配置的 Android 虚拟设备供您选择。当有多个可用目标时，自动化流程必须提供`--device`或`--emulator`。

附加应用级日志流并监控目标，直到应用退出或您按下 Ctrl+C：

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

附加运行会跟随应用重启，并报告目标断开连接或转为离线状态。默认情况下，模拟器会继续运行，以便快速迭代。`--stop-emulator`必须显式指定且要求提供`--logs`；附加会话结束时，它只会停止作为目标的模拟器。如果 Wails 启动了新模拟器，但等待启动就绪失败或被取消，它会自动清理该进程。

## 运行受限的生命周期钩子

实验性钩子接口会在每个已声明的阶段调用一个项目自有脚本。钩子不是 shell 命令：Wails 会直接执行该文件，除非设置了`directory`，否则从项目根目录执行。

在 Unix 主机上，脚本必须可执行，并以 shebang 开头。Windows 主机通过固定的解释器适配器接受`.cmd`、`.bat`和`.ps1`脚本；Wails 在任一平台上都不会解析内联 shell 命令。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

支持的阶段为`before_build`、`after_build`、`before_package`、`after_package`、`before_sign`和`after_sign`。构建钩子会根据情况为共享项目运行一次，或为每个目标各运行一次。软件包和签名钩子会在某个目标所请求的所有格式外围形成屏障，而屏障内针对各格式的操作仍会并发执行。

钩子会收到一个由 Wails 管理的环境变量：`WAILS_HOOK_CONTEXT_FILE`。它指向一个受保护的 JSON 文件，其中包含带版本的调用上下文：

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

上下文包含解析后的路径，但绝不会复制继承的进程环境或密钥。Wails 会在`.wails/hooks/context/`下以仅所有者可访问的文件权限原子创建该文件。操作成功或取消后，该文件会被删除。操作失败时，该文件会保留以供诊断，其路径也会包含在错误中。非零退出码会导致常规 Plan 步骤失败，取消操作会终止钩子进程组。

默认情况下，钩子会在每次调用时运行。仅当脚本具有完整声明的输入和输出时，才选择启用缓存：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

脚本及其可执行模式是隐式输入。多个输出必须位于同一个非根目录中，该目录就是 Wails 记录和恢复的有界产物。输出根目录不得包含脚本或任何已声明的输入；缓存的输出树只能包含普通文件和目录，绝不能包含符号链接或特殊文件。始终运行的钩子会使其下游操作不可缓存，因为无法安全地为未声明的副作用生成指纹。可缓存的钩子通过其声明的产物标识维持下游缓存。

## 迁移现有项目

创建并检查一份尚未启用的草稿：

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

不带参数、行为保守且调用项目自有脚本的生命周期任务会转换为钩子。其他可达的自定义 Taskfile 行为会被报告为阻塞项。审核
`wails.migrated.yaml`后，请显式激活它：

```bash
wails3 migrate --activate
```

激活操作会以原子方式创建`wails.yaml`。它不会删除、重命名或编辑 Taskfile、`build/config.yml`、脚本或资源。无效的 YAML 会以 YAML 错误告终；绝不会回退到旧版路径。

## 具体化所有默认值

若要检查或复制完整的已解析配置，请写入一个未激活的快照：

```bash
wails3 eject
```

此操作会创建`wails.ejected.yaml`，绝不会重写活动配置。

## 兼容期

`wails3 package`和`wails3 sign`是同一规划器的已弃用别名；新的自动化流程应使用带配置文件或格式的`wails3 build`。推出过程受版本发布约束：首个版本将 YAML 设为新建项目的默认路径，下一个版本会针对旧版 Taskfile 路由发出警告，再下一个版本则会移除 Taskfile 执行功能及这些别名。每次过渡都要求相应版本的原生主机验收矩阵全部通过。

有关每个字段、类型、默认值和适用规则，请参阅生成的[`wails.yaml`字段参考](/reference/wails-yaml-fields/)。

[`wails3 run`配置参考](/reference/wails-yaml-run/)介绍了本地启动默认值、平台覆盖项以及无清单时的`go run .`回退机制。
