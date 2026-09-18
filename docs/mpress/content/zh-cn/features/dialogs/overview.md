---
title: "对话框概述"
description: "在应用程序中显示原生系统对话框"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## 原生对话框

Wails 提供可在所有平台上使用的<strong>原生系统对话框</strong>：消息对话框（信息、警告、错误、问题）、文件对话框（打开、保存、文件夹），以及具有平台原生外观和行为的自定义对话框窗口。

![macOS 上带有“取消”和“丢弃”按钮的 Wails 问题对话框](/assets/screenshots/dialog-question-macos.png)

同一 API 会按照各个受支持平台的惯例进行呈现。此 macOS 示例展示了一个附加到其 Wails 窗口的问题对话框，其中包括默认按钮和取消按钮。

## 快速入门

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

<strong>就是这么简单！</strong>只需极少量代码即可使用原生对话框。

## 访问对话框

通过`app.Dialog`管理器访问对话框：

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## 对话框类型

### 信息对话框

显示简单消息：

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**使用场景：**

- 成功消息
- 信息通知
- 完成确认

### 警告对话框

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

### 错误对话框

显示错误：

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**使用场景：**

- 错误消息
- 失败通知
- 异常处理

### 问题对话框

向用户提问，并通过按钮回调处理响应：

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

**使用场景：**

- 确认操作
- 是/否问题
- 多项选择

## 文件对话框

### 打开文件对话框

选择要打开的文件：

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

**多选：**

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

### 保存文件对话框

选择保存位置：

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

### 选择文件夹对话框

启用目录选择后，使用打开文件对话框选择目录：

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

## 对话框选项

### 标题和消息

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### 按钮

**简单对话框的默认按钮：**

信息、警告和错误对话框会显示默认的“确定”按钮：

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**问题对话框的自定义按钮：**

使用`AddButton()`添加按钮。它会返回一个`*Button`，你可以为其配置回调：

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
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### 附加到窗口

将对话框附加到指定窗口：

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**行为：**

- 对话框显示在父窗口中央
- 显示对话框期间禁用父窗口
- 对话框随父窗口移动（macOS）

## 平台行为

@tabs{sync-key="platform"}
[macOS]
**macOS 对话框：**

- 原生 NSAlert 外观
- 跟随系统主题（浅色/深色）
- 支持键盘导航
- 标准快捷键（按 ⌘. 取消）
- 内置无障碍功能
- 附加到窗口时采用片式对话框样式

**示例：**

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
**Windows 对话框：**

- 原生 TaskDialog 外观
- 跟随系统主题
- 支持键盘导航
- 标准快捷键（按 Esc 取消）
- 内置无障碍功能
- 相对于父窗口为模态

**示例：**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Linux 对话框：**

- GTK 对话框外观
- 遵循桌面主题
- 支持键盘导航
- 与桌面环境集成
- 因桌面环境而异（GNOME、KDE 等）

**示例：**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## 常见模式

### 在执行破坏性操作前确认

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

### 使用对话框处理错误

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

### 选择文件并进行验证

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

### 多步骤对话框流程

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

## 最佳实践

### ✅ 应该做

- **使用原生对话框**——用户体验优于自定义对话框
- **提供清晰的消息**——内容要具体
- **设置恰当的标题**——上下文很重要
- **合理使用默认按钮**——将安全选项设为默认选项
- **处理取消操作**——用户可能会取消
- **验证所选文件**——检查文件类型

### ❌ 不应该做

- **不要过度使用对话框**——这会打断工作流程
- **不要用对话框显示频繁出现的消息**——应使用通知
- **不要忘记处理错误**——用户可能会取消
- **不要进行不必要的阻塞**——考虑其他方案
- **不要使用笼统的消息**——内容要具体
- **不要忽视平台差异**——在所有平台上进行测试

## 后续步骤

@cards{cols="2"}
ℹ 消息对话框
信息、警告和错误对话框。

[了解更多 →](/features/dialogs/message/)

---
📖 文件对话框
打开、保存和选择文件夹。

[了解更多 →](/features/dialogs/file/)

---
◆ 自定义对话框
创建自定义对话框窗口。

[了解更多 →](/features/dialogs/custom/)

---
▣ 窗口
了解窗口管理。

[了解更多 →](/features/windows/basics/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[对话框示例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
