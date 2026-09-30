---
title: "Notes"
description: "Créez une application de prise de notes avec des opérations sur les fichiers et des boîtes de dialogue natives"
slug: "tutorials/03-notes-vanilla"
sourcePath: "tutorials/03-notes-vanilla.md"
---

Dans ce tutoriel, vous allez créer avec Wails v3 une application de prise de notes qui illustre les opérations sur les fichiers, les boîtes de dialogue natives et les modèles modernes de conception d’applications de bureau.

![Capture d’écran de l’application de prise de notes](/assets/notes-app.png)

## Ce que vous allez créer

- Une application de prise de notes complète permettant de créer, modifier et supprimer des notes
- Des boîtes de dialogue natives d’ouverture et d’enregistrement de fichiers pour importer et exporter des notes
- Un enregistrement automatique pendant la saisie, avec temporisation, afin de réduire les mises à jour inutiles
- Une mise en page professionnelle à deux colonnes (barre latérale et éditeur) inspirée d’Apple Notes

## Ce que vous allez apprendre

- Utiliser les boîtes de dialogue de fichiers natives dans Wails (`SaveFileDialog`, `OpenFileDialog`, `InfoDialog`)
- Utiliser des fichiers JSON pour la persistance des données
- Mettre en œuvre des modèles d’enregistrement automatique temporisé
- Créer des interfaces utilisateur professionnelles pour applications de bureau avec du CSS moderne
- Sérialiser correctement des structures Go en JSON

## Configuration du projet

@steps
### Créer un projet Wails
```bash
wails3 init -n notes-app -t vanilla
cd notes-app
```

### Créer le NotesService
Créez un fichier `notesservice.go` à la racine du projet :

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

**Voici ce qui se passe :**

- **Structure Note** : définit la structure des données avec des balises JSON en minuscules pour assurer une sérialisation correcte
- **Opérations CRUD** : GetAll, Create, Update et Delete permettent de gérer les notes en mémoire
- **Boîtes de dialogue de fichiers** : utilisent `application.Get().Dialog.SaveFile()` et `application.Get().Dialog.OpenFile()` pour accéder aux boîtes de dialogue natives
- **Boîtes de dialogue d’information** : affichent des messages de réussite à l’aide de `application.Get().Dialog.Info()`
- **Génération des identifiants** : générateur d’identifiants simple fondé sur l’horodatage

### Mettre à jour main.go
Remplacez le contenu de `main.go` :

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

**Voici ce qui se passe :**

- Enregistre `NotesService` auprès de l’application
- Crée une fenêtre aux dimensions 1000x700, inspirée d’Apple Notes
- Configure le comportement approprié sous macOS afin de quitter l’application lorsque la dernière fenêtre se ferme

### Créer la structure HTML
Remplacez `frontend/index.html` :

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

**Voici ce qui se passe :**

- **Mise en page à deux colonnes** : une barre latérale pour la liste des notes et une zone principale pour l’éditeur
- **État vide** : s’affiche lorsqu’aucune note n’est sélectionnée
- **Environnement d’exécution Wails** : doit être chargé avant le script de module

### Ajouter les styles CSS
Remplacez `frontend/public/style.css` :

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

**Voici ce qui se passe :**

- Design dans le style d’Apple, avec une typographie et des couleurs épurées
- Mise en page Flexbox pour une barre latérale et un éditeur adaptatifs
- Mise en évidence de la note active sur un arrière-plan bleu
- Transitions fluides au survol

### Implémenter la logique JavaScript
Remplacez `frontend/src/main.js` :

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

**Voici ce qui se passe :**

- **Enregistrement automatique** : une temporisation de 500 ms évite les appels excessifs au backend pendant la saisie
- **Accès aux propriétés** : utilise des noms de propriétés en minuscules (`.id`, `.title`), conformément aux balises JSON de Go
- **Gestion du focus** : place automatiquement le focus sur le titre et le sélectionne lors de la création d’une note
- **Suppression sans confirmation** : la fonction `confirm()` du navigateur ne fonctionne pas dans les vues web de Wails
- **Opérations sur les fichiers** : les boîtes de dialogue natives gèrent l’enregistrement et le chargement avec une gestion appropriée des erreurs

### Exécuter l’application
```bash
wails3 dev
```

L’application démarre et vous pouvez alors :

- Cliquer sur « + New Note » pour créer des notes
- Modifier le titre et le contenu, qui sont enregistrés automatiquement après 500 ms
- Cliquer sur les notes dans la barre latérale pour passer de l’une à l’autre
- Cliquer sur « Delete » pour supprimer la note active
- Cliquer sur « Save » pour exporter les notes au format JSON
- Cliquer sur « Load » pour importer des notes précédemment enregistrées

@end

## Concepts clés

### Boîtes de dialogue de fichiers et de messages

Dans Wails v3, il n’existe **aucun** constructeur de boîte de dialogue au niveau du paquet : chaque boîte de dialogue est créée au moyen du gestionnaire `app.Dialog`. Depuis un service, récupérez l’application avec `application.Get()`, qui renvoie l’instance `*application.App` en cours d’exécution :

```go
// Correct — manager-based dialogs.
app := application.Get()
path, err := app.Dialog.SaveFile().
    SetFilename("notes.json").
    AddFilter("JSON Files", "*.json").
    PromptForSingleSelection()
```

Les méthodes correspondantes sont `app.Dialog.Info() / Question() / Warning() / Error()` pour les boîtes de dialogue de messages et `app.Dialog.OpenFile() / SaveFile()`, ainsi que les variantes `*WithOptions`, pour les boîtes de dialogue de fichiers.

### Correspondance des balises JSON

Les balises JSON des structures Go doivent être en minuscules pour correspondre à l’accès aux propriétés en JavaScript :

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

### Enregistrement automatique temporisé

La temporisation de 500 ms réduit les appels inutiles au backend :

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

## Étapes suivantes

- Ajoutez des catégories ou des étiquettes pour organiser les notes
- Implémentez la recherche et le filtrage
- Ajoutez l’édition de texte enrichi à l’aide d’un éditeur WYSIWYG
- Synchronisez les notes avec un stockage dans le cloud
- Ajoutez des raccourcis clavier pour les opérations courantes
