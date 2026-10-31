---
title: "Build YAML"
description: "Konfigurasikan dan jalankan sistem build native Wails v3 yang eksperimental dengan satu file wails.yaml."
sidebar:
  order: 2
  badge:
    text: "Eksperimental"
    variant: caution
---

@note{type="caution" title="Pratinjau eksperimental"}
Rilis pertama yang diterima hanya mendukung konfigurasi literal. Eksperimen ini juga membuat prototipe hook siklus hidup berbatas agar bentuknya dapat dievaluasi sebelum diusulkan untuk dirilis. Tahap kustom, ekspresi YAML, dan pemanggilan alat Wails bertipe masih ditangguhkan. Proyek yang sudah ada tetap menggunakan Taskfile sampai Anda secara eksplisit mengaktifkan migrasi yang telah ditinjau.
@end

## Aktifkan build YAML

Tetapkan `WAILS_EXP_USE_WAKE` di setiap terminal yang menjalankan perintah eksperimental. Untuk Bash atau Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Untuk PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Variabel ini mengaktifkan fitur berdasarkan keberadaannya: bahkan nilai kosong atau `false` akan mengaktifkannya. Hapus variabel tersebut (`unset WAILS_EXP_USE_WAKE` di Bash/Zsh atau `Remove-Item Env:WAILS_EXP_USE_WAKE` di PowerShell) untuk menonaktifkan eksperimen. Setelah `wails.yaml` aktif, perintah mengharuskan Wake diaktifkan kembali dan tidak pernah beralih kembali ke Taskfile.
Tanpa variabel tersebut, perutean build/dev YAML dan pembuatan proyek YAML dinonaktifkan, sedangkan perintah `migrate`, `config`, dan `eject` tidak tersedia. Proyek yang belum dimigrasikan tetap diproses melalui runtime Task tersemat yang terpisah.

## Konfigurasikan proyek

Setelah opsi tersebut diaktifkan, tambahkan `wails.yaml` di root proyek. Keberadaannya memilih sistem build native dan membuat Wails mengabaikan setiap Taskfile dalam proyek tersebut. Biarkan opsi ini aktif untuk semua perintah di bawah. Modeline skema mengaktifkan pelengkapan dan validasi sebaris pada editor yang mendukung server bahasa YAML.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Jalankan build host default atau permintaan rilis bernama yang lengkap:

```bash
wails3 build
wails3 build release
```

Untuk build sekali pakai, pilih target dan format yang dipisahkan koma tanpa menambahkan profil:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Kompilasi silang aplikasi desktop

Untuk backend desktop pure-Go eksperimental, tetapkan `CGO_ENABLED: "0"` di lingkungan build. Wake memilih tag `purego` untuk program GUI macOS dan Linux serta menggunakan Go secara langsung pada host Windows, macOS, dan Linux. Dalam mode ini, pembuatan executable desktop tidak memerlukan Docker, emulasi, compiler C, ataupun SDK platform target. Dependensi aplikasi juga harus mendukung `CGO_ENABLED=0`.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

Perintah kedua menggabungkan kedua arsitektur macOS dan membuat DMG tanpa tanda tangan menggunakan Go. Perintah ini dapat dijalankan pada host desktop apa pun yang didukung. Penandatanganan kode dan notarisasi Apple tetap merupakan operasi macOS yang terpisah.

Mode hasil resolusi dipilih secara independen untuk setiap executable. `build.environment.CGO_ENABLED` yang dideklarasikan, lingkungan target, atau lingkungan biner menggantikan nilai terminal yang menjalankan perintah. Jika tidak, tag `purego` yang dipilih secara eksplisit, program CLI/server, atau target Windows secara default menggunakan `0`; program GUI macOS dan Linux secara default menggunakan `1`. `--plan --json` melaporkan compiler, mode CGo, dan tag hasil resolusi untuk setiap executable. Kompilasi silang membuat file target; menjalankan target desktop yang tidak kompatibel tetap memerlukan host yang sesuai.

Dengan `CGO_ENABLED=1`, Wake terlebih dahulu menggunakan toolchain native yang kompatibel atau cross-compiler `CC` yang dikonfigurasi, kemudian Zig untuk build silang Windows jika tersedia, lalu fallback Docker/Podman yang sudah ada jika didukung. `CC` yang ditetapkan secara eksplisit harus menyediakan compiler, pustaka, dan sysroot target; Zig saja tidak menyediakan sysroot GTK/WebKit Linux. Pemaketan dan penandatanganan mungkin memerlukan alat tambahan sebagaimana dijelaskan di bawah. Target seluler memerlukan CGo dan SDK platformnya.

Kompilasi khusus Go tidak membundel pustaka sistem milik platform tujuan.
Aplikasi GUI Linux, termasuk executable pendamping GUI, tetap memerlukan runtime
GTK/WebKitGTK yang dipilih oleh tag build-nya. Dependensi paket atau lingkungan
deployment yang sesuai harus menyediakannya. Secara default, paket yang hanya
berisi CLI/server tidak menambahkan dependensi GUI tersebut; penggantian
dependensi secara eksplisit tetap tersedia.

## Build beberapa executable

Proyek menentukan program utama. Tetapkan `project.package` ke paket Go-nya (default `.`) dan `project.kind` ke `gui`, `cli`, atau `server` (default `gui`). Program tambahan menggunakan entri bernama di bawah `binaries`:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Setiap program memiliki input kompilasi dan cache sendiri. Paket proyek yang diimpor dan file tersemat dilacak; perubahan pada program independen tidak menyebabkan program lain dibangun ulang. Executable tambahan didukung di desktop serta disertakan dalam bundel aplikasi dan penginstal default. File paket pengganti milik pengguna mengendalikan kontennya sendiri.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Profil bernama merupakan permintaan lengkap; gunakan penggantian pilihan ini dengan build anonim. `--jobs` membatasi worker eksekusi. Memilih hanya program CLI/server akan melewatkan pekerjaan frontend.

Untuk membaca versi tanpa menjalankan kode aplikasi, konfigurasikan `project.version_file` dan `project.version_symbol`. File tersebut harus berisi deklarasi string Go statis. Literal `project.version` kemudian dapat dihilangkan. Tetapkan `commit_symbol: main.Commit` pada proyek atau biner bernama untuk menyisipkan revisi repositori ke dalam variabel string executable tersebut. Tidak adanya metadata repositori merupakan kesalahan ketika penyisipan commit diminta.

`upx_level: 1` hingga `9` secara opsional mengompresi executable rilis Linux dan Windows; nilai nol menonaktifkannya. UPX harus terinstal. Build pengembangan dan macOS tidak menggunakan kompresi. Kegagalan kompilasi atau kompresi mempertahankan program yang sebelumnya dipublikasikan.

## Paketkan dan tandatangani untuk distribusi

Pembuatan paket tanpa tanda tangan memiliki persyaratan host/alat yang berbeda dari penandatanganan:

| Keluaran | Host dan alat pembuatan | Penandatanganan |
| --- | --- | --- |
| MSIX Windows | Windows, macOS, atau Linux; penulis Go, tanpa memerlukan Windows SDK. | Host Windows dan alat penandatanganan. |
| Penginstal NSIS Windows | Windows, macOS, atau Linux; memerlukan `makensis` native untuk host. | Host Windows dan alat penandatanganan. |
| Aplikasi macOS dan DMG | Windows, macOS, atau Linux; perakitan dengan Go, penggabungan biner universal, dan penulis DMG portabel. | macOS dengan alat dan kredensial Apple. |
| DEB, RPM, Arch, dan Void Linux | Windows, macOS, atau Linux; penulis paket native Go. | DEB/RPM menggunakan alat penandatanganan eksternalnya; Arch menggunakan penanda tangan Go di bawah. Keluaran Void tidak ditandatangani. |
| AppImage Linux | Host Linux dan arsitektur target yang cocok; linuxdeploy eksternal dan plugin deployment GTK. | Tidak ada alur penandatanganan AppImage bawaan. |

Arsip MSIX Windows dibuat langsung dengan Go; pembuatan paket tanpa tanda tangan
tidak memerlukan Windows SDK atau Docker. Penandatanganan Windows tetap memerlukan
host Windows beserta alat penandatanganannya. NSIS menggunakan compiler NSIS
eksternal dan menerima definisi literal melalui `packages.nsis.defines`. Di macOS, Wails
menyediakan locale UTF-8 hanya kepada proses turunan compiler NSIS jika locale
yang diwarisinya bukan UTF-8; lingkungan terminal dan aplikasi tidak berubah.
Installer NSIS bawaan tidak menyertakan bootstrapper WebView2 jika setiap
executable yang dikemas adalah CLI atau server. Executable utama atau pendamping
berbasis GUI tetap menyertakan WebView2. Skrip NSIS kustom tetap menyertakan
bootstrapper demi kompatibilitas dengan logika instalasi runtime miliknya sendiri.

Linux mendukung paket AppImage, DEB, RPM, Arch Linux, dan Void. AppImage tetap mempertahankan batas integrasi dengan linuxdeploy eksternal dan plugin GTK. Tetapkan `packages.appimage.update_info` untuk meminta file zsync dan `packages.appimage.checksum: true` untuk file SHA256. File sidecar dipulihkan dari cache paket dan dicantumkan dalam bukti artefak bersama AppImage. Pemaketan AppImage saat ini memerlukan Linux serta arsitektur host/target yang cocok. Deteksi GTK-nya tidak menemukan pustaka yang dimuat secara dinamis oleh GUI pure-Go, sehingga adaptor ini belum mendukung deployment GUI pure-Go. Gunakan salah satu format pengelola paket Linux untuk mode build tersebut.

Paket Arch mempertahankan prarilis dalam metadata paket native; misalnya, `1.2.3-beta.2` menjadi `1.2.3beta.2-1`, yang diurutkan sebelum rilis final. Dependensi Void menerima nama paket polos seperti `gtk4`, yang dinormalisasi menjadi `gtk4>=0`; perbandingan versi eksplisit dan nilai `name-version_revision` yang lengkap dipertahankan. Gunakan nama paket Void, bukan nama file pustaka.

Untuk membuat paket Arch dan tanda tangan terpisah, gunakan file kunci privat PGP berformat armored yang diekspor dan, jika perlu, tentukan nama variabel lingkungan frasa sandi:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Tetapkan variabel lingkungan yang disebutkan melalui lingkungan lokal atau penyimpanan rahasia CI Anda, lalu jalankan `wails3 build release`. Hilangkan `credential` untuk kunci yang tidak dienkripsi. Penanda tangan Go berjalan pada setiap host desktop yang didukung dan tidak memerlukan GPG. Penanda tangan tersebut menerbitkan `.pkg.tar.zst` asli dan `.pkg.tar.zst.sig` di lokasi yang sama; keduanya dicatat secara terpisah dalam tanda terima artefak. Isi kunci tidak ditempatkan sementara atau disimpan dalam cache build. Tanda tangan dan penerbitannya tidak pernah digunakan kembali dari cache. Jika penandatanganan gagal, pasangan yang sebelumnya diterbitkan tetap dipertahankan, dan penerbitan mengembalikan kedua berkas jika pemasangan salah satunya gagal. Hal ini tidak membuat dua operasi penggantian nama pada sistem berkas menjadi atomik bagi pembaca serentak.

Bundel aplikasi macOS menyertakan program pendamping. Tetapkan `darwin.microphone_usage_description` untuk menambahkan teks penggunaan dan, jika berkas entitlement khusus tidak disediakan, entitlement input audio. Sumber plist dan entitlement khusus tetap tidak berubah. Penandatanganan dan notarisasi menggunakan alat serta kredensial Apple yang sudah ada di macOS; notarisasi memerlukan profil notaris yang sebenarnya. Output yang ditandatangani mempertahankan ekstensi platformnya, misalnya `notes.signed.app` dan `notes.signed.dmg`.

## Kontrol pengembangan

Selama `wails3 dev` berlangsung, masukkan perintah yang diikuti Enter:

| Perintah | Tindakan |
| --- | --- |
| `r` | Paksa build ulang dan mulai ulang, termasuk backend yang tidak berubah. |
| `b` | Buat ulang binding dan lakukan build ulang. |
| `f` | Mulai ulang frontend. |
| `c` | Bersihkan terminal. |
| `q` | Hentikan sesi beserta proses turunannya. |
| `h` | Tampilkan kontrol. |

Build yang gagal membiarkan aplikasi saat ini tetap berjalan. Kontrol memerlukan terminal interaktif; input yang disalurkan melalui pipe tidak mengubah perilaku sesi.

## Periksa pekerjaan sebelum menjalankannya

Pemeriksaan rencana me-resolve graf tetap yang sama dengan yang digunakan oleh eksekusi. Teks ditujukan bagi manusia; JSON merupakan kontrak yang dapat dibaca mesin dan memiliki versi.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Validasi konfigurasi tanpa melakukan build

Gunakan perintah preflight khusus di CI, dalam tugas editor, atau sebelum melakukan commit konfigurasi. Tanpa profil, perintah tersebut memvalidasi manifes dasar dan setiap profil bernama; dengan profil, perintah tersebut hanya memeriksa permintaan lengkap itu:

```bash
wails3 config check
wails3 config check release
```

Kegagalan sintaks, skema, semantik, jalur, target, format, penandatanganan, perencana, dan kemampuan host yang ditentukan konfigurasi menggunakan perender diagnostik yang sama dengan Build dan Dev. Diagnostik mengidentifikasi bidang dan rentang sumber, menampilkan baris sumber beserta tanda sisipan, serta menyertakan petunjuk perbaikan jika berguna. Nilai lingkungan disamarkan dari kutipan sumber diagnostik. Kegagalan host yang hanya disebabkan oleh pilihan baris perintah tetap dilaporkan sebagai galat perintah biasa, bukan secara keliru sebagai YAML yang tidak valid.

## Gunakan aset milik proyek

Jalur manifes merujuk ke berkas milik Anda, seperti ikon, manifes platform, templat paket, dan skrip pengelola. Wails tidak pernah mengedit sumber tersebut. Wails menyalinnya ke ruang kerja `.wails/` sementara sebelum digunakan oleh suatu alat. Artefak akhir ditempatkan di bawah `build.output` dan tanda terima terverifikasi ditulis ke `.wails/artifacts/receipt.json`.

## Deploy ke Android

Profil Android produksi menghasilkan AAB. Untuk deployment lokal, Wails secara internal membuat APK debug khusus ABI, memasangnya dengan ADB, lalu meluncurkan aplikasi:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Tanpa pemilih eksplisit, terminal interaktif menawarkan perangkat yang terhubung dan Android Virtual Device yang telah dikonfigurasi. Otomatisasi dengan beberapa target yang tersedia harus menyediakan `--device` atau `--emulator`.

Lampirkan aliran log dalam cakupan aplikasi dan pantau target hingga aplikasi berhenti atau Anda menekan Ctrl+C:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Proses yang dijalankan dengan lampiran akan mengikuti mulai ulang aplikasi dan melaporkan terputusnya koneksi target atau peralihan ke status offline. Secara default, emulator tetap berjalan agar iterasi dapat dilakukan dengan cepat. `--stop-emulator` bersifat eksplisit, memerlukan `--logs`, dan hanya menghentikan target emulator ketika sesi terlampir berakhir. Jika Wails memulai emulator baru tetapi kesiapan boot gagal atau dibatalkan, Wails otomatis membersihkan proses tersebut.

## Jalankan hook siklus hidup yang dibatasi

Antarmuka hook eksperimental menjalankan satu skrip milik proyek pada setiap fase yang dideklarasikan. Hook bukan perintah shell: Wails mengeksekusi berkas tersebut secara langsung dari root proyek, kecuali jika `directory` ditetapkan.

Pada host Unix, skrip harus dapat dieksekusi dan diawali dengan shebang. Host Windows menerima skrip `.cmd`, `.bat`, dan `.ps1` melalui adaptor interpreter tetap; Wails tidak mengurai perintah shell inline pada kedua platform.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

Fase yang didukung adalah `before_build`, `after_build`, `before_package`, `after_package`, `before_sign`, dan `after_sign`. Hook build berjalan satu kali untuk proyek bersama atau satu kali per target, sesuai kebutuhan. Hook paket dan penandatanganan membentuk penghalang di sekeliling semua format yang diminta untuk suatu target, sementara operasi format di dalam penghalang tetap berjalan secara serentak.

Hook menerima satu variabel lingkungan milik Wails: `WAILS_HOOK_CONTEXT_FILE`. Variabel ini menyebutkan berkas JSON terlindungi yang berisi konteks pemanggilan berversi:

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

Konteks tersebut berisi jalur yang telah di-resolve, tetapi tidak pernah menyalin lingkungan proses yang diwarisi maupun rahasia. Wails membuatnya secara atomik dengan izin berkas khusus pemilik di bawah `.wails/hooks/context/`. Berkas tersebut dihapus setelah berhasil atau dibatalkan. Jika terjadi kegagalan, berkas dipertahankan untuk diagnosis dan jalurnya disertakan dalam galat. Status keluar bukan nol menggagalkan langkah Plan normal, dan pembatalan menghentikan grup proses hook.

Secara default, hook berjalan pada setiap pemanggilan. Aktifkan caching hanya jika skrip memiliki input dan output lengkap yang dideklarasikan:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

Skrip dan mode executable-nya merupakan input implisit. Beberapa output harus menggunakan satu direktori non-root yang sama, yang menjadi artefak terbatas yang dicatat dan dipulihkan oleh Wails. Root output tidak boleh berisi skrip atau input yang dideklarasikan, dan pohon output yang di-cache hanya boleh berisi berkas serta direktori biasa, tidak pernah tautan simbolis atau berkas khusus. Hook yang selalu dijalankan membuat operasi turunannya tidak dapat di-cache karena efek samping yang tidak dideklarasikan tidak dapat diberi sidik jari secara aman. Hook yang dapat di-cache mempertahankan caching turunan melalui identitas artefak yang dideklarasikan.

## Migrasikan proyek yang sudah ada

Buat dan tinjau draf yang tidak aktif:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Tugas siklus hidup konservatif tanpa argumen yang menjalankan skrip milik proyek diterjemahkan menjadi hook. Perilaku Taskfile khusus lain yang dapat dijangkau dilaporkan sebagai penghambat. Setelah meninjau `wails.migrated.yaml`, aktifkan secara eksplisit:

```bash
wails3 migrate --activate
```

Aktivasi membuat `wails.yaml` secara atomik. Aktivasi tidak menghapus, mengganti nama, atau mengedit Taskfile, `build/config.yml`, skrip, maupun aset. YAML yang tidak valid gagal sebagai YAML; proses tidak pernah kembali ke rute lama.

## Materialisasikan semua nilai default

Untuk memeriksa atau menyalin konfigurasi lengkap yang telah di-resolve, tulis snapshot yang tidak aktif:

```bash
wails3 eject
```

Tindakan ini membuat `wails.ejected.yaml`. Tindakan ini tidak pernah menulis ulang konfigurasi aktif.

## Jendela kompatibilitas

`wails3 package` dan `wails3 sign` adalah alias usang untuk perencana yang sama; otomatisasi baru sebaiknya menggunakan `wails3 build` dengan profil atau format. Peluncuran dibatasi berdasarkan rilis: rilis pertama menjadikan YAML sebagai jalur proyek baru, rilis berikutnya memperingatkan tentang perutean Taskfile lama, dan rilis setelahnya menghapus eksekusi Taskfile beserta alias tersebut. Setiap transisi mengharuskan matriks penerimaan host native untuk rilis tersebut berstatus lulus.

Lihat referensi bidang [`wails.yaml` yang dihasilkan](/reference/wails-yaml-fields/) untuk setiap bidang, tipe, nilai default, dan aturan penerapan.

Referensi konfigurasi [`wails3 run`](/reference/wails-yaml-run/) menjelaskan nilai default peluncuran lokal, penggantian khusus platform, dan fallback `go run .` tanpa manifes.
