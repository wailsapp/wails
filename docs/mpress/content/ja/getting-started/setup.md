---
title: "セットアップ"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="試験運用中"}
セットアップウィザードは新しい機能で、主に Linux でテストされています。問題が発生した場合は、[報告](https://github.com/wailsapp/wails/issues/4904)したうえで、代わりに[手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従ってください。

@end

## クイックスタート

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

ウィザードがブラウザーで開き、依存関係の確認、プロジェクトのデフォルト設定、オプションのクロスプラットフォームビルド設定を案内します。

これで、最初のプロジェクトを作成する準備が整いました。

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 実行される内容

- **依存関係を確認** — Go、npm、プラットフォーム固有のツールを検証します
- **デフォルト値を設定** — 作者情報、バンドル ID のプレフィックス、優先テンプレートを設定します
- **クロスプラットフォームビルド** — 任意のホストからビルドするための Docker 環境を任意で設定します
- **コード署名** — macOS、Windows、Linux 向けに任意で設定します

設定は `~/.config/wails/config.yaml` に保存され、`wails3 init` で使用されます。

## サブコマンド

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## 問題が発生した場合

1. 問題を診断するには、`wails3 doctor` を実行してください
2. [手動インストール手順](/getting-started/installation/#platform-specific-dependencies)に従ってください
3. `wails3 doctor` の出力を添えて、[問題を報告](https://github.com/wailsapp/wails/issues/4904)してください
