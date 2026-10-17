---
title: "端對端測試"
description: "測試完整的使用者工作流程"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 概觀

端對端測試可驗證應用程式中的完整使用者工作流程。

## 使用 Playwright

### 設定

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### 組態設定

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

### 撰寫測試

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

## 測試對話方塊

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

## 測試視窗行為

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

## 最佳實務

### ✅ 建議做法

- 測試關鍵使用者流程
- 使用 data-testid 屬性
- 清理測試資料
- 在 CI/CD 中執行測試
- 測試錯誤情境
- 讓各項測試保持獨立

### ❌ 請勿這樣做

- 不要測試實作細節
- 不要撰寫脆弱的選取器
- 不要略過清理步驟
- 不要忽略不穩定的測試
- 不要測試所有內容

## 執行端對端測試

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

## CI/CD 整合

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

## 後續步驟

- [測試](/guides/testing/) - 瞭解單元測試
- [建置](/guides/build/building/) - 建置您的應用程式
