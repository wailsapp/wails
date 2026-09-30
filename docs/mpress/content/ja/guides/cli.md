---
title: "CLI リファレンス"
description: "Wails CLI コマンドの完全なリファレンス"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

Wails CLI には、Wails アプリケーションの開発、ビルド、保守に役立つ包括的なコマンドセットが用意されています。

## コアコマンド

コアコマンドは、プロジェクトの作成、開発、ビルドに使用する主要なコマンドです。

すべての CLI コマンドは、`wails3 <command>` という形式です。

### `init`

新しい Wails プロジェクトを初期化します。この初期化中に `go mod tidy` コマンドが実行され、プロジェクトのパッケージが最新の状態になります。`init` コマンドで `-skipgomodtidy` フラグを使用すると、この処理を省略できます。

```bash
wails3 init [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-p` | Go パッケージ名 | `main` |
| `-t` | テンプレート名または URL | `vanilla` |
| `-n` | プロジェクト名 |  |
| `-d` | プロジェクトディレクトリ | `.` |
| `-q` | 出力を抑制 | `false` |
| `-l` | テンプレートを一覧表示 | `false` |
| `-mod` | Go モジュールパス（省略した場合は `-git` から算出） |  |
| `-git` | Git リポジトリの URL |  |
| `-s` | リモートテンプレートを使用する際の警告を省略 | `false` |
| `-productname` | 製品名 | `My Product` |
| `-productdescription` | 製品の説明 | `My Product Description` |
| `-productversion` | 製品バージョン | `0.1.0` |
| `-productcompany` | 会社名 | `My Company` |
| `-productcopyright` | 著作権表示 | `© now, My Company` |
| `-productcomments` | ファイルのコメント | `This is a comment` |
| `-productidentifier` | 製品識別子 |  |
| `-skipgomodtidy` | go mod tidy を省略 | `false` |

`-git` フラグでは、さまざまな形式の Git URL を指定できます。

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project` または `ssh://git@github.com/username/project`
- Git プロトコル：`git://github.com/username/project`
- ファイルシステム：`file:///path/to/project.git`

このフラグを指定すると、次の処理が行われます。

1. プロジェクトディレクトリで Git リポジトリを初期化する
2. 指定した URL をリモートの origin に設定する
3. `go.mod` 内のモジュール名をリポジトリの URL に合わせて更新する
4. すべてのファイルを追加する

### `dev`

アプリケーションを開発モードで実行します。フロントエンドコードをリアルタイムで確認でき、アプリケーション全体を再ビルドせずに変更を実行中のアプリケーションへ反映できます。Go コードの変更も検出され、アプリケーションは自動的に再ビルドされて再起動します。

```bash
wails3 dev [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-config` | 設定ファイルのパス | `./build/config.yml` |
| `-port` | Vite 開発サーバーのポート | `9245` |
| `-s` | HTTPS を有効化 | `false` |

@note{type="info"}
これは `wails3 task dev` の実行と同等であり、プロジェクトのメイン Taskfile にある `dev` タスクを実行します。`Taskfile.yml` ファイルを編集すると、この動作をカスタマイズできます。

@end

### `build`

アプリケーションのデバッグ版をビルドします。デフォルトでは、現在のプラットフォームおよびアーキテクチャ向けにビルドされます。

```bash
wails3 build [flags] [CLI variables...]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-tags` | 追加の Go ビルドタグ（カンマ区切り） |  |

CLI 変数を渡してビルドをカスタマイズできます。

```bash
wails3 build PLATFORM=linux CONFIG=production
```

カスタム Go ビルドタグを渡すには、`-tags` フラグを使用します。

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

タグは `EXTRA_TAGS` として基盤となる Taskfile に転送されます。

@note{type="info"}
これは `wails3 task build` の実行と同等で、プロジェクトのメイン Taskfile にある `build` タスクを実行します。`build` に渡した CLI 変数はすべて基盤となるタスクに転送されます。`Taskfile.yml` ファイルを編集すると、ビルドプロセスをカスタマイズできます。

@end

### `package`

配布用のプラットフォーム固有パッケージを作成します。

```bash
wails3 package [CLI variables...]
```

CLI 変数を渡してパッケージ作成をカスタマイズできます。

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

#### パッケージの種類

各プラットフォームでは、次の種類のパッケージを利用できます。

| プラットフォーム | パッケージの種類 |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
これは `wails3 task package` と同等で、プロジェクトのメイン Taskfile にある `package` タスクを実行します。`package` に渡した CLI 変数はすべて基盤となるタスクに転送されます。`Taskfile.yml` ファイルを編集すると、パッケージ作成プロセスをカスタマイズできます。

@end

### `task`

プロジェクトの Taskfile.yml に定義されたタスクを実行します。これは [Taskfile](https://taskfile.dev) の組み込み版で、ビルド、テスト、デプロイ用のカスタムタスクを定義して実行できます。

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI 変数

`KEY=VALUE` 形式でタスクに変数を渡せます。

```bash
wails3 task build PLATFORM=linux CONFIG=production
wails3 task deploy ENV=staging VERSION=1.2.3
```

これらの変数には、Taskfile.yml 内から Go テンプレート構文を使用してアクセスできます。

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-h` | Task の使用方法を表示します | `false` |
| `-i` | 新しい Taskfile.yml を作成します | `false` |
| `-list` | 説明のあるタスクのみを一覧表示します | `false` |
| `-list-all` | すべてのタスクを一覧表示します（説明の有無を問いません） | `false` |
| `-json` | タスク一覧を JSON 形式で出力します | `false` |
| `-status` | タスクが最新でない場合は、ゼロ以外の終了コードで終了します | `false` |
| `-f` | タスクが最新でも強制的に実行します | `false` |
| `-w` | 指定したタスクの監視モードを有効にします | `false` |
| `-v` | 詳細出力モードを有効にします | `false` |
| `-version` | Task のバージョンを表示します | `false` |
| `-s` | コマンドのエコー表示を無効にします | `false` |
| `-p` | タスクを並列実行します | `false` |
| `-dry` | タスクを実行せずにコンパイルして表示します | `false` |
| `-summary` | タスクの概要を表示します | `false` |
| `-x` | タスクの終了コードをそのまま返します | `false` |
| `-dir` | 実行ディレクトリを設定します |  |
| `-taskfile` | 実行する Taskfile を選択します |  |
| `-output` | 出力形式を設定します：[interleaved|group|prefixed] |  |
| `-c` | カラー出力（デフォルトで有効） | `true` |
| `-C` | 同時に実行するタスク数を制限します |  |
| `-interval` | 変更を監視する間隔（秒） |  |

#### 例

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build PLATFORM=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

エージェント支援によるプロジェクト管理用の Wails プロジェクト MCP サーバーを起動します。これは、実行中のアプリケーションに組み込まれた MCP サーバーとは別のものです。`wails3 mcp`はプロジェクトファイルとライフサイクルコマンドを管理し、アプリケーションの MCP サーバーは実行中の WebView を制御します。

```bash
wails3 mcp [flags]
```

トランスポートは自動的に選択されます：

- MCP ホストがパイプ接続された標準入力／標準出力を使用して Wails を起動した場合、サーバーは<strong>stdio</strong>を使用します。
- ターミナルで対話的に実行した場合、サーバーは`127.0.0.1`上の<strong>Streamable HTTP</strong>を使用し、空きポートをオペレーティングシステムに要求します。

トランスポートを明示的に選択するには、`--stdio`または`--http`を使用します。HTTP モードで空いているループバックポートを選択するには、`--port 0`を使用します。

#### MCP フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `--root` | 許可するプロジェクトルート。この範囲外のパスとシンボリックリンクは拒否されます。 | 現在のディレクトリ |
| `--token` | 変更操作およびプロセス制御ツール用のセッション／Bearer トークン。指定されていない場合は`WAILS_MCP_TOKEN`を使用します。 | 安全に生成 |
| `--stdio` | stdio トランスポートを強制的に使用します。 | 自動 |
| `--http` | Streamable HTTP トランスポートを強制的に使用します。 | 自動 |
| `--port` | HTTP ポート。`0`を指定すると、空いているループバックポートが選択されます。 | `0` |

HTTP モードでは、Wails はエンドポイントと Bearer トークンを標準エラー出力に出力します。stdio モードでは、トークンは MCP の初期化指示に含まれます。サーバーは任意のシェルコマンドを実行する機能を公開しません。リモートテンプレートと Git リモートを使用するには、ツールの`allowExternal`入力による明示的な承認が必要です。

### `doctor`

システムチェックを実行し、ステータスレポートを表示します。

```bash
wails3 doctor
```

## 生成コマンド

生成コマンドを使用すると、バインディング、アイコン、ビルドファイルなど、さまざまなプロジェクトアセットを作成できます。すべての生成コマンドは、ベースコマンド`wails3 generate <command>`を使用します。

### `generate bindings`

Go コード用のバインディングとモデルを生成します。

```bash
wails3 generate bindings [flags] [patterns...]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-f` | 追加の Go ビルドフラグ |  |
| `-d` | 出力ディレクトリ | `frontend/bindings` |
| `-models` | モデルのファイル名 | `models` |
| `-index` | インデックスのファイル名 | `index` |
| `-ts` | TypeScript を生成します | `false` |
| `-i` | TypeScript インターフェースを使用します | `false` |
| `-b` | 同梱のランタイムを使用します | `false` |
| `-names` | ID の代わりに名前を使用します | `false` |
| `-noindex` | インデックスファイルを生成しません | `false` |
| `-noevents` | イベント関連のバインディングの生成をスキップ | `false` |
| `-dry` | ドライラン | `false` |
| `-silent` | サイレントモード | `false` |
| `-v` | デバッグ出力 | `false` |
| `-clean` | 生成前に出力ディレクトリをクリーンアップ | `true` |

### `generate build-assets`

アプリケーション用のビルドアセットを生成します。

```bash
wails3 generate build-assets [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-name` | プロジェクト名 |  |
| `-dir` | 出力ディレクトリ | `build` |
| `-silent` | 出力を抑制 | `false` |
| `-company` | 会社名 |  |
| `-productname` | 製品名 |  |
| `-description` | 製品の説明 |  |
| `-version` | 製品バージョン |  |
| `-identifier` | 製品識別子 | `com.wails.[name]` |
| `-copyright` | 著作権表示 |  |
| `-comments` | ファイルのコメント |  |

### `generate icons`

アプリケーションアイコンを生成します。

```bash
wails3 generate icons [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-input` | 入力PNGファイル | 必須 |
| `-windowsfilename` | Windows用の出力ファイル名 |  |
| `-macfilename` | macOS用の出力ファイル名 |  |
| `-sizes` | アイコンサイズ（カンマ区切り） | `256,128,64,48,32,16` |
| `-example` | サンプルアイコンを生成 | `false` |
| `-iconcomposerinput` | 入力Icon Composerファイル（`.icon`） |  |
| `-macassetdir` | Mac用アセット（Assets.car + icns）の出力ディレクトリ |  |

#### Icon Composer（macOS）

macOS 26以降では、Icon Composerの`.icon`ファイルを使用して`Assets.car`と`icons.icns`を生成できます：

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Appleの`actool`コマンドを使用して`.icon`ファイルをコンパイルします。`actool`のバージョンが26以降であるXcodeが必要です。

Icon Composerを使用する場合は、`build/config.yml`内の`cfBundleIconName`を`.icon`のファイル名（拡張子なし）と一致するように設定します：

```yaml
info:
  cfBundleIconName: "appicon"
```

未設定で`Assets.car`が存在する場合、デフォルトは`"appicon"`です。

### `generate syso`

Windows用の.sysoファイルを生成します。

```bash
wails3 generate syso [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-manifest` | マニフェストファイルへのパス | 必須 |
| `-icon` | アイコンファイルへのパス | 必須 |
| `-info` | バージョン情報ファイルへのパス |  |
| `-arch` | ターゲットアーキテクチャ | 現在のGOARCH |
| `-out` | 出力ファイル名 | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Linux用の.desktopファイルを生成します。

```bash
wails3 generate .desktop [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-name` | アプリケーション名 | 必須 |
| `-exec` | 実行ファイルへのパス | 必須 |
| `-icon` | アイコンへのパス |  |
| `-categories` | アプリケーションのカテゴリ | `Utility` |
| `-comment` | アプリケーションのコメント |  |
| `-terminal` | ターミナルで実行 | `false` |
| `-keywords` | 検索キーワード |  |
| `-version` | アプリケーションのバージョン |  |
| `-genericname` | 汎用名 |  |
| `-startupnotify` | 起動通知を表示 | `false` |
| `-mimetype` | サポートするMIMEタイプ |  |
| `-output` | 出力ファイル名 | `[name].desktop` |

### `generate runtime`

ビルド済みのランタイムを生成します。

```bash
wails3 generate runtime
```

### `generate constants`

GoコードからJavaScript定数を生成します。

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

配布用のWindows WebView2ブートストラップインストーラーを生成します。

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

新しいプロジェクトテンプレートのディレクトリを生成します。

```bash
wails3 generate template [flags]
```

### `generate appimage`

Linux用のAppImageを生成します。

```bash
wails3 generate appimage [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-binary` | バイナリへのパス | 必須 |
| `-icon` | アイコンファイルへのパス | 必須 |
| `-desktop` | .desktop ファイルへのパス | 必須 |
| `-builddir` | ビルドディレクトリ | 一時ディレクトリ |
| `-output` | 出力ディレクトリ | `.` |

## サービスコマンド

サービスコマンドは、Wails サービスの管理に使用します。すべてのサービスコマンドでは、ベースコマンド `wails3 service <command>` を使用します。

### `service init`

新しいサービスを初期化します。

```bash
wails3 service init [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-n` | サービス名 | `example_service` |
| `-d` | サービスの説明 | `Example service` |
| `-p` | パッケージ名 |  |
| `-o` | 出力ディレクトリ | `.` |
| `-q` | 出力を抑制 | `false` |
| `-a` | 作成者名 |  |
| `-v` | バージョン |  |
| `-w` | Web サイトの URL |  |
| `-r` | リポジトリの URL |  |
| `-l` | ライセンス |  |

## ツールコマンド

ツールコマンドは、開発とデバッグに役立つユーティリティを提供します。すべてのツールコマンドでは、ベースコマンド `wails3 tool <command>` を使用します。

### `tool checkport`

ポートが開いているかどうかを確認します。vite が実行中かどうかをテストする際に便利です。

```bash
wails3 tool checkport [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-port` | 確認するポート | `9245` |
| `-host` | 確認するホスト | `localhost` |

### `tool watcher`

ファイルを監視し、変更されたときにコマンドを実行します。

```bash
wails3 tool watcher [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-config` | 設定ファイルのパス | `./build/config.yml` |
| `-ignore` | 無視するパターン |  |
| `-include` | 含めるパターン |  |

### `tool cp`

ファイルをコピーします。

```bash
wails3 tool cp
```

### `tool buildinfo`

アプリケーションのビルド情報を表示します。

```bash
wails3 tool buildinfo
```

### `tool version`

指定されたフラグに基づいてセマンティックバージョンを更新します。

```bash
wails3 tool version [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-v` | 更新する現在のバージョン |  |
| `-major` | メジャーバージョンを上げる | `false` |
| `-minor` | マイナーバージョンを上げる | `false` |
| `-patch` | パッチバージョンを上げる | `false` |
| `-prerelease` | プレリリースバージョンを上げる（例：alpha.5 から alpha.6） | `false` |

このコマンドは、major > minor > patch > prerelease の優先順位に従います。入力バージョンに「v」プレフィックスがある場合は、プレリリースおよびメタデータの各要素とともに保持します。

使用例：

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.0.0-alpha.5 -prerelease # Output: v3.0.0-alpha.6
```

### `tool package`

Linux パッケージ（deb、rpm、archlinux）を生成します。

```bash
wails3 tool package [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-format` | パッケージ形式（deb、rpm、archlinux） | `deb` |
| `-name` | 実行ファイル名 | `myapp` |
| `-config` | 設定ファイルのパス |  |
| `-out` | 出力ディレクトリ | `.` |

### `tool lipo`

アーキテクチャ固有のバイナリを結合して、macOS ユニバーサルバイナリを作成します。

```bash
wails3 tool lipo [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-output` | 出力バイナリのパス |  |

### `tool capabilities`

システムのビルド機能（Linux で GTK4／GTK3 が利用可能かどうか）を確認します。

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

クロスコンパイル用の Docker ボリュームマウントフラグを生成します。Go モジュールキャッシュと、`go.mod` 内のローカルな `replace` ディレクティブを対象とする `-v` フラグを出力し、Taskfile の `docker run` コマンドで使用できるようにします。

```bash
wails3 tool docker-mounts
```

### `tool has`

ツールまたは機能が利用可能かどうかを確認し、標準出力に `true` または `false` を出力します。`command -v` に代わるクロスプラットフォーム対応の手段として、Taskfile の `sh:` 変数で使用することを想定しています。

複数の候補のうち、いずれか1つが利用可能かどうかを確認するには、`|` を使用します。

```bash
wails3 tool has <tool>
```

#### 例

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Taskfile での使用方法

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="非推奨"}
`wails3 tool has-cc` は非推奨です。代わりに `wails3 tool has gcc|clang` を使用するように Taskfile を更新してください。

@end

`wails3 tool has gcc|clang` の後方互換エイリアスです。`gcc` または `clang` が PATH で利用可能かどうかを確認し、`true` または `false` を出力します。

```bash
wails3 tool has-cc
```

## 更新コマンド

更新コマンドは、プロジェクトアセットの管理と更新に役立ちます。すべての更新コマンドでは、ベースコマンドとして `wails3 update <command>` を使用します。

### `update cli`

Wails CLI を新しいバージョンに更新します。

```bash
wails3 update cli [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-pre` | 最新のプレリリースに更新する | `false` |
| `-version` | 指定したバージョンに更新する |  |
| `-nocolour` | 色付き出力を無効にする | `false` |

update cli コマンドを使用すると、インストール済みの Wails CLI を更新できます。デフォルトでは、最新の安定版リリースに更新します。 最新のプレリリースバージョンに更新するには `-pre` フラグを使用します。また、`-version` フラグを使用して特定のバージョンを指定することもできます。

更新後は、プロジェクトの go.mod ファイルでも同じバージョンを使用するように更新してください：

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

指定した設定ファイルを使用してビルドアセットを更新します。

```bash
wails3 update build-assets [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-config` | 設定ファイルのパス |  |
| `-dir` | 出力ディレクトリ | `build` |
| `-silent` | 出力を抑制 | `false` |
| `-company` | 会社名 |  |
| `-productname` | 製品名 |  |
| `-description` | 製品の説明 |  |
| `-version` | 製品バージョン |  |
| `-identifier` | 製品識別子 |  |
| `-copyright` | 著作権表示 |  |
| `-comments` | ファイルのコメント |  |

## ユーティリティコマンド

ユーティリティコマンドは、一般的なタスクに便利なショートカットを提供します。これらのコマンドは、ベースコマンド `wails3 <command>` で直接使用します。

### `docs`

デフォルトのブラウザーで Wails のドキュメントを開きます。

```bash
wails3 docs
```

### `releasenotes`

現在のバージョンまたは指定したバージョンのリリースノートを表示します。

```bash
wails3 releasenotes [flags]
```

#### フラグ

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-v` | リリースノートを表示するバージョン |  |
| `-n` | カラー出力を無効化 | `false` |

### `version`

Wails の現在のバージョンを出力します。

```bash
wails3 version
```

### `sponsor`

デフォルトのブラウザーで Wails のスポンサーシップページを開きます。

```bash
wails3 sponsor

```
