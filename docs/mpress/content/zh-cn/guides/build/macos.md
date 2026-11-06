---
title: "macOS 打包"
description: "打包 Wails 应用以在 macOS 上分发"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## 私有 macOS API

Wails v3默认使用公开的 macOS API。若要启用需要使用未公开 Apple API 的功能，请使用单个 Go 构建标签`private_mac_apis`构建应用：

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

直接使用 Go 构建时，请使用`go build -tags private_mac_apis .`（生产环境则使用`-tags production,private_mac_apis`）。依赖私有行为的现有应用必须添加此标签，才能保留这些行为。该标签仅适用于 macOS 桌面构建。

有关受影响功能及选项值的完整清单、公开 API 构建的确切回退行为、Liquid Glass 样式映射以及检查器构建组合，请参阅[私有 macOS API](/guides/build/private-macos-apis/)。如果没有该标签，仅限私有 API 的操作不会执行任何操作；公开的 Go API 保持不变。

## 应用程序包

将应用打包为标准 macOS `.app`包：

```bash
wails3 package GOOS=darwin
```

这会创建包含以下内容的`bin/<AppName>.app`：

- 位于`Contents/MacOS/`中的已编译二进制文件
- 位于`Contents/Resources/`中的应用图标（来自`icons.icns`；如果存在资产目录`Assets.car`，则来自该目录）
- 包含应用元数据的`Info.plist`

### 同步文件夹中的项目

某些 macOS 文件提供程序会向`.app`目录添加 Finder 元数据。这会导致`codesign`拒绝该程序包，或在签名后使其签名失效。如果提供程序会再次添加元数据，仅删除元数据并不足以解决问题。

生成的 macOS 构建任务会检查此行为。检测到后，Wails 会将项目生成的输出目录移至`~/Library/Application Support/Wails/build-output/<project>-<hash>`，并用符号链接替换原始输出目录。现有输出文件会予以保留，而`bin/<AppName>.app`仍是打包、签名和开发模式所使用的路径。CLI 执行移动时会输出本地存储位置。

设置`WAILS_MACOS_OUTPUT_ROOT`可选择另一个未同步的本地存储位置。迁移根目录必须与项目位于同一文件系统。对于项目外部显式指定的`BIN_DIR`，将按配置使用。现有项目需要更新后的`build/darwin/Taskfile.yml`才能使用此处理机制。

分发时，请复制`bin/<AppName>.app`或创建 DMG。仅复制`bin`符号链接不会包含构建产物。如果直接运行`wails3 tool sign`命令时，其输入文件夹被重新附加元数据，该命令会报告错误；请先将该程序包移至未同步的文件夹，再对其签名。

## 程序包资源

`Contents/Resources/`是存放随 macOS 应用分发的只读文件的标准位置。对于大型模板、种子数据、媒体、语言包或其他应按需打开而非通过`embed`编译进 Go 可执行文件的有效负载，请使用此位置。

Wails 已将应用图标放在此目录中。若要添加自己的文件，请将其放入`build/resources/`等源目录，然后在`build/darwin/Taskfile.yml`的`create:app:bundle`任务中添加复制步骤：

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

如果使用 Taskfile 的`darwin:run`任务，请将等效命令添加到其`run`任务中，并将目标设为`{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`。

### 从 Go 读取资源

导入 macOS 平台包：

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

对于小文件，请使用`LoadResource`：

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

对于较大的文件，请使用`ResourceFS`。它会返回一个以`Contents/Resources`为根目录的`io/fs.FS`，因此调用方无需先将整个资源加载到 Go 字节切片中，即可打开并流式读取资源：

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

资源名称是相对于`Contents/Resources`、以斜杠分隔的路径。除非可执行文件从`.app/Contents/MacOS`运行，否则`ResourceFS`和`LoadResource`会返回`mac.ErrNotInAppBundle`。

应将程序包资源视为不可变内容。修改已签名应用程序中的文件会使其代码签名失效；请改为将下载、生成或可由用户编辑的数据存储在用户的 Application Support 目录中。

### 通用二进制文件

为 Apple 芯片和 Intel Mac 同时构建：

```bash
wails3 task darwin:package:universal
```

这会创建一个可在两种架构上原生运行的`.app`。可以在任何平台上构建通用二进制文件；在 Linux 和 Windows 上会自动使用`wails3 tool lipo`。

## 自定义程序包

编辑`build/darwin/Info.plist`以自定义：

- 程序包标识符（`CFBundleIdentifier`）
- 应用名称和版本
- 最低 macOS 版本
- 文件关联
- URL 方案

应用图标根据`build/`目录中的资产生成。请使用`generate:icons`任务：

```bash
wails3 task common:generate:icons
```

此任务使用`build/appicon.png`生成`darwin/icons.icns`和`windows/icon.ico`。在 macOS 上，还可以提供`build/appicon.icon`（Icon Composer 格式）：该任务会传递`-iconcomposerinput appicon.icon -macassetdir darwin`，从`.icon`文件生成`Assets.car`和`darwin/icons.icns`（在非 macOS 平台上跳过）。如果存在`Assets.car`，请运行`update:build-assets`任务，以相应更新`Info.plist`和`CFBundleIconName`：

```bash
wails3 task common:update:build-assets
```

若要从`build/`目录手动运行图标命令：

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## 代码签名

对应用进行签名以供分发：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

在`build/darwin/Taskfile.yml`中配置签名：

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

首先，存储凭据。可以运行交互式向导（`wails3 signing`），也可以直接调用`notarytool`：

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

在`build/darwin/Taskfile.yml`中配置：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

有关详细信息，请参阅[为应用程序签名](/guides/build/signing/)。

## DMG 安装程序

Wails 随附的3模板提供了`wails3 task darwin:package:dmg`。它会先创建`.app`，再使用 DMG 库构建带样式的 DMG。默认情况下，DMG 使用带有红色龙形标志和 WAILS 字标的 Wails 品牌渐变背景。

```bash
wails3 task darwin:package:dmg
```

较低层级的`darwin:create:dmg`任务会从现有`.app`程序包创建 DMG，并且可以直接在 Taskfile 中进行配置：

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

- 左侧的应用程序包
- 右侧的`Applications`链接
- 尺寸为540×380像素的 Finder 窗口
- 图标大小为96点，每个图标下方均有标签
- 来自`build/darwin/dmg-background.png`的 Wails 品牌背景

应用和`Applications`图标的位置会根据配置的窗口尺寸确定，因此更改`DMG_WINDOW_WIDTH`或`DMG_WINDOW_HEIGHT`后，默认的双图标布局仍会保持成比例的间距。为获得最佳效果，请使用与 Finder 窗口像素尺寸相同的背景图像。

### 替换 DMG 资产

`build/darwin/`下生成的文件是普通项目资产，可以替换：

- `DMG_BACKGROUND`控制 Finder 窗口内容后方显示的图像。
- `DMG_VOLUME_ICON`控制已挂载卷显示的图标。
- `DMG_FILE_ICON`控制生成的`.dmg`文件在 Finder 中显示的图标。

卷宗图标和 DMG 文件图标是两个独立的资源。替换应用程序图标不会自动替换其中任何一个。

### 添加额外文件

使用`DMG_FILES`在应用程序旁附加安装程序脚本、发行说明、许可证或其他资源。该值是由逗号分隔的`name=path`对列表：

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=`之前的名称是在 DMG 内显示的文件名。`=`之后的路径是项目中的源文件。开头和结尾的空白字符会被忽略。

每个显示名称都必须唯一。额外文件不能替换打包工具已经创建的条目，包括应用程序包或`Applications`条目。如果名称冲突，打包会失败并报错，而不会生成损坏的 DMG。

@note{type="note"}
DMG 创建仅受 macOS 支持，因为它使用 macOS 的磁盘映像和 Finder 工具。可以在其他平台上创建交叉编译的`.app`包，但最终的 DMG 必须在 Mac 上生成。

@end

## 故障排除

### “App 已损坏，无法打开”

该应用未签名。可以使用 Developer ID 证书为其签名，用户也可以绕过 Gatekeeper：

```bash
xattr -cr /path/to/YourApp.app
```

### 公证失败

常见问题：

- **凭据无效**：重新运行`xcrun notarytool store-credentials`（或`wails3 signing`）
- **需要强化运行时**：如有需要，请确保授权配置中包含`com.apple.security.cs.allow-unsigned-executable-memory`
- **缺少时间戳**：签名过程应自动包含时间戳

### 交叉编译的应用无法运行

交叉编译的 macOS 二进制文件未签名。请将其传输到 Mac，并在测试前签名：

```bash
codesign --force --deep --sign - YourApp.app
```
