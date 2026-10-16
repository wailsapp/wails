---
title: "Arsitektur Wails v3"
description: "Diagram dan penjelasan mendalam tentang setiap komponen yang bekerja di dalam Wails v3"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 adalah **framework desktop full-stack** yang terdiri atas runtime Go, bridge JavaScript, toolchain berbasis tugas, dan sekumpulan templat yang memungkinkan Anda mendistribusikan aplikasi native yang didukung teknologi web modern.

Halaman ini menyajikan *gambaran umum* dalam empat diagram:

1. **Arsitektur Keseluruhan** – cara setiap subsistem saling terhubung\
2. **Alur Runtime** – apa yang terjadi saat JS memanggil Go dan sebaliknya\
3. **Pengembangan vs Produksi** – dua mode server aset\
4. **Implementasi Platform** – lokasi kode khusus OS\

---

## 1 · Arsitektur Keseluruhan

**Wails v3 – Stack Tingkat Tinggi**

**[Placeholder Diagram Stack Tingkat Tinggi]**

---

## 2 · Alur Pemanggilan Runtime

**Runtime – Jalur Pemanggilan JavaScript ⇄ Go**

**[Placeholder Diagram Alur Pemanggilan Runtime]**

Poin-poin utama:

- **Tanpa HTTP / IPC** – bridge menggunakan kanal dalam memori milik WebView native\
- **ID metode** – hash FNV deterministik memungkinkan pencarian O(1) di Go\
- **Promise** – galat diteruskan sebagai rejection beserta stack & kode

---

## 3 · Alur Aset Pengembangan vs Produksi

**Server Aset Dev ↔ Prod**

**[Placeholder Diagram Alur Aset]**

- Dalam mode **dev**, server meneruskan path yang tidak dikenal melalui proksi ke server live-reload milik framework dan menyajikan aset statis dari disk.
- Dalam mode **prod**, API yang sama didukung oleh `go:embed` sehingga menghasilkan biner tanpa dependensi.

---

## 4 · Pemisahan Runtime Khusus Platform

**File Runtime per OS**

**[Placeholder Diagram Pemisahan Platform]**

Setiap fitur mengikuti pola ini:

1. **Antarmuka umum** di `pkg/application`\
2. Entri **pemroses pesan** di `pkg/application/messageprocessor_*.go`\
3. **Implementasi per OS** di `pkg/application/*_{darwin,linux,windows}.go` (misalnya `webview_window_darwin.go`, `clipboard_linux.go`, `dialogs_windows.go`, `systemtray_*.go`, `mainthread_*.go`), yang dilindungi oleh build tag. Linux juga memiliki bridge cgo di `linux_cgo.go` / `linux_cgo_gtk4.{go,c,h}`.

`internal/runtime/` hanya memuat sedikit kode penghubung build tag `runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` serta runtime JS tersemat di bawah `internal/runtime/desktop/`.

`internal/capabilities/` tersedia untuk mendeklarasikan kumpulan kapabilitas per platform, tetapi tidak ada sentinel `ErrCapability` — pembatasan fitur dilakukan melalui build tag biasa dan nilai balik stub khusus platform (misalnya `nil` atau galat khusus fitur).

---

## Ringkasan

Diagram-diagram ini menguraikan **lokasi kode**, **cara data berpindah**, dan **lapisan yang bertanggung jawab atas setiap fungsi**. Gunakan diagram ini sebagai acuan saat menjelajahi halaman terperinci berikutnya – diagram ini merupakan peta Anda untuk menelusuri pohon sumber Wails v3.
