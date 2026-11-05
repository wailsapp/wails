---
title: "端對端測試"
description: "測試完整的使用者工作流程"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 本指南的測試範圍

使用 Playwright 對**伺服器組建**執行測試，以便在瀏覽器中測試前端和實際的 Go 服務。連接埠9245上的 Vite 預覽會提供前端檔案；僅在一般瀏覽器中開啟該預覽，並不會提供 Wails 服務繫結。

此測試架構涵蓋瀏覽器 UI 和服務呼叫，但不會自動化原生 WebView 視窗、作業系統選單、系統匣圖示或原生檔案對話方塊。請在每個支援的桌面平台上分別測試這些項目。

## 準備入門應用程式

建立 JavaScript 入門專案，然後在其目錄中操作：

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

在`main.go`中，將此欄位新增至現有的`application.Options`值：

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

保留現有的資產處理常式和服務註冊。使用`server`標記建置應用程式時，伺服器選項才會生效。如需瞭解其功能與限制，請參閱[伺服器組建](/guides/server-build/)。

啟動測試主機之前，請先產生繫結並建置正式環境前端：

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

提交產生的套件資訊清單和鎖定檔，讓 CI 可以使用`npm ci`。

## 啟動 Go 主機並等待其就緒

建立`frontend/playwright.config.js`：

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

Playwright 會從專案根目錄啟動 Go 應用程式，並等待其 HTTP 端點就緒。每個測試作業請執行一個隔離的主機。新增持久性測試時，請使用暫存儲存空間或測試資料庫；絕不可讓此測試架構存取真實的使用者資料。

## 驗證實際的服務呼叫

建立`frontend/e2e/app.spec.js`：

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

上述選取器與產生的`vanilla-js`入門專案相符。替換入門專案後，請依據你的 UI 調整這些選取器。結果斷言會等待 Go 呼叫完成，而不是假設固定的延遲時間。

從`frontend`執行：

```bash
npx playwright test
```

此設定使用 Playwright 的預設 Chromium 瀏覽器。

## 瀏覽器互動與原生功能邊界

針對瀏覽器`window.confirm`，請在開啟它的點擊動作**之前**註冊處理常式：

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

針對 HTML 上傳輸入欄位，請使用`page.locator('input[type=file]').setInputFiles('fixtures/note.json')`。這會將檔案提供給瀏覽器輸入欄位；它不會測試 Wails 原生儲存對話方塊。請透過應用程式的服務邊界測試儲存結果，再手動或使用適當的平台自動化測試架構驗證原生對話方塊。

設定的800×600檢視區會測試瀏覽器的響應式版面配置。`page.viewportSize()`不會測量原生 Wails 視窗。同樣地，Playwright 的瀏覽器對話方塊處理常式也不會處理 Wails 原生訊息對話方塊。

## CI

此工作流程假設應用程式位於儲存庫根目錄，且已提交上述 Playwright 相依套件和鎖定檔：

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

此設定中的`webServer`會在 CI 和本機環境中啟動後端。請將伺服器測試與原生桌面應用程式的封裝及啟動檢查分開。

## 後續步驟

- [測試](/guides/testing/) — 單元測試與服務層級測試
- [伺服器組建](/guides/server-build/) — 瀏覽器主機行為
- [建置](/guides/build/building/) — 原生組建
