---
title: "端到端测试"
description: "测试完整的用户工作流"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 本指南的测试范围

使用 Playwright 针对**服务器构建**进行测试，以便在浏览器中测试前端和真实的 Go 服务。端口9245上的 Vite 预览服务器提供前端文件；仅在普通浏览器中打开该地址并不会提供 Wails 服务绑定。

此测试框架涵盖浏览器 UI 和服务调用，但不会自动测试原生 WebView 窗口、操作系统菜单、托盘图标或原生文件对话框。请在每个受支持的桌面平台上分别测试这些功能。

## 准备入门应用

创建一个 JavaScript 入门应用，然后进入其目录：

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

在`main.go`中，将此字段添加到现有的`application.Options`值中：

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

保留现有的资源处理程序和服务注册。使用`server`标签构建应用时，服务器选项才会生效。有关其功能和限制，请参阅[服务器构建](/guides/server-build/)。

启动测试宿主之前，请先生成绑定并构建生产版前端：

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

提交生成的软件包清单和锁文件，以便 CI 使用`npm ci`。

## 启动 Go 宿主并等待其就绪

创建`frontend/playwright.config.js`：

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

Playwright 从项目根目录启动 Go 应用，并等待其 HTTP 端点就绪。每个测试作业应运行一个相互隔离的宿主。添加持久化测试时，请使用临时存储或测试数据库；切勿让此测试框架连接真实用户数据。

## 验证真实的服务调用

创建`frontend/e2e/app.spec.js`：

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

上述选择器与生成的`vanilla-js`入门应用匹配。替换入门应用后，请根据你的 UI 调整这些选择器。结果断言会等待 Go 调用完成，而不是假定一个固定延迟。

从`frontend`运行：

```bash
npx playwright test
```

此配置使用 Playwright 默认的 Chromium 浏览器。

## 浏览器交互与原生功能边界

对于浏览器`window.confirm`，请在触发打开操作的点击**之前**注册处理程序：

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

对于 HTML 上传输入框，请使用`page.locator('input[type=file]').setInputFiles('fixtures/note.json')`。这会向浏览器输入框提供一个文件；它不会测试 Wails 原生保存对话框。请通过应用的服务边界测试保存结果，然后手动验证原生对话框，或使用适当的平台自动化测试框架进行验证。

配置的800×600视口用于测试响应式浏览器布局。`page.viewportSize()`不会测量原生 Wails 窗口。同样，Playwright 的浏览器对话框处理程序无法处理 Wails 原生消息对话框。

## CI

此工作流假定应用位于仓库根目录，并且已提交上述 Playwright 依赖项和锁文件：

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

配置中的`webServer`会在 CI 和本地环境中启动后端。请将服务器测试与原生桌面应用的打包和启动检查分开。

## 后续步骤

- [测试](/guides/testing/) — 单元测试和服务级测试
- [服务器构建](/guides/server-build/) — 浏览器宿主行为
- [构建](/guides/build/building/) — 原生构建
