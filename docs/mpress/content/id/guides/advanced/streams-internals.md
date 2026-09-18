---
title: "Stream — internal"
description: "Cara kerja transport stream, alasan arsitekturnya dibuat seperti ini, arti konstanta buffer, dan bagian yang belum selesai"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

Referensi bagi siapa pun—manusia maupun agen—yang mengubah transport stream. API untuk pengguna terdapat di [Stream](/guides/streams/); halaman ini membahas mekanisme di baliknya dan pertimbangan yang mendasarinya, karena beberapa keputusan tampak arbitrer sebelum Anda mengetahui hal yang ingin dihindari oleh keputusan tersebut.

## File

| file | peran |
| --- | --- |
| `v3/pkg/application/stream.go` | API publik, `StreamConn`, `streamSink`, pengelola, dan registrinya |
| `v3/pkg/application/stream_session.go` | satu pemuatan halaman dalam satu jendela: antrean keluar, jenis frame, tabel koneksi |
| `v3/pkg/application/stream_transport.go` | dua endpoint HTTP, framing biner, perakitan ulang chunk, prelude runtime |
| `v3/pkg/application/stream_server.go` | khusus `-tags server`: sink WebSocket yang sebenarnya |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | memilih transport klien saat bundle disajikan |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | klien berbentuk `WebSocket` |
| `v3/tests/stream-performance/` | harness pengujian beban (`-upload`, `-reloads`, rangkaian skenario) |

## Bentuk arsitektur

Go→JS dan JS→Go menggunakan mekanisme yang berbeda, dan asimetri itulah keseluruhan desainnya.

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

**Go→JS menggunakan polling yang ditahan.** Permintaan menunggu hingga ada sesuatu yang dapat dikirim. Sengaja tidak ada interval polling ataupun mekanisme adaptif: server menahan permintaan hingga tersedia sebuah frame, sehingga latensi pengiriman sudah sekitar ~0, dan interval apa pun di sisi klien hanya akan menambah latensi tersebut. Frame yang tiba ketika respons sedang diproses akan terkumpul dan ikut dalam respons berikutnya, sehingga perjalanan bolak-balik itu sendiri menjadi jendela batching—jendela ini melebar saat beban meningkat tanpa perlu ada mekanisme yang mengukurnya. Hasil pengukuran: 1.0 frame per respons pada 100/dtk, masih 1.0 pada 5000/dtk, 3.4 pada 20000/dtk, dan latensi p99 *menurun* seiring meningkatnya laju.

**JS→Go menggunakan POST biasa.** Pengiriman diserialisasi per koneksi dengan rantai promise, karena pemanggilan `fetch` secara bersamaan tidak mempertahankan urutan, sedangkan Go mengandalkan urutan pengiriman sebagai urutan yang diamatinya. Frame yang terkumpul di belakang permintaan yang sedang diproses digabungkan ke dalam POST berikutnya. Go menambahkan frame yang diterima atau prefiks batch ke kotak masuk koneksi *sebelum* merespons, sehingga klien tidak dapat melanjutkan melewati byte yang belum dimasukkan Go ke antrean.

**Satu polling yang sedang diproses per jendela, dengan semua koneksi dimultipleks.** Inilah yang menjamin urutan benar secara inheren—satu antrean, satu penguras, tanpa jalur pengiriman kedua yang dapat mendahului jalur pertama. Hal ini juga menghindari batas enam koneksi per host pada HTTP/1.1 di Windows, tempat permintaan tersebut merupakan permintaan jaringan Chromium yang sebenarnya ke `http://wails.localhost`.

## Alasan di balik keputusan khusus ini

Setiap keputusan ini merupakan bekas dari pekerjaan pada transport peristiwa. Menghapus salah satunya akan memunculkan kembali bug yang telah terukur.

**Tidak ada bagian jalur Go→JS yang menyentuh thread utama.** `Send` menambahkan data dengan perlindungan mutex lalu kembali. Sebelumnya, peristiwa menjalankan eval secara inline ketika dipancarkan dari thread utama sementara pemancaran terdahulu dari goroutine masih berada dalam antrean—urutan 4.4% peristiwa terbalik pada ketiga platform. Satu antrean dengan satu penguras tidak dapat mengalami hal itu.

**Tidak ada yang menyentuh `evaluateJavaScript`, berapa pun ukurannya.** Menyisipkan payload ke sumber eval mempertahankan memori host di atas ambang yang berbeda untuk setiap platform: 11.6 GB di macOS dan 6.2 GB di WebKitGTK pada 100 × 1 MB/dtk. Stream sama sekali tidak mendekatinya; karena itu, rangkaian pengujian dengan laju byte konstan tetap datar pada setiap ukuran frame.

**Data kontrol dikirim melalui header, tidak pernah melalui body atau string kueri.** WebKitGTK 6.0 dapat mengirimkan body POST sebagai parameter kueri untuk skema URI khusus (`transport_http.go` menyediakan fallback khusus untuk kasus tersebut), sedangkan WebView2 membatasi pengiriman body pada sekitar 2 MB.

**Respons polling berbentuk biner, bukan JSON.** Frame adalah `[]byte`; base64 di dalam pembungkus JSON akan menambah biaya sebesar 33% pada setiap frame, ditambah parsing pada thread UI.

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind` adalah data / open / close / error. Tidak ada nomor urut maupun ack: WebSocket tidak melakukan replay, dan koneksi yang terputus kehilangan data yang sedang dikirim. Menirukan perilaku tersebut lebih sederhana dan lebih jujur daripada menggunakan kursor yang tidak selalu dapat dipenuhi oleh buffer terbatas.

**Menahan permintaan aman dilakukan** karena setiap permintaan webview sudah mendapatkan goroutine tersendiri. `dispatchWorkers` di `assetserver_webview.go` ditetapkan pada 0 dengan komentar yang secara khusus menyebutkan kasus ini; mengaktifkan pool tersebut terlebih dahulu memerlukan batas masa aktif permintaan.

## Konstanta buffer

Semuanya berada di `stream.go`. **Konstanta tersebut merupakan konstanta waktu kompilasi, bukan opsi**—tidak ada `Options.Streams` maupun pengaturan per stream. Untuk mengubahnya, edit file tersebut.

| konstanta | nilai | hal yang dibatasi |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 MB | byte yang disimpan dalam buffer per jendela sambil menunggu diambil |
| `streamOutQueueDepth` | 256 | frame yang disimpan dalam buffer per jendela |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 MB / 8192 | data keluar yang disimpan dalam buffer di seluruh aplikasi |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 MB / 8192 | data masuk yang menunggu `Receive` di seluruh aplikasi |
| `streamMaxConnections` | 256 | koneksi aktif ditambah penutupan yang mengantre dalam satu sesi |
| `streamMaxConnectionsGlobal` | 4096 | koneksi aktif di seluruh aplikasi |
| `streamOutCloseDepthGlobal` | 4096 | notifikasi penutupan yang belum dikirim di seluruh aplikasi |
| `streamMaxSessionsPerWindow` | 16 | jumlah sesi yang dapat dipertahankan oleh satu jendela sebelum generasi yang lebih baru harus menggantikan generasi yang lebih lama |
| `streamMaxSessions` | 1024 | sesi di seluruh aplikasi |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | frame kontrol non-close yang mengantre, per sesi dan di seluruh aplikasi |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | unggahan yang belum lengkap per sesi, dan bagian dalam satu unggahan |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 MB / 4096 | payload chunk dan metadata bagian di seluruh aplikasi |
| `streamMaxChunkIDLen` | 64 byte | satu pengidentifikasi kumpulan chunk yang diberikan klien |
| `streamMaxResponseBytes` | 1 MB | satu respons polling |
| `streamHoldTimeout` | 20 dtk | durasi polling kosong menunggu |
| `streamSessionTTL` | 60 dtk | tidak ada polling selama ini dan tidak ada koneksi aktif ⇒ sesi mati |
| `streamSessionGrace` | 10 mnt | tidak ada polling selama ini meskipun ada koneksi aktif ⇒ sesi mati |
| `streamSessionSweep` | 20 dtk | seberapa sering pembersih mencari sesi yang mati |
| `streamMaxFrameBytes` | 64 MB | satu frame di kedua arah |
| `streamMaxNameLen` | 256 byte | satu nama stream yang didaftarkan atau diminta |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 MB | frame yang diterima dan belum diambil oleh `Receive` |

### Cara memilih nilainya

**`streamOutQueueDepth` sengaja tidak bernilai `eventQueueCapacity` (64).** Konstanta tersebut diukur untuk antrean yang dikosongkan satu eval setiap kali, sehingga penambahan kedalaman hanya meningkatkan latensi ekor. Polling mengosongkan antrean secara berkelompok, jadi kedalaman di sini harus mencakup produksi selama satu perjalanan bolak-balik — pada 5000 frame/dtk dan perjalanan bolak-balik selama 5 mdtk, sekitar 25 frame. 256 menyediakan ruang untuk lonjakan tanpa membuat produsen terhenti.

**`streamOutQueueBytes` adalah batas yang benar-benar penting**, karena 256 frame berukuran 1 MB sama dengan 256 MB. Batas ini menjadi pengaman terakhir bagi memori host ketika frontend berhenti mengambil data.

Ada dua aturan yang saling berinteraksi di sini, dan aturan kedua mudah rusak secara tidak sengaja:

- Batas kedalaman dan byte membatasi *akumulasi*.
- **Antrean kosong selalu menerima satu frame, berapa pun ukurannya.** Penerapan batas byte tanpa pengecualian membuat frame yang lebih besar daripada batas sama sekali tidak mungkin dikirim — kondisi tunggu tidak akan pernah terpenuhi, sehingga `Send` terblokir selamanya dan `TrySend` terus melaporkan kondisi penuh. Ukuran frame tidak selalu dapat ditentukan oleh pemanggil; struct dengan field `[]byte` akan di-marshal sesuai hasil marshal-nya.

**`streamMaxResponseBytes` ada karena Windows.** Penulis respons WebView2 mengakumulasi seluruh isi di memori dan baru menyerahkannya di `Finish`, sehingga respons tanpa batas menyebabkan alokasi tanpa batas di sana. Menaikkan nilainya *tidak* meningkatkan throughput Windows — berdasarkan pengukuran, hambatan Windows terjadi per byte, bukan per respons: respons/dtk bervariasi 4× sepanjang pengujian berbagai ukuran frame, sementara MB/dtk tetap stabil pada sekitar 90.

**Batas masuklah yang membuat frontend menunggu.** Di desktop, `deliver` melaporkan kondisi penuh, endpoint merespons dengan `429`, dan klien mencoba kembali frame yang sama atau bagian akhir batch yang belum diterima dengan backoff terbatas. Hal ini menghindari penggunaan slot permintaan webview saat handler mengejar ketertinggalan. Dalam mode server, pompa pembacaan soket menunggu dan membiarkan TCP menerapkan backpressure. Tanpa batas tersebut, handler yang lambat memanggil `Receive` dapat membuat penggunaan memori host bertambah tanpa batas.

**Frame kontrol melewati batas data, tetapi memiliki batas siklus hidup tersendiri.** Kehilangan frame data akibat backpressure menyebabkan perlambatan; kehilangan konfirmasi pembukaan membuat frontend terus berada dalam `CONNECTING` selamanya, sedangkan kehilangan penutupan membuatnya mengira koneksi yang mati masih aktif. Karena itu, kontrol non-close memiliki antrean terbatas sendiri, yang terpisah dari antrean yang digunakan oleh close, agar lonjakan pembukaan yang ditolak tidak dapat menghabiskan kapasitas yang dibutuhkan koneksi yang diterima untuk melaporkan bahwa koneksi tersebut telah berakhir. Setiap sesi juga mencadangkan satu slot close untuk setiap koneksi yang diterima. Ketika kapasitas itu terpakai, pembukaan baru menerima backpressure yang dapat dicoba ulang sebelum didaftarkan.

**Batas per sesi juga memiliki batas padanan di seluruh aplikasi.** Tanpanya, setiap sesi atau koneksi yang diterima dapat sekaligus mempertahankan seluruh jatah lokalnya. Karena itu, data keluar dan masuk berbagi anggaran terpisah sebesar 256 MiB / 8192 frame di seluruh transport desktop dan server. Koneksi aktif memiliki anggaran 4096 entri, dan notifikasi close yang belum dikirim memiliki anggaran kedua dengan ukuran yang sama. Ketika jatah bersama tercapai, perilaku `Send` yang memblokir atau `TrySend` yang tidak memblokir berlaku sama seperti ketika jatah lokal tercapai, dan setiap jalur pengosongan, penerimaan, penutupan, penulisan yang gagal, serta penghentian mengembalikan reservasinya.

Kedua anggaran tersebut sengaja dipisahkan, bukan dijadikan satu jatah yang diserahkan koneksi kepada frame close miliknya. Setiap reservasi dilepaskan oleh tepat satu pemilik: slot koneksi oleh `shutdown`, yang dijalankan sekali, dan slot frame close oleh mekanisme apa pun yang membuang frame tersebut — pengosongan antrean atau pembongkaran sesinya. Kepemilikan yang berpindah di antara dua pihak harus ditransfer secara atomik, dan revisi sebelumnya yang memungkinkan close mewarisi slot koneksi membocorkan satu slot secara permanen setiap kali pembongkaran berlangsung di antara percobaan close dan kegagalan percobaan tersebut.

**Frame Go mentransfer kepemilikan; snapshot frame JavaScript dibuat.** `Send` Go mempertahankan slice milik pemanggil hingga transport menuliskannya, sehingga pemanggil tidak boleh memutasi atau menggunakan kembali penyimpanan tersebut setelah panggilan berhasil. `send()` JavaScript menyalin input biner yang dapat dimutasi sebelum kembali, sesuai dengan semantik kepemilikan WebSocket native. Aturan asimetris ini menghindari penyalinan seluruh frame untuk kedua kalinya di dalam Go sekaligus menjaga agar API yang digunakan browser tetap sesuai ekspektasi.

**Pengiriman JavaScript mengikuti kontrak buffering WebSocket.** `send()` tidak dapat memblokir, sehingga aplikasi dapat mengantrekan data lebih cepat daripada kemampuan kanal permintaan desktop untuk menerimanya, sebagaimana aplikasi dapat melampaui kemampuan WebSocket native. `bufferedAmount` mencakup setiap byte yang dipertahankan oleh soket tersebut dan menjadi sinyal backpressure bagi pemanggil; antrean di sisi host tetap dibatasi secara independen oleh batas di atas. Kegagalan terminal, penutupan oleh peer, atau `close()` lokal melepaskan payload yang dipertahankan. Penutupan lokal juga membatalkan permintaan pembukaan atau data yang sedang menunggu di `429` sebelum mengirim kontrol close yang telah dicadangkan, sehingga backpressure penerimaan atau penerima tidak dapat membuat soket macet dalam `CLOSING`.

**Perakitan ulang potongan memiliki alokasi memori host bersama.** Setiap sesi dapat merakit satu frame hingga 64 MiB, tetapi alokasi tersebut tidak boleh dikalikan untuk setiap sesi yang diterima. Karena itu, kumpulan potongan yang belum lengkap dan dapat dicoba ulang berbagi anggaran payload yang diterima sebesar 128 MiB. Penyelesaian suatu kumpulan mempertahankan bagian-bagiannya sekaligus frame rakitan kontigunya untuk sementara, sehingga menggandakan alokasi logis tersebut masih tetap berada dalam batas atas memori efektif sebesar 256 MiB. Bagian yang dipertahankan juga berbagi alokasi metadata sebanyak 4096 entri, sehingga potongan kecil atau kosong tidak dapat memperbesar map dan pembukuan slice tanpa mendekati batas byte. Permintaan yang akan melampaui salah satu alokasi menerima backpressure yang dapat dicoba ulang; pengiriman, penolakan, kedaluwarsa, atau penghentian sesi mengembalikan byte dan entri bagian ke anggaran bersama.

**Polling hanya dicoba ulang jika kegagalannya dapat dipulihkan.** Kesalahan jaringan, respons waktu permintaan habis (`408`), respons data awal (`425`), backpressure (`429`), dan kesalahan server (`5xx`) menggunakan backoff eksponensial mulai dari 250 ms hingga 5 detik. Respons `4xx` lainnya merupakan kegagalan protokol atau kepemilikan dan segera menutup Streams milik halaman; `410` adalah sinyal terminal bersih untuk sesi yang telah dihentikan. Menutup koneksi terakhir membatalkan polling yang sedang berlangsung atau timer backoff, sedangkan koneksi yang dibuka selama proses penghentian tersebut memulai satu loop polling pengganti.

**`streamSessionTTL` harus tetap jauh di atas `streamHoldTimeout`**; jika tidak, sesi akan dibersihkan saat polling-nya sendiri secara sah sedang menunggu.

Jika Anda melakukan penyetelan untuk beban kerja dengan banyak pesan kecil, batas kedalaman tercapai lebih dahulu; untuk payload besar, batas byte yang tercapai lebih dahulu. Keduanya tidak perlu diubah untuk aplikasi biasa—di macOS, nilai default mampu menangani 634000 frame/dtk dan 2100 MB/dtk.

## Siklus hidup koneksi dan sesi

Sebuah **sesi** adalah satu pemuatan halaman dalam satu jendela, dengan id yang dibuat oleh klien sebagai kunci (seperti `clientId` milik runtime). Sesi dibuat secara lazy oleh permintaan mana pun yang tiba lebih dahulu. Ketika platform tidak dapat mengidentifikasi jendela yang membuat permintaan (`windowID == 0`), jumlah id sesi tetap dibatasi secara global, tetapi generasinya sengaja tidak dibandingkan: sesi-sesi tersebut mungkin berasal dari klien browser independen dengan penghitung generasi yang tidak saling terkait. Sesi tersebut kedaluwarsa melalui penutupan atau TTL, bukan dengan saling menggantikan.

Tiga mekanisme menutup komponen, diurutkan berdasarkan seberapa cepat mekanisme tersebut mendeteksinya:

1. **Polling sesi yang lebih baru untuk suatu jendela menggantikan generasi yang lebih lama.** Pemuatan ulang memberi halaman id sesi baru dan menaikkan generasi yang disimpan dalam `sessionStorage` milik jendela tersebut. Nilai yang sama dicerminkan dalam `window.name`, yang tetap bertahan setelah pemuatan ulang ketika penyimpanan dinonaktifkan, dan ditambatkan ke `performance.timeOrigin` (atau `Date.now()` pada engine lama) agar penghapusan kedua penyimpanan tidak memulai ulang dari satu. Setiap permintaan membawa id sesi dan generasi. Polling hanya menghentikan generasi halaman yang lebih rendah, sehingga penjadwalan server tidak dapat membuat permintaan tertunda dari halaman sebelumnya tampak lebih baru daripada penggantinya. Jika kebijakan memblokir penyimpanan dan `window.name`, pengurutan kembali menggunakan jam halaman sehingga bergantung pada halaman yang lebih baru untuk menerima titik awal waktu yang lebih kemudian. Pengelola mempertahankan watermark generasi yang telah dihentikan untuk setiap jendela agar permintaan lama yang sudah berlangsung tidak dapat membuat ulang halaman lama tanpa perlu mempertahankan setiap id sesi historis. Koneksi sesi sebelumnya segera ditutup.
2. **Penghancuran jendela** menghapus setiap sesi untuk jendela tersebut, seperti `eventPayloadStore.dropWindow`.
3. **Pembersihan berdasarkan TTL** menangani semua kasus lainnya—renderer yang mengalami crash atau mesin yang tertidur.

Hanya dua mekanisme pertama yang menghentikan generasi halaman. Pembersihan TTL menghapus sesi yang tidak aktif tanpa memajukan watermark generasi yang telah dihentikan: halaman berhenti melakukan polling setiap kali koneksi terakhirnya ditutup, tetapi halaman sama yang masih dimuat tersebut harus dapat membuka stream lain nanti. Generasi yang benar-benar telah digantikan tetap diblokir karena polling halaman yang lebih baru memajukan watermark sebelum sesi lama dihapus.

**WebView Apple melaporkan permintaan yang dibatalkan.** Di macOS dan iOS, callback `stopURLSchemeTask` milik WebKit membatalkan konteks permintaan yang cocok, sehingga polling milik halaman yang telah berpindah segera berhenti menunggu. Registry menggunakan identitas task native yang dipertahankan sebagai kunci dan menghapus entri ketika pemrosesan permintaan menutupnya. Linux dan Windows masih belum menyediakan callback pembatalan dini yang setara dalam bridge saat ini; di sana, permintaan yang sedang menunggu tetap bertahan hingga periode tunggunya berakhir. Aturan 1 membuat *koneksi* segera ditutup dalam semua kasus. Jika tidak, pembatalan muncul sebagai `EPIPE` di Linux dan baru pada `Finish` di Windows.

## Pemilihan transport

`Stream(name)` memeriksa `window._wails.streamFactory`. Build server memasang factory yang mengembalikan `WebSocket` nyata; build webview membiarkannya tidak ditetapkan dan menggunakan klien polling.

Factory **harus** dipasang sebelum body modul mana pun dijalankan karena binding yang dihasilkan akan membuat stream dalam cakupan modul. `custom.js` tidak dapat melakukannya—`loadOptionalScript` menjalankan permintaan HEAD lalu menambahkan tag `<script>`, sehingga proses tersebut terjadi jauh terlambat. Sebagai gantinya, factory ditambahkan di awal bundle runtime saat bundle disajikan (`stream_prelude_server.go`), yang secara konstruksi bersifat sinkron: dependensi modul ES dievaluasi sebelum modul yang mengimpornya.

Jika Anda menambahkan transport ketiga, letakkan juga di prelude. Jangan tergoda untuk kembali menggunakan `custom.js`.

## Yang belum selesai

|  | status |
| --- | --- |
| Pembatalan permintaan dari lapisan platform | **Apple selesai; Linux/Windows tertunda**—lihat di atas |
| Konstanta buffer sebagai opsi | belum selesai; hanya pada waktu kompilasi |
| Stream bertipe | sengaja belum selesai—berdasarkan keputusan, frame berupa `[]byte` |
| Pipelining (polling kedua sedang berlangsung) | belum selesai; memerlukan perakitan ulang terurut di JS |
| Keadilan per koneksi | belum selesai—koneksi dalam satu jendela berbagi satu antrean, sehingga koneksi yang membanjiri antrean memperlambat koneksi di sekitarnya |
| Penggabungan frame JS→Go | **selesai**—frame yang terkumpul di belakang permintaan yang sedang berlangsung dikirim dalam batch terbatas; koneksi dengan beban ringan tetap mengirim satu frame per POST |
| Throughput Windows | ~100 MB/dtk, dibatasi oleh marshalling `WebResourceRequested`. Buffer bersama (`PostSharedBufferToScript`) adalah kandidat perbaikannya; binding tersedia di bawah `internal/webview2/pkg/webview2/`, tetapi belum dihubungkan ke `pkg/edge` |
| `wails3 dev` / Vite | **berfungsi**—diverifikasi dengan proyek `vanilla-js` yang dihasilkan: server pengembangan Vite melakukan proxy di `/`, dan `/wails/stream/*` dicocokkan oleh middleware server aset sebelum proxy, sehingga stream tidak terpengaruh |
| Multi-jendela | belum diuji dalam kondisi berbeban, meskipun secara konstruksi cakupan sesi dibatasi per jendela |

## Paket frontend dalam mode pengembangan

Proyek yang dihasilkan mengimpor `@wailsio/runtime` dari **npm**, bukan dari `/wails/runtime.js` yang dibundel dan disajikan oleh server aset. Di bawah `wails3 dev`, Vite me-resolve-nya dari `node_modules`, sehingga aplikasi yang dibuat menggunakan runtime yang telah dipublikasikan tidak akan melihat penambahan sisi klien yang dibuat pada suatu branch.

Selama streams belum dirilis, arahkan aplikasi pengujian ke paket dari working copy ini:

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

Tindakan tersebut terlebih dahulu membangun ulang `dist/`, sehingga sumber terkini selalu terinstal. Batalkan dengan menjalankan  
`npm install @wailsio/runtime@latest` di direktori yang sama.

Perhatikan bahwa ada **dua** keluaran build klien, dan mudah untuk membangun ulang salah satunya tetapi tidak yang  
lain: `task v3:runtime:build:package` menghasilkan `dist/` milik paket npm (yang diimpor oleh  
frontend aplikasi), sedangkan `task v3:runtime:build:assets` menghasilkan  
`bundledassets/runtime.js` (yang dimuat webview dari server aset). Perubahan pada  
`stream.ts` memerlukan keduanya.

## Pengujian

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

Pengujian urutan adalah yang terpenting: delapan goroutine mengirim secara bersamaan berdasarkan  
penghitung yang dibagikan saat lock antrean ditahan, dan urutan hasil pengurasan harus sama persis dengan urutan  
penerimaan. Jika pengujian ini sampai gagal, invarian penguras tunggal telah dilanggar.

Harness pengujian beban:

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

Di Windows, pengujian ini harus dijalankan dalam sesi konsol interaktif—pemanggilan SSH biasa berhenti di  
sesi 0 dengan keluaran sepanjang nol—dan berkas biner harus ditempatkan di lokasi yang dapat dibaca oleh akun SSH  
maupun akun konsol, karena ACL `C:\Users\<user>` membatasi akses hanya kepada pemiliknya.
