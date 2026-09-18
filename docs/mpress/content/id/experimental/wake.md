---
title: "Wake"
description: "Runner build eksperimental yang memahami Wails dan menjalankan Taskfile Anda yang sudah ada dengan build inkremental yang lebih cepat, keluaran terstruktur, dan eksekusi paralel secara default."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="Fitur Eksperimental"}
Wake harus diaktifkan secara khusus melalui `WAILS_USE_WAKE=true` dan **bukan** runner default. Jika variabel tersebut tidak ditetapkan, `wails3 build / package / sign / task` berperilaku persis seperti sebelumnya. Cakupan dan perilaku fitur dapat berubah antar-rilis.

@end

Wake adalah **runner build alternatif eksperimental** untuk `wails3`. Wake membaca `Taskfile.yml` yang sama dengan yang sudah dimiliki proyek Anda—dengan sintaks task, dep, var, template, include, dan namespace platform yang sama—lalu menjalankannya melalui eksekutor yang memahami Wails, bukan runtime [Task](https://taskfile.dev) serbaguna.

Tujuannya bukan menggantikan Task, melainkan menyediakan runner yang dibuat khusus untuk cara proyek Wails melakukan build, dengan semantik, keluaran, dan nilai default yang selaras dengan bagian lain CLI `wails3`. **Jika Anda hanya menggunakan Wake, Taskfile Anda tidak berubah.**

## Alasan keberadaannya

Wake dan runtime Task sama-sama dikompilasi ke dalam `wails3`—keduanya tidak memerlukan instalasi biner terpisah. Perbedaannya adalah Wake **memahami domainnya**. Runner serbaguna menjalankan langkah apa pun yang tercantum dalam Taskfile sesuai urutan yang ditentukan. Wake memahami apa *sebenarnya* build Wails itu—bundel frontend disematkan ke dalam biner, biner dikemas menjadi artefak khusus platform, serta ikon dan binding dibuat bersamaan—dan menggunakan pengetahuan tersebut untuk mengoptimalkan build dengan cara yang tidak dapat dilakukan runner generik.

- **Wake hanya mengerjakan hal yang benar-benar diperlukan oleh build.** Wake melacak sendiri masukan dan keluaran aktual setiap langkah. Untuk build Go, masukan dan keluaran tersebut adalah graf modul beserta keluaran langkah-langkah yang menjadi dependensinya. Jadi, jika tidak ada perubahan yang relevan, Wake melewati compiler dan linker sepenuhnya alih-alih menjalankannya kembali. Runner serbaguna hanya dapat melewati suatu langkah jika Taskfile telah mencantumkan terlebih dahulu berkas mana tepatnya yang harus dipantau; Wake menentukannya berdasarkan hal yang sudah diketahuinya tentang build. Untuk build ulang tanpa perubahan, waktunya sekitar **~20 ms (Wake) dibandingkan ~316 ms (Task)**. Build dari awal memiliki waktu tempuh yang sama karena sebagian besar waktunya digunakan oleh `npm install`, Vite, dan compiler Go.

- **Wake mengetahui langkah mana yang dapat dijalankan bersamaan.** Karena memahami langkah mana yang independen, Wake menjalankannya secara paralel secara default, dan baris hasil melaporkan peningkatan kecepatan yang diperoleh. Nonaktifkan dengan `WAKE_SERIAL=true` jika keluaran yang berselang-seling dari langkah-langkah sejajar akan mengaburkan penyelidikan.

- **Keluaran terstruktur yang dikendalikan oleh wails3.** Wake merender melalui pelapor milik wails3: satu baris untuk setiap langkah yang direncanakan, status langsung, perincian fase berkode warna di bagian akhir, dan tautan `file:line` yang dapat diklik di dalam panel kegagalan. `NO_COLOR` dan lingkungan non-TTY (log CI) ditangani dengan baik.

- **Terintegrasi sehingga dapat berkembang bersama Wails.** Karena Wake merupakan bagian dari `wails3` dan bukan alat pihak ketiga, kemampuan build baru dapat ditambahkan langsung tanpa harus menunggu proyek terpisah mengimplementasikannya. Hal ini juga membuka peluang untuk menjalankan skrip dan alat lintas platform secara native, yang saat ini dijalankan Taskfile melalui biner `wails3` (dengan membuat proses untuk setiap panggilan). Menjalankan pekerjaan tersebut dalam proses yang sama mengurangi overhead dan memungkinkan peningkatan kecepatan lebih lanjut pada masa mendatang.

## Mengaktifkan Wake

Wake sepenuhnya dibatasi oleh variabel lingkungan `WAILS_USE_WAKE=true`. Jika variabel tersebut tidak ditetapkan (atau ditetapkan ke nilai selain `true`), setiap perintah `wails3` menggunakan runtime Task tertanam persis seperti sebelumnya.

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

Flag ini mencakup `wails3 build`, `wails3 package`, `wails3 sign`, dan `wails3 task <name>`. `wails3 dev` **belum** terpengaruh—pemantau pengembangan masih menggunakan pipeline-nya sendiri.

@note{type="tip" title="Mengaktifkan Wake itu aman"}
Jika Wake menemukan fitur Taskfile yang belum diimplementasikannya, Wake menyerahkan seluruh proses eksekusi kepada runtime Task tertanam dalam proses yang sama—tidak ada biner `task` eksternal yang perlu diinstal. Dalam kondisi terburuk, Anda memperoleh perilaku yang persis sama seperti jika flag tersebut tidak digunakan.

@end

## Override lokal berlapis

Wake mendukung **Taskfile dasar beserta override lokal**. Tempatkan sebuah berkas di samping `Taskfile.yml` Anda, dan definisi di dalamnya akan diprioritaskan:

| Berkas | Tujuan | Prioritas |
| --- | --- | --- |
| `Taskfile.yml` | dasar, di-commit | terendah |
| `Taskfile.override.yml` / `.yaml` | override seluruh tim yang di-commit | menengah |
| `Taskfile.local.yml` / `.yaml` | pribadi, biasanya diabaikan oleh Git | tertinggi |

**Semantik penggabungan (lokal diprioritaskan):**

- Task dengan **nama yang sama** menimpa task dasar. Field daftar (`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`, `aliases`) **menggantikan** field dasar ketika disediakan oleh override; field yang tidak disertakan oleh override tetap diambil dari definisi dasar.
- `env` dan `vars` **digabungkan per kunci**, dan override diprioritaskan jika terjadi benturan.
- Task yang **hanya** ada dalam berkas override akan **ditambahkan**.

Misalnya, jika `Taskfile.yml` yang di-commit melakukan build dengan flag pengembangan, tetapi mesin Anda seharusnya selalu melakukan build produksi:

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

Sekarang `build` menjalankan perintah produksi Anda dan `smoke` tersedia, tanpa perubahan pada Taskfile yang di-commit.

@note{type="note" title="Model kepercayaan"}
Berkas override ditemukan dan diterapkan secara otomatis tanpa meminta konfirmasi. Hal ini tidak memberikan kemampuan baru—Taskfile sudah dapat menjalankan perintah shell arbitrer, sehingga override tidak dapat melakukan apa pun yang tidak dapat dilakukan dengan mengedit `Taskfile.yml`. `Taskfile.override.*` yang di-commit akan terlihat dalam diff PR; `Taskfile.local.*` dibuat di mesin Anda sendiri. Override yang tidak valid akan membatalkan proses eksekusi, bukan dilewati secara diam-diam. Tetapkan `WAILS_NO_OVERRIDES=true` untuk sepenuhnya melewati penemuan override dalam build CI yang deterministik.

@end

## Fallback otomatis

Jika Wake menemukan fitur Taskfile yang belum diimplementasikannya, Wake menyerahkan seluruh proses eksekusi kepada runtime Task tertanam. Hal-hal berikut saat ini memicu fallback tersebut:

- `dotenv` pada tingkat Taskfile
- mode `output` selain `interleaved`
- blok `requires`
- `interval` (pada tingkat Taskfile atau task)
- mode `run` selain `always`
- `short` dalam sebuah task
- `defer` dalam sebuah tugas

## Variabel lingkungan

| Variabel | Efek |
| --- | --- |
| `WAILS_USE_WAKE` | `true` mengaktifkan Wake untuk verba `wails3` yang dapat dirutekan; yang lainnya menggunakan runtime Task |
| `WAILS_NO_OVERRIDES` | `true` melewati penemuan `Taskfile.local.*` / `.override.*` (build deterministik) |
| `WAKE_VERBOSE` | Streaming stdout/stderr subproses secara langsung, alih-alih merekamnya hanya untuk ditampilkan saat terjadi kegagalan |
| `WAKE_SILENT` | Menyembunyikan seluruh output tugas |
| `WAKE_SERIAL` | `true` menonaktifkan fan-out paralel `deps:` (paralel digunakan secara default) |
| `WAKE_FORCE` | `true` melewati semua cache untuk build ulang yang benar-benar bersih |
| `WAKE_DEBUG` | Mencatat detail internal resolver ke log (DAG, dependensi, referensi variabel, perutean eksekusi) |
| `WAKE_NOTICE` | `off` untuk menyembunyikan pemberitahuan "wake (experimental)" pada setiap eksekusi |

Cache build berada di `.wake/cache.json` (Task menggunakan `.task/`).

## Umpan balik

Wake adalah sebuah eksperimen, dan umpan balik Anda menentukan arah pengembangannya. Jika Anda mencobanya, kami ingin mengetahui apakah Wake lebih cepat, lebih jelas, dan apakah ada sesuatu yang bermasalah — laporan yang paling berguna menyebutkan apa yang Anda jalankan, hasil yang Anda harapkan, dan apa yang sebenarnya terjadi. Sampaikan kepada kami dalam [diskusi umpan balik Wake](https://github.com/wailsapp/wails/discussions/5679).
