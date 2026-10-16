---
title: "クリップボード操作"
description: "システムクリップボードを使用してテキストをコピー＆ペーストする"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## クリップボード操作

Wails は、すべてのプラットフォームで動作する<strong>統一クリップボード API</strong>を提供します。Windows、macOS、Linux で、シンプルかつ一貫したメソッドを使用してテキストをコピー＆ペーストできます。

## クイックスタート

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

<strong>これだけです！</strong>クロスプラットフォームでクリップボードにアクセスできます。

## テキストのコピー

### 基本的なコピー

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**戻り値：**`bool` - 成功した場合は`true`、それ以外の場合は`false`

### サービスからのコピー

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**JavaScript から呼び出す：**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### フィードバック付きのコピー

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

## テキストの貼り付け

### 基本的な貼り付け

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**戻り値：**`(string, bool)` - テキストと成功フラグ

### サービスからの貼り付け

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**JavaScript から呼び出す：**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### 検証付きの貼り付け

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

## 完全な例

### コピーボタン

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

### 貼り付けと処理

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

### 複数形式でのコピー

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

### クリップボードの監視

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

### 履歴付きのコピー

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

## フロントエンドとの統合

### ブラウザーの Clipboard API の使用

単純なテキストには、ブラウザーの Clipboard API を使用できます：

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

<strong>注：</strong>ブラウザーの Clipboard API には、HTTPS または localhost と、ユーザーの許可が必要です。

### Wails クリップボードの使用

システム全体のクリップボードにアクセスするには：

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

## ベストプラクティス

### ✅ 推奨事項

- **戻り値を確認する** - 失敗を適切に処理する
- **フィードバックを提示する** - コピーが成功したことをユーザーに知らせる
- **貼り付けたテキストを検証する** - 形式とサイズを確認する
- **適切な方法を使用する** - ブラウザー API と Wails API を使い分ける
- **クリップボードが空の場合に対処する** - 使用前に確認する
- **前後の空白を除去する** - 貼り付けたテキストを整える

### ❌ 禁止事項

- **失敗を無視しない** - 必ず成功したか確認する
- **機密データをコピーしない** - クリップボードは共有される
- **形式を決めつけない** - 貼り付けたデータを検証する
- **頻繁にポーリングしすぎない** - クリップボードを監視する場合
- **大容量のデータをコピーしない** - 代わりにファイルを使用する
- **セキュリティを忘れない** - 貼り付けた内容をサニタイズする

## プラットフォームごとの違い

### macOS

- NSPasteboard を使用
- リッチテキストをサポート（将来対応）
- システム全体のクリップボード
- クリップボード履歴（システム機能）

### Windows

- Windows Clipboard API を使用
- 複数形式をサポート（将来対応）
- システム全体のクリップボード
- クリップボード履歴（Windows 10以降）

### Linux

- X11/Wayland のクリップボードを使用
- プライマリ選択とクリップボード選択
- デスクトップ環境によって異なる
- クリップボードマネージャーが必要な場合がある

## 制限事項

### 現在のバージョン

- **テキストのみ** - 画像はまだサポートされていない
- **形式検出なし** - プレーンテキストのみ
- **クリップボードイベントなし** - 変更をポーリングする必要があります
- **クリップボード履歴なし** - 独自に実装する必要があります

### 今後の機能

- 画像のサポート
- リッチテキストのサポート
- 複数形式
- クリップボード変更イベント
- クリップボード履歴 API

## 次のステップ

@cards{cols="2"}
🚀 バインディング
JavaScript から Go 関数を呼び出します。

[詳細を見る →](/features/bindings/methods/)

---
★ イベント
クリップボードの通知にはイベントを使用します。

[詳細を見る →](/features/events/system/)

---
ℹ ダイアログ
コピー／貼り付けの結果を表示します。

[詳細を見る →](/features/dialogs/message/)

---
◆ サービス
クリップボード関連のコードを整理します。

[詳細を見る →](/features/bindings/services/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[クリップボードのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
