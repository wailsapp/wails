---
title: "YAML ビルド"
description: "1 つの wails.yaml ファイルで、実験的な Wails v3 ネイティブビルドシステムを構成して実行します。"
sidebar:
  order: 2
  badge:
    text: "実験的機能"
    variant: caution
---

@note{type="caution" title="実験的プレビュー"}
最初のリリースで受け付けるのはリテラルによる構成のみです。この実験では、リリースへの提案前に形式を評価できるよう、実行時間が制限されたライフサイクルフックのプロトタイプも提供します。カスタムステージ、YAML 式、型付きの Wails ツール呼び出しは引き続き先送りされています。既存のプロジェクトでは、レビュー済みの移行を明示的に有効化するまで Taskfile を使い続けます。
@end

## YAML ビルドを有効にする

実験的なコマンドを実行するすべてのターミナルで `WAILS_EXP_USE_WAKE` を設定します。Bash または Zsh の場合:

```bash
export WAILS_EXP_USE_WAKE=1
```

PowerShell の場合:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

この変数は、存在するだけでオプトインになる方式です。空の値や `false` でも有効になります。実験を無効にするには、変数の設定を解除します（Bash/Zsh では `unset WAILS_EXP_USE_WAKE`、PowerShell では `Remove-Item Env:WAILS_EXP_USE_WAKE`）。`wails.yaml` が有効になると、コマンドの実行には Wake の再有効化が必要になり、Taskfile にフォールバックすることはありません。
この変数がない場合、YAML による build/dev のルーティングと YAML プロジェクトの生成は無効になり、`migrate`、`config`、`eject` の各コマンドは使用できません。移行していないプロジェクトでは、引き続き別個の組み込み Task ランタイムを使用します。

## プロジェクトを構成する

オプトインを有効にして、プロジェクトルートに `wails.yaml` を追加します。このファイルが存在するとネイティブビルドシステムが選択され、Wails はそのプロジェクト内のすべての Taskfile を無視します。以下のすべてのコマンドを実行する際も、オプトインを有効にしておいてください。スキーマのモードラインにより、YAML 言語サーバーをサポートするエディターで補完とインライン検証が有効になります。

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

デフォルトのホスト向けビルド、または名前付きの完全なリリース要求を実行します:

```bash
wails3 build
wails3 build release
```

一度限りのビルドでは、プロファイルを追加せずに、カンマ区切りでターゲットと形式を選択します:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## デスクトップアプリケーションをクロスコンパイルする

実験的な pure-Go デスクトップバックエンドを使用するには、ビルド環境で `CGO_ENABLED: "0"` を設定します。Wake は macOS および Linux の GUI プログラムに `purego` タグを選択し、Windows、macOS、Linux の各ホストで Go を直接使用します。このモードでデスクトップ実行可能ファイルをビルドする場合、Docker、エミュレーション、C コンパイラー、ターゲットプラットフォームの SDK は不要です。アプリケーションの依存関係も `CGO_ENABLED=0` をサポートしている必要があります。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

2 番目のコマンドは両方の macOS アーキテクチャを統合し、Go を使用して未署名の DMG を作成します。サポート対象の任意のデスクトップホストで実行できます。コード署名と Apple の公証は、引き続き macOS で別途行う処理です。

解決後のモードは、実行可能ファイルごとに個別に選択されます。宣言された `build.environment.CGO_ENABLED`、ターゲット環境、またはバイナリ環境は、コマンドを起動したターミナルの値より優先されます。それ以外の場合、明示的に選択された `purego` タグ、CLI/サーバープログラム、または Windows ターゲットでは、デフォルトで `0` が使用されます。macOS および Linux の GUI プログラムでは、デフォルトで `1` が使用されます。`--plan --json` は、各実行可能ファイルについて、解決されたコンパイラー、CGo モード、タグを報告します。クロスコンパイルによってターゲット用ファイルは作成されますが、互換性のないデスクトップターゲットを実行するには、引き続き適切なホストが必要です。

`CGO_ENABLED=1` を使用すると、Wake はまず互換性のあるネイティブツールチェーンまたは構成済みの `CC` クロスコンパイラーを使用し、次に利用可能であれば Windows のクロスビルドに Zig を使用し、その後、サポートされている場合は既存の Docker/Podman フォールバックを使用します。明示的な `CC` では、ターゲット用のコンパイラー、ライブラリ、sysroot を指定する必要があります。Zig だけでは Linux GTK/WebKit の sysroot は提供されません。以下で説明するように、パッケージ化と署名には追加のツールが必要になる場合があります。モバイルターゲットには CGo と各プラットフォームの SDK が必要です。

Go のみのコンパイルでは、ターゲット環境のシステムライブラリは同梱されません。
GUI コンパニオン実行可能ファイルを含む Linux GUI アプリケーションでは、ビルドタグで選択された GTK/WebKitGTK ランタイムが引き続き必要です。パッケージの依存関係または適切なデプロイ環境から、このランタイムを提供する必要があります。CLI/サーバー専用パッケージには、デフォルトではこれらの GUI 依存関係は追加されません。依存関係を明示的に上書きすることは引き続き可能です。

## 複数の実行可能ファイルをビルドする

プロジェクトではプライマリープログラムを定義します。その Go パッケージを `project.package`（デフォルトは `.`）に設定し、`project.kind` を `gui`、`cli`、または `server`（デフォルトは `gui`）に設定します。追加のプログラムには、`binaries` 以下の名前付きエントリーを使用します:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

各プログラムには、それぞれ独自のコンパイル入力とキャッシュ入力が割り当てられます。インポートされたプロジェクトパッケージと埋め込みファイルは追跡されます。独立したプログラムを編集しても、ほかのプログラムは再ビルドされません。追加の実行可能ファイルはデスクトップでサポートされ、デフォルトのアプリバンドルとインストーラーに含まれます。ユーザーが所有する置換用パッケージファイルでは、その内容をユーザー自身が管理します。

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

名前付きプロファイルは完全な要求です。これらの選択オーバーライドは匿名ビルドで使用してください。`--jobs` は実行ワーカー数の上限を設定します。CLI/サーバープログラムのみを選択すると、フロントエンド処理は省略されます。

アプリケーションコードを実行せずにバージョンを読み取るには、`project.version_file` と `project.version_symbol` の両方を構成します。ファイルには静的な Go 文字列宣言が含まれている必要があります。これにより、リテラルの `project.version` を省略できます。リポジトリのリビジョンを実行可能ファイルの文字列変数へ注入するには、プロジェクトまたは名前付きバイナリに `commit_symbol: main.Commit` を設定します。コミットの注入を要求したときにリポジトリのメタデータがない場合はエラーになります。

`upx_level: 1` から `9` の値を指定すると、リリース用の Linux および Windows 実行可能ファイルを必要に応じて圧縮できます。ゼロを指定すると無効になります。UPX がインストールされている必要があります。開発用ビルドと macOS ビルドでは圧縮を省略します。コンパイルまたは圧縮に失敗した場合、以前に公開されたプログラムは保持されます。

## 配布用にパッケージ化して署名する

未署名パッケージの作成に必要なホストとツールは、署名の場合とは異なります:

| 出力 | 作成に使用できるホストと必要なツール | 署名 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS、Linux。Go ライターを使用し、Windows SDK は不要です。 | Windows ホストと署名ツール。 |
| Windows NSIS インストーラー | Windows、macOS、Linux。ホストネイティブの `makensis` が必要です。 | Windows ホストと署名ツール。 |
| macOS アプリと DMG | Windows、macOS、Linux。Go によるアセンブリ、ユニバーサルバイナリの統合、ポータブル DMG ライターを使用します。 | Apple のツールと資格情報を備えた macOS。 |
| Linux DEB、RPM、Arch、Void | Windows、macOS、Linux。ネイティブの Go パッケージライターを使用します。 | DEB/RPM はそれぞれ外部の署名ツールを使用します。Arch は以下の Go 署名機能を使用します。Void の出力は未署名です。 |
| Linux AppImage | ターゲットアーキテクチャと一致する Linux ホスト。外部の linuxdeploy と GTK デプロイプラグインが必要です。 | AppImage には組み込みの署名フローがありません。 |

Windows MSIX アーカイブは Go で直接作成されます。署名なしパッケージの作成に Windows SDK や Docker は必要ありません。ただし、Windows での署名には引き続き Windows ホストとその署名ツールが必要です。NSIS は外部の NSIS コンパイラを使用し、`packages.nsis.defines` を通じてリテラル定義を受け付けます。macOS では、継承されたロケールが UTF-8 でない場合に限り、Wails は NSIS コンパイラの子プロセスにのみ UTF-8 ロケールを設定します。ターミナルとアプリケーションの環境は変更されません。標準の NSIS インストーラーは、パッケージ化されたすべての実行可能ファイルが CLI またはサーバーである場合、WebView2 ブートストラッパーを含めません。プライマリ実行可能ファイルまたはコンパニオン実行可能ファイルが GUI の場合は、WebView2 が維持されます。カスタム NSIS スクリプトでは、独自のランタイムインストールロジックとの互換性を保つため、ブートストラッパーが維持されます。

Linux は AppImage、DEB、RPM、Arch Linux、Void の各パッケージをサポートします。AppImage では、外部の linuxdeploy と GTK プラグインを使用するという境界が維持されます。zsync ファイルを要求するには `packages.appimage.update_info` を、SHA256 ファイルを要求するには `packages.appimage.checksum: true` を設定します。サイドカーファイルはパッケージキャッシュから復元され、AppImage とともにアーティファクトレシートに記載されます。現在、AppImage のパッケージ化には Linux と、ホストとターゲットで一致するアーキテクチャが必要です。GTK の検出では、pure-Go GUI によって動的に読み込まれるライブラリを検出できないため、このアダプターはまだ pure-Go GUI のデプロイをサポートしていません。そのビルドモードでは、Linux のパッケージマネージャー形式のいずれかを使用してください。

Arch パッケージでは、プレリリースがネイティブのパッケージメタデータに保持されます。たとえば、`1.2.3-beta.2` は `1.2.3beta.2-1` になり、最終リリースより前に並びます。Void の依存関係では、`gtk4` のようなパッケージ名だけの指定を受け付け、`gtk4>=0` に正規化します。明示的なバージョン比較と完全な `name-version_revision` 値は保持されます。ライブラリのファイル名ではなく、Void のパッケージ名を使用してください。

Arch パッケージと分離署名を作成するには、エクスポートされた ASCII Armor 形式の PGP 秘密鍵ファイルを使用し、必要に応じてパスフレーズを格納した環境変数の名前を指定します:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

ローカル環境または CI のシークレットストアで指定した環境変数を設定してから、`wails3 build release` を実行します。暗号化されていない鍵の場合は、`credential` を省略します。Go 署名ツールはサポート対象のすべてのデスクトップホストで動作し、GPG を必要としません。元の `.pkg.tar.zst` と、同じディレクトリの `.pkg.tar.zst.sig` を公開し、両方をアーティファクトレシートに個別に記録します。鍵の内容がステージングされたり、ビルドキャッシュに保存されたりすることはありません。署名とその公開結果がキャッシュから再利用されることもありません。署名に失敗した場合は以前に公開されたファイルの組み合わせが保持され、いずれかのファイルのインストールに失敗した場合は両方の公開がロールバックされます。ただし、これによって 2 回のファイルシステム上の名前変更が、同時アクセスする読み取り側に対してアトミックになるわけではありません。

macOS アプリバンドルには付属プログラムが含まれます。`darwin.microphone_usage_description` を設定すると、用途説明文に加え、カスタムの entitlements ファイルが指定されていない場合は音声入力の entitlement も追加されます。カスタムの plist および entitlements のソースは変更されません。署名と公証には、macOS 上の既存の Apple ツールおよび認証情報が使用されます。公証には実際の notary profile が必要です。署名済みの出力には、`notes.signed.app` や `notes.signed.dmg` など、プラットフォーム固有の拡張子が維持されます。

## 開発時の操作

`wails3 dev` の実行中に、コマンドを入力して Enter キーを押します。

| コマンド | 動作 |
| --- | --- |
| `r` | バックエンドに変更がない場合も含め、強制的に再ビルドして再起動します。 |
| `b` | バインディングを再生成して再ビルドします。 |
| `f` | フロントエンドを再起動します。 |
| `c` | ターミナルをクリアします。 |
| `q` | セッションとその子プロセスを停止します。 |
| `h` | 操作方法を表示します。 |

ビルドに失敗しても、現在のアプリケーションは実行を継続します。これらの操作には対話型ターミナルが必要です。入力をパイプしてもセッションの動作は変わりません。

## 実行前に処理内容を確認する

プランの確認では、実行時と同じ不変グラフが解決されます。テキストは人が読むためのものであり、JSON はバージョン管理された機械可読の契約です。

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## ビルドせずに設定を検証する

CI やエディタータスク、または設定をコミットする前に、専用の事前検証コマンドを使用します。プロファイルを指定しない場合は、ベースマニフェストと名前付きプロファイルをすべて検証します。プロファイルを指定した場合は、その完全なリクエストだけを検証します。

```bash
wails3 config check
wails3 config check release
```

構文、スキーマ、セマンティクス、パス、ターゲット、形式、署名、プランナー、および設定に起因するホスト機能のエラーには、Build および Dev と同じ診断レンダラーが使用されます。診断にはフィールドとソース範囲が示され、ソース行とキャレットが表示されます。有用な場合は修正のヒントも含まれます。診断のソース抜粋では環境変数の値が伏せられます。コマンドラインの選択だけが原因のホストエラーは、無効な YAML として誤って報告されず、通常のコマンドエラーとして扱われます。

## プロジェクト所有のアセットを使用する

マニフェストのパスは、アイコン、プラットフォームマニフェスト、パッケージテンプレート、メンテナースクリプトなど、ユーザーが所有するファイルを参照します。Wails がそれらのソースを編集することはありません。ツールが使用する前に、一時的な `.wails/` ワークスペースへコピーします。最終アーティファクトは `build.output` 配下に配置され、検証済みのレシートが `.wails/artifacts/receipt.json` に書き込まれます。

## Android にデプロイする

本番用 Android プロファイルは AAB を生成します。ローカルデプロイでは、Wails が ABI 固有のデバッグ APK を内部でビルドし、ADB でインストールしてアプリを起動します。

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

セレクターを明示的に指定しない場合、対話型ターミナルには接続済みデバイスと設定済みの Android Virtual Device が表示され、選択できます。複数のターゲットを利用できる環境で自動化する場合は、`--device` または `--emulator` を指定する必要があります。

アプリケーション単位のログストリームを接続し、アプリが終了するか Ctrl+C を押すまでターゲットを監視します。

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

接続実行ではアプリケーションの再起動後も追跡を続け、ターゲットの切断やオフラインへの移行を報告します。素早く反復作業できるよう、デフォルトではエミュレーターを実行したままにします。`--stop-emulator` は明示的な指定であり、`--logs` が必要です。接続セッションの終了時に停止するのは、ターゲットがエミュレーターの場合だけです。Wails が新しいエミュレーターを起動したものの、起動準備の確認に失敗した場合やキャンセルされた場合は、そのプロセスを自動的にクリーンアップします。

## 範囲を限定したライフサイクルフックを実行する

試験的なフック機能は、宣言された各フェーズでプロジェクト所有のスクリプトを 1 つ呼び出します。フックはシェルコマンドではありません。Wails はファイルを直接実行し、`directory` が設定されていない場合はプロジェクトルートを作業ディレクトリとして使用します。

Unix ホストでは、スクリプトに実行権限があり、先頭に shebang が必要です。Windows ホストでは、固定されたインタープリターアダプターを介して `.cmd`、`.bat`、`.ps1` のスクリプトを実行できます。どちらのプラットフォームでも、Wails はインラインのシェルコマンドを解析しません。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

サポートされるフェーズは、`before_build`、`after_build`、`before_package`、`after_package`、`before_sign`、`after_sign` です。ビルドフックは、必要に応じて共有プロジェクト全体で 1 回、またはターゲットごとに 1 回実行されます。パッケージフックと署名フックは、ターゲットに要求されたすべての形式を囲むバリアを形成しますが、バリア内の形式ごとの処理は並行して実行されます。

フックには、Wails が所有する環境変数 `WAILS_HOOK_CONTEXT_FILE` が 1 つ渡されます。この変数は、バージョン管理された呼び出しコンテキストを格納する保護された JSON ファイルを指します。

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

コンテキストには解決済みのパスが含まれますが、継承したプロセス環境やシークレットがコピーされることはありません。Wails は、所有者だけがアクセスできるファイル権限を設定し、`.wails/hooks/context/` 配下にこのファイルをアトミックに作成します。成功時またはキャンセル時には削除されます。失敗時には診断用に保持され、そのパスがエラーに含まれます。ゼロ以外の終了コードにより通常の Plan ステップは失敗し、キャンセルするとフックのプロセスグループが終了します。

デフォルトでは、呼び出しのたびにフックが実行されます。スクリプトの入力と出力をすべて宣言できる場合に限り、キャッシュを有効にしてください。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

スクリプトとその実行モードは暗黙の入力です。複数の出力は、ルートではない 1 つのディレクトリを共有する必要があります。このディレクトリが、Wails によって記録および復元される範囲限定のアーティファクトになります。出力ルートにスクリプトや宣言済みの入力を含めることはできません。また、キャッシュされる出力ツリーに含められるのは通常のファイルとディレクトリだけであり、シンボリックリンクや特殊ファイルは含められません。常時実行フックの下流処理はキャッシュできません。宣言されていない副作用のフィンガープリントを安全に作成できないためです。キャッシュ可能なフックでは、宣言したアーティファクト ID を通じて下流のキャッシュが維持されます。

## 既存のプロジェクトを移行する

無効状態のドラフトを作成して確認します。

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

プロジェクト所有のスクリプトを引数なしで呼び出す、保守的に変換可能なライフサイクルタスクはフックに変換されます。それ以外の到達可能なカスタム Taskfile 動作は、ブロッカーとして報告されます。`wails.migrated.yaml` を確認した後、明示的に有効化します。

```bash
wails3 migrate --activate
```

有効化すると、`wails.yaml` がアトミックに作成されます。Taskfile、`build/config.yml`、スクリプト、アセットが削除、名前変更、編集されることはありません。無効な YAML は YAML エラーとして失敗し、従来のルートへフォールバックすることはありません。

## すべてのデフォルト値を具体化する

解決済みの完全な設定を確認またはコピーするには、無効状態のスナップショットを書き出します。

```bash
wails3 eject
```

これにより `wails.ejected.yaml` が作成されます。有効な設定が書き換えられることはありません。

## 互換性維持期間

`wails3 package` と `wails3 sign` は、同じプランナーを指す非推奨のエイリアスです。新しい自動化では、プロファイルまたは形式を指定して `wails3 build` を使用してください。ロールアウトはリリース単位で段階的に行われます。最初のリリースでは新規プロジェクトに YAML を使用し、次のリリースでは従来の Taskfile ルーティングに警告を出し、その次のリリースではエイリアスとともに Taskfile の実行を削除します。各移行には、そのリリースのネイティブホスト受け入れテストマトリックスがすべて成功している必要があります。

すべてのフィールド、型、デフォルト値、適用条件については、生成された [`wails.yaml`フィールドリファレンス](/reference/wails-yaml-fields/)を参照してください。

[`wails3 run`設定リファレンス](/reference/wails-yaml-run/)では、ローカル起動のデフォルト値、プラットフォーム固有の上書き設定、マニフェストがない場合の `go run .` フォールバックについて説明しています。
