---
title: "Aplikasi Wails yang Memperbarui Diri Sendiri"
description: "Bangun aplikasi Wails v3 yang memperbarui dirinya sendiri dari GitHub Releases — mulai dari `wails3 init` hingga verifikasi rilis bertanda tangan dan pertukaran melalui mode pembantu."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

Dalam tutorial ini, Anda akan menambahkan pembaru dalam aplikasi ke aplikasi Wails v3 baru. Setelah selesai, aplikasi akan:

- Memeriksa GitHub Releases sesuai permintaan (dan secara opsional berdasarkan pengatur waktu).
- Mengunduh aset yang tepat untuk OS dan arsitektur yang sedang berjalan.
- Memverifikasi digest SHA-256 (dan secara opsional tanda tangan Ed25519) terhadap byte yang diunduh.
- Menampilkan catatan rilis di jendela pembaruan bawaan framework.
- Mengganti biner yang sedang berjalan dan meluncurkannya kembali — semuanya tanpa mendistribusikan executable pembantu terpisah.

Kita akan menggunakan **GitHub Releases** sebagai sumber pembaruan karena gratis dan tidak memerlukan infrastruktur. Pola yang sama dapat digunakan dengan [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) dan [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast) — setelah selesai, lihat [panduan Updater](/guides/updater/).

@note{type="tip" title="Prasyarat"}
- Go 1.25 atau yang lebih baru
- CLI `wails3` telah diinstal (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Repositori GitHub tempat Anda dapat mengirim rilis
- Pemahaman tentang [tutorial QR Code Service](/tutorials/01-creating-a-service/) akan membantu, tetapi tidak diwajibkan

@end

<br/>

@steps
### Mulai dengan aplikasi Wails baru
Buat kerangka proyek baru dengan templat vanilla:

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Sekarang Anda seharusnya memiliki direktori yang berisi `main.go`, `frontend/`, dan sebuah `Taskfile.yml`. Pastikan proyek dapat dibangun dan diluncurkan:

```bash
wails3 task dev
```

Jendela Wails kosong seharusnya terbuka. Tutup aplikasi, lalu lanjutkan.

### Tambahkan impor pembaru
Buka `main.go`, lalu tambahkan dua paket pembaru ke daftar impor:

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

Paket-paket ini menyertakan Updater itu sendiri dan penyedia GitHub Releases.

### Konfigurasikan Updater
`app.Updater` sudah terhubung ke setiap `*application.App` — Anda hanya perlu memanggil `Init`:

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

Letakkan ini setelah `application.New` dan sebelum `app.Run()`.

@note{type="note" title="Format string versi"}
Teruskan versi yang sama dengan tag rilis Anda, **tanpa** awalan `v`. Penyedia menghapus `v` dari nama tag di sisinya. `1.0.0` di sini ↔ `v1.0.0` di GitHub.

@end

### Tambahkan item menu yang memicu pembaruan
Dalam `main.go` yang sama, tambahkan entri menu "Periksa Pembaruan…":

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall` membuka jendela pembaruan framework, menjalankan `Check`, dan jika menemukan rilis, menjalankan `DownloadAndInstall` secara otomatis. Jika tidak ada yang baru, jendela tetap terbuka dalam status "Sudah Terbaru" — pengguna menutupnya dengan tombol **Tutup**.

@note{type="caution" title="Jalankan dalam goroutine"}
`CheckAndInstall` memblokir hingga verifikasi dan penginstalan selesai. Memanggilnya secara langsung dari klik menu akan memblokir thread UI. Bungkus dengan `go func()`.

@end

### Jalankan sekali tanpa rilis
```bash
wails3 task dev
```

Klik **App → Periksa Pembaruan…**. Jendela pembaruan seharusnya terbuka sebentar, mengakses API GitHub, tidak menemukan rilis yang lebih baru daripada `1.0.0`, lalu menetap pada status **Sudah Terbaru** dengan tanda ✓ berwarna hijau.

Jika Anda mendapatkan galat di sini, biasanya penyebabnya adalah salah satu hal berikut:

| Gejala | Perbaikan |
| --- | --- |
| `404 Not Found` | Kolom `Repository` salah — nilainya harus `owner/repo` |
| `403 rate-limited` | Tambahkan `Token: "ghp_…"` ke github.Config (gunakan PAT dengan cakupan `public_repo`) |
| Galat jaringan | Pastikan aplikasi yang sedang berjalan dapat mengakses `api.github.com` |

### Publikasikan rilis pengujian
Naikkan `currentVersion` dalam `main.go` menjadi `1.0.0` (atau biarkan tetap). Bangun untuk satu platform guna memperoleh biner yang dapat Anda lampirkan ke rilis:

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

Buat file `SHA256SUMS` di samping biner:

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

Anda seharusnya melihat satu atau beberapa baris seperti berikut:

```
abc123…  updater-tutorial-darwin-arm64.zip
```

Sekarang publikasikan ini sebagai **v2.0.0** di repositori GitHub Anda:

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="Penamaan aset"}
Pencocok aset bawaan memilih berdasarkan substring `GOOS` + `GOARCH` pada nama file. Selama nama aset Anda menyertakan `darwin` (atau `linux` / `windows`) dan `arm64` (atau `amd64` / `386`), pencocok akan menemukannya. Lihat [panduan Updater](/guides/updater/#github-releases--updaterprovidersgithub) untuk pencocok khusus.

@end

### Jalankan aplikasi dan verifikasi pembaruan
Dengan `currentVersion` yang masih bernilai `1.0.0`, jalankan kembali aplikasi:

```bash
wails3 task dev
```

Klik **App → Periksa Pembaruan…**. Kali ini Anda seharusnya melihat sesuatu seperti berikut:

![Jendela pembaru bawaan dalam status Pembaruan Siap, menampilkan label versi, catatan rilis yang dirender dari Markdown, dan tombol utama Mulai Ulang & Terapkan.](/assets/updater/default-window-ready.png)

- Ikon utama berubah dari ↓ biru ("Pembaruan Tersedia") menjadi ✓ hijau ("Pembaruan Siap").
- Subjudul menampilkan `v1.0.0 → v2.0.0 · <size>`.
- Panel catatan rilis merender Markdown Anda beserta teks tebal, rentang kode, dan tabel.
- Bilah progres terisi selama pengunduhan (prosesnya akan cepat — ukuran binernya kecil).

Updater menyiapkan biner baru dalam direktori sementara. Untuk menyelesaikan pembaruan:

- Klik **Mulai Ulang & Terapkan**.
- Aplikasi Anda berhenti, pembantu mengganti biner, lalu biner baru diluncurkan kembali.
- Aplikasi yang diluncurkan kembali melaporkan `currentVersion = "1.0.0"` (karena kita menetapkan nilainya secara hardcode), tetapi byte pada disk cocok dengan build v2.0.0.

Dalam aplikasi sungguhan, `currentVersion` akan ditetapkan saat build melalui `-ldflags` agar biner baru mengetahui bahwa versinya kini v2.0.0 dan pemeriksaan berikutnya tidak menemukan pembaruan.

### Hubungkan `currentVersion` ke build
Ganti konstanta dengan variabel waktu build:

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

Kemudian, dalam perintah build Anda:

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

Atau tambahkan `-ldflags` ke `Taskfile.yml` agar nilainya diambil dari `git describe --tags`.

### Tambahkan penandatanganan kriptografis (direkomendasikan untuk produksi)
Jalur SHA256SUMS memverifikasi *integritas* (byte cocok dengan yang disimpan GitHub), tetapi tidak memverifikasi *keaslian* (bahwa byte tersebut dihasilkan oleh pipeline rilis Anda, bukan oleh akun pengelola yang telah disusupi). Agar tahan terhadap manipulasi, tandatangani setiap rilis dengan kunci Ed25519:

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

Untuk setiap rilis, tanda tangani digest SHA-256 setiap aset dengan kunci privat Anda. Berikut pembantu Go sederhana:

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

Penyedia GitHub bawaan saat ini tidak mengambil berkas tanda tangan terpisah — Anda dapat [menulis penyedia khusus](/guides/updater/#writing-your-own-provider) yang melakukannya, atau beralih ke **keygen.sh**, yang menandatangani setiap artefak di sisi server dan menyediakan digest serta tanda tangannya melalui API.

Sematkan kunci publik dalam aplikasi Anda:

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

Jika `PublicKey` ditetapkan, setiap rilis yang menyertakan `Signature` harus berhasil diverifikasi dengan kunci ini. Sumber rilis tidak dapat menggantinya dengan kuncinya sendiri — itulah tujuan utama penyematan di luar saluran saat build.

### Sesuaikan jendela
Jendela bawaan mencakup kasus penggunaan umum. Tersedia tiga opsi jika Anda memerlukan kontrol lebih besar — pilih salah satu berdasarkan tingkat penyesuaian yang Anda inginkan:

@tabs
[Hanya CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

Lihat bagian [Tema melalui variabel CSS](/guides/updater/#theme-via-css-variables) untuk daftar lengkap variabel.

[HTML khusus]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

HTML Anda harus berlangganan peristiwa `updater:*` dan mengirim tindakan `updater:user:*` melalui saluran peristiwa Wails. Lihat [Ganti templat](/guides/updater/#replace-the-template) untuk shim JS.

[Gunakan jendela Anda sendiri]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Berguna jika Anda sudah memiliki infrastruktur jendela sendiri dan ingin pembaru mengendalikannya alih-alih membuka jendela lain. Templat HTML yang sepenuhnya khusus — dikendalikan oleh peristiwa pembaru yang sama seperti templat bawaan — terlihat seperti ini:

![Jendela pembaru milik sendiri dengan latar belakang gradasi merah muda-jingga dan tata letak kartu membulat khusus, yang menunjukkan bahwa UI bawaan dapat diganti sepenuhnya.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` wajib digunakan"}
Shim HTML khusus milik pembaru menjalankan Install / Skip / Remind / Restart melalui pintasan postMessage `wails:event:emit:`, dan demi keamanan, pintasan tersebut hanya diizinkan jika bidang ini diaktifkan. Jika Anda lupa mengaktifkannya, tombol-tombol tersebut tidak akan melakukan apa pun tanpa menampilkan kesalahan. Jangan aktifkan pada jendela yang memuat HTML yang tidak sepenuhnya Anda kendalikan — lihat bagian [Gunakan jendela Anda sendiri](/guides/updater/#bring-your-own-window) dalam panduan untuk memahami model ancamannya.

@end

@end

### Jalankan pemeriksaan otomatis di latar belakang
Untuk melakukan pemeriksaan menggunakan pewaktu sebagai pengganti (atau selain) klik menu:

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

Setiap interval menjalankan alur `CheckAndInstall` yang sama seperti klik manual. Tetapkan `Window: updater.WindowNone` jika Anda ingin pemeriksaan berkala tetap senyap hingga sesuatu benar-benar ditemukan — lalu berlanggananlah sendiri ke `EventUpdateAvailable` untuk menentukan UX yang akan ditampilkan.

@end

## Selesai

Sekarang Anda memiliki aplikasi Wails yang:

- Memeriksa pembaruan di GitHub Releases sesuai permintaan dan menggunakan pewaktu.
- Merender catatan rilis sebagai Markdown dalam jendela bawaan yang rapi.
- Memverifikasi unduhan terhadap digest SHA-256 yang Anda publikasikan.
- Secara opsional memverifikasi tanda tangan Ed25519 terhadap kunci publik yang Anda sematkan saat build.
- Mengganti biner yang sedang berjalan di lokasi yang sama dan meluncurkannya kembali secara otomatis.

## Langkah berikutnya

- [Panduan pembaru](/guides/updater/) berisi referensi API lengkap, setiap peristiwa, setiap opsi konfigurasi, dan mekanisme penggantian dalam mode pembantu.
- Lihat [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater) untuk contoh lengkap yang berfungsi dan dapat Anda klona.
- Repositori target pengujian [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo) menunjukkan tata letak aset rilis yang direkomendasikan.

## Hal yang perlu diperhatikan dalam produksi

- **Penandatanganan kode di macOS** — Gatekeeper mengharuskan biner pengganti ditandatangani dan dinotarisasi. Tanda tangani bundel `.app` Anda *sebelum* mengompresinya sebagai ZIP untuk rilis. Pembaru mempertahankan setiap byte apa adanya; pembaru tidak menandatangani ulang apa pun.
- **Antivirus di Windows** — berkas `.exe` tanpa tanda tangan yang diunduh dari internet dapat memicu peringatan SmartScreen. Tanda tangani biner Anda dengan sertifikat Authenticode, atau terima bahwa pengguna pada mesin dengan pembatasan ketat mungkin perlu memasukkan aplikasi Anda ke daftar yang diizinkan.
- **Rilis atomik** — publikasikan `SHA256SUMS` (beserta biner Anda) secara bersamaan, bukan dalam commit terpisah. Pembaru mengambil sidecar secara terpisah dari biner; jika keduanya tidak sinkron, pemeriksaan digest akan gagal dengan aman.
- **Versi yang dilewati** — tombol "Lewati Versi Ini" pada jendela bawaan mencatat versi yang dilewati secara lokal. Jika Anda merilis pembaruan keamanan penting, berikan nomor versi baru agar pembaruan tersebut tidak otomatis dilewati oleh pengguna yang mengabaikan rilis sebelumnya.
