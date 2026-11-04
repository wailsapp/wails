---
title: "End-to-End Testing"
description: "Test complete user workflows"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## What this guide tests

Use Playwright against a **server build** to exercise your frontend and real Go services in a browser. The Vite preview at port 9245 serves frontend files; opening it in an ordinary browser does not by itself provide Wails service bindings.

This harness covers browser UI and service calls. It does not automate native WebView windows, operating-system menus, tray icons or native file dialogs. Test those separately on each supported desktop platform.

## Prepare a starter application

Create a JavaScript starter, then work in its directory:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

In `main.go`, add this field to the existing `application.Options` value:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Keep the existing asset handler and service registration. The server options take effect when the application is built with the `server` tag. See [Server Build](/guides/server-build/) for its capabilities and limits.

Generate bindings and build the production frontend before starting the test host:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Commit the resulting package manifest and lockfile so CI can use `npm ci`.

## Start and wait for the Go host

Create `frontend/playwright.config.js`:

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

Playwright starts the Go application from the project root and waits for its HTTP endpoint. Run one isolated host per test job. Use temporary storage or a test database when adding persistence tests; never aim this harness at real user data.

## Verify a real service call

Create `frontend/e2e/app.spec.js`:

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

The selectors above match the generated `vanilla-js` starter. Adapt them to your UI after replacing the starter. The result assertion waits for the Go call rather than assuming a fixed delay.

Run from `frontend`:

```bash
npx playwright test
```

This configuration uses Playwright’s default Chromium browser.

## Browser interactions and native boundaries

For a browser `window.confirm`, register the handler **before** the click that opens it:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

For an HTML upload input, use `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. This supplies a file to the browser input; it does not test a Wails native save dialog. Test your application’s save outcome through its service boundary, then verify the native dialog manually or with an appropriate platform automation harness.

The configured 800×600 viewport tests responsive browser layout. `page.viewportSize()` does not measure a native Wails window. Similarly, Playwright browser-dialog handlers do not handle Wails native message dialogs.

## CI

This workflow assumes the application is at the repository root and that the Playwright dependency and lockfile above are committed:

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

The configuration's `webServer` starts the backend in CI as well as locally. Keep server tests separate from native desktop packaging and launch checks.

## Next steps

- [Testing](/guides/testing/) — unit and service-level tests
- [Server Build](/guides/server-build/) — browser-host behaviour
- [Building](/guides/build/building/) — native builds
