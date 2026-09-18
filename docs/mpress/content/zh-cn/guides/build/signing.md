---
title: "代码签名"
description: "在所有平台上为 Wails 应用程序签名的指南"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## 为应用程序进行代码签名

本指南介绍如何为 macOS、Windows 和 Linux 上的 Wails 应用程序签名。Wails v3 提供内置 CLI 工具，用于代码签名、公证和 PGP 密钥管理。

- **macOS** — 为 macOS 应用程序签名并进行公证
- **Windows** — 为 Windows 可执行文件和软件包签名
- **Linux** — 使用 PGP 密钥为 DEB 和 RPM 软件包签名

## 跨平台签名矩阵

此矩阵显示了可以从各源平台签名的目标格式：

| 目标格式 | 从 Windows | 从 macOS | 从 Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS .app 包 | ❌ | ✅ | ❌ |
| macOS 公证 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
可以从<strong>任何平台</strong>为 Windows 和 Linux 软件包签名。由于 Apple 工具的要求，macOS 签名必须在 Mac 上进行。

@end

### 签名后端

Wails 会自动选择最合适的可用签名后端：

| 平台 | 原生后端 | 跨平台后端 |
| --- | --- | --- |
| Windows | `signtool.exe`（Windows SDK） | 内置 |
| macOS | `codesign`（Xcode） | 不可用 |
| Linux | 不适用 | 内置 |

在原生平台上运行时，Wails 使用原生工具以获得最佳兼容性。进行交叉编译时，它会使用内置的签名支持。

## 快速开始

配置签名最快捷的方法是使用设置向导：

```bash
wails3 setup
```

这会将共享签名配置写入`~/.config/wails/defaults.yaml`，并且 **在每个平台执行签名时都会采用该配置**（请参阅[配置优先级](#heading-4)）。 其签名步骤会：

- 可从<strong>任何主机</strong>检测当前主机上已安装哪些可用于目标平台的签名工具，并显示适用于你的操作系统的安装命令（例如`brew install gnupg`、`sudo apt install osslsigncode`、`winget install GnuPG.Gpg4win`）。
- 在 macOS 上，列出钥匙串中的 Developer ID 证书。
- 对于 Linux，列出你的 GPG 密钥，并可<strong>生成和导出</strong>新密钥。
- 对于 Windows，可通过 OpenSSL<strong>生成自签名证书</strong>（用于测试）。
- **将密码安全地存储在系统钥匙串中**（而不是 Taskfile 中）。

@note{type="tip"}
密码存储在系统的原生凭据存储区中（macOS 钥匙串、Windows 凭据管理器或 Linux Secret Service）。这意味着你的签名配置是安全的，并且可用于所有 Wails 项目。

@end

## 配置优先级

运行签名任务时，每个签名选项按以下顺序解析（首个 匹配项优先）：

1. 显式传递给`wails3 tool sign`的标志（例如`--pgp-key`、`--certificate`、`--identity`）。
2. 匹配的<strong>项目 Taskfile 变量</strong>（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`等）。
3. `~/.config/wails/defaults.yaml`中的<strong>全局配置</strong>（由`wails3 setup`写入）。

这意味着 Taskfile 变量是<strong>可选的覆盖项</strong>：如果未设置某个变量，则使用 全局配置的密钥、证书或身份。如果这三者均未提供 值，签名命令会报告明确的错误，说明如何进行配置。

## 按项目配置

若要仅为单个项目而非全局配置签名，请在项目内运行按项目配置 向导；它会将`vars`写入该项目的 `build/<platform>/Taskfile.yml`文件：

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### 手动配置

你也可以手动编辑各平台专用的 Taskfile。编辑每个文件顶部的`vars`部分：

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

密码从系统钥匙串中检索（运行`wails3 setup signing`进行配置）。

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

密码从系统钥匙串中获取（运行`wails3 setup signing`进行配置）。

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

要以交互方式设置所有平台的签名配置，请使用向导：

```bash
wails3 setup signing
```

## macOS 代码签名

### 前提条件

- Apple Developer 账户（每年 $99）
- Developer ID Application 证书
- 已安装 Xcode Command Line Tools

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
要在 App Store 之外分发应用，需要<strong>Developer ID Application</strong>证书。

@end

### 配置

编辑`build/darwin/Taskfile.yml`并设置签名变量：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 变量 | 必需 | 说明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | 是 | 您的 Developer ID（例如，“Developer ID Application: Your Company (TEAMID)”） |
| `KEYCHAIN_PROFILE` | 公证时需要 | 存有凭据的钥匙串配置文件名称 |
| `ENTITLEMENTS` | 否 | 授权文件的路径 |

然后运行：

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### 授权

授权用于控制应用可以使用哪些功能。Wails 应用在开发环境和生产环境中通常需要不同的授权：

- **开发**：需要 JIT、未签名内存和调试授权
- **生产**：最少授权（仅网络访问）

使用交互式设置向导生成这两个文件：

```bash
wails3 setup entitlements
```

这会创建：

- `build/darwin/entitlements.dev.plist` - 用于开发构建
- `build/darwin/entitlements.plist` - 用于生产/签名构建

**可用预设：**\

| 预设 | 说明 |
| --- | --- |
| 开发 | JIT、未签名内存、调试、网络 |
| 生产 | 仅网络（最少授权，最安全） |
| 两者 | 同时创建开发文件和生产文件（推荐） |
| App Store | 启用沙盒，并允许网络和文件访问 |
| 自定义 | 逐项选择授权 |

@note{type="note"}
darwin Taskfile 中的`run`任务会自动使用`entitlements.dev.plist`。`sign`任务会对生产构建使用`entitlements.plist`。

@end

然后在 Taskfile 的 vars 中设置`ENTITLEMENTS`，使其指向相应的文件。

### 公证

Apple 要求所有分发的应用都必须经过公证。

@steps
### **将凭据存储在钥匙串中**（仅需设置一次）。可以运行 `wails3 setup signing`（它会提示输入这些值，并在内部调用 `notarytool`），也可以直接调用 `notarytool`：
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

### **验证公证结果**：
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公证通常需要1-2分钟。票据会自动装订到应用中。

@end

## Windows 代码签名

### 前提条件

- 代码签名证书（由 DigiCert、Sectigo 等机构颁发）
- 要在 Windows 上进行原生签名：需安装 Windows SDK（用于`signtool.exe`）
- 要从 macOS/Linux 进行跨平台签名：[`osslsigncode`](https://github.com/mtrojnar/osslsigncode)（`wails3 setup`中的签名步骤会显示适用于当前主机的安装命令）

### 生成自签名证书（测试）

如果只需测试签名流程，请运行`wails3 setup`，打开签名步骤中的<strong>Windows</strong>选项卡，然后选择<strong>生成自签名证书</strong>。这会使用 OpenSSL 创建一个用于代码签名的`.pfx`，并将其路径记录到全局配置中。

@note{type="caution"}
自签名证书<strong>仅适用于测试和内部分发</strong>，它们会触发面向最终用户的 SmartScreen 警告。公开发布需要使用由受信任 CA 签发的证书。

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

| 变量 | 必需 | 说明 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 覆盖 | .pfx/.p12 证书文件的路径（未设置时回退到全局配置） |
| `SIGN_THUMBPRINT` | 覆盖 | Windows 证书存储区中的证书指纹（可替代`SIGN_CERTIFICATE`） |
| `TIMESTAMP_SERVER` | 否 | 时间戳服务器 URL（默认：http://timestamp.digicert.com） |

@note{type="note"}
这些变量是可选的<strong>覆盖项</strong>。如果未设置，则使用通过`wails3 setup`配置的全局证书。请参阅[配置优先级](#heading-4)。

@end

@note{type="note"}
证书密码存储在系统钥匙串中，而不是 Taskfile 中。运行`wails3 setup signing`进行配置，或在 CI 中设置`WAILS_WINDOWS_CERT_PASSWORD`环境变量。

@end

然后运行：

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### 跨平台签名

可以在任意平台上为 Windows 可执行文件签名。同一套 Taskfile 配置和命令可在 macOS 和 Linux 上使用。

### 支持的 Windows 格式

| 格式 | 扩展名 | 备注 |
| --- | --- | --- |
| 可执行文件 | .exe | 标准 PE 签名 |
| 安装程序 | .msi | Windows Installer 软件包 |
| 应用包 | .msix, .appx | 现代 Windows 应用 |

## Linux 软件包签名

Linux 软件包（DEB 和 RPM）使用 PGP/GPG 密钥签名。与 Windows 和 macOS 代码签名不同，Linux 软件包签名证明软件包来自受信任的来源，而不是证明操作系统信任其中的代码。

### 前提条件

- PGP 密钥对（可使用 Wails 生成）

### 生成 PGP 密钥

最简单的方法是使用设置向导：运行`wails3 setup`，打开签名步骤中的<strong>Linux</strong>选项卡，然后选择<strong>创建新的 GPG 密钥</strong>。该向导会：

- 在 GPG 密钥环中生成一个 RSA 4096密钥（如需适用于无人值守 CI 的密钥，请将密码短语留空），
- 将其<strong>导出到`~/.wails/signing/<keyid>.asc`</strong>（构建进行签名时使用的文件），并且
- 将密钥 ID 和导出路径都记录在`~/.config/wails/defaults.yaml`中，以便在签名时自动使用。

@note{type="note"}
如果某个密钥配置时<strong>仅指定了 ID</strong>（例如，该密钥创建于向导开始自动导出密钥之前），那么下次访问 Linux 选项卡时，该密钥会导出到文件并自动填写路径，无需手动操作。构建使用密钥<em>文件</em>签名，因此路径才是关键。

@end

也可以使用`gpg`手动完成：

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

推荐设置：RSA 4096 位、1年后过期，并使用强密码保护。

@note{type="caution"}
请妥善保护私钥！加密存储私钥，并安全备份。

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
| `PGP_KEY` | 覆盖 | 导出的 PGP 私钥文件路径（如未设置，则回退到通过`wails3 setup`配置的全局密钥） |
| `SIGN_ROLE` | 否 | DEB 签名角色（默认：builder） |

@note{type="note"}
PGP 密钥密码存储在系统钥匙串中，而非 Taskfile 中。运行`wails3 setup signing`进行配置，或在 CI 中设置`WAILS_PGP_PASSWORD`环境变量。

@end

然后运行：

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 签名角色

对于 DEB 软件包，可以通过`SIGN_ROLE`指定签名角色：

- `origin`：来自软件包来源方的签名
- `maint`：来自软件包维护者的签名
- `archive`：来自归档维护者的签名
- `builder`：来自软件包构建者的签名（默认）

### 跨平台签名

可以在任意平台上为 Linux 软件包签名。同一套 Taskfile 配置和命令也适用于 Windows 和 macOS。

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

### 验证 Linux 软件包

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### 分发公钥

用户需要使用您的公钥来验证软件包：

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

在 CI 环境中，密码通过环境变量提供，而不使用系统钥匙串：

| 环境变量 | 说明 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 证书密码 |
| `WAILS_PGP_PASSWORD` | Linux 软件包的 PGP 密钥密码 |

也可以直接传入 Taskfile 变量：

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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

使用单个 Linux 运行器为 Windows 和 Linux 软件包签名：

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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
使用 Linux 运行器进行跨平台签名，无需单独使用 Windows 运行器，从而简化 CI/CD。由于 Apple 原生工具的要求，macOS 签名仍需使用 macOS 运行器。

@end

## CLI 参考

### wails3 setup signing

用于配置项目签名的交互式向导。

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

该向导将引导您完成以下操作：

- **macOS**：选择 Developer ID 证书并配置公证凭据（调用`xcrun notarytool store-credentials`）。
- **Windows**：选择使用证书文件或指纹，并设置密码和时间戳服务器。
- **Linux**：使用现有 PGP 密钥或生成新密钥（调用`gpg`），并配置签名角色。

### wails3 setup entitlements

用于配置 macOS 权限的交互式向导。

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**预设：**

- **开发**：创建包含 JIT、调试和网络权限的`entitlements.dev.plist`
- **生产**：创建仅含最低必要权限的`entitlements.plist`
- **两者**：创建这两个文件（推荐）
- **App Store**：为 Mac App Store 创建沙盒权限配置
- **自定义**：选择各项权限和目标文件

### wails3 sign

为当前平台或指定平台的二进制文件和软件包签名。此命令是一个包装器，会调用相应平台专用的签名任务。

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
| `--output` | 输出路径（可选，默认原地写入） |
| `--verbose` | 启用详细输出 |

**Windows/macOS 标志：**\

| 标志 | 说明 |
| --- | --- |
| `--certificate` | PKCS#12 证书（.pfx/.p12）的路径 |
| `--password` | 证书密码 |
| `--timestamp` | 时间戳服务器 URL |

**macOS 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--identity` | 签名身份（临时签名请使用“-”） |
| `--entitlements` | 授权 plist 文件的路径 |
| `--hardened-runtime` | 启用强化运行时（默认值：true） |
| `--notarize` | 提交公证 |
| `--keychain-profile` | 用于公证的钥匙串配置文件 |

**Windows 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--thumbprint` | Windows 证书存储中的证书指纹 |

**Linux 专用标志：**\

| 标志 | 说明 |
| --- | --- |
| `--pgp-key` | PGP 私钥的路径 |
| `--pgp-password` | PGP 密钥密码 |
| `--role` | DEB 签名角色（origin/maint/archive/builder） |

### 检查签名状态（原生工具）

v3 中<strong>没有</strong> `wails3 signing` 命令。要检查签名状态，请直接使用原生工具：

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

- 确保你的证书已安装到钥匙串中
- 使用`security find-identity -v -p codesigning`检查证书是否已过期
- 确保你拥有“Developer ID Application”证书（不能只有“Apple Development”证书）

**“公证失败”**

- 检查公证日志：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- 确保已启用强化运行时
- 确认你的应用不包含未签名的二进制文件

**“Codesign 失败”**

- 确保钥匙串已解锁：`security unlock-keychain`
- 检查应用包的文件权限

### Windows 问题

**“找不到证书”**

- 确认证书路径正确
- 检查证书密码
- 确保证书有效（未过期且未被吊销）

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

- 确保公钥已正确导入
- 检查软件包在签名后是否被修改

## 其他资源

### 官方文档

- [Apple 代码签名指南](https://developer.apple.com/support/code-signing/)
- [Apple 公证文档](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft 代码签名](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian 软件包签名](https://wiki.debian.org/SecureApt)
- [RPM 软件包签名](https://rpm-software-management.github.io/rpm/manual/signatures.html)
