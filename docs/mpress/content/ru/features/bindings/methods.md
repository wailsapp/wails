---
title: "Привязки методов"
description: "Типобезопасный вызов методов Go из JavaScript"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## Типобезопасные привязки Go и JavaScript

Wails **автоматически создаёт типобезопасные привязки JavaScript/TypeScript** для ваших методов Go. Напишите код Go, выполните одну команду и получите полностью типизированные функции фронтенда без накладных расходов HTTP, ручной работы и шаблонного кода.

## Быстрый старт

**1. Напишите сервис Go:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. Зарегистрируйте сервис:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. Создайте привязки:**

```bash
wails3 generate bindings
```

**4. Используйте в JavaScript:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**Готово!** Типобезопасные вызовы из Go в JavaScript.

## Создание сервисов

### Базовый сервис

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

**Зарегистрируйте:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Основные моменты:**

- Привязки создаются только для **экспортируемых методов** (PascalCase)
- Методы могут возвращать значения или `(value, error)`
- Сервисы являются **одиночками** (по одному экземпляру на приложение)

### Сервис с состоянием

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

**Важно:** сервисы совместно используются всеми окнами. Для потокобезопасности применяйте мьютексы.

### Сервис с зависимостями

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

**Зарегистрируйте с зависимостями:**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## Создание привязок

### Базовое создание

```bash
wails3 generate bindings
```

**Результат:**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**Созданная структура:**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### Создание привязок TypeScript

```bash
wails3 generate bindings -ts
```

**Создаёт файлы `.ts`** с полными типами TypeScript.

### Пользовательский каталог вывода

```bash
wails3 generate bindings -d ./src/bindings
```

### Режим наблюдения (разработка)

```bash
wails3 dev
```

**Автоматически пересоздаёт привязки** при изменении кода Go.

## Использование привязок

### JavaScript

**Созданная привязка:**

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

**Использование:**

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

**Созданная привязка:**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**Использование:**

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

**Преимущества:**

- Полная проверка типов
- Автодополнение в IDE
- Ошибки на этапе компиляции
- Более удобный рефакторинг

### Индексные файлы

**Созданный индекс:**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**Упрощённый импорт:**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## Сопоставление типов

### Примитивные типы

| Тип Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Составные типы

| Тип Go | JavaScript/TypeScript | Примечания |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | map со строковыми ключами |
| `map[K]V` | `{ [_ in K]?: V }` | Если `K` — нестроковый тип, используется сопоставленный тип, а **не** `Map` JavaScript. |
| `[]byte` | `string` | в кодировке base64 |
| `struct` | `class` / `interface` | с полями |
| `time.Time` | `any` | во время выполнения сериализуется в строку RFC3339Nano |
| `*T` | `T \| null` | указатель означает, что допустимо значение null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Исключение при использовании в качестве возвращаемого значения, иначе — any |

### Неподдерживаемые типы

Эти типы **нельзя** передавать через мост:

- `chan T` (каналы)
- `func()` (функции)
- Сложные интерфейсы (кроме `interface{}`)
- Неэкспортируемые поля (со строчной буквы)

**Обходное решение:** используйте идентификаторы или дескрипторы:

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

## Обработка ошибок

### На стороне Go

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

### На стороне JavaScript

Если связанный метод завершается с ошибкой, возвращённый промис отклоняется с объектом JavaScript `Error`:

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

Среда выполнения выбрасывает ошибку разного типа в зависимости от причины сбоя:

| Тип ошибки | Когда выбрасывается |
| --- | --- |
| `TypeError` | В вызове указано неверное количество аргументов или аргумент невозможно преобразовать в тип Go |
| `RuntimeError` | Метод вернул ошибку или во время выполнения произошла паника |
| `Error` | Любой другой сбой, например вызов несуществующего метода |

Каждая ошибка содержит:

- `name`: тип ошибки из приведённой выше таблицы
- `message`: сообщение об ошибке Go
- `cause`: ошибка Go, сериализованная в JSON, если она доступна. Если метод вернул несколько ошибок, `cause` представляет собой массив с одним элементом для каждой ошибки.

Класс `RuntimeError` экспортируется пакетом `@wailsio/runtime`, поэтому ошибки, возвращённые вашим кодом Go, можно отличать от других сбоев:

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
В предыдущих версиях среды выполнения каждый неудачный вызов приводил к отклонению с обычным объектом `Error`, сообщение которого содержало необработанный JSON исходной ошибки.

@end

### Структурированные данные об ошибке

Если вернуть из Go ошибку пользовательского типа, её представление в формате JSON будет доступно в свойстве `cause` выброшенной ошибки:

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

Ошибки сериализуются стандартным пакетом `encoding/json`, поэтому включаются только экспортируемые поля. Ошибки, созданные с помощью `errors.New` или `fmt.Errorf`, не имеют экспортируемых полей и сериализуются в пустой объект.

Чтобы полностью контролировать сериализацию ошибок, укажите функцию `MarshalError` в параметрах сервиса:

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

`MarshalError` должна возвращать допустимый JSON или `nil`, чтобы использовать стандартную сериализацию.

## Производительность

### Накладные расходы вызова

**Типичный вызов:** &lt;1 мс

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**Сравнение с альтернативами:**

- HTTP/REST: 5-50 мс
- IPC: 1-10 мс
- Wails: &lt;1 мс

### Советы по оптимизации

**✅ Объединяйте операции в пакеты:**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ Кешируйте результаты:**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ Используйте события для потоковой передачи:**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## Полный пример

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

## Рекомендации

### ✅ Что следует делать

- **Делайте методы простыми** — каждый метод должен решать одну задачу
- **Возвращайте ошибки** — не вызывайте панику
- **Обеспечивайте потокобезопасность состояния** — защищайте общие данные мьютексами
- **Объединяйте операции в пакеты** — сократите количество вызовов через мост
- **Кешируйте на стороне Go** — избегайте повторного выполнения одной и той же работы
- **Документируйте методы** — комментарии преобразуются в JSDoc

### ❌ Чего не следует делать

- **Не блокируйте выполнение** — используйте горутины для длительных операций
- **Не возвращайте каналы** — вместо них используйте события
- **Не возвращайте функции** — это не поддерживается
- **Не игнорируйте ошибки** — всегда обрабатывайте их
- **Не используйте неэкспортируемые поля** — для них не будут созданы привязки

## Дальнейшие шаги

@cards{cols="2"}
◆ Сервисы
Подробно изучите систему сервисов.

[Подробнее →](/features/bindings/services/)

---
📖 Модели
Создавайте привязки для сложных структур данных.

[Подробнее →](/features/bindings/models/)

---
🚀 Мост между Go и фронтендом
Разберитесь в механизме работы моста.

[Подробнее →](/concepts/bridge/)

---
★ События
Используйте события для взаимодействия по модели публикации и подписки.

[Подробнее →](/features/events/system/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами привязок](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
