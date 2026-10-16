---
title: "Vinculações de métodos"
description: "Chame métodos Go pelo JavaScript com segurança de tipos"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## Vinculações Go–JavaScript com segurança de tipos

O Wails **gera automaticamente vinculações JavaScript/TypeScript com segurança de tipos** para seus métodos Go. Escreva o código Go, execute um comando e obtenha funções de frontend totalmente tipadas, sem a sobrecarga de HTTP, sem trabalho manual e sem código repetitivo.

## Início rápido

**1. Escreva o serviço Go:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. Registre o serviço:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. Gere as vinculações:**

```bash
wails3 generate bindings
```

**4. Use no JavaScript:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**Pronto!** Chamadas do Go para o JavaScript com segurança de tipos.

## Criação de serviços

### Serviço básico

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

**Registre:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Pontos principais:**

- Somente métodos **exportados** (PascalCase) são vinculados
- Os métodos podem retornar valores ou `(value, error)`
- Os serviços são **singletons** (uma instância por aplicativo)

### Serviço com estado

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

**Importante:** os serviços são compartilhados entre todas as janelas. Use mutexes para garantir a segurança entre threads.

### Serviço com dependências

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

**Registre com as dependências:**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## Geração de vinculações

### Geração básica

```bash
wails3 generate bindings
```

**Saída:**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**Estrutura gerada:**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### Geração para TypeScript

```bash
wails3 generate bindings -ts
```

**Gera arquivos `.ts`** com tipagem TypeScript completa.

### Diretório de saída personalizado

```bash
wails3 generate bindings -d ./src/bindings
```

### Modo de observação (desenvolvimento)

```bash
wails3 dev
```

**Regenera automaticamente as vinculações** quando o código Go é alterado.

## Uso das vinculações

### JavaScript

**Vinculação gerada:**

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

**Uso:**

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

**Vinculação gerada:**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**Uso:**

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

**Benefícios:**

- Verificação completa de tipos
- Preenchimento automático no IDE
- Erros em tempo de compilação
- Refatoração aprimorada

### Arquivos de índice

**Índice gerado:**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**Importações simplificadas:**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## Mapeamento de tipos

### Tipos primitivos

| Tipo Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Tipos complexos

| Tipo Go | JavaScript/TypeScript | Observações |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | map com chaves do tipo string |
| `map[K]V` | `{ [_ in K]?: V }` | um `K` que não seja string é renderizado como um tipo mapeado, e **não** como um `Map` do JS |
| `[]byte` | `string` | codificado em base64 |
| `struct` | `class` / `interface` | com campos |
| `time.Time` | `any` | serializado no runtime como uma string RFC3339Nano |
| `*T` | `T \| null` | o ponteiro indica que o valor pode ser nulo |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Exception se usado como valor de retorno; caso contrário, any |

### Tipos não compatíveis

Estes tipos **não podem** ser passados pela ponte:

- `chan T` (canais)
- `func()` (funções)
- Interfaces complexas (exceto `interface{}`)
- Campos não exportados (com inicial minúscula)

**Solução alternativa:** use IDs ou identificadores:

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

## Tratamento de erros

### No Go

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

### No JavaScript

Quando um método vinculado falha, a promise retornada é rejeitada com um objeto JavaScript `Error`:

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

O runtime lança um tipo de erro diferente conforme o que deu errado:

| Tipo de erro | Lançado quando |
| --- | --- |
| `TypeError` | A chamada tem um número incorreto de argumentos ou um argumento não pode ser convertido para o tipo Go |
| `RuntimeError` | O método retornou um erro ou entrou em pânico durante a execução |
| `Error` | Qualquer outra falha, como uma chamada a um método inexistente |

Cada erro fornece:

- `name`: o tipo de erro da tabela acima
- `message`: a mensagem do erro do Go
- `cause`: o erro do Go serializado como JSON, quando disponível. Se o método retornar vários erros, `cause` será um array com uma entrada para cada erro.

A classe `RuntimeError` é exportada pelo pacote `@wailsio/runtime`, permitindo distinguir os erros retornados pelo seu código Go de outras falhas:

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
Em versões anteriores do runtime, toda chamada com falha era rejeitada com um `Error` simples, cuja mensagem continha o JSON bruto do erro subjacente.

@end

### Dados de erro estruturados

Ao retornar um tipo de erro personalizado do Go, sua representação em JSON fica disponível na propriedade `cause` do erro lançado:

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

Os erros são serializados com o pacote padrão `encoding/json`, portanto apenas os campos exportados são incluídos. Erros criados com `errors.New` ou `fmt.Errorf` não têm campos exportados e são serializados como um objeto vazio.

Para ter controle total sobre como os erros são serializados, forneça uma função `MarshalError` nas opções do serviço:

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

`MarshalError` deve retornar um JSON válido ou `nil` para usar a serialização padrão como alternativa.

## Desempenho

### Sobrecarga das chamadas

**Chamada típica:** &lt;1 ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**Em comparação com as alternativas:**

- HTTP/REST: 5-50 ms
- IPC: 1-10 ms
- Wails: &lt;1 ms

### Dicas de otimização

**✅ Agrupe operações:**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ Armazene resultados em cache:**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ Use eventos para streaming:**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## Exemplo completo

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

## Práticas recomendadas

### ✅ Faça

- **Mantenha os métodos simples** — Responsabilidade única
- **Retorne erros** — Não cause panic
- **Use estado seguro para concorrência** — Use mutexes para dados compartilhados
- **Agrupe operações** — Reduza as chamadas pela ponte
- **Use cache no lado do Go** — Evite trabalho repetido
- **Documente os métodos** — Os comentários se tornam JSDoc

### ❌ Não faça

- **Não bloqueie** — Use goroutines para operações demoradas
- **Não retorne canais** — Use eventos em vez disso
- **Não retorne funções** — Não há suporte
- **Não ignore erros** — Sempre trate-os
- **Não use campos não exportados** — Eles não serão vinculados

## Próximas etapas

@cards{cols="2"}
◆ Serviços
Explore em detalhes o sistema de serviços.

[Saiba mais →](/features/bindings/services/)

---
📖 Modelos
Vincule estruturas de dados complexas.

[Saiba mais →](/features/bindings/models/)

---
🚀 Ponte Go-frontend
Entenda o mecanismo da ponte.

[Saiba mais →](/concepts/bridge/)

---
★ Eventos
Use eventos para comunicação de publicação/assinatura.

[Saiba mais →](/features/events/system/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de bindings](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
