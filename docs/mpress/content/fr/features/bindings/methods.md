---
title: "Liaisons de méthodes"
description: "Appelez des méthodes Go depuis JavaScript avec un typage sûr"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## Liaisons Go-JavaScript à typage sûr

Wails **génère automatiquement des liaisons JavaScript/TypeScript à typage sûr** pour vos méthodes Go. Écrivez du code Go, exécutez une commande et obtenez des fonctions frontend entièrement typées, sans surcoût HTTP, sans intervention manuelle et sans aucun code répétitif.

## Démarrage rapide

**1. Écrivez le service Go :**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. Enregistrez le service :**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. Générez les liaisons :**

```bash
wails3 generate bindings
```

**4. Utilisez-les dans JavaScript :**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**C’est tout !** Vous disposez maintenant d’appels de Go vers JavaScript à typage sûr.

## Création de services

### Service de base

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

**Enregistrez-le :**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Points clés :**

- Seules les **méthodes exportées** (PascalCase) sont liées
- Les méthodes peuvent renvoyer des valeurs ou `(value, error)`
- Les services sont des **singletons** (une instance par application)

### Service avec état

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

**Important :** les services sont partagés entre toutes les fenêtres. Utilisez des mutex pour garantir la sûreté des accès concurrents.

### Service avec dépendances

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

**Enregistrez-le avec ses dépendances :**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## Génération des liaisons

### Génération de base

```bash
wails3 generate bindings
```

**Sortie :**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**Structure générée :**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### Génération TypeScript

```bash
wails3 generate bindings -ts
```

**Génère des fichiers `.ts`** avec tous les types TypeScript.

### Répertoire de sortie personnalisé

```bash
wails3 generate bindings -d ./src/bindings
```

### Mode surveillance (développement)

```bash
wails3 dev
```

**Régénère automatiquement les liaisons** lorsque le code Go change.

## Utilisation des liaisons

### JavaScript

**Liaison générée :**

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

**Utilisation :**

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

**Liaison générée :**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**Utilisation :**

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

**Avantages :**

- Vérification complète des types
- Autocomplétion dans l’IDE
- Erreurs détectées à la compilation
- Refactorisation facilitée

### Fichiers d’index

**Index généré :**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**Imports simplifiés :**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## Correspondance des types

### Types primitifs

| Type Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Types complexes

| Type Go | JavaScript/TypeScript | Remarques |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | table associative à clés de type chaîne |
| `map[K]V` | `{ [_ in K]?: V }` | un `K` qui n’est pas de type chaîne est représenté par un type mappé, et **non** par un `Map` JS |
| `[]byte` | `string` | encodé en base64 |
| `struct` | `class` / `interface` | avec des champs |
| `time.Time` | `any` | sérialisé à l’exécution sous forme de chaîne RFC3339Nano |
| `*T` | `T \| null` | un pointeur indique une valeur nullable |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Exception si utilisé comme valeur de retour, sinon any |

### Types non pris en charge

Ces types **ne peuvent pas** être transmis via le pont :

- `chan T` (canaux)
- `func()` (fonctions)
- Interfaces complexes (sauf `interface{}`)
- Champs non exportés (en minuscules)

**Solution de contournement :** utilisez des identifiants ou des handles :

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

## Gestion des erreurs

### Côté Go

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

### Côté JavaScript

Lorsqu’une méthode liée échoue, la promesse renvoyée est rejetée avec un objet JavaScript `Error` :

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

Le runtime lève un type d’erreur différent selon la nature du problème :

| Type d’erreur | Levée lorsque |
| --- | --- |
| `TypeError` | L’appel comporte un nombre incorrect d’arguments, ou un argument ne peut pas être converti vers le type Go |
| `RuntimeError` | La méthode a renvoyé une erreur ou a déclenché une panique pendant son exécution |
| `Error` | Tout autre échec, par exemple un appel à une méthode inexistante |

Chaque erreur fournit :

- `name` : le type d’erreur indiqué dans le tableau ci-dessus
- `message` : le message de l’erreur Go
- `cause` : l’erreur Go sérialisée en JSON, lorsqu’elle est disponible. Si la méthode a renvoyé plusieurs erreurs, `cause` est un tableau contenant une entrée par erreur.

La classe `RuntimeError` est exportée par le package `@wailsio/runtime`, ce qui vous permet de distinguer les erreurs renvoyées par votre code Go des autres échecs :

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
Dans les versions antérieures du runtime, chaque appel ayant échoué rejetait la promesse avec un simple `Error` dont le message contenait le JSON brut de l’erreur sous-jacente.

@end

### Données d’erreur structurées

Lorsque Go renvoie un type d’erreur personnalisé, sa représentation JSON est disponible dans la propriété `cause` de l’erreur levée :

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

Les erreurs sont sérialisées avec le package standard `encoding/json` ; seuls les champs exportés sont donc inclus. Les erreurs créées avec `errors.New` ou `fmt.Errorf` ne comportent aucun champ exporté et sont sérialisées sous forme d’objet vide.

Pour contrôler entièrement la sérialisation des erreurs, fournissez une fonction `MarshalError` dans les options du service :

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

`MarshalError` doit renvoyer du JSON valide, ou `nil` pour revenir à la sérialisation par défaut.

## Performances

### Surcoût des appels

**Appel standard :** &lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**Comparaison avec les autres solutions :**

- HTTP/REST : 5-50ms
- IPC : 1-10ms
- Wails : &lt;1ms

### Conseils d’optimisation

**✅ Regroupez les opérations :**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ Mettez les résultats en cache :**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ Utilisez des événements pour la diffusion en continu :**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## Exemple complet

**Go :**

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

**JavaScript :**

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

## Bonnes pratiques

### ✅ À faire

- **Gardez les méthodes simples** — Une seule responsabilité
- **Renvoyez les erreurs** — Ne déclenchez pas de panique
- **Utilisez un état sûr pour les accès concurrents** — Protégez les données partagées avec des mutex
- **Regroupez les opérations** — Réduisez les appels à la passerelle
- **Mettez les données en cache côté Go** — Évitez de répéter le même travail
- **Documentez les méthodes** — Les commentaires deviennent de la JSDoc

### ❌ À ne pas faire

- **Ne bloquez pas l’exécution** — Utilisez des goroutines pour les opérations longues
- **Ne renvoyez pas de canaux** — Utilisez plutôt des événements
- **Ne renvoyez pas de fonctions** — Ce n’est pas pris en charge
- **N’ignorez pas les erreurs** — Gérez-les systématiquement
- **N’utilisez pas de champs non exportés** — Aucune liaison ne sera générée pour eux

## Étapes suivantes

@cards{cols="2"}
◆ Services
Explorez en détail le système de services.

[En savoir plus →](/features/bindings/services/)

---
📖 Modèles
Créez des liaisons pour des structures de données complexes.

[En savoir plus →](/features/bindings/models/)

---
🚀 Passerelle Go-frontend
Comprenez le mécanisme de la passerelle.

[En savoir plus →](/concepts/bridge/)

---
★ Événements
Utilisez des événements pour les communications par publication/abonnement.

[En savoir plus →](/features/events/system/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de liaisons](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
