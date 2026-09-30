---
title: "模板"
description: "Wails 项目模板和入门套件"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
此页面对于 Wails v3 可能已过时。

@end

<!-- TODO: Update this link -->

此页面列出了由社区支持的模板。要构建自己的 模板，请参阅[模板](https://wails.io/docs/guides/templates) 指南。

@note{type="tip" title="如何提交模板"}
你可以点击底部的`Edit this page`来添加你的模板。

@end

要使用这些模板，请运行 `wails init -n "Your Project Name" -t [the link below[@version]]`

如果没有版本后缀，则默认使用主分支中的代码模板。 如果有版本后缀，则使用与该版本标签对应的代码模板。

示例： `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="注意"}
**Wails项目不维护第3方模板，也不对其负责或承担法律责任！**

如果你对某个模板不确定，请检查`package.json`和`wails.json`，了解会运行哪些脚本以及会安装哪些软件包。

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - 基于 Vue 生态的 Wails 模板（集成 TypeScript、深色主题、国际化、单页路由和 TailwindCSS）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) - 使用 JavaScript + Quasar V2 的模板（Vue 3、Vite、Sass、Pinia、ESLint、Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) - 使用 TypeScript + Quasar V2 的模板（Vue 3、Vite、Sass、Pinia、ESLint、Prettier，以及采用&lt;script setup&gt;的组合式 API）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) - 基于 Naive UI（一个 Vue 3组件库）的 Wails 模板
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - 使用纯净 Nuxt3 和 TypeScript，并可自动导入 Wails JS 运行时的 Wails 模板
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - 使用 Vue+TypeScript+Vite+Element-plus 的 Wails 模板（仿网易云）

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) - Angular 15+ 功能丰富，可直接投入生产。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) - 集成 TypeScript、Sass、热重载、代码拆分和国际化的 Angular 模板

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) - 使用 ReactJS 的模板
- [wails-react-template](https://github.com/flin7/wails-react-template) - 支持实时开发的精简 React 模板
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - 使用 Next.js 和 TypeScript 的模板
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) - 使用 Next.js、TypeScript 和 App Router 的模板
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) - 使用 Next.js、TypeScript、App Router src 并附带示例的模板
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) - 使用 React + TypeScript + Vite + TailwindCSS 的模板
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) - 使用 Vite、React、TypeScript、TailwindCSS 和 shadcn/ui 的模板

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) - 使用 Svelte 的模板
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) - 使用 Svelte 和 Vite 的模板
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) - 使用 Svelte、Vite 和 TailwindCSS v3 的模板
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) - 更新后的模板，使用 Svelte v4.2.0、Vite 和 TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) - 使用 SvelteKit 的模板
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) - 使用 SvelteKit 和 Shadcn-Svelte 的模板

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) - 使用 Solid + TypeScript + Vite 的模板
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) - 使用 Solid + JavaScript + Vite 的模板

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) - 使用函数式编程和<strong>响应迅速的</strong>热重载配置开发 GUI 应用:tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) - 汇聚 Elm + Tailwind CSS + Wails 的强大能力 :muscle:！支持热重载。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) - 采用独特组合：使用纯 htmx 实现交互，并使用 templ 创建组件和表单

## 纯 JavaScript（原生）

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 仅包含基础 JavaScript、HTML 和 CSS 的模板

## Lit（Web 组件）

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) - 提供基于 Lit 和 Shoelace 组件库的前端，并预先配置 Prettier 和 TypeScript 的 Wails 模板
