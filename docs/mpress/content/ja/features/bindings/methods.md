---
title: "メソッドバインディング"
description: "型安全に JavaScript から Go メソッドを呼び出す"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## 型安全な Go-JavaScript バインディング

Wails は、Go メソッド用の<strong>型安全な JavaScript/TypeScript バインディングを自動生成します</strong>。Go コードを記述してコマンドを 1 つ実行するだけで、HTTP のオーバーヘッドも手作業もボイラープレートもなく、完全に型付けされたフロントエンド関数を利用できます。

## クイックスタート

**1. Go サービスを記述します：**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. サービスを登録します：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. バインディングを生成します：**

```bash
wails3 generate bindings
```

**4. JavaScript で使用します：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**これだけです！** 型安全に Go から JavaScript を呼び出せます。

## サービスの作成

### 基本的なサービス

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

**登録：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**要点：**

- **エクスポートされたメソッド**（PascalCase）のみがバインドされます
- メソッドは値または `(value, error)` を返せます
- サービスは<strong>シングルトン</strong>です（アプリケーションごとに 1 インスタンス）

### 状態を持つサービス

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

<strong>重要：</strong>サービスはすべてのウィンドウで共有されます。スレッドセーフにするために mutex を使用してください。

### 依存関係を持つサービス

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

**依存関係とともに登録：**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## バインディングの生成

### 基本的な生成方法

```bash
wails3 generate bindings
```

**出力：**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**生成される構造：**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### TypeScript の生成

```bash
wails3 generate bindings -ts
```

**完全な TypeScript 型を備えた `.ts` ファイルを生成します**。

### 出力先ディレクトリのカスタマイズ

```bash
wails3 generate bindings -d ./src/bindings
```

### 監視モード（開発用）

```bash
wails3 dev
```

Go コードが変更されると、**バインディングを自動的に再生成します**。

## バインディングの使用

### JavaScript

**生成されたバインディング：**

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

**使用方法：**

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

**生成されたバインディング：**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**使用方法：**

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

**利点：**

- 完全な型チェック
- IDE の自動補完
- コンパイル時のエラー検出
- リファクタリングのしやすさ

### インデックスファイル

**生成されたインデックス：**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**簡潔なインポート：**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## 型のマッピング

### プリミティブ型

| Go の型 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`、`int8`、`int16`、`int32`、`int64` | `number` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `number` |
| `float32`、`float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 複合型

| Go の型 | JavaScript/TypeScript | 注記 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 文字列をキーとするマップ |
| `map[K]V` | `{ [_ in K]?: V }` | 文字列以外の`K`は、JS の`Map`では<strong>なく</strong>、マップ型として生成される |
| `[]byte` | `string` | Base64 エンコード済み |
| `struct` | `class` / `interface` | フィールドを持つ |
| `time.Time` | `any` | 実行時に RFC3339Nano 文字列としてシリアライズされる |
| `*T` | `T \| null` | ポインターは null 許容を意味する |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 戻り値の場合は例外、それ以外は any |

### サポートされていない型

次の型はブリッジを介して渡すことが<strong>できません</strong>。

- `chan T`（チャネル）
- `func()`（関数）
- 複雑なインターフェース（`interface{}`を除く）
- 非公開フィールド（小文字で始まるもの）

<strong>回避策：</strong>ID またはハンドルを使用します。

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

## エラー処理

### Go 側

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

### JavaScript 側

バインドされたメソッドが失敗すると、返された Promise は JavaScript の`Error`オブジェクトで拒否されます。

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

ランタイムは、発生した問題に応じて異なるエラー型をスローします。

| エラー型 | スローされる条件 |
| --- | --- |
| `TypeError` | 呼び出しの引数の数が正しくないか、引数を Go の型に変換できない場合 |
| `RuntimeError` | メソッドがエラーを返したか、実行中にパニックが発生した場合 |
| `Error` | 存在しないメソッドの呼び出しなど、その他の障害が発生した場合 |

すべてのエラーには次の情報が含まれます。

- `name`：上の表に記載されたエラー型
- `message`：Go エラーのメッセージ
- `cause`：利用可能な場合は、JSON としてシリアライズされた Go エラー。メソッドが複数のエラーを返した場合、`cause`はエラーごとに 1 つの要素を持つ配列になります。

`RuntimeError`クラスは`@wailsio/runtime`パッケージからエクスポートされるため、Go コードが返したエラーをその他の障害と区別できます。

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
以前のバージョンのランタイムでは、失敗した呼び出しはすべて、基となるエラーの未加工の JSON をメッセージに含む通常の`Error`で拒否されていました。

@end

### 構造化エラーデータ

Go からカスタムエラー型を返すと、スローされたエラーの`cause`プロパティで、その JSON 形式を利用できます。

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

エラーは標準の`encoding/json`パッケージでシリアライズされるため、エクスポートされたフィールドのみが含まれます。`errors.New`または`fmt.Errorf`で作成したエラーにはエクスポートされたフィールドがないため、空のオブジェクトとしてシリアライズされます。

エラーのシリアライズ方法を完全に制御するには、サービスオプションに`MarshalError`関数を指定します。

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

`MarshalError`は有効な JSON を返す必要があります。デフォルトのシリアライズにフォールバックする場合は、`nil`を返す必要があります。

## パフォーマンス

### 呼び出しのオーバーヘッド

**一般的な呼び出し：** &lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**他の方式との比較：**

- HTTP/REST：5-50ms
- IPC：1-10ms
- Wails：&lt;1ms

### 最適化のヒント

**✅ 操作をバッチ処理する：**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ 結果をキャッシュする：**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ ストリーミングにはイベントを使用する：**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## 完全な例

**Go：**

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

**JavaScript：**

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

## ベストプラクティス

### ✅ 推奨事項

- **メソッドをシンプルに保つ** — 単一責任にする
- **エラーを返す** — パニックを発生させない
- **スレッドセーフな状態を使用する** — 共有データにはミューテックスを使用する
- **操作をバッチ処理する** — ブリッジ呼び出しを減らす
- **Go 側でキャッシュする** — 同じ処理の繰り返しを避ける
- **メソッドを文書化する** — コメントは JSDoc に変換される

### ❌ 禁止事項

- **ブロックしない** — 時間のかかる処理には goroutine を使用する
- **チャネルを返さない** — 代わりにイベントを使用する
- **関数を返さない** — サポートされていない
- **エラーを無視しない** — 必ず処理する
- **エクスポートされていないフィールドを使用しない** — バインドされない

## 次のステップ

@cards{cols="2"}
◆ サービス
サービスシステムについて詳しく学びます。

[詳しく見る →](/features/bindings/services/)

---
📖 モデル
複雑なデータ構造をバインドします。

[詳しく見る →](/features/bindings/models/)

---
🚀 Go-フロントエンドブリッジ
ブリッジの仕組みを理解します。

[詳しく見る →](/concepts/bridge/)

---
★ イベント
イベントを使用してパブリッシュ／サブスクライブ通信を行います。

[詳しく見る →](/features/events/system/)

@end

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[バインディングの例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)を確認してください。
