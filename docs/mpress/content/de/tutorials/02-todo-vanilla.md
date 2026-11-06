---
title: "TODO-Liste"
description: "Erstellen Sie eine vollständige TODO-Listen-Anwendung mit CRUD-Operationen"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

In diesem Tutorial erstellen Sie eine voll funktionsfähige TODO-Listen-Anwendung. Es baut auf dem Tutorial zum QR-Code-Dienst auf: Sie lernen, wie Sie Zustand verwalten, mehrere Operationen verarbeiten und eine ausgereifte Benutzeroberfläche erstellen.

**Was Sie erstellen:**

- Eine vollständige TODO-App mit Funktionen zum Hinzufügen, Erledigen und Löschen
- Threadsichere Zustandsverwaltung (wichtig für Desktop-Apps)
- Modernes Benutzeroberflächendesign im Glassmorphismus-Stil
- Alles mit Vanilla JavaScript – keine Frameworks erforderlich

**Was Sie lernen:**

- CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen)
- Sichere Verwaltung veränderlichen Zustands in Go
- Verarbeiten und Validieren von Benutzereingaben
- Erstellen responsiver Benutzeroberflächen, die sich wie native Anwendungen anfühlen

![TODO-Listen-Anwendung](/assets/todo-app.png)

**Benötigte Zeit:** 20 Minuten

## Projekt erstellen

@steps
### Projekt generieren
Erstellen Sie zunächst ein neues Wails-Projekt. Wir verwenden die standardmäßige Vanilla-Vorlage, die uns einen übersichtlichen Ausgangspunkt bietet:

```bash
wails3 init -n todo-app -t vanilla-js
cd todo-app
```

Dadurch wird ein neues Projekt mit der grundlegenden Struktur erstellt: das Go-Backend im Stammverzeichnis und der Frontend-Code im Verzeichnis `frontend/`.

### TODO-Dienst erstellen
Der TODO-Dienst verwaltet den Zustand unserer Anwendung und stellt Methoden für CRUD-Operationen bereit. Anders als bei einem Webserver, auf dem jede Anfrage isoliert ist, können in Desktop-Apps mehrere Operationen gleichzeitig ausgeführt werden. Daher benötigen wir eine threadsichere Zustandsverwaltung.

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
- Jedes Feld ist exportiert (beginnt mit einem Großbuchstaben), damit der Bindings-Generator es erkennen kann

**Die Struktur `TodoService`:**

- `todos []Todo` – ein Slice, das alle unsere TODO-Einträge enthält
- `nextID int` – speichert die nächste zu vergebende ID (simuliert die automatische Inkrementierung)
- `mu sync.RWMutex` – ein Lese-/Schreib-Mutex für threadsicheren Zugriff

**Threadsicherheit mit `sync.RWMutex`:**

- In Desktop-Apps können über die Benutzeroberfläche mehrere Operationen gleichzeitig ausgeführt werden
- `RLock()` ermöglicht mehrere gleichzeitige Lesezugriffe (z. B. mehrere Aufrufe von `GetAll`)
- `Lock()` gewährt exklusiven Zugriff für Schreibvorgänge (z. B. `Add`, `Toggle`, `Delete`)
- `defer` stellt sicher, dass Sperren selbst dann aufgehoben werden, wenn die Funktion vorzeitig zurückkehrt oder eine Panic auslöst

**Die Methoden:**

- `GetAll()` – gibt alle TODO-Einträge zurück (verwendet eine Lesesperre, da keine Daten geändert werden)
- `Add(title)` – erstellt einen neuen TODO-Eintrag, validiert die Eingabe und erhöht die ID
- `Toggle(id)` – kehrt den Erledigt-Status eines TODO-Eintrags um
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
Erstellen wir nun das Frontend. Hier rufen wir unsere Go-Methoden auf und zeigen die Benutzeroberfläche an. Wir verwenden Vanilla JavaScript, um alles einfach zu halten und Ihnen unmittelbar zu zeigen, wie die Bindings funktionieren.

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
- Verwendet `onclick`-Attribute, um Schaltflächen mit unseren Funktionen zu verbinden
- Fügt das gesamte HTML zusammen und bindet es in das DOM ein

**Die CRUD-Funktionen:**

- `addTodo()` – validiert die Eingabe, ruft die Go-Methode `Add` auf und aktualisiert die Liste
- `toggleTodo(id)` – ruft die Go-Methode `Toggle` auf und aktualisiert die Liste
- `deleteTodo(id)` – ruft die Go-Methode `Delete` auf und aktualisiert die Liste
- Alle Funktionen sind asynchron, da Go-Aufrufe Promises zurückgeben

**Warum an window anhängen:**

- `window.addTodo = ...` macht Funktionen über HTML-Attribute vom Typ `onclick` zugänglich
- Dies ist ein einfaches Muster für Vanilla JavaScript (Frameworks handhaben dies anders)
- In einer Produktionsanwendung könnten Sie stattdessen eine geeignete Ereignisdelegation verwenden

**Das Aktualisierungsmuster:**

- Nach jeder Änderung (Hinzufügen/Umschalten/Löschen) rufen wir `loadTodos()` erneut auf
- Dadurch bleibt die Benutzeroberfläche mit dem Go-Zustand synchronisiert
- Alternative: Lassen Sie die Go-Methoden den neuen Zustand zurückgeben, um den zweiten Aufruf zu vermeiden

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

**Was geschieht hier?**

**Die Struktur:**

- `container` – zentriert unsere App und begrenzt ihre Breite
- `card` – die weiße Hauptkarte, die alle Elemente enthält
- `input-box` – Flex-Container für das Eingabefeld und die Schaltfläche „Hinzufügen“
- `todo-list` – hier fügt JavaScript die einzelnen TODO-Einträge ein

**Ereignisbehandlung:**

- `onkeypress="if(event.key==='Enter') addTodo()"` – fügt beim Drücken der Eingabetaste einen TODO-Eintrag hinzu
- `onclick="addTodo()"` – fügt beim Anklicken der Schaltfläche einen TODO-Eintrag hinzu
- Inline-Ereignishandler eignen sich gut für einfache Apps mit Vanilla JavaScript

**Modulskript:**

- Mit `<script type="module">` können wir ES6-Importe verwenden
- Unsere Datei `main.js` kann die Bindings importieren und modernes JavaScript verwenden

### App gestalten
Das CSS erzeugt ein modernes Glassmorphism-Design mit fließenden Übergängen. Das ausgereifte Erscheinungsbild soll dafür sorgen, dass die App angenehm zu bedienen ist.

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

**Was geschieht hier?**

**Glassmorphism-Design:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` – violetter Farbverlauf im Hintergrund
- `backdrop-filter: blur(10px)` – erzeugt den Milchglaseffekt der Karte
- `rgba(255, 255, 255, 0.95)` – halbtransparentes Weiß für den Glaseffekt

**Benutzerdefinierte Gestaltung der Kontrollkästchen:**

- `appearance: none` entfernt das standardmäßige Kontrollkästchen des Browsers
- Mit `::after` erstellen wir ein eigenes abgerundetes Quadrat mit Häkchen
- Bei `checked` erscheint das Häkchen als Unicode-Zeichen ✓

**Hover-Interaktionen:**

- TODO-Einträge gleiten beim Daraufzeigen nach rechts (`transform: translateX(4px)`)
- Die Löschschaltfläche bleibt bis zum Daraufzeigen ausgeblendet (`opacity: 0` → `opacity: 1`)
- Schaltflächen werden beim Daraufzeigen für eine spürbare Rückmeldung leicht vergrößert

**Leerzustand:**

- `#todo-list:empty::before` zeigt eine Meldung an, wenn keine TODO-Einträge vorhanden sind
- Reine CSS-Lösung – kein JavaScript erforderlich

### App ausführen
Sehen wir uns die App in Aktion an! Starten Sie den Entwicklungsserver:

```bash
wails3 dev
```

Die App wird kompiliert und geöffnet. Probieren Sie sie aus:

- Geben Sie einen TODO-Eintrag ein und drücken Sie die Eingabetaste oder klicken Sie auf „Hinzufügen“
- Klicken Sie auf das Kontrollkästchen, um den Eintrag als erledigt zu markieren
- Zeigen Sie auf einen TODO-Eintrag, damit die Löschschaltfläche erscheint
- Beachten Sie, wie die Benutzeroberfläche sofort aktualisiert wird – hier arbeitet unser Aktualisierungsmuster

**Was geschieht hier?**

- Wails hat automatisch Bindings für die Methoden Ihres TodoService generiert
- Der Entwicklungsmodus unterstützt Hot Reload – ändern Sie das CSS und beobachten Sie die Aktualisierung
- Ihr Go-Code wird nativ ausgeführt – eine Übersetzung oder Interpretation ist nicht erforderlich

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

**Warum ist das wichtig?**

- Mehrere Aufrufe aus dem Frontend können gleichzeitig erfolgen
- Leseoperationen blockieren einander nicht
- Schreiboperationen erhalten exklusiven Zugriff
- `defer` stellt sicher, dass Sperren immer freigegeben werden

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

Im Frontend können Sie diese abfangen:

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

**Alternative Vorgehensweise:** Geben Sie aus jeder Methode die aktualisierte Liste zurück, um den zweiten Aufruf zu vermeiden.

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

Fügen Sie für Produktions-Apps Persistenz mit der Speicherlösung hinzu, die am besten zu Ihrer Anwendung passt, beispielsweise SQLite, PostgreSQL oder ein gehosteter Dienst.

## Für den Produktionseinsatz bauen

Wenn Sie Ihre TODO-App verteilen möchten, bauen Sie sie für den Produktionseinsatz:

```bash
wails3 build
```

Dadurch wird eine optimierte native ausführbare Datei in `bin/` erstellt:

- Kompiliert Ihren Go-Code mit Optimierungen
- Baut Ihr Frontend für den Produktionseinsatz
- Bündelt alles in einer einzigen ausführbaren Datei
- Die resultierende App ist üblicherweise 10-20 MB groß (im Vergleich zu 150 MB oder mehr bei Electron)

Sie können die ausführbare Datei direkt starten – es ist keine Laufzeitumgebung erforderlich und es müssen keine Server gestartet werden. Es handelt sich um eine echte native Anwendung.

## Was Sie erstellt haben

Sie haben gerade eine vollständige TODO-Anwendung mit folgenden Funktionen erstellt:

**Vollständige CRUD-Implementierung:**

- Einen Dienst mit Operationen zum Erstellen, Lesen, Aktualisieren und Löschen erstellt
- Eingabevalidierung und Fehlerbehandlung hinzugefügt
- Gelernt, wie Go-Fehler zu JavaScript-Ausnahmen werden

**Threadsichere Zustandsverwaltung:**

- `sync.RWMutex` verwendet, um gleichzeitige Zugriffe sicher zu handhaben
- Den Unterschied zwischen Lesesperren (RLock) und Schreibsperren (Lock) verstanden
- Gesehen, wie `defer` durch garantiertes Aufräumen Deadlocks verhindert

**Moderne, ausgereifte Benutzeroberfläche:**

- Eine Glassmorphism-Oberfläche mit Farbverläufen und Weichzeichnungseffekten erstellt
- Individuell gestaltete Kontrollkästchen ohne Framework erstellt
- Hover-Interaktionen und Übergänge für ein natives Bediengefühl hinzugefügt
- Einen leeren Zustand ausschließlich mit CSS implementiert

**Wails-Grundlagen:**

- Dienstregistrierung und automatische Generierung von Bindings
- Aufrufen von Go-Methoden aus JavaScript mit async/await
- Synchronisierung des Zustands zwischen Go und dem Frontend
- Erstellen und Paketieren einer nativen Desktop-App

## Nächste Schritte

Nachdem Sie nun CRUD-Operationen und Zustandsverwaltung kennen, können Sie Folgendes ausprobieren:

- **Persistenz hinzufügen:** Sorgen Sie mit SQLite oder einer anderen Speicherlösung dafür, dass die Aufgaben auch nach einem Neustart der App erhalten bleiben.
- **Weitere Funktionen hinzufügen:** Filterung (alle/aktive/erledigte Aufgaben), Bearbeitung vorhandener Aufgaben und Massenoperationen
- **Das Notes-Tutorial erkunden:** Erfahren Sie unter [Notes](/tutorials/03-notes-vanilla/), wie Dateioperationen funktionieren.
- **Etwas Eigenes entwickeln:** Nutzen Sie diese Konzepte, um Ihre eigene App zu entwickeln!
