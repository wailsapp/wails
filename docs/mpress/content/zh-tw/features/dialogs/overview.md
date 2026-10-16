---
title: "對話方塊概觀"
description: "在應用程式中顯示原生系統對話方塊"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## 原生對話方塊

Wails 提供可跨所有平台運作的<strong>原生系統對話方塊</strong>：訊息對話方塊（資訊、警告、錯誤、問題）、檔案對話方塊（開啟、儲存、資料夾），以及具有平台原生外觀與行為的自訂對話方塊視窗。

![macOS 上具有「取消」和「捨棄」按鈕的 Wails 問題對話方塊](/assets/screenshots/dialog-question-macos.png)

同一套 API 會依各支援平台的慣例呈現。此 macOS 範例顯示附加至其 Wails 視窗的問題對話方塊，其中包含預設按鈕和取消按鈕。

## 快速入門

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

<strong>就是這麼簡單！</strong>只需極少量程式碼即可使用原生對話方塊。

## 存取對話方塊

透過`app.Dialog`管理器存取對話方塊：

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## 對話方塊類型

### 資訊對話方塊

顯示簡短訊息：

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**使用情境：**

- 成功訊息
- 資訊通知
- 完成確認

### 警告對話方塊

顯示警告：

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**使用情境：**

- 非重大警告
- 棄用通知
- 注意訊息

### 錯誤對話方塊

顯示錯誤：

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**使用情境：**

- 錯誤訊息
- 失敗通知
- 例外處理

### 問題對話方塊

向使用者提問，並透過按鈕回呼處理回應：

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

**使用情境：**

- 確認動作
- 是／否問題
- 多選一

## 檔案對話方塊

### 開啟檔案對話方塊

選取要開啟的檔案：

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

**多選：**

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

### 儲存檔案對話方塊

選擇儲存位置：

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

### 選取資料夾對話方塊

啟用開啟檔案對話方塊的目錄選取功能，以選擇目錄：

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

## 對話方塊選項

### 標題與訊息

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### 按鈕

**簡易對話方塊的預設按鈕：**

資訊、警告和錯誤對話方塊會顯示預設的「確定」按鈕：

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**問題對話方塊的自訂按鈕：**

使用`AddButton()`新增按鈕。此操作會傳回可設定回呼的`*Button`：

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

**預設按鈕和取消按鈕：**

使用`SetDefaultButton()`指定要醒目提示並在按下 Enter 時觸發的按鈕。 使用`SetCancelButton()`指定在按下 Escape 時觸發的按鈕。

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

### 附加至視窗

將對話方塊附加至特定視窗：

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**行為：**

- 對話方塊會顯示在父視窗中央
- 顯示對話方塊時，父視窗會停用
- 對話方塊會隨父視窗移動（macOS）

## 平台行為

@tabs{sync-key="platform"}
[macOS]
**macOS 對話方塊：**

- 原生 NSAlert 外觀
- 遵循系統主題（淺色／深色）
- 支援鍵盤導覽
- 標準快速鍵（按 ⌘. 取消）
- 內建無障礙功能
- 附加至視窗時採用附屬面板樣式

**範例：**

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
**Windows 對話方塊：**

- 原生 TaskDialog 外觀
- 遵循系統主題
- 支援鍵盤導覽
- 標準快速鍵（按 Esc 取消）
- 內建無障礙功能
- 相對於父視窗的模態視窗

**範例：**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Linux 對話方塊：**

- GTK 對話方塊外觀
- 遵循桌面佈景主題
- 支援鍵盤導覽
- 整合桌面環境
- 依桌面環境而異（GNOME、KDE 等）

**範例：**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## 常見模式

### 執行破壞性操作前先確認

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

### 使用對話方塊處理錯誤

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

### 選取檔案並進行驗證

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

### 多步驟對話方塊流程

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

## 最佳實務

### ✅ 建議做法

- **使用原生對話方塊** — 使用者體驗優於自訂對話方塊
- **提供清楚的訊息** — 內容要具體
- **設定適當的標題** — 情境很重要
- **審慎使用預設按鈕** — 將安全選項設為預設值
- **處理取消操作** — 使用者可能會取消
- **驗證選取的檔案** — 檢查檔案類型

### ❌ 避免的做法

- **不要過度使用對話方塊** — 這會中斷工作流程
- **不要用於頻繁出現的訊息** — 請改用通知
- **不要忘記處理錯誤** — 使用者可能會取消
- **不要在不必要時阻塞操作** — 考慮其他做法
- **不要使用籠統的訊息** — 內容要具體
- **不要忽略平台差異** — 在所有平台上進行測試

## 後續步驟

@cards{cols="2"}
ℹ 訊息對話方塊
資訊、警告和錯誤對話方塊。

[深入瞭解 →](/features/dialogs/message/)

---
📖 檔案對話方塊
開啟、儲存和資料夾選取。

[深入瞭解 →](/features/dialogs/file/)

---
◆ 自訂對話方塊
建立自訂對話方塊視窗。

[深入瞭解 →](/features/dialogs/custom/)

---
▣ 視窗
瞭解視窗管理。

[深入瞭解 →](/features/windows/basics/)

@end

---

<strong>有問題嗎？</strong>請前往[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[對話方塊範例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
