---
title: "範本"
description: "Wails 的專案範本與入門套件"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
對 Wails v3而言，此頁面的內容可能已過時。

@end

<!-- TODO: Update this link -->

此頁面列出由社群支援的範本。若要建立自己的範本，請參閱[範本](https://wails.io/docs/guides/templates)指南。

@note{type="tip" title="如何提交範本"}
您可以按一下底部的`Edit this page`，以加入您的範本。

@end

若要使用這些範本，請執行  
`wails init -n "Your Project Name" -t [the link below[@version]]`

若沒有版本後綴，預設會使用主分支的程式碼範本。  
若有版本後綴，則會使用與該版本標籤對應的程式碼範本。

範例：  
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="注意"}
**Wails 專案不維護第3方範本， 亦不對其負責或承擔任何法律責任！**

若您對某個範本有疑慮，請檢查`package.json`和`wails.json`，確認 會執行哪些指令碼以及安裝哪些套件。

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue)－以 Vue 生態系為基礎的 Wails 範本（整合 TypeScript、深色主題、國際化、單頁路由及 TailwindCSS）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js)－使用 JavaScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts)－使用 TypeScript + Quasar V2 的範本（Vue 3、Vite、Sass、Pinia、ESLint、Prettier，以及搭配&lt;script setup&gt;的 Composition API）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive)－以 Naive UI（Vue 3元件庫）為基礎的 Wails 範本
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt)－使用純淨 Nuxt3 與 TypeScript，並可自動匯入 Wails JS 執行階段的 Wails 範本
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template)－使用 Vue+TypeScript+Vite+Element-plus 的 Wails 範本（仿網易雲）

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular)－功能豐富的 Angular 15+ 範本，已可投入正式環境。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template)－整合 TypeScript、Sass、熱重新載入、程式碼分割及 i18n 的 Angular 範本

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template)－使用 ReactJS 的範本
- [wails-react-template](https://github.com/flin7/wails-react-template)－支援即時開發的精簡 React 範本
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs)－使用 Next.js 與 TypeScript 的範本
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router)－使用 Next.js、TypeScript 與 App Router 的範本
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router)－使用 Next.js、TypeScript 與 App Router src，並附帶範例的範本
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template)－適用於 React + TypeScript + Vite + TailwindCSS 的範本
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts)－整合 Vite、React、TypeScript、TailwindCSS 及 shadcn/ui 的範本

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template)－使用 Svelte 的範本
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template)－使用 Svelte 與 Vite 的範本
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template)－使用 Svelte、Vite 與 TailwindCSS v3 的範本
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master)－已更新的範本，使用 Svelte v4.2.0、Vite 與 TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template)－使用 SvelteKit 的範本
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte)－使用 SvelteKit 與 Shadcn-Svelte 的範本

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts)－使用 Solid + TypeScript + Vite 的範本
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js)－使用 Solid + JavaScript + Vite 的範本

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template)－使用函數式程式設計與<strong>反應迅速的</strong>熱重新載入設定來開發 GUI 應用程式:tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind)－結合 Elm + Tailwind CSS + Wails 的強大能力 :muscle:！支援熱重新載入。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template)－採用獨特組合：以純 htmx 提供互動功能，並以 templ 建立元件與表單

## 純 JavaScript（Vanilla）

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template)－僅包含基本 JavaScript、HTML 與 CSS 的範本

## Lit（Web 元件）

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template)－以前端 Lit、Shoelace 元件庫及預先設定的 Prettier 和 TypeScript 為基礎的 Wails 範本
