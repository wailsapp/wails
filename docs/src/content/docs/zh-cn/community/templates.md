---
title: 模板
description: Wails 的项目模板和入门套件
---

:::caution

此页面可能已过时，不适用于 Wails v3。

:::

<!-- TODO: Update this link -->

本页面列出了由社区支持的模板。若要构建您自己的模板，请参阅 [模板](https://wails.io/docs/guides/templates)
指南。

:::tip[如何提交模板]

您可以点击底部的 `编辑此页面` 来添加您的模板。

:::

要使用这些模板，请运行
`wails init -n "Your Project Name" -t [下方的链接[@version]]`

如果没有版本后缀，则默认使用主分支的代码模板。
如果有版本后缀，则使用该版本标签对应的代码模板。

示例：
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

:::danger[注意]

**Wails 项目不维护、也不对第三方模板负责或承担责任！**

如果您不确定某个模板，请检查 `package.json` 和 `wails.json` 以了解运行了哪些脚本以及安装了哪些包。

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - 基于 Vue 生态的 Wails
  模板（集成 TypeScript、深色主题、
  国际化、单页路由、TailwindCSS）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  使用 JavaScript + Quasar V2 的模板（Vue 3、Vite、Sass、Pinia、ESLint、
  Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  使用 TypeScript + Quasar V2 的模板（Vue 3、Vite、Sass、Pinia、ESLint、
  Prettier、Composition API 与 &lt;script setup&gt;）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  基于 Naive UI（一个 Vue 3 组件库）的 Wails 模板
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - 使用简洁的 Nuxt3 和 TypeScript 的 Wails
  模板，并为 wails js
  运行时提供自动导入功能
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - 使用 Vue+TypeScript+Vite+Element-plus(仿网易云) 的 Wails
  模板

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ 功能丰富且已准备好投入生产。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  带有 TypeScript、Sass、热重载、代码分割和本地化的 Angular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  使用 reactjs 的模板
- [wails-react-template](https://github.com/flin7/wails-react-template) - 支持实时开发的
  最小化 React 模板
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - 使用 Next.js 和 TypeScript 的
  模板
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  使用带有 App 路由器的 Next.js 和 TypeScript 的
  模板
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  使用带有 App 路由器和 src 目录的 Next.js 和 TypeScript 的
  模板 + 示例
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  用于 React + TypeScript + Vite + TailwindCSS 的
  模板
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  带有 Vite、React、TypeScript、TailwindCSS 和 shadcn/ui 的
  模板

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  使用 Svelte 的模板
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  使用 Svelte 和 Vite 的模板
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  使用带有 TailwindCSS v3 的 Svelte 和 Vite 的
  模板
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  使用 Svelte v4.2.0 和 Vite 并带有 TailwindCSS v3.3.3 的
  更新版模板
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  使用 SvelteKit 的模板
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  使用 Sveltekit 和 Shadcn-Svelte 的模板

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  使用 Solid + Ts + Vite 的模板
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  使用 Solid + Js + Vite 的模板

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  使用函数式编程和**快速**的热重载设置来开发您的 GUI 应用
  :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  结合 Elm + Tailwind CSS + Wails 的力量 :muscle:！支持热重载。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  使用纯 htmx 实现交互性，并结合 templ 来创建组件和表单的独特组合

## 纯 JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 一个仅包含基本 JavaScript、HTML 和 CSS 的
  模板

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  为前端提供 lit、Shoelace 组件库 +
  预配置的 prettier 和 typescript 的 Wails 模板。