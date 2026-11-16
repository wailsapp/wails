---
title: "Build YAML"
description: "Konfigurasikan dan jalankan sistem build native Wails v3 eksperimental dengan satu file wails.yaml."
sidebar:
  order: 2
  badge:
    text: "Eksperimental"
    variant: caution
---

@note{type="caution" title="Pratinjau eksperimental"}
Rilis pertama yang diterima hanya mendukung konfigurasi literal. Eksperimen ini juga membuat purwarupa hook siklus hidup berbatas agar bentuknya dapat dievaluasi sebelum diusulkan untuk dirilis. Tahap kustom, ekspresi YAML, dan pemanggilan alat Wails bertipe masih ditangguhkan. Proyek yang sudah ada tetap menggunakan Taskfile sampai Anda secara eksplisit mengaktifkan migrasi yang telah ditinjau.
@end

Baru mengenal Wake? Mulailah dengan [tutorial memulai](/experimental/wake/).
Tutorial tersebut membahas batasan eksperimen, setiap layar penyiapan proyek dan build, build lokal pertama, serta migrasi yang telah ditinjau. Halaman ini membahas alur kerja build, pengemasan, dan konfigurasi secara terpisah dengan lebih terperinci.

## Mengaktifkan build YAML

Tetapkan `WAILS_EXP_USE_WAKE` di setiap terminal yang menjalankan perintah eksperimental. Untuk Bash atau Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Untuk PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Variabel ini mengaktifkan fitur berdasarkan keberadaannya: nilai kosong sekalipun atau `false` akan mengaktifkannya. Hapus variabel tersebut (`unset WAILS_EXP_USE_WAKE` di Bash/Zsh atau `Remove-Item Env:WAILS_EXP_USE_WAKE` di PowerShell) untuk menonaktifkan eksperimen. Setelah `wails.yaml` aktif, perintah mengharuskan Wake diaktifkan kembali dan tidak pernah kembali menggunakan Taskfile.
Tanpa variabel tersebut, perutean build/dev YAML dan pembuatan proyek YAML dinonaktifkan, serta perintah `migrate`, `config`, dan `eject` tidak tersedia. Proyek yang belum dimigrasikan tetap diproses melalui runtime Task tertanam yang terpisah.

## Mengonfigurasi proyek

Untuk membuat dan mengonfigurasi proyek di browser, jalankan `wails3 init -ui` dengan opsi keikutsertaan diaktifkan, lalu ikuti [tutorial Wake](/experimental/wake/). Pilihan proyek tersimpan di `wails.yaml`; edit file tersebut untuk mengubah proyek yang sudah ada. `wails3 setup` mengonfigurasi preferensi bersama, alat, dan kredensial penandatanganan.

Setelah pilihan ikut serta diaktifkan, tambahkan `wails.yaml` di akar proyek. Keberadaan file tersebut memilih sistem build native dan membuat Wails mengabaikan setiap Taskfile dalam proyek itu. Pertahankan pilihan ikut serta tetap aktif untuk semua perintah di bawah ini. Modeline skema mengaktifkan pelengkapan dan validasi sebaris di editor yang mendukung server bahasa YAML.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Jalankan build host default atau permintaan rilis bernama yang lengkap:

```bash
wails3 build
wails3 build release
```

Untuk build sekali jalan, pilih target dan format yang dipisahkan koma tanpa menambahkan profil:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Mengompilasi silang aplikasi desktop

Untuk backend desktop eksperimental yang sepenuhnya menggunakan Go, tetapkan `CGO_ENABLED: "0"` di lingkungan build. Wake memilih tag `purego` untuk program GUI macOS dan Linux serta menggunakan Go secara langsung pada host Windows, macOS, dan Linux. Dalam mode ini, membangun executable desktop tidak memerlukan Docker, emulasi, compiler C, atau SDK platform target. Dependensi aplikasi juga harus mendukung `CGO_ENABLED=0`.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

Perintah kedua menggabungkan kedua arsitektur macOS dan membuat DMG tanpa tanda tangan menggunakan Go. Perintah ini dapat dijalankan pada host desktop apa pun yang didukung.

Mode hasil resolusi dipilih secara independen untuk setiap executable. `build.environment.CGO_ENABLED` yang dideklarasikan, lingkungan target, atau lingkungan biner menggantikan nilai terminal yang menjalankan perintah. Jika tidak, tag `purego` yang dipilih secara eksplisit, program CLI/server, atau target Windows menggunakan `0` secara default; program GUI macOS dan Linux menggunakan `1` secara default. `--plan --json` melaporkan compiler, mode CGo, dan tag hasil resolusi untuk setiap executable. Kompilasi silang membuat file target; menjalankan target desktop yang tidak kompatibel tetap memerlukan host yang sesuai.

Dengan `CGO_ENABLED=1`, Wake terlebih dahulu menggunakan toolchain native yang kompatibel atau cross-compiler `CC` yang dikonfigurasi, kemudian Zig untuk build silang Windows jika tersedia, lalu fallback Docker/Podman yang sudah ada jika didukung. `CC` yang ditetapkan secara eksplisit harus menyediakan compiler target, pustaka, dan sysroot; Zig saja tidak menyediakan sysroot GTK/WebKit Linux. Pemaketan dan penandatanganan mungkin memerlukan alat tambahan seperti dijelaskan di bawah. Target seluler memerlukan CGo dan SDK platformnya.

Kompilasi khusus Go tidak membundel pustaka sistem milik tujuan. Aplikasi GUI Linux, termasuk executable pendamping GUI, tetap memerlukan runtime GTK/WebKitGTK yang dipilih oleh tag build-nya. Dependensi paket atau lingkungan deployment yang sesuai harus menyediakannya. Secara default, paket yang hanya berisi CLI/server tidak menambahkan dependensi GUI tersebut; penggantian dependensi secara eksplisit tetap tersedia.

## Membangun beberapa executable

Proyek mendefinisikan program utama. Tetapkan `project.package` ke paket Go-nya (default `.`) dan `project.kind` ke `gui`, `cli`, atau `server` (default `gui`). Program tambahan menggunakan entri bernama di bawah `binaries`:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Setiap program memiliki input kompilasi dan cache tersendiri. Paket proyek yang diimpor dan file tertanam dilacak; pengeditan pada program independen tidak membangun ulang program lainnya. Executable tambahan didukung di desktop serta disertakan dalam bundel aplikasi dan penginstal default. File paket pengganti milik pengguna mengendalikan kontennya sendiri.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Profil bernama merupakan permintaan lengkap; gunakan penggantian pilihan ini dengan build anonim. `--jobs` membatasi worker eksekusi. Memilih hanya program CLI/server akan melewatkan pekerjaan frontend.

Untuk membaca versi tanpa menjalankan kode aplikasi, konfigurasikan `project.version_file` dan `project.version_symbol`. File tersebut harus berisi deklarasi string Go statis. Literal `project.version` kemudian dapat dihilangkan. Tetapkan `commit_symbol: main.Commit` pada proyek atau biner bernama untuk menyuntikkan revisi repositori ke variabel string executable tersebut. Ketiadaan metadata repositori merupakan kesalahan ketika penyuntikan commit diminta.

`upx_level: 1` hingga `9` secara opsional mengompresi executable rilis Linux dan Windows; nol menonaktifkannya. UPX harus terinstal. Build pengembangan dan macOS tidak menggunakan kompresi. Kegagalan kompilasi atau kompresi mempertahankan program yang sebelumnya dipublikasikan.

## Memaketkan dan menandatangani untuk distribusi

### Menggunakan executable desktop yang telah dibangun sebelumnya

Gunakan `--skip-build` jika sistem build atau tugas CI lain telah menghasilkan executable:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Tempatkan setiap executable di `build.output` (nilai defaultnya `bin`), dengan menggunakan `project.binary_name` untuk executable utama dan nama setiap entri `binaries` untuk executable pendamping. Nama untuk Windows diakhiri dengan `.exe`. Pemanggilan dengan beberapa target membaca file setiap target dari subdirektori `<os>-<arch>`. Untuk GUI macOS, Wake membaca executable mentah di `bin/<name>` jika tersedia; jika tidak, dari `bin/<name>.app/Contents/MacOS/<name>`. Input `darwin/universal` harus berisi slice executable amd64 dan arm64.

Wake memvalidasi jenis file dan arsitektur target, menyalin executable ke ruang kerja yang dihasilkannya, serta melacak kontennya untuk cache paket. Input executable yang tidak ada, tidak kompatibel, atau berupa symlink menyebabkan kegagalan alih-alih membangun ulang atau memulihkan input lama. Executable yang diberikan tidak ditimpa atau diubah modenya dengan chmod. Perubahan pada executable mana pun membatalkan cache paket yang menggunakannya.

Mode ini melewatkan instalasi/build frontend, binding, kompilasi Go, UPX, serta hook `before_build`/`after_build`. Mode ini tetap menghasilkan aset paket, menyusun bundel aplikasi macOS, dan menjalankan hook pemaketan/penandatanganan yang dikonfigurasi. Alat pemaketan dan penandatanganan tetap tunduk pada tabel di bawah; Go dan alat build frontend tidak diperlukan. Metadata paket berasal dari manifes saat ini, termasuk `version_file` jika dikonfigurasi; versi tidak disimpulkan dari biner. Gunakan pengaturan build dan paket yang sesuai, termasuk jenis GUI versus CLI.

`wails3 package --skip-build` juga didukung. Dengan `wails3 run`, argumen peluncuran dan lingkungan tetap berlaku, dan executable harus dapat berjalan pada host. Mode ini memerlukan proyek YAML aktif dan hanya mendukung target desktop produksi; mode ini tidak berlaku untuk sesi pengembangan atau pemaketan seluler.

Pembuatan paket tanpa tanda tangan memiliki persyaratan host/alat yang berbeda dari penandatanganan:

Payload IPA dan ZIP notarisasi macOS ditulis dengan Go, dengan mempertahankan jalur bundle, izin executable, dan symlink relatif. Kedua operasi tersebut tidak memerlukan executable `zip` atau `ditto` eksternal. Arsip ini berisi file aplikasi, bukan metadata Finder atau resource fork; penandatanganan macOS memang sudah menghapus metadata tersebut karena metadata itu membuat bundle yang ditandatangani menjadi tidak valid.

| Output | Host dan alat pembuatan | Penandatanganan |
| --- | --- | --- |
| MSIX Windows | Windows, macOS, atau Linux; penulis Go, tidak memerlukan Windows SDK. | Penandatanganan dengan Go menggunakan PFX RSA yang diekspor pada host desktop apa pun; identitas dari penyimpanan sertifikat/penyedia memerlukan Windows SignTool. |
| Penginstal NSIS Windows | Windows, macOS, atau Linux; memerlukan `makensis` native host. | Penandatanganan PE dengan Go menggunakan PFX RSA yang diekspor pada host desktop apa pun; identitas dari penyimpanan sertifikat/penyedia memerlukan Windows SignTool. |
| Aplikasi macOS dan DMG | Windows, macOS, atau Linux; penyusunan dengan Go, penggabungan biner universal, dan penulis DMG portabel. | Penandatanganan dengan Go menggunakan identitas PKCS#12 RSA yang diekspor, atau penandatanganan ad hoc; identitas keychain tetap menggunakan alat macOS. Notarisasi dengan kunci API dan penempelan tiket dijalankan dengan Go. |
| DEB, RPM, Arch, dan Void Linux | Windows, macOS, atau Linux; penulis paket Go native. | Kunci OpenPGP yang diekspor menggunakan penanda tangan bawaan DEB/RPM/Arch. ID kunci DEB/RPM lama memerlukan alat eksternal. Output Void tidak ditandatangani. |
| AppImage Linux | Windows, macOS, atau Linux; deployment dengan Go dan penulis SquashFS, ditambah runtime AppImage target. Pemaketan silang GUI memerlukan pohon runtime Linux yang disediakan. | Tidak ada alur bawaan untuk menandatangani AppImage. |

Untuk DMG yang ditandatangani, Wake biasanya menandatangani salinan privat aplikasi sebelum penulis Go mengompresi image, lalu menandatangani image tersebut. Ini menghindari pemasangan, perluasan, dan pengompresian ulang DMG dengan `hdiutil`. Penandatanganan dijalankan pada setiap permintaan; paket yang ditandatangani tidak digunakan kembali dari cache build. Proyek dengan hook `after_package` atau `before_sign` mempertahankan alur asli berupa pemaketan lalu penandatanganan, termasuk `hdiutil`, sehingga perubahan yang dibuat hook tersebut pada image tanpa tanda tangan tetap dipertahankan. Jalur pengeditan image yang kompatibel dengan hook ini memerlukan macOS dan `hdiutil`.

Untuk penandatanganan Apple portabel, tetapkan `darwin.signing.identity: "-"` untuk tanda tangan ad hoc, atau gunakan `darwin.signing.certificate` untuk identitas PKCS#12 RSA yang diekspor. `darwin.signing.credential` menentukan nama variabel lingkungan kata sandinya. Tanda tangan sertifikat menggunakan implementasi PFX/CMS Go bersama; tidak ada kunci privat yang diekspor secara otomatis dari keychain. Executable tipis dan universal, bundle aplikasi bertingkat, framework, serta DMG telah lolos pemeriksaan tanda tangan native Apple yang ketat, termasuk penolakan terhadap manipulasi, menggunakan identitas pengujian yang dibuat.

Untuk notarisasi, identitas Developer ID yang diekspor juga memerlukan 3161 `darwin.signing.timestamp_server` RFC yang kompatibel. Verifikator bersama menerima profil stempel waktu SHA-2 yang didukungnya dan menolak algoritma lama yang tidak didukung. Tanda tangan ad hoc tidak dapat dinotarisasi. Validasi format sertifikat native terpisah dari kepercayaan Developer ID yang sebenarnya dan penerimaan kiriman oleh Apple.

Penempel Go mengautentikasi tanda tangan tiket Apple dan memverifikasi bahwa tiket tersebut mencakup setiap arsitektur sebelum menempelkannya ke salinan privat. Penempel ini menggunakan protokol pengiriman tiket yang diamati, yang tidak dipublikasikan Apple sebagai API yang didukung. Ini tidak menggantikan pemeriksaan pencabutan terkini oleh Gatekeeper. Output aplikasi dan DMG yang ditandatangani telah lolos pemeriksaan native `stapler validate` dan pemeriksaan tanda tangan yang ketat.

Arsip MSIX Windows dibuat langsung dengan Go; pembuatan paket tanpa tanda tangan tidak memerlukan Windows SDK maupun Docker. Identitas RSA PFX yang diekspor menandatangani berkas executable Windows dan setiap paket MSIX secara langsung dengan Go pada semua host desktop. Tetapkan `windows.signing.certificate` ke PFX yang lokasinya relatif terhadap proyek dan `windows.signing.credential` ke variabel lingkungan yang memuat kata sandinya. Ekspor AES/PBKDF2 dan 3DES lama didukung; ekspor RC2 harus diekspor ulang atau menggunakan Windows SignTool native sebagai fallback. Identitas dari penyimpanan/penyedia sertifikat dan penandatanganan MSI memerlukan Windows SignTool. Layanan stempel waktu default adalah DigiCert; stempel waktu diverifikasi sebelum output yang ditandatangani dipublikasikan. Lihat [panduan penandatanganan](/guides/build/signing/) untuk mengetahui batasan format.
NSIS menggunakan compiler NSIS eksternal dan menerima definisi literal melalui `packages.nsis.defines`. Di macOS, Wails memberikan locale UTF-8 hanya kepada proses turunan compiler NSIS jika locale yang diwarisinya bukan UTF-8; lingkungan terminal dan aplikasi tidak berubah. Installer NSIS bawaan tidak menyertakan bootstrapper WebView2 jika semua executable yang dikemas merupakan CLI atau server. Executable utama atau pendamping yang memiliki GUI tetap menyertakan WebView2. Skrip NSIS khusus tetap menyertakan bootstrapper demi kompatibilitas dengan logika instalasi runtime miliknya sendiri.

Linux mendukung paket AppImage, DEB, RPM, Arch Linux, dan Void. Wake merakit AppImage secara langsung dengan Go, termasuk deployment GTK/WebKit, kompresi SquashFS, dan pembuatan zsync. Wake tidak menjalankan linuxdeploy, plugin GTK, appimagetool, ataupun mksquashfs. Runtime AppImage tipe-2 yang versinya dipatok dan checksum-nya diverifikasi diunduh untuk arsitektur target; `packages.appimage.runtime_path` digunakan untuk memilih ELF runtime yang jalurnya relatif terhadap proyek.

Di Linux, sumber pustaka default adalah sistem file root host. Untuk memaketkan GUI bagi arsitektur lain atau dari Windows/macOS, tetapkan `packages.appimage.runtime_root` ke pohon root Linux target yang jalurnya relatif terhadap proyek dan berisi pustaka, skema terkompilasi, serta cache modulnya. Nilai ini harus berupa seluruh tata letak root, termasuk `usr/`, bukan direktori `usr/` itu sendiri. Sediakan runtime dari distribusi Linux tertua yang ingin Anda dukung. Wake membaca impor ELF dan profil GTK/WebKit eksplisit tanpa menjalankan program target. Wake membundel closure yang dipilih, bukan glibc atau loader dinamis. Hash konten pustaka dan resource membatalkan cache paket setelah runtime diperbarui.

Build GUI Go murni didukung. Wake memilih GTK 3 untuk executable GUI yang dibuat dengan `gtk3`, dan GTK 4 untuk executable lainnya, termasuk executable GUI pendamping CLI/server. Semua executable GUI dalam satu image harus menggunakan versi GTK yang sama. Untuk biner siap pakai, `packages.appimage.gtk_version` dapat secara eksplisit memilih `3` atau `4`; pustaka CGo yang diimpor harus sesuai dengan pilihan tersebut.

Tetapkan `packages.appimage.update_info` untuk meminta file zsync dan `packages.appimage.checksum: true` untuk file SHA256. Sidecar dipulihkan dari cache paket dan dicantumkan dalam tanda terima artefak bersama AppImage. Pemaketan di Windows dan macOS telah diverifikasi dengan pohon GTK/WebKit yang disediakan; image GUI Go murni yang dihasilkan lolos pemeriksaan runtime di Linux. Dalam satu build GUI CGo yang sebanding, jalur native memerlukan 8.77 detik dibandingkan 26.19 detik untuk linuxdeploy. Image gzip native berukuran 133 MiB dibandingkan 114 MiB dengan jalur zstd lama, sehingga pengukuran ini mencakup seluruh jalur pemaketan, bukan hanya kompresi.

Paket Arch mempertahankan prarilis dalam metadata paket native. Misalnya, `1.2.3-beta.2` menjadi `1.2.3beta.2-1`, yang diurutkan sebelum rilis final. Dependensi Void menerima nama paket tanpa keterangan tambahan seperti `gtk4`, yang dinormalisasi menjadi `gtk4>=0`; perbandingan versi eksplisit dan nilai `name-version_revision` lengkap dipertahankan. Gunakan nama paket Void, bukan nama file pustaka.

Penandatanganan DEB, RPM, dan Arch dapat menggunakan kunci privat OpenPGP armored atau biner yang diekspor pada host desktop apa pun yang didukung, tanpa GnuPG, `dpkg-sig`, atau `rpmsign`. Konfigurasikan `linux.signing.certificate` dengan file kunci dan `credential` dengan nama variabel lingkungan frasa sandinya. DEB dan RPM menyematkan tanda tangan; Arch menerbitkan tanda tangan terpisah. Kunci privat hanya dibaca pada langkah penandatanganan dan tidak disalin ke aset paket maupun artefak cache yang dapat digunakan kembali.

Konfigurasi DEB/RPM yang sudah ada dan memuat pengidentifikasi kunci GnuPG polos tetap menggunakan penanda tangan eksternal. Untuk menghapus dependensi tersebut, ekspor kunci rahasia dan konfigurasikan jalurnya sebagai gantinya. Gunakan jalur seperti `./signing/release.asc`; file kunci yang tidak ada, tidak dapat dibaca, atau tidak valid akan menyebabkan kegagalan, bukan fallback ke keyring. Tanda tangan DEB tetap menggunakan format tanda tangan paket `dpkg-sig`; ini tidak menggantikan penandatanganan metadata rilis repositori apt.

Untuk membuat paket Arch dan tanda tangan terpisah, gunakan file kunci privat PGP berformat armored yang telah diekspor dan, jika perlu, tentukan nama variabel lingkungan untuk frasa sandi:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Tetapkan variabel lingkungan bernama tersebut melalui lingkungan lokal atau penyimpanan rahasia CI Anda, lalu jalankan `wails3 build release`. Hilangkan `credential` untuk kunci yang tidak dienkripsi. Penanda tangan Go berjalan di setiap host desktop yang didukung dan tidak memerlukan GPG. Penanda tangan ini menerbitkan `.pkg.tar.zst` asli dan file `.pkg.tar.zst.sig` di lokasi yang sama; keduanya dicatat secara terpisah dalam tanda terima artefak. Isi kunci tidak ditempatkan di area staging atau disimpan dalam cache build. Tanda tangan dan penerbitannya tidak pernah digunakan kembali dari cache. Jika penandatanganan gagal, pasangan yang diterbitkan sebelumnya tetap dipertahankan. Jika penginstalan salah satu file gagal, penerbitan akan mengembalikan kedua file ke keadaan sebelumnya. Hal ini tidak membuat dua operasi penggantian nama pada sistem file menjadi atomik bagi pembaca konkuren.

Bundle aplikasi macOS menyertakan program pendamping. Tetapkan `darwin.microphone_usage_description` untuk menambahkan teks penggunaan dan, jika tidak ada file entitlement khusus yang disediakan, entitlement input audio. Sumber plist dan entitlement khusus tetap tidak berubah. Identitas PKCS#12 RSA yang diekspor menggunakan penandatanganan Go pada setiap host desktop; `identity: "-"` memilih penandatanganan ad hoc portabel. Identitas keychain tetap menggunakan `codesign` di macOS. Penempelan tiket dijalankan dengan Go. Pengiriman notarisasi dapat menggunakan profil keychain notaris macOS yang sudah ada, atau kredensial API tim App Store Connect:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

Mode API mengirimkan, mengunggah, dan melakukan polling dari Go tanpa `notarytool`. Tetapkan ketiga kolom dan hilangkan `credential`; jangan simpan kunci privat dalam kontrol versi. Kunci hanya dibaca pada saat penandatanganan dan tidak disalin ke aset yang dihasilkan ataupun cache. Mode ini mendukung kunci API tim dan arsip hingga 5 GiB; unggahan multipart dan kunci API individual belum diimplementasikan. Kiriman yang ditolak tetap menyimpan ID kiriman dan menyimpan log diagnostik Apple ke file sementara privat yang jalurnya dilaporkan dalam pesan kesalahan. Kiriman yang dibatalkan atau ditolak tidak menggantikan output yang sebelumnya telah dipublikasikan. Setelah diterima, Wake mengambil dan memverifikasi tiket Apple, lalu menempelkannya dengan Go. Dengan identitas yang diekspor dan kredensial API tim, penandatanganan, pengiriman, serta penempelan dapat dijalankan di Linux, Windows, atau macOS. Profil keychain tetap memerlukan macOS dan `notarytool`.

Penandatanganan sertifikat portabel mendukung identitas RSA dalam file PKCS#12 AES/PBKDF2 atau 3DES lama. Ekspor yang dienkripsi dengan RC2 harus diekspor ulang menggunakan cipher yang didukung. Tanda tangan portabel yang dinotarisasi memerlukan layanan stempel waktu SHA-2 RFC 3161 yang ditentukan secara eksplisit; verifikator Go menolak tanda tangan stempel waktu SHA-1 lama. Layanan tiket mengikuti protokol Apple yang diamati dan menyematkan Apple Root G3; Gatekeeper tetap menjadi otoritas untuk kepercayaan distribusi dan pencabutan.

Output yang ditandatangani tetap mempertahankan ekstensi platformnya, misalnya
`notes.signed.app` dan `notes.signed.dmg`.

## Kontrol pengembangan

Selama `wails3 dev` berjalan, masukkan perintah lalu tekan Enter:

| Perintah | Tindakan |
| --- | --- |
| `r` | Paksa build ulang dan mulai ulang, termasuk backend yang tidak berubah. |
| `b` | Buat ulang binding dan lakukan build ulang. |
| `f` | Mulai ulang frontend. |
| `c` | Bersihkan terminal. |
| `q` | Hentikan sesi beserta proses turunannya. |
| `h` | Tampilkan kontrol. |

Build yang gagal membiarkan aplikasi saat ini tetap berjalan. Kontrol memerlukan terminal interaktif; input yang disalurkan melalui pipe tidak mengubah perilaku sesi.

## Periksa pekerjaan sebelum menjalankannya

Pemeriksaan rencana me-resolve graf imutabel yang sama dengan yang digunakan oleh eksekusi. Teks ditujukan untuk manusia; JSON adalah kontrak terbaca mesin yang memiliki versi.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Validasi konfigurasi tanpa melakukan build

Gunakan perintah preflight khusus di CI, sebagai tugas editor, atau sebelum melakukan commit konfigurasi. Tanpa profil, perintah ini memvalidasi manifes dasar dan setiap profil bernama; dengan profil, perintah ini hanya memeriksa permintaan lengkap tersebut:

```bash
wails3 config check
wails3 config check release
```

Kegagalan sintaks, skema, semantik, jalur, target, format, penandatanganan, perencana, dan kemampuan host yang ditentukan konfigurasi menggunakan perender diagnostik yang sama dengan Build dan Dev. Diagnostik mengidentifikasi kolom dan rentang sumber, menampilkan baris sumber beserta tanda sisipan, serta menyertakan petunjuk koreksi jika berguna. Nilai lingkungan disamarkan dalam cuplikan sumber diagnostik. Kegagalan host yang hanya disebabkan oleh pilihan baris perintah tetap dilaporkan sebagai kesalahan perintah biasa, bukan secara keliru sebagai YAML yang tidak valid.

## Gunakan aset milik proyek

Jalur manifes merujuk pada file milik Anda, seperti ikon, manifes platform, templat paket, dan skrip maintainer. Wails tidak pernah mengedit sumber tersebut. Wails menyalinnya ke ruang kerja sementara `.wails/` sebelum digunakan oleh suatu alat. Artefak final ditempatkan di bawah `build.output` dan tanda terima terverifikasi ditulis ke `.wails/artifacts/receipt.json`.

## Deploy ke Android

Profil Android produksi menghasilkan AAB. Untuk deployment lokal, Wails membuat APK debug khusus ABI secara internal, menginstalnya dengan ADB, lalu meluncurkan aplikasi:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Tanpa pemilih eksplisit, terminal interaktif menawarkan perangkat yang terhubung dan Android Virtual Device yang telah dikonfigurasi. Otomatisasi dengan beberapa target yang tersedia harus memberikan `--device` atau `--emulator`.

Lampirkan aliran log dengan cakupan aplikasi dan pantau target hingga aplikasi berhenti atau Anda menekan Ctrl+C:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Proses yang dijalankan dalam mode terlampir akan mengikuti aplikasi ketika dimulai ulang dan melaporkan terputusnya koneksi target atau transisi ke status offline. Secara default, emulator tetap berjalan agar iterasi dapat dilakukan dengan cepat. `--stop-emulator` harus ditentukan secara eksplisit, memerlukan `--logs`, dan hanya menghentikan target emulator saat sesi terlampir berakhir. Jika Wails memulai emulator baru tetapi emulator gagal mencapai kesiapan boot atau proses dibatalkan, Wails otomatis membersihkan proses tersebut.

## Jalankan hook siklus hidup dengan batasan

Antarmuka hook eksperimental menjalankan satu skrip milik proyek pada setiap fase yang dideklarasikan. Hook bukan perintah shell: Wails mengeksekusi file secara langsung dari root proyek, kecuali jika `directory` ditetapkan.

Pada host Unix, skrip harus dapat dieksekusi dan diawali dengan shebang. Host Windows menerima skrip `.cmd`, `.bat`, dan `.ps1` melalui adaptor interpreter tetap; Wails tidak mengurai perintah shell inline pada kedua platform.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

Fase yang didukung adalah `before_build`, `after_build`, `before_package`, `after_package`, `before_sign`, dan `after_sign`. Hook build dijalankan satu kali untuk proyek bersama atau satu kali per target, sesuai kebutuhan. Hook paket dan penandatanganan membentuk penghalang yang mengapit semua format yang diminta untuk suatu target, sementara operasi format di dalam penghalang tetap berjalan secara konkuren.

Hook menerima satu variabel lingkungan milik Wails: `WAILS_HOOK_CONTEXT_FILE`. Variabel ini menunjuk ke file JSON terlindungi yang berisi konteks pemanggilan berversi:

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

Konteks tersebut memuat jalur yang telah di-resolve, tetapi tidak pernah menyalin lingkungan proses yang diwariskan maupun rahasia. Wails membuatnya secara atomik dengan izin file khusus pemilik di bawah `.wails/hooks/context/`. File tersebut dihapus setelah berhasil atau dibatalkan. Jika terjadi kegagalan, file dipertahankan untuk diagnosis dan jalurnya disertakan dalam pesan kesalahan. Kode keluar bukan nol menggagalkan langkah Plan normal, dan pembatalan menghentikan grup proses hook.

Secara default, hook berjalan pada setiap pemanggilan. Aktifkan caching hanya jika skrip memiliki semua input dan output yang dideklarasikan:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

Skrip dan mode executable-nya merupakan input implisit. Beberapa output harus berbagi satu direktori non-root, yang menjadi artefak terbatas yang dicatat dan dipulihkan oleh Wails. Root output tidak boleh memuat skrip atau input yang dideklarasikan, dan pohon output yang di-cache hanya boleh memuat file biasa dan direktori, tidak pernah tautan simbolis atau file khusus. Hook yang selalu dijalankan membuat operasi setelahnya tidak dapat di-cache karena efek samping yang tidak dideklarasikan tidak dapat dibuat sidik jarinya secara aman. Hook yang dapat di-cache mempertahankan caching untuk operasi setelahnya melalui identitas artefak yang dideklarasikan.

## Migrasikan proyek yang sudah ada

Buat dan tinjau draf yang belum aktif:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Tugas siklus hidup konservatif tanpa argumen yang menjalankan skrip milik proyek
diterjemahkan menjadi hook. Perilaku Taskfile kustom lain yang dapat dijangkau
dilaporkan sebagai penghambat. Setelah meninjau
`wails.migrated.yaml`, aktifkan secara eksplisit:

```bash
wails3 migrate --activate
```

Aktivasi membuat `wails.yaml` secara atomik. Aktivasi tidak menghapus, mengganti nama, atau mengedit
Taskfile, `build/config.yml`, skrip, maupun aset. YAML yang tidak valid akan gagal sebagai YAML;
sistem tidak pernah beralih kembali ke jalur lama.

## Materialisasikan semua nilai default

Untuk memeriksa atau menyalin konfigurasi lengkap yang telah di-resolve, tulis snapshot
tidak aktif:

```bash
wails3 eject
```

Tindakan ini membuat `wails.ejected.yaml`. Tindakan ini tidak pernah menulis ulang konfigurasi aktif.

## Periode kompatibilitas

`wails3 package` dan `wails3 sign` adalah alias usang untuk perencana yang sama;
otomatisasi baru sebaiknya menggunakan `wails3 build` dengan profil atau format. Peluncurannya
dibatasi oleh rilis: rilis pertama menjadikan YAML sebagai jalur untuk proyek baru, rilis berikutnya
memperingatkan tentang perutean Taskfile lama, dan rilis selanjutnya menghapus eksekusi Taskfile
beserta alias tersebut. Setiap transisi mengharuskan matriks pengujian penerimaan host native
untuk rilis tersebut lulus.

Lihat referensi kolom [`wails.yaml` yang dihasilkan](/reference/wails-yaml-fields/)
untuk setiap kolom, tipe, nilai default, dan aturan penerapan.

Referensi konfigurasi [`wails3 run`](/reference/wails-yaml-run/) menjelaskan
nilai default peluncuran lokal, penggantian khusus platform, dan fallback `go run .` tanpa manifes.
