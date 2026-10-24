---
title: "Memutar Audio dan Video Lokal"
description: "Putar media bawaan di Linux menggunakan URL blob dengan ukuran terbatas dan lepaskan setelah selesai."
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

Gunakan `Media.SetSource` untuk memutar klip lokal pendek dalam aplikasi Wails v3. Di Linux, WebKitGTK menyerahkan pemutaran media kepada GStreamer, yang tidak dapat memuat URL `wails://` secara langsung. Fungsi pembantu menerima klip melalui stream Wails dan menetapkan URL blob ke pemutar.

Transfer desktop menggunakan transport aset Wails yang ada tanpa membuka socket pendengar. API bekerja dengan elemen audio dan video. Media HTTP/HTTPS biasa dapat langsung menggunakan `src` pemutar.

## Mendaftarkan berkas media

Sediakan sistem berkas berisi klip yang boleh diputar frontend, lalu daftarkan handler media pada stream bernama:

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

Panggil `registerMedia(app)` setelah membuat aplikasi dan sebelum memanggil `app.Run()`.

Untuk berkas pada disk, gunakan `os.OpenRoot(directory)` dan berikan `root.FS()` ke `media.NewHandler`. Biarkan root terbuka hingga `app.Run()` selesai, lalu tutup. Pilih direktori yang hanya berisi berkas yang boleh dibaca frontend; `os.Root` membatasi akses meskipun symlink mengarah ke luar direktori tersebut.

## Memuat klip

Buat pemutar:

```html
<video id="player" controls></video>
```

Untuk frontend yang menggunakan runtime npm:

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

Untuk aplikasi yang menggunakan runtime bawaan, ubah impor menjadi:

```javascript
import { Media } from '/wails/runtime.js';
```

Argumen kedua adalah nama stream terdaftar. Argumen ketiga adalah jalur berkas dengan pemisah garis miring, relatif terhadap root sistem berkasnya, seperti `welcome.mp4` atau `tutorials/intro.mp4`. Ini bukan URL atau jalur sistem operasi.

Promise diselesaikan ketika sumber ditetapkan. Pemutar kemudian mendekodenya. Tangani peristiwa `error` pada pemutar untuk mendeteksi codec yang tidak didukung, dan panggil `player.play()` dari interaksi pengguna jika perlu memulai pemutaran sendiri.

## Mengganti atau melepaskan klip

Panggil `Media.SetSource` lagi untuk mengganti klip. Ini membatalkan pemuatan sebelumnya yang masih tertunda untuk pemutar tersebut, sehingga respons lambat tidak menimpa pilihan terbaru. Klip sebelumnya tetap tersedia hingga penggantinya berhasil dimuat. URL blob sebelumnya kemudian dicabut.

Panggil `Media.ClearSource(player)` saat menutup pemutar atau melepas komponennya:

```javascript
Media.ClearSource(player);
```

Ini membatalkan pemuatan tertunda, mengatur ulang pemutar, dan melepaskan URL blob. Panggil sebelum menghapus elemen dari dokumen. Komponen framework harus memanggilnya dalam hook unmount atau disposal. Menghapus elemen saja tidak melepaskan URL blob. Gunakan fungsi pembantu ini secara konsisten untuk sumber pemutar.

Untuk markup `<video><source ...></video>`, berikan nama berkas pilihan ke `Media.SetSource(video, "media", name)`. Fungsi pembantu mengatur `src` pemutar induk, yang diprioritaskan atas anak `<source>`. Menghapus atribut tersebut memungkinkan browser mempertimbangkan anak-anak itu lagi; gunakan pemutar kosong jika mengelola semua sumber melalui API ini.

Objek audio yang tidak terhubung ke dokumen juga berfungsi:

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## Membatasi unduhan dan membatalkan pemuatan

Batas default adalah **32 MiB per sumber**. Anda dapat memilih batas bilangan bulat positif yang lebih kecil atau lebih besar dalam byte, selama masih dalam batas yang dikonfigurasi di Go:

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

Berkas yang terlalu besar ditolak dengan `RangeError`. Handler Go memeriksa ukuran berkas sebelum membaca isinya dan mentransfer paling banyak nilai terkecil antara batas konfigurasinya dan batas frontend. Frontend juga memeriksa ukuran yang diterima dan menolak transfer tidak lengkap. Pembatalan menolak dengan `AbortError`, atau alasan yang diberikan ke `AbortController.abort(reason)`.

Berkas dikirim dalam frame 64 KiB. Transport stream desktop Wails membatasi respons polling hingga 1 MiB, sehingga buffering respons lengkap oleh WebView2 tidak menampung seluruh berkas media dalam satu respons. Antrean stream memiliki mekanisme backpressure terbatas sendiri. Batas transport ini tidak mengubah perilaku blob lengkap pada fungsi pembantu pemutar.

**Seluruh berkas diunduh sebelum diputar.** Batas byte berlaku per berkas, bukan batas memori total aplikasi. Beberapa pemutar, klip sebelumnya saat penggantian, pembuatan blob, dan media yang didekode dapat menggunakan memori tambahan. Fungsi pembantu mentransfer saat dipanggil, terlepas dari pengaturan `preload` pemutar. Panggil ketika pengguna memilih untuk memuat klip.

Untuk berkas lokal besar, fungsi pembantu ini bukan solusi streaming. Menaikkan batas juga meningkatkan penggunaan memori. Media yang sudah dihosting melalui HTTP/HTTPS harus menggunakan pemuatan media native agar browser dapat melakukan streaming dan berpindah posisi dengan permintaan rentang.

## Memecahkan masalah pemutaran di Linux

- Jika pemutaran lokal langsung melaporkan **No URI handler implemented for "wails"**, muat klip dengan `Media.SetSource`. Ini memengaruhi stack GTK4 default maupun stack lama `-tags gtk3`.
- Jika pemuatan berhasil tetapi dekode gagal, periksa codec GStreamer yang terinstal pada sistem target. MP4 biasanya memerlukan dukungan video H.264 dan audio AAC; MP3 memerlukan dekoder MP3. Uji format yang Anda distribusikan pada distribusi yang didukung.
- Jika transfer gagal, periksa nama stream terdaftar, nama berkas relatif, dan izin sistem berkas. Mengubah berkas selama transfer dapat menyebabkan transfer tidak lengkap; coba lagi setelah penulisan berkas selesai.
- Jika aplikasi menetapkan Content Security Policy, izinkan origin aset Wails dalam `connect-src` dan `blob:` dalam `media-src`. Untuk kebijakan khusus lokal, direktif ini dapat berupa `connect-src 'self'; media-src 'self' blob:`. Pertahankan direktif lainnya.
- Jika berkas melebihi batas, pilih klip yang lebih pendek atau kecil, atau tetapkan batas eksplisit yang sesuai dengan anggaran memori aplikasi.

Jalankan [contoh audio-video](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video) dengan `go run .` untuk memeriksa sampel MP3 dan MP4 bawaan pada komputer Anda.
