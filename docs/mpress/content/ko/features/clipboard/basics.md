---
title: "클립보드 작업"
description: "시스템 클립보드로 텍스트 복사 및 붙여넣기"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## 클립보드 작업

Wails는 모든 플랫폼에서 작동하는 <strong>통합 클립보드 API</strong>를 제공합니다. 간단하고 일관된 메서드를 사용하여 Windows, macOS 및 Linux에서 텍스트를 복사하고 붙여넣을 수 있습니다.

## 빠른 시작

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**이것으로 끝입니다!** 크로스 플랫폼 클립보드에 접근할 수 있습니다.

## 텍스트 복사

### 기본 복사

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**반환값:** `bool` - 성공하면 `true`, 그렇지 않으면 `false`

### 서비스에서 복사

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**JavaScript에서 호출:**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### 피드백을 제공하는 복사

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

## 텍스트 붙여넣기

### 기본 붙여넣기

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**반환값:** `(string, bool)` - 텍스트와 성공 여부 플래그

### 서비스에서 붙여넣기

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**JavaScript에서 호출:**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### 유효성 검사를 포함한 붙여넣기

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

## 전체 예제

### 복사 버튼

**Go:**

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

**JavaScript:**

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

### 붙여넣기 및 처리

**Go:**

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

**JavaScript:**

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

### 여러 형식으로 복사

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

### 클립보드 모니터링

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

### 기록을 포함한 복사

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

## 프런트엔드 통합

### 브라우저 클립보드 API 사용

간단한 텍스트에는 브라우저의 클립보드 API를 사용할 수 있습니다:

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

**참고:** 브라우저 클립보드 API를 사용하려면 HTTPS 또는 localhost 환경과 사용자 권한이 필요합니다.

### Wails 클립보드 사용

시스템 전역 클립보드에 접근하려면 다음을 사용합니다:

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

## 권장 사례

### ✅ 권장 사항

- **반환값 확인** - 실패를 적절하게 처리하세요
- **피드백 제공** - 복사에 성공했음을 사용자에게 알리세요
- **붙여넣은 텍스트 검증** - 형식과 크기를 확인하세요
- **적절한 메서드 사용** - 브라우저 API와 Wails API 중에서 선택하세요
- **빈 클립보드 처리** - 사용하기 전에 확인하세요
- **앞뒤 공백 제거** - 붙여넣은 텍스트를 정리하세요

### ❌ 금지 사항

- **실패를 무시하지 마세요** - 항상 성공 여부를 확인하세요
- **민감한 데이터를 복사하지 마세요** - 클립보드는 공유됩니다
- **형식을 단정하지 마세요** - 붙여넣은 데이터의 유효성을 검사하세요
- **너무 자주 폴링하지 마세요** - 클립보드를 모니터링하는 경우에 해당합니다
- **대용량 데이터를 복사하지 마세요** - 대신 파일을 사용하세요
- **보안을 잊지 마세요** - 붙여넣은 콘텐츠를 정제하세요

## 플랫폼별 차이점

### macOS

- NSPasteboard 사용
- 서식 있는 텍스트 지원(향후)
- 시스템 전역 클립보드
- 클립보드 기록(시스템 기능)

### Windows

- Windows Clipboard API 사용
- 여러 형식 지원(향후)
- 시스템 전역 클립보드
- 클립보드 기록(Windows 10 이상)

### Linux

- X11/Wayland 클립보드 사용
- 기본 선택 영역 및 클립보드 선택 영역
- 데스크톱 환경에 따라 다름
- 클립보드 관리자가 필요할 수 있음

## 제한 사항

### 현재 버전

- **텍스트만 지원** - 이미지는 아직 지원되지 않습니다
- **형식 감지 미지원** - 일반 텍스트만 지원합니다
- **클립보드 이벤트 없음** - 변경 사항을 폴링해야 합니다
- **클립보드 기록 없음** - 직접 구현해야 합니다

### 향후 기능

- 이미지 지원
- 서식 있는 텍스트 지원
- 여러 형식 지원
- 클립보드 변경 이벤트
- 클립보드 기록 API

## 다음 단계

@cards{cols="2"}
🚀 바인딩
JavaScript에서 Go 함수를 호출합니다.

[자세히 알아보기 →](/features/bindings/methods/)

---
★ 이벤트
클립보드 알림에 이벤트를 사용합니다.

[자세히 알아보기 →](/features/events/system/)

---
ℹ 대화 상자
복사 및 붙여넣기 결과를 표시합니다.

[자세히 알아보기 →](/features/dialogs/message/)

---
◆ 서비스
클립보드 코드를 체계적으로 구성합니다.

[자세히 알아보기 →](/features/bindings/services/)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [클립보드 예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
