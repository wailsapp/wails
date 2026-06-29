---
title: 範本
description: Wails 的專案範本與起始套件
---

:::caution

此頁面可能已過時，不適用於 Wails v3。

:::

<!-- TODO: Update this link -->

此頁面提供社群支援的範本列表。若要建立您自己的
範本，請參閱 [Templates](https://wails.io/docs/guides/templates)
指南。

:::tip[如何提交範本]

您可以點擊底部的 `Edit this page` 來新增您的範本。

:::

要使用這些範本，請執行
`wails init -n "Your Project Name" -t [下方的連結[@version]]`

如果沒有版本後綴，則預設使用 main 分支的程式碼範本。
如果有版本後綴，則使用對應此版本標籤的程式碼範本。

範例：
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

:::danger[注意]

**Wails 專案不維護、不負責亦不對第三方範本承擔任何責任！**

如果您對某個範本不確定，請檢查 `package.json` 和 `wails.json`，
以了解執行了哪些腳本以及安裝了哪些套件。

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - 基於 Vue 生態系的 Wails
  範本（整合 TypeScript、深色主題、
  國際化、單頁路由、TailwindCSS）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  使用 JavaScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、
  Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  使用 TypeScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、
  Prettier、Composition API 搭配 &lt;script setup&gt;）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  基於 Naive UI（一個 Vue 3 元件庫）的 Wails 範本
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - 使用乾淨的 Nuxt3 和 TypeScript 的 Wails
  範本，並為 wails js
  runtime 提供自動導入功能
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - 使用 Vue+TypeScript+Vite+Element-plus(仿网易云) 的 Wails
  範本

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ 功能豐富且已準備好投入生產環境。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  使用 TypeScript、Sass、熱重載、程式碼分割和 i18n 的 Angular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  使用 reactjs 的範本
- [wails-react-template](https://github.com/flin7/wails-react-template) - 支援即時開發的
  React 最小化範本
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - 使用 Next.js 和 TypeScript 的
  範本
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  使用 Next.js 和 TypeScript 並搭配 App router 的
  範本
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  使用 Next.js 和 TypeScript 並搭配 App router src + 範例的
  範本
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  用於 React + TypeScript + Vite + TailwindCSS 的
  範本
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  使用 Vite、React、TypeScript、TailwindCSS 和 shadcn/ui 的
  範本

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  使用 Svelte 的範本
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  使用 Svelte 和 Vite 的範本
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  使用 Svelte 和 Vite 並搭配 TailwindCSS v3 的
  範本
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  使用 Svelte v4.2.0 和 Vite 並搭配 TailwindCSS v3.3.3 的更新版範本
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  使用 SvelteKit 的範本
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  使用 Sveltekit 和 Shadcn-Svelte 的範本

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  使用 Solid + Ts + Vite 的範本
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  使用 Solid + Js + Vite 的範本

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  使用函數式程式設計和**快速**熱重載
  設定來開發您的 GUI 應用程式 :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  結合 Elm + Tailwind CSS + Wails 的力量 :muscle:！支援熱重載。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  使用純 htmx 的獨特組合來實現互動性，並使用 templ 來
  建立元件和表單

## 純 JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 僅包含基本 JavaScript、HTML 和 CSS 的
  範本

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  Wails 範本，為前端提供 lit、Shoelace 元件庫 +
  預先設定的 prettier 和 typescript。