---
title: "アプリケーションのビルド"
description: "Wails アプリケーションをビルドしてパッケージ化する"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 は、ビルドシステムとして [Task](https://taskfile.dev) を使用します。`wails3 build` コマンドと `wails3 package` コマンドは、Task を簡単に使用するためのラッパーです。

## ビルド

現在のプラットフォーム向けにビルドします：

```bash
wails3 build
```

特定のプラットフォーム向けにビルドします：

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

出力は `bin/` ディレクトリに生成されます。

@note{type="tip"}
別のプラットフォームから macOS または Linux 向けにクロスコンパイルするには、Docker が必要です。セットアップについては、[クロスプラットフォームビルド](/guides/build/cross-platform/)を参照してください。

@end

## 開発

ホットリロードを有効にしてアプリケーションを実行します：

```bash
wails3 dev
```

ファイル監視が開始され、変更があるとアプリケーションが再ビルドされ、再起動されます。フロントエンド開発サーバーは、デフォルトでポート 9245 で動作します。

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## パッケージ化

配布用にアプリケーションをパッケージ化します：

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

プラットフォーム固有のパッケージが作成されます：

- **Windows**：NSIS インストーラー — [Windows のパッケージ化](/guides/build/windows/)を参照してください
- **macOS**：アプリケーションバンドル（`.app`）— [macOS のパッケージ化](/guides/build/macos/)を参照してください
- **Linux**：AppImage、deb、rpm — [Linux のパッケージ化](/guides/build/linux/)を参照してください

## カスタムビルドタグ

`-tags` フラグを使用して、カスタム Go ビルドタグを渡します：

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

タグは、基盤となる Taskfile に `EXTRA_TAGS` として渡されます。詳細については、[サーバーのビルド](/guides/server-build/)および[Linux のパッケージ化 - 従来の GTK3 サポート](/guides/build/linux/#legacy-gtk3-support)を参照してください。

## Task の直接使用

より細かく制御するには、Task を直接使用します：

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

`linux:create:deb` や `darwin:build:universal` などのプラットフォーム固有のタスクは、Task からのみ利用できます。

## アセットの生成

アイコンを再生成するか、ビルド設定を更新します：

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
