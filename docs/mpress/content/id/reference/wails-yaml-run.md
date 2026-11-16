---
title: "Konfigurasi run YAML"
description: "Default run lokal, penimpaan khusus platform, dan perilaku peluncuran untuk wails3 run."
---

## Ketersediaan

Referensi ini menjelaskan ekstensi run eksperimental. Ini bukan jaminan bahwa ekstensi tersebut tersedia dalam CLI yang telah dirilis.

`wails3 run` membangun dan meluncurkan aplikasi menggunakan default manifes lokal. Eksekusi berbasis manifes mengharuskan pengaktifan [`WAILS_EXP_USE_WAKE=1`](/experimental/yaml-builds#enable-yaml-builds). Jika manifes tidak ditemukan, perintah menjalankan `go run .` di direktori saat ini tanpa mengharuskan pengaktifan tersebut. Manifes yang tidak valid atau tidak dapat dibaca merupakan kesalahan; kondisi tersebut tidak pernah memilih fallback. Jika manifes tersedia saat YAML dinonaktifkan, perintah akan menjelaskan cara mengaktifkannya.

Referensi kolom [`wails.yaml` yang dihasilkan](/reference/wails-yaml-fields) mencantumkan tipe dan nilai default kolom. Penerimaan paket native untuk perangkat seluler dan macOS harus diverifikasi pada host yang sesuai sebelum suatu rilis menyatakan platform tersebut siap.

## Kontrak perintah

`wails3 run` menemukan `wails.yaml` terdekat, menentukan konfigurasi untuk satu target, membangun melalui pipeline cache yang ada, lalu meluncurkan hasilnya. Penemuan menelusuri dari direktori kerja ke atas, memeriksa manifes sebelum berhenti pada batas `go.mod`. Direktori manifes menjadi root proyek dan direktori kerja aplikasi desktop.

Target default adalah sistem operasi dan arsitektur host. Target seluler harus dipilih secara eksplisit. Kompilasi silang tidak berarti bahwa executable desktop dapat dijalankan pada host; target desktop yang tidak kompatibel akan gagal sebelum pembangunan dimulai. Hal ini mencakup executable yang dibangun dengan `CGO_ENABLED=0`: keberhasilan build silang khusus Go tidak menyediakan emulator atau mengubah kontrak peluncuran lokal. Gunakan [`wails3 build` untuk kompilasi silang](/experimental/yaml-builds/#cross-compile-desktop-applications), lalu jalankan atau instal file hasilnya pada sistem target. Gunakan `--target android/arm64 --device <adb-serial>` atau `--target android/arm64 --emulator <avd-name>` untuk Android. Arsitektur yang dipilih harus sesuai dengan perangkat atau emulator. Gunakan `--target ios/arm64 --device <udid>` untuk simulator iOS dan tambahkan `--destination device` untuk perangkat fisik. Perencanaan dan eksekusi iOS memerlukan macOS dan Xcode.

Dengan YAML, build menggunakan aset tersemat dan jalur build produksi, yang secara otomatis menyediakan tag `production`. `wails3 dev` tetap menangani pemuatan ulang langsung; `dev.tags` tidak berkontribusi pada `wails3 run`.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan` menentukan dan menampilkan target yang dipilih, tag build efektif, artefak executable atau yang dapat diinstal, argumen, serta direktori kerja tanpa membangun, menginstal, atau meluncurkan. Rencana mencantumkan nama kunci penimpaan lingkungan dan tidak menampilkan nilainya.

## Opsi CLI

| Opsi | Arti |
| --- | --- |
| `--tags <tags>` | Tag build Go tambahan yang dipisahkan dengan koma. |
| `--skip-build` | Impor executable desktop yang sudah ada dari direktori keluaran YAML, lalu jalankan tanpa melakukan kompilasi. |
| `--target <os>/<arch>` | Pilih satu target; default-nya adalah OS dan arsitektur host. |
| `--plan` | Tampilkan pengaturan yang telah ditentukan tanpa membangun atau meluncurkan. Rencana desktop lintas platform diperbolehkan. |
| `--device <id>` | Nomor seri adb Android atau pengenal simulator/perangkat iOS. |
| `--emulator <name>` | Android Virtual Device yang akan dimulai atau digunakan kembali; tidak dapat digabungkan dengan `--device`. |
| `--destination <kind>` | Tujuan iOS: `simulator` (default) atau `device`. |
| `-- <args...>` | Ganti argumen aplikasi yang dikonfigurasi. `--` tanpa nilai akan menghapusnya. |
| `--appargs <string>` | Ejaan kompatibilitas: uraikan satu string menggunakan pengutipan bergaya POSIX. Tidak dapat diulang atau digabungkan dengan argumen setelah `--`. |

Pemilih seluler memerlukan manifes dan target seluler yang sesuai.

`--skip-build` memerlukan proyek YAML yang aktif dan executable desktop yang kompatibel. Opsi ini melewati pekerjaan frontend, binding, kompilasi, UPX, dan hook build; argumen peluncuran serta lingkungan yang dikonfigurasi tetap berlaku. Executable yang diimpor divalidasi dan disalin ke ruang kerja yang dihasilkan Wake tanpa mengubah file sumbernya. Aplikasi GUI macOS dirakit dari masukan tersebut dan aset paket saat ini. Lihat [jalur dan persyaratan masukan bawaan](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Tag build tidak dapat mengubah executable bawaan, sehingga `--tags` dan tag build yang dikonfigurasi tidak berpengaruh dalam mode ini.

## Kolom

Pemetaan opsional `run` menyediakan default peluncuran bersama. Pemetaan opsional `run` di dalam `targets["<os>"]` atau `targets["<os>/<arch>"]` menyediakan default khusus platform atau arsitektur. Pengaturan khusus arsitektur lebih diprioritaskan daripada pengaturan tingkat OS, terlepas dari urutan pemetaan.

| Kolom | Tipe | Default | Arti |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | Tag Go tambahan yang hanya digunakan oleh build run. |
| `run.args` | `list(string)` | `[]` | Argumen yang diteruskan ke aplikasi. |
| `run.environment` | `map(string)` | `{}` | Penimpaan lingkungan runtime, terpisah dari lingkungan compiler. |
| `target["<os>"].run.tags` | `list(string)` | `[]` | Tag run tambahan untuk platform yang dipilih. |
| `target["<os>"].run.args` | `list(string)` | diwarisi | Mengganti argumen bersama jika ditetapkan secara eksplisit; `[]` akan menghapusnya. |
| `target["<os>"].run.environment` | `map(string)` | `{}` | Menimpa nilai lingkungan runtime bersama berdasarkan kunci. |

Argumen diteruskan sebagai vektor argumen tanpa ekspansi shell. Nilai lingkungan merupakan string literal berdasarkan kontrak YAML yang hanya mengizinkan nilai literal.

## Contoh platform

Fragmen ini menyertai konfigurasi proyek normal. `example_private_api` adalah placeholder untuk tag API privat aktual yang diperlukan oleh contoh.

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Pengaturan lingkungan Linux bersifat khusus untuk contoh ini dan memerlukan sesi Wayland; pengaturan tersebut bukan default Wails. Tag API privat hanya dipilih di macOS. `target.tags` yang ada berlaku untuk build produksi, termasuk build run. `target.run.tags` membatasi tag tambahan agar hanya digunakan oleh build run.

Penimpaan khusus arsitektur dapat menyempurnakan default tingkat OS:

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

Jika digabungkan dengan pemetaan Linux sebelumnya, konfigurasi ini mempertahankan `GDK_BACKEND`, menambahkan tag khusus arsitektur, dan menghapus argumen aplikasi bersama.

## Aturan resolusi

| Pengaturan | Urutan resolusi, dari prioritas terendah hingga tertinggi |
| --- | --- |
| Tag | Gabungkan `build.tags`, `target.tags` yang dipilih, `production` otomatis, `run.tags`, `target.run.tags` yang dipilih, dan `--tags` CLI; pertahankan kemunculan pertama setiap tag. |
| Lingkungan runtime | Warisi shell yang meluncurkan, lapisi dengan `run.environment`, lalu lapisi dengan `target.run.environment` yang dipilih. Nilai kosong menetapkan string kosong. |
| Argumen aplikasi | `run.args` bersama, digantikan oleh argumen target yang ditetapkan secara eksplisit, lalu digantikan oleh argumen eksplisit setelah `--`. `--` tanpa nilai menghapus argumen yang dikonfigurasi. |

Dalam target yang dipilih, pengaturan tingkat OS ditentukan sebelum pengaturan khusus arsitektur: tag diakumulasikan, nilai lingkungan digabungkan berdasarkan kunci, dan argumen yang diberikan secara eksplisit menggantikan daftar yang kurang spesifik. Build run juga menyediakan `android` atau `ios` untuk target seluler yang sesuai.

Pembuatan binding dan kompilasi Go menerima tag hasil resolusi yang sama. Kunci cache build menyertakan tag tersebut. Perubahan hanya pada argumen aplikasi atau lingkungan runtime tidak membatalkan kompilasi; peluncuran berikutnya menggunakan nilai baru.

Pemetaan platform menyesuaikan suatu platform; keberadaannya tidak menjadikan platform lain tidak didukung. `project.supported_platforms` membatasi platform yang diterima oleh build dan `run`; daftar yang tidak dicantumkan atau kosong mengizinkan semua platform. Sebagai contoh, `supported_platforms: [darwin]` di dalam `project` menyebabkan contoh khusus macOS gagal lebih awal di Linux.

## Perilaku peluncuran platform

| Platform | Perilaku peluncuran |
| --- | --- |
| Linux dan Windows | Luncurkan executable dengan input/output standar yang terhubung, penanganan terminasi, dan penerusan status keluar aplikasi. |
| macOS | Susun bundle aplikasi dan jalankan biner `Contents/MacOS`-nya dengan input/output standar yang terhubung. Daftarkan bundle tersebut pada Launch Services untuk skema URL dan asosiasi file. Penandatanganan dan penerimaan siklus hidup native memerlukan verifikasi macOS. |
| Android | Build APK, pilih perangkat atau emulator, lalu instal dan luncurkan aplikasi. Laporkan kegagalan build, instalasi, dan peluncuran secara terpisah. |
| iOS | Build bundle aplikasi, pilih simulator atau perangkat, penuhi persyaratan penandatanganan dan provisioning, lalu instal dan luncurkan. Build untuk perangkat memerlukan kredensial dan toolchain Apple yang sesuai. |

Android menolak argumen aplikasi dan penggantian lingkungan runtime sebelum deployment. Jika Android juga didukung, pengaturan bersama yang hanya berlaku untuk desktop harus ditempatkan dalam pemetaan target desktop. Proses yang dijalankan di Android melampirkan log aplikasi dan menghentikan aplikasi saat pelampiran berakhir. iOS meneruskan argumen ke launcher yang dipilih. Simulator menerima penggantian lingkungan melalui variabel `SIMCTL_CHILD_`; perangkat fisik menerima peta JSON melalui `devicectl --environment-variables`. Status keluar pada perangkat seluler menjelaskan hasil launcher atau pemantauan log, bukan kode keluar aplikasi yang portabel.

## Contoh khusus Go dan fallback

Contoh khusus Go dengan aset yang sudah ada atau disematkan menonaktifkan tahapan frontend:

```yaml
frontend:
  disabled: true
```

Tindakan ini melewati instalasi dependensi frontend, pembuatan binding, dan build frontend. Sebagai gantinya, contoh yang memerlukan kelas binding yang dihasilkan dapat menetapkan `interfaces: false` di bawah `frontend.bindings`.

Fallback tanpa manifes menjalankan `go run [-tags <tags>] . [application arguments]` tanpa menambahkan `production` atau `devtools`. Fallback ini mengikuti perilaku status keluar Go sendiri. Argumen setelah `--` diteruskan ke aplikasi, termasuk flag seperti `--help`. Nilai runtime yang ditampilkan oleh `--plan` dibatasi hanya pada nama kunci lingkungan; nilainya tidak dicetak. Rencana fallback menampilkan perintah Go dan direktori kerja.

## Argumen aplikasi pengembangan

Dengan proyek YAML yang aktif, `wails3 dev` menerima argumen aplikasi secara independen dari nilai default `run` produksi. `dev.args` menyediakan nilai default bersama; pemetaan `dev` bertingkat dalam `targets["<os>"]` atau `targets["<os>/<arch>"]` menggantikannya untuk target tersebut. Pengaturan arsitektur lebih diprioritaskan daripada pengaturan tingkat OS, terlepas dari urutan peta. `args: []` yang ditetapkan secara eksplisit menghapus argumen yang diwarisi.

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args` dan `run.tags` tidak digunakan dalam peluncuran pengembangan. Kebijakan build `dev.tags` yang ada tidak berubah. Path aplikasi relatif ditafsirkan oleh aplikasi dari direktori manifes.

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

Argumen setelah `--` menggantikan argumen yang dikonfigurasi, termasuk saat daftarnya kosong. `--help` setelah pembatas tersebut merupakan milik aplikasi. `--appargs=""` juga menghapus nilai default. Opsi kompatibilitas menerima satu string yang dikutip; menggunakannya bersama `--`, argumen aplikasi posisional, atau `--appargs` kedua merupakan kesalahan. Pengutipan yang tidak valid ditolak sebelum sesi dimulai.

Parser kompatibilitas menggunakan aturan pengutipan bergaya POSIX yang sama pada setiap host. Parser ini tidak menjalankan perintah shell atau memperluas variabel; shell pemanggil mungkin telah memperluas teks sebelum meneruskannya ke Wails. Argumen setelah `--` sudah berupa nilai terpisah dan tidak mengalami tokenisasi tambahan. Bentuk ini lebih disarankan untuk path Windows, misalnya di PowerShell:

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

Sesi mempertahankan vektor argumen yang telah di-resolve selama build ulang sumber. Mengubah hanya `dev.args` atau argumen target yang dipilih akan memulai ulang backend menggunakan executable yang di-cache; tindakan ini tidak memulai ulang frontend. Penggantian dari CLI tetap berlaku setelah manifes dimuat ulang. Penggantian yang gagal mempertahankan, atau memulihkan, executable sebelumnya beserta argumen sebelumnya. Argumen aplikasi bukan flag compiler dan tidak menjadi bagian dari kunci cache kompilasi. `dev --plan` menampilkan argumen efektif bersama rencana build terbatas tanpa melakukan peluncuran.

Linux, Windows, dan macOS meneruskan vektor secara langsung ke executable. iOS meneruskannya melalui launcher simulator atau perangkat. Android menolak argumen aplikasi pengembangan yang tidak kosong sebelum penemuan perangkat atau build karena tidak tersedia saluran argumen setara yang didukung. Nilai default khusus desktop dapat ditempatkan dalam pemetaan target desktop, atau dihapus dengan `args: []` di bawah `targets.android.dev`.

Opsi pengembangan ini memerlukan `wails.yaml` yang aktif dan keikutsertaan melalui YAML. Opsi ini tidak menulis ulang perintah shell Taskfile atau watcher dalam jalur pengembangan lama.
