---
title: "エンドツーエンドテスト"
description: "ユーザーワークフロー全体をテストする"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## このガイドでテストする内容

ブラウザーでフロントエンドと実際の Go サービスを動作させるには、Playwright を **サーバービルド**に対して使用します。ポート 9245 の Vite プレビューはフロントエンドファイルを配信しますが、通常のブラウザーで開くだけでは Wails のサービスバインディングは利用できません。

このテストハーネスは、ブラウザー UI とサービス呼び出しを対象とします。ネイティブ WebView ウィンドウ、OS のメニュー、トレイアイコン、ネイティブファイルダイアログは自動化しません。これらは、サポート対象の各デスクトッププラットフォームで個別にテストしてください。

## スターターアプリケーションを準備する

JavaScript スターターを作成し、そのディレクトリで作業します。

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

`main.go` で、既存の `application.Options` 値に次のフィールドを追加します。

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

既存のアセットハンドラーとサービス登録は維持してください。サーバーオプションは、アプリケーションを `server` タグ付きでビルドした場合に有効になります。機能と制限については、[サーバービルド](/guides/server-build/)を参照してください。

テストホストを起動する前に、バインディングを生成し、本番用フロントエンドをビルドします。

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

CI で `npm ci` を使用できるように、生成されたパッケージマニフェストとロックファイルをコミットしてください。

## Go ホストを起動して待機する

`frontend/playwright.config.js` を作成します。

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

Playwright はプロジェクトルートから Go アプリケーションを起動し、その HTTP エンドポイントが利用可能になるまで待機します。テストジョブごとに、分離したホストを 1 つ実行してください。永続化テストを追加する場合は、一時ストレージまたはテスト用データベースを使用し、このテストハーネスを実際のユーザーデータに接続しないでください。

## 実際のサービス呼び出しを検証する

`frontend/e2e/app.spec.js` を作成します。

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

上記のセレクターは、生成された `vanilla-js` スターターに対応しています。スターターを置き換えた後は、実際の UI に合わせて調整してください。結果のアサーションは固定時間の遅延を前提とせず、Go 呼び出しの完了を待機します。

`frontend` から実行します。

```bash
npx playwright test
```

この構成では、Playwright のデフォルトの Chromium ブラウザーを使用します。

## ブラウザー操作とネイティブ機能の境界

ブラウザーの `window.confirm` では、それを開くクリックの**前に**ハンドラーを登録します。

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

HTML のアップロード入力には、`page.locator('input[type=file]').setInputFiles('fixtures/note.json')` を使用します。これはブラウザーの入力要素にファイルを渡すものであり、Wails のネイティブ保存ダイアログをテストするものではありません。アプリケーションの保存結果はサービス境界を通じてテストし、その後、ネイティブダイアログを手動または適切なプラットフォーム自動化ハーネスで検証してください。

設定された 800×600 のビューポートは、ブラウザーのレスポンシブレイアウトをテストします。`page.viewportSize()` はネイティブ Wails ウィンドウを測定しません。同様に、Playwright のブラウザーダイアログハンドラーは、Wails のネイティブメッセージダイアログを処理しません。

## CI

このワークフローでは、アプリケーションがリポジトリルートにあり、上記の Playwright 依存関係とロックファイルがコミットされていることを前提とします。

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

構成内の `webServer` は、ローカル環境だけでなく CI でもバックエンドを起動します。サーバーテストは、ネイティブデスクトップのパッケージ化および起動確認とは分けてください。

## 次のステップ

- [テスト](/guides/testing/) — ユニットテストとサービスレベルテスト
- [サーバービルド](/guides/server-build/) — ブラウザーホストの動作
- [ビルド](/guides/build/building/) — ネイティブビルド
