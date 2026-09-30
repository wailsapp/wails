---
title: "Memigrasikan WebSocket ke Stream"
description: "Konversi langkah demi langkah dari implementasi WebSocket yang ada ke stream Wails, termasuk perbedaan yang dapat menyebabkan kegagalan tanpa disadari"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

Panduan mekanis untuk mengonversi aplikasi yang saat ini menjalankan server WebSocket ke [Stream](/guides/streams/). Ditulis agar dapat diikuti secara harfiah, termasuk oleh agen.

Hasilnya adalah listener dapat dihapus: tidak ada port TCP yang diikat, pemeriksaan origin, token, ataupun hal lain yang dapat ditolak oleh firewall atau produk keamanan endpoint. API-nya cukup mirip sehingga sebagian besar kode frontend tidak perlu diubah — tetapi **tiga perbedaan dapat menyebabkan kegagalan tanpa disadari**, dan perbedaan tersebut dicantumkan lebih dahulu karena dapat menghabiskan waktu Anda sepanjang sore.

## Kasus umum: server HTTP lokal sebagai solusi sementara

Alasan umum aplikasi Wails memiliki WebSocket adalah karena sebelumnya tidak ada cara lain untuk mengirim umpan berkelanjutan ke frontend, sehingga aplikasi menjalankan `http.Server` sendiri pada port lokal dan frontend terhubung kembali ke server tersebut. Jika arsitektur aplikasi Anda seperti ini, migrasi ini akan menghapus server sepenuhnya — beserta beberapa hal yang Anda bangun *di sekitar* server tersebut.

**Mekanisme penemuan port tidak lagi diperlukan.** Sesuatu harus memberi tahu frontend port mana yang harus dihubungi: `GetServerPort()` yang diikat, port tetap dengan port cadangan jika port tersebut telah digunakan, global yang diinjeksi, atau nilai yang disimpan di `localStorage`. Semuanya tidak lagi diperlukan — stream dialamatkan berdasarkan nama, dan nama tersebut merupakan konstanta waktu kompilasi di kedua sisi.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**Konfigurasi CORS tidak lagi diperlukan.** Origin webview adalah `wails://` atau `http://wails.localhost`, tergantung platformnya, sehingga server lokal memerlukan `CheckOrigin`, header `Access-Control-Allow-Origin`, atau keduanya. Stream menggunakan server aset yang memuat halaman tersebut, sehingga tidak ada permintaan lintas origin yang perlu diizinkan.

**Token autentikasi apa pun yang Anda buat tidak lagi diperlukan.** Port yang diikat pada localhost dapat diakses oleh setiap proses di mesin, sehingga implementasi yang cermat menambahkan token atau nonce untuk mencegah perangkat lunak lain terhubung. Kini tidak ada lagi port yang dapat diakses.

**Endpoint non-WebSocket dipindahkan ke middleware server aset.** Server seperti ini jarang tetap murni — unduhan file, endpoint gambar, dan pemeriksaan kesehatan cenderung bertambah di samping socket. Stream tidak menggantikannya, tetapi Anda juga tidak memerlukan server kedua untuk endpoint tersebut. Pasang handler yang sama pada server aset:

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

Frontend kemudian memanggil `/api/...` sebagai URL relatif dengan origin yang sama — tanpa host, port, atau CORS. Dengan stream dan perubahan ini, server lokal tidak lagi memiliki tugas apa pun.

## Baca ini sebelum memulai

### 1. `ev.data` adalah `ArrayBuffer`, bukan string

Ini adalah perbedaan terpenting. WebSocket mengirimkan pesan teks sebagai string; stream mengirimkan setiap pesan sebagai byte. Kode seperti ini **dapat dikompilasi dan dijalankan, tetapi berperilaku keliru**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Perbaiki di batas sistem, bukan di setiap lokasi pemanggilan:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Atau, jika lalu lintas Anda berupa JSON — dan biasanya memang demikian — gunakan `JSONStream` sebagai pengganti `Stream` dan hindari masalah ini sepenuhnya:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

Itulah migrasi tersingkat untuk WebSocket JSON: ganti konstruktor, hapus pemanggilan `JSON.parse` dan `JSON.stringify`, lalu bagian handler lainnya tetap sama. Untuk lalu lintas non-JSON, cukup buat satu pembungkus dan biarkan setiap handler tetap utuh — lihat [shim kompatibilitas](#shim-kompatibilitas).

### 2. Pengiriman kompatibel; penerimaan tidak

`send()` menerima string dan mengodekannya sebagai UTF-8, sehingga `s.send(JSON.stringify(x))` tetap berfungsi tanpa perubahan. Hanya jalur penerimaan yang perlu diedit. Asimetri ini mudah terlewat karena separuh kode Anda tetap berfungsi.

### 3. Tidak ada URL

WebSocket membawa parameter koneksi dalam URL-nya — path, string kueri, subprotokol, dan token autentikasi. Stream hanya memiliki nama. Semua yang sebelumnya Anda teruskan dalam URL harus dipindahkan ke frame pertama, atau ke metode terikat yang dipanggil sebelum membuat koneksi.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Sisi Go

Hapus server HTTP, upgrader, dan registri koneksi. Masing-masing digantikan oleh handler.

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Stream |
| --- | --- |
| `upgrader.Upgrade` / `websocket.Accept` | *(tidak ada — `HandleStream` sudah mencakup seluruh pendaftaran)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`, atau cukup keluar dari handler |
| registri koneksi untuk broadcast | pertahankan registri Anda sendiri — lihat [Broadcast](#penyiaran) |
| `http.ListenAndServe`, mux, pemeriksaan origin, token | **hapus** |
| keepalive ping/pong | **hapus** — tidak ada socket menganggur yang harus dipertahankan tetap aktif |
| `r.Context()` | `c.Context()` |

Masa hidup goroutine handler sama dengan masa hidup koneksi, persis seperti handler gorilla, sehingga struktur loop yang ada dapat digunakan tanpa perubahan.

## Sisi frontend

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`, keempat konstanta status, `addEventListener`, `close(code, reason)`, dan `bufferedAmount` semuanya berperilaku sama seperti pada `WebSocket`.

### Shim kompatibilitas

Jika Anda memilih untuk tidak mengubah handler sama sekali, cukup bungkus konstruktornya sekali. Kode yang ada dan mengharapkan pesan string kemudian dapat digunakan tanpa perubahan:

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

Dengan demikian, `const ws = TextStream("feed")` dapat langsung menggantikan `new WebSocket(url)`.

## Penyiaran

Server WebSocket biasanya menyimpan registri agar dapat menyebarkan pesan. Stream tidak memiliki fitur penyiaran bawaan—pertahankan registrinya, tetapi simpan `*StreamConn` sebagai pengganti `*websocket.Conn`:

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

Dua aturan yang sebaiknya dipertahankan: lepaskan kunci sebelum mengirim, dan gunakan `TrySend` untuk penyebaran agar satu konsumen yang lambat tidak menghambat semua klien lainnya.

## Kasus lainnya: frontend berkomunikasi langsung dengan broker

Kasus ini lebih jarang terjadi dalam aplikasi Wails, tetapi penting untuk diketahui. Jika frontend membuka WebSocket **ke broker, bukan ke aplikasi Anda**—misalnya `nats.ws` ke server NATS atau MQTT melalui WebSocket—stream tidak dapat langsung menggantikannya karena stream menghubungkan frontend ke *kode Go Anda*, bukan ke pihak ketiga.

Migrasi ini merupakan perubahan arsitektur, dan biasanya merupakan perubahan yang baik:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Pindahkan klien broker ke Go, tempat pustaka native lebih baik daripada pustaka browser, lalu ekspos bagian yang dibutuhkan frontend melalui stream:

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

Perhatikan `TrySend` di dalam callback langganan: callback tersebut berjalan pada goroutine milik klien broker, dan memblokirnya akan menghambat pengiriman untuk setiap langganan pada koneksi itu.

Manfaat yang diperoleh: kredensial broker tidak pernah sampai ke frontend, tidak ada port WebSocket yang diekspos pada mesin, dan penyambungan kembali/backoff ditangani oleh klien Go yang matang, bukan oleh klien browser.

## Daftar periksa migrasi

- [ ] `HandleStream` didaftarkan untuk setiap endpoint WebSocket yang sebelumnya Anda miliki
- [ ] Loop pembacaan dikonversi: `ReadMessage`/`Read` → `c.Receive()`
- [ ] Operasi penulisan dikonversi: `WriteMessage` → `c.Send()`, atau `TrySend` dalam setiap operasi penyebaran atau callback broker
- [ ] Server HTTP, entri mux, upgrader, pemeriksaan origin, dan token autentikasi **dihapus**
- [ ] Keepalive ping/pong **dihapus**
- [ ] Parameter URL dipindahkan ke frame pertama atau metode terikat
- [ ] **Setiap hasil pembacaan `ev.data` didekode**—`new TextDecoder().decode(ev.data)`—atau shim digunakan
- [ ] Logika penyambungan kembali dipertahankan apa adanya (secara sengaja tidak ada penyambungan kembali bawaan)
- [ ] Klien broker dipindahkan ke Go jika frontend sebelumnya berkomunikasi langsung dengannya
- [ ] Periksa jendela `InitialHTML`—jendela tersebut sama sekali tidak dapat menggunakan stream

Hal-hal yang sebaiknya dicari dengan grep selama konversi:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Perbedaan perilaku yang perlu diperkirakan

|  | WebSocket | Stream |
| --- | --- | --- |
| Jenis pesan | teks atau biner | hanya byte |
| `ev.data` | string atau `Blob`/`ArrayBuffer` | selalu `ArrayBuffer` (kecuali jika `binaryType = "blob"`) |
| `binaryType` bawaan | `"blob"` | `"arraybuffer"` |
| Subprotokol, `extensions` | dinegosiasikan | tidak didukung; selalu `""` |
| Parameter koneksi | URL dan kueri | frame pertama atau panggilan terikat |
| Autentikasi | token atau cookie | tidak diperlukan—aplikasi adalah satu-satunya pemanggil |
| Keepalive | ping/pong | tidak diperlukan |
| Penyambungan kembali otomatis | tidak ada | tidak ada (sama) |
| Kode penutupan | rentang lengkap | `1000` normal, `1001` sesi ditutup, `1002` ketidakcocokan framing, `1006` kesalahan |
| Tekanan balik | buffer soket kernel | 8 MB / 256 frame per jendela, lalu `Send` memblokir |
| Banyak koneksi, satu endpoint | ya | ya |

## Setelah migrasi

Pemeriksaan kewajaran untuk mendeteksi kesalahan umum:

1. Muat ulang halaman berulang kali — handler seharusnya berhenti dan handler baru dimulai setiap kali, tanpa pernah menumpuk.
2. Kirim pesan yang lebih besar dari 512 KB ke setiap arah.
3. Biarkan tidak aktif selama lebih dari satu menit; lalu lintas seharusnya kembali berjalan tanpa menyambung ulang.
4. Buka devtools dan jeda pada breakpoint saat ada beban, lalu lanjutkan — produsen seharusnya menunggu dalam keadaan terblokir lalu pulih, bukan kehilangan data atau bertambah tanpa batas.
