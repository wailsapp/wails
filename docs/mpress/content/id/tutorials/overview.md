---
title: "Tutorial"
description: "Pelajari Wails dengan membangun aplikasi"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Tutorial langkah demi langkah yang mengajarkan konsep Wails melalui pembangunan aplikasi lengkap. Setiap tutorial menyertakan kode yang berfungsi, penjelasan, dan pola praktis.

@note{type="tip" title="Baru mengenal Go?"}
Selesaikan [Tur Go](https://go.dev/tour/) sebelum memulai tutorial.

@end

## Layanan Kode QR

![Contoh Kode QR](/assets/qr1.png)

Pelajari dasar-dasar layanan Wails dengan membangun generator kode QR. Tutorial ini memperkenalkan konsep inti untuk menata logika aplikasi Anda menjadi layanan yang dapat digunakan kembali.

**Yang akan Anda pelajari:**

- Cara membuat dan menyusun layanan Wails
- Mengelola dependensi Go eksternal
- Mengikat metode Go ke frontend Anda
- Meneruskan data antara Go dan JavaScript
- Menata kode agar mudah dipelihara

**Paling cocok untuk:** Pengguna Wails pemula yang ingin memahami arsitektur layanan

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Mulai</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Daftar TODO

![Aplikasi Daftar TODO](/assets/todo-app.png)

Bangun aplikasi daftar TODO yang lengkap dengan antarmuka modern dan menarik. Tutorial praktik ini mengajarkan pola-pola inti Wails melalui aplikasi praktis yang sesuai dengan penggunaan nyata menggunakan JavaScript murni.

**Yang akan Anda pelajari:**

- Arsitektur berbasis layanan dengan pengelolaan state yang aman untuk thread
- Operasi CRUD (Buat, Baca, Perbarui, Hapus)
- Binding yang aman secara tipe antara Go dan JavaScript
- Membangun UI modern tanpa kerumitan framework
- Pola penanganan error dan validasi yang tepat

**Waktu penyelesaian:** ~20 menit

**Paling cocok untuk:** Aplikasi Wails lengkap pertama Anda—ideal untuk memahami dasar-dasarnya sebelum menambahkan kerumitan framework

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Mulai</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Catatan

![Aplikasi Catatan](/assets/notes-app.png)

Bangun aplikasi bergaya Apple Notes dengan dialog file native dan fungsi penyimpanan otomatis. Tutorial ini mendemonstrasikan fitur khusus desktop seperti operasi file, dialog native, dan pola UI profesional.

**Yang akan Anda pelajari:**

- Dialog file native (Simpan, Buka, Info)
- Persistensi data berbasis JSON
- Pola penyimpanan otomatis dengan debounce
- Tata letak desktop dua kolom yang profesional
- Menggunakan operasi sistem file di Go

**Waktu penyelesaian:** ~30 menit

**Paling cocok untuk:** Mempelajari fitur khusus desktop seperti operasi file dan dialog native OS

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Mulai</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Aplikasi Wails yang Memperbarui Diri Sendiri

Tambahkan pembaruan mandiri dalam aplikasi ke aplikasi Wails mulai dari `wails3 init` baru hingga verifikasi rilis bertanda tangan dan penggantian biner dalam mode helper. Menggunakan GitHub Releases sebagai sumber pembaruan.

**Yang akan Anda pelajari:**

- Cara `app.Updater` diintegrasikan ke aplikasi Wails
- Mengonfigurasi penyedia GitHub Releases
- Menerbitkan rilis dengan `SHA256SUMS` untuk verifikasi digest
- Menambahkan penandatanganan Ed25519 agar tahan terhadap manipulasi
- Menyesuaikan jendela default melalui CSS, HTML khusus, atau BYO
- Pemeriksaan latar belakang berkala dengan `CheckInterval`

**Waktu penyelesaian:** ~25 menit

**Paling cocok untuk:** Mendistribusikan aplikasi desktop yang dapat diperbarui—mencakup seluruh pipeline rilis, bukan hanya API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Mulai</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
