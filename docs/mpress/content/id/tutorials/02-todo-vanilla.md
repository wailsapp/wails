---
title: "Daftar TODO"
description: "Buat aplikasi daftar TODO lengkap dengan operasi CRUD"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

Dalam tutorial ini, Anda akan membuat aplikasi daftar TODO yang berfungsi sepenuhnya. Tutorial ini merupakan kelanjutan dari tutorial Layanan Kode QR. Anda akan mempelajari cara mengelola state, menangani beberapa operasi, dan membuat antarmuka pengguna yang dipoles.

**Yang akan Anda buat:**

- Aplikasi TODO lengkap dengan fungsi untuk menambahkan, menyelesaikan, dan menghapus item
- Pengelolaan state yang aman untuk thread (penting bagi aplikasi desktop)
- Desain UI modern bergaya glassmorphism
- Semuanya menggunakan JavaScript murni—tanpa memerlukan framework

**Yang akan Anda pelajari:**

- Operasi CRUD (Buat, Baca, Perbarui, Hapus)
- Mengelola state yang dapat berubah secara aman di Go
- Menangani input pengguna dan validasi
- Membuat UI responsif yang terasa seperti aplikasi native

![Aplikasi Daftar TODO](/assets/todo-app.png)

**Waktu penyelesaian:** 20 menit

## Buat Proyek Anda

@steps
### Buat proyek
Pertama, buat proyek Wails baru. Kita akan menggunakan templat vanilla bawaan agar memiliki titik awal yang bersih:

```bash
wails3 init -n todo-app -t vanilla-js
cd todo-app
```

Perintah ini membuat proyek baru dengan struktur dasar: backend Go di direktori root dan kode frontend di direktori `frontend/`.

### Buat layanan TODO
Layanan TODO akan mengelola state aplikasi dan menyediakan metode untuk operasi CRUD. Tidak seperti server web yang mengisolasi setiap permintaan, aplikasi desktop dapat menjalankan beberapa operasi secara bersamaan sehingga kita memerlukan pengelolaan state yang aman untuk thread.

Hapus `greetservice.go` dan buat file baru `todoservice.go`:

```go {title="todoservice.go"}
package main

import (
    "errors"
    "sync"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

type TodoService struct {
    todos  []Todo
    nextID int
    mu     sync.RWMutex
}

func NewTodoService() *TodoService {
    return &TodoService{
        todos:  []Todo{},
        nextID: 1,
    }
}

func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }

    t.mu.Lock()
    defer t.mu.Unlock()

    todo := Todo{
        ID:        t.nextID,
        Title:     title,
        Completed: false,
    }
    t.todos = append(t.todos, todo)
    t.nextID++

    return &todo, nil
}

func (t *TodoService) Toggle(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos[i].Completed = !t.todos[i].Completed
            return nil
        }
    }
    return errors.New("todo not found")
}

func (t *TodoService) Delete(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i, todo := range t.todos {
        if todo.ID == id {
            t.todos = append(t.todos[:i], t.todos[i+1:]...)
            return nil
        }
    }
    return errors.New("todo not found")
}
```

**Yang terjadi di sini:**

**Struct `Todo`:**

- Menentukan struktur data kita dengan field ID, Title, dan Completed
- Tag `json:` memberi tahu Go cara mengonversi struct ini menjadi JSON untuk frontend
- Setiap field diekspor (diawali huruf kapital) agar dapat dikenali oleh generator binding

**Struct `TodoService`:**

- `todos []Todo`—slice yang menyimpan semua item TODO kita
- `nextID int`—melacak ID berikutnya yang akan ditetapkan (menyimulasikan peningkatan otomatis)
- `mu sync.RWMutex`—mutex baca/tulis untuk akses yang aman bagi thread

**Keamanan thread dengan `sync.RWMutex`:**

- Aplikasi desktop dapat menerima beberapa operasi serentak dari UI
- `RLock()` memungkinkan beberapa pembaca mengakses secara bersamaan (misalnya, beberapa panggilan `GetAll`)
- `Lock()` memberikan akses eksklusif untuk operasi tulis (misalnya, `Add`, `Toggle`, dan `Delete`)
- `defer` memastikan lock dilepaskan meskipun fungsi kembali lebih awal atau mengalami panic

**Metode-metodenya:**

- `GetAll()`—mengembalikan semua TODO (menggunakan read lock karena kita tidak mengubah data)
- `Add(title)`—membuat TODO baru, memvalidasi input, dan menaikkan ID
- `Toggle(id)`—membalik status selesai suatu TODO
- `Delete(id)`—menghapus TODO dari slice

**Penanganan kesalahan:**

- Kita mengembalikan `error` sebagai nilai terakhir sesuai konvensi Go
- Judul kosong ditolak
- Operasi pada TODO yang tidak ada akan mengembalikan kesalahan
- Kesalahan ini menjadi exception JavaScript di frontend

### Perbarui main.go
Daftarkan layanan TODO pada aplikasi Wails Anda. Temukan bagian `Services` dalam `main.go`, lalu ganti GreetService dengan TodoService kita:

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**Yang terjadi di sini:**

- Kita menghapus GreetService bawaan dan menggantinya dengan TodoService
- `application.NewService()` membungkus layanan kita agar dapat dikelola oleh Wails
- Wails akan otomatis menghasilkan binding JavaScript untuk semua metode publik pada layanan ini

### Buat UI frontend
Sekarang, mari buat frontend. Di sinilah kita akan memanggil metode Go dan menampilkan UI. Kita menggunakan JavaScript murni agar semuanya tetap sederhana dan untuk menunjukkan secara langsung cara kerja binding.

Ganti `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {TodoService} from "../bindings/changeme";

async function loadTodos() {
    const todos = await TodoService.GetAll();
    const list = document.getElementById('todo-list');

    list.innerHTML = todos.map(todo => `
        <div class="todo ${todo.completed ? 'completed' : ''}">
            <input type="checkbox"
                   ${todo.completed ? 'checked' : ''}
                   onchange="toggleTodo(${todo.id})">
            <span>${todo.title}</span>
            <button onclick="deleteTodo(${todo.id})">Delete</button>
        </div>
    `).join('');
}

window.addTodo = async () => {
    const input = document.getElementById('todo-input');
    const title = input.value.trim();

    if (title) {
        await TodoService.Add(title);
        input.value = '';
        await loadTodos();
    }
}

window.toggleTodo = async (id) => {
    await TodoService.Toggle(id);
    await loadTodos();
}

window.deleteTodo = async (id) => {
    await TodoService.Delete(id);
    await loadTodos();
}

// Load todos on startup
loadTodos();
```

**Yang terjadi di sini:**

**Mengimpor binding:**

- `import {TodoService} from "../bindings/changeme"`—mengimpor binding Go yang dibuat secara otomatis
- Catatan: `changeme` akan menjadi nama modul Anda yang sebenarnya dari `go.mod`

**Fungsi `loadTodos()`:**

- Memanggil `TodoService.GetAll()` untuk mengambil semua TODO dari Go
- Menyusun HTML untuk setiap TODO menggunakan literal templat
- Menambahkan atau menghapus class `completed` secara dinamis untuk penerapan gaya
- Menggunakan atribut `onclick` untuk menghubungkan tombol dengan fungsi kita
- Menggabungkan seluruh HTML dan menyisipkannya ke dalam DOM

**Fungsi CRUD:**

- `addTodo()`—memvalidasi input, memanggil metode Go `Add`, lalu memuat ulang daftar
- `toggleTodo(id)`—memanggil metode `Toggle` milik Go, lalu memuat ulang daftar
- `deleteTodo(id)`—memanggil metode `Delete` milik Go, lalu memuat ulang daftar
- Semua fungsi bersifat asinkron karena panggilan Go mengembalikan Promise

**Alasan melampirkannya ke window:**

- `window.addTodo = ...` membuat fungsi dapat diakses dari atribut HTML `onclick`
- Ini adalah pola sederhana untuk JavaScript murni (framework menanganinya dengan cara berbeda)
- Dalam produksi, Anda dapat menggunakan delegasi event yang tepat sebagai gantinya

**Pola pemuatan ulang:**

- Setelah setiap mutasi (tambah/ubah status/hapus), kita memanggil `loadTodos()` lagi
- Hal ini memastikan UI tetap sinkron dengan state Go
- Alternatif: Buat metode Go mengembalikan state baru untuk menghindari panggilan kedua

### Perbarui HTML
HTML menyediakan struktur untuk aplikasi TODO kita. Strukturnya minimal dan semantik—keajaibannya terjadi di JavaScript dan CSS.

Ganti `frontend/index.html`:

```html {title="frontend/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>TODO App</title>
    <link rel="stylesheet" href="./style.css"/>
</head>
<body>
    <div class="container">
        <h1>My TODOs</h1>
        <div class="card">
            <div class="input-box">
                <input type="text"
                       id="todo-input"
                       class="input"
                       placeholder="Add a new todo..."
                       onkeypress="if(event.key==='Enter') addTodo()">
                <button class="btn" onclick="addTodo()">Add</button>
            </div>
            <div id="todo-list"></div>
        </div>
    </div>
    <script type="module" src="./src/main.js"></script>
</body>
</html>
```

**Yang terjadi di sini:**

**Strukturnya:**

- `container`—memusatkan aplikasi kita dan membatasi lebarnya
- `card`—kartu putih utama yang memuat semuanya
- `input-box`—kontainer flex untuk kolom input dan tombol Tambah
- `todo-list`—tempat setiap todo akan disisipkan oleh JavaScript

**Penanganan peristiwa:**

- `onkeypress="if(event.key==='Enter') addTodo()"`—menambahkan todo saat tombol Enter ditekan
- `onclick="addTodo()"`—menambahkan todo saat tombol diklik
- Penangan peristiwa inline cocok untuk aplikasi JavaScript murni yang sederhana

**Skrip modul:**

- `<script type="module">` memungkinkan kita menggunakan impor ES6
- File `main.js` kita dapat mengimpor binding dan menggunakan JavaScript modern

### Atur gaya aplikasi
CSS menghasilkan desain glassmorphism modern dengan transisi yang halus. Kita ingin menciptakan tampilan rapi yang membuat aplikasi menyenangkan untuk digunakan.

Ganti `frontend/public/style.css`:

```css {title="frontend/public/style.css"}
:root {
    font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
    "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
    font-size: 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: rgba(255, 255, 255, 0.87);
}

body {
    margin: 0;
    display: flex;
    place-items: center;
    justify-content: center;
    min-height: 100vh;
}

.container {
    width: 100%;
    max-width: 600px;
    padding: 20px;
}

h1 {
    text-align: center;
    color: white;
    font-size: 2.5em;
    font-weight: 300;
    margin: 0 0 30px 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 30px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.input-box {
    display: flex;
    gap: 10px;
    margin-bottom: 25px;
}

.input {
    flex: 1;
    border: 2px solid #e0e0e0;
    border-radius: 12px;
    height: 50px;
    padding: 0 20px;
    font-size: 16px;
    transition: all 0.3s ease;
}

.input:focus {
    border-color: #667eea;
    outline: none;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.btn {
    height: 50px;
    padding: 0 30px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

#todo-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.todo {
    display: flex;
    align-items: center;
    padding: 18px 20px;
    background: white;
    border: 2px solid #f0f0f0;
    border-radius: 12px;
    transition: all 0.3s ease;
    gap: 15px;
}

.todo:hover {
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
    transform: translateX(4px);
}

.todo.completed {
    opacity: 0.6;
}

.todo.completed span {
    text-decoration: line-through;
    color: #999;
}

.todo input[type="checkbox"] {
    width: 24px;
    height: 24px;
    cursor: pointer;
    appearance: none;
    -webkit-appearance: none;
    border: 2px solid #667eea;
    border-radius: 6px;
    position: relative;
    transition: all 0.3s ease;
    flex-shrink: 0;
}

.todo input[type="checkbox"]:hover {
    background: rgba(102, 126, 234, 0.1);
}

.todo input[type="checkbox"]:checked {
    background: #667eea;
    border-color: #667eea;
}

.todo input[type="checkbox"]:checked::after {
    content: '✓';
    position: absolute;
    color: white;
    font-size: 16px;
    font-weight: bold;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

.todo span {
    flex: 1;
    font-size: 16px;
    color: #333;
}

.todo button {
    padding: 8px 16px;
    background: #ff4757;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    opacity: 0;
    flex-shrink: 0;
}

.todo:hover button {
    opacity: 1;
}

.todo button:hover {
    background: #ee5a6f;
    transform: scale(1.05);
}

#todo-list:empty::before {
    content: "No todos yet. Add one above!";
    display: block;
    text-align: center;
    padding: 40px 20px;
    color: #999;
}
```

**Yang terjadi di sini:**

**Desain glassmorphism:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)`—latar belakang dengan gradien ungu
- `backdrop-filter: blur(10px)`—menghasilkan efek kaca buram pada kartu
- `rgba(255, 255, 255, 0.95)`—warna putih semitransparan untuk efek kaca

**Gaya kotak centang khusus:**

- `appearance: none` menghapus kotak centang bawaan browser
- Kita membuat kotak bersudut bundar khusus dengan tanda centang menggunakan `::after`
- Tanda centang muncul saat `checked` dengan menggunakan karakter Unicode ✓

**Interaksi saat penunjuk diarahkan:**

- Todo bergeser ke kanan saat penunjuk diarahkan ke atasnya (`transform: translateX(4px)`)
- Tombol hapus disembunyikan hingga penunjuk diarahkan ke atasnya (`opacity: 0` → `opacity: 1`)
- Tombol sedikit membesar saat penunjuk diarahkan ke atasnya untuk memberikan umpan balik taktil

**Keadaan kosong:**

- `#todo-list:empty::before` menampilkan pesan saat tidak ada todo
- Solusi yang hanya menggunakan CSS—JavaScript tidak diperlukan

### Jalankan aplikasi
Mari kita lihat cara kerjanya! Jalankan server pengembangan:

```bash
wails3 dev
```

Aplikasi akan dikompilasi dan dibuka. Cobalah:

- Ketik todo, lalu tekan Enter atau klik Tambah
- Klik kotak centang untuk menandainya sebagai selesai
- Arahkan penunjuk ke todo untuk menampilkan tombol hapus
- Perhatikan bahwa UI langsung diperbarui—ini menunjukkan pola penyegaran kita berfungsi

**Yang terjadi:**

- Wails secara otomatis menghasilkan binding untuk metode TodoService Anda
- Mode pengembangan menyertakan hot reload—coba ubah CSS dan perhatikan pembaruannya
- Kode Go Anda berjalan secara native—tidak memerlukan penerjemahan atau interpretasi

@end

## Cara Kerjanya

### Pengelolaan State yang Aman untuk Thread

`sync.RWMutex` menyediakan akses bersamaan yang aman:

```go
func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()  // Read lock - multiple readers allowed
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    t.mu.Lock()  // Write lock - exclusive access
    defer t.mu.Unlock()
    // ... mutations
}
```

**Mengapa hal ini penting:**

- Beberapa panggilan frontend dapat berlangsung secara bersamaan
- Operasi baca tidak saling memblokir
- Operasi tulis memperoleh akses eksklusif
- `defer` memastikan penguncian selalu dilepas

### Penanganan Kesalahan

Layanan mengembalikan kesalahan untuk operasi yang tidak valid:

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

Di frontend, Anda dapat menangkap kesalahan tersebut:

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### Sinkronisasi State

Setelah setiap mutasi, kita memuat ulang seluruh daftar:

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**Pendekatan alternatif:** Kembalikan daftar yang telah diperbarui dari setiap metode agar panggilan kedua tidak diperlukan.

## Penyempurnaan

### Tambahkan Statistik

Tambahkan ini ke `todoservice.go`:

```go
type TodoStats struct {
    Total     int `json:"total"`
    Completed int `json:"completed"`
    Active    int `json:"active"`
}

func (t *TodoService) GetStats() TodoStats {
    t.mu.RLock()
    defer t.mu.RUnlock()

    stats := TodoStats{
        Total: len(t.todos),
    }

    for _, todo := range t.todos {
        if todo.Completed {
            stats.Completed++
        } else {
            stats.Active++
        }
    }

    return stats
}
```

Tampilkan di frontend:

```javascript
async function loadTodos() {
    const [todos, stats] = await Promise.all([
        TodoService.GetAll(),
        TodoService.GetStats()
    ]);

    // Display stats
    document.getElementById('stats').textContent =
        `${stats.active} active, ${stats.completed} completed`;

    // ... render todos
}
```

### Tambahkan "Hapus yang Selesai"

```go
func (t *TodoService) ClearCompleted() int {
    t.mu.Lock()
    defer t.mu.Unlock()

    removed := 0
    newTodos := []Todo{}

    for _, todo := range t.todos {
        if !todo.Completed {
            newTodos = append(newTodos, todo)
        } else {
            removed++
        }
    }

    t.todos = newTodos
    return removed
}
```

### Tambahkan Persistensi

Untuk aplikasi produksi, tambahkan persistensi menggunakan solusi penyimpanan yang paling sesuai dengan aplikasi Anda, seperti SQLite, PostgreSQL, atau layanan terkelola.

## Build untuk Produksi

Saat siap mendistribusikan aplikasi TODO, buat build produksi:

```bash
wails3 build
```

Tindakan ini menghasilkan executable native yang dioptimalkan di `bin/`:

- Mengompilasi kode Go Anda dengan pengoptimalan
- Membuat build produksi frontend Anda
- Mengemas semuanya ke dalam satu executable
- Ukuran aplikasi yang dihasilkan biasanya 10-20MB (bandingkan dengan Electron yang berukuran 150MB+)

Anda dapat langsung menjalankan executable tersebut—tidak memerlukan runtime dan tidak ada server yang harus dijalankan. Ini benar-benar aplikasi native.

## Yang Telah Anda Buat

Anda baru saja membuat aplikasi TODO lengkap dengan:

**Implementasi CRUD lengkap:**

- Membuat layanan dengan operasi Create, Read, Update, dan Delete
- Menambahkan validasi input dan penanganan kesalahan
- Mempelajari cara kesalahan Go menjadi pengecualian JavaScript

**Pengelolaan state yang aman untuk thread:**

- Menggunakan `sync.RWMutex` untuk menangani akses bersamaan dengan aman
- Memahami perbedaan antara penguncian baca (RLock) dan penguncian tulis (Lock)
- Melihat cara `defer` mencegah deadlock dengan menjamin pembersihan

**UI modern dan rapi:**

- Membuat antarmuka glassmorphism dengan gradien dan efek buram
- Membuat kotak centang bergaya khusus tanpa framework apa pun
- Menambahkan interaksi saat penunjuk diarahkan dan transisi untuk memberikan nuansa native
- Menerapkan keadaan kosong dengan CSS murni

**Dasar-dasar Wails:**

- Pendaftaran layanan dan pembuatan binding otomatis
- Memanggil metode Go dari JavaScript dengan async/await
- Sinkronisasi state antara Go dan frontend
- Membangun dan mengemas aplikasi desktop native

## Langkah Berikutnya

Setelah memahami operasi CRUD dan pengelolaan state, cobalah:

- **Tambahkan persistensi:** Pertahankan data tugas agar tetap tersedia setelah aplikasi dimulai ulang dengan SQLite atau solusi penyimpanan lainnya
- **Tambahkan lebih banyak fitur:** Pemfilteran (semua/aktif/selesai), pengeditan tugas yang sudah ada, operasi massal
- **Pelajari tutorial Notes:** Lihat cara kerja operasi file di [Notes](/tutorials/03-notes-vanilla/)
- **Bangun sesuatu yang nyata:** Gunakan konsep-konsep ini untuk membangun aplikasi Anda sendiri!
