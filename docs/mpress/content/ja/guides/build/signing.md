---
title: "コード署名"
description: "あらゆるプラットフォームで Wails アプリケーションに署名するためのガイド"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## アプリケーションのコード署名

このガイドでは、macOS、Windows、Linux 向けの Wails アプリケーションに署名する方法について説明します。Wails v3 には、コード署名、公証、PGP キー管理のための CLI ツールが組み込まれています。

- **macOS** - macOS アプリケーションへの署名と公証
- **Windows** - Windows の実行ファイルとパッケージへの署名
- **Linux** - PGP キーを使用した DEB および RPM パッケージへの署名

## クロスプラットフォーム署名マトリックス

次のマトリックスは、各ソースプラットフォームから署名できる対象を示しています。

| 対象形式 | Windows から | macOS から | Linux から |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS .app バンドル | ❌ | ✅ | ❌ |
| macOS 公証 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Windows および Linux のパッケージには、**どのプラットフォームからでも**署名できます。macOS の署名には Apple のツールが必要なため、Mac が必要です。

@end

### 署名バックエンド

Wails は、利用可能な最適な署名バックエンドを自動的に選択します。

| プラットフォーム | ネイティブバックエンド | クロスプラットフォームバックエンド |
| --- | --- | --- |
| Windows | `signtool.exe`（Windows SDK） | 組み込み |
| macOS | `codesign`（Xcode） | 利用不可 |
| Linux | 該当なし | 組み込み |

ネイティブプラットフォーム上で実行する場合、Wails は互換性を最大限に確保するため、ネイティブツールを使用します。クロスコンパイルする場合は、組み込みの署名機能を使用します。

## クイックスタート

署名を最もすばやく構成するには、セットアップウィザードを使用します。

```bash
wails3 setup
```

これにより、`~/.config/wails/defaults.yaml` に共有署名構成が書き込まれ、その構成は**すべてのプラットフォームで署名時に適用されます**（[構成の優先順位](#heading-4)を参照）。このウィザードの署名ステップでは、次の処理を行います。

- 対象とするプラットフォーム用にインストールされている署名ツールを、**どのホストからでも**検出し、使用中の OS 向けのインストールコマンド（例：`brew install gnupg`、`sudo apt install osslsigncode`、`winget install GnuPG.Gpg4win`）を表示します。
- macOS では、キーチェーン内の Developer ID 証明書を一覧表示します。
- Linux では、GPG キーを一覧表示し、新しいキーを**生成してエクスポート**できます。
- Windows では、OpenSSL を使用して**自己署名証明書を生成**できます（テスト用）。
- パスワードは**システムのキーチェーンに安全に保存されます**（Taskfile には保存されません）。

@note{type="tip"}
パスワードは、システムのネイティブ資格情報ストア（macOS Keychain、Windows Credential Manager、または Linux Secret Service）に保存されます。そのため、署名構成は安全に保たれ、すべての Wails プロジェクトで使用できます。

@end

## 構成の優先順位

署名タスクを実行すると、各署名オプションは次の順序で解決されます（最初に一致した値が使用されます）。

1. `wails3 tool sign` に明示的に渡されたフラグ（例：`--pgp-key`、`--certificate`、`--identity`）。
2. 対応する**プロジェクトの Taskfile 変数**（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`など）。
3. `~/.config/wails/defaults.yaml` 内の**グローバル構成**（`wails3 setup` によって書き込まれます）。

つまり、Taskfile 変数は**任意のオーバーライド**です。変数が設定されていない場合は、グローバルに構成されたキー、証明書、または署名 ID が使用されます。3 つのいずれからも値が指定されていない場合、署名コマンドは構成方法を示す明確なエラーを報告します。

## プロジェクト単位の構成

グローバルではなく単一のプロジェクトに署名を構成するには、そのプロジェクト内からプロジェクト単位のウィザードを実行します。このウィザードは、`vars` をそのプロジェクトの `build/<platform>/Taskfile.yml` ファイルに書き込みます。

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### 手動構成

別の方法として、プラットフォーム固有の Taskfile を手動で編集できます。各ファイルの先頭にある `vars` セクションを編集します。

@tabs
[macOS]
`build/darwin/Taskfile.yml` を編集します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

次に、以下を実行します。

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
`build/windows/Taskfile.yml` を編集します。

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

パスワードはシステムのキーチェーンから取得されます（構成するには `wails3 signing` を実行します）。

次に、以下を実行します。

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
`build/linux/Taskfile.yml` を編集します。

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

パスワードはシステムのキーチェーンから取得されます（構成するには `wails3 signing` を実行します）。

次に、以下を実行します。

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

システムから署名状態を直接確認することもできます。

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

すべてのプラットフォームの署名構成を対話形式でセットアップするには、ウィザードを使用します。

```bash
wails3 signing
```

## macOS のコード署名

### 前提条件

- Apple Developer Account（年間 $99）
- Developer ID Application 証明書
- Xcode Command Line Tools がインストール済みであること

### 署名 ID

利用可能な署名 ID を確認します。

```bash
security find-identity -v -p codesigning
```

出力：

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
App Store 外で配布するには、**Developer ID Application** 証明書が必要です。

@end

### 構成

`build/darwin/Taskfile.yml` を編集し、署名変数を設定します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | はい | Developer ID（例：「Developer ID Application: Your Company (TEAMID)」） |
| `KEYCHAIN_PROFILE` | 公証時に必須 | 保存済みの資格情報を含むキーチェーンプロファイル名 |
| `ENTITLEMENTS` | いいえ | entitlements ファイルへのパス |

次に、以下を実行します。

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### エンタイトルメント

エンタイトルメントは、アプリがアクセスできる機能を制御します。通常、Wails アプリでは開発用と本番用で異なるエンタイトルメントが必要です。

- **開発**：JIT、未署名メモリ、デバッグのエンタイトルメントが必要
- **本番**：最小限のエンタイトルメント（ネットワークアクセスのみ）

対話形式のセットアップウィザードを使用して、両方のファイルを生成します。

```bash
wails3 entitlements
```

これにより、以下が作成されます。

- `build/darwin/entitlements.dev.plist` - 開発ビルド用
- `build/darwin/entitlements.plist` - 本番／署名済みビルド用

**利用可能なプリセット：**\

| プリセット | 説明 |
| --- | --- |
| 開発 | JIT、未署名メモリ、デバッグ、ネットワーク |
| 本番 | ネットワークのみ（最小限で最も安全） |
| 両方 | 開発用と本番用の両方のファイルを作成（推奨） |
| App Store | ネットワークおよびファイルへのアクセスを許可したサンドボックスを有効化 |
| カスタム | 個々のエンタイトルメントを選択 |

@note{type="note"}
darwin Taskfile の `run` タスクは、`entitlements.dev.plist` を自動的に使用します。本番ビルドでは、`sign` タスクが `entitlements.plist` を使用します。

@end

次に、Taskfile の変数 `ENTITLEMENTS` に適切なファイルを指定します。

### 公証

Apple では、配布するすべてのアプリを公証する必要があります。

@steps
### **認証情報をキーチェーンに保存します**（初回のみのセットアップ）。`wails3 signing` を実行するか（値の入力を求め、内部で `notarytool` を呼び出します）、`notarytool` を直接呼び出します。
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Taskfile で KEYCHAIN_PROFILE を設定し**、上記のプロファイル名と一致させます。
### **アプリに署名して公証します**。
```bash
wails3 task darwin:sign:notarize
```

### **公証を検証します**。
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公証には通常 1-2 分かかります。チケットはアプリに自動的にステープルされます。

@end

## Windows コード署名

### 前提条件

- コード署名証明書（DigiCert、Sectigo などから取得）
- Windows でネイティブ署名する場合：Windows SDK がインストールされていること（`signtool.exe` 用）
- macOS/Linux からクロスプラットフォーム署名する場合：[`osslsigncode`](https://github.com/mtrojnar/osslsigncode)（`wails3 setup` の署名手順に、ホスト固有のインストールコマンドが表示されます）

### 自己署名証明書の生成（テスト用）

署名パイプラインのテストだけが必要な場合は、`wails3 setup` を実行し、署名手順の **Windows** タブを開いて、**自己署名証明書を生成**を選択します。OpenSSL を使用してコード署名用の `.pfx` を作成し、そのパスをグローバル設定に記録します。

@note{type="caution"}
自己署名証明書は**テストおよび内部配布専用**です。エンドユーザーには SmartScreen の警告が表示されます。一般公開リリースには、信頼された CA が発行した証明書が必要です。

@end

### 設定

`build/windows/Taskfile.yml` を編集し、署名用の変数を設定します。

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 上書き指定 | .pfx/.p12 証明書ファイルへのパス（未設定の場合はグローバル設定を使用） |
| `SIGN_THUMBPRINT` | 上書き指定 | Windows 証明書ストア内の証明書サムプリント（`SIGN_CERTIFICATE` の代替） |
| `TIMESTAMP_SERVER` | いいえ | タイムスタンプサーバーの URL（デフォルト：http://timestamp.digicert.com） |

@note{type="note"}
これらの変数は、任意の**上書き設定**です。未設定の場合は、`wails3 setup` でグローバルに設定された証明書が使用されます。[設定の優先順位](#heading-4)を参照してください。

@end

@note{type="note"}
証明書のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。設定するには `wails3 signing` を実行します。CI では、環境変数 `WAILS_WINDOWS_CERT_PASSWORD` を設定することもできます。

@end

次に、以下を実行します。

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### クロスプラットフォーム署名

Windows 実行ファイルは、どのプラットフォームからでも署名できます。同じ Taskfile の設定とコマンドが macOS と Linux でも機能します。

### サポートされる Windows 形式

| 形式 | 拡張子 | 備考 |
| --- | --- | --- |
| 実行ファイル | .exe | 標準的な PE 署名 |
| インストーラー | .msi | Windows Installer パッケージ |
| アプリパッケージ | .msix, .appx | モダン Windows アプリ |

## Linux パッケージ署名

Linux パッケージ（DEB および RPM）は、PGP/GPG 鍵を使用して署名されます。Windows や macOS のコード署名とは異なり、Linux のパッケージ署名は、コードが OS に信頼されていることではなく、パッケージが信頼できる配布元から提供されたことを証明します。

### 前提条件

- PGP 鍵ペア（Wails で生成可能）

### PGP 鍵の生成

最も簡単な方法はセットアップウィザードです。`wails3 setup` を実行し、署名手順の **Linux** タブを開いて、**新しい GPG 鍵を作成**を選択します。ウィザードは以下を行います。

- GPG キーリングに RSA 4096 鍵を生成します（無人実行が可能で CI に適した鍵にするには、パスフレーズを空欄のままにします）。
- **それを `~/.wails/signing/<keyid>.asc`**（ビルドが署名に使用するファイル）へエクスポートします。
- 鍵 ID とエクスポート先のパスを `~/.config/wails/defaults.yaml` に記録し、署名時に自動的に使用されるようにします。

@note{type="note"}
鍵が**ID のみ**で設定されている場合（たとえば、ウィザードが鍵を自動的にエクスポートするようになる前に作成された場合）、次に Linux タブを開いたときに鍵がファイルへエクスポートされ、パスが入力されます。手動操作は不要です。ビルドは鍵*ファイル*を使用して署名するため、重要なのはパスです。

@end

`gpg` を使用して手動で行うこともできます。

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

推奨設定：RSA 4096 ビット、有効期限 1 年、強力なパスワードで保護。

@note{type="caution"}
秘密鍵は安全に管理してください。暗号化して保存し、安全な場所にバックアップしてください。

@end

### 設定

`build/linux/Taskfile.yml`を編集し、署名変数を設定します。

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `PGP_KEY` | 上書き | エクスポートした PGP 秘密鍵ファイルへのパス（未設定の場合は、`wails3 setup`でグローバルに設定した鍵を使用） |
| `SIGN_ROLE` | いいえ | DEB の署名ロール（デフォルト：builder） |

@note{type="note"}
PGP 鍵のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。設定するには`wails3 signing`を実行するか、CI で環境変数`WAILS_PGP_PASSWORD`を設定します。

@end

次に、以下を実行します。

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB の署名ロール

DEB パッケージでは、`SIGN_ROLE`を使用して署名ロールを指定できます。

- `origin`：パッケージの配布元による署名
- `maint`：パッケージメンテナーによる署名
- `archive`：アーカイブメンテナーによる署名
- `builder`：パッケージビルダーによる署名（デフォルト）

### クロスプラットフォーム署名

Linux パッケージは、どのプラットフォームからでも署名できます。同じ Taskfile の設定とコマンドが Windows と macOS でも動作します。

### 鍵情報の表示

```bash
gpg --show-keys signing-key.asc
```

出力：

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Linux パッケージの検証

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### 公開鍵の配布

ユーザーがパッケージを検証するには、公開鍵が必要です。

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

## GitHub Actions との統合

CI 環境では、システムのキーチェーンではなく環境変数を使用してパスワードを指定します。

| 環境変数 | 説明 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 証明書のパスワード |
| `WAILS_PGP_PASSWORD` | Linux パッケージ用 PGP 鍵のパスワード |

Taskfile の変数を直接渡すこともできます。

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### macOS ワークフロー

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

### Windows ワークフロー

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

### クロスプラットフォームワークフロー（Linux ランナー）

単一の Linux ランナーから Windows パッケージと Linux パッケージに署名します。

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
クロスプラットフォーム署名に Linux ランナーを使用すると、Windows ランナーを別途用意する必要がなくなり、CI/CD を簡素化できます。macOS の署名には Apple のネイティブツールが必要なため、引き続き macOS ランナーが必要です。

@end

## CLI リファレンス

### wails3 signing

プロジェクトの署名を設定するための対話形式ウィザードです。

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

ウィザードでは、以下の設定を順に行います。

- **macOS**：Developer ID 証明書を選択し、公証用の認証情報を設定します（`xcrun notarytool store-credentials`を呼び出します）。
- **Windows**：証明書ファイルまたは拇印のどちらを使用するか選択し、パスワードとタイムスタンプサーバーを設定します。
- **Linux**：既存の PGP 鍵を使用するか、新しい鍵を生成し（`gpg`を呼び出します）、署名ロールを設定します。

### wails3 entitlements

macOS のエンタイトルメントを設定するための対話形式ウィザードです。

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**プリセット：**

- **Development**：JIT、デバッグ、ネットワーク用のエンタイトルメントを含む`entitlements.dev.plist`を作成します
- **Production**：必要最小限のエンタイトルメントを含む`entitlements.plist`を作成します
- **Both**：両方のファイルを作成します（推奨）
- **App Store**：Mac App Store 向けにサンドボックス化されたエンタイトルメントを作成します
- **Custom**：個々のエンタイトルメントと出力先ファイルを選択します

### wails3 sign

現在のプラットフォームまたは指定したプラットフォーム向けのバイナリとパッケージに署名します。適切なプラットフォーム固有の署名タスクを呼び出すラッパーです。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

これにより、Taskfile の署名設定を使用する、対応する`<platform>:sign`タスクが実行されます。

### wails3 tool sign

指定したファイルに直接署名するための低レベルコマンドです。Taskfile が内部で使用します。

```bash
wails3 tool sign [flags]
```

**共通フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--input` | 署名するファイルへのパス |
| `--output` | 出力先のパス（省略可能。デフォルトでは元のファイルを直接更新） |
| `--verbose` | 詳細出力を有効にする |

**Windows/macOS 用フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--certificate` | PKCS#12（.pfx/.p12）証明書へのパス |
| `--password` | 証明書のパスワード |
| `--timestamp` | タイムスタンプサーバーの URL |

**macOS 固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--identity` | 署名 ID（アドホック署名には「-」を使用） |
| `--entitlements` | エンタイトルメント plist へのパス |
| `--hardened-runtime` | Hardened Runtime を有効にする（デフォルト：true） |
| `--notarize` | 公証のために送信する |
| `--keychain-profile` | 公証用のキーチェーンプロファイル |

**Windows 固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--thumbprint` | Windows ストア内の証明書の拇印 |

**Linux 固有のフラグ:**\

| フラグ | 説明 |
| --- | --- |
| `--pgp-key` | PGP 秘密鍵へのパス |
| `--pgp-password` | PGP 鍵のパスワード |
| `--role` | DEB 署名ロール（origin/maint/archive/builder） |

### 署名状態の確認（ネイティブツール）

`wails3 signing` は署名用認証情報を設定します。署名の状態を確認するには、ネイティブツールを直接使用してください。

| タスク | コマンド |
| --- | --- |
| macOS のコード署名 ID を一覧表示する | `security find-identity -v -p codesigning` |
| 公証用の認証情報を保存する | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP 鍵ファイルを確認する | `gpg --show-keys <key.asc>` |
| PGP 鍵ペアを生成する | `gpg --full-generate-key` |
| 公開鍵をエクスポートする | `gpg --armor --export <email>` |

## トラブルシューティング

### macOS の問題

**「Developer ID 証明書が見つからない」**

- 証明書がキーチェーンにインストールされていることを確認してください
- `security find-identity -v -p codesigning`を使用して、証明書の有効期限が切れていないことを確認してください
- 「Apple Development」証明書だけでなく、「Developer ID Application」証明書があることを確認してください

**「公証に失敗した」**

- 公証ログを確認してください: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Hardened Runtime が有効になっていることを確認してください
- アプリに未署名のバイナリが含まれていないことを確認してください

**「コード署名に失敗した」**

- キーチェーンがロック解除されていることを確認してください: `security unlock-keychain`
- アプリバンドルのファイル権限を確認してください

### Windows の問題

**「証明書が見つからない」**

- 証明書のパスが正しいことを確認してください
- 証明書のパスワードを確認してください
- 証明書が有効であること（有効期限切れまたは失効していないこと）を確認してください

**「タイムスタンプサーバーエラー」**

- 別のタイムスタンプサーバーを試してください:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux の問題

**「無効な PGP 鍵」**

- 鍵ファイルが ASCII Armor 形式であることを確認してください
- `gpg --show-keys <key.asc>`を使用して、鍵の有効期限が切れていないことを確認してください
- パスワードが正しいことを確認してください

**「署名の検証に失敗した」**

- 公開鍵が正しくインポートされていることを確認してください
- 署名後にパッケージが変更されていないことを確認してください

## その他のリソース

### 公式ドキュメント

- [Apple コード署名ガイド](https://developer.apple.com/support/code-signing/)
- [Apple 公証ドキュメント](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft コード署名](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian パッケージ署名](https://wiki.debian.org/SecureApt)
- [RPM パッケージ署名](https://rpm-software-management.github.io/rpm/manual/signatures.html)
