---
title: "Enum"
description: "Pembuatan enum otomatis dari konstanta Go"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Binding Enum

Generator binding Wails v3 **secara otomatis mendeteksi tipe konstanta Go dan menghasilkan enum TypeScript atau objek const JavaScript**. Tidak perlu registrasi maupun konfigurasi — cukup definisikan tipe dan konstanta Anda di Go, lalu generator akan menangani sisanya.

@note{type="info"}
Tidak seperti Wails v2, Anda **tidak perlu memanggil `EnumBind`** atau mendaftarkan enum secara manual. Generator menemukannya secara otomatis dari kode sumber Anda.

@end

## Mulai Cepat

**Definisikan tipe bernama dengan konstanta di Go:**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**Gunakan tipe tersebut dalam struct atau metode layanan:**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**Buat binding:**

```bash
wails3 generate bindings
```

Output generator akan melaporkan jumlah enum bersama jumlah model:

```
3 Enums, 5 Models
```

**Gunakan di frontend Anda:**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**Selesai!** Tipe enum diberlakukan di Go maupun JavaScript/TypeScript.

## Mendefinisikan Enum

Enum di Wails adalah **tipe bernama** dengan tipe dasar sederhana, yang dipadukan dengan **deklarasi const** bertipe tersebut.

### Enum String

```go
// Title is a title
type Title string

const (
    // Mister is a title
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
    Mrs    Title = "Mrs"
    Dr     Title = "Dr"
)
```

**TypeScript yang dihasilkan:**

```typescript
/**
 * Title is a title
 */
export enum Title {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    /**
     * Mister is a title
     */
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
    Mrs = "Mrs",
    Dr = "Dr",
}
```

**JavaScript yang dihasilkan:**

```javascript
/**
 * Title is a title
 * @readonly
 * @enum {string}
 */
export const Title = {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero: "",

    /**
     * Mister is a title
     */
    Mister: "Mr",
    Miss: "Miss",
    Ms: "Ms",
    Mrs: "Mrs",
    Dr: "Dr",
};
```

### Enum Bilangan Bulat

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**TypeScript yang dihasilkan:**

```typescript
export enum Priority {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = 0,

    PriorityLow = 0,
    PriorityMedium = 1,
    PriorityHigh = 2,
}
```

### Enum Alias Tipe

Alias tipe Go (`=`) juga didukung, tetapi menghasilkan output yang sedikit berbeda — definisi tipe beserta objek const, bukan `enum` TypeScript native:

```go
// Age is an integer with some predefined values
type Age = int

const (
    NewBorn    Age = 0
    Teenager   Age = 12
    YoungAdult Age = 18

    // Oh no, some grey hair!
    MiddleAged Age = 50
    Mathusalem Age = 1000 // Unbelievable!
)
```

**TypeScript yang dihasilkan:**

```typescript
/**
 * Age is an integer with some predefined values
 */
export type Age = number;

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

**JavaScript yang dihasilkan:**

```javascript
/**
 * Age is an integer with some predefined values
 * @typedef {number} Age
 */

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

@note{type="tip"}
**Tipe bernama** (`type Title string`) menghasilkan deklarasi `enum` TypeScript native dengan anggota `$zero`. **Alias tipe** (`type Age = int`) menghasilkan pasangan namespace `type` + `const` tanpa `$zero`.

@end

## Nilai `$zero`

Setiap enum bertipe bernama menyertakan anggota khusus `$zero` yang merepresentasikan **nilai nol Go** untuk tipe dasarnya:

| Tipe Dasar | Nilai `$zero` |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

Saat sebuah field struct menggunakan tipe enum dan tidak ada nilai yang diberikan, constructor secara default menggunakan `$zero`:

```typescript
export class Person {
    "Title": Title;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;  // defaults to ""
        }
        Object.assign(this, $$source);
    }
}
```

Hal ini memastikan inisialisasi yang aman secara tipe saat menghasilkan class — field enum tidak pernah bernilai `undefined`. Saat menghasilkan interface TypeScript (menggunakan `-i`), tidak ada constructor dan field dapat tidak disertakan seperti biasa.

## Menggunakan Enum dalam Struct

Saat sebuah field struct memiliki tipe enum, kode yang dihasilkan **mempertahankan tipe tersebut** alih-alih kembali menggunakan tipe primitif:

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**TypeScript yang dihasilkan:**

```typescript
export class Person {
    "Title": Title;
    "Name": string;
    "Age": Age;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Age" in $$source)) {
            this["Age"] = 0;
        }

        Object.assign(this, $$source);
    }
}
```

Field `Title` diberi tipe `Title`, bukan `string`. Dengan demikian, IDE Anda menyediakan pelengkapan otomatis dan pemeriksaan tipe sepenuhnya untuk nilai enum.

## Enum dari Paket yang Diimpor

Enum yang didefinisikan dalam paket terpisah didukung sepenuhnya. Enum tersebut dihasilkan dalam direktori paket yang sesuai:

```go
// services/types.go
package services

type Title string

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
)
```

```go
// main.go
package main

import "myapp/services"

func (*GreetService) Greet(name string, title services.Title) string {
    return "Hello " + string(title) + " " + name
}
```

Enum `Title` dihasilkan dalam file model `services`, dan path impor ditentukan secara otomatis:

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## Metode Enum

Anda dapat menambahkan metode ke tipe enum di Go. Metode tersebut tidak memengaruhi pembuatan binding, tetapi menyediakan fungsionalitas sisi server yang berguna:

```go
type Title string

func (t Title) String() string {
    return string(t)
}

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
)
```

Enum yang dihasilkan tetap identik, baik ada maupun tidak ada metode Go pada tipe tersebut.

## Komentar dan Dokumentasi

Generator mempertahankan komentar Go sebagai JSDoc dalam output yang dihasilkan:

- **Komentar tipe** menjadi komentar dokumentasi enum
- **Komentar grup const** menjadi pemisah bagian
- **Komentar pada masing-masing const** menjadi komentar dokumentasi anggota
- **Komentar inline** dipertahankan jika memungkinkan

Dengan demikian, IDE Anda akan menampilkan dokumentasi nilai enum saat penunjuk diarahkan ke nilai tersebut.

## Tipe Dasar yang Didukung

Generator binding mendukung enum dengan tipe dasar Go berikut:

| Tipe Go | Berfungsi sebagai Enum |
| --- | :---: |
| `string` | Ya |
| `int`, `int8`, `int16`, `int32`, `int64` | Ya |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Ya |
| `float32`, `float64` | Ya |
| `byte` (`uint8`) | Ya |
| `rune` (`int32`) | Ya |
| `bool` | Ya |
| `complex64`, `complex128` | Tidak |

## Keterbatasan

Berikut ini **tidak** didukung untuk pembuatan enum:

- **Tipe generik** — Parameter tipe mencegah pendeteksian konstanta
- **Tipe dengan `json.Marshaler` atau `encoding.TextMarshaler`** kustom — Serialisasi kustom berarti nilai yang dihasilkan mungkin tidak sesuai dengan perilaku saat runtime, sehingga generator melewati tipe tersebut
- **Konstanta yang nilainya tidak dapat dievaluasi atau direpresentasikan secara statis** — Konstanta harus memiliki nilai yang diketahui dan dapat direpresentasikan dalam tipe dasarnya. Pola `iota` standar dapat digunakan tanpa masalah karena compiler menguraikannya menjadi nilai konkret
- **Tipe bilangan kompleks** — `complex64` dan `complex128` tidak dapat menjadi tipe dasar enum

## Contoh Lengkap

**Go:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

// BackgroundType defines the type of background
type BackgroundType string

const (
    BackgroundSolid    BackgroundType = "solid"
    BackgroundGradient BackgroundType = "gradient"
    BackgroundImage    BackgroundType = "image"
)

type BackgroundConfig struct {
    Type  BackgroundType `json:"type"`
    Value string         `json:"value"`
}

type ThemeService struct{}

func (*ThemeService) GetBackground() BackgroundConfig {
    return BackgroundConfig{
        Type:  BackgroundSolid,
        Value: "#ffffff",
    }
}

func (*ThemeService) SetBackground(config BackgroundConfig) error {
    // Apply background
    return nil
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&ThemeService{}),
        },
    })
    app.Window.New()
    app.Run()
}
```

**Frontend (TypeScript):**

```typescript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Get current background
const bg = await GetBackground()

// Check the type using enum values
if (bg.type === BackgroundType.BackgroundSolid) {
    console.log("Solid background:", bg.value)
}

// Set a new background
await SetBackground(new BackgroundConfig({
    type: BackgroundType.BackgroundGradient,
    value: "linear-gradient(to right, #000, #fff)"
}))
```

**Frontend (JavaScript):**

```javascript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Use enum values for type-safe comparisons
const bg = await GetBackground()

switch (bg.type) {
    case BackgroundType.BackgroundSolid:
        applySolid(bg.value)
        break
    case BackgroundType.BackgroundGradient:
        applyGradient(bg.value)
        break
    case BackgroundType.BackgroundImage:
        applyImage(bg.value)
        break
}
```

## Langkah Berikutnya

@cards{cols="2"}
📖 Model Data
Struct, pemetaan tipe, dan pembuatan model.

[Pelajari Lebih Lanjut →](/features/bindings/models/)

---
🚀 Binding Metode
Hubungkan metode Go ke frontend.

[Pelajari Lebih Lanjut →](/features/bindings/methods/)

---
⚙ Binding Lanjutan
Direktif, injeksi kode, dan ID kustom.

[Pelajari Lebih Lanjut →](/features/bindings/advanced/)

---
✓ Praktik Terbaik
Pola desain binding.

[Pelajari Lebih Lanjut →](/features/bindings/best-practices/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh binding](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
