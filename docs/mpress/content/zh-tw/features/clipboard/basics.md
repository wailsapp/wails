---
title: "剪貼簿操作"
description: "透過系統剪貼簿複製及貼上文字"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## 剪貼簿操作

Wails 提供可在所有平台上運作的<strong>統一剪貼簿 API</strong>。在 Windows、macOS 及 Linux 上使用簡單且一致的方法複製及貼上文字。

## 快速入門

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

<strong>就這麼簡單！</strong>跨平台存取剪貼簿。

## 複製文字

### 基本複製

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**傳回值：**`bool`—成功時為`true`，否則為`false`

### 從服務複製

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**從 JavaScript 呼叫：**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### 複製並提供回饋

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

## 貼上文字

### 基本貼上

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**傳回值：**`(string, bool)`—文字及成功旗標

### 從服務貼上

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**從 JavaScript 呼叫：**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### 貼上並驗證

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

## 完整範例

### 複製按鈕

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

### 貼上並處理

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

### 複製多種格式

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

### 剪貼簿監控

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

### 複製並保留歷程記錄

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

## 前端整合

### 使用瀏覽器剪貼簿 API

若只需處理純文字，可以使用瀏覽器的剪貼簿 API：

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

<strong>注意：</strong>瀏覽器剪貼簿 API 需要 HTTPS 或 localhost，並須取得使用者權限。

### 使用 Wails 剪貼簿

若要存取全系統剪貼簿：

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

## 最佳實務

### ✅ 建議做法

- **檢查傳回值**—妥善處理失敗情況
- **提供回饋**—讓使用者知道複製成功
- **驗證貼上的文字**—檢查格式及大小
- **使用適當的方法**—依情況選用瀏覽器 API 或 Wails API
- **處理剪貼簿為空的情況**—使用前先檢查
- **移除頭尾空白**—清理貼上的文字

### ❌ 請勿

- **請勿忽略失敗情況**—一律檢查操作是否成功
- **請勿複製敏感資料**—剪貼簿由系統共用
- **請勿預設資料格式**—驗證貼上的資料
- **請勿過度頻繁地輪詢**—監控剪貼簿時尤其如此
- **請勿複製大量資料**—改用檔案
- **請勿忽略安全性**—清理貼上的內容

## 平台差異

### macOS

- 使用 NSPasteboard
- 支援格式化文字（未來功能）
- 全系統剪貼簿
- 剪貼簿歷程記錄（系統功能）

### Windows

- 使用 Windows Clipboard API
- 將支援多種格式（未來功能）
- 全系統剪貼簿
- 剪貼簿歷程記錄（Windows 10+）

### Linux

- 使用 X11/Wayland 剪貼簿
- 主要選取區及剪貼簿選取區
- 依桌面環境而異
- 可能需要剪貼簿管理程式

## 限制

### 目前版本

- **僅支援文字**—尚未支援圖片
- **不支援格式偵測**—僅支援純文字
- **無剪貼簿事件**－必須輪詢變更
- **無剪貼簿歷程記錄**－必須自行實作

### 未來功能

- 圖片支援
- 格式化文字支援
- 多種格式
- 剪貼簿變更事件
- 剪貼簿歷程記錄 API

## 後續步驟

@cards{cols="2"}
🚀 繫結
從 JavaScript 呼叫 Go 函式。

[深入瞭解 →](/features/bindings/methods/)

---
★ 事件
使用事件發出剪貼簿通知。

[深入瞭解 →](/features/events/system/)

---
ℹ 對話方塊
顯示複製／貼上的操作結果。

[深入瞭解 →](/features/dialogs/message/)

---
◆ 服務
整理剪貼簿程式碼。

[深入瞭解 →](/features/bindings/services/)

@end

---

<strong>有問題嗎？</strong>請在 [Discord](https://discord.gg/JDdSxwjhGf) 中提問，或查看[剪貼簿範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
