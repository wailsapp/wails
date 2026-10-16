---
title: "Instalasi"
description: "Instal Wails dan siapkan untuk membangun aplikasi"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## Instalasi Cepat (5 Menit)

@note{type="tip" title="Ringkasnya - Pengembang Berpengalaman"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

Atau verifikasi secara manual dengan `wails3 doctor`. [Lanjut ke Aplikasi Pertama →](/quick-start/first-app/)

@end

## Instalasi Langkah demi Langkah

@steps
### Instal Go (Wajib)
Wails memerlukan Go 1.25 atau yang lebih baru.

@tabs{sync-key="os"}
[Windows]
Unduh penginstal Windows dari **[go.dev/dl](https://go.dev/dl/)**, lalu jalankan.

**Verifikasi instalasi:**

```powershell
go version  # Should show 1.25 or later
```

**Periksa PATH:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

Jika kosong, tambahkan `C:\Users\YourName\go\bin` ke PATH Anda.

[macOS]
**Opsi 1: Penginstal Resmi**

Unduh penginstal macOS (file .pkg) dari **[go.dev/dl](https://go.dev/dl/)**, lalu jalankan.

**Opsi 2: Homebrew**

```bash
brew install go
```

**Verifikasi instalasi:**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

Jika `~/go/bin` tidak ada di PATH, tambahkan ke `~/.zshrc` atau `~/.bash_profile`:

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**Opsi 1: Tarball Resmi**

Unduh tarball Linux dari **[go.dev/dl](https://go.dev/dl/)**, lalu:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**Opsi 2: Manajer Paket**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Tambahkan ke PATH** (tambahkan ke `~/.bashrc` atau `~/.zshrc`):

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**Verifikasi:**

```bash
go version
echo $PATH | grep go/bin
```

@end

### Instal Dependensi Platform
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime** (biasanya sudah terinstal)

Windows 10/11 menyertakan WebView2 secara default. Jika tidak tersedia:

- Unduh dari [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Atau jalankan `wails3 doctor` nanti—perintah tersebut akan memandu Anda

**Selesai!** Tidak diperlukan dependensi lain.

@note{type="tip" title="Kiat Performa untuk Windows 11"}
Pertimbangkan untuk menggunakan [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) sebagai tempat menyimpan proyek Anda. Dev Drive dioptimalkan untuk beban kerja pengembangan dan dapat mempersingkat waktu build serta meningkatkan kecepatan akses disk secara signifikan hingga 30%.

@end

[macOS]
**Xcode Command Line Tools** (wajib)

```bash
xcode-select --install
```

Klik "Instal" pada dialog yang muncul.

**Verifikasi:**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**Selesai!** macOS menyertakan WebKit secara default.

[Linux]
**Alat build dan WebKit**

@note{type="caution" title="Versi minimum distribusi"}
Secara default, Wails v3 memerlukan **WebKitGTK 6.0**. Distribusi yang hanya menyediakan WebKit2GTK 4.1 — Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x — harus di-build dengan opsi keikutsertaan `-tags gtk3` lama. Rilis yang lebih lama dan hanya menyediakan WebKit2GTK 4.0 (Ubuntu 20.04, Debian 11, RHEL 8) tidak didukung.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
Memerlukan Ubuntu 24.04+ atau Debian 13+ untuk stack GTK4 default.

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
Tambahkan ke `shell.nix` atau `devShell` Anda:

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[Lainnya]
Jalankan `wails3 doctor` setelah menginstal Wails—perintah tersebut akan menampilkan paket yang tepat untuk distribusi Anda.

@end

@note{type="info" title="Stack GTK3 lama"}
Jika distribusi target Anda belum menyediakan WebKitGTK 6.0 (misalnya Ubuntu 22.04 LTS, Debian 12), instal pustaka pengembangan GTK3 + WebKit2GTK 4.1 sebagai gantinya (`libgtk-3-dev libwebkit2gtk-4.1-dev` di Debian/Ubuntu; paket setara di distribusi lain), lalu lakukan build dengan `wails3 build -tags gtk3`. Jalur lama didukung hingga lini v3.0.x dan akan dihapus pada v3.1. Lihat [Pemaketan Linux - Dukungan GTK3 Lama](/guides/build/linux/#legacy-gtk3-support) untuk detailnya.

@end

@end

### Instal CLI Wails
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Langkah ini menginstal perintah `wails3` ke `~/go/bin` (atau `%USERPROFILE%\go\bin` di Windows).

### Jalankan Wisaya Penyiapan (Disarankan)
```bash
wails3 setup
```

Wisaya penyiapan akan memeriksa dependensi Anda, membantu menginstal dependensi yang belum tersedia, dan mengonfigurasi default proyek.

@note{type="caution" title="Eksperimental"}
Wisaya penyiapan masih baru dan terutama telah diuji di Linux. Jika mengalami masalah, [laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) dan gunakan `wails3 doctor` sebagai gantinya.

@end

### Verifikasi Instalasi
```bash
wails3 doctor
```

**Keluaran yang diharapkan (atau serupa):**

```
Wails (v3.0.0-dev)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | Your installed version |
| Go Version  | go1.25.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="Jika perintah `wails3` tidak ditemukan"}
`~/go/bin` Anda tidak ada di PATH. Lihat langkah 1 di atas untuk memperbaikinya, lalu mulai ulang terminal Anda.

@end

### Instal npm (Opsional tetapi Disarankan)
Sebagian besar templat Wails menggunakan npm untuk alat frontend.

@tabs{sync-key="os"}
[Windows]
Unduh dari [nodejs.org](https://nodejs.org/), lalu jalankan penginstal.

**Verifikasi:**

```powershell
npm --version
```

[macOS]
**Opsi 1: Penginstal Resmi** Unduh dari [nodejs.org](https://nodejs.org/)

**Opsi 2: Homebrew**

```bash
brew install node
```

**Verifikasi:**

```bash
npm --version
```

[Linux]
**Opsi 1: NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**Opsi 2: Manajer Paket**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**Verifikasi:**

```bash
npm --version
```

@end

@note{type="tip" title="Manajer Paket Alternatif"}
Lebih memilih `pnpm`, `yarn`, atau `bun`? Tidak masalah! Cukup perbarui `Taskfile.yml` dalam proyek agar menggunakan alat pilihan Anda.

@end

@end

## Pemecahan Masalah

### Perintah `wails3` tidak ditemukan

**Penyebab:** `~/go/bin` (atau `%USERPROFILE%\go\bin`) tidak tercantum dalam PATH Anda.

**Solusi:**

@tabs{sync-key="os"}
[Windows]
1. Buka "Variabel Lingkungan" (cari di menu Mulai)
2. Di bagian "Variabel pengguna", cari `Path`
3. Klik "Edit" → "Baru"
4. Tambahkan: `C:\Users\YourName\go\bin` (ganti `YourName`)
5. Klik "OK" pada semua dialog
6. **Mulai ulang terminal Anda**

**Verifikasi:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
Tambahkan ke `~/.zshrc` (macOS) atau `~/.bashrc` (Linux):

```bash
export PATH=$PATH:~/go/bin
```

Muat ulang:

```bash
source ~/.zshrc  # or ~/.bashrc
```

**Verifikasi:**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` melaporkan dependensi yang belum tersedia

**Linux:** Output memberi tahu Anda paket mana saja yang harus diinstal. Contoh:

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows:** Jika WebView2 belum tersedia:

- Unduh dari [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Atau, WebView2 akan diinstal secara otomatis saat Anda menjalankan aplikasi pertama

**macOS:** Jika alat Xcode belum tersedia:

```bash
xcode-select --install
```

---

#### Versi Go terlalu lama

Wails v3 memerlukan Go 1.25+. Jika Anda menggunakan versi yang lebih lama:

@tabs{sync-key="os"}
[Windows/macOS]
Unduh versi terbaru dari [go.dev/dl](https://go.dev/dl/), lalu instal ulang.

[Linux]
Unduh tarball terbaru dari [go.dev/dl](https://go.dev/dl/), lalu:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## Versi Pengembangan (Paling Mutakhir)

Ingin menggunakan kode paling mutakhir dari cabang pengembangan utama? Dengan versi ini, Anda dapat mengakses fitur dan perbaikan baru sebelum dirilis, tetapi ada risiko bug dan perubahan yang merusak kompatibilitas. Hanya disarankan bagi kontributor atau pengguna yang perlu menguji fitur mendatang.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="Versi Pengembangan"}
- Mungkin mengandung bug atau perubahan yang merusak kompatibilitas
- Proyek yang dibuat akan menggunakan direktif `replace` untuk merujuk ke Wails lokal
- Hanya disarankan bagi kontributor atau untuk menguji fitur baru

@end

## Langkah Berikutnya

**Instalasi Selesai!** Sistem Anda siap untuk pengembangan dengan Wails.

@cards{cols="1"}
🚀 Buat Aplikasi Pertama Anda
Buat aplikasi yang berfungsi dalam 10 menit.

[Tutorial Aplikasi Pertama →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 Jelajahi Templat
Lihat apa saja yang langsung tersedia.

```bash
wails3 init -l  # List templates
```

@end

---

**Mengalami masalah?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau [buat issue](https://github.com/wailsapp/wails/issues).
