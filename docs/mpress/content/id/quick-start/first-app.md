---
title: "Aplikasi Pertama Anda"
description: "Bangun aplikasi Wails yang berfungsi dalam 10 menit"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Kita akan membangun aplikasi sapaan sederhana yang menunjukkan konsep inti Wails:

- Backend Go yang mengelola logika
- Frontend yang memanggil fungsi Go
- Binding yang aman secara tipe
- Hot reload selama pengembangan

**Waktu penyelesaian:** 10 menit

@note{type="tip" title="Kiat Performa untuk Pengguna Windows 11"}
Pertimbangkan untuk menggunakan [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) guna menyimpan proyek Anda. Dev Drive dioptimalkan untuk beban kerja pengembang dan dapat mempersingkat waktu build serta meningkatkan kecepatan akses disk secara signifikan hingga 30% dibandingkan dengan drive NTFS biasa.

@end

## Buat Proyek Anda

@steps
### Buat proyek
```bash
wails3 init -n myapp
cd myapp
```

Tindakan ini membuat proyek baru dengan templat Vanilla + Vite bawaan (HTML/CSS/TypeScript dengan bundler Vite).

@note{type="tip" title="Templat Lainnya"}
Coba `-t react`, `-t vue`, atau `-t svelte` untuk framework pilihan Anda. Secara bawaan, semuanya menggunakan TypeScript; untuk JavaScript biasa, gunakan `-t vanilla-js` atau `-t react-js`. Jalankan `wails3 init -l` untuk melihat semua templat yang tersedia, atau [gunakan framework frontend Anda sendiri](/guides/dev/frontend-frameworks/).

@end

### Pahami struktur proyek
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### Jalankan aplikasi
```bash
wails3 dev
```

@note{type="info" title="Jalankan untuk Pertama Kalinya"}
Proses pertama kali dijalankan mungkin memerlukan waktu lebih lama dari perkiraan karena menginstal dependensi frontend, menghasilkan binding, dan sebagainya. Proses berikutnya akan jauh lebih cepat.

@end

Aplikasi terbuka dan menampilkan antarmuka sapaan. Masukkan nama Anda lalu klik "Greet"—backend Go akan memproses input Anda dan mengembalikan sapaan.

@end

## Cara Kerjanya

Mari pahami kode yang membuatnya berfungsi.

### Backend Go

Buka `greetservice.go`:

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**Konsep utama:**

1. **Layanan**—Struct Go dengan metode yang diekspor
2. **Metode yang diekspor**—`Greet` diawali huruf kapital sehingga tersedia bagi frontend
3. **Logika sederhana**—Menerima nama dan mengembalikan sapaan
4. **Keamanan tipe**—Tipe input dan output telah ditentukan

@note{type="tip" title="Memahami Layanan dan Binding"}
**Layanan** adalah modul Go mandiri yang menyediakan fungsionalitas bagi frontend Anda. Layanan hanyalah struct Go biasa dengan metode yang diekspor, yang Anda daftarkan dalam bidang `Services` pada konfigurasi aplikasi.

**Binding** adalah SDK TypeScript/JavaScript yang dibuat secara otomatis agar frontend Anda dapat memanggil layanan tersebut. Saat Anda menjalankan `wails3 dev` atau `wails3 build`, Wails menganalisis layanan yang telah didaftarkan dan menghasilkan binding yang aman secara tipe di `frontend/bindings/`.

Anggap layanan sebagai API backend Anda dan binding sebagai pustaka klien yang berkomunikasi dengannya.

@end

### Mendaftarkan Layanan

Buka `main.go` dan temukan pendaftaran layanan:

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Tindakan ini mendaftarkan `GreetService` Anda ke Wails sehingga semua metode yang diekspor tersedia bagi frontend.

### Frontend

Buka `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**Konsep utama:**

1. **Binding yang dibuat secara otomatis**—`GreetService` diimpor dari kode yang dihasilkan
2. **Pemanggilan yang aman secara tipe**—Nama dan signature metode sesuai dengan kode Go Anda
3. **Asinkron secara bawaan**—Semua pemanggilan Go mengembalikan Promise
4. **Penanganan kesalahan**—Kesalahan dari Go ditangkap dalam try/catch

@note{type="info" title="Di mana binding berada?"}
Binding yang dihasilkan berada di `frontend/bindings/`. Binding tersebut dibuat secara otomatis saat Anda menjalankan `wails3 dev` atau `wails3 build`.

**Jangan pernah mengedit berkas ini secara manual**—berkas tersebut dibuat ulang pada setiap build.

@end

## Sesuaikan Aplikasi Anda

Mari tambahkan fitur baru untuk memahami alur kerjanya.

### Tambahkan Fitur "Sapa Banyak Orang"

@steps
### Tambahkan metode ke GreetService
Tambahkan kode berikut ke `greetservice.go`:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### Aplikasi akan di-build ulang secara otomatis
Simpan berkas tersebut dan `wails3 dev` akan secara otomatis mem-build ulang kode Go Anda serta memulai ulang aplikasi.

@note{type="info" title="Build Ulang Otomatis"}
Perubahan pada kode Go memicu build ulang dan pemulaian ulang secara otomatis. Perubahan frontend dimuat ulang secara langsung tanpa memulai ulang aplikasi.

@end

### Gunakan di frontend
Tambahkan kode berikut ke `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Buka konsol browser dan panggil `greetMany()`—Anda akan melihat array sapaan.

@end

## Build untuk Produksi

Saat Anda siap mendistribusikan aplikasi:

```bash
wails3 build
```

**Yang dilakukan proses ini:**

- Mengompilasi kode Go dengan pengoptimalan
- Mem-build frontend untuk produksi (diperkecil)
- Membuat executable native di `bin/`

@tabs{sync-key="os"}
[Windows]
**Output:** `bin/myapp.exe`

Klik dua kali untuk menjalankannya. Tidak memerlukan dependensi (WebView2 merupakan bagian dari Windows).

[macOS]
**Output:** `bin/myapp.app`

Seret ke folder Applications atau klik dua kali untuk menjalankannya.

[Linux]
**Output:** `bin/myapp`

Jalankan dengan `./bin/myapp` atau buat file `.desktop` untuk peluncur Anda.

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

- Buat struct Go dengan method yang diekspor
- Daftarkan dengan `application.NewService()`
- Method tersedia secara otomatis di frontend

**Binding**

- Definisi TypeScript yang dibuat secara otomatis
- Pemanggilan fungsi yang aman terhadap tipe
- Asinkron secara default (Promise)

**Alur Kerja Pengembangan**

- `wails3 dev` untuk hot reload
- Perubahan Go secara otomatis memicu build ulang dan mulai ulang
- Perubahan frontend langsung dimuat ulang dengan hot reload

---

**Ada pertanyaan?** Bergabunglah dengan [Discord](https://discord.gg/JDdSxwjhGf) dan tanyakan kepada komunitas.
