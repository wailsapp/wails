---
title: "Methoden-Bindings"
description: "Go-Methoden typsicher aus JavaScript aufrufen"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## Typsichere Go-JavaScript-Bindings

Wails **generiert automatisch typsichere JavaScript-/TypeScript-Bindings** für Ihre Go-Methoden. Schreiben Sie Go-Code, führen Sie einen einzigen Befehl aus und erhalten Sie vollständig typisierte Frontend-Funktionen – ohne HTTP-Overhead, manuelle Arbeit oder Boilerplate-Code.

## Schnellstart

**1. Go-Service schreiben:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. Service registrieren:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. Bindings generieren:**

```bash
wails3 generate bindings
```

**4. In JavaScript verwenden:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**Das ist alles!** Typsichere Aufrufe von Go nach JavaScript.

## Services erstellen

### Einfacher Service

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}

func (c *CalculatorService) Multiply(a, b int) int {
    return a * b
}

func (c *CalculatorService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

**Registrieren:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Wichtige Punkte:**

- Nur **exportierte Methoden** (PascalCase) werden gebunden
- Methoden können Werte oder `(value, error)` zurückgeben
- Services sind **Singletons** (eine Instanz pro Anwendung)

### Service mit Zustand

```go
type CounterService struct {
    count int
    mu    sync.Mutex
}

func (c *CounterService) Increment() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    return c.count
}

func (c *CounterService) Decrement() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count--
    return c.count
}

func (c *CounterService) GetCount() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

func (c *CounterService) Reset() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count = 0
}
```

**Wichtig:** Services werden von allen Fenstern gemeinsam genutzt. Verwenden Sie Mutexe, um Threadsicherheit zu gewährleisten.

### Service mit Abhängigkeiten

```go
type DatabaseService struct {
    db *sql.DB
}

func NewDatabaseService(db *sql.DB) *DatabaseService {
    return &DatabaseService{db: db}
}

func (d *DatabaseService) GetUser(id int) (*User, error) {
    var user User
    err := d.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    return &user, err
}
```

**Mit Abhängigkeiten registrieren:**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## Bindings generieren

### Einfache Generierung

```bash
wails3 generate bindings
```

**Ausgabe:**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**Generierte Struktur:**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### TypeScript-Generierung

```bash
wails3 generate bindings -ts
```

**Generiert `.ts`-Dateien** mit vollständigen TypeScript-Typen.

### Benutzerdefiniertes Ausgabeverzeichnis

```bash
wails3 generate bindings -d ./src/bindings
```

### Überwachungsmodus (Entwicklung)

```bash
wails3 dev
```

**Generiert Bindings automatisch neu**, wenn sich der Go-Code ändert.

## Bindings verwenden

### JavaScript

**Generiertes Binding:**

```javascript
// frontend/bindings/<full-go-import-path>/calculatorservice.js
// (Real generated output — imports $Call from /wails/runtime.js and calls $Call.ByID
// with a numeric method ID. Generate with `wails3 generate bindings -names` to get
// $Call.ByName("<package>.<Struct>.<Method>", ...) instead.)
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {number} $0
 * @param {number} $1
 * @returns {Promise<number>}
 */
export function Add($0, $1) {
    return $Call.ByID(1234567890, $0, $1); // numeric ID assigned by the generator
}
```

**Verwendung:**

```javascript
import { Add, Subtract, Multiply, Divide } from './bindings/changeme/calculatorservice'

// Simple calls
const sum = await Add(5, 3)        // 8
const diff = await Subtract(10, 4)  // 6
const product = await Multiply(7, 6) // 42

// Error handling
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Error:", error)  // "division by zero"
}
```

### TypeScript

**Generiertes Binding:**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**Verwendung:**

```typescript
import { Add, Divide } from './bindings/changeme/calculatorservice'

const sum: number = await Add(5, 3)

try {
    const result = await Divide(10, 0)
} catch (error: unknown) {
    if (error instanceof Error) {
        console.error(error.message)
    }
}
```

**Vorteile:**

- Vollständige Typprüfung
- Automatische Vervollständigung in der IDE
- Fehler zur Kompilierzeit
- Einfacheres Refactoring

### Indexdateien

**Generierter Index:**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**Vereinfachte Importe:**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## Typzuordnung

### Primitive Typen

| Go-Typ | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Komplexe Typen

| Go-Typ | JavaScript/TypeScript | Hinweise |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | Map mit Zeichenfolgenschlüsseln |
| `map[K]V` | `{ [_ in K]?: V }` | `K`, die keine Zeichenfolgen sind, werden als gemappter Typ dargestellt und **nicht** als JavaScript-`Map` |
| `[]byte` | `string` | Base64-codiert |
| `struct` | `class` / `interface` | mit Feldern |
| `time.Time` | `any` | wird zur Laufzeit als RFC3339Nano-Zeichenfolge serialisiert |
| `*T` | `T \| null` | Zeiger bedeutet, dass null zulässig ist |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Als Rückgabewert Exception, andernfalls any |

### Nicht unterstützte Typen

Diese Typen können **nicht** über die Bridge übergeben werden:

- `chan T` (Kanäle)
- `func()` (Funktionen)
- Komplexe Schnittstellen (außer `interface{}`)
- Nicht exportierte Felder (kleingeschrieben)

**Problemumgehung:** Verwenden Sie IDs oder Handles:

```go
// ❌ Can't return file handle
func OpenFile(path string) (*os.File, error)

// ✅ Return file ID instead
var files = make(map[string]*os.File)

func OpenFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    id := generateID()
    files[id] = file
    return id, nil
}

func ReadFile(id string) ([]byte, error) {
    file := files[id]
    return io.ReadAll(file)
}

func CloseFile(id string) error {
    file := files[id]
    delete(files, id)
    return file.Close()
}
```

## Fehlerbehandlung

### Go-Seite

```go
func (d *DatabaseService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    var user User
    err := d.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user %d not found", id)
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    return &user, nil
}
```

### JavaScript-Seite

Wenn eine gebundene Methode fehlschlägt, wird das zurückgegebene Promise mit einem JavaScript-`Error`-Objekt abgelehnt:

```javascript
import { GetUser } from './bindings/changeme/databaseservice'

try {
    const user = await GetUser(123)
    console.log("User:", user)
} catch (error) {
    console.error(error.message)
    // "user 123 not found"
}
```

Die Runtime löst abhängig von der Fehlerursache einen anderen Fehlertyp aus:

| Fehlertyp | Ausgelöst, wenn |
| --- | --- |
| `TypeError` | der Aufruf die falsche Anzahl an Argumenten enthält oder ein Argument nicht in den Go-Typ konvertiert werden kann |
| `RuntimeError` | die Methode einen Fehler zurückgegeben hat oder während ihrer Ausführung eine Panic aufgetreten ist |
| `Error` | ein anderer Fehler auftritt, beispielsweise beim Aufruf einer nicht vorhandenen Methode |

Jeder Fehler enthält:

- `name`: den Fehlertyp aus der obigen Tabelle
- `message`: die Meldung des Go-Fehlers
- `cause`: den als JSON serialisierten Go-Fehler, sofern verfügbar. Hat die Methode mehrere Fehler zurückgegeben, ist `cause` ein Array mit einem Eintrag pro Fehler.

Die Klasse `RuntimeError` wird vom Paket `@wailsio/runtime` exportiert. Dadurch können Sie von Ihrem Go-Code zurückgegebene Fehler von anderen Fehlern unterscheiden:

```javascript
import { Call } from '@wailsio/runtime'
import { GetUser } from './bindings/changeme/databaseservice'

try {
    const user = await GetUser(123)
} catch (error) {
    if (error instanceof Call.RuntimeError) {
        // GetUser returned an error
    } else {
        // The call itself failed
    }
}
```

@note{type="info"}
In früheren Versionen der Runtime wurde jeder fehlgeschlagene Aufruf mit einem einfachen `Error` abgelehnt, dessen Meldung das unverarbeitete JSON des zugrunde liegenden Fehlers enthielt.

@end

### Strukturierte Fehlerdaten

Wenn Sie aus Go einen benutzerdefinierten Fehlertyp zurückgeben, steht dessen JSON-Darstellung in der Eigenschaft `cause` des ausgelösten Fehlers zur Verfügung:

```go
type ValidationError struct {
    Field  string `json:"field"`
    Reason string `json:"reason"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func (s *UserService) UpdateEmail(id int, email string) error {
    if !strings.Contains(email, "@") {
        return &ValidationError{Field: "email", Reason: "invalid email address"}
    }
    // ...
    return nil
}
```

```javascript
import { UpdateEmail } from './bindings/changeme/userservice'

try {
    await UpdateEmail(1, "not-an-email")
} catch (error) {
    console.log(error.message)  // "email: invalid email address"
    console.log(error.cause)    // { field: "email", reason: "invalid email address" }
}
```

Fehler werden mit dem Standardpaket `encoding/json` serialisiert, daher werden nur exportierte Felder einbezogen. Mit `errors.New` oder `fmt.Errorf` erstellte Fehler haben keine exportierten Felder und werden als leeres Objekt serialisiert.

Um vollständig zu steuern, wie Fehler serialisiert werden, geben Sie in den Dienstoptionen eine `MarshalError`-Funktion an:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&UserService{}, application.ServiceOptions{
            MarshalError: func(err error) []byte {
                var validationErr *ValidationError
                if errors.As(err, &validationErr) {
                    data, _ := json.Marshal(map[string]string{
                        "type":  "validation",
                        "field": validationErr.Field,
                    })
                    return data
                }
                return nil // fall back to the default serialisation
            },
        }),
    },
})
```

`MarshalError` muss gültiges JSON zurückgeben oder, um auf die Standardserialisierung zurückzugreifen, `nil`.

## Leistung

### Aufruf-Overhead

**Typischer Aufruf:** &lt;1 ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**Im Vergleich zu Alternativen:**

- HTTP/REST: 5-50 ms
- IPC: 1-10 ms
- Wails: &lt;1 ms

### Optimierungstipps

**✅ Operationen bündeln:**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ Ergebnisse zwischenspeichern:**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ Events für Streaming verwenden:**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## Vollständiges Beispiel

**Go:**

```go
package main

import (
    "fmt"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type TodoService struct {
    todos []Todo
}

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func (t *TodoService) GetAll() []Todo {
    return t.todos
}

func (t *TodoService) Add(title string) Todo {
    todo := Todo{
        ID:        len(t.todos) + 1,
        Title:     title,
        Completed: false,
    }
    t.todos = append(t.todos, todo)
    return todo
}

func (t *TodoService) Toggle(id int) error {
    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos[i].Completed = !t.todos[i].Completed
            return nil
        }
    }
    return fmt.Errorf("todo %d not found", id)
}

func (t *TodoService) Delete(id int) error {
    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos = append(t.todos[:i], t.todos[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("todo %d not found", id)
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&TodoService{}),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**JavaScript:**

```javascript
import { GetAll, Add, Toggle, Delete } from './bindings/changeme/todoservice'

class TodoApp {
    async loadTodos() {
        const todos = await GetAll()
        this.renderTodos(todos)
    }
    
    async addTodo(title) {
        try {
            const todo = await Add(title)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to add todo:", error)
        }
    }
    
    async toggleTodo(id) {
        try {
            await Toggle(id)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to toggle todo:", error)
        }
    }
    
    async deleteTodo(id) {
        try {
            await Delete(id)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to delete todo:", error)
        }
    }
    
    renderTodos(todos) {
        const list = document.getElementById('todo-list')
        list.innerHTML = todos.map(todo => `
            <div class="todo ${todo.Completed ? 'completed' : ''}">
                <input type="checkbox" 
                       ${todo.Completed ? 'checked' : ''}
                       onchange="app.toggleTodo(${todo.ID})">
                <span>${todo.Title}</span>
                <button onclick="app.deleteTodo(${todo.ID})">Delete</button>
            </div>
        `).join('')
    }
}

const app = new TodoApp()
app.loadTodos()
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Methoden einfach halten** – Nur eine Verantwortlichkeit
- **Fehler zurückgeben** – Keine Panik auslösen
- **Threadsicheren Zustand verwenden** – Mutexe für gemeinsam genutzte Daten
- **Operationen bündeln** – Aufrufe über die Bridge reduzieren
- **Auf der Go-Seite zwischenspeichern** – Wiederholte Arbeit vermeiden
- **Methoden dokumentieren** – Kommentare werden zu JSDoc

### ❌ Nicht empfohlen

- **Nicht blockieren** – Goroutinen für lang laufende Operationen verwenden
- **Keine Channels zurückgeben** – Stattdessen Events verwenden
- **Keine Funktionen zurückgeben** – Wird nicht unterstützt
- **Fehler nicht ignorieren** – Immer behandeln
- **Keine nicht exportierten Felder verwenden** – Für sie werden keine Bindings erstellt

## Nächste Schritte

@cards{cols="2"}
◆ Services
Vertiefen Sie Ihr Wissen über das Service-System.

[Mehr erfahren →](/features/bindings/services/)

---
📖 Modelle
Binden Sie komplexe Datenstrukturen ein.

[Mehr erfahren →](/features/bindings/models/)

---
🚀 Go-Frontend-Bridge
Machen Sie sich mit dem Bridge-Mechanismus vertraut.

[Mehr erfahren →](/concepts/bridge/)

---
★ Events
Verwenden Sie Events für die Publish/Subscribe-Kommunikation.

[Mehr erfahren →](/features/events/system/)

@end

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Binding-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/binding) an.
