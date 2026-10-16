---
title: "Membuat Templat Kustom"
description: "Cara membuat, menyesuaikan, dan meng-host templat proyek Wails v3 Anda sendiri"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails menyediakan serangkaian templat bawaan, tetapi Anda dapat membuat templat sendiri dan membagikannya kepada komunitas. Templat kustom hanyalah repositori Git — setelah di-host secara publik, siapa pun dapat membuat kerangka proyek darinya dengan satu perintah.

## Membuat kerangka templat

Perintah `wails3 generate template` menghasilkan direktori templat yang siap disesuaikan:

```bash
wails3 generate template -name MyTemplate
```

Semua flag:

| Flag | Deskripsi | Nilai bawaan |
| --- | --- | --- |
| `-name` | Nama templat (wajib) | — |
| `-author` | Nama penulis | — |
| `-description` | Deskripsi singkat yang ditampilkan di CLI | — |
| `-helpurl` | URL dokumentasi untuk templat ini | — |
| `-version` | Versi awal | `v0.0.1` |
| `-frontend` | Salin direktori frontend yang sudah ada ke dalam templat | — |
| `-dir` | Lokasi untuk menulis direktori templat | Direktori saat ini |

Contoh dengan semua flag:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

Direktori yang dihasilkan terlihat seperti ini:

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="Baca NEXTSTEPS.md"}
`NEXTSTEPS.md` yang dihasilkan berisi panduan mendetail untuk setiap bagian templat. Bacalah sebelum menyesuaikan templat. Hapus file tersebut sebelum memublikasikan — file itu tidak boleh muncul dalam proyek yang dibuat dari templat Anda.

@end

## Mengonfigurasi metadata templat

Buka `template.yaml` untuk mengatur metadata templat Anda:

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

Kolom `wailsVersion` **wajib** diisi dan harus berupa `3`. Komentar `# yaml-language-server` di bagian atas mengaktifkan pelengkapan otomatis dan validasi sebaris di VS Code (dengan [ekstensi YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) serta IDE JetBrains — Anda dapat membiarkannya atau menghapusnya; komentar tersebut tidak berpengaruh saat runtime.

## Menyesuaikan templat

### Frontend

Direktori `frontend/` disalin apa adanya ke setiap proyek yang dibuat dari templat Anda. Ganti konten placeholder dengan frontend Anda yang sebenarnya:

@tabs
[Mulai dari awal]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Ikuti petunjuk yang muncul, lalu instal dependensi:

```bash
npm install
```

[Gunakan proyek yang sudah ada]
Sertakan `-frontend` saat membuat templat untuk menyalin frontend yang sudah ada dalam satu langkah:

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Atau salin secara manual ke direktori `frontend/` setelahnya.

@end

### Tugas build

`Taskfile.tmpl.yml` mendefinisikan alur kerja build. Perbarui tugas `install:frontend:deps` dan `build:frontend` agar sesuai dengan toolchain frontend Anda:

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Aplikasi Go

File `main.go.tmpl` adalah titik masuk aplikasi. File ini diproses oleh mesin templat Wails saat proyek dibuat — variabel templat seperti `{{.ProductName}}` diganti dengan nilai yang diberikan pengguna.

Untuk mengeditnya sebagai file Go sungguhan (dengan dukungan IDE), ganti namanya sementara menjadi `main.go`, buat perubahan Anda, lalu kembalikan namanya menjadi `main.go.tmpl` sebelum melakukan commit.

#### Variabel templat

Variabel berikut tersedia di semua file `.tmpl`:

| Variabel | Deskripsi | Contoh |
| --- | --- | --- |
| `{{.ProjectName}}` | Nama proyek yang diberikan pengguna | `"MyApp"` |
| `{{.BinaryName}}` | Nama file biner | `"myapp"` |
| `{{.ProductName}}` | Nama tampilan produk | `"My Application"` |
| `{{.ProductDescription}}` | Deskripsi produk | `"An awesome application"` |
| `{{.ProductVersion}}` | Versi produk | `"1.0.0"` |
| `{{.ProductCompany}}` | Nama perusahaan / penulis | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | String hak cipta | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Komentar tambahan tentang produk | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Pengidentifikasi produk dengan format DNS terbalik | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Jalur modul Go | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Versi Wails yang digunakan untuk membuat proyek | `"3.0.0"` |
| `{{.Typescript}}` | `true` jika nama templat diakhiri dengan `-ts` | `true` |
| `{{.Opn}}` | `{{` literal — lakukan escape di dalam templat | `{{` |
| `{{.Cls}}` | `}}` literal — lakukan escape di dalam templat | `}}` |

@note{type="tip"}
File apa pun dalam templat Anda dapat berupa file `.tmpl` — termasuk file HTML, JSON, dan YAML. File tanpa akhiran `.tmpl` disalin apa adanya.

@end

## Uji templat Anda secara lokal

Sebelum memublikasikannya, uji templat dengan membuat proyek dari jalur lokal:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Kemudian, pastikan proyek berfungsi:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Pastikan bahwa:

- Hot reload frontend berfungsi
- Perubahan kode Go membangun ulang dan menjalankan kembali aplikasi
- Biner produksi di `bin/` berjalan dengan benar

## Publikasikan di GitHub

@steps
### **Buat repositori GitHub publik** untuk templat Anda. Root repositori harus berisi `template.yaml`.
### **Hapus `NEXTSTEPS.md`** — file ini merupakan panduan bagi pembuat templat dan tidak boleh muncul dalam proyek yang dibuat pengguna dari templat Anda.
### **Commit dan push** isi direktori templat sebagai root repositori:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Beri tag pada rilis** menggunakan versioning semantik:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Pengguna kini dapat membuat proyek dari templat Anda:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Peringatan templat pihak ketiga"}
Saat pengguna menginstal templat jarak jauh, Wails menampilkan peringatan yang menjelaskan bahwa templat tersebut adalah kode pihak ketiga dan proyek Wails tidak bertanggung jawab atas isinya. Pengguna harus memberikan konfirmasi secara eksplisit sebelum proyek dibuat.

Sebagai pembuat templat, Anda bertanggung jawab atas keamanan dan kebenaran seluruh kode dalam templat Anda.

@end

## Praktik terbaik

- **Tulis `README.md`** yang jelas — ini ditampilkan kepada pengguna setelah mereka membuat proyek. Jelaskan cara menjalankan, membangun, dan menyesuaikan proyek.
- **Isi `helpurl`** — tautkan ke repositori atau dokumentasi khusus Anda. Pengguna melihatnya dalam daftar templat Wails CLI.
- **Kunci versi dependensi frontend** dalam `package.json` untuk mencegah kegagalan instalasi akibat pembaruan upstream.
- **Uji sebelum memberi tag** — buat proyek baru dari rilis yang telah diberi tag sebelum mengumumkannya kepada komunitas.
- **Pertahankan `wailsVersion: 3`** — kolom ini memberi tahu Wails versi mayor yang ditargetkan templat. Jangan mengubahnya.
- **Perbarui secara berkala** — selalu mutakhirkan dependensi dan uji dengan rilis Wails baru.
