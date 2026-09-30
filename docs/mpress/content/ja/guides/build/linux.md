---
title: "Linux パッケージング"
description: "Wails アプリケーションを Linux で配布するためにパッケージ化します"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## パッケージ形式

アプリを Linux で配布するためにパッケージ化します。

```bash
wails3 package GOOS=linux
```

これにより、`bin/` ディレクトリに複数の形式が作成されます。

- **AppImage**：そのまま持ち運べて、どの Linux ディストリビューションでも実行可能
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

`.desktop` ファイルは、アプリケーションメニューでのアプリの表示方法を制御します。このファイルは、`build/linux/Taskfile.yml` の値から生成されます。

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

署名は `build/linux/Taskfile.yml` で設定します。

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

キーのパスワードを保存します。

```bash
wails3 setup signing
```

詳細については、[アプリケーションへの署名](/guides/build/signing/)を参照してください。

## ARM 向けのビルド

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
x86_64 ホストからの ARM64 ビルドでは、CGO のクロスコンパイルに Docker を使用します。

@end

## 従来の GTK3 サポート

Wails v3 は、デフォルトでは <strong>GTK4 と WebKitGTK 6.0</strong>を使用してビルドされます。WebKitGTK 6.0 がまだ提供されていないディストリビューション（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x）向けに、従来の GTK3 / WebKit2GTK 4.1 パスも引き続き利用できます。この従来のパスはビルドタグで明示的に有効にする必要があり、v3.1 で削除される予定です。

@note{type="caution" title="従来のパス"}
GTK3 / WebKit2GTK 4.1 パスは、v3.0.x 系列を通じてサポートされます。対象ディストリビューションでの GTK4 / WebKitGTK 6.0 の提供状況に合わせて、GTK4 への移行を計画してください。`-tags gtk3` は v3.1 で削除されます。

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

- **ファイルダイアログ**：GTK4 は、デフォルトでファイルダイアログに `xdg-desktop-portal` を使用します。そのため、デフォルトディレクトリやカスタムフィルターの表示など、一部のダイアログオプションの動作が GTK3 とは異なります。詳細については、[ダイアログリファレンス - Linux でのダイアログの動作](/reference/dialogs/#linux-dialog-behavior)を参照してください。
- **メニュースタイル**：GTK4 は、GNOME HIG に従ってヘッダーバーにハンバーガーボタン（☰）を表示する `LinuxMenuStylePrimaryMenu` オプションをサポートしています。このオプションは `-tags gtk3` ビルドには影響しません。[Window API - Linux MenuStyle](/reference/window/#linux)を参照してください。
- **DPI スケーリング**：GTK4 は、フラクショナルスケーリングのサポートに `gdk_monitor_get_scale`（GTK 4.14+）を使用します。

### ビルドの確認

`wails3 doctor` を実行して、セットアップを確認します。フラグを指定しない場合は、デフォルトの GTK4 / WebKitGTK 6.0 を確認します。従来の GTK3 / WebKit2GTK 4.1 パッケージは、オプションとして一覧表示されます。

## トラブルシューティング

### AppImage を実行できない

実行可能にします。

```bash
chmod +x MyApp-x86_64.AppImage
```

### 依存関係が不足している

アプリを起動できない場合は、不足している WebKit の依存関係を確認します。

```bash
# Debian/Ubuntu
sudo apt install libwebkit2gtk-4.1-0

# Fedora
sudo dnf install webkit2gtk4.1

# Arch
sudo pacman -S webkit2gtk-4.1
```

### C コンパイラーが見つからない

ビルドシステムでは、CGO のために GCC または Clang が必要です。

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

または、`wails3 task setup:docker` を実行すると、ビルドシステムが自動的に Docker を使用します。

### NVIDIA GPU でウィンドウが空白または白くなる

NVIDIA のプロプライエタリドライバーを使用する Linux では、Wails アプリの起動時にウィンドウが空白または白くなることがあります。これは、NVIDIA のプロプライエタリドライバーに対して `gbm_bo_map()` を使用すると DMA-BUF レンダラーが失敗する WebKitGTK のバグが原因です（X11 と Wayland、ドライバーバージョン 377～580+、10 シリーズおよびそれ以前の GT 710 GPU に影響します）。

<strong>Wails は NVIDIA カーネルモジュール（`/sys/module/nvidia`）を検出すると、`WEBKIT_DISABLE_DMABUF_RENDERER=1` を自動的に適用する</strong>ため、ほとんどのユーザーは何もする必要がありません。

それでも空白のウィンドウが表示される場合（たとえば、モジュールのパスが見えないコンテナー内など）は、アプリを起動する前に環境変数を手動で設定します。

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

関連するアップストリームのバグ：[WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607)、[WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739)。

### AppImage の strip 互換性

最新の Linux ディストリビューション（Arch Linux、Fedora 39+、Ubuntu 24.04+）では、より効率的な再配置のために、システムライブラリが `.relr.dyn` ELF セクションを使用してコンパイルされています。AppImage の作成に使用される `linuxdeploy` ツールには、これらの最新のセクションを処理できない古い `strip` バイナリが同梱されています。

Wails は、AppImage をビルドする前にシステムの GTK ライブラリを確認し、この状況を自動的に検出します。検出された場合は、互換性を確保するためにストリップ処理が無効になります（`NO_STRIP=1`）。

**これが意味すること：**

- 影響を受けるシステムでは、AppImage のサイズがわずかに大きくなります（約 20-40%）
- アプリケーションの機能には影響しません
- これは自動的に処理されるため、対応は不要です

最新のシステムで AppImage のサイズを小さくする必要がある場合は、より新しい `strip` バイナリをインストールし、同梱版の代わりにそのバイナリを使用するように `linuxdeploy` を設定できます。
