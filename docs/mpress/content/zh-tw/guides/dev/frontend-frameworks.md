---
title: "使用其他前端框架"
description: "如何將您自己的 Vite 專案放入 frontend 目錄，以使用沒有內建範本的框架"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails 刻意只為少數幾種框架內建入門範本：

| 範本 | 語言 |
| --- | --- |
| `vanilla` | TypeScript（預設） |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

但這不代表您只能選擇這些框架。Wails 應用程式的前端<strong>就只是一個 Web 專案</strong>——任何能建置為靜態 HTML/CSS/JS 的框架都能使用。如果您選擇的框架（Solid、Preact、Lit、Qwik、SvelteKit、Angular……）沒有範本，只要幾分鐘就能自行建立初始專案。

## `frontend/`目錄的用途

Wails 不在意`frontend/`中使用哪一種框架。它只依賴一組精簡且與框架無關的約定：

- **`frontend/dist/`是最終隨應用程式發布的內容。**`main.go`會使用`//go:embed all:frontend/dist`嵌入建置完成的前端，並由資產伺服器提供。您的建置必須將靜態套件輸出至`frontend/dist/`（Vite 的預設輸出目錄）。
- <strong>建置由`frontend/package.json`驅動。</strong>執行`wails3 build`時，Wails 會執行前端的`build`指令碼；執行`wails3 dev`時，則會執行`dev`並代理 Vite 開發伺服器，以進行熱重新載入。
- <strong>繫結會產生至`frontend/bindings/`。</strong>Wails 會檢查您註冊的 Go 服務，並在該處寫入型別安全的 SDK。您可以像匯入其他模組一樣匯入它：
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **開發伺服器在固定連接埠上執行。**`wails3 dev`會在`WAILS_VITE_PORT`所指定的連接埠（預設為`9245`）上代理 Vite，因此請將`server.port`設為該連接埠，並同時設定`strictPort: true`。這只是一般的 Vite 設定，不涉及任何 Wails 外掛程式。
- <strong>選用功能——具型別的自訂事件。</strong>內建範本也會註冊`@wailsio/runtime/plugins/vite`外掛程式。只有在使用<em>具型別的</em>自訂事件時才需要它：它會將產生的事件型別定義注入執行階段，而且在繫結產生前會使建置失敗。如果您只使用以字串為基礎的`Events.On("time", …)` API，可以省略它。

其餘一切——元件、路由、狀態與樣式——完全由您的框架決定。

## 使用 Vite 建立任何框架的初始專案

最快的方法是從內建範本開始（如此即可取得`main.go`、`Taskfile`、建置資產及可運作的 Go 服務），然後將`frontend/`替換成您框架的全新 Vite 專案。

@steps
### 從預設範本建立專案
```bash
wails3 init -n myapp
cd myapp
```

### 將 `frontend/` 替換成您框架的 Vite 應用程式
Vite 只需一個命令，就能為大多數框架建立初始專案。請選擇一個範本：

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

可將`solid`替換成任何 Vite 範本：`preact`、`lit`、`svelte`、`vue`、`react`、`vanilla`，或它們的`-ts`變體（`solid-ts`、`preact-ts`……）。

### 安裝執行階段，並將 Vite 指向 Wails 開發伺服器
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime`提供 JS API（`Events`、`Browser`、對話方塊等）。唯一<em>必要的</em>`vite.config`變更是開發伺服器的連接埠，如此`wails3 dev`才能找到它：

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

只有在您打算使用<strong>具型別的自訂事件</strong>時，才需另外加入此外掛程式——它會注入產生的事件型別，並要求建置前必須已有繫結：

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### 呼叫您的 Go 服務
先產生一次繫結，之後便可在元件中的任何位置匯入：

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### 執行專案
```bash
wails3 dev
```

@end

@note{type="tip" title="產生純 JavaScript 專案"}
同一個命令也能建立不使用 TypeScript 的專案——只需使用非`-ts`的 Vite 範本：

```bash
npm create vite@latest frontend -- --template solid
```

@end

## 具備自有初始專案建立工具的框架

少數框架不是透過 Vite 的`create`範本建立，而是有自己的工具。它們仍然可以使用——只要以其原生命令建立初始專案，再加入 Wails 外掛程式即可：

- **SvelteKit：**`npx sv create frontend`。請使用靜態配接器（`@sveltejs/adapter-static`），使其建置為靜態套件，並停用 SSR。
- **Qwik：**`npm create qwik@latest`。請使用靜態（SSG）配接器。
- <strong>Angular：</strong>使用`ng new`建立初始專案，將`outputPath`設為`dist`，並讓建置指令碼指向`ng build`。

規則一律相同：在`frontend/dist/`中產生靜態建置，保留`@wailsio/runtime` Vite 外掛程式（或直接匯入執行階段），並從`frontend/bindings/`匯入您的 Go 繫結。

@note{type="info"}
如果您為某個框架打造出完善的設定，請考慮將其發布為[自訂範本](/guides/advanced/custom-templates/)，讓其他人可以直接`wails3 init -t`它。

@end
