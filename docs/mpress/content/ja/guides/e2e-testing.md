---
title: "エンドツーエンドテスト"
description: "ユーザーワークフロー全体をテストする"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 概要

エンドツーエンドテストでは、アプリケーション内のユーザーワークフロー全体を検証します。

## Playwright の使用

### セットアップ

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### 設定

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

### テストの作成

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

## ダイアログのテスト

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

## ウィンドウ動作のテスト

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

## ベストプラクティス

### ✅ 推奨事項

- 重要なユーザーフローをテストする
- data-testid 属性を使用する
- テストデータをクリーンアップする
- CI/CD でテストを実行する
- エラーシナリオをテストする
- 各テストを独立させる

### ❌ 禁止事項

- 実装の詳細をテストしない
- 壊れやすいセレクターを作成しない
- クリーンアップを省略しない
- 不安定なテストを放置しない
- すべてをテストしようとしない

## E2E テストの実行

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

## CI/CD との統合

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

## 次のステップ

- [テスト](/guides/testing/) - 単体テストについて学ぶ
- [ビルド](/guides/build/building/) - アプリケーションをビルドする
