---
title: "テンプレート"
description: "Wails 用のプロジェクトテンプレートとスターターキット"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
このページは Wails v3 に関して古い可能性があります。

@end

<!-- TODO: Update this link -->

このページでは、コミュニティによってサポートされているテンプレートを一覧で紹介します。独自の テンプレートを作成する方法については、[テンプレート](https://wails.io/docs/guides/templates) ガイドを参照してください。

@note{type="tip" title="テンプレートを登録する方法"}
下部にある `Edit this page` をクリックすると、作成したテンプレートを追加できます。

@end

これらのテンプレートを使用するには、次を実行します。 `wails init -n "Your Project Name" -t [the link below[@version]]`

バージョンサフィックスがない場合は、デフォルトでメインブランチのコードテンプレートが使用されます。 バージョンサフィックスがある場合は、そのバージョンのタグに対応する コードテンプレートが使用されます。

例： `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="注意"}
**Wails プロジェクトは、第3者製テンプレートの保守を行わず、 これらに関する責任も法的責任も負いません！**

テンプレートについて不明な点がある場合は、どのスクリプトが実行され、どのパッケージが インストールされるかを `package.json` と `wails.json` で確認してください。

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Vue エコシステムをベースとする Wails テンプレート（TypeScript、ダークテーマ、国際化、シングルページルーティング、TailwindCSS を統合）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) - JavaScript + Quasar V2 を使用するテンプレート（Vue 3、Vite、Sass、Pinia、ESLint、Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) - TypeScript + Quasar V2 を使用するテンプレート（Vue 3、Vite、Sass、Pinia、ESLint、Prettier、&lt;script setup&gt; を使用した Composition API）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) - Naive UI（Vue 3 コンポーネントライブラリ）をベースとする Wails テンプレート
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - シンプルな Nuxt3 と TypeScript を使用し、wails js runtime の自動インポートに対応した Wails テンプレート
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Vue+TypeScript+Vite+Element-plus を使用する Wails テンプレート（网易云風）

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) - 機能が充実し、すぐに本番環境へ移行できる Angular 15 以降。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) - TypeScript、Sass、ホットリロード、コード分割、i18n に対応した Angular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) - reactjs を使用するテンプレート
- [wails-react-template](https://github.com/flin7/wails-react-template) - ライブ開発に対応した最小構成の React テンプレート
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Next.js と TypeScript を使用するテンプレート
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) - Next.js、TypeScript、App router を使用するテンプレート
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) - Next.js、TypeScript、App router src を使用し、サンプルを含むテンプレート
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) - React + TypeScript + Vite + TailwindCSS 用テンプレート
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) - Vite、React、TypeScript、TailwindCSS、shadcn/ui を使用するテンプレート

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) - Svelte を使用するテンプレート
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) - Svelte と Vite を使用するテンプレート
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) - Svelte、Vite、TailwindCSS v3 を使用するテンプレート
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) - Svelte v4.2.0、Vite、TailwindCSS v3.3.3 を使用する更新版テンプレート
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) - SvelteKit を使用するテンプレート
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) - Sveltekit と Shadcn-Svelte を使用するテンプレート

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) - Solid + Ts + Vite を使用するテンプレート
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) - Solid + Js + Vite を使用するテンプレート

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) - 関数型プログラミングと<strong>軽快な</strong>ホットリロード環境を使用して GUI アプリを開発できます :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) - Elm + Tailwind CSS + Wails の力を結集 :muscle:！ホットリロードに対応しています。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) - インタラクティブ機能には純粋な htmx を使用し、コンポーネントとフォームの作成には templ を使用する独自の組み合わせ

## ピュア JavaScript（Vanilla）

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 基本的な JavaScript、HTML、CSS だけで構成されたテンプレート

## Lit（Web Components）

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) - lit と Shoelace コンポーネントライブラリを使用し、prettier と typescript が事前設定されたフロントエンドを提供する Wails テンプレート
