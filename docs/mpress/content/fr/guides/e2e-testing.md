---
title: "Tests de bout en bout"
description: "Testez des parcours utilisateur complets"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Vue d’ensemble

Les tests de bout en bout valident des parcours utilisateur complets dans votre application.

## Utilisation de Playwright

### Installation

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### Configuration

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

### Écriture des tests

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

## Test des boîtes de dialogue

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

## Test du comportement des fenêtres

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

## Bonnes pratiques

### ✅ À faire

- Testez les parcours utilisateur critiques
- Utilisez des attributs data-testid
- Nettoyez les données de test
- Exécutez les tests dans le pipeline CI/CD
- Testez les scénarios d’erreur
- Veillez à ce que les tests restent indépendants

### ❌ À ne pas faire

- Ne testez pas les détails d’implémentation
- N’écrivez pas de sélecteurs fragiles
- Ne négligez pas le nettoyage
- N’ignorez pas les tests instables
- Ne testez pas tout

## Exécution des tests de bout en bout

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

## Intégration CI/CD

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

## Étapes suivantes

- [Tests](/guides/testing/) – Découvrez les tests unitaires
- [Compilation](/guides/build/building/) – Compilez votre application
