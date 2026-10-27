---
title: "Сквозное тестирование"
description: "Тестирование полных пользовательских сценариев"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Что проверяется в этом руководстве

Используйте Playwright со **серверной сборкой**, чтобы тестировать интерфейс и реальные сервисы Go в браузере. Предварительный просмотр Vite на порте 9245 раздаёт файлы интерфейса, однако простое открытие этой страницы в обычном браузере само по себе не предоставляет привязки к сервисам Wails.

Эта тестовая среда охватывает браузерный интерфейс и вызовы сервисов. Она не автоматизирует нативные окна WebView, меню операционной системы, значки в области уведомлений и нативные диалоговые окна выбора файлов. Тестируйте их отдельно на каждой поддерживаемой настольной платформе.

## Подготовка стартового приложения

Создайте стартовое приложение на JavaScript, затем перейдите в его каталог:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

В `main.go` добавьте это поле в существующее значение `application.Options`:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Сохраните существующие обработчик ресурсов и регистрацию сервисов. Параметры сервера вступают в силу, когда приложение собирается с тегом `server`. Возможности и ограничения описаны в разделе [«Серверная сборка»](/guides/server-build/).

Перед запуском тестового хоста сгенерируйте привязки и соберите интерфейс для рабочей среды:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Зафиксируйте полученные манифест пакета и файл блокировки, чтобы CI мог использовать `npm ci`.

## Запуск и ожидание готовности хоста Go

Создайте `frontend/playwright.config.js`:

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

Playwright запускает приложение Go из корня проекта и ожидает готовности его конечной точки HTTP. Запускайте по одному изолированному хосту для каждой тестовой задачи. При добавлении тестов постоянного хранения используйте временное хранилище или тестовую базу данных; никогда не направляйте эту тестовую среду на реальные пользовательские данные.

## Проверка реального вызова сервиса

Создайте `frontend/e2e/app.spec.js`:

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

Приведённые выше селекторы соответствуют сгенерированному стартовому приложению `vanilla-js`. После замены стартового приложения адаптируйте их к своему интерфейсу. Проверка результата ожидает завершения вызова Go, а не предполагает фиксированную задержку.

Запустите из `frontend`:

```bash
npx playwright test
```

В этой конфигурации используется браузер Chromium, выбранный Playwright по умолчанию.

## Взаимодействие в браузере и границы нативной среды

Для браузерного `window.confirm` зарегистрируйте обработчик **до** щелчка, который его открывает:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

Для HTML-поля загрузки используйте `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. Так файл передаётся полю ввода в браузере; нативное диалоговое окно сохранения Wails при этом не тестируется. Проверьте результат сохранения в приложении через границу его сервиса, а затем проверьте нативное диалоговое окно вручную или с помощью подходящей среды автоматизации для соответствующей платформы.

Настроенная область просмотра 800×600 проверяет адаптивную компоновку в браузере. `page.viewportSize()` не измеряет нативное окно Wails. Аналогично обработчики браузерных диалоговых окон Playwright не обрабатывают нативные диалоговые окна сообщений Wails.

## CI

Этот рабочий процесс предполагает, что приложение находится в корне репозитория, а указанная выше зависимость Playwright и файл блокировки зафиксированы в репозитории:

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

Параметр `webServer` в конфигурации запускает серверную часть как в CI, так и локально. Выполняйте серверные тесты отдельно от упаковки нативного приложения для настольных систем и проверок его запуска.

## Дальнейшие действия

- [Тестирование](/guides/testing/) — модульные тесты и тесты на уровне сервисов
- [Серверная сборка](/guides/server-build/) — поведение хоста в браузере
- [Сборка](/guides/build/building/) — нативные сборки
