---
title: "セットアップ"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="試験的機能"}
セットアップウィザードは新しい機能であり、主に Linux でテストされています。問題が発生した場合は、[報告](https://github.com/wailsapp/wails/issues/4904)し、代わりに[手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従ってください。

@end

## クイックスタート

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

ウィザードがブラウザーで開き、依存関係の確認、プロジェクトのグローバルデフォルト、オプションのクロスプラットフォームツール、署名認証情報の設定を案内します。これらの設定は複数のプロジェクトで再利用できます。

これで、最初のプロジェクトを作成する準備が整いました。

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## セットアップの機能

- **依存関係の確認** - Go、npm、プラットフォーム用ツールを検証します
- **既定値の設定** - 作成者情報、バンドル ID のプレフィックス、優先テンプレートを設定します
- **クロスプラットフォームビルド** - 任意のホストからビルドするための Docker を必要に応じて設定します
- **コード署名** - macOS、Windows、Linux 向けの署名を必要に応じて設定します

設定は`$XDG_CONFIG_HOME/wails/defaults.yaml`に保存され、使用できない場合は`~/.config/wails/defaults.yaml`に保存されます。この設定は`wails3 init`によって使用されます。これは、プロジェクトの`build/config.yml`とは別のものです。

## Wake プロジェクトを作成する

Wake の実験機能を有効にすると、`wails3 init -ui` はプロジェクト固有のビルド設定を収集し、作成に成功した時点で `wails.yaml` に保存します。デスクトッププロジェクトの作成時にビルドバックエンド、ターゲット、パッケージ形式を選択するには、[Wake 入門チュートリアル](/experimental/wake/)に従ってください。`wails3 setup` は引き続きグローバル設定を管理します。

既存の Wake プロジェクトは`wails.yaml`で設定します。プロファイルやパッケージ設定を変更するには、[YAML ビルドガイド](/experimental/yaml-builds/)を参照してください。既存の Taskfile プロジェクトで Wake を使用するには、事前に[レビュー済みの一方向移行](/experimental/wake/#migrate-an-existing-project)が必要です。セットアップでは Taskfile の変換や実行は行いません。

## 署名設定を再利用する

署名付き Wake ビルドでは、プロジェクトにない署名設定をユーザーレベルの既定値から再利用します。これには、プロジェクト外の証明書ファイルも含まれます。プロジェクトで明示した設定が優先されます。秘密鍵、パスワード、グローバルな認証情報のパスは`wails.yaml`にコピーされません。

グローバルに設定した Windows および Linux の鍵については、既存の認証情報ストア、または`WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`を使用してパスワードを指定します。暗号化された Apple の署名用 ID については、そのパスワードを格納した環境変数の名前を`darwin.signing.credential`に設定します。

## サブコマンド

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 問題が発生した場合

1. 問題を診断するには、`wails3 doctor`を実行します
2. [手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従います
3. `wails3 doctor`の出力を添えて、[問題を報告](https://github.com/wailsapp/wails/issues/4904)します
