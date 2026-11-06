---
title: "构建与打包流水线"
description: "介绍运行`wails3 build`时的底层过程、如何生成跨平台二进制文件，以及如何为各操作系统生成安装程序。"
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build`有意保持**轻量**：它是一个Taskfile包装器，将额外的构建标签转发给宿主项目的`build`任务。繁重的工作由项目自身的`build/Taskfile.yml`（由`wails3 init`生成）、`internal/commands/build-assets.go`（管理烘焙时资源）、`internal/packager`（Linux nfpm打包）以及`internal/commands/appimage.go`、`internal/commands/msix.go`、`internal/commands/dmg/dmg.go`和`internal/commands/dot_desktop.go`（各平台的安装程序）完成。

本页涵盖：

1. 实际的CLI入口点
2. 由Taskfile驱动的构建流程
3. 资源烘焙与构建信息注入
4. 各平台的打包后端
5. 自定义流水线
6. 故障排除

---

## 1. 实际的CLI入口点

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

`flags.Build`仅公开<strong>一个</strong>标志——`--tags`（转发为`EXTRA_TAGS=`）。`wails3 build`<strong>没有</strong>`-platform`、`-o`、`-skipbindings`、`-skip-package`、`-package`、`-ldflags`、`-verbose`、`-debug`、`-devbuild`、`-icon`或`-clean`标志。交叉编译、输出路径和图标等在<strong>`Taskfile.yml`</strong>、<strong>`build/config.yml`</strong>以及配套的`wails3 generate icons`/`wails3 generate build-assets`命令中配置。

`build/build.json`**不**属于v3——配置由`Taskfile.yml`和`build/config.yml`组成。

---

## 2. 由Taskfile驱动的构建流程

新初始化的项目附带一个`build/Taskfile.yml`，其命名空间大致如下：

| 命名空间 | 任务（节选） |
| --- | --- |
| `darwin:` | `build`、`build:universal`、`package`、`run`、`dev` |
| `windows:` | `build`、`package`、`run`、`dev` |
| `linux:` | `build`、`package`、`run`、`dev` |
| `common:` | `update:build-assets`、`generate:icons`、`generate:syso` |

默认情况下，`wails3 build`会调用宿主操作系统的`build`命名空间；项目的Taskfile随后使用宿主平台专用标志执行`go build`。若要为其他操作系统构建，应直接运行对应任务（例如`wails3 task darwin:build:universal`），而不是向`wails3 build`传递标志。

默认输出目录为<strong>`bin/<APP_NAME>`</strong>（没有`build/bin/`前缀）。

---

## 3. 烘焙时资源与构建信息

| 事项 | 文件 |
| --- | --- |
| 构建资源的生成/更新 | `internal/commands/build-assets.go` |
| 构建信息输出器（CLI：`wails3 tool buildinfo`） | `internal/commands/tool_buildinfo.go`——输出信息；它**不是**`ldflags`注入器 |
| 生产环境存根 | `internal/assetserver/build_production.go`——`//go:build production` |
| 前端捆绑包 | 通过应用自身包中的`//go:embed`嵌入（例如位于`main.go`旁边） |
| Windows资源（`.syso`） | `internal/commands/syso.go`——生成`rsrc_windows_<arch>.syso` |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG输入 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI不会为你的应用自动烘焙`bundled_assetserver.go`——`internal/assetserver/bundled_assetserver.go`是**手写的**，用于包装`bundledassets/`下嵌入的JavaScript运行时。

---

## 4. 打包后端

### Linux

Linux打包由**nfpm**驱动（而非`fpm`）：

- `internal/packager/packager.go`包装了`github.com/goreleaser/nfpm/v2`，并公开`CreatePackageFromConfig(pkgType, configPath, output)`/`CreatePackageFromConfigWriter(...)`。
- 生成的项目在`internal/commands/`下附带nfpm风格的`myapp.DEB`、`myapp.RPM`和`myapp.ARCHLINUX`配置（由`wails3 tool package`使用）。
- AppImage生成逻辑位于`internal/commands/appimage.go`中，它会调用`linuxdeploy`和`linuxdeploy-plugin-gtk`（插件位于`internal/commands/linuxdeploy-plugin-gtk.sh`）。

`wails3 build`上**没有**`-package deb`/`rpm`标志。请使用`wails3 tool package`或平台专用的Taskfile目标。

### macOS

- 项目Taskfile中的`darwin:package`会生成`.app`捆绑包。
- DMG资源位于`internal/commands/dmg/`下；`darwin:package`完成后，项目可使用`hdiutil`将捆绑包封装成DMG（较新模板中的Taskfile包含`dmg`辅助任务）。
- CFBundle标识符、版本和版权信息来自`wails3 init`期间使用的`-product*`标志以及`build/config.yml`。

### Windows

- Windows打包以**MSIX**为目标（而非WiX/MSI）。完整工作流请参阅`internal/commands/msix.go`和`internal/commands/webview2/`。
- **没有**`internal/commands/packager.go`，并且**没有**`internal/commands/windows_resources/`目录。
- 可选的可执行文件代码签名通过`wails3 tool sign`（Authenticode）执行——请参阅`internal/commands/sign.go`。

---

## 5. 自定义流水线

| 需求 | 方法 |
| --- | --- |
| 额外的构建标签 | `wails3 build --tags myFeature,otherTag` |
| 代码检查器/构建前步骤 | 向`build/Taskfile.yml`添加任务，并让操作系统专用的`build`任务依赖该任务 |
| 交叉编译 | 运行对应的操作系统任务（例如`wails3 task linux:build`）——不存在`-platform`标志 |
| 跳过打包 | 只运行`build`任务；`package`是独立任务 |
| 自定义打包器 | 在`internal/commands/myapp.*`下放置配置文件，并使用`-config <file>`调用`wails3 tool package` |
| 剥离符号 | 编辑`darwin:/windows:/linux:`中的`build`任务，将`-ldflags "-s -w"`直接传递给`go build`——`wails3 build`本身没有`-ldflags`标志 |

所有 Taskfile 目标都遵循 Wails 发布的环境变量（`APP_NAME`、`WAILS_VITE_PORT`、`FRONTEND_DEVSERVER_URL`等），因此自定义任务可以依赖这些变量。

---

## 6. 故障排除

| 症状 | 可能的原因 | 解决方法 |
| --- | --- | --- |
| **`ld: framework not found WebKit`（mac）** | 缺少 Xcode CLI 工具 | `xcode-select --install` |
| **生产构建中出现空白窗口** | 前端构建失败或 SPA 路由问题 | 检查`frontend/dist/index.html`是否存在，并确保资源处理程序会回退到该文件 |
| **缺少 MSIX 打包工具** | 未安装`WebView2` SDK / MSIX 工具 | 运行`wails3 task install:msix:tools` |
| 找不到<strong>`linuxdeploy`</strong> | PATH 中缺少插件 | 安装`linuxdeploy`，并通过 CLI 的自动安装步骤运行`internal/commands/linuxdeploy-plugin-gtk.sh` |

`wails3 build`没有`-verbose`标志。设置`TASK_X_VERBOSE=1`（Taskfile），或直接检查任务目标以查看正在执行的命令。

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

排查构建失败时，请随时参考此表。

---

现在，你已经全面了解从**源代码**到**安装程序**的整个流程。简而言之，`wails3 build`本身只是一个轻量包装器——几乎所有自定义操作都在项目的`Taskfile.yml` / `build/config.yml`中完成，或通过明确的`wails3 generate …` / `wails3 tool …`子命令完成。祝发布顺利！
