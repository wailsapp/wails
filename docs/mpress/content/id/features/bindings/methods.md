---
title: "Binding Metode"
description: "Panggil metode Go dari JavaScript dengan keamanan tipe"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## Binding Go-JavaScript yang Aman secara Tipe

Wails **secara otomatis menghasilkan binding JavaScript/TypeScript yang aman secara tipe** untuk metode Go Anda. Tulis kode Go, jalankan satu perintah, lalu dapatkan fungsi frontend dengan tipe lengkap tanpa overhead HTTP, tanpa pekerjaan manual, dan tanpa kode boilerplate.

## Mulai Cepat

**1. Tulis layanan Go:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. Daftarkan layanan:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. Hasilkan binding:**

```bash
wails3 generate bindings
```

**4. Gunakan di JavaScript:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

**Selesai!** Panggilan dari Go ke JavaScript yang aman secara tipe.

## Membuat Layanan

### Layanan Dasar

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

**Daftarkan:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Poin penting:**

- Hanya **metode yang diekspor** (PascalCase) yang dibuatkan binding
- Metode dapat mengembalikan nilai atau `(value, error)`
- Layanan merupakan **singleton** (satu instans per aplikasi)

### Layanan dengan State

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

**Penting:** Layanan digunakan bersama oleh semua jendela. Gunakan mutex untuk keamanan thread.

### Layanan dengan Dependensi

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

**Daftarkan dengan dependensi:**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## Menghasilkan Binding

### Pembuatan Dasar

```bash
wails3 generate bindings
```

**Output:**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**Struktur yang dihasilkan:**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### Pembuatan TypeScript

```bash
wails3 generate bindings -ts
```

**Menghasilkan file `.ts`** dengan tipe TypeScript lengkap.

### Direktori Output Khusus

```bash
wails3 generate bindings -d ./src/bindings
```

### Mode Pemantauan (Pengembangan)

```bash
wails3 dev
```

**Secara otomatis menghasilkan ulang binding** ketika kode Go berubah.

## Menggunakan Binding

### JavaScript

**Binding yang dihasilkan:**

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

**Penggunaan:**

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

**Binding yang dihasilkan:**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**Penggunaan:**

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

**Manfaat:**

- Pemeriksaan tipe lengkap
- Pelengkapan otomatis IDE
- Kesalahan pada waktu kompilasi
- Refactoring yang lebih baik

### File Indeks

**Indeks yang dihasilkan:**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**Impor yang disederhanakan:**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## Pemetaan Tipe

### Tipe Primitif

| Tipe Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Tipe Kompleks

| Tipe Go | JavaScript/TypeScript | Catatan |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | map dengan kunci string |
| `map[K]V` | `{ [_ in K]?: V }` | `K` non-string dirender sebagai tipe terpetakan, **bukan** `Map` JS |
| `[]byte` | `string` | dikodekan dengan base64 |
| `struct` | `class` / `interface` | dengan field |
| `time.Time` | `any` | diserialisasi sebagai string RFC3339Nano saat runtime |
| `*T` | `T \| null` | pointer berarti dapat bernilai null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Exception jika digunakan sebagai nilai kembalian; jika tidak, any |

### Tipe yang Tidak Didukung

Tipe-tipe ini **tidak dapat** diteruskan melalui bridge:

- `chan T` (channel)
- `func()` (fungsi)
- Interface kompleks (kecuali `interface{}`)
- Field yang tidak diekspor (huruf kecil)

**Solusi alternatif:** Gunakan ID atau handle:

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

## Penanganan Error

### Sisi Go

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

### Sisi JavaScript

Ketika metode yang terikat gagal, promise yang dikembalikan ditolak dengan objek `Error` JavaScript:

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

Runtime melempar tipe error yang berbeda bergantung pada penyebab kegagalan:

| Tipe error | Dilempar ketika |
| --- | --- |
| `TypeError` | Pemanggilan memiliki jumlah argumen yang salah, atau suatu argumen tidak dapat dikonversi ke tipe Go |
| `RuntimeError` | Metode mengembalikan error atau mengalami panic saat dijalankan |
| `Error` | Kegagalan lainnya, misalnya pemanggilan metode yang tidak ada |

Setiap error menyediakan:

- `name`: tipe error dari tabel di atas
- `message`: pesan error Go
- `cause`: error Go yang diserialisasi sebagai JSON, jika tersedia. Jika metode mengembalikan beberapa error, `cause` berupa array dengan satu entri untuk setiap error.

Kelas `RuntimeError` diekspor oleh paket `@wailsio/runtime`, sehingga Anda dapat membedakan error yang dikembalikan oleh kode Go Anda dari kegagalan lainnya:

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
Pada versi runtime sebelumnya, setiap pemanggilan yang gagal menolak promise dengan `Error` biasa yang pesannya berisi JSON mentah dari error yang mendasarinya.

@end

### Data Error Terstruktur

Mengembalikan tipe error khusus dari Go membuat bentuk JSON-nya tersedia pada properti `cause` dari error yang dilempar:

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

Error diserialisasi dengan paket standar `encoding/json`, sehingga hanya field yang diekspor yang disertakan. Error yang dibuat dengan `errors.New` atau `fmt.Errorf` tidak memiliki field yang diekspor dan diserialisasi sebagai objek kosong.

Untuk mengontrol sepenuhnya cara error diserialisasi, sediakan fungsi `MarshalError` dalam opsi layanan:

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

`MarshalError` harus mengembalikan JSON yang valid, atau `nil` untuk kembali menggunakan serialisasi default.

## Performa

### Overhead Pemanggilan

**Panggilan umum:** &lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**Dibandingkan dengan alternatif:**

- HTTP/REST: 5-50ms
- IPC: 1-10ms
- Wails: &lt;1ms

### Tips Optimasi

**✅ Kelompokkan operasi:**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ Simpan hasil dalam cache:**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ Gunakan event untuk streaming:**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## Contoh Lengkap

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

## Praktik Terbaik

### ✅ Lakukan

- **Buat metode tetap sederhana** - Satu tanggung jawab
- **Kembalikan error** - Jangan memicu panic
- **Gunakan state yang aman untuk thread** - Gunakan mutex untuk data bersama
- **Kelompokkan operasi** - Kurangi panggilan melalui bridge
- **Simpan dalam cache di sisi Go** - Hindari pekerjaan berulang
- **Dokumentasikan metode** - Komentar menjadi JSDoc

### ❌ Jangan Lakukan

- **Jangan memblokir** - Gunakan goroutine untuk operasi yang berlangsung lama
- **Jangan kembalikan channel** - Gunakan event sebagai gantinya
- **Jangan kembalikan fungsi** - Tidak didukung
- **Jangan abaikan error** - Selalu tangani error tersebut
- **Jangan gunakan field yang tidak diekspor** - Field tersebut tidak akan dibuatkan binding

## Langkah Berikutnya

@cards{cols="2"}
◆ Layanan
Pelajari sistem layanan secara mendalam.

[Pelajari Lebih Lanjut →](/features/bindings/services/)

---
📖 Model
Buat binding untuk struktur data yang kompleks.

[Pelajari Lebih Lanjut →](/features/bindings/models/)

---
🚀 Bridge Go-Frontend
Pahami mekanisme bridge.

[Pelajari Lebih Lanjut →](/concepts/bridge/)

---
★ Event
Gunakan event untuk komunikasi pub/sub.

[Pelajari Lebih Lanjut →](/features/events/system/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh binding](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
