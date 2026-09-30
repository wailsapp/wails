---
title: "筆記"
description: "建立具備檔案操作與原生對話方塊的筆記應用程式"
slug: "tutorials/03-notes-vanilla"
sourcePath: "tutorials/03-notes-vanilla.md"
---

在本教學中，您將使用 Wails v3建立一個示範檔案操作、原生對話方塊及現代桌面應用程式模式的筆記應用程式。

![筆記應用程式螢幕截圖](/assets/notes-app.png)

## 您將建立的內容

- 具備建立、編輯及刪除功能的完整筆記應用程式
- 用於匯入及匯出筆記的原生檔案儲存／開啟對話方塊
- 輸入時使用防抖自動儲存，以減少不必要的更新
- 仿照 Apple 備忘錄的專業雙欄版面配置（側邊欄 + 編輯器）

## 您將學到的內容

- 在 Wails 中使用原生檔案對話方塊（`SaveFileDialog`、`OpenFileDialog`、`InfoDialog`）
- 使用 JSON 檔案持久保存資料
- 實作防抖自動儲存模式
- 使用現代 CSS 建立專業桌面使用者介面
- 將 Go 結構正確序列化為 JSON

## 專案設定

@steps
### 建立新的 Wails 專案
```bash
wails3 init -n notes-app -t vanilla
cd notes-app
```

### 建立 NotesService
在專案根目錄中建立新檔案`notesservice.go`：

```go
package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NotesService struct {
	notes []Note
}

func NewNotesService() *NotesService {
	return &NotesService{
		notes: make([]Note, 0),
	}
}

// GetAll returns all notes
func (n *NotesService) GetAll() []Note {
	return n.notes
}

// Create creates a new note
func (n *NotesService) Create(title, content string) Note {
	note := Note{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	n.notes = append(n.notes, note)
	return note
}

// Update updates an existing note
func (n *NotesService) Update(id, title, content string) error {
	for i := range n.notes {
		if n.notes[i].ID == id {
			n.notes[i].Title = title
			n.notes[i].Content = content
			n.notes[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("note not found")
}

// Delete deletes a note
func (n *NotesService) Delete(id string) error {
	for i := range n.notes {
		if n.notes[i].ID == id {
			n.notes = append(n.notes[:i], n.notes[i+1:]...)
			return nil
		}
	}
	return errors.New("note not found")
}

// SaveToFile saves notes to a file
func (n *NotesService) SaveToFile() error {
	path, err := application.Get().Dialog.SaveFile().
		SetFilename("notes.json").
		AddFilter("JSON Files", "*.json").
		PromptForSingleSelection()

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(n.notes, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	application.Get().Dialog.Info().
		SetTitle("Success").
		SetMessage("Notes saved successfully!").
		Show()

	return nil
}

// LoadFromFile loads notes from a file
func (n *NotesService) LoadFromFile() error {
	path, err := application.Get().Dialog.OpenFile().
		AddFilter("JSON Files", "*.json").
		PromptForSingleSelection()

	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return err
	}

	n.notes = notes

	application.Get().Dialog.Info().
		SetTitle("Success").
		SetMessage("Notes loaded successfully!").
		Show()

	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
```

**這裡的運作方式：**

- **Note 結構**：使用 JSON 標籤（小寫）定義資料結構，以便正確序列化
- **CRUD 操作**：使用 GetAll、Create、Update 及 Delete 管理記憶體中的筆記
- **檔案對話方塊**：使用`application.Get().Dialog.SaveFile()`及`application.Get().Dialog.OpenFile()`存取原生對話方塊
- **資訊對話方塊**：使用`application.Get().Dialog.Info()`顯示成功訊息
- **ID 產生**：簡易的時間戳記式 ID 產生器

### 更新 main.go
取代`main.go`的內容：

```go
package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Notes App",
		Description: "A simple notes application",
		Services: []application.Service{
			application.NewService(NewNotesService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Notes App",
		Width:            1000,
		Height:           700,
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

**這裡的運作方式：**

- 向應用程式註冊`NotesService`
- 建立尺寸為（1000x700）的視窗，仿照 Apple 備忘錄
- 設定正確的 macOS 行為，在最後一個視窗關閉時結束應用程式

### 建立 HTML 結構
取代`frontend/index.html`：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notes App</title>
    <link rel="stylesheet" href="./style.css">
</head>
<body>
    <div class="app">
        <!-- Sidebar -->
        <div class="sidebar">
            <div class="sidebar-header">
                <h1>Notes</h1>
                <button id="new-note-btn" class="btn-primary">+ New Note</button>
            </div>
            <div id="notes-list" class="notes-list"></div>
            <div class="sidebar-footer">
                <button id="save-btn" class="btn-secondary">Save</button>
                <button id="load-btn" class="btn-secondary">Load</button>
            </div>
        </div>

        <!-- Editor -->
        <div class="editor">
            <div id="empty-state" class="empty-state">
                <h2>No note selected</h2>
                <p>Select a note from the list or create a new one</p>
            </div>
            <div id="note-editor" class="note-editor" style="display: none;">
                <input type="text" id="note-title" placeholder="Note title" class="title-input">
                <textarea id="note-content" placeholder="Start typing..." class="content-input"></textarea>
                <div class="editor-footer">
                    <button id="delete-btn" class="btn-danger">Delete</button>
                    <span id="last-updated" class="last-updated"></span>
                </div>
            </div>
        </div>
    </div>

    <script src="/wails/runtime.js"></script>
    <script type="module" src="./src/main.js"></script>
</body>
</html>
```

**這裡的運作方式：**

- **雙欄版面配置**：側邊欄顯示筆記清單，主要區域則用作編輯器
- **空白狀態**：未選取任何筆記時顯示
- **Wails 執行階段**：必須在模組指令碼之前載入

### 新增 CSS 樣式
取代`frontend/public/style.css`：

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    height: 100vh;
    overflow: hidden;
}

.app {
    display: flex;
    height: 100vh;
}

/* Sidebar */
.sidebar {
    width: 300px;
    background: #f5f5f5;
    border-right: 1px solid #e0e0e0;
    display: flex;
    flex-direction: column;
}

.sidebar-header {
    padding: 20px;
    border-bottom: 1px solid #e0e0e0;
}

.sidebar-header h1 {
    font-size: 24px;
    margin-bottom: 16px;
}

.notes-list {
    flex: 1;
    overflow-y: auto;
}

.note-item {
    padding: 16px 20px;
    border-bottom: 1px solid #e0e0e0;
    cursor: pointer;
    transition: background 0.2s;
}

.note-item:hover {
    background: #e8e8e8;
}

.note-item.active {
    background: #007aff;
    color: white;
}

.note-item h3 {
    font-size: 16px;
    margin-bottom: 4px;
}

.note-item p {
    font-size: 14px;
    opacity: 0.7;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.sidebar-footer {
    padding: 16px 20px;
    border-top: 1px solid #e0e0e0;
    display: flex;
    gap: 8px;
}

/* Editor */
.editor {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #999;
}

.note-editor {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 20px;
}

.title-input {
    font-size: 32px;
    font-weight: bold;
    border: none;
    outline: none;
    margin-bottom: 16px;
    padding: 8px 0;
}

.content-input {
    flex: 1;
    font-size: 16px;
    border: none;
    outline: none;
    resize: none;
    font-family: inherit;
    line-height: 1.6;
}

.editor-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 16px;
    border-top: 1px solid #e0e0e0;
}

.last-updated {
    font-size: 14px;
    color: #999;
}

/* Buttons */
.btn-primary {
    background: #007aff;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    width: 100%;
}

.btn-primary:hover {
    background: #0056b3;
}

.btn-secondary {
    background: white;
    color: #333;
    border: 1px solid #e0e0e0;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    flex: 1;
}

.btn-secondary:hover {
    background: #f5f5f5;
}

.btn-danger {
    background: #ff3b30;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
}

.btn-danger:hover {
    background: #cc0000;
}
```

**這裡的運作方式：**

- 採用簡潔字型排印與配色的 Apple 風格設計
- 使用 Flexbox 為側邊欄與編輯器建立響應式版面配置
- 以藍色背景醒目標示目前的筆記
- 平順的游標懸停轉場效果

### 實作 JavaScript 邏輯
取代`frontend/src/main.js`：

```javascript
import { NotesService } from '../bindings/changeme'

let notes = []
let currentNote = null

// Load notes on startup
async function loadNotes() {
    notes = await NotesService.GetAll()
    renderNotesList()
}

// Render notes list
function renderNotesList() {
    const notesList = document.getElementById('notes-list')

    if (notes.length === 0) {
        notesList.innerHTML = '<div style="padding: 20px; text-align: center; color: #999;">No notes yet</div>'
        return
    }

    notesList.innerHTML = notes.map(note => `
        <div class="note-item ${currentNote?.id === note.id ? 'active' : ''}" data-id="${note.id}">
            <h3>${note.title || 'Untitled'}</h3>
            <p>${note.content || 'No content'}</p>
        </div>
    `).join('')

    // Add click handlers
    document.querySelectorAll('.note-item').forEach(item => {
        item.addEventListener('click', () => {
            const id = item.dataset.id
            selectNote(id)
        })
    })
}

// Select a note
function selectNote(id) {
    currentNote = notes.find(n => n.id === id)
    if (currentNote) {
        document.getElementById('empty-state').style.display = 'none'
        document.getElementById('note-editor').style.display = 'flex'
        document.getElementById('note-title').value = currentNote.title
        document.getElementById('note-content').value = currentNote.content
        document.getElementById('last-updated').textContent =
            `Last updated: ${new Date(currentNote.updatedAt).toLocaleString()}`
        renderNotesList()
    }
}

// Create new note
document.getElementById('new-note-btn').addEventListener('click', async () => {
    const note = await NotesService.Create('Untitled', '')
    notes.push(note)
    selectNote(note.id)
    // Focus the title input and select all text so user can immediately type
    const titleInput = document.getElementById('note-title')
    titleInput.focus()
    titleInput.select()
})

// Update note on input
let updateTimeout
function scheduleUpdate() {
    clearTimeout(updateTimeout)
    updateTimeout = setTimeout(async () => {
        if (currentNote) {
            const title = document.getElementById('note-title').value
            const content = document.getElementById('note-content').value

            await NotesService.Update(currentNote.id, title, content)

            // Update local copy
            const note = notes.find(n => n.id === currentNote.id)
            if (note) {
                note.title = title
                note.content = content
                note.updatedAt = new Date().toISOString()
            }

            renderNotesList()
            document.getElementById('last-updated').textContent =
                `Last updated: ${new Date().toLocaleString()}`
        }
    }, 500)
}

document.getElementById('note-title').addEventListener('input', scheduleUpdate)
document.getElementById('note-content').addEventListener('input', scheduleUpdate)

// Delete note
document.getElementById('delete-btn').addEventListener('click', async () => {
    if (!currentNote) return

    try {
        await NotesService.Delete(currentNote.id)
        notes = notes.filter(n => n.id !== currentNote.id)
        currentNote = null
        document.getElementById('empty-state').style.display = 'flex'
        document.getElementById('note-editor').style.display = 'none'
        renderNotesList()
    } catch (error) {
        console.error('Delete failed:', error)
    }
})

// Save to file
document.getElementById('save-btn').addEventListener('click', async () => {
    try {
        await NotesService.SaveToFile()
    } catch (error) {
        if (error) console.error('Save failed:', error)
    }
})

// Load from file
document.getElementById('load-btn').addEventListener('click', async () => {
    try {
        await NotesService.LoadFromFile()
        notes = await NotesService.GetAll()
        currentNote = null
        document.getElementById('empty-state').style.display = 'flex'
        document.getElementById('note-editor').style.display = 'none'
        renderNotesList()
    } catch (error) {
        if (error) console.error('Load failed:', error)
    }
})

// Initialize
loadNotes()
```

**這裡的運作方式：**

- **自動儲存**：500毫秒的防抖可避免輸入時過度呼叫後端
- **屬性存取**：使用與 Go JSON 標籤相符的小寫屬性名稱（`.id`、`.title`）
- **焦點管理**：建立新筆記時，自動將焦點移至標題並選取標題
- **不經確認即刪除**：瀏覽器的`confirm()`無法在 Wails WebView 中運作
- **檔案操作**：原生對話方塊負責儲存／載入，並提供適當的錯誤處理

### 執行應用程式
```bash
wails3 dev
```

應用程式啟動後，您可以：

- 按一下「+ New Note」建立筆記
- 編輯標題與內容（500毫秒後自動儲存）
- 按一下側邊欄中的筆記以切換筆記
- 按一下「Delete」移除目前的筆記
- 按一下「Save」將筆記匯出為 JSON
- 按一下「Load」匯入先前儲存的筆記

@end

## 重要概念

### 檔案與訊息對話方塊

Wails v3中<strong>沒有</strong>套件層級的對話方塊建構函式；每個對話方塊都透過`app.Dialog`管理器建立。在服務中，請透過`application.Get()`取得應用程式（它會傳回執行中的`*application.App`）：

```go
// Correct — manager-based dialogs.
app := application.Get()
path, err := app.Dialog.SaveFile().
    SetFilename("notes.json").
    AddFilter("JSON Files", "*.json").
    PromptForSingleSelection()
```

對應的方法為用於訊息對話方塊的`app.Dialog.Info() / Question() / Warning() / Error()`，以及用於檔案對話方塊的`app.Dialog.OpenFile() / SaveFile()`（另有`*WithOptions`變體）。

### JSON 標籤對應

Go 結構的 JSON 標籤必須為小寫，才能與 JavaScript 屬性存取相符：

```go
type Note struct {
    ID string `json:"id"` // Must be lowercase
}
```

```javascript
// JavaScript accesses with lowercase
const noteId = note.id // Correct
const noteId = note.ID // Would be undefined
```

### 防抖自動儲存

500毫秒的防抖可減少不必要的後端呼叫：

```javascript
let updateTimeout
function scheduleUpdate() {
    clearTimeout(updateTimeout) // Cancel previous timer
    updateTimeout = setTimeout(async () => {
        // Only saves if user stops typing for 500ms
        await NotesService.Update(currentNote.id, title, content)
    }, 500)
}
```

## 後續步驟

- 新增分類或標籤來整理筆記
- 實作搜尋與篩選功能
- 使用所見即所得編輯器新增富文字編輯功能
- 將筆記同步至雲端儲存空間
- 為常用操作新增鍵盤快速鍵
