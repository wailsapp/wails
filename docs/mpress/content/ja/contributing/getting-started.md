---
title: "はじめに"
description: "Wails v3へのコントリビューションを始める方法"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## コントリビューターの皆さん、ようこそ！

Wails へのコントリビューションに関心をお寄せいただき、ありがとうございます。このガイドでは、初めてのコントリビューションの手順を説明します。

## 前提条件

始める前に、以下が準備できていることを確認してください。

- **Go 1.26+** がインストールされていること（[ダウンロード](https://go.dev/dl/)）
- **Node.js 20+** および **npm**（[ダウンロード](https://nodejs.org/)）
- GitHub アカウント用に設定された **Git**
- Go および JavaScript/TypeScript の基礎知識

### プラットフォーム固有の要件

**macOS：**

- Xcode Command Line Tools：`xcode-select --install`

**Windows：**

- MSYS2 または同様の Unix 互換環境を推奨
- WebView2 ランタイム（通常は Windows 11 にプリインストール済み）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（デフォルトの GTK4 スタック）
- インストール方法：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 従来の `-tags gtk3` ビルドパスを使用する場合は、`libgtk-3-dev` と `libwebkit2gtk-4.1-dev` もインストールしてください

## コントリビューション手順の概要

一般的なコントリビューションのワークフローは次のとおりです。

1. **フォークとクローン** - Wails リポジトリの自分用のコピーを作成します
2. **セットアップ** - Wails CLI をビルドし、環境を確認します
3. **ブランチ** - 変更用の機能ブランチを作成します
4. **開発** - コーディング規約に従って変更を加えます
5. **テスト** - テストを実行して、すべてが動作することを確認します
6. **コミット** - 明確で規約に沿ったコミットメッセージを付けてコミットします
7. **提出** - レビュー用のプルリクエストを作成します
8. **反復** - フィードバックに対応し、調整します
9. **マージ** - 承認されると、変更が Wails の一部になります！

## ステップ別ガイド

コントリビューションの種類を選択してください。

@tabs
[バグ修正]
@steps
### バグを探す、または報告する
- そのバグが [GitHub Issues](https://github.com/wailsapp/wails/issues) ですでに報告されていないか確認します
- 報告されていない場合は、再現手順を記載した新しい Issue を作成します
- 作業を開始する前に確認を待ちます

### フォークとクローン
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) でリポジトリをフォークします

フォークをクローンします。

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### ビルドと確認
Wails をビルドし、バグを再現できることを確認します。

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### バグ修正ブランチの作成
修正用のブランチを作成します。

```bash
git checkout -b fix/issue-123-window-crash
```

### バグの修正
- バグの修正に必要な最小限の変更を加えます
- 無関係なコードをリファクタリングしないでください
- リグレッションを防ぐため、テストを追加または更新します

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

PR の説明には、以下を記載します。

- バグと根本原因を説明します
- 修正内容を説明します
- Issue を参照します："Fixes #123"
- 修正前と修正後の動作を記載します

### フィードバックへの対応
レビューコメントに対応し、必要に応じて PR を更新します。

@end

[WEP（機能強化）]
@steps
### WEP の作成
- [WEP（Wails Enhancement Proposal）のプロセス](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)を確認します
- WEP テンプレートを `v3/wep/proposals/<name>/proposal.md` にコピーします
- WEP のみを含み、タイトルを `[WEP] <title>` としたドラフト PR を作成します
- 実装を開始する前に、メンテナーの決定を待ちます

### フォークとクローン
[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) でリポジトリをフォークします

フォークをクローンします。

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 開発環境のセットアップ
Wails をビルドし、環境を確認します。

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 機能ブランチの作成
内容が分かる名前のブランチを作成します。

```bash
git checkout -b feat/window-transparency-support
```

### 機能の実装
- [コーディング規約](/contributing/standards/)に従います
- 変更を機能に関係する範囲に限定します
- 明確で、適切に文書化されたコードを記述します
- 包括的なテストを追加します

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 十分なテスト
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

### 機能の文書化
- すべての公開 API にドキュメンテーションコメントを追加します
- `/docs/mpress/content/` 内の関連ドキュメントを更新します
- 該当する場合は例を追加します

### 規約に沿ったコミット
Conventional Commits を使用します。

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### プルリクエストの提出
プッシュして PR を作成します。

```bash
git push origin feat/window-transparency-support
```

PR には、以下を記載します。

- 機能とユースケースを説明します
- 例またはスクリーンショットを提示します
- 破壊的変更がある場合はすべて列挙します
- 承認された WEP の PR を参照します

### レビューに基づく改善
メンテナーから変更を求められる場合があります。辛抱強く、協力的に対応してください。

@end

[ドキュメント]
修正 PR は、事前に Issue を作成しなくても歓迎します。ドキュメントのみの  
修正では、失敗するコードテストは必要ありません。M-Press のインストール、  
ソースパス、プレビュー、検証、PR の手順については、[M-Press でドキュメントを作成する](/contributing/documentation/)に従ってください。

@end

## 取り組む Issue を探す

- [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue) ラベルを探す
- [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted) Issue を確認する
- [未解決の Issue](https://github.com/wailsapp/wails/issues)を確認し、担当への割り当てを依頼する

## サポートを受ける

- **Discord：**[Wails Discord](https://discord.gg/JDdSxwjhGf) に参加する
- **Discussions：**[GitHub Discussions](https://github.com/wailsapp/wails/discussions) に投稿する
- <strong>Issue：</strong>再現可能なバグは Issue を作成し、質問には Discussions、機能強化には WEP PR を使用する

## 行動規範

互いを尊重し、建設的かつ歓迎する姿勢で接してください。私たちは、優れたソフトウェアを共に作ることに注力する、親しみやすいコミュニティを築いています。

## 次のステップ

- [開発環境](/contributing/setup/)をセットアップする
- [コーディング規約](/contributing/standards/)を確認する
- [技術ドキュメント](/contributing/overview/)を確認する
