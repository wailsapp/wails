---
title: "Pertanyaan yang Sering Diajukan"
description: "Jawaban atas pertanyaan paling umum tentang pembuatan aplikasi dengan Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Umum

### Apa itu Wails?

Wails adalah framework untuk membuat aplikasi desktop dengan Go dan teknologi web. Logika aplikasi Anda ditulis dalam Go, antarmukanya dibuat dengan HTML, CSS, dan JavaScript (atau framework frontend apa pun), lalu Wails merendernya dalam webview bawaan sistem operasi. Hasilnya adalah aplikasi berukuran kecil, cepat, dan terasa native: tanpa browser yang dibundel, dengan penggunaan memori rendah, serta satu berkas biner yang biasanya berukuran sekitar 10MB.

### Platform apa saja yang didukung Wails?

| Platform | Persyaratan |
| --- | --- |
| Windows | AMD64 dan ARM64. Menggunakan runtime [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15+ pada Intel (aplikasi dapat menargetkan 10.13+), 11.0+ pada Apple Silicon. Biner universal didukung. |
| Linux | AMD64 dan ARM64. Stack default adalah GTK4 dengan WebKitGTK 6.0 (Ubuntu 24.04+, Debian 13+, Fedora 40+, dan distribusi serupa). Distribusi yang hanya menyediakan WebKit2GTK 4.1, seperti Ubuntu 22.04, Debian 12, dan RHEL 9, didukung melalui build lama `-tags gtk3` (tersedia hingga v3.1). Distribusi yang hanya menyediakan WebKit2GTK 4.0 tidak didukung. Lihat [panduan build Linux](/guides/build/linux/). |
| iOS dan Android | Eksperimental. Lihat [panduan perangkat seluler](/guides/mobile/). |

Anda juga dapat menyajikan aplikasi sebagai aplikasi web biasa menggunakan [build server](/guides/server-build/).

Jalankan `wails3 doctor` kapan saja untuk memeriksa sistem Anda dan mendapatkan petunjuk penginstalan khusus platform.

### Apa yang saya perlukan untuk memulai?

- Go 1.26 atau yang lebih baru
- Node.js dan npm (untuk build frontend)
- Toolchain platform: WebView2 pada Windows (sudah terinstal di 10/11), Xcode Command Line Tools pada macOS, serta `gcc` dan paket pengembangan GTK/WebKit pada Linux

`wails3 doctor` memeriksa semua ini untuk Anda dan memberi tahu dengan tepat apa yang belum tersedia. Lihat [Penginstalan](/quick-start/installation/) untuk panduan lengkapnya.

### Apakah Wails v3 siap digunakan dalam produksi?

Wails v3.0.0 adalah rilis stabil untuk aplikasi desktop pada Windows, macOS, dan Linux. Dukungan Android dan iOS masih bersifat eksperimental dan tidak termasuk dalam kontrak kompatibilitas desktop stabil. Lihat [halaman status proyek](/status/).

## Pengembangan

### Apakah saya perlu menguasai Go?

Pengetahuan dasar tentang Go akan membantu, tetapi Anda tidak perlu menjadi ahli. Logika aplikasi berada dalam metode Go biasa, dan [tutorial](/tutorials/overview/) akan memandu Anda mempelajari semua hal lainnya. Banyak pengembang mempelajari Go sambil membuat aplikasi Wails pertama mereka.

### Dapatkah saya menggunakan framework frontend favorit saya?

Ya. Jika dapat dibangun menjadi HTML, CSS, dan JavaScript, framework tersebut dapat digunakan dengan Wails. Tersedia templat untuk React, Vue, Svelte, dan JavaScript biasa (masing-masing dengan varian TypeScript), dan framework lainnya dapat diintegrasikan dalam hitungan menit. Lihat [Framework Frontend](/guides/dev/frontend-frameworks/).

### Bagaimana cara memanggil fungsi Go dari JavaScript?

Daftarkan sebuah layanan, lalu Wails akan menghasilkan binding bertipe untuk layanan tersebut:

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

Binding dibuat ulang secara otomatis selama `wails3 dev`, atau sesuai kebutuhan dengan `wails3 generate bindings`. Lihat [Layanan](/features/bindings/services/).

### Dapatkah saya menggunakan TypeScript?

Ya. Generator binding menghasilkan definisi TypeScript untuk layanan dan tipe Anda, sehingga pemanggilan ke Go memiliki tipe yang lengkap.

### Bagaimana cara mengirim event antara Go dan JavaScript?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

Nama event harus sama persis. Lihat [Referensi Event](/guides/events-reference/).

### Bagaimana cara men-debug aplikasi saya?

Jalankan `wails3 dev`, lalu klik kanan di dalam jendela untuk membuka alat pengembang browser, persis seperti yang Anda lakukan di web. Server pengembangan juga mendukung hot reload untuk frontend Anda. Lihat [Debugging](/guides/dev/debugging/).

## Build & Distribusi

### Bagaimana cara membuat build untuk produksi?

```bash
wails3 build
```

Berkas biner Anda akan ditempatkan di `bin/`. Build produksi sudah menerapkan default yang sesuai (tag build, `-trimpath`, simbol yang dihapus), sehingga tidak diperlukan flag tambahan untuk menghasilkan berkas biner yang ringkas.

### Dapatkah saya melakukan kompilasi silang?

Bisa, dengan batasan tertentu. Kompilasi silang Go murni tidak dapat digunakan karena setiap platform memakai pustaka webview native, tetapi tersedia dukungan yang baik untuk kasus-kasus umum:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

Pembuatan build Linux dari sistem operasi lain menggunakan toolchain berbasis Docker. Lihat [Build Lintas Platform](/guides/build/cross-platform/) untuk matriks lengkapnya.

### Bagaimana cara membuat penginstal atau paket?

```bash
wails3 package
```

Proses ini menghasilkan format native platform, dan [panduan Penginstal](/guides/installers/) membahas NSIS pada Windows, bundle `.app` dan DMG pada macOS, serta paket Linux.

### Bagaimana cara menandatangani kode aplikasi saya?

Penandatanganan pada Windows dan macOS (termasuk notarisasi) dibahas langkah demi langkah dalam [panduan Penandatanganan](/guides/build/signing/).

## Fitur

### Dapatkah saya membuat beberapa jendela?

Ya, dukungan untuk beberapa jendela tersedia secara native di v3:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Lihat [Beberapa Jendela](/features/windows/multiple/).

### Apakah Wails mendukung baki sistem?

Ya, termasuk menu dan handler klik:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Lihat [Baki Sistem](/features/menus/systray/).

### Dapatkah saya menggunakan dialog native?

Ya. Dialog berkas, dialog pesan, dan dialog pertanyaan semuanya menggunakan implementasi native:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Lihat [Dialog](/features/dialogs/overview/).

### Apakah Wails mendukung pembaruan otomatis?

Ya. Wails v3 menyertakan pembaru mandiri bawaan (`app.Updater`) dengan penyedia yang dapat dipasang untuk GitHub Releases, keygen.sh, dan Sparkle AppCast, verifikasi tanda tangan kriptografis, serta antarmuka default yang dapat Anda sesuaikan temanya atau ganti. Lihat panduan [Pembaruan Dalam Aplikasi](/guides/updater/) dan tutorial [Aplikasi Wails yang Memperbarui Diri](/tutorials/04-self-update-a-wails-app/).

## Pemecahan Masalah

### Ada sesuatu yang tidak berfungsi. Dari mana saya harus memulai?

```bash
wails3 doctor
```

Fitur ini memverifikasi toolchain Anda, mencantumkan dependensi yang belum tersedia beserta perintah penginstalannya, dan menampilkan informasi versi yang harus Anda sertakan dalam setiap laporan bug.

### Build saya gagal

Cobalah solusi umum berikut secara berurutan:

1. `go mod tidy`
2. `cd frontend && npm install` (`node_modules` yang tidak tersedia adalah penyebab paling umum)
3. Perbarui CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. Pada Linux, periksa paket GTK/WebKit yang belum tersedia dengan `wails3 doctor`

### Binding saya hilang atau kedaluwarsa

```bash
wails3 generate bindings
```

Binding dibuat ulang secara otomatis dalam mode pengembangan; jika Anda menambahkan layanan baru atau mengubah signature metode di luar `wails3 dev`, buat ulang binding tersebut secara manual.

### Peristiwa tidak dipicu

Nama peristiwa harus sama persis antara `app.Event.Emit("name", ...)` di Go dan `Events.On("name", ...)` di JavaScript. Periksa terlebih dahulu kesalahan ketik dan perbedaan huruf besar-kecil.

### Saya menemukan bug

Silakan [buat issue](https://github.com/wailsapp/wails/issues) dan sertakan output `wails3 doctor` Anda. [Panduan pemberian masukan](/feedback/) menjelaskan hal-hal yang membuat laporan mudah ditindaklanjuti.

## Bermigrasi dari v2

### Haruskah saya bermigrasi dari v2 ke v3?

v3 menghadirkan dukungan multi-jendela, API berbasis layanan yang lebih rapi, pemutakhiran bawaan, sistem build yang jauh lebih fleksibel, dan performa yang lebih baik. Proyek baru sebaiknya dimulai dengan v3. Untuk proyek yang sudah ada, [Panduan Migrasi](/migration/v2-to-v3/) menjelaskan perbedaannya langkah demi langkah.

### Apakah v2 akan tetap dipelihara?

Ya. Wails v2 tetap menerima pembaruan penting. Proyek baru sebaiknya menggunakan Wails v3.0.0.

### Dapatkah saya menjalankan v2 dan v3 secara berdampingan?

Ya. CLI-nya merupakan biner terpisah (`wails` dan `wails3`), sedangkan modulnya memiliki jalur impor yang berbeda. Karena itu, proyek dengan versi mayor yang berbeda dapat berjalan berdampingan tanpa masalah pada satu mesin.

## Komunitas

### Bagaimana cara mendapatkan bantuan?

- [Discord](https://discord.gg/JDdSxwjhGf) untuk pertanyaan singkat dan diskusi
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) untuk pertanyaan yang memerlukan pembahasan lebih panjang
- [GitHub Issues](https://github.com/wailsapp/wails/issues) untuk bug

### Bagaimana cara berkontribusi?

Lihat [Panduan Kontribusi](/contributing/). Perbaikan bug selalu diterima; fungsionalitas baru dan perubahan perilaku publik menggunakan draf PR [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md). Diskusi informal di Discord atau GitHub Discussions bersifat opsional.

### Di mana saya dapat menemukan contoh?

Repositori ini menyediakan lebih dari 60 contoh yang dapat dijalankan, mencakup jendela, dialog, peristiwa, baki sistem, layanan, dan lainnya: [v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).

## Masih Punya Pertanyaan?

Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau [buat diskusi](https://github.com/wailsapp/wails/discussions).
