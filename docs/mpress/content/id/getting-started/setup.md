---
title: "Penyiapan"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Eksperimental"}
Wizard penyiapan ini masih baru dan terutama telah diuji di Linux. Jika Anda mengalami masalah, harap [laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) dan ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies) sebagai gantinya.

@end

## Mulai Cepat

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

Wizard akan terbuka di browser dan memandu Anda melalui pemeriksaan dependensi, pengaturan default proyek global, alat lintas platform opsional, serta kredensial penandatanganan. Preferensi ini dapat digunakan kembali di berbagai proyek.

Setelah itu, Anda siap membuat proyek pertama:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fungsinya

- **Memeriksa dependensi** - Memverifikasi Go, npm, dan alat platform
- **Mengonfigurasi default** - Informasi pembuat, prefiks ID bundel, dan templat pilihan
- **Build lintas platform** - Penyiapan Docker opsional untuk melakukan build dari host apa pun
- **Penandatanganan kode** - Penyiapan opsional untuk macOS, Windows, dan Linux

Konfigurasi disimpan ke `$XDG_CONFIG_HOME/wails/defaults.yaml`, dengan `~/.config/wails/defaults.yaml` sebagai fallback, dan digunakan oleh `wails3 init`. Konfigurasi ini terpisah dari `build/config.yml` milik proyek.

## Membuat proyek Wake

Jika eksperimen Wake diaktifkan, `wails3 init -ui` mengumpulkan pilihan build khusus proyek Anda dan menyimpannya di `wails.yaml` setelah pembuatan berhasil. Ikuti [tutorial memulai Wake](/experimental/wake/) untuk memilih backend build, target, dan format paket saat membuat proyek desktop. `wails3 setup` tetap berperan menyimpan preferensi global.

Proyek Wake yang sudah ada dikonfigurasi di `wails.yaml`. Gunakan [panduan build YAML](/experimental/yaml-builds/) untuk mengubah profil dan pengaturan pemaketannya. Proyek Taskfile yang sudah ada memerlukan [migrasi satu arah yang ditinjau](/experimental/wake/#migrate-an-existing-project) sebelum dapat menggunakan Wake; penyiapan tidak menerjemahkan atau menjalankan Taskfile proyek tersebut.

## Menggunakan Kembali Preferensi Penandatanganan

Build Wake yang ditandatangani menggunakan kembali pengaturan penandatanganan yang belum ditetapkan dari default tingkat pengguna Anda, termasuk file sertifikat di luar proyek. Pengaturan proyek yang ditetapkan secara eksplisit akan diprioritaskan. Kunci privat, kata sandi, dan jalur kredensial global tidak disalin ke `wails.yaml`.

Untuk kunci Windows dan Linux yang dikonfigurasi secara global, berikan kata sandi melalui penyimpanan kredensial yang sudah ada atau `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Untuk identitas Apple terenkripsi, atur `darwin.signing.credential` ke nama variabel lingkungan yang berisi kata sandinya.

## Subperintah

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Mengalami Masalah?

1. Jalankan `wails3 doctor` untuk mendiagnosis masalah
2. Ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies)
3. [Laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) beserta output `wails3 doctor` Anda
