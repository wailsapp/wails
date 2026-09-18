---
title: "Testes de ponta a ponta"
description: "Teste fluxos de trabalho completos do usuário"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Visão geral

Os testes de ponta a ponta validam fluxos de trabalho completos do usuário no seu aplicativo.

## Como usar o Playwright

### Configuração inicial

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### Configuração

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

### Como escrever testes

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

## Como testar caixas de diálogo

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

## Como testar o comportamento da janela

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

## Boas práticas

### ✅ Faça

- Teste os fluxos críticos do usuário
- Use atributos data-testid
- Limpe os dados de teste
- Execute os testes no CI/CD
- Teste cenários de erro
- Mantenha os testes independentes

### ❌ Não faça

- Não teste detalhes de implementação
- Não escreva seletores frágeis
- Não pule a limpeza
- Não ignore testes instáveis
- Não teste tudo

## Como executar testes E2E

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

## Integração com CI/CD

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

## Próximas etapas

- [Testes](/guides/testing/) — Aprenda a realizar testes unitários
- [Compilação](/guides/build/building/) — Compile seu aplicativo
