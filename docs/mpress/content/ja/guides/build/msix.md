---
title: "MSIX パッケージ化"
description: "Wails v3 アプリケーションを MSIX パッケージとしてパッケージ化する"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX は、Windows の最新のアプリケーションパッケージ形式です。Wails では、Windows ビルドの一環として MSIX パッケージを生成できます。

MSIX パッケージ化の手順については、[Windows パッケージ化](/guides/build/windows/#msix-package)ガイドを参照してください。

## CLI による直接パッケージ化

CI を含め、MSIX ツールは Windows 上で実行してください。既定のバックエンドは `MakeAppx.exe` を使用し、署名には `signtool.exe` も必要です。どちらも Windows SDK に含まれています。インストールヘルパーは Microsoft Store と、必要に応じて SDK のダウンロードページを開きます。パッケージ化の前にインストールを完了してください。

識別情報を `build/config.yml` に設定し、実行ファイルをビルドしてパッケージ化します。

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

直接実行するコマンドは現在のディレクトリに `MyApp.msix` を出力します。一方、[Windows パッケージ化タスク](/guides/build/windows/#msix-package) は独自の出力パスを指定します。

## CLI オプション

| オプション | 意味 |
| --- | --- |
| `--config` | 設定ファイル。既定値は `build/config.yml`。 |
| `--executable`, `--name` | 既存の実行ファイルとパッケージ内でのファイル名。どちらも必須。 |
| `--out` | 出力ファイル。既定値は `<ProductName>.msix`。 |
| `--arch` | パッケージのアーキテクチャ：`x64`（既定値）、`x86`、`arm`、`arm64`、`x86a64`、または `neutral`。Go の別名 `amd64` と `386` も使用可能。実行ファイルのアーキテクチャに合わせてください。 |
| `--publisher` | 発行元の識別情報。既定値は `CN=<companyName>`。 |
| `--cert`, `--cert-password` | 署名用 PFX 証明書のパスとパスワード。 |
| `--use-makeappx` | Windows SDK の既定のパッケージ化ツールを使用。 |
| `--use-msix-tool` | `MsixPackagingTool.exe` を明示的に選択。`PATH` に存在する必要があります。 |

## 署名と CI

Store 以外で配布する場合は、配布先のマシンで信頼される証明書で署名してください。証明書の Subject は `--publisher` と完全に一致する必要があります。MakeAppx バックエンドは `--cert` が指定されると SHA256 を使用して SignTool を呼び出します。[Microsoft の署名ガイド](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide)を参照してください。

この Windows ワークフローのステップは、Wails と SDK がインストール済みで、先行するステップが PFX ファイルを `CERT_PATH` に安全に配置済みであることを前提としています。パスを格納したシークレットだけでは証明書はアップロードされません。

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## ファイルの関連付けとアセット

拡張子を `build/config.yml` に先頭のドットを除いて追加します。生成されるマニフェストがドットを付加します。

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

実行時のファイルオープンは[ファイルの関連付け](/guides/file-associations/)に従って処理してください。現在、MakeAppx バックエンドは実行ファイルのみをコピーし、透明なプレースホルダー画像を生成します。プロジェクトの `Assets/` ファイルの取り込みや `iconName` のアイコン変換は行いません。独自のブランド画像や追加 DLL にはカスタムのパッケージ化ワークフローを使用してください。

| 生成されるアセット | サイズ（ピクセル） |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` はファイルの関連付けが設定されている場合にのみ生成されます。これらのファイルはパッケージの `Assets/` ディレクトリに格納されます。

## Store への提出とトラブルシューティング

[Partner Center](https://partner.microsoft.com/dashboard) でアプリケーションを予約し、提出の準備ではそこで指定されるパッケージの識別情報と発行元を使用してください。Store は提出時に MSIX パッケージへ署名するため、この配布方法では署名証明書を購入する必要はありません。[Microsoft のパッケージ要件](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements)を参照してください。

ツールの `MakeAppx.exe` または `signtool.exe` が見つからない場合は、Windows SDK をインストールまたは修復してください。Wails は `PATH` と SDK の標準の場所を検索します。署名に失敗した場合は、証明書の Subject、有効期限、配布先のマシンでの信頼を確認してください。[MSIX のトラブルシューティング](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide)を参照してください。
