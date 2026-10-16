---
title: "Cara Kerja Wails"
description: "Memahami arsitektur Wails dan caranya mencapai performa native"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails adalah kerangka kerja untuk membangun aplikasi desktop dengan **Go sebagai backend** dan **teknologi web sebagai frontend**. Namun, tidak seperti Electron, Wails tidak menyertakan browser—Wails menggunakan **WebView native milik sistem operasi**.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Aplikasi Wails

  frontend: Frontend
  backend: Backend Go
  os: Sistem Operasi

  Initialisation: Inisialisasi {
    shape: sequence_diagram
    backend."Serves Static Web App": Menyajikan Aplikasi Web Statis
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": Merender Situs melalui WebView Bawaan OS
  }
  Regular Communication: Komunikasi Reguler {
    shape: sequence_diagram
    frontend."Make API-style call": Melakukan Panggilan Bergaya API
    frontend -> backend.a: JSON
    backend.a."Service processes request": Layanan memproses permintaan
    backend.a -> os: Memanggil API Sistem
    backend.a."Generate Response": Menghasilkan Respons
    backend.a -> frontend: JSON
    frontend."Process response": Memproses respons
  }
  backend.a.label: a
}
```

**Perbedaan utama dari Electron:**

| Aspek | Wails | Electron |
| --- | --- | --- |
| **Browser** | WebView yang disediakan OS | Chromium yang disertakan (~100MB) |
| **Backend** | Go (dikompilasi) | Node.js (diinterpretasikan) |
| **Komunikasi** | Bridge dalam memori | IPC (antarproses) |
| **Ukuran Bundel** | ~15MB | ~150MB |
| **Memori** | ~10MB | ~100MB+ |
| **Waktu mulai** | &lt;0.5s | 2-3s |

## Komponen Inti

### 1. WebView Native

Wails menggunakan mesin perender web bawaan sistem operasi:

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Berbasis Chromium (sama seperti browser Edge)
- Sudah terinstal di Windows 10/11
- Pembaruan otomatis melalui Windows Update
- Dukungan penuh untuk standar web modern

[macOS]
**WebKit** (mesin perender Safari)

- Terintegrasi dalam macOS
- Mesin yang sama dengan browser Safari
- Performa dan daya tahan baterai yang sangat baik
- Dukungan penuh untuk standar web modern

[Linux]
**WebKitGTK** (port GTK dari WebKit)

- Diinstal melalui pengelola paket
- Mesin yang sama dengan GNOME Web (Epiphany)
- Dukungan standar yang baik
- Ringan dan berperforma tinggi

@end

**Mengapa hal ini penting:**

- **Tanpa browser yang disertakan** → Ukuran aplikasi lebih kecil
- **Native pada OS** → Integrasi dan performa lebih baik
- **Pembaruan otomatis** → Patch keamanan dari pembaruan OS
- **Perenderan yang familier** → Sama seperti browser sistem

### 2. Bridge Wails

Bridge adalah inti Wails—komponen ini memungkinkan **komunikasi langsung** antara Go dan JavaScript.

```d2
direction: down

Frontend: Frontend (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Bridge Wails {
  Encoder: Enkoder JSON {
    shape: rectangle
  }

  Router: Router Metode {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Dekoder JSON {
    shape: rectangle
  }
}

Backend: Backend (Go) {
  Services: Layanan Terdaftar {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Memanggil metode Go\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Mengodekan ke JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. Merutekan ke layanan\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Mengembalikan hasil\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Mendekodekan ke JS\nPromise diselesaikan"
```

**Cara kerjanya:**

1. **Frontend memanggil metode Go** (melalui binding yang dibuat otomatis)
2. **Bridge mengodekan pemanggilan** ke JSON (nama metode + argumen)
3. **Router menemukan metode Go** dalam layanan yang terdaftar
4. **Metode Go dijalankan** dan mengembalikan nilai
5. **Bridge mendekode hasil** dan mengirimkannya kembali ke frontend
6. **Promise diselesaikan** di JavaScript dengan hasil tersebut

**Karakteristik performa:**

- **Dalam memori**: Tanpa overhead jaringan, tanpa HTTP
- **Tanpa penyalinan** jika memungkinkan (untuk data berukuran besar)
- **Asinkron secara default**: Tidak memblokir di kedua sisi
- **Aman secara tipe**: Definisi TypeScript dibuat otomatis

### 3. Sistem Layanan

Layanan adalah cara yang direkomendasikan untuk mengekspos fungsionalitas Go ke frontend.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Penemuan layanan:**

- Wails **memindai struct Anda** saat mulai dijalankan
- **Metode yang diekspor** dapat dipanggil dari frontend
- **Informasi tipe** diekstrak untuk binding TypeScript
- **Penanganan kesalahan** berlangsung otomatis (error Go → exception JS)

**Binding TypeScript yang dihasilkan:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Mengapa menggunakan layanan?**

- **Aman terhadap tipe**: Dukungan penuh untuk TypeScript
- **Penemuan otomatis**: Tidak perlu mendaftarkan metode secara manual
- **Terorganisasi**: Kelompokkan fungsionalitas yang terkait
- **Dapat diuji**: Layanan hanyalah struct Go

[Pelajari layanan lebih lanjut →](/features/bindings/services/)

### 4. Sistem Peristiwa

Peristiwa memungkinkan **komunikasi publikasi/langganan** antarkomponen.

```d2
direction: left

Wails Event System: Sistem Peristiwa Wails {
  shape: sequence_diagram

  window1: Jendela 1
  window2: Jendela 2
  backend: Backend Go

  Event Driver: Penggerak Peristiwa {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "Berlangganan peristiwa 'data-updated'"
    window2."Subscribe to 'data-updated' events": "Berlangganan peristiwa 'data-updated'"
    backend.a."App Emit('data-updated', data)": "Aplikasi Emit('data-updated', data)"
    backend.a -> window1.a: Bus Peristiwa JSON
    backend.a -> window2: Bus Peristiwa JSON
    window1.a."Subscriber processes On('data-updated', handler)": "Pelanggan memproses On('data-updated', handler)"
    window2."Subscriber processes On('data-updated', handler)": "Pelanggan memproses On('data-updated', handler)"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**Kasus penggunaan:**

- **Komunikasi antarjendela**: Satu jendela memberi tahu jendela lainnya
- **Tugas latar belakang**: Layanan Go memberi tahu UI tentang progres
- **Sinkronisasi status**: Jaga agar beberapa jendela tetap sinkron
- **Kopling longgar**: Komponen tidak memerlukan referensi langsung

**Contoh:**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[Pelajari peristiwa lebih lanjut →](/features/events/system/)

## Siklus Hidup Aplikasi

Memahami siklus hidup membantu Anda mengetahui kapan harus menginisialisasi dan membersihkan sumber daya.

```d2
direction: down

Start: Aplikasi Dimulai {
  shape: oval
  style.fill: "#10B981"
}

Init: Inisialisasi {
  Create: Membuat Aplikasi {
    shape: rectangle
  }

  Register: Mendaftarkan Layanan {
    shape: rectangle
  }

  Setup: Menyiapkan Jendela/Menu {
    shape: rectangle
  }
}

Run: Perulangan Peristiwa {
  Events: Memproses Peristiwa {
    shape: rectangle
  }

  Messages: Menangani Pesan {
    shape: rectangle
  }

  Render: Memperbarui UI {
    shape: rectangle
  }
}

Shutdown: Penghentian {
  Cleanup: Membersihkan Sumber Daya {
    shape: rectangle
  }

  Save: Menyimpan Status {
    shape: rectangle
  }
}

End: Aplikasi Berakhir {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: Ulangi
Run.Events -> Shutdown.Cleanup: Sinyal keluar
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Hook siklus hidup:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

Tidak ada bidang `OnStartup` pada `application.Options`. Pekerjaan saat startup harus ditempatkan dalam `ServiceStartup(ctx, options)` milik layanan, dalam callback yang didaftarkan melalui `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`, atau cukup dijalankan sebelum `app.Run()`.

[Pelajari siklus hidup lebih lanjut →](/concepts/lifecycle/)

## Proses Build

Memahami cara Wails mem-build aplikasi Anda:

```d2
direction: down

Source: Kode Sumber {
  Go: "Kode Go\n(main.go, layanan)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "Kode Frontend\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: Proses Build {
  AnalyseGo: Menganalisis Kode Go {
    shape: rectangle
  }

  GenerateBindings: Menghasilkan Binding {
    shape: rectangle
  }

  BuildFrontend: Mem-build Frontend {
    shape: rectangle
  }

  CompileGo: Mengompilasi Go {
    shape: rectangle
  }

  Embed: Menyematkan Aset {
    shape: rectangle
  }
}

Output: Keluaran {
  Binary: "Biner Native\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Mengekstrak tipe
Build.GenerateBindings -> Source.Frontend: Binding TypeScript
Source.Frontend -> Build.BuildFrontend: Mengompilasi (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Aset yang Dibundel
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Langkah-langkah build:**

1. **Analisis kode Go**
  - Pindai metode yang diekspor dalam layanan
  - Ekstrak tipe parameter dan nilai kembalian
  - Buat signature metode


2. **Buat binding TypeScript**
  - Buat file `.ts` untuk setiap layanan
  - Sertakan definisi tipe lengkap
  - Tambahkan komentar JSDoc


3. **Bangun frontend**
  - Jalankan bundler Anda (Vite, webpack, dan sebagainya)
  - Minifikasi dan optimalkan
  - Keluarkan hasil ke `frontend/dist/`


4. **Kompilasi Go**
  - Kompilasi dengan pengoptimalan (`-ldflags="-s -w"`)
  - Sertakan metadata build
  - Kompilasi khusus platform


5. **Sematkan aset**
  - Sematkan file frontend ke dalam biner Go
  - Kompres aset
  - Buat satu file yang dapat dieksekusi


**Hasil:** Satu file native yang dapat dieksekusi dengan semua komponen tertanam di dalamnya.

[Pelajari proses build lebih lanjut →](/guides/build/building/)

## Pengembangan vs Produksi

Wails berperilaku berbeda dalam lingkungan pengembangan dan produksi:

@tabs{sync-key="mode"}
[Pengembangan (wails3 dev)]
**Karakteristik:**

- **Hot reload**: Perubahan frontend langsung dimuat ulang
- **Source map**: Lakukan debugging menggunakan kode sumber asli
- **DevTools**: DevTools browser tersedia
- **Pencatatan log**: Pencatatan log mendetail diaktifkan
- **Frontend eksternal**: Disajikan dari server pengembangan (Vite)

**Cara kerjanya:**

```d2
direction: right

WailsApp: Aplikasi Wails {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Server Pengembangan Vite\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: Mem-proxy permintaan
DevServer -> WebView: Menyajikan dengan HMR
WebView -> WailsApp: Memanggil metode Go
```

**Manfaat:**

- Umpan balik seketika atas perubahan
- Kemampuan debugging lengkap
- Iterasi lebih cepat

[Produksi (wails3 build)]
**Karakteristik:**

- **Aset tertanam**: Frontend disertakan dalam biner saat build
- **Dioptimalkan**: Diminifikasi dan dikompresi
- **Tanpa DevTools**: Dinonaktifkan secara default
- **Pencatatan log minimal**: Hanya kesalahan
- **Satu file**: Semuanya terdapat dalam satu file yang dapat dieksekusi

**Cara kerjanya:**

```d2
direction: right

Binary: "Biner Tunggal\n(myapp.exe)" {
  GoCode: Go yang Dikompilasi {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "Aset Tertanam\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: Menyajikan dari memori
WebView -> Binary.GoCode: Memanggil metode Go
```

**Manfaat:**

- Distribusi dalam satu file
- Ukuran lebih kecil (diminifikasi)
- Performa lebih baik
- Tanpa dependensi eksternal

@end

## Model Memori

Memahami penggunaan memori membantu Anda membangun aplikasi yang efisien.

**Wilayah memori:**

1. **Heap Go**
  - Layanan dan status aplikasi Anda
  - Dikelola oleh pengumpul sampah Go
  - Biasanya 5-10 MB untuk aplikasi sederhana


2. **Memori WebView**
  - DOM, heap JavaScript, CSS
  - Dikelola oleh mesin WebView
  - Biasanya 10-20 MB untuk aplikasi sederhana


3. **Memori Bridge**
  - Buffer pesan untuk komunikasi
  - Overhead minimal (<1 MB)
  - Tanpa penyalinan untuk data berukuran besar jika memungkinkan


**Tips pengoptimalan:**

- **Hindari transfer data berukuran besar**: Teruskan ID, ambil detail sesuai kebutuhan
- **Gunakan event untuk pembaruan**: Jangan lakukan polling dari frontend
- **Streaming file berukuran besar**: Jangan muat seluruhnya ke dalam memori
- **Bersihkan listener**: Hapus event listener setelah selesai

[Pelajari lebih lanjut tentang performa →](/guides/performance/)

## Model Keamanan

Wails menyediakan arsitektur yang aman secara bawaan:

```d2
direction: down

Frontend: Frontend (Tidak Tepercaya) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Bridge Wails (Validasi) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Backend (Tepercaya) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Memanggil metode
Bridge -> Bridge: "Validasi:\n- Metode tersedia?\n- Tipe sudah benar?\n- Akses diizinkan?"
Bridge -> Backend: Jalankan jika valid
Backend -> Bridge: Kembalikan hasil
Bridge -> Frontend: Kirim respons
```

**Fitur keamanan:**

1. **Daftar putih metode**
  - Hanya metode yang diekspor yang dapat dipanggil
  - Metode privat tidak dapat diakses
  - Pendaftaran layanan secara eksplisit diperlukan


2. **Validasi tipe**
  - Argumen diperiksa berdasarkan tipe Go
  - Tipe yang tidak valid ditolak
  - Mencegah serangan injeksi


3. **Tanpa eval()**
  - Frontend tidak dapat menjalankan kode Go arbitrer
  - Hanya metode yang telah ditentukan sebelumnya yang dapat dipanggil
  - Tidak ada eksekusi kode dinamis


4. **Isolasi konteks**
  - Setiap jendela memiliki konteksnya sendiri
  - Layanan dapat memeriksa konteks pemanggil
  - Izin per jendela dapat diterapkan


**Praktik terbaik:**

- **Validasi input pengguna** di Go (jangan percayai frontend)
- **Gunakan konteks** untuk autentikasi/otorisasi
- **Sanitasi jalur file** sebelum operasi file
- **Batasi laju** operasi yang mahal

[Pelajari lebih lanjut tentang keamanan →](/guides/security/)

## Langkah Berikutnya

**Siklus Hidup Aplikasi** - Pahami proses startup, shutdown, dan hook siklus hidup [Pelajari Lebih Lanjut →](/concepts/lifecycle/)

**Bridge Go-Frontend** - Pelajari secara mendalam cara kerja bridge [Pelajari Lebih Lanjut →](/concepts/bridge/)

**Sistem Build** - Pahami cara Wails membangun aplikasi Anda [Pelajari Lebih Lanjut →](/concepts/build-system/)

**Mulai Membangun** - Terapkan yang telah Anda pelajari dalam tutorial [Tutorial →](/tutorials/03-notes-vanilla/)

---

**Ada pertanyaan tentang arsitektur?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [referensi API](/reference/overview/).
