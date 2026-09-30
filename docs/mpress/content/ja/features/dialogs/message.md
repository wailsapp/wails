---
title: "メッセージダイアログ"
description: "情報、警告、エラー、質問を表示する"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## メッセージダイアログ

Wails は、プラットフォームに適した外観の<strong>ネイティブメッセージダイアログ</strong>を提供します。情報、警告、エラー、質問の各ダイアログで、タイトル、メッセージ、ボタンをカスタマイズできます。API はシンプルで、ネイティブの動作を備え、デフォルトでアクセシブルです。

## ダイアログの作成

メッセージダイアログには、`app.Dialog`マネージャーを介してアクセスします。

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

すべてのメソッドは`*MessageDialog`を返し、メソッドチェーンを使用して設定できます。

## 情報ダイアログ

情報メッセージを表示します。

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**ユースケース：**

- 成功の確認
- 完了通知
- 情報メッセージ
- ステータスの更新

**例 — 保存の確認：**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## 警告ダイアログ

警告を表示します。

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**ユースケース：**

- 重大ではない警告
- 非推奨化の通知
- 注意メッセージ
- 潜在的な問題

**例 — ディスク容量の警告：**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## エラーダイアログ

エラーを表示します。

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**ユースケース：**

- エラーメッセージ
- 失敗の通知
- 例外処理
- 重大な問題

**例 — ネットワークエラー：**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## 質問ダイアログ

ユーザーに質問し、ボタンのコールバックを介して応答を処理します。

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**ユースケース：**

- 操作の確認
- 「はい」または「いいえ」で答える質問
- 選択式の質問
- ユーザーによる判断

**例 — 未保存の変更：**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## ダイアログのオプション

### タイトルとメッセージ

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**ベストプラクティス：**

- <strong>タイトル：</strong>短く、内容が分かるもの（2-5語）
- <strong>メッセージ：</strong>明確かつ具体的で、取るべき対応が分かるもの
- <strong>専門用語を避ける：</strong>平易な言葉を使用する

### ボタン

**単一ボタン（情報／警告／エラー）：**

情報、警告、エラーの各ダイアログには、デフォルトの「OK」ボタンが表示されます。

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

カスタムボタンも追加できます。

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**複数ボタン（質問）：**

`AddButton()`を使用してボタンを追加します。このメソッドは、設定可能な`*Button`を返します。

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

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
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

ボタンでは、流れるように記述できる`SetAsDefault()`メソッドと`SetAsCancel()`メソッドも使用できます。

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**ベストプラクティス：**

- <strong>1-3個のボタン：</strong>選択肢を増やしすぎてユーザーを困惑させない
- **明確なラベル：**「OK」ではなく「保存」
- <strong>安全なデフォルト：</strong>破壊的でない操作
- <strong>順序が重要：</strong>最も選ばれる可能性が高い操作を先頭にする（「キャンセル」を除く）

### カスタムアイコン

ダイアログにカスタムアイコンを設定します。

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### ウィンドウへの関連付け

特定のウィンドウに関連付けます。

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**利点：**

- 正しいウィンドウにダイアログが表示される
- 表示中は親ウィンドウが無効になる
- 複数ウィンドウでの UX が向上する

## 完全な例

### 破壊的な操作の確認

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### 終了の確認

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### ダウンロードオプション付きの更新ダイアログ

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### カスタムアイコン付きの質問ダイアログ

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## ベストプラクティス

### ✅ 推奨事項

- **具体的に記述する** - 「成功」ではなく「ファイルをDocumentsに保存しました」
- **適切な種類を使用する** - エラーにはError、警告にはWarningを使用する
- **コンテキストを示す** - 関連する詳細を含める
- **明確なボタンラベルを使用する** - 「OK」ではなく「削除」を使用する
- **安全な既定値を設定する** - 非破壊的な操作を既定にする
- **キャンセルを処理する** - ユーザーがダイアログを閉じる可能性がある

### ❌ 避けること

- **多用しない** - 作業の流れを中断する
- **頻繁な更新には使用しない** - 代わりに通知を使用する
- **漠然としたメッセージを使用しない** - 「エラー」だけでは何も伝わらない
- **エラーを無視しない** - dialog.Show()のエラーを処理する
- **不必要にブロックしない** - 非同期の代替手段を検討する
- **専門用語を使用しない** - 平易な言葉を使う

## プラットフォームごとの違い

### macOS

- ウィンドウに関連付けるとシート形式で表示
- 標準のキーボードショートカット（キャンセルは⌘.）
- システムテーマに自動的に従う
- アクセシビリティ機能を標準搭載

### Windows

- モーダルダイアログ
- TaskDialogの外観
- Escでキャンセル
- システムテーマに従う

### Linux

- GTKダイアログ
- デスクトップ環境によって異なる
- デスクトップテーマに従う
- 標準のキーボードナビゲーション

## 次のステップ

@cards{cols="2"}
📖 ファイルダイアログ
ファイルを開く、保存する、フォルダーを選択する。

[詳細を見る →](/features/dialogs/file/)

---
◆ カスタムダイアログ
カスタムダイアログウィンドウを作成する。

[詳細を見る →](/features/dialogs/custom/)

---
● 通知
作業を妨げない通知。

[詳細を見る →](/features/notifications/overview/)

---
★ イベント
ノンブロッキング通信にはイベントを使用する。

[詳細を見る →](/features/events/system/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[ダイアログのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)を確認してください。
