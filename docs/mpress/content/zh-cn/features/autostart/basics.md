---
title: "自动启动"
description: "注册应用，使其在 macOS、Windows 和 Linux 上于用户登录时启动"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 自动启动

`app.Autostart`可注册应用，使其在用户登录时自动启动。它会为每个平台选择合适的原生机制，并解析符号链接指向的安装路径（Homebrew、Scoop），因此二进制文件升级后注册不会失效。

注册会在<strong>下次登录</strong>时生效，而不是立即生效。

## 快速开始

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

使用默认选项注册应用，使其在登录时启动。

```go
func (m *AutostartManager) Enable() error
```

反复调用`Enable`是安全的——每次都会覆盖注册，因此，如果已持久保存用户的偏好设置，可以在每次启动时调用它。

### `EnableWithOptions`

使用自定义选项进行注册。

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`：**

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Identifier` | `string` | 覆盖自动派生的注册 ID。请参阅下文的“标识符”。 |
| `Arguments` | `[]string` | 登录时启动应用时追加到可执行文件路径后的额外参数（例如`--hidden`）。 |

### `Disable`

移除自动启动注册。如果应用尚未注册，则返回`nil`——禁用操作是幂等的。

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

报告注册是否存在。此操作速度很快，但不会验证已注册的路径。

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

返回完整的注册状态。

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`：**

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Enabled` | `bool` | 注册是否存在。 |
| `Path` | `string` | 注册产物在磁盘上的位置（plist 路径、`.desktop`路径或注册表子项）。当`Enabled`为 false 时为空。 |
| `Strategy` | `AutostartStrategy` | 注册应用时使用的机制（请参阅[平台行为](#heading-2)）。 |

## 平台行为

@tabs{sync-key="platform"}
[macOS]
根据应用的打包方式，将使用以下两种机制之一：

- **macOS 13+，已捆绑的`.app`**：`SMAppService.mainAppService`。适用于沙盒应用和 Mac App Store 构建。不会出现 TCC 自动化提示（过去的 AppleScript 方案会触发该提示）。
- **早于13的 macOS，或未捆绑的二进制文件**：将包含`RunAtLoad=true`的 LaunchAgent plist 写入`~/Library/LaunchAgents/<identifier>.plist`。

`Status()`返回`AutostartStrategySMAppService`或`AutostartStrategyLaunchAgent`，以便调用方判断使用了哪条路径。

当应用从未捆绑版本升级到已捆绑版本时，`Status()`会检查这两条路径，`Disable()`会进行清理，以免遗留的 LaunchAgent 继续启动旧版本。

[Windows]
在`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`下添加一个注册表值，将值名称设为自动启动`Identifier`，并将数据设为带引号的可执行文件路径及任何`Arguments`。

参数引号处理遵循`CommandLineToArgvW`规则（引号前的反斜杠会加倍），因此包含空格或引号的路径可以正确往返转换。

`Status().Strategy`返回`AutostartStrategyRegistryRun`。

[Linux]
将 XDG 自动启动项写入`$XDG_CONFIG_HOME/autostart/<identifier>.desktop`（默认为`~/.config/autostart/`），其中包含：

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec`字段按照[freedesktop.org Desktop Entry 规范](https://specifications.freedesktop.org/desktop-entry-spec/)进行转义——保留字符（`"`、`` ` ``、`$`、`\\`）使用反斜杠转义；当值包含空白字符时，使用双引号将其括起。

`Status().Strategy`返回`AutostartStrategyXDGAutostart`。

[iOS / Android / 服务器]
不支持。所有方法均返回`ErrAutostartNotSupported`。使用`errors.Is(err, application.ErrAutostartNotSupported)`可明确检测这种情况：

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## 标识符

如果`Options.Identifier`为空，则根据应用名称派生默认值：

| 平台 | 默认标识符 |
| --- | --- |
| macOS（已捆绑） | 应用的捆绑包标识符，例如`com.example.MyApp` |
| macOS（未捆绑） | `wails.autostart.<slug>`，其中`<slug>`派生自`application.Options.Name` |
| Windows | `application.Options.Name`的 slug（转换为小写、移除非`A-Za-z0-9._-`字符，并将空格替换为连字符） |
| Linux | 与 Windows 使用相同的 slug |

标识符必须匹配`^[A-Za-z0-9._-]+$`，且长度不得超过200个字符。建议在 macOS 上使用反向 DNS 格式（这与 launchd Label 的惯常写法一致）。

覆盖`AutostartOptions.Identifier`后，同一标识符将复用于 Windows 上的注册表值名称和 Linux 上的`.desktop`文件名，因此一个字符串即可跨平台标识该注册项。

## 失效项检测

`Disable()`和`Status()`通过<strong>将已注册的可执行文件路径与`os.Executable()`（解析所有符号链接后）</strong>进行匹配来定位注册项，而不是查找标识符。这意味着：

- <strong>在不同版本之间更改标识符是安全的。</strong>只要可执行文件路径相同，仍可通过`Status()`找到旧注册项，并由`Disable()`将其清理。
- <strong>位于不同路径的第二个应用副本不会覆盖第一个副本的注册项。</strong>每个二进制文件位置都会单独跟踪。
- <strong>通过符号链接安装的应用（Homebrew、Scoop）可保持稳定。</strong>匹配前会对`os.Executable()`应用`filepath.EvalSymlinks`，因此 Homebrew 升级时即使替换链接目标，也不会遗留该注册项。

以下情况<em>不在</em>此机制的处理范围内：如果用户将二进制文件移动到无关路径或重命名，旧注册项将成为孤立项（它会指向当前已不存在的文件）。从固定安装路径发布的应用无需担心这一点；以便携式单文件二进制形式发布的应用应在移动自身之前调用`Disable()`，或者始终通过固定的符号链接启动。

## 示例

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

[`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart)中提供了一个包含状态、启用和禁用按钮的完整可运行示例。
