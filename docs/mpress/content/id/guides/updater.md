---
title: "Pembaruan"
description: "Pembaruan mandiri dalam aplikasi untuk Wails v3 — penyedia yang dapat dipasang, verifikasi kriptografis, penggantian atomik, dan UI bawaan yang dapat Anda sesuaikan temanya atau ganti."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

Pembaruan mengirimkan pembaruan perangkat lunak dalam aplikasi tanpa mengharuskan Anda membuat sendiri alur pengunduhan, verifikasi, dan penggantian. Komponen ini berjalan di atas `app.Updater`, menerima satu atau beberapa `Provider` yang dapat dipasang (GitHub Releases, keygen.sh, Sparkle AppCast, protokol terbuka Wails Update Manifest, atau milik Anda sendiri), mengautentikasi unduhan menggunakan kunci publik yang dikonfigurasi, mengganti biner yang sedang berjalan dengan aman, dan menyampaikan setiap transisi melalui bus peristiwa Wails standar.

![Jendela pembaruan bawaan dalam status Pembaruan Siap — ikon yang menyesuaikan status, label versi (v1.0.0 → v2.0.1 · 8.8 MB), catatan rilis yang dirender dari Markdown termasuk tabel GFM, dan satu tindakan utama.](/assets/updater/default-window-ready.png)

## Mulai cepat

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

Tindakan tersebut membuka jendela pembaruan framework, memeriksa GitHub, mengunduh artefak platform, memverifikasinya, mengganti biner, lalu menunggu pengguna memulai ulang aplikasi.

## Siklus hidup

`app.Updater` adalah mesin status dengan status berikut (`updater.State`):

| Status | Kondisi |
| --- | --- |
| `unconfigured` | Sebelum `Init` dipanggil |
| `idle` | Setelah `Init`, sebelum pemeriksaan apa pun |
| `checking` | Sebuah `Check` sedang berlangsung |
| `up-to-date` | Respons penyedia terbaru menyatakan bahwa pemanggil sudah menggunakan versi terkini |
| `available` | Rilis baru ditemukan; pengunduhan belum dimulai |
| `downloading` | Byte sedang dialirkan dari penyedia |
| `verifying` | Pengunduhan selesai; tanda tangan/digest sedang diperiksa |
| `installing` | Byte yang telah diverifikasi sedang dibongkar dan diganti namanya ke direktori staging |
| `ready` | Pembaruan telah ditempatkan di staging; panggil `Restart` untuk menerapkannya |
| `error` | Salah satu langkah sebelumnya gagal |

Anda dapat membaca status saat ini kapan saja dengan `app.Updater.State()`. Setiap transisi juga memancarkan peristiwa Wails (lihat [Peristiwa](#event)).

`Restart` menunggu helper mencapai `application.New` sebelum meminta aplikasi yang sedang berjalan untuk berhenti. Batas waktu startup bawaan adalah 30 detik. Jika aplikasi Anda melakukan inisialisasi yang lama sebelum `application.New`, atur `Config.HelperReadyTimeout` ke durasi yang lebih panjang, misalnya `time.Minute`. Nilai nol memilih nilai bawaan; durasi negatif ditolak. Jika startup melewati batas waktu, `Restart` mengembalikan `updater.ErrHelperNotReady` dan membiarkan aplikasi yang sedang berjalan tetap terbuka.

Pembaruan ditempatkan di staging dalam direktori `wails-update-*` di samping biner yang terinstal (di macOS, di samping bundel `.app`), sehingga direktori tersebut harus dapat ditulisi.

Di dalam AppImage, pembaruan mengganti dan menjalankan ulang berkas `.AppImage` yang ditentukan oleh `$APPIMAGE`, bukan mount hanya-baca.

Jendela bawaan secara otomatis mencerminkan status saat ini — misalnya, ketika `Check` menyatakan tidak ada peningkatan, pengguna akan melihat tampilan ini dan menutupnya dengan **Tutup**:

![Jendela pembaruan bawaan dalam status Sudah Terkini — tanda centang hijau, judul 'Versi Anda Sudah Terkini', dan satu tombol Tutup.](/assets/updater/default-window-up-to-date.png)

## Penyedia

`Provider` adalah apa pun yang memenuhi antarmuka ini:

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

Empat implementasi disertakan dalam repositori.

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

Pencocok artefak bawaan memilih berdasarkan substring `GOOS` + `GOARCH` pada nama berkas dan mengenali alias umum (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Untuk skema penamaan khusus:

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset` adalah nama artefak rilis pendamping yang isinya berupa baris-baris `<sha256>  <filename>` (format yang dihasilkan oleh `sha256sum` dan `shasum -a 256`). Penyedia mengambilnya selama `Check`, menemukan baris yang cocok dengan artefak terpilih, lalu mengisi `Release.Verification.Digest` agar framework memverifikasi unduhan.

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

Penyedia secara otomatis memetakan checksum SHA-512 dan tanda tangan Ed25519ph per artefak milik keygen.sh ke blok `Release.Verification` milik framework — tidak memerlukan konfigurasi tambahan.

**Format token:** token keygen.sh memiliki prefiks peran (`admi-` / `prod-` / `envi-` / `user-`). UUID mentah yang Anda lihat di dasbor adalah *pengidentifikasi* token, bukan nilai rahasianya — rahasia tersebut hanya terlihat saat token dibuat. Lihat [dokumentasi autentikasi](https://keygen.sh/docs/api/authentication/) keygen.sh untuk detailnya.

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

Dapat langsung digunakan dengan infrastruktur Sparkle / WinSparkle yang sudah ada tanpa perubahan. Membaca `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>`, dan `sparkle:channel` dari feed.

Tanda tangan DSA milik Sparkle 1 (`sparkle:dsaSignature`) tidak didukung — proyek yang menggunakan skema penandatanganan tersebut sebaiknya beralih ke EdDSA (Sparkle 2).

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Menggunakan [protokol terbuka Wails Update Manifest](/reference/update-manifest/): satu dokumen JSON yang menjelaskan rilis terbaru beserta artefak per platformnya, dengan checksum dan tanda tangan sebaris. Dokumen yang sama dapat digunakan dari host berkas statis (S3, GitHub Pages, CDN apa pun — publikasikan satu manifes per kanal yang mencantumkan semua platform) atau dari server pembaruan dinamis (penyedia mengirim `platform`, `arch`, `version`, dan `channel` pada setiap pemeriksaan, sehingga server dapat mengembalikan tepat satu artefak atau membatasi akses berdasarkan lisensi).

Placeholder URL memungkinkan tata letak statis dikonfigurasi dalam satu baris:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Header yang dikonfigurasi dikirim pada setiap permintaan manifes; unduhan artefak menggunakannya kembali hanya pada host manifes itu sendiri tanpa downgrade dari `https` ke `http`, dan header `Authorization` dihapus pada setiap pengalihan lintas origin atau pengalihan yang melakukan downgrade.

CLI menangani sisi publikasi: `wails3 updater manifest` menghitung digest, menandatangani, dan mendeskripsikan berkas rilis Anda dalam satu perintah, sedangkan `wails3 updater verify` memeriksa ulang hasilnya sebelum Anda mengunggahnya. Lihat [Publikasi dengan CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Rantai fallback

`Config.Providers` memiliki urutan. Pembaruan menelusurinya secara berurutan: penyedia pertama yang mengembalikan rilis akan digunakan, sedangkan penyedia pertama yang melaporkan "sudah terkini" segera menghentikan rantai (fallback ditujukan untuk kondisi "penyedia utama tidak dapat dijangkau", bukan "para penyedia tidak sepakat"). Kesalahan akan melanjutkan proses ke penyedia berikutnya.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### Menulis penyedia Anda sendiri

Tiga metode, sekitar 150 baris untuk implementasi biasa. Updater menangani verifikasi, staging atomik, penggantian, dan jendela — kode penyedia menentukan rilis berikutnya dan mengalirkan byte:

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

Gunakan penyedia dalam repositori sebagai referensi — masing-masing hanya terdiri dari satu berkas Go.

## Verifikasi kriptografis

Rilis diautentikasi oleh pemverifikasi framework menggunakan `Config.PublicKey` sebagai akar kepercayaan:

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

Algoritma yang didukung (`Release.Verification.SignatureAlgo`):

| Algoritma | Yang ditandatangani | Catatan |
| --- | --- | --- |
| `ed25519` | Digest SHA-256 artefak | Digunakan oleh Sparkle EdDSA |
| `ed25519ph` | Artefak lengkap melalui pra-hash Ed25519ph (secara internal menggunakan SHA-512) | Digunakan oleh keygen.sh |
| `ecdsa-p256` | Digest SHA-256 artefak | Tanda tangan mentah `r∥s` maupun DER diterima |

Selain itu, tersedia mode khusus digest (`DigestAlgo`: `sha256` / `sha512`) saat rilis menyertakan hash tetapi tidak menyertakan tanda tangan.

`Config.PublicKey` adalah SATU-SATUNYA jangkar kepercayaan untuk verifikasi tanda tangan—sumber rilis tidak dapat menggantinya dengan kunci miliknya sendiri. Rilis yang membawa `Signature` tanpa `Config.PublicKey` yang dikonfigurasi akan ditolak secara aman. Pemverifikasi menghitung digest dalam satu kali proses streaming selama pengunduhan, sehingga verifikasi tidak memerlukan pembacaan disk tambahan, bahkan untuk pembaruan berukuran beberapa GB.

@note{type="caution" title="Khusus digest ≠ verifikasi kriptografis"}
Rilis yang hanya memiliki `Digest` diautentikasi berdasarkan TLS registri serta jaminan integritas apa pun yang disediakan oleh registri itu sendiri—bukan berdasarkan akar kriptografis yang Anda kendalikan. Gunakan mode khusus digest untuk mendeteksi kerusakan bit; gunakan tanda tangan agar tahan terhadap manipulasi jika pipeline rilis disusupi.

@end

### Membuat kunci penandatanganan

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

Kunci privat menggunakan PEM PKCS#8, sedangkan kunci publik menggunakan PEM PKIX; `Config.PublicKey` menerima berkas `.pub` apa adanya (juga menerima kunci mentah 32 bita atau bentuk base64-nya, yang dicetak oleh `genkey` untuk disematkan). Tandatangani rilis dengan `wails3 updater manifest -key updater.key ...` atau `wails3 updater sign`; lihat [Menerbitkan dengan CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

Atau dalam Go:

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## Format artefak

Penyedia mengalirkan bita untuk berkas apa pun yang Anda terbitkan; framework kemudian mengekstraknya sebelum pertukaran:

- **Biner tunggal** (misalnya `myapp-linux-amd64`)—digunakan apa adanya. Umum digunakan di Linux.
- **`.zip`**—diekstrak di lokasi tujuan. Arsip harus berisi tepat satu entri tingkat teratas (biasanya bundel `.app` macOS atau satu biner). Ini adalah kemasan yang disarankan untuk macOS.
- **`.tar.gz`** / **`.tgz`**—diekstrak di lokasi tujuan dengan aturan satu entri tingkat teratas yang sama. Berguna untuk distribusi Linux yang menyertakan pohon runtime bersama biner.

Arsip yang berisi lebih dari satu entri tingkat teratas akan ditolak: framework menukar satu target pada disk, sehingga perintah "tukar arsip ini ke lokasi tujuan" menjadi ambigu jika arsip berisi beberapa hal. `.dmg` dan `.pkg` (macOS), serta `.msi` (Windows), tidak didukung dalam v1—sebagai gantinya, distribusikan `.zip` dari bundel tersebut. Proses ekstraksi menerapkan perlindungan zip-slip, menolak symlink yang keluar dari akar arsip, serta membatasi ukuran total setelah dekompresi (2 GiB) dan jumlah entri (50 000).

## Jendela bawaan

`app.Updater.CheckAndInstall(ctx)` membuka jendela berukuran 520×540 yang dikelola framework dengan:

- Ikon utama berbasis status (↓ biru untuk tersedia/sedang diunduh, ✓ hijau untuk siap/sudah terbaru, ! merah untuk kesalahan)
- Label versi: `v1.0.0 → v2.0.1 · 8.8 MB`
- Panel catatan rilis yang dapat digulir dengan **Markdown yang dirender** (paragraf, tebal/miring, daftar, tabel GFM, kode sebaris, blok kode berpagar, h1–h3, tautan)
- Satu tindakan utama untuk setiap status (Instal / Mulai Ulang & Terapkan / Coba Lagi)
- Tindakan sekunder bergaya transparan (Lewati Versi Ini / Ingatkan Saya Nanti)
- Mode gelap/terang melalui `prefers-color-scheme`
- Efek kilau progres tak tentu saat ukuran total tidak diketahui

Jendela ini mendengarkan peristiwa `updater:*` pada bus peristiwa Wails dan mengirim tindakan `updater:user:*` kembali ke Go.

### Tema melalui variabel CSS

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

Stylesheet bawaan menyediakan variabel berikut—Anda dapat menimpa variabel mana pun:

| Variabel | Bawaan (terang) | Bawaan (gelap) |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | tumpukan sistem | — |

### Mengganti templat

Sediakan HTML Anda sendiri; HTML tersebut hanya perlu mendengarkan peristiwa `wails:updater:*` dan mengirim tindakan `wails:updater:user:*`:

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

Jendela InitialHTML dimuat tanpa origin server aset sehingga tidak dapat mengambil `/wails/runtime.js` secara dinamis. Ada dua cara untuk berkomunikasi dengan host dari jendela tersebut:

1. **Cukup tulis HTML.** Framework secara otomatis menyuntikkan shim `window.wails.Events` minimal ke setiap jendela yang dibuka dengan `WebviewWindowOptions.AllowSimpleEventEmit = true` dan `HTML` yang telah ditetapkan—persis seperti yang dilakukan jalur bawaan dan BYO milik updater. Tidak diperlukan langkah build. Contoh di bawah menggunakan jalur ini.
2. **Bundelkan `@wailsio/runtime` dengan bundler pilihan Anda** (Vite, esbuild, Rollup), lalu impor ke HTML kustom Anda saat build. `Events.On` langsung berfungsi karena sepenuhnya berjalan di sisi klien; `Events.Emit` melewati transpor fetch milik runtime, yang tidak dapat digunakan karena origin null — jadi pasang transpor postMessage kecil melalui hook [`setTransport`](/guides/custom-transport/) milik runtime yang merutekannya melalui `window._wails.invoke("wails:event:emit:<name>")`. Injeksi framework tidak melakukan apa pun jika `window.wails.Events` sudah berada dalam scope, sehingga kedua pendekatan tersebut tidak saling mengganggu.

Apa pun caranya, JS yang Anda tulis dalam HTML kustom akan memanggil API `Events.On` / `Events.Emit` yang sama:

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

Shim tersebut menyediakan subset runtime modern yang diperlukan oleh event dengan nama polos: `Events.On(name, cb)` mengembalikan fungsi untuk berhenti berlangganan, sedangkan `Events.Emit(nameOrEventObject)` merutekan event ke host melalui jalur postMessage `wails:event:emit:` yang dilindungi. Shim dipasang satu kali saat halaman dimuat, sebelum skrip inline Anda dijalankan.

Jika Anda *ingin* mengganti shim tersebut (atau memuat runtime lengkap dengan cara lain), tetapkan `window.wails.Events` sebelum tag `<script>` pertama pada halaman dijalankan agar injeksi melewati dirinya sendiri.

### Bingkai jendela

Timpa opsi jendela (ukuran, tanpa bingkai, selalu di atas) tanpa mengubah HTML:

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### Gunakan jendela Anda sendiri

Jalankan alur pembaruan pada `*application.WebviewWindow` yang Anda buat sendiri. Updater memanggil `Show()` / `Close()` / `EmitEvent()` pada jendela Anda — HTML Anda menentukan apa yang akan dirender:

![Jendela updater “Gunakan jendela Anda sendiri” dengan latar belakang gradasi merah muda-oranye, satu kartu putih bersudut membulat, tipografi kustom, dan event updater yang sama untuk mengendalikan status yang terlihat. Contoh ini menunjukkan bahwa UI bawaan dapat diganti sepenuhnya.](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

HTML Anda menggunakan `window.wails.Events.On` / `Events.Emit` seperti templat bawaan — shim yang diinjeksi otomatis oleh framework disertakan dalam setiap jendela dengan `AllowSimpleEventEmit: true`, baik jendela tersebut dimiliki framework maupun Anda. Lihat [Ganti templat](#mengganti-templat) untuk API-nya.

@note{type="caution" title="`AllowSimpleEventEmit` wajib untuk jendela updater BYO"}
Demi keamanan, framework membatasi pintasan postMessage `wails:event:emit:` berdasarkan kolom ini: jendela yang tidak mengaktifkannya tidak dapat membuat event kustom pada sisi host. Shim HTML kustom milik updater mengirim event `updater:user:*` melalui pintasan tersebut, sehingga jendela BYO yang lupa mengaktifkan kolom ini akan mengabaikan setiap klik tombol tanpa pemberitahuan — pengguna mengeklik Instal dan tidak terjadi apa-apa.

Pastikan `AllowSimpleEventEmit` tetap **nonaktif** untuk setiap jendela yang memuat HTML yang tidak sepenuhnya Anda kendalikan (URL jarak jauh, konten yang disediakan pengguna). Saat aktif, JavaScript apa pun di halaman — termasuk sink XSS — dapat memicu handler `app.Event.On(name, …)` apa pun. Pintasan ini hanya menyampaikan nama polos (tanpa payload) dan tidak dapat menjangkau jalur binding/Call, tetapi tetap dapat memicu handler event kustom berhak istimewa dalam kode Go Anda jika handler tersebut bertindak hanya berdasarkan nama event.

Jendela updater *bawaan* framework mengaktifkan ini secara internal — hanya pemanggil BYO yang perlu mengingatnya.

@end

### Tanpa antarmuka

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Tidak ada jendela yang dibuka. Berlangganan event `updater:*` dari UI Anda sendiri (atau jendela utama yang sudah ada), lalu panggil `app.Updater.CheckAndInstall(ctx)` dari handler tombol. Berguna untuk pemeriksaan latar belakang berkala yang hanya perlu ditampilkan saat sesuatu ditemukan, atau untuk aplikasi yang mengintegrasikan alur pembaruan ke panel pengaturan kustom.

## Event

Go dan JavaScript sama-sama berlangganan melalui bus event Wails standar. **Jangan ketik string wire secara manual** — gunakan konstanta yang diekspor dari paket updater (Go) atau paket runtime (JS). Kedua lapisan menggunakan kumpulan nama yang sama dan dijaga agar tetap sinkron oleh pengujian regresi.

### Dari Go

Konstanta berada di `github.com/wailsapp/wails/v3/pkg/updater`. Berlangganan melalui `app.Event.On(name, fn)`; callback menerima `*application.CustomEvent` dengan kolom `Data` berupa payload bertipe yang tercantum dalam [referensi event](#referensi-event) — gunakan type assertion, bukan dekode JSON:

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

Semua konstanta Go yang tersedia:

| Konstanta | String wire |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### Dari JavaScript

Konstanta berada di bawah `Updater.Events` dalam `@wailsio/runtime`. Namanya sama seperti di Go dan diatur berdasarkan sub-namespace (`User.*`, `Window.*`) agar mudah ditemukan melalui pelengkapan otomatis:

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

Event tindakan pengguna yang dikirim *kembali* oleh HTML kustom Anda ke host berada di bawah `Updater.Events.User`:

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### Referensi event

Sisi pelanggan (host → halaman):

| Konstanta (Go) | Konstanta (JS) | Payload | Kapan |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | tidak ada | Sebelum setiap perjalanan pulang-pergi `Check` |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` menemukan rilis yang lebih baru |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | tidak ada | `Check` mengonfirmasi bahwa aplikasi sudah mutakhir |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Byte mulai dialirkan |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ~10 Hz selama pengunduhan |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Semua byte telah ditulis, sebelum verifikasi |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | Pemeriksaan tanda tangan/digest dimulai |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | Ekstraksi + staging dimulai |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Menunggu dimulai ulang |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Tahap mana pun gagal |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Sekali per sesi sebelum pemutaran ulang snapshot |

Sisi halaman (halaman → host) — kode Anda berlangganan jika Anda membuat templat khusus:

| Konstanta (Go) | Konstanta (JS) | Kapan |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | Jendela selesai dimuat; host memutar ulang status saat ini |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Tindakan utama dalam status `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Tindakan utama dalam status `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | "Lewati Versi Ini" |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | "Ingatkan Saya Nanti" |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Tombol tutup |

## Referensi API

### `updater.Config`

| Bidang | Tipe | Catatan |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Wajib.** String yang sama dengan yang Anda gunakan untuk menandai rilis (tanpa prefiks `v`) |
| `Providers` | `[]updater.Provider` | **Wajib.** Rantai fallback berurutan |
| `PublicKey` | `[]byte` | PEM atau byte mentah. Opsional, tetapi rilis bertanda tangan akan gagal secara aman tanpa ini |
| `CheckInterval` | `time.Duration` | Nilai bukan nol memulai loop polling latar belakang yang memanggil `CheckAndInstall` |
| `Platform` | `string` | Timpa `runtime.GOOS` untuk pemilihan aset |
| `Arch` | `string` | Timpa `runtime.GOARCH` untuk pemilihan aset |
| `Channel` | `string` | Saat ini hanya bersifat informatif; pemfilteran kanal khusus penyedia |
| `Window` | `updater.WindowOption` | `nil` (default bawaan), `&BuiltinWindow{…}`, `BYOWindow(handle)`, atau `WindowNone` |

### Metode pada `*updater.Updater`

| Tanda tangan | Tujuan |
| --- | --- |
| `Init(cfg Config) error` | Mengonfigurasi. Mengembalikan `ErrAlreadyConfigured` pada pemanggilan kedua |
| `State() State` | Fase siklus hidup saat ini |
| `CurrentVersion() string` | Versi yang diteruskan ke `Init` |
| `Check(ctx) (*Release, error)` | Menelusuri rantai penyedia. `(rel, nil)` = ditemukan, `(nil, nil)` = sudah mutakhir, `(nil, err)` = semuanya gagal |
| `DownloadAndInstall(ctx) error` | Mengalirkan, memverifikasi, mengekstrak (jika berupa arsip), lalu melakukan staging. Memerlukan `Check` sebelumnya |
| `CheckAndInstall(ctx) error` | Kemudahan: membuka jendela, menjalankan `Check`, lalu `DownloadAndInstall` jika ditemukan |
| `Restart(ctx) error` | Menjalankan helper, memanggil `Host.Quit`, lalu keluar; helper menukar + meluncurkan ulang |
| `DownloadedPath() string` | Lokasi pembaruan yang telah di-staging pada disk, atau `""` jika tidak ada |
| `SkipVersion(v string)` | Mencatat `v` sebagai dilewati; `Check` berikutnya menganggapnya sudah mutakhir |
| `SkippedVersion() string` | Membaca versi yang saat ini dilewati |
| `StopPeriodicCheck()` | Membatalkan timer yang dimulai oleh `Config.CheckInterval` dan menunggu hingga loop selesai |

### Kesalahan

| Sentinel | Dikembalikan oleh |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` setelah keberhasilan pertama |
| `ErrNotConfigured` | Operasi apa pun sebelum `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` tanpa didahului `Check` |
| `ErrDownloadInProgress` | `DownloadAndInstall` dipanggil saat proses lain sedang berjalan |
| `ErrNotReady` | `Restart` tanpa pembaruan yang telah disiapkan |

## Cara kerja penggantian

`Restart` menjalankan ulang biner saat ini dengan variabel lingkungan sentinel yang telah ditetapkan. `application.New` mendeteksinya saat aplikasi dimulai dan mengalihkan eksekusi ke mode pembantu:

1. Pembantu menunggu hingga 30 detik sampai PID induk berhenti (`platformIsAlive` melakukan polling melalui `syscall.OpenProcess` + `GetExitCodeProcess` di Windows, serta `os.FindProcess` + `proc.Signal(syscall.Signal(0))` di Unix).
2. Pembantu mencadangkan target (menyalin file, atau menyalin secara rekursif direktori bundel `.app` macOS).
3. Pembantu mengganti target dengan artefak yang telah disiapkan dan mencoba ulang hingga 20 kali dengan jeda 500 milidetik di antara setiap percobaan:
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Deskriptor file yang terbuka dan mengacu pada inode lama tetap valid.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows mengizinkan penggantian nama file yang image-nya masih dipetakan, tetapi tidak mengizinkan penghapusannya; pada pembaruan berikutnya, pembantu membersihkan file saudara `.old.*` yang tersisa setelah pemetaan kernel pemiliknya dilepaskan.

4. Pembantu memulihkan mode executable asli pada biner baru (file unduhan dibuat dengan umask bawaan, yang menghapus `+x` di Unix; di Windows, operasi ini tidak melakukan apa pun).
5. Pembantu menghapus variabel lingkungan mode pembantu dan menjalankan kembali biner yang kini telah diganti.
6. Pembantu berhenti.

Jika peluncuran gagal, pembantu memulihkan cadangan. Jika proses induk belum berhenti dalam 30 detik, pembantu membatalkan operasi sebelum menyentuh target (sehingga pengguna tetap memiliki aplikasi yang berfungsi meskipun dialog penutupan memblokir `Quit`).

Untuk bundel `.app` macOS yang didistribusikan sebagai `.zip` (kemasan yang disarankan), arsip dibongkar di antara tahap verifikasi dan siap agar pembantu memiliki direktori nyata untuk ditempatkan melalui penggantian.

## Pemeriksaan berkala

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Saat `CheckInterval > 0`, goroutine latar belakang memanggil `CheckAndInstall` sesuai interval yang dikonfigurasi. Tick yang tiba ketika alur lain masih berlangsung (memeriksa / mengunduh / memverifikasi / menginstal) akan diabaikan—mesin status konkuren tidak didukung.

Untuk polling latar belakang senyap yang hanya ditampilkan ketika sesuatu ditemukan, tetapkan `Window: updater.WindowNone` dan tanggapi `EventUpdateAvailable` dari UI Anda sendiri.

## Lewati & ingatkan

Tombol "Lewati Versi Ini" pada jendela bawaan mencatat versi yang tersedia melalui `SkipVersion(rel.Version)`. Pemanggilan `Check` berikutnya menemukan versi yang sama dan menganggap aplikasi sudah terbaru (hingga pengguna memperbarui `CurrentVersion`, yang terjadi secara otomatis setelah `Restart` berhasil). "Ingatkan Saya Nanti" hanya menutup jendela tanpa mencatat apa pun.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## Daftar periksa distribusi

Sebelum menerbitkan rilis yang akan diinstal oleh pemutakhir:

1. **Pilih format arsip yang tepat.** macOS: `.zip` dari bundel `.app`. Linux: satu biner atau `.tar.gz`. Windows: satu `.exe` atau `.zip`. `.dmg` / `.msi` / `.pkg` tidak didukung.
2. **Tandatangani artefak** dengan kunci privat yang cocok dengan `Config.PublicKey`. Untuk feed yang diterbitkan penyedia (keygen.sh, AppCast), ikuti alur kerja penandatanganan masing-masing penyedia. Untuk GitHub Releases dengan `ChecksumAsset`, buat file `SHA256SUMS` menggunakan `sha256sum` / `shasum -a 256`.
3. **Cocokkan string versi.** `Config.CurrentVersion` dan tag versi rilis harus sama persis (misalnya `1.0.0` ↔ tag `v1.0.0`; awalan `v` dihapus di sisi penyedia).
4. **Uji penggantian di platform target** setidaknya sekali sebelum merilis—penandatanganan kode, notarisasi, dan penanganan Gatekeeper bersifat khusus platform dan tidak ditangani oleh pemutakhir itu sendiri.

## Pemecahan masalah

**"signature requires a public key but none configured"** — rilis memiliki bidang `Signature`, tetapi `Config.PublicKey` kosong. Tetapkan kunci publik atau ubah pipeline rilis agar tidak menyertakan tanda tangan.

**"digest mismatch"** — byte yang diunduh tidak cocok dengan yang dijanjikan penyedia. Biasanya penyebabnya adalah unduhan parsial (gangguan jaringan) atau artefak yang rusak. Menjalankannya kembali sering kali mengatasi masalah ini.

**Jendela terbuka tetapi langsung menghilang, tanpa markdown dan tanpa progres** — HTML kustom Anda tidak memanggil `wails:runtime:ready`. Lihat shim [Ganti templat](#mengganti-templat).

**Pembaruan Windows tidak pernah selesai; log pembantu menampilkan "remove old (attempt N): Access is denied"** — ini hanya terjadi pada versi PR ini sebelum `de764fb`; implementasi saat ini menggunakan penggantian nama ke samping sehingga tidak mengalami masalah tersebut. Lakukan pemutakhiran.

**Gatekeeper macOS memblokir biner yang telah diganti** — penandatanganan kode harus dipertahankan dari awal hingga akhir. Tandatangani `.app` asli *dan* tandatangani ulang biner yang dijalankan kembali jika pipeline build Anda mengubah entitlement pada saat pembaruan.

## Lihat juga

- Contoh yang dapat dijalankan: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- Repositori demo pengujian: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Tutorial: [Menambahkan pembaruan mandiri ke aplikasi Wails](/tutorials/04-self-update-a-wails-app/)
