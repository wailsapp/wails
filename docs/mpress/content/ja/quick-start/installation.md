---
title: "インストール"
description: "Wails をインストールし、アプリケーションをビルドできる状態にします"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## クイックインストール（5 分）

@note{type="tip" title="要点 — 経験豊富な開発者向け"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

または、`wails3 doctor` で手動確認します。[最初のアプリへ進む →](/quick-start/first-app/)

@end

## ステップごとのインストール手順

@steps
### Go のインストール（必須）
Wails には Go 1.25 以降が必要です。

@tabs{sync-key="os"}
[Windows]
**[go.dev/dl](https://go.dev/dl/)** から Windows インストーラーをダウンロードして実行します。

**インストールを確認：**

```powershell
go version  # Should show 1.25 or later
```

**PATH を確認：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

空の場合は、`C:\Users\YourName\go\bin` を PATH に追加します。

[macOS]
**方法 1：公式インストーラー**

**[go.dev/dl](https://go.dev/dl/)** から macOS インストーラー（.pkg ファイル）をダウンロードして実行します。

**方法 2：Homebrew**

```bash
brew install go
```

**インストールを確認：**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

`~/go/bin` が PATH にない場合は、`~/.zshrc` または `~/.bash_profile` に追加します：

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**方法 1：公式 tarball**

**[go.dev/dl](https://go.dev/dl/)** から Linux の tarball をダウンロードし、次を実行します：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**方法 2：パッケージマネージャー**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**PATH に追加**（`~/.bashrc` または `~/.zshrc` に追加）：

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**確認：**

```bash
go version
echo $PATH | grep go/bin
```

@end

### プラットフォーム依存関係のインストール
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime**（通常はプリインストール済み）

Windows 10/11 には、WebView2 が標準で含まれています。含まれていない場合：

- [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/) からダウンロードします
- または、後で `wails3 doctor` を実行します。手順が案内されます

**これで完了です！** ほかの依存関係は必要ありません。

@note{type="tip" title="Windows 11 のパフォーマンス向上のヒント"}
プロジェクトの保存先として [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) の利用を検討してください。Dev Drive は開発者向けワークロードに最適化されており、ビルド時間とディスクアクセス速度を最大 30% 大幅に改善できます。

@end

[macOS]
**Xcode Command Line Tools**（必須）

```bash
xcode-select --install
```

表示されたダイアログで「Install」をクリックします。

**確認：**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**これで完了です！** macOS には WebKit が標準で含まれています。

[Linux]
**ビルドツールと WebKit**

@note{type="caution" title="ディストリビューションの最小バージョン"}
Wails v3 では、デフォルトで **WebKitGTK 6.0** が必要です。WebKit2GTK 4.1 のみを提供するディストリビューション（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x）では、レガシーの `-tags gtk3` を明示的に有効にしてビルドする必要があります。WebKit2GTK 4.0 のみを提供するさらに古いリリース（Ubuntu 20.04、Debian 11、RHEL 8）はサポートされていません。

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
デフォルトの GTK4 スタックには Ubuntu 24.04 以降または Debian 13 以降が必要です。

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
使用している `shell.nix` または `devShell` に追加します：

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[その他]
Wails のインストール後に `wails3 doctor` を実行します。使用しているディストリビューションに必要なパッケージが正確に表示されます。

@end

@note{type="info" title="レガシー GTK3 スタック"}
対象のディストリビューションがまだ WebKitGTK 6.0 を提供していない場合（Ubuntu 22.04 LTS、Debian 12 など）は、代わりに GTK3 と WebKit2GTK 4.1 の開発ライブラリ（Debian/Ubuntu では `libgtk-3-dev libwebkit2gtk-4.1-dev`、その他のディストリビューションでは相当するもの）をインストールし、`wails3 build -tags gtk3` を指定してビルドします。レガシー方式は v3.0.x 系までサポートされ、v3.1 で削除されます。詳細については、[Linux パッケージング — レガシー GTK3 サポート](/guides/build/linux/#legacy-gtk3-support)を参照してください。

@end

@end

### Wails CLI のインストール
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

これにより、`wails3` コマンドが `~/go/bin`（Windows では `%USERPROFILE%\go\bin`）にインストールされます。

### セットアップウィザードの実行（推奨）
```bash
wails3 setup
```

セットアップウィザードは、依存関係を確認し、不足しているもののインストールを支援して、プロジェクトのデフォルト設定を構成します。

@note{type="caution" title="試験的機能"}
セットアップウィザードは新しい機能で、主に Linux でテストされています。問題が発生した場合は、[報告](https://github.com/wailsapp/wails/issues/4904)し、代わりに `wails3 doctor` を使用してください。

@end

### インストールの確認
```bash
wails3 doctor
```

**想定される出力（または同様の出力）：**

```
Wails (v3.0.0-dev)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | Your installed version |
| Go Version  | go1.25.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="`wails3` コマンドが見つからない場合"}
`~/go/bin` が PATH にありません。上記の手順 1 を参照して修正し、ターミナルを再起動してください。

@end

### npm のインストール（任意、ただし推奨）
ほとんどの Wails テンプレートでは、フロントエンドのツールに npm を使用します。

@tabs{sync-key="os"}
[Windows]
[nodejs.org](https://nodejs.org/) からダウンロードし、インストーラーを実行します。

**確認：**

```powershell
npm --version
```

[macOS]
**オプション 1：公式インストーラー** [nodejs.org](https://nodejs.org/) からダウンロードします

**オプション 2：Homebrew**

```bash
brew install node
```

**確認：**

```bash
npm --version
```

[Linux]
**オプション 1：NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**オプション 2：パッケージマネージャー**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**確認：**

```bash
npm --version
```

@end

@note{type="tip" title="その他のパッケージマネージャー"}
`pnpm`、`yarn`、または `bun` を使いたい場合でも問題ありません。プロジェクト内の `Taskfile.yml` を更新し、使用するツールを指定してください。

@end

@end

## トラブルシューティング

### `wails3` コマンドが見つからない

**原因：** `~/go/bin`（または `%USERPROFILE%\go\bin`）が PATH に含まれていません。

**解決方法：**

@tabs{sync-key="os"}
[Windows]
1. 「環境変数」を開きます（スタートメニューで検索してください）
2. 「ユーザー環境変数」で `Path` を探します
3. 「編集」→「新規」をクリックします
4. `C:\Users\YourName\go\bin` を追加します（`YourName` は置き換えてください）
5. すべてのダイアログで「OK」をクリックします
6. **ターミナルを再起動します**

**確認：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
`~/.zshrc`（macOS）または `~/.bashrc`（Linux）に追加します：

```bash
export PATH=$PATH:~/go/bin
```

再読み込み：

```bash
source ~/.zshrc  # or ~/.bashrc
```

**確認：**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` で依存関係の不足が報告される

<strong>Linux：</strong>出力には、インストールが必要なパッケージが正確に示されます。例：

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

<strong>Windows：</strong>WebView2 が見つからない場合：

- [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/) からダウンロードします
- または、最初のアプリを実行すると自動的にインストールされます

<strong>macOS：</strong>Xcode ツールが見つからない場合：

```bash
xcode-select --install
```

---

#### Go のバージョンが古すぎる

Wails v3 には Go 1.25 以降が必要です。これより古いバージョンを使用している場合：

@tabs{sync-key="os"}
[Windows/macOS]
[go.dev/dl](https://go.dev/dl/) から最新版をダウンロードし、再インストールします。

[Linux]
[go.dev/dl](https://go.dev/dl/) から最新の tarball をダウンロードし、次の操作を行います：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## 開発版（最新開発版）

メインの開発ブランチにある最新コードを使用したい場合は、リリース前の新機能や修正を利用できます。ただし、バグや破壊的変更が含まれる可能性があります。コントリビューター、または今後追加される機能をテストする必要がある方にのみ推奨します。

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="開発版"}
- バグや破壊的変更が含まれる可能性があります
- 作成されるプロジェクトでは、ローカルの Wails を参照するために `replace` ディレクティブが使用されます
- コントリビューター、または新機能のテストにのみ推奨します

@end

## 次のステップ

<strong>インストールが完了しました！</strong>お使いのシステムで Wails を使って開発する準備が整いました。

@cards{cols="1"}
🚀 最初のアプリを作成する
10 分で動作するアプリケーションを作成します。

[最初のアプリのチュートリアル →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 テンプレートを見る
すぐに使えるテンプレートを確認します。

```bash
wails3 init -l  # List templates
```

@end

---

**問題が発生しましたか？**[Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[Issue を作成](https://github.com/wailsapp/wails/issues)してください。
