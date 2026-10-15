---
title: "自定义对话框"
description: "在 Wails 应用程序中创建自定义对话框"
slug: "features/dialogs/custom"
sourcePath: "features/dialogs/custom.md"
---

## 自定义对话框

使用具有对话框式行为的常规 Wails 窗口创建<strong>自定义对话框窗口</strong>。在保持用户熟悉的对话框交互模式的同时，构建自定义表单、复杂的输入验证、品牌化外观和富媒体内容（图像、视频）。

## 快速开始

```go
// Create custom dialog window
dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:       "Custom dialog",
    Width:       400,
    Height:      300,
    AlwaysOnTop: true,
    Frameless:   true,
    Hidden:      true,
})

// Load custom UI
dialog.SetURL("http://wails.localhost/dialog.html")

// Show as modal
dialog.Show()
dialog.Focus()
```

<strong>就是这样！</strong>您现在拥有了具备对话框行为的自定义 UI。

## 创建自定义对话框

### 基本自定义对话框

```go
type Customdialog struct {
    window *application.WebviewWindow
    result chan string
}

func NewCustomdialog(app *application.App) *Customdialog {
    dialog := &Customdialog{
        result: make(chan string, 1),
    }
    
    dialog.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         "Custom dialog",
        Width:         400,
        Height:        300,
        AlwaysOnTop:   true,
        DisableResize: true,
        Hidden:        true,
    })
    
    return dialog
}

func (d *Customdialog) Show() string {
    d.window.Show()
    d.window.Focus()
    
    // Wait for result
    return <-d.result
}

func (d *Customdialog) Close(result string) {
    d.result <- result
    d.window.Close()
}
```

### 模态对话框

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) string {
    // Create dialog
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         title,
        Width:         400,
        Height:        200,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    // Attach as a child modal of the parent.
    parent.AttachModal(dialog)

    // Disable parent
    parent.SetEnabled(false)

    // Re-enable parent on close. RegisterHook so the cleanup runs
    // before the window is actually torn down.
    dialog.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        parent.SetEnabled(true)
        parent.Focus()
    })

    dialog.Show()

    return waitForResult(dialog)
}
```

### 表单对话框

```go
type Formdialog struct {
    window *application.WebviewWindow
    data   map[string]interface{}
    done   chan bool
}

func NewFormdialog(app *application.App) *Formdialog {
    fd := &Formdialog{
        data: make(map[string]interface{}),
        done: make(chan bool, 1),
    }
    
    fd.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Enter Information",
        Width:     500,
        Height:    400,
        Frameless: true,
        Hidden:    true,
    })
    
    return fd
}

func (fd *Formdialog) Show() (map[string]interface{}, bool) {
    fd.window.Show()
    fd.window.Focus()
    
    ok := <-fd.done
    return fd.data, ok
}

func (fd *Formdialog) Submit(data map[string]interface{}) {
    fd.data = data
    fd.done <- true
    fd.window.Close()
}

func (fd *Formdialog) Cancel() {
    fd.done <- false
    fd.window.Close()
}
```

## 对话框模式

### 确认对话框

```go
func ShowConfirmdialog(message string) bool {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:       "Confirm",
        Width:       400,
        Height:      150,
        AlwaysOnTop: true,
        Frameless:   true,
    })
    
    // Pass message to dialog once the in-window runtime is ready.
    dialog.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
        dialog.EmitEvent("set-message", message)
    })
    
    result := make(chan bool, 1)
    
    // Handle responses
    app.Event.On("confirm-yes", func(e *application.CustomEvent) {
        result <- true
        dialog.Close()
    })
    
    app.Event.On("confirm-no", func(e *application.CustomEvent) {
        result <- false
        dialog.Close()
    })
    
    dialog.Show()
    return <-result
}
```

**前端（HTML/JS）：**

```html
<div class="dialog">
    <h2 id="message"></h2>
    <div class="buttons">
        <button onclick="confirm(true)">Yes</button>
        <button onclick="confirm(false)">No</button>
    </div>
</div>

<script>
import { Events } from '@wailsio/runtime'

Events.On("set-message", (message) => {
    document.getElementById("message").textContent = message
})

function confirm(result) {
    Events.Emit(result ? "confirm-yes" : "confirm-no")
}
</script>
```

### 输入对话框

```go
func ShowInputdialog(prompt string, defaultValue string) (string, bool) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Input",
        Width:     400,
        Height:    150,
        Frameless: true,
    })
    
    result := make(chan struct {
        value string
        ok    bool
    }, 1)
    
    dialog.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
        dialog.EmitEvent("set-prompt", map[string]string{
            "prompt":  prompt,
            "default": defaultValue,
        })
    })
    
    app.Event.On("input-submit", func(e *application.CustomEvent) {
        result <- struct {
            value string
            ok    bool
        }{e.Data.(string), true}
        dialog.Close()
    })
    
    app.Event.On("input-cancel", func(e *application.CustomEvent) {
        result <- struct {
            value string
            ok    bool
        }{"", false}
        dialog.Close()
    })
    
    dialog.Show()
    r := <-result
    return r.value, r.ok
}
```

### 进度对话框

```go
type Progressdialog struct {
    window *application.WebviewWindow
}

func NewProgressdialog(title string) *Progressdialog {
    pd := &Progressdialog{}
    
    pd.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     title,
        Width:     400,
        Height:    150,
        Frameless: true,
    })
    
    return pd
}

func (pd *Progressdialog) Show() {
    pd.window.Show()
}

func (pd *Progressdialog) UpdateProgress(current, total int, message string) {
    pd.window.EmitEvent("progress-update", map[string]interface{}{
        "current": current,
        "total":   total,
        "message": message,
    })
}

func (pd *Progressdialog) Close() {
    pd.window.Close()
}
```

**用法：**

```go
func processFiles(files []string) {
    progress := NewProgressdialog("Processing Files")
    progress.Show()
    
    for i, file := range files {
        progress.UpdateProgress(i+1, len(files), 
            fmt.Sprintf("Processing %s...", filepath.Base(file)))
        
        processFile(file)
    }
    
    progress.Close()
}
```

## 完整示例

### 登录对话框

**Go：**

```go
type Logindialog struct {
    window *application.WebviewWindow
    result chan struct {
        username string
        password string
        ok       bool
    }
}

func NewLogindialog(app *application.App) *Logindialog {
    ld := &Logindialog{
        result: make(chan struct {
            username string
            password string
            ok       bool
        }, 1),
    }
    
    ld.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Login",
        Width:     400,
        Height:    250,
        Frameless: true,
    })
    
    return ld
}

func (ld *Logindialog) Show() (string, string, bool) {
    ld.window.Show()
    ld.window.Focus()
    
    r := <-ld.result
    return r.username, r.password, r.ok
}

func (ld *Logindialog) Submit(username, password string) {
    ld.result <- struct {
        username string
        password string
        ok       bool
    }{username, password, true}
    ld.window.Close()
}

func (ld *Logindialog) Cancel() {
    ld.result <- struct {
        username string
        password string
        ok       bool
    }{"", "", false}
    ld.window.Close()
}
```

**前端：**

```html
<div class="login-dialog">
    <h2>Login</h2>
    <form id="login-form">
        <input type="text" id="username" placeholder="Username" required>
        <input type="password" id="password" placeholder="Password" required>
        <div class="buttons">
            <button type="submit">Login</button>
            <button type="button" onclick="cancel()">Cancel</button>
        </div>
    </form>
</div>

<script>
import { Events } from '@wailsio/runtime'

document.getElementById('login-form').addEventListener('submit', (e) => {
    e.preventDefault()
    const username = document.getElementById('username').value
    const password = document.getElementById('password').value
    Events.Emit('login-submit', { username, password })
})

function cancel() {
    Events.Emit('login-cancel')
}
</script>
```

### 设置对话框

**Go：**

```go
type Settingsdialog struct {
    window   *application.WebviewWindow
    settings map[string]interface{}
    done     chan bool
}

func NewSettingsdialog(app *application.App, current map[string]interface{}) *Settingsdialog {
    sd := &Settingsdialog{
        settings: current,
        done:     make(chan bool, 1),
    }
    
    sd.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Settings",
        Width:  600,
        Height: 500,
    })
    
    sd.window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
        sd.window.EmitEvent("load-settings", current)
    })
    
    return sd
}

func (sd *Settingsdialog) Show() (map[string]interface{}, bool) {
    sd.window.Show()
    
    ok := <-sd.done
    return sd.settings, ok
}

func (sd *Settingsdialog) Save(settings map[string]interface{}) {
    sd.settings = settings
    sd.done <- true
    sd.window.Close()
}

func (sd *Settingsdialog) Cancel() {
    sd.done <- false
    sd.window.Close()
}
```

### 向导对话框

```go
type Wizarddialog struct {
    window      *application.WebviewWindow
    currentStep int
    data        map[string]interface{}
    done        chan bool
}

func NewWizarddialog(app *application.App) *Wizarddialog {
    wd := &Wizarddialog{
        currentStep: 0,
        data:        make(map[string]interface{}),
        done:        make(chan bool, 1),
    }
    
    wd.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         "Setup Wizard",
        Width:         600,
        Height:        400,
        DisableResize: true,
    })
    
    return wd
}

func (wd *Wizarddialog) Show() (map[string]interface{}, bool) {
    wd.window.Show()
    
    ok := <-wd.done
    return wd.data, ok
}

func (wd *Wizarddialog) NextStep(stepData map[string]interface{}) {
    // Merge step data
    for k, v := range stepData {
        wd.data[k] = v
    }
    
    wd.currentStep++
    wd.window.EmitEvent("next-step", wd.currentStep)
}

func (wd *Wizarddialog) PreviousStep() {
    if wd.currentStep > 0 {
        wd.currentStep--
        wd.window.EmitEvent("previous-step", wd.currentStep)
    }
}

func (wd *Wizarddialog) Finish(finalData map[string]interface{}) {
    for k, v := range finalData {
        wd.data[k] = v
    }
    
    wd.done <- true
    wd.window.Close()
}

func (wd *Wizarddialog) Cancel() {
    wd.done <- false
    wd.window.Close()
}
```

## 最佳实践

### ✅ 应该做

- **使用适当的窗口选项**——例如 AlwaysOnTop、Frameless 等。
- **处理取消操作**——始终提供取消方式
- **验证输入**——接受数据前先进行检查
- **提供反馈**——例如加载状态和错误信息
- **使用事件进行通信**——保持清晰的职责分离
- **清理资源**——关闭窗口并移除监听器

### ❌ 不应该做

- **不要阻塞主线程**——使用通道传递结果
- **不要忘记关闭窗口**——否则会导致内存泄漏
- **不要跳过验证**——始终验证输入
- **不要忽略错误**——处理所有错误情况
- **不要使其过于复杂**——保持对话框简洁
- **不要忘记无障碍支持**——支持键盘导航

## 设置自定义对话框的样式

### 现代对话框样式

```css
.dialog {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: white;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.dialog-header {
    --wails-draggable: drag;
    padding: 16px;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}

.dialog-content {
    flex: 1;
    padding: 24px;
    overflow: auto;
}

.dialog-footer {
    padding: 16px;
    background: #f5f5f5;
    border-top: 1px solid #e0e0e0;
    display: flex;
    justify-content: flex-end;
    gap: 8px;
}

button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
}

button.primary {
    background: #007aff;
    color: white;
}

button.secondary {
    background: #e0e0e0;
    color: #333;
}
```

## 后续步骤

@cards{cols="2"}
ℹ 消息对话框
标准信息、警告和错误对话框。

[了解更多 →](/features/dialogs/message/)

---
📖 文件对话框
打开、保存和文件夹选择。

[了解更多 →](/features/dialogs/file/)

---
▣ 窗口
了解窗口管理。

[了解更多 →](/features/windows/basics/)

---
★ 事件
使用事件实现对话框通信。

[了解更多 →](/features/events/system/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[自定义对话框示例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
