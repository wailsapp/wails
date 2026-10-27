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

ウィザードがブラウザーで開き、依存関係の確認、プロジェクトのデフォルト設定、任意のクロスプラットフォームビルド設定を順に案内します。

これで、最初のプロジェクトを作成する準備が整いました。

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 実行される処理

- **依存関係を確認** - Go、npm、プラットフォーム用ツールを検証します
- **デフォルトを設定** - 作成者情報、バンドル ID のプレフィックス、優先テンプレートを設定します
- **クロスプラットフォームビルド** - 任意のホストからビルドするための Docker を必要に応じて設定します
- **コード署名** - macOS、Windows、Linux 向けのコード署名を必要に応じて設定します

設定は `$XDG_CONFIG_HOME/wails/defaults.yaml` に保存され、利用できない場合は `~/.config/wails/defaults.yaml` が使用されます。この設定は `wails3 init` によって使用されます。これは、プロジェクトの `build/config.yml` とは別のものです。

## サブコマンド

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 問題が発生した場合

1. 問題を診断するには、`wails3 doctor` を実行します
2. [手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従います
3. `wails3 doctor` の出力を添えて、[問題を報告](https://github.com/wailsapp/wails/issues/4904)します
