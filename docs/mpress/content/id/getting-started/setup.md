---
title: "Penyiapan"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Eksperimental"}
Wizard penyiapan ini masih baru dan terutama telah diuji di Linux. Jika Anda mengalami masalah, silakan [laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904), lalu ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies) sebagai gantinya.

@end

## Mulai Cepat

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

Wizard akan terbuka di browser dan memandu Anda memeriksa dependensi, menentukan nilai default proyek, serta menyiapkan build lintas platform secara opsional.

Setelah itu, Anda siap membuat proyek pertama:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fungsinya

- **Memeriksa dependensi** - Memverifikasi Go, npm, dan alat platform
- **Mengonfigurasi nilai default** - Informasi pembuat, prefiks ID bundel, dan templat pilihan
- **Build lintas platform** - Penyiapan Docker opsional untuk menjalankan build dari host apa pun
- **Penandatanganan kode** - Penyiapan opsional untuk macOS, Windows, dan Linux

Konfigurasi disimpan ke `~/.config/wails/config.yaml` dan digunakan oleh `wails3 init`.

## Subperintah

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## Mengalami Masalah?

1. Jalankan `wails3 doctor` untuk mendiagnosis masalah
2. Ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies)
3. [Laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) beserta output `wails3 doctor` Anda
