---
title: "Konfigurasi UAC Windows"
description: "Konfigurasikan Kontrol Akun Pengguna (UAC) untuk aplikasi Wails Windows Anda"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

Platform yang Relevan: <span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Kontrol Akun Pengguna (UAC) Windows menentukan hak istimewa eksekusi aplikasi Wails Anda. Secara default, aplikasi Wails v3 menyertakan konfigurasi UAC eksplisit dalam manifes Windows-nya, sehingga memastikan perilaku yang konsisten di berbagai komputer.

## Tingkat Eksekusi UAC

Aplikasi Windows dapat meminta tingkat eksekusi yang berbeda melalui file manifesnya. Wails v3 secara otomatis menyertakan konfigurasi UAC dengan tingkat eksekusi default yang dapat Anda sesuaikan berdasarkan kebutuhan aplikasi Anda.

### Tingkat Eksekusi yang Tersedia

| Tingkat | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `asInvoker` | Berjalan dengan hak istimewa yang sama seperti proses induk | Default untuk sebagian besar aplikasi |
| `highestAvailable` | Berjalan dengan hak istimewa tertinggi yang tersedia bagi pengguna | Aplikasi yang mungkin memerlukan akses dengan hak istimewa lebih tinggi |
| `requireAdministrator` | Selalu memerlukan hak istimewa administrator | Utilitas sistem, penginstal |

### Konfigurasi Default

Aplikasi Wails v3 menyertakan konfigurasi UAC default dalam manifes Windows-nya:

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

Konfigurasi ini memastikan aplikasi Anda:

- Berjalan dengan hak istimewa yang sama seperti proses yang menjalankannya
- Secara default tidak memerlukan peningkatan hak istimewa
- Berfungsi secara konsisten di berbagai komputer
- Tidak memicu dialog UAC bagi pengguna biasa

## Menyesuaikan Konfigurasi UAC

Karena Wails v3 mendorong pengguna untuk menyesuaikan aset build mereka, Anda dapat mengubah konfigurasi UAC dengan mengedit templat manifes Windows secara langsung.

### Menemukan Templat Manifes

Templat manifes Windows berada di:

```
build/windows/wails.exe.manifest
```

### Mengubah Tingkat Eksekusi

Untuk mengubah tingkat eksekusi, edit atribut `level` dalam elemen `requestedExecutionLevel`:

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### Contoh

#### Aplikasi Standar (Default)

Sebagian besar aplikasi sebaiknya menggunakan tingkat `asInvoker` default:

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### Utilitas Sistem

Aplikasi yang memerlukan akses dengan hak istimewa lebih tinggi jika tersedia:

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### Alat Administratif

Aplikasi yang selalu memerlukan hak istimewa administrator:

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## Akses UI

Atribut `uiAccess` mengontrol apakah aplikasi Anda dapat berinteraksi dengan elemen UI yang memiliki hak istimewa lebih tinggi. Dalam sebagian besar kasus, nilainya sebaiknya tetap `false`.

Atur ke `true` hanya jika aplikasi Anda perlu:

- Mengirim input ke aplikasi lain
- Mengendalikan UI aplikasi lain
- Mengakses elemen UI milik proses dengan hak istimewa lebih tinggi

@note{type="caution" title="Persyaratan Akses UI"}
Jika diatur ke `uiAccess="true"`, aplikasi Anda harus:

- Ditandatangani secara digital dengan sertifikat dari otoritas sertifikat tepercaya
- Diinstal di lokasi aman (Program Files atau Windows\System32)

@end

## Mem-build dengan Pengaturan UAC Khusus

Setelah mengubah templat manifes, build aplikasi Anda seperti biasa:

```bash
wails3 build
```

Proses build akan secara otomatis menyematkan konfigurasi UAC khusus Anda ke dalam file yang dapat dieksekusi.

## Memverifikasi Konfigurasi UAC

Anda dapat memverifikasi bahwa pengaturan UAC telah disematkan dengan benar menggunakan alat `go-winres`:

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

Kemudian, periksa file manifes yang diekstrak untuk memastikan konfigurasi UAC Anda tercantum di dalamnya.

@note{type="tip" title="Persistensi Manifes"}
Tidak seperti beberapa framework lain, konfigurasi UAC Wails v3 disematkan langsung ke dalam file yang dapat dieksekusi saat kompilasi, sehingga konfigurasi tersebut tetap ada ketika aplikasi disalin ke komputer lain.

@end

## Pemecahan Masalah

### Dialog UAC Tidak Muncul

Jika Anda mengatur `requireAdministrator` tetapi dialog UAC tidak muncul:

- Pastikan manifes disematkan dengan benar ke dalam file yang dapat dieksekusi
- Pastikan Anda tidak menjalankannya dari proses yang hak istimewanya sudah ditingkatkan
- Pastikan sintaks manifes merupakan XML yang valid

### Aplikasi Tidak Dapat Dimulai

Jika aplikasi Anda gagal dimulai setelah perubahan UAC:

- Periksa apakah sintaks manifes mengandung kesalahan XML
- Pastikan nilai tingkat eksekusi valid
- Coba kembalikan ke `asInvoker` untuk mengisolasi masalah

### Perilaku yang Tidak Konsisten di Berbagai Komputer

Jika perilaku UAC berbeda di antara komputer:

- Pastikan manifes disematkan dalam berkas yang dapat dieksekusi (bukan sebagai berkas eksternal)
- Pastikan berkas yang dapat dieksekusi tidak dimodifikasi setelah proses build
- Pastikan pengaturan UAC Windows diaktifkan pada komputer target
