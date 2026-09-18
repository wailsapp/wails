---
title: "TODO-Liste"
description: "Erstellen Sie eine vollständige TODO-Listenanwendung mit CRUD-Operationen"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

In diesem Tutorial erstellen Sie eine voll funktionsfähige TODO-Listenanwendung. Es baut auf dem Tutorial zum QR-Code-Dienst auf: Sie lernen, Zustand zu verwalten, mehrere Operationen zu verarbeiten und eine ansprechende Benutzeroberfläche zu erstellen.

**Was Sie erstellen:**

- Eine vollständige TODO-App mit Funktionen zum Hinzufügen, Abschließen und Löschen
- Threadsichere Zustandsverwaltung (wichtig für Desktop-Apps)
- Modernes, glasmorphes UI-Design
- Alles mit Vanilla JavaScript – keine Frameworks erforderlich

**Was Sie lernen:**

- CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen)
- Veränderlichen Zustand in Go sicher verwalten
- Benutzereingaben verarbeiten und validieren
- Responsive Benutzeroberflächen erstellen, die sich nativ anfühlen

![TODO-Listenanwendung](/assets/todo-app.png)

**Benötigte Zeit:** 20 Minuten

## Projekt erstellen

@steps
### Projekt generieren
Erstellen Sie zunächst ein neues Wails-Projekt. Wir verwenden die standardmäßige Vanilla-Vorlage, die uns einen aufgeräumten Ausgangspunkt bietet:

```bash
wails3 init -n todo-app
cd todo-app
```

Dadurch entsteht ein neues Projekt mit der grundlegenden Struktur: Das Go-Backend liegt im Stammverzeichnis, der Frontend-Code im Verzeichnis `frontend/`.

### TODO-Dienst erstellen
Der TODO-Dienst verwaltet den Zustand unserer Anwendung und stellt Methoden für CRUD-Operationen bereit. Anders als bei einem Webserver, bei dem jede Anfrage isoliert ist, können in Desktop-Apps mehrere Operationen gleichzeitig ausgeführt werden. Daher benötigen wir eine threadsichere Zustandsverwaltung.

Löschen Sie `greetservice.go` und erstellen Sie die neue Datei `todoservice.go`:

```go {title="todoservice.go"}
package main

import (
    "errors"
    "sync"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

type TodoService struct {
    todos  []Todo
    nextID int
    mu     sync.RWMutex
}

func NewTodoService() *TodoService {
    return &TodoService{
        todos:  []Todo{},
        nextID: 1,
    }
}

func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }

    t.mu.Lock()
    defer t.mu.Unlock()

    todo := Todo{
        ID:        t.nextID,
        Title:     title,
        Completed: false,
    }
    t.todos = append(t.todos, todo)
    t.nextID++

    return &todo, nil
}

func (t *TodoService) Toggle(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos[i].Completed = !t.todos[i].Completed
            return nil
        }
    }
    return errors.New("todo not found")
}

func (t *TodoService) Delete(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i, todo := range t.todos {
        if todo.ID == id {
            t.todos = append(t.todos[:i], t.todos[i+1:]...)
            return nil
        }
    }
    return errors.New("todo not found")
}
```

**Was hier geschieht:**

**Die Struktur `Todo`:**

- Definiert die Struktur unserer Daten mit den Feldern ID, Title und Completed
- `json:`-Tags geben Go vor, wie diese Struktur für das Frontend in JSON konvertiert wird
- Jedes Feld wird exportiert (beginnt mit einem Großbuchstaben), damit der Bindings-Generator es erkennen kann

**Die Struktur `TodoService`:**

- `todos []Todo` – ein Slice, das alle unsere TODO-Einträge enthält
- `nextID int` – verwaltet die nächste zu vergebende ID (simuliert Auto-Inkrementierung)
- `mu sync.RWMutex` – ein Lese-/Schreib-Mutex für threadsicheren Zugriff

**Threadsicherheit mit `sync.RWMutex`:**

- Desktop-Apps können mehrere gleichzeitige Operationen über die Benutzeroberfläche ausführen
- `RLock()` ermöglicht mehrere gleichzeitige Lesezugriffe (z. B. mehrere Aufrufe von `GetAll`)
- `Lock()` gewährt exklusiven Schreibzugriff (z. B. für `Add`, `Toggle` und `Delete`)
- `defer` stellt sicher, dass Sperren auch dann freigegeben werden, wenn die Funktion vorzeitig zurückkehrt oder eine Panic auftritt

**Die Methoden:**

- `GetAll()` – gibt alle TODO-Einträge zurück (verwendet eine Lesesperre, da keine Daten geändert werden)
- `Add(title)` – erstellt einen neuen TODO-Eintrag, validiert die Eingabe und erhöht die ID
- `Toggle(id)` – kehrt den Abschlussstatus eines TODO-Eintrags um
- `Delete(id)` – entfernt einen TODO-Eintrag aus dem Slice

**Fehlerbehandlung:**

- Gemäß den Go-Konventionen geben wir `error` als letzten Wert zurück
- Leere Titel werden abgelehnt
- Operationen für nicht vorhandene TODO-Einträge geben Fehler zurück
- Diese Fehler werden im Frontend zu JavaScript-Ausnahmen

### main.go aktualisieren
Registrieren Sie den TODO-Dienst bei Ihrer Wails-Anwendung. Suchen Sie in `main.go` den Abschnitt `Services` und ersetzen Sie den GreetService durch unseren TodoService:

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**Was hier geschieht:**

- Wir entfernen den standardmäßigen GreetService und fügen stattdessen unseren TodoService hinzu
- `application.NewService()` umschließt unseren Dienst, damit Wails ihn verwalten kann
- Wails generiert automatisch JavaScript-Bindings für alle öffentlichen Methoden dieses Dienstes

### Frontend-Benutzeroberfläche erstellen
Erstellen wir nun das Frontend. Hier rufen wir unsere Go-Methoden auf und stellen die Benutzeroberfläche dar. Wir verwenden Vanilla JavaScript, um alles einfach zu halten und Ihnen direkt zu zeigen, wie die Bindings funktionieren.

Ersetzen Sie `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {TodoService} from "../bindings/changeme";

async function loadTodos() {
    const todos = await TodoService.GetAll();
    const list = document.getElementById('todo-list');

    list.innerHTML = todos.map(todo => `
        <div class="todo ${todo.completed ? 'completed' : ''}">
            <input type="checkbox"
                   ${todo.completed ? 'checked' : ''}
                   onchange="toggleTodo(${todo.id})">
            <span>${todo.title}</span>
            <button onclick="deleteTodo(${todo.id})">Delete</button>
        </div>
    `).join('');
}

window.addTodo = async () => {
    const input = document.getElementById('todo-input');
    const title = input.value.trim();

    if (title) {
        await TodoService.Add(title);
        input.value = '';
        await loadTodos();
    }
}

window.toggleTodo = async (id) => {
    await TodoService.Toggle(id);
    await loadTodos();
}

window.deleteTodo = async (id) => {
    await TodoService.Delete(id);
    await loadTodos();
}

// Load todos on startup
loadTodos();
```

**Was hier geschieht:**

**Bindings importieren:**

- `import {TodoService} from "../bindings/changeme"` – importiert die automatisch generierten Go-Bindings
- Hinweis: `changeme` entspricht Ihrem tatsächlichen Modulnamen aus `go.mod`

**Die Funktion `loadTodos()`:**

- Ruft `TodoService.GetAll()` auf, um alle TODO-Einträge aus Go abzurufen
- Erstellt mithilfe von Template-Literalen HTML für jeden TODO-Eintrag
- Fügt die Klasse `completed` für die Formatierung dynamisch hinzu oder entfernt sie
- Verwendet `onclick`-Attribute, um Schaltflächen mit unseren Funktionen zu verknüpfen
- Fügt das gesamte HTML zusammen und bindet es in das DOM ein

**Die CRUD-Funktionen:**

- `addTodo()` – validiert die Eingabe, ruft die Go-Methode `Add` auf und aktualisiert die Liste
- `toggleTodo(id)` – ruft die Go-Methode `Toggle` auf und aktualisiert die Liste
- `deleteTodo(id)` – ruft die Go-Methode `Delete` auf und aktualisiert die Liste
- Alle Funktionen sind asynchron, da Go-Aufrufe Promises zurückgeben

**Warum an window anhängen:**

- `window.addTodo = ...` macht Funktionen über `onclick`-Attribute im HTML zugänglich
- Dies ist ein einfaches Muster für Vanilla JS (Frameworks handhaben dies anders).
- In der Produktion könnten Sie stattdessen eine geeignete Ereignisdelegation verwenden.

**Das Aktualisierungsmuster:**

- Nach jeder Änderung (Hinzufügen/Umschalten/Löschen) rufen wir `loadTodos()` erneut auf.
- Dadurch bleibt die Benutzeroberfläche mit dem Go-Zustand synchron.
- Alternative: Lassen Sie die Go-Methoden den neuen Zustand zurückgeben, um den zweiten Aufruf zu vermeiden.

### HTML aktualisieren
Das HTML stellt die Struktur unserer TODO-App bereit. Es ist minimalistisch und semantisch – die eigentliche Funktionalität steckt in JavaScript und CSS.

Ersetzen Sie `frontend/index.html`:

```html {title="frontend/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>TODO App</title>
    <link rel="stylesheet" href="./style.css"/>
</head>
<body>
    <div class="container">
        <h1>My TODOs</h1>
        <div class="card">
            <div class="input-box">
                <input type="text"
                       id="todo-input"
                       class="input"
                       placeholder="Add a new todo..."
                       onkeypress="if(event.key==='Enter') addTodo()">
                <button class="btn" onclick="addTodo()">Add</button>
            </div>
            <div id="todo-list"></div>
        </div>
    </div>
    <script type="module" src="./src/main.js"></script>
</body>
</html>
```

**Was hier geschieht:**

**Die Struktur:**

- `container` – zentriert unsere App und begrenzt ihre Breite
- `card` – die weiße Hauptkarte, die alle Inhalte enthält
- `input-box` – Flex-Container für das Eingabefeld und die Schaltfläche „Hinzufügen“
- `todo-list` – hier fügt JavaScript die einzelnen Todos ein

**Ereignisbehandlung:**

- `onkeypress="if(event.key==='Enter') addTodo()"` – fügt beim Drücken der Eingabetaste ein Todo hinzu
- `onclick="addTodo()"` – fügt beim Klicken auf die Schaltfläche ein Todo hinzu
- Inline-Ereignishandler eignen sich gut für einfache Vanilla-JS-Apps.

**Modulskript:**

- `<script type="module">` ermöglicht die Verwendung von ES6-Importen.
- Unsere Datei `main.js` kann die Bindings importieren und modernes JavaScript verwenden.

### App gestalten
Das CSS erzeugt ein modernes Glassmorphism-Design mit fließenden Übergängen. Unser Ziel ist eine hochwertige Anmutung, die eine angenehme Nutzung der App ermöglicht.

Ersetzen Sie `frontend/public/style.css`:

```css {title="frontend/public/style.css"}
:root {
    font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
    "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
    font-size: 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: rgba(255, 255, 255, 0.87);
}

body {
    margin: 0;
    display: flex;
    place-items: center;
    justify-content: center;
    min-height: 100vh;
}

.container {
    width: 100%;
    max-width: 600px;
    padding: 20px;
}

h1 {
    text-align: center;
    color: white;
    font-size: 2.5em;
    font-weight: 300;
    margin: 0 0 30px 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 30px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.input-box {
    display: flex;
    gap: 10px;
    margin-bottom: 25px;
}

.input {
    flex: 1;
    border: 2px solid #e0e0e0;
    border-radius: 12px;
    height: 50px;
    padding: 0 20px;
    font-size: 16px;
    transition: all 0.3s ease;
}

.input:focus {
    border-color: #667eea;
    outline: none;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.btn {
    height: 50px;
    padding: 0 30px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

#todo-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.todo {
    display: flex;
    align-items: center;
    padding: 18px 20px;
    background: white;
    border: 2px solid #f0f0f0;
    border-radius: 12px;
    transition: all 0.3s ease;
    gap: 15px;
}

.todo:hover {
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
    transform: translateX(4px);
}

.todo.completed {
    opacity: 0.6;
}

.todo.completed span {
    text-decoration: line-through;
    color: #999;
}

.todo input[type="checkbox"] {
    width: 24px;
    height: 24px;
    cursor: pointer;
    appearance: none;
    -webkit-appearance: none;
    border: 2px solid #667eea;
    border-radius: 6px;
    position: relative;
    transition: all 0.3s ease;
    flex-shrink: 0;
}

.todo input[type="checkbox"]:hover {
    background: rgba(102, 126, 234, 0.1);
}

.todo input[type="checkbox"]:checked {
    background: #667eea;
    border-color: #667eea;
}

.todo input[type="checkbox"]:checked::after {
    content: '✓';
    position: absolute;
    color: white;
    font-size: 16px;
    font-weight: bold;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

.todo span {
    flex: 1;
    font-size: 16px;
    color: #333;
}

.todo button {
    padding: 8px 16px;
    background: #ff4757;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    opacity: 0;
    flex-shrink: 0;
}

.todo:hover button {
    opacity: 1;
}

.todo button:hover {
    background: #ee5a6f;
    transform: scale(1.05);
}

#todo-list:empty::before {
    content: "No todos yet. Add one above!";
    display: block;
    text-align: center;
    padding: 40px 20px;
    color: #999;
}
```

**Was hier geschieht:**

**Glassmorphism-Design:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` – violetter Farbverlauf im Hintergrund
- `backdrop-filter: blur(10px)` – erzeugt den Milchglaseffekt auf der Karte
- `rgba(255, 255, 255, 0.95)` – halbtransparentes Weiß für den Glaseffekt

**Benutzerdefinierte Gestaltung des Kontrollkästchens:**

- `appearance: none` entfernt das standardmäßige Kontrollkästchen des Browsers.
- Mit `::after` erstellen wir ein benutzerdefiniertes abgerundetes Quadrat mit einem Häkchen.
- Das Häkchen erscheint bei `checked` als Unicode-Zeichen ✓.

**Interaktionen beim Darüberfahren:**

- Todos verschieben sich beim Darüberfahren nach rechts (`transform: translateX(4px)`).
- Die Schaltfläche zum Löschen bleibt bis zum Darüberfahren ausgeblendet (`opacity: 0` → `opacity: 1`).
- Schaltflächen werden beim Darüberfahren leicht vergrößert, um eine taktile Rückmeldung zu geben.

**Leerzustand:**

- `#todo-list:empty::before` zeigt eine Meldung an, wenn keine Todos vorhanden sind.
- Reine CSS-Lösung – kein JavaScript erforderlich

### App ausführen
Sehen wir uns die App in Aktion an! Starten Sie den Entwicklungsserver:

```bash
wails3 dev
```

Die App wird kompiliert und geöffnet. Probieren Sie sie aus:

- Geben Sie ein Todo ein und drücken Sie die Eingabetaste oder klicken Sie auf „Hinzufügen“.
- Klicken Sie auf das Kontrollkästchen, um das Todo als erledigt zu markieren.
- Fahren Sie mit der Maus über ein Todo, damit die Schaltfläche zum Löschen erscheint.
- Die Benutzeroberfläche wird sofort aktualisiert – hier ist unser Aktualisierungsmuster am Werk.

**Was geschieht:**

- Wails hat automatisch Bindings für Ihre TodoService-Methoden generiert.
- Der Entwicklungsmodus umfasst Hot Reload – ändern Sie das CSS und beobachten Sie die Aktualisierung.
- Ihr Go-Code wird nativ ausgeführt – eine Übersetzung oder Interpretation ist nicht erforderlich.

@end

## Funktionsweise

### Threadsichere Zustandsverwaltung

`sync.RWMutex` ermöglicht sicheren gleichzeitigen Zugriff:

```go
func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()  // Read lock - multiple readers allowed
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    t.mu.Lock()  // Write lock - exclusive access
    defer t.mu.Unlock()
    // ... mutations
}
```

**Warum dies wichtig ist:**

- Mehrere Frontend-Aufrufe können gleichzeitig erfolgen.
- Leseoperationen blockieren einander nicht.
- Schreiboperationen erhalten exklusiven Zugriff.
- `defer` stellt sicher, dass Sperren immer freigegeben werden.

### Fehlerbehandlung

Der Dienst gibt bei ungültigen Operationen Fehler zurück:

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

Im Frontend können Sie diese Fehler abfangen:

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### Zustandssynchronisierung

Nach jeder Änderung laden wir die vollständige Liste neu:

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**Alternativer Ansatz:** Geben Sie aus jeder Methode die aktualisierte Liste zurück, um den zweiten Aufruf zu vermeiden.

## Erweiterungen

### Statistiken hinzufügen

Fügen Sie Folgendes zu `todoservice.go` hinzu:

```go
type TodoStats struct {
    Total     int `json:"total"`
    Completed int `json:"completed"`
    Active    int `json:"active"`
}

func (t *TodoService) GetStats() TodoStats {
    t.mu.RLock()
    defer t.mu.RUnlock()

    stats := TodoStats{
        Total: len(t.todos),
    }

    for _, todo := range t.todos {
        if todo.Completed {
            stats.Completed++
        } else {
            stats.Active++
        }
    }

    return stats
}
```

Im Frontend anzeigen:

```javascript
async function loadTodos() {
    const [todos, stats] = await Promise.all([
        TodoService.GetAll(),
        TodoService.GetStats()
    ]);

    // Display stats
    document.getElementById('stats').textContent =
        `${stats.active} active, ${stats.completed} completed`;

    // ... render todos
}
```

### „Erledigte löschen“ hinzufügen

```go
func (t *TodoService) ClearCompleted() int {
    t.mu.Lock()
    defer t.mu.Unlock()

    removed := 0
    newTodos := []Todo{}

    for _, todo := range t.todos {
        if !todo.Completed {
            newTodos = append(newTodos, todo)
        } else {
            removed++
        }
    }

    t.todos = newTodos
    return removed
}
```

### Persistenz hinzufügen

Fügen Sie für Produktions-Apps Persistenz mit der Speicherlösung hinzu, die am besten zu Ihrer Anwendung passt, beispielsweise SQLite, PostgreSQL oder einem gehosteten Dienst.

## Für die Produktion erstellen

Wenn du deine TODO-App veröffentlichen möchtest, erstelle sie für den Produktivbetrieb:

```bash
wails3 build
```

Dadurch wird eine optimierte native ausführbare Datei in `bin/` erstellt:

- Kompiliert deinen Go-Code mit Optimierungen
- Erstellt dein Frontend für den Produktivbetrieb
- Bündelt alles in einer einzigen ausführbaren Datei
- Die resultierende App ist normalerweise 10-20 MB groß (im Vergleich zu über 150 MB bei Electron)

Du kannst die ausführbare Datei direkt starten – ohne erforderliche Laufzeitumgebung und ohne zu startende Server. Es handelt sich um eine echte native Anwendung.

## Was du erstellt hast

Du hast gerade eine vollständige TODO-Anwendung mit folgenden Funktionen erstellt:

**Vollständige CRUD-Implementierung:**

- Einen Dienst mit Operationen zum Erstellen, Lesen, Aktualisieren und Löschen erstellt
- Eingabevalidierung und Fehlerbehandlung hinzugefügt
- Gelernt, wie Go-Fehler zu JavaScript-Ausnahmen werden

**Threadsichere Zustandsverwaltung:**

- `sync.RWMutex` verwendet, um gleichzeitige Zugriffe sicher zu handhaben
- Den Unterschied zwischen Lesesperren (RLock) und Schreibsperren (Lock) verstanden
- Gesehen, wie `defer` durch garantiertes Aufräumen Deadlocks verhindert

**Moderne, ausgereifte Benutzeroberfläche:**

- Eine Benutzeroberfläche im Glassmorphism-Stil mit Farbverläufen und Weichzeichnungseffekten erstellt
- Individuell gestaltete Kontrollkästchen ohne Framework erstellt
- Hover-Interaktionen und Übergänge für ein natives Bediengefühl hinzugefügt
- Einen Leerzustand ausschließlich mit CSS implementiert

**Grundlagen von Wails:**

- Dienstregistrierung und automatische Generierung von Bindings
- Aufrufen von Go-Methoden aus JavaScript mit async/await
- Zustandssynchronisierung zwischen Go und dem Frontend
- Erstellen und Paketieren einer nativen Desktop-App

## Nächste Schritte

Nachdem du CRUD-Operationen und Zustandsverwaltung verstanden hast, probiere Folgendes aus:

- **Persistenz hinzufügen:** Sorge mit SQLite oder einer anderen Speicherlösung dafür, dass TODOs auch nach einem Neustart der App erhalten bleiben
- **Weitere Funktionen hinzufügen:** Filterung (alle/aktive/abgeschlossene), Bearbeitung vorhandener TODOs und Massenoperationen
- **Das Notes-Tutorial erkunden:** Erfahre unter [Notes](/tutorials/03-notes-vanilla/), wie Dateioperationen funktionieren
- **Etwas Reales entwickeln:** Nutze diese Konzepte und entwickle deine eigene App!
