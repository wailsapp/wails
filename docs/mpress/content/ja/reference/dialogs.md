---
title: "ダイアログ API"
description: "ネイティブダイアログ API の完全なリファレンス"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## 概要

Dialogs API には、ネイティブのファイルダイアログとメッセージダイアログを表示するためのメソッドが用意されています。ダイアログには `app.Dialog` マネージャーを介してアクセスします。

**ダイアログの種類：**

- **ファイルダイアログ** - 開くダイアログと保存ダイアログ
- **メッセージダイアログ** - 情報、エラー、警告、質問の各ダイアログ

すべてのダイアログは、プラットフォームの外観と操作感に合った<strong>ネイティブ OS ダイアログ</strong>です。

## ダイアログへのアクセス

ダイアログには `app.Dialog` マネージャーを介してアクセスします：

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## ファイルダイアログ

### OpenFile()

ファイルを開くダイアログを作成します。

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**例：**

```go
dialog := app.Dialog.OpenFile()
```

### OpenFileDialogStruct のメソッド

#### SetTitle()

ダイアログのタイトルを設定します。

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**例：**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

ファイル形式フィルターを追加します。

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**パラメーター：**

- `displayName` - ユーザーに表示するフィルターの説明（例：「画像」、「ドキュメント」）
- `pattern` - セミコロンで区切った拡張子のリスト（例：「*.png;*.jpg」）

**例：**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

初期ディレクトリを設定します。

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**例：**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

ディレクトリの選択を有効または無効にします。

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**例（フォルダーの選択）：**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

ファイルの選択を有効または無効にします。

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

新しいディレクトリの作成を有効または無効にします。

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

隠しファイルの表示または非表示を切り替えます。

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

ダイアログを指定したウィンドウに関連付けます。

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

ダイアログを表示し、選択されたファイルを返します。

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**戻り値：**

- `string` - 選択されたファイルのパス。空文字列は「選択なし」として扱ってください（キャンセル時には、OS によってプラットフォーム実装が空文字列または nil ではないエラーのいずれかを返す場合があります）。
- `error` - ダイアログ自体を表示できなかった場合は nil ではない値

**例：**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

ダイアログを表示し、選択された複数のファイルを返します。

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**戻り値：**

- `[]string` - 選択されたファイルパスの配列
- `error` - ダイアログの表示に失敗した場合のエラー

**例：**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

ファイル保存ダイアログを作成します。

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**例：**

```go
dialog := app.Dialog.SaveFile()
```

### SaveFileDialogStruct のメソッド

#### SetTitle()

ダイアログのタイトルを設定します。

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

デフォルトのファイル名を設定します。

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**例：**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

ファイル形式フィルターを追加します。

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**例：**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

初期ディレクトリを設定します。

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

ダイアログを指定したウィンドウに関連付けます。

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

ダイアログを表示し、保存先のパスを返します。

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**例：**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### フォルダーの選択

専用の`SelectFolderDialog`はありません。ディレクトリ用のオプションを指定して`OpenFile()`を使用します：

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## メッセージダイアログ

すべてのメッセージダイアログは`*MessageDialog`を返し、同じメソッドを共有します。

### Info()

情報ダイアログを作成します。

```go
func (dm *DialogManager) Info() *MessageDialog
```

**例：**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

エラーダイアログを作成します。

```go
func (dm *DialogManager) Error() *MessageDialog
```

**例：**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

警告ダイアログを作成します。

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**例：**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

カスタムボタンを備えた質問ダイアログを作成します。

```go
func (dm *DialogManager) Question() *MessageDialog
```

**例：**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### MessageDialogのメソッド

#### SetTitle()

ダイアログのタイトルを設定します。

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

ダイアログのメッセージを設定します。

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

ダイアログにカスタムアイコンを設定します。

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

ダイアログにボタンを追加し、設定用のボタンを返します。

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**戻り値：** `*Button` - 追加設定に使用するボタンのインスタンス

**例：**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

デフォルトのボタン（Enterキーを押すと実行されるボタン）を設定します。

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**例：**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

キャンセルボタン（Escapeキーを押すと実行されるボタン）を設定します。

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**例：**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

ダイアログを特定のウィンドウに関連付けます。

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

ダイアログを表示します。ユーザーの応答はボタンのコールバックで処理します。

```go
func (d *MessageDialog) Show()
```

**注：** `Show()`は値を返しません。ユーザーの応答はボタンのコールバックで処理してください。

### ボタンのメソッド

#### OnClick()

ボタンがクリックされたときに実行するコールバック関数を設定します。

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

このボタンをデフォルトボタンとして設定します。

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

このボタンをキャンセルボタンとして設定します。

```go
func (b *Button) SetAsCancel() *Button
```

## 完全な例

### ファイル選択の例

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### 確認ダイアログの例

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### 変更保存ダイアログ

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### 複数ファイルの処理

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### ダイアログを使用したエラー処理

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### プラットフォーム固有のデフォルト設定

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## ベストプラクティス

### 推奨事項

- **ネイティブダイアログを使用する** - プラットフォームの外観と操作感に合います
- **分かりやすいタイトルを付ける** - ユーザーが目的を理解しやすくなります
- **適切なフィルターを設定する** - ユーザーを正しいファイル形式へ誘導します
- **キャンセルを処理する** - エラーを確認してください（ユーザーがキャンセルする場合があります）
- **破壊的な操作では確認を表示する** - Questionダイアログを使用します
- **フィードバックを提供する** - 成功メッセージにはInfoダイアログを使用します
- **適切なデフォルト値を設定する** - デフォルトのディレクトリやファイル名などを設定します
- **ボタン操作にコールバックを使用する** - ユーザーの応答を適切に処理します

### 避けるべきこと

- **エラーを無視しない** - ユーザーがキャンセルするとエラーが返されます
- **曖昧なボタンラベルを使用しない** - 「保存」や「キャンセル」のように具体的にします
- **ダイアログを多用しない** - 作業の流れを中断します
- **キャンセル時にエラーを表示しない** - キャンセルは通常の操作です
- **ファイルフィルターを忘れない** - ユーザーが適切なファイルを見つけやすくなります
- **パスをハードコードしない** - os.UserHomeDir() などを使用してください

## プラットフォーム別のダイアログ形式

### macOS

- ダイアログはタイトルバーから下方向にスライドして表示されます
- 親ウィンドウに付随する「シート」形式
- macOS ネイティブの外観

### Windows

- 標準の Windows ダイアログ
- Windows のデザインガイドラインに準拠
- モダンな Windows 10/11 の外観

### Linux

- GTK ベースのシステムでは GTK ダイアログを使用
- Qt ベースのシステムでは Qt ダイアログを使用
- デスクトップ環境に適合

#### Linux でのダイアログの動作

Linux では、デフォルトの GTK4 ビルドはファイルダイアログに **xdg-desktop-portal** を使用します。これによりデスクトップとのネイティブな統合が提供されますが、一部のオプションは効果がありません。従来の GTK3 パス（`-tags gtk3`）では、これらのオプションをプログラムから完全に制御できます：

| オプション | GTK3（`-tags gtk3`） | GTK4（デフォルト） | 備考 |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ 動作する | ❌ 効果なし | ユーザーがダイアログの UI にある切り替え機能（Ctrl+H またはメニュー）で制御 |
| `CanCreateDirectories()` | ✅ 動作する | ❌ 効果なし | ポータルでは常に有効 |
| `ResolvesAliases()` | ✅ 動作する | ❌ 効果なし | シンボリックリンクの解決はポータルが処理 |
| `SetButtonText()` | ✅ 動作する | ✅ 動作する | 承認ボタンのテキストをカスタマイズ可能 |

**これらの制限が存在する理由：** GTK4 のポータルベースのダイアログは、UI の制御をデスクトップ環境（GNOME、KDE など）に委ねます。これは仕様です。ポータルはアプリケーション間で一貫した UX を提供し、ユーザー設定を尊重します。

@note{type="info"}
デフォルトの GTK4 ビルドでは、ポータルを利用するダイアログが使用されます。アプリケーションで上記のダイアログオプションをプログラムから完全に制御する必要がある場合は、従来の `-tags gtk3` パスでビルドしてください（v3.0.x までサポートされ、v3.1 で削除）— [Linux のパッケージ化 - 従来の GTK3 サポート](/guides/build/linux/#legacy-gtk3-support)を参照してください。

@end

## 一般的なパターン

### 「名前を付けて保存」パターン

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### 「最近使ったファイルを開く」パターン

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
