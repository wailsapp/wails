---
title: "Opérations sur le presse-papiers"
description: "Copier et coller du texte avec le presse-papiers système"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## Opérations sur le presse-papiers

Wails fournit une **API de presse-papiers unifiée** qui fonctionne sur toutes les plateformes. Copiez et collez du texte à l’aide de méthodes simples et cohérentes sous Windows, macOS et Linux.

## Démarrage rapide

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**C’est tout !** Vous disposez maintenant d’un accès multiplateforme au presse-papiers.

## Copie de texte

### Copie simple

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**Renvoie :** `bool` — `true` en cas de réussite, `false` sinon

### Copie depuis un service

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**Appelez depuis JavaScript :**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### Copie avec confirmation

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

## Collage de texte

### Collage simple

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**Renvoie :** `(string, bool)` — le texte et un indicateur de réussite

### Collage depuis un service

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**Appelez depuis JavaScript :**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### Collage avec validation

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

## Exemples complets

### Bouton de copie

**Go :**

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

**JavaScript :**

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

### Collage et traitement

**Go :**

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

**JavaScript :**

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

### Copie dans plusieurs formats

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

### Surveillance du presse-papiers

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

### Copie avec historique

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

## Intégration au frontend

### Utilisation de l’API de presse-papiers du navigateur

Pour du texte simple, vous pouvez utiliser l’API de presse-papiers du navigateur :

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

**Remarque :** l’API de presse-papiers du navigateur nécessite HTTPS ou localhost, ainsi que l’autorisation de l’utilisateur.

### Utilisation du presse-papiers de Wails

Pour accéder au presse-papiers à l’échelle du système :

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

## Bonnes pratiques

### ✅ À faire

- **Vérifiez les valeurs de retour** — gérez les échecs de manière appropriée
- **Fournissez un retour** — informez les utilisateurs que la copie a réussi
- **Validez le texte collé** — vérifiez son format et sa taille
- **Utilisez la méthode appropriée** — API du navigateur ou API de Wails
- **Gérez un presse-papiers vide** — effectuez une vérification avant toute utilisation
- **Supprimez les caractères blancs en début et en fin de texte** — nettoyez le texte collé

### ❌ À ne pas faire

- **N’ignorez pas les échecs** — vérifiez toujours la réussite de l’opération
- **Ne copiez pas de données sensibles** — le presse-papiers est partagé
- **Ne présumez pas du format** — validez les données collées
- **N’effectuez pas d’interrogations trop fréquentes** — si vous surveillez le presse-papiers
- **Ne copiez pas de données volumineuses** — utilisez plutôt des fichiers
- **Ne négligez pas la sécurité** — assainissez le contenu collé

## Différences entre les plateformes

### macOS

- Utilise NSPasteboard
- Prise en charge future du texte enrichi
- Presse-papiers à l’échelle du système
- Historique du presse-papiers (fonctionnalité système)

### Windows

- Utilise l’API Presse-papiers de Windows
- Prise en charge future de plusieurs formats
- Presse-papiers à l’échelle du système
- Historique du presse-papiers (Windows 10+)

### Linux

- Utilise le presse-papiers de X11/Wayland
- Sélections primaire et du presse-papiers
- Varie selon l’environnement de bureau
- Peut nécessiter un gestionnaire de presse-papiers

## Limitations

### Version actuelle

- **Texte uniquement** — les images ne sont pas encore prises en charge
- **Aucune détection du format** — texte brut uniquement
- **Aucun événement du presse-papiers** — Vous devez rechercher régulièrement les modifications
- **Aucun historique du presse-papiers** — Implémentez-le vous-même

### Fonctionnalités futures

- Prise en charge des images
- Prise en charge du texte enrichi
- Formats multiples
- Événements de modification du presse-papiers
- API d’historique du presse-papiers

## Étapes suivantes

@cards{cols="2"}
🚀 Liaisons
Appelez des fonctions Go depuis JavaScript.

[En savoir plus →](/features/bindings/methods/)

---
★ Événements
Utilisez des événements pour recevoir les notifications du presse-papiers.

[En savoir plus →](/features/events/system/)

---
ℹ Boîtes de dialogue
Affichez une confirmation après une opération de copie ou de collage.

[En savoir plus →](/features/dialogs/message/)

---
◆ Services
Organisez le code du presse-papiers.

[En savoir plus →](/features/bindings/services/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples relatifs au presse-papiers](https://github.com/wailsapp/wails/tree/master/v3/examples).
