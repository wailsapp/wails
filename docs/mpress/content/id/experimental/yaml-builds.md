---
title: "Build YAML"
description: "Konfigurasikan dan jalankan sistem build native Wails v3 eksperimental dengan satu file wails.yaml."
sidebar:
  order: 2
  badge:
    text: "Eksperimental"
    variant: caution
---

@note{type="caution" title="Pratinjau eksperimental"}
Rilis pertama yang diterima hanya mendukung konfigurasi literal. Eksperimen ini juga membuat purwarupa hook siklus hidup berbatas agar bentuknya dapat dievaluasi sebelum diusulkan untuk dirilis. Tahap kustom, ekspresi YAML, dan pemanggilan alat Wails bertipe masih ditangguhkan. Proyek yang sudah ada tetap menggunakan Taskfile sampai Anda secara eksplisit mengaktifkan migrasi yang telah ditinjau.
@end

## Mengaktifkan build YAML

Tetapkan `WAILS_EXP_USE_WAKE` di setiap terminal yang menjalankan perintah eksperimental. Untuk Bash atau Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Untuk PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Variabel ini mengaktifkan fitur berdasarkan keberadaannya: nilai kosong sekalipun atau `false` akan mengaktifkannya. Hapus variabel tersebut (`unset WAILS_EXP_USE_WAKE` di Bash/Zsh atau `Remove-Item Env:WAILS_EXP_USE_WAKE` di PowerShell) untuk menonaktifkan eksperimen. Setelah `wails.yaml` aktif, perintah mengharuskan Wake diaktifkan kembali dan tidak pernah kembali menggunakan Taskfile.
Tanpa variabel tersebut, perutean build/dev YAML dan pembuatan proyek YAML dinonaktifkan, serta perintah `migrate`, `config`, dan `eject` tidak tersedia. Proyek yang belum dimigrasikan tetap diproses melalui runtime Task tertanam yang terpisah.

## Mengonfigurasi proyek

Setelah pilihan ikut serta diaktifkan, tambahkan `wails.yaml` di akar proyek. Keberadaan file tersebut memilih sistem build native dan membuat Wails mengabaikan setiap Taskfile dalam proyek itu. Pertahankan pilihan ikut serta tetap aktif untuk semua perintah di bawah ini. Modeline skema mengaktifkan pelengkapan dan validasi sebaris di editor yang mendukung server bahasa YAML.

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

Untuk build sekali jalan, pilih target dan format yang dipisahkan koma tanpa menambahkan profil:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Mengompilasi silang aplikasi desktop

Untuk backend desktop eksperimental yang sepenuhnya menggunakan Go, tetapkan `CGO_ENABLED: "0"` di lingkungan build. Wake memilih tag `purego` untuk program GUI macOS dan Linux serta menggunakan Go secara langsung pada host Windows, macOS, dan Linux. Dalam mode ini, membangun executable desktop tidak memerlukan Docker, emulasi, compiler C, atau SDK platform target. Dependensi aplikasi juga harus mendukung `CGO_ENABLED=0`.

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

Mode hasil resolusi dipilih secara independen untuk setiap executable. `build.environment.CGO_ENABLED` yang dideklarasikan, lingkungan target, atau lingkungan biner menggantikan nilai terminal yang menjalankan perintah. Jika tidak, tag `purego` yang dipilih secara eksplisit, program CLI/server, atau target Windows menggunakan `0` secara default; program GUI macOS dan Linux menggunakan `1` secara default. `--plan --json` melaporkan compiler, mode CGo, dan tag hasil resolusi untuk setiap executable. Kompilasi silang membuat file target; menjalankan target desktop yang tidak kompatibel tetap memerlukan host yang sesuai.

Dengan `CGO_ENABLED=1`, Wake terlebih dahulu menggunakan toolchain native yang kompatibel atau cross-compiler `CC` yang dikonfigurasi, kemudian Zig untuk build silang Windows jika tersedia, lalu fallback Docker/Podman yang sudah ada jika didukung. `CC` yang ditetapkan secara eksplisit harus menyediakan compiler target, pustaka, dan sysroot; Zig saja tidak menyediakan sysroot GTK/WebKit Linux. Pemaketan dan penandatanganan mungkin memerlukan alat tambahan seperti dijelaskan di bawah. Target seluler memerlukan CGo dan SDK platformnya.

Kompilasi khusus Go tidak membundel pustaka sistem milik tujuan. Aplikasi GUI Linux, termasuk executable pendamping GUI, tetap memerlukan runtime GTK/WebKitGTK yang dipilih oleh tag build-nya. Dependensi paket atau lingkungan deployment yang sesuai harus menyediakannya. Secara default, paket yang hanya berisi CLI/server tidak menambahkan dependensi GUI tersebut; penggantian dependensi secara eksplisit tetap tersedia.

## Membangun beberapa executable

Proyek mendefinisikan program utama. Tetapkan `project.package` ke paket Go-nya (default `.`) dan `project.kind` ke `gui`, `cli`, atau `server` (default `gui`). Program tambahan menggunakan entri bernama di bawah `binaries`:

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

Setiap program memiliki input kompilasi dan cache tersendiri. Paket proyek yang diimpor dan file tertanam dilacak; pengeditan pada program independen tidak membangun ulang program lainnya. Executable tambahan didukung di desktop serta disertakan dalam bundel aplikasi dan penginstal default. File paket pengganti milik pengguna mengendalikan kontennya sendiri.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Profil bernama merupakan permintaan lengkap; gunakan penggantian pilihan ini dengan build anonim. `--jobs` membatasi worker eksekusi. Memilih hanya program CLI/server akan melewatkan pekerjaan frontend.

Untuk membaca versi tanpa menjalankan kode aplikasi, konfigurasikan `project.version_file` dan `project.version_symbol`. File tersebut harus berisi deklarasi string Go statis. Literal `project.version` kemudian dapat dihilangkan. Tetapkan `commit_symbol: main.Commit` pada proyek atau biner bernama untuk menyuntikkan revisi repositori ke variabel string executable tersebut. Ketiadaan metadata repositori merupakan kesalahan ketika penyuntikan commit diminta.

`upx_level: 1` hingga `9` secara opsional mengompresi executable rilis Linux dan Windows; nol menonaktifkannya. UPX harus terinstal. Build pengembangan dan macOS tidak menggunakan kompresi. Kegagalan kompilasi atau kompresi mempertahankan program yang sebelumnya dipublikasikan.

## Memaketkan dan menandatangani untuk distribusi

### Menggunakan executable desktop yang telah dibangun sebelumnya

Gunakan `--skip-build` jika sistem build atau tugas CI lain telah menghasilkan executable:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Tempatkan setiap executable di `build.output` (nilai defaultnya `bin`), dengan menggunakan `project.binary_name` untuk executable utama dan nama setiap entri `binaries` untuk executable pendamping. Nama untuk Windows diakhiri dengan `.exe`. Pemanggilan dengan beberapa target membaca file setiap target dari subdirektori `<os>-<arch>`. Untuk GUI macOS, Wake membaca executable mentah di `bin/<name>` jika tersedia; jika tidak, dari `bin/<name>.app/Contents/MacOS/<name>`. Input `darwin/universal` harus berisi slice executable amd64 dan arm64.

Wake memvalidasi jenis file dan arsitektur target, menyalin executable ke ruang kerja yang dihasilkannya, serta melacak kontennya untuk cache paket. Input executable yang tidak ada, tidak kompatibel, atau berupa symlink menyebabkan kegagalan alih-alih membangun ulang atau memulihkan input lama. Executable yang diberikan tidak ditimpa atau diubah modenya dengan chmod. Perubahan pada executable mana pun membatalkan cache paket yang menggunakannya.

Mode ini melewatkan instalasi/build frontend, binding, kompilasi Go, UPX, serta hook `before_build`/`after_build`. Mode ini tetap menghasilkan aset paket, menyusun bundel aplikasi macOS, dan menjalankan hook pemaketan/penandatanganan yang dikonfigurasi. Alat pemaketan dan penandatanganan tetap tunduk pada tabel di bawah; Go dan alat build frontend tidak diperlukan. Metadata paket berasal dari manifes saat ini, termasuk `version_file` jika dikonfigurasi; versi tidak disimpulkan dari biner. Gunakan pengaturan build dan paket yang sesuai, termasuk jenis GUI versus CLI.

`wails3 package --skip-build` juga didukung. Dengan `wails3 run`, argumen peluncuran dan lingkungan tetap berlaku, dan executable harus dapat berjalan pada host. Mode ini memerlukan proyek YAML aktif dan hanya mendukung target desktop produksi; mode ini tidak berlaku untuk sesi pengembangan atau pemaketan seluler.

Pembuatan paket tanpa tanda tangan memiliki persyaratan host/alat yang berbeda dari penandatanganan:

| Output | Host dan alat pembuatan | Penandatanganan |
| --- | --- | --- |
| MSIX Windows | Windows, macOS, atau Linux; penulis Go, tidak memerlukan Windows SDK. | Host Windows dan alat penandatanganan. |
| Penginstal NSIS Windows | Windows, macOS, atau Linux; memerlukan `makensis` native host. | Host Windows dan alat penandatanganan. |
| Aplikasi macOS dan DMG | Windows, macOS, atau Linux; penyusunan dengan Go, penggabungan biner universal, dan penulis DMG portabel. | macOS dengan alat dan kredensial Apple. |
| DEB, RPM, Arch, dan Void Linux | Windows, macOS, atau Linux; penulis paket Go native. | DEB/RPM menggunakan alat penandatanganan eksternalnya; Arch menggunakan penanda tangan Go di bawah ini. Output Void tidak ditandatangani. |
| AppImage Linux | Host Linux dan arsitektur target yang cocok; linuxdeploy eksternal dan plugin deployment GTK. | Tidak ada alur bawaan untuk menandatangani AppImage. |

Arsip MSIX Windows dibuat langsung di Go; pembuatan paket yang tidak ditandatangani tidak memerlukan Windows SDK atau Docker. Penandatanganan Windows tetap memerlukan host Windows beserta alat penandatanganannya. NSIS menggunakan compiler NSIS eksternal dan menerima definisi literal melalui `packages.nsis.defines`. Di macOS, Wails memberikan locale UTF-8 hanya kepada proses turunan compiler NSIS jika locale yang diwarisinya bukan UTF-8; lingkungan terminal dan aplikasi tidak berubah. Penginstal NSIS standar menghilangkan bootstrapper WebView2 jika setiap executable dalam paket merupakan CLI atau server. Executable utama atau pendamping yang memiliki GUI tetap menyertakan WebView2. Skrip NSIS khusus tetap menyertakan bootstrapper demi kompatibilitas dengan logika penginstalan runtime-nya sendiri.

Linux mendukung paket AppImage, DEB, RPM, Arch Linux, dan Void. AppImage tetap menggunakan batas integrasi dengan linuxdeploy eksternal dan plugin GTK. Tetapkan `packages.appimage.update_info` untuk meminta file zsync dan `packages.appimage.checksum: true` untuk file SHA256. File sidecar dipulihkan dari cache paket dan dicantumkan dalam tanda terima artefak bersama AppImage. Saat ini, pemaketan AppImage memerlukan Linux serta arsitektur host dan target yang cocok. Deteksi GTK-nya tidak menemukan pustaka yang dimuat secara dinamis oleh GUI Go murni, sehingga adaptor ini belum mendukung deployment GUI Go murni. Gunakan salah satu format pengelola paket Linux untuk mode build tersebut.

Paket Arch mempertahankan prarilis dalam metadata paket native. Misalnya, `1.2.3-beta.2` menjadi `1.2.3beta.2-1`, yang diurutkan sebelum rilis final. Dependensi Void menerima nama paket tanpa keterangan tambahan seperti `gtk4`, yang dinormalisasi menjadi `gtk4>=0`; perbandingan versi eksplisit dan nilai `name-version_revision` lengkap dipertahankan. Gunakan nama paket Void, bukan nama file pustaka.

Untuk membuat paket Arch dan tanda tangan terpisah, gunakan file kunci privat PGP berformat armored yang telah diekspor dan, jika perlu, tentukan nama variabel lingkungan untuk frasa sandi:

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

Tetapkan variabel lingkungan bernama tersebut melalui lingkungan lokal atau penyimpanan rahasia CI Anda, lalu jalankan `wails3 build release`. Hilangkan `credential` untuk kunci yang tidak dienkripsi. Penanda tangan Go berjalan di setiap host desktop yang didukung dan tidak memerlukan GPG. Penanda tangan ini menerbitkan `.pkg.tar.zst` asli dan file `.pkg.tar.zst.sig` di lokasi yang sama; keduanya dicatat secara terpisah dalam tanda terima artefak. Isi kunci tidak ditempatkan di area staging atau disimpan dalam cache build. Tanda tangan dan penerbitannya tidak pernah digunakan kembali dari cache. Jika penandatanganan gagal, pasangan yang diterbitkan sebelumnya tetap dipertahankan. Jika penginstalan salah satu file gagal, penerbitan akan mengembalikan kedua file ke keadaan sebelumnya. Hal ini tidak membuat dua operasi penggantian nama pada sistem file menjadi atomik bagi pembaca konkuren.

Bundel aplikasi macOS menyertakan program pendamping. Tetapkan `darwin.microphone_usage_description` untuk menambahkan teks penggunaan dan, jika tidak ada file entitlement khusus yang diberikan, entitlement input audio. Sumber plist dan entitlement khusus tetap tidak berubah. Penandatanganan dan notarisasi menggunakan alat serta kredensial Apple yang sudah ada di macOS; notarisasi memerlukan profil notaris yang nyata. Output yang ditandatangani mempertahankan ekstensi platformnya, misalnya `notes.signed.app` dan `notes.signed.dmg`.

## Kontrol pengembangan

Selama `wails3 dev` berjalan, masukkan perintah lalu tekan Enter:

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

Pemeriksaan rencana me-resolve graf imutabel yang sama dengan yang digunakan oleh eksekusi. Teks ditujukan untuk manusia; JSON adalah kontrak terbaca mesin yang memiliki versi.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Validasi konfigurasi tanpa melakukan build

Gunakan perintah preflight khusus di CI, sebagai tugas editor, atau sebelum melakukan commit konfigurasi. Tanpa profil, perintah ini memvalidasi manifes dasar dan setiap profil bernama; dengan profil, perintah ini hanya memeriksa permintaan lengkap tersebut:

```bash
wails3 config check
wails3 config check release
```

Kegagalan sintaks, skema, semantik, jalur, target, format, penandatanganan, perencana, dan kemampuan host yang ditentukan konfigurasi menggunakan perender diagnostik yang sama dengan Build dan Dev. Diagnostik mengidentifikasi kolom dan rentang sumber, menampilkan baris sumber beserta tanda sisipan, serta menyertakan petunjuk koreksi jika berguna. Nilai lingkungan disamarkan dalam cuplikan sumber diagnostik. Kegagalan host yang hanya disebabkan oleh pilihan baris perintah tetap dilaporkan sebagai kesalahan perintah biasa, bukan secara keliru sebagai YAML yang tidak valid.

## Gunakan aset milik proyek

Jalur manifes merujuk pada file milik Anda, seperti ikon, manifes platform, templat paket, dan skrip maintainer. Wails tidak pernah mengedit sumber tersebut. Wails menyalinnya ke ruang kerja sementara `.wails/` sebelum digunakan oleh suatu alat. Artefak final ditempatkan di bawah `build.output` dan tanda terima terverifikasi ditulis ke `.wails/artifacts/receipt.json`.

## Deploy ke Android

Profil Android produksi menghasilkan AAB. Untuk deployment lokal, Wails membuat APK debug khusus ABI secara internal, menginstalnya dengan ADB, lalu meluncurkan aplikasi:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Tanpa pemilih eksplisit, terminal interaktif menawarkan perangkat yang terhubung dan Android Virtual Device yang telah dikonfigurasi. Otomatisasi dengan beberapa target yang tersedia harus memberikan `--device` atau `--emulator`.

Lampirkan aliran log dengan cakupan aplikasi dan pantau target hingga aplikasi berhenti atau Anda menekan Ctrl+C:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Proses yang dijalankan dalam mode terlampir akan mengikuti aplikasi ketika dimulai ulang dan melaporkan terputusnya koneksi target atau transisi ke status offline. Secara default, emulator tetap berjalan agar iterasi dapat dilakukan dengan cepat. `--stop-emulator` harus ditentukan secara eksplisit, memerlukan `--logs`, dan hanya menghentikan target emulator saat sesi terlampir berakhir. Jika Wails memulai emulator baru tetapi emulator gagal mencapai kesiapan boot atau proses dibatalkan, Wails otomatis membersihkan proses tersebut.

## Jalankan hook siklus hidup dengan batasan

Antarmuka hook eksperimental menjalankan satu skrip milik proyek pada setiap fase yang dideklarasikan. Hook bukan perintah shell: Wails mengeksekusi file secara langsung dari root proyek, kecuali jika `directory` ditetapkan.

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

Fase yang didukung adalah `before_build`, `after_build`, `before_package`, `after_package`, `before_sign`, dan `after_sign`. Hook build dijalankan satu kali untuk proyek bersama atau satu kali per target, sesuai kebutuhan. Hook paket dan penandatanganan membentuk penghalang yang mengapit semua format yang diminta untuk suatu target, sementara operasi format di dalam penghalang tetap berjalan secara konkuren.

Hook menerima satu variabel lingkungan milik Wails: `WAILS_HOOK_CONTEXT_FILE`. Variabel ini menunjuk ke file JSON terlindungi yang berisi konteks pemanggilan berversi:

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

Konteks tersebut memuat jalur yang telah di-resolve, tetapi tidak pernah menyalin lingkungan proses yang diwariskan maupun rahasia. Wails membuatnya secara atomik dengan izin file khusus pemilik di bawah `.wails/hooks/context/`. File tersebut dihapus setelah berhasil atau dibatalkan. Jika terjadi kegagalan, file dipertahankan untuk diagnosis dan jalurnya disertakan dalam pesan kesalahan. Kode keluar bukan nol menggagalkan langkah Plan normal, dan pembatalan menghentikan grup proses hook.

Secara default, hook berjalan pada setiap pemanggilan. Aktifkan caching hanya jika skrip memiliki semua input dan output yang dideklarasikan:

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

Skrip dan mode executable-nya merupakan input implisit. Beberapa output harus berbagi satu direktori non-root, yang menjadi artefak terbatas yang dicatat dan dipulihkan oleh Wails. Root output tidak boleh memuat skrip atau input yang dideklarasikan, dan pohon output yang di-cache hanya boleh memuat file biasa dan direktori, tidak pernah tautan simbolis atau file khusus. Hook yang selalu dijalankan membuat operasi setelahnya tidak dapat di-cache karena efek samping yang tidak dideklarasikan tidak dapat dibuat sidik jarinya secara aman. Hook yang dapat di-cache mempertahankan caching untuk operasi setelahnya melalui identitas artefak yang dideklarasikan.

## Migrasikan proyek yang sudah ada

Buat dan tinjau draf yang belum aktif:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Tugas siklus hidup konservatif tanpa argumen yang menjalankan skrip milik proyek
diterjemahkan menjadi hook. Perilaku Taskfile kustom lain yang dapat dijangkau
dilaporkan sebagai penghambat. Setelah meninjau
`wails.migrated.yaml`, aktifkan secara eksplisit:

```bash
wails3 migrate --activate
```

Aktivasi membuat `wails.yaml` secara atomik. Aktivasi tidak menghapus, mengganti nama, atau mengedit
Taskfile, `build/config.yml`, skrip, maupun aset. YAML yang tidak valid akan gagal sebagai YAML;
sistem tidak pernah beralih kembali ke jalur lama.

## Materialisasikan semua nilai default

Untuk memeriksa atau menyalin konfigurasi lengkap yang telah di-resolve, tulis snapshot
tidak aktif:

```bash
wails3 eject
```

Tindakan ini membuat `wails.ejected.yaml`. Tindakan ini tidak pernah menulis ulang konfigurasi aktif.

## Periode kompatibilitas

`wails3 package` dan `wails3 sign` adalah alias usang untuk perencana yang sama;
otomatisasi baru sebaiknya menggunakan `wails3 build` dengan profil atau format. Peluncurannya
dibatasi oleh rilis: rilis pertama menjadikan YAML sebagai jalur untuk proyek baru, rilis berikutnya
memperingatkan tentang perutean Taskfile lama, dan rilis selanjutnya menghapus eksekusi Taskfile
beserta alias tersebut. Setiap transisi mengharuskan matriks pengujian penerimaan host native
untuk rilis tersebut lulus.

Lihat referensi kolom [`wails.yaml` yang dihasilkan](/reference/wails-yaml-fields/)
untuk setiap kolom, tipe, nilai default, dan aturan penerapan.

Referensi konfigurasi [`wails3 run`](/reference/wails-yaml-run/) menjelaskan
nilai default peluncuran lokal, penggantian khusus platform, dan fallback `go run .` tanpa manifes.
