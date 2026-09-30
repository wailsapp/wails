---
title: "End-to-End-Tests"
description: "Vollständige Benutzerabläufe testen"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Überblick

End-to-End-Tests validieren vollständige Benutzerabläufe in Ihrer Anwendung.

## Playwright verwenden

### Einrichtung

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### Konfiguration

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

### Tests schreiben

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

## Dialoge testen

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

## Fensterverhalten testen

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- Kritische Benutzerabläufe testen
- data-testid-Attribute verwenden
- Testdaten bereinigen
- Tests in CI/CD ausführen
- Fehlerszenarien testen
- Tests unabhängig voneinander halten

### ❌ Nicht empfohlen

- Keine Implementierungsdetails testen
- Keine fragilen Selektoren schreiben
- Die Bereinigung nicht überspringen
- Flaky Tests nicht ignorieren
- Nicht alles testen

## E2E-Tests ausführen

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

## CI/CD-Integration

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

## Nächste Schritte

- [Tests](/guides/testing/) – Unit-Tests kennenlernen
- [Build](/guides/build/building/) – Anwendung erstellen
