---
title: "Сквозное тестирование"
description: "Тестирование полных пользовательских сценариев"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Обзор

Сквозное тестирование проверяет полные пользовательские сценарии в вашем приложении.

## Использование Playwright

### Настройка

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### Конфигурация

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

### Написание тестов

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

## Тестирование диалоговых окон

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

## Тестирование поведения окна

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

## Рекомендации

### ✅ Рекомендуется

- Тестируйте критически важные пользовательские сценарии
- Используйте атрибуты data-testid
- Очищайте тестовые данные
- Запускайте тесты в CI/CD
- Тестируйте сценарии обработки ошибок
- Обеспечивайте независимость тестов

### ❌ Не рекомендуется

- Не тестируйте детали реализации
- Не используйте ненадёжные селекторы
- Не пропускайте очистку
- Не игнорируйте нестабильные тесты
- Не пытайтесь протестировать всё

## Запуск E2E-тестов

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

## Интеграция с CI/CD

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

## Дальнейшие шаги

- [Тестирование](/guides/testing/) — Узнайте о модульном тестировании
- [Сборка](/guides/build/building/) — Соберите приложение
