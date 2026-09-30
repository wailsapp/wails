---
title: "更新器"
description: "Wails v3 的应用内自更新功能——支持可插拔提供程序、加密验证、原子替换，以及可自定义主题或替换的默认界面。"
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

更新器可在应用内发布软件更新，无需自行构建下载、验证和替换流水线。它基于`app.Updater`运行，接受一个或多个可插拔的`Provider`（GitHub Releases、keygen.sh、Sparkle AppCast、开放的 Wails Update Manifest 协议或你自己的实现），使用配置的公钥验证下载内容，安全替换正在运行的二进制文件，并通过标准 Wails 事件总线公开每次状态转换。

![处于“更新已就绪”状态的默认更新器窗口——包括随状态变化的图标、版本标记（v1.0.0 → v2.0.1 · 8.8 MB）、经 Markdown 渲染且包含 GFM 表格的发行说明，以及一个主要操作按钮。](/assets/updater/default-window-ready.png)

## 快速开始

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

这会打开框架的更新窗口、检查 GitHub、下载对应平台的构件、进行验证、替换二进制文件，然后等待用户重启。

## 生命周期

`app.Updater`是一个状态机，包含以下状态（`updater.State`）：

| 状态 | 何时进入 |
| --- | --- |
| `unconfigured` | 调用`Init`之前 |
| `idle` | 调用`Init`之后、进行任何检查之前 |
| `checking` | `Check`正在进行 |
| `up-to-date` | 最近一次提供程序响应表明调用方已是最新版本 |
| `available` | 已发现新版本，但尚未开始下载 |
| `downloading` | 正在以流式方式从提供程序接收字节 |
| `verifying` | 下载已完成，正在检查签名或摘要 |
| `installing` | 正在解包已验证的字节，并通过重命名将其移入暂存目录 |
| `ready` | 更新已暂存；调用`Restart`以应用更新 |
| `error` | 此前任一步骤失败 |

你可以随时通过`app.Updater.State()`读取当前状态。每次状态转换也会发出一个 Wails 事件（参见[事件](#heading-13)）。

`Restart` 会等待辅助进程执行到 `application.New`，然后才请求正在运行的应用退出。默认启动超时为 30 秒。如果应用在 `application.New` 之前需要较长时间进行初始化，请将 `Config.HelperReadyTimeout` 设置为更长的时长，例如 `time.Minute`。零值使用默认值；负时长会被拒绝。如果启动超时，`Restart` 会返回 `updater.ErrHelperNotReady`，并保持当前应用运行。

默认窗口会自动反映当前状态——例如，当`Check`返回没有可用升级时，用户会看到以下内容，并通过<strong>关闭</strong>将其关闭：

![处于“已是最新版本”状态的默认更新器窗口——绿色对勾、“你使用的是最新版本”标题，以及唯一的“关闭”按钮。](/assets/updater/default-window-up-to-date.png)

## 提供程序

任何满足此接口的实现都可以作为`Provider`：

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

源代码树中随附了四种实现。

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

默认构件匹配器根据文件名中包含的`GOOS`和`GOARCH`子字符串进行选择，并可识别常见别名（`amd64` / `x86_64` / `x64`、`arm64` / `aarch64`、`386` / `i386` / `x86` / `ia32`）。如需使用自定义命名方案：

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset`是同一发行版本中另一个构件的名称，其内容为`<sha256>  <filename>`格式的行（即`sha256sum`和`shasum -a 256`生成的格式）。提供程序会在`Check`期间获取该构件，找到与所选构件匹配的行，并填充`Release.Verification.Digest`，以便框架验证下载内容。

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

提供程序会自动将 keygen.sh 针对每个构件的 SHA-512校验和及 Ed25519ph 签名映射到框架的`Release.Verification`块，无需额外配置。

<strong>令牌格式：</strong>keygen.sh 令牌带有角色前缀（`admi-` / `prod-` / `envi-` / `user-`）。你在仪表板中看到的原始 UUID 是令牌的<em>标识符</em>，并非其密钥值；密钥仅在创建令牌时可见。详情请参阅 keygen.sh 的[身份验证文档](https://keygen.sh/docs/api/authentication/)。

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

无需更改即可接入现有的 Sparkle / WinSparkle 基础设施。它从 Feed 中读取`sparkle:shortVersionString`、`<enclosure url type length sparkle:os sparkle:edSignature>`和`sparkle:channel`。

不支持 Sparkle 1的 DSA 签名（`sparkle:dsaSignature`）；使用该签名方案的项目应轮换为 EdDSA（Sparkle 2）。

### Wails 更新清单 — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

使用开放的[Wails Update Manifest 协议](/reference/update-manifest/)：通过一个 JSON 文档描述最新发行版本及其各平台构件，并内嵌校验和与签名。同一文档既可托管在静态文件主机上（S3、GitHub Pages 或任意 CDN——每个渠道发布一份列出所有平台的清单），也可由动态更新服务器提供（提供程序会在每次检查时发送`platform`、`arch`、`version`和`channel`，因此服务器可以只返回一个构件，或根据许可证决定是否提供更新）。

借助 URL 占位符，只需一行配置即可定义静态布局：

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

每次清单请求都会发送所配置的请求头；仅当构件位于清单自身所在的主机上，并且未从`https`降级到`http`时，构件下载才会复用这些请求头。任何跨源或降级重定向都会移除`Authorization`请求头。

发布端由 CLI 负责：`wails3 updater manifest`可通过一条命令计算摘要、签名并描述发行文件，而`wails3 updater verify`会在上传前重新检查结果。参见[使用 wails3 CLI 发布](/reference/update-manifest/#publishing-with-the-wails3-cli)。

### 回退链

`Config.Providers`是有序的。更新器会依次遍历：第一个返回发行版本的提供程序胜出；第一个报告“已是最新版本”的提供程序会使该链短路（回退用于“主要提供程序不可访问”，而非“提供程序意见不一致”）。发生错误时会继续尝试下一个提供程序。

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### 编写自己的提供程序

典型实现只需三个方法、约150行代码。Updater 负责验证、原子暂存、替换和窗口；提供程序代码负责确定下一个版本并以流式方式传输其字节：

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

请参考项目内置的提供程序——每个提供程序都仅由一个 Go 文件构成。

## 密码学验证

框架的验证器使用`Config.PublicKey`作为信任根来验证发布版本：

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

支持的算法（`Release.Verification.SignatureAlgo`）：

| 算法 | 签名内容 | 说明 |
| --- | --- | --- |
| `ed25519` | 制品的 SHA-256摘要 | 供 Sparkle EdDSA 使用 |
| `ed25519ph` | 通过 Ed25519ph 的预哈希对完整制品签名（内部使用 SHA-512） | 供 keygen.sh 使用 |
| `ecdsa-p256` | 制品的 SHA-256摘要 | 同时接受原始`r∥s`签名和 DER 签名 |

如果发布版本附带哈希但没有签名，还支持仅摘要验证（`DigestAlgo`：`sha256` / `sha512`）。

`Config.PublicKey`是签名验证的唯一信任锚——发布源无法替换为自己的密钥。若发布版本带有`Signature`，但未配置`Config.PublicKey`，验证将以安全方式失败。验证器会在下载期间以流式方式计算摘要，因此即使更新包大小达到数 GB，验证也不会增加额外的磁盘遍历。

@note{type="caution" title="仅摘要验证 ≠ 密码学验证"}
仅包含`Digest`的发布版本，其真实性依赖注册表的 TLS 以及注册表自身提供的完整性保证，而不是由您控制的加密信任根。仅摘要验证适合检测数据衰变；如需抵御发布流水线遭入侵后的篡改，请使用签名。

@end

### 生成签名密钥

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

私钥采用 PKCS#8 PEM 格式，公钥采用 PKIX PEM 格式；`Config.PublicKey`可直接接受`.pub`文件（也接受原始32字节密钥或其 base64 编码，`genkey`会输出后者以便内联）。请使用`wails3 updater manifest -key updater.key ...`或`wails3 updater sign`为发布版本签名；请参阅[使用 wails3 CLI 发布](/reference/update-manifest/#publishing-with-the-wails3-cli)。

也可以使用 Go：

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## 制品格式

提供程序以流式方式传输您所发布文件的字节；框架随后会在替换前将其解包：

- **单个二进制文件**（例如`myapp-linux-amd64`）——直接使用，无需处理。在 Linux 上很常见。
- **`.zip`**——就地解压。归档必须恰好包含一个顶层条目（通常是 macOS `.app`包或单个二进制文件）。这是 macOS 的推荐打包格式。
- **`.tar.gz`** / **`.tgz`**——按照同样的单一顶层条目规则就地解压。适用于随二进制文件一同提供运行时目录树的 Linux 发行版。

包含多个顶层条目的归档会被拒绝：框架只替换一个磁盘目标，因此当归档包含多个对象时，“将此归档替换到目标位置”的含义并不明确。v1 不支持`.dmg`和`.pkg`（macOS）以及`.msi`（Windows）；请改为分发该包的`.zip`。解压过程会防止 zip-slip 攻击、拒绝指向归档根目录之外的符号链接，并限制解压后的总大小（2 GiB）和条目数（50 000）。

## 默认窗口

`app.Updater.CheckAndInstall(ctx)`会打开一个由框架管理、尺寸为520×540的窗口，其中包含：

- 随状态变化的主图标（有可用更新/正在下载时为蓝色 ↓，准备就绪/已是最新版本时为绿色 ✓，发生错误时为红色 !）
- 版本标签：`v1.0.0 → v2.0.1 · 8.8 MB`
- 可滚动的发布说明面板，其中显示<strong>渲染后的 Markdown</strong>（段落、粗体/斜体、列表、GFM 表格、行内代码、围栏代码块、h1–h3、链接）
- 每种状态对应一个主要操作（安装 / 重启并应用 / 重试）
- 幽灵样式的次要操作（跳过此版本 / 稍后提醒我）
- 通过`prefers-color-scheme`支持深色/浅色模式
- 总大小未知时显示不确定进度的流光效果

该窗口监听 Wails 事件总线上的`updater:*`事件，并将`updater:user:*`操作发送回 Go。

### 通过 CSS 变量设置主题

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

默认样式表公开以下变量，您可以覆盖其中任意变量：

| 变量 | 默认值（浅色） | 默认值（深色） |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | 系统字体栈 | — |

### 替换模板

提供你自己的 HTML；它只需监听`wails:updater:*`事件并发出`wails:updater:user:*`操作：

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML 窗口在加载时没有资产服务器源，因此无法动态获取`/wails/runtime.js`。在这种窗口中与宿主通信有两种方式：

1. <strong>直接编写 HTML。</strong>对于任何在打开时设置了`WebviewWindowOptions.AllowSimpleEventEmit = true`和`HTML`的窗口，框架都会自动注入一个精简的`window.wails.Events`垫片——更新程序的内置路径和 BYO 路径正是如此。无需构建步骤。下面的示例采用此方式。
2. **使用你偏好的打包工具**（Vite、esbuild、Rollup）打包`@wailsio/runtime`，并在构建时将其导入自定义 HTML。`Events.On`完全在客户端运行，因此可以开箱即用；`Events.Emit`则通过运行时的 fetch 传输层，而 null 源会使其失效——因此，请通过运行时的[`setTransport`](https://wails.io/wails/runtime.js)钩子安装一个精简的 postMessage 传输层，经由`window._wails.invoke("wails:event:emit:<name>")`进行路由。如果`window.wails.Events`已在作用域中，框架注入不会执行任何操作，因此两种方式不会冲突。

无论采用哪种方式，自定义 HTML 中编写的 JS 都会调用相同的`Events.On`/`Events.Emit` API：

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

该垫片公开了使用裸名称事件所需的现代运行时功能子集：`Events.On(name, cb)`返回取消订阅函数，`Events.Emit(nameOrEventObject)`则通过受控的`wails:event:emit:` postMessage 路径路由到宿主。它会在页面加载时安装一次，并且先于你的任何内联脚本运行。

如果你<em>希望</em>覆盖该垫片（或者正通过其他方式加载完整运行时），请在页面的第一个`<script>`标签执行前设置`window.wails.Events`，注入便会自行跳过。

### 窗口外观

无需修改 HTML 即可覆盖窗口选项（尺寸、无边框、始终置顶）：

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### 使用你自己的窗口

使用你自行创建的`*application.WebviewWindow`来驱动更新流程。更新程序会在你的窗口上调用`Show()`/`Close()`/`EmitEvent()`——由你的 HTML 决定渲染内容：

![一个“使用你自己的窗口”的更新程序窗口，采用粉橙渐变背景、单张白色圆角卡片和自定义字体排印，并由相同的更新程序事件驱动可见状态。它展示了默认 UI 可以被多么彻底地替换。](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

你的 HTML 与内置模板一样使用`window.wails.Events.On`/`Events.Emit`——只要窗口设置了`AllowSimpleEventEmit: true`，无论它归框架还是你所有，框架自动注入的垫片都会随之提供。有关该 API，请参阅[替换模板](#heading-9)。

@note{type="caution" title="BYO 更新程序窗口必须设置`AllowSimpleEventEmit`"}
出于安全考虑，框架通过此字段控制`wails:event:emit:` postMessage 快捷路径：未设置该字段的窗口无法伪造宿主端自定义事件。更新程序的自定义 HTML 垫片会通过此快捷路径发出`updater:user:*`事件，因此 BYO 窗口如果忘记设置该字段，就会悄无声息地丢弃每次按钮点击——用户点击“安装”后不会有任何反应。

对于任何会加载你无法完全控制的 HTML（远程 URL、用户提供的内容）的窗口，请保持`AllowSimpleEventEmit`处于<strong>关闭</strong>状态。启用后，页面中的任何 JavaScript（包括 XSS 注入点）都能触发任意`app.Event.On(name, …)`处理程序。此快捷路径只能传递裸名称（不含载荷），也无法访问绑定/Call 路径；但如果 Go 代码中的特权自定义事件处理程序仅依据事件名称执行操作，它仍可触发这些处理程序。

框架的<em>内置</em>更新程序窗口已在内部设置此字段——只有 BYO 调用方需要记得设置。

@end

### 无界面模式

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

整个过程中都不会打开窗口。请从你自己的 UI（或现有主窗口）订阅`updater:*`事件，并从按钮处理程序调用`app.Updater.CheckAndInstall(ctx)`。这适用于仅在发现更新时才应显示提示的定期后台检查，也适用于将更新流程集成到自定义设置面板中的应用。

## 事件

Go 和 JavaScript 都通过标准 Wails 事件总线订阅。**不要手动输入传输字符串**——请使用 updater 包（Go）或 runtime 包（JS）导出的常量。两层使用同一组名称，并通过回归测试保持同步。

### 从 Go 使用

这些常量位于`github.com/wailsapp/wails/v3/pkg/updater`中。通过`app.Event.On(name, fn)`订阅；回调会收到一个`*application.CustomEvent`，其`Data`字段是[事件参考](#heading-14)中列出的强类型载荷——请使用类型断言，不要进行 JSON 解码：

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

所有可用的 Go 常量：

| 常量 | 传输字符串 |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### 从 JavaScript

这些常量位于`@wailsio/runtime`的`Updater.Events`下。它们与 Go 中的名称相同，并按子命名空间（`User.*`、`Window.*`）组织，便于通过自动补全发现：

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

自定义 HTML 向宿主<em>回传</em>的用户操作事件位于`Updater.Events.User`下：

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### 事件参考

订阅端（宿主 → 页面）：

| 常量（Go） | 常量（JS） | 载荷 | 触发时机 |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | 无 | 每次`Check`往返之前 |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check`发现了较新的版本 |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | 无 | `Check`确认当前已是最新版本 |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | 开始传输字节数据 |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | 下载期间约为10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | 所有字节均已写入，尚未验证 |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | 开始检查签名/摘要 |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | 开始解包和暂存 |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | 等待重启 |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | 任一阶段失败 |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | 每个会话一次，在重放快照之前 |

页面端（页面 → 宿主）——如果编写自定义模板，你的代码需订阅这些事件：

| 常量（Go） | 常量（JS） | 触发时机 |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | 窗口加载完成；宿主端会重放当前状态 |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | 处于`available`状态时的主要操作 |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | 处于`ready`状态时的主要操作 |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | “跳过此版本” |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | “稍后提醒我” |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | 关闭按钮 |

## API 参考

### `updater.Config`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `CurrentVersion` | `string` | <strong>必填。</strong>与发布版本所用标签相同的字符串（不带`v`前缀） |
| `Providers` | `[]updater.Provider` | <strong>必填。</strong>按顺序尝试的回退链 |
| `PublicKey` | `[]byte` | PEM 或原始字节。可选，但签名发布版本在未提供该值时会以安全方式失败 |
| `CheckInterval` | `time.Duration` | 非零值会启动后台轮询循环，并调用`CheckAndInstall` |
| `Platform` | `string` | 覆盖`runtime.GOOS`以选择构件 |
| `Arch` | `string` | 覆盖`runtime.GOARCH`以选择构件 |
| `Channel` | `string` | 目前仅供参考；通道筛选取决于具体提供程序 |
| `Window` | `updater.WindowOption` | `nil`（内置默认值）、`&BuiltinWindow{…}`、`BYOWindow(handle)`或`WindowNone` |

### `*updater.Updater`的方法

| 签名 | 用途 |
| --- | --- |
| `Init(cfg Config) error` | 进行配置。第二次调用时返回`ErrAlreadyConfigured` |
| `State() State` | 当前生命周期阶段 |
| `CurrentVersion() string` | 传给`Init`的版本 |
| `Check(ctx) (*Release, error)` | 依次遍历提供程序链。`(rel, nil)` = 发现更新，`(nil, nil)` = 已是最新版本，`(nil, err)` = 全部失败 |
| `DownloadAndInstall(ctx) error` | 以流式方式传输、验证、解压（如果是归档文件）并暂存。需要先调用`Check` |
| `CheckAndInstall(ctx) error` | 便捷操作：打开窗口，调用`Check`；如果发现更新，再调用`DownloadAndInstall` |
| `Restart(ctx) error` | 启动辅助程序，调用`Host.Quit`，然后退出；辅助程序会替换文件并重新启动应用 |
| `DownloadedPath() string` | 暂存更新在磁盘上的位置；如果没有暂存更新，则为`""` |
| `SkipVersion(v string)` | 将`v`记录为已跳过；后续`Check`会将其视为已是最新版本 |
| `SkippedVersion() string` | 读取当前跳过的版本 |
| `StopPeriodicCheck()` | 取消由`Config.CheckInterval`启动的计时器，并等待循环返回 |

### 错误

| 哨兵值 | 返回方 |
| --- | --- |
| `ErrAlreadyConfigured` | 首次成功后调用`Init` |
| `ErrNotConfigured` | 在`Init`之前执行的任何操作 |
| `ErrNoPendingRelease` | 未先调用`Check`便调用`DownloadAndInstall` |
| `ErrDownloadInProgress` | 已有另一个流程运行时调用`DownloadAndInstall` |
| `ErrNotReady` | 没有已暂存的更新时调用`Restart` |

## 交换过程的工作原理

`Restart`会设置哨兵环境变量并重新执行当前二进制文件。`application.New`在启动时检测这些变量，并转入辅助程序模式：

1. 辅助程序最多等待30秒，让父进程 PID 退出（在 Windows 上，`platformIsAlive`通过`syscall.OpenProcess` + `GetExitCodeProcess`轮询；在 Unix 上则通过`os.FindProcess` + `proc.Signal(syscall.Signal(0))`轮询）。
2. 辅助程序会备份目标（文件采用复制，macOS `.app`包目录采用递归复制）。
3. 辅助程序用已暂存的制品替换目标，最多重试20次，每次尝试之间退避500毫秒：
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`。指向旧 inode 的已打开文件描述符仍然有效。
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`。Windows 允许重命名映像仍处于映射状态的文件，但不允许删除它们；下次更新时，辅助程序会清理所有遗留的`.old.*`同级文件，前提是拥有它们的内核映射已释放。

4. 辅助程序会在新二进制文件上恢复原可执行文件的模式（下载的文件使用默认 umask 创建，这会在 Unix 上去除`+x`；在 Windows 上此操作不执行任何动作）。
5. 辅助程序会清除辅助程序模式的环境变量，并重新启动已完成替换的二进制文件。
6. 辅助程序退出。

如果启动失败，辅助程序会恢复备份。如果父进程未在30秒内退出，辅助程序会在接触目标前中止（因此，即使应用退出对话框阻塞了`Quit`，用户仍可保留可用的应用）。

对于以`.zip`形式分发的 macOS `.app`包（推荐的打包方式），系统会在验证与就绪之间解压归档，以便辅助程序获得一个真实目录并将其交换到目标位置。

## 定期检查

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

当`CheckInterval > 0`时，后台 goroutine 会按配置的时间间隔调用`CheckAndInstall`。如果定时信号到达时已有另一个流程正在进行（检查、下载、验证或安装），该信号会被丢弃，因为不支持并发状态机。

若要在后台静默轮询且仅在发现内容时显示提示，请设置`Window: updater.WindowNone`，并在你自己的 UI 中响应`EventUpdateAvailable`。

## 跳过与稍后提醒

默认窗口中的“跳过此版本”按钮会通过`SkipVersion(rel.Version)`记录可用版本。后续`Check`会发现同一版本并将其视为已是最新版本（直到用户更新`CurrentVersion`；成功执行`Restart`后会自动更新）。“稍后提醒我”只会关闭窗口，不记录任何内容。

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## 分发检查清单

发布将由更新程序安装的版本前：

1. <strong>选择正确的归档格式。</strong>macOS：包含`.app`包的`.zip`。Linux：单个二进制文件或`.tar.gz`。Windows：单个`.exe`或`.zip`。不支持`.dmg`、`.msi`或`.pkg`。
2. 使用与`Config.PublicKey`匹配的私钥<strong>签名制品</strong>。对于由提供程序发布的源（keygen.sh、AppCast），请遵循各提供程序的签名流程。对于使用`ChecksumAsset`的 GitHub Releases，请使用`sha256sum` / `shasum -a 256`生成`SHA256SUMS`文件。
3. **确保版本字符串匹配。**`Config.CurrentVersion`必须与发布版本的版本标签完全匹配（例如`1.0.0` ↔ 标签`v1.0.0`；提供程序端会移除开头的`v`）。
4. 发布前，至少在目标平台上<strong>测试一次交换过程</strong>——代码签名、公证和 Gatekeeper 处理因平台而异，更新程序本身不涵盖这些操作。

## 故障排除

**“签名需要公钥，但未配置公钥”** — 发布版本包含`Signature`字段，但`Config.PublicKey`为空。请设置公钥，或修改发布流水线，使其不包含签名。

**“摘要不匹配”** — 下载的字节与提供程序所承诺的内容不匹配。通常是下载不完整（网络短暂故障）或制品损坏。重新运行通常可以解决。

**窗口打开后立即消失，既无 Markdown，也无进度** — 你的自定义 HTML 没有调用`wails:runtime:ready`。请参阅[替换模板](#heading-9)中的垫片代码。

**Windows 更新始终无法完成；辅助程序日志显示“remove old (attempt N): Access is denied”** — 这只会发生在此 PR 的`de764fb`之前版本中；当前实现采用先重命名到一旁的方式，不会遇到此问题。请升级。

**macOS Gatekeeper 阻止已交换的二进制文件** — 必须在端到端流程中保留代码签名。请为原始`.app`签名；如果你的构建流水线在更新时修改权利配置，还要<em>并且</em>重新签名将要重启的二进制文件。

## 另请参阅

- 可运行示例：[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)
- 测试演示仓库：[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- 教程：[为 Wails 应用添加自更新](/tutorials/04-self-update-a-wails-app/)
