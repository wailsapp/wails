---
title: "Protokol Manifes Pembaruan"
description: "Protokol JSON terbuka yang digunakan aplikasi Wails untuk menemukan dan memverifikasi pembaruan mandiri, serta dapat disajikan dari host file statis atau server pembaruan dinamis apa pun."
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Protokol Manifes Pembaruan Wails adalah kontrak JSON terbuka berukuran kecil antara aplikasi Wails dan sumber pembaruan. Apa pun yang dapat menyajikan file JSON melalui HTTPS dapat menyajikan pembaruan Wails: bucket S3, GitHub Pages, CDN, atau server pembaruan dinamis yang membatasi rilis berdasarkan lisensi.

Sisi klien disertakan dalam framework sebagai penyedia `endpoint` (`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`). Halaman ini merupakan referensi format komunikasi bagi siapa pun yang mengimplementasikan sisi server.

## Tujuan desain

1. **Ramah terhadap host statis.** Satu file manifes per kanal, yang mencantumkan artefak untuk setiap platform, merupakan implementasi lengkap. Tidak memerlukan kode server.
2. **Ramah terhadap server dinamis.** Klien mengirimkan `platform`, `arch`, `version`, dan `channel` pada setiap pemeriksaan sehingga server dapat merespons dengan tepat satu artefak, menerapkan aturan lisensi, atau mengembalikan `204 No Content` ketika pemanggil sudah menggunakan versi terbaru.
3. **Utamakan verifikasi.** Manifes memuat checksum dan tanda tangan untuk setiap artefak, dan pemutakhir Wails memverifikasinya dengan kunci publik yang disematkan dalam biner aplikasi pada waktu build. Sumber pembaruan tidak pernah memilih akar kepercayaannya sendiri.

## Permintaan

Klien mengirimkan `GET` ke URL manifes yang dikonfigurasi dengan `Accept: application/json` beserta header apa pun yang dikonfigurasi oleh aplikasi (misalnya `Authorization: License <key>`).

URL dapat menyertakan placeholder yang diganti oleh klien pada setiap pemeriksaan:

| Placeholder | Diganti dengan |
| --- | --- |
| `{{platform}}` | OS yang sedang berjalan sebagai nilai `GOOS` Go (`darwin`, `windows`, `linux`) |
| `{{arch}}` | Arsitektur yang sedang berjalan sebagai nilai `GOARCH` Go (`amd64`, `arm64`, ...) |
| `{{version}}` | Versi yang saat ini terinstal |
| `{{channel}}` | Kanal rilis yang dikonfigurasi, jika ditetapkan |

Setiap nilai dari keempat nilai tersebut yang tidak digunakan oleh placeholder ditambahkan sebagai parameter kueri dengan nama yang sama (`channel` hanya jika dikonfigurasi). Oleh karena itu, kedua konfigurasi berikut valid dan setara:

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

Host statis cukup mengabaikan parameter kueri yang diterimanya.

## Respons

| Status | Arti |
| --- | --- |
| `200 OK` | Manifes disertakan setelahnya. Klien menentukan apakah manifes tersebut merupakan pemutakhiran. |
| `204 No Content` | Server telah membandingkan versi dan versi pemanggil sudah terbaru. |
| `404 Not Found` | Tidak ada yang dipublikasikan (diperlakukan sama seperti sudah menggunakan versi terbaru). |
| Status lainnya | Terjadi kesalahan. Pemutakhir beralih ke penyedia berikutnya yang dikonfigurasi. |

Isi `200` adalah dokumen manifes:

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### Bidang tingkat teratas

| Bidang | Tipe | Wajib | Catatan |
| --- | --- | --- | --- |
| `schemaVersion` | int | tidak | Versi protokol. Jika dihilangkan, berarti `1`. Klien menolak nilai yang lebih baru daripada versi yang dipahaminya. |
| `version` | string | **ya** | SemVer 2.0.0, dengan atau tanpa awalan `v`. |
| `channel` | string | tidak | Bersifat informatif. Klien yang dikonfigurasi untuk kanal lain memperlakukan manifes tersebut seolah-olah tidak ada pembaruan. |
| `name` | string | tidak | Judul rilis yang mudah dibaca manusia, ditampilkan di jendela pembaruan. |
| `notes` | string | tidak | Catatan rilis dalam Markdown, dirender di jendela pembaruan. |
| `publishedAt` | string | tidak | Stempel waktu RFC 3339. |
| `artifacts` | array | **ya** | Satu entri untuk setiap artefak yang dapat diunduh. Urutan menunjukkan preferensi penerbit. |
| `metadata` | object | tidak | Data kunci/nilai berformat bebas yang diteruskan ke aplikasi. |

Klien mengabaikan bidang yang tidak dikenal, sehingga server dapat menambahkan bidangnya sendiri tanpa merusak kompatibilitas. Penambahan khusus server ditempatkan di `metadata`.

### Bidang artefak

| Bidang | Tipe | Wajib | Catatan |
| --- | --- | --- | --- |
| `url` | string | **ya** | Absolut atau relatif terhadap URL manifes. Hanya `http(s)`. |
| `platform` | string | tidak | Nilai Go `GOOS`. Alias umum (`macos`, `win`, ...) dapat digunakan. Nilai kosong cocok dengan semua platform. |
| `arch` | string | tidak | Nilai Go `GOARCH`. Alias umum (`x86_64`, `aarch64`, ...) dapat digunakan. Nilai kosong cocok dengan semua arsitektur. |
| `filename` | string | tidak | Secara default menggunakan segmen jalur terakhir dari `url`. |
| `filetype` | string | tidak | Secara default menggunakan ekstensi nama file. |
| `size` | int | tidak | Byte, digunakan untuk menampilkan progres pengunduhan. |
| `digestAlgo` / `digest` | string / base64 | tidak | `sha256` atau `sha512`. |
| `signatureAlgo` / `signature` | string / base64 | tidak | `ed25519`, `ed25519ph`, atau `ecdsa-p256`. `signatureAlgo` wajib ada setiap kali `signature` disertakan. Lihat [panduan updater](/guides/updater/#cryptographic-verification) untuk mengetahui apa yang ditandatangani oleh setiap algoritma. |

Klien memilih artefak **pertama** yang `platform` dan `arch`-nya cocok dengan sistem yang sedang berjalan. Nilai Base64 dapat menggunakan padding maupun tidak.

### Perbandingan versi

Penentuan apakah manifes merupakan pemutakhiran selalu dilakukan di sisi klien berdasarkan presedensi SemVer 2.0.0: `version` manifes harus benar-benar lebih baru daripada versi yang terinstal. Dengan demikian, hosting statis dapat berfungsi dengan benar tanpa pengaturan khusus (manifes selalu mendeskripsikan rilis terbaru dan klien yang sudah mutakhir tidak melakukan apa pun), sedangkan `204` tetap tersedia bagi server dinamis untuk menghemat bandwidth.

## Verifikasi dan kepercayaan

Checksum dan tanda tangan disertakan dalam manifes, tetapi akar kepercayaan tidak: tanda tangan diverifikasi terhadap kunci publik yang disematkan oleh aplikasi melalui `updater.Config.PublicKey` pada waktu build. Sumber pemutakhiran yang disusupi atau diganti tidak dapat menyediakan kuncinya sendiri. Artefak yang menyertakan tanda tangan ketika aplikasi tidak memiliki kunci yang disematkan akan ditolak secara aman; begitu pula tanda tangan tanpa `signatureAlgo` yang dideklarasikan atau yang tidak dapat didekode. Klien tidak pernah diam-diam beralih ke verifikasi berbasis digest saja.

Artefak yang hanya menggunakan digest akan diinstal setelah pemeriksaan digest. Pemeriksaan ini melindungi dari kerusakan data, tetapi mengandalkan TLS beserta integritas host itu sendiri agar tahan terhadap manipulasi. Sertakan tanda tangan untuk semua hal yang sensitif terhadap keamanan.

Penandatanganan artefak dengan skema `ed25519ph` milik framework hanya memerlukan beberapa baris kode Go:

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

Dalam praktiknya, Anda jarang perlu menulis kode tersebut: CLI melakukannya untuk Anda.

## Publikasi dengan CLI wails3

Grup perintah `wails3 updater` mencakup seluruh alur publikasi. Sebuah rilis hanya memerlukan tiga perintah:

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest` menerima file atau direktori (materi kunci, `.json`, serta file pendamping checksum dan catatan dilewati secara otomatis), mengalirkan setiap artefak melalui SHA-512, menandatangani digest dengan Ed25519ph jika `-key` diberikan, serta menyimpulkan `platform` dan `arch` dari nama file konvensional seperti `MyApp-2.1.0-darwin-arm64.zip` (alias umum seperti `macOS`, `win64`, `x86_64`, dan `aarch64` dikenali; peringatan ditampilkan untuk nilai yang tidak dapat disimpulkan, yang kemudian dianggap cocok dengan semua platform). Hilangkan `-url-prefix` untuk menghasilkan URL relatif dan unggah manifes di lokasi yang sama dengan artefak.

`verify` keluar dengan kode bukan nol jika ada ketidakcocokan, sehingga cocok digunakan sebagai gerbang CI antara proses build dan publikasi. Untuk server yang menyusun manifes sendiri, `wails3 updater sign -key updater.key <files...>` mencetak bidang `digest`/`signature` setiap file sebagai JSON yang siap digabungkan ke dalam dokumen Anda sendiri.

## Autentikasi

Autentikasi merupakan urusan server; protokol hanya membawa header. Klien mengirim ulang header yang dikonfigurasi pada setiap permintaan manifes. Untuk pengunduhan artefak, header `Authorization` hanya dikirim jika URL artefak berada di host yang sama dengan manifes dan tidak diturunkan dari `https` ke `http`. Header tersebut juga dihapus pada setiap pengalihan lintas origin atau pengalihan yang menurunkan protokol, sehingga kredensial tidak pernah bocor ke CDN atau penyimpanan objek dan tidak pernah dikirim sebagai teks biasa.

Contoh yang dibatasi lisensi dan cocok dipadukan dengan layanan lisensi terkelola:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## Konfigurasi klien

Lihat [panduan updater](/guides/updater/#providers) untuk referensi lengkap `endpoint.Config` dan cara menempatkan penyedia tersebut dalam rantai fallback bersama penyedia GitHub, keygen.sh, dan AppCast.
