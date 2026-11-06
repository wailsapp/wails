---
title: "Aplikasi Pertama Anda"
description: "Buat aplikasi Wails yang berfungsi dalam 10 menit"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Kita akan membuat aplikasi sapaan sederhana yang memperagakan konsep inti Wails:

- Backend Go yang menangani logika
- Frontend yang memanggil fungsi Go
- Binding yang aman secara tipe
- Hot reload selama pengembangan

**Waktu penyelesaian:** 10 menit

## Buat Proyek Anda

@steps
### Buat proyek
```bash
wails3 init -n myapp -t vanilla-js
cd myapp
```

Pelajaran ini secara eksplisit memilih Vanilla + Vite dengan JavaScript agar berkas sumber di bawah sesuai dengan proyek yang dihasilkan. Templat dasar `vanilla` menggunakan TypeScript.

@note{type="tip" title="Templat Lain"}
Cobalah `-t react`, `-t vue`, atau `-t svelte` untuk framework pilihan Anda. Secara default, templat ini menggunakan TypeScript; untuk JavaScript biasa, gunakan `-t vanilla-js` atau `-t react-js`. Jalankan `wails3 init -l` untuk melihat semua templat yang tersedia, atau [gunakan framework frontend Anda sendiri](/guides/dev/frontend-frameworks/).

@end

### Pahami struktur proyek
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.js      # Frontend JavaScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # Editor and JavaScript tooling configuration
│   └── vite.config.js   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### Jalankan aplikasi
```bash
wails3 dev
```

@note{type="info" title="Menjalankan untuk Pertama Kalinya"}
Proses pertama kali mungkin memerlukan waktu lebih lama dari yang diperkirakan karena aplikasi memasang dependensi frontend, menghasilkan binding, dan sebagainya. Proses berikutnya jauh lebih cepat.

@end

Aplikasi akan terbuka dan menampilkan antarmuka sapaan. Masukkan nama Anda, lalu klik "Sapa"—backend Go akan memproses masukan Anda dan mengembalikan sapaan.

@end

@note{type="tip" title="Kiat Performa untuk Pengguna Windows 11"}
Pertimbangkan untuk menggunakan [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) guna menyimpan proyek Anda. Dev Drive dioptimalkan untuk beban kerja pengembang dan dapat meningkatkan waktu build serta kecepatan akses disk secara signifikan, hingga 30% dibandingkan drive NTFS biasa.

@end

## Cara Kerjanya

Mari pahami kode yang membuat aplikasi ini berfungsi.

### Backend Go

Buka `greetservice.go`:

```go {title="greetservice.go"}
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
```

**Konsep utama:**

1. **Layanan** - Struct Go dengan method yang diekspor
2. **Method yang diekspor** - `Greet` diawali huruf kapital sehingga dapat diakses oleh frontend
3. **Logika sederhana** - Menerima nama dan mengembalikan sapaan
4. **Keamanan tipe** - Tipe masukan dan keluaran ditentukan

@note{type="tip" title="Memahami Layanan dan Binding"}
**Layanan** adalah modul Go mandiri yang menyediakan fungsionalitas bagi frontend Anda. Layanan hanyalah struct Go biasa dengan method yang diekspor, yang Anda daftarkan di bidang `Services` pada konfigurasi aplikasi.

**Binding** adalah SDK TypeScript/JavaScript yang dihasilkan secara otomatis dan memungkinkan frontend Anda memanggil layanan tersebut. Saat Anda menjalankan `wails3 dev` atau `wails3 build`, Wails menganalisis layanan yang telah didaftarkan dan menghasilkan binding yang aman secara tipe di `frontend/bindings/`.

Anggap layanan sebagai API backend Anda dan binding sebagai pustaka klien yang berkomunikasi dengannya.

@end

### Mendaftarkan Layanan

Buka `main.go` dan temukan pendaftaran layanan:

```go {title="main.go" highlight="4-6"}
app := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Tindakan ini mendaftarkan `GreetService` Anda ke Wails sehingga semua method yang diekspor dapat diakses oleh frontend.

### Frontend

Buka `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.doGreet = async () => {
    const nameElement = document.getElementById('name');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        showToast(result);
    } catch (err) {
        console.error(err);
    }
};
```

Handler yang disederhanakan ini menggantikan `window.doGreet` milik proyek awal; pertahankan helper `showToast` yang ada dan kode startup lainnya. Tombol HTML yang dihasilkan memanggil `doGreet`, lalu `showToast` menampilkan hasilnya. Jalur impor mengikuti nama modul Go (`changeme`, kecuali jika diganti).

**Konsep utama:**

1. **Binding yang dihasilkan otomatis** - `GreetService` diimpor dari kode yang dihasilkan
2. **Pemanggilan yang aman secara tipe** - Nama dan signature method sesuai dengan kode Go Anda
3. **Asinkron secara default** - Semua pemanggilan Go mengembalikan Promise
4. **Penanganan kesalahan** - Kesalahan dari Go ditangkap dalam try/catch

@note{type="info" title="Di mana binding berada?"}
Binding yang dihasilkan berada di `frontend/bindings/`. Binding tersebut dibuat secara otomatis saat Anda menjalankan `wails3 dev` atau `wails3 build`.

**Jangan pernah mengedit berkas ini secara manual**—berkas tersebut dibuat ulang pada setiap build.

@end

## Sesuaikan Aplikasi Anda

Mari tambahkan fitur baru untuk memahami alur kerjanya.

### Tambahkan Fitur "Sapa Banyak Orang"

@steps
### Tambahkan method ke GreetService
Tambahkan kode ini ke `greetservice.go`:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = "Hello " + name + "!"
    }
    return greetings
}
```

### Aplikasi akan di-build ulang secara otomatis
Simpan berkas, lalu `wails3 dev` akan otomatis mem-build ulang kode Go Anda dan memulai ulang aplikasi.

@note{type="info" title="Build Ulang Otomatis"}
Perubahan kode Go memicu build ulang dan mulai ulang secara otomatis. Perubahan frontend dimuat ulang dengan hot reload tanpa memulai ulang aplikasi.

@end

### Gunakan di frontend
Tambahkan kode ini ke `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Buka konsol browser dan panggil `greetMany()`—Anda akan melihat array berisi sapaan.

@end

## Build untuk Produksi

Saat Anda siap mendistribusikan aplikasi:

```bash
wails3 build
```

**Yang dilakukan proses ini:**

- Mengompilasi kode Go dengan pengoptimalan
- Mem-build frontend untuk produksi (diminifikasi)
- Membuat berkas executable native di `bin/`

@tabs{sync-key="os"}
[Windows]
**Keluaran:** `bin/myapp.exe`

Klik dua kali untuk menjalankannya. Tidak memerlukan dependensi (WebView2 merupakan bagian dari Windows).

[macOS]
**Keluaran:** `bin/myapp.app`

Seret ke folder Applications atau klik dua kali untuk menjalankannya.

[Linux]
**Keluaran:** `bin/myapp`

Jalankan dengan `./bin/myapp` atau buat berkas `.desktop` untuk launcher Anda.

@end

@note{type="tip" title="Build Lintas Platform"}
Ingin membuat build untuk platform lain? Lihat [Build Lintas Platform →](/guides/build/cross-platform/)

@end

## Yang telah kita pelajari

**Struktur Proyek**

- `main.go` untuk backend Go
- `frontend/` untuk kode UI
- `Taskfile.yml` untuk tugas build

**Layanan**

- Buat struct Go dengan metode yang diekspor
- Daftarkan dengan `application.NewService()`
- Metode tersedia secara otomatis di frontend

**Binding**

- Definisi TypeScript yang dibuat secara otomatis
- Pemanggilan fungsi yang aman secara tipe
- Asinkron secara default (Promise)

**Alur Kerja Pengembangan**

- `wails3 dev` untuk hot reload
- Perubahan Go memicu build ulang dan mulai ulang secara otomatis
- Perubahan frontend langsung dimuat ulang secara otomatis

---

**Ada pertanyaan?** Bergabunglah ke [Discord](https://discord.gg/JDdSxwjhGf) dan tanyakan kepada komunitas.
