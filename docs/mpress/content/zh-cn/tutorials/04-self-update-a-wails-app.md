---
title: "可自更新的 Wails 应用"
description: "构建一个可通过 GitHub Releases 自行更新的 Wails v3应用——涵盖从`wails3 init`到签名版本验证和辅助模式替换的全过程。"
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

本教程将为一个全新的 Wails v3应用添加应用内更新程序。完成后，该应用将能够：

- 按需检查 GitHub Releases（也可选择定时检查）。
- 下载适用于当前操作系统和架构的资源。
- 根据下载的字节验证 SHA-256摘要（也可选择验证 Ed25519 签名）。
- 在框架的默认更新窗口中显示发行说明。
- 替换正在运行的二进制文件并重新启动，而且无需分发单独的辅助可执行文件。

我们将使用<strong>GitHub Releases</strong>作为更新源，因为它免费且不需要任何基础设施。同样的模式也适用于[keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen)和[Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast)——完成本教程后，请参阅[更新程序指南](/guides/updater/)。

@note{type="tip" title="前提条件"}
- Go 1.25或更高版本
- 已安装`wails3` CLI（`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`）
- 一个可向其推送发行版本的 GitHub 仓库
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

现在应当已有一个包含`main.go`、`frontend/`和`Taskfile.yml`的目录。确认项目能够构建并启动：

```bash
wails3 task dev
```

此时应打开一个空白的 Wails 窗口。退出应用，然后继续。

### 添加更新程序导入项
打开`main.go`，将两个更新程序包添加到导入项中：

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

这两个包分别引入 Updater 本身和 GitHub Releases 提供程序。

### 配置 Updater
`app.Updater`已接入每个`*application.App`，你只需调用`Init`：

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
传入与发行版本标签相同的版本号，但<strong>不要</strong>包含开头的`v`。提供程序会在其一侧移除标签名称中的`v`。此处的`1.0.0` ↔ GitHub 上的`v1.0.0`。

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

`CheckAndInstall`会打开框架的更新窗口，运行`Check`；如果发现发行版本，则自动运行`DownloadAndInstall`。没有新版本时，窗口会保持打开并显示“已是最新版本”状态，用户可通过<strong>关闭</strong>按钮将其关闭。

@note{type="caution" title="在 goroutine 中运行"}
`CheckAndInstall`会阻塞，直到验证和安装完成。如果直接在菜单点击处理程序中调用它，将阻塞 UI 线程。请用`go func()`将其包装起来。

@end

### 在没有发行版本时运行一次
```bash
wails3 task dev
```

点击<strong>应用 → 检查更新…</strong>。此时应看到更新窗口短暂打开、访问 GitHub API、未找到比`1.0.0`更新的发行版本，最后停留在带有绿色 ✓ 的<strong>已是最新版本</strong>状态。

如果此处出现错误，通常是以下原因之一：

| 症状 | 修复方法 |
| --- | --- |
| `404 Not Found` | `Repository`字段不正确——必须为`owner/repo` |
| `403 rate-limited` | 将`Token: "ghp_…"`添加到 github.Config（使用具有`public_repo`作用域的 PAT） |
| 网络错误 | 确认正在运行的应用可以访问`api.github.com` |

### 发布测试版本
将`main.go`中的`currentVersion`提升至`1.0.0`（也可保持不变）。针对一个平台进行构建，以获得可附加到发行版本的二进制文件：

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

你应当会看到一行或多行类似以下内容：

```
abc123…  updater-tutorial-darwin-arm64.zip
```

现在将其作为<strong>v2.0.0</strong>发布到你的 GitHub 仓库：

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

@note{type="note" title="资源命名"}
默认资源匹配器根据文件名中的`GOOS`和`GOARCH`子字符串进行选择。只要资源名称包含`darwin`（或`linux`/`windows`）以及`arm64`（或`amd64`/`386`），匹配器就能找到它。有关自定义匹配器，请参阅[更新程序指南](/guides/updater/#github-releases--updaterprovidersgithub)。

@end

### 运行应用并验证更新
保持`currentVersion`仍为`1.0.0`，再次运行应用：

```bash
wails3 task dev
```

点击<strong>应用 → 检查更新…</strong>。这一次，你应当会看到类似以下内容：

![处于“更新就绪”状态的默认更新程序窗口，其中显示版本徽标、经 Markdown 渲染的发行说明，以及“重新启动并应用”主按钮。](/assets/updater/default-window-ready.png)

- 主图标会从蓝色 ↓（“有可用更新”）变为绿色 ✓（“更新就绪”）。
- 副标题显示`v1.0.0 → v2.0.0 · <size>`。
- 发行说明面板会渲染你的 Markdown，包括粗体、行内代码和表格。
- 下载期间进度条会逐渐填满（速度会很快，因为二进制文件很小）。

Updater 会将新二进制文件暂存到临时目录中。要完成更新：

- 点击<strong>重新启动并应用</strong>。
- 应用退出，辅助程序替换二进制文件，然后重新启动新的二进制文件。
- 重新启动的应用会报告`currentVersion = "1.0.0"`（因为我们对其进行了硬编码），但磁盘上的字节与v2.0.0构建版本一致。

在实际应用中，`currentVersion`应在构建时通过`-ldflags`设置，使新二进制文件知道自己现在是v2.0.0，这样后续检查就不会发现更新。

### 将`currentVersion`接入构建流程
将常量替换为构建时变量：

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
SHA256SUMS 路径验证的是<em>完整性</em>（这些字节与 GitHub 存储的内容一致），而非<em>真实性</em>（即这些字节由你的发布流水线生成，而不是来自已遭入侵的维护者账户）。为防止篡改，请使用 Ed25519 密钥为每个发行版本签名：

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

每次发布时，都使用私钥对每个资源的 SHA-256 摘要进行签名。下面是一个简单的 Go 辅助程序：

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

默认的 GitHub 提供程序目前不会获取单独的签名文件——你可以[编写自定义提供程序](/guides/updater/#writing-your-own-provider)来获取该文件，也可以改用<strong>keygen.sh</strong>；它会在服务器端签署每个制品，并通过其 API 同时提供摘要和签名。

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

设置`PublicKey`后，任何附带`Signature`的发行版都必须通过此密钥的验证。发行源无法替换为自己的密钥——这正是在构建时通过带外方式固定密钥的意义所在。

### 自定义窗口
默认窗口适用于常见情况。如果需要更多控制，可以使用以下三种扩展方式——根据所需的自定义程度选择一种：

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

你的 HTML 必须订阅`updater:*`事件，并通过 Wails 事件通道发出`updater:user:*`操作。有关 JS 适配层，请参阅[替换模板](/guides/updater/#replace-the-template)。

[使用你自己的窗口]
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

如果你已有自己的窗口基础设施，并希望更新程序驱动该窗口而不是再打开一个窗口，此选项会很有用。下面是一个完全自定义的 HTML 模板，它由与默认模板相同的更新程序事件驱动：

![一个自带的更新程序窗口，采用粉橙渐变背景和自定义圆角卡片布局，展示默认 UI 可被完全替换。](/assets/updater/byo-custom-window.png)

@note{type="caution" title="必须启用 `AllowSimpleEventEmit`"}
更新程序的自定义 HTML 适配层通过`wails:event:emit:` postMessage 快捷方式驱动“安装 / 跳过 / 稍后提醒 / 重启”操作；出于安全考虑，该快捷方式受此字段控制。忘记启用此字段会导致按钮点击后无任何反应。不要在加载了你无法完全控制的 HTML 的窗口上启用它——有关威胁模型，请参阅本指南的[使用你自己的窗口](/guides/updater/#bring-your-own-window)一节。

@end

@end

### 在后台运行自动检查
要按定时器检查，而不是通过菜单点击检查（或在菜单点击检查之外同时进行定时检查），请使用：

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

每次定时触发都会运行与手动点击相同的`CheckAndInstall`流程。如果希望定期检查在实际发现内容之前保持静默，请设置`Window: updater.WindowNone`，然后自行订阅`EventUpdateAvailable`，以决定显示何种用户体验。

@end

## 大功告成

现在，你已经拥有一个具备以下功能的 Wails 应用：

- 按需和按定时器检查 GitHub Releases 中的更新。
- 在精美的默认窗口中将发行说明渲染为 Markdown。
- 使用你发布的 SHA-256 摘要验证下载内容。
- 可选择使用构建时嵌入的公钥验证 Ed25519 签名。
- 原地替换正在运行的二进制文件，并自动重新启动。

## 后续步骤

- [更新程序指南](/guides/updater/)包含完整的 API 参考、所有事件、所有配置选项以及辅助模式下的替换机制。
- 请查看[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)，其中提供了可供克隆的完整可运行示例。
- 测试目标仓库[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)展示了推荐的发行资源布局。

## 生产环境中的注意事项

- **macOS 上的代码签名**——Gatekeeper 要求替换后的二进制文件经过签名和公证。请在将`.app`包压缩为发行文件<em>之前</em>对其签名。更新程序会原样保留所有字节，不会重新签署任何内容。
- **Windows 上的防病毒软件**——从互联网下载的未签名`.exe`文件可能会触发 SmartScreen 警告。请使用 Authenticode 证书对二进制文件进行签名；否则，使用受严格管控计算机的用户可能需要将你的应用加入白名单。
- **原子化发布**——请一起发布`SHA256SUMS`和二进制文件，不要分成不同的提交。更新程序会分别获取伴随文件和二进制文件；如果两者不同步，摘要检查将以安全方式失败并拒绝继续。
- **跳过版本**——默认窗口中的“跳过此版本”按钮会在本地记录跳过操作。如果要发布关键安全更新，请使用新的版本号，以免曾忽略较早发行版的用户自动跳过该更新。
