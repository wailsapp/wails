---
title: "macOS パッケージング"
description: "Wails アプリケーションを macOS 向けに配布できるようパッケージ化します"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## macOS の非公開 API

Wails v3 はデフォルトで macOS の公開 API を使用します。文書化されていない Apple API を必要とする機能を有効にするには、単一の Go ビルドタグ `private_mac_apis` を指定してアプリをビルドします。

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Go で直接ビルドする場合は、`go build -tags private_mac_apis .`（本番用には `-tags production,private_mac_apis`）を使用します。非公開の動作に依存する既存のアプリケーションでその動作を維持するには、このタグを追加する必要があります。このタグは macOS デスクトップビルドにのみ適用されます。

影響を受ける機能とオプション値の完全な一覧、公開 API を使用したビルドでの正確なフォールバック動作、Liquid Glass スタイルのマッピング、およびインスペクターのビルド構成については、[macOS の非公開 API](/guides/build/private-macos-apis/)を参照してください。非公開 API 専用の操作は、タグがない場合は何も行いません。公開 Go API に変更はありません。

## アプリケーションバンドル

アプリを標準の macOS `.app` バンドルとしてパッケージ化します。

```bash
wails3 package GOOS=darwin
```

これにより、以下を含む `bin/<AppName>.app` が作成されます。

- `Contents/MacOS/` 内のコンパイル済みバイナリ
- `Contents/Resources/` 内のアプリアイコン（`icons.icns`、またはアセットカタログ `Assets.car` が存在する場合はそこから生成）
- アプリのメタデータを含む `Info.plist`

## バンドルリソース

`Contents/Resources/` は、macOS アプリに同梱する読み取り専用ファイルの標準的な格納場所です。大きなテンプレート、初期データ、メディア、言語パックなど、`embed` を使用して Go 実行ファイルにコンパイルするのではなく、必要に応じて開くペイロードに使用します。

Wails はすでにアプリケーションアイコンをこのディレクトリに配置しています。独自のファイルを追加するには、`build/resources/` などのソースディレクトリに配置し、`build/darwin/Taskfile.yml` 内の `create:app:bundle` タスクにコピーステップを追加します。

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

Taskfile の `darwin:run` タスクを使用する場合は、それに対応するコマンドを `run` タスクに追加し、コピー先を `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/` に設定します。

### Go からのリソースの読み取り

macOS プラットフォームパッケージをインポートします。

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

小さいファイルには `LoadResource` を使用します。

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

大きいファイルには `ResourceFS` を使用します。これは `Contents/Resources` をルートとする `io/fs.FS` を返すため、呼び出し元はリソース全体を最初に Go のバイトスライスへ読み込むことなく、開いてストリーミングできます。

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

リソース名には、`Contents/Resources` を基準とするスラッシュ区切りのパスを指定します。実行ファイルが `.app/Contents/MacOS` から実行されていない限り、`ResourceFS` と `LoadResource` は `mac.ErrNotInAppBundle` を返します。

バンドルリソースは不変として扱ってください。署名済みアプリケーション内のファイルを変更するとコード署名が無効になります。ダウンロードしたデータ、生成したデータ、またはユーザーが編集できるデータは、代わりにユーザーの Application Support ディレクトリへ保存してください。

### ユニバーサルバイナリ

Apple Silicon と Intel Mac の両方に対応するようビルドします。

```bash
wails3 task darwin:package:universal
```

これにより、両方のアーキテクチャでネイティブに動作する単一の `.app` が作成されます。ユニバーサルバイナリはどのプラットフォームでもビルドできます。Linux と Windows では `wails3 tool lipo` が自動的に使用されます。

## バンドルのカスタマイズ

`build/darwin/Info.plist` を編集して、以下をカスタマイズします。

- バンドル識別子（`CFBundleIdentifier`）
- アプリ名とバージョン
- macOS の最小バージョン
- ファイルの関連付け
- URL スキーム

アプリアイコンは `build/` ディレクトリ内のアセットから生成されます。`generate:icons` タスクを使用します。

```bash
wails3 task common:generate:icons
```

これは `build/appicon.png` を使用して、`darwin/icons.icns` と `windows/icon.ico` を生成します。macOS では `build/appicon.icon`（Icon Composer 形式）も指定できます。このタスクは `-iconcomposerinput appicon.icon -macassetdir darwin` を渡し、`.icon` ファイルから `Assets.car` と `darwin/icons.icns` を生成します（macOS 以外のプラットフォームではスキップされます）。`Assets.car` が存在する場合は、`Info.plist` と `CFBundleIconName` がそれに応じて更新されるよう、`update:build-assets` タスクを実行します。

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

まず、認証情報を保存します。対話型ウィザード（`wails3 setup signing`）を実行するか、`notarytool` を直接呼び出します。

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

Wails 3 に同梱されるテンプレートには、`wails3 task darwin:package:dmg` が用意されています。最初に `.app` を作成し、次に DMG ライブラリを使用してスタイル設定済みの DMG をビルドします。デフォルトでは、DMG の背景に Wails ブランドのグラデーション画像が使用され、赤いドラゴンのマークと WAILS のワードマークが表示されます。

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
- 540×380 ピクセルの Finder ウインドウ
- 96 ポイントのアイコンサイズと、各アイコンの下に表示されるラベル
- `build/darwin/dmg-background.png` の Wails ブランド背景

アプリケーションアイコンと `Applications` アイコンは、設定されたウインドウサイズを基準に配置されます。そのため、`DMG_WINDOW_WIDTH` または `DMG_WINDOW_HEIGHT` を変更しても、デフォルトの 2 アイコンレイアウトでは相対的な間隔が維持されます。最適な結果を得るには、Finder ウインドウと同じピクセルサイズの背景画像を使用してください。

### DMG アセットの置き換え

`build/darwin/` 以下に生成されたファイルは通常のプロジェクトアセットであり、置き換えることができます。

- `DMG_BACKGROUND` は、Finderウインドウの内容の背後に表示される画像を制御します。
- `DMG_VOLUME_ICON` は、マウントされたボリュームに表示されるアイコンを制御します。
- `DMG_FILE_ICON` は、生成された `.dmg` ファイルにFinderで表示されるアイコンを制御します。

ボリュームアイコンとDMGファイルアイコンは別々のリソースです。アプリケーションアイコンを置き換えても、どちらも自動的には置き換わりません。

### 追加ファイルの追加

インストーラースクリプト、リリースノート、ライセンス、その他のリソースをアプリケーションと一緒に含めるには、`DMG_FILES` を使用します。値には、`name=path` ペアをカンマで区切ったリストを指定します。

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=` の前の名前は、DMG内に表示されるファイル名です。`=` の後のパスは、プロジェクト内のソースファイルです。先頭と末尾の空白は無視されます。

表示名はそれぞれ一意である必要があります。追加ファイルで、アプリケーションバンドルや `Applications` エントリなど、パッケージャーがすでに作成したエントリを置き換えることはできません。名前が競合すると、壊れたDMGが生成されるのではなく、エラーによりパッケージ化が失敗します。

@note{type="note"}
DMGの作成ではmacOSのディスクイメージツールとFinderツールを使用するため、macOSでのみサポートされます。クロスコンパイルした `.app` バンドルは他の環境でも作成できますが、最終的なDMGはMac上で生成する必要があります。

@end

## トラブルシューティング

### 「アプリが破損しているため開けません」

アプリが署名されていません。Developer ID証明書で署名するか、ユーザーがGatekeeperを回避できます。

```bash
xattr -cr /path/to/YourApp.app
```

### 公証に失敗する

よくある問題：

- **認証情報が無効**：`xcrun notarytool store-credentials`（または `wails3 setup signing`）を再実行します
- **Hardened Runtimeが必要**：必要に応じて、entitlementsに `com.apple.security.cs.allow-unsigned-executable-memory` が含まれていることを確認します
- **タイムスタンプがない**：署名プロセスではタイムスタンプが自動的に付与されるはずです

### クロスコンパイルしたアプリが実行できない

クロスコンパイルしたmacOSバイナリは署名されていません。テストする前にMacへ転送して署名します。

```bash
codesign --force --deep --sign - YourApp.app
```
