---
title: "セットアップ"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="実験的機能"}
セットアップウィザードは新しい機能であり、主に Linux でテストされています。問題が発生した場合は、[報告](https://github.com/wailsapp/wails/issues/4904)し、代わりに[手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従ってください。

@end

## クイックスタート

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

ウィザードはブラウザで開き、依存関係の確認、プロジェクトのデフォルト設定、任意のクロスプラットフォームビルド設定を案内します。上記の手順ではデフォルトのセットアップフローを使用しています。実験的な Wake プロジェクトのフローについては、以下で説明します。

これで、最初のプロジェクトを作成する準備が整いました。

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 実行される処理

- **依存関係を確認** - Go、npm、およびプラットフォーム用ツールを検証します
- **デフォルトを設定** - 作成者情報、バンドル ID のプレフィックス、優先テンプレートを設定します
- **クロスプラットフォームビルド** - 任意のホストからビルドするための Docker を必要に応じて設定します
- **コード署名** - macOS、Windows、および Linux 向けの署名を必要に応じて設定します

設定は `$XDG_CONFIG_HOME/wails/defaults.yaml` に保存され、利用できない場合は `~/.config/wails/defaults.yaml` にフォールバックし、`wails3 init` によって使用されます。これはプロジェクトの `build/config.yml` とは別のものです。

## Wake プロジェクトを設定する

最初のアプリケーションを一通り作成するには、[Wake 入門チュートリアル](/experimental/wake/)に従ってください。
このチュートリアルでは、プロジェクトの作成と各セットアップ画面を順に説明し、この実験で想定される動作と既存のビルドからの移行について解説します。

プロジェクトディレクトリから、[Wake の実験的機能を有効にして](/experimental/yaml-builds/#enable-yaml-builds)セットアップを実行します。

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

PowerShell の場合：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

プロジェクト向けのフローが開きます。

1. **プロジェクト** — ネイティブマニフェストが存在しない場合に、プロジェクトのメタデータを入力します。
2. **アプリケーション** — GUI、CLI、サーバーのいずれか、ビルドバックエンド、またはビルド済み実行ファイルを選択します。Pure Go は、明示的に選択する実験的なオプションです。
3. **ターゲット** — オペレーティングシステムとアーキテクチャを選択します。iOS と Android は引き続き実験的機能です。これらの署名設定は YAML で構成してください。
4. **パッケージ** — パッケージ形式と、署名するかどうかを選択します。
5. **ツール** — 選択したビルドおよびパッケージング方法の要件を確認します。
6. **署名** — 署名を選択した場合に、既存の署名設定を使用して構成します。
7. **保存** — プロジェクト設定を書き込み、実行するコマンドを表示します。

必要なツールは選択内容によって異なります。Pure Go によるデスクトップ向けコンパイルには、Docker も C コンパイラーも必要ありません。CGo のクロスビルドには、互換性のあるコンパイラーとターゲットライブラリ、またはサポート対象のコンテナツールチェーンが必要です。NSIS には引き続き `makensis` が必要です。ホストとは異なる環境で GUI AppImage をビルドするには、提供された Linux ランタイムツリーが必要です。エクスポートされた署名 ID では Wails の Go 署名機能を使用できますが、ネイティブのキーチェーンまたは証明書ストアの ID では、引き続きプラットフォーム用ツールを使用します。[パッケージングと署名の要件](/experimental/yaml-builds/#package-and-sign-for-distribution)を参照してください。

ビルド済み実行ファイルを選択すると、Go のコンパイル、フロントエンドのインストールとビルド、バインディング、および UPX はスキップされます。それでも、選択したパッケージング／署名ツールと、有効な入力バイナリは必要です。[ビルド済み入力のパス](/experimental/yaml-builds/#use-prebuilt-desktop-executables)を参照してください。この選択は保存後に表示されるコマンドに適用されるものであり、永続的な YAML ビルド設定ではありません。

Save を実行すると、`wails.yaml`が書き込まれ、アプリケーションの種類、フロントエンドの有効化設定、選択したバックエンド設定が更新され、専用の`profiles.setup`ターゲット／パッケージ／署名の選択内容が置き換えられます。その他のプロファイル、フロントエンドコマンド、コンパニオン実行ファイル、フック、および無関係なプロジェクト設定は保持されます。バックエンド設定はプロジェクト全体に適用されるため、その設定を継承する他のプロファイルでも、選択したバックエンドが使用されます。新しいプロジェクトメタデータは、Save を実行するまでウィザード内にのみ保持されます。不足しているツールは、ビルド前にインストールできるよう報告されます。有効な構成を保存しても、ツールのインストールやビルドの実行は行われません。

署名付きビルドでは、プロジェクト外の証明書ファイルを含め、ユーザーレベルのデフォルト設定から不足している署名設定を再利用します。プロジェクトで明示した設定が優先されます。セットアップでは、秘密鍵、パスワード、またはグローバルな認証情報のパスを `wails.yaml` にコピーしません。グローバルに設定された Windows および Linux の鍵については、既存の認証情報ストア、`WAILS_WINDOWS_CERT_PASSWORD`、または `WAILS_PGP_PASSWORD` を使用してパスワードを指定してください。暗号化された Apple ID については、そのパスワードを格納した環境変数の名前を `darwin.signing.credential` に設定してください。

保存したプロファイルは `wails3 build setup` で実行します。ビルド済みオプションを選択した場合は `--skip-build` が追加されます。

表示されたコマンドを実行する際は、Wake を有効なままにしてください。フラグが設定されていない場合、有効な `wails.yaml` が Taskfile にフォールバックすることはありません。セットアップでは、既存の Taskfile やレガシーメタデータは変換されません。該当するプロジェクトを Wake で設定する前に、[レビュー済みの移行ワークフロー](/experimental/yaml-builds/#migrate-an-existing-project)を使用してください。

## サブコマンド

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 問題が発生した場合

1. 問題を診断するには `wails3 doctor` を実行します
2. [手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従います
3. `wails3 doctor` の出力を添えて、[問題を報告](https://github.com/wailsapp/wails/issues/4904)します
