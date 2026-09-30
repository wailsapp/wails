---
title: "对话框 API"
description: "原生对话框 API 完整参考"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## 概述

对话框 API 提供用于显示原生文件对话框和消息对话框的方法。通过`app.Dialog`管理器访问对话框。

**对话框类型：**

- **文件对话框** - 打开和保存对话框
- **消息对话框** - 信息、错误、警告和询问对话框

所有对话框都是符合相应平台外观和体验的<strong>操作系统原生对话框</strong>。

## 访问对话框

通过`app.Dialog`管理器访问对话框：

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## 文件对话框

### OpenFile()

创建文件打开对话框。

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**示例：**

```go
dialog := app.Dialog.OpenFile()
```

### OpenFileDialogStruct 方法

#### SetTitle()

设置对话框标题。

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**示例：**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

添加文件类型筛选器。

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**参数：**

- `displayName` - 向用户显示的筛选器说明（例如“图像”“文档”）
- `pattern` - 以分号分隔的扩展名列表（例如“*.png;*.jpg”）

**示例：**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

设置初始目录。

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**示例：**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

启用或禁用目录选择。

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**示例（选择文件夹）：**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

启用或禁用文件选择。

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

启用或禁用新建目录。

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

显示或隐藏隐藏文件。

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

将对话框附加到指定窗口。

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

显示对话框并返回选中的文件。

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**返回值：**

- `string` - 选中文件的路径。空字符串应视为“未选择”（用户取消时，平台实现可能返回空字符串，也可能返回非 nil 错误，具体取决于操作系统）。
- `error` - 如果对话框本身未能显示，则为非 nil。

**示例：**

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

显示对话框并返回选中的多个文件。

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**返回值：**

- `[]string` - 选中文件路径的数组
- `error` - 对话框显示失败时返回的错误

**示例：**

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

创建文件保存对话框。

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**示例：**

```go
dialog := app.Dialog.SaveFile()
```

### SaveFileDialogStruct 方法

#### SetTitle()

设置对话框标题。

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

设置默认文件名。

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**示例：**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

添加文件类型筛选器。

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**示例：**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

设置初始目录。

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

将对话框附加到指定窗口。

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

显示对话框并返回保存路径。

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**示例：**

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

### 文件夹选择

没有单独的`SelectFolderDialog`。请使用`OpenFile()`并设置目录选项：

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

## 消息对话框

所有消息对话框都返回`*MessageDialog`，并且共享相同的方法。

### Info()

创建信息对话框。

```go
func (dm *DialogManager) Info() *MessageDialog
```

**示例：**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

创建错误对话框。

```go
func (dm *DialogManager) Error() *MessageDialog
```

**示例：**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

创建警告对话框。

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**示例：**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

创建带有自定义按钮的询问对话框。

```go
func (dm *DialogManager) Question() *MessageDialog
```

**示例：**

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

设置对话框标题。

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

设置对话框消息。

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

为对话框设置自定义图标。

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

向对话框添加按钮，并返回该按钮以供配置。

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**返回值：** `*Button`— 可供进一步配置的按钮实例

**示例：**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

设置默认按钮（按 Enter 键时会激活该按钮）。

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**示例：**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

设置取消按钮（按 Escape 键时会激活该按钮）。

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**示例：**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

将对话框附加到指定窗口。

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

显示对话框。按钮回调负责处理用户响应。

```go
func (d *MessageDialog) Show()
```

**注意：**`Show()`不返回值。请使用按钮回调处理用户响应。

### 按钮方法

#### OnClick()

设置单击按钮时调用的回调函数。

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

将此按钮标记为默认按钮。

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

将此按钮标记为取消按钮。

```go
func (b *Button) SetAsCancel() *Button
```

## 完整示例

### 文件选择示例

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

### 确认对话框示例

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

### 保存更改对话框

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

### 多文件处理

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

### 使用对话框处理错误

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

### 平台特定的默认设置

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

## 最佳实践

### 推荐做法

- **使用原生对话框**— 与平台的外观和体验保持一致
- **提供清晰的标题**— 帮助用户理解对话框的用途
- **设置适当的筛选器**— 引导用户选择正确的文件类型
- **处理取消操作**— 检查错误（用户可能会取消）
- **对破坏性操作显示确认提示**— 使用 Question 对话框
- **提供反馈**— 使用 Info 对话框显示成功消息
- **设置合理的默认值**— 例如默认目录、文件名等
- **使用回调处理按钮操作**— 正确处理用户响应

### 避免的做法

- **不要忽略错误**— 用户取消操作时会返回错误
- **不要使用含义模糊的按钮标签**— 应明确标注：“保存”/“取消”
- **不要过度使用对话框**——它们会打断工作流程
- **不要在用户取消操作时显示错误**——取消是正常操作
- **不要忘记设置文件筛选器**——帮助用户找到正确的文件
- **不要硬编码路径**——请使用os.UserHomeDir()或类似函数

## 各平台的对话框类型

### macOS

- 对话框从标题栏向下滑出
- 附加到父窗口的“表单”样式
- 原生macOS外观

### Windows

- 标准Windows对话框
- 遵循Windows设计准则
- 现代Windows 10/11外观

### Linux

- 基于GTK的系统使用GTK对话框
- 基于Qt的系统使用Qt对话框
- 与桌面环境保持一致

#### Linux对话框行为

在Linux上，默认的GTK4构建使用<strong>xdg-desktop-portal</strong>实现文件对话框。这提供了原生桌面集成，但也意味着某些选项不会生效。旧版GTK3路径（`-tags gtk3`）仍可对这些选项进行完整的编程控制：

| 选项 | GTK3（`-tags gtk3`） | GTK4（默认） | 备注 |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ 生效 | ❌ 不生效 | 用户通过对话框的界面开关控制（Ctrl+H或菜单） |
| `CanCreateDirectories()` | ✅ 生效 | ❌ 不生效 | 在门户中始终启用 |
| `ResolvesAliases()` | ✅ 生效 | ❌ 不生效 | 门户负责解析符号链接 |
| `SetButtonText()` | ✅ 生效 | ✅ 生效 | 自定义接受按钮文本可正常生效 |

<strong>存在这些限制的原因：</strong>GTK4基于门户的对话框将界面控制权交由桌面环境（GNOME、KDE等）。这是有意的设计——门户可在不同应用程序之间提供一致的用户体验，并遵循用户偏好。

@note{type="info"}
默认的GTK4构建使用由门户支持的对话框。如果应用程序需要对上述对话框选项进行完整的编程控制，请使用旧版`-tags gtk3`路径进行构建（支持至v3.0.x；已在v3.1中移除）——请参阅[Linux打包——旧版GTK3支持](/guides/build/linux/#legacy-gtk3-support)。

@end

## 常见模式

### “另存为”模式

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

### “打开最近使用的文件”模式

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
