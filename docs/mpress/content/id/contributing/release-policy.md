---
title: "Kebijakan Rilis"
description: "Cara Wails v3 menetapkan versi rilis nightly dan weekly, cara pengelola membuat rilis tersebut, serta alasan pengguna harus menyematkan versi tertentu."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## Ringkasan model dalam satu tabel

Wails v3 merilis berdasarkan **jadwal, bukan fitur**. Setelah versi dasar `v3.0.0`:

| Rilis | Tag | Frekuensi | Rilis GitHub | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| Nightly | `v3.<next>.0-nightly.<YYYYMMDD>` (tanggal UTC; `.N` ditambahkan jika dijalankan ulang pada hari yang sama) | otomatis, setiap hari, hanya jika `UNRELEASED_CHANGELOG.md` memiliki konten | prarilis, **bukan** Latest | versi yang sama tanpa `v`, dist-tag `nightly` |
| Weekly | `v3.<next>.0` | otomatis, Minggu 20:00 waktu Australia/Sydney **(hari belum dikonfirmasi)** | rilis normal, Latest | versi yang sama, dist-tag `latest` |
| Hotfix | `v3.<current>.<N>`, `N > 0` | jarang, manual | rilis normal, Terbaru | versi yang sama, dist-tag `latest` |
| Mayor | `v4.0.0` | hanya untuk perubahan yang merusak kompatibilitas | rilis normal | `latest` |

`<next>` adalah minor **yang telah dirilis** v3 tertinggi ditambah satu. Nightly adalah pratinjau **rilis mingguan berikutnya**. Karena itu, versi dasarnya adalah tag rilis tertinggi (prarilis diabaikan): setelah hotfix `v3.1.1`, nightly berikutnya tetap `v3.2.0-nightly.*`. Nightly harus menggunakan nomor versi yang lebih tinggi karena prarilis diurutkan sebelum rilisnya sendiri: `v3.0.0-nightly.<date>` akan diurutkan sebelum `v3.0.0`.

**Kenaikan versi minor bukan tonggak pencapaian fitur.** `v3.4.0` berarti "rilis weekly setelah `v3.3.0`". Untuk mengetahui perubahan yang terjadi, baca catatan rilis, bukan nomor versinya.

Kebijakan dipilih berdasarkan **bentuk** versi: akhiran `-nightly.` menunjukkan nightly; akhiran lainnya (`-alpha`, `-beta`, `-rc`) menunjukkan build lama sebelum GA; tanpa akhiran berarti rilis.

### Versi yang selalu selaras

CLI, modul Go, dan paket npm `@wailsio/runtime` menggunakan **versi yang sama**, yang diambil dari `v3/internal/version/version.txt` (`3.2.0-nightly.20261102` di npm untuk tag `v3.2.0-nightly.20261102`). Proyek yang dibuat dengan `wails3 init` menyematkan `@wailsio/runtime` ke versi CLI itu sendiri; build pengembangan CLI menggunakan `latest` sebagai nilai fallback.

## Contoh lengkap

| Tanggal | Peristiwa | Tag | `go install ...@latest` di-resolve menjadi |
| --- | --- | --- | --- |
| Min 25 Okt | GA, rilis manual | `v3.0.0` | `v3.0.0` |
| Sen 26–Sab 31 Okt | nightly | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| Min 1 Nov 20:00 waktu Sydney | weekly | `v3.1.0` | `v3.1.0` |
| Sen 2 Nov | nightly (pratinjau weekly berikutnya) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| Rab 4 Nov | hotfix, manual, dari cabang rilis | `v3.1.1` | `v3.1.1` |
| Min 8 Nov | weekly | `v3.2.0` | `v3.2.0` |

Weekly selalu memiliki urutan lebih tinggi daripada setiap nightly dengan nomor yang sama, dan nightly tidak pernah menjadi `@latest`.

## Yang dijalankan pengguna

| Tujuan | Perintah |
| --- | --- |
| Rilis terbaru (weekly atau hotfix) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` atau `wails3 update cli` |
| **Dapat direproduksi (disarankan)** | sematkan tag tertentu: `v3.1.0` di `go.mod`, CI, dan dokumentasi |
| Nightly tertentu | tag persisnya (`@v3.2.0-nightly.20261102`), atau `wails3 update cli -nightly` |
| Build terbaru dari jenis apa pun | `wails3 update cli -pre` |
| Paket npm runtime | `@wailsio/runtime` (`latest` = mingguan/hotfix, `nightly` = rilis nightly) |

Alasan mekanisme ini berfungsi: Go me-resolve `@latest` ke tag tertinggi **tanpa akhiran prarilis** dan hanya mempertimbangkan prarilis jika tidak ada rilis. Kotak centang "pre-release" di GitHub tidak terlihat oleh proksi Go. Karena itu, nightly menggunakan akhiran prarilis semver, bukan mengandalkan kotak centang tersebut. Selalu sematkan versi tertentu dalam produksi.

## Cabang

| Cabang | Tujuan |
| --- | --- |
| `master` | Pengembangan dan nightly. |
| `releases/**` | Cabang berumur pendek untuk rilis satu kali: GA `v3.0.0` dan hotfix. Memerlukan CI yang sama dengan `master`. |

Weekly dibuat dari `master`, tetapi **bukan dari commit terbarunya**: alur kerja memilih commit terbaru sejak weekly sebelumnya yang telah lulus semua pemeriksaan wajib (`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`), lalu menandai tepat commit tersebut dengan tag. Jika tidak ada commit yang berstatus hijau, tidak ada rilis. Nightly dibuat dari `master` dan hanya memberikan peringatan ketika HEAD tidak berstatus hijau.

## Cara membuat rilis

Hak merilis hanya diberikan kepada Lead proyek (lihat dokumen tata kelola). Kontributor tidak pernah membuat tag atau menerbitkan rilis; mereka menambahkan entri ke `v3/UNRELEASED_CHANGELOG.md`.

Semuanya dijalankan dari **Nightly + Weekly Release v3** (`.github/workflows/nightly-release-v3.yml`):

- **Nightly (otomatis)** pada pukul 15:00 UTC. Dilewati jika changelog yang belum dirilis kosong.
- **Mingguan (otomatis)** pada 20:00 waktu Sydney. Cron GitHub hanya mendukung UTC, sehingga dua cron (09:00 dan 10:00 UTC) dijadwalkan, dan job pengaman hanya dilanjutkan saat waktu di Sydney menunjukkan 20:00 (dengan memperhitungkan DST). Rilis mingguan harus lolos gerbang rilis, memilih commit terbaru yang berhasil melewati semua pemeriksaan, menaikkan versi minor, memberi tag pada commit tersebut tanpa mengubah changelog, lalu hanya memperbarui kolom versi pada `master` (`version.txt` dan versi paket runtime), sehingga pembaruan dependensi dan lockfile yang digabungkan setelah commit terpilih tetap dipertahankan. npm menerbitkan dari tag, bukan dari `master`.
- **Manual.** Jalankan alur kerja dengan `channel`, `dry_run` (aktif secara default), dan jika perlu `release_version`. Dry run menampilkan versi, jenis, flag GitHub, dist-tag npm, dan catatan tanpa mengubah apa pun di GitHub atau npm.
- **GA atau hotfix.** Tetapkan `release_version` (`v3.0.0`, `v3.1.1`). Hotfix untuk lini versi lama juga dapat menetapkan `ref: releases/...`; commit rilisnya di-cherry-pick kembali ke `master`.

Uji coba lokal: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`; `go run release.go --print-policy v3.2.0-nightly.20261102` menampilkan flag yang diturunkan.

### Mekanisme pengaman

- Versi otomatis harus benar-benar lebih tinggi (berdasarkan urutan semver, prarilis < rilis) daripada setiap tag yang sudah ada; jika tidak, alat akan menolaknya.
- Hanya satu tugas rilis yang dapat berjalan pada satu waktu (grup konkurensi bersama).
- Publikasi npm dipicu oleh **tag** rilis dan membangun commit yang diberi tag tersebut.

## Gerbang rilis

Proses **mingguan** (bukan proses nightly) menolak membuat tag selama salah satu hal berikut masih ada:

1. issue atau pull request terbuka yang diberi label `P0` atau `release-blocker`;
2. issue terbuka yang diberi label `regression` dan dibuat sejak rilis mingguan terakhir;
3. peringatan Dependabot terbuka dengan tingkat keparahan kritis atau tinggi, peringatan pemindaian kode terbuka dengan tingkat keparahan kritis atau tinggi, serta setiap peringatan pemindaian rahasia yang masih terbuka, **hanya untuk kode yang didistribusikan**. Peringatan yang hanya berada di bawah `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/`, atau direktori `examples/` mana pun akan diabaikan (`SHIPPED_ALERT_IGNORE` di `v3/tasks/releasegate/main.go`, **perlu ditinjau oleh pemilik**);
4. advisori keamanan repositori yang masih berstatus `triage`.

**Siapa yang dapat menetapkan pemblokir.** Sebuah label hanya diperhitungkan jika peristiwa `labeled` terbarunya (tanpa `unlabeled` setelahnya) dilakukan oleh anggota tim `wailsapp/maintainers` atau kolaborator dengan izin `maintain`/`admin`. Label yang sama jika diterapkan oleh pihak lain akan diabaikan dan dicatat dalam log, sehingga pihak luar tidak dapat memblokir atau membuka blokir rilis. Label `high-priority` yang diterapkan secara otomatis sengaja **tidak** digunakan. Akun bot seperti `taliesin-ai` lolos daftar izin jika memiliki `maintain`; hal ini memang disengaja, dan agen triase dapat menerapkan label pemblokir dengan bukti.

**Perilaku saat terjadi kegagalan.** Gerbang menerapkan prinsip **fail-closed**: jika API peringatan atau daftar izin tidak dapat dibaca, proses mingguan akan berhenti. Jika gagal, proses ini tidak membuat tag apa pun, menuliskan daftar tersebut (beserta tautannya) ke ringkasan tugas, dan memperbarui komentar pada satu-satunya issue "Status rilis mingguan" jika variabel repositori `WEEKLY_STATUS_ISSUE` telah ditetapkan.

**Penyiapan yang harus dilakukan pemilik di GitHub** (tidak ada yang dibuat secara otomatis): buat label `P0`, `release-blocker`, dan `regression`; buat issue "Status rilis mingguan" dan tetapkan `WEEKLY_STATUS_ISSUE` ke nomornya; berikan `WAILS_REPO_TOKEN` akses baca ke Issues, Pull requests, Administration (izin kolaborator), Dependabot alerts, Code scanning alerts, Secret scanning alerts, Repository security advisories, dan Members organisasi, serta akses tulis ke Issues untuk komentar status. `GITHUB_TOKEN` bawaan tidak dapat membacanya.

**Pengabaian.** `override_reason` (input workflow-dispatch) melewati gerbang. Alasan, pelaku, dan setiap pemblokir yang dilewati dituliskan ke ringkasan dan issue status. Hanya Lead rilis yang boleh menggunakannya. GitHub tidak dapat membatasi input hanya untuk satu orang, jadi jalankan tugas rilis dalam environment dengan peninjau wajib (disarankan), agar pengabaian memerlukan persetujuan eksplisit.

## Terkait

- [Pengujian & Integrasi Berkelanjutan](/contributing/testing-ci/)
- [Triase Pengelola](/contributing/maintainer-triage/)
