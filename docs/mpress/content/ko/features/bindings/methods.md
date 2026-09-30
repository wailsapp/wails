---
title: "메서드 바인딩"
description: "타입 안전성을 유지하면서 JavaScript에서 Go 메서드 호출"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## 타입 안전 Go-JavaScript 바인딩

Wails는 Go 메서드에 사용할 **타입 안전 JavaScript/TypeScript 바인딩을 자동으로 생성합니다**. Go 코드를 작성하고 명령 하나만 실행하면 HTTP 오버헤드, 수작업, 상용구 코드 없이 타입이 완전히 지정된 프런트엔드 함수를 얻을 수 있습니다.

## 빠른 시작

**1. Go 서비스를 작성합니다.**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. 서비스를 등록합니다.**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. 바인딩을 생성합니다.**

```bash
wails3 generate bindings
```

**4. JavaScript에서 사용합니다.**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**이것으로 끝입니다!** 타입 안전 Go-JavaScript 호출을 사용할 수 있습니다.

## 서비스 생성

### 기본 서비스

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

**등록:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**핵심 사항:**

- **내보낸 메서드**(PascalCase)만 바인딩됩니다
- 메서드는 값 또는 `(value, error)`을 반환할 수 있습니다
- 서비스는 <strong>싱글턴</strong>입니다(애플리케이션당 인스턴스 하나)

### 상태가 있는 서비스

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

**중요:** 서비스는 모든 창에서 공유됩니다. 스레드 안전성을 위해 뮤텍스를 사용하세요.

### 종속성이 있는 서비스

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

**종속성과 함께 등록:**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## 바인딩 생성

### 기본 생성

```bash
wails3 generate bindings
```

**출력:**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**생성된 구조:**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### TypeScript 생성

```bash
wails3 generate bindings -ts
```

완전한 TypeScript 타입이 포함된 <strong>`.ts` 파일</strong>을 생성합니다.

### 사용자 지정 출력 디렉터리

```bash
wails3 generate bindings -d ./src/bindings
```

### 감시 모드(개발)

```bash
wails3 dev
```

Go 코드가 변경되면 **바인딩을 자동으로 다시 생성합니다**.

## 바인딩 사용

### JavaScript

**생성된 바인딩:**

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

**사용법:**

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

**생성된 바인딩:**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**사용법:**

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

**장점:**

- 완전한 타입 검사
- IDE 자동 완성
- 컴파일 시간 오류
- 더 나은 리팩터링

### 인덱스 파일

**생성된 인덱스:**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**간소화된 가져오기:**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## 타입 매핑

### 기본 타입

| Go 타입 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 복합 타입

| Go 타입 | JavaScript/TypeScript | 참고 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 문자열 키 맵 |
| `map[K]V` | `{ [_ in K]?: V }` | 문자열이 아닌 `K`은 JS `Map`이 **아니라** 매핑된 타입으로 렌더링됨 |
| `[]byte` | `string` | base64로 인코딩됨 |
| `struct` | `class` / `interface` | 필드 포함 |
| `time.Time` | `any` | 런타임에서 RFC3339Nano 문자열로 직렬화됨 |
| `*T` | `T \| null` | 포인터는 null 허용을 의미함 |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 반환 값이면 Exception, 아니면 any |

### 지원되지 않는 타입

다음 타입은 브리지를 통해 전달할 수 **없습니다**:

- `chan T`(채널)
- `func()`(함수)
- 복잡한 인터페이스(`interface{}` 제외)
- 내보내지 않은 필드(소문자로 시작)

**해결 방법:** ID 또는 핸들을 사용하세요:

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

## 오류 처리

### Go 측

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

### JavaScript 측

바인딩된 메서드가 실패하면 반환된 Promise가 JavaScript `Error` 객체와 함께 거부됩니다:

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

런타임은 발생한 문제에 따라 서로 다른 오류 타입을 throw합니다:

| 오류 타입 | 발생 조건 |
| --- | --- |
| `TypeError` | 호출에 전달된 인수의 개수가 잘못되었거나 인수를 Go 타입으로 변환할 수 없는 경우 |
| `RuntimeError` | 메서드가 오류를 반환했거나 실행 중 패닉이 발생한 경우 |
| `Error` | 존재하지 않는 메서드를 호출한 경우 등 그 밖의 오류가 발생한 경우 |

모든 오류는 다음 정보를 제공합니다:

- `name`: 위 표에 나온 오류 타입
- `message`: Go 오류 메시지
- `cause`: 가능한 경우 JSON으로 직렬화한 Go 오류입니다. 메서드가 여러 오류를 반환한 경우 `cause`은 오류별 항목이 하나씩 포함된 배열입니다.

`RuntimeError` 클래스는 `@wailsio/runtime` 패키지에서 내보내므로, Go 코드가 반환한 오류를 다른 실패와 구분할 수 있습니다:

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
이전 버전의 런타임에서는 모든 실패한 호출이 내부 오류의 원시 JSON을 메시지에 포함한 일반 `Error`과 함께 거부되었습니다.

@end

### 구조화된 오류 데이터

Go에서 사용자 정의 오류 타입을 반환하면 해당 JSON 형식을 throw된 오류의 `cause` 속성에서 사용할 수 있습니다:

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

오류는 표준 `encoding/json` 패키지로 직렬화되므로 내보낸 필드만 포함됩니다. `errors.New` 또는 `fmt.Errorf`로 생성한 오류에는 내보낸 필드가 없으므로 빈 객체로 직렬화됩니다.

오류 직렬화 방식을 완전히 제어하려면 서비스 옵션에 `MarshalError` 함수를 제공하세요:

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

`MarshalError`은 유효한 JSON을 반환해야 하며, 기본 직렬화 방식을 사용하려면 `nil`을 반환해야 합니다.

## 성능

### 호출 오버헤드

**일반적인 호출:** &lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**다른 방식과 비교:**

- HTTP/REST: 5-50ms
- IPC: 1-10ms
- Wails: &lt;1ms

### 최적화 팁

**✅ 작업을 일괄 처리하세요:**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ 결과를 캐시하세요:**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ 스트리밍에는 이벤트를 사용하세요:**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## 전체 예제

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

## 권장 사례

### ✅ 권장 사항

- **메서드를 단순하게 유지하세요** - 단일 책임 원칙을 따르세요
- **오류를 반환하세요** - 패닉을 일으키지 마세요
- **스레드 안전 상태를 사용하세요** - 공유 데이터에는 뮤텍스를 사용하세요
- **작업을 일괄 처리하세요** - 브리지 호출을 줄이세요
- **Go 측에서 캐시하세요** - 반복 작업을 피하세요
- **메서드를 문서화하세요** - 주석이 JSDoc으로 변환됩니다

### ❌ 금지 사항

- **블로킹하지 마세요** - 오래 걸리는 작업에는 고루틴을 사용하세요
- **채널을 반환하지 마세요** - 대신 이벤트를 사용하세요
- **함수를 반환하지 마세요** - 지원되지 않습니다
- **오류를 무시하지 마세요** - 항상 처리하세요
- **내보내지 않은 필드를 사용하지 마세요** - 바인딩되지 않습니다

## 다음 단계

@cards{cols="2"}
◆ 서비스
서비스 시스템을 자세히 알아보세요.

[자세히 알아보기 →](/features/bindings/services/)

---
📖 모델
복잡한 데이터 구조를 바인딩하세요.

[자세히 알아보기 →](/features/bindings/models/)

---
🚀 Go-프런트엔드 브리지
브리지 메커니즘을 이해하세요.

[자세히 알아보기 →](/concepts/bridge/)

---
★ 이벤트
게시/구독 통신에 이벤트를 사용하세요.

[자세히 알아보기 →](/features/events/system/)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [바인딩 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)를 확인하세요.
