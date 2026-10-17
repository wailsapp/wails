---
title: "コード署名"
description: "あらゆるプラットフォームで Wails アプリケーションに署名するためのガイド"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## アプリケーションのコード署名

このガイドでは、macOS、Windows、Linux 向けの Wails アプリケーションに署名する方法を説明します。Wails v3 には、コード署名、公証、PGP キー管理のための CLI ツールが組み込まれています。

- **macOS** - macOS アプリケーションへの署名と公証
- **Windows** - Windows の実行ファイルとパッケージへの署名
- **Linux** - PGP キーによる DEB および RPM パッケージへの署名

## クロスプラットフォーム署名マトリックス

次のマトリックスは、各ソースプラットフォームから署名できる対象を示しています。

| 対象形式 | Windows から | macOS から | Linux から |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS .app バンドル | ❌ | ✅ | ❌ |
| macOS の公証 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Windows および Linux のパッケージには、<strong>どのプラットフォームからでも</strong>署名できます。macOS の署名には Apple のツールが必要なため、Mac が必要です。

@end

### 署名バックエンド

Wails は、利用可能な最適な署名バックエンドを自動的に選択します。

| プラットフォーム | ネイティブバックエンド | クロスプラットフォームバックエンド |
| --- | --- | --- |
| Windows | `signtool.exe`（Windows SDK） | 組み込み |
| macOS | `codesign`（Xcode） | 利用不可 |
| Linux | 該当なし | 組み込み |

ネイティブプラットフォーム上で実行する場合、Wails は互換性を最大限に確保するため、ネイティブツールを使用します。クロスコンパイル時には、組み込みの署名サポートを使用します。

## クイックスタート

署名を最も短時間で設定する方法は、セットアップウィザードを使用することです。

```bash
wails3 setup
```

これにより、共有の署名設定が `~/.config/wails/defaults.yaml` に書き込まれ、 **すべてのプラットフォームで署名時に使用されます**（[設定の優先順位](#heading-4)を参照）。 この署名手順では、次の処理を行います。

- 対象とするプラットフォーム向けの署名ツールがインストールされているかを、<strong>どのホストからでも</strong>検出し、使用中の OS に対応するインストールコマンド（例：`brew install gnupg`、`sudo apt install osslsigncode`、`winget install GnuPG.Gpg4win`）を表示します。
- macOS では、キーチェーン内の Developer ID 証明書を一覧表示します。
- Linux 向けの署名では、GPG キーを一覧表示し、新しいキーを<strong>生成してエクスポート</strong>できます。
- Windows 向けの署名では、OpenSSL を使用して（テスト用の）<strong>自己署名証明書を生成</strong>できます。
- **パスワードをシステムのキーチェーンに安全に保存します**（Taskfile には保存しません）。

@note{type="tip"}
パスワードは、システムのネイティブな資格情報ストア（macOS Keychain、Windows Credential Manager、または Linux Secret Service）に保存されます。そのため、署名設定を安全に保ちながら、すべての Wails プロジェクトで使用できます。

@end

## 設定の優先順位

署名タスクを実行すると、各署名オプションは次の順序で解決されます（最初に 一致した値が使用されます）。

1. `wails3 tool sign` に明示的に渡されたフラグ（例：`--pgp-key`、`--certificate`、`--identity`）。
2. 対応する<strong>プロジェクトの Taskfile 変数</strong>（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`、…）。
3. `~/.config/wails/defaults.yaml` 内の<strong>グローバル設定</strong>（`wails3 setup` によって書き込まれます）。

つまり、Taskfile 変数は<strong>任意の上書き設定</strong>です。変数が未設定の場合は、 グローバルに設定されたキー、証明書、または ID が使用されます。3 つのいずれにも 値がない場合、sign コマンドは設定方法を示す明確なエラーを報告します。

## プロジェクト単位の設定

グローバルではなく単一のプロジェクトに署名を設定するには、プロジェクト内からプロジェクト単位の ウィザードを実行します。これにより、`vars` がそのプロジェクトの `build/<platform>/Taskfile.yml` ファイルに書き込まれます。

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### 手動設定

または、プラットフォーム固有の Taskfile を手動で編集できます。各ファイルの先頭にある `vars` セクションを編集します。

@tabs
[macOS]
`build/darwin/Taskfile.yml` を編集します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

次に実行します。

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

パスワードはシステムのキーチェーンから取得されます（設定するには `wails3 setup signing` を実行してください）。

次に実行します：

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
`build/linux/Taskfile.yml`を編集します：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

パスワードはシステムのキーチェーンから取得されます（設定するには`wails3 setup signing`を実行してください）。

次に実行します：

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

システムから署名の状態を直接確認することもできます：

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

すべてのプラットフォームの署名設定を対話形式で行うには、ウィザードを使用します：

```bash
wails3 setup signing
```

## macOSコード署名

### 前提条件

- Apple Developerアカウント（年間$99）
- Developer ID Application証明書
- Xcode Command Line Toolsがインストール済みであること

### 署名ID

利用可能な署名IDを確認します：

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
App Store以外で配布するには、<strong>Developer ID Application</strong>証明書が必要です。

@end

### 設定

`build/darwin/Taskfile.yml`を編集し、署名用の変数を設定します：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | はい | Developer ID（例：「Developer ID Application: Your Company (TEAMID)」） |
| `KEYCHAIN_PROFILE` | 公証に必要 | 認証情報が保存されたキーチェーンプロファイル名 |
| `ENTITLEMENTS` | いいえ | エンタイトルメントファイルへのパス |

次に実行します：

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### エンタイトルメント

エンタイトルメントは、アプリがアクセスできる機能を制御します。Wailsアプリでは通常、開発環境と本番環境で異なるエンタイトルメントが必要です：

- **開発**：JIT、未署名メモリ、およびデバッグ用のエンタイトルメントが必要です
- **本番**：最小限のエンタイトルメント（ネットワークアクセスのみ）

対話形式のセットアップウィザードを使用して、両方のファイルを生成します：

```bash
wails3 setup entitlements
```

次のファイルが作成されます：

- `build/darwin/entitlements.dev.plist` - 開発ビルド用
- `build/darwin/entitlements.plist` - 本番／署名済みビルド用

**利用可能なプリセット：**\

| プリセット | 説明 |
| --- | --- |
| 開発 | JIT、未署名メモリ、デバッグ、ネットワーク |
| 本番 | ネットワークのみ（最小限で最も安全） |
| 両方 | 開発用と本番用の両方のファイルを作成（推奨） |
| App Store | ネットワークおよびファイルへのアクセスを許可したサンドボックスを有効化 |
| カスタム | 個別のエンタイトルメントを選択 |

@note{type="note"}
darwinのTaskfileにある`run`タスクは、`entitlements.dev.plist`を自動的に使用します。`sign`タスクは、本番ビルドに`entitlements.plist`を使用します。

@end

次に、Taskfileのvarsにある`ENTITLEMENTS`を、適切なファイルを参照するように設定します。

### 公証

Appleでは、配布するすべてのアプリを公証する必要があります。

@steps
### **認証情報をキーチェーンに保存します**（初回のみ）。`wails3 setup signing`を実行する（これらの値の入力を求め、内部で`notarytool`を呼び出します）か、`notarytool`を直接呼び出します：
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **TaskfileのKEYCHAIN_PROFILEを設定し**、上記のプロファイル名と一致させます。
### **アプリに署名して公証します**：
```bash
wails3 task darwin:sign:notarize
```

### **公証を検証します**：
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公証には通常1-2分かかります。チケットはアプリに自動的にステープルされます。

@end

## Windowsコード署名

### 前提条件

- コード署名証明書（DigiCert、Sectigoなどから取得）
- Windowsでネイティブ署名を行う場合：Windows SDKがインストール済みであること（`signtool.exe`用）
- macOS/Linux からクロスプラットフォーム署名を行う場合：[`osslsigncode`](https://github.com/mtrojnar/osslsigncode)（`wails3 setup` の署名設定ステップで、ホスト固有のインストールコマンドが表示されます）

### 自己署名証明書の生成（テスト用）

署名パイプラインをテストするだけであれば、`wails3 setup` を実行し、署名手順の **Windows** タブを開いて、<strong>自己署名証明書を生成</strong>を選択します。これにより、OpenSSL を使用してコード署名用の `.pfx` が作成され、そのパスがグローバル設定に記録されます。

@note{type="caution"}
自己署名証明書は、<strong>テストおよび社内配布専用</strong>です。エンドユーザーには SmartScreen の警告が表示されます。一般公開するリリースには、信頼された CA が発行した証明書が必要です。

@end

### 設定

`build/windows/Taskfile.yml` を編集し、署名変数を設定します：

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
| `SIGN_THUMBPRINT` | 上書き指定 | Windows 証明書ストア内の証明書の拇印（`SIGN_CERTIFICATE` の代替） |
| `TIMESTAMP_SERVER` | いいえ | タイムスタンプサーバーの URL（デフォルト：http://timestamp.digicert.com） |

@note{type="note"}
これらの変数は省略可能な<strong>上書き指定</strong>です。未設定の場合は、`wails3 setup` でグローバルに設定した証明書が使用されます。[設定の優先順位](#heading-4)を参照してください。

@end

@note{type="note"}
証明書のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。設定するには `wails3 setup signing` を実行します。CI では `WAILS_WINDOWS_CERT_PASSWORD` 環境変数を設定することもできます。

@end

次に、以下を実行します：

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### クロスプラットフォーム署名

Windows 実行可能ファイルは、どのプラットフォームからでも署名できます。同じ Taskfile の設定とコマンドを macOS および Linux でも使用できます。

### サポートされる Windows 形式

| 形式 | 拡張子 | 備考 |
| --- | --- | --- |
| 実行可能ファイル | .exe | 標準の PE 署名 |
| インストーラー | .msi | Windows Installer パッケージ |
| アプリパッケージ | .msix, .appx | モダン Windows アプリ |

## Linux パッケージの署名

Linux パッケージ（DEB および RPM）は、PGP/GPG 鍵を使用して署名します。Windows や macOS のコード署名とは異なり、Linux パッケージ署名はコードが OS によって信頼されていることではなく、パッケージが信頼できる提供元から配布されたことを証明します。

### 前提条件

- PGP 鍵ペア（Wails で生成可能）

### PGP 鍵の生成

最も簡単なのはセットアップウィザードを使用する方法です。`wails3 setup` を実行し、署名手順の **Linux** タブを開いて、<strong>新しい GPG 鍵を作成</strong>を選択します。ウィザードは次の処理を行います：

- GPG キーリングに RSA 4096 鍵を生成します（無人環境や CI で使用しやすい鍵にするには、パスフレーズを空のままにします）。
- **鍵を `~/.wails/signing/<keyid>.asc`** にエクスポートします（ビルドではこのファイルを使用して署名します）。
- 鍵 ID とエクスポート先のパスの両方を `~/.config/wails/defaults.yaml` に記録し、署名時に自動的に使用されるようにします。

@note{type="note"}
鍵が<strong>ID のみ</strong>で設定されている場合（たとえば、ウィザードが鍵を自動的にエクスポートするようになる前に作成された場合）、次に Linux タブを開いたときに鍵がファイルへエクスポートされ、パスが入力されます。手動での操作は不要です。ビルドでは鍵の<em>ファイル</em>を使用して署名するため、重要なのはパスです。

@end

`gpg` を使用して手動で行うこともできます：

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

推奨設定：RSA 4096 ビット、有効期限 1 年、強力なパスワードで保護。

@note{type="caution"}
秘密鍵は厳重に管理してください。暗号化して保存し、安全にバックアップしてください。

@end

### 設定

`build/linux/Taskfile.yml` を編集し、署名変数を設定します：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `PGP_KEY` | 上書き | エクスポートした PGP 秘密鍵ファイルへのパス（未設定の場合は、`wails3 setup` でグローバルに設定された鍵を使用） |
| `SIGN_ROLE` | いいえ | DEB 署名ロール（デフォルト：builder） |

@note{type="note"}
PGP 鍵のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。設定するには `wails3 setup signing` を実行します。CI では `WAILS_PGP_PASSWORD` 環境変数を設定してください。

@end

次に、以下を実行します。

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 署名ロール

DEB パッケージでは、`SIGN_ROLE` を使用して署名ロールを指定できます。

- `origin`：パッケージの提供元による署名
- `maint`：パッケージのメンテナーによる署名
- `archive`：アーカイブのメンテナーによる署名
- `builder`：パッケージのビルダーによる署名（デフォルト）

### クロスプラットフォーム署名

Linux パッケージは、どのプラットフォームからでも署名できます。同じ Taskfile の設定とコマンドを Windows および macOS でも使用できます。

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

CI 環境では、システムのキーチェーンではなく環境変数を使用してパスワードを渡します。

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
クロスプラットフォーム署名に Linux ランナーを使用すると、個別の Windows ランナーが不要になり、CI/CD が簡素化されます。macOS の署名には Apple のネイティブツールが必要なため、引き続き macOS ランナーが必要です。

@end

## CLI リファレンス

### wails3 setup signing

プロジェクトの署名を設定する対話型ウィザードです。

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

ウィザードでは、以下の設定を順に行います。

- **macOS**：Developer ID 証明書の選択と、公証用認証情報の設定（`xcrun notarytool store-credentials` を呼び出します）。
- **Windows**：証明書ファイルまたは拇印の選択と、パスワードおよびタイムスタンプサーバーの設定。
- **Linux**：既存の PGP 鍵の使用または新しい鍵の生成（`gpg` を呼び出します）と、署名ロールの設定。

### wails3 setup entitlements

macOS のエンタイトルメントを設定する対話型ウィザードです。

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**プリセット：**

- **開発**：JIT、デバッグ、ネットワーク用のエンタイトルメントを含む `entitlements.dev.plist` を作成します
- **本番**：最小限のエンタイトルメントを含む `entitlements.plist` を作成します
- **両方**：両方のファイルを作成します（推奨）
- **App Store**：Mac App Store 用のサンドボックス化されたエンタイトルメントを作成します
- **カスタム**：個別のエンタイトルメントと対象ファイルを選択します

### wails3 sign

現在のプラットフォームまたは指定したプラットフォーム向けのバイナリとパッケージに署名します。これは、プラットフォームに対応する署名タスクを呼び出すラッパーです。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

これにより、Taskfile の署名設定を使用する、対応する `<platform>:sign` タスクが実行されます。

### wails3 tool sign

指定したファイルに直接署名するための低レベルコマンドです。Taskfile から内部的に使用されます。

```bash
wails3 tool sign [flags]
```

**共通フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--input` | 署名するファイルのパス |
| `--output` | 出力パス（省略可能。デフォルトでは元のファイルを直接更新） |
| `--verbose` | 詳細出力を有効にする |

**Windows/macOS用フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--certificate` | PKCS#12（.pfx/.p12）証明書のパス |
| `--password` | 証明書のパスワード |
| `--timestamp` | タイムスタンプサーバーのURL |

**macOS固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--identity` | 署名ID（アドホック署名には「-」を使用） |
| `--entitlements` | エンタイトルメントplistのパス |
| `--hardened-runtime` | Hardened Runtimeを有効にする（デフォルト：true） |
| `--notarize` | 公証のために提出する |
| `--keychain-profile` | 公証用のキーチェーンプロファイル |

**Windows固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--thumbprint` | Windows証明書ストア内の証明書サムプリント |

**Linux固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--pgp-key` | PGP秘密鍵のパス |
| `--pgp-password` | PGP鍵のパスワード |
| `--role` | DEB署名ロール（origin/maint/archive/builder） |

### 署名状態の確認（ネイティブツール）

v3 には `wails3 signing` コマンドは<strong>ありません</strong>。署名状態を確認するには、ネイティブツールを直接使用してください：

| タスク | コマンド |
| --- | --- |
| macOSのコード署名IDを一覧表示する | `security find-identity -v -p codesigning` |
| 公証用の認証情報を保存する | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP鍵ファイルを確認する | `gpg --show-keys <key.asc>` |
| PGP鍵ペアを生成する | `gpg --full-generate-key` |
| 公開鍵をエクスポートする | `gpg --armor --export <email>` |

## トラブルシューティング

### macOSの問題

**「Developer ID証明書が見つかりません」**

- 証明書がキーチェーンにインストールされていることを確認してください
- `security find-identity -v -p codesigning`を使用して、証明書の有効期限が切れていないことを確認してください
- 「Apple Development」証明書だけでなく、「Developer ID Application」証明書があることを確認してください

**「公証に失敗しました」**

- 公証ログを確認してください：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Hardened Runtimeが有効になっていることを確認してください
- アプリに未署名のバイナリが含まれていないことを確認してください

**「コード署名に失敗しました」**

- キーチェーンがロック解除されていることを確認してください：`security unlock-keychain`
- アプリバンドルのファイル権限を確認してください

### Windows の問題

**「証明書が見つかりません」**

- 証明書のパスが正しいことを確認してください
- 証明書のパスワードを確認してください
- 証明書が有効であること（有効期限切れまたは失効済みでないこと）を確認してください

**「タイムスタンプサーバーエラー」**

- 別のタイムスタンプサーバーを試してください：
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux の問題

**「無効な PGP キー」**

- キーファイルが ASCII Armor 形式であることを確認してください
- `gpg --show-keys <key.asc>`を使用して、キーの有効期限が切れていないことを確認してください
- パスワードが正しいことを確認してください

**「署名の検証に失敗しました」**

- 公開鍵が正しくインポートされていることを確認してください
- 署名後にパッケージが変更されていないことを確認してください

## その他のリソース

### 公式ドキュメント

- [Apple コード署名ガイド](https://developer.apple.com/support/code-signing/)
- [Apple 公証ドキュメント](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft コード署名](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian パッケージ署名](https://wiki.debian.org/SecureApt)
- [RPM パッケージ署名](https://rpm-software-management.github.io/rpm/manual/signatures.html)
