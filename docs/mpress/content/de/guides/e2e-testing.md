---
title: "End-to-End-Tests"
description: "Vollständige Benutzerabläufe testen"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Was in diesem Leitfaden getestet wird

Verwende Playwright mit einem **Server-Build**, um dein Frontend und echte Go-Dienste in einem Browser zu testen. Die Vite-Vorschau an Port 9245 stellt die Frontend-Dateien bereit. Wenn du sie in einem gewöhnlichen Browser öffnest, stehen dadurch allein noch keine Wails-Dienstbindungen zur Verfügung.

Dieses Testsystem deckt die Browser-Benutzeroberfläche und Dienstaufrufe ab. Es automatisiert weder native WebView-Fenster noch Betriebssystemmenüs, Taskleistensymbole oder native Dateidialoge. Teste diese separat auf jeder unterstützten Desktopplattform.

## Eine Starter-Anwendung vorbereiten

Erstelle eine JavaScript-Starter-Anwendung und arbeite anschließend in ihrem Verzeichnis:

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

Füge in `main.go` dem vorhandenen Wert `application.Options` dieses Feld hinzu:

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Behalte den vorhandenen Asset-Handler und die Dienstregistrierung bei. Die Serveroptionen werden wirksam, wenn die Anwendung mit dem Tag `server` gebaut wird. Informationen zu Funktionen und Einschränkungen findest du unter [Server-Build](/guides/server-build/).

Generiere die Bindungen und baue das Produktions-Frontend, bevor du den Test-Host startest:

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Committe das resultierende Paketmanifest und die Lockdatei, damit CI `npm ci` verwenden kann.

## Den Go-Host starten und darauf warten

Erstelle `frontend/playwright.config.js`:

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

Playwright startet die Go-Anwendung aus dem Projektstammverzeichnis und wartet auf ihren HTTP-Endpunkt. Führe pro Testjob einen isolierten Host aus. Verwende temporären Speicher oder eine Testdatenbank, wenn du Persistenztests hinzufügst. Richte dieses Testsystem niemals auf echte Benutzerdaten.

## Einen echten Dienstaufruf überprüfen

Erstelle `frontend/e2e/app.spec.js`:

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

Die obigen Selektoren entsprechen der generierten Starter-Anwendung `vanilla-js`. Passe sie an deine Benutzeroberfläche an, nachdem du die Starter-Anwendung ersetzt hast. Die Ergebnisprüfung wartet auf den Go-Aufruf, statt eine feste Verzögerung anzunehmen.

Führe den Befehl aus `frontend` aus:

```bash
npx playwright test
```

Diese Konfiguration verwendet den standardmäßigen Chromium-Browser von Playwright.

## Browserinteraktionen und native Grenzen

Registriere für einen Browser-`window.confirm` den Handler **vor** dem Klick, der ihn öffnet:

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

Verwende für ein HTML-Upload-Eingabefeld `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. Dadurch wird dem Eingabefeld des Browsers eine Datei bereitgestellt; ein nativer Speicherdialog von Wails wird damit nicht getestet. Teste das Speicherergebnis deiner Anwendung über ihre Dienstschnittstelle und überprüfe den nativen Dialog anschließend manuell oder mit einem geeigneten Automatisierungssystem für die jeweilige Plattform.

Der konfigurierte Viewport von 800×600 testet das responsive Browserlayout. `page.viewportSize()` misst kein natives Wails-Fenster. Ebenso behandeln Playwright-Handler für Browserdialoge keine nativen Meldungsdialoge von Wails.

## CI

Dieser Workflow setzt voraus, dass sich die Anwendung im Stammverzeichnis des Repositorys befindet und dass die oben genannte Playwright-Abhängigkeit sowie die Lockdatei committet sind:

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

Das konfigurierte `webServer` startet das Backend sowohl in CI als auch lokal. Halte Servertests von der nativen Desktop-Paketierung und den Startprüfungen getrennt.

## Nächste Schritte

- [Tests](/guides/testing/) — Unit-Tests und Tests auf Dienstebene
- [Server-Build](/guides/server-build/) — Verhalten des Browser-Hosts
- [Bauen](/guides/build/building/) — native Builds
