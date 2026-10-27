---
title: "Tests de bout en bout"
description: "Testez des parcours utilisateur complets"
slug: "guides/e2e-testing"
sourcePath: "guides/e2e-testing.md"
---

## Ce que teste ce guide

Utilisez Playwright avec une **version serveur** pour tester votre frontend et de véritables services Go dans un navigateur. La prévisualisation Vite sur le port 9245 sert les fichiers du frontend ; son ouverture dans un navigateur ordinaire ne fournit pas à elle seule les liaisons de services Wails.

Ce banc de test couvre l’interface utilisateur dans le navigateur et les appels de services. Il n’automatise ni les fenêtres WebView natives, ni les menus du système d’exploitation, ni les icônes de la zone de notification, ni les boîtes de dialogue natives de sélection de fichiers. Testez-les séparément sur chaque plateforme de bureau prise en charge.

## Préparer une application de démarrage

Créez une application de démarrage JavaScript, puis travaillez dans son répertoire :

```bash
wails3 init -n e2e-app -t vanilla-js
cd e2e-app
```

Dans `main.go`, ajoutez ce champ à la valeur `application.Options` existante :

```go
Server: application.ServerOptions{
    Host: "127.0.0.1",
    Port: 8080,
},
```

Conservez le gestionnaire de ressources et l’enregistrement des services existants. Les options du serveur prennent effet lorsque l’application est compilée avec l’étiquette `server`. Consultez [Version serveur](/guides/server-build/) pour connaître ses possibilités et ses limites.

Générez les liaisons et compilez le frontend de production avant de démarrer l’hôte de test :

```bash
wails3 generate bindings -f '-tags=server'
cd frontend
npm install
npm run build
npm install -D @playwright/test
npx playwright install chromium
```

Validez dans le dépôt le manifeste du paquet et le fichier de verrouillage obtenus afin que l’intégration continue puisse utiliser `npm ci`.

## Démarrer l’hôte Go et attendre qu’il soit prêt

Créez `frontend/playwright.config.js` :

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

Playwright démarre l’application Go depuis la racine du projet et attend que son point de terminaison HTTP soit disponible. Exécutez un hôte isolé par tâche de test. Utilisez un stockage temporaire ou une base de données de test lorsque vous ajoutez des tests de persistance ; ne dirigez jamais ce banc de test vers de vraies données utilisateur.

## Vérifier un véritable appel de service

Créez `frontend/e2e/app.spec.js` :

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

Les sélecteurs ci-dessus correspondent à l’application de démarrage `vanilla-js` générée. Adaptez-les à votre interface utilisateur après avoir remplacé cette application. L’assertion sur le résultat attend l’appel Go au lieu de supposer un délai fixe.

Exécutez la commande depuis `frontend` :

```bash
npx playwright test
```

Cette configuration utilise le navigateur Chromium par défaut de Playwright.

## Interactions dans le navigateur et limites des fonctionnalités natives

Pour une `window.confirm` du navigateur, enregistrez le gestionnaire **avant** le clic qui l’ouvre :

```javascript
page.once('dialog', dialog => dialog.accept())
await page.getByRole('button', { name: 'Delete' }).click()
```

Pour un champ HTML de téléversement, utilisez `page.locator('input[type=file]').setInputFiles('fixtures/note.json')`. Cette méthode fournit un fichier au champ du navigateur ; elle ne teste pas une boîte de dialogue d’enregistrement native de Wails. Testez le résultat de l’enregistrement de votre application via son interface de service, puis vérifiez la boîte de dialogue native manuellement ou à l’aide d’un banc d’automatisation adapté à la plateforme.

La fenêtre d’affichage configurée de 800 × 600 teste la mise en page adaptative dans le navigateur. `page.viewportSize()` ne mesure pas une fenêtre Wails native. De même, les gestionnaires de boîtes de dialogue du navigateur de Playwright ne prennent pas en charge les boîtes de dialogue de message natives de Wails.

## Intégration continue

Ce workflow suppose que l’application se trouve à la racine du dépôt et que la dépendance Playwright ainsi que le fichier de verrouillage mentionnés ci-dessus sont validés dans le dépôt :

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

Le paramètre `webServer` de la configuration démarre le backend aussi bien dans l’intégration continue qu’en local. Séparez les tests du serveur des contrôles de création des paquets natifs pour ordinateur et de leur lancement.

## Étapes suivantes

- [Tests](/guides/testing/) — tests unitaires et au niveau des services
- [Version serveur](/guides/server-build/) — comportement de l’hôte dans le navigateur
- [Compilation](/guides/build/building/) — compilations natives
