---
title: "CLI リファレンス"
description: "Wails CLI コマンドの完全なリファレンス"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## 概要

Wails CLI（`wails3`）は、Wails 3 アプリケーションの作成、開発、ビルド、署名、パッケージ化、検査を行うためのコマンドラインエントリーポイントです。ビルドのオーケストレーションの大部分は、プロジェクトごとの Taskfile（プロジェクト内の `build/` 以下）に委譲されます。多くの `wails3` コマンドは、特定のタスクを呼び出す薄いラッパーです。

各コマンドの最新のヘルプを確認するには、次を実行します。

```bash
wails3 --help
wails3 <command> --help
```

## プロジェクトのライフサイクル

| コマンド | 説明 |
| --- | --- |
| `wails3 init` | テンプレートから新しいプロジェクトを作成します。フラグ：`-n`（プロジェクト名）、`-t`（テンプレート、デフォルトは `vanilla`）、`-p`（Go パッケージ名、デフォルトは `main`）、`-d`（プロジェクトディレクトリ、デフォルトは `.`）、`-q`（出力を抑制）、`-l`（テンプレート一覧を表示）、`-mod`（Go モジュールパス）、`--git`（Git リポジトリ URL）、`--skipgomodtidy`、`-s`（リモートテンプレートの警告を省略）、`--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`。 |
| `wails3 dev` | フロントエンドのホットリロードを有効にして、アプリケーションを開発モードで実行します。フラグ：`--config`（デフォルトは `./build/config.yml`）、`--port`（Vite 開発ポート）、`-s`（HTTPS を有効化）。 |
| `wails3 build` | プロジェクトをビルドします。Taskfile の `build` タスクを呼び出す薄いラッパーです。フラグ：`--tags`（`EXTRA_TAGS=` として転送）、`--obfuscated`（Garble を使用してビルド。[難読化ビルド](/guides/build/obfuscation/)を参照）、`--garbleargs`（`build` サブコマンドの前で `garble` に転送する追加フラグ）。 |
| `wails3 package` | プラットフォーム固有の `package` Taskfile タスクを実行します。 |
| `wails3 task [name]` | 任意の Taskfile タスクを実行します。名前を指定しない場合、`--list` は登録済みのすべてのタスクを表示します。 |
| `wails3 mcp` | プロジェクトの MCP サーバーを起動します。エージェントから起動されたプロセスでは stdio を、対話型ターミナルでの使用時にはループバックの Streamable HTTP を自動的に使用します。 |
| `wails3 doctor` | 環境の診断レポートを出力します。 |
| `wails3 doctor-ng` | `doctor` の新しい TUI 版です。 |
| `wails3 version` | CLI のバージョンを出力します。 |
| `wails3 releasenotes` | 最近のリリースノートを出力します。 |
| `wails3 docs` | ブラウザーでドキュメントサイトを開きます。 |
| `wails3 sponsor` | スポンサー用ページを開きます。 |

## 生成

`wails3 generate <subcommand>`：

| サブコマンド | 説明 |
| --- | --- |
| `generate bindings` | Go からフロントエンドへのバインディングを生成します。フラグ：`-d`（出力ディレクトリ）、`-models`、`-index`、`-ts`、`-i`（インターフェース）、`-b`（バンドル）、`-names`（`Call.ByName` を出力）、`-noevents`、`-noindex`、`-dry`、`-silent`、`-v`、`-clean`（デフォルトは `true`）、`-f`、`-obfuscated`（Garble ビルド用の安定したバインディング ID を含む `wails_obfuscated.gen.go` を生成。[難読化ビルド](/guides/build/obfuscation/)を参照）、`-obfuscated-output`（生成ファイルの格納先ディレクトリ。デフォルトは main パッケージのディレクトリ）。パッケージパターン（例：`./...`）を指定できます。何も指定しない場合は、現在のディレクトリを使用します。 |
| `generate icons` | ソース PNG をプラットフォーム用のアイコン形式に変換します。フラグ：`-input`、`-windowsfilename`、`-macfilename`、`-iconcomposerinput`、`-macassetdir`。 |
| `generate build-assets` | `build/config.yml` から `build/` ディレクトリの内容（Taskfile スニペット、NSIS ファイル、`Info.plist`、`.desktop` テンプレートなど）を生成します。 |
| `generate runtime` | WebView に提供されるビルド済みの `/wails/runtime.js` を再生成します。 |
| `generate syso` | Windows 用の `.syso` リソースファイル（アイコン、マニフェスト、バージョン情報）を生成します。 |
| `generate webview2bootstrapper` | Windows 用の WebView2 ブートストラップインストーラーを生成します。 |
| `generate constants` | Go のイベント型から JS のイベント名定数を生成します。 |
| `generate template` | 新しいプロジェクトテンプレートのひな形を生成します。 |
| `generate .desktop` | Linux 用の `.desktop` ファイル（AppImage/DEB/RPM で使用）を生成します。 |
| `generate appimage` | AppImage のビルドディレクトリを生成します。 |

## 更新

`wails3 update <subcommand>`：

| サブコマンド | 説明 |
| --- | --- |
| `update build-assets` | `build/config.yml` から `build/` ディレクトリを更新します（可能な限りユーザーによる編集を保持します）。 |
| `update cli` | `wails3` バイナリ自体を更新します。 |

## コード署名とパッケージ化

| コマンド | 説明 |
| --- | --- |
| `wails3 setup signing` | `build/` で検出されたプラットフォーム向けの署名を構成する対話型ウィザードです。フラグ：`--platform`（複数回指定可能。デフォルトではビルドディレクトリから自動検出）。 |
| `wails3 setup entitlements` | macOS のエンタイトルメントを設定する対話型ウィザード。フラグ：`--output`（パス、デフォルト：`build/darwin/entitlements.plist`）。 |
| `wails3 sign [GOOS=…]` | 現在の OS（または `GOOS` で指定した OS）向けのプラットフォーム固有の `*:sign` Taskfile タスクを実行するラッパー。 |
| `wails3 tool sign` | 低レベルの直接署名エントリーポイント。フラグ：`--input`、`--output`、`--verbose`、`--certificate`、`--password`、`--thumbprint`、`--timestamp`、`--identity`、`--entitlements`、`--hardened-runtime`、`--notarize`、`--keychain-profile`、`--pgp-key`、`--pgp-password`、`--role`。 |

`wails3 signing` サブコマンドは<strong>ありません</strong>。キーチェーンの資格情報には `xcrun notarytool store-credentials` を、PGP キーには `gpg` を直接使用してください（`wails3 setup signing` ウィザードでは両方を自動化できます）。

## ツール

`wails3 tool <subcommand>`：

| サブコマンド | 説明 |
| --- | --- |
| `tool checkport` | TCP ポートが開いているか確認します（Vite の起動を待機する場合に便利です）。 |
| `tool watcher` | 監視対象ファイルが変更されるたびにコマンドを実行します。 |
| `tool cp` | クロスプラットフォームでファイルをコピーします。 |
| `tool buildinfo` | バイナリに埋め込まれた Go のビルド情報を出力します。 |
| `tool package` | `build/linux/nfpm` から Linux パッケージ（`deb`、`rpm`、`archlinux`）をビルドします。 |
| `tool version` | プロジェクトのセマンティックバージョンを更新します。 |
| `tool lipo` | 複数の macOS アーキテクチャ向けバイナリを結合し、ユニバーサルバイナリを作成します。 |
| `tool capabilities` | システムを調査し、GTK3/GTK4 および WebKit が利用可能か確認します。 |
| `tool sign` | （[コード署名とパッケージ化](#heading-4)を参照してください。） |

## サービス

`wails3 service <subcommand>`：

| サブコマンド | 説明 |
| --- | --- |
| `service init` | 新しいサービスパッケージのひな形を生成します。 |

## iOS

`wails3 ios <subcommand>`：

| サブコマンド | 説明 |
| --- | --- |
| `ios overlay:gen` | iOS ブリッジシム用の Go オーバーレイを生成します。 |
| `ios xcode:gen` | 出力ディレクトリに Xcode プロジェクトを生成します。 |

## ビルド出力パス

- ネイティブバイナリは `bin/<APP_NAME>`（Windows では `bin/<APP_NAME>.exe`）に出力されます。`build/bin/`は存在しません。
- パッケージ化された出力（`.app`、`.dmg`、NSIS インストーラー、MSIX、DEB/RPM/AppImage）も同様に `bin/` に出力されます（または、該当する Taskfile タスクによって作成されたプラットフォーム固有のサブディレクトリに出力されます）。

## グローバルフラグ

| フラグ | 適用対象 | 説明 |
| --- | --- | --- |
| `--no-colour` | すべてのコマンド | CLI 出力で ANSI カラーを無効にします。 |

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
