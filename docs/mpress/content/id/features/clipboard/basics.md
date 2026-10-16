---
title: "Operasi Papan Klip"
description: "Salin dan tempel teks menggunakan papan klip sistem"
slug: "features/clipboard/basics"
sourcePath: "features/clipboard/basics.md"
---

## Operasi Papan Klip

Wails menyediakan **API papan klip terpadu** yang berfungsi di semua platform. Salin dan tempel teks dengan metode yang sederhana dan konsisten di Windows, macOS, dan Linux.

## Mulai Cepat

```go
// Copy text to clipboard
app.Clipboard.SetText("Hello, World!")

// Get text from clipboard
text, ok := app.Clipboard.Text()
if ok {
    fmt.Println("Clipboard:", text)
}
```

**Selesai!** Akses papan klip lintas platform.

## Menyalin Teks

### Penyalinan Dasar

```go
success := app.Clipboard.SetText("Text to copy")
if !success {
    app.Logger.Error("Failed to copy to clipboard")
}
```

**Mengembalikan:** `bool` - `true` jika berhasil, `false` jika tidak

### Menyalin dari Layanan

```go
type ClipboardService struct {
    app *application.App
}

func (c *ClipboardService) CopyToClipboard(text string) bool {
    return c.app.Clipboard.SetText(text)
}
```

**Panggil dari JavaScript:**

```javascript
import { CopyToClipboard } from './bindings/changeme/clipboardservice'

await CopyToClipboard("Text to copy")
```

### Menyalin dengan Umpan Balik

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

## Menempelkan Teks

### Penempelan Dasar

```go
text, ok := app.Clipboard.Text()
if !ok {
    app.Logger.Error("Failed to read clipboard")
    return
}

fmt.Println("Clipboard text:", text)
```

**Mengembalikan:** `(string, bool)` - Teks dan penanda keberhasilan

### Menempelkan dari Layanan

```go
func (c *ClipboardService) PasteFromClipboard() string {
    text, ok := c.app.Clipboard.Text()
    if !ok {
        return ""
    }
    return text
}
```

**Panggil dari JavaScript:**

```javascript
import { PasteFromClipboard } from './bindings/changeme/clipboardservice'

const text = await PasteFromClipboard()
console.log("Pasted:", text)
```

### Menempelkan dengan Validasi

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

## Contoh Lengkap

### Tombol Salin

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

### Tempel dan Proses

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

### Menyalin Beberapa Format

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

### Pemantau Papan Klip

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

### Menyalin dengan Riwayat

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

## Integrasi Frontend

### Menggunakan API Papan Klip Browser

Untuk teks sederhana, Anda dapat menggunakan API papan klip browser:

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

**Catatan:** API papan klip browser memerlukan HTTPS atau localhost serta izin pengguna.

### Menggunakan Papan Klip Wails

Untuk mengakses papan klip di seluruh sistem:

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

## Praktik Terbaik

### ✅ Lakukan

- **Periksa nilai kembalian** - Tangani kegagalan dengan baik
- **Berikan umpan balik** - Beri tahu pengguna bahwa penyalinan berhasil
- **Validasi teks yang ditempelkan** - Periksa format dan ukurannya
- **Gunakan metode yang sesuai** - API browser atau API Wails
- **Tangani papan klip kosong** - Periksa sebelum digunakan
- **Hapus spasi kosong di awal dan akhir** - Bersihkan teks yang ditempelkan

### ❌ Jangan Lakukan

- **Jangan abaikan kegagalan** - Selalu periksa keberhasilan
- **Jangan salin data sensitif** - Papan klip digunakan bersama
- **Jangan berasumsi tentang format** - Validasi data yang ditempelkan
- **Jangan lakukan polling terlalu sering** - Saat memantau papan klip
- **Jangan salin data berukuran besar** - Gunakan file sebagai gantinya
- **Jangan abaikan keamanan** - Sanitasi konten yang ditempelkan

## Perbedaan Antarplatform

### macOS

- Menggunakan NSPasteboard
- Mendukung teks kaya (di masa mendatang)
- Papan klip seluruh sistem
- Riwayat papan klip (fitur sistem)

### Windows

- Menggunakan Windows Clipboard API
- Mendukung beberapa format (di masa mendatang)
- Papan klip seluruh sistem
- Riwayat papan klip (Windows 10+)

### Linux

- Menggunakan papan klip X11/Wayland
- Pilihan primer dan papan klip
- Bervariasi menurut lingkungan desktop
- Mungkin memerlukan pengelola papan klip

## Keterbatasan

### Versi Saat Ini

- **Hanya teks** - Gambar belum didukung
- **Tidak ada deteksi format** - Hanya teks biasa
- **Tidak ada peristiwa papan klip** - Perubahan harus diperiksa secara berkala
- **Tidak ada riwayat papan klip** - Implementasikan sendiri

### Fitur Mendatang

- Dukungan gambar
- Dukungan teks kaya
- Berbagai format
- Peristiwa perubahan papan klip
- API riwayat papan klip

## Langkah Berikutnya

@cards{cols="2"}
🚀 Binding
Panggil fungsi Go dari JavaScript.

[Pelajari Lebih Lanjut →](/features/bindings/methods/)

---
★ Peristiwa
Gunakan peristiwa untuk notifikasi papan klip.

[Pelajari Lebih Lanjut →](/features/events/system/)

---
ℹ dialog
Tampilkan umpan balik penyalinan atau penempelan.

[Pelajari Lebih Lanjut →](/features/dialogs/message/)

---
◆ Layanan
Atur kode papan klip.

[Pelajari Lebih Lanjut →](/features/bindings/services/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh papan klip](https://github.com/wailsapp/wails/tree/master/v3/examples).
