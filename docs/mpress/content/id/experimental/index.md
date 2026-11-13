---
title: "Fitur Eksperimental"
description: "Tempat untuk eksperimen Wails v3 yang masih dikembangkan—mencakup apa itu, mengapa eksperimen tersebut dibuat, dan cara memberikan umpan balik."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Di sinilah eksperimen berada"}
Semua yang ada di bagian ini, menurut definisinya, merupakan eksperimen. Fitur-fitur ini bersifat opsional, dinonaktifkan secara default, dan bentuk atau namanya dapat berubah, bahkan dapat dihapus sepenuhnya, dari satu rilis ke rilis berikutnya. Jangan mengandalkannya untuk komponen penting kecuali Anda siap menghadapi perubahan yang berulang.

@end

## Arti "Eksperimental"

Wails merilis eksperimen secara terbuka. Eksperimen adalah gagasan yang menurut kami cukup menjanjikan untuk diberikan kepada Anda lebih awal, tetapi belum sepenuhnya kami tetapkan. Kami memublikasikannya *karena* ingin mempelajari penggunaan nyata sebelum memutuskan apakah eksperimen tersebut akan menjadi bagian permanen Wails yang didukung.

Artinya, beberapa hal berikut berlaku untuk semua yang ada di bagian ini:

- **Fitur ini harus diaktifkan secara eksplisit.** Eksperimen tidak pernah mengubah perilaku default `wails3`. Anda mengaktifkannya dengan sengaja (biasanya menggunakan variabel lingkungan atau flag build), dan jika dinonaktifkan, alur kerja Anda yang sudah ada tidak berubah.
- **Fitur ini mungkin tidak dipertahankan.** Beberapa eksperimen berkembang menjadi fitur stabil. Eksperimen lainnya dirombak hingga sama sekali berbeda atau dihentikan. Kami lebih memilih mencoba berbagai hal secara terbuka dan belajar dengan cepat daripada hanya merilis hal yang sudah kami yakini sepenuhnya.
- **API belum dibekukan.** Nama, flag, nilai default, dan perilaku dapat berubah antar-rilis selama bentuk eksperimen masih disempurnakan. Catatan rilis akan menyebutkan perubahan tersebut, tetapi jangan mengharapkan jaminan stabilitas seperti yang diberikan untuk fitur stabil.

## Kami menginginkan umpan balik Anda

Inilah bagian yang penting. Kelangsungan eksperimen bergantung pada tanggapan dari orang-orang yang benar-benar menggunakannya. Jika Anda mencoba salah satunya, kami sungguh ingin mengetahui:

- Apakah eksperimen tersebut berfungsi untuk proyek Anda? Di bagian mana eksperimen itu gagal?
- Apakah eksperimen tersebut lebih cepat, lebih jelas, atau lebih nyaman—atau tidak sepadan dengan peralihannya?
- Apa yang harus terpenuhi agar Anda menggunakannya secara default?

Umpan balik yang paling berguna bersifat konkret: apa yang Anda jalankan, apa yang Anda harapkan, dan apa yang sebenarnya terjadi. Setiap eksperimen memiliki utas tersendiri dalam kategori **Eksperimen** di GitHub Discussions:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Diskusi Eksperimen" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Temukan utas untuk eksperimen yang Anda gunakan, lalu beri tahu kami hasilnya, apa yang rusak, atau apa yang masih kurang."}
@end

## Eksperimen saat ini

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Cobalah sistem build YAML native: buat dan konfigurasikan proyek dalam satu wizard, build dan paketkan proyek tersebut, atau tinjau migrasi satu arah dari Taskfile."}
@linkcard{title="Build tanpa CGO" href="/experimental/#cgo-free-builds-purego" description="Build aplikasi macOS dan Linux dengan CGO_ENABLED=0: tanpa toolchain C dan dapat melakukan kompilasi silang tanpa Docker. Aktifkan dengan WAILS_PUREGO=1 atau tag build purego."}
@linkcard{title="Kontrol LLM (MCP)" href="/guides/mcp-service/" description="Server Model Context Protocol bawaan yang memungkinkan agen LLM memeriksa, menguji, dan mengendalikan aplikasi Wails yang sedang berjalan—tanpa memerlukan kode pengguna dan diaktifkan dengan tag build."}
@end

## Build tanpa CGO (purego)

@note{type="caution" title="Fitur Eksperimental"}
Backend tanpa CGO harus diaktifkan secara eksplisit melalui tag build `purego` (atau `WAILS_PUREGO=1`) dan **bukan** merupakan pilihan default. Tanpanya, build macOS dan Linux menggunakan backend cgo persis seperti sebelumnya. Perilakunya dapat berubah antar-rilis.

@end

### Apa itu

Di macOS dan Linux, Wails biasanya berkomunikasi dengan pustaka UI sistem (AppKit dan WKWebView di macOS, serta GTK dan WebKitGTK di Linux) melalui cgo. Backend tanpa CGO melakukan pemanggilan yang sama dari Go murni menggunakan [purego](https://github.com/ebitengine/purego), yang memuat pustaka sistem saat runtime.

API aplikasi publik tidak berubah. Pilih backend saat melakukan build; kode aplikasi tidak perlu diubah. Uji aplikasi Anda pada sistem tujuan sebelum mendistribusikannya.

Manfaat yang Anda peroleh:

- **Tidak memerlukan toolchain C untuk kompilasi Go.** Build Go tidak memerlukan alat baris perintah Xcode, `gcc`, atau header pengembangan Linux. Build frontend tetap memerlukan alat JavaScript yang biasanya digunakan.
- **Kompilasi silang tanpa Docker.** Build biner macOS atau Linux dari sistem operasi apa pun hanya dengan `go build`.

Windows sudah dapat di-build tanpa cgo dan tidak terpengaruh.

### Sebelum mengaktifkannya

Gunakan versi CLI dan framework yang sama-sama menyertakan eksperimen ini. Memperbarui CLI saja tidak menambahkan backend ke dependensi framework lama dalam `go.mod` aplikasi Anda. Untuk mencoba versi pengembangan, jalankan perintah berikut dari direktori aplikasi Anda:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

Kedua perintah mengikuti branch `master` yang terus berubah. Untuk build yang dapat direproduksi, gunakan commit teruji yang sama untuk CLI dan framework. Ikuti [panduan instalasi](/getting-started/installation/) untuk menyiapkan toolchain Go dan CLI.

Backend eksperimental macOS dan Linux mendukung `amd64` dan `arm64`. Kompilasi silang menghasilkan executable untuk target yang diminta; proses tersebut tidak menjalankan executable itu pada mesin build. Uji executable tersebut pada sistem operasi dan arsitektur target.

### Cara menggunakannya

Pilih salah satu cara berikut untuk menentukan backend.

#### Dengan Wails CLI

Dalam proyek yang dibuat menggunakan versi CLI ini, `build`, `package`, dan `dev` menggunakan backend eksperimental ketika `WAILS_PUREGO` diaktifkan. Dalam shell POSIX:

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

Dalam PowerShell:

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

CLI mengenali `1`, `true`, `on`, dan `yes` tanpa membedakan huruf besar-kecil serta mengabaikan spasi di sekitarnya. Nilai kosong tidak mengaktifkan eksperimen ini.

Atau teruskan tag tersebut ke satu perintah build atau pemaketan:

```sh
wails3 build -tags purego
wails3 package -tags purego
```

Saat tag ditetapkan, Taskfile standar yang telah diperbarui melakukan build dengan `CGO_ENABLED=0` menggunakan toolchain Go milik host. Sebagai contoh, build executable Linux ARM64 dari macOS atau Windows:

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

Proyek lama tetap menggunakan Taskfile yang sudah ada. Bandingkan tugas build khusus dengan templat yang telah diperbarui, atau gunakan Go biasa untuk menguji backend tanpa mengubah tugas tersebut. Variabel lingkungan memilih tag dalam perintah Wails CLI; `go build` biasa memerlukan tag tersebut secara eksplisit.

#### Dengan Go biasa

Build frontend Anda terlebih dahulu agar semua aset yang dirujuk oleh `go:embed` tersedia, lalu kompilasi aplikasinya:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego` juga dapat digunakan: backend-nya tetap sama, sedangkan cgo tetap tersedia untuk bagian lain program Anda.

Untuk build silang, tetapkan juga `GOOS` dan `GOARCH`:

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

Di Linux, pilih stack GTK3 dengan `-tags purego,gtk3,production`.

Untuk kembali ke backend default, hapus tag `purego` dan batalkan penetapan `WAILS_PUREGO` (`unset WAILS_PUREGO` dalam shell POSIX atau `Remove-Item Env:WAILS_PUREGO` dalam PowerShell). Hapus juga penetapan eksplisit `CGO_ENABLED=0` saat kembali ke backend cgo macOS atau Linux.

Jika Anda melakukan build dengan `CGO_ENABLED=0` tetapi lupa menyertakan tag, build akan berhenti dengan satu pesan kesalahan yang menyebutkan cara memperbaikinya:

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### Persyaratan runtime

Aplikasi memuat pustaka sistem saat dimulai, sehingga pustaka tersebut harus terinstal pada mesin target, sama seperti pada build cgo:

| Platform | Pustaka | Versi paling lama yang didukung |
|---|---|---|
| macOS | AppKit, WebKit (bagian dari OS) | Sama seperti backend cgo |
| Linux (default) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40 (Debian 12) |
| Linux (`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36 (Ubuntu 22.04) |

Fungsi yang lebih baru daripada versi minimum ini dicari saat aplikasi dimulai dan hanya digunakan jika tersedia di pustaka yang terinstal. Dengan demikian, satu berkas biner dapat berjalan pada sistem lama maupun baru.

Aplikasi tetap menggunakan pustaka UI native dan dapat menautkan pustaka sistem secara dinamis, seperti libc di Linux. `CGO_ENABLED=0` tidak berarti berkas executable yang sepenuhnya statis. Jika pustaka runtime yang diperlukan tidak tersedia, instal stack GTK/WebKit yang sesuai pada sistem tujuan; menginstal header C pada mesin build tidak memperbaiki ketiadaan pustaka runtime.

### Pemaketan dan penandatanganan

Backend ini menghilangkan kebutuhan akan toolchain C saat kompilasi. Setiap format paket tetap memiliki persyaratannya sendiri. Secara khusus, memilih `purego` tidak membuat alat deployment AppImage yang ada menjadi lintas platform. Lihat [pemaketan Linux](/guides/build/linux/).

Aplikasi macOS yang di-build secara silang tetap memerlukan penandatanganan dan notarisasi yang sesuai sebelum didistribusikan. Installer Windows tetap memerlukan alat pembuatan installernya. Lihat [pemaketan macOS](/guides/build/macos/) dan [pemaketan Windows](/guides/build/windows/). Untuk dependensi yang menggunakan cgo, `CGO_ENABLED=1` juga memerlukan compiler C serta semua header atau pustaka untuk target tersebut; memilih backend pure-Go Wails tidak menghilangkan persyaratan dependensi tersebut.

### Validasi

Pemeriksaan binding yang dihasilkan dan `puregocheck` memeriksa simbol native, guard versi opsional, dan cakupan backend. Alur kerja paritas Linux menguji GTK4 dan GTK3 dengan kedua backend, termasuk permintaan aset, isi permintaan biner, dan dispatch native; alur ini juga melakukan kompilasi silang untuk target Linux ARM64. Pemeriksaan kompilasi silang memastikan bahwa kode dapat di-build, sedangkan pemeriksaan runtime menguji pustaka native yang terinstal pada mesin pengujian.

[Panduan pengujian backend](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md) menjelaskan probe paritas dan pengujian input menyeluruh. Saat melaporkan masalah, bedakan kegagalan kompilasi dari kegagalan runtime atau pemaketan.

### Perbedaan yang diketahui

- **Akselerator menu di Linux menggunakan pendaftaran pada tingkat aplikasi.** Hindari menetapkan akselerator yang berkonflik ke menu yang berbeda. Menghapus menu sebuah jendela secara eksplisit dengan `SetMenu(nil)` juga mencegah pintasan menu yang telah dihapus tersebut aktif di jendela itu; jendela lain tetap mempertahankan menu dan pintasannya.
- **Panel buka dan simpan macOS** berjalan dalam proses sistem terpisah, seperti halnya dengan cgo; pengujian memastikan panel tersebut dapat dibuka dan dibatalkan dengan benar.

Menemukan hal lain? Beri tahu kami.

### Umpan balik

Ceritakan pengalaman Anda di utas **build tanpa CGO** pada bagian [Eksperimen](https://github.com/wailsapp/wails/discussions/categories/experiments): OS beserta versinya, cara Anda melakukan build, serta apa yang berhasil atau tidak. Kasus kegagalan yang disertai output `wails3 doctor` merupakan laporan paling berguna yang dapat kami terima.
