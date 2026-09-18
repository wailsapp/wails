---
title: "Debugging"
description: "Menyelidiki masalah dan membuat profil performa aplikasi"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

Panduan ini menunjukkan berbagai alat untuk memeriksa dan menyelidiki kemungkinan masalah performa dalam aplikasi Wails Anda menggunakan

- [`runtime/trace`](https://pkg.go.dev/runtime/trace) untuk membuat grafik performa dan memeriksanya di browser

## Membuat trace performa

@steps
### Siapkan aplikasi Anda untuk tracing
Di dekat titik masuk program Anda, pastikan terdapat kode seperti berikut

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

Output akan disimpan ke `trace.out` dalam direktori kerja Anda, dengan metrik yang mencakup seluruh waktu berjalan aplikasi. Trace bawaan cukup mentah; Anda sangat disarankan untuk mempelajari trace lebih lanjut agar dapat menambahkan lebih banyak informasi kontekstual dan membatasi data yang benar-benar direkam. Contoh: `WithRegion, NewTask, Log`

### Jalankan aplikasi Anda untuk membuat trace
Saat aplikasi berjalan, trace akan direkam terus-menerus hingga aplikasi ditutup. Lakukan beberapa tindakan dalam aplikasi Anda, lalu tutup aplikasi setelah selesai

### Instal alat bantu visualisasi
Beberapa tampilan memerlukan [Graphviz](https://graphviz.org/). Anda dapat memverifikasi bahwa Graphviz telah terinstal dengan menjalankan `dot -V`

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### Visualisasikan trace Anda
Setelah data trace tersedia, kita dapat memulai antarmuka web

```bash
go tool trace trace.out
```

Tindakan ini seharusnya membuka browser bawaan Anda (disarankan menggunakan browser berbasis Chrome) dan menampilkan halaman beranda penampil peristiwa trace.

Bagi pengguna baru, layar `Syscall profile` mungkin paling berguna. Layar ini memberikan perincian tentang apa yang dilakukan program secara tepat beserta informasi waktunya

@end
