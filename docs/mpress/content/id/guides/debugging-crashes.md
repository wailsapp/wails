---
title: "Men-debug Crash"
description: "Mengumpulkan laporan diagnostik dan minidump dari aplikasi Wails"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

Saat terjadi masalah di produksi, log saja jarang cukup. Panduan ini membahas paket `github.com/wailsapp/wails/v3/pkg/debug`, yang menghasilkan laporan diagnostik terstruktur dan — di Windows — minidump yang dapat dibuka di WinDbg atau Visual Studio.

Paket ini sengaja dibuat kecil dan memiliki dua titik masuk:

- `debug.Dump(...)` — menulis minidump Windows dari proses saat ini.
- `debug.Report(...)` — mengumpulkan snapshot diagnostik terperinci, dengan minidump opsional.

Keduanya berjalan menggunakan token pengguna pemanggil. Tanpa peningkatan hak akses, tanpa `SeDebugPrivilege`, tanpa `OpenProcess`: pembuatan dump proses saat ini menggunakan pseudo-handle `GetCurrentProcess`, yang melewati pemeriksaan ACL untuk akses ke proses sendiri.

## Referensi singkat

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` selalu mengembalikan `*CrashReport`, bahkan saat sebagian pengumpulan gagal — pemanggil tetap mendapatkan informasi yang berhasil dikumpulkan sebelum kesalahan, sehingga laporan tidak pernah `nil` ketika `err != nil`.

## Struktur `CrashReport`

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

Serialisasikan dengan `encoding/json` untuk layanan pelaporan crash, atau baca field secara langsung jika Anda hanya memerlukan sebagian informasi.

## Integrasi dengan `PanicHandler`

Pola paling berguna adalah mengumpulkan laporan dan minidump opsional ketika Wails menangkap panic, lalu menyerahkan data gabungan tersebut ke alur pelaporan Anda sendiri:

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails **tidak** memanggil `debug.Report` secara otomatis — pelaporan crash sering memerlukan persetujuan pengguna atau penyamaran informasi identitas pribadi, sehingga keputusan tersebut menjadi tanggung jawab handler Anda. Lihat [panduan penanganan panic](/guides/panic-handling/) untuk mengetahui kapan Wails memanggil `PanicHandler` dan cara menangani goroutine buatan pengguna.

## Debugging manual

`debug.Dump` dan `debug.Report` juga berguna di luar kejadian panic:

- Deteksi hang: ketika watchdog terpicu, buat dump proses agar Anda dapat membukanya di WinDbg dan melihat apa yang memblokir setiap thread.
- Laporan keadaan abnormal: misalnya jumlah goroutine terus meningkat tanpa batas — kumpulkan laporan, biarkan aplikasi tetap berjalan, lalu selidiki nanti.
- Paket dukungan sesuai permintaan: hubungkan item menu "Laporkan bug" ke `debug.Report(debug.WithDump())` dan lampirkan hasilnya pada tiket dukungan pengguna.

## Penggunaan minidump Windows

Secara default, minidump menggunakan kumpulan flag yang kaya informasi tetapi ringkas:

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

Ini menghasilkan dump berukuran sekitar 5–50 MB yang mempertahankan semua stack thread, tabel handle, dan daftar modul — cukup untuk merekonstruksi posisi setiap goroutine saat dump dibuat.

`debug.WithFullMemory()` atau `debug.WithDumpFullMemory()` pada `Report` menambahkan `MiniDumpWithFullMemory`, yang menangkap seluruh ruang alamat proses. Berguna untuk menyelidiki mengapa pointer tertentu rusak, tetapi ukuran file dapat mencapai ratusan MB atau GB. Gunakan hanya untuk penyelidikan mendalam.

Buka `.dmp` yang dihasilkan di WinDbg dengan `windbg -z C:\path\to\your.dmp` atau File → Open Crash Dump di Visual Studio. Simbol biasanya dihapus dari build rilis — pasangkan dump dengan `.pdb` dari build yang sama, atau biner Go itu sendiri jika dibangun dengan informasi debug, untuk mendapatkan stack trace yang berguna.

## Platform selain Windows

`debug.Dump` mengembalikan kesalahan "belum diimplementasikan" pada platform selain Windows. Linux dan macOS memiliki perangkat analisis pascakegagalan yang berbeda, seperti core dump melalui `GOTRACEBACK=crash`, `rr`, atau pelapor crash khusus platform, yang tidak dapat dipetakan langsung ke API `MiniDumpWriteDump`. `debug.Report` tetap berfungsi — Anda mendapatkan semuanya kecuali file dump.

Kesalahan dari `Dump` dapat diperiksa agar aplikasi dapat melanjutkan dengan fitur yang tersedia:

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## Catatan keamanan

- **Tidak memerlukan peningkatan hak akses.** Paket ini sengaja menghindari `SeDebugPrivilege` dan jalur `OpenProcess` / PID jarak jauh yang digunakan alat pencuri kredensial. Paket hanya membuat dump proses pemanggil.
- **File dump berisi semua isi memori.** Termasuk kredensial yang sedang digunakan, token autentikasi, data pengguna, dan kunci enkripsi. Perlakukan dump seperti citra memori lengkap: lindungi saat dikirim, bersihkan data sensitif yang tersimpan, dan pertimbangkan penyamaran sebelum mengunggah ke pihak ketiga.
- **Lokasi default adalah `os.TempDir()`.** Di Windows, lokasi ini mengarah ke `%LOCALAPPDATA%\Temp` — khusus pengguna, tidak dapat dibaca semua orang, tetapi file tetap tersimpan. Pindahkan dump ke lokasi aman jika ingin menyimpannya.
