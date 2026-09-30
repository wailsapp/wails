---
title: "Binding Tingkat Lanjut"
description: "Teknik binding tingkat lanjut, termasuk direktif, injeksi kode, dan ID khusus"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

Panduan ini membahas teknik tingkat lanjut untuk menyesuaikan dan mengoptimalkan proses pembuatan binding di Wails v3.

## Menyesuaikan Kode yang Dihasilkan dengan Direktif

### Menginjeksikan Kode Khusus

Direktif `//wails:inject` memungkinkan Anda menginjeksikan kode JavaScript/TypeScript khusus ke dalam binding yang dihasilkan:

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

Tindakan ini akan menginjeksikan kode yang ditentukan ke dalam file JavaScript/TypeScript yang dihasilkan untuk layanan `MyService`.

Anda juga dapat menggunakan injeksi kondisional untuk menargetkan format keluaran tertentu:

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### Menyertakan File Tambahan

Direktif `//wails:include` memungkinkan Anda menyertakan file tambahan bersama binding yang dihasilkan:

```go
//wails:include js/*.js
package mypackage
```

Direktif ini biasanya digunakan dalam komentar dokumentasi paket untuk menyertakan file JavaScript/TypeScript tambahan bersama binding yang dihasilkan.

### Menandai Tipe dan Metode Internal

Direktif `//wails:internal` menandai tipe atau metode sebagai internal agar tidak diekspor ke frontend:

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Ini berguna untuk tipe dan metode yang hanya digunakan secara internal oleh kode Go Anda dan sebaiknya tidak diekspos ke frontend.

### Mengabaikan Metode

Direktif `//wails:ignore` sepenuhnya mengabaikan suatu metode selama pembuatan binding:

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

Ini serupa dengan `//wails:internal`, tetapi sepenuhnya mengabaikan metode tersebut alih-alih menandainya sebagai internal.

### ID Metode Khusus

Direktif `//wails:id` menentukan ID khusus untuk suatu metode, menggantikan ID bawaan berbasis hash:

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Ini dapat berguna untuk mempertahankan kompatibilitas saat melakukan pemfaktoran ulang kode.

## Menggunakan Tipe Kompleks

### Struct Bersarang

Generator binding menangani struct bersarang secara otomatis:

```go
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Person struct {
    Name    string
    Address Address
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Address: Address{
            Street: "123 Main St",
            City:   "Anytown",
            State:  "CA",
            Zip:    "12345",
        },
    }
}
```

Kode JavaScript/TypeScript yang dihasilkan akan menyertakan kelas untuk `Person` dan `Address`.

### Map dan Slice

Map dan slice juga ditangani secara otomatis:

```go
type Person struct {
    Name       string
    Attributes map[string]string
    Friends    []string
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Attributes: map[string]string{
            "hair": "brown",
            "eyes": "blue",
        },
        Friends: []string{"Jane", "Bob", "Alice"},
    }
}
```

Dalam JavaScript, map direpresentasikan sebagai objek dan slice sebagai array. Dalam TypeScript, map direpresentasikan sebagai `Record<K, V>` dan slice sebagai `T[]`.

### Tipe Generik

Generator binding mendukung tipe generik:

```go
type Result[T any] struct {
    Data  T
    Error string
}

func (s *MyService) GetResult() Result[string] {
    return Result[string]{
        Data:  "Hello, World!",
        Error: "",
    }
}
```

Kode TypeScript yang dihasilkan akan menyertakan kelas generik untuk `Result`:

```typescript
export class Result<T> {
    "Data": T;
    "Error": string;

    constructor(source: Partial<Result<T>> = {}) {
        if (!("Data" in source)) {
            this["Data"] = null as any;
        }
        if (!("Error" in source)) {
            this["Error"] = "";
        }

        Object.assign(this, source);
    }

    static createFrom<T>(source: string | object = {}): Result<T> {
        let parsedSource = typeof source === "string" ? JSON.parse(source) : source;
        return new Result<T>(parsedSource as Partial<Result<T>>);
    }
}
```

### Antarmuka

Generator binding dapat menghasilkan antarmuka TypeScript sebagai pengganti kelas dengan menggunakan flag `-i`:

```bash
wails3 generate bindings -ts -i
```

Tindakan ini akan menghasilkan antarmuka TypeScript untuk semua model:

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## Mengoptimalkan Pembuatan Binding

### Menggunakan Nama sebagai Pengganti ID

Secara bawaan, generator binding menggunakan ID berbasis hash untuk pemanggilan metode. Anda dapat menggunakan flag `-names` untuk menggunakan nama sebagai gantinya:

```bash
wails3 generate bindings -names
```

Tindakan ini akan menghasilkan kode yang memanggil metode terikat berdasarkan nama **lengkap** (`<package>.<Service>.<Method>`) dan menggunakan nama parameter posisional `$0`, `$1`, …:

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Ini dapat membuat kode yang dihasilkan lebih mudah dibaca dan di-debug, tetapi mungkin sedikit kurang efisien.

### Membundel Runtime

Secara bawaan, kode yang dihasilkan mengimpor runtime Wails dari paket npm `@wailsio/runtime`. Anda dapat menggunakan flag `-b` untuk membundel runtime bersama kode yang dihasilkan:

```bash
wails3 generate bindings -b
```

Tindakan ini akan menyertakan kode runtime secara langsung dalam file yang dihasilkan sehingga paket npm tidak lagi diperlukan.

### Menonaktifkan File Indeks

Jika tidak memerlukan file indeks, Anda dapat menggunakan flag `-noindex` untuk menonaktifkan pembuatannya:

```bash
wails3 generate bindings -noindex
```

Ini dapat berguna jika Anda lebih memilih mengimpor layanan dan model secara langsung dari file masing-masing.

## Contoh Dunia Nyata

### Layanan Autentikasi

Berikut adalah contoh layanan autentikasi dengan direktif khusus:

```go
package auth

//wails:inject console.log("Auth service initialized");
type AuthService struct {
    // Private fields
    users map[string]User
}

type User struct {
    Username string
    Email    string
    Role     string
}

type LoginRequest struct {
    Username string
    Password string
}

type LoginResponse struct {
    Success bool
    User    User
    Token   string
    Error   string
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) LoginResponse {
    // Implementation...
}

// GetCurrentUser returns the current user
func (s *AuthService) GetCurrentUser() User {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *AuthService) validateCredentials(username, password string) bool {
    // Implementation...
}
```

### Layanan Pemrosesan Data

Berikut adalah contoh layanan pemrosesan data dengan tipe generik:

```go
package data

type ProcessingResult[T any] struct {
    Data  T
    Error string
}

type DataService struct {}

// Process processes data and returns a result
func (s *DataService) Process(data string) ProcessingResult[map[string]int] {
    // Implementation...
}

// ProcessBatch processes multiple data items
func (s *DataService) ProcessBatch(data []string) ProcessingResult[[]map[string]int] {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *DataService) parseData(data string) (map[string]int, error) {
    // Implementation...
}
```

### Injeksi Kode Kondisional

Berikut adalah contoh injeksi kode kondisional untuk format keluaran yang berbeda:

```go
//wails:inject j*:/**
//wails:inject j*: * @param {string} arg
//wails:inject j*: * @returns {Promise<void>}
//wails:inject j*: */
//wails:inject j*:export async function CustomMethod(arg) {
//wails:inject t*:export async function CustomMethod(arg: string): Promise<void> {
//wails:inject     await InternalMethod("Hello " + arg + "!");
//wails:inject }
type Service struct{}
```

Ini menyisipkan kode yang berbeda untuk keluaran JavaScript dan TypeScript, dengan menyediakan anotasi tipe yang sesuai untuk setiap bahasa.
