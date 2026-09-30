---
title: "Kontrol LLM (MCP)"
description: "Izinkan agen LLM menguji dan mengontrol aplikasi Anda melalui Model Context Protocol"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Fitur Eksperimental"}
Server MCP bawaan masih bersifat eksperimental dan API-nya dapat berubah dalam rilis mendatang.

@end

Wails v3 memiliki server [Model Context Protocol](https://modelcontextprotocol.io) (MCP) bawaan yang memungkinkan agen LLM—Claude Code, asisten IDE, atau klien MCP apa pun—memeriksa, menguji, dan mengendalikan aplikasi Wails yang **sedang berjalan**.

Untuk otomatisasi siklus hidup proyek, gunakan server CLI `wails3 mcp` yang terpisah. Server ini memungkinkan agen memeriksa dan menginisialisasi proyek, menjalankan diagnostik, memulai build dan tugas pengembangan, menghasilkan binding, menjalankan tugas Taskfile bernama, serta mengambil keluaran tugas dalam batas tertentu. Secara default, server CLI dibatasi pada direktori saat ini dan tidak menyediakan eksekusi shell arbitrer. Lihat [dokumentasi MCP CLI](/guides/cli/#mcp) untuk detail transpor, autentikasi, dan alat.

Jika diaktifkan, agen yang terhubung ke aplikasi Anda dapat:

- **Mencantumkan dan mengontrol jendela** — ukuran, posisi, fokus, layar penuh, devtools, pemuatan ulang, …
- **Memeriksa DOM** — mencari elemen, mendapatkan HTML, mengambil snapshot struktural
- **Mengevaluasi JavaScript** — menjalankan kode arbitrer di dalam jendela mana pun dan mendapatkan hasilnya
- **Menyimulasikan input pengguna** — gerakan mouse, klik, seret, dan gulir yang ditampilkan dengan **kursor layar beranimasi** sehingga Anda dapat melihat agen bekerja
- **Mengetik dan menekan tombol** — peristiwa per karakter yang realistis dan berfungsi dengan input terkontrol React
- **Memanggil metode Go yang terikat** serta memancarkan/menunggu peristiwa aplikasi

## Cara kerjanya

Server MCP dikompilasi ke dalam aplikasi Anda hanya jika tag build **`mcp`** tersedia. Tanpa tag tersebut, kode server sama sekali tidak terdapat dalam biner—tanpa overhead runtime, tanpa port terbuka, dan tanpa permukaan serangan.

Jika tag tersebut tersedia, server dimulai secara otomatis di dalam `App.Run()`, secara default melakukan bind ke `127.0.0.1:9099`, dan mencatat endpoint-nya dalam log. Tidak diperlukan kode pengguna.

## Tutorial

### Langkah 1 — tulis aplikasi Wails biasa

MCP tidak memerlukan impor atau pendaftaran. Buat aplikasi seperti biasa:

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Langkah 2 — lakukan build atau jalankan dengan tag `mcp`

@tabs
[Wails CLI (disarankan)]
Tetapkan `WAILS_MCP=1` agar Wails CLI menambahkan tag `mcp` untuk Anda:

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go Langsung]
Teruskan tag secara langsung ke `go run` atau `go build`:

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows (PowerShell)]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

Saat dimulai, aplikasi mencatat endpoint MCP dalam log:

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### Langkah 3 — hubungkan klien

Server menggunakan **transpor HTTP streamable MCP**. Hubungkan dengan klien apa pun yang kompatibel dengan MCP.

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

Kemudian, minta Claude berinteraksi dengan aplikasi Anda:

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code (GitHub Copilot)]
Tambahkan ke `.vscode/settings.json`:

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[Klien lain]
Arahkan klien MCP apa pun yang mendukung transpor HTTP streamable ke:

```
http://127.0.0.1:9099/mcp
```

@end

### Langkah 4 — jalankan sesi pengujian

Minta agen menguji aplikasi Anda. Berikut beberapa contoh prompt:

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## Konfigurasi

Semua konfigurasi dilakukan melalui variabel lingkungan—tidak diperlukan perubahan kode.

| Variabel lingkungan | Default | Deskripsi |
| --- | --- | --- |
| `WAILS_MCP` | (tidak ditetapkan) | Tetapkan ke `1`, `true`, `on`, atau `yes` untuk menambahkan tag build `mcp` secara otomatis saat menggunakan Wails CLI. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Antarmuka yang akan menjadi tujuan bind. Bind selain loopback memerlukan `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | tidak ditetapkan | Token bearer bersifat opsional pada loopback, tetapi wajib pada alamat bind lainnya. Klien mengirim `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Port untuk mendengarkan koneksi. Tetapkan ke `0` agar port kosong ditetapkan secara acak (dicetak dalam log). |
| `WAILS_MCP_TIMEOUT` | `30000` | Batas waktu default evaluasi JS dalam **milidetik**. |
| `WAILS_MCP_HIDE_CURSOR` | (tidak ditetapkan) | Tetapkan ke `1` atau `true` untuk menonaktifkan overlay kursor beranimasi. |

Contoh—port khusus dan batas waktu 60 detik:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Alat yang tersedia

| Alat | Tujuan |
| --- | --- |
| `app_info` | Informasi aplikasi: platform, arsitektur, semua jendela, endpoint MCP |
| `windows_list` | Mencantumkan semua jendela beserta geometri dan statusnya |
| `window_control` | Memfokuskan, mengubah ukuran, memindahkan, menampilkan dalam layar penuh, membuka devtools, memuat ulang, menetapkan URL, … (22 tindakan) |
| `js_eval` | Mengevaluasi JavaScript dalam jendela (isi async, `return` untuk nilai) |
| `dom_html` | Dapatkan HTML halaman atau elemen tertentu |
| `dom_query` | Temukan elemen berdasarkan selektor CSS — tag, teks, batas, visibilitas |
| `screenshot_dom` | Snapshot struktural halaman yang terlihat (berbasis DOM, tanpa piksel) |
| `mouse_move` | Animasikan kursor ke suatu titik atau selektor CSS |
| `mouse_click` | Klik dengan kursor beranimasi (kiri/kanan/tengah, klik ganda, tombol pengubah) |
| `mouse_drag` | Seret dengan kursor beranimasi (mendukung elemen seret dan lepas HTML5) |
| `mouse_scroll` | Gulir pada suatu titik atau elemen |
| `keyboard_type` | Ketik teks karakter demi karakter dengan peristiwa yang realistis |
| `keyboard_press` | Tekan satu tombol (Enter, Tab, Escape, ArrowDown, …) dengan tombol pengubah opsional |
| `call_bound_method` | Panggil metode layanan Go yang terikat, misalnya `main.GreetService.Greet` |
| `emit_event` | Kirim peristiwa aplikasi Wails |
| `wait_for_event` | Tunggu peristiwa aplikasi Wails dan kembalikan datanya |

### Dukungan multi-jendela

Semua alat yang bekerja pada jendela menerima argumen opsional `window` yang berisi **nama** jendela (ditetapkan melalui `WebviewWindowOptions.Name`). Jika argumen ini dihilangkan, alat menargetkan jendela yang sedang memiliki fokus, atau jendela pertama jika tidak ada yang memiliki fokus.

```
List all windows, then click the "New" button in the window named "editor".
```

### Memilih elemen

Alat tetikus dan papan ketik menerima salah satu dari:

- **Selektor CSS** — `selector: "#submit-btn"` (elemen otomatis digulir hingga terlihat)
- **Koordinat** — `x: 400, y: 300` (piksel CSS relatif terhadap viewport)

Untuk operasi seret, awali dengan `from_` dan `to_`:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Keamanan

@note{type="caution"}
Server MCP memberikan kendali terprogram penuh atas aplikasi Anda. Siapa pun yang dapat mengakses alatnya dapat membaca DOM, mengevaluasi JavaScript, mengeklik tombol, dan memanggil metode Go.

@end

- Secara default, server melakukan bind ke `127.0.0.1`. Origin browser harus berupa origin loopback HTTP(S); origin opak (`null`), yang formatnya tidak valid, dan origin asing akan ditolak.
- Klien MCP native tanpa header tetap didukung. Tanpa `WAILS_MCP_TOKEN`, proses lokal dan origin browser lokal yang diizinkan dianggap tepercaya; pemeriksaan origin bukanlah autentikasi. Tetapkan token dengan entropi tinggi untuk mewajibkan autentikasi bearer bagi semua panggilan `/mcp`. Konfigurasikan token yang sama dalam header `Authorization: Bearer <token>` milik klien. Permintaan preflight tidak memerlukan token tersebut.
- Callback `/eval-result` menggunakan ID per evaluasi yang tidak dapat diprediksi, bukan token bearer klien, sehingga pengiriman hasil webview tetap kompatibel.
- Build produksi sebaiknya **tidak** menyertakan tag `mcp`. Wails CLI hanya menambahkannya jika `WAILS_MCP=1` ditetapkan secara eksplisit, dan `wails3 build` default sama sekali tidak memuat kode server.
- Jika Anda perlu mengekspos server pada antarmuka non-loopback (misalnya untuk pengujian LAN), tetapkan `WAILS_MCP_HOST=0.0.0.0` dan `WAILS_MCP_TOKEN` dengan entropi tinggi; proses startup akan gagal tanpa token. Gunakan tunnel terenkripsi atau proksi yang mengakhiri TLS pada jaringan yang tidak tepercaya karena listener bawaan menggunakan HTTP.

## Aplikasi contoh

Aplikasi playground lengkap yang mendemonstrasikan semua alat tersedia di [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp). Aplikasi ini mencakup:

- Penghitung dengan tombol untuk menaikkan nilai dan mereset
- Input nama dengan metode terikat Greet, Add, dan Shout
- Sumber dan target seret dan lepas HTML5
- Daftar yang dapat digulir (50 item)
- Log peristiwa

Jalankan dengan:

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Kemudian hubungkan Claude Code atau klien MCP apa pun ke `http://127.0.0.1:9099/mcp` dan minta klien tersebut menguji UI.

## Umpan balik

Server MCP bawaan merupakan eksperimen, dan umpan balik Anda menentukan arah pengembangannya. Jika Anda mencobanya, kami ingin mengetahui klien dan alat yang Anda gunakan, apa yang Anda harapkan dan apa yang sebenarnya terjadi, serta apakah membiarkan agen mengendalikan aplikasi Anda bermanfaat — laporan yang paling berguna menyebutkan dengan tepat apa yang Anda jalankan. Sampaikan kepada kami dalam [diskusi umpan balik MCP Server](https://github.com/wailsapp/wails/discussions/5692).
