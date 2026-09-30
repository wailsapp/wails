---
title: "消息对话框"
description: "显示信息、警告、错误和问题"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## 消息对话框

Wails 提供具有平台原生外观的<strong>原生消息对话框</strong>：包括信息、警告、错误和问题对话框，并可自定义标题、消息和按钮。API 简洁、行为原生，并且默认支持无障碍访问。

## 创建对话框

通过`app.Dialog`管理器访问消息对话框：

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

所有方法都会返回一个`*MessageDialog`，可通过方法链进行配置。

## 信息对话框

显示信息性消息：

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**使用场景：**

- 成功确认
- 完成通知
- 信息性消息
- 状态更新

**示例——保存确认：**

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

## 警告对话框

显示警告：

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**使用场景：**

- 非严重警告
- 弃用通知
- 注意消息
- 潜在问题

**示例——磁盘空间警告：**

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

## 错误对话框

显示错误：

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**使用场景：**

- 错误消息
- 失败通知
- 异常处理
- 严重问题

**示例——网络错误：**

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

## 问题对话框

向用户提问，并通过按钮回调处理响应：

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

**使用场景：**

- 确认操作
- 是/否问题
- 多项选择
- 用户决策

**示例——未保存的更改：**

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

## 对话框选项

### 标题和消息

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**最佳实践：**

- <strong>标题：</strong>简短且含义明确（2-5个词）
- <strong>消息：</strong>清晰、具体且可据此采取行动
- <strong>避免使用行话：</strong>使用浅显易懂的语言

### 按钮

**单个按钮（信息/警告/错误）：**

信息、警告和错误对话框会显示默认的“确定”按钮：

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

也可以添加自定义按钮：

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**多个按钮（问题）：**

使用`AddButton()`添加按钮；该方法会返回一个可供配置的`*Button`：

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

**默认按钮和取消按钮：**

使用`SetDefaultButton()`指定要突出显示且按 Enter 键时触发的按钮。 使用`SetCancelButton()`指定按 Escape 键时触发的按钮。

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

也可以对按钮使用流式`SetAsDefault()`和`SetAsCancel()`方法：

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**最佳实践：**

- <strong>1-3个按钮：</strong>不要让过多选项困扰用户
- <strong>明确的标签：</strong>使用“保存”，而不是“确定”
- <strong>安全的默认选项：</strong>非破坏性操作
- <strong>顺序很重要：</strong>最可能执行的操作放在最前面（“取消”除外）

### 自定义图标

为对话框设置自定义图标：

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### 关联窗口

关联到特定窗口：

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**优点：**

- 对话框显示在正确的窗口上
- 对话框显示期间禁用父窗口
- 改善多窗口用户体验

## 完整示例

### 确认破坏性操作

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

### 退出确认

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

### 带下载选项的更新对话框

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

### 带自定义图标的问题对话框

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

## 最佳实践

### ✅ 推荐做法

- **具体明确** — 使用“文件已保存到文档文件夹”，而不是“成功”
- **使用适当的类型** — 错误使用 Error，警告使用 Warning
- **提供上下文** — 包含相关详细信息
- **使用明确的按钮标签** — 使用“删除”，而不是“确定”
- **设置安全的默认操作** — 使用非破坏性操作
- **处理取消操作** — 用户可能会关闭对话框

### ❌ 不要这样做

- **不要过度使用** — 这会打断工作流程
- **不要用于频繁更新** — 应改用通知
- **不要使用笼统的消息** — “错误”无法提供任何有用信息
- **不要忽略错误** — 处理 dialog.Show() 返回的错误
- **不要进行不必要的阻塞** — 考虑异步替代方案
- **不要使用技术术语** — 使用通俗易懂的语言

## 平台差异

### macOS

- 附加到窗口时采用表单式对话框
- 标准键盘快捷键（按 ⌘. 取消）
- 自动跟随系统主题
- 内置无障碍功能

### Windows

- 模态对话框
- TaskDialog 外观
- 按 Esc 取消
- 跟随系统主题

### Linux

- GTK 对话框
- 因桌面环境而异
- 跟随桌面主题
- 标准键盘导航

## 后续步骤

@cards{cols="2"}
📖 文件对话框
打开、保存和选择文件夹。

[了解更多 →](/features/dialogs/file/)

---
◆ 自定义对话框
创建自定义对话框窗口。

[了解更多 →](/features/dialogs/custom/)

---
● 通知
非侵入式通知。

[了解更多 →](/features/notifications/overview/)

---
★ 事件
使用事件进行非阻塞通信。

[了解更多 →](/features/events/system/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[对话框示例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
