---
title: "コード署名"
description: "すべてのプラットフォームで Wails アプリケーションに署名するためのガイド"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## アプリケーションのコード署名

このガイドでは、macOS、Windows、Linux 向けの Wails アプリケーションに署名する方法を説明します。Wails v3 には、コード署名、公証、PGP キー管理用の CLI ツールが組み込まれています。

- **macOS** - macOS アプリケーションへの署名と公証
- **Windows** - Windows の実行ファイルとパッケージへの署名
- **Linux** - PGP キーを使用した DEB および RPM パッケージへの署名

## クロスプラットフォーム署名マトリックス

次のマトリックスは、各ソースプラットフォームから署名できる対象を示しています。

| 対象形式 | Windows から | macOS から | Linux から |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX、エクスポート済み PFX | ✅ | ✅ | ✅ |
| Windows MSI または証明書ストア／プロバイダーの ID | ✅ | ❌ | ❌ |
| macOS .app / Mach-O / DMG、エクスポート済み ID | ✅ | ✅ | ✅ |
| macOS 公証、チーム API 認証情報 | ✅ | ✅ | ✅ |
| macOS キーチェーン ID または公証プロファイル | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
エクスポート済みの Apple ID とチーム API 認証情報を使用すると、ポータブルな署名と公証が可能になります。macOS キーチェーン ID と公証プロファイルには Mac が必要です。形式固有の制限については、以下で説明します。

@end

### 署名バックエンド

Wails は、利用可能な最適な署名バックエンドを自動的に選択します。

| プラットフォーム | ネイティブバックエンド | クロスプラットフォームバックエンド |
| --- | --- | --- |
| Windows | ストア／プロバイダーの ID、MSI、または従来の RC2 PFX には `signtool.exe` | エクスポート済み RSA PFX および PE/MSIX には組み込みバックエンド |
| macOS | キーチェーン ID には `codesign` | エクスポート済み PKCS#12 ID またはアドホック署名には組み込みバックエンド |
| Linux | 該当なし | 組み込み |

どのホストでも、認証情報の種類によって Apple 署名バックエンドが選択されます。エクスポート済み ID には Go を使用し、macOS キーチェーン ID には `codesign` を使用します。チケットの取得、検証、ステープルには Go を使用します。

## クイックスタート

署名を最も速く設定するには、セットアップウィザードを使用します。

```bash
wails3 setup
```

これにより、共有の署名設定が `~/.config/wails/defaults.yaml` に書き込まれ、**すべてのプラットフォームで署名時に適用されます**（[設定の優先順位](#heading-4)を参照）。このウィザードの署名手順では、次の処理を行います。

- OpenSSL や GnuPG など、証明書やキーを生成またはエクスポートするためのセットアップツールを表示します。エクスポート済みの Apple および Windows の ID には組み込み署名ツールを使用します。Windows の証明書ストア／プロバイダーの ID と MSI パッケージには Windows SignTool が必要です。
- macOS では、キーチェーン内の Developer ID 証明書を一覧表示します。
- Linux では、GPG キーを一覧表示し、新しいキーを**生成してエクスポート**できます。
- Windows では、OpenSSL を使用して（テスト用の）**自己署名証明書を生成**できます。
- **パスワードをシステムのキーチェーンに安全に保存します**（Taskfile には保存しません）。

@note{type="tip"}
パスワードは、システムのネイティブな認証情報ストア（macOS Keychain、Windows Credential Manager、または Linux Secret Service）に保存されます。これにより、署名設定の安全性が確保され、すべての Wails プロジェクトで使用できます。

@end

## 設定の優先順位

署名タスクを実行すると、各署名オプションは次の順序で解決されます（最初に一致した値が優先されます）。

1. `wails3 tool sign` に明示的に渡したフラグ（例：`--pgp-key`、`--certificate`、`--identity`）。
2. 対応する**プロジェクト Taskfile 変数**（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`、…）。
3. `~/.config/wails/defaults.yaml` 内の**グローバル設定**（`wails3 setup` によって書き込まれます）。

つまり、Taskfile 変数は**任意の上書き設定**です。変数が未設定の場合は、グローバルに設定されたキー、証明書、または ID が使用されます。3 つのいずれにも値がない場合、署名コマンドは設定方法を示す明確なエラーを報告します。

## プロジェクト単位の設定

署名をグローバルではなく単一のプロジェクトに設定するには、そのプロジェクト内でプロジェクト単位のウィザードを実行します。このウィザードは、そのプロジェクトの `build/<platform>/Taskfile.yml` ファイルに `vars` を書き込みます。

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
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

パスワードはシステムのキーチェーンから取得されます（設定するには `wails3 signing` を実行します）。

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

パスワードはシステムのキーチェーンから取得されます（設定するには `wails3 signing` を実行します）。

次に、以下を実行します。

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

署名の状態をシステムから直接確認することもできます。

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

すべてのプラットフォームの署名設定を対話形式でセットアップするには、ウィザードを使用します。

```bash
wails3 signing
```

## macOS のコード署名

### 前提条件

- Apple Developer Account（年間 $99）
- Developer ID Application 証明書
- macOS キーチェーン ID を使用する場合：Xcode Command Line Tools がインストールされていること
- ポータブル署名を使用する場合：エクスポート済みの RSA PKCS#12 ID。Apple の署名ツールは不要です

### Wake を使用したポータブル署名

Wake は `wails.yaml` を直接読み取ります。Taskfile は実行しません。エクスポート済みの Developer ID と App Store Connect チーム API キーを設定します。

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

`credential` には、証明書のパスワードを保持する環境変数の名前を指定します。両方の秘密キーファイルをバージョン管理の対象外にしてください。[YAML ビルドガイド](/experimental/yaml-builds/)で説明されているリリースプロファイルで、署名と公証を有効にします。

ポータブル署名ツールは、64 ビットの Intel および Apple Silicon Mach-O ファイル、ユニバーサルバイナリ、ネストされたコードを含む標準的なアプリバンドル、フレームワーク、ならびに 4 GiB 未満の UDIF DMG をサポートします。未対応のレイアウトは、出力を公開する前に拒否されます。アドホック署名では `identity: "-"` を使用するため、公証できません。PKCS#12 ID では、AES/PBKDF2 または従来の 3DES 暗号化を使用できます。RC2 でエクスポートされたものは、サポート対象の暗号方式で再エクスポートする必要があります。

ポータブルな公証済み署名には、明示的な SHA-2 RFC 3161 タイムスタンプサービスが必要です。従来の SHA-1 タイムスタンプ署名は拒否されます。Team API による送信では、最大 5 GiB の単一アップロードをサポートします。個人用 API キーとマルチパートアップロードはサポートされません。チケットの検証とステープルには、Apple Root G3 を固定した Apple の観測済みサービスプロトコルを使用し、`xcrun stapler` は使用しません。配布時の信頼性と失効については、引き続き Gatekeeper が最終的な判断を行います。

以降のセクションでは、従来のプロジェクトと、キーチェーン ID を選択する Wake 構成で使用される既存の macOS キーチェーンワークフローについて説明します。

### 署名 ID

利用可能な署名 ID を確認します。

```bash
security find-identity -v -p codesigning
```

出力:

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

`build/darwin/Taskfile.yml` を編集し、署名用変数を設定します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | はい | Developer ID（例: 「Developer ID Application: Your Company (TEAMID)」） |
| `KEYCHAIN_PROFILE` | 公証時に必須 | 保存済みの認証情報を持つキーチェーンプロファイル名 |
| `ENTITLEMENTS` | いいえ | エンタイトルメントファイルへのパス |

次に、以下を実行します。

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### エンタイトルメント

エンタイトルメントは、アプリがアクセスできる機能を制御します。通常、Wails アプリでは開発環境と本番環境で異なるエンタイトルメントが必要です。

- **開発**: JIT、未署名メモリ、およびデバッグ用のエンタイトルメントが必要
- **本番**: 最小限のエンタイトルメント（ネットワークアクセスのみ）

対話形式のセットアップウィザードを使用して、両方のファイルを生成します。

```bash
wails3 entitlements
```

次のファイルが作成されます。

- `build/darwin/entitlements.dev.plist` - 開発ビルド用
- `build/darwin/entitlements.plist` - 本番／署名済みビルド用

**利用可能なプリセット:**\

| プリセット | 説明 |
| --- | --- |
| 開発 | JIT、未署名メモリ、デバッグ、ネットワーク |
| 本番 | ネットワークのみ（最小限で最も安全） |
| 両方 | 開発用と本番用の両方のファイルを作成（推奨） |
| App Store | ネットワークおよびファイルへのアクセスを許可したサンドボックスを有効化 |
| カスタム | 個別のエンタイトルメントを選択 |

@note{type="note"}
darwin Taskfile の `run` タスクは、`entitlements.dev.plist` を自動的に使用します。`sign` タスクは、本番ビルドに `entitlements.plist` を使用します。

@end

次に、Taskfile の変数 `ENTITLEMENTS` に適切なファイルを指定します。

### 公証

Apple は、配布するすべてのアプリに公証を義務付けています。

@steps
### **認証情報をキーチェーンに保存します**（初回のみのセットアップ）。`wails3 signing` を実行するか（これらの値の入力を求め、内部で `notarytool` を呼び出します）、`notarytool` を直接呼び出します。
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Taskfile の KEYCHAIN_PROFILE を設定し**、上記のプロファイル名と一致させます。
### **アプリに署名して公証します**:
```bash
wails3 task darwin:sign:notarize
```

### **公証を検証します**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
通常、公証には 1-2 分かかります。チケットはアプリに自動的にステープルされます。

@end

## Windows コード署名

### 前提条件

- コード署名証明書（DigiCert、Sectigo などから取得）
- ポータブル署名の場合: AES/PBKDF2 または従来の 3DES で暗号化してエクスポートした RSA PFX ID。外部署名ツールは不要です
- 証明書ストア／プロバイダーの ID、MSI パッケージ、または従来の RC2 PFX の場合: `PATH` 上に Windows SDK の署名ツールを備えた Windows

### 自己署名証明書の生成（テスト用）

署名パイプラインをテストするだけであれば、`wails3 setup` を実行し、署名手順の **Windows** タブを開いて、**自己署名証明書を生成** を選択します。これにより、OpenSSL を使用してコード署名用の `.pfx` が作成され、そのパスがグローバル構成に記録されます。

@note{type="caution"}
自己署名証明書は、**テストおよび社内配布専用**です。エンドユーザーには SmartScreen の警告が表示されます。一般公開するリリースには、信頼できる CA が発行した証明書が必要です。

@end

### 構成

`build/windows/Taskfile.yml` を編集し、署名用変数を設定します。

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 上書き | .pfx/.p12 証明書ファイルへのパス（未設定の場合はグローバル構成を使用） |
| `SIGN_THUMBPRINT` | 上書き | Windows 証明書ストア内の証明書サムプリント（`SIGN_CERTIFICATE` の代替） |
| `TIMESTAMP_SERVER` | いいえ | タイムスタンプサーバーの URL（デフォルト: http://timestamp.digicert.com） |

@note{type="note"}
これらの変数は任意の**上書き設定**です。未設定の場合は、`wails3 setup` でグローバルに構成した証明書が使用されます。[構成の優先順位](#heading-4)を参照してください。

@end

@note{type="note"}
証明書のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。`wails3 signing` を実行して構成するか、CI で環境変数 `WAILS_WINDOWS_CERT_PASSWORD` を設定します。

@end

次に、以下を実行します。

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### クロスプラットフォーム署名

エクスポートされた RSA PFX ID を使用すると、Windows、macOS、Linux 上で Go により EXE/DLL ファイル、NSIS インストーラー実行ファイル、個別の MSIX/APPX パッケージに署名できます。既存の署名は置き換え可能です。署名処理では出力をステージングし、署名とタイムスタンプの検証が成功した後にのみ公開するため、失敗時には元のファイルが保持されます。署名済みファイルを別ファイルとして書き込むには、`--output` を使用します。

Wake は、次の設定を `wails.yaml` から直接読み取ります。

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential` にはパスワード自体ではなく、パスワードを格納する環境変数の名前を指定します。リリースターゲットでは `sign: true` を有効にしてください。[YAML ビルドガイド](/experimental/yaml-builds/)を参照してください。Wake は Taskfile を呼び出しません。従来のプロジェクトでは、上記の Taskfile 設定を引き続き使用できます。

ポータブルプロファイルは RSA/SHA-256 を使用し、タイムスタンプ署名者の信頼された証明書チェーンを含む RFC 3161 タイムスタンプを検証します。デフォルトのサービスは DigiCert です。MSIX Publisher は署名証明書のサブジェクトと一致する必要があります。パッケージバンドルと、4 GiB 以上のアーカイブはサポートされません。PFX ファイルは 16 MiB までに制限されます。ほかの証明書アルゴリズムと RC2 暗号化は拒否されます。ポータブル署名に使用する RC2 PFX は、AES または 3DES を使用して再エクスポートしてください。Windows では、認証済みの RC2 PFX に代わりに SignTool が使用されます。パスワードが正しくない場合や PFX が破損している場合、このフォールバックが実行されることはありません。

証明書ストア、ハードウェアトークン、プロバイダーにより管理される鍵では、`thumbprint` を使用する Windows ネイティブの SignTool ワークフローが維持されます。MSI への署名にも、Windows 上の SignTool が必要です。エクスポート済み証明書を設定しても、これらのワークフローがポータブルになるわけではありません。

### サポートされる Windows 形式

| 形式 | 拡張子 | 備考 |
| --- | --- | --- |
| 実行ファイルとライブラリ | .exe, .dll | NSIS インストーラーを含むポータブル PE 署名 |
| インストーラー | .msi | Windows SignTool が必要 |
| アプリパッケージ | .msix, .appx | 個別パッケージへのポータブル署名。バンドルは非対応 |

## Linux パッケージへの署名

Linux パッケージ（DEB と RPM）には、PGP/GPG 鍵を使用して署名します。Windows や macOS のコード署名とは異なり、Linux パッケージへの署名は、コードが OS によって信頼されていることではなく、パッケージが信頼できる提供元から入手されたことを証明します。

エクスポートされた ASCII Armor 形式またはバイナリ形式の秘密鍵ファイルでは、Windows、macOS、Linux 上の Go ネイティブ署名機能が使用されます。既存の DEB または RPM への署名では、展開や再圧縮を行わずにパッケージの内容が保持され、GnuPG、`dpkg-sig`、`rpmsign` のインストールは不要です。Arch パッケージでは、同じキーローダーを使用して分離署名を書き込みます。GnuPG は、鍵の生成やエクスポート、および署名の独立した検証には引き続き有用です。鍵ファイルではなくキーリング内の ID を指定する従来の DEB/RPM 設定では、引き続き外部署名ツールが使用されます。

### 前提条件

- PGP 鍵ペア（Wails で生成可能）

### PGP 鍵の生成

最も簡単な方法はセットアップウィザードです。`wails3 setup` を実行し、署名手順の **Linux** タブを開いて、**新しい GPG 鍵を作成**を選択します。ウィザードでは、次の処理が行われます。

- GPG キーリングに RSA 4096 鍵を生成します（無人実行できる CI 向けの鍵にするには、パスフレーズを空欄のままにします）。
- その鍵を **`~/.wails/signing/<keyid>.asc`** にエクスポートします（ビルドが署名に使用するファイルです）。
- 鍵 ID とエクスポート先のパスの両方を `~/.config/wails/defaults.yaml` に記録し、署名時に自動的に使用されるようにします。

@note{type="note"}
鍵が **ID のみ**で設定されていた場合（たとえば、ウィザードが鍵を自動的にエクスポートするようになる前に作成された場合）、次回 Linux タブを開いたときに鍵がファイルへエクスポートされ、パスが入力されます。手動操作は不要です。ビルドでは鍵の*ファイル*を使用して署名するため、重要なのはパスです。

@end

`gpg` を使用して手動で行うこともできます。

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

推奨設定：RSA 4096 ビット、有効期限 1 年、強力なパスワードによる保護。

@note{type="caution"}
秘密鍵は厳重に管理してください。暗号化して保存し、安全にバックアップしてください。

@end

### 設定

`build/linux/Taskfile.yml` を編集し、署名用の変数を設定します。

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 変数 | 必須 | 説明 |
| --- | --- | --- |
| `PGP_KEY` | 上書き指定 | エクスポートされた PGP 秘密鍵ファイルへのパス（未設定の場合は、`wails3 setup` でグローバルに設定された鍵にフォールバック） |
| `SIGN_ROLE` | いいえ | DEB 署名ロール（デフォルト：builder） |

@note{type="note"}
PGP 鍵のパスワードは Taskfile ではなく、システムのキーチェーンに保存されます。設定するには `wails3 signing` を実行するか、CI で `WAILS_PGP_PASSWORD` 環境変数を設定します。

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
- `maint`：パッケージメンテナーによる署名
- `archive`：アーカイブメンテナーによる署名
- `builder`：パッケージビルダーによる署名（デフォルト）

### クロスプラットフォーム署名

Linux パッケージには、どのプラットフォームからでも署名できます。同じ Taskfile 設定とコマンドが Windows と macOS でも機能します。

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

Taskfile 変数を直接渡すこともできます。

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
クロスプラットフォーム署名に Linux ランナーを使用すると、個別の Windows ランナーが不要になり、CI/CD が簡素化されます。Apple キーチェーン内の ID と notary プロファイルには、macOS ランナーが必要です。エクスポートされた ID とチーム API 認証情報は、上記のポータブル Wake ワークフローをサポートします。

@end

## CLI リファレンス

### wails3 signing

プロジェクトの署名を設定するための対話型ウィザードです。

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

ウィザードでは、次の手順を案内します。

- **macOS**：Developer ID 証明書の選択、Notarization 認証情報の設定（`xcrun notarytool store-credentials`を呼び出します）。
- **Windows**：証明書ファイルまたはサムプリントの選択、パスワードとタイムスタンプサーバーの設定。
- **Linux**：既存の PGP 鍵の使用または新しい鍵の生成（`gpg`を呼び出します）、署名ロールの設定。

### wails3 entitlements

macOS のエンタイトルメントを設定する対話型ウィザード。

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**プリセット：**

- **開発**：JIT、デバッグ、ネットワーク用の`entitlements.dev.plist`を作成します
- **本番**：必要最小限のエンタイトルメントを含む`entitlements.plist`を作成します
- **両方**：両方のファイルを作成します（推奨）
- **App Store**：Mac App Store 向けのサンドボックス化されたエンタイトルメントを作成します
- **カスタム**：個々のエンタイトルメントと出力先ファイルを選択します

### wails3 sign

現在のプラットフォームまたは指定したプラットフォーム向けのバイナリとパッケージに署名します。これは、プラットフォームに対応する署名タスクを呼び出すラッパーです。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Taskfile の署名設定を使用する、対応する`<platform>:sign`タスクを実行します。

### wails3 tool sign

特定のファイルに直接署名するための低レベルコマンドです。Taskfile が内部で使用します。

```bash
wails3 tool sign [flags]
```

**共通フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--input` | 署名するファイルのパス |
| `--output` | 出力パス（省略可能。デフォルトでは元のファイルを直接更新） |
| `--verbose` | 詳細出力を有効にする |

**Windows/macOS フラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--certificate` | PKCS#12（.pfx/.p12）証明書のパス |
| `--password` | 証明書のパスワード |
| `--timestamp` | タイムスタンプサーバーの URL |

**macOS 固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--identity` | 署名 ID（アドホック署名には「-」を使用） |
| `--entitlements` | エンタイトルメント plist のパス |
| `--hardened-runtime` | Hardened Runtime を有効にする（デフォルト：true） |
| `--notarize` | Notarization のために提出する |
| `--keychain-profile` | Notarization 用のキーチェーンプロファイル |
| `--notary-key` | App Store Connect チーム API の秘密鍵`.p8`ファイル。Go Notary API クライアントを使用します |
| `--notary-key-id` | チーム API キー ID。`--notary-key`と併用する場合は必須です |
| `--notary-issuer-id` | チーム発行者 ID。`--notary-key`と併用する場合は必須です |

**Windows 固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--thumbprint` | Windows 証明書ストア内の証明書サムプリント |

**Linux 固有のフラグ：**\

| フラグ | 説明 |
| --- | --- |
| `--pgp-key` | PGP 秘密鍵のパス |
| `--pgp-password` | PGP 鍵のパスワード |
| `--role` | DEB 署名ロール（origin/maint/archive/builder） |

### 署名状態の確認（ネイティブツール）

`wails3 signing`は署名用の認証情報を設定します。署名状態を確認するには、ネイティブツールを直接使用してください。

| タスク | コマンド |
| --- | --- |
| macOS コード署名 ID の一覧表示 | `security find-identity -v -p codesigning` |
| Notarization 認証情報の保存 | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP 鍵ファイルの確認 | `gpg --show-keys <key.asc>` |
| PGP 鍵ペアの生成 | `gpg --full-generate-key` |
| 公開鍵のエクスポート | `gpg --armor --export <email>` |

## トラブルシューティング

### macOS の問題

**「Developer ID 証明書が見つかりません」**

- 証明書がキーチェーンにインストールされていることを確認してください
- `security find-identity -v -p codesigning`を使用して、有効期限が切れていないことを確認してください
- 「Apple Development」証明書だけでなく、「Developer ID Application」証明書があることを確認してください

**「Notarization に失敗しました」**

- Notarization ログを確認してください：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Hardened Runtime が有効になっていることを確認してください
- アプリに未署名のバイナリが含まれていないことを確認してください

**「コード署名に失敗しました」**

- キーチェーンがロック解除されていることを確認してください：`security unlock-keychain`
- アプリバンドルのファイル権限を確認してください

### Windows の問題

**「証明書が見つかりません」**

- 証明書のパスが正しいことを確認してください
- 証明書のパスワードを確認してください
- 証明書が有効であること（有効期限切れまたは失効していないこと）を確認します

**「タイムスタンプサーバーエラー」**

- 別のタイムスタンプサーバーを試します：
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux に関する問題

**「無効な PGP キー」**

- キーファイルが ASCII Armor 形式であることを確認します
- `gpg --show-keys <key.asc>`を使用して、キーの有効期限が切れていないことを確認します
- パスワードが正しいことを確認します

**「署名の検証に失敗しました」**

- 公開鍵が正しくインポートされていることを確認します
- 署名後にパッケージが変更されていないことを確認します

## その他のリソース

### 公式ドキュメント

- [Apple コード署名ガイド](https://developer.apple.com/support/code-signing/)
- [Apple 公証ドキュメント](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft コード署名](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian パッケージ署名](https://wiki.debian.org/SecureApt)
- [RPM パッケージ署名](https://rpm-software-management.github.io/rpm/manual/signatures.html)
