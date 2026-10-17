---
title: "Windows 向けパッケージ化"
description: "Wails アプリケーションを Windows で配布するためにパッケージ化します"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## NSIS インストーラー

デフォルトのパッケージ形式では、NSIS インストーラーが作成されます：

```bash
wails3 package GOOS=windows
```

これにより `wails3 task windows:package` が実行され、次の処理が行われます：

1. アプリケーションのビルド
2. WebView2 ブートストラッパーの生成
3. NSIS インストーラーの作成

出力：`build/windows/nsis/<AppName>-installer.exe`

### MSIX パッケージ

Microsoft Store での配布または最新の Windows 環境への展開には、次を使用します：

```bash
wails3 package GOOS=windows FORMAT=msix
```

出力：`bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX には、`makeappx.exe`（Windows SDK）またはスタンドアロンの MSIX ツールが必要です。Windows Taskfile では、インストールタスクが `wails3 task install:msix:tools` として提供されています。

@end

## インストーラーのカスタマイズ

NSIS の設定は `build/windows/nsis/project.nsi` にあります。このファイルを編集すると、次の項目をカスタマイズできます：

- インストーラーの UI とブランディング
- インストール先ディレクトリ
- スタートメニューとデスクトップのショートカット
- ファイルの関連付け
- 使用許諾契約

アプリケーションのメタデータは `build/windows/info.json` から取得されます：

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## コード署名

SmartScreen の警告を回避するには、実行ファイルとインストーラーに署名します：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

署名は `build/windows/Taskfile.yml` で設定します：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

証明書のパスワードは安全に保管してください：

```bash
wails3 setup signing
```

詳細については、[アプリケーションへの署名](/guides/build/signing/)を参照してください。

## ARM 向けのビルド

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## トラブルシューティング

### makensis が見つからない

NSIS をインストールします：

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### SmartScreen の警告

実行ファイルが署名されていません。上記の[コード署名](#heading-1)を参照してください。

### WebView2 が見つからない

インストーラーには、必要に応じてランタイムをダウンロードする WebView2 ブートストラッパーが含まれています。オフラインインストールが必要な場合は、Microsoft から Evergreen Standalone Installer をダウンロードしてください。
