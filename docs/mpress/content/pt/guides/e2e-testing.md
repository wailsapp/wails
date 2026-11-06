---
title: "Testes de ponta a ponta"
description: "Teste fluxos de usuário completos"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## O que este guia testa

Use o Playwright com uma **compilação para servidor** para testar o frontend e serviços Go reais em um navegador. A prévia do Vite na porta 9245 disponibiliza os arquivos do frontend; abri-la em um navegador comum não fornece, por si só, os bindings de serviço do Wails.

Esse ambiente de testes abrange a interface no navegador e as chamadas de serviço. Ele não automatiza janelas WebView nativas, menus do sistema operacional, ícones da área de notificação nem caixas de diálogo de arquivos nativas. Teste esses elementos separadamente em cada plataforma de desktop compatível.

## Prepare um aplicativo inicial

Crie um projeto inicial em JavaScript e trabalhe no diretório dele:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

Em `main.go`, adicione este campo ao valor existente de `application.Options`:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Mantenha o manipulador de ativos e o registro de serviços existentes. As opções do servidor entram em vigor quando o aplicativo é compilado com a tag `server`. Consulte [Compilação para servidor](/guides/server-build/) para conhecer seus recursos e limites.

Gere os bindings e compile o frontend de produção antes de iniciar o host de teste:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Faça commit do manifesto do pacote e do arquivo de lock resultantes para que a CI possa usar `npm ci`.

## Inicie e aguarde o host Go

Crie `frontend/playwright.config.js`:

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

O Playwright inicia o aplicativo Go a partir da raiz do projeto e aguarda seu endpoint HTTP. Execute um host isolado por tarefa de teste. Use armazenamento temporário ou um banco de dados de teste ao adicionar testes de persistência; nunca direcione esse ambiente de testes a dados reais de usuários.

## Verifique uma chamada de serviço real

Crie `frontend/e2e/app.spec.js`:

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

Os seletores acima correspondem ao projeto inicial `vanilla-js` gerado. Adapte-os à sua interface depois de substituir o projeto inicial. A asserção do resultado aguarda a chamada Go, em vez de presumir um atraso fixo.

Execute a partir de `frontend`:

```bash
npx playwright test
```

Essa configuração usa o navegador Chromium padrão do Playwright.

## Interações no navegador e limites nativos

Para uma `window.confirm` do navegador, registre o manipulador **antes** do clique que a abre:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

Para um campo HTML de upload, use `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. Isso fornece um arquivo ao campo do navegador; não testa uma caixa de diálogo nativa de salvamento do Wails. Teste o resultado da operação de salvamento do aplicativo por meio do limite do serviço e depois verifique a caixa de diálogo nativa manualmente ou com um ambiente de automação adequado à plataforma.

A viewport configurada de 800 × 600 testa o layout responsivo no navegador. `page.viewportSize()` não mede uma janela nativa do Wails. Da mesma forma, os manipuladores de caixas de diálogo do navegador do Playwright não tratam caixas de diálogo de mensagens nativas do Wails.

## CI

Este fluxo de trabalho pressupõe que o aplicativo esteja na raiz do repositório e que a dependência do Playwright e o arquivo de lock mencionados acima tenham sido incluídos em um commit:

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

O `webServer` da configuração inicia o backend tanto na CI quanto localmente. Mantenha os testes de servidor separados do empacotamento e das verificações de inicialização do aplicativo de desktop nativo.

## Próximas etapas

- [Testes](/guides/testing/) — testes de unidade e no nível de serviço
- [Compilação para servidor](/guides/server-build/) — comportamento do host no navegador
- [Compilação](/guides/build/building/) — compilações nativas
