---
title: "Cara Kerja Wails"
description: "Memahami arsitektur Wails dan caranya mencapai performa native"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails adalah framework untuk membangun aplikasi desktop menggunakan **Go untuk backend** dan **teknologi web untuk frontend**. Namun, tidak seperti Electron, Wails tidak menyertakan browser—Wails menggunakan **WebView native milik sistem operasi**.

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
    frontend."Render Site via OS-native WebView": Merender Situs melalui WebView Native OS
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
| **Ukuran Paket** | ~15MB | ~150MB |
| **Memori** | ~10MB | ~100MB+ |
| **Waktu Mulai** | &lt;0.5 dtk | 2-3 dtk |

## Komponen Inti

### 1. WebView Native

Wails menggunakan mesin rendering web bawaan sistem operasi:

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Berbasis Chromium (sama seperti browser Edge)
- Sudah terinstal di Windows 10/11
- Pembaruan otomatis melalui Windows Update
- Dukungan penuh untuk standar web modern

[macOS]
**WebKit** (mesin rendering Safari)

- Terintegrasi dalam macOS
- Mesin yang sama dengan browser Safari
- Performa dan masa pakai baterai yang sangat baik
- Dukungan penuh untuk standar web modern

[Linux]
**WebKitGTK** (port GTK dari WebKit)

- Diinstal melalui pengelola paket
- Mesin yang sama dengan GNOME Web (Epiphany)
- Dukungan standar yang baik
- Ringan dan berperforma tinggi

@end

**Mengapa hal ini penting:**

- **Tidak menyertakan browser** → Ukuran aplikasi lebih kecil
- **Native OS** → Integrasi dan performa lebih baik
- **Pembaruan otomatis** → Patch keamanan dari pembaruan OS
- **Rendering yang familier** → Sama seperti browser sistem

### 2. Bridge Wails

Bridge adalah inti Wails—komponen ini memungkinkan **komunikasi langsung** antara Go dan JavaScript.

```d2
direction: down

Frontend: Frontend (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Bridge Wails {
  Encoder: Encoder JSON {
    shape: rectangle
  }

  Router: Router Metode {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Decoder JSON {
    shape: rectangle
  }
}

Backend: Backend (Go) {
  Services: Layanan Terdaftar {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Panggil metode Go\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Enkode ke JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. Arahkan ke layanan\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Kembalikan hasil\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Dekode ke JS\nPromise diselesaikan"
```

**Cara kerjanya:**

1. **Frontend memanggil metode Go** (melalui binding yang dibuat otomatis)
2. **Bridge mengenkode panggilan** ke JSON (nama metode + argumen)
3. **Router menemukan metode Go** dalam layanan terdaftar
4. **Metode Go dijalankan** dan mengembalikan nilai
5. **Bridge mendekode hasil** dan mengirimkannya kembali ke frontend
6. **Promise diselesaikan** di JavaScript dengan hasil tersebut

**Karakteristik performa:**

- **Dalam memori**: Tanpa overhead jaringan dan tanpa HTTP
- **Tanpa penyalinan** jika memungkinkan (untuk data berukuran besar)
- **Asinkron secara default**: Tidak memblokir di kedua sisi
- **Aman secara tipe**: Definisi TypeScript dibuat secara otomatis

### 3. Sistem Layanan

Layanan adalah cara yang direkomendasikan untuk menyediakan fungsionalitas Go bagi frontend.

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

- Wails **memindai struct Anda** saat dimulai
- **Method yang diekspor** dapat dipanggil dari frontend
- **Informasi tipe** diekstrak untuk binding TypeScript
- **Penanganan error** berlangsung otomatis (error Go → exception JS)

**Binding TypeScript yang dihasilkan:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Mengapa menggunakan layanan?**

- **Aman terhadap tipe**: Dukungan penuh untuk TypeScript
- **Penemuan otomatis**: Tidak perlu mendaftarkan method secara manual
- **Terorganisasi**: Kelompokkan fungsionalitas yang berkaitan
- **Dapat diuji**: Layanan hanyalah struct Go

[Pelajari layanan lebih lanjut →](/features/bindings/services/)

### 4. Sistem Event

Event memungkinkan **komunikasi pub/sub** antarkomponen.

```d2
direction: left

Wails Event System: Sistem Event Wails {
  shape: sequence_diagram

  window1: Jendela 1
  window2: Jendela 2
  backend: Backend Go

  Event Driver: Penggerak Event {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "Berlangganan event 'data-updated'"
    window2."Subscribe to 'data-updated' events": "Berlangganan event 'data-updated'"
    backend.a."App Emit('data-updated', data)": "Aplikasi Emit('data-updated', data)"
    backend.a -> window1.a: Bus Event JSON
    backend.a -> window2: Bus Event JSON
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

[Pelajari event lebih lanjut →](/features/events/system/)

## Siklus Hidup Aplikasi

Memahami siklus hidup membantu Anda mengetahui kapan harus menginisialisasi dan membersihkan sumber daya.

```d2
direction: down

Start: Aplikasi Dimulai {
  shape: oval
  style.fill: "#10B981"
}

Init: Inisialisasi {
  Create: Buat Aplikasi {
    shape: rectangle
  }

  Register: Daftarkan Layanan {
    shape: rectangle
  }

  Setup: Siapkan Jendela/Menu {
    shape: rectangle
  }
}

Run: Loop Event {
  Events: Proses Event {
    shape: rectangle
  }

  Messages: Tangani Pesan {
    shape: rectangle
  }

  Render: Perbarui UI {
    shape: rectangle
  }
}

Shutdown: Penonaktifan {
  Cleanup: Bersihkan Sumber Daya {
    shape: rectangle
  }

  Save: Simpan Status {
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
Run.Render -> Run.Events: Loop
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

Tidak ada field `OnStartup` pada `application.Options`. Pekerjaan saat startup harus ditempatkan dalam `ServiceStartup(ctx, options)` milik layanan, dalam callback yang didaftarkan melalui `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`, atau cukup sebelum `app.Run()`.

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
  AnalyseGo: Analisis Kode Go {
    shape: rectangle
  }

  GenerateBindings: Hasilkan Binding {
    shape: rectangle
  }

  BuildFrontend: Build Frontend {
    shape: rectangle
  }

  CompileGo: Kompilasi Go {
    shape: rectangle
  }

  Embed: Sematkan Aset {
    shape: rectangle
  }
}

Output: Output {
  Binary: "Biner Native\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Ekstrak tipe
Build.GenerateBindings -> Source.Frontend: Binding TypeScript
Source.Frontend -> Build.BuildFrontend: Kompilasi (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Aset yang dibundel
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Langkah-langkah build:**

1. **Analisis kode Go**
  - Pindai layanan untuk menemukan method yang diekspor
  - Ekstrak tipe parameter dan nilai kembalian
  - Hasilkan signature method


2. **Hasilkan binding TypeScript**
  - Buat file `.ts` untuk setiap layanan
  - Sertakan definisi tipe lengkap
  - Tambahkan komentar JSDoc


3. **Build frontend**
  - Jalankan bundler Anda (Vite, webpack, dan sebagainya)
  - Minifikasi dan optimalkan
  - Tulis output ke `frontend/dist/`


4. **Kompilasi Go**
  - Kompilasi dengan optimasi (`-ldflags="-s -w"`)
  - Sertakan metadata build
  - Kompilasi khusus platform


5. **Sematkan aset**
  - Sematkan file frontend ke dalam biner Go
  - Kompres aset
  - Buat satu file executable


**Hasil:** Satu file executable native dengan semua komponen disematkan di dalamnya.

[Pelajari proses build lebih lanjut →](/guides/build/building/)

## Pengembangan vs Produksi

Wails berperilaku berbeda dalam lingkungan pengembangan dan produksi:

@tabs{sync-key="mode"}
[Pengembangan (wails3 dev)]
**Karakteristik:**

- **Hot reload**: Perubahan frontend langsung dimuat ulang
- **Peta sumber**: Lakukan debug menggunakan sumber asli
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

WailsApp -> DevServer: Proksikan permintaan
DevServer -> WebView: Sajikan dengan HMR
WebView -> WailsApp: Panggil metode Go
```

**Manfaat:**

- Umpan balik langsung atas perubahan
- Kemampuan debug lengkap
- Iterasi lebih cepat

[Produksi (wails3 build)]
**Karakteristik:**

- **Aset tertanam**: Frontend disertakan dalam biner
- **Dioptimalkan**: Diminifikasi dan dikompresi
- **Tanpa DevTools**: Dinonaktifkan secara default
- **Pencatatan log minimal**: Hanya kesalahan
- **Satu berkas**: Semuanya berada dalam satu berkas yang dapat dieksekusi

**Cara kerjanya:**

```d2
direction: right

Binary: "Biner Tunggal\n(myapp.exe)" {
  GoCode: Go Terkompilasi {
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

Binary.Assets -> WebView: Sajikan dari memori
WebView -> Binary.GoCode: Panggil metode Go
```

**Manfaat:**

- Distribusi dalam satu berkas
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
  - Biasanya 5-10MB untuk aplikasi sederhana


2. **Memori WebView**
  - DOM, heap JavaScript, CSS
  - Dikelola oleh mesin WebView
  - Biasanya 10-20MB untuk aplikasi sederhana


3. **Memori Bridge**
  - Buffer pesan untuk komunikasi
  - Overhead minimal (<1MB)
  - Tanpa penyalinan untuk data berukuran besar jika memungkinkan


**Kiat pengoptimalan:**

- **Hindari transfer data berukuran besar**: Teruskan ID, lalu ambil detail saat diperlukan
- **Gunakan event untuk pembaruan**: Jangan lakukan polling dari frontend
- **Lakukan streaming berkas besar**: Jangan muat seluruh berkas ke dalam memori
- **Bersihkan listener**: Hapus event listener setelah selesai

[Pelajari performa lebih lanjut →](/guides/performance/)

## Model Keamanan

Wails menyediakan arsitektur yang aman secara default:

```d2
direction: down

Frontend: Frontend (Tidak Tepercaya) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails Bridge (Validasi) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Backend (Tepercaya) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Panggil metode
Bridge -> Bridge: "Validasi:\n- Apakah metode tersedia?\n- Apakah tipenya benar?\n- Apakah akses diizinkan?"
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
  - Argumen diperiksa terhadap tipe Go
  - Tipe yang tidak valid ditolak
  - Mencegah serangan injeksi


3. **Tanpa eval()**
  - Frontend tidak dapat menjalankan kode Go arbitrer
  - Hanya metode yang telah ditentukan yang dapat dipanggil
  - Tidak ada eksekusi kode dinamis


4. **Isolasi konteks**
  - Setiap jendela memiliki konteksnya sendiri
  - Layanan dapat memeriksa konteks pemanggil
  - Izin per jendela dapat diterapkan


**Praktik terbaik:**

- **Validasi input pengguna** di Go (jangan percayai frontend)
- **Gunakan konteks** untuk autentikasi/otorisasi
- **Sanitasi jalur berkas** sebelum operasi berkas
- **Batasi laju** operasi yang mahal

[Pelajari keamanan lebih lanjut →](/guides/security/)

## Langkah Berikutnya

**Siklus Hidup Aplikasi** - Pahami startup, shutdown, dan hook siklus hidup [Pelajari Lebih Lanjut →](/concepts/lifecycle/)

**Bridge Go-Frontend** - Pelajari secara mendalam cara kerja bridge [Pelajari Lebih Lanjut →](/concepts/bridge/)

**Sistem Build** - Pahami cara Wails membangun aplikasi Anda [Pelajari Lebih Lanjut →](/concepts/build-system/)

**Mulai Membangun** - Terapkan hal yang telah Anda pelajari dalam tutorial [Tutorial →](/tutorials/03-notes-vanilla/)

---

**Ada pertanyaan tentang arsitektur?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [referensi API](/reference/overview/).
