---
title: "Enum"
description: "Pembuatan enum otomatis dari konstanta Go"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Binding Enum

Generator binding Wails v3 **secara otomatis mendeteksi tipe konstanta Go dan menghasilkan enum TypeScript atau objek const JavaScript**. Tidak perlu registrasi atau konfigurasi—cukup definisikan tipe dan konstanta Anda di Go, lalu generator akan menangani sisanya.

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

Keluaran generator akan melaporkan jumlah enum bersama dengan model:

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

Enum di Wails adalah **tipe bernama** dengan tipe dasar primitif, yang dipadukan dengan **deklarasi const** bertipe tersebut.

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

### Enum Integer

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

Alias tipe Go (`=`) juga didukung, tetapi menghasilkan keluaran yang sedikit berbeda—definisi tipe beserta objek const, bukan `enum` TypeScript native:

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

Jika sebuah field struct menggunakan tipe enum dan tidak ada nilai yang diberikan, konstruktor secara default menggunakan `$zero`:

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

Hal ini memastikan inisialisasi yang aman terhadap tipe saat menghasilkan class—field enum tidak pernah bernilai `undefined`. Saat menghasilkan interface TypeScript (menggunakan `-i`), tidak ada konstruktor dan field dapat tidak disertakan seperti biasa.

## Menggunakan Enum dalam Struct

Jika sebuah field struct memiliki tipe enum, kode yang dihasilkan **mempertahankan tipe tersebut** dan tidak menggantinya dengan tipe primitif:

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

Field `Title` bertipe `Title`, bukan `string`. Dengan demikian, IDE Anda dapat menyediakan pelengkapan otomatis dan pemeriksaan tipe sepenuhnya untuk nilai enum.

## Enum dari Paket yang Diimpor

Enum yang didefinisikan dalam paket terpisah didukung sepenuhnya. Enum tersebut dibuat di direktori paket yang sesuai:

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

Enum `Title` dibuat dalam file model `services`, dan path impor diselesaikan secara otomatis:

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

Anda dapat menambahkan metode ke tipe enum di Go. Metode ini tidak memengaruhi pembuatan binding, tetapi menyediakan fungsionalitas sisi server yang berguna:

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

Enum yang dihasilkan tetap sama, terlepas dari ada atau tidaknya metode Go pada tipe tersebut.

## Komentar dan Dokumentasi

Generator mempertahankan komentar Go sebagai JSDoc dalam keluaran yang dihasilkan:

- **Komentar tipe** menjadi komentar dokumentasi enum
- **Komentar grup const** menjadi pemisah bagian
- **Komentar setiap const** menjadi komentar dokumentasi anggota
- **Komentar inline** dipertahankan jika memungkinkan

Artinya, IDE Anda akan menampilkan dokumentasi untuk nilai enum saat kursor diarahkan ke nilai tersebut.

## Tipe Dasar yang Didukung

Generator binding mendukung enum dengan tipe dasar Go berikut:

| Tipe Go | Dapat Digunakan sebagai Enum |
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

Hal berikut **tidak** didukung untuk pembuatan enum:

- **Tipe generik**—Parameter tipe menghalangi pendeteksian konstanta
- **Tipe dengan `json.Marshaler` atau `encoding.TextMarshaler`** kustom—Serialisasi kustom berarti nilai yang dihasilkan mungkin tidak cocok dengan perilaku runtime, sehingga generator melewati tipe ini
- **Konstanta yang nilainya tidak dapat dievaluasi secara statis atau direpresentasikan** — Konstanta harus memiliki nilai yang diketahui dan dapat direpresentasikan dalam tipe dasarnya. Pola `iota` standar dapat digunakan karena compiler mengubahnya menjadi nilai konkret
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

[Pelajari Selengkapnya →](/features/bindings/models/)

---
🚀 Binding Metode
Hubungkan metode Go ke frontend.

[Pelajari Selengkapnya →](/features/bindings/methods/)

---
⚙ Binding Tingkat Lanjut
Direktif, injeksi kode, dan ID khusus.

[Pelajari Selengkapnya →](/features/bindings/advanced/)

---
✓ Praktik Terbaik
Pola desain binding.

[Pelajari Selengkapnya →](/features/bindings/best-practices/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh binding](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/binding).
