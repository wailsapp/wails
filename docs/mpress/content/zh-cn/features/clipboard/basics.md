---
title: "剪贴板操作"
description: "使用系统剪贴板复制和粘贴文本"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## 剪贴板操作

Wails 提供了一个<strong>统一的剪贴板 API</strong>，可在所有平台上使用。通过简单且一致的方法，在 Windows、macOS 和 Linux 上复制和粘贴文本。

## 快速开始

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

<strong>就是这么简单！</strong>实现跨平台剪贴板访问。

## 复制文本

### 基本复制

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**返回值：**`bool`——成功时为`true`，否则为`false`

### 从服务复制

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**从 JavaScript 调用：**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### 复制并提供反馈

```go
func copyWithFeedback(text string) {
    if app.Clipboard.SetText(text) {
        app.Dialog.Info().
            SetTitle("Copied").
            SetMessage("Text copied to clipboard!").
            Show()
    } else {
        app.Dialog.Error().
            SetTitle("Copy Failed").
            SetMessage("Failed to copy to clipboard.").
            Show()
    }
}
```

## 粘贴文本

### 基本粘贴

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**返回值：**`(string, bool)`——文本和成功标志

### 从服务粘贴

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**从 JavaScript 调用：**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### 粘贴并验证

```go
func pasteText() (string, error) {
    text, ok := app.Clipboard.Text()
    if !ok {
        return "", errors.New("clipboard empty or unavailable")
    }
    
    // Validate
    if len(text) == 0 {
        return "", errors.New("clipboard is empty")
    }
    
    if len(text) > 10000 {
        return "", errors.New("clipboard text too large")
    }
    
    return text, nil
}
```

## 完整示例

### 复制按钮

**Go：**

```go
type TextService struct {
    app *application.App
}

func (t *TextService) CopyText(text string) error {
    if !t.app.Clipboard.SetText(text) {
        return errors.New("failed to copy")
    }
    return nil
}
```

**JavaScript：**

```javascript
import { CopyText } from './bindings/changeme/textservice'

async function copyToClipboard(text) {
    try {
        await CopyText(text)
        showNotification("Copied to clipboard!")
    } catch (error) {
        showError("Failed to copy: " + error)
    }
}

// Usage
document.getElementById('copy-btn').addEventListener('click', () => {
    const text = document.getElementById('text').value
    copyToClipboard(text)
})
```

### 粘贴并处理

**Go：**

```go
type DataService struct {
    app *application.App
}

func (d *DataService) PasteAndProcess() (string, error) {
    // Get clipboard text
    text, ok := d.app.Clipboard.Text()
    if !ok {
        return "", errors.New("clipboard unavailable")
    }
    
    // Process text
    processed := strings.TrimSpace(text)
    processed = strings.ToUpper(processed)
    
    return processed, nil
}
```

**JavaScript：**

```javascript
import { PasteAndProcess } from './bindings/changeme/dataservice'

async function pasteAndProcess() {
    try {
        const result = await PasteAndProcess()
        document.getElementById('output').value = result
    } catch (error) {
        showError("Failed to paste: " + error)
    }
}
```

### 复制多种格式

```go
type CopyService struct {
    app *application.App
}

func (c *CopyService) CopyAsPlainText(text string) bool {
    return c.app.Clipboard.SetText(text)
}

func (c *CopyService) CopyAsJSON(data interface{}) bool {
    jsonBytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return false
    }
    return c.app.Clipboard.SetText(string(jsonBytes))
}

func (c *CopyService) CopyAsCSV(rows [][]string) bool {
    var buf bytes.Buffer
    writer := csv.NewWriter(&buf)
    
    for _, row := range rows {
        if err := writer.Write(row); err != nil {
            return false
        }
    }
    
    writer.Flush()
    return c.app.Clipboard.SetText(buf.String())
}
```

### 剪贴板监视器

```go
type ClipboardMonitor struct {
    app          *application.App
    lastText     string
    ticker       *time.Ticker
    stopChan     chan bool
}

func NewClipboardMonitor(app *application.App) *ClipboardMonitor {
    return &ClipboardMonitor{
        app:      app,
        stopChan: make(chan bool),
    }
}

func (cm *ClipboardMonitor) Start() {
    cm.ticker = time.NewTicker(1 * time.Second)
    
    go func() {
        for {
            select {
            case <-cm.ticker.C:
                cm.checkClipboard()
            case <-cm.stopChan:
                return
            }
        }
    }()
}

func (cm *ClipboardMonitor) Stop() {
    if cm.ticker != nil {
        cm.ticker.Stop()
    }
    cm.stopChan <- true
}

func (cm *ClipboardMonitor) checkClipboard() {
    text, ok := cm.app.Clipboard.Text()
    if !ok {
        return
    }
    
    if text != cm.lastText {
        cm.lastText = text
        cm.app.Event.Emit("clipboard-changed", text)
    }
}
```

### 复制并记录历史

```go
type ClipboardHistory struct {
    app     *application.App
    history []string
    maxSize int
}

func NewClipboardHistory(app *application.App) *ClipboardHistory {
    return &ClipboardHistory{
        app:     app,
        history: make([]string, 0),
        maxSize: 10,
    }
}

func (ch *ClipboardHistory) Copy(text string) bool {
    if !ch.app.Clipboard.SetText(text) {
        return false
    }
    
    // Add to history
    ch.history = append([]string{text}, ch.history...)
    
    // Limit size
    if len(ch.history) > ch.maxSize {
        ch.history = ch.history[:ch.maxSize]
    }
    
    return true
}

func (ch *ClipboardHistory) GetHistory() []string {
    return ch.history
}

func (ch *ClipboardHistory) RestoreFromHistory(index int) bool {
    if index < 0 || index >= len(ch.history) {
        return false
    }
    
    return ch.app.Clipboard.SetText(ch.history[index])
}
```

## 前端集成

### 使用浏览器剪贴板 API

对于简单文本，可以使用浏览器的剪贴板 API：

```javascript
// Copy
async function copyText(text) {
    try {
        await navigator.clipboard.writeText(text)
        console.log("Copied!")
    } catch (error) {
        console.error("Copy failed:", error)
    }
}

// Paste
async function pasteText() {
    try {
        const text = await navigator.clipboard.readText()
        return text
    } catch (error) {
        console.error("Paste failed:", error)
        return ""
    }
}
```

<strong>注意：</strong>浏览器剪贴板 API 需要使用 HTTPS 或 localhost，并获得用户权限。

### 使用 Wails 剪贴板

如需访问系统范围的剪贴板：

```javascript
import { CopyToClipboard, PasteFromClipboard } from './bindings/changeme/clipboardservice'

// Copy
async function copy(text) {
    const success = await CopyToClipboard(text)
    if (success) {
        console.log("Copied!")
    }
}

// Paste
async function paste() {
    const text = await PasteFromClipboard()
    return text
}
```

## 最佳实践

### ✅ 应该做

- **检查返回值**——妥善处理失败情况
- **提供反馈**——让用户知道复制已成功
- **验证粘贴的文本**——检查格式和大小
- **使用适当的方法**——根据情况选择浏览器 API 或 Wails API
- **处理剪贴板为空的情况**——使用前先检查
- **去除首尾空白字符**——清理粘贴的文本

### ❌ 不应该做

- **不要忽略失败**——始终检查操作是否成功
- **不要复制敏感数据**——剪贴板是共享的
- **不要假定格式**——验证粘贴的数据
- **不要过于频繁地轮询**——监视剪贴板时尤其如此
- **不要复制大量数据**——改用文件
- **不要忽视安全性**——净化粘贴的内容

## 平台差异

### macOS

- 使用 NSPasteboard
- 支持富文本（未来）
- 系统范围的剪贴板
- 剪贴板历史记录（系统功能）

### Windows

- 使用 Windows Clipboard API
- 支持多种格式（未来）
- 系统范围的剪贴板
- 剪贴板历史记录（Windows 10+）

### Linux

- 使用 X11/Wayland 剪贴板
- 主选区和剪贴板选区
- 因桌面环境而异
- 可能需要剪贴板管理器

## 限制

### 当前版本

- **仅支持文本**——尚不支持图像
- **不支持格式检测**——仅支持纯文本
- **无剪贴板事件**——必须轮询更改
- **无剪贴板历史记录**——需要自行实现

### 未来功能

- 图像支持
- 富文本支持
- 多格式支持
- 剪贴板更改事件
- 剪贴板历史记录 API

## 后续步骤

@cards{cols="2"}
🚀 绑定
从 JavaScript 调用 Go 函数。

[了解更多 →](/features/bindings/methods/)

---
★ 事件
使用事件发送剪贴板通知。

[了解更多 →](/features/events/system/)

---
ℹ 对话框
显示复制/粘贴操作的反馈。

[了解更多 →](/features/dialogs/message/)

---
◆ 服务
组织剪贴板代码。

[了解更多 →](/features/bindings/services/)

@end

---

<strong>有问题？</strong>请在 [Discord](https://discord.gg/JDdSxwjhGf) 中提问，或查看[剪贴板示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
