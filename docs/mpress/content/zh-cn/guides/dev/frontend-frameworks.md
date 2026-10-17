---
title: "使用其他前端框架"
description: "如何通过将自己的 Vite 项目放入 frontend 目录，使用没有内置模板的框架"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails 特意只为少数几种框架提供了内置入门模板：

| 模板 | 语言 |
| --- | --- |
| `vanilla` | TypeScript（默认） |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

但这并不意味着你只能使用这些选项。Wails 应用的前端<strong>只是一个 Web 项目</strong>——任何能够构建为静态 HTML/CSS/JS 的项目都可以使用。如果你选择的框架（Solid、Preact、Lit、Qwik、SvelteKit、Angular……）没有模板，只需几分钟即可自行搭建。

## `frontend/`目录的用途

Wails 不关心`frontend/`中使用哪种框架。它只依赖一套简单且与框架无关的约定：

- **`frontend/dist/`是最终随应用分发的内容。**`main.go`使用`//go:embed all:frontend/dist`嵌入构建后的前端，并通过资源服务器提供这些内容。构建必须将静态包输出到`frontend/dist/`（Vite 的默认输出目录）。
- <strong>构建由`frontend/package.json`驱动。</strong>执行`wails3 build`时，Wails 会运行前端的`build`脚本；执行`wails3 dev`时，它会运行`dev`并代理 Vite 开发服务器，以支持热重载。
- <strong>绑定会生成到`frontend/bindings/`中。</strong>Wails 会检查已注册的 Go 服务，并在该目录中写入类型安全的 SDK。你可以像导入其他模块一样导入它：
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **开发服务器在固定端口上运行。**`wails3 dev`会在`WAILS_VITE_PORT`指定的端口（默认为`9245`）上代理 Vite，因此请将`server.port`设为该端口，同时设置`strictPort: true`。这只是普通的 Vite 配置，不涉及 Wails 插件。
- <strong>可选——带类型的自定义事件。</strong>内置模板还会注册`@wailsio/runtime/plugins/vite`插件。仅当你使用<em>带类型的</em>自定义事件时才需要它：该插件会将生成的事件类型定义注入运行时，并使构建在绑定生成之前失败。如果你只使用基于字符串的`Events.On("time", …)` API，可以不添加该插件。

其他一切——组件、路由、状态和样式——完全由你所用的框架决定。

## 使用 Vite 搭建任意框架

最快的方法是从内置模板开始（这样可以获得`main.go`、`Taskfile`、构建资源以及可正常工作的 Go 服务），然后用适用于所选框架的新 Vite 项目替换`frontend/`。

@steps
### 使用默认模板创建项目
```bash
wails3 init -n myapp
cd myapp
```

### 将`frontend/`替换为适用于所选框架的 Vite 应用
Vite 可以用一条命令搭建大多数框架。请选择一个模板：

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

可以将`solid`替换为任意 Vite 模板：`preact`、`lit`、`svelte`、`vue`、`react`、`vanilla`，或者它们的`-ts`变体（`solid-ts`、`preact-ts`……）。

### 安装运行时并让 Vite 指向 Wails 开发服务器
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime`提供 JS API（`Events`、`Browser`、对话框等）。唯一<em>必需的</em>`vite.config`更改是开发服务器端口，这样`wails3 dev`才能找到它：

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

仅当你打算使用<strong>带类型的自定义事件</strong>时，才需要另外添加该插件——它会注入生成的事件类型，并要求在构建前已存在绑定：

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### 调用 Go 服务
先生成一次绑定，然后就可以在组件中的任意位置导入它们：

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### 运行项目
```bash
wails3 dev
```

@end

@note{type="tip" title="生成纯 JavaScript 项目"}
同一条命令也可以搭建不使用 TypeScript 的项目——只需使用非`-ts` Vite 模板：

```bash
npm create vite@latest frontend -- --template solid
```

@end

## 拥有自有脚手架的框架

少数框架不通过 Vite 的`create`模板创建，而是使用自己的工具。它们仍然可以正常使用——只需用框架原生的命令搭建项目，然后添加 Wails 插件：

- **SvelteKit：**`npx sv create frontend`。使用静态适配器（`@sveltejs/adapter-static`），使其构建为静态包，并禁用 SSR。
- **Qwik：**`npm create qwik@latest`。使用静态（SSG）适配器。
- <strong>Angular：</strong>使用`ng new`搭建项目，将`outputPath`设为`dist`，并让构建脚本指向`ng build`。

规则始终相同：在`frontend/dist/`中生成静态构建产物，保留`@wailsio/runtime` Vite 插件（或直接导入运行时），并从`frontend/bindings/`导入 Go 绑定。

@note{type="info"}
如果你为某个框架制作了一套完善的配置，可以考虑将其发布为[自定义模板](/guides/advanced/custom-templates/)，以便其他人直接`wails3 init -t`它。

@end
