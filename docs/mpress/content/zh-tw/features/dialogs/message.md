---
title: "訊息對話方塊"
description: "顯示資訊、警告、錯誤與問題"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## 訊息對話方塊

Wails 提供具有平台原生外觀的<strong>原生訊息對話方塊</strong>：包括資訊、警告、錯誤與問題對話方塊，並可自訂標題、訊息和按鈕。API 簡單、行為符合平台原生慣例，且預設具備無障礙支援。

## 建立對話方塊

透過`app.Dialog`管理器存取訊息對話方塊：

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

所有方法都會傳回`*MessageDialog`，可使用方法鏈進行設定。

## 資訊對話方塊

顯示資訊訊息：

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**使用情境：**

- 成功確認
- 完成通知
- 資訊訊息
- 狀態更新

**範例——儲存確認：**

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

## 警告對話方塊

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
- 潛在問題

**範例——磁碟空間警告：**

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

## 錯誤對話方塊

顯示錯誤：

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**使用情境：**

- 錯誤訊息
- 失敗通知
- 例外處理
- 重大問題

**範例——網路錯誤：**

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

## 問題對話方塊

向使用者提問，並透過按鈕回呼處理回應：

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

**使用情境：**

- 確認動作
- 是／否問題
- 多選一
- 使用者決策

**範例——未儲存的變更：**

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

## 對話方塊選項

### 標題與訊息

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**最佳實務：**

- <strong>標題：</strong>簡短且具描述性（2-5個單字）
- <strong>訊息：</strong>清楚、具體且可據以採取行動
- <strong>避免術語：</strong>使用淺白易懂的語言

### 按鈕

**單一按鈕（資訊／警告／錯誤）：**

資訊、警告與錯誤對話方塊會顯示預設的「確定」按鈕：

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

您也可以新增自訂按鈕：

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**多個按鈕（問題）：**

使用`AddButton()`新增按鈕。此方法會傳回可供設定的`*Button`：

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

**預設按鈕與取消按鈕：**

使用`SetDefaultButton()`指定反白顯示且按 Enter 鍵時觸發的按鈕。 使用`SetCancelButton()`指定按 Escape 鍵時觸發的按鈕。

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

您也可以對按鈕使用流暢介面的`SetAsDefault()`與`SetAsCancel()`方法：

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**最佳實務：**

- <strong>1-3個按鈕：</strong>不要讓使用者無所適從
- <strong>清楚的標籤：</strong>使用「儲存」，而非「確定」
- <strong>安全的預設選項：</strong>非破壞性動作
- <strong>順序很重要：</strong>最可能執行的動作排在最前面（取消除外）

### 自訂圖示

設定對話方塊的自訂圖示：

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### 附加至視窗

附加至指定視窗：

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**優點：**

- 對話方塊會顯示在正確的視窗上
- 顯示對話方塊期間會停用父視窗
- 提供更好的多視窗使用者體驗

## 完整範例

### 確認破壞性動作

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

### 結束確認

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

### 含下載選項的更新對話方塊

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

### 使用自訂圖示的問題對話方塊

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

## 最佳實務

### ✅ 建議做法

- **具體明確**—使用「檔案已儲存至『文件』」而非「成功」
- **使用適當的類型**—錯誤使用 Error，警告使用 Warning
- **提供脈絡**—包含相關詳細資訊
- **使用清楚的按鈕標籤**—使用「刪除」而非「確定」
- **設定安全的預設選項**—使用非破壞性操作
- **處理取消操作**—使用者可能會關閉對話方塊

### ❌ 不該做的事

- **不要過度使用**—這會中斷工作流程
- **不要用於頻繁更新**—改用通知
- **不要使用籠統的訊息**—「錯誤」無法提供任何資訊
- **不要忽略錯誤**—處理dialog.Show()錯誤
- **不要在非必要時阻塞**—考慮非同步替代方案
- **不要使用技術術語**—使用淺白易懂的語言

## 平台差異

### macOS

- 附加至視窗時會以表單式對話方塊顯示
- 標準鍵盤快速鍵（按⌘.取消）
- 自動遵循系統佈景主題
- 內建無障礙功能

### Windows

- 強制回應對話方塊
- TaskDialog外觀
- 按Esc取消
- 遵循系統佈景主題

### Linux

- GTK對話方塊
- 依桌面環境而異
- 遵循桌面佈景主題
- 標準鍵盤導覽

## 後續步驟

@cards{cols="2"}
📖 檔案對話方塊
開啟、儲存及選取資料夾。

[深入瞭解 →](/features/dialogs/file/)

---
◆ 自訂對話方塊
建立自訂對話方塊視窗。

[深入瞭解 →](/features/dialogs/custom/)

---
● 通知
不干擾操作的通知。

[深入瞭解 →](/features/notifications/overview/)

---
★ 事件
使用事件進行非阻塞式通訊。

[深入瞭解 →](/features/events/system/)

@end

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[對話方塊範例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
