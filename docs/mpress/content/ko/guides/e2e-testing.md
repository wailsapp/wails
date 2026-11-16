---
title: "엔드투엔드 테스트"
description: "전체 사용자 워크플로 테스트"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## 이 가이드의 테스트 범위

브라우저에서 프런트엔드와 실제 Go 서비스를 테스트하려면 **서버 빌드**에 Playwright를 사용하세요. 포트 9245의 Vite 미리보기는 프런트엔드 파일을 제공하지만, 일반 브라우저에서 연다고 해서 Wails 서비스 바인딩이 자동으로 제공되지는 않습니다.

이 테스트 하네스는 브라우저 UI와 서비스 호출을 다룹니다. 네이티브 WebView 창, 운영 체제 메뉴, 트레이 아이콘 또는 네이티브 파일 대화 상자는 자동화하지 않습니다. 이러한 항목은 지원되는 각 데스크톱 플랫폼에서 별도로 테스트하세요.

## 스타터 애플리케이션 준비

JavaScript 스타터를 만든 다음 해당 디렉터리에서 작업하세요:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

`main.go`에서 기존 `application.Options` 값에 다음 필드를 추가하세요:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

기존 애셋 핸들러와 서비스 등록을 유지하세요. 서버 옵션은 애플리케이션을 `server` 태그로 빌드할 때 적용됩니다. 기능과 제한 사항은 [서버 빌드](/guides/server-build/)를 참조하세요.

테스트 호스트를 시작하기 전에 바인딩을 생성하고 프로덕션 프런트엔드를 빌드하세요:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

CI에서 `npm ci`을 사용할 수 있도록 생성된 패키지 매니페스트와 잠금 파일을 커밋하세요.

## Go 호스트 시작 및 대기

`frontend/playwright.config.js`을 만드세요:

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

Playwright는 프로젝트 루트에서 Go 애플리케이션을 시작하고 HTTP 엔드포인트가 준비될 때까지 기다립니다. 테스트 작업마다 격리된 호스트 하나를 실행하세요. 영속성 테스트를 추가할 때는 임시 저장소나 테스트 데이터베이스를 사용하고, 이 테스트 하네스가 실제 사용자 데이터를 대상으로 하도록 설정해서는 안 됩니다.

## 실제 서비스 호출 확인

`frontend/e2e/app.spec.js`을 만드세요:

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

위 선택자는 생성된 `vanilla-js` 스타터에 맞춰져 있습니다. 스타터를 교체한 후에는 UI에 맞게 수정하세요. 결과 어설션은 고정된 지연 시간을 가정하지 않고 Go 호출이 완료될 때까지 기다립니다.

`frontend`에서 실행하세요:

```bash
npx playwright test
```

이 구성은 Playwright의 기본 Chromium 브라우저를 사용합니다.

## 브라우저 상호작용과 네이티브 경계

브라우저 `window.confirm`의 경우, 이를 여는 클릭 **전에** 핸들러를 등록하세요:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

HTML 업로드 입력에는 `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`을 사용하세요. 이 방식은 브라우저 입력에 파일을 제공하지만 Wails 네이티브 저장 대화 상자를 테스트하지는 않습니다. 서비스 경계를 통해 애플리케이션의 저장 결과를 테스트한 다음, 네이티브 대화 상자는 수동으로 또는 적절한 플랫폼 자동화 하네스로 확인하세요.

구성된 800×600 뷰포트는 반응형 브라우저 레이아웃을 테스트합니다. `page.viewportSize()`은 네이티브 Wails 창을 측정하지 않습니다. 마찬가지로 Playwright의 브라우저 대화 상자 핸들러는 Wails 네이티브 메시지 대화 상자를 처리하지 않습니다.

## CI

이 워크플로에서는 애플리케이션이 저장소 루트에 있고 위의 Playwright 종속성과 잠금 파일이 커밋되어 있다고 가정합니다:

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

구성의 `webServer`은 로컬뿐 아니라 CI에서도 백엔드를 시작합니다. 서버 테스트는 네이티브 데스크톱 패키징 및 실행 검사와 분리하세요.

## 다음 단계

- [테스트](/guides/testing/) — 단위 및 서비스 수준 테스트
- [서버 빌드](/guides/server-build/) — 브라우저 호스트 동작
- [빌드](/guides/build/building/) — 네이티브 빌드
