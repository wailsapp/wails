---
title: "インストール"
description: "Wails をインストールして開発環境をセットアップする"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## サポート対象プラットフォーム

- Windows AMD64/ARM64
- macOS 10.15 以降 AMD64（macOS 10.13 以降にデプロイ可能）
- macOS 11.0 以降 ARM64
- Ubuntu 24.04 AMD64/ARM64（その他の Linux でも動作する可能性があります）

## 依存関係

Wails をインストールする前に、いくつかの共通の依存関係を用意する必要があります。

@note{type="tip"}
Wails CLI のインストール後に `wails3 setup` を実行すると、これらの依存関係を自動的に確認し、インストールを支援できます。

@end

@tabs
[Go（1.24 以上）]
[Go ダウンロードページ](https://go.dev/dl/)から Go をダウンロードしてください。

必ず公式の[Go インストール手順](https://go.dev/doc/install)に従ってください。また、`PATH` 環境変数に `~/go/bin` ディレクトリへのパスが含まれていることも確認してください。ターミナルを再起動し、次の項目を確認します。

- Go が正しくインストールされていることを確認：`go version`
- `~/go/bin` が PATH 環境変数に含まれていることを確認
  - Mac / Linux：`echo $PATH | grep go/bin`
  - Windows：`$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm（任意）]
Wails 自体は npm を必要としませんが、同梱されているほとんどのテンプレートでは npm が必要です。

[Node ダウンロードページ](https://nodejs.org/en/download/)から最新の Node インストーラーをダウンロードしてください。通常は最新リリースを対象にテストしているため、最新リリースの使用を推奨します。

確認するには `npm --version` を実行します。

@note{type="info"}
npm 以外のパッケージマネージャーを使用したい場合は、自由に使用できます。そのパッケージマネージャーを使用するように、プロジェクトの Taskfile を更新する必要があります。

@end

@end

## プラットフォーム固有の依存関係

プラットフォーム固有の依存関係もインストールする必要があります。

@tabs{sync-key="platform"}
[Mac]
Wails では、Xcode コマンドラインツールがインストールされている必要があります。次のコマンドを実行してインストールできます。

```sh
xcode-select --install
```

[Windows]
Wails では、[WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) がインストールされている必要があります。ほぼすべての Windows 環境には、すでにインストールされています。`wails doctor` コマンドで確認できます。

[Linux]
Linux では、標準の `gcc` ビルドツールに加えて、`gtk4` と `webkitgtk-6.0` が必要です。インストール後に <code>wails3 doctor</code> を実行すると、依存関係のインストール方法が表示されます。従来の GTK3 / WebKit2GTK 4.1 スタックも、v3.1 までは `-tags gtk3` で引き続き利用できます（[Linux パッケージング - 従来の GTK3 サポート](/guides/build/linux/#legacy-gtk3-support)を参照）。お使いのディストリビューションまたはパッケージマネージャーがサポートされていない場合は、Discord でお知らせください。

@end

## インストール

Go Modules を使用して Wails CLI をインストールするには、次のコマンドを実行します。

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

最新の開発版をインストールする場合は、次のコマンドを実行します。

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

開発版を使用する場合、生成されるすべてのプロジェクトでは、開発版の Wails が使用されるように Go の [replace](https://go.dev/ref/mod#go-mod-file-replace) ディレクティブが使用されます。

## 次のステップ

CLI のインストール後、セットアップウィザードを実行して開発環境を構成します。

```shell
wails3 setup
```

@note{type="caution" title="試験的機能"}
セットアップウィザードは新しい機能であり、主に Linux でテストされています。問題が発生した場合は、[問題を報告](https://github.com/wailsapp/wails/issues/4904)し、以下の手動による依存関係のインストール手順を使用してください。

@end

セットアップウィザードでは、次の処理を行います。

- プラットフォーム固有の依存関係を確認し、インストールを支援する
- プロジェクトのデフォルト設定（作成者情報、バンドル ID のプレフィックス）を構成する
- 必要に応じて、クロスプラットフォームビルド用に Docker をセットアップする
- 必要に応じてコード署名を構成する

詳細については、[セットアップガイド](/getting-started/setup/)を参照してください。

## 依存関係の手動インストール

依存関係を手動でインストールする場合、またはお使いのシステムでセットアップウィザードが動作しない場合は、上記のプラットフォーム固有の手順に従ってから、次のコマンドを実行します。

```shell
wails3 doctor
```

必要な依存関係が正しくインストールされているかを確認し、不足しているものを通知します。

## `wails3` コマンドが見つからない場合

`wails3` コマンドが見つからないとシステムから報告された場合は、次の項目を確認してください。

- 上記の<strong>Go インストールガイド</strong>に正しく従っていること、および `go/bin` ディレクトリが `PATH` 環境変数に含まれていることを確認してください。
- 新しい `PATH` 環境変数を反映するため、現在開いているターミナルを閉じてから再度開いてください。
