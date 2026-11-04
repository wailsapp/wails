---
title: "ビルドのカスタマイズ"
description: "Task と Taskfile.yml を使用してビルドプロセスをカスタマイズする"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 概要

Wails のビルドシステムは、Wails アプリケーションのビルドプロセスを効率化するために設計された、柔軟で強力なツールです。タスクを簡単に定義して実行できるタスクランナー、[Task](https://taskfile.dev) を活用しています。デフォルトでは v3 ビルドシステムを使用しますが、Wails は「好みのツールを持ち込む」アプローチを推奨しており、開発者は必要に応じてビルドプロセスをカスタマイズできます。

Task の使用方法について詳しくは、[公式ドキュメント](https://taskfile.dev/usage/)を参照してください。

## Task：ビルドシステムの中核

[Task](https://taskfile.dev) は Go で記述された、Make に代わるモダンなツールです。YAML ファイルを使用してタスクとその依存関係を定義します。Wails のビルドシステムでは、[Task](https://taskfile.dev) がビルドプロセスの統括において中心的な役割を果たします。

メインの `Taskfile.yml` はプロジェクトルートにあり、プラットフォーム固有のタスクは `build/<platform>/Taskfile.yml` ファイルで定義されています。`build` ディレクトリにある共通の `Taskfile.yml` ファイルには、プラットフォーム間で共有される共通タスクが含まれています。

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

プロジェクトルートにある `Taskfile.yml` ファイルは、ビルドシステムのメインエントリーポイントです。タスクとその依存関係を定義します。デフォルトの `Taskfile.yml` ファイルは次のとおりです。

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

各プラットフォームには、それぞれ専用の Taskfile があり、`build` ディレクトリ配下の各プラットフォーム用ディレクトリに配置されています。これらのファイルは、そのプラットフォームの主要なタスクを定義します。各 Taskfile には、`build/Taskfile.yml` ファイルの共通タスクが含まれています。

### Windows

場所：`build/windows/Taskfile.yml`

Windows 固有の Taskfile には、Windows 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- 任意のプロダクション用フラグを指定したビルド
- `.ico` アイコンファイルの生成
- Windows `.syso` ファイルの生成
- パッケージ化用の NSIS インストーラーの作成

### Linux

場所：`build/linux/Taskfile.yml`

Linux 固有の Taskfile には、Linux 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- 任意のプロダクション用フラグを指定したビルド
- AppImage、deb、rpm、および Arch Linux パッケージの作成
- Linux アプリケーション用の `.desktop` ファイルの生成

### macOS

場所：`build/darwin/Taskfile.yml`

macOS 固有の Taskfile には、macOS 上でアプリケーションをビルド、パッケージ化、実行するためのタスクが含まれています。主な機能は次のとおりです。

- amd64、arm64、およびユニバーサル（両方）アーキテクチャ向けバイナリのビルド
- `.icns` アイコンファイルの生成
- 配布用の `.app` バンドルの作成
- `.app` バンドルのアドホック署名
- macOS 固有のビルドフラグと環境変数の設定

## タスクの実行とコマンドエイリアス

`wails3 task` コマンドは、`Taskfile.yml` で定義されたタスクを実行する、[Taskfile](https://taskfile.dev) の組み込み版です。

`wails3 build` コマンドと `wails3 package` コマンドは、それぞれ `wails3 task build` と `wails3 task package` のエイリアスです。これらのコマンドを実行すると、Wails は内部で適切なタスク実行に変換します。

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### タスクへのパラメーターの受け渡し

`KEY=VALUE` 形式を使用して、CLI 変数をタスクに渡すことができます。これらの変数はエイリアスコマンドを介して転送されます。

```bash
# These are equivalent:
wails3 build GOOS=linux
wails3 task build GOOS=linux
```

標準タスクでは、`GOOS`、`ARCH`、`DEV` を使用します。`VERSION`、`OUTPUT`、`PLATFORM`、`CONFIG` などの変数は、独自の Taskfile でその変数を読み取る場合にのみタスクへ影響します。次に示すのはカスタム Taskfile の例であり、Wails の標準ビルドタスクではありません。

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

カスタムのパッケージングタスクが `VERSION` と `OUTPUT` を読み取る場合も、同じ方法で渡せます。

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

## 共通のビルドプロセス

すべてのプラットフォームで、ビルドプロセスには通常、次の手順が含まれます。

1. Go モジュールの整理
2. フロントエンドのビルド
3. アイコンの生成
4. プラットフォーム固有のフラグを指定した Go コードのコンパイル
5. アプリケーションのパッケージ化（プラットフォーム固有）

## ビルドプロセスのカスタマイズ

v3 ビルドシステムには堅実なデフォルト設定が用意されていますが、プロジェクトのニーズに合わせて簡単にカスタマイズできます。`Taskfile.yml` とプラットフォーム固有の Taskfile を変更することで、次のことができます。

- 新しいタスクの追加
- 既存タスクの変更
- タスク実行順序の変更
- ほかのツールやスクリプトとの統合

この柔軟性により、Wails ビルドシステムが提供する構造を活用しながら、個別の要件に合わせてビルドプロセスを調整できます。

@note{type="tip" title="Taskfile を学ぶ"}
Taskfile を効果的に使用する方法を理解するために、[Taskfile](https://taskfile.dev) のドキュメントを一読することを強く推奨します。Wails CLI に組み込まれている Taskfile のバージョンは、`wails3 task --version` を実行して確認できます。

@end

## 開発モード

Wails のビルドシステムには、ライブリロードとホットモジュールリプレースメントを提供し、開発者体験を向上させる強力な開発モードが用意されています。このモードは、`wails3 dev` コマンドを使用して有効にします。

### 仕組み

`wails3 dev` を実行すると、次の処理が行われます。

1. コマンドは利用可能なポートを確認し、ポートが指定されていない場合は 9245 をデフォルトとして使用します。
2. フロントエンド開発サーバー（Vite）用の環境変数を設定します。
3. [refresh](https://github.com/atterpac/refresh) ライブラリを使用してファイル監視を開始します。

[refresh](https://github.com/atterpac/refresh) ライブラリは、ファイルの変更を監視し、再ビルドをトリガーします。`./build/config.yml` ファイルの `dev_mode` キーで定義された設定を使用します。特定のディレクトリやファイルを無視するように設定できるほか、監視対象のファイルや、変更を検出したときに実行するアクションも指定できます。デフォルト設定で十分に機能しますが、必要に応じて自由にカスタマイズしてください。

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

- ファイル監視のルートパスの設定
- ログレベルの設定
- ファイル変更イベントのデバウンス時間の設定
- 特定のディレクトリ、ファイル、またはファイル拡張子の除外
- ファイル変更時に実行するコマンドの定義

### 開発モードのカスタマイズ

`config.yml` ファイル内のこれらの値を変更することで、開発モードの動作をカスタマイズできます。

カスタマイズ方法には、次のようなものがあります。

1. 監視するディレクトリまたはファイルの変更
2. 変更に対するシステムの応答速度を制御するためのデバウンス時間の調整
3. プロジェクトのニーズに合わせた実行コマンドの追加または変更

### 開発でのブラウザーの使用

Wails v2 では、開発時のブラウザー使用が完全にサポートされていましたが、多くの混乱を招きました。すべてのブラウザー API が WebView で利用できるわけではないため、ブラウザーでは動作するアプリケーションがデスクトップアプリケーションでも必ず動作するとは限りません。

UI を中心とした開発作業では、引き続き v3 でブラウザーを使用できます。開発モードで `http://localhost:9245` の Vite URL にアクセスしてください。これにより、スタイルやレイアウトの作業中に強力なブラウザー開発ツールを利用できます。ただし、このモードでは Go バインディングが *動作しません*。バインディングやイベントなどの機能をテストする準備ができたら、デスクトップ表示に切り替えて、本番環境ですべてが正しく動作することを確認してください。
