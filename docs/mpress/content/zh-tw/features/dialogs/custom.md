---
title: "自訂對話方塊"
description: "在 Wails 應用程式中建立自訂對話方塊"
slug: "features/dialogs/custom"
sourcePath: "features/dialogs/custom.md"
---

## 自訂對話方塊

使用一般的 Wails 視窗建立具有對話方塊行為的<strong>自訂對話方塊視窗</strong>。您可以建立自訂表單、複雜的輸入驗證、品牌化外觀和豐富內容（圖片、影片），同時維持使用者熟悉的對話方塊模式。

## 快速入門

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

<strong>就是這麼簡單！</strong>具有對話方塊行為的自訂 UI。

## 建立自訂對話方塊

### 基本自訂對話方塊

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

### 強制回應對話方塊

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

### 表單對話方塊

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

## 對話方塊模式

### 確認對話方塊

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

### 輸入對話方塊

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

### 進度對話方塊

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

## 完整範例

### 登入對話方塊

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

### 設定對話方塊

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

### 精靈對話方塊

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

## 最佳實務

### ✅ 建議做法

- **使用適當的視窗選項**——例如 AlwaysOnTop、Frameless 等。
- **處理取消操作**——一律提供取消方式
- **驗證輸入**——接受資料前先檢查
- **提供回饋**——例如載入狀態和錯誤
- **使用事件進行通訊**——保持明確分離
- **清理資源**——關閉視窗並移除監聽器

### ❌ 避免事項

- **不要阻塞主執行緒**——使用通道傳回結果
- **不要忘記關閉**——否則會造成記憶體洩漏
- **不要略過驗證**——一律驗證輸入
- **不要忽略錯誤**——處理所有錯誤情況
- **不要設計得過於複雜**——讓對話方塊保持簡潔
- **不要忽略無障礙功能**——支援鍵盤導覽

## 設定自訂對話方塊的樣式

### 現代化對話方塊樣式

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

## 後續步驟

@cards{cols="2"}
ℹ 訊息對話方塊
標準的資訊、警告和錯誤對話方塊。

[深入瞭解 →](/features/dialogs/message/)

---
📖 檔案對話方塊
開啟、儲存及選取資料夾。

[深入瞭解 →](/features/dialogs/file/)

---
▣ 視窗
瞭解視窗管理。

[深入瞭解 →](/features/windows/basics/)

---
★ 事件
使用事件進行對話方塊通訊。

[深入瞭解 →](/features/events/system/)

@end

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[自訂對話方塊範例](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)。
