---
title: "Operações da área de transferência"
description: "Copie e cole texto usando a área de transferência do sistema"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## Operações da área de transferência

O Wails fornece uma **API unificada para a área de transferência** que funciona em todas as plataformas. Copie e cole texto com métodos simples e consistentes no Windows, macOS e Linux.

## Início rápido

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**É só isso!** Acesso multiplataforma à área de transferência.

## Como copiar texto

### Cópia básica

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**Retorna:** `bool` — `true` em caso de sucesso; caso contrário, `false`

### Copiar a partir de um serviço

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**Chamada a partir do JavaScript:**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### Copiar com confirmação

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

## Como colar texto

### Colagem básica

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**Retorna:** `(string, bool)` — texto e indicador de sucesso

### Colar a partir de um serviço

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**Chamada a partir do JavaScript:**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### Colar com validação

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

## Exemplos completos

### Botão de cópia

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

### Colar e processar

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

### Copiar em vários formatos

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

### Monitoramento da área de transferência

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

### Copiar com histórico

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

## Integração com o frontend

### Como usar a API da área de transferência do navegador

Para texto simples, você pode usar a API da área de transferência do navegador:

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

**Observação:** a API da área de transferência do navegador requer HTTPS ou localhost, além da permissão do usuário.

### Como usar a área de transferência do Wails

Para acessar a área de transferência em todo o sistema:

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

## Práticas recomendadas

### ✅ Faça

- **Verifique os valores retornados** — trate as falhas de forma adequada
- **Forneça uma confirmação** — informe aos usuários que a cópia foi bem-sucedida
- **Valide o texto colado** — verifique o formato e o tamanho
- **Use o método adequado** — API do navegador ou API do Wails
- **Trate a área de transferência vazia** — verifique antes de usar
- **Remova os espaços em branco no início e no fim** — limpe o texto colado

### ❌ Não faça

- **Não ignore falhas** — sempre verifique se a operação foi bem-sucedida
- **Não copie dados confidenciais** — a área de transferência é compartilhada
- **Não presuma o formato** — valide os dados colados
- **Não faça consultas com frequência excessiva** — ao monitorar a área de transferência
- **Não copie grandes volumes de dados** — use arquivos em vez disso
- **Não se esqueça da segurança** — higienize o conteúdo colado

## Diferenças entre plataformas

### macOS

- Usa NSPasteboard
- Compatibilidade futura com rich text
- Área de transferência disponível em todo o sistema
- Histórico da área de transferência (recurso do sistema)

### Windows

- Usa a API da área de transferência do Windows
- Compatibilidade futura com vários formatos
- Área de transferência disponível em todo o sistema
- Histórico da área de transferência (Windows 10+)

### Linux

- Usa a área de transferência do X11/Wayland
- Seleções primária e da área de transferência
- Varia conforme o ambiente de desktop
- Pode exigir um gerenciador da área de transferência

## Limitações

### Versão atual

- **Somente texto** — imagens ainda não são compatíveis
- **Sem detecção de formato** — somente texto sem formatação
- **Sem eventos da área de transferência** — É necessário verificar periodicamente se há alterações
- **Sem histórico da área de transferência** — Implemente-o por conta própria

### Recursos futuros

- Suporte a imagens
- Suporte a rich text
- Vários formatos
- Eventos de alteração da área de transferência
- API de histórico da área de transferência

## Próximas etapas

@cards{cols="2"}
🚀 Bindings
Chame funções Go pelo JavaScript.

[Saiba mais →](/features/bindings/methods/)

---
★ Eventos
Use eventos para receber notificações da área de transferência.

[Saiba mais →](/features/events/system/)

---
ℹ Caixas de diálogo
Mostre uma confirmação ao copiar ou colar.

[Saiba mais →](/features/dialogs/message/)

---
◆ Serviços
Organize o código da área de transferência.

[Saiba mais →](/features/bindings/services/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de uso da área de transferência](https://github.com/wailsapp/wails/tree/master/v3/examples).
