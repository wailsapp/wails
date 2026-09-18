---
title: "開発環境のセットアップ"
description: "Wails v3の開発環境をセットアップします"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## 開発環境のセットアップ

このガイドでは、Wails v3の開発に必要な一通りの開発環境をセットアップする手順を説明します。

## 必要なツール

### Go開発環境

1. **Go 1.25以降をインストールします：**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Go環境を設定します：**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **便利なGoツールをインストールします：**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.jsとnpm

フロントエンド統合の例を使用する場合にのみ必要です。

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### プラットフォーム固有の依存関係

**macOS：**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows：**

1. Unixに似た環境を利用するために[MSYS2](https://www.msys2.org/)をインストールします
2. WebView2 Runtime（Windows 11にはプリインストールされています。Windows 10向けには[ダウンロード](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)してください）
3. 任意：[Git for Windows](https://git-scm.com/download/win)をインストールします

**Linux（Debian/Ubuntu）：**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux（Fedora/RHEL）：**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux（Arch）：**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## リポジトリのセットアップ

### クローンと設定

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Wails CLIのビルド

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### PATHへの追加（任意）

**Linux/macOS：**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows：**

システムのプロパティから、WailsディレクトリをPATH環境変数に追加します。

## IDEのセットアップ

### VS Code（推奨）

1. **VS Codeをインストールします：** [ダウンロード](https://code.visualstudio.com/)

2. **拡張機能をインストールします：**
  - Go（GoogleのGo Team製）
  - ESLint
  - Prettier
  - MDX（ドキュメント用）


3. **ワークスペース設定**（`.vscode/settings.json`）を構成します：
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **GoLandをインストールします：** [ダウンロード](https://www.jetbrains.com/go/)

2. **次の項目を設定します：**
  - Go modulesのサポートを有効にします
  - `goimports`用のファイルウォッチャーを設定します
  - プロジェクトの規約に合わせてコードスタイルを設定します


## セットアップの確認

次のコマンドを実行して、すべてが正常に動作することを確認します：

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

テストアプリをビルドして実行できれば、環境の準備は完了です。

## テストの実行

### ユニットテスト

```bash
cd v3
go test ./...
cd ..
```

### 特定のパッケージのテスト

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### カバレッジを有効にして実行

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### レースディテクターを有効にして実行

```bash
cd v3
go test ./... -race
```

## ドキュメントの作業

Wails v3 ドキュメントは M-Press で記述されています。英語のソースファイルは  
`docs/mpress/content/` にあり、翻訳ファイルは `fr/` や  
`id/` などの言語別ディレクトリにあります。

リポジトリのルートからドキュメントの変更をプレビューし、検証します。

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

本番サイトは静的サイトです。ローカルでドキュメントを編集する場合、Node.js、翻訳プロバイダー、Cloudflare  
の認証情報は必要ありません。

## デバッグ

### Goコードのデバッグ

**VS Code：**

`.vscode/launch.json`を作成します：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**コマンドライン：**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### プラットフォームコードのデバッグ

プラットフォーム固有のデバッグには、各プラットフォーム用のツールが必要です：

- **macOS：** Xcode Instruments
- **Windows：** Visual Studio Debugger
- **Linux：** GDB

## よくある問題

### 「command not found: wails3」

Wails ディレクトリを PATH に追加するか、プロジェクトのルートから `./wails3` を使用してください。

### 「webkitgtk-6.0 not found」または「webkit2gtk not found」（Linux）

ビルド対象のスタックに対応する開発パッケージをインストールしてください。

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Go モジュールのエラーによりビルドに失敗する

```bash
cd v3
go mod tidy
go mod download
```

### Windows での「CGO_ENABLED」エラー

C コンパイラ（MSYS2 経由の MinGW-w64）が PATH に含まれていることを確認してください。

## 次のステップ

- [コーディング規約](/contributing/standards/)を確認する
- [技術ドキュメント](/contributing/)を参照する
- 取り組む Issue を探す：[初心者向け Issue](https://github.com/wailsapp/wails/labels/good%20first%20issue)
