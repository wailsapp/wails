---
title: "ダイアログの概要"
description: "アプリケーションにネイティブシステムダイアログを表示する"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## ネイティブダイアログ

Wails は、すべてのプラットフォームで動作する<strong>ネイティブシステムダイアログ</strong>を提供します。メッセージダイアログ（情報、警告、エラー、質問）、ファイルダイアログ（開く、保存、フォルダー）、および各プラットフォームに固有の外観と動作を備えたカスタムダイアログウィンドウを利用できます。

![「キャンセル」ボタンと「破棄」ボタンがある macOS 上の Wails 質問ダイアログ](/assets/screenshots/dialog-question-macos.png)

同じ API でも、サポート対象の各プラットフォームの規則に従って表示されます。この macOS の例では、デフォルトボタンとキャンセルボタンを含む質問ダイアログが Wails ウィンドウにアタッチされています。

## クイックスタート

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

<strong>これだけです。</strong>最小限のコードでネイティブダイアログを使用できます。

## ダイアログへのアクセス

ダイアログには、`app.Dialog`マネージャーを介してアクセスします。

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## ダイアログの種類

### 情報ダイアログ

簡単なメッセージを表示します。

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**使用例：**

- 成功メッセージ
- 情報通知
- 完了確認

### 警告ダイアログ

警告を表示します。

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**使用例：**

- 重大ではない警告
- 非推奨化の通知
- 注意メッセージ

### エラーダイアログ

エラーを表示します。

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**使用例：**

- エラーメッセージ
- 失敗通知
- 例外処理

### 質問ダイアログ

ユーザーに質問し、ボタンのコールバックで応答を処理します。

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**使用例：**

- 操作の確認
- 「はい」または「いいえ」で答える質問
- 複数選択肢

## ファイルダイアログ

### ファイルを開くダイアログ

開くファイルを選択します。

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**複数選択：**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### ファイルを保存するダイアログ

保存先を選択します。

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### フォルダー選択ダイアログ

ディレクトリ選択を有効にしたファイルを開くダイアログを使用して、ディレクトリを選択します。

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## ダイアログのオプション

### タイトルとメッセージ

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### ボタン

**シンプルなダイアログのデフォルトボタン：**

情報、警告、エラーの各ダイアログには、デフォルトの「OK」ボタンが表示されます。

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**質問ダイアログのカスタムボタン：**

`AddButton()`を使用してボタンを追加します。このメソッドは、コールバックを設定できる`*Button`を返します。

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**デフォルトボタンとキャンセルボタン：**

`SetDefaultButton()`を使用して、強調表示され、Enter キーで実行されるボタンを指定します。 `SetCancelButton()`を使用して、Escape キーで実行されるボタンを指定します。

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### ウィンドウへのアタッチ

ダイアログを特定のウィンドウにアタッチします。

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**動作：**

- ダイアログは親ウィンドウの中央に表示される
- ダイアログの表示中は親ウィンドウが無効になる
- ダイアログは親ウィンドウとともに移動する（macOS）

## プラットフォームごとの動作

@tabs{sync-key="platform"}
[macOS]
**macOS のダイアログ：**

- ネイティブの NSAlert の外観
- システムテーマ（ライト／ダーク）に従う
- キーボード操作をサポート
- 標準ショートカット（キャンセルは ⌘.）
- アクセシビリティ機能を搭載
- ウィンドウへのアタッチ時はシート形式で表示

**例：**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Windows のダイアログ：**

- ネイティブの TaskDialog の外観
- システムテーマに従う
- キーボード操作をサポート
- 標準ショートカット（キャンセルは Esc）
- アクセシビリティ機能を搭載
- 親ウィンドウに対してモーダル表示

**例：**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Linux のダイアログ：**

- GTK ダイアログの外観
- デスクトップテーマに準拠
- キーボード操作に対応
- デスクトップ環境との統合
- デスクトップ環境（GNOME、KDE など）によって異なる

**例：**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## 一般的なパターン

### 破壊的な操作の実行前に確認する

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### ダイアログによるエラー処理

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### 検証を伴うファイル選択

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### 複数ステップのダイアログフロー

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## ベストプラクティス

### ✅ 推奨事項

- **ネイティブダイアログを使用する** — カスタムダイアログより優れた UX を提供できる
- **明確なメッセージを表示する** — 具体的に記述する
- **適切なタイトルを設定する** — コンテキストが重要
- **デフォルトボタンを適切に使用する** — 安全な選択肢をデフォルトにする
- **キャンセルを処理する** — ユーザーがキャンセルする可能性がある
- **選択されたファイルを検証する** — ファイル形式を確認する

### ❌ 禁止事項

- **ダイアログを多用しない** — ワークフローを中断してしまう
- **頻繁に表示するメッセージには使用しない** — 通知を使用する
- **エラー処理を忘れない** — ユーザーがキャンセルする可能性がある
- **不必要に処理をブロックしない** — 代替手段を検討する
- **曖昧なメッセージを使用しない** — 具体的に記述する
- **プラットフォーム間の違いを無視しない** — すべてのプラットフォームでテストする

## 次のステップ

@cards{cols="2"}
ℹ メッセージダイアログ
情報、警告、エラーの各ダイアログ。

[詳しく見る →](/features/dialogs/message/)

---
📖 ファイルダイアログ
ファイルを開く、保存する、およびフォルダーを選択する。

[詳しく見る →](/features/dialogs/file/)

---
◆ カスタムダイアログ
カスタムダイアログウィンドウを作成する。

[詳しく見る →](/features/dialogs/custom/)

---
▣ ウィンドウ
ウィンドウ管理について説明します。

[詳しく見る →](/features/windows/basics/)

@end

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[ダイアログの例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)を確認してください。
