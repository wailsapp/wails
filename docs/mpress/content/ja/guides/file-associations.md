---
title: "ファイルの関連付け"
description: "Wails アプリケーションのファイルの関連付けを設定する"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

対応プラットフォーム：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

ファイルの関連付けを使用すると、ユーザーが特定の種類のファイルを開いたときに、アプリケーションでそのファイルを処理できます。これは、テキストエディターや画像ビューアーなど、特定のファイル形式を扱うアプリケーションで特に便利です。このガイドでは、Wails v3 アプリケーションにファイルの関連付けを実装する方法について説明します。

## 概要

Wails v3 では現在、次のプラットフォームでファイルの関連付けがサポートされています：

- Windows（NSIS インストーラーパッケージ）
- macOS（アプリケーションバンドル）

## 設定

ファイルの関連付けは、プロジェクトの `build` ディレクトリにある `config.yml` ファイルで設定します。

### 基本設定

ファイルの関連付けを設定するには：

1. `build/config.yml` を開きます
2. `fileAssociations` セクションにファイルの関連付けを追加します
3. `wails3 update build-assets` を実行してビルドアセットを更新します
4. アプリケーションオプションの `FileAssociations` フィールドを設定します
5. `wails3 package` を使用してアプリケーションをパッケージ化します

設定例を次に示します：

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### 設定プロパティ

| プロパティ | 説明 | プラットフォーム |
| --- | --- | --- |
| ext | 先頭のピリオドを除いたファイル拡張子（例：`txt`） | すべて |
| name | ファイル形式の表示名 | すべて |
| description | ファイルのプロパティに表示される説明 | Windows |
| iconName | build フォルダー内にあるアイコンファイルの名前（拡張子なし） | すべて |
| role | このファイル形式に対するアプリケーションの役割（例：`Editor`、`Viewer`） | macOS |
| mimeType | ファイルの MIME タイプ（例：`image/jpeg`） | macOS |

## ファイルオープンイベントのリッスン

アプリケーションでファイルオープンイベントを処理するには、`events.Common.ApplicationOpenedWithFile` イベントをリッスンします：

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## ステップ別チュートリアル

簡単なテキストエディターにファイルの関連付けを設定する手順を説明します：

@steps
### アイコンを作成する
- ファイル形式用のアイコンを作成します（推奨サイズ：16x16、32x32、48x48、256x256）
- プロジェクトの `build` フォルダーにアイコンを保存します
- `iconName` の設定に従ってアイコンに名前を付けます（例：`textFileIcon.png`）

@note{type="tip"}
`wails3 generate icons` を使用して、必要なアイコンを生成できます。詳細については、`wails3 generate icons --help` を実行してください。

@end

- macOS の場合は、`create:app:bundle:` タスクに `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` のようなコピー文を追加します。

### ファイルの関連付けを設定する
`build/config.yml` ファイルを編集し、ファイルの関連付けを追加します：

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### ビルドアセットを更新する
次のコマンドを実行してビルドアセットを更新します：

```bash
wails3 update build-assets
```

### アプリケーションオプションにファイルの関連付けを設定する
`main.go` ファイルで、アプリケーションオプションの `FileAssociations` フィールドを設定します：

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="アプリケーション設定と config.yml の両方でファイル拡張子が必要なのはなぜですか？"}
Windows では、ファイルの関連付けによってファイルを開くと、そのファイル名がアプリケーションの最初の引数として渡され、アプリケーションが起動します。アプリケーションには、最初の引数がファイルなのかコマンドライン引数なのかを判別する方法がありません。そのため、アプリケーションオプションの `FileAssociations` フィールドを使用して、最初の引数が関連付けられたファイルかどうかを判定します。

@end

### アプリケーションをパッケージ化する
次のコマンドを使用してアプリケーションをパッケージ化します：

```bash
wails3 package
```

パッケージ化されたアプリケーションは `bin` ディレクトリに作成されます。その後、アプリケーションをインストールしてテストできます。

## 補足事項

- アイコンは PNG 形式で build フォルダーに配置することを推奨します。
- ファイルの関連付けをテストするには、パッケージ化されたアプリケーションをインストールする必要があります

@end
