---
title: テンプレート
description: Wails用のプロジェクトテンプレートとスターターキット
---

:::caution

このページはWails v3では古くなっている可能性があります。

:::

<!-- TODO: Update this link -->

このページはコミュニティがサポートするテンプレートのリストです。独自のテンプレートを作成するには、[テンプレート](https://wails.io/docs/guides/templates)
ガイドを参照してください。

:::tip[テンプレートの提出方法]

下部の `Edit this page` をクリックして、あなたのテンプレートを含めることができます。

:::

これらのテンプレートを使用するには、
`wails init -n "Your Project Name" -t [以下のリンク[@version]]` を実行します。

バージョン接尾辞がない場合、デフォルトでmainブランチのコードテンプレートが使用されます。
バージョン接尾辞がある場合、そのバージョンのタグに対応するコードテンプレートが使用されます。

例:
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

:::danger[注意]

**Wailsプロジェクトは、サードパーティのテンプレートを維持せず、それらに対して責任を負いません！**

テンプレートについて確信が持てない場合は、`package.json` と `wails.json` を確認して、
どのスクリプトが実行され、どのパッケージがインストールされているかを確認してください。

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Vueエコシステムに基づくWails
  テンプレート (TypeScript統合、ダークテーマ、
  国際化、シングルページルーティング、TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  JavaScript + Quasar V2を使用するテンプレート (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  TypeScript + Quasar V2を使用するテンプレート (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, &lt;script setup&gt;付きのComposition API)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Naive UI (Vue 3のコンポーネントライブラリ)に基づくWailsテンプレート
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Wails
  テンプレートはクリーンなNuxt3とTypeScriptを使用し、wails js
  ランタイムの自動インポートに対応
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Wails
  テンプレートはVue+TypeScript+Vite+Element-plus(NetEase Cloud Music風)を使用

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ で機能満載、本番環境への展開準備完了。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  TypeScript、Sass、ホットリロード、コード分割、i18nを備えたAngular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  reactjsを使用するテンプレート
- [wails-react-template](https://github.com/flin7/wails-react-template) - リモート開発に対応した
  React用のミニマルなテンプレート
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Next.jsとTypeScriptを使用する
  テンプレート
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  App routerを備えたNext.jsとTypeScriptを使用するテンプレート
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  App router srcと例を備えたNext.jsとTypeScriptを使用するテンプレート
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  React + TypeScript + Vite + TailwindCSS用のテンプレート
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Vite、React、TypeScript、TailwindCSS、shadcn/uiを備えたテンプレート

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Svelteを使用するテンプレート
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  SvelteとViteを使用するテンプレート
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  TailwindCSS v3を備えたSvelteとViteを使用するテンプレート
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Svelte v4.2.0とVite、TailwindCSS v3.3.3を使用する更新されたテンプレート
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  SvelteKitを使用するテンプレート
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  SveltekitとShadcn-Svelteを使用するテンプレート

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Solid + Ts + Viteを使用するテンプレート
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Solid + Js + Viteを使用するテンプレート

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  関数型プログラミングと**高速な**ホットリロードセットアップでGUIアプリを開発しましょう
  :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Elm + Tailwind CSS + Wailsの力を組み合わせましょう :muscle:! ホットリロードに対応。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  インタラクティブ性のために純粋なhtmxを使用し、コンポーネントとフォームの作成にはtemplを使用する
  独自の組み合わせ

## 純粋なJavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 基本的なJavaScript、HTML、CSS
  のみを使用したテンプレート

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  フロントエンドにlit、Shoelaceコンポーネントライブラリ、
  事前に設定されたprettierとtypescriptを提供するWailsテンプレート。