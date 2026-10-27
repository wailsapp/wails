---
title: "Pemaketan Linux"
description: "Paketkan aplikasi Wails Anda untuk distribusi Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Format Paket

Paketkan aplikasi Anda untuk distribusi Linux:

```bash
wails3 package GOOS=linux
```

Perintah ini membuat beberapa format di direktori `bin/`:

- **AppImage**: Portabel, berjalan di semua distribusi Linux
- **DEB**: Untuk Debian, Ubuntu, dan turunannya
- **RPM**: Untuk Fedora, RHEL, dan turunannya
- **Arch**: Untuk Arch Linux dan turunannya

### Format Individual

Buat format tertentu:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Menyesuaikan Paket

### Entri Desktop

File `.desktop` mengatur tampilan aplikasi Anda dalam menu aplikasi. File ini dibuat dari nilai-nilai dalam `build/linux/Taskfile.yml`:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Metadata Paket

Edit `build/linux/nfpm/nfpm.yaml` untuk menyesuaikan paket DEB dan RPM:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

Konfigurasi AppImage berada di `build/linux/appimage/`. Ikon aplikasi berasal dari `build/appicon.png`.

AppImage menyertakan WebKit beserta proses pembantunya dan menambal pustaka yang disertakan agar proses tersebut dimulai dari dalam AppImage. `AppRun` miliknya memulai aplikasi dari `usr/lib` di dalam AppImage; jalur yang telah ditambal dapat ditemukan baik di luar maupun di dalam sandbox bubblewrap WebKit, sehingga sandbox tetap aktif. Sediakan `AppRun` Anda sendiri hanya jika direktori kerja tersebut tetap dipertahankan.

## Menandatangani Paket

Tandatangani paket DEB dan RPM dengan kunci PGP:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Konfigurasikan penandatanganan di `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Simpan kata sandi kunci Anda:

```bash
wails3 signing
```

Lihat [Menandatangani Aplikasi](/guides/build/signing/) untuk detail selengkapnya.

## Membangun untuk ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
Build ARM64 dari host x86_64 menggunakan Docker untuk kompilasi silang CGO.

@end

## Dukungan GTK3 Lama

Secara default, Wails v3 dibangun pada **GTK4 dengan WebKitGTK 6.0**. Jalur GTK3 / WebKit2GTK 4.1 lama masih tersedia untuk distribusi yang belum menyediakan WebKitGTK 6.0 (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x). Jalur lama ini harus diaktifkan secara eksplisit melalui tag build dan dijadwalkan untuk dihapus dalam v3.1.

@note{type="caution" title="Jalur lama"}
Jalur GTK3 / WebKit2GTK 4.1 didukung hingga lini v3.0.x. Rencanakan migrasi ke GTK4 sesuai dengan ketersediaan GTK4 / WebKitGTK 6.0 pada distribusi target Anda—`-tags gtk3` akan dihapus dalam v3.1.

@end

### Dependensi

Instal pustaka pengembangan GTK3 dan WebKit2GTK 4.1:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Paket pkg-config yang diperlukan adalah `gtk+-3.0` dan `webkit2gtk-4.1`.

### Membangun dengan GTK3

Gunakan flag `-tags gtk3`:

```bash
wails3 build -tags gtk3
```

Atau langsung dengan Go:

```bash
go build -tags gtk3 -o myapp .
```

### Perbedaan yang Diketahui dari GTK4

- **Dialog file**: GTK4 menggunakan `xdg-desktop-portal` untuk dialog file (secara default), sehingga beberapa opsi dialog (seperti direktori default dan tampilan filter khusus) berperilaku berbeda dari GTK3. Lihat [Referensi Dialog - Perilaku Dialog Linux](/reference/dialogs/#linux-dialog-behavior) untuk detail selengkapnya.
- **Gaya menu**: GTK4 mendukung opsi `LinuxMenuStylePrimaryMenu` yang menampilkan tombol hamburger (☰) di bilah header sesuai dengan GNOME HIG. Opsi ini tidak berpengaruh pada build `-tags gtk3`. Lihat [API Window - Linux MenuStyle](/reference/window/#linux).
- **Penskalaan DPI**: GTK4 menggunakan `gdk_monitor_get_scale` (GTK 4.14+) untuk mendukung penskalaan fraksional.

### Memeriksa Build Anda

Jalankan `wails3 doctor` untuk memverifikasi penyiapan Anda. Tanpa flag, perintah ini memeriksa GTK4 / WebKitGTK 6.0 (default). Paket GTK3 / WebKit2GTK 4.1 lama dicantumkan sebagai opsional.

## Pemecahan Masalah

### AppImage tidak dapat dijalankan

Jadikan file dapat dieksekusi:

```bash
chmod +x MyApp-x86_64.AppImage
```

### Dependensi tidak tersedia

Sesuaikan pustaka dengan varian build biner. Untuk build GTK4/WebKitGTK 6.0 default:

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

Untuk biner yang secara eksplisit dibangun dengan `-tags gtk3`, instal pustaka GTK3/WebKit2GTK 4.1 lama sebagai gantinya:

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Ketersediaan dan nama paket bergantung pada rilis distribusi. Gunakan rilis yang didukung dan menyediakan pustaka yang diperlukan, atau gunakan jalur build silang/container yang didokumentasikan; menginstal pustaka lama tidak memenuhi dependensi build default. Portabilitas AppImage bergantung pada komponen yang benar-benar disertakan oleh proses pemaketannya dan pustaka host yang masih diperlukannya.

### Compiler C tidak ditemukan

Sistem build memerlukan GCC atau Clang untuk CGO:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Atau, jalankan `wails3 task setup:docker` dan sistem build akan menggunakan Docker secara otomatis.

### Jendela kosong atau putih pada GPU NVIDIA

Di Linux dengan driver proprietari NVIDIA, aplikasi Wails mungkin menampilkan jendela kosong atau putih saat dimulai. Hal ini disebabkan oleh bug WebKitGTK yang membuat perender DMA-BUF gagal dengan `gbm_bo_map()` saat digunakan bersama driver proprietari NVIDIA (memengaruhi X11 dan Wayland, driver versi 377–580+, GPU dari seri 10 dan GT 710 yang lebih lama).

**Wails menerapkan `WEBKIT_DISABLE_DMABUF_RENDERER=1` secara otomatis** ketika mendeteksi modul kernel NVIDIA (`/sys/module/nvidia`), sehingga sebagian besar pengguna tidak perlu melakukan apa pun.

Jika jendela masih kosong (misalnya dalam container yang tidak dapat melihat jalur modul), atur variabel lingkungan secara manual sebelum menjalankan aplikasi Anda:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Bug upstream terkait: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Kompatibilitas strip AppImage

Pada distribusi Linux modern (Arch Linux, Fedora 39+, Ubuntu 24.04+), pustaka sistem dikompilasi dengan bagian ELF `.relr.dyn` agar relokasi lebih efisien. Alat `linuxdeploy` yang digunakan untuk membuat AppImage menyertakan biner `strip` lama yang tidak dapat memproses bagian modern ini.

Wails mendeteksi situasi ini secara otomatis dengan memeriksa pustaka GTK sistem sebelum membangun AppImage. Jika terdeteksi, stripping dinonaktifkan (`NO_STRIP=1`) untuk memastikan kompatibilitas.

**Artinya:**

- AppImage akan sedikit lebih besar (~20-40%) pada sistem yang terdampak
- Fungsionalitas aplikasi tidak terpengaruh
- Hal ini ditangani secara otomatis—tidak diperlukan tindakan apa pun

Jika memerlukan AppImage yang lebih kecil pada sistem modern, Anda dapat menginstal biner `strip` yang lebih baru dan mengonfigurasi `linuxdeploy` untuk menggunakannya sebagai pengganti versi bawaannya.
