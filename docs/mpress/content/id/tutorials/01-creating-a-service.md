---
title: "Layanan Kode QR"
description: "Buat layanan kode QR untuk mempelajari layanan Wails"
slug: "tutorials/01-creating-a-service"
sourcePath: "tutorials/01-creating-a-service.md"
---

**Layanan** di Wails adalah struct Go yang berisi logika bisnis yang ingin Anda sediakan untuk frontend. Layanan membantu menjaga kode tetap terorganisasi dengan mengelompokkan fungsionalitas terkait.

Anggap layanan sebagai kumpulan metode yang dapat dipanggil oleh kode JavaScript Anda. Setiap metode publik pada layanan dapat dipanggil dari frontend setelah binding dibuat.

Dalam tutorial ini, kita akan membuat layanan pembuat kode QR untuk mendemonstrasikan konsep-konsep tersebut. Setelah selesai, Anda akan memahami cara membuat layanan, mengelola dependensi, dan menghubungkan kode Go ke frontend.

<br/>

@steps
### Buat File Layanan QR
Buat file baru bernama `qrservice.go` di direktori aplikasi Anda:

```go {title="qrservice.go"}
package main

import (
    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}
```

**Yang terjadi di sini:**

- `QRService` adalah struct kosong yang akan menampung metode pembuatan kode QR kita
- `NewQRService()` adalah fungsi konstruktor yang membuat instans baru layanan kita
- `Generate()` adalah metode yang menerima teks dan ukuran, lalu mengembalikan kode QR sebagai array byte PNG
- Metode tersebut mengembalikan `([]byte, error)` sesuai konvensi Go yang menempatkan error sebagai nilai terakhir yang dikembalikan
- Kita menggunakan paket `github.com/skip2/go-qrcode` untuk menangani proses pembuatan kode QR

 <br/>

### Daftarkan Layanan
Membuat layanan saja tidak cukup—kita perlu **mendaftarkannya** ke aplikasi Wails agar aplikasi mengetahui keberadaan layanan tersebut dan dapat membuat binding untuknya.

Pendaftaran dilakukan di `main.go` saat Anda membuat aplikasi. Teruskan instans layanan Anda ke opsi `Services`:

```go {title="main.go" ins="7-9"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewService(NewQRService()),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Yang terjadi di sini:**

- `application.NewService()` membungkus layanan Anda agar dapat dikelola oleh Wails
- Kita memanggil `NewQRService()` untuk membuat instans layanan kita
- Layanan ditambahkan ke slice `Services` dalam opsi aplikasi
- Wails sekarang akan memindai metode publik dalam layanan ini agar tersedia untuk frontend

 <br/>

### Instal Dependensi
Kita telah merujuk paket `github.com/skip2/go-qrcode` dalam kode, tetapi belum benar-benar mengunduhnya. Go perlu mengetahui dependensi ini dan mengunduhnya ke proyek Anda.

Jalankan perintah ini di terminal dari direktori proyek Anda:

```bash
go mod tidy
```

**Yang terjadi di sini:**

- `go mod tidy` memindai file Go Anda untuk mencari pernyataan import
- Perintah tersebut mengunduh paket yang belum tersedia (seperti `go-qrcode`) dan menambahkannya ke `go.mod`
- Perintah tersebut juga menghapus dependensi yang tidak lagi digunakan
- Hal ini memastikan proyek Anda memiliki semua kode yang diperlukan agar berhasil dikompilasi

Anda seharusnya melihat output yang menunjukkan bahwa paket kode QR telah diunduh dan ditambahkan ke proyek Anda.

 <br/>

### Buat Binding
Agar metode ini dapat dipanggil dari frontend, kita perlu membuat binding. Anda dapat melakukannya dengan menjalankan `wails generate bindings` di direktori root proyek Anda.

@note{type="info"}
Saat pertama kali menjalankannya dalam suatu proyek, generator binding akan menganalisis kode dan dependensi Anda secara menyeluruh. Proses ini terkadang memerlukan waktu sedikit lebih lama dari perkiraan, tetapi proses berikutnya akan jauh lebih cepat.

@end

Setelah menjalankannya, Anda seharusnya melihat output yang mirip dengan berikut ini di terminal:

```bash
 % wails3 generate bindings
 INFO  Processed: 337 Packages, 1 Service, 1 Method, 0 Enums, 0 Models in 740.196125ms.
 INFO  Output directory: /Users/leaanthony/myproject/frontend/bindings
```

Di dalam direktori frontend, kini akan ada direktori baru bernama `bindings`:

```bash
frontend/
└── bindings
    └── changeme
        ├── index.js
        └── qrservice.js
```

@note{type="tip" title="Kiat Profesional"}
Saat membangun aplikasi menggunakan `wails3 build`, binding akan dibuat secara otomatis dan selalu diperbarui.

@end

 <br/>

### Memahami Binding
Mari kita lihat binding yang dibuat di `bindings/changeme/qrservice.js`:

```js {title="bindings/changeme/qrservice.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 /**
  * QRService handles QR code generation
  * @module
  */

 // eslint-disable-next-line @typescript-eslint/ban-ts-comment
 // @ts-ignore: Unused imports
 import {Call as $Call, Create as $Create} from "@wailsio/runtime";

 /**
  * Generate creates a QR code from the given text
  * @param {string} text
  * @param {number} size
  * @returns {Promise<string> & { cancel(): void }}
  */
 export function Generate(text, size) {
     let $resultPromise = /** @type {any} */($Call.ByID(3576998831, text, size));
     let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
         return $Create.ByteSlice($result);
     }));
     $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
     return $typingPromise;
 }
```

Kita dapat melihat bahwa binding dibuat untuk metode `Generate`. Nama parameter dipertahankan, begitu pula komentarnya. JSDoc juga dibuat untuk metode tersebut guna menyediakan informasi tipe bagi IDE Anda.

@note{type="info"}
Anda tidak perlu memahami binding yang dibuat secara menyeluruh, tetapi penting untuk memahami cara kerjanya.

@end

Binding menyediakan:

- Fungsi yang setara dengan metode Go Anda
- Konversi otomatis antara tipe Go dan JavaScript
- Operasi asinkron berbasis Promise
- Informasi tipe dalam bentuk komentar JSDoc

@note{type="tip" title="TypeScript"}
Generator binding juga mendukung pembuatan binding TypeScript. Anda dapat melakukannya dengan menjalankan `wails3 generate bindings -ts`.

@end

Layanan yang dibuat diekspor ulang oleh file `index.js`:

```js {title="bindings/changeme/index.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 import * as QRService from "./qrservice.js";
 export {
     QRService
 };
```

Selanjutnya, Anda dapat mengaksesnya melalui jalur import sederhana `./bindings/changeme` yang hanya terdiri dari jalur paket Go Anda, tanpa menentukan nama file apa pun.

@note{type="info"}
Jalur import sederhana hanya tersedia saat menggunakan bundler frontend. Jika Anda lebih memilih frontend vanilla yang tidak menggunakan bundler, Anda harus mengimpor `index.js` atau `qrservice.js` secara manual.

@end

 <br/>

### Gunakan Binding di Frontend
Sekarang kita dapat memanggil layanan Go dari JavaScript! Binding yang dibuat memudahkan hal ini sekaligus menjaga keamanan tipe.

Perbarui `frontend/src/main.js` agar menggunakan binding baru:

```js {title="frontend/src/main.js"}
 import { QRService } from './bindings/changeme';

 async function generateQR() {
     const text = document.getElementById('text').value;
     if (!text) {
         alert('Please enter some text');
         return;
     }

     try {
         // Generate QR code as base64
         const qrCodeBase64 = await QRService.Generate(text, 256);

         // Display the QR code
         const qrDiv = document.getElementById('qrcode');
         qrDiv.src = `data:image/png;base64,${qrCodeBase64}`;

     } catch (err) {
         console.error('Failed to generate QR code:', err);
         alert('Failed to generate QR code: ' + err);
     }
 }

 export function initializeQRGenerator() {
     const button = document.getElementById('generateButton');
     button.addEventListener('click', generateQR);
 }
```

**Yang terjadi di sini:**

- Kita mengimpor `QRService` dari binding yang dibuat
- `QRService.Generate()` memanggil metode Go kita—pemanggilan ini mengembalikan Promise, sehingga kita menggunakan `await`
- Metode Go mengembalikan `[]byte`, yang secara otomatis dikonversi oleh Wails menjadi string base64 untuk JavaScript
- Kita membuat URL data dengan string base64 untuk menampilkan gambar PNG
- Blok `try/catch` menangani setiap error dari sisi Go (seperti input yang tidak valid)
- Jika kode Go kita mengembalikan error, Promise akan ditolak dan kita menangkapnya di sini

Sekarang perbarui `index.html` agar menggunakan binding baru dalam fungsi `initializeQRGenerator`:

```html {title="frontend/src/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>QR Code Generator</title>
            <style>
                body {
                font-family: Arial, sans-serif;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                height: 100vh;
                margin: 0;
            }
                #qrcode {
                margin-bottom: 20px;
                width: 256px;
                height: 256px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
                #controls {
                display: flex;
                gap: 10px;
            }
                #text {
                padding: 5px;
            }
                #generateButton {
                padding: 5px 10px;
                cursor: pointer;
            }
            </style>
</head>
<body>
<img id="qrcode"/>
<div id="controls">
    <input type="text" id="text" placeholder="Enter text">
        <button id="generateButton">Generate QR Code</button>
</div>

<script type="module">
    import { initializeQRGenerator } from './main.js';
    document.addEventListener('DOMContentLoaded', initializeQRGenerator);
</script>
</body>
</html>
```

Jalankan `wails3 dev` untuk memulai server pengembangan. Setelah beberapa detik, aplikasi seharusnya terbuka.

Masukkan teks, lalu klik tombol "Buat Kode QR". Anda seharusnya melihat kode QR di bagian tengah halaman:

![Kode QR](/assets/qr1.png)

 <br/>

 <br/>

### Pendekatan Alternatif: Handler HTTP
Sejauh ini, kita telah membahas area berikut:

- Membuat Layanan baru
- Membuat Binding
- Menggunakan Binding dalam kode Frontend kita

**Mengapa menggunakan handler HTTP?**

Binding metode sangat cocok untuk operasi data, tetapi ada pendekatan alternatif untuk menyajikan file, gambar, atau media lainnya. Alih-alih mengonversi semuanya ke base64 dan mengirimkannya melalui binding, Anda dapat membuat layanan bertindak seperti server web mini.

Ini berguna ketika:

- Anda menyajikan gambar, video, atau file berukuran besar
- Anda ingin menggunakan tag HTML `<img>` atau `<video>` standar dengan atribut `src`
- Anda memerlukan akses URL langsung ke sumber daya

Jika layanan Anda mengimplementasikan metode `ServeHTTP(w http.ResponseWriter, r *http.Request)` standar Go, Wails dapat membuatnya dapat diakses sebagai endpoint HTTP. Mari perluas layanan kode QR kita agar mendukung hal ini:

```go {title="qrservice.go" ins="4-5,37-65"}
package main

import (
    "net/http"
    "strconv"

    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Yang terjadi di sini:**

- `ServeHTTP` adalah antarmuka standar Go untuk menangani permintaan HTTP
- Kita mengurai parameter kueri dari URL (`?text=hello&size=256`)
- Kita memanggil metode `Generate()` yang sudah ada untuk membuat kode QR
- Kita menetapkan jenis konten ke `image/png` agar browser mengetahui bahwa konten tersebut adalah gambar
- Kita menulis byte PNG mentah langsung ke respons—tanpa memerlukan base64!

Sekarang perbarui `main.go` untuk menentukan rute tempat layanan kode QR dapat diakses:

```go {title="main.go" ins="8-10"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                 Route: "/qrservice",
             }),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Yang terjadi di sini:**

- Kita menambahkan `application.ServiceOptions` untuk mengonfigurasi cara layanan diekspos
- `Route: "/qrservice"` membuat handler HTTP dapat diakses di `/qrservice`
- Sekarang setiap permintaan ke `/qrservice?text=hello` akan memanggil metode `ServeHTTP` kita
- Tanpa menetapkan `Route`, fungsi handler HTTP dinonaktifkan

@note{type="info"}
Jika Anda tidak menetapkan opsi `Route` secara eksplisit, handler HTTP tidak akan dapat diakses dari frontend.

@end

Terakhir, perbarui `main.js` agar menggunakan `src` gambar sederhana, bukan pengodean base64:

```js {title="frontend/src/main.js"}
async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Make the image source the path to the QR code service, passing the text
    img.src = `/qrservice?text=${encodeURIComponent(text)}`
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Yang terjadi di sini:**

- Kita menghapus impor dan pemanggilan `await QRService.Generate()`
- Sebagai gantinya, kita cukup menetapkan `img.src` agar mengarah ke endpoint HTTP kita
- `encodeURIComponent()` meng-escape karakter khusus dalam URL dengan aman
- Browser secara otomatis membuat permintaan HTTP GET saat kita menetapkan `src`
- Cara ini lebih sederhana dan efisien untuk gambar—tanpa memerlukan konversi base64!

Menjalankan aplikasi kembali seharusnya menghasilkan kode QR yang sama:

![Kode QR](/assets/qr1.png)

 <br/>

 <br/>

### Mendukung Konfigurasi Dinamis
**Masalah pada rute yang di-hardcode:**

Dalam contoh di atas, kita menggunakan rute `/qrservice` yang di-hardcode dalam kode JavaScript. Hal ini menciptakan keterikatan yang erat antara konfigurasi Go dan kode frontend Anda.

Jika Anda mengedit `main.go` dan mengubah opsi `Route` tanpa memperbarui `main.js`, aplikasi akan berhenti berfungsi:

```go {title="main.go" ins="3"}
        // ...
            application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                Route: "/services/qr",
            }),
        // ...
```

Rute yang di-hardcode dapat digunakan untuk aplikasi sederhana, tetapi membuat kode Anda rapuh dan lebih sulit dipelihara.

**Solusinya: Konfigurasi dinamis**

Binding metode dan handler HTTP dapat digunakan bersama! Kita dapat menggunakan binding untuk memberi tahu frontend rute yang harus digunakan, sehingga konfigurasi menjadi dinamis dan jalur yang di-hardcode dapat dihilangkan.

Berikut cara kerjanya:

1. Metode siklus hidup `ServiceStartup` dijalankan saat aplikasi Anda dimulai
2. Kita menyimpan rute yang dikonfigurasi dari opsi
3. Kita menambahkan metode `URL()` yang dapat dipanggil frontend untuk mendapatkan rute yang benar
4. Sekarang frontend meminta rutenya kepada layanan Go, bukan menebaknya

Pertama, implementasikan antarmuka `ServiceStartup` dan tambahkan metode `URL` baru:

```go {title="qrservice.go" ins="4,6,10,15,23-27,46-55"}
package main

import (
    "context"
    "net/http"
    "net/url"
    "strconv"

    "github.com/skip2/go-qrcode"
    "github.com/wailsapp/wails/v3/pkg/application"
)

// QRService handles QR code generation
type QRService struct {
    route string
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// ServiceStartup runs at application startup.
func (s *QRService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.route = options.Route
    return nil
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

// URL returns an URL that may be used to fetch
// a QR code with the given text and size.
// It returns an error if the HTTP handler is not available.
func (s *QRService) URL(text string, size int) (string, error) {
    if s.route == "" {
        return "", errors.New("http handler unavailable")
    }

    return fmt.Sprintf("%s?text=%s&size=%d", s.route, url.QueryEscape(text), size), nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Yang terjadi di sini:**

- Kita menambahkan field `route` untuk menyimpan rute yang dikonfigurasi dari `ServiceStartup`
- `ServiceStartup(ctx, options)` dipanggil saat aplikasi dimulai—kita menyimpan rute di sini
- Metode `URL()` menyusun URL lengkap beserta parameter kueri
- Jika tidak ada rute yang dikonfigurasi (rute kosong), kita mengembalikan error
- `url.QueryEscape()` mengodekan teks dengan aman untuk digunakan dalam URL
- Metode ini akan tersedia bagi frontend melalui binding

Sekarang perbarui `main.js` agar menggunakan metode `URL` sebagai pengganti jalur yang di-hardcode:

```js {title="frontend/src/main.js" ins="1,11-12"}
import { QRService } from "./bindings/changeme";

async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Invoke the URL method to obtain an URL for the given text.
    img.src = await QRService.URL(text, 256);
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Yang terjadi di sini:**

- Kita mengimpor `QRService` agar dapat menggunakan binding kembali
- Alih-alih melakukan hardcode terhadap `/qrservice`, kita memanggil `await QRService.URL(text, 256)`
- Layanan Go menyusun URL dengan rute dan parameter yang benar
- Sekarang, jika Anda mengubah rute di `main.go`, frontend otomatis menggunakan rute baru
- Tidak perlu lagi melakukan sinkronisasi manual antara konfigurasi Go dan kode frontend!

Semua seharusnya berfungsi seperti pada contoh sebelumnya, namun mengubah rute layanan di `main.go` tidak akan lagi menyebabkan frontend berhenti berfungsi.

@note{type="info"}
Jika metode Go mengembalikan error yang bukan nil, promise di sisi JS akan ditolak dan pernyataan await akan melempar pengecualian.

@end

 <br/>

@end
