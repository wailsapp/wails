---
title: "Pengujian \u0026 Integrasi Berkelanjutan"
description: "Cara Wails v3 memastikan kualitas melalui pengujian unit, rangkaian pengujian integrasi, deteksi race, dan CI GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Framework desktop yang tangguh memerlukan pengujian yang benar-benar andal. Wails v3 menggunakan **strategi berlapis**:

| Lapisan | Tujuan | Peralatan |
| --- | --- | --- |
| Pengujian unit | Umpan balik cepat untuk fungsi yang terisolasi | `go test ./...` |
| Pengujian generator/CLI | Memvalidasi `wails3 generate bindings` dan infrastruktur CLI | `task test:generator`, `task test:cli` |
| Pengujian templat | Memastikan setiap templat yang didistribusikan tetap dapat dibangun | `task test:templates` |
| Deteksi race | Mendeteksi data race pada runtime dan bridge | `go test -race ./...` |
| Matriks CI | Keyakinan lintas OS pada setiap PR | GitHub Actions |

Dokumen ini menjelaskan **lokasi pengujian**, **cara menjalankannya**, dan **hal-hal yang diorkestrasi oleh Taskfile**.

`v3/TESTING.md` mendokumentasikan build contoh lintas platform; ini bukan katalog penekanan laporan race.

---

## 1. Konvensi Direktori

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

Panduan:

- **Letakkan pengujian unit di sebelah kode** (`foo.go` ↔ `foo_test.go`).
- Gunakan **gaya black-box** untuk paket `pkg/` (`package application_test`) jika hal tersebut meningkatkan kebersihan API.
- Fixture bersama ditempatkan di lokasi penggunaannya (tidak ada paket `internal/testutil/` terpusat dalam pohon ini—hubungkan helper secara terpisah untuk setiap paket).

---

## 2. Pengujian Unit

### Menulis Pengujian

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Rekomendasi:

- Gunakan [`stretchr/testify`](https://github.com/stretchr/testify)—sudah tersedia di `go.mod`.
- Utamakan pengujian **berbasis tabel** setiap kali beberapa input, kasus tepi, atau hasil yang diharapkan menguji perilaku yang sama. Berikan nama deskriptif untuk setiap kasus.
- Jika diperlukan, buat stub untuk kekhususan platform di balik tag build (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Target cakupan

Logika baru dan yang diubah diharapkan memiliki cakupan pernyataan Go sebesar 100%. Ukur paket yang Anda ubah, bukan mengandalkan persentase seluruh repositori:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Beberapa jalur memang tidak dapat diuji secara wajar dalam lingkungan pengujian normal—misalnya kegagalan khusus platform, perilaku yang bergantung pada perangkat keras, atau fallback defensif yang tidak dapat dipicu dengan aman. Batasi pengecualian tersebut dan jelaskan setiap jalur yang tidak tercakup dalam deskripsi PR.

### Menjalankan Secara Lokal

```bash
cd v3
go test ./... -cover
```

Atau melalui Taskfile (target sebenarnya—tidak ada pintasan `task test`):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. Pengujian Integrasi

`v3/tests/` menampung harness integrasi lintas paket. Target Taskfile `test:example:*` dan `test:examples:*` menjalankan pemeriksaan kompilasi pada darwin/windows/linux (termasuk matriks GTK3/GTK4 berbasis Docker di Linux).

> Contoh yang dapat dijalankan tersedia di `v3/examples/`. Target pengujian memilih dan membangun
>
> contoh yang sesuai untuk platform host atau matriks CI.

Jalankan rangkaian smoke test untuk platform host dengan:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Deteksi Race

Data race berakibat fatal pada runtime GUI.

### Menyelidiki race

Periksa kedua akses yang berkonflik dan reproduksi masalah pada platform yang terdampak. Jangan menganggap laporan sebagai sesuatu yang tidak berbahaya atau menekannya tanpa bukti. `v3/TESTING.md` menjelaskan kompilasi contoh; keberhasilan build tersebut tidak membuktikan ketiadaan race ataupun perilaku runtime native.

### Rangkaian Pengujian Race Lokal

```
go test -race ./...
```

> `wails3 dev` tidak memiliki flag `-race`—flag CLI-nya adalah `--config`, `--port`,
>
> dan `-s` (mengaktifkan HTTPS). Untuk menjalankan runtime di bawah race detector,
>
> bangun aplikasi pengujian dengan `go build -race` lalu jalankan secara langsung.

---

## 5. Alur Kerja GitHub Actions

File alur kerja yang sebenarnya di bawah `.github/workflows/` (telah diverifikasi terhadap pohon repositori):

| File | Tujuan |
| --- | --- |
| `build-and-test-v3.yml` | Matriks Go 1.26/1.27, pemeriksaan runtime dan generator, kompilasi templat/contoh, serta pemeriksaan build aplikasi. Secara default, Linux menggunakan GTK4/WebKitGTK 6.0; varian lama menetapkan `BUILD_TAGS=gtk3`. Peluncuran dan interaksi native harus diperiksa secara terpisah. |
| `cross-compile-test-v3.yml` | Pemeriksaan kewajaran kompilasi silang |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Otomatisasi changelog |
| `nightly-release-v3.yml` | Artefak rilis nightly v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Pengelolaan dependensi/rilis WebView2 |
| `build-cross-image.yml` | Membangun image container untuk kompilator silang |
| `publish-npm.yml` | Menerbitkan runtime JS `@wailsio/runtime` tertanam ke npm |
| `pr-master.yml` | Pemeriksaan pada sisi PR terhadap branch `master` |
| `semgrep.yml` | Analisis statis Semgrep |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Pemeliharaan repositori dan alur kerja v2; alur kerja sinkronisasi/pengunggahan Crowdin yang lama telah dihentikan. |

Langkah-langkah CI mencerminkan target Taskfile di atas (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), sehingga Anda dapat mereproduksi CI secara lokal dengan kesetaraan satu banding satu. Langkah smoke `wails3 build` di `build-and-test-v3.yml` dipanggil dengan **tanpa flag tambahan**—`wails3 build` tidak memiliki flag `-skip-package`.

---

## 6. Kesetaraan CI Lokal

Tidak ada satu target payung `task ci`. Reproduksi CI dengan merangkai target yang sebenarnya:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Memecahkan Masalah Kegagalan Pengujian

| Gejala | Kemungkinan Penyebab | Perbaikan |
| --- | --- | --- |
| **Race pada `webview_window_darwin.go`** | Mengubah status jendela di luar thread utama | Teruskan melalui `application.InvokeAsync` / `Invoke` agar pemanggilan berjalan di thread utama |
| **Pengujian Linux macet di CI tanpa antarmuka grafis** | GTK memerlukan display | Jalankan dengan `xvfb-run`, misalnya `xvfb-run task test:examples:linux` |
| **Build templat gagal** | Lockfile frontend sudah tidak mutakhir | Jalankan ulang `wails3 init` pada direktori bersih untuk memperbarui templat |
| **Kesalahan coverpkg** | Pengujian integrasi mengimpor `main` | Beralihlah ke tag build `//go:build integration` dan batasi impor tersebut dengannya |

---

## 8. Menambahkan Pengujian Baru

1. **Unit** — buat `*_test.go`, lalu jalankan `go test ./...`
2. **Generator / CLI** — perluas kasus di bawah `internal/generator/testcases/` atau `internal/commands/*_test.go`, lalu jalankan ulang `task test:generator` / `task test:cli`
3. **Templat / contoh** — pastikan templat yang didistribusikan tetap dapat di-build dengan `task test:templates`

---

## 9. Peta Berkas Utama

| Apa | Path |
| --- | --- |
| Pengujian round-trip generator | `internal/generator/generate_test.go` |
| Pengujian build-assets | `internal/commands/build-assets_test.go` |
| Build contoh lintas platform | `v3/TESTING.md` |
| Target pengujian Taskfile | `v3/Taskfile.yaml` |
| Generator konstanta event | `v3/tasks/events/generate.go` |
| Alur kerja CI | `.github/workflows/build-and-test-v3.yml` (1.26/1.27 Go melalui `actions/setup-go@v5`) |
| Analisis statis | `.github/workflows/semgrep.yml` |
| Publikasi runtime npm | `.github/workflows/publish-npm.yml` |

---

Kualitas bukan sekadar pertimbangan belakangan di Wails v3. Dengan pengujian unit, rangkaian pengujian generator/templat, deteksi race, dan matriks CI lintas platform, Anda dapat berkontribusi dengan yakin karena perubahan Anda berhasil dijalankan pada setiap sistem operasi yang kami dukung. Selamat menguji!
