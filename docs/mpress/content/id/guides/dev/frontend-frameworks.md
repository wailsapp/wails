---
title: "Menggunakan Framework Frontend Lain"
description: "Cara menggunakan framework yang tidak memiliki templat bawaan dengan menempatkan proyek Vite Anda sendiri ke dalam direktori frontend"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails menyediakan templat awal bawaan untuk sejumlah kecil framework yang dipilih secara sengaja:

| Templat | Bahasa |
| --- | --- |
| `vanilla` | TypeScript (bawaan) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

Namun, pilihan Anda tidak terbatas pada framework tersebut. Frontend aplikasi Wails **hanyalah proyek web** — apa pun yang dapat dibangun menjadi HTML/CSS/JS statis bisa digunakan. Jika framework pilihan Anda (Solid, Preact, Lit, Qwik, SvelteKit, Angular, …) tidak memiliki templat, Anda dapat melakukan scaffolding sendiri dalam beberapa menit.

## Cara direktori `frontend/` digunakan

Wails tidak mempermasalahkan framework apa yang berada di `frontend/`. Wails hanya mengandalkan kontrak kecil yang tidak bergantung pada framework:

- **`frontend/dist/` adalah artefak yang didistribusikan.** `main.go` menyematkan frontend yang telah dibangun dengan `//go:embed all:frontend/dist` dan menyajikannya dari server aset. Proses build Anda harus menghasilkan bundel statis di `frontend/dist/` (direktori keluaran bawaan Vite).
- **Build dikendalikan oleh `frontend/package.json`.** Selama `wails3 build`, Wails menjalankan skrip `build` milik frontend; selama `wails3 dev`, Wails menjalankan `dev` dan memproksikan server pengembangan Vite untuk hot reload.
- **Binding dibuat di `frontend/bindings/`.** Wails memeriksa layanan Go yang Anda daftarkan dan menulis SDK yang aman secara tipe di sana. Anda dapat mengimpornya seperti modul lain:
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **Server pengembangan berjalan pada port tetap.** `wails3 dev` memproksikan Vite pada port yang tercantum di `WAILS_VITE_PORT` (bawaan: `9245`), jadi tetapkan `server.port` ke port tersebut dengan `strictPort: true`. Ini adalah konfigurasi Vite biasa — tidak melibatkan plugin Wails.
- **Opsional — peristiwa khusus bertipe.** Templat bawaan juga mendaftarkan plugin `@wailsio/runtime/plugins/vite`. Plugin ini hanya diperlukan jika Anda menggunakan peristiwa khusus *bertipe*: plugin tersebut menyuntikkan definisi tipe peristiwa yang dihasilkan ke dalam runtime dan membuat build gagal hingga binding dibuat. Jika Anda hanya menggunakan API `Events.On("time", …)` berbasis string, plugin ini dapat dihilangkan.

Selebihnya — komponen, perutean, state, dan gaya — sepenuhnya bergantung pada framework Anda.

## Lakukan scaffolding framework apa pun dengan Vite

Cara tercepat adalah memulai dari templat bawaan (agar Anda mendapatkan `main.go`, `Taskfile`, aset build, dan layanan Go yang berfungsi), lalu mengganti `frontend/` dengan proyek Vite baru untuk framework Anda.

@steps
### Buat proyek dari templat bawaan
```bash
wails3 init -n myapp
cd myapp
```

### Ganti `frontend/` dengan aplikasi Vite untuk framework Anda
Vite dapat melakukan scaffolding untuk sebagian besar framework dengan satu perintah. Pilih templat:

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

Ganti `solid` dengan templat Vite apa pun: `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla`, atau varian `-ts` masing-masing (`solid-ts`, `preact-ts`, …).

### Instal runtime dan arahkan Vite ke server pengembangan Wails
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` menyediakan API JS (`Events`, `Browser`, dialog, …). Satu-satunya perubahan `vite.config` yang *wajib* adalah port server pengembangan, agar `wails3 dev` dapat menemukannya:

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

Hanya jika Anda bermaksud menggunakan **peristiwa khusus bertipe**, tambahkan juga plugin tersebut — plugin ini menyuntikkan tipe peristiwa yang dihasilkan dan mengharuskan binding sudah tersedia sebelum build:

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Panggil layanan Go Anda
Buat binding satu kali, lalu impor binding tersebut di mana pun dalam komponen Anda:

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### Jalankan
```bash
wails3 dev
```

@end

@note{type="tip" title="Buat proyek JavaScript biasa"}
Perintah yang sama melakukan scaffolding proyek tanpa TypeScript — cukup gunakan templat Vite non-`-ts`:

```bash
npm create vite@latest frontend -- --template solid
```

@end

## Framework dengan alat scaffolding sendiri

Beberapa framework tidak dibuat melalui templat `create` Vite dan memiliki alatnya sendiri. Framework tersebut tetap dapat digunakan — cukup lakukan scaffolding dengan perintah bawaannya, lalu tambahkan plugin Wails setelahnya:

- **SvelteKit:** `npx sv create frontend`. Gunakan adaptor statis (`@sveltejs/adapter-static`) agar menghasilkan bundel statis, dan nonaktifkan SSR.
- **Qwik:** `npm create qwik@latest`. Gunakan adaptor statis (SSG).
- **Angular:** lakukan scaffolding dengan `ng new`, tetapkan `outputPath` ke `dist`, lalu arahkan skrip build ke `ng build`.

Aturannya selalu sama: hasilkan build statis di `frontend/dist/`, pertahankan plugin Vite `@wailsio/runtime` (atau impor runtime secara langsung), dan impor binding Go Anda dari `frontend/bindings/`.

@note{type="info"}
Jika Anda membuat konfigurasi framework yang siap digunakan, pertimbangkan untuk menerbitkannya sebagai [templat khusus](/guides/advanced/custom-templates/) agar orang lain dapat melakukan `wails3 init -t` secara langsung.

@end
