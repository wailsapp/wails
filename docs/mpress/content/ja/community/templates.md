---
title: "テンプレート"
description: "Wails向けのプロジェクトテンプレートとスターターキット"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
このページはWails v3について最新ではない可能性があります。

@end

<!-- TODO: Update this link -->

このページでは、コミュニティがサポートするテンプレートを一覧にまとめています。独自のテンプレートを作成するには、[テンプレート](https://v2.wails.io/docs/guides/templates)ガイドを参照してください。

@note{type="tip" title="テンプレートを登録する方法"}
独自のテンプレートを掲載するには、ページ下部の`Edit this page`をクリックしてください。

@end

これらのテンプレートを使用するには、`wails init -n "Your Project Name" -t [the link below[@version]]`を実行します。

バージョンサフィックスを指定しない場合、デフォルトでmainブランチのコードテンプレートが使用されます。バージョンサフィックスを指定した場合、そのバージョンのタグに対応するコードテンプレートが使用されます。

例：`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="注意"}
**Wailsプロジェクトは、サードパーティー製テンプレートの保守を行っておらず、それらに対して責任を負いません。3**

テンプレートに不安がある場合は、実行されるスクリプトとインストールされるパッケージを`package.json`と`wails.json`で確認してください。

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - VueエコシステムをベースにしたWailsテンプレート（TypeScript、ダークテーマ、国際化、シングルページルーティング、TailwindCSSを統合）
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) - JavaScript + Quasar V2を使用するテンプレート（Vue 3、Vite、Sass、Pinia、ESLint、Prettier）
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) - TypeScript + Quasar V2を使用するテンプレート（Vue 3、Vite、Sass、Pinia、ESLint、Prettier、&lt;script setup&gt;を使用するComposition API）
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) - Naive UI（Vue 3向けコンポーネントライブラリ）をベースにしたWailsテンプレート
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - クリーンなNuxt3とTypeScriptを使用し、Wails JavaScriptランタイムの自動インポートに対応したWailsテンプレート
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Vue+TypeScript+Vite+Element-plusを使用するWailsテンプレート（NetEase Cloud Music風）

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) - Angular 15+ 機能が充実し、すぐに本番環境へ導入可能。
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) - TypeScript、Sass、ホットリロード、コード分割、i18nに対応したAngular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) - reactjsを使用するテンプレート
- [wails-react-template](https://github.com/flin7/wails-react-template) - ライブ開発に対応したReact向けの最小構成テンプレート
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Next.jsとTypeScriptを使用するテンプレート
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) - Next.jsとTypeScriptをApp Routerとともに使用するテンプレート
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) - Next.jsとTypeScriptを使用し、App Routerのsrc構成とサンプルを含むテンプレート
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) - React + TypeScript + Vite + TailwindCSS向けテンプレート
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) - Vite、React、TypeScript、TailwindCSS、shadcn/uiを使用するテンプレート

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) - Svelteを使用するテンプレート
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) - SvelteとViteを使用するテンプレート
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) - SvelteとViteをTailwindCSS v3とともに使用するテンプレート
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) - Svelte v4.2.0とViteをTailwindCSS v3.3.3とともに使用する更新版テンプレート
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) - SvelteKitを使用するテンプレート
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) - SvelteKitとShadcn-Svelteを使用するテンプレート

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) - Solid + Ts + Viteを使用するテンプレート
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) - Solid + Js + Viteを使用するテンプレート

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) - 関数型プログラミングと**軽快な**ホットリロード環境でGUIアプリを開発できます :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) - Elm + Tailwind CSS + Wailsの力を組み合わせます :muscle:！ホットリロードに対応しています。

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) - インタラクティブ機能には純粋なhtmxを、コンポーネントとフォームの作成にはtemplを使う独自の組み合わせ

## 純粋なJavaScript（Vanilla）

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 基本的なJavaScript、HTML、CSSだけで構成されたテンプレート

## Lit（Web Components）

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) - litとShoelaceコンポーネントライブラリを使用し、prettierとtypescriptが事前設定されたフロントエンドを提供するWailsテンプレート
