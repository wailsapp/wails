---
title: "Como usar outros frameworks de frontend"
description: "Como usar um framework que não tem um template integrado colocando seu próprio projeto Vite no diretório frontend"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

O Wails inclui templates iniciais integrados para um conjunto deliberadamente pequeno de frameworks:

| Template | Linguagem |
| --- | --- |
| `vanilla` | TypeScript (padrão) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

Isso não significa que essas sejam suas únicas opções. O frontend de um aplicativo Wails é **apenas um projeto web** — qualquer projeto que gere HTML/CSS/JS estático funcionará. Se o framework de sua preferência (Solid, Preact, Lit, Qwik, SvelteKit, Angular, …) não tiver um template, você poderá criar a estrutura inicial por conta própria em poucos minutos.

## Como o diretório `frontend/` é usado

Para o Wails, não importa qual framework esteja em `frontend/`. Ele depende apenas de um contrato pequeno e independente de framework:

- **`frontend/dist/` é o que será distribuído.** `main.go` incorpora o frontend compilado com `//go:embed all:frontend/dist` e o disponibiliza pelo servidor de ativos. Sua compilação deve gerar um pacote estático em `frontend/dist/` (o diretório de saída padrão do Vite).
- **A compilação é controlada por `frontend/package.json`.** Durante `wails3 build`, o Wails executa o script `build` do frontend; durante `wails3 dev`, ele executa `dev` e encaminha o servidor de desenvolvimento do Vite por proxy para permitir o recarregamento automático.
- **Os bindings são gerados em `frontend/bindings/`.** O Wails inspeciona os serviços Go registrados e grava nesse local um SDK com tipagem segura. Você pode importá-lo como qualquer outro módulo:
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **O servidor de desenvolvimento é executado em uma porta fixa.** `wails3 dev` encaminha o Vite por proxy na porta definida em `WAILS_VITE_PORT` (por padrão, `9245`); portanto, configure `server.port` com essa porta usando `strictPort: true`. Essa é uma configuração comum do Vite — nenhum plugin do Wails está envolvido.
- **Opcional — eventos personalizados tipados.** Os templates integrados também registram o plugin `@wailsio/runtime/plugins/vite`. Ele só é necessário se você usar eventos personalizados *tipados*: o plugin injeta no runtime as definições de tipos de eventos geradas e faz a compilação falhar enquanto os bindings não tiverem sido gerados. Se você usar somente a API `Events.On("time", …)` baseada em strings, poderá omiti-lo.

Todo o restante — componentes, roteamento, estado e estilos — fica inteiramente a cargo do seu framework.

## Crie a estrutura inicial de qualquer framework com o Vite

O caminho mais rápido é começar com um template integrado (para obter `main.go`, o `Taskfile`, os ativos de compilação e um serviço Go funcional) e depois substituir `frontend/` por um novo projeto Vite para seu framework.

@steps
### Crie um projeto usando o template padrão
```bash
wails3 init -n myapp
cd myapp
```

### Substitua `frontend/` por um aplicativo Vite para seu framework
O Vite pode criar a estrutura inicial da maioria dos frameworks com um único comando. Escolha um template:

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

Substitua `solid` por qualquer template do Vite: `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla` ou suas variantes `-ts` (`solid-ts`, `preact-ts`, …).

### Instale o runtime e direcione o Vite ao servidor de desenvolvimento do Wails
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` fornece as APIs JS (`Events`, `Browser`, caixas de diálogo, …). A única alteração *obrigatória* em `vite.config` é a porta do servidor de desenvolvimento, para que `wails3 dev` possa encontrá-lo:

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

Somente se você pretende usar **eventos personalizados tipados**, adicione também o plugin — ele injeta os tipos de eventos gerados e exige que os bindings existam antes da compilação:

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Chame seus serviços Go
Gere os bindings uma vez e depois importe-os em qualquer lugar dos seus componentes:

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### Execute o projeto
```bash
wails3 dev
```

@end

@note{type="tip" title="Gere um projeto JavaScript puro"}
O mesmo comando cria a estrutura inicial de um projeto sem TypeScript — basta usar um template do Vite que não seja `-ts`:

```bash
npm create vite@latest frontend -- --template solid
```

@end

## Frameworks com sua própria ferramenta de criação de estrutura inicial

Alguns frameworks não são criados por meio dos templates `create` do Vite e têm suas próprias ferramentas. Eles também funcionam — basta criar a estrutura inicial com o comando nativo correspondente e depois adicionar o plugin do Wails:

- **SvelteKit:** `npx sv create frontend`. Use o adaptador estático (`@sveltejs/adapter-static`) para gerar um pacote estático. Desative também a SSR.
- **Qwik:** `npm create qwik@latest`. Use o adaptador estático (SSG).
- **Angular:** crie a estrutura inicial com `ng new`, defina `outputPath` como `dist` e direcione o script de compilação para `ng build`.

A regra é sempre a mesma: gere uma compilação estática em `frontend/dist/`, mantenha o plugin `@wailsio/runtime` do Vite (ou importe o runtime diretamente) e importe seus bindings Go de `frontend/bindings/`.

@note{type="info"}
Se você criar uma configuração bem-acabada para um framework, considere publicá-la como um [template personalizado](/guides/advanced/custom-templates/) para que outras pessoas possam executá-la diretamente com `wails3 init -t`.

@end
