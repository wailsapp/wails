---
title: "Standar Pengodean"
description: "Gaya kode, konvensi, dan praktik terbaik untuk Wails v3"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## Gaya dan Konvensi Kode

Penerapan standar pengodean yang konsisten membuat basis kode lebih mudah dibaca dan dipelihara, serta memudahkan pengembang untuk berkontribusi pada basis kode tersebut.

## Standar Kode Go

### Pemformatan Kode

Gunakan alat pemformatan Go standar:

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**Wajib:** Semua kode Go harus lolos `gofmt` dan `goimports` sebelum di-commit.

### Konvensi Penamaan

**Paket:**

- Gunakan huruf kecil dan satu kata jika memungkinkan
- `package application`, `package events`
- Hindari garis bawah atau campuran huruf besar dan kecil

**Nama yang diekspor:**

- Gunakan PascalCase untuk tipe, fungsi, dan konstanta
- `type WebviewWindow struct`, `func NewApplication()`

**Nama yang tidak diekspor:**

- Gunakan camelCase untuk tipe, fungsi, dan variabel internal
- `type windowImpl struct`, `func createWindow()`

**Antarmuka:**

- Beri nama berdasarkan perilaku: `Reader`, `Writer`, `Handler`
- Untuk antarmuka dengan satu metode, gunakan nama dengan akhiran `-er`

```go
// Good
type Closer interface {
    Close() error
}

// Avoid
type CloseInterface interface {
    Close() error
}
```

### Penanganan Kesalahan

**Selalu periksa kesalahan:**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**Gunakan pembungkusan kesalahan:**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**Buat tipe kesalahan khusus bila diperlukan:**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### Komentar dan Dokumentasi

**Komentar paket:**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**Deklarasi yang diekspor:**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**Komentar implementasi:**

```go
// processEvent handles incoming events from the runtime.
// It dispatches to registered handlers and manages event lifecycle.
func (a *Application) processEvent(event *Event) {
    // Validate event before processing
    if event == nil {
        return
    }

    // Find and invoke handlers
    // ...
}
```

### Struktur Fungsi dan Metode

**Pastikan fungsi tetap terfokus:**

```go
// Good - single responsibility
func (w *Window) setTitle(title string) {
    w.title = title
    w.updateNativeTitle()
}

// Bad - doing too much
func (w *Window) updateEverything() {
    w.setTitle(w.title)
    w.setSize(w.width, w.height)
    w.setPosition(w.x, w.y)
    // ... 20 more operations
}
```

**Gunakan pengembalian awal:**

```go
// Good
func validate(input string) error {
    if input == "" {
        return errors.New("empty input")
    }

    if len(input) > 100 {
        return errors.New("input too long")
    }

    return nil
}

// Avoid deep nesting
```

### Konkurensi

**Gunakan context untuk pembatalan:**

```go
func (a *Application) RunWithContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-a.done:
        return nil
    }
}
```

**Lindungi status bersama dengan mutex:**

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**Hindari kebocoran goroutine:**

```go
// Good - goroutine has exit condition
func (a *Application) startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            case work := <-a.workChan:
                a.process(work)
            }
        }
    }()
}
```

### Pengujian

**Penamaan berkas pengujian:**

```go
// Implementation: window.go
// Tests: window_test.go
```

**Pengujian berbasis tabel:**

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty input", "", true},
        {"valid input", "hello", false},
        {"too long", strings.Repeat("a", 101), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Standar JavaScript/TypeScript

### Pemformatan Kode

Gunakan Prettier untuk pemformatan yang konsisten:

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### Konvensi Penamaan

**Variabel dan fungsi:**

- camelCase: `const userName = "John"`

**Kelas dan tipe:**

- PascalCase: `class WindowManager`

**Konstanta:**

- UPPER<em>SNAKE</em>CASE: `const MAX_RETRIES = 3`

### TypeScript

**Gunakan tipe eksplisit:**

```typescript
// Good
function greet(name: string): string {
    return `Hello, ${name}`
}

// Avoid implicit any
function process(data) {  // Bad
    return data
}
```

**Definisikan antarmuka:**

```typescript
interface WindowOptions {
    title: string
    width: number
    height: number
}

function createWindow(options: WindowOptions): void {
    // ...
}
```

## Format Pesan Commit

Gunakan [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Tipe:**

- `feat`: Fitur baru
- `fix`: Perbaikan bug
- `docs`: Perubahan dokumentasi
- `refactor`: Refaktorisasi kode
- `test`: Penambahan atau pembaruan pengujian
- `chore`: Tugas pemeliharaan

**Contoh:**

```
feat(window): add SetAlwaysOnTop method

Implement SetAlwaysOnTop for keeping windows above others.
Adds platform implementations for macOS, Windows, and Linux.

Closes #123
```

```
fix(events): prevent event handler memory leak

Event listeners were not being properly cleaned up when
windows were closed. This adds explicit cleanup in the
window destructor.
```

## Pedoman Pull Request

### Sebelum Mengirimkan

- [ ] Kode lolos `gofmt` dan `goimports`
- [ ] Semua pengujian lulus (`go test ./...`)
- [ ] Kode baru dilengkapi pengujian
- [ ] Dokumentasi diperbarui jika diperlukan
- [ ] Pesan commit mengikuti konvensi
- [ ] Tidak ada konflik penggabungan dengan `master`

### Templat Deskripsi PR

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## Proses Peninjauan Kode

### Sebagai Peninjau

- Bersikaplah konstruktif dan penuh hormat
- Fokuslah pada kualitas kode, bukan preferensi pribadi
- Jelaskan alasan perubahan disarankan
- Setujui setelah merasa puas

### Sebagai Penulis

- Tanggapi semua komentar
- Mintalah klarifikasi jika diperlukan
- Lakukan perubahan yang diminta atau jelaskan alasan Anda tidak melakukannya
- Terbukalah terhadap umpan balik

## Praktik Terbaik

### Performa

- Hindari optimasi prematur
- Lakukan profiling sebelum mengoptimalkan
- Gunakan benchmark untuk kode yang sangat memengaruhi performa

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### Keamanan

- Validasi semua masukan pengguna
- Sanitasi data sebelum ditampilkan
- Gunakan `crypto/rand` untuk data acak
- Jangan pernah mencatat informasi sensitif dalam log

### Dokumentasi

- Dokumentasikan API yang diekspor
- Sertakan contoh dalam dokumentasi
- Perbarui dokumentasi saat mengubah API
- Pastikan file README selalu mutakhir

## Kode Khusus Platform

### Penamaan File

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### Tag Build

```go
//go:build darwin

package application

// macOS-specific code
```

## Linting

Jalankan linter sebelum melakukan commit:

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## Ada Pertanyaan?

Jika Anda tidak yakin tentang standar apa pun:

- Periksa kode yang ada untuk melihat contoh
- Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf)
- Buka diskusi di GitHub
