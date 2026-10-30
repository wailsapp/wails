---
title: "Pengujian End-to-End"
description: "Uji alur kerja pengguna secara menyeluruh"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Yang diuji dalam panduan ini

Gunakan Playwright pada **build server** untuk menguji frontend dan layanan Go yang sebenarnya di browser. Pratinjau Vite pada port 9245 menyajikan file frontend; membukanya di browser biasa tidak serta-merta menyediakan binding layanan Wails.

Harness ini mencakup UI browser dan panggilan layanan. Harness ini tidak mengotomatiskan jendela WebView native, menu sistem operasi, ikon baki sistem, atau dialog file native. Uji semua itu secara terpisah pada setiap platform desktop yang didukung.

## Siapkan aplikasi awal

Buat aplikasi awal JavaScript, lalu kerjakan di direktorinya:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

Di `main.go`, tambahkan bidang ini ke nilai `application.Options` yang sudah ada:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Pertahankan handler aset dan pendaftaran layanan yang sudah ada. Opsi server berlaku saat aplikasi dibangun dengan tag `server`. Lihat [Build Server](/guides/server-build/) untuk mengetahui kemampuan dan batasannya.

Buat binding dan bangun frontend produksi sebelum memulai host pengujian:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Commit manifes paket dan lockfile yang dihasilkan agar CI dapat menggunakan `npm ci`.

## Mulai dan tunggu host Go

Buat `frontend/playwright.config.js`:

```javascript
import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  use: {
    baseURL: 'http://127.0.0.1:8080',
    viewport: { width: 800, height: 600 },
  },
  webServer: {
    command: 'go run -tags server,production .',
    cwd: '..',
    url: 'http://127.0.0.1:8080',
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
})
```

Playwright memulai aplikasi Go dari root proyek dan menunggu endpoint HTTP-nya. Jalankan satu host terisolasi untuk setiap job pengujian. Gunakan penyimpanan sementara atau basis data pengujian saat menambahkan pengujian persistensi; jangan pernah arahkan harness ini ke data pengguna yang sebenarnya.

## Verifikasi panggilan layanan yang sebenarnya

Buat `frontend/e2e/app.spec.js`:

```javascript
import { test, expect } from '@playwright/test'

test('greet through the Go service', async ({ page }) => {
  await page.goto('/')
  await page.locator('#name').fill('Ada')
  // The starter's aria-label supplies the button's accessible name.
  await page.getByRole('button', { name: 'greet-btn' }).click()
  await expect(page.locator('#result')).toContainText('Ada')
})
```

Selektor di atas cocok dengan aplikasi awal `vanilla-js` yang dihasilkan. Sesuaikan selektor tersebut dengan UI Anda setelah mengganti aplikasi awal. Pernyataan hasil menunggu panggilan Go, bukan mengasumsikan jeda tetap.

Jalankan dari `frontend`:

```bash
npx playwright test
```

Konfigurasi ini menggunakan browser Chromium default Playwright.

## Interaksi browser dan batas native

Untuk `window.confirm` browser, daftarkan handler **sebelum** klik yang membukanya:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

Untuk input unggahan HTML, gunakan `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. Tindakan ini memberikan file ke input browser; tindakan ini tidak menguji dialog simpan native Wails. Uji hasil penyimpanan aplikasi melalui batas layanannya, lalu verifikasi dialog native secara manual atau dengan harness otomatisasi platform yang sesuai.

Viewport 800×600 yang dikonfigurasi menguji tata letak browser responsif. `page.viewportSize()` tidak mengukur jendela Wails native. Demikian pula, handler dialog browser Playwright tidak menangani dialog pesan native Wails.

## CI

Alur kerja ini mengasumsikan aplikasi berada di root repositori serta dependensi Playwright dan lockfile di atas telah di-commit:

```yaml
name: Browser integration tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: npm
          cache-dependency-path: frontend/package-lock.json
      - run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
      - run: wails3 generate bindings -f '-tags=server'
      - run: npm ci
        working-directory: frontend
      - run: npm run build
        working-directory: frontend
      - run: npx playwright install --with-deps chromium
        working-directory: frontend
      - run: npx playwright test
        working-directory: frontend
```

`webServer` pada konfigurasi memulai backend baik di CI maupun secara lokal. Pisahkan pengujian server dari pemeriksaan pengemasan dan peluncuran desktop native.

## Langkah berikutnya

- [Pengujian](/guides/testing/) — pengujian tingkat unit dan layanan
- [Build Server](/guides/server-build/) — perilaku host browser
- [Pembangunan](/guides/build/building/) — build native
