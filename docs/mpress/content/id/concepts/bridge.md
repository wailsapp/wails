---
title: "Bridge Go-Frontend"
description: "Pembahasan mendalam tentang cara Wails memungkinkan komunikasi langsung antara Go dan JavaScript"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Komunikasi Langsung Go-JavaScript

Wails menyediakan **bridge langsung dalam memori** antara Go dan JavaScript, sehingga keduanya dapat berkomunikasi dengan lancar tanpa overhead HTTP, batas proses, atau hambatan serialisasi.

## Gambaran Besar

```d2
direction: right

Frontend: Frontend (JavaScript) {
  UI: React/Vue/Vanilla {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: Binding yang Dibuat Otomatis {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Bridge Wails {
  Encoder: Encoder JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: Router Metode {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Decoder JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: Generator Tipe {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: Backend (Go) {
  Services: Layanan Anda {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: Registri Layanan {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Panggil Method('arg')"
Bridge.Encoder -> Bridge.Router: Enkode ke JSON
Bridge.Router -> Backend.Registry: Cari layanan
Backend.Registry -> Backend.Services: Panggil metode
Backend.Services -> Bridge.Decoder: Kembalikan hasil
Bridge.Decoder -> Frontend.Bindings: Dekode ke JS
Frontend.Bindings -> Frontend.UI: Promise diselesaikan
Bridge.TypeGen -> Frontend.Bindings: Buat tipe
```

**Inti penting:** Tanpa HTTP, tanpa IPC, tanpa batas proses. Hanya **pemanggilan fungsi langsung** dengan **keamanan tipe**.

## Cara Kerjanya: Langkah demi Langkah

### 1. Pendaftaran Layanan (Saat Mulai)

Saat aplikasi dimulai, Wails memindai layanan Anda:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) Add(a, b int) int {
    return a + b
}

// Register service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Yang dilakukan Wails:**

1. **Memindai struct** untuk mencari metode yang diekspor
2. **Mengekstrak informasi tipe** (parameter, tipe nilai kembalian)
3. **Membuat registri** yang memetakan nama metode ke fungsi
4. **Membuat binding TypeScript** dengan definisi tipe lengkap

### 2. Pembuatan Binding (Saat Build)

Wails membuat binding TypeScript secara otomatis:

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**Pemetaan tipe:**

| Tipe Go | Tipe TypeScript |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | Pengecualian (dilempar) |

### 3. Pemanggilan Frontend (Runtime)

Developer memanggil metode Go dari JavaScript:

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**Yang terjadi:**

1. **Fungsi binding dipanggil** - `Greet("World")`
2. **Pesan dibuat** - `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **Dikirim ke bridge** - Melalui bridge JavaScript WebView
4. **Promise dikembalikan** - Menunggu respons

### 4. Pemrosesan Bridge (Runtime)

Bridge menerima dan memproses pesan:

```d2
direction: down

Receive: Terima Pesan {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: Urai JSON {
  shape: rectangle
}

Validate: Validasi {
  Check: Layanan tersedia? {
    shape: diamond
  }

  CheckMethod: Metode tersedia? {
    shape: diamond
  }

  CheckTypes: Tipe sudah benar? {
    shape: diamond
  }
}

Invoke: Panggil Metode Go {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: Enkode Hasil {
  shape: rectangle
}

Send: Kirim Respons {
  shape: rectangle
  style.fill: "#10B981"
}

Error: Kirim Kesalahan {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: Ya
Validate.Check -> Error: Tidak
Validate.CheckMethod -> Validate.CheckTypes: Ya
Validate.CheckMethod -> Error: Tidak
Validate.CheckTypes -> Invoke: Ya
Validate.CheckTypes -> Error: Tidak
Invoke -> Encode: Berhasil
Invoke -> Error: Kesalahan
Encode -> Send
```

**Keamanan:** Hanya layanan yang terdaftar dan metode yang diekspor yang dapat dipanggil.

### 5. Eksekusi Go (Runtime)

Metode Go dieksekusi:

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**Konteks eksekusi:**

- Berjalan dalam sebuah **goroutine** (non-blocking)
- Memiliki akses ke **semua fitur Go** (sistem berkas, jaringan, basis data)
- Dapat memanggil **kode Go lainnya** secara bebas
- Mengembalikan hasil atau galat

### 6. Respons (Runtime)

Hasil dikirim kembali ke JavaScript:

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**Penanganan galat:**

```go
func (g *GreetService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```javascript
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Go error:", error)  // "division by zero"
}
```

## Karakteristik Performa

### Kecepatan

**Overhead panggilan pada umumnya:** &lt;1ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**Dibandingkan dengan alternatif:**

- **HTTP/REST:** 5-50ms (tumpukan jaringan, serialisasi)
- **IPC:** 1-10ms (batas proses, marshalling)
- **Wails Bridge:** &lt;1ms (dalam memori, panggilan langsung)

### Memori

**Overhead per panggilan:** ~1KB (buffer pesan)

**Optimasi tanpa penyalinan:** Data berukuran besar (>1MB) menggunakan memori bersama jika memungkinkan.

### Konkurensi

**Panggilan berlangsung secara konkuren:**

- Setiap panggilan berjalan dalam goroutine tersendiri
- Beberapa panggilan dapat dieksekusi secara bersamaan
- Tidak ada pemblokiran antar-panggilan

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## Sistem Tipe

### Tipe yang Didukung

#### Tipe Primitif

```go
// Go
func Example(
    s string,
    i int,
    f float64,
    b bool,
) (string, int, float64, bool) {
    return s, i, f, b
}
```

```typescript
// TypeScript (auto-generated)
function Example(
    s: string,
    i: number,
    f: number,
    b: boolean,
): Promise<[string, number, number, boolean]>
```

#### Slice dan Array

```go
// Go
func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```typescript
// TypeScript
function Sum(numbers: number[]): Promise<number>

// Usage
const total = await Sum([1, 2, 3, 4, 5])  // 15
```

#### Map

```go
// Go
func GetConfig() map[string]interface{} {
    return map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
        "enabled": true,
    }
}
```

```typescript
// TypeScript
function GetConfig(): Promise<Record<string, any>>

// Usage
const config = await GetConfig()
console.log(config.theme)  // "dark"
```

#### Struct

```go
// Go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(id int) (*User, error) {
    return &User{
        ID:    id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}
```

```typescript
// TypeScript (auto-generated)
interface User {
    id: number
    name: string
    email: string
}

function GetUser(id: number): Promise<User>

// Usage
const user = await GetUser(1)
console.log(user.name)  // "Alice"
```

**Tag JSON:** Gunakan tag `json:` untuk mengatur nama field di TypeScript.

#### Waktu

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<Date>

// Usage
const timestamp = await GetTimestamp()
console.log(timestamp.toISOString())
```

#### Galat

```go
// Go
func Validate(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}
```

```typescript
// TypeScript
function Validate(input: string): Promise<void>

// Usage
try {
    await Validate("")
} catch (error) {
    console.error(error)  // "input cannot be empty"
}
```

### Tipe yang Tidak Didukung

Tipe-tipe ini **tidak dapat** diteruskan melalui bridge:

- **Channel** (`chan T`)
- **Fungsi** (`func()`)
- **Antarmuka** (kecuali `interface{}` / `any`)
- **Pointer** (kecuali pointer ke struct)
- **Field yang tidak diekspor** (huruf kecil)

**Solusi alternatif:** Gunakan ID atau handle:

```go
// ❌ Can't pass file handle
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

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

## Pola Tingkat Lanjut

### Penerusan Context

Layanan dapat mengakses context panggilan:

```go
type UserService struct{}

func (s *UserService) GetCurrentUser(ctx context.Context) (*User, error) {
    // Access the calling window via the context value
    window, _ := ctx.Value(application.WindowKey).(application.Window)
    _ = window

    // Access the application
    app := application.Get()
    _ = app

    // Your logic
    return getCurrentUser(), nil
}
```

**Context menyediakan:**

- Jendela yang melakukan panggilan
- Instans aplikasi
- Metadata permintaan

### Aliran data

Untuk data berukuran besar, gunakan event sebagai pengganti nilai kembalian:

```go
func ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        // Emit progress events
        app.Event.Emit("file-progress", map[string]interface{}{
            "line": lineNum,
            "text": scanner.Text(),
        })
    }

    return scanner.Err()
}
```

```javascript
import { Events } from '@wailsio/runtime'
import { ProcessLargeFile } from './bindings/FileService'

// Listen for progress
Events.On('file-progress', (data) => {
    console.log(`Line ${data.line}: ${data.text}`)
})

// Start processing
await ProcessLargeFile('/path/to/large/file.txt')
```

### Pembatalan

Gunakan context untuk operasi yang dapat dibatalkan:

```go
func LongRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check if cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

**Catatan:** Pembatalan context saat frontend terputus berlangsung secara otomatis.

### Operasi Batch

Kurangi overhead bridge dengan mengelompokkan operasi dalam batch:

```go
// ❌ Inefficient: N bridge calls
for _, item := range items {
    await ProcessItem(item)
}

// ✅ Efficient: 1 bridge call
await ProcessItems(items)
```

```go
func ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = processItem(item)
    }
    return results, nil
}
```

## Men-debug Bridge

### Aktifkan Logging Debug

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**Output menampilkan:**

- Panggilan metode
- Parameter
- Nilai kembalian
- Error
- Informasi waktu eksekusi

### Periksa Binding yang Dihasilkan

Periksa `frontend/bindings/` untuk melihat TypeScript yang dihasilkan:

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### Uji Layanan Secara Langsung

Uji layanan Go tanpa frontend:

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## Tips Performa

### ✅ Lakukan

- **Kelompokkan operasi dalam batch** - Kurangi panggilan bridge
- **Gunakan event untuk streaming** - Jangan kembalikan array berukuran besar
- **Pastikan metode berjalan cepat** - Idealnya &lt;100 ms
- **Gunakan goroutine** - Untuk operasi yang berlangsung lama
- **Gunakan cache di sisi Go** - Hindari perhitungan berulang

### ❌ Jangan Lakukan

- **Jangan lakukan panggilan secara berlebihan** - Kelompokkan dalam batch jika memungkinkan
- **Jangan kembalikan data berukuran sangat besar** - Gunakan paginasi atau streaming
- **Jangan memblokir** - Gunakan goroutine untuk operasi yang berlangsung lama
- **Jangan teruskan tipe kompleks** - Gunakan tipe yang sederhana
- **Jangan abaikan error** - Selalu tangani error tersebut

## Keamanan

Bridge aman secara default:

1. **Hanya daftar yang diizinkan** - Hanya layanan terdaftar yang dapat dipanggil
2. **Validasi tipe** - Argumen diperiksa berdasarkan tipe Go
3. **Tanpa eval()** - Frontend tidak dapat mengeksekusi kode Go arbitrer
4. **Tidak ada penyalahgunaan refleksi** - Hanya metode yang diekspor yang dapat diakses

**Praktik terbaik:**

- **Validasi input** di Go (jangan percayai frontend)
- **Gunakan context** untuk autentikasi/otorisasi
- **Batasi laju** operasi yang mahal
- **Sanitasi** path file dan input pengguna

## Langkah Berikutnya

**Sistem Build** - Pelajari cara Wails membangun dan membundel aplikasi Anda [Pelajari Selengkapnya →](/concepts/build-system/)

**Layanan** - Pelajari sistem layanan secara mendalam [Pelajari Selengkapnya →](/features/bindings/services/)

**Peristiwa** - Gunakan peristiwa untuk komunikasi publikasi/langganan [Pelajari Lebih Lanjut →](/features/events/system/)

---

**Ada pertanyaan tentang bridge?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh binding](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
