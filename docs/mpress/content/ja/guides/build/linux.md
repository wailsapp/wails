---
title: "Linux パッケージング"
description: "Wails アプリケーションを Linux 向けにパッケージ化して配布する"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## パッケージ形式

アプリを Linux 向けにパッケージ化して配布します。

```bash
wails3 package GOOS=linux
```

これにより、`bin/` ディレクトリに複数の形式が作成されます。

- **AppImage**：ポータブルで、あらゆる Linux ディストリビューション上で実行可能
- **DEB**：Debian、Ubuntu、およびその派生ディストリビューション向け
- **RPM**：Fedora、RHEL、およびその派生ディストリビューション向け
- **Arch**：Arch Linux およびその派生ディストリビューション向け

### 個別の形式

特定の形式をビルドします。

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## パッケージのカスタマイズ

### デスクトップエントリ

アプリケーションメニューでのアプリの表示方法は、`.desktop` ファイルによって制御されます。このファイルは `build/linux/Taskfile.yml` の値から生成されます。

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### パッケージのメタデータ

DEB および RPM パッケージをカスタマイズするには、`build/linux/nfpm/nfpm.yaml` を編集します。

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

AppImage の設定は `build/linux/appimage/` にあります。アプリアイコンには `build/appicon.png` が使用されます。

AppImage には WebKit とそのヘルパープロセスが同梱されており、同梱ライブラリには、それらを AppImage 内から起動するためのパッチが適用されています。AppImage の `AppRun` は、AppImage 内の `usr/lib` からアプリを起動します。パッチ適用済みのパスは WebKit の bubblewrap サンドボックスの外部と内部の両方で解決されるため、サンドボックスは有効なままです。この作業ディレクトリを維持できる場合に限り、独自の `AppRun` を指定してください。

## パッケージへの署名

PGP キーを使用して DEB および RPM パッケージに署名します。

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

`build/linux/Taskfile.yml` で署名を設定します。

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

キーのパスワードを保存します。

```bash
wails3 signing
```

詳しくは、[アプリケーションへの署名](/guides/build/signing/)を参照してください。

## ARM 向けのビルド

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
x86_64 ホストからの ARM64 ビルドでは、CGO のクロスコンパイルに Docker を使用します。

@end

## 従来の GTK3 サポート

Wails v3 は、デフォルトでは **GTK4 と WebKitGTK 6.0** を使用してビルドされます。WebKitGTK 6.0 をまだ提供していないディストリビューション（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x）向けには、従来の GTK3 / WebKit2GTK 4.1 パスも引き続き利用できます。従来のパスはビルドタグで明示的に有効化する必要があり、v3.1 で削除される予定です。

@note{type="caution" title="従来のパス"}
GTK3 / WebKit2GTK 4.1 パスは、v3.0.x 系列までサポートされます。対象ディストリビューションでの GTK4 / WebKitGTK 6.0 の提供状況に合わせて、GTK4 への移行を計画してください。`-tags gtk3` は v3.1 で削除されます。

@end

### 依存関係

GTK3 および WebKit2GTK 4.1 の開発用ライブラリをインストールします。

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

必要な pkg-config パッケージは `gtk+-3.0` と `webkit2gtk-4.1` です。

### GTK3 を使用したビルド

`-tags gtk3` フラグを使用します。

```bash
wails3 build -tags gtk3
```

または、Go で直接ビルドします。

```bash
go build -tags gtk3 -o myapp .
```

### GTK4 との既知の相違点

- **ファイルダイアログ**：GTK4 では、デフォルトでファイルダイアログに `xdg-desktop-portal` を使用します。そのため、一部のダイアログオプション（デフォルトディレクトリやカスタムフィルターの表示など）の動作が GTK3 とは異なります。詳しくは、[ダイアログリファレンス - Linux でのダイアログの動作](/reference/dialogs/#linux-dialog-behavior)を参照してください。
- **メニュースタイル**：GTK4 は、GNOME HIG に準拠してヘッダーバーにハンバーガーボタン（☰）を表示する `LinuxMenuStylePrimaryMenu` オプションをサポートしています。このオプションは `-tags gtk3` ビルドには影響しません。[Window API - Linux MenuStyle](/reference/window/#linux) を参照してください。
- **DPI スケーリング**：GTK4 は、小数倍率のスケーリングをサポートするために `gdk_monitor_get_scale`（GTK 4.14 以降）を使用します。

### ビルドの確認

セットアップを確認するには、`wails3 doctor` を実行します。フラグを指定しない場合は、デフォルトの GTK4 / WebKitGTK 6.0 を確認します。従来の GTK3 / WebKit2GTK 4.1 パッケージは任意の依存関係として表示されます。

## トラブルシューティング

### AppImage が実行できない

実行可能にします。

```bash
chmod +x MyApp-x86_64.AppImage
```

### 依存関係が不足している

バイナリのビルドバリアントに合ったライブラリを使用してください。デフォルトの GTK4/WebKitGTK 6.0 ビルドの場合：

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

`-tags gtk3` を指定して明示的にビルドしたバイナリの場合は、代わりに従来の GTK3/WebKit2GTK 4.1 ライブラリをインストールします。

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

パッケージの提供状況と名前は、ディストリビューションのリリースによって異なります。必要なライブラリを備えたサポート対象リリース、または文書化されているクロスビルド／コンテナの手順を使用してください。従来のライブラリをインストールしても、デフォルトビルドの要件は満たされません。AppImage の可搬性は、実際にパッケージへ同梱されるものと、引き続き必要となるホストライブラリによって決まります。

### C コンパイラが見つからない

ビルドシステムでは、CGO 用に GCC または Clang が必要です。

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

代わりに `wails3 task setup:docker` を実行すると、ビルドシステムが自動的に Docker を使用します。

### NVIDIA GPU でウィンドウが空白または白く表示される

Linux で NVIDIA のプロプライエタリドライバーを使用している場合、Wails アプリの起動時にウィンドウが空白または白く表示されることがあります。これは、NVIDIA のプロプライエタリドライバーに対して DMA-BUF レンダラーが `gbm_bo_map()` で失敗する WebKitGTK のバグが原因です（X11 と Wayland、ドライバーバージョン 377～580 以降、10 シリーズおよびそれ以前の GT 710 GPU に影響します）。

**Wails は NVIDIA カーネルモジュール（`/sys/module/nvidia`）を検出すると、`WEBKIT_DISABLE_DMABUF_RENDERER=1` を自動的に適用する**ため、ほとんどのユーザーは何もする必要がありません。

それでもウィンドウが空白になる場合（モジュールのパスが見えないコンテナ内など）は、アプリを起動する前に環境変数を手動で設定してください。

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

関連するアップストリームのバグ：[WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607)、[WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739)。

### AppImage の strip 互換性

最新の Linux ディストリビューション（Arch Linux、Fedora 39 以降、Ubuntu 24.04 以降）では、より効率的な再配置を実現するため、システムライブラリが `.relr.dyn` ELF セクションを使用してコンパイルされています。AppImage の作成に使用する `linuxdeploy` ツールには、これらの新しいセクションを処理できない古い `strip` バイナリが同梱されています。

Wails は AppImage をビルドする前にシステムの GTK ライブラリを確認し、この状況を自動的に検出します。検出された場合は、互換性を確保するためにストリップ処理が無効化されます（`NO_STRIP=1`）。

**影響：**

- 影響を受けるシステムでは、AppImage のサイズが若干大きくなります（約 20-40%）
- アプリケーションの機能には影響しません
- 自動的に処理されるため、対応は不要です

最新のシステムで AppImage を小さくする必要がある場合は、より新しい `strip` バイナリをインストールし、同梱版の代わりにそのバイナリを使用するよう `linuxdeploy` を設定できます。
