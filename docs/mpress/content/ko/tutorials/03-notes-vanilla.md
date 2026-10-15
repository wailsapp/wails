---
title: "메모"
description: "파일 작업과 네이티브 대화 상자를 사용하는 메모 애플리케이션 만들기"
slug: "tutorials/03-notes-vanilla"
sourcePath: "tutorials/03-notes-vanilla.md"
---

이 튜토리얼에서는 Wails v3을 사용하여 파일 작업, 네이티브 대화 상자, 최신 데스크톱 앱 패턴을 보여 주는 메모 애플리케이션을 만듭니다.

![메모 앱 스크린샷](/assets/notes-app.png)

## 만들어 볼 항목

- 메모 작성, 편집, 삭제 기능을 갖춘 완전한 메모 앱
- 메모를 가져오고 내보내기 위한 네이티브 파일 저장/열기 대화 상자
- 불필요한 업데이트를 줄이기 위해 디바운스를 적용한 입력 시 자동 저장
- Apple Notes를 본뜬 전문적인 2열 레이아웃(사이드바 + 편집기)

## 학습할 내용

- Wails에서 네이티브 파일 대화 상자 사용하기(`SaveFileDialog`, `OpenFileDialog`, `InfoDialog`)
- 데이터 영속화를 위한 JSON 파일 사용하기
- 디바운스가 적용된 자동 저장 패턴 구현하기
- 최신 CSS로 전문적인 데스크톱 UI 만들기
- Go 구조체를 올바르게 JSON으로 직렬화하기

## 프로젝트 설정

@steps
### 새 Wails 프로젝트 만들기
```bash
wails3 init -n notes-app -t vanilla
cd notes-app
```

### NotesService 만들기
프로젝트 루트에 새 파일 `notesservice.go`을 만드세요:

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

**여기서 수행되는 작업:**

- **Note 구조체**: 올바른 직렬화를 위해 소문자 JSON 태그가 지정된 데이터 구조를 정의합니다.
- **CRUD 작업**: 메모리에서 메모를 관리하기 위한 GetAll, Create, Update, Delete 작업입니다.
- **파일 대화 상자**: 네이티브 대화 상자에 접근하기 위해 `application.Get().Dialog.SaveFile()` 및 `application.Get().Dialog.OpenFile()`을 사용합니다.
- **정보 대화 상자**: `application.Get().Dialog.Info()`을 사용하여 성공 메시지를 표시합니다.
- **ID 생성**: 타임스탬프를 기반으로 하는 간단한 ID 생성기입니다.

### main.go 업데이트하기
`main.go`의 내용을 교체하세요:

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

**여기서 수행되는 작업:**

- 애플리케이션에 `NotesService`을 등록합니다.
- Apple Notes를 본뜬 크기(1000x700)의 창을 만듭니다.
- 마지막 창이 닫히면 종료되도록 올바른 macOS 동작을 설정합니다.

### HTML 구조 만들기
`frontend/index.html`을 교체하세요:

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

**여기서 수행되는 작업:**

- **2열 레이아웃**: 사이드바에는 메모 목록을, 기본 영역에는 편집기를 배치합니다.
- **빈 상태**: 선택된 메모가 없을 때 표시됩니다.
- **Wails 런타임**: 모듈 스크립트보다 먼저 로드해야 합니다.

### CSS 스타일 추가하기
`frontend/public/style.css`을 교체하세요:

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

**여기서 수행되는 작업:**

- 깔끔한 타이포그래피와 색상을 사용한 Apple 스타일 디자인
- 반응형 사이드바와 편집기를 위한 Flexbox 레이아웃
- 활성 메모를 파란색 배경으로 강조 표시
- 부드러운 마우스 오버 전환 효과

### JavaScript 로직 구현하기
`frontend/src/main.js`을 교체하세요:

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

**여기서 수행되는 작업:**

- **자동 저장**: 500ms 디바운스를 적용하여 입력 중 백엔드 호출이 지나치게 많이 발생하지 않도록 합니다.
- **속성 접근**: Go JSON 태그와 일치하는 소문자 속성 이름(`.id`, `.title`)을 사용합니다.
- **포커스 관리**: 새 메모를 만들 때 제목에 자동으로 포커스를 두고 제목을 선택합니다.
- **확인 없는 삭제**: 브라우저의 `confirm()`은 Wails WebView에서 작동하지 않습니다.
- **파일 작업**: 네이티브 대화 상자에서 적절한 오류 처리와 함께 저장 및 불러오기를 수행합니다.

### 애플리케이션 실행하기
```bash
wails3 dev
```

앱이 시작되면 다음 작업을 할 수 있습니다:

- "+ New Note"를 클릭하여 메모를 만드세요.
- 제목과 내용을 편집하세요(500ms 후 자동 저장됩니다).
- 사이드바에서 메모를 클릭하여 메모 사이를 전환하세요.
- "Delete"를 클릭하여 현재 메모를 삭제하세요.
- "Save"를 클릭하여 메모를 JSON으로 내보내세요.
- "Load"를 클릭하여 이전에 저장한 메모를 가져오세요.

@end

## 핵심 개념

### 파일 및 메시지 대화 상자

Wails v3에는 패키지 수준의 대화 상자 생성자가 **없으며**, 모든 대화 상자는 `app.Dialog` 관리자를 통해 생성됩니다. 서비스에서는 `application.Get()`을 통해 앱을 가져오세요. 이 함수는 실행 중인 `*application.App`을 반환합니다:

```go
// Correct — manager-based dialogs.
app := application.Get()
path, err := app.Dialog.SaveFile().
    SetFilename("notes.json").
    AddFilter("JSON Files", "*.json").
    PromptForSingleSelection()
```

대응하는 메서드는 메시지 대화 상자의 경우 `app.Dialog.Info() / Question() / Warning() / Error()`이고, 파일 대화 상자의 경우 `app.Dialog.OpenFile() / SaveFile()` 및 `*WithOptions` 변형입니다.

### JSON 태그 매핑

JavaScript 속성 접근과 일치하도록 Go 구조체의 JSON 태그는 소문자여야 합니다:

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

### 디바운스가 적용된 자동 저장

500ms 디바운스를 적용하면 불필요한 백엔드 호출이 줄어듭니다:

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

## 다음 단계

- 노트를 정리할 수 있도록 카테고리 또는 태그를 추가하세요
- 검색 및 필터링 기능을 구현하세요
- WYSIWYG 편집기를 사용한 서식 있는 텍스트 편집 기능을 추가하세요
- 노트를 클라우드 스토리지와 동기화하세요
- 자주 사용하는 작업에 대한 키보드 단축키를 추가하세요
