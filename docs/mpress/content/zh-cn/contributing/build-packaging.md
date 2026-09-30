---
title: "构建与打包流水线"
description: "介绍运行`wails3 build`时的底层流程、如何生成跨平台二进制文件，以及如何为各操作系统生成安装程序。"
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build`有意保持<strong>轻量</strong>：它是一个 Taskfile 包装器，会将额外的构建标签转发给宿主项目的`build`任务。繁重的工作由项目自身的`build/Taskfile.yml`（由`wails3 init`生成）、`internal/commands/build-assets.go`（管理烘焙时资源）、`internal/packager`（Linux nfpm 打包）以及`internal/commands/appimage.go`、`internal/commands/msix.go`、`internal/commands/dmg/dmg.go`、`internal/commands/dot_desktop.go`（各平台的安装程序）完成。

本页涵盖：

1. 实际的 CLI 入口点
2. 由 Taskfile 驱动的构建流程
3. 资源烘焙与构建信息注入
4. 各平台的打包后端
5. 自定义流水线
6. 故障排除

---

## 1. 实际的 CLI 入口点

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`：

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build`只公开<strong>一个</strong>标志——`--tags`（转发为`EXTRA_TAGS=`）。`wails3 build`上<strong>没有</strong>`-platform`、`-o`、`-skipbindings`、`-skip-package`、`-package`、`-ldflags`、`-verbose`、`-debug`、`-devbuild`、`-icon`或`-clean`标志。交叉编译、输出路径、图标等在<strong>`Taskfile.yml`</strong>、<strong>`build/config.yml`</strong>以及配套的`wails3 generate icons`/`wails3 generate build-assets`命令中配置。

`build/build.json`<strong>不</strong>属于 v3——其配置由`Taskfile.yml`和`build/config.yml`组成。

---

## 2. 由 Taskfile 驱动的构建流程

新初始化的项目附带一个`build/Taskfile.yml`，其中的命名空间大致如下：

| 命名空间 | 任务（节选） |
| --- | --- |
| `darwin:` | `build`、`build:universal`、`package`、`run`、`dev` |
| `windows:` | `build`、`package`、`run`、`dev` |
| `linux:` | `build`、`package`、`run`、`dev` |
| `common:` | `update:build-assets`、`generate:icons`、`generate:syso` |

默认情况下，`wails3 build`会调用宿主操作系统的`build`命名空间；随后，项目的 Taskfile 会使用宿主平台专用标志通过 shell 调用`go build`。为其他操作系统构建时，需要直接运行相应任务（例如`wails3 task darwin:build:universal`），而不是向`wails3 build`传递标志。

默认输出目录为<strong>`bin/<APP_NAME>`</strong>（没有`build/bin/`前缀）。

---

## 3. 烘焙时资源与构建信息

| 事项 | 文件 |
| --- | --- |
| 构建资源的生成/更新 | `internal/commands/build-assets.go` |
| 构建信息输出器（CLI：`wails3 tool buildinfo`） | `internal/commands/tool_buildinfo.go`——输出信息；它<strong>不是</strong>`ldflags`注入器 |
| 生产环境存根 | `internal/assetserver/build_production.go`——`//go:build production` |
| 前端捆绑包 | 通过应用自身包中的`//go:embed`嵌入（例如位于`main.go`旁边） |
| Windows 资源（`.syso`） | `internal/commands/syso.go`——生成`rsrc_windows_<arch>.syso` |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG 输入 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI 不会自动为应用烘焙`bundled_assetserver.go`——`internal/assetserver/bundled_assetserver.go`是<strong>手写的</strong>，用于包装`bundledassets/`下嵌入的 JS 运行时。

---

## 4. 打包后端

### Linux

Linux 打包由<strong>nfpm</strong>（而非`fpm`）驱动：

- `internal/packager/packager.go`包装`github.com/goreleaser/nfpm/v2`并公开`CreatePackageFromConfig(pkgType, configPath, output)`/`CreatePackageFromConfigWriter(...)`。
- 生成的项目在`internal/commands/`下附带 nfpm 风格的`myapp.DEB`、`myapp.RPM`和`myapp.ARCHLINUX`配置（由`wails3 tool package`使用）。
- AppImage 生成功能位于`internal/commands/appimage.go`，它会调用`linuxdeploy` + `linuxdeploy-plugin-gtk`（插件随附于`internal/commands/linuxdeploy-plugin-gtk.sh`）。

`wails3 build`上<strong>没有</strong>`-package deb`/`rpm`标志。请使用`wails3 tool package`或平台专用的 Taskfile 目标。

### macOS

- 项目 Taskfile 中的`darwin:package`会生成`.app`捆绑包。
- DMG 资源位于`internal/commands/dmg/`下；`darwin:package`完成后，项目可以使用`hdiutil`将应用程序包封装为 DMG（较新模板中的 Taskfile 包含`dmg`辅助任务）。
- CFBundle 标识符、版本和版权信息来自执行`wails3 init`时使用的`-product*`标志以及`build/config.yml`。

### Windows

- Windows 打包以<strong>MSIX</strong>为目标（不是 WiX/MSI）。完整工作流请参阅`internal/commands/msix.go`和`internal/commands/webview2/`。
- 既<strong>没有</strong>`internal/commands/packager.go`，也<strong>没有</strong>`internal/commands/windows_resources/`目录。
- 可选的可执行文件代码签名通过`wails3 tool sign`（Authenticode）运行——请参阅`internal/commands/sign.go`。

---

## 5. 自定义流水线

| 需求 | 方法 |
| --- | --- |
| 额外的构建标签 | `wails3 build --tags myFeature,otherTag` |
| 代码检查器／构建前步骤 | 向`build/Taskfile.yml`添加任务，并让特定于操作系统的`build`任务依赖该任务 |
| 交叉编译 | 运行相应的操作系统任务（例如`wails3 task linux:build`）——不存在`-platform`标志 |
| 跳过打包 | 只需运行`build`任务；`package`是独立任务 |
| 自定义打包工具 | 将配置放在`internal/commands/myapp.*`下，并使用`-config <file>`调用`wails3 tool package` |
| 剥离符号 | 编辑`darwin:/windows:/linux:`中的`build`任务，将`-ldflags "-s -w"`直接传给`go build`——`wails3 build`本身没有`-ldflags`标志 |

所有 Taskfile 目标都遵循 Wails 发布的环境变量（`APP_NAME`、 `WAILS_VITE_PORT`、`FRONTEND_DEVSERVER_URL`……），因此自定义任务可以依赖这些变量。

---

## 6. 故障排除

| 症状 | 可能的原因 | 解决方法 |
| --- | --- | --- |
| **`ld: framework not found WebKit`（mac）** | 缺少 Xcode CLI 工具 | `xcode-select --install` |
| **生产构建中出现空白窗口** | 前端构建失败或 SPA 路由问题 | 检查`frontend/dist/index.html`是否存在，并确认资源处理程序会回退到该文件 |
| **缺少 MSIX 打包工具** | 未安装`WebView2` SDK／MSIX 工具 | 运行`wails3 task install:msix:tools` |
| **`linuxdeploy`未找到** | PATH 中缺少插件 | 安装`linuxdeploy`，并通过 CLI 的自动安装步骤运行`internal/commands/linuxdeploy-plugin-gtk.sh` |

`wails3 build`没有`-verbose`标志。设置`TASK_X_VERBOSE=1`（Taskfile），或 直接检查任务目标以查看正在执行的命令。

---

## 7. 关键源代码索引

| 关注点 | 文件 |
| --- | --- |
| 构建包装器 | `internal/commands/task_wrapper.go`（`Build`、`Package`、`SignWrapper`、`wrapTask`） |
| 构建资源生成 | `internal/commands/build-assets.go`（`GenerateBuildAssets`、`UpdateBuildAssets`） |
| 构建信息输出程序 | `internal/commands/tool_buildinfo.go` |
| AppImage 构建器 | `internal/commands/appimage.go` |
| Linux 打包（nfpm） | `internal/packager/packager.go`、`internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows MSIX | `internal/commands/msix.go`、`internal/commands/webview2/` |
| Windows 资源生成器 | `internal/commands/syso.go` |
| macOS DMG 资源 | `internal/commands/dmg/` |
| `.desktop`生成器 | `internal/commands/dot_desktop.go` |
| 版本常量 | `internal/version/version.go` |

排查构建失败时，请将此表留作参考。

---

现在，你已经全面了解了从<strong>源代码</strong>到<strong>安装程序</strong>的整个过程。简而言之， `wails3 build`本身只是一个轻量包装器——几乎所有自定义都在项目的 `Taskfile.yml`／`build/config.yml`中完成，或通过明确的`wails3 generate …`／ `wails3 tool …`子命令完成。祝发布顺利！
