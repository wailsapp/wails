---
title: "macOS パッケージング"
description: "Wails アプリケーションを macOS 向けに配布できるようパッケージ化します"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## macOS の非公開 API

Wails v3 は、デフォルトで macOS の公開 API を使用します。文書化されていない Apple API を必要とする機能を有効にするには、単一の Go ビルドタグ `private_mac_apis` を指定してアプリをビルドします。

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Go で直接ビルドする場合は、`go build -tags private_mac_apis .` を使用します（本番環境では `-tags production,private_mac_apis`）。非公開の動作に依存する既存のアプリケーションでは、その動作を維持するためにこのタグを追加する必要があります。このタグは macOS デスクトップ向けビルドにのみ適用されます。

影響を受ける機能とオプション値の完全な一覧、公開 API のみを使用するビルドでの正確なフォールバック動作、Liquid Glass のスタイルマッピング、およびインスペクターのビルド構成については、[macOS の非公開 API](/guides/build/private-macos-apis/)を参照してください。このタグがない場合、非公開 API 専用の操作は何も行いません。公開 Go API に変更はありません。

## アプリケーションバンドル

アプリを標準の macOS `.app` バンドルとしてパッケージ化します。

```bash
wails3 package GOOS=darwin
```

これにより、以下を含む `bin/<AppName>.app` が作成されます。

- `Contents/MacOS/` 内のコンパイル済みバイナリ
- `Contents/Resources/` 内のアプリアイコン（`icons.icns`、またはアセットカタログ `Assets.car` が存在する場合はそこから取得）
- アプリのメタデータを含む `Info.plist`

### 同期フォルダー内のプロジェクト

一部の macOS ファイルプロバイダーは、`.app` ディレクトリに Finder メタデータを追加します。そのため、`codesign` がバンドルを拒否したり、署名後に署名が無効になったりします。プロバイダーがメタデータを再度追加する場合、メタデータを削除するだけでは不十分です。

生成された macOS ビルドタスクは、この動作を確認します。検出すると、Wails はプロジェクトの生成出力ディレクトリを `~/Library/Application Support/Wails/build-output/<project>-<hash>` に移動し、元の出力ディレクトリをシンボリックリンクに置き換えます。既存の出力ファイルは保持され、パッケージ化、署名、開発モードでは引き続き `bin/<AppName>.app` がパスとして使用されます。移動を実行すると、CLI にローカルの保存先が表示されます。

同期されない別のローカル保存先を選択するには、`WAILS_MACOS_OUTPUT_ROOT` を設定します。移動先のルートは、プロジェクトと同じファイルシステム上に存在する必要があります。プロジェクト外を明示的に指定した `BIN_DIR` は、設定どおりに使用されます。既存のプロジェクトでこの処理を使用するには、更新された `build/darwin/Taskfile.yml` が必要です。

配布するには、`bin/<AppName>.app` をコピーするか、DMG を作成します。`bin` シンボリックリンクだけをコピーしても、ビルド成果物は含まれません。直接実行した `wails3 tool sign` コマンドは、入力フォルダーにメタデータが再付加される場合にエラーを報告します。そのバンドルを同期されないフォルダーへ移動してから署名してください。

## バンドルリソース

`Contents/Resources/` は、macOS アプリに同梱する読み取り専用ファイルの標準的な保存場所です。大きなテンプレート、初期データ、メディア、言語パックなど、`embed` で Go 実行ファイルにコンパイルするのではなく、必要に応じて開くペイロードに使用します。

Wails は、すでにこのディレクトリにアプリケーションアイコンを配置しています。独自のファイルを追加するには、それらを `build/resources/` などのソースディレクトリに配置し、`build/darwin/Taskfile.yml` の `create:app:bundle` タスクにコピーステップを追加します。

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Taskfile の `darwin:run` タスクを使用する場合は、その `run` タスクに、`{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/` をコピー先とする同等のコマンドを追加します。

### Go からのリソースの読み込み

macOS プラットフォームパッケージをインポートします。

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

小さなファイルには、`LoadResource` を使用します。

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

大きなファイルには、`ResourceFS` を使用します。これは `Contents/Resources` をルートとする `io/fs.FS` を返すため、呼び出し元はリソース全体を最初に Go のバイトスライスへ読み込まずに、リソースを開いてストリーミングできます。

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

リソース名は、`Contents/Resources` を基準とするスラッシュ区切りのパスです。実行ファイルが `.app/Contents/MacOS` から実行されている場合を除き、`ResourceFS` と `LoadResource` は `mac.ErrNotInAppBundle` を返します。

バンドルリソースは変更不可として扱ってください。署名済みアプリケーション内のファイルを変更すると、コード署名が無効になります。ダウンロードしたデータ、生成したデータ、またはユーザーが編集できるデータは、代わりにユーザーの Application Support ディレクトリへ保存してください。

### ユニバーサルバイナリ

Apple Silicon と Intel Mac の両方に対応するようビルドします。

```bash
wails3 task darwin:package:universal
```

これにより、両方のアーキテクチャでネイティブに動作する単一の `.app` が作成されます。ユニバーサルバイナリはどのプラットフォームでもビルドできます。Linux と Windows では、`wails3 tool lipo` が自動的に使用されます。

## バンドルのカスタマイズ

以下をカスタマイズするには、`build/darwin/Info.plist` を編集します。

- バンドル識別子（`CFBundleIdentifier`）
- アプリ名とバージョン
- macOS の最小バージョン
- ファイルの関連付け
- URL スキーム

アプリアイコンは、`build/` ディレクトリ内のアセットから生成されます。`generate:icons` タスクを使用します。

```bash
wails3 task common:generate:icons
```

これは `build/appicon.png` を使用して、`darwin/icons.icns` と `windows/icon.ico` を生成します。macOS では、`build/appicon.icon`（Icon Composer 形式）も指定できます。このタスクは `-iconcomposerinput appicon.icon -macassetdir darwin` を渡し、`.icon` ファイルから `Assets.car` と `darwin/icons.icns` を生成します（macOS 以外のプラットフォームではスキップされます）。`Assets.car` が存在する場合は、`Info.plist` と `CFBundleIconName` も適切に更新されるよう、`update:build-assets` タスクを実行します。

```bash
wails3 task common:update:build-assets
```

`build/` ディレクトリからアイコンコマンドを手動で実行するには、次のようにします。

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## コード署名

配布用にアプリへ署名します。

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

`build/darwin/Taskfile.yml` で署名を設定します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### 公証

Mac App Store 以外で配布するアプリには、Apple による公証が必要です。

```bash
wails3 task darwin:sign:notarize
```

まず、認証情報を保存します。対話形式のウィザード（`wails3 signing`）を実行するか、`notarytool` を直接呼び出します。

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

`build/darwin/Taskfile.yml` で設定します。

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

詳細については、[アプリケーションへの署名](/guides/build/signing/)を参照してください。

## DMG インストーラー

Wails に同梱されている 3 テンプレートには、`wails3 task darwin:package:dmg` が用意されています。最初に `.app` を作成し、次に DMG ライブラリを使用してデザインされた DMG をビルドします。デフォルトでは、DMG は赤いドラゴンのマークと WAILS ワードマークが入った Wails ブランドのグラデーション背景を使用します。

```bash
wails3 task darwin:package:dmg
```

下位レベルの `darwin:create:dmg` タスクは、既存の `.app` バンドルから DMG を作成します。このタスクは Taskfile から直接設定できます。

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### デフォルトレイアウト

生成される DMG には、以下が含まれます。

- 左側のアプリケーションバンドル
- 右側の `Applications` リンク
- 540×380 ピクセルに設定された Finder ウインドウ
- 各アイコンの下にラベルを表示する、96 ポイントのアイコンサイズ
- `build/darwin/dmg-background.png` の Wails ブランド背景

アプリケーションと `Applications` のアイコンは、設定されたウインドウ寸法を基準に配置されます。そのため、`DMG_WINDOW_WIDTH` または `DMG_WINDOW_HEIGHT` を変更しても、デフォルトの 2 アイコンレイアウトでは比例した間隔が維持されます。最適な結果を得るには、Finder ウインドウと同じピクセル寸法の背景画像を使用してください。

### DMG アセットの置き換え

`build/darwin/` 以下に生成されるファイルは通常のプロジェクトアセットであり、置き換えることができます。

- `DMG_BACKGROUND` は、Finder ウインドウの内容の背後に表示される画像を制御します。
- `DMG_VOLUME_ICON` は、マウントされたボリュームに表示されるアイコンを制御します。
- `DMG_FILE_ICON` は、生成された `.dmg` ファイルに Finder で表示されるアイコンを制御します。

ボリュームアイコンとDMGファイルアイコンは別々のリソースです。アプリケーションアイコンを置き換えても、どちらも自動的には置き換わりません。

### 追加ファイルの追加

`DMG_FILES`を使用すると、インストーラースクリプト、リリースノート、ライセンス、その他のリソースをアプリケーションとともに収録できます。値には、`name=path`ペアをカンマで区切ったリストを指定します。

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=`の前の名前は、DMG内に表示されるファイル名です。`=`の後のパスは、プロジェクト内のソースファイルを示します。先頭と末尾の空白は無視されます。

表示名はそれぞれ一意でなければなりません。追加ファイルで、アプリケーションバンドルや`Applications`エントリなど、パッケージャーがすでに作成したエントリを置き換えることはできません。名前が競合すると、破損したDMGが生成されるのではなく、エラーが発生してパッケージ作成に失敗します。

@note{type="note"}
DMGの作成ではmacOSのディスクイメージツールとFinderツールを使用するため、macOSでのみサポートされます。クロスコンパイルした`.app`バンドルはほかの環境でも作成できますが、最終的なDMGはMac上で生成する必要があります。

@end

## トラブルシューティング

### 「App is damaged and can't be opened」と表示される

アプリが署名されていません。Developer ID証明書で署名するか、ユーザーがGatekeeperを回避できます。

```bash
xattr -cr /path/to/YourApp.app
```

### 公証に失敗する

よくある問題：

- **認証情報が無効**：`xcrun notarytool store-credentials`（または`wails3 signing`）を再実行します
- **Hardened Runtimeが必要**：必要に応じて、entitlementsに`com.apple.security.cs.allow-unsigned-executable-memory`が含まれていることを確認します
- **タイムスタンプがない**：署名プロセスではタイムスタンプが自動的に付与されるはずです

### クロスコンパイルしたアプリが実行できない

クロスコンパイルしたmacOSバイナリは署名されていません。Macに転送し、テスト前に署名してください。

```bash
codesign --force --deep --sign - YourApp.app
```
