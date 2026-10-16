---
title: "Pengujian \u0026 Integrasi Berkelanjutan"
description: "Cara Wails v3 memastikan kualitas melalui pengujian unit, rangkaian pengujian integrasi, deteksi race, dan CI GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Framework desktop yang tangguh memerlukan pengujian yang sangat andal. Wails v3 menerapkan **strategi berlapis**:

| Lapisan | Tujuan | Peralatan |
| --- | --- | --- |
| Pengujian unit | Umpan balik cepat untuk fungsi yang terisolasi | `go test ./...` |
| Pengujian generator/CLI | Memvalidasi `wails3 generate bindings` dan mekanisme penghubung CLI | `task test:generator`, `task test:cli` |
| Pengujian templat | Memastikan setiap templat yang disertakan tetap dapat dibangun | `task test:templates` |
| Deteksi race | Mendeteksi data race dalam runtime & bridge | `go test -race ./...` |
| Matriks CI | Keyakinan akan keandalan lintas sistem operasi pada setiap PR | GitHub Actions |

Dokumen ini menjelaskan **lokasi pengujian**, **cara menjalankannya**, dan **hal-hal yang diorkestrasi oleh Taskfile**.

> Panduan race yang dirujuk oleh draf lama sebagai `pkg/application/RACE.md`
>
> kini berada di `v3/TESTING.md`.

---

## 1. Konvensi Direktori

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

Panduan:

- **Letakkan pengujian unit di sebelah kode** (`foo.go` ↔ `foo_test.go`).
- Gunakan **gaya black-box** untuk paket `pkg/` (`package application_test`) jika hal itu meningkatkan kebersihan API.
- Fixture bersama berada di tempat penggunaannya (tidak ada paket `internal/testutil/` terpusat dalam working tree ini—hubungkan helper untuk setiap paket sebagai gantinya).

---

## 2. Pengujian Unit

### Menulis Pengujian

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Rekomendasi:

- Gunakan [`stretchr/testify`](https://github.com/stretchr/testify)—sudah tersedia dalam `go.mod`.
- Utamakan pengujian **berbasis tabel** setiap kali beberapa input, kasus batas, atau hasil yang diharapkan menguji perilaku yang sama. Beri setiap kasus nama yang deskriptif.
- Jika diperlukan, buat stub untuk kekhususan platform di balik build tag (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Target cakupan

Logika baru dan yang diubah diharapkan memiliki cakupan pernyataan Go sebesar 100%. Ukur paket yang Anda ubah, bukan mengandalkan persentase seluruh repositori:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Beberapa jalur memang tidak dapat diuji secara wajar dalam lingkungan pengujian normal—misalnya kegagalan yang hanya terjadi pada platform tertentu, perilaku yang bergantung pada perangkat keras, atau fallback defensif yang tidak dapat dipicu dengan aman. Batasi pengecualian tersebut secara ketat dan jelaskan setiap jalur yang tidak tercakup dalam deskripsi PR.

### Menjalankan Secara Lokal

```bash
cd v3
go test ./... -cover
```

Atau melalui Taskfile (target yang sebenarnya—tidak ada pintasan `task test`):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. Pengujian Integrasi

`v3/tests/` menampung harness integrasi lintas paket. Target Taskfile `test:example:*` dan `test:examples:*` menjalankan pemeriksaan build dan peluncuran di darwin/windows/linux (termasuk matriks GTK3/GTK4 berbasis Docker di Linux).

> Contoh yang dapat dijalankan berada di `v3/examples/`. Target pengujian memilih dan membangun
>
> contoh yang sesuai untuk platform host atau matriks CI.

Jalankan rangkaian smoke test untuk platform host dengan:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Deteksi Race

Data race berakibat fatal dalam runtime GUI.

### Panduan Race

Lihat `v3/TESTING.md` untuk:

- Race jinak yang telah diketahui dan alasan penyisihannya
- Cara menafsirkan stack trace yang melintasi batas Cgo (Linux GTK + WebKit2GTK)

### Rangkaian Pengujian Race Lokal

```
go test -race ./...
```

> `wails3 dev` tidak memiliki flag `-race`—flag CLI-nya adalah `--config`, `--port`,
>
> dan `-s` (mengaktifkan HTTPS). Untuk menguji runtime dengan race detector,
>
> bangun aplikasi pengujian dengan `go build -race` lalu jalankan secara langsung.

---

## 5. Alur Kerja GitHub Actions

File alur kerja aktual di bawah `.github/workflows/` (telah diverifikasi terhadap working tree):

| File | Tujuan |
| --- | --- |
| `build-and-test-v3.yml` | Matriks build + pengujian utama v3. Menggunakan `actions/setup-go@v5` dengan `go-version: 1.25`. Menjalankan `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples` (serta `BUILD_TAGS=gtk4 task test:examples` untuk jalur GTK4), `task generator:test:check`, `task install`, lalu `wails3 build` sebagai smoke check. Job Linux menginstal `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk` dan menjalankan rangkaian pengujian di bawah `dbus-run-session -- xvfb-run`. |
| `cross-compile-test-v3.yml` | Pemeriksaan kewajaran kompilasi silang |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Otomatisasi changelog |
| `nightly-release-v3.yml` | Artefak rilis nightly v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Pengelolaan dependensi / rilis WebView2 |
| `build-cross-image.yml` | Membangun image container cross-compiler |
| `publish-npm.yml` | Menerbitkan runtime JS `@wailsio/runtime` tersemat ke npm |
| `pr-master.yml` | Pemeriksaan pada sisi PR terhadap cabang `master` |
| `semgrep.yml` | Analisis statis Semgrep |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Alur pemeliharaan repositori / sisi v2 |

Di pohon ini, **tidak ada** `qodana.yaml` dan **tidak ada** `runtime.yml` — draf lama halaman ini menyebut keduanya, tetapi hanya `semgrep.yml` yang mencakup analisis statis, dan paket JS runtime diterbitkan melalui `publish-npm.yml`.

Langkah-langkah CI mencerminkan target Taskfile di atas (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), sehingga Anda dapat mereproduksi setiap langkah CI secara lokal dengan padanan yang persis sama. Langkah smoke `wails3 build` di `build-and-test-v3.yml` dijalankan **tanpa flag tambahan** — tidak ada flag `-skip-package` pada `wails3 build`.

---

## 6. Kesetaraan Lokal dengan CI

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
| **Race pada `webview_window_darwin.go`** | Mengubah state jendela di luar thread utama | Lakukan marshaling melalui `application.InvokeAsync` / `Invoke` agar panggilan berjalan di thread utama |
| **Pengujian Linux macet di CI headless** | GTK memerlukan display | Jalankan di bawah `xvfb-run`, misalnya `xvfb-run task test:examples:linux` |
| **Build templat gagal** | Lockfile frontend sudah tidak mutakhir | Jalankan kembali `wails3 init` pada direktori bersih untuk memperbarui templat |
| **Kesalahan Coverpkg** | Pengujian integrasi mengimpor `main` | Beralihlah ke build tag `//go:build integration` dan batasi impor dengan kondisi tersebut |

---

## 8. Menambahkan Pengujian Baru

1. **Unit** — buat `*_test.go`, lalu jalankan `go test ./...`
2. **Generator / CLI** — perluas kasus di bawah `internal/generator/testcases/` atau `internal/commands/*_test.go`, lalu jalankan kembali `task test:generator` / `task test:cli`
3. **Templat / contoh** — pastikan templat yang didistribusikan tetap dapat dibangun dengan `task test:templates`

---

## 9. Peta File Utama

| Apa | Path |
| --- | --- |
| Pengujian round-trip generator | `internal/generator/generate_test.go` |
| Pengujian build-assets | `internal/commands/build-assets_test.go` |
| Panduan race / Cgo | `v3/TESTING.md` |
| Target pengujian Taskfile | `v3/Taskfile.yaml` |
| Generator konstanta event | `v3/tasks/events/generate.go` |
| Workflow CI | `.github/workflows/build-and-test-v3.yml` (Go 1.25 melalui `actions/setup-go@v5`) |
| Analisis statis | `.github/workflows/semgrep.yml` |
| Publikasi runtime ke npm | `.github/workflows/publish-npm.yml` |

---

Kualitas bukan sekadar pertimbangan belakangan di Wails v3. Dengan pengujian unit, rangkaian pengujian generator/templat, deteksi race, dan matriks CI lintas platform, Anda dapat berkontribusi dengan yakin karena mengetahui bahwa perubahan Anda berhasil dijalankan di setiap OS yang kami dukung. Selamat menguji!
