---
title: "Runtime Frontend"
description: "Paket runtime JavaScript Wails untuk integrasi frontend"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Runtime frontend Wails adalah pustaka standar untuk aplikasi Wails. Runtime ini menyediakan sejumlah fitur yang dapat digunakan dalam aplikasi Anda, termasuk:

- Pengelolaan jendela
- Dialog
- Integrasi browser
- Papan klip
- Menu
- Informasi sistem
- Peristiwa
- Menu konteks
- Layar
- WML (Bahasa Markah Wails)

Runtime diperlukan untuk integrasi antara Go dan frontend. Ada 2 cara untuk mengintegrasikan runtime:

- Menggunakan paket `@wailsio/runtime`
- Menggunakan bundel siap pakai

## Menggunakan paket npm

Paket `@wailsio/runtime` adalah paket JavaScript yang menyediakan akses ke runtime Wails dari frontend. Paket ini digunakan oleh semua templat standar dan merupakan cara yang disarankan untuk mengintegrasikan runtime ke dalam aplikasi Anda. Dengan menggunakan paket `@wailsio/runtime`, Anda hanya akan menyertakan bagian runtime yang digunakan.

Paket ini tersedia di npm dan dapat diinstal menggunakan:

```shell
npm install --save @wailsio/runtime
```

## Menggunakan bundel siap pakai

Beberapa proyek tidak menggunakan bundler JavaScript dan mungkin lebih memilih versi runtime dalam bentuk bundel siap pakai. Versi ini dapat dibuat secara lokal menggunakan perintah berikut:

```shell
wails3 generate runtime
```

Perintah tersebut akan menghasilkan file `runtime.js` (dan `runtime.debug.js`) di direktori saat ini. File ini merupakan modul ES yang dapat diimpor oleh skrip aplikasi Anda seperti paket npm, tetapi API-nya juga diekspor ke objek window global. Jadi, untuk aplikasi yang lebih sederhana, Anda dapat menggunakannya sebagai berikut:

```html
<html>
    <head>
        <script type="module" src="./runtime.js"></script>
        <script>
            window.onload = function () {
                wails.Window.SetTitle("A new window title");
            }
        </script>
    </head>
    <!--- ... -->
</html>
```

@note{type="caution"}
Penting untuk menyertakan atribut `type="module"` pada tag `<script>` yang memuat runtime dan menunggu hingga halaman dimuat sepenuhnya sebelum memanggil API karena skrip dengan atribut `type="module"` berjalan secara asinkron.

@end

## Inisialisasi

Selain fungsi API, runtime menyediakan dukungan untuk menu konteks dan penyeretan jendela. Fitur-fitur ini hanya akan berfungsi sebagaimana mestinya setelah runtime diinisialisasi. Meskipun Anda tidak menggunakan API, pastikan untuk menyertakan pernyataan impor efek samping di suatu tempat dalam kode frontend Anda:

```javascript
import "@wailsio/runtime";
```

Bundler Anda seharusnya mendeteksi adanya efek samping dan menyertakan semua kode inisialisasi yang diperlukan dalam build.

@note{type="info"}
Jika Anda lebih memilih bundel siap pakai, cukup tambahkan tag skrip seperti yang ditunjukkan di atas.

@end

## Plugin Vite untuk Peristiwa Bertipe

Runtime menyertakan plugin Vite yang mengaktifkan dukungan HMR (Hot Module Replacement) untuk peristiwa bertipe selama pengembangan.

### Penyiapan

Tambahkan plugin ke `vite.config.ts` Anda:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### Manfaat

- **Pemuatan Ulang Otomatis**: Binding peristiwa dibuat ulang dan dimuat ulang secara otomatis saat Anda menjalankan `wails3 generate bindings`
- **Mode Pengembangan**: Terintegrasi lancar dengan `wails3 dev` untuk pembaruan seketika
- **Keamanan Tipe**: Dukungan penuh TypeScript dengan pelengkapan otomatis dan pemeriksaan tipe

### Penggunaan dengan Pendaftaran Peristiwa

Daftarkan peristiwa Anda di Go:

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

Buat binding:

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

Gunakan peristiwa bertipe di frontend Anda:

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## Referensi API

Runtime disusun dalam beberapa modul yang masing-masing menyediakan fungsi tertentu. Impor hanya yang Anda perlukan:

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### Peristiwa

Sistem peristiwa untuk komunikasi antara Go dan JavaScript.

#### On()

Daftarkan callback untuk suatu peristiwa.

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**Mengembalikan:** Fungsi untuk berhenti berlangganan

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event listening
const unsubscribe = Events.On('user-logged-in', (event) => {
    console.log('User:', event.data.username)
})

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log('User:', event.data.username)
})

// Later: unsubscribe()
```

#### Once()

Daftarkan callback yang hanya dijalankan satu kali.

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Pancarkan peristiwa ke backend Go atau jendela lain.

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

**Mengembalikan:** Promise yang di-resolve menjadi `true` jika peristiwa dibatalkan, atau `false` jika tidak

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
Nilai yang dikembalikan menunjukkan apakah peristiwa dibatalkan oleh hook. Sebagian besar peristiwa tidak dapat dibatalkan dan akan selalu mengembalikan `false`.

@end

#### Off()

Hapus listener peristiwa.

```typescript
function Off(...eventNames: string[]): void
```

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

Hapus semua listener peristiwa.

```typescript
function OffAll(): void
```

### Jendela

Metode pengelolaan jendela. Ekspor default adalah jendela saat ini.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### Visibilitas

**Show()** - Menampilkan jendela

```typescript
function Show(): Promise<void>
```

**Hide()** - Menyembunyikan jendela

```typescript
function Hide(): Promise<void>
```

**Close()** - Menutup jendela

```typescript
function Close(): Promise<void>
```

#### Ukuran dan Posisi

**SetSize(width, height)** - Mengatur ukuran jendela

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - Mendapatkan ukuran jendela

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - Mengatur posisi absolut

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - Mendapatkan posisi absolut

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - Menempatkan jendela di tengah

```typescript
function Center(): Promise<void>
```

**Contoh:**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### Status Jendela

**Minimise()** - Meminimalkan jendela

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - Memaksimalkan jendela

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - Beralih ke layar penuh

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - Memulihkan jendela dari keadaan diminimalkan, dimaksimalkan, atau layar penuh

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - Memeriksa apakah jendela diminimalkan

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - Memeriksa apakah jendela dimaksimalkan

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - Memeriksa apakah jendela dalam mode layar penuh

```typescript
function IsFullscreen(): Promise<boolean>
```

#### Properti Jendela

**SetTitle(title)** - Mengatur judul jendela

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - Mendapatkan nama jendela

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - Mengatur warna latar belakang

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - Menjaga jendela tetap berada di atas

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - Memungkinkan ukuran jendela diubah

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### Fokus dan Layar

**Focus()** - Memfokuskan jendela

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - Memeriksa apakah jendela memiliki fokus

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - Mendapatkan layar tempat jendela berada

```typescript
function GetScreen(): Promise<Screen>
```

#### Konten

**Reload()** - Memuat ulang halaman

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - Memaksa halaman dimuat ulang (menghapus cache)

```typescript
function ForceReload(): Promise<void>
```

#### Zoom

**SetZoom(level)** - Mengatur tingkat zoom

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - Mendapatkan tingkat zoom

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - Memperbesar tampilan

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - Memperkecil tampilan

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - Mengatur ulang zoom ke 100%

```typescript
function ZoomReset(): Promise<void>
```

#### Pencetakan

**Print()** - Membuka dialog pencetakan native

```typescript
function Print(): Promise<void>
```

**Contoh:**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**Catatan:** Tindakan ini membuka dialog pencetakan native dari sistem operasi sehingga pengguna dapat memilih pengaturan printer dan mencetak konten jendela saat ini. Tidak seperti `window.print()` yang mungkin tidak berfungsi di webview, fungsi ini menggunakan API pencetakan native platform.

### Papan Klip

Operasi papan klip.

#### SetText()

Mengatur teks papan klip.

```typescript
function SetText(text: string): Promise<void>
```

**Contoh:**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

Mendapatkan teks papan klip.

```typescript
function Text(): Promise<string>
```

**Contoh:**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### Sistem

Metode sistem tingkat rendah untuk berkomunikasi langsung dengan backend.

#### invoke()

Mengirim pesan mentah langsung ke backend. Tindakan ini melewati sistem binding standar dan ditangani oleh `RawMessageHandler` dalam opsi aplikasi Anda.

```typescript
function invoke(message: any): void
```

**Contoh:**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
Fungsi ini bersifat kirim-dan-lupakan tanpa nilai kembalian. Gunakan event untuk menerima respons dari backend.

@end

Untuk detail selengkapnya, lihat [Panduan Pesan Mentah](/guides/raw-messages/).

### Aplikasi

Metode tingkat aplikasi.

#### Show()

Menampilkan semua jendela aplikasi.

```typescript
function Show(): Promise<void>
```

#### Hide()

Menyembunyikan semua jendela aplikasi.

```typescript
function Hide(): Promise<void>
```

#### Quit()

Keluar dari aplikasi.

```typescript
function Quit(): Promise<void>
```

**Contoh:**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### Browser

Buka URL di browser default.

#### OpenURL()

Membuka URL di browser sistem.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**Contoh:**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### Layar

Informasi dan pengelolaan layar.

#### GetAll()

Mendapatkan semua layar.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

Mendapatkan layar utama.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

Mendapatkan layar yang sedang aktif.

```typescript
function GetCurrent(): Promise<Screen>
```

**Antarmuka Screen:**

```typescript
interface Screen {
    ID: string
    Name: string
    ScaleFactor: number
    X: number
    Y: number
    Size: { Width: number, Height: number }
    Bounds: { X: number, Y: number, Width: number, Height: number }
    WorkArea: { X: number, Y: number, Width: number, Height: number }
    IsPrimary: boolean
    Rotation: number
}
```

**Contoh:**

```javascript
import { Screens } from '@wailsio/runtime'

// List all screens
const screens = await Screens.GetAll()
screens.forEach(screen => {
    console.log(`${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})

// Get primary screen
const primary = await Screens.GetPrimary()
console.log('Primary screen:', primary.Name)
```

### Dialog

Dialog bawaan sistem operasi dari JavaScript.

#### Info()

Menampilkan dialog informasi.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**Contoh:**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

Menampilkan dialog kesalahan.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

Menampilkan dialog peringatan.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

Menampilkan dialog pertanyaan dengan tombol khusus.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**Contoh:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const result = await Dialogs.Question({
    Title: 'Confirm Delete',
    Message: 'Are you sure you want to delete this file?',
    Buttons: [
        { Label: 'Delete', IsDefault: false },
        { Label: 'Cancel', IsDefault: true }
    ]
})

if (result === 'Delete') {
    // Delete the file
}
```

#### OpenFile()

Menampilkan dialog untuk membuka file.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**Contoh:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const file = await Dialogs.OpenFile({
    Title: 'Select Image',
    Filters: [
        { DisplayName: 'Images', Pattern: '*.png;*.jpg;*.jpeg' },
        { DisplayName: 'All Files', Pattern: '*.*' }
    ]
})

if (file) {
    console.log('Selected:', file)
}
```

#### SaveFile()

Menampilkan dialog untuk menyimpan file.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML (Wails Markup Language)

WML menyediakan atribut deklaratif untuk tindakan umum. Tambahkan atribut ke elemen HTML:

#### Atribut

**wml-event** - Memancarkan peristiwa saat diklik

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - Memanggil metode jendela

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - Menentukan jendela target untuk wml-window

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - Membuka URL di browser

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - Menampilkan dialog konfirmasi sebelum tindakan dijalankan

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**Contoh:**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## Contoh Lengkap

```javascript
import { Events, Window, Clipboard, Dialogs, Screens } from '@wailsio/runtime'

// Listen for events from Go
Events.On('data-updated', (event) => {
    console.log('Data:', event.data)
    updateUI(event.data)
})

// Window management
document.getElementById('center-btn').addEventListener('click', async () => {
    await Window.Center()
})

document.getElementById('fullscreen-btn').addEventListener('click', async () => {
    const isFullscreen = await Window.IsFullscreen()
    if (isFullscreen) {
        await Window.UnFullscreen()
    } else {
        await Window.Fullscreen()
    }
})

// Clipboard operations
document.getElementById('copy-btn').addEventListener('click', async () => {
    await Clipboard.SetText('Copied from Wails!')
})

// Dialog with confirmation
document.getElementById('delete-btn').addEventListener('click', async () => {
    const result = await Dialogs.Question({
        Title: 'Confirm',
        Message: 'Delete this item?',
        Buttons: [
            { Label: 'Delete' },
            { Label: 'Cancel', IsDefault: true }
        ]
    })

    if (result === 'Delete') {
        await Events.Emit('delete-item', { id: currentItemId })
    }
})

// Screen information
const screens = await Screens.GetAll()
console.log(`Detected ${screens.length} screen(s)`)
screens.forEach(screen => {
    console.log(`- ${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})
```

## Praktik Terbaik

### ✅ Lakukan

- **Impor secara selektif** - Impor hanya yang Anda perlukan
- **Tangani promise** - Semua metode mengembalikan promise
- **Gunakan WML untuk tindakan sederhana** - Lebih ringkas daripada JavaScript
- **Periksa nilai kembalian** - Terutama untuk dialog
- **Berhenti berlangganan peristiwa** - Bersihkan setelah selesai

### ❌ Jangan

- **Jangan lupa menggunakan await** - Sebagian besar metode bersifat asinkron
- **Jangan blokir UI** - Gunakan async/await dengan benar
- **Jangan abaikan kesalahan** - Selalu tangani penolakan promise

## Dukungan TypeScript

Runtime menyertakan definisi TypeScript yang lengkap:

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
