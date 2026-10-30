---
title: "可自我更新的 Wails 应用"
description: "构建一个可从 GitHub Releases 自我更新的 Wails v3应用——涵盖从`wails3 init`到签名发布验证和辅助模式替换的全过程。"
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

在本教程中，你将为一个全新的 Wails v3应用添加应用内更新程序。完成后，该应用将能够：

- 按需检查 GitHub Releases（也可选择通过定时器检查）。
- 下载与当前操作系统和架构匹配的资产。
- 根据下载的字节验证 SHA-256摘要（也可选择验证 Ed25519 签名）。
- 在框架的默认更新窗口中显示发布说明。
- 替换正在运行的二进制文件并重新启动，且无需分发单独的辅助可执行文件。

我们将使用**GitHub Releases**作为更新源，因为它免费且不需要任何基础设施。同样的模式也适用于[keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen)和[Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast)。完成本教程后，请参阅[更新程序指南](/guides/updater/)。

@note{type="tip" title="前提条件"}
- Go 1.26或更高版本
- 已安装`wails3` CLI（`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`）
- 一个你可以向其推送发布版本的 GitHub 仓库
- 熟悉[二维码服务教程](/tutorials/01-creating-a-service/)会有所帮助，但并非必需

@end

<br/>

@steps
### 从全新的 Wails 应用开始
使用 vanilla 模板搭建新项目：

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

现在应该已经生成了一个包含`main.go`、`frontend/`和`Taskfile.yml`的目录。确认它能够构建并启动：

```bash
wails3 task dev
```

此时应打开一个空白的 Wails 窗口。退出应用，然后继续。

### 添加更新程序导入
打开`main.go`，将两个更新程序包添加到导入列表中：

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

它们会引入 Updater 本身和 GitHub Releases 提供程序。

### 配置 Updater
每个`*application.App`都已接入`app.Updater`，你只需调用`Init`：

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

将此代码放在`application.New`之后、`app.Run()`之前。

@note{type="note" title="版本字符串格式"}
传入与发布标签相同的版本，但**不要**包含开头的`v`。提供程序会在其一侧移除标签名称中的`v`。此处的`1.0.0` ↔ GitHub 上的`v1.0.0`。

@end

### 添加触发更新的菜单项
在同一个`main.go`中添加“检查更新…”菜单项：

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall`会打开框架的更新窗口，运行`Check`；如果找到发布版本，则自动运行`DownloadAndInstall`。没有新版本时，窗口会保持打开并显示“已是最新”状态，用户可通过**关闭**按钮将其关闭。

@note{type="caution" title="在 goroutine 中运行"}
`CheckAndInstall`会一直阻塞到验证和安装完成。如果从菜单点击事件中直接调用它，将阻塞 UI 线程。请用`go func()`将其包装起来。

@end

### 在没有发布版本时运行一次
```bash
wails3 task dev
```

点击**应用 → 检查更新…**。你应看到更新窗口短暂打开、访问 GitHub API、未找到比`1.0.0`更新的发布版本，最后停留在带有绿色 ✓ 的**已是最新**状态。

如果此处出现错误，通常是以下原因之一：

| 症状 | 修复方法 |
| --- | --- |
| `404 Not Found` | `Repository`字段不正确——必须为`owner/repo` |
| `403 rate-limited` | 将`Token: "ghp_…"`添加到 github.Config（使用具有`public_repo`作用域的 PAT） |
| 网络错误 | 确认正在运行的应用可以访问`api.github.com` |

### 发布测试版本
将`main.go`中的`currentVersion`递增到`1.0.0`（也可以保持不变）。为一个平台执行构建，以获得可附加到发布版本的二进制文件：

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

在二进制文件旁生成一个`SHA256SUMS`文件：

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

你应该会看到一行或多行类似以下内容的输出：

```
abc123…  updater-tutorial-darwin-arm64.zip
```

现在将其作为**v2.0.0**发布到你的 GitHub 仓库：

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="资产命名"}
默认资产匹配器通过文件名中的`GOOS`和`GOARCH`子字符串进行选择。只要资产名称包含`darwin`（或`linux`/`windows`）以及`arm64`（或`amd64`/`386`），匹配器就能找到它。有关自定义匹配器的信息，请参阅[更新程序指南](/guides/updater/#github-releases--updaterprovidersgithub)。

@end

### 运行应用并验证更新
在`currentVersion`仍为`1.0.0`时，再次运行应用：

```bash
wails3 task dev
```

点击**应用 → 检查更新…**。这一次，你应该会看到类似以下内容：

![处于“更新就绪”状态的默认更新程序窗口，其中显示版本标记、由 Markdown 渲染的发布说明，以及“重启并应用”主按钮。](/assets/updater/default-window-ready.png)

- 主图标从蓝色 ↓（“有可用更新”）变为绿色 ✓（“更新就绪”）。
- 副标题显示`v1.0.0 → v2.0.0 · <size>`。
- 发布说明面板会渲染你的 Markdown，包括粗体、行内代码和表格。
- 下载期间进度条会逐渐填满（速度会很快，因为二进制文件很小）。

Updater 会将新二进制文件暂存到临时目录中。要完成更新：

- 点击**重启并应用**。
- 你的应用退出，辅助程序替换二进制文件，然后重新启动新的二进制文件。
- 重新启动的应用会报告`currentVersion = "1.0.0"`（因为我们对其进行了硬编码），但磁盘上的字节与v2.0.0构建一致。

在实际应用中，`currentVersion`会在构建时通过`-ldflags`设置，因此新二进制文件会知道自己当前的版本是v2.0.0，后续检查也就不会发现更新。

### 将`currentVersion`接入构建流程
将该常量替换为构建时变量：

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

然后在构建命令中：

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

或者，将`-ldflags`添加到`Taskfile.yml`，使其从`git describe --tags`获取值。

### 添加加密签名（建议用于生产环境）
SHA256SUMS 路径可验证*完整性*（这些字节与 GitHub 存储的内容一致），但不能验证*真实性*（即这些字节由你的发布流水线生成，而不是来自遭入侵的维护者账户）。为防止篡改，请使用 Ed25519 密钥对每个发布版本签名：

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

对于每个版本，请使用私钥为每项资产的 SHA-256 摘要签名。下面是一个简短的 Go 辅助程序：

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

默认的 GitHub 提供程序目前不会获取单独的签名文件——你可以[编写自定义提供程序](/guides/updater/#writing-your-own-provider)来获取签名文件，也可以改用**keygen.sh**；它会在服务器端为每个构件签名，并通过其 API 提供摘要和签名。

将公钥嵌入应用：

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

设置`PublicKey`后，任何附带`Signature`的版本都必须使用此密钥通过验证。发布源无法替换成自己的密钥——这正是在构建时通过带外方式固定密钥的意义所在。

### 自定义窗口
默认窗口可满足常见需求。如需更多控制，有三种自定义方式可供选择——请根据所需的自定义程度选择一种：

@tabs
[仅使用 CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

有关完整的变量列表，请参阅[通过 CSS 变量设置主题](/guides/updater/#theme-via-css-variables)一节。

[自定义 HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

你的 HTML 必须订阅`updater:*`事件，并通过 Wails 事件通道发出`updater:user:*`操作。有关 JS 适配代码，请参阅[替换模板](/guides/updater/#replace-the-template)。

[使用自己的窗口]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

如果你已有自己的窗口基础设施，并希望由更新程序驱动该窗口，而不是另行打开一个窗口，此方式会很有用。下面展示了一个完全自定义的 HTML 模板，它由与默认窗口相同的更新程序事件驱动：

![一个自带的更新程序窗口，采用粉橙渐变背景和自定义圆角卡片布局，展示了默认 UI 可以被完全替换。](/assets/updater/byo-custom-window.png)

@note{type="caution" title="必须启用 `AllowSimpleEventEmit`"}
更新程序的自定义 HTML 适配代码通过`wails:event:emit:`这一 postMessage 快捷方式驱动“安装 / 跳过 / 提醒 / 重启”操作；出于安全考虑，只有启用此字段后才能使用该快捷方式。如果忘记启用，按钮将不会执行任何操作，也不会显示错误。不要在加载非完全受你控制的 HTML 的窗口中启用它——有关威胁模型，请参阅指南中的[使用自己的窗口](/guides/updater/#bring-your-own-window)一节。

@end

@end

### 在后台自动检查
如果要使用定时器进行检查，以取代菜单点击或作为其补充：

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

定时器每次触发时，都会运行与手动点击相同的`CheckAndInstall`流程。如果希望定期检查在实际发现更新前保持静默，请设置`Window: updater.WindowNone`，然后自行订阅`EventUpdateAvailable`，以决定显示何种用户体验。

@end

## 大功告成

现在，你已经拥有一个具备以下功能的 Wails 应用：

- 按需及定时检查 GitHub Releases 中的更新。
- 在精美的默认窗口中将发行说明渲染为 Markdown。
- 根据你发布的 SHA-256 摘要验证下载内容。
- 可以选择使用构建时嵌入的公钥验证 Ed25519 签名。
- 就地替换正在运行的二进制文件并自动重新启动。

## 后续步骤

- [更新程序指南](/guides/updater/)包含完整的 API 参考、所有事件、所有配置选项，以及辅助程序模式的替换机制。
- 请查看[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/updater)，其中提供了可供克隆的完整工作示例。
- 测试目标仓库[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)展示了推荐的发布资产布局。

## 生产环境中的注意事项

- **macOS 上的代码签名**——Gatekeeper 要求替换后的二进制文件经过签名和公证。请在将`.app`包压缩为发布文件*之前*完成签名。更新程序会逐字节保留原始内容，不会重新签名任何文件。
- **Windows 上的防病毒软件**——从互联网下载的未签名`.exe`文件可能触发 SmartScreen 警告。请使用 Authenticode 证书为二进制文件签名，否则必须接受这样一种情况：使用管控严格的计算机的用户可能需要将你的应用加入白名单。
- **原子化发布**——请同时发布`SHA256SUMS`和二进制文件，不要分开提交。更新程序会分别获取辅助文件和二进制文件；如果两者不一致，摘要检查将以安全方式失败。
- **已跳过的版本**——默认窗口中的“跳过此版本”按钮会在本地记录该选择。如果发布关键安全更新，请使用新的版本号，以免曾忽略早期版本的用户自动跳过该更新。
