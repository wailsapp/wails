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

第一次使用 Wake？请先阅读[入门教程](/experimental/wake/)。
其中介绍了该实验的限制、所有项目和构建设置界面、首次本地构建，以及经过审核的迁移。本页将更详细地介绍各项构建、打包和配置工作流。

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

要在浏览器中创建并配置项目，请在启用选择加入功能后运行`wails3 init -ui`，并按照[Wake 教程](/experimental/wake/)操作。项目选项保存在`wails.yaml`中；要更改现有项目，请编辑该文件。`wails3 setup`用于配置共享首选项、工具和签名凭据。

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

第二条命令会合并两种 macOS 架构，并使用 Go 创建未签名的 DMG。它可以在任何受支持的桌面主机上运行。

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

IPA 有效负载和用于 macOS 公证的 ZIP 均由 Go 写入，并保留 bundle 路径、可执行权限和相对符号链接。这两项操作都不需要外部`zip`或`ditto`可执行文件。这些归档包含应用程序文件，而非 Finder 元数据或资源分支；macOS 签名本身就会移除这些元数据，因为它们会导致已签名的 bundle 失效。

| 输出 | 创建所需的宿主机和工具 | 签名 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS 或 Linux；使用 Go 写入器，无需 Windows SDK。 | 可在任何桌面主机上使用导出的 RSA PFX 通过 Go 签名；证书存储区/提供程序身份需要 Windows SignTool。 |
| Windows NSIS 安装程序 | Windows、macOS 或 Linux；需要宿主机原生的`makensis`。 | 可在任何桌面主机上使用导出的 RSA PFX 通过 Go 为 PE 签名；证书存储区/提供程序身份需要 Windows SignTool。 |
| macOS 应用程序和 DMG | Windows、macOS 或 Linux；使用 Go 进行组装、通用二进制文件合并，并通过可移植 DMG 写入器生成 DMG。 | 使用导出的 RSA PKCS#12身份通过 Go 签名，或进行临时签名；钥匙串身份仍使用 macOS 工具。API 密钥公证和票据装订通过 Go 运行。 |
| Linux DEB、RPM、Arch 和 Void | Windows、macOS 或 Linux；使用原生 Go 包写入器。 | 导出的 OpenPGP 密钥使用原生 DEB/RPM/Arch 签名器。旧版 DEB/RPM 密钥 ID 需要外部工具。Void 输出不签名。 |
| Linux AppImage | Windows、macOS 或 Linux；通过 Go 部署并写入 SquashFS，另需目标 AppImage 运行时。跨平台打包 GUI 需要提供 Linux 运行时根目录树。 | 没有内置的 AppImage 签名流程。 |

对于已签名的 DMG，Wake 通常先为应用程序的私有副本签名，再由 Go 写入器压缩映像，最后为映像签名。这样可避免使用`hdiutil`挂载、展开并重新压缩 DMG。每次请求都会执行签名；不会复用构建缓存中的已签名软件包。具有`after_package`或`before_sign`钩子的项目仍沿用原有的先打包后签名流程，包括`hdiutil`，因此这些钩子对未签名映像所做的修改会得到保留。这个兼容钩子的映像编辑路径需要 macOS 和`hdiutil`。

要进行可移植的 Apple 签名，请设置`darwin.signing.identity: "-"`以使用临时签名，或使用`darwin.signing.certificate`指定导出的 RSA PKCS#12身份。`darwin.signing.credential`指定其密码环境变量的名称。证书签名使用共享的 Go PFX/CMS 实现；不会自动从钥匙串导出私钥。使用生成的测试身份时，瘦可执行文件和通用可执行文件、嵌套 app bundle、框架及 DMG 均已通过 Apple 严格的原生签名检查，包括篡改拒绝检查。

进行公证时，导出的 Developer ID 身份还需要兼容 RFC 3161 的`darwin.signing.timestamp_server`。共享验证器接受其支持的 SHA-2时间戳配置，并拒绝不受支持的旧版算法。临时签名无法通过公证。原生证书格式验证与真正的 Developer ID 信任及 Apple 是否接受提交相互独立。

Go 装订器会验证 Apple 票据签名，并确认其涵盖所有架构，然后再将票据附加到私有副本。它使用观测到的票据交付协议，而 Apple 并未将该协议作为受支持的 API 发布。它不能取代 Gatekeeper 当前的吊销检查。App 和已签名的 DMG 输出均已通过原生`stapler validate`检查及严格签名检查。

Windows MSIX 归档直接使用 Go 创建；创建未签名的软件包无需 Windows SDK 或 Docker。导出的 RSA PFX 身份可在任何桌面主机上直接使用 Go 对 Windows 可执行文件和各个 MSIX 软件包进行签名。将 `windows.signing.certificate` 设置为项目相对路径下的 PFX，并将 `windows.signing.credential` 设置为其密码环境变量。支持 AES/PBKDF2 和旧版 3DES 导出；RC2 导出需要重新导出，或回退到原生 Windows SignTool。证书存储区/提供程序身份和 MSI 签名需要 Windows SignTool。默认时间戳服务为 DigiCert；发布已签名的输出前会验证时间戳。有关格式限制，请参阅[签名指南](/guides/build/signing/)。
NSIS 使用外部 NSIS 编译器，并通过 `packages.nsis.defines` 接受字面量定义。在 macOS 上，如果 NSIS 编译器继承的区域设置不是 UTF-8，Wails 只会为该编译器子进程提供 UTF-8 区域设置；终端和应用程序环境不会改变。当所有打包的可执行文件都是 CLI 或服务器时，标准 NSIS 安装程序会省略 WebView2 引导程序。若主可执行文件或配套可执行文件是 GUI，则会保留 WebView2。自定义 NSIS 脚本会保留该引导程序，以兼容其自身的运行时安装逻辑。

Linux 支持 AppImage、DEB、RPM、Arch Linux 和 Void 软件包。Wake 直接通过 Go 组装 AppImage，包括部署 GTK/WebKit、压缩 SquashFS 和生成 zsync。它不会调用 linuxdeploy、GTK 插件、appimagetool 或 mksquashfs。系统会为目标架构下载固定版本且经过校验和验证的 AppImage 类型 2运行时；也可通过`packages.appimage.runtime_path`改为选择相对于项目的运行时 ELF。

在 Linux 上，默认库来源是主机的根文件系统。要为其他架构打包 GUI，或从 Windows/macOS 打包 GUI，请将`packages.appimage.runtime_root`设为相对于项目的目标 Linux 根目录树，其中应包含相应的库、已编译 schema 和模块缓存。这里指完整的根目录布局（包括`usr/`），而不是`usr/`目录本身。请提供来自计划支持的最旧 Linux 发行版的运行时。Wake 无需执行目标程序即可读取 ELF 导入及显式 GTK/WebKit 配置文件。它会捆绑选定的依赖闭包，但不包括 glibc 或动态加载器。运行时更新后，库和资源的内容哈希会使软件包缓存失效。

支持纯 Go GUI 构建。对于使用`gtk3`构建的 GUI 可执行文件，Wake 选择 GTK 3；其他情况则选择 GTK 4，包括 CLI/服务器的 GUI 伴随程序。同一映像中的所有 GUI 可执行文件必须使用同一 GTK 版本。对于预构建二进制文件，可通过`packages.appimage.gtk_version`明确选择`3`或`4`；导入的 CGo 库必须与此选择一致。

设置`packages.appimage.update_info`可请求生成 zsync 文件，设置`packages.appimage.checksum: true`可请求生成 SHA256 文件。Sidecar 文件会从软件包缓存恢复，并与 AppImage 一同列入制品收据。使用提供的 GTK/WebKit 目录树时，已验证可在 Windows 和 macOS 上进行打包；生成的纯 Go GUI 映像通过了 Linux 上的运行时检查。在一次条件匹配的 CGo GUI 构建中，原生路径耗时8.77秒，而 linuxdeploy 耗时26.19秒。原生 gzip 映像为133 MiB，而旧 zstd 路径生成的映像为114 MiB，因此这里衡量的是整个打包路径，而不仅仅是压缩。

Arch 软件包会在原生软件包元数据中保留预发布版本，例如`1.2.3-beta.2`会变为`1.2.3beta.2-1`，其排序位置位于最终版本之前。Void 依赖项接受诸如`gtk4`的裸软件包名称，并将其规范化为`gtk4>=0`；显式版本比较和完整的`name-version_revision`值会予以保留。请使用 Void 软件包名称，而不是库文件名。

在任何受支持的桌面主机上，DEB、RPM 和 Arch 签名均可使用导出的 ASCII 装甲格式或二进制 OpenPGP 私钥，而无需 GnuPG、`dpkg-sig`或`rpmsign`。通过`linux.signing.certificate`配置密钥文件，并通过`credential`配置其密码短语环境变量的名称。DEB 和 RPM 会嵌入签名；Arch 则发布分离式签名。私钥仅由签名步骤读取，不会复制到软件包资源或可复用的缓存制品中。

包含裸 GnuPG 密钥标识符的现有 DEB/RPM 配置会继续使用外部签名器。要消除此依赖，请导出私钥并改为配置其路径。请使用类似`./signing/release.asc`的路径；如果密钥文件缺失、不可读或无效，操作将失败，而不会回退到密钥环。DEB 签名继续使用`dpkg-sig`软件包签名格式；这不能取代对 apt 仓库发布元数据的签名。

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

macOS app bundle 包含伴随程序。设置`darwin.microphone_usage_description`可添加用途说明文本；若未提供自定义 entitlements 文件，还会添加音频输入 entitlement。自定义 plist 和 entitlements 源文件保持不变。导出的 RSA PKCS#12身份可在任何桌面主机上使用 Go 签名；`identity: "-"`用于选择可移植的临时签名。钥匙串身份在 macOS 上继续使用`codesign`。票据装订通过 Go 运行。公证提交可以使用现有的 macOS 公证钥匙串配置文件，也可以使用 App Store Connect 团队 API 凭据：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

API 模式通过 Go 提交、上传和轮询，无需`notarytool`。请设置全部三个字段并省略`credential`；不要将私钥纳入版本控制。仅在签名时读取密钥，且不会将其复制到生成的资源或缓存中。此模式支持团队 API 密钥和最大5 GiB 的归档；尚未实现分段上传和个人 API 密钥。提交被拒绝时会保留提交 ID，并将 Apple 的诊断日志保存到私有临时文件，错误中会报告该文件的路径。提交被取消或拒绝时，不会替换此前发布的输出。Apple 接受提交后，Wake 会检索并验证其票据，再通过 Go 进行装订。使用导出的身份和团队 API 凭据时，可在 Linux、Windows 或 macOS 上完成签名、提交和装订。钥匙串配置文件仍需要 macOS 和`notarytool`。

可移植证书签名支持采用 AES/PBKDF2 或旧版3DES 加密的 PKCS#12文件中的 RSA 身份。使用 RC2 加密的导出文件必须以受支持的密码重新导出。用于公证的可移植签名需要明确指定 SHA-2 RFC 3161时间戳服务；Go 验证器会拒绝旧版 SHA-1时间戳签名。票据服务遵循观测到的 Apple 协议，并固定信任 Apple Root G3；Gatekeeper 仍是分发信任和吊销的最终判定方。

签名后的输出文件会保留其平台扩展名，例如
`notes.signed.app` 和 `notes.signed.dmg`。

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
