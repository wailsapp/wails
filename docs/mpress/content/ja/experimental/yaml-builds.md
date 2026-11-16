---
title: "YAML ビルド"
description: "1 つの wails.yaml ファイルで、実験的な Wails v3 ネイティブビルドシステムを設定して実行します。"
sidebar:
  order: 2
  badge:
    text: "実験的機能"
    variant: caution
---

@note{type="caution" title="実験的プレビュー"}
最初のリリースで受け付けるのは、リテラルによる設定のみです。この実験では、制限付きライフサイクルフックの形式をリリースに向けて提案する前に評価できるよう、そのプロトタイプも提供します。カスタムステージ、YAML 式、型付き Wails ツール呼び出しへの対応は引き続き保留されています。既存のプロジェクトでは、レビュー済みの移行を明示的に有効化するまで Taskfile を使用し続けます。
@end

Wake を初めて使用する場合は、[入門チュートリアル](/experimental/wake/)から始めてください。
このチュートリアルでは、実験の制限事項、プロジェクトおよびビルドの各設定画面、最初のローカルビルド、レビュー済みの移行について説明します。このページでは、個別のビルド、パッケージ化、構成のワークフローについて、さらに詳しく説明します。

## YAML ビルドを有効にする

実験的なコマンドを実行するすべてのターミナルで、`WAILS_EXP_USE_WAKE` を設定します。Bash または Zsh の場合：

```bash
export WAILS_EXP_USE_WAKE=1
```

PowerShell の場合：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

この変数は存在するだけでオプトインになります。空の値や `false` でも有効になります。実験を無効にするには、設定を解除します（Bash/Zsh では `unset WAILS_EXP_USE_WAKE`、PowerShell では `Remove-Item Env:WAILS_EXP_USE_WAKE`）。`wails.yaml` が有効になると、コマンドを実行するには Wake を再度有効にする必要があり、Taskfile にフォールバックすることはありません。この変数がない場合、YAML による build/dev のルーティングと YAML プロジェクトの生成は無効になり、`migrate`、`config`、`eject` の各コマンドは使用できません。未移行のプロジェクトでは、独立した組み込み Task ランタイムが引き続き使用されます。

## プロジェクトを設定する

ブラウザーでプロジェクトを作成して構成するには、オプトインを有効にして `wails3 init -ui` を実行し、[Wake チュートリアル](/experimental/wake/)に従ってください。プロジェクトの選択内容は `wails.yaml` に保存されます。既存のプロジェクトを変更するには、このファイルを編集してください。`wails3 setup` では、共有設定、ツール、署名用の認証情報を構成します。

オプトインを有効にした状態で、プロジェクトルートに `wails.yaml` を追加します。このファイルが存在するとネイティブビルドシステムが選択され、Wails はそのプロジェクト内のすべての Taskfile を無視します。以降のすべてのコマンドでも、オプトインを有効にしたままにしてください。スキーマのモードラインにより、YAML 言語サーバーをサポートするエディターで補完とインライン検証を利用できます。

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

デフォルトのホスト向けビルド、または名前付きの完全なリリースリクエストを実行します：

```bash
wails3 build
wails3 build release
```

一度限りのビルドでは、プロファイルを追加せずに、カンマ区切りでターゲットと形式を選択します：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## デスクトップアプリケーションをクロスコンパイルする

実験的な Pure Go デスクトップバックエンドを使用するには、ビルド環境で `CGO_ENABLED: "0"` を設定します。Wake は macOS および Linux の GUI プログラムに `purego` タグを選択し、Windows、macOS、Linux の各ホストで Go を直接使用します。このモードでデスクトップ実行ファイルをビルドする際、Docker、エミュレーション、C コンパイラー、ターゲットプラットフォームの SDK は不要です。アプリケーションの依存関係も `CGO_ENABLED=0` をサポートしている必要があります。

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

2 番目のコマンドは、両方の macOS アーキテクチャを結合し、Go を使用して署名なしの DMG を作成します。サポート対象のどのデスクトップホストでも実行できます。

解決後のモードは、実行ファイルごとに個別に選択されます。宣言された `build.environment.CGO_ENABLED`、ターゲット環境、またはバイナリ環境は、コマンドを呼び出したターミナルの値より優先されます。それ以外の場合、`purego` タグが明示的に選択されているプログラム、CLI／サーバープログラム、または Windows ターゲットでは、デフォルトで `0` が使用されます。macOS および Linux の GUI プログラムでは、デフォルトで `1` が使用されます。`--plan --json` は、各実行ファイルについて解決されたコンパイラー、CGo モード、タグを報告します。クロスコンパイルではターゲット用ファイルが作成されますが、互換性のないデスクトップターゲットを実行するには、引き続き適切なホストが必要です。

`CGO_ENABLED=1` を使用すると、Wake はまず互換性のあるネイティブツールチェーンまたは設定済みの `CC` クロスコンパイラーを使用します。次に、利用可能であれば Windows のクロスビルドに Zig を使用し、その後、サポートされている環境では既存の Docker/Podman フォールバックを使用します。`CC` を明示的に指定する場合は、ターゲットのコンパイラー、ライブラリ、sysroot を用意する必要があります。Zig だけでは Linux GTK/WebKit の sysroot は提供されません。パッケージ化と署名には、後述する追加ツールが必要な場合があります。モバイルターゲットには CGo と各プラットフォームの SDK が必要です。

Go のみでコンパイルしても、出力先のシステムライブラリは同梱されません。GUI のコンパニオン実行ファイルを含む Linux GUI アプリケーションには、ビルドタグで選択された GTK/WebKitGTK ランタイムが引き続き必要です。パッケージの依存関係または適切なデプロイ環境から提供する必要があります。CLI／サーバーのみのパッケージでは、デフォルトでこれらの GUI 依存関係は追加されません。依存関係の明示的なオーバーライドは引き続き使用できます。

## 複数の実行ファイルをビルドする

プロジェクトではプライマリプログラムを定義します。その Go パッケージを `project.package`（デフォルトは `.`）に設定し、`project.kind` を `gui`、`cli`、`server` のいずれか（デフォルトは `gui`）に設定します。追加のプログラムには、`binaries` の下にある名前付きエントリを使用します：

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

各プログラムには、それぞれ独自のコンパイル入力とキャッシュ入力があります。インポートされたプロジェクトパッケージと埋め込みファイルは追跡されるため、独立したプログラムを編集しても他のプログラムは再ビルドされません。追加の実行ファイルはデスクトップでサポートされ、デフォルトのアプリバンドルとインストーラーに含まれます。ユーザーが所有する置換用パッケージファイルの内容は、ユーザー自身で管理します。

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

名前付きプロファイルは完全なリクエストです。以下の選択オーバーライドは匿名ビルドで使用してください。`--jobs` は実行ワーカー数を制限します。CLI／サーバープログラムだけを選択すると、フロントエンド処理は省略されます。

アプリケーションコードを実行せずにバージョンを読み取るには、`project.version_file` と `project.version_symbol` の両方を設定します。ファイルには静的な Go 文字列宣言が含まれている必要があります。その場合、リテラルの `project.version` は省略できます。リポジトリのリビジョンを実行ファイルの文字列変数に注入するには、プロジェクトまたは名前付きバイナリに `commit_symbol: main.Commit` を設定します。コミットの注入を要求したときにリポジトリのメタデータが存在しない場合はエラーになります。

`upx_level: 1` から `9` までの値により、リリース用の Linux および Windows 実行ファイルを必要に応じて圧縮できます。ゼロを指定すると無効になります。UPX がインストールされている必要があります。開発ビルドと macOS ビルドでは圧縮を行いません。コンパイルまたは圧縮に失敗した場合、以前に公開されたプログラムは保持されます。

## 配布用にパッケージ化して署名する

### ビルド済みのデスクトップ実行ファイルを使用する

別のビルドシステムまたは CI ジョブですでに実行ファイルを生成している場合は、`--skip-build` を使用します：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

各実行可能ファイルを `build.output`（デフォルトは `bin`）に配置します。メインの実行可能ファイルには `project.binary_name` を使用し、付随する実行可能ファイルには各 `binaries` エントリの名前を使用します。Windows のファイル名は `.exe` で終わります。複数ターゲットを指定した呼び出しでは、各ターゲットのファイルを `<os>-<arch>` サブディレクトリから読み取ります。macOS GUI の場合、Wake は `bin/<name>` に未加工の実行可能ファイルがあればそこから読み取り、なければ `bin/<name>.app/Contents/MacOS/<name>` から読み取ります。`darwin/universal` の入力には、amd64 と arm64 の両方の実行可能ファイルスライスが含まれている必要があります。

Wake はファイル形式とターゲットアーキテクチャを検証し、生成したワークスペースに実行ファイルをコピーして、パッケージキャッシュ用にその内容を追跡します。実行ファイルの入力が存在しない、互換性がない、またはシンボリックリンクである場合、再ビルドや古い入力の復元は行わず、処理に失敗します。指定した実行ファイルは上書きされず、chmod も実行されません。いずれかの実行ファイルを変更すると、それを使用するパッケージのキャッシュは無効になります。

このモードでは、フロントエンドのインストールとビルド、バインディング、Go のコンパイル、UPX、`before_build`／`after_build` フックを省略します。一方、パッケージアセットの生成、macOS アプリバンドルの組み立て、設定済みのパッケージ化／署名フックの実行は引き続き行います。パッケージ化ツールと署名ツールには、引き続き下表の要件が適用されます。Go とフロントエンドのビルドツールは不要です。パッケージのメタデータは現在のマニフェストから取得され、設定されている場合は `version_file` も含まれます。バイナリからバージョンを推測することはありません。GUI と CLI の種別を含め、ビルド設定とパッケージ設定を一致させてください。

`wails3 package --skip-build` もサポートされています。`wails3 run` を使用する場合も、起動引数と環境設定が適用され、実行ファイルはホスト上で実行できる必要があります。このモードには有効な YAML プロジェクトが必要で、本番用デスクトップターゲットのみをサポートします。開発セッションやモバイルのパッケージ化には適用されません。

署名なしパッケージの作成と署名では、必要なホストとツールが異なります：

IPA ペイロードと macOS 公証用 ZIP は、バンドル内のパス、実行権限、相対シンボリックリンクを維持したまま Go で書き込まれます。どちらの処理にも、外部の `zip` または `ditto` 実行ファイルは必要ありません。これらのアーカイブに含まれるのはアプリケーションファイルであり、Finder メタデータやリソースフォークではありません。これらのメタデータは署名済みバンドルを無効にするため、macOS の署名処理によってすでに削除されます。

| 出力 | 作成に使用できるホストとツール | 署名 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS、Linux。Go ライターを使用し、Windows SDK は不要です。 | サポート対象のどのデスクトップホストでも、エクスポート済み RSA PFX を使用して Go で署名できます。証明書ストアまたはプロバイダーの ID には Windows SignTool が必要です。 |
| Windows NSIS インストーラー | Windows、macOS、Linux。ホストネイティブの `makensis` が必要です。 | サポート対象のどのデスクトップホストでも、エクスポート済み RSA PFX を使用して Go で PE 署名できます。証明書ストアまたはプロバイダーの ID には Windows SignTool が必要です。 |
| macOS アプリと DMG | Windows、macOS、Linux。Go による組み立て、ユニバーサルバイナリの結合、ポータブル DMG ライターを使用します。 | エクスポート済み RSA PKCS#12 ID を使用した Go 署名、またはアドホック署名。キーチェーン ID では引き続き macOS ツールを使用します。API キーによる公証とチケットのステープル処理は Go で実行されます。 |
| Linux DEB、RPM、Arch、Void | Windows、macOS、Linux。ネイティブの Go パッケージライターを使用します。 | エクスポート済み OpenPGP キーでは、ネイティブの DEB/RPM/Arch 署名機能を使用します。従来の DEB/RPM キー ID には外部ツールが必要です。Void の出力は署名されません。 |
| Linux AppImage | Windows、macOS、または Linux。Go によるデプロイ機能と SquashFS ライターに加え、ターゲット用の AppImage ランタイムが必要です。GUI のクロスパッケージには、指定された Linux ランタイムツリーが必要です。 | AppImage の組み込み署名フローはありません。 |

署名済み DMG の場合、通常 Wake はアプリケーションの非公開コピーに署名し、Go ライターがイメージを圧縮してから、そのイメージに署名します。これにより、`hdiutil` を使用した DMG のマウント、展開、再圧縮を回避できます。署名はリクエストのたびに実行され、署名済みパッケージがビルドキャッシュから再利用されることはありません。`after_package` または `before_sign` フックを持つプロジェクトでは、`hdiutil` を含む、従来のパッケージ化後に署名するフローが維持されるため、それらのフックによる署名前イメージへの変更も保持されます。このフック互換のイメージ編集処理には、macOS と `hdiutil` が必要です。

ポータブルな Apple 署名を行うには、アドホック署名の場合は `darwin.signing.identity: "-"` を設定し、エクスポート済み RSA PKCS#12 ID を使用する場合は `darwin.signing.certificate` を設定します。`darwin.signing.credential` には、そのパスワードを格納する環境変数の名前を指定します。証明書署名では共通の Go PFX/CMS 実装を使用し、キーチェーンから秘密鍵が自動的にエクスポートされることはありません。生成したテスト用 ID を使用した検証では、thin 実行ファイルと universal 実行ファイル、ネストされたアプリバンドル、フレームワーク、DMG が、改ざん拒否を含む Apple の厳格なネイティブ署名チェックに合格しています。

公証には、エクスポート済み Developer ID ID に加えて、互換性のある RFC 3161 `darwin.signing.timestamp_server` も必要です。共通の検証機能は、サポート対象の SHA-2 タイムスタンププロファイルを受け入れ、サポート対象外の従来アルゴリズムを拒否します。アドホック署名は公証できません。ネイティブの証明書形式の検証は、実際の Developer ID に対する信頼や、Apple による申請の受理とは別のものです。

Go のステープラーは Apple のチケット署名を認証し、そのチケットがすべてのアーキテクチャを対象としていることを確認してから、非公開コピーに添付します。この処理では、Apple がサポート対象 API として公開していない、観測に基づくチケット配信プロトコルを使用します。Gatekeeper による現在の失効チェックを置き換えるものではありません。アプリと署名済み DMG の出力は、ネイティブの `stapler validate` と厳格な署名チェックに合格しています。

Windows MSIX アーカイブは Go で直接作成されます。署名なしパッケージの作成には、Windows SDK も Docker も必要ありません。エクスポートした RSA PFX ID を使用すると、どのデスクトップホスト上でも、Windows 実行可能ファイルと個々の MSIX パッケージに Go で直接署名できます。`windows.signing.certificate` にはプロジェクトからの相対パスで PFX を指定し、`windows.signing.credential` にはそのパスワードを格納する環境変数を指定します。AES/PBKDF2 および従来の 3DES エクスポートがサポートされています。RC2 エクスポートは、再エクスポートするか、ネイティブの Windows SignTool にフォールバックする必要があります。証明書ストア／プロバイダーの ID と MSI の署名には、Windows SignTool が必要です。デフォルトのタイムスタンプサービスは DigiCert です。署名済みの出力を公開する前に、タイムスタンプが検証されます。形式上の制限については、[署名ガイド](/guides/build/signing/)を参照してください。
NSIS は外部の NSIS コンパイラを使用し、`packages.nsis.defines` を通じてリテラル定義を受け取ります。macOS では、継承されたロケールが UTF-8 でない場合に限り、Wails は NSIS コンパイラの子プロセスだけに UTF-8 ロケールを設定します。ターミナルとアプリケーションの環境は変更されません。標準の NSIS インストーラーは、パッケージ化されるすべての実行可能ファイルが CLI またはサーバーである場合、WebView2 ブートストラッパーを含めません。プライマリまたはコンパニオンが GUI の場合は、WebView2 が維持されます。カスタム NSIS スクリプトでは、独自のランタイムインストールロジックとの互換性を保つため、ブートストラッパーが維持されます。

Linux では、AppImage、DEB、RPM、Arch Linux、Void の各パッケージがサポートされます。Wake は、GTK/WebKit のデプロイ、SquashFS 圧縮、zsync の生成を含め、AppImage を Go で直接組み立てます。linuxdeploy、GTK プラグイン、appimagetool、mksquashfs は呼び出しません。ターゲットアーキテクチャ用としてバージョンを固定し、チェックサムを検証した AppImage type-2 ランタイムがダウンロードされます。代わりに、`packages.appimage.runtime_path` でプロジェクト相対のランタイム ELF を選択できます。

Linux では、デフォルトのライブラリソースはホストのルートファイルシステムです。別のアーキテクチャ向けに、または Windows/macOS から GUI をパッケージ化するには、`packages.appimage.runtime_root` に、ライブラリ、コンパイル済みスキーマ、モジュールキャッシュを含む、プロジェクト相対のターゲット Linux ルートツリーを設定します。指定するのは `usr/` ディレクトリ自体ではなく、`usr/` を含むルート全体のレイアウトです。サポート対象とする最も古い Linux ディストリビューションのランタイムを指定してください。Wake はターゲットプログラムを実行せずに、ELF のインポートと明示的な GTK/WebKit プロファイルを読み取ります。選択された依存関係の閉包は同梱されますが、glibc と動的ローダーは同梱されません。ランタイムの更新後は、ライブラリとリソースのコンテンツハッシュによってパッケージキャッシュが無効になります。

Pure Go の GUI ビルドがサポートされています。Wake は、`gtk3` でビルドされた GUI 実行ファイルには GTK 3 を選択し、それ以外には GTK 4 を選択します。CLI/サーバーに付随する GUI コンパニオンも同様です。1 つのイメージに含まれるすべての GUI 実行ファイルは、同じ GTK バージョンを使用する必要があります。ビルド済みバイナリでは、`packages.appimage.gtk_version` によって `3` または `4` を明示的に選択できます。インポートされた CGo ライブラリも、その選択と一致している必要があります。

zsync ファイルを生成するには `packages.appimage.update_info` を、SHA256 ファイルを生成するには `packages.appimage.checksum: true` を設定します。サイドカーはパッケージキャッシュから復元され、AppImage とともに成果物レシートに記載されます。指定された GTK/WebKit ツリーを使用した Windows と macOS でのパッケージ化が検証されており、生成された Pure Go GUI イメージは Linux 上のランタイムチェックに合格しています。条件を揃えたある CGo GUI ビルドでは、ネイティブ処理に 8.77 秒、linuxdeploy に 26.19 秒かかりました。ネイティブの gzip イメージは 133 MiB、従来の zstd 処理では 114 MiB だったため、この測定値は圧縮だけでなく、パッケージ化処理全体を対象としています。

Arch パッケージでは、プレリリースがネイティブのパッケージメタデータに保持されます。たとえば `1.2.3-beta.2` は `1.2.3beta.2-1` となり、最終リリースより前に並びます。Void の依存関係には、`gtk4` のようなパッケージ名のみを指定でき、`gtk4>=0` に正規化されます。明示的なバージョン比較と完全な `name-version_revision` 値は維持されます。ライブラリのファイル名ではなく、Void のパッケージ名を使用してください。

DEB、RPM、Arch の署名では、GnuPG、`dpkg-sig`、`rpmsign` を使用せず、サポート対象のどのデスクトップホストでも、エクスポート済みの ASCII Armor 形式またはバイナリ形式の OpenPGP 秘密鍵を使用できます。`linux.signing.certificate` にキーファイルを、`credential` にそのパスフレーズを格納する環境変数の名前を設定します。DEB と RPM には署名が埋め込まれ、Arch では分離署名が公開されます。秘密鍵は署名手順でのみ読み取られ、パッケージアセットや再利用可能なキャッシュ成果物にはコピーされません。

GnuPG のキー ID だけを含む既存の DEB/RPM 設定では、引き続き外部署名機能が使用されます。この依存関係をなくすには、秘密鍵をエクスポートし、代わりにそのパスを設定します。`./signing/release.asc` のようなパスを使用してください。キーファイルが存在しない、読み取れない、または無効な場合は、キーリングへフォールバックせずに失敗します。DEB の署名では `dpkg-sig` パッケージ署名形式が維持されます。これは、apt リポジトリのリリースメタデータへの署名を置き換えるものではありません。

Arch パッケージと分離署名を作成するには、エクスポートした ASCII Armor 形式の PGP 秘密鍵ファイルを使用し、必要に応じてパスフレーズを格納する環境変数名を指定します。

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

ローカル環境または CI のシークレットストアで指定した環境変数を設定してから、`wails3 build release` を実行します。暗号化されていない鍵では `credential` を省略します。Go 署名ツールは、サポートされるすべてのデスクトップホストで動作し、GPG を必要としません。元の `.pkg.tar.zst` と同じディレクトリに `.pkg.tar.zst.sig` を公開し、両方を成果物レシートに個別に記録します。鍵の内容がステージングされたり、ビルドキャッシュに保存されたりすることはありません。署名とその公開結果がキャッシュから再利用されることもありません。署名に失敗した場合、以前に公開されたファイルのペアは維持されます。いずれかのファイルのインストールに失敗した場合、公開処理によって両方のファイルがロールバックされます。ただし、これによって同時アクセスする読み取り側に対して、2 回のファイルシステム上の名前変更がアトミックになるわけではありません。

macOS アプリバンドルにはコンパニオンプログラムが含まれます。`darwin.microphone_usage_description` を設定すると、用途を説明するテキストが追加され、カスタムのエンタイトルメントファイルが指定されていない場合は、オーディオ入力のエンタイトルメントも追加されます。カスタムの plist とエンタイトルメントのソースは変更されません。エクスポート済み RSA PKCS#12 ID では、サポート対象のどのデスクトップホストでも Go 署名が使用されます。`identity: "-"` を指定すると、ポータブルなアドホック署名が選択されます。キーチェーン ID では、引き続き macOS 上で `codesign` を使用します。チケットのステープル処理は Go で実行されます。公証の申請には、既存の macOS 公証キーチェーンプロファイル、または App Store Connect のチーム API 認証情報を使用できます。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

API モードでは、`notarytool` を使用せず、Go から申請、アップロード、ポーリングを実行します。3 つのフィールドをすべて設定し、`credential` は省略してください。秘密鍵はバージョン管理に含めないでください。鍵は署名時にのみ読み取られ、生成されたアセットやキャッシュにはコピーされません。このモードではチーム API キーと最大 5 GiB のアーカイブがサポートされます。マルチパートアップロードと個人 API キーは実装されていません。申請が拒否された場合は申請 ID が保持され、Apple の診断ログが非公開の一時ファイルに保存されます。そのパスはエラー内で報告されます。キャンセルまたは拒否された申請によって、以前に公開された出力が置き換えられることはありません。受理後、Wake は Apple のチケットを取得して検証し、Go でステープル処理を行います。エクスポート済み ID とチーム API 認証情報を使用すれば、Linux、Windows、macOS のいずれでも署名、申請、ステープル処理を実行できます。キーチェーンプロファイルには、引き続き macOS と `notarytool` が必要です。

ポータブルな証明書署名では、AES/PBKDF2 または従来の 3DES を使用する PKCS#12 ファイル内の RSA ID がサポートされます。RC2 で暗号化されたエクスポートは、サポート対象の暗号方式を使用して再エクスポートする必要があります。公証済みのポータブル署名には、明示的な SHA-2 RFC 3161 タイムスタンプサービスが必要です。Go の検証機能は、従来の SHA-1 タイムスタンプ署名を拒否します。チケットサービスは Apple の観測済みプロトコルに従い、Apple Root G3 を固定して使用します。配布における信頼性と失効については、引き続き Gatekeeper が最終的な判断を行います。

署名済みの出力には、たとえば
`notes.signed.app` や `notes.signed.dmg` のように、プラットフォーム固有の拡張子が維持されます。

## 開発時の制御

`wails3 dev` の実行中に、コマンドを入力して Enter キーを押します。

| コマンド | 操作 |
| --- | --- |
| `r` | 変更されていないバックエンドも含め、強制的に再ビルドして再起動します。 |
| `b` | バインディングを再生成して再ビルドします。 |
| `f` | フロントエンドを再起動します。 |
| `c` | ターミナルをクリアします。 |
| `q` | セッションとその子プロセスを停止します。 |
| `h` | 制御コマンドを表示します。 |

ビルドに失敗しても、現在のアプリケーションは実行を続けます。制御コマンドには対話型ターミナルが必要です。パイプから入力しても、セッションの動作は変わりません。

## 実行前に処理内容を確認する

プランの検査では、実行時と同じ不変グラフが解決されます。テキストは人が読むためのものであり、JSON はバージョン管理された機械可読の契約です。

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## ビルドせずに構成を検証する

CI やエディターのタスクで、または構成をコミットする前に、専用の事前検証コマンドを使用してください。プロファイルを指定しない場合は、ベースマニフェストと名前付きの全プロファイルを検証します。プロファイルを指定した場合は、その完全なリクエストだけを検証します。

```bash
wails3 config check
wails3 config check release
```

構文、スキーマ、セマンティクス、パス、ターゲット、形式、署名、プランナー、および構成によって指定されたホスト機能に関するエラーには、Build および Dev と同じ診断レンダラーが使用されます。診断ではフィールドとソース範囲を特定し、該当するソース行とキャレットを表示し、有用な場合は修正のヒントも示します。診断のソース抜粋では環境変数の値が伏せられます。コマンドラインでの選択だけが原因のホストエラーは、無効な YAML として誤って報告されず、通常のコマンドエラーとして扱われます。

## プロジェクト所有のアセットを使用する

マニフェストのパスは、アイコン、プラットフォームマニフェスト、パッケージテンプレート、メンテナースクリプトなど、ユーザーが所有するファイルを参照します。Wails がそれらのソースを編集することはありません。ツールが使用する前に、使い捨ての `.wails/` ワークスペースへコピーします。最終成果物は `build.output` の下に配置され、検証済みのレシートが `.wails/artifacts/receipt.json` に書き込まれます。

## Android にデプロイする

本番用の Android プロファイルでは AAB が生成されます。ローカルデプロイでは、Wails が ABI 固有のデバッグ APK を内部でビルドし、ADB でインストールしてアプリを起動します。

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

セレクターを明示的に指定しない場合、対話型ターミナルには、接続済みデバイスと構成済みの Android Virtual Device が提示されます。複数のターゲットを利用できる状態で自動化する場合は、`--device` または `--emulator` を指定する必要があります。

アプリケーション単位のログストリームを接続し、アプリが終了するか Ctrl+C を押すまでターゲットを監視します。

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

接続実行ではアプリケーションの再起動にも追従し、ターゲットの切断やオフラインへの移行を報告します。素早く反復できるよう、デフォルトではエミュレーターを実行したままにします。`--stop-emulator` は明示的に指定するオプションであり、`--logs` が必要です。接続セッションの終了時に、ターゲットがエミュレーターの場合に限って停止します。Wails が新しいエミュレーターを起動したものの、起動準備の確認に失敗した場合やキャンセルされた場合は、そのプロセスを自動的にクリーンアップします。

## 範囲を限定したライフサイクルフックを実行する

試験的なフック機能では、宣言された各フェーズでプロジェクト所有のスクリプトを 1 つ呼び出します。フックはシェルコマンドではありません。`directory` が設定されていない限り、Wails はプロジェクトルートからファイルを直接実行します。

Unix ホストでは、スクリプトに実行権限があり、先頭が shebang でなければなりません。Windows ホストでは、固定のインタープリターアダプターを介して `.cmd`、`.bat`、`.ps1` の各スクリプトを使用できます。どちらのプラットフォームでも、Wails はインラインのシェルコマンドを解析しません。

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

サポートされるフェーズは、`before_build`、`after_build`、`before_package`、`after_package`、`before_sign`、`after_sign` です。ビルドフックは、必要に応じて共有プロジェクト全体で 1 回、またはターゲットごとに 1 回実行されます。パッケージ化フックと署名フックは、ターゲットに対して要求された全形式の前後にバリアを形成しますが、バリア内の各形式の処理は並行して実行されます。

フックには、Wails が所有する環境変数 `WAILS_HOOK_CONTEXT_FILE` が 1 つ渡されます。この変数は、バージョン管理された呼び出しコンテキストを含む、保護された JSON ファイルを指します。

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

コンテキストには解決済みのパスが含まれますが、継承したプロセス環境やシークレットがコピーされることはありません。Wails は、所有者だけがアクセスできるファイル権限を設定し、`.wails/hooks/context/` の下にこのファイルをアトミックに作成します。成功またはキャンセルの後に削除されます。失敗時には診断用に保持され、そのパスがエラーに含まれます。0 以外の終了コードでは通常の Plan ステップが失敗し、キャンセル時にはフックのプロセスグループが終了されます。

デフォルトでは、フックは呼び出しのたびに実行されます。スクリプトのすべての入力と出力を宣言できる場合に限り、キャッシュを明示的に有効にしてください。

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

スクリプトとその実行権限モードは暗黙の入力です。複数の出力は、ルート以外の 1 つのディレクトリを共有する必要があります。このディレクトリが、Wails によって記録および復元される範囲限定の成果物となります。出力ルートには、スクリプトや宣言済みの入力を含めることはできません。また、キャッシュされる出力ツリーに含められるのは通常のファイルとディレクトリだけであり、シンボリックリンクや特殊ファイルは一切含められません。常時実行フックでは、宣言されていない副作用のフィンガープリントを安全に生成できないため、後続の処理をキャッシュできなくなります。キャッシュ可能なフックでは、宣言した成果物の識別情報を通じて、後続処理のキャッシュ可能性が維持されます。

## 既存のプロジェクトを移行する

無効な状態のドラフトを作成し、確認します。

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

引数を取らず、プロジェクト所有のスクリプトを呼び出す保守的なライフサイクルタスクは、フックに変換されます。到達可能なその他のカスタム Taskfile 動作は、ブロッカーとして報告されます。
`wails.migrated.yaml`を確認した後、明示的に有効化します。

```bash
wails3 migrate --activate
```

有効化すると、`wails.yaml`がアトミックに作成されます。Taskfile、`build/config.yml`、スクリプト、アセットが削除、名前変更、編集されることはありません。無効な YAML は YAML エラーとして失敗し、従来の経路にフォールバックすることはありません。

## すべてのデフォルト値を具体化する

解決済みの完全な設定を確認またはコピーするには、非アクティブなスナップショットを書き出します。

```bash
wails3 eject
```

これにより、`wails.ejected.yaml`が作成されます。アクティブな設定が書き換えられることはありません。

## 互換性期間

`wails3 package`と`wails3 sign`は、同じプランナーを指す非推奨のエイリアスです。新しい自動化では、プロファイルまたはフォーマットを指定して`wails3 build`を使用してください。ロールアウトはリリース単位で制御されます。最初のリリースでは、新規プロジェクトで YAML を使用するようにし、次のリリースでは従来の Taskfile ルーティングについて警告し、さらに次のリリースではエイリアスとともに Taskfile の実行を削除します。各移行には、そのリリースのネイティブホスト受け入れマトリクスがすべて成功している必要があります。

すべてのフィールド、型、デフォルト値、適用条件については、生成された[`wails.yaml`フィールドリファレンス](/reference/wails-yaml-fields/)を参照してください。

[`wails3 run`設定リファレンス](/reference/wails-yaml-run/)では、ローカル起動のデフォルト、プラットフォームごとのオーバーライド、およびマニフェストがない場合の`go run .`フォールバックについて説明しています。
