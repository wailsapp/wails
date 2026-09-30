---
title: "Perjalanan Menuju Wails v3"
description: "Catatan rilis dan pengumuman untuk Wails"
authors: ["leaanthony"]
tags: ["wails","v3"]
date: "2023-01-17"
slug: "blog/the-road-to-wails-v3"
image: "/assets/blog-images/multiwindow.webp"
sourcePath: "blog/the-road-to-wails-v3.md"
---

![tangkapan layar multi-jendela](/assets/blog-images/multiwindow.webp)

## Pendahuluan

Wails adalah proyek yang mempermudah pembuatan aplikasi desktop lintas platform menggunakan Go. Wails menggunakan komponen webview native untuk frontend (bukan browser tersemat), sehingga menghadirkan kemampuan sistem UI paling populer di dunia ke Go sekaligus tetap ringan.

Versi 2 dirilis pada 22 September 2022 dan menghadirkan banyak penyempurnaan, termasuk:

- Pengembangan langsung dengan memanfaatkan proyek Vite yang populer
- Fungsionalitas lengkap untuk mengelola jendela dan membuat menu
- Komponen WebView2 dari Microsoft
- Pembuatan model TypeScript yang mencerminkan struct Go Anda
- Pembuatan penginstal NSIS
- Build yang diobfusikasi

Saat ini, Wails v2 menyediakan alat canggih untuk membuat aplikasi desktop lintas platform yang kaya fitur.

Tulisan blog ini bertujuan meninjau kondisi proyek saat ini dan hal-hal yang dapat kita tingkatkan ke depannya.

## Di mana posisi kita sekarang?

Luar biasa melihat popularitas Wails terus meningkat sejak rilis v2. Saya selalu takjub melihat kreativitas komunitas dan berbagai hal mengagumkan yang dibuat dengannya. Seiring meningkatnya popularitas, semakin banyak orang memperhatikan proyek ini. Akibatnya, semakin banyak pula permintaan fitur dan laporan bug.

Seiring waktu, saya berhasil mengidentifikasi beberapa masalah paling mendesak yang dihadapi proyek ini. Saya juga berhasil mengidentifikasi beberapa hal yang menghambat kemajuannya.

## Masalah saat ini

Saya telah mengidentifikasi beberapa area berikut yang menurut saya menghambat kemajuan proyek ini:

- API
- Pembuatan binding
- Sistem Build

### API

API untuk membuat aplikasi Wails saat ini terdiri atas 2 bagian:

- API Aplikasi
- API Runtime

API Aplikasi dikenal hanya memiliki 1 fungsi: `Run()`, yang menerima banyak opsi untuk mengatur cara kerja aplikasi. Meskipun sangat mudah digunakan, pendekatan ini juga sangat membatasi. Ini adalah pendekatan "deklaratif" yang menyembunyikan banyak kompleksitas yang mendasarinya. Misalnya, tidak ada handle untuk jendela utama sehingga Anda tidak dapat berinteraksi langsung dengannya. Untuk itu, Anda harus menggunakan API Runtime. Hal ini menjadi masalah ketika Anda ingin melakukan hal yang lebih kompleks, seperti membuat beberapa jendela.

API Runtime menyediakan banyak fungsi utilitas bagi pengembang. Ini mencakup:

- Pengelolaan jendela
- Dialog
- Menu
- Peristiwa
- Log

Ada sejumlah hal yang tidak saya sukai dari API Runtime. Pertama, API ini mengharuskan "context" untuk diteruskan dari satu tempat ke tempat lain. Hal ini membuat frustrasi sekaligus membingungkan bagi pengembang baru yang meneruskan context, tetapi kemudian mendapatkan error runtime.

Masalah terbesar pada API Runtime adalah bahwa API ini dirancang untuk aplikasi yang hanya menggunakan satu jendela. Seiring waktu, kebutuhan akan beberapa jendela meningkat, dan API ini tidak cocok untuk kebutuhan tersebut.

### Pemikiran tentang API v3

Bukankah akan sangat baik jika kita dapat melakukan sesuatu seperti ini?

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

Pendekatan programatik ini jauh lebih intuitif dan memungkinkan pengembang berinteraksi langsung dengan elemen aplikasi. Semua metode runtime yang ada saat ini untuk jendela cukup dijadikan metode pada objek jendela. Untuk metode runtime lainnya, kita dapat memindahkannya ke objek aplikasi seperti berikut:

```go
app := wails.NewApplication(options.App{})
app.NewInfoDialog(options.InfoDialog{})
app.Log.Info("Hello World")
```

Ini adalah API yang jauh lebih canggih dan memungkinkan pembuatan aplikasi yang lebih kompleks. API ini juga memungkinkan pembuatan beberapa jendela, [fitur dengan suara terbanyak di GitHub](https://github.com/wailsapp/wails/issues/1480):

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    myWindow2 := app.NewWindow(options.Window{})
    myWindow2.SetTitle("My Window 2")
    myWindow2.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

### Pembuatan binding

Salah satu fitur utama Wails adalah membuat binding untuk metode Go Anda agar dapat dipanggil dari JavaScript. Metode yang digunakan saat ini untuk melakukannya agak bersifat tambal sulam. Prosesnya mencakup pembuatan build aplikasi dengan flag khusus, lalu menjalankan biner hasilnya, yang menggunakan refleksi untuk menentukan apa saja yang telah di-binding. Hal ini menimbulkan situasi seperti ayam dan telur: Anda tidak dapat membuat build aplikasi tanpa binding, dan Anda tidak dapat membuat binding tanpa membuat build aplikasi. Ada banyak cara untuk mengatasinya, tetapi cara terbaik adalah tidak menggunakan pendekatan ini sama sekali.

Ada sejumlah upaya untuk menulis penganalisis statis bagi proyek Wails, tetapi upaya tersebut tidak berkembang jauh. Belakangan ini, hal itu menjadi sedikit lebih mudah dilakukan karena tersedia lebih banyak materi mengenai topik tersebut.

Dibandingkan dengan refleksi, pendekatan AST jauh lebih cepat, tetapi juga jauh lebih rumit. Sebagai langkah awal, kita mungkin perlu menetapkan batasan tertentu pada cara menentukan binding di dalam kode. Tujuannya adalah mendukung kasus penggunaan yang paling umum, lalu memperluas dukungannya nanti.

### Sistem Build

Seperti pendekatan deklaratif pada API, sistem build dibuat untuk menyembunyikan kerumitan pembuatan build aplikasi desktop. Saat Anda menjalankan `wails build`, sistem ini melakukan banyak hal di balik layar:

- Membuat biner backend untuk binding dan menghasilkan binding
- Menginstal dependensi frontend
- Membuat aset frontend
- Memeriksa apakah ikon aplikasi tersedia dan, jika tersedia, menyematkannya
- Membuat biner akhir
- Jika build ditujukan untuk `darwin/universal`, sistem membuat 2 biner, satu untuk `darwin/amd64` dan satu untuk `darwin/arm64`, lalu membuat fat binary menggunakan `lipo`
- Jika kompresi diperlukan, sistem mengompresi biner dengan UPX
- Memeriksa apakah biner ini akan dikemas dan, jika ya:
  - Memastikan ikon dan manifes aplikasi dikompilasi ke dalam biner (Windows)
  - Membuat bundel aplikasi, menghasilkan bundel ikon dan menyalinnya beserta biner dan Info.plist ke bundel aplikasi (Mac)

- Jika penginstal NSIS diperlukan, sistem membuatnya

Seluruh proses ini, meskipun sangat andal, juga sangat tidak transparan. Proses ini sangat sulit disesuaikan dan sangat sulit di-debug.

Untuk mengatasi hal ini di v3, saya ingin beralih ke sistem build yang berada di luar Wails. Setelah menggunakan [Task](https://taskfile.dev/) selama beberapa waktu, saya menjadi penggemar beratnya. Ini adalah alat yang sangat baik untuk mengonfigurasi sistem build dan seharusnya cukup familier bagi siapa pun yang pernah menggunakan Makefile.

Sistem build akan dikonfigurasi menggunakan sebuah file `Taskfile.yml` yang secara default akan dibuat bersama setiap templat yang didukung. File ini akan memuat semua langkah yang diperlukan untuk menjalankan seluruh tugas yang tersedia saat ini, seperti mem-build atau mengemas aplikasi, sehingga mudah disesuaikan.

Tidak akan ada persyaratan eksternal untuk alat ini karena alat tersebut akan menjadi bagian dari Wails CLI. Artinya, Anda tetap dapat menggunakan `wails build` dan perintah tersebut akan melakukan semua yang dilakukannya saat ini. Namun, jika ingin menyesuaikan proses build, Anda dapat melakukannya dengan mengedit file `Taskfile.yml`. Ini juga berarti Anda dapat dengan mudah memahami langkah-langkah build dan menggunakan sistem build sendiri jika menginginkannya.

Bagian yang masih belum tersedia dalam rangkaian proses build adalah operasi atomik dalam proses tersebut, seperti pembuatan ikon, kompresi, dan pengemasan. Mengharuskan penggunaan sekumpulan alat eksternal bukanlah pengalaman yang baik bagi pengembang. Untuk mengatasinya, Wails CLI akan menyediakan semua kemampuan tersebut sebagai bagian dari CLI. Artinya, proses build tetap berjalan sebagaimana mestinya tanpa alat eksternal tambahan, tetapi Anda dapat mengganti langkah mana pun dalam proses build dengan alat apa pun yang Anda inginkan.

Sistem build ini akan jauh lebih transparan, lebih mudah disesuaikan, dan mengatasi banyak masalah yang telah dilaporkan terkait sistem tersebut.

## Hasilnya

Perubahan positif ini akan memberikan manfaat besar bagi proyek:

- API baru akan jauh lebih intuitif dan memungkinkan pembuatan aplikasi yang lebih kompleks.
- Penggunaan analisis statis untuk menghasilkan binding akan jauh lebih cepat dan mengurangi banyak kerumitan dalam proses saat ini.
- Penggunaan sistem build eksternal yang telah mapan akan membuat proses build sepenuhnya transparan sehingga dapat disesuaikan secara luas.

Manfaat bagi pengelola proyek meliputi:

- API baru akan jauh lebih mudah dipelihara dan diadaptasi untuk fitur serta platform baru.
- Sistem build baru akan jauh lebih mudah dipelihara dan diperluas. Saya berharap hal ini akan menghasilkan ekosistem baru berupa pipeline build yang digerakkan oleh komunitas.
- Pemisahan tanggung jawab yang lebih baik di dalam proyek. Hal ini akan memudahkan penambahan fitur dan platform baru.

## Rencana

Banyak eksperimen untuk hal ini telah dilakukan dan hasilnya terlihat baik. Saat ini belum ada jadwal untuk pekerjaan ini, tetapi saya berharap pada akhir Q1 2023 akan tersedia rilis alfa untuk Mac agar komunitas dapat menguji, bereksperimen, dan memberikan masukan.

## Ringkasan

- API v2 bersifat deklaratif, menyembunyikan banyak hal dari pengembang, dan tidak cocok untuk fitur seperti beberapa jendela. API baru akan dibuat dengan desain yang lebih sederhana, intuitif, dan kemampuan yang lebih besar.
- Sistem build tidak transparan dan sulit disesuaikan, sehingga kami akan beralih ke sistem build eksternal yang akan membuat seluruh prosesnya terbuka.
- Pembuatan binding berlangsung lambat dan rumit, sehingga kami akan beralih ke analisis statis yang akan menghilangkan banyak kerumitan pada metode saat ini.

Banyak upaya telah dicurahkan pada bagian inti v2 dan fondasinya sudah kokoh. Kini saatnya membenahi lapisan di atasnya dan memberikan pengalaman yang jauh lebih baik bagi pengembang.

Saya harap Anda sama antusiasnya dengan saya mengenai hal ini. Saya menantikan pendapat dan masukan Anda.

Salam,

&dash; Lea

PS: Jika Anda atau perusahaan Anda merasakan manfaat Wails, harap pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!

PPS: Ya, itu benar-benar tangkapan layar aplikasi berjendela majemuk yang dibuat dengan Wails. Itu bukan mockup. Itu nyata. Itu luar biasa. Segera hadir.
