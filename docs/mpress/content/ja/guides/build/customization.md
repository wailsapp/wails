---
title: "ビルドのカスタマイズ"
description: "Task と Taskfile.yml を使用してビルドプロセスをカスタマイズする"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 概要

Wails のビルドシステムは、Wails アプリケーションのビルドプロセスを効率化するために設計された、柔軟で強力なツールです。タスクを簡単に定義して実行できるタスクランナー、[Task](https://taskfile.dev)を活用しています。v3 のビルドシステムがデフォルトですが、Wails は「独自のツールを持ち込む」アプローチを推奨しており、開発者は必要に応じてビルドプロセスをカスタマイズできます。

Task の使用方法について詳しくは、[公式ドキュメント](https://taskfile.dev/usage/)を参照してください。

## Task：ビルドシステムの中核

[Task](https://taskfile.dev) は、Go で記述された Make のモダンな代替ツールです。YAML ファイルを使用して、タスクとその依存関係を定義します。Wails のビルドシステムでは、[Task](https://taskfile.dev) がビルドプロセスのオーケストレーションにおいて中心的な役割を果たします。

メインの `Taskfile.yml` はプロジェクトルートにあり、プラットフォーム固有のタスクは `build/<platform>/Taskfile.yml` ファイルで定義されます。`build` ディレクトリ内の共通 `Taskfile.yml` ファイルには、プラットフォーム間で共有される共通タスクが含まれています。

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

プロジェクトルートにある `Taskfile.yml` ファイルは、ビルドシステムの主要なエントリーポイントです。このファイルでタスクとその依存関係を定義します。デフォルトの `Taskfile.yml` ファイルは次のとおりです。

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## プラットフォーム固有の Taskfile

各プラットフォームには、`build` ディレクトリ配下の各プラットフォーム用ディレクトリに専用の Taskfile があります。これらのファイルでは、そのプラットフォームの主要なタスクを定義します。各 Taskfile は、`build/Taskfile.yml` ファイルの共通タスクを取り込みます。

### Windows

場所：`build/windows/Taskfile.yml`

Windows 固有の Taskfile には、Windows 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- 任意の本番環境向けフラグを使用したビルド
- `.ico` アイコンファイルの生成
- Windows の `.syso` ファイルの生成
- パッケージ化用 NSIS インストーラーの作成

### Linux

場所：`build/linux/Taskfile.yml`

Linux 固有の Taskfile には、Linux 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- 任意の本番環境向けフラグを使用したビルド
- AppImage、deb、rpm、および Arch Linux パッケージの作成
- Linux アプリケーション用 `.desktop` ファイルの生成

### macOS

場所：`build/darwin/Taskfile.yml`

macOS 固有の Taskfile には、macOS 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- amd64、arm64、および universal（両方）の各アーキテクチャ向けバイナリのビルド
- `.icns` アイコンファイルの生成
- 配布用 `.app` バンドルの作成
- `.app` バンドルのアドホック署名
- macOS 固有のビルドフラグと環境変数の設定

## タスクの実行とコマンドエイリアス

`wails3 task` コマンドは [Taskfile](https://taskfile.dev) の組み込み版であり、`Taskfile.yml` に定義されているタスクを実行します。

`wails3 build` コマンドと `wails3 package` コマンドは、それぞれ `wails3 task build` と `wails3 task package` のエイリアスです。これらのコマンドを実行すると、Wails は内部で適切なタスク実行に変換します。

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### タスクへのパラメーターの受け渡し

`KEY=VALUE` 形式を使用して、CLI 変数をタスクに渡せます。これらの変数はエイリアスコマンドを介して転送されます。

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

`Taskfile.yml` では、Go テンプレート構文を使用してこれらの変数にアクセスできます。

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

## 共通のビルドプロセス

すべてのプラットフォームに共通して、ビルドプロセスには通常、次の手順が含まれます。

1. Go モジュールの整理
2. フロントエンドのビルド
3. アイコンの生成
4. プラットフォーム固有のフラグを使用した Go コードのコンパイル
5. アプリケーションのパッケージ化（プラットフォーム固有）

## ビルドプロセスのカスタマイズ

v3 のビルドシステムには実用的なデフォルト設定が用意されていますが、プロジェクトのニーズに合わせて簡単にカスタマイズできます。`Taskfile.yml` とプラットフォーム固有の Taskfile を変更することで、次のことができます。

- 新しいタスクの追加
- 既存のタスクの変更
- タスク実行順序の変更
- 他のツールやスクリプトとの統合

この柔軟性により、Wails ビルドシステムが提供する構造を活用しながら、個別の要件に合わせてビルドプロセスを調整できます。

@note{type="tip" title="Taskfile の学習"}
Taskfile を効果的に使用する方法を理解するために、[Taskfile](https://taskfile.dev) のドキュメントを読むことを強く推奨します。Wails CLI に組み込まれている Taskfile のバージョンは、`wails3 task --version` を実行すると確認できます。

@end

## 開発モード

Wails のビルドシステムには、ライブリロードとホットモジュール置換によって開発者体験を向上させる強力な開発モードが含まれています。このモードは、`wails3 dev` コマンドを使用して有効にします。

### 仕組み

`wails3 dev`を実行すると、次の処理が行われます。

1. コマンドは利用可能なポートを確認し、指定されていない場合はデフォルトで9245を使用します。
2. フロントエンド開発サーバー（Vite）用の環境変数を設定します。
3. [refresh](https://github.com/atterpac/refresh)ライブラリを使用してファイル監視を開始します。

[refresh](https://github.com/atterpac/refresh)ライブラリは、ファイルの変更を監視し、再ビルドをトリガーします。`./build/config.yml`ファイルの`dev_mode`キーに定義された設定を使用します。  
特定のディレクトリやファイルを無視する設定、監視対象とするファイルの指定、変更を検出した際に実行する処理の指定が可能です。  
デフォルト設定でも十分に機能しますが、必要に応じて自由にカスタマイズできます。

### 設定

構造の例を次に示します。

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

この設定ファイルでは、次のことができます。

- ファイル監視のルートパスを設定する
- ログレベルを設定する
- ファイル変更イベントのデバウンス時間を設定する
- 特定のディレクトリ、ファイル、またはファイル拡張子を無視する
- ファイル変更時に実行するコマンドを定義する

### 開発モードのカスタマイズ

`config.yml`ファイル内のこれらの値を変更することで、開発モードをカスタマイズできます。

カスタマイズ方法には、次のようなものがあります。

1. 監視するディレクトリやファイルを変更する
2. システムが変更に応答するまでの速さを制御するため、デバウンス時間を調整する
3. プロジェクトの要件に合わせて実行コマンドを追加または変更する

### 開発でのブラウザの使用

Wails v2では開発時のブラウザ使用を完全にサポートしていましたが、多くの混乱を招いていました。ブラウザで動作するアプリケーションがデスクトップアプリケーションでも必ず動作するとは限りません。これは、ブラウザAPIの一部がWebViewでは利用できないためです。

UIを中心とする開発作業では、v3でも開発モードで`http://localhost:9245`のVite URLにアクセスすることで、ブラウザを柔軟に使用できます。これにより、スタイルやレイアウトの作業中に高機能なブラウザ開発者ツールを利用できます。ただし、このモードではGoバインディングが<em>動作しない</em>ことに注意してください。  
バインディングやイベントなどの機能をテストする段階になったら、デスクトップ表示に切り替え、本番環境ですべてが正しく動作することを確認してください。
