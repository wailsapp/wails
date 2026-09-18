---
title: "엔드 투 엔드 테스트"
description: "전체 사용자 워크플로 테스트"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 개요

엔드 투 엔드 테스트는 애플리케이션의 전체 사용자 워크플로가 올바르게 작동하는지 검증합니다.

## Playwright 사용하기

### 설정

```bash
# Install Playwright
npm install -D @playwright/test

# Initialize
npx playwright install
```

### 구성

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

### 테스트 작성하기

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

## 대화 상자 테스트하기

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

## 창 동작 테스트하기

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

## 권장 사례

### ✅ 권장 사항

- 중요한 사용자 흐름을 테스트하세요
- data-testid 속성을 사용하세요
- 테스트 데이터를 정리하세요
- CI/CD에서 테스트를 실행하세요
- 오류 시나리오를 테스트하세요
- 각 테스트를 독립적으로 유지하세요

### ❌ 피해야 할 사항

- 구현 세부 사항을 테스트하지 마세요
- 쉽게 깨지는 선택자를 작성하지 마세요
- 정리 작업을 생략하지 마세요
- 불안정한 테스트를 무시하지 마세요
- 모든 것을 테스트하려고 하지 마세요

## E2E 테스트 실행하기

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

## CI/CD 통합

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

## 다음 단계

- [테스트](/guides/testing/) - 단위 테스트 알아보기
- [빌드](/guides/build/building/) - 애플리케이션 빌드하기
