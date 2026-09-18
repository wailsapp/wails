---
title: "Заметки"
description: "Создание приложения для заметок с файловыми операциями и нативными диалоговыми окнами"
slug: "tutorials/03-notes-vanilla"
sourcePath: "tutorials/03-notes-vanilla.md"
---

В этом руководстве вы создадите с помощью Wails v3 приложение для заметок, демонстрирующее файловые операции, нативные диалоговые окна и современные подходы к разработке настольных приложений.

![Снимок экрана приложения «Заметки»](/assets/notes-app.png)

## Что вы создадите

- Полноценное приложение для создания, редактирования и удаления заметок
- Нативные диалоговые окна сохранения и открытия файлов для импорта и экспорта заметок
- Автосохранение при вводе с устранением дребезга для сокращения количества лишних обновлений
- Профессиональный двухколоночный макет (боковая панель и редактор) в стиле Apple Notes

## Чему вы научитесь

- Использовать нативные файловые диалоговые окна в Wails (`SaveFileDialog`, `OpenFileDialog`, `InfoDialog`)
- Работать с файлами JSON для постоянного хранения данных
- Реализовывать автосохранение с устранением дребезга
- Создавать профессиональные интерфейсы настольных приложений с помощью современного CSS
- Правильно сериализовать структуры Go в JSON

## Настройка проекта

@steps
### Создайте новый проект Wails
```bash
wails3 init -n notes-app -t vanilla
cd notes-app
```

### Создайте NotesService
Создайте в корне проекта новый файл `notesservice.go`:

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

**Что здесь происходит:**

- **Структура Note**: определяет структуру данных с тегами JSON в нижнем регистре для правильной сериализации
- **Операции CRUD**: GetAll, Create, Update и Delete для управления заметками в памяти
- **Файловые диалоговые окна**: для доступа к нативным диалоговым окнам используются `application.Get().Dialog.SaveFile()` и `application.Get().Dialog.OpenFile()`
- **Информационные диалоговые окна**: с помощью `application.Get().Dialog.Info()` отображают сообщения об успешном выполнении
- **Генерация идентификаторов**: простой генератор идентификаторов на основе меток времени

### Обновите main.go
Замените содержимое файла `main.go`:

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

**Что здесь происходит:**

- `NotesService` регистрируется в приложении
- Создаётся окно размером 1000x700 в стиле Apple Notes
- Настраивается корректное поведение в macOS: приложение завершается после закрытия последнего окна

### Создайте структуру HTML
Замените содержимое `frontend/index.html`:

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

**Что здесь происходит:**

- **Двухколоночный макет**: боковая панель со списком заметок и основная область с редактором
- **Пустое состояние**: отображается, если ни одна заметка не выбрана
- **Среда выполнения Wails**: должна быть загружена до скрипта модуля

### Добавьте стили CSS
Замените содержимое `frontend/public/style.css`:

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

**Что здесь происходит:**

- Дизайн в стиле Apple с аккуратной типографикой и цветовой палитрой
- Макет на основе Flexbox с адаптивными боковой панелью и редактором
- Выделение активной заметки синим фоном
- Плавные переходы при наведении указателя

### Реализуйте логику JavaScript
Замените содержимое `frontend/src/main.js`:

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

**Что здесь происходит:**

- **Автосохранение**: задержка устранения дребезга длительностью 500 мс предотвращает избыточные вызовы серверной части при вводе
- **Доступ к свойствам**: используются имена свойств в нижнем регистре (`.id`, `.title`), соответствующие тегам JSON в Go
- **Управление фокусом**: при создании новой заметки фокус автоматически устанавливается на заголовок, а его текст выделяется
- **Удаление без подтверждения**: браузерная функция `confirm()` не работает в WebView Wails
- **Файловые операции**: нативные диалоговые окна обеспечивают сохранение и загрузку с надлежащей обработкой ошибок

### Запустите приложение
```bash
wails3 dev
```

После запуска приложения вы сможете:

- Нажать «+ Новая заметка», чтобы создать заметку
- Редактировать заголовок и содержимое (автосохранение выполняется через 500 мс)
- Нажимать на заметки в боковой панели, чтобы переключаться между ними
- Нажать «Удалить», чтобы удалить текущую заметку
- Нажать «Сохранить», чтобы экспортировать заметки в формате JSON
- Нажать «Загрузить», чтобы импортировать ранее сохранённые заметки

@end

## Ключевые понятия

### Файловые диалоговые окна и окна сообщений

В Wails v3 **нет** конструкторов диалоговых окон уровня пакета — каждое диалоговое окно создаётся через диспетчер `app.Dialog`. В сервисе получите приложение с помощью `application.Get()` (эта функция возвращает запущенный экземпляр `*application.App`):

```go
// Correct — manager-based dialogs.
app := application.Get()
path, err := app.Dialog.SaveFile().
    SetFilename("notes.json").
    AddFilter("JSON Files", "*.json").
    PromptForSingleSelection()
```

Соответствующие методы: `app.Dialog.Info() / Question() / Warning() / Error()` для диалоговых окон сообщений и `app.Dialog.OpenFile() / SaveFile()` (а также варианты `*WithOptions`) для файловых диалоговых окон.

### Сопоставление тегов JSON

Теги JSON в структурах Go должны быть записаны в нижнем регистре, чтобы соответствовать обращениям к свойствам в JavaScript:

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

### Автосохранение с устранением дребезга

Задержка устранения дребезга длительностью 500 мс сокращает количество лишних вызовов серверной части:

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

## Дальнейшие шаги

- Добавьте категории или теги для упорядочивания заметок
- Реализуйте поиск и фильтрацию
- Добавьте редактирование форматированного текста с помощью WYSIWYG-редактора
- Синхронизируйте заметки с облачным хранилищем
- Добавьте сочетания клавиш для часто используемых операций
