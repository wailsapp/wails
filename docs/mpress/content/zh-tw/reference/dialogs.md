---
title: "對話方塊 API"
description: "原生對話方塊 API 完整參考"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## 概觀

對話方塊 API 提供顯示原生檔案對話方塊與訊息對話方塊的方法。請透過 `app.Dialog` 管理器存取對話方塊。

**對話方塊類型：**

- **檔案對話方塊**：開啟與儲存對話方塊
- **訊息對話方塊**：資訊、錯誤、警告與詢問對話方塊

所有對話方塊都是符合平台外觀與操作體驗的<strong>原生作業系統對話方塊</strong>。

## 存取對話方塊

請透過 `app.Dialog` 管理器存取對話方塊：

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## 檔案對話方塊

### OpenFile()

建立檔案開啟對話方塊。

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**範例：**

```go
dialog := app.Dialog.OpenFile()
```

### OpenFileDialogStruct 方法

#### SetTitle()

設定對話方塊標題。

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**範例：**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

新增檔案類型篩選條件。

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**參數：**

- `displayName`：向使用者顯示的篩選條件說明（例如「圖片」、「文件」）
- `pattern`：以分號分隔的副檔名清單（例如「*.png;*.jpg」）

**範例：**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

設定初始目錄。

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**範例：**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

啟用或停用目錄選取功能。

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**範例（選取資料夾）：**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

啟用或停用檔案選取功能。

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

啟用或停用建立新目錄的功能。

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

顯示或隱藏隱藏檔案。

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

將對話方塊附加至特定視窗。

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

顯示對話方塊並傳回選取的檔案。

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**傳回值：**

- `string`：所選檔案的路徑。空字串應視為「未選取」（依作業系統而定，取消時各平台實作可能傳回空字串或非 nil 錯誤）。
- `error`：若對話方塊本身無法顯示，則為非 nil。

**範例：**

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

顯示對話方塊並傳回多個選取的檔案。

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**傳回值：**

- `[]string`：所選檔案路徑的陣列
- `error`：對話方塊失敗時的錯誤

**範例：**

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

建立檔案儲存對話方塊。

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**範例：**

```go
dialog := app.Dialog.SaveFile()
```

### SaveFileDialogStruct 方法

#### SetTitle()

設定對話方塊標題。

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

設定預設檔名。

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**範例：**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

新增檔案類型篩選條件。

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**範例：**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

設定初始目錄。

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

將對話方塊附加至特定視窗。

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

顯示對話方塊並傳回儲存路徑。

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**範例：**

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

### 選取資料夾

沒有獨立的`SelectFolderDialog`。請使用`OpenFile()`並設定目錄選項：

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

## 訊息對話方塊

所有訊息對話方塊都會傳回`*MessageDialog`，並共用相同的方法。

### Info()

建立資訊對話方塊。

```go
func (dm *DialogManager) Info() *MessageDialog
```

**範例：**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

建立錯誤對話方塊。

```go
func (dm *DialogManager) Error() *MessageDialog
```

**範例：**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

建立警告對話方塊。

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**範例：**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

建立具有自訂按鈕的詢問對話方塊。

```go
func (dm *DialogManager) Question() *MessageDialog
```

**範例：**

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

### MessageDialog 方法

#### SetTitle()

設定對話方塊標題。

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

設定對話方塊訊息。

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

設定對話方塊的自訂圖示。

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

將按鈕新增至對話方塊，並傳回該按鈕以供設定。

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**傳回值：**`*Button`－可進一步設定的按鈕執行個體

**範例：**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

設定預設按鈕（按下 Enter 時觸發該按鈕的動作）。

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**範例：**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

設定取消按鈕（按下 Escape 時觸發該按鈕的動作）。

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**範例：**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

將對話方塊附加至特定視窗。

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

顯示對話方塊。按鈕回呼會處理使用者的回應。

```go
func (d *MessageDialog) Show()
```

**注意：**`Show()`不會傳回值。請使用按鈕回呼來處理使用者的回應。

### 按鈕方法

#### OnClick()

設定按鈕遭點選時執行的回呼函式。

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

將此按鈕標記為預設按鈕。

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

將此按鈕標記為取消按鈕。

```go
func (b *Button) SetAsCancel() *Button
```

## 完整範例

### 檔案選取範例

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

### 確認對話方塊範例

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

### 儲存變更對話方塊

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

### 多檔案處理

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

### 使用對話方塊處理錯誤

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

### 平台特定的預設值

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

## 最佳實務

### 建議事項

- **使用原生對話方塊**－使其外觀與操作體驗符合平台風格
- **提供清楚的標題**－協助使用者瞭解用途
- **設定適當的篩選條件**－引導使用者選擇正確的檔案類型
- **處理取消操作**－檢查錯誤（使用者可能會取消）
- **對破壞性操作顯示確認提示**－使用 Question 對話方塊
- **提供回饋**－使用 Info 對話方塊顯示成功訊息
- **設定合理的預設值**－例如預設目錄、檔案名稱等
- **使用回呼處理按鈕操作**－正確處理使用者的回應

### 避免事項

- **不要忽略錯誤**－使用者取消操作時會傳回錯誤
- **不要使用語意不明的按鈕標籤**－請明確標示：「儲存」／「取消」
- **不要過度使用對話方塊**——它們會中斷工作流程
- **不要將取消操作顯示為錯誤**——這是正常操作
- **不要忘記設定檔案篩選條件**——協助使用者找到正確的檔案
- **不要將路徑寫死**——請使用 os.UserHomeDir() 或類似函式

## 各平台的對話方塊類型

### macOS

- 對話方塊從標題列向下滑出
- 附屬於父視窗的「表單」樣式
- 原生 macOS 外觀

### Windows

- 標準 Windows 對話方塊
- 遵循 Windows 設計準則
- 現代 Windows 10/11 外觀

### Linux

- 在以 GTK 為基礎的系統上使用 GTK 對話方塊
- 在以 Qt 為基礎的系統上使用 Qt 對話方塊
- 與桌面環境一致

#### Linux 對話方塊行為

在 Linux 上，預設的 GTK4 組建會使用 **xdg-desktop-portal** 顯示檔案對話方塊。這可提供原生桌面整合，但也表示某些選項不會生效。舊版 GTK3 路徑（`-tags gtk3`）仍可透過程式完整控制以下選項：

| 選項 | GTK3（`-tags gtk3`） | GTK4（預設） | 備註 |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ 可運作 | ❌ 不生效 | 由使用者透過對話方塊中的 UI 切換控制（Ctrl+H 或選單） |
| `CanCreateDirectories()` | ✅ 可運作 | ❌ 不生效 | 在入口服務中一律啟用 |
| `ResolvesAliases()` | ✅ 可運作 | ❌ 不生效 | 由入口服務處理符號連結解析 |
| `SetButtonText()` | ✅ 可運作 | ✅ 可運作 | 可使用自訂的接受按鈕文字 |

<strong>為何會有這些限制：</strong>GTK4 以入口服務為基礎的對話方塊會將 UI 控制權交由桌面環境（GNOME、KDE 等）處理。這是刻意的設計：入口服務可在不同應用程式間提供一致的使用者體驗，並遵循使用者偏好設定。

@note{type="info"}
預設的 GTK4 組建使用入口服務後端的對話方塊。如果您的應用程式需要透過程式完整控制上述對話方塊選項，請使用舊版 `-tags gtk3` 路徑進行組建（支援至 v3.0.x；已在 v3.1 移除）——請參閱[Linux 封裝——舊版 GTK3 支援](/guides/build/linux/#legacy-gtk3-support)。

@end

## 常見模式

### 「另存新檔」模式

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

### 「開啟最近使用的項目」模式

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
