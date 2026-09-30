---
title: "Memperluas Wails"
description: "Panduan praktis untuk menambahkan fitur dan platform baru ke Wails v3"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails dirancang agar **mudah diutak-atik**.
>
> Setiap subsistem utama berada dalam kode Go yang dapat Anda baca, ubah, dan distribusikan.
>
> Halaman ini menunjukkan *dari mana* harus memulai dan *cara* mempertahankan kompatibilitas lintas platform saat Anda:

- Menambahkan **Layanan** (notifikasi, penyimpanan KV, IPC khusus, …)
- Membuat **perintah CLI baru** (`wails3 <foo>`)
- Memperluas **Runtime** (API jendela, dialog, peristiwa)
- Memperkenalkan **kapabilitas platform** (Wayland, …)
- Mempertahankan **kompatibilitas lintas platform** tanpa kewalahan oleh tag `//go:build`

---

## 1. Menambahkan Layanan

"Layanan" di v3 adalah tipe Go yang disediakan pengguna, didaftarkan melalui `application.Options.Services`, dan diekspos ke JS melalui binding yang dihasilkan. Basis kode v3 menyediakan:

- `internal/service/` — kerangka awal untuk `wails3 generate service`:
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — layanan siap pakai yang dapat Anda daftarkan sekarang (notifications, kvstore, sqlite, log, fileserver, dock, …).

File generator dan CLI yang dirujuk draf lama sebagai `internal/service/template/template.go` dan `internal/generator/collect/services.go` tidak ada — pembuat kerangka awalnya adalah `internal/service/service.go` (titik masuk `service.Install`), sedangkan metadata binding untuk layanan dikumpulkan di `internal/generator/collect/service.go`.

### 1.1 Definisikan Layanan

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 Implementasikan antarmuka siklus hidup (opsional)

Layanan dapat secara opsional memenuhi antarmuka berikut (dari `pkg/application`):

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **Penting:** `ServiceShutdown` **tidak menerima argumen**. Metode dengan
>
> signature `ServiceShutdown(ctx context.Context) error` **tidak** memenuhi
>
> antarmuka tersebut dan diam-diam tidak akan pernah dipanggil.

### 1.3 Daftarkan layanan ke aplikasi

Tidak ada pemanggilan `services.Register(...)` global. Layanan didaftarkan saat runtime melalui `application.Options.Services`:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

Setelah didaftarkan, `wails3 generate bindings` menghasilkan modul ES di bawah `frontend/bindings/<your import path>/...` yang membungkus metode yang diekspor.

### 1.4 Memanggil dari JS

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

Tidak ada global `window.backend.*` di v3 — pemanggilan dilakukan melalui modul ES yang dihasilkan, yang kemudian memanggil `Call.ByID(...)` dari `/wails/runtime.js`.

---

## 2. Menulis Perintah CLI Baru

CLI v3 menggunakan **`github.com/leaanthony/clir`** (bukan cobra). Pengkabelannya berada di `v3/cmd/wails3/main.go`:

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

Tidak ada pendaftaran otomatis berbasis `init()`. Tambahkan subperintah baru Anda ke `cmd/wails3/main.go` beserta fungsi pendukung di `internal/commands/` (dan struct flag di bawah `internal/flags/` jika menerima opsi). Bangun ulang CLI:

```
cd v3
go install ./cmd/wails3
wails3 hello
```

Jika perintah Anda memerlukan penghubung Taskfile, gunakan kembali helper di `internal/commands/task_wrapper.go` (`wrapTask("yourtask", args)`).

---

## 3. Memodifikasi Runtime

Alasan umum:

- Fitur jendela baru (`SetOpacity`, `Shake`, …)
- Dialog tambahan (`ColorPicker`)
- API tingkat sistem (kecerahan layar)

### 3.1 API Publik

Tambahkan metode ke `pkg/application/webview_window.go` (antarmukanya berada di `window.go`):

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

Gunakan helper `InvokeSync`/`InvokeAsync` yang sudah ada untuk memastikan pemanggilan berjalan di thread utama.

### 3.2 Pemroses Pesan

Jika JS perlu memanggil metode baru tersebut, perluas file `pkg/application/messageprocessor_*.go` yang sesuai. Pemroses pesan menggunakan metode berbasis switch pada `MessageProcessor`, bukan pemanggilan `register(...)` global:

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

Dalam basis kode, **tidak** ada file `messageprocessor_window_opacity.go` maupun pola `register(MsgSetOpacity, ...)` berbasis `init()`.

### 3.3 Implementasi Platform

Tambahkan implementasi ke setiap file khusus OS di bawah `pkg/application/`:

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

Jika suatu platform tidak dapat mendukung fitur tersebut, tulis stub no-op. Framework ini tidak memiliki sentinel `ErrCapability` — tampilkan dukungan melalui dokumentasi dan, jika diperlukan, melalui field boolean yang relevan pada `Options` / struct opsi khusus platform.

### 3.4 Flag Kapabilitas (opsional)

Paket `internal/capabilities/` tersedia untuk mendeklarasikan kumpulan kapabilitas per platform. Tidak ada API publik `application.HasCapability` / `application.CapOpacity`. Jika Anda menginginkan kapabilitas yang dapat diperiksa saat runtime, tambahkan kapabilitas tersebut di bawah `internal/capabilities/` dan ekspos getter bertipe dari `pkg/application`.

---

## 4. Menambahkan Kapabilitas Platform Baru

Contoh: dukungan Wayland opsional di Linux.

1. Pisahkan file `pkg/application/*_linux.go` yang relevan menjadi `*_linux_x11.go` (`//go:build linux && !wayland`) dan `*_linux_wayland.go` (`//go:build linux && wayland`).
2. Minta pengguna mengaktifkannya dengan `wails3 build --tags wayland`. Teruskan tag tambahan melalui mekanisme `EXTRA_TAGS` yang sudah ada di `internal/commands/task_wrapper.go`. Tidak ada flag `--tags wayland` tingkat `dev` — `wails3 dev` hanya menerima `--config`, `--port`, dan `-s`.
3. Perbarui dokumentasi dan README khusus platform apa pun di bawah `pkg/application/`.

> Pertahankan tag build default seminimal mungkin; gunakan tag yang harus diaktifkan secara khusus hanya untuk fitur khusus.

---

## 5. Daftar Periksa Kompatibilitas Lintas Platform

| ✅ Langkah | Alasan |
| --- | --- |
| Sediakan **setiap** metode publik di semua berkas platform (meskipun hanya berupa stub) | Menjaga agar build berhasil di setiap OS |
| Dokumentasikan penurunan fungsi secara wajar untuk setiap OS | Aplikasi dapat membuat percabangan berdasarkan `runtime.GOOS` tanpa kesalahan tersembunyi |
| Utamakan **Go murni**, gunakan Cgo hanya jika diperlukan | Menyederhanakan kompilasi silang (Linux sudah menanggung beban Cgo) |
| Jalankan `task test:cli`, `task test:generator`, dan `task test:templates` | Mereproduksi CI secara lokal |
| Dokumentasikan tag build baru dalam dokumentasi kontributor / README templat | Pengguna harus mengetahui fitur yang perlu diaktifkan secara eksplisit |

---

## 6. Build Debug & Kecepatan Iterasi

- Gunakan `Options.LogLevel = slog.LevelDebug` (`Options.Logger = slog.Default()`) untuk menampilkan aktivitas runtime secara mendetail. Tidak ada variabel lingkungan `WAILS_LOG_LEVEL`.
- Flag `wails3 dev` adalah `--config`, `--port`, dan `-s`. Tidak ada flag `-race` atau `-verbose` — jalankan race detector dengan `go test -race ./...` atau dengan melakukan `go build -race` pada aplikasi Anda lalu menjalankannya secara langsung.
- Panduan pengujian race / Cgo tersedia di `v3/TESTING.md` (draf lama merujuk ke `pkg/application/RACE.md`, yang tidak ada).

---

## 7. Kontribusi ke Upstream

1. Untuk fungsionalitas baru atau perubahan perilaku publik, buka PR draf **WEP (Wails Enhancement Proposal)** guna mendiskusikan gagasan dan desainnya. Gunakan issue hanya untuk bug yang dapat direproduksi atau masalah dokumentasi.
2. Ikuti teknik di atas untuk mengimplementasikannya.
3. Tambahkan:
  - Pengujian unit (`*_test.go`)
  - Dokumentasi (berkas ini atau halaman `docs/...` yang relevan)
  - Pengujian regresi di bawah `internal/generator/testcases/` jika Anda mengubah generator binding

4. Jalankan `task precommit` beserta target `task test:*` yang relevan secara lokal sebelum melakukan push.

---

### Tautan Cepat

| Area | Lokasi |
| --- | --- |
| Layanan bawaan | `pkg/services/` |
| Pembuat kerangka layanan | `internal/service/` |
| Pendaftaran perintah CLI | `v3/cmd/wails3/main.go` |
| Implementasi perintah CLI | `internal/commands/` |
| Runtime per OS | `pkg/application/*_{darwin,linux,windows}.go` |
| Deklarasi kapabilitas | `internal/capabilities/` |
| DSL Taskfile | `v3/Taskfile.yaml` |
| Generator konstanta peristiwa | `v3/tasks/events/generate.go` |

---

Kini Anda memiliki **peta jalan** untuk menyesuaikan Wails sesuka Anda — tambahkan layanan, sisipkan keajaiban CLI, utak-atik runtime, atau hadirkan fitur OS yang benar-benar baru. Selamat mengembangkan!
