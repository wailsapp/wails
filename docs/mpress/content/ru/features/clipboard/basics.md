---
title: "Операции с буфером обмена"
description: "Копирование и вставка текста через системный буфер обмена"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## Операции с буфером обмена

Wails предоставляет **унифицированный API буфера обмена**, который работает на всех платформах. Копируйте и вставляйте текст с помощью простых и единообразных методов в Windows, macOS и Linux.

## Быстрый старт

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**Вот и всё!** Кроссплатформенный доступ к буферу обмена.

## Копирование текста

### Базовое копирование

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**Возвращает:** `bool` — `true` при успешном выполнении, иначе `false`

### Копирование из сервиса

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**Вызов из JavaScript:**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### Копирование с уведомлением

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

## Вставка текста

### Базовая вставка

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**Возвращает:** `(string, bool)` — текст и флаг успешного выполнения

### Вставка из сервиса

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**Вызов из JavaScript:**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### Вставка с проверкой

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

## Полные примеры

### Кнопка копирования

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

### Вставка и обработка

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

### Копирование в нескольких форматах

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

### Мониторинг буфера обмена

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

### Копирование с сохранением истории

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

## Интеграция с фронтендом

### Использование API буфера обмена браузера

Для простого текста можно использовать API буфера обмена браузера:

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

**Примечание:** для API буфера обмена браузера требуется HTTPS или localhost, а также разрешение пользователя.

### Использование буфера обмена Wails

Для общесистемного доступа к буферу обмена:

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

## Рекомендации

### ✅ Рекомендуется

- **Проверяйте возвращаемые значения** — корректно обрабатывайте ошибки
- **Предоставляйте обратную связь** — сообщайте пользователям об успешном копировании
- **Проверяйте вставленный текст** — проверяйте его формат и размер
- **Используйте подходящий метод** — API браузера или API Wails
- **Обрабатывайте пустой буфер обмена** — проверяйте его перед использованием
- **Удаляйте пробельные символы по краям** — очищайте вставленный текст

### ❌ Не рекомендуется

- **Не игнорируйте ошибки** — всегда проверяйте успешность операции
- **Не копируйте конфиденциальные данные** — буфер обмена является общим
- **Не полагайтесь на предполагаемый формат** — проверяйте вставленные данные
- **Не выполняйте опрос слишком часто** — при мониторинге буфера обмена
- **Не копируйте большие объёмы данных** — вместо этого используйте файлы
- **Не забывайте о безопасности** — очищайте вставляемое содержимое от потенциально опасных данных

## Различия между платформами

### macOS

- Использует NSPasteboard
- Поддержка форматированного текста (в будущем)
- Общесистемный буфер обмена
- История буфера обмена (системная функция)

### Windows

- Использует Windows Clipboard API
- Поддержка нескольких форматов (в будущем)
- Общесистемный буфер обмена
- История буфера обмена (Windows 10+)

### Linux

- Использует буфер обмена X11/Wayland
- Основное выделение и выделение буфера обмена
- Зависит от окружения рабочего стола
- Может потребоваться менеджер буфера обмена

## Ограничения

### Текущая версия

- **Только текст** — изображения пока не поддерживаются
- **Без определения формата** — поддерживается только простой текст
- **Нет событий буфера обмена** — изменения необходимо отслеживать с помощью опроса
- **Нет истории буфера обмена** — реализуйте её самостоятельно

### Будущие возможности

- Поддержка изображений
- Поддержка форматированного текста
- Несколько форматов
- События изменения буфера обмена
- API истории буфера обмена

## Дальнейшие шаги

@cards{cols="2"}
🚀 Привязки
Вызывайте функции Go из JavaScript.

[Подробнее →](/features/bindings/methods/)

---
★ События
Используйте события для уведомлений о буфере обмена.

[Подробнее →](/features/events/system/)

---
ℹ Диалоговые окна
Показывайте уведомления о копировании и вставке.

[Подробнее →](/features/dialogs/message/)

---
◆ Сервисы
Организуйте код для работы с буфером обмена.

[Подробнее →](/features/bindings/services/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами работы с буфером обмена](https://github.com/wailsapp/wails/tree/master/v3/examples).
