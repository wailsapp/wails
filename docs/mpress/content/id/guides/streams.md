---
title: "Stream"
description: "Stream byte dua arah antara Go dan JavaScript, dengan model pemrograman WebSocket tanpa soket yang mendengarkan"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

Stream menyediakan kanal byte dua arah bernama dan berurutan antara Go dan frontend Anda, dengan model pemrograman yang sama seperti WebSocket — **tanpa mengikat port TCP**.

WebSocket tidak dapat digunakan melalui skema URL khusus, sehingga satu-satunya cara untuk menggunakannya di dalam webview adalah dengan menjalankan server HTTP sungguhan dan mendengarkan pada suatu port. Dalam aplikasi desktop, ini berarti ada port lokal terbuka yang dapat dijangkau oleh proses lain mana pun di mesin tersebut, sehingga memerlukan pemeriksaan origin dan token agar aman, serta terlihat oleh setiap produk firewall dan keamanan endpoint yang digunakan pengguna Anda. Stream menghindari semua itu: stream menggunakan server aset yang sudah disediakan oleh aplikasi Anda, yang sudah terikat pada origin.

Sedang memigrasikan implementasi WebSocket yang ada? Ikuti [Memigrasikan WebSocket ke Stream](/guides/streams-from-websockets/) — panduan tersebut dirancang agar dapat diterapkan secara mekanis dan diawali dengan tiga perbedaan yang dapat menyebabkan kerusakan tanpa gejala yang jelas.

## Mulai cepat

Deklarasikan stream di Go. Handler dijalankan sekali untuk setiap koneksi, dalam goroutine tersendiri:

```go
app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        _ = c.Send(process(frame))  // blocks like a socket write
    }
})
```

Hubungkan dari frontend berdasarkan nama. Objek tersebut mengimplementasikan antarmuka `WebSocket`:

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` mengembalikan **secara sinkron** dengan `readyState === CONNECTING`, persis seperti `new WebSocket(url)`, sehingga Anda dapat membuatnya dalam lingkup modul:

```js
export const Telemetry = Stream("telemetry");
```

## Frame berupa byte

Setiap frame berupa `[]byte` di Go dan `ArrayBuffer` di JavaScript. Tidak ada skema atau pengodean yang dipaksakan kepada Anda — gunakan JSON, protobuf, CBOR, atau byte mentah sesuai kebutuhan.

Frame adalah **pesan, bukan stream byte**: frame tiba secara utuh atau tidak sama sekali, dan panjangnya dikirim bersamanya. Kedua sisi tidak perlu mengetahui ukurannya terlebih dahulu, sehingga struct dengan bidang `[]byte` akan dimarshal menjadi hasil apa pun yang dihasilkan proses marshal tersebut dan dikirim sebagai satu frame.

## Mengirim objek

Frame berupa byte, tetapi Anda jarang perlu memikirkannya sebagai byte. Kedua sisi menyediakan kemudahan JSON yang saling berpasangan:

```go
type Reading struct {
    Sensor string  `json:"sensor"`
    Value  float64 `json:"value"`
}

app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    var cmd map[string]any
    if err := c.ReceiveJSON(&cmd); err != nil {
        return
    }

    _ = c.SendJSON(Reading{Sensor: "cpu", Value: 42.5})
})
```

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("telemetry");
s.onopen    = () => s.send({ subscribe: "cpu" });   // stringified for you
s.onmessage = (ev) => console.log(ev.data.value);   // already an object
```

`JSONStream` adalah objek yang sama dengan `Stream`, dengan pengodean dilakukan pada batas sistem — tidak ada protokol terpisah, dan handler Go tidak dapat membedakannya. Frame yang bukan JSON valid memicu peristiwa `error` dan dibuang, bukan memutus koneksi.

Gunakan `Stream` biasa saat Anda menginginkan byte: protobuf, CBOR, format biner, atau apa pun yang lebih Anda pilih untuk dikodekan sendiri.

## API Go

```go
// Register a handler. Runs once per connection, on its own goroutine.
func (a *App) HandleStream(name string, handler StreamHandler)

// The connection.
func (c *StreamConn) Send(data []byte) error        // blocks when the buffer is full
func (c *StreamConn) TrySend(data []byte) error     // ErrStreamFull instead of blocking
func (c *StreamConn) Receive() ([]byte, error)      // blocks until a frame or close
func (c *StreamConn) SendJSON(v any) error          // marshal and send as one frame
func (c *StreamConn) ReceiveJSON(v any) error       // receive one frame and unmarshal
func (c *StreamConn) Context() context.Context      // cancelled on disconnect
func (c *StreamConn) Window() Window                // nil in server mode
func (c *StreamConn) Name() string
func (c *StreamConn) Close() error
```

**Masa aktif goroutine handler sama dengan masa aktif koneksi.** Kembali dari handler akan menutup koneksi, jadi blokir pada `Receive` (atau pada `c.Context()`) selama Anda ingin koneksi tetap terbuka. Polanya sama seperti handler WebSocket `gorilla`/`coder`.

Error yang mungkin terjadi adalah `ErrStreamClosed` (peer sudah tidak tersedia) dan `ErrStreamFull` (hanya dari `TrySend`).

## API JavaScript

`Stream(name)` mengembalikan objek yang mengimplementasikan subset berguna dari `WebSocket`:

| didukung | catatan |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`, `onmessage`, `onclose`, `onerror` | ditambah `addEventListener` |
| `send(data)` | string, `ArrayBuffer`, typed array, atau `Blob` — lihat kepemilikan di bawah |
| `JSONStream(name)` | objek yang sama, objek sebagai masukan dan keluaran |
| `close(code, reason)` |  |
| `binaryType` | **secara default bernilai `"arraybuffer"`**, bukan `"blob"` |
| `bufferedAmount` | byte yang diantrekan oleh `send` dan belum mencapai Go |
| `protocol`, `extensions` | selalu `""` — tidak dinegosiasikan |

Nilai default `binaryType` adalah satu-satunya penyimpangan yang disengaja dari standar: frame selalu berupa biner, dan `Blob` akan memaksakan satu tahap asinkron tambahan untuk membaca setiap pesan. Atur nilainya menjadi `"blob"` jika Anda menginginkan perilaku standar.

Beberapa koneksi ke nama stream yang sama diperbolehkan, baik dari satu maupun beberapa jendela. Masing-masing mendapatkan `StreamConn` dan goroutine handler tersendiri.

**Kepemilikan buffer berbeda menurut arah.** `send()` JavaScript mengambil snapshot masukan biner yang dapat diubah secara sinkron, sesuai dengan perilaku WebSocket native, sehingga pemanggil dapat menggunakannya kembali segera setelah `send()` kembali. `Send` Go mengalihkan kepemilikan slice-nya ke transport dan tidak menyalinnya; jangan mengubah atau menggunakan kembali slice tersebut setelah pemanggilan berhasil. Serahkan slice baru jika produsen perlu menggunakan kembali penyimpanannya.

## Siklus hidup

Stream berperilaku seperti soket, dan peristiwa yang menutup soket juga menutup stream:

| peristiwa | yang terjadi |
| --- | --- |
| Pemuatan ulang atau navigasi halaman | koneksi ditutup, `Receive` milik handler mengembalikan error, halaman baru membuat koneksi baru |
| `window.close()` / jendela dimusnahkan | semua koneksi milik jendela tersebut ditutup |
| `s.close()` di JS | `Receive` milik handler mengembalikan `ErrStreamClosed` |
| Handler kembali | frontend menerima `onclose` |
| Aplikasi dimatikan | context setiap koneksi dibatalkan |

Tidak ada **koneksi ulang otomatis**, sesuai dengan `WebSocket`. Jika aplikasi Anda memerlukannya, logika koneksi ulang yang sudah Anda miliki untuk WebSocket akan berfungsi tanpa perubahan — buat ulang stream di `onclose`.

## Backpressure

`Send` akan memblokir ketika frontend tidak mampu mengimbangi, sebagaimana penulisan ke soket memblokir saat buffer pengiriman penuh. Gunakan `TrySend` jika Anda lebih memilih membuang data daripada menunggu:

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

Frontend yang dijeda — akibat breakpoint devtools, jendela tersembunyi, atau App Nap — berhenti mengumpulkan data, lalu batas buffer memblokir produsen. Ini disengaja: batas tersebut membatasi penggunaan memori agar stream yang tidak dibaca tidak terus bertambah tanpa batas.

## Mode server

Membangun dengan `-tags server` mengganti transport dengan **WebSocket sungguhan** di `/wails/stream/ws`, karena mode server sudah memiliki listener yang mendukung upgrade koneksi HTTP ke WebSocket. Handler Go dan kode frontend tetap sama—tidak ada yang berubah dalam aplikasi Anda. Runtime memilih transport untuk Anda sebelum kode modul apa pun dijalankan. Secara default, koneksi WebSocket berasal dari origin yang sama. Server yang sengaja menghosting frontend-nya di origin tepercaya lain dapat menambahkan host tersebut dengan `ServerOptions.WebSocketOriginPatterns`.

## Performa

Diukur dengan `v3/tests/stream-performance`, tanpa pembatasan kecepatan, terdapat 0 frame yang hilang dan 0 frame yang berubah urutan di antara ~41 juta frame:

|  | Puncak Go→JS | Puncak JS→Go |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/dtk | 727 MB/dtk |
| Windows / WebView2 | 100 MB/dtk | 99 MB/dtk |

Karakteristiknya lebih penting daripada nilai puncaknya:

- **Go→JS jauh lebih cepat untuk frame kecil** — 634000 frame/dtk di macOS dibandingkan dengan ~6200/dtk untuk arah sebaliknya. Satu respons menggabungkan hingga 256 frame; JS→Go juga mengelompokkan frame yang terakumulasi di belakang permintaan yang sedang berlangsung, tetapi setiap koneksi tetap menserialkan rantai POST-nya sendiri. Jika Anda mengirim banyak pesan kecil, pilih Go→JS atau kelompokkan pesan tersebut pada tingkat aplikasi sebelum mengirimkannya ke Go.
- **Di Windows, 512 KB merupakan ukuran optimal untuk unggahan.** Frame yang lebih besar akan dipecah menjadi beberapa permintaan, dan hasil pengukuran menunjukkan bahwa frame berukuran 4 MB *lebih lambat* daripada frame berukuran 512 KB.
- **Latensinya rendah dan tetap rendah**: p99 sekitar 1–2 ms di macOS, dan tidak memburuk saat beban meningkat — p99 yang diukur pada 20000 frame/dtk *lebih rendah* daripada pada 100 frame/dtk.

Tabel lengkap untuk setiap platform dan metode pengukurannya tersedia dalam catatan pengukuran yang menyertai fitur ini.

## Batas

|  | batas | yang terjadi ketika batas tercapai |
| --- | --- | --- |
| Di-buffer per jendela, menunggu diambil | 8 MB atau 256 frame, mana pun yang tercapai lebih dahulu | `Send` memblokir; `TrySend` mengembalikan `ErrStreamFull` |
| Di-buffer di seluruh aplikasi, menunggu diambil/ditulis | 256 MB atau 8192 frame data | sama |
| Diterima per koneksi, menunggu `Receive` | 8 MB atau 256 frame | `send()` frontend dicoba ulang secara otomatis hingga handler dapat mengejar ketertinggalan |
| Diterima di seluruh aplikasi, menunggu `Receive` | 256 MB atau 8192 frame | sama |
| Koneksi per jendela | 256 | pembukaan dicoba ulang secara otomatis hingga tersedia slot |
| Koneksi aktif di seluruh aplikasi | 4096 | sama |
| Sesi per jendela | 16 | pemuatan ulang menggantikan sesi lamanya sendiri; jika tidak, pembukaan dicoba ulang |
| Satu frame, dalam kedua arah | 64 MB | Go mengembalikan `ErrStreamTooLarge`; dari JS, stream memunculkan `error` lalu ditutup |
| Nama stream | 256 byte UTF-8 | pembukaan ditolak dan stream memunculkan `error` |
| Waktu tunggu polling saat tidak aktif | 20 dtk | polling mengembalikan hasil kosong dan runtime segera menjalankannya kembali |

Tidak satu pun batas di atas menyebabkan data hilang secara diam-diam. Dua baris yang menyatakan *dicoba ulang secara otomatis* merupakan backpressure biasa — runtime menahan frame dan mencoba kembali setelah jeda singkat, sehingga kode Anda mengalami stream yang lebih lambat, bukan error. Baris yang memunculkan `error` menunjukkan kesalahan pemrograman, bukan beban, dan kesalahan tersebut ditampilkan alih-alih disembunyikan.

Saat ini, nilai-nilai tersebut merupakan konstanta waktu kompilasi, bukan opsi. Lihat panduan internal jika Anda perlu mengubahnya.

## Kapan tidak menggunakan stream

- **Untuk permintaan/respons, gunakan binding.** Stream ditujukan untuk data berkelanjutan atau data yang tidak diminta; panggilan yang mengembalikan nilai lebih sederhana jika dibuat sebagai metode terikat.
- **Untuk peristiwa aplikasi, gunakan `Emit`/`On`.** Peristiwa disebarkan ke setiap listener dan menggunakan sistem terpisah yang telah teruji. Stream bersifat titik-ke-titik.
- **Tidak tersedia untuk jendela `InitialHTML`.** Jendela tersebut dimuat dengan `origin === "null"`, sehingga sama sekali tidak dapat menjangkau server aset.
