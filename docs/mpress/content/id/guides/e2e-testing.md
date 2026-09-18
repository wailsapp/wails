---
title: "Pengujian Menyeluruh"
description: "Uji alur kerja pengguna secara lengkap"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Ikhtisar

Pengujian menyeluruh memvalidasi alur kerja pengguna secara lengkap dalam aplikasi Anda.

## Menggunakan Playwright

### Penyiapan

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### Konfigurasi

```javascript
// playwright.config.js
import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  use: {
    baseURL: 'http://localhost:9245', // Wails v3 dev server (default port)
  },
})
```

### Menulis Pengujian

```javascript
// e2e/app.spec.js
import { test, expect } from '@playwright/test'

test('create note', async ({ page }) => {
  await page.goto('/')
  
  // Click new note button
  await page.click('#new-note-btn')
  
  // Fill in title
  await page.fill('#note-title', 'Test Note')
  
  // Fill in content
  await page.fill('#note-content', 'Test content')
  
  // Verify note appears in list
  await expect(page.locator('.note-item')).toContainText('Test Note')
})

test('delete note', async ({ page }) => {
  await page.goto('/')
  
  // Create a note first
  await page.click('#new-note-btn')
  await page.fill('#note-title', 'To Delete')
  
  // Delete it
  await page.click('#delete-btn')
  
  // Confirm dialog
  page.on('dialog', dialog => dialog.accept())
  
  // Verify it's gone
  await expect(page.locator('.note-item')).not.toContainText('To Delete')
})
```

## Menguji dialog

```javascript
test('file save dialog', async ({ page }) => {
  await page.goto('/')
  
  // Intercept file dialog
  page.on('filechooser', async (fileChooser) => {
    await fileChooser.setFiles('/path/to/test/file.json')
  })
  
  // Trigger save
  await page.click('#save-btn')
  
  // Verify success message
  await expect(page.locator('.success-message')).toBeVisible()
})
```

## Menguji Perilaku Jendela

```javascript
test('window state', async ({ page }) => {
  await page.goto('/')
  
  // Test window title
  await expect(page).toHaveTitle('My App')
  
  // Test window size
  const size = await page.viewportSize()
  expect(size.width).toBe(800)
  expect(size.height).toBe(600)
})
```

## Praktik Terbaik

### ✅ Lakukan

- Uji alur pengguna yang penting
- Gunakan atribut data-testid
- Bersihkan data pengujian
- Jalankan pengujian dalam CI/CD
- Uji skenario kesalahan
- Pastikan setiap pengujian tetap independen

### ❌ Jangan Lakukan

- Jangan menguji detail implementasi
- Jangan menulis pemilih yang rapuh
- Jangan melewatkan pembersihan
- Jangan mengabaikan pengujian yang tidak stabil
- Jangan menguji semuanya

## Menjalankan Pengujian E2E

```bash
# Run all tests
npx playwright test

# Run in headed mode
npx playwright test --headed

# Run specific test
npx playwright test e2e/app.spec.js

# Debug mode
npx playwright test --debug
```

## Integrasi CI/CD

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on: [push]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - name: Install dependencies
        run: npm ci
      - name: Install Playwright
        run: npx playwright install --with-deps
      - name: Run tests
        run: npx playwright test
```

## Langkah Berikutnya

- [Pengujian](/guides/testing/) - Pelajari pengujian unit
- [Membangun](/guides/build/building/) - Bangun aplikasi Anda
