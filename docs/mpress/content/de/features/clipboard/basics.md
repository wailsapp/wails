---
title: "Zwischenablageoperationen"
description: "Text über die Systemzwischenablage kopieren und einfügen"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## Zwischenablageoperationen

Wails bietet eine **einheitliche Zwischenablage-API**, die auf allen Plattformen funktioniert. Kopieren und fügen Sie Text unter Windows, macOS und Linux mit einfachen, einheitlichen Methoden ein.

## Schnelleinstieg

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**Das ist alles!** Plattformübergreifender Zugriff auf die Zwischenablage.

## Text kopieren

### Einfaches Kopieren

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**Rückgabewert:** `bool` – bei Erfolg `true`, andernfalls `false`

### Aus einem Service kopieren

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**Aufruf aus JavaScript:**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### Kopieren mit Rückmeldung

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

## Text einfügen

### Einfaches Einfügen

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**Rückgabewert:** `(string, bool)` – Text und Erfolgskennzeichen

### Aus einem Service einfügen

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**Aufruf aus JavaScript:**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### Einfügen mit Validierung

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

## Vollständige Beispiele

### Kopierschaltfläche

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

### Einfügen und verarbeiten

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

### Mehrere Formate kopieren

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

### Zwischenablage überwachen

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

### Mit Verlauf kopieren

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

## Frontend-Integration

### Zwischenablage-API des Browsers verwenden

Für einfachen Text können Sie die Zwischenablage-API des Browsers verwenden:

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

**Hinweis:** Die Zwischenablage-API des Browsers erfordert HTTPS oder localhost sowie die Berechtigung des Benutzers.

### Wails-Zwischenablage verwenden

Für systemweiten Zugriff auf die Zwischenablage:

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

## Bewährte Methoden

### ✅ Empfohlen

- **Rückgabewerte prüfen** – Fehler ordnungsgemäß behandeln
- **Rückmeldung geben** – Benutzer über erfolgreiches Kopieren informieren
- **Eingefügten Text validieren** – Format und Größe prüfen
- **Geeignete Methode verwenden** – Browser-API oder Wails-API
- **Leere Zwischenablage berücksichtigen** – Vor der Verwendung prüfen
- **Führenden und nachgestellten Leerraum entfernen** – Eingefügten Text bereinigen

### ❌ Nicht empfohlen

- **Fehler nicht ignorieren** – Erfolg immer prüfen
- **Keine vertraulichen Daten kopieren** – Die Zwischenablage wird gemeinsam genutzt
- **Format nicht voraussetzen** – Eingefügte Daten validieren
- **Nicht zu häufig abfragen** – Beim Überwachen der Zwischenablage
- **Keine großen Datenmengen kopieren** – Stattdessen Dateien verwenden
- **Sicherheit nicht vernachlässigen** – Eingefügte Inhalte bereinigen

## Plattformunterschiede

### macOS

- Verwendet NSPasteboard
- Unterstützt Rich Text (zukünftig)
- Systemweite Zwischenablage
- Zwischenablageverlauf (Systemfunktion)

### Windows

- Verwendet die Windows Clipboard API
- Unterstützt mehrere Formate (zukünftig)
- Systemweite Zwischenablage
- Zwischenablageverlauf (Windows 10+)

### Linux

- Verwendet die X11-/Wayland-Zwischenablage
- Primär- und Zwischenablageauswahl
- Abhängig von der Desktop-Umgebung
- Möglicherweise ist ein Zwischenablage-Manager erforderlich

## Einschränkungen

### Aktuelle Version

- **Nur Text** – Bilder werden noch nicht unterstützt
- **Keine Formaterkennung** – Nur reiner Text
- **Keine Zwischenablageereignisse** – Änderungen müssen regelmäßig abgefragt werden
- **Kein Zwischenablageverlauf** – selbst implementieren

### Zukünftige Funktionen

- Bildunterstützung
- Unterstützung für formatierten Text
- Mehrere Formate
- Ereignisse bei Änderungen der Zwischenablage
- API für den Zwischenablageverlauf

## Nächste Schritte

@cards{cols="2"}
🚀 Bindings
Go-Funktionen aus JavaScript aufrufen.

[Mehr erfahren →](/features/bindings/methods/)

---
★ Ereignisse
Ereignisse für Benachrichtigungen zur Zwischenablage verwenden.

[Mehr erfahren →](/features/events/system/)

---
ℹ Dialoge
Rückmeldungen zu Kopier- und Einfügevorgängen anzeigen.

[Mehr erfahren →](/features/dialogs/message/)

---
◆ Dienste
Code für die Zwischenablage organisieren.

[Mehr erfahren →](/features/bindings/services/)

@end

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Zwischenablagebeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
