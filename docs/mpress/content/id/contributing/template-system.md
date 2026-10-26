---
title: "Sistem Templat"
description: "Cara Wails v3 membuat kerangka proyek baru, cara templat ditata, dan cara membuat templat Anda sendiri."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails menyediakan **sistem templat** yang memungkinkan `wails3 init` menghasilkan proyek yang siap dijalankan. Sejumlah kecil framework sengaja disediakan sebagai templat bawaan (Vanilla, React, Vue, Svelte); framework lain apa pun dapat digunakan dengan [membawa frontend Anda sendiri](/guides/dev/frontend-frameworks/) atau menerbitkan [templat khusus](/guides/advanced/custom-templates/).

Halaman ini membahas:

1. Tata letak direktori templat
2. Cara CLI memilih dan merender templat
3. Membuat templat baru langkah demi langkah
4. Memperbarui atau mengganti templat yang ada
5. Pemecahan masalah dan praktik terbaik

---

## 1. Lokasi Templat

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — boilerplate universal (Taskfile, direktori `build/`, dan infrastruktur bersama) yang digabungkan ke setiap proyek.
- **`base/`** — kerangka frontend minimal untuk templat khusus, bukan lapisan backend otomatis. Pembuatan salinan kerja proyek secara normal menumpangtindihkan `_common/` dengan templat framework yang dipilih.
- **Folder framework** — berisi frontend (`frontend/`), konfigurasi framework, dan `template.yaml` yang menjelaskan metadata templat.
- Nama folder sama dengan **ID templat** yang Anda teruskan ke CLI (`wails3 init -t react`).
- **Konvensi bahasa:** TypeScript merupakan bahasa default dan menggunakan nama tanpa akhiran (`react`); varian JavaScript, jika tersedia, diberi akhiran `-js` (`react-js`). Templat bawaan menyatakan bahasanya secara eksplisit melalui `typescript: true|false` dalam `template.yaml`. Templat komunitas mungkin masih menggunakan akhiran lama `-ts`, yang tetap diterima sebagai fallback.

> Seluruh direktori `internal/templates/` dikompilasi ke dalam biner CLI
>
> melalui `//go:embed *`, sehingga pengguna dapat membuat kerangka proyek secara luring.

---

Pemuat juga menerima metadata lama `template.json`. Templat bawaan baru menggunakan YAML, termasuk pengaturan `typescript` yang eksplisit.

## 2. Cara `wails3 init` Menggunakan Templat

Rantai pemanggilan (tanpa `cmd/wails3/init.go` — CLI dirangkai secara langsung dalam `cmd/wails3/main.go`):

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

Tidak ada API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` — pembuatan salinan kerja ditangani oleh `gosod` (`github.com/leaanthony/gosod`) menggunakan `fs.FS` yang disematkan.

### Flag `wails3 init`

Didefinisikan dalam `internal/flags/init.go`:

| Flag | Kegunaan | Default |
| --- | --- | --- |
| `-p` | Nama paket | `main` |
| `-t` | Nama templat bawaan, jalur lokal, atau URL | `vanilla` |
| `-n` | Nama proyek | (kosong) |
| `-d` | Direktori proyek | `.` |
| `-q` | Sembunyikan keluaran konsol | false |
| `-l` | Tampilkan daftar templat | false |
| `-skipgomodtidy` | Lewati eksekusi `go mod tidy` setelah membuat salinan kerja | false |
| `-git` | URL repositori Git yang akan diinisialisasi | (kosong) |
| `-mod` | Jalur modul Go (diturunkan dari `-git` jika tidak ditetapkan) | (kosong) |
| `-s` | Lewati peringatan saat menggunakan templat jarak jauh | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Metadata yang disematkan ke dalam aset build yang dihasilkan | nilai default yang sesuai |

Tidak ada alias panjang `-list` (hanya `-l`) dan tidak ada `--help` per templat.

### Substitusi

Placeholder merupakan direktif templat Go standar — `.` di awal merupakan bagian dari pengakses bidang:

| Placeholder | Contoh | Sumber |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Flag `-n` / nama direktori |
| `{{.ModulePath}}` | `github.com/me/myapp` | Flag `-mod` atau diturunkan dari `-git` |
| `{{.WailsVersion}}` | `v3.0.0` | Konstanta bawaan hasil kompilasi dari `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Metadata pada saat build | Flag `-product*` yang bersesuaian |

Jika Anda memerlukan placeholder baru, tambahkan bidang ke data templat dalam `internal/templates/templates.go` dan bidang/flag yang sesuai dalam `internal/flags/init.go` (atau tetapkan nilainya dari `internal/commands/init.go`).

### Hook Setelah Penyalinan

Setelah `gosod` selesai membuat salinan kerja dari templat, CLI menjalankan:

```
go mod tidy
```

kecuali jika Anda meneruskan `-skipgomodtidy`. Tidak ada langkah `task deps`.

---

## 3. Membuat Templat Baru

> Contoh: Tambahkan templat **Solid**

### 3.1 Folder dan ID

```
internal/templates/solid/
```

Nama folder = ID templat. Gunakan format **kebab-case**.

### 3.2 Kumpulan File Minimal

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

Mulailah dengan menyalin `react` lalu menghapus file yang tidak diperlukan. Jangan lupa membuat `template.yaml` — tetapkan `typescript: true` untuk templat TypeScript — `base/` adalah satu-satunya folder yang tidak memilikinya.

### 3.3 Perbarui Placeholder

Cari dan ganti nilai contoh literal dengan direktif templat Go, misalnya:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Hubungkan

Karena `templates.go` menelusuri sistem file tersemat saat inisialisasi, umumnya cukup tambahkan folder baru di bawah `internal/templates/<id>/` — tidak diperlukan pemanggilan pendaftaran manual. Jika memerlukan logika tambahan (validasi khusus atau langkah pascapenyalinan), tambahkan logika tersebut ke `templates.Install` dalam `internal/templates/templates.go`.

### 3.5 Uji

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Pastikan:

- Server pengembangan dimulai pada port yang dipublikasikan dalam `WAILS_VITE_PORT`
- Binding yang dihasilkan muncul di bawah `frontend/bindings/...`
- Hot reload berfungsi

---

## 4. Mengubah Templat yang Ada

1. Edit file di bawah `internal/templates/<id>/`.
2. Bangun ulang CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`); direktif `//go:embed *` akan menyertakan konten baru.
3. Naikkan setiap **versi dependensi** dalam `frontend/package.json` dan `Taskfile.yml`.
4. Perbarui deskripsi `template.yaml` milik templat jika perilakunya berubah.

### Penyesuaian Umum

| Tugas | Lokasi |
| --- | --- |
| Ubah port server pengembangan | `frontend/vite.config.ts` — baca `WAILS_VITE_PORT` |
| Tambahkan variabel lingkungan | `build/Taskfile.yml` atau `frontend/.env` |
| Ganti manajer paket JS | Ganti `npm` → `pnpm`/`bun` dalam `build/Taskfile.yml` |

---

## 5. Kiat Membuat Templat

- **Pertahankan frontend tetap generik** — hindari merujuk ke variabel global khusus Wails; `/wails/runtime.js` disajikan oleh server aset saat runtime.
- **Jangan sertakan artefak hasil kompilasi** — kecualikan `node_modules`, `dist`, dan `.DS_Store` dari direktori tersemat (atau `.gitignore` agar artefak tersebut tidak pernah di-commit).
- **Dokumentasikan prasyarat** — versi Node, alat CLI tambahan, dan sebagainya, dalam `template.yaml` atau `NEXTSTEPS.md`.
- **Hindari perubahan yang merusak kompatibilitas** — jika perombakannya besar, buat ID templat baru alih-alih mengubah templat yang sudah ada.

---

## 6. Pemecahan Masalah

| Gejala | Penyebab | Perbaikan |
| --- | --- | --- |
| `unknown template name` | Kesalahan ketik dalam `-t` atau templat tidak disematkan | Jalankan `wails3 init -l` untuk mencantumkan templat yang tersedia |
| Placeholder tidak diganti | Menggunakan `{{ProjectName}}` alih-alih `{{.ProjectName}}` | Tambahkan `.` di awal (akses field templat Go) |
| Server pengembangan membuka halaman kosong | Konfigurasi Vite tidak membaca `WAILS_VITE_PORT` | Periksa `vite.config.ts` Anda |
| Frontend gagal dibangun untuk produksi | Lupa menetapkan jalur `base` Vite | Tetapkan `base: "./"` dalam `vite.config.ts` |

---

## 7. Peta File Sumber Utama

| File | Tanggung Jawab |
| --- | --- |
| `internal/templates/templates.go` | Menyematkan sistem file templat serta mengekspos `Install(options *flags.Init) error`, `GetDefaultTemplates()`, dan `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Konten templat yang sebenarnya |
| `internal/commands/init.go` | Kode penghubung CLI: memilih templat, mengisi metadata, lalu memanggil `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — utilitas untuk *mengekspor* proyek aktif kembali menjadi templat (berguna untuk pembaruan) |
| `internal/flags/init.go` | Definisi flag untuk `wails3 init` |

---

## 8. Ringkasan

- Templat berada di **`internal/templates/`** dan disertakan ke dalam CLI melalui `//go:embed *`.
- `wails3 init -t <id>` mengekstrak templat melalui `gosod` dan menjalankan `go mod tidy` (dapat dilewati dengan `-skipgomodtidy`).
- Membuat templat cukup dengan **membuat folder**, menambahkan file beserta `template.yaml`, dan menggunakan placeholder bergaya `{{.ProjectName}}`.
- Sistem ini **dapat diperluas** dan **mandiri** — sangat cocok untuk membagikan stack khusus kepada tim atau komunitas Anda.

Selamat membuat templat!
