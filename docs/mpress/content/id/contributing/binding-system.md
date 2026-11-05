---
title: "Sistem Binding"
description: "Cara Wails v3 memungkinkan Go dan JavaScript saling memanggil tanpa kode boilerplate"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> "Binding" adalah **kontrak aman tipe** yang memungkinkan Anda menulis:

```go
msg, err := chatService.Send("Hello")
```

di Go *dan*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

di TypeScript **tanpa menulis kode penghubung IPC secara manual**. Dokumen ini menguraikan *cara* hal tersebut berlangsung, mulai dari **analisis statis** pada waktu build, melalui **pembuatan kode**, hingga **jembatan runtime** yang memindahkan byte melalui WebView.

> Lihat [`contributing/architecture/bindings`](/contributing/architecture/bindings/) untuk
>
> pembahasan mendalam dan otoritatif tentang pipeline generator — halaman ini merupakan
>
> ikhtisar yang berfokus pada kontributor.

---

## 1. Ikhtisar 30 Detik

| Tahap | Komponen | Output |
| --- | --- | --- |
| **Pengumpulan/Analisis** | `internal/generator/collect/`, `internal/generator/analyse.go` | Model dalam memori untuk layanan, metode, parameter, tipe nilai kembalian, dan model Go yang diekspor |
| **Pembuatan** | `internal/generator/render/templates/*.tmpl` (`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | Modul ES per layanan di bawah `frontend/bindings/<full Go import path>/...` |
| **Runtime** | `pkg/application/messageprocessor*.go` + runtime JS tertanam di bawah `internal/runtime/desktop/@wailsio/runtime/src/` (`calls.ts`, `events.ts`, …) | Pesan panggilan/peristiwa melalui jembatan native WebView |

Alur ini diatur oleh perintah `wails3 generate bindings`, yang menjalankan `generator.Generate` (didefinisikan dalam `internal/generator/generate.go`) pada sekumpulan paket Go.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. Analisis Statis

### Titik Masuk

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

Tahap kolektor menelusuri setiap paket yang dimuat dan mencatat:

- `collect.ServiceInfo` — satu untuk setiap struct Go yang diekspor dan di-binding.
- `collect.ServiceMethodInfo` / `collect.MethodInfo` — informasi signature per metode (nama, parameter, hasil, posisi error, receiver, dokumentasi).
- `collect.ModelInfo` / `collect.StructInfo` — dihasilkan sebagai model TS/JS.
- Komentar direktif seperti `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore`, `//wails:id <hex>` (lihat `internal/generator/collect/directive.go`).

Tipe yang tidak didukung menghasilkan error generator sehingga kesalahan terungkap pada waktu build, bukan saat runtime.

### Pengidentifikasi Model

Envelope panggilan runtime mengidentifikasi metode menggunakan **hash FNV-1a deterministik** dari nama lengkapnya (`pkg.Struct.Method`). Ini akan muncul sebagai `$Call.ByID(<numeric-id>, …)` dalam binding yang dihasilkan atau sebagai `$Call.ByName("pkg.Struct.Method", …)` ketika pembuatan dijalankan dengan `-names`.

---

## 3. Pembuatan Kode

### Templat

`internal/generator/render/templates/`:

| Templat | Tujuan |
| --- | --- |
| `service.js.tmpl` | Satu modul JS per layanan yang di-binding |
| `service.ts.tmpl` | Pendamping TypeScript (dengan `-ts`) |
| `models.js.tmpl` | Output kelas model (per paket) |
| `models.ts.tmpl` | Output `.d.ts` model (per paket) |
| `index.tmpl` | Ekspor ulang barrel `index.{js,ts}` per paket |
| `eventcreate.js.tmpl` / `eventdata.d.ts.tmpl` | Konstruktor peristiwa / deklarasi tipe payload |
| `newline.tmpl` | Penormal baris baru di akhir |

Output ditempatkan di bawah `frontend/bindings/<full Go import path>/...` — misalnya, layanan yang didefinisikan dalam `github.com/you/yourapp/services/chat` ditempatkan di bawah `frontend/bindings/github.com/you/yourapp/services/chat/`. Tidak ada direktori `frontend/src/wailsjs/` di v3.

### Output JavaScript

Binding yang dihasilkan berupa modul ES yang mengimpor helper runtime dari `/wails/runtime.js`:

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

Ketika pembuatan dijalankan dengan `-names`, `$Call.ByName("pkg.Struct.Method", ...)` yang dihasilkan sebagai gantinya — selalu **bernama lengkap**, tidak pernah hanya `"Method"`.

Kelas model yang dihasilkan menggunakan pola konstruktor `$$source` dengan nilai default `if (!("X" in $$source))` per field, nama field dalam tanda kutip, serta `static createFrom(...)` yang menjalankan `JSON.parse` pada input string.

### Sorotan Pemetaan Tipe

Diverifikasi terhadap `internal/generator/render/`:

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V` (`K` non-string) | `{ [_ in K]?: V }` (bukan `Map<K, V>`, bukan `Record<K, V>`) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string` (JSON ISO 8601) |
| `error` (posisi nilai kembalian) | promise ditolak |

### Catatan Refleksi

`pkg/application/bindings.go` **ditulis secara manual** dan menggunakan `reflect` untuk mengatur dispatch metode dari registry `BoundMethod`. Jangan menafsirkan klaim lama tentang "tanpa refleksi saat runtime" terlalu harfiah—generator menghindari refleksi, tetapi dispatcher runtime menggunakannya.

---

## 4. Protokol Pemanggilan Runtime

### Sisi JavaScript

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

Helper runtime berada di `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` (dispatch panggilan), `events.ts` (event), dan file terkait—tidak ada `invoke.ts` atau `errors.ts` dalam working tree ini. Envelope wire yang tepat dikodekan oleh `calls.ts` di sisi JS dan didekodekan oleh `pkg/application/messageprocessor_call.go` di sisi Go; periksa kedua file tersebut secara bersamaan saat men-debug bridge.

### Sisi Go

1. `pkg/application/messageprocessor_call.go` menerima pesan panggilan.
2. Mencari metode yang diikat berdasarkan ID atau nama di `pkg/application/bindings.go` (dikendalikan oleh `reflect`).
3. Memanggil metode yang diikat dan membuat serialisasi `{result, error}` untuk dikirim kembali ke JS.

### Pemetaan Error

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise` diselesaikan dengan hasil |
| `error != nil` | `Promise` ditolak dengan `Error` yang `message`-nya berisi string error Go |

---

## 5. Memanggil JavaScript dari Go

Generator binding bersifat satu arah (metode Go diekspos ke JS). Untuk komunikasi Go → JS, gunakan bus event atau jalankan JS di sebuah jendela:

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

Di sisi JS, berlanggananlah dengan `Events.On(name, cb)` dari `/wails/runtime.js`.

---

## 6. Ekstensi & Pemecahan Masalah

### Error Tipe yang Tidak Didukung

```
error: field "Client" uses unsupported type: chan struct{}
```

→ bungkus channel di balik API metode, atau tandai field dengan `//wails:internal` agar generator melewatinya.

### Binding Kedaluwarsa

Output yang dihasilkan akan ditimpa pada setiap `wails3 generate bindings` / `wails3 dev` / `wails3 build`. Jika IntelliSense IDE menampilkan stub kedaluwarsa, hapus `frontend/bindings/` dan jalankan kembali generator. Flag `-clean` (nilai default `true` dalam build saat ini) menghapus isi direktori binding sebelum setiap proses dijalankan.

### Kiat Performa

- Hindari streaming slice byte berukuran besar melalui bridge—sajikan melalui server aset sebagai gantinya.
- Gabungkan beberapa panggilan cepat ke dalam satu metode jika latensi penting.
- Utamakan receiver nilai untuk struct parameter berukuran kecil guna mengurangi alokasi.

---

## 7. Peta File Utama

| Aspek | File |
| --- | --- |
| Orkestrasi generator | `internal/generator/generate.go` |
| Pemeriksaan semantik | `internal/generator/analyse.go` |
| Pengumpulan (layanan, metode, model) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| Template rendering | `internal/generator/render/templates/*.tmpl` |
| Lokasi binding yang dihasilkan | `frontend/bindings/<full Go import path>/...` |
| Dispatcher sisi Go | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| Runtime JS | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

Simpan lembar referensi ringkas ini agar mudah diakses saat Anda melacak bug pada bridge.

---

## 8. Ringkasan

1. **Collector** memindai kode Go Anda → model semantik dalam memori.
2. **Templates** menghasilkan modul ES per layanan serta file model/indeks per paket.
3. **Message Processor** melakukan dispatch panggilan di sisi Go melalui registry binding.
4. **JS Runtime** membungkus semuanya dalam promise idiomatis yang mendukung pembatalan.

Semuanya tanpa perlu menulis satu baris pun kode boilerplate IPC. Itulah sistem binding Wails v3. Silakan mulai membuat binding!
