---
title: "範本"
description: "Wails 的專案範本與入門套件"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
此頁面對 Wails v3 而言可能已過時。

@end

<!-- TODO: Update this link -->

此頁面列出社群支援的範本。若要建立自己的範本，請參閱[範本](https://v2.wails.io/docs/guides/templates)指南。

@note{type="tip" title="如何提交範本"}
您可以按一下底部的`Edit this page`，將您的範本加入其中。

@end

若要使用這些範本，請執行`wails init -n "Your Project Name" -t [the link below[@version]]`

若未指定版本後綴，預設會使用主分支的程式碼範本。若指定版本後綴，則會使用與該版本標籤對應的程式碼範本。

範例：`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="注意"}
**Wails 專案不維護任何第三方範本，亦不對其負責或承擔任何責任！3rd party templates!**

若您對某個範本有疑慮，請檢查`package.json`和`wails.json`，確認它會執行哪些指令碼及安裝哪些套件。

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue)－以 Vue 生態系為基礎的 Wails 範本（整合 TypeScript、深色主題、國際化、單頁路由及 TailwindCSS）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js)－使用 JavaScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts)－使用 TypeScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、Prettier，以及搭配&lt;script setup&gt;的 Composition API）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive)－以 Naive UI（Vue 3元件庫）為基礎的 Wails 範本
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt)－使用純淨 Nuxt3 與 TypeScript，並支援自動匯入 Wails JS 執行階段的 Wails 範本
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template)－使用 Vue+TypeScript+Vite+Element-plus 的 Wails 範本（仿網易雲）

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular)－功能豐富且可直接投入正式環境的 Angular 15+ 範本。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template)－搭配 TypeScript、Sass、熱重新載入、程式碼分割和 i18n 的 Angular 範本

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template)－使用 ReactJS 的範本
- [wails-react-template](https://github.com/flin7/wails-react-template)－支援即時開發的精簡 React 範本
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs)－使用 Next.js 和 TypeScript 的範本
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router)－使用 Next.js、TypeScript 和 App Router 的範本
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router)－使用 Next.js、TypeScript 和 App Router src，並附帶範例的範本
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template)－React + TypeScript + Vite + TailwindCSS 範本
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts)－搭配 Vite、React、TypeScript、TailwindCSS 和 shadcn/ui 的範本

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template)－使用 Svelte 的範本
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template)－使用 Svelte 和 Vite 的範本
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template)－使用 Svelte、Vite 和 TailwindCSS v3 的範本
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master)－使用 Svelte v4.2.0、Vite 和 TailwindCSS v3.3.3 的新版範本
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template)－使用 SvelteKit 的範本
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte)－使用 SvelteKit 和 Shadcn-Svelte 的範本

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts)－使用 Solid + TS + Vite 的範本
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js)－使用 Solid + JS + Vite 的範本

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template)－透過函數式程式設計及**反應迅速的**熱重新載入設定來開發 GUI 應用程式 :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind)－結合 Elm + Tailwind CSS + Wails 的強大功能 :muscle:！支援熱重新載入。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template)－採用獨特組合：以純 htmx 提供互動功能，並以 templ 建立元件和表單

## 純 JavaScript（Vanilla）

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template)－僅使用基本 JavaScript、HTML 和 CSS 的範本

## Lit（Web Components）

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template)－使用 Lit 和 Shoelace 元件庫提供前端，並預先設定 Prettier 與 TypeScript 的 Wails 範本。
