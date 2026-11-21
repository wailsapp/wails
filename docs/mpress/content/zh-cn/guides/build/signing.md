---
title: "代码签名"
description: "在所有平台上为 Wails 应用签名的指南"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## 为应用进行代码签名

本指南介绍如何为 macOS、Windows 和 Linux 上的 Wails 应用签名。Wails v3提供用于代码签名、公证和 PGP 密钥管理的内置 CLI 工具。

- **macOS** - 为 macOS 应用签名并进行公证
- **Windows** - 为 Windows 可执行文件和软件包签名
- **Linux** - 使用 PGP 密钥为 DEB 和 RPM 软件包签名

## 跨平台签名矩阵

此矩阵展示了可从各源平台签名的目标格式：

| 目标格式 | 从 Windows | 从 macOS | 从 Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX、导出的 PFX | ✅ | ✅ | ✅ |
| Windows MSI 或证书存储区/提供程序身份 | ✅ | ❌ | ❌ |
| macOS .app / Mach-O / DMG、导出的身份 | ✅ | ✅ | ✅ |
| macOS 公证、团队 API 凭据 | ✅ | ✅ | ✅ |
| macOS 钥匙串身份或公证配置文件 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
导出的 Apple 身份和团队 API 凭据支持可移植签名与公证。macOS 钥匙串身份和公证配置文件需要 Mac。下文介绍了各格式特有的限制。

@end

### 签名后端

Wails 会自动选择最佳的可用签名后端：

| 平台 | 原生后端 | 跨平台后端 |
| --- | --- | --- |
| Windows | 对于存储区/提供程序身份、MSI 或旧版 RC2 PFX，使用`signtool.exe` | 对于导出的 RSA PFX 和 PE/MSIX，使用内置后端 |
| macOS | 对于钥匙串身份，使用`codesign` | 对于导出的 PKCS#12 身份或临时签名，使用内置后端 |
| Linux | 不适用 | 内置 |

在每种主机上，凭据类型决定使用哪个 Apple 签名后端。导出的身份使用 Go；macOS 钥匙串身份使用`codesign`。票据的获取、验证和装订使用 Go。

## 快速开始

配置签名最快的方法是使用设置向导：

```bash
wails3 setup
```

该向导会将共享签名配置写入`~/.config/wails/defaults.yaml`，并且**在每个平台执行签名时都会采用该配置**（请参阅[配置优先级](#heading-4)）。其签名步骤会：

- 显示用于生成或导出证书与密钥的设置工具，例如 OpenSSL 和 GnuPG。导出的 Apple 和 Windows 身份使用内置签名器；Windows 证书存储区/提供程序身份和 MSI 软件包需要 Windows SignTool。
- 在 macOS 上，列出钥匙串中的 Developer ID 证书。
- 对于 Linux，列出 GPG 密钥，并可**生成并导出**新密钥。
- 对于 Windows，可通过 OpenSSL**生成自签名证书**（用于测试）。
- **将密码安全地存储在系统钥匙串中**（而非 Taskfile 中）。

@note{type="tip"}
密码存储在系统的原生凭据存储区中（macOS 钥匙串、Windows 凭据管理器或 Linux Secret Service）。这意味着签名配置既安全，又可用于所有 Wails 项目。

@end

## 配置优先级

运行签名任务时，每个签名选项均按以下顺序解析（首个匹配项生效）：

1. 传递给`wails3 tool sign`的显式标志（例如`--pgp-key`、`--certificate`、`--identity`）。
2. 匹配的**项目 Taskfile 变量**（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`等）。
3. `~/.config/wails/defaults.yaml`中的**全局配置**（由`wails3 setup`写入）。

这意味着 Taskfile 变量是**可选的覆盖项**：如果未设置某个变量，则使用全局配置的密钥、证书或身份。如果这三种来源均未提供值，签名命令会报告明确的错误，并说明如何进行配置。

## 单项目配置

如果只想为单个项目配置签名，而不是进行全局配置，请在项目内运行项目级向导；该向导会将`vars`写入该项目的`build/<platform>/Taskfile.yml`文件：

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### 手动配置

也可以手动编辑各平台专用的 Taskfile。编辑每个文件顶部的`vars`部分：

@tabs
[macOS]
编辑`build/darwin/Taskfile.yml`：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

然后运行：

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
编辑`build/windows/Taskfile.yml`：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

密码从系统钥匙串中获取（运行`wails3 signing`进行配置）。

然后运行：

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
编辑`build/linux/Taskfile.yml`：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

密码从系统钥匙串中获取（运行`wails3 signing`进行配置）。

然后运行：

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

也可以直接从系统中检查签名状态：

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

要以交互方式设置所有平台的签名配置，请使用该向导：

```bash
wails3 signing
```

## macOS 代码签名

### 前提条件

- Apple Developer 账户（每年 $99）
- Developer ID Application 证书
- 对于 macOS 钥匙串身份：已安装 Xcode 命令行工具
- 对于可移植签名：导出的 RSA PKCS#12 身份；无需 Apple 签名工具

### 使用 Wake 进行可移植签名

Wake 直接读取`wails.yaml`，不会执行 Taskfile。请配置导出的 Developer ID 身份和 App Store Connect 团队 API 密钥：

```yaml
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

`credential`指定保存证书密码的环境变量。请将两个私钥文件都排除在版本控制之外。在[YAML 构建指南](/experimental/yaml-builds/)所述的发布配置文件中启用签名和公证。

可移植签名器支持64位 Intel 和 Apple Silicon Mach-O 文件、通用二进制文件、包含嵌套代码和框架的规范应用程序包，以及小于4 GiB 的 UDIF DMG。它会在发布输出之前拒绝不受支持的布局。临时签名使用`identity: "-"`，无法进行公证。PKCS#12身份可使用 AES/PBKDF2 或旧版3DES 加密；使用 RC2 的导出文件必须使用受支持的密码重新导出。

可移植的已公证签名需要明确指定 SHA-2 RFC 3161时间戳服务。旧版 SHA-1时间戳签名会被拒绝。团队 API 提交支持最大5 GiB 的单次上传；不支持个人 API 密钥和分段上传。票据验证和装订使用观察到的 Apple 服务协议并固定使用 Apple Root G3，无需`xcrun stapler`。Gatekeeper 仍是判断分发信任和吊销状态的权威机制。

以下各节介绍现有的 macOS 钥匙串工作流，该工作流供旧版项目以及选择钥匙串身份的 Wake 配置使用。

### 签名身份

检查可用的签名身份：

```bash
security find-identity -v -p codesigning
```

输出：

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
若要在 App Store 之外分发，你需要**Developer ID Application**证书。

@end

### 配置

编辑`build/darwin/Taskfile.yml`并设置签名变量：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 变量 | 是否必需 | 说明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | 是 | 你的 Developer ID（例如“Developer ID Application: Your Company (TEAMID)”） |
| `KEYCHAIN_PROFILE` | 公证时必需 | 存有凭据的钥匙串配置文件名称 |
| `ENTITLEMENTS` | 否 | 权利文件的路径 |

然后运行：

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### 权利

权利控制应用可以访问哪些功能。Wails 应用在开发环境和生产环境中通常需要不同的权利：

- **开发**：需要 JIT、未签名内存和调试权利
- **生产**：最少的权利（仅网络访问）

使用交互式设置向导生成这两个文件：

```bash
wails3 entitlements
```

这会创建：

- `build/darwin/entitlements.dev.plist` - 用于开发构建
- `build/darwin/entitlements.plist` - 用于生产/已签名构建

**可用预设：**\

| 预设 | 说明 |
| --- | --- |
| 开发 | JIT、未签名内存、调试、网络 |
| 生产 | 仅网络（最少、最安全） |
| 两者 | 同时创建开发和生产文件（推荐） |
| App Store | 启用沙盒，并允许网络和文件访问 |
| 自定义 | 选择各项权利 |

@note{type="note"}
darwin Taskfile 中的`run`任务会自动使用`entitlements.dev.plist`。`sign`任务对生产构建使用`entitlements.plist`。

@end

然后在 Taskfile 变量中设置`ENTITLEMENTS`，使其指向相应文件。

### 公证

Apple 要求所有分发的应用都经过公证。

@steps
### **将凭据存入钥匙串**（一次性设置）。可以运行`wails3 signing`（它会提示输入这些值，并在内部调用`notarytool`），也可以直接调用`notarytool`：
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **在 Taskfile 中设置 KEYCHAIN_PROFILE**，使其与上面的配置文件名称一致。
### **为应用签名并进行公证**：
```bash
wails3 task darwin:sign:notarize
```

### **验证公证**：
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公证通常需要1-2分钟。票据会自动装订到你的应用中。

@end

## Windows 代码签名

### 前提条件

- 代码签名证书（来自 DigiCert、Sectigo 等）
- 对于可移植签名：使用 AES/PBKDF2 或旧版3DES 加密的已导出 RSA PFX 身份；不需要外部签名工具
- 对于证书存储区/提供程序身份、MSI 包或旧版 RC2 PFX：在`PATH`上安装带有签名工具的 Windows SDK

### 生成自签名证书（测试）

如果只需测试签名流水线，请运行`wails3 setup`，打开签名步骤的**Windows**选项卡，然后选择**生成自签名证书**。这会使用 OpenSSL 创建用于代码签名的`.pfx`，并将其路径记录到全局配置中。

@note{type="caution"}
自签名证书**仅用于测试和内部分发**——它们会触发面向最终用户的 SmartScreen 警告。公开发布需要使用受信任 CA 颁发的证书。

@end

### 配置

编辑`build/windows/Taskfile.yml`并设置签名变量：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| 变量 | 是否必需 | 说明 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 覆盖项 | .pfx/.p12 证书文件的路径（如未设置，则回退到全局配置） |
| `SIGN_THUMBPRINT` | 覆盖项 | Windows 证书存储区中的证书指纹（`SIGN_CERTIFICATE`的替代方式） |
| `TIMESTAMP_SERVER` | 否 | 时间戳服务器 URL（默认值：http://timestamp.digicert.com） |

@note{type="note"}
这些变量是可选的**覆盖项**。如未设置，则使用通过`wails3 setup`配置的全局证书。请参阅[配置优先级](#heading-4)。

@end

@note{type="note"}
证书密码存储在系统钥匙串中，而不是 Taskfile 中。运行`wails3 signing`进行配置，或在 CI 中设置`WAILS_WINDOWS_CERT_PASSWORD`环境变量。

@end

然后运行：

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### 跨平台签名

导出的 RSA PFX 身份可在 Windows、macOS 和 Linux 上使用 Go 对 EXE/DLL 文件、NSIS 安装程序可执行文件以及单个 MSIX/APPX 包进行签名。可以替换现有签名。签名过程会暂存输出，并且仅在签名和时间戳验证成功后才发布；如果失败，则保留原始文件。使用`--output`写入单独的已签名文件。

Wake 直接从`wails.yaml`读取这些设置：

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential`指定密码环境变量的名称，绝不是密码本身。为发布目标启用`sign: true`；请参阅[YAML 构建指南](/experimental/yaml-builds/)。Wake 不会调用 Taskfile。旧项目可以继续使用上面的 Taskfile 配置。

可移植配置使用 RSA/SHA-256，并验证 RFC 3161时间戳，包括时间戳签名者的受信任证书链。默认服务为 DigiCert。MSIX Publisher 必须与签名证书的主题匹配；不支持包捆绑和大小不低于4 GiB 的归档文件。PFX 文件上限为16 MiB。其他证书算法和 RC2 加密会被拒绝。若要进行可移植签名，请使用 AES 或3DES 重新导出 RC2 PFX；在 Windows 上，经过身份验证的 RC2 PFX 改用 SignTool。密码错误或 PFX 损坏绝不会触发这种回退。

证书存储、硬件令牌和提供程序支持的密钥继续使用原生 Windows SignTool 工作流，并通过`thumbprint`进行配置。MSI 签名也需要在 Windows 上使用 SignTool。配置导出的证书不会使这些工作流变为可移植工作流。

### 支持的 Windows 格式

| 格式 | 扩展名 | 说明 |
| --- | --- | --- |
| 可执行文件和库 | .exe、.dll | 可移植 PE 签名，包括 NSIS 安装程序 |
| 安装程序 | .msi | 需要 Windows SignTool |
| 应用包 | .msix、.appx | 可移植的单包签名；不支持包捆绑 |

## Linux 包签名

Linux 包（DEB 和 RPM）使用 PGP/GPG 密钥签名。与 Windows 和 macOS 代码签名不同，Linux 包签名证明包来自可信来源，而不是证明操作系统信任其中的代码。

导出的 ASCII 装甲或二进制私钥文件可在 Windows、macOS 和 Linux 上使用原生 Go 签名器。对现有 DEB 或 RPM 进行签名时，无需解包或重新压缩即可保留其包内容，也不需要安装 GnuPG、`dpkg-sig`或`rpmsign`。Arch 包使用同一密钥加载器写入分离签名。GnuPG 对生成和导出密钥以及独立验证签名仍然很有用。旧版 DEB/RPM 配置若指定的是密钥环身份而非密钥文件，则继续使用外部签名工具。

### 前提条件

- PGP 密钥对（可使用 Wails 生成）

### 生成 PGP 密钥

最简单的方法是使用设置向导：运行`wails3 setup`，打开签名步骤中的**Linux**选项卡，然后选择**创建新的 GPG 密钥**。该向导会：

- 在 GPG 密钥环中生成一个 RSA 4096密钥（若要创建适合无人值守和 CI 使用的密钥，请将密码短语留空），
- **将其导出到`~/.wails/signing/<keyid>.asc`**（构建用于签名的文件），并且
- 在`~/.config/wails/defaults.yaml`中同时记录密钥 ID 和导出路径，以便签名时自动使用。

@note{type="note"}
如果此前仅通过**ID**配置密钥（例如，该密钥是在向导能够自动导出密钥之前创建的），则下次打开 Linux 选项卡时，它会将密钥导出到文件并填写路径，无需任何手动步骤。构建使用密钥*文件*进行签名，因此真正重要的是路径。

@end

也可以使用`gpg`手动完成：

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

推荐设置：RSA 4096位、1年后到期，并使用强密码保护。

@note{type="caution"}
请妥善保护私钥！加密存储并安全备份。

@end

### 配置

编辑`build/linux/Taskfile.yml`并设置签名变量：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 变量 | 必需 | 说明 |
| --- | --- | --- |
| `PGP_KEY` | 覆盖 | 导出的 PGP 私钥文件路径（若未设置，则回退到通过`wails3 setup`全局配置的密钥） |
| `SIGN_ROLE` | 否 | DEB 签名角色（默认值：builder） |

@note{type="note"}
PGP 密钥密码存储在系统钥匙串中，而不是 Taskfile 中。运行`wails3 signing`进行配置，或在 CI 中设置`WAILS_PGP_PASSWORD`环境变量。

@end

然后运行：

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 签名角色

对于 DEB 包，可以通过`SIGN_ROLE`指定签名角色：

- `origin`：来自包来源方的签名
- `maint`：来自包维护者的签名
- `archive`：来自归档维护者的签名
- `builder`：来自包构建者的签名（默认）

### 跨平台签名

可以从任意平台对 Linux 包进行签名。同一套 Taskfile 配置和命令也适用于 Windows 和 macOS。

### 查看密钥信息

```bash
gpg --show-keys signing-key.asc
```

输出：

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### 验证 Linux 包

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### 分发公钥

用户需要使用你的公钥来验证包：

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## GitHub Actions 集成

在 CI 环境中，密码通过环境变量提供，而不是通过系统钥匙串：

| 环境变量 | 说明 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 证书密码 |
| `WAILS_PGP_PASSWORD` | Linux 包的 PGP 密钥密码 |

也可以直接传递 Taskfile 变量：

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### macOS 工作流

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Windows 工作流

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### 跨平台工作流（Linux 运行器）

使用一个 Linux 运行器对 Windows 和 Linux 包进行签名：

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
使用 Linux 运行器进行跨平台签名，无需单独的 Windows 运行器，从而简化 CI/CD。Apple 钥匙串身份和公证配置文件需要 macOS 运行器。导出的身份和团队 API 凭据支持上述可移植 Wake 工作流。

@end

## CLI 参考

### wails3 signing

用于配置项目签名的交互式向导。

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

该向导会引导你完成：

- **macOS**：选择 Developer ID 证书并配置公证凭据（调用`xcrun notarytool store-credentials`）。
- **Windows**：选择使用证书文件还是指纹，并设置密码和时间戳服务器。
- **Linux**：使用现有 PGP 密钥或生成新密钥（调用`gpg`），并配置签名角色。

### wails3 entitlements

用于配置 macOS 权限的交互式向导。

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**预设：**

- **开发**：创建包含 JIT、调试和网络权限的`entitlements.dev.plist`
- **生产**：创建仅含最少权限的`entitlements.plist`
- **两者**：创建这两个文件（推荐）
- **App Store**：为 Mac App Store 创建启用沙盒的权限配置
- **自定义**：选择各项权限和目标文件

### wails3 sign

为当前平台或指定平台的二进制文件和软件包签名。这是一个封装命令，会调用相应平台专用的签名任务。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

此命令会运行对应的`<platform>:sign`任务，该任务使用 Taskfile 中的签名配置。

### wails3 tool sign

用于直接为指定文件签名的底层命令。Taskfile 在内部使用此命令。

```bash
wails3 tool sign [flags]
```

**通用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--input` | 待签名文件的路径 |
| `--output` | 输出路径（可选，默认原地修改） |
| `--verbose` | 启用详细输出 |

**Windows/macOS 标志：**\

| 标志 | 说明 |
| --- | --- |
| `--certificate` | PKCS#12（.pfx/.p12）证书的路径 |
| `--password` | 证书密码 |
| `--timestamp` | 时间戳服务器 URL |

**macOS 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--identity` | 签名身份（使用“-”进行临时签名） |
| `--entitlements` | 权限 plist 的路径 |
| `--hardened-runtime` | 启用强化运行时（默认值：true） |
| `--notarize` | 提交公证 |
| `--keychain-profile` | 用于公证的钥匙串配置文件 |
| `--notary-key` | App Store Connect 团队 API 私钥`.p8`文件；使用 Go Notary API 客户端 |
| `--notary-key-id` | 团队 API 密钥 ID；与`--notary-key`一起使用时为必填项 |
| `--notary-issuer-id` | 团队颁发者 ID；与`--notary-key`一起使用时为必填项 |

**Windows 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--thumbprint` | Windows 证书存储区中的证书指纹 |

**Linux 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--pgp-key` | PGP 私钥的路径 |
| `--pgp-password` | PGP 密钥密码 |
| `--role` | DEB 签名角色（origin/maint/archive/builder） |

### 检查签名状态（原生工具）

`wails3 signing`用于配置签名凭据。要检查签名状态，请直接使用原生工具：

| 任务 | 命令 |
| --- | --- |
| 列出 macOS 代码签名身份 | `security find-identity -v -p codesigning` |
| 存储公证凭据 | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| 检查 PGP 密钥文件 | `gpg --show-keys <key.asc>` |
| 生成 PGP 密钥对 | `gpg --full-generate-key` |
| 导出公钥 | `gpg --armor --export <email>` |

## 故障排除

### macOS 问题

**“未找到 Developer ID 证书”**

- 确保已将证书安装到钥匙串中
- 使用`security find-identity -v -p codesigning`检查证书是否已过期
- 确保拥有“Developer ID Application”证书（而不只是“Apple Development”证书）

**“公证失败”**

- 检查公证日志：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- 确保已启用强化运行时
- 确认应用中不包含未签名的二进制文件

**“代码签名失败”**

- 确保钥匙串已解锁：`security unlock-keychain`
- 检查应用程序包的文件权限

### Windows 问题

**“未找到证书”**

- 确认凭据路径正确
- 检查证书密码
- 确保证书有效（未过期或被吊销）

**“时间戳服务器错误”**

- 尝试使用其他时间戳服务器：
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux 问题

**“PGP 密钥无效”**

- 确保密钥文件采用 ASCII 装甲格式
- 使用`gpg --show-keys <key.asc>`检查密钥是否已过期
- 确认密码正确

**“签名验证失败”**

- 确保已正确导入公钥
- 检查软件包在签名后是否被修改

## 其他资源

### 官方文档

- [Apple 代码签名指南](https://developer.apple.com/support/code-signing/)
- [Apple 公证文档](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft 代码签名](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian 软件包签名](https://wiki.debian.org/SecureApt)
- [RPM 软件包签名](https://rpm-software-management.github.io/rpm/manual/signatures.html)
