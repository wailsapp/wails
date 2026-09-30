---
title: "はじめに"
description: "Wails v3 へのコントリビューションを始める方法"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## コントリビューターの皆さん、ようこそ！

Wails へのコントリビューションに関心をお寄せいただき、ありがとうございます。このガイドでは、初めてのコントリビューションを行う方法を説明します。

## 前提条件

始める前に、以下が準備できていることを確認してください。

- **Go 1.25+** がインストール済みであること（[ダウンロード](https://go.dev/dl/)）
- **Node.js 20+** および **npm**（[ダウンロード](https://nodejs.org/)）
- **Git** が GitHub アカウント用に設定済みであること
- Go および JavaScript/TypeScript の基礎知識

### プラットフォーム固有の要件

**macOS：**

- Xcode Command Line Tools：`xcode-select --install`

**Windows：**

- MSYS2 または同様の Unix 系環境を推奨
- WebView2 ランタイム（通常は Windows 11 にプリインストール済み）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（デフォルトの GTK4 スタック）
- 次のコマンドでインストール：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 従来の `-tags gtk3` ビルドパスを使用する場合は、`libgtk-3-dev` と `libwebkit2gtk-4.1-dev` もインストールしてください

## コントリビューション手順の概要

一般的なコントリビューションのワークフローは、次の手順で進みます。

1. **フォークとクローン** - Wails リポジトリの自分用のコピーを作成します
2. **セットアップ** - Wails CLI をビルドし、環境を検証します
3. **ブランチ作成** - 変更用の機能ブランチを作成します
4. **開発** - コーディング規約に従って変更を加えます
5. **テスト** - テストを実行し、すべてが正しく動作することを確認します
6. **コミット** - 明確で Conventional Commits に準拠したメッセージを付けてコミットします
7. **提出** - レビュー用のプルリクエストを作成します
8. **修正の反復** - フィードバックに対応し、必要な調整を行います
9. **マージ** - 承認されると、変更が Wails に取り込まれます！

## ステップ別ガイド

コントリビューションの種類を選択してください。

@tabs
[バグ修正]
@steps
### バグを探す、または報告する
- そのバグが [GitHub Issues](https://github.com/wailsapp/wails/issues) ですでに報告されていないか確認してください
- 報告されていない場合は、再現手順を記載した新しい Issue を作成してください
- 作業を開始する前に確認を待ってください

### フォークとクローン
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) でリポジトリをフォークしてください

フォークをクローンします。

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### ビルドと検証
Wails をビルドし、バグを再現できることを確認します。

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### バグ修正用ブランチの作成
修正用のブランチを作成します。

```bash
git checkout -b fix/issue-123-window-crash
```

### バグの修正
- バグの修正に必要な最小限の変更を加えてください
- 無関係なコードをリファクタリングしないでください
- リグレッションを防ぐため、テストを追加または更新してください

```bash
# Make your changes
# Add tests in *_test.go files
```

### 修正のテスト
テストを実行し、修正が機能することを確認します。

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 修正のコミット
明確なメッセージを付けてコミットします。

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### プルリクエストの提出
プッシュして PR を作成します。

```bash
git push origin fix/issue-123-window-crash
```

PR の説明には、以下を記載してください。

- バグと根本原因を説明する
- 修正内容を説明する
- Issue を参照する："Fixes #123"
- 修正前と修正後の動作を記載する

### フィードバックへの対応
レビューコメントに対応し、必要に応じて PR を更新してください。

@end

[WEP（機能拡張）]
@steps
### WEP の作成
- [WEP（Wails Enhancement Proposal）のプロセス](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)を確認してください
- WEP テンプレートを `v3/wep/proposals/<name>/proposal.md` にコピーしてください
- タイトルを `[WEP] <title>` とし、WEP のみを含むドラフト PR を作成してください
- 実装を開始する前に、メンテナーの決定を待ってください

### フォークとクローン
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) でリポジトリをフォークしてください

フォークをクローンします。

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 開発環境をセットアップする
Wails をビルドし、環境を検証します。

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 機能ブランチを作成する
内容が分かりやすい名前のブランチを作成します。

```bash
git checkout -b feat/window-transparency-support
```

### 機能を実装する
- [コーディング規約](/contributing/standards/)に従う
- 変更を対象の機能に限定する
- 明確で、適切に文書化されたコードを書く
- 包括的なテストを追加する

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 十分にテストする
機能をテストします。

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### 機能を文書化する
- すべての公開 API に docstring を追加する
- `/docs/mpress/content/`内の関連ドキュメントを更新する
- 該当する場合は例を追加する

### 規約に従ってコミットする
Conventional Commits を使用します。

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### プルリクエストを提出する
プッシュして PR を作成します。

```bash
git push origin feat/window-transparency-support
```

PR には次の内容を含めます。

- 機能とユースケースを説明する
- 例またはスクリーンショットを提示する
- 破壊的変更があればすべて列挙する
- 承認済みの WEP PR を参照する

### レビューに基づいて修正を重ねる
メンテナーから変更を求められる場合があります。辛抱強く、協力的に対応してください。

@end

[ドキュメント]
修正 PR は、事前に issue を作成しなくても歓迎します。ドキュメントのみの修正では、失敗するコードテストは必要ありません。M-Press のインストール、ソースパス、プレビュー、検証、PR の各手順については、[ドキュメントを修正する](/contributing/documentation/)に従ってください。

@end

## 取り組む issue を探す

- [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)ラベルを探す
- [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)issue を確認する
- [未解決の issue](https://github.com/wailsapp/wails/issues)を確認し、担当への割り当てを依頼する

## サポートを受ける

- **Discord：**[Wails Discord](https://discord.gg/JDdSxwjhGf) に参加する
- **Discussions：**[GitHub Discussions](https://github.com/wailsapp/wails/discussions) に投稿する
- <strong>Issues：</strong>再現可能なバグについては issue を作成する。質問には Discussions、機能強化には WEP PR を使用する

## 行動規範

相手を尊重し、建設的かつ友好的に接してください。私たちは、優れたソフトウェアを共に作ることに注力する、親しみやすいコミュニティを築いています。

## 次のステップ

- [開発環境](/contributing/setup/)をセットアップする
- [コーディング規約](/contributing/standards/)を確認する
- [技術ドキュメント](/contributing/overview/)を参照する
