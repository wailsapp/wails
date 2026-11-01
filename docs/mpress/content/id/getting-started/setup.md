---
title: "Penyiapan"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Eksperimental"}
Wizard penyiapan masih baru dan terutama telah diuji di Linux. Jika mengalami masalah, silakan [laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) dan ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies) sebagai gantinya.

@end

## Mulai Cepat

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

Wizard terbuka di browser dan memandu Anda memeriksa dependensi, menetapkan nilai default proyek, serta menyiapkan build lintas platform secara opsional. Petunjuk di atas menggunakan alur penyiapan default; alur eksperimental untuk proyek Wake dijelaskan di bawah ini.

Setelah itu, Anda siap membuat proyek pertama:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fungsinya

- **Memeriksa dependensi** - Memverifikasi Go, npm, dan alat platform
- **Mengonfigurasi nilai default** - Informasi pembuat, prefiks ID bundle, dan templat pilihan
- **Build lintas platform** - Penyiapan Docker opsional untuk melakukan build dari host mana pun
- **Penandatanganan kode** - Penyiapan opsional untuk macOS, Windows, dan Linux

Konfigurasi disimpan ke `$XDG_CONFIG_HOME/wails/defaults.yaml`, dengan `~/.config/wails/defaults.yaml` sebagai alternatif, dan digunakan oleh `wails3 init`. Konfigurasi ini terpisah dari `build/config.yml` milik proyek.

## Mengonfigurasi proyek Wake

Untuk membuat aplikasi pertama secara lengkap, ikuti [tutorial memulai Wake](/experimental/wake/).
Tutorial ini memandu Anda melalui pembuatan proyek dan setiap layar penyiapan, menjelaskan hal-hal yang dapat diharapkan dari fitur eksperimental tersebut, serta membahas migrasi dari build yang sudah ada.

Jalankan penyiapan dari direktori proyek dengan [eksperimen Wake diaktifkan](/experimental/yaml-builds/#enable-yaml-builds):

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

Di PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

Tindakan ini membuka alur yang berfokus pada proyek:

1. **Proyek** — masukkan metadata proyek jika tidak ada manifes native.
2. **Aplikasi** — pilih GUI, CLI, atau server; backend build; atau berkas executable yang sudah di-build. Go murni merupakan pilihan eksperimental yang harus dipilih secara eksplisit.
3. **Target** — pilih sistem operasi dan arsitektur. iOS dan Android tetap bersifat eksperimental; konfigurasikan pengaturan penandatanganannya dalam YAML.
4. **Paket** — pilih format paket dan apakah paket tersebut akan ditandatangani.
5. **Alat** — tinjau persyaratan untuk jalur build dan pemaketan yang dipilih.
6. **Penandatanganan** — jika dipilih, konfigurasikan penandatanganan menggunakan kontrol penandatanganan yang sudah ada.
7. **Simpan** — tulis konfigurasi proyek dan tampilkan perintah yang harus dijalankan.

Persyaratan alat bergantung pada pilihan Anda. Kompilasi desktop dengan Go murni tidak memerlukan Docker atau compiler C. Build lintas target dengan CGo memerlukan compiler dan pustaka target yang kompatibel, atau toolchain container yang didukung. NSIS tetap memerlukan `makensis`; AppImage GUI yang di-build lintas host memerlukan pohon runtime Linux yang disediakan. Identitas penandatanganan yang diekspor dapat menggunakan penanda tangan Go milik Wails, sedangkan identitas dari keychain atau penyimpanan sertifikat native tetap menggunakan alat platform. Lihat [persyaratan pemaketan dan penandatanganan](/experimental/yaml-builds/#package-and-sign-for-distribution).

Berkas executable yang sudah di-build melewati kompilasi Go, instalasi dan build frontend, binding, serta UPX. Berkas tersebut tetap memerlukan alat pemaketan dan penandatanganan yang dipilih serta berkas biner input yang valid. Lihat [jalur input hasil build sebelumnya](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Pilihan ini berlaku untuk perintah yang ditampilkan setelah penyimpanan; pilihan ini bukan pengaturan build YAML permanen.

Simpan menulis `wails.yaml`, memperbarui jenis aplikasi, pengaktifan frontend, dan pengaturan backend yang dipilih, serta mengganti pilihan target/paket/penandatanganan khusus `profiles.setup`. Profil lain, perintah frontend, executable pendamping, hook, dan pengaturan proyek lain yang tidak terkait tetap dipertahankan. Pengaturan backend berlaku untuk seluruh proyek, sehingga profil lain yang mewarisinya juga menggunakan backend yang dipilih. Metadata proyek baru tetap berada di wizard hingga Anda memilih Simpan. Alat yang belum tersedia akan dilaporkan agar Anda dapat menginstalnya sebelum melakukan build; menyimpan konfigurasi yang valid tidak menginstal alat atau menjalankan build.

Build yang ditandatangani menggunakan kembali pengaturan penandatanganan yang tidak tersedia dari nilai default tingkat pengguna, termasuk berkas sertifikat di luar proyek. Pengaturan proyek yang ditentukan secara eksplisit diprioritaskan. Penyiapan tidak menyalin kunci privat, kata sandi, atau jalur kredensial global ke `wails.yaml`. Untuk kunci Windows dan Linux yang dikonfigurasi secara global, berikan kata sandi melalui penyimpanan kredensial yang ada atau `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Untuk identitas Apple terenkripsi, tetapkan `darwin.signing.credential` ke nama variabel lingkungan yang berisi kata sandinya.

Profil yang disimpan dijalankan dengan `wails3 build setup`; opsi berkas yang sudah di-build menambahkan `--skip-build`.

Pastikan Wake tetap aktif saat menjalankan perintah yang ditampilkan. `wails.yaml` yang aktif tidak pernah beralih ke Taskfile sebagai alternatif ketika flag tidak ditetapkan. Penyiapan tidak mengonversi Taskfile yang sudah ada atau metadata lama: gunakan [alur kerja migrasi yang telah ditinjau](/experimental/yaml-builds/#migrate-an-existing-project) sebelum mengonfigurasi proyek tersebut dengan Wake.

## Subperintah

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Mengalami Masalah?

1. Jalankan `wails3 doctor` untuk mendiagnosis masalah
2. Ikuti [langkah-langkah instalasi manual](/getting-started/installation/#platform-specific-dependencies)
3. [Laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) beserta output `wails3 doctor` Anda
