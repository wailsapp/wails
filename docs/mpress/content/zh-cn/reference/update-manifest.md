---
title: "更新清单协议"
description: "Wails 应用用于发现并验证自更新的开放 JSON 协议，可由任何静态文件托管服务或动态更新服务器提供。"
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Wails 更新清单协议是 Wails 应用与更新源之间一份小型、开放的 JSON 约定。任何能够通过 HTTPS 提供 JSON 文件的服务都可以提供 Wails 更新，例如 S3 存储桶、GitHub Pages、CDN，或根据许可证控制版本发布的动态更新服务器。

客户端以 `endpoint` 提供程序（`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`）的形式随框架发布。本页是供服务器端实现者参考的线路格式规范。

## 设计目标

1. <strong>适合静态托管。</strong>每个发布通道只需一个清单文件，其中列出各个平台的制品，即可构成完整实现。无需服务器端代码。
2. <strong>适合动态服务器。</strong>客户端每次检查都会发送 `platform`、`arch`、`version` 和 `channel`，因此服务器可以仅返回一个制品、应用许可规则，或在调用方已是最新版本时返回 `204 No Content`。
3. <strong>验证优先。</strong>清单为每个制品提供校验和与签名，Wails 更新程序会使用构建时固定在应用二进制文件中的公钥进行验证。更新源绝不能自行选择信任根。

## 请求

客户端向配置的清单 URL 发出 `GET`，并附带 `Accept: application/json` 以及应用配置的所有请求头（例如 `Authorization: License <key>`）。

URL 可以包含占位符，客户端会在每次检查时替换这些占位符：

| 占位符 | 替换为 |
| --- | --- |
| `{{platform}}` | 当前运行的操作系统，以 Go `GOOS` 值表示（`darwin`、`windows`、`linux`） |
| `{{arch}}` | 当前运行的体系结构，以 Go `GOARCH` 值表示（`amd64`、`arm64`……） |
| `{{version}}` | 当前安装的版本 |
| `{{channel}}` | 已配置的发布通道（如果已设置） |

这四个值中，任何未被占位符使用的值都会作为同名查询参数追加到 URL 中（仅在已配置时追加 `channel`）。因此，以下两种配置均有效且等效：

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

静态托管服务只需忽略收到的查询参数。

## 响应

| 状态 | 含义 |
| --- | --- |
| `200 OK` | 响应中包含清单。由客户端判断其是否为升级版本。 |
| `204 No Content` | 服务器已比较版本，调用方当前已是最新版本。 |
| `404 Not Found` | 没有已发布的内容（与已是最新版本同等处理）。 |
| 其他任何状态 | 发生错误。更新程序会继续尝试下一个已配置的提供程序。 |

`200` 响应正文是一个清单文档：

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### 顶层字段

| 字段 | 类型 | 必需 | 说明 |
| --- | --- | --- | --- |
| `schemaVersion` | int | 否 | 协议版本。省略时表示 `1`。客户端会拒绝高于其所能理解版本的值。 |
| `version` | string | **是** | SemVer 2.0.0，可以带或不带前导 `v`。 |
| `channel` | string | 否 | 仅供参考。配置为其他通道的客户端会将该清单视为没有更新。 |
| `name` | string | 否 | 供用户阅读的发布标题，显示在更新窗口中。 |
| `notes` | string | 否 | 以 Markdown 编写的发行说明，呈现在更新窗口中。 |
| `publishedAt` | string | 否 | RFC 3339 时间戳。 |
| `artifacts` | array | **是** | 每个可下载制品对应一个条目。条目顺序表示发布者的偏好。 |
| `metadata` | object | 否 | 自由格式的键值数据，将原样传递给应用。 |

客户端会忽略未知字段，因此服务器可以添加自己的字段，而不会破坏任何客户端。服务器特有的附加数据应放在`metadata`中。

### 制品字段

| 字段 | 类型 | 必需 | 说明 |
| --- | --- | --- | --- |
| `url` | string | **是** | 绝对 URL，或相对于清单 URL 的 URL。仅限`http(s)`。 |
| `platform` | string | 否 | Go 的`GOOS`值。接受常用别名（`macos`、`win`等）。空值匹配所有平台。 |
| `arch` | string | 否 | Go 的`GOARCH`值。接受常用别名（`x86_64`、`aarch64`等）。空值匹配所有架构。 |
| `filename` | string | 否 | 默认为`url`的最后一个路径段。 |
| `filetype` | string | 否 | 默认为文件扩展名。 |
| `size` | int | 否 | 字节数，用于显示下载进度。 |
| `digestAlgo` / `digest` | string / base64 | 否 | `sha256`或`sha512`。 |
| `signatureAlgo` / `signature` | string / base64 | 否 | `ed25519`、`ed25519ph`或`ecdsa-p256`。只要存在`signature`，就必须提供`signatureAlgo`。有关各算法所签名的内容，请参阅[更新程序指南](/guides/updater/#cryptographic-verification)。 |

客户端会选择<strong>第一个</strong>其`platform`和`arch`与当前运行系统匹配的制品。Base64 值可以带填充，也可以不带填充。

### 版本比较

清单是否表示升级始终由客户端按照 SemVer 2.0.0优先级规则决定：清单中的`version`必须严格高于已安装版本。这样一来，静态托管自然能够正确工作（清单始终描述最新版本，而已经是最新版本的客户端不会执行任何操作），同时动态服务器仍可使用`204`来节省带宽。

## 验证与信任

校验和与签名包含在清单中，但信任根不在其中：签名会根据应用在构建时通过`updater.Config.PublicKey`固定的公钥进行验证。遭入侵或被替换的更新源无法提供自己的密钥。如果制品带有签名，但应用没有固定公钥，则验证会以安全失败方式终止；签名未声明`signatureAlgo`或无法解码时也是如此：客户端绝不会悄然回退到仅验证摘要。

仅含摘要的制品会在摘要检查通过后安装。这可以防止数据损坏，但其防篡改能力依赖 TLS 和托管方自身的完整性保障。任何安全敏感的内容都应附带签名。

使用框架的`ed25519ph`方案为制品签名只需编写几行 Go 代码：

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

实际使用中，你很少需要编写这些代码：CLI 会替你完成。

## 使用 wails3 CLI 发布

`wails3 updater`命令组覆盖整个发布流水线。发布一个版本只需三条命令：

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest`接受文件或目录作为输入（会自动跳过密钥材料、`.json`、校验和及发行说明附属文件），以流式方式用 SHA-512处理每个制品；如果提供了`-key`，则使用 Ed25519ph 对摘要签名；还会从`MyApp-2.1.0-darwin-arm64.zip`等遵循惯例的文件名推断`platform`和`arch`（可识别`macOS`、`win64`、`x86_64`和`aarch64`等常用别名；对于无法推断的内容会输出警告，随后该内容将匹配所有平台）。省略`-url-prefix`可生成相对 URL，并将清单上传到制品所在位置。

`verify`在出现任何不匹配时都会以非零状态退出，因此很适合作为构建与发布之间的 CI 门禁。对于自行组装清单的服务器，`wails3 updater sign -key updater.key <files...>`会以 JSON 格式输出每个文件的`digest`/`signature`字段，可直接合并到你自己的文档中。

## 身份验证

身份验证由服务器负责；协议只负责传递标头。客户端会在每次请求清单时重新发送已配置的标头。下载制品时，仅当制品 URL 与清单位于同一主机上，并且没有从`https`降级到`http`时，才会发送`Authorization`标头；发生任何跨源重定向或降级重定向时，该标头都会被移除，因此凭据绝不会泄露给 CDN 或对象存储，也绝不会以明文传输。

受许可证限制的示例，可自然地与托管式许可服务配合使用：

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## 客户端配置

有关完整的`endpoint.Config`参考，以及如何将该提供程序与 GitHub、keygen.sh 和 AppCast 提供程序一起纳入回退链，请参阅[更新程序指南](/guides/updater/#providers)。
