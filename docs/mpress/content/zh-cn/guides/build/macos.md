---
title: "macOS 打包"
description: "打包 Wails 应用，以便在 macOS 上分发"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## macOS 私有 API

Wails v3 默认使用 macOS 公共 API。若要启用需要使用 Apple 未公开 API 的功能，只需在构建应用时添加`private_mac_apis`这一个 Go 构建标签：

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

若直接使用 Go 构建，请使用 `go build -tags private_mac_apis .`（生产环境则使用 `-tags production,private_mac_apis`）。依赖私有行为的现有应用必须添加此标签，才能保留这些行为。该标签仅适用于 macOS 桌面构建。

有关受影响功能和选项值的完整清单、公共构建的确切回退行为、Liquid Glass 样式映射以及检查器构建组合，请参阅[macOS 私有 API](/guides/build/private-macos-apis/)。如果没有此标签，仅限私有 API 的操作不会执行任何动作；公共 Go API 保持不变。

## 应用包

将应用打包为标准 macOS `.app` 应用包：

```bash
wails3 package GOOS=darwin
```

这会创建包含以下内容的 `bin/<AppName>.app`：

- 位于 `Contents/MacOS/` 中的已编译二进制文件
- 位于 `Contents/Resources/` 中的应用图标（来自 `icons.icns`；如果存在资产目录 `Assets.car`，则来自该目录）
- 包含应用元数据的 `Info.plist`

## 应用包资源

`Contents/Resources/` 是存放随 macOS 应用分发的只读文件的标准位置。可在其中存放较大的模板、种子数据、媒体、语言包或其他应按需打开而非通过 `embed` 编译进 Go 可执行文件的载荷。

Wails 已将应用图标放在此目录中。若要添加自己的文件，请将其放入 `build/resources/` 等源目录，然后在 `build/darwin/Taskfile.yml` 中为 `create:app:bundle` 任务添加复制步骤：

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

如果使用 Taskfile 的 `darwin:run` 任务，请将等效命令添加到其 `run` 任务，并将目标设为 `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`。

### 从 Go 读取资源

导入 macOS 平台包：

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

对于小文件，请使用 `LoadResource`：

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

对于较大的文件，请使用 `ResourceFS`。它返回一个根目录为 `Contents/Resources` 的 `io/fs.FS`，因此调用方无需先将整个资源加载到 Go 字节切片中，即可打开并以流式方式读取资源：

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

资源名称是相对于 `Contents/Resources`、以斜杠分隔的路径。除非可执行文件从 `.app/Contents/MacOS` 运行，否则 `ResourceFS` 和 `LoadResource` 会返回 `mac.ErrNotInAppBundle`。

应将应用包资源视为不可变内容。更改已签名应用内的文件会使其代码签名失效；下载、生成或可由用户编辑的数据应改为存储在用户的 Application Support 目录中。

### 通用二进制文件

为 Apple 芯片和 Intel Mac 同时构建：

```bash
wails3 task darwin:package:universal
```

这会创建一个可在两种架构上原生运行的 `.app`。可在任何平台上构建通用二进制文件——在 Linux 和 Windows 上会自动使用 `wails3 tool lipo`。

## 自定义应用包

编辑 `build/darwin/Info.plist` 可自定义以下内容：

- 应用包标识符（`CFBundleIdentifier`）
- 应用名称和版本
- 最低 macOS 版本
- 文件关联
- URL 方案

应用图标根据 `build/` 目录中的资产生成。请使用 `generate:icons` 任务：

```bash
wails3 task common:generate:icons
```

此任务使用 `build/appicon.png` 生成 `darwin/icons.icns` 和 `windows/icon.ico`。在 macOS 上还可以提供 `build/appicon.icon`（Icon Composer 格式）：该任务会传递 `-iconcomposerinput appicon.icon -macassetdir darwin`，以便从 `.icon` 文件生成 `Assets.car` 和 `darwin/icons.icns`（在非 macOS 平台上会跳过）。存在 `Assets.car` 时，请运行 `update:build-assets` 任务，以相应更新 `Info.plist` 和 `CFBundleIconName`：

```bash
wails3 task common:update:build-assets
```

若要从 `build/` 目录手动运行图标命令：

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## 代码签名

对应用签名以供分发：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

在 `build/darwin/Taskfile.yml` 中配置签名：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### 公证

对于在 Mac App Store 之外分发的应用，Apple 要求进行公证：

```bash
wails3 task darwin:sign:notarize
```

首先存储凭据。可以运行交互式向导（`wails3 setup signing`），也可以直接调用 `notarytool`：

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

在 `build/darwin/Taskfile.yml` 中配置：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

有关详细信息，请参阅[为应用签名](/guides/build/signing/)。

## DMG 安装程序

Wails 3随附的模板提供`wails3 task darwin:package:dmg`。它会先创建`.app`，然后使用 DMG 库构建带样式的 DMG。默认情况下，DMG 使用带有红色龙形标志和 WAILS 字标的 Wails 品牌渐变背景。

```bash
wails3 task darwin:package:dmg
```

较底层的 `darwin:create:dmg` 任务从现有的 `.app` 应用包创建 DMG，并可直接通过 Taskfile 配置：

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### 默认布局

生成的 DMG 包含：

- 左侧的应用包
- 右侧的 `Applications` 链接
- 尺寸为 540×380 像素的 Finder 窗口
- 图标大小为 96 点，每个图标下方均有标签
- 来自 `build/darwin/dmg-background.png` 的 Wails 品牌背景

应用图标和 `Applications` 图标的位置取决于配置的窗口尺寸，因此更改 `DMG_WINDOW_WIDTH` 或 `DMG_WINDOW_HEIGHT` 后，默认的双图标布局仍会保持成比例的间距。为获得最佳效果，请使用像素尺寸与 Finder 窗口相同的背景图像。

### 替换 DMG 资产

`build/darwin/` 下生成的文件是普通的项目资源，可以替换：

- `DMG_BACKGROUND` 控制 Finder 窗口内容背后显示的图像。
- `DMG_VOLUME_ICON` 控制已挂载宗卷所显示的图标。
- `DMG_FILE_ICON` 控制生成的 `.dmg` 文件在 Finder 中显示的图标。

宗卷图标和 DMG 文件图标是两项独立资源。替换应用程序图标不会自动替换其中任何一项。

### 添加额外文件

使用 `DMG_FILES` 将安装程序脚本、发行说明、许可证或其他资源与应用程序一起包含进来。其值是以逗号分隔的 `name=path` 对列表：

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=` 前的名称是在 DMG 内显示的文件名。`=` 后的路径是项目中的源文件。开头和结尾的空白字符会被忽略。

每个显示名称都必须唯一。额外文件不能替换打包工具已创建的条目，包括应用程序包或 `Applications` 条目。如果名称冲突，打包会报错失败，而不会生成损坏的 DMG。

@note{type="note"}
由于 DMG 创建过程使用 macOS 磁盘映像和 Finder 工具，因此仅支持在 macOS 上创建 DMG。可以在其他平台上创建交叉编译的 `.app` 包，但最终的 DMG 必须在 Mac 上生成。

@end

## 故障排除

### “应用程序已损坏，无法打开”

该应用程序未经签名。请使用 Developer ID 证书对其签名，或者让用户绕过 Gatekeeper：

```bash
xattr -cr /path/to/YourApp.app
```

### 公证失败

常见问题：

- **凭据无效**：重新运行 `xcrun notarytool store-credentials`（或 `wails3 setup signing`）
- **需要强化运行时**：如有需要，请确保授权文件中包含 `com.apple.security.cs.allow-unsigned-executable-memory`
- **缺少时间戳**：签名过程应自动包含时间戳

### 交叉编译的应用程序无法运行

交叉编译的 macOS 二进制文件未经签名。请将其传输到 Mac，并在测试前签名：

```bash
codesign --force --deep --sign - YourApp.app
```
